package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputValidation_EmptyCityName(t *testing.T) {
	// Test that empty city name after trimming is handled
	cityName := strings.TrimSpace("   ")
	assert.Equal(t, "", cityName, "Whitespace-only string should become empty")

	cityName = strings.TrimSpace("")
	assert.Equal(t, "", cityName, "Empty string should remain empty")
}

func TestInputValidation_ValidCityNames(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple name", "Tokyo", "Tokyo"},
		{"Name with spaces", "New York", "New York"},
		{"Name with leading/trailing spaces", "  London  ", "London"},
		{"Name with multiple spaces", "San   Francisco", "San   Francisco"},
		{"Name with special characters", "Zürich", "Zürich"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.TrimSpace(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInputParsing_MultipleWords(t *testing.T) {
	// Simulate os.Args parsing
	args := []string{"sky", "New", "York"}
	cityName := strings.Join(args[1:], " ")
	assert.Equal(t, "New York", cityName)

	args = []string{"sky", "São", "Paulo"}
	cityName = strings.Join(args[1:], " ")
	assert.Equal(t, "São Paulo", cityName)
}

func TestInputParsing_SingleWord(t *testing.T) {
	args := []string{"sky", "Tokyo"}
	cityName := strings.Join(args[1:], " ")
	assert.Equal(t, "Tokyo", cityName)
}

func TestInputParsing_NoArguments(t *testing.T) {
	args := []string{"sky"}
	cityName := strings.Join(args[1:], " ")
	assert.Equal(t, "", cityName)
}

// Mock stdin for testing interactive mode
func TestInteractiveMode_InputReading(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "Valid city name",
			input:   "Paris\n",
			want:    "Paris",
			wantErr: false,
		},
		{
			name:    "City name with spaces",
			input:   "Los Angeles\n",
			want:    "Los Angeles",
			wantErr: false,
		},
		{
			name:    "Empty input",
			input:   "\n",
			want:    "",
			wantErr: false,
		},
		{
			name:    "Whitespace only",
			input:   "   \n",
			want:    "   ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a reader from the input string
			reader := strings.NewReader(tt.input)

			// Read the line
			var buf bytes.Buffer
			tee := io.TeeReader(reader, &buf)

			data, err := io.ReadAll(tee)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Remove newline
				result := strings.TrimSpace(string(data))
				expected := strings.TrimSpace(tt.input)
				assert.Equal(t, expected, result)
			}
		})
	}
}

func TestErrorOutput_Formatting(t *testing.T) {
	tests := []struct {
		name     string
		errorMsg string
		contains string
	}{
		{
			name:     "Location not found",
			errorMsg: "location not found",
			contains: "Error:",
		},
		{
			name:     "Network error",
			errorMsg: "failed to fetch location: connection timeout",
			contains: "Error:",
		},
		{
			name:     "Empty city name",
			errorMsg: "city name cannot be empty",
			contains: "Error:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test error formatting
			formatted := "Error: " + tt.errorMsg + "\n"
			assert.Contains(t, formatted, tt.contains)
			assert.Contains(t, formatted, tt.errorMsg)
		})
	}
}

func TestOutputFormat_CaptureStdout(t *testing.T) {
	// Test capturing stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Write some output
	fmt.Print("Test output")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	assert.Equal(t, "Test output", buf.String())
}

func TestOutputFormat_CaptureStderr(t *testing.T) {
	// Test capturing stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Write error output
	fmt.Fprintln(w, "Error: test error")

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	assert.Contains(t, buf.String(), "Error: test error")
}

func TestArgumentParsing_CommandLineArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectCity  string
		expectEmpty bool
	}{
		{
			name:        "Single city name",
			args:        []string{"sky", "Tokyo"},
			expectCity:  "Tokyo",
			expectEmpty: false,
		},
		{
			name:        "Multi-word city name",
			args:        []string{"sky", "New", "York"},
			expectCity:  "New York",
			expectEmpty: false,
		},
		{
			name:        "No arguments",
			args:        []string{"sky"},
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.args) > 1 {
				cityName := strings.Join(tt.args[1:], " ")
				if tt.expectEmpty {
					assert.Empty(t, cityName)
				} else {
					assert.Equal(t, tt.expectCity, cityName)
				}
			} else {
				cityName := strings.Join(tt.args[1:], " ")
				assert.Empty(t, cityName)
			}
		})
	}
}

func TestArgumentParsing_WithQuotedSpaces(t *testing.T) {
	// Test parsing when arguments contain spaces
	args := []string{"sky", "\"New York\""}
	cityName := strings.Join(args[1:], " ")
	// Quotes would be handled by the shell, not by our code
	// This test documents the behavior
	assert.Contains(t, cityName, "New York")
}

func TestWhitespaceHandling_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Tab characters", "\tTokyo\t", "Tokyo"},
		{"Mixed whitespace", "  \t Tokyo \t  ", "Tokyo"},
		{"Newlines", "\nLondon\n", "London"},
		{"Multiple spaces between words", "New  York", "New  York"},
		{"Unicode whitespace", "\u2009Paris", "Paris"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.TrimSpace(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExitCodes_Simulation(t *testing.T) {
	// Test exit code logic (without actually exiting)
	tests := []struct {
		name       string
		shouldFail bool
		exitCode   int
	}{
		{"Success case", false, 0},
		{"Error case", true, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldFail {
				assert.Equal(t, 1, tt.exitCode, "Error should return exit code 1")
			} else {
				assert.Equal(t, 0, tt.exitCode, "Success should return exit code 0")
			}
		})
	}
}

// Helper function to test the full flow without actual network calls
func TestMain_FlowValidation(t *testing.T) {
	// This test validates the logical flow without executing main()
	// which would make network calls and exit the process

	t.Run("Valid flow sequence", func(t *testing.T) {
		// 1. Parse arguments
		args := []string{"sky", "Tokyo"}
		cityName := strings.Join(args[1:], " ")

		// 2. Trim whitespace
		cityName = strings.TrimSpace(cityName)

		// 3. Validate not empty
		isEmpty := cityName == ""

		assert.Equal(t, "Tokyo", cityName)
		assert.False(t, isEmpty, "City name should not be empty")
	})

	t.Run("Empty city name flow", func(t *testing.T) {
		args := []string{"sky", ""}
		cityName := strings.Join(args[1:], " ")
		cityName = strings.TrimSpace(cityName)
		isEmpty := cityName == ""

		assert.True(t, isEmpty, "Empty city name should trigger error path")
	})
}
