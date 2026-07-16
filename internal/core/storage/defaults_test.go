package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLockMaxRetries(t *testing.T) {
	assert.Equal(t, 100, DefaultLockMaxRetries)
}

func TestDefaultVersionMaxRetries(t *testing.T) {
	assert.Equal(t, 100, DefaultVersionMaxRetries)
}

func TestDefaultLockRetryDelay(t *testing.T) {
	assert.Equal(t, 50*time.Millisecond, DefaultLockRetryDelay)
}

func TestPrefixConstants(t *testing.T) {
	assert.Equal(t, "__bucket__:", BucketPrefix)
	assert.Equal(t, "__locks__:", LockPrefix)
	assert.Equal(t, "__indexes__:", IndexPrefix)
	assert.Equal(t, "__versioned__:", VersionedPrefix)
	assert.Equal(t, "__cond__:", ConditionalPrefix)
	assert.Equal(t, "__idempotency__:", IdempotencyPrefix)
	assert.Equal(t, "__meta_", MetaPrefix)
	assert.Equal(t, "__counter__", CounterPrefix)
}
