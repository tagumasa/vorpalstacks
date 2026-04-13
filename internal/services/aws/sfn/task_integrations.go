package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/common/endpoint"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	snsstore "vorpalstacks/internal/store/aws/sns"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

func (e *Executor) executeLambdaTask(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, input string) (string, error) {
	if e.lambdaInvoker == nil {
		return "", fmt.Errorf("lambda invoker not configured")
	}

	parts := strings.Split(state.Resource, ":")
	if len(parts) < 7 {
		return "", fmt.Errorf("invalid Lambda ARN: %s", state.Resource)
	}

	statusCode, payload, err := e.lambdaInvoker.InvokeForGateway(ctx, state.Resource, []byte(input))
	if err != nil {
		return "", fmt.Errorf("failed to invoke Lambda function: %w", err)
	}

	if statusCode != 200 {
		return "", fmt.Errorf("lambda function returned status code: %d, error: %s", statusCode, string(payload))
	}

	return string(payload), nil
}

func (e *Executor) executeSQSTask(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, input string) (string, error) {
	if e.sqsStore == nil {
		return "", fmt.Errorf("SQS store not configured")
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

		if q, qErr := e.sqsStore.GetQueueByName(queueName); qErr == nil {
			queueURL = q.URL
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

	msg := &sqsstore.Message{
		Body:              messageBody,
		MessageAttributes: make(map[string]*sqsstore.MessageAttributeValue),
		Attributes:        make(map[string]string),
	}

	sentMsg, err := e.sqsStore.SendMessage(queueURL, msg)
	if err != nil {
		return "", fmt.Errorf("failed to send SQS message: %w", err)
	}

	result := map[string]interface{}{
		"MessageId":        sentMsg.ID,
		"MD5OfMessageBody": sentMsg.MD5OfBody,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal SQS task result: %w", err)
	}
	return string(resultJSON), nil
}

func (e *Executor) executeSNSTask(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.TaskState, input string) (string, error) {
	if e.snsStore == nil {
		return "", fmt.Errorf("SNS store not configured")
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

		_, _, topicRegion, _, _ := svcarn.SplitARN(state.Resource)
		if topicRegion == "" {
			topicRegion = e.region
		}
		topicArn = fmt.Sprintf("arn:aws:sns:%s:%s:%s", topicRegion, e.accountID, topicName)
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

	if err := e.snsStore.Put(topicArn+":messages:"+msg.MessageId, msg); err != nil {
		return "", fmt.Errorf("failed to store SNS message: %w", err)
	}

	if e.bus != nil {
		snsEvt := &eventbus.SNSDeliveryEvent{
			TopicARN:  topicArn,
			MessageID: msg.MessageId,
			Message:   message,
			Subject:   subject,
		}
		_, _, evtRegion, _, _ := svcarn.SplitARN(topicArn)
		if evtRegion == "" {
			evtRegion = e.region
		}
		snsEvt.Region = evtRegion
		e.bus.Publish(context.Background(), snsEvt)
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
	if e.eventsStore == nil {
		return "", fmt.Errorf("events store not configured")
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

	_, _, eventsRegion, _, _ := svcarn.SplitARN(state.Resource)
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
		if err := e.eventsStore.Put(key, event); err != nil {
			return "", fmt.Errorf("failed to store event: %w", err)
		}

		if e.bus != nil {
			eventJSON, err := json.Marshal(event)
			if err == nil {
				ebEvt := &eventbus.EventBridgeDeliveryEvent{
					TargetARN: fmt.Sprintf("arn:aws:events:%s:%s:event-bus/%s", eventsRegion, e.accountID, eventBusName),
					Input:     eventJSON,
				}
				ebEvt.Region = eventsRegion
				e.bus.Publish(context.Background(), ebEvt)
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
	activityArn := fmt.Sprintf("arn:aws:states:%s:%s:activity:%s", e.region, e.accountID, activityName)

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

func (e *Executor) isActivityResource(resource string) bool {
	if strings.HasPrefix(resource, "arn:aws:states:::activity:") {
		return true
	}
	if strings.HasPrefix(resource, "arn:aws:states:") && strings.Contains(resource, ":activity:") {
		return true
	}
	return false
}
