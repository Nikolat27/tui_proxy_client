package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

// exportConfig exports the current configuration to a user-selected file
func (tui *TUI) exportConfig() {
	configText := tui.configText.GetText(true)
	if configText == "" || configText == "Configuration will appear here..." {
		tui.updateStatus("Error: No configuration to export", tcell.ColorRed)
		return
	}

	tui.showFileDialog(configText)
}

// showFileDialog displays the file export dialog
func (tui *TUI) showFileDialog(configText string) {
	tui.app.SetRoot(tui.fileExplorer, true)
}

// performExport performs the actual file export
func (tui *TUI) performExport() {
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

	tui.fileList.AddItem(".. (Parent Directory)", "", 0, func() {
		parent := filepath.Dir(path)
		fmt.Printf("Current path: %s, Parent path: %s\n", path, parent)
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

	for _, file := range files {
		if file.IsDir() {
			dirName := file.Name()
			localDirName := dirName
			tui.fileList.AddItem("📁 "+dirName, "Directory", 0, func() {
				newPath := filepath.Join(path, localDirName)
				fmt.Printf("Navigating to: %s\n", newPath)
				tui.updateStatus(fmt.Sprintf("Navigating to directory: %s", newPath), tcell.ColorBlue)
				tui.loadDirectory(newPath)
			})
		}
	}

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
