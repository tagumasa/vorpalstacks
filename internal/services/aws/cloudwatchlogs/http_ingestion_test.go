package cloudwatchlogs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// stubIngestionAuthenticator accepts exactly the tokens in its set.
type stubIngestionAuthenticator struct{ valid map[string]bool }

func (s *stubIngestionAuthenticator) AuthenticateLogIngestionToken(token string) error {
	if s.valid[token] {
		return nil
	}
	return errors.New("invalid bearer token")
}

// newIngestionEndpointTest builds the ingestion handler over temporary
// storage carrying one group/stream pair, the region pinned through the
// logs.<region> host form.
func newIngestionEndpointTest(t *testing.T, authenticator IngestionTokenAuthenticator) (http.Handler, *LogsService, string, string) {
	t.Helper()
	svc, store := newTestService(t)
	const group, stream = "ingest-http-group", "ingest-http-stream"
	createTestLogGroup(t, store, group)
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, group)); err != nil {
		t.Fatalf("create log stream: %v", err)
	}
	return NewHTTPIngestionHandler(svc, authenticator), svc, group, stream
}

// postIngestion sends one ingestion request through the handler and
// returns the recorder.
func postIngestion(t *testing.T, handler http.Handler, path, body string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Host = "logs.us-east-1.amazonaws.com"
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

// readIngestedEvents reads the whole stream back through the read plane
// (its scoped pagination tokens page across calls; GetLogEvents terminates
// by re-offering the same token).
func readIngestedEvents(t *testing.T, svc *LogsService, group, stream string) []*logsstore.OutputLogEvent {
	t.Helper()
	var events []*logsstore.OutputLogEvent
	var token string
	// The read's parameters must stay identical across pages: the scoped
	// pagination token rejects a replay whose window shifted.
	endTime := time.Now().Add(time.Hour).UnixMilli()
	for {
		result, err := svc.getLogEventsCore(GetLogEventsInput{
			LogGroupName:  group,
			LogStreamName: stream,
			EndTime:       endTime,
			Limit:         10000,
			StartFromHead: true,
			NextToken:     token,
			Region:        "us-east-1",
		})
		if err != nil {
			t.Fatalf("read events: %v", err)
		}
		events = append(events, result.Events...)
		if result.NextForwardToken == "" || result.NextForwardToken == token {
			return events
		}
		token = result.NextForwardToken
	}
}

func enableBearer(t *testing.T, svc *LogsService, group string) {
	t.Helper()
	if err := svc.putBearerTokenAuthenticationCore(PutBearerTokenAuthenticationInput{
		LogGroupIdentifier:               group,
		BearerTokenAuthenticationEnabled: true,
		BearerTokenAuthenticationSet:     true,
		Region:                           "us-east-1",
	}); err != nil {
		t.Fatalf("enable bearer: %v", err)
	}
}

func TestHTTPIngestionNDJSON(t *testing.T) {
	handler, svc, group, stream := newIngestionEndpointTest(t, nil)
	now := time.Now().UnixMilli()
	body := fmt.Sprintf(
		"{\"timestamp\":%d,\"message\":\"event one\"}\n"+
			"this is not valid json\n"+
			"\"a plain string\"\n"+
			"42\n"+
			"\n"+
			"{\"timestamp\":%d,\"error\":\"timeout\"}\n",
		now-2000, now-1000)
	rec := postIngestion(t, handler, "/ingest/bulk?logGroup="+group+"&logStream="+stream, body,
		map[string]string{"Content-Type": "application/x-ndjson"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("response body: %v", err)
	}
	if _, has := resp["partialSuccess"]; !has {
		t.Fatalf("expected partialSuccess for a request with an invalid line, got %v", resp)
	}
	ps := resp["partialSuccess"].(map[string]interface{})
	if ps["rejectedLogRecords"].(float64) != 1 {
		t.Fatalf("rejectedLogRecords = %v, want 1 (the invalid line)", ps["rejectedLogRecords"])
	}

	events := readIngestedEvents(t, svc, group, stream)
	if len(events) != 4 {
		t.Fatalf("ingested %d events, want 4", len(events))
	}
	messages := map[string]bool{}
	for _, e := range events {
		messages[e.Message] = true
	}
	// Object messages are the object's compact JSON; primitives render
	// their own literal (a bare string without quotes).
	for _, want := range []string{
		fmt.Sprintf(`{"message":"event one","timestamp":%d}`, now-2000),
		"a plain string",
		"42",
		fmt.Sprintf(`{"error":"timeout","timestamp":%d}`, now-1000),
	} {
		if !messages[want] {
			t.Fatalf("message %q missing; ingested %v", want, messages)
		}
	}
	// The two timestamped objects keep their own timestamps.
	if events[0].Timestamp != now-2000 {
		t.Fatalf("first timestamp = %d, want %d", events[0].Timestamp, now-2000)
	}
}

func TestHTTPIngestionNDJSONAllInvalid(t *testing.T) {
	handler, _, group, stream := newIngestionEndpointTest(t, nil)
	rec := postIngestion(t, handler, "/ingest/bulk?logGroup="+group+"&logStream="+stream,
		"not json\nalso not json\n", map[string]string{"Content-Type": "application/x-ndjson"})
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "All events were invalid") {
		t.Fatalf("all-invalid: status %d body %s", rec.Code, rec.Body.String())
	}
}

func TestHTTPIngestionStructuredJSON(t *testing.T) {
	handler, svc, group, stream := newIngestionEndpointTest(t, nil)
	now := time.Now().UnixMilli()
	body := fmt.Sprintf(`[
		{"timestamp":%d,"message":"valid object"},
		"just a string",
		42,
		{"timestamp":%d,"message":"another object"}
	]`, now, now)
	rec := postIngestion(t, handler, "/ingest/json?logGroup="+group+"&logStream="+stream, body,
		map[string]string{"Content-Type": "application/json"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	events := readIngestedEvents(t, svc, group, stream)
	if len(events) != 2 {
		t.Fatalf("ingested %d events, want 2 (objects only)", len(events))
	}
	// Same-timestamp events order by message digest on the read plane, so
	// the assertion is the message set, not the arrival order.
	messages := map[string]bool{}
	for _, e := range events {
		messages[e.Message] = true
	}
	if !messages[fmt.Sprintf(`{"message":"valid object","timestamp":%d}`, now)] ||
		!messages[fmt.Sprintf(`{"message":"another object","timestamp":%d}`, now)] {
		t.Fatalf("ingested messages %v", messages)
	}

	// Concatenated objects parse only the first; a top-level primitive is
	// skipped, producing zero events and an empty success body.
	rec = postIngestion(t, handler, "/ingest/json?logGroup="+group+"&logStream="+stream,
		`{"message":"first"}{"message":"second"}`, map[string]string{"Content-Type": "application/json"})
	if rec.Code != http.StatusOK {
		t.Fatalf("concatenated status %d", rec.Code)
	}
	// A top-level primitive is skipped — counted as one rejected record
	// through the partialSuccess body.
	rec = postIngestion(t, handler, "/ingest/json?logGroup="+group+"&logStream="+stream,
		`"hello"`, map[string]string{"Content-Type": "application/json"})
	if rec.Code != http.StatusOK {
		t.Fatalf("primitive status %d body %s", rec.Code, rec.Body.String())
	}
	var primitiveResp struct {
		PartialSuccess struct {
			RejectedLogRecords int `json:"rejectedLogRecords"`
		} `json:"partialSuccess"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &primitiveResp); err != nil ||
		primitiveResp.PartialSuccess.RejectedLogRecords != 1 {
		t.Fatalf("primitive body %s: want partialSuccess.rejectedLogRecords 1", rec.Body.String())
	}
}

func TestHTTPIngestionHLC(t *testing.T) {
	handler, svc, group, stream := newIngestionEndpointTest(t, nil)
	// Single event, concatenated events, an array, an object without the
	// event member (skipped) and the three time forms — float, integer
	// and numeric string seconds (the documented example times are 2017
	// stamps; the fixtures use current seconds so the age window keeps
	// them).
	base := time.Now().Unix() - 10
	body := fmt.Sprintf(`{"event":"Hello world!","time":%d.5}`, base) +
		fmt.Sprintf(`{"event":"msg2","time":%d}`, base+1) +
		fmt.Sprintf(`{"event":"msg3","time":"%d.25"}`, base+2) +
		`{"message":"no event member"}` +
		`[{"event":"arr1"},{"event":{"structured":true}}]`
	rec := postIngestion(t, handler, "/services/collector/event?logGroup="+group+"&logStream="+stream, body,
		map[string]string{"Content-Type": "application/json"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	events := readIngestedEvents(t, svc, group, stream)
	if len(events) != 5 {
		t.Fatalf("ingested %d events, want 5", len(events))
	}
	if events[0].Message != "Hello world!" || events[0].Timestamp != base*1000+500 {
		t.Fatalf("first event = %q @%d", events[0].Message, events[0].Timestamp)
	}
	if events[1].Timestamp != (base+1)*1000 {
		t.Fatalf("integer seconds event timestamp = %d", events[1].Timestamp)
	}
	if events[2].Timestamp != (base+2)*1000+250 {
		t.Fatalf("string seconds event timestamp = %d", events[2].Timestamp)
	}
	if events[4].Message != `{"structured":true}` {
		t.Fatalf("structured event message = %q", events[4].Message)
	}

	// HLC carries no partial-success vocabulary (the endpoint comparison
	// table): age-rejected records drop without a report.
	oldBody := fmt.Sprintf(`{"event":"ancient","time":%d}`, time.Now().Unix()-15*24*3600)
	rec = postIngestion(t, handler, "/services/collector/event?logGroup="+group+"&logStream="+stream, oldBody,
		map[string]string{"Content-Type": "application/json"})
	if rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != "{}" {
		t.Fatalf("age-rejected HLC status %d body %s", rec.Code, rec.Body.String())
	}
	if events = readIngestedEvents(t, svc, group, stream); len(events) != 5 {
		t.Fatalf("age-rejected HLC ingested %d events, want 5", len(events))
	}
}

func TestHTTPIngestionOTLP(t *testing.T) {
	handler, svc, group, stream := newIngestionEndpointTest(t, nil)
	// The documented example's timeUnixNano is a 2025 stamp — outside the
	// age window — so the fixture pins the form with a current
	// nanosecond stamp.
	nanos := time.Now().Add(-2 * time.Second).UnixNano()
	body := fmt.Sprintf(`{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my-service"}}]},"scopeLogs":[{"scope":{"name":"my-library"},"logRecords":[`+
		`{"timeUnixNano":"%d","severityNumber":9,"severityText":"INFO","body":{"stringValue":"User logged in successfully"},"attributes":[{"key":"user.id","value":{"stringValue":"12345"}}]},`+
		`{"body":{"intValue":"42"}},`+
		`{"severityText":"NOBODY"}`+
		`]}]}]}`, nanos)
	headers := map[string]string{
		"Content-Type":     "application/json",
		"x-aws-log-group":  group,
		"x-aws-log-stream": stream,
	}
	rec := postIngestion(t, handler, "/v1/logs", body, headers)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	// The record without a body is the request's one rejected record.
	var resp struct {
		PartialSuccess struct {
			RejectedLogRecords int `json:"rejectedLogRecords"`
		} `json:"partialSuccess"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil || resp.PartialSuccess.RejectedLogRecords != 1 {
		t.Fatalf("response %s: want partialSuccess.rejectedLogRecords 1", rec.Body.String())
	}
	events := readIngestedEvents(t, svc, group, stream)
	if len(events) != 2 {
		t.Fatalf("ingested %d events, want 2 (the record without a body is skipped)", len(events))
	}
	if events[0].Message != "User logged in successfully" || events[0].Timestamp != nanos/1e6 {
		t.Fatalf("first event = %q @%d", events[0].Message, events[0].Timestamp)
	}
	if events[1].Message != "42" {
		t.Fatalf("intValue event message = %q", events[1].Message)
	}

	// OTLP is headers-only: query parameters are not supported, so a
	// request addressing the target by query parameter alone lacks it.
	rec = postIngestion(t, handler, "/v1/logs?logGroup="+group+"&logStream="+stream, "{}",
		map[string]string{"Content-Type": "application/json"})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("query-parameter OTLP status %d, want 400", rec.Code)
	}

	// The protobuf encoding is not implemented on this platform: it
	// rejects rather than accepting and dropping.
	rec = postIngestion(t, handler, "/v1/logs", "", map[string]string{
		"Content-Type":     "application/x-protobuf",
		"x-aws-log-group":  group,
		"x-aws-log-stream": stream,
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("protobuf OTLP status %d, want 400", rec.Code)
	}
}

func TestHTTPIngestionBearerAuthentication(t *testing.T) {
	handler, svc, group, stream := newIngestionEndpointTest(t, &stubIngestionAuthenticator{valid: map[string]bool{"good-token": true}})
	path := "/ingest/bulk?logGroup=" + group + "&logStream=" + stream
	const body = `{"message":"bearer hello"}`

	// An invalid token answers 401 before anything else.
	rec := postIngestion(t, handler, path, body, map[string]string{"Authorization": "Bearer bad-token", "Content-Type": "application/json"})
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("invalid token status %d", rec.Code)
	}

	// A valid token against a group without the switch answers 403: the
	// documented rule that bearer token authentication must be enabled on
	// the log group.
	rec = postIngestion(t, handler, path, body, map[string]string{"Authorization": "Bearer good-token", "Content-Type": "application/json"})
	if rec.Code != http.StatusForbidden {
		t.Fatalf("flag-disabled status %d", rec.Code)
	}

	// With the switch on, the same request ingests.
	enableBearer(t, svc, group)
	rec = postIngestion(t, handler, path, body, map[string]string{"Authorization": "Bearer good-token", "Content-Type": "application/json"})
	if rec.Code != http.StatusOK {
		t.Fatalf("enabled status %d: %s", rec.Code, rec.Body.String())
	}
	if events := readIngestedEvents(t, svc, group, stream); len(events) != 1 || events[0].Message != `{"message":"bearer hello"}` {
		t.Fatalf("ingested %v", events)
	}

	// Disabling the switch closes the endpoint again.
	if err := svc.putBearerTokenAuthenticationCore(PutBearerTokenAuthenticationInput{
		LogGroupIdentifier:               group,
		BearerTokenAuthenticationEnabled: false,
		BearerTokenAuthenticationSet:     true,
		Region:                           "us-east-1",
	}); err != nil {
		t.Fatalf("disable bearer: %v", err)
	}
	rec = postIngestion(t, handler, path, body, map[string]string{"Authorization": "Bearer good-token", "Content-Type": "application/json"})
	if rec.Code != http.StatusForbidden {
		t.Fatalf("re-disabled status %d", rec.Code)
	}
}

func TestHTTPIngestionTargetRules(t *testing.T) {
	handler, _, group, stream := newIngestionEndpointTest(t, nil)
	const body = `{"message":"x"}`

	// A parameter supplied both ways rejects; each missing parameter
	// rejects; an unknown group or stream answers 404.
	rec := postIngestion(t, handler, "/ingest/bulk?logGroup="+group+"&logStream="+stream, body,
		map[string]string{"x-aws-log-group": group, "Content-Type": "application/json"})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("both-ways logGroup status %d", rec.Code)
	}
	rec = postIngestion(t, handler, "/ingest/bulk?logStream="+stream, body, nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("missing logGroup status %d", rec.Code)
	}
	rec = postIngestion(t, handler, "/ingest/bulk?logGroup=no-such-group&logStream="+stream, body, nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("unknown group status %d", rec.Code)
	}
	rec = postIngestion(t, handler, "/ingest/bulk?logGroup="+group+"&logStream=no-such-stream", body, nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("unknown stream status %d", rec.Code)
	}

	// Header addressing works on the query-parameter endpoints too.
	rec = postIngestion(t, handler, "/ingest/bulk", body, map[string]string{
		"x-aws-log-group": group, "x-aws-log-stream": stream, "Content-Type": "application/json"})
	if rec.Code != http.StatusOK {
		t.Fatalf("header addressing status %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHTTPIngestionLimitsAndPartialSuccess(t *testing.T) {
	handler, svc, group, stream := newIngestionEndpointTest(t, nil)
	path := "/ingest/json?logGroup=" + group + "&logStream=" + stream

	// A request over 1 MB rejects.
	oversized := "[" + strings.Repeat(`{"message":"pad"}`, 210000) + "]"
	rec := postIngestion(t, handler, path, oversized, nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("oversized request status %d", rec.Code)
	}

	// An event over 256 KB rejects.
	large := `{"message":"` + strings.Repeat("x", 300*1024) + `"}`
	rec = postIngestion(t, handler, path, large, nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("oversized event status %d", rec.Code)
	}

	// Age rejections report through the documented partialSuccess body:
	// one too old (past the 14-day window), one too new (beyond two
	// hours), one valid.
	now := time.Now().UnixMilli()
	body := fmt.Sprintf(`[{"timestamp":%d,"message":"old"},{"timestamp":%d,"message":"new"},{"timestamp":%d,"message":"ok"}]`,
		now-15*24*3600*1000, now+3*3600*1000, now-1000)
	rec = postIngestion(t, handler, path, body, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("partial status %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		PartialSuccess struct {
			RejectedLogRecords int    `json:"rejectedLogRecords"`
			ErrorMessage       string `json:"errorMessage"`
		} `json:"partialSuccess"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("partial body: %v", err)
	}
	if resp.PartialSuccess.RejectedLogRecords != 2 {
		t.Fatalf("rejectedLogRecords = %d, want 2", resp.PartialSuccess.RejectedLogRecords)
	}
	var breakdown struct {
		TooOld int `json:"tooOldLogEventCount"`
		TooNew int `json:"tooNewLogEventCount"`
	}
	if err := json.Unmarshal([]byte(resp.PartialSuccess.ErrorMessage), &breakdown); err != nil {
		t.Fatalf("breakdown: %v", err)
	}
	if breakdown.TooOld != 1 || breakdown.TooNew != 1 {
		t.Fatalf("breakdown = %+v, want tooOld=1 tooNew=1", breakdown)
	}
	events := readIngestedEvents(t, svc, group, stream)
	if len(events) != 1 || events[0].Message != `{"message":"ok","timestamp":`+fmt.Sprint(now-1000)+`}` {
		t.Fatalf("ingested after partial: %v", events)
	}

	// A request within the 1 MB body limit whose message bytes plus the
	// seam's per-event accounting overflow the PutLogEvents batch bound
	// still ingests whole: the handler chunks at the seam's boundary
	// (1,000 events of 1,032 message bytes: body 1,048,000, accounting
	// 1,058,000). Every message is unique: the read engine orders by the
	// message digest when timestamp and ingestion time tie, so identical
	// messages would share one pagination coordinate.
	var chunkLines []string
	for i := 0; i < 1000; i++ {
		chunkLines = append(chunkLines, fmt.Sprintf(`{"message":"%04d:%s"}`, i, strings.Repeat("y", 1026)))
	}
	rec = postIngestion(t, handler, "/ingest/bulk?logGroup="+group+"&logStream="+stream,
		strings.Join(chunkLines, "\n"), nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("chunked batch status %d: %s", rec.Code, rec.Body.String())
	}
	all := readIngestedEvents(t, svc, group, stream)
	if len(all) != 1001 { // 1 from the partial batch + 1,000 chunked
		t.Fatalf("total ingested %d, want 1001", len(all))
	}
}

func TestHTTPIngestionRegionFromHost(t *testing.T) {
	// The logs.<region> host form selects the region; a non-matching host
	// falls back to the platform default region. The second store (the
	// default region's) stays empty for a request addressed to
	// us-west-2.
	handler, svc, group, stream := newIngestionEndpointTest(t, nil)
	req := httptest.NewRequest(http.MethodPost, "/ingest/bulk?logGroup="+group+"&logStream="+stream,
		bytes.NewReader([]byte(`{"message":"west"}`)))
	req.Host = "logs.us-west-2.amazonaws.com"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("us-west-2 request status %d, want 404 (the group lives in us-east-1)", rec.Code)
	}
	west, err := svc.getLogsStoreByRegion("us-west-2")
	if err != nil {
		t.Fatalf("west store: %v", err)
	}
	if _, err := west.GetLogGroup(group); err == nil {
		t.Fatal("the us-west-2 store must not carry the us-east-1 group")
	}
}

// A request whose records span more than 24 hours is chunked at the span
// boundary: each line is an independent record, and a whole-request 400
// after earlier chunks already committed would both break that contract
// and duplicate the committed prefix on retry.
func TestHTTPIngestionSpansBeyondOneDay(t *testing.T) {
	handler, svc, group, stream := newIngestionEndpointTest(t, nil)
	now := time.Now().UnixMilli()
	tenDays := int64(10 * 24 * time.Hour / time.Millisecond)
	body := fmt.Sprintf("{\"timestamp\":%d,\"message\":\"old record\"}\n{\"timestamp\":%d,\"message\":\"new record\"}\n",
		now-tenDays, now)
	rec := postIngestion(t, handler, "/ingest/bulk?logGroup="+group+"&logStream="+stream, body,
		map[string]string{"Content-Type": "application/x-ndjson"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	events := readIngestedEvents(t, svc, group, stream)
	if len(events) != 2 {
		t.Fatalf("a >24h-span request must ingest both records, got %d", len(events))
	}
}

// The HLC batched request body carries the array inside a wrapper
// object's "event" member ("The request body should be in JSON format
// with an array of events", the HLC page's Request-format section):
// each element is its own event, on the element's own "time" when it
// carries one and the wrapper's otherwise.
func TestHTTPIngestionHLCBatchedWrapper(t *testing.T) {
	handler, svc, group, stream := newIngestionEndpointTest(t, nil)
	base := time.Now().Unix() - 10
	body := fmt.Sprintf(`{"time":%d,"event":[`, base) +
		fmt.Sprintf(`{"event":"a","time":%d},`, base+1) +
		`{"event":{"structured":2}},` +
		`"bare-string-element"` +
		`]}`
	rec := postIngestion(t, handler, "/services/collector/event?logGroup="+group+"&logStream="+stream, body,
		map[string]string{"Content-Type": "application/json"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	events := readIngestedEvents(t, svc, group, stream)
	if len(events) != 3 {
		t.Fatalf("ingested %d events, want 3 (one per array element)", len(events))
	}
	timestamps := map[string]int64{}
	for _, evt := range events {
		timestamps[evt.Message] = evt.Timestamp
	}
	// The wrapper-time elements share the wrapper's timestamp; the
	// own-time element carries its own.
	if timestamps["bare-string-element"] != base*1000 || timestamps[`{"structured":2}`] != base*1000 {
		t.Fatalf("wrapper-time elements = %v", timestamps)
	}
	if timestamps["a"] != (base+1)*1000 {
		t.Fatalf("own-time element = %v", timestamps)
	}
}

// Content-Type matches on the media type alone: a charset parameter is
// standard content negotiation and does not reject the request, while a
// different media type still does.
func TestHTTPIngestionContentTypeParameters(t *testing.T) {
	handler, _, group, stream := newIngestionEndpointTest(t, nil)
	query := "?logGroup=" + group + "&logStream=" + stream

	for _, path := range []string{"/ingest/json", "/ingest/bulk"} {
		if rec := postIngestion(t, handler, path+query, `{"x":1}`,
			map[string]string{"Content-Type": "application/json; charset=utf-8"}); rec.Code == http.StatusBadRequest {
			t.Errorf("%s: charset parameter rejected: %s", path, rec.Body.String())
		}
	}
	// The OTLP endpoint rides the same media-type parse; its header-only
	// target resolution stays outside this pin.
	if mt := mediaTypeOf("application/json; charset=utf-8"); mt != "application/json" {
		t.Errorf("mediaTypeOf charset: %q", mt)
	}
	if mt := mediaTypeOf("APPLICATION/JSON ; v=1"); mt != "application/json" {
		t.Errorf("mediaTypeOf normalisation: %q", mt)
	}
	if rec := postIngestion(t, handler, "/ingest/bulk"+query, `{"x":1}`,
		map[string]string{"Content-Type": "application/x-ndjson; charset=utf-8"}); rec.Code == http.StatusBadRequest {
		t.Errorf("ndjson charset parameter rejected: %s", rec.Body.String())
	}
	if rec := postIngestion(t, handler, "/ingest/json"+query, `{"x":1}`,
		map[string]string{"Content-Type": "text/plain"}); rec.Code != http.StatusBadRequest {
		t.Errorf("wrong media type accepted: %d", rec.Code)
	}
}

// The OTLP AnyValue edges the JSON encoding alone expresses: non-finite
// doubles travel as doubleValue's string form (JSON numbers cannot carry
// NaN or the infinities) and bytesValue is the base64 rendering of the
// value's bytes — both render their value instead of skipping the record.
func TestHTTPIngestionOTLPAnyValueEdges(t *testing.T) {
	handler, svc, group, stream := newIngestionEndpointTest(t, nil)
	nanos := time.Now().Add(-2 * time.Second).UnixNano()
	body := fmt.Sprintf(`{"resourceLogs":[{"scopeLogs":[{"logRecords":[`+
		`{"timeUnixNano":"%d","body":{"doubleValue":"NaN"}},`+
		`{"timeUnixNano":"%d","body":{"doubleValue":"Infinity"}},`+
		`{"timeUnixNano":"%d","body":{"bytesValue":"aGVsbG8="}}`+
		`]}]}]}`, nanos, nanos, nanos)
	rec := postIngestion(t, handler, "/v1/logs", body, map[string]string{
		"Content-Type":     "application/json",
		"x-aws-log-group":  group,
		"x-aws-log-stream": stream,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		PartialSuccess struct {
			RejectedLogRecords int `json:"rejectedLogRecords"`
		} `json:"partialSuccess"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil || resp.PartialSuccess.RejectedLogRecords != 0 {
		t.Fatalf("response %s: want partialSuccess.rejectedLogRecords 0", rec.Body.String())
	}
	events := readIngestedEvents(t, svc, group, stream)
	if len(events) != 3 {
		t.Fatalf("ingested %d events, want 3", len(events))
	}
	// The three records share one timestamp, so the store's tiebreak
	// order — not the body order — serves them; the rendered set is the
	// contract under test.
	got := []string{events[0].Message, events[1].Message, events[2].Message}
	sort.Strings(got)
	want := []string{"NaN", "Infinity", "hello"}
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("edge messages = %v, want %v", got, want)
	}
}
