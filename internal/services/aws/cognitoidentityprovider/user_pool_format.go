package cognitoidentityprovider

import (
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

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
		result["DeletionProtection"] = deletionProtectionInactive
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
	// The wire member is SchemaAttributes (the model member name) and the
	// value is the projected schema: the standard attribute set plus the
	// pool's custom attributes with their custom:/dev: prefix.
	schemaAttrs := schemaAttributesForDescribe(pool)
	schema := make([]map[string]interface{}, 0, len(schemaAttrs))
	for _, sa := range schemaAttrs {
		entry := map[string]interface{}{
			"Name":                   sa.Name,
			"AttributeDataType":      sa.AttributeDataType,
			"DeveloperOnlyAttribute": sa.DeveloperOnlyAttribute,
			"Mutable":                sa.Mutable,
			"Required":               sa.Required,
		}
		if sa.NumberAttributeConstraints != nil {
			num := map[string]interface{}{}
			if sa.NumberAttributeConstraints.MinValue != "" {
				num["MinValue"] = sa.NumberAttributeConstraints.MinValue
			}
			if sa.NumberAttributeConstraints.MaxValue != "" {
				num["MaxValue"] = sa.NumberAttributeConstraints.MaxValue
			}
			entry["NumberAttributeConstraints"] = num
		}
		if sa.StringAttributeConstraints != nil {
			str := map[string]interface{}{}
			if sa.StringAttributeConstraints.MinLength != "" {
				str["MinLength"] = sa.StringAttributeConstraints.MinLength
			}
			if sa.StringAttributeConstraints.MaxLength != "" {
				str["MaxLength"] = sa.StringAttributeConstraints.MaxLength
			}
			entry["StringAttributeConstraints"] = str
		}
		schema = append(schema, entry)
	}
	result["SchemaAttributes"] = schema
	if pool.PasswordPolicy != nil || (pool.SignInPolicy != nil && len(pool.SignInPolicy.AllowedFirstAuthFactors) > 0) {
		policies := map[string]interface{}{}
		if pool.PasswordPolicy != nil {
			policies["PasswordPolicy"] = map[string]interface{}{
				"MinimumLength":                 pool.PasswordPolicy.MinimumLength,
				"RequireUppercase":              pool.PasswordPolicy.RequireUppercase,
				"RequireLowercase":              pool.PasswordPolicy.RequireLowercase,
				"RequireNumbers":                pool.PasswordPolicy.RequireNumbers,
				"RequireSymbols":                pool.PasswordPolicy.RequireSymbols,
				"TemporaryPasswordValidityDays": pool.PasswordPolicy.TemporaryPasswordValidityDays,
				"PasswordHistorySize":           pool.PasswordPolicy.PasswordHistorySize,
			}
		}
		if pool.SignInPolicy != nil && len(pool.SignInPolicy.AllowedFirstAuthFactors) > 0 {
			policies["SignInPolicy"] = map[string]interface{}{
				"AllowedFirstAuthFactors": pool.SignInPolicy.AllowedFirstAuthFactors,
			}
		}
		result["Policies"] = policies
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
		addOns := map[string]interface{}{
			"AdvancedSecurityMode": pool.UserPoolAddOns.AdvancedSecurityMode,
		}
		if fl := pool.UserPoolAddOns.AdvancedSecurityAdditionalFlows; fl != nil {
			addOns["AdvancedSecurityAdditionalFlows"] = map[string]interface{}{
				"CustomAuthMode": fl.CustomAuthMode,
			}
		}
		result["UserPoolAddOns"] = addOns
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
	if config.EumsSms != nil {
		eums := map[string]interface{}{}
		if config.EumsSms.CallerArn != "" {
			eums["CallerArn"] = config.EumsSms.CallerArn
		}
		if config.EumsSms.ExternalId != "" {
			eums["ExternalId"] = config.EumsSms.ExternalId
		}
		if config.EumsSms.OriginationIdentity != "" {
			eums["OriginationIdentity"] = config.EumsSms.OriginationIdentity
		}
		if config.EumsSms.ConfigurationSetName != "" {
			eums["ConfigurationSetName"] = config.EumsSms.ConfigurationSetName
		}
		if config.EumsSms.InEntityId != "" {
			eums["InEntityId"] = config.EumsSms.InEntityId
		}
		if config.EumsSms.InTemplateId != "" {
			eums["InTemplateId"] = config.EumsSms.InTemplateId
		}
		if config.EumsSms.Region != "" {
			eums["Region"] = config.EumsSms.Region
		}
		result["EumsSms"] = eums
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
		result["Attributes"] = userAttributeList(user)
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

// formatMFAOptions projects the legacy MFA options list (MFAOptionType:
// DeliveryMedium/AttributeName) shared by the user responses.
func formatMFAOptions(user *cognitostore.User) []map[string]interface{} {
	mfaOpts := make([]map[string]interface{}, 0, len(user.MFAOptions))
	for _, opt := range user.MFAOptions {
		mfaOpts = append(mfaOpts, map[string]interface{}{
			"DeliveryMedium": opt.DeliveryMedium,
			"AttributeName":  opt.AttributeName,
		})
	}
	return mfaOpts
}

// userMFAPreferences derives the activated-MFA names the model documents for
// UserMFASettingList — SMS_MFA, EMAIL_OTP and SOFTWARE_TOKEN_MFA — together
// with the user's preferred setting. A software token counts as activated
// only once enrolled and verified, matching the sign-in challenge
// selection. The legacy MFAOptions SMS delivery counts as SMS_MFA.
func userMFAPreferences(user *cognitostore.User) (preferred string, settings []string) {
	has := func(name string) bool {
		for _, s := range settings {
			if s == name {
				return true
			}
		}
		return false
	}
	add := func(name string) {
		if !has(name) {
			settings = append(settings, name)
		}
	}

	software := user.SoftwareTokenMfa != nil && user.SoftwareTokenMfa.Enabled && user.SoftwareTokenMfa.Verified
	sms := user.SmsMfa != nil && user.SmsMfa.Enabled
	email := user.EmailMfa != nil && user.EmailMfa.Enabled

	if software {
		add("SOFTWARE_TOKEN_MFA")
		if preferred == "" && user.SoftwareTokenMfa.PreferredMfa {
			preferred = "SOFTWARE_TOKEN_MFA"
		}
	}
	if sms {
		add("SMS_MFA")
		if preferred == "" && user.SmsMfa.PreferredMfa {
			preferred = "SMS_MFA"
		}
	}
	if email {
		add("EMAIL_OTP")
		if preferred == "" && user.EmailMfa.PreferredMfa {
			preferred = "EMAIL_OTP"
		}
	}
	for _, opt := range user.MFAOptions {
		if opt.DeliveryMedium == "SMS" {
			add("SMS_MFA")
		}
	}
	return preferred, settings
}

// formatGetUserResponse projects the GetUser response: the model defines
// Username, UserAttributes, MFAOptions, PreferredMfaSetting and
// UserMFASettingList — the UserType status members do not belong here.
func formatGetUserResponse(user *cognitostore.User) map[string]interface{} {
	result := map[string]interface{}{
		"Username":       user.Username,
		"UserAttributes": userAttributeList(user),
	}
	if len(user.MFAOptions) > 0 {
		result["MFAOptions"] = formatMFAOptions(user)
	}
	preferred, settings := userMFAPreferences(user)
	if preferred != "" {
		result["PreferredMfaSetting"] = preferred
	}
	if len(settings) > 0 {
		result["UserMFASettingList"] = settings
	}
	return result
}

// formatAdminGetUserResponse projects the AdminGetUser response: the model
// defines the UserType status members plus UserAttributes,
// PreferredMfaSetting and UserMFASettingList.
func formatAdminGetUserResponse(user *cognitostore.User) map[string]interface{} {
	result := map[string]interface{}{
		"Username":             user.Username,
		"Enabled":              user.Enabled,
		"UserStatus":           user.UserStatus,
		"UserCreateDate":       user.CreatedDate.Unix(),
		"UserLastModifiedDate": user.LastModifiedDate.Unix(),
		"UserAttributes":       userAttributeList(user),
	}
	if len(user.MFAOptions) > 0 {
		result["MFAOptions"] = formatMFAOptions(user)
	}
	preferred, settings := userMFAPreferences(user)
	if preferred != "" {
		result["PreferredMfaSetting"] = preferred
	}
	if len(settings) > 0 {
		result["UserMFASettingList"] = settings
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
