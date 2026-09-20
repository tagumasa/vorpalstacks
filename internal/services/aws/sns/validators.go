package sns

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	snsstore "vorpalstacks/internal/store/aws/sns"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// Numeric limits and attribute-value caps live in the store package's
// limits.go — the single definition site for every documented bound. Only
// patterns and vocabularies that already have one definition site live
// here.

// platformApplicationNamePattern is the documented platform application
// name charset: "Application names must be made up of only uppercase and
// lowercase ASCII letters, numbers, underscores, hyphens, and periods"
// (CreatePlatformApplication member documentation).
var platformApplicationNamePattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

// kmsKeyIDRegex validates a bare KMS key ID in UUID hex format
// (8-4-4-4-12 lowercase hex digits, case-insensitive).
var kmsKeyIDRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// messageAttrNamePattern validates message attribute names per AWS docs:
// alphanumeric, underscore, hyphen, and period; 1-256 characters (the
// length bound is the limits file's MaxMessageAttributeNameLength).
var messageAttrNamePattern = regexp.MustCompile(fmt.Sprintf(`^[a-zA-Z0-9_.-]{1,%d}$`, snsstore.MaxMessageAttributeNameLength))

// fifoIdentifierPattern validates MessageGroupId and MessageDeduplicationId
// values, whose documented charset is "up to 128 alphanumeric characters
// (a-z, A-Z, 0-9) and punctuation (!"#$%&'()*+,-./:;<=>?@[\]^_`{|}~)"
// (API_Publish, both members). The punctuation run is the printable ASCII
// set minus space and the alphanumerics, expressed as the four ranges
// !-/ :;-@ [–` {|-~.
var fifoIdentifierPattern = regexp.MustCompile(`^[a-zA-Z0-9!-/:-@\[-` + "`" + `{-~]+$`)

// validMessageAttributeDataTypes is the complete set of DataType values
// accepted by SNS message attributes.
var validMessageAttributeDataTypes = map[string]bool{
	"String":       true,
	"Number":       true,
	"String.Array": true,
	"Binary":       true,
}

// ---------------------------------------------------------------------------
// Protocol validators
// ---------------------------------------------------------------------------

// validateProtocol returns an error when the protocol is not in the
// registry. The accepted-values list in the message is generated from the
// registry itself, so it can never disagree with the acceptance check.
func validateProtocol(protocol string) error {
	if _, ok := protocolRegistry[protocol]; !ok {
		return NewInvalidParameter(fmt.Sprintf(
			"Invalid protocol: %s. Valid values: %s",
			protocol, sortedVocabulary(protocolRegistry)))
	}
	return nil
}

// ---------------------------------------------------------------------------
// Topic validators
// ---------------------------------------------------------------------------

// isAlnumHyphenUnderscore reports whether r is drawn from the
// alphanumeric-plus-hyphen-underscore charset the topic-name member and the
// Lambda function-name endpoint share — the single definition of the class
// both validators enforce.
func isAlnumHyphenUnderscore(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_'
}

// validateTopicName validates the full set of SNS topic name constraints:
// non-empty, max 256 chars, alphanumeric + hyphen + underscore, .fifo
// suffix handled for FIFO topics. The CreateTopic documentation states no
// reserved-prefix rule ("Topic names must be made up of only uppercase and
// lowercase ASCII letters, numbers, underscores, and hyphens, and must be
// between 1 and 256 characters long"), so aws-/amazon-prefixed names are
// legal and accepted.
func validateTopicName(name string) error {
	if name == "" {
		return NewInvalidParameter("Topic name is required")
	}
	if len(name) > snsstore.MaxTopicNameLength {
		return NewInvalidParameter(fmt.Sprintf("Topic name must not exceed %d characters", snsstore.MaxTopicNameLength))
	}

	// Character validation (allow .fifo suffix).
	baseName := name
	if strings.HasSuffix(name, ".fifo") {
		baseName = strings.TrimSuffix(name, ".fifo")
	}
	for _, c := range baseName {
		if !isAlnumHyphenUnderscore(c) {
			return NewInvalidParameter("Topic name can only contain alphanumeric characters, hyphens, and underscores")
		}
	}

	return nil
}

// validateDataProtectionPolicy validates the DataProtectionPolicy parameter
// accepted by PutDataProtectionPolicy and inline by CreateTopic. AWS caps the
// policy at 30,720 bytes and requires a valid JSON document.
func validateDataProtectionPolicy(policy string) error {
	if len(policy) > snsstore.MaxTopicAttributeValueLength {
		return NewInvalidParameter(fmt.Sprintf("DataProtectionPolicy value too long: %d characters (maximum %d)", len(policy), snsstore.MaxTopicAttributeValueLength))
	}
	var policyCheck interface{}
	if err := json.Unmarshal([]byte(policy), &policyCheck); err != nil {
		return NewInvalidParameter(fmt.Sprintf("Invalid DataProtectionPolicy: not valid JSON: %s", err.Error()))
	}
	return nil
}

// validateTopicAttribute validates well-known topic attributes against
// their documented value constraints; unknown attributes pass through
// without validation (forward-compatible with future AWS additions). All
// attribute values are capped at the documented topic-attribute ceiling for
// DoS protection. isFifo gates the FIFO-only attributes
// ("The following attributes apply only to FIFO topics").
func validateTopicAttribute(name, value string, isFifo bool) error {
	// The DataProtectionPolicy attribute key is reserved: the policy is set
	// through the CreateTopic input parameter or PutDataProtectionPolicy
	// only, never through the generic attribute map (SetTopicAttributes or
	// CreateTopic Attributes), matching the documented attribute sets of
	// those APIs.
	if name == snsstore.AttrDataProtectionPolicy {
		return NewInvalidParameter("DataProtectionPolicy cannot be set via topic attributes; use the DataProtectionPolicy parameter of CreateTopic or the PutDataProtectionPolicy API")
	}

	// General DoS cap for all topic attribute values.
	if len(value) > snsstore.MaxTopicAttributeValueLength {
		return NewInvalidParameter(fmt.Sprintf("%s value too long: %d characters (maximum %d)", name, len(value), snsstore.MaxTopicAttributeValueLength))
	}

	switch name {
	case snsstore.AttrDeliveryPolicy:
		// The topic-level document wraps the policies under "http" with
		// default-prefixed names; the parser validates every documented
		// field and constraint.
		_, err := parseTopicDeliveryPolicy(value)
		return err
	case snsstore.AttrPolicy:
		return validateJSONAttribute(name, value)
	case snsstore.AttrDisplayName:
		if len(value) > snsstore.MaxDisplayNameLength {
			return NewInvalidParameter(fmt.Sprintf("DisplayName too long: %d characters (maximum %d)", len(value), snsstore.MaxDisplayNameLength))
		}
	case snsstore.AttrKmsMasterKeyId:
		if value != "" && !strings.HasPrefix(value, "arn:") && !isValidKmsKeyId(value) {
			return NewInvalidParameter(fmt.Sprintf("Invalid KmsMasterKeyId: %s", value))
		}
	case snsstore.AttrSignatureVersion:
		// "By default, SignatureVersion is set to 1"; version 2 switches the
		// signature to SHA256withRSA. No other version exists.
		if value != "1" && value != "2" {
			return NewInvalidParameter(fmt.Sprintf("Invalid SignatureVersion: %s. Valid values: 1, 2", value))
		}
	case snsstore.AttrMaximumMessageSize:
		// Documented range 1024..1048576; the platform honours the
		// attribute up to its flat transport ceiling — values above it are
		// rejected as a recorded platform bound rather than accepted and
		// unenforced (limits.go, MinTopicMessageSize comment).
		size, err := strconv.Atoi(value)
		if err != nil {
			return NewInvalidParameter(fmt.Sprintf("Invalid MaximumMessageSize: %s", value))
		}
		if size < snsstore.MinTopicMessageSize || size > snsstore.MaxMessageSize {
			return NewInvalidParameter(fmt.Sprintf(
				"Invalid MaximumMessageSize: %d (this platform accepts %d to %d)",
				size, snsstore.MinTopicMessageSize, snsstore.MaxMessageSize))
		}
	case snsstore.AttrContentBasedDedup:
		if !isFifo {
			return NewInvalidParameter("ContentBasedDeduplication applies only to FIFO topics")
		}
		if value != "true" && value != "false" {
			return NewInvalidParameter(fmt.Sprintf("Invalid ContentBasedDeduplication: %s. Valid values: true, false", value))
		}
	case snsstore.AttrFifoThroughputScope:
		if !isFifo {
			return NewInvalidParameter("FifoThroughputScope applies only to FIFO topics")
		}
		if value != "Topic" && value != "MessageGroup" {
			return NewInvalidParameter(fmt.Sprintf("Invalid FifoThroughputScope: %s. Valid values: Topic, MessageGroup", value))
		}
	case snsstore.AttrArchivePolicy:
		if !isFifo {
			return NewInvalidParameter("ArchivePolicy applies only to FIFO topics")
		}
		return validateJSONAttribute(name, value)
	}
	return nil
}

// rejectNonSettableTopicAttribute rejects the topic-plane keys
// SetTopicAttributes must refuse: FifoTopic is create-only (the topic type
// is fixed by the name at creation — a .fifo name makes a FIFO topic, any
// other name a standard topic, and no later write may contradict the
// predicate), and the synthesised or derived keys (EffectiveDeliveryPolicy,
// TopicArn, Owner, the Subscriptions* counters) are read-only members of
// GetTopicAttributes, not writable state.
func rejectNonSettableTopicAttribute(name string) error {
	switch name {
	case snsstore.AttrFifoTopic:
		return NewInvalidParameter("FifoTopic cannot be set via SetTopicAttributes; the topic type is fixed at creation by the topic name")
	case snsstore.AttrEffectiveDelivery:
		return NewInvalidParameter("EffectiveDeliveryPolicy is a read-only attribute synthesised from DeliveryPolicy and the system defaults")
	case snsstore.AttrTopicArn, snsstore.AttrOwner,
		snsstore.AttrSubscriptionsConfirmed, snsstore.AttrSubscriptionsDeleted, snsstore.AttrSubscriptionsPending:
		return NewInvalidParameter(fmt.Sprintf("%s is a read-only topic attribute", name))
	}
	return nil
}

// validateJSONAttribute validates that the value is valid JSON.
func validateJSONAttribute(name, value string) error {
	if value == "" {
		return nil
	}
	var v interface{}
	if err := json.Unmarshal([]byte(value), &v); err != nil {
		return NewInvalidParameter(fmt.Sprintf("Invalid %s: not valid JSON: %s", name, err.Error()))
	}
	return nil
}

// ---------------------------------------------------------------------------
// Subscription validators
// ---------------------------------------------------------------------------

// validateSubscriptionAttribute validates well-known subscription attributes
// against their documented value constraints, and rejects the read-only
// keys of the documented GetSubscriptionAttributes set. protocol gates the
// per-protocol attributes (RawMessageDelivery). FilterPolicy is validated
// separately by validateFilterPolicyScopePair — its grammar depends on the
// scope it will match under, which the per-attribute switch cannot see.
// Unknown attribute names pass through (forward-compatible with future
// AWS additions) under the plane-wide value cap; the reserved internal
// keys below are refused.
func validateSubscriptionAttribute(name, value, protocol string) error {
	switch name {
	// AuthenticateOnUnsubscribe is set exclusively through the
	// ConfirmSubscription input parameter; it is not a writable attribute.
	case snsstore.AttrAuthenticateOnUnsubscribe:
		return NewInvalidParameter("AuthenticateOnUnsubscribe cannot be set via SetSubscriptionAttributes; it is set when confirming the subscription")
	case snsstore.AttrFilterPolicyScope:
		return validateFilterPolicyScope(value)
	case snsstore.AttrRedrivePolicy:
		return validateRedrivePolicy(value)
	case snsstore.AttrDeliveryPolicy:
		// The subscription-level document names the three policies without
		// the "http" wrapper; the parser validates every documented field
		// and constraint.
		_, err := parseSubscriptionDeliveryPolicy(value)
		return err
	case snsstore.AttrRawMessageDelivery:
		// "When set to true, enables raw message delivery to Amazon SQS or
		// HTTP/S endpoints" (SetSubscriptionAttributes, Subscribe) — the
		// value is a boolean literal and the attribute exists only on the
		// protocols that sentence names.
		if value != "true" && value != "false" {
			return NewInvalidParameter(fmt.Sprintf("Invalid RawMessageDelivery: %s. Valid values: true, false", value))
		}
		if !protocolSupportsRawDelivery(protocol) {
			return NewInvalidParameter(fmt.Sprintf(
				"RawMessageDelivery applies only to Amazon SQS and HTTP/S endpoints, not protocol %s", protocol))
		}
	case snsstore.AttrPendingConfirmation, snsstore.AttrConfirmationWasAuthenticated,
		snsstore.AttrSubscriptionArn, snsstore.AttrTopicArn, snsstore.AttrOwner,
		snsstore.AttrProtocol, snsstore.AttrEndpoint, snsstore.AttrReplayStatus:
		return NewInvalidParameter(fmt.Sprintf("%s is a read-only subscription attribute", name))
	}
	// No AWS-documented value bound exists for subscription attributes;
	// the platform caps every value — known and unknown names alike — at
	// the topic plane's bound so the two attribute planes carry the same
	// storage-exposure ceiling.
	if len(value) > snsstore.MaxSubscriptionAttributeValueLength {
		return NewInvalidParameter(fmt.Sprintf("%s value too long: %d characters (maximum %d)", name, len(value), snsstore.MaxSubscriptionAttributeValueLength))
	}
	return nil
}

// validateFilterPolicy validates the filter policy document against the
// scope it will match under: its size, leaf-key count, total combination,
// wildcard complexity and every operator operand against the documented
// constraints, plus the scope's nesting rule — a nested policy is
// payload-based filtering alone. The grammar and its matcher share
// filter_policy.go, so a policy is accepted only when delivery-time
// matching can evaluate it.
func validateFilterPolicy(value, scope string) error {
	return validateFilterPolicyDocument(value, scope)
}

// validateFilterPolicyScopePair validates the FilterPolicy and
// FilterPolicyScope a subscription will carry together. The two attributes
// are coupled: "Amazon SNS accepts a nested filter policy for
// payload-based filtering" while attribute-based filtering "doesn't accept
// a nested filter policy", so the policy document is validated against the
// scope that will apply once the request lands — an absent scope is the
// documented default, MessageAttributes. Both cores call this with the
// request's paired values resolved against the stored subscription, which
// covers both directions: a FilterPolicy write under the effective scope,
// and a FilterPolicyScope write re-validated against the stored policy.
func validateFilterPolicyScopePair(policy, scope string) error {
	if scope == "" {
		scope = snsstore.FilterPolicyScopeAttributes
	}
	return validateFilterPolicy(policy, scope)
}

// validateFilterPolicyScope validates the FilterPolicyScope attribute.
// AWS accepts only "MessageAttributes" (default) or "MessageBody".
func validateFilterPolicyScope(value string) error {
	switch value {
	case snsstore.FilterPolicyScopeAttributes, snsstore.FilterPolicyScopeBody:
		return nil
	default:
		return NewInvalidParameter(
			fmt.Sprintf("Invalid FilterPolicyScope: %s. Valid values: MessageAttributes, MessageBody", value))
	}
}

// validateRedrivePolicy validates the JSON structure of a subscription
// redrive policy, parsing into the store's single RedrivePolicy type. AWS
// requires deadLetterTargetArn.
func validateRedrivePolicy(value string) error {
	if value == "" {
		return nil
	}

	var rp snsstore.RedrivePolicy
	if err := json.Unmarshal([]byte(value), &rp); err != nil {
		return NewInvalidParameter(fmt.Sprintf("Invalid redrive policy: %s", err.Error()))
	}

	if strings.TrimSpace(rp.DeadLetterTargetArn) == "" {
		return NewInvalidParameter("Invalid redrive policy: deadLetterTargetArn is required")
	}

	return nil
}

// ---------------------------------------------------------------------------
// Publish validators
// ---------------------------------------------------------------------------

// validatePublishParams validates all top-level Publish parameters that have
// AWS-documented constraints. The FIFO state is passed as primitives to keep
// this validator free of store-layer dependencies; maxMessageBytes is the
// topic's effective size ceiling — the MaximumMessageSize attribute when the
// topic sets one, the platform's flat cap otherwise — and the ceiling counts
// "the combined size of the message body and message attributes against this
// value" (Publish), so the parsed attributes ride along for the measure.
func validatePublishParams(isFifo, isContentBasedDedup bool, maxMessageBytes int, message, subject, messageStructure, messageGroupId, messageDeduplicationId string, attrs map[string]*snsstore.MessageAttribute) error {
	if size := len(message) + attributeBytes(attrs); size > maxMessageBytes {
		return NewInvalidParameter(fmt.Sprintf("Message and message attributes too long: %d bytes (maximum %d)", size, maxMessageBytes))
	}

	// Subject is documented as "UTF-8 text … less than 100 characters
	// long", so the ceiling counts Unicode characters.
	if n := utf8.RuneCountInString(subject); n > snsstore.MaxSubjectLength {
		return NewInvalidParameter(fmt.Sprintf("Subject too long: %d characters (maximum %d)", n, snsstore.MaxSubjectLength))
	}
	// "Subjects must be UTF-8 text with no line breaks or control
	// characters" (Publish) — line breaks are themselves control
	// characters, so one sweep over Unicode's control category covers
	// both.
	for _, r := range subject {
		if unicode.IsControl(r) {
			return NewInvalidParameter(fmt.Sprintf("Invalid Subject: control character U+%04X is not allowed (Subjects must be UTF-8 text with no line breaks or control characters)", r))
		}
	}

	if messageStructure != "" && messageStructure != "json" {
		return NewInvalidParameter(fmt.Sprintf("Invalid MessageStructure: %s. Valid value: json", messageStructure))
	}

	if messageStructure == "json" {
		var msgMap map[string]interface{}
		if err := json.Unmarshal([]byte(message), &msgMap); err != nil {
			return NewInvalidParameter(fmt.Sprintf("Invalid parameter: MessageStructure is json but message body is not valid JSON: %s", err.Error()))
		}
		if _, ok := msgMap["default"]; !ok {
			return NewInvalidParameter("Invalid parameter: MessageStructure is json but message body does not contain a 'default' key")
		}
	}

	// Both FIFO identifiers carry the documented charset and length
	// constraints on every topic type — for standard topics the Publish
	// documentation states the same validation rules apply to a forwarded
	// MessageGroupId.
	if messageGroupId != "" {
		if err := validateFifoIdentifier("MessageGroupId", messageGroupId); err != nil {
			return err
		}
	}
	if messageDeduplicationId != "" {
		if err := validateFifoIdentifier("MessageDeduplicationId", messageDeduplicationId); err != nil {
			return err
		}
	}

	if isFifo {
		if messageGroupId == "" {
			return NewInvalidParameter("MessageGroupId is required for FIFO topics")
		}
		// When ContentBasedDeduplication is disabled the publisher must supply
		// a MessageDeduplicationId. When enabled, the caller auto-generates
		// one from the message body hash.
		if !isContentBasedDedup && messageDeduplicationId == "" {
			return NewInvalidParameter("MessageDeduplicationId is required when ContentBasedDeduplication is false")
		}
	} else {
		// MessageGroupId is optional on standard topics: it "is forwarded
		// only to Amazon SQS standard subscriptions to activate fair
		// queues" (API_Publish) and reaches no other endpoint type.
		// MessageDeduplicationId "applies only to FIFO topics".
		if messageDeduplicationId != "" {
			return NewInvalidParameter("MessageDeduplicationId is only valid for FIFO topics")
		}
	}

	return nil
}

// validateFifoIdentifier enforces the documented MessageGroupId and
// MessageDeduplicationId constraints — at most 128 characters drawn from
// the alphanumeric-plus-punctuation set (the charset is ASCII, so byte
// length is character length).
func validateFifoIdentifier(name, value string) error {
	if len(value) > snsstore.MaxFifoIdentifierLength {
		return NewInvalidParameter(fmt.Sprintf("%s too long: %d characters (maximum %d)", name, len(value), snsstore.MaxFifoIdentifierLength))
	}
	if !fifoIdentifierPattern.MatchString(value) {
		return NewInvalidParameter(fmt.Sprintf("Invalid %s %q: must contain up to 128 alphanumeric characters (a-z, A-Z, 0-9) and punctuation", name, value))
	}
	return nil
}

// validateMessageAttributeName validates the name of a message attribute per
// AWS docs: alphanumeric, underscore, hyphen, and period; 1-256 characters.
func validateMessageAttributeName(name string) error {
	if !messageAttrNamePattern.MatchString(name) {
		return NewInvalidParameter(fmt.Sprintf("Invalid message attribute name %q: must match [a-zA-Z0-9_.-] and be 1-256 characters", name))
	}
	return nil
}

// validateMessageAttributeLimits enforces the AWS-documented limits on message
// attributes: maximum 10 attributes, String values up to 256 chars, Binary
// values up to 256 bytes.
func validateMessageAttributeLimits(name, stringValue string, binaryValue []byte) error {
	if len(stringValue) > snsstore.MaxMessageAttributeStringValue {
		return NewInvalidParameter(fmt.Sprintf("Message attribute %q StringValue too long: %d characters (maximum %d)", name, len(stringValue), snsstore.MaxMessageAttributeStringValue))
	}
	if len(binaryValue) > snsstore.MaxMessageAttributeBinaryValue {
		return NewInvalidParameter(fmt.Sprintf("Message attribute %q BinaryValue too long: %d bytes (maximum %d)", name, len(binaryValue), snsstore.MaxMessageAttributeBinaryValue))
	}
	return nil
}

// ---------------------------------------------------------------------------
// Platform validators
// ---------------------------------------------------------------------------

// validatePlatformApplicationName enforces the non-empty requirement, the
// documented charset, and the 256-character ceiling on
// CreatePlatformApplication names, counted in Unicode characters.
func validatePlatformApplicationName(name string) error {
	if name == "" {
		return NewInvalidParameter("Name is required")
	}
	if n := utf8.RuneCountInString(name); n > snsstore.MaxPlatformApplicationNameLength {
		return NewInvalidParameter(fmt.Sprintf("Name too long: %d characters (maximum %d)", n, snsstore.MaxPlatformApplicationNameLength))
	}
	if !platformApplicationNamePattern.MatchString(name) {
		return NewInvalidParameter("Invalid parameter: Name must contain only letters, numbers, underscores, hyphens, and periods")
	}
	return nil
}

// validatePlatformApplicationArn validates the structure of a platform
// application ARN: the "arn:" prefix, the "sns" service namespace, and the
// platform-application resource form app/<platform>/<name> — the form the
// store's endpoint-ARN construction extracts, checked here so the store
// side carries no ARN-format validation of its own.
func validatePlatformApplicationArn(arn string) error {
	if arn == "" {
		return NewInvalidParameter("PlatformApplicationArn is required")
	}
	if !strings.HasPrefix(arn, "arn:") {
		return NewInvalidParameter(fmt.Sprintf("Invalid PlatformApplicationArn: %s", arn))
	}
	_, service, _, _, _ := svcarn.SplitARN(arn)
	if service == "" {
		return NewInvalidParameter(fmt.Sprintf("Invalid PlatformApplicationArn format: %s", arn))
	}
	if service != "sns" {
		return NewInvalidParameter(fmt.Sprintf("PlatformApplicationArn must be an SNS ARN: %s", arn))
	}
	if platform, _ := svcarn.ExtractPlatformApplicationFromARN(arn); platform == "" {
		return NewInvalidParameter(fmt.Sprintf(
			"Invalid PlatformApplicationArn: the resource must name a platform application as app/<platform>/<name>: %s", arn))
	}
	return nil
}

// validatePlatformAttributeValue enforces a length cap on platform application
// and endpoint attribute values for DoS protection.
func validatePlatformAttributeValue(name, value string) error {
	if len(value) > snsstore.MaxPlatformAttributeValueLength {
		return NewInvalidParameter(fmt.Sprintf("%s value too long: %d characters (maximum %d)", name, len(value), snsstore.MaxPlatformAttributeValueLength))
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helper validators
// ---------------------------------------------------------------------------

// isValidKmsKeyId checks whether a value is a valid bare KMS key identifier:
// a UUID-format key ID or an alias name (prefixed with "alias/").
// Key ARNs and alias ARNs are handled by the caller via the "arn:" prefix
// check.
func isValidKmsKeyId(value string) bool {
	if strings.HasPrefix(value, "alias/") {
		return len(value) > len("alias/")
	}
	return kmsKeyIDRegex.MatchString(value)
}
