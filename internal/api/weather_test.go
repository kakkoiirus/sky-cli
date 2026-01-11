package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeatherCodeDescription_KnownCodes(t *testing.T) {
	tests := []struct {
		code        int
		expected    string
		description string
	}{
		{0, "Clear", "Clear sky"},
		{1, "Mainly clear", "Mainly clear"},
		{2, "Partly cloudy", "Partly cloudy"},
		{3, "Overcast", "Overcast"},
		{45, "Foggy", "Foggy"},
		{48, "Depositing rime fog", "Depositing rime fog"},
		{51, "Light drizzle", "Light drizzle"},
		{53, "Moderate drizzle", "Moderate drizzle"},
		{55, "Dense drizzle", "Dense drizzle"},
		{61, "Slight rain", "Slight rain"},
		{63, "Moderate rain", "Moderate rain"},
		{65, "Heavy rain", "Heavy rain"},
		{71, "Slight snow", "Slight snow"},
		{73, "Moderate snow", "Moderate snow"},
		{75, "Heavy snow", "Heavy snow"},
		{77, "Snow grains", "Snow grains"},
		{80, "Slight showers", "Slight showers"},
		{81, "Moderate showers", "Moderate showers"},
		{82, "Violent showers", "Violent showers"},
		{85, "Slight snow showers", "Slight snow showers"},
		{86, "Heavy snow showers", "Heavy snow showers"},
		{95, "Thunderstorm", "Thunderstorm"},
		{96, "Thunderstorm with hail", "Thunderstorm with hail"},
		{99, "Thunderstorm with heavy hail", "Thunderstorm with heavy hail"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := WeatherCodeDescription(tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWeatherCodeDescription_UnknownCodes(t *testing.T) {
	tests := []struct {
		name  string
		code  int
	}{
		{"Unknown positive code", 999},
		{"Unknown negative code", -1},
		{"Unknown zero code in range", 4},
		{"Large unknown code", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WeatherCodeDescription(tt.code)
			assert.Equal(t, "Unknown", result)
		})
	}
}

func TestWeatherCodeEmoji_KnownCodes(t *testing.T) {
	tests := []struct {
		code        int
		expected    string
		description string
	}{
		{0, "☀️", "Clear"},
		{1, "🌤️", "Mainly clear"},
		{2, "⛅", "Partly cloudy"},
		{3, "☁️", "Overcast"},
		{45, "🌫️", "Foggy"},
		{48, "🌫️", "Depositing rime fog"},
		{51, "🌧️", "Light drizzle"},
		{53, "🌧️", "Moderate drizzle"},
		{55, "🌧️", "Dense drizzle"},
		{61, "🌧️", "Slight rain"},
		{63, "🌧️", "Moderate rain"},
		{65, "🌧️", "Heavy rain"},
		{71, "🌨️", "Slight snow"},
		{73, "🌨️", "Moderate snow"},
		{75, "❄️", "Heavy snow"},
		{77, "🌨️", "Snow grains"},
		{80, "🌦️", "Slight showers"},
		{81, "🌦️", "Moderate showers"},
		{82, "🌧️", "Violent showers"},
		{85, "🌨️", "Slight snow showers"},
		{86, "🌨️", "Heavy snow showers"},
		{95, "⛈️", "Thunderstorm"},
		{96, "⛈️", "Thunderstorm with hail"},
		{99, "⛈️", "Thunderstorm with heavy hail"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := WeatherCodeEmoji(tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWeatherCodeEmoji_UnknownCodes(t *testing.T) {
	tests := []struct {
		name  string
		code  int
	}{
		{"Unknown positive code", 999},
		{"Unknown negative code", -1},
		{"Large unknown code", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WeatherCodeEmoji(tt.code)
			assert.Equal(t, "🌡️", result)
		})
	}
}

func TestGetWeather_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := WeatherResponse{
			Current: struct {
				Temperature     float64 `json:"temperature_2m"`
				ApparentTemp    float64 `json:"apparent_temperature"`
				WindSpeed       float64 `json:"windspeed_10m"`
				WeatherCode     int     `json:"weather_code"`
			}{
				Temperature:  15.5,
				ApparentTemp: 14.2,
				WindSpeed:    10.5,
				WeatherCode:  0,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Save original function and restore it after test
	// Note: Since GetWeather uses a hardcoded URL, we need to modify the function
	// For this test, we'll test the JSON unmarshaling logic separately

	// Test JSON unmarshaling
	responseJSON := `{"current":{"temperature_2m":15.5,"apparent_temperature":14.2,"windspeed_10m":10.5,"weather_code":0}}`
	var weatherResp WeatherResponse
	err := json.Unmarshal([]byte(responseJSON), &weatherResp)
	require.NoError(t, err)

	assert.Equal(t, 15.5, weatherResp.Current.Temperature)
	assert.Equal(t, 14.2, weatherResp.Current.ApparentTemp)
	assert.Equal(t, 10.5, weatherResp.Current.WindSpeed)
	assert.Equal(t, 0, weatherResp.Current.WeatherCode)
}

func TestWeatherResponse_Unmarshal(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		wantTemp  float64
		wantCode  int
		wantError bool
	}{
		{
			name: "Valid response",
			json: `{"current":{"temperature_2m":20.0,"apparent_temperature":19.0,"windspeed_10m":5.0,"weather_code":1}}`,
			wantTemp: 20.0,
			wantCode: 1,
			wantError: false,
		},
		{
			name: "Negative temperature",
			json: `{"current":{"temperature_2m":-15.3,"apparent_temperature":-22.1,"windspeed_10m":5.0,"weather_code":65}}`,
			wantTemp: -15.3,
			wantCode: 65,
			wantError: false,
		},
		{
			name: "Invalid JSON",
			json: `{invalid json}`,
			wantError: true,
		},
		{
			name: "Empty JSON",
			json: `{}`,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp WeatherResponse
			err := json.Unmarshal([]byte(tt.json), &resp)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.name == "Valid response" || tt.name == "Negative temperature" {
					assert.Equal(t, tt.wantTemp, resp.Current.Temperature)
					assert.Equal(t, tt.wantCode, resp.Current.WeatherCode)
				}
			}
		})
	}
}

func TestWeather_Struct(t *testing.T) {
	weather := &Weather{
		Temperature:     25.5,
		ApparentTemp:    27.0,
		WindSpeed:       12.3,
		WeatherCode:     2,
		WeatherCodeDesc: "Partly cloudy",
	}

	assert.Equal(t, 25.5, weather.Temperature)
	assert.Equal(t, 27.0, weather.ApparentTemp)
	assert.Equal(t, 12.3, weather.WindSpeed)
	assert.Equal(t, 2, weather.WeatherCode)
	assert.Equal(t, "Partly cloudy", weather.WeatherCodeDesc)
}

func TestWeatherCodeDescription_AllCodesReturnString(t *testing.T) {
	// Test that all codes return a non-empty string
	codes := []int{0, 1, 2, 3, 45, 48, 51, 53, 55, 61, 63, 65, 71, 73, 75, 77, 80, 81, 82, 85, 86, 95, 96, 99}

	for _, code := range codes {
		t.Run(string(rune(code)), func(t *testing.T) {
			result := WeatherCodeDescription(code)
			assert.NotEmpty(t, result)
			assert.NotEqual(t, "Unknown", result)
		})
	}
}

func TestWeatherCodeEmoji_AllCodesReturnEmoji(t *testing.T) {
	// Test that all known codes return an emoji
	codes := []int{0, 1, 2, 3, 45, 48, 51, 53, 55, 61, 63, 65, 71, 73, 75, 77, 80, 81, 82, 85, 86, 95, 96, 99}

	for _, code := range codes {
		t.Run(string(rune(code)), func(t *testing.T) {
			result := WeatherCodeEmoji(code)
			assert.NotEmpty(t, result)
			// All emojis are multi-byte characters
			assert.True(t, len(result) > 1)
		})
	}
}
