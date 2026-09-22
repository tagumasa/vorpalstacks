package cloudwatchlogs

import (
	"errors"
	"fmt"
	"os"

	"vorpalstacks/internal/core/logs"
	pb "vorpalstacks/internal/pb/storage/storage_cloudwatchlogs"
	"vorpalstacks/internal/store/aws/common"
)

// CreateLogStream creates a new CloudWatch Logs log stream. The
// existence checks and the write run under the group's lock, the same
// lock the deletes hold: a create interleaving with
// DeleteLogStream/DeleteLogGroup cannot resurrect the stream family, and
// two concurrent creators of one name admit exactly one record.
func (s *Store) CreateLogStream(ls *LogStream) error {
	gLock := s.groupLock(ls.LogGroupName)
	gLock.Lock()
	defer gLock.Unlock()

	if _, err := s.GetLogGroup(ls.LogGroupName); err != nil {
		if !errors.Is(err, ErrLogGroupNotFound) {
			return err
		}
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

// GetLogStream retrieves a CloudWatch Logs log stream by group and stream
// name. A missing record reports the not-found sentinel; any other
// storage failure propagates as a storage error.
func (s *Store) GetLogStream(logGroupName, logStreamName string) (*LogStream, error) {
	key := s.logStreamKey(logGroupName, logStreamName)
	var p pb.LogStream
	if err := s.GetProto(key, &p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrLogStreamNotFound
		}
		return nil, err
	}
	return ProtoToLogStream(&p), nil
}

// DeleteLogStream deletes a CloudWatch Logs log stream — its chunk files,
// index records and StoredBytes share — under the group's lock: an
// in-flight PutLogEvents either commits before the delete (and is then
// removed with the rest) or fails the stream existence check once it
// acquires the lock, so a write can no longer resurrect a deleted
// stream.
func (s *Store) DeleteLogStream(logGroupName, logStreamName string) error {
	gLock := s.groupLock(logGroupName)
	gLock.Lock()
	defer gLock.Unlock()
	return s.deleteLogStreamLocked(logGroupName, logStreamName)
}

func (s *Store) deleteLogStreamLocked(logGroupName, logStreamName string) error {
	if _, err := s.GetLogStream(logGroupName, logStreamName); err != nil {
		return err
	}

	// Without the chunk list the teardown would delete the stream record
	// while its chunks survive as permanent orphans; fail here so the
	// stream stays and a retry reaches them.
	chunks, err := s.ListChunksForStream(logGroupName, logStreamName)
	if err != nil {
		return fmt.Errorf("list chunks: %w", err)
	}
	var bytesRemoved int64
	for _, chunk := range chunks {
		cp, err := s.safeChunkPath(chunk.ChunkPath)
		if err != nil {
			logs.Warn("Path traversal attempt in chunk", logs.String("chunkPath", chunk.ChunkPath), logs.Err(err))
			continue
		}
		if err := os.Remove(cp); err != nil && !os.IsNotExist(err) {
			logs.Error("Failed to remove chunk file", logs.String("chunkPath", chunk.ChunkPath), logs.Err(err))
		}
		deleteKey := s.chunkIndexKey(logGroupName, logStreamName, chunk.ChunkID)
		if err := s.Delete(deleteKey); err != nil {
			// The index record survives, so its bytes stay counted: a
			// later purge or teardown retry re-reaches this chunk (the
			// file removal above already tolerates its absence) and
			// decrements there — counting here would double-decrement on
			// the retry.
			logs.Error("Failed to delete chunk index", logs.String("key", deleteKey), logs.Err(err))
			continue
		}
		bytesRemoved += chunk.ByteSize
	}
	// The stream's bytes leave the group's accounting on the same basis
	// they entered it (the chunk records' ingested message bytes). The
	// decrement precedes the stream-record delete so a failure here
	// leaves the stream deletable again; an orphaned stream whose group
	// record is already gone has no counter left to update.
	if bytesRemoved > 0 {
		if err := s.mutateLogGroupLocked(logGroupName, func(lg *LogGroup) error {
			lg.StoredBytes -= bytesRemoved
			if lg.StoredBytes < 0 {
				lg.StoredBytes = 0
			}
			return nil
		}); err != nil && !errors.Is(err, ErrLogGroupNotFound) {
			return err
		}
	}
	key := s.logStreamKey(logGroupName, logStreamName)
	return s.Delete(key)
}

// ListLogStreams lists CloudWatch Logs log streams for a given log group.
// The marker is a stream NAME — the boundary of a previous page — and is
// translated to its key form so the scan resumes strictly after that
// stream; callers never handle raw store keys.
func (s *Store) ListLogStreams(logGroupName, prefix, marker string, maxItems int) ([]*LogStream, string, error) {
	if maxItems <= 0 {
		maxItems = 50
	}

	streamPrefix := s.logStreamKey(logGroupName, "")
	if prefix != "" {
		// The stream component of the stored key is raw, so the caller's
		// prefix must match it byte-for-byte: stream names legally contain
		// '/' and '\', which an escaped prefix could never match.
		streamPrefix += prefix
	}

	keyMarker := ""
	if marker != "" {
		keyMarker = s.logStreamKey(logGroupName, marker)
	}

	opts := common.ListOptions{
		Prefix:   streamPrefix,
		Marker:   keyMarker,
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
	// The caller-facing marker is the boundary stream NAME, never the raw
	// store key pageIterator reports: a key fed back through the marker
	// translation above would double-prefix and skip almost every later
	// stream, silently truncating marker-driven walks (stream teardown,
	// fetch-all consumers).
	nextMarker := ""
	if result.NextMarker != "" && len(streams) > 0 {
		nextMarker = streams[len(streams)-1].Name
	}
	return streams, nextMarker, nil
}
