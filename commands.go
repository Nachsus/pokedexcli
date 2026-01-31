package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/Nachsus/pokedexcli/internal/pokeapi"
	"github.com/Nachsus/pokedexcli/internal/pokedex"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

var supportedCommands map[string]cliCommand
var userPokedex = pokedex.NewPokedex()

func init() {
	supportedCommands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Lists 20 maps incrementing",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Lists 20 maps decrementing",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Lists pokemon in given location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a Pokemon based on its base experience",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects the data for a Pokemon you have in your Pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists names of all caught Pokemon",
			callback:    commandPokedex,
		},
	}
}

func commandExit(args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range supportedCommands {
		fmt.Println(command.name + ": " + command.description)
	}
	return nil
}

func commandMap(args []string) error {
	areaNames, err := pokeapi.MapsForward(&pokeapi.Conf)
	if err != nil {
		return err
	}

	for _, name := range areaNames {
		fmt.Println(name)
	}

	return nil
}

func commandMapB(args []string) error {
	areaNames, err := pokeapi.MapsBackward(&pokeapi.Conf)
	if err != nil {
		return err
	}

	for _, name := range areaNames {
		fmt.Println(name)
	}

	return nil
}

func commandExplore(args []string) error {
	if len(args) == 0 {
		return errors.New("please provide a location area name")
	}

	areaName := args[0]
	fmt.Printf("Exploring %s...\n", areaName)

	pokemonNames, err := pokeapi.GetPokemonFromArea(areaName, &pokeapi.Conf)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, name := range pokemonNames {
		fmt.Println(" - " + name)
	}

	return nil
}

func commandCatch(args []string) error {
	if len(args) == 0 {
		return errors.New("please provide a pokemon name")
	}

	pokemonName := args[0]

	if userPokedex.Has(pokemonName) {
		fmt.Printf("You already caught %s!\n", pokemonName)
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := pokeapi.GetPokemon(pokemonName, &pokeapi.Conf)
	if err != nil {
		return err
	}

	const minChance = 30.0
	const maxChance = 80.0
	const maxBaseXP = 608.0

	catchChance := maxChance - ((float64(pokemon.BaseExperience) / maxBaseXP) * (maxChance - minChance))
	if catchChance < minChance {
		catchChance = minChance
	}
	if catchChance > maxChance {
		catchChance = maxChance
	}

	roll := rand.Float64() * 100.0
	if roll <= catchChance {
		userPokedex.Add(*pokemon)
		fmt.Printf("%s was caught!\n", pokemonName)
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func commandInspect(args []string) error {
	if len(args) == 0 {
		return errors.New("please provide a pokemon name")
	}

	pokemonName := args[0]

	if !userPokedex.Has(pokemonName) {
		return errors.New("You have not caught " + pokemonName)
	}

	pokemon, ok := userPokedex.Get(pokemonName)
	if !ok {
		return errors.New("error getting pokemon from pokedex")
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	fmt.Printf("  -hp: %d\n", pokemon.Stats["hp"])
	fmt.Printf("  -attack: %d\n", pokemon.Stats["attack"])
	fmt.Printf("  -defense: %d\n", pokemon.Stats["defense"])
	fmt.Printf("  -special-attack: %d\n", pokemon.Stats["special-attack"])
	fmt.Printf("  -special-defense: %d\n", pokemon.Stats["special-defense"])
	fmt.Printf("  -speed: %d\n", pokemon.Stats["speed"])
	fmt.Println("Types:")
	for _, typeName := range pokemon.Types {
		fmt.Printf("  - %s\n", typeName)
	}

	return nil
}

func commandPokedex(args []string) error {
	pokemon := userPokedex.GetAll()
	if len(pokemon) < 1 {
		fmt.Println("No pokemon in your pokedex")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, poke := range pokemon {
		fmt.Println(" - " + poke.Name)
	}

	return nil
}
