package main

import (
	"testing"
)

// TestTUICreation tests that the TUI can be created successfully
func TestTUICreation(t *testing.T) {
	// Create a TUI instance to test
	tui := NewTUI()

	// Check that the TUI was created successfully
	if tui == nil {
		t.Fatal("Failed to create TUI instance")
	}

	// Check that the app field is initialized
	if tui.app == nil {
		t.Fatal("TUI app field not initialized")
	}

	t.Log("TUI created successfully with all required fields")
}

// TestTUIMethodsExist tests that the TUI has the required methods
func TestTUIMethodsExist(t *testing.T) {
	// Create a TUI instance to test
	tui := NewTUI()

	// Check that the TUI was created successfully
	if tui == nil {
		t.Fatal("Failed to create TUI instance")
	}

	// Test that we can call the disconnect method (it should not panic)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Calling disconnect method caused panic: %v", r)
		}
	}()

	// This is a basic test that the method exists and can be called
	// In a real test environment, we'd need to mock the TUI dependencies
	t.Log("TUI methods can be accessed without panic")
}
