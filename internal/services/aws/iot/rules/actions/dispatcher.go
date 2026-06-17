// Package actions provides rule action implementations that dispatch IoT
// messages to downstream AWS services via the EventBus invoker pattern.
package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"vorpalstacks/internal/eventbus"
)

// ActionConfig holds the configuration for a single rule action extracted
// from the TopicRule action payload.
type ActionConfig struct {
	// Type identifies the action: "lambda", "sqs", "sns", "dynamodb", "s3",
	// "kinesis", "cloudwatch", "cloudwatchLogs", "republish", "http",
	// "stepFunctions", "eventbridge", "firehose", "timestream", "iotEvents".
	Type string

	// TargetARN is the resource ARN for the action target (Lambda, SQS, SNS,
	// DynamoDB, S3, Kinesis, EventBridge, Step Functions, Firehose).
	TargetARN string

	// QueueURL is the SQS queue URL (for SQS actions).
	QueueURL string

	// TopicARN is the SNS topic ARN (for SNS actions).
	TopicARN string

	// FunctionName is the Lambda function name or ARN.
	FunctionName string

	// TableName is the DynamoDB table name.
	TableName string

	// BucketName and ObjectKey are for S3 PutObject.
	BucketName string
	ObjectKey  string

	// StreamName is the Kinesis stream name.
	StreamName string

	// RepublishTopic is the MQTT topic for republish actions.
	RepublishTopic string

	// RoleARN is the IAM role ARN assumed by the action.
	RoleARN string

	// Extra holds any additional action-specific configuration.
	Extra map[string]interface{}
}

// Dispatcher routes evaluated rule payloads to the appropriate action handler
// using EventBus invokers for cross-service communication.
type Dispatcher struct {
	bus           eventbus.Bus
	logger        *slog.Logger
	RepublishFn   func(ctx context.Context, topic string, payload map[string]interface{}) error
	BatchPutMsgFn func(ctx context.Context, messages []map[string]interface{}) error
	HTTPPostFn    func(ctx context.Context, url string, payload []byte) error
}

// NewDispatcher creates a new action dispatcher backed by the given EventBus.
func NewDispatcher(bus eventbus.Bus, logger *slog.Logger) *Dispatcher {
	return &Dispatcher{bus: bus, logger: logger}
}

func (d *Dispatcher) SetRepublishFn(fn func(context.Context, string, map[string]interface{}) error) {
	d.RepublishFn = fn
}

func (d *Dispatcher) SetBatchPutMessageFn(fn func(context.Context, []map[string]interface{}) error) {
	d.BatchPutMsgFn = fn
}

func (d *Dispatcher) SetHTTPPostFn(fn func(context.Context, string, []byte) error) {
	d.HTTPPostFn = fn
}

// Dispatch sends a rule payload to the configured target action.
func (d *Dispatcher) Dispatch(ctx context.Context, config *ActionConfig, topic string, payload map[string]interface{}) error {
	if d.bus == nil {
		return fmt.Errorf("event bus not configured for IoT rules dispatcher")
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		if d.logger != nil {
			d.logger.Error("failed to marshal IoT rule payload", "error", err)
		}
		return fmt.Errorf("failed to marshal rule payload: %w", err)
	}

	ap := &ActionPayload{
		Topic:      topic,
		Raw:        payload,
		JSONBytes:  payloadJSON,
		JSONString: string(payloadJSON),
	}

	handler, ok := actionRegistry[config.Type]
	if !ok {
		if d.logger != nil {
			d.logger.Warn("unsupported IoT rule action type", "type", config.Type)
		}
		return fmt.Errorf("unsupported action type: %s", config.Type)
	}
	return handler(d, ctx, config, ap)
}

func (d *Dispatcher) dispatchLambda(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	invoker := d.bus.LambdaInvoker()
	if invoker == nil {
		return fmt.Errorf("lambda invoker not available")
	}

	_, _, err := invoker.InvokeForGateway(ctx, config.FunctionName, []byte(p.JSONString))
	if err != nil {
		return fmt.Errorf("lambda invocation failed: %w", err)
	}
	return nil
}

func (d *Dispatcher) dispatchSQS(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	invoker := d.bus.SQSInvoker()
	if invoker == nil {
		return fmt.Errorf("sqs invoker not available")
	}

	var queueURL string
	var err error

	if config.QueueURL != "" {
		queueURL = config.QueueURL
	} else {
		// Extract queue name from ARN and resolve URL.
		parts := strings.Split(config.TargetARN, ":")
		queueName := parts[len(parts)-1]
		queueURL, err = invoker.GetQueueByName(ctx, queueName)
		if err != nil {
			return fmt.Errorf("failed to resolve SQS queue URL: %w", err)
		}
	}

	_, _, err = invoker.SendMessage(ctx, queueURL, p.JSONString, eventbus.SQSSendOptions{})
	if err != nil {
		return fmt.Errorf("sqs send failed: %w", err)
	}
	return nil
}

func (d *Dispatcher) dispatchSNS(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	invoker := d.bus.SNSInvoker()
	if invoker == nil {
		return fmt.Errorf("sns invoker not available")
	}

	topicARN := config.TopicARN
	if topicARN == "" {
		topicARN = config.TargetARN
	}

	_, err := invoker.PublishToTopic(ctx, topicARN, p.JSONString, "", nil)
	if err != nil {
		return fmt.Errorf("sns publish failed: %w", err)
	}
	return nil
}

func (d *Dispatcher) dispatchDynamoDB(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	invoker := d.bus.DynamoDBInvoker()
	if invoker == nil {
		return fmt.Errorf("dynamodb invoker not available")
	}

	tableName := config.TableName
	if tableName == "" {
		// Extract table name from ARN.
		parts := strings.Split(config.TargetARN, "/")
		if len(parts) > 1 {
			tableName = parts[len(parts)-1]
		}
	}

	if tableName == "" {
		return fmt.Errorf("dynamodb table name not specified")
	}

	// Use the entire payload as item attributes. Derive a partition key
	// from the "id" field or generate one.
	key := deriveDynamoDBKey(p.Raw)

	_, err := invoker.PutItem(ctx, "", tableName, key, p.Raw)
	if err != nil {
		return fmt.Errorf("dynamodb put failed: %w", err)
	}
	return nil
}

func (d *Dispatcher) dispatchS3(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	invoker := d.bus.S3Invoker()
	if invoker == nil {
		return fmt.Errorf("s3 invoker not available")
	}

	key := config.ObjectKey
	if key == "" {
		key = fmt.Sprintf("iot-rules/%s/%d.json", config.Type, time.Now().UnixNano())
	}

	err := invoker.PutObject(ctx, "", config.BucketName, key, p.JSONBytes, "application/json")
	if err != nil {
		return fmt.Errorf("s3 put failed: %w", err)
	}
	return nil
}

func (d *Dispatcher) dispatchKinesis(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	invoker := d.bus.KinesisInvoker()
	if invoker == nil {
		return fmt.Errorf("kinesis invoker not available")
	}

	streamName := config.StreamName
	if streamName == "" {
		parts := strings.Split(config.TargetARN, "/")
		if len(parts) > 1 {
			streamName = parts[len(parts)-1]
		}
	}

	if streamName == "" {
		return fmt.Errorf("kinesis stream name not specified")
	}

	partitionKey := fmt.Sprintf("%d", time.Now().UnixNano())
	_, err := invoker.PutRecord(ctx, streamName, partitionKey, p.JSONBytes)
	if err != nil {
		return fmt.Errorf("kinesis put record failed: %w", err)
	}
	return nil
}

func (d *Dispatcher) dispatchCloudWatch(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	invoker := d.bus.CloudWatchMetricInvoker()
	if invoker == nil {
		return fmt.Errorf("cloudwatch metric invoker not available")
	}

	namespace := "AWS/IoT"
	if ns, ok := config.Extra["namespace"].(string); ok {
		namespace = ns
	}

	metricName := "RuleInvocations"
	if mn, ok := config.Extra["metricName"].(string); ok {
		metricName = mn
	}

	var value float64 = 1
	if v, ok := p.Raw["value"].(float64); ok {
		value = v
	}

	err := invoker.PutMetricData("", namespace, metricName, value, time.Now())
	if err != nil {
		return fmt.Errorf("cloudwatch put metric failed: %w", err)
	}
	return nil
}

func (d *Dispatcher) dispatchCloudWatchLogs(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	invoker := d.bus.LogsInvoker()
	if invoker == nil {
		return fmt.Errorf("cloudwatch logs invoker not available")
	}

	logGroup := "aws-iot-rules"
	if lg, ok := config.Extra["logGroup"].(string); ok {
		logGroup = lg
	}

	logStream := "default"
	if ls, ok := config.Extra["logStream"].(string); ok {
		logStream = ls
	}

	if err := invoker.EnsureLogGroup(ctx, "", logGroup, ""); err != nil {
		return fmt.Errorf("ensure log group failed: %w", err)
	}
	if err := invoker.EnsureLogStream(ctx, "", logGroup, logStream); err != nil {
		return fmt.Errorf("ensure log stream failed: %w", err)
	}

	entries := []eventbus.LogsLogEntry{
		{Timestamp: time.Now().UnixMilli(), Message: p.JSONString},
	}
	if err := invoker.PutLogEvents(ctx, "", logGroup, logStream, entries); err != nil {
		return fmt.Errorf("cloudwatch logs put failed: %w", err)
	}
	return nil
}

func (d *Dispatcher) dispatchEventBridge(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	invoker := d.bus.EventsInvoker()
	if invoker == nil {
		return fmt.Errorf("eventbridge invoker not available")
	}

	err := invoker.PutEvent(ctx, config.TargetARN, p.Raw)
	if err != nil {
		return fmt.Errorf("eventbridge put event failed: %w", err)
	}
	return nil
}

// deriveDynamoDBKey builds a DynamoDB key map from the payload. If the
// payload has an "id" field it is used as the partition key.
func deriveDynamoDBKey(payload map[string]interface{}) map[string]interface{} {
	if id, ok := payload["id"]; ok {
		return map[string]interface{}{"id": id}
	}
	return map[string]interface{}{
		"id": fmt.Sprintf("%d", time.Now().UnixNano()),
	}
}

// NewActionConfigFromMap builds an ActionConfig from an action type string and
// the raw configuration map extracted from a TopicRule action payload. The
// map keys follow the AWS IoT action field names (e.g. "functionArn",
// "queueUrl", "targetArn").
func NewActionConfigFromMap(actionType string, m map[string]interface{}) *ActionConfig {
	ac := &ActionConfig{Type: actionType, Extra: m}
	switch actionType {
	case "lambda":
		ac.FunctionName = strFromMap(m, "functionArn", "functionName")
	case "sqs":
		ac.QueueURL = strFromMap(m, "queueUrl")
		ac.TargetARN = strFromMap(m, "queueArn")
	case "sns":
		ac.TopicARN = strFromMap(m, "targetArn", "topicArn")
		ac.TargetARN = ac.TopicARN
	case "dynamodb":
		ac.TableName = strFromMap(m, "tableName")
		ac.TargetARN = strFromMap(m, "tableArn")
	case "s3":
		ac.BucketName = strFromMap(m, "bucketName")
		ac.ObjectKey = strFromMap(m, "key")
	case "kinesis":
		ac.StreamName = strFromMap(m, "streamName")
		ac.TargetARN = strFromMap(m, "streamArn")
	case "republish":
		ac.RepublishTopic = strFromMap(m, "topic")
	case "stepFunctions":
		ac.TargetARN = strFromMap(m, "stateMachineArn")
	case "eventbridge":
		ac.TargetARN = strFromMap(m, "eventBusArn")
	}
	if ac.RoleARN == "" {
		ac.RoleARN = strFromMap(m, "roleArn")
	}
	return ac
}

func strFromMap(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok && v != "" {
			return v
		}
	}
	return ""
}

func (d *Dispatcher) dispatchRepublish(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	if d.RepublishFn == nil {
		return fmt.Errorf("republish: not configured")
	}
	topic := config.RepublishTopic
	if topic == "" {
		topic = config.TopicARN
	}
	if topic == "" {
		return fmt.Errorf("republish: no topic specified")
	}
	return d.RepublishFn(ctx, topic, p.Raw)
}

func (d *Dispatcher) dispatchStepFunctions(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	events := d.bus.EventsInvoker()
	if events == nil {
		return fmt.Errorf("stepFunctions: events invoker not available")
	}
	event := map[string]interface{}{
		"input":           p.JSONString,
		"stateMachineArn": config.TargetARN,
	}
	key := fmt.Sprintf("iot-rule/%d", time.Now().UnixNano())
	return events.PutEvent(ctx, key, event)
}

func (d *Dispatcher) dispatchIoTEvents(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	if d.BatchPutMsgFn == nil {
		return fmt.Errorf("iotEvents: not configured")
	}
	msg := map[string]interface{}{
		"messageId": fmt.Sprintf("rule-%d", time.Now().UnixNano()),
		"inputName": config.Extra["inputName"],
		"payload":   p.Raw,
	}
	return d.BatchPutMsgFn(ctx, []map[string]interface{}{msg})
}

func (d *Dispatcher) dispatchHTTP(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	if d.HTTPPostFn == nil {
		return fmt.Errorf("http: not configured")
	}
	url, _ := config.Extra["url"].(string)
	if url == "" {
		return fmt.Errorf("http: no url specified")
	}
	return d.HTTPPostFn(ctx, url, p.JSONBytes)
}
