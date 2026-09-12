package cognitoidentityprovider

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/common/errors"
	tagutil "vorpalstacks/internal/common/tags"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// applyListLimitDefaults clamps an already-parsed list limit to the
// service's defaulting rule: an omitted or out-of-range value lists at the
// maximum page size.
func applyListLimitDefaults(maxResults int) int {
	if maxResults <= 0 || maxResults > listLimitMax {
		return listLimitMax
	}
	return maxResults
}

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
// Applied by the tag Cores (tag_core.go) on both the admin and the HTTP
// plane. The ArnType length and pattern constraints are enforced by
// validateArnType further down in this file.
// Reference: cognito-identity-provider-2016-04-18.json shapes ArnType
// (@length(min=20, max=2048)), TagKeysType (@length(min=1, max=128)),
// TagValueType (@length(min=0, max=256)).
// ---------------------------------------------------------------------------

// minTagKeyLength is the TagKeysType length minimum; the tag-shape maxima
// live in the shared cross-service tag limits, which the validators below
// reference rather than duplicating.
const minTagKeyLength = 1

// validateCognitoTagKey validates a single tag key against the Smithy
// TagKeysType length constraint [1, 128] (counted in Unicode characters).
func validateCognitoTagKey(key string) error {
	if n := utf8.RuneCountInString(key); n < minTagKeyLength || n > tagutil.MaxTagKeyLength {
		return awserrors.NewInvalidParameterException(
			fmt.Sprintf("Tag key length must be %d-%d: got %d", minTagKeyLength, tagutil.MaxTagKeyLength, n))
	}
	return nil
}

// validateCognitoTagValue validates a single tag value against the Smithy
// TagValueType length constraint [0, 256] (counted in Unicode characters).
func validateCognitoTagValue(value string) error {
	if n := utf8.RuneCountInString(value); n > tagutil.MaxTagValueLength {
		return awserrors.NewInvalidParameterException(
			fmt.Sprintf("Tag value length must not exceed %d: got %d", tagutil.MaxTagValueLength, n))
	}
	return nil
}

// cognitoTagViolation maps a tag-limit check verdict onto the Cognito
// InvalidParameterException shape. Both tag entry points — the map form and
// the HTTP tag-handler slice form — adjudicate through here so the two
// surfaces cannot drift apart in error semantics or message texts. The
// offender finders are lazy: only the verdict's own finder runs.
func cognitoTagViolation(v tagutil.Violation, firstKeyErr, firstValueErr func() error) error {
	switch v {
	case tagutil.TooManyTags:
		return awserrors.NewInvalidParameterException(
			fmt.Sprintf("Number of tags must not exceed %d", tagutil.MaxTagsPerResource))
	case tagutil.TagKeyTooShort, tagutil.TagKeyTooLong:
		return firstKeyErr()
	case tagutil.TagValueTooLong:
		return firstValueErr()
	case tagutil.ReservedTagKey:
		return awserrors.NewInvalidParameterException(
			"Tag keys cannot start with 'aws:' because the prefix is reserved for AWS use")
	}
	return nil
}

// validateCognitoTags validates a tag map against the Cognito tag limits:
// at most 50 tags per user pool or identity pool, keys of 1-128 characters,
// values of at most 256 characters and the aws: key prefix reserved for AWS
// use.
func validateCognitoTags(tags map[string]string) error {
	v, _ := tagutil.CheckStringTags(tags, tagutil.StandardLimits())
	return cognitoTagViolation(v,
		func() error { return cognitoTagKeyError(tags) },
		func() error { return cognitoTagValueError(tags) })
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
	v, _ := tagutil.CheckTags(tagList, tagutil.StandardLimits())
	firstKey := func() error {
		for _, t := range tagList {
			if err := validateCognitoTagKey(t.Key); err != nil {
				return err
			}
		}
		return nil
	}
	firstValue := func() error {
		for _, t := range tagList {
			if err := validateCognitoTagValue(t.Value); err != nil {
				return err
			}
		}
		return nil
	}
	return cognitoTagViolation(v, firstKey, firstValue)
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

// validDeviceRememberedStatuses is the Smithy DeviceRememberedStatusType
// enum (lowercase wire values).
var validDeviceRememberedStatuses = map[string]bool{
	"remembered":     true,
	"not_remembered": true,
}

// validateDeviceRememberedStatus returns true if the value is a recognised
// Smithy DeviceRememberedStatusType enum value.
func validateDeviceRememberedStatus(status string) bool {
	return validDeviceRememberedStatuses[status]
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

// Token-validity bounds from the Smithy model: AccessTokenValidityType and
// IdTokenValidityType share the maximum 86400; RefreshTokenValidityType
// maxes at 315360000.
const (
	maxTokenValidity        = 86400
	maxRefreshTokenValidity = 315360000
)

// validateAccessTokenValidity returns true if the value is within the Smithy
// AccessTokenValidityType range [1, 86400].
func validateAccessTokenValidity(v int) bool {
	return v >= 1 && v <= maxTokenValidity
}

// validateIdTokenValidity returns true if the value is within the Smithy
// IdTokenValidityType range [1, 86400].
func validateIdTokenValidity(v int) bool {
	return v >= 1 && v <= maxTokenValidity
}

// validateRefreshTokenValidity returns true if the value is within the Smithy
// RefreshTokenValidityType range [0, 315360000] — the only validity member
// whose minimum is 0.
func validateRefreshTokenValidity(v int) bool {
	return v >= 0 && v <= maxRefreshTokenValidity
}

// validatePrecedence returns true if the value satisfies the Smithy
// PrecedenceType minimum constraint (>= 0).
func validatePrecedence(v int) bool {
	return v >= 0
}

// standardSchemaAttributeNames is the set of standard user pool attribute
// names that may appear in the schema to modify their properties (data
// type, mutability, requiredness). AWS documents modifying standard
// attribute properties through the CreateUserPool Schema parameter and
// returns the standard set — including the 21-character
// phone_number_verified — in DescribeUserPool responses, so the custom
// attribute name constraints do not apply to them.
var standardSchemaAttributeNames = map[string]bool{
	"address": true, "birthdate": true, "email": true,
	"email_verified": true, "family_name": true, "gender": true,
	"given_name": true, "locale": true, "middle_name": true,
	"name": true, "nickname": true, "phone_number": true,
	"phone_number_verified": true, "picture": true,
	"preferred_username": true, "profile": true, "sub": true,
	"updated_at": true, "website": true, "zoneinfo": true,
	// identities is the reserved federation-populated standard attribute;
	// Amazon Cognito returns it in the DescribeUserPool standard set.
	"identities": true,
}

// customAttributeNamePattern is the Smithy CustomAttributeNameType pattern:
// ^[\p{L}\p{M}\p{S}\p{N}\p{P}]+$
var customAttributeNamePattern = regexp.MustCompile(`^[\p{L}\p{M}\p{S}\p{N}\p{P}]+$`)

// CustomAttributeNameType length range from the Smithy model.
const (
	minCustomAttributeNameLength = 1
	maxCustomAttributeNameLength = 20
)

// validateCustomAttributeName validates a custom attribute name against the
// Smithy CustomAttributeNameType length [1, 20] (counted in Unicode
// characters — the pattern admits multibyte categories) and pattern
// constraints.
func validateCustomAttributeName(name string) error {
	if n := utf8.RuneCountInString(name); n < minCustomAttributeNameLength || n > maxCustomAttributeNameLength {
		return awserrors.NewInvalidParameterException(
			fmt.Sprintf("Custom attribute name length must be %d-%d: got %d",
				minCustomAttributeNameLength, maxCustomAttributeNameLength, n))
	}
	if !customAttributeNamePattern.MatchString(name) {
		return awserrors.NewInvalidParameterException(
			fmt.Sprintf("Custom attribute name has invalid characters: %s", name))
	}
	return nil
}

// totpCodePattern matches totpCodeDigits-digit TOTP codes per RFC 6238.
var totpCodePattern = regexp.MustCompile(fmt.Sprintf(`^[0-9]{%d}$`, totpCodeDigits))

// validateMFADeliveryMedium returns true if the value is a recognised Smithy
// DeliveryMediumType enum value (SMS/EMAIL).
func validateMFADeliveryMedium(m string) bool {
	return validDeliveryMediums[m]
}

// RegionNameType length range from the Smithy model.
const (
	minRegionNameLength = 5
	maxRegionNameLength = 32
)

// validateRegionName validates a region name against the Smithy
// RegionNameType length constraint [5, 32] counted in Unicode characters.
func validateRegionName(name string) error {
	if n := utf8.RuneCountInString(name); n < minRegionNameLength || n > maxRegionNameLength {
		return awserrors.NewInvalidParameterException(
			fmt.Sprintf("RegionName length must be %d-%d: got %d", minRegionNameLength, maxRegionNameLength, n))
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

// deletionProtectionInactive is the wire value of the inactive
// DeletionProtectionLevelType enum member.
const deletionProtectionInactive = "INACTIVE"

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

// validExplicitAuthFlows is the ExplicitAuthFlowsType enum. Nine values come
// from the Smithy model; ALLOW_CLIENT_TOKEN_AUTH is recorded by the AWS
// developer guide's machine-to-machine authorisation pages — the vendored
// model enum lags the documented value — and is accepted on the wire.
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
	"ALLOW_CLIENT_TOKEN_AUTH":        true,
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

// validCustomAuthModes is the Smithy AdvancedSecurityEnabledModeType enum
// (the CustomAuthMode member of AdvancedSecurityAdditionalFlows) — unlike
// AdvancedSecurityModeType it carries no OFF value.
var validCustomAuthModes = map[string]bool{
	"AUDIT": true, "ENFORCED": true,
}

func validateCustomAuthMode(v string) bool { return validCustomAuthModes[v] }

// validAuthFactors is the Smithy AuthFactorType enum less SOFTWARE_TOKEN:
// the model documents SOFTWARE_TOKEN as not currently supported as a
// first authentication factor ("Do not include this value in
// AllowedFirstAuthFactors").
var validAuthFactors = map[string]bool{
	"PASSWORD": true, "EMAIL_OTP": true, "SMS_OTP": true, "WEB_AUTHN": true,
}

func validateAuthFactor(v string) bool { return validAuthFactors[v] }

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

// Length bounds of the name-shaped string types, which all carry the
// Smithy length range [1, 128]: UsernameType and GroupNameType (validated
// by validateUsernamePattern), UserPoolNameType, and ClientNameType.
const (
	minNameLength = 1
	maxNameLength = 128
)

// usernamePattern is the Smithy pattern for UsernameType and GroupNameType:
// ^[\p{L}\p{M}\p{S}\p{N}\p{P}]+$ (length 1-128).
var usernamePattern = regexp.MustCompile(`^[\p{L}\p{M}\p{S}\p{N}\p{P}]+$`)

// validateUsernamePattern returns true if the value matches the Smithy
// pattern and length constraint (1-128, counted in Unicode characters) for
// UsernameType / GroupNameType.
func validateUsernamePattern(v string) bool {
	n := utf8.RuneCountInString(v)
	return n >= minNameLength && n <= maxNameLength && usernamePattern.MatchString(v)
}

// userPoolNamePattern is the Smithy pattern for UserPoolNameType:
// ^[\w\s+=,.@-]+$ (length 1-128).
var userPoolNamePattern = regexp.MustCompile(`^[\w\s+=,.@-]+$`)

// validateUserPoolNamePattern returns true if the value matches the Smithy
// pattern and length constraint (1-128) for UserPoolNameType.
func validateUserPoolNamePattern(v string) bool {
	return len(v) >= minNameLength && len(v) <= maxNameLength && userPoolNamePattern.MatchString(v)
}

// Password-policy range bounds from the Smithy model:
// PasswordPolicyMinLengthType {6, 99}, TemporaryPasswordValidityDaysType
// {0, 365}, PasswordHistorySizeType {0, 24}, and
// AdminCreateUserUnusedAccountValidityDaysType {0, 365}.
const (
	MinPasswordPolicyMinimumLength   = 6
	MaxPasswordPolicyMinimumLength   = 99
	MaxTemporaryPasswordValidityDays = 365
	MaxPasswordHistorySize           = 24
	MaxUnusedAccountValidityDays     = 365
)

// UserImportJobNameType bounds from the Smithy model: length 1-128 and
// the pattern ^[\w\s+=,.@-]+$.
const (
	MinImportJobNameLength = 1
	MaxImportJobNameLength = 128
)

// importJobNamePattern is the UserImportJobNameType Smithy pattern.
var importJobNamePattern = regexp.MustCompile(`^[\w\s+=,.@-]+$`)

// validateImportJobName checks an import job name against the Smithy
// UserImportJobNameType constraints.
func validateImportJobName(name string) bool {
	if len(name) < MinImportJobNameLength || len(name) > MaxImportJobNameLength {
		return false
	}
	return importJobNamePattern.MatchString(name)
}

func validatePasswordHistorySize(v int) bool {
	return v >= 0 && v <= MaxPasswordHistorySize
}

// validateClientNamePattern returns true if the value matches the Smithy
// ClientNameType constraints (length 1-128, pattern ^[\w\s+=,.@-]+$, which is
// shared with UserPoolNameType).
func validateClientNamePattern(v string) bool {
	return len(v) >= minNameLength && len(v) <= maxNameLength && userPoolNamePattern.MatchString(v)
}

// ArnType length range from the Smithy model.
const (
	minArnLength = 20
	maxArnLength = 2048
)

// arnTypePattern is the Smithy ArnType pattern.
var arnTypePattern = regexp.MustCompile(`^arn:[\w+=/,.@-]+:[\w+=/,.@-]+:([\w+=/,.@-]*)?:[0-9]+:[\w+=/,.@-]+(:[\w+=/,.@-]+)?(:[\w+=/,.@-]+)?$`)

// validateArnType returns true if the value matches the Smithy ArnType
// constraints (length 20-2048 plus the generic ARN pattern).
func validateArnType(v string) bool {
	return len(v) >= minArnLength && len(v) <= maxArnLength && arnTypePattern.MatchString(v)
}

// validateExplicitMinimumLength checks the Smithy
// PasswordPolicyMinimumLengthType range for a value explicitly supplied in
// a request. Unlike validatePasswordPolicyRanges, zero is not treated as
// "unset" here because the member was present.
func validateExplicitMinimumLength(v int) error {
	if v < MinPasswordPolicyMinimumLength || v > MaxPasswordPolicyMinimumLength {
		return awserrors.NewInvalidParameterException(fmt.Sprintf(
			"MinimumLength must be between %d and %d", MinPasswordPolicyMinimumLength, MaxPasswordPolicyMinimumLength))
	}
	return nil
}

// validatePasswordPolicyRanges checks the Smithy range traits of the
// password-policy members. A zero MinimumLength means the member was never
// set (the store default is applied instead), so it is accepted here;
// explicitly supplied values are range-checked at parse time.
func validatePasswordPolicyRanges(p *cognitostore.PasswordPolicy) error {
	if p == nil {
		return nil
	}
	if p.MinimumLength != 0 && (p.MinimumLength < MinPasswordPolicyMinimumLength || p.MinimumLength > MaxPasswordPolicyMinimumLength) {
		return awserrors.NewInvalidParameterException(fmt.Sprintf(
			"MinimumLength must be between %d and %d", MinPasswordPolicyMinimumLength, MaxPasswordPolicyMinimumLength))
	}
	if p.TemporaryPasswordValidityDays < 0 || p.TemporaryPasswordValidityDays > MaxTemporaryPasswordValidityDays {
		return awserrors.NewInvalidParameterException(fmt.Sprintf(
			"TemporaryPasswordValidityDays must be between 0 and %d", MaxTemporaryPasswordValidityDays))
	}
	if !validatePasswordHistorySize(p.PasswordHistorySize) {
		return awserrors.NewInvalidParameterException(fmt.Sprintf(
			"PasswordHistorySize must be between 0 and %d", MaxPasswordHistorySize))
	}
	return nil
}

// validateUserPoolConfig performs the model-derived validation that applies to
// a whole user pool, on both the create and the update path. It is the single
// validation entry point shared by the HTTP handler, the admin handler, and
// the Core functions, so that no transport bypasses it.
func validateUserPoolConfig(pool *cognitostore.UserPool) error {
	if pool == nil {
		return ErrInvalidParameter
	}
	if !validateUserPoolNamePattern(pool.Name) {
		return ErrInvalidParameter
	}
	// AliasAttributes and UsernameAttributes are mutually exclusive.
	if len(pool.AliasAttributes) > 0 && len(pool.UsernameAttributes) > 0 {
		return ErrInvalidParameter
	}
	if pool.MfaConfiguration != "" && !validateUserPoolMfaConfig(pool.MfaConfiguration) {
		return ErrInvalidParameter
	}
	if pool.DeletionProtection != "" && !validateDeletionProtection(pool.DeletionProtection) {
		return ErrInvalidParameter
	}
	for _, a := range pool.AutoVerifiedAttributes {
		if !validateVerifiedAttribute(a) {
			return ErrInvalidParameter
		}
	}
	for _, a := range pool.AliasAttributes {
		if !validateAliasAttribute(a) {
			return ErrInvalidParameter
		}
	}
	for _, a := range pool.UsernameAttributes {
		if !validateUsernameAttribute(a) {
			return ErrInvalidParameter
		}
	}
	if err := validatePasswordPolicyRanges(pool.PasswordPolicy); err != nil {
		return err
	}
	if pool.LambdaConfig != nil {
		for _, arn := range []string{
			pool.LambdaConfig.PreSignUp,
			pool.LambdaConfig.CustomMessage,
			pool.LambdaConfig.PostConfirmation,
			pool.LambdaConfig.PreAuthentication,
			pool.LambdaConfig.PostAuthentication,
			pool.LambdaConfig.DefineAuthChallenge,
			pool.LambdaConfig.CreateAuthChallenge,
			pool.LambdaConfig.VerifyAuthChallengeResponse,
			pool.LambdaConfig.PreTokenGeneration,
			pool.LambdaConfig.UserMigration,
		} {
			if arn != "" && !validateArnType(arn) {
				return ErrInvalidParameter
			}
		}
	}
	if pool.EmailConfiguration != nil {
		if pool.EmailConfiguration.EmailSendingAccount != "" && !validateEmailSendingAccount(pool.EmailConfiguration.EmailSendingAccount) {
			return ErrInvalidParameter
		}
		if pool.EmailConfiguration.SourceArn != "" && !validateArnType(pool.EmailConfiguration.SourceArn) {
			return ErrInvalidParameter
		}
	}
	if pool.SmsConfiguration != nil && pool.SmsConfiguration.SnsCallerArn != "" && !validateArnType(pool.SmsConfiguration.SnsCallerArn) {
		return ErrInvalidParameter
	}
	if pool.AdminCreateUserConfig != nil {
		days := pool.AdminCreateUserConfig.UnusedAccountValidityDays
		if days < 0 || days > MaxUnusedAccountValidityDays {
			return ErrInvalidParameter
		}
	}
	for _, sa := range pool.SchemaAttributes {
		// The schema name domain is the union of the standard attribute
		// names and custom attribute names; the custom-attribute
		// constraints apply only to the latter.
		if !standardSchemaAttributeNames[sa.Name] {
			if err := validateCustomAttributeName(sa.Name); err != nil {
				return err
			}
		}
		if sa.AttributeDataType != "" && !validateAttributeDataType(sa.AttributeDataType) {
			return ErrInvalidParameter
		}
	}
	return nil
}

// listLimitMax is the upper bound shared by the Cognito list-limit shapes
// (QueryLimitType, PoolQueryLimitType, QueryLimit — all Smithy max 60).
// The store package owns the value; this alias is the package-local name
// the parsers and Core defaults reference.
const listLimitMax = cognitostore.MaxListLimit

// parseListLimit extracts a list-limit parameter typed QueryLimitType in the
// Smithy model (min 0, e.g. ListUsers.Limit). Absent or zero selects
// maxValue as the page size — the AWS-documented default equals the shape
// maximum; an explicitly provided value outside 0-maxValue is rejected.
func parseListLimit(params map[string]interface{}, key string, maxValue int) (int, error) {
	if err := rejectNonNumericLimit(params, key, "0", maxValue); err != nil {
		return 0, err
	}
	raw, present := params[key]
	if !present {
		return maxValue, nil
	}
	v, _ := parseJSONInt(raw)
	if v == 0 {
		return maxValue, nil
	}
	if v < 0 || v > maxValue {
		return 0, awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("%s must be between 0 and %d", key, maxValue), http.StatusBadRequest)
	}
	return v, nil
}

// parseStrictListLimit extracts a list-limit parameter whose Smithy shape
// has a minimum of 1 (PoolQueryLimitType / QueryLimit /
// ListResourceServersLimitType, e.g. ListUserPools, ListUserPoolClients,
// ListResourceServers). Absent selects maxValue as the page size — the
// AWS-documented default equals the shape maximum; an explicitly provided
// value outside 1-maxValue is rejected, including an explicit zero.
func parseStrictListLimit(params map[string]interface{}, key string, maxValue int) (int, error) {
	if err := rejectNonNumericLimit(params, key, "1", maxValue); err != nil {
		return 0, err
	}
	raw, present := params[key]
	if !present {
		return maxValue, nil
	}
	v, _ := parseJSONInt(raw)
	if v < 1 || v > maxValue {
		return 0, awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("%s must be between 1 and %d", key, maxValue), http.StatusBadRequest)
	}
	return v, nil
}

// rejectNonNumericLimit fails closed when a limit parameter is present but
// cannot be read as a JSON number. Presence at the exact awsJson1_1 member
// name is the only form the parser recognises, so an uninterpretable value
// is a wire-shape violation, never a silently omitted member.
func rejectNonNumericLimit(params map[string]interface{}, key, lowerBound string, maxValue int) error {
	if v, ok := params[key]; ok {
		if _, isInt := parseJSONInt(v); !isInt {
			return awserrors.NewAWSError("InvalidParameterException",
				fmt.Sprintf("%s must be a number between %s and %d", key, lowerBound, maxValue), http.StatusBadRequest)
		}
	}
	return nil
}

// validateAccountTakeoverAction validates against the Smithy enum
// AccountTakeoverEventActionType: BLOCK, MFA_IF_CONFIGURED, MFA_REQUIRED, NO_ACTION.
func validateAccountTakeoverAction(action string) bool {
	switch action {
	case "BLOCK", "MFA_IF_CONFIGURED", "MFA_REQUIRED", "NO_ACTION":
		return true
	}
	return false
}
