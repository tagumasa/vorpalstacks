package eventbridge

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

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

// retriesExhausted reports whether the retry budget is spent after the
// given failed attempt. MaximumRetryAttempts counts retries after the
// initial attempt, so the total attempt budget is maxRetries+1: zero
// permits a single attempt with no retries.
func retriesExhausted(attempt, maxRetries int32) bool {
	return attempt > maxRetries
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
			deliverErr = s.deliverToKinesis(ctx, region, event.ID, target, payload)
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
		if retriesExhausted(attempt, maxRetries) || time.Now().After(deadline) {
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
		if err := s.deliverToSQS(ctx, region, eventsstore.Target{ARN: dlqArn}, payload); err != nil {
			// The event is now lost for this target: the primary delivery
			// already failed and the DLQ copy failed too. Log loudly so
			// the loss is visible instead of silently reported as routed.
			logs.Error("failed to route event to DLQ (SQS)",
				logs.String("dlqArn", dlqArn),
				logs.String("eventId", event.ID),
				logs.String("originalTarget", target.ARN),
				logs.Err(err))
			return
		}
		logs.Info("event routed to DLQ (SQS)",
			logs.String("dlqArn", dlqArn),
			logs.String("eventId", event.ID),
			logs.String("originalTarget", target.ARN))
	case "sns":
		if err := s.deliverToSNS(ctx, region, dlqArn, payload); err != nil {
			logs.Error("failed to route event to DLQ (SNS)",
				logs.String("dlqArn", dlqArn),
				logs.String("eventId", event.ID),
				logs.String("originalTarget", target.ARN),
				logs.Err(err))
			return
		}
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

	_, _, sqsRegion, _, _ := arnutil.SplitARN(arnStr)

	queueURL, qErr := s.bus.SQSInvoker().GetQueueByName(ctx, sqsRegion, queueName)
	if qErr != nil {
		return fmt.Errorf("queue not found for SQS delivery %s: %w", queueName, qErr)
	}

	queueARN, arnErr := s.bus.SQSInvoker().GetQueueARN(ctx, sqsRegion, queueURL)
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

	if _, _, err := s.bus.SQSInvoker().SendMessage(ctx, sqsRegion, queueURL, string(payload), opts); err != nil {
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

// deliverToKinesis puts the event to a Kinesis stream. The partition key
// defaults to the event ID — the behaviour documented for Kinesis targets
// in the EventBridge API reference (KinesisParameters) — and stays stable
// across delivery retries, so a retried event lands on the same shard
// instead of scattering a single logical event across the stream.
func (s *EventsService) deliverToKinesis(ctx context.Context, region string, eventID string, target eventsstore.Target, payload []byte) error {
	if s.bus == nil || s.bus.KinesisInvoker() == nil {
		return fmt.Errorf("Kinesis invoker not configured")
	}

	targetArn := target.ARN
	_, _, _, _, resource := arnutil.SplitARN(targetArn)

	streamName := resource
	if idx := strings.Index(resource, "stream/"); idx != -1 {
		streamName = resource[idx+len("stream/"):]
	}

	partitionKey := eventID
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
