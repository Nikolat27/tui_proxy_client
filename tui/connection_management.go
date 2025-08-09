package tui

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	parser "go_v2ray_client/parser"
)

// connectClient is now a high-level wrapper that delegates to helper methods
func (tui *TUI) connectClient(
	clientType string,
	parser func(string) (interface{}, error),
	command []string,
) {
	config, ok := tui.getSelectedConfig()
	if !ok {
		return
	}

	if !tui.prepareConfigFile(config, parser) {
		return
	}

	tui.updateStatus(fmt.Sprintf("Starting %s with configuration: %s...", clientType, config.Name), tcell.ColorBlue)
	tui.configText.SetText(fmt.Sprintf("Starting %s...\n", clientType))

	go tui.startClientProcess(clientType, config.Name, command)
}

// getSelectedConfig validates selection and returns the chosen config
func (tui *TUI) getSelectedConfig() (Config, bool) {
	currentIndex := tui.configList.GetCurrentItem()

	if len(tui.configs.Configurations) == 0 {
		tui.updateStatus("Error: No configurations to connect. Please add a configuration first.", tcell.ColorRed)
		return Config{}, false
	}
	if currentIndex < 0 || currentIndex >= len(tui.configs.Configurations) {
		tui.updateStatus("Error: Please select a configuration to connect.", tcell.ColorYellow)
		return Config{}, false
	}

	return tui.configs.Configurations[currentIndex], true
}

// prepareConfigFile parses, saves JSON, and updates last used time
func (tui *TUI) prepareConfigFile(config Config, parser func(string) (interface{}, error)) bool {
	parsedConfig, err := parser(config.Link)
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error parsing VMess: %v", err), tcell.ColorRed)
		return false
	}

	configJSON, err := json.MarshalIndent(parsedConfig, "", "  ")
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error marshaling config: %v", err), tcell.ColorRed)
		return false
	}

	if err := os.WriteFile("config.json", configJSON, 0644); err != nil {
		tui.updateStatus(fmt.Sprintf("Error saving config: %v", err), tcell.ColorRed)
		return false
	}

	config.LastUsed = time.Now().Format(time.RFC3339)
	tui.saveConfigsToFile()
	return true
}

// startClientProcess launches the command and manages its lifecycle
func (tui *TUI) startClientProcess(clientType, configName string, command []string) {
	cmd := exec.Command(command[0], command[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		tui.showPipeError("stdout", err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		tui.showPipeError("stderr", err)
		return
	}

	if err := cmd.Start(); err != nil {
		tui.showStartError(clientType, err)
		return
	}

	tui.app.QueueUpdateDraw(func() {
		tui.isConnected = true
		tui.clientType = clientType
		tui.connectedConfig = configName
		tui.updateStatus(fmt.Sprintf("%s started successfully with config: %s! Check your proxy settings (127.0.0.1:1080)", clientType, configName), tcell.ColorGreen)
		tui.updateConnectionStatus()
	})

	go tui.streamOutput(stdout, "")
	go tui.streamOutput(stderr, "[ERROR] ")

	if err := cmd.Wait(); err != nil {
		tui.handleClientExit(clientType, err)
	}
}

// showPipeError displays an error if pipe creation fails
func (tui *TUI) showPipeError(pipeName string, err error) {
	tui.app.QueueUpdateDraw(func() {
		tui.configText.SetText(fmt.Sprintf("Error creating %s pipe: %v", pipeName, err))
	})
}

// showStartError displays an error if starting the process fails
func (tui *TUI) showStartError(clientType string, err error) {
	tui.app.QueueUpdateDraw(func() {
		tui.configText.SetText(fmt.Sprintf("Error starting %s: %v", clientType, err))
		tui.updateStatus(fmt.Sprintf("%s error: %v", clientType, err), tcell.ColorRed)
	})
}

// handleClientExit updates UI when process ends
func (tui *TUI) handleClientExit(clientType string, err error) {
	tui.app.QueueUpdateDraw(func() {
		tui.isConnected = false
		tui.clientType = ""
		tui.connectedConfig = ""
		tui.updateStatus(fmt.Sprintf("%s stopped with error: %v", clientType, err), tcell.ColorRed)
		tui.updateConnectionStatus()
	})
}

// streamOutput reads process output and updates UI
func (tui *TUI) streamOutput(pipe io.Reader, prefix string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		tui.app.QueueUpdateDraw(func() {
			currentText := tui.configText.GetText(true)
			if len(currentText) > 10000 {
				lines := strings.Split(currentText, "\n")
				if len(lines) > 100 {
					currentText = strings.Join(lines[len(lines)-100:], "\n")
				}
			}
			tui.configText.SetText(currentText + "\n" + prefix + line)
			tui.configText.ScrollToEnd()
		})
	}
}

// connectToConfig shows the client selection modal
func (tui *TUI) connectToConfig() {
	config, ok := tui.getSelectedConfig()
	if !ok {
		return
	}

	clientModal := tview.NewModal().
		SetText(fmt.Sprintf("Choose client for configuration: %s", config.Name)).
		AddButtons([]string{"V2Ray", "SingBox", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "V2Ray":
				tui.connectClient("v2ray", func(link string) (interface{}, error) {
					return parser.VMessToV2ray(link)
				}, []string{"v2ray", "run", "config.json"})
			case "SingBox":
				tui.connectClient("singbox", func(link string) (interface{}, error) {
					return parser.VMessToSingBox(link)
				}, []string{"sing-box", "run", "-c", "config.json"})
			}
			tui.app.SetRoot(tui.mainFlex, true)
		})

	tui.app.SetRoot(clientModal, true)
}
