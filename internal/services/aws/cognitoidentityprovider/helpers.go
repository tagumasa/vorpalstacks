package cognitoidentityprovider

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"github.com/google/uuid"
)

// cognitoIdpHost returns the Cognito IDP hostname for the given region,
// accounting for partition-specific suffixes (aws-cn uses .amazonaws.com.cn).
func cognitoIdpHost(region string) string {
	if strings.HasPrefix(region, "cn-") {
		return "cognito-idp." + region + ".amazonaws.com.cn"
	}
	return "cognito-idp." + region + ".amazonaws.com"
}

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
	v, _ := parseIntParam(req, key)
	return v
}

func getIntParamOK(req *request.ParsedRequest, key string) (int, bool) {
	return parseIntParam(req, key)
}

func parseIntParam(req *request.ParsedRequest, key string) (int, bool) {
	tryKey := func(k string) (int, bool) {
		if v, ok := req.Parameters[k]; ok {
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
	if v, ok := tryKey(key); ok {
		return v, true
	}
	lowerKey := strings.ToLower(key[:1]) + key[1:]
	return tryKey(lowerKey)
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
	return parseNamedAttributeList(req, "UserAttributes")
}

// parseValidationData extracts ValidationData from the request (CIPD-8).
func parseValidationData(req *request.ParsedRequest) map[string]string {
	return parseNamedAttributeList(req, "ValidationData")
}

func parseNamedAttributeList(req *request.ParsedRequest, key string) map[string]string {
	attributes := make(map[string]string)

	if attrs, ok := req.Parameters[key].([]interface{}); ok {
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
		nameKey := key + "." + idx + ".Name"
		valueKey := key + "." + idx + ".Value"
		name := req.GetParam(nameKey)
		if name == "" {
			break
		}
		value := req.GetParam(valueKey)
		attributes[name] = value
	}
	return attributes
}

// parseClientMetadata extracts ClientMetadata from the request (CIPD-8).
func parseClientMetadata(req *request.ParsedRequest) map[string]string {
	metadata := make(map[string]string)
	if cm, ok := req.Parameters["ClientMetadata"].(map[string]interface{}); ok {
		for k, v := range cm {
			if vs, ok := v.(string); ok {
				metadata[k] = vs
			}
		}
	}
	for k := range req.Parameters {
		if strings.HasPrefix(k, "ClientMetadata.") {
			attrKey := strings.TrimPrefix(k, "ClientMetadata.")
			metadata[attrKey] = req.GetParam(k)
		}
	}
	return metadata
}

func parsePasswordPolicyWithBase(req *request.ParsedRequest, base *cognitostore.PasswordPolicy) *cognitostore.PasswordPolicy {
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

func parsePasswordPolicy(req *request.ParsedRequest) *cognitostore.PasswordPolicy {
	return parsePasswordPolicyWithBase(req, nil)
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

func parseInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	var result int
	for _, c := range s {
		if c < '0' || c > '9' {
			return -1
		}
		result = result*10 + int(c-'0')
	}
	return result
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

func applyUserPoolUpdates(pool *cognitostore.UserPool, req *request.ParsedRequest) {
	if v := req.GetParam("PoolName"); v != "" {
		pool.Name = v
	}
	if v := req.GetParam("MfaConfiguration"); v != "" {
		pool.MfaConfiguration = v
	}
	if v := req.GetParam("DeletionProtection"); v != "" {
		pool.DeletionProtection = v
	}
	if v := req.GetParam("EmailVerificationMessage"); v != "" {
		pool.EmailVerificationMessage = v
	}
	if v := req.GetParam("EmailVerificationSubject"); v != "" {
		pool.EmailVerificationSubject = v
	}
	if v := req.GetParam("SmsVerificationMessage"); v != "" {
		pool.SmsVerificationMessage = v
	}
	if v := req.GetParam("SmsAuthenticationMessage"); v != "" {
		pool.SmsAuthenticationMessage = v
	}
	if v := request.GetStringList(req.Parameters, "AliasAttributes"); v != nil {
		pool.AliasAttributes = v
	}
	if v := request.GetStringList(req.Parameters, "UsernameAttributes"); v != nil {
		pool.UsernameAttributes = v
	}
	if v := request.GetStringList(req.Parameters, "AutoVerifiedAttributes"); v != nil {
		pool.AutoVerifiedAttributes = v
	}
	if schemaAttrs := parseSchemaAttributes(req); len(schemaAttrs) > 0 {
		pool.SchemaAttributes = schemaAttrs
	}
	if v := parsePasswordPolicyWithBase(req, pool.PasswordPolicy); v != nil {
		pool.PasswordPolicy = v
	}
	if v := parseLambdaConfigWithBase(req, pool.LambdaConfig); v != nil {
		pool.LambdaConfig = v
	}
	if v := parseEmailConfiguration(req); v != nil {
		pool.EmailConfiguration = v
	}
	if v := parseSmsConfiguration(req); v != nil {
		pool.SmsConfiguration = v
	}
	if v := parseAdminCreateUserConfig(req); v != nil {
		pool.AdminCreateUserConfig = v
	}
	if v := parseVerificationMessageTemplate(req); v != nil {
		pool.VerificationMessageTemplate = v
	}
	if v := parseUserAttributeUpdateSettings(req); v != nil {
		pool.UserAttributeUpdateSettings = v
	}
	if v := parseUserPoolAddOns(req); v != nil {
		pool.UserPoolAddOns = v
	}
	if v := parseAccountRecoverySetting(req); v != nil {
		pool.AccountRecoverySetting = v
	}
	if v := parseUsernameConfiguration(req); v != nil {
		pool.UsernameConfiguration = v
	}
	if v := parseDeviceConfiguration(req); v != nil {
		pool.DeviceConfiguration = v
	}
	if m, ok := req.Parameters["IssuerConfiguration"].(map[string]interface{}); ok {
		pool.IssuerConfiguration = &cognitostore.IssuerConfiguration{
			Type: getStringParam(m, "Type"),
		}
	}
	if m, ok := req.Parameters["KeyConfiguration"].(map[string]interface{}); ok {
		pool.KeyConfiguration = &cognitostore.KeyConfiguration{
			KeyType:   getStringParam(m, "KeyType"),
			KmsKeyArn: getStringParam(m, "KmsKeyArn"),
		}
	}
	if v := req.GetParam("UserPoolTier"); v != "" {
		pool.UserPoolTier = v
	}
}

func formatUserPool(pool *cognitostore.UserPool) map[string]interface{} {
	mfaConfig := pool.MfaConfiguration
	if mfaConfig == "" {
		mfaConfig = "OFF"
	}
	result := map[string]interface{}{
		"Id":                 pool.ID,
		"Name":               pool.Name,
		"Arn":                pool.Arn,
		"Status":             pool.Status,
		"CreationDate":       pool.CreationDate.Unix(),
		"LastModifiedDate":   pool.LastModifiedDate.Unix(),
		"MfaConfiguration":   mfaConfig,
		"DeletionProtection": pool.DeletionProtection,
	}
	if pool.DeletionProtection == "" {
		result["DeletionProtection"] = "INACTIVE"
	}

	if pool.AliasAttributes != nil {
		result["AliasAttributes"] = pool.AliasAttributes
	}
	if pool.UsernameAttributes != nil {
		result["UsernameAttributes"] = pool.UsernameAttributes
	}
	if pool.AutoVerifiedAttributes != nil {
		result["AutoVerifiedAttributes"] = pool.AutoVerifiedAttributes
	}
	if len(pool.SchemaAttributes) > 0 {
		schema := make([]map[string]interface{}, 0, len(pool.SchemaAttributes))
		for _, sa := range pool.SchemaAttributes {
			entry := map[string]interface{}{
				"Name":                  sa.Name,
				"AttributeDataType":     sa.AttributeDataType,
				"DeveloperOnlyAttribute": sa.DeveloperOnlyAttribute,
				"Mutable":               sa.Mutable,
				"Required":              sa.Required,
			}
			if sa.NumberAttributeConstraints != nil {
				entry["NumberAttributeConstraints"] = map[string]interface{}{
					"MinValue": sa.NumberAttributeConstraints.MinValue,
					"MaxValue": sa.NumberAttributeConstraints.MaxValue,
				}
			}
			if sa.StringAttributeConstraints != nil {
				entry["StringAttributeConstraints"] = map[string]interface{}{
					"MinLength": sa.StringAttributeConstraints.MinLength,
					"MaxLength": sa.StringAttributeConstraints.MaxLength,
				}
			}
			schema = append(schema, entry)
		}
		result["Schema"] = schema
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
	if pool.EmailConfiguration != nil {
		result["EmailConfiguration"] = formatEmailConfiguration(pool.EmailConfiguration)
	}
	if pool.SmsConfiguration != nil {
		result["SmsConfiguration"] = formatSmsConfiguration(pool.SmsConfiguration)
	}
	if pool.AdminCreateUserConfig != nil {
		result["AdminCreateUserConfig"] = formatAdminCreateUserConfig(pool.AdminCreateUserConfig)
	}
	if pool.VerificationMessageTemplate != nil {
		result["VerificationMessageTemplate"] = formatVerificationMessageTemplate(pool.VerificationMessageTemplate)
	}
	if pool.EmailVerificationMessage != "" {
		result["EmailVerificationMessage"] = pool.EmailVerificationMessage
	}
	if pool.EmailVerificationSubject != "" {
		result["EmailVerificationSubject"] = pool.EmailVerificationSubject
	}
	if pool.SmsVerificationMessage != "" {
		result["SmsVerificationMessage"] = pool.SmsVerificationMessage
	}
	if pool.SmsAuthenticationMessage != "" {
		result["SmsAuthenticationMessage"] = pool.SmsAuthenticationMessage
	}
	if pool.UserAttributeUpdateSettings != nil {
		result["UserAttributeUpdateSettings"] = map[string]interface{}{
			"AttributesRequireVerificationBeforeUpdate": pool.UserAttributeUpdateSettings.AttributesRequireVerificationBeforeUpdate,
		}
	}
	if pool.UserPoolAddOns != nil {
		result["UserPoolAddOns"] = map[string]interface{}{
			"AdvancedSecurityMode": pool.UserPoolAddOns.AdvancedSecurityMode,
		}
	}
	if pool.AccountRecoverySetting != nil && len(pool.AccountRecoverySetting.RecoveryMechanisms) > 0 {
		mechanisms := make([]map[string]interface{}, 0, len(pool.AccountRecoverySetting.RecoveryMechanisms))
		for _, rm := range pool.AccountRecoverySetting.RecoveryMechanisms {
			mechanisms = append(mechanisms, map[string]interface{}{
				"Priority": rm.Priority,
				"Name":     rm.Name,
			})
		}
		result["AccountRecoverySetting"] = map[string]interface{}{
			"RecoveryMechanisms": mechanisms,
		}
	}
	if pool.UsernameConfiguration != nil {
		result["UsernameConfiguration"] = map[string]interface{}{
			"CaseSensitive": pool.UsernameConfiguration.CaseSensitive,
		}
	}
	if pool.DeviceConfiguration != nil {
		result["DeviceConfiguration"] = map[string]interface{}{
			"ChallengeRequiredOnNewDevice":     pool.DeviceConfiguration.ChallengeRequiredOnNewDevice,
			"DeviceOnlyRememberedOnUserPrompt": pool.DeviceConfiguration.DeviceOnlyRememberedOnUserPrompt,
		}
	}
	if pool.IssuerConfiguration != nil {
		result["IssuerConfiguration"] = map[string]interface{}{
			"Type": pool.IssuerConfiguration.Type,
		}
	}
	if pool.KeyConfiguration != nil {
		kc := map[string]interface{}{}
		if pool.KeyConfiguration.KeyType != "" {
			kc["KeyType"] = pool.KeyConfiguration.KeyType
		}
		if pool.KeyConfiguration.KmsKeyArn != "" {
			kc["KmsKeyArn"] = pool.KeyConfiguration.KmsKeyArn
		}
		result["KeyConfiguration"] = kc
	}
	if pool.UserPoolTier != "" {
		result["UserPoolTier"] = pool.UserPoolTier
	}

	return result
}

func formatEmailConfiguration(config *cognitostore.EmailConfiguration) map[string]interface{} {
	result := make(map[string]interface{})
	if config.SourceArn != "" {
		result["SourceArn"] = config.SourceArn
	}
	if config.ReplyToEmailAddress != "" {
		result["ReplyToEmailAddress"] = config.ReplyToEmailAddress
	}
	if config.EmailSendingAccount != "" {
		result["EmailSendingAccount"] = config.EmailSendingAccount
	}
	if config.From != "" {
		result["From"] = config.From
	}
	if config.ConfigurationSet != "" {
		result["ConfigurationSet"] = config.ConfigurationSet
	}
	return result
}

func formatSmsConfiguration(config *cognitostore.SmsConfiguration) map[string]interface{} {
	result := make(map[string]interface{})
	if config.SnsCallerArn != "" {
		result["SnsCallerArn"] = config.SnsCallerArn
	}
	if config.ExternalId != "" {
		result["ExternalId"] = config.ExternalId
	}
	if config.SnsRegion != "" {
		result["SnsRegion"] = config.SnsRegion
	}
	return result
}

func formatAdminCreateUserConfig(config *cognitostore.AdminCreateUserConfig) map[string]interface{} {
	result := map[string]interface{}{
		"AllowAdminCreateUserOnly":  config.AllowAdminCreateUserOnly,
		"UnusedAccountValidityDays": config.UnusedAccountValidityDays,
	}
	if config.InviteMessageTemplate != nil {
		tmpl := make(map[string]interface{})
		if config.InviteMessageTemplate.SMSMessage != "" {
			tmpl["SMSMessage"] = config.InviteMessageTemplate.SMSMessage
		}
		if config.InviteMessageTemplate.EmailMessage != "" {
			tmpl["EmailMessage"] = config.InviteMessageTemplate.EmailMessage
		}
		if config.InviteMessageTemplate.EmailSubject != "" {
			tmpl["EmailSubject"] = config.InviteMessageTemplate.EmailSubject
		}
		result["InviteMessageTemplate"] = tmpl
	}
	return result
}

func formatVerificationMessageTemplate(config *cognitostore.VerificationMessageTemplate) map[string]interface{} {
	result := make(map[string]interface{})
	if config.SmsMessage != "" {
		result["SmsMessage"] = config.SmsMessage
	}
	if config.EmailMessage != "" {
		result["EmailMessage"] = config.EmailMessage
	}
	if config.EmailSubject != "" {
		result["EmailSubject"] = config.EmailSubject
	}
	if config.EmailMessageByLink != "" {
		result["EmailMessageByLink"] = config.EmailMessageByLink
	}
	if config.EmailSubjectByLink != "" {
		result["EmailSubjectByLink"] = config.EmailSubjectByLink
	}
	if config.DefaultEmailOption != "" {
		result["DefaultEmailOption"] = config.DefaultEmailOption
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
		{"KMSKeyID", config.KMSKeyID},
	}
	for _, f := range fields {
		if f.value != "" {
			result[f.key] = f.value
		}
	}
	if config.CustomEmailSender != nil {
		result["CustomEmailSender"] = map[string]interface{}{
			"LambdaArn":     config.CustomEmailSender.LambdaArn,
			"LambdaVersion": config.CustomEmailSender.LambdaVersion,
		}
	}
	if config.CustomSMSSender != nil {
		result["CustomSMSSender"] = map[string]interface{}{
			"LambdaArn":     config.CustomSMSSender.LambdaArn,
			"LambdaVersion": config.CustomSMSSender.LambdaVersion,
		}
	}
	if config.PreTokenGenerationConfig != nil {
		result["PreTokenGenerationConfig"] = map[string]interface{}{
			"LambdaArn":     config.PreTokenGenerationConfig.LambdaArn,
			"LambdaVersion": config.PreTokenGenerationConfig.LambdaVersion,
		}
	}
	if config.InboundFederation != nil {
		result["InboundFederation"] = map[string]interface{}{
			"LambdaArn":     config.InboundFederation.LambdaArn,
			"LambdaVersion": config.InboundFederation.LambdaVersion,
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
	if group.Precedence != nil {
		result["Precedence"] = *group.Precedence
	}

	return result
}

func generateSessionID() string {
	return "SESSION_" + uuid.New().String()
}

var userFilterRe = regexp.MustCompile(`^"?(\w[\w:.\-+]*)\s*(=|\^=)\s*"?([^"]+?)"?\s*$`)

func matchUserFilter(user *cognitostore.User, filter string) bool {
	f := strings.TrimSpace(filter)
	m := userFilterRe.FindStringSubmatch(f)
	if m == nil {
		return false
	}
	attrName := m[1]
	op := m[2]
	attrValue := m[3]

	var actual string
	switch strings.ToLower(attrName) {
	case "username":
		actual = user.Username
	case "cognito:user_status":
		actual = user.UserStatus
	case "status":
		if user.Enabled {
			actual = "Enabled"
		} else {
			actual = "Disabled"
		}
	default:
		if user.Attributes != nil {
			actual = user.Attributes[attrName]
		}
	}

	switch op {
	case "=":
		return strings.EqualFold(actual, attrValue)
	case "^=":
		return strings.HasPrefix(strings.ToLower(actual), strings.ToLower(attrValue))
	}
	return false
}

// computeSrpVerifier derives a fresh random salt and the matching SRP verifier
// for the supplied password. It must be invoked at every site that stores a
// password hash so that the user can later authenticate via USER_SRP_AUTH.
//
// userPoolID is the full Cognito pool ID (e.g. "us-east-1_abc123"); the part
// after the underscore (poolName) is required by Cognito's SRP variant inner
// hash. The returned saltHex and verifierHex are lowercase hex strings suitable
// for direct assignment to User.SrpSalt and User.SrpVerifier.
func computeSrpVerifier(userPoolID, username, password string) (saltHex, verifierHex string, err error) {
	idx := strings.Index(userPoolID, "_")
	if idx < 0 || idx == len(userPoolID)-1 {
		return "", "", fmt.Errorf("invalid user pool ID %q: missing region prefix", userPoolID)
	}
	poolName := userPoolID[idx+1:]
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", "", err
	}
	saltHex = hex.EncodeToString(salt)
	v := ComputeVerifier(saltHex, poolName, username, password)
	return saltHex, hex.EncodeToString(v.Bytes()), nil
}

// poolNameFromID extracts the portion of a Cognito user pool ID after the
// underscore (e.g. "us-east-1_abc123" => "abc123"). The pool name is used as
// part of the Cognito SRP inner hash and the claim message. The boolean is
// false when the ID does not contain a valid region/name separator.
func poolNameFromID(userPoolID string) (string, bool) {
	idx := strings.Index(userPoolID, "_")
	if idx < 0 || idx == len(userPoolID)-1 {
		return "", false
	}
	return userPoolID[idx+1:], true
}
