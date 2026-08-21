package sesv2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMessageTags(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		params := map[string]interface{}{
			"MessageTags": []interface{}{
				map[string]interface{}{"Name": "n1", "Value": "v1"},
			},
		}
		result := ParseMessageTags(params, "MessageTags")
		require.Len(t, result, 1)
		assert.Equal(t, "n1", result[0].Name)
	})

	t.Run("missing", func(t *testing.T) {
		assert.Nil(t, ParseMessageTags(map[string]interface{}{}, "MessageTags"))
	})

	t.Run("wrong type", func(t *testing.T) {
		params := map[string]interface{}{"MessageTags": "not a list"}
		assert.Nil(t, ParseMessageTags(params, "MessageTags"))
	})
}

func TestParseMessageTagsFromList(t *testing.T) {
	t.Run("mixed types", func(t *testing.T) {
		list := []interface{}{
			map[string]interface{}{"Name": "n1", "Value": "v1"},
			"invalid",
		}
		result := ParseMessageTagsFromList(list)
		require.Len(t, result, 1)
	})
}
