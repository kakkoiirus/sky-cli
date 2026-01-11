package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

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

	// Create context with timeout for API calls
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Get location
	location, err := api.GetLocation(ctx, cityName)
	if err != nil {
		fmt.Fprintln(os.Stderr, ui.FormatError(err))
		os.Exit(1)
		return
	}

	// Get weather
	weather, err := api.GetWeather(ctx, location.Latitude, location.Longitude)
	if err != nil {
		fmt.Fprintln(os.Stderr, ui.FormatError(err))
		os.Exit(1)
		return
	}

	// Display result
	fmt.Print(ui.FormatWeather(location, weather))
}
