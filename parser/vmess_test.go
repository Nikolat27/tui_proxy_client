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
