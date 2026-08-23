package sqs

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	tagutil "vorpalstacks/internal/common/tags"
)

var (
	queueNameRegex      = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	fifoQueueNameRegex  = regexp.MustCompile(`^[a-zA-Z0-9_-]+\.fifo$`)
	batchEntryIdRegex   = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	validAttributeNames = map[string]bool{
		"All":                                   true,
		"QueueArn":                              true,
		"ApproximateNumberOfMessages":           true,
		"ApproximateNumberOfMessagesDelayed":    true,
		"ApproximateNumberOfMessagesNotVisible": true,
		"CreatedTimestamp":                      true,
		"LastModifiedTimestamp":                 true,
		"VisibilityTimeout":                     true,
		"MaximumMessageSize":                    true,
		"MessageRetentionPeriod":                true,
		"DelaySeconds":                          true,
		"ReceiveMessageWaitTimeSeconds":         true,
		"Policy":                                true,
		"RedrivePolicy":                         true,
		"FifoQueue":                             true,
		"ContentBasedDeduplication":             true,
		"KmsMasterKeyId":                        true,
		"KmsDataKeyReusePeriodSeconds":          true,
		"DeduplicationScope":                    true,
		"FifoThroughputLimit":                   true,
		"RedriveAllowPolicy":                    true,
		"SqsManagedSseEnabled":                  true,
	}
)

func validateQueueName(name string) error {
	if len(name) == 0 {
		return ErrInvalidQueueName
	}
	if len(name) > MaxQueueNameLength {
		return ErrInvalidQueueName
	}
	if strings.HasSuffix(name, ".fifo") {
		if !fifoQueueNameRegex.MatchString(name) {
			return ErrInvalidQueueName
		}
	} else {
		if !queueNameRegex.MatchString(name) {
			return ErrInvalidQueueName
		}
	}
	return nil
}

// ValidateBatchEntryId validates a batch entry ID.
func ValidateBatchEntryId(id string) error {
	if len(id) == 0 || len(id) > maxBatchEntryIdLength {
		return ErrInvalidBatchEntryId
	}
	if !batchEntryIdRegex.MatchString(id) {
		return ErrInvalidBatchEntryId
	}
	return nil
}

// IsValidAttributeName checks if an attribute name is valid for SQS queues.
func IsValidAttributeName(name string) bool {
	return validAttributeNames[name]
}

func validateVisibilityTimeout(value int32) error {
	if value < MinVisibilityTimeout || value > MaxVisibilityTimeout {
		return ErrInvalidParameterValue
	}
	return nil
}

func validateDelaySeconds(value int32) error {
	if value < MinDelaySeconds || value > MaxDelaySeconds {
		return ErrInvalidParameterValue
	}
	return nil
}

func validateMessageRetentionPeriod(value int32) error {
	if value < MinMessageRetentionPeriod || value > MaxMessageRetentionPeriod {
		return ErrInvalidParameterValue
	}
	return nil
}

func validateMaximumMessageSize(value int32) error {
	if value < MinMaximumMessageSize || value > MaxMaximumMessageSize {
		return ErrInvalidParameterValue
	}
	return nil
}

func validateReceiveMessageWaitTimeSeconds(value int32) error {
	if value < MinReceiveMessageWaitTimeSeconds || value > MaxReceiveMessageWaitTimeSeconds {
		return ErrInvalidParameterValue
	}
	return nil
}

// validateTags validates queue tags against the SQS tag limits: at most 50
// tags per queue, keys of 1-128 characters, values of at most 256 characters
// and the aws: key prefix reserved for AWS use.
func validateTags(tags map[string]string) error {
	switch v, _ := tagutil.CheckStringTags(tags, tagutil.StandardLimits()); v {
	case tagutil.TooManyTags:
		return ErrTooManyTags
	case tagutil.TagKeyTooShort, tagutil.TagKeyTooLong:
		return ErrInvalidTagKey
	case tagutil.TagValueTooLong:
		return ErrInvalidTagValue
	case tagutil.ReservedTagKey:
		return ErrInvalidTagKey
	}
	return nil
}

// ---------------------------------------------------------------------------
// Newer SQS attribute validation
// ---------------------------------------------------------------------------

var (
	awsAccountRegex = regexp.MustCompile(`^[0-9]{12}$`)
	dataTypeRegex   = regexp.MustCompile(`^(String|Number|Binary)(\..+)?$`)
)

const (
	maxPermissionLabels   = 10
	maxPermissionAccounts = 10
	maxPermissionLabelLen = 80
	// MaxActionsPerStatement is the maximum number of actions per
	// AddPermission statement (7 per the AWS SQS API Reference: "An Amazon
	// SQS policy can have a maximum of seven actions per statement.").
	MaxActionsPerStatement = 7
)

var validDeduplicationScopes = map[string]bool{
	"messageGroup": true,
	"queue":        true,
}

var validFifoThroughputLimits = map[string]bool{
	"perMessageGroupId": true,
	"perQueue":          true,
}

var validRedrivePermissions = map[string]bool{
	"allowAll": true,
	"byQueue":  true,
	"denyAll":  true,
}

// sqsActionNameRegex matches SQS action names: "the name of any action or
// `*`" (AWS SQS API Reference, AddPermission Actions parameter). Action names
// are alphanumeric, so anything outside [A-Za-z0-9*] is rejected without
// keeping a closed list that would go stale as new actions are added.
var sqsActionNameRegex = regexp.MustCompile(`^[A-Za-z0-9*]+$`)

// validateMessageBody enforces the documented message-body character set:
// "#x9 | #xA | #xD | #x20 to #xD7FF | #xE000 to #xFFFD | #x10000 to
// #x10FFFF" — "If a message contains characters outside the allowed set,
// Amazon SQS rejects the message and returns an InvalidMessageContents
// error." (AWS SQS API Reference).
func validateMessageBody(body string) error {
	for _, r := range body {
		if r == 0x9 || r == 0xA || r == 0xD {
			continue
		}
		if r >= 0x20 && r <= 0xD7FF {
			continue
		}
		if r >= 0xE000 && r <= 0xFFFD {
			continue
		}
		if r >= 0x10000 && r <= 0x10FFFF {
			continue
		}
		return ErrInvalidMessageContents
	}
	return nil
}

// validateFifoIdentifier enforces the documented MessageGroupId and
// MessageDeduplicationId rules: at most 128 characters of "alphanumeric
// characters (a-z, A-Z, 0-9) and punctuation
// (!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~)" — every printable non-space ASCII
// character (AWS SQS API Reference). An empty identifier means the field was
// not provided and is validated by the FIFO presence rules instead.
func validateFifoIdentifier(value string) error {
	if value == "" {
		return nil
	}
	if len(value) > MaxFifoIdLength {
		return ErrInvalidParameterValue
	}
	for i := 0; i < len(value); i++ {
		if value[i] < 0x21 || value[i] > 0x7E {
			return ErrInvalidParameterValue
		}
	}
	return nil
}

func validateKmsDataKeyReusePeriod(value int32) error {
	if value < MinKmsDataKeyReusePeriodSeconds || value > MaxKmsDataKeyReusePeriodSeconds {
		return ErrInvalidParameterValue
	}
	return nil
}

func validateDeduplicationScope(value string) error {
	if !validDeduplicationScopes[value] {
		return ErrInvalidParameterValue
	}
	return nil
}

func validateFifoThroughputLimit(value string) error {
	if !validFifoThroughputLimits[value] {
		return ErrInvalidParameterValue
	}
	return nil
}

// ValidateDeduplicationScope applies the store deduplication-scope rule
// (valid values "messageGroup" and "queue") for callers outside the store
// package.
func ValidateDeduplicationScope(value string) error {
	return validateDeduplicationScope(value)
}

// ValidateFifoThroughputLimit applies the store FIFO throughput-limit rule
// (valid values "perQueue" and "perMessageGroupId") for callers outside the
// store package.
func ValidateFifoThroughputLimit(value string) error {
	return validateFifoThroughputLimit(value)
}

// validateHighThroughputFifo enforces the documented cross-rule on the merged
// attribute view of a queue: "The perMessageGroupId value is allowed only
// when the value for DeduplicationScope is messageGroup." An absent
// DeduplicationScope behaves as the documented default "queue".
func validateHighThroughputFifo(attrs map[string]string) error {
	if attrs["FifoThroughputLimit"] == "perMessageGroupId" && attrs["DeduplicationScope"] != "messageGroup" {
		return ErrInvalidParameterValue
	}
	return nil
}

// validateSSEExclusion enforces the documented server-side-encryption rule on
// the merged attribute view of a queue: "Only one server-side encryption
// option is supported per queue (for example, SSE-KMS or SSE-SQS)." An empty
// KmsMasterKeyId behaves as unset.
func validateSSEExclusion(attrs map[string]string) error {
	if attrs["KmsMasterKeyId"] != "" && attrs["SqsManagedSseEnabled"] == "true" {
		return ErrInvalidParameterValue
	}
	return nil
}

func validateSqsManagedSseEnabled(value string) error {
	if _, err := strconv.ParseBool(value); err != nil {
		return ErrInvalidParameterValue
	}
	return nil
}

// ValidateRedriveAllowPolicyJSON validates a RedriveAllowPolicy JSON string.
func ValidateRedriveAllowPolicyJSON(data string) error {
	return validateRedriveAllowPolicyJSON(data)
}

// validateRedriveAllowPolicyJSON validates the structure of a RedriveAllowPolicy
// JSON string. Valid fields: redrivePermission (enum), sourceQueueArns (list).
func validateRedriveAllowPolicyJSON(data string) error {
	if data == "" {
		return ErrInvalidParameterValue
	}
	var raw struct {
		RedrivePermission string   `json:"redrivePermission"`
		SourceQueueArns   []string `json:"sourceQueueArns"`
	}
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		return ErrInvalidParameterValue
	}
	if raw.RedrivePermission != "" && !validRedrivePermissions[raw.RedrivePermission] {
		return ErrInvalidParameterValue
	}
	if len(raw.SourceQueueArns) > 10 {
		return ErrInvalidParameterValue
	}
	return nil
}

// ValidatePolicyJSON validates that a Policy string is valid JSON.
func ValidatePolicyJSON(policy string) error {
	return validatePolicyJSON(policy)
}

// validatePolicyJSON validates that a Policy string is valid JSON.
func validatePolicyJSON(policy string) error {
	if policy == "" {
		return nil
	}
	var v interface{}
	if err := json.Unmarshal([]byte(policy), &v); err != nil {
		return ErrInvalidParameterValue
	}
	return nil
}

// ---------------------------------------------------------------------------
// AddPermission validation
// ---------------------------------------------------------------------------

var permissionLabelCharRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func validatePermissionLabel(label string) error {
	if label == "" || len(label) > maxPermissionLabelLen {
		return ErrInvalidParameterValue
	}
	if !permissionLabelCharRegex.MatchString(label) {
		return ErrInvalidParameterValue
	}
	return nil
}

// ValidatePermissionLabel applies the store permission-label rules (length
// and character set) so callers outside the store package share a single
// definition.
func ValidatePermissionLabel(label string) error {
	return validatePermissionLabel(label)
}

func validateAWSAccountIDs(ids []string) error {
	if len(ids) == 0 || len(ids) > maxPermissionAccounts {
		return ErrInvalidParameterValue
	}
	for _, id := range ids {
		if !awsAccountRegex.MatchString(id) {
			return ErrInvalidParameterValue
		}
	}
	return nil
}

func validateSQSActionList(actions []string) error {
	if len(actions) == 0 {
		return ErrInvalidParameterValue
	}
	if len(actions) > MaxActionsPerStatement {
		return ErrOverLimit
	}
	seen := make(map[string]bool, len(actions))
	for _, a := range actions {
		if !sqsActionNameRegex.MatchString(a) {
			return ErrInvalidParameterValue
		}
		if seen[a] {
			return ErrInvalidParameterValue
		}
		seen[a] = true
	}
	return nil
}

// ValidateMessageAttributeDataType validates the DataType field of a message
// attribute value. Must be String, Number, or Binary, optionally followed by
// a custom suffix (e.g. "Number.int").
func ValidateMessageAttributeDataType(dataType string) error {
	if !dataTypeRegex.MatchString(dataType) {
		return ErrInvalidDataType
	}
	return nil
}

// MaxMessageAttributes is the maximum number of message attributes per
// message (Amazon SQS Developer Guide: message metadata).
const MaxMessageAttributes = 10

// MaxMessageAttributeNameLength is the maximum length of a message attribute
// name in characters (Amazon SQS Developer Guide: message metadata).
const MaxMessageAttributeNameLength = 256

// MaxMessageAttributeDataTypeLength is the maximum length of a message
// attribute data type in characters (Amazon SQS Developer Guide: message
// metadata).
const MaxMessageAttributeDataTypeLength = 256

// messageAttributeNameRegex matches the allowed characters for a message
// attribute name: A-Z, a-z, 0-9, underscore, hyphen, and period.
var messageAttributeNameRegex = regexp.MustCompile(`^[A-Za-z0-9_\-.]+$`)

// ValidateMessageAttributes validates the full set of message attributes of
// an outbound message: the count cap, each attribute name (length, character
// set, reserved AWS./Amazon. prefixes, period placement rules), and each
// DataType (format plus length). Per the Developer Guide, all components of
// a message attribute count towards the message size restriction, which is
// enforced separately against the queue's MaximumMessageSize.
func ValidateMessageAttributes(attrs map[string]*MessageAttributeValue) error {
	if len(attrs) == 0 {
		return nil
	}
	if len(attrs) > MaxMessageAttributes {
		return ErrInvalidParameterValue
	}
	for name, attr := range attrs {
		if name == "" {
			return ErrInvalidParameterValue
		}
		// The charset check comes first so the length check below only ever
		// sees single-byte names: the allowed charset is ASCII-only, which
		// makes the byte-length check equivalent to the documented
		// character count.
		if !messageAttributeNameRegex.MatchString(name) {
			return ErrInvalidParameterValue
		}
		// Reserved prefixes are rejected in any casing variation.
		lower := strings.ToLower(name)
		if strings.HasPrefix(lower, "aws.") || strings.HasPrefix(lower, "amazon.") {
			return ErrInvalidParameterValue
		}
		// Names must not start or end with a period nor contain
		// consecutive periods.
		if strings.HasPrefix(name, ".") || strings.HasSuffix(name, ".") || strings.Contains(name, "..") {
			return ErrInvalidParameterValue
		}
		if len(name) > MaxMessageAttributeNameLength {
			return ErrInvalidParameterValue
		}
		if attr == nil || attr.DataType == "" || len(attr.DataType) > MaxMessageAttributeDataTypeLength {
			return ErrInvalidParameterValue
		}
		if err := ValidateMessageAttributeDataType(attr.DataType); err != nil {
			return err
		}
	}
	return nil
}

// validateKmsMasterKeyId validates a KMS key ID or alias. Accepted forms:
//   - UUID format (36 chars, e.g. "12345678-1234-1234-1234-123456789012")
//   - alias/ prefix (e.g. "alias/my-key")
//   - key/ prefix (e.g. "key/12345678-1234-1234-1234-123456789012")
//   - arn:aws:kms:... full ARN
func validateKmsMasterKeyId(v string) error {
	if v == "" {
		return nil
	}
	if len(v) > 256 {
		return ErrInvalidParameterValue
	}
	uuidRe := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	if uuidRe.MatchString(v) {
		return nil
	}
	if strings.HasPrefix(v, "alias/") || strings.HasPrefix(v, "key/") {
		return nil
	}
	if strings.HasPrefix(v, "arn:aws:kms:") {
		return nil
	}
	return ErrInvalidParameterValue
}
