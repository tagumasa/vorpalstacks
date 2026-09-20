package sns

// Package sns provides SNS (Simple Notification Service) data store implementations
// for vorpalstacks.

// Attribute-map keys for the topic, subscription and platform-endpoint
// attribute planes. The string maps are both the persistence format and
// the wire format (SetTopicAttributes / SetSubscriptionAttributes /
// SetEndpointAttributes inputs and the Get*Attributes responses), so every
// reference to a known key goes through these definitions.
//
// Writability classification (decided against the AWS documentation):
// topic keys settable at any time are Policy,
// DeliveryPolicy, DisplayName, KmsMasterKeyId, SignatureVersion (1|2),
// TracingConfig, MaximumMessageSize and the FIFO-only
// ContentBasedDeduplication, FifoThroughputScope and ArchivePolicy; FifoTopic
// is create-only (the topic type is fixed by the name at creation); the
// synthesised and derived keys (EffectiveDeliveryPolicy, Owner, TopicArn,
// SubscriptionsConfirmed/Deleted/Pending) are read-only. Subscription keys
// settable at any time are DeliveryPolicy, FilterPolicy, FilterPolicyScope,
// RawMessageDelivery (true|false, sqs/http/https only), RedrivePolicy and
// SubscriptionRoleArn; the lifecycle keys (PendingConfirmation,
// ConfirmationWasAuthenticated) and the identity keys (SubscriptionArn,
// TopicArn, Owner, Protocol, Endpoint) are read-only, as is ReplayStatus.
// ArchivePolicy, KmsMasterKeyId and TracingConfig describe features the
// platform does not run (message archiving, SSE, X-Ray tracing) and pass
// through with that exclusion recorded in docs/services.md. A key not
// listed here is unknown to the implementation and passes through
// unvalidated by design (forward-compatible with future AWS additions).
// The definitions below cover every key code references; the classified
// pass-through keys no code path touches (TracingConfig,
// SubscriptionRoleArn) stay named in this classification comment only — a
// definition nothing consumes is dead vocabulary, not documentation.
const (
	// Topic plane.
	AttrPolicy              = "Policy"
	AttrDeliveryPolicy      = "DeliveryPolicy"
	AttrDisplayName         = "DisplayName"
	AttrKmsMasterKeyId      = "KmsMasterKeyId"
	AttrSignatureVersion    = "SignatureVersion"
	AttrMaximumMessageSize  = "MaximumMessageSize"
	AttrFifoTopic           = "FifoTopic"
	AttrContentBasedDedup   = "ContentBasedDeduplication"
	AttrFifoThroughputScope = "FifoThroughputScope"
	AttrArchivePolicy       = "ArchivePolicy"
	AttrEffectiveDelivery   = "EffectiveDeliveryPolicy"

	// Subscription plane.
	AttrRawMessageDelivery           = "RawMessageDelivery"
	AttrFilterPolicy                 = "FilterPolicy"
	AttrFilterPolicyScope            = "FilterPolicyScope"
	AttrRedrivePolicy                = "RedrivePolicy"
	AttrPendingConfirmation          = "PendingConfirmation"
	AttrConfirmationWasAuthenticated = "ConfirmationWasAuthenticated"

	// Platform endpoint plane.
	AttrEnabled        = "Enabled"
	AttrToken          = "Token"
	AttrCustomUserData = "CustomUserData"

	// Internal keys: stored in an attribute map but excluded from the API
	// attribute surfaces — set through their owning operations only, never
	// through the generic attribute paths.
	AttrDataProtectionPolicy      = "DataProtectionPolicy"
	AttrAuthenticateOnUnsubscribe = "AuthenticateOnUnsubscribe"
)

// Synthesised and identity keys: members the Get*Attributes responses
// synthesise from the record (or that no write path may set), each defined
// once because both the read planes emit them and the write planes reject
// them. TopicArn and Owner serve both the topic and subscription planes —
// one key string, one definition.
const (
	// Topic plane: derived counter members and shared identity members.
	AttrSubscriptionsConfirmed = "SubscriptionsConfirmed"
	AttrSubscriptionsDeleted   = "SubscriptionsDeleted"
	AttrSubscriptionsPending   = "SubscriptionsPending"
	AttrTopicArn               = "TopicArn"
	AttrOwner                  = "Owner"

	// Subscription plane: the record's identity members and the FIFO
	// replay status, read-only on the attribute plane.
	AttrSubscriptionArn = "SubscriptionArn"
	AttrProtocol        = "Protocol"
	AttrEndpoint        = "Endpoint"
	AttrReplayStatus    = "ReplayStatus"
)

// Filter-policy scope values — the FilterPolicyScope attribute's
// documented value set. Attribute-based matching is the documented default
// when the attribute is unset; payload-based matching alone may carry a
// nested policy. Both the store's delivery-time resolution and the
// service's set-time validation read these definitions, so the default
// cannot drift between the two layers.
const (
	FilterPolicyScopeAttributes = "MessageAttributes"
	FilterPolicyScopeBody       = "MessageBody"
)

// mergeAttrs merges src into dst's attribute map, initialising dst when nil,
// and returns the merged map — the single merge every attribute-set path
// (topic, subscription, platform application, platform endpoint and the
// endpoint merge-on-create) applies.
func mergeAttrs(dst, src map[string]string) map[string]string {
	if dst == nil {
		dst = make(map[string]string, len(src))
	}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
