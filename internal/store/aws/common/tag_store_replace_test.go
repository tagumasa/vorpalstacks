package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Replace swaps the full tag set in one call: removed keys disappear, changed
// keys take the new value, and the inverted index follows the main entry.
func TestTagStoreReplaceSwapsFullTagSet(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	require.NoError(t, ts.Tag("res-1", map[string]string{
		"env":    "prod",
		"team":   "core",
		"legacy": "gone",
	}))

	require.NoError(t, ts.Replace("res-1", map[string]string{
		"env":  "dev",
		"team": "core",
	}))

	tags, err := ts.List("res-1")
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"env": "dev", "team": "core"}, tags)

	// The stale value entry for the changed key must be gone and the removed
	// key must no longer resolve to the resource.
	prodResources, err := ts.FindByTagValue("env", "prod")
	require.NoError(t, err)
	assert.Empty(t, prodResources)
	legacyResources, err := ts.FindByTag("legacy")
	require.NoError(t, err)
	assert.Empty(t, legacyResources)
	devResources, err := ts.FindByTagValue("env", "dev")
	require.NoError(t, err)
	assert.Equal(t, []string{"res-1"}, devResources)
}

// Replacing with an empty set clears the resource's tags entirely, matching
// Untag's behaviour when the last key is removed.
func TestTagStoreReplaceWithEmptySetClearsTags(t *testing.T) {
	ts, cleanup := newTestTagStore(t)
	defer cleanup()

	require.NoError(t, ts.Tag("res-2", map[string]string{"env": "prod"}))
	require.NoError(t, ts.Replace("res-2", map[string]string{}))

	tags, err := ts.List("res-2")
	require.NoError(t, err)
	assert.Empty(t, tags)
	resources, err := ts.FindByTag("env")
	require.NoError(t, err)
	assert.Empty(t, resources)
}
