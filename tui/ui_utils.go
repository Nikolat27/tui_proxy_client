package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"go_v2ray_client/parser"

	"github.com/gdamore/tcell/v2"
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

// parseProxyLink converts a proxy link into a SingBox JSON config and displays it
func (tui *TUI) parseProxyLink() {
	proxyLink := strings.TrimSpace(tui.vmessInput.GetText())
	if proxyLink == "" {
		tui.updateStatus("Error: Please enter a proxy link", tcell.ColorRed)
		return
	}

	// Detect protocol and parse accordingly
	var config interface{}
	var protocol string
	var err error

	switch {
	case strings.HasPrefix(proxyLink, "vmess://"):
		protocol = "vmess"
		config, err = parser.VMessToSingBox(proxyLink)
	case strings.HasPrefix(proxyLink, "ss://"):
		protocol = "shadowsocks"
		config, err = parser.SSToSingBox(proxyLink)
	case strings.HasPrefix(proxyLink, "vless://"):
		protocol = "vless"
		config, err = parser.VLESSToSingBox(proxyLink)
	default:
		tui.updateStatus("Error: Invalid proxy link format. Must start with 'vmess://', 'ss://', or 'vless://'", tcell.ColorRed)
		return
	}

	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error parsing %s: %v", protocol, err), tcell.ColorRed)
		return
	}

	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error marshaling config: %v", err), tcell.ColorRed)
		return
	}

	tui.configText.SetText(string(configJSON))
	tui.updateStatus(fmt.Sprintf("%s configuration parsed successfully!", strings.Title(protocol)), tcell.ColorGreen)
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

// copyToClipboard copies text to clipboard depending on OS
func (tui *TUI) copyToClipboard(text string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xclip", "-i", "-selection", "clipboard")
	case "darwin":
		cmd = exec.Command("pbcopy")
	default:
		tui.updateStatus("Copy not supported on this system", tcell.ColorYellow)
		return
	}

	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		tui.updateStatus(fmt.Sprintf("Copy to clipboard failed: %v", err), tcell.ColorRed)
	} else {
		tui.updateStatus("Text copied to clipboard successfully!", tcell.ColorGreen)
	}
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
