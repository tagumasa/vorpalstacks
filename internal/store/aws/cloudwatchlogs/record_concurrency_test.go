package cloudwatchlogs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

// Configuration mutations and ingestion write the same LogGroup record:
// when both sides run through the mutate seam / the ingestion mutex, the
// StoredBytes accounting of every successful PutLogEvents survives the
// interleaving — a lost update is a conservation violation.
func TestLogGroupMutationSerialisesWithIngestion(t *testing.T) {
	s := newLogsTestStore(t)

	if err := s.CreateLogGroup(&LogGroup{Name: "group"}); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateLogStream(&LogStream{LogGroupName: "group", Name: "stream"}); err != nil {
		t.Fatal(err)
	}

	const writers = 4
	const batches = 25
	message := strings.Repeat("x", 64)

	var succeeded atomic.Int64
	var wg sync.WaitGroup
	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < batches; i++ {
				_, err := s.PutLogEvents("group", "stream", []LogEntry{{
					Message:   message,
					Timestamp: time.Now().UnixMilli() + int64(i),
				}})
				if err == nil {
					succeeded.Add(1)
				}
			}
		}()
	}
	for w := 0; w < 2; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < batches; i++ {
				_ = s.MutateLogGroup("group", func(lg *LogGroup) error {
					lg.KmsKeyId = fmt.Sprintf("key-%d", id)
					return nil
				})
			}
		}(w)
	}
	wg.Wait()

	lg, err := s.GetLogGroup("group")
	if err != nil {
		t.Fatal(err)
	}
	want := int64(succeeded.Load()) * int64(len(message))
	if lg.StoredBytes != want {
		t.Fatalf("lost update: StoredBytes=%d, want %d (%d successful puts)",
			lg.StoredBytes, want, succeeded.Load())
	}
	if succeeded.Load() == 0 {
		t.Fatal("no PutLogEvents succeeded; the oracle is vacuous")
	}
}

// A PutLogEvents that is in flight while DeleteLogStream runs must not
// resurrect the deleted stream: the delete holds the ingestion mutex, so
// a write either commits before the delete (and is removed with the
// stream) or fails the stream existence check once it acquires the mutex.
// After the delete returns, the stream record must be absent and further
// writes rejected.
func TestDeleteLogStreamBlocksInFlightIngestionResurrection(t *testing.T) {
	s := newLogsTestStore(t)

	if err := s.CreateLogGroup(&LogGroup{Name: "group"}); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateLogStream(&LogStream{LogGroupName: "group", Name: "stream"}); err != nil {
		t.Fatal(err)
	}

	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		for {
			select {
			case <-stop:
				return
			default:
			}
			_, _ = s.PutLogEvents("group", "stream", []LogEntry{{
				Message:   "concurrent-write",
				Timestamp: time.Now().UnixMilli() + int64(i),
			}})
			i++
		}
	}()

	time.Sleep(5 * time.Millisecond)
	if err := s.DeleteLogStream("group", "stream"); err != nil {
		t.Fatal(err)
	}
	close(stop)
	wg.Wait()

	if _, err := s.GetLogStream("group", "stream"); err == nil {
		t.Fatal("deleted log stream resurrected by an in-flight PutLogEvents")
	}
	if _, err := s.PutLogEvents("group", "stream", []LogEntry{{
		Message:   "after-delete",
		Timestamp: time.Now().UnixMilli(),
	}}); err == nil {
		t.Fatal("PutLogEvents after DeleteLogStream succeeded")
	}
}

// MetricFilterCount is maintained transactionally against the same
// LogGroup record the configuration mutations rewrite: with both sides
// under the ingestion mutex, the count equals the net number of live
// filters after the interleaving.
func TestMetricFilterCountConservedUnderConfigMutation(t *testing.T) {
	s := newLogsTestStore(t)

	if err := s.CreateLogGroup(&LogGroup{Name: "group"}); err != nil {
		t.Fatal(err)
	}

	const created = 8
	const deleted = 4

	var wg sync.WaitGroup
	for i := 0; i < created; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_ = s.PutMetricFilter(NewMetricFilter("group", fmt.Sprintf("filter-%d", id), "ERROR", nil))
		}(i)
	}
	for i := 0; i < deleted; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// The create goroutines run concurrently; retry until the
			// filter exists so a delete is not lost to ordering.
			for attempt := 0; attempt < 200; attempt++ {
				if err := s.DeleteMetricFilter("group", fmt.Sprintf("filter-%d", id)); err == nil {
					return
				}
				time.Sleep(time.Millisecond)
			}
		}(i)
	}
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_ = s.MutateLogGroup("group", func(lg *LogGroup) error {
				lg.KmsKeyId = fmt.Sprintf("key-%d", id)
				return nil
			})
		}(i)
	}
	wg.Wait()

	lg, err := s.GetLogGroup("group")
	if err != nil {
		t.Fatal(err)
	}
	if lg.MetricFilterCount != int32(created-deleted) {
		t.Fatalf("MetricFilterCount=%d, want %d", lg.MetricFilterCount, created-deleted)
	}
}

// Concurrent creators of one log group name must admit exactly one
// record: the existence check and the write share the ingestion mutex,
// so the second creator to acquire it observes the first record.
func TestCreateLogGroupAdmitsExactlyOneDuplicate(t *testing.T) {
	s := newLogsTestStore(t)

	const creators = 32
	var succeeded atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < creators; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.CreateLogGroup(&LogGroup{Name: "dupe"}); err == nil {
				succeeded.Add(1)
			} else if !errors.Is(err, ErrLogGroupAlreadyExists) {
				t.Errorf("unexpected create error: %v", err)
			}
		}()
	}
	wg.Wait()

	if got := succeeded.Load(); got != 1 {
		t.Fatalf("duplicate admission: %d creates succeeded, want exactly 1", got)
	}
}

// The same single-admission contract holds on the stream plane.
func TestCreateLogStreamAdmitsExactlyOneDuplicate(t *testing.T) {
	s := newLogsTestStore(t)

	if err := s.CreateLogGroup(&LogGroup{Name: "group"}); err != nil {
		t.Fatal(err)
	}

	const creators = 32
	var succeeded atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < creators; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.CreateLogStream(&LogStream{LogGroupName: "group", Name: "dupe"}); err == nil {
				succeeded.Add(1)
			} else if !errors.Is(err, ErrLogStreamAlreadyExists) {
				t.Errorf("unexpected create error: %v", err)
			}
		}()
	}
	wg.Wait()

	if got := succeeded.Load(); got != 1 {
		t.Fatalf("duplicate admission: %d creates succeeded, want exactly 1", got)
	}
}

// A stream created while the group teardown is in flight cannot
// resurrect the family: the create's existence checks share the
// ingestion mutex with the delete, so the create either commits before
// the teardown (and is removed with the rest) or fails the group check
// once the delete has committed.
func TestCreateDuringGroupTeardownCannotResurrect(t *testing.T) {
	s := newLogsTestStore(t)

	if err := s.CreateLogGroup(&LogGroup{Name: "group"}); err != nil {
		t.Fatal(err)
	}

	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		for {
			select {
			case <-stop:
				return
			default:
			}
			_ = s.CreateLogStream(&LogStream{
				LogGroupName: "group",
				Name:         fmt.Sprintf("stream-%d", i),
			})
			i++
		}
	}()

	time.Sleep(5 * time.Millisecond)
	if err := s.DeleteLogGroup("group"); err != nil {
		t.Fatal(err)
	}
	close(stop)
	wg.Wait()

	if _, err := s.GetLogGroup("group"); err == nil {
		t.Fatal("deleted log group resurrected by a concurrent CreateLogStream")
	}
	for i := 0; i < 512; i++ {
		if _, err := s.GetLogStream("group", fmt.Sprintf("stream-%d", i)); err == nil {
			t.Fatalf("stream-%d survived the group teardown", i)
		}
	}
}

// Deletion protection must hold even when it is enabled while a delete
// is in flight: the guarded variant decides inside the delete's critical
// section. The plain variant is the create-path rollback's deliberate
// bypass and keeps deleting protected records.
func TestGuardedDeleteHonoursDeletionProtection(t *testing.T) {
	s := newLogsTestStore(t)

	if err := s.CreateLogGroup(&LogGroup{Name: "group"}); err != nil {
		t.Fatal(err)
	}
	if err := s.MutateLogGroup("group", func(lg *LogGroup) error {
		lg.DeletionProtectionEnabled = true
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteLogGroupIfUnprotected("group"); !errors.Is(err, ErrLogGroupDeletionProtected) {
		t.Fatalf("guarded delete error = %v, want ErrLogGroupDeletionProtected", err)
	}
	if _, err := s.GetLogGroup("group"); err != nil {
		t.Fatal("protected log group deleted by the guarded path")
	}

	// The rollback contract: the plain delete bypasses the guard for a
	// group the caller just created (CreateLogGroup's tag-failure path).
	if err := s.DeleteLogGroup("group"); err != nil {
		t.Fatalf("plain delete on a protected group failed: %v", err)
	}
}

// A terminal task record keeps its status against a later writer: the
// mutate seam decides on the persisted record under the task mutex, so
// a canceller's CANCELLED cannot be overwritten by the worker's stale
// COMPLETED snapshot.
func TestMutateExportTaskTerminalStatusWins(t *testing.T) {
	s := newLogsTestStore(t)

	if err := s.PutExportTask(&ExportTask{TaskId: "t", Status: "CANCELLED"}); err != nil {
		t.Fatal(err)
	}
	err := s.MutateExportTask("t", func(task *ExportTask) error {
		if task.Status == "COMPLETED" || task.Status == "FAILED" || task.Status == "CANCELLED" {
			return nil
		}
		task.Status = "COMPLETED"
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	task, err := s.GetExportTask("t")
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != "CANCELLED" {
		t.Fatalf("terminal status overwritten: %s", task.Status)
	}
}

// The scheduled execution plane's StopQuery contract: the atomic
// RUNNING→CANCELLED flip and the finalise gate together keep a late
// worker terminal write from overwriting the canceller's status.
func TestScheduledExecutionCancellationWinsOverLateWorkerWrite(t *testing.T) {
	s := newLogsTestStore(t)

	exec := &ScheduledQueryExecution{
		ScheduledQueryId: "sq-1", TriggerTime: 100, QueryId: "q-1",
		Status: ScheduledExecutionStatusRunning,
	}
	if err := s.PutScheduledQueryExecution(exec); err != nil {
		t.Fatal(err)
	}

	cancelled, err := s.CancelScheduledQueryExecutionIfRunning(exec)
	if err != nil {
		t.Fatal(err)
	}
	if !cancelled {
		t.Fatal("running execution was not cancelled")
	}
	// A second cancel on the now-terminal record is a no-op.
	cancelled, err = s.CancelScheduledQueryExecutionIfRunning(exec)
	if err != nil {
		t.Fatal(err)
	}
	if cancelled {
		t.Fatal("terminal execution flipped again")
	}

	late := *exec
	late.Status = ScheduledExecutionStatusSuccess
	late.RecordsMatched = 7
	if err := s.FinaliseScheduledQueryExecution(&late); err != nil {
		t.Fatal(err)
	}

	execs, err := s.ListScheduledQueryExecutions("sq-1", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(execs) != 1 {
		t.Fatalf("execution count = %d, want 1", len(execs))
	}
	if execs[0].Status != ScheduledExecutionStatusCancelled {
		t.Fatalf("canceller's status overwritten by the late worker write: %s", execs[0].Status)
	}
}

// The delivery record's shaping fields and the engine cursor are
// maintained by different writers: through the mutate seam both survive
// an interleaving — a lost update is a conservation violation.
func TestDeliveryShapingSurvivesConcurrentCursorAdvances(t *testing.T) {
	s := newLogsTestStore(t)

	delivery := &Delivery{
		Id: "d-1", DeliverySourceName: "src", DeliveryDestinationArn: "arn:dst",
		RecordFields: []string{"message"},
	}
	if existed, err := s.PutDeliveryIfPairAbsent(delivery); err != nil || existed {
		t.Fatalf("seed: existed=%v err=%v", existed, err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = s.MutateDelivery("d-1", func(d *Delivery) error {
				d.CursorTime = int64(n)
				d.CursorStream = fmt.Sprintf("stream-%d", n)
				return nil
			})
		}(i)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.MutateDelivery("d-1", func(d *Delivery) error {
			d.FieldDelimiter = ","
			return nil
		})
		if err != nil {
			t.Errorf("shaping update: %v", err)
		}
	}()
	wg.Wait()

	got, err := s.GetDelivery("d-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.FieldDelimiter != "," {
		t.Fatal("shaping update lost under concurrent cursor advances")
	}
}

// Two concurrent creations of the same (source, destination) delivery
// pair admit exactly one record.
func TestPutDeliveryIfPairAbsentAdmitsExactlyOne(t *testing.T) {
	s := newLogsTestStore(t)

	var admitted atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d := &Delivery{
				Id:                     fmt.Sprintf("d-%d", time.Now().UnixNano()),
				DeliverySourceName:     "src",
				DeliveryDestinationArn: "arn:dst",
			}
			existed, err := s.PutDeliveryIfPairAbsent(d)
			if err != nil {
				t.Errorf("create: %v", err)
				return
			}
			if !existed {
				admitted.Add(1)
			}
		}()
	}
	wg.Wait()

	if got := admitted.Load(); got != 1 {
		t.Fatalf("duplicate pair admission: %d creates admitted, want exactly 1", got)
	}
}

// A delete that lands between a mutation's read and write cannot
// resurrect the record: the delete and the mutate share the family
// mutex, so once DeleteLookupTable returns no concurrent mutation's
// write survives — the mutation either completed before the delete or
// fails to load the deleted record.
func TestDeleteLookupTableSurvivesConcurrentMutations(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.PutLookupTable(&LookupTable{Name: "raced"}); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_ = s.MutateLookupTable("raced", func(lt *LookupTable) error {
				lt.Description = "refreshed"
				return nil
			})
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		if err := s.DeleteLookupTable("raced"); err != nil && err != ErrResourceNotFound {
			t.Errorf("delete: %v", err)
		}
	}()
	close(start)
	wg.Wait()

	if _, err := s.GetLookupTable("raced"); err == nil {
		t.Fatal("the deleted lookup table was resurrected by a concurrent mutation")
	}
}

// The JSON-regime getters distinguish a missing record from a broken
// one: not-found answers the sentinel, a record the getter cannot decode
// propagates the storage error — a corrupt record or a storage outage
// must not surface as ResourceNotFoundException.
func TestJSONGetterDistinguishesNotFoundFromStorageError(t *testing.T) {
	s := newLogsTestStore(t)

	if err := s.PutRaw(s.lookupTableKey("corrupt"), []byte("{not json")); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetLookupTable("corrupt"); err == nil || err == ErrResourceNotFound {
		t.Fatalf("corrupt record: %v, want the storage error to propagate", err)
	}
	if _, err := s.GetLookupTable("absent"); err != ErrResourceNotFound {
		t.Fatalf("absent record: %v, want ErrResourceNotFound", err)
	}

	// The ingestion plane keeps the same distinction: a group record
	// that cannot decode is not "log group not found".
	if err := s.CreateLogGroup(NewLogGroup("corrupt-group", "us-east-1", "000000000000")); err != nil {
		t.Fatal(err)
	}
	if err := s.PutRaw(s.logGroupKey("corrupt-group"), []byte("\xff\xfe not proto")); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := s.GetLogEvents("corrupt-group", "any", 0, 0, 10, true, ""); err == nil || err == ErrLogGroupNotFound {
		t.Fatalf("corrupt group record: %v, want the storage error to propagate", err)
	}
}

// The JSON-regime listings share one scan body: a record that fails to
// decode is skipped without failing the listing, and the skip leaves a
// corruption warning behind — the per-family drift the shared body
// exists to make impossible would otherwise let a corrupt record vanish
// from a listing with no trace.
func TestListJSONRecordsLogsCorruption(t *testing.T) {
	s := newLogsTestStore(t)

	if err := s.PutScheduledQuery(&ScheduledQuery{Id: "healthy", Name: "q", State: "ENABLED"}); err != nil {
		t.Fatal(err)
	}
	if err := s.PutRaw(s.scheduledQueryKey("corrupt"), []byte("{not json")); err != nil {
		t.Fatal(err)
	}

	captured := captureProcessStdout(t, func() {
		queries, err := s.ListScheduledQueries("")
		if err != nil {
			t.Fatalf("listing with a corrupt sibling: %v", err)
		}
		if len(queries) != 1 || queries[0].Id != "healthy" {
			t.Fatalf("healthy record must survive a corrupt sibling: %+v", queries)
		}
	})
	if !strings.Contains(captured, "Corrupt scheduled query record") {
		t.Fatalf("corrupt record skipped without a corruption warning, captured: %q", captured)
	}
}

// captureProcessStdout redirects the process's standard output to a pipe
// for the duration of fn and returns what was written. The package-level
// logger binds the os.Stdout file at first use, so the descriptor — not
// the os.Stdout variable — is the only capture point that works
// regardless of when the logger was created.
func captureProcessStdout(t *testing.T, fn func()) (captured string) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	saved, err := unix.Dup(1)
	if err != nil {
		t.Fatalf("duplicate stdout: %v", err)
	}
	if err := unix.Dup3(int(w.Fd()), 1, 0); err != nil {
		t.Fatalf("redirect stdout: %v", err)
	}
	defer func() {
		if err := unix.Dup3(saved, 1, 0); err != nil {
			t.Errorf("restore stdout: %v", err)
		}
		if err := unix.Close(saved); err != nil {
			t.Errorf("close saved stdout: %v", err)
		}
		if err := w.Close(); err != nil {
			t.Errorf("close pipe writer: %v", err)
		}
		data, readErr := io.ReadAll(r)
		if readErr != nil {
			t.Errorf("read captured stdout: %v", readErr)
		}
		captured = string(data)
	}()
	fn()
	return
}

// A per-group configuration put racing the group's teardown must not
// leave its record outliving the deleted group: every family's put holds
// the teardown lock and re-checks the group inside it, so whatever the
// interleaving, a successful DeleteLogGroup leaves no family record
// behind — the put either ran before the teardown (and was torn down
// with everything else) or failed the re-check.
func TestConfigRecordsNeverOutliveRacingGroupDelete(t *testing.T) {
	families := []struct {
		name      string
		put       func(s *Store, group string) error
		remaining func(s *Store, group string) int
	}{
		{"subscription filter", func(s *Store, group string) error {
			return s.PutSubscriptionFilterWithLimitCheck(&SubscriptionFilter{
				LogGroupName:   group,
				FilterName:     "race",
				DestinationArn: s.arnBuilder.CloudWatch().LogGroup(group),
			}, MaxSubscriptionFiltersPerLogGroup)
		}, func(s *Store, group string) int {
			filters, err := s.ListSubscriptionFilters(group, "")
			if err != nil {
				t.Error(err)
			}
			return len(filters)
		}},
		{"bearer token switch", func(s *Store, group string) error {
			// The switch rides the group record: the racing mutate either
			// lands before the teardown (and dies with the record) or
			// fails the re-read.
			return s.MutateLogGroup(group, func(lg *LogGroup) error {
				lg.BearerTokenAuthenticationEnabled = true
				return nil
			})
		}, func(s *Store, group string) int {
			if _, err := s.GetLogGroup(group); err == nil {
				return 1
			}
			return 0
		}},
		{"data protection policy", func(s *Store, group string) error {
			return s.PutDataProtectionPolicy(&DataProtectionPolicy{
				LogGroupIdentifier: group,
				PolicyDocument:     `{"Name":"race"}`,
			})
		}, func(s *Store, group string) int {
			if s.Exists(s.dataProtectionPolicyKey(group)) {
				return 1
			}
			return 0
		}},
		{"transformer", func(s *Store, group string) error {
			return s.PutTransformer(&Transformer{LogGroupName: group})
		}, func(s *Store, group string) int {
			if s.Exists(s.transformerKey(group)) {
				return 1
			}
			return 0
		}},
	}

	for _, family := range families {
		for iter := 0; iter < 30; iter++ {
			s := newLogsTestStore(t)
			group := fmt.Sprintf("race-%s-%d", family.name, iter)
			if err := s.CreateLogGroup(&LogGroup{Name: group}); err != nil {
				t.Fatal(err)
			}

			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				// The racing put's outcome is the interleaving's choice;
				// only the invariant after the delete is fixed.
				_ = family.put(s, group)
			}()
			go func() {
				defer wg.Done()
				if err := s.DeleteLogGroup(group); err != nil {
					t.Errorf("family %s iteration %d: DeleteLogGroup: %v", family.name, iter, err)
				}
			}()
			wg.Wait()

			if left := family.remaining(s, group); left != 0 {
				t.Fatalf("family %s left %d record(s) outliving deleted group %s",
					family.name, left, group)
			}
		}
	}
}

// Every per-group configuration family rejects a group the teardown
// already removed: the store-level re-check closes the window the
// service layer's earlier group check leaves open.
func TestConfigPutsRejectDeletedGroup(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.CreateLogGroup(&LogGroup{Name: "gone"}); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteLogGroup("gone"); err != nil {
		t.Fatal(err)
	}

	if err := s.PutSubscriptionFilterWithLimitCheck(&SubscriptionFilter{
		LogGroupName:   "gone",
		FilterName:     "f",
		DestinationArn: "arn:aws:logs:us-east-1:000000000000:destination:f",
	}, MaxSubscriptionFiltersPerLogGroup); !errors.Is(err, ErrLogGroupNotFound) {
		t.Fatalf("subscription filter put on deleted group: %v, want ErrLogGroupNotFound", err)
	}
	if err := s.MutateLogGroup("gone", func(lg *LogGroup) error {
		lg.BearerTokenAuthenticationEnabled = true
		return nil
	}); !errors.Is(err, ErrLogGroupNotFound) {
		t.Fatalf("bearer token put on deleted group: %v, want ErrLogGroupNotFound", err)
	}
	if err := s.PutDataProtectionPolicy(&DataProtectionPolicy{
		LogGroupIdentifier: "gone",
		PolicyDocument:     `{"Name":"p"}`,
	}); !errors.Is(err, ErrLogGroupNotFound) {
		t.Fatalf("data protection policy put on deleted group: %v, want ErrLogGroupNotFound", err)
	}
	if err := s.PutTransformer(&Transformer{LogGroupName: "gone"}); !errors.Is(err, ErrLogGroupNotFound) {
		t.Fatalf("transformer put on deleted group: %v, want ErrLogGroupNotFound", err)
	}
}

// A GetLogEvents page served while DeleteLogStream runs must never be
// torn: the gather holds the chunk plane's read lock and the teardown its
// write side, so every page observes the complete pre-delete set or the
// post-delete state — never a partial listing whose files the delete
// already removed. The fixture spreads the events over several chunks so
// a torn gather would surface as a partial count, not a whole-stream
// skip.
func TestGetLogEventsPageNeverInterleavesStreamDelete(t *testing.T) {
	s := newLogsTestStore(t)
	const group, stream = "torn-group", "torn-stream"
	if err := s.CreateLogGroup(&LogGroup{Name: group}); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateLogStream(&LogStream{LogGroupName: group, Name: stream}); err != nil {
		t.Fatal(err)
	}

	const batches = 5
	const total = batches * 2
	ts := time.Now().UnixMilli()
	for b := 0; b < batches; b++ {
		batch := []LogEntry{
			{Timestamp: ts + int64(2*b), Message: fmt.Sprintf("m%02d-a", b)},
			{Timestamp: ts + int64(2*b+1), Message: fmt.Sprintf("m%02d-b", b)},
		}
		if _, err := s.PutLogEvents(group, stream, batch); err != nil {
			t.Fatal(err)
		}
	}

	stop := make(chan struct{})
	var wg sync.WaitGroup
	for r := 0; r < 4; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				events, _, _, err := s.GetLogEvents(group, stream, 0, 0, 0, true, "")
				if err != nil {
					if errors.Is(err, ErrLogStreamNotFound) || errors.Is(err, ErrLogGroupNotFound) {
						return // the delete completed
					}
					t.Error(err)
					return
				}
				// An empty page with no error is the post-delete chunk
				// set observed by a reader whose stream existence check
				// passed before the teardown — a consistent view. A
				// partial count is the torn page this pin forbids.
				if len(events) != 0 && len(events) != total {
					t.Errorf("torn page: %d events, want 0 or %d", len(events), total)
					return
				}
			}
		}()
	}
	time.Sleep(2 * time.Millisecond)
	if err := s.DeleteLogStream(group, stream); err != nil {
		t.Fatal(err)
	}
	close(stop)
	wg.Wait()
}

// The mutate seam conserves concurrent field updates: two mutator
// families (StoredBytes counters and configuration overwrites) ride the
// same locked read-modify-write, so neither side's fields are lost. A
// full-record overwrite that reads its snapshot outside the seam is a
// stale-snapshot writer the family mutex cannot and does not save —
// overwrites belong to the mutate seam, which is what this pins.
func TestMutateSeamConservesConcurrentFieldWrites(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.CreateLogGroup(&LogGroup{Name: "group"}); err != nil {
		t.Fatal(err)
	}

	const mutators = 2
	const batches = 50
	var wg sync.WaitGroup
	for w := 0; w < mutators; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < batches; i++ {
				_ = s.MutateLogGroup("group", func(lg *LogGroup) error {
					lg.StoredBytes++
					return nil
				})
			}
		}(w)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < batches; i++ {
			_ = s.MutateLogGroup("group", func(lg *LogGroup) error {
				lg.KmsKeyId = fmt.Sprintf("overwrite-%d", i)
				return nil
			})
		}
	}()
	wg.Wait()

	lg, err := s.GetLogGroup("group")
	if err != nil {
		t.Fatal(err)
	}
	want := int64(mutators * batches)
	if lg.StoredBytes != want {
		t.Fatalf("lost increments: StoredBytes=%d, want %d", lg.StoredBytes, want)
	}
	if lg.KmsKeyId == "" {
		t.Fatal("the overwrite writer never landed; the oracle is vacuous")
	}
}

// The per-group lock shards the mutation plane: ingestions to different
// groups run under different lock domains and must all land — no lost
// StoredBytes on any group's record, no dropped batch on any stream —
// while every same-group interleaving keeps the conservation the
// single-lock design guaranteed. The race detector watches the shared
// shard table through the whole fan-out.
func TestConcurrentIngestAcrossGroups(t *testing.T) {
	s := newLogsTestStore(t)

	const groups = 8
	const streams = 2
	const batches = 4

	msg := func(group, stream string, b int) string {
		return fmt.Sprintf("%s/%s/%d", group, stream, b)
	}

	for g := 0; g < groups; g++ {
		name := fmt.Sprintf("group-%d", g)
		if err := s.CreateLogGroup(&LogGroup{Name: name}); err != nil {
			t.Fatal(err)
		}
		for st := 0; st < streams; st++ {
			if err := s.CreateLogStream(&LogStream{
				LogGroupName: name, Name: fmt.Sprintf("stream-%d", st),
			}); err != nil {
				t.Fatal(err)
			}
		}
	}

	var wg sync.WaitGroup
	for g := 0; g < groups; g++ {
		group := fmt.Sprintf("group-%d", g)
		for st := 0; st < streams; st++ {
			stream := fmt.Sprintf("stream-%d", st)
			wg.Add(1)
			go func() {
				defer wg.Done()
				for b := 0; b < batches; b++ {
					if _, err := s.PutLogEvents(group, stream, []LogEntry{{
						Timestamp: int64(1700000000000 + b*1000),
						Message:   msg(group, stream, b),
					}}); err != nil {
						t.Errorf("put %s/%s batch %d: %v", group, stream, b, err)
					}
				}
			}()
		}
	}
	wg.Wait()

	for g := 0; g < groups; g++ {
		group := fmt.Sprintf("group-%d", g)
		lg, err := s.GetLogGroup(group)
		if err != nil {
			t.Fatal(err)
		}
		var wantBytes int64
		for st := 0; st < streams; st++ {
			stream := fmt.Sprintf("stream-%d", st)
			for b := 0; b < batches; b++ {
				wantBytes += int64(len(msg(group, stream, b)))
			}
			events, _, _, err := s.GetLogEvents(group, stream, 0, 0, 0, true, "")
			if err != nil {
				t.Fatal(err)
			}
			if len(events) != batches {
				t.Fatalf("stream %s/%s: got %d events, want %d", group, stream, len(events), batches)
			}
		}
		if lg.StoredBytes != wantBytes {
			t.Fatalf("group %s: StoredBytes=%d, want %d", group, lg.StoredBytes, wantBytes)
		}
	}
}
