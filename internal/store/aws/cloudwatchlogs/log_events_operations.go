package cloudwatchlogs

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/pkg/filterpattern"
)

// PutLogEvents puts log events to a CloudWatch Logs log stream.
func (s *Store) PutLogEvents(logGroupName, logStreamName string, events []LogEntry) (string, error) {
	// The existence checks run under the group's lock, the same lock the
	// commit takes: a DeleteLogStream/DeleteLogGroup of this group that
	// acquires the lock first makes this call fail here instead of
	// resurrecting the deleted records at the commit below.
	gLock := s.groupLock(logGroupName)
	gLock.Lock()
	defer gLock.Unlock()

	lg, err := s.GetLogGroup(logGroupName)
	if err != nil {
		if !errors.Is(err, ErrLogGroupNotFound) {
			return "", err
		}
		return "", ErrLogGroupNotFound
	}

	ls, err := s.GetLogStream(logGroupName, logStreamName)
	if err != nil {
		if !errors.Is(err, ErrLogStreamNotFound) {
			return "", err
		}
		return "", ErrLogStreamNotFound
	}

	// Fallback (non-transactional) commits tear: the sequential writes can
	// fail midway with the mutated records already on disk, so the
	// pre-mutation snapshots are kept for the restore below (the shallow
	// copies are safe — the mutations touch scalars only, and the Tags
	// map is only ever read here).
	lgBefore, lsBefore := *lg, *ls

	ingestionTs := time.Now().UnixMilli()
	for i := range events {
		if events[i].IngestionTime == 0 {
			events[i].IngestionTime = ingestionTs
		}
	}

	var minTs, maxTs int64
	var bytesAdded int64
	for i, e := range events {
		if i == 0 || e.Timestamp < minTs {
			minTs = e.Timestamp
		}
		if i == 0 || e.Timestamp > maxTs {
			maxTs = e.Timestamp
		}
		bytesAdded += int64(len(e.Message))
	}
	ls.UpdateEventTimestamps(minTs, maxTs)
	ls.LastIngestionTs = time.Now().UnixMilli()

	newToken, err := incrementToken(ls.UploadSequenceToken)
	if err != nil {
		return "", err
	}

	// Every acknowledged PutLogEvents call persists its own chunk. AWS
	// treats a successful PutLogEvents as durable, so buffering events in
	// memory across calls (and losing the buffer on restart) would break
	// that contract. Collect the chunk index entry that must be committed
	// atomically with the LogGroup/LogStream metadata update; if the
	// metadata transaction fails, the chunk file is removed to prevent
	// orphans. The sequence token is advanced before the chunk is written
	// so a token failure cannot orphan an already-persisted chunk.
	var pendingChunks []pendingChunkIndex
	if len(events) > 0 {
		pi, err := s.prepareChunkFlush(logGroupName, logStreamName, events)
		if err != nil {
			return "", err
		}
		if pi.meta != nil {
			pendingChunks = append(pendingChunks, pi)
		}
	}

	ls.UploadSequenceToken = newToken
	lg.StoredBytes += bytesAdded

	// Atomically commit the LogGroup, LogStream, and any pending chunk
	// indexes in a single transaction. Without this, a failure between
	// the chunk index write and the metadata update leaves orphaned
	// chunks with inconsistent StoredBytes and timestamp metadata.
	if err := s.updateLogGroupStreamAndChunks(logGroupName, lg, logStreamName, ls, pendingChunks); err != nil {
		// The transaction failed; remove the orphaned chunk files so they
		// do not leak storage with no index pointing to them.
		for _, pc := range pendingChunks {
			if rmErr := os.Remove(pc.chunkPath); rmErr != nil && !os.IsNotExist(rmErr) {
				logs.Error("Failed to remove orphaned chunk file after transaction failure",
					logs.String("path", pc.chunkPath), logs.Err(rmErr))
			}
		}
		// On the non-transactional path the writes that landed before the
		// failure are compensated: the chunk index records by
		// updateLogGroupStreamAndChunks itself, and the mutated
		// LogGroup/LogStream records by rewriting the pre-mutation
		// snapshots taken above the mutations.
		if s.ts == nil {
			if restoreErr := s.putLogGroupLocked(&lgBefore); restoreErr != nil {
				logs.Error("Failed to restore LogGroup after fallback commit failure",
					logs.String("logGroupName", logGroupName), logs.Err(restoreErr))
			}
			if restoreErr := s.PutProto(s.logStreamKey(logGroupName, logStreamName), LogStreamToProto(&lsBefore)); restoreErr != nil {
				logs.Error("Failed to restore LogStream after fallback commit failure",
					logs.String("logStreamName", logStreamName), logs.Err(restoreErr))
			}
		}
		return "", err
	}

	return ls.UploadSequenceToken, nil
}

func incrementToken(token string) (string, error) {
	if token == "" {
		return "1", nil
	}
	val, err := strconv.Atoi(token)
	if err != nil {
		return "", fmt.Errorf("corrupted upload sequence token %q: %w", token, err)
	}
	return strconv.Itoa(val + 1), nil
}

func matchFilterPattern(message, pattern string) bool {
	matcher := filterpattern.NewMatcher()
	return matcher.Matches(pattern, message)
}

// --- The shared event read engine ---
//
// GetLogEvents and FilterLogEvents are two scopes over one engine: a
// single-stream read and a group-wide (or explicit stream set) read with
// an optional filter pattern. Both order their results by the same
// documented key and page with the same scoped-token vocabulary, so the
// direction/cursor/slicing logic exists once.

// eventReadQuery describes one event read through the shared engine.
type eventReadQuery struct {
	Group string
	// Stream scopes the read to one log stream (the GetLogEvents shape).
	Stream string
	// Streams restricts a group-wide read to an explicit set (the
	// FilterLogEvents logStreamNames shape); empty means the whole group.
	Streams   []string
	StartTime int64
	EndTime   int64
	// EndTimeExclusive selects the endTime boundary semantics: GetLogEvents
	// documents endTime exclusively ("Events with a timestamp equal to or
	// later than this time are not included"), while FilterLogEvents
	// documents it inclusively ("events with a timestamp later than this
	// time are not returned"). The default keeps the inclusive comparison.
	EndTimeExclusive bool
	Pattern          string
	// Limit caps one page; <= 0 serves the whole scope (the internal
	// fetch-all layer's contract).
	Limit int
	// MaxResponseBytes caps one page by serialized size (the documented
	// GetLogEvents default, "as many log events as can fit in a response
	// size of 1 MB"); zero disables the cap. The trim keeps at least one
	// event so pagination always makes progress.
	MaxResponseBytes int
	// StartFromHead sets the direction of the FIRST request only; a
	// presented token's direction wins on subsequent requests.
	StartFromHead bool
	NextToken     string
}

// indexedEvent pairs an output event with its cursor tiebreak digest and
// ordinal so the sort and the cursor comparison see one key.
type indexedEvent struct {
	event  *OutputLogEvent
	digest string
	// ordinal is the chunk identifier plus the event's position inside
	// the chunk: the stand-in for the documented per-put request ID that
	// keeps byte-identical duplicate events distinct.
	ordinal string
}

// eventOrdinal derives the cursor's per-event disambiguator from the
// immutable chunk identity: the chunk ID is fixed-width, and so is the
// position, so the composite string orders exactly as (chunk, position).
func eventOrdinal(chunkID string, position int) string {
	return chunkID + "/" + fmt.Sprintf("%06d", position)
}

func cursorOfEvent(ie indexedEvent) eventCursor {
	return eventCursor{
		Timestamp:     ie.event.Timestamp,
		IngestionTime: ie.event.IngestionTime,
		LogStream:     ie.event.LogStreamName,
		MessageDigest: ie.digest,
		Ordinal:       ie.ordinal,
	}
}

// cursorLess orders two cursors by the read engine's total key:
// (timestamp, ingestion time, log stream, message digest, ordinal).
func cursorLess(a, b eventCursor) bool {
	if a.Timestamp != b.Timestamp {
		return a.Timestamp < b.Timestamp
	}
	if a.IngestionTime != b.IngestionTime {
		return a.IngestionTime < b.IngestionTime
	}
	if a.LogStream != b.LogStream {
		return a.LogStream < b.LogStream
	}
	if a.MessageDigest != b.MessageDigest {
		return a.MessageDigest < b.MessageDigest
	}
	return a.Ordinal < b.Ordinal
}

// gatherEventScope reads every event in the query scope into the
// engine's deterministic total order. A chunk-listing failure
// propagates: the caller answers an error, never a shortened page.
// The gather runs under the group's read lock: chunk files are removed
// and their index records deleted under the write side of the same
// lock (stream deletes, group teardown, retention purges), so a page
// served here can never interleave a mid-flight delete and serve a
// torn listing whose files are already gone.
//
// A resumed walk narrows the gather to the unserved side: the cursor's
// timestamp is the sort's leading key, so a chunk whose bounds sit
// entirely on the served side (forward: MaxTs below the cursor;
// backward: MinTs above it) holds only events the page's binary search
// would exclude, and re-reading and re-sorting it on every page is
// pure repeated work. Chunks whose bounds merely touch the cursor stay
// gathered: their same-timestamp events order by the tie-break members
// (ingestion time, stream, digest, ordinal), which the chunk bounds do
// not see. The served side's has-more hint is unaffected — the paging
// envelope reads the hint of the walking direction only, and that side
// is never skipped.
func (s *Store) gatherEventScope(q eventReadQuery, direction byte, resume eventCursor, haveResume bool) ([]indexedEvent, error) {
	gLock := s.groupLock(q.Group)
	gLock.RLock()
	defer gLock.RUnlock()

	var chunks []*ChunkMeta
	switch {
	case q.Stream != "":
		listed, err := s.ListChunksForStream(q.Group, q.Stream)
		if err != nil {
			return nil, err
		}
		chunks = listed
	case len(q.Streams) > 0:
		names := make([]string, len(q.Streams))
		copy(names, q.Streams)
		sort.Strings(names)
		// The name list is a set on the wire, not an ordered sequence:
		// a repeated name must not list its stream's chunks twice (every
		// event would return once per repetition), so the sorted walk
		// skips a name equal to its predecessor.
		for i, name := range names {
			if i > 0 && names[i-1] == name {
				continue
			}
			listed, err := s.ListChunksForStream(q.Group, name)
			if err != nil {
				return nil, err
			}
			chunks = append(chunks, listed...)
		}
	default:
		listed, err := s.ListChunksForLogGroup(q.Group)
		if err != nil {
			return nil, err
		}
		chunks = listed
	}

	var gathered []indexedEvent
	for _, chunk := range chunks {
		if chunk == nil {
			continue
		}
		if q.EndTime > 0 && chunk.MinTs > q.EndTime {
			continue
		}
		if q.StartTime > 0 && chunk.MaxTs < q.StartTime {
			continue
		}
		if haveResume {
			if direction == PageForward && chunk.MaxTs < resume.Timestamp {
				continue
			}
			if direction == PageBackward && chunk.MinTs > resume.Timestamp {
				continue
			}
		}

		entries, err := s.readChunkFile(chunk.ChunkPath)
		if err != nil {
			// A chunk that fails to read is invisible to the caller;
			// log it so corruption is discoverable instead of silently
			// shrinking the result.
			logs.Warn("Failed to read log events chunk",
				logs.String("chunkPath", chunk.ChunkPath), logs.Err(err))
			continue
		}

		for idx, e := range entries {
			if q.StartTime > 0 && e.Timestamp < q.StartTime {
				continue
			}
			if q.EndTime > 0 && (e.Timestamp > q.EndTime || (q.EndTimeExclusive && e.Timestamp == q.EndTime)) {
				continue
			}
			if q.Pattern != "" && !matchFilterPattern(e.Message, q.Pattern) {
				continue
			}
			ingestionTime := e.IngestionTime
			if ingestionTime == 0 {
				ingestionTime = e.Timestamp
			}
			ordinal := eventOrdinal(chunk.ChunkID, idx)
			gathered = append(gathered, indexedEvent{
				event: &OutputLogEvent{
					Timestamp:     e.Timestamp,
					Message:       e.Message,
					IngestionTime: ingestionTime,
					LogStreamName: chunk.LogStream,
					Ordinal:       ordinal,
				},
				digest:  messageDigest(e.Message),
				ordinal: ordinal,
			})
		}
	}

	sort.Slice(gathered, func(i, j int) bool {
		return cursorLess(cursorOfEvent(gathered[i]), cursorOfEvent(gathered[j]))
	})
	return gathered, nil
}

// eventPage is one served page plus the boundary cursors the caller needs
// to mint continuation tokens in either direction.
type eventPage struct {
	Events []*OutputLogEvent
	// Direction is the direction this page was served in.
	Direction byte
	// HasMoreForward/HasMoreBackward report whether the scope holds
	// events beyond the served window in each direction.
	HasMoreForward  bool
	HasMoreBackward bool
	// ForwardCursor is the cursor a following forward page resumes after
	// (the max-key served event, or the page's entry boundary when the
	// page is empty); BackwardCursor is the cursor a following backward
	// page resumes before (the min-key served event, or the boundary).
	ForwardCursor  eventCursor
	BackwardCursor eventCursor
	// MintedMs is the walk's token mint: the presented token's mint when
	// one decoded, zero for a fresh walk. The response tokens carry it so
	// the walk's tokens share one expiry clock and the end-of-stream
	// re-offer equals the presented token.
	MintedMs int64
}

// eventResponseBytes estimates the serialized cost of one output log event
// against the response-size budget: its message bytes plus the fixed
// JSON-envelope estimate.
func eventResponseBytes(e *OutputLogEvent) int {
	return len(e.Message) + EventResponseEnvelopeBytes
}

// readEventPage gathers the query scope and serves one direction-aware
// page. A presented token must decode into the request's own scope; a
// token from any other request shape is rejected.
func (s *Store) readEventPage(q eventReadQuery) (*eventPage, error) {
	scope := eventPageToken{
		Version:     eventPageTokenVersion,
		Group:       q.Group,
		Stream:      q.Stream,
		StreamSetID: streamSetID(q.Streams),
		StartTime:   q.StartTime,
		EndTime:     q.EndTime,
		Pattern:     q.Pattern,
	}

	direction := PageForward
	var resume eventCursor
	haveResume := false
	if q.NextToken != "" {
		tok, err := decodeEventPageToken(q.NextToken)
		if err != nil {
			return nil, err
		}
		if tok.Group != scope.Group || tok.Stream != scope.Stream || tok.StreamSetID != scope.StreamSetID ||
			tok.StartTime != scope.StartTime || tok.EndTime != scope.EndTime || tok.Pattern != scope.Pattern {
			return nil, ErrInvalidPaginationToken
		}
		direction = tok.Direction
		resume = tok.Cursor
		haveResume = true
		// The walk's mint rides the presented token; the response tokens
		// carry it onward so one walk's tokens share one expiry clock.
		scope.MintedMs = tok.MintedMs
	} else if !q.StartFromHead {
		direction = PageBackward
	}

	events, err := s.gatherEventScope(q, direction, resume, haveResume)
	if err != nil {
		return nil, err
	}

	lo, hi := 0, len(events)
	if direction == PageForward {
		if haveResume {
			lo = sort.Search(len(events), func(i int) bool {
				return cursorLess(resume, cursorOfEvent(events[i]))
			})
		}
		hi = lo + q.Limit
		if q.Limit <= 0 || hi > len(events) {
			hi = len(events)
		}
	} else {
		if haveResume {
			hi = sort.Search(len(events), func(i int) bool {
				return !cursorLess(cursorOfEvent(events[i]), resume)
			})
		}
		lo = hi - q.Limit
		if q.Limit <= 0 || lo < 0 {
			lo = 0
		}
	}

	// The response-size budget trims the count-bounded window from its
	// leading edge (the direction the page serves), keeping one event so a
	// page of oversized events still advances the cursor.
	if q.MaxResponseBytes > 0 && lo < hi {
		budget := q.MaxResponseBytes
		if direction == PageForward {
			end := lo + 1
			budget -= eventResponseBytes(events[lo].event)
			for end < hi {
				cost := eventResponseBytes(events[end].event)
				if budget-cost < 0 {
					break
				}
				budget -= cost
				end++
			}
			hi = end
		} else {
			start := hi - 1
			budget -= eventResponseBytes(events[hi-1].event)
			for start > lo {
				cost := eventResponseBytes(events[start-1].event)
				if budget-cost < 0 {
					break
				}
				budget -= cost
				start--
			}
			lo = start
		}
	}

	page := &eventPage{Direction: direction}
	if lo < hi {
		if direction == PageForward {
			for i := lo; i < hi; i++ {
				page.Events = append(page.Events, events[i].event)
			}
		} else {
			// Backward pages serve their window newest-first.
			for i := hi - 1; i >= lo; i-- {
				page.Events = append(page.Events, events[i].event)
			}
		}
		page.ForwardCursor = cursorOfEvent(events[hi-1])
		page.BackwardCursor = cursorOfEvent(events[lo])
	} else if haveResume {
		// An empty page re-offers its entry cursor, which encodes to the
		// very token the caller presented (GetLogEvents' documented
		// end-of-stream rule).
		page.ForwardCursor = resume
		page.BackwardCursor = resume
	} else if direction == PageForward {
		page.ForwardCursor = cursorBeforeZero()
		page.BackwardCursor = cursorBeforeZero()
	} else {
		page.ForwardCursor = cursorAfterAll()
		page.BackwardCursor = cursorAfterAll()
	}
	page.HasMoreForward = hi < len(events)
	page.HasMoreBackward = lo > 0
	page.MintedMs = scope.MintedMs
	return page, nil
}

// GetLogEvents retrieves log events from one CloudWatch Logs log stream.
// Both returned tokens are always present (the model documents them as
// never null); at the end of a stream a token equals the one the caller
// presented, which is how clients detect that pagination is finished.
// The endTime boundary is exclusive per this operation's documentation
// ("Events with a timestamp equal to or later than this time are not
// included").
func (s *Store) GetLogEvents(logGroupName, logStreamName string, startTime, endTime int64, limit int, startFromHead bool, nextToken string) ([]*OutputLogEvent, string, string, error) {
	return s.getLogEventsPage(logGroupName, logStreamName, startTime, endTime, limit, startFromHead, nextToken, true)
}

// GetLogEventsEndInclusive reads one page of the same single-stream scope
// with an inclusive endTime — the boundary the internal consumers that
// define their windows through documented inclusive ends use ("The range
// is inclusive, so the specified end time is included in the query",
// StartQuery; "Events with a timestamp later than this time are not
// exported", CreateExportTask). The wire-facing GetLogEvents keeps its
// exclusive end; only the paging boundary differs between the two.
func (s *Store) GetLogEventsEndInclusive(logGroupName, logStreamName string, startTime, endTime int64, limit int, startFromHead bool, nextToken string) ([]*OutputLogEvent, string, string, error) {
	return s.getLogEventsPage(logGroupName, logStreamName, startTime, endTime, limit, startFromHead, nextToken, false)
}

// getLogEventsPage is the single GetLogEvents envelope: existence checks,
// one engine page and both direction tokens; endTimeExclusive selects the
// window's end boundary semantics.
func (s *Store) getLogEventsPage(logGroupName, logStreamName string, startTime, endTime int64, limit int, startFromHead bool, nextToken string, endTimeExclusive bool) ([]*OutputLogEvent, string, string, error) {
	if _, err := s.GetLogGroup(logGroupName); err != nil {
		if !errors.Is(err, ErrLogGroupNotFound) {
			return nil, "", "", err
		}
		return nil, "", "", ErrLogGroupNotFound
	}
	if _, err := s.GetLogStream(logGroupName, logStreamName); err != nil {
		if !errors.Is(err, ErrLogStreamNotFound) {
			return nil, "", "", err
		}
		return nil, "", "", ErrLogStreamNotFound
	}

	page, err := s.readEventPage(eventReadQuery{
		Group:            logGroupName,
		Stream:           logStreamName,
		StartTime:        startTime,
		EndTime:          endTime,
		EndTimeExclusive: endTimeExclusive,
		Limit:            limit,
		MaxResponseBytes: MaxReadResponseBytes,
		StartFromHead:    startFromHead,
		NextToken:        nextToken,
	})
	if err != nil {
		return nil, "", "", err
	}

	scope := eventPageToken{
		Version:   eventPageTokenVersion,
		Group:     logGroupName,
		Stream:    logStreamName,
		StartTime: startTime,
		EndTime:   endTime,
		MintedMs:  page.MintedMs,
	}
	scope.Direction = PageForward
	scope.Cursor = page.ForwardCursor
	nextForwardToken, err := encodeEventPageToken(scope)
	if err != nil {
		return nil, "", "", err
	}
	scope.Direction = PageBackward
	scope.Cursor = page.BackwardCursor
	nextBackwardToken, err := encodeEventPageToken(scope)
	if err != nil {
		return nil, "", "", err
	}

	return page.Events, nextForwardToken, nextBackwardToken, nil
}

// FilterLogEvents filters log events from a log group based on a filter
// pattern. The next token is present only while events remain in the
// page's direction (the model documents its absence as the end of
// pagination).
func (s *Store) FilterLogEvents(logGroupName string, logStreamNames []string, startTime, endTime int64, filterPattern string, limit int, startFromHead bool, nextToken string) ([]*OutputLogEvent, string, error) {
	if _, err := s.GetLogGroup(logGroupName); err != nil {
		if !errors.Is(err, ErrLogGroupNotFound) {
			return nil, "", err
		}
		return nil, "", ErrLogGroupNotFound
	}

	page, err := s.readEventPage(eventReadQuery{
		Group:            logGroupName,
		Streams:          logStreamNames,
		StartTime:        startTime,
		EndTime:          endTime,
		Pattern:          filterPattern,
		Limit:            limit,
		MaxResponseBytes: MaxReadResponseBytes,
		StartFromHead:    startFromHead,
		NextToken:        nextToken,
	})
	if err != nil {
		return nil, "", err
	}

	var newNextToken string
	if (page.Direction == PageForward && page.HasMoreForward) ||
		(page.Direction == PageBackward && page.HasMoreBackward) {
		scope := eventPageToken{
			Version:     eventPageTokenVersion,
			Group:       logGroupName,
			StreamSetID: streamSetID(logStreamNames),
			StartTime:   startTime,
			EndTime:     endTime,
			Pattern:     filterPattern,
			Direction:   page.Direction,
			MintedMs:    page.MintedMs,
		}
		if page.Direction == PageForward {
			scope.Cursor = page.ForwardCursor
		} else {
			scope.Cursor = page.BackwardCursor
		}
		newNextToken, err = encodeEventPageToken(scope)
		if err != nil {
			return nil, "", err
		}
	}

	return page.Events, newNextToken, nil
}
