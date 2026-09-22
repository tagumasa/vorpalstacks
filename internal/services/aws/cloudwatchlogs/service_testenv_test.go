package cloudwatchlogs

// This file carries the shared test-environment construction of the
// service tests: the temporary storage manager, the service built over
// it (alone or over a caller-owned manager, the restart-fixture shape),
// and the group-creation tail the group-owning builders compose on. The
// per-file builders keep only their varying tails (groups, streams,
// invokers, buses, request contexts) on top of these.

import (
	"testing"

	"vorpalstacks/internal/core/storage"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// newTestStorageManager builds the temporary per-test storage manager the
// test services run over.
func newTestStorageManager(t *testing.T) *storage.RegionStorageManager {
	t.Helper()
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("region storage manager: %v", err)
	}
	return sm
}

// newTestServiceOnManager builds a service over a caller-owned storage
// manager — the shape the restart tests use to construct a second service
// over the first one's storage.
func newTestServiceOnManager(t *testing.T, sm *storage.RegionStorageManager) *LogsService {
	t.Helper()
	svc := NewLogsService(sm, "000000000000", t.TempDir())
	t.Cleanup(svc.Stop)
	return svc
}

// newTestService builds a LogsService over its own temporary storage,
// paired with the service's us-east-1 log store — the region every test
// builder composes onto.
func newTestService(t *testing.T) (*LogsService, *logsstore.Store) {
	t.Helper()
	svc := newTestServiceOnManager(t, newTestStorageManager(t))
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}
	return svc, store
}

// createTestLogGroup creates a log group in the us-east-1 store — the
// fixture step every group-owning builder shares.
func createTestLogGroup(t *testing.T, store *logsstore.Store, group string) {
	t.Helper()
	if err := store.CreateLogGroup(logsstore.NewLogGroup(group, "us-east-1", "000000000000")); err != nil {
		t.Fatalf("create log group: %v", err)
	}
}
