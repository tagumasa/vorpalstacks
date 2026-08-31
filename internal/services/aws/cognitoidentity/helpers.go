package cognitoidentity

import (
	"vorpalstacks/internal/common/request"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
)

func getBoolParam(req *request.ParsedRequest, key string) bool {
	return request.GetBoolParam(req.Parameters, key)
}

func getStringSliceParam(req *request.ParsedRequest, key string) []string {
	if val, ok := req.Parameters[key]; ok {
		if slice, ok := val.([]interface{}); ok {
			result := make([]string, 0, len(slice))
			for _, v := range slice {
				if s, ok := v.(string); ok {
					result = append(result, s)
				}
			}
			return result
		}
		if slice, ok := val.([]string); ok {
			return slice
		}
	}
	return nil
}

func formatLoginKeys(logins map[string]string) []string {
	keys := make([]string, 0, len(logins))
	for k := range logins {
		keys = append(keys, k)
	}
	return keys
}

func parseMapParam(req *request.ParsedRequest, key string) map[string]string {
	if val, ok := req.Parameters[key]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			return request.CopyStringMap(m)
		}
	}
	return nil
}

func formatRoleMappings(mappings map[string]cognitoidentitystore.RoleMapping) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range mappings {
		mapping := map[string]interface{}{
			"Type": v.Type,
		}
		if v.AmbiguousRoleResolution != "" {
			mapping["AmbiguousRoleResolution"] = v.AmbiguousRoleResolution
		}
		if v.RulesConfiguration != nil {
			mapping["RulesConfiguration"] = formatRulesConfiguration(v.RulesConfiguration)
		}
		result[k] = mapping
	}
	return result
}

func formatRulesConfiguration(config *cognitoidentitystore.RulesConfiguration) map[string]interface{} {
	rules := make([]map[string]interface{}, 0, len(config.Rules))
	for _, r := range config.Rules {
		rule := map[string]interface{}{
			"Claim":     r.Claim,
			"MatchType": r.MatchType,
			"Value":     r.Value,
			"RoleARN":   r.RoleARN,
		}
		rules = append(rules, rule)
	}
	return map[string]interface{}{
		"Rules": rules,
	}
}

// roleMappingInputToStore converts a service-layer RoleMappingInput DTO into
// the store-layer RoleMapping type required by SetIdentityPoolRoles.
func roleMappingInputToStore(in RoleMappingInput) cognitoidentitystore.RoleMapping {
	rm := cognitoidentitystore.RoleMapping{
		Type:                    in.Type,
		AmbiguousRoleResolution: in.AmbiguousRoleResolution,
	}
	if in.RulesConfiguration != nil {
		rules := make([]cognitoidentitystore.MappingRule, len(in.RulesConfiguration.Rules))
		for i, r := range in.RulesConfiguration.Rules {
			rules[i] = cognitoidentitystore.MappingRule{
				Claim:     r.Claim,
				MatchType: r.MatchType,
				Value:     r.Value,
				RoleARN:   r.RoleARN,
			}
		}
		rm.RulesConfiguration = &cognitoidentitystore.RulesConfiguration{Rules: rules}
	}
	return rm
}

// roleMappingMapToStore converts a map of RoleMappingInput DTOs into the
// store-layer map required by SetIdentityPoolRoles.
func roleMappingMapToStore(in map[string]RoleMappingInput) map[string]cognitoidentitystore.RoleMapping {
	if len(in) == 0 {
		return nil
	}
	result := make(map[string]cognitoidentitystore.RoleMapping, len(in))
	for k, v := range in {
		result[k] = roleMappingInputToStore(v)
	}
	return result
}

// providerOutsToStore converts a slice of service-layer ProviderOut DTOs into
// the store-layer CognitoIdentityProvider slice.
func providerOutsToStore(outs []ProviderOut) []cognitoidentitystore.CognitoIdentityProvider {
	providers := make([]cognitoidentitystore.CognitoIdentityProvider, 0, len(outs))
	for _, p := range outs {
		providers = append(providers, cognitoidentitystore.CognitoIdentityProvider{
			ProviderName:         p.ProviderName,
			ClientID:             p.ClientID,
			ServerSideTokenCheck: p.ServerSideTokenCheck,
		})
	}
	return providers
}
