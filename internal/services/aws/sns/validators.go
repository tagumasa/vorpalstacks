package sns

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/common/errors"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// AWS-documented constants
// ---------------------------------------------------------------------------

// SNS message and batch size limits (AWS documented constants). The
// Publish member documentation states subjects must be "less than 100
// characters long", so the longest legal subject is 99 characters.
const (
	maxMessageSize    = 256 * 1024 // 256 KB
	maxSubjectLength  = 99
	maxBatchTotalSize = 256 * 1024 // 256 KB total for all entries
	maxBatchEntries   = 10
)

// AWS-documented attribute value caps (DoS protection).
// Topic Policy and DeliveryPolicy are documented at 30,720 bytes.
// Platform application and endpoint attributes use a generous cap.
const (
	maxTopicAttributeValueLength    = 30720
	maxPlatformAttributeValueLength = 8192
)

// SNS message attribute limits per AWS docs.
const (
	maxMessageAttributes      = 10
	maxMessageAttrStringValue = 256 // chars
	maxMessageAttrBinaryValue = 256 // bytes
)

// maxPlatformApplicationNameLength is the documented platform application
// name ceiling: names "must be between 1 and 256 characters long"
// (CreatePlatformApplication member documentation).
const maxPlatformApplicationNameLength = 256

// platformApplicationNamePattern is the documented platform application
// name charset: "Application names must be made up of only uppercase and
// lowercase ASCII letters, numbers, underscores, hyphens, and periods"
// (CreatePlatformApplication member documentation).
var platformApplicationNamePattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

// validatePlatformApplicationName enforces the non-empty requirement, the
// documented charset, and the 256-character ceiling on
// CreatePlatformApplication names, counted in Unicode characters.
func validatePlatformApplicationName(name string) error {
	if name == "" {
		return NewInvalidParameter("Name is required")
	}
	if n := utf8.RuneCountInString(name); n > maxPlatformApplicationNameLength {
		return NewInvalidParameter(fmt.Sprintf("Name too long: %d characters (maximum %d)", n, maxPlatformApplicationNameLength))
	}
	if !platformApplicationNamePattern.MatchString(name) {
		return NewInvalidParameter("Invalid parameter: Name must contain only letters, numbers, underscores, hyphens, and periods")
	}
	return nil
}

// Endpoint URL length cap for http/https protocols.
const maxEndpointURLLength = 2048

// ---------------------------------------------------------------------------
// Regex patterns
// ---------------------------------------------------------------------------

// kmsKeyIDRegex validates a bare KMS key ID in UUID hex format
// (8-4-4-4-12 lowercase hex digits, case-insensitive).
var kmsKeyIDRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// messageAttrNamePattern validates message attribute names per AWS docs:
// alphanumeric, underscore, hyphen, and period; 1-256 characters.
var messageAttrNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]{1,256}$`)

// ---------------------------------------------------------------------------
// Enum maps
// ---------------------------------------------------------------------------

// validProtocols lists the nine AWS-supported subscription protocols.
// Any protocol not in this map is rejected at Subscribe time.
var validProtocols = map[string]bool{
	"http":        true,
	"https":       true,
	"email":       true,
	"email-json":  true,
	"sms":         true,
	"sqs":         true,
	"application": true,
	"lambda":      true,
	"firehose":    true,
}

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

// validateProtocol returns an error when the protocol is not one of the nine
// AWS-supported values.
func validateProtocol(protocol string) error {
	if !validProtocols[protocol] {
		return NewInvalidParameter(fmt.Sprintf(
			"Invalid protocol: %s. Valid values: http, https, email, email-json, sms, sqs, application, lambda, firehose",
			protocol))
	}
	return nil
}

// validateEndpointForProtocol validates the endpoint format against the
// protocol-specific requirements. This catches grossly invalid endpoints at
// Subscribe time rather than silently failing at delivery time.
func validateEndpointForProtocol(protocol, endpoint string) error {
	switch protocol {
	case "http":
		if !strings.HasPrefix(endpoint, "http://") {
			return NewInvalidParameter("Endpoint must be a valid HTTP URL starting with http://")
		}
		if len(endpoint) > maxEndpointURLLength {
			return NewInvalidParameter(fmt.Sprintf("Endpoint URL too long: %d characters (maximum %d)", len(endpoint), maxEndpointURLLength))
		}
		if _, err := url.Parse(endpoint); err != nil {
			return NewInvalidParameter(fmt.Sprintf("Invalid endpoint URL: %s", err.Error()))
		}

	case "https":
		if !strings.HasPrefix(endpoint, "https://") {
			return NewInvalidParameter("Endpoint must be a valid HTTPS URL starting with https://")
		}
		if len(endpoint) > maxEndpointURLLength {
			return NewInvalidParameter(fmt.Sprintf("Endpoint URL too long: %d characters (maximum %d)", len(endpoint), maxEndpointURLLength))
		}
		if _, err := url.Parse(endpoint); err != nil {
			return NewInvalidParameter(fmt.Sprintf("Invalid endpoint URL: %s", err.Error()))
		}

	case "sqs":
		if !strings.HasPrefix(endpoint, "http") && !strings.HasPrefix(endpoint, "arn:") {
			return NewInvalidParameter("Endpoint must be a valid SQS queue URL or ARN for protocol sqs")
		}

	case "lambda":
		if !strings.HasPrefix(endpoint, "arn:") {
			if endpoint == "" {
				return NewInvalidParameter("Endpoint must be a valid Lambda function ARN or name")
			}
			for _, c := range endpoint {
				if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
					return NewInvalidParameter("Endpoint must be a valid Lambda function ARN or name")
				}
			}
		}

	case "email", "email-json":
		if !strings.Contains(endpoint, "@") {
			return NewInvalidParameter("Endpoint must be a valid email address for protocol " + protocol)
		}

	case "application":
		if !strings.HasPrefix(endpoint, "arn:") {
			return NewInvalidParameter("Endpoint must be a valid platform endpoint ARN for protocol application")
		}

	case "firehose":
		if !strings.HasPrefix(endpoint, "arn:") {
			return NewInvalidParameter("Endpoint must be a valid Firehose delivery stream ARN for protocol firehose")
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Topic validators
// ---------------------------------------------------------------------------

// validateTopicName validates the full set of SNS topic name constraints:
// non-empty, max 256 chars, alphanumeric + hyphen + underscore, no reserved
// AWS prefixes, .fifo suffix handled for FIFO topics.
func validateTopicName(name string) error {
	if name == "" {
		return NewInvalidParameter("Topic name is required")
	}
	if len(name) > 256 {
		return NewInvalidParameter("Topic name must not exceed 256 characters")
	}

	// Reject AWS-reserved prefixes (case-insensitive).
	lower := strings.ToLower(name)
	if strings.HasPrefix(lower, "aws") || strings.HasPrefix(lower, "amazon") {
		return NewInvalidParameter(fmt.Sprintf("Topic name %q starts with a reserved prefix (aws/amazon)", name))
	}

	// Character validation (allow .fifo suffix).
	baseName := name
	if strings.HasSuffix(name, ".fifo") {
		baseName = strings.TrimSuffix(name, ".fifo")
	}
	for _, c := range baseName {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return NewInvalidParameter("Topic name can only contain alphanumeric characters, hyphens, and underscores")
		}
	}

	return nil
}

// validateDataProtectionPolicy validates the DataProtectionPolicy parameter
// accepted by PutDataProtectionPolicy and inline by CreateTopic. AWS caps the
// policy at 30,720 bytes and requires a valid JSON document.
func validateDataProtectionPolicy(policy string) error {
	if len(policy) > maxTopicAttributeValueLength {
		return NewInvalidParameter(fmt.Sprintf("DataProtectionPolicy value too long: %d characters (maximum %d)", len(policy), maxTopicAttributeValueLength))
	}
	var policyCheck interface{}
	if err := json.Unmarshal([]byte(policy), &policyCheck); err != nil {
		return NewInvalidParameter(fmt.Sprintf("Invalid DataProtectionPolicy: not valid JSON: %s", err.Error()))
	}
	return nil
}

// validateTopicAttribute validates well-known topic attributes that have
// structured values. Unknown attributes pass through without validation
// (forward-compatible with future AWS additions). All attribute values are
// capped at maxTopicAttributeValueLength for DoS protection.
func validateTopicAttribute(name, value string) error {
	// The DataProtectionPolicy attribute key is reserved: the policy is set
	// through the CreateTopic input parameter or PutDataProtectionPolicy
	// only, never through the generic attribute map (SetTopicAttributes or
	// CreateTopic Attributes), matching the documented attribute sets of
	// those APIs.
	if name == "DataProtectionPolicy" {
		return NewInvalidParameter("DataProtectionPolicy cannot be set via topic attributes; use the DataProtectionPolicy parameter of CreateTopic or the PutDataProtectionPolicy API")
	}

	// General DoS cap for all topic attribute values.
	if len(value) > maxTopicAttributeValueLength {
		return NewInvalidParameter(fmt.Sprintf("%s value too long: %d characters (maximum %d)", name, len(value), maxTopicAttributeValueLength))
	}

	switch name {
	case "DeliveryPolicy":
		return validateJSONAttribute(name, value)
	case "Policy":
		return validateJSONAttribute(name, value)
	case "DisplayName":
		if len(value) > 100 {
			return NewInvalidParameter(fmt.Sprintf("DisplayName too long: %d characters (maximum 100)", len(value)))
		}
	case "KmsMasterKeyId":
		if value != "" && !strings.HasPrefix(value, "arn:") && !isValidKmsKeyId(value) {
			return NewInvalidParameter(fmt.Sprintf("Invalid KmsMasterKeyId: %s", value))
		}
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
// that have structured values. Unknown attributes pass through without
// validation (forward-compatible with future AWS additions), except for the
// reserved internal keys below.
func validateSubscriptionAttribute(name, value string) error {
	switch name {
	// AuthenticateOnUnsubscribe is set exclusively through the
	// ConfirmSubscription input parameter; it is not a writable attribute.
	case "AuthenticateOnUnsubscribe":
		return NewInvalidParameter("AuthenticateOnUnsubscribe cannot be set via SetSubscriptionAttributes; it is set when confirming the subscription")
	case "FilterPolicy":
		return validateFilterPolicy(value)
	case "FilterPolicyScope":
		return validateFilterPolicyScope(value)
	case "RedrivePolicy":
		return validateRedrivePolicy(value)
	}
	return nil
}

// validateFilterPolicy validates the JSON structure of a subscription filter
// policy. AWS allows at most 100 attribute names in a single policy.
func validateFilterPolicy(value string) error {
	if value == "" {
		return nil
	}

	var policy map[string]interface{}
	if err := json.Unmarshal([]byte(value), &policy); err != nil {
		return NewInvalidParameter(fmt.Sprintf("Invalid filter policy: %s", err.Error()))
	}

	if len(policy) > 100 {
		return ErrFilterLimitExceeded
	}

	for attrName, raw := range policy {
		if attrName == "" {
			return NewInvalidParameter("Invalid filter policy: attribute name must not be empty")
		}
		if err := validateFilterPolicyValue(raw); err != nil {
			return err
		}
	}

	return nil
}

// validateFilterPolicyValue validates the value array for one attribute in
// the filter policy. Values must be an array of strings, numbers, booleans,
// or condition objects.
func validateFilterPolicyValue(raw interface{}) error {
	values, ok := raw.([]interface{})
	if !ok {
		return NewInvalidParameter("Invalid filter policy: value must be an array")
	}

	for _, v := range values {
		switch cond := v.(type) {
		case string, float64, bool:
			continue
		case map[string]interface{}:
			if len(cond) != 1 {
				return NewInvalidParameter(
					"Invalid filter policy: a condition object must have exactly one operator")
			}
			for operator := range cond {
				switch operator {
				case "prefix", "anything-but", "numeric", "exists":
				default:
					return NewInvalidParameter(
						fmt.Sprintf("Invalid filter policy: unknown operator %q", operator))
				}
			}
		default:
			return NewInvalidParameter("Invalid filter policy: unsupported value type")
		}
	}

	return nil
}

// validateFilterPolicyScope validates the FilterPolicyScope attribute.
// AWS accepts only "MessageAttributes" (default) or "MessageBody".
func validateFilterPolicyScope(value string) error {
	switch value {
	case "MessageAttributes", "MessageBody":
		return nil
	default:
		return NewInvalidParameter(
			fmt.Sprintf("Invalid FilterPolicyScope: %s. Valid values: MessageAttributes, MessageBody", value))
	}
}

// validateRedrivePolicy validates the JSON structure of a subscription
// redrive policy. AWS requires deadLetterTargetArn.
func validateRedrivePolicy(value string) error {
	if value == "" {
		return nil
	}

	var rp struct {
		DeadLetterTargetArn string `json:"deadLetterTargetArn"`
	}
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
// this validator free of store-layer dependencies.
func validatePublishParams(isFifo, isContentBasedDedup bool, message, subject, messageStructure, messageGroupId, messageDeduplicationId string) error {
	if len(message) > maxMessageSize {
		return awserrors.NewAWSError("InvalidParameter", fmt.Sprintf("Message too long: %d bytes (maximum %d)", len(message), maxMessageSize), 400)
	}

	// Subject is documented as "UTF-8 text … less than 100 characters
	// long", so the ceiling counts Unicode characters.
	if n := utf8.RuneCountInString(subject); n > maxSubjectLength {
		return NewInvalidParameter(fmt.Sprintf("Subject too long: %d characters (maximum %d)", n, maxSubjectLength))
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
		if messageGroupId != "" {
			return NewInvalidParameter("MessageGroupId is only valid for FIFO topics")
		}
		if messageDeduplicationId != "" {
			return NewInvalidParameter("MessageDeduplicationId is only valid for FIFO topics")
		}
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
	if len(stringValue) > maxMessageAttrStringValue {
		return NewInvalidParameter(fmt.Sprintf("Message attribute %q StringValue too long: %d characters (maximum %d)", name, len(stringValue), maxMessageAttrStringValue))
	}
	if len(binaryValue) > maxMessageAttrBinaryValue {
		return NewInvalidParameter(fmt.Sprintf("Message attribute %q BinaryValue too long: %d bytes (maximum %d)", name, len(binaryValue), maxMessageAttrBinaryValue))
	}
	return nil
}

// ---------------------------------------------------------------------------
// Platform validators
// ---------------------------------------------------------------------------

// validatePlatformApplicationArn validates the structure of a platform
// application ARN. The ARN must start with "arn:", use the "sns" service,
// and have at least six colon-separated parts.
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
	return nil
}

// validatePlatformAttributeValue enforces a length cap on platform application
// and endpoint attribute values for DoS protection.
func validatePlatformAttributeValue(name, value string) error {
	if len(value) > maxPlatformAttributeValueLength {
		return NewInvalidParameter(fmt.Sprintf("%s value too long: %d characters (maximum %d)", name, len(value), maxPlatformAttributeValueLength))
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
