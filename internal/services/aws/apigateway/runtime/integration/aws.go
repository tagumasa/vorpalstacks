// Package integration provides API Gateway integration functionality for vorpalstacks.
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common"
	"vorpalstacks/internal/common/endpoint"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	sns "vorpalstacks/internal/store/aws/sns"
	sqs "vorpalstacks/internal/store/aws/sqs"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/pkg/vtl"
)

var (
	lambdaFunctionRegex = regexp.MustCompile(`functions/(.+?)/invocations`)
	sqsPathRegex        = regexp.MustCompile(`sqs:path/[^/]+/([^/]+)`)
	sqsActionRegex      = regexp.MustCompile(`sqs:action/[^/]+/([^/]+)`)
	snsPathRegex        = regexp.MustCompile(`sns:path/[^/]+/([^/]+)`)
	snsActionRegex      = regexp.MustCompile(`sns:action/[^/]+/([^/]+)`)

	snsDeliverySem = make(chan struct{}, 10)
)

// AWSExecutor executes AWS integration requests for API Gateway.
type AWSExecutor struct {
	lambdaInvoker common.LambdaInvoker
	sqsStore      sqs.SQSStoreInterface
	snsStore      sns.SNSStoreInterface
	accountID     string
	region        string
	bus           eventbus.Bus
	deliveryWg    *sync.WaitGroup
}

// NewAWSExecutor creates a new AWSExecutor with the given Lambda invoker.
func NewAWSExecutor(lambdaInvoker common.LambdaInvoker) *AWSExecutor {
	return &AWSExecutor{
		lambdaInvoker: lambdaInvoker,
	}
}

// NewAWSExecutorWithStores creates a new AWSExecutor with the given Lambda invoker and store dependencies.
func NewAWSExecutorWithStores(lambdaInvoker common.LambdaInvoker, sqsStore sqs.SQSStoreInterface, snsStore sns.SNSStoreInterface, accountID, region string, bus eventbus.Bus, deliveryWg *sync.WaitGroup) *AWSExecutor {
	return &AWSExecutor{
		lambdaInvoker: lambdaInvoker,
		sqsStore:      sqsStore,
		snsStore:      snsStore,
		accountID:     accountID,
		region:        region,
		bus:           bus,
		deliveryWg:    deliveryWg,
	}
}

// Execute executes an AWS integration request.
func (e *AWSExecutor) Execute(ctx context.Context, req *IntegrationRequest) (*IntegrationResponse, error) {
	if isLambdaURI(req.URI) {
		return e.executeLambda(ctx, req)
	}

	if isSQSURI(req.URI) {
		return e.executeSQS(ctx, req)
	}

	if isSNSURI(req.URI) {
		return e.executeSNS(ctx, req)
	}

	return nil, &IntegrationError{
		Message:  "AWS integration not yet implemented for: " + req.URI,
		Type:     "InternalServerError",
		HTTPCode: http.StatusNotImplemented,
	}
}

func (e *AWSExecutor) executeLambda(ctx context.Context, req *IntegrationRequest) (*IntegrationResponse, error) {
	if e.lambdaInvoker == nil {
		return nil, &IntegrationError{
			Message:  "Lambda client not configured",
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
		}
	}

	functionRef, err := extractFunctionRefFromURI(req.URI)
	if err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Invalid Lambda URI: %v", err),
			Type:     "BadRequestException",
			HTTPCode: http.StatusBadRequest,
		}
	}

	isProxy := req.IntegrationType == "AWS_PROXY"

	var eventJSON []byte
	if isProxy {
		event := e.buildLambdaProxyEvent(req)
		eventJSON, err = json.Marshal(event)
		if err != nil {
			return nil, &IntegrationError{
				Message:  fmt.Sprintf("Failed to marshal Lambda event: %v", err),
				Type:     "InternalServerError",
				HTTPCode: http.StatusInternalServerError,
			}
		}
	} else {
		contentType := req.Headers["Content-Type"]
		tmpl := req.RequestTemplates[contentType]
		if tmpl == "" {
			tmpl = req.RequestTemplates["application/json"]
		}
		if tmpl != "" {
			eventJSON, err = applyRequestTemplate(tmpl, req)
			if err != nil {
				return nil, &IntegrationError{
					Message:  fmt.Sprintf("Failed to apply request template: %v", err),
					Type:     "InternalServerError",
					HTTPCode: http.StatusInternalServerError,
				}
			}
		} else if len(req.Body) > 0 {
			eventJSON = req.Body
		} else {
			eventJSON = []byte("{}")
		}
	}

	statusCode, payload, err := e.lambdaInvoker.InvokeForGateway(ctx, functionRef, eventJSON)
	if err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Lambda invocation failed: %v", err),
			Type:     "IntegrationFailure",
			HTTPCode: http.StatusBadGateway,
		}
	}

	if isProxy {
		lambdaResp, err := parseLambdaResponse(payload)
		if err != nil {
			return nil, &IntegrationError{
				Message:  fmt.Sprintf("Failed to parse Lambda response: %v", err),
				Type:     "InternalServerError",
				HTTPCode: http.StatusInternalServerError,
			}
		}
		return lambdaResp, nil
	}

	resp := &IntegrationResponse{
		StatusCode:      int(statusCode),
		Headers:         make(map[string]string),
		Body:            payload,
		IsBase64Encoded: false,
	}

	if req.IntegrationResponses != nil {
		respConfig := matchIntegrationResponse(req.IntegrationResponses, string(payload), int(statusCode))
		if respConfig != nil {
			if respConfig.ResponseHeaders != nil {
				for k, v := range respConfig.ResponseHeaders {
					resp.Headers[k] = v
				}
			}

			if respConfig.ResponseTemplates != nil {
				contentType := "application/json"
				if tmpl, ok := respConfig.ResponseTemplates[contentType]; ok && tmpl != "" {
					transformed, err := applyResponseTemplate(tmpl, payload, req)
					if err != nil {
						return nil, &IntegrationError{
							Message:  fmt.Sprintf("Failed to apply response template: %v", err),
							Type:     "InternalServerError",
							HTTPCode: http.StatusInternalServerError,
						}
					}
					resp.Body = transformed
				}
			}

			if respConfig.StatusCode != "" {
				_, _ = fmt.Sscanf(respConfig.StatusCode, "%d", &resp.StatusCode)
			}
		}
	}

	return resp, nil
}

func (e *AWSExecutor) buildLambdaProxyEvent(req *IntegrationRequest) *LambdaProxyEvent {
	now := time.Now()
	requestID := fmt.Sprintf("%x", now.UnixNano())

	return &LambdaProxyEvent{
		Resource:                        req.Path,
		Path:                            req.Path,
		HTTPMethod:                      req.Method,
		Headers:                         req.Headers,
		MultiValueHeaders:               req.MultiValueHeaders,
		QueryStringParameters:           req.QueryParams,
		MultiValueQueryStringParameters: req.MultiValueQueryStringParams,
		PathParameters:                  req.PathParams,
		StageVariables:                  req.StageVariables,
		RequestContext: map[string]interface{}{
			"accountID":    e.accountID,
			"apiId":        req.RestApiId,
			"requestId":    requestID,
			"resourcePath": req.Path,
			"stage":        req.StageName,
			"region":       e.region,
			"http": map[string]interface{}{
				"method":    req.Method,
				"path":      req.Path,
				"protocol":  "HTTP/1.1",
				"sourceIp":  req.SourceIP,
				"userAgent": req.Headers["User-Agent"],
			},
			"requestTime":      now.Format("02/Jan/2006:15:04:05 -0700"),
			"requestTimeEpoch": now.UnixMilli(),
		},
		Body:            string(req.Body),
		IsBase64Encoded: false,
	}
}

// applyRequestTemplate applies a basic subset of VTL substitutions to the given
// request template string. Supports $input.json('$'), $input.body,
// $input.params('name'), and $context.stage/$context.apiId/$context.requestId.
func applyRequestTemplate(tmpl string, req *IntegrationRequest) ([]byte, error) {
	engine := vtl.NewEngine()

	engine.GetContext().Body = string(req.Body)
	if len(req.Body) > 0 {
		var bodyObj interface{}
		if err := json.Unmarshal(req.Body, &bodyObj); err == nil {
			engine.GetContext().JSONBody = bodyObj
		}
	}

	params := make(map[string]string)
	for k, v := range req.PathParams {
		params[k] = v
	}
	for k, v := range req.QueryParams {
		params[k] = v
	}
	for k, v := range req.Headers {
		params[k] = v
	}
	engine.GetContext().Params = params

	engine.GetContext().Context = map[string]interface{}{
		"stage":     req.StageName,
		"apiId":     req.RestApiId,
		"requestId": fmt.Sprintf("%x", time.Now().UnixNano()),
	}

	engine.GetContext().StageVars = req.StageVariables

	result, err := engine.Transform(tmpl)
	if err != nil {
		return nil, err
	}

	return []byte(result), nil
}

func matchIntegrationResponse(responses map[string]*IntegrationResponseConfig, responseBody string, statusCode int) *IntegrationResponseConfig {
	var defaultResp *IntegrationResponseConfig

	for _, resp := range responses {
		if resp.SelectionPattern == "" {
			if defaultResp == nil {
				defaultResp = resp
			}
			continue
		}

		matched, err := regexp.MatchString(resp.SelectionPattern, responseBody)
		if err == nil && matched {
			return resp
		}
	}

	if defaultResp != nil {
		return defaultResp
	}

	statusStr := fmt.Sprintf("%d", statusCode)
	return responses[statusStr]
}

func applyResponseTemplate(tmpl string, responseBody []byte, req *IntegrationRequest) ([]byte, error) {
	engine := vtl.NewEngine()

	engine.GetContext().Body = string(responseBody)
	if len(responseBody) > 0 {
		var bodyObj interface{}
		if err := json.Unmarshal(responseBody, &bodyObj); err == nil {
			engine.GetContext().JSONBody = bodyObj
		}
	}

	params := make(map[string]string)
	for k, v := range req.PathParams {
		params[k] = v
	}
	for k, v := range req.QueryParams {
		params[k] = v
	}
	for k, v := range req.Headers {
		params[k] = v
	}
	engine.GetContext().Params = params

	engine.GetContext().Context = map[string]interface{}{
		"stage":     req.StageName,
		"apiId":     req.RestApiId,
		"requestId": fmt.Sprintf("%x", time.Now().UnixNano()),
	}

	engine.GetContext().StageVars = req.StageVariables

	result, err := engine.Transform(tmpl)
	if err != nil {
		return nil, err
	}

	return []byte(result), nil
}

func extractFunctionRefFromURI(uri string) (string, error) {
	matches := lambdaFunctionRegex.FindStringSubmatch(uri)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid Lambda URI format")
	}

	return matches[1], nil
}

func isLambdaURI(uri string) bool {
	return strings.Contains(uri, "lambda:")
}

// LambdaProxyResponse represents the response from a Lambda proxy integration.
type LambdaProxyResponse struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

// LambdaProxyEvent represents the event structure passed to a Lambda function
// when using API Gateway Lambda proxy integration.
type LambdaProxyEvent struct {
	Resource                        string                 `json:"resource"`
	Path                            string                 `json:"path"`
	HTTPMethod                      string                 `json:"httpMethod"`
	Headers                         map[string]string      `json:"headers"`
	MultiValueHeaders               map[string][]string    `json:"multiValueHeaders"`
	QueryStringParameters           map[string]string      `json:"queryStringParameters"`
	MultiValueQueryStringParameters map[string][]string    `json:"multiValueQueryStringParameters"`
	PathParameters                  map[string]string      `json:"pathParameters"`
	StageVariables                  map[string]string      `json:"stageVariables"`
	RequestContext                  map[string]interface{} `json:"requestContext"`
	Body                            string                 `json:"body"`
	IsBase64Encoded                 bool                   `json:"isBase64Encoded"`
}

func parseLambdaResponse(body []byte) (*IntegrationResponse, error) {
	var lambdaResp LambdaProxyResponse
	if err := json.Unmarshal(body, &lambdaResp); err != nil {
		return nil, fmt.Errorf("failed to parse Lambda response: %w", err)
	}

	return &IntegrationResponse{
		StatusCode:      lambdaResp.StatusCode,
		Headers:         lambdaResp.Headers,
		Body:            []byte(lambdaResp.Body),
		IsBase64Encoded: lambdaResp.IsBase64Encoded,
	}, nil
}

func isSQSURI(uri string) bool {
	return strings.Contains(uri, ":sqs:")
}

func isSNSURI(uri string) bool {
	return strings.Contains(uri, ":sns:")
}

// SQSIntegrationRequest represents a request for SQS integration in API Gateway.
type SQSIntegrationRequest struct {
	Action      string `json:"Action"`
	MessageBody string `json:"MessageBody"`
	QueueUrl    string `json:"QueueUrl,omitempty"`
}

func (e *AWSExecutor) executeSQS(ctx context.Context, req *IntegrationRequest) (*IntegrationResponse, error) {
	if e.sqsStore == nil {
		return nil, &IntegrationError{
			Message:  "SQS store not configured",
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
		}
	}

	queueName, err := extractQueueNameFromURI(req.URI)
	if err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Invalid SQS URI: %v", err),
			Type:     "BadRequestException",
			HTTPCode: http.StatusBadRequest,
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
			HTTPCode: http.StatusBadRequest,
		}
	}
}

func (e *AWSExecutor) executeSQSSendMessage(ctx context.Context, queueURL string, req *IntegrationRequest) (*IntegrationResponse, error) {
	message := &sqs.Message{
		Body: string(req.Body),
	}

	if message.Body == "" {
		message.Body = req.Headers["MessageBody"]
	}

	message.MessageAttributes = extractSQSMessageAttributes(req.Headers, req.QueryParams, req.Body)

	_, err := e.sqsStore.SendMessage(queueURL, message)
	if err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Failed to send SQS message: %v", err),
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
		}
	}

	response := map[string]interface{}{
		"SendMessageResponse": map[string]interface{}{
			"SendMessageResult": map[string]string{
				"MD5OfMessageBody":       message.MD5OfBody,
				"MessageId":              message.ID,
				"MD5OfMessageAttributes": message.MD5OfMessageAttributes,
			},
			"ResponseMetadata": map[string]string{
				"RequestId": fmt.Sprintf("%x", time.Now().UnixNano()),
			},
		},
	}

	responseJSON, _ := json.Marshal(response)
	return &IntegrationResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            responseJSON,
		IsBase64Encoded: false,
	}, nil
}

func extractSQSMessageAttributes(headers, queryParams map[string]string, body []byte) map[string]*sqs.MessageAttributeValue {
	attrs := make(map[string]*sqs.MessageAttributeValue)
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
			attrs[name] = &sqs.MessageAttributeValue{
				DataType:    val.DataType,
				StringValue: &val.StringValue,
			}
		}
	}

	return attrs
}

func parseFormData(data string) map[string]string {
	result := make(map[string]string)
	for _, pair := range strings.Split(data, "&") {
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			val := parts[1]
			result[key] = val
		}
	}
	return result
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

	messages, err := e.sqsStore.ReceiveMessage(queueURL, maxMessages, &visibilityTimeout, waitTime)
	if err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Failed to receive SQS messages: %v", err),
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
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

	responseJSON, _ := json.Marshal(response)
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

// SNSPublishRequest represents a request for SNS integration in API Gateway.
type SNSPublishRequest struct {
	Message   string `json:"Message"`
	TopicArn  string `json:"TopicArn,omitempty"`
	TargetArn string `json:"TargetArn,omitempty"`
	Subject   string `json:"Subject,omitempty"`
}

func (e *AWSExecutor) executeSNS(ctx context.Context, req *IntegrationRequest) (*IntegrationResponse, error) {
	if e.snsStore == nil {
		return nil, &IntegrationError{
			Message:  "SNS store not configured",
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
		}
	}

	topicName, topicRegion, err := extractTopicFromURI(req.URI)
	if err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Invalid SNS URI: %v", err),
			Type:     "BadRequestException",
			HTTPCode: http.StatusBadRequest,
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
			HTTPCode: http.StatusBadRequest,
		}
	}
}

func (e *AWSExecutor) executeSNSPublish(ctx context.Context, topicArn string, req *IntegrationRequest) (*IntegrationResponse, error) {
	topic, err := e.snsStore.GetTopic(topicArn)
	if err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("SNS topic not found: %s", topicArn),
			Type:     "NotFoundException",
			HTTPCode: http.StatusNotFound,
		}
	}

	message := string(req.Body)
	if message == "" {
		message = req.Headers["Message"]
	}

	messageID := fmt.Sprintf("%x", time.Now().UnixNano())

	now := time.Now().UTC()
	notification := &sns.Message{
		MessageId:          messageID,
		TopicArn:           topicArn,
		Subject:            req.Headers["Subject"],
		Message:            message,
		PublishedTimestamp: now,
	}

	if err := e.snsStore.Put(topicArn+":messages:"+messageID, notification); err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Failed to store SNS message: %v", err),
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
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
		if err := e.bus.Publish(context.Background(), snsEvt); err != nil {
			logs.Warn("Failed to publish SNS delivery event to event bus; message is stored but subscribers may not be notified",
				logs.String("topicArn", topicArn),
				logs.String("messageId", messageID),
				logs.Err(err))
		}
	} else {
		if e.deliveryWg != nil {
			e.deliveryWg.Add(1)
		}
		go func() {
			if e.deliveryWg != nil {
				defer e.deliveryWg.Done()
			}
			defer func() {
				if r := recover(); r != nil {
					logs.Error("Panic in SNS delivery goroutine",
						logs.String("topicArn", topicArn),
						logs.Any("panic", r))
				}
			}()
			snsDeliverySem <- struct{}{}
			defer func() { <-snsDeliverySem }()
			e.deliverToSNSSubscribers(topicArn, notification)
		}()
	}

	_ = topic

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

	responseJSON, _ := json.Marshal(response)
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

func (e *AWSExecutor) deliverToSNSSubscribers(topicArn string, notification *sns.Message) {
	result, err := e.snsStore.ListSubscriptionsByTopic(topicArn, storecommon.ListOptions{})
	if err != nil {
		logs.Warn("failed to list SNS subscribers",
			logs.String("topicArn", topicArn),
			logs.String("messageId", notification.MessageId),
			logs.Err(err))
		return
	}

	for _, sub := range result.Items {
		if sub.PendingConfirmation {
			continue
		}

		switch {
		case strings.HasPrefix(sub.Endpoint, "arn:aws:sqs"):
			e.deliverToSQSSubscriber(sub.Endpoint, notification)
		case strings.HasPrefix(sub.Endpoint, "arn:aws:lambda"):
			e.deliverToLambdaSubscriber(sub.Endpoint, notification)
		}
	}
}

func (e *AWSExecutor) deliverToSQSSubscriber(queueArn string, notification *sns.Message) {
	if e.sqsStore == nil {
		return
	}

	parts := strings.Split(queueArn, ":")
	if len(parts) < 6 {
		return
	}
	queueName := parts[5]
	queueRegion := parts[3]
	if queueRegion == "" {
		queueRegion = e.region
	}
	queueURL := fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", queueRegion, e.accountID, queueName)

	body, err := json.Marshal(map[string]interface{}{
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

	msg := &sqs.Message{
		Body:              string(body),
		MessageAttributes: convertSNSAttrsToSQS(notification.MessageAttributes),
	}
	if _, err := e.sqsStore.SendMessage(queueURL, msg); err != nil {
		logs.Error("Failed to send SQS message to subscriber", logs.String("queue", queueURL), logs.Err(err))
	}
}

func convertSNSAttrsToSQS(snsAttrs map[string]*sns.MessageAttribute) map[string]*sqs.MessageAttributeValue {
	if snsAttrs == nil {
		return nil
	}
	sqsAttrs := make(map[string]*sqs.MessageAttributeValue)
	for k, v := range snsAttrs {
		attr := &sqs.MessageAttributeValue{
			DataType: v.Type,
		}
		if v.StringValue != "" {
			attr.StringValue = &v.StringValue
		}
		if v.BinaryValue != nil {
			attr.BinaryValue = v.BinaryValue
		}
		sqsAttrs[k] = attr
	}
	return sqsAttrs
}

func (e *AWSExecutor) deliverToLambdaSubscriber(functionArn string, notification *sns.Message) {
	if e.lambdaInvoker == nil {
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

	eventJSON, _ := json.Marshal(event)
	ctx := context.Background()
	if _, _, err := e.lambdaInvoker.InvokeForGateway(ctx, functionArn, eventJSON); err != nil {
		logs.Error("apigateway: SNS-to-Lambda invocation failed", logs.String("function", functionArn), logs.Err(err))
	}
}
