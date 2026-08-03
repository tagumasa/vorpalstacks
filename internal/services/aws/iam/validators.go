// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"encoding/json"
	"regexp"
	"strconv"
	"unicode"

	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/aws/types"
)

// ---------------------------------------------------------------------------
// Entity name patterns (Smithy trait source)
// ---------------------------------------------------------------------------

// entityNamePattern validates IAM entity names for users and roles.
// Per Smithy userNameType/roleNameType: length 1-64, pattern ^[\w+=,.@-]+$.
var entityNamePattern = regexp.MustCompile(`^[\w+=,.@-]{1,64}$`)

// entityNamePattern128 validates entity names that allow up to 128
// characters: groups, instance profiles, server certificates, policies,
// and inline policy names. Per Smithy groupNameType/instanceProfileNameType/
// serverCertificateNameType/policyNameType: length 1-128, pattern ^[\w+=,.@-]+$.
var entityNamePattern128 = regexp.MustCompile(`^[\w+=,.@-]{1,128}$`)

// samlProviderNamePattern validates SAML provider names. Per Smithy
// samlProviderNameType: length 1-128, pattern ^[\w._-]+$.
var samlProviderNamePattern = regexp.MustCompile(`^[\w._-]{1,128}$`)

// samlPrivateKeyIdPattern validates SAML private key IDs. Per Smithy
// privateKeyIdType: length [22,64] (checked separately), pattern ^[A-Z0-9]+$.
var samlPrivateKeyIdPattern = regexp.MustCompile(`^[A-Z0-9]+$`)

// accountAliasPattern validates account alias names. Per Smithy
// accountAliasType: lowercase alphanumeric with no consecutive hyphens.
// Length 3-63 is enforced separately.
var accountAliasPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// tagKeyPattern validates tag keys. Per Smithy tagKeyType: length 1-128,
// pattern ^[\p{L}\p{Z}\p{N}_.:/=+\-@]+$.
var tagKeyPattern = regexp.MustCompile(`^[\p{L}\p{Z}\p{N}_.:/=+\-@]{1,128}$`)

// iamPolicyArnPattern validates IAM policy ARNs for permissions boundary
// operations.  Only customer-managed policies (numeric account ID) are
// accepted — AWS-managed policies cannot be set as permissions boundaries.
var iamPolicyArnPattern = regexp.MustCompile(`^arn:aws:iam::\d+:policy/.+$`)

// attachPolicyArnPattern validates IAM policy ARNs for Attach/Detach*Policy
// operations.  Accepts both customer-managed (numeric account ID) and
// AWS-managed policies (literal "aws" account ID, e.g.
// arn:aws:iam::aws:policy/AdministratorAccess).
var attachPolicyArnPattern = regexp.MustCompile(`^arn:aws:iam::(?:\d+|aws):policy/.+$`)

// awsServiceNamePattern validates AWSServiceName for service-linked roles.
// Must be a dotted service name like "ec2.amazonaws.com".
var awsServiceNamePattern = regexp.MustCompile(`^[A-Za-z0-9_-]+(\.[A-Za-z0-9_-]+)+$`)

// ---------------------------------------------------------------------------
// Path validation (Smithy pathType)
// ---------------------------------------------------------------------------

// pathTypePattern implements the Smithy pathType pattern:
//
//	^(\u002F)|(\u002F[\u0021-\u007E]+\u002F)$
//
// i.e. either "/" alone or a path starting and ending with "/" using
// printable ASCII characters.
var pathTypePattern = regexp.MustCompile(`^(/)|(/[!-~]+/)$`)

// validatePath checks that a path conforms to the Smithy pathType trait:
// pattern ^(\u002F)|(\u002F[\u0021-\u007E]+\u002F)$ and length 1-512.
func validatePath(path string) bool {
	if len(path) < 1 || len(path) > 512 {
		return false
	}
	return pathTypePattern.MatchString(path)
}

// ---------------------------------------------------------------------------
// Thumbprint validation (Smithy thumbprintType)
// ---------------------------------------------------------------------------

// thumbprintPattern matches a 40-character hex-encoded SHA-1 hash.
// Per Smithy thumbprintType: length 40; AWS docs specify hex-encoded.
var thumbprintPattern = regexp.MustCompile(`^[0-9A-Fa-f]{40}$`)

// validateThumbprint checks that a thumbprint is a 40-character hex string.
func validateThumbprint(tp string) bool {
	return thumbprintPattern.MatchString(tp)
}

// ---------------------------------------------------------------------------
// Client ID validation (Smithy clientIDType)
// ---------------------------------------------------------------------------

// validateClientID checks that a client ID conforms to Smithy clientIDType:
// length 1-255.
func validateClientID(id string) bool {
	return len(id) >= 1 && len(id) <= 255
}

// ---------------------------------------------------------------------------
// Service namespace validation (Smithy serviceNamespaceType)
// ---------------------------------------------------------------------------

// validateServiceNamespace checks the length of a service namespace.
// Per Smithy serviceNamespaceType: length 1-64.  The Smithy pattern
// trait ^[\w-]*$ is not enforced because AWS service names in practice
// include dots (e.g. "codecommit.amazonaws.com") which the pattern
// would reject.
func validateServiceNamespace(ns string) bool {
	return len(ns) >= 1 && len(ns) <= 64
}

// ---------------------------------------------------------------------------
// Custom suffix validation (Smithy customSuffixType)
// ---------------------------------------------------------------------------

// customSuffixPattern implements Smithy customSuffixType:
// pattern ^[\w+=,.@-]+$ and length 1-64.
var customSuffixPattern = regexp.MustCompile(`^[\w+=,.@-]+$`)

// validateCustomSuffix checks that a custom suffix conforms to
// Smithy customSuffixType: length 1-64, pattern ^[\w+=,.@-]+$.
func validateCustomSuffix(suffix string) bool {
	if len(suffix) < 1 || len(suffix) > 64 {
		return false
	}
	return customSuffixPattern.MatchString(suffix)
}

// ---------------------------------------------------------------------------
// Role MaxSessionDuration validation (Smithy roleMaxSessionDurationType)
// ---------------------------------------------------------------------------

// validateRoleMaxSessionDuration checks that a MaxSessionDuration value
// conforms to Smithy roleMaxSessionDurationType: range 3600-43200.
func validateRoleMaxSessionDuration(v int) bool {
	return v >= 3600 && v <= 43200
}

// ---------------------------------------------------------------------------
// Role description validation (Smithy roleDescriptionType)
// ---------------------------------------------------------------------------

// roleDescriptionPattern implements Smithy roleDescriptionType:
//
//	pattern ^[\u0009\u000A\u000D\u0020-\u007E\u00A1-\u00FF]*$
//
// i.e. tab, LF, CR, printable ASCII, or Latin-1 supplement characters.
var roleDescriptionPattern = regexp.MustCompile(`^[\t\n\r\x20-\x7E\xA1-\xFF]*$`)

// validateRoleDescription checks that a role description conforms to
// Smithy roleDescriptionType: length 0-1000 and the pattern above.
func validateRoleDescription(s string) bool {
	if len(s) > 1000 {
		return false
	}
	return roleDescriptionPattern.MatchString(s)
}

// ---------------------------------------------------------------------------
// Policy scope validation (Smithy policyScopeType)
// ---------------------------------------------------------------------------

// validatePolicyScope checks that a ListPolicies Scope parameter conforms
// to Smithy policyScopeType: enum {All, AWS, Local}.
func validatePolicyScope(s string) bool {
	switch s {
	case "All", "AWS", "Local":
		return true
	default:
		return false
	}
}

// ---------------------------------------------------------------------------
// SAML metadata certificate extraction patterns
// ---------------------------------------------------------------------------

// x509CertDataPattern extracts X509Certificate elements from SAML metadata
// XML, supporting both bare and ds:-prefixed element names.
var x509CertDataPattern = regexp.MustCompile(`(?s)<(?:ds:)?X509Certificate>([^<]+)</(?:ds:)?X509Certificate>`)

// whitespacePattern matches any run of whitespace characters, used to
// normalise extracted certificate data before Base64 decoding.
var whitespacePattern = regexp.MustCompile(`\s+`)

// ---------------------------------------------------------------------------
// Password policy validation
// ---------------------------------------------------------------------------

// symbolPattern matches a single symbol character for password policy
// RequireSymbols enforcement.
var symbolPattern = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':",\.<>?/\\|~]`)

// numberPattern matches a single digit character for password policy
// RequireNumbers enforcement.
var numberPattern = regexp.MustCompile(`[0-9]`)

// validatePasswordAgainstPolicy checks that a password satisfies every
// applicable requirement of the account password policy.  Each boolean
// flag on the policy is checked independently, matching AWS behaviour.
func validatePasswordAgainstPolicy(password string, policy *iamstore.AccountPasswordPolicy) bool {
	if len(password) < policy.MinimumPasswordLength {
		return false
	}

	if policy.RequireSymbols {
		if !symbolPattern.MatchString(password) {
			return false
		}
	}

	if policy.RequireNumbers {
		if !numberPattern.MatchString(password) {
			return false
		}
	}

	if policy.RequireUppercaseCharacters {
		hasUpper := false
		for _, r := range password {
			if unicode.IsUpper(r) {
				hasUpper = true
				break
			}
		}
		if !hasUpper {
			return false
		}
	}

	if policy.RequireLowercaseCharacters {
		hasLower := false
		for _, r := range password {
			if unicode.IsLower(r) {
				hasLower = true
				break
			}
		}
		if !hasLower {
			return false
		}
	}

	return true
}

// ---------------------------------------------------------------------------
// Policy document validation
// ---------------------------------------------------------------------------

// policyValidationMode controls which fields are required for each statement.
type policyValidationMode int

const (
	// policyModeManaged is for identity-based policies (managed and inline).
	// Requires Effect + Action/NotAction + Resource/NotResource per statement.
	policyModeManaged policyValidationMode = iota
	// policyModeTrust is for AssumeRolePolicyDocument (resource-based trust
	// policies). Requires Effect + Action/NotAction + Principal/NotPrincipal.
	policyModeTrust
)

// validatePolicyDocument checks if a policy document is valid JSON and has
// the minimum required structure for an IAM identity-based policy: a
// top-level object with a "Statement" field containing at least one
// statement object with Effect, Action/NotAction, and Resource/NotResource.
func validatePolicyDocument(document string) bool {
	return validatePolicyDocumentMode(document, policyModeManaged)
}

// validateTrustPolicyDocument validates an AssumeRolePolicyDocument (trust
// policy).  Each statement must have Effect, Action/NotAction, and
// Principal/NotPrincipal.
func validateTrustPolicyDocument(document string) bool {
	return validatePolicyDocumentMode(document, policyModeTrust)
}

func validatePolicyDocumentMode(document string, mode policyValidationMode) bool {
	if document == "" {
		return false
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(document), &raw); err != nil {
		return false
	}

	// Version is a top-level policy-document field (alongside Statement).
	// When present it must be a recognised policy language version; an
	// unknown value would silently change evaluation semantics.
	if rawVersion, ok := raw["Version"]; ok {
		var version string
		if err := json.Unmarshal(rawVersion, &version); err != nil {
			return false
		}
		if version != "2012-10-17" && version != "2008-10-17" {
			return false
		}
	}

	statementsRaw, ok := raw["Statement"]
	if !ok {
		return false
	}

	// Statement can be a single object or an array of objects.
	var singleStmt map[string]interface{}
	if err := json.Unmarshal(statementsRaw, &singleStmt); err == nil {
		return validateStatement(singleStmt, mode)
	}

	var stmtArray []map[string]interface{}
	if err := json.Unmarshal(statementsRaw, &stmtArray); err != nil {
		return false
	}
	if len(stmtArray) == 0 {
		return false
	}
	for _, stmt := range stmtArray {
		if !validateStatement(stmt, mode) {
			return false
		}
	}
	return true
}

// validateStatement checks that a single policy statement has the required
// members for the given validation mode:
//   - Effect must be "Allow" or "Deny"
//   - Action or NotAction must be present (but not both — they are mutually
//     exclusive per the IAM JSON policy specification)
//   - For managed policies: Resource or NotResource must be present (mutually
//     exclusive)
//   - For trust policies: Principal or NotPrincipal must be present
//
// Empty strings ("Action":""), empty arrays ("Action":[]), and JSON null
// are all treated as absent to prevent fail-OPEN acceptance of malformed
// policies that could silently grant unintended permissions.
func validateStatement(stmt map[string]interface{}, mode policyValidationMode) bool {
	effect, ok := stmt["Effect"].(string)
	if !ok {
		return false
	}
	if effect != "Allow" && effect != "Deny" {
		return false
	}

	// Action and NotAction are mutually exclusive — at least one must be
	// present, but not both.
	hasAction := hasPolicyKey(stmt, "Action")
	hasNotAction := hasPolicyKey(stmt, "NotAction")
	if !hasAction && !hasNotAction {
		return false
	}
	if hasAction && hasNotAction {
		return false
	}

	switch mode {
	case policyModeManaged:
		// Identity-based policies require Resource or NotResource — mutually
		// exclusive per the IAM JSON policy specification.
		hasResource := hasPolicyKey(stmt, "Resource")
		hasNotResource := hasPolicyKey(stmt, "NotResource")
		if !hasResource && !hasNotResource {
			return false
		}
		if hasResource && hasNotResource {
			return false
		}
	case policyModeTrust:
		// Trust policies require Principal or NotPrincipal.
		if !hasPolicyKey(stmt, "Principal") && !hasPolicyKey(stmt, "NotPrincipal") {
			return false
		}
	}

	return true
}

// hasPolicyKey returns true if the statement map contains the given key
// with a meaningful value. The following are all treated as absent:
//   - key not present
//   - JSON null
//   - empty string ""
//   - empty array []interface{}{}
//   - empty map map[string]interface{}{}
//
// This prevents fail-OPEN acceptance of malformed policies where a
// required field is technically present but semantically empty.
func hasPolicyKey(stmt map[string]interface{}, key string) bool {
	val, ok := stmt[key]
	if !ok {
		return false
	}
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case string:
		return v != ""
	case []interface{}:
		return len(v) > 0
	case map[string]interface{}:
		return len(v) > 0
	default:
		return true
	}
}

// ---------------------------------------------------------------------------
// Tag validation
// ---------------------------------------------------------------------------

// validateTagEntries validates the key and value length limits for each
// individual tag entry.  It does NOT check the total tag count — the caller
// is responsible for that because the acceptable count depends on context
// (new tags on Create vs. merged tags on TagResource).
func validateTagEntries(newTags []types.Tag) error {
	for _, t := range newTags {
		if len(t.Key) == 0 || len(t.Key) > MaxTagKeyLength {
			return NewInvalidInputError("TagKey", "must be 1 to "+strconv.Itoa(MaxTagKeyLength)+" characters")
		}
		if !tagKeyPattern.MatchString(t.Key) {
			return NewInvalidInputError("TagKey", "contains invalid characters")
		}
		if len(t.Value) > MaxTagValueLength {
			return NewInvalidInputError("TagValue", "must be 0 to "+strconv.Itoa(MaxTagValueLength)+" characters")
		}
	}
	return nil
}

// validateNewTags validates both per-tag entry limits and the total tag
// count for resources being created.  On Create operations there are no
// pre-existing tags, so the total count is simply len(newTags).
func validateNewTags(newTags []types.Tag) error {
	if err := validateTagEntries(newTags); err != nil {
		return err
	}
	if len(newTags) > MaxTagsPerResource {
		return NewInvalidInputError("Tags", "exceeds maximum of "+strconv.Itoa(MaxTagsPerResource)+" tags per resource")
	}
	return nil
}

// validateAssertionEncryptionMode validates the Smithy
// assertionEncryptionModeType enum.  Allowed values: "Required",
// "Allowed".  Empty string is valid (means the field was not specified).
func validateAssertionEncryptionMode(mode string) bool {
	return mode == "" || mode == "Required" || mode == "Allowed"
}

// ---------------------------------------------------------------------------
// Centralised regex-backed validators (M5 refactor).
//
// These helpers wrap the bare regex variables above so that callers do not
// have to repeat the pattern+error-message dance at every call site.  Each
// helper returns a fully-formed *awserrors.AWSError so callers can return
// it directly.  Keeping the regex variables exported (package-level) lets
// ad-hoc callers (tests, validators.go itself) still reach them, but all
// production callers should go through these helpers to keep error messages
// uniform.
// ---------------------------------------------------------------------------

// validateEntityName checks names conforming to the Smithy userNameType /
// roleNameType pattern (length 1-64, [\w+=,.@-]).
func validateEntityName(name, paramName string) error {
	if !entityNamePattern.MatchString(name) {
		return NewInvalidInputError(paramName, "must be 1 to 64 alphanumeric characters or any of +=,.@-_")
	}
	return nil
}

// validateEntityName128 checks names conforming to the Smithy
// groupNameType / instanceProfileNameType / serverCertificateNameType /
// policyNameType pattern (length 1-128, [\w+=,.@-]).
func validateEntityName128(name, paramName string) error {
	if !entityNamePattern128.MatchString(name) {
		return NewInvalidInputError(paramName, "must be 1 to 128 alphanumeric characters or any of +=,.@-_")
	}
	return nil
}

// validateIAMPolicyArn checks ARNs acceptable as a permissions boundary.
// Only customer-managed policies (numeric account ID) are accepted — AWS
// managed policies (the literal "aws" account) cannot serve as a
// permissions boundary.
func validateIAMPolicyArn(arn string) error {
	if !iamPolicyArnPattern.MatchString(arn) {
		return NewInvalidInputError("PermissionsBoundary", "must be a valid customer-managed IAM policy ARN")
	}
	return nil
}

// validateAttachPolicyArn checks ARNs acceptable for Attach/Detach*Policy
// operations.  Both customer-managed (numeric account ID) and AWS-managed
// (literal "aws") policy ARNs are accepted.
func validateAttachPolicyArn(arn string) error {
	if !attachPolicyArnPattern.MatchString(arn) {
		return NewInvalidInputError("PolicyArn", "must be a valid IAM policy ARN")
	}
	return nil
}

// validateAccountAlias checks the Smithy accountAliasType pattern.
// Length 3-63 is enforced here (not in the regex).
func validateAccountAlias(alias string) error {
	if len(alias) < 3 || len(alias) > 63 || !accountAliasPattern.MatchString(alias) {
		return NewInvalidInputError("AccountAlias", "must be 3 to 63 characters; lowercase letters, digits, and hyphens; no consecutive hyphens")
	}
	return nil
}

// validateSAMLProviderName checks the Smithy samlProviderNameType pattern.
// Length 1-128 is enforced by the regex.
func validateSAMLProviderName(name string) error {
	if !samlProviderNamePattern.MatchString(name) {
		return NewInvalidInputError("SAMLProviderName", "must be 1 to 128 characters; allowed: word characters, dot, underscore, hyphen")
	}
	return nil
}

// validateSAMLPrivateKeyId checks the Smithy privateKeyIdType pattern.
// Length bounds are checked separately by callers when relevant.
func validateSAMLPrivateKeyId(id string) error {
	if !samlPrivateKeyIdPattern.MatchString(id) {
		return NewInvalidInputError("PrivateKeyId", "must be uppercase alphanumeric characters only")
	}
	return nil
}

// validateAWSServiceName checks the dotted service-name pattern used by
// CreateServiceLinkedRole's AWSServiceName parameter.
func validateAWSServiceName(name string) error {
	if !awsServiceNamePattern.MatchString(name) {
		return NewInvalidInputError("AWSServiceName", "must be a dotted service name (e.g. ec2.amazonaws.com)")
	}
	return nil
}
