package main

import (
	"log"
	"tui_proxy_client/tui"
)

func main() {
	// Create and run the TUI
	tui := tui.NewTUI()
	if err := tui.Run(); err != nil {
		log.Fatal(err)
	}
}
