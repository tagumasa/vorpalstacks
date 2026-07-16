package sns

// Package sns provides SNS (Simple Notification Service) data store implementations
// for vorpalstacks.

import (
	"encoding/json"
	"strings"
	"time"
)

// Topic represents an SNS topic.
type Topic struct {
	Name                      string            `json:"name"`
	Arn                       string            `json:"arn"`
	DisplayName               string            `json:"display_name,omitempty"`
	Policy                    string            `json:"policy,omitempty"`
	DeliveryPolicy            string            `json:"delivery_policy,omitempty"`
	EffectiveDeliveryPolicy   string            `json:"effective_delivery_policy,omitempty"`
	Owner                     string            `json:"owner,omitempty"`
	SubscriptionsConfirmed    int32             `json:"subscriptions_confirmed"`
	SubscriptionsDeleted      int32             `json:"subscriptions_deleted"`
	SubscriptionsPending      int32             `json:"subscriptions_pending"`
	KmsMasterKeyId            string            `json:"kms_master_key_id,omitempty"`
	FifoTopic                 bool              `json:"fifo_topic"`
	ContentBasedDeduplication bool              `json:"content_based_deduplication"`
	SignatureVersion          string            `json:"signature_version,omitempty"`
	CreatedDate               time.Time         `json:"created_date"`
	LastModifiedTime          time.Time         `json:"last_modified_time"`
	Attributes                map[string]string `json:"attributes,omitempty"`
	Tags                      map[string]string `json:"tags,omitempty"`
	DataProtectionPolicy      string            `json:"data_protection_policy,omitempty"`
	Permissions               []Permission      `json:"permissions,omitempty"`
}

// Subscription represents an SNS subscription.
// Subscription attributes (RawMessageDelivery, FilterPolicy, RedrivePolicy,
// etc.) are stored exclusively in the Attributes map — the single source of
// truth — matching the AWS SNS SetSubscriptionAttributes API. Typed accessors
// below provide convenient, type-safe reads without duplicating state.
type Subscription struct {
	SubscriptionArn              string            `json:"subscription_arn"`
	TopicArn                     string            `json:"topic_arn"`
	Protocol                     string            `json:"protocol"`
	Endpoint                     string            `json:"endpoint"`
	Owner                        string            `json:"owner"`
	ConfirmationWasAuthenticated bool              `json:"confirmation_was_authenticated"`
	TopicOwner                   string            `json:"topic_owner"`
	PendingConfirmation          bool              `json:"pending_confirmation"`
	ConfirmationToken            string            `json:"confirmation_token,omitempty"`
	SubscriptionPrincipal        string            `json:"subscription_principal,omitempty"`
	Attributes                   map[string]string `json:"attributes,omitempty"`
	CreatedDate                  time.Time         `json:"created_date"`
}

// IsRawMessageDelivery returns whether raw message delivery is enabled.
func (s *Subscription) IsRawMessageDelivery() bool {
	return strings.EqualFold(s.Attributes["RawMessageDelivery"], "true")
}

// GetFilterPolicy returns the subscription's filter policy JSON string.
func (s *Subscription) GetFilterPolicy() string {
	return s.Attributes["FilterPolicy"]
}

// GetFilterPolicyScope returns the filter policy scope, defaulting to
// "MessageAttributes" when unset (matching AWS behaviour).
func (s *Subscription) GetFilterPolicyScope() string {
	scope := s.Attributes["FilterPolicyScope"]
	if scope == "" {
		return "MessageAttributes"
	}
	return scope
}

// GetRedrivePolicy parses the subscription's redrive policy from the
// Attributes map. Returns nil when no redrive policy is set.
func (s *Subscription) GetRedrivePolicy() (*RedrivePolicy, error) {
	raw := s.Attributes["RedrivePolicy"]
	if raw == "" {
		return nil, nil
	}
	var rp RedrivePolicy
	if err := json.Unmarshal([]byte(raw), &rp); err != nil {
		return nil, err
	}
	return &rp, nil
}

// RedrivePolicy represents the redrive policy for an SNS subscription.
// SNS subscription RedrivePolicy only contains deadLetterTargetArn.
// (maxReceiveCount is an SQS queue concept, not SNS.)
type RedrivePolicy struct {
	DeadLetterTargetArn string `json:"deadLetterTargetArn"`
}

// Message represents an SNS message.
type Message struct {
	MessageId              string                       `json:"message_id"`
	TopicArn               string                       `json:"topic_arn"`
	Subject                string                       `json:"subject,omitempty"`
	Message                string                       `json:"message"`
	MessageStructure       string                       `json:"message_structure,omitempty"`
	MessageAttributes      map[string]*MessageAttribute `json:"message_attributes,omitempty"`
	ReceivedTimestamp      time.Time                    `json:"received_timestamp"`
	PublishedTimestamp     time.Time                    `json:"published_timestamp"`
	MessageGroupId         string                       `json:"message_group_id,omitempty"`
	MessageDeduplicationId string                       `json:"message_deduplication_id,omitempty"`
}

// MessageAttribute represents an SNS message attribute.
type MessageAttribute struct {
	Type        string `json:"type"`
	StringValue string `json:"string_value,omitempty"`
	BinaryValue []byte `json:"binary_value,omitempty"`
}

// Permission represents an SNS topic permission.
type Permission struct {
	Label       string   `json:"label"`
	AccountId   string   `json:"account_id,omitempty"`
	ActionName  string   `json:"action_name,omitempty"`
	AWAccountId string   `json:"aws_account_id,omitempty"`
	Principals  []string `json:"principals,omitempty"`
	Actions     []string `json:"actions,omitempty"`
}

// PlatformApplication represents an SNS platform application.
type PlatformApplication struct {
	PlatformApplicationArn string            `json:"platform_application_arn"`
	Name                   string            `json:"name,omitempty"`
	Platform               string            `json:"platform,omitempty"`
	Attributes             map[string]string `json:"attributes,omitempty"`
}

// PlatformEndpoint represents an SNS platform endpoint.
type PlatformEndpoint struct {
	EndpointArn            string            `json:"endpoint_arn"`
	PlatformApplicationArn string            `json:"platform_application_arn,omitempty"`
	Token                  string            `json:"token,omitempty"`
	CustomUserData         string            `json:"custom_user_data,omitempty"`
	Attributes             map[string]string `json:"attributes,omitempty"`
}

// NewTopic creates a new Topic with the specified name and ARN.
func NewTopic(name, arn string) *Topic {
	return &Topic{
		Name:       name,
		Arn:        arn,
		Attributes: make(map[string]string),
	}
}

// NewSubscription creates a new Subscription with the specified parameters.
func NewSubscription(subscriptionArn, topicArn, protocol, endpoint string) *Subscription {
	return &Subscription{
		SubscriptionArn:     subscriptionArn,
		TopicArn:            topicArn,
		Protocol:            protocol,
		Endpoint:            endpoint,
		Attributes:          make(map[string]string),
		PendingConfirmation: true,
	}
}

// NewMessage creates a new Message with the specified message ID, topic ARN, and message body.
func NewMessage(messageId, topicArn, message string) *Message {
	return &Message{
		MessageId:         messageId,
		TopicArn:          topicArn,
		Message:           message,
		MessageAttributes: make(map[string]*MessageAttribute),
	}
}
