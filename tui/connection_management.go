package tui

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// connectV2Ray connects using V2Ray with the selected configuration
func (tui *TUI) connectV2Ray() {
	currentIndex := tui.configList.GetCurrentItem()

	if len(tui.configs.Configurations) == 0 {
		tui.updateStatus("Error: No configurations to connect. Please add a configuration first.", tcell.ColorRed)
		return
	}

	if currentIndex < 0 || currentIndex >= len(tui.configs.Configurations) {
		tui.updateStatus("Error: Please select a configuration to connect.", tcell.ColorYellow)
		return
	}

	config := tui.configs.Configurations[currentIndex]
	vmessLink := config.Link

	v2rayConfig, err := VMessToV2ray(vmessLink)
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error parsing VMess: %v", err), tcell.ColorRed)
		return
	}

	configJSON, err := json.MarshalIndent(v2rayConfig, "", "  ")
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error marshaling config: %v", err), tcell.ColorRed)
		return
	}

	err = os.WriteFile("config.json", configJSON, 0644)
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error saving config: %v", err), tcell.ColorRed)
		return
	}

	tui.configs.Configurations[currentIndex].LastUsed = time.Now().Format(time.RFC3339)
	tui.saveConfigsToFile()

	tui.updateStatus(fmt.Sprintf("Starting V2Ray with configuration: %s...", config.Name), tcell.ColorBlue)
	tui.configText.SetText("Starting V2Ray...\n")

	go func() {
		cmd := exec.Command("v2ray", "run", "config.json")

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			tui.app.QueueUpdateDraw(func() {
				tui.configText.SetText(fmt.Sprintf("Error creating stdout pipe: %v", err))
			})
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			tui.app.QueueUpdateDraw(func() {
				tui.configText.SetText(fmt.Sprintf("Error creating stderr pipe: %v", err))
			})
			return
		}

		if err := cmd.Start(); err != nil {
			tui.app.QueueUpdateDraw(func() {
				tui.configText.SetText(fmt.Sprintf("Error starting V2Ray: %v", err))
				tui.updateStatus(fmt.Sprintf("V2Ray error: %v", err), tcell.ColorRed)
			})
			return
		}

		tui.app.QueueUpdateDraw(func() {
			tui.isConnected = true
			tui.clientType = "v2ray"
			tui.connectedConfig = config.Name
			tui.updateStatus(fmt.Sprintf("V2Ray started successfully with config: %s! Check your proxy settings (127.0.0.1:1080)", config.Name), tcell.ColorGreen)
			tui.updateConnectionStatus()
		})

		go func() {
			scanner := bufio.NewScanner(stdout)
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
					tui.configText.SetText(currentText + "\n" + line)
					tui.configText.ScrollToEnd()
				})
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderr)
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
					tui.configText.SetText(currentText + "\n" + "[ERROR] " + line)
					tui.configText.ScrollToEnd()
				})
			}
		}()

		if err := cmd.Wait(); err != nil {
			tui.app.QueueUpdateDraw(func() {
				tui.isConnected = false
				tui.clientType = ""
				tui.connectedConfig = ""
				tui.updateStatus(fmt.Sprintf("V2Ray stopped with error: %v", err), tcell.ColorRed)
				tui.updateConnectionStatus()
			})
		}
	}()
}

// connectSingBox connects using sing-box with the selected configuration
func (tui *TUI) connectSingBox() {
	currentIndex := tui.configList.GetCurrentItem()

	if len(tui.configs.Configurations) == 0 {
		tui.updateStatus("Error: No configurations to connect. Please add a configuration first.", tcell.ColorRed)
		return
	}

	if currentIndex < 0 || currentIndex >= len(tui.configs.Configurations) {
		tui.updateStatus("Error: Please select a configuration to connect.", tcell.ColorYellow)
		return
	}

	config := tui.configs.Configurations[currentIndex]
	vmessLink := config.Link

	singboxConfig, err := VMessToSingBox(vmessLink)
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error parsing VMess: %v", err), tcell.ColorRed)
		return
	}

	configJSON, err := json.MarshalIndent(singboxConfig, "", "  ")
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error marshaling config: %v", err), tcell.ColorRed)
		return
	}

	err = os.WriteFile("config.json", configJSON, 0644)
	if err != nil {
		tui.updateStatus(fmt.Sprintf("Error saving config: %v", err), tcell.ColorRed)
		return
	}

	tui.configs.Configurations[currentIndex].LastUsed = time.Now().Format(time.RFC3339)
	tui.saveConfigsToFile()

	tui.updateStatus(fmt.Sprintf("Starting sing-box with configuration: %s...", config.Name), tcell.ColorBlue)
	tui.configText.SetText("Starting sing-box...\n")

	go func() {
		cmd := exec.Command("sing-box", "run", "-c", "config.json")

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			tui.app.QueueUpdateDraw(func() {
				tui.configText.SetText(fmt.Sprintf("Error creating stdout pipe: %v", err))
			})
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			tui.app.QueueUpdateDraw(func() {
				tui.configText.SetText(fmt.Sprintf("Error creating stderr pipe: %v", err))
			})
			return
		}

		if err := cmd.Start(); err != nil {
			tui.app.QueueUpdateDraw(func() {
				tui.configText.SetText(fmt.Sprintf("Error starting sing-box: %v", err))
				tui.updateStatus(fmt.Sprintf("sing-box error: %v", err), tcell.ColorRed)
			})
			return
		}

		tui.app.QueueUpdateDraw(func() {
			tui.isConnected = true
			tui.clientType = "singbox"
			tui.connectedConfig = config.Name
			tui.updateStatus(fmt.Sprintf("sing-box started successfully with config: %s! Check your proxy settings (127.0.0.1:1080)", config.Name), tcell.ColorGreen)
			tui.updateConnectionStatus()
		})

		go func() {
			scanner := bufio.NewScanner(stdout)
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
					tui.configText.SetText(currentText + "\n" + line)
					tui.configText.ScrollToEnd()
				})
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderr)
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
					tui.configText.SetText(currentText + "\n" + "[ERROR] " + line)
					tui.configText.ScrollToEnd()
				})
			}
		}()

		if err := cmd.Wait(); err != nil {
			tui.app.QueueUpdateDraw(func() {
				tui.isConnected = false
				tui.clientType = ""
				tui.connectedConfig = ""
				tui.updateStatus(fmt.Sprintf("sing-box stopped with error: %v", err), tcell.ColorRed)
				tui.updateConnectionStatus()
			})
		}
	}()
}

// connectToConfig shows a client selection modal and connects to the selected configuration
func (tui *TUI) connectToConfig() {
	currentIndex := tui.configList.GetCurrentItem()

	if len(tui.configs.Configurations) == 0 {
		tui.updateStatus("Error: No configurations to connect. Please add a configuration first.", tcell.ColorRed)
		return
	}

	if currentIndex < 0 || currentIndex >= len(tui.configs.Configurations) {
		tui.updateStatus("Error: Please select a configuration to connect.", tcell.ColorYellow)
		return
	}

	config := tui.configs.Configurations[currentIndex]

	clientModal := tview.NewModal().
		SetText(fmt.Sprintf("Choose client for configuration: %s", config.Name)).
		AddButtons([]string{"V2Ray", "SingBox", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "V2Ray":
				tui.connectV2Ray()
			case "SingBox":
				tui.connectSingBox()
			case "Cancel":
			}
			tui.app.SetRoot(tui.mainFlex, true)
		})

	tui.app.SetRoot(clientModal, true)
}
