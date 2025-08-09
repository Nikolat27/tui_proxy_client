package tui

import (
	"github.com/rivo/tview"
)

// Config represents a single configuration entry
type Config struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Protocol  string `json:"protocol"`
	Link      string `json:"link"`
	CreatedAt string `json:"created_at"`
	LastUsed  string `json:"last_used"`
}

// ConfigStorage represents the configuration storage structure
type ConfigStorage struct {
	Configurations []Config `json:"configurations"`
	Metadata       struct {
		Version      string `json:"version"`
		TotalConfigs int    `json:"total_configs"`
		LastUpdated  string `json:"last_updated"`
	} `json:"metadata"`
}

// TUI represents the terminal user interface
type TUI struct {
	app              *tview.Application
	mainFlex         *tview.Flex
	vmessInput       *tview.InputField
	statusText       *tview.TextView
	configText       *tview.TextView
	buttons          *tview.Flex
	configList       *tview.List
	fileDialog       *tview.Modal
	fileExplorer     *tview.Flex
	fileList         *tview.List
	pathInput        *tview.InputField
	currentPath      string
	configs          ConfigStorage
	isConnected      bool
	clientType       string
	connectedConfig  string
	connectionStatus *tview.TextView
}

// UIComponents holds references to UI elements for easier access
type UIComponents struct {
	Title            *tview.TextView
	ConnectionStatus *tview.TextView
	VMessInput       *tview.InputField
	StatusText       *tview.TextView
	ConfigText       *tview.TextView
	ConfigList       *tview.List
	Buttons          *tview.Flex
	FileExplorer     *tview.Flex
	FileList         *tview.List
	PathInput        *tview.InputField
}
