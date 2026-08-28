package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// RevokeToken revokes a refresh token.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_RevokeToken.html
func (s *CognitoService) RevokeToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.revokeTokenCore(reqCtx, RevokeTokenInput{
		Token:        req.GetParam("Token"),
		ClientID:     req.GetParam("ClientId"),
		ClientSecret: req.GetParam("ClientSecret"),
	})
}

// GetTokensFromRefreshToken exchanges a refresh token for new access and ID tokens.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetTokensFromRefreshToken.html
func (s *CognitoService) GetTokensFromRefreshToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.getTokensFromRefreshTokenCore(reqCtx, GetTokensFromRefreshTokenInput{
		RefreshToken:   req.GetParam("RefreshToken"),
		ClientID:       req.GetParam("ClientId"),
		DeviceKey:      req.GetParam("DeviceKey"),
		ClientMetadata: parseClientMetadata(req),
	})
}

// GetUserAttributeVerificationCode generates a verification code for a user attribute.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetUserAttributeVerificationCode.html
func (s *CognitoService) GetUserAttributeVerificationCode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.getUserAttributeVerificationCodeCore(reqCtx, GetUserAttributeVerificationCodeInput{
		AccessToken:   getAccessToken(req),
		AttributeName: req.GetParam("AttributeName"),
	})
}

// VerifyUserAttribute verifies a user attribute with a confirmation code.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_VerifyUserAttribute.html
func (s *CognitoService) VerifyUserAttribute(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.verifyUserAttributeCore(reqCtx, VerifyUserAttributeInput{
		AccessToken:   getAccessToken(req),
		AttributeName: req.GetParam("AttributeName"),
		Code:          req.GetParam("Code"),
	})
}

// ResendConfirmationCode resends the confirmation code for user registration.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ResendConfirmationCode.html
func (s *CognitoService) ResendConfirmationCode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.resendConfirmationCodeCore(reqCtx, ResendConfirmationCodeInput{
		ClientID: getClientId(req),
		Username: getUsername(req),
	})
}

// GetUserAuthFactors returns the configured authentication factors for a user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetUserAuthFactors.html
func (s *CognitoService) GetUserAuthFactors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.getUserAuthFactorsCore(reqCtx, GetUserAuthFactorsInput{AccessToken: getAccessToken(req)})
}

// determineDeliveryMedium picks the appropriate delivery medium and attribute
// name for confirmation codes, based on the pool's AutoVerifiedAttributes and
// the user's contact information.
func determineDeliveryMedium(pool *cognitostore.UserPool, user *cognitostore.User) (medium, attrName string) {
	// Check AutoVerifiedAttributes first — these take priority
	for _, attr := range pool.AutoVerifiedAttributes {
		if attr == "phone_number" {
			if phone, ok := user.Attributes["phone_number"]; ok && phone != "" {
				return "SMS", "phone_number"
			}
		}
		if attr == "email" {
			if email, ok := user.Attributes["email"]; ok && email != "" {
				return "EMAIL", "email"
			}
		}
	}
	// Fallback: use any available contact attribute
	if email, ok := user.Attributes["email"]; ok && email != "" {
		return "EMAIL", "email"
	}
	if phone, ok := user.Attributes["phone_number"]; ok && phone != "" {
		return "SMS", "phone_number"
	}
	// Default to EMAIL when the user has neither phone nor email — this
	// matches Cognito behaviour which always returns a delivery medium
	// rather than an empty CodeDeliveryDetails.
	return "EMAIL", "email"
}

// AdminGetUserAuthFactors returns the configured authentication factors for a
// user as viewed by an administrator.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminGetUserAuthFactors.html
func (s *CognitoService) AdminGetUserAuthFactors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminGetUserAuthFactorsCore(reqCtx, AdminGetUserAuthFactorsInput{
		UserPoolID: getUserPoolID(req),
		Username:   getUsername(req),
	})
}

// computeUserAuthFactors builds the auth factors response for a user, shared
// by GetUserAuthFactors (access-token based) and AdminGetUserAuthFactors
// (admin based).
func computeUserAuthFactors(user *cognitostore.User) map[string]interface{} {
	result := map[string]interface{}{
		"Username": user.Username,
	}

	var configuredFactors []string
	var preferredFactor string

	if user.PasswordHash != "" {
		configuredFactors = append(configuredFactors, "PASSWORD")
	}

	if user.SmsMfa != nil && user.SmsMfa.Enabled {
		configuredFactors = append(configuredFactors, "SMS_OTP")
		if user.SmsMfa.PreferredMfa && preferredFactor == "" {
			preferredFactor = "SMS_OTP"
		}
	}

	if user.EmailMfa != nil && user.EmailMfa.Enabled {
		configuredFactors = append(configuredFactors, "EMAIL_OTP")
		if user.EmailMfa.PreferredMfa && preferredFactor == "" {
			preferredFactor = "EMAIL_OTP"
		}
	}

	if user.SoftwareTokenMfa != nil && user.SoftwareTokenMfa.Enabled {
		configuredFactors = append(configuredFactors, "SOFTWARE_TOKEN")
		if user.SoftwareTokenMfa.PreferredMfa && preferredFactor == "" {
			preferredFactor = "SOFTWARE_TOKEN"
		}
	}

	if user.WebAuthnMfaEnabled {
		configuredFactors = append(configuredFactors, "WEB_AUTHN")
	}

	if len(user.MFAOptions) > 0 {
		for _, opt := range user.MFAOptions {
			if opt.DeliveryMedium == "SMS" {
				alreadyHas := false
				for _, f := range configuredFactors {
					if f == "SMS_OTP" {
						alreadyHas = true
						break
					}
				}
				if !alreadyHas {
					configuredFactors = append(configuredFactors, "SMS_OTP")
				}
			}
		}
	}

	if preferredFactor != "" {
		result["PreferredMfaSetting"] = preferredFactor
	}

	if len(configuredFactors) == 0 {
		configuredFactors = []string{}
	}
	result["ConfiguredUserAuthFactors"] = configuredFactors

	if len(user.MFAOptions) > 0 {
		mfaSettings := make([]map[string]interface{}, 0, len(user.MFAOptions))
		for _, opt := range user.MFAOptions {
			entry := map[string]interface{}{
				"DeliveryMedium": opt.DeliveryMedium,
			}
			if opt.AttributeName != "" {
				entry["AttributeName"] = opt.AttributeName
			}
			mfaSettings = append(mfaSettings, entry)
		}
		result["UserMFASettingList"] = mfaSettings
	}

	return result
}
