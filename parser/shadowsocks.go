package parser

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// decodeBase64String safely decodes base64 with or without padding
func decodeBase64String(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	if m := len(s) % 4; m != 0 {
		s += strings.Repeat("=", 4-m)
	}
	if data, err := base64.StdEncoding.DecodeString(s); err == nil {
		return data, nil
	}
	return base64.RawStdEncoding.DecodeString(s)
}

// SS -> Shadow Socks
// parseSSCredentials extracts method, password, host, and port from an ss:// link
func parseSSCredentials(ssLink string) (method, password, host string, port int, err error) {
	raw := strings.TrimSpace(strings.TrimPrefix(ssLink, "ss://"))

	// Case 1: base64(method:password)@host:port
	if strings.Contains(raw, "@") {
		parts := strings.SplitN(raw, "@", 2)
		left := strings.TrimSpace(parts[0]) // base64 part only
		right := parts[1]

		// URL-decode before base64 decode (important for Shadowsocks 2022 links)
		// Handle multiple layers of encoding
		for {
			if decodedURL, errURL := url.QueryUnescape(left); errURL == nil && decodedURL != left {
				left = decodedURL
			} else {
				break
			}
		}

		// Decode left part
		decoded, decErr := decodeBase64String(left)
		if decErr != nil {
			err = fmt.Errorf("base64 decode error: %w", decErr)
			return
		}

		// Split credentials
		credParts := strings.SplitN(string(decoded), ":", 2)
		if len(credParts) != 2 {
			err = fmt.Errorf("invalid credentials format in ss:// link")
			return
		}
		method, password = credParts[0], credParts[1]

		// Strip tag/path from right
		right = strings.SplitN(right, "#", 2)[0]
		right = strings.SplitN(right, "/", 2)[0]

		hp := strings.SplitN(right, ":", 2)
		if len(hp) != 2 {
			err = fmt.Errorf("invalid host:port in ss:// link")
			return
		}
		host = hp[0]
		port, err = strconv.Atoi(hp[1])
		return
	}

	// Case 2: fully base64(method:password@host:port)
	// URL-decode before base64 decode - handle multiple layers
	for {
		if decodedURL, errURL := url.QueryUnescape(raw); errURL == nil && decodedURL != raw {
			raw = decodedURL
		} else {
			break
		}
	}

	decoded, decErr := decodeBase64String(raw)
	if decErr != nil {
		err = fmt.Errorf("base64 decode error: %w", decErr)
		return
	}

	link := strings.SplitN(string(decoded), "#", 2)[0]
	link = strings.SplitN(link, "/", 2)[0]

	u, parseErr := url.Parse("ss://" + link)
	if parseErr != nil {
		err = fmt.Errorf("invalid parsed link: %w", parseErr)
		return
	}

	userParts := strings.SplitN(u.User.String(), ":", 2)
	if len(userParts) != 2 {
		err = fmt.Errorf("invalid user info in ss:// link")
		return
	}
	method, password = userParts[0], userParts[1]
	host = u.Hostname()
	port, err = strconv.Atoi(u.Port())
	return
}

// buildSSV2rayConfig builds the final V2Ray JSON config
func buildSSV2rayConfig(method, password, host string, port int) map[string]any {
	return map[string]any{
		"log": map[string]any{"loglevel": "none"},
		"inbounds": []map[string]any{
			{
				"tag":      "socks",
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
				"tag":      "proxy",
				"protocol": "shadowsocks",
				"settings": map[string]any{
					"servers": []map[string]any{
						{
							"address":  host,
							"port":     port,
							"method":   method,
							"password": password,
							"level":    0,
							"ota":      false,
						},
					},
				},
			},
		},
	}
}

// SSToV2ray converts an ss:// link to a V2Ray Shadowsocks config
func SSToV2ray(ssLink string) (map[string]any, error) {
	if !strings.HasPrefix(ssLink, "ss://") {
		return nil, fmt.Errorf("invalid link: must start with ss://")
	}
	method, password, host, port, err := parseSSCredentials(ssLink)
	if err != nil {
		return nil, err
	}
	return buildSSV2rayConfig(method, password, host, port), nil
}

// SSToSingBox converts an ss:// link to a sing-box Shadowsocks config
func SSToSingBox(ssLink string) (map[string]any, error) {
	if !strings.HasPrefix(ssLink, "ss://") {
		return nil, fmt.Errorf("invalid link: must start with ss://")
	}
	method, password, host, port, err := parseSSCredentials(ssLink)
	if err != nil {
		return nil, err
	}

	cfg := map[string]any{
		"log": map[string]any{
			"level": "info",
		},
		"inbounds": []map[string]any{
			{
				"type":        "socks",
				"tag":         "socks-in",
				"listen":      "127.0.0.1",
				"listen_port": 1080,
				"sniff":       true,
			},
		},
		"outbounds": []map[string]any{
			{
				"type":        "shadowsocks",
				"tag":         "proxy",
				"server":      host,
				"server_port": port,
				"method":      method,
				"password":    password,
			},
			{
				"type": "direct",
				"tag":  "direct",
			},
		},
	}

	return cfg, nil
}
