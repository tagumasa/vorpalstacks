package timestreamquery

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"vorpalstacks/internal/common/bucketname"
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
// the Smithy shapes: length {32,128} counted in Unicode characters (the
// shape carries no pattern). Only called when the client explicitly
// provides a token (empty tokens are replaced with a generated UUID).
func validateClientToken(token string) error {
	if n := utf8.RuneCountInString(token); n < minClientTokenLen || n > maxClientTokenLen {
		return awserrors.NewAWSError("ValidationException",
			"ClientToken must be between 32 and 128 characters.", 400)
	}
	return nil
}

// validateQueryString validates a QueryString against the Smithy
// QueryString shape: length {1,262144} counted in Unicode characters (the
// shape carries no pattern, so multibyte SQL literals are valid input).
func validateQueryString(qs string) error {
	if n := utf8.RuneCountInString(qs); n < 1 || n > maxQueryStringLen {
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
// AmazonResourceName shape: length {1,2048} counted in Unicode characters
// (the shape carries no pattern).
func validateAmazonResourceName(arn string) error {
	if n := utf8.RuneCountInString(arn); n < 1 || n > maxAmazonResourceName {
		return awserrors.NewAWSError("ValidationException",
			"ResourceArn must be between 1 and 2048 characters.", 400)
	}
	return nil
}

// validateS3BucketName validates an S3 bucket name against the AWS
// general-purpose bucket naming rules (the Smithy S3BucketName shape's
// documented form): length 3-63, charset, adjacency, IP-form and
// reserved prefix/suffix restrictions.
func validateS3BucketName(name string) error {
	if !bucketname.Validate(name) {
		return awserrors.NewAWSError("ValidationException",
			"BucketName must be a valid S3 bucket name: 3-63 characters of lowercase letters, numbers, hyphens and periods, starting and ending with a letter or number, without adjacent dots, IP-address form or reserved prefixes and suffixes.", 400)
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
// length {1,128} counted in Unicode characters (the shape carries no
// pattern, unlike the general AWS tag guidance).
func validateTagKey(key string) error {
	if n := utf8.RuneCountInString(key); n < 1 || n > maxTagKeyLen {
		return awserrors.NewAWSError("ValidationException",
			"Tag key must be between 1 and 128 characters.", 400)
	}
	return nil
}

// validateTagValue validates a tag value against the Smithy TagValue shape:
// length {0,256} counted in Unicode characters (no pattern).
func validateTagValue(val string) error {
	if utf8.RuneCountInString(val) > maxTagValueLen {
		return awserrors.NewAWSError("ValidationException",
			"Tag value must be between 0 and 256 characters.", 400)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Integer range validators
// ---------------------------------------------------------------------------

// validateMaxResultsInRange validates a strict page-size parameter against
// a Smithy @range upper bound: zero and out-of-range values are rejected
// (no default-on-zero). paramName names the parameter in the error
// message; the bound comes from the operation's Smithy limit constant.
func validateMaxResultsInRange(val int, paramName string, max int) error {
	if val < 1 || val > max {
		return awserrors.NewAWSError("ValidationException",
			fmt.Sprintf("%s must be between 1 and %d.", paramName, max), 400)
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
