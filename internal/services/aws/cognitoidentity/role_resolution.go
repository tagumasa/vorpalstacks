package cognitoidentity

import (
	"context"
	"errors"
	"maps"
	"slices"
	"strings"

	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
)

// role_resolution.go — authenticated-role selection for
// GetCredentialsForIdentity. The pool's RoleMappings govern the providers
// linked to the identity: Token mappings read the cognito:preferred_role and
// cognito:roles claims, Rules mappings match each rule's claim against the
// claim set (Equals/Contains/StartsWith/NotEqual per the model enum), and an
// empty or ambiguous result falls back to the mapping's
// AmbiguousRoleResolution (AuthenticatedRole or Deny). CustomRoleArn is
// documented as "the role to be assumed when multiple roles were received in
// the token", so it may only select among the roles the mapping grants;
// without a governing mapping the pool grants only its authenticated role.

const (
	preferredRoleClaim = "cognito:preferred_role"
	rolesClaim         = "cognito:roles"
)

// resolveAuthenticatedRole returns the role an authenticated identity may
// assume, or ErrNotAuthorized when the pool's configuration denies the
// request — a Deny ambiguous resolution, or a CustomRoleArn outside the
// multiple roles a mapping granted.
func resolveAuthenticatedRole(authRole string, mappings map[string]cognitoidentitystore.RoleMapping, linkedLogins map[string]string, claims map[string]string, customRoleARN string) (string, error) {
	// Providers are visited in sorted order: with several linked providers
	// the first governing mapping decides, and Go map iteration order must
	// not make that decision vary between requests.
	for _, provider := range slices.Sorted(maps.Keys(mappings)) {
		mapping := mappings[provider]
		if _, linked := linkedLogins[provider]; !linked {
			continue
		}
		granted := mappingGrantedRoles(mapping, claims)
		switch {
		case len(granted) == 1:
			// CustomRoleArn selects "when multiple roles were received in
			// the token": with a single granted role there is nothing to
			// select among and the parameter is ignored rather than
			// enforced.
			return granted[0], nil
		case len(granted) > 1 && customRoleARN != "":
			for _, role := range granted {
				if role == customRoleARN {
					return customRoleARN, nil
				}
			}
			return "", ErrNotAuthorized
		}
		// No role granted, or several granted without a CustomRoleArn to
		// pick one: the mapping's AmbiguousRoleResolution decides.
		if mapping.AmbiguousRoleResolution != "AuthenticatedRole" {
			return "", ErrNotAuthorized
		}
		return authRole, nil
	}

	// No mapping governs a linked provider: the pool grants only its
	// authenticated role — a single granted role, so CustomRoleArn is
	// ignored like every other single-role grant.
	return authRole, nil
}

// mappingGrantedRoles evaluates a role mapping against the claim set and
// returns every role it grants.
func mappingGrantedRoles(mapping cognitoidentitystore.RoleMapping, claims map[string]string) []string {
	switch mapping.Type {
	case "Token":
		if role := claims[preferredRoleClaim]; role != "" {
			return []string{role}
		}
		if list := claims[rolesClaim]; list != "" {
			roles := make([]string, 0, 4)
			for _, role := range strings.Split(list, ",") {
				if role = strings.TrimSpace(role); role != "" {
					roles = append(roles, role)
				}
			}
			return roles
		}
		return nil
	case "Rules":
		if mapping.RulesConfiguration == nil {
			return nil
		}
		for _, rule := range mapping.RulesConfiguration.Rules {
			if ruleMatchesClaim(claims, rule) {
				// "Rules are evaluated in order. The first one to match
				// specifies the role." — later matches do not add a second
				// granted role.
				return []string{rule.RoleARN}
			}
		}
		return nil
	}
	return nil
}

// ruleMatchesClaim applies the rule's match condition to the claim value.
func ruleMatchesClaim(claims map[string]string, rule cognitoidentitystore.MappingRule) bool {
	value, present := claims[rule.Claim]
	if !present {
		return false
	}
	switch rule.MatchType {
	case "Equals":
		return value == rule.Value
	case "Contains":
		return strings.Contains(value, rule.Value)
	case "StartsWith":
		return strings.HasPrefix(value, rule.Value)
	case "NotEqual":
		return value != rule.Value
	}
	return false
}

// tokenClaimsForLogins resolves the claim set role mappings evaluate: every
// linked login whose provider name is a platform user-pool issuer
// (cognito-idp.<region>.amazonaws.com/<poolID>) contributes the string
// claims of its validated ID token — cognito:roles, cognito:preferred_role
// and the custom: attributes Token and Rules mappings read. A presented
// platform token that fails validation fails the whole request closed. A
// provider name outside the issuer form belongs to an external identity
// provider, which the platform does not implement: it carries no verifiable
// token, contributes no claims, and the mapping's AmbiguousRoleResolution
// governs.
func (s *CognitoIdentityService) tokenClaimsForLogins(linkedLogins map[string]string) (map[string]string, error) {
	var claims map[string]string
	for _, provider := range slices.Sorted(maps.Keys(linkedLogins)) {
		region, poolID, ok := platformUserPoolProvider(provider)
		if !ok {
			continue
		}
		if s.idTokenClaims == nil {
			return nil, ErrInternalError
		}
		tokenClaims, err := s.idTokenClaims.IDTokenClaimsForPool(context.Background(), region, poolID, linkedLogins[provider])
		if err != nil {
			return nil, ErrNotAuthorized
		}
		if claims == nil {
			claims = make(map[string]string, len(tokenClaims))
		}
		for k, v := range tokenClaims {
			claims[k] = v
		}
	}
	return claims, nil
}

// principalTagClaims resolves the session tags the pool's principal-tag
// attribute maps attach to an authenticated identity's credentials. A
// custom mapping takes its tag value from the named claim of the provider
// token; a map with UseDefaults applies the default mappings, whose values
// are the aud and sub claims of a user-pool ID token (the app client ID
// and the user ID). Only platform user-pool providers carry verifiable
// claims on this platform — a map for any other provider contributes
// nothing under the external-IdP exclusion. A failing map read fails the
// request closed: silently under-tagging a session would change what its
// IAM policies evaluate.
func (s *CognitoIdentityService) principalTagClaims(store cognitoidentitystore.CognitoIdentityStoreInterface, poolID string, linkedLogins, claims map[string]string) (map[string]string, error) {
	var tags map[string]string
	apply := func(key, value string) {
		if value == "" {
			return
		}
		if tags == nil {
			tags = make(map[string]string)
		}
		tags[key] = value
	}
	for _, provider := range slices.Sorted(maps.Keys(linkedLogins)) {
		if _, _, ok := platformUserPoolProvider(provider); !ok {
			continue
		}
		ptam, err := store.GetPrincipalTagAttributeMap(poolID, provider)
		if err != nil {
			if errors.Is(err, cognitoidentitystore.ErrIdentityNotFound) {
				continue
			}
			return nil, err
		}
		if ptam.UseDefaults {
			apply("aud", claims["aud"])
			apply("sub", claims["sub"])
		}
		for tagKey, claimName := range ptam.PrincipalTags {
			apply(tagKey, claims[claimName])
		}
	}
	return tags, nil
}

// platformUserPoolProvider reports the region and user-pool ID a provider
// name identifies when it carries the platform issuer form
// cognito-idp.<region>.amazonaws.com/<poolID>, with or without the https://
// scheme. Any other name identifies an external identity provider.
func platformUserPoolProvider(provider string) (region, poolID string, ok bool) {
	name := strings.TrimPrefix(provider, "https://")
	const idpPrefix = "cognito-idp."
	const hostSuffix = ".amazonaws.com/"
	if !strings.HasPrefix(name, idpPrefix) {
		return "", "", false
	}
	rest := strings.TrimPrefix(name, idpPrefix)
	idx := strings.Index(rest, hostSuffix)
	if idx <= 0 || idx+len(hostSuffix) >= len(rest) {
		return "", "", false
	}
	return rest[:idx], rest[idx+len(hostSuffix):], true
}
