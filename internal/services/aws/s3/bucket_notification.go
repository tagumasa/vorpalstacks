package s3

import (
	"context"
	"fmt"
	"strings"

	"vorpalstacks/internal/common/request"
)

// PutBucketNotificationInput contains the request parameters for the PutBucketNotificationConfiguration operation.
type PutBucketNotificationInput struct {
	Bucket                    string
	NotificationConfiguration *NotificationConfigurationInput
}

// NotificationConfigurationInput defines the notification configuration for a bucket.
type NotificationConfigurationInput struct {
	TopicConfigurations  []TopicConfigurationInput  `xml:"TopicConfiguration,omitempty"`
	QueueConfigurations  []QueueConfigurationInput  `xml:"QueueConfiguration,omitempty"`
	LambdaConfigurations []LambdaConfigurationInput `xml:"CloudFunctionConfiguration,omitempty"`
}

// TopicConfigurationInput defines a topic notification configuration.
type TopicConfigurationInput struct {
	Id       string                   `xml:"Id,omitempty"`
	TopicArn string                   `xml:"Topic"`
	Events   []string                 `xml:"Event"`
	Filter   *NotificationFilterInput `xml:"Filter,omitempty"`
}

// QueueConfigurationInput defines a queue notification configuration.
type QueueConfigurationInput struct {
	Id       string                   `xml:"Id,omitempty"`
	QueueArn string                   `xml:"Queue"`
	Events   []string                 `xml:"Event"`
	Filter   *NotificationFilterInput `xml:"Filter,omitempty"`
}

// LambdaConfigurationInput defines a Lambda notification configuration.
type LambdaConfigurationInput struct {
	Id                string                   `xml:"Id,omitempty"`
	LambdaFunctionArn string                   `xml:"CloudFunction"`
	Events            []string                 `xml:"Event"`
	Filter            *NotificationFilterInput `xml:"Filter,omitempty"`
}

// NotificationFilterInput defines filter criteria for notifications.
type NotificationFilterInput struct {
	S3Key *S3KeyFilterInput `xml:"S3Key,omitempty"`
}

// S3KeyFilterInput defines filter rules for S3 object keys.
type S3KeyFilterInput struct {
	FilterRules []FilterRuleInput `xml:"FilterRule,omitempty"`
}

// FilterRuleInput defines a single filter rule.
type FilterRuleInput struct {
	Name  string `xml:"Name"`
	Value string `xml:"Value"`
}

// PutBucketNotificationConfiguration configures event notifications for a bucket.
func (o *BucketOperations) PutBucketNotificationConfiguration(ctx *request.RequestContext, input *PutBucketNotificationInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.putBucketNotificationConfigurationCore(ctx, store.buckets, input)
}

// GetBucketNotificationInput contains the request parameters for the GetBucketNotificationConfiguration operation.
type GetBucketNotificationInput struct {
	Bucket string
}

// GetBucketNotificationOutput contains the result of the GetBucketNotificationConfiguration operation.
type GetBucketNotificationOutput struct {
	NotificationConfiguration *NotificationConfigurationOutput `xml:"NotificationConfiguration"`
}

// NotificationConfigurationOutput defines the notification configuration for a bucket.
type NotificationConfigurationOutput struct {
	TopicConfigurations  []TopicConfigurationOutput  `xml:"TopicConfiguration,omitempty"`
	QueueConfigurations  []QueueConfigurationOutput  `xml:"QueueConfiguration,omitempty"`
	LambdaConfigurations []LambdaConfigurationOutput `xml:"CloudFunctionConfiguration,omitempty"`
}

// TopicConfigurationOutput defines a topic notification configuration.
type TopicConfigurationOutput struct {
	Id       string                    `xml:"Id,omitempty"`
	TopicArn string                    `xml:"Topic"`
	Events   []string                  `xml:"Event"`
	Filter   *NotificationFilterOutput `xml:"Filter,omitempty"`
}

// QueueConfigurationOutput defines a queue notification configuration.
type QueueConfigurationOutput struct {
	Id       string                    `xml:"Id,omitempty"`
	QueueArn string                    `xml:"Queue"`
	Events   []string                  `xml:"Event"`
	Filter   *NotificationFilterOutput `xml:"Filter,omitempty"`
}

// LambdaConfigurationOutput defines a Lambda notification configuration.
type LambdaConfigurationOutput struct {
	Id                string                    `xml:"Id,omitempty"`
	LambdaFunctionArn string                    `xml:"CloudFunction"`
	Events            []string                  `xml:"Event"`
	Filter            *NotificationFilterOutput `xml:"Filter,omitempty"`
}

// NotificationFilterOutput defines filter criteria for notifications.
type NotificationFilterOutput struct {
	S3Key *S3KeyFilterOutput `xml:"S3Key,omitempty"`
}

// S3KeyFilterOutput defines filter rules for S3 object keys.
type S3KeyFilterOutput struct {
	FilterRules []FilterRuleOutput `xml:"FilterRule,omitempty"`
}

// FilterRuleOutput defines a single filter rule.
type FilterRuleOutput struct {
	Name  string `xml:"Name"`
	Value string `xml:"Value"`
}

// GetBucketNotificationConfiguration retrieves the notification configuration for a bucket.
func (o *BucketOperations) GetBucketNotificationConfiguration(ctx *request.RequestContext, input *GetBucketNotificationInput) (*GetBucketNotificationOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	return o.svc.getBucketNotificationConfigurationCore(store.buckets, input)
}

// validateNotificationTarget checks whether the target resource (SNS topic,
// SQS queue, or Lambda function) exists via the event bus invokers.  If the
// event bus or the relevant invoker is not available, the check is skipped.
func (s *S3Service) validateNotificationTarget(ctx context.Context, arn, serviceType string) error {
	if s.bus == nil {
		return nil
	}

	switch serviceType {
	case "sns":
		snsInvoker := s.bus.SNSInvoker()
		if snsInvoker == nil {
			return nil
		}
		if _, err := snsInvoker.GetTopic(ctx, arn); err != nil {
			return NewInvalidArgumentError(fmt.Sprintf("SNS topic does not exist: %s", arn))
		}
	case "sqs":
		sqsInvoker := s.bus.SQSInvoker()
		if sqsInvoker == nil {
			return nil
		}
		parts := strings.SplitN(arn, ":", 6)
		if len(parts) < 6 {
			return nil
		}
		if _, err := sqsInvoker.GetQueueByName(ctx, parts[3], parts[5]); err != nil {
			return NewInvalidArgumentError(fmt.Sprintf("SQS queue does not exist: %s", arn))
		}
	case "lambda":
		lambdaInvoker := s.bus.LambdaInvoker()
		if lambdaInvoker == nil {
			return nil
		}
		if _, err := lambdaInvoker.GetFunctionARN(ctx, arn); err != nil {
			return NewInvalidArgumentError(fmt.Sprintf("Lambda function does not exist: %s", arn))
		}
	}
	return nil
}
