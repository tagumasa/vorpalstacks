package s3

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// === Validation constants ===========================================

const (
	// Bucket name constraints (AWS: 3-63 chars, DNS-compatible).
	minBucketNameLength = 3
	maxBucketNameLength = 63

	// Object key constraints (AWS: 1-1024 bytes).
	maxObjectKeyLength  = 1024
	maxObjectSize       = 5 * 1024 * 1024 * 1024 * 1024 // 5 TiB
	maxSingleUploadSize = 5 * 1024 * 1024 * 1024        // 5 GiB

	// Tag constraints (AWS: max 50 tags, key 128, value 256).
	maxTags           = 50
	maxTagKeyLength   = 128
	maxTagValueLength = 256

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

var (
	bucketNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`)
	ipAddressRegex  = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
)

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

// validateBucketName validates an S3 bucket name per AWS DNS naming rules.
func validateBucketName(name string) error {
	if len(name) < minBucketNameLength || len(name) > maxBucketNameLength {
		return NewInvalidBucketNameError(name)
	}
	if !bucketNameRegex.MatchString(name) {
		return NewInvalidBucketNameError(name)
	}
	if strings.HasPrefix(name, "xn--") {
		return NewInvalidBucketNameError(name)
	}
	if strings.HasSuffix(name, "-s3alias") {
		return NewInvalidBucketNameError(name)
	}
	if strings.HasSuffix(name, "--ol-s3") {
		return NewInvalidBucketNameError(name)
	}
	if strings.HasSuffix(name, ".mrap") {
		return NewInvalidBucketNameError(name)
	}
	if ipAddressRegex.MatchString(name) {
		return NewInvalidBucketNameError(name)
	}
	if strings.Contains(name, "..") {
		return NewInvalidBucketNameError(name)
	}
	if strings.Contains(name, ".-") || strings.Contains(name, "-.") {
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
	if len(tags) > maxTags {
		return NewInvalidArgumentError(fmt.Sprintf("too many tags (maximum %d)", maxTags))
	}
	for _, tag := range tags {
		if len(tag.Key) == 0 {
			return NewInvalidArgumentError("tag key cannot be empty")
		}
		if len(tag.Key) > maxTagKeyLength {
			return NewInvalidArgumentError("tag key cannot exceed 128 characters")
		}
		if len(tag.Value) > maxTagValueLength {
			return NewInvalidArgumentError("tag value cannot exceed 256 characters")
		}
		if strings.HasPrefix(strings.ToLower(tag.Key), "aws:") {
			return NewInvalidArgumentError("tag key cannot start with 'aws:' (reserved prefix)")
		}
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
			if t.StorageClass != "" && !validStorageClasses[t.StorageClass] {
				return NewInvalidArgumentError(fmt.Sprintf("invalid StorageClass: %s", t.StorageClass))
			}
		}

		for _, t := range rule.NoncurrentVersionTransitions {
			if t.NoncurrentDays != nil {
				nd := *t.NoncurrentDays
				if nd < 1 || nd > maxLifecycleDays {
					return NewInvalidArgumentError(fmt.Sprintf("NoncurrentVersionTransition NoncurrentDays must be between 1 and %d, got %d", maxLifecycleDays, nd))
				}
			}
			if t.StorageClass != "" && !validStorageClasses[t.StorageClass] {
				return NewInvalidArgumentError(fmt.Sprintf("invalid NoncurrentVersionTransition StorageClass: %s", t.StorageClass))
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
