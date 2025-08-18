package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	// This is a simple test to ensure the main package compiles correctly
	// We can't easily test the actual main() function since it starts the TUI
	// But we can verify that the package imports work correctly

	// The test passes if this compiles and runs
	t.Log("Main package compiles and imports correctly")
}
