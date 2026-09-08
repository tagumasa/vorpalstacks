package lambda

import (
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// EventSourceMappingCreateInput carries the wire members of a
// CreateEventSourceMapping request. The Has* flags distinguish explicitly
// provided values from omitted members; the raw map and slice members keep
// their wire form so the Core applies the documented parse-and-validate
// order.
type EventSourceMappingCreateInput struct {
	FunctionNameRaw string
	EventSourceArn  string

	HasBatchSize     bool
	BatchSize        int32
	StartingPosition string

	// StartingPositionTimestampRaw is the raw wire value; the Core applies
	// the float64 type assertion exactly as the wire carrier does.
	HasStartingPositionTimestamp bool
	StartingPositionTimestampRaw interface{}

	HasMaximumBatchingWindowInSeconds bool
	MaximumBatchingWindowInSeconds    int32

	HasParallelizationFactor bool
	ParallelizationFactor    int32

	HasEnabled bool
	Enabled    bool

	HasMaximumRecordAgeInSeconds bool
	MaximumRecordAgeInSeconds    int32

	HasMaximumRetryAttempts bool
	MaximumRetryAttempts    int32

	HasTumblingWindowInSeconds bool
	TumblingWindowInSeconds    int32

	HasBisectBatchOnFunctionError bool
	BisectBatchOnFunctionError    bool

	DestinationConfigRaw     map[string]interface{}
	FilterCriteriaRaw        map[string]interface{}
	KMSKeyArn                string
	FunctionResponseTypesRaw []interface{}

	Tags map[string]string
}

// EventSourceMappingUpdateInput carries the wire members of an
// UpdateEventSourceMapping request with the same presence-flag semantics.
type EventSourceMappingUpdateInput struct {
	UUID string

	HasBatchSize bool
	BatchSize    int32

	HasEnabled bool
	Enabled    bool

	HasMaximumBatchingWindowInSeconds bool
	MaximumBatchingWindowInSeconds    int32

	HasParallelizationFactor bool
	ParallelizationFactor    int32

	HasMaximumRecordAgeInSeconds bool
	MaximumRecordAgeInSeconds    int32

	HasMaximumRetryAttempts bool
	MaximumRetryAttempts    int32

	HasTumblingWindowInSeconds bool
	TumblingWindowInSeconds    int32

	HasBisectBatchOnFunctionError bool
	BisectBatchOnFunctionError    bool

	// FunctionNameRaw is the optional FunctionName member of the update
	// request — any reference form Create accepts (name, name:qualifier,
	// full or partial ARN) — which repoints the mapping at another
	// function or qualifier.
	FunctionNameRaw          string
	DestinationConfigRaw     map[string]interface{}
	FilterCriteriaRaw        map[string]interface{}
	KMSKeyArn                string
	FunctionResponseTypesRaw []interface{}
}

// createEventSourceMappingCore creates a mapping between an event source
// and a Lambda function.
func (s *LambdaService) createEventSourceMappingCore(reqCtx *request.RequestContext, in *EventSourceMappingCreateInput) (*lambdastore.EventSourceMapping, error) {
	if in.FunctionNameRaw == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	if err := validateNamespacedFunctionName(in.FunctionNameRaw); err != nil {
		return nil, err
	}
	functionName, embeddedQualifier := resolveNamespacedFunctionRef(in.FunctionNameRaw)

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// The FunctionName reference may address a version or alias via its
	// ":qualifier" suffix; the mapping records the qualified ARN so the
	// poller invokes the addressed qualifier.
	function, _, _, err := s.resolveQualifier(stores.Functions, functionName, embeddedQualifier)
	if err != nil {
		return nil, err
	}
	functionArn := function.FunctionArn
	if embeddedQualifier != "" && embeddedQualifier != "$LATEST" {
		functionArn = function.FunctionArn + ":" + embeddedQualifier
	}

	if in.EventSourceArn == "" {
		return nil, NewInvalidParameter("EventSourceArn", "Event source ARN is required")
	}

	// Validate that the event source ARN refers to a supported service.
	// Only SQS, Kinesis, and DynamoDB streams are polled by the ESM
	// poller; accepting other services would silently create a mapping
	// that never processes events.
	if err := validateEventSourceArn(in.EventSourceArn); err != nil {
		return nil, err
	}

	batchSize := in.BatchSize
	if batchSize <= 0 {
		batchSize = defaultESMBatchSize(in.EventSourceArn)
	}
	if err := validateESMBatchSizeForSource(batchSize, in.EventSourceArn); err != nil {
		return nil, err
	}

	startingPosition := in.StartingPosition
	if startingPosition != "" {
		if err := validateStartingPosition(startingPosition); err != nil {
			return nil, err
		}
	}
	if err := validateStartingPositionForStream(startingPosition, in.EventSourceArn); err != nil {
		return nil, err
	}
	if err := validateStartingPositionTimestamp(startingPosition, in.HasStartingPositionTimestamp); err != nil {
		return nil, err
	}
	var startingPositionTimestamp time.Time
	if in.HasStartingPositionTimestamp {
		// The SDK serialises timestamps as epoch floats; any other shape
		// is malformed on the wire and rejected rather than silently
		// zeroed into an AT_TIMESTAMP mapping.
		ts, ok := in.StartingPositionTimestampRaw.(float64)
		if !ok {
			return nil, NewInvalidParameter("StartingPositionTimestamp", "StartingPositionTimestamp must be a timestamp")
		}
		startingPositionTimestamp = time.Unix(int64(ts), 0).UTC()
	}

	// Validate MaximumBatchingWindowInSeconds if explicitly provided.
	batchingWindow := in.MaximumBatchingWindowInSeconds
	if in.HasMaximumBatchingWindowInSeconds {
		if err := validateESMBatchingWindow(batchingWindow); err != nil {
			return nil, err
		}
	}
	// The pairing rule applies "when you set BatchSize to a value greater
	// than 10" — an explicitly requested size, not the service default.
	if in.HasBatchSize {
		if err := validateESMBatchWindowPair(batchSize, batchingWindow); err != nil {
			return nil, err
		}
	}
	// Validate ParallelizationFactor if explicitly provided.
	parallelizationFactor := int32(1)
	if in.HasParallelizationFactor {
		parallelizationFactor = in.ParallelizationFactor
		if err := validateESMParallelFactor(parallelizationFactor); err != nil {
			return nil, err
		}
	}

	mapping := &lambdastore.EventSourceMapping{
		FunctionArn:                    functionArn,
		FunctionName:                   function.FunctionName,
		EventSourceArn:                 in.EventSourceArn,
		BatchSize:                      batchSize,
		MaximumBatchingWindowInSeconds: batchingWindow,
		ParallelizationFactor:          parallelizationFactor,
		StartingPosition:               startingPosition,
		StartingPositionTimestamp:      startingPositionTimestamp,
		State:                          "Enabled",
	}

	// Fields shared with UpdateEventSourceMapping — accept on Create
	// for API surface parity with AWS.
	if in.HasEnabled {
		if in.Enabled {
			mapping.State = "Enabled"
		} else {
			mapping.State = "Disabled"
		}
	}
	if in.HasMaximumRecordAgeInSeconds {
		if err := validateESMMaxRecordAge(in.MaximumRecordAgeInSeconds); err != nil {
			return nil, err
		}
		mapping.MaximumRecordAgeInSeconds = in.MaximumRecordAgeInSeconds
	} else {
		mapping.MaximumRecordAgeInSeconds = -1
	}
	if in.HasMaximumRetryAttempts {
		if err := validateESMMaxRetry(in.MaximumRetryAttempts); err != nil {
			return nil, err
		}
		mapping.MaximumRetryAttempts = in.MaximumRetryAttempts
	} else {
		mapping.MaximumRetryAttempts = -1
	}
	if in.HasTumblingWindowInSeconds {
		if err := validateESMTumblingWindow(in.TumblingWindowInSeconds); err != nil {
			return nil, err
		}
		mapping.TumblingWindowInSeconds = in.TumblingWindowInSeconds
	}
	if in.HasBisectBatchOnFunctionError {
		mapping.BisectBatchOnFunctionError = in.BisectBatchOnFunctionError
	}
	if in.DestinationConfigRaw != nil {
		destConfig := parseDestinationConfig(in.DestinationConfigRaw)
		if err := validateESMDestinationConfig(destConfig); err != nil {
			return nil, err
		}
		mapping.DestinationConfig = destConfig
	}
	if in.FilterCriteriaRaw != nil {
		mapping.FilterCriteria = parseFilterCriteria(in.FilterCriteriaRaw)
	}
	if in.KMSKeyArn != "" {
		mapping.KMSKeyArn = in.KMSKeyArn
	}
	if in.FunctionResponseTypesRaw != nil {
		parsed, err := parseFunctionResponseTypes(in.FunctionResponseTypesRaw)
		if err != nil {
			return nil, err
		}
		mapping.FunctionResponseTypes = parsed
	}

	created, err := stores.EventSources.Create(mapping)
	if err != nil {
		return nil, mapStoreError(err)
	}

	// "A list of tags to apply to the event source mapping." The tags
	// share the function tag store under a namespaced key; the mapping
	// response shape carries no Tags member, so ListTags is the only
	// reader.
	if len(in.Tags) > 0 {
		if err := stores.Functions.TagStore.Tag(esmTagResourceKey(created.UUID), in.Tags); err != nil {
			return nil, mapStoreError(err)
		}
	}

	return created, nil
}

// deleteEventSourceMappingCore deletes the specified event source mapping
// and drops its tags from the shared tag store.
func (s *LambdaService) deleteEventSourceMappingCore(reqCtx *request.RequestContext, uuid string) (*lambdastore.EventSourceMapping, error) {
	if uuid == "" {
		return nil, NewInvalidParameter("UUID", "UUID is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	mapping, err := stores.EventSources.Get(uuid)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if err := stores.EventSources.Delete(uuid); err != nil {
		return nil, mapStoreError(err)
	}

	// The mapping's tags live in the shared tag store; dropping them here
	// keeps a same-UUID recreation from inheriting stale tags.
	if err := stores.Functions.TagStore.Delete(esmTagResourceKey(uuid)); err != nil {
		return nil, mapStoreError(err)
	}

	return mapping, nil
}

// getEventSourceMappingCore retrieves the specified event source mapping.
func (s *LambdaService) getEventSourceMappingCore(reqCtx *request.RequestContext, uuid string) (*lambdastore.EventSourceMapping, error) {
	if uuid == "" {
		return nil, NewInvalidParameter("UUID", "UUID is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	mapping, err := stores.EventSources.Get(uuid)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return mapping, nil
}

// updateEventSourceMappingCore applies the requested member updates to an
// existing event source mapping. The read-modify-write runs inside the
// store lock, so a poll-cycle state write can never be clobbered by the
// update, and a FunctionName reference resolves with Create's rules
// before the mapping is repointed.
func (s *LambdaService) updateEventSourceMappingCore(reqCtx *request.RequestContext, in *EventSourceMappingUpdateInput) (*lambdastore.EventSourceMapping, error) {
	if in.UUID == "" {
		return nil, NewInvalidParameter("UUID", "UUID is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := stores.EventSources.Get(in.UUID); err != nil {
		return nil, mapStoreError(err)
	}

	// A FunctionName reference, when provided, resolves exactly like
	// Create's so a typo'd name cannot silently break polling and the
	// function-name filter of the list API.
	var newFunctionName, newFunctionArn string
	if in.FunctionNameRaw != "" {
		if err := validateNamespacedFunctionName(in.FunctionNameRaw); err != nil {
			return nil, err
		}
		name, qualifier := resolveNamespacedFunctionRef(in.FunctionNameRaw)
		function, _, _, err := s.resolveQualifier(stores.Functions, name, qualifier)
		if err != nil {
			return nil, err
		}
		newFunctionName = function.FunctionName
		newFunctionArn = function.FunctionArn
		if qualifier != "" && qualifier != "$LATEST" {
			newFunctionArn = function.FunctionArn + ":" + qualifier
		}
	}

	mapping, err := stores.EventSources.UpdateAtomically(in.UUID, func(m *lambdastore.EventSourceMapping) error {
		if in.HasBatchSize {
			if err := validateESMBatchSizeForSource(in.BatchSize, m.EventSourceArn); err != nil {
				return err
			}
			m.BatchSize = in.BatchSize
		}
		if in.HasEnabled {
			if in.Enabled {
				m.State = "Enabled"
			} else {
				m.State = "Disabled"
			}
		}
		if in.HasMaximumBatchingWindowInSeconds {
			if err := validateESMBatchingWindow(in.MaximumBatchingWindowInSeconds); err != nil {
				return err
			}
			m.MaximumBatchingWindowInSeconds = in.MaximumBatchingWindowInSeconds
		}
		if in.HasParallelizationFactor {
			if err := validateESMParallelFactor(in.ParallelizationFactor); err != nil {
				return err
			}
			m.ParallelizationFactor = in.ParallelizationFactor
		}
		// The batch-size/window pairing is validated on the merged mapping when
		// the request sets either member, so a stored explicitly-set BatchSize
		// above 10 also constrains a window-lowering update. The rule binds the
		// setter: requests that set neither member leave defaults alone.
		if in.HasBatchSize || in.HasMaximumBatchingWindowInSeconds {
			if err := validateESMBatchWindowPair(m.BatchSize, m.MaximumBatchingWindowInSeconds); err != nil {
				return err
			}
		}
		if in.HasMaximumRecordAgeInSeconds {
			if err := validateESMMaxRecordAge(in.MaximumRecordAgeInSeconds); err != nil {
				return err
			}
			m.MaximumRecordAgeInSeconds = in.MaximumRecordAgeInSeconds
		}
		if in.HasMaximumRetryAttempts {
			if err := validateESMMaxRetry(in.MaximumRetryAttempts); err != nil {
				return err
			}
			m.MaximumRetryAttempts = in.MaximumRetryAttempts
		}
		if in.HasTumblingWindowInSeconds {
			if err := validateESMTumblingWindow(in.TumblingWindowInSeconds); err != nil {
				return err
			}
			m.TumblingWindowInSeconds = in.TumblingWindowInSeconds
		}
		if in.HasBisectBatchOnFunctionError {
			m.BisectBatchOnFunctionError = in.BisectBatchOnFunctionError
		}
		if newFunctionArn != "" {
			m.FunctionName = newFunctionName
			m.FunctionArn = newFunctionArn
		}
		if in.DestinationConfigRaw != nil {
			destConfig := parseDestinationConfig(in.DestinationConfigRaw)
			if err := validateESMDestinationConfig(destConfig); err != nil {
				return err
			}
			m.DestinationConfig = destConfig
		}
		if in.FilterCriteriaRaw != nil {
			m.FilterCriteria = parseFilterCriteria(in.FilterCriteriaRaw)
		}
		if in.KMSKeyArn != "" {
			m.KMSKeyArn = in.KMSKeyArn
		}
		if in.FunctionResponseTypesRaw != nil {
			parsed, err := parseFunctionResponseTypes(in.FunctionResponseTypesRaw)
			if err != nil {
				return err
			}
			m.FunctionResponseTypes = parsed
		}
		return nil
	})
	if err != nil {
		return nil, mapStoreError(err)
	}

	return mapping, nil
}

// listEventSourceMappingsCore lists the event source mappings for the
// given function or event source filters. A filtered list walks every
// mapping and paginates client-side; an unfiltered list pages in the
// store.
func (s *LambdaService) listEventSourceMappingsCore(stores *lambdaStore, functionNameOrArn, eventSourceArn, marker string, maxItems int) ([]*lambdastore.EventSourceMapping, string, bool, error) {
	hasFilter := functionNameOrArn != "" || eventSourceArn != ""

	if !hasFilter {
		result, err := stores.EventSources.List(storecommon.ListOptions{Marker: marker, MaxItems: maxItems})
		if err != nil {
			return nil, "", false, err
		}
		return result.Items, result.NextMarker, result.IsTruncated, nil
	}

	allMappings, err := stores.EventSources.ListAllMappings()
	if err != nil {
		return nil, "", false, err
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

	return pageResult.Items, pageResult.NextMarker, pageResult.IsTruncated, nil
}
