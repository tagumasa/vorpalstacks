package sfn

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

func (e *Executor) executeLambdaTask(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, input string) (string, error) {
	if e.bus == nil || e.bus.LambdaInvoker() == nil {
		return "", &taskFailure{code: "Lambda.ServiceException", cause: "lambda invoker not configured", failedToStart: true}
	}

	if parsed, err := arnutil.ParseARN(state.Resource); err != nil || parsed.Service != "lambda" {
		return "", &taskFailure{code: "Lambda.ServiceException", cause: "invalid Lambda ARN: " + state.Resource, failedToStart: true}
	}

	statusCode, payload, err := e.bus.LambdaInvoker().InvokeForGateway(ctx, state.Resource, []byte(input))
	if err != nil {
		// The invocation transport itself failed; the function never ran.
		return "", &taskFailure{code: "Lambda.ServiceException", cause: "failed to invoke Lambda function: " + err.Error(), failedToStart: true}
	}

	if statusCode != 200 {
		// A non-200 response carries the function's error: the error name
		// is the payload's errorType member (Lambda.Unknown when the
		// payload carries none), the cause its errorMessage.
		code, cause := lambdaFunctionError(payload)
		return "", &taskFailure{code: code, cause: cause}
	}

	return string(payload), nil
}

// lambdaFunctionError extracts the error name and cause from a failed Lambda
// invocation payload ({"errorType": ..., "errorMessage": ...}); an error
// without an errorType reports the documented Lambda.Unknown name.
func lambdaFunctionError(payload []byte) (string, string) {
	var fnErr struct {
		ErrorType    string `json:"errorType"`
		ErrorMessage string `json:"errorMessage"`
	}
	cause := string(payload)
	if err := json.Unmarshal(payload, &fnErr); err == nil {
		if fnErr.ErrorMessage != "" {
			cause = fnErr.ErrorMessage
		}
		if fnErr.ErrorType != "" {
			return fnErr.ErrorType, cause
		}
	}
	return "Lambda.Unknown", cause
}

// regionFromQueueURL extracts the region from an SQS queue URL host
// (https://sqs.<region>.amazonaws.com/<account>/<name>).
func regionFromQueueURL(queueURL string) string {
	u, err := url.Parse(queueURL)
	if err != nil {
		return ""
	}
	parts := strings.Split(u.Hostname(), ".")
	if len(parts) >= 2 && parts[0] == "sqs" {
		return parts[1]
	}
	return ""
}

// integrationAction extracts the API action segment from an integration
// resource: the segment after the last colon of the suffix-stripped
// identifier, at either optimised depth (…:::sqs:sendMessage) or AWS SDK
// depth (…:::aws-sdk:sqs:sendMessage).
func integrationAction(resource string) string {
	if i := strings.LastIndex(resource, ":"); i >= 0 {
		return resource[i+1:]
	}
	return resource
}

// regionOfARN returns the ARN's region segment, falling back to the state
// machine's own region when the ARN is region-less — the states
// integration pseudo-ARNs carry no region.
func (e *Executor) regionOfARN(arn string) string {
	_, _, region, _, _ := arnutil.SplitARN(arn)
	if region == "" {
		region = e.region
	}
	return region
}

// resolvedInvocation carries what the attempt's Credentials authorisation
// already resolved, so the dispatch reuses the authorisation's resource
// resolution instead of repeating its lookups. A nil carrier (dispatches
// without a Credentials authorisation, or direct dispatch tests) resolves
// on the spot.
type resolvedInvocation struct {
	sqsQueueURL string
	sqsRegion   string
}

// sqsQueueParams returns the queue this attempt's SendMessage addresses:
// the authorisation's resolution when it ran, otherwise a fresh resolution.
func (r *resolvedInvocation) sqsQueueParams(ctx context.Context, e *Executor, params map[string]interface{}) (string, string, error) {
	if r != nil && r.sqsQueueURL != "" {
		return r.sqsQueueURL, r.sqsRegion, nil
	}
	return e.resolveSqsQueueURL(ctx, params)
}

// resolveSqsQueueURL resolves the queue a SendMessage addresses. The queue
// URL is the only member that names the queue's region; the integration
// resource is a region-less states pseudo-ARN, so a QueueUrl member carries
// the region and a QueueName member resolves through the invoker in the
// state machine's own region. Both the dispatch and the Credentials
// authorisation resolve through this one path, so the authorised resource
// cannot drift from the sent one.
func (e *Executor) resolveSqsQueueURL(ctx context.Context, params map[string]interface{}) (queueURL, region string, err error) {
	queueURL = getStr(params, "QueueUrl")
	region = regionFromQueueURL(queueURL)
	if region == "" {
		region = e.region
	}
	if queueURL == "" {
		queueName := getStr(params, "QueueName")
		if queueName == "" {
			return "", "", fmt.Errorf("SQS SendMessage requires QueueUrl or QueueName")
		}
		resolved, rerr := e.bus.SQSInvoker().GetQueueByName(ctx, region, queueName)
		if rerr != nil {
			return "", "", fmt.Errorf("SQS queue not found: %s: %w", queueName, rerr)
		}
		queueURL = resolved
	}
	return queueURL, region, nil
}

// resolveSnsTopicArn resolves the topic a Publish addresses: the TopicArn
// member names the topic directly, otherwise the TopicName form carries no
// region, so the topic lives in the state machine's own region. The
// dispatch and the Credentials authorisation share this one resolution.
func (e *Executor) resolveSnsTopicArn(params map[string]interface{}) (string, error) {
	if topicArn := getStr(params, "TopicArn"); topicArn != "" {
		return topicArn, nil
	}
	topicName := getStr(params, "TopicName")
	if err := validateSNSTopicInput("", topicName); err != nil {
		return "", err
	}
	return arnutil.NewARNBuilder(e.accountID, e.region).SNS().Topic(topicName), nil
}

// entryEventBusName returns the bus a PutEvents entry addresses: the
// EventBusName member, or the account's default bus.
func entryEventBusName(entry map[string]interface{}) string {
	if ebn, ok := entry["EventBusName"].(string); ok && ebn != "" {
		return ebn
	}
	return "default"
}

// eventBusNames returns the distinct bus names the PutEvents entries
// address, in first-seen order — every bus the invocation writes to must be
// authorised.
func eventBusNames(params map[string]interface{}) []string {
	seen := map[string]bool{}
	var names []string
	if rawEntries, ok := params["Entries"].([]interface{}); ok {
		for _, entry := range rawEntries {
			entryMap, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}
			busName := entryEventBusName(entryMap)
			if seen[busName] {
				continue
			}
			seen[busName] = true
			names = append(names, busName)
		}
	}
	return names
}

func (e *Executor) executeSQSTask(ctx context.Context, resource, input string, resolved *resolvedInvocation) (string, error) {
	if e.bus == nil || e.bus.SQSInvoker() == nil {
		return "", fmt.Errorf("SQS invoker not configured")
	}

	if parsed, err := arnutil.ParseARN(resource); err != nil || parsed.Service != "states" {
		return "", fmt.Errorf("invalid SQS resource ARN: %s", resource)
	}

	action := integrationAction(resource)
	if action == "sendMessage" {
		return e.executeSQSSendMessage(ctx, input, resolved)
	}
	return "", fmt.Errorf("unsupported SQS action: %s", action)
}

func (e *Executor) executeSQSSendMessage(ctx context.Context, input string, resolved *resolvedInvocation) (string, error) {
	var inputData map[string]interface{}
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		inputData = map[string]interface{}{"MessageBody": input}
	}

	queueURL, sqsRegion, qErr := resolved.sqsQueueParams(ctx, e, inputData)
	if qErr != nil {
		return "", qErr
	}

	// SendMessage's MessageBody is a required member; a structured body
	// (the callback pattern's documented shape nests the task token
	// inside it) serialises that object — never the whole parameter set.
	messageBody := ""
	if mb, ok := inputData["MessageBody"].(string); ok {
		messageBody = mb
	} else if mb, ok := inputData["MessageBody"]; ok {
		if bodyBytes, err := json.Marshal(mb); err == nil {
			messageBody = string(bodyBytes)
		}
	}
	if messageBody == "" {
		return "", fmt.Errorf("SQS SendMessage requires a MessageBody")
	}

	messageID, md5OfBody, err := e.bus.SQSInvoker().SendMessage(ctx, sqsRegion, queueURL, messageBody, invokers.SQSSendOptions{
		DelaySeconds:           getInt64FromInput(inputData, "DelaySeconds"),
		MessageGroupID:         getStr(inputData, "MessageGroupId"),
		MessageDeduplicationID: getStr(inputData, "MessageDeduplicationId"),
		TypedMessageAttributes: extractTypedAttrs(inputData["MessageAttributes"]),
	})
	if err != nil {
		return "", fmt.Errorf("failed to send SQS message: %w", err)
	}

	result := map[string]interface{}{
		"MessageId":        messageID,
		"MD5OfMessageBody": md5OfBody,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal SQS task result: %w", err)
	}
	return string(resultJSON), nil
}

func (e *Executor) executeSNSTask(ctx context.Context, resource, input string) (string, error) {
	if e.bus == nil || e.bus.SNSInvoker() == nil {
		return "", fmt.Errorf("SNS invoker not configured")
	}

	if parsed, err := arnutil.ParseARN(resource); err != nil || parsed.Service != "states" {
		return "", fmt.Errorf("invalid SNS resource ARN: %s", resource)
	}

	action := integrationAction(resource)
	if action == "publish" {
		return e.executeSNSPublish(ctx, resource, input)
	}
	return "", fmt.Errorf("unsupported SNS action: %s", action)
}

func (e *Executor) executeSNSPublish(ctx context.Context, resource, input string) (string, error) {
	var inputData map[string]interface{}
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		inputData = map[string]interface{}{"Message": input}
	}

	topicArn, tErr := e.resolveSnsTopicArn(inputData)
	if tErr != nil {
		return "", tErr
	}

	// Publish's Message is a required member; a structured message
	// serialises that value — never the whole parameter set.
	message := ""
	if m, ok := inputData["Message"].(string); ok {
		message = m
	} else if m, ok := inputData["Message"]; ok {
		if bodyBytes, err := json.Marshal(m); err == nil {
			message = string(bodyBytes)
		}
	}
	if message == "" {
		return "", fmt.Errorf("SNS Publish requires a Message")
	}

	subject := ""
	if s, ok := inputData["Subject"].(string); ok {
		subject = s
	}

	// Message attributes keep their typed form: DataType plus
	// StringValue or the Base64-decoded BinaryValue. The delivery event
	// carries them as raw JSON in the SNS attribute transport shape.
	typedAttrs := extractTypedAttrs(inputData["MessageAttributes"])
	attrTransport := make(map[string]json.RawMessage, len(typedAttrs))
	storedAttrs := make(map[string]interface{}, len(typedAttrs))
	for k, v := range typedAttrs {
		transport := snsAttributeTransport{Type: v.DataType, StringValue: v.StringValue, BinaryValue: v.BinaryValue}
		raw, terr := json.Marshal(transport)
		if terr != nil {
			continue
		}
		attrTransport[k] = raw
		storedAttrs[k] = transport
	}

	msgID := uuid.New().String()
	msg := map[string]interface{}{
		"MessageId":         msgID,
		"TopicArn":          topicArn,
		"Subject":           subject,
		"Message":           message,
		"MessageAttributes": storedAttrs,
	}

	if err := e.bus.SNSInvoker().StoreMessage(ctx, topicArn+":messages:"+msgID, msg); err != nil {
		return "", fmt.Errorf("failed to store SNS message: %w", err)
	}

	// The bus was nil-checked by the caller (executeSNSTask); publishing
	// is fire-and-forget — a delivery failure logs, it does not fail the
	// state, the message is already stored.
	snsEvt := &eventbus.SNSDeliveryEvent{
		TopicARN:          topicArn,
		MessageID:         msgID,
		Message:           message,
		Subject:           subject,
		MessageAttributes: attrTransport,
	}
	snsEvt.Region = e.regionOfARN(topicArn)
	if err := e.bus.Publish(context.Background(), snsEvt); err != nil {
		logs.Warn("failed to publish SNS event from Step Functions", logs.Err(err))
	}

	result := map[string]interface{}{
		"MessageId": msgID,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal SNS task result: %w", err)
	}
	return string(resultJSON), nil
}

func (e *Executor) executeEventsTask(ctx context.Context, execCtx *ExecutionContext, resource, input string) (string, error) {
	if e.bus == nil || e.bus.EventsInvoker() == nil {
		return "", fmt.Errorf("events invoker not configured")
	}

	resourceParts := strings.Split(resource, ":")
	if len(resourceParts) < 7 {
		return "", fmt.Errorf("invalid Events resource ARN: %s", resource)
	}

	action := integrationAction(resource)
	if action == "putEvents" {
		return e.executeEventsPutEvents(ctx, execCtx, resource, input)
	}
	return "", fmt.Errorf("unsupported Events action: %s", action)
}

func (e *Executor) executeEventsPutEvents(ctx context.Context, execCtx *ExecutionContext, resource, input string) (string, error) {
	var inputData map[string]interface{}
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		inputData = map[string]interface{}{"Entries": []interface{}{}}
	}

	eventsRegion := e.regionOfARN(resource)

	var entries []map[string]interface{}
	if rawEntries, ok := inputData["Entries"].([]interface{}); ok {
		for _, entry := range rawEntries {
			if entryMap, ok := entry.(map[string]interface{}); ok {
				entries = append(entries, entryMap)
			}
		}
	}
	// PutEvents takes the entries array and nothing else; a payload
	// without Entries is an invocation failure, not a licence to invent
	// an event.
	if len(entries) == 0 {
		return "", fmt.Errorf("PutEvents requires at least one entry")
	}

	var results []map[string]interface{}
	failedCount := 0
	firstFailure := ""
	now := time.Now().UTC()

	for _, entry := range entries {
		eventBusName := entryEventBusName(entry)

		detail := map[string]interface{}{}
		if d, ok := entry["Detail"]; ok {
			if detailMap, ok := d.(map[string]interface{}); ok {
				detail = detailMap
			} else if detailStr, ok := d.(string); ok {
				if jsonErr := json.Unmarshal([]byte(detailStr), &detail); jsonErr != nil {
					detail = map[string]interface{}{"raw": detailStr}
				}
			}
		}

		// Entry Resources: the user's values first, then "the execution
		// ARN and the state machine ARN are automatically appended to the
		// Resources field of each PutEventsRequestEntry".
		resources := []string{}
		if raw, ok := entry["Resources"].([]interface{}); ok {
			for _, r := range raw {
				if s, ok := r.(string); ok && s != "" {
					resources = append(resources, s)
				}
			}
		}
		resources = append(resources, execCtx.Execution.ExecutionArn, execCtx.Execution.StateMachineArn)

		// An entry Time is honoured when supplied (an RFC3339
		// timestamp); otherwise the publication instant applies.
		entryTime := now
		if ts, ok := entry["Time"].(string); ok {
			if parsed, perr := time.Parse(time.RFC3339, ts); perr == nil {
				entryTime = parsed
			}
		}

		event := map[string]interface{}{
			// EventBridge event IDs are UUIDs; entries published in the
			// same tick must still receive distinct IDs (and storage
			// keys), which a shared clock value cannot guarantee.
			"ID":           uuid.New().String(),
			"EventBusName": eventBusName,
			"Source":       getStr(entry, "Source"),
			"DetailType":   getStr(entry, "DetailType"),
			"Time":         entryTime,
			"Region":       eventsRegion,
			"Account":      e.accountID,
			"Detail":       detail,
			"Resources":    resources,
		}

		key := fmt.Sprintf("events:%s:%s", eventBusName, event["ID"])
		if err := e.bus.EventsInvoker().PutEvent(ctx, key, event); err != nil {
			// "PutEvents returns the number of failed entries in the
			// FailedEntryCount field" — a failed entry is a per-entry
			// result, not an invocation abort.
			failedCount++
			errMsg := fmt.Sprintf("failed to store event: %s", err.Error())
			if firstFailure == "" {
				firstFailure = errMsg
			}
			results = append(results, map[string]interface{}{
				"EventId":      event["ID"],
				"ErrorCode":    "InternalFailure",
				"ErrorMessage": errMsg,
			})
			continue
		}

		// The bus was nil-checked by the caller (executeEventsTask);
		// publishing is fire-and-forget — the event is already stored.
		if eventJSON, merr := json.Marshal(event); merr == nil {
			ebEvt := &eventbus.EventBridgeDeliveryEvent{
				TargetARN: arnutil.NewARNBuilder(e.accountID, eventsRegion).Events().EventBus(eventBusName),
				Input:     eventJSON,
			}
			ebEvt.Region = eventsRegion
			if err := e.bus.Publish(context.Background(), ebEvt); err != nil {
				logs.Warn("failed to publish EventBridge event from Step Functions", logs.Err(err))
			}
		}

		results = append(results, map[string]interface{}{"EventId": event["ID"]})
	}

	// "Step Functions checks whether the FailedEntryCount is greater than
	// zero. If it is greater than zero, Step Functions fails the state
	// with the error EventBridge.FailedEntry." — the retry note applies:
	// a retried task resubmits the original entries array.
	if failedCount > 0 {
		return "", &taskFailure{
			code:  "EventBridge.FailedEntry",
			cause: fmt.Sprintf("PutEvents failed %d of %d entries: %s", failedCount, len(entries), firstFailure),
		}
	}

	result := map[string]interface{}{
		"FailedEntryCount": failedCount,
		"Entries":          results,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Events task result: %w", err)
	}
	return string(resultJSON), nil
}

// executeActivityTask schedules the activity task for one attempt. The
// token is minted by executeTask for this specific attempt and passed
// explicitly, so concurrent Map/Parallel branches and retry attempts never
// share or clobber each other's tokens.
func (e *Executor) executeActivityTask(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, input string, taskToken string, timeoutSeconds, heartbeatSeconds int32) (string, error) {
	parts := strings.Split(state.Resource, ":")
	if len(parts) < 7 {
		return "", fmt.Errorf("invalid activity ARN: %s", state.Resource)
	}

	activityName := parts[6]
	activityArn := arnutil.NewARNBuilder(e.accountID, e.region).StepFunctions().Activity(activityName)

	_, err := e.store.GetActivity(ctx, activityArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrActivityNotFound) {
			e.logActivityScheduleFailed(ctx, execCtx, "activity not found: "+activityArn)
			return "", fmt.Errorf("activity not found: %s", activityArn)
		}
		e.logActivityScheduleFailed(ctx, execCtx, "failed to read the activity record: "+err.Error())
		return "", fmt.Errorf("failed to read the activity record: %w", err)
	}

	task := &sfnstore.ActivityTask{
		ActivityArn:  activityArn,
		ExecutionArn: execCtx.Execution.ExecutionArn,
		Input:        input,
		// The attempt's token — the same value $$.Task.Token resolved to
		// while building the input, so the worker can return it with
		// SendTaskSuccess/SendTaskFailure.
		TaskToken: taskToken,
	}

	if err := e.store.CreateActivityTask(task); err != nil {
		e.logActivityScheduleFailed(ctx, execCtx, "failed to create activity task: "+err.Error())
		return "", fmt.Errorf("failed to create activity task: %w", err)
	}

	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn: execCtx.Execution.ExecutionArn,
		EventId:      execCtx.nextEventId(),
		Type:         "ActivityScheduled",
		ActivityTaskScheduledEventDetails: &sfnstore.ActivityTaskScheduledEventDetails{
			Resource:         state.Resource,
			Input:            input,
			TimeoutInSeconds: timeoutSeconds,
			HeartbeatSeconds: heartbeatSeconds,
		},
	})

	// An activity attempt without TimeoutSeconds waits the documented
	// Task-state default, not a transport-style short cap.
	timeout := taskWaitTimeout(timeoutSeconds)

	hbTimeout := time.Duration(heartbeatSeconds) * time.Second

	result, err := e.store.WaitForTaskResult(ctx, task.TaskToken, timeout, hbTimeout)
	if err != nil {
		if err == sfnstore.ErrTaskTimeout || err == sfnstore.ErrHeartbeatTimeout {
			eventId := execCtx.nextEventId()
			e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
				ExecutionArn:    execCtx.Execution.ExecutionArn,
				EventId:         eventId,
				PreviousEventId: eventId - 1,
				Type:            "ActivityTimedOut",
				Timestamp:       time.Now().UTC(),
				ActivityTaskTimedOutEventDetails: &sfnstore.ActivityTaskTimedOutEventDetails{
					Error: "States.Timeout",
					Cause: "Task timed out",
				},
			})
			return "", &taskFailure{code: "States.Timeout", cause: "Task timed out", timedOut: true}
		}
		return "", err
	}

	if result.Error != nil {
		e.logActivityStarted(ctx, execCtx, result.WorkerName)
		eventId := execCtx.nextEventId()
		e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
			ExecutionArn:    execCtx.Execution.ExecutionArn,
			EventId:         eventId,
			PreviousEventId: eventId - 1,
			Type:            "ActivityFailed",
			Timestamp:       time.Now().UTC(),
			ActivityTaskFailedEventDetails: &sfnstore.ActivityTaskFailedEventDetails{
				Error: result.Error.Error(),
				Cause: result.Cause,
			},
		})
		// The worker's report keeps its own error name so the state's
		// Catch/Retry can match it.
		code := result.Error.Error()
		if code == "" {
			code = "States.TaskFailed"
		}
		return "", &taskFailure{code: code, cause: result.Cause}
	}

	e.logActivityStarted(ctx, execCtx, result.WorkerName)
	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "ActivitySucceeded",
		Timestamp:       time.Now().UTC(),
		ActivityTaskSucceededEventDetails: &sfnstore.ActivityTaskSucceededEventDetails{
			Output: result.Output,
		},
	})

	return result.Output, nil
}

// logActivityStarted records the modelled ActivityStarted event once the
// executor learns a worker claimed the task. The store stamps the claim at
// poll time; the event is appended when the attempt's result arrives so it
// precedes the terminal event in the history.
func (e *Executor) logActivityStarted(ctx context.Context, execCtx *ExecutionContext, workerName string) {
	if workerName == "" {
		return
	}
	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "ActivityStarted",
		Timestamp:       time.Now().UTC(),
		ActivityTaskStartedEventDetails: &sfnstore.ActivityTaskStartedEventDetails{
			WorkerName: workerName,
		},
	})
}

// logActivityScheduleFailed records the modelled ActivityScheduleFailed
// event when creating an activity task record fails; the error keeps the
// documented wildcard name for task failures, the cause carries the
// schedule failure.
func (e *Executor) logActivityScheduleFailed(ctx context.Context, execCtx *ExecutionContext, cause string) {
	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "ActivityScheduleFailed",
		Timestamp:       time.Now().UTC(),
		ActivityScheduleFailedEventDetails: &sfnstore.ActivityScheduleFailedEventDetails{
			Error: "States.TaskFailed",
			Cause: cause,
		},
	})
}

func (e *Executor) executeDynamoDBTask(ctx context.Context, resource, input string) (string, error) {
	if e.bus == nil || e.bus.DynamoDBInvoker() == nil {
		return "", fmt.Errorf("DynamoDB invoker not configured")
	}

	resourceParts := strings.Split(resource, ":")
	if len(resourceParts) < 7 {
		return "", fmt.Errorf("invalid DynamoDB resource ARN: %s", resource)
	}

	action := integrationAction(resource)

	var inputData map[string]interface{}
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		return "", fmt.Errorf("invalid input for DynamoDB task: %w", err)
	}

	tableName := getStr(inputData, "TableName")
	if tableName == "" {
		return "", fmt.Errorf("TableName is required for DynamoDB task")
	}

	region := e.region

	switch action {
	case "getItem":
		return e.executeDynamoDBGetItem(ctx, region, tableName, inputData)
	case "putItem":
		return e.executeDynamoDBPutItem(ctx, region, tableName, inputData)
	case "deleteItem":
		return e.executeDynamoDBDeleteItem(ctx, region, tableName, inputData)
	case "updateItem":
		return e.executeDynamoDBUpdateItem(ctx, region, tableName, inputData)
	default:
		return "", fmt.Errorf("unsupported DynamoDB action: %s", action)
	}
}

func (e *Executor) executeDynamoDBGetItem(ctx context.Context, region, tableName string, inputData map[string]interface{}) (string, error) {
	keyRaw, ok := inputData["Key"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("key is required for getItem")
	}
	key := awsAttrValuesToPlain(keyRaw)

	item, err := e.bus.DynamoDBInvoker().GetItem(ctx, region, tableName, key)
	if err != nil {
		return "", fmt.Errorf("DynamoDB GetItem failed: %w", err)
	}

	// A GetItem that matches nothing omits the Item element entirely —
	// an explicit null would read as present to a downstream
	// IsPresent-style check.
	if item == nil {
		return "{}", nil
	}
	result := map[string]interface{}{"Item": plainMapToAWSAttrValues(item)}
	resultJSON, _ := json.Marshal(result)
	return string(resultJSON), nil
}

func (e *Executor) executeDynamoDBPutItem(ctx context.Context, region, tableName string, inputData map[string]interface{}) (string, error) {
	itemRaw, ok := inputData["Item"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("item is required for putItem")
	}
	plainItem := awsAttrValuesToPlain(itemRaw)

	// The DynamoDB store's buildItemKey extracts PK/SK from the key map
	// using the table schema, so passing the entire item as key is safe.
	// Attributes are merged with key during storage.
	result, err := e.bus.DynamoDBInvoker().PutItem(ctx, region, tableName, plainItem, nil)
	if err != nil {
		return "", fmt.Errorf("DynamoDB PutItem failed: %w", err)
	}

	resultJSON, _ := json.Marshal(map[string]interface{}{"Attributes": plainMapToAWSAttrValues(result)})
	return string(resultJSON), nil
}

func (e *Executor) executeDynamoDBDeleteItem(ctx context.Context, region, tableName string, inputData map[string]interface{}) (string, error) {
	keyRaw, ok := inputData["Key"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("key is required for deleteItem")
	}
	key := awsAttrValuesToPlain(keyRaw)

	if err := e.bus.DynamoDBInvoker().DeleteItem(ctx, region, tableName, key); err != nil {
		return "", fmt.Errorf("DynamoDB DeleteItem failed: %w", err)
	}

	return "{}", nil
}

// updateExprResult captures the parsed result of a DynamoDB UpdateExpression,
// separating SET, ADD, REMOVE and DELETE operations so that the caller can
// apply them to an existing item via read-modify-write.
type updateExprResult struct {
	Sets    map[string]interface{}
	Adds    map[string]interface{}
	Deletes map[string]interface{}
	Removes []string
}

// tokenizeUpdateExpression splits a DynamoDB UpdateExpression into tokens,
// normalising whitespace and ensuring that '=' and ',' are standalone tokens.
// This handles tabs, newlines, consecutive spaces, and expressions written
// without spaces around '=' (e.g. "SET attr=:val").
func tokenizeUpdateExpression(expr string) []string {
	expr = strings.ReplaceAll(expr, "=", " = ")
	expr = strings.ReplaceAll(expr, ",", " , ")
	return strings.Fields(expr)
}

func (e *Executor) executeDynamoDBUpdateItem(ctx context.Context, region, tableName string, inputData map[string]interface{}) (string, error) {
	keyRaw, ok := inputData["Key"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("key is required for updateItem")
	}
	key := awsAttrValuesToPlain(keyRaw)

	if updateExpr, ok := inputData["UpdateExpression"].(string); ok && updateExpr != "" {
		exprAttrValues := map[string]interface{}{}
		if av, ok := inputData["ExpressionAttributeValues"].(map[string]interface{}); ok {
			exprAttrValues = awsAttrValuesToPlain(av)
		}
		exprAttrNames := map[string]string{}
		if en, ok := inputData["ExpressionAttributeNames"].(map[string]interface{}); ok {
			for placeholder, val := range en {
				if s, ok := val.(string); ok {
					exprAttrNames[placeholder] = s
				}
			}
		}
		parsed, err := resolveUpdateExpression(updateExpr, exprAttrValues, exprAttrNames)
		if err != nil {
			return "", fmt.Errorf("DynamoDB UpdateItem expression error: %w", err)
		}

		existing, err := e.bus.DynamoDBInvoker().GetItem(ctx, region, tableName, key)
		if err != nil {
			return "", fmt.Errorf("DynamoDB UpdateItem read failed: %w", err)
		}

		merged := make(map[string]interface{})
		for k, v := range existing {
			merged[k] = v
		}
		for k := range key {
			delete(merged, k)
		}
		for name, val := range parsed.Sets {
			merged[name] = val
		}
		for name, addVal := range parsed.Adds {
			sum, aerr := applyDynamoDBAdd(merged[name], addVal)
			if aerr != nil {
				return "", aerr
			}
			merged[name] = sum
		}
		for name, delVal := range parsed.Deletes {
			remaining, derr := applyDynamoDBDelete(merged[name], delVal)
			if derr != nil {
				return "", derr
			}
			if remaining == nil {
				delete(merged, name)
			} else {
				merged[name] = remaining
			}
		}
		for _, name := range parsed.Removes {
			delete(merged, name)
		}

		if err := e.bus.DynamoDBInvoker().UpdateItem(ctx, region, tableName, key, merged); err != nil {
			return "", fmt.Errorf("DynamoDB UpdateItem failed: %w", err)
		}
		return "{}", nil
	}

	if av, ok := inputData["AttributeValues"].(map[string]interface{}); ok {
		attrValues := awsAttrValuesToPlain(av)
		if err := e.bus.DynamoDBInvoker().UpdateItem(ctx, region, tableName, key, attrValues); err != nil {
			return "", fmt.Errorf("DynamoDB UpdateItem failed: %w", err)
		}
		return "{}", nil
	}

	if ev, ok := inputData["ExpressionAttributeValues"].(map[string]interface{}); ok && len(ev) > 0 {
		return "", fmt.Errorf("ExpressionAttributeValues provided without UpdateExpression")
	}

	return "{}", nil
}

// resolveUpdateExpression parses a DynamoDB UpdateExpression and returns the
// parsed SET, ADD, DELETE and REMOVE operations. Placeholder values (e.g.
// :val) are resolved from exprAttrValues, and placeholder names (e.g. #attr)
// from exprAttrNames; an unresolvable placeholder or a clause that does not
// parse is a validation error, as at DynamoDB itself — a malformed update
// never silently commits no operations. The caller is responsible for
// applying these operations to the existing item via read-modify-write.
func resolveUpdateExpression(expr string, exprAttrValues map[string]interface{}, exprAttrNames map[string]string) (*updateExprResult, error) {
	result := &updateExprResult{
		Sets:    map[string]interface{}{},
		Adds:    map[string]interface{}{},
		Deletes: map[string]interface{}{},
		Removes: []string{},
	}

	tokens := tokenizeUpdateExpression(expr)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("Invalid UpdateExpression: the expression is empty")
	}

	// The token stream splits at the clause keywords; everything between
	// two keywords is the earlier clause's operand body.
	type updateClause struct {
		keyword string
		body    []string
	}
	var clauses []updateClause
	for _, token := range tokens {
		upper := strings.ToUpper(token)
		if upper == "SET" || upper == "ADD" || upper == "REMOVE" || upper == "DELETE" {
			clauses = append(clauses, updateClause{keyword: upper})
			continue
		}
		if len(clauses) == 0 {
			return nil, fmt.Errorf("Invalid UpdateExpression: expected SET, ADD, REMOVE or DELETE, found %q", token)
		}
		clauses[len(clauses)-1].body = append(clauses[len(clauses)-1].body, token)
	}

	// operandPairs splits a name-value clause body (ADD, DELETE) at its
	// comma separators.
	operandPairs := func(body []string, clause string) ([][2]string, error) {
		var pairs [][2]string
		i := 0
		for i < len(body) {
			if i+1 >= len(body) {
				return nil, fmt.Errorf("Invalid UpdateExpression: %s expects a name and a value", clause)
			}
			pairs = append(pairs, [2]string{body[i], body[i+1]})
			i += 2
			if i < len(body) {
				if body[i] != "," {
					return nil, fmt.Errorf("Invalid UpdateExpression: expected , between %s actions, found %q", clause, body[i])
				}
				i++
			}
		}
		return pairs, nil
	}

	for _, clause := range clauses {
		body := clause.body
		switch clause.keyword {
		case "SET":
			if len(body) == 0 {
				return nil, fmt.Errorf("Invalid UpdateExpression: SET names no action")
			}
			i := 0
			for i < len(body) {
				if i+2 > len(body)-1 || body[i+1] != "=" {
					return nil, fmt.Errorf("Invalid UpdateExpression: SET expects name = value, found %q", strings.Join(body[i:], " "))
				}
				name := resolveAttrName(body[i], exprAttrNames)
				val := resolveAttrValue(body[i+2], exprAttrValues)
				if name == "" || val == nil {
					return nil, fmt.Errorf("Invalid UpdateExpression: SET operand %q does not resolve", body[i])
				}
				result.Sets[name] = val
				i += 3
				if i < len(body) {
					if body[i] != "," {
						return nil, fmt.Errorf("Invalid UpdateExpression: expected , between SET actions, found %q", body[i])
					}
					i++
				}
			}
		case "REMOVE":
			if len(body) == 0 {
				return nil, fmt.Errorf("Invalid UpdateExpression: REMOVE names no attribute")
			}
			i := 0
			for i < len(body) {
				name := resolveAttrName(body[i], exprAttrNames)
				if name == "" {
					return nil, fmt.Errorf("Invalid UpdateExpression: REMOVE operand %q does not resolve", body[i])
				}
				result.Removes = append(result.Removes, name)
				i++
				if i < len(body) {
					if body[i] != "," {
						return nil, fmt.Errorf("Invalid UpdateExpression: expected , between REMOVE attributes, found %q", body[i])
					}
					i++
				}
			}
		case "ADD", "DELETE":
			// DELETE removes elements from a set: the clause names the
			// attribute and the set of elements to remove.
			if len(body) == 0 {
				return nil, fmt.Errorf("Invalid UpdateExpression: %s names no action", clause.keyword)
			}
			actions, perr := operandPairs(body, clause.keyword)
			if perr != nil {
				return nil, perr
			}
			for _, action := range actions {
				name := resolveAttrName(action[0], exprAttrNames)
				val := resolveAttrValue(action[1], exprAttrValues)
				if name == "" || val == nil {
					return nil, fmt.Errorf("Invalid UpdateExpression: %s operand %q does not resolve", clause.keyword, action[0])
				}
				if clause.keyword == "ADD" {
					result.Adds[name] = val
				} else {
					result.Deletes[name] = val
				}
			}
		}
	}

	return result, nil
}

// applyDynamoDBAdd applies one ADD operation per the DynamoDB update
// semantics: numbers add to numbers, sets union with sets, an absent
// attribute takes the value outright, and any other combination is a
// validation error rather than a silent replacement.
func applyDynamoDBAdd(existing, addVal interface{}) (interface{}, error) {
	if existing == nil {
		return addVal, nil
	}
	if existingNum, ok := toFloat64(existing); ok {
		if addNum, ok := toFloat64(addVal); ok {
			return existingNum + addNum, nil
		}
		return nil, fmt.Errorf("ADD expects a number for a numeric attribute")
	}
	existingSet, ok1 := existing.([]interface{})
	addSet, ok2 := addVal.([]interface{})
	if ok1 && ok2 {
		// A DynamoDB set holds no duplicates: ADD unions the sets, keeping
		// the existing order and appending only elements not yet present.
		seen := make(map[interface{}]bool, len(existingSet))
		for _, v := range existingSet {
			seen[v] = true
		}
		union := append([]interface{}{}, existingSet...)
		for _, v := range addSet {
			if !seen[v] {
				union = append(union, v)
				seen[v] = true
			}
		}
		return union, nil
	}
	return nil, fmt.Errorf("ADD only supports number and set attributes")
}

// applyDynamoDBDelete applies one DELETE operation: the clause removes the
// given elements from a set attribute. Deleting from an absent attribute is
// a no-op; a non-set attribute is a validation error.
func applyDynamoDBDelete(existing, delVal interface{}) (interface{}, error) {
	if existing == nil {
		return nil, nil
	}
	existingSet, ok1 := existing.([]interface{})
	delSet, ok2 := delVal.([]interface{})
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("DELETE only supports set attributes")
	}
	remove := map[interface{}]bool{}
	for _, v := range delSet {
		remove[v] = true
	}
	remaining := make([]interface{}, 0, len(existingSet))
	for _, v := range existingSet {
		if !remove[v] {
			remaining = append(remaining, v)
		}
	}
	if len(remaining) == 0 {
		return nil, nil
	}
	return remaining, nil
}

func resolveAttrName(token string, exprAttrNames map[string]string) string {
	token = strings.TrimSuffix(token, ",")
	if strings.HasPrefix(token, "#") {
		if name, ok := exprAttrNames[token]; ok {
			return name
		}
		return ""
	}
	return token
}

func resolveAttrValue(token string, exprAttrValues map[string]interface{}) interface{} {
	token = strings.TrimSuffix(token, ",")
	if strings.HasPrefix(token, ":") {
		if val, ok := exprAttrValues[token]; ok {
			return val
		}
		return nil
	}
	return token
}

func awsAttrValuesToPlain(av map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(av))
	for k, v := range av {
		out[k] = awsAttrValueToPlain(v)
	}
	return out
}

func awsAttrValueToPlain(v interface{}) interface{} {
	m, ok := v.(map[string]interface{})
	if !ok {
		return v
	}
	if s, ok := m["S"].(string); ok {
		return s
	}
	if n, ok := m["N"].(string); ok {
		// Numbers stay numbers through the plain form so a read-back
		// serialises them as N again instead of drifting to S.
		if f, err := strconv.ParseFloat(n, 64); err == nil {
			return f
		}
		return n
	}
	if b, ok := m["BOOL"].(bool); ok {
		return b
	}
	if _, ok := m["NULL"]; ok {
		return nil
	}
	if s, ok := m["SS"].([]interface{}); ok {
		strs := make([]string, len(s))
		for i, v := range s {
			strs[i] = fmt.Sprintf("%v", v)
		}
		return strs
	}
	if ns, ok := m["NS"].([]interface{}); ok {
		nums := make([]interface{}, len(ns))
		for i, v := range ns {
			if f, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64); err == nil {
				nums[i] = f
				continue
			}
			nums[i] = v
		}
		return nums
	}
	return v
}

func plainMapToAWSAttrValues(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return nil
	}
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = plainToAWSAttrValue(v)
	}
	return out
}

func plainToAWSAttrValue(v interface{}) interface{} {
	switch val := v.(type) {
	case string:
		return map[string]interface{}{"S": val}
	case bool:
		return map[string]interface{}{"BOOL": val}
	case nil:
		return map[string]interface{}{"NULL": true}
	case float64:
		return map[string]interface{}{"N": strconv.FormatFloat(val, 'f', -1, 64)}
	case int:
		return map[string]interface{}{"N": strconv.Itoa(val)}
	case int64:
		return map[string]interface{}{"N": strconv.FormatInt(val, 10)}
	case json.Number:
		return map[string]interface{}{"N": val.String()}
	case []string:
		ss := make([]interface{}, len(val))
		for i, s := range val {
			ss[i] = s
		}
		return map[string]interface{}{"SS": ss}
	case []interface{}:
		allNumeric := len(val) > 0
		for _, e := range val {
			if _, ok := e.(float64); !ok {
				allNumeric = false
				break
			}
		}
		if allNumeric {
			ns := make([]interface{}, len(val))
			for i, e := range val {
				ns[i] = strconv.FormatFloat(e.(float64), 'f', -1, 64)
			}
			return map[string]interface{}{"NS": ns}
		}
		ss := make([]interface{}, len(val))
		for i, s := range val {
			ss[i] = fmt.Sprintf("%v", s)
		}
		return map[string]interface{}{"SS": ss}
	default:
		return map[string]interface{}{"S": fmt.Sprintf("%v", val)}
	}
}

func getStr(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok && v != "" {
		return v
	}
	return ""
}

// isActivityResource reports whether the integration resource addresses a
// Step Functions activity in any partition: the parsed ARN's service is
// states and its resource path starts at activity:.
func isActivityResource(resource string) bool {
	parsed, err := arnutil.ParseARN(resource)
	return err == nil && parsed.Service == "states" && strings.HasPrefix(parsed.Resource, "activity:")
}

// snsAttributeTransport is the raw-JSON attribute shape the SNS delivery
// event carries on the bus (the SNS store's attribute serialisation).
type snsAttributeTransport struct {
	Type        string `json:"type"`
	StringValue string `json:"string_value,omitempty"`
	BinaryValue []byte `json:"binary_value,omitempty"`
}

// extractTypedAttrs renders integration MessageAttributes parameters to
// the typed attribute value: DataType ("String", "Number" or "Binary"),
// StringValue for the value types, and BinaryValue decoded from the
// Base64 transport the APIs define. A plain string value is the String
// shorthand.
func extractTypedAttrs(raw interface{}) map[string]invokers.SQSMessageAttribute {
	if raw == nil {
		return nil
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	result := make(map[string]invokers.SQSMessageAttribute, len(m))
	for k, v := range m {
		switch attr := v.(type) {
		case map[string]interface{}:
			typed := invokers.SQSMessageAttribute{DataType: "String"}
			if dt, ok := attr["DataType"].(string); ok && dt != "" {
				typed.DataType = dt
			}
			if sv, ok := attr["StringValue"].(string); ok {
				typed.StringValue = sv
			}
			if bv, ok := attr["BinaryValue"].(string); ok && bv != "" {
				if decoded, derr := base64.StdEncoding.DecodeString(bv); derr == nil {
					typed.BinaryValue = decoded
				}
			}
			result[k] = typed
		case string:
			result[k] = invokers.SQSMessageAttribute{DataType: "String", StringValue: attr}
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func getInt64FromInput(m map[string]interface{}, key string) int64 {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int64(n)
	case int64:
		return n
	case int:
		return int64(n)
	}
	return 0
}
