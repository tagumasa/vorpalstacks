package scheduler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// Target delivery: dispatch by ARN service and the per-family deliverers.

// deliverToTarget dispatches a schedule delivery to the appropriate target
// type and returns an error on failure. This is the single entry point for
// all target deliveries, used by both the direct path and the bus path.
func (e *Engine) deliverToTarget(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	_, service, _, _, _ := svcarn.SplitARN(target.Arn)
	// The universal-target ARN form (arn:aws:scheduler:::aws-sdk:{svc}:{action})
	// dispatches by its SDK service identifier and API action instead of the
	// templated-service switch below.
	if service == universalTargetARNService {
		return e.deliverUniversalTarget(ctx, schedule, target)
	}
	switch service {
	case "lambda":
		return e.invokeLambda(ctx, schedule, target)
	case "sqs":
		return e.sendToSQS(ctx, schedule, target)
	case "sns":
		return e.publishToSNS(ctx, schedule, target)
	case "kinesis":
		return e.sendToKinesis(ctx, schedule, target)
	case "states":
		return e.startStepFunctionExecution(ctx, schedule, target)
	case "events":
		return e.sendToEventBridge(ctx, schedule, target)
	case "ecs":
		// ECS is an AWS templated target. The ECS service is not yet
		// available on this platform, so delivery fails and the schedule
		// engine's retry/DLQ path takes over. EcsParameters validation
		// is fully implemented, so the schedule itself is valid.
		logs.Error("ECS target delivery failed: ECS service is not available",
			logs.String("targetArn", target.Arn))
		return fmt.Errorf("ecs delivery target %s is not available in this deployment", target.Arn)
	case "firehose":
		// Firehose is an AWS templated target with no sub-parameters.
		// The Firehose service is not yet available on this platform,
		// so delivery fails and the retry/DLQ path takes over.
		logs.Error("Firehose target delivery failed: Firehose service is not available",
			logs.String("targetArn", target.Arn))
		return fmt.Errorf("firehose delivery target %s is not available in this deployment", target.Arn)
	default:
		return fmt.Errorf("unsupported target type: %s", target.Arn)
	}
}

func scheduleInput(target *schedulerstore.Target, scheduleName string) string {
	if target.Input != "" {
		return target.Input
	}
	msgPayload := map[string]interface{}{
		"schedule":  scheduleName,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	if msgBytes, err := json.Marshal(msgPayload); err == nil {
		return string(msgBytes)
	}
	return "{}"
}

// sqsFifoSendOptions resolves the send options of an SQS delivery: the
// target's SqsParameters.MessageGroupId when set, otherwise the schedule
// name when the destination queue is FIFO — AWS requires a MessageGroupId
// on every FIFO SendMessage, and the schedule name keeps each schedule's
// deliveries inside one ordered group. Both SQS writers (the target
// deliverer and the DLQ router) resolve through here so the fallback rule
// cannot drift between them.
func sqsFifoSendOptions(target *schedulerstore.Target, queueName, scheduleName string) invokers.SQSSendOptions {
	sendOpts := invokers.SQSSendOptions{}
	if target.SqsParameters != nil && target.SqsParameters.MessageGroupId != "" {
		sendOpts.MessageGroupID = target.SqsParameters.MessageGroupId
	} else if strings.HasSuffix(queueName, ".fifo") {
		sendOpts.MessageGroupID = scheduleName
	}
	return sendOpts
}

func (e *Engine) invokeLambda(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured for Lambda invocation", logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	input := scheduleInput(target, schedule.Name)

	functionName := svcarn.ExtractFunctionNameFromARN(target.Arn)
	if functionName == "" {
		logs.Debug("Failed to extract function name from ARN",
			logs.String("schedule", schedule.Name),
			logs.String("arn", target.Arn))
		return fmt.Errorf("invalid Lambda ARN: %s", target.Arn)
	}

	logs.Debug("Invoking Lambda for schedule",
		logs.String("schedule", schedule.Name),
		logs.String("function", functionName))

	lambdaInvoker := e.bus.LambdaInvoker()
	if lambdaInvoker == nil {
		return fmt.Errorf("Lambda invoker not available")
	}

	// The trigger invocation distinguishes a failed function execution
	// (FunctionError, with the invoke transport succeeding) from a
	// transport failure. AWS counts a failed invocation against the
	// schedule's RetryPolicy, so both outcomes must reach the retry/DLQ
	// lifecycle as delivery errors — never as a silent success.
	invocation, err := lambdaInvoker.InvokeForTrigger(ctx, target.Arn, []byte(input))
	if err != nil {
		logs.Debug("Failed to invoke Lambda",
			logs.String("schedule", schedule.Name),
			logs.String("function", functionName),
			logs.String("error", err.Error()))
		return err
	}
	if invocation.FunctionError != "" {
		logs.Debug("Lambda function execution failed",
			logs.String("schedule", schedule.Name),
			logs.String("function", functionName),
			logs.String("functionError", invocation.FunctionError))
		return fmt.Errorf("lambda function %s execution failed: %s", functionName, invocation.FunctionError)
	}
	if invocation.StatusCode < 200 || invocation.StatusCode > 299 {
		logs.Debug("Lambda invocation returned a non-success status",
			logs.String("schedule", schedule.Name),
			logs.String("function", functionName),
			logs.Int("statusCode", int(invocation.StatusCode)))
		return fmt.Errorf("lambda invocation of %s returned status %d", functionName, invocation.StatusCode)
	}

	logs.Debug("Lambda invocation completed",
		logs.String("schedule", schedule.Name),
		logs.String("function", functionName),
		logs.Int("statusCode", int(invocation.StatusCode)))
	return nil
}

func (e *Engine) sendToSQS(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured for SQS delivery", logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	sqsInvoker := e.bus.SQSInvoker()
	if sqsInvoker == nil {
		return fmt.Errorf("SQS invoker not available")
	}

	queueName := svcarn.ExtractQueueNameFromARN(target.Arn)
	if queueName == "" {
		logs.Debug("Invalid SQS ARN", logs.String("arn", target.Arn))
		return fmt.Errorf("invalid SQS ARN: %s", target.Arn)
	}

	_, _, sqsRegion, _, _ := svcarn.SplitARN(target.Arn)

	queueURL, qErr := sqsInvoker.GetQueueByName(ctx, sqsRegion, queueName)
	if qErr != nil {
		logs.Debug("SQS queue not found", logs.String("queue", queueName), logs.Err(qErr))
		return qErr
	}

	messageBody := scheduleInput(target, schedule.Name)

	sendOpts := sqsFifoSendOptions(target, queueName, schedule.Name)

	logs.Debug("Sending to SQS for schedule",
		logs.String("schedule", schedule.Name),
		logs.String("queue", queueName))

	if _, _, err := sqsInvoker.SendMessage(ctx, sqsRegion, queueURL, messageBody, sendOpts); err != nil {
		logs.Debug("Failed to send to SQS",
			logs.String("schedule", schedule.Name),
			logs.String("queue", queueName),
			logs.Err(err))
		return err
	}
	return nil
}

func (e *Engine) publishToSNS(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured for SNS delivery", logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	snsInvoker := e.bus.SNSInvoker()
	if snsInvoker == nil {
		return fmt.Errorf("SNS invoker not available")
	}

	message := scheduleInput(target, schedule.Name)

	// Delegate to the SNS service's Publish API. This ensures the SNS
	// service handles topic policy evaluation, subscription filtering,
	// message persistence, fan-out to all subscription endpoints (SQS,
	// Lambda, HTTP, etc.), and EventBridge bus event publication.
	messageID, err := snsInvoker.PublishToTopic(ctx, target.Arn, message, "", nil)
	if err != nil {
		logs.Debug("Failed to publish to SNS topic",
			logs.String("schedule", schedule.Name),
			logs.String("topicArn", target.Arn),
			logs.String("error", err.Error()))
		return err
	}

	logs.Debug("SNS delivery completed",
		logs.String("schedule", schedule.Name),
		logs.String("topic", target.Arn),
		logs.String("messageId", messageID))
	return nil
}

func (e *Engine) sendToKinesis(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured for Kinesis delivery", logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	kinesisInvoker := e.bus.KinesisInvoker()
	if kinesisInvoker == nil {
		logs.Debug("Kinesis invoker not available", logs.String("schedule", schedule.Name))
		return fmt.Errorf("Kinesis invoker not available")
	}

	// Extract stream name from ARN: arn:aws:kinesis:<region>:<account>:stream/<name>
	_, _, kRegion, _, _ := svcarn.SplitARN(target.Arn)
	streamName := svcarn.ExtractStreamNameFromARN(target.Arn)
	if streamName == "" {
		logs.Debug("Failed to extract stream name from Kinesis ARN",
			logs.String("schedule", schedule.Name),
			logs.String("arn", target.Arn))
		return fmt.Errorf("invalid Kinesis ARN: %s", target.Arn)
	}

	// Use PartitionKey from KinesisParameters if provided, otherwise fall
	// back to the schedule name (AWS uses a similar default behaviour).
	partitionKey := schedule.Name
	if target.KinesisParameters != nil && target.KinesisParameters.PartitionKey != "" {
		partitionKey = target.KinesisParameters.PartitionKey
	}

	// The local Kinesis service stores data as-is and GetRecords returns it
	// without additional encoding. SDK clients expect base64-encoded Data in
	// the GetRecords response, so cross-service callers must pre-encode the
	// payload to match the format that the Kinesis SDK PutRecord would send.
	data := []byte(base64.StdEncoding.EncodeToString([]byte(scheduleInput(target, schedule.Name))))

	logs.Debug("Sending to Kinesis for schedule",
		logs.String("schedule", schedule.Name),
		logs.String("stream", streamName))

	if _, _, err := kinesisInvoker.PutRecord(ctx, kRegion, streamName, partitionKey, data); err != nil {
		logs.Debug("Failed to send to Kinesis",
			logs.String("schedule", schedule.Name),
			logs.String("stream", streamName),
			logs.Err(err))
		return err
	}

	logs.Debug("Kinesis delivery completed",
		logs.String("schedule", schedule.Name),
		logs.String("stream", streamName))
	return nil
}

func (e *Engine) startStepFunctionExecution(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured, skipping Step Functions delivery",
			logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	// The target ARN's region is the delivery region: creation-time
	// validation rejects an empty-region ARN, so no fallback applies.
	_, _, smRegion, _, _ := svcarn.SplitARN(target.Arn)

	input := scheduleInput(target, schedule.Name)

	evt := &eventbus.StepFunctionsStartExecutionEvent{
		StateMachineArn: target.Arn,
		Input:           input,
	}
	evt.Region = smRegion
	evt.AccountID = e.accountID

	// PublishSync so the execution-start failure propagates back and the
	// schedule's retry policy and dead-letter routing apply to it, exactly
	// as they already do for the invoker-based deliveries. The bus
	// contract reports a handler failure in HandlerResult.Error with a
	// nil error, so both must be folded into the returned error.
	result, err := e.bus.PublishSync(ctx, evt)
	if err == nil {
		err = result.Error
	}
	if err != nil {
		logs.Debug("Failed to start Step Functions execution from schedule",
			logs.String("schedule", schedule.Name),
			logs.String("stateMachineArn", target.Arn),
			logs.String("error", err.Error()))
		return err
	}

	logs.Debug("Schedule delivered to Step Functions",
		logs.String("schedule", schedule.Name),
		logs.String("stateMachineArn", target.Arn))
	return nil
}

func (e *Engine) sendToEventBridge(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured, skipping EventBridge delivery",
			logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	// The target ARN's region is the delivery region: creation-time
	// validation rejects an empty-region ARN, so no fallback applies.
	_, _, ebRegion, _, _ := svcarn.SplitARN(target.Arn)

	// The target ARN must resolve to an event bus: the shared extractor
	// maps both event-bus/<name> and rule/<bus>/<rule> (whose first
	// resource segment names its own bus) to a bus name; every other
	// events resource yields "" and surfaces as a delivery error rather
	// than silently landing on the default bus.
	eventBusName := svcarn.ExtractEventBusNameFromARN(target.Arn)
	if eventBusName == "" {
		logs.Debug("Target ARN does not name an EventBridge event bus",
			logs.String("schedule", schedule.Name),
			logs.String("arn", target.Arn))
		return fmt.Errorf("target ARN does not name an EventBridge event bus: %s", target.Arn)
	}

	input := scheduleInput(target, schedule.Name)

	evt := &eventbus.EventBridgePutEventsEvent{
		EventBusName: eventBusName,
		Input:        input,
	}
	// Populate DetailType and Source from EventBridgeParameters so that
	// EventBridge rules can match on them.
	if target.EventBridgeParameters != nil {
		evt.DetailType = target.EventBridgeParameters.DetailType
		evt.Source = target.EventBridgeParameters.Source
	}
	evt.Region = ebRegion
	evt.AccountID = e.accountID

	// PublishSync so the PutEvents failure propagates back and the
	// schedule's retry policy and dead-letter routing apply to it, exactly
	// as they already do for the invoker-based deliveries. The bus
	// contract reports a handler failure in HandlerResult.Error with a
	// nil error, so both must be folded into the returned error.
	result, err := e.bus.PublishSync(ctx, evt)
	if err == nil {
		err = result.Error
	}
	if err != nil {
		logs.Debug("Failed to deliver schedule to EventBridge",
			logs.String("schedule", schedule.Name),
			logs.String("eventBus", eventBusName),
			logs.String("error", err.Error()))
		return err
	}

	logs.Debug("Schedule delivered to EventBridge",
		logs.String("schedule", schedule.Name),
		logs.String("eventBus", eventBusName))
	return nil
}
