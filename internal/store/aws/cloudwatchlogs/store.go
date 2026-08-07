package cloudwatchlogs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/core/storage/chunk"
	pb "vorpalstacks/internal/pb/storage/storage_cloudwatchlogs"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/naming"
	"vorpalstacks/pkg/filterpattern"
)

// Store provides CloudWatch Logs storage operations.
type Store struct {
	*common.BaseStore
	tagStore     *common.TagStore
	ts           storage.TransactionalStorage
	arnBuilder   *svcarn.ARNBuilder
	region       string
	bucketName   string
	chunksDir    string
	chunkMutex   sync.Mutex
	activeChunks map[string]*activeChunk
	chunkCounter uint64
	subFilterMu  sync.Mutex
}

type activeChunk struct {
	entries []LogEntry
}

// NewStore creates a new CloudWatch Logs store.
func NewStore(store storage.BasicStorage, bucket storage.Bucket, accountID, region, dataPath string) (*Store, error) {
	baseStore := common.NewBaseStore(bucket, "logs")
	chunksDir := filepath.Join(dataPath, "logs-chunks")
	if err := os.MkdirAll(chunksDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs chunks directory: %w", err)
	}

	var ts storage.TransactionalStorage
	if txnStorage, ok := store.(storage.TransactionalStorage); ok {
		ts = txnStorage
	}

	return &Store{
		BaseStore:    baseStore,
		tagStore:     common.NewTagStoreWithRegion(store, "cloudwatchlogs", region),
		ts:           ts,
		arnBuilder:   svcarn.NewARNBuilder(accountID, region),
		region:       region,
		bucketName:   "logs-" + region,
		chunksDir:    chunksDir,
		activeChunks: make(map[string]*activeChunk),
	}, nil
}

// ARNBuilder returns the ARN builder for the store.
func (s *Store) ARNBuilder() *svcarn.ARNBuilder {
	return s.arnBuilder
}

// Tags returns the tag store for CloudWatch Logs log groups.
func (s *Store) Tags() *common.TagStore {
	return s.tagStore
}

func (s *Store) safeChunkPath(chunkPath string) (string, error) {
	return naming.ValidatePathWithinDir(s.chunksDir, chunkPath)
}

// CreateLogGroup creates a new CloudWatch Logs log group.
func (s *Store) CreateLogGroup(lg *LogGroup) error {
	key := s.logGroupKey(lg.Name)
	if s.Exists(key) {
		return ErrLogGroupAlreadyExists
	}
	lg.ARN = s.arnBuilder.CloudWatch().LogGroup(lg.Name)
	return s.PutProto(key, LogGroupToProto(lg))
}

// GetLogGroup retrieves a CloudWatch Logs log group by name.
func (s *Store) GetLogGroup(name string) (*LogGroup, error) {
	key := s.logGroupKey(name)
	var p pb.LogGroup
	if err := s.GetProto(key, &p); err != nil {
		return nil, ErrLogGroupNotFound
	}
	return ProtoToLogGroup(&p), nil
}

// PutLogGroup creates or updates a CloudWatch Logs log group.
func (s *Store) PutLogGroup(lg *LogGroup) error {
	key := s.logGroupKey(lg.Name)
	return s.PutProto(key, LogGroupToProto(lg))
}

// DeleteLogGroup deletes a CloudWatch Logs log group.
func (s *Store) DeleteLogGroup(name string) error {
	if _, err := s.GetLogGroup(name); err != nil {
		return err
	}

	if err := s.deleteAllLogStreams(name); err != nil {
		return err
	}

	chunks := s.ListChunksForLogGroup(name)
	for _, chunk := range chunks {
		cp, err := s.safeChunkPath(chunk.ChunkPath)
		if err != nil {
			return err
		}
		if err := os.Remove(cp); err != nil {
			return err
		}
	}
	metricFilters, _, _ := s.ListMetricFilters(name, "", "", 1000)
	for _, mf := range metricFilters {
		if err := s.Delete(s.metricFilterKey(name, mf.Name)); err != nil {
			logs.Error("Failed to delete metric filter", logs.String("name", mf.Name), logs.Err(err))
		}
	}
	subFilters, _ := s.ListSubscriptionFilters(name, "")
	for _, sf := range subFilters {
		if err := s.Delete(s.subscriptionFilterKey(name, sf.FilterName)); err != nil {
			logs.Error("Failed to delete subscription filter", logs.String("name", sf.FilterName), logs.Err(err))
		}
	}
	arn := s.arnBuilder.CloudWatch().LogGroup(name)
	if err := s.tagStore.Delete(arn); err != nil {
		return fmt.Errorf("failed to delete tags for log group %s: %w", name, err)
	}
	key := s.logGroupKey(name)
	return s.Delete(key)
}

func (s *Store) deleteAllLogStreams(logGroupName string) error {
	marker := ""
	for {
		streams, nextMarker, err := s.ListLogStreams(logGroupName, "", marker, 1000)
		if err != nil {
			return err
		}
		for _, stream := range streams {
			if err := s.DeleteLogStream(logGroupName, stream.Name); err != nil {
				return err
			}
		}
		if nextMarker == "" {
			return nil
		}
		marker = nextMarker
	}
}

// ListLogGroups lists CloudWatch Logs log groups with optional prefix and pagination.
func (s *Store) ListLogGroups(prefix, marker string, maxItems int) ([]*LogGroup, string, error) {
	if maxItems <= 0 {
		maxItems = 50
	}

	opts := common.ListOptions{
		Prefix:   "log-group:",
		Marker:   marker,
		MaxItems: maxItems,
	}

	result, err := common.ListProto[*pb.LogGroup](s.BaseStore, opts, func() *pb.LogGroup { return &pb.LogGroup{} }, func(lg *pb.LogGroup) bool {
		if prefix != "" && lg.Name != prefix && !strings.HasPrefix(lg.Name, prefix) {
			return false
		}
		return true
	})
	if err != nil {
		return nil, "", err
	}

	groups := make([]*LogGroup, len(result.Items))
	for i, p := range result.Items {
		groups[i] = ProtoToLogGroup(p)
	}
	return groups, result.NextMarker, nil
}

// CreateLogStream creates a new CloudWatch Logs log stream.
func (s *Store) CreateLogStream(ls *LogStream) error {
	if _, err := s.GetLogGroup(ls.LogGroupName); err != nil {
		return ErrLogGroupNotFound
	}

	key := s.logStreamKey(ls.LogGroupName, ls.Name)
	if s.Exists(key) {
		return ErrLogStreamAlreadyExists
	}

	ls.ARN = s.arnBuilder.CloudWatch().LogStream(ls.LogGroupName, ls.Name)
	ls.UploadSequenceToken = "0"
	return s.PutProto(key, LogStreamToProto(ls))
}

// GetLogStream retrieves a CloudWatch Logs log stream by group and stream name.
func (s *Store) GetLogStream(logGroupName, logStreamName string) (*LogStream, error) {
	key := s.logStreamKey(logGroupName, logStreamName)
	var p pb.LogStream
	if err := s.GetProto(key, &p); err != nil {
		return nil, ErrLogStreamNotFound
	}
	return ProtoToLogStream(&p), nil
}

// DeleteLogStream deletes a CloudWatch Logs log stream.
func (s *Store) DeleteLogStream(logGroupName, logStreamName string) error {
	if _, err := s.GetLogStream(logGroupName, logStreamName); err != nil {
		return err
	}

	streamKey := fmt.Sprintf("%s:%s", logGroupName, logStreamName)
	s.chunkMutex.Lock()
	delete(s.activeChunks, streamKey)
	s.chunkMutex.Unlock()

	chunks := s.ListChunksForStream(logGroupName, logStreamName)
	for _, chunk := range chunks {
		cp, err := s.safeChunkPath(chunk.ChunkPath)
		if err != nil {
			logs.Warn("Path traversal attempt in chunk", logs.String("chunkPath", chunk.ChunkPath), logs.Err(err))
			continue
		}
		if err := os.Remove(cp); err != nil {
			logs.Error("Failed to remove chunk file", logs.String("chunkPath", chunk.ChunkPath), logs.Err(err))
		}
		deleteKey := s.chunkIndexKey(logGroupName, logStreamName, chunk.ChunkID)
		if err := s.Delete(deleteKey); err != nil {
			logs.Error("Failed to delete chunk index", logs.String("key", deleteKey), logs.Err(err))
		}
	}
	key := s.logStreamKey(logGroupName, logStreamName)
	return s.Delete(key)
}

// ListLogStreams lists CloudWatch Logs log streams for a given log group.
func (s *Store) ListLogStreams(logGroupName, prefix, marker string, maxItems int) ([]*LogStream, string, error) {
	if maxItems <= 0 {
		maxItems = 50
	}

	streamPrefix := s.logStreamKey(logGroupName, "")
	if prefix != "" {
		streamPrefix += escapePath(prefix)
	}

	opts := common.ListOptions{
		Prefix:   streamPrefix,
		Marker:   marker,
		MaxItems: maxItems,
	}

	result, err := common.ListProto[*pb.LogStream](s.BaseStore, opts, func() *pb.LogStream { return &pb.LogStream{} }, nil)
	if err != nil {
		return nil, "", err
	}

	streams := make([]*LogStream, len(result.Items))
	for i, p := range result.Items {
		streams[i] = ProtoToLogStream(p)
	}
	return streams, result.NextMarker, nil
}

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
	// fails, the chunk files are removed to prevent orphans (M-S1).
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

// pendingChunkIndex holds a chunk index entry that has been written to
// disk but not yet committed to Pebble. The caller must write the index
// entry to Pebble (ideally inside a transaction together with the
// LogGroup/LogStream update) or remove chunkPath on failure (M-S1).
type pendingChunkIndex struct {
	meta      *ChunkMeta
	indexKey  string
	chunkPath string
}

// updateLogGroupStreamAndChunks atomically writes the LogGroup, LogStream,
// and any pending chunk indexes using a storage transaction. Falls back to
// sequential writes if the storage backend does not support transactions.
// If this returns an error, the caller is responsible for removing the
// orphaned chunk files referenced by pendingChunks (M-S1).
func (s *Store) updateLogGroupStreamAndChunks(logGroupName string, lg *LogGroup, logStreamName string, ls *LogStream, pendingChunks []pendingChunkIndex) error {
	if s.ts == nil {
		if err := s.PutLogGroup(lg); err != nil {
			return err
		}
		if err := s.PutProto(s.logStreamKey(logGroupName, logStreamName), LogStreamToProto(ls)); err != nil {
			return err
		}
		for _, pc := range pendingChunks {
			if err := s.PutProto(pc.indexKey, ChunkMetaToProto(pc.meta)); err != nil {
				return err
			}
		}
		return nil
	}
	ctx := context.Background()
	return s.ts.Update(ctx, func(txn storage.Transaction) error {
		lgData, err := proto.Marshal(LogGroupToProto(lg))
		if err != nil {
			return err
		}
		if err := txn.Bucket(s.bucketName).Put([]byte(s.logGroupKey(logGroupName)), lgData); err != nil {
			return err
		}
		lsData, err := proto.Marshal(LogStreamToProto(ls))
		if err != nil {
			return err
		}
		if err := txn.Bucket(s.bucketName).Put([]byte(s.logStreamKey(logGroupName, logStreamName)), lsData); err != nil {
			return err
		}
		for _, pc := range pendingChunks {
			idxData, err := proto.Marshal(ChunkMetaToProto(pc.meta))
			if err != nil {
				return err
			}
			if err := txn.Bucket(s.bucketName).Put([]byte(pc.indexKey), idxData); err != nil {
				return err
			}
		}
		return nil
	})
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

// prepareChunkFlush writes the chunk file to disk and prepares the index
// entry WITHOUT committing it to Pebble. The caller must commit the
// returned pendingChunkIndex to Pebble (ideally inside a transaction) or
// remove chunkPath on failure (M-S1).
func (s *Store) prepareChunkFlush(logGroupName, logStreamName string, ac *activeChunk) (pendingChunkIndex, error) {
	if len(ac.entries) == 0 {
		return pendingChunkIndex{}, nil
	}

	chunkSeq := atomic.AddUint64(&s.chunkCounter, 1)
	chunkID := fmt.Sprintf("%d-%d-%d", ac.entries[0].Timestamp, len(ac.entries), chunkSeq)

	actualPath, header, err := s.writeChunkFile(ac.entries)
	if err != nil {
		return pendingChunkIndex{}, err
	}

	relPath := actualPath
	if filepath.IsAbs(actualPath) && s.chunksDir != "" {
		if rel, err := filepath.Rel(s.chunksDir, actualPath); err == nil {
			relPath = rel
		}
	}

	meta := &ChunkMeta{
		ChunkID:    chunkID,
		LogGroup:   logGroupName,
		LogStream:  logStreamName,
		MinTs:      header.MinTs,
		MaxTs:      header.MaxTs,
		EntryCount: int(header.EntryCount),
		ChunkPath:  relPath,
	}

	return pendingChunkIndex{
		meta:      meta,
		indexKey:  s.chunkIndexKey(logGroupName, logStreamName, chunkID),
		chunkPath: actualPath,
	}, nil
}

// flushChunk writes a chunk to disk and commits its index entry to Pebble.
// Used by callers that do not need transactional atomicity with metadata
// updates (e.g. flushIfNeeded before reads). PutLogEvents uses
// prepareChunkFlush instead so the index can be committed inside the
// metadata transaction (M-S1).
func (s *Store) flushChunk(logGroupName, logStreamName string, ac *activeChunk) error {
	pi, err := s.prepareChunkFlush(logGroupName, logStreamName, ac)
	if err != nil {
		return err
	}
	if pi.meta == nil {
		return nil
	}
	if err := s.PutProto(pi.indexKey, ChunkMetaToProto(pi.meta)); err != nil {
		// Remove the orphaned chunk file so it does not leak storage
		// with no index entry pointing to it.
		if rmErr := os.Remove(pi.chunkPath); rmErr != nil && !os.IsNotExist(rmErr) {
			logs.Error("Failed to remove orphaned chunk file after index failure",
				logs.String("path", pi.chunkPath), logs.Err(rmErr))
		}
		return err
	}
	return nil
}

func (s *Store) writeChunkFile(entries []LogEntry) (string, *chunk.Header, error) {
	chunkEntries := make([]chunk.Entry, len(entries))
	for i, e := range entries {
		chunkEntries[i] = chunk.SimpleEntry{
			Ts:          e.Timestamp,
			IngestionTs: e.IngestionTime,
			Msg:         []byte(e.Message),
		}
	}

	opts := &chunk.WriterOptions{
		ChunksDir:   s.chunksDir,
		Encoding:    chunk.EncodingZstd,
		SyncOnWrite: false,
	}

	w := chunk.NewWriter(opts)
	if err := w.WriteBatch(chunkEntries); err != nil {
		return "", nil, err
	}

	actualPath, err := w.Flush()
	if err != nil {
		return "", nil, err
	}

	header, err := chunk.ReadHeader(actualPath)
	if err != nil {
		return "", nil, err
	}

	return actualPath, header, nil
}

func (s *Store) readChunkFile(chunkPath string) ([]LogEntry, error) {
	fullPath := chunkPath
	if !filepath.IsAbs(chunkPath) {
		p, err := s.safeChunkPath(chunkPath)
		if err != nil {
			return nil, err
		}
		fullPath = p
	}

	r := chunk.NewReader(&chunk.ReaderOptions{ChunksDir: s.chunksDir})
	chunkEntries, err := r.Read(fullPath)
	if err != nil {
		return nil, err
	}

	entries := make([]LogEntry, len(chunkEntries))
	for i, ce := range chunkEntries {
		var ingestionTs int64
		if ig, ok := ce.(chunk.Ingestible); ok {
			ingestionTs = ig.IngestionTimeUnixMilli()
		}
		entries[i] = LogEntry{
			Timestamp:     ce.Timestamp(),
			Message:       string(ce.Message()),
			IngestionTime: ingestionTs,
		}
	}

	return entries, nil
}

// ListChunksForStream lists chunk metadata for a specific log stream.
func (s *Store) ListChunksForStream(logGroupName, logStreamName string) []*ChunkMeta {
	prefix := s.chunkIndexKey(logGroupName, logStreamName, "")
	var chunks []*ChunkMeta

	if err := s.ScanPrefix(prefix, func(key string, value []byte) error {
		if !bytes.HasPrefix([]byte(key), []byte(prefix)) {
			return nil
		}
		var p pb.ChunkMeta
		if err := proto.Unmarshal(value, &p); err == nil {
			chunks = append(chunks, ProtoToChunkMeta(&p))
		}
		return nil
	}); err != nil {
		logs.Error("Failed to scan chunks for stream", logs.String("logGroup", logGroupName), logs.String("logStream", logStreamName), logs.Err(err))
	}

	return chunks
}

// ListChunksForLogGroup lists chunk metadata for all streams in a log group.
func (s *Store) ListChunksForLogGroup(logGroupName string) []*ChunkMeta {
	prefix := "chunk:" + escapePath(logGroupName) + ":"
	var chunks []*ChunkMeta

	if err := s.ScanPrefix(prefix, func(key string, value []byte) error {
		if !bytes.HasPrefix([]byte(key), []byte(prefix)) {
			return nil
		}
		var p pb.ChunkMeta
		if err := proto.Unmarshal(value, &p); err == nil {
			chunks = append(chunks, ProtoToChunkMeta(&p))
		}
		return nil
	}); err != nil {
		logs.Error("Failed to scan chunks for log group", logs.String("logGroup", logGroupName), logs.Err(err))
	}

	return chunks
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

func (s *Store) logGroupKey(name string) string {
	return "log-group:" + name
}

func (s *Store) logStreamKey(logGroupName, logStreamName string) string {
	return "log-stream:" + escapePath(logGroupName) + ":" + logStreamName
}

func (s *Store) chunkIndexKey(logGroupName, logStreamName, chunkID string) string {
	return "chunk:" + escapePath(logGroupName) + ":" + escapePath(logStreamName) + ":" + chunkID
}

func escapePath(path string) string {
	result := make([]byte, 0, len(path)*2)
	for i := 0; i < len(path); i++ {
		c := path[i]
		if c == '/' || c == ':' || c == '\\' {
			result = append(result, '%')
			result = append(result, hexChar(c>>4))
			result = append(result, hexChar(c&0x0F))
		} else {
			result = append(result, c)
		}
	}
	return string(result)
}

func hexChar(b byte) byte {
	if b < 10 {
		return '0' + b
	}
	return 'a' + b - 10
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

// PutSubscriptionFilter creates or updates a subscription filter.
func (s *Store) PutSubscriptionFilter(filter *SubscriptionFilter) error {
	key := s.subscriptionFilterKey(filter.LogGroupName, filter.FilterName)
	var existing pb.SubscriptionFilter
	if s.GetProto(key, &existing) == nil {
		filter.CreationTime = time.UnixMilli(existing.CreatedAt)
	} else {
		filter.CreationTime = time.Now().UTC()
	}
	return s.PutProto(key, SubscriptionFilterToProto(filter))
}

// PutSubscriptionFilterWithLimitCheck atomically checks the per-log-group
// subscription filter limit (max 2 distinct filter names) and creates or
// updates the filter. This prevents race conditions where concurrent
// PutSubscriptionFilter calls could exceed the limit.
func (s *Store) PutSubscriptionFilterWithLimitCheck(filter *SubscriptionFilter, maxPerGroup int) error {
	s.subFilterMu.Lock()
	defer s.subFilterMu.Unlock()

	filters, err := s.ListSubscriptionFilters(filter.LogGroupName, "")
	if err != nil {
		return err
	}

	existingCount := 0
	isUpdate := false
	for _, f := range filters {
		if f.FilterName == filter.FilterName {
			isUpdate = true
			continue
		}
		existingCount++
	}

	if !isUpdate && existingCount >= maxPerGroup {
		return ErrLimitExceeded
	}

	return s.PutSubscriptionFilter(filter)
}

// DeleteSubscriptionFilter deletes a subscription filter.
func (s *Store) DeleteSubscriptionFilter(logGroupName, filterName string) error {
	key := s.subscriptionFilterKey(logGroupName, filterName)
	if !s.Exists(key) {
		return ErrSubscriptionFilterNotFound
	}
	return s.Delete(key)
}

// GetSubscriptionFilter retrieves a subscription filter by name.
func (s *Store) GetSubscriptionFilter(logGroupName, filterName string) (*SubscriptionFilter, error) {
	key := s.subscriptionFilterKey(logGroupName, filterName)
	var p pb.SubscriptionFilter
	if err := s.GetProto(key, &p); err != nil {
		return nil, ErrSubscriptionFilterNotFound
	}
	return ProtoToSubscriptionFilter(&p), nil
}

// ListSubscriptionFilters lists subscription filters for a log group.
func (s *Store) ListSubscriptionFilters(logGroupName, filterNamePrefix string) ([]*SubscriptionFilter, error) {
	prefix := s.subscriptionFilterKey(logGroupName, "")
	var filters []*SubscriptionFilter

	if err := s.ScanPrefix(prefix, func(key string, value []byte) error {
		var p pb.SubscriptionFilter
		if err := proto.Unmarshal(value, &p); err == nil {
			filter := ProtoToSubscriptionFilter(&p)
			if filterNamePrefix == "" || strings.HasPrefix(filter.FilterName, filterNamePrefix) {
				filters = append(filters, filter)
			}
		}
		return nil
	}); err != nil {
		logs.Error("Failed to scan subscription filters", logs.String("logGroup", logGroupName), logs.Err(err))
	}

	return filters, nil
}

func (s *Store) subscriptionFilterKey(logGroupName, filterName string) string {
	return "subscription-filter:" + escapePath(logGroupName) + ":" + filterName
}

func (s *Store) destinationKey(name string) string {
	return "destination:" + name
}

// PutDestination creates a new CloudWatch Logs destination.
func (s *Store) PutDestination(dest *Destination) error {
	key := s.destinationKey(dest.Name)
	if dest.CreationTime == 0 {
		if existing, err := s.GetDestination(dest.Name); err == nil {
			dest.CreationTime = existing.CreationTime
		} else {
			dest.CreationTime = time.Now().UnixMilli()
		}
	}
	return s.PutProto(key, DestinationToProto(dest))
}

// GetDestination retrieves a CloudWatch Logs destination by name.
func (s *Store) GetDestination(name string) (*Destination, error) {
	key := s.destinationKey(name)
	var p pb.Destination
	if err := s.GetProto(key, &p); err != nil {
		return nil, ErrDestinationNotFound
	}
	return ProtoToDestination(&p), nil
}

// DeleteDestination deletes a CloudWatch Logs destination by name.
func (s *Store) DeleteDestination(name string) error {
	key := s.destinationKey(name)
	if !s.Exists(key) {
		return ErrDestinationNotFound
	}
	return s.Delete(key)
}

// PutDestinationPolicy sets the resource-based access policy for a CloudWatch Logs destination.
func (s *Store) PutDestinationPolicy(name, accessPolicy string) error {
	dest, err := s.GetDestination(name)
	if err != nil {
		return err
	}
	dest.AccessPolicy = accessPolicy
	return s.PutProto(s.destinationKey(name), DestinationToProto(dest))
}

// ListDestinations returns all CloudWatch Logs destinations, optionally filtered by name prefix.
func (s *Store) ListDestinations(prefix string) ([]*Destination, error) {
	destPrefix := "destination:"
	var destinations []*Destination

	if err := s.ScanPrefix(destPrefix, func(key string, value []byte) error {
		var p pb.Destination
		if err := proto.Unmarshal(value, &p); err != nil {
			return nil
		}
		dest := ProtoToDestination(&p)
		if prefix == "" || strings.HasPrefix(dest.Name, prefix) {
			destinations = append(destinations, dest)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return destinations, nil
}

// PurgeExpiredChunks removes expired log chunks from a log group based on retention policy.
func (s *Store) PurgeExpiredChunks(logGroupName string, cutoffTime int64) (int64, error) {
	lg, err := s.GetLogGroup(logGroupName)
	if err != nil {
		return 0, err
	}

	chunks := s.ListChunksForLogGroup(logGroupName)
	var totalBytesRemoved int64

	for _, chunk := range chunks {
		if chunk.MaxTs >= cutoffTime {
			continue
		}

		chunkPath, err := s.safeChunkPath(chunk.ChunkPath)
		if err != nil {
			continue
		}
		var fileSize int64
		if info, err := os.Stat(chunkPath); err == nil {
			fileSize = info.Size()
		}

		if err := os.Remove(chunkPath); err != nil && !os.IsNotExist(err) {
			continue
		}

		indexKey := s.chunkIndexKey(logGroupName, chunk.LogStream, chunk.ChunkID)
		if err := s.Delete(indexKey); err != nil {
			logs.Error("Failed to delete chunk index", logs.String("key", indexKey), logs.Err(err))
		}

		totalBytesRemoved += fileSize
	}

	if totalBytesRemoved > 0 {
		lg.StoredBytes -= totalBytesRemoved
		if lg.StoredBytes < 0 {
			lg.StoredBytes = 0
		}
		if err := s.PutLogGroup(lg); err != nil {
			return totalBytesRemoved, err
		}
	}

	return totalBytesRemoved, nil
}

// PurgeAllExpiredChunks purges expired chunks from all log groups based on their retention policies.
func (s *Store) PurgeAllExpiredChunks() error {
	now := time.Now().UnixMilli()

	marker := ""
	for {
		groups, nextMarker, err := s.ListLogGroups("", marker, 1000)
		if err != nil {
			return err
		}

		for _, lg := range groups {
			if lg.RetentionInDays <= 0 {
				continue
			}

			cutoffTime := now - int64(lg.RetentionInDays)*24*60*60*1000
			if _, err := s.PurgeExpiredChunks(lg.Name, cutoffTime); err != nil {
				continue
			}
		}

		if nextMarker == "" {
			break
		}
		marker = nextMarker
	}

	return nil
}

// --- Resource Policy ---

func (s *Store) resourcePolicyKey(policyName string) string {
	return "resource-policy:" + policyName
}

func (s *Store) PutResourcePolicy(policy *ResourcePolicy) error {
	policy.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.resourcePolicyKey(policy.PolicyName), policy)
}

func (s *Store) GetResourcePolicy(policyName string) (*ResourcePolicy, error) {
	var p ResourcePolicy
	if err := s.Get(s.resourcePolicyKey(policyName), &p); err != nil {
		return nil, ErrResourceNotFound
	}
	return &p, nil
}

func (s *Store) DeleteResourcePolicy(policyName string) error {
	if !s.Exists(s.resourcePolicyKey(policyName)) {
		return ErrResourceNotFound
	}
	return s.Delete(s.resourcePolicyKey(policyName))
}

func (s *Store) ListResourcePolicies(resourceArn string) ([]*ResourcePolicy, error) {
	var policies []*ResourcePolicy
	if err := s.ScanPrefix("resource-policy:", func(key string, value []byte) error {
		var p ResourcePolicy
		if err := json.Unmarshal(value, &p); err != nil {
			return nil
		}
		if resourceArn == "" || p.ResourceArn == resourceArn {
			policies = append(policies, &p)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return policies, nil
}

// --- Account Policy ---

func (s *Store) accountPolicyKey(policyType, policyName string) string {
	return "account-policy:" + policyType + ":" + policyName
}

func (s *Store) PutAccountPolicy(policy *AccountPolicy) error {
	policy.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.accountPolicyKey(policy.PolicyType, policy.PolicyName), policy)
}

func (s *Store) DeleteAccountPolicyEntry(policyType, policyName string) error {
	key := s.accountPolicyKey(policyType, policyName)
	if !s.Exists(key) {
		return ErrResourceNotFound
	}
	return s.Delete(key)
}

func (s *Store) ListAccountPolicies(policyType, policyName string) ([]*AccountPolicy, error) {
	var policies []*AccountPolicy
	scanPrefix := "account-policy:"
	if policyType != "" {
		scanPrefix = "account-policy:" + policyType + ":"
	}
	if err := s.ScanPrefix(scanPrefix, func(key string, value []byte) error {
		var p AccountPolicy
		if err := json.Unmarshal(value, &p); err != nil {
			return nil
		}
		if policyName != "" && p.PolicyName != policyName {
			return nil
		}
		policies = append(policies, &p)
		return nil
	}); err != nil {
		return nil, err
	}
	return policies, nil
}

// --- Data Protection Policy ---

func (s *Store) dataProtectionPolicyKey(logGroupIdentifier string) string {
	return "data-protection-policy:" + escapePath(logGroupIdentifier)
}

func (s *Store) PutDataProtectionPolicy(policy *DataProtectionPolicy) error {
	policy.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.dataProtectionPolicyKey(policy.LogGroupIdentifier), policy)
}

func (s *Store) GetDataProtectionPolicy(logGroupIdentifier string) (*DataProtectionPolicy, error) {
	var p DataProtectionPolicy
	if err := s.Get(s.dataProtectionPolicyKey(logGroupIdentifier), &p); err != nil {
		return nil, ErrResourceNotFound
	}
	return &p, nil
}

func (s *Store) DeleteDataProtectionPolicy(logGroupIdentifier string) error {
	key := s.dataProtectionPolicyKey(logGroupIdentifier)
	if !s.Exists(key) {
		return ErrResourceNotFound
	}
	return s.Delete(key)
}

// --- Query Definition ---

func (s *Store) queryDefinitionKey(id string) string {
	return "query-definition:" + id
}

func (s *Store) PutQueryDefinitionEntry(qd *QueryDefinition) error {
	qd.LastModified = time.Now().UTC().UnixMilli()
	return s.Put(s.queryDefinitionKey(qd.QueryDefinitionId), qd)
}

func (s *Store) DeleteQueryDefinitionEntry(id string) error {
	key := s.queryDefinitionKey(id)
	if !s.Exists(key) {
		return ErrResourceNotFound
	}
	return s.Delete(key)
}

func (s *Store) ListQueryDefinitions(namePrefix string) ([]*QueryDefinition, error) {
	var defs []*QueryDefinition
	if err := s.ScanPrefix("query-definition:", func(key string, value []byte) error {
		var qd QueryDefinition
		if err := json.Unmarshal(value, &qd); err != nil {
			return nil
		}
		if namePrefix == "" || strings.HasPrefix(qd.Name, namePrefix) {
			defs = append(defs, &qd)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return defs, nil
}

// --- Export Task ---

func (s *Store) exportTaskKey(taskId string) string {
	return "export-task:" + taskId
}

func (s *Store) PutExportTask(task *ExportTask) error {
	return s.Put(s.exportTaskKey(task.TaskId), task)
}

func (s *Store) GetExportTask(taskId string) (*ExportTask, error) {
	var t ExportTask
	if err := s.Get(s.exportTaskKey(taskId), &t); err != nil {
		return nil, ErrResourceNotFound
	}
	return &t, nil
}

func (s *Store) DeleteExportTask(taskId string) error {
	return s.Delete(s.exportTaskKey(taskId))
}

func (s *Store) ListExportTasks(statusCode string) ([]*ExportTask, error) {
	var tasks []*ExportTask
	if err := s.ScanPrefix("export-task:", func(key string, value []byte) error {
		var t ExportTask
		if err := json.Unmarshal(value, &t); err != nil {
			return nil
		}
		if statusCode == "" || t.Status == statusCode {
			tasks = append(tasks, &t)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return tasks, nil
}

// --- Import Task ---

func (s *Store) importTaskKey(importId string) string {
	return "import-task:" + importId
}

func (s *Store) PutImportTask(task *ImportTask) error {
	task.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.importTaskKey(task.ImportId), task)
}

func (s *Store) GetImportTask(importId string) (*ImportTask, error) {
	var t ImportTask
	if err := s.Get(s.importTaskKey(importId), &t); err != nil {
		return nil, ErrResourceNotFound
	}
	return &t, nil
}

func (s *Store) ListImportTasks(importId, status, sourceArn string) ([]*ImportTask, error) {
	var tasks []*ImportTask
	if err := s.ScanPrefix("import-task:", func(key string, value []byte) error {
		var t ImportTask
		if err := json.Unmarshal(value, &t); err != nil {
			return nil
		}
		if importId != "" && t.ImportId != importId {
			return nil
		}
		if status != "" && t.ImportStatus != status {
			return nil
		}
		if sourceArn != "" && t.ImportSourceArn != sourceArn {
			return nil
		}
		tasks = append(tasks, &t)
		return nil
	}); err != nil {
		return nil, err
	}
	return tasks, nil
}

// --- Scheduled Query ---

func (s *Store) scheduledQueryKey(id string) string {
	return "scheduled-query:" + id
}

func (s *Store) PutScheduledQuery(sq *ScheduledQuery) error {
	sq.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.scheduledQueryKey(sq.Id), sq)
}

func (s *Store) GetScheduledQuery(id string) (*ScheduledQuery, error) {
	var sq ScheduledQuery
	if err := s.Get(s.scheduledQueryKey(id), &sq); err != nil {
		return nil, ErrResourceNotFound
	}
	return &sq, nil
}

func (s *Store) DeleteScheduledQuery(id string) error {
	key := s.scheduledQueryKey(id)
	if !s.Exists(key) {
		return ErrResourceNotFound
	}
	return s.Delete(key)
}

func (s *Store) ListScheduledQueries(state string) ([]*ScheduledQuery, error) {
	var queries []*ScheduledQuery
	if err := s.ScanPrefix("scheduled-query:", func(key string, value []byte) error {
		var sq ScheduledQuery
		if err := json.Unmarshal(value, &sq); err != nil {
			return nil
		}
		if state == "" || sq.State == state {
			queries = append(queries, &sq)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return queries, nil
}

// --- Scheduled Query Execution History ---

func (s *Store) scheduledQueryExecutionKey(sqId string, triggerTime int64) string {
	return "sq-execution:" + sqId + ":" + strconv.FormatInt(triggerTime, 10)
}

func (s *Store) PutScheduledQueryExecution(exec *ScheduledQueryExecution) error {
	return s.Put(s.scheduledQueryExecutionKey(exec.ScheduledQueryId, exec.TriggerTime), exec)
}

func (s *Store) ListScheduledQueryExecutions(sqId string, startTime, endTime int64) ([]*ScheduledQueryExecution, error) {
	var execs []*ScheduledQueryExecution
	prefix := "sq-execution:" + sqId + ":"
	if err := s.ScanPrefix(prefix, func(key string, value []byte) error {
		var exec ScheduledQueryExecution
		if err := json.Unmarshal(value, &exec); err != nil {
			return nil
		}
		if startTime > 0 && exec.TriggerTime < startTime {
			return nil
		}
		if endTime > 0 && exec.TriggerTime > endTime {
			return nil
		}
		execs = append(execs, &exec)
		return nil
	}); err != nil {
		return nil, err
	}
	return execs, nil
}
