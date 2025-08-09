package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

// exportConfig exports the current configuration to a default file
func (tui *TUI) exportConfig() {
	if !tui.hasConfigToExport() {
		tui.updateStatus("Error: No configuration to export", tcell.ColorRed)
		return
	}
	tui.showFileDialog(tui.configText.GetText(true))
}

// showFileDialog switches to file explorer for export location
func (tui *TUI) showFileDialog(configText string) {
	// TODO: implement actual file picker selection, currently just shows explorer
	tui.app.SetRoot(tui.fileExplorer, true)
}

// performExport writes the config to a fixed filename in current working directory
func (tui *TUI) performExport() {
	if !tui.hasConfigToExport() {
		tui.updateStatus("Error: No configuration to export", tcell.ColorRed)
		return
	}

	filename := "exported_config.json"
	if err := tui.writeConfigToFile(filename); err != nil {
		tui.updateStatus(fmt.Sprintf("Error exporting file: %v", err), tcell.ColorRed)
		return
	}

	tui.updateStatus(fmt.Sprintf("Configuration exported to %s successfully!", filename), tcell.ColorGreen)
}

// loadDirectory lists files and folders in a path for file explorer
func (tui *TUI) loadDirectory(path string) {
	tui.fileList.Clear()
	tui.currentPath = path
	tui.pathInput.SetText(path)

	// Parent dir
	tui.fileList.AddItem(".. (Parent Directory)", "", 0, func() {
		parent := filepath.Dir(path)
		if parent != path {
			tui.updateStatus(fmt.Sprintf("Navigating to parent directory: %s", parent), tcell.ColorBlue)
			tui.loadDirectory(parent)
		}
	})

	files, err := os.ReadDir(path)
	if err != nil {
		tui.fileList.AddItem("Error reading directory", err.Error(), 0, nil)
		return
	}

	// Directories
	for _, file := range files {
		if file.IsDir() {
			dirName := file.Name()
			tui.fileList.AddItem("📁 "+dirName, "Directory", 0, func() {
				newPath := filepath.Join(path, dirName)
				tui.updateStatus(fmt.Sprintf("Navigating to directory: %s", newPath), tcell.ColorBlue)
				tui.loadDirectory(newPath)
			})
		}
	}

	// Files
	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			label := "📄 " + fileName
			desc := "File"
			if strings.HasSuffix(fileName, ".json") {
				desc = "JSON file"
			}
			tui.fileList.AddItem(label, desc, 0, nil)
		}
	}
}

// exportToCurrentPath saves the config to the currently opened directory
func (tui *TUI) exportToCurrentPath() {
	if !tui.hasConfigToExport() {
		tui.updateStatus("Error: No configuration to export", tcell.ColorRed)
		return
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(tui.currentPath, fmt.Sprintf("v2ray_config_%s.json", timestamp))

	if err := tui.writeConfigToFile(filename); err != nil {
		tui.updateStatus(fmt.Sprintf("Error exporting file: %v", err), tcell.ColorRed)
		return
	}

	tui.updateStatus(fmt.Sprintf("Configuration exported to %s successfully!", filename), tcell.ColorGreen)
	tui.app.SetRoot(tui.mainFlex, true)
}

// hasConfigToExport checks if there's an actual config to save
func (tui *TUI) hasConfigToExport() bool {
	text := tui.configText.GetText(true)
	return text != "" && text != "Configuration will appear here..."
}

// writeConfigToFile handles saving the current config text to a file
func (tui *TUI) writeConfigToFile(path string) error {
	return os.WriteFile(path, []byte(tui.configText.GetText(true)), 0644)
}
