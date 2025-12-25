package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LocationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

const defaultUrl = "https://pokeapi.co/api/v2/location-area"

func getLocationAreas(url string) (LocationAreaResponse, error) {
	if url == "" {
		url = defaultUrl
	}

	resp, err := http.Get(url)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	defer resp.Body.Close()

	var locationAreaData LocationAreaResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&locationAreaData)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	return locationAreaData, nil
}

func mapLocationAreas(url string, config *config) error {
	// cachedResp, exists := 

	resp, err := getLocationAreas(url)
	if err != nil {
		return err
	}

	config.Next = resp.Next
	config.Previous = resp.Previous

	for _, location := range resp.Results {
		fmt.Println(location.Name)
	}

	return nil
}