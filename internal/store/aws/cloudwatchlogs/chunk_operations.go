package cloudwatchlogs

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
		// Fallback (non-transactional) commits are sequential writes, so a
		// failure midway must undo the writes that already landed: the
		// index records written before the failure are removed here (the
		// caller removes the chunk files), and the caller restores the
		// LogGroup/LogStream records to their pre-mutation state — without
		// both compensations a failed call leaves StoredBytes and
		// sequence-token drift that no later operation reconciles.
		if err := s.putLogGroupLocked(lg); err != nil {
			return err
		}
		if err := s.PutProto(s.logStreamKey(logGroupName, logStreamName), LogStreamToProto(ls)); err != nil {
			return err
		}
		var written []string
		for _, pc := range pendingChunks {
			if err := s.PutProto(pc.indexKey, ChunkMetaToProto(pc.meta)); err != nil {
				for _, key := range written {
					if rmErr := s.Delete(key); rmErr != nil {
						logs.Error("Failed to roll back chunk index record",
							logs.String("key", key), logs.Err(rmErr))
					}
				}
				return err
			}
			written = append(written, pc.indexKey)
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
func (s *Store) prepareChunkFlush(logGroupName, logStreamName string, entries []LogEntry) (pendingChunkIndex, error) {
	if len(entries) == 0 {
		return pendingChunkIndex{}, nil
	}

	chunkSeq := atomic.AddUint64(&s.chunkCounter, 1)
	// Fixed-width components in timestamp-sequence-entryCount order keep
	// the lexicographic order of the chunk index keys equal to the
	// chronological order of the chunks: with variable-width numbers, or
	// with the entry count ahead of the sequence, same-timestamp reads
	// return chunks out of ingestion order.
	chunkID := fmt.Sprintf("%013d-%010d-%06d", entries[0].Timestamp, chunkSeq, len(entries))

	actualPath, header, err := s.writeChunkFile(entries)
	if err != nil {
		return pendingChunkIndex{}, err
	}

	// The index records the relative form when the chunks directory
	// allows the conversion; a server started with a relative data path
	// cannot convert an absolute writer path against it, so both forms
	// occur in stored records and both are first-class at read time.
	relPath := actualPath
	if filepath.IsAbs(actualPath) && s.chunksDir != "" {
		if rel, err := filepath.Rel(s.chunksDir, actualPath); err == nil {
			relPath = rel
		}
	}

	// The recorded byte size is the ingestion accounting basis (the sum of
	// message lengths PutLogEvents added to StoredBytes), not the file
	// size: the file is compressed, so removal paths that decremented by
	// file size would drift against the increment.
	var msgBytes int64
	var maxIngestion int64
	for _, e := range entries {
		msgBytes += int64(len(e.Message))
		if e.IngestionTime > maxIngestion {
			maxIngestion = e.IngestionTime
		}
	}

	meta := &ChunkMeta{
		ChunkID:        chunkID,
		LogGroupName:   logGroupName,
		LogStream:      logStreamName,
		MinTs:          header.MinTs,
		MaxTs:          header.MaxTs,
		MaxIngestionTs: maxIngestion,
		EntryCount:     int(header.EntryCount),
		ChunkPath:      relPath,
		ByteSize:       msgBytes,
	}

	return pendingChunkIndex{
		meta:      meta,
		indexKey:  s.chunkIndexKey(logGroupName, logStreamName, chunkID),
		chunkPath: actualPath,
	}, nil
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

	// Each acknowledged PutLogEvents must be durable across a crash, not
	// only a process restart: the chunk file is fsynced before the chunk
	// index transaction commits, so the index can never point at data
	// that a crash lost. The cost is one fsync per chunk, and a chunk is
	// written exactly once per PutLogEvents call. ChunkSize must exceed
	// the maximal batch: the writer auto-flushes once its buffer reaches
	// ChunkSize, and with a maximal batch that auto-flush would empty the
	// buffer before the caller's explicit Flush, which then returns an
	// empty path and the header read fails.
	opts := &chunk.WriterOptions{
		ChunksDir:   s.chunksDir,
		Encoding:    chunk.EncodingZstd,
		SyncOnWrite: true,
		ChunkSize:   MaxChunkSize + 1,
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
	// Every read resolves through the safe-path check, absolute forms
	// included: a stored ChunkPath is the absolute path the chunk writer
	// records, and the check validates it in place against the chunks
	// directory, so a traversal or fabricated record cannot answer
	// through the chunk plane while the legitimate file opens as named.
	fullPath, err := s.safeChunkPath(chunkPath)
	if err != nil {
		return nil, err
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

// scanPrefix is indirected so tests can exercise the scan-failure paths
// that ride on scanChunkMetas.
var scanPrefix = (*Store).ScanPrefix

// scanChunkMetas walks the chunk metadata records under one prefix —
// the shared body of the stream- and group-scoped chunk listings. A
// record that fails to decode is invisible to the caller; it is logged
// so chunk corruption stays discoverable instead of silently shrinking
// the chunk list. A scan failure propagates: the read engine must
// answer an error rather than a short page, and a teardown must leave
// the records it can no longer see in place for a retry.
func (s *Store) scanChunkMetas(prefix string) ([]*ChunkMeta, error) {
	var chunks []*ChunkMeta
	if err := scanPrefix(s, prefix, func(key string, value []byte) error {
		if !bytes.HasPrefix([]byte(key), []byte(prefix)) {
			return nil
		}
		var p pb.ChunkMeta
		if err := proto.Unmarshal(value, &p); err != nil {
			logs.Warn("Corrupt chunk metadata record",
				logs.String("key", key), logs.Err(err))
			return nil
		}
		chunks = append(chunks, ProtoToChunkMeta(&p))
		return nil
	}); err != nil {
		return nil, err
	}
	return chunks, nil
}

// ListChunksForStream lists chunk metadata for a specific log stream.
func (s *Store) ListChunksForStream(logGroupName, logStreamName string) ([]*ChunkMeta, error) {
	return s.scanChunkMetas(s.chunkIndexKey(logGroupName, logStreamName, ""))
}

// LateIngestionEvents returns the stream's events that landed below a
// delivery cursor: timestamp strictly before cursorTime, ingestion time
// strictly after ingestionMark. The chunk selection rides the index alone
// (a chunk qualifies when its greatest ingestion time exceeds the mark
// and its least timestamp falls below the cursor), so a stream with no
// late arrivals costs one meta scan and no chunk-file reads. The listing
// runs under the group's read lock — the same lock a PutLogEvents commit
// holds from ingestion-time stamp to index commit — which is what lets a
// caller that stamps its mark at listing time treat everything the
// listing could not see as ingested after that mark.
func (s *Store) LateIngestionEvents(logGroupName, logStreamName string, ingestionMark, cursorTime int64) ([]LogEntry, error) {
	if cursorTime <= 0 {
		return nil, nil
	}
	gLock := s.groupLock(logGroupName)
	gLock.RLock()
	defer gLock.RUnlock()

	chunks, err := s.ListChunksForStream(logGroupName, logStreamName)
	if err != nil {
		return nil, err
	}
	var late []LogEntry
	for _, chunk := range chunks {
		if chunk.MaxIngestionTs <= ingestionMark || chunk.MinTs >= cursorTime {
			continue
		}
		entries, err := s.readChunkFile(chunk.ChunkPath)
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if e.Timestamp < cursorTime && e.IngestionTime > ingestionMark {
				late = append(late, e)
			}
		}
	}
	sort.SliceStable(late, func(i, j int) bool {
		if late[i].Timestamp != late[j].Timestamp {
			return late[i].Timestamp < late[j].Timestamp
		}
		return late[i].Message < late[j].Message
	})
	return late, nil
}

// ListChunksForLogGroup lists chunk metadata for all streams in a log group.
func (s *Store) ListChunksForLogGroup(logGroupName string) ([]*ChunkMeta, error) {
	return s.scanChunkMetas(keyPrefixChunk + escapePath(logGroupName) + ":")
}
