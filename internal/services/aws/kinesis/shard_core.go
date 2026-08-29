package kinesis

import (
	"strconv"

	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// The MaxResults input window and the effective page follow two different
// documented rules: the ListShardsInputLimit range trait accepts 1-10000 on
// the wire, while the MaxResults member documentation fixes the default
// page at 1000 and returns at most 1000 results per call.
const (
	DefaultListShardsResults = 1000
	MaxListShardsResults     = 1000
)

// ListShardsResult carries a page of shards plus the effective page size
// the Core applied, so the transport can emit the resumption token.
type ListShardsResult struct {
	Shards         []*kinesisstore.Shard
	EffectiveLimit int
}

// ListShardsInput is the transport-agnostic input for ListShards. The stream
// identification is optional; NextToken is the decoded resumption shard ID.
type ListShardsInput struct {
	StreamName              string
	StreamARN               string
	StreamCreationTimestamp string
	ShardFilter             *kinesisstore.ShardFilter
	MaxResults              int
	HasMaxResults           bool
	NextToken               string
}

// SplitShardInput is the transport-agnostic input for SplitShard.
type SplitShardInput struct {
	StreamName         string
	StreamARN          string
	ShardToSplit       string
	NewStartingHashKey string
}

// MergeShardsInput is the transport-agnostic input for MergeShards.
type MergeShardsInput struct {
	StreamName           string
	StreamARN            string
	ShardToMerge         string
	AdjacentShardToMerge string
}

// UpdateShardCountInput is the transport-agnostic input for UpdateShardCount.
type UpdateShardCountInput struct {
	StreamName       string
	StreamARN        string
	TargetShardCount int32
	ScalingType      string
}

// UpdateShardCountResult carries the reshaping receipt.
type UpdateShardCountResult struct {
	StreamName        string
	CurrentShardCount int32
	TargetShardCount  int32
	StreamARN         string
}

// listShardsCore validates the MaxResults window, resolves the optional
// stream, verifies the creation timestamp when provided, and lists the
// shards through the store filter.
func (s *KinesisService) listShardsCore(store *kinesisstore.KinesisStore, input ListShardsInput) (ListShardsResult, error) {
	// The range trait window is enforced on the provided value, then the
	// effective page applies the documented cap.
	limit := input.MaxResults
	if input.HasMaxResults {
		if !validateListShardsLimit(limit) {
			return ListShardsResult{}, ErrInvalidArgument
		}
		if limit > MaxListShardsResults {
			limit = MaxListShardsResults
		}
	} else {
		limit = DefaultListShardsResults
	}

	streamName, err := s.resolveStreamNameOptionalCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return ListShardsResult{}, err
	}

	if streamName != "" {
		stream, err := store.GetStream(streamName)
		if err != nil {
			return ListShardsResult{}, s.mapStoreError(err)
		}
		// Verify StreamCreationTimestamp matches when provided
		// (used to disambiguate deleted+recreated streams)
		if input.StreamCreationTimestamp != "" {
			if unixTs, err := strconv.ParseFloat(input.StreamCreationTimestamp, 64); err == nil {
				// Compare at second precision: stream.CreatedAt has nanosecond
				// resolution but the client timestamp is epoch seconds.
				if stream.CreatedAt.Unix() != int64(unixTs) {
					return ListShardsResult{}, s.mapStoreError(kinesisstore.ErrStreamNotFound)
				}
			}
		}
	}

	shards, err := store.ListShards(streamName, input.ShardFilter, input.NextToken, limit)
	if err != nil {
		return ListShardsResult{}, s.mapStoreError(err)
	}

	return ListShardsResult{Shards: shards, EffectiveLimit: limit}, nil
}

// splitShardCore splits a shard at the given starting hash key.
func (s *KinesisService) splitShardCore(store *kinesisstore.KinesisStore, input SplitShardInput) (string, error) {
	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return "", err
	}

	if input.ShardToSplit == "" {
		return "", ErrInvalidArgument
	}

	if err := store.SplitShard(streamName, input.ShardToSplit, input.NewStartingHashKey); err != nil {
		return "", s.mapStoreError(err)
	}

	return store.BuildStreamARN(streamName), nil
}

// mergeShardsCore merges two adjacent shards.
func (s *KinesisService) mergeShardsCore(store *kinesisstore.KinesisStore, input MergeShardsInput) (string, error) {
	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return "", err
	}

	if input.ShardToMerge == "" || input.AdjacentShardToMerge == "" {
		return "", ErrInvalidArgument
	}

	if err := store.MergeShards(streamName, input.ShardToMerge, input.AdjacentShardToMerge); err != nil {
		return "", s.mapStoreError(err)
	}

	return store.BuildStreamARN(streamName), nil
}

// updateShardCountCore reshapes the stream to the target shard count.
func (s *KinesisService) updateShardCountCore(store *kinesisstore.KinesisStore, input UpdateShardCountInput) (UpdateShardCountResult, error) {
	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return UpdateShardCountResult{}, err
	}

	if !validateShardCount(input.TargetShardCount) {
		return UpdateShardCountResult{}, ErrInvalidArgument
	}

	if input.ScalingType != "" && input.ScalingType != "UNIFORM_SCALING" {
		return UpdateShardCountResult{}, ErrInvalidArgument
	}

	if err := store.UpdateShardCount(streamName, input.TargetShardCount); err != nil {
		return UpdateShardCountResult{}, s.mapStoreError(err)
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return UpdateShardCountResult{}, s.mapStoreError(err)
	}

	return UpdateShardCountResult{
		StreamName:        streamName,
		CurrentShardCount: stream.ShardCount,
		TargetShardCount:  input.TargetShardCount,
		StreamARN:         stream.StreamARN,
	}, nil
}
