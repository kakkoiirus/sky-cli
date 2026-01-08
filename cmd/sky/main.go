package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kakkoiirus/sky-cli/internal/api"
	"github.com/kakkoiirus/sky-cli/internal/ui"
)

func main() {
	var cityName string

	// Check if city name is provided as argument
	if len(os.Args) > 1 {
		cityName = strings.Join(os.Args[1:], " ")
	} else {
		// Interactive mode
		fmt.Print("Enter city name: ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, ui.FormatError(fmt.Errorf("failed to read input")))
			os.Exit(1)
			return
		}
		cityName = scanner.Text()
	}

	// Trim whitespace
	cityName = strings.TrimSpace(cityName)
	if cityName == "" {
		fmt.Fprintln(os.Stderr, ui.FormatError(fmt.Errorf("city name cannot be empty")))
		os.Exit(1)
		return
	}

	// Get location
	location, err := api.GetLocation(cityName)
	if err != nil {
		fmt.Fprintln(os.Stderr, ui.FormatError(err))
		os.Exit(1)
		return
	}

	// Get weather
	weather, err := api.GetWeather(location.Latitude, location.Longitude)
	if err != nil {
		fmt.Fprintln(os.Stderr, ui.FormatError(err))
		os.Exit(1)
		return
	}

	// Display result
	fmt.Print(ui.FormatWeather(location, weather))
}
