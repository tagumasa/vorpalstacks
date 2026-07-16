package cloudtrail

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"vorpalstacks/internal/core/storage"
)

func TestEventIndexKey_EncodePrefix(t *testing.T) {
	t.Run("Time index with all segments", func(t *testing.T) {
		key := NewTimeIndexKey("acc123", "us-east-1", "2024-02-25:10", "evt456")
		prefix := key.EncodePrefix()
		assert.Equal(t, "ct_idx_time:acc123:us-east-1:2024-02-25:10", prefix)
	})

	t.Run("Event name index with segment2", func(t *testing.T) {
		key := NewEventNameIndexKey("acc123", "us-east-1", "CreateTrail", 1708867200, "evt456")
		prefix := key.EncodePrefix()
		assert.Equal(t, "ct_idx_event:acc123:us-east-1:CreateTrail", prefix)
	})

	t.Run("Username index", func(t *testing.T) {
		key := NewUsernameIndexKey("acc123", "ap-northeast-1", "testuser", 1708867200, "evt456")
		prefix := key.EncodePrefix()
		assert.Equal(t, "ct_idx_user:acc123:ap-northeast-1:testuser", prefix)
	})

	t.Run("Event source index", func(t *testing.T) {
		key := NewEventSourceIndexKey("acc123", "eu-west-1", "cloudtrail.amazonaws.com", 1708867200, "evt456")
		prefix := key.EncodePrefix()
		assert.Equal(t, "ct_idx_source:acc123:eu-west-1:cloudtrail.amazonaws.com", prefix)
	})

	t.Run("Prefix only without segment1", func(t *testing.T) {
		key := &EventIndexKey{
			IndexType: IndexByTime,
			AccountID: "acc123",
			Region:    "us-east-1",
		}
		prefix := key.EncodePrefix()
		assert.Equal(t, "ct_idx_time:acc123:us-east-1", prefix)
	})
}

func TestEventIndexManager_AddAndRemoveIndex(t *testing.T) {
	tmpDir := "./tmp/cloudtrail-index-test"
	defer os.RemoveAll(tmpDir)

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	defer s.Close()

	manager := NewEventIndexManager(s, "acc123", "us-east-1")

	event := &Event{
		EventID:     "evt-001",
		EventName:   "CreateTrail",
		EventSource: "cloudtrail.amazonaws.com",
		EventTime:   time.Date(2024, 2, 25, 10, 30, 0, 0, time.UTC),
		UserIdentity: &UserIdentity{
			Type:     "IAMUser",
			UserName: "testuser",
		},
	}

	t.Run("Add index creates all index entries", func(t *testing.T) {
		err := manager.AddIndex(event)
		require.NoError(t, err)

		ids, err := manager.QueryByEventName([]string{"CreateTrail"}, 10)
		require.NoError(t, err)
		assert.Contains(t, ids, "evt-001")

		ids, err = manager.QueryByUsername("testuser", 10)
		require.NoError(t, err)
		assert.Contains(t, ids, "evt-001")

		ids, err = manager.QueryByEventSource("cloudtrail.amazonaws.com", 10)
		require.NoError(t, err)
		assert.Contains(t, ids, "evt-001")
	})

	t.Run("Remove index deletes all index entries", func(t *testing.T) {
		err := manager.RemoveIndex(event)
		require.NoError(t, err)

		ids, err := manager.QueryByEventName([]string{"CreateTrail"}, 10)
		require.NoError(t, err)
		assert.NotContains(t, ids, "evt-001")
	})
}

func TestEventIndexManager_QueryByTime(t *testing.T) {
	tmpDir := "./tmp/cloudtrail-time-query-test"
	defer os.RemoveAll(tmpDir)

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	defer s.Close()

	manager := NewEventIndexManager(s, "acc123", "us-east-1")

	baseTime := time.Date(2024, 2, 25, 10, 0, 0, 0, time.UTC)

	events := []*Event{
		{EventID: "evt-001", EventName: "CreateTrail", EventSource: "cloudtrail.amazonaws.com", EventTime: baseTime},
		{EventID: "evt-002", EventName: "DeleteTrail", EventSource: "cloudtrail.amazonaws.com", EventTime: baseTime.Add(30 * time.Minute)},
		{EventID: "evt-003", EventName: "StartLogging", EventSource: "cloudtrail.amazonaws.com", EventTime: baseTime.Add(2 * time.Hour)},
	}

	for _, e := range events {
		err := manager.AddIndex(e)
		require.NoError(t, err)
	}

	t.Run("Query with time range", func(t *testing.T) {
		start := baseTime
		end := baseTime.Add(time.Hour)
		ids, err := manager.QueryByTime(&start, &end, 10)
		require.NoError(t, err)
		assert.Len(t, ids, 2)
		assert.Contains(t, ids, "evt-001")
		assert.Contains(t, ids, "evt-002")
	})

	t.Run("Query with maxResults limit", func(t *testing.T) {
		start := baseTime
		end := baseTime.Add(3 * time.Hour)
		ids, err := manager.QueryByTime(&start, &end, 2)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(ids), 2)
	})

	t.Run("Query with startTime only", func(t *testing.T) {
		ids, err := manager.QueryByTime(&baseTime, nil, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(ids), 1)
	})

	t.Run("Query with endTime only", func(t *testing.T) {
		ids, err := manager.QueryByTime(nil, &baseTime, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(ids), 1)
	})
}

func TestEventIndexManager_QueryByEventName(t *testing.T) {
	tmpDir := "./tmp/cloudtrail-eventname-query-test"
	defer os.RemoveAll(tmpDir)

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	defer s.Close()

	manager := NewEventIndexManager(s, "acc123", "us-east-1")

	now := time.Now()

	events := []*Event{
		{EventID: "evt-001", EventName: "CreateTrail", EventSource: "cloudtrail.amazonaws.com", EventTime: now},
		{EventID: "evt-002", EventName: "CreateTrail", EventSource: "cloudtrail.amazonaws.com", EventTime: now.Add(time.Second)},
		{EventID: "evt-003", EventName: "DeleteTrail", EventSource: "cloudtrail.amazonaws.com", EventTime: now.Add(2 * time.Second)},
	}

	for _, e := range events {
		err := manager.AddIndex(e)
		require.NoError(t, err)
	}

	t.Run("Query single event name", func(t *testing.T) {
		ids, err := manager.QueryByEventName([]string{"CreateTrail"}, 10)
		require.NoError(t, err)
		assert.Len(t, ids, 2)
	})

	t.Run("Query multiple event names", func(t *testing.T) {
		ids, err := manager.QueryByEventName([]string{"CreateTrail", "DeleteTrail"}, 10)
		require.NoError(t, err)
		assert.Len(t, ids, 3)
	})

	t.Run("Query with maxResults", func(t *testing.T) {
		ids, err := manager.QueryByEventName([]string{"CreateTrail"}, 1)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(ids), 1)
	})
}

func TestEventIndexManager_QueryByUsername(t *testing.T) {
	tmpDir := "./tmp/cloudtrail-username-query-test"
	defer os.RemoveAll(tmpDir)

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	defer s.Close()

	manager := NewEventIndexManager(s, "acc123", "us-east-1")

	now := time.Now()

	events := []*Event{
		{EventID: "evt-001", EventName: "CreateTrail", EventSource: "cloudtrail.amazonaws.com", EventTime: now, UserIdentity: &UserIdentity{Type: "IAMUser", UserName: "alice"}},
		{EventID: "evt-002", EventName: "DeleteTrail", EventSource: "cloudtrail.amazonaws.com", EventTime: now.Add(time.Second), UserIdentity: &UserIdentity{Type: "IAMUser", UserName: "bob"}},
		{EventID: "evt-003", EventName: "StartLogging", EventSource: "cloudtrail.amazonaws.com", EventTime: now.Add(2 * time.Second), UserIdentity: &UserIdentity{Type: "IAMUser", UserName: "alice"}},
	}

	for _, e := range events {
		err := manager.AddIndex(e)
		require.NoError(t, err)
	}

	t.Run("Query by username returns correct events", func(t *testing.T) {
		ids, err := manager.QueryByUsername("alice", 10)
		require.NoError(t, err)
		assert.Len(t, ids, 2)
	})

	t.Run("Query non-existent username returns empty", func(t *testing.T) {
		ids, err := manager.QueryByUsername("charlie", 10)
		require.NoError(t, err)
		assert.Empty(t, ids)
	})
}

func TestEventIndexManager_QueryByEventSource(t *testing.T) {
	tmpDir := "./tmp/cloudtrail-source-query-test"
	defer os.RemoveAll(tmpDir)

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	defer s.Close()

	manager := NewEventIndexManager(s, "acc123", "us-east-1")

	now := time.Now()

	events := []*Event{
		{EventID: "evt-001", EventName: "CreateBucket", EventSource: "s3.amazonaws.com", EventTime: now},
		{EventID: "evt-002", EventName: "CreateTrail", EventSource: "cloudtrail.amazonaws.com", EventTime: now.Add(time.Second)},
		{EventID: "evt-003", EventName: "DeleteBucket", EventSource: "s3.amazonaws.com", EventTime: now.Add(2 * time.Second)},
	}

	for _, e := range events {
		err := manager.AddIndex(e)
		require.NoError(t, err)
	}

	t.Run("Query by event source returns correct events", func(t *testing.T) {
		ids, err := manager.QueryByEventSource("s3.amazonaws.com", 10)
		require.NoError(t, err)
		assert.Len(t, ids, 2)
	})

	t.Run("Query non-existent source returns empty", func(t *testing.T) {
		ids, err := manager.QueryByEventSource("lambda.amazonaws.com", 10)
		require.NoError(t, err)
		assert.Empty(t, ids)
	})
}

func TestEventIndexManager_ClearIndexes(t *testing.T) {
	tmpDir := "./tmp/cloudtrail-clear-test"
	defer os.RemoveAll(tmpDir)

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	defer s.Close()

	manager := NewEventIndexManager(s, "acc123", "us-east-1")

	now := time.Now()

	events := []*Event{
		{EventID: "evt-001", EventName: "CreateTrail", EventSource: "cloudtrail.amazonaws.com", EventTime: now, UserIdentity: &UserIdentity{Type: "IAMUser", UserName: "testuser"}},
		{EventID: "evt-002", EventName: "DeleteTrail", EventSource: "cloudtrail.amazonaws.com", EventTime: now.Add(time.Second), UserIdentity: &UserIdentity{Type: "IAMUser", UserName: "testuser"}},
	}

	for _, e := range events {
		err := manager.AddIndex(e)
		require.NoError(t, err)
	}

	t.Run("ClearIndexes removes all indexes", func(t *testing.T) {
		err := manager.ClearIndexes("acc123", "us-east-1")
		require.NoError(t, err)

		ids, err := manager.QueryByEventName([]string{"CreateTrail"}, 10)
		require.NoError(t, err)
		assert.Empty(t, ids)

		ids, err = manager.QueryByUsername("testuser", 10)
		require.NoError(t, err)
		assert.Empty(t, ids)
	})
}

func TestEventIndexManager_EventWithoutOptionalFields(t *testing.T) {
	tmpDir := "./tmp/cloudtrail-minimal-event-test"
	defer os.RemoveAll(tmpDir)

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	defer s.Close()

	manager := NewEventIndexManager(s, "acc123", "us-east-1")

	t.Run("Event with minimal fields only creates time index", func(t *testing.T) {
		event := &Event{
			EventID:     "evt-minimal",
			EventName:   "",
			EventSource: "",
			EventTime:   time.Now(),
			UserIdentity: &UserIdentity{
				Type: "IAMUser",
			},
		}

		err := manager.AddIndex(event)
		require.NoError(t, err)

		start := time.Now().Add(-time.Hour)
		end := time.Now().Add(time.Hour)
		ids, err := manager.QueryByTime(&start, &end, 10)
		require.NoError(t, err)
		assert.Contains(t, ids, "evt-minimal")
	})
}

func TestEventIndexManager_AddIndexInTxn(t *testing.T) {
	tmpDir := "./tmp/cloudtrail-txn-index-test"
	defer os.RemoveAll(tmpDir)

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	defer s.Close()

	tstore, ok := s.(storage.TransactionalStorageWith2PC)
	require.True(t, ok, "storage must support transactions")

	manager := NewEventIndexManager(s, "acc123", "us-east-1")

	event := &Event{
		EventID:     "evt-txn-001",
		EventName:   "CreateTrail",
		EventSource: "cloudtrail.amazonaws.com",
		EventTime:   time.Date(2024, 2, 25, 10, 30, 0, 0, time.UTC),
		UserIdentity: &UserIdentity{
			Type:     "IAMUser",
			UserName: "testuser",
		},
	}

	err = tstore.Update(context.Background(), func(txn storage.Transaction) error {
		return manager.AddIndexInTxn(txn, event)
	})
	require.NoError(t, err)

	t.Run("Transaction-committed event index is visible via QueryByEventName", func(t *testing.T) {
		ids, err := manager.QueryByEventName([]string{"CreateTrail"}, 10)
		require.NoError(t, err)
		assert.Contains(t, ids, "evt-txn-001")
	})

	t.Run("Transaction-committed event index is visible via QueryByUsername", func(t *testing.T) {
		ids, err := manager.QueryByUsername("testuser", 10)
		require.NoError(t, err)
		assert.Contains(t, ids, "evt-txn-001")
	})

	t.Run("Transaction-committed event index is visible via QueryByEventSource", func(t *testing.T) {
		ids, err := manager.QueryByEventSource("cloudtrail.amazonaws.com", 10)
		require.NoError(t, err)
		assert.Contains(t, ids, "evt-txn-001")
	})
}
