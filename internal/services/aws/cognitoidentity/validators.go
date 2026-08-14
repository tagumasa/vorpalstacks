package cognitoidentity

import "regexp"

// ---------------------------------------------------------------------------
// Smithy-derived patterns
// ---------------------------------------------------------------------------

var (
	// IdentityPoolName: ^[\w\s+=,.@-]+$, length 1-128
	identityPoolNamePattern = regexp.MustCompile(`^[\w\s+=,.@-]+$`)

	// DeveloperProviderName: ^[\w._-]+$, length 1-128
	developerProviderNamePattern = regexp.MustCompile(`^[\w._-]+$`)

	// CognitoIdentityProviderName: ^[\w._:/-]+$, length 1-128
	providerNamePattern = regexp.MustCompile(`^[\w._:/-]+$`)

	// CognitoIdentityProviderClientId: ^[\w_]+$, length 1-128
	providerClientIdPattern = regexp.MustCompile(`^[\w_]+$`)

	// ClaimName: ^[\p{L}\p{M}\p{S}\p{N}\p{P}]+$, length 1-64
	claimNamePattern = regexp.MustCompile(`^[\p{L}\p{M}\p{S}\p{N}\p{P}]+$`)

	// IdentityId / IdentityPoolId: ^[\w-]+:[0-9a-f-]+$, length 1-55
	identityIdPattern = regexp.MustCompile(`^[\w-]+:[0-9a-f-]+$`)

	// PaginationKey: ^[\S]+$, length 1-65535
	paginationKeyPattern = regexp.MustCompile(`^[\S]+$`)

	// AccountId: ^\d+$, length 1-15
	accountIdPattern = regexp.MustCompile(`^\d+$`)
)

// ---------------------------------------------------------------------------
// Smithy-derived enum sets
// ---------------------------------------------------------------------------

var validRoleTypes = map[string]bool{
	"authenticated":   true,
	"unauthenticated": true,
}

var validRoleMappingTypes = map[string]bool{
	"Token": true,
	"Rules": true,
}

var validAmbiguousRoleResolutions = map[string]bool{
	"AuthenticatedRole": true,
	"Deny":              true,
}

var validMappingRuleMatchTypes = map[string]bool{
	"Equals":     true,
	"Contains":   true,
	"StartsWith": true,
	"NotEqual":   true,
}

// ---------------------------------------------------------------------------
// Validators
// ---------------------------------------------------------------------------

func validateIdentityPoolName(name string) bool {
	if len(name) < 1 || len(name) > 128 {
		return false
	}
	return identityPoolNamePattern.MatchString(name)
}

func validateDeveloperProviderName(name string) bool {
	if len(name) < 1 || len(name) > 128 {
		return false
	}
	return developerProviderNamePattern.MatchString(name)
}

func validateProviderName(name string) bool {
	if len(name) < 1 || len(name) > 128 {
		return false
	}
	return providerNamePattern.MatchString(name)
}

func validateProviderClientId(id string) bool {
	if len(id) < 1 || len(id) > 128 {
		return false
	}
	return providerClientIdPattern.MatchString(id)
}

// validateQueryLimit enforces the Smithy QueryLimit range (min=1, max=60)
// for MaxResults on ListIdentityPools, ListIdentities, and
// LookupDeveloperIdentity. AWS rejects values outside this range server-side
// with InvalidParameterException.
func validateQueryLimit(n int) bool {
	return n >= 1 && n <= 60
}

func validateRoleKeys(authRole, unauthRole string) bool {
	return authRole != "" || unauthRole != ""
}

func validateRoleMappings(mappings map[string]RoleMappingInput) bool {
	for key, rm := range mappings {
		if len(key) < 1 || len(key) > 128 {
			return false
		}
		if !validRoleMappingTypes[rm.Type] {
			return false
		}
		// Per Smithy RoleMapping documentation, AmbiguousRoleResolution is
		// required whenever Type is Token or Rules.
		if rm.AmbiguousRoleResolution == "" {
			return false
		}
		if !validAmbiguousRoleResolutions[rm.AmbiguousRoleResolution] {
			return false
		}
		if rm.Type == "Rules" && rm.RulesConfiguration == nil {
			return false
		}
		if rm.Type == "Rules" && rm.RulesConfiguration != nil {
			if len(rm.RulesConfiguration.Rules) == 0 {
				return false
			}
			if len(rm.RulesConfiguration.Rules) > 400 {
				return false
			}
			for _, rule := range rm.RulesConfiguration.Rules {
				if !validateMappingRule(rule) {
					return false
				}
			}
		}
	}
	return true
}

func validateMappingRule(rule MappingRuleInput) bool {
	if len(rule.Claim) < 1 || len(rule.Claim) > 64 {
		return false
	}
	if !claimNamePattern.MatchString(rule.Claim) {
		return false
	}
	if !validMappingRuleMatchTypes[rule.MatchType] {
		return false
	}
	if !validateClaimValue(rule.Value) {
		return false
	}
	if !validateRoleARN(rule.RoleARN) {
		return false
	}
	return true
}

// ---------------------------------------------------------------------------
// Smithy-derived scalar validators
// ---------------------------------------------------------------------------

// validateIdentityId enforces the Smithy IdentityId constraints: length 1-55
// and pattern ^[\w-]+:[0-9a-f-]+$ (region:guid format).
func validateIdentityId(id string) bool {
	if len(id) < 1 || len(id) > 55 {
		return false
	}
	return identityIdPattern.MatchString(id)
}

// validateIdentityPoolId enforces the Smithy IdentityPoolId constraints:
// length 1-55 and pattern ^[\w-]+:[0-9a-f-]+$.
func validateIdentityPoolId(id string) bool {
	if len(id) < 1 || len(id) > 55 {
		return false
	}
	return identityIdPattern.MatchString(id)
}

// validatePaginationKey enforces the Smithy PaginationKey constraints:
// length 1-65535 and pattern ^[\S]+$. An empty token is valid (first page).
func validatePaginationKey(token string) bool {
	if token == "" {
		return true
	}
	if len(token) > 65535 {
		return false
	}
	return paginationKeyPattern.MatchString(token)
}

// validateRoleARN enforces the Smithy ARNString constraints: length 20-2048.
func validateRoleARN(arn string) bool {
	return len(arn) >= 20 && len(arn) <= 2048
}

// validateAccountId enforces the Smithy AccountId constraints: length 1-15
// and pattern ^\d+$ (digits only).
func validateAccountId(id string) bool {
	if len(id) < 1 || len(id) > 15 {
		return false
	}
	return accountIdPattern.MatchString(id)
}

// validateTokenDuration enforces the Smithy TokenDuration range [1, 86400].
func validateTokenDuration(d int64) bool {
	return d >= 1 && d <= 86400
}

// validateClaimValue enforces the Smithy ClaimValue constraints: length 1-128.
func validateClaimValue(v string) bool {
	return len(v) >= 1 && len(v) <= 128
}

// validatePrincipalTagValue enforces the Smithy PrincipalTagValue constraints:
// length 1-256.
func validatePrincipalTagValue(v string) bool {
	return len(v) >= 1 && len(v) <= 256
}

// validateLoginTokenValue enforces the Smithy IdentityProviderToken
// constraints: length 1-50000.
func validateLoginTokenValue(v string) bool {
	return len(v) >= 1 && len(v) <= 50000
}

// validateMapSize returns true when the map size does not exceed the Smithy
// @length max constraint.
func validateMapSize(size, max int) bool {
	return size <= max
}

// validateLoginsValues checks that every value in a Logins map satisfies the
// Smithy IdentityProviderToken length constraint [1, 50000].
func validateLoginsValues(logins map[string]string) bool {
	for _, v := range logins {
		if !validateLoginTokenValue(v) {
			return false
		}
	}
	return true
}

// validateTagValues checks that every value in a tag map satisfies the Smithy
// TagValueType length constraint [0, 256].
func validateTagValues(tags map[string]string) bool {
	for _, v := range tags {
		if len(v) > 256 {
			return false
		}
	}
	return true
}

// validateDeveloperUserIdentifier enforces the Smithy DeveloperUserIdentifier
// @length(min=1, max=1024) constraint.
func validateDeveloperUserIdentifier(id string) bool {
	return len(id) >= 1 && len(id) <= 1024
}

// validateIdentityProviderNameLength enforces the Smithy
// IdentityProviderName @length(min=1, max=128) constraint used by
// GetPrincipalTagAttributeMap and SetPrincipalTagAttributeMap.
func validateIdentityProviderNameLength(name string) bool {
	return len(name) >= 1 && len(name) <= 128
}

// validateLoginsKeys checks that every key in a Logins map satisfies the
// Smithy IdentityProviderName @length(min=1, max=128) constraint.
func validateLoginsKeys(logins map[string]string) bool {
	for k := range logins {
		if len(k) < 1 || len(k) > 128 {
			return false
		}
	}
	return true
}
