package cognitoidentityprovider

import (
	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// applyCreateUserPoolRequest applies a CreateUserPool request onto a
// fresh store record: the create-only members first, then the members the
// create and update requests share. It is pure member application; the
// whole-pool validation runs once in createUserPoolCore, the single
// persistence path.
func applyCreateUserPoolRequest(pool *cognitostore.UserPool, req *request.ParsedRequest, poolRegion string) error {
	applyCreateOnlyUserPoolMembers(pool, req)
	return applyUserPoolConfigMembers(pool, req, poolRegion)
}

// applyCreateOnlyUserPoolMembers applies the members CreateUserPool carries
// but UpdateUserPool does not: the request model of the update operation has
// no Schema, AliasAttributes, UsernameAttributes or UsernameConfiguration
// member, so those parts of a pool are create-only and survive every update.
func applyCreateOnlyUserPoolMembers(pool *cognitostore.UserPool, req *request.ParsedRequest) {
	// The Schema parameter is a CreateUserPool-only member; the caller
	// applies it before this helper so the whole-pool validation still sees
	// the schema definitions.
	if v := request.GetStringList(req.Parameters, "AliasAttributes"); v != nil {
		pool.AliasAttributes = v
	}
	if v := request.GetStringList(req.Parameters, "UsernameAttributes"); v != nil {
		pool.UsernameAttributes = v
	}
	if v := parseUsernameConfiguration(req); v != nil {
		pool.UsernameConfiguration = v
	}
}

// applyUserPoolConfigMembers applies the configuration members present in the
// request onto pool. It is member-presence application, deliberately free of
// reset semantics: the create path feeds it a defaults record, and
// rebuildUserPoolFromUpdate feeds it the defaults record it rebuilt, so in
// both cases an absent member leaves the member at its default value. The
// per-member value validation runs here; the whole-pool validation belongs to
// the callers. poolRegion is the pool's home Region, named so the
// region-bound members (the End User Messaging SMS configuration) can be
// checked against it.
func applyUserPoolConfigMembers(pool *cognitostore.UserPool, req *request.ParsedRequest, poolRegion string) error {
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
	if v := request.GetStringList(req.Parameters, "AutoVerifiedAttributes"); v != nil {
		pool.AutoVerifiedAttributes = v
	}
	if v, err := parsePasswordPolicyWithBase(req, pool.PasswordPolicy); err != nil {
		return err
	} else if v != nil {
		pool.PasswordPolicy = v
	}
	if v, err := parseSignInPolicy(req); err != nil {
		return err
	} else if v != nil {
		pool.SignInPolicy = v
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
	if v, err := parseSmsConfiguration(req, poolRegion); err != nil {
		return err
	} else if v != nil {
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
		if fl := v.AdvancedSecurityAdditionalFlows; fl != nil && fl.CustomAuthMode != "" && !validateCustomAuthMode(fl.CustomAuthMode) {
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
	return nil
}

// rebuildUserPoolFromUpdate implements the model's UpdateUserPool contract:
// the request replaces the pool configuration — a member omitted from the
// request reverts to its default value, exactly as the same request members
// would shape a freshly created pool. Only what no update request can carry
// survives from the stored record: identity (ID, ARN, status, name,
// timestamps), the signing keys, the create-only members (schema, alias and
// username attributes, username configuration) and the MFA settings owned by
// SetUserPoolMfaConfig.
func rebuildUserPoolFromUpdate(stored *cognitostore.UserPool, req *request.ParsedRequest, poolRegion string) (*cognitostore.UserPool, error) {
	rebuilt := &cognitostore.UserPool{
		ID:                            stored.ID,
		Name:                          stored.Name,
		Arn:                           stored.Arn,
		Status:                        stored.Status,
		CreationDate:                  stored.CreationDate,
		LastModifiedDate:              stored.LastModifiedDate,
		MfaConfiguration:              "OFF",
		SchemaAttributes:              stored.SchemaAttributes,
		AliasAttributes:               stored.AliasAttributes,
		UsernameAttributes:            stored.UsernameAttributes,
		UsernameConfiguration:         stored.UsernameConfiguration,
		JwtPrivateKey:                 stored.JwtPrivateKey,
		JwtPublicKey:                  stored.JwtPublicKey,
		JwtKeyID:                      stored.JwtKeyID,
		MfaConfigurationSms:           stored.MfaConfigurationSms,
		MfaConfigurationSoftwareToken: stored.MfaConfigurationSoftwareToken,
		EmailMfaConfig:                stored.EmailMfaConfig,
		WebAuthnConfiguration:         stored.WebAuthnConfiguration,
		EstimatedNumberOfUsers:        stored.EstimatedNumberOfUsers,
	}
	if err := applyUserPoolConfigMembers(rebuilt, req, poolRegion); err != nil {
		return nil, err
	}
	return rebuilt, validateUserPoolConfig(rebuilt)
}
