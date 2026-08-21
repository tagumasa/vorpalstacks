package cloudtrail

import (
	"net/http"
	"regexp"

	"vorpalstacks/internal/common/bucketname"
	awserrors "vorpalstacks/internal/common/errors"
)

// ---------------------------------------------------------------------------
// Smithy / AWS-docs derived patterns
// ---------------------------------------------------------------------------

// trailNamePattern validates trail names per AWS docs: "Contains only ASCII
// letters (a-z, A-Z), numbers (0-9), periods (.), underscores (_), or dashes
// (-), must start with a letter or number, between 3 and 128 characters."
var trailNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{2,127}$`)

// edsNamePattern matches the Smithy model constraint for EventDataStoreName:
// ^[a-zA-Z0-9._\-]+$, length 3-128.
var edsNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._\-]+$`)

// s3KeyPrefixPattern validates S3 key prefixes: 0-200 chars, printable ASCII
// excluding leading slashes.
var s3KeyPrefixPattern = regexp.MustCompile(`^[a-zA-Z0-9!_\-.()*'&$@:?+/]+$`)

// snsTopicNamePattern validates SNS topic names: alphanumeric, hyphens, and
// underscores, 1-256 characters.
var snsTopicNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,256}$`)

// kmsKeyIDPattern accepts either a UUID-style key id (e.g. "abcd1234-...")
// or a full KMS key ARN.
var kmsKeyIDPattern = regexp.MustCompile(`^(arn:aws:kms:[a-z0-9-]+:\d{12}:key/)?[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// cwlLogGroupArnPattern validates CloudWatch Logs log group ARNs.
var cwlLogGroupArnPattern = regexp.MustCompile(`^arn:aws:logs:[a-z0-9-]+:\d{12}:log-group:[^:]+$`)

// cwlRoleArnPattern validates IAM role ARNs used for CloudWatch Logs
// delivery.
var cwlRoleArnPattern = regexp.MustCompile(`^arn:aws:iam::\d{12}:role/[^:]+$`)

// ---------------------------------------------------------------------------
// Smithy-derived enum sets
// ---------------------------------------------------------------------------

// validReadWriteTypes contains the Smithy ReadWriteType enum values.
var validReadWriteTypes = map[string]bool{
	"ReadOnly":  true,
	"WriteOnly": true,
	"All":       true,
}

// validInsightTypes contains the Smithy InsightType enum values.
var validInsightTypes = map[string]bool{
	"ApiCallRateInsight":  true,
	"ApiErrorRateInsight": true,
}

// validLookupAttributeKeys contains the Smithy LookupAttributeKey enum values.
var validLookupAttributeKeys = map[string]bool{
	"EventId":      true,
	"EventName":    true,
	"ReadOnly":     true,
	"Username":     true,
	"ResourceType": true,
	"ResourceName": true,
	"EventSource":  true,
	"AccessKeyId":  true,
}

// validQueryStatuses contains the Smithy QueryStatus enum values.
var validQueryStatuses = map[string]bool{
	"QUEUED":    true,
	"RUNNING":   true,
	"FINISHED":  true,
	"FAILED":    true,
	"CANCELLED": true,
	"TIMED_OUT": true,
}

// validEventDataStoreStatuses contains the Smithy EventDataStoreStatus enum
// values used by Start/StopEventDataStoreIngestion state checks.
var validEventDataStoreStatuses = map[string]bool{
	"CREATED":            true,
	"ENABLED":            true,
	"PENDING_DELETION":   true,
	"STARTING_INGESTION": true,
	"STOPPING_INGESTION": true,
	"STOPPED_INGESTION":  true,
}

// validBillingModes contains the Smithy BillingMode enum values.
var validBillingModes = map[string]bool{
	"EXTENDABLE_RETENTION_PRICING": true,
	"FIXED_RETENTION_PRICING":      true,
}

// ---------------------------------------------------------------------------
// Validation functions
// ---------------------------------------------------------------------------

// validateTrailName validates a trail name against AWS docs constraints.
func validateTrailName(name string) error {
	if len(name) < 3 || len(name) > 128 {
		return awserrors.NewAWSError("InvalidTrailNameException",
			"Trail name must be between 3 and 128 characters", http.StatusBadRequest)
	}
	if !trailNamePattern.MatchString(name) {
		return awserrors.NewAWSError("InvalidTrailNameException",
			"Trail name must contain only letters, numbers, periods, underscores, and dashes, and must start with a letter or number",
			http.StatusBadRequest)
	}
	return nil
}

// validateS3BucketName validates an S3 bucket name against the AWS
// general-purpose bucket naming rules.
func validateS3BucketName(name string) error {
	if !bucketname.Validate(name) {
		return awserrors.NewAWSError("InvalidS3BucketNameException",
			"S3 bucket name must be 3-63 characters of lowercase letters, numbers, hyphens and periods, must start and end with a letter or number, must not contain adjacent dots or dot-hyphen pairs, must not look like an IP address, and must not use a reserved prefix or suffix",
			http.StatusBadRequest)
	}
	return nil
}

func validateS3KeyPrefix(prefix string) error {
	if len(prefix) > 200 {
		return awserrors.NewAWSError("InvalidParameterException",
			"S3KeyPrefix must be 200 characters or fewer", http.StatusBadRequest)
	}
	if prefix != "" && !s3KeyPrefixPattern.MatchString(prefix) {
		return awserrors.NewAWSError("InvalidParameterException",
			"S3KeyPrefix contains invalid characters", http.StatusBadRequest)
	}
	return nil
}

func validateSnsTopicName(name string) error {
	if !snsTopicNamePattern.MatchString(name) {
		return awserrors.NewAWSError("InvalidSnsTopicNameException",
			"SnsTopicName must be 1-256 alphanumeric characters, hyphens, or underscores",
			http.StatusBadRequest)
	}
	return nil
}

func validateKMSKeyID(keyID string) error {
	if !kmsKeyIDPattern.MatchString(keyID) {
		return awserrors.NewAWSError("InvalidKmsKeyIdException",
			"KmsKeyId must be a valid UUID or KMS key ARN", http.StatusBadRequest)
	}
	return nil
}

func validateCloudWatchLogsLogGroupARN(arn string) error {
	if !cwlLogGroupArnPattern.MatchString(arn) {
		return awserrors.NewAWSError("InvalidCloudWatchLogsLogGroupArnException",
			"CloudWatchLogsLogGroupArn must be a valid CloudWatch Logs log group ARN",
			http.StatusBadRequest)
	}
	return nil
}

func validateCloudWatchLogsRoleARN(arn string) error {
	if !cwlRoleArnPattern.MatchString(arn) {
		return awserrors.NewAWSError("InvalidCloudWatchLogsRoleArnException",
			"CloudWatchLogsRoleArn must be a valid IAM role ARN",
			http.StatusBadRequest)
	}
	return nil
}

// validateEventDataStoreName validates the EDS name against AWS spec
// (Smithy model: length 3-128, pattern ^[a-zA-Z0-9._\-]+$).
func validateEventDataStoreName(name string) error {
	if len(name) < 3 || len(name) > 128 {
		return awserrors.NewAWSError("InvalidParameterException",
			"Event data store name must be between 3 and 128 characters", http.StatusBadRequest)
	}
	if !edsNamePattern.MatchString(name) {
		return awserrors.NewAWSError("InvalidParameterException",
			"Event data store name contains invalid characters", http.StatusBadRequest)
	}
	return nil
}

// validateReadWriteType returns an error when the value is not one of the
// Smithy ReadWriteType enum values.
func validateReadWriteType(v string) error {
	if !validReadWriteTypes[v] {
		return awserrors.NewAWSError("InvalidEventSelectorsException",
			"ReadWriteType must be one of: ReadOnly, WriteOnly, All", http.StatusBadRequest)
	}
	return nil
}

// validateInsightType returns an error when the value is not one of the
// Smithy InsightType enum values.
func validateInsightType(v string) error {
	if !validInsightTypes[v] {
		return awserrors.NewAWSError("InvalidInsightSelectorsException",
			"InsightType must be one of: ApiCallRateInsight, ApiErrorRateInsight", http.StatusBadRequest)
	}
	return nil
}

// validateLookupAttributeKey returns an error when the key is not one of the
// Smithy LookupAttributeKey enum values.
func validateLookupAttributeKey(k string) error {
	if !validLookupAttributeKeys[k] {
		return awserrors.NewAWSError("InvalidLookupAttributesException",
			"Invalid AttributeKey: "+k, http.StatusBadRequest)
	}
	return nil
}

// validateBillingMode returns an error when the value is not one of the
// Smithy BillingMode enum values.
func validateBillingMode(v string) error {
	if !validBillingModes[v] {
		return awserrors.NewAWSError("InvalidParameterException",
			"BillingMode must be EXTENDABLE_RETENTION_PRICING or FIXED_RETENTION_PRICING", http.StatusBadRequest)
	}
	return nil
}

func validateQueryStatus(v string) error {
	if !validQueryStatuses[v] {
		return awserrors.NewAWSError("InvalidParameterException",
			"QueryStatus must be one of: QUEUED, RUNNING, FINISHED, FAILED, CANCELLED, TIMED_OUT",
			http.StatusBadRequest)
	}
	return nil
}

func validateEventDataStoreStatus(v string) error {
	if !validEventDataStoreStatuses[v] {
		return awserrors.NewAWSError("InvalidEventDataStoreStatusException",
			"EventDataStoreStatus must be one of: CREATED, ENABLED, PENDING_DELETION, STARTING_INGESTION, STOPPING_INGESTION, STOPPED_INGESTION",
			http.StatusBadRequest)
	}
	return nil
}
