package cognitoidentityprovider

import (
	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

func parsePasswordPolicyWithBase(req *request.ParsedRequest, base *cognitostore.PasswordPolicy) (*cognitostore.PasswordPolicy, error) {
	hasPolicy := false
	policy := &cognitostore.PasswordPolicy{}
	if base != nil {
		*policy = *base
	} else {
		*policy = cognitostore.DefaultPasswordPolicy()
	}

	if policiesMap, ok := req.Parameters["Policies"].(map[string]interface{}); ok {
		if ppMap, ok := policiesMap["PasswordPolicy"].(map[string]interface{}); ok {
			// Every member is read at its awsJson1_1 shape — a typed JSON
			// value under its PascalCase name. A present member of any
			// other form is a wire-shape violation and rejected outright.
			if val, ok := ppMap["MinimumLength"]; ok {
				n, numeric := parseJSONInt(val)
				if !numeric {
					return nil, ErrInvalidParameter
				}
				policy.MinimumLength = n
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
				b, isBool := val.(bool)
				if !isBool {
					return nil, ErrInvalidParameter
				}
				policy.RequireUppercase = b
				hasPolicy = true
			}
			if val, ok := ppMap["RequireLowercase"]; ok {
				b, isBool := val.(bool)
				if !isBool {
					return nil, ErrInvalidParameter
				}
				policy.RequireLowercase = b
				hasPolicy = true
			}
			if val, ok := ppMap["RequireNumbers"]; ok {
				b, isBool := val.(bool)
				if !isBool {
					return nil, ErrInvalidParameter
				}
				policy.RequireNumbers = b
				hasPolicy = true
			}
			if val, ok := ppMap["RequireSymbols"]; ok {
				b, isBool := val.(bool)
				if !isBool {
					return nil, ErrInvalidParameter
				}
				policy.RequireSymbols = b
				hasPolicy = true
			}
			if val, ok := ppMap["TemporaryPasswordValidityDays"]; ok {
				n, numeric := parseJSONInt(val)
				if !numeric {
					return nil, ErrInvalidParameter
				}
				policy.TemporaryPasswordValidityDays = n
				if err := validatePasswordPolicyRanges(policy); err != nil {
					return nil, err
				}
				hasPolicy = true
			}
			if val, ok := ppMap["PasswordHistorySize"]; ok {
				n, numeric := parseJSONInt(val)
				if !numeric {
					return nil, ErrInvalidParameter
				}
				policy.PasswordHistorySize = n
				if err := validatePasswordPolicyRanges(policy); err != nil {
					return nil, err
				}
				hasPolicy = true
			}
		}
	}

	if !hasPolicy {
		return nil, nil
	}
	return policy, nil
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

func parseSmsConfiguration(req *request.ParsedRequest, poolRegion string) (*cognitostore.SmsConfiguration, error) {
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
		if eums, ok := m["EumsSms"].(map[string]interface{}); ok {
			// The End User Messaging SMS configuration is the documented
			// alternative to SNS delivery: a pool carries SnsCallerArn or
			// this structure, never both, and the structure itself always
			// names the role Cognito assumes (CallerArn is required). Its
			// Region, when present, must be the pool's own Region.
			if config.SnsCallerArn != "" {
				return nil, ErrInvalidParameter
			}
			eumsCfg := &cognitostore.EumsSmsConfiguration{
				CallerArn:            getStringParam(eums, "CallerArn"),
				ExternalId:           getStringParam(eums, "ExternalId"),
				OriginationIdentity:  getStringParam(eums, "OriginationIdentity"),
				ConfigurationSetName: getStringParam(eums, "ConfigurationSetName"),
				InEntityId:           getStringParam(eums, "InEntityId"),
				InTemplateId:         getStringParam(eums, "InTemplateId"),
				Region:               getStringParam(eums, "Region"),
			}
			if eumsCfg.CallerArn == "" {
				return nil, ErrInvalidParameter
			}
			if eumsCfg.Region != "" && eumsCfg.Region != poolRegion {
				return nil, ErrInvalidParameter
			}
			config.EumsSms = eumsCfg
			hasConfig = true
		}
	}
	if !hasConfig {
		return nil, nil
	}
	return config, nil
}

// parseSignInPolicy reads the Policies.SignInPolicy member: the first
// authentication factors the pool permits under choice-based sign-in.
// A present list must hold only AuthFactorType values within the model's
// length bounds; SOFTWARE_TOKEN is excluded because the model documents
// it as unsupported as a first factor.
func parseSignInPolicy(req *request.ParsedRequest) (*cognitostore.SignInPolicy, error) {
	policiesMap, ok := req.Parameters["Policies"].(map[string]interface{})
	if !ok {
		return nil, nil
	}
	signInMap, ok := policiesMap["SignInPolicy"].(map[string]interface{})
	if !ok {
		return nil, nil
	}
	rawFactors, ok := signInMap["AllowedFirstAuthFactors"].([]interface{})
	if !ok {
		// Absence leaves the policy empty; a present member of the wrong
		// type is malformed and must not fail open to an all-factors
		// policy.
		if _, present := signInMap["AllowedFirstAuthFactors"]; present {
			return nil, ErrInvalidParameter
		}
		return &cognitostore.SignInPolicy{}, nil
	}
	factors := make([]string, 0, len(rawFactors))
	for _, v := range rawFactors {
		s, isString := v.(string)
		if !isString || !validateAuthFactor(s) {
			return nil, ErrInvalidParameter
		}
		factors = append(factors, s)
	}
	if len(factors) < cognitostore.MinSignInPolicyFirstAuthFactors ||
		len(factors) > cognitostore.MaxSignInPolicyFirstAuthFactors {
		return nil, ErrInvalidParameter
	}
	return &cognitostore.SignInPolicy{AllowedFirstAuthFactors: factors}, nil
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
	if v, ok := m["AdvancedSecurityAdditionalFlows"].(map[string]interface{}); ok {
		addOns.AdvancedSecurityAdditionalFlows = &cognitostore.AdvancedSecurityAdditionalFlows{
			CustomAuthMode: getStringParam(v, "CustomAuthMode"),
		}
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
	// Like every sibling parser, a member present but carrying no mechanisms
	// parses as "not set" — it must not overwrite the stored setting with an
	// empty one.
	if len(setting.RecoveryMechanisms) == 0 {
		return nil
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

// parseSchemaAttributes reads the create-only Schema member. A present list
// must sit inside the SchemaAttributesListType bounds and every member must
// be a well-formed attribute definition; an absent member leaves the schema
// empty (the standard attributes apply).
func parseSchemaAttributes(req *request.ParsedRequest) ([]cognitostore.SchemaAttributeType, error) {
	var result []cognitostore.SchemaAttributeType
	if rawList, ok := req.Parameters["Schema"].([]interface{}); ok {
		if len(rawList) < cognitostore.MinSchemaAttributesPerPool ||
			len(rawList) > cognitostore.MaxSchemaAttributesPerPool {
			return nil, ErrInvalidParameter
		}
		for _, item := range rawList {
			m, ok := item.(map[string]interface{})
			if !ok {
				return nil, ErrInvalidParameter
			}
			name := getStringParam(m, "Name")
			// A standard attribute definition starts from the documented
			// defaults so that members the request leaves unset keep the
			// properties Amazon Cognito reports in describe responses;
			// supplied members overwrite them below.
			sa := cognitostore.SchemaAttributeType{Name: name}
			if def, ok := standardSchemaAttributeDefault(name); ok {
				sa = def
				sa.Name = name
			}
			if v := getStringParam(m, "AttributeDataType"); v != "" {
				sa.AttributeDataType = v
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
	return result, nil
}
