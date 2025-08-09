package tui

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// VMessToSingBox converts a vmess:// link into sing-box JSON config
func VMessToSingBox(vmessLink string) (map[string]any, error) {
	if !strings.HasPrefix(vmessLink, "vmess://") {
		return nil, fmt.Errorf("invalid link: must start with vmess://")
	}

	raw := strings.TrimPrefix(vmessLink, "vmess://")
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %w", err)
	}

	var v map[string]string
	if err := json.Unmarshal(decoded, &v); err != nil {
		return nil, fmt.Errorf("vmess JSON parse error: %w", err)
	}

	cfg := map[string]any{
		"log": map[string]any{
			"level": "info",
		},
		"dns": map[string]any{
			"servers": []map[string]any{
				{
					"tag":     "google",
					"address": "tcp://1.1.1.1",
				},
				{
					"tag":     "cloudflare",
					"address": "tcp://1.0.0.1",
				},
			},
		},
		"inbounds": []map[string]any{
			{
				"type": "socks",
				"tag":  "socks-in",
				"listen": map[string]any{
					"address": "127.0.0.1",
					"port":    1080,
				},
				"users": []map[string]any{
					{
						"name": "default",
					},
				},
			},
		},
		"outbounds": []map[string]any{
			{
				"type": "vmess",
				"tag":  "proxy",
				"server": map[string]any{
					"address": v["add"],
					"port":    atoiSafe(v["port"]),
				},
				"uuid": v["id"],
				"security": map[string]any{
					"type": v["scy"],
				},
				"alterId": atoiSafe(v["aid"]),
				"transport": map[string]any{
					"type": v["net"],
					"path": v["path"],
					"host": v["host"],
				},
				"tls": map[string]any{
					"enabled":     v["tls"] == "tls",
					"serverName":  v["host"],
					"fingerprint": v["fp"],
				},
			},
			{
				"type": "direct",
				"tag":  "direct",
			},
		},
	}

	return cfg, nil
}

// VMessToV2ray converts a vmess:// link into V2Ray JSON config
func VMessToV2ray(vmessLink string) (map[string]any, error) {
	if !strings.HasPrefix(vmessLink, "vmess://") {
		return nil, fmt.Errorf("invalid link: must start with vmess://")
	}

	raw := strings.TrimPrefix(vmessLink, "vmess://")
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %w", err)
	}

	var v map[string]string
	if err := json.Unmarshal(decoded, &v); err != nil {
		return nil, fmt.Errorf("vmess JSON parse error: %w", err)
	}

	cfg := map[string]any{
		"log": map[string]any{"loglevel": "none"},
		"dns": map[string]any{
			"servers": []string{"1.1.1.1", "9.9.9.10", "8.8.8.8"},
		},
		"inbounds": []map[string]any{
			{
				"tag":      "socks",
				"port":     1080,
				"listen":   "127.0.0.1",
				"protocol": "socks",
				"sniffing": map[string]any{
					"enabled":      true,
					"destOverride": []string{"http", "tls"},
					"routeOnly":    false,
				},
				"settings": map[string]any{
					"auth":             "noauth",
					"udp":              true,
					"allowTransparent": false,
				},
			},
		},
		"outbounds": []map[string]any{
			{
				"tag":      "proxy",
				"protocol": "vmess",
				"settings": map[string]any{
					"vnext": []map[string]any{
						{
							"address": v["add"],
							"port":    atoiSafe(v["port"]),
							"users": []map[string]any{
								{
									"id":       v["id"],
									"alterId":  atoiSafe(v["aid"]),
									"email":    "t@t.tt",
									"security": v["scy"],
								},
							},
						},
					},
				},
				"streamSettings": map[string]any{
					"network":  v["net"],
					"security": v["tls"],
					"tlsSettings": map[string]any{
						"allowInsecure": false,
						"serverName":    v["host"],
						"fingerprint":   v["fp"],
					},
					"wsSettings": map[string]any{
						"path": v["path"],
						"host": v["host"],
						"headers": map[string]any{
							"Host": v["host"],
						},
					},
				},
				"mux": map[string]any{
					"enabled":     false,
					"concurrency": -1,
				},
			},
			{
				"tag":      "direct",
				"protocol": "freedom",
				"settings": map[string]any{
					"domainStrategy": "UseIP",
					"userLevel":      0,
				},
			},
		},
	}

	return cfg, nil
}

// atoiSafe safely converts string to int
func atoiSafe(s string) int {
	if s == "" {
		return 0
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	
	return i
} 