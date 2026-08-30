package athena

import (
	"fmt"
	"net/http"
	"regexp"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/paramvalidation"
	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
)

// Smithy-conformant validators for Athena API parameters.
// Reference: third_party/api-models-aws/models/athena/service/2017-05-18/athena-2017-05-18.json

// ---------------------------------------------------------------------------
// Generic helpers
// ---------------------------------------------------------------------------

// validateQueryStringSize enforces the Smithy QueryString @length(1, 262144),
// counted in Unicode characters (paramvalidation counts runes).
func validateQueryStringSize(qs string) error {
	return paramvalidation.StringLength("QueryString", qs, 1, maxQueryStringSize,
		func(field string, length, min, max int) error {
			return ErrInvalidRequestException
		})
}

// validateStringLength checks that a string's length falls within [min, max].
// If the value is empty and min == 0 the check passes (optional field).
func validateStringLength(fieldName, value string, min, max int) error {
	return paramvalidation.StringLength(fieldName, value, min, max,
		func(field string, length, min, max int) error {
			return invalidRequestParameter(
				fmt.Sprintf("%s length must be between %d and %d (got %d)", field, min, max, length))
		})
}

// ---------------------------------------------------------------------------
// Smithy string shape: pattern + length validators
// ---------------------------------------------------------------------------

// Smithy WorkGroupName: ^[a-zA-Z0-9._-]{1,128}$
var workGroupNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]{1,128}$`)

// Smithy StatementName: ^[a-zA-Z_][a-zA-Z0-9_@:]{0,255}$
// Note: the Smithy model's pattern says {1,256} which requires a minimum of
// 2 characters, but the length trait says min 1. We use {0,255} so the
// effective range is 1-256 total characters (1 first char + 0-255 rest),
// matching the length constraint exactly.
var statementNamePattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_@:]{0,255}$`)

// Smithy CapacityReservationName: ^[a-zA-Z0-9._-]+$, length [1, 128]
var capacityReservationNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// validateWorkGroupName checks a WorkGroup name against the Smithy pattern.
func validateWorkGroupName(name string) error {
	if !workGroupNamePattern.MatchString(name) {
		return invalidRequestParameter(
			fmt.Sprintf("WorkGroup name %q does not match pattern ^[a-zA-Z0-9._-]{1,128}$", name))
	}
	return nil
}

// validateStatementName checks a prepared-statement name against the Smithy pattern.
func validateStatementName(name string) error {
	if !statementNamePattern.MatchString(name) {
		return invalidRequestParameter(
			fmt.Sprintf("StatementName %q does not match pattern ^[a-zA-Z_][a-zA-Z0-9_@:]{0,255}$", name))
	}
	return nil
}

// validateCapacityReservationName validates the Smithy CapacityReservationName
// shape (pattern ^[a-zA-Z0-9._-]+$, length [1, 128]).
func validateCapacityReservationName(name string) error {
	if err := validateStringLength("Name", name, 1, 128); err != nil {
		return err
	}
	if !capacityReservationNamePattern.MatchString(name) {
		return invalidRequestParameter(
			fmt.Sprintf("Name %q does not match pattern ^[a-zA-Z0-9._-]+$", name))
	}
	return nil
}

// ---------------------------------------------------------------------------
// Smithy string shape: length-only validators
// ---------------------------------------------------------------------------

// validateDescriptionString validates Smithy DescriptionString (length [1, 1024]).
// Used by CreateDataCatalog.Description, CreatePreparedStatement.Description.
func validateDescriptionString(value string) error {
	return validateStringLength("Description", value, 1, 1024)
}

// validateWorkGroupDescriptionString validates Smithy WorkGroupDescriptionString
// (length [0, 1024]). Empty is valid (optional field).
func validateWorkGroupDescriptionString(value string) error {
	return validateStringLength("Description", value, 0, 1024)
}

// validateNamedQueryDescriptionString validates Smithy NamedQueryDescriptionString
// (length [0, 1024]). Empty is valid (optional field).
func validateNamedQueryDescriptionString(value string) error {
	return validateStringLength("Description", value, 0, 1024)
}

// validateDatabaseString validates Smithy DatabaseString (length [1, 255]).
func validateDatabaseString(value string) error {
	return validateStringLength("Database", value, 1, 255)
}

// validateCatalogNameString validates Smithy CatalogNameString (length [1, 256]).
func validateCatalogNameString(value string) error {
	return validateStringLength("Name", value, 1, 256)
}

// validateNameString validates Smithy NameString (length [1, 128]).
// Used by NamedQuery.Name.
func validateNameString(value string) error {
	return validateStringLength("Name", value, 1, 128)
}

// ---------------------------------------------------------------------------
// Smithy tag shape validators
// ---------------------------------------------------------------------------

// validateTags validates a tag map against the Athena tag limits: at most
// 50 tags per resource, keys of 1-128 characters, values of at most 256
// characters and the aws: key prefix reserved for AWS use.
func validateTags(tags map[string]string) error {
	violation, key := tagutil.CheckStringTags(tags, tagutil.StandardLimits())
	switch violation {
	case tagutil.TooManyTags:
		return invalidRequestParameter(
			fmt.Sprintf("Number of tags must not exceed %d", tagutil.MaxTagsPerResource))
	case tagutil.TagKeyTooShort, tagutil.TagKeyTooLong:
		return validateStringLength("TagKey", key, 1, tagutil.MaxTagKeyLength)
	case tagutil.TagValueTooLong:
		return validateStringLength("TagValue", tags[key], 0, tagutil.MaxTagValueLength)
	case tagutil.ReservedTagKey:
		return invalidRequestParameter(
			"Tag keys cannot start with 'aws:' because the prefix is reserved for AWS use")
	}
	return nil
}

// ---------------------------------------------------------------------------
// Smithy string shape: AWS-docs-sourced validators
// ---------------------------------------------------------------------------

// validateClientRequestToken validates the ClientRequestToken parameter
// of StartQueryExecution. AWS docs: "Length Constraints: Minimum length of
// 32. Maximum length of 128."
func validateClientRequestToken(token string) error {
	return validateStringLength("ClientRequestToken", token, 32, 128)
}

// validateAdditionalConfiguration validates Smithy AdditionalConfiguration
// (length [1, 128]). AWS docs: "Length Constraints: Minimum length of 1.
// Maximum length of 128." Note: this differs from DescriptionString which
// allows up to 1024.
func validateAdditionalConfiguration(value string) error {
	return validateStringLength("AdditionalConfiguration", value, 1, 128)
}

// Smithy ExecutionRole: ^arn:aws[a-z\-]*:iam::\d{12}:role/?[a-zA-Z_0-9+=,.@\-_/]+$
// AWS docs: "Length Constraints: Minimum length of 20. Maximum length of 2048."
var executionRolePattern = regexp.MustCompile(`^arn:aws[a-z\-]*:iam::\d{12}:role/?[a-zA-Z_0-9+=,.@\-_/]+$`)

// validateExecutionRole validates the ExecutionRole field (IAM role ARN).
func validateExecutionRole(role string) error {
	if err := validateStringLength("ExecutionRole", role, 20, 2048); err != nil {
		return err
	}
	if !executionRolePattern.MatchString(role) {
		return invalidRequestParameter(
			fmt.Sprintf("ExecutionRole %q does not match IAM role ARN pattern", role))
	}
	return nil
}

// ---------------------------------------------------------------------------
// Smithy integer shape: range validators
// ---------------------------------------------------------------------------

// validateTargetDpus validates Smithy TargetDpusInteger (range [4, ∞)).
// CapacityReservation operations declare only InternalServerException and
// InvalidRequestException, so the violation is an InvalidRequestException.
func validateTargetDpus(dpus int32) error {
	if dpus < 4 {
		return awserrors.NewAWSError("InvalidRequestException",
			fmt.Sprintf("TargetDpus must be at least 4 (got %d)", dpus), http.StatusBadRequest)
	}
	return nil
}

// validateBytesScannedCutoff validates Smithy BytesScannedCutoffValue
// (range [10000000, ∞)). AWS docs confirm: "Valid Range: Minimum value of
// 10000000" with no upper bound. A value of 0 is NOT valid — the correct
// way to remove the cutoff is RemoveBytesScannedCutoffPerQuery in
// WorkGroupConfigurationUpdates.
func validateBytesScannedCutoff(value int64) error {
	if value < 10000000 {
		return invalidRequestParameter(
			fmt.Sprintf("BytesScannedCutoffPerQuery must be at least 10000000 (got %d)", value))
	}
	return nil
}

// ---------------------------------------------------------------------------
// Enum validators
// ---------------------------------------------------------------------------

// validateDataCatalogType validates the Smithy DataCatalogType enum.
// Values: LAMBDA | GLUE | HIVE | FEDERATED
func validateDataCatalogType(catalogType string) error {
	switch catalogType {
	case "LAMBDA", "GLUE", "HIVE", "FEDERATED":
		return nil
	default:
		return invalidRequestParameter(
			fmt.Sprintf("Type %q is not a valid DataCatalogType (LAMBDA, GLUE, HIVE, FEDERATED)", catalogType))
	}
}

// validateWorkGroupState validates the Smithy WorkGroupState enum.
// Values: ENABLED | DISABLED
func validateWorkGroupState(state string) error {
	switch state {
	case "ENABLED", "DISABLED":
		return nil
	default:
		return invalidRequestParameter(
			fmt.Sprintf("State %q is not a valid WorkGroupState (ENABLED, DISABLED)", state))
	}
}

// ---------------------------------------------------------------------------
// Capacity-reservation status validators
// ---------------------------------------------------------------------------

// validateCapacityReservationStatusForCancel ensures only ACTIVE reservations
// can be cancelled, matching the Athena API behaviour.
func validateCapacityReservationStatusForCancel(currentStatus string) error {
	if currentStatus != "ACTIVE" {
		return awserrors.NewAWSError("InvalidRequestException",
			fmt.Sprintf("Capacity reservation in state %s cannot be cancelled; only ACTIVE reservations can be cancelled", currentStatus),
			400)
	}
	return nil
}

// validateCapacityReservationStatusForDelete ensures only CANCELLED reservations
// can be deleted, matching the Athena API behaviour.
func validateCapacityReservationStatusForDelete(currentStatus string) error {
	if currentStatus != "CANCELLED" {
		return awserrors.NewAWSError("InvalidRequestException",
			fmt.Sprintf("Capacity reservation in state %s cannot be deleted; only CANCELLED reservations can be deleted", currentStatus),
			400)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Pagination validators
// ---------------------------------------------------------------------------

// Smithy MaxWorkGroupsCount range [1, 50]. Shared by HTTP and admin handlers.
const athenaMaxWorkGroupsResults = 50

// validateMaxResults extracts and validates a pagination limit against the
// Smithy per-operation range. If the parameter is absent the default is
// returned. A value of 0 is treated as "use service default" to prevent
// infinite pagination loops (MaxResults=0 would produce an empty page
// with NextToken pointing to the same offset). Out-of-range values produce
// InvalidRequestException, the error every operation declares.
func validateMaxResults(params map[string]interface{}, defaultVal, minVal, maxVal int) (int, error) {
	val, ok := request.GetIntParamCaseInsensitive(params, "MaxResults")
	if !ok {
		return defaultVal, nil
	}
	if val < minVal || val > maxVal {
		return 0, invalidRequestParameter(
			fmt.Sprintf("MaxResults must be between %d and %d (got %d)", minVal, maxVal, val))
	}
	if val == 0 {
		return defaultVal, nil
	}
	return val, nil
}

// resolveMaxResults applies the same window semantics as validateMaxResults
// to a pre-extracted MaxResults value on the Core layer: absent → default,
// out-of-window → InvalidRequestException, explicit zero → default (the
// zero-as-default rule prevents empty pages with self-pointing tokens).
func resolveMaxResults(val int, present bool, defaultVal, minVal, maxVal int) (int, error) {
	if !present {
		return defaultVal, nil
	}
	if val < minVal || val > maxVal {
		return 0, invalidRequestParameter(
			fmt.Sprintf("MaxResults must be between %d and %d (got %d)", minVal, maxVal, val))
	}
	if val == 0 {
		return defaultVal, nil
	}
	return val, nil
}

// clampMaxResults normalises a MaxResults integer value: non-positive values
// receive the default, values above maxVal are clamped to maxVal. Used by
// the admin gRPC-Web handler where protobuf fields arrive as int64 (not
// map[string]interface{}), so validateMaxResults is not applicable.
func clampMaxResults(value, defaultVal, maxVal int) int {
	return pagination.ClampMaxItems(value, defaultVal, maxVal)
}
