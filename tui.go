package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
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
	app          *tview.Application
	mainFlex     *tview.Flex
	vmessInput   *tview.InputField
	statusText   *tview.TextView
	configText   *tview.TextView
	buttons      *tview.Flex
	configList   *tview.List
	fileDialog   *tview.Modal
	fileExplorer *tview.Flex
	fileList     *tview.List
	pathInput    *tview.InputField
	currentPath  string
	configs      ConfigStorage
}

// NewTUI creates a new TUI instance
func NewTUI() *TUI {
	tui := &TUI{
		app: tview.NewApplication(),
	}

	// Enable mouse support
	tui.app.EnableMouse(true)

	tui.setupUI()
	tui.setupKeybindings()

	return tui
}

// setupUI initializes the user interface components
func (tui *TUI) setupUI() {
	// Create main title
	title := tview.NewTextView().
		SetText("V2Ray Client Configuration Generator").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorYellow).
		SetDynamicColors(true)

	// Create VMess input field
	tui.vmessInput = tview.NewInputField()
	tui.vmessInput.SetLabel("VMess Link: ")
	tui.vmessInput.SetPlaceholder("vmess://...")
	tui.vmessInput.SetFieldWidth(80)
	tui.vmessInput.SetBorder(true)
	tui.vmessInput.SetTitle(" Enter VMess Configuration ")

	// Create status text
	tui.statusText = tview.NewTextView()
	tui.statusText.SetText("Ready to parse VMess configuration")
	tui.statusText.SetTextAlign(tview.AlignCenter)
	tui.statusText.SetTextColor(tcell.ColorGreen)
	tui.statusText.SetBorder(true)
	tui.statusText.SetTitle(" Status ")

	// Create config preview area
	tui.configText = tview.NewTextView()
	tui.configText.SetText("Configuration will appear here...")
	tui.configText.SetTextAlign(tview.AlignLeft)
	tui.configText.SetTextColor(tcell.ColorWhite)
	tui.configText.SetBorder(true)
	tui.configText.SetTitle(" Generated Configuration ")
	tui.configText.SetScrollable(true)

	// Create config list
	tui.configList = tview.NewList()
	tui.configList.SetBorder(true)
	tui.configList.SetTitle(" Saved Configurations ")
	tui.configList.SetMainTextColor(tcell.ColorWhite)
	tui.loadConfigList() // Load existing configurations

	// Create buttons with shortcut hints
	addConfigBtn := tview.NewButton("Add Config (Ctrl+A)").
		SetSelectedFunc(func() {
			tui.addConfig()
		})

	exportBtn := tview.NewButton("Export Config (Ctrl+S)").
		SetSelectedFunc(func() {
			tui.exportConfig()
		})

	deleteBtn := tview.NewButton("Delete Config (Ctrl+D)").
		SetSelectedFunc(func() {
			tui.deleteSelectedConfig()
		})

	renameBtn := tview.NewButton("Rename Config (Ctrl+R)").
		SetSelectedFunc(func() {
			tui.renameSelectedConfig()
		})

	refreshBtn := tview.NewButton("Refresh (Ctrl+F)").
		SetSelectedFunc(func() {
			tui.refreshConfigurations()
		})

	clearBtn := tview.NewButton("Clear (Ctrl+L)").
		SetSelectedFunc(func() {
			tui.clearUI()
		})

	quitBtn := tview.NewButton("Quit (Ctrl+C)").
		SetSelectedFunc(func() {
			tui.app.Stop()
		})

	// Arrange buttons horizontally
	tui.buttons = tview.NewFlex().
		AddItem(addConfigBtn, 0, 1, false).
		AddItem(exportBtn, 0, 1, false).
		AddItem(deleteBtn, 0, 1, false).
		AddItem(renameBtn, 0, 1, false).
		AddItem(refreshBtn, 0, 1, false).
		AddItem(clearBtn, 0, 1, false).
		AddItem(quitBtn, 0, 1, false)

	// Create main layout with config list
	configSection := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tui.configText, 0, 2, false).
		AddItem(tui.configList, 0, 1, false)

	tui.mainFlex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 3, 0, false).
		AddItem(tui.vmessInput, 3, 0, false).
		AddItem(tui.statusText, 3, 0, false).
		AddItem(configSection, 0, 1, false).
		AddItem(tui.buttons, 3, 0, false)

	// Initialize file dialog
	tui.fileDialog = tview.NewModal().
		SetText("Enter filename to export:").
		AddButtons([]string{"Export", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Export" {
				tui.performExport()
			}
			tui.app.SetRoot(tui.mainFlex, true)
		})

	// Initialize file explorer
	tui.currentPath, _ = os.Getwd() // Get current working directory
	tui.pathInput = tview.NewInputField()
	tui.pathInput.SetLabel("Path: ")
	tui.pathInput.SetText(tui.currentPath)
	tui.pathInput.SetFieldWidth(50)
	tui.pathInput.SetBorder(true)
	tui.pathInput.SetTitle(" Current Directory ")

	tui.fileList = tview.NewList()
	tui.fileList.SetBorder(true)
	tui.fileList.SetTitle(" Files and Folders ")

	// Create file explorer layout
	explorerButtons := tview.NewFlex().
		AddItem(tview.NewButton("Export Here").SetSelectedFunc(func() {
			tui.exportToCurrentPath()
		}), 0, 1, false).
		AddItem(tview.NewButton("Back to Main").SetSelectedFunc(func() {
			tui.app.SetRoot(tui.mainFlex, true)
		}), 0, 1, false)

	tui.fileExplorer = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tui.pathInput, 3, 0, false).
		AddItem(tui.fileList, 0, 1, false).
		AddItem(explorerButtons, 3, 0, false)

	// Load initial directory
	tui.loadDirectory(tui.currentPath)

	tui.app.SetRoot(tui.mainFlex, true)
}

// setupKeybindings sets up keyboard shortcuts
func (tui *TUI) setupKeybindings() {
	tui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlA:
			tui.addConfig()
			return nil
		case tcell.KeyCtrlS:
			tui.exportConfig()
			return nil
		case tcell.KeyCtrlD:
			tui.deleteSelectedConfig()
			return nil
		case tcell.KeyCtrlR:
			tui.renameSelectedConfig()
			return nil
		case tcell.KeyCtrlF:
			tui.refreshConfigurations()
			return nil
		case tcell.KeyCtrlL:
			tui.clearUI()
			return nil
		case tcell.KeyCtrlC:
			tui.app.Stop()
			return nil
		}
		return event
	})

	// Set up custom paste handling for the VMess input field
	tui.vmessInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlV {
			// Handle paste operation
			tui.handlePaste()
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			// Parse VMess link when Enter is pressed
			tui.parseVMess()
			return nil
		}
		return event
	})
}

// parseVMess parses the VMess link and displays the configuration
func (tui *TUI) parseVMess() {
	vmessLink := tui.vmessInput.GetText()
	if vmessLink == "" {
		tui.updateStatus("Error: Please enter a VMess link", tcell.ColorRed)
		return
	}

	// Trim whitespace
	vmessLink = strings.TrimSpace(vmessLink)
	if vmessLink == "" {
		tui.updateStatus("Error: Please enter a VMess link", tcell.ColorRed)
		return
	}

	// Validate that it starts with vmess://
	if !strings.HasPrefix(vmessLink, "vmess://") {
		tui.updateStatus("Error: Invalid VMess link format. Must start with 'vmess://'", tcell.ColorRed)
		return
	}

	config, err := VMessToSingBox(vmessLink)
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error: %v", err), tcell.ColorRed)
		return
	}

	// Convert config to JSON string
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error marshaling config: %v", err), tcell.ColorRed)
		return
	}

	tui.configText.SetText(string(configJSON))
	tui.updateStatus("VMess configuration parsed successfully!", tcell.ColorGreen)
}

// exportConfig exports the current configuration to a user-selected file
func (tui *TUI) exportConfig() {
	// Get the current configuration text
	configText := tui.configText.GetText(true)
	if configText == "" || configText == "Configuration will appear here..." {
		tui.updateStatus("Error: No configuration to export", tcell.ColorRed)
		return
	}

	// Create file dialog UI
	tui.showFileDialog(configText)
}

// showFileDialog displays the file export dialog
func (tui *TUI) showFileDialog(configText string) {
	// Show the file explorer instead of simple modal
	tui.app.SetRoot(tui.fileExplorer, true)
}

// performExport performs the actual file export
func (tui *TUI) performExport() {
	// For now, export to a default filename
	// In a full implementation, you'd get the filename from user input
	filename := "exported_config.json"

	configText := tui.configText.GetText(true)
	if configText == "" || configText == "Configuration will appear here..." {
		tui.updateStatus("Error: No configuration to export", tcell.ColorRed)
		return
	}

	err := os.WriteFile(filename, []byte(configText), 0644)
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error exporting file: %v", err), tcell.ColorRed)
		return
	}

	tui.updateStatus(fmt.Sprintf("Configuration exported to %s successfully!", filename), tcell.ColorGreen)
}

// loadDirectory loads the contents of a directory into the file list
func (tui *TUI) loadDirectory(path string) {
	tui.fileList.Clear()
	tui.currentPath = path
	tui.pathInput.SetText(path)

	// Add parent directory option
	tui.fileList.AddItem(".. (Parent Directory)", "", 0, func() {
		parent := filepath.Dir(path)
		// Debug: print the paths to see what's happening
		fmt.Printf("Current path: %s, Parent path: %s\n", path, parent)
		if parent != path {
			tui.updateStatus(fmt.Sprintf("Navigating to parent directory: %s", parent), tcell.ColorBlue)
			tui.loadDirectory(parent)
		}
	})

	// Read directory contents
	files, err := os.ReadDir(path)
	if err != nil {
		tui.fileList.AddItem("Error reading directory", err.Error(), 0, nil)
		return
	}

	// Add directories first
	for _, file := range files {
		if file.IsDir() {
			dirName := file.Name()
			// Create a local copy to avoid closure issues
			localDirName := dirName
			tui.fileList.AddItem("📁 "+dirName, "Directory", 0, func() {
				newPath := filepath.Join(path, localDirName)
				fmt.Printf("Navigating to: %s\n", newPath)
				tui.updateStatus(fmt.Sprintf("Navigating to directory: %s", newPath), tcell.ColorBlue)
				tui.loadDirectory(newPath)
			})
		}
	}

	// Add files
	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			if strings.HasSuffix(fileName, ".json") {
				tui.fileList.AddItem("📄 "+fileName, "JSON file", 0, func() {
					// Could implement file selection here
				})
			} else {
				tui.fileList.AddItem("📄 "+fileName, "File", 0, func() {
					// Could implement file selection here
				})
			}
		}
	}
}

// exportToCurrentPath exports the configuration to the current directory
func (tui *TUI) exportToCurrentPath() {
	configText := tui.configText.GetText(true)
	if configText == "" || configText == "Configuration will appear here..." {
		tui.updateStatus("Error: No configuration to export", tcell.ColorRed)
		return
	}

	// Generate filename based on current time
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(tui.currentPath, fmt.Sprintf("v2ray_config_%s.json", timestamp))

	err := os.WriteFile(filename, []byte(configText), 0644)
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error exporting file: %v", err), tcell.ColorRed)
		return
	}

	tui.updateStatus(fmt.Sprintf("Configuration exported to %s successfully!", filename), tcell.ColorGreen)
	tui.app.SetRoot(tui.mainFlex, true)
}

// clearUI clears the input and config display
func (tui *TUI) clearUI() {
	tui.vmessInput.SetText("")
	tui.configText.SetText("Configuration will appear here...")
	tui.updateStatus("UI cleared", tcell.ColorBlue)
}

// updateStatus updates the status text with a message and color
func (tui *TUI) updateStatus(message string, color tcell.Color) {
	tui.statusText.SetText(message)
	tui.statusText.SetTextColor(color)
}

// handlePaste handles paste operations by reading from clipboard
func (tui *TUI) handlePaste() {
	// Try to read from clipboard using xclip (Linux) or pbcopy (macOS)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xclip", "-o", "-selection", "clipboard")
	case "darwin":
		cmd = exec.Command("pbpaste")
	default:
		// For Windows or other systems, just show a message
		tui.updateStatus("Paste not supported on this system. Use right-click or manual typing.", tcell.ColorYellow)
		return
	}

	output, err := cmd.Output()
	if err != nil {
		tui.updateStatus("Failed to read clipboard. Use right-click or manual typing.", tcell.ColorYellow)
		return
	}

	// Set the pasted text to the input field
	pastedText := strings.TrimSpace(string(output))
	if pastedText != "" {
		tui.vmessInput.SetText(pastedText)
		tui.updateStatus("Text pasted successfully!", tcell.ColorGreen)
	}
}

// addConfig adds the current configuration to the saved list
func (tui *TUI) addConfig() {
	vmessLink := tui.vmessInput.GetText()
	if vmessLink == "" {
		tui.updateStatus("Error: No VMess link to add", tcell.ColorRed)
		return
	}

	// Trim whitespace
	vmessLink = strings.TrimSpace(vmessLink)
	if vmessLink == "" {
		tui.updateStatus("Error: No VMess link to add", tcell.ColorRed)
		return
	}

	// Validate that it starts with vmess://
	if !strings.HasPrefix(vmessLink, "vmess://") {
		tui.updateStatus("Error: Invalid VMess link format. Must start with 'vmess://'", tcell.ColorRed)
		return
	}

	// Parse the VMess link behind the scenes
	config, err := VMessToSingBox(vmessLink)
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error parsing VMess: %v", err), tcell.ColorRed)
		return
	}

	// Display the parsed configuration
	configJSON, _ := json.MarshalIndent(config, "", "  ")
	tui.configText.SetText(string(configJSON))

	// Generate a name for the config
	configName := fmt.Sprintf("Config %d", len(tui.configs.Configurations)+1)

	// Create new config entry
	newConfig := Config{
		ID:        fmt.Sprintf("%d", len(tui.configs.Configurations)+1),
		Name:      configName,
		Protocol:  "vmess",
		Link:      vmessLink,
		CreatedAt: time.Now().Format(time.RFC3339),
		LastUsed:  time.Now().Format(time.RFC3339),
	}

	// Add to configs storage
	tui.configs.Configurations = append(tui.configs.Configurations, newConfig)

	// Save to file
	err = tui.saveConfigsToFile()
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error saving config: %v", err), tcell.ColorRed)
		return
	}

	// Refresh the config list
	tui.refreshConfigList()

	tui.updateStatus(fmt.Sprintf("Configuration '%s' added, parsed, and saved successfully", configName), tcell.ColorGreen)
}

// loadConfigList loads existing configurations from storage
func (tui *TUI) loadConfigList() {
	// Load configurations from configs.json
	tui.loadConfigsFromFile()

	// Refresh the config list display
	tui.refreshConfigList()

	// Show initial status
	if len(tui.configs.Configurations) == 0 {
		tui.updateStatus("Ready to add configurations. Use Ctrl+A to add a new VMess configuration.", tcell.ColorBlue)
	} else {
		tui.updateStatus(fmt.Sprintf("Ready! %d configuration(s) loaded from configs.json", len(tui.configs.Configurations)), tcell.ColorGreen)
	}
}

// refreshConfigList refreshes the configuration list to show the latest saved configs
func (tui *TUI) refreshConfigList() {
	tui.configList.Clear()
	if len(tui.configs.Configurations) == 0 {
		tui.configList.AddItem("No configurations yet", "Add your first configuration", 0, nil)
	} else {
		for i, config := range tui.configs.Configurations {
			// Create a local copy to avoid closure issues
			localIndex := i

			// Format the display name with more details
			displayName := fmt.Sprintf("%s (%s)", config.Name, config.Protocol)

			// Format the secondary text with creation date and last used
			secondaryText := fmt.Sprintf("Created: %s | Last Used: %s",
				config.CreatedAt[:10], config.LastUsed[:10])

			tui.configList.AddItem(displayName, secondaryText, 0, func() {
				tui.viewConfig(localIndex)
			})
		}
	}

	// Update status with configuration count
	tui.updateConfigCount()
}

// updateConfigCount updates the status to show the current number of configurations
func (tui *TUI) updateConfigCount() {
	count := len(tui.configs.Configurations)
	if count == 0 {
		tui.updateStatus("No configurations saved", tcell.ColorYellow)
	} else {
		tui.updateStatus(fmt.Sprintf("Loaded %d configuration(s) from configs.json", count), tcell.ColorGreen)
	}
}

// refreshConfigurations manually reloads configurations from configs.json
func (tui *TUI) refreshConfigurations() {
	// Reload from file
	tui.loadConfigsFromFile()

	// Refresh the display
	tui.refreshConfigList()

	// Show status
	if len(tui.configs.Configurations) == 0 {
		tui.updateStatus("No configurations found in configs.json", tcell.ColorYellow)
	} else {
		tui.updateStatus(fmt.Sprintf("Refreshed! %d configuration(s) loaded from configs.json", len(tui.configs.Configurations)), tcell.ColorGreen)
	}
}

// viewConfig displays the details of a selected configuration
func (tui *TUI) viewConfig(configIndex int) {
	if configIndex < 0 || configIndex >= len(tui.configs.Configurations) {
		tui.updateStatus("Invalid configuration index", tcell.ColorRed)
		return
	}

	config := tui.configs.Configurations[configIndex]

	// Parse the VMess link to show the parsed configuration
	parsedConfig, err := VMessToSingBox(config.Link)
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error parsing config: %v", err), tcell.ColorRed)
		return
	}

	// Convert config to JSON string
	configJSON, err := json.MarshalIndent(parsedConfig, "", "  ")
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error marshaling config: %v", err), tcell.ColorRed)
		return
	}

	// Display the configuration
	tui.configText.SetText(string(configJSON))

	// Update the VMess input field to show the selected config
	tui.vmessInput.SetText(config.Link)

	// Update status
	tui.updateStatus(fmt.Sprintf("Viewing configuration: %s", config.Name), tcell.ColorBlue)

	// Update last used timestamp
	tui.configs.Configurations[configIndex].LastUsed = time.Now().Format(time.RFC3339)
	tui.saveConfigsToFile()
}

// deleteSelectedConfig deletes the currently selected configuration
func (tui *TUI) deleteSelectedConfig() {
	currentIndex := tui.configList.GetCurrentItem()

	if len(tui.configs.Configurations) == 0 {
		tui.updateStatus("No configurations to delete", tcell.ColorYellow)
		return
	}

	if currentIndex < 0 || currentIndex >= len(tui.configs.Configurations) {
		tui.updateStatus("Please select a configuration to delete", tcell.ColorYellow)
		return
	}

	config := tui.configs.Configurations[currentIndex]

	// Remove the configuration from the slice
	tui.configs.Configurations = append(tui.configs.Configurations[:currentIndex],
		tui.configs.Configurations[currentIndex+1:]...)

	// Save to file
	err := tui.saveConfigsToFile()
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error saving after deletion: %v", err), tcell.ColorRed)
		return
	}

	// Refresh the config list
	tui.refreshConfigList()

	// Clear the display if we deleted the currently viewed config
	tui.clearUI()

	tui.updateStatus(fmt.Sprintf("Configuration '%s' deleted successfully", config.Name), tcell.ColorGreen)
}

// renameSelectedConfig allows the user to rename the currently selected configuration
func (tui *TUI) renameSelectedConfig() {
	currentIndex := tui.configList.GetCurrentItem()

	if len(tui.configs.Configurations) == 0 {
		tui.updateStatus("No configurations to rename", tcell.ColorYellow)
		return
	}

	if currentIndex < 0 || currentIndex >= len(tui.configs.Configurations) {
		tui.updateStatus("Please select a configuration to rename", tcell.ColorYellow)
		return
	}

	config := tui.configs.Configurations[currentIndex]

	// Create input field for new name
	nameInput := tview.NewInputField()
	nameInput.SetLabel("New Name: ")
	nameInput.SetText(config.Name)
	nameInput.SetFieldWidth(30)
	nameInput.SetBorder(true)
	nameInput.SetTitle(" Rename Configuration ")

	// Create a simple modal for confirmation
	renameModal := tview.NewModal().
		SetText("Press Enter to save, Esc to cancel").
		AddButtons([]string{"Cancel"})

	// Set up input handling
	nameInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			newName := strings.TrimSpace(nameInput.GetText())
			if newName == "" {
				tui.updateStatus("Error: Name cannot be empty", tcell.ColorRed)
				tui.app.SetRoot(tui.mainFlex, true)
				return
			}

			// Update the configuration name
			tui.configs.Configurations[currentIndex].Name = newName

			// Save to file
			err := tui.saveConfigsToFile()
			if err != nil {
				tui.updateStatus(fmt.Sprintf("Error saving after rename: %v", err), tcell.ColorRed)
				tui.app.SetRoot(tui.mainFlex, true)
				return
			}

			// Refresh the config list
			tui.refreshConfigList()

			tui.updateStatus(fmt.Sprintf("Configuration renamed to '%s' successfully", newName), tcell.ColorGreen)
			tui.app.SetRoot(tui.mainFlex, true)
		}
	})

	// Set up modal button handling
	renameModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		tui.app.SetRoot(tui.mainFlex, true)
	})

	// Create layout for the rename interface
	renameFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(renameModal, 0, 1, false).
		AddItem(nameInput, 3, 0, false)

	// Show the rename interface
	tui.app.SetRoot(renameFlex, true)

	// Focus on the input field
	tui.app.SetFocus(nameInput)
}

// loadConfigsFromFile loads configurations from configs.json
func (tui *TUI) loadConfigsFromFile() {
	data, err := os.ReadFile("configs.json")
	if err != nil {
		// If file doesn't exist, create default structure
		tui.configs = ConfigStorage{
			Configurations: []Config{},
			Metadata: struct {
				Version      string `json:"version"`
				TotalConfigs int    `json:"total_configs"`
				LastUpdated  string `json:"last_updated"`
			}{
				Version:      "1.0",
				TotalConfigs: 0,
				LastUpdated:  time.Now().Format(time.RFC3339),
			},
		}
		return
	}

	err = json.Unmarshal(data, &tui.configs)
	if err != nil {
		tui.updateStatus("Error loading configs.json", tcell.ColorRed)
		tui.configs = ConfigStorage{
			Configurations: []Config{},
			Metadata: struct {
				Version      string `json:"version"`
				TotalConfigs int    `json:"total_configs"`
				LastUpdated  string `json:"last_updated"`
			}{
				Version:      "1.0",
				TotalConfigs: 0,
				LastUpdated:  time.Now().Format(time.RFC3339),
			},
		}
	}
}

// saveConfigsToFile saves configurations to configs.json
func (tui *TUI) saveConfigsToFile() error {
	tui.configs.Metadata.TotalConfigs = len(tui.configs.Configurations)
	tui.configs.Metadata.LastUpdated = time.Now().Format(time.RFC3339)

	data, err := json.MarshalIndent(tui.configs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("configs.json", data, 0644)
}

// Run starts the TUI application
func (tui *TUI) Run() error {
	return tui.app.Run()
}
