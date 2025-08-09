package tui

import (
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// setupUI initializes the user interface components
func (tui *TUI) setupUI() {
	title := tui.createTitle()
	connectionStatus := tui.createConnectionStatus()
	tui.vmessInput = tui.createVMessInput()
	tui.statusText = tui.createStatusText()
	tui.configText = tui.createConfigText()
	tui.configList = tui.createConfigList()
	tui.buttons = tui.createButtons()
	tui.fileDialog = tui.createFileDialog()
	tui.fileExplorer = tui.createFileExplorer()

	configSection := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tui.configText, 0, 2, false).
		AddItem(tui.configList, 0, 1, false)

	tui.mainFlex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 3, 0, false).
		AddItem(connectionStatus, 3, 0, false).
		AddItem(tui.vmessInput, 3, 0, false).
		AddItem(tui.statusText, 3, 0, false).
		AddItem(configSection, 0, 1, false).
		AddItem(tui.buttons, 3, 0, false)

	tui.connectionStatus = connectionStatus
	tui.app.SetRoot(tui.mainFlex, true)
}

func (tui *TUI) createTitle() *tview.TextView {
	return tview.NewTextView().
		SetText("V2Ray Client Configuration Generator").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorYellow).
		SetDynamicColors(true)
}

func (tui *TUI) createConnectionStatus() *tview.TextView {
	status := tview.NewTextView()
	status.SetText("Status: Not Connected")
	status.SetTextAlign(tview.AlignCenter)
	status.SetTextColor(tcell.ColorRed)
	status.SetBorder(true)
	status.SetTitle(" Connection Status ")
	return status
}

func (tui *TUI) createVMessInput() *tview.InputField {
	input := tview.NewInputField()
	input.SetLabel("Proxy Link: ")
	input.SetPlaceholder("Supported: vmess, vless and ss")
	input.SetFieldWidth(80)
	input.SetBorder(true)
	input.SetTitle(" Enter Proxy Configuration (VMess/SS/VLESS) ")
	return input
}

func (tui *TUI) createStatusText() *tview.TextView {
	text := tview.NewTextView()
	text.SetText("Ready to parse proxy configuration (VMess/SS/VLESS)")
	text.SetTextAlign(tview.AlignCenter)
	text.SetTextColor(tcell.ColorGreen)
	text.SetBorder(true)
	text.SetTitle(" Status ")
	return text
}

func (tui *TUI) createConfigText() *tview.TextView {
	text := tview.NewTextView()
	text.SetText("Logs will appear here...\n\nTip: Press 'c' to copy all logs to clipboard")
	text.SetTextAlign(tview.AlignLeft)
	text.SetTextColor(tcell.ColorWhite)
	text.SetBorder(true)
	text.SetTitle(" Logs (Press 'c' to copy) ")
	text.SetScrollable(true)
	text.SetChangedFunc(func() {
		text.ScrollToEnd()
	})

	// Simple copy functionality
	text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'c' || event.Rune() == 'C' {
			// Copy all text to clipboard
			fullText := text.GetText(true)
			if fullText != "" {
				tui.copyToClipboard(fullText)
				text.SetTitle(" Logs (Copied! Press 'c' again to copy) ")
				// Reset title after 2 seconds
				go func() {
					time.Sleep(2 * time.Second)
					tui.app.QueueUpdateDraw(func() {
						text.SetTitle(" Logs (Press 'c' to copy) ")
					})
				}()
			}
			return nil
		}
		return event
	})

	return text
}

func (tui *TUI) createConfigList() *tview.List {
	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle(" Saved Configurations ")
	list.SetMainTextColor(tcell.ColorWhite)
	return list
}

func (tui *TUI) createButtons() *tview.Flex {
	addConfigBtn := tview.NewButton("Add Config\n(Ctrl+A)").
		SetSelectedFunc(func() {
			tui.addConfig()
		})

	exportBtn := tview.NewButton("Export Config\n(Ctrl+S)").
		SetSelectedFunc(func() {
			tui.exportConfig()
		})

	connectBtn := tview.NewButton("Connect").
		SetSelectedFunc(func() {
			tui.connectToConfig()
		})

	disconnectBtn := tview.NewButton("Disconnect\n(Ctrl+X)").
		SetSelectedFunc(func() {
			tui.disconnect()
		})

	deleteBtn := tview.NewButton("Delete Config\n(Ctrl+D)").
		SetSelectedFunc(func() {
			tui.deleteSelectedConfig()
		})

	renameBtn := tview.NewButton("Rename Config\n(Ctrl+R)").
		SetSelectedFunc(func() {
			tui.renameSelectedConfig()
		})

	refreshBtn := tview.NewButton("Refresh\n(Ctrl+F)").
		SetSelectedFunc(func() {
			tui.refreshConfigurations()
		})

	clearBtn := tview.NewButton("Clear\n(Ctrl+L)").
		SetSelectedFunc(func() {
			tui.clearUI()
		})

	quitBtn := tview.NewButton("Quit\n(Ctrl+C)").
		SetSelectedFunc(func() {
			tui.app.Stop()
		})

	return tview.NewFlex().
		AddItem(addConfigBtn, 0, 1, false).
		AddItem(exportBtn, 0, 1, false).
		AddItem(connectBtn, 0, 1, false).
		AddItem(disconnectBtn, 0, 1, false).
		AddItem(deleteBtn, 0, 1, false).
		AddItem(tview.NewBox(), 1, 0, false). // <--- Spacer, 1 column wide

		AddItem(renameBtn, 0, 1, false).
		AddItem(refreshBtn, 0, 1, false).
		AddItem(clearBtn, 0, 1, false).
		AddItem(quitBtn, 0, 1, false)
}

func (tui *TUI) createFileDialog() *tview.Modal {
	return tview.NewModal().
		SetText("Enter filename to export:").
		AddButtons([]string{"Export", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Export" {
				tui.performExport()
			}
			tui.app.SetRoot(tui.mainFlex, true)
		})
}

func (tui *TUI) createFileExplorer() *tview.Flex {
	tui.currentPath, _ = os.Getwd()
	tui.pathInput = tview.NewInputField()
	tui.pathInput.SetLabel("Path: ")
	tui.pathInput.SetText(tui.currentPath)
	tui.pathInput.SetFieldWidth(50)
	tui.pathInput.SetBorder(true)
	tui.pathInput.SetTitle(" Current Directory ")

	tui.fileList = tview.NewList()
	tui.fileList.SetBorder(true)
	tui.fileList.SetTitle(" Files and Folders ")

	exportHereBtn := tview.NewButton("Export Here").SetSelectedFunc(func() {
		tui.exportToCurrentPath()
	})

	backToMainBtn := tview.NewButton("Back to Main").SetSelectedFunc(func() {
		tui.app.SetRoot(tui.mainFlex, true)
	})

	explorerButtons := tview.NewFlex().
		AddItem(exportHereBtn, 0, 1, false).
		AddItem(backToMainBtn, 0, 1, false)

	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tui.pathInput, 3, 0, false).
		AddItem(tui.fileList, 0, 1, false).
		AddItem(explorerButtons, 3, 0, false)
}
