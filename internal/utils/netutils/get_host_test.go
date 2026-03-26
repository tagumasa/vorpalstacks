package netutils

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"testing"
)

func TestGetHost(t *testing.T) {
	tests := []struct {
		name           string
		host           string
		forwardedHost  string
		trustedProxies []string
		remoteAddr     string
		expected       string
	}{
		{
			name:     "no forwarded host",
			host:     "example.com",
			expected: "example.com",
		},
		{
			name:          "with forwarded host but no trusted proxy",
			host:          "example.com",
			forwardedHost: "cdn.example.com",
			expected:      "example.com",
		},
		{
			name:           "with forwarded host and trusted proxy",
			host:           "example.com",
			forwardedHost:  "cdn.example.com",
			trustedProxies: []string{"192.168.1.1"},
			remoteAddr:     "192.168.1.1:8080",
			expected:       "cdn.example.com",
		},
		{
			name:           "trusted proxy with wildcard",
			host:           "example.com",
			forwardedHost:  "cdn.example.com",
			trustedProxies: []string{"*"},
			remoteAddr:     "192.168.1.1:8080",
			expected:       "cdn.example.com",
		},
		{
			name:           "untrusted proxy",
			host:           "example.com",
			forwardedHost:  "cdn.example.com",
			trustedProxies: []string{"192.168.1.1"},
			remoteAddr:     "10.0.0.1:8080",
			expected:       "example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetTrustedProxies(tt.trustedProxies)
			defer SetTrustedProxies(nil)

			req := &http.Request{
				Host:       tt.host,
				RemoteAddr: tt.remoteAddr,
			}
			if tt.forwardedHost != "" {
				req.Header = http.Header{
					"X-Forwarded-Host": []string{tt.forwardedHost},
				}
			}

			result := GetHost(req)
			if result != tt.expected {
				t.Errorf("GetHost() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetScheme(t *testing.T) {
	tests := []struct {
		name           string
		forwardedProto string
		trustedProxies []string
		remoteAddr     string
		hasTLS         bool
		expected       string
	}{
		{
			name:     "http without TLS",
			expected: "http",
		},
		{
			name:     "https with TLS",
			hasTLS:   true,
			expected: "https",
		},
		{
			name:           "with forwarded proto but no trusted proxy",
			forwardedProto: "https",
			expected:       "http",
		},
		{
			name:           "with forwarded proto and trusted proxy",
			forwardedProto: "https",
			trustedProxies: []string{"192.168.1.1"},
			remoteAddr:     "192.168.1.1:8080",
			expected:       "https",
		},
		{
			name:           "trusted proxy with wildcard",
			forwardedProto: "https",
			trustedProxies: []string{"*"},
			remoteAddr:     "any:8080",
			expected:       "https",
		},
		{
			name:           "untrusted proxy uses TLS",
			forwardedProto: "https",
			trustedProxies: []string{"192.168.1.1"},
			remoteAddr:     "10.0.0.1:8080",
			hasTLS:         true,
			expected:       "https",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetTrustedProxies(tt.trustedProxies)
			defer SetTrustedProxies(nil)

			req := &http.Request{
				RemoteAddr: tt.remoteAddr,
			}
			if tt.hasTLS {
				req.TLS = &tls.ConnectionState{}
			}
			if tt.forwardedProto != "" {
				req.Header = http.Header{
					"X-Forwarded-Proto": []string{tt.forwardedProto},
				}
			}

			result := GetScheme(req)
			if result != tt.expected {
				t.Errorf("GetScheme() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestIsLocalhost(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{"ipv4 loopback", "127.0.0.1", true},
		{"ipv6 loopback", "::1", true},
		{"ipv4 127.x.x.x", "127.0.0.2", true},
		{"ipv4 127.255.255.255", "127.255.255.255", true},
		{"ipv6 full localhost", "0:0:0:0:0:0:0:1", true},
		{"private IP", "192.168.1.1", false},
		{"public IP", "8.8.8.8", false},
		{"empty string", "", false},
		{"localhost string", "localhost", false},
		{"docker IP", "172.17.0.1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsLocalhost(tt.ip)
			if result != tt.expected {
				t.Errorf("IsLocalhost(%q) = %v, want %v", tt.ip, result, tt.expected)
			}
		})
	}
}

func TestURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid URL", "http://example.com/path?query=1", false},
		{"valid HTTPS", "https://example.com", false},
		{"valid with port", "http://example.com:8080", false},
		{"invalid URL", "://example.com", true},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := url.Parse(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("url.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
