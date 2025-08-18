package parser

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestVMessToSingBox(t *testing.T) {
	// Sample VMess link (you can replace with real working one for real testing)
	vmessJSON := `{"v":"2","ps":"Test","add":"example.com","port":"443","id":"11111111-1111-1111-1111-111111111111","aid":"0","net":"ws","type":"","host":"example.com","path":"/ws","tls":"tls","sni":"example.com","fp":"chrome","scy":"auto"}`
	vmessLink := "vmess://" + base64.StdEncoding.EncodeToString([]byte(vmessJSON))

	cfg, err := VMessToSingBox(vmessLink)
	if err != nil {
		t.Fatalf("VMessToSingBox failed: %v", err)
	}

	// Save generated config
	fileName := "test.json"
	fileData, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.WriteFile(fileName, fileData, 0644); err != nil {
		t.Fatalf("failed to write test.json: %v", err)
	}

	// Run sing-box with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sing-box", "run", "-c", fileName)
	output, err := cmd.CombinedOutput()
	t.Logf("Sing-box output:\n%s", string(output))

	if ctx.Err() == context.DeadlineExceeded {
		// OK, no fatal config error detected
	} else if err != nil {
		t.Fatalf("sing-box run failed: %v", err)
	}
}

func TestVMessToV2ray(t *testing.T) {
	// Same sample VMess link
	vmessJSON := `{"v":"2","ps":"Test","add":"example.com","port":"443","id":"11111111-1111-1111-1111-111111111111","aid":"0","net":"ws","type":"","host":"example.com","path":"/ws","tls":"tls","sni":"example.com","fp":"chrome","scy":"auto"}`
	vmessLink := "vmess://" + base64.StdEncoding.EncodeToString([]byte(vmessJSON))

	cfg, err := VMessToV2ray(vmessLink)
	if err != nil {
		t.Fatalf("VMessToV2ray failed: %v", err)
	}

	// Save generated config
	fileName := "test.json"
	data, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.WriteFile(fileName, data, 0644); err != nil {
		t.Fatalf("failed to write %s: %v", fileName, err)
	}

	// Run v2ray with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "v2ray", "run", fileName)
	output, err := cmd.CombinedOutput()
	t.Logf("V2Ray output:\n%s", string(output))

	if ctx.Err() == context.DeadlineExceeded {
		// OK, config loaded without fatal error
	} else if err != nil {
		t.Fatalf("v2ray run failed: %v", err)
	}
}

// TestVMessToSingBox_InvalidInput tests error handling for invalid inputs
func TestVMessToSingBox_InvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid prefix",
			input:   "http://example.com",
			wantErr: true,
		},
		{
			name:    "invalid base64",
			input:   "vmess://invalid-base64!@#",
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			input:   "vmess://" + base64.StdEncoding.EncodeToString([]byte("invalid json")),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VMessToSingBox(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("VMessToSingBox() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestVMessToSingBox_ValidConfig tests valid VMess configuration parsing
func TestVMessToSingBox_ValidConfig(t *testing.T) {
	vmessJSON := `{
		"v": "2",
		"ps": "Test Config",
		"add": "test.example.com",
		"port": "8080",
		"id": "12345678-1234-1234-1234-123456789012",
		"aid": "0",
		"net": "tcp",
		"type": "",
		"host": "test.example.com",
		"path": "",
		"tls": "",
		"sni": "",
		"fp": "",
		"scy": "auto"
	}`
	vmessLink := "vmess://" + base64.StdEncoding.EncodeToString([]byte(vmessJSON))

	cfg, err := VMessToSingBox(vmessLink)
	if err != nil {
		t.Fatalf("VMessToSingBox() unexpected error: %v", err)
	}

	// Verify required fields exist
	if cfg == nil {
		t.Fatal("VMessToSingBox() returned nil config")
	}

	// Check for required top-level keys
	requiredKeys := []string{"log", "inbounds", "outbounds"}
	for _, key := range requiredKeys {
		if _, exists := cfg[key]; !exists {
			t.Errorf("VMessToSingBox() missing required key: %s", key)
		}
	}

	// Verify outbounds structure
	outbounds, ok := cfg["outbounds"].([]map[string]any)
	if !ok {
		t.Fatal("VMessToSingBox() outbounds is not a slice")
	}

	if len(outbounds) < 2 {
		t.Fatal("VMessToSingBox() expected at least 2 outbounds (proxy + direct)")
	}

	// Check proxy outbound
	proxy := outbounds[0]
	if proxy["type"] != "vmess" {
		t.Errorf("VMessToSingBox() proxy type = %v, want vmess", proxy["type"])
	}

	if proxy["server"] != "test.example.com" {
		t.Errorf("VMessToSingBox() server = %v, want test.example.com", proxy["server"])
	}

	if proxy["server_port"] != 8080 {
		t.Errorf("VMessToSingBox() server_port = %v, want 8080", proxy["server_port"])
	}
}

// TestVMessToV2ray_InvalidInput tests error handling for invalid inputs
func TestVMessToV2ray_InvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid prefix",
			input:   "http://example.com",
			wantErr: true,
		},
		{
			name:    "invalid base64",
			input:   "vmess://invalid-base64!@#",
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			input:   "vmess://" + base64.StdEncoding.EncodeToString([]byte("invalid json")),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VMessToV2ray(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("VMessToV2ray() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestVMessToV2ray_ValidConfig tests valid VMess configuration parsing
func TestVMessToV2ray_ValidConfig(t *testing.T) {
	vmessJSON := `{
		"v": "2",
		"ps": "Test Config",
		"add": "test.example.com",
		"port": "8080",
		"id": "12345678-1234-1234-1234-123456789012",
		"aid": "0",
		"net": "tcp",
		"type": "",
		"host": "test.example.com",
		"path": "",
		"tls": "",
		"sni": "",
		"fp": "",
		"scy": "auto"
	}`
	vmessLink := "vmess://" + base64.StdEncoding.EncodeToString([]byte(vmessJSON))

	cfg, err := VMessToV2ray(vmessLink)
	if err != nil {
		t.Fatalf("VMessToV2ray() unexpected error: %v", err)
	}

	// Verify required fields exist
	if cfg == nil {
		t.Fatal("VMessToV2ray() returned nil config")
	}

	// Check for required top-level keys
	requiredKeys := []string{"log", "dns", "inbounds", "outbounds"}
	for _, key := range requiredKeys {
		if _, exists := cfg[key]; !exists {
			t.Errorf("VMessToV2ray() missing required key: %s", key)
		}
	}

	// Verify outbounds structure
	outbounds, ok := cfg["outbounds"].([]map[string]any)
	if !ok {
		t.Fatal("VMessToV2ray() outbounds is not a slice")
	}

	if len(outbounds) < 1 {
		t.Fatal("VMessToV2ray() expected at least 1 outbound")
	}

	// Check proxy outbound
	proxy := outbounds[0]
	if proxy["protocol"] != "vmess" {
		t.Errorf("VMessToV2ray() protocol = %v, want vmess", proxy["protocol"])
	}

	settings, ok := proxy["settings"].(map[string]any)
	if !ok {
		t.Fatal("VMessToV2ray() settings is not a map")
	}

	vnext, ok := settings["vnext"].([]map[string]any)
	if !ok || len(vnext) == 0 {
		t.Fatal("VMessToV2ray() vnext is not a valid slice")
	}

	server := vnext[0]["address"]
	if server != "test.example.com" {
		t.Errorf("VMessToV2ray() server = %v, want test.example.com", server)
	}
}

// TestAtoiSafe tests the atoiSafe helper function
func TestAtoiSafe(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"valid number", "123", 123},
		{"zero", "0", 0},
		{"large number", "65535", 65535},
		{"empty string", "", 0},
		{"invalid string", "abc", 0},
		{"mixed string", "123abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := atoiSafe(tt.input)
			if result != tt.expected {
				t.Errorf("atoiSafe(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}
