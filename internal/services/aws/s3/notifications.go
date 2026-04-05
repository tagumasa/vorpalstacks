// Package s3 provides S3 service operations for vorpalstacks.
package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/server/eventbus"
	s3store "vorpalstacks/internal/store/aws/s3"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
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
	eventRecord := buildS3EventRecord(event)

	for _, tc := range config.TopicConfigurations {
		if !matchS3Event(tc.Events, eventName) {
			continue
		}
		if tc.Filter != nil && !matchS3KeyFilter(event.Key, tc.Filter) {
			continue
		}
		s.dispatchToSNS(ctx, tc.TopicArn, eventRecord)
	}

	for _, qc := range config.QueueConfigurations {
		if !matchS3Event(qc.Events, eventName) {
			continue
		}
		if qc.Filter != nil && !matchS3KeyFilter(event.Key, qc.Filter) {
			continue
		}
		s.dispatchToSQS(ctx, qc.QueueArn, eventRecord)
	}

	for _, lc := range config.LambdaConfigurations {
		if !matchS3Event(lc.Events, eventName) {
			continue
		}
		if lc.Filter != nil && !matchS3KeyFilter(event.Key, lc.Filter) {
			continue
		}
		s.dispatchToLambda(ctx, lc.LambdaFunctionArn, eventRecord)
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

// buildS3EventRecord constructs the AWS S3 event notification JSON payload
// matching the format documented at:
// https://docs.aws.amazon.com/AmazonS3/latest/userguide/EventBridge.html
// The output is a JSON envelope with a single-element "Records" array.
func buildS3EventRecord(event *eventbus.S3ObjectEvent) []byte {
	record := map[string]interface{}{
		"eventVersion": "2.1",
		"eventSource":  "aws:s3",
		"awsRegion":    event.EventRegion(),
		"eventTime":    event.EventTimestamp().UTC().Format(time.RFC3339Nano),
		"eventName":    string(event.Op),
		"userIdentity": map[string]string{
			"principalId": event.EventAccountID(),
		},
		"requestParameters": map[string]string{
			"sourceIPAddress": "127.0.0.1",
		},
		"s3": map[string]interface{}{
			"s3SchemaVersion": "1.0",
			"configurationId": "",
			"bucket": map[string]interface{}{
				"name": event.Bucket,
				"arn":  fmt.Sprintf("arn:aws:s3:::%s", event.Bucket),
				"ownerIdentity": map[string]string{
					"principalId": event.EventAccountID(),
				},
			},
			"object": map[string]interface{}{
				"key":       event.Key,
				"size":      event.Size,
				"eTag":      event.ETag,
				"versionId": event.VersionID,
			},
		},
	}

	data, _ := json.Marshal(map[string]interface{}{
		"Records": []interface{}{record},
	})
	return data
}

// dispatchToSNS publishes the S3 event record to an SNS topic via the
// event bus. Region and account ID are extracted from the topic ARN to
// support cross-region notification configurations.
func (s *S3Service) dispatchToSNS(ctx context.Context, topicArn string, payload []byte) {
	if s.bus == nil {
		return
	}

	allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, topicArn, "sns", "s3.amazonaws.com", "sns:Publish", topicArn)
	if evalErr != nil {
		fmt.Printf("[s3:notification] policy evaluation failed for SNS topic=%s: %v\n", topicArn, evalErr)
		return
	}
	if !allowed {
		return
	}

	_, _, region, accountID, _ := svcarn.SplitARN(topicArn)
	messageID := fmt.Sprintf("%d", time.Now().UnixNano())
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
		TopicARN:  topicArn,
		MessageID: messageID,
		Message:   string(payload),
	}
	_ = s.bus.Publish(ctx, snsEvt)
}

// dispatchToSQS sends the S3 event record directly to an SQS queue.
// Region and account ID are extracted from the queue ARN so that the
// correct queue URL is constructed for the target region.
func (s *S3Service) dispatchToSQS(ctx context.Context, queueArn string, payload []byte) {
	if s.sqsStore == nil {
		return
	}

	queueName := svcarn.ExtractQueueNameFromARN(queueArn)
	if queueName == "" {
		return
	}

	queue, qErr := s.sqsStore.GetQueueByName(queueName)
	if qErr != nil {
		fmt.Printf("[s3:notification] queue not found for SQS delivery queue=%s: %v\n", queueName, qErr)
		return
	}

	if s.bus != nil {
		allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, queue.ARN, "sqs", "s3.amazonaws.com", "sqs:SendMessage", queue.ARN)
		if evalErr != nil {
			fmt.Printf("[s3:notification] policy evaluation failed for SQS queue=%s: %v\n", queueArn, evalErr)
			return
		}
		if !allowed {
			return
		}
	}

	message := &sqsstore.Message{
		Body: string(payload),
	}
	if _, err := s.sqsStore.SendMessage(queue.URL, message); err != nil {
		fmt.Printf("[s3:notification] sqs dispatch failed queue=%s: %v\n", queueName, err)
	}
}

// dispatchToLambda invokes a Lambda function with the S3 event record
// as payload. The function name is extracted from the function ARN.
func (s *S3Service) dispatchToLambda(ctx context.Context, functionArn string, payload []byte) {
	if s.lambdaInvoker == nil {
		return
	}

	if s.bus != nil {
		allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, functionArn, "lambda", "s3.amazonaws.com", "lambda:InvokeFunction", functionArn)
		if evalErr != nil {
			fmt.Printf("[s3:notification] policy evaluation failed for Lambda function=%s: %v\n", functionArn, evalErr)
			return
		}
		if !allowed {
			return
		}
	}

	fnName := svcarn.ExtractFunctionNameFromARN(functionArn)
	_, _, err := s.lambdaInvoker.InvokeForGateway(ctx, functionArn, payload)
	if err != nil {
		fmt.Printf("[s3:notification] lambda dispatch failed function=%s: %v\n", fnName, err)
	}
}
