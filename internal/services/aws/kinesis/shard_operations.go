package kinesis

import (
	"context"
	"encoding/base64"
	"sort"
	"strconv"
	"time"

	"vorpalstacks/internal/common/request"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// ListShards lists the shards in a Kinesis stream.
func (s *KinesisService) ListShards(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
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

	// The wire member is MaxResults; the value and its presence travel to
	// the Core so an explicit out-of-window value is rejected instead of
	// being folded into the default page.
	maxResults := 0
	hasMaxResults := false
	if _, ok := req.Parameters["MaxResults"]; ok {
		maxResults = request.GetIntParam(req.Parameters, "MaxResults")
		hasMaxResults = true
	}

	var nextToken string
	if encoded := request.GetStringParam(req.Parameters, "NextToken"); encoded != "" {
		if decoded, err := base64.StdEncoding.DecodeString(encoded); err == nil {
			nextToken = string(decoded)
		}
	}

	result, err := s.listShardsCore(store, ListShardsInput{
		StreamName:              request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:               request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		StreamCreationTimestamp: request.GetParamLowerFirst(req.Parameters, "StreamCreationTimestamp"),
		ShardFilter:             filter,
		MaxResults:              maxResults,
		HasMaxResults:           hasMaxResults,
		NextToken:               nextToken,
	})
	if err != nil {
		return nil, err
	}

	shards := result.Shards

	// Apply ShardOrder (ASCENDING is default, DESCENDING reverses)
	shardOrder := request.GetParamLowerFirst(req.Parameters, "ShardOrder")
	if shardOrder == "DESCENDING" {
		sort.Slice(shards, func(i, j int) bool {
			return shards[i].ShardID > shards[j].ShardID
		})
	}

	resp := map[string]interface{}{
		"Shards": formatShards(shards),
	}

	if len(shards) == result.EffectiveLimit {
		lastShard := shards[len(shards)-1]
		resp["NextToken"] = base64.StdEncoding.EncodeToString([]byte(lastShard.ShardID))
	}

	return resp, nil
}

// SplitShard splits a shard in a Kinesis stream.
func (s *KinesisService) SplitShard(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	streamARN, err := s.splitShardCore(store, SplitShardInput{
		StreamName:         request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:          request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		ShardToSplit:       request.GetParamLowerFirst(req.Parameters, "ShardToSplit"),
		NewStartingHashKey: request.GetParamLowerFirst(req.Parameters, "NewStartingHashKey"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"StreamARN": streamARN,
	}, nil
}

// MergeShards merges two adjacent shards in a Kinesis stream.
func (s *KinesisService) MergeShards(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	streamARN, err := s.mergeShardsCore(store, MergeShardsInput{
		StreamName:           request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:            request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		ShardToMerge:         request.GetParamLowerFirst(req.Parameters, "ShardToMerge"),
		AdjacentShardToMerge: request.GetParamLowerFirst(req.Parameters, "AdjacentShardToMerge"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"StreamARN": streamARN,
	}, nil
}

// UpdateShardCount updates the shard count of a Kinesis stream.
func (s *KinesisService) UpdateShardCount(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.updateShardCountCore(store, UpdateShardCountInput{
		StreamName:       request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:        request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		TargetShardCount: int32(request.GetIntParam(req.Parameters, "TargetShardCount")),
		ScalingType:      request.GetParamLowerFirst(req.Parameters, "ScalingType"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"StreamName":        result.StreamName,
		"CurrentShardCount": result.CurrentShardCount,
		"TargetShardCount":  result.TargetShardCount,
		"StreamARN":         result.StreamARN,
	}, nil
}
