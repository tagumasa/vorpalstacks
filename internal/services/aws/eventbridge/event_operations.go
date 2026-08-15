package eventbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

const maxPutEventsEntries = 10

// PutEvents delivers one or more events to EventBridge.
// Validates required fields (Source, DetailType, Detail) and delivers
// events to matching rules on the specified event bus.
func (s *EventsService) PutEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	entries, ok := req.Parameters["Entries"].([]interface{})
	if !ok {
		entries, ok = req.Parameters["entries"].([]interface{})
	}
	if !ok || len(entries) == 0 {
		return nil, awserrors.NewValidationException("Entries are required")
	}
	if len(entries) > maxPutEventsEntries {
		return nil, awserrors.NewValidationException("Maximum 10 entries allowed per request")
	}

	// EndpointId routes the request through a global endpoint. Global
	// endpoints are out of scope for this edge platform, but we still
	// accept the parameter so SDK clients do not receive an unexpected
	// ValidationException when populating it.
	if endpointID, _ := req.Parameters["EndpointId"].(string); endpointID != "" {
		logs.Debug("PutEvents EndpointId ignored (global endpoints out of scope)",
			logs.String("endpointId", endpointID))
	}

	resultEntries := make([]map[string]interface{}, 0)
	failedCount := int32(0)

	// Acquire the store once before iterating entries. s.store() is
	// cached, but hoisting it keeps the loop body focused on per-entry
	// validation and delivery.
	busStore, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		entryMap, ok := e.(map[string]interface{})
		if !ok {
			failedCount++
			resultEntries = append(resultEntries, map[string]interface{}{
				"ErrorCode":    "ValidationException",
				"ErrorMessage": "Invalid entry format",
			})
			continue
		}

		source, _ := entryMap["Source"].(string)
		detailType, _ := entryMap["DetailType"].(string)
		detailStr, _ := entryMap["Detail"].(string)
		traceHeader, _ := entryMap["TraceHeader"].(string)
		eventBusName := "default"
		if rawBus, busPresent := entryMap["EventBusName"]; busPresent {
			busStr, ok := rawBus.(string)
			if !ok || busStr == "" {
				resultEntries = append(resultEntries, map[string]interface{}{
					"ErrorCode":    "ValidationException",
					"ErrorMessage": "EventBusName must not be empty",
				})
				failedCount++
				continue
			}
			eventBusName = busStr
		}

		// DetailType max 128 chars per AWS EventBridge PutEvents API reference.
		if !validateDetailType(detailType) {
			resultEntries = append(resultEntries, map[string]interface{}{
				"ErrorCode":    "ValidationException",
				"ErrorMessage": "DetailType must be between 1 and 128 characters",
			})
			failedCount++
			continue
		}

		// Source max 256 chars (AWS PutEventsRequestEntry).
		if !validateSource(source) {
			resultEntries = append(resultEntries, map[string]interface{}{
				"ErrorCode":    "ValidationException",
				"ErrorMessage": "Source must be between 1 and 256 characters",
			})
			failedCount++
			continue
		}

		// TraceHeader max 500 chars per Smithy @length trait.
		if !validateTraceHeader(traceHeader) {
			resultEntries = append(resultEntries, map[string]interface{}{
				"ErrorCode":    "ValidationException",
				"ErrorMessage": "TraceHeader must be at most 500 characters",
			})
			failedCount++
			continue
		}

		// Validate the event bus exists before attempting delivery so
		// that callers receive a per-entry error code instead of a
		// silent no-op when targeting a non-existent bus.
		if _, err := busStore.GetEventBus(ctx, eventBusName); err != nil {
			resultEntries = append(resultEntries, map[string]interface{}{
				"ErrorCode":    "ResourceNotFoundException",
				"ErrorMessage": "Event bus '" + eventBusName + "' does not exist",
			})
			failedCount++
			continue
		}

		if source == "" || detailType == "" || detailStr == "" {
			resultEntries = append(resultEntries, map[string]interface{}{
				"ErrorCode":    "ValidationException",
				"ErrorMessage": "Source, DetailType, and Detail are required",
			})
			failedCount++
			continue
		}

		var detail map[string]interface{}
		if err := json.Unmarshal([]byte(detailStr), &detail); err != nil {
			resultEntries = append(resultEntries, map[string]interface{}{
				"ErrorCode":    "ValidationException",
				"ErrorMessage": "Detail must be a valid JSON object",
			})
			failedCount++
			continue
		}

		eventTime := time.Now().UTC()
		if timeVal, ok := entryMap["Time"]; ok {
			switch tv := timeVal.(type) {
			case float64:
				eventTime = time.Unix(int64(tv), 0).UTC()
			case string:
				parsed, err := time.Parse(time.RFC3339, tv)
				if err != nil {
					resultEntries = append(resultEntries, map[string]interface{}{
						"ErrorCode":    "ValidationException",
						"ErrorMessage": "Time must be a valid RFC3339 timestamp",
					})
					failedCount++
					continue
				}
				eventTime = parsed.UTC()
			default:
				resultEntries = append(resultEntries, map[string]interface{}{
					"ErrorCode":    "ValidationException",
					"ErrorMessage": "Time must be a string or numeric timestamp",
				})
				failedCount++
				continue
			}
		}

		event := &eventsstore.Event{
			ID:           generateEventID(),
			Version:      "0",
			DetailType:   detailType,
			Source:       source,
			Account:      s.accountID,
			Time:         eventTime,
			Region:       reqCtx.GetRegion(),
			Detail:       detail,
			EventBusName: eventBusName,
			TraceHeader:  traceHeader,
		}

		if rv, ok := entryMap["Resources"]; ok {
			resources, ok := rv.([]interface{})
			if !ok {
				resultEntries = append(resultEntries, map[string]interface{}{
					"ErrorCode":    "ValidationException",
					"ErrorMessage": "Resources must be an array of strings",
				})
				failedCount++
				continue
			}
			valid := true
			for _, r := range resources {
				rStr, ok := r.(string)
				if !ok {
					valid = false
					break
				}
				event.Resources = append(event.Resources, rStr)
			}
			if !valid {
				resultEntries = append(resultEntries, map[string]interface{}{
					"ErrorCode":    "ValidationException",
					"ErrorMessage": "Resources must be an array of strings",
				})
				failedCount++
				continue
			}
		}

		if err := s.deliverEvent(ctx, reqCtx, event, eventBusName); err != nil {
			resultEntries = append(resultEntries, map[string]interface{}{
				"ErrorCode":    "InternalFailure",
				"ErrorMessage": err.Error(),
			})
			failedCount++
			continue
		}

		resultEntries = append(resultEntries, map[string]interface{}{
			"EventId": event.ID,
		})
	}

	return map[string]interface{}{
		"FailedEntryCount": failedCount,
		"Entries":          resultEntries,
	}, nil
}

func (s *EventsService) deliverEvent(ctx context.Context, reqCtx *request.RequestContext, event *eventsstore.Event, eventBusName string) error {
	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	s.archiveEvent(ctx, store, event, eventBusName)

	return s.deliverEventWithStore(ctx, reqCtx.GetRegion(), event, eventBusName, store)
}

func (s *EventsService) deliverEventWithStore(ctx context.Context, region string, event *eventsstore.Event, eventBusName string, store *eventsstore.EventsStore) error {
	var allRules []*eventsstore.Rule
	nextToken := ""
	for {
		result, err := store.ListRules(ctx, eventBusName, "", 100, nextToken)
		if err != nil {
			return err
		}
		allRules = append(allRules, result.Rules...)
		if result.NextToken == "" {
			break
		}
		nextToken = result.NextToken
	}

	var targetWg sync.WaitGroup
	for _, rule := range allRules {
		if rule.State != eventsstore.RuleStateEnabled && rule.State != eventsstore.RuleStateEnabledWithAllCloudtrailManagementEvents {
			continue
		}

		if rule.EventPattern != "" {
			if !s.matchEventPattern(event, rule.EventPattern) {
				continue
			}
		}

		// Paginate through all targets for the matched rule.
		var allTargets []*eventsstore.Target
		targetsNextToken := ""
		for {
			targetsResult, err := store.ListTargetsByRule(ctx, eventBusName, rule.Name, 100, targetsNextToken)
			if err != nil {
				// Proceeding with a partial target list would silently
				// drop the remaining targets; fail the event so PutEvents
				// reports it as a failed entry.
				return fmt.Errorf("failed to list targets for rule %s: %w", rule.Name, err)
			}
			allTargets = append(allTargets, targetsResult.Targets...)
			if targetsResult.NextToken == "" {
				break
			}
			targetsNextToken = targetsResult.NextToken
		}

		for _, target := range allTargets {
			targetCopy := *target
			payload := s.buildTargetPayload(event, targetCopy)
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				logs.Error("failed to marshal event payload",
					logs.String("eventId", event.ID),
					logs.String("targetId", targetCopy.ID),
					logs.Err(err))
				continue
			}

			if s.bus != nil {
				ebEvt := &eventbus.EventBridgeDeliveryEvent{
					RuleARN:   rule.ARN,
					TargetID:  targetCopy.ID,
					TargetARN: targetCopy.ARN,
					Input:     payloadBytes,
				}
				if targetCopy.DeadLetterConfig != nil {
					ebEvt.DeadLetterConfigArn = targetCopy.DeadLetterConfig.Arn
				}
				if targetCopy.RetryPolicy != nil {
					ebEvt.MaximumRetryAttempts = targetCopy.RetryPolicy.MaximumRetryAttempts
					ebEvt.MaximumEventAgeInSeconds = targetCopy.RetryPolicy.MaximumEventAgeInSeconds
				}
				if targetCopy.SqsParameters != nil {
					ebEvt.SqsMessageGroupId = targetCopy.SqsParameters.MessageGroupId
				}
				if targetCopy.KinesisParameters != nil {
					ebEvt.KinesisPartitionKeyPath = targetCopy.KinesisParameters.PartitionKeyPath
				}
				ebEvt.Region = region
				if err := s.bus.Publish(context.Background(), ebEvt); err != nil {
					logs.Warn("failed to publish event directly", logs.String("targetArn", targetCopy.ARN), logs.Err(err))
				}
			} else {
				targetWg.Add(1)
				select {
				case s.targetSemaphore <- struct{}{}:
					go func(evt *eventsstore.Event, tgt eventsstore.Target) {
						defer func() {
							<-s.targetSemaphore
							targetWg.Done()
							if r := recover(); r != nil {
								logs.Error("eventbridge: panic delivering to target", logs.String("arn", tgt.ARN), logs.Any("panic", r))
							}
						}()
						s.deliverToTarget(ctx, region, evt, tgt)
					}(event, targetCopy)
				case <-ctx.Done():
					targetWg.Done()
					goto done
				}
			}
		}
	}

done:
	targetWg.Wait()
	return nil
}

func (s *EventsService) buildTargetPayload(event *eventsstore.Event, target eventsstore.Target) map[string]interface{} {
	payload := map[string]interface{}{
		"version":     event.Version,
		"id":          event.ID,
		"detail-type": event.DetailType,
		"source":      event.Source,
		"account":     event.Account,
		"time":        event.Time,
		"region":      event.Region,
		"resources":   event.Resources,
		"detail":      event.Detail,
	}

	if target.Input != "" {
		var inputPayload map[string]interface{}
		if err := json.Unmarshal([]byte(target.Input), &inputPayload); err == nil {
			payload = inputPayload
		}
	} else if target.InputPath != "" {
		extracted := s.extractInputPath(payload, target.InputPath)
		if extracted != nil {
			payload = extracted
		}
	} else if target.InputTransformer != nil {
		transformed := s.applyInputTransformer(payload, target.InputTransformer)
		if transformed != nil {
			payload = transformed
		}
	}

	return payload
}

func (s *EventsService) archiveEvent(ctx context.Context, store *eventsstore.EventsStore, event *eventsstore.Event, eventBusName string) {
	eventBus, err := store.GetEventBus(ctx, eventBusName)
	if err != nil {
		return
	}

	archives, err := store.ListArchivesForEventBus(ctx, eventBusName)
	if err != nil {
		return
	}

	for _, archive := range archives {
		if archive.State != eventsstore.ArchiveStateEnabled {
			continue
		}

		if archive.EventPattern != "" {
			if !s.matchEventPattern(event, archive.EventPattern) {
				continue
			}
		}

		eventMap := map[string]interface{}{
			"version":     event.Version,
			"id":          event.ID,
			"detail-type": event.DetailType,
			"source":      event.Source,
			"account":     event.Account,
			"time":        event.Time.Format(time.RFC3339),
			"region":      event.Region,
			"resources":   event.Resources,
			"detail":      event.Detail,
		}

		archivedEvent := &eventsstore.ArchivedEvent{
			ID:          event.ID,
			EventBusARN: eventBus.ARN,
			Event:       eventMap,
			Timestamp:   event.Time,
		}

		if err := store.StoreArchiveEvent(ctx, archive.Name, archivedEvent); err != nil {
			logs.Warn("failed to archive event",
				logs.String("eventId", event.ID),
				logs.String("archiveName", archive.Name),
				logs.Err(err))
			continue
		}

		eventSize := int64(0)
		if eventBytes, err := json.Marshal(eventMap); err == nil {
			eventSize = int64(len(eventBytes))
		}
		if err := store.IncrementArchiveCounters(ctx, archive.Name, eventSize); err != nil {
			logs.Warn("failed to update archive counters",
				logs.String("archiveName", archive.Name),
				logs.Err(err))
		}
	}
}

func (s *EventsService) TestEventPattern(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	patternStr := request.GetStringParam(req.Parameters, "EventPattern")
	eventStr := request.GetStringParam(req.Parameters, "Event")

	if patternStr == "" {
		return nil, awserrors.NewValidationException("Parameter EventPattern is required")
	}
	if !validateEventPatternLength(patternStr) {
		return nil, awserrors.NewValidationException("EventPattern must be at most 4096 characters")
	}
	if eventStr == "" {
		return nil, awserrors.NewValidationException("Parameter Event is required")
	}

	var patternMap, eventMap map[string]interface{}
	if err := json.Unmarshal([]byte(patternStr), &patternMap); err != nil {
		return nil, awserrors.NewInvalidEventPatternException(fmt.Sprintf("EventPattern is not valid JSON: %s", err))
	}
	if err := json.Unmarshal([]byte(eventStr), &eventMap); err != nil {
		return nil, awserrors.NewValidationException(fmt.Sprintf("Event is not valid JSON: %s", err))
	}

	result := true
	for key, patternValue := range patternMap {
		eventValue, exists := eventMap[key]
		if !exists {
			if isExistsFalsePattern(patternValue) {
				continue
			}
			result = false
			break
		}
		if !s.matchValue(eventValue, patternValue) {
			result = false
			break
		}
	}

	return map[string]interface{}{
		"Result": result,
	}, nil
}

func generateEventID() string {
	return uuid.New().String()
}
