package cloudtrail

import (
	"fmt"
	"net"
	"regexp"

	"vorpalstacks/internal/common/bucketname"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// ---------------------------------------------------------------------------
// Smithy / AWS-docs derived patterns
// ---------------------------------------------------------------------------

// trailNamePattern matches the trail-name charset: "Contain only ASCII
// letters (a-z, A-Z), numbers (0-9), periods (.), underscores (_), or
// dashes (-)" (InvalidTrailNameException). The remaining documented rules —
// start and end with a letter or number, no adjacent periods/underscores/
// dashes, and not an IP address — are enforced procedurally in
// validateTrailName because Go's RE2 has no lookahead.
var trailNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// edsNamePattern matches the Smithy model constraint for EventDataStoreName:
// ^[a-zA-Z0-9._\-]+$, length 3-128. ChannelName carries the same constraint.
var edsNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._\-]+$`)

// destinationLocationPattern matches the Smithy model constraint for
// Destination.Location: ^[a-zA-Z0-9._/\-:*]+$, length 3-1024.
var destinationLocationPattern = regexp.MustCompile(`^[a-zA-Z0-9._/\-:*]+$`)

// edsKmsKeyIDPattern matches the Smithy model constraint for
// EventDataStoreKmsKeyId: ^[a-zA-Z0-9._/\-:]+$, length 1-350.
var edsKmsKeyIDPattern = regexp.MustCompile(`^[a-zA-Z0-9._/\-:]+$`)

// advancedFieldPattern matches the Smithy model constraint for
// AdvancedFieldSelector.Field (SelectorField):
// ^[\w|\d|\.|_]+$, length 1-1000.
var advancedFieldPattern = regexp.MustCompile(`^[\w|\d|\.|_]+$`)

// promptPattern matches the Smithy model pattern for Prompt and
// SearchSampleQueriesSearchPhrase: printable ASCII plus newline.
var promptPattern = regexp.MustCompile(`^[ -~\n]*$`)

// s3KeyPrefixPattern validates S3 key prefixes: 0-200 chars, printable ASCII
// excluding leading slashes.
var s3KeyPrefixPattern = regexp.MustCompile(`^[a-zA-Z0-9!_\-.()*'&$@:?+/]+$`)

// snsTopicNamePattern validates SNS topic names: alphanumeric, hyphens, and
// underscores, 1-256 characters.
var snsTopicNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,256}$`)

// kmsKeyIDPattern accepts the four documented KmsKeyId forms for trails:
// "The value can be an alias name prefixed by alias/, a fully specified ARN
// to an alias, a fully specified ARN to a key, or a globally unique
// identifier" (CreateTrail/UpdateTrail KmsKeyId).
var kmsKeyIDPattern = regexp.MustCompile(`^(alias/[a-zA-Z0-9/_-]+|arn:aws:kms:[a-z0-9-]+:[0-9]{12}:(alias/[a-zA-Z0-9/_-]+|key/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})|[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})$`)

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

// validDestinationTypes contains the Smithy DestinationType enum values.
var validDestinationTypes = map[string]bool{
	"EVENT_DATA_STORE": true,
	"AWS_SERVICE":      true,
}

// validImportStatuses contains the Smithy ImportStatus enum values that
// ListImports accepts as a filter.
var validImportStatuses = map[string]bool{
	"INITIALIZING": true,
	"IN_PROGRESS":  true,
	"FAILED":       true,
	"STOPPED":      true,
	"COMPLETED":    true,
}

// validAggregationTemplates contains the Smithy Template enum values for
// AggregationConfiguration.Templates.
var validAggregationTemplates = map[string]bool{
	"API_ACTIVITY":    true,
	"RESOURCE_ACCESS": true,
	"USER_ACTIONS":    true,
}

// validContextKeyTypes contains the Smithy Type enum values for
// ContextKeySelector.Type.
var validContextKeyTypes = map[string]bool{
	"TagContext":     true,
	"RequestContext": true,
}

// basicDataResourceTypes are the resource types a basic event selector can
// select: "You can use EventSelectors to log management events and data
// events for the following resource types: AWS::DynamoDB::Table,
// AWS::Lambda::Function, AWS::S3::Object" (PutEventSelectors).
var basicDataResourceTypes = map[string]bool{
	"AWS::DynamoDB::Table":  true,
	"AWS::Lambda::Function": true,
	"AWS::S3::Object":       true,
}

// ---------------------------------------------------------------------------
// Validation functions
// ---------------------------------------------------------------------------

// validateTrailName validates a trail name against the AWS rules: the
// charset pattern, the alphanumeric start/end anchors, the
// no-adjacent-separator rule, and the separate "Not be in IP address
// format" rule.
func validateTrailName(name string) error {
	if len(name) < 3 || len(name) > 128 {
		return newInvalidTrailNameException(
			"Trail name must be between 3 and 128 characters")
	}
	if !trailNamePattern.MatchString(name) {
		return newInvalidTrailNameException(
			"Trail name must contain only letters, numbers, periods, underscores, and dashes")
	}
	if !isTrailNameAlnum(name[0]) || !isTrailNameAlnum(name[len(name)-1]) {
		return newInvalidTrailNameException(
			"Trail name must start and end with a letter or number")
	}
	for i := 1; i < len(name); i++ {
		if isTrailNameSeparator(name[i-1]) && isTrailNameSeparator(name[i]) {
			return newInvalidTrailNameException(
				"Trail name must not have adjacent periods, underscores, or dashes")
		}
	}
	if net.ParseIP(name) != nil {
		return newInvalidTrailNameException(
			"Trail name must not be in IP address format")
	}
	return nil
}

func isTrailNameAlnum(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9'
}

func isTrailNameSeparator(c byte) bool {
	return c == '.' || c == '_' || c == '-'
}

// validateS3BucketName validates an S3 bucket name against the AWS
// general-purpose bucket naming rules.
func validateS3BucketName(name string) error {
	if !bucketname.Validate(name) {
		return newInvalidS3BucketNameException(
			"S3 bucket name must be 3-63 characters of lowercase letters, numbers, hyphens and periods, must start and end with a letter or number, must not contain adjacent dots or dot-hyphen pairs, must not look like an IP address, and must not use a reserved prefix or suffix")
	}
	return nil
}

// validateS3KeyPrefix enforces the S3KeyPrefix bounds; both violations are
// the model's InvalidS3PrefixException ("This exception is thrown when the
// provided S3 prefix is not valid", declared on CreateTrail/UpdateTrail).
func validateS3KeyPrefix(prefix string) error {
	if len(prefix) > cloudtrailstore.MaxS3KeyPrefixLength {
		return newInvalidS3PrefixException(
			fmt.Sprintf("S3KeyPrefix must be %d characters or fewer", cloudtrailstore.MaxS3KeyPrefixLength))
	}
	if prefix != "" && !s3KeyPrefixPattern.MatchString(prefix) {
		return newInvalidS3PrefixException(
			"S3KeyPrefix contains invalid characters")
	}
	return nil
}

func validateSnsTopicName(name string) error {
	if !snsTopicNamePattern.MatchString(name) {
		return newInvalidSnsTopicNameException(
			"SnsTopicName must be 1-256 alphanumeric characters, hyphens, or underscores")
	}
	return nil
}

func validateKMSKeyID(keyID string) error {
	if !kmsKeyIDPattern.MatchString(keyID) {
		return newInvalidKmsKeyIdException(
			"KmsKeyId must be an alias name, a KMS key or alias ARN, or a key ID")
	}
	return nil
}

// validateEventDataStoreKMSKeyID validates an EDS KmsKeyId against the
// model's EventDataStoreKmsKeyId constraint (length 1-350, pattern
// ^[a-zA-Z0-9._/\-:]+$), which is the documented superset admitting the
// alias, alias-ARN, key-ARN and key-ID forms.
func validateEventDataStoreKMSKeyID(keyID string) error {
	if len(keyID) < 1 || len(keyID) > 350 || !edsKmsKeyIDPattern.MatchString(keyID) {
		return newInvalidKmsKeyIdException(
			"KmsKeyId must be a valid KMS key identifier")
	}
	return nil
}

func validateCloudWatchLogsLogGroupARN(arn string) error {
	if !cwlLogGroupArnPattern.MatchString(arn) {
		return newInvalidCloudWatchLogsLogGroupArnException(
			"CloudWatchLogsLogGroupArn must be a valid CloudWatch Logs log group ARN")
	}
	return nil
}

func validateCloudWatchLogsRoleARN(arn string) error {
	if !cwlRoleArnPattern.MatchString(arn) {
		return newInvalidCloudWatchLogsRoleArnException(
			"CloudWatchLogsRoleArn must be a valid IAM role ARN")
	}
	return nil
}

// validateEventDataStoreName validates the EDS name against AWS spec
// (Smithy model: length 3-128, pattern ^[a-zA-Z0-9._\-]+$).
func validateEventDataStoreName(name string) error {
	if len(name) < 3 || len(name) > 128 {
		return newInvalidParameterException(
			"Event data store name must be between 3 and 128 characters")
	}
	if !edsNamePattern.MatchString(name) {
		return newInvalidParameterException(
			"Event data store name contains invalid characters")
	}
	return nil
}

// validateReadWriteType returns an error when the value is not one of the
// Smithy ReadWriteType enum values.
func validateReadWriteType(v string) error {
	if !validReadWriteTypes[v] {
		return newInvalidEventSelectorsException(
			"ReadWriteType must be one of: ReadOnly, WriteOnly, All")
	}
	return nil
}

// validateInsightType returns an error when the value is not one of the
// Smithy InsightType enum values.
func validateInsightType(v string) error {
	if !validInsightTypes[v] {
		return newInvalidInsightSelectorsException(
			"InsightType must be one of: ApiCallRateInsight, ApiErrorRateInsight")
	}
	return nil
}

// validateLookupAttributeKey returns an error when the key is not one of the
// Smithy LookupAttributeKey enum values.
func validateLookupAttributeKey(k string) error {
	if !validLookupAttributeKeys[k] {
		return newInvalidLookupAttributesException(
			"Invalid AttributeKey: " + k)
	}
	return nil
}

// validateBillingMode returns an error when the value is not one of the
// Smithy BillingMode enum values.
func validateBillingMode(v string) error {
	if !validBillingModes[v] {
		return newInvalidParameterException(
			"BillingMode must be EXTENDABLE_RETENTION_PRICING or FIXED_RETENTION_PRICING")
	}
	return nil
}

// validateQueryStatus returns an error when the value is not one of the
// Smithy QueryStatus enum values; ListQueries declares
// InvalidQueryStatusException for a bad status filter.
func validateQueryStatus(v string) error {
	if !validQueryStatuses[v] {
		return newInvalidQueryStatusException(
			"QueryStatus must be one of: QUEUED, RUNNING, FINISHED, FAILED, CANCELLED, TIMED_OUT")
	}
	return nil
}

func validateEventDataStoreStatus(v string) error {
	if !validEventDataStoreStatuses[v] {
		return newInvalidEventDataStoreStatusException(
			"EventDataStoreStatus must be one of: CREATED, ENABLED, PENDING_DELETION, STARTING_INGESTION, STOPPING_INGESTION, STOPPED_INGESTION")
	}
	return nil
}

// validateChannelName validates a channel name against the model's
// ChannelName constraint: length 3-128, pattern ^[a-zA-Z0-9._\-]+$. The
// channel operations declare no name-specific error shape, so violations
// surface as their generic InvalidParameterException.
func validateChannelName(name string) error {
	if len(name) < 3 || len(name) > 128 || !edsNamePattern.MatchString(name) {
		return newInvalidParameterException(
			"Channel name must be 3-128 characters of letters, numbers, periods, underscores, and dashes")
	}
	return nil
}

// validateChannelSource validates a channel source against the model's
// Source constraint (length 1-256). The Source-specific violations surface
// as the operation's declared InvalidSourceException.
func validateChannelSource(source string) error {
	if len(source) < 1 || len(source) > 256 {
		return newInvalidSourceException(
			"Source must be between 1 and 256 characters")
	}
	return nil
}

// validateDestinations validates a parsed destination list against the
// model's Destinations and Destination constraints: 1-200 entries, each
// carrying the DestinationType enum in Type and a 3-1024 character
// Location matching the model pattern.
func validateDestinations(dests []cloudtrailstore.Destination) error {
	if len(dests) < 1 || len(dests) > 200 {
		return newInvalidParameterException(
			"Destinations must contain between 1 and 200 entries")
	}
	for _, d := range dests {
		if !validDestinationTypes[d.Type] {
			return newInvalidParameterException(
				"Destination.Type must be EVENT_DATA_STORE or AWS_SERVICE")
		}
		if len(d.Location) < 3 || len(d.Location) > 1024 || !destinationLocationPattern.MatchString(d.Location) {
			return newInvalidParameterException(
				"Destination.Location must be 3-1024 characters of letters, numbers, periods, slashes, underscores, dashes, colons, and asterisks")
		}
	}
	return nil
}

// validateImportStatusFilter validates a ListImports ImportStatus filter
// value against the model enum.
func validateImportStatusFilter(v string) error {
	if !validImportStatuses[v] {
		return newInvalidParameterException(
			"ImportStatus must be one of: INITIALIZING, IN_PROGRESS, FAILED, STOPPED, COMPLETED")
	}
	return nil
}
