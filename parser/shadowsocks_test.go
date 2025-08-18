package parser

import (
	"testing"
)

func TestSSToSingBox(t *testing.T) {
	// Test valid Shadowsocks link
	// Format: ss://method:password@server:port#name
	ssLink := "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@example.com:8388#Test%20Config"

	cfg, err := SSToSingBox(ssLink)
	if err != nil {
		t.Fatalf("SSToSingBox failed: %v", err)
	}

	// Verify required fields exist
	if cfg == nil {
		t.Fatal("SSToSingBox returned nil config")
	}

	// Check for required top-level keys
	requiredKeys := []string{"log", "inbounds", "outbounds"}
	for _, key := range requiredKeys {
		if _, exists := cfg[key]; !exists {
			t.Errorf("SSToSingBox missing required key: %s", key)
		}
	}

	// Verify outbounds structure
	outbounds, ok := cfg["outbounds"].([]map[string]any)
	if !ok {
		t.Fatal("SSToSingBox outbounds is not a slice")
	}

	if len(outbounds) < 2 {
		t.Fatal("SSToSingBox expected at least 2 outbounds (proxy + direct)")
	}

	// Check proxy outbound
	proxy := outbounds[0]
	if proxy["type"] != "shadowsocks" {
		t.Errorf("SSToSingBox proxy type = %v, want shadowsocks", proxy["type"])
	}

	if proxy["server"] != "example.com" {
		t.Errorf("SSToSingBox server = %v, want example.com", proxy["server"])
	}

	if proxy["server_port"] != 8388 {
		t.Errorf("SSToSingBox server_port = %v, want 8388", proxy["server_port"])
	}

	// Check method and password
	if proxy["method"] != "aes-256-gcm" {
		t.Errorf("SSToSingBox method = %v, want aes-256-gcm", proxy["method"])
	}

	if proxy["password"] != "password" {
		t.Errorf("SSToSingBox password = %v, want password", proxy["password"])
	}
}

func TestSSToSingBox_InvalidInput(t *testing.T) {
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
			name:    "missing method",
			input:   "ss://@example.com:8388",
			wantErr: true,
		},
		{
			name:    "missing password",
			input:   "ss://aes-256-gcm@example.com:8388",
			wantErr: true,
		},
		{
			name:    "missing server",
			input:   "ss://aes-256-gcm:password@:8388",
			wantErr: true,
		},
		{
			name:    "missing port",
			input:   "ss://aes-256-gcm:password@example.com",
			wantErr: true,
		},
		{
			name:    "invalid port",
			input:   "ss://aes-256-gcm:password@example.com:invalid",
			wantErr: true,
		},
		{
			name:    "invalid base64",
			input:   "ss://invalid-base64!@#@example.com:8388",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SSToSingBox(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SSToSingBox() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSSToSingBox_ValidConfigs(t *testing.T) {
	tests := []struct {
		name     string
		ssLink   string
		expected map[string]interface{}
	}{
		{
			name:   "basic config",
			ssLink: "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@example.com:8388#Test%20Config",
			expected: map[string]interface{}{
				"server":      "example.com",
				"server_port": 8388,
				"method":      "aes-256-gcm",
				"password":    "password",
			},
		},
		{
			name:   "chacha20 config",
			ssLink: "ss://Y2hhY2hhMjAtcG9seTEzMDU6cGFzc3dvcmQ=@example.com:8388#Test%20Config",
			expected: map[string]interface{}{
				"server":      "example.com",
				"server_port": 8388,
				"method":      "chacha20-poly1305",
				"password":    "password",
			},
		},
		{
			name:   "aes-128-gcm config",
			ssLink: "ss://YWVzLTEyOC1nY206cGFzc3dvcmQ=@example.com:8388#Test%20Config",
			expected: map[string]interface{}{
				"server":      "example.com",
				"server_port": 8388,
				"method":      "aes-128-gcm",
				"password":    "password",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := SSToSingBox(tt.ssLink)
			if err != nil {
				t.Fatalf("SSToSingBox() unexpected error: %v", err)
			}

			outbounds, ok := cfg["outbounds"].([]map[string]any)
			if !ok || len(outbounds) == 0 {
				t.Fatal("SSToSingBox() outbounds is not valid")
			}

			proxy := outbounds[0]
			for key, expectedValue := range tt.expected {
				if proxy[key] != expectedValue {
					t.Errorf("SSToSingBox() %s = %v, want %v", key, proxy[key], expectedValue)
				}
			}
		})
	}
}

func TestSSToV2ray(t *testing.T) {
	// Test valid Shadowsocks link
	ssLink := "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@example.com:8388#Test%20Config"

	cfg, err := SSToV2ray(ssLink)
	if err != nil {
		t.Fatalf("SSToV2ray failed: %v", err)
	}

	// Verify required fields exist
	if cfg == nil {
		t.Fatal("SSToV2ray returned nil config")
	}

	// Check for required top-level keys
	requiredKeys := []string{"log", "inbounds", "outbounds"}
	for _, key := range requiredKeys {
		if _, exists := cfg[key]; !exists {
			t.Errorf("SSToV2ray missing required key: %s", key)
		}
	}

	// Verify outbounds structure
	outbounds, ok := cfg["outbounds"].([]map[string]any)
	if !ok {
		t.Fatal("SSToV2ray outbounds is not a slice")
	}

	if len(outbounds) < 1 {
		t.Fatal("SSToV2ray expected at least 1 outbound")
	}

	// Check proxy outbound
	proxy := outbounds[0]
	if proxy["protocol"] != "shadowsocks" {
		t.Errorf("SSToV2ray protocol = %v, want shadowsocks", proxy["protocol"])
	}

	settings, ok := proxy["settings"].(map[string]any)
	if !ok {
		t.Fatal("SSToV2ray settings is not a map")
	}

	servers, ok := settings["servers"].([]map[string]any)
	if !ok || len(servers) == 0 {
		t.Fatal("SSToV2ray servers is not a valid slice")
	}

	server := servers[0]
	if server["address"] != "example.com" {
		t.Errorf("SSToV2ray server address = %v, want example.com", server["address"])
	}

	if server["port"] != 8388 {
		t.Errorf("SSToV2ray server port = %v, want 8388", server["port"])
	}
}

func TestSSToV2ray_InvalidInput(t *testing.T) {
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
			name:    "missing method",
			input:   "ss://@example.com:8388",
			wantErr: true,
		},
		{
			name:    "missing password",
			input:   "ss://aes-256-gcm@example.com:8388",
			wantErr: true,
		},
		{
			name:    "missing server",
			input:   "ss://aes-256-gcm:password@:8388",
			wantErr: true,
		},
		{
			name:    "missing port",
			input:   "ss://aes-256-gcm:password@example.com",
			wantErr: true,
		},
		{
			name:    "invalid port",
			input:   "ss://aes-256-gcm:password@example.com:invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SSToV2ray(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SSToV2ray() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseSSCredentials(t *testing.T) {
	tests := []struct {
		name        string
		ssLink      string
		expectError bool
		expected    struct {
			method   string
			password string
			host     string
			port     int
		}
	}{
		{
			name:        "valid ss link",
			ssLink:      "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@example.com:8388#Test%20Config",
			expectError: false,
			expected: struct {
				method   string
				password string
				host     string
				port     int
			}{
				method:   "aes-256-gcm",
				password: "password",
				host:     "example.com",
				port:     8388,
			},
		},
		{
			name:        "ss link without name",
			ssLink:      "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@example.com:8388",
			expectError: false,
			expected: struct {
				method   string
				password string
				host     string
				port     int
			}{
				method:   "aes-256-gcm",
				password: "password",
				host:     "example.com",
				port:     8388,
			},
		},
		{
			name:        "empty string",
			ssLink:      "",
			expectError: true,
		},
		{
			name:        "invalid prefix",
			ssLink:      "http://example.com",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method, password, host, port, err := parseSSCredentials(tt.ssLink)
			if tt.expectError {
				if err == nil {
					t.Errorf("parseSSCredentials() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("parseSSCredentials() unexpected error: %v", err)
			}

			if method != tt.expected.method {
				t.Errorf("parseSSCredentials() method = %v, want %v", method, tt.expected.method)
			}

			if password != tt.expected.password {
				t.Errorf("parseSSCredentials() password = %v, want %v", password, tt.expected.password)
			}

			if host != tt.expected.host {
				t.Errorf("parseSSCredentials() host = %v, want %v", host, tt.expected.host)
			}

			if port != tt.expected.port {
				t.Errorf("parseSSCredentials() port = %v, want %v", port, tt.expected.port)
			}
		})
	}
}

func TestDecodeBase64String(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    string
	}{
		{
			name:        "valid base64",
			input:       "cGFzc3dvcmQ=",
			expectError: false,
			expected:    "password",
		},
		{
			name:        "base64 without padding",
			input:       "cGFzc3dvcmQ",
			expectError: false,
			expected:    "password",
		},
		{
			name:        "empty string",
			input:       "",
			expectError: false,
			expected:    "",
		},
		{
			name:        "invalid base64",
			input:       "invalid-base64!@#",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := decodeBase64String(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("decodeBase64String() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("decodeBase64String() unexpected error: %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("decodeBase64String() = %v, want %v", string(result), tt.expected)
			}
		})
	}
}
