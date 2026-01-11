package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GeocodingResponse represents the response from Open-Meteo Geocoding API
type GeocodingResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Country   string  `json:"country_code"`
	} `json:"results"`
}

// Location represents a geographic location
type Location struct {
	Name      string
	Latitude  float64
	Longitude float64
	Country   string
}

// GetLocation retrieves coordinates for a given city name
func GetLocation(city string) (*Location, error) {
	encodedCity := url.QueryEscape(city)
	apiURL := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json", encodedCity)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch location: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var geoResp GeocodingResponse
	if err := json.Unmarshal(body, &geoResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(geoResp.Results) == 0 {
		return nil, fmt.Errorf("location not found")
	}

	result := geoResp.Results[0]
	return &Location{
		Name:      result.Name,
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
		Country:   result.Country,
	}, nil
}
