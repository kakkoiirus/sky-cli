package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// WeatherResponse represents the response from Open-Meteo Weather API
type WeatherResponse struct {
	Current struct {
		Temperature     float64 `json:"temperature_2m"`
		ApparentTemp    float64 `json:"apparent_temperature"`
		WindSpeed       float64 `json:"windspeed_10m"`
		WeatherCode     int     `json:"weather_code"`
	} `json:"current"`
}

// Weather represents current weather conditions
type Weather struct {
	Temperature     float64
	ApparentTemp    float64
	WindSpeed       float64
	WeatherCode     int
	WeatherCodeDesc string
}

// WeatherCodeDescription returns a human-readable description and emoji for weather codes
// Based on WMO codes: https://open-meteo.com/en/docs
func WeatherCodeDescription(code int) string {
	descriptions := map[int]string{
		0:  "Clear",
		1:  "Mainly clear",
		2:  "Partly cloudy",
		3:  "Overcast",
		45: "Foggy",
		48: "Depositing rime fog",
		51: "Light drizzle",
		53: "Moderate drizzle",
		55: "Dense drizzle",
		61: "Slight rain",
		63: "Moderate rain",
		65: "Heavy rain",
		71: "Slight snow",
		73: "Moderate snow",
		75: "Heavy snow",
		77: "Snow grains",
		80: "Slight showers",
		81: "Moderate showers",
		82: "Violent showers",
		85: "Slight snow showers",
		86: "Heavy snow showers",
		95: "Thunderstorm",
		96: "Thunderstorm with hail",
		99: "Thunderstorm with heavy hail",
	}

	if desc, ok := descriptions[code]; ok {
		return desc
	}
	return "Unknown"
}

// WeatherCodeEmoji returns an emoji for a given weather code
func WeatherCodeEmoji(code int) string {
	emojis := map[int]string{
		0:  "☀️",
		1:  "🌤️",
		2:  "⛅",
		3:  "☁️",
		45: "🌫️",
		48: "🌫️",
		51: "🌧️",
		53: "🌧️",
		55: "🌧️",
		61: "🌧️",
		63: "🌧️",
		65: "🌧️",
		71: "🌨️",
		73: "🌨️",
		75: "❄️",
		77: "🌨️",
		80: "🌦️",
		81: "🌦️",
		82: "🌧️",
		85: "🌨️",
		86: "🌨️",
		95: "⛈️",
		96: "⛈️",
		99: "⛈️",
	}

	if emoji, ok := emojis[code]; ok {
		return emoji
	}
	return "🌡️"
}

// GetWeather retrieves current weather for a given location
func GetWeather(ctx context.Context, lat, lon float64) (*Weather, error) {
	apiURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,apparent_temperature,windspeed_10m,weather_code&temperature_unit=celsius&windspeed_unit=kmh&timezone=auto", lat, lon)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var weatherResp WeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &Weather{
		Temperature:     weatherResp.Current.Temperature,
		ApparentTemp:    weatherResp.Current.ApparentTemp,
		WindSpeed:       weatherResp.Current.WindSpeed,
		WeatherCode:     weatherResp.Current.WeatherCode,
		WeatherCodeDesc: WeatherCodeDescription(weatherResp.Current.WeatherCode),
	}, nil
}
