package logs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLogConstants(t *testing.T) {
	assert.Equal(t, int64(104857600), int64(DefaultMaxLogSize))
	assert.Equal(t, int64(1000000), int64(DefaultMaxLogEntries))
	assert.Equal(t, 7*24*time.Hour, DefaultLogMaxAge)
	assert.Equal(t, time.Hour, DefaultRotationInterval)
}
