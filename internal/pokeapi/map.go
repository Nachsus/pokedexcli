package pokeapi

import (
	"encoding/json"
	"errors"
	"net/http"
)

type LocationAreaResponse struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func GetMaps(c *config) ([]string, error) {
	var url string
	if c.mapNextUrl == "" {
		url = c.mapBaseUrl
	} else {
		url = c.mapNextUrl
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch data from PokeAPI")
	}

	var response LocationAreaResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	c.mapNextUrl = response.Next
	c.mapPrevUrl = response.Previous

	var areaNames []string
	for _, area := range response.Results {
		areaNames = append(areaNames, area.Name)
	}

	return areaNames, nil
}
