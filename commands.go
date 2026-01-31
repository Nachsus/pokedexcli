package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Nachsus/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

var supportedCommands map[string]cliCommand

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
