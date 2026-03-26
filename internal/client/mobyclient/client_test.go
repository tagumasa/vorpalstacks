package mobyclient

import (
	"testing"
	"vorpalstacks/internal/core/logs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	logger := logs.NewLogger(&logs.Config{})

	cfg := Config{
		Host:    "unix:///var/run/docker.sock",
		Version: "1.41",
	}

	client, err := NewClient(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, "unix:///var/run/docker.sock", cfg.Host)
	assert.Equal(t, "1.41", cfg.Version)
	assert.Equal(t, false, cfg.TLSVerify)
	assert.Equal(t, 30, cfg.Timeout)
}
