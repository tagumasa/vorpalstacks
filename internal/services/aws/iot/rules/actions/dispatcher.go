// Package actions provides rule action implementations that dispatch IoT
// messages to downstream AWS services via the EventBus invoker pattern.
package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/services/aws/iot/iotutil"
)

// ActionConfig holds the configuration for a single rule action extracted
// from the TopicRule action payload.
type ActionConfig struct {
	// Type identifies the action: "lambda", "sqs", "sns", "dynamodb", "s3",
	// "kinesis", "cloudwatch", "cloudwatchLogs", "republish", "http",
	// "stepFunctions", "eventbridge", "firehose", "timestream".
	Type string

	// TargetARN is the resource ARN for the action target (Lambda, SQS, SNS,
	// DynamoDB, S3, Kinesis, EventBridge, Step Functions, Firehose).
	TargetARN string

	// QueueURL is the SQS queue URL (for SQS actions).
	QueueURL string

	// Region is the AWS region for the action target, used for
	// region-aware SQS delivery when QueueURL is provided without ARN.
	Region string

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
	bus         eventbus.Bus
	logger      *slog.Logger
	RepublishFn func(ctx context.Context, topic string, payload map[string]interface{}) error
	HTTPPostFn  func(ctx context.Context, url string, payload []byte) error
}

// NewDispatcher creates a new action dispatcher backed by the given EventBus.
func NewDispatcher(bus eventbus.Bus, logger *slog.Logger) *Dispatcher {
	return &Dispatcher{bus: bus, logger: logger}
}

func (d *Dispatcher) SetRepublishFn(fn func(context.Context, string, map[string]interface{}) error) {
	d.RepublishFn = fn
}

func (d *Dispatcher) SetHTTPPostFn(fn func(context.Context, string, []byte) error) {
	d.HTTPPostFn = fn
}

// Dispatch sends a rule payload to the configured target action.
func (d *Dispatcher) Dispatch(ctx context.Context, config *ActionConfig, topic string, payload map[string]interface{}) error {
	if d.bus == nil {
		return fmt.Errorf("event bus not configured for IoT rules dispatcher")
	}

	// Resolve AWS IoT substitution templates (${topic()}, ${payload.x},
	// ${timestamp()}, etc.) in action config string fields before dispatch.
	d.resolveTemplates(config, topic, payload)

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

// resolveTemplates applies AWS IoT substitution templates to action config
// string fields that support them (DynamoDB hashKeyValue/rangeKeyValue,
// S3 key, Kinesis partitionKey). Non-string values are left unchanged.
func (d *Dispatcher) resolveTemplates(config *ActionConfig, topic string, payload map[string]interface{}) {
	// Resolve string fields on ActionConfig that support substitution
	// templates per AWS IoT documentation.
	//
	// Reference: https://docs.aws.amazon.com/iot/latest/developerguide/iot-substitution-templates.html

	// Republish topic — "Supports substitution templates: Yes"
	if config.RepublishTopic != "" {
		config.RepublishTopic = iotutil.ResolveTemplate(config.RepublishTopic, topic, "", payload)
	}

	// S3 key — "Supports substitution templates: Yes"
	if config.ObjectKey != "" {
		config.ObjectKey = iotutil.ResolveTemplate(config.ObjectKey, topic, "", payload)
	}

	// S3 bucket — "Supports substitution templates: API and AWS CLI only"
	if config.BucketName != "" {
		config.BucketName = iotutil.ResolveTemplate(config.BucketName, topic, "", payload)
	}

	// SQS queueUrl — "Supports substitution templates: API and AWS CLI only"
	if config.QueueURL != "" {
		config.QueueURL = iotutil.ResolveTemplate(config.QueueURL, topic, "", payload)
	}

	// SNS targetArn — "Supports substitution templates: API and AWS CLI only"
	if config.TopicARN != "" {
		config.TopicARN = iotutil.ResolveTemplate(config.TopicARN, topic, "", payload)
	}

	// Kinesis streamName — "Supports substitution templates: API and AWS CLI only"
	if config.StreamName != "" {
		config.StreamName = iotutil.ResolveTemplate(config.StreamName, topic, "", payload)
	}

	// Extra is a reference to the persistent rule action map. Copy it so
	// that resolved values do not overwrite the templates for subsequent
	// dispatches.
	if len(config.Extra) > 0 {
		copied := make(map[string]interface{}, len(config.Extra))
		for k, v := range config.Extra {
			copied[k] = v
		}
		config.Extra = copied
	}

	// Resolve substitution templates in Extra map string fields.
	//
	// Fields marked "Supports substitution templates: Yes" (fully supported):
	//   hashKeyValue, rangeKeyValue, payloadField, operation (DynamoDB)
	//   partitionKey (Kinesis)
	//   metricName, metricNamespace, metricUnit, metricValue, metricTimestamp (CloudWatch Metric)
	//
	// Fields marked "Supports substitution templates: API and AWS CLI only":
	//   tableName, hashKeyField, hashKeyType, rangeKeyField, rangeKeyType (DynamoDB)
	//   stream (Kinesis — already handled via StreamName above)
	//
	// metricValue is also resolved in dispatchCloudWatch but we resolve it
	// here centrally so the per-action handler does not need to duplicate
	// the logic.
	templateFields := []string{
		"hashKeyValue", "rangeKeyValue", "payloadField", "operation",
		"partitionKey",
		"metricName", "metricNamespace", "metricUnit", "metricValue", "metricTimestamp",
		"tableName", "hashKeyField", "hashKeyType", "rangeKeyField", "rangeKeyType",
		"inputName", "messageId",
	}
	for _, field := range templateFields {
		if v, ok := config.Extra[field].(string); ok && v != "" {
			config.Extra[field] = iotutil.ResolveTemplate(v, topic, "", payload)
		}
	}
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
	var sqsRegion string
	var err error

	if config.QueueURL != "" {
		queueURL = config.QueueURL
		sqsRegion = config.Region
	} else if config.TargetARN != "" {
		// Extract queue name and region from ARN and resolve URL.
		parts := strings.Split(config.TargetARN, ":")
		if len(parts) < 6 {
			return fmt.Errorf("sqs: malformed queue ARN: %s", config.TargetARN)
		}
		queueName := parts[len(parts)-1]
		sqsRegion = parts[3]
		queueURL, err = invoker.GetQueueByName(ctx, sqsRegion, queueName)
		if err != nil {
			return fmt.Errorf("failed to resolve SQS queue URL: %w", err)
		}
	} else {
		return fmt.Errorf("sqs: no queueUrl or queueArn specified")
	}

	_, _, err = invoker.SendMessage(ctx, sqsRegion, queueURL, p.JSONString, eventbus.SQSSendOptions{})
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
		parts := strings.Split(config.TargetARN, "/")
		if len(parts) > 1 {
			tableName = parts[len(parts)-1]
		}
	}

	if tableName == "" {
		return fmt.Errorf("dynamodb table name not specified")
	}

	// AWS IoT DynamoDBAction names the partition/sort keys and payload field
	// explicitly; honour them rather than forcing a hardcoded "id" key.
	key, attributes := buildDynamoDBItem(p.Raw, config.Extra)

	// AWS IoT DynamoDB action supports an optional "operation" field:
	// INSERT (default), UPDATE, or DELETE. INSERT uses PutItem, UPDATE uses
	// UpdateItem, DELETE uses DeleteItem.
	operation := iotutil.StrFromMap(config.Extra, "operation")
	switch strings.ToUpper(operation) {
	case "", "INSERT":
		_, err := invoker.PutItem(ctx, "", tableName, key, attributes)
		if err != nil {
			return fmt.Errorf("dynamodb put failed: %w", err)
		}
	case "UPDATE":
		if err := invoker.UpdateItem(ctx, "", tableName, key, attributes); err != nil {
			return fmt.Errorf("dynamodb update failed: %w", err)
		}
	case "DELETE":
		if err := invoker.DeleteItem(ctx, "", tableName, key); err != nil {
			return fmt.Errorf("dynamodb delete failed: %w", err)
		}
	default:
		return fmt.Errorf("dynamodb: unsupported operation %q (expected INSERT, UPDATE, or DELETE)", operation)
	}
	return nil
}

// dispatchDynamoDBv2 writes the entire message payload as a DynamoDB item.
// Unlike dynamoDB which extracts specific fields, dynamoDBv2 puts every
// payload attribute into its own column (Smithy DynamoDBv2Action.putItem).
func (d *Dispatcher) dispatchDynamoDBv2(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	invoker := d.bus.DynamoDBInvoker()
	if invoker == nil {
		return fmt.Errorf("dynamodb invoker not available")
	}
	// tableName comes from putItem.tableName per Smithy DynamoDBv2Action.
	tableName := ""
	if putItem, ok := config.Extra["putItem"].(map[string]interface{}); ok {
		tableName, _ = putItem["tableName"].(string)
	}
	if tableName == "" {
		tableName = iotutil.StrFromMap(config.Extra, "tableName")
	}
	if tableName == "" {
		return fmt.Errorf("dynamoDBv2: table name not specified")
	}
	// The entire payload becomes the item attributes. An empty key map
	// signals the invoker to use the payload directly.
	_, err := invoker.PutItem(ctx, "", tableName, p.Raw, p.Raw)
	if err != nil {
		return fmt.Errorf("dynamoDBv2 put failed: %w", err)
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

	partitionKey := iotutil.StrFromMap(config.Extra, "partitionKey")
	if partitionKey == "" {
		partitionKey = fmt.Sprintf("%d", time.Now().UnixNano())
	}
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

	// All substitution-template fields (metricNamespace, metricName,
	// metricUnit, metricValue, metricTimestamp) are already resolved by
	// resolveTemplates before Dispatch calls this handler.
	namespace := iotutil.StrFromMap(config.Extra, "metricNamespace", "namespace")
	if namespace == "" {
		namespace = "AWS/IoT"
	}

	metricName := iotutil.StrFromMap(config.Extra, "metricName")
	if metricName == "" {
		metricName = "RuleInvocations"
	}

	var value float64 = 1
	if mv := iotutil.StrFromMap(config.Extra, "metricValue"); mv != "" {
		if parsed, err := strconv.ParseFloat(mv, 64); err == nil {
			value = parsed
		}
	} else if v, ok := p.Raw["value"].(float64); ok {
		value = v
	}

	ts := time.Now()
	if mts := iotutil.StrFromMap(config.Extra, "metricTimestamp"); mts != "" {
		if epoch, err := strconv.ParseInt(mts, 10, 64); err == nil {
			ts = time.Unix(epoch, 0)
		}
	}

	err := invoker.PutMetricData("", namespace, metricName, value, ts)
	if err != nil {
		return fmt.Errorf("cloudwatch put metric failed: %w", err)
	}
	return nil
}

// dispatchCloudWatchAlarm sets a CloudWatch alarm state (Smithy
// CloudwatchAlarmAction: alarmName, stateReason, stateValue).
func (d *Dispatcher) dispatchCloudWatchAlarm(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	invoker := d.bus.CloudWatchAlarmInvoker()
	if invoker == nil {
		return fmt.Errorf("cloudwatch alarm invoker not available")
	}
	alarmName := iotutil.StrFromMap(config.Extra, "alarmName")
	if alarmName == "" {
		return fmt.Errorf("cloudwatchAlarm: alarmName not specified")
	}
	stateReason := iotutil.ResolveTemplate(iotutil.StrFromMap(config.Extra, "stateReason"), p.Topic, "", p.Raw)
	stateValue := iotutil.ResolveTemplate(iotutil.StrFromMap(config.Extra, "stateValue"), p.Topic, "", p.Raw)
	if stateValue == "" {
		stateValue = "ALARM"
	}
	err := invoker.SetAlarmState("", alarmName, stateValue, stateReason)
	if err != nil {
		return fmt.Errorf("cloudwatch set alarm state failed: %w", err)
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

// buildDynamoDBItem constructs the DynamoDB key map and attribute map from
// the rule payload using the AWS IoT DynamoDBAction field names
// (hashKeyField/hashKeyValue/rangeKeyField/rangeKeyValue/payloadField).
func buildDynamoDBItem(payload map[string]interface{}, extra map[string]interface{}) (map[string]interface{}, map[string]interface{}) {
	hashField := iotutil.StrFromMap(extra, "hashKeyField")
	if hashField == "" {
		hashField = "id"
	}
	hashValue := iotutil.StrFromMap(extra, "hashKeyValue")
	if hashValue == "" {
		if id, ok := payload[hashField]; ok {
			hashValue = fmt.Sprintf("%v", id)
		} else {
			hashValue = fmt.Sprintf("%d", time.Now().UnixNano())
		}
	}

	key := map[string]interface{}{hashField: hashValue}

	rangeField := iotutil.StrFromMap(extra, "rangeKeyField")
	if rangeField != "" {
		rangeValue := iotutil.StrFromMap(extra, "rangeKeyValue")
		if rangeValue == "" {
			rangeValue = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		key[rangeField] = rangeValue
	}

	attributes := map[string]interface{}{}
	for k, v := range payload {
		attributes[k] = v
	}
	for k, v := range key {
		attributes[k] = v
	}
	if payloadField := iotutil.StrFromMap(extra, "payloadField"); payloadField != "" {
		attributes[payloadField] = string(mustJSON(payload))
	}

	return key, attributes
}

func mustJSON(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return b
}

// NewActionConfigFromMap builds an ActionConfig from an action type string and
// the raw configuration map extracted from a TopicRule action payload. The
// map keys follow the AWS IoT action field names (e.g. "functionArn",
// "queueUrl", "targetArn").
func NewActionConfigFromMap(actionType string, m map[string]interface{}) *ActionConfig {
	ac := &ActionConfig{Type: actionType, Extra: m}
	switch actionType {
	case "lambda":
		ac.FunctionName = iotutil.StrFromMap(m, "functionArn", "functionName")
	case "sqs":
		ac.QueueURL = iotutil.StrFromMap(m, "queueUrl")
		ac.TargetARN = iotutil.StrFromMap(m, "queueArn")
	case "sns":
		ac.TopicARN = iotutil.StrFromMap(m, "targetArn", "topicArn")
		ac.TargetARN = ac.TopicARN
	case "dynamoDB":
		ac.TableName = iotutil.StrFromMap(m, "tableName")
		ac.TargetARN = iotutil.StrFromMap(m, "tableArn")
	case "s3":
		ac.BucketName = iotutil.StrFromMap(m, "bucketName")
		ac.ObjectKey = iotutil.StrFromMap(m, "key")
	case "kinesis":
		ac.StreamName = iotutil.StrFromMap(m, "streamName")
		ac.TargetARN = iotutil.StrFromMap(m, "streamArn")
	case "republish":
		ac.RepublishTopic = iotutil.StrFromMap(m, "topic")
	case "stepFunctions":
		ac.TargetARN = iotutil.StrFromMap(m, "stateMachineName", "stateMachineArn")
	}
	if ac.RoleARN == "" {
		ac.RoleARN = iotutil.StrFromMap(m, "roleArn")
	}
	return ac
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
	// Copy the payload so the iteration counter does not pollute the original
	// map shared with other actions in the same rule dispatch.
	republishPayload := make(map[string]interface{}, len(p.Raw)+1)
	for k, v := range p.Raw {
		republishPayload[k] = v
	}
	// Increment iteration counter for republish chain tracking.
	// The Executor's OnMessage extracts this field and enforces the 25-iteration limit.
	if iterRaw, ok := republishPayload["_iotRuleIteration"]; ok {
		republishPayload["_iotRuleIteration"] = iotutil.ToInt(iterRaw) + 1
	} else {
		republishPayload["_iotRuleIteration"] = 1
	}
	return d.RepublishFn(ctx, topic, republishPayload)
}

func (d *Dispatcher) dispatchStepFunctions(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	// Publish a start-execution event; do NOT call EventBridge PutEvent (that
	// only inserts a bus entry rather than starting a state machine run).
	stateMachineArn := config.TargetARN
	stateMachineName := ""
	if stateMachineArn == "" {
		// Smithy StepFunctionsAction has stateMachineName (not Arn).
		// Safe type assertion to avoid panic when key is absent.
		if v, ok := config.Extra["stateMachineArn"].(string); ok && v != "" {
			stateMachineArn = v
		} else if v, ok := config.Extra["stateMachineName"].(string); ok && v != "" {
			stateMachineName = v
		}
	}
	if stateMachineArn == "" && stateMachineName == "" {
		return fmt.Errorf("stepFunctions: state machine identifier not specified")
	}
	evt := &eventbus.StepFunctionsStartExecutionEvent{
		StateMachineArn:  stateMachineArn,
		StateMachineName: stateMachineName,
		Input:            p.JSONString,
	}
	return d.bus.Publish(ctx, evt)
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

// dispatchTimestream writes IoT message data to a Timestream table (Smithy
// TimestreamAction: databaseName, tableName, dimensions, timestamp).
func (d *Dispatcher) dispatchTimestream(ctx context.Context, config *ActionConfig, p *ActionPayload) error {
	invoker := d.bus.TimestreamInvoker()
	if invoker == nil {
		return fmt.Errorf("timestream invoker not available")
	}
	databaseName := iotutil.StrFromMap(config.Extra, "databaseName")
	tableName := iotutil.StrFromMap(config.Extra, "tableName")
	if databaseName == "" || tableName == "" {
		return fmt.Errorf("timestream: databaseName and tableName are required")
	}
	// Build dimensions from the action config. AWS IoT supports substitution
	// templates in dimension values, so resolve them here.
	dimensions := map[string]string{}
	if dims, ok := config.Extra["dimensions"].([]interface{}); ok {
		for _, d := range dims {
			if dm, ok := d.(map[string]interface{}); ok {
				name, _ := dm["name"].(string)
				value, _ := dm["value"].(string)
				if name != "" {
					dimensions[name] = iotutil.ResolveTemplate(value, p.Topic, "", p.Raw)
				}
			}
		}
	}
	// Determine the timestamp. AWS IoT allows a user-defined timestamp via
	// the "timestamp" field, which contains value (epoch) and unit.
	ts := time.Now()
	if tsRaw, ok := config.Extra["timestamp"].(map[string]interface{}); ok {
		tsValue := iotutil.StrFromMap(tsRaw, "value")
		tsUnit := iotutil.StrFromMap(tsRaw, "unit")
		if tsValue != "" {
			resolved := iotutil.ResolveTemplate(tsValue, p.Topic, "", p.Raw)
			if epoch, err := strconv.ParseInt(resolved, 10, 64); err == nil {
				ts = epochToTime(epoch, tsUnit)
			}
		}
	}
	// AWS IoT Timestream action writes EVERY top-level payload field as a
	// separate measure record, sharing the same dimensions and timestamp.
	for k, v := range p.Raw {
		measureValue, measureType := measureForValue(v)
		if measureValue == "" {
			continue
		}
		if err := invoker.WriteRecords("", databaseName, tableName, dimensions, k, measureValue, measureType, ts); err != nil {
			return fmt.Errorf("timestream write failed for measure %q: %w", k, err)
		}
	}
	return nil
}

func epochToTime(epoch int64, unit string) time.Time {
	switch unit {
	case "SECONDS":
		return time.Unix(epoch, 0)
	case "MICROSECONDS":
		return time.Unix(0, epoch*1000)
	case "NANOSECONDS":
		return time.Unix(0, epoch)
	default: // MILLISECONDS
		return time.Unix(0, epoch*int64(time.Millisecond))
	}
}

func measureForValue(v interface{}) (value, typ string) {
	switch fv := v.(type) {
	case float64:
		return strconv.FormatFloat(fv, 'f', -1, 64), "DOUBLE"
	case int:
		return strconv.Itoa(fv), "BIGINT"
	case int64:
		return strconv.FormatInt(fv, 10), "BIGINT"
	case bool:
		if fv {
			return "true", "BOOLEAN"
		}
		return "false", "BOOLEAN"
	case string:
		return fv, "VARCHAR"
	default:
		return "", ""
	}
}
