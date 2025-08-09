package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func VLESSToSingBox(vlessLink string) (map[string]any, error) {
	if !strings.HasPrefix(vlessLink, "vless://") {
		return nil, fmt.Errorf("invalid link: must start with vless://")
	}

	u, err := url.Parse(vlessLink)
	if err != nil {
		return nil, fmt.Errorf("invalid vless link: %w", err)
	}

	uuid := u.User.Username()
	if uuid == "" {
		return nil, fmt.Errorf("missing UUID in VLESS link")
	}

	host := u.Hostname()
	portStr := u.Port()

	q := u.Query()
	security := strings.ToLower(q.Get("security"))
	if portStr == "" {
		if security == "tls" {
			portStr = "443"
		} else {
			portStr = "80"
		}
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	sni := q.Get("sni")
	fp := q.Get("fp")
	network := strings.ToLower(q.Get("type"))
	path := q.Get("path")
	hostHeader := q.Get("host")

	// Map or validate transport type
	switch network {
	case "xhttp":
		network = "http" // map unsupported to supported
	case "tcp", "ws", "grpc", "http", "quic", "h2", "splithttp":
		// valid
	default:
		return nil, fmt.Errorf("unsupported transport type: %s", network)
	}

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
				"type":        "vless",
				"tag":         "proxy-out",
				"server":      host,
				"server_port": port,
				"uuid":        uuid,
				"flow":        q.Get("flow"),
				"tls": map[string]any{
					"enabled":     security == "tls",
					"server_name": sni,
					"utls": map[string]any{
						"enabled":     fp != "",
						"fingerprint": fp,
					},
				},
				"transport": map[string]any{
					"type": network,
					"path": path,
					"headers": map[string]any{
						"Host": hostHeader,
					},
				},
			},
		},
	}

	return cfg, nil
}

func VLESSToV2Ray(vlessLink string) (map[string]any, error) {
	if !strings.HasPrefix(vlessLink, "vless://") {
		return nil, fmt.Errorf("invalid link: must start with vless://")
	}

	u, err := url.Parse(vlessLink)
	if err != nil {
		return nil, fmt.Errorf("invalid VLESS link: %w", err)
	}

	uuid := u.User.Username()
	if uuid == "" {
		return nil, fmt.Errorf("missing UUID in VLESS link")
	}

	host := u.Hostname()
	portStr := u.Port()
	if portStr == "" {
		portStr = "443"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	q := u.Query()
	security := q.Get("security")
	network := q.Get("type")
	path := q.Get("path")
	hostHeader := q.Get("host")
	sni := q.Get("sni")

	cfg := map[string]any{
		"log": map[string]any{
			"loglevel": "info",
		},
		"inbounds": []map[string]any{
			{
				"tag":      "socks-in",
				"port":     1080,
				"listen":   "127.0.0.1",
				"protocol": "socks",
				"sniffing": map[string]any{
					"enabled":      true,
					"destOverride": []string{"http", "tls"},
				},
				"settings": map[string]any{
					"auth": "noauth",
					"udp":  true,
				},
			},
		},
		"outbounds": []map[string]any{
			{
				"tag":      "proxy-out",
				"protocol": "vless",
				"settings": map[string]any{
					"vnext": []map[string]any{
						{
							"address": host,
							"port":    port,
							"users": []map[string]any{
								{
									"id":         uuid,
									"encryption": "none",
									"flow":       q.Get("flow"),
								},
							},
						},
					},
				},
				"streamSettings": map[string]any{
					"network":  network,
					"security": security,
					"tlsSettings": map[string]any{
						"serverName": sni,
					},
					"wsSettings": map[string]any{
						"path": path,
						"headers": map[string]any{
							"Host": hostHeader,
						},
					},
				},
			},
		},
	}

	return cfg, nil
}
