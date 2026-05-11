package kinesis

import (
	"context"
	"encoding/base64"
	"strconv"
	"time"

	"vorpalstacks/internal/common/request"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// ListShards lists the shards in a Kinesis stream.
func (s *KinesisService) ListShards(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, streamName, err := s.resolveStreamNameOptional(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	if streamName != "" {
		if _, err := store.GetStream(streamName); err != nil {
			return nil, s.mapStoreError(err)
		}
	}

	var filter *kinesisstore.ShardFilter
	shardFilterMap := request.GetMapParam(req.Parameters, "ShardFilter")
	if shardFilterMap == nil {
		shardFilterMap = request.GetMapParam(req.Parameters, "shardFilter")
	}
	if shardFilterMap != nil {
		filter = &kinesisstore.ShardFilter{}
		if ft, ok := shardFilterMap["Type"].(string); ok {
			filter.Type = ft
		} else if ft, ok := shardFilterMap["type"].(string); ok {
			filter.Type = ft
		}
		if shardID, ok := shardFilterMap["ShardId"].(string); ok {
			filter.ShardID = shardID
		} else if shardID, ok := shardFilterMap["shardId"].(string); ok {
			filter.ShardID = shardID
		}
		if ts, ok := shardFilterMap["Timestamp"].(string); ok {
			if unixTs, err := strconv.ParseInt(ts, 10, 64); err == nil {
				t := time.Unix(unixTs, 0).UTC()
				filter.Timestamp = &t
			}
		} else if ts, ok := shardFilterMap["timestamp"].(string); ok {
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

	resp := map[string]interface{}{
		"Shards": formatShards(shards),
	}

	if len(shards) == limit {
		lastShard := shards[len(shards)-1]
		resp["NextToken"] = base64.StdEncoding.EncodeToString([]byte(lastShard.ShardID))
	}

	return resp, nil
}

// SplitShard splits a shard in a Kinesis stream.
func (s *KinesisService) SplitShard(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, streamName, err := s.resolveStreamName(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	shardID := request.GetParamLowerFirst(req.Parameters, "ShardToSplit")
	newHashKey := request.GetParamLowerFirst(req.Parameters, "NewStartingHashKey")
	if shardID == "" {
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
	store, streamName, err := s.resolveStreamName(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	shardID1 := request.GetParamLowerFirst(req.Parameters, "ShardToMerge")
	shardID2 := request.GetParamLowerFirst(req.Parameters, "AdjacentShardToMerge")
	if shardID1 == "" || shardID2 == "" {
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
	store, streamName, err := s.resolveStreamName(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	targetCount := int32(request.GetIntParam(req.Parameters, "TargetShardCount"))
	if targetCount <= 0 {
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
