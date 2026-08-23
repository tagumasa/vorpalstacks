package cloudwatchlogs

import (
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
	lg, err := s.GetLogGroup(logGroupName)
	if err != nil {
		return "", ErrLogGroupNotFound
	}

	ls, err := s.GetLogStream(logGroupName, logStreamName)
	if err != nil {
		return "", ErrLogStreamNotFound
	}

	s.chunkMutex.Lock()
	defer s.chunkMutex.Unlock()

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

func sortEventsByTimestamp(events []*OutputLogEvent) {
	// Stable: events with identical timestamps keep their relative
	// (ingestion) order, matching the order AWS returns them in.
	sort.SliceStable(events, func(i, j int) bool {
		return events[i].Timestamp < events[j].Timestamp
	})
}

// GetLogEvents retrieves log events from a CloudWatch Logs log stream.
func (s *Store) GetLogEvents(logGroupName, logStreamName string, startTime, endTime int64, limit int, startFromHead bool, nextToken string) ([]*OutputLogEvent, string, string, error) {
	if _, err := s.GetLogGroup(logGroupName); err != nil {
		return nil, "", "", ErrLogGroupNotFound
	}

	if _, err := s.GetLogStream(logGroupName, logStreamName); err != nil {
		return nil, "", "", ErrLogStreamNotFound
	}

	chunks := s.ListChunksForStream(logGroupName, logStreamName)
	var allEvents []*OutputLogEvent

	for _, chunk := range chunks {
		if chunk == nil {
			continue
		}
		if endTime > 0 && chunk.MinTs > endTime {
			continue
		}
		if startTime > 0 && chunk.MaxTs < startTime {
			continue
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

		for _, e := range entries {
			if startTime > 0 && e.Timestamp < startTime {
				continue
			}
			if endTime > 0 && e.Timestamp > endTime {
				continue
			}
			ingestionTime := e.IngestionTime
			if ingestionTime == 0 {
				ingestionTime = e.Timestamp
			}
			allEvents = append(allEvents, &OutputLogEvent{
				Timestamp:     e.Timestamp,
				Message:       e.Message,
				IngestionTime: ingestionTime,
			})
		}
	}

	// Always sort ascending internally; direction-aware slicing handles ordering.
	sortEventsByTimestamp(allEvents)

	// Parse the direction-aware token. For the first request (empty token),
	// derive direction from startFromHead and use the appropriate boundary.
	direction, offset, err := ParsePaginationToken(nextToken)
	if err != nil {
		return nil, "", "", err
	}
	if nextToken == "" {
		if startFromHead {
			direction = PaginationForward
			offset = 0
		} else {
			direction = PaginationBackward
			offset = len(allEvents)
		}
	}

	totalLen := len(allEvents)
	if offset > totalLen {
		offset = totalLen
	}

	var result []*OutputLogEvent
	var nextFwdOffset, nextBwdOffset int

	if direction == PaginationForward {
		endIdx := offset + limit
		if limit <= 0 || endIdx > totalLen {
			endIdx = totalLen
		}
		if offset < totalLen {
			result = allEvents[offset:endIdx]
		}
		nextFwdOffset = endIdx
		nextBwdOffset = offset
	} else {
		// Backward: read [max(0, offset-limit), offset) in descending order.
		startIdx := offset - limit
		if limit <= 0 || startIdx < 0 {
			startIdx = 0
		}
		result = make([]*OutputLogEvent, 0, offset-startIdx)
		for i := offset - 1; i >= startIdx; i-- {
			if i >= 0 && i < totalLen {
				result = append(result, allEvents[i])
			}
		}
		nextBwdOffset = startIdx
		nextFwdOffset = offset
	}

	// Return empty tokens at boundaries so clients can detect end-of-data.
	// Forward: no more events when next offset reaches total.
	// Backward: no more events when next offset reaches 0.
	var nextForwardToken, nextBackwardToken string
	if nextFwdOffset < totalLen {
		nextForwardToken = EncodePaginationToken(PaginationForward, nextFwdOffset)
	}
	if nextBwdOffset > 0 {
		nextBackwardToken = EncodePaginationToken(PaginationBackward, nextBwdOffset)
	}

	return result, nextForwardToken, nextBackwardToken, nil
}

// FilterLogEvents filters log events from a log group based on filter pattern.
func (s *Store) FilterLogEvents(logGroupName string, logStreamNames []string, startTime, endTime int64, filterPattern string, limit int, startFromHead bool, nextToken string) ([]*OutputLogEvent, map[string]bool, string, error) {
	if _, err := s.GetLogGroup(logGroupName); err != nil {
		return nil, nil, "", ErrLogGroupNotFound
	}

	var chunks []*ChunkMeta
	if len(logStreamNames) > 0 {
		for _, streamName := range logStreamNames {
			chunks = append(chunks, s.ListChunksForStream(logGroupName, streamName)...)
		}
	} else {
		chunks = s.ListChunksForLogGroup(logGroupName)
	}

	var allEvents []*OutputLogEvent
	searchedStreams := make(map[string]bool)

	for _, chunk := range chunks {
		if endTime > 0 && chunk.MinTs > endTime {
			continue
		}
		if startTime > 0 && chunk.MaxTs < startTime {
			continue
		}

		searchedStreams[chunk.LogStream] = true

		entries, err := s.readChunkFile(chunk.ChunkPath)
		if err != nil {
			// A chunk that fails to read is invisible to the caller;
			// log it so corruption is discoverable instead of silently
			// shrinking the result.
			logs.Warn("Failed to read log events chunk",
				logs.String("chunkPath", chunk.ChunkPath), logs.Err(err))
			continue
		}

		for _, e := range entries {
			if startTime > 0 && e.Timestamp < startTime {
				continue
			}
			if endTime > 0 && e.Timestamp > endTime {
				continue
			}
			if filterPattern != "" && !matchFilterPattern(e.Message, filterPattern) {
				continue
			}
			ingestionTime := e.IngestionTime
			if ingestionTime == 0 {
				ingestionTime = e.Timestamp
			}
			allEvents = append(allEvents, &OutputLogEvent{
				Timestamp:     e.Timestamp,
				Message:       e.Message,
				IngestionTime: ingestionTime,
				LogStreamName: chunk.LogStream,
			})
		}
	}
	// Always sort ascending internally; direction-aware slicing handles ordering.
	sortEventsByTimestamp(allEvents)

	// Parse the direction-aware token. For the first request (empty token),
	// derive direction from startFromHead and use the appropriate boundary.
	direction, offset, err := ParsePaginationToken(nextToken)
	if err != nil {
		return nil, nil, "", err
	}
	if nextToken == "" {
		if startFromHead {
			direction = PaginationForward
			offset = 0
		} else {
			direction = PaginationBackward
			offset = len(allEvents)
		}
	}

	totalLen := len(allEvents)
	var result []*OutputLogEvent
	var nextOffset int

	if offset > totalLen {
		offset = totalLen
	}

	if direction == PaginationForward {
		endIdx := offset + limit
		if limit <= 0 || endIdx > totalLen {
			endIdx = totalLen
		}
		if offset < totalLen {
			result = allEvents[offset:endIdx]
		}
		nextOffset = endIdx
	} else {
		// Backward: read [max(0, offset-limit), offset) in descending order.
		startIdx := offset - limit
		if limit <= 0 || startIdx < 0 {
			startIdx = 0
		}
		result = make([]*OutputLogEvent, 0, offset-startIdx)
		for i := offset - 1; i >= startIdx; i-- {
			if i >= 0 && i < totalLen {
				result = append(result, allEvents[i])
			}
		}
		nextOffset = startIdx
	}
	// Return an empty token at the boundary so clients can detect
	// end-of-data and stop paginating. Without this check the encoded
	// token (e.g. "f-<totalLen>" or "b-0") is non-empty and the AWS SDK
	// client loops forever requesting the next page.
	var newNextToken string
	if (direction == PaginationForward && nextOffset < totalLen) ||
		(direction == PaginationBackward && nextOffset > 0) {
		newNextToken = EncodePaginationToken(direction, nextOffset)
	}

	return result, searchedStreams, newNextToken, nil
}
