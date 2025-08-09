package parser

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestVLESSToSingBox(t *testing.T) {
	vlessLink := "vless://82a15da6-96cf-45b4-87f4-fe85a36113b5@filmnet.filimnet.com:80?encryption=none&security=none&type=xhttp&host=abriconf.global.ssl.fastly.net&path=%2Fcdn%2Fassets%3Fv%3D12&mode=packet-up"

	cfg, err := VLESSToSingBox(vlessLink)
	if err != nil {
		t.Fatalf("VLESSToSingBox failed: %v", err)
	}

	fileName := "test.json"
	fileData, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.WriteFile(fileName, fileData, 0644); err != nil {
		t.Fatalf("failed to write test.json: %v", err)
	}

	// Run with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sing-box", "run", "-c", fileName)
	output, err := cmd.CombinedOutput()
	t.Logf("Sing-box output:\n%s", string(output))

	if ctx.Err() == context.DeadlineExceeded {
		// Process ran without fatal config errors for 2 seconds → config is valid
	} else if err != nil {
		t.Fatalf("sing-box run failed: %v", err)
	}

	// Basic field checks
	outbounds := cfg["outbounds"].([]map[string]any)
	ob := outbounds[0]
	if ob["type"] != "vless" {
		t.Errorf("expected type 'vless', got %v", ob["type"])
	}
}

func TestVLESSToV2Ray(t *testing.T) {
	vlessLink := "vless://82a15da6-96cf-45b4-87f4-fe85a36113b5@filmnet.filimnet.com:80?encryption=none&security=none&type=ws&host=abriconf.global.ssl.fastly.net&path=%2Fcdn%2Fassets%3Fv%3D12"

	cfg, err := VLESSToV2Ray(vlessLink)
	if err != nil {
		t.Fatalf("VLESSToV2Ray failed: %v", err)
	}

	fileName := "test.json"
	data, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.WriteFile(fileName, data, 0644); err != nil {
		t.Fatalf("failed to write %s: %v", fileName, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "v2ray", "run", fileName)
	output, err := cmd.CombinedOutput()
	t.Logf("V2Ray output:\n%s", string(output))

	if ctx.Err() == context.DeadlineExceeded {
		// No fatal config error detected in 2 seconds
	} else if err != nil {
		t.Fatalf("v2ray run failed: %v", err)
	}
}
