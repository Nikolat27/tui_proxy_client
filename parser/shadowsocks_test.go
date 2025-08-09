package parser

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestSSToV2ray(t *testing.T) {
	// Your provided Shadowsocks link
	ssLink := "ss://MjAyMi1ibGFrZTMtYWVzLTI1Ni1nY206dHEwdE1CVHVCV0t5SWMrMjJCUUNNTW84SjJoeHZTcVJBWDZlb0JsSVZTVT06aWZDMDh4azBEL1NQaWowbVptZld0bXFkUUZYRzljeHdIcWZWUUNadVBSaz0%3D@89.44.242.222:20379#%40InvoProxy"

	cfg, err := SSToV2ray(ssLink)
	if err != nil {
		t.Fatalf("SSToV2ray failed: %v", err)
	}

	// Save config to file
	fileName := "test.json"
	data, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.WriteFile(fileName, data, 0644); err != nil {
		t.Fatalf("failed to write %s: %v", fileName, err)
	}
	t.Logf("Saved config to %s", fileName)

	// Run v2ray with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "v2ray", "run", fileName)
	output, err := cmd.CombinedOutput()
	t.Logf("V2Ray output:\n%s", string(output))

	if ctx.Err() == context.DeadlineExceeded {
		// OK — config ran without immediate fatal errors
	} else if err != nil {
		t.Fatalf("v2ray run failed: %v", err)
	}
}
