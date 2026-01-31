package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type PokemonDetails struct {
	Name           string         `json:"name"`
	BaseExperience int            `json:"base_experience"`
	Height         int            `json:"height"`
	Weight         int            `json:"weight"`
	Stats          map[string]int `json:"-"`
	Types          []string       `json:"-"`
}

type pokemonAPIResponse struct {
	Name           string           `json:"name"`
	BaseExperience int              `json:"base_experience"`
	Height         int              `json:"height"`
	Weight         int              `json:"weight"`
	Stats          []pokemonStatAPI `json:"stats"`
	Types          []pokemonTypeAPI `json:"types"`
}

type pokemonStatAPI struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type pokemonTypeAPI struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

func GetPokemon(pokemonName string, c *config) (*PokemonDetails, error) {
	url := c.pokemonBaseUrl + pokemonName

	if data, ok := cache.Get(url); ok {
		var apiResponse pokemonAPIResponse
		if err := json.Unmarshal(data, &apiResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached data: %w", err)
		}
		return convertToPokemonDetails(apiResponse), nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, errors.New("pokemon not found")
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch data from PokeAPI")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	cache.Add(url, body)

	var apiResponse pokemonAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, err
	}

	return convertToPokemonDetails(apiResponse), nil
}

func convertToPokemonDetails(apiResponse pokemonAPIResponse) *PokemonDetails {
	// Convert stats to map
	stats := make(map[string]int)
	for _, stat := range apiResponse.Stats {
		stats[stat.Stat.Name] = stat.BaseStat
	}

	// Convert types to string slice
	types := make([]string, len(apiResponse.Types))
	for i, t := range apiResponse.Types {
		types[i] = t.Type.Name
	}

	return &PokemonDetails{
		Name:           apiResponse.Name,
		BaseExperience: apiResponse.BaseExperience,
		Height:         apiResponse.Height,
		Weight:         apiResponse.Weight,
		Stats:          stats,
		Types:          types,
	}
}
