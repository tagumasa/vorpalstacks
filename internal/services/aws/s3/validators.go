package s3

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/bucketname"
	tagutil "vorpalstacks/internal/common/tags"
)

// === Validation constants ===========================================

const (
	// Object key constraints (AWS: 1-1024 bytes).
	maxObjectKeyLength  = 1024
	maxObjectSize       = 5 * 1024 * 1024 * 1024 * 1024 // 5 TiB
	maxSingleUploadSize = 5 * 1024 * 1024 * 1024        // 5 GiB

	// Multipart upload constraints (AWS: 1-10000 parts, min 5 MiB except last).
	minPartNumber = 1
	maxPartNumber = 10000
	minPartSize   = 5 * 1024 * 1024 // 5 MiB

	// Lifecycle rule constraints (AWS: max 1000 rules per bucket).
	maxLifecycleRules = 1000
	maxLifecycleDays  = 3650

	// Notification filter rule constraints (AWS: value 1-1024 chars).
	maxFilterRuleValueLength = 1024

	// Restore object constraints (AWS: Days 1-3653).
	maxRestoreDays = 3653
)

// === Regex patterns ================================================

var validCORSMethods = map[string]bool{
	"GET":    true,
	"PUT":    true,
	"HEAD":   true,
	"POST":   true,
	"DELETE": true,
}

// validStorageClasses contains all S3 storage class enum values.
var validStorageClasses = map[string]bool{
	"STANDARD":            true,
	"REDUCED_REDUNDANCY":  true,
	"STANDARD_IA":         true,
	"ONEZONE_IA":          true,
	"INTELLIGENT_TIERING": true,
	"GLACIER":             true,
	"DEEP_ARCHIVE":        true,
	"GLACIER_IR":          true,
	"OUTPOSTS":            true,
	"EXPRESS_ONEZONE":     true,
	"SNOW":                true,
}

// === Existing validators (consolidated from bucket_operations.go,
//      object_operations.go) ========================================

// validateBucketName validates an S3 bucket name per the AWS
// general-purpose bucket naming rules (shared implementation in the
// bucketname package: charset, adjacency, IP form, reserved prefixes
// and suffixes).
func validateBucketName(name string) error {
	if !bucketname.Validate(name) {
		return NewInvalidBucketNameError(name)
	}
	return nil
}

// validateObjectKey validates an S3 object key (1-1024 bytes, no control
// characters except TAB/LF/CR, no null bytes).
func validateObjectKey(key string) error {
	if len(key) == 0 {
		return NewInvalidArgumentError("object key cannot be empty")
	}
	if len(key) > maxObjectKeyLength {
		return NewInvalidArgumentError("object key cannot exceed 1024 bytes")
	}
	if strings.Contains(key, "..") {
		return NewInvalidArgumentError("invalid object key: path traversal detected")
	}
	if strings.Contains(key, "\x00") {
		return NewInvalidArgumentError("invalid object key: null byte detected")
	}
	for i, r := range key {
		if r < 0x20 && r != 0x09 && r != 0x0A && r != 0x0D {
			return NewInvalidArgumentError(fmt.Sprintf("object key contains invalid control character at position %d", i))
		}
	}
	return nil
}

// validateTags validates a list of S3 tags (max 50, key ≤ 128, value ≤ 256,
// no "aws:" prefix).
func validateTags(tags []Tag) error {
	switch v, _ := tagutil.CheckTags(TagsToCommon(tags), tagutil.StandardLimits()); v {
	case tagutil.TooManyTags:
		return NewInvalidArgumentError(fmt.Sprintf("too many tags (maximum %d)", tagutil.MaxTagsPerResource))
	case tagutil.TagKeyTooShort:
		return NewInvalidArgumentError("tag key cannot be empty")
	case tagutil.TagKeyTooLong:
		return NewInvalidArgumentError(fmt.Sprintf("tag key cannot exceed %d characters", tagutil.MaxTagKeyLength))
	case tagutil.TagValueTooLong:
		return NewInvalidArgumentError(fmt.Sprintf("tag value cannot exceed %d characters", tagutil.MaxTagValueLength))
	case tagutil.ReservedTagKey:
		return NewInvalidArgumentError("tag key cannot start with 'aws:' (reserved prefix)")
	}
	return nil
}

// validatePartNumber validates a multipart upload part number (1-10000).
func validatePartNumber(partNumber int) error {
	if partNumber < minPartNumber || partNumber > maxPartNumber {
		return NewInvalidArgumentError(fmt.Sprintf("part number must be between %d and %d", minPartNumber, maxPartNumber))
	}
	return nil
}

// isPublicCannedACL returns true when the canned ACL grants public or
// authenticated-users access. Used by both PutBucketAcl and PutObjectAcl
// to enforce BlockPublicAcls consistently.
func isPublicCannedACL(acl string) bool {
	switch acl {
	case "public-read", "public-read-write", "authenticated-read":
		return true
	default:
		return false
	}
}

// validateAccelerateStatus validates that the Transfer Acceleration status
// is either "Enabled" or "Suspended".
func validateAccelerateStatus(status string) error {
	if status != "Enabled" && status != "Suspended" {
		return NewInvalidArgumentError(fmt.Sprintf("invalid accelerate status: %s (must be Enabled or Suspended)", status))
	}
	return nil
}

// validateRetentionMode validates an Object Lock retention mode.
// Must be "GOVERNANCE" or "COMPLIANCE".
func validateRetentionMode(mode string) error {
	if mode != "GOVERNANCE" && mode != "COMPLIANCE" {
		return NewInvalidArgumentError(fmt.Sprintf("invalid retention mode: %s (must be GOVERNANCE or COMPLIANCE)", mode))
	}
	return nil
}

// validateRetainUntilDate validates that the retain-until date is in the
// future. AWS requires RetainUntilDate to be a timestamp later than the
// current time.
func validateRetainUntilDate(date *time.Time) error {
	if date == nil {
		return NewInvalidArgumentError("RetainUntilDate is required when Mode is specified")
	}
	if !date.After(time.Now()) {
		return NewInvalidArgumentError("RetainUntilDate must be a future timestamp")
	}
	return nil
}

// validateS3EventName validates that an event name is a recognised S3
// event type. The full list is derived from the AWS S3 notification
// specification.
var validS3Events = map[string]bool{
	"s3:ObjectCreated:*":                               true,
	"s3:ObjectCreated:Put":                             true,
	"s3:ObjectCreated:Post":                            true,
	"s3:ObjectCreated:Copy":                            true,
	"s3:ObjectCreated:CompleteMultipartUpload":         true,
	"s3:ObjectRemoved:*":                               true,
	"s3:ObjectRemoved:Delete":                          true,
	"s3:ObjectRemoved:DeleteMarkerCreated":             true,
	"s3:ObjectRestore:*":                               true,
	"s3:ObjectRestore:Post":                            true,
	"s3:ObjectRestore:Completed":                       true,
	"s3:ObjectRestore:Delete":                          true,
	"s3:ReducedRedundancyLostObject":                   true,
	"s3:Replication:*":                                 true,
	"s3:Replication:OperationFailedReplication":        true,
	"s3:Replication:OperationCompletedReplication":     true,
	"s3:Replication:OperationMissedThreshold":          true,
	"s3:Replication:OperationReplicatedAfterThreshold": true,
	"s3:Replication:OperationNotTracked":               true,
	"s3:LifecycleExpiration:*":                         true,
	"s3:LifecycleExpiration:Delete":                    true,
	"s3:LifecycleExpiration:DeleteMarkerCreated":       true,
	"s3:LifecycleTransition":                           true,
	"s3:IntelligentTiering":                            true,
	"s3:ObjectAcl:Put":                                 true,
	"s3:ObjectTagging:*":                               true,
	"s3:ObjectTagging:Put":                             true,
	"s3:ObjectTagging:Delete":                          true,
}

// validateS3EventName checks that each event in the list is a valid S3
// event type. Wildcard suffixes (e.g. "s3:ObjectCreated:*") are accepted.
func validateS3EventNames(events []string) error {
	if len(events) == 0 {
		return NewInvalidArgumentError("at least one event must be specified")
	}
	for _, e := range events {
		if !validS3Events[e] {
			return NewInvalidArgumentError(fmt.Sprintf("invalid event name: %s", e))
		}
	}
	return nil
}

// validateFilterRule validates a notification filter rule. Name must be
// "prefix" or "suffix"; Value must be 1-1024 characters.
func validateFilterRule(name, value string) error {
	if name != "prefix" && name != "suffix" {
		return NewInvalidArgumentError(fmt.Sprintf("invalid filter rule name: %s (must be prefix or suffix)", name))
	}
	if len(value) == 0 {
		return NewInvalidArgumentError("filter rule value cannot be empty")
	}
	if len(value) > maxFilterRuleValueLength {
		return NewInvalidArgumentError(fmt.Sprintf("filter rule value exceeds maximum length of %d characters", maxFilterRuleValueLength))
	}
	return nil
}

// validatePolicyDocument performs basic structural validation of a JSON
// bucket policy document. It checks:
//   - The document is valid JSON
//   - A "Statement" array exists and is non-empty
//   - Each statement has an "Effect" of "Allow" or "Deny"
//
// Full IAM policy syntax validation (Action/Resource format, principal
// format, condition operators) is intentionally out of scope for this
// edge implementation.
func validatePolicyDocument(policy string) error {
	var doc struct {
		Statement []struct {
			Effect string `json:"Effect"`
		} `json:"Statement"`
	}
	if err := json.Unmarshal([]byte(policy), &doc); err != nil {
		return NewInvalidArgumentError(fmt.Sprintf("invalid policy JSON: %s", err.Error()))
	}
	if len(doc.Statement) == 0 {
		return NewInvalidArgumentError("policy must contain at least one statement")
	}
	for i, stmt := range doc.Statement {
		if stmt.Effect != "Allow" && stmt.Effect != "Deny" {
			return NewInvalidArgumentError(fmt.Sprintf("statement %d has invalid Effect: %s (must be Allow or Deny)", i, stmt.Effect))
		}
	}
	return nil
}

// validateLifecycleRules performs comprehensive validation of lifecycle
// rules including count limit, ID uniqueness, and per-field range checks.
func validateLifecycleRules(rules []LifecycleRuleInput) error {
	if len(rules) == 0 {
		return NewInvalidArgumentError("lifecycle configuration must contain at least one rule")
	}
	if len(rules) > maxLifecycleRules {
		return NewInvalidArgumentError(fmt.Sprintf("too many lifecycle rules (maximum %d)", maxLifecycleRules))
	}

	seenIDs := make(map[string]bool)
	for _, rule := range rules {
		if rule.Status != "Enabled" && rule.Status != "Disabled" {
			return NewInvalidArgumentError("rule status must be Enabled or Disabled")
		}

		if rule.ID != "" {
			if seenIDs[rule.ID] {
				return NewInvalidArgumentError(fmt.Sprintf("duplicate rule ID: %s", rule.ID))
			}
			seenIDs[rule.ID] = true
		}

		if rule.Expiration != nil {
			hasDays := rule.Expiration.Days != nil
			hasDate := rule.Expiration.Date != nil
			if hasDays && hasDate {
				return NewInvalidArgumentError("Expiration Days and Date are mutually exclusive")
			}
			if hasDays && (*rule.Expiration.Days < 1 || *rule.Expiration.Days > maxLifecycleDays) {
				return NewInvalidArgumentError(fmt.Sprintf("Expiration Days must be between 1 and %d, got %d", maxLifecycleDays, *rule.Expiration.Days))
			}
		}

		for _, t := range rule.Transitions {
			hasDays := t.Days != nil
			hasDate := t.Date != nil
			if hasDays && hasDate {
				return NewInvalidArgumentError("Transition Days and Date are mutually exclusive")
			}
			if hasDays && (*t.Days < 0 || *t.Days > maxLifecycleDays) {
				return NewInvalidArgumentError(fmt.Sprintf("Transition Days must be between 0 and %d, got %d", maxLifecycleDays, *t.Days))
			}
			if t.StorageClass != "" {
				if err := validateStorageClass(t.StorageClass); err != nil {
					return err
				}
			}
		}

		for _, t := range rule.NoncurrentVersionTransitions {
			if t.NoncurrentDays != nil {
				nd := *t.NoncurrentDays
				if nd < 1 || nd > maxLifecycleDays {
					return NewInvalidArgumentError(fmt.Sprintf("NoncurrentVersionTransition NoncurrentDays must be between 1 and %d, got %d", maxLifecycleDays, nd))
				}
			}
			if t.StorageClass != "" {
				if err := validateStorageClass(t.StorageClass); err != nil {
					return err
				}
			}
		}

		if rule.AbortIncompleteMultipartUpload != nil && rule.AbortIncompleteMultipartUpload.DaysAfterInitiation != nil {
			d := *rule.AbortIncompleteMultipartUpload.DaysAfterInitiation
			if d < 1 || d > maxLifecycleDays {
				return NewInvalidArgumentError(fmt.Sprintf("AbortIncompleteMultipartUpload DaysAfterInitiation must be between 1 and %d, got %d", maxLifecycleDays, d))
			}
		}

		if rule.NoncurrentVersionExpiration != nil && rule.NoncurrentVersionExpiration.NoncurrentDays != nil {
			d := *rule.NoncurrentVersionExpiration.NoncurrentDays
			if d < 1 || d > maxLifecycleDays {
				return NewInvalidArgumentError(fmt.Sprintf("NoncurrentVersionExpiration NoncurrentDays must be between 1 and %d, got %d", maxLifecycleDays, d))
			}
		}
	}
	return nil
}

// validateRestoreDays validates the Days parameter of a RestoreObject
// request (AWS: 1-3653).
func validateRestoreDays(days int) error {
	if days < 1 || days > maxRestoreDays {
		return NewInvalidArgumentError(fmt.Sprintf("Days must be between 1 and %d", maxRestoreDays))
	}
	return nil
}

// === Bucket configuration validators ===============================

const (
	// Replication constraints (AWS spec).
	maxReplicationIDLength = 255
	maxLogTargetGrants     = 100

	// Ownership controls: AWS allows exactly one rule.
	maxOwnershipRules = 1
)

// validateStorageClass checks that the value is either empty (defaults to
// STANDARD) or a recognised S3 storage class enum.
func validateStorageClass(sc string) error {
	if sc == "" {
		return nil
	}
	if !validStorageClasses[sc] {
		return NewInvalidArgumentError(fmt.Sprintf("invalid StorageClass: %s", sc))
	}
	return nil
}

// validatePayer checks that the Payer value is one of the two AWS-allowed
// values for the RequestPaymentConfiguration.
func validatePayer(payer string) error {
	if payer != "BucketOwner" && payer != "Requester" {
		return NewInvalidArgumentError(fmt.Sprintf("invalid Payer: %s (must be BucketOwner or Requester)", payer))
	}
	return nil
}

// arnRegex matches the general AWS ARN format.
var arnRegex = regexp.MustCompile(`^arn:aws[a-zA-Z-]*:iam::[0-9]{12}:role/[\w+=,.@/-]+$`)

// validateIAMRoleARN validates that the string is a well-formed IAM role
// ARN. S3 replication requires an IAM role ARN in the Role field.
func validateIAMRoleARN(role string) error {
	if role == "" {
		return NewInvalidArgumentError("Role is required for replication configuration")
	}
	if !arnRegex.MatchString(role) {
		return NewInvalidArgumentError(fmt.Sprintf("invalid Role ARN: %s (must be a valid IAM role ARN)", role))
	}
	return nil
}

// validateMFADelete checks the MFADelete value. Empty string is allowed
// (not specified); non-empty must be "Enabled" or "Disabled".
func validateMFADelete(mfa string) error {
	if mfa == "" {
		return nil
	}
	if mfa != "Enabled" && mfa != "Disabled" {
		return NewInvalidArgumentError(fmt.Sprintf("invalid MFADelete: %s (must be Enabled or Disabled)", mfa))
	}
	return nil
}

// KMS key identifier patterns.  AWS accepts four forms for KMSMasterKeyID:
//  1. Key ID (UUID):            12345678-1234-1234-1234-123456789012
//  2. Key ARN:                  arn:aws:kms:<region>:<account>:key/<uuid>
//  3. Alias name:               alias/my-key
//  4. Alias ARN:                arn:aws:kms:<region>:<account>:alias/my-key
var (
	kmsKeyUUIDPattern   = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	kmsKeyArnPattern    = regexp.MustCompile(`^arn:aws[a-zA-Z-]*:kms:[a-z0-9-]+:[0-9]{12}:key/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	kmsAliasNamePattern = regexp.MustCompile(`^alias/[a-zA-Z0-9/_-]{1,128}$`)
	kmsAliasArnPattern  = regexp.MustCompile(`^arn:aws[a-zA-Z-]*:kms:[a-z0-9-]+:[0-9]{12}:alias/[a-zA-Z0-9/_-]{1,128}$`)
)

// validateKMSMasterKeyID validates that when the SSE algorithm is aws:kms
// or aws:kms:dsse, the KMSMasterKeyID is a valid KMS key identifier in any
// of the four accepted forms (key ID, key ARN, alias name, alias ARN),
// or empty (meaning use the default KMS key).
func validateKMSMasterKeyID(keyID, algorithm string) error {
	if algorithm != "aws:kms" && algorithm != "aws:kms:dsse" {
		return nil
	}
	if keyID == "" {
		return nil
	}
	for _, re := range []*regexp.Regexp{kmsKeyUUIDPattern, kmsKeyArnPattern, kmsAliasNamePattern, kmsAliasArnPattern} {
		if re.MatchString(keyID) {
			return nil
		}
	}
	return NewInvalidArgumentError(fmt.Sprintf("invalid KMSMasterKeyID: %s (must be a valid KMS key ID, key ARN, alias name, or alias ARN)", keyID))
}

// validateKMSKeyArn validates that the value is a full KMS key ARN.
// UpdateObjectEncryption accepts only the full key ARN; key IDs and aliases
// that PutObject tolerates are rejected here per the S3 contract.
func validateKMSKeyArn(arn string) error {
	if arn == "" {
		return NewInvalidRequestError("Requests that modify an object's encryption type to SSE-KMS require an Amazon Web Services KMS key Amazon Resource Name (ARN). Modify the request to specify a KMS key ARN, and then try again.")
	}
	if kmsKeyArnPattern.MatchString(arn) {
		return nil
	}
	return NewInvalidRequestError("Requests that modify an object's encryption type to SSE-KMS require a valid Amazon Web Services KMS key Amazon Resource Name (ARN). Confirm that you have a correctly formatted KMS key ARN in your request, and then try again.")
}

// validateGranteeType checks the Grantee type for logging grants.
func validateGranteeType(t string) error {
	switch t {
	case "CanonicalUser", "AmazonCustomerByEmail", "Group":
		return nil
	default:
		return NewInvalidArgumentError(fmt.Sprintf("invalid Grantee type: %s (must be CanonicalUser, AmazonCustomerByEmail, or Group)", t))
	}
}

// validateLogPermission checks the Permission value for logging target grants.
func validateLogPermission(p string) error {
	switch p {
	case "FULL_CONTROL", "READ", "WRITE", "READ_ACP", "WRITE_ACP":
		return nil
	default:
		return NewInvalidArgumentError(fmt.Sprintf("invalid logging Permission: %s (must be FULL_CONTROL, READ, WRITE, READ_ACP, or WRITE_ACP)", p))
	}
}

// validateReplicationStatus checks the Status value used in DeleteMarkerReplication.
func validateReplicationStatus(status string) error {
	if status != "Enabled" && status != "Disabled" {
		return NewInvalidArgumentError(fmt.Sprintf("invalid replication Status: %s (must be Enabled or Disabled)", status))
	}
	return nil
}

// validateObjectLockEnabled checks that the ObjectLockEnabled value is "Enabled".
func validateObjectLockEnabled(val string) error {
	if val != "Enabled" {
		return NewInvalidArgumentError(fmt.Sprintf("invalid ObjectLockEnabled: %s (must be Enabled)", val))
	}
	return nil
}

// validateWebsiteConfig validates a website configuration: IndexDocument
// and RedirectAllRequestsTo are mutually exclusive, and when IndexDocument
// is present its Suffix must be non-empty.
func validateWebsiteConfig(config *WebsiteConfigurationInput) error {
	if config == nil {
		return NewInvalidArgumentError("website configuration is required")
	}
	hasIndex := config.IndexDocument != nil
	hasRedirect := config.RedirectAllRequestsTo != nil
	if hasIndex && hasRedirect {
		return NewInvalidArgumentError("IndexDocument and RedirectAllRequestsTo are mutually exclusive")
	}
	if hasRedirect {
		if config.RedirectAllRequestsTo.HostName == "" {
			return NewInvalidArgumentError("RedirectAllRequestsTo.HostName is required")
		}
	}
	if hasIndex {
		if config.IndexDocument.Suffix == "" {
			return NewInvalidArgumentError("IndexDocument.Suffix is required")
		}
	}
	for _, rule := range config.RoutingRules {
		if rule.Redirect == nil {
			return NewInvalidArgumentError("RoutingRule.Redirect is required")
		}
		if rule.Condition != nil {
			if rule.Condition.HTTPErrorCodeReturnedEquals != nil {
				code := *rule.Condition.HTTPErrorCodeReturnedEquals
				n, err := strconv.Atoi(code)
				if err != nil || n < 400 || n > 499 {
					return NewInvalidArgumentError(fmt.Sprintf("HttpErrorCodeReturnedEquals must be a 4xx HTTP code, got: %s", code))
				}
			}
		}
	}
	return nil
}

// validateOwnershipControls checks that exactly one rule is present and
// the ObjectOwnership value is valid.
func validateOwnershipControls(rules []OwnershipControlsRuleInput) error {
	if len(rules) != maxOwnershipRules {
		return NewInvalidArgumentError(fmt.Sprintf("exactly %d ownership rule is required, got %d", maxOwnershipRules, len(rules)))
	}
	switch rules[0].ObjectOwnership {
	case "BucketOwnerEnforced", "BucketOwnerPreferred", "ObjectWriter":
		return nil
	default:
		return NewInvalidArgumentError(fmt.Sprintf("invalid ObjectOwnership: %s (must be BucketOwnerEnforced, BucketOwnerPreferred, or ObjectWriter)", rules[0].ObjectOwnership))
	}
}

// sseCustomerRequested reports whether a request carries any SSE-C customer
// header. AWS requires the algorithm, key, and key MD5 as a set, so the
// presence of any one marks the request as an SSE-C request; incomplete sets
// are rejected when the customer key is parsed.
func sseCustomerRequested(algorithm, key, keyMD5 string) bool {
	return algorithm != "" || key != "" || keyMD5 != ""
}
