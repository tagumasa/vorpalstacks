package kinesis

import (
	"context"
	"encoding/base64"
	"strconv"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// ListShards lists the shards in a Kinesis stream.
func (s *KinesisService) ListShards(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if streamARN != "" {
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName != "" {
		if _, err := store.GetStream(streamName); err != nil {
			return nil, s.mapStoreError(err)
		}
	}

	var filter *kinesisstore.ShardFilter
	if ft := request.GetParamLowerFirst(req.Parameters, "ShardFilter.Type"); ft != "" {
		filter = &kinesisstore.ShardFilter{
			Type: ft,
		}
		if shardID := request.GetParamLowerFirst(req.Parameters, "ShardFilter.ShardId"); shardID != "" {
			filter.ShardID = shardID
		}
		if ts := request.GetParamLowerFirst(req.Parameters, "ShardFilter.Timestamp"); ts != "" {
			if unixTs, err := strconv.ParseInt(ts, 10, 64); err == nil {
				t := time.Unix(unixTs, 0).UTC()
				filter.Timestamp = &t
			}
		}
	}

	limit := int(request.GetIntParam(req.Parameters, "Limit"))
	if limit <= 0 {
		limit = 1000
	}

	var nextToken string
	if encoded := request.GetStringParam(req.Parameters, "NextToken"); encoded != "" {
		if decoded, err := base64.StdEncoding.DecodeString(encoded); err == nil {
			nextToken = string(decoded)
		}
	}

	shards, err := store.ListShards(streamName, filter, nextToken, limit)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	response := map[string]interface{}{
		"Shards": s.formatShards(shards),
	}

	if len(shards) == limit {
		lastShard := shards[len(shards)-1]
		response["NextToken"] = base64.StdEncoding.EncodeToString([]byte(lastShard.ShardID))
	}

	return response, nil
}

// SplitShard splits a shard in a Kinesis stream.
func (s *KinesisService) SplitShard(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	shardID := request.GetParamLowerFirst(req.Parameters, "ShardToSplit")
	newHashKey := request.GetParamLowerFirst(req.Parameters, "NewStartingHashKey")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if streamARN != "" {
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" || shardID == "" {
		return nil, ErrInvalidArgument
	}

	if err := store.SplitShard(streamName, shardID, newHashKey); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"StreamARN": store.BuildStreamARN(streamName),
	}, nil
}

// MergeShards merges two adjacent shards in a Kinesis stream.
func (s *KinesisService) MergeShards(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	shardID1 := request.GetParamLowerFirst(req.Parameters, "ShardToMerge")
	shardID2 := request.GetParamLowerFirst(req.Parameters, "AdjacentShardToMerge")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if streamARN != "" {
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" || shardID1 == "" || shardID2 == "" {
		return nil, ErrInvalidArgument
	}

	if err := store.MergeShards(streamName, shardID1, shardID2); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"StreamARN": store.BuildStreamARN(streamName),
	}, nil
}

// UpdateShardCount updates the shard count of a Kinesis stream.
func (s *KinesisService) UpdateShardCount(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	targetCount := int32(request.GetIntParam(req.Parameters, "TargetShardCount"))

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if streamARN != "" {
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" || targetCount <= 0 {
		return nil, ErrInvalidArgument
	}

	if err := store.UpdateShardCount(streamName, targetCount); err != nil {
		return nil, s.mapStoreError(err)
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"StreamName":        streamName,
		"CurrentShardCount": stream.ShardCount,
		"TargetShardCount":  targetCount,
		"StreamARN":         stream.StreamARN,
	}, nil
}
