package grpcweb

import "testing"

func TestConfig_DefaultPort(t *testing.T) {
	cfg := &Config{}

	port := cfg.DefaultPort()
	if port != "9090" {
		t.Errorf("expected '9090', got '%s'", port)
	}
}

func TestConfig_DefaultPort_WithPort(t *testing.T) {
	cfg := &Config{Port: "8080"}

	port := cfg.DefaultPort()
	if port != "8080" {
		t.Errorf("expected '8080', got '%s'", port)
	}
}
