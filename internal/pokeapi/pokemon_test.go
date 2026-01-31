package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Nachsus/pokedexcli/internal/pokecache"
)

func TestGetPokemon(t *testing.T) {
	// Create a test config
	testConfig := &config{
		pokemonBaseUrl: "",
	}

	// Create mock Pokemon data
	mockPokemon := PokemonDetails{
		Name:           "pikachu",
		BaseExperience: 112,
		Height:         4,
		Weight:         60,
	}
	mockData, _ := json.Marshal(mockPokemon)

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockData)
	}))
	defer server.Close()

	testConfig.pokemonBaseUrl = server.URL + "/"

	// Reset cache for test
	cache = pokecache.NewCache(5 * time.Minute)

	// Test getting Pokemon
	pokemon, err := GetPokemon("pikachu", testConfig)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if pokemon.Name != "pikachu" {
		t.Errorf("Expected name 'pikachu', got %s", pokemon.Name)
	}

	if pokemon.BaseExperience != 112 {
		t.Errorf("Expected base experience 112, got %d", pokemon.BaseExperience)
	}

	if pokemon.Height != 4 {
		t.Errorf("Expected height 4, got %d", pokemon.Height)
	}

	if pokemon.Weight != 60 {
		t.Errorf("Expected weight 60, got %d", pokemon.Weight)
	}
}

func TestGetPokemon_NotFound(t *testing.T) {
	testConfig := &config{
		pokemonBaseUrl: "",
	}

	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	testConfig.pokemonBaseUrl = server.URL + "/"

	// Reset cache for test
	cache = pokecache.NewCache(5 * time.Minute)

	// Test getting non-existent Pokemon
	_, err := GetPokemon("invalidpokemon", testConfig)
	if err == nil {
		t.Fatal("Expected error for non-existent Pokemon, got nil")
	}
}

func TestGetPokemon_InvalidJSON(t *testing.T) {
	testConfig := &config{
		pokemonBaseUrl: "",
	}

	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	testConfig.pokemonBaseUrl = server.URL + "/"

	// Reset cache for test
	cache = pokecache.NewCache(5 * time.Minute)

	// Test getting Pokemon with invalid JSON
	_, err := GetPokemon("pikachu", testConfig)
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}

func TestGetPokemon_UsesCache(t *testing.T) {
	testConfig := &config{
		pokemonBaseUrl: "",
	}

	mockPokemon := PokemonDetails{
		Name:           "charizard",
		BaseExperience: 240,
		Height:         17,
		Weight:         905,
	}
	mockData, _ := json.Marshal(mockPokemon)

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
		w.Write(mockData)
	}))
	defer server.Close()

	testConfig.pokemonBaseUrl = server.URL + "/"

	// Reset cache for test
	cache = pokecache.NewCache(5 * time.Minute)

	// First call - should hit the API
	pokemon1, err := GetPokemon("charizard", testConfig)
	if err != nil {
		t.Fatalf("Expected no error on first call, got %v", err)
	}

	// Second call - should use cache
	pokemon2, err := GetPokemon("charizard", testConfig)
	if err != nil {
		t.Fatalf("Expected no error on second call, got %v", err)
	}

	// Verify both results are the same
	if pokemon1.Name != pokemon2.Name || pokemon1.BaseExperience != pokemon2.BaseExperience {
		t.Error("Expected both calls to return the same data")
	}

	// Verify API was only called once
	if callCount != 1 {
		t.Errorf("Expected API to be called once, got %d calls", callCount)
	}
}

func TestGetPokemon_DifferentPokemon(t *testing.T) {
	testConfig := &config{
		pokemonBaseUrl: "",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var mockPokemon PokemonDetails
		if r.URL.Path == "/bulbasaur" {
			mockPokemon = PokemonDetails{
				Name:           "bulbasaur",
				BaseExperience: 64,
				Height:         7,
				Weight:         69,
			}
		} else if r.URL.Path == "/squirtle" {
			mockPokemon = PokemonDetails{
				Name:           "squirtle",
				BaseExperience: 63,
				Height:         5,
				Weight:         90,
			}
		}
		mockData, _ := json.Marshal(mockPokemon)
		w.WriteHeader(http.StatusOK)
		w.Write(mockData)
	}))
	defer server.Close()

	testConfig.pokemonBaseUrl = server.URL + "/"

	// Reset cache for test
	cache = pokecache.NewCache(5 * time.Minute)

	// Get different Pokemon
	bulbasaur, err := GetPokemon("bulbasaur", testConfig)
	if err != nil {
		t.Fatalf("Expected no error for bulbasaur, got %v", err)
	}

	squirtle, err := GetPokemon("squirtle", testConfig)
	if err != nil {
		t.Fatalf("Expected no error for squirtle, got %v", err)
	}

	// Verify they're different
	if bulbasaur.Name == squirtle.Name {
		t.Error("Expected different Pokemon to have different names")
	}

	if bulbasaur.BaseExperience == squirtle.BaseExperience {
		t.Error("Expected different Pokemon to have different base experience")
	}
}
