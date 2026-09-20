package sns

// Package sns provides SNS (Simple Notification Service) data store implementations
// for vorpalstacks.

// Documented bounds of the SNS API surface. Every bound the service layer
// enforces is defined here and referenced at each check site; the raw
// number appears only at this definition site, so declaration drift is
// impossible. Values are docs-derived where the Smithy model declares no
// trait (the model carries @length only on ARN/TagKey/TagValue members);
// the source of each value is stated in its comment.
const (
	// MaxMessageSize: "With the exception of SMS, messages must be UTF-8
	// encoded strings and at most 256 KB in size (262,144 bytes)"
	// (API_Publish, Message constraints).
	MaxMessageSize = 256 * 1024

	// MaxSubjectLength: subjects must be "less than 100 characters long"
	// (API_Publish, Subject constraints), so the longest legal subject is
	// 99 characters.
	MaxSubjectLength = 99

	// MaxBatchTotalSize: the total PublishBatch request payload is capped
	// at 256 KB, counted across every entry's message, subject and
	// message attributes.
	MaxBatchTotalSize = 256 * 1024

	// MaxBatchEntries: PublishBatch accepts at most 10 entries per request
	// (TooManyEntriesInBatchRequest beyond).
	MaxBatchEntries = 10

	// MaxTopicNameLength: topic names "must be between 1 and 256
	// characters long" (API_CreateTopic, Name constraints).
	MaxTopicNameLength = 256

	// MaxDisplayNameLength: the DisplayName topic attribute is capped at
	// 100 characters (SetTopicAttributes documentation).
	MaxDisplayNameLength = 100

	// MaxFilterPolicyKeys: "A filter policy can have a maximum of five
	// keys" (filter policy constraints, SNS developer guide).
	MaxFilterPolicyKeys = 5

	// MaxFilterPolicyCombination: "The total combination of values in a
	// filter policy must not exceed 150. To calculate the total
	// combination, multiply the number of values in each array" (filter
	// policy constraints, SNS developer guide).
	MaxFilterPolicyCombination = 150

	// MaxFilterPolicySizeBytes: "The maximum size of a filter policy is
	// 256 KB" (filter policy constraints, SNS developer guide).
	MaxFilterPolicySizeBytes = 256 * 1024

	// MaxFilterPolicyNumericValue: for numeric matching "the value can
	// range from -10^9 to 10^9 (-1 billion to 1 billion), with five digits
	// of accuracy after the decimal point" (filter policy constraints,
	// SNS developer guide). Five digits of accuracy means a value that
	// needs more than five fractional digits is out of range.
	MaxFilterPolicyNumericValue = 1000000000

	// MaxFilterPolicyNumericAccuracyDigits: the same sentence's accuracy
	// half — "five digits of accuracy after the decimal point" — as the
	// digit count the numeric validator scales by.
	MaxFilterPolicyNumericAccuracyDigits = 5

	// MaxWildcardFiltersPerPattern and MaxWildcardPatternComplexity: the
	// wildcard usage guidelines — "Maximum 3 wildcards per pattern" and
	// "Complexity level of 100 or less per wildcard pattern" (string value
	// matching, SNS developer guide).
	MaxWildcardFiltersPerPattern = 3
	MaxWildcardPatternComplexity = 100

	// MaxStandardTopicsPerAccount / MaxFifoTopicsPerAccount: "Topics —
	// Standard: 100,000 per account; FIFO: 1,000 per account" (SNS
	// endpoints and quotas, AWS general reference). Enforced per Region
	// store — AWS account resources live in a Region.
	MaxStandardTopicsPerAccount = 100000
	MaxFifoTopicsPerAccount     = 1000

	// MaxStandardSubscriptionsPerTopic / MaxFifoSubscriptionsPerTopic:
	// "Subscriptions — Standard: 12,500,000 per topic; FIFO: 100 per topic"
	// (SNS endpoints and quotas, AWS general reference).
	MaxStandardSubscriptionsPerTopic = 12500000
	MaxFifoSubscriptionsPerTopic     = 100

	// MaxFilterPoliciesPerTopic: "Subscription filter policies — 200 filter
	// policies per topic" (SNS endpoints and quotas, AWS general
	// reference). The account-scope companion (10,000 per account) is not
	// enforced: it aggregates across Region stores and the platform keeps
	// no cross-Region account view.
	MaxFilterPoliciesPerTopic = 200

	// MaxFirehoseSubscriptionsPerTopic: "For Firehose delivery streams,
	// 5 per topic, per subscription owner" (SNS endpoints and quotas, AWS
	// general reference); the platform is single-account, so the
	// per-topic count is the per-owner count.
	MaxFirehoseSubscriptionsPerTopic = 5

	// MaxResourceArnLength: the tag operations' ResourceArn member targets
	// the model's AmazonResourceName shape, "@length": {"min": 1,
	// "max": 1011}.
	MaxResourceArnLength = 1011

	// MaxTopicAttributeValueLength: topic attribute values documented at
	// 30,720 bytes (Policy, DeliveryPolicy, DataProtectionPolicy —
	// CreateTopic/SetTopicAttributes documentation); applied to every
	// topic attribute value as a DoS cap.
	MaxTopicAttributeValueLength = 30720

	// MaxPlatformAttributeValueLength: platform application and endpoint
	// attribute-value DoS cap (platform bound; AWS documents no single
	// value for this family).
	MaxPlatformAttributeValueLength = 8192

	// MaxSubscriptionAttributeValueLength: subscription attribute-value DoS
	// cap (platform bound; AWS documents no value bound for this family) —
	// the topic plane's documented bound, so the two attribute planes
	// carry the same storage-exposure ceiling.
	MaxSubscriptionAttributeValueLength = MaxTopicAttributeValueLength

	// MaxPlatformApplicationNameLength: platform application names "must
	// be between 1 and 256 characters long"
	// (CreatePlatformApplication, Name constraints).
	MaxPlatformApplicationNameLength = 256

	// MaxPlatformTokenLength: platform endpoint tokens are capped at 2048
	// bytes (CreatePlatformEndpoint, Token constraints).
	MaxPlatformTokenLength = 2048

	// MaxCustomUserDataLength: platform endpoint CustomUserData is capped
	// at 2048 bytes (CreatePlatformEndpoint documentation).
	MaxCustomUserDataLength = 2048

	// MaxEndpointURLLength: http/https subscription endpoint URL DoS cap
	// (platform bound; AWS documents no explicit URL length).
	MaxEndpointURLLength = 2048

	// MaxFifoIdentifierLength: MessageGroupId and MessageDeduplicationId
	// "can contain up to 128 alphanumeric characters (a-z, A-Z, 0-9) and
	// punctuation" (API_Publish, both members).
	MaxFifoIdentifierLength = 128

	// MinTopicMessageSize: the MaximumMessageSize topic attribute accepts
	// "valid values ... 1024 to 1048576 (1 MiB)" (SetTopicAttributes). The
	// platform honours the attribute up to its flat transport ceiling
	// (MaxMessageSize): the platform's SQS and Lambda invocations cap every
	// payload at 256 KiB, so values above that ceiling are rejected as a
	// recorded platform bound instead of being accepted and unenforced.
	MinTopicMessageSize = 1024

	// MinDeliveryPolicyDelaySeconds / MaxDeliveryPolicySeconds: the delivery
	// policy's delay bounds — "minDelayTarget ... 1 to maximum delay" and
	// "maxDelayTarget ... Minimum delay to 3,600", with "the total policy
	// retry time for an HTTP/S endpoint cannot be greater than 3,600 seconds"
	// as the same documented ceiling (message delivery retries documentation).
	MinDeliveryPolicyDelaySeconds = 1
	MaxDeliveryPolicySeconds      = 3600

	// MaxDeliveryPolicyRetries: "numRetries ... 0 to 100" (message delivery
	// retries documentation).
	MaxDeliveryPolicyRetries = 100

	// MinThrottleReceivesPerSecond: "maxReceivesPerSecond ... 1 or greater"
	// (message delivery retries documentation).
	MinThrottleReceivesPerSecond = 1

	// DefaultDeliveryPolicyMinDelaySeconds, DefaultDeliveryPolicyMaxDelaySeconds
	// and DefaultDeliveryPolicyNumRetries: the system defaults of the HTTP/S
	// delivery policy when the policy document omits the fields (constraint
	// table, message delivery retries documentation) — three retries, twenty
	// seconds apart, with the linear backoff the service layer defines.
	DefaultDeliveryPolicyMinDelaySeconds = 20
	DefaultDeliveryPolicyMaxDelaySeconds = 20
	DefaultDeliveryPolicyNumRetries      = 3
)

// Message-attribute bounds (Publish MessageAttributes documentation):
// at most 10 attributes per message, String values up to 256 characters,
// Binary values up to 256 bytes, names matching [a-zA-Z0-9_.-]{1,256}.
const (
	MaxMessageAttributes           = 10
	MaxMessageAttributeStringValue = 256
	MaxMessageAttributeBinaryValue = 256
	MaxMessageAttributeNameLength  = 256
)
