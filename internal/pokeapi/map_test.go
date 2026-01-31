package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMapsForward(t *testing.T) {
	t.Run("uses base URL when mapNextUrl is empty", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := LocationAreaResponse{
				Count:    20,
				Next:     "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
				Previous: "",
				Results: []LocationArea{
					{Name: "canalave-city-area", URL: "https://pokeapi.co/api/v2/location-area/1/"},
					{Name: "eterna-city-area", URL: "https://pokeapi.co/api/v2/location-area/2/"},
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		c := &config{
			mapBaseUrl: server.URL,
			mapNextUrl: "",
			mapPrevUrl: "",
		}

		areas, err := MapsForward(c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(areas) != 2 {
			t.Errorf("expected 2 areas, got %d", len(areas))
		}

		if areas[0] != "canalave-city-area" {
			t.Errorf("expected first area to be 'canalave-city-area', got '%s'", areas[0])
		}

		if c.mapNextUrl != "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20" {
			t.Errorf("expected mapNextUrl to be updated, got '%s'", c.mapNextUrl)
		}
	})

	t.Run("uses mapNextUrl when not empty", func(t *testing.T) {
		requestCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			if r.URL.Path == "/next" {
				response := LocationAreaResponse{
					Count:    20,
					Next:     "https://pokeapi.co/api/v2/location-area/?offset=40&limit=20",
					Previous: "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
					Results: []LocationArea{
						{Name: "pastoria-city-area", URL: "https://pokeapi.co/api/v2/location-area/3/"},
					},
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response)
			}
		}))
		defer server.Close()

		c := &config{
			mapBaseUrl: server.URL,
			mapNextUrl: server.URL + "/next",
			mapPrevUrl: "",
		}

		areas, err := MapsForward(c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(areas) != 1 {
			t.Errorf("expected 1 area, got %d", len(areas))
		}

		if areas[0] != "pastoria-city-area" {
			t.Errorf("expected first area to be 'pastoria-city-area', got '%s'", areas[0])
		}
	})
}

func TestMapsBackward(t *testing.T) {
	t.Run("uses base URL when mapPrevUrl is empty", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := LocationAreaResponse{
				Count:    20,
				Next:     "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
				Previous: "",
				Results: []LocationArea{
					{Name: "canalave-city-area", URL: "https://pokeapi.co/api/v2/location-area/1/"},
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		c := &config{
			mapBaseUrl: server.URL,
			mapNextUrl: "",
			mapPrevUrl: "",
		}

		areas, err := MapsBackward(c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(areas) != 1 {
			t.Errorf("expected 1 area, got %d", len(areas))
		}

		if areas[0] != "canalave-city-area" {
			t.Errorf("expected first area to be 'canalave-city-area', got '%s'", areas[0])
		}
	})

	t.Run("uses mapPrevUrl when not empty", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/prev" {
				response := LocationAreaResponse{
					Count:    20,
					Next:     "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
					Previous: "",
					Results: []LocationArea{
						{Name: "jubilife-city-area", URL: "https://pokeapi.co/api/v2/location-area/4/"},
					},
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response)
			}
		}))
		defer server.Close()

		c := &config{
			mapBaseUrl: server.URL,
			mapNextUrl: "",
			mapPrevUrl: server.URL + "/prev",
		}

		areas, err := MapsBackward(c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(areas) != 1 {
			t.Errorf("expected 1 area, got %d", len(areas))
		}

		if areas[0] != "jubilife-city-area" {
			t.Errorf("expected first area to be 'jubilife-city-area', got '%s'", areas[0])
		}
	})
}

func TestGetMaps(t *testing.T) {
	t.Run("successfully fetches and parses location areas", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := LocationAreaResponse{
				Count:    1281,
				Next:     "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
				Previous: "",
				Results: []LocationArea{
					{Name: "canalave-city-area", URL: "https://pokeapi.co/api/v2/location-area/1/"},
					{Name: "eterna-city-area", URL: "https://pokeapi.co/api/v2/location-area/2/"},
					{Name: "pastoria-city-area", URL: "https://pokeapi.co/api/v2/location-area/3/"},
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		c := &config{
			mapBaseUrl: server.URL,
			mapNextUrl: "",
			mapPrevUrl: "",
		}

		areas, err := GetMaps(server.URL, c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(areas) != 3 {
			t.Errorf("expected 3 areas, got %d", len(areas))
		}

		expectedAreas := []string{"canalave-city-area", "eterna-city-area", "pastoria-city-area"}
		for i, expected := range expectedAreas {
			if areas[i] != expected {
				t.Errorf("expected area %d to be '%s', got '%s'", i, expected, areas[i])
			}
		}

		if c.mapNextUrl != "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20" {
			t.Errorf("expected mapNextUrl to be updated, got '%s'", c.mapNextUrl)
		}

		if c.mapPrevUrl != "" {
			t.Errorf("expected mapPrevUrl to be empty, got '%s'", c.mapPrevUrl)
		}
	})

	t.Run("updates config with next and previous URLs", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := LocationAreaResponse{
				Count:    1281,
				Next:     "https://pokeapi.co/api/v2/location-area/?offset=40&limit=20",
				Previous: "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
				Results: []LocationArea{
					{Name: "sunyshore-city-area", URL: "https://pokeapi.co/api/v2/location-area/4/"},
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		c := &config{
			mapBaseUrl: server.URL,
			mapNextUrl: "",
			mapPrevUrl: "",
		}

		_, err := GetMaps(server.URL, c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if c.mapNextUrl != "https://pokeapi.co/api/v2/location-area/?offset=40&limit=20" {
			t.Errorf("expected mapNextUrl to be updated, got '%s'", c.mapNextUrl)
		}

		if c.mapPrevUrl != "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20" {
			t.Errorf("expected mapPrevUrl to be updated, got '%s'", c.mapPrevUrl)
		}
	})

	t.Run("returns error on non-200 status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		c := &config{
			mapBaseUrl: server.URL,
			mapNextUrl: "",
			mapPrevUrl: "",
		}

		_, err := GetMaps(server.URL, c)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		expectedError := "failed to fetch data from PokeAPI"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("returns error on invalid JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		c := &config{
			mapBaseUrl: server.URL,
			mapNextUrl: "",
			mapPrevUrl: "",
		}

		_, err := GetMaps(server.URL, c)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns error on network failure", func(t *testing.T) {
		c := &config{
			mapBaseUrl: "http://invalid-url-that-does-not-exist.local",
			mapNextUrl: "",
			mapPrevUrl: "",
		}

		_, err := GetMaps("http://invalid-url-that-does-not-exist.local", c)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("handles empty results", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := LocationAreaResponse{
				Count:    0,
				Next:     "",
				Previous: "",
				Results:  []LocationArea{},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		c := &config{
			mapBaseUrl: server.URL,
			mapNextUrl: "",
			mapPrevUrl: "",
		}

		areas, err := GetMaps(server.URL, c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(areas) != 0 {
			t.Errorf("expected 0 areas, got %d", len(areas))
		}
	})
}
