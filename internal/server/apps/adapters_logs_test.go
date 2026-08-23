package apps

import (
	"context"
	"testing"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	svclogs "vorpalstacks/internal/services/aws/cloudwatchlogs"
)

func newLogsInvokerTestService(t *testing.T) *svclogs.LogsService {
	t.Helper()
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("region storage manager: %v", err)
	}
	svc := svclogs.NewLogsService(sm, "000000000000", t.TempDir())
	t.Cleanup(svc.Stop)
	return svc
}

// TestLogsInvokerWritesVisibleToServingStore pins read-your-writes between
// the LogsInvoker adapter and the store instance that serves CloudWatch Logs
// API reads: the adapter resolves the LogsService-owned store, so writes are
// immediately visible through the API read plane.
func TestLogsInvokerWritesVisibleToServingStore(t *testing.T) {
	svc := newLogsInvokerTestService(t)
	adapter := &logsInvokerAdapter{provider: svc}
	ctx := context.Background()
	if err := adapter.EnsureLogGroup(ctx, "us-east-1", "group", "000000000000"); err != nil {
		t.Fatalf("EnsureLogGroup: %v", err)
	}
	if err := adapter.EnsureLogStream(ctx, "us-east-1", "group", "stream"); err != nil {
		t.Fatalf("EnsureLogStream: %v", err)
	}
	if err := adapter.PutLogEvents(ctx, "us-east-1", "group", "stream", []eventbus.LogsLogEntry{
		{Timestamp: 1000, Message: "first"},
		{Timestamp: 2000, Message: "second"},
	}); err != nil {
		t.Fatalf("PutLogEvents: %v", err)
	}

	serving, err := svc.GetStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("serving store: %v", err)
	}
	events, _, _, err := serving.GetLogEvents("group", "stream", 0, 0, 0, true, "")
	if err != nil {
		t.Fatalf("GetLogEvents: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("serving store must see both adapter-written events immediately, got %d", len(events))
	}
	if events[0].Message != "first" || events[1].Message != "second" {
		t.Fatalf("unexpected event order: %q then %q", events[0].Message, events[1].Message)
	}
}

// TestLogsInvokerEnsureIdempotent pins that a second Ensure call for an
// existing group and stream succeeds (concurrent creation races resolve to
// "already exists", which the adapter tolerates).
func TestLogsInvokerEnsureIdempotent(t *testing.T) {
	svc := newLogsInvokerTestService(t)
	adapter := &logsInvokerAdapter{provider: svc}
	ctx := context.Background()
	for i := 0; i < 2; i++ {
		if err := adapter.EnsureLogGroup(ctx, "us-east-1", "group", "000000000000"); err != nil {
			t.Fatalf("EnsureLogGroup call %d: %v", i+1, err)
		}
		if err := adapter.EnsureLogStream(ctx, "us-east-1", "group", "stream"); err != nil {
			t.Fatalf("EnsureLogStream call %d: %v", i+1, err)
		}
	}
}

// TestLogsInvokerRegionIsolated pins that a group written in one region is
// not visible in another region's store.
func TestLogsInvokerRegionIsolated(t *testing.T) {
	svc := newLogsInvokerTestService(t)
	adapter := &logsInvokerAdapter{provider: svc}
	ctx := context.Background()
	if err := adapter.EnsureLogGroup(ctx, "us-west-2", "group", "000000000000"); err != nil {
		t.Fatalf("EnsureLogGroup: %v", err)
	}
	other, err := svc.GetStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("other region store: %v", err)
	}
	if _, _, _, err := other.GetLogEvents("group", "stream", 0, 0, 0, true, ""); err == nil {
		t.Fatal("us-west-2 group must not be visible in the us-east-1 store")
	}
}
