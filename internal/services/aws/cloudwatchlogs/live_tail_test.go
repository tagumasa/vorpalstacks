package cloudwatchlogs

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/pkg/filterpattern"
)

// Events ingested in one batch share an ingestion millisecond: the
// session's dedupe baseline must be the watermark as of the fetch, so
// the whole batch is delivered — comparing against the running maximum
// dropped every event after the first of each millisecond — and the
// advanced watermark dedupes the re-read window on the next collect.
func TestLiveTailCollectDeliversWholeSharedIngestionBatch(t *testing.T) {
	svc, store := newIndexPolicyTestEnv(t, "lt-cursor")
	// The events sit comfortably below collect's internal now: the fetch
	// window's end boundary is exclusive, so an event timestamped on the
	// boundary millisecond would be dropped whenever the collect clock
	// lands on it — a race, not the sharing behaviour under test.
	base := time.Now().UnixMilli() - 60_000
	putIndexEvents(t, svc, "lt-cursor", base, "m1-live-tail-marker", "m2-live-tail-marker")

	session := &liveTailSession{
		store:            store,
		groupNames:       []string{"lt-cursor"},
		groupARNs:        []string{store.ARNBuilder().CloudWatch().LogGroup("lt-cursor")},
		input:            StartLiveTailInput{LogEventFilterPattern: "live-tail-marker"},
		startedIngestion: 0,
		matcher:          filterpattern.NewMatcher(),
	}
	cursors := map[string]liveTailCursor{}
	got := session.collect(cursors)
	if len(got) != 2 {
		t.Fatalf("a batch sharing one ingestion millisecond must deliver both events, got %d", len(got))
	}
	if again := session.collect(cursors); len(again) != 0 {
		t.Fatalf("re-collect over the advanced watermark must deliver nothing, got %d", len(again))
	}
}

// TestLiveTailCursorAdvancesPastFilterNonMatches pins the consumption
// cursor's contract: every scanned event advances the fetch floor and
// the ingestion watermark, so a filter matching nothing still moves the
// second tick's scan to the first tick's end — a non-matching tail is
// consumed, not re-fetched and re-matched every second.
func TestLiveTailCursorAdvancesPastFilterNonMatches(t *testing.T) {
	svc, store := newIndexPolicyTestEnv(t, "lt-advance")
	base := time.Now().UnixMilli() - 60_000
	putIndexEvents(t, svc, "lt-advance", base, "nm-1", "nm-2", "nm-3")

	var starts []int64
	orig := liveTailFetchEvents
	liveTailFetchEvents = func(st *logsstore.Store, group, stream string, start, end int64, max int) ([]*logsstore.OutputLogEvent, error) {
		starts = append(starts, start)
		return orig(st, group, stream, start, end, max)
	}
	t.Cleanup(func() { liveTailFetchEvents = orig })

	session := &liveTailSession{
		store:            store,
		groupNames:       []string{"lt-advance"},
		groupARNs:        []string{store.ARNBuilder().CloudWatch().LogGroup("lt-advance")},
		input:            StartLiveTailInput{LogEventFilterPattern: "never-matches-anything"},
		startedIngestion: 0,
		matcher:          filterpattern.NewMatcher(),
	}
	cursors := map[string]liveTailCursor{}
	if got := session.collect(cursors); len(got) != 0 {
		t.Fatalf("filter matching nothing must emit no events, got %d", len(got))
	}
	if len(starts) == 0 || starts[0] != 0 {
		t.Fatalf("first tick's fetch floor = %v, want the initial cursor 0", starts)
	}
	key := "lt-advance\x00s1"
	if cursors[key].ts != base+2 {
		t.Fatalf("cursor.ts after the first tick = %d, want the last scanned timestamp %d", cursors[key].ts, base+2)
	}
	before := len(starts)
	if again := session.collect(cursors); len(again) != 0 {
		t.Fatalf("second tick over the advanced cursor must emit nothing, got %d", len(again))
	}
	if n := len(starts) - before; n != 1 {
		t.Fatalf("second tick fetched %d times, want 1", n)
	}
	if second := starts[before]; second != base+2 {
		t.Fatalf("second tick's fetch floor = %d, want the first tick's last scanned timestamp %d — the rescan must be O(1), not the full history from 0", second, base+2)
	}
}

// TestLiveTailRunReturnsWhenEnqueueBlockedOnStalledClient pins the
// enqueue select's cancellation arm: a client that stays connected but
// stops reading fills the frame queue and parks the writer on the pipe,
// so no frame drain and no writer exit ever arrive — the session loop
// must return on the context's cancellation rather than holding one of
// the concurrent-session slots indefinitely.
func TestLiveTailRunReturnsWhenEnqueueBlockedOnStalledClient(t *testing.T) {
	_, store := newIndexPolicyTestEnv(t, "lt-stall")
	session := &liveTailSession{
		store:            store,
		groupNames:       nil,
		groupARNs:        nil,
		input:            StartLiveTailInput{},
		startedIngestion: time.Now().UnixMilli(),
		matcher:          filterpattern.NewMatcher(),
	}
	// The stalled client: the queue is already full and the writer never
	// consumes another frame nor exits.
	frames := make(chan *liveTailFrame, logsstore.LiveTailBufferUpdatesMax)
	for i := 0; i < logsstore.LiveTailBufferUpdatesMax; i++ {
		frames <- &liveTailFrame{kind: "update"}
	}
	writerDone := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		session.run(ctx, frames, writerDone)
		close(done)
	}()
	// One full ticker period with margin: the first tick's collect (no
	// groups) returns at once and the loop parks in the enqueue select.
	time.Sleep(2 * time.Second)
	cancel()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("run did not return after cancellation while the frame queue was full and the writer stalled")
	}
}

// ltGroupARN addresses a group in the member's documented form
// ("Specify each log group by its ARN", StartLiveTail).
func ltGroupARN(name string) string {
	return "arn:aws:logs:us-east-1:000000000000:log-group:" + name
}

// TestStartLiveTailValidationRows pins the session-request matrix: the
// required identifier array with its ten-group ceiling, the identifier
// alphabet and the ARN-only addressing rule, the trailing asterisk
// rule, the stream-filter pairing rules with the stream-name alphabet,
// the filter-pattern length bound, the unknown group, the Standard-class
// scope and the concurrent-session quota.
func TestStartLiveTailValidationRows(t *testing.T) {
	svc, _ := newReadTestService(t, "lt-group")

	rows := []struct {
		name  string
		input StartLiveTailInput
		code  string
	}{
		{"no identifiers", StartLiveTailInput{Region: "us-east-1"}, "InvalidParameterException"},
		{"eleven groups", StartLiveTailInput{
			Region:              "us-east-1",
			LogGroupIdentifiers: repeaters(ltGroupARN("lt-group"), 11),
		}, "InvalidParameterException"},
		{"trailing asterisk", StartLiveTailInput{
			Region:              "us-east-1",
			LogGroupIdentifiers: []string{"arn:aws:logs:us-east-1:000000000000:log-group:lt-*"},
		}, "InvalidParameterException"},
		{"bare name", StartLiveTailInput{
			Region:              "us-east-1",
			LogGroupIdentifiers: []string{"lt-group"},
		}, "InvalidParameterException"},
		{"identifier alphabet", StartLiveTailInput{
			Region:              "us-east-1",
			LogGroupIdentifiers: []string{ltGroupARN("lt group")},
		}, "InvalidParameterException"},
		{"both stream filters", StartLiveTailInput{
			Region:                "us-east-1",
			LogGroupIdentifiers:   []string{ltGroupARN("lt-group")},
			LogStreamNames:        []string{"s1"},
			LogStreamNamePrefixes: []string{"s"},
		}, "InvalidParameterException"},
		{"stream filters with two groups", StartLiveTailInput{
			Region:              "us-east-1",
			LogGroupIdentifiers: []string{ltGroupARN("lt-group"), ltGroupARN("lt-other")},
			LogStreamNames:      []string{"s1"},
		}, "InvalidParameterException"},
		{"stream name alphabet", StartLiveTailInput{
			Region:              "us-east-1",
			LogGroupIdentifiers: []string{ltGroupARN("lt-group")},
			LogStreamNames:      []string{"s:1"},
		}, "InvalidParameterException"},
		{"stream prefix alphabet", StartLiveTailInput{
			Region:                "us-east-1",
			LogGroupIdentifiers:   []string{ltGroupARN("lt-group")},
			LogStreamNamePrefixes: []string{"s*"},
		}, "InvalidParameterException"},
		{"unknown group", StartLiveTailInput{
			Region:              "us-east-1",
			LogGroupIdentifiers: []string{ltGroupARN("lt-nope")},
		}, "ResourceNotFoundException"},
		// InputLogStreamNames carries @length {1,100} ("Array Members:
		// Minimum number of 1 item."): an explicitly empty list (the
		// non-nil empty slice the parse marks present) is a member
		// violation, not an absent member.
		{"empty stream names", StartLiveTailInput{
			Region:              "us-east-1",
			LogGroupIdentifiers: []string{ltGroupARN("lt-group")},
			LogStreamNames:      []string{},
		}, "InvalidParameterException"},
		{"empty stream prefixes", StartLiveTailInput{
			Region:                "us-east-1",
			LogGroupIdentifiers:   []string{ltGroupARN("lt-group")},
			LogStreamNamePrefixes: []string{},
		}, "InvalidParameterException"},
		{"long filter pattern", StartLiveTailInput{
			Region:                "us-east-1",
			LogGroupIdentifiers:   []string{ltGroupARN("lt-group")},
			LogEventFilterPattern: strings.Repeat("a", logsstore.LiveTailFilterPatternMax+1),
		}, "InvalidParameterException"},
	}
	for _, row := range rows {
		if _, err := svc.startLiveTailCore(context.Background(), row.input); logsErrorCode(err) != row.code {
			t.Fatalf("%s: err = %v, want %s", row.name, err, row.code)
		}
	}

	// The trailing-asterisk row carries its documented message: the rule
	// reads before the alphabet rule, which excludes '*' and would
	// otherwise shadow the asterisk wording with the generic one.
	if _, err := svc.startLiveTailCore(context.Background(), StartLiveTailInput{
		Region:              "us-east-1",
		LogGroupIdentifiers: []string{ltGroupARN("lt-group") + "*"},
	}); logsErrorCode(err) != "InvalidParameterException" || !strings.Contains(err.Error(), "asterisk") {
		t.Fatalf("trailing asterisk = %v, want the documented asterisk message", err)
	}

	// Live Tail is Standard-class only; the non-Standard group answers
	// the operation's InvalidOperationException.
	svcIA, storeIA := newReadTestService(t, "lt-standard")
	ia := logsstore.NewLogGroup("lt-ia", "us-east-1", "000000000000")
	ia.LogGroupClass = "INFREQUENT_ACCESS"
	if err := storeIA.CreateLogGroup(ia); err != nil {
		t.Fatal(err)
	}
	if _, err := svcIA.startLiveTailCore(context.Background(), StartLiveTailInput{
		Region: "us-east-1", LogGroupIdentifiers: []string{ltGroupARN("lt-ia")},
	}); logsErrorCode(err) != "InvalidOperationException" {
		t.Fatalf("non-Standard group = %v, want InvalidOperationException", err)
	}

	// The concurrent-session quota: fifteen sessions stay open, the
	// sixteenth rejects, and closing one frees its slot.
	svcQ, _ := newReadTestService(t, "lt-quota")
	var open []io.ReadCloser
	for i := 0; i < logsstore.LiveTailConcurrentSessionsMax; i++ {
		stream, err := svcQ.startLiveTailCore(context.Background(), StartLiveTailInput{
			Region: "us-east-1", LogGroupIdentifiers: []string{ltGroupARN("lt-quota")},
		})
		if err != nil {
			t.Fatalf("session %d: %v", i, err)
		}
		rc := stream.reader.(io.ReadCloser)
		open = append(open, rc)
	}
	if _, err := svcQ.startLiveTailCore(context.Background(), StartLiveTailInput{
		Region: "us-east-1", LogGroupIdentifiers: []string{ltGroupARN("lt-quota")},
	}); logsErrorCode(err) != "LimitExceededException" {
		t.Fatalf("sixteenth session = %v, want LimitExceededException", err)
	}
	// Closing a session's stream ends it (the writer fails against the
	// closed pipe and the driver releases the slot).
	open[0].Close()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		svcQ.liveTailMu.Lock()
		sessions := svcQ.liveTailSessions
		svcQ.liveTailMu.Unlock()
		if sessions == logsstore.LiveTailConcurrentSessionsMax-1 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	svcQ.liveTailMu.Lock()
	sessions := svcQ.liveTailSessions
	svcQ.liveTailMu.Unlock()
	if sessions != logsstore.LiveTailConcurrentSessionsMax-1 {
		t.Fatalf("open sessions after close = %d, want %d", sessions, logsstore.LiveTailConcurrentSessionsMax-1)
	}
	for _, rc := range open[1:] {
		rc.Close()
	}
}

// TestStartLiveTailSessionStream pins the session's wire contract: the
// initial-response message, the sessionStart event carrying the resolved
// ARN identifiers, the per-second sessionUpdate carrying only events
// ingested after the session opened (no replay), the filter pattern's
// effect, and the update vocabulary's exact key set.
func TestStartLiveTailSessionStream(t *testing.T) {
	svc, store := newIndexPolicyTestEnv(t, "lt-flow")
	now := time.Now().UnixMilli()
	putIndexEvents(t, svc, "lt-flow", now, "before-session-event")

	stream, err := svc.startLiveTailCore(context.Background(), StartLiveTailInput{
		Region:                "us-east-1",
		LogGroupIdentifiers:   []string{ltGroupARN("lt-flow")},
		LogStreamNamePrefixes: []string{"s"},
		LogEventFilterPattern: "live-tail-marker",
	})
	if err != nil {
		t.Fatal(err)
	}
	rc := stream.reader.(io.ReadCloser)
	defer rc.Close()
	groupARN := store.ARNBuilder().CloudWatch().LogGroup("lt-flow")

	initial, err := readEventStreamFrame(rc)
	if err != nil {
		t.Fatal(err)
	}
	if initial.headers[":event-type"] != "initial-response" {
		t.Fatalf("first frame = %v, want initial-response", initial.headers)
	}
	start, err := readEventStreamFrame(rc)
	if err != nil {
		t.Fatal(err)
	}
	if start.headers[":event-type"] != "sessionStart" {
		t.Fatalf("second frame = %v, want sessionStart", start.headers)
	}
	var startBody struct {
		RequestID             string   `json:"requestId"`
		SessionID             string   `json:"sessionId"`
		LogGroupIdentifiers   []string `json:"logGroupIdentifiers"`
		LogStreamNamePrefixes []string `json:"logStreamNamePrefixes"`
		LogEventFilterPattern string   `json:"logEventFilterPattern"`
	}
	if err := json.Unmarshal(start.payload, &startBody); err != nil {
		t.Fatal(err)
	}
	if startBody.RequestID == "" || startBody.SessionID == "" {
		t.Fatalf("session ids = %+v", startBody)
	}
	if len(startBody.LogGroupIdentifiers) != 1 || startBody.LogGroupIdentifiers[0] != groupARN {
		t.Fatalf("sessionStart identifiers = %v, want the resolved ARN", startBody.LogGroupIdentifiers)
	}
	if len(startBody.LogStreamNamePrefixes) != 1 || startBody.LogStreamNamePrefixes[0] != "s" {
		t.Fatalf("sessionStart prefixes = %v", startBody.LogStreamNamePrefixes)
	}
	if startBody.LogEventFilterPattern != "live-tail-marker" {
		t.Fatalf("sessionStart filter pattern = %q", startBody.LogEventFilterPattern)
	}
	// The start event's key vocabulary is the model's member set.
	var rawStart map[string]interface{}
	if err := json.Unmarshal(start.payload, &rawStart); err != nil {
		t.Fatal(err)
	}
	for key := range rawStart {
		switch key {
		case "requestId", "sessionId", "logGroupIdentifiers", "logStreamNames", "logStreamNamePrefixes", "logEventFilterPattern":
		default:
			t.Fatalf("sessionStart emitted %q, outside the modelled member set", key)
		}
	}

	// The session opened: an event ingested now (matching the pattern)
	// streams in a sessionUpdate; the pre-session event never does, and
	// a non-matching event stays out. The sleep keeps the put strictly
	// after the session-open millisecond.
	time.Sleep(20 * time.Millisecond)
	putIndexEvents(t, svc, "lt-flow", now, "live-tail-marker hits", "no-marker-here")

	deadline := time.Now().Add(6 * time.Second)
	var seen map[string]interface{}
	for time.Now().Before(deadline) {
		frame, err := readEventStreamFrame(rc)
		if err != nil {
			t.Fatalf("read update frame: %v", err)
		}
		if frame.headers[":event-type"] != "sessionUpdate" {
			t.Fatalf("frame = %v, want sessionUpdate", frame.headers)
		}
		var update struct {
			SessionMetadata struct {
				Sampled bool `json:"sampled"`
			} `json:"sessionMetadata"`
			SessionResults []map[string]interface{} `json:"sessionResults"`
		}
		if err := json.Unmarshal(frame.payload, &update); err != nil {
			t.Fatal(err)
		}
		if update.SessionMetadata.Sampled {
			t.Fatal("sampled = true under the 500-event ceiling")
		}
		for _, result := range update.SessionResults {
			if result["message"] == "live-tail-marker hits" {
				seen = result
			}
			if msg, _ := result["message"].(string); msg == "before-session-event" || msg == "no-marker-here" {
				t.Fatalf("streamed %q: replay and non-match both stay out", msg)
			}
		}
		if seen != nil {
			break
		}
	}
	if seen == nil {
		t.Fatal("the marked event never streamed within the window")
	}
	if seen["logStreamName"] != "s1" || seen["logGroupIdentifier"] != groupARN {
		t.Fatalf("event identifiers = %v", seen)
	}
	if ts, ok := seen["timestamp"].(float64); !ok || ts <= 0 {
		t.Fatalf("timestamp = %v, want epoch milliseconds", seen["timestamp"])
	}
	if it, ok := seen["ingestionTime"].(float64); !ok || it <= 0 {
		t.Fatalf("ingestionTime = %v, want epoch milliseconds", seen["ingestionTime"])
	}
	for key := range seen {
		switch key {
		case "logStreamName", "logGroupIdentifier", "message", "timestamp", "ingestionTime":
		default:
			t.Fatalf("sessionResults entry emitted %q, outside the modelled member set", key)
		}
	}
}

// repeaters builds n copies of one identifier (eleven copies of one
// group still violate the ten-group ceiling).
func repeaters(id string, n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = id
	}
	return out
}

// esFrame is one decoded event-stream message.
type esFrame struct {
	headers map[string]string
	payload []byte
}

// readEventStreamFrame reads and decodes one AWS event-stream message
// (the writer's own framing: string-valued headers only).
func readEventStreamFrame(r io.Reader) (*esFrame, error) {
	prelude := make([]byte, 12)
	if _, err := io.ReadFull(r, prelude); err != nil {
		return nil, err
	}
	total := binary.BigEndian.Uint32(prelude[0:4])
	headersLen := binary.BigEndian.Uint32(prelude[4:8])
	rest := make([]byte, total-12)
	if _, err := io.ReadFull(r, rest); err != nil {
		return nil, err
	}
	headers := make(map[string]string)
	hb := rest[:headersLen]
	for len(hb) > 0 {
		nameLen := int(hb[0])
		hb = hb[1:]
		name := string(hb[:nameLen])
		hb = hb[nameLen:]
		valueType := hb[0]
		hb = hb[1:]
		if valueType != 7 {
			return nil, fmt.Errorf("unsupported event-stream header value type %d", valueType)
		}
		valueLen := int(binary.BigEndian.Uint16(hb[:2]))
		hb = hb[2:]
		headers[name] = string(hb[:valueLen])
		hb = hb[valueLen:]
	}
	return &esFrame{headers: headers, payload: rest[headersLen : len(rest)-4]}, nil
}

// The wire's explicitly empty array reaches the Core as a present member:
// the handler marks the parsed slice present-but-empty (HasListParam),
// so the list's @length(min 1) trait — unenforceable through the
// nil-collapsing shared accessor alone — rejects on the HTTP plane.
func TestStartLiveTailEmptyArrayWireForm(t *testing.T) {
	svc, _ := newTestService(t)
	reqCtx := request.NewRequestContext(context.Background(), svc.storageManager, "000000000000", "us-east-1")

	_, err := svc.StartLiveTail(context.Background(), reqCtx, vocabRequest(map[string]interface{}{
		"logGroupIdentifiers": []interface{}{ltGroupARN("lt-empty-array")},
		"logStreamNames":      []interface{}{},
	}))
	if logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("empty logStreamNames array: err = %v, want InvalidParameterException", err)
	}
}
