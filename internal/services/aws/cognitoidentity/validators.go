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
	for _, rm := range mappings {
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
	if rule.Value == "" {
		return false
	}
	if len(rule.RoleARN) < 20 || len(rule.RoleARN) > 2048 {
		return false
	}
	return true
}
