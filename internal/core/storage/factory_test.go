package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenPebbleDirectly(t *testing.T) {
	storage, err := Open(t.TempDir())
	require.NoError(t, err)
	defer storage.Close()

	assert.NotNil(t, storage)
}

func TestConfigOptions(t *testing.T) {
	cfg := &Config{}
	opt := WithCacheSize(512 << 20)
	opt(cfg)
	assert.Equal(t, int64(512<<20), cfg.CacheSizeBytes)

	opt = WithMaxOpenFiles(500)
	opt(cfg)
	assert.Equal(t, 500, cfg.MaxOpenFiles)

	opt = WithBytesPerSync(2 << 20)
	opt(cfg)
	assert.Equal(t, 2<<20, cfg.BytesPerSync)
}
