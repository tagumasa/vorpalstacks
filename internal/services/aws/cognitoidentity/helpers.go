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

func parseCognitoIdentityProviders(req *request.ParsedRequest) ([]cognitoidentitystore.CognitoIdentityProvider, error) {
	val, ok := req.Parameters["CognitoIdentityProviders"]
	if !ok {
		return nil, nil
	}
	slice, ok := val.([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	providers := make([]cognitoidentitystore.CognitoIdentityProvider, 0)
	for _, v := range slice {
		m, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}
		provider := cognitoidentitystore.CognitoIdentityProvider{}
		if name, ok := m["ProviderName"].(string); ok {
			if err := validateProviderName(name); err != nil {
				return nil, err
			}
			provider.ProviderName = name
		}
		if clientID, ok := m["ClientId"].(string); ok {
			if err := validateProviderClientId(clientID); err != nil {
				return nil, err
			}
			provider.ClientID = clientID
		}
		if check, ok := m["ServerSideTokenCheck"].(bool); ok {
			provider.ServerSideTokenCheck = check
		}
		providers = append(providers, provider)
	}
	return providers, nil
}

func parseMapParam(req *request.ParsedRequest, key string) map[string]string {
	if val, ok := req.Parameters[key]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			return request.CopyStringMap(m)
		}
	}
	return nil
}

func parseRoleMappings(req *request.ParsedRequest) (map[string]cognitoidentitystore.RoleMapping, error) {
	val, ok := req.Parameters["RoleMappings"]
	if !ok {
		return nil, nil
	}
	m, ok := val.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	result := make(map[string]cognitoidentitystore.RoleMapping)
	for k, v := range m {
		mapping, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}
		rm := cognitoidentitystore.RoleMapping{}
		if t, ok := mapping["Type"].(string); ok {
			rm.Type = t
		}
		if arr, ok := mapping["AmbiguousRoleResolution"].(string); ok {
			rm.AmbiguousRoleResolution = arr
		}
		if rules, ok := mapping["RulesConfiguration"].(map[string]interface{}); ok {
			rc, err := parseRulesConfiguration(rules)
			if err != nil {
				return nil, err
			}
			rm.RulesConfiguration = rc
		}
		result[k] = rm
	}
	if err := validateRoleMappings(result); err != nil {
		return nil, err
	}
	return result, nil
}

func parseRulesConfiguration(m map[string]interface{}) (*cognitoidentitystore.RulesConfiguration, error) {
	rules, ok := m["Rules"].([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	config := &cognitoidentitystore.RulesConfiguration{
		Rules: make([]cognitoidentitystore.MappingRule, 0),
	}
	for _, r := range rules {
		rule, ok := r.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}
		mr := cognitoidentitystore.MappingRule{}
		if claim, ok := rule["Claim"].(string); ok {
			mr.Claim = claim
		}
		if matchType, ok := rule["MatchType"].(string); ok {
			mr.MatchType = matchType
		}
		if value, ok := rule["Value"].(string); ok {
			mr.Value = value
		}
		if roleArn, ok := rule["RoleARN"].(string); ok {
			mr.RoleARN = roleArn
		}
		config.Rules = append(config.Rules, mr)
	}
	return config, nil
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
