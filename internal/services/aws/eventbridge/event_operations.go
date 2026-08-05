package eventbridge

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	arnutil "vorpalstacks/internal/utils/aws/arn"
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
		eventBusName, _ := entryMap["EventBusName"].(string)
		if eventBusName == "" {
			eventBusName = "default"
		}

		// M10: DetailType max 128 chars (AWS API Reference).
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

		// L1: TraceHeader max 500 chars (Smithy @length).
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
				break
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

func (s *EventsService) matchEventPattern(event *eventsstore.Event, pattern string) bool {
	var patternMap map[string]interface{}
	if err := json.Unmarshal([]byte(pattern), &patternMap); err != nil {
		return false
	}

	eventMap := map[string]interface{}{
		"version":     event.Version,
		"id":          event.ID,
		"source":      event.Source,
		"detail-type": event.DetailType,
		"time":        event.Time.Format(time.RFC3339),
		"region":      event.Region,
		"resources":   event.Resources,
		"detail":      event.Detail,
		"account":     event.Account,
	}

	for key, patternValue := range patternMap {
		eventValue, exists := eventMap[key]
		if !exists {
			if isExistsFalsePattern(patternValue) {
				continue
			}
			return false
		}

		if !s.matchValue(eventValue, patternValue) {
			return false
		}
	}

	return true
}

// isExistsFalsePattern checks whether a pattern value is equivalent to
// {"exists": false} (either directly or as a single-element array),
// which should match when the field is absent from the event.
func isExistsFalsePattern(patternValue interface{}) bool {
	check := func(obj map[string]interface{}) bool {
		if len(obj) != 1 {
			return false
		}
		if op, ok := obj["exists"]; ok {
			if b, ok := op.(bool); ok && !b {
				return true
			}
		}
		return false
	}
	switch p := patternValue.(type) {
	case map[string]interface{}:
		return check(p)
	case []interface{}:
		if len(p) == 1 {
			if item, ok := p[0].(map[string]interface{}); ok {
				return check(item)
			}
		}
	}
	return false
}

func (s *EventsService) matchValue(eventValue, patternValue interface{}) bool {
	switch p := patternValue.(type) {
	case []interface{}:
		for _, item := range p {
			if s.matchValue(eventValue, item) {
				return true
			}
		}
		return false
	case string:
		evStr, ok := eventValue.(string)
		if !ok {
			return false
		}
		return evStr == p
	case map[string]interface{}:
		if len(p) == 1 {
			for key, operand := range p {
				if isKnownOperator(key) {
					return s.matchOperator(eventValue, key, operand)
				}
				break
			}
		}
		evMap, ok := eventValue.(map[string]interface{})
		if !ok {
			return false
		}
		for k, v := range p {
			if !s.matchValue(evMap[k], v) {
				return false
			}
		}
		return true
	default:
		return fmt.Sprintf("%v", eventValue) == fmt.Sprintf("%v", patternValue)
	}
}

func (s *EventsService) matchOperator(eventValue interface{}, op string, operand interface{}) bool {
	switch op {
	case "prefix":
		evStr, ok := eventValue.(string)
		if !ok {
			return false
		}
		prefix, ok := operand.(string)
		if !ok {
			return false
		}
		return strings.HasPrefix(evStr, prefix)
	case "numeric":
		evNum, ok := toFloat64(eventValue)
		if !ok {
			return false
		}
		operands, ok := operand.([]interface{})
		if !ok || len(operands) < 2 {
			return false
		}
		for i := 0; i < len(operands)-1; i++ {
			compOp, ok := operands[i].(string)
			if !ok {
				return false
			}
			compVal, ok := toFloat64(operands[i+1])
			if !ok {
				return false
			}
			if !compareNumeric(evNum, compOp, compVal) {
				return false
			}
			i++
		}
		return true
	case "anything-but":
		return !s.matchValue(eventValue, operand)
	case "exists":
		existsVal, ok := operand.(bool)
		if !ok {
			return false
		}
		if existsVal {
			return eventValue != nil
		}
		return eventValue == nil
	case "suffix":
		evStr, ok := eventValue.(string)
		if !ok {
			return false
		}
		suffix, ok := operand.(string)
		if !ok {
			return false
		}
		return strings.HasSuffix(evStr, suffix)
	case "equals-ignore-case":
		evStr, ok := eventValue.(string)
		if !ok {
			return false
		}
		operandStr, ok := operand.(string)
		if !ok {
			return false
		}
		return strings.EqualFold(evStr, operandStr)
	case "wildcard":
		evStr, ok := eventValue.(string)
		if !ok {
			return false
		}
		pattern, ok := operand.(string)
		if !ok {
			return false
		}
		return matchWildcardPattern(evStr, pattern)
	case "cidr":
		evStr, ok := eventValue.(string)
		if !ok {
			return false
		}
		cidr, ok := operand.(string)
		if !ok {
			return false
		}
		return matchCIDRBlock(evStr, cidr)
	default:
		return false
	}
}

func matchWildcardPattern(s, pattern string) bool {
	for len(pattern) > 0 && pattern[0] == '*' {
		pattern = pattern[1:]
	}
	if pattern == "" {
		return true
	}
	for len(pattern) > 0 {
		starIdx := -1
		for i, c := range pattern {
			if c == '*' {
				starIdx = i
				break
			}
		}
		if starIdx == -1 {
			return strings.HasSuffix(s, pattern)
		}
		fragment := pattern[:starIdx]
		pattern = pattern[starIdx+1:]
		for len(pattern) > 0 && pattern[0] == '*' {
			pattern = pattern[1:]
		}
		idx := strings.Index(s, fragment)
		if idx == -1 {
			return false
		}
		if len(pattern) == 0 {
			// Pattern ended with '*' after this fragment — everything
			// following the fragment matches. We already confirmed the
			// fragment exists in s via Index, so the match succeeds.
			return true
		}
		s = s[idx+len(fragment):]
	}
	return true
}

func matchCIDRBlock(ipStr, cidr string) bool {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ipNet.Contains(ip)
}

func isKnownOperator(key string) bool {
	switch key {
	case "prefix", "numeric", "anything-but", "exists", "suffix", "equals-ignore-case", "wildcard", "cidr":
		return true
	default:
		return false
	}
}

func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	default:
		return 0, false
	}
}

func compareNumeric(val float64, op string, comp float64) bool {
	switch op {
	case "<":
		return val < comp
	case "<=":
		return val <= comp
	case ">":
		return val > comp
	case ">=":
		return val >= comp
	case "=":
		return val == comp
	default:
		return false
	}
}

func (s *EventsService) deliverToTarget(ctx context.Context, region string, event *eventsstore.Event, target eventsstore.Target) {
	payload := s.buildTargetPayload(event, target)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logs.Error("failed to marshal event payload",
			logs.String("eventId", event.ID),
			logs.String("targetId", target.ID),
			logs.Err(err))
		return
	}

	s.dispatchToTarget(ctx, region, event, target, payloadBytes)
}

const (
	// AWS production defaults: 185 retries / 86400 s (24 h) deadline,
	// 1 s initial backoff, 60 s max backoff. These match the AWS
	// EventBridge retry behaviour for target delivery.
	prodMaxRetryAttempts     = 185
	prodMaxEventAgeInSeconds = 86400
	prodRetryInitialBackoff  = 1 * time.Second
	prodRetryMaxBackoff      = 60 * time.Second

	// TEST_MODE defaults: reduced limits so integration tests are not
	// blocked for hours when a target is temporarily unreachable. The
	// delivery semaphore (targetConcurrencyLimit = 100) can be
	// exhausted by retrying goroutines under the production defaults,
	// stalling the entire delivery pipeline during tests.
	testMaxRetryAttempts     = 3
	testMaxEventAgeInSeconds = 60
	testRetryInitialBackoff  = 500 * time.Millisecond
	testRetryMaxBackoff      = 5 * time.Second
)

// retryDefaults returns the retry constants appropriate for the current
// runtime mode. In TEST_MODE, reduced limits prevent goroutine
// exhaustion during integration tests. In production, AWS-compatible
// defaults ensure transient target failures are retried for up to 24 h.
func retryDefaults() (maxRetry int32, maxAge int32, initialBackoff, maxBackoff time.Duration) {
	if os.Getenv("TEST_MODE") == "true" {
		return testMaxRetryAttempts, testMaxEventAgeInSeconds, testRetryInitialBackoff, testRetryMaxBackoff
	}
	return prodMaxRetryAttempts, prodMaxEventAgeInSeconds, prodRetryInitialBackoff, prodRetryMaxBackoff
}

func (s *EventsService) dispatchToTarget(ctx context.Context, region string, event *eventsstore.Event, target eventsstore.Target, payload []byte) {
	targetType := s.parseTargetType(target.ARN)

	maxRetries, maxAge, defaultBackoff, maxBackoff := retryDefaults()
	if target.RetryPolicy != nil {
		if target.RetryPolicy.MaximumRetryAttempts >= 0 && target.RetryPolicy.MaximumRetryAttempts <= 185 {
			maxRetries = target.RetryPolicy.MaximumRetryAttempts
		}
		if target.RetryPolicy.MaximumEventAgeInSeconds >= 60 && target.RetryPolicy.MaximumEventAgeInSeconds <= 86400 {
			maxAge = target.RetryPolicy.MaximumEventAgeInSeconds
		}
	}
	deadline := time.Now().Add(time.Duration(maxAge) * time.Second)

	attempt := int32(0)
	backoff := defaultBackoff
	var deliverErr error

	for {
		deliverErr = nil
		switch targetType {
		case "lambda":
			deliverErr = s.deliverToLambda(ctx, region, event.ID, target.ARN, payload)
		case "sqs":
			deliverErr = s.deliverToSQS(ctx, region, target, payload)
		case "sns":
			deliverErr = s.deliverToSNS(ctx, region, target.ARN, payload)
		case "logs":
			deliverErr = s.deliverToCloudWatchLogs(ctx, region, event.ID, target.ARN, payload)
		case "states":
			deliverErr = s.deliverToStepFunctions(ctx, region, target.ARN, payload)
		case "kinesis":
			deliverErr = s.deliverToKinesis(ctx, region, target, payload)
		case "firehose":
			deliverErr = s.deliverToFirehose(ctx, region, target.ARN, payload)
		case "ecs":
			deliverErr = s.deliverToECS(ctx, region, target.ARN, payload)
		case "events":
			deliverErr = s.deliverToEventBus(ctx, region, event, target.ARN)
		default:
			deliverErr = fmt.Errorf("target type %q not implemented", targetType)
		}

		if deliverErr == nil {
			return
		}

		attempt++
		if attempt >= maxRetries || time.Now().After(deadline) {
			break
		}

		// Exponential backoff with jitter.
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
		jitter := time.Duration(rand.Int63n(int64(backoff) / 2))
		sleepDur := backoff + jitter
		select {
		case <-time.After(sleepDur):
		case <-ctx.Done():
			return
		}
		backoff *= 2
	}

	logs.Error("event delivery to target failed after retries",
		logs.String("targetArn", target.ARN),
		logs.String("eventId", event.ID),
		logs.Int("attempts", int(attempt)),
		logs.Err(deliverErr))
	s.routeToDeadLetter(ctx, region, event, target, payload)
}

// routeToDeadLetter delivers the event payload to the configured dead-letter
// queue (SQS or SNS) when the primary target delivery fails.
func (s *EventsService) routeToDeadLetter(ctx context.Context, region string, event *eventsstore.Event, target eventsstore.Target, payload []byte) {
	if target.DeadLetterConfig == nil || target.DeadLetterConfig.Arn == "" {
		return
	}
	dlqArn := target.DeadLetterConfig.Arn
	dlqType := s.parseTargetType(dlqArn)
	switch dlqType {
	case "sqs":
		s.deliverToSQS(ctx, region, eventsstore.Target{ARN: dlqArn}, payload)
		logs.Info("event routed to DLQ (SQS)",
			logs.String("dlqArn", dlqArn),
			logs.String("eventId", event.ID),
			logs.String("originalTarget", target.ARN))
	case "sns":
		s.deliverToSNS(ctx, region, dlqArn, payload)
		logs.Info("event routed to DLQ (SNS)",
			logs.String("dlqArn", dlqArn),
			logs.String("eventId", event.ID),
			logs.String("originalTarget", target.ARN))
	default:
		logs.Warn("DLQ target type not supported",
			logs.String("dlqArn", dlqArn),
			logs.String("dlqType", dlqType),
			logs.String("eventId", event.ID))
	}
}

func (s *EventsService) parseTargetType(arnStr string) string {
	_, service, _, _, _ := arnutil.SplitARN(arnStr)
	return service
}

func (s *EventsService) deliverToLambda(ctx context.Context, region string, eventID string, targetArn string, payload []byte) error {
	if s.bus == nil || s.bus.LambdaInvoker() == nil {
		return fmt.Errorf("lambda invoker not configured")
	}

	if s.bus != nil {
		allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, targetArn, "lambda", "events.amazonaws.com", "lambda:InvokeFunction", targetArn)
		if evalErr != nil {
			return fmt.Errorf("resource policy evaluation failed for Lambda target: %w", evalErr)
		}
		if !allowed {
			return fmt.Errorf("resource policy denied Lambda invocation")
		}
	}

	functionName := arnutil.ExtractFunctionNameFromARN(targetArn)
	if functionName == "" {
		return fmt.Errorf("failed to extract function name from ARN %s", targetArn)
	}

	statusCode, result, err := s.bus.LambdaInvoker().InvokeForGateway(ctx, targetArn, payload)
	if err != nil {
		return fmt.Errorf("failed to invoke Lambda function %s: %w", functionName, err)
	}

	if statusCode != 200 {
		return fmt.Errorf("Lambda invocation returned status %d: %s", statusCode, string(result))
	}

	logs.Debug("event delivered to Lambda successfully",
		logs.String("eventId", eventID),
		logs.String("functionName", functionName))
	return nil
}

func (s *EventsService) extractInputPath(payload map[string]interface{}, inputPath string) map[string]interface{} {
	if inputPath == "" || inputPath == "$" {
		return payload
	}

	path := strings.TrimPrefix(inputPath, "$.")
	parts := strings.Split(path, ".")

	current := interface{}(payload)
	for _, part := range parts {
		if part == "" {
			continue
		}
		if current == nil {
			return payload
		}
		m, ok := current.(map[string]interface{})
		if !ok {
			return payload
		}
		val, exists := m[part]
		if !exists {
			return payload
		}
		current = val
	}

	if result, ok := current.(map[string]interface{}); ok {
		return result
	}
	return map[string]interface{}{"value": current}
}

func (s *EventsService) applyInputTransformer(payload map[string]interface{}, transformer *eventsstore.InputTransformer) map[string]interface{} {
	if transformer == nil || transformer.InputTemplate == "" {
		return payload
	}

	// Build values from InputPathsMap
	values := make(map[string]interface{})
	if transformer.InputPathsMap != nil {
		for key, path := range transformer.InputPathsMap {
			values[key] = s.extractValueByPath(payload, path)
		}
	}

	// Apply template - simple replacement of <key> placeholders
	template := transformer.InputTemplate
	for key, value := range values {
		placeholder := "<" + key + ">"
		var valueStr string
		switch v := value.(type) {
		case string:
			valueStr = v
		case nil:
			valueStr = "null"
		default:
			b, _ := json.Marshal(v)
			valueStr = string(b)
		}
		template = strings.ReplaceAll(template, placeholder, valueStr)
	}

	// Try to parse as JSON, otherwise return as raw template
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(template), &parsed); err == nil {
		return parsed
	}

	return map[string]interface{}{"message": template}
}

func (s *EventsService) extractValueByPath(payload map[string]interface{}, path string) interface{} {
	if path == "" || path == "$" {
		return payload
	}

	path = strings.TrimPrefix(path, "$.")
	parts := strings.Split(path, ".")

	current := interface{}(payload)
	for _, part := range parts {
		if part == "" {
			continue
		}
		if current == nil {
			return nil
		}
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current = m[part]
	}
	return current
}

func (s *EventsService) deliverToSQS(ctx context.Context, region string, target eventsstore.Target, payload []byte) error {
	if s.bus == nil || s.bus.SQSInvoker() == nil {
		return fmt.Errorf("SQS invoker not configured")
	}

	arnStr := target.ARN
	queueName := arnutil.ExtractQueueNameFromARN(arnStr)
	if queueName == "" {
		return fmt.Errorf("failed to extract queue name from ARN %s", arnStr)
	}

	queueURL, qErr := s.bus.SQSInvoker().GetQueueByName(ctx, queueName)
	if qErr != nil {
		return fmt.Errorf("queue not found for SQS delivery %s: %w", queueName, qErr)
	}

	queueARN, arnErr := s.bus.SQSInvoker().GetQueueARN(ctx, queueURL)
	if arnErr != nil {
		return fmt.Errorf("failed to get queue ARN: %w", arnErr)
	}

	allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, queueARN, "sqs", "events.amazonaws.com", "sqs:SendMessage", queueARN)
	if evalErr != nil {
		return fmt.Errorf("resource policy evaluation failed for SQS target: %w", evalErr)
	}
	if !allowed {
		return fmt.Errorf("resource policy denied SQS SendMessage")
	}

	opts := eventbus.SQSSendOptions{}
	if target.SqsParameters != nil && target.SqsParameters.MessageGroupId != "" {
		opts.MessageGroupID = target.SqsParameters.MessageGroupId
	}

	if _, _, err := s.bus.SQSInvoker().SendMessage(ctx, queueURL, string(payload), opts); err != nil {
		return fmt.Errorf("failed to deliver event to SQS %s: %w", queueName, err)
	}

	logs.Debug("Event delivered to SQS successfully",
		logs.String("queue", queueName))
	return nil
}

func (s *EventsService) deliverToSNS(ctx context.Context, region string, arnStr string, payload []byte) error {
	if s.bus == nil || s.bus.SNSInvoker() == nil {
		return fmt.Errorf("SNS invoker not configured")
	}

	allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, arnStr, "sns", "events.amazonaws.com", "sns:Publish", arnStr)
	if evalErr != nil {
		return fmt.Errorf("resource policy evaluation failed for SNS target: %w", evalErr)
	}
	if !allowed {
		return fmt.Errorf("resource policy denied SNS Publish")
	}

	_, _, _, resource, _ := arnutil.SplitARN(arnStr)
	if resource == "" {
		return fmt.Errorf("failed to extract topic name from ARN %s", arnStr)
	}

	topicName := strings.TrimPrefix(resource, "topic/")

	_, err := s.bus.SNSInvoker().PublishToTopic(ctx, arnStr, string(payload), "", nil)
	if err != nil {
		return fmt.Errorf("failed to deliver event to SNS %s: %w", arnStr, err)
	}

	logs.Debug("Event delivered to SNS successfully",
		logs.String("topic", topicName),
		logs.String("arn", arnStr))
	return nil
}

func (s *EventsService) deliverToCloudWatchLogs(ctx context.Context, region string, eventID string, arnStr string, payload []byte) error {
	if s.bus == nil {
		return fmt.Errorf("event bus not configured")
	}

	_, _, _, _, resource := arnutil.SplitARN(arnStr)
	logGroup := arnutil.ExtractLogGroupNameFromARN(arnStr)
	if logGroup == "" {
		return fmt.Errorf("failed to extract log group from CloudWatch Logs ARN %s", arnStr)
	}

	logStream := resource
	if idx := strings.LastIndex(resource, ":log-stream:"); idx != -1 {
		logStream = resource[idx+12:]
	}
	if logStream == resource {
		logStream = fmt.Sprintf("eventbridge-%s", eventID)
	}

	evt := &eventbus.CloudWatchLogsPutEvent{
		LogGroup:  logGroup,
		LogStream: logStream,
		LogEvents: []eventbus.LogEntry{
			{Timestamp: time.Now().UnixMilli(), Message: string(payload)},
		},
	}
	evt.Region = region
	evt.AccountID = s.accountID

	if err := s.bus.Publish(ctx, evt); err != nil {
		return fmt.Errorf("failed to deliver event to CloudWatch Logs %s: %w", logGroup, err)
	}

	logs.Debug("Event delivered to CloudWatch Logs successfully",
		logs.String("logGroup", logGroup),
		logs.String("logStream", logStream))
	return nil
}

func (s *EventsService) deliverToStepFunctions(ctx context.Context, region string, targetArn string, payload []byte) error {
	if s.bus == nil {
		return fmt.Errorf("event bus not configured")
	}

	_, _, smRegion, _, _ := arnutil.SplitARN(targetArn)
	if smRegion == "" {
		smRegion = region
	}

	evt := &eventbus.StepFunctionsStartExecutionEvent{
		StateMachineArn: targetArn,
		Input:           string(payload),
	}
	evt.Region = smRegion
	evt.AccountID = s.accountID
	if err := s.bus.Publish(ctx, evt); err != nil {
		return fmt.Errorf("failed to publish Step Functions start event for %s: %w", targetArn, err)
	}

	logs.Debug("event delivered to Step Functions successfully",
		logs.String("targetArn", targetArn))
	return nil
}

func (s *EventsService) deliverToKinesis(ctx context.Context, region string, target eventsstore.Target, payload []byte) error {
	if s.bus == nil || s.bus.KinesisInvoker() == nil {
		return fmt.Errorf("Kinesis invoker not configured")
	}

	targetArn := target.ARN
	_, _, _, _, resource := arnutil.SplitARN(targetArn)

	streamName := resource
	if idx := strings.Index(resource, "stream/"); idx != -1 {
		streamName = resource[idx+len("stream/"):]
	}

	partitionKey := fmt.Sprintf("eventbridge-%s", generateEventID())
	if target.KinesisParameters != nil && target.KinesisParameters.PartitionKeyPath != "" {
		var payloadMap map[string]interface{}
		if err := json.Unmarshal(payload, &payloadMap); err == nil {
			if val := s.extractValueByPath(payloadMap, target.KinesisParameters.PartitionKeyPath); val != nil {
				if valStr, ok := val.(string); ok && valStr != "" {
					partitionKey = valStr
				}
			}
		}
	}

	// The local Kinesis service stores data as-is and GetRecords returns it
	// without additional encoding. SDK clients expect base64-encoded Data in
	// the GetRecords response, so cross-service callers must pre-encode the
	// payload to match the format that the Kinesis SDK PutRecord would send.
	encodedPayload := base64.StdEncoding.EncodeToString(payload)
	_, err := s.bus.KinesisInvoker().PutRecord(ctx, streamName, partitionKey, []byte(encodedPayload))
	if err != nil {
		return fmt.Errorf("failed to put record to Kinesis stream %s: %w", streamName, err)
	}

	logs.Debug("event delivered to Kinesis successfully",
		logs.String("stream", streamName))
	return nil
}

func (s *EventsService) deliverToFirehose(ctx context.Context, region string, targetArn string, payload []byte) error {
	logs.Error("Firehose target delivery failed: Firehose service is not available",
		logs.String("targetArn", targetArn),
		logs.String("region", region))
	return fmt.Errorf("firehose delivery target %s is not available in this deployment", targetArn)
}

func (s *EventsService) deliverToECS(ctx context.Context, region string, targetArn string, payload []byte) error {
	logs.Error("ECS target delivery failed: ECS service is not available",
		logs.String("targetArn", targetArn),
		logs.String("region", region))
	return fmt.Errorf("ecs delivery target %s is not available in this deployment", targetArn)
}

// deliveryDepthKey is used to track cross-bus delivery depth via context,
// preventing infinite loops when bus A targets bus B and vice versa.
type deliveryDepthKey struct{}

const maxCrossBusDepth = 10

// deliverToEventBus delivers an event to a cross-account or cross-region
// event bus target.  The target ARN identifies the destination bus.  A
// depth counter (propagated via context) prevents infinite delivery loops.
func (s *EventsService) deliverToEventBus(ctx context.Context, sourceRegion string, event *eventsstore.Event, targetARN string) error {
	depth := 0
	if v, ok := ctx.Value(deliveryDepthKey{}).(int); ok {
		depth = v
	}
	if depth >= maxCrossBusDepth {
		return fmt.Errorf("maximum cross-bus delivery depth (%d) exceeded for event %s", maxCrossBusDepth, event.ID)
	}

	_, _, targetRegion, _, resource := arnutil.SplitARN(targetARN)
	busName := strings.TrimPrefix(resource, "event-bus/")
	if busName == "" {
		return fmt.Errorf("invalid event bus target ARN: %s", targetARN)
	}

	store, err := s.GetStoreForRegion(targetRegion)
	if err != nil {
		return fmt.Errorf("failed to get store for region %s: %w", targetRegion, err)
	}

	childCtx := context.WithValue(ctx, deliveryDepthKey{}, depth+1)

	s.archiveEvent(childCtx, store, event, busName)

	return s.deliverEventWithStore(childCtx, targetRegion, event, busName, store)
}

// TestEventPattern tests whether an event pattern matches a given event.
func (s *EventsService) TestEventPattern(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	patternStr := request.GetStringParam(req.Parameters, "EventPattern")
	eventStr := request.GetStringParam(req.Parameters, "Event")

	if patternStr == "" {
		return nil, awserrors.NewValidationException("Parameter EventPattern is required")
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
