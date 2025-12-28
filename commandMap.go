package main

import "fmt"

func commandMap(config *config, args ...string) error {
	return mapLocationAreas(config.Next, config)
}

func commandMapBack(config *config, args ...string) error {
	return mapLocationAreas(config.Previous, config)
}

func mapLocationAreas(url *string, config *config) error {
	resp, err := config.pokeapiClient.ListLocations(url)
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