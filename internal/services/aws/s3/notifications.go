package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	s3store "vorpalstacks/internal/store/aws/s3"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// handleS3Notification is the async bus handler for S3ObjectEvent.
// It reads the bucket's NotificationConfiguration, matches the event type
// and key filter rules, then dispatches the S3 event record to each
// configured destination (SNS topic, SQS queue, or Lambda function).
func (s *S3Service) handleS3Notification(ctx context.Context, event *eventbus.S3ObjectEvent) eventbus.HandlerResult {
	if s.s3Store == nil {
		return eventbus.HandlerResult{}
	}

	buckets := s.s3Store.Buckets(event.EventRegion())
	if buckets == nil {
		return eventbus.HandlerResult{}
	}

	config, err := buckets.GetNotificationConfiguration(event.Bucket)
	if err != nil || config == nil {
		return eventbus.HandlerResult{}
	}

	eventName := "s3:" + string(event.Op)
	attrs := s3NotificationAttributes(event)

	for _, tc := range config.TopicConfigurations {
		if !matchS3Event(tc.Events, eventName) {
			continue
		}
		if tc.Filter != nil && !matchS3KeyFilter(event.Key, tc.Filter) {
			continue
		}
		payload := buildS3EventRecord(event, tc.Id)
		s.dispatchToSNS(ctx, tc.TopicArn, payload, attrs)
	}

	for _, qc := range config.QueueConfigurations {
		if !matchS3Event(qc.Events, eventName) {
			continue
		}
		if qc.Filter != nil && !matchS3KeyFilter(event.Key, qc.Filter) {
			continue
		}
		payload := buildS3EventRecord(event, qc.Id)
		s.dispatchToSQS(ctx, qc.QueueArn, payload, attrs)
	}

	for _, lc := range config.LambdaConfigurations {
		if !matchS3Event(lc.Events, eventName) {
			continue
		}
		if lc.Filter != nil && !matchS3KeyFilter(event.Key, lc.Filter) {
			continue
		}
		payload := buildS3EventRecord(event, lc.Id)
		s.dispatchToLambda(ctx, lc.LambdaFunctionArn, payload)
	}

	if config.EventBridgeConfiguration != nil {
		s.dispatchToEventBridge(ctx, event, eventName)
	}

	return eventbus.HandlerResult{}
}

// matchS3Event checks whether the actual event name matches any of the
// configured event patterns. Wildcard patterns ending in ":*" match all
// events sharing the same prefix (e.g. "s3:ObjectCreated:*").
func matchS3Event(configuredEvents []string, eventName string) bool {
	for _, cfg := range configuredEvents {
		if cfg == eventName {
			return true
		}
		if strings.HasSuffix(cfg, ":*") {
			prefix := cfg[:len(cfg)-1]
			if strings.HasPrefix(eventName, prefix) {
				return true
			}
		}
	}
	return false
}

// matchS3KeyFilter applies the S3 key filter rules (prefix/suffix) from
// a NotificationConfigurationFilter against the given object key.
// Returns true when all filter rules are satisfied or when no filter is set.
func matchS3KeyFilter(key string, filter *s3store.NotificationConfigurationFilter) bool {
	if filter == nil || filter.Key == nil {
		return true
	}
	for _, rule := range filter.Key.FilterRules {
		switch strings.ToLower(rule.Name) {
		case "prefix":
			if !strings.HasPrefix(key, rule.Value) {
				return false
			}
		case "suffix":
			if !strings.HasSuffix(key, rule.Value) {
				return false
			}
		}
	}
	return true
}

// s3NotificationAttributes returns the MessageAttributes that AWS S3 attaches
// to every SNS and SQS notification delivery. They allow subscribers to filter
// or route messages without parsing the full JSON body.
func s3NotificationAttributes(event *eventbus.S3ObjectEvent) map[string]string {
	return map[string]string{
		"s3.bucket.name": event.Bucket,
		"s3.object.key":  event.Key,
		"s3.eventName":   "s3:" + string(event.Op),
	}
}

// s3SNSMessageAttributes wraps the flat string attributes into the
// json.RawMessage format required by SNSDeliveryEvent.
func s3SNSMessageAttributes(attrs map[string]string) map[string]json.RawMessage {
	result := make(map[string]json.RawMessage, len(attrs))
	for k, v := range attrs {
		result[k], _ = json.Marshal(map[string]interface{}{
			"Type":  "String",
			"Value": v,
		})
	}
	return result
}

// buildS3EventRecord constructs the AWS S3 event notification JSON payload
// matching the format documented at:
// https://docs.aws.amazon.com/AmazonS3/latest/userguide/EventBridge.html
// The output is a JSON envelope with a single-element "Records" array.
// configurationId is populated from the matching notification configuration.
func buildS3EventRecord(event *eventbus.S3ObjectEvent, configurationId string) []byte {
	record := map[string]interface{}{
		"eventVersion": "2.1",
		"eventSource":  "aws:s3",
		"awsRegion":    event.EventRegion(),
		"eventTime":    event.EventTimestamp().UTC().Format(time.RFC3339Nano),
		"eventName":    "s3:" + string(event.Op),
		"userIdentity": map[string]string{
			"principalId": event.EventAccountID(),
		},
		"requestParameters": map[string]string{
			"sourceIPAddress": event.SourceIP,
		},
		"s3": map[string]interface{}{
			"s3SchemaVersion": "1.0",
			"configurationId": configurationId,
			"bucket": map[string]interface{}{
				"name": event.Bucket,
				"arn":  svcarn.NewARNBuilder("", "").S3().Bucket(event.Bucket),
				"ownerIdentity": map[string]string{
					"principalId": event.EventAccountID(),
				},
			},
			"object": map[string]interface{}{
				"key":       event.Key,
				"size":      event.Size,
				"eTag":      event.ETag,
				"versionId": event.VersionID,
				"sequencer": fmt.Sprintf("%016X", event.EventTimestamp().UnixNano()),
			},
		},
	}

	// The restore completion record carries the glacierEventData extension
	// with the restored copy's expiry and storage class; the notification
	// structure documents it as visible only for s3:ObjectRestore:Completed
	// events.
	if event.Op == eventbus.S3ObjectRestoreCompleted && event.RestoreExpiry != nil {
		record["glacierEventData"] = map[string]interface{}{
			"restoreEventData": map[string]interface{}{
				"lifecycleRestorationExpiryTime": event.RestoreExpiry.UTC().Format("2006-01-02T15:04:05.000Z"),
				"lifecycleRestoreStorageClass":   event.RestoreStorageClass,
			},
		}
	}

	data, _ := json.Marshal(map[string]interface{}{
		"Records": []interface{}{record},
	})
	return data
}

// dispatchToSNS publishes the S3 event record to an SNS topic via the
// event bus. MessageAttributes are included so that SNS subscribers can
// filter by bucket, key, or event type without parsing the body.
func (s *S3Service) dispatchToSNS(ctx context.Context, topicArn string, payload []byte, attrs map[string]string) {
	if s.bus == nil {
		return
	}

	allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, topicArn, "sns", "s3.amazonaws.com", "sns:Publish", topicArn)
	if evalErr != nil || !allowed {
		return
	}

	_, _, region, accountID, _ := svcarn.SplitARN(topicArn)
	messageID := uuid.New().String()
	snsEvt := &eventbus.SNSDeliveryEvent{
		EventBase: eventbus.EventBase{
			Timestamp: time.Now().UTC(),
			Source:    "aws:s3",
			Region:    region,
			AccountID: accountID,
			Caller: eventbus.CallerContext{
				ServicePrincipal: "s3.amazonaws.com",
				AccountID:        accountID,
			},
		},
		TopicARN:          topicArn,
		MessageID:         messageID,
		Message:           string(payload),
		MessageAttributes: s3SNSMessageAttributes(attrs),
	}
	if err := s.bus.Publish(ctx, snsEvt); err != nil {
		logs.Warn("s3: failed to publish notification to SNS", logs.String("topicArn", topicArn), logs.Err(err))
	}
}

// dispatchToSQS sends the S3 event record directly to an SQS queue.
// MessageAttributes are included so that SQS consumers can filter by
// bucket, key, or event type.
func (s *S3Service) dispatchToSQS(ctx context.Context, queueArn string, payload []byte, attrs map[string]string) {
	if s.bus == nil {
		return
	}

	sqsInvoker := s.bus.SQSInvoker()
	if sqsInvoker == nil {
		return
	}

	queueName := svcarn.ExtractQueueNameFromARN(queueArn)
	if queueName == "" {
		return
	}

	_, _, sqsRegion, _, _ := svcarn.SplitARN(queueArn)

	queueURL, qErr := sqsInvoker.GetQueueByName(ctx, sqsRegion, queueName)
	if qErr != nil {
		return
	}

	allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, queueArn, "sqs", "s3.amazonaws.com", "sqs:SendMessage", queueArn)
	if evalErr != nil || !allowed {
		return
	}

	if _, _, err := sqsInvoker.SendMessage(ctx, sqsRegion, queueURL, string(payload), invokers.SQSSendOptions{
		MessageAttributes: attrs,
	}); err != nil {
		logs.Warn("Failed to send S3 event to SQS queue", logs.String("queue", queueURL), logs.Err(err))
	}
}

// dispatchToLambda invokes a Lambda function with the S3 event record
// as payload. The function name is extracted from the function ARN.
func (s *S3Service) dispatchToLambda(ctx context.Context, functionArn string, payload []byte) {
	if s.bus == nil {
		return
	}

	lambdaInvoker := s.bus.LambdaInvoker()
	if lambdaInvoker == nil {
		return
	}

	allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, functionArn, "lambda", "s3.amazonaws.com", "lambda:InvokeFunction", functionArn)
	if evalErr != nil || !allowed {
		return
	}

	if _, _, err := lambdaInvoker.InvokeForGateway(ctx, functionArn, payload); err != nil {
		logs.Warn("Failed to invoke Lambda for S3 event", logs.String("function", functionArn), logs.Err(err))
	}
}

// dispatchToEventBridge sends the S3 event to the default EventBridge event bus
// for the account. AWS S3 delivers events to EventBridge when
// EventBridgeConfiguration is set on the bucket notification configuration.
func (s *S3Service) dispatchToEventBridge(ctx context.Context, event *eventbus.S3ObjectEvent, eventName string) {
	if s.bus == nil {
		return
	}

	reason := eventName
	deletionType := ""
	if strings.HasPrefix(eventName, "s3:ObjectRemoved:") {
		parts := strings.SplitN(eventName, ":", 3)
		if len(parts) == 3 {
			deletionType = parts[2]
		}
	}

	sourceIP := event.SourceIP
	if sourceIP == "" {
		sourceIP = "127.0.0.1"
	}

	detail := map[string]interface{}{
		"version":       "0",
		"bucket":        map[string]string{"name": event.Bucket},
		"object":        map[string]interface{}{"key": event.Key, "size": event.Size, "etag": event.ETag, "version-id": event.VersionID},
		"request-id":    fmt.Sprintf("%016X", event.EventTimestamp().UnixNano()),
		"requester":     event.EventAccountID(),
		"source-ip":     sourceIP,
		"reason":        reason,
		"deletion-type": deletionType,
	}
	detailBytes, _ := json.Marshal(detail)

	ebEvt := &eventbus.EventBridgePutEventsEvent{
		EventBase: eventbus.EventBase{
			Timestamp: time.Now().UTC(),
			Source:    "aws.s3",
			Region:    event.EventRegion(),
			AccountID: event.EventAccountID(),
		},
		EventBusName: "default",
		Input:        string(detailBytes),
	}
	if err := s.bus.Publish(ctx, ebEvt); err != nil {
		logs.Warn("s3: failed to publish EventBridge event", logs.String("bucket", event.Bucket), logs.Err(err))
	}
}
