package ssm

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	ssmstore "vorpalstacks/internal/store/aws/ssm"
)

// storeErrorMappings maps store-level sentinel errors to SSM AWS errors.
var storeErrorMappings = []awserrors.StoreErrorMapping{
	{Store: ssmstore.ErrParameterNotFound, AWS: ErrParameterNotFound},
	{Store: ssmstore.ErrParameterAlreadyExists, AWS: ErrParameterAlreadyExists},
	{Store: ssmstore.ErrInvalidParameterName, AWS: ErrInvalidParameterName},
	{Store: ssmstore.ErrInvalidParameterValue, AWS: ErrInvalidParameterValue},
	{Store: ssmstore.ErrInvalidParameterType, AWS: ErrInvalidParameterType},
	{Store: ssmstore.ErrInvalidParameterVersion, AWS: ErrInvalidParameterVersion},
	{Store: ssmstore.ErrParameterVersionNotFound, AWS: ErrParameterVersionNotFound},
	{Store: ssmstore.ErrParameterPatternMismatch, AWS: ErrParameterPatternMismatch},
	{Store: ssmstore.ErrInvalidAllowedPattern, AWS: ErrInvalidAllowedPattern},
	{Store: ssmstore.ErrReservedParameterName, AWS: ErrParameterPatternMismatch},
	{Store: ssmstore.ErrHierarchyLevelLimitExceeded, AWS: ErrHierarchyLevelLimitExceeded},
	{Store: ssmstore.ErrInvalidNextToken, AWS: ErrInvalidNextToken},
}

// toSSMError converts a generic error to an *awserrors.AWSError,
// mapping store-level sentinels via storeErrorMappings.
func toSSMError(err error) error {
	if err == nil {
		return nil
	}
	if awsErr, ok := err.(*awserrors.AWSError); ok {
		return awsErr
	}
	mapped := awserrors.MapStoreError(err, storeErrorMappings)
	if mapped != nil {
		return mapped
	}
	return err
}

var (
	// ErrParameterNotFound is returned when the specified parameter does not exist.
	ErrParameterNotFound = awserrors.NewAWSError("ParameterNotFound", "Parameter not found", http.StatusNotFound)
	// ErrParameterAlreadyExists is returned when attempting to create a parameter that already exists.
	ErrParameterAlreadyExists = awserrors.NewAWSError("ParameterAlreadyExists", "Parameter already exists", http.StatusBadRequest)
	// ErrInvalidParameterName is returned when the parameter name is invalid.
	ErrInvalidParameterName = awserrors.NewAWSError("InvalidParameter", "Invalid parameter name", http.StatusBadRequest)
	// ErrInvalidParameterValue is returned when the parameter value is invalid.
	ErrInvalidParameterValue = awserrors.NewAWSError("InvalidParameter", "Invalid parameter value", http.StatusBadRequest)
	// ErrInvalidParameterType is returned when the parameter type is invalid.
	ErrInvalidParameterType = awserrors.NewAWSError("InvalidParameterType", "Invalid parameter type", http.StatusBadRequest)
	// ErrInvalidParameterVersion is returned when the parameter version is invalid.
	ErrInvalidParameterVersion = awserrors.NewAWSError("InvalidParameterVersion", "Invalid parameter version", http.StatusBadRequest)
	// ErrInvalidParameterLabel is returned when the parameter label is invalid.
	ErrInvalidParameterLabel = awserrors.NewAWSError("InvalidParameter", "Invalid parameter label", http.StatusBadRequest)
	// ErrParameterVersionNotFound is returned when the specified parameter version does not exist.
	ErrParameterVersionNotFound = awserrors.NewAWSError("ParameterVersionNotFound", "Parameter version not found", http.StatusBadRequest)
	// ErrParameterPatternMismatch is returned when the parameter name uses a reserved prefix.
	ErrParameterPatternMismatch = awserrors.NewAWSError("ParameterPatternMismatchException", "The parameter name is not valid", http.StatusBadRequest)
	// ErrInvalidAllowedPattern is returned when the AllowedPattern is not a valid regular expression.
	ErrInvalidAllowedPattern = awserrors.NewAWSError("InvalidAllowedPatternException", "The regular expression is invalid", http.StatusBadRequest)
	// ErrInvalidResourceType is returned when a tag operation is given a
	// ResourceType that this implementation does not support.
	ErrInvalidResourceType = awserrors.NewAWSError("InvalidResourceType", "The resource type is not supported", http.StatusBadRequest)
	// ErrInvalidFilterKey is returned when a DescribeParameters/GetParametersByPath
	// filter Key is not in the Smithy ParameterStringFilterKey pattern.
	ErrInvalidFilterKey = awserrors.NewAWSError("InvalidFilterKey", "The filter key is not valid", http.StatusBadRequest)
	// ErrInvalidFilterOption is returned when a filter Option is not one of
	// Equals/BeginsWith/Contains.
	ErrInvalidFilterOption = awserrors.NewAWSError("InvalidFilterOption", "The filter option is not valid", http.StatusBadRequest)
	// ErrInvalidFilterValue is returned when a ParameterFilters value is
	// malformed (not a list of objects, or an entry without Values).
	ErrInvalidFilterValue = awserrors.NewAWSError("InvalidFilterValue", "The filter value is not valid", http.StatusBadRequest)
	// ErrSerializationException is the awsJson1_1 deserialisation error for
	// wire requests that never parse into the modelled shape (for example a
	// non-string entry inside a ParameterFilters Values list); AWS returns it
	// before any operation-level validation runs.
	ErrSerializationException = awserrors.NewAWSError("SerializationException", "The request payload couldn't be deserialized.", http.StatusBadRequest)
	// ErrHierarchyLevelLimitExceeded is returned when a parameter name
	// hierarchy exceeds the maximum depth of fifteen levels.
	ErrHierarchyLevelLimitExceeded = awserrors.NewAWSError("HierarchyLevelLimitExceededException", "A parameter name hierarchy can have a maximum depth of fifteen levels", http.StatusBadRequest)
	// ErrInvalidNextToken is returned when a pagination marker is not valid.
	ErrInvalidNextToken = awserrors.NewAWSError("InvalidNextToken", "The specified token is not valid", http.StatusBadRequest)
	// ErrValidationException is returned when a required request member is
	// missing or a member violates its modelled constraint. AWS rejects null
	// required members with the Smithy ValidationException shape (HTTP 400).
	ErrValidationException = awserrors.NewAWSError("ValidationException", "The request isn't valid. Verify that you entered valid contents for the command and try again.", http.StatusBadRequest)
)
