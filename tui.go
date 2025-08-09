package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TUI represents the terminal user interface
type TUI struct {
	app        *tview.Application
	mainFlex   *tview.Flex
	vmessInput *tview.InputField
	statusText *tview.TextView
	configText *tview.TextView
	buttons    *tview.Flex
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

	// Create buttons with shortcut hints
	parseBtn := tview.NewButton("Parse VMess (Ctrl+P)").
		SetSelectedFunc(func() {
			tui.parseVMess()
		})

	saveBtn := tview.NewButton("Save Config (Ctrl+S)").
		SetSelectedFunc(func() {
			tui.saveConfig()
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
		AddItem(parseBtn, 0, 1, false).
		AddItem(saveBtn, 0, 1, false).
		AddItem(clearBtn, 0, 1, false).
		AddItem(quitBtn, 0, 1, false)

	// Create main layout
	tui.mainFlex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 3, 0, false).
		AddItem(tui.vmessInput, 3, 0, false).
		AddItem(tui.statusText, 3, 0, false).
		AddItem(tui.configText, 0, 1, false).
		AddItem(tui.buttons, 3, 0, false)

	tui.app.SetRoot(tui.mainFlex, true)
}

// setupKeybindings sets up keyboard shortcuts
func (tui *TUI) setupKeybindings() {
	tui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlP:
			tui.parseVMess()
			return nil
		case tcell.KeyCtrlS:
			tui.saveConfig()
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

// saveConfig saves the configuration to a file
func (tui *TUI) saveConfig() {
	vmessLink := tui.vmessInput.GetText()
	if vmessLink == "" {
		tui.updateStatus("Error: No VMess link to save", tcell.ColorRed)
		return
	}

	config, err := VMessToSingBox(vmessLink)
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error: %v", err), tcell.ColorRed)
		return
	}

	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error marshaling config: %v", err), tcell.ColorRed)
		return
	}

	if err := os.WriteFile("config.json", data, 0644); err != nil {
		tui.updateStatus(fmt.Sprintf("Error saving file: %v", err), tcell.ColorRed)
		return
	}

	tui.updateStatus("Configuration saved to config.json successfully!", tcell.ColorGreen)
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
	if runtime.GOOS == "linux" {
		cmd = exec.Command("xclip", "-o", "-selection", "clipboard")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("pbpaste")
	} else {
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

// Run starts the TUI application
func (tui *TUI) Run() error {
	return tui.app.Run()
}
