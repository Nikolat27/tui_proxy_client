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
	parser "go_v2ray_client/parser"
)

// clearUI resets the VMess input and config display
func (tui *TUI) clearUI() {
	tui.vmessInput.SetText("")
	tui.configText.SetText("Configuration will appear here...")
	tui.updateStatus("UI cleared", tcell.ColorBlue)
}

// updateStatus sets a status message with color, including connection info if active
func (tui *TUI) updateStatus(message string, color tcell.Color) {
	if tui.isConnected {
		message = fmt.Sprintf("%s | Connected: %s (%s)", message, tui.connectedConfig, tui.clientType)
	}
	tui.statusText.SetText(message).SetTextColor(color)
}

// updateConnectionStatus checks port 1080 and updates UI connection state
func (tui *TUI) updateConnectionStatus() {
	portInUse := tui.isPort1080InUse()

	switch {
	case tui.isConnected && portInUse:
		tui.connectionStatus.SetText(fmt.Sprintf("Status: Connected to %s (%s) - Port 1080 Active", tui.connectedConfig, tui.clientType)).
			SetTextColor(tcell.ColorGreen)

	case tui.isConnected && !portInUse:
		tui.isConnected = false
		tui.clientType, tui.connectedConfig = "", ""
		tui.connectionStatus.SetText("Status: Not Connected (Port 1080 Free)").
			SetTextColor(tcell.ColorRed)

	default:
		tui.connectionStatus.SetText("Status: Not Connected (Port 1080 Free)").
			SetTextColor(tcell.ColorRed)
	}
}

// parseVMess converts a VMess link into a SingBox JSON config and displays it
func (tui *TUI) parseVMess() {
	vmessLink := strings.TrimSpace(tui.vmessInput.GetText())
	if vmessLink == "" {
		tui.updateStatus("Error: Please enter a VMess link", tcell.ColorRed)
		return
	}

	if !strings.HasPrefix(vmessLink, "vmess://") {
		tui.updateStatus("Error: Invalid VMess link format. Must start with 'vmess://'", tcell.ColorRed)
		return
	}

	config, err := parser.VMessToSingBox(vmessLink)
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error parsing VMess: %v", err), tcell.ColorRed)
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

// handlePaste gets clipboard text depending on OS
func (tui *TUI) handlePaste() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xclip", "-o", "-selection", "clipboard")
	case "darwin":
		cmd = exec.Command("pbpaste")
	default:
		tui.updateStatus("Paste not supported on this system", tcell.ColorYellow)
		return
	}

	output, err := cmd.Output()
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Clipboard read failed: %v", err), tcell.ColorYellow)
		return
	}

	if pasted := strings.TrimSpace(string(output)); pasted != "" {
		tui.vmessInput.SetText(pasted)
		tui.updateStatus("Text pasted successfully!", tcell.ColorGreen)
	}
}

// isPort1080InUse checks if TCP port 1080 is active
func (tui *TUI) isPort1080InUse() bool {
	return tui.runPortCheck("lsof -ti:1080 2>/dev/null")
}

// runPortCheck executes a shell command to check port usage
func (tui *TUI) runPortCheck(cmdStr string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", cmdStr)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	portOutput := strings.TrimSpace(string(output))

	return len(portOutput) > 0
}

// periodicStatusCheck runs every 10s to update connection status
func (tui *TUI) periodicStatusCheck() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if tui.app == nil {
			continue
		}

		tui.app.QueueUpdateDraw(func() {
			if tui.connectionStatus != nil {
				tui.updateConnectionStatus()
			}
		})

	}
}
