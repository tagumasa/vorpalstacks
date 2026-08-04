package cognitoidentity

import (
	"fmt"
	"regexp"

	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
)

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

func validateIdentityPoolName(name string) error {
	if len(name) < 1 || len(name) > 128 {
		return fmt.Errorf("%w: IdentityPoolName must be 1-128 characters", ErrInvalidParameter)
	}
	if !identityPoolNamePattern.MatchString(name) {
		return fmt.Errorf("%w: IdentityPoolName contains invalid characters", ErrInvalidParameter)
	}
	return nil
}

func validateDeveloperProviderName(name string) error {
	if len(name) < 1 || len(name) > 128 {
		return fmt.Errorf("%w: DeveloperProviderName must be 1-128 characters", ErrInvalidParameter)
	}
	if !developerProviderNamePattern.MatchString(name) {
		return fmt.Errorf("%w: DeveloperProviderName contains invalid characters", ErrInvalidParameter)
	}
	return nil
}

func validateProviderName(name string) error {
	if len(name) < 1 || len(name) > 128 {
		return fmt.Errorf("%w: ProviderName must be 1-128 characters", ErrInvalidParameter)
	}
	if !providerNamePattern.MatchString(name) {
		return fmt.Errorf("%w: ProviderName contains invalid characters", ErrInvalidParameter)
	}
	return nil
}

func validateProviderClientId(id string) error {
	if len(id) < 1 || len(id) > 128 {
		return fmt.Errorf("%w: ClientId must be 1-128 characters", ErrInvalidParameter)
	}
	if !providerClientIdPattern.MatchString(id) {
		return fmt.Errorf("%w: ClientId contains invalid characters", ErrInvalidParameter)
	}
	return nil
}

// validateQueryLimit enforces the Smithy QueryLimit range (min=1, max=60)
// for MaxResults on ListIdentityPools, ListIdentities, and
// LookupDeveloperIdentity. AWS rejects values outside this range server-side
// with InvalidParameterException.
func validateQueryLimit(n int) error {
	if n < 1 || n > 60 {
		return fmt.Errorf("%w: MaxResults must be between 1 and 60", ErrInvalidParameter)
	}
	return nil
}

func validateRoleKeys(roles map[string]interface{}) error {
	if len(roles) > 2 {
		return fmt.Errorf("%w: Roles must not exceed 2 entries", ErrInvalidParameter)
	}
	for k := range roles {
		if !validRoleTypes[k] {
			return fmt.Errorf("%w: Role type must be one of [authenticated, unauthenticated]", ErrInvalidParameter)
		}
	}
	return nil
}

func validateRoleMappings(mappings map[string]cognitoidentitystore.RoleMapping) error {
	for _, rm := range mappings {
		if !validRoleMappingTypes[rm.Type] {
			return fmt.Errorf("%w: RoleMapping Type must be Token or Rules", ErrInvalidParameter)
		}
		// Per Smithy RoleMapping documentation, AmbiguousRoleResolution is
		// required whenever Type is Token or Rules.
		if rm.AmbiguousRoleResolution == "" {
			return fmt.Errorf("%w: AmbiguousRoleResolution is required when Type is Token or Rules", ErrInvalidParameter)
		}
		if !validAmbiguousRoleResolutions[rm.AmbiguousRoleResolution] {
			return fmt.Errorf("%w: AmbiguousRoleResolution must be AuthenticatedRole or Deny", ErrInvalidParameter)
		}
		if rm.Type == "Rules" && rm.RulesConfiguration == nil {
			return fmt.Errorf("%w: RulesConfiguration is required when Type is Rules", ErrInvalidParameter)
		}
		if rm.Type == "Rules" && rm.RulesConfiguration != nil {
			if len(rm.RulesConfiguration.Rules) == 0 {
				return fmt.Errorf("%w: RulesConfiguration must contain at least 1 rule", ErrInvalidParameter)
			}
			for _, rule := range rm.RulesConfiguration.Rules {
				if err := validateMappingRule(rule); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func validateMappingRule(rule cognitoidentitystore.MappingRule) error {
	if len(rule.Claim) < 1 || len(rule.Claim) > 64 {
		return fmt.Errorf("%w: Claim must be 1-64 characters", ErrInvalidParameter)
	}
	if !claimNamePattern.MatchString(rule.Claim) {
		return fmt.Errorf("%w: Claim contains invalid characters", ErrInvalidParameter)
	}
	if !validMappingRuleMatchTypes[rule.MatchType] {
		return fmt.Errorf("%w: MatchType must be Equals, Contains, StartsWith, or NotEqual", ErrInvalidParameter)
	}
	if rule.Value == "" {
		return fmt.Errorf("%w: Value is required", ErrInvalidParameter)
	}
	if len(rule.RoleARN) < 20 || len(rule.RoleARN) > 2048 {
		return fmt.Errorf("%w: RoleARN must be 20-2048 characters", ErrInvalidParameter)
	}
	return nil
}
