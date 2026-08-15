package cloudwatchlogs

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	pb "vorpalstacks/internal/pb/storage/storage_cloudwatchlogs"
	"vorpalstacks/internal/store/aws/common"
)

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

// DeleteLogGroup deletes a CloudWatch Logs log group together with every
// record and chunk file it owns. The teardown continues past individual
// failures and tolerates files already removed by an earlier partial
// delete, so a retry always makes progress. The group record itself is
// removed only when every step succeeded: deleting it while sub-resources
// survive would orphan them permanently, because a retry stops at the
// missing group before it can reach them.
func (s *Store) DeleteLogGroup(name string) error {
	if _, err := s.GetLogGroup(name); err != nil {
		return err
	}

	var errs []error

	if err := s.deleteAllLogStreams(name); err != nil {
		errs = append(errs, fmt.Errorf("delete log streams: %w", err))
	}

	// Backstop for chunk index records whose stream record is already
	// gone: DeleteLogStream covers only chunks listed under a live
	// stream. The file may already have been removed by an earlier
	// partial delete, which is not an error here.
	chunks := s.ListChunksForLogGroup(name)
	for _, chunk := range chunks {
		cp, err := s.safeChunkPath(chunk.ChunkPath)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if err := os.Remove(cp); err != nil && !os.IsNotExist(err) {
			errs = append(errs, err)
		}
		if err := s.Delete(s.chunkIndexKey(name, chunk.LogStream, chunk.ChunkID)); err != nil {
			errs = append(errs, err)
		}
	}
	metricFilters, _, listErr := s.ListMetricFilters(name, "", "", 1000)
	if listErr != nil {
		// Without the list the filters cannot be deleted; report the
		// failure so the group record survives and a retry can reach them.
		errs = append(errs, fmt.Errorf("list metric filters: %w", listErr))
	}
	for _, mf := range metricFilters {
		if err := s.Delete(s.metricFilterKey(name, mf.Name)); err != nil {
			errs = append(errs, fmt.Errorf("delete metric filter %s: %w", mf.Name, err))
		}
	}
	subFilters, listErr := s.ListSubscriptionFilters(name, "")
	if listErr != nil {
		errs = append(errs, fmt.Errorf("list subscription filters: %w", listErr))
	}
	for _, sf := range subFilters {
		if err := s.Delete(s.subscriptionFilterKey(name, sf.FilterName)); err != nil {
			errs = append(errs, fmt.Errorf("delete subscription filter %s: %w", sf.FilterName, err))
		}
	}
	arn := s.arnBuilder.CloudWatch().LogGroup(name)
	if err := s.tagStore.Delete(arn); err != nil {
		errs = append(errs, fmt.Errorf("delete tags for log group %s: %w", name, err))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	if err := s.Delete(s.logGroupKey(name)); err != nil {
		return err
	}
	return nil
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
