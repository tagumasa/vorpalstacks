package sns

// Package sns provides SNS (Simple Notification Service) data store implementations
// for vorpalstacks.

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// Topic represents an SNS topic. Tags are owned by the TagStore and never
// ride this record; the create-time tag set travels on the CreateTopic
// input instead.
type Topic struct {
	Name                   string            `json:"name"`
	Arn                    string            `json:"arn"`
	Owner                  string            `json:"owner,omitempty"`
	SubscriptionsConfirmed int32             `json:"subscriptions_confirmed"`
	SubscriptionsDeleted   int32             `json:"subscriptions_deleted"`
	SubscriptionsPending   int32             `json:"subscriptions_pending"`
	Attributes             map[string]string `json:"attributes,omitempty"`
	Permissions            []Permission      `json:"permissions,omitempty"`
}

// IsFifoTopic returns whether the topic is a FIFO topic.
func (t *Topic) IsFifoTopic() bool {
	return strings.HasSuffix(t.Name, ".fifo")
}

// IsContentBasedDeduplication returns whether content-based deduplication is enabled.
func (t *Topic) IsContentBasedDeduplication() bool {
	return strings.EqualFold(t.Attributes[AttrContentBasedDedup], "true")
}

// PerGroupDeduplication reports whether FIFO deduplication is scoped to
// each individual message group: "Message deduplication applies to an
// entire Amazon SNS FIFO topic when the topic attribute FifoThroughputScope
// is set to Topic. When the topic attribute FifoThroughputScope is set to
// MessageGroup, message deduplication applies to each individual message
// group" (message deduplication for FIFO topics).
func (t *Topic) PerGroupDeduplication() bool {
	return t.Attributes[AttrFifoThroughputScope] == "MessageGroup"
}

// MaximumMessageSizeBytes returns the topic's effective message-size
// ceiling: the MaximumMessageSize attribute when the topic sets one, the
// platform's flat cap otherwise (the attribute's documented default is
// 262144 — the flat cap — so an unset attribute is the flat cap).
func (t *Topic) MaximumMessageSizeBytes() int {
	if v, ok := t.Attributes[AttrMaximumMessageSize]; ok {
		if size, err := strconv.Atoi(v); err == nil && size > 0 {
			return size
		}
	}
	return MaxMessageSize
}

// GetPolicy returns the topic's access policy JSON.
func (t *Topic) GetPolicy() string {
	return t.Attributes[AttrPolicy]
}

// GetDataProtectionPolicy returns the topic's data protection policy.
func (t *Topic) GetDataProtectionPolicy() string {
	return t.Attributes[AttrDataProtectionPolicy]
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
	PendingConfirmation          bool              `json:"pending_confirmation"`
	ConfirmationToken            string            `json:"confirmation_token,omitempty"`
	Attributes                   map[string]string `json:"attributes,omitempty"`
}

// IsRawMessageDelivery returns whether raw message delivery is enabled.
func (s *Subscription) IsRawMessageDelivery() bool {
	return strings.EqualFold(s.Attributes[AttrRawMessageDelivery], "true")
}

// GetFilterPolicy returns the subscription's filter policy JSON string.
func (s *Subscription) GetFilterPolicy() string {
	return s.Attributes[AttrFilterPolicy]
}

// GetFilterPolicyScope returns the filter policy scope, defaulting to
// "MessageAttributes" when unset (matching AWS behaviour).
func (s *Subscription) GetFilterPolicyScope() string {
	scope := s.Attributes[AttrFilterPolicyScope]
	if scope == "" {
		return FilterPolicyScopeAttributes
	}
	return scope
}

// GetRedrivePolicy parses the subscription's redrive policy from the
// Attributes map. Returns nil when no redrive policy is set.
func (s *Subscription) GetRedrivePolicy() (*RedrivePolicy, error) {
	raw := s.Attributes[AttrRedrivePolicy]
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

// Message represents an SNS message in flight through the delivery
// engine; messages are not persisted (the S3-era StoreMessage keyspace is
// gone), so the record carries exactly the fields a delivery consumes.
type Message struct {
	MessageId              string                       `json:"message_id"`
	TopicArn               string                       `json:"topic_arn"`
	Subject                string                       `json:"subject,omitempty"`
	Message                string                       `json:"message"`
	MessageStructure       string                       `json:"message_structure,omitempty"`
	MessageAttributes      map[string]*MessageAttribute `json:"message_attributes,omitempty"`
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
	Label      string   `json:"label"`
	Principals []string `json:"principals,omitempty"`
	Actions    []string `json:"actions,omitempty"`
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
