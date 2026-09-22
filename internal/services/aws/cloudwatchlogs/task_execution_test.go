package cloudwatchlogs

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/eventbus"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// recordingS3Invoker captures the objects the task runners write, keyed by
// bucket/key, and serves them back through GetObject. Seeded objects can
// be marked unreadable to pin per-object failure accounting.
type recordingS3Invoker struct {
	mu      sync.Mutex
	objects map[string][]byte
	buckets map[string]bool
}

func newRecordingS3Invoker() *recordingS3Invoker {
	return &recordingS3Invoker{objects: map[string][]byte{}, buckets: map[string]bool{}}
}

func (r *recordingS3Invoker) PutObject(_ context.Context, _, bucket, key string, data []byte, _ string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.objects[bucket+"/"+key] = append([]byte(nil), data...)
	return nil
}

func (r *recordingS3Invoker) PutObjectWithMetadata(ctx context.Context, region, bucket, key string, data []byte, contentType string, _ map[string]string) error {
	return r.PutObject(ctx, region, bucket, key, data, contentType)
}

func (r *recordingS3Invoker) GetObject(_ context.Context, _, bucket, key string, _ int64) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if data, ok := r.objects[bucket+"/"+key]; ok {
		return append([]byte(nil), data...), nil
	}
	return nil, logsstore.ErrResourceNotFound
}

func (r *recordingS3Invoker) GetObjectVersion(ctx context.Context, region, bucket, key, _ string, maxBytes int64) ([]byte, error) {
	return r.GetObject(ctx, region, bucket, key, maxBytes)
}

func (r *recordingS3Invoker) ListObjects(_ context.Context, _, bucket, prefix string, _ int) ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var keys []string
	for k := range r.objects {
		if strings.HasPrefix(k, bucket+"/"+prefix) {
			keys = append(keys, strings.TrimPrefix(k, bucket+"/"))
		}
	}
	return keys, nil
}

func (r *recordingS3Invoker) ListObjectEntries(_ context.Context, _, _, _ string, _ int) ([]invokers.S3ObjectEntry, error) {
	return nil, nil
}

func (r *recordingS3Invoker) BucketExists(_ context.Context, _, bucket string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.buckets[bucket], nil
}

func (r *recordingS3Invoker) GetBucketPolicy(_ context.Context, _, _ string) (string, error) {
	return "", nil
}

func (r *recordingS3Invoker) EnsureBucket(_ context.Context, _, bucket string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buckets[bucket] = true
	return nil
}

func (r *recordingS3Invoker) DeleteObject(_ context.Context, _, bucket, key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.objects, bucket+"/"+key)
	return nil
}

func (r *recordingS3Invoker) object(bucket, key string) ([]byte, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	data, ok := r.objects[bucket+"/"+key]
	return data, ok
}

// An export over a stream holding more events than one GetLogEvents page
// serves writes every event to the destination object.
func TestExecuteExportTaskFetchesAllEvents(t *testing.T) {
	const group, stream = "export-group", "export-stream"
	svc, store := newReadTestService(t, group)
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, group)); err != nil {
		t.Fatal(err)
	}

	const totalEvents = 10005
	base := time.Now().UnixMilli() - int64(totalEvents)
	putReadEvents(t, store, group, stream, base, totalEvents)

	s3 := newRecordingS3Invoker()
	bus := eventbus.NewEventBus()
	bus.SetS3Invoker(s3)
	busCtx, busCancel := context.WithCancel(context.Background())
	if err := bus.Start(busCtx); err != nil {
		t.Fatalf("bus start: %v", err)
	}
	t.Cleanup(busCancel)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatalf("SetEventBus: %v", err)
	}

	taskId := "export-fetchall-1"
	if err := store.PutExportTask(&logsstore.ExportTask{
		TaskId: taskId, LogGroupName: group, Status: "RUNNING",
		From: base - 1000, To: base + int64(totalEvents) + 1000,
		Destination: "export-bucket", CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}

	svc.executeExportTask(context.Background(), "us-east-1", group, "", base-1000, base+int64(totalEvents)+1000, "export-bucket", "", taskId)

	task, err := store.GetExportTask(taskId)
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != "COMPLETED" {
		t.Fatalf("task status = %q (%s), want COMPLETED", task.Status, task.StatusMessage)
	}

	// The documented default destinationPrefix applies: "If you don't
	// specify a value, the default is exportedlogs."
	data, ok := s3.object("export-bucket", "exportedlogs/"+taskId+"/exportedlogs.gz")
	if !ok {
		t.Fatal("export object not written")
	}
	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("export object is not gzip: %v", err)
	}
	plain, err := io.ReadAll(gr)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(string(plain), "\n"); got != totalEvents {
		t.Fatalf("exported %d records, want %d", got, totalEvents)
	}
}

// An export with no delivery path reports FAILED, not a silent COMPLETED
// with nothing delivered.
func TestExecuteExportTaskWithoutInvokerFails(t *testing.T) {
	const group, stream = "export-noinv-group", "export-noinv-stream"
	svc, store := newReadTestService(t, group)
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, group)); err != nil {
		t.Fatal(err)
	}
	putReadEvents(t, store, group, stream, time.Now().UnixMilli()-10, 3)

	taskId := "export-noinv-1"
	if err := store.PutExportTask(&logsstore.ExportTask{
		TaskId: taskId, LogGroupName: group, Status: "RUNNING",
		From: 0, To: time.Now().UnixMilli() + 1000,
		Destination: "export-bucket", CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}

	svc.executeExportTask(context.Background(), "us-east-1", group, "", 0, time.Now().UnixMilli()+1000, "export-bucket", "", taskId)

	task, err := store.GetExportTask(taskId)
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != "FAILED" {
		t.Fatalf("status = %q, want FAILED (message: %q)", task.Status, task.StatusMessage)
	}
}

// A task the canceller already moved to CANCELLED keeps that state: the
// worker skips the upload and the finaliser refuses to overwrite a
// terminal status.
// The cancel flow rides the status vocabulary's transitional state: the
// canceller writes PENDING_CANCEL (a client sees a live task, never a
// terminal status while the worker may still be mid-read), the worker
// honours the request at its next checkpoint and writes the terminal
// CANCELLED without delivering, and a repeat cancel is idempotent.
func TestExportCancelRidesPendingCancelTransition(t *testing.T) {
	const group, stream = "export-pending-cancel-group", "export-pending-cancel-stream"
	svc, store := newReadTestService(t, group)
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, group)); err != nil {
		t.Fatal(err)
	}
	putReadEvents(t, store, group, stream, time.Now().UnixMilli()-10, 3)

	s3 := newRecordingS3Invoker()
	bus := eventbus.NewEventBus()
	bus.SetS3Invoker(s3)
	busCtx, busCancel := context.WithCancel(context.Background())
	if err := bus.Start(busCtx); err != nil {
		t.Fatalf("bus start: %v", err)
	}
	t.Cleanup(busCancel)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatalf("SetEventBus: %v", err)
	}

	taskId := "export-pending-cancel-1"
	if err := store.PutExportTask(&logsstore.ExportTask{
		TaskId: taskId, LogGroupName: group, Status: "RUNNING",
		From: 0, To: time.Now().UnixMilli() + 1000,
		Destination: "export-bucket", CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}

	// The canceller's write is the transitional state, not the terminal
	// one: the worker has not run yet.
	if err := svc.cancelExportTaskCore(store, taskId); err != nil {
		t.Fatalf("cancel: %v", err)
	}
	task, err := store.GetExportTask(taskId)
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != logsstore.ExportStatusPendingCancel {
		t.Fatalf("after cancel: status = %q, want PENDING_CANCEL", task.Status)
	}
	// A repeat cancel keeps the transitional state (idempotent request).
	if err := svc.cancelExportTaskCore(store, taskId); err != nil {
		t.Fatalf("repeat cancel: %v", err)
	}
	if task, err = store.GetExportTask(taskId); err != nil || task.Status != logsstore.ExportStatusPendingCancel {
		t.Fatalf("after repeat cancel: status = %q err %v, want PENDING_CANCEL", task.Status, err)
	}

	// The worker honours the request at its checkpoint: terminal CANCELLED,
	// no delivered object.
	svc.executeExportTask(context.Background(), "us-east-1", group, "", 0, time.Now().UnixMilli()+1000, "export-bucket", "", taskId)
	task, err = store.GetExportTask(taskId)
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != logsstore.ExportStatusCancelled {
		t.Fatalf("after worker: status = %q, want CANCELLED", task.Status)
	}
	if _, written := s3.object("export-bucket", "exportedlogs/"+taskId+"/exportedlogs.gz"); written {
		t.Fatal("cancelled task delivered its export object")
	}

	// A terminal task refuses a further cancel.
	if err := svc.cancelExportTaskCore(store, taskId); logsErrorCode(err) != "InvalidOperationException" {
		t.Fatalf("cancel of a terminal task: code=%q err=%v", logsErrorCode(err), err)
	}
}

func TestExecuteExportTaskCancelledStaysCancelled(t *testing.T) {
	const group, stream = "export-cancel-group", "export-cancel-stream"
	svc, store := newReadTestService(t, group)
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, group)); err != nil {
		t.Fatal(err)
	}
	putReadEvents(t, store, group, stream, time.Now().UnixMilli()-10, 3)

	s3 := newRecordingS3Invoker()
	bus := eventbus.NewEventBus()
	bus.SetS3Invoker(s3)
	busCtx, busCancel := context.WithCancel(context.Background())
	if err := bus.Start(busCtx); err != nil {
		t.Fatalf("bus start: %v", err)
	}
	t.Cleanup(busCancel)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatalf("SetEventBus: %v", err)
	}

	taskId := "export-cancel-1"
	if err := store.PutExportTask(&logsstore.ExportTask{
		TaskId: taskId, LogGroupName: group, Status: "CANCELLED",
		StatusMessage: "Cancelled by user",
		From:          0, To: time.Now().UnixMilli() + 1000,
		Destination: "export-bucket", CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}

	svc.executeExportTask(context.Background(), "us-east-1", group, "", 0, time.Now().UnixMilli()+1000, "export-bucket", "", taskId)

	task, err := store.GetExportTask(taskId)
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != "CANCELLED" {
		t.Fatalf("status = %q, want CANCELLED to survive the worker", task.Status)
	}
	if _, written := s3.object("export-bucket", "exportedlogs/"+taskId+"/exportedlogs.gz"); written {
		t.Fatal("cancelled task delivered its export object")
	}
}

// The documented quota: one active (RUNNING or PENDING) export task per
// account; a second creation rejects with LimitExceededException until
// the first reaches a terminal status.
func TestCreateExportTaskSingleActiveLimit(t *testing.T) {
	const group = "export-limit-group"
	svc, store := newReadTestService(t, group)

	active := &logsstore.ExportTask{
		TaskId: "export-active-1", LogGroupName: group, Status: "RUNNING",
		From: 1, To: 2, Destination: "export-bucket", CreationTime: time.Now().UnixMilli(),
	}
	if err := store.PutExportTask(active); err != nil {
		t.Fatal(err)
	}

	// The window rides the group's lifetime: "You must specify a time
	// that is not earlier than when this log group was created", so the
	// epoch-based fixture values the older rows used would now reject.
	windowTo := int(time.Now().UnixMilli() + 1000)
	windowFrom := windowTo - 2000

	_, err := svc.createExportTaskCore(store, &CreateExportTaskInput{
		LogGroupName: group, Destination: "export-bucket",
		From: windowFrom, To: windowTo, Region: "us-east-1",
	})
	if err == nil || !strings.Contains(err.Error(), "LimitExceededException") {
		t.Fatalf("second active export: got %v, want LimitExceededException", err)
	}

	active.Status = "COMPLETED"
	if err := store.PutExportTask(active); err != nil {
		t.Fatal(err)
	}
	taskId, err := svc.createExportTaskCore(store, &CreateExportTaskInput{
		LogGroupName: group, Destination: "export-bucket",
		From: windowFrom, To: windowTo, Region: "us-east-1",
	})
	if err != nil {
		t.Fatalf("creation after completion: %v", err)
	}
	// The spawned worker runs under the service lifecycle; with no bus
	// wired in this test it reaches its terminal FAILED state (the
	// invoker guard), which proves the launch is live and finished
	// before the store is torn down.
	waitFor(t, 5*time.Second, "export worker to finish", func() bool {
		task, err := store.GetExportTask(taskId)
		return err == nil && isTerminalExportStatus(task.Status)
	})
}

// The export quota's registry row is "Active export task — Each
// supported Region: 1": the slot is the Region's, so an active export in
// another Region leaves this Region's open while the same Region's own
// active task still blocks.
func TestCreateExportTaskQuotaIsRegionScoped(t *testing.T) {
	const group = "export-region-scope-group"
	svc, store := newReadTestService(t, group)

	otherStore, err := svc.getLogsStoreByRegion("us-west-2")
	if err != nil {
		t.Fatalf("second region store: %v", err)
	}
	if err := otherStore.PutExportTask(&logsstore.ExportTask{
		TaskId: "export-other-region", LogGroupName: group, Status: "RUNNING",
		From: 1, To: 2, Destination: "export-bucket", CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}

	// The window rides the group's lifetime: "You must specify a time
	// that is not earlier than when this log group was created".
	windowTo := int(time.Now().UnixMilli() + 1000)
	windowFrom := windowTo - 2000

	taskId, err := svc.createExportTaskCore(store, &CreateExportTaskInput{
		LogGroupName: group, Destination: "export-bucket",
		From: windowFrom, To: windowTo, Region: "us-east-1",
	})
	if err != nil {
		t.Fatalf("creation with another Region's export active: %v", err)
	}
	waitFor(t, 5*time.Second, "export worker to finish", func() bool {
		task, err := store.GetExportTask(taskId)
		return err == nil && isTerminalExportStatus(task.Status)
	})

	if err := store.PutExportTask(&logsstore.ExportTask{
		TaskId: "export-same-region", LogGroupName: group, Status: "RUNNING",
		From: 1, To: 2, Destination: "export-bucket", CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.createExportTaskCore(store, &CreateExportTaskInput{
		LogGroupName: group, Destination: "export-bucket",
		From: windowFrom, To: windowTo, Region: "us-east-1",
	}); err == nil || !strings.Contains(err.Error(), "LimitExceededException") {
		t.Fatalf("second active export in the same Region: got %v, want LimitExceededException", err)
	}
}

// The documented account-level quota: no more than three active
// (IN_PROGRESS) imports at a time; the fourth creation rejects with
// ThrottlingException — the quota-shaped error the operation declares —
// and a slot freed by a terminal status admits the next creation.
func TestCreateImportTaskActiveImportQuota(t *testing.T) {
	const group = "import-quota-group"
	svc, store := newReadTestService(t, group)

	for i := 0; i < logsstore.MaxActiveImportTasks; i++ {
		if err := store.PutImportTask(&logsstore.ImportTask{
			ImportId:        fmt.Sprintf("import-active-%d", i),
			ImportSourceArn: "arn:aws:s3:::import-quota-bucket",
			LogGroupName:    group,
			ImportStatus:    "IN_PROGRESS",
			CreationTime:    time.Now().UnixMilli(),
		}); err != nil {
			t.Fatal(err)
		}
	}

	_, err := svc.createImportTaskCore(store, &CreateImportTaskInput{
		ImportSourceArn: "arn:aws:s3:::import-quota-bucket",
		ImportRoleArn:   "arn:aws:iam::000000000000:role/import",
		Region:          "us-east-1",
		AccountID:       "000000000000",
	})
	if err == nil || !strings.Contains(err.Error(), "ThrottlingException") {
		t.Fatalf("fourth active import: got %v, want ThrottlingException", err)
	}

	// A terminal status frees the slot; the next creation is admitted and
	// its worker runs to a terminal state (no S3 invoker is wired in this
	// test, so the import fails its objects), proving the launch finished
	// before the store is torn down.
	first, err := store.GetImportTask("import-active-0")
	if err != nil {
		t.Fatal(err)
	}
	first.ImportStatus = "COMPLETED"
	if err := store.PutImportTask(first); err != nil {
		t.Fatal(err)
	}

	task, err := svc.createImportTaskCore(store, &CreateImportTaskInput{
		ImportSourceArn: "arn:aws:s3:::import-quota-bucket",
		ImportRoleArn:   "arn:aws:iam::000000000000:role/import",
		Region:          "us-east-1",
		AccountID:       "000000000000",
	})
	if err != nil {
		t.Fatalf("creation after a slot freed: %v", err)
	}
	waitFor(t, 5*time.Second, "import worker to finish", func() bool {
		saved, err := store.GetImportTask(task.ImportId)
		return err == nil && isTerminalImportStatus(saved.ImportStatus)
	})
}

// A multi-object import keeps every object's events (one stream per
// source object), applies the importFilter window, records the batches
// the DescribeImportTaskBatches shape serves, and reports bytesImported.
func TestExecuteImportTaskMultiObjectCompleteness(t *testing.T) {
	const group = "import-multi-group"
	svc, store := newReadTestService(t, group)

	hour := int64(1721237155000 / (60 * 60 * 1000) * (60 * 60 * 1000))
	objectA := fmt.Sprintf("{\"timestamp\": %d, \"message\": \"a1\"}\n{\"timestamp\": %d, \"message\": \"a2\"}\n", hour+1000, hour+2000)
	objectB := fmt.Sprintf("{\"timestamp\": %d, \"message\": \"b1\"}\n", hour+3600000) // next hour bucket
	var gz bytes.Buffer
	zw := gzip.NewWriter(&gz)
	zw.Write([]byte(objectA))
	zw.Close()

	s3 := newRecordingS3Invoker()
	if err := s3.PutObject(context.Background(), "us-east-1", "import-bucket", "obj/a.log", gz.Bytes(), "application/x-gzip"); err != nil {
		t.Fatal(err)
	}
	if err := s3.PutObject(context.Background(), "us-east-1", "import-bucket", "obj/b.log", []byte(objectB), "text/plain"); err != nil {
		t.Fatal(err)
	}
	bus := eventbus.NewEventBus()
	bus.SetS3Invoker(s3)
	busCtx, busCancel := context.WithCancel(context.Background())
	if err := bus.Start(busCtx); err != nil {
		t.Fatalf("bus start: %v", err)
	}
	t.Cleanup(busCancel)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatalf("SetEventBus: %v", err)
	}

	importId := "import-multi-1"
	if err := store.PutImportTask(&logsstore.ImportTask{
		ImportId: importId, ImportSourceArn: "arn:aws:s3:::import-bucket",
		LogGroupName: group, ImportStatus: "IN_PROGRESS",
		ImportFilter: map[string]interface{}{
			"startEventTime": float64(hour),
			"endEventTime":   float64(hour + 7200000),
		},
		CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}

	svc.executeImportTask(context.Background(), "us-east-1", importId, "import-bucket", "", group)

	task, err := store.GetImportTask(importId)
	if err != nil {
		t.Fatal(err)
	}
	if task.ImportStatus != "COMPLETED" {
		t.Fatalf("status = %q (%s), want COMPLETED", task.ImportStatus, task.ErrorMessage)
	}
	if task.ImportStatistics == nil {
		t.Fatalf("statistics missing, want bytesImported > 0")
	}
	if imported, ok := toInt64(task.ImportStatistics["bytesImported"]); !ok || imported <= 0 {
		t.Fatalf("statistics = %v, want bytesImported > 0", task.ImportStatistics)
	}

	// Every object's events survive: one stream per source key.
	eventsA, _, _, err := store.GetLogEvents(group, "obj/a.log", 0, 0, 0, true, "")
	if err != nil || len(eventsA) != 2 {
		t.Fatalf("obj/a.log events = %d (err %v), want 2 (the gzipped object)", len(eventsA), err)
	}
	eventsB, _, _, err := store.GetLogEvents(group, "obj/b.log", 0, 0, 0, true, "")
	if err != nil || len(eventsB) != 1 {
		t.Fatalf("obj/b.log events = %d (err %v), want 1", len(eventsB), err)
	}

	// Batches group by event-time hour and serve through the core with
	// the status filter and pagination the operation documents.
	if len(task.ImportBatches) != 2 {
		t.Fatalf("batches = %+v, want two hour buckets", task.ImportBatches)
	}
	if task.ImportBatches[0].HourStartMs >= task.ImportBatches[1].HourStartMs {
		t.Fatalf("batches not ordered by bucket: %+v", task.ImportBatches)
	}
	for _, b := range task.ImportBatches {
		if b.Status != "COMPLETED" {
			t.Fatalf("batch %s status = %q, want COMPLETED", b.BatchId, b.Status)
		}
	}
	page, err := svc.describeImportTaskBatchesCore(store, &DescribeImportTaskBatchesInput{
		ImportId: importId, BatchImportStatus: []string{"COMPLETED"}, Limit: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.ImportBatches) != 1 || page.NextToken == "" {
		t.Fatalf("filtered batch page = %+v, token %q", page.ImportBatches, page.NextToken)
	}
	// The continuation request repeats the filter member: the scoped
	// marker keys on the request identity, so a filtered walk's second
	// page presents the same filter (an unfiltered continuation is a
	// different listing and its marker rejects).
	page2, err := svc.describeImportTaskBatchesCore(store, &DescribeImportTaskBatchesInput{
		ImportId: importId, BatchImportStatus: []string{"COMPLETED"},
		NextToken: page.NextToken, Limit: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(page2.ImportBatches) != 1 || page2.NextToken != "" {
		t.Fatalf("second batch page = %+v, token %q", page2.ImportBatches, page2.NextToken)
	}
}

// The importFilter window excludes out-of-window events instead of
// importing everything.
func TestExecuteImportTaskFilterWindowHonoured(t *testing.T) {
	const group = "import-window-group"
	svc, store := newReadTestService(t, group)

	inWindow := int64(1721237155000)
	outBefore := inWindow - 7200000
	outAfter := inWindow + 7200000
	content := fmt.Sprintf("{\"timestamp\": %d, \"message\": \"before\"}\n{\"timestamp\": %d, \"message\": \"inside\"}\n{\"timestamp\": %d, \"message\": \"after\"}\n",
		outBefore, inWindow, outAfter)

	s3 := newRecordingS3Invoker()
	if err := s3.PutObject(context.Background(), "us-east-1", "import-bucket", "win.log", []byte(content), "text/plain"); err != nil {
		t.Fatal(err)
	}
	bus := eventbus.NewEventBus()
	bus.SetS3Invoker(s3)
	busCtx, busCancel := context.WithCancel(context.Background())
	if err := bus.Start(busCtx); err != nil {
		t.Fatalf("bus start: %v", err)
	}
	t.Cleanup(busCancel)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatalf("SetEventBus: %v", err)
	}

	importId := "import-window-1"
	if err := store.PutImportTask(&logsstore.ImportTask{
		ImportId: importId, ImportSourceArn: "arn:aws:s3:::import-bucket/win.log",
		LogGroupName: group, ImportStatus: "IN_PROGRESS",
		ImportFilter: map[string]interface{}{
			"startEventTime": float64(inWindow - 1000),
			"endEventTime":   float64(inWindow + 1000),
		},
		CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}

	svc.executeImportTask(context.Background(), "us-east-1", importId, "import-bucket", "win.log", group)

	events, _, _, err := store.GetLogEvents(group, "win.log", 0, 0, 0, true, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].Message != "inside" {
		t.Fatalf("imported events = %+v, want only the in-window event", events)
	}
}

// An import where every object failed reports FAILED with the failure
// count, not COMPLETED.
func TestExecuteImportTaskAllObjectsFailed(t *testing.T) {
	const group = "import-fail-group"
	svc, store := newReadTestService(t, group)

	s3 := newRecordingS3Invoker()
	bus := eventbus.NewEventBus()
	bus.SetS3Invoker(s3)
	busCtx, busCancel := context.WithCancel(context.Background())
	if err := bus.Start(busCtx); err != nil {
		t.Fatalf("bus start: %v", err)
	}
	t.Cleanup(busCancel)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatalf("SetEventBus: %v", err)
	}

	importId := "import-fail-1"
	if err := store.PutImportTask(&logsstore.ImportTask{
		ImportId: importId, ImportSourceArn: "arn:aws:s3:::import-bucket/missing.log",
		LogGroupName: group, ImportStatus: "IN_PROGRESS",
		CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}

	// The single named object does not exist: GetObject fails for it.
	svc.executeImportTask(context.Background(), "us-east-1", importId, "import-bucket", "missing.log", group)

	task, err := store.GetImportTask(importId)
	if err != nil {
		t.Fatal(err)
	}
	if task.ImportStatus != "FAILED" {
		t.Fatalf("status = %q, want FAILED for an all-objects-failed import", task.ImportStatus)
	}
	if !strings.Contains(task.ErrorMessage, "failed") {
		t.Fatalf("errorMessage = %q, want the failure accounting", task.ErrorMessage)
	}
}

// A restart fails export and import records an interrupted process left
// in an active state, so no RUNNING task survives without a worker.
func TestReconcileInterruptedTasks(t *testing.T) {
	const group = "reconcile-group"
	svc, store := newReadTestService(t, group)

	if err := store.PutExportTask(&logsstore.ExportTask{
		TaskId: "export-orphan-1", LogGroupName: group, Status: "RUNNING",
		From: 1, To: 2, Destination: "b", CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.PutExportTask(&logsstore.ExportTask{
		TaskId: "export-done-1", LogGroupName: group, Status: "COMPLETED",
		From: 1, To: 2, Destination: "b", CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.PutImportTask(&logsstore.ImportTask{
		ImportId: "import-orphan-1", ImportSourceArn: "arn:aws:s3:::b/k",
		LogGroupName: group, ImportStatus: "IN_PROGRESS",
		CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}

	svc.reconcileInterruptedTasks()

	export, err := store.GetExportTask("export-orphan-1")
	if err != nil {
		t.Fatal(err)
	}
	if export.Status != "FAILED" || !strings.Contains(export.StatusMessage, "interrupted") {
		t.Fatalf("orphaned export = %q/%q, want FAILED with the interruption reason", export.Status, export.StatusMessage)
	}
	done, err := store.GetExportTask("export-done-1")
	if err != nil {
		t.Fatal(err)
	}
	if done.Status != "COMPLETED" {
		t.Fatalf("completed export touched by reconciliation: %q", done.Status)
	}
	imp, err := store.GetImportTask("import-orphan-1")
	if err != nil {
		t.Fatal(err)
	}
	if imp.ImportStatus != "FAILED" || !strings.Contains(imp.ErrorMessage, "interrupted") {
		t.Fatalf("orphaned import = %q/%q, want FAILED with the interruption reason", imp.ImportStatus, imp.ErrorMessage)
	}
}

// Import cancellation and worker finalisation are mutually exclusive
// terminal writers: both run through MutateImportTask's critical section
// with a terminal gate, so at most ONE racer can flip a non-terminal
// task to its terminal status — the unlocked read-then-put canceller
// could observe IN_PROGRESS, sleep past the worker's COMPLETED write,
// and revert it from a stale snapshot.
func TestImportCancelAndFinaliseAdmitOneTerminalWriter(t *testing.T) {
	svc, store := newReadTestService(t, "import-gate-group")
	if err := store.PutImportTask(&logsstore.ImportTask{
		ImportId:     "import-gated",
		LogGroupName: "import-gate-group",
		ImportStatus: logsstore.ImportStatusInProgress,
	}); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	var flips int32
	start := make(chan struct{})
	cancel := func() {
		defer wg.Done()
		<-start
		if _, err := svc.cancelImportTaskCore(store, "import-gated"); err == nil {
			atomic.AddInt32(&flips, 1)
		}
	}
	finalise := func() {
		defer wg.Done()
		<-start
		_ = store.MutateImportTask("import-gated", func(task *logsstore.ImportTask) error {
			if task.ImportStatus != logsstore.ImportStatusInProgress {
				return nil
			}
			task.ImportStatus = logsstore.ImportStatusCompleted
			task.ImportStatistics = map[string]interface{}{"importedLogEventCount": 7}
			atomic.AddInt32(&flips, 1)
			return nil
		})
	}
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go cancel()
		wg.Add(1)
		go finalise()
	}
	close(start)
	wg.Wait()

	if flips > 1 {
		t.Fatalf("terminal admission must admit at most one writer, saw %d", flips)
	}
	task, err := store.GetImportTask("import-gated")
	if err != nil {
		t.Fatal(err)
	}
	switch task.ImportStatus {
	case logsstore.ImportStatusCancelled:
	case logsstore.ImportStatusCompleted:
		// The store round-trips the statistics map through JSON, so the
		// count deserialises as float64.
		got, ok := task.ImportStatistics["importedLogEventCount"].(float64)
		if !ok || got != 7 {
			t.Fatalf("the finaliser's statistics must survive, got %v", task.ImportStatistics["importedLogEventCount"])
		}
	default:
		t.Fatalf("the task must reach a terminal status, got %s", task.ImportStatus)
	}
}

// The managed destination group is named after the source bucket: a
// bucket-only ARN carries no slash and must not leak the whole ARN —
// colons and all — into a name the log-group charset forbids, and a
// keyed ARN of the same bucket shares the one managed group.
// The documented CloudTrail Lake source form is accepted: the eventdatastore
// ARN creates the managed group under the store's id, a foreign cloudtrail
// resource or a foreign service rejects, and the executor walks the store
// through the cross-service invoker — every record's own payload lands as
// one log event, batched by event-time hour.
func TestCreateImportTaskEventDataStoreSource(t *testing.T) {
	svc, store := newReadTestService(t, "import-eds-group")
	newInput := func(source string) *CreateImportTaskInput {
		return &CreateImportTaskInput{
			ImportSourceArn: source,
			ImportRoleArn:   "arn:aws:iam::000000000000:role/import",
			Region:          "us-east-1",
			AccountID:       "000000000000",
		}
	}

	// The create-path rejection rows: a foreign cloudtrail resource, a
	// missing store id and a foreign service are parameter errors.
	for name, source := range map[string]string{
		"foreign cloudtrail resource": "arn:aws:cloudtrail:us-east-1:000000000000:trail/my-trail",
		"bare eventdatastore prefix":  "arn:aws:cloudtrail:us-east-1:000000000000:eventdatastore/",
		"foreign service":             "arn:aws:dynamodb:us-east-1:000000000000:table/t",
		"not an ARN":                  "some-bucket",
	} {
		_, err := svc.createImportTaskCore(store, newInput(source))
		if code := logsErrorCode(err); code != "InvalidParameterException" {
			t.Fatalf("%s: code=%q err=%v", name, code, err)
		}
	}

	// The documented form creates: the managed group names the store's id.
	task, err := svc.createImportTaskCore(store, newInput("arn:aws:cloudtrail:us-east-1:000000000000:eventdatastore/6c5af0ea-1234-4f2a-9d3b-111111111111"))
	if err != nil {
		t.Fatal(err)
	}
	wantGroup := "/aws/imported/cloudtrail-lake/6c5af0ea-1234-4f2a-9d3b-111111111111"
	if task.LogGroupName != wantGroup {
		t.Fatalf("managed group name = %q, want %q", task.LogGroupName, wantGroup)
	}
	if _, err := store.GetLogGroup(wantGroup); err != nil {
		t.Fatalf("managed group was not created: %v", err)
	}
}

// fakeEDSInvoker serves the event data store walk from an in-memory page
// set, recording the bounds it was asked for.
type fakeEDSInvoker struct {
	events  []invokers.CloudTrailEDSEvent
	start   *time.Time
	end     *time.Time
	pageFmt string // format verb to render the walk's cursor
}

func (f *fakeEDSInvoker) LookupEvents(ctx context.Context, region, username, nextToken string, startTime, endTime time.Time, maxResults int32) ([]invokers.CloudTrailEventInfo, string, error) {
	return nil, "", nil
}

func (f *fakeEDSInvoker) LookupEDSEvents(ctx context.Context, edsArn string, start, end *time.Time, maxResults int, nextToken string) ([]invokers.CloudTrailEDSEvent, string, error) {
	if f.start == nil {
		f.start, f.end = start, end
	}
	begin := 0
	if nextToken != "" {
		fmt.Sscanf(nextToken, f.pageFmt, &begin)
	}
	if begin >= len(f.events) {
		return nil, "", nil
	}
	stop := begin + maxResults
	if stop > len(f.events) {
		stop = len(f.events)
	}
	token := ""
	if stop < len(f.events) {
		token = fmt.Sprintf(f.pageFmt, stop)
	}
	return f.events[begin:stop], token, nil
}

func TestExecuteImportTaskFromEDSWalksStore(t *testing.T) {
	const group = "import-eds-walk-group"
	svc, store := newReadTestService(t, group)

	hour := int64(1721237155000 / (60 * 60 * 1000) * (60 * 60 * 1000))
	fake := &fakeEDSInvoker{
		pageFmt: "%d",
		events: []invokers.CloudTrailEDSEvent{
			{EventTime: time.UnixMilli(hour + 1000), Payload: `{"eventID":"e1"}`},
			{EventTime: time.UnixMilli(hour + 2000), Payload: `{"eventID":"e2"}`},
			{EventTime: time.UnixMilli(hour + 3600000), Payload: `{"eventID":"e3"}`},
		},
	}
	bus := eventbus.NewEventBus()
	bus.SetCloudTrailInvoker(fake)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatalf("SetEventBus: %v", err)
	}

	importId := "import-eds-1"
	if err := store.PutImportTask(&logsstore.ImportTask{
		ImportId: importId, ImportSourceArn: "arn:aws:cloudtrail:us-east-1:000000000000:eventdatastore/eds-1",
		LogGroupName: group, ImportStatus: "IN_PROGRESS",
		ImportFilter: map[string]interface{}{
			"startEventTime": float64(hour),
			"endEventTime":   float64(hour + 7200000),
		},
		CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}

	svc.executeImportTaskFromEDS(context.Background(), "us-east-1", importId,
		"arn:aws:cloudtrail:us-east-1:000000000000:eventdatastore/eds-1", group)

	task, err := store.GetImportTask(importId)
	if err != nil {
		t.Fatal(err)
	}
	if task.ImportStatus != "COMPLETED" {
		t.Fatalf("status = %q (%s), want COMPLETED", task.ImportStatus, task.ErrorMessage)
	}
	// The walk carried the task's filter window as its bounds.
	if fake.start == nil || fake.start.UnixMilli() != hour {
		t.Fatalf("walk start bound = %v, want the filter's startEventTime", fake.start)
	}
	if fake.end == nil || fake.end.UnixMilli() != hour+7200000 {
		t.Fatalf("walk end bound = %v, want the filter's endEventTime", fake.end)
	}
	// Every record landed in the one stream the walk owns, payloads intact.
	streamName := importStreamName("eventdatastore/eds-1")
	events, _, _, err := store.GetLogEvents(group, streamName, 0, 0, 0, true, "")
	if err != nil || len(events) != 3 {
		t.Fatalf("walked events = %d (err %v), want 3", len(events), err)
	}
	if events[0].Message != `{"eventID":"e1"}` || events[2].Message != `{"eventID":"e3"}` {
		t.Fatalf("event payloads = %+v, want the records' own JSON", events)
	}
	// The hour buckets the events landed in report COMPLETED.
	if len(task.ImportBatches) != 2 {
		t.Fatalf("batches = %+v, want two hour buckets", task.ImportBatches)
	}
	for _, b := range task.ImportBatches {
		if b.Status != "COMPLETED" {
			t.Fatalf("batch %s status = %q, want COMPLETED", b.BatchId, b.Status)
		}
	}
	if task.ImportStatistics == nil {
		t.Fatalf("statistics missing")
	}
	if imported, ok := toInt64(task.ImportStatistics["bytesImported"]); !ok || imported != int64(len(`{"eventID":"e1"}`)*3) {
		t.Fatalf("bytesImported = %v, want the summed payload bytes", task.ImportStatistics)
	}
}

func TestImportManagedGroupNameDerivesFromBucket(t *testing.T) {
	svc, store := newReadTestService(t, "import-name-group")
	new := func(source string) *CreateImportTaskInput {
		return &CreateImportTaskInput{
			ImportSourceArn: source,
			ImportRoleArn:   "arn:aws:iam::000000000000:role/import",
			Region:          "us-east-1",
			AccountID:       "000000000000",
		}
	}
	task, err := svc.createImportTaskCore(store, new("arn:aws:s3:::my-import-bucket"))
	if err != nil {
		t.Fatal(err)
	}
	if task.LogGroupName != "/aws/imported/cloudtrail-lake/my-import-bucket" {
		t.Fatalf("managed group name = %q", task.LogGroupName)
	}
	if strings.Contains(task.LogGroupName, ":") {
		t.Fatalf("managed group name left the log-group charset: %q", task.LogGroupName)
	}
	if _, err := store.GetLogGroup(task.LogGroupName); err != nil {
		t.Fatalf("managed group was not created: %v", err)
	}
	waitFor(t, 5*time.Second, "import worker to finish", func() bool {
		saved, err := store.GetImportTask(task.ImportId)
		return err == nil && isTerminalImportStatus(saved.ImportStatus)
	})
}

// "Specifying a task ID filters the results to one or zero export
// tasks" — an unknown taskId serves an empty list, not the
// undeclared ResourceNotFoundException.
func TestDescribeExportTasksUnknownTaskIdIsEmpty(t *testing.T) {
	svc, store := newReadTestService(t, "dxt-unknown-group")
	tasks, next, err := svc.describeExportTasksCore(store, &DescribeExportTasksInput{TaskId: "no-such-export-task"})
	if err != nil {
		t.Fatalf("unknown taskId must not error: %v", err)
	}
	if len(tasks) != 0 || next != "" {
		t.Fatalf("unknown taskId served %+v (token %q), want an empty list", tasks, next)
	}

	known := "dxt-known-1"
	if err := store.PutExportTask(&logsstore.ExportTask{
		TaskId: known, LogGroupName: "dxt-unknown-group", Status: "COMPLETED",
		Destination: "some-bucket", CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}
	tasks, _, err = svc.describeExportTasksCore(store, &DescribeExportTasksInput{TaskId: known})
	if err != nil || len(tasks) != 1 || tasks[0].TaskId != known {
		t.Fatalf("known taskId served %+v (err %v)", tasks, err)
	}
}

// destinationPrefix: "The prefix used as the start of the key for every
// object exported. If you don't specify a value, the default is
// exportedlogs." — an explicit prefix lands under it verbatim.
func TestExportDestinationPrefixes(t *testing.T) {
	const group, stream = "export-prefix-group", "export-prefix-stream"
	svc, store := newReadTestService(t, group)
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, group)); err != nil {
		t.Fatal(err)
	}
	putReadEvents(t, store, group, stream, time.Now().UnixMilli()-10, 3)

	s3 := newRecordingS3Invoker()
	bus := eventbus.NewEventBus()
	bus.SetS3Invoker(s3)
	busCtx, busCancel := context.WithCancel(context.Background())
	if err := bus.Start(busCtx); err != nil {
		t.Fatalf("bus start: %v", err)
	}
	t.Cleanup(busCancel)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatalf("SetEventBus: %v", err)
	}

	taskId := "export-prefix-1"
	if err := store.PutExportTask(&logsstore.ExportTask{
		TaskId: taskId, LogGroupName: group, Status: "RUNNING",
		From: 0, To: time.Now().UnixMilli() + 1000,
		Destination: "export-bucket", DestinationPrefix: "custom-prefix",
		CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}

	svc.executeExportTask(context.Background(), "us-east-1", group, "", 0, time.Now().UnixMilli()+1000, "export-bucket", "custom-prefix", taskId)

	if task, err := store.GetExportTask(taskId); err != nil || task.Status != "COMPLETED" {
		t.Fatalf("status = %+v (err %v)", task, err)
	}
	if _, ok := s3.object("export-bucket", "custom-prefix/"+taskId+"/exportedlogs.gz"); !ok {
		t.Fatal("explicit-prefix export object not written")
	}
	if _, ok := s3.object("export-bucket", "exportedlogs/"+taskId+"/exportedlogs.gz"); ok {
		t.Fatal("default prefix leaked into an explicit-prefix export")
	}
}

// DescribeImportTaskBatches validates its members: limit "Valid Range:
// Minimum value of 1. Maximum value of 50. Default: 10" and
// batchImportStatus "Valid Values: IN_PROGRESS | CANCELLED | COMPLETED
// | FAILED".
func TestDescribeImportTaskBatchesValidation(t *testing.T) {
	svc, store := newReadTestService(t, "ditb-val-group")
	if err := store.PutImportTask(&logsstore.ImportTask{
		ImportId: "ditb-val-1", ImportSourceArn: "arn:aws:s3:::some-bucket",
		LogGroupName: "ditb-val-group", ImportStatus: "COMPLETED",
		CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.describeImportTaskBatchesCore(store, &DescribeImportTaskBatchesInput{
		ImportId: "ditb-val-1", Limit: 51,
	}); err == nil || logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("limit 51: %v", err)
	}
	if _, err := svc.describeImportTaskBatchesCore(store, &DescribeImportTaskBatchesInput{
		ImportId: "ditb-val-1", BatchImportStatus: []string{"BOGUS"},
	}); err == nil || logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("bogus batchImportStatus: %v", err)
	}
	// The default limit applies when unset, and the legal boundary
	// serves.
	if _, err := svc.describeImportTaskBatchesCore(store, &DescribeImportTaskBatchesInput{
		ImportId: "ditb-val-1",
	}); err != nil {
		t.Fatalf("default limit: %v", err)
	}
	if _, err := svc.describeImportTaskBatchesCore(store, &DescribeImportTaskBatchesInput{
		ImportId: "ditb-val-1", Limit: 50,
	}); err != nil {
		t.Fatalf("limit 50: %v", err)
	}
}

// The import batch lifecycle: batches persist IN_PROGRESS while the
// task runs (a mid-flight cancellation observes progress statistics), a
// bucket whose write failed stamps FAILED with its errorMessage, and
// the untouched buckets follow the task's terminal status.
func TestImportBatchLifecycle(t *testing.T) {
	const group = "import-batch-group"
	svc, store := newReadTestService(t, group)

	importId := "import-batch-1"
	if err := store.PutImportTask(&logsstore.ImportTask{
		ImportId: importId, ImportSourceArn: "arn:aws:s3:::import-bucket",
		LogGroupName: group, ImportStatus: "IN_PROGRESS",
		CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}

	hour := int64(1721237155000 / (60 * 60 * 1000) * (60 * 60 * 1000))
	batchIndex := map[int64]*logsstore.ImportBatch{
		hour: {BatchId: "batch-good", Status: "IN_PROGRESS", HourStartMs: hour},
	}
	// Progress persistence: the in-flight task serves live batches and
	// statistics a cancellation would report.
	svc.persistImportProgress("us-east-1", importId, 1234, batchIndex)
	live, err := store.GetImportTask(importId)
	if err != nil {
		t.Fatal(err)
	}
	if live.ImportStatus != "IN_PROGRESS" || len(live.ImportBatches) != 1 || live.ImportBatches[0].Status != "IN_PROGRESS" {
		t.Fatalf("in-flight batches = %+v, status %q", live.ImportBatches, live.ImportStatus)
	}
	if v, ok := toInt64(live.ImportStatistics["bytesImported"]); !ok || v != 1234 {
		t.Fatalf("in-flight statistics = %v", live.ImportStatistics)
	}

	// Finalisation: the failed bucket carries FAILED + errorMessage,
	// the good bucket follows the task's terminal status.
	nextHour := hour + 3600000
	batchIndex[nextHour] = &logsstore.ImportBatch{BatchId: "batch-bad", Status: "IN_PROGRESS", HourStartMs: nextHour}
	svc.finalizeImportTask("us-east-1", importId, logsstore.ImportStatusCompleted,
		"1 of 2 source objects failed to import", 2345, batchIndex,
		map[int64]string{nextHour: "stream s write error: rejected"})
	final, err := store.GetImportTask(importId)
	if err != nil {
		t.Fatal(err)
	}
	if final.ImportStatus != "COMPLETED" || final.ErrorMessage == "" {
		t.Fatalf("task = %+v", final)
	}
	byHour := map[string]*logsstore.ImportBatch{}
	for _, b := range final.ImportBatches {
		byHour[b.BatchId] = b
	}
	if byHour["batch-bad"].Status != "FAILED" || byHour["batch-bad"].ErrorMessage == "" {
		t.Fatalf("failed bucket = %+v", byHour["batch-bad"])
	}
	if byHour["batch-good"].Status != "COMPLETED" || byHour["batch-good"].ErrorMessage != "" {
		t.Fatalf("good bucket = %+v", byHour["batch-good"])
	}
}

// CreateExportTask's member traits: taskName 1-512 ("Length Constraints:
// Minimum length of 1. Maximum length of 512.", CreateExportTask API
// reference), the logStreamNamePrefix stream-name traits, from/to range
// min 0 (from == 0 is a valid window start), and "You must specify a time
// that is not earlier than when this log group was created."
func TestCreateExportTaskMemberTraitsAndWindowRules(t *testing.T) {
	const group = "export-traits-group"
	svc, store := newReadTestService(t, group)

	base := &CreateExportTaskInput{
		LogGroupName: group, Destination: "export-bucket",
		From: 0, To: int(time.Now().UnixMilli() + 1000), Region: "us-east-1",
	}
	with := func(mut func(in *CreateExportTaskInput)) *CreateExportTaskInput {
		in := *base
		mut(&in)
		return &in
	}

	cases := []struct {
		name string
		in   *CreateExportTaskInput
	}{
		{"taskName over 512", with(func(in *CreateExportTaskInput) {
			in.TaskName = strings.Repeat("n", logsstore.MaxExportTaskNameLength+1)
		})},
		{"prefix with colon", with(func(in *CreateExportTaskInput) { in.LogStreamNamePrefix = "bad:prefix" })},
		{"prefix with asterisk", with(func(in *CreateExportTaskInput) { in.LogStreamNamePrefix = "bad*" })},
		{"negative from", with(func(in *CreateExportTaskInput) { in.From = -1 })},
		{"negative to", with(func(in *CreateExportTaskInput) { in.To = -1 })},
		{"to before group creation", with(func(in *CreateExportTaskInput) { in.From = 1; in.To = 2 })},
	}
	for _, tc := range cases {
		_, err := svc.createExportTaskCore(store, tc.in)
		if code := logsErrorCode(err); code != "InvalidParameterException" {
			t.Fatalf("%s: code=%q want InvalidParameterException (err=%v)", tc.name, code, err)
		}
	}

	// from == 0 (the epoch start) sits inside the Timestamp range: the
	// task is admitted and its worker reaches a terminal state (no S3
	// invoker is wired in this test, so the invoker guard fails it).
	taskId, err := svc.createExportTaskCore(store, base)
	if err != nil {
		t.Fatalf("epoch-start window rejected: %v", err)
	}
	waitFor(t, 5*time.Second, "export worker to finish", func() bool {
		task, err := store.GetExportTask(taskId)
		return err == nil && isTerminalExportStatus(task.Status)
	})
}

// The taskId and statusCode filters compose: "Specifying a task ID
// filters the results to one or zero export tasks" — the unknown id and
// the status mismatch are the same zero side of that filter.
func TestDescribeExportTasksTaskIdHonoursStatusCode(t *testing.T) {
	const group = "export-filter-group"
	svc, store := newReadTestService(t, group)
	if err := store.PutExportTask(&logsstore.ExportTask{
		TaskId: "export-filter-1", LogGroupName: group, Status: "RUNNING",
		From: 1, To: 2, Destination: "b", CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}

	tasks, _, err := svc.describeExportTasksCore(store, &DescribeExportTasksInput{TaskId: "export-filter-1"})
	if err != nil || len(tasks) != 1 {
		t.Fatalf("taskId alone: %v (%d tasks)", err, len(tasks))
	}
	tasks, _, err = svc.describeExportTasksCore(store, &DescribeExportTasksInput{TaskId: "export-filter-1", StatusCode: "COMPLETED"})
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 0 {
		t.Fatalf("status-mismatched taskId: %d tasks, want 0", len(tasks))
	}
	tasks, _, err = svc.describeExportTasksCore(store, &DescribeExportTasksInput{TaskId: "export-filter-1", StatusCode: "RUNNING"})
	if err != nil || len(tasks) != 1 {
		t.Fatalf("status-matched taskId: %v (%d tasks)", err, len(tasks))
	}
	tasks, _, err = svc.describeExportTasksCore(store, &DescribeExportTasksInput{TaskId: "never-was"})
	if err != nil || len(tasks) != 0 {
		t.Fatalf("unknown taskId: %v (%d tasks)", err, len(tasks))
	}
}

// "Export tasks time out after 24 hours." (Exporting log data to Amazon
// S3) — a RUNNING task past the deadline fails on the next read or
// creation census instead of reporting RUNNING forever or holding the
// account's single active slot.
func TestExportTasksTimeOutAfter24Hours(t *testing.T) {
	const group = "export-timeout-group"
	svc, store := newReadTestService(t, group)

	if err := store.PutExportTask(&logsstore.ExportTask{
		TaskId: "export-stale-1", LogGroupName: group, Status: "RUNNING",
		From: 1, To: 2, Destination: "b",
		CreationTime: time.Now().Add(-25 * time.Hour).UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}

	tasks, _, err := svc.describeExportTasksCore(store, &DescribeExportTasksInput{TaskId: "export-stale-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 || tasks[0].Status != logsstore.ExportStatusFailed {
		t.Fatalf("timed-out task after describe: %+v", tasks)
	}
	if tasks[0].StatusMessage != exportTaskTimeoutMessage {
		t.Fatalf("timeout message: %q", tasks[0].StatusMessage)
	}

	// The expired task no longer holds the account's single active slot.
	taskId, err := svc.createExportTaskCore(store, &CreateExportTaskInput{
		LogGroupName: group, Destination: "export-bucket",
		From: 0, To: int(time.Now().UnixMilli() + 1000), Region: "us-east-1",
	})
	if err != nil {
		t.Fatalf("creation after the timeout sweep: %v", err)
	}
	waitFor(t, 5*time.Second, "export worker to finish", func() bool {
		task, err := store.GetExportTask(taskId)
		return err == nil && isTerminalExportStatus(task.Status)
	})
}

// The executor's own deadline: past the runtime, the read legs stop and
// the task fails with the timeout message instead of running without
// bound.
func TestExecuteExportTaskDeadline(t *testing.T) {
	const group, stream = "export-deadline-group", "export-deadline-stream"
	svc, store := newReadTestService(t, group)
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, group)); err != nil {
		t.Fatal(err)
	}
	putReadEvents(t, store, group, stream, time.Now().UnixMilli()-10, 3)

	s3 := newRecordingS3Invoker()
	bus := eventbus.NewEventBus()
	bus.SetS3Invoker(s3)
	busCtx, busCancel := context.WithCancel(context.Background())
	if err := bus.Start(busCtx); err != nil {
		t.Fatalf("bus start: %v", err)
	}
	t.Cleanup(busCancel)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatalf("SetEventBus: %v", err)
	}

	if err := store.PutExportTask(&logsstore.ExportTask{
		TaskId: "export-deadline-1", LogGroupName: group, Status: "RUNNING",
		From: 0, To: time.Now().UnixMilli() + 1000,
		Destination: "export-bucket", CreationTime: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}

	orig := exportTaskRuntime
	exportTaskRuntime = time.Nanosecond
	t.Cleanup(func() { exportTaskRuntime = orig })

	svc.executeExportTask(context.Background(), "us-east-1", group, "", 0, time.Now().UnixMilli()+1000, "export-bucket", "", "export-deadline-1")

	task, err := store.GetExportTask("export-deadline-1")
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != logsstore.ExportStatusFailed || task.StatusMessage != exportTaskTimeoutMessage {
		t.Fatalf("deadline task: %q / %q", task.Status, task.StatusMessage)
	}
}

// The importFilter window members are Timestamps (range min 0): a
// non-numeric or negative value rejects instead of silently becoming the
// unbounded window the executor would treat it as. ImportId carries its
// 1-256 [a-zA-Z0-9-] traits on every operation that addresses an import:
// the malformed id is a parameter error, the well-formed unknown id keeps
// its not-found / empty-list identity.
func TestImportFilterWindowAndImportIdValidation(t *testing.T) {
	const group = "import-traits-group"
	svc, store := newReadTestService(t, group)
	const role = "arn:aws:iam::000000000000:role/import"

	_, err := svc.createImportTaskCore(store, &CreateImportTaskInput{
		ImportSourceArn: "arn:aws:s3:::import-traits-bucket/import.log", ImportRoleArn: role,
		ImportFilter: map[string]interface{}{"startEventTime": "not-a-number"},
		Region:       "us-east-1", AccountID: "000000000000",
	})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("non-numeric startEventTime: code=%q err=%v", code, err)
	}
	_, err = svc.createImportTaskCore(store, &CreateImportTaskInput{
		ImportSourceArn: "arn:aws:s3:::import-traits-bucket/import.log", ImportRoleArn: role,
		ImportFilter: map[string]interface{}{"endEventTime": float64(-5)},
		Region:       "us-east-1", AccountID: "000000000000",
	})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("negative endEventTime: code=%q err=%v", code, err)
	}

	task, err := svc.createImportTaskCore(store, &CreateImportTaskInput{
		ImportSourceArn: "arn:aws:s3:::import-traits-bucket/import.log", ImportRoleArn: role,
		ImportFilter: map[string]interface{}{"startEventTime": float64(0), "endEventTime": float64(time.Now().UnixMilli())},
		Region:       "us-east-1", AccountID: "000000000000",
	})
	if err != nil {
		t.Fatalf("numeric window rejected: %v", err)
	}
	waitFor(t, 5*time.Second, "import worker to finish", func() bool {
		saved, err := store.GetImportTask(task.ImportId)
		return err == nil && isTerminalImportStatus(saved.ImportStatus)
	})

	if _, err := svc.cancelImportTaskCore(store, "bad id!"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("malformed importId cancel: %v", err)
	}
	if _, err := svc.cancelImportTaskCore(store, strings.Repeat("i", logsstore.MaxImportIdLength+1)); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("oversize importId cancel: %v", err)
	}
	if _, err := svc.cancelImportTaskCore(store, "import-unknown-ok"); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("well-formed unknown importId cancel: %v", err)
	}
	if _, _, err := svc.describeImportTasksCore(store, "bad id!", "", "", "", 0); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("malformed importId describe: %v", err)
	}
	if _, err := svc.describeImportTaskBatchesCore(store, &DescribeImportTaskBatchesInput{ImportId: "bad id!"}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("malformed importId batches: %v", err)
	}
	list, _, err := svc.describeImportTasksCore(store, "import-unknown-ok", "", "", "", 0)
	if err != nil || len(list) != 0 {
		t.Fatalf("well-formed unknown importId describe: %v (%d tasks)", err, len(list))
	}
}
