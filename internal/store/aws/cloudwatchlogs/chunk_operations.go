package cloudwatchlogs

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/core/storage/chunk"
	pb "vorpalstacks/internal/pb/storage/storage_cloudwatchlogs"
)

// pendingChunkIndex holds a chunk index entry that has been written to
// disk but not yet committed to Pebble. The caller must write the index
// entry to Pebble (ideally inside a transaction together with the
// LogGroup/LogStream update) or remove chunkPath on failure.
type pendingChunkIndex struct {
	meta      *ChunkMeta
	indexKey  string
	chunkPath string
}

// updateLogGroupStreamAndChunks atomically writes the LogGroup, LogStream,
// and any pending chunk indexes using a storage transaction. Falls back to
// sequential writes if the storage backend does not support transactions.
// If this returns an error, the caller is responsible for removing the
// orphaned chunk files referenced by pendingChunks.
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

// prepareChunkFlush writes the chunk file to disk and prepares the index
// entry WITHOUT committing it to Pebble. The caller must commit the
// returned pendingChunkIndex to Pebble (ideally inside a transaction) or
// remove chunkPath on failure.
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
// metadata transaction.
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

// readChunkHeader is indirected so tests can exercise the ReadHeader
// failure path of writeChunkFile.
var readChunkHeader = chunk.ReadHeader

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

	header, err := readChunkHeader(actualPath)
	if err != nil {
		// The file is on disk but will never be indexed; remove it so it
		// cannot leak storage as an orphan no sweep knows about.
		if rmErr := os.Remove(actualPath); rmErr != nil && !os.IsNotExist(rmErr) {
			logs.Error("Failed to remove orphaned chunk file after header read failure",
				logs.String("path", actualPath), logs.Err(rmErr))
		}
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
		if err := proto.Unmarshal(value, &p); err != nil {
			// A record that fails to decode is invisible to the caller;
			// log it so chunk corruption is discoverable instead of
			// silently shrinking the chunk list.
			logs.Warn("Corrupt chunk metadata record",
				logs.String("key", key), logs.Err(err))
			return nil
		}
		chunks = append(chunks, ProtoToChunkMeta(&p))
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
		if err := proto.Unmarshal(value, &p); err != nil {
			// See ListChunksForStream: corruption must stay visible.
			logs.Warn("Corrupt chunk metadata record",
				logs.String("key", key), logs.Err(err))
			return nil
		}
		chunks = append(chunks, ProtoToChunkMeta(&p))
		return nil
	}); err != nil {
		logs.Error("Failed to scan chunks for log group", logs.String("logGroup", logGroupName), logs.Err(err))
	}

	return chunks
}
