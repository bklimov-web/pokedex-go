package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

type LocationAreaResponse struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}


func (c *Client) ListLocations(pageUrl *string) (LocationAreaResponse, error) {
	url := baseUrl + locationAreaPath
	
	if pageUrl != nil {
		url = *pageUrl
	}

	if cachedData, exists := c.cache.Get(url); exists {
		locationAreaData := LocationAreaResponse{}

		err := json.Unmarshal(cachedData, &locationAreaData)
		if err != nil {
			return LocationAreaResponse{}, err
		}

		return locationAreaData, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	locationAreaData := LocationAreaResponse{}

	c.cache.Add(url, data)
	
	err = json.Unmarshal(data, &locationAreaData)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	return locationAreaData, nil
}

