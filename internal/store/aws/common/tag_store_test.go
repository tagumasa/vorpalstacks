package common

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	awserrors "vorpalstacks/internal/common/errors"
	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
)

func newTestTagStore(t *testing.T) (*TagStore, func()) {
	t.Helper()
	tmpDir := "./tmp/tag-store-test-" + t.Name()
	err := os.MkdirAll(tmpDir, 0o755)
	require.NoError(t, err)

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)

	ts := NewTagStore(s, "test", TagBudget{})
	cleanup := func() {
		s.Close()
		os.RemoveAll(tmpDir)
	}
	return ts, cleanup
}

func TestNewTagStore(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()
	assert.NotNil(t, ts)
	assert.NotNil(t, ts.main)
	assert.NotNil(t, ts.index)
}

func TestNewTagStoreWithRegion(t *testing.T) {
	tmpDir := "./tmp/tag-store-test-region"
	err := os.MkdirAll(tmpDir, 0o755)
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	defer s.Close()

	ts := NewTagStoreWithRegion(s, "svc", "us-east-1", StandardTagBudget("ValidationException"))
	assert.NotNil(t, ts)
}

func TestTagStore_List_NonExistent(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	tags, err := ts.List("nonexistent")
	require.NoError(t, err)
	assert.Empty(t, tags)
}

func TestTagStore_List_NilValue(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	// Store a value that is not a valid tag map (json.Marshal of []byte
	// produces a JSON string, which cannot unmarshal into map[string]string).
	// List should surface this as an error rather than silently returning
	// an empty map and hiding the data corruption.
	err := ts.main.Put("res1", []byte("null"))
	require.NoError(t, err)

	tags, err := ts.List("res1")
	assert.Error(t, err)
	assert.Nil(t, tags)
}

func TestTagStore_List_HappyPath(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.Tag("res1", map[string]string{"env": "prod", "team": "backend"})
	require.NoError(t, err)

	tags, err := ts.List("res1")
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"env": "prod", "team": "backend"}, tags)
}

func TestTagStore_ListAsSlice(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.Tag("res1", map[string]string{"k": "v"})
	require.NoError(t, err)

	slice, err := ts.ListAsSlice("res1")
	require.NoError(t, err)
	require.Len(t, slice, 1)
	assert.Equal(t, "k", slice[0].Key)
	assert.Equal(t, "v", slice[0].Value)
}

func TestTagStore_ListAsSlice_Empty(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	slice, err := ts.ListAsSlice("nonexistent")
	require.NoError(t, err)
	assert.Empty(t, slice)
}

func TestTagStore_Tag_NewResource(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.Tag("res1", map[string]string{"env": "prod"})
	require.NoError(t, err)

	tags, err := ts.List("res1")
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"env": "prod"}, tags)
}

func TestTagStore_Tag_EmptyMap(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.Tag("res1", map[string]string{})
	require.NoError(t, err)
}

func TestTagStore_Tag_OverwriteExisting(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.Tag("res1", map[string]string{"env": "prod", "team": "backend"})
	require.NoError(t, err)

	err = ts.Tag("res1", map[string]string{"env": "staging", "owner": "devops"})
	require.NoError(t, err)

	tags, err := ts.List("res1")
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"env": "staging", "team": "backend", "owner": "devops"}, tags)
}

func TestTagStore_TagFromSlice(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.TagFromSlice("res1", []types.Tag{{Key: "k1", Value: "v1"}, {Key: "k2", Value: "v2"}})
	require.NoError(t, err)

	tags, err := ts.List("res1")
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"k1": "v1", "k2": "v2"}, tags)
}

func TestTagStore_TagFromSlice_Empty(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.TagFromSlice("res1", nil)
	require.NoError(t, err)
}

func TestTagStore_Untag(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.Tag("res1", map[string]string{"env": "prod", "team": "backend"})
	require.NoError(t, err)

	err = ts.Untag("res1", []string{"env"})
	require.NoError(t, err)

	tags, err := ts.List("res1")
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"team": "backend"}, tags)
}

func TestTagStore_Untag_RemoveNonExistentKey(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.Tag("res1", map[string]string{"env": "prod"})
	require.NoError(t, err)

	err = ts.Untag("res1", []string{"nonexistent"})
	require.NoError(t, err)

	tags, err := ts.List("res1")
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"env": "prod"}, tags)
}

func TestTagStore_Untag_EmptyKeys(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.Untag("res1", []string{})
	require.NoError(t, err)
}

func TestTagStore_Untag_RemovesAllTags(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.Tag("res1", map[string]string{"env": "prod"})
	require.NoError(t, err)

	err = ts.Untag("res1", []string{"env"})
	require.NoError(t, err)

	tags, err := ts.List("res1")
	require.NoError(t, err)
	assert.Empty(t, tags)

	exists := ts.main.Exists("res1")
	assert.False(t, exists)
}

func TestTagStore_Delete(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.Tag("res1", map[string]string{"env": "prod", "team": "backend"})
	require.NoError(t, err)

	err = ts.Delete("res1")
	require.NoError(t, err)

	tags, err := ts.List("res1")
	require.NoError(t, err)
	assert.Empty(t, tags)
}

func TestTagStore_Delete_NonExistent(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.Delete("nonexistent")
	require.NoError(t, err)
}

func TestTagStore_Delete_ClearsIndexEntries(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.Tag("res1", map[string]string{"env": "prod"})
	require.NoError(t, err)
	err = ts.Tag("res2", map[string]string{"env": "staging"})
	require.NoError(t, err)

	resources, err := ts.FindByTag("env")
	require.NoError(t, err)
	assert.Len(t, resources, 2)

	err = ts.Delete("res1")
	require.NoError(t, err)

	resources, err = ts.FindByTag("env")
	require.NoError(t, err)
	assert.Len(t, resources, 1)
	assert.Equal(t, "res2", resources[0])
}

func TestTagStore_FindByTag(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.Tag("res1", map[string]string{"env": "prod"})
	require.NoError(t, err)
	err = ts.Tag("res2", map[string]string{"env": "staging", "team": "backend"})
	require.NoError(t, err)

	resources, err := ts.FindByTag("env")
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"res1", "res2"}, resources)

	resources, err = ts.FindByTag("team")
	require.NoError(t, err)
	assert.Equal(t, []string{"res2"}, resources)
}

func TestTagStore_FindByTag_NoMatches(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	resources, err := ts.FindByTag("nonexistent")
	require.NoError(t, err)
	assert.Empty(t, resources)
}

func TestTagStore_FindByTagValue(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.Tag("res1", map[string]string{"env": "prod"})
	require.NoError(t, err)
	err = ts.Tag("res2", map[string]string{"env": "staging"})
	require.NoError(t, err)

	resources, err := ts.FindByTagValue("env", "prod")
	require.NoError(t, err)
	assert.Equal(t, []string{"res1"}, resources)

	resources, err = ts.FindByTagValue("env", "nonexistent")
	require.NoError(t, err)
	assert.Empty(t, resources)
}

func TestTagStore_RebuildIndex(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.Tag("res1", map[string]string{"env": "prod"})
	require.NoError(t, err)
	err = ts.Tag("res2", map[string]string{"env": "staging", "team": "backend"})
	require.NoError(t, err)

	resources, err := ts.FindByTag("env")
	require.NoError(t, err)
	assert.Len(t, resources, 2)

	err = ts.RebuildIndex()
	require.NoError(t, err)

	resources, err = ts.FindByTag("env")
	require.NoError(t, err)
	assert.Len(t, resources, 2)
	assert.ElementsMatch(t, []string{"res1", "res2"}, resources)

	resources, err = ts.FindByTag("team")
	require.NoError(t, err)
	assert.Equal(t, []string{"res2"}, resources)

	tags, err := ts.List("res1")
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"env": "prod"}, tags)
}

func TestTagStore_Tag_OverwriteClearsStaleIndex(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	err := ts.Tag("res1", map[string]string{"env": "prod"})
	require.NoError(t, err)
	err = ts.Tag("res2", map[string]string{"env": "staging"})
	require.NoError(t, err)

	resources, err := ts.FindByTagValue("env", "staging")
	require.NoError(t, err)
	assert.Equal(t, []string{"res2"}, resources)

	err = ts.Tag("res2", map[string]string{"env": "prod"})
	require.NoError(t, err)

	resources, err = ts.FindByTagValue("env", "staging")
	require.NoError(t, err)
	assert.Empty(t, resources)

	resources, err = ts.FindByTagValue("env", "prod")
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"res1", "res2"}, resources)
}

// The merge seam enforces the per-resource ceiling against the cumulative
// total, not the incoming batch alone: a resource already holding 30 tags
// admits a 20-tag write, the write that would carry it past 50 rejects
// with the construction budget's identity and leaves the stored set
// untouched, and the full-swap form carries the same ceiling.
func TestTagStore_Tag_CumulativeCeiling(t *testing.T) {
	ts := newBudgetTagStore(t, StandardTagBudget("TooManyTagsException"))

	batch := func(n, offset int) map[string]string {
		tags := make(map[string]string, n)
		for i := 0; i < n; i++ {
			tags[fmt.Sprintf("key-%d", offset+i)] = "v"
		}
		return tags
	}

	require.NoError(t, ts.Tag("res1", batch(30, 0)))
	require.NoError(t, ts.Tag("res1", batch(20, 30)))
	tags, err := ts.List("res1")
	require.NoError(t, err)
	assert.Len(t, tags, types.MaxTagsPerResource)

	err = ts.Tag("res1", map[string]string{"over": "v"})
	apiErr, ok := err.(*awserrors.AWSError)
	require.True(t, ok, "budget rejection must carry the injected identity, got %T", err)
	assert.Equal(t, "TooManyTagsException", apiErr.Code)
	// The rejected write leaves the stored set untouched.
	tags, err = ts.List("res1")
	require.NoError(t, err)
	assert.Len(t, tags, types.MaxTagsPerResource)
	assert.NotContains(t, tags, "over")

	// The full-swap form carries the same ceiling.
	err = ts.Replace("res1", batch(types.MaxTagsPerResource+1, 0))
	require.Error(t, err)
	err = ts.Replace("res1", batch(types.MaxTagsPerResource, 0))
	require.NoError(t, err)
}

// ListAsSlice answers in ascending key order: every marker-paginated
// consumer rides a stable order (pagination presupposes one), and the
// map-backed store must not leak Go's random iteration order to the
// listing surfaces.
func TestListAsSliceIsOrderedAscendingByKey(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()
	if err := ts.Tag("res-1", map[string]string{
		"zebra": "1", "alpha": "2", "mike": "3", "beta": "4",
	}); err != nil {
		t.Fatalf("tag: %v", err)
	}
	for i := 0; i < 20; i++ {
		tags, err := ts.ListAsSlice("res-1")
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(tags) != 4 {
			t.Fatalf("want 4 tags, got %d", len(tags))
		}
		for j := 1; j < len(tags); j++ {
			if tags[j-1].Key >= tags[j].Key {
				t.Fatalf("iteration %d: keys not ascending: %v", i, tags)
			}
		}
	}
}
