package route53

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	awserrors "vorpalstacks/internal/common/errors"
	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
)

// Route 53's tag budget is its own: the API caps a hosted zone or health
// check at ten tags (not the AWS-wide fifty), and an overflow answers
// InvalidInput — the identity the change-tags core already uses.
func TestRoute53TagBudgetIsTenWithInvalidInput(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	require.NoError(t, err)
	t.Cleanup(func() { st.Close() })
	ts := NewTagStore(st)

	at := func(n int) []types.Tag {
		out := make([]types.Tag, 0, n)
		for i := 0; i < n; i++ {
			out = append(out, types.Tag{Key: fmt.Sprintf("k-%02d", i), Value: "v"})
		}
		return out
	}

	require.NoError(t, ts.Tag("/hostedzone/Z123", at(MaxTagsPerResource)))

	err = ts.Tag("/hostedzone/Z123", []types.Tag{{Key: "over", Value: "v"}})
	require.Error(t, err)
	apiErr, ok := err.(*awserrors.AWSError)
	require.True(t, ok, "overflow must carry the wire identity, got %T", err)
	require.Equal(t, "InvalidInput", apiErr.Code)

	require.NoError(t, ts.Tag("/hostedzone/Z123", []types.Tag{{Key: "k-00", Value: "rewritten"}}))
	tags, err := ts.ListTagsForResource("/hostedzone/Z123")
	require.NoError(t, err)
	require.Len(t, tags, MaxTagsPerResource)
}
