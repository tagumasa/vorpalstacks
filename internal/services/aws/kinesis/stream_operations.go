package kinesis

import (
	"context"
	"errors"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// CreateStream creates a new Kinesis stream.
func (s *KinesisService) CreateStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	shardCount, hasShardCount, err := strictIntParam(req.Parameters, "ShardCount")
	if err != nil {
		return nil, err
	}
	maxRecordSizeInKiB, hasMaxRecordSize, err := strictIntParam(req.Parameters, "MaxRecordSizeInKiB")
	if err != nil {
		return nil, err
	}
	warmThroughputMiBps, hasWarmThroughput, err := strictIntParam(req.Parameters, "WarmThroughputMiBps")
	if err != nil {
		return nil, err
	}

	streamMode, hasStreamModeDetails := parseStreamModeDetails(req.Parameters)

	stream, err := s.createStreamCore(reqCtx, CreateStreamInput{
		StreamName:             request.GetParamLowerFirst(req.Parameters, "StreamName"),
		ShardCount:             int32(shardCount),
		HasShardCount:          hasShardCount,
		StreamMode:             streamMode,
		HasStreamModeDetails:   hasStreamModeDetails,
		MaxRecordSizeInKiB:     int32(maxRecordSizeInKiB),
		HasMaxRecordSizeInKiB:  hasMaxRecordSize,
		WarmThroughputMiBps:    int32(warmThroughputMiBps),
		HasWarmThroughputMiBps: hasWarmThroughput,
		Tags:                   tags.ParseTags(req.Parameters, "Tags"),
	})
	if err != nil {
		return nil, err
	}

	// The model types this operation's output as Unit — no members.
	_ = stream
	return response.EmptyResponse(), nil
}

// DeleteStream deletes a Kinesis stream.
func (s *KinesisService) DeleteStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.deleteStreamCore(reqCtx, DeleteStreamInput{
		StreamName: request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:  request.GetParamLowerFirst(req.Parameters, "StreamARN"),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeStream returns detailed information about a Kinesis stream.
func (s *KinesisService) DescribeStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit, hasLimit, err := strictIntParam(req.Parameters, "Limit")
	if err != nil {
		return nil, err
	}
	result, err := s.describeStreamCore(reqCtx, DescribeStreamInput{
		StreamName:            request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:             request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		Limit:                 limit,
		HasLimit:              hasLimit,
		ExclusiveStartShardId: request.GetParamLowerFirst(req.Parameters, "ExclusiveStartShardId"),
	})
	if err != nil {
		return nil, err
	}

	stream := result.Stream
	shards := result.Shards

	description := map[string]interface{}{
		"StreamName":              stream.StreamName,
		"StreamARN":               stream.StreamARN,
		"StreamStatus":            stream.StreamStatus,
		"StreamModeDetails":       formatStreamModeDetails(stream.StreamModeDetails),
		"Shards":                  formatShards(shards),
		"HasMoreShards":           result.HasMoreShards,
		"RetentionPeriodHours":    stream.RetentionPeriodHours,
		"StreamCreationTimestamp": formatEpochSeconds(stream.CreatedAt),
		"EnhancedMonitoring":      formatEnhancedMonitoring(stream.EnhancedMonitoring),
		"EncryptionType":          resolveEncryptionType(stream),
	}
	// KeyId is optional: it travels only with a configured KMS key, and an
	// unencrypted stream omits the member rather than answering an empty
	// identifier.
	if stream.KeyID != "" {
		description["KeyId"] = stream.KeyID
	}
	return map[string]interface{}{
		"StreamDescription": description,
	}, nil
}

// DescribeStreamSummary returns summary information about a Kinesis stream.
func (s *KinesisService) DescribeStreamSummary(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.describeStreamSummaryCore(reqCtx, DescribeStreamSummaryInput{
		StreamName: request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:  request.GetParamLowerFirst(req.Parameters, "StreamARN"),
	})
	if err != nil {
		return nil, err
	}

	stream := result.Stream

	summary := map[string]interface{}{
		"StreamName":              stream.StreamName,
		"StreamARN":               stream.StreamARN,
		"StreamStatus":            stream.StreamStatus,
		"StreamModeDetails":       formatStreamModeDetails(stream.StreamModeDetails),
		"ConsumerCount":           stream.ConsumerCount,
		"OpenShardCount":          stream.ShardCount,
		"RetentionPeriodHours":    stream.RetentionPeriodHours,
		"StreamCreationTimestamp": formatEpochSeconds(stream.CreatedAt),
		"EnhancedMonitoring":      formatEnhancedMonitoring(stream.EnhancedMonitoring),
		"EncryptionType":          resolveEncryptionType(stream),
		"MaxRecordSizeInKiB":      stream.MaxRecordSizeInKiB,
		// The Channel family is not implemented on this platform: the
		// count is reported for shape coherence with the model.
		"ChannelCount": 0,
	}
	// KeyId is optional: it travels only with a configured KMS key.
	if stream.KeyID != "" {
		summary["KeyId"] = stream.KeyID
	}
	// WarmThroughput is optional in the model: the summary reports the
	// configured figure. Zero is the documented release floor ("To release
	// excess capacity, call the API again and set the warm throughput to
	// the same or a lower value"; the update response echoes the accepted
	// target, zero included), so the summary reports the released state as
	// the absence of a configured figure — the stored int32 carries no
	// presence bit, and a released target and a never-configured stream
	// read identically here.
	if stream.WarmThroughputMiBps > 0 {
		summary["WarmThroughput"] = map[string]interface{}{
			"CurrentMiBps": stream.WarmThroughputMiBps,
			"TargetMiBps":  stream.WarmThroughputMiBps,
		}
	}
	return map[string]interface{}{
		"StreamDescriptionSummary": summary,
	}, nil
}

// ListStreams lists the Kinesis streams.
func (s *KinesisService) ListStreams(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit, hasLimit, err := strictIntParam(req.Parameters, "Limit")
	if err != nil {
		return nil, err
	}

	result, err := s.listStreamsCore(reqCtx, ListStreamsInput{
		ExclusiveStartStreamName: request.GetStringParam(req.Parameters, "ExclusiveStartStreamName"),
		Limit:                    limit,
		HasLimit:                 hasLimit,
		NextToken:                request.GetStringParam(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	streamNames := make([]string, 0, len(result.Streams))
	streamSummaries := make([]map[string]interface{}, 0, len(result.Streams))
	for _, stream := range result.Streams {
		streamNames = append(streamNames, stream.StreamName)
		streamSummaries = append(streamSummaries, map[string]interface{}{
			"StreamName":              stream.StreamName,
			"StreamARN":               stream.StreamARN,
			"StreamStatus":            stream.StreamStatus,
			"StreamModeDetails":       formatStreamModeDetails(stream.StreamModeDetails),
			"StreamCreationTimestamp": formatEpochSeconds(stream.CreatedAt),
		})
	}

	resp := map[string]interface{}{
		"StreamNames":     streamNames,
		"StreamSummaries": streamSummaries,
		"HasMoreStreams":  result.IsTruncated,
	}
	if result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}
	return resp, nil
}

// UpdateStreamMode updates the stream mode of a Kinesis stream.
func (s *KinesisService) UpdateStreamMode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	warmThroughputMiBps, hasWarmThroughput, err := strictIntParam(req.Parameters, "WarmThroughputMiBps")
	if err != nil {
		return nil, err
	}

	streamMode, _ := parseStreamModeDetails(req.Parameters)

	if err := s.updateStreamModeCore(reqCtx, UpdateStreamModeInput{
		StreamARN:           request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		StreamMode:          streamMode,
		WarmThroughputMiBps: int32(warmThroughputMiBps),
		HasWarmThroughput:   hasWarmThroughput,
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func formatShards(shards []*kinesisstore.Shard) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(shards))
	for _, shard := range shards {
		// A shard record without a range or hash-key submessage reads as
		// empty-string bounds — the emitted shape stays well-formed
		// instead of panicking on the absent record, the same
		// absent-submessage reading the store's own guards take.
		startHash, endHash := "", ""
		if shard.HashKeyRange != nil {
			startHash, endHash = shard.HashKeyRange.StartingHashKey, shard.HashKeyRange.EndingHashKey
		}
		startSeq, endSeq := "", ""
		if shard.SequenceNumberRange != nil {
			startSeq, endSeq = shard.SequenceNumberRange.StartingSequenceNumber, shard.SequenceNumberRange.EndingSequenceNumber
		}
		m := map[string]interface{}{
			"ShardId": shard.ShardID,
			"HashKeyRange": map[string]interface{}{
				"StartingHashKey": startHash,
				"EndingHashKey":   endHash,
			},
			"SequenceNumberRange": map[string]interface{}{
				"StartingSequenceNumber": startSeq,
			},
		}
		if shard.ParentShardID != "" {
			m["ParentShardId"] = shard.ParentShardID
		}
		if shard.AdjacentParentShardID != "" {
			m["AdjacentParentShardId"] = shard.AdjacentParentShardID
		}
		if endSeq != "" {
			m["SequenceNumberRange"].(map[string]interface{})["EndingSequenceNumber"] = endSeq
		}
		result = append(result, m)
	}
	return result
}

func formatEnhancedMonitoring(em []kinesisstore.EnhancedMonitoring) []map[string]interface{} {
	result := make([]map[string]interface{}, len(em))
	for i, m := range em {
		metrics := m.ShardLevelMetrics
		if metrics == nil {
			metrics = []string{}
		}
		result[i] = map[string]interface{}{
			"ShardLevelMetrics": metrics,
		}
	}
	return result
}

func formatStreamModeDetails(smd *kinesisstore.StreamModeDetails) map[string]interface{} {
	if smd != nil {
		return map[string]interface{}{
			"StreamMode": string(smd.StreamMode),
		}
	}
	return map[string]interface{}{
		"StreamMode": "PROVISIONED",
	}
}

func mergeMetrics(current, added []string) []string {
	seen := make(map[string]bool, len(current)+len(added))
	for _, m := range current {
		seen[m] = true
	}
	result := make([]string, 0, len(current)+len(added))
	result = append(result, current...)
	for _, m := range added {
		if !seen[m] {
			result = append(result, m)
			seen[m] = true
		}
	}
	return result
}

func subtractMetrics(current, removed []string) []string {
	removeSet := make(map[string]bool, len(removed))
	for _, m := range removed {
		removeSet[m] = true
	}
	result := make([]string, 0, len(current))
	for _, m := range current {
		if !removeSet[m] {
			result = append(result, m)
		}
	}
	return result
}

func (s *KinesisService) mapStoreError(err error) error {
	if err == nil {
		return nil
	}
	// An already wire-shaped error — a Core guard refusal surfaced through
	// a store call — passes through unchanged.
	var awsErr *awserrors.AWSError
	if errors.As(err, &awsErr) {
		return err
	}
	mapped := awserrors.MapStoreError(err, storeErrorMappings)
	if mapped != err {
		return mapped
	}
	// An unmapped store error is an infrastructure failure — a wrapped
	// storage error, a proto marshal failure — never a client fault.
	return ErrInternalFailure
}
