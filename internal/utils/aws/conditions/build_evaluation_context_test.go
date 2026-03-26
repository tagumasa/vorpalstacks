package conditions

import (
	"net/http"
	"net/url"
	"testing"
)

func TestBuildEvaluationContextFromMap(t *testing.T) {
	tests := []struct {
		name           string
		values         map[string]string
		checkSourceIP  bool
		checkUserAgent bool
	}{
		{
			name:          "nil values",
			values:        nil,
			checkSourceIP: false,
		},
		{
			name:          "empty values",
			values:        map[string]string{},
			checkSourceIP: false,
		},
		{
			name: "with source IP",
			values: map[string]string{
				"SourceIP": "192.168.1.1",
			},
			checkSourceIP: true,
		},
		{
			name: "with user agent",
			values: map[string]string{
				"UserAgent": "Mozilla/5.0",
			},
			checkUserAgent: true,
		},
		{
			name: "with all fields",
			values: map[string]string{
				"SourceIP":     "192.168.1.1",
				"UserAgent":    "Mozilla/5.0",
				"Referer":      "https://example.com",
				"PrincipalARN": "arn:aws:iam::123456789012:user/test",
				"PrincipalID":  "AIDAEXAMPLE",
				"ResourceARN":  "arn:aws:s3:::bucket/*",
				"Action":       "s3:GetObject",
			},
			checkSourceIP:  true,
			checkUserAgent: true,
		},
		{
			name: "with unknown fields",
			values: map[string]string{
				"UnknownField": "value",
				"SourceIP":     "192.168.1.1",
			},
			checkSourceIP: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := BuildEvaluationContextFromMap(tt.values)

			if ctx == nil {
				t.Fatal("BuildEvaluationContextFromMap() returned nil")
			}

			if ctx.RequestHeaders == nil {
				t.Error("RequestHeaders should not be nil")
			}

			if ctx.Environment == nil {
				t.Error("Environment should not be nil")
			}

			if tt.checkSourceIP && ctx.SourceIP != "192.168.1.1" {
				t.Errorf("SourceIP = %q, want %q", ctx.SourceIP, "192.168.1.1")
			}

			if tt.checkUserAgent && ctx.UserAgent != "Mozilla/5.0" {
				t.Errorf("UserAgent = %q, want %q", ctx.UserAgent, "Mozilla/5.0")
			}

			if tt.values != nil {
				for k, v := range tt.values {
					if ctx.RequestHeaders[k] != v {
						t.Errorf("RequestHeaders[%q] = %q, want %q", k, ctx.RequestHeaders[k], v)
					}
				}
			}
		})
	}
}

func TestBuildEvaluationContext(t *testing.T) {
	tests := []struct {
		name          string
		request       *http.Request
		wantSourceIP  string
		wantUserAgent string
	}{
		{
			name:    "nil request",
			request: nil,
		},
		{
			name: "basic request",
			request: &http.Request{
				RemoteAddr: "192.168.1.1:12345",
				Header: http.Header{
					"User-Agent": []string{"TestAgent/1.0"},
				},
			},
			wantSourceIP:  "192.168.1.1",
			wantUserAgent: "TestAgent/1.0",
		},
		{
			name: "request with X-Forwarded-For",
			request: &http.Request{
				RemoteAddr: "10.0.0.1:12345",
				Header: http.Header{
					"X-Forwarded-For": []string{"203.0.113.1, 70.41.3.18"},
				},
			},
			wantSourceIP: "203.0.113.1",
		},
		{
			name: "request with X-Real-IP",
			request: &http.Request{
				RemoteAddr: "10.0.0.1:12345",
				Header: http.Header{
					"X-Real-Ip": []string{"198.51.100.1"},
				},
			},
			wantSourceIP: "198.51.100.1",
		},
		{
			name: "IPv6 address",
			request: &http.Request{
				RemoteAddr: "[2001:db8::1]:12345",
			},
			wantSourceIP: "2001:db8::1",
		},
		{
			name: "request with referer",
			request: &http.Request{
				RemoteAddr: "192.168.1.1:12345",
				Header: http.Header{
					"Referer": []string{"https://example.com/page"},
				},
			},
			wantSourceIP: "192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := BuildEvaluationContext(tt.request)

			if ctx == nil {
				t.Fatal("BuildEvaluationContext() returned nil")
			}

			if ctx.CurrentTime.IsZero() {
				t.Error("CurrentTime should not be zero")
			}

			if ctx.RequestHeaders == nil {
				t.Error("RequestHeaders should not be nil")
			}

			if ctx.Environment == nil {
				t.Error("Environment should not be nil")
			}

			if tt.wantSourceIP != "" && ctx.SourceIP != tt.wantSourceIP {
				t.Errorf("SourceIP = %q, want %q", ctx.SourceIP, tt.wantSourceIP)
			}

			if tt.wantUserAgent != "" && ctx.UserAgent != tt.wantUserAgent {
				t.Errorf("UserAgent = %q, want %q", ctx.UserAgent, tt.wantUserAgent)
			}

			if tt.request != nil && tt.request.Header.Get("Referer") != "" {
				if ctx.Referer != tt.request.Header.Get("Referer") {
					t.Errorf("Referer = %q, want %q", ctx.Referer, tt.request.Header.Get("Referer"))
				}
			}
		})
	}
}

func TestExtractSourceIP(t *testing.T) {
	tests := []struct {
		name     string
		request  *http.Request
		expected string
	}{
		{
			name: "X-Forwarded-For single IP",
			request: &http.Request{
				RemoteAddr: "10.0.0.1:12345",
				Header: http.Header{
					"X-Forwarded-For": []string{"203.0.113.1"},
				},
			},
			expected: "203.0.113.1",
		},
		{
			name: "X-Forwarded-For multiple IPs",
			request: &http.Request{
				RemoteAddr: "10.0.0.1:12345",
				Header: http.Header{
					"X-Forwarded-For": []string{"203.0.113.1, 70.41.3.18, 150.172.238.178"},
				},
			},
			expected: "203.0.113.1",
		},
		{
			name: "X-Real-IP takes precedence over RemoteAddr",
			request: &http.Request{
				RemoteAddr: "10.0.0.1:12345",
				Header: http.Header{
					"X-Real-Ip": []string{"198.51.100.1"},
				},
			},
			expected: "198.51.100.1",
		},
		{
			name: "X-Forwarded-For takes precedence over X-Real-IP",
			request: &http.Request{
				RemoteAddr: "10.0.0.1:12345",
				Header: http.Header{
					"X-Forwarded-For": []string{"203.0.113.1"},
					"X-Real-Ip":       []string{"198.51.100.1"},
				},
			},
			expected: "203.0.113.1",
		},
		{
			name: "RemoteAddr fallback",
			request: &http.Request{
				RemoteAddr: "192.168.1.100:54321",
			},
			expected: "192.168.1.100",
		},
		{
			name: "IPv6 RemoteAddr",
			request: &http.Request{
				RemoteAddr: "[2001:db8:85a3::8a2e:370:7334]:443",
			},
			expected: "2001:db8:85a3::8a2e:370:7334",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSourceIP(tt.request)
			if result != tt.expected {
				t.Errorf("extractSourceIP() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestExtractHeaders(t *testing.T) {
	tests := []struct {
		name     string
		headers  http.Header
		expected map[string]string
	}{
		{
			name:     "nil headers",
			headers:  nil,
			expected: map[string]string{},
		},
		{
			name:     "empty headers",
			headers:  http.Header{},
			expected: map[string]string{},
		},
		{
			name: "single value headers",
			headers: http.Header{
				"Content-Type":  []string{"application/json"},
				"Authorization": []string{"Bearer token123"},
			},
			expected: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer token123",
			},
		},
		{
			name: "multiple values - first used",
			headers: http.Header{
				"Accept": []string{"application/json", "text/html"},
			},
			expected: map[string]string{
				"Accept": "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractHeaders(tt.headers)
			if len(result) != len(tt.expected) {
				t.Errorf("extractHeaders() returned %d headers, want %d", len(result), len(tt.expected))
			}
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("extractHeaders()[%q] = %q, want %q", k, result[k], v)
				}
			}
		})
	}
}

func TestBuildEvaluationContextIntegration(t *testing.T) {
	req := &http.Request{
		Method:     "POST",
		URL:        &url.URL{Path: "/api/v1/resource"},
		RemoteAddr: "192.168.1.50:8080",
		Header: http.Header{
			"User-Agent":      []string{"MyClient/2.0"},
			"Referer":         []string{"https://example.com/origin"},
			"X-Forwarded-For": []string{"203.0.113.50"},
			"Content-Type":    []string{"application/json"},
			"Authorization":   []string{"Bearer token"},
		},
	}

	ctx := BuildEvaluationContext(req)

	if ctx.SourceIP != "203.0.113.50" {
		t.Errorf("SourceIP = %q, want %q", ctx.SourceIP, "203.0.113.50")
	}

	if ctx.UserAgent != "MyClient/2.0" {
		t.Errorf("UserAgent = %q, want %q", ctx.UserAgent, "MyClient/2.0")
	}

	if ctx.Referer != "https://example.com/origin" {
		t.Errorf("Referer = %q, want %q", ctx.Referer, "https://example.com/origin")
	}

	if ctx.RequestHeaders["Content-Type"] != "application/json" {
		t.Errorf("RequestHeaders[Content-Type] = %q, want %q", ctx.RequestHeaders["Content-Type"], "application/json")
	}

	if ctx.RequestHeaders["Authorization"] != "Bearer token" {
		t.Errorf("RequestHeaders[Authorization] = %q, want %q", ctx.RequestHeaders["Authorization"], "Bearer token")
	}
}
