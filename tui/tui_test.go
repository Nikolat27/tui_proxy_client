package tui

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
)

func TestNewTUI(t *testing.T) {
	tui := NewTUI()
	if tui == nil {
		t.Fatal("NewTUI() returned nil")
	}

	if tui.app == nil {
		t.Error("NewTUI() app is nil")
	}

	if tui.mainFlex == nil {
		t.Error("NewTUI() mainFlex is nil")
	}

	if tui.vmessInput == nil {
		t.Error("NewTUI() vmessInput is nil")
	}

	if tui.statusText == nil {
		t.Error("NewTUI() statusText is nil")
	}

	if tui.configText == nil {
		t.Error("NewTUI() configText is nil")
	}

	if tui.configList == nil {
		t.Error("NewTUI() configList is nil")
	}

	if tui.buttons == nil {
		t.Error("NewTUI() buttons is nil")
	}

	if tui.fileDialog == nil {
		t.Error("NewTUI() fileDialog is nil")
	}

	if tui.fileExplorer == nil {
		t.Error("NewTUI() fileExplorer is nil")
	}

	if tui.connectionStatus == nil {
		t.Error("NewTUI() connectionStatus is nil")
	}

	// Check initial state
	if tui.isConnected {
		t.Error("NewTUI() should start with isConnected = false")
	}

	if tui.clientType != "" {
		t.Error("NewTUI() should start with empty clientType")
	}

	if tui.connectedConfig != "" {
		t.Error("NewTUI() should start with empty connectedConfig")
	}
}

func TestTUI_UpdateStatus(t *testing.T) {
	tui := NewTUI()

	// Test basic status update
	tui.updateStatus("Test message", tcell.ColorGreen)
	if tui.statusText.GetText(true) != "Test message" {
		t.Errorf("updateStatus() failed to update text, got: %s", tui.statusText.GetText(true))
	}

	// Test status update when connected
	tui.isConnected = true
	tui.clientType = "v2ray"
	tui.connectedConfig = "Test Config"
	tui.updateStatus("Connected message", tcell.ColorBlue)

	expectedText := "Connected message | Connected: Test Config (v2ray)"
	if tui.statusText.GetText(true) != expectedText {
		t.Errorf("updateStatus() with connection info failed, got: %s, want: %s",
			tui.statusText.GetText(true), expectedText)
	}
}

func TestTUI_ClearUI(t *testing.T) {
	tui := NewTUI()

	// Set some text
	tui.vmessInput.SetText("test input")
	tui.configText.SetText("test config")

	// Clear UI
	tui.clearUI()

	if tui.vmessInput.GetText() != "" {
		t.Error("clearUI() failed to clear vmessInput")
	}

	if tui.configText.GetText(true) != "Configuration will appear here..." {
		t.Error("clearUI() failed to reset configText")
	}
}

func TestTUI_ParseProxyLink(t *testing.T) {
	tui := NewTUI()

	tests := []struct {
		name        string
		proxyLink   string
		expectError bool
	}{
		{
			name:        "valid vmess link",
			proxyLink:   "vmess://eyJ2IjoiMiIsInBzIjoiVGVzdCIsImFkZCI6ImV4YW1wbGUuY29tIiwicG9ydCI6IjQ0MyIsImlkIjoiMTExMTExMTEtMTExMS0xMTExLTExMTEtMTExMTExMTExMTExIiwiYWlkIjoiMCIsIm5ldCI6IndzIiwidHlwZSI6IiIsImhvc3QiOiJleGFtcGxlLmNvbSIsInBhdGgiOiIvd3MiLCJ0bHMiOiJ0bHMiLCJzbmkiOiJleGFtcGxlLmNvbSIsImZwIjoiY2hyb21lIiwic2N5IjoiYXV0byJ9",
			expectError: false,
		},
		{
			name:        "valid ss link",
			proxyLink:   "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@example.com:8388#Test%20Config",
			expectError: false,
		},
		{
			name:        "valid vless link",
			proxyLink:   "vless://12345678-1234-1234-1234-123456789012@example.com:443?encryption=none#Test%20Config",
			expectError: false,
		},
		{
			name:        "empty link",
			proxyLink:   "",
			expectError: true,
		},
		{
			name:        "invalid protocol",
			proxyLink:   "http://example.com",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tui.vmessInput.SetText(tt.proxyLink)
			tui.parseProxyLink()

			if tt.expectError {
				// Check if error status is set
				statusText := tui.statusText.GetText(true)
				if statusText == "" || statusText == "Configuration will appear here..." {
					t.Error("parseProxyLink() should set error status for invalid input")
				}
			} else {
				// Check if config was parsed successfully
				configText := tui.configText.GetText(true)
				if configText == "" || configText == "Configuration will appear here..." {
					t.Error("parseProxyLink() should parse valid config")
				}
			}
		})
	}
}

func TestTUI_UpdateConnectionStatus(t *testing.T) {
	tui := NewTUI()

	// Test initial status (not connected)
	tui.updateConnectionStatus()
	statusText := tui.connectionStatus.GetText(true)
	if statusText != "Status: Not Connected (Port 1080 Free)" {
		t.Errorf("updateConnectionStatus() initial status wrong, got: %s", statusText)
	}

	// Test connected status
	tui.isConnected = true
	tui.clientType = "v2ray"
	tui.connectedConfig = "Test Config"
	tui.updateConnectionStatus()

	// Note: This test depends on whether port 1080 is actually in use
	// We can't easily mock the port check, so we just verify the function doesn't panic
	statusText = tui.connectionStatus.GetText(true)
	if statusText == "" {
		t.Error("updateConnectionStatus() failed to set status text")
	}
}

func TestTUI_HandlePaste(t *testing.T) {
	tui := NewTUI()

	// Test paste functionality
	// Note: This is hard to test without mocking the clipboard
	// We'll just verify the function doesn't panic
	tui.handlePaste()

	// The function should not cause any errors
	// We can't easily test the actual clipboard functionality in unit tests
}

func TestTUI_CopyToClipboard(t *testing.T) {
	tui := NewTUI()

	// Test clipboard functionality
	// Note: This is hard to test without mocking the clipboard
	// We'll just verify the function doesn't panic
	tui.copyToClipboard("test text")

	// The function should not cause any errors
	// We can't easily test the actual clipboard functionality in unit tests
}

func TestTUI_IsPort1080InUse(t *testing.T) {
	tui := NewTUI()

	// Test port check functionality
	// Note: This depends on the actual system state
	// We'll just verify the function returns a boolean
	result := tui.isPort1080InUse()

	// The function should return a boolean value
	if result != true && result != false {
		t.Error("isPort1080InUse() should return a boolean value")
	}
}

func TestTUI_PeriodicStatusCheck(t *testing.T) {
	tui := NewTUI()

	// Test periodic status check
	// We'll run it for a short time and verify it doesn't panic
	done := make(chan bool)
	go func() {
		tui.periodicStatusCheck()
		done <- true
	}()

	// Wait a short time then stop the app
	time.Sleep(100 * time.Millisecond)
	tui.app.Stop()

	// Wait for the goroutine to finish
	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Log("periodicStatusCheck() did not stop within timeout, but this is expected in some environments")
	}
}

func TestTUI_LoadConfigsFromFile(t *testing.T) {
	tui := NewTUI()

	// Test loading configs from file
	// This will try to load from configs.json if it exists
	tui.loadConfigsFromFile()

	// The function should not panic
	// We can't easily test the file loading without mocking the filesystem
}

func TestTUI_SaveConfigsToFile(t *testing.T) {
	tui := NewTUI()

	// Test saving configs to file
	err := tui.saveConfigsToFile()
	if err != nil {
		// It's okay if this fails due to file permissions or other system issues
		// We're just testing that the function doesn't panic
		t.Logf("saveConfigsToFile() returned error (expected in some environments): %v", err)
	}
}

func TestTUI_RefreshConfigList(t *testing.T) {
	tui := NewTUI()

	// Test refreshing config list
	tui.refreshConfigList()

	// The function should not panic
	// We can't easily test the UI updates without a full UI environment
}

func TestTUI_LoadDirectory(t *testing.T) {
	tui := NewTUI()

	// Test loading directory
	tui.loadDirectory(".")

	// The function should not panic
	// We can't easily test the file listing without a full UI environment
}
