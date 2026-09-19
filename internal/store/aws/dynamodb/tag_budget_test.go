package dynamodb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
)

// The table tag store enforces the documented fifty on the merged set and
// answers an overflow with DynamoDB's coral ValidationException identity,
// injected at construction — no per-path mapping.
func TestTableTagBudgetBoundsMergedSet(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	require.NoError(t, err)
	t.Cleanup(func() { st.Close() })
	store := NewTableStore(st, "123456789012", "us-east-1")

	at := func(prefix string, n int) map[string]string {
		m := make(map[string]string, n)
		for i := 0; i < n; i++ {
			m[fmt.Sprintf("%s-%02d", prefix, i)] = "v"
		}
		return m
	}

	require.NoError(t, store.Tags().Tag("budget-table", at("a", 30)))
	require.NoError(t, store.Tags().Tag("budget-table", at("b", 20)))

	err = store.Tags().Tag("budget-table", map[string]string{"extra": "v"})
	require.Error(t, err)
	apiErr, ok := err.(*awserrors.AWSError)
	require.True(t, ok, "overflow must carry the injected identity, got %T", err)
	require.Equal(t, "com.amazon.coral.validate#ValidationException", apiErr.Code)

	require.NoError(t, store.Tags().Tag("budget-table", map[string]string{"a-00": "rewritten"}))
	tags, err := store.Tags().List("budget-table")
	require.NoError(t, err)
	require.Len(t, tags, 50)
}
