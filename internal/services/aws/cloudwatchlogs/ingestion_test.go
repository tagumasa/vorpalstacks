package cloudwatchlogs

import (
	"context"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/eventbus"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The PutLogEvents wire parsers must distinguish an explicit timestamp
// of zero — a legal member value, the Timestamp shape's range starts at
// zero — from an absent member, and must keep neither-member and
// non-object array elements in the batch so the seam's member
// validation rejects the request instead of the elements being
// silently dropped.
func TestParseLogEventsMemberPresence(t *testing.T) {
	flat := parseLogEvents(&request.ParsedRequest{Parameters: map[string]interface{}{
		"LogEvents.1.Timestamp": 0,
		"LogEvents.1.Message":   "zero",
	}})
	if len(flat) != 1 || !flat[0].TimestampSet || flat[0].Message != "zero" {
		t.Fatalf("flat explicit zero timestamp: got %+v", flat)
	}

	flatAbsent := parseLogEvents(&request.ParsedRequest{Parameters: map[string]interface{}{
		"LogEvents.1.Message": "msg only",
	}})
	if len(flatAbsent) != 1 || flatAbsent[0].TimestampSet {
		t.Fatalf("flat absent timestamp: got %+v", flatAbsent)
	}

	mapped := parseLogEvents(&request.ParsedRequest{Parameters: map[string]interface{}{
		"logEvents": []interface{}{
			map[string]interface{}{"message": "ok", "timestamp": 1.0},
			map[string]interface{}{"foo": 1},
			"bare string",
		},
	}})
	if len(mapped) != 3 {
		t.Fatalf("map form dropped elements: got %+v", mapped)
	}
	if mapped[0].Message != "ok" || !mapped[0].TimestampSet {
		t.Fatalf("map form first element: got %+v", mapped[0])
	}
	if mapped[1].TimestampSet || mapped[1].Message != "" {
		t.Fatalf("map form neither-member element: got %+v", mapped[1])
	}
	if mapped[2].TimestampSet || mapped[2].Message != "" {
		t.Fatalf("map form non-object element: got %+v", mapped[2])
	}
}

// The seam's timestamp edges: an explicit zero is a legal member the
// age rules then reject per-event (tooOld, reported not thrown), while
// a negative timestamp violates the shape's range and rejects the
// whole request.
func TestPutLogEventsTimestampEdges(t *testing.T) {
	svc := newIngestionTestServiceBare(t)
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}
	const group, stream = "tsedges-group", "tsedges-stream"
	if err := store.CreateLogGroup(logsstore.NewLogGroup(group, "us-east-1", "000000000000")); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, group)); err != nil {
		t.Fatal(err)
	}

	res, err := svc.putLogEventsCore(PutLogEventsInput{
		LogGroupName: group, LogStreamName: stream, Region: "us-east-1",
		Events: []PutLogEvent{{LogEntry: logsstore.LogEntry{Timestamp: 0, Message: "epoch"}, TimestampSet: true}},
	})
	if err != nil {
		t.Fatalf("explicit zero timestamp must not reject the request: %v", err)
	}
	if res.RejectedLogEvents == nil {
		t.Fatal("explicit zero timestamp must be age-rejected per event")
	}
	if got := res.RejectedLogEvents["tooOldLogEventEndIndex"]; got != 1 {
		t.Fatalf("tooOldLogEventEndIndex = %v, want 1", got)
	}

	_, err = svc.putLogEventsCore(PutLogEventsInput{
		LogGroupName: group, LogStreamName: stream, Region: "us-east-1",
		Events: []PutLogEvent{{LogEntry: logsstore.LogEntry{Timestamp: -1, Message: "before epoch"}, TimestampSet: true}},
	})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("negative timestamp: code=%q err=%v", code, err)
	}
}

// recordingMetricInvoker captures every metric emission the filter fan-out
// makes, with the region it was addressed to (and the dimensions the
// dimensioned form carries).
type recordingMetricInvoker struct {
	mu    sync.Mutex
	calls []string // "region|namespace|metric"
	dims  []map[string]string
}

func (r *recordingMetricInvoker) PutMetricData(region, namespace, metricName string, _ float64, _ time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, region+"|"+namespace+"|"+metricName)
	r.dims = append(r.dims, nil)
	return nil
}

func (r *recordingMetricInvoker) PutMetricDataWithDimensions(region, namespace, metricName string, dimensions map[string]string, _ float64, _ time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, region+"|"+namespace+"|"+metricName)
	r.dims = append(r.dims, dimensions)
	return nil
}

func (r *recordingMetricInvoker) PutMetricDataWithDimensionsAndUnit(region, namespace, metricName string, dimensions map[string]string, _ string, _ float64, _ time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, region+"|"+namespace+"|"+metricName)
	r.dims = append(r.dims, dimensions)
	return nil
}

func (r *recordingMetricInvoker) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.calls)
}

// waitFor polls cond until it holds or the deadline passes, failing the
// test with message on timeout. The filter fan-out is asynchronous by
// design, so its effects are observed this way rather than asserted
// immediately.
func waitFor(t *testing.T, deadline time.Duration, message string, cond func() bool) {
	t.Helper()
	deadlineCh := time.After(deadline)
	for {
		if cond() {
			return
		}
		select {
		case <-deadlineCh:
			t.Fatalf("timed out waiting for %s", message)
		case <-time.After(5 * time.Millisecond):
		}
	}
}

// newIngestionTestService builds a service backed by temporary storage with
// the metric invoker recording and a bus carrying the recording Kinesis
// invoker, plus the group/stream pair the tests ingest into.
func newIngestionTestService(t *testing.T) (*LogsService, *recordingMetricInvoker, *regionRecordingKinesisInvoker, string, string) {
	t.Helper()
	svc, _ := newTestService(t)

	metrics := &recordingMetricInvoker{}
	svc.SetCloudWatchMetricInvoker(metrics)

	kinesis := &regionRecordingKinesisInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(kinesis)
	busCtx, busCancel := context.WithCancel(context.Background())
	if err := bus.Start(busCtx); err != nil {
		t.Fatalf("bus start: %v", err)
	}
	t.Cleanup(busCancel)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatalf("SetEventBus: %v", err)
	}

	const group, stream = "ingest-group", "ingest-stream"
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}
	createTestLogGroup(t, store, group)
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, group)); err != nil {
		t.Fatalf("create log stream: %v", err)
	}
	if err := store.PutMetricFilter(&logsstore.MetricFilter{
		Name:          "errors",
		LogGroupName:  group,
		FilterPattern: "ERROR",
		MetricTransformations: []logsstore.MetricTransformation{{
			MetricName:      "ErrorCount",
			MetricNamespace: "App",
			MetricValue:     "1",
		}},
	}); err != nil {
		t.Fatalf("put metric filter: %v", err)
	}
	if err := store.PutSubscriptionFilter(&logsstore.SubscriptionFilter{
		LogGroupName:   group,
		FilterName:     "deliver",
		FilterPattern:  "ERROR",
		DestinationArn: "arn:aws:kinesis:us-east-1:000000000000:stream/logs-dest",
	}); err != nil {
		t.Fatalf("put subscription filter: %v", err)
	}
	return svc, metrics, kinesis, group, stream
}

// matchingEntry returns one log entry that matches both test filters.
func matchingEntry() eventbus.LogEntry {
	return eventbus.LogEntry{Timestamp: time.Now().UnixMilli(), Message: "ERROR something failed"}
}

// TestEveryIngestionPathAppliesMetricAndSubscriptionFilters pins the
// single-ingestion-seam property: the HTTP core, the three bus handlers and
// the cross-service invoker entry all evaluate metric and subscription
// filters. Before the seam, only the HTTP path and the Lambda bus path did,
// so a subscription consumer received Lambda-written events but not
// EventBridge/Scheduler/APIGW/invoker-written events in the same group.
func TestEveryIngestionPathAppliesMetricAndSubscriptionFilters(t *testing.T) {
	svc, metrics, kinesis, group, stream := newIngestionTestService(t)
	now := time.Now().UnixMilli()

	// HTTP plane (the core the PutLogEvents handler calls).
	if _, err := svc.putLogEventsCore(PutLogEventsInput{
		LogGroupName: group, LogStreamName: stream, Region: "us-east-1",
		Events: []PutLogEvent{{LogEntry: logsstore.LogEntry{Timestamp: now, Message: "ERROR via http"}, TimestampSet: true}},
	}); err != nil {
		t.Fatalf("putLogEventsCore: %v", err)
	}

	// Bus planes: Lambda log write, direct put, API Gateway access log.
	lambdaEvt := &eventbus.LambdaLogWriteEvent{
		LogGroup: group, LogStream: stream, LogEvents: []eventbus.LogEntry{matchingEntry()},
	}
	lambdaEvt.Region = "us-east-1"
	svc.handleLambdaLogWrite(context.Background(), lambdaEvt)
	directEvt := &eventbus.CloudWatchLogsPutEvent{
		LogGroup: group, LogStream: stream, LogEvents: []eventbus.LogEntry{matchingEntry()},
	}
	directEvt.Region = "us-east-1"
	directEvt.AccountID = "000000000000"
	svc.handleDirectPutLogEvents(context.Background(), directEvt)
	apigwEvt := &eventbus.APIGatewayAccessLogEvent{LogGroup: group, LogStream: stream, FormattedLog: "ERROR via apigw"}
	apigwEvt.Region = "us-east-1"
	svc.handleAPIGatewayAccessLog(context.Background(), apigwEvt)

	// Invoker plane (the adapter's entry).
	if err := svc.IngestLogEvents("us-east-1", group, stream,
		[]logsstore.LogEntry{{Timestamp: now, Message: "ERROR via invoker"}}); err != nil {
		t.Fatalf("IngestLogEvents: %v", err)
	}

	waitFor(t, 5*time.Second, "all five ingestion paths to evaluate the metric filter", func() bool {
		return metrics.count() >= 5
	})
	waitFor(t, 5*time.Second, "all five ingestion paths to deliver to the subscription destination", func() bool {
		return kinesis.putCount() >= 5
	})
}

// TestIngestLogEventsValidatesInvokerWrites pins that the invoker entry no
// longer bypasses validation: a batch with a stale timestamp is rejected
// (the age window every PutLogEvents write is subject to) and writes
// nothing.
func TestIngestLogEventsValidatesInvokerWrites(t *testing.T) {
	svc, metrics, _, group, stream := newIngestionTestService(t)

	stale := time.Now().Add(-30 * 24 * time.Hour).UnixMilli()
	err := svc.IngestLogEvents("us-east-1", group, stream,
		[]logsstore.LogEntry{{Timestamp: stale, Message: "ERROR too old"}})
	if err != nil {
		t.Fatalf("stale batch reports through rejectedLogEventsInfo, not an error: %v", err)
	}

	store, serr := svc.getLogsStoreByRegion("us-east-1")
	if serr != nil {
		t.Fatalf("logs store: %v", serr)
	}
	events, _, _, gerr := store.GetLogEvents(group, stream, 0, 0, 0, true, "")
	if gerr != nil {
		t.Fatalf("GetLogEvents: %v", gerr)
	}
	if len(events) != 0 {
		t.Fatalf("stale batch must not be written, got %d events", len(events))
	}
	time.Sleep(100 * time.Millisecond)
	if n := metrics.count(); n != 0 {
		t.Fatalf("stale batch must not fan out to filters, got %d metric emissions", n)
	}
}

// TestIngestBusEventsAutoCreatesGroupAndStream pins the platform-internal
// writers' precondition: bus ingestion creates the missing group and
// stream instead of failing like the API plane.
func TestIngestBusEventsAutoCreatesGroupAndStream(t *testing.T) {
	svc := newIngestionTestServiceBare(t)
	svc.ingestBusEvents("test source", "us-east-1", "fresh-group", "fresh-stream", "000000000000",
		[]eventbus.LogEntry{{Timestamp: time.Now().UnixMilli(), Message: "created"}})

	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}
	events, _, _, gerr := store.GetLogEvents("fresh-group", "fresh-stream", 0, 0, 0, true, "")
	if gerr != nil {
		t.Fatalf("GetLogEvents after bus ingestion: %v", gerr)
	}
	if len(events) != 1 || events[0].Message != "created" {
		t.Fatalf("bus ingestion must auto-create and write, got %d events", len(events))
	}
}

// newIngestionTestServiceBare builds a service with storage only — no
// filters, no invokers.
func newIngestionTestServiceBare(t *testing.T) *LogsService {
	t.Helper()
	svc, _ := newTestService(t)
	return svc
}

// TestDispatchFilterFanOutBoundedSubmitOnFullQueue pins the bound's
// safety property: a saturated queue makes submission wait its bounded
// window and then drop one batch's evaluation instead of deadlocking
// the write path (whose callers can include goroutines the pool itself
// provoked through subscription delivery).
func TestDispatchFilterFanOutBoundedSubmitOnFullQueue(t *testing.T) {
	svc := newIngestionTestServiceBare(t)
	svc.filterJobs = make(chan filterFanOutJob, 1) // no workers: the queued job stays put
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}
	events := []logsstore.LogEntry{{Timestamp: time.Now().UnixMilli(), Message: "x"}}

	done := make(chan struct{})
	go func() {
		defer close(done)
		svc.dispatchFilterFanOut(store, "us-east-1", "g", "s", events, nil) // fills the queue
		svc.dispatchFilterFanOut(store, "us-east-1", "g", "s", events, nil) // waits the window, then drops
	}()
	select {
	case <-done:
	case <-time.After(filterFanOutSubmitTimeout + 2*time.Second):
		t.Fatal("dispatchFilterFanOut blocked past its submit window on a full queue")
	}
	// With the queued job drained, a submission accepts immediately: the
	// bound only engages on saturation.
	<-svc.filterJobs
	accepted := make(chan struct{})
	go func() {
		defer close(accepted)
		svc.dispatchFilterFanOut(store, "us-east-1", "g", "s", events, nil)
	}()
	select {
	case <-accepted:
	case <-time.After(2 * time.Second):
		t.Fatal("dispatchFilterFanOut did not accept with queue capacity")
	}
}

// TestDispatchFilterFanOutDropDeliversSubscriptionsInline pins the drop
// path's durability: a batch dropped past the submit window (the queue
// saturated) still carries its subscription deliveries — the drop path
// evaluates the subscription filters inline, so a dispatch failure rides
// the retry window instead of vanishing with a warning.
func TestDispatchFilterFanOutDropDeliversSubscriptionsInline(t *testing.T) {
	svc, invoker, store := newRetryTestService(t)
	const dest = "arn:aws:kinesis:us-east-1:000000000000:stream/drop-dest"
	if err := store.CreateLogGroup(logsstore.NewLogGroup("drop-group", "us-east-1", "000000000000")); err != nil {
		t.Fatal(err)
	}
	if err := store.PutSubscriptionFilter(&logsstore.SubscriptionFilter{
		LogGroupName:   "drop-group",
		FilterName:     "deliver",
		FilterPattern:  "ERROR",
		DestinationArn: dest,
		Distribution:   "Random",
	}); err != nil {
		t.Fatal(err)
	}

	// A saturated queue with no workers: the first submit fills the slot,
	// the second waits its window and takes the drop path.
	svc.filterJobs = make(chan filterFanOutJob, 1)
	invoker.mu.Lock()
	invoker.fail = true
	invoker.mu.Unlock()
	events := []logsstore.LogEntry{{Timestamp: time.Now().UnixMilli(), Message: "ERROR dropped"}}
	done := make(chan struct{})
	go func() {
		defer close(done)
		svc.dispatchFilterFanOut(store, "us-east-1", "drop-group", "s", events, nil)
		svc.dispatchFilterFanOut(store, "us-east-1", "drop-group", "s", events, nil)
	}()
	select {
	case <-done:
	case <-time.After(filterFanOutSubmitTimeout + 5*time.Second):
		t.Fatal("dispatchFilterFanOut blocked past its submit window on a full queue")
	}

	pending, err := store.ListPendingDeliveries()
	if err != nil || len(pending) != 1 {
		t.Fatalf("the dropped batch's failed delivery must ride the retry window, got %d pending (err %v)", len(pending), err)
	}
	if pending[0].DestArn != dest || pending[0].LogGroup != "drop-group" {
		t.Fatalf("the pending record must carry the subscription addressing, got %+v", pending[0])
	}
}

// TestFilterFanOutWorkerDrainsQueuedJobsOnShutdown pins the shutdown
// drain: a worker that sees the context cancelled evaluates the jobs the
// queue already accepted instead of abandoning them — their metric
// emissions and subscription deliveries leave no retry trail unless the
// evaluation runs.
func TestFilterFanOutWorkerDrainsQueuedJobsOnShutdown(t *testing.T) {
	svc, metrics, _, group, _ := newIngestionTestService(t)
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatal(err)
	}

	// Swap the dispatch channel: the service's own workers keep the
	// original, and the worker under test serves the swapped one.
	queued := make(chan filterFanOutJob, 2)
	svc.filterJobs = queued
	events := []logsstore.LogEntry{{Timestamp: time.Now().UnixMilli(), Message: "ERROR drained"}}
	job := filterFanOutJob{store: store, region: "us-east-1", logGroup: group, logStream: "ingest-stream", events: events}
	queued <- job
	queued <- job

	svc.cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		svc.runFilterFanOutWorker(queued)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("the worker did not return after shutdown")
	}
	if got := metrics.count(); got < 2 {
		t.Fatalf("the worker must drain both queued jobs on shutdown, evaluated %d", got)
	}
}

// TestRunFilterFanOutJobSurvivesPanic pins the per-job recovery: a panicking
// evaluation is logged and the worker lives on to process the next batch.
func TestRunFilterFanOutJobSurvivesPanic(t *testing.T) {
	svc, metrics, _, group, _ := newIngestionTestService(t)
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}

	// A nil store makes the evaluation calls panic before any guard can
	// reject the job.
	svc.runFilterFanOutJob(filterFanOutJob{logGroup: "g", events: []logsstore.LogEntry{{
		Timestamp: time.Now().UnixMilli(), Message: "ERROR panicking",
	}}})

	events := []logsstore.LogEntry{{Timestamp: time.Now().UnixMilli(), Message: "ERROR after panic"}}
	svc.runFilterFanOutJob(filterFanOutJob{
		store: store, region: "us-east-1", logGroup: group, logStream: "ingest-stream", events: events,
	})
	waitFor(t, 5*time.Second, "post-panic job to evaluate the metric filter", func() bool {
		return metrics.count() >= 1
	})
}

// PutLogEvents rejects events preceding the group's retention period
// ("Events older than 14 days or preceding the log group's retention
// period are rejected while processing remaining valid events") and
// reports them through rejectedLogEventsInfo's expiredLogEventEndIndex;
// an event inside the retention window still ingests.
func TestPutLogEventsRejectsEventsPrecedingRetention(t *testing.T) {
	const group, stream = "retention-group", "retention-stream"
	svc, store := newReadTestService(t, group)
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, group)); err != nil {
		t.Fatal(err)
	}
	lg, err := store.GetLogGroup(group)
	if err != nil {
		t.Fatal(err)
	}
	lg.SetRetention(1)
	if err := store.PutLogGroup(lg); err != nil {
		t.Fatal(err)
	}

	now := time.Now().UnixMilli()
	res, err := svc.putLogEventsCore(PutLogEventsInput{
		LogGroupName: group, LogStreamName: stream, Region: "us-east-1",
		Events: []PutLogEvent{
			{LogEntry: logsstore.LogEntry{Timestamp: now - 2*24*60*60*1000, Message: "retention-preceding"}, TimestampSet: true},
			{LogEntry: logsstore.LogEntry{Timestamp: now, Message: "current"}, TimestampSet: true},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.RejectedLogEvents == nil {
		t.Fatal("retention-preceding event was not reported")
	}
	if got := res.RejectedLogEvents["expiredLogEventEndIndex"]; got != 1 {
		t.Fatalf("expiredLogEventEndIndex = %v, want 1", got)
	}
	saved, _, _, err := store.GetLogEvents(group, stream, 0, 0, 10, true, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(saved) != 1 || saved[0].Message != "current" {
		t.Fatalf("stored events = %v, want the current one alone", saved)
	}
}
