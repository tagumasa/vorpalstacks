package validators

import "testing"

func TestIsPreflightRequest(t *testing.T) {
	tests := []struct {
		name   string
		method string
		valid  bool
	}{
		{"OPTIONS", "OPTIONS", true},
		{"GET", "GET", false},
		{"POST", "POST", false},
		{"PUT", "PUT", false},
		{"DELETE", "DELETE", false},
		{"Empty", "", false},
		{"Lowercase options", "options", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPreflightRequest(tt.method); got != tt.valid {
				t.Errorf("IsPreflightRequest(%q) = %v, want %v", tt.method, got, tt.valid)
			}
		})
	}
}
