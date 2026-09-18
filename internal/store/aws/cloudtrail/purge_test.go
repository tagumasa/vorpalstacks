package cloudtrail

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"vorpalstacks/internal/core/storage"
)

func newPurgeTestStore(t *testing.T) (*CloudTrailStore, *storage.PebbleStorage) {
	t.Helper()
	tmpDir := "./tmp/cloudtrail-purge-test"
	os.RemoveAll(tmpDir)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	t.Cleanup(func() { s.Close() })

	return NewCloudTrailStore(s, "acc123", "us-east-1"), s
}

func purgeTestEvent(id, name, user string, at time.Time) *Event {
	return &Event{
		EventID:     id,
		EventName:   name,
		EventSource: "cloudtrail.amazonaws.com",
		EventTime:   at,
		UserIdentity: &UserIdentity{
			Type:     "IAMUser",
			UserName: user,
		},
	}
}

// TestPurgeEventHistoryBefore verifies the retention purge end to end:
// entries below the cutoff disappear from every bucket together with their
// five companion entries, retained entries stay visible through every
// lookup path, and a re-run removes nothing.
func TestPurgeEventHistoryBefore(t *testing.T) {
	store, s := newPurgeTestStore(t)

	oldHour := time.Now().UTC().Add(-48 * time.Hour).Truncate(time.Hour)
	freshTime := time.Now().UTC().Add(-1 * time.Hour).Truncate(time.Hour)

	oldEvents := []*Event{
		purgeTestEvent("old-001", "CreateTrail", "alice", oldHour),
		purgeTestEvent("old-002", "DeleteTrail", "bob", oldHour.Add(11*time.Minute)),
		purgeTestEvent("old-003", "StartLogging", "carol", oldHour.Add(22*time.Minute)),
	}
	freshEvents := []*Event{
		purgeTestEvent("fresh-001", "CreateTrail", "alice", freshTime),
		purgeTestEvent("fresh-002", "StopLogging", "dave", freshTime.Add(5*time.Minute)),
	}
	for _, e := range append(append([]*Event{}, oldEvents...), freshEvents...) {
		require.NoError(t, store.PutEvent(e))
	}

	cutoff := time.Now().UTC().Add(-24 * time.Hour)
	result, err := store.PurgeEventHistoryBefore(cutoff, 0)
	require.NoError(t, err)

	assert.Equal(t, 3, result.Events, "three old records deleted")
	assert.Equal(t, 3, result.EventIDIndex)
	assert.Equal(t, 3, result.TimeIndex)
	assert.Equal(t, 3, result.EventNameIndex)
	assert.Equal(t, 3, result.UsernameIndex)
	assert.Equal(t, 3, result.EventSourceIndex)
	assert.Equal(t, 1, result.Batches)

	// Bucket cardinalities balance exactly: only the retained events
	// remain, and every bucket family agrees on the count.
	assert.Equal(t, 2, store.eventsStore.Count(), "events bucket")
	assert.Equal(t, 2, store.eventIDIndexStore.Count(), "eventID index bucket")

	oldHourStr := oldHour.Format("2006-01-02:15")
	assert.Equal(t, 0, s.Bucket("ct_idx_time:acc123:us-east-1:"+oldHourStr).Count(),
		"purged hour bucket empty")
	assert.Equal(t, 2, s.Bucket("ct_idx_time:acc123:us-east-1:"+freshTime.Format("2006-01-02:15")).Count(),
		"retained hour bucket intact")

	// Retained events stay visible through the identifier path; purged
	// identifiers are gone.
	fresh, err := store.GetEventByID("fresh-001")
	require.NoError(t, err)
	assert.Equal(t, "CreateTrail", fresh.EventName)
	_, err = store.GetEventByID("old-001")
	assert.ErrorIs(t, err, ErrEventNotFound)

	// The scan path returns only retained events.
	events, nextToken, err := store.LookupEvents(EventQuery{MaxResults: 50})
	require.NoError(t, err)
	assert.Empty(t, nextToken)
	assert.Len(t, events, 2)

	// The indexer paths no longer resolve purged events.
	ids, _, err := store.indexer.QueryByEventName([]string{"DeleteTrail"}, 10, IndexCursor{})
	require.NoError(t, err)
	assert.Empty(t, ids, "purged event name index entry removed")
	ids, _, err = store.indexer.QueryByUsername("bob", 10, IndexCursor{})
	require.NoError(t, err)
	assert.Empty(t, ids, "purged username index entry removed")
	oldHourEnd := oldHour.Add(time.Hour)
	ids, _, err = store.indexer.QueryByTime(&oldHour, &oldHourEnd, 10, IndexCursor{})
	require.NoError(t, err)
	assert.Empty(t, ids, "purged hour bucket removed")

	// A second run over the same cutoff is a no-op.
	result, err = store.PurgeEventHistoryBefore(cutoff, 0)
	require.NoError(t, err)
	assert.Equal(t, &PurgeResult{}, result, "second run removes nothing")
	assert.Equal(t, 2, store.eventsStore.Count())
}

// TestPurgeEventHistoryBeforeBatchResume drives the batch loop across
// multiple transactions: every batch removes its own slice and the next
// batch resumes strictly after the last deleted key.
func TestPurgeEventHistoryBeforeBatchResume(t *testing.T) {
	store, _ := newPurgeTestStore(t)

	base := time.Now().UTC().Add(-48 * time.Hour).Truncate(time.Hour)
	for i := 0; i < 5; i++ {
		e := purgeTestEvent(fmt.Sprintf("batch-%03d", i), "CreateTrail", "alice", base.Add(time.Duration(i)*time.Second))
		require.NoError(t, store.PutEvent(e))
	}
	require.NoError(t, store.PutEvent(purgeTestEvent("batch-fresh", "CreateTrail", "alice", time.Now().UTC().Add(-time.Minute))))

	result, err := store.PurgeEventHistoryBefore(time.Now().UTC().Add(-24*time.Hour), 2)
	require.NoError(t, err)

	assert.Equal(t, 5, result.Events)
	assert.Equal(t, 3, result.Batches, "batches of 2, 2 and 1")
	assert.Equal(t, 1, store.eventsStore.Count())
	remaining, _, err := store.LookupEvents(EventQuery{MaxResults: 10})
	require.NoError(t, err)
	require.Len(t, remaining, 1)
	assert.Equal(t, "batch-fresh", remaining[0].EventID)
}

// TestPurgeEventHistoryBeforeNanosecondPrecision pins the key-vs-record
// precision contract: recorded event keys carry full nanosecond timestamps
// while the record proto stores milliseconds, so the purge must rebuild
// index keys from the storage key (not the record) or the index entries
// would survive the purge.
func TestPurgeEventHistoryBeforeNanosecondPrecision(t *testing.T) {
	store, s := newPurgeTestStore(t)

	// Sub-millisecond component guarantees the record's millisecond
	// rounding differs from the key's nanosecond value.
	at := time.Now().UTC().Add(-48 * time.Hour).Truncate(time.Hour).Add(123456789 * time.Nanosecond)
	require.NoError(t, store.PutEvent(purgeTestEvent("nano-001", "CreateTrail", "alice", at)))

	_, err := store.PurgeEventHistoryBefore(time.Now().UTC().Add(-24*time.Hour), 0)
	require.NoError(t, err)

	assert.Equal(t, 0, store.eventsStore.Count())
	assert.Equal(t, 0, store.eventIDIndexStore.Count())
	assert.Equal(t, 0, s.Bucket("ct_idx_time:acc123:us-east-1:"+at.Format("2006-01-02:15")).Count())
	assert.Equal(t, 0, s.Bucket("ct_idx_event:acc123:us-east-1:CreateTrail").Count())
	assert.Equal(t, 0, s.Bucket("ct_idx_user:acc123:us-east-1:alice").Count())
	assert.Equal(t, 0, s.Bucket("ct_idx_source:acc123:us-east-1:cloudtrail.amazonaws.com").Count())
}
