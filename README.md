# V2Ray Client TUI

A Terminal User Interface (TUI) application for generating V2Ray client configurations from VMess links.

## Features

- **Interactive TUI**: Beautiful terminal-based user interface using the `tview` library
- **VMess Parser**: Converts VMess links to both V2Ray and sing-box configurations
- **Real-time Preview**: See the generated configuration before saving
- **Easy Configuration**: Simple input field for VMess links
- **File Export**: Save configurations to `config.json`

## Installation

1. Make sure you have Go 1.24.4 or later installed
2. Clone or download this repository
3. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

### Running the Application

```bash
go run .
```

Or build and run:
```bash
go build
./go_v2ray_client
```

### Using the TUI

1. **Enter VMess Link**: Type or paste your VMess link in the input field
2. **Parse Configuration**: Click "Parse VMess" or press `Ctrl+P`
3. **Preview**: View the generated configuration in the preview area
4. **Save**: Click "Save Config" or press `Ctrl+S` to save to `config.json`
5. **Clear**: Use "Clear" button to reset the interface
6. **Quit**: Click "Quit" or press `Ctrl+C` to exit

### Keyboard Shortcuts

- `Ctrl+P`: Parse VMess link
- `Ctrl+S`: Save configuration
- `Ctrl+C`: Quit application
- `Tab`: Navigate between elements
- `Enter`: Activate buttons
- Arrow keys: Navigate and scroll

### Example VMess Link

```
vmess://ewogICJ2IjogIjIiLAogICJwcyI6ICJTZXJ2ZXIgQiIsCiAgImFkZCI6ICIxMDQuMjEuMzAuMjI0IiwKICAicG9ydCI6ICI0NDMiLAogICJpZCI6ICJhNmY4YzNhMS02OWE0LTRjN2UtOGFkNi0xYjdhMmQ3ZjliNGMiLAogICJhaWQiOiAiMCIsCiAgInNjeSI6ICJhdXRvIiwKICAibmV0IjogIndzIiwKICAidHlwZSI6ICJub25lIiwKICAiaG9zdCI6ICJzZXJ2ZXItYi50YWJhdGVsZWNvbS5kZXYiLAogICJwYXRoIjogIi8iLAogICJ0bHMiOiAidGxzIiwKICAic25pIjogIiIsCiAgImFscG4iOiAiIiwKICAiZnAiOiAiIgp9
```

## Configuration Output

The application generates configurations compatible with:
- **sing-box**: Modern proxy client with enhanced features
- **V2Ray**: Traditional V2Ray client (function available but not used in TUI)

## Project Structure

```
go_v2ray_client/
├── main.go          # Main application entry point
├── tui.go           # Terminal User Interface implementation
├── vmess_parser.go  # VMess link parsing and configuration generation
├── config.json      # Generated configuration file
├── go.mod           # Go module dependencies
└── README.md        # This file
```

## Dependencies

- `github.com/rivo/tview`: Terminal UI library
- `github.com/gdamore/tcell/v2`: Terminal cell library (dependency of tview)

## Building

```bash
# Build for current platform
go build

# Build for specific platform
GOOS=linux GOARCH=amd64 go build
GOOS=windows GOARCH=amd64 go build
GOOS=darwin GOARCH=amd64 go build
```

## License

This project is open source and available under the MIT License. 