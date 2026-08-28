package neptune

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The documented Describe* page-size window is 20-100, so the paginator's
// page sizes in these tests stay inside or beyond that window; sub-20
// requests are covered by the explicit floor cases below.
func buildPaginateItems(n int) ([]interface{}, func(interface{}) string) {
	items := make([]interface{}, n)
	for i := 0; i < n; i++ {
		items[i] = map[string]string{"id": fmt.Sprintf("id%02d", i)}
	}
	keyFn := func(item interface{}) string {
		return item.(map[string]string)["id"]
	}
	return items, keyFn
}

func TestPaginateItems(t *testing.T) {
	items, keyFn := buildPaginateItems(30)

	t.Run("no marker returns first page", func(t *testing.T) {
		result, marker, truncated := paginateItems(items, "", 25, keyFn)
		assert.Equal(t, 25, len(result))
		assert.Equal(t, "id24", marker)
		assert.True(t, truncated)
	})

	t.Run("valid marker returns subsequent page", func(t *testing.T) {
		result, marker, truncated := paginateItems(items, "id04", 25, keyFn)
		assert.Equal(t, 25, len(result))
		assert.Equal(t, "", marker)
		assert.False(t, truncated)
	})

	t.Run("partial final page not truncated", func(t *testing.T) {
		result, marker, truncated := paginateItems(items, "id24", 25, keyFn)
		assert.Equal(t, 5, len(result))
		assert.Equal(t, "", marker)
		assert.False(t, truncated)
	})

	// A pagination marker that does not match any key must yield an empty
	// page rather than iterating from the start of the slice. This
	// preserves cursor semantics across paginated list APIs.
	t.Run("invalid marker returns empty", func(t *testing.T) {
		result, marker, truncated := paginateItems(items, "nonexistent", 25, keyFn)
		assert.Equal(t, 0, len(result))
		assert.Equal(t, "", marker)
		assert.False(t, truncated)
	})

	t.Run("empty items returns empty", func(t *testing.T) {
		result, marker, truncated := paginateItems([]interface{}{}, "", 25, keyFn)
		assert.Equal(t, 0, len(result))
		assert.Equal(t, "", marker)
		assert.False(t, truncated)
	})

	t.Run("marker at last item returns empty", func(t *testing.T) {
		result, marker, truncated := paginateItems(items, "id29", 25, keyFn)
		assert.Equal(t, 0, len(result))
		assert.Equal(t, "", marker)
		assert.False(t, truncated)
	})

	t.Run("maxRecords zero defaults to 100", func(t *testing.T) {
		result, _, truncated := paginateItems(items, "", 0, keyFn)
		assert.Equal(t, 30, len(result))
		assert.False(t, truncated)
	})

	t.Run("maxRecords below minimum clamps to 20", func(t *testing.T) {
		result, marker, truncated := paginateItems(items, "", 5, keyFn)
		assert.Equal(t, 20, len(result))
		assert.Equal(t, "id19", marker)
		assert.True(t, truncated)
	})

	t.Run("maxRecords capped at 100", func(t *testing.T) {
		big, bigKeyFn := buildPaginateItems(150)
		result, _, truncated := paginateItems(big, "", 200, bigKeyFn)
		assert.Equal(t, 100, len(result))
		assert.True(t, truncated)
	})

	t.Run("single item exact fit", func(t *testing.T) {
		single := []interface{}{map[string]string{"id": "only"}}
		result, marker, truncated := paginateItems(single, "", 20, keyFn)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "", marker)
		assert.False(t, truncated)
	})
}
