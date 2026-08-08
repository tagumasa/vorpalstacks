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

// defaultAsyncMaxRetries is the AWS default for MaximumRetryAttempts on
// asynchronous invocations when no EventInvokeConfig is set.
const defaultAsyncMaxRetries = 2

// defaultAsyncMaxEventAge is the AWS default MaximumEventAgeInSeconds (6h).
const defaultAsyncMaxEventAge = int32(21600)

// invokeAsyncWithRetry executes an asynchronous Lambda invocation with
// retry and destination delivery.  It reads the function's
// EventInvokeConfig (if any) to determine MaximumRetryAttempts,
// MaximumEventAgeInSeconds, and DestinationConfig.
//
// On success: if OnSuccess destination is configured, the invocation
// record is delivered to it.
// On failure after all retries: if OnFailure destination is configured,
// the invocation record is delivered to it.
func (s *LambdaService) invokeAsyncWithRetry(
	ctx context.Context,
	function *lambdastore.Function,
	ver *lambdastore.Version,
	store *lambdastore.FunctionStore,
	region string,
	payload []byte,
	qualifier string,
) {
	// Load EventInvokeConfig for this function+qualifier.
	maxRetries := defaultAsyncMaxRetries
	maxEventAge := defaultAsyncMaxEventAge
	var destConfig *lambdastore.DestinationConfig

	if store != nil {
		eic, err := store.GetEventInvokeConfig(function.FunctionName, qualifier)
		if err == nil && eic != nil {
			if eic.MaximumRetryAttempts >= 0 {
				maxRetries = int(eic.MaximumRetryAttempts)
			}
			if eic.MaximumEventAgeInSeconds > 0 {
				maxEventAge = eic.MaximumEventAgeInSeconds
			}
			destConfig = eic.DestinationConfig
		}
	}

	startTime := time.Now()
	var lastResult *lambdastore.InvocationResult
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Check MaximumEventAge before retrying.
		if attempt > 0 {
			elapsed := int32(time.Since(startTime).Seconds())
			if elapsed >= maxEventAge {
				break
			}
			// Exponential backoff: 1s, 2s, 4s... capped at 60s.
			backoff := time.Duration(1) << uint(attempt-1)
			if backoff > 60 {
				backoff = 60
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff * time.Second):
			}
		}

		result, err := s.invokeFunction(function, ver, store, region, payload, "")
		lastResult = result
		lastErr = err

		// Determine success: no infrastructure error and no function error.
		if err == nil && result != nil && result.FunctionError == "" {
			// Success — deliver to OnSuccess destination if configured.
			if destConfig != nil && destConfig.OnSuccess != nil && destConfig.OnSuccess.Destination != "" {
				deliverDestination(ctx, s, destConfig.OnSuccess.Destination, true,
					payload, result, function, region, attempt+1)
			}
			return
		}

		logs.Warn("async invocation failed, will retry",
			logs.String("function", function.FunctionName),
			logs.Int("attempt", attempt+1),
			logs.Int("maxRetries", maxRetries+1),
			logs.Err(err))
	}

	// All retries exhausted — deliver to OnFailure destination if configured.
	if destConfig != nil && destConfig.OnFailure != nil && destConfig.OnFailure.Destination != "" {
		errPayload := []byte(`{"errorMessage":"internal error"}`)
		if lastErr != nil {
			errPayload = []byte(fmt.Sprintf(`{"errorMessage":%q}`, lastErr.Error()))
		} else if lastResult != nil && lastResult.FunctionError != "" {
			errPayload = lastResult.Payload
		}
		deliverDestination(ctx, s, destConfig.OnFailure.Destination, false,
			payload, &lambdastore.InvocationResult{
				StatusCode:    200,
				Payload:       errPayload,
				FunctionError: "Unhandled",
			}, function, region, maxRetries+1)
	}
}

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
