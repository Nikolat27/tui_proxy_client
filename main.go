package main

import (
	"log"
)

func main() {
	// Create and run the TUI
	tui := NewTUI()
	if err := tui.Run(); err != nil {
		log.Fatal(err)
	}
}
