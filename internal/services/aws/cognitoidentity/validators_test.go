package cognitoidentity

import (
	"strings"
	"testing"
)

// TestValidateTagKeysAndValuesUnicodeLengths pins that identity pool tag
// key and value limits are counted in Unicode characters per the Smithy
// TagKeysType and TagValueType @length traits.
func TestValidateTagKeysAndValuesUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	tags := map[string]string{strings.Repeat(cjk, 128): strings.Repeat(cjk, 256)}
	if !validateTagKeys(tags) {
		t.Error("128-character CJK key rejected by validateTagKeys")
	}
	if !validateTagValues(tags) {
		t.Error("256-character CJK value rejected by validateTagValues")
	}
	if validateTagKeys(map[string]string{strings.Repeat(cjk, 129): "v"}) {
		t.Error("129-character CJK key accepted by validateTagKeys")
	}
	if validateTagKeys(map[string]string{"": "v"}) {
		t.Error("empty key accepted by validateTagKeys")
	}
	if validateTagValues(map[string]string{"k": strings.Repeat(cjk, 257)}) {
		t.Error("257-character CJK value accepted by validateTagValues")
	}
}

// TestValidateLoginsKeysUnicodeLengths pins that Logins map keys are limited
// to 128 Unicode characters per the Smithy IdentityProviderName @length
// trait, not 128 bytes.
func TestValidateLoginsKeysUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if !validateLoginsKeys(map[string]string{strings.Repeat(cjk, 128): "token"}) {
		t.Error("128-character CJK login key rejected by validateLoginsKeys")
	}
	if validateLoginsKeys(map[string]string{strings.Repeat(cjk, 129): "token"}) {
		t.Error("129-character CJK login key accepted by validateLoginsKeys")
	}
	if validateLoginsKeys(map[string]string{"": "token"}) {
		t.Error("empty login key accepted by validateLoginsKeys")
	}
}

// TestValidateRoleMappingUnicodeLengths pins that role-mapping keys, mapping
// rule claims and claim values are counted in Unicode characters per the
// Smithy IdentityProviderName, ClaimName and ClaimValue @length traits.
func TestValidateRoleMappingUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"
	validRule := MappingRuleInput{
		Claim:     "isAdmin",
		MatchType: "Equals",
		Value:     "yes",
		RoleARN:   "arn:aws:iam::123456789012:role/valid",
	}

	keyMappings := map[string]RoleMappingInput{
		strings.Repeat(cjk, 100): {
			Type:                    "Rules",
			AmbiguousRoleResolution: "AuthenticatedRole",
			RulesConfiguration:      &RulesConfigInput{Rules: []MappingRuleInput{validRule}},
		},
	}
	if !validateRoleMappings(keyMappings) {
		t.Error("100-character CJK role-mapping key rejected by validateRoleMappings")
	}

	claimRule := validRule
	claimRule.Claim = strings.Repeat(cjk, 64)
	if !validateMappingRule(claimRule) {
		t.Error("64-character CJK claim rejected by validateMappingRule")
	}
	claimRule.Claim = strings.Repeat(cjk, 65)
	if validateMappingRule(claimRule) {
		t.Error("65-character CJK claim accepted by validateMappingRule")
	}

	valueRule := validRule
	valueRule.Value = strings.Repeat(cjk, 128)
	if !validateMappingRule(valueRule) {
		t.Error("128-character CJK claim value rejected by validateMappingRule")
	}
	valueRule.Value = strings.Repeat(cjk, 129)
	if validateMappingRule(valueRule) {
		t.Error("129-character CJK claim value accepted by validateMappingRule")
	}
}

// TestValidateScalarUnicodeLengths pins the remaining no-pattern @length
// gates: PrincipalTagValue (1-256), IdentityProviderToken (1-50000),
// DeveloperUserIdentifier (1-1024) and IdentityProviderName (1-128), all
// counted in Unicode characters.
func TestValidateScalarUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if !validatePrincipalTagValue(strings.Repeat(cjk, 256)) {
		t.Error("256-character CJK principal tag value rejected")
	}
	if validatePrincipalTagValue(strings.Repeat(cjk, 257)) {
		t.Error("257-character CJK principal tag value accepted")
	}

	if !validateLoginTokenValue(strings.Repeat(cjk, 50000)) {
		t.Error("50000-character CJK login token rejected")
	}
	if validateLoginTokenValue(strings.Repeat(cjk, 50001)) {
		t.Error("50001-character CJK login token accepted")
	}

	if !validateDeveloperUserIdentifier(strings.Repeat(cjk, 1024)) {
		t.Error("1024-character CJK developer user identifier rejected")
	}
	if validateDeveloperUserIdentifier(strings.Repeat(cjk, 1025)) {
		t.Error("1025-character CJK developer user identifier accepted")
	}

	if !validateIdentityProviderNameLength(strings.Repeat(cjk, 128)) {
		t.Error("128-character CJK identity provider name rejected")
	}
	if validateIdentityProviderNameLength(strings.Repeat(cjk, 129)) {
		t.Error("129-character CJK identity provider name accepted")
	}
}
