package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type PokemonDetails struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
}

func GetPokemon(pokemonName string, c *config) (*PokemonDetails, error) {
	url := c.pokemonBaseUrl + pokemonName

	if data, ok := cache.Get(url); ok {
		var pokemon PokemonDetails
		if err := json.Unmarshal(data, &pokemon); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached data: %w", err)
		}
		return &pokemon, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch data from PokeAPI")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	cache.Add(url, body)

	var pokemon PokemonDetails
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return nil, err
	}

	return &pokemon, nil
}
