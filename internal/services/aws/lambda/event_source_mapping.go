package lambda

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/store/aws/common"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateEventSourceMapping creates a mapping between an event source and a Lambda function.
func (s *LambdaService) CreateEventSourceMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, err := s.validateAndGetFunction(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	eventSourceArn := request.GetStringParam(req.Parameters, "EventSourceArn")
	if eventSourceArn == "" {
		return nil, NewInvalidParameter("EventSourceArn", "Event source ARN is required")
	}

	// Validate that the event source ARN refers to a supported service.
	// Only SQS, Kinesis, and DynamoDB streams are polled by the ESM
	// poller; accepting other services would silently create a mapping
	// that never processes events.
	if err := validateEventSourceArn(eventSourceArn); err != nil {
		return nil, err
	}

	batchSize := int32(request.GetIntParam(req.Parameters, "BatchSize"))
	if batchSize <= 0 {
		batchSize = 100
	}

	startingPosition := request.GetStringParam(req.Parameters, "StartingPosition")
	if startingPosition != "" {
		if err := validateStartingPosition(startingPosition); err != nil {
			return nil, err
		}
	}
	if err := validateStartingPositionForStream(startingPosition, eventSourceArn); err != nil {
		return nil, err
	}
	_, hasTimestamp := req.Parameters["StartingPositionTimestamp"]
	if err := validateStartingPositionTimestamp(startingPosition, hasTimestamp); err != nil {
		return nil, err
	}

	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:                    function.FunctionArn,
		FunctionName:                   function.FunctionName,
		EventSourceArn:                 eventSourceArn,
		BatchSize:                      batchSize,
		MaximumBatchingWindowInSeconds: int32(request.GetIntParam(req.Parameters, "MaximumBatchingWindowInSeconds")),
		ParallelizationFactor:          int32(request.GetIntParam(req.Parameters, "ParallelizationFactor")),
		StartingPosition:               startingPosition,
		State:                          "Enabled",
	}

	// Fields shared with UpdateEventSourceMapping — accept on Create
	// for API surface parity with AWS.
	if _, ok := req.Parameters["Enabled"]; ok {
		if request.GetBoolParam(req.Parameters, "Enabled") {
			mapping.State = "Enabled"
		} else {
			mapping.State = "Disabled"
		}
	}
	if _, ok := req.Parameters["MaximumRecordAgeInSeconds"]; ok {
		mapping.MaximumRecordAgeInSeconds = int32(request.GetIntParam(req.Parameters, "MaximumRecordAgeInSeconds"))
	} else {
		mapping.MaximumRecordAgeInSeconds = -1
	}
	if _, ok := req.Parameters["MaximumRetryAttempts"]; ok {
		mapping.MaximumRetryAttempts = int32(request.GetIntParam(req.Parameters, "MaximumRetryAttempts"))
	} else {
		mapping.MaximumRetryAttempts = -1
	}
	if _, ok := req.Parameters["TumblingWindowInSeconds"]; ok {
		mapping.TumblingWindowInSeconds = int32(request.GetIntParam(req.Parameters, "TumblingWindowInSeconds"))
	}
	if _, ok := req.Parameters["BisectBatchOnFunctionError"]; ok {
		mapping.BisectBatchOnFunctionError = request.GetBoolParam(req.Parameters, "BisectBatchOnFunctionError")
	}
	if destMap := request.GetMapParam(req.Parameters, "DestinationConfig"); destMap != nil {
		mapping.DestinationConfig = parseDestinationConfig(destMap)
	}
	if filterMap := request.GetMapParam(req.Parameters, "FilterCriteria"); filterMap != nil {
		mapping.FilterCriteria = parseFilterCriteria(filterMap)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := store.EventSources.Create(mapping)
	if err != nil {
		return nil, err
	}

	return s.toEventSourceMapping(created), nil
}

// DeleteEventSourceMapping deletes the specified event source mapping.
func (s *LambdaService) DeleteEventSourceMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	uuid := request.GetStringParam(req.Parameters, "UUID")
	if uuid == "" {
		return nil, NewInvalidParameter("UUID", "UUID is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	mapping, err := store.EventSources.Get(uuid)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if err := store.EventSources.Delete(uuid); err != nil {
		return nil, err
	}

	return s.toEventSourceMapping(mapping), nil
}

// GetEventSourceMapping retrieves the specified event source mapping.
func (s *LambdaService) GetEventSourceMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	uuid := request.GetStringParam(req.Parameters, "UUID")
	if uuid == "" {
		return nil, NewInvalidParameter("UUID", "UUID is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	mapping, err := store.EventSources.Get(uuid)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return s.toEventSourceMapping(mapping), nil
}

// UpdateEventSourceMapping updates the specified event source mapping.
func (s *LambdaService) UpdateEventSourceMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	uuid := request.GetStringParam(req.Parameters, "UUID")
	if uuid == "" {
		return nil, NewInvalidParameter("UUID", "UUID is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	mapping, err := store.EventSources.Get(uuid)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if batchSize := request.GetIntParam(req.Parameters, "BatchSize"); batchSize > 0 {
		mapping.BatchSize = int32(batchSize)
	}
	if _, ok := req.Parameters["Enabled"]; ok {
		if enabled := request.GetBoolParam(req.Parameters, "Enabled"); enabled {
			mapping.State = "Enabled"
		} else {
			mapping.State = "Disabled"
		}
	}
	if _, ok := req.Parameters["MaximumBatchingWindowInSeconds"]; ok {
		mapping.MaximumBatchingWindowInSeconds = int32(request.GetIntParam(req.Parameters, "MaximumBatchingWindowInSeconds"))
	}
	if parallelFactor := request.GetIntParam(req.Parameters, "ParallelizationFactor"); parallelFactor > 0 {
		mapping.ParallelizationFactor = int32(parallelFactor)
	}
	if _, ok := req.Parameters["MaximumRecordAgeInSeconds"]; ok {
		mapping.MaximumRecordAgeInSeconds = int32(request.GetIntParam(req.Parameters, "MaximumRecordAgeInSeconds"))
	}
	if _, ok := req.Parameters["MaximumRetryAttempts"]; ok {
		mapping.MaximumRetryAttempts = int32(request.GetIntParam(req.Parameters, "MaximumRetryAttempts"))
	}
	if _, ok := req.Parameters["TumblingWindowInSeconds"]; ok {
		mapping.TumblingWindowInSeconds = int32(request.GetIntParam(req.Parameters, "TumblingWindowInSeconds"))
	}
	if _, ok := req.Parameters["BisectBatchOnFunctionError"]; ok {
		mapping.BisectBatchOnFunctionError = request.GetBoolParam(req.Parameters, "BisectBatchOnFunctionError")
	}
	if functionArn := request.GetStringParam(req.Parameters, "FunctionArn"); functionArn != "" {
		mapping.FunctionArn = functionArn
	}
	if destMap := request.GetMapParam(req.Parameters, "DestinationConfig"); destMap != nil {
		mapping.DestinationConfig = parseDestinationConfig(destMap)
	}
	if filterMap := request.GetMapParam(req.Parameters, "FilterCriteria"); filterMap != nil {
		mapping.FilterCriteria = parseFilterCriteria(filterMap)
	}

	if err := store.EventSources.Update(mapping); err != nil {
		return nil, err
	}

	return s.toEventSourceMapping(mapping), nil
}

// ListEventSourceMappings lists all event source mappings for the specified function or event source.
func (s *LambdaService) ListEventSourceMappings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionNameOrArn := request.GetStringParam(req.Parameters, "FunctionName")
	eventSourceArn := request.GetStringParam(req.Parameters, "EventSourceArn")
	marker := request.GetStringParam(req.Parameters, "Marker")

	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	if maxItems <= 0 {
		maxItems = 50
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	hasFilter := functionNameOrArn != "" || eventSourceArn != ""

	var allMappings []*lambdastore.EventSourceMapping
	if hasFilter {
		allMappings, err = store.EventSources.ListAllMappings()
		if err != nil {
			return nil, err
		}
	} else {
		result, err := store.EventSources.List(common.ListOptions{Marker: marker, MaxItems: maxItems})
		if err != nil {
			return nil, err
		}
		allMappings = result.Items

		mappings := make([]interface{}, 0, len(allMappings))
		for _, m := range allMappings {
			mappings = append(mappings, s.toEventSourceMapping(m))
		}

		response := map[string]interface{}{
			"EventSourceMappings": mappings,
		}
		if result.IsTruncated {
			response["NextMarker"] = result.NextMarker
		}
		return response, nil
	}

	filtered := make([]*lambdastore.EventSourceMapping, 0)
	for _, m := range allMappings {
		if functionNameOrArn != "" && m.FunctionName != functionNameOrArn && m.FunctionArn != functionNameOrArn {
			continue
		}
		if eventSourceArn != "" && m.EventSourceArn != eventSourceArn {
			continue
		}
		filtered = append(filtered, m)
	}

	pageResult := pagination.PaginateSlice(filtered, marker, maxItems, func(m *lambdastore.EventSourceMapping) string {
		return m.UUID
	})

	mappings := make([]interface{}, 0, len(pageResult.Items))
	for _, m := range pageResult.Items {
		mappings = append(mappings, s.toEventSourceMapping(m))
	}

	response := map[string]interface{}{
		"EventSourceMappings": mappings,
	}

	if pageResult.IsTruncated {
		response["NextMarker"] = pageResult.NextMarker
	}

	return response, nil
}

func (s *LambdaService) toEventSourceMapping(m *lambdastore.EventSourceMapping) map[string]interface{} {
	result := map[string]interface{}{
		"UUID":                           m.UUID,
		"FunctionArn":                    m.FunctionArn,
		"FunctionName":                   m.FunctionName,
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
