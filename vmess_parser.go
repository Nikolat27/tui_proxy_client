package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// VMessToV2ray converts a vmess:// link into your JSON config struct (v2ray socks5)
func VMessToV2ray(vmessLink string) (map[string]any, error) {
	if !strings.HasPrefix(vmessLink, "vmess://") {
		return nil, fmt.Errorf("invalid link: must start with vmess://")
	}

	// Strip prefix & decode Base64
	raw := strings.TrimPrefix(vmessLink, "vmess://")
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %w", err)
	}

	// Parse VMess JSON
	var v map[string]string
	if err := json.Unmarshal(decoded, &v); err != nil {
		return nil, fmt.Errorf("vmess JSON parse error: %w", err)
	}

	// Build your JSON config
	cfg := map[string]any{
		"log": map[string]any{"loglevel": "none"},
		"dns": map[string]any{
			"servers": []string{"1.1.1.1", "9.9.9.10", "8.8.8.8"},
		},
		"inbounds": []map[string]any{
			{
				"tag":      "socks",
				"port":     1088,
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
			{
				"tag":      "block",
				"protocol": "blackhole",
			},
		},
		"routing": map[string]any{
			"domainStrategy": "IPOnDemand",
			"rules": []map[string]any{
				{"type": "field", "inboundTag": []string{"api"}, "outboundTag": "api"},
				{"type": "field", "outboundTag": "direct", "ip": []string{"8.8.8.8"}},
				{"type": "field", "port": "443", "network": "udp", "outboundTag": "block"},
				{"type": "field", "outboundTag": "direct", "protocol": []string{"bittorrent"}},
				{"type": "field", "port": "0-65535", "outboundTag": "proxy"},
			},
		},
	}

	return cfg, nil
}

// VMessToSingBox converts a vmess:// link into your JSON config struct (sing-box socks5)
func VMessToSingBox(vmessLink string) (map[string]any, error) {
	if !strings.HasPrefix(vmessLink, "vmess://") {
		return nil, fmt.Errorf("invalid link: must start with vmess://")
	}

	// Decode Base64 part
	raw := strings.TrimPrefix(vmessLink, "vmess://")
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %w", err)
	}

	// Parse the vmess JSON inside
	var v map[string]string
	if err := json.Unmarshal(decoded, &v); err != nil {
		return nil, fmt.Errorf("vmess JSON parse error: %w", err)
	}

	// Build config exactly as given
	cfg := map[string]any{
		"log": map[string]any{
			"level": "info",
		},
		"inbounds": []map[string]any{
			{
				"type":                       "socks",
				"tag":                        "socks-in",
				"listen":                     "127.0.0.1",
				"listen_port":                1080,
				"sniff":                      true,
				"sniff_override_destination": true,
			},
		},
		"outbounds": []map[string]any{
			{
				"type":        "vmess",
				"tag":         "proxy-out",
				"server":      v["add"],
				"server_port": atoiSafe(v["port"]),
				"uuid":        v["id"],
				"security":    v["scy"],
				"tls": map[string]any{
					"enabled":     strings.ToLower(v["tls"]) == "tls",
					"server_name": v["host"],
					"utls": map[string]any{
						"enabled":     v["fp"] != "",
						"fingerprint": v["fp"],
					},
				},
				"transport": map[string]any{
					"type": v["net"],
					"path": v["path"],
					"headers": map[string]any{
						"Host": v["host"],
					},
				},
			},
			{
				"type": "direct",
				"tag":  "direct",
			},
			{
				"type": "block",
				"tag":  "block",
			},
		},
		"route": map[string]any{
			"rules": []map[string]any{
				{
					"outbound": "direct",
					"ip_cidr":  []string{"8.8.8.8/32"},
				},
				{
					"outbound": "block",
					"port":     []int{443},
					"network":  "udp",
				},
				{
					"outbound": "direct",
					"protocol": []string{"bittorrent"},
				},
				{
					"outbound": "proxy-out",
				},
			},
		},
	}

	return cfg, nil
}

// atoiSafe converts string to int safely
func atoiSafe(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
