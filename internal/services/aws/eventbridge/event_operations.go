package eventbridge

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
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

	resultEntries := make([]map[string]interface{}, 0)
	failedCount := int32(0)

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
		eventBusName, _ := entryMap["EventBusName"].(string)
		if eventBusName == "" {
			eventBusName = "default"
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

		event := &eventsstore.Event{
			ID:           generateEventID(),
			Version:      "0",
			DetailType:   detailType,
			Source:       source,
			Account:      s.accountID,
			Time:         time.Now().UTC(),
			Region:       reqCtx.GetRegion(),
			Detail:       detail,
			EventBusName: eventBusName,
		}

		if resources, ok := entryMap["Resources"].([]interface{}); ok {
			for _, r := range resources {
				if rStr, ok := r.(string); ok {
					event.Resources = append(event.Resources, rStr)
				}
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
		if rule.State != eventsstore.RuleStateEnabled {
			continue
		}

		if rule.EventPattern != "" {
			if !s.matchEventPattern(event, rule.EventPattern) {
				continue
			}
		}

		targetsResult, err := store.ListTargetsByRule(ctx, eventBusName, rule.Name, 100, "")
		if err != nil {
			continue
		}

		for _, target := range targetsResult.Targets {
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
			return false
		}

		if !s.matchValue(eventValue, patternValue) {
			return false
		}
	}

	return true
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
				continue
			}
			compVal, ok := toFloat64(operands[i+1])
			if !ok {
				continue
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
			return s == pattern
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
			return strings.HasSuffix(s, fragment)
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
	targetType := s.parseTargetType(target.ARN)

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

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logs.Error("failed to marshal event payload",
			logs.String("eventId", event.ID),
			logs.String("targetId", target.ID),
			logs.Err(err))
		return
	}

	switch targetType {
	case "lambda":
		s.deliverToLambda(ctx, region, event.ID, target.ARN, payloadBytes)
	case "sqs":
		s.deliverToSQS(ctx, region, target.ARN, payloadBytes)
	case "sns":
		s.deliverToSNS(ctx, region, target.ARN, payloadBytes)
	case "logs":
		s.deliverToCloudWatchLogs(ctx, region, event.ID, target.ARN, payloadBytes)
	case "states":
		s.deliverToStepFunctions(ctx, region, target.ARN, payloadBytes)
	case "kinesis":
		s.deliverToKinesis(ctx, region, target.ARN, payloadBytes)
	case "firehose":
		s.deliverToFirehose(ctx, region, target.ARN, payloadBytes)
	case "ecs":
		s.deliverToECS(ctx, region, target.ARN, payloadBytes)
	default:
		logs.Warn("target type not yet implemented, event not delivered",
			logs.String("targetType", targetType),
			logs.String("targetArn", target.ARN),
			logs.String("eventId", event.ID))
	}
}

func (s *EventsService) dispatchToTarget(ctx context.Context, region string, event *eventsstore.Event, target eventsstore.Target, targetType string, payload []byte) {
	switch targetType {
	case "lambda":
		s.deliverToLambda(ctx, region, event.ID, target.ARN, payload)
	case "sqs":
		s.deliverToSQS(ctx, region, target.ARN, payload)
	case "sns":
		s.deliverToSNS(ctx, region, target.ARN, payload)
	case "logs":
		s.deliverToCloudWatchLogs(ctx, region, event.ID, target.ARN, payload)
	case "states":
		s.deliverToStepFunctions(ctx, region, target.ARN, payload)
	case "kinesis":
		s.deliverToKinesis(ctx, region, target.ARN, payload)
	case "firehose":
		s.deliverToFirehose(ctx, region, target.ARN, payload)
	case "ecs":
		s.deliverToECS(ctx, region, target.ARN, payload)
	default:
		logs.Warn("target type not yet implemented, event not delivered",
			logs.String("targetType", targetType),
			logs.String("targetArn", target.ARN),
			logs.String("eventId", event.ID))
	}
}

func (s *EventsService) parseTargetType(arnStr string) string {
	_, service, _, _, _ := arnutil.SplitARN(arnStr)
	return service
}

func (s *EventsService) deliverToLambda(ctx context.Context, region string, eventID string, targetArn string, payload []byte) {
	if s.bus == nil || s.bus.LambdaInvoker() == nil {
		logs.Warn("lambda invoker not configured, skipping Lambda delivery",
			logs.String("eventId", eventID),
			logs.String("targetArn", targetArn))
		return
	}

	if s.bus != nil {
		allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, targetArn, "lambda", "events.amazonaws.com", "lambda:InvokeFunction", targetArn)
		if evalErr != nil {
			logs.Warn("resource policy evaluation failed for Lambda target, dropping delivery",
				logs.String("targetArn", targetArn),
				logs.String("error", evalErr.Error()))
			return
		}
		if !allowed {
			return
		}
	}

	functionName := arnutil.ExtractFunctionNameFromARN(targetArn)
	if functionName == "" {
		logs.Error("failed to extract function name from ARN",
			logs.String("eventId", eventID),
			logs.String("targetArn", targetArn))
		return
	}

	statusCode, result, err := s.bus.LambdaInvoker().InvokeForGateway(ctx, targetArn, payload)
	if err != nil {
		logs.Error("failed to invoke Lambda function",
			logs.String("eventId", eventID),
			logs.String("functionName", functionName),
			logs.Err(err))
		return
	}

	if statusCode != 200 {
		logs.Warn("Lambda invocation returned non-200 status",
			logs.String("eventId", eventID),
			logs.String("functionName", functionName),
			logs.Int("statusCode", int(statusCode)),
			logs.String("response", string(result)))
		return
	}

	logs.Debug("event delivered to Lambda successfully",
		logs.String("eventId", eventID),
		logs.String("functionName", functionName))
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
		template = strings.ReplaceAll(template, placeholder, fmt.Sprintf("%v", value))
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

func (s *EventsService) deliverToSQS(ctx context.Context, region string, arnStr string, payload []byte) {
	if s.bus == nil || s.bus.SQSInvoker() == nil {
		logs.Warn("SQS invoker not configured, skipping SQS delivery",
			logs.String("arn", arnStr))
		return
	}

	queueName := arnutil.ExtractQueueNameFromARN(arnStr)
	if queueName == "" {
		logs.Error("failed to extract queue name from ARN",
			logs.String("arn", arnStr))
		return
	}

	queueURL, qErr := s.bus.SQSInvoker().GetQueueByName(ctx, queueName)
	if qErr != nil {
		logs.Warn("queue not found for SQS delivery",
			logs.String("arn", arnStr),
			logs.String("queueName", queueName),
			logs.String("error", qErr.Error()))
		return
	}

	queueARN, arnErr := s.bus.SQSInvoker().GetQueueARN(ctx, queueURL)
	if arnErr != nil {
		logs.Warn("failed to get queue ARN for policy evaluation",
			logs.String("arn", arnStr),
			logs.String("error", arnErr.Error()))
		return
	}

	allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, queueARN, "sqs", "events.amazonaws.com", "sqs:SendMessage", queueARN)
	if evalErr != nil {
		logs.Warn("resource policy evaluation failed for SQS target, dropping delivery",
			logs.String("targetArn", arnStr),
			logs.String("error", evalErr.Error()))
		return
	}
	if !allowed {
		return
	}

	if _, _, err := s.bus.SQSInvoker().SendMessage(ctx, queueURL, string(payload), 0, nil); err != nil {
		logs.Error("Failed to deliver event to SQS",
			logs.String("queue", queueName),
			logs.String("error", err.Error()))
	} else {
		logs.Debug("Event delivered to SQS successfully",
			logs.String("queue", queueName))
	}
}

func (s *EventsService) deliverToSNS(ctx context.Context, region string, arnStr string, payload []byte) {
	if s.bus == nil || s.bus.SNSInvoker() == nil {
		logs.Warn("SNS invoker not configured, skipping SNS delivery",
			logs.String("arn", arnStr))
		return
	}

	allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, arnStr, "sns", "events.amazonaws.com", "sns:Publish", arnStr)
	if evalErr != nil {
		logs.Warn("resource policy evaluation failed for SNS target, dropping delivery",
			logs.String("targetArn", arnStr),
			logs.String("error", evalErr.Error()))
		return
	}
	if !allowed {
		return
	}

	_, _, _, resource, _ := arnutil.SplitARN(arnStr)
	if resource == "" {
		logs.Error("failed to extract topic name from ARN",
			logs.String("arn", arnStr))
		return
	}

	topicName := strings.TrimPrefix(resource, "topic/")

	_, err := s.bus.SNSInvoker().PublishToTopic(ctx, arnStr, string(payload), "", nil)
	if err != nil {
		logs.Error("Failed to deliver event to SNS",
			logs.String("arn", arnStr),
			logs.String("error", err.Error()))
		return
	}

	logs.Debug("Event delivered to SNS successfully",
		logs.String("topic", topicName),
		logs.String("arn", arnStr))
}

func (s *EventsService) deliverToCloudWatchLogs(ctx context.Context, region string, eventID string, arnStr string, payload []byte) {
	if s.bus == nil {
		logs.Warn("event bus not configured, skipping CloudWatch Logs delivery",
			logs.String("arn", arnStr))
		return
	}

	_, _, _, _, resource := arnutil.SplitARN(arnStr)
	logGroup := arnutil.ExtractLogGroupNameFromARN(arnStr)
	if logGroup == "" {
		logs.Error("failed to extract log group from CloudWatch Logs ARN",
			logs.String("arn", arnStr))
		return
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
		logs.Error("Failed to deliver event to CloudWatch Logs",
			logs.String("logGroup", logGroup),
			logs.String("error", err.Error()))
		return
	}

	logs.Debug("Event delivered to CloudWatch Logs successfully",
		logs.String("logGroup", logGroup),
		logs.String("logStream", logStream))
}

func (s *EventsService) deliverToStepFunctions(ctx context.Context, region string, targetArn string, payload []byte) {
	if s.bus == nil {
		logs.Warn("event bus not configured, skipping Step Functions delivery",
			logs.String("arn", targetArn))
		return
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
		logs.Error("failed to publish Step Functions start event",
			logs.String("targetArn", targetArn),
			logs.Err(err))
		return
	}

	logs.Debug("event delivered to Step Functions successfully",
		logs.String("targetArn", targetArn))
}

func (s *EventsService) deliverToKinesis(ctx context.Context, region string, targetArn string, payload []byte) {
	if s.bus == nil || s.bus.KinesisInvoker() == nil {
		logs.Warn("Kinesis invoker not configured, skipping Kinesis delivery",
			logs.String("arn", targetArn))
		return
	}

	_, _, kRegion, _, resource := arnutil.SplitARN(targetArn)
	if kRegion == "" {
		kRegion = region
	}

	streamName := resource
	if idx := strings.Index(resource, "stream/"); idx != -1 {
		streamName = resource[idx+len("stream/"):]
	}

	partitionKey := fmt.Sprintf("eventbridge-%s", generateEventID())
	encodedPayload := base64.StdEncoding.EncodeToString(payload)
	_, err := s.bus.KinesisInvoker().PutRecord(ctx, streamName, partitionKey, []byte(encodedPayload))
	if err != nil {
		logs.Error("failed to put record to Kinesis stream",
			logs.String("stream", streamName),
			logs.Err(err))
		return
	}

	logs.Debug("event delivered to Kinesis successfully",
		logs.String("stream", streamName))
}

func (s *EventsService) deliverToFirehose(ctx context.Context, region string, targetArn string, payload []byte) {
	logs.Warn("Firehose target delivery is not yet implemented",
		logs.String("targetArn", targetArn),
		logs.String("region", region))
}

func (s *EventsService) deliverToECS(ctx context.Context, region string, targetArn string, payload []byte) {
	logs.Warn("ECS target delivery is not yet implemented",
		logs.String("targetArn", targetArn),
		logs.String("region", region))
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
		return nil, awserrors.NewValidationException(fmt.Sprintf("EventPattern is not valid JSON: %s", err))
	}
	if err := json.Unmarshal([]byte(eventStr), &eventMap); err != nil {
		return nil, awserrors.NewValidationException(fmt.Sprintf("Event is not valid JSON: %s", err))
	}

	result := true
	for key, patternValue := range patternMap {
		eventValue, exists := eventMap[key]
		if !exists {
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
