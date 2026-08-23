package cognitoidentityprovider

import (
	"fmt"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

func applyUserPoolUpdates(pool *cognitostore.UserPool, req *request.ParsedRequest) error {
	if v := req.GetParam("PoolName"); v != "" {
		pool.Name = v
	}
	if v := req.GetParam("MfaConfiguration"); v != "" {
		if !validateUserPoolMfaConfig(v) {
			return ErrInvalidParameter
		}
		pool.MfaConfiguration = v
	}
	if v := req.GetParam("DeletionProtection"); v != "" {
		if !validateDeletionProtection(v) {
			return ErrInvalidParameter
		}
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
		for _, a := range v {
			if !validateAliasAttribute(a) {
				return fmt.Errorf("invalid AliasAttribute: %s", a)
			}
		}
		pool.AliasAttributes = v
	}
	if v := request.GetStringList(req.Parameters, "UsernameAttributes"); v != nil {
		for _, a := range v {
			if !validateUsernameAttribute(a) {
				return fmt.Errorf("invalid UsernameAttribute: %s", a)
			}
		}
		pool.UsernameAttributes = v
	}
	if v := request.GetStringList(req.Parameters, "AutoVerifiedAttributes"); v != nil {
		for _, a := range v {
			if !validateVerifiedAttribute(a) {
				return fmt.Errorf("invalid AutoVerifiedAttribute: %s", a)
			}
		}
		pool.AutoVerifiedAttributes = v
	}
	// The Schema parameter is a CreateUserPool-only member; the caller
	// applies it on the create path. UpdateUserPool has no such member in
	// the model, so an update request never mutates the pool schema.
	if v, err := parsePasswordPolicyWithBase(req, pool.PasswordPolicy); err != nil {
		return err
	} else if v != nil {
		pool.PasswordPolicy = v
	}
	if v := parseLambdaConfigWithBase(req, pool.LambdaConfig); v != nil {
		pool.LambdaConfig = v
	}
	if v := parseEmailConfiguration(req); v != nil {
		if v.EmailSendingAccount != "" && !validateEmailSendingAccount(v.EmailSendingAccount) {
			return ErrInvalidParameter
		}
		pool.EmailConfiguration = v
	}
	if v := parseSmsConfiguration(req); v != nil {
		pool.SmsConfiguration = v
	}
	if v := parseAdminCreateUserConfig(req); v != nil {
		pool.AdminCreateUserConfig = v
	}
	if v := parseVerificationMessageTemplate(req); v != nil {
		if v.DefaultEmailOption != "" && !validateDefaultEmailOption(v.DefaultEmailOption) {
			return ErrInvalidParameter
		}
		pool.VerificationMessageTemplate = v
	}
	if v := parseUserAttributeUpdateSettings(req); v != nil {
		pool.UserAttributeUpdateSettings = v
	}
	if v := parseUserPoolAddOns(req); v != nil {
		if v.AdvancedSecurityMode != "" && !validateAdvancedSecurityMode(v.AdvancedSecurityMode) {
			return ErrInvalidParameter
		}
		pool.UserPoolAddOns = v
	}
	if v := parseAccountRecoverySetting(req); v != nil {
		for _, rm := range v.RecoveryMechanisms {
			if rm.Name != "" && !validateRecoveryOptionName(rm.Name) {
				return ErrInvalidParameter
			}
		}
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
		if !validateUserPoolTier(v) {
			return ErrInvalidParameter
		}
		pool.UserPoolTier = v
	}

	// Re-run the whole-pool validation so updates are held to the same
	// model-derived constraints as creation.
	return validateUserPoolConfig(pool)
}
