package ui

import (
	"errors"
	"testing"

	"github.com/kakkoiirus/sky-cli/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestFormatWeather_Typical(t *testing.T) {
	location := &api.Location{
		Name:      "Tokyo",
		Country:   "JP",
		Latitude:  35.6762,
		Longitude: 139.6503,
	}

	weather := &api.Weather{
		Temperature:     15.5,
		ApparentTemp:    14.2,
		WeatherCode:     0,
		WeatherCodeDesc: "Clear",
	}

	output := FormatWeather(location, weather)

	assert.Contains(t, output, "Tokyo, JP")
	assert.Contains(t, output, "Clear ☀️")
	assert.Contains(t, output, "Temp: 15.5°C")
	assert.Contains(t, output, "Feels like: 14.2°C")
}

func TestFormatWeather_NegativeTemperature(t *testing.T) {
	location := &api.Location{
		Name:      "Moscow",
		Country:   "RU",
		Latitude:  55.7558,
		Longitude: 37.6173,
	}

	weather := &api.Weather{
		Temperature:     -15.3,
		ApparentTemp:    -22.1,
		WeatherCode:     65,
		WeatherCodeDesc: "Heavy rain",
	}

	output := FormatWeather(location, weather)

	assert.Contains(t, output, "Moscow, RU")
	assert.Contains(t, output, "Temp: -15.3°C")
	assert.Contains(t, output, "Feels like: -22.1°C")
}

func TestFormatWeather_VeryLongCityName(t *testing.T) {
	location := &api.Location{
		Name:      "Llanfairpwllgwyngyllgogerychwyrndrobwllllantysiliogogogoch",
		Country:   "GB",
		Latitude:  53.2224,
		Longitude: -4.2179,
	}

	weather := &api.Weather{
		Temperature:     10.0,
		ApparentTemp:    9.0,
		WeatherCode:     3,
		WeatherCodeDesc: "Overcast",
	}

	output := FormatWeather(location, weather)

	assert.Contains(t, output, "Llanfairpwllgwyngyllgogerychwyrndrobwllllantysiliogogogoch, GB")
	assert.Contains(t, output, "Temp: 10.0°C")
}

func TestFormatWeather_SpecialCharacters(t *testing.T) {
	location := &api.Location{
		Name:      "Zürich",
		Country:   "CH",
		Latitude:  47.3769,
		Longitude: 8.5417,
	}

	weather := &api.Weather{
		Temperature:     5.5,
		ApparentTemp:    2.1,
		WeatherCode:     61,
		WeatherCodeDesc: "Slight rain",
	}

	output := FormatWeather(location, weather)

	assert.Contains(t, output, "Zürich, CH")
	assert.Contains(t, output, "Slight rain 🌧️")
}

func TestFormatWeather_DifferentWeatherCodes(t *testing.T) {
	tests := []struct {
		name         string
		weatherCode  int
		description  string
		expectedEmoji string
	}{
		{"Clear sky", 0, "Clear", "☀️"},
		{"Partly cloudy", 2, "Partly cloudy", "⛅"},
		{"Thunderstorm", 95, "Thunderstorm", "⛈️"},
		{"Snow", 75, "Heavy snow", "❄️"},
	}

	location := &api.Location{
		Name:      "Test City",
		Country:   "TC",
		Latitude:  0.0,
		Longitude: 0.0,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weather := &api.Weather{
				Temperature:     20.0,
				ApparentTemp:    19.0,
				WeatherCode:     tt.weatherCode,
				WeatherCodeDesc: tt.description,
			}

			output := FormatWeather(location, weather)

			assert.Contains(t, output, tt.expectedEmoji)
		})
	}
}

func TestFormatError_Typical(t *testing.T) {
	err := errors.New("location not found")

	output := FormatError(err)

	assert.Equal(t, "Error: location not found\n", output)
}

func TestFormatError_LongMessage(t *testing.T) {
	longErr := errors.New("failed to fetch location: connection timeout after 30 seconds while attempting to connect to api.open-meteo.com")

	output := FormatError(longErr)

	assert.Contains(t, output, "Error:")
	assert.Contains(t, output, "connection timeout")
}

func TestFormatError_WrappedError(t *testing.T) {
	wrappedErr := errors.New("failed to fetch location: network unreachable")

	output := FormatError(wrappedErr)

	assert.Equal(t, "Error: failed to fetch location: network unreachable\n", output)
}
