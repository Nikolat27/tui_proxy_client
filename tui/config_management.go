package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// addConfig adds the current configuration to the saved list
func (tui *TUI) addConfig() {
	vmessLink := strings.TrimSpace(tui.vmessInput.GetText())
	if vmessLink == "" {
		tui.updateStatus("Error: No VMess link to add. Please paste your VMess link first.", tcell.ColorRed)
		return
	}

	if !strings.HasPrefix(vmessLink, "vmess://") {
		tui.updateStatus("Error: Invalid VMess link format. Must start with 'vmess://'", tcell.ColorRed)
		return
	}

	config, err := VMessToSingBox(vmessLink)
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error parsing VMess: %v", err), tcell.ColorRed)
		return
	}

	if configJSON, err := json.MarshalIndent(config, "", "  "); err == nil {
		tui.configText.SetText(string(configJSON))
	} else {
		tui.updateStatus(fmt.Sprintf("Error serializing config: %v", err), tcell.ColorRed)
		return
	}

	configName := fmt.Sprintf("Config %d", len(tui.configs.Configurations)+1)

	newConfig := Config{
		ID:        fmt.Sprintf("%d", len(tui.configs.Configurations)+1),
		Name:      configName,
		Protocol:  "vmess",
		Link:      vmessLink,
		CreatedAt: time.Now().Format(time.RFC3339),
		LastUsed:  time.Now().Format(time.RFC3339),
	}

	tui.configs.Configurations = append(tui.configs.Configurations, newConfig)

	if err := tui.saveConfigsToFile(); err != nil {
		tui.updateStatus(fmt.Sprintf("Error saving config: %v", err), tcell.ColorRed)
		return
	}

	tui.refreshConfigList()
	tui.updateStatus(fmt.Sprintf("Configuration '%s' added and saved successfully", configName), tcell.ColorGreen)
}

// loadConfigList loads existing configurations from storage
func (tui *TUI) loadConfigList() {
	tui.loadConfigsFromFile()
	tui.refreshConfigList()

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
			localIndex := i
			displayName := fmt.Sprintf("%s (%s)", config.Name, config.Protocol)
			secondaryText := fmt.Sprintf("Created: %s | Last Used: %s",
				config.CreatedAt[:10], config.LastUsed[:10])

			tui.configList.AddItem(displayName, secondaryText, 0, func() {
				tui.viewConfig(localIndex)
			})
		}
	}

	tui.updateConfigCount()
}

// updateConfigCount updates the status to show the current number of configurations
func (tui *TUI) updateConfigCount() {
	if count := len(tui.configs.Configurations); count == 0 {
		tui.updateStatus("No configurations saved", tcell.ColorYellow)
	} else {
		tui.updateStatus(fmt.Sprintf("Loaded %d configuration(s) from configs.json", count), tcell.ColorGreen)
	}
}

// refreshConfigurations manually reloads configurations from configs.json
func (tui *TUI) refreshConfigurations() {
	tui.loadConfigsFromFile()
	tui.refreshConfigList()

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

	parsedConfig, err := VMessToSingBox(config.Link)
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error parsing config: %v", err), tcell.ColorRed)
		return
	}

	configJSON, err := json.MarshalIndent(parsedConfig, "", "  ")
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error marshaling config: %v", err), tcell.ColorRed)
		return
	}

	tui.configText.SetText(fmt.Sprintf(
		"Configuration: %s\nProtocol: %s\nLink: %s\nCreated: %s\nLast Used: %s\n\nParsed Configuration:\n%s\n\nReady to connect - click Connect button to start",
		config.Name, config.Protocol, config.Link, config.CreatedAt[:10], config.LastUsed[:10], string(configJSON)))

	tui.vmessInput.SetText(config.Link)
	tui.updateStatus(fmt.Sprintf("Selected configuration: %s (Ready to connect)", config.Name), tcell.ColorBlue)

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

	tui.configs.Configurations = append(
		tui.configs.Configurations[:currentIndex],
		tui.configs.Configurations[currentIndex+1:]...,
	)

	if err := tui.saveConfigsToFile(); err != nil {
		tui.updateStatus(fmt.Sprintf("Error saving after deletion: %v", err), tcell.ColorRed)
		return
	}

	tui.refreshConfigList()
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

	// Keep the type as *tview.InputField
	nameInput := tview.NewInputField()
	nameInput.SetLabel("New Name: ")
	nameInput.SetText(config.Name)
	nameInput.SetFieldWidth(30)
	nameInput.SetBorder(true)
	nameInput.SetTitle(" Rename Configuration ")

	renameModal := tview.NewModal().
		SetText("Press Enter to save, Esc to cancel").
		AddButtons([]string{"Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			tui.app.SetRoot(tui.mainFlex, true)
		})

	// Handle Enter key
	nameInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			newName := strings.TrimSpace(nameInput.GetText())
			if newName == "" {
				tui.updateStatus("Error: Name cannot be empty", tcell.ColorRed)
				tui.app.SetRoot(tui.mainFlex, true)
				return
			}

			tui.configs.Configurations[currentIndex].Name = newName

			if err := tui.saveConfigsToFile(); err != nil {
				tui.updateStatus(fmt.Sprintf("Error saving after rename: %v", err), tcell.ColorRed)
				tui.app.SetRoot(tui.mainFlex, true)
				return
			}

			tui.refreshConfigList()
			tui.updateStatus(fmt.Sprintf("Configuration renamed to '%s' successfully", newName), tcell.ColorGreen)
			tui.app.SetRoot(tui.mainFlex, true)
		}
	})

	renameFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(renameModal, 0, 1, false).
		AddItem(nameInput, 3, 0, true)

	tui.app.SetRoot(renameFlex, true)
	tui.app.SetFocus(nameInput)
}

// loadConfigsFromFile loads configurations from configs.json
func (tui *TUI) loadConfigsFromFile() {
	data, err := os.ReadFile("configs.json")
	if err != nil {
		tui.configs = ConfigStorage{
			Configurations: []Config{},
			Metadata: struct {
				Version      string `json:"version"`
				TotalConfigs int    `json:"total_configs"`
				LastUpdated  string `json:"last_updated"`
			}{"1.0", 0, time.Now().Format(time.RFC3339)},
		}
		return
	}

	if err := json.Unmarshal(data, &tui.configs); err != nil {
		tui.updateStatus("Error loading configs.json", tcell.ColorRed)
		tui.configs = ConfigStorage{
			Configurations: []Config{},
			Metadata: struct {
				Version      string `json:"version"`
				TotalConfigs int    `json:"total_configs"`
				LastUpdated  string `json:"last_updated"`
			}{"1.0", 0, time.Now().Format(time.RFC3339)},
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
