package cognitoidentityprovider

import (
	"strings"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

func parsePasswordPolicyWithBase(req *request.ParsedRequest, base *cognitostore.PasswordPolicy) (*cognitostore.PasswordPolicy, error) {
	hasPolicy := false
	policy := &cognitostore.PasswordPolicy{}
	if base != nil {
		*policy = *base
	} else {
		policy.MinimumLength = 8
		policy.RequireUppercase = true
		policy.RequireLowercase = true
		policy.RequireNumbers = true
		policy.RequireSymbols = true
		policy.TemporaryPasswordValidityDays = 7
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
				// The member was explicitly present, so zero is a rejected
				// out-of-range value, not the "unset" marker the
				// stored-policy check tolerates.
				if err := validateExplicitMinimumLength(policy.MinimumLength); err != nil {
					return nil, err
				}
				if err := validatePasswordPolicyRanges(policy); err != nil {
					return nil, err
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
				if err := validatePasswordPolicyRanges(policy); err != nil {
					return nil, err
				}
				hasPolicy = true
			}
			if val, ok := ppMap["PasswordHistorySize"]; ok {
				switch v := val.(type) {
				case int:
					policy.PasswordHistorySize = v
				case float64:
					policy.PasswordHistorySize = int(v)
				case string:
					policy.PasswordHistorySize = parseInt(v)
				}
				if err := validatePasswordPolicyRanges(policy); err != nil {
					return nil, err
				}
				hasPolicy = true
			}
		}
	}

	if val := req.GetParam("Policies.PasswordPolicy.MinimumLength"); val != "" {
		policy.MinimumLength = parseInt(val)
		// The member was explicitly present, so zero is a rejected
		// out-of-range value, not the "unset" marker the stored-policy
		// check tolerates.
		if err := validateExplicitMinimumLength(policy.MinimumLength); err != nil {
			return nil, err
		}
		if err := validatePasswordPolicyRanges(policy); err != nil {
			return nil, err
		}
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
	if val := req.GetParam("Policies.PasswordPolicy.PasswordHistorySize"); val != "" {
		policy.PasswordHistorySize = parseInt(val)
		if err := validatePasswordPolicyRanges(policy); err != nil {
			return nil, err
		}
		hasPolicy = true
	}
	if val := req.GetParam("Policies.PasswordPolicy.TemporaryPasswordValidityDays"); val != "" {
		policy.TemporaryPasswordValidityDays = parseInt(val)
		if err := validatePasswordPolicyRanges(policy); err != nil {
			return nil, err
		}
		hasPolicy = true
	}

	if !hasPolicy {
		return nil, nil
	}
	return policy, nil
}

func parseLambdaConfig(req *request.ParsedRequest) *cognitostore.LambdaConfig {
	return parseLambdaConfigWithBase(req, nil)
}

func parseLambdaConfigWithBase(req *request.ParsedRequest, base *cognitostore.LambdaConfig) *cognitostore.LambdaConfig {
	hasConfig := false
	config := &cognitostore.LambdaConfig{}
	if base != nil {
		*config = *base
	}

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
		if val, ok := lambdaConfigMap["KMSKeyID"].(string); ok && val != "" {
			config.KMSKeyID = val
			hasConfig = true
		}
		if m, ok := lambdaConfigMap["CustomEmailSender"].(map[string]interface{}); ok {
			config.CustomEmailSender = parseLambdaVersionConfig(m)
			hasConfig = true
		}
		if m, ok := lambdaConfigMap["CustomSMSSender"].(map[string]interface{}); ok {
			config.CustomSMSSender = parseLambdaVersionConfig(m)
			hasConfig = true
		}
		if m, ok := lambdaConfigMap["PreTokenGenerationConfig"].(map[string]interface{}); ok {
			config.PreTokenGenerationConfig = parseLambdaVersionConfig(m)
			hasConfig = true
		}
		if m, ok := lambdaConfigMap["InboundFederation"].(map[string]interface{}); ok {
			config.InboundFederation = parseLambdaVersionConfig(m)
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
		{"LambdaConfig.KMSKeyID", &config.KMSKeyID},
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

func parseLambdaVersionConfig(m map[string]interface{}) *cognitostore.LambdaVersionConfig {
	return &cognitostore.LambdaVersionConfig{
		LambdaArn:     getStringParam(m, "LambdaArn"),
		LambdaVersion: getStringParam(m, "LambdaVersion"),
	}
}

func parseEmailConfiguration(req *request.ParsedRequest) *cognitostore.EmailConfiguration {
	hasConfig := false
	config := &cognitostore.EmailConfiguration{}
	if m, ok := req.Parameters["EmailConfiguration"].(map[string]interface{}); ok {
		if v, ok := m["SourceArn"].(string); ok && v != "" {
			config.SourceArn = v
			hasConfig = true
		}
		if v, ok := m["ReplyToEmailAddress"].(string); ok && v != "" {
			config.ReplyToEmailAddress = v
			hasConfig = true
		}
		if v, ok := m["EmailSendingAccount"].(string); ok && v != "" {
			config.EmailSendingAccount = v
			hasConfig = true
		}
		if v, ok := m["From"].(string); ok && v != "" {
			config.From = v
			hasConfig = true
		}
		if v, ok := m["ConfigurationSet"].(string); ok && v != "" {
			config.ConfigurationSet = v
			hasConfig = true
		}
	}
	if !hasConfig {
		return nil
	}
	return config
}

func parseSmsConfiguration(req *request.ParsedRequest) *cognitostore.SmsConfiguration {
	hasConfig := false
	config := &cognitostore.SmsConfiguration{}
	if m, ok := req.Parameters["SmsConfiguration"].(map[string]interface{}); ok {
		if v, ok := m["SnsCallerArn"].(string); ok && v != "" {
			config.SnsCallerArn = v
			hasConfig = true
		}
		if v, ok := m["ExternalId"].(string); ok && v != "" {
			config.ExternalId = v
			hasConfig = true
		}
		if v, ok := m["SnsRegion"].(string); ok && v != "" {
			config.SnsRegion = v
			hasConfig = true
		}
	}
	if !hasConfig {
		return nil
	}
	return config
}

func parseAdminCreateUserConfig(req *request.ParsedRequest) *cognitostore.AdminCreateUserConfig {
	hasConfig := false
	config := &cognitostore.AdminCreateUserConfig{
		UnusedAccountValidityDays: 7,
	}
	if m, ok := req.Parameters["AdminCreateUserConfig"].(map[string]interface{}); ok {
		if v, ok := m["AllowAdminCreateUserOnly"].(bool); ok {
			config.AllowAdminCreateUserOnly = v
			hasConfig = true
		}
		if v, ok := m["UnusedAccountValidityDays"]; ok {
			switch n := v.(type) {
			case int:
				config.UnusedAccountValidityDays = n
			case float64:
				config.UnusedAccountValidityDays = int(n)
			}
			hasConfig = true
		}
		if tmplMap, ok := m["InviteMessageTemplate"].(map[string]interface{}); ok {
			tmpl := &cognitostore.MessageTemplate{}
			if v, ok := tmplMap["SMSMessage"].(string); ok {
				tmpl.SMSMessage = v
			}
			if v, ok := tmplMap["EmailMessage"].(string); ok {
				tmpl.EmailMessage = v
			}
			if v, ok := tmplMap["EmailSubject"].(string); ok {
				tmpl.EmailSubject = v
			}
			config.InviteMessageTemplate = tmpl
			hasConfig = true
		}
	}
	if !hasConfig {
		return nil
	}
	return config
}

func parseVerificationMessageTemplate(req *request.ParsedRequest) *cognitostore.VerificationMessageTemplate {
	hasConfig := false
	config := &cognitostore.VerificationMessageTemplate{}
	if m, ok := req.Parameters["VerificationMessageTemplate"].(map[string]interface{}); ok {
		if v, ok := m["SmsMessage"].(string); ok && v != "" {
			config.SmsMessage = v
			hasConfig = true
		}
		if v, ok := m["EmailMessage"].(string); ok && v != "" {
			config.EmailMessage = v
			hasConfig = true
		}
		if v, ok := m["EmailSubject"].(string); ok && v != "" {
			config.EmailSubject = v
			hasConfig = true
		}
		if v, ok := m["EmailMessageByLink"].(string); ok && v != "" {
			config.EmailMessageByLink = v
			hasConfig = true
		}
		if v, ok := m["EmailSubjectByLink"].(string); ok && v != "" {
			config.EmailSubjectByLink = v
			hasConfig = true
		}
		if v, ok := m["DefaultEmailOption"].(string); ok && v != "" {
			config.DefaultEmailOption = v
			hasConfig = true
		}
	}
	if !hasConfig {
		return nil
	}
	return config
}

func parseUserAttributeUpdateSettings(req *request.ParsedRequest) *cognitostore.UserAttributeUpdateSettings {
	m, ok := req.Parameters["UserAttributeUpdateSettings"].(map[string]interface{})
	if !ok {
		return nil
	}
	arr, ok := m["AttributesRequireVerificationBeforeUpdate"].([]interface{})
	if !ok {
		return nil
	}
	var attrs []string
	for _, v := range arr {
		if s, ok := v.(string); ok {
			attrs = append(attrs, s)
		}
	}
	if len(attrs) == 0 {
		return nil
	}
	return &cognitostore.UserAttributeUpdateSettings{
		AttributesRequireVerificationBeforeUpdate: attrs,
	}
}

func parseUserPoolAddOns(req *request.ParsedRequest) *cognitostore.UserPoolAddOns {
	m, ok := req.Parameters["UserPoolAddOns"].(map[string]interface{})
	if !ok {
		return nil
	}
	addOns := &cognitostore.UserPoolAddOns{}
	if v, ok := m["AdvancedSecurityMode"].(string); ok {
		addOns.AdvancedSecurityMode = v
	}
	return addOns
}

func parseAccountRecoverySetting(req *request.ParsedRequest) *cognitostore.AccountRecoverySetting {
	m, ok := req.Parameters["AccountRecoverySetting"].(map[string]interface{})
	if !ok {
		return nil
	}
	setting := &cognitostore.AccountRecoverySetting{}
	if mechs, ok := m["RecoveryMechanisms"].([]interface{}); ok {
		for _, mech := range mechs {
			if mm, ok := mech.(map[string]interface{}); ok {
				rm := cognitostore.RecoveryMechanism{}
				if v, ok := mm["Priority"].(float64); ok {
					rm.Priority = int(v)
				}
				if v, ok := mm["Name"].(string); ok {
					rm.Name = v
				}
				setting.RecoveryMechanisms = append(setting.RecoveryMechanisms, rm)
			}
		}
	}
	return setting
}

func parseUsernameConfiguration(req *request.ParsedRequest) *cognitostore.UsernameConfiguration {
	m, ok := req.Parameters["UsernameConfiguration"].(map[string]interface{})
	if !ok {
		return nil
	}
	cfg := &cognitostore.UsernameConfiguration{}
	if v, ok := m["CaseSensitive"].(bool); ok {
		cfg.CaseSensitive = v
	}
	return cfg
}

func parseDeviceConfiguration(req *request.ParsedRequest) *cognitostore.DeviceConfiguration {
	m, ok := req.Parameters["DeviceConfiguration"].(map[string]interface{})
	if !ok {
		return nil
	}
	cfg := &cognitostore.DeviceConfiguration{}
	if v, ok := m["ChallengeRequiredOnNewDevice"].(bool); ok {
		cfg.ChallengeRequiredOnNewDevice = v
	}
	if v, ok := m["DeviceOnlyRememberedOnUserPrompt"].(bool); ok {
		cfg.DeviceOnlyRememberedOnUserPrompt = v
	}
	return cfg
}

func parseSchemaAttributes(req *request.ParsedRequest) []cognitostore.SchemaAttributeType {
	var result []cognitostore.SchemaAttributeType
	if rawList, ok := req.Parameters["Schema"].([]interface{}); ok {
		for _, item := range rawList {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			sa := cognitostore.SchemaAttributeType{
				Name:              getStringParam(m, "Name"),
				AttributeDataType: getStringParam(m, "AttributeDataType"),
			}
			if v, ok := m["DeveloperOnlyAttribute"].(bool); ok {
				sa.DeveloperOnlyAttribute = v
			}
			if v, ok := m["Mutable"].(bool); ok {
				sa.Mutable = v
			}
			if v, ok := m["Required"].(bool); ok {
				sa.Required = v
			}
			if nac, ok := m["NumberAttributeConstraints"].(map[string]interface{}); ok {
				sa.NumberAttributeConstraints = &cognitostore.NumberAttributeConstraints{
					MinValue: getStringParam(nac, "MinValue"),
					MaxValue: getStringParam(nac, "MaxValue"),
				}
			}
			if sac, ok := m["StringAttributeConstraints"].(map[string]interface{}); ok {
				sa.StringAttributeConstraints = &cognitostore.StringAttributeConstraints{
					MinLength: getStringParam(sac, "MinLength"),
					MaxLength: getStringParam(sac, "MaxLength"),
				}
			}
			if sa.Name != "" {
				result = append(result, sa)
			}
		}
	}
	return result
}
