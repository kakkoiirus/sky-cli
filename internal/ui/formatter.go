package ui

import (
	"fmt"

	"github.com/kakkoiirus/sky-cli/internal/api"
)

// FormatWeather formats the weather data for display
func FormatWeather(location *api.Location, weather *api.Weather) string {
	emoji := api.WeatherCodeEmoji(weather.WeatherCode)

	return fmt.Sprintf("%s, %s\n%s %s\nTemp: %.1f°C\nFeels like: %.1f°C\n",
		location.Name,
		location.Country,
		weather.WeatherCodeDesc,
		emoji,
		weather.Temperature,
		weather.ApparentTemp,
	)
}

// FormatError formats an error message for stderr
func FormatError(err error) string {
	return fmt.Sprintf("Error: %s\n", err.Error())
}
