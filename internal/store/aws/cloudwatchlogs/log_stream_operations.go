package cloudwatchlogs

import (
	"fmt"
	"os"

	"vorpalstacks/internal/core/logs"
	pb "vorpalstacks/internal/pb/storage/storage_cloudwatchlogs"
	"vorpalstacks/internal/store/aws/common"
)

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
		if err := os.Remove(cp); err != nil && !os.IsNotExist(err) {
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
