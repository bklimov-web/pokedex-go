package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next string
	Previous string
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
		},
		"help": {
				name:        "help",
				description: "Displays a help message",
				callback:    commandHelp,
		},
		"map": {
				name:        "map",
				description: "Displays a list of locations",
				callback:    commandMap,
		},
		"mapb": {
				name:        "mapb",
				description: "Displays a list of locations backwards",
				callback:    commandMapBack,
		},
	}
}


func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	config := config{}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		commandName := strings.ToLower(strings.Split(scanner.Text(), " ")[0])
		
		command, exists := getCommands()[commandName]

		if exists {
			err := command.callback(&config)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func commandMap(config *config) error {
	return mapLocationAreas(config.Next, config)
}

func commandMapBack(config *config) error {
	return mapLocationAreas(config.Previous, config)
}

func commandExit(_ *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")

	for _, cmd := range getCommands() {
		fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func cleanInput(text string) []string {
	result := []string{}

	trimmed := strings.Trim(text, " ")
	
	for _, word := range strings.Split(trimmed, " ") {
		if word != "" {
			result = append(result, strings.ToLower(word))
		}
	}

	return result
}