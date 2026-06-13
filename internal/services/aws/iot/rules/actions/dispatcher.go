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
	bus    eventbus.Bus
	logger *slog.Logger
}

// NewDispatcher creates a new action dispatcher backed by the given EventBus.
func NewDispatcher(bus eventbus.Bus, logger *slog.Logger) *Dispatcher {
	return &Dispatcher{bus: bus, logger: logger}
}

// Dispatch sends a rule payload to the configured target action.
func (d *Dispatcher) Dispatch(ctx context.Context, config *ActionConfig, topic string, payload map[string]interface{}) error {
	if d.bus == nil {
		return fmt.Errorf("event bus not configured for IoT rules dispatcher")
	}

	payloadJSON, _ := json.Marshal(payload)
	payloadStr := string(payloadJSON)

	switch config.Type {
	case "lambda":
		return d.dispatchLambda(ctx, config, payloadStr)
	case "sqs":
		return d.dispatchSQS(ctx, config, payloadStr)
	case "sns":
		return d.dispatchSNS(ctx, config, topic, payloadStr)
	case "dynamodb":
		return d.dispatchDynamoDB(ctx, config, payload)
	case "s3":
		return d.dispatchS3(ctx, config, payloadJSON)
	case "kinesis":
		return d.dispatchKinesis(ctx, config, payloadJSON)
	case "cloudwatch":
		return d.dispatchCloudWatch(ctx, config, payload)
	case "cloudwatchLogs":
		return d.dispatchCloudWatchLogs(ctx, config, payloadStr)
	case "eventbridge":
		return d.dispatchEventBridge(ctx, config, payload)
	case "republish":
		return fmt.Errorf("republish action requires broker access, use inline dispatch")
	default:
		if d.logger != nil {
			d.logger.Warn("unsupported IoT rule action type", "type", config.Type)
		}
		return fmt.Errorf("unsupported action type: %s", config.Type)
	}
}

func (d *Dispatcher) dispatchLambda(ctx context.Context, config *ActionConfig, payload string) error {
	invoker := d.bus.LambdaInvoker()
	if invoker == nil {
		return fmt.Errorf("lambda invoker not available")
	}

	_, _, err := invoker.InvokeForGateway(ctx, config.FunctionName, []byte(payload))
	if err != nil {
		return fmt.Errorf("lambda invocation failed: %w", err)
	}
	return nil
}

func (d *Dispatcher) dispatchSQS(ctx context.Context, config *ActionConfig, payload string) error {
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

	_, _, err = invoker.SendMessage(ctx, queueURL, payload, eventbus.SQSSendOptions{})
	if err != nil {
		return fmt.Errorf("sqs send failed: %w", err)
	}
	return nil
}

func (d *Dispatcher) dispatchSNS(ctx context.Context, config *ActionConfig, topic string, payload string) error {
	invoker := d.bus.SNSInvoker()
	if invoker == nil {
		return fmt.Errorf("sns invoker not available")
	}

	topicARN := config.TopicARN
	if topicARN == "" {
		topicARN = config.TargetARN
	}

	_, err := invoker.PublishToTopic(ctx, topicARN, payload, "", nil)
	if err != nil {
		return fmt.Errorf("sns publish failed: %w", err)
	}
	return nil
}

func (d *Dispatcher) dispatchDynamoDB(ctx context.Context, config *ActionConfig, payload map[string]interface{}) error {
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
	key := deriveDynamoDBKey(payload)

	_, err := invoker.PutItem(ctx, "", tableName, key, payload)
	if err != nil {
		return fmt.Errorf("dynamodb put failed: %w", err)
	}
	return nil
}

func (d *Dispatcher) dispatchS3(ctx context.Context, config *ActionConfig, payload []byte) error {
	invoker := d.bus.S3Invoker()
	if invoker == nil {
		return fmt.Errorf("s3 invoker not available")
	}

	key := config.ObjectKey
	if key == "" {
		key = fmt.Sprintf("iot-rules/%s/%d.json", config.Type, time.Now().UnixNano())
	}

	err := invoker.PutObject(ctx, "", config.BucketName, key, payload, "application/json")
	if err != nil {
		return fmt.Errorf("s3 put failed: %w", err)
	}
	return nil
}

func (d *Dispatcher) dispatchKinesis(ctx context.Context, config *ActionConfig, payload []byte) error {
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
	_, err := invoker.PutRecord(ctx, streamName, partitionKey, payload)
	if err != nil {
		return fmt.Errorf("kinesis put record failed: %w", err)
	}
	return nil
}

func (d *Dispatcher) dispatchCloudWatch(ctx context.Context, config *ActionConfig, payload map[string]interface{}) error {
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
	if v, ok := payload["value"].(float64); ok {
		value = v
	}

	err := invoker.PutMetricData("", namespace, metricName, value, time.Now())
	if err != nil {
		return fmt.Errorf("cloudwatch put metric failed: %w", err)
	}
	return nil
}

func (d *Dispatcher) dispatchCloudWatchLogs(ctx context.Context, config *ActionConfig, payload string) error {
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
		{Timestamp: time.Now().UnixMilli(), Message: payload},
	}
	if err := invoker.PutLogEvents(ctx, "", logGroup, logStream, entries); err != nil {
		return fmt.Errorf("cloudwatch logs put failed: %w", err)
	}
	return nil
}

func (d *Dispatcher) dispatchEventBridge(ctx context.Context, config *ActionConfig, payload map[string]interface{}) error {
	invoker := d.bus.EventsInvoker()
	if invoker == nil {
		return fmt.Errorf("eventbridge invoker not available")
	}

	err := invoker.PutEvent(ctx, config.TargetARN, payload)
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
