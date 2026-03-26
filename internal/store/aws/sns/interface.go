package sns

import (
	"vorpalstacks/internal/store/aws/common"
)

// SNSStoreInterface defines operations for managing SNS topics and subscriptions.
type SNSStoreInterface interface {
	CreateTopic(topic *Topic) (*Topic, error)
	GetTopic(topicArn string) (*Topic, error)
	GetTopicByName(topicName string) (*Topic, error)
	UpdateTopic(topic *Topic) error
	DeleteTopic(topicArn string) error
	ListTopics(opts common.ListOptions) (*common.ListResult[Topic], error)
	SetTopicAttributes(topicArn string, attributes map[string]string) error

	CreateSubscription(subscription *Subscription) (*Subscription, error)
	GetSubscription(subscriptionArn string) (*Subscription, error)
	DeleteSubscription(subscriptionArn string) error
	UpdateSubscription(subscription *Subscription) error
	AutoConfirmSubscription(subscription *Subscription) error
	ConfirmSubscription(subscriptionArn, token string) (*Subscription, error)
	ListSubscriptions(opts common.ListOptions) (*common.ListResult[Subscription], error)
	ListSubscriptionsByTopic(topicArn string, opts common.ListOptions) (*common.ListResult[Subscription], error)
	SetSubscriptionAttributes(subscriptionArn string, attributes map[string]string) error
	GetSubscriptionAttributes(subscriptionArn string) (map[string]string, error)

	ListTagsForResource(resourceArn string) ([]common.Tag, error)
	TagResource(resourceArn string, tags []common.Tag) error
	UntagResource(resourceArn string, tagKeys []string) error

	CheckDeduplication(topicArn, messageDeduplicationId string) (string, bool)
	RecordDeduplication(topicArn, messageDeduplicationId, messageID string)
	GetNextSequenceNumber(topicArn, messageGroupId string) string

	GetDataProtectionPolicy(topicArn string) (string, error)
	PutDataProtectionPolicy(topicArn, policy string) error

	AddPermission(topicArn string, permission *Permission) error
	RemovePermission(topicArn, label string) error

	CreatePlatformApplication(app *PlatformApplication) (*PlatformApplication, error)
	GetPlatformApplication(arn string) (*PlatformApplication, error)
	DeletePlatformApplication(arn string) error
	ListPlatformApplications(opts common.ListOptions) (*common.ListResult[PlatformApplication], error)
	SetPlatformApplicationAttributes(arn string, attributes map[string]string) error
	GetPlatformApplicationAttributes(arn string) (map[string]string, error)

	CreatePlatformEndpoint(endpoint *PlatformEndpoint) (*PlatformEndpoint, error)
	GetEndpoint(arn string) (*PlatformEndpoint, error)
	DeleteEndpoint(arn string) error
	GetEndpointAttributes(arn string) (map[string]string, error)
	SetEndpointAttributes(arn string, attributes map[string]string) error
	ListEndpointsByPlatformApplication(platformAppArn string, opts common.ListOptions) (*common.ListResult[PlatformEndpoint], error)

	Put(key string, data interface{}) error

	Close()
}
