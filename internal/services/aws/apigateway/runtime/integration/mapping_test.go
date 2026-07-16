package integration

import (
	"testing"
)

func TestApplyResponseParameterMapping(t *testing.T) {
	tests := []struct {
		name           string
		responseParams map[string]string
		currentHeaders map[string]string
		responseBody   string
		expected       map[string]string
	}{
		{
			name: "integration header to method header",
			responseParams: map[string]string{
				"method.response.header.Location": "integration.response.header.Location",
			},
			currentHeaders: map[string]string{"Location": "https://example.com"},
			expected:       map[string]string{"Location": "https://example.com"},
		},
		{
			name: "integration header with different name",
			responseParams: map[string]string{
				"method.response.header.X-Backend-Id": "integration.response.header.RequestId",
			},
			currentHeaders: map[string]string{"RequestId": "abc-123"},
			expected:       map[string]string{"RequestId": "abc-123", "X-Backend-Id": "abc-123"},
		},
		{
			name: "static value in single quotes",
			responseParams: map[string]string{
				"method.response.header.Content-Type": "'application/xml'",
			},
			currentHeaders: map[string]string{},
			expected:       map[string]string{"Content-Type": "application/xml"},
		},
		{
			name: "body JSON extraction",
			responseParams: map[string]string{
				"method.response.header.X-Id": "integration.response.body.id",
			},
			currentHeaders: map[string]string{},
			responseBody:   `{"id": "12345"}`,
			expected:       map[string]string{"X-Id": "12345"},
		},
		{
			name: "body nested JSON extraction",
			responseParams: map[string]string{
				"method.response.header.X-Region": "integration.response.body.data.region",
			},
			currentHeaders: map[string]string{},
			responseBody:   `{"data": {"region": "us-east-1"}}`,
			expected:       map[string]string{"X-Region": "us-east-1"},
		},
		{
			name: "false value removes header",
			responseParams: map[string]string{
				"method.response.header.Location": "false",
			},
			currentHeaders: map[string]string{"Location": "https://example.com"},
			expected:       map[string]string{},
		},
		{
			name:           "nil params returns current headers unchanged",
			responseParams: nil,
			currentHeaders: map[string]string{"X-Custom": "val"},
			expected:       map[string]string{"X-Custom": "val"},
		},
		{
			name: "existing headers preserved",
			responseParams: map[string]string{
				"method.response.header.X-New": "integration.response.header.X-Source",
			},
			currentHeaders: map[string]string{"X-Source": "src-val", "X-Existing": "keep"},
			expected:       map[string]string{"X-Source": "src-val", "X-Existing": "keep", "X-New": "src-val"},
		},
		{
			name: "non-mapping key is skipped",
			responseParams: map[string]string{
				"invalid.key": "some-value",
			},
			currentHeaders: map[string]string{"X-Keep": "1"},
			expected:       map[string]string{"X-Keep": "1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyResponseParameterMapping(tt.responseParams, tt.currentHeaders, tt.responseBody)
			if len(result) != len(tt.expected) {
				t.Errorf("result has %d keys, expected %d (got=%v, want=%v)",
					len(result), len(tt.expected), result, tt.expected)
				return
			}
			for k, v := range tt.expected {
				if got, ok := result[k]; !ok {
					t.Errorf("missing key %q", k)
				} else if got != v {
					t.Errorf("key %q: got %q, want %q", k, got, v)
				}
			}
		})
	}
}

func TestIsBinaryContentType(t *testing.T) {
	tests := []struct {
		contentType      string
		binaryMediaTypes []string
		expected         bool
	}{
		{"image/png", []string{"image/png"}, true},
		{"image/jpeg", []string{"image/png"}, false},
		{"image/jpeg", []string{"image/*"}, true},
		{"image/png", []string{"image/*"}, true},
		{"application/octet-stream", []string{"*/*"}, true},
		{"text/html", []string{"*/*"}, true},
		{"image/gif", []string{"image/*", "application/pdf"}, true},
		{"application/pdf", []string{"image/*", "application/pdf"}, true},
		{"text/plain", []string{"image/*"}, false},
		{"image/png", []string{}, false},
		{"", []string{"image/*"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.contentType, func(t *testing.T) {
			got := isBinaryContentType(tt.contentType, tt.binaryMediaTypes)
			if got != tt.expected {
				t.Errorf("isBinaryContentType(%q, %v) = %v, want %v",
					tt.contentType, tt.binaryMediaTypes, got, tt.expected)
			}
		})
	}
}

func TestSelectResponseContentType(t *testing.T) {
	templates := map[string]string{
		"application/json": "{\"id\":$input.path('$.id')}",
		"application/xml":  "<id>$input.path('$.id')</id>",
		"text/html":        "<p>$input.path('$.id')</p>",
	}

	tests := []struct {
		name            string
		templates       map[string]string
		responseHeaders map[string]string
		requestHeaders  map[string]string
		expected        string
	}{
		{
			name:            "response Content-Type takes priority",
			templates:       templates,
			responseHeaders: map[string]string{"Content-Type": "application/xml"},
			requestHeaders:  map[string]string{"Accept": "text/html"},
			expected:        "application/xml",
		},
		{
			name:            "Accept header fallback",
			templates:       templates,
			responseHeaders: map[string]string{},
			requestHeaders:  map[string]string{"Accept": "text/html"},
			expected:        "text/html",
		},
		{
			name:            "default application/json",
			templates:       templates,
			responseHeaders: map[string]string{},
			requestHeaders:  map[string]string{},
			expected:        "application/json",
		},
		{
			name:            "response Content-Type not in templates, falls to Accept",
			templates:       templates,
			responseHeaders: map[string]string{"Content-Type": "text/csv"},
			requestHeaders:  map[string]string{"Accept": "application/xml"},
			expected:        "application/xml",
		},
		{
			name:            "no match, defaults to application/json",
			templates:       templates,
			responseHeaders: map[string]string{"Content-Type": "text/csv"},
			requestHeaders:  map[string]string{"Accept": "text/csv"},
			expected:        "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := selectResponseContentType(tt.templates, tt.responseHeaders, tt.requestHeaders)
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}
