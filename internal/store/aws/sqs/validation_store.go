package sqs

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
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
	if len(name) > maxQueueNameLength {
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
	if value < minVisibilityTimeout || value > maxVisibilityTimeout {
		return ErrInvalidParameterValue
	}
	return nil
}

func validateDelaySeconds(value int32) error {
	if value < minDelaySeconds || value > maxDelaySeconds {
		return ErrInvalidParameterValue
	}
	return nil
}

func validateMessageRetentionPeriod(value int32) error {
	if value < minMessageRetentionPeriod || value > maxMessageRetentionPeriod {
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
	if value < minReceiveMessageWaitTime || value > maxReceiveMessageWaitTime {
		return ErrInvalidParameterValue
	}
	return nil
}

func validateTags(tags map[string]string) error {
	if len(tags) > maxTagsPerQueue {
		return ErrTooManyTags
	}
	for key, value := range tags {
		if len(key) > maxTagKeyLength {
			return ErrInvalidTagKey
		}
		if len(value) > maxTagValueLength {
			return ErrInvalidTagValue
		}
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
	maxPermissionLabels    = 10
	maxPermissionAccounts  = 10
	maxPermissionLabelLen  = 80
	maxActionsPerStatement = 7
)

var validDeduplicationScopes = map[string]bool{
	"queueMessageGroup": true,
	"queue":             true,
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

var validSQSActions = map[string]bool{
	"*":                            true,
	"SendMessage":                  true,
	"SendMessageBatch":             true,
	"ReceiveMessage":               true,
	"DeleteMessage":                true,
	"DeleteMessageBatch":           true,
	"ChangeMessageVisibility":      true,
	"ChangeMessageVisibilityBatch": true,
	"GetQueueAttributes":           true,
	"GetQueueUrl":                  true,
	"PurgeQueue":                   true,
	"SetQueueAttributes":           true,
	"AddPermission":                true,
	"RemovePermission":             true,
	"ListQueueTags":                true,
	"TagQueue":                     true,
	"UntagQueue":                   true,
	"ListDeadLetterSourceQueues":   true,
	"StartMessageMoveTask":         true,
	"CancelMessageMoveTask":        true,
	"ListMessageMoveTasks":         true,
}

func validateKmsDataKeyReusePeriod(value int32) error {
	if value < 60 || value > 86400 {
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
	if len(actions) > maxActionsPerStatement {
		return ErrOverLimit
	}
	seen := make(map[string]bool, len(actions))
	for _, a := range actions {
		if !validSQSActions[a] {
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
