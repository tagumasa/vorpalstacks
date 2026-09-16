package eventbridge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
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

// parseEntryTimestamp reads a PutEvents entry Time member: an RFC 3339
// string or a numeric epoch, anything else being an InvalidArgument. Both
// ingress planes — the API entry loop here and the internal bus handler —
// share it, so the member carries the same semantics whichever plane a
// publisher uses.
func parseEntryTimestamp(v interface{}) (time.Time, error) {
	switch tv := v.(type) {
	case float64:
		return time.Unix(int64(tv), 0).UTC(), nil
	case string:
		parsed, err := time.Parse(time.RFC3339, tv)
		if err != nil {
			return time.Time{}, errors.New("Time must be a valid RFC3339 timestamp")
		}
		return parsed.UTC(), nil
	default:
		return time.Time{}, errors.New("Time must be a string or numeric timestamp")
	}
}

// putEventsCore validates and delivers one or more events to EventBridge,
// matching rules on each entry's event bus. Validations of required fields
// (Source, DetailType, Detail) and the member bounds run per entry so
// callers receive per-entry error codes; the entry-count, request-size and
// none-complete-entry rules fail the entire request.
func (s *EventsService) putEventsCore(ctx context.Context, store *eventsstore.EventsStore, input PutEventsInput) (*PutEventsResult, error) {
	// Request-level bounds. The entry count is the model @length(1,10)
	// trait on the Entries list; the request size is the user guide's
	// summed-entry ceiling ("the total request size ... must be less than
	// 1 MB ... This limit applies to the request as a whole, not to
	// individual entries"); and the none-complete-entry rule is the model
	// member documentation: "If you submit a request in which none of the
	// entries have each of these properties, EventBridge fails the entire
	// request."
	if len(input.Entries) == 0 {
		return nil, awserrors.NewValidationException("Entries are required")
	}
	if len(input.Entries) > eventsstore.MaxPutEventsEntries {
		return nil, awserrors.NewValidationException(
			fmt.Sprintf("Maximum %d entries allowed per request", eventsstore.MaxPutEventsEntries))
	}

	totalSize := 0
	anyEntryComplete := false
	for _, e := range input.Entries {
		entryMap, ok := e.(map[string]interface{})
		if !ok {
			continue
		}
		if _, specified := entryMap["Time"]; specified {
			totalSize += eventsstore.PutEventsTimeFieldSizeBytes
		}
		source, _ := entryMap["Source"].(string)
		detailType, _ := entryMap["DetailType"].(string)
		detailStr, _ := entryMap["Detail"].(string)
		totalSize += len(source) + len(detailType) + len(detailStr)
		if resources, ok := entryMap["Resources"].([]interface{}); ok {
			for _, r := range resources {
				if rStr, ok := r.(string); ok {
					totalSize += len(rStr)
				}
			}
		}
		if source != "" && detailType != "" && detailStr != "" {
			anyEntryComplete = true
		}
	}
	if !anyEntryComplete {
		return nil, awserrors.NewValidationException(
			"Source, DetailType, and Detail are required for EventBridge to send an event to an event bus; no entry in the request includes all three")
	}
	if totalSize >= eventsstore.MaxPutEventsRequestSizeBytes {
		return nil, awserrors.NewValidationException(
			fmt.Sprintf("The total entry size of the request must be less than %d bytes", eventsstore.MaxPutEventsRequestSizeBytes))
	}

	resultEntries := make([]map[string]interface{}, 0)
	failedCount := int32(0)

	// The per-entry failure codes stay inside the documented ErrorCode
	// vocabulary (the model's PutEventsResultEntry.ErrorCode
	// documentation): InvalidArgument for an invalid parameter value,
	// MalformedDetail for invalid Detail JSON, InternalFailure as the
	// retryable delivery failure. ValidationException and
	// ResourceNotFoundException are not per-entry codes and do not appear
	// in this loop.
	for _, e := range input.Entries {
		entryMap, ok := e.(map[string]interface{})
		if !ok {
			failedCount++
			resultEntries = append(resultEntries, map[string]interface{}{
				"ErrorCode":    "InvalidArgument",
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
					"ErrorCode":    "InvalidArgument",
					"ErrorMessage": "EventBusName must not be empty",
				})
				failedCount++
				continue
			}
			// EventBusName accepts "the name or ARN of the event bus"
			// (model member documentation); rules and archives are keyed
			// by the canonical name, so the ARN form is resolved before
			// any lookup.
			eventBusName = resolveEventBusName(busStr)
		}

		// DetailType max 128 chars per AWS EventBridge PutEvents API reference.
		if !validateDetailType(detailType) {
			resultEntries = append(resultEntries, map[string]interface{}{
				"ErrorCode":    "InvalidArgument",
				"ErrorMessage": "DetailType must be between 1 and 128 characters",
			})
			failedCount++
			continue
		}

		// Source max 256 chars (AWS PutEventsRequestEntry).
		if !validateSource(source) {
			resultEntries = append(resultEntries, map[string]interface{}{
				"ErrorCode":    "InvalidArgument",
				"ErrorMessage": "Source must be between 1 and 256 characters",
			})
			failedCount++
			continue
		}

		// TraceHeader max 500 chars per Smithy @length trait.
		if !validateTraceHeader(traceHeader) {
			resultEntries = append(resultEntries, map[string]interface{}{
				"ErrorCode":    "InvalidArgument",
				"ErrorMessage": "TraceHeader must be at most 500 characters",
			})
			failedCount++
			continue
		}

		// A non-existent bus is not an entry failure: "If you use PutEvents
		// to publish an event to an event bus that does not exist,
		// EventBridge event matching will not find a corresponding rule and
		// will drop the event. Although EventBridge will send a 200
		// response, it will not fail the request or include the event in
		// the FailedEntryCount" (user guide, Handling failures with
		// PutEvents). The event therefore proceeds to rule matching on the
		// named bus, matches nothing, and reports success; the internal
		// bus ingress plane (handlePutEventsEvent) reports the same
		// condition as a handler failure instead, because internal
		// publishers depend on the delivery verdict for retry and
		// dead-letter routing.

		if source == "" || detailType == "" || detailStr == "" {
			// The model documentation says "EventBridge fails that entry"
			// without naming a code; InvalidArgument ("A specified
			// parameter is not valid") is the vocabulary's parameter-invalid
			// code.
			resultEntries = append(resultEntries, map[string]interface{}{
				"ErrorCode":    "InvalidArgument",
				"ErrorMessage": "Source, DetailType, and Detail are required",
			})
			failedCount++
			continue
		}

		var detail map[string]interface{}
		parsedDetail, detailErr := parseDetailObject(detailStr)
		if detailErr != nil {
			resultEntries = append(resultEntries, map[string]interface{}{
				"ErrorCode":    "MalformedDetail",
				"ErrorMessage": "Detail must be a valid JSON object",
			})
			failedCount++
			continue
		}
		detail = parsedDetail

		eventTime := time.Now().UTC()
		if timeVal, ok := entryMap["Time"]; ok {
			parsed, err := parseEntryTimestamp(timeVal)
			if err != nil {
				resultEntries = append(resultEntries, map[string]interface{}{
					"ErrorCode":    "InvalidArgument",
					"ErrorMessage": err.Error(),
				})
				failedCount++
				continue
			}
			eventTime = parsed
		}

		event := &eventsstore.Event{
			ID:         generateEventID(),
			Version:    "0",
			DetailType: detailType,
			Source:     source,
			Account:    s.accountID,
			Time:       eventTime,
			// The receipt stamp is the request's arrival at the platform,
			// distinct from the publisher-suppliable Time above.
			IngestionTime: time.Now().UTC(),
			Region:        input.Region,
			Detail:        detail,
			EventBusName:  eventBusName,
			TraceHeader:   traceHeader,
		}

		if rv, ok := entryMap["Resources"]; ok {
			resources, ok := rv.([]interface{})
			if !ok {
				resultEntries = append(resultEntries, map[string]interface{}{
					"ErrorCode":    "InvalidArgument",
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
					"ErrorCode":    "InvalidArgument",
					"ErrorMessage": "Resources must be an array of strings",
				})
				failedCount++
				continue
			}
		}

		if err := s.deliverEvent(ctx, store, event, eventBusName, input.Region, 0); err != nil {
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

// parseDetailObject enforces the model's Detail contract on the string
// form: the member is "a valid JSON object", so the value must unmarshal to
// a JSON object — a scalar, an array or malformed JSON is rejected. The
// null literal is a scalar too: it unmarshals into a nil map without an
// error, so the nil result is rejected explicitly. Both ingress planes
// (the PutEvents API's per-entry loop and the internal bus ingress)
// resolve the string form through this single site, so the acceptance
// vocabulary cannot drift between them.
func parseDetailObject(detailStr string) (map[string]interface{}, error) {
	var detail map[string]interface{}
	if err := json.Unmarshal([]byte(detailStr), &detail); err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, fmt.Errorf("Detail must be a valid JSON object")
	}
	return detail, nil
}

// deliverEvent is the unified event ingress: it archives the event onto the
// bus's enabled archives and delivers it to the matching rules. Every path
// an event enters a bus funnels through here — the PutEvents API plane, the
// internal bus ingress, bus-target delivery, scheduled fires — so the
// archive pattern filter is the only stated ingress discriminator ("EventBridge
// sends events that match the event pattern to the archive", the archives
// user guide), never the ingress path. Replayed events carry the replay-name
// stamp and are excluded by archiveEvent itself, so the replay path can
// reuse this ingress unchanged. hopDepth carries the cross-bus hop count as
// in deliverEventWithStore.
func (s *EventsService) deliverEvent(ctx context.Context, store *eventsstore.EventsStore, event *eventsstore.Event, eventBusName, region string, hopDepth int) error {
	s.archiveEvent(ctx, store, event, eventBusName)

	return s.deliverEventWithStore(ctx, region, event, eventBusName, store, hopDepth)
}

// deliverEventWithStore delivers an event to every enabled rule on the event
// bus whose event pattern matches, fanning the matched rules' targets out to
// the delivery plane. hopDepth is the cross-bus hop count the event has
// already accumulated; ingress paths pass 0, the bus-target delivery path
// passes its received depth plus one, and the counter guards against
// mutually-targeting bus cycles.
func (s *EventsService) deliverEventWithStore(ctx context.Context, region string, event *eventsstore.Event, eventBusName string, store *eventsstore.EventsStore, hopDepth int) error {
	return s.deliverEventToBusRules(ctx, region, event, eventBusName, store, hopDepth, nil)
}

// deliverEventToBusRules is the rule fan-out behind deliverEventWithStore.
// filterRuleARNs, when non-nil, restricts the fan-out to the named rules —
// the replay destination's FilterArns ("A list of ARNs for rules to replay
// events to", Smithy ReplayDestination.FilterArns member documentation): a
// filtered replay fires only the chosen rules' targets while pattern
// matching, target listing and dispatch stay the unfiltered path's.
func (s *EventsService) deliverEventToBusRules(ctx context.Context, region string, event *eventsstore.Event, eventBusName string, store *eventsstore.EventsStore, hopDepth int, filterRuleARNs map[string]bool) error {
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

	for _, rule := range allRules {
		if rule.State != eventsstore.RuleStateEnabled && rule.State != eventsstore.RuleStateEnabledWithAllCloudtrailManagementEvents {
			continue
		}

		if filterRuleARNs != nil && !filterRuleARNs[rule.ARN] {
			continue
		}

		// A rule without an EventPattern is a scheduled rule (PutRule
		// requires at least one of EventPattern and ScheduleExpression):
		// a scheduled rule delivers on its own boundaries and never
		// matches bus events, so the fan-out skips it here — the
		// scheduler dispatches a scheduled fire's own targets directly.
		if rule.EventPattern == "" {
			continue
		}
		if !s.matchEventPattern(event, rule.EventPattern) {
			continue
		}

		if err := s.dispatchRuleTargets(ctx, region, event, rule, store, hopDepth); err != nil {
			return err
		}
	}

	return nil
}

// dispatchRuleTargets fans one matched rule's targets out to the delivery
// plane. It is the single target-dispatch path for both rule matching
// (deliverEventWithStore) and the scheduler (fireScheduledRule); hopDepth
// carries the cross-bus hop count as in deliverEventWithStore.
func (s *EventsService) dispatchRuleTargets(ctx context.Context, region string, event *eventsstore.Event, rule *eventsstore.Rule, store *eventsstore.EventsStore, hopDepth int) error {
	// Paginate through all targets for the rule.
	var allTargets []*eventsstore.Target
	targetsNextToken := ""
	for {
		targetsResult, err := store.ListTargetsByRule(ctx, rule.EventBusName, rule.Name, 100, targetsNextToken)
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

	var targetWg sync.WaitGroup
	var dispatchErrOnce sync.Once
	var dispatchErr error
	for _, target := range allTargets {
		targetCopy := *target
		payloadBytes := s.buildTargetPayload(rule.ARN, rule.Name, event, targetCopy)

		if s.bus != nil {
			ebEvt := &eventbus.EventBridgeDeliveryEvent{
				RuleARN:            rule.ARN,
				TargetID:           targetCopy.ID,
				TargetARN:          targetCopy.ARN,
				Input:              payloadBytes,
				EventBridgeEventID: event.ID,
				EventIngestionTime: event.IngestionTime,
				HopDepth:           hopDepth,
				TraceHeader:        event.TraceHeader,
			}
			if targetCopy.DeadLetterConfig != nil {
				ebEvt.DeadLetterConfigArn = targetCopy.DeadLetterConfig.Arn
			}
			if targetCopy.RetryPolicy != nil {
				ebEvt.MaximumRetryAttempts = targetCopy.RetryPolicy.MaximumRetryAttempts
				ebEvt.MaximumEventAgeInSeconds = targetCopy.RetryPolicy.MaximumEventAgeInSeconds
				ebEvt.RetryPolicySet = true
			}
			if targetCopy.SqsParameters != nil {
				ebEvt.SqsMessageGroupId = targetCopy.SqsParameters.MessageGroupId
			}
			// The partition-key path references the original event
			// ("dynamic path parameters must reference the original
			// event, not the transformed event"), so it is resolved
			// here where the full event is in hand and the resolved
			// key rides the bus.
			if key := kinesisPartitionKeyFor(event, &targetCopy); key != "" {
				ebEvt.KinesisPartitionKey = key
			}
			if targetCopy.AppSyncParameters != nil {
				ebEvt.AppSyncGraphQLOperation = targetCopy.AppSyncParameters.GraphQLOperation
			}
			if targetCopy.HttpParameters != nil {
				if encoded, err := json.Marshal(targetCopy.HttpParameters); err == nil {
					ebEvt.TargetHttpParameters = encoded
				}
			}
			ebEvt.Region = region
			if err := s.bus.Publish(context.Background(), ebEvt); err != nil {
				logs.Warn("failed to publish event directly", logs.String("targetArn", targetCopy.ARN), logs.Err(err))
			}
		} else {
			// The synchronous arm honours both the caller's context
			// and the engine's shutdown — whichever ends first aborts
			// the retry sleeps inside dispatchToTarget, so service
			// shutdown does not wait out a delivery deadline.
			dctx, dcancel := eitherContext(ctx, s.delivery.ctx)
			targetWg.Add(1)
			select {
			case s.targetSemaphore <- struct{}{}:
				go func(evt *eventsstore.Event, tgt eventsstore.Target) {
					defer dcancel()
					defer func() {
						<-s.targetSemaphore
						targetWg.Done()
						if r := recover(); r != nil {
							logs.Error("eventbridge: panic delivering to target", logs.String("arn", tgt.ARN), logs.Any("panic", r))
						}
					}()
					if err := s.deliverToTarget(dctx, region, rule.ARN, rule.Name, evt, tgt); err != nil {
						// The bus-less arm has no transport verdict to carry
						// the loss, so the terminal error is reported here:
						// the first failure becomes the dispatch result the
						// callers (the PutEvents entry, the scheduled fire)
						// surface, and every failed target is logged.
						logs.Error("eventbridge: synchronous target delivery failed",
							logs.String("targetArn", tgt.ARN),
							logs.String("eventId", evt.ID),
							logs.Err(err))
						dispatchErrOnce.Do(func() { dispatchErr = err })
					}
				}(event, targetCopy)
			case <-dctx.Done():
				dcancel()
				targetWg.Done()
				goto done
			}
		}
	}

done:
	targetWg.Wait()
	return dispatchErr
}

// archiveEvent stores the event into every enabled archive of the event bus
// whose event pattern matches, updating the archives' counters. Replayed
// events never enter an archive: the replay-name metadata field is the
// documented replay discriminator (the archives user guide's managed
// {"replay-name":[{"exists":false}]} pattern), and excluding by the stamp
// holds for every ingress — including a replay delivered onto its own source
// bus, which must not grow the archive it replays from.
func (s *EventsService) archiveEvent(ctx context.Context, store *eventsstore.EventsStore, event *eventsstore.Event, eventBusName string) {
	if event.ReplayName != "" {
		return
	}

	if _, err := store.GetEventBus(ctx, eventBusName); err != nil {
		// Archiving silently stops here; without the warning the stop is
		// invisible until a later DescribeArchive shows a stale EventCount.
		logs.Warn("event archiving skipped: event bus lookup failed",
			logs.String("eventBusName", eventBusName),
			logs.Err(err))
		return
	}

	archives, err := store.ListArchivesForEventBus(ctx, eventBusName)
	if err != nil {
		logs.Warn("event archiving skipped: archive listing failed",
			logs.String("eventBusName", eventBusName),
			logs.Err(err))
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

		eventMap := eventEnvelope(event)

		archivedEvent := &eventsstore.ArchivedEvent{
			ID:        event.ID,
			Event:     eventMap,
			Timestamp: event.Time,
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
