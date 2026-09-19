package scheduler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// Universal-target delivery. A universal target's Arn carries the form
// arn:aws:scheduler:::aws-sdk:{service}:{action} and its Input carries the
// invoked API's request JSON ("Using universal targets in EventBridge
// Scheduler", AWS User Guide). The {service} value is the AWS SDK service
// identifier, which can differ from the endpoint prefix: sfn, not states;
// eventbridge, not events.
//
// The operations dispatched below are the ones the platform's delivery
// families implement — the same operations the templated targets deliver
// through. Every other universal ARN was accepted at validation per the AWS
// contract (the ARN form is valid for any SDK service) and fails HERE with
// a cause naming what is missing, so the schedule's retry policy and
// dead-letter routing engage exactly as they do for the ecs/firehose
// templated stubs.

// deliverUniversalTarget dispatches a universal-target delivery by its SDK
// service identifier and API action.
func (e *Engine) deliverUniversalTarget(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	sdkService, action, ok := parseUniversalTargetARN(target.Arn)
	if !ok {
		// Unreachable through validated schedules; a malformed stored
		// record still fails with a cause instead of dispatching blindly.
		return fmt.Errorf("malformed universal target ARN: %s", target.Arn)
	}

	operation := sdkService + ":" + action
	switch operation {
	case "lambda:invoke":
		return e.deliverUniversalLambda(ctx, schedule, target)
	case "sqs:sendMessage":
		return e.deliverUniversalSQS(ctx, schedule, target)
	case "sns:publish":
		return e.deliverUniversalSNS(ctx, schedule, target)
	case "kinesis:putRecord":
		return e.deliverUniversalKinesis(ctx, schedule, target)
	case "sfn:startExecution":
		return e.deliverUniversalStepFunctions(ctx, schedule, target)
	case "eventbridge:putEvents":
		return e.deliverUniversalEventBridge(ctx, schedule, target)
	}

	// Accept-and-fail: the ARN form is valid for any SDK service, so the
	// schedule was created; the delivery names what is missing.
	logs.Error("Universal target delivery failed: operation not available",
		logs.String("operation", operation),
		logs.String("targetArn", target.Arn))
	switch sdkService {
	case "sagemaker", "codebuild", "codepipeline", "inspector":
		return fmt.Errorf(
			"universal target operation %s cannot be delivered: the %s service is permanently excluded on this platform (the recorded templated-target exclusion set)",
			operation, sdkService)
	case "lambda", "sqs", "sns", "kinesis", "sfn", "eventbridge":
		return fmt.Errorf(
			"universal target operation %s is not available in this deployment: only the delivery family's primary operation is dispatched",
			operation)
	default:
		return fmt.Errorf(
			"universal target operation %s cannot be delivered: the %s service is not implemented on this platform",
			operation, sdkService)
	}
}

// scheduleDeliveryRegion resolves the region a universal request that
// carries no region of its own addresses: the schedule's region, falling
// back to the server default exactly as the engine's store resolution
// does (the bus delivery path already carries the firing region on the
// reconstructed schedule).
func scheduleDeliveryRegion(schedule *schedulerstore.Schedule) string {
	if schedule.Region != "" {
		return schedule.Region
	}
	return defaults.DefaultRegion
}

// regionFromQueueURL extracts the region an AWS-form SQS queue URL embeds
// in its host (https://sqs.<region>.amazonaws.com/...). Platform queue
// URLs are host-local and embed none; the caller then falls back to the
// schedule's region.
func regionFromQueueURL(queueURL string) string {
	u, err := url.Parse(queueURL)
	if err != nil {
		return ""
	}
	host := u.Hostname()
	if !strings.HasPrefix(host, "sqs.") {
		return ""
	}
	rest := strings.TrimSuffix(strings.TrimPrefix(host, "sqs."), ".amazonaws.com")
	if rest == "" || strings.Contains(rest, ".") {
		return ""
	}
	return rest
}

// universalRequest unmarshals the universal target's Input into the
// invoked API's request members. The AWS universal-target configuration
// specifies Input with the request parameters; an absent Input can carry
// none, so it fails with a cause instead of invoking the API empty.
func universalRequest(target *schedulerstore.Target, operation string) (map[string]json.RawMessage, error) {
	if target.Input == "" {
		return nil, fmt.Errorf("universal target %s requires Target.Input carrying the request parameters", operation)
	}
	var req map[string]json.RawMessage
	if err := json.Unmarshal([]byte(target.Input), &req); err != nil {
		return nil, fmt.Errorf("universal target %s Target.Input is not a JSON object: %w", operation, err)
	}
	return req, nil
}

// requiredRequestString extracts a required string request member. The
// error names the member the way the target API's own missing-member
// rejection would; none of the addressed APIs accept an empty value for
// these members.
func requiredRequestString(req map[string]json.RawMessage, member, operation string) (string, error) {
	raw, ok := req[member]
	if !ok {
		return "", fmt.Errorf("%s request member %s is required", operation, member)
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", fmt.Errorf("%s request member %s must be a string", operation, member)
	}
	if s == "" {
		return "", fmt.Errorf("%s request member %s must not be empty", operation, member)
	}
	return s, nil
}

// optionalRequestString extracts an optional string request member; an
// absent member or a non-string value yields "".
func optionalRequestString(req map[string]json.RawMessage, member string) string {
	raw, ok := req[member]
	if !ok {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return ""
	}
	return s
}

// optionalRequestInt64 extracts an optional integer request member; an
// absent member or a non-integer value yields 0.
func optionalRequestInt64(req map[string]json.RawMessage, member string) int64 {
	raw, ok := req[member]
	if !ok {
		return 0
	}
	var n int64
	if err := json.Unmarshal(raw, &n); err != nil {
		return 0
	}
	return n
}

// deliverUniversalLambda delivers aws-sdk:lambda:invoke: Input carries the
// Invoke request (FunctionName, optional InvocationType and Payload). The
// trigger contract of the templated family applies to RequestResponse
// invocations — a FunctionError counts against the RetryPolicy; an Event
// invocation is executed synchronously like every cross-service invoke,
// but only its transport verdict fails it (a FunctionError in Event mode
// is the async-delivery outcome AWS reports through the function's own
// logs, not a retryable delivery failure).
func (e *Engine) deliverUniversalLambda(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		return fmt.Errorf("event bus not configured")
	}
	lambdaInvoker := e.bus.LambdaInvoker()
	if lambdaInvoker == nil {
		return fmt.Errorf("Lambda invoker not available")
	}

	req, err := universalRequest(target, "lambda:invoke")
	if err != nil {
		return err
	}
	functionName, err := requiredRequestString(req, "FunctionName", "lambda:invoke")
	if err != nil {
		return err
	}

	var payload []byte
	if raw, ok := req["Payload"]; ok {
		// Invoke's Payload is a blob: the JSON request carries it as a
		// string (the documented universal-target example escapes the JSON
		// payload into one); any other JSON value passes through as its
		// own bytes.
		var s string
		if json.Unmarshal(raw, &s) == nil {
			payload = []byte(s)
		} else {
			payload = raw
		}
	}

	invocationType := optionalRequestString(req, "InvocationType")
	switch invocationType {
	case "", "RequestResponse":
		invocation, err := lambdaInvoker.InvokeForTrigger(ctx, functionName, payload)
		if err != nil {
			return err
		}
		if invocation.FunctionError != "" {
			return fmt.Errorf("lambda function %s execution failed: %s", functionName, invocation.FunctionError)
		}
		if invocation.StatusCode < 200 || invocation.StatusCode > 299 {
			return fmt.Errorf("lambda invocation of %s returned status %d", functionName, invocation.StatusCode)
		}
	case "Event":
		statusCode, _, err := lambdaInvoker.InvokeForGateway(ctx, functionName, payload)
		if err != nil {
			return err
		}
		if statusCode < 200 || statusCode > 299 {
			return fmt.Errorf("lambda invocation of %s returned status %d", functionName, statusCode)
		}
	case "DryRun":
		return fmt.Errorf("universal target lambda:invoke InvocationType DryRun is not supported on this platform")
	default:
		return fmt.Errorf("lambda:invoke request member InvocationType has an invalid value %q", invocationType)
	}

	logs.Debug("Universal lambda invocation completed",
		logs.String("schedule", schedule.Name),
		logs.String("function", functionName))
	return nil
}

// deliverUniversalSQS delivers aws-sdk:sqs:sendMessage: Input carries the
// SendMessage request (QueueUrl, MessageBody, optional DelaySeconds,
// MessageGroupId and MessageAttributes). The queue is addressed in the
// region the QueueUrl's host embeds (the AWS form) or, for host-local
// platform URLs, the schedule's region.
func (e *Engine) deliverUniversalSQS(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		return fmt.Errorf("event bus not configured")
	}
	sqsInvoker := e.bus.SQSInvoker()
	if sqsInvoker == nil {
		return fmt.Errorf("SQS invoker not available")
	}

	req, err := universalRequest(target, "sqs:sendMessage")
	if err != nil {
		return err
	}
	queueURL, err := requiredRequestString(req, "QueueUrl", "sqs:sendMessage")
	if err != nil {
		return err
	}
	messageBody, err := requiredRequestString(req, "MessageBody", "sqs:sendMessage")
	if err != nil {
		return err
	}

	sendOpts := invokers.SQSSendOptions{
		DelaySeconds:           optionalRequestInt64(req, "DelaySeconds"),
		MessageGroupID:         optionalRequestString(req, "MessageGroupId"),
		MessageDeduplicationID: optionalRequestString(req, "MessageDeduplicationId"),
	}
	if raw, ok := req["MessageAttributes"]; ok {
		attrs, err := universalSQSMessageAttributes(raw)
		if err != nil {
			return err
		}
		sendOpts.TypedMessageAttributes = attrs
	}

	region := regionFromQueueURL(queueURL)
	if region == "" {
		region = scheduleDeliveryRegion(schedule)
	}

	logs.Debug("Universal SQS sendMessage for schedule",
		logs.String("schedule", schedule.Name),
		logs.String("queueUrl", queueURL))

	if _, _, err := sqsInvoker.SendMessage(ctx, region, queueURL, messageBody, sendOpts); err != nil {
		return err
	}
	return nil
}

// universalSQSMessageAttributes translates the SendMessage wire form of
// MessageAttributes (name → {DataType, StringValue | BinaryValue}) into the
// typed invoker attribute form; the Binary value is base64 on the wire.
func universalSQSMessageAttributes(raw json.RawMessage) (map[string]invokers.SQSMessageAttribute, error) {
	var wire map[string]struct {
		DataType    string `json:"DataType"`
		StringValue string `json:"StringValue"`
		BinaryValue string `json:"BinaryValue"`
	}
	if err := json.Unmarshal(raw, &wire); err != nil {
		return nil, fmt.Errorf("sqs:sendMessage request member MessageAttributes must be a map of message attributes: %w", err)
	}
	attrs := make(map[string]invokers.SQSMessageAttribute, len(wire))
	for name, attr := range wire {
		switch attr.DataType {
		case "String", "Number":
			attrs[name] = invokers.SQSMessageAttribute{DataType: attr.DataType, StringValue: attr.StringValue}
		case "Binary":
			decoded, err := base64.StdEncoding.DecodeString(attr.BinaryValue)
			if err != nil {
				return nil, fmt.Errorf("sqs:sendMessage MessageAttributes[%s].BinaryValue must be base64-encoded: %w", name, err)
			}
			attrs[name] = invokers.SQSMessageAttribute{DataType: attr.DataType, BinaryValue: decoded}
		default:
			return nil, fmt.Errorf("sqs:sendMessage MessageAttributes[%s].DataType %q is not valid", name, attr.DataType)
		}
	}
	return attrs, nil
}

// deliverUniversalSNS delivers aws-sdk:sns:publish: Input carries the
// Publish request (TopicArn, Message, optional Subject and
// MessageAttributes).
func (e *Engine) deliverUniversalSNS(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		return fmt.Errorf("event bus not configured")
	}
	snsInvoker := e.bus.SNSInvoker()
	if snsInvoker == nil {
		return fmt.Errorf("SNS invoker not available")
	}

	req, err := universalRequest(target, "sns:publish")
	if err != nil {
		return err
	}
	topicArn, err := requiredRequestString(req, "TopicArn", "sns:publish")
	if err != nil {
		return err
	}
	message, err := requiredRequestString(req, "Message", "sns:publish")
	if err != nil {
		return err
	}

	var attributes map[string]string
	if raw, ok := req["MessageAttributes"]; ok {
		attributes, err = universalSNSMessageAttributes(raw)
		if err != nil {
			return err
		}
	}

	messageID, err := snsInvoker.PublishToTopic(ctx, topicArn, message, optionalRequestString(req, "Subject"), attributes)
	if err != nil {
		return err
	}
	logs.Debug("Universal SNS publish completed",
		logs.String("schedule", schedule.Name),
		logs.String("topicArn", topicArn),
		logs.String("messageId", messageID))
	return nil
}

// universalSNSMessageAttributes translates the Publish wire form of
// MessageAttributes for the string-valued invoker seam: String and Number
// types translate, a Binary attribute fails the delivery with a cause
// rather than being dropped silently (the seam carries no binary form).
func universalSNSMessageAttributes(raw json.RawMessage) (map[string]string, error) {
	var wire map[string]struct {
		DataType    string `json:"DataType"`
		StringValue string `json:"StringValue"`
	}
	if err := json.Unmarshal(raw, &wire); err != nil {
		return nil, fmt.Errorf("sns:publish request member MessageAttributes must be a map of message attributes: %w", err)
	}
	attrs := make(map[string]string, len(wire))
	for name, attr := range wire {
		switch attr.DataType {
		case "String", "Number":
			attrs[name] = attr.StringValue
		case "Binary":
			return nil, fmt.Errorf("sns:publish MessageAttributes[%s] uses the Binary data type, which this platform's publish path does not carry", name)
		default:
			return nil, fmt.Errorf("sns:publish MessageAttributes[%s].DataType %q is not valid", name, attr.DataType)
		}
	}
	return attrs, nil
}

// deliverUniversalKinesis delivers aws-sdk:kinesis:putRecord: Input
// carries the PutRecord request (Data as base64, PartitionKey, and
// StreamName or StreamARN). The stream is addressed in the StreamARN's
// region when one is given, otherwise the schedule's region. Data passes
// through in its base64 wire form: the local Kinesis service stores the
// data string as-is and GetRecords returns it without additional
// encoding, and SDK clients expect base64-encoded Data in that response —
// the same pre-encode convention the cross-service EventBridge dispatcher
// applies.
func (e *Engine) deliverUniversalKinesis(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		return fmt.Errorf("event bus not configured")
	}
	kinesisInvoker := e.bus.KinesisInvoker()
	if kinesisInvoker == nil {
		return fmt.Errorf("Kinesis invoker not available")
	}

	req, err := universalRequest(target, "kinesis:putRecord")
	if err != nil {
		return err
	}
	dataB64, err := requiredRequestString(req, "Data", "kinesis:putRecord")
	if err != nil {
		return err
	}
	partitionKey, err := requiredRequestString(req, "PartitionKey", "kinesis:putRecord")
	if err != nil {
		return err
	}

	region := scheduleDeliveryRegion(schedule)
	streamName := optionalRequestString(req, "StreamName")
	if streamName == "" {
		streamARN := optionalRequestString(req, "StreamARN")
		if streamARN == "" {
			return fmt.Errorf("kinesis:putRecord requires request member StreamName or StreamARN")
		}
		_, _, arnRegion, _, _ := svcarn.SplitARN(streamARN)
		if arnRegion != "" {
			region = arnRegion
		}
		streamName = svcarn.ExtractStreamNameFromARN(streamARN)
	}

	logs.Debug("Universal Kinesis putRecord for schedule",
		logs.String("schedule", schedule.Name),
		logs.String("stream", streamName))

	if _, _, err := kinesisInvoker.PutRecord(ctx, region, streamName, partitionKey, []byte(dataB64)); err != nil {
		return err
	}
	return nil
}

// deliverUniversalStepFunctions delivers aws-sdk:sfn:startExecution:
// Input carries the StartExecution request (StateMachineArn, optional
// Input). The execution Name member is not carried — the platform's
// execution-start seam generates the execution name itself.
func (e *Engine) deliverUniversalStepFunctions(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		return fmt.Errorf("event bus not configured")
	}

	req, err := universalRequest(target, "sfn:startExecution")
	if err != nil {
		return err
	}
	stateMachineArn, err := requiredRequestString(req, "StateMachineArn", "sfn:startExecution")
	if err != nil {
		return err
	}

	_, _, smRegion, _, _ := svcarn.SplitARN(stateMachineArn)
	if smRegion == "" {
		smRegion = scheduleDeliveryRegion(schedule)
	}

	// StartExecution's Input member is optional and defaults to the empty
	// JSON object, which is also the faithful encoding of an omitted one.
	input := optionalRequestString(req, "Input")
	if input == "" {
		input = "{}"
	}

	evt := &eventbus.StepFunctionsStartExecutionEvent{
		StateMachineArn: stateMachineArn,
		Input:           input,
	}
	evt.Region = smRegion
	evt.AccountID = e.accountID

	// PublishSync so the execution-start failure propagates back and the
	// schedule's retry policy and dead-letter routing apply to it. The bus
	// contract reports a handler failure in HandlerResult.Error with a nil
	// error, so both fold into the returned error.
	result, err := e.bus.PublishSync(ctx, evt)
	if err == nil {
		err = result.Error
	}
	if err != nil {
		return err
	}

	logs.Debug("Universal Step Functions execution started",
		logs.String("schedule", schedule.Name),
		logs.String("stateMachineArn", stateMachineArn))
	return nil
}

// deliverUniversalEventBridge delivers aws-sdk:eventbridge:putEvents:
// Input carries the PutEvents request (Entries). Each entry is delivered
// as its own putEvents bus event carrying the entry's members — Source,
// DetailType, Detail, EventBusName, and the Resources/Time members the
// handler applies with the PutEvents semantics. The entry-count bound of
// the PutEvents API plane is that plane's request rule and is not
// re-imposed here. Any entry's failure fails the delivery, so the retry
// policy re-drives the whole request (at-least-once, the same posture AWS
// applies to a retried partially-failed PutEvents).
func (e *Engine) deliverUniversalEventBridge(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		return fmt.Errorf("event bus not configured")
	}

	req, err := universalRequest(target, "eventbridge:putEvents")
	if err != nil {
		return err
	}
	rawEntries, ok := req["Entries"]
	if !ok {
		return fmt.Errorf("eventbridge:putEvents request member Entries is required")
	}
	var entries []map[string]json.RawMessage
	if err := json.Unmarshal(rawEntries, &entries); err != nil {
		return fmt.Errorf("eventbridge:putEvents request member Entries must be an array of entry objects: %w", err)
	}
	if len(entries) == 0 {
		return fmt.Errorf("eventbridge:putEvents request member Entries must carry at least one entry")
	}

	for _, entry := range entries {
		entryBytes, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("eventbridge:putEvents entry re-encoding failed: %w", err)
		}

		region := ""
		if busName := optionalRequestString(entry, "EventBusName"); busName != "" {
			// EventBusName accepts the name or ARN form; an ARN carries
			// the region the bus's store lives in.
			_, service, busRegion, _, _ := svcarn.SplitARN(busName)
			if service == "events" {
				region = busRegion
			}
		}
		if region == "" {
			region = scheduleDeliveryRegion(schedule)
		}

		evt := &eventbus.EventBridgePutEventsEvent{
			EventBusName: optionalRequestString(entry, "EventBusName"),
			Input:        string(entryBytes),
			DetailType:   optionalRequestString(entry, "DetailType"),
			Source:       optionalRequestString(entry, "Source"),
		}
		evt.Region = region
		evt.AccountID = e.accountID

		result, err := e.bus.PublishSync(ctx, evt)
		if err == nil {
			err = result.Error
		}
		if err != nil {
			return err
		}
	}

	logs.Debug("Universal EventBridge putEvents completed",
		logs.String("schedule", schedule.Name),
		logs.Int("entries", len(entries)))
	return nil
}
