package grpcweb

import "testing"

func TestConfig_DefaultPort(t *testing.T) {
	cfg := &Config{}

	port := cfg.DefaultPort()
	if port != 50090 {
		t.Errorf("expected 50090, got %d", port)
	}
}

func TestConfig_DefaultPort_WithPort(t *testing.T) {
	cfg := &Config{Port: 50090}

	port := cfg.DefaultPort()
	if port != 50090 {
		t.Errorf("expected 50090, got %d", port)
	}
}
