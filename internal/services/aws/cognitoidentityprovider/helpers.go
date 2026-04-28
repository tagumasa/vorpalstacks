package cognitoidentityprovider

import (
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"github.com/google/uuid"
)

func getBoolParam(req *request.ParsedRequest, key string) bool {
	lowerKey := strings.ToLower(key[:1]) + key[1:]

	if v, ok := req.Parameters[lowerKey]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	if v, ok := req.Parameters[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}

	val := request.GetParamLowerFirst(req.Parameters, key)
	return val == "true" || val == "True" || val == "1"
}

func getIntParam(req *request.ParsedRequest, key string) int {
	if v, ok := req.Parameters[key]; ok {
		switch n := v.(type) {
		case int:
			return n
		case int64:
			return int(n)
		case float64:
			return int(n)
		case string:
			return parseInt(n)
		}
	}
	lowerKey := strings.ToLower(key[:1]) + key[1:]
	if v, ok := req.Parameters[lowerKey]; ok {
		switch n := v.(type) {
		case int:
			return n
		case int64:
			return int(n)
		case float64:
			return int(n)
		case string:
			return parseInt(n)
		}
	}
	return 0
}

func getIntParamOK(req *request.ParsedRequest, key string) (int, bool) {
	if v, ok := req.Parameters[key]; ok {
		switch n := v.(type) {
		case int:
			return n, true
		case int64:
			return int(n), true
		case float64:
			return int(n), true
		case string:
			if n != "" {
				return parseInt(n), true
			}
		}
	}
	lowerKey := strings.ToLower(key[:1]) + key[1:]
	if v, ok := req.Parameters[lowerKey]; ok {
		switch n := v.(type) {
		case int:
			return n, true
		case int64:
			return int(n), true
		case float64:
			return int(n), true
		case string:
			if n != "" {
				return parseInt(n), true
			}
		}
	}
	return 0, false
}

func getUserPoolID(req *request.ParsedRequest) string {
	return req.GetParam("UserPoolId")
}

func getUsername(req *request.ParsedRequest) string {
	return req.GetParam("Username")
}

func getGroupName(req *request.ParsedRequest) string {
	return req.GetParam("GroupName")
}

func getPassword(req *request.ParsedRequest) string {
	return req.GetParam("Password")
}

func getNewPassword(req *request.ParsedRequest) string {
	if v := req.GetParam("NewPassword"); v != "" {
		return v
	}
	return req.GetParam("ProposedPassword")
}

func getPreviousPassword(req *request.ParsedRequest) string {
	return req.GetParam("PreviousPassword")
}

func getAccessToken(req *request.ParsedRequest) string {
	return req.GetParam("AccessToken")
}

func getConfirmationCode(req *request.ParsedRequest) string {
	return req.GetParam("ConfirmationCode")
}

func getClientId(req *request.ParsedRequest) string {
	return req.GetParam("ClientId")
}

func parseUserAttributes(req *request.ParsedRequest) map[string]string {
	attributes := make(map[string]string)

	if attrs, ok := req.Parameters["UserAttributes"].([]interface{}); ok {
		for _, attr := range attrs {
			if m, ok := attr.(map[string]interface{}); ok {
				name, _ := m["Name"].(string)
				value, _ := m["Value"].(string)
				if name != "" {
					attributes[name] = value
				}
			}
		}
		return attributes
	}

	for i := 1; ; i++ {
		idx := strconv.Itoa(i)
		nameKey := "UserAttributes." + idx + ".Name"
		valueKey := "UserAttributes." + idx + ".Value"
		name := req.GetParam(nameKey)
		if name == "" {
			break
		}
		value := req.GetParam(valueKey)
		attributes[name] = value
	}
	return attributes
}

func parsePasswordPolicy(req *request.ParsedRequest) *cognitostore.PasswordPolicy {
	hasPolicy := false
	policy := &cognitostore.PasswordPolicy{
		MinimumLength:                 8,
		RequireUppercase:              true,
		RequireLowercase:              true,
		RequireNumbers:                true,
		RequireSymbols:                true,
		TemporaryPasswordValidityDays: 7,
	}

	if policiesMap, ok := req.Parameters["Policies"].(map[string]interface{}); ok {
		if ppMap, ok := policiesMap["PasswordPolicy"].(map[string]interface{}); ok {
			if val, ok := ppMap["MinimumLength"]; ok {
				switch v := val.(type) {
				case int:
					policy.MinimumLength = v
				case float64:
					policy.MinimumLength = int(v)
				case string:
					policy.MinimumLength = parseInt(v)
				}
				hasPolicy = true
			}
			if val, ok := ppMap["RequireUppercase"]; ok {
				if b, ok := val.(bool); ok {
					policy.RequireUppercase = b
				} else if s, ok := val.(string); ok {
					policy.RequireUppercase = strings.ToLower(s) == "true"
				}
				hasPolicy = true
			}
			if val, ok := ppMap["RequireLowercase"]; ok {
				if b, ok := val.(bool); ok {
					policy.RequireLowercase = b
				} else if s, ok := val.(string); ok {
					policy.RequireLowercase = strings.ToLower(s) == "true"
				}
				hasPolicy = true
			}
			if val, ok := ppMap["RequireNumbers"]; ok {
				if b, ok := val.(bool); ok {
					policy.RequireNumbers = b
				} else if s, ok := val.(string); ok {
					policy.RequireNumbers = strings.ToLower(s) == "true"
				}
				hasPolicy = true
			}
			if val, ok := ppMap["RequireSymbols"]; ok {
				if b, ok := val.(bool); ok {
					policy.RequireSymbols = b
				} else if s, ok := val.(string); ok {
					policy.RequireSymbols = strings.ToLower(s) == "true"
				}
				hasPolicy = true
			}
			if val, ok := ppMap["TemporaryPasswordValidityDays"]; ok {
				switch v := val.(type) {
				case int:
					policy.TemporaryPasswordValidityDays = v
				case float64:
					policy.TemporaryPasswordValidityDays = int(v)
				case string:
					policy.TemporaryPasswordValidityDays = parseInt(v)
				}
				hasPolicy = true
			}
		}
	}

	if val := req.GetParam("Policies.PasswordPolicy.MinimumLength"); val != "" {
		policy.MinimumLength = parseInt(val)
		hasPolicy = true
	}
	if val := req.GetParam("Policies.PasswordPolicy.RequireUppercase"); val != "" {
		policy.RequireUppercase = strings.ToLower(val) == "true"
		hasPolicy = true
	}
	if val := req.GetParam("Policies.PasswordPolicy.RequireLowercase"); val != "" {
		policy.RequireLowercase = strings.ToLower(val) == "true"
		hasPolicy = true
	}
	if val := req.GetParam("Policies.PasswordPolicy.RequireNumbers"); val != "" {
		policy.RequireNumbers = strings.ToLower(val) == "true"
		hasPolicy = true
	}
	if val := req.GetParam("Policies.PasswordPolicy.RequireSymbols"); val != "" {
		policy.RequireSymbols = strings.ToLower(val) == "true"
		hasPolicy = true
	}
	if val := req.GetParam("Policies.PasswordPolicy.TemporaryPasswordValidityDays"); val != "" {
		policy.TemporaryPasswordValidityDays = parseInt(val)
		hasPolicy = true
	}

	if !hasPolicy {
		return nil
	}
	return policy
}

func parseLambdaConfig(req *request.ParsedRequest) *cognitostore.LambdaConfig {
	hasConfig := false
	config := &cognitostore.LambdaConfig{}

	if lambdaConfigMap, ok := req.Parameters["LambdaConfig"].(map[string]interface{}); ok {
		if val, ok := lambdaConfigMap["PreSignUp"].(string); ok && val != "" {
			config.PreSignUp = val
			hasConfig = true
		}
		if val, ok := lambdaConfigMap["CustomMessage"].(string); ok && val != "" {
			config.CustomMessage = val
			hasConfig = true
		}
		if val, ok := lambdaConfigMap["PostConfirmation"].(string); ok && val != "" {
			config.PostConfirmation = val
			hasConfig = true
		}
		if val, ok := lambdaConfigMap["PreAuthentication"].(string); ok && val != "" {
			config.PreAuthentication = val
			hasConfig = true
		}
		if val, ok := lambdaConfigMap["PostAuthentication"].(string); ok && val != "" {
			config.PostAuthentication = val
			hasConfig = true
		}
		if val, ok := lambdaConfigMap["DefineAuthChallenge"].(string); ok && val != "" {
			config.DefineAuthChallenge = val
			hasConfig = true
		}
		if val, ok := lambdaConfigMap["CreateAuthChallenge"].(string); ok && val != "" {
			config.CreateAuthChallenge = val
			hasConfig = true
		}
		if val, ok := lambdaConfigMap["VerifyAuthChallengeResponse"].(string); ok && val != "" {
			config.VerifyAuthChallengeResponse = val
			hasConfig = true
		}
		if val, ok := lambdaConfigMap["PreTokenGeneration"].(string); ok && val != "" {
			config.PreTokenGeneration = val
			hasConfig = true
		}
		if val, ok := lambdaConfigMap["UserMigration"].(string); ok && val != "" {
			config.UserMigration = val
			hasConfig = true
		}
	}

	fields := []struct {
		param string
		field *string
	}{
		{"LambdaConfig.PreSignUp", &config.PreSignUp},
		{"LambdaConfig.CustomMessage", &config.CustomMessage},
		{"LambdaConfig.PostConfirmation", &config.PostConfirmation},
		{"LambdaConfig.PreAuthentication", &config.PreAuthentication},
		{"LambdaConfig.PostAuthentication", &config.PostAuthentication},
		{"LambdaConfig.DefineAuthChallenge", &config.DefineAuthChallenge},
		{"LambdaConfig.CreateAuthChallenge", &config.CreateAuthChallenge},
		{"LambdaConfig.VerifyAuthChallengeResponse", &config.VerifyAuthChallengeResponse},
		{"LambdaConfig.PreTokenGeneration", &config.PreTokenGeneration},
		{"LambdaConfig.UserMigration", &config.UserMigration},
	}
	for _, f := range fields {
		if val := req.GetParam(f.param); val != "" {
			*f.field = val
			hasConfig = true
		}
	}
	if !hasConfig {
		return nil
	}
	return config
}

func parseInt(s string) int {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	return result
}

func formatUserPool(pool *cognitostore.UserPool) map[string]interface{} {
	result := map[string]interface{}{
		"Id":                 pool.ID,
		"Name":               pool.Name,
		"Arn":                pool.Arn,
		"Status":             pool.Status,
		"CreationDate":       pool.CreationDate.Unix(),
		"LastModifiedDate":   pool.LastModifiedDate.Unix(),
		"MfaConfiguration":   pool.MfaConfiguration,
		"DeletionProtection": "INACTIVE",
	}

	if pool.MfaConfiguration == "" {
		result["MfaConfiguration"] = "OFF"
	}

	if pool.AliasAttributes != nil {
		result["AliasAttributes"] = pool.AliasAttributes
	}
	if pool.UsernameAttributes != nil {
		result["UsernameAttributes"] = pool.UsernameAttributes
	}
	if pool.Schema != "" {
		result["Schema"] = pool.Schema
	}
	if pool.MfaConfiguration != "" {
		result["MfaConfiguration"] = pool.MfaConfiguration
	}
	if pool.PasswordPolicy != nil {
		result["Policies"] = map[string]interface{}{
			"PasswordPolicy": map[string]interface{}{
				"MinimumLength":                 pool.PasswordPolicy.MinimumLength,
				"RequireUppercase":              pool.PasswordPolicy.RequireUppercase,
				"RequireLowercase":              pool.PasswordPolicy.RequireLowercase,
				"RequireNumbers":                pool.PasswordPolicy.RequireNumbers,
				"RequireSymbols":                pool.PasswordPolicy.RequireSymbols,
				"TemporaryPasswordValidityDays": pool.PasswordPolicy.TemporaryPasswordValidityDays,
			},
		}
	}
	if pool.LambdaConfig != nil {
		result["LambdaConfig"] = formatLambdaConfig(pool.LambdaConfig)
	}
	if pool.EstimatedNumberOfUsers > 0 {
		result["EstimatedNumberOfUsers"] = pool.EstimatedNumberOfUsers
	}
	if len(pool.Tags) > 0 {
		tags := make(map[string]string, len(pool.Tags))
		for _, t := range pool.Tags {
			tags[t.Key] = t.Value
		}
		result["UserPoolTags"] = tags
	}

	return result
}

func formatLambdaConfig(config *cognitostore.LambdaConfig) map[string]interface{} {
	result := make(map[string]interface{})
	fields := []struct {
		key   string
		value string
	}{
		{"PreSignUp", config.PreSignUp},
		{"CustomMessage", config.CustomMessage},
		{"PostConfirmation", config.PostConfirmation},
		{"PreAuthentication", config.PreAuthentication},
		{"PostAuthentication", config.PostAuthentication},
		{"DefineAuthChallenge", config.DefineAuthChallenge},
		{"CreateAuthChallenge", config.CreateAuthChallenge},
		{"VerifyAuthChallengeResponse", config.VerifyAuthChallengeResponse},
		{"PreTokenGeneration", config.PreTokenGeneration},
		{"UserMigration", config.UserMigration},
	}
	for _, f := range fields {
		if f.value != "" {
			result[f.key] = f.value
		}
	}
	return result
}

func formatUser(user *cognitostore.User) map[string]interface{} {
	result := map[string]interface{}{
		"Username":             user.Username,
		"UserCreateDate":       user.CreatedDate.Unix(),
		"UserLastModifiedDate": user.LastModifiedDate.Unix(),
		"Enabled":              user.Enabled,
		"UserStatus":           user.UserStatus,
	}

	if user.Attributes != nil {
		attrs := make([]map[string]string, 0)
		for name, value := range user.Attributes {
			attrs = append(attrs, map[string]string{
				"Name":  name,
				"Value": value,
			})
		}
		result["Attributes"] = attrs
	}

	if len(user.MFAOptions) > 0 {
		mfaOpts := make([]map[string]interface{}, 0)
		for _, opt := range user.MFAOptions {
			mfaOpts = append(mfaOpts, map[string]interface{}{
				"DeliveryMedium": opt.DeliveryMedium,
				"AttributeName":  opt.AttributeName,
			})
		}
		result["MFAOptions"] = mfaOpts
	}

	return result
}

func formatGroup(group *cognitostore.Group) map[string]interface{} {
	result := map[string]interface{}{
		"GroupName":        group.Name,
		"UserPoolId":       group.UserPoolID,
		"CreationDate":     group.CreationDate.Unix(),
		"LastModifiedDate": group.LastModifiedDate.Unix(),
	}

	if group.Description != "" {
		result["Description"] = group.Description
	}
	if group.RoleArn != "" {
		result["RoleArn"] = group.RoleArn
	}
	if group.Precedence != 0 {
		result["Precedence"] = group.Precedence
	}

	return result
}

func generateSessionID() string {
	return "SESSION_" + uuid.New().String()[:16]
}
