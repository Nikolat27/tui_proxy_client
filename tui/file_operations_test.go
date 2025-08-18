package tui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTUI_ExportConfigFileOps(t *testing.T) {
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

func TestTUI_ShowFileDialog(t *testing.T) {
	tui := NewTUI()

	// Test showing file dialog
	tui.showFileDialog("test config content")

	// The function should not panic
	// We can't easily test the UI state changes without a full UI environment
}

func TestTUI_PerformExport(t *testing.T) {
	tui := NewTUI()

	// Test export with no config
	tui.performExport()

	// Should show error status
	statusText := tui.statusText.GetText(true)
	if statusText == "" || statusText == "Configuration will appear here..." {
		t.Error("performExport() should set error status when no config to export")
	}

	// Test export with config
	tui.configText.SetText("test config content")
	tui.performExport()

	// Should show success status
	statusText = tui.statusText.GetText(true)
	if statusText == "" || statusText == "Configuration will appear here..." {
		t.Error("performExport() should set success status when config is available")
	}
}

func TestTUI_LoadDirectoryBasic(t *testing.T) {
	tui := NewTUI()

	// Test loading current directory
	tui.loadDirectory(".")

	// The function should not panic
	// We can't easily test the file listing without a full UI environment
	// But we can verify the path is set correctly
	if tui.currentPath == "" {
		t.Error("loadDirectory() should set currentPath")
	}

	// Test loading a specific directory
	testPath := "/tmp"
	tui.loadDirectory(testPath)
	if tui.currentPath != testPath {
		t.Errorf("loadDirectory() should set currentPath to %s, got: %s", testPath, tui.currentPath)
	}
}

func TestTUI_ExportToCurrentPath(t *testing.T) {
	tui := NewTUI()

	// Test export with no config
	tui.exportToCurrentPath()

	// Should show error status
	statusText := tui.statusText.GetText(true)
	if statusText == "" || statusText == "Configuration will appear here..." {
		t.Error("exportToCurrentPath() should set error status when no config to export")
	}

	// Test export with config
	tui.configText.SetText("test config content")
	tui.currentPath = "/tmp"
	tui.exportToCurrentPath()

	// Should show success status
	statusText = tui.statusText.GetText(true)
	if statusText == "" || statusText == "Configuration will appear here..." {
		t.Error("exportToCurrentPath() should set success status when config is available")
	}
}

func TestTUI_HasConfigToExport(t *testing.T) {
	tui := NewTUI()

	// Test with no config
	if tui.hasConfigToExport() {
		t.Error("hasConfigToExport() should return false when no config")
	}

	// Test with empty config text
	tui.configText.SetText("")
	if tui.hasConfigToExport() {
		t.Error("hasConfigToExport() should return false with empty config text")
	}

	// Test with default config text
	tui.configText.SetText("Logs will appear here...\n\nTip: Press 'c' to copy all logs to clipboard")
	if tui.hasConfigToExport() {
		t.Error("hasConfigToExport() should return false with default config text")
	}

	// Test with valid config
	tui.configText.SetText("test config content")
	if !tui.hasConfigToExport() {
		t.Error("hasConfigToExport() should return true with valid config")
	}
}

func TestTUI_WriteConfigToFile(t *testing.T) {
	tui := NewTUI()

	// Set test config content
	testContent := "test config content"
	tui.configText.SetText(testContent)

	// Test writing to a temporary file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_config.json")

	err := tui.writeConfigToFile(testFile)
	if err != nil {
		t.Fatalf("writeConfigToFile() failed: %v", err)
	}

	// Verify file was created and contains correct content
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("writeConfigToFile() wrote wrong content, got: %s, want: %s", string(content), testContent)
	}

	// Test writing to invalid path
	err = tui.writeConfigToFile("/invalid/path/test.json")
	if err == nil {
		t.Error("writeConfigToFile() should return error for invalid path")
	}
}

func TestTUI_LoadDirectoryWithFiles(t *testing.T) {
	tui := NewTUI()

	// Create a temporary directory with some test files
	tempDir := t.TempDir()

	// Create some test files
	testFiles := []string{"test1.txt", "test2.json", "test3.txt"}
	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Create a test subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test subdirectory: %v", err)
	}

	// Test loading the directory
	tui.loadDirectory(tempDir)

	// Verify current path is set
	if tui.currentPath != tempDir {
		t.Errorf("loadDirectory() should set currentPath to %s, got: %s", tempDir, tui.currentPath)
	}

	// Verify path input is updated
	if tui.pathInput.GetText() != tempDir {
		t.Errorf("loadDirectory() should update pathInput to %s, got: %s", tempDir, tui.pathInput.GetText())
	}

	// The function should not panic
	// We can't easily test the file list UI updates without a full UI environment
}

func TestTUI_LoadDirectoryWithParentNavigation(t *testing.T) {
	tui := NewTUI()

	// Create a temporary directory structure
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test subdirectory: %v", err)
	}

	// Load the subdirectory
	tui.loadDirectory(subDir)
	if tui.currentPath != subDir {
		t.Errorf("loadDirectory() should set currentPath to %s, got: %s", subDir, tui.currentPath)
	}

	// The function should not panic
	// We can't easily test the parent directory navigation without a full UI environment
}

func TestTUI_ExportToCurrentPathWithTimestamp(t *testing.T) {
	tui := NewTUI()

	// Set test config content
	testContent := "test config content"
	tui.configText.SetText(testContent)

	// Set current path to temporary directory
	tempDir := t.TempDir()
	tui.currentPath = tempDir

	// Test export
	tui.exportToCurrentPath()

	// Should show success status
	statusText := tui.statusText.GetText(true)
	if statusText == "" || statusText == "Configuration will appear here..." {
		t.Error("exportToCurrentPath() should set success status when config is available")
	}

	// Check if file was created (it should have a timestamp in the name)
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read temp directory: %v", err)
	}

	found := false
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			found = true
			break
		}
	}

	if !found {
		t.Error("exportToCurrentPath() should create a JSON file with timestamp")
	}
}
