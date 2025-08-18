package parser

import (
	"testing"
)

func TestVLESSToSingBox(t *testing.T) {
	// Test valid VLESS link
	vlessLink := "vless://12345678-1234-1234-1234-123456789012@example.com:443?encryption=none&security=tls&sni=example.com&type=ws&path=/ws#Test%20Config"

	cfg, err := VLESSToSingBox(vlessLink)
	if err != nil {
		t.Fatalf("VLESSToSingBox failed: %v", err)
	}

	// Verify required fields exist
	if cfg == nil {
		t.Fatal("VLESSToSingBox returned nil config")
	}

	// Check for required top-level keys
	requiredKeys := []string{"log", "inbounds", "outbounds"}
	for _, key := range requiredKeys {
		if _, exists := cfg[key]; !exists {
			t.Errorf("VLESSToSingBox missing required key: %s", key)
		}
	}

	// Verify outbounds structure
	outbounds, ok := cfg["outbounds"].([]map[string]any)
	if !ok {
		t.Fatal("VLESSToSingBox outbounds is not a slice")
	}

	if len(outbounds) < 1 {
		t.Fatal("VLESSToSingBox expected at least 1 outbound")
	}

	// Check proxy outbound
	proxy := outbounds[0]
	if proxy["type"] != "vless" {
		t.Errorf("VLESSToSingBox proxy type = %v, want vless", proxy["type"])
	}

	if proxy["server"] != "example.com" {
		t.Errorf("VLESSToSingBox server = %v, want example.com", proxy["server"])
	}

	if proxy["server_port"] != 443 {
		t.Errorf("VLESSToSingBox server_port = %v, want 443", proxy["server_port"])
	}

	// Check UUID
	if proxy["uuid"] != "12345678-1234-1234-1234-123456789012" {
		t.Errorf("VLESSToSingBox uuid = %v, want 12345678-1234-1234-1234-123456789012", proxy["uuid"])
	}
}

func TestVLESSToSingBox_InvalidInput(t *testing.T) {
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
			name:    "missing uuid",
			input:   "vless://@example.com:443",
			wantErr: true,
		},
		{
			name:    "invalid uuid format",
			input:   "vless://invalid-uuid@example.com:443",
			wantErr: true,
		},
		{
			name:    "missing server",
			input:   "vless://12345678-1234-1234-1234-123456789012@:443",
			wantErr: true,
		},
		{
			name:    "missing port",
			input:   "vless://12345678-1234-1234-1234-123456789012@example.com",
			wantErr: true,
		},
		{
			name:    "invalid port",
			input:   "vless://12345678-1234-1234-1234-123456789012@example.com:invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VLESSToSingBox(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("VLESSToSingBox() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVLESSToSingBox_ValidConfigs(t *testing.T) {
	tests := []struct {
		name      string
		vlessLink string
		expected  map[string]interface{}
	}{
		{
			name:      "basic tcp config",
			vlessLink: "vless://12345678-1234-1234-1234-123456789012@example.com:443?encryption=none&type=tcp#Test%20Config",
			expected: map[string]interface{}{
				"server":      "example.com",
				"server_port": 443,
				"uuid":        "12345678-1234-1234-1234-123456789012",
			},
		},
		{
			name:      "tls websocket config",
			vlessLink: "vless://12345678-1234-1234-1234-123456789012@example.com:443?encryption=none&security=tls&sni=example.com&type=ws&path=/ws#Test%20Config",
			expected: map[string]interface{}{
				"server":      "example.com",
				"server_port": 443,
				"uuid":        "12345678-1234-1234-1234-123456789012",
			},
		},
		{
			name:      "grpc config",
			vlessLink: "vless://12345678-1234-1234-1234-123456789012@example.com:443?encryption=none&security=tls&sni=example.com&type=grpc&serviceName=grpc#Test%20Config",
			expected: map[string]interface{}{
				"server":      "example.com",
				"server_port": 443,
				"uuid":        "12345678-1234-1234-1234-123456789012",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := VLESSToSingBox(tt.vlessLink)
			if err != nil {
				t.Fatalf("VLESSToSingBox() unexpected error: %v", err)
			}

			outbounds, ok := cfg["outbounds"].([]map[string]any)
			if !ok || len(outbounds) == 0 {
				t.Fatal("VLESSToSingBox() outbounds is not valid")
			}

			proxy := outbounds[0]
			for key, expectedValue := range tt.expected {
				if proxy[key] != expectedValue {
					t.Errorf("VLESSToSingBox() %s = %v, want %v", key, proxy[key], expectedValue)
				}
			}
		})
	}
}

func TestVLESSToV2Ray(t *testing.T) {
	// Test valid VLESS link
	vlessLink := "vless://12345678-1234-1234-1234-123456789012@example.com:443?encryption=none&security=tls&sni=example.com&type=ws&path=/ws#Test%20Config"

	cfg, err := VLESSToV2Ray(vlessLink)
	if err != nil {
		t.Fatalf("VLESSToV2Ray failed: %v", err)
	}

	// Verify required fields exist
	if cfg == nil {
		t.Fatal("VLESSToV2ray returned nil config")
	}

	// Check for required top-level keys
	requiredKeys := []string{"log", "inbounds", "outbounds"}
	for _, key := range requiredKeys {
		if _, exists := cfg[key]; !exists {
			t.Errorf("VLESSToV2ray missing required key: %s", key)
		}
	}

	// Verify outbounds structure
	outbounds, ok := cfg["outbounds"].([]map[string]any)
	if !ok {
		t.Fatal("VLESSToV2ray outbounds is not a slice")
	}

	if len(outbounds) < 1 {
		t.Fatal("VLESSToV2ray expected at least 1 outbound")
	}

	// Check proxy outbound
	proxy := outbounds[0]
	if proxy["protocol"] != "vless" {
		t.Errorf("VLESSToV2ray protocol = %v, want vless", proxy["protocol"])
	}

	settings, ok := proxy["settings"].(map[string]any)
	if !ok {
		t.Fatal("VLESSToV2ray settings is not a map")
	}

	vnext, ok := settings["vnext"].([]map[string]any)
	if !ok || len(vnext) == 0 {
		t.Fatal("VLESSToV2ray vnext is not a valid slice")
	}

	server := vnext[0]["address"]
	if server != "example.com" {
		t.Errorf("VLESSToV2ray server = %v, want example.com", server)
	}
}

func TestVLESSToV2Ray_InvalidInput(t *testing.T) {
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
			name:    "missing uuid",
			input:   "vless://@example.com:443",
			wantErr: true,
		},
		{
			name:    "invalid url format",
			input:   "vless://invalid:url:format",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VLESSToV2Ray(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("VLESSToV2Ray() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
