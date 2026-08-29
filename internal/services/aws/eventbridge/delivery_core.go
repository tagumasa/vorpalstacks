package eventbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// PutEventsInput carries the parameters for PutEvents. Entries holds the raw
// wire list because every parse or validation failure of a single entry
// becomes a per-entry failure record rather than a request error, so the
// per-entry loop must run at the Core layer in its original order.
type PutEventsInput struct {
	Entries []interface{}
	Region  string
}

// PutEventsResult holds the outcome of putEventsCore.
type PutEventsResult struct {
	FailedEntryCount int32
	Entries          []map[string]interface{}
}

// putEventsCore validates and delivers one or more events to EventBridge,
// matching rules on each entry's event bus. Validations of required fields
// (Source, DetailType, Detail) and the per-entry event bus existence check
// run per entry so callers receive per-entry error codes.
func (s *EventsService) putEventsCore(ctx context.Context, store *eventsstore.EventsStore, input PutEventsInput) (*PutEventsResult, error) {
	resultEntries := make([]map[string]interface{}, 0)
	failedCount := int32(0)

	for _, e := range input.Entries {
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
		if _, err := store.GetEventBus(ctx, eventBusName); err != nil {
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
			Region:       input.Region,
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

		if err := s.deliverEvent(ctx, store, event, eventBusName, input.Region); err != nil {
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

	return &PutEventsResult{
		FailedEntryCount: failedCount,
		Entries:          resultEntries,
	}, nil
}

// deliverEvent archives the event and delivers it to the matching rules on
// the event bus.
func (s *EventsService) deliverEvent(ctx context.Context, store *eventsstore.EventsStore, event *eventsstore.Event, eventBusName, region string) error {
	s.archiveEvent(ctx, store, event, eventBusName)

	return s.deliverEventWithStore(ctx, region, event, eventBusName, store)
}

// deliverEventWithStore delivers an event to every enabled rule on the event
// bus whose event pattern matches, fanning the matched rules' targets out to
// the delivery plane.
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

// archiveEvent stores the event into every enabled archive of the event bus
// whose event pattern matches, updating the archives' counters.
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
