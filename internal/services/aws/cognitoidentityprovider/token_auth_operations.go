package cognitoidentityprovider

import (
	"context"
	"crypto/subtle"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// RevokeToken revokes a refresh token.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_RevokeToken.html
func (s *CognitoService) RevokeToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	token := req.GetParam("Token")
	clientID := req.GetParam("ClientId")
	clientSecret := req.GetParam("ClientSecret")

	if token == "" || clientID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// RevokeToken always returns 200 OK per AWS spec, even for invalid,
	// expired, or already-revoked tokens. The sole exception is a
	// ClientSecret mismatch for a client that has a secret configured.
	rt, err := store.GetRefreshTokenByValue(token)
	if err != nil {
		return response.EmptyResponse(), nil
	}

	// Silently succeed when the token belongs to a different client —
	// do not delete, do not reveal the token's existence.
	if rt.ClientID != clientID {
		return response.EmptyResponse(), nil
	}

	// Verify ClientSecret when the client has one configured.
	if clientSecret != "" {
		client, err := store.GetUserPoolClient(rt.UserPoolID, clientID)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		if client.ClientSecret != "" && subtle.ConstantTimeCompare([]byte(client.ClientSecret), []byte(clientSecret)) != 1 {
			return nil, ErrNotAuthorized
		}
	}

	// Best-effort deletion; a concurrent revocation may have already removed
	// the token, which is not an error.
	_ = store.DeleteRefreshToken(rt.UserPoolID, rt.UserID, token)

	return response.EmptyResponse(), nil
}

// GetTokensFromRefreshToken exchanges a refresh token for new access and ID tokens.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetTokensFromRefreshToken.html
func (s *CognitoService) GetTokensFromRefreshToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	refreshToken := req.GetParam("RefreshToken")
	clientID := req.GetParam("ClientId")
	deviceKey := req.GetParam("DeviceKey")
	clientMetadata := parseClientMetadata(req)

	if refreshToken == "" || clientID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	rt, err := store.GetRefreshTokenByValue(refreshToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	if time.Now().After(rt.Expires) {
		return nil, ErrNotAuthorized
	}

	if rt.ClientID != clientID {
		return nil, ErrNotAuthorized
	}

	user, err := store.GetUserByID(rt.UserID)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// When DeviceKey is provided, verify the device is remembered for this
	// user and include it in the token claims for device-aware sessions.
	if deviceKey != "" {
		device, err := store.GetDevice(rt.UserPoolID, user.ID, deviceKey)
		if err != nil || device.DeviceRememberedStatus != "remembered" {
			return nil, ErrNotAuthorized
		}
	}

	// CreateTokens fires TokenGenerationRefreshTokens internally and applies
	// the trigger result to the token claims. ClientMetadata is forwarded.
	accessToken, idToken, _, expiresIn, err := s.CreateTokens(reqCtx, rt.UserPoolID, user.ID, clientID, TokenGenerationRefreshTokens, clientMetadata)
	if err != nil {
		return nil, err
	}

	return authResultNoRefresh(accessToken, idToken, expiresIn), nil
}

// GetUserAttributeVerificationCode generates a verification code for a user attribute.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetUserAttributeVerificationCode.html
func (s *CognitoService) GetUserAttributeVerificationCode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	attributeName := req.GetParam("AttributeName")

	if accessToken == "" || attributeName == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.Attributes[attributeName] == "" {
		return nil, ErrInvalidParameter
	}

	code, err := generateConfirmationCode()
	if err != nil {
		return nil, ErrInternalError
	}
	if user.AttributeVerificationCodes == nil {
		user.AttributeVerificationCodes = make(map[string]*cognitostore.AttributeVerification)
	}
	user.AttributeVerificationCodes[attributeName] = &cognitostore.AttributeVerification{
		Code:   code,
		Expiry: time.Now().UTC().Add(verificationCodeTTL),
	}
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	deliveryMedium := "EMAIL"
	if attributeName == "phone_number" {
		deliveryMedium = "SMS"
	}

	return map[string]interface{}{
		"CodeDeliveryDetails": map[string]interface{}{
			"Destination":    "***",
			"DeliveryMedium": deliveryMedium,
			"AttributeName":  attributeName,
		},
	}, nil
}

// VerifyUserAttribute verifies a user attribute with a confirmation code.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_VerifyUserAttribute.html
func (s *CognitoService) VerifyUserAttribute(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	attributeName := req.GetParam("AttributeName")
	code := req.GetParam("Code")

	if accessToken == "" || attributeName == "" || code == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	av := user.AttributeVerificationCodes[attributeName]
	if av == nil || av.Code == "" || subtle.ConstantTimeCompare([]byte(av.Code), []byte(code)) != 1 {
		return nil, ErrCodeMismatch
	}

	if time.Now().After(av.Expiry) {
		return nil, ErrExpiredCode
	}

	delete(user.AttributeVerificationCodes, attributeName)
	if user.Attributes == nil {
		user.Attributes = make(map[string]string)
	}
	user.Attributes[attributeName+"_verified"] = "true"
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// ResendConfirmationCode resends the confirmation code for user registration.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ResendConfirmationCode.html
func (s *CognitoService) ResendConfirmationCode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	clientID := getClientId(req)
	username := getUsername(req)

	if clientID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	userPool, err := store.GetUserPoolByClientID(clientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	user, err := store.GetUser(userPool.ID, username)
	if err != nil {
		// Return a masked CodeDeliveryDetails to prevent user enumeration.
		return map[string]interface{}{
			"CodeDeliveryDetails": map[string]interface{}{
				"Destination":    "***",
				"DeliveryMedium": "EMAIL",
				"AttributeName":  "email",
			},
		}, nil
	}

	code, err := generateConfirmationCode()
	if err != nil {
		return nil, ErrInternalError
	}
	user.ConfirmationCode = code
	user.ConfirmationCodeExpiry = time.Now().UTC().Add(verificationCodeTTL)
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	deliveryMedium, deliveryAttr := determineDeliveryMedium(userPool, user)

	return map[string]interface{}{
		"CodeDeliveryDetails": map[string]interface{}{
			"Destination":    "***",
			"DeliveryMedium": deliveryMedium,
			"AttributeName":  deliveryAttr,
		},
	}, nil
}

// GetUserAuthFactors returns the configured authentication factors for a user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetUserAuthFactors.html
func (s *CognitoService) GetUserAuthFactors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	if accessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return computeUserAuthFactors(user), nil
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
	userPoolID := getUserPoolID(req)
	username := getUsername(req)

	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return computeUserAuthFactors(user), nil
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
