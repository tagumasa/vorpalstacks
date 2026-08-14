package timestreamwrite

import (
	"regexp"
	"strings"

	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// Smithy model pagination limits (timestream-write-2018-11-01.json):
// PaginationLimit range {min:1, max:20} for ListDatabases/ListTables.
// PageLimit range {min:1, max:100} for ListBatchLoadTasks.
const (
	maxListDatabasesResults      = 20
	maxListTablesResults         = 20
	maxListBatchLoadTasksResults = 100
)

// maxBatchLoadObjectSize caps the size of a single S3 object read during batch
// load processing. AWS Timestream BatchLoadTask processes CSV files from S3;
// capping the read prevents OOM on oversized objects.
const maxBatchLoadObjectSize = 256 * 1024 * 1024 // 256 MiB

// timestreamNameRegex validates DatabaseName and TableName per AWS docs:
// "Length Constraints: Minimum length of 3. Maximum length of 64."
// Pattern: must start with a letter, followed by letters, digits, underscores,
// periods, or hyphens.
var timestreamNameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.-]{2,63}$`)

// batchLoadTaskIdRegex validates BatchLoadTaskId per Smithy: ^[A-Z0-9]+$ (3-32 chars).
var batchLoadTaskIdRegex = regexp.MustCompile(`^[A-Z0-9]{3,32}$`)

// s3BucketNameRegex validates S3 bucket names per AWS S3 naming rules:
// 3-63 characters, lowercase letters, digits, hyphens, and periods.
var s3BucketNameRegex = regexp.MustCompile(`^[a-z0-9][\.\-a-z0-9]{1,61}[a-z0-9]$`)

// kmsKeyArnRegex validates KMS key ARN format:
// arn:aws:kms:<region>:<account-id>:key/<uuid>
var kmsKeyArnRegex = regexp.MustCompile(`^arn:aws(-cn|-us-gov)?:kms:[a-z0-9-]+:\d{12}:key/[0-9a-f-]+$`)

func isValidTimestreamName(name string) bool {
	return timestreamNameRegex.MatchString(name)
}

// validateKmsKeyId checks that KmsKeyId is within Smithy StringValue2048 bounds
// (1-2048) and, when it looks like an ARN, validates the ARN format.
func validateKmsKeyId(kmsKeyId string) bool {
	if len(kmsKeyId) < 1 || len(kmsKeyId) > 2048 {
		return false
	}
	if strings.HasPrefix(kmsKeyId, "arn:") {
		return kmsKeyArnRegex.MatchString(kmsKeyId)
	}
	return true
}

func validateClientToken(token string) bool {
	return len(token) >= 1 && len(token) <= 64
}

func validateS3BucketName(bucket string) bool {
	if len(bucket) < 3 || len(bucket) > 63 {
		return false
	}
	return s3BucketNameRegex.MatchString(bucket)
}

// s3ObjectKeyPattern is shared by the S3ObjectKey and S3ObjectKeyPrefix
// Smithy shapes, which carry the identical pattern trait
// ^[a-zA-Z0-9|!\-_*'\(\\)]([a-zA-Z0-9]|[!\-_*'\(\\)/.])+$; they differ
// only in their length bounds (1-1024 vs 1-928).
var s3ObjectKeyRegex = regexp.MustCompile(`^[a-zA-Z0-9|!\-_*'\(\\)]([a-zA-Z0-9]|[!\-_*'\(\\)/.])+$`)

// maxS3ObjectKeyLen and maxS3ObjectKeyPrefixLen bound the two S3 key
// shapes per the Smithy length traits (S3ObjectKey 1-1024,
// S3ObjectKeyPrefix 1-928).
const (
	maxS3ObjectKeyLen       = 1024
	maxS3ObjectKeyPrefixLen = 928
)

// validateS3ObjectKey validates a DataModelS3Configuration.ObjectKey value
// against the Smithy S3ObjectKey shape (length 1-1024 plus pattern).
func validateS3ObjectKey(key string) bool {
	if len(key) < 1 || len(key) > maxS3ObjectKeyLen {
		return false
	}
	return s3ObjectKeyRegex.MatchString(key)
}

// validateS3ObjectKeyPrefix validates an ObjectKeyPrefix value against the
// Smithy S3ObjectKeyPrefix shape: the same pattern as S3ObjectKey with
// length 1-928.
func validateS3ObjectKeyPrefix(prefix string) bool {
	if len(prefix) < 1 || len(prefix) > maxS3ObjectKeyPrefixLen {
		return false
	}
	return s3ObjectKeyRegex.MatchString(prefix)
}

func validateEncryptionOption(option string) bool {
	switch option {
	case "SSE_S3", "SSE_KMS":
		return true
	default:
		return false
	}
}

func validatePartitionKeyType(t string) bool {
	switch t {
	case string(tsstore.PartitionKeyTypeDimension),
		string(tsstore.PartitionKeyTypeMeasure):
		return true
	default:
		return false
	}
}

func validateEnforcementInRecord(e string) bool {
	switch e {
	case string(tsstore.EnforcementInRecordRequired),
		string(tsstore.EnforcementInRecordOptional):
		return true
	default:
		return false
	}
}

func validateMeasureValueType(t string) bool {
	switch t {
	case string(tsstore.MeasureValueTypeDouble),
		string(tsstore.MeasureValueTypeBigint),
		string(tsstore.MeasureValueTypeVarchar),
		string(tsstore.MeasureValueTypeBoolean),
		string(tsstore.MeasureValueTypeTimestamp),
		string(tsstore.MeasureValueTypeMulti):
		return true
	default:
		return false
	}
}

func validateTimeUnit(u string) bool {
	switch u {
	case string(tsstore.TimeUnitMilliseconds),
		string(tsstore.TimeUnitSeconds),
		string(tsstore.TimeUnitMicroseconds),
		string(tsstore.TimeUnitNanoseconds):
		return true
	default:
		return false
	}
}

func validateDimensionValueType(t string) bool {
	switch t {
	case string(tsstore.DimensionValueTypeVarchar):
		return true
	default:
		return false
	}
}
