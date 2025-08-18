package tui

import (
	"github.com/rivo/tview"
)

// NewTUI creates a new TUI instance
func NewTUI() *TUI {
	tui := &TUI{
		app: tview.NewApplication(),
	}

	tui.app.EnableMouse(true)
	tui.setupUI()
	tui.setupKeybindings()
	tui.loadConfigList()
	tui.loadDirectory(tui.currentPath)

	go tui.periodicStatusCheck()

	return tui
}

// Run starts the TUI application
func (tui *TUI) Run() error {
	return tui.app.Run()
}
