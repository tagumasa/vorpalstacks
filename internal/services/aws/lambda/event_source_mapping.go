package lambda

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateEventSourceMapping creates a mapping between an event source and a Lambda function.
func (s *LambdaService) CreateEventSourceMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionNameRaw := request.GetStringParam(req.Parameters, "FunctionName")
	if functionNameRaw == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	functionName, embeddedQualifier := resolveFunctionRef(functionNameRaw)
	if err := validateFunctionName(functionName); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// The FunctionName reference may address a version or alias via its
	// ":qualifier" suffix; the mapping records the qualified ARN so the
	// poller invokes the addressed qualifier.
	function, _, _, err := s.resolveQualifier(store.Functions, functionName, embeddedQualifier)
	if err != nil {
		return nil, err
	}
	functionArn := function.FunctionArn
	if embeddedQualifier != "" && embeddedQualifier != "$LATEST" {
		functionArn = function.FunctionArn + ":" + embeddedQualifier
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
		batchSize = defaultESMBatchSize(eventSourceArn)
	}
	if err := validateESMBatchSizeForSource(batchSize, eventSourceArn); err != nil {
		return nil, err
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
	var startingPositionTimestamp time.Time
	if ts, ok := req.Parameters["StartingPositionTimestamp"].(float64); ok {
		startingPositionTimestamp = time.Unix(int64(ts), 0).UTC()
	}

	// Validate MaximumBatchingWindowInSeconds if explicitly provided.
	batchingWindow := int32(request.GetIntParam(req.Parameters, "MaximumBatchingWindowInSeconds"))
	if _, ok := req.Parameters["MaximumBatchingWindowInSeconds"]; ok {
		if err := validateESMBatchingWindow(batchingWindow); err != nil {
			return nil, err
		}
	}
	// The pairing rule applies "when you set BatchSize to a value greater
	// than 10" — an explicitly requested size, not the service default.
	if _, ok := req.Parameters["BatchSize"]; ok {
		if err := validateESMBatchWindowPair(batchSize, batchingWindow); err != nil {
			return nil, err
		}
	}
	// Validate ParallelizationFactor if explicitly provided.
	parallelizationFactor := int32(1)
	if _, ok := req.Parameters["ParallelizationFactor"]; ok {
		parallelizationFactor = int32(request.GetIntParam(req.Parameters, "ParallelizationFactor"))
		if err := validateESMParallelFactor(parallelizationFactor); err != nil {
			return nil, err
		}
	}

	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:                    functionArn,
		FunctionName:                   function.FunctionName,
		EventSourceArn:                 eventSourceArn,
		BatchSize:                      batchSize,
		MaximumBatchingWindowInSeconds: batchingWindow,
		ParallelizationFactor:          parallelizationFactor,
		StartingPosition:               startingPosition,
		StartingPositionTimestamp:      startingPositionTimestamp,
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
		val := int32(request.GetIntParam(req.Parameters, "MaximumRecordAgeInSeconds"))
		if err := validateESMMaxRecordAge(val); err != nil {
			return nil, err
		}
		mapping.MaximumRecordAgeInSeconds = val
	} else {
		mapping.MaximumRecordAgeInSeconds = -1
	}
	if _, ok := req.Parameters["MaximumRetryAttempts"]; ok {
		val := int32(request.GetIntParam(req.Parameters, "MaximumRetryAttempts"))
		if err := validateESMMaxRetry(val); err != nil {
			return nil, err
		}
		mapping.MaximumRetryAttempts = val
	} else {
		mapping.MaximumRetryAttempts = -1
	}
	if _, ok := req.Parameters["TumblingWindowInSeconds"]; ok {
		val := int32(request.GetIntParam(req.Parameters, "TumblingWindowInSeconds"))
		if err := validateESMTumblingWindow(val); err != nil {
			return nil, err
		}
		mapping.TumblingWindowInSeconds = val
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
	if kmsKeyArn := request.GetStringParam(req.Parameters, "KMSKeyArn"); kmsKeyArn != "" {
		mapping.KMSKeyArn = kmsKeyArn
	}
	if frts, ok := req.Parameters["FunctionResponseTypes"].([]interface{}); ok {
		parsed, err := parseFunctionResponseTypes(frts)
		if err != nil {
			return nil, err
		}
		mapping.FunctionResponseTypes = parsed
	}

	created, err := store.EventSources.Create(mapping)
	if err != nil {
		return nil, mapStoreError(err)
	}

	// "A list of tags to apply to the event source mapping." The tags
	// share the function tag store under a namespaced key; the mapping
	// response shape carries no Tags member, so ListTags is the only
	// reader.
	if tm, ok := req.Parameters["Tags"].(map[string]interface{}); ok {
		if tags := tagutil.ToMap(tagutil.MapInterfaceToTags(tm)); len(tags) > 0 {
			if terr := store.Functions.TagStore.Tag(esmTagResourceKey(created.UUID), tags); terr != nil {
				return nil, mapStoreError(terr)
			}
		}
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
		return nil, mapStoreError(err)
	}

	// The mapping's tags live in the shared tag store; dropping them here
	// keeps a same-UUID recreation from inheriting stale tags.
	if err := store.Functions.TagStore.Delete(esmTagResourceKey(uuid)); err != nil {
		return nil, mapStoreError(err)
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
		if err := validateESMBatchSizeForSource(int32(batchSize), mapping.EventSourceArn); err != nil {
			return nil, err
		}
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
		val := int32(request.GetIntParam(req.Parameters, "MaximumBatchingWindowInSeconds"))
		if err := validateESMBatchingWindow(val); err != nil {
			return nil, err
		}
		mapping.MaximumBatchingWindowInSeconds = val
	}
	if _, ok := req.Parameters["ParallelizationFactor"]; ok {
		val := int32(request.GetIntParam(req.Parameters, "ParallelizationFactor"))
		if err := validateESMParallelFactor(val); err != nil {
			return nil, err
		}
		mapping.ParallelizationFactor = val
	}
	// The batch-size/window pairing is validated on the merged mapping when
	// the request sets either member, so a stored explicitly-set BatchSize
	// above 10 also constrains a window-lowering update. The rule binds the
	// setter: requests that set neither member leave defaults alone.
	_, batchSet := req.Parameters["BatchSize"]
	_, windowSet := req.Parameters["MaximumBatchingWindowInSeconds"]
	if batchSet || windowSet {
		if err := validateESMBatchWindowPair(mapping.BatchSize, mapping.MaximumBatchingWindowInSeconds); err != nil {
			return nil, err
		}
	}
	if _, ok := req.Parameters["MaximumRecordAgeInSeconds"]; ok {
		val := int32(request.GetIntParam(req.Parameters, "MaximumRecordAgeInSeconds"))
		if err := validateESMMaxRecordAge(val); err != nil {
			return nil, err
		}
		mapping.MaximumRecordAgeInSeconds = val
	}
	if _, ok := req.Parameters["MaximumRetryAttempts"]; ok {
		val := int32(request.GetIntParam(req.Parameters, "MaximumRetryAttempts"))
		if err := validateESMMaxRetry(val); err != nil {
			return nil, err
		}
		mapping.MaximumRetryAttempts = val
	}
	if _, ok := req.Parameters["TumblingWindowInSeconds"]; ok {
		val := int32(request.GetIntParam(req.Parameters, "TumblingWindowInSeconds"))
		if err := validateESMTumblingWindow(val); err != nil {
			return nil, err
		}
		mapping.TumblingWindowInSeconds = val
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
	if kmsKeyArn := request.GetStringParam(req.Parameters, "KMSKeyArn"); kmsKeyArn != "" {
		mapping.KMSKeyArn = kmsKeyArn
	}
	if frts, ok := req.Parameters["FunctionResponseTypes"].([]interface{}); ok {
		parsed, err := parseFunctionResponseTypes(frts)
		if err != nil {
			return nil, err
		}
		mapping.FunctionResponseTypes = parsed
	}

	if err := store.EventSources.Update(mapping); err != nil {
		return nil, mapStoreError(err)
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
