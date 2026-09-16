package eventbridge

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// deliverToTarget builds the target's payload and runs the synchronous
// delivery lifecycle, returning dispatchToTarget's terminal error: nil on
// success, otherwise the last delivery error once the retry budget is spent
// (or the failure is permanent) and dead-letter routing has been attempted.
func (s *EventsService) deliverToTarget(ctx context.Context, region, ruleARN, ruleName string, event *eventsstore.Event, target eventsstore.Target) error {
	payloadBytes := s.buildTargetPayload(ruleARN, ruleName, event, target)
	return s.dispatchToTarget(ctx, region, ruleARN, event, target, payloadBytes)
}

const (
	// AWS production defaults: the retry ceiling and the 24 h deadline
	// coincide with the RetryPolicy bounds ("By default, EventBridge
	// retries sending the event for 24 hours and up to 185 times with an
	// exponential back off and jitter", EventBridge user guide), with a
	// 1 s initial backoff growing to a 60 s cap.
	prodMaxRetryAttempts     = eventsstore.RetryPolicyMaxRetryAttempts
	prodMaxEventAgeInSeconds = eventsstore.RetryPolicyMaxEventAgeSeconds
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

// newDeliveryJob builds the retry state for one target delivery from the
// target's explicit RetryPolicy (validated at PutTargets time; out-of-range
// values keep the deployment defaults) or the production defaults.
func (s *EventsService) newDeliveryJob(region string, event *eventsstore.Event, target eventsstore.Target, payload []byte) *deliveryJob {
	maxRetries, maxAge, initialBackoff, maxBackoff := retryDefaults()
	if target.RetryPolicy != nil {
		if v := target.RetryPolicy.MaximumRetryAttempts; v >= 0 && v <= eventsstore.RetryPolicyMaxRetryAttempts {
			maxRetries = v
		}
		if v := target.RetryPolicy.MaximumEventAgeInSeconds; v >= eventsstore.RetryPolicyMinEventAgeSeconds && v <= eventsstore.RetryPolicyMaxEventAgeSeconds {
			maxAge = v
		}
	}
	job := &deliveryJob{
		region:     region,
		event:      event,
		target:     target,
		payload:    payload,
		maxRetries: maxRetries,
		deadline:   time.Now().Add(time.Duration(maxAge) * time.Second),
		backoff:    initialBackoff,
		maxBackoff: maxBackoff,
	}
	// The partition-key path references the original event, which this
	// constructor's callers hold in full (the bus handler overrides the
	// field with the publisher's resolution for the stub event it carries).
	if key := kinesisPartitionKeyFor(event, &target); key != "" {
		job.partitionKey = key
	}
	return job
}

// kinesisPartitionKeyFor resolves the target's Kinesis partition-key path
// against the ORIGINAL event — "dynamic path parameters must reference the
// original event, not the transformed event" — so the direct-dispatch
// constructor and the bus-delivery publisher resolve it identically. An
// unresolved or empty key returns "" so the caller leaves the field unset.
func kinesisPartitionKeyFor(event *eventsstore.Event, target *eventsstore.Target) string {
	if target.KinesisParameters == nil || target.KinesisParameters.PartitionKeyPath == "" {
		return ""
	}
	val, ok := resolveEventPathString(eventEnvelope(event), target.KinesisParameters.PartitionKeyPath)
	if !ok {
		return ""
	}
	valStr, ok := val.(string)
	if !ok || valStr == "" {
		return ""
	}
	return valStr
}

// attemptDelivery runs a single delivery attempt and returns its outcome.
// Failures classified permanent (errPermanentDelivery) are never retried;
// every other failure is transient and consumes the retry budget.
func (s *EventsService) attemptDelivery(ctx context.Context, job *deliveryJob) error {
	targetType := s.parseTargetType(job.target.ARN)
	switch targetType {
	case "lambda":
		return s.deliverToLambda(ctx, job.region, job.event.ID, job.target.ARN, job.payload)
	case "sqs":
		return s.deliverToSQS(ctx, job.region, job.target, job.payload, nil)
	case "sns":
		return s.deliverToSNS(ctx, job.region, job.target.ARN, job.payload)
	case "logs":
		return s.deliverToCloudWatchLogs(ctx, job.region, job.event.ID, job.target.ARN, job.payload)
	case "states":
		return s.deliverToStepFunctions(ctx, job.region, job.target.ARN, job.payload)
	case "kinesis":
		return s.deliverToKinesis(ctx, job.region, job.event.ID, job.partitionKey, job.target, job.payload)
	case "appsync":
		return s.deliverToAppSync(ctx, job.region, job.event.ID, job.target, job.payload)
	case "firehose":
		return s.deliverToFirehose(ctx, job.region, job.target.ARN, job.payload)
	case "ecs":
		return s.deliverToECS(ctx, job.region, job.target.ARN, job.payload)
	case "events":
		_, _, _, _, resource := arnutil.SplitARN(job.target.ARN)
		switch {
		case strings.HasPrefix(resource, "api-destination/"):
			return s.deliverToApiDestination(ctx, job)
		case strings.HasPrefix(resource, "event-bus/"):
			return s.deliverToEventBus(ctx, job.region, job.event, job.target.ARN)
		default:
			// Every other events-service resource form (rule, archive, a
			// bare name) would otherwise deliver to a phantom bus that no
			// rule listing ever matches — a silent drop reported as success.
			return permanentDeliveryError("unsupported events-service target resource %q: only event-bus/ and api-destination/ targets have a delivery path", job.target.ARN)
		}
	default:
		return permanentDeliveryError("target type %q not implemented", targetType)
	}
}

// terminalDelivery handles a delivery whose failure is permanent or whose
// retry budget is spent. AWS drops the exhausted event ("If an event isn't
// delivered after all retry attempts are exhausted, the event is dropped
// and EventBridge doesn't continue to process it", user guide); a
// configured dead-letter queue receives it first. The returned error is
// non-nil only when a configured DLQ write itself failed, so the caller
// can report the loss instead of recording a completed delivery.
func (s *EventsService) terminalDelivery(ctx context.Context, job *deliveryJob, deliverErr error) error {
	logs.Error("event delivery to target failed after retries",
		logs.String("targetArn", job.target.ARN),
		logs.String("eventId", job.event.ID),
		logs.Int("attempts", int(job.attempts)),
		logs.Err(deliverErr))
	return s.routeToDeadLetter(ctx, job, deliverErr)
}

// dispatchToTarget runs the full delivery lifecycle synchronously — the
// loop form of the retry engine, used by the non-bus delivery arm and
// direct callers. It returns the terminal error: nil on success, otherwise
// the last delivery error once the retry budget is spent (or the failure
// is permanent) and dead-letter routing has been attempted. A caller
// context cancelled mid-retry dead-letters the same way, through a
// detached bounded context so the cancellation that ends the retry cannot
// also kill the dead-letter write.
func (s *EventsService) dispatchToTarget(ctx context.Context, region, ruleARN string, event *eventsstore.Event, target eventsstore.Target, payload []byte) error {
	job := s.newDeliveryJob(region, event, target, payload)
	job.ruleARN = ruleARN
	for {
		job.attempts++
		err := s.attemptDelivery(ctx, job)
		if err == nil {
			return nil
		}
		if errors.Is(err, errPermanentDelivery) || job.exhausted() {
			if dlqErr := s.terminalDelivery(ctx, job, err); dlqErr != nil {
				return dlqErr
			}
			return err
		}
		select {
		case <-time.After(backoffSleep(job.backoff, job.maxBackoff)):
		case <-ctx.Done():
			// The cancelled delivery gets the same terminal handling the
			// engine path gives; the write runs on its own short-lived
			// context because the caller's is gone.
			dlqCtx, dlqCancel := context.WithTimeout(context.Background(), dlqShutdownWriteTimeout)
			_ = s.terminalDelivery(dlqCtx, job, ctx.Err())
			dlqCancel()
			return ctx.Err()
		}
		job.backoff *= 2
	}
}

// dlqShutdownWriteTimeout bounds the detached dead-letter write a
// cancelled synchronous dispatch performs before returning the context
// error.
const dlqShutdownWriteTimeout = 5 * time.Second

// dlqWriteAttempts bounds the dead-letter SendMessage retries: the DLQ
// copy is the last durable form of an event whose primary budget is
// already spent, so a transient queue fault gets a short bounded retry
// rather than a single shot.
const dlqWriteAttempts = 3

// routeToDeadLetter delivers the event payload to the configured SQS
// dead-letter queue when the primary target delivery fails terminally.
// The message carries the attribute envelope AWS documents for EventBridge
// DLQs — RULE_ARN, TARGET_ARN, ERROR_CODE, ERROR_MESSAGE,
// EXHAUSTED_RETRY_CONDITION (MaximumRetryAttempts or
// MaximumEventAgeInSeconds, present only when a retry condition ended the
// delivery) and RETRY_ATTEMPTS — plus the trace header as the
// AWSTraceHeader message attribute when the inbound event carried one.
// It returns nil when no DLQ is configured (AWS drops the exhausted
// event) or the write succeeded; a non-nil error means the event is lost
// for this target and the caller must report the loss rather than record
// a completed delivery.
func (s *EventsService) routeToDeadLetter(ctx context.Context, job *deliveryJob, deliverErr error) error {
	if job.target.DeadLetterConfig == nil || job.target.DeadLetterConfig.Arn == "" {
		return nil
	}
	dlqArn := job.target.DeadLetterConfig.Arn
	if dlqType := s.parseTargetType(dlqArn); dlqType != "sqs" {
		return fmt.Errorf("dead-letter queue %s: EventBridge DLQs are SQS queues, got service %q", dlqArn, dlqType)
	}

	attrs := s.deadLetterAttributes(job, deliverErr)
	dlqTarget := eventsstore.Target{ARN: dlqArn}
	var lastErr error
	for attempt := 0; attempt < dlqWriteAttempts; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(backoffSleep(time.Duration(attempt)*100*time.Millisecond, 400*time.Millisecond)):
			case <-ctx.Done():
				return fmt.Errorf("dead-letter write for event %s cancelled: %w", job.event.ID, ctx.Err())
			}
		}
		if err := s.deliverToSQS(ctx, job.region, dlqTarget, job.payload, attrs); err != nil {
			lastErr = fmt.Errorf("failed to route event %s to DLQ %s (attempt %d): %w", job.event.ID, dlqArn, attempt+1, err)
			continue
		}
		logs.Info("event routed to DLQ (SQS)",
			logs.String("dlqArn", dlqArn),
			logs.String("eventId", job.event.ID),
			logs.String("originalTarget", job.target.ARN))
		return nil
	}
	return lastErr
}

// deadLetterAttributes builds the documented DLQ message-attribute
// envelope. ERROR_CODE uses the AWS DLQ error-code vocabulary; the
// EXHAUSTED_RETRY_CONDITION attribute is emitted only when a retry
// condition (attempt budget or age deadline) ended the delivery —
// permanent failures and permission denials go to the DLQ without
// retries, so neither condition applies.
func (s *EventsService) deadLetterAttributes(job *deliveryJob, deliverErr error) map[string]invokers.SQSMessageAttribute {
	str := func(v string) invokers.SQSMessageAttribute {
		return invokers.SQSMessageAttribute{DataType: "String", StringValue: v}
	}
	errorCode := "ERROR_FROM_TARGET"
	if errors.Is(deliverErr, errPermanentDelivery) {
		errorCode = "NO_RESOURCE"
	} else if strings.Contains(deliverErr.Error(), "resource policy denied") {
		errorCode = "NO_PERMISSIONS"
	}

	attrs := map[string]invokers.SQSMessageAttribute{
		"RULE_ARN":       str(job.ruleARN),
		"TARGET_ARN":     str(job.target.ARN),
		"ERROR_CODE":     str(errorCode),
		"ERROR_MESSAGE":  str(deliverErr.Error()),
		"RETRY_ATTEMPTS": str(strconv.FormatInt(int64(job.attempts-1), 10)),
	}
	if job.attempts > job.maxRetries {
		attrs["EXHAUSTED_RETRY_CONDITION"] = str("MaximumRetryAttempts")
	} else if time.Now().After(job.deadline) {
		attrs["EXHAUSTED_RETRY_CONDITION"] = str("MaximumEventAgeInSeconds")
	}
	if job.traceHeader != "" {
		attrs["AWSTraceHeader"] = str(job.traceHeader)
	}
	return attrs
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

// deliverToSQS sends the payload to the queue named by the target ARN.
// extraAttributes carries the dead-letter envelope when this write is a
// DLQ routing; regular target delivery passes nil.
func (s *EventsService) deliverToSQS(ctx context.Context, region string, target eventsstore.Target, payload []byte, extraAttributes map[string]invokers.SQSMessageAttribute) error {
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

	opts := invokers.SQSSendOptions{}
	if target.SqsParameters != nil && target.SqsParameters.MessageGroupId != "" {
		opts.MessageGroupID = target.SqsParameters.MessageGroupId
	}
	if extraAttributes != nil {
		opts.TypedMessageAttributes = extraAttributes
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
// instead of scattering a single logical event across the stream. The
// explicit key arrives already resolved from the original event ("dynamic
// path parameters must reference the original event, not the transformed
// event").
func (s *EventsService) deliverToKinesis(ctx context.Context, region, eventID, partitionKey string, target eventsstore.Target, payload []byte) error {
	if s.bus == nil || s.bus.KinesisInvoker() == nil {
		return fmt.Errorf("Kinesis invoker not configured")
	}

	targetArn := target.ARN
	_, _, _, _, resource := arnutil.SplitARN(targetArn)

	streamName := resource
	if idx := strings.Index(resource, "stream/"); idx != -1 {
		streamName = resource[idx+len("stream/"):]
	}

	if partitionKey == "" {
		partitionKey = eventID
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

// deliverToAppSync invokes the GraphQL mutation configured for the target.
// The ARN names the GraphQL endpoint of the API
// (arn:...:appsync:region:account:apis/<apiId>[/endpoints/GRAPHQL]);
// AppSyncParameters.GraphQLOperation carries the mutation document and the
// transformed event payload becomes the operation's variables — the shapes
// documented on the AWS AppSync target page.
func (s *EventsService) deliverToAppSync(ctx context.Context, region string, eventID string, target eventsstore.Target, payload []byte) error {
	operation := ""
	if target.AppSyncParameters != nil {
		operation = target.AppSyncParameters.GraphQLOperation
	}
	if operation == "" {
		return permanentDeliveryError("AppSync target %s carries no GraphQLOperation", target.ARN)
	}
	apiID := extractAppSyncApiIDFromARN(target.ARN)
	if apiID == "" {
		return permanentDeliveryError("AppSync target ARN %s does not name a GraphQL API endpoint (apis/<apiId>)", target.ARN)
	}
	if s.bus == nil || s.bus.AppSyncInvoker() == nil {
		return fmt.Errorf("AppSync invoker not configured")
	}

	allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, target.ARN, "appsync", "events.amazonaws.com", "appsync:GraphQL", target.ARN)
	if evalErr != nil {
		return fmt.Errorf("resource policy evaluation failed for AppSync target: %w", evalErr)
	}
	if !allowed {
		return fmt.Errorf("resource policy denied AppSync GraphQL invocation")
	}

	_, _, arnRegion, _, _ := arnutil.SplitARN(target.ARN)
	if arnRegion == "" {
		arnRegion = region
	}

	if err := s.bus.AppSyncInvoker().ExecuteGraphQLMutation(ctx, arnRegion, apiID, operation, payload); err != nil {
		return fmt.Errorf("failed to invoke AppSync mutation on %s: %w", target.ARN, err)
	}

	logs.Debug("event delivered to AppSync successfully",
		logs.String("eventId", eventID),
		logs.String("apiId", apiID))
	return nil
}

// extractAppSyncApiIDFromARN pulls the API ID out of an AppSync endpoint
// ARN resource: apis/<apiId> or apis/<apiId>/endpoints/GRAPHQL.
func extractAppSyncApiIDFromARN(arnStr string) string {
	_, _, _, _, resource := arnutil.SplitARN(arnStr)
	rest, ok := strings.CutPrefix(resource, "apis/")
	if !ok {
		return ""
	}
	if idx := strings.Index(rest, "/"); idx != -1 {
		rest = rest[:idx]
	}
	return rest
}

// deliverToFirehose fails permanently until the Firehose service exists on
// this platform (release blocker 4): the target ARN stays accepted, the
// delivery fails fast to its terminal handling instead of burning the
// retry budget on a service that cannot answer.
func (s *EventsService) deliverToFirehose(ctx context.Context, region string, targetArn string, payload []byte) error {
	logs.Error("Firehose target delivery failed: Firehose service is not available",
		logs.String("targetArn", targetArn),
		logs.String("region", region))
	return permanentDeliveryError("firehose delivery target %s is not available in this deployment", targetArn)
}

// deliverToECS fails permanently: the ECS service is out of scope for this
// platform (Basic Policy future-expansion list), so a stored ECS target
// terminates immediately instead of burning the retry budget.
func (s *EventsService) deliverToECS(ctx context.Context, region string, targetArn string, payload []byte) error {
	logs.Error("ECS target delivery failed: ECS service is not available",
		logs.String("targetArn", targetArn),
		logs.String("region", region))
	return permanentDeliveryError("ecs delivery target %s is not available in this deployment", targetArn)
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

	// The unified ingress: archive onto the target bus's enabled archives,
	// then rule matching with the carried hop depth.
	return s.deliverEvent(childCtx, store, event, busName, targetRegion, depth+1)
}
