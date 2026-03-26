package netutils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name           string
		remoteAddr     string
		xForwardedFor  string
		trustedProxies []string
		expected       string
	}{
		{
			name:       "no headers, use remote addr",
			remoteAddr: "192.168.1.1:8080",
			expected:   "192.168.1.1",
		},
		{
			name:           "trusted proxy with X-Forwarded-For",
			remoteAddr:     "10.0.0.1:8080",
			xForwardedFor:  "1.2.3.4, 5.6.7.8",
			trustedProxies: []string{"10.0.0.1"},
			expected:       "1.2.3.4",
		},
		{
			name:          "untrusted proxy ignores X-Forwarded-For",
			remoteAddr:    "10.0.0.1:8080",
			xForwardedFor: "1.2.3.4",
			expected:      "10.0.0.1",
		},
		{
			name:       "invalid remote addr returns as-is",
			remoteAddr: "invalid-addr",
			expected:   "invalid-addr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetTrustedProxies(tt.trustedProxies)
			defer SetTrustedProxies(nil)

			req := &http.Request{
				RemoteAddr: tt.remoteAddr,
			}
			if tt.xForwardedFor != "" {
				req.Header = http.Header{
					"X-Forwarded-For": []string{tt.xForwardedFor},
				}
			}

			result := GetClientIP(req)
			if result != tt.expected {
				t.Errorf("GetClientIP() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSetCORSHeaders(t *testing.T) {
	tests := []struct {
		name             string
		allowedOrigins   string
		allowedMethods   string
		allowedHeaders   string
		allowCredentials bool
		wantErr          bool
		checkOrigin      string
		checkMethod      string
		checkHeader      string
		checkCredentials string
	}{
		{
			name:             "basic cors",
			allowedOrigins:   "*",
			allowedMethods:   "GET,POST",
			allowedHeaders:   "Content-Type",
			wantErr:          false,
			checkOrigin:      "*",
			checkMethod:      "GET,POST",
			checkHeader:      "Content-Type",
			checkCredentials: "",
		},
		{
			name:             "with credentials",
			allowedOrigins:   "https://example.com",
			allowedMethods:   "GET,POST",
			allowedHeaders:   "Content-Type",
			allowCredentials: true,
			wantErr:          false,
			checkOrigin:      "https://example.com",
			checkMethod:      "GET,POST",
			checkHeader:      "Content-Type",
			checkCredentials: "true",
		},
		{
			name:             "credentials with wildcard error",
			allowedOrigins:   "*",
			allowedMethods:   "GET,POST",
			allowedHeaders:   "Content-Type",
			allowCredentials: true,
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := SetCORSHeaders(w, tt.allowedOrigins, tt.allowedMethods, tt.allowedHeaders, tt.allowCredentials)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetCORSHeaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if w.Header().Get("Access-Control-Allow-Origin") != tt.checkOrigin {
					t.Errorf("Origin = %v, want %v", w.Header().Get("Access-Control-Allow-Origin"), tt.checkOrigin)
				}
				if w.Header().Get("Access-Control-Allow-Methods") != tt.checkMethod {
					t.Errorf("Methods = %v, want %v", w.Header().Get("Access-Control-Allow-Methods"), tt.checkMethod)
				}
				if w.Header().Get("Access-Control-Allow-Headers") != tt.checkHeader {
					t.Errorf("Headers = %v, want %v", w.Header().Get("Access-Control-Allow-Headers"), tt.checkHeader)
				}
				if tt.checkCredentials != "" && w.Header().Get("Access-Control-Allow-Credentials") != tt.checkCredentials {
					t.Errorf("Credentials = %v, want %v", w.Header().Get("Access-Control-Allow-Credentials"), tt.checkCredentials)
				}
			}
		})
	}
}
