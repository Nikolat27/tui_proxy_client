package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

// clearUI clears the input and config display
func (tui *TUI) clearUI() {
	tui.vmessInput.SetText("")
	tui.configText.SetText("Configuration will appear here...")
	tui.updateStatus("UI cleared", tcell.ColorBlue)
}

// updateStatus updates the status text with a message and color
func (tui *TUI) updateStatus(message string, color tcell.Color) {
	if tui.isConnected {
		statusMessage := fmt.Sprintf("%s | Connected: %s (%s)", message, tui.connectedConfig, tui.clientType)
		tui.statusText.SetText(statusMessage)
	} else {
		tui.statusText.SetText(message)
	}
	tui.statusText.SetTextColor(color)
}

// updateConnectionStatus updates the connection status display
func (tui *TUI) updateConnectionStatus() {
	portInUse := tui.isPort1080InUse()

	if tui.isConnected && portInUse {
		tui.connectionStatus.SetText(fmt.Sprintf("Status: Connected to %s (%s) - Port 1080 Active", tui.connectedConfig, tui.clientType))
		tui.connectionStatus.SetTextColor(tcell.ColorGreen)
	} else if tui.isConnected && !portInUse {
		tui.isConnected = false
		tui.clientType = ""
		tui.connectedConfig = ""
		tui.connectionStatus.SetText("Status: Not Connected (Port 1080 Free)")
		tui.connectionStatus.SetTextColor(tcell.ColorRed)
	} else {
		tui.connectionStatus.SetText("Status: Not Connected (Port 1080 Free)")
		tui.connectionStatus.SetTextColor(tcell.ColorRed)
	}
}

// parseVMess parses the VMess link and displays the configuration
func (tui *TUI) parseVMess() {
	vmessLink := tui.vmessInput.GetText()
	if vmessLink == "" {
		tui.updateStatus("Error: Please enter a VMess link", tcell.ColorRed)
		return
	}

	vmessLink = strings.TrimSpace(vmessLink)
	if vmessLink == "" {
		tui.updateStatus("Error: Please enter a VMess link", tcell.ColorRed)
		return
	}

	if !strings.HasPrefix(vmessLink, "vmess://") {
		tui.updateStatus("Error: Invalid VMess link format. Must start with 'vmess://'", tcell.ColorRed)
		return
	}

	config, err := VMessToSingBox(vmessLink)
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error: %v", err), tcell.ColorRed)
		return
	}

	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error marshaling config: %v", err), tcell.ColorRed)
		return
	}

	tui.configText.SetText(string(configJSON))
	tui.updateStatus("VMess configuration parsed successfully!", tcell.ColorGreen)
}

// handlePaste handles paste operations by reading from clipboard
func (tui *TUI) handlePaste() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xclip", "-o", "-selection", "clipboard")
	case "darwin":
		cmd = exec.Command("pbpaste")
	default:
		tui.updateStatus("Paste not supported on this system. Use right-click or manual typing.", tcell.ColorYellow)
		return
	}

	output, err := cmd.Output()
	if err != nil {
		tui.updateStatus("Failed to read clipboard. Use right-click or manual typing.", tcell.ColorYellow)
		return
	}

	pastedText := strings.TrimSpace(string(output))
	if pastedText != "" {
		tui.vmessInput.SetText(pastedText)
		tui.updateStatus("Text pasted successfully!", tcell.ColorGreen)
	}
}

// isPort1080InUse checks if port 1080 is currently being used
func (tui *TUI) isPort1080InUse() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	cmd := exec.CommandContext(ctx, "sh", "-c", "netstat -tlnp 2>/dev/null | grep :1080")
	output, err := cmd.Output()
	cancel()
	if err == nil && len(strings.TrimSpace(string(output))) > 0 {
		return true
	}

	altCtx, altCancel := context.WithTimeout(context.Background(), 3*time.Second)
	cmd = exec.CommandContext(altCtx, "sh", "-c", "lsof -ti:1080 2>/dev/null")
	output, err = cmd.Output()
	altCancel()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(output))) > 0
}

// periodicStatusCheck periodically checks the connection status and updates the UI
func (tui *TUI) periodicStatusCheck() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 24*time.Hour)
	defer cancel()

	for {
		select {
		case <-ticker.C:
			if tui.app != nil {
				tui.app.QueueUpdateDraw(func() {
					if tui.connectionStatus != nil {
						tui.updateConnectionStatus()
					}
				})
			}
		case <-ctx.Done():
			return
		}
	}
} 