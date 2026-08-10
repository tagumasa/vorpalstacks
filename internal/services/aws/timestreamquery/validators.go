package timestreamquery

import (
	"regexp"

	awserrors "vorpalstacks/internal/common/errors"
)

// ----------------------------------------------------------------------------
// Smithy shape constraints (timestream-query-2018-11-01.json)
// ----------------------------------------------------------------------------

// String shapes with length and/or pattern traits.
//
//   ScheduledQueryName   length {1,64}   pattern ^[a-zA-Z0-9|!\-_*'()]([a-zA-Z0-9]|[!\-_*'()/.])+$
//   QueryId              length {1,64}   pattern ^[a-zA-Z0-9]+$
//   ClientToken          length {32,128}
//   ClientRequestToken   length {32,128}
//   QueryString          length {1,262144}
//   ScheduleExpression   length {1,256}
//   AmazonResourceName   length {1,2048}
//   S3BucketName         length {3,63}    pattern ^[a-z0-9][.\-a-z0-9]{1,61}[a-z0-9]$
//   S3ObjectKeyPrefix    length {1,896}   pattern ^[a-zA-Z0-9|!\-_*'()]([a-zA-Z0-9]|[!\-_*'()/.])+$
//   PaginationToken      length {1,2048}
//   TagKey               length {1,128}
//   TagValue             length {0,256}
//
// Integer shapes with range traits.
//
//   MaxQueryResults            range {1,1000}
//   MaxScheduledQueriesResults range {1,1000}
//   MaxTagsForResourceResult   range {1,200}

const (
	maxScheduledQueryName = 64
	maxQueryIDLen         = 64
	minClientTokenLen     = 32
	maxClientTokenLen     = 128
	maxQueryStringLen     = 262144
	maxScheduleExpression = 256
	maxAmazonResourceName = 2048
	minS3BucketName       = 3
	maxS3BucketName       = 63
	maxS3ObjectKeyPrefix  = 896
	maxPaginationToken    = 2048
	maxTagKeyLen          = 128
	maxTagValueLen        = 256

	rangeMaxQueryResults            = 1000
	rangeMaxScheduledQueriesResults = 1000
	rangeMaxTagsForResourceResult   = 200
)

var (
	scheduledQueryNamePattern = regexp.MustCompile(`^[a-zA-Z0-9|!\-_*'()]([a-zA-Z0-9]|[!\-_*'()/.])+$`)
	queryIDPattern            = regexp.MustCompile(`^[a-zA-Z0-9]{1,64}$`)
	s3BucketNamePattern       = regexp.MustCompile(`^[a-z0-9][.\-a-z0-9]{1,61}[a-z0-9]$`)
)

// ---------------------------------------------------------------------------
// Enum validators
// ---------------------------------------------------------------------------

var (
	validScheduledQueryStates = map[string]bool{
		"ENABLED":  true,
		"DISABLED": true,
	}

	validQueryPricingModels = map[string]bool{
		"BYTES_SCANNED": true,
		"COMPUTE_UNITS": true,
	}

	validComputeModes = map[string]bool{
		"ON_DEMAND":   true,
		"PROVISIONED": true,
	}

	validDimensionValueTypes = map[string]bool{
		"VARCHAR": true,
	}

	validMeasureValueTypes = map[string]bool{
		"BIGINT":  true,
		"BOOLEAN": true,
		"DOUBLE":  true,
		"VARCHAR": true,
		"MULTI":   true,
	}

	validScalarMeasureValueTypes = map[string]bool{
		"BIGINT":    true,
		"BOOLEAN":   true,
		"DOUBLE":    true,
		"VARCHAR":   true,
		"TIMESTAMP": true,
	}

	validS3EncryptionOptions = map[string]bool{
		"SSE_S3":  true,
		"SSE_KMS": true,
	}

	validQueryInsightsModes = map[string]bool{
		"ENABLED_WITH_RATE_CONTROL": true,
		"DISABLED":                  true,
	}
)

// ---------------------------------------------------------------------------
// String validators
// ---------------------------------------------------------------------------

// validateScheduledQueryName validates the Name parameter against the Smithy
// ScheduledQueryName shape: length {1,64} and pattern above.
func validateScheduledQueryName(name string) error {
	if len(name) < 1 || len(name) > maxScheduledQueryName {
		return awserrors.NewAWSError("ValidationException",
			"ScheduledQueryName must be between 1 and 64 characters.", 400)
	}
	if !scheduledQueryNamePattern.MatchString(name) {
		return awserrors.NewAWSError("ValidationException",
			"ScheduledQueryName does not match the required pattern.", 400)
	}
	return nil
}

// validateQueryID validates a QueryId against the Smithy QueryId shape:
// length {1,64} and pattern ^[a-zA-Z0-9]+$.
func validateQueryID(id string) error {
	if len(id) < 1 || len(id) > maxQueryIDLen {
		return awserrors.NewAWSError("ValidationException",
			"QueryId must be between 1 and 64 characters.", 400)
	}
	if !queryIDPattern.MatchString(id) {
		return awserrors.NewAWSError("ValidationException",
			"QueryId does not match the required pattern.", 400)
	}
	return nil
}

// validateClientToken validates a ClientToken or ClientRequestToken against
// the Smithy shapes: length {32,128}. Only called when the client explicitly
// provides a token (empty tokens are replaced with a generated UUID).
func validateClientToken(token string) error {
	if len(token) < minClientTokenLen || len(token) > maxClientTokenLen {
		return awserrors.NewAWSError("ValidationException",
			"ClientToken must be between 32 and 128 characters.", 400)
	}
	return nil
}

// validateQueryString validates a QueryString against the Smithy
// QueryString shape: length {1,262144}.
func validateQueryString(qs string) error {
	if len(qs) < 1 || len(qs) > maxQueryStringLen {
		return awserrors.NewAWSError("ValidationException",
			"QueryString must be between 1 and 262144 characters.", 400)
	}
	return nil
}

// validateScheduleExpression validates a ScheduleExpression against the
// Smithy ScheduleExpression shape: length {1,256}.
func validateScheduleExpression(expr string) error {
	if len(expr) < 1 || len(expr) > maxScheduleExpression {
		return awserrors.NewAWSError("ValidationException",
			"ScheduleExpression must be between 1 and 256 characters.", 400)
	}
	return nil
}

// validateAmazonResourceName validates an ARN against the Smithy
// AmazonResourceName shape: length {1,2048}.
func validateAmazonResourceName(arn string) error {
	if len(arn) < 1 || len(arn) > maxAmazonResourceName {
		return awserrors.NewAWSError("ValidationException",
			"ResourceArn must be between 1 and 2048 characters.", 400)
	}
	return nil
}

// validateS3BucketName validates an S3 bucket name against the Smithy
// S3BucketName shape: length {3,63} and pattern.
func validateS3BucketName(name string) error {
	if len(name) < minS3BucketName || len(name) > maxS3BucketName {
		return awserrors.NewAWSError("ValidationException",
			"BucketName must be between 3 and 63 characters.", 400)
	}
	if !s3BucketNamePattern.MatchString(name) {
		return awserrors.NewAWSError("ValidationException",
			"BucketName does not match the required pattern.", 400)
	}
	return nil
}

// validateS3ObjectKeyPrefix validates an S3 object key prefix against the
// Smithy S3ObjectKeyPrefix shape: length {1,896} and pattern.
func validateS3ObjectKeyPrefix(prefix string) error {
	if len(prefix) < 1 || len(prefix) > maxS3ObjectKeyPrefix {
		return awserrors.NewAWSError("ValidationException",
			"ObjectKeyPrefix must be between 1 and 896 characters.", 400)
	}
	if !scheduledQueryNamePattern.MatchString(prefix) {
		return awserrors.NewAWSError("ValidationException",
			"ObjectKeyPrefix does not match the required pattern.", 400)
	}
	return nil
}

// validateTagKey validates a tag key against the Smithy TagKey shape:
// length {1,128}.
func validateTagKey(key string) error {
	if len(key) < 1 || len(key) > maxTagKeyLen {
		return awserrors.NewAWSError("ValidationException",
			"Tag key must be between 1 and 128 characters.", 400)
	}
	return nil
}

// validateTagValue validates a tag value against the Smithy TagValue shape:
// length {0,256}.
func validateTagValue(val string) error {
	if len(val) > maxTagValueLen {
		return awserrors.NewAWSError("ValidationException",
			"Tag value must be between 0 and 256 characters.", 400)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Integer range validators
// ---------------------------------------------------------------------------

// validateMaxResultsQuery validates MaxRows for the Query operation against
// the Smithy MaxQueryResults shape: range {1,1000}.
func validateMaxResultsQuery(val int) error {
	if val < 1 || val > rangeMaxQueryResults {
		return awserrors.NewAWSError("ValidationException",
			"MaxRows must be between 1 and 1000.", 400)
	}
	return nil
}

// validateMaxResultsScheduledQueries validates MaxResults for the
// ListScheduledQueries operation against the Smithy
// MaxScheduledQueriesResults shape: range {1,1000}.
func validateMaxResultsScheduledQueries(val int) error {
	if val < 1 || val > rangeMaxScheduledQueriesResults {
		return awserrors.NewAWSError("ValidationException",
			"MaxResults must be between 1 and 1000.", 400)
	}
	return nil
}

// validateMaxResultsTags validates MaxResults for the ListTagsForResource
// operation against the Smithy MaxTagsForResourceResult shape: range {1,200}.
func validateMaxResultsTags(val int) error {
	if val < 1 || val > rangeMaxTagsForResourceResult {
		return awserrors.NewAWSError("ValidationException",
			"MaxResults must be between 1 and 200.", 400)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Enum validators
// ---------------------------------------------------------------------------

// validateScheduledQueryState validates the State parameter against the
// Smithy ScheduledQueryState enum: ENABLED, DISABLED.
func validateScheduledQueryState(state string) bool {
	return validScheduledQueryStates[state]
}

// validateMeasureValueType validates a MeasureValueType against the Smithy
// MeasureValueType enum.
func validateMeasureValueType(t string) bool {
	return validMeasureValueTypes[t]
}

// validateScalarMeasureValueType validates a ScalarMeasureValueType against
// the Smithy ScalarMeasureValueType enum.
func validateScalarMeasureValueType(t string) bool {
	return validScalarMeasureValueTypes[t]
}

// validateDimensionValueType validates a DimensionValueType against the
// Smithy DimensionValueType enum.
func validateDimensionValueType(t string) bool {
	return validDimensionValueTypes[t]
}

// validateS3EncryptionOption validates an S3EncryptionOption against the
// Smithy S3EncryptionOption enum.
func validateS3EncryptionOption(opt string) bool {
	return validS3EncryptionOptions[opt]
}
