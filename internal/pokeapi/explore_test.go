package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPokemonFromArea(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := LocationAreaDetail{
			PokemonEncounters: []PokemonEncounter{
				{
					Pokemon: Pokemon{
						Name: "pikachu",
						URL:  "https://pokeapi.co/api/v2/pokemon/25/",
					},
				},
				{
					Pokemon: Pokemon{
						Name: "bulbasaur",
						URL:  "https://pokeapi.co/api/v2/pokemon/1/",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create a test config
	testConfig := &config{
		mapBaseUrl: server.URL + "/",
		mapNextUrl: "",
		mapPrevUrl: "",
	}

	// Test getting pokemon from area
	pokemonNames, err := GetPokemonFromArea("test-area", testConfig)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(pokemonNames) != 2 {
		t.Errorf("expected 2 pokemon, got %d", len(pokemonNames))
	}

	if pokemonNames[0] != "pikachu" {
		t.Errorf("expected pikachu, got %s", pokemonNames[0])
	}

	if pokemonNames[1] != "bulbasaur" {
		t.Errorf("expected bulbasaur, got %s", pokemonNames[1])
	}
}

func TestGetPokemonFromArea_EmptyEncounters(t *testing.T) {
	// Create a test server with no pokemon
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := LocationAreaDetail{
			PokemonEncounters: []PokemonEncounter{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	testConfig := &config{
		mapBaseUrl: server.URL + "/",
		mapNextUrl: "",
		mapPrevUrl: "",
	}

	pokemonNames, err := GetPokemonFromArea("empty-area", testConfig)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(pokemonNames) != 0 {
		t.Errorf("expected 0 pokemon, got %d", len(pokemonNames))
	}
}

func TestGetPokemonFromArea_HTTPError(t *testing.T) {
	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	testConfig := &config{
		mapBaseUrl: server.URL + "/",
		mapNextUrl: "",
		mapPrevUrl: "",
	}

	_, err := GetPokemonFromArea("nonexistent-area", testConfig)
	if err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestGetPokemonFromArea_InvalidJSON(t *testing.T) {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	testConfig := &config{
		mapBaseUrl: server.URL + "/",
		mapNextUrl: "",
		mapPrevUrl: "",
	}

	_, err := GetPokemonFromArea("invalid-area", testConfig)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestGetPokemonFromArea_Cache(t *testing.T) {
	callCount := 0

	// Create a test server that tracks calls
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		response := LocationAreaDetail{
			PokemonEncounters: []PokemonEncounter{
				{
					Pokemon: Pokemon{
						Name: "charmander",
						URL:  "https://pokeapi.co/api/v2/pokemon/4/",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	testConfig := &config{
		mapBaseUrl: server.URL + "/",
		mapNextUrl: "",
		mapPrevUrl: "",
	}

	// First call - should hit the API
	pokemonNames1, err := GetPokemonFromArea("cache-test-area", testConfig)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected 1 API call, got %d", callCount)
	}

	// Second call - should use cache
	pokemonNames2, err := GetPokemonFromArea("cache-test-area", testConfig)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected still 1 API call (cached), got %d", callCount)
	}

	// Verify both results are the same
	if len(pokemonNames1) != len(pokemonNames2) {
		t.Error("cached result different from original")
	}

	if pokemonNames1[0] != pokemonNames2[0] {
		t.Error("cached pokemon name differs from original")
	}
}
