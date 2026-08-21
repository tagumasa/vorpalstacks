package secretsmanager

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/common/errors"
	tagutil "vorpalstacks/internal/common/tags"

	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

// Smithy-derived constraints (single source of truth).
// Source: secrets-manager-2017-10-17.json Smithy model.
const (
	maxSecretNameLength        = 512                        // NameType @length(min=1, max=512)
	maxDescriptionLength       = 2048                       // DescriptionType @length(min=0, max=2048)
	minClientRequestTokenLen   = 32                         // ClientRequestTokenType @length(min=32, max=64)
	maxClientRequestTokenLen   = 64                         // ClientRequestTokenType @length(max=64)
	maxSecretValueBytes        = 65536                      // SecretStringType / SecretBinaryType @length(min=1, max=65536)
	minRotationTokenLen        = 36                         // RotationTokenType @length(min=36, max=256)
	maxRotationTokenLen        = 256                        // RotationTokenType @length(max=256)
	minAutomaticallyAfterDays  = 1                          // AutomaticallyRotateAfterDaysType @range(min=1, max=1000)
	maxAutomaticallyAfterDays  = 1000                       // AutomaticallyRotateAfterDaysType @range(max=1000)
	maxResourcePolicyBytes     = 20480                      // NonEmptyResourcePolicyType @length(min=1, max=20480)
	maxFilters                 = 10                         // FiltersListType @length(min=0, max=10)
	maxSecretIdLength          = 2048                       // SecretIdType @length(min=1, max=2048)
	maxSecretIdListItems       = 20                         // SecretIdListType @length(min=1, max=20)
	maxTagsPerSecret           = tagutil.MaxTagsPerResource // AWS tag quota
	maxTagKeyLength            = tagutil.MaxTagKeyLength    // TagKeyType @length(min=1, max=128)
	maxTagValueLength          = tagutil.MaxTagValueLength  // TagValueType @length(min=0, max=256)
	maxKmsKeyIdLength          = 2048                       // KmsKeyIdType @length(min=0, max=2048)
	maxRotationLambdaARNLength = 2048                       // RotationLambdaARNType @length(min=0, max=2048)
	maxExcludeCharactersLength = 4096                       // ExcludeCharactersType @length(min=0, max=4096)
)

var rotationTokenPattern = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)

var (
	durationPattern     = regexp.MustCompile(`^[0-9]+h$`)
	scheduleExprPattern = regexp.MustCompile(`^[0-9A-Za-z()#?*/, -]+$`)
	regionPattern       = regexp.MustCompile(`^([a-z]+-)+[0-9]+$`)
)

// validateSecretName validates the secret name against the Smithy
// NameType @length(1-512) constraint.
func validateSecretName(name string) error {
	if name == "" {
		return awserrors.ErrMissingParameter
	}
	if len(name) > maxSecretNameLength {
		return awserrors.NewAWSError("InvalidParameterException",
			"Secret name must be between 1 and 512 characters long.", http.StatusBadRequest)
	}
	return nil
}

// validateClientRequestToken validates the ClientRequestToken against the
// Smithy ClientRequestTokenType @length(32-64) constraint.
func validateClientRequestToken(token string) error {
	if token == "" {
		return nil
	}
	if len(token) < minClientRequestTokenLen || len(token) > maxClientRequestTokenLen {
		return awserrors.NewAWSError("InvalidParameterException",
			"ClientRequestToken must be 32 to 64 characters long.", http.StatusBadRequest)
	}
	return nil
}

// validateDescription validates the Description against the Smithy
// DescriptionType @length(0-2048) constraint.
func validateDescription(desc string) error {
	if len(desc) > maxDescriptionLength {
		return awserrors.NewAWSError("InvalidParameterException",
			"Description must be between 0 and 2048 characters long.", http.StatusBadRequest)
	}
	return nil
}

// validateSecretStringLength validates the SecretString against the Smithy
// SecretStringType @length(1-65536) constraint.
func validateSecretStringLength(s string) error {
	if len(s) > maxSecretValueBytes {
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("SecretString must not exceed %d bytes.", maxSecretValueBytes), http.StatusBadRequest)
	}
	return nil
}

// validateSecretBinaryLength validates the decoded SecretBinary against the
// Smithy SecretBinaryType @length(1-65536) constraint.
func validateSecretBinaryLength(b []byte) error {
	if len(b) > maxSecretValueBytes {
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("SecretBinary must not exceed %d bytes.", maxSecretValueBytes), http.StatusBadRequest)
	}
	return nil
}

// validateSecretValueMutex enforces the AWS rule that SecretString and
// SecretBinary are mutually exclusive in a single request.
func validateSecretValueMutex(secretString, secretBinaryStr string) error {
	if secretString != "" && secretBinaryStr != "" {
		return awserrors.NewAWSError("InvalidParameterException",
			"You can't specify both SecretString and SecretBinary in the same request.", http.StatusBadRequest)
	}
	return nil
}

// validateRecoveryWindow validates the RecoveryWindowInDays parameter
// against the AWS constraint of 7-30 days (inclusive).
func validateRecoveryWindow(days int, hasWindow bool, forceDelete bool) error {
	if forceDelete && hasWindow {
		return awserrors.NewAWSError("InvalidParameterException",
			"You can't use ForceDeleteWithoutRecovery in conjunction with RecoveryWindowInDays.", http.StatusBadRequest)
	}
	if hasWindow && !forceDelete {
		if days < 7 || days > 30 {
			return awserrors.NewAWSError("InvalidParameterException",
				"RecoveryWindowInDays must be between 7 and 30 days.", http.StatusBadRequest)
		}
	}
	return nil
}

// validateRotationToken validates the RotationToken against the Smithy
// RotationTokenType @length(36-256) and @pattern("^[a-zA-Z0-9\-]+$")
// constraints.
func validateRotationToken(token string) error {
	if token == "" {
		return nil
	}
	if len(token) < minRotationTokenLen || len(token) > maxRotationTokenLen {
		return awserrors.NewAWSError("InvalidParameterException",
			"RotationToken must be between 36 and 256 characters long.", http.StatusBadRequest)
	}
	if !rotationTokenPattern.MatchString(token) {
		return awserrors.NewAWSError("InvalidParameterException",
			"RotationToken must contain only alphanumeric characters and hyphens.", http.StatusBadRequest)
	}
	return nil
}

// validateAutomaticallyAfterDays validates the RotationRules
// AutomaticallyAfterDays against the Smithy
// AutomaticallyRotateAfterDaysType @range(1-1000) constraint.
func validateAutomaticallyAfterDays(days int) error {
	if days < minAutomaticallyAfterDays || days > maxAutomaticallyAfterDays {
		return awserrors.NewAWSError("InvalidParameterException",
			"AutomaticallyAfterDays must be between 1 and 1000.", http.StatusBadRequest)
	}
	return nil
}

// validateResourcePolicyLength validates the resource policy JSON string
// against the Smithy NonEmptyResourcePolicyType @length(1-20480) constraint.
func validateResourcePolicyLength(policy string) error {
	if len(policy) > maxResourcePolicyBytes {
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("ResourcePolicy must not exceed %d bytes.", maxResourcePolicyBytes), http.StatusBadRequest)
	}
	return nil
}

// validateMaxFilters validates the Filters list against the Smithy
// FiltersListType @length(max=10) constraint.
func validateMaxFilters(n int) error {
	if n > maxFilters {
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("You can include up to %d filters.", maxFilters), http.StatusBadRequest)
	}
	return nil
}

// validateSecretIdList validates each SecretId in a SecretIdList against
// the Smithy SecretIdType @length(1-2048) constraint and the list against
// SecretIdListType @length(max=20).
func validateSecretIdList(ids []string) error {
	if len(ids) > maxSecretIdListItems {
		return awserrors.NewAWSError("ValidationException",
			fmt.Sprintf("You can include up to %d secrets in a batch.", maxSecretIdListItems), http.StatusBadRequest)
	}
	for _, id := range ids {
		if len(id) > maxSecretIdLength {
			return awserrors.NewAWSError("InvalidParameterException",
				fmt.Sprintf("SecretId must not exceed %d characters.", maxSecretIdLength), http.StatusBadRequest)
		}
	}
	return nil
}

// validateSecretTags validates tag count, key length, and value length
// against AWS Secrets Manager quotas.  Tag count overflow uses
// InvalidParameterException — it is documented for both CreateSecret and
// TagResource, whereas LimitExceededException is only documented for
// CreateSecret.
func validateSecretTags(tags []tagutil.Tag) error {
	switch v, _ := tagutil.CheckTags(tags, tagutil.StandardLimits()); v {
	case tagutil.TooManyTags:
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("You can't have more than %d tags on a secret.", tagutil.MaxTagsPerResource), http.StatusBadRequest)
	case tagutil.TagKeyTooShort, tagutil.TagKeyTooLong:
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("Tag key length must be between 1 and %d characters.", tagutil.MaxTagKeyLength), http.StatusBadRequest)
	case tagutil.TagValueTooLong:
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("Tag value length must be between 0 and %d characters.", tagutil.MaxTagValueLength), http.StatusBadRequest)
	case tagutil.ReservedTagKey:
		return awserrors.NewAWSError("InvalidParameterException",
			"Tag keys cannot start with 'aws:' because the prefix is reserved for AWS use.", http.StatusBadRequest)
	}
	return nil
}

// validatePasswordRequirements checks that the number of required character
// types does not exceed the requested password length.  AWS rejects this
// combination with InvalidParameterException because the password cannot
// satisfy the requirement.
func validatePasswordRequirements(requiredCount, length int) error {
	if requiredCount > length {
		return awserrors.NewAWSError("InvalidParameterException",
			"PasswordLength is too short to include all required character types.", http.StatusBadRequest)
	}
	return nil
}

// checkNotDeleted rejects write operations on secrets that are scheduled
// for deletion.  AWS returns InvalidRequestException: "The secret is
// scheduled for deletion."
func checkNotDeleted(secret *secretsmanagerstore.Secret) error {
	if secret.DeletedDate != nil {
		return awserrors.NewAWSError("InvalidRequestException",
			"The secret is scheduled for deletion.", http.StatusBadRequest)
	}
	return nil
}

// decodeAndValidateSecretBinary decodes a base64-encoded SecretBinary
// string and validates the decoded length against the Smithy
// SecretBinaryType @length(max=65536) constraint.
func decodeAndValidateSecretBinary(secretBinaryStr string) ([]byte, error) {
	if secretBinaryStr == "" {
		return nil, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(secretBinaryStr)
	if err != nil {
		return nil, awserrors.NewValidationException(fmt.Sprintf("invalid SecretBinary encoding: %v", err))
	}
	if err := validateSecretBinaryLength(decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

// validateKmsKeyId validates the KmsKeyId against the Smithy
// KmsKeyIdType @length(max=2048) constraint.  KmsKeyId may be an ARN,
// key ID, or alias — the Smithy model defines no pattern, only length.
func validateKmsKeyId(id string) error {
	if id == "" {
		return nil
	}
	if len(id) > maxKmsKeyIdLength {
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("KmsKeyId must not exceed %d characters.", maxKmsKeyIdLength), http.StatusBadRequest)
	}
	return nil
}

// validateRotationLambdaARN validates the RotationLambdaARN against the
// Smithy RotationLambdaARNType @length(max=2048) constraint.
func validateRotationLambdaARN(arn string) error {
	if arn == "" {
		return nil
	}
	if len(arn) > maxRotationLambdaARNLength {
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("RotationLambdaARN must not exceed %d characters.", maxRotationLambdaARNLength), http.StatusBadRequest)
	}
	return nil
}

// validateExcludeCharacters validates ExcludeCharacters against the Smithy
// ExcludeCharactersType @length(max=4096) constraint.
func validateExcludeCharacters(s string) error {
	if len(s) > maxExcludeCharactersLength {
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("ExcludeCharacters must not exceed %d characters.", maxExcludeCharactersLength), http.StatusBadRequest)
	}
	return nil
}

// validateSortOrder validates SortOrder against the Smithy SortOrderType
// enum (asc, desc).  Empty string is valid (defaults to asc).
func validateSortOrder(order string) error {
	switch order {
	case "", "asc", "desc":
		return nil
	default:
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("SortOrder must be 'asc' or 'desc', got '%s'.", order), http.StatusBadRequest)
	}
}

// validateSortBy validates SortBy against the Smithy SortByType enum
// (name, created-date, last-accessed-date, last-changed-date).  Empty
// string is valid (defaults to name).
func validateSortBy(by string) error {
	switch by {
	case "", "name", "created-date", "last-accessed-date", "last-changed-date":
		return nil
	default:
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("SortBy must be one of 'name', 'created-date', 'last-accessed-date', 'last-changed-date', got '%s'.", by), http.StatusBadRequest)
	}
}

// validateUntagKeys validates each tag key against the Smithy
// TagKeyType @length(min=1, max=128) constraint.
func validateUntagKeys(keys []string) error {
	for _, k := range keys {
		if k == "" || utf8.RuneCountInString(k) > maxTagKeyLength {
			return awserrors.NewAWSError("InvalidParameterException",
				fmt.Sprintf("Tag key length must be between 1 and %d characters.", maxTagKeyLength), http.StatusBadRequest)
		}
	}
	return nil
}

// validateDuration validates the RotationRules Duration against the Smithy
// DurationType @length(min=2, max=3) and @pattern("^[0-9]+h$") constraints.
func validateDuration(d string) error {
	if d == "" {
		return nil
	}
	if len(d) < 2 || len(d) > 3 {
		return awserrors.NewAWSError("InvalidParameterException",
			"Duration must be between 2 and 3 characters long (e.g. '24h').", http.StatusBadRequest)
	}
	if !durationPattern.MatchString(d) {
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("Duration must match the format '<number>h' (e.g. '24h'), got '%s'.", d), http.StatusBadRequest)
	}
	return nil
}

// validateScheduleExpression validates the RotationRules ScheduleExpression
// against the Smithy ScheduleExpressionType @length(min=1, max=256) and
// @pattern constraints.
func validateScheduleExpression(expr string) error {
	if expr == "" {
		return nil
	}
	if len(expr) > 256 {
		return awserrors.NewAWSError("InvalidParameterException",
			"ScheduleExpression must not exceed 256 characters.", http.StatusBadRequest)
	}
	if !scheduleExprPattern.MatchString(expr) {
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("ScheduleExpression contains invalid characters: '%s'.", expr), http.StatusBadRequest)
	}
	return nil
}

// validateRotationRules enforces the complete RotationRulesType contract.
// Per the Smithy model documentation, the rotation schedule can be set with
// AutomaticallyAfterDays or ScheduleExpression, but not both. Each field is
// additionally validated against its own Smithy type constraints.
func validateRotationRules(automaticallyAfterDays int, scheduleExpression, duration string) error {
	if automaticallyAfterDays > 0 {
		if err := validateAutomaticallyAfterDays(automaticallyAfterDays); err != nil {
			return err
		}
	}
	if err := validateScheduleExpression(scheduleExpression); err != nil {
		return err
	}
	if err := validateDuration(duration); err != nil {
		return err
	}
	if automaticallyAfterDays > 0 && scheduleExpression != "" {
		return awserrors.NewAWSError("InvalidParameterException",
			"You can set the rotation schedule with RotationRules.AutomaticallyAfterDays or RotationRules.ScheduleExpression, but not both.",
			http.StatusBadRequest)
	}
	return nil
}

// validateRegion validates a region code against the Smithy RegionType
// @length(min=1, max=128) and @pattern("^([a-z]+-)+[0-9]+$") constraints.
func validateRegion(region string) error {
	if region == "" {
		return nil
	}
	if len(region) > 128 {
		return awserrors.NewAWSError("InvalidParameterException",
			"Region must not exceed 128 characters.", http.StatusBadRequest)
	}
	if !regionPattern.MatchString(region) {
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("Region '%s' does not match the expected format (e.g. 'us-east-1').", region), http.StatusBadRequest)
	}
	return nil
}

// validateSecretId validates the SecretId against the Smithy
// SecretIdType @length(min=1, max=2048) constraint.
func validateSecretId(id string) error {
	if id == "" {
		return awserrors.ErrMissingParameter
	}
	if len(id) > maxSecretIdLength {
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("SecretId must not exceed %d characters.", maxSecretIdLength), http.StatusBadRequest)
	}
	return nil
}
