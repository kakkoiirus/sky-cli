# sky

A lightweight CLI weather utility that shows current weather conditions in your terminal.

![](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![](https://img.shields.io/badge/license-MIT-blue)

## Features

- No API key required (uses [Open-Meteo](https://open-meteo.com/))
- Cross-platform (Windows, macOS, Linux)
- Two modes: interactive and command-line arguments
- Weather emoji indicators
- Metric units (°C, km/h)
- Single binary, no dependencies

## Installation

### From pre-built binaries

Download the latest release from the [Releases](https://github.com/kakkoiirus/sky-cli/releases) page.

**macOS/Linux:**
```bash
# Download and make executable
chmod +x sky

# Move to PATH
sudo mv sky /usr/local/bin/
```

**Windows:**
Rename `sky.exe` and add to your PATH.

### From source

```bash
git clone https://github.com/kakkoiirus/sky-cli.git
cd sky-cli
go build -o sky cmd/sky/main.go
```

## Usage

### Command-line mode

```bash
sky Tokyo
sky London
sky "New York"
```

Output:
```
Tokyo, JP
Clear ☀️
Temp: 0.9°C
Feels like: -3.3°C
```

### Interactive mode

```bash
sky
```

Then enter your city name when prompted:
```
Enter city name: Paris
Paris, FR
Partly cloudy ⛅
Temp: 12°C
Feels like: 10°C
```

## API Data

Uses [Open-Meteo](https://open-meteo.com/) API:
- Geocoding API for city lookup
- Weather API for current conditions
- Apparent temperature calculation

## License

MIT

## Author

kakkoiirus
