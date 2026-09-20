package sns

import (
	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
)

// SNSStoreInterface defines operations for managing SNS topics and subscriptions.
type SNSStoreInterface interface {
	// CreateTopic commits a new topic record and its initial tag set as
	// one transaction. tags ride the input only — the Topic record itself
	// carries no tag state (the TagStore owns tags).
	CreateTopic(topic *Topic, tags map[string]string) (*Topic, error)
	GetTopic(topicArn string) (*Topic, error)
	DeleteTopic(topicArn string) error
	ListTopics(opts common.ListOptions) (*common.ListResult[Topic], error)
	SetTopicAttributes(topicArn string, attributes map[string]string) error

	// CreateSubscription is idempotent on the subscription natural key
	// (TopicArn+Protocol+Endpoint+Owner): a matching live subscription is
	// returned with created=false instead of a duplicate record.
	CreateSubscription(subscription *Subscription) (sub *Subscription, created bool, err error)
	GetSubscription(subscriptionArn string) (*Subscription, error)
	DeleteSubscription(subscriptionArn string) error
	AutoConfirmSubscription(subscription *Subscription) error
	ConfirmSubscription(subscriptionArn, token string, authenticateOnUnsubscribe *bool) (*Subscription, error)
	ListSubscriptions(opts common.ListOptions) (*common.ListResult[Subscription], error)
	ListSubscriptionsByTopic(topicArn string, opts common.ListOptions) (*common.ListResult[Subscription], error)
	ListAllSubscriptionsByTopic(topicArn string) ([]*Subscription, error)
	FindSubscriptionByToken(topicArn, token string) (*Subscription, error)
	SetSubscriptionAttributes(subscriptionArn string, attributes map[string]string) error
	GetSubscriptionAttributes(subscriptionArn string) (map[string]string, error)

	ListTagsForResource(resourceArn string) ([]types.Tag, error)
	Tag(resourceArn string, tags []types.Tag) error
	Untag(resourceArn string, tagKeys []string) error

	AcquireFifoGroupLock(topicArn, messageGroupId string) func()
	CheckFifoDeduplication(topicArn, messageGroupId, messageDeduplicationId string, perGroupDedup bool) (messageID, sequenceNumber string, hit bool, err error)
	CheckAndRecordFifoDeduplication(topicArn, messageGroupId, messageDeduplicationId, messageID, sequenceNumber string, perGroupDedup bool) (string, string, bool, error)
	AllocateFifoSequence(topicArn, messageGroupId string) (string, error)

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

	Close()
}
