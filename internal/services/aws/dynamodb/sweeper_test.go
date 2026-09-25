package dynamodb

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// Pins for the background sweeper family's per-store bodies: what each
// sweep does to one regional store, independent of the interval shell
// that schedules it (the shell's contract is pinned separately against
// the shared starter).

// journalCount reads the journal of one table through the replay surface.
func journalCount(t *testing.T, store dbstore.DynamoDBStoreInterface, table string) int {
	t.Helper()
	n := 0
	if err := store.Journal().ReverseReplay(table, time.Time{}, func(*dbstore.JournalChange) error {
		n++
		return nil
	}); err != nil {
		t.Fatalf("replay journal of %s: %v", table, err)
	}
	return n
}

// TestSweepStoreJournals pins the journal pruner's state contract: a
// table without recovery has no use for a journal and the sweep drops it
// whole, while a recovery-enabled table's journal survives the sweep —
// records inside the restorable window are never pruned.
func TestSweepStoreJournals(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	ctx := context.Background()
	store, err := svc.GetCachedStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("region store: %v", err)
	}
	setPitr := func(enabled bool) {
		t.Helper()
		if _, err := svc.UpdateContinuousBackups(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":                        "LegacyTable",
			"PointInTimeRecoverySpecification": map[string]interface{}{"PointInTimeRecoveryEnabled": enabled},
		}}); err != nil {
			t.Fatalf("set pitr %v: %v", enabled, err)
		}
	}
	put := func(id string) {
		t.Helper()
		if _, err := svc.PutItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "LegacyTable",
			"Item": map[string]interface{}{
				"id": map[string]interface{}{"S": id},
				"sk": map[string]interface{}{"S": id},
			},
		}}); err != nil {
			t.Fatalf("put %s: %v", id, err)
		}
	}

	// Recovery-enabled writes are journaled; disabling recovery makes the
	// journal useless and the sweep must drop it whole.
	setPitr(true)
	put("drop1")
	put("drop2")
	if n := journalCount(t, store, "LegacyTable"); n != 2 {
		t.Fatalf("journaled writes on recovery-enabled table = %d records, want 2", n)
	}
	setPitr(false)
	svc.sweepStoreJournals(store)
	if n := journalCount(t, store, "LegacyTable"); n != 0 {
		t.Fatalf("journal after sweep of recovery-disabled table = %d records, want 0", n)
	}

	// With recovery enabled again the fresh record sits well inside the
	// restorable window and the sweep must keep it.
	setPitr(true)
	put("keep")
	if n := journalCount(t, store, "LegacyTable"); n != 1 {
		t.Fatalf("journaled write after re-enable = %d records, want 1", n)
	}
	svc.sweepStoreJournals(store)
	if n := journalCount(t, store, "LegacyTable"); n != 1 {
		t.Fatalf("journal after sweep of recovery-enabled table = %d records, want 1 (kept)", n)
	}
}

// TestSweepStoreRetentions pins the retention pruner: stream records of
// a streaming table older than the retention window are trimmed. The
// sweeper's clock is injected past the window so the records, stamped at
// the real now, fall before the cutoff without waiting.
func TestSweepStoreRetentions(t *testing.T) {
	svc, store, _, _ := streamPlaneFixture(t)

	// The fixture's two seeded puts captured two stream records.
	seq, err := store.Streams().GetLatestSequenceForStream("StreamsPinTable", "")
	if err != nil || seq < 1 {
		t.Fatalf("seeded stream records: seq = %d, err = %v, want seq >= 1", seq, err)
	}

	restore := streamTimeNow
	streamTimeNow = func() time.Time { return time.Now().Add(25 * time.Hour) }
	defer func() { streamTimeNow = restore }()

	svc.sweepStoreRetentions(store)

	trimmed, err := store.Streams().GetLatestSequenceForStream("StreamsPinTable", "")
	if err != nil || trimmed != 0 {
		t.Fatalf("stream extent after retention sweep = %d, err = %v, want 0 (trimmed)", trimmed, err)
	}
}

// TestSweepStoreSystemBackups pins the system-backup pruner's keep half:
// an on-demand backup is outside the sweep by type, and a delete-time
// system backup inside its retention window survives — only expired
// system backups leave the store through the sweep.
func TestSweepStoreSystemBackups(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()
	store, err := svc.GetCachedStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("region store: %v", err)
	}

	createRestorableBackup(t, svc, reqCtx, "SweepBkTbl", "user-bk")
	if _, err := svc.UpdateContinuousBackups(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                        "SweepBkTbl",
		"PointInTimeRecoverySpecification": map[string]interface{}{"PointInTimeRecoveryEnabled": true},
	}}); err != nil {
		t.Fatalf("enable pitr: %v", err)
	}
	if _, err := svc.DeleteTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "SweepBkTbl",
	}}); err != nil {
		t.Fatalf("delete table: %v", err)
	}
	systemName := "SweepBkTbl" + dbstore.DeletedTableBackupSuffix
	names := listBackupNames(t, svc, reqCtx, "SweepBkTbl")
	if !containsName(names, "user-bk") || !containsName(names, systemName) {
		t.Fatalf("precondition backups = %v, want user-bk and %s", names, systemName)
	}

	svc.sweepStoreSystemBackups(store)

	after := listBackupNames(t, svc, reqCtx, "SweepBkTbl")
	if !containsName(after, "user-bk") {
		t.Fatalf("on-demand backup swept away: %v", after)
	}
	if !containsName(after, systemName) {
		t.Fatalf("unexpired system backup swept away: %v", after)
	}
}

// sweeperTestService builds a bare service with two regional stores and a
// cancellable background context, the minimum the shared starter needs.
func sweeperTestService(t *testing.T) (*DynamoDBService, context.CancelFunc) {
	t.Helper()
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	svc := &DynamoDBService{}
	svc.SetStorageManager(sm)
	bgCtx, bgCancel := context.WithCancel(context.Background())
	svc.bgCtx, svc.bgCancel = bgCtx, bgCancel
	t.Cleanup(func() { svc.Close(); sm.Close() })
	for _, region := range []string{"us-east-1", "eu-west-1"} {
		if _, err := svc.GetCachedStoreForRegion(region); err != nil {
			t.Fatalf("region store %s: %v", region, err)
		}
	}
	return svc, bgCancel
}

// TestStartIntervalSweeper pins the shared starter's contract: the sweep
// runs once per regional store on every tick, a second start behind the
// same spent once never registers a second goroutine, and cancelling the
// service context stops the sweeper.
func TestStartIntervalSweeper(t *testing.T) {
	svc, bgCancel := sweeperTestService(t)

	var mu sync.Mutex
	hits := map[dbstore.DynamoDBStoreInterface]int{}
	sweep := func(store dbstore.DynamoDBStoreInterface) {
		mu.Lock()
		hits[store]++
		mu.Unlock()
	}
	secondRuns := 0
	other := func(dbstore.DynamoDBStoreInterface) {
		mu.Lock()
		secondRuns++
		mu.Unlock()
	}

	var once sync.Once
	svc.startIntervalSweeper(&once, 2*time.Millisecond, "test sweep", sweep)
	svc.startIntervalSweeper(&once, 2*time.Millisecond, "test sweep", other)

	deadline := time.Now().Add(5 * time.Second)
	for {
		mu.Lock()
		both := len(hits) == 2
		for _, n := range hits {
			if n < 3 {
				both = false
			}
		}
		mu.Unlock()
		if both {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("sweeper did not reach every regional store: hits = %v", hits)
		}
		time.Sleep(2 * time.Millisecond)
	}

	mu.Lock()
	second := secondRuns
	mu.Unlock()
	if second != 0 {
		t.Fatalf("a second start behind the spent once ran %d times, want 0", second)
	}

	// Cancellation stops the sweeper: the goroutine leaves the service's
	// background wait group, and Wait returning is the proof — no in-flight
	// sweep remains to grow the counts, so no sleep window is needed.
	bgCancel()
	svc.bgWg.Wait()
	mu.Lock()
	total := 0
	for _, n := range hits {
		total += n
	}
	mu.Unlock()
	if total == 0 {
		t.Fatal("sweeper never ran before cancellation")
	}
}

// TestStartIntervalSweeperContainsPanics pins the containment half: a
// panicking sweep invocation is recovered and logged at that invocation
// alone — the ticker keeps scheduling, later ticks run the same pruner
// again, and the panic never escapes to the process. A pruner's own bug
// costs one interval, not the pruner.
func TestStartIntervalSweeperContainsPanics(t *testing.T) {
	svc, _ := sweeperTestService(t)

	var ran int32
	boom := func(dbstore.DynamoDBStoreInterface) {
		atomic.AddInt32(&ran, 1)
		panic("sweeper test panic")
	}

	var once sync.Once
	svc.startIntervalSweeper(&once, 2*time.Millisecond, "test sweep panic", boom)

	// The panicking sweep runs, is contained, and the ticker delivers
	// further invocations: reaching four panicking runs proves three
	// survivals after the first.
	deadline := time.Now().Add(5 * time.Second)
	for {
		if n := atomic.LoadInt32(&ran); n >= 4 {
			break
		} else if time.Now().After(deadline) {
			t.Fatalf("panicking sweep ran %d time(s); a panic stopped the ticker after the first", n)
		}
		time.Sleep(2 * time.Millisecond)
	}

	// The goroutine is still alive: the wait group has not released its
	// slot, and cancellation stops it promptly.
	svc.bgCancel()
	svc.bgWg.Wait()
}
