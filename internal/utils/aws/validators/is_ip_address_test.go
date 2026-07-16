package validators

import "testing"

func TestIsIPAddress(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"Valid IPv4", "192.168.1.1", true},
		{"Valid IPv4 localhost", "127.0.0.1", true},
		{"Valid IPv6", "::1", true},
		{"Valid IPv6 full", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", true},
		{"Invalid empty", "", false},
		{"Invalid random", "random", false},
		{"Invalid partial", "192.168.1", false},
		{"Invalid too many octets", "192.168.1.1.1", false},
		{"Invalid out of range", "256.256.256.256", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIPAddress(tt.input); got != tt.valid {
				t.Errorf("IsIPAddress(%q) = %v, want %v", tt.input, got, tt.valid)
			}
		})
	}
}
