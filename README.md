# V2Ray Client TUI

A Terminal User Interface (TUI) application for generating V2Ray client configurations from VMess links with built-in configuration management.

## Features

- **Interactive TUI**: Beautiful terminal-based user interface using the `tview` library
- **VMess Parser**: Converts VMess links to both V2Ray and sing-box configurations
- **Configuration Management**: Store and manage multiple configurations in `configs.json`
- **Real-time Preview**: See the generated configuration before saving
- **Easy Configuration**: Simple input field for VMess links
- **File Export**: Save configurations to custom locations
- **Configuration Operations**: Add, delete, rename, and view saved configurations
- **Auto-save**: Configurations are automatically saved to `configs.json`

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
2. **Parse Configuration**: Press `Enter` in the input field to parse the VMess link
3. **Preview**: View the generated configuration in the preview area
4. **Add Configuration**: Click "Add Config" or press `Ctrl+A` to save to `configs.json`
5. **Manage Configurations**: Use the configuration list on the right to view, rename, or delete saved configs
6. **Export**: Click "Export Config" or press `Ctrl+S` to export to a custom location
7. **Refresh**: Click "Refresh" or press `Ctrl+F` to reload configurations from `configs.json`
8. **Clear**: Use "Clear" button or press `Ctrl+L` to reset the interface
9. **Quit**: Click "Quit" or press `Ctrl+C` to exit

### Keyboard Shortcuts

- `Enter`: Parse VMess link (when in input field)
- `Ctrl+A`: Add configuration to `configs.json`
- `Ctrl+S`: Export configuration to file
- `Ctrl+D`: Delete selected configuration
- `Ctrl+R`: Rename selected configuration
- `Ctrl+F`: Refresh configurations from `configs.json`
- `Ctrl+L`: Clear interface
- `Ctrl+C`: Quit application
- `Tab`: Navigate between elements
- Arrow keys: Navigate and scroll

### Configuration Management

The application automatically manages configurations in `configs.json`:

- **Add Config**: Automatically parses and saves VMess configurations
- **View Config**: Click on any saved configuration to view its details
- **Rename Config**: Give your configurations meaningful names
- **Delete Config**: Remove unwanted configurations
- **Auto-save**: All changes are automatically saved to `configs.json`

### Example VMess Link

```
vmess://ewogICJ2IjogIjIiLAogICJwcyI6ICJTZXJ2ZXIgQiIsCiAgImFkZCI6ICIxMDQuMjEuMzAuMjI0IiwKICAicG9ydCI6ICI0NDMiLAogICJpZCI6ICJhNmY4YzNhMS02OWE0LTRjN2UtOGFkNi0xYjdhMmQ3ZjliNGMiLAogICJhaWQiOiAiMCIsCiAgInNjeSI6ICJhdXRvIiwKICAibmV0IjogIndzIiwKICAidHlwZSI6ICJub25lIiwKICAiaG9zdCI6ICJzZXJ2ZXItYi50YWJhdGVsZWNvbS5kZXYiLAogICJwYXRoIjogIi8iLAogICJ0bHMiOiAidGxzIiwKICAic25pIjogIiIsCiAgImFscG4iOiAiIiwKICAiZnAiOiAiIgp9
```

## Configuration Files

### configs.json
Stores all your saved configurations with metadata:
```json
{
    "configurations": [
        {
            "id": "1",
            "name": "Example VMess Server",
            "protocol": "vmess",
            "link": "vmess://...",
            "created_at": "2024-12-20T14:30:00Z",
            "last_used": "2024-12-20T14:30:00Z"
        }
    ],
    "metadata": {
        "version": "1.0",
        "total_configs": 1,
        "last_updated": "2024-12-20T14:30:00Z"
    }
}
```

### Exported Configurations
Individual configuration files exported to your chosen location.

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
├── configs.json     # Configuration storage file
├── config.json      # Generated configuration file (legacy)
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