package main

import (
	"time"

	"github.com/bklimov-web/pokedex-go/internal/pokeapi"
)


func main() {
	// TODO: change interval to minutes
	client := pokeapi.NewClient(5 * time.Second, 5 * time.Second)

	config := &config{
		pokeapiClient: client,
		caughtPokemon: map[string]pokeapi.PokemonResponse{},
	}

	startRepl(config)
}
