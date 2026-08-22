package cognitoidentityprovider

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
)

// validators.go — input validation functions derived from Smithy model traits.
// All validators return bool (true = valid) following the API Gateway pattern.

// validAuthFlows is the set of Smithy AuthFlowType enum values (8 total).
// Reference: cognito-identity-provider-2016-04-18.json shape AuthFlowType.
var validAuthFlows = map[string]bool{
	"USER_SRP_AUTH":            true,
	"REFRESH_TOKEN_AUTH":       true,
	"REFRESH_TOKEN":            true,
	"CUSTOM_AUTH":              true,
	"ADMIN_NO_SRP_AUTH":        true,
	"USER_PASSWORD_AUTH":       true,
	"ADMIN_USER_PASSWORD_AUTH": true,
	"USER_AUTH":                true,
}

// validInitiateAuthFlows is the subset valid for InitiateAuth (client-side).
// ADMIN_NO_SRP_AUTH and ADMIN_USER_PASSWORD_AUTH are AdminInitiateAuth-only.
var validInitiateAuthFlows = map[string]bool{
	"USER_SRP_AUTH":      true,
	"REFRESH_TOKEN_AUTH": true,
	"REFRESH_TOKEN":      true,
	"CUSTOM_AUTH":        true,
	"USER_PASSWORD_AUTH": true,
	"USER_AUTH":          true,
}

// validAdminInitiateAuthFlows is the subset valid for AdminInitiateAuth.
var validAdminInitiateAuthFlows = map[string]bool{
	"ADMIN_NO_SRP_AUTH":        true,
	"ADMIN_USER_PASSWORD_AUTH": true,
	"REFRESH_TOKEN_AUTH":       true,
	"REFRESH_TOKEN":            true,
	"CUSTOM_AUTH":              true,
}

// validChallengeNames is the set of Smithy ChallengeNameType enum values
// (16 total).
var validChallengeNames = map[string]bool{
	"SMS_MFA":                  true,
	"EMAIL_OTP":                true,
	"SOFTWARE_TOKEN_MFA":       true,
	"SELECT_MFA_TYPE":          true,
	"MFA_SETUP":                true,
	"PASSWORD_VERIFIER":        true,
	"CUSTOM_CHALLENGE":         true,
	"SELECT_CHALLENGE":         true,
	"DEVICE_SRP_AUTH":          true,
	"DEVICE_PASSWORD_VERIFIER": true,
	"ADMIN_NO_SRP_AUTH":        true,
	"NEW_PASSWORD_REQUIRED":    true,
	"SMS_OTP":                  true,
	"PASSWORD":                 true,
	"WEB_AUTHN":                true,
	"PASSWORD_SRP":             true,
}

// validProviderTypes is the set of Smithy IdentityProviderTypeType enum
// values (6 total).
var validProviderTypes = map[string]bool{
	"SAML":            true,
	"Facebook":        true,
	"Google":          true,
	"LoginWithAmazon": true,
	"SignInWithApple": true,
	"OIDC":            true,
}

// validAuthFactors is the set of Smithy AuthFactorType enum values (5 total).
// Used by GetUserAuthFactors response (ConfiguredUserAuthFactors field).
var validAuthFactors = map[string]bool{
	"PASSWORD":       true,
	"EMAIL_OTP":      true,
	"SMS_OTP":        true,
	"WEB_AUTHN":      true,
	"SOFTWARE_TOKEN": true,
}

// validDeliveryMediums is the set of Smithy DeliveryMediumType enum values.
var validDeliveryMediums = map[string]bool{
	"SMS":   true,
	"EMAIL": true,
}

// validateInitiateAuthFlow returns true if the auth flow is a recognised
// Smithy AuthFlowType value valid for the client-side InitiateAuth API.
func validateInitiateAuthFlow(flow string) bool {
	return validInitiateAuthFlows[flow]
}

// validateAdminInitiateAuthFlow returns true if the auth flow is a recognised
// Smithy AuthFlowType value valid for the server-side AdminInitiateAuth API.
func validateAdminInitiateAuthFlow(flow string) bool {
	return validAdminInitiateAuthFlows[flow]
}

// validateChallengeName returns true if the challenge name is a recognised
// Smithy ChallengeNameType enum value.
func validateChallengeName(name string) bool {
	return validChallengeNames[name]
}

// validateProviderType returns true if the provider type is a recognised
// Smithy IdentityProviderTypeType enum value.
func validateProviderType(t string) bool {
	return validProviderTypes[t]
}

// validatePasswordHistorySize returns true if the value is within the Smithy
// PasswordHistorySizeType range [0, 24].
func validatePasswordHistorySize(v int) bool {
	return v >= 0 && v <= 24
}

// maxImageFileSize is the maximum allowed ImageFile blob size in bytes.
// Smithy ImageFileType trait: @length(min=0, max=131072).
const maxImageFileSize = 131072

// validateImageFileSize returns true if the decoded blob does not exceed the
// 131072-byte (128 KiB) Smithy length constraint.
func validateImageFileSize(data []byte) bool {
	return len(data) <= maxImageFileSize
}

// ---------------------------------------------------------------------------
// Tag validation — Smithy ArnType, TagKeysType, TagValueType constraints.
// Shared by both the Core layer (tag_core.go, admin handler path) and the
// HTTP tag handler (tag_operations.go, via tagutil callbacks).
// Reference: cognito-identity-provider-2016-04-18.json shapes ArnType
// (@length(min=20, max=2048)), TagKeysType (@length(min=1, max=128)),
// TagValueType (@length(min=0, max=256)).
// ---------------------------------------------------------------------------

// validateCognitoResourceArn validates the ResourceArn parameter against the
// Smithy ArnType length constraint [20, 2048] and the @required trait.
func validateCognitoResourceArn(arn string) error {
	if arn == "" {
		return ErrInvalidParameter
	}
	if len(arn) < 20 || len(arn) > 2048 {
		return awserrors.NewInvalidParameterException(
			fmt.Sprintf("ResourceArn length must be 20-2048: got %d", len(arn)))
	}
	return nil
}

// validateCognitoTagKey validates a single tag key against the Smithy
// TagKeysType length constraint [1, 128] (counted in Unicode characters).
func validateCognitoTagKey(key string) error {
	if n := utf8.RuneCountInString(key); n < 1 || n > 128 {
		return awserrors.NewInvalidParameterException(
			fmt.Sprintf("Tag key length must be 1-128: got %d", n))
	}
	return nil
}

// validateCognitoTagValue validates a single tag value against the Smithy
// TagValueType length constraint [0, 256] (counted in Unicode characters).
func validateCognitoTagValue(value string) error {
	if n := utf8.RuneCountInString(value); n > 256 {
		return awserrors.NewInvalidParameterException(
			fmt.Sprintf("Tag value length must not exceed 256: got %d", n))
	}
	return nil
}

// validateCognitoTags validates a tag map against the Cognito tag limits:
// at most 50 tags per user pool or identity pool, keys of 1-128 characters,
// values of at most 256 characters and the aws: key prefix reserved for AWS
// use.
func validateCognitoTags(tags map[string]string) error {
	switch v, _ := tagutil.CheckStringTags(tags, tagutil.StandardLimits()); v {
	case tagutil.TooManyTags:
		return awserrors.NewInvalidParameterException(
			fmt.Sprintf("Number of tags must not exceed %d", tagutil.MaxTagsPerResource))
	case tagutil.TagKeyTooShort, tagutil.TagKeyTooLong:
		return cognitoTagKeyError(tags)
	case tagutil.TagValueTooLong:
		return cognitoTagValueError(tags)
	case tagutil.ReservedTagKey:
		return awserrors.NewInvalidParameterException(
			"Tag keys cannot start with 'aws:' because the prefix is reserved for AWS use")
	}
	return nil
}

// cognitoTagKeyError reports the first key outside the 1-128 range in the
// validateCognitoTagKey message shape. Keys are walked in sorted order so
// the reported offender is deterministic.
func cognitoTagKeyError(tags map[string]string) error {
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if err := validateCognitoTagKey(k); err != nil {
			return err
		}
	}
	return nil
}

// cognitoTagValueError reports the first value above 256 characters in the
// validateCognitoTagValue message shape. Keys are walked in sorted order so
// the reported offender is deterministic.
func cognitoTagValueError(tags map[string]string) error {
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if err := validateCognitoTagValue(tags[k]); err != nil {
			return err
		}
	}
	return nil
}

// validateCognitoTagKeys validates all keys in a tag-key slice (UntagResource).
func validateCognitoTagKeys(keys []string) error {
	for _, k := range keys {
		if err := validateCognitoTagKey(k); err != nil {
			return err
		}
	}
	return nil
}

// validateCognitoTagsFromTypes validates []tagutil.Tag — the slice form used by
// the HTTP tag handler framework via the ValidateTagsFunc callback. It applies
// the same limits as the map form (count, aws: reservation, key and value
// lengths) so both entry points enforce one contract.
func validateCognitoTagsFromTypes(tagList []tagutil.Tag) error {
	switch v, _ := tagutil.CheckTags(tagList, tagutil.StandardLimits()); v {
	case tagutil.TooManyTags:
		return awserrors.NewInvalidParameterException(
			fmt.Sprintf("Number of tags must not exceed %d", tagutil.MaxTagsPerResource))
	case tagutil.TagKeyTooShort, tagutil.TagKeyTooLong:
		for _, t := range tagList {
			if err := validateCognitoTagKey(t.Key); err != nil {
				return err
			}
		}
		return nil
	case tagutil.TagValueTooLong:
		for _, t := range tagList {
			if err := validateCognitoTagValue(t.Value); err != nil {
				return err
			}
		}
		return nil
	case tagutil.ReservedTagKey:
		return awserrors.NewInvalidParameterException(
			"Tag keys cannot start with 'aws:' because the prefix is reserved for AWS use")
	}
	return nil
}

// ---------------------------------------------------------------------------
// Enum / pattern / range validators for Smithy-constrained input types.
// ---------------------------------------------------------------------------

// validTermsNames is the set of Smithy TermsNameType enum values.
// Pattern: ^(terms-of-use|privacy-policy)$
var validTermsNames = map[string]bool{
	"terms-of-use":   true,
	"privacy-policy": true,
}

// validateTermsName returns true if the value matches the Smithy
// TermsNameType pattern.
func validateTermsName(name string) bool {
	return validTermsNames[name]
}

// validUpdateReplicaStatuses is the set of Smithy UpdateReplicaStatusType
// enum values (wire format is uppercase).
var validUpdateReplicaStatuses = map[string]bool{
	"ACTIVE":   true,
	"INACTIVE": true,
}

// validateUpdateReplicaStatus returns true if the value is a recognised
// Smithy UpdateReplicaStatusType enum value.
func validateUpdateReplicaStatus(status string) bool {
	return validUpdateReplicaStatuses[status]
}

// validAttributeDataTypes is the set of Smithy AttributeDataType enum values
// (wire format is PascalCase: String/Number/DateTime/Boolean).
var validAttributeDataTypes = map[string]bool{
	"String":   true,
	"Number":   true,
	"DateTime": true,
	"Boolean":  true,
}

// validateAttributeDataType returns true if the value is a recognised
// Smithy AttributeDataType enum value.
func validateAttributeDataType(t string) bool {
	return validAttributeDataTypes[t]
}

// validateAccessTokenValidity returns true if the value is within the Smithy
// AccessTokenValidityType range [0, 86400].
func validateAccessTokenValidity(v int) bool {
	return v >= 0 && v <= 86400
}

// validateIdTokenValidity returns true if the value is within the Smithy
// IdTokenValidityType range [0, 86400].
func validateIdTokenValidity(v int) bool {
	return v >= 0 && v <= 86400
}

// validateRefreshTokenValidity returns true if the value is within the Smithy
// RefreshTokenValidityType range [0, 315360000].
func validateRefreshTokenValidity(v int) bool {
	return v >= 0 && v <= 315360000
}

// validatePrecedence returns true if the value satisfies the Smithy
// PrecedenceType minimum constraint (>= 0).
func validatePrecedence(v int) bool {
	return v >= 0
}

// customAttributeNamePattern is the Smithy CustomAttributeNameType pattern:
// ^[\p{L}\p{M}\p{S}\p{N}\p{P}]+$
var customAttributeNamePattern = regexp.MustCompile(`^[\p{L}\p{M}\p{S}\p{N}\p{P}]+$`)

// validateCustomAttributeName validates a custom attribute name against the
// Smithy CustomAttributeNameType length [1, 20] (counted in Unicode
// characters — the pattern admits multibyte categories) and pattern
// constraints.
func validateCustomAttributeName(name string) error {
	if n := utf8.RuneCountInString(name); n < 1 || n > 20 {
		return awserrors.NewInvalidParameterException(
			fmt.Sprintf("Custom attribute name length must be 1-20: got %d", n))
	}
	if !customAttributeNamePattern.MatchString(name) {
		return awserrors.NewInvalidParameterException(
			fmt.Sprintf("Custom attribute name has invalid characters: %s", name))
	}
	return nil
}

// totpCodePattern matches 6-digit TOTP codes per RFC 6238.
var totpCodePattern = regexp.MustCompile(`^[0-9]{6}$`)

// validateMFADeliveryMedium returns true if the value is a recognised Smithy
// DeliveryMediumType enum value (SMS/EMAIL).
func validateMFADeliveryMedium(m string) bool {
	return validDeliveryMediums[m]
}

// validateRegionName validates a region name against the Smithy
// RegionNameType length constraint [5, 32] counted in Unicode characters.
func validateRegionName(name string) error {
	if n := utf8.RuneCountInString(name); n < 5 || n > 32 {
		return awserrors.NewInvalidParameterException(
			fmt.Sprintf("RegionName length must be 5-32: got %d", n))
	}
	return nil
}

// ---------------------------------------------------------------------------
// Smithy enum validators.
// All maps derived from cognito-identity-provider-2016-04-18.json enum shapes.
// ---------------------------------------------------------------------------

// validTimeUnits is the Smithy TimeUnitsType enum (lowercase wire format).
var validTimeUnits = map[string]bool{
	"seconds": true,
	"minutes": true,
	"hours":   true,
	"days":    true,
}

func validateTimeUnit(v string) bool { return validTimeUnits[v] }

// validUserPoolMfaConfigs is the Smithy UserPoolMfaType enum.
var validUserPoolMfaConfigs = map[string]bool{
	"OFF": true, "ON": true, "OPTIONAL": true,
}

func validateUserPoolMfaConfig(v string) bool { return validUserPoolMfaConfigs[v] }

// validDeletionProtections is the Smithy DeletionProtectionType enum.
var validDeletionProtections = map[string]bool{
	"ACTIVE": true, "INACTIVE": true,
}

func validateDeletionProtection(v string) bool { return validDeletionProtections[v] }

// validEmailSendingAccounts is the Smithy EmailSendingAccountType enum.
var validEmailSendingAccounts = map[string]bool{
	"COGNITO_DEFAULT": true, "DEVELOPER": true,
}

func validateEmailSendingAccount(v string) bool { return validEmailSendingAccounts[v] }

// validUserPoolTiers is the Smithy UserPoolTierType enum.
var validUserPoolTiers = map[string]bool{
	"LITE": true, "ESSENTIALS": true, "PLUS": true,
}

func validateUserPoolTier(v string) bool { return validUserPoolTiers[v] }

// validExplicitAuthFlows is the Smithy ExplicitAuthFlowsType enum (9 values).
var validExplicitAuthFlows = map[string]bool{
	"ADMIN_NO_SRP_AUTH":              true,
	"CUSTOM_AUTH_FLOW_ONLY":          true,
	"USER_PASSWORD_AUTH":             true,
	"ALLOW_ADMIN_USER_PASSWORD_AUTH": true,
	"ALLOW_CUSTOM_AUTH":              true,
	"ALLOW_USER_PASSWORD_AUTH":       true,
	"ALLOW_USER_SRP_AUTH":            true,
	"ALLOW_REFRESH_TOKEN_AUTH":       true,
	"ALLOW_USER_AUTH":                true,
}

func validateExplicitAuthFlow(v string) bool { return validExplicitAuthFlows[v] }

// validOAuthFlows is the Smithy OAuthFlowType enum (lowercase wire format).
var validOAuthFlows = map[string]bool{
	"code":               true,
	"implicit":           true,
	"client_credentials": true,
}

func validateOAuthFlow(v string) bool { return validOAuthFlows[v] }

// validSecurityPolicies is the Smithy SecurityPolicyType enum.
var validSecurityPolicies = map[string]bool{
	"TLS_V1":        true,
	"TLS_V1_2_2021": true,
	"TLS_V1_3_2025": true,
}

func validateSecurityPolicy(v string) bool { return validSecurityPolicies[v] }

// validAdvancedSecurityModes is the Smithy AdvancedSecurityModeType enum.
var validAdvancedSecurityModes = map[string]bool{
	"OFF": true, "AUDIT": true, "ENFORCED": true,
}

func validateAdvancedSecurityMode(v string) bool { return validAdvancedSecurityModes[v] }

// validDefaultEmailOptions is the Smithy DefaultEmailOptionType enum.
var validDefaultEmailOptions = map[string]bool{
	"CONFIRM_WITH_LINK": true,
	"CONFIRM_WITH_CODE": true,
}

func validateDefaultEmailOption(v string) bool { return validDefaultEmailOptions[v] }

// validFeatures is the Smithy FeatureType enum.
var validFeatures = map[string]bool{
	"ENABLED": true, "DISABLED": true,
}

func validateFeatureType(v string) bool { return validFeatures[v] }

// validRecoveryOptionNames is the Smithy RecoveryOptionNameType enum
// (lowercase wire format).
var validRecoveryOptionNames = map[string]bool{
	"verified_email":        true,
	"verified_phone_number": true,
	"admin_only":            true,
}

func validateRecoveryOptionName(v string) bool { return validRecoveryOptionNames[v] }

// validEventFilters is the Smithy EventFilterType enum (uppercase wire format).
var validEventFilters = map[string]bool{
	"SIGN_IN":         true,
	"PASSWORD_CHANGE": true,
	"SIGN_UP":         true,
}

func validateEventFilter(v string) bool { return validEventFilters[v] }

// validMessageActions is the Smithy MessageActionType enum.
var validMessageActions = map[string]bool{
	"RESEND":   true,
	"SUPPRESS": true,
}

func validateMessageAction(v string) bool { return validMessageActions[v] }

// validPasswordHashingAlgorithms is the Smithy PasswordHashingAlgorithmType enum.
var validPasswordHashingAlgorithms = map[string]bool{
	"BCRYPT":        true,
	"SCRYPT":        true,
	"ARGON2ID":      true,
	"PBKDF2_SHA256": true,
}

func validatePasswordHashingAlgorithm(v string) bool {
	return validPasswordHashingAlgorithms[v]
}

// validPreventUserExistenceErrors is the Smithy PreventUserExistenceErrorTypes enum.
var validPreventUserExistenceErrors = map[string]bool{
	"LEGACY":  true,
	"ENABLED": true,
}

func validatePreventUserExistenceErrors(v string) bool {
	return validPreventUserExistenceErrors[v]
}

// validAliasAttributes is the Smithy AliasAttributeType enum (lowercase wire format).
var validAliasAttributes = map[string]bool{
	"phone_number":       true,
	"email":              true,
	"preferred_username": true,
}

func validateAliasAttribute(v string) bool { return validAliasAttributes[v] }

// validUsernameAttributes is the Smithy UsernameAttributeType enum (lowercase wire format).
var validUsernameAttributes = map[string]bool{
	"phone_number": true,
	"email":        true,
}

func validateUsernameAttribute(v string) bool { return validUsernameAttributes[v] }

// validVerifiedAttributes is the Smithy VerifiedAttributeType enum (lowercase wire format).
var validVerifiedAttributes = map[string]bool{
	"phone_number": true,
	"email":        true,
}

func validateVerifiedAttribute(v string) bool { return validVerifiedAttributes[v] }

// usernamePattern is the Smithy pattern for UsernameType and GroupNameType:
// ^[\p{L}\p{M}\p{S}\p{N}\p{P}]+$ (length 1-128).
var usernamePattern = regexp.MustCompile(`^[\p{L}\p{M}\p{S}\p{N}\p{P}]+$`)

// validateUsernamePattern returns true if the value matches the Smithy
// pattern and length constraint (1-128, counted in Unicode characters) for
// UsernameType / GroupNameType.
func validateUsernamePattern(v string) bool {
	n := utf8.RuneCountInString(v)
	return n >= 1 && n <= 128 && usernamePattern.MatchString(v)
}

// userPoolNamePattern is the Smithy pattern for UserPoolNameType:
// ^[\w\s+=,.@-]+$ (length 1-128).
var userPoolNamePattern = regexp.MustCompile(`^[\w\s+=,.@-]+$`)

// validateUserPoolNamePattern returns true if the value matches the Smithy
// pattern and length constraint (1-128) for UserPoolNameType.
func validateUserPoolNamePattern(v string) bool {
	return len(v) >= 1 && len(v) <= 128 && userPoolNamePattern.MatchString(v)
}

// listLimitMax is the upper bound shared by the Cognito list-limit shapes
// (QueryLimitType, PoolQueryLimitType, QueryLimit — all Smithy max 60).
const listLimitMax = 60

// parseListLimit extracts a list-limit parameter typed QueryLimitType in the
// Smithy model (range 0-60, e.g. ListUsers.Limit). Absent or zero selects
// defaultValue, matching the AWS behaviour of returning the documented
// default; an explicitly provided value outside 0-60 is rejected.
func parseListLimit(params map[string]interface{}, key string, defaultValue int) (int, error) {
	if err := rejectNonNumericLimit(params, key, "0"); err != nil {
		return 0, err
	}
	v, present := request.GetIntParamCaseInsensitive(params, key)
	if !present || v == 0 {
		return defaultValue, nil
	}
	if v < 0 || v > listLimitMax {
		return 0, awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("%s must be between 0 and %d", key, listLimitMax), http.StatusBadRequest)
	}
	return v, nil
}

// parseStrictListLimit extracts a list-limit parameter whose Smithy shape
// has a minimum of 1 (PoolQueryLimitType / QueryLimit, e.g. ListUserPools,
// ListUserPoolClients). Absent selects defaultValue; an explicitly provided
// value outside 1-60 is rejected, including an explicit zero.
func parseStrictListLimit(params map[string]interface{}, key string, defaultValue int) (int, error) {
	if err := rejectNonNumericLimit(params, key, "1"); err != nil {
		return 0, err
	}
	v, present := request.GetIntParamCaseInsensitive(params, key)
	if !present {
		return defaultValue, nil
	}
	if v < 1 || v > listLimitMax {
		return 0, awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("%s must be between 1 and %d", key, listLimitMax), http.StatusBadRequest)
	}
	return v, nil
}

// rejectNonNumericLimit fails closed when a limit parameter is present but
// cannot be interpreted as an integer. GetIntParamCaseInsensitive cannot
// distinguish that case from an absent parameter, so presence is checked
// against the key variants directly.
func rejectNonNumericLimit(params map[string]interface{}, key, lowerBound string) error {
	for _, k := range []string{key, request.LowerFirst(key), strings.ToLower(key)} {
		if _, ok := params[k]; !ok {
			continue
		}
		if _, isInt := request.GetIntParamCaseInsensitive(params, key); !isInt {
			return awserrors.NewAWSError("InvalidParameterException",
				fmt.Sprintf("%s must be a number between %s and %d", key, lowerBound, listLimitMax), http.StatusBadRequest)
		}
		break
	}
	return nil
}
