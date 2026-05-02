package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/endpoint"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	snsstore "vorpalstacks/internal/store/aws/sns"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

func (e *Executor) executeLambdaTask(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, input string) (string, error) {
	if e.bus == nil || e.bus.LambdaInvoker() == nil {
		return "", fmt.Errorf("lambda invoker not configured")
	}

	parts := strings.Split(state.Resource, ":")
	if len(parts) < 7 {
		return "", fmt.Errorf("invalid Lambda ARN: %s", state.Resource)
	}

	statusCode, payload, err := e.bus.LambdaInvoker().InvokeForGateway(ctx, state.Resource, []byte(input))
	if err != nil {
		return "", fmt.Errorf("failed to invoke Lambda function: %w", err)
	}

	if statusCode != 200 {
		return "", fmt.Errorf("lambda function returned status code: %d, error: %s", statusCode, string(payload))
	}

	return string(payload), nil
}

func (e *Executor) executeSQSTask(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, input string) (string, error) {
	if e.bus == nil || e.bus.SQSInvoker() == nil {
		return "", fmt.Errorf("SQS invoker not configured")
	}

	resourceParts := strings.Split(state.Resource, ":")
	if len(resourceParts) < 7 {
		return "", fmt.Errorf("invalid SQS resource ARN: %s", state.Resource)
	}

	action := resourceParts[6]
	if action == "sendMessage" || strings.HasSuffix(action, ".sendMessage") {
		return e.executeSQSSendMessage(ctx, execCtx, state, input)
	}

	return "", fmt.Errorf("unsupported SQS action: %s", action)
}

func (e *Executor) executeSQSSendMessage(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, input string) (string, error) {
	var inputData map[string]interface{}
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		inputData = map[string]interface{}{"MessageBody": input}
	}

	queueURL := ""
	if qu, ok := inputData["QueueUrl"].(string); ok && qu != "" {
		queueURL = qu
	} else if qu, ok := inputData["queueUrl"].(string); ok && qu != "" {
		queueURL = qu
	}

	if queueURL == "" {
		queueName := ""
		if qn, ok := inputData["QueueName"].(string); ok {
			queueName = qn
		} else if qn, ok := inputData["queueName"].(string); ok {
			queueName = qn
		}

		if queueName == "" {
			queueName = "default-queue"
		}

		if qURL, qErr := e.bus.SQSInvoker().GetQueueByName(ctx, queueName); qErr == nil {
			queueURL = qURL
		} else {
			queueURL = endpoint.SQSQueueURL(e.accountID, queueName)
		}
	}

	messageBody := input
	if mb, ok := inputData["MessageBody"].(string); ok {
		messageBody = mb
	} else if mb, ok := inputData["messageBody"].(string); ok {
		messageBody = mb
	} else {
		if bodyBytes, err := json.Marshal(inputData); err == nil {
			messageBody = string(bodyBytes)
		}
	}

	messageID, md5OfBody, err := e.bus.SQSInvoker().SendMessage(ctx, queueURL, messageBody, 0, nil)
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

func (e *Executor) executeSNSTask(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, input string) (string, error) {
	if e.bus == nil || e.bus.SNSInvoker() == nil {
		return "", fmt.Errorf("SNS invoker not configured")
	}

	resourceParts := strings.Split(state.Resource, ":")
	if len(resourceParts) < 7 {
		return "", fmt.Errorf("invalid SNS resource ARN: %s", state.Resource)
	}

	action := resourceParts[6]
	if action == "publish" || strings.HasSuffix(action, ".publish") {
		return e.executeSNSPublish(ctx, execCtx, state, input)
	}

	return "", fmt.Errorf("unsupported SNS action: %s", action)
}

func (e *Executor) executeSNSPublish(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, input string) (string, error) {
	var inputData map[string]interface{}
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		inputData = map[string]interface{}{"Message": input}
	}

	topicArn := ""
	if ta, ok := inputData["TopicArn"].(string); ok && ta != "" {
		topicArn = ta
	} else {
		topicName := ""
		if tn, ok := inputData["TopicName"].(string); ok {
			topicName = tn
		} else if tn, ok := inputData["topicName"].(string); ok {
			topicName = tn
		}
		if topicName == "" {
			topicName = "default-topic"
		}

		_, _, topicRegion, _, _ := arnutil.SplitARN(state.Resource)
		if topicRegion == "" {
			topicRegion = e.region
		}
		topicArn = arnutil.NewARNBuilder(e.accountID, topicRegion).SNS().Topic(topicName)
	}

	message := input
	if m, ok := inputData["Message"].(string); ok {
		message = m
	} else if m, ok := inputData["message"].(string); ok {
		message = m
	} else {
		if bodyBytes, err := json.Marshal(inputData); err == nil {
			message = string(bodyBytes)
		}
	}

	subject := ""
	if s, ok := inputData["Subject"].(string); ok {
		subject = s
	}

	msg := &snsstore.Message{
		MessageId:         fmt.Sprintf("%x", time.Now().UnixNano()),
		TopicArn:          topicArn,
		Subject:           subject,
		Message:           message,
		MessageAttributes: make(map[string]*snsstore.MessageAttribute),
	}

	msg.PublishedTimestamp = time.Now().UTC()
	msg.ReceivedTimestamp = time.Now().UTC()

	if err := e.bus.SNSInvoker().StoreMessage(ctx, topicArn+":messages:"+msg.MessageId, msg); err != nil {
		return "", fmt.Errorf("failed to store SNS message: %w", err)
	}

	if e.bus != nil {
		snsEvt := &eventbus.SNSDeliveryEvent{
			TopicARN:  topicArn,
			MessageID: msg.MessageId,
			Message:   message,
			Subject:   subject,
		}
		_, _, evtRegion, _, _ := arnutil.SplitARN(topicArn)
		if evtRegion == "" {
			evtRegion = e.region
		}
		snsEvt.Region = evtRegion
		if err := e.bus.Publish(context.Background(), snsEvt); err != nil {
			logs.Warn("failed to publish SNS event from Step Functions", logs.Err(err))
		}
	}

	result := map[string]interface{}{
		"MessageId": msg.MessageId,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal SNS task result: %w", err)
	}
	return string(resultJSON), nil
}

func (e *Executor) executeEventsTask(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, input string) (string, error) {
	if e.bus == nil || e.bus.EventsInvoker() == nil {
		return "", fmt.Errorf("events invoker not configured")
	}

	resourceParts := strings.Split(state.Resource, ":")
	if len(resourceParts) < 7 {
		return "", fmt.Errorf("invalid Events resource ARN: %s", state.Resource)
	}

	action := resourceParts[6]
	if action == "putEvents" || strings.HasSuffix(action, ".putEvents") {
		return e.executeEventsPutEvents(ctx, execCtx, state, input)
	}

	return "", fmt.Errorf("unsupported Events action: %s", action)
}

func (e *Executor) executeEventsPutEvents(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, input string) (string, error) {
	var inputData map[string]interface{}
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		inputData = map[string]interface{}{"Entries": []interface{}{}}
	}

	_, _, eventsRegion, _, _ := arnutil.SplitARN(state.Resource)
	if eventsRegion == "" {
		eventsRegion = e.region
	}

	var entries []map[string]interface{}
	if e, ok := inputData["Entries"].([]interface{}); ok {
		for _, entry := range e {
			if entryMap, ok := entry.(map[string]interface{}); ok {
				entries = append(entries, entryMap)
			}
		}
	} else if e, ok := inputData["entries"].([]interface{}); ok {
		for _, entry := range e {
			if entryMap, ok := entry.(map[string]interface{}); ok {
				entries = append(entries, entryMap)
			}
		}
	}

	if len(entries) == 0 {
		entries = []map[string]interface{}{
			{
				"Source":       "aws.states",
				"DetailType":   "StepFunctionsExecution",
				"Detail":       input,
				"EventBusName": "default",
			},
		}
	}

	var eventIds []string
	now := time.Now().UTC()

	for _, entry := range entries {
		eventBusName := "default"
		if ebn, ok := entry["EventBusName"].(string); ok && ebn != "" {
			eventBusName = ebn
		} else if ebn, ok := entry["eventBusName"].(string); ok && ebn != "" {
			eventBusName = ebn
		}

		detail := map[string]interface{}{}
		if d, ok := entry["Detail"]; ok {
			if detailMap, ok := d.(map[string]interface{}); ok {
				detail = detailMap
			} else if detailStr, ok := d.(string); ok {
				if jsonErr := json.Unmarshal([]byte(detailStr), &detail); jsonErr != nil {
					detail = map[string]interface{}{"raw": detailStr}
				}
			}
		} else if d, ok := entry["detail"]; ok {
			if detailMap, ok := d.(map[string]interface{}); ok {
				detail = detailMap
			} else if detailStr, ok := d.(string); ok {
				if jsonErr := json.Unmarshal([]byte(detailStr), &detail); jsonErr != nil {
					detail = map[string]interface{}{"raw": detailStr}
				}
			}
		}

		event := &eventsstore.Event{
			ID:           fmt.Sprintf("%x", now.UnixNano()),
			EventBusName: eventBusName,
			Source:       getStringOrEmpty(entry, "Source", "source"),
			DetailType:   getStringOrEmpty(entry, "DetailType", "detailType"),
			Time:         now,
			Region:       eventsRegion,
			Account:      e.accountID,
			Detail:       detail,
		}

		key := fmt.Sprintf("events:%s:%s", eventBusName, event.ID)
		if err := e.bus.EventsInvoker().PutEvent(ctx, key, event); err != nil {
			return "", fmt.Errorf("failed to store event: %w", err)
		}

		if e.bus != nil {
			eventJSON, err := json.Marshal(event)
			if err == nil {
				ebEvt := &eventbus.EventBridgeDeliveryEvent{
					TargetARN: arnutil.NewARNBuilder(e.accountID, eventsRegion).Events().EventBus(eventBusName),
					Input:     eventJSON,
				}
				ebEvt.Region = eventsRegion
				if err := e.bus.Publish(context.Background(), ebEvt); err != nil {
					logs.Warn("failed to publish EventBridge event from Step Functions", logs.Err(err))
				}
			}
		}

		eventIds = append(eventIds, event.ID)
	}

	result := map[string]interface{}{
		"FailedEntryCount": 0,
		"Entries":          eventIds,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Events task result: %w", err)
	}
	return string(resultJSON), nil
}

func getStringOrEmpty(m map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if v, ok := m[key].(string); ok && v != "" {
			return v
		}
	}
	return ""
}

func (e *Executor) executeActivityTask(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, input string, timeoutSeconds, heartbeatSeconds int32) (string, string, error) {
	parts := strings.Split(state.Resource, ":")
	if len(parts) < 7 {
		return "", "", fmt.Errorf("invalid activity ARN: %s", state.Resource)
	}

	activityName := parts[6]
	activityArn := arnutil.NewARNBuilder(e.accountID, e.region).StepFunctions().Activity(activityName)

	_, err := e.store.GetActivity(ctx, activityArn)
	if err != nil {
		return "", "", fmt.Errorf("activity not found: %s", activityArn)
	}

	task := &sfnstore.ActivityTask{
		ActivityArn:  activityArn,
		ExecutionArn: execCtx.Execution.ExecutionArn,
		Input:        input,
	}

	if err := e.store.CreateActivityTask(task); err != nil {
		return "", "", fmt.Errorf("failed to create activity task: %w", err)
	}

	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn: execCtx.Execution.ExecutionArn,
		EventId:      execCtx.nextEventId(),
		Type:         "ActivityTaskScheduled",
		ActivityTaskScheduledEventDetails: &sfnstore.ActivityTaskScheduledEventDetails{
			Resource:         state.Resource,
			Input:            input,
			TaskToken:        task.TaskToken,
			HeartbeatSeconds: heartbeatSeconds,
		},
	})

	timeout := time.Duration(timeoutSeconds) * time.Second
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	result, err := e.store.WaitForTaskResult(ctx, task.TaskToken, timeout)
	if err != nil {
		if err == sfnstore.ErrTaskTimeout {
			e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
				ExecutionArn: execCtx.Execution.ExecutionArn,
				EventId:      execCtx.nextEventId(),
				Type:         "ActivityTaskTimedOut",
				ActivityTaskTimedOutEventDetails: &sfnstore.ActivityTaskTimedOutEventDetails{
					Error: "States.Timeout",
					Cause: "Task timed out",
				},
			})
			return "", "", fmt.Errorf("States.Timeout: Task timed out")
		}
		return "", "", err
	}

	if result.Error != nil {
		e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
			ExecutionArn: execCtx.Execution.ExecutionArn,
			EventId:      execCtx.nextEventId(),
			Type:         "ActivityTaskFailed",
			ActivityTaskFailedEventDetails: &sfnstore.ActivityTaskFailedEventDetails{
				Error: result.Error.Error(),
				Cause: "",
			},
		})
		return "", "", result.Error
	}

	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn: execCtx.Execution.ExecutionArn,
		EventId:      execCtx.nextEventId(),
		Type:         "ActivityTaskSucceeded",
		ActivityTaskSucceededEventDetails: &sfnstore.ActivityTaskSucceededEventDetails{
			Output: result.Output,
		},
	})

	return result.Output, state.Next, nil
}

func (e *Executor) executeDynamoDBTask(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, input string) (string, error) {
	if e.bus == nil || e.bus.DynamoDBInvoker() == nil {
		return "", fmt.Errorf("DynamoDB invoker not configured")
	}

	resourceParts := strings.Split(state.Resource, ":")
	if len(resourceParts) < 7 {
		return "", fmt.Errorf("invalid DynamoDB resource ARN: %s", state.Resource)
	}

	action := resourceParts[6]

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
		return "", fmt.Errorf("Key is required for getItem")
	}
	key := awsAttrValuesToPlain(keyRaw)

	item, err := e.bus.DynamoDBInvoker().GetItem(ctx, region, tableName, key)
	if err != nil {
		return "", fmt.Errorf("DynamoDB GetItem failed: %w", err)
	}

	result := map[string]interface{}{"Item": plainMapToAWSAttrValues(item)}
	resultJSON, _ := json.Marshal(result)
	return string(resultJSON), nil
}

func (e *Executor) executeDynamoDBPutItem(ctx context.Context, region, tableName string, inputData map[string]interface{}) (string, error) {
	itemRaw, ok := inputData["Item"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Item is required for putItem")
	}
	plainItem := awsAttrValuesToPlain(itemRaw)

	key := extractKeyFromItem(plainItem)
	attrs := removeKeys(plainItem, key)

	result, err := e.bus.DynamoDBInvoker().PutItem(ctx, region, tableName, key, attrs)
	if err != nil {
		return "", fmt.Errorf("DynamoDB PutItem failed: %w", err)
	}

	resultJSON, _ := json.Marshal(map[string]interface{}{"Attributes": plainMapToAWSAttrValues(result)})
	return string(resultJSON), nil
}

func (e *Executor) executeDynamoDBDeleteItem(ctx context.Context, region, tableName string, inputData map[string]interface{}) (string, error) {
	keyRaw, ok := inputData["Key"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Key is required for deleteItem")
	}
	key := awsAttrValuesToPlain(keyRaw)

	if err := e.bus.DynamoDBInvoker().DeleteItem(ctx, region, tableName, key); err != nil {
		return "", fmt.Errorf("DynamoDB DeleteItem failed: %w", err)
	}

	return "{}", nil
}

func (e *Executor) executeDynamoDBUpdateItem(ctx context.Context, region, tableName string, inputData map[string]interface{}) (string, error) {
	keyRaw, ok := inputData["Key"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Key is required for updateItem")
	}
	key := awsAttrValuesToPlain(keyRaw)

	attrValues := map[string]interface{}{}
	if av, ok := inputData["AttributeValues"].(map[string]interface{}); ok {
		attrValues = awsAttrValuesToPlain(av)
	} else if av, ok := inputData["ExpressionAttributeValues"].(map[string]interface{}); ok {
		attrValues = awsAttrValuesToPlain(av)
	}

	if err := e.bus.DynamoDBInvoker().UpdateItem(ctx, region, tableName, key, attrValues); err != nil {
		return "", fmt.Errorf("DynamoDB UpdateItem failed: %w", err)
	}

	return "{}", nil
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
	case []string:
		ss := make([]interface{}, len(val))
		for i, s := range val {
			ss[i] = s
		}
		return map[string]interface{}{"SS": ss}
	case []interface{}:
		ss := make([]interface{}, len(val))
		for i, s := range val {
			ss[i] = fmt.Sprintf("%v", s)
		}
		return map[string]interface{}{"SS": ss}
	default:
		return map[string]interface{}{"S": fmt.Sprintf("%v", val)}
	}
}

func extractKeyFromItem(item map[string]interface{}) map[string]interface{} {
	key := make(map[string]interface{})
	for k, v := range item {
		key[k] = v
	}
	return key
}

func removeKeys(item map[string]interface{}, keys map[string]interface{}) map[string]interface{} {
	attrs := make(map[string]interface{})
	for k, v := range item {
		if _, isKey := keys[k]; !isKey {
			attrs[k] = v
		}
	}
	return attrs
}

func getStr(m map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if v, ok := m[key].(string); ok && v != "" {
			return v
		}
	}
	return ""
}

func (e *Executor) isActivityResource(resource string) bool {
	if strings.HasPrefix(resource, "arn:aws:states:::activity:") {
		return true
	}
	if strings.HasPrefix(resource, "arn:aws:states:") && strings.Contains(resource, ":activity:") {
		return true
	}
	return false
}
