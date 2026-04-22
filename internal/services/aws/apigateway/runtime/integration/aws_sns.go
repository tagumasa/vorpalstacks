package integration

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"vorpalstacks/internal/common/endpoint"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

var (
	sqsPathRegex   = regexp.MustCompile(`sqs:path/[^/]+/([^/]+)`)
	sqsActionRegex = regexp.MustCompile(`sqs:action/[^/]+/([^/]+)`)
)

// SQSIntegrationRequest represents a request for SQS integration in API Gateway.
type SQSIntegrationRequest struct {
	Action      string `json:"Action"`
	MessageBody string `json:"MessageBody"`
	QueueUrl    string `json:"QueueUrl,omitempty"`
}

func isSQSURI(uri string) bool {
	return strings.Contains(uri, ":sqs:")
}

func (e *AWSExecutor) executeSQS(ctx context.Context, req *IntegrationRequest) (*IntegrationResponse, error) {
	if e.bus == nil || e.bus.SQSInvoker() == nil {
		return nil, &IntegrationError{
			Message:  "SQS store not configured",
			Type:     "InternalServerError",
			HTTPCode: 500,
		}
	}

	queueName, err := extractQueueNameFromURI(req.URI)
	if err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Invalid SQS URI: %v", err),
			Type:     "BadRequestException",
			HTTPCode: 400,
		}
	}

	queueURL := endpoint.SQSQueueURL(e.accountID, queueName)

	action := req.QueryParams["Action"]
	if action == "" {
		action = req.Headers["Action"]
	}
	if action == "" {
		action = "SendMessage"
	}

	switch action {
	case "SendMessage":
		return e.executeSQSSendMessage(ctx, queueURL, req)
	case "ReceiveMessage":
		return e.executeSQSReceiveMessage(ctx, queueURL, req)
	default:
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Unsupported SQS action: %s", action),
			Type:     "BadRequestException",
			HTTPCode: 400,
		}
	}
}

func (e *AWSExecutor) executeSQSSendMessage(ctx context.Context, queueURL string, req *IntegrationRequest) (*IntegrationResponse, error) {
	messageBody := string(req.Body)
	if messageBody == "" {
		messageBody = req.Headers["MessageBody"]
	}

	messageAttributes := convertToSQSInvokerAttrs(extractSQSMessageAttributes(req.Headers, req.QueryParams, req.Body))

	messageID, md5OfBody, err := e.bus.SQSInvoker().SendMessage(ctx, queueURL, messageBody, 0, messageAttributes)
	if err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Failed to send SQS message: %v", err),
			Type:     "InternalServerError",
			HTTPCode: 500,
		}
	}

	response := map[string]interface{}{
		"SendMessageResponse": map[string]interface{}{
			"SendMessageResult": map[string]string{
				"MD5OfMessageBody":       md5OfBody,
				"MessageId":              messageID,
				"MD5OfMessageAttributes": "",
			},
			"ResponseMetadata": map[string]string{
				"RequestId": fmt.Sprintf("%x", time.Now().UnixNano()),
			},
		},
	}

	responseJSON, _ := jsonMarshal(response)
	return &IntegrationResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            responseJSON,
		IsBase64Encoded: false,
	}, nil
}

func (e *AWSExecutor) executeSQSReceiveMessage(ctx context.Context, queueURL string, req *IntegrationRequest) (*IntegrationResponse, error) {
	maxMessages := int32(1)
	if val := req.Headers["MaxNumberOfMessages"]; val != "" {
		_, _ = fmt.Sscanf(val, "%d", &maxMessages)
	}

	waitTime := int32(0)
	if val := req.Headers["WaitTimeSeconds"]; val != "" {
		_, _ = fmt.Sscanf(val, "%d", &waitTime)
	}

	visibilityTimeout := int32(30)
	if val := req.Headers["VisibilityTimeout"]; val != "" {
		_, _ = fmt.Sscanf(val, "%d", &visibilityTimeout)
	}

	messages, err := e.bus.SQSInvoker().ReceiveMessage(ctx, queueURL, maxMessages, &visibilityTimeout, waitTime)
	if err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Failed to receive SQS messages: %v", err),
			Type:     "InternalServerError",
			HTTPCode: 500,
		}
	}

	response := map[string]interface{}{
		"ReceiveMessageResponse": map[string]interface{}{
			"ReceiveMessageResult": map[string]interface{}{
				"messages": messages,
			},
			"ResponseMetadata": map[string]string{
				"RequestId": fmt.Sprintf("%x", time.Now().UnixNano()),
			},
		},
	}

	responseJSON, _ := jsonMarshal(response)
	return &IntegrationResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            responseJSON,
		IsBase64Encoded: false,
	}, nil
}

func extractQueueNameFromURI(uri string) (string, error) {
	matches := sqsPathRegex.FindStringSubmatch(uri)
	if len(matches) < 2 {
		matches = sqsActionRegex.FindStringSubmatch(uri)
		if len(matches) < 2 {
			return "", fmt.Errorf("invalid SQS URI format: %s", uri)
		}
	}
	return matches[1], nil
}

func extractSQSMessageAttributes(headers, queryParams map[string]string, body []byte) map[string]struct {
	StringValue string
	DataType    string
} {
	attrs := make(map[string]struct {
		StringValue string
		DataType    string
	})
	attrNames := make(map[int]string)
	attrValues := make(map[int]struct {
		StringValue string
		DataType    string
	})

	extractFromMap := func(m map[string]string) {
		for k, v := range m {
			if matches := regexp.MustCompile(`^MessageAttribute\.(\d+)\.Name$`).FindStringSubmatch(k); matches != nil {
				var idx int
				if _, err := fmt.Sscanf(matches[1], "%d", &idx); err != nil {
					continue
				}
				attrNames[idx] = v
			}
			if matches := regexp.MustCompile(`^MessageAttribute\.(\d+)\.Value\.StringValue$`).FindStringSubmatch(k); matches != nil {
				var idx int
				if _, err := fmt.Sscanf(matches[1], "%d", &idx); err != nil {
					continue
				}
				val := attrValues[idx]
				val.StringValue = v
				attrValues[idx] = val
			}
			if matches := regexp.MustCompile(`^MessageAttribute\.(\d+)\.Value\.DataType$`).FindStringSubmatch(k); matches != nil {
				var idx int
				if _, err := fmt.Sscanf(matches[1], "%d", &idx); err != nil {
					continue
				}
				val := attrValues[idx]
				val.DataType = v
				attrValues[idx] = val
			}
		}
	}

	extractFromMap(headers)
	extractFromMap(queryParams)

	if len(body) > 0 {
		contentType := headers["Content-Type"]
		if strings.Contains(contentType, "application/x-www-form-urlencoded") || contentType == "" {
			bodyParams := parseFormData(string(body))
			extractFromMap(bodyParams)
		}
	}

	for idx, name := range attrNames {
		val := attrValues[idx]
		if name != "" && val.DataType != "" {
			attrs[name] = val
		}
	}

	return attrs
}

func convertToSQSInvokerAttrs(attrs map[string]struct {
	StringValue string
	DataType    string
}) map[string]string {
	if attrs == nil {
		return nil
	}
	result := make(map[string]string, len(attrs))
	for k, v := range attrs {
		if v.DataType == "String" {
			result[k] = v.StringValue
		}
	}
	return result
}

func parseFormData(data string) map[string]string {
	result := make(map[string]string)
	for _, pair := range strings.Split(data, "&") {
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

type snsNotification struct {
	MessageId          string    `json:"messageId"`
	TopicArn           string    `json:"topicArn"`
	Subject            string    `json:"subject,omitempty"`
	Message            string    `json:"message"`
	PublishedTimestamp time.Time `json:"publishedTimestamp"`
}

var (
	snsPathRegex   = regexp.MustCompile(`sns:path/[^/]+/([^/]+)`)
	snsActionRegex = regexp.MustCompile(`sns:action/[^/]+/([^/]+)`)
)

func isSNSURI(uri string) bool {
	return strings.Contains(uri, ":sns:")
}

func (e *AWSExecutor) executeSNS(ctx context.Context, req *IntegrationRequest) (*IntegrationResponse, error) {
	if e.bus == nil || e.bus.SNSInvoker() == nil {
		return nil, &IntegrationError{
			Message:  "SNS store not configured",
			Type:     "InternalServerError",
			HTTPCode: 500,
		}
	}

	topicName, topicRegion, err := extractTopicFromURI(req.URI)
	if err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Invalid SNS URI: %v", err),
			Type:     "BadRequestException",
			HTTPCode: 400,
		}
	}

	if topicRegion == "" {
		topicRegion = e.region
	}
	topicArn := arnutil.NewARNBuilder(e.accountID, topicRegion).SNS().Topic(topicName)

	action := req.QueryParams["Action"]
	if action == "" {
		action = req.Headers["Action"]
	}
	if action == "" && len(req.Body) > 0 {
		if bodyParams, err := url.ParseQuery(string(req.Body)); err == nil {
			action = bodyParams.Get("Action")
		}
	}
	if action == "" {
		action = "Publish"
	}

	switch action {
	case "Publish":
		return e.executeSNSPublish(ctx, topicArn, req)
	default:
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Unsupported SNS action: %s", action),
			Type:     "BadRequestException",
			HTTPCode: 400,
		}
	}
}

func (e *AWSExecutor) executeSNSPublish(ctx context.Context, topicArn string, req *IntegrationRequest) (*IntegrationResponse, error) {
	_, err := e.bus.SNSInvoker().GetTopic(ctx, topicArn)
	if err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("SNS topic not found: %s", topicArn),
			Type:     "NotFoundException",
			HTTPCode: 404,
		}
	}

	message := string(req.Body)
	if message == "" {
		message = req.Headers["Message"]
	}

	messageID := fmt.Sprintf("%x", time.Now().UnixNano())

	now := time.Now().UTC()
	notification := &snsNotification{
		MessageId:          messageID,
		TopicArn:           topicArn,
		Subject:            req.Headers["Subject"],
		Message:            message,
		PublishedTimestamp: now,
	}

	if err := e.bus.SNSInvoker().StoreMessage(ctx, topicArn+":messages:"+messageID, notification); err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Failed to store SNS message: %v", err),
			Type:     "InternalServerError",
			HTTPCode: 500,
		}
	}

	if e.bus != nil {
		_, _, snsRegion, _, _ := arnutil.SplitARN(topicArn)
		if snsRegion == "" {
			snsRegion = e.region
		}
		snsEvt := &eventbus.SNSDeliveryEvent{
			TopicARN:  topicArn,
			MessageID: messageID,
			Message:   message,
			Subject:   req.Headers["Subject"],
		}
		snsEvt.Region = snsRegion
		if err := e.bus.Publish(ctx, snsEvt); err != nil {
			logs.Warn("Failed to publish SNS delivery event to event bus; message is stored but subscribers may not be notified",
				logs.String("topicArn", topicArn),
				logs.String("messageId", messageID),
				logs.Err(err))
		}
	}

	response := map[string]interface{}{
		"PublishResponse": map[string]interface{}{
			"PublishResult": map[string]string{
				"MessageId": messageID,
			},
			"ResponseMetadata": map[string]string{
				"RequestId": fmt.Sprintf("%x", time.Now().UnixNano()),
			},
		},
	}

	responseJSON, _ := jsonMarshal(response)
	return &IntegrationResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            responseJSON,
		IsBase64Encoded: false,
	}, nil
}

func extractTopicFromURI(uri string) (topicName, region string, err error) {
	matches := snsPathRegex.FindStringSubmatch(uri)
	if len(matches) < 2 {
		matches = snsActionRegex.FindStringSubmatch(uri)
		if len(matches) < 2 {
			return "", "", fmt.Errorf("invalid SNS URI format: %s", uri)
		}
	}
	return matches[1], "", nil
}

func (e *AWSExecutor) deliverToSNSSubscribers(topicArn string, notification *snsNotification) {
	subs, err := e.bus.SNSInvoker().ListSubscriptionsByTopic(context.Background(), topicArn)
	if err != nil {
		logs.Warn("failed to list SNS subscribers",
			logs.String("topicArn", topicArn),
			logs.String("messageId", notification.MessageId),
			logs.Err(err))
		return
	}

	for _, sub := range subs {
		if sub.PendingConfirmation {
			continue
		}

		switch {
		case arnutil.IsSQSARN(sub.Endpoint):
			e.deliverToSQSSubscriber(sub.Endpoint, notification)
		case arnutil.IsLambdaARN(sub.Endpoint):
			e.deliverToLambdaSubscriber(sub.Endpoint, notification)
		}
	}
}

func (e *AWSExecutor) deliverToSQSSubscriber(queueArn string, notification *snsNotification) {
	if e.bus == nil || e.bus.SQSInvoker() == nil {
		return
	}

	_, _, queueRegion, _, _ := arnutil.SplitARN(queueArn)
	queueName := arnutil.ExtractQueueNameFromARN(queueArn)
	if queueRegion == "" {
		queueRegion = e.region
	}
	queueURL := fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", queueRegion, e.accountID, queueName)

	body, err := jsonMarshal(map[string]interface{}{
		"Type":      "Notification",
		"MessageId": notification.MessageId,
		"TopicArn":  notification.TopicArn,
		"Subject":   notification.Subject,
		"Message":   notification.Message,
		"Timestamp": notification.PublishedTimestamp.Format(time.RFC3339),
	})
	if err != nil {
		return
	}

	if _, _, err := e.bus.SQSInvoker().SendMessage(context.Background(), queueURL, string(body), 0, nil); err != nil {
		logs.Error("Failed to send SQS message to subscriber", logs.String("queue", queueURL), logs.Err(err))
	}
}

func (e *AWSExecutor) deliverToLambdaSubscriber(functionArn string, notification *snsNotification) {
	if e.bus == nil || e.bus.LambdaInvoker() == nil {
		return
	}

	event := map[string]interface{}{
		"Records": []interface{}{
			map[string]interface{}{
				"EventSource":          "aws:sns",
				"EventVersion":         "1.0",
				"EventSubscriptionArn": notification.TopicArn,
				"Sns": map[string]interface{}{
					"Type":      "Notification",
					"MessageId": notification.MessageId,
					"TopicArn":  notification.TopicArn,
					"Subject":   notification.Subject,
					"Message":   notification.Message,
					"Timestamp": notification.PublishedTimestamp.Format(time.RFC3339),
				},
			},
		},
	}

	eventJSON, _ := jsonMarshal(event)
	if _, _, err := e.bus.LambdaInvoker().InvokeForGateway(context.Background(), functionArn, eventJSON); err != nil {
		logs.Error("apigateway: SNS-to-Lambda invocation failed", logs.String("function", functionArn), logs.Err(err))
	}
}
