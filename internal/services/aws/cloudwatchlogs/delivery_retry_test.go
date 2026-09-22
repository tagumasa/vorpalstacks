package cloudwatchlogs

import (
	"context"
	"encoding/base64"
	"errors"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/eventbus"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// flakyKinesisInvoker wraps the recording invoker with a fail switch on
// PutRecord, so a test can simulate a throttled destination and recover
// it.
type flakyKinesisInvoker struct {
	regionRecordingKinesisInvoker
	mu   sync.Mutex
	fail bool
}

func (f *flakyKinesisInvoker) failing() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.fail
}

func (f *flakyKinesisInvoker) PutRecord(ctx context.Context, region, streamName, partitionKey string, data []byte) (string, string, error) {
	if f.failing() {
		return "", "", errors.New("simulated throttling")
	}
	return f.regionRecordingKinesisInvoker.PutRecord(ctx, region, streamName, partitionKey, data)
}

// newRetryTestService builds a service whose Kinesis invoker can fail on
// demand; the tests drive the shared dispatch and the retry drain
// directly, as both delivery legs' call sites do.
func newRetryTestService(t *testing.T) (*LogsService, *flakyKinesisInvoker, *logsstore.Store) {
	t.Helper()
	svc, store := newTestService(t)
	invoker := &flakyKinesisInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	svc.bus = bus
	return svc, invoker, store
}

// rewritePending mutates the stored pending delivery through the store so
// a test can force the due or expired state deterministically.
func rewritePending(t *testing.T, store *logsstore.Store, fn func(*logsstore.PendingDelivery)) {
	t.Helper()
	pending, err := store.ListPendingDeliveries()
	if err != nil || len(pending) != 1 {
		t.Fatalf("expected exactly one pending delivery, got %d (err %v)", len(pending), err)
	}
	fn(pending[0])
	if err := store.PutPendingDelivery(pending[0]); err != nil {
		t.Fatalf("rewrite pending delivery: %v", err)
	}
}

// TestSubscriptionDeliveryRidesRetryWindowThenRecovers pins the
// documented retry window end to end: "Throttled deliverables are retried
// for up to 24 hours. After 24 hours, the failed deliverables are
// dropped" — a failed batch persists with the window's deadline, a due
// re-drive against a recovered destination delivers and clears the
// record.
func TestSubscriptionDeliveryRidesRetryWindowThenRecovers(t *testing.T) {
	svc, invoker, store := newRetryTestService(t)
	const dest = "arn:aws:kinesis:us-east-1:000000000000:stream/retry-dest"
	payload, err := compressJSON(map[string]interface{}{
		"owner":       "000000000000",
		"logGroup":    "retry-group",
		"logStream":   "stream-1",
		"messageType": "DATA_MESSAGE",
		"logEvents":   []map[string]interface{}{{"id": "1", "timestamp": 1700000000000, "message": "ERROR boom"}},
	})
	if err != nil {
		t.Fatalf("compressJSON: %v", err)
	}

	invoker.mu.Lock()
	invoker.fail = true
	invoker.mu.Unlock()
	if err := svc.dispatchSubscriptionDelivery("us-east-1", dest, "retry-group", "stream-1", "ByLogStream", payload); err == nil {
		t.Fatal("a throttled destination must fail the dispatch")
	} else {
		svc.retryFailedDelivery(store, "us-east-1", dest, "retry-group", "stream-1", "ByLogStream", payload, err)
	}

	pending, err := store.ListPendingDeliveries()
	if err != nil || len(pending) != 1 {
		t.Fatalf("failed delivery must persist one pending record, got %d (err %v)", len(pending), err)
	}
	pd := pending[0]
	window := time.Duration(pd.Deadline-pd.FirstAttempt) * time.Millisecond
	if window < 24*time.Hour-5*time.Minute || window > 24*time.Hour+5*time.Minute {
		t.Fatalf("deadline must sit the documented 24 hours past the first attempt, got %v", window)
	}
	if pd.DestArn != dest || pd.LogGroup != "retry-group" || string(pd.Payload) != string(payload) {
		t.Fatalf("pending record must carry the batch addressing and payload, got %+v", pd)
	}

	// Recover the destination and force the record due.
	invoker.mu.Lock()
	invoker.fail = false
	invoker.mu.Unlock()
	rewritePending(t, store, func(p *logsstore.PendingDelivery) { p.NextAttempt = time.Now().Add(-time.Second).UnixMilli() })
	svc.drainDueSubscriptionDeliveries(store)

	if n := invoker.putCount(); n != 1 {
		t.Fatalf("the re-drive must deliver the batch, got %d PutRecord calls", n)
	}
	if pending, err := store.ListPendingDeliveries(); err != nil || len(pending) != 0 {
		t.Fatalf("a delivered batch must clear its pending record, got %d (err %v)", len(pending), err)
	}
	_, _, _, data := invoker.snapshot()
	var redriven map[string]interface{}
	if err := gunzipJSON(wireBlob(t, data[0]), &redriven); err != nil {
		t.Fatalf("re-driven record must gunzip to JSON: %v", err)
	}
	if redriven["messageType"] != "DATA_MESSAGE" {
		t.Fatalf("re-driven record must carry the original batch payload, got %v", redriven["messageType"])
	}
	events, _ := redriven["logEvents"].([]interface{})
	if len(events) != 1 {
		t.Fatalf("re-driven record must carry the matched event, got %v", redriven["logEvents"])
	}
}

// TestSubscriptionDeliveryBackoffReschedules pins the failure half of the
// loop: a still-failing re-drive increments the attempt count and pushes
// the next attempt out under the backoff, and a not-yet-due record is not
// redispatched.
func TestSubscriptionDeliveryBackoffReschedules(t *testing.T) {
	svc, invoker, store := newRetryTestService(t)
	const dest = "arn:aws:kinesis:us-east-1:000000000000:stream/backoff-dest"
	payload := []byte("compressed-batch")

	invoker.mu.Lock()
	invoker.fail = true
	invoker.mu.Unlock()
	if err := svc.dispatchSubscriptionDelivery("us-east-1", dest, "backoff-group", "stream-1", "ByLogStream", payload); err == nil {
		t.Fatal("a throttled destination must fail the dispatch")
	} else {
		svc.retryFailedDelivery(store, "us-east-1", dest, "backoff-group", "stream-1", "ByLogStream", payload, err)
	}

	rewritePending(t, store, func(p *logsstore.PendingDelivery) { p.NextAttempt = time.Now().Add(-time.Second).UnixMilli() })
	svc.drainDueSubscriptionDeliveries(store)

	pending, err := store.ListPendingDeliveries()
	if err != nil || len(pending) != 1 {
		t.Fatalf("a still-failing batch must stay pending, got %d (err %v)", len(pending), err)
	}
	if pending[0].Attempts != 1 {
		t.Fatalf("the re-driven failure must count its attempt, got %d", pending[0].Attempts)
	}
	backoff := time.Duration(pending[0].NextAttempt-time.Now().UnixMilli()) * time.Millisecond
	if backoff <= 0 || backoff > 3*time.Second {
		t.Fatalf("the next attempt must sit one backoff step out, got %v", backoff)
	}

	// Not yet due: a second drain must not redispatch.
	svc.drainDueSubscriptionDeliveries(store)
	if pending, err := store.ListPendingDeliveries(); err != nil || len(pending) != 1 || pending[0].Attempts != 1 {
		t.Fatalf("a not-yet-due record must not redispatch, got %+v (err %v)", pending, err)
	}
}

// TestSubscriptionDeliveryDroppedAfterWindow pins the expiry half: a
// record past its deadline is dropped without another dispatch attempt.
func TestSubscriptionDeliveryDroppedAfterWindow(t *testing.T) {
	svc, invoker, store := newRetryTestService(t)
	if err := store.PutPendingDelivery(&logsstore.PendingDelivery{
		ID:           "expired",
		Region:       "us-east-1",
		DestArn:      "arn:aws:kinesis:us-east-1:000000000000:stream/expired-dest",
		LogGroup:     "g",
		LogStream:    "s",
		Distribution: "ByLogStream",
		Payload:      []byte("compressed"),
		FirstAttempt: time.Now().Add(-25 * time.Hour).UnixMilli(),
		Deadline:     time.Now().Add(-time.Hour).UnixMilli(),
		NextAttempt:  time.Now().Add(-time.Minute).UnixMilli(),
		Attempts:     3,
	}); err != nil {
		t.Fatalf("seed pending delivery: %v", err)
	}

	svc.drainDueSubscriptionDeliveries(store)

	if pending, err := store.ListPendingDeliveries(); err != nil || len(pending) != 0 {
		t.Fatalf("an expired batch must drop, got %d (err %v)", len(pending), err)
	}
	if n := invoker.putCount(); n != 0 {
		t.Fatalf("an expired batch must not redispatch, got %d PutRecord calls", n)
	}
}

// TestPutSubscriptionFilterEmitsControlMessage pins the documented
// reachability record: "Sometimes CloudWatch Logs may emit Amazon Kinesis
// Data Streams records with a 'CONTROL_MESSAGE' type, mainly for checking
// if the destination is reachable" — the probe arrives at filter creation
// with the CONTROL_MESSAGE type and no log events, ahead of any data
// batch.
func TestPutSubscriptionFilterEmitsControlMessage(t *testing.T) {
	svc, invoker, store := newRetryTestService(t)
	svc.bus = eventbus.NewEventBus()
	svc.bus.SetKinesisInvoker(invoker)
	if err := store.CreateLogGroup(logsstore.NewLogGroup("probe-group", "us-east-1", "000000000000")); err != nil {
		t.Fatalf("create log group: %v", err)
	}
	if err := svc.putSubscriptionFilterCore(context.Background(), store, &PutSubscriptionFilterInput{
		LogGroupName:     "probe-group",
		FilterName:       "probe",
		FilterPattern:    "",
		FilterPatternSet: true,
		DestinationArn:   "arn:aws:kinesis:us-east-1:000000000000:stream/probe-dest",
		Region:           "us-east-1",
	}); err != nil {
		t.Fatalf("put subscription filter: %v", err)
	}

	_, _, _, data := invoker.snapshot()
	if len(data) != 1 {
		t.Fatalf("the probe must arrive as one record, got %d", len(data))
	}
	var probe map[string]interface{}
	if err := gunzipJSON(wireBlob(t, data[0]), &probe); err != nil {
		t.Fatalf("probe record must gunzip to JSON: %v", err)
	}
	if probe["messageType"] != "CONTROL_MESSAGE" {
		t.Fatalf("probe record must carry the CONTROL_MESSAGE type, got %v", probe["messageType"])
	}
	if probe["logGroup"] != "probe-group" {
		t.Fatalf("probe record must carry the filter's log group, got %v", probe["logGroup"])
	}
	if filters, _ := probe["subscriptionFilters"].([]interface{}); len(filters) != 1 || filters[0] != "probe" {
		t.Fatalf("probe record must name its filter, got %v", probe["subscriptionFilters"])
	}
	if events, _ := probe["logEvents"].([]interface{}); len(events) != 0 {
		t.Fatalf("probe record must carry no log events, got %v", probe["logEvents"])
	}

	// Data batches keep their own shape after the probe.
	if err := svc.dispatchSubscriptionDelivery("us-east-1", "arn:aws:kinesis:us-east-1:000000000000:stream/probe-dest", "probe-group", "stream-1", "ByLogStream",
		mustCompress(t, map[string]interface{}{
			"owner":       "000000000000",
			"logGroup":    "probe-group",
			"logStream":   "stream-1",
			"messageType": "DATA_MESSAGE",
			"logEvents":   []map[string]interface{}{{"id": "1", "timestamp": 1700000000000, "message": "data row"}},
		})); err != nil {
		t.Fatalf("data batch after the probe must deliver: %v", err)
	}
	_, _, _, data = invoker.snapshot()
	if len(data) != 2 {
		t.Fatalf("a data batch after the probe must deliver, got %d records", len(data))
	}
	var batch map[string]interface{}
	if err := gunzipJSON(wireBlob(t, data[1]), &batch); err != nil {
		t.Fatalf("data record must gunzip to JSON: %v", err)
	}
	if batch["messageType"] != "DATA_MESSAGE" {
		t.Fatalf("data record must be the DATA_MESSAGE batch, got %v", batch["messageType"])
	}
	if events, _ := batch["logEvents"].([]interface{}); len(events) != 1 {
		t.Fatalf("data record must carry the matched event, got %v", batch["logEvents"])
	}
}

// wireBlob decodes the record's stored wire form (the base64 text the
// platform keeps for every SDK-written record) to the gzip bytes.
func wireBlob(t *testing.T, data []byte) []byte {
	t.Helper()
	blob, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		t.Fatalf("record data must be the base64 wire presentation: %v", err)
	}
	return blob
}

// mustCompress gzips one subscription payload for the tests.
func mustCompress(t *testing.T, payload map[string]interface{}) []byte {
	t.Helper()
	compressed, err := compressJSON(payload)
	if err != nil {
		t.Fatalf("compressJSON: %v", err)
	}
	return compressed
}
