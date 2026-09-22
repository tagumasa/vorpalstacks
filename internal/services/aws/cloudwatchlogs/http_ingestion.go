package cloudwatchlogs

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/defaults"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The plain-HTTP log ingestion endpoints (developer guide, "Log ingestion
// through HTTP endpoints"): four request shapes — OTLP (/v1/logs), HLC
// (/services/collector/event), ND-JSON (/ingest/bulk) and Structured JSON
// (/ingest/json) — that write through the single ingestion seam
// (putLogEventsCore), so HTTP-ingested events carry the same validation,
// transformer application and filter fan-out as every other write path.
//
// Authentication is either SigV4 (verified by the platform's signature
// middleware when enabled) or an ACWL bearer token: an IAM
// service-specific credential of service logs.amazonaws.com whose
// ServiceCredentialSecret travels in the Authorization: Bearer header. The
// documented gate applies to the bearer leg alone: bearer token
// authentication must be enabled on the log group (PutBearerToken
// Authentication) before the group accepts bearer-authenticated logs.

// IngestionTokenAuthenticator validates an ACWL bearer token presented to
// the HTTP ingestion endpoints. The composition root injects the IAM
// service implementation; a nil authenticator fails every token closed.
type IngestionTokenAuthenticator interface {
	AuthenticateLogIngestionToken(token string) error
}

// httpIngestionHandler serves the four ingestion paths.
type httpIngestionHandler struct {
	svc           *LogsService
	authenticator IngestionTokenAuthenticator
}

// NewHTTPIngestionHandler builds the plain-HTTP ingestion endpoints. The
// returned handler is registered on the platform HTTP server ahead of the
// classified AWS API catch-all.
func NewHTTPIngestionHandler(svc *LogsService, authenticator IngestionTokenAuthenticator) http.Handler {
	h := &httpIngestionHandler{svc: svc, authenticator: authenticator}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/logs", h.serveOTLP)
	mux.HandleFunc("/services/collector/event", h.serveHLC)
	mux.HandleFunc("/ingest/bulk", h.serveNDJSON)
	mux.HandleFunc("/ingest/json", h.serveStructuredJSON)
	return mux
}

// IngestionEndpointPaths are the four ingestion route paths, exported for
// the server's route registration and signature-verification middleware.
var IngestionEndpointPaths = []string{
	"/v1/logs",
	"/services/collector/event",
	"/ingest/bulk",
	"/ingest/json",
}

// HasBearerScheme reports whether an Authorization header carries the
// bearer auth-scheme. RFC 7235 compares auth-schemes case-insensitively,
// so "bearer", "Bearer" and "BEARER" are one scheme; the signature
// middleware and the ingestion handler share this predicate.
func HasBearerScheme(auth string) bool {
	scheme, _, _ := strings.Cut(auth, " ")
	return strings.EqualFold(scheme, "Bearer")
}

// IsIngestionEndpointPath reports whether a request path is one of the
// plain-HTTP ingestion endpoints — a platform endpoint, not a classified
// AWS API call.
func IsIngestionEndpointPath(path string) bool {
	for _, p := range IngestionEndpointPaths {
		if path == p {
			return true
		}
	}
	return false
}

// ingestedRecord is one parsed request record before the seam.
type ingestedRecord struct {
	timestamp int64 // epoch milliseconds
	message   string
}

// ingestionCounts carries the per-reason rejection accounting of one
// request. Parse skips join rejectedLogRecords alongside the age reasons
// the errorMessage breakdown names.
type ingestionCounts struct {
	skipped int // records the format's rules skip (invalid JSON, wrong shape)
	tooOld  int
	tooNew  int
	expired int
}

func (c *ingestionCounts) rejected() int {
	return c.skipped + c.tooOld + c.tooNew + c.expired
}

// ingestionTarget is the request's resolved destination.
type ingestionTarget struct {
	group  string
	stream string
	region string
}

// serveIngestion is the shared endpoint body: authentication, target
// resolution, existence and flag checks, body limits, then the format
// parser and the seam. parse must consume the body and return the parsed
// records; format-level rejects answer 400 themselves.
func (h *httpIngestionHandler) serveIngestion(w http.ResponseWriter, r *http.Request, allowQueryParams bool, parse func(body []byte) ([]ingestedRecord, *ingestionCounts, *ingestionError)) {
	token, bearer, err := bearerTokenOf(r)
	if err != nil {
		writeIngestionJSON(w, http.StatusUnauthorized, map[string]string{"message": err.Error()})
		return
	}
	if bearer {
		if h.authenticator == nil {
			writeIngestionJSON(w, http.StatusUnauthorized, map[string]string{"message": "bearer token authentication is unavailable"})
			return
		}
		if err := h.authenticator.AuthenticateLogIngestionToken(token); err != nil {
			writeIngestionJSON(w, http.StatusUnauthorized, map[string]string{"message": err.Error()})
			return
		}
	}

	target, errCode, errMsg := h.resolveTarget(r, allowQueryParams)
	if errCode != 0 {
		writeIngestionJSON(w, errCode, map[string]string{"message": errMsg})
		return
	}

	store, err := h.svc.getLogsStoreByRegion(target.region)
	if err != nil {
		writeIngestionJSON(w, http.StatusInternalServerError, map[string]string{"message": "storage unavailable"})
		return
	}
	group, err := store.GetLogGroup(target.group)
	if err != nil {
		writeIngestionJSON(w, http.StatusNotFound, map[string]string{"message": "Log group does not exist"})
		return
	}
	if _, err := store.GetLogStream(target.group, target.stream); err != nil {
		writeIngestionJSON(w, http.StatusNotFound, map[string]string{"message": "Log stream does not exist"})
		return
	}
	if bearer {
		// The switch rides the group record the gate already holds.
		if !group.BearerTokenAuthenticationEnabled {
			writeIngestionJSON(w, http.StatusForbidden, map[string]string{"message": "Bearer token authentication is not enabled on the log group"})
			return
		}
	}

	body, errCode, errMsg := readIngestionBody(r)
	if errCode != 0 {
		writeIngestionJSON(w, errCode, map[string]string{"message": errMsg})
		return
	}

	records, counts, perr := parse(body)
	if perr != nil {
		writeIngestionJSON(w, perr.code, map[string]string{"message": perr.message})
		return
	}

	for _, rec := range records {
		if len(rec.message) > logsstore.HTTPIngestionMaxEventBytes {
			writeIngestionJSON(w, http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("Each log event can be no larger than %d bytes", logsstore.HTTPIngestionMaxEventBytes)})
			return
		}
	}
	if len(records) > logsstore.HTTPIngestionMaxEvents {
		writeIngestionJSON(w, http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("Maximum number of events per request is %d", logsstore.HTTPIngestionMaxEvents)})
		return
	}

	// The endpoints' rejection vocabulary has no ordering reason — each
	// line is an independent record — so the platform accepts events in
	// any arrival order and sorts stably by timestamp before the seam.
	sort.SliceStable(records, func(i, j int) bool { return records[i].timestamp < records[j].timestamp })

	// The destination group's retention horizon rejects records that
	// precede it, the ingestion plane's form of the PutLogEvents rule
	// ("Events older than 14 days or preceding the log group's retention
	// period are rejected while processing remaining valid events"). The
	// horizon rides the group record the gate already holds.
	retentionCutoff := int64(0)
	if group.RetentionInDays > 0 {
		retentionCutoff = time.Now().UnixMilli() - int64(group.RetentionInDays)*24*60*60*1000
	}
	now := time.Now().UnixMilli()
	valid := make([]ingestedRecord, 0, len(records))
	for _, rec := range records {
		if rec.timestamp > now+tooNewThreshold {
			counts.tooNew++
			continue
		}
		if rec.timestamp < now-tooOldThreshold {
			counts.tooOld++
			continue
		}
		if retentionCutoff > now-tooOldThreshold && rec.timestamp < retentionCutoff {
			counts.expired++
			continue
		}
		valid = append(valid, rec)
	}

	if err := h.ingestRecords(store, target, valid); err != nil {
		writeIngestionJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	if counts.rejected() == 0 || !h.formatSupportsPartialSuccess(r.URL.Path) {
		writeIngestionJSON(w, http.StatusOK, map[string]interface{}{})
		return
	}
	breakdown, _ := json.Marshal(map[string]int{
		"tooOldLogEventCount":  counts.tooOld,
		"tooNewLogEventCount":  counts.tooNew,
		"expiredLogEventCount": counts.expired,
	})
	writeIngestionJSON(w, http.StatusOK, map[string]interface{}{
		"partialSuccess": map[string]interface{}{
			"rejectedLogRecords": counts.rejected(),
			"errorMessage":       string(breakdown),
		},
	})
}

// formatSupportsPartialSuccess: the endpoint comparison table gives
// ND-JSON, Structured JSON and OTLP a partial-success response; HLC has
// none, so its skipped and age-rejected records are dropped without a
// report (recorded adjudication — the format documents no accounting of
// its own).
func (h *httpIngestionHandler) formatSupportsPartialSuccess(path string) bool {
	return path != "/services/collector/event"
}

// ingestRecords writes the valid records through the single ingestion
// seam. The HTTP request is not a PutLogEvents batch: a documented 1 MB
// request of maximal events can exceed the seam's per-batch byte
// accounting, so the handler chunks at that accounting boundary.
func (h *httpIngestionHandler) ingestRecords(store *logsstore.Store, target ingestionTarget, records []ingestedRecord) error {
	// The seam rejects an empty message and a negative timestamp
	// whole-request; a record reaching either rejection after an earlier
	// chunk flushed would 400 the request with part of it persisted —
	// the same mid-commit hazard the span pre-check inside the loop
	// exists to prevent — so both are checked before any flush, with
	// the seam's own vocabulary.
	for _, rec := range records {
		if rec.message == "" {
			return NewLogsError("InvalidParameterException",
				"The message member of every log event must be at least 1 character", 400)
		}
		if rec.timestamp < 0 {
			return NewLogsError("InvalidParameterException",
				"The timestamp member of every log event must not be negative", 400)
		}
	}
	batch := make([]PutLogEvent, 0, len(records))
	batchBytes := 0
	for _, rec := range records {
		eventBytes := len(rec.message) + logsstore.PutLogEventOverheadBytes
		// Records are timestamp-sorted, so a candidate's span over the
		// batch is its timestamp minus the batch's first event. The seam
		// rejects spans beyond 24 hours, and a chunk rejected there would
		// 400 the request after earlier chunks already committed — the
		// span is a chunking dimension, so every emitted chunk satisfies
		// the seam's contract and "each line is an independent record"
		// holds for requests spanning more than a day.
		if len(batch) > 0 && (batchBytes+eventBytes > logsstore.MaxPutLogEventsBatchBytes ||
			len(batch) >= logsstore.MaxBatchLogEvents ||
			rec.timestamp-batch[0].Timestamp > maxEventsTimeSpan) {
			if err := h.flushIngestionBatch(target, batch); err != nil {
				return err
			}
			batch = batch[:0]
			batchBytes = 0
		}
		batch = append(batch, PutLogEvent{
			LogEntry:     logsstore.LogEntry{Timestamp: rec.timestamp, Message: rec.message},
			TimestampSet: true,
		})
		batchBytes += eventBytes
	}
	if len(batch) > 0 {
		return h.flushIngestionBatch(target, batch)
	}
	return nil
}

func (h *httpIngestionHandler) flushIngestionBatch(target ingestionTarget, batch []PutLogEvent) error {
	_, err := h.svc.putLogEventsCore(PutLogEventsInput{
		LogGroupName:  target.group,
		LogStreamName: target.stream,
		Events:        batch,
		Region:        target.region,
	})
	return err
}

// resolveTarget reads the log group and stream: the x-aws-log-group and
// x-aws-log-stream headers, or the logGroup and logStream query parameters
// (every endpoint except OTLP). A parameter supplied both ways rejects;
// the region comes from the logs.<region> host form, defaulting to the
// platform region.
func (h *httpIngestionHandler) resolveTarget(r *http.Request, allowQueryParams bool) (ingestionTarget, int, string) {
	headerGroup := r.Header.Get("x-aws-log-group")
	headerStream := r.Header.Get("x-aws-log-stream")
	queryGroup := r.URL.Query().Get("logGroup")
	queryStream := r.URL.Query().Get("logStream")
	if !allowQueryParams {
		queryGroup, queryStream = "", ""
	}
	if headerGroup != "" && queryGroup != "" {
		return ingestionTarget{}, http.StatusBadRequest, "logGroup was provided in both a header and a query parameter; use one"
	}
	if headerStream != "" && queryStream != "" {
		return ingestionTarget{}, http.StatusBadRequest, "logStream was provided in both a header and a query parameter; use one"
	}
	group := headerGroup
	if group == "" {
		group = queryGroup
	}
	stream := headerStream
	if stream == "" {
		stream = queryStream
	}
	if group == "" {
		return ingestionTarget{}, http.StatusBadRequest, "logGroup is required (x-aws-log-group header or logGroup query parameter)"
	}
	if stream == "" {
		return ingestionTarget{}, http.StatusBadRequest, "logStream is required (x-aws-log-stream header or logStream query parameter)"
	}
	return ingestionTarget{group: group, stream: stream, region: ingestionRegion(r.Host)}, 0, ""
}

// ingestionRegion derives the region from the endpoint host form
// logs.<region>.amazonaws.com; a host of any other shape (a local
// deployment address, say) falls back to the platform's default region.
func ingestionRegion(host string) string {
	host = strings.SplitN(host, ":", 2)[0]
	parts := strings.Split(host, ".")
	if len(parts) >= 3 && parts[0] == "logs" && parts[1] != "" {
		return parts[1]
	}
	return defaults.DefaultRegion
}

// bearerTokenOf extracts the Authorization: Bearer token. The scheme is
// matched case-insensitively (RFC 7235 auth-scheme comparison); an
// Authorization header of any other shape (SigV4, none) is the non-bearer
// leg: verification is the signature middleware's when enabled, and this
// handler treats the request as authenticated on that plane.
func bearerTokenOf(r *http.Request) (string, bool, error) {
	auth := r.Header.Get("Authorization")
	if !HasBearerScheme(auth) {
		return "", false, nil
	}
	_, rest, _ := strings.Cut(auth, " ")
	token := strings.TrimSpace(rest)
	if token == "" {
		return "", true, errors.New("missing bearer token")
	}
	return token, true, nil
}

func readIngestionBody(r *http.Request) ([]byte, int, string) {
	// The reader is capped before the read (and an advertised
	// Content-Length over the ceiling rejects without reading): a body
	// beyond the documented request size must not be buffered whole to be
	// rejected.
	if r.ContentLength > logsstore.HTTPIngestionMaxRequestBytes {
		return nil, http.StatusBadRequest, fmt.Sprintf("Maximum request size is %d bytes", logsstore.HTTPIngestionMaxRequestBytes)
	}
	r.Body = http.MaxBytesReader(nil, r.Body, logsstore.HTTPIngestionMaxRequestBytes+1)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusBadRequest, "failed to read the request body"
	}
	if len(body) > logsstore.HTTPIngestionMaxRequestBytes {
		return nil, http.StatusBadRequest, fmt.Sprintf("Maximum request size is %d bytes", logsstore.HTTPIngestionMaxRequestBytes)
	}
	return body, 0, ""
}

type ingestionError struct {
	code    int
	message string
}

func writeIngestionJSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

// mediaTypeOf parses a Content-Type header's media type — the type/
// subtype alone, before any parameter (charset and friends are standard
// content negotiation, RFC 9110 §8.3; the endpoints document the media
// type they accept, not an exact-header match).
func mediaTypeOf(contentType string) string {
	mediaType, _, _ := strings.Cut(strings.TrimSpace(strings.ToLower(contentType)), ";")
	return strings.TrimSpace(mediaType)
}

// --- ND-JSON endpoint (/ingest/bulk) ---

func (h *httpIngestionHandler) serveNDJSON(w http.ResponseWriter, r *http.Request) {
	if ct := mediaTypeOf(r.Header.Get("Content-Type")); ct != "" && ct != "application/json" && ct != "application/x-ndjson" {
		writeIngestionJSON(w, http.StatusBadRequest, map[string]string{"message": "Content-Type must be application/json or application/x-ndjson"})
		return
	}
	h.serveIngestion(w, r, true, parseNDJSON)
}

// parseNDJSON parses one JSON value per line. Objects and flattened array
// elements carry their own timestamp; every other value becomes an event
// with the server's time. Lines that are not valid JSON are skipped and
// counted; empty lines are ignored; a request whose lines are all invalid
// rejects.
func parseNDJSON(body []byte) ([]ingestedRecord, *ingestionCounts, *ingestionError) {
	counts := &ingestionCounts{}
	var records []ingestedRecord
	lines := 0
	for _, line := range bytes.Split(body, []byte("\n")) {
		line = bytes.TrimSuffix(line, []byte("\r"))
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		lines++
		var value interface{}
		if err := json.Unmarshal(line, &value); err != nil {
			counts.skipped++
			continue
		}
		records = appendNDJSONValue(records, value, counts)
	}
	if lines > 0 && len(records) == 0 && counts.skipped == lines {
		return nil, nil, &ingestionError{code: http.StatusBadRequest, message: "All events were invalid"}
	}
	return records, counts, nil
}

// appendNDJSONValue maps one ND-JSON value to events: an object is one
// event (timestamp member if present and numeric), an array's elements
// each become one event, and a primitive is an event on the server's
// clock. The message of a structured value is its compact JSON
// serialisation; a bare string is the string itself.
func appendNDJSONValue(records []ingestedRecord, value interface{}, counts *ingestionCounts) []ingestedRecord {
	switch v := value.(type) {
	case map[string]interface{}:
		return append(records, objectRecord(v, "timestamp"))
	case []interface{}:
		for _, element := range v {
			records = appendNDJSONValue(records, element, counts)
		}
		return records
	default:
		return append(records, primitiveRecord(value))
	}
}

// --- Structured JSON endpoint (/ingest/json) ---

// requireJSONContentType answers 400 unless the request's Content-Type
// is empty or application/json — the shared guard of the two JSON-body
// ingestion endpoints (the media type alone; parameters pass). The
// return reports whether handling may proceed.
func requireJSONContentType(w http.ResponseWriter, r *http.Request) bool {
	if ct := mediaTypeOf(r.Header.Get("Content-Type")); ct != "" && ct != "application/json" {
		writeIngestionJSON(w, http.StatusBadRequest, map[string]string{"message": "Content-Type must be application/json"})
		return false
	}
	return true
}

func (h *httpIngestionHandler) serveStructuredJSON(w http.ResponseWriter, r *http.Request) {
	if !requireJSONContentType(w, r) {
		return
	}
	h.serveIngestion(w, r, true, parseStructuredJSON)
}

// parseStructuredJSON accepts a single JSON object or an array of
// objects. Non-object values are skipped; concatenated objects parse
// only the first (the decoder stops after one value) — the recorded
// platform contract of the /ingest endpoints, pinned by test.
func parseStructuredJSON(body []byte) ([]ingestedRecord, *ingestionCounts, *ingestionError) {
	counts := &ingestionCounts{}
	if len(bytes.TrimSpace(body)) == 0 {
		return nil, nil, &ingestionError{code: http.StatusBadRequest, message: "Request body must be a JSON object or an array of objects"}
	}
	var value interface{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&value); err != nil {
		return nil, nil, &ingestionError{code: http.StatusBadRequest, message: "Request body is not valid JSON"}
	}
	var records []ingestedRecord
	switch v := value.(type) {
	case map[string]interface{}:
		records = append(records, objectRecord(v, "timestamp"))
	case []interface{}:
		for _, element := range v {
			if obj, ok := element.(map[string]interface{}); ok {
				records = append(records, objectRecord(obj, "timestamp"))
			} else {
				counts.skipped++
			}
		}
	default:
		counts.skipped++
	}
	return records, counts, nil
}

// --- HLC endpoint (/services/collector/event) ---

func (h *httpIngestionHandler) serveHLC(w http.ResponseWriter, r *http.Request) {
	if !requireJSONContentType(w, r) {
		return
	}
	h.serveIngestion(w, r, true, parseHLC)
}

// parseHLC accepts a single event object, an array of event objects, or
// concatenated event objects. Each object must carry an "event" member of
// any JSON type (its value is the message; objects without the member are
// skipped); the optional "time" member is epoch seconds — number or
// numeric string, fractional allowed. Primitives are skipped. The
// entityName/entityEnvironment query parameters are accepted and ignored:
// entity association is the Explore related-telemetry surface, which the
// platform does not carry (recorded adjudication).
func parseHLC(body []byte) ([]ingestedRecord, *ingestionCounts, *ingestionError) {
	counts := &ingestionCounts{}
	decoder := json.NewDecoder(bytes.NewReader(body))
	var records []ingestedRecord
	for {
		var value interface{}
		err := decoder.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, &ingestionError{code: http.StatusBadRequest, message: "Request body is not valid JSON"}
		}
		switch v := value.(type) {
		case map[string]interface{}:
			records = appendHLCObject(records, v, counts)
		case []interface{}:
			for _, element := range v {
				if obj, ok := element.(map[string]interface{}); ok {
					records = appendHLCObject(records, obj, counts)
				} else {
					counts.skipped++
				}
			}
		default:
			counts.skipped++
		}
	}
	return records, counts, nil
}

func appendHLCObject(records []ingestedRecord, obj map[string]interface{}, counts *ingestionCounts) []ingestedRecord {
	return appendHLCEvent(records, obj, counts, time.Now().UnixMilli())
}

// appendHLCEvent maps one HLC wrapper object to events with a fallback
// timestamp for members and elements that carry no "time" of their own.
func appendHLCEvent(records []ingestedRecord, obj map[string]interface{}, counts *ingestionCounts, fallbackMs int64) []ingestedRecord {
	event, ok := obj["event"]
	if !ok {
		counts.skipped++
		return records
	}
	timestamp := hlcTimestampOr(obj["time"], fallbackMs)
	switch v := event.(type) {
	case string:
		return append(records, ingestedRecord{timestamp: timestamp, message: v})
	case []interface{}:
		// The batched request body carries the array in a wrapper
		// object's "event" member ("The request body should be in JSON
		// format with an array of events", the HLC page's Request-format
		// section): each element is its own event, with the element's
		// own "time" when it carries one and the wrapper's otherwise.
		for _, element := range v {
			switch e := element.(type) {
			case string:
				records = append(records, ingestedRecord{timestamp: timestamp, message: e})
			case map[string]interface{}:
				records = appendHLCEvent(records, e, counts, timestamp)
			default:
				counts.skipped++
			}
		}
		return records
	default:
		message, err := json.Marshal(event)
		if err != nil || len(message) == 0 {
			counts.skipped++
			return records
		}
		return append(records, ingestedRecord{timestamp: timestamp, message: string(message)})
	}
}

// hlcTimestampOr converts the HLC "time" member — epoch seconds as a
// number or a numeric string, fractional part allowed — to milliseconds;
// anything absent or unparseable reads as the fallback (the server's
// clock at the wrapper level, the wrapper object's time for a batched
// element that carries none).
func hlcTimestampOr(value interface{}, fallbackMs int64) int64 {
	var seconds float64
	switch v := value.(type) {
	case float64:
		seconds = v
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fallbackMs
		}
		seconds = parsed
	default:
		return fallbackMs
	}
	return int64(seconds * 1000)
}

// --- OTLP endpoint (/v1/logs) ---

func (h *httpIngestionHandler) serveOTLP(w http.ResponseWriter, r *http.Request) {
	if mediaType := mediaTypeOf(r.Header.Get("Content-Type")); mediaType != "" && mediaType != "application/json" {
		writeIngestionJSON(w, http.StatusBadRequest, map[string]string{"message": "Only the OTLP JSON encoding (application/json) is supported; OTLP protobuf ingestion is not implemented on this platform"})
		return
	}
	h.serveIngestion(w, r, false, parseOTLPJSON)
}

// parseOTLPJSON walks resourceLogs[].scopeLogs[].logRecords[]: each log
// record is one event whose message renders the body value (a bare
// stringValue stays the string; every other AnyValue form renders its
// value) and whose timestamp is timeUnixNano — observedTimeUnixNano when
// the record omits it, the server's clock when both are absent. A record
// without a body is skipped (the ingestion seam requires a non-empty
// message). Attributes, resource attributes and severity fields have no
// message representation in the endpoint's documented contract and stay
// off the stored message (recorded adjudication).
func parseOTLPJSON(body []byte) ([]ingestedRecord, *ingestionCounts, *ingestionError) {
	counts := &ingestionCounts{}
	if len(bytes.TrimSpace(body)) == 0 {
		return nil, counts, nil
	}
	var request struct {
		ResourceLogs []struct {
			ScopeLogs []struct {
				LogRecords []struct {
					TimeUnixNano         string                 `json:"timeUnixNano"`
					ObservedTimeUnixNano string                 `json:"observedTimeUnixNano"`
					Body                 map[string]interface{} `json:"body"`
				} `json:"logRecords"`
			} `json:"scopeLogs"`
		} `json:"resourceLogs"`
	}
	if err := json.Unmarshal(body, &request); err != nil {
		return nil, nil, &ingestionError{code: http.StatusBadRequest, message: "Request body is not valid OTLP JSON"}
	}
	var records []ingestedRecord
	now := time.Now().UnixNano()
	for _, rl := range request.ResourceLogs {
		for _, sl := range rl.ScopeLogs {
			for _, rec := range sl.LogRecords {
				message, ok := otlpBodyMessage(rec.Body)
				if !ok {
					counts.skipped++
					continue
				}
				timestamp := now
				if nanos, err := strconv.ParseInt(rec.TimeUnixNano, 10, 64); err == nil && nanos > 0 {
					timestamp = nanos
				} else if nanos, err := strconv.ParseInt(rec.ObservedTimeUnixNano, 10, 64); err == nil && nanos > 0 {
					timestamp = nanos
				}
				records = append(records, ingestedRecord{timestamp: timestamp / 1e6, message: message})
			}
		}
	}
	return records, counts, nil
}

// otlpBodyMessage renders an OTLP AnyValue body to the event message. A
// bare stringValue is the message verbatim; the scalar forms render their
// value; arrayValue and kvlistValue render as compact JSON. An absent or
// empty body carries no message.
func otlpBodyMessage(body map[string]interface{}) (string, bool) {
	if len(body) == 0 {
		return "", false
	}
	if v, ok := body["stringValue"].(string); ok {
		return v, v != ""
	}
	if v, ok := body["boolValue"].(bool); ok {
		return strconv.FormatBool(v), true
	}
	if v, ok := body["intValue"].(string); ok {
		return v, v != ""
	}
	// OTLP/JSON carries NaN and the infinities as doubleValue's string
	// form (JSON numbers cannot express them): the string is the value's
	// own textual rendering, and a record carrying it must not vanish
	// from ingestion.
	if v, ok := body["doubleValue"].(string); ok {
		return v, v != ""
	}
	if v, ok := body["doubleValue"].(float64); ok {
		return strconv.FormatFloat(v, 'g', -1, 64), true
	}
	for _, key := range []string{"arrayValue", "kvlistValue"} {
		if v, ok := body[key]; ok {
			rendered, err := json.Marshal(v)
			if err != nil || len(rendered) == 0 {
				return "", false
			}
			return string(rendered), true
		}
	}
	// bytesValue is the base64 rendering of the value's bytes (OTLP/JSON):
	// the message renders the decoded bytes, and a value that does not
	// decode keeps its raw text rather than dropping the record.
	if v, ok := body["bytesValue"].(string); ok {
		if decoded, err := base64.StdEncoding.DecodeString(v); err == nil {
			return string(decoded), len(decoded) > 0
		}
		return v, v != ""
	}
	return "", false
}

// --- shared record builders ---

// objectRecord maps a JSON object to one event: the named timestamp
// member when it is numeric (epoch milliseconds), the server's clock
// otherwise; the message is the object's compact JSON serialisation.
func objectRecord(obj map[string]interface{}, timestampMember string) ingestedRecord {
	timestamp := time.Now().UnixMilli()
	if v, ok := obj[timestampMember].(float64); ok {
		timestamp = int64(v)
	}
	message, err := json.Marshal(obj)
	if err != nil {
		return ingestedRecord{timestamp: timestamp, message: ""}
	}
	return ingestedRecord{timestamp: timestamp, message: string(message)}
}

// primitiveRecord maps a non-object JSON value to one event on the
// server's clock: a bare string is the message itself, every other
// primitive its JSON literal.
func primitiveRecord(value interface{}) ingestedRecord {
	if s, ok := value.(string); ok {
		return ingestedRecord{timestamp: time.Now().UnixMilli(), message: s}
	}
	rendered, err := json.Marshal(value)
	if err != nil {
		return ingestedRecord{timestamp: time.Now().UnixMilli(), message: ""}
	}
	return ingestedRecord{timestamp: time.Now().UnixMilli(), message: string(rendered)}
}
