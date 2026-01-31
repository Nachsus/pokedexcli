package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type LocationAreaDetail struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func GetPokemonFromArea(area string, c *config) ([]string, error) {
	url := c.mapBaseUrl + area
	if data, ok := cache.Get(url); ok {
		var response LocationAreaDetail
		if err := json.Unmarshal(data, &response); err == nil {
			var pokemonNames []string
			for _, encounter := range response.PokemonEncounters {
				pokemonNames = append(pokemonNames, encounter.Pokemon.Name)
			}
			return pokemonNames, nil
		}
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to getch data from PokeAPI")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	cache.Add(url, body)

	var response LocationAreaDetail
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	var pokemonNames []string
	for _, encounter := range response.PokemonEncounters {
		pokemonNames = append(pokemonNames, encounter.Pokemon.Name)
	}

	return pokemonNames, nil
}
