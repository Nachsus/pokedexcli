package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func MapsForward(c *config) ([]string, error) {
	var url string
	switch c.mapNextUrl {
	case "":
		url = c.mapBaseUrl
	default:
		url = c.mapNextUrl
	}

	return GetMaps(url, c)
}

func MapsBackward(c *config) ([]string, error) {
	var url string
	switch c.mapPrevUrl {
	case "":
		url = c.mapBaseUrl
	default:
		url = c.mapPrevUrl
	}

	return GetMaps(url, c)
}

func GetMaps(url string, c *config) ([]string, error) {
	fmt.Printf("Requesting URL: %s\n", url)
	if data, ok := cache.Get(url); ok {
		var response LocationAreaResponse
		if err := json.Unmarshal(data, &response); err == nil {
			c.mapNextUrl = response.Next
			c.mapPrevUrl = response.Previous

			var areaNames []string
			for _, area := range response.Results {
				areaNames = append(areaNames, area.Name)
			}
			fmt.Println("got from cache")
			return areaNames, nil
		}
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
	fmt.Println("Added to cache")

	var response LocationAreaResponse
	if err := json.Unmarshal(body, &response); err != nil {
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
