package kinesis

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// ListShards lists the shards in a Kinesis stream.
func (s *KinesisService) ListShards(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// The wire member is MaxResults; the value and its presence travel to
	// the Core so an explicit out-of-window value is rejected instead of
	// being folded into the default page.
	maxResults, hasMaxResults, err := strictIntParam(req.Parameters, "MaxResults")
	if err != nil {
		return nil, err
	}

	// The NextToken member travels raw: the envelope — including the
	// ExpiredNextTokenException a malformed one answers — parses in the
	// Core, the definition both planes share.
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
		for _, tsKey := range []string{"Timestamp", "timestamp"} {
			if raw, ok := shardFilterMap[tsKey]; ok && raw != nil {
				// The member's wire value is read through the shared
				// timestamp normaliser: the SDKs serialise it as an
				// epoch-second JSON number, which a string-only read
				// drops and silently degrades the filter.
				s, err := request.NormalizeTimestampValue(raw)
				if err != nil {
					return nil, ErrInvalidArgument
				}
				t, err := parseTimestampMember(s)
				if err != nil {
					return nil, err
				}
				filter.Timestamp = &t
				break
			}
		}
	}

	// The SDKs serialise StreamCreationTimestamp as an epoch-second JSON
	// number; the strict reader carries both that form and the documented
	// string notations, so the Core's generation check actually operates.
	streamCreationTimestamp, _, err := strictTimestampParam(req.Parameters, "StreamCreationTimestamp")
	if err != nil {
		return nil, err
	}

	result, err := s.listShardsCore(reqCtx, ListShardsInput{
		StreamName:              request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:               request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		StreamCreationTimestamp: streamCreationTimestamp,
		ShardFilter:             filter,
		MaxResults:              maxResults,
		HasMaxResults:           hasMaxResults,
		NextToken:               request.GetParamLowerFirst(req.Parameters, "NextToken"),
		ExclusiveStartShardId:   request.GetParamLowerFirst(req.Parameters, "ExclusiveStartShardId"),
	})
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"Shards": formatShards(result.Shards),
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}
	return resp, nil
}

// SplitShard splits a shard in a Kinesis stream.
func (s *KinesisService) SplitShard(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamARN, err := s.splitShardCore(reqCtx, SplitShardInput{
		StreamName:         request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:          request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		ShardToSplit:       request.GetParamLowerFirst(req.Parameters, "ShardToSplit"),
		NewStartingHashKey: request.GetParamLowerFirst(req.Parameters, "NewStartingHashKey"),
	})
	if err != nil {
		return nil, err
	}

	// The model types this operation's output as Unit — no members.
	_ = streamARN
	return response.EmptyResponse(), nil
}

// MergeShards merges two adjacent shards in a Kinesis stream.
func (s *KinesisService) MergeShards(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamARN, err := s.mergeShardsCore(reqCtx, MergeShardsInput{
		StreamName:           request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:            request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		ShardToMerge:         request.GetParamLowerFirst(req.Parameters, "ShardToMerge"),
		AdjacentShardToMerge: request.GetParamLowerFirst(req.Parameters, "AdjacentShardToMerge"),
	})
	if err != nil {
		return nil, err
	}

	// The model types this operation's output as Unit — no members.
	_ = streamARN
	return response.EmptyResponse(), nil
}

// UpdateShardCount updates the shard count of a Kinesis stream.
func (s *KinesisService) UpdateShardCount(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	targetShardCount, _, err := strictIntParam(req.Parameters, "TargetShardCount")
	if err != nil {
		return nil, err
	}

	result, err := s.updateShardCountCore(reqCtx, UpdateShardCountInput{
		StreamName:       request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:        request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		TargetShardCount: int32(targetShardCount),
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
