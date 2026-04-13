package cognitoidentity

import (
	"vorpalstacks/internal/common/request"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
)

func getParam(req *request.ParsedRequest, key string) string {
	return request.GetParamLowerFirst(req.Parameters, key)
}

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

func parseCognitoIdentityProviders(req *request.ParsedRequest) []cognitoidentitystore.CognitoIdentityProvider {
	if val, ok := req.Parameters["CognitoIdentityProviders"]; ok {
		if slice, ok := val.([]interface{}); ok {
			providers := make([]cognitoidentitystore.CognitoIdentityProvider, 0)
			for _, v := range slice {
				if m, ok := v.(map[string]interface{}); ok {
					provider := cognitoidentitystore.CognitoIdentityProvider{}
					if name, ok := m["ProviderName"].(string); ok {
						provider.ProviderName = name
					}
					if clientID, ok := m["ClientId"].(string); ok {
						provider.ClientID = clientID
					}
					if check, ok := m["ServerSideTokenCheck"].(bool); ok {
						provider.ServerSideTokenCheck = check
					}
					providers = append(providers, provider)
				}
			}
			return providers
		}
	}
	return nil
}

func parseSupportedLoginProviders(req *request.ParsedRequest) map[string]string {
	if val, ok := req.Parameters["SupportedLoginProviders"]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			result := make(map[string]string)
			for k, v := range m {
				if s, ok := v.(string); ok {
					result[k] = s
				}
			}
			return result
		}
	}
	return nil
}

func parseMapParam(req *request.ParsedRequest, key string) map[string]string {
	if val, ok := req.Parameters[key]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			result := make(map[string]string)
			for k, v := range m {
				if s, ok := v.(string); ok {
					result[k] = s
				}
			}
			return result
		}
	}
	return nil
}

func parseRoleMappings(req *request.ParsedRequest) map[string]cognitoidentitystore.RoleMapping {
	if val, ok := req.Parameters["RoleMappings"]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			result := make(map[string]cognitoidentitystore.RoleMapping)
			for k, v := range m {
				if mapping, ok := v.(map[string]interface{}); ok {
					rm := cognitoidentitystore.RoleMapping{}
					if t, ok := mapping["Type"].(string); ok {
						rm.Type = t
					}
					if arr, ok := mapping["AmbiguousRoleResolution"].(string); ok {
						rm.AmbiguousRoleResolution = arr
					}
					if rules, ok := mapping["RulesConfiguration"].(map[string]interface{}); ok {
						rm.RulesConfiguration = parseRulesConfiguration(rules)
					}
					result[k] = rm
				}
			}
			return result
		}
	}
	return nil
}

func parseRulesConfiguration(m map[string]interface{}) *cognitoidentitystore.RulesConfiguration {
	if rules, ok := m["Rules"].([]interface{}); ok {
		config := &cognitoidentitystore.RulesConfiguration{
			Rules: make([]cognitoidentitystore.MappingRule, 0),
		}
		for _, r := range rules {
			if rule, ok := r.(map[string]interface{}); ok {
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
		}
		return config
	}
	return nil
}

func formatIdentityPool(pool *cognitoidentitystore.IdentityPool) map[string]interface{} {
	result := map[string]interface{}{
		"IdentityPoolId":                 pool.ID,
		"IdentityPoolName":               pool.Name,
		"AllowUnauthenticatedIdentities": pool.AllowUnauthenticatedIdentities,
		"Arn":                            pool.Arn,
		"CreationDate":                   pool.CreationDate.Unix(),
		"LastModifiedDate":               pool.LastModifiedDate.Unix(),
	}

	if pool.AllowClassicFlow {
		result["AllowClassicFlow"] = pool.AllowClassicFlow
	}
	if len(pool.CognitoIdentityProviders) > 0 {
		result["CognitoIdentityProviders"] = formatCognitoIdentityProviders(pool.CognitoIdentityProviders)
	}
	if pool.DeveloperProviderName != "" {
		result["DeveloperProviderName"] = pool.DeveloperProviderName
	}
	if len(pool.SupportedLoginProviders) > 0 {
		result["SupportedLoginProviders"] = pool.SupportedLoginProviders
	}
	if len(pool.OpenIdConnectProviderARNs) > 0 {
		result["OpenIdConnectProviderARNs"] = pool.OpenIdConnectProviderARNs
	}
	if len(pool.SamlProviderARNs) > 0 {
		result["SamlProviderARNs"] = pool.SamlProviderARNs
	}

	return result
}

func formatIdentityPoolWithTags(pool *cognitoidentitystore.IdentityPool, tags map[string]string) map[string]interface{} {
	result := formatIdentityPool(pool)
	if len(tags) > 0 {
		result["IdentityPoolTags"] = tags
	}
	return result
}

func formatCognitoIdentityProviders(providers []cognitoidentitystore.CognitoIdentityProvider) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(providers))
	for _, p := range providers {
		provider := map[string]interface{}{
			"ProviderName": p.ProviderName,
		}
		if p.ClientID != "" {
			provider["ClientId"] = p.ClientID
		}
		if p.ServerSideTokenCheck {
			provider["ServerSideTokenCheck"] = p.ServerSideTokenCheck
		}
		result = append(result, provider)
	}
	return result
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
