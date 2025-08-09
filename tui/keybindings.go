package tui

import (
	"github.com/gdamore/tcell/v2"
)

// setupKeybindings registers all global and input-specific shortcuts
func (tui *TUI) setupKeybindings() {
	// Global keybindings
	tui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyCtrlA:
			tui.addConfig()
		case event.Key() == tcell.KeyCtrlS:
			tui.exportConfig()
		case event.Key() == tcell.KeyCtrlD:
			tui.deleteSelectedConfig()
		case event.Key() == tcell.KeyCtrlR:
			tui.renameSelectedConfig()
		case event.Key() == tcell.KeyCtrlF:
			tui.refreshConfigurations()
		case event.Key() == tcell.KeyCtrlL:
			tui.clearUI()
		case event.Key() == tcell.KeyCtrlX:
			tui.disconnect()
		case event.Key() == tcell.KeyCtrlC:
			tui.app.Stop()
		default:
			return event // let other keys pass through
		}
		return nil // handled key, no further processing
	})

	// Input field specific keybindings
	tui.vmessInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlV:
			tui.handlePaste()
			return nil
		case tcell.KeyEnter:
			tui.parseProxyLink()
			return nil
		default:
			return event
		}
	})
}
