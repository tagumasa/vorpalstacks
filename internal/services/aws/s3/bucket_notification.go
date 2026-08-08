package s3

import (
	"context"
	"fmt"
	"strings"

	"vorpalstacks/internal/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
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

	// Track Id uniqueness across all configurations.
	seenIds := make(map[string]bool)

	validateId := func(id string) error {
		if id == "" {
			return nil
		}
		if seenIds[id] {
			return NewInvalidArgumentError(fmt.Sprintf("duplicate notification configuration Id: %s", id))
		}
		seenIds[id] = true
		return nil
	}

	validateArn := func(arn, expectedService string) error {
		if arn == "" {
			return NewInvalidArgumentError("notification ARN is required")
		}
		parts := strings.SplitN(arn, ":", 6)
		if len(parts) < 6 || parts[0] != "arn" {
			return NewInvalidArgumentError(fmt.Sprintf("invalid ARN format: %s", arn))
		}
		if parts[2] != expectedService {
			return NewInvalidArgumentError(fmt.Sprintf("expected ARN service %s, got %s in: %s", expectedService, parts[2], arn))
		}
		return nil
	}

	validateEvents := func(events []string) error {
		return validateS3EventNames(events)
	}

	config := &s3store.NotificationConfiguration{}

	for _, tc := range input.NotificationConfiguration.TopicConfigurations {
		if err := validateId(tc.Id); err != nil {
			return err
		}
		if err := validateArn(tc.TopicArn, "sns"); err != nil {
			return err
		}
		if err := o.svc.validateNotificationTarget(ctx, tc.TopicArn, "sns"); err != nil {
			return err
		}
		if err := validateEvents(tc.Events); err != nil {
			return err
		}
		topicConfig := s3store.TopicNotificationConfiguration{
			Id:       tc.Id,
			TopicArn: tc.TopicArn,
			Events:   tc.Events,
		}
		if tc.Filter != nil && tc.Filter.S3Key != nil {
			topicConfig.Filter = &s3store.NotificationConfigurationFilter{
				Key: &s3store.S3KeyFilter{},
			}
			for _, fr := range tc.Filter.S3Key.FilterRules {
				if err := validateFilterRule(fr.Name, fr.Value); err != nil {
					return err
				}
				topicConfig.Filter.Key.FilterRules = append(topicConfig.Filter.Key.FilterRules, s3store.FilterRule{
					Name:  fr.Name,
					Value: fr.Value,
				})
			}
		}
		config.TopicConfigurations = append(config.TopicConfigurations, topicConfig)
	}

	for _, qc := range input.NotificationConfiguration.QueueConfigurations {
		if err := validateId(qc.Id); err != nil {
			return err
		}
		if err := validateArn(qc.QueueArn, "sqs"); err != nil {
			return err
		}
		if err := o.svc.validateNotificationTarget(ctx, qc.QueueArn, "sqs"); err != nil {
			return err
		}
		if err := validateEvents(qc.Events); err != nil {
			return err
		}
		queueConfig := s3store.QueueNotificationConfiguration{
			Id:       qc.Id,
			QueueArn: qc.QueueArn,
			Events:   qc.Events,
		}
		if qc.Filter != nil && qc.Filter.S3Key != nil {
			queueConfig.Filter = &s3store.NotificationConfigurationFilter{
				Key: &s3store.S3KeyFilter{},
			}
			for _, fr := range qc.Filter.S3Key.FilterRules {
				if err := validateFilterRule(fr.Name, fr.Value); err != nil {
					return err
				}
				queueConfig.Filter.Key.FilterRules = append(queueConfig.Filter.Key.FilterRules, s3store.FilterRule{
					Name:  fr.Name,
					Value: fr.Value,
				})
			}
		}
		config.QueueConfigurations = append(config.QueueConfigurations, queueConfig)
	}

	for _, lc := range input.NotificationConfiguration.LambdaConfigurations {
		if err := validateId(lc.Id); err != nil {
			return err
		}
		if err := validateArn(lc.LambdaFunctionArn, "lambda"); err != nil {
			return err
		}
		if err := o.svc.validateNotificationTarget(ctx, lc.LambdaFunctionArn, "lambda"); err != nil {
			return err
		}
		if err := validateEvents(lc.Events); err != nil {
			return err
		}
		lambdaConfig := s3store.LambdaNotificationConfiguration{
			Id:                lc.Id,
			LambdaFunctionArn: lc.LambdaFunctionArn,
			Events:            lc.Events,
		}
		if lc.Filter != nil && lc.Filter.S3Key != nil {
			lambdaConfig.Filter = &s3store.NotificationConfigurationFilter{
				Key: &s3store.S3KeyFilter{},
			}
			for _, fr := range lc.Filter.S3Key.FilterRules {
				if err := validateFilterRule(fr.Name, fr.Value); err != nil {
					return err
				}
				lambdaConfig.Filter.Key.FilterRules = append(lambdaConfig.Filter.Key.FilterRules, s3store.FilterRule{
					Name:  fr.Name,
					Value: fr.Value,
				})
			}
		}
		config.LambdaConfigurations = append(config.LambdaConfigurations, lambdaConfig)
	}

	return store.buckets.SetNotificationConfiguration(input.Bucket, config)
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

	config, err := store.buckets.GetNotificationConfiguration(input.Bucket)
	if err != nil {
		return nil, err
	}

	output := &GetBucketNotificationOutput{
		NotificationConfiguration: &NotificationConfigurationOutput{},
	}

	if config == nil {
		return output, nil
	}

	for _, tc := range config.TopicConfigurations {
		topicOut := TopicConfigurationOutput{
			Id:       tc.Id,
			TopicArn: tc.TopicArn,
			Events:   tc.Events,
		}
		if tc.Filter != nil && tc.Filter.Key != nil {
			topicOut.Filter = &NotificationFilterOutput{S3Key: &S3KeyFilterOutput{}}
			for _, fr := range tc.Filter.Key.FilterRules {
				topicOut.Filter.S3Key.FilterRules = append(topicOut.Filter.S3Key.FilterRules, FilterRuleOutput{
					Name:  fr.Name,
					Value: fr.Value,
				})
			}
		}
		output.NotificationConfiguration.TopicConfigurations = append(output.NotificationConfiguration.TopicConfigurations, topicOut)
	}

	for _, qc := range config.QueueConfigurations {
		queueOut := QueueConfigurationOutput{
			Id:       qc.Id,
			QueueArn: qc.QueueArn,
			Events:   qc.Events,
		}
		if qc.Filter != nil && qc.Filter.Key != nil {
			queueOut.Filter = &NotificationFilterOutput{S3Key: &S3KeyFilterOutput{}}
			for _, fr := range qc.Filter.Key.FilterRules {
				queueOut.Filter.S3Key.FilterRules = append(queueOut.Filter.S3Key.FilterRules, FilterRuleOutput{
					Name:  fr.Name,
					Value: fr.Value,
				})
			}
		}
		output.NotificationConfiguration.QueueConfigurations = append(output.NotificationConfiguration.QueueConfigurations, queueOut)
	}

	for _, lc := range config.LambdaConfigurations {
		lambdaOut := LambdaConfigurationOutput{
			Id:                lc.Id,
			LambdaFunctionArn: lc.LambdaFunctionArn,
			Events:            lc.Events,
		}
		if lc.Filter != nil && lc.Filter.Key != nil {
			lambdaOut.Filter = &NotificationFilterOutput{S3Key: &S3KeyFilterOutput{}}
			for _, fr := range lc.Filter.Key.FilterRules {
				lambdaOut.Filter.S3Key.FilterRules = append(lambdaOut.Filter.S3Key.FilterRules, FilterRuleOutput{
					Name:  fr.Name,
					Value: fr.Value,
				})
			}
		}
		output.NotificationConfiguration.LambdaConfigurations = append(output.NotificationConfiguration.LambdaConfigurations, lambdaOut)
	}

	return output, nil
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
