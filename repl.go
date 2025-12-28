package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"

	"github.com/bklimov-web/pokedex-go/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

type config struct {
	caughtPokemon map[string]pokeapi.PokemonResponse
	pokeapiClient pokeapi.Client
	Next *string
	Previous *string
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
		"explore": {
				name:        "explore <location-name>",
				description: "Explore a location",
				callback:    commandExplore,
		},
		"catch": {
				name:        "catch <pokemon-name>",
				description: "Attempt to catch a pokemon",
				callback:    commandCatch,
		},
		"inspect": {
				name:        "inspect <pokemon-name>",
				description: "Inspect a caught pokemon",
				callback:    commandInspect,
		},
		"pokedex": {
				name:        "pokedex",
				description: "Check out your pokedex",
				callback:    commandPokedex,
		},
	}
}


func startRepl(config *config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}
		
		command, exists := getCommands()[words[0]]

		args := []string{}

		if len(words) > 1 {
			args = words[1:]
		}

		if exists {
			err := command.callback(config, args...)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}


func commandExit(_ *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range getCommands() {
		fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandExplore(config *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("You must provide a location name")
	}

	locationName := args[0]

	location, err	:= config.pokeapiClient.ListPokemonsInLocation(locationName)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", location.Name)
	fmt.Println("Found pokemon: ")

	for _, encounter := range location.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("You must provide a pokemon name")
	}

	name := args[0]

	pokemon, err	:= config.pokeapiClient.GetPokemon(name)
	if err != nil {
		return err
	}

	difficulty := 0.01
	chanceToCatch := math.Exp(-difficulty * float64(pokemon.BaseExperience))
	isCaught := rand.Float64() < chanceToCatch
	fmt.Println(isCaught, chanceToCatch, "chance")

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	if isCaught {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		fmt.Println("You may now inspect it with the inspect command.")
		config.caughtPokemon[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandPokedex(config *config, args ...string) error {
	if len(config.caughtPokemon) == 0 {
		return errors.New("You haven't caught any pokemon yet")
	}

	fmt.Println("Your pokedex:")
	for _, pokemon := range config.caughtPokemon {
		fmt.Printf("  - %s\n", pokemon.Name)
	}

	return nil
}

func commandInspect(config *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("You must provide a pokemon name")
	}

	name := args[0]
	if pokemon, exists := config.caughtPokemon[name]; exists {
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf(" -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}

		fmt.Println("Types:")
		for _, t := range pokemon.Types {
			fmt.Printf("  %s\n", t.Type.Name)
		}
	} else {
		return errors.New("you have not caught that pokemon")
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