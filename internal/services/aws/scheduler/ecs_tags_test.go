package scheduler

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awserrors "vorpalstacks/internal/common/errors"
)

func TestParseEcsTags(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		data := []interface{}{
			map[string]interface{}{"key": "k", "value": "v"},
		}
		result, err := parseEcsTags(data)
		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "v", result[0]["value"])
	})

	t.Run("nil input", func(t *testing.T) {
		result, err := parseEcsTags(nil)
		require.NoError(t, err)
		assert.Nil(t, result)
	})

	// A non-map list entry is a wire-format violation, never a silently
	// skipped item the Core validator cannot see.
	t.Run("non-map entry rejected", func(t *testing.T) {
		_, err := parseEcsTags([]interface{}{"not-a-map"})
		require.Error(t, err)
		assertValidationMessage(t, err, "must be a map")
	})

	// A non-string pair value is likewise reported, not dropped.
	t.Run("non-string value rejected", func(t *testing.T) {
		_, err := parseEcsTags([]interface{}{map[string]interface{}{"k": 123}})
		require.Error(t, err)
		assertValidationMessage(t, err, ".k must be a string")
	})
}

// assertValidationMessage pins the ValidationException identity and a
// message fragment in one call.
func assertValidationMessage(t *testing.T, err error, fragment string) {
	t.Helper()
	var awsErr *awserrors.AWSError
	if !assert.ErrorAs(t, err, &awsErr) {
		return
	}
	assert.Equal(t, "ValidationException", awsErr.Code)
	assert.True(t, strings.Contains(awsErr.Message, fragment),
		"message %q does not contain %q", awsErr.Message, fragment)
}
