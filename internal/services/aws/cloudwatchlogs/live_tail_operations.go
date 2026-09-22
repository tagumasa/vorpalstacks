package cloudwatchlogs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	corelogs "vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/eventstream"
	"vorpalstacks/pkg/filterpattern"
)

// StartLiveTail: the Live Tail streaming session. The response is an
// awsJson1.1 event stream (the SubscribeToShard transport): one
// initial-response message, one sessionStart event, then one
// sessionUpdate event every second carrying the events ingested in the
// past second, until the three-hour session timeout, the client closes
// the stream, or the service shuts down. Updates carry at most 500
// events and flag sampling when more matched; a slow client gets the
// documented bounded buffer — 10 updates, 5000 events — after which the
// oldest frames drop.

// LiveTailEventStream is the streamed response of StartLiveTail: the
// event-stream reader the dispatcher serves with the eventstream content
// type.
type LiveTailEventStream struct {
	reader io.Reader
}

// GetStream returns the event-stream reader.
func (e *LiveTailEventStream) GetStream() io.Reader {
	return e.reader
}

// GetStreamHeaders returns the HTTP headers of the event stream.
func (e *LiveTailEventStream) GetStreamHeaders() http.Header {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/vnd.amazon.eventstream")
	return headers
}

// StartLiveTailInput is the parsed request of StartLiveTail.
type StartLiveTailInput struct {
	LogGroupIdentifiers   []string
	LogStreamNames        []string
	LogStreamNamePrefixes []string
	LogEventFilterPattern string
	Region                string
}

// liveTailFrame is one queued stream message: a per-second update or a
// terminal exception.
type liveTailFrame struct {
	kind    string // "update", "timeout", "failure"
	events  []liveTailEvent
	sampled bool
}

// liveTailEvent is one log event of a sessionUpdate.
type liveTailEvent struct {
	logGroupIdentifier string
	logStreamName      string
	message            string
	timestamp          int64
	ingestionTime      int64
}

// liveTailCursor is one stream's read position inside a session: the
// event-timestamp scan floor and the ingestion stamp below which events
// are already served.
type liveTailCursor struct {
	ts        int64
	ingestion int64
}

// liveTailFetchEvents is the collect loop's fetch seam: tests replace
// it to observe the scan bounds each tick drives.
var liveTailFetchEvents = fetchEventsUpTo

// startLiveTailCore validates the session request and opens the session
// stream. Validation failures answer through the operation's declared
// error identities before any frame is written; a passing request
// reserves one of the account's concurrent-session slots, released when
// the session ends.
func (s *LogsService) startLiveTailCore(ctx context.Context, input StartLiveTailInput) (*LiveTailEventStream, error) {
	if err := validateLiveTailInput(input); err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, err
	}
	groupNames, groupARNs, err := resolveLiveTailGroups(store, input)
	if err != nil {
		return nil, err
	}

	if !s.reserveLiveTailSession() {
		return nil, NewLogsError("LimitExceededException",
			fmt.Sprintf("The account has reached its limit of %d concurrent Live Tail sessions", logsstore.LiveTailConcurrentSessionsMax), 400)
	}

	sessionId, err := newLiveTailID()
	if err != nil {
		s.releaseLiveTailSession()
		return nil, err
	}
	requestId, err := newLiveTailID()
	if err != nil {
		s.releaseLiveTailSession()
		return nil, err
	}

	pr, pw := io.Pipe()
	frames := make(chan *liveTailFrame, logsstore.LiveTailBufferUpdatesMax)
	writerDone := make(chan struct{})

	writer := eventstream.NewEncoder(pw)
	go func() {
		defer close(writerDone)
		if !s.writeLiveTailOpening(writer, requestId, sessionId, groupARNs, input) {
			return
		}
		for frame := range frames {
			if !s.writeLiveTailFrame(writer, frame) {
				return
			}
			if frame.kind != "update" {
				return
			}
		}
	}()

	session := &liveTailSession{
		store:            store,
		groupNames:       groupNames,
		groupARNs:        groupARNs,
		input:            input,
		startedIngestion: time.Now().UnixMilli(),
		matcher:          filterpattern.NewMatcher(),
	}
	s.spawnTask(func(taskCtx context.Context) {
		defer s.releaseLiveTailSession()
		// Closing the pipe from the driver side too: a writer blocked on
		// a client that stopped reading must unblock when the session
		// ends, not wait for the peer.
		defer pw.Close()
		defer close(frames)
		session.run(taskCtx, frames, writerDone)
	})

	return &LiveTailEventStream{reader: pr}, nil
}

// validateLiveTailInput enforces the session request's member rules: the
// required identifier array with its ten-group ceiling, the
// per-identifier rules (the trailing-asterisk and alphabet and
// ARN-addressing forms), the stream-filter pairing rules with the
// stream-name alphabet and filter ceilings, and the filter-pattern
// length bound.
func validateLiveTailInput(input StartLiveTailInput) error {
	if len(input.LogGroupIdentifiers) == 0 {
		return errRequiredMember("logGroupIdentifiers")
	}
	if len(input.LogGroupIdentifiers) > logsstore.LiveTailSessionLogGroupsMax {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("A Live Tail session can include as many as %d log groups; the request carries %d",
				logsstore.LiveTailSessionLogGroupsMax, len(input.LogGroupIdentifiers)), 400)
	}
	for _, id := range input.LogGroupIdentifiers {
		// The trailing-asterisk rule reads before the alphabet rule: the
		// member's own documented case (an ARN ending in *) carries its
		// documented message, and the alphabet below excludes '*' — in
		// the other order this rule would be unreachable wording.
		if strings.HasSuffix(id, "*") {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Log group identifier %s: an ARN can't end with an asterisk (*)", id), 400)
		}
		if id == "" || len(id) > logsstore.MaxLogGroupIdentifierLength || !logGroupIdentifierElementPattern.MatchString(id) {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Each log group identifier must be 1 to %d characters of the allowed alphabet",
					logsstore.MaxLogGroupIdentifierLength), 400)
		}
		// "Specify each log group by its ARN. If you specify an ARN, the
		// ARN can't end with an asterisk (*)" (StartLiveTail
		// logGroupIdentifiers member documentation): a bare name is not
		// the documented addressing form.
		if svcarn.ExtractLogGroupNameFromARN(id) == "" {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Log group identifier %s: specify each log group by its ARN", id), 400)
		}
	}
	// InputLogStreamNames carries @length {1,100} ("Array Members:
	// Minimum number of 1 item."): the optional stream lists cannot arrive
	// as the empty list — a present-but-empty parse (marked by the non-nil
	// empty slice) rejects instead of silently meaning "absent".
	if (input.LogStreamNames != nil && len(input.LogStreamNames) == 0) ||
		(input.LogStreamNamePrefixes != nil && len(input.LogStreamNamePrefixes) == 0) {
		return NewLogsError("InvalidParameterException",
			"logStreamNames and logStreamNamePrefixes must contain at least 1 entry", 400)
	}
	if len(input.LogStreamNames) > 0 && len(input.LogStreamNamePrefixes) > 0 {
		return NewLogsError("InvalidParameterException",
			"logStreamNames and logStreamNamePrefixes are mutually exclusive; only one of these parameters can be passed", 400)
	}
	if (len(input.LogStreamNames) > 0 || len(input.LogStreamNamePrefixes) > 0) && len(input.LogGroupIdentifiers) != 1 {
		return NewLogsError("InvalidParameterException",
			"logStreamNames and logStreamNamePrefixes can be specified only when the session includes exactly one log group", 400)
	}
	for _, name := range append(append([]string{}, input.LogStreamNames...), input.LogStreamNamePrefixes...) {
		if name == "" || len(name) > logsstore.MaxLogStreamNameLength {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Each log stream name or prefix must be 1 to %d characters", logsstore.MaxLogStreamNameLength), 400)
		}
		// The stream members carry the LogStreamName pattern "[^:*]":
		// neither separator the wildcard forms on reaches the session.
		if strings.ContainsAny(name, ":*") {
			return NewLogsError("InvalidParameterException",
				"Each log stream name or prefix must not contain ':' or '*'", 400)
		}
	}
	if len(input.LogStreamNames) > logsstore.LiveTailStreamFiltersMax || len(input.LogStreamNamePrefixes) > logsstore.LiveTailStreamFiltersMax {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("A Live Tail session can filter on as many as %d log streams", logsstore.LiveTailStreamFiltersMax), 400)
	}
	if len(input.LogEventFilterPattern) > logsstore.LiveTailFilterPatternMax {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("logEventFilterPattern can be as many as %d characters", logsstore.LiveTailFilterPatternMax), 400)
	}
	return nil
}

// resolveLiveTailGroups resolves the request's identifiers to their group
// records — the unknown group answers the operation's declared
// ResourceNotFoundException and the non-Standard class its
// InvalidOperationException — returning the resolved names and the
// ARN-form identifiers the session echoes.
func resolveLiveTailGroups(store *logsstore.Store, input StartLiveTailInput) ([]string, []string, error) {
	groupNames := make([]string, 0, len(input.LogGroupIdentifiers))
	groupARNs := make([]string, 0, len(input.LogGroupIdentifiers))
	for _, id := range input.LogGroupIdentifiers {
		name := resolveLogGroupIdentifier(id)
		group, err := store.GetLogGroup(name)
		if err != nil {
			return nil, nil, NewLogsError("ResourceNotFoundException",
				fmt.Sprintf("The specified log group does not exist: %s", name), 400)
		}
		// "Live Tail is supported only for log groups in the Standard
		// log class" — the developer guide's scope rule surfaces as the
		// operation's declared InvalidOperationException.
		if class := group.LogGroupClass; class != "" && class != "STANDARD" {
			return nil, nil, NewLogsError("InvalidOperationException",
				fmt.Sprintf("Live Tail is not supported for log group %s: Live Tail is supported only for log groups in the Standard log class", name), 400)
		}
		groupNames = append(groupNames, name)
		groupARNs = append(groupARNs, store.ARNBuilder().CloudWatch().LogGroup(name))
	}
	return groupNames, groupARNs, nil
}

// liveTailSession carries one open session's read state.
type liveTailSession struct {
	store            *logsstore.Store
	groupNames       []string
	groupARNs        []string
	input            StartLiveTailInput
	startedIngestion int64
	matcher          *filterpattern.Matcher
}

// run drives the per-second collection loop until the session ends: the
// timeout, the service shutdown, or the writer's exit (a client that
// closed the stream or stopped reading — the bounded frame queue drops
// the oldest updates past its capacity, and a writer failure ends the
// session at once).
func (session *liveTailSession) run(ctx context.Context, frames chan *liveTailFrame, writerDone <-chan struct{}) {
	defer func() {
		if r := recover(); r != nil {
			resilience.LogPanic("cloudwatchlogs.StartLiveTail.session", r)
			select {
			case frames <- &liveTailFrame{kind: "failure"}:
			default:
			}
		}
	}()

	cursors := make(map[string]liveTailCursor)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeout := time.NewTimer(logsstore.LiveTailSessionTimeout)
	defer timeout.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-writerDone:
			return
		case <-timeout.C:
			// The session ends at its timeout even when the client
			// stopped reading: the frame is enqueued under the same
			// slow-client policy as the update arm (drop the oldest
			// buffered frame, never block on a stalled writer), so run
			// returns and the session's resources are released.
		enqueueTimeout:
			for {
				select {
				case frames <- &liveTailFrame{kind: "timeout"}:
					break enqueueTimeout
				default:
				}
				select {
				case <-frames:
				case <-writerDone:
					break enqueueTimeout
				case <-ctx.Done():
					break enqueueTimeout
				}
			}
			return
		case <-ticker.C:
			events := session.collect(cursors)
			frame := &liveTailFrame{kind: "update", events: events}
			if len(events) > logsstore.LiveTailUpdateEventsMax {
				frame.events = events[:logsstore.LiveTailUpdateEventsMax]
				frame.sampled = true
			}
			// The documented slow-client policy: buffer up to the frame
			// bound, then drop the oldest. A client that stays connected
			// but stops reading parks the writer on the pipe with the
			// frame queue full — no frame drain and no writer exit ever
			// come — so the cancellation arm ends the session and frees
			// its slot instead of holding it for the stalled peer.
		enqueue:
			for {
				select {
				case frames <- frame:
					break enqueue
				default:
				}
				select {
				case <-frames:
				case <-writerDone:
					return
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

// collect reads one second's worth of newly ingested events across the
// session's groups. Each stream's scan advances by event timestamp with
// an ingestion floor at the session's start; the deterministic result
// order is the read engine's key (timestamp, ingestion time, stream
// name) with the group ahead of it.
func (session *liveTailSession) collect(cursors map[string]liveTailCursor) []liveTailEvent {
	var out []liveTailEvent
	now := time.Now().UnixMilli()
	streamNames := make(map[string]bool, len(session.input.LogStreamNames))
	for _, name := range session.input.LogStreamNames {
		streamNames[name] = true
	}
	for i, group := range session.groupNames {
		streams, err := fetchAllLogStreams(session.store, group, "")
		if err != nil {
			corelogs.Error("Live Tail failed to list log streams",
				corelogs.String("logGroup", group), corelogs.Err(err))
			continue
		}
		for _, ls := range streams {
			if len(streamNames) > 0 && !streamNames[ls.Name] {
				continue
			}
			if len(session.input.LogStreamNamePrefixes) > 0 && !hasAnyPrefix(ls.Name, session.input.LogStreamNamePrefixes) {
				continue
			}
			key := group + "\x00" + ls.Name
			cursor := cursors[key]
			events, err := liveTailFetchEvents(session.store, group, ls.Name, cursor.ts, now, 0)
			if err != nil {
				corelogs.Error("Live Tail failed to read log events",
					corelogs.String("logGroup", group), corelogs.String("logStream", ls.Name), corelogs.Err(err))
				continue
			}
			// The dedupe baseline is the watermark as of this fetch, not
			// the running maximum: events sharing an ingestion
			// millisecond within one batch must all pass (comparing
			// against the in-loop maximum would drop every one after the
			// first), while the next fetch's baseline is the advanced
			// watermark, so the re-read window emits nothing twice.
			batchFloor := cursor.ingestion
			for _, evt := range events {
				// The consumption cursor advances past every scanned
				// event — pre-session history and filter non-matches
				// included — so the next fetch starts where this one
				// ended: the filter selects what is emitted, not what is
				// consumed, and a non-matching tail re-fetched and
				// re-matched on every tick would grow the per-second
				// scan without bound. The fetch floor never rises past
				// the collect window: a future-stamped event
				// (PutLogEvents accepts +2h) must not push the floor
				// beyond every later normal event and mute the stream —
				// dedupe is the ingestion watermark's job.
				if evt.Timestamp > cursor.ts && evt.Timestamp <= now {
					cursor.ts = evt.Timestamp
				}
				if evt.IngestionTime > cursor.ingestion {
					cursor.ingestion = evt.IngestionTime
				}
				if evt.IngestionTime < session.startedIngestion || evt.IngestionTime <= batchFloor {
					continue
				}
				if session.input.LogEventFilterPattern != "" &&
					!session.matcher.Matches(session.input.LogEventFilterPattern, evt.Message) {
					continue
				}
				out = append(out, liveTailEvent{
					logGroupIdentifier: session.groupARNs[i],
					logStreamName:      ls.Name,
					message:            evt.Message,
					timestamp:          evt.Timestamp,
					ingestionTime:      evt.IngestionTime,
				})
			}
			cursors[key] = cursor
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].timestamp != out[j].timestamp {
			return out[i].timestamp < out[j].timestamp
		}
		if out[i].ingestionTime != out[j].ingestionTime {
			return out[i].ingestionTime < out[j].ingestionTime
		}
		if out[i].logGroupIdentifier != out[j].logGroupIdentifier {
			return out[i].logGroupIdentifier < out[j].logGroupIdentifier
		}
		return out[i].logStreamName < out[j].logStreamName
	})
	return out
}

// writeLiveTailOpening writes the initial-response message and the
// sessionStart event. The start event echoes the session's selectors —
// the identifiers in the ARN form the request resolved them to.
func (s *LogsService) writeLiveTailOpening(encoder *eventstream.Encoder, requestId, sessionId string, groupARNs []string, input StartLiveTailInput) bool {
	if err := encoder.WriteInitialResponse([]byte("{}")); err != nil {
		corelogs.Error("Live Tail failed to write initial response", corelogs.Err(err))
		return false
	}
	start := map[string]interface{}{
		"requestId":           requestId,
		"sessionId":           sessionId,
		"logGroupIdentifiers": groupARNs,
	}
	if len(input.LogStreamNames) > 0 {
		start["logStreamNames"] = input.LogStreamNames
	}
	if len(input.LogStreamNamePrefixes) > 0 {
		start["logStreamNamePrefixes"] = input.LogStreamNamePrefixes
	}
	if input.LogEventFilterPattern != "" {
		start["logEventFilterPattern"] = input.LogEventFilterPattern
	}
	return s.writeLiveTailEvent(encoder, "sessionStart", start)
}

// writeLiveTailFrame writes one queued message.
func (s *LogsService) writeLiveTailFrame(encoder *eventstream.Encoder, frame *liveTailFrame) bool {
	switch frame.kind {
	case "timeout":
		return s.writeLiveTailException(encoder, "SessionTimeoutException",
			"The Live Tail session timed out: sessions time out after three hours")
	case "failure":
		return s.writeLiveTailException(encoder, "SessionStreamingException",
			"An unknown error occurred on the server side")
	default:
		results := make([]map[string]interface{}, 0, len(frame.events))
		for _, evt := range frame.events {
			results = append(results, map[string]interface{}{
				"logStreamName":      evt.logStreamName,
				"logGroupIdentifier": evt.logGroupIdentifier,
				"message":            evt.message,
				"timestamp":          evt.timestamp,
				"ingestionTime":      evt.ingestionTime,
			})
		}
		return s.writeLiveTailEvent(encoder, "sessionUpdate", map[string]interface{}{
			"sessionMetadata": map[string]interface{}{"sampled": frame.sampled},
			"sessionResults":  results,
		})
	}
}

func (s *LogsService) writeLiveTailEvent(encoder *eventstream.Encoder, eventType string, payload map[string]interface{}) bool {
	body, err := json.Marshal(payload)
	if err != nil {
		corelogs.Error("Live Tail failed to marshal event", corelogs.String("eventType", eventType), corelogs.Err(err))
		return false
	}
	if err := encoder.WriteEvent(eventType, "application/json", body); err != nil {
		corelogs.Error("Live Tail failed to write event", corelogs.String("eventType", eventType), corelogs.Err(err))
		return false
	}
	return true
}

func (s *LogsService) writeLiveTailException(encoder *eventstream.Encoder, exceptionType, message string) bool {
	if err := encoder.WriteErrorEvent(exceptionType, message); err != nil {
		corelogs.Error("Live Tail failed to write streaming exception",
			corelogs.String("exceptionType", exceptionType), corelogs.Err(err))
		return false
	}
	return true
}

// reserveLiveTailSession takes one concurrent-session slot, reporting
// whether the documented quota still had room.
func (s *LogsService) reserveLiveTailSession() bool {
	s.liveTailMu.Lock()
	defer s.liveTailMu.Unlock()
	if s.liveTailSessions >= logsstore.LiveTailConcurrentSessionsMax {
		return false
	}
	s.liveTailSessions++
	return true
}

func (s *LogsService) releaseLiveTailSession() {
	s.liveTailMu.Lock()
	defer s.liveTailMu.Unlock()
	if s.liveTailSessions > 0 {
		s.liveTailSessions--
	}
}

func newLiveTailID() (string, error) {
	return newRandomHexID(16)
}

// --- HTTP handler ---

func (s *LogsService) StartLiveTail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// An explicitly empty array is a present member, not an omitted one:
	// the non-nil empty slice carries that distinction to the Core, where
	// the stream lists' @length(min 1) trait rejects it.
	streamNames := request.GetStringList(req.Parameters, "LogStreamNames")
	if streamNames == nil && request.HasListParam(req.Parameters, "LogStreamNames") {
		streamNames = []string{}
	}
	streamNamePrefixes := request.GetStringList(req.Parameters, "LogStreamNamePrefixes")
	if streamNamePrefixes == nil && request.HasListParam(req.Parameters, "LogStreamNamePrefixes") {
		streamNamePrefixes = []string{}
	}
	return s.startLiveTailCore(ctx, StartLiveTailInput{
		LogGroupIdentifiers:   request.GetStringList(req.Parameters, "LogGroupIdentifiers"),
		LogStreamNames:        streamNames,
		LogStreamNamePrefixes: streamNamePrefixes,
		LogEventFilterPattern: request.GetParamLowerFirst(req.Parameters, "LogEventFilterPattern"),
		Region:                reqCtx.GetRegion(),
	})
}
