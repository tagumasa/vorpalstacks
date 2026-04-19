package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/utils/aws/types"
)

func newTestTagStore(t *testing.T) (*TagStore, func()) {
	t.Helper()
	tmpDir := "./tmp/tag-store-test-" + t.Name()
	err := os.MkdirAll(tmpDir, 0o755)
	require.NoError(t, err)

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)

	ts := NewTagStore(s, "test")
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

	ts := NewTagStoreWithRegion(s, "svc", "us-east-1")
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

	err := ts.main.Put("res1", []byte("null"))
	require.NoError(t, err)

	tags, err := ts.List("res1")
	require.NoError(t, err)
	assert.Empty(t, tags)
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
