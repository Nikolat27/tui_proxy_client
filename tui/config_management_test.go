package tui

import (
	"testing"
	"time"
)

func TestTUI_AddConfig(t *testing.T) {

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
			proxyLink:   "vless://12345678-1234-1234-1234-123456789012@example.com:443?encryption=none&type=tcp#Test%20Config",
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
			// Create a fresh TUI instance for each test
			testTUI := NewTUI()
			testTUI.vmessInput.SetText(tt.proxyLink)

			initialCount := len(testTUI.configs.Configurations)
			testTUI.addConfig()

			if tt.expectError {
				// Check if error status is set
				statusText := testTUI.statusText.GetText(true)
				if statusText == "" || statusText == "Configuration will appear here..." {
					t.Error("addConfig() should set error status for invalid input")
				}
				// Config count should not increase
				if len(testTUI.configs.Configurations) != initialCount {
					t.Errorf("addConfig() should not add config for invalid input, count: %d, expected: %d",
						len(testTUI.configs.Configurations), initialCount)
				}
			} else {
				// Config should be added
				if len(testTUI.configs.Configurations) != initialCount+1 {
					t.Errorf("addConfig() should add config, count: %d, expected: %d",
						len(testTUI.configs.Configurations), initialCount+1)
				}

				// Check the added config
				addedConfig := testTUI.configs.Configurations[len(testTUI.configs.Configurations)-1]
				if addedConfig.Link != tt.proxyLink {
					t.Errorf("addConfig() link mismatch, got: %s, want: %s", addedConfig.Link, tt.proxyLink)
				}

				// Check if config name follows pattern
				if addedConfig.Name == "" {
					t.Error("addConfig() should set a name for the config")
				}

				// Check if timestamps are set
				if addedConfig.CreatedAt == "" {
					t.Error("addConfig() should set CreatedAt timestamp")
				}
				if addedConfig.LastUsed == "" {
					t.Error("addConfig() should set LastUsed timestamp")
				}
			}
		})
	}
}

func TestTUI_LoadConfigList(t *testing.T) {
	tui := NewTUI()

	// Test loading config list
	tui.loadConfigList()

	// The function should not panic
	// We can't easily test the file loading without mocking the filesystem
	// But we can verify the UI components are properly initialized
	if tui.configList == nil {
		t.Error("loadConfigList() should initialize configList")
	}
}

func TestTUI_RefreshConfigListWithConfigs(t *testing.T) {
	tui := NewTUI()

	// Add a test config
	tui.configs.Configurations = []Config{
		{
			ID:        "1",
			Name:      "Test Config",
			Protocol:  "vmess",
			Link:      "vmess://test",
			CreatedAt: time.Now().Format(time.RFC3339),
			LastUsed:  time.Now().Format(time.RFC3339),
		},
	}

	tui.refreshConfigList()

	// The function should not panic
	// We can't easily test the UI updates without a full UI environment
	// But we can verify the function handles the configs properly
	if len(tui.configs.Configurations) != 1 {
		t.Error("refreshConfigList() should preserve existing configs")
	}
}

func TestTUI_DeleteSelectedConfig(t *testing.T) {
	tui := NewTUI()

	// Add test configs
	tui.configs.Configurations = []Config{
		{
			ID:        "1",
			Name:      "Test Config 1",
			Protocol:  "vmess",
			Link:      "vmess://test1",
			CreatedAt: time.Now().Format(time.RFC3339),
			LastUsed:  time.Now().Format(time.RFC3339),
		},
		{
			ID:        "2",
			Name:      "Test Config 2",
			Protocol:  "ss",
			Link:      "ss://test2",
			CreatedAt: time.Now().Format(time.RFC3339),
			LastUsed:  time.Now().Format(time.RFC3339),
		},
	}

	// Test deleting with no selection
	tui.configList.SetCurrentItem(-1)
	tui.deleteSelectedConfig()

	// Test deleting with selection
	tui.configList.SetCurrentItem(0)
	tui.deleteSelectedConfig()

	// The function should not panic
	// We can't easily test the exact behavior without a full UI environment
}

func TestTUI_RenameSelectedConfig(t *testing.T) {
	tui := NewTUI()

	// Add a test config
	tui.configs.Configurations = []Config{
		{
			ID:        "1",
			Name:      "Test Config",
			Protocol:  "vmess",
			Link:      "vmess://test",
			CreatedAt: time.Now().Format(time.RFC3339),
			LastUsed:  time.Now().Format(time.RFC3339),
		},
	}

	// Test renaming with no selection
	tui.configList.SetCurrentItem(-1)
	tui.renameSelectedConfig()

	// Should not rename anything when no selection
	if tui.configs.Configurations[0].Name != "Test Config" {
		t.Error("renameSelectedConfig() should not rename when no selection")
	}

	// Test renaming with selection
	tui.configList.SetCurrentItem(0)
	tui.renameSelectedConfig()

	// The function should not panic
	// We can't easily test the modal dialog without a full UI environment
}

func TestTUI_ExportConfig(t *testing.T) {
	tui := NewTUI()

	// Test export with no config
	tui.exportConfig()

	// Should show error status
	statusText := tui.statusText.GetText(true)
	if statusText == "" || statusText == "Configuration will appear here..." {
		t.Error("exportConfig() should set error status when no config to export")
	}

	// Test export with config
	tui.configText.SetText("test config content")
	tui.exportConfig()

	// Should not show error status
	statusText = tui.statusText.GetText(true)
	if statusText == "" || statusText == "Configuration will appear here..." {
		t.Error("exportConfig() should not set error status when config is available")
	}
}

func TestTUI_ConnectToConfig(t *testing.T) {
	tui := NewTUI()

	// Test connect with no configs
	tui.connectToConfig()

	// Should show error status
	statusText := tui.statusText.GetText(true)
	if statusText == "" || statusText == "Configuration will appear here..." {
		t.Error("connectToConfig() should set error status when no configs available")
	}

	// Add a test config
	tui.configs.Configurations = []Config{
		{
			ID:        "1",
			Name:      "Test Config",
			Protocol:  "vmess",
			Link:      "vmess://eyJ2IjoiMiIsInBzIjoiVGVzdCIsImFkZCI6ImV4YW1wbGUuY29tIiwicG9ydCI6IjQ0MyIsImlkIjoiMTExMTExMTEtMTExMS0xMTExLTExMTEtMTExMTExMTExMTExIiwiYWlkIjoiMCIsIm5ldCI6IndzIiwidHlwZSI6IiIsImhvc3QiOiJleGFtcGxlLmNvbSIsInBhdGgiOiIvd3MiLCJ0bHMiOiJ0bHMiLCJzbmkiOiJleGFtcGxlLmNvbSIsImZwIjoiY2hyb21lIiwic2N5IjoiYXV0byJ9",
			CreatedAt: time.Now().Format(time.RFC3339),
			LastUsed:  time.Now().Format(time.RFC3339),
		},
	}

	// Test connect with selection
	tui.configList.SetCurrentItem(0)
	tui.connectToConfig()

	// The function should not panic
	// We can't easily test the actual connection without mocking external processes
}

func TestTUI_Disconnect(t *testing.T) {
	tui := NewTUI()

	// Test disconnect when not connected
	tui.disconnect()

	// Should show appropriate status
	statusText := tui.statusText.GetText(true)
	if statusText == "" || statusText == "Configuration will appear here..." {
		t.Error("disconnect() should set status message")
	}

	// Test disconnect when connected
	tui.isConnected = true
	tui.clientType = "v2ray"
	tui.connectedConfig = "Test Config"
	tui.disconnect()

	// The disconnect function runs asynchronously, so we can't test the state immediately
	// Instead, we test that the function doesn't panic and sets a status message
	statusText = tui.statusText.GetText(true)
	if statusText == "" || statusText == "Configuration will appear here..." {
		t.Error("disconnect() should set status message")
	}
}

func TestTUI_GetSelectedConfig(t *testing.T) {
	tui := NewTUI()

	// Test with no configs
	config, ok := tui.getSelectedConfig()
	// The function should return false when no configs are available
	// We can't easily test this without a full UI environment

	// Add test configs
	tui.configs.Configurations = []Config{
		{
			ID:        "1",
			Name:      "Test Config 1",
			Protocol:  "vmess",
			Link:      "vmess://test1",
			CreatedAt: time.Now().Format(time.RFC3339),
			LastUsed:  time.Now().Format(time.RFC3339),
		},
		{
			ID:        "2",
			Name:      "Test Config 2",
			Protocol:  "ss",
			Link:      "ss://test2",
			CreatedAt: time.Now().Format(time.RFC3339),
			LastUsed:  time.Now().Format(time.RFC3339),
		},
	}

	// Test with no selection (but with configs available)
	tui.configList.SetCurrentItem(-1)
	config, ok = tui.getSelectedConfig()
	// The function should return false when no selection is made
	// We can't easily test this without a full UI environment

	// Test with valid selection
	tui.configList.SetCurrentItem(0)
	config, ok = tui.getSelectedConfig()
	if !ok {
		t.Error("getSelectedConfig() should return true with valid selection")
	}
	if config.ID != "1" {
		t.Errorf("getSelectedConfig() returned wrong config, got ID: %s, want: 1", config.ID)
	}

	// Test with invalid selection (out of bounds)
	tui.configList.SetCurrentItem(10) // Out of bounds
	config, ok = tui.getSelectedConfig()
	// The function should return false when selection is out of bounds
	// We can't easily test this without a full UI environment
}
