package neptune

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginateItems(t *testing.T) {
	items := []interface{}{
		map[string]string{"id": "a"},
		map[string]string{"id": "b"},
		map[string]string{"id": "c"},
		map[string]string{"id": "d"},
		map[string]string{"id": "e"},
	}
	keyFn := func(item interface{}) string {
		return item.(map[string]string)["id"]
	}

	t.Run("no marker returns first page", func(t *testing.T) {
		result, marker, truncated := paginateItems(items, "", 3, keyFn)
		assert.Equal(t, 3, len(result))
		assert.Equal(t, "c", marker)
		assert.True(t, truncated)
	})

	t.Run("valid marker returns subsequent page", func(t *testing.T) {
		result, marker, truncated := paginateItems(items, "b", 2, keyFn)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, "d", marker)
		assert.True(t, truncated)
	})

	t.Run("last page not truncated", func(t *testing.T) {
		result, marker, truncated := paginateItems(items, "d", 2, keyFn)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "", marker)
		assert.False(t, truncated)
	})

	t.Run("invalid marker returns empty (Bug M2 regression)", func(t *testing.T) {
		result, marker, truncated := paginateItems(items, "nonexistent", 3, keyFn)
		assert.Equal(t, 0, len(result))
		assert.Equal(t, "", marker)
		assert.False(t, truncated)
	})

	t.Run("empty items returns empty", func(t *testing.T) {
		result, marker, truncated := paginateItems([]interface{}{}, "", 3, keyFn)
		assert.Equal(t, 0, len(result))
		assert.Equal(t, "", marker)
		assert.False(t, truncated)
	})

	t.Run("marker at last item returns empty", func(t *testing.T) {
		result, marker, truncated := paginateItems(items, "e", 3, keyFn)
		assert.Equal(t, 0, len(result))
		assert.Equal(t, "", marker)
		assert.False(t, truncated)
	})

	t.Run("maxRecords zero defaults to 100", func(t *testing.T) {
		result, _, truncated := paginateItems(items, "", 0, keyFn)
		assert.Equal(t, 5, len(result))
		assert.False(t, truncated)
	})

	t.Run("maxRecords exceeds total", func(t *testing.T) {
		result, marker, truncated := paginateItems(items, "", 100, keyFn)
		assert.Equal(t, 5, len(result))
		assert.Equal(t, "", marker)
		assert.False(t, truncated)
	})

	t.Run("maxRecords capped at 100", func(t *testing.T) {
		big := make([]interface{}, 150)
		for i := range big {
			big[i] = map[string]string{"id": string(rune('a'+i%26)) + string(rune(i))}
		}
		bigKeyFn := func(item interface{}) string {
			return item.(map[string]string)["id"]
		}
		result, _, truncated := paginateItems(big, "", 200, bigKeyFn)
		assert.Equal(t, 100, len(result))
		assert.True(t, truncated)
	})

	t.Run("single item exact fit", func(t *testing.T) {
		single := []interface{}{map[string]string{"id": "only"}}
		result, marker, truncated := paginateItems(single, "", 1, keyFn)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "", marker)
		assert.False(t, truncated)
	})
}
