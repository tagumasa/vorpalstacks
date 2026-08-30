package lambda

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateEventSourceMapping creates a mapping between an event source and a Lambda function.
func (s *LambdaService) CreateEventSourceMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := &EventSourceMappingCreateInput{
		FunctionNameRaw:  request.GetStringParam(req.Parameters, "FunctionName"),
		EventSourceArn:   request.GetStringParam(req.Parameters, "EventSourceArn"),
		StartingPosition: request.GetStringParam(req.Parameters, "StartingPosition"),
		KMSKeyArn:        request.GetStringParam(req.Parameters, "KMSKeyArn"),
	}
	if _, ok := req.Parameters["BatchSize"]; ok {
		in.HasBatchSize = true
		in.BatchSize = int32(request.GetIntParam(req.Parameters, "BatchSize"))
	}
	if raw, ok := req.Parameters["StartingPositionTimestamp"]; ok {
		in.HasStartingPositionTimestamp = true
		in.StartingPositionTimestampRaw = raw
	}
	if _, ok := req.Parameters["MaximumBatchingWindowInSeconds"]; ok {
		in.HasMaximumBatchingWindowInSeconds = true
		in.MaximumBatchingWindowInSeconds = int32(request.GetIntParam(req.Parameters, "MaximumBatchingWindowInSeconds"))
	}
	if _, ok := req.Parameters["ParallelizationFactor"]; ok {
		in.HasParallelizationFactor = true
		in.ParallelizationFactor = int32(request.GetIntParam(req.Parameters, "ParallelizationFactor"))
	}
	if _, ok := req.Parameters["Enabled"]; ok {
		in.HasEnabled = true
		in.Enabled = request.GetBoolParam(req.Parameters, "Enabled")
	}
	if _, ok := req.Parameters["MaximumRecordAgeInSeconds"]; ok {
		in.HasMaximumRecordAgeInSeconds = true
		in.MaximumRecordAgeInSeconds = int32(request.GetIntParam(req.Parameters, "MaximumRecordAgeInSeconds"))
	}
	if _, ok := req.Parameters["MaximumRetryAttempts"]; ok {
		in.HasMaximumRetryAttempts = true
		in.MaximumRetryAttempts = int32(request.GetIntParam(req.Parameters, "MaximumRetryAttempts"))
	}
	if _, ok := req.Parameters["TumblingWindowInSeconds"]; ok {
		in.HasTumblingWindowInSeconds = true
		in.TumblingWindowInSeconds = int32(request.GetIntParam(req.Parameters, "TumblingWindowInSeconds"))
	}
	if _, ok := req.Parameters["BisectBatchOnFunctionError"]; ok {
		in.HasBisectBatchOnFunctionError = true
		in.BisectBatchOnFunctionError = request.GetBoolParam(req.Parameters, "BisectBatchOnFunctionError")
	}
	in.DestinationConfigRaw = request.GetMapParam(req.Parameters, "DestinationConfig")
	in.FilterCriteriaRaw = request.GetMapParam(req.Parameters, "FilterCriteria")
	if frts, ok := req.Parameters["FunctionResponseTypes"].([]interface{}); ok {
		in.FunctionResponseTypesRaw = frts
	}
	if tm, ok := req.Parameters["Tags"].(map[string]interface{}); ok {
		in.Tags = tagutil.ToMap(tagutil.MapInterfaceToTags(tm))
	}

	created, err := s.createEventSourceMappingCore(reqCtx, in)
	if err != nil {
		return nil, err
	}

	return s.toEventSourceMapping(created), nil
}

// DeleteEventSourceMapping deletes the specified event source mapping.
func (s *LambdaService) DeleteEventSourceMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	mapping, err := s.deleteEventSourceMappingCore(reqCtx, request.GetStringParam(req.Parameters, "UUID"))
	if err != nil {
		return nil, err
	}

	return s.toEventSourceMapping(mapping), nil
}

// GetEventSourceMapping retrieves the specified event source mapping.
func (s *LambdaService) GetEventSourceMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	mapping, err := s.getEventSourceMappingCore(reqCtx, request.GetStringParam(req.Parameters, "UUID"))
	if err != nil {
		return nil, err
	}

	return s.toEventSourceMapping(mapping), nil
}

// UpdateEventSourceMapping updates the specified event source mapping.
func (s *LambdaService) UpdateEventSourceMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := &EventSourceMappingUpdateInput{
		UUID:        request.GetStringParam(req.Parameters, "UUID"),
		FunctionArn: request.GetStringParam(req.Parameters, "FunctionArn"),
		KMSKeyArn:   request.GetStringParam(req.Parameters, "KMSKeyArn"),
	}
	if _, ok := req.Parameters["BatchSize"]; ok {
		in.HasBatchSize = true
		in.BatchSize = int32(request.GetIntParam(req.Parameters, "BatchSize"))
	}
	if _, ok := req.Parameters["Enabled"]; ok {
		in.HasEnabled = true
		in.Enabled = request.GetBoolParam(req.Parameters, "Enabled")
	}
	if _, ok := req.Parameters["MaximumBatchingWindowInSeconds"]; ok {
		in.HasMaximumBatchingWindowInSeconds = true
		in.MaximumBatchingWindowInSeconds = int32(request.GetIntParam(req.Parameters, "MaximumBatchingWindowInSeconds"))
	}
	if _, ok := req.Parameters["ParallelizationFactor"]; ok {
		in.HasParallelizationFactor = true
		in.ParallelizationFactor = int32(request.GetIntParam(req.Parameters, "ParallelizationFactor"))
	}
	if _, ok := req.Parameters["MaximumRecordAgeInSeconds"]; ok {
		in.HasMaximumRecordAgeInSeconds = true
		in.MaximumRecordAgeInSeconds = int32(request.GetIntParam(req.Parameters, "MaximumRecordAgeInSeconds"))
	}
	if _, ok := req.Parameters["MaximumRetryAttempts"]; ok {
		in.HasMaximumRetryAttempts = true
		in.MaximumRetryAttempts = int32(request.GetIntParam(req.Parameters, "MaximumRetryAttempts"))
	}
	if _, ok := req.Parameters["TumblingWindowInSeconds"]; ok {
		in.HasTumblingWindowInSeconds = true
		in.TumblingWindowInSeconds = int32(request.GetIntParam(req.Parameters, "TumblingWindowInSeconds"))
	}
	if _, ok := req.Parameters["BisectBatchOnFunctionError"]; ok {
		in.HasBisectBatchOnFunctionError = true
		in.BisectBatchOnFunctionError = request.GetBoolParam(req.Parameters, "BisectBatchOnFunctionError")
	}
	in.DestinationConfigRaw = request.GetMapParam(req.Parameters, "DestinationConfig")
	in.FilterCriteriaRaw = request.GetMapParam(req.Parameters, "FilterCriteria")
	if frts, ok := req.Parameters["FunctionResponseTypes"].([]interface{}); ok {
		in.FunctionResponseTypesRaw = frts
	}

	mapping, err := s.updateEventSourceMappingCore(reqCtx, in)
	if err != nil {
		return nil, err
	}

	return s.toEventSourceMapping(mapping), nil
}

// ListEventSourceMappings lists all event source mappings for the specified function or event source.
func (s *LambdaService) ListEventSourceMappings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionNameOrArn := request.GetStringParam(req.Parameters, "FunctionName")
	eventSourceArn := request.GetStringParam(req.Parameters, "EventSourceArn")
	marker := request.GetStringParam(req.Parameters, "Marker")

	maxItems := validateMaxItemsCapped(request.GetIntParam(req.Parameters, "MaxItems"), maxEventSourceMappingListItemsCap)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	mappings, nextMarker, isTruncated, err := s.listEventSourceMappingsCore(store, functionNameOrArn, eventSourceArn, marker, maxItems)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(mappings))
	for _, m := range mappings {
		items = append(items, s.toEventSourceMapping(m))
	}

	resp := map[string]interface{}{
		"EventSourceMappings": items,
	}
	if isTruncated {
		resp["NextMarker"] = nextMarker
	}

	return resp, nil
}

// eventSourceMappingArn derives the ESM ARN from the mapping's function
// ARN region and account: arn:aws:lambda:<region>:<account>:event-source-mapping:<uuid>.
func eventSourceMappingArn(m *lambdastore.EventSourceMapping) string {
	_, _, region, accountID, _ := arnutil.SplitARN(m.FunctionArn)
	return fmt.Sprintf("arn:aws:lambda:%s:%s:event-source-mapping:%s", region, accountID, m.UUID)
}

func (s *LambdaService) toEventSourceMapping(m *lambdastore.EventSourceMapping) map[string]interface{} {
	result := map[string]interface{}{
		"UUID":                           m.UUID,
		"EventSourceMappingArn":          eventSourceMappingArn(m),
		"FunctionArn":                    m.FunctionArn,
		"EventSourceArn":                 m.EventSourceArn,
		"BatchSize":                      m.BatchSize,
		"MaximumBatchingWindowInSeconds": m.MaximumBatchingWindowInSeconds,
		"ParallelizationFactor":          m.ParallelizationFactor,
		"LastModified":                   timeutils.FormatEpochSeconds(m.LastModified),
		"LastProcessingResult":           m.LastProcessingResult,
		"State":                          m.State,
		"StateTransitionReason":          m.StateTransitionReason,
		"StartingPosition":               m.StartingPosition,
	}
	if !m.StartingPositionTimestamp.IsZero() {
		result["StartingPositionTimestamp"] = float64(m.StartingPositionTimestamp.Unix())
	}
	if m.KMSKeyArn != "" {
		result["KMSKeyArn"] = m.KMSKeyArn
	}
	if len(m.FunctionResponseTypes) > 0 {
		result["FunctionResponseTypes"] = m.FunctionResponseTypes
	}

	if m.DestinationConfig != nil {
		result["DestinationConfig"] = toDestinationConfig(m.DestinationConfig)
	}
	if m.FilterCriteria != nil {
		result["FilterCriteria"] = toFilterCriteria(m.FilterCriteria)
	}
	// AWS always includes MaximumRecordAgeInSeconds and MaximumRetryAttempts
	// in the ESM response (they default to -1 when not explicitly configured).
	result["MaximumRecordAgeInSeconds"] = m.MaximumRecordAgeInSeconds
	result["MaximumRetryAttempts"] = m.MaximumRetryAttempts
	if m.TumblingWindowInSeconds > 0 {
		result["TumblingWindowInSeconds"] = m.TumblingWindowInSeconds
	}
	if m.BisectBatchOnFunctionError {
		result["BisectBatchOnFunctionError"] = m.BisectBatchOnFunctionError
	}

	return result
}
