package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// The asynchronous retry engine itself (invokeAsyncWithRetry) lives in
// invocation_core.go; this file carries the destination delivery helpers
// it dispatches to.

// deliverDestination sends an invocation record to the configured
// destination ARN.  The destination can be an SQS queue, SNS topic,
// Lambda function, or EventBridge event bus.
func deliverDestination(
	ctx context.Context,
	s *LambdaService,
	destinationArn string,
	onSuccess bool,
	requestPayload []byte,
	result *lambdastore.InvocationResult,
	function *lambdastore.Function,
	region string,
	invokeCount int,
) {
	if s.bus == nil {
		return
	}

	condition := "RetriesExhausted"
	if onSuccess {
		condition = "EventDeliverySuccess"
	}

	record := map[string]interface{}{
		"version":   "1.0",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"requestContext": map[string]interface{}{
			"requestId":              fmt.Sprintf("dest-%d", time.Now().UnixNano()),
			"functionArn":            function.FunctionArn,
			"condition":              condition,
			"approximateInvokeCount": invokeCount,
		},
		"requestPayload": json.RawMessage(requestPayload),
	}

	if result != nil {
		record["responseContext"] = map[string]interface{}{
			"statusCode":      result.StatusCode,
			"executedVersion": result.ExecutedVersion,
			"functionError":   result.FunctionError,
		}
		record["responsePayload"] = json.RawMessage(result.Payload)
	}

	payloadBytes, err := json.Marshal(record)
	if err != nil {
		logs.Warn("failed to marshal destination payload", logs.Err(err))
		return
	}
	payloadStr := string(payloadBytes)

	service := arnutil.GetServiceFromARN(destinationArn)

	switch service {
	case "sqs":
		deliverToSQS(ctx, s, destinationArn, payloadStr, region)
	case "sns":
		deliverToSNS(ctx, s, destinationArn, payloadStr)
	case "lambda":
		deliverToLambda(ctx, s, destinationArn, payloadBytes)
	case "events":
		deliverToEventBridge(ctx, s, destinationArn, payloadStr)
	default:
		logs.Warn("unsupported destination type", logs.String("arn", destinationArn))
	}
}

func deliverToSQS(ctx context.Context, s *LambdaService, arn, payload, region string) {
	sqsInvoker := s.bus.SQSInvoker()
	if sqsInvoker == nil {
		logs.Warn("destination: SQS invoker not configured", logs.String("arn", arn))
		return
	}
	queueName := arnutil.ExtractQueueNameFromARN(arn)
	if queueName == "" {
		return
	}
	queueURL, err := sqsInvoker.GetQueueByName(ctx, region, queueName)
	if err != nil {
		logs.Warn("destination: failed to resolve SQS queue URL", logs.String("arn", arn), logs.Err(err))
		return
	}
	if _, _, err := sqsInvoker.SendMessage(ctx, region, queueURL, payload, eventbus.SQSSendOptions{}); err != nil {
		logs.Warn("destination: failed to send to SQS", logs.String("arn", arn), logs.Err(err))
	}
}

func deliverToSNS(ctx context.Context, s *LambdaService, arn, payload string) {
	snsInvoker := s.bus.SNSInvoker()
	if snsInvoker == nil {
		logs.Warn("destination: SNS invoker not configured", logs.String("arn", arn))
		return
	}
	if _, err := snsInvoker.PublishToTopic(ctx, arn, payload, "", nil); err != nil {
		logs.Warn("destination: failed to publish to SNS", logs.String("arn", arn), logs.Err(err))
	}
}

func deliverToLambda(ctx context.Context, s *LambdaService, arn string, payload []byte) {
	lambdaInvoker := s.bus.LambdaInvoker()
	if lambdaInvoker == nil {
		logs.Warn("destination: Lambda invoker not configured", logs.String("arn", arn))
		return
	}
	functionName := arnutil.ExtractFunctionNameFromARN(arn)
	if functionName == "" {
		// Fallback: extract from resource portion.
		resource := arnutil.ExtractResourceFromARN(arn)
		if idx := strings.Index(resource, "function:"); idx != -1 {
			functionName = resource[idx+len("function:"):]
		}
	}
	if functionName == "" {
		return
	}
	if _, _, err := lambdaInvoker.InvokeForGateway(ctx, functionName, payload); err != nil {
		logs.Warn("destination: failed to invoke Lambda", logs.String("arn", arn), logs.Err(err))
	}
}

func deliverToEventBridge(ctx context.Context, s *LambdaService, arn, payload string) {
	eventsInvoker := s.bus.EventsInvoker()
	if eventsInvoker == nil {
		logs.Warn("destination: EventBridge invoker not configured", logs.String("arn", arn))
		return
	}
	busName := arnutil.ExtractEventBusNameFromARN(arn)
	if busName == "" {
		logs.Warn("destination: failed to extract event bus name", logs.String("arn", arn))
		return
	}
	eventID := fmt.Sprintf("dest-%d", time.Now().UnixNano())
	key := fmt.Sprintf("events:%s:%s", busName, eventID)
	if err := eventsInvoker.PutEvent(ctx, key, payload); err != nil {
		logs.Warn("destination: failed to put EventBridge event", logs.String("arn", arn), logs.Err(err))
	}
}
