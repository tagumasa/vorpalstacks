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
	cfg := &Config{Port: 8080}

	port := cfg.DefaultPort()
	if port != 8080 {
		t.Errorf("expected 8080, got %d", port)
	}
}
