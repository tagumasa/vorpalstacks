package s3

import (
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

	config := &s3store.NotificationConfiguration{}

	for _, tc := range input.NotificationConfiguration.TopicConfigurations {
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
				topicConfig.Filter.Key.FilterRules = append(topicConfig.Filter.Key.FilterRules, s3store.FilterRule{
					Name:  fr.Name,
					Value: fr.Value,
				})
			}
		}
		config.TopicConfigurations = append(config.TopicConfigurations, topicConfig)
	}

	for _, qc := range input.NotificationConfiguration.QueueConfigurations {
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
				queueConfig.Filter.Key.FilterRules = append(queueConfig.Filter.Key.FilterRules, s3store.FilterRule{
					Name:  fr.Name,
					Value: fr.Value,
				})
			}
		}
		config.QueueConfigurations = append(config.QueueConfigurations, queueConfig)
	}

	for _, lc := range input.NotificationConfiguration.LambdaConfigurations {
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
