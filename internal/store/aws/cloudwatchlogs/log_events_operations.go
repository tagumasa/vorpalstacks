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

	streamKey := fmt.Sprintf("%s:%s", logGroupName, logStreamName)
	ac, exists := s.activeChunks[streamKey]
	if !exists {
		ac = &activeChunk{
			entries: make([]LogEntry, 0, MaxChunkSize),
		}
		s.activeChunks[streamKey] = ac
	}

	// Collect chunk index entries that must be committed atomically with
	// the LogGroup/LogStream metadata update. If the metadata transaction
	// fails, the chunk files are removed to prevent orphans.
	var pendingChunks []pendingChunkIndex

	// Flush existing entries first if appending the new batch would
	// exceed MaxChunkSize, preventing a single chunk from growing
	// up to ~2× MaxChunkSize.
	if len(ac.entries)+len(events) > MaxChunkSize && len(ac.entries) > 0 {
		pi, err := s.prepareChunkFlush(logGroupName, logStreamName, ac)
		if err != nil {
			return "", err
		}
		if pi.meta != nil {
			pendingChunks = append(pendingChunks, pi)
		}
		ac.entries = ac.entries[:0]
	}

	ac.entries = append(ac.entries, events...)

	if len(ac.entries) >= MaxChunkSize {
		pi, err := s.prepareChunkFlush(logGroupName, logStreamName, ac)
		if err != nil {
			return "", err
		}
		if pi.meta != nil {
			pendingChunks = append(pendingChunks, pi)
		}
		delete(s.activeChunks, streamKey)
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
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp < events[j].Timestamp
	})
}

func (s *Store) flushIfNeeded(logGroupName, logStreamName string) {
	streamKey := fmt.Sprintf("%s:%s", logGroupName, logStreamName)
	s.chunkMutex.Lock()
	defer s.chunkMutex.Unlock()

	if ac, exists := s.activeChunks[streamKey]; exists && len(ac.entries) > 0 {
		if err := s.flushChunk(logGroupName, logStreamName, ac); err != nil {
			logs.Error("Failed to flush chunk", logs.String("logGroup", logGroupName), logs.String("logStream", logStreamName), logs.Err(err))
		}
		ac.entries = ac.entries[:0]
	}
}

func (s *Store) flushAllForLogGroup(logGroupName string) {
	streams, _, err := s.ListLogStreams(logGroupName, "", "", 10000)
	if err != nil {
		return
	}
	for _, stream := range streams {
		s.flushIfNeeded(logGroupName, stream.Name)
	}
}

// GetLogEvents retrieves log events from a CloudWatch Logs log stream.
func (s *Store) GetLogEvents(logGroupName, logStreamName string, startTime, endTime int64, limit int, startFromHead bool, nextToken string) ([]*OutputLogEvent, string, string, error) {
	if _, err := s.GetLogGroup(logGroupName); err != nil {
		return nil, "", "", ErrLogGroupNotFound
	}

	if _, err := s.GetLogStream(logGroupName, logStreamName); err != nil {
		return nil, "", "", ErrLogStreamNotFound
	}

	s.flushIfNeeded(logGroupName, logStreamName)

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

	for _, streamName := range logStreamNames {
		s.flushIfNeeded(logGroupName, streamName)
	}

	if len(logStreamNames) == 0 {
		s.flushAllForLogGroup(logGroupName)
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
