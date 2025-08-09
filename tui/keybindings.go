package tui

import (
	"github.com/gdamore/tcell/v2"
)

// setupKeybindings sets up keyboard shortcuts
func (tui *TUI) setupKeybindings() {
	tui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlA:
			tui.addConfig()
			return nil
		case tcell.KeyCtrlS:
			tui.exportConfig()
			return nil
		case tcell.KeyCtrlD:
			tui.deleteSelectedConfig()
			return nil
		case tcell.KeyCtrlR:
			tui.renameSelectedConfig()
			return nil
		case tcell.KeyCtrlF:
			tui.refreshConfigurations()
			return nil
		case tcell.KeyCtrlL:
			tui.clearUI()
			return nil
		case tcell.KeyCtrlX:
			tui.disconnect()
			return nil
		case tcell.KeyCtrlC:
			tui.app.Stop()
			return nil
		}
		return event
	})

	tui.vmessInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlV {
			tui.handlePaste()
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			tui.parseVMess()
			return nil
		}
		return event
	})
}
