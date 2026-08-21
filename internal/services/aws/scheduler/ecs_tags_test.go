package scheduler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEcsTags(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		data := []interface{}{
			map[string]interface{}{"key": "k", "value": "v"},
		}
		result := parseEcsTags(data)
		require.Len(t, result, 1)
		assert.Equal(t, "v", result[0]["value"])
	})

	t.Run("nil input", func(t *testing.T) {
		assert.Nil(t, parseEcsTags(nil))
	})
}
