package cognitoidentityprovider

import (
	"crypto/subtle"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// RevokeTokenInput carries the wire parameters of RevokeToken.
type RevokeTokenInput struct {
	Token        string
	ClientID     string
	ClientSecret string
}

// GetTokensFromRefreshTokenInput carries the wire parameters of
// GetTokensFromRefreshToken.
type GetTokensFromRefreshTokenInput struct {
	RefreshToken   string
	ClientID       string
	DeviceKey      string
	ClientMetadata map[string]string
}

// GetUserAttributeVerificationCodeInput carries the wire parameters of
// GetUserAttributeVerificationCode.
type GetUserAttributeVerificationCodeInput struct {
	AccessToken   string
	AttributeName string
}

// VerifyUserAttributeInput carries the wire parameters of
// VerifyUserAttribute.
type VerifyUserAttributeInput struct {
	AccessToken   string
	AttributeName string
	Code          string
}

// ResendConfirmationCodeInput carries the wire parameters of
// ResendConfirmationCode.
type ResendConfirmationCodeInput struct {
	ClientID string
	Username string
}

// GetUserAuthFactorsInput carries the wire parameters of GetUserAuthFactors.
type GetUserAuthFactorsInput struct {
	AccessToken string
}

// AdminGetUserAuthFactorsInput carries the wire parameters of
// AdminGetUserAuthFactors.
type AdminGetUserAuthFactorsInput struct {
	UserPoolID string
	Username   string
}

// revokeTokenCore revokes a refresh token. RevokeToken always returns
// 200 OK per AWS spec, even for invalid, expired, or already-revoked
// tokens. The sole exception is a ClientSecret mismatch for a client that
// has a secret configured.
func (s *CognitoService) revokeTokenCore(reqCtx *request.RequestContext, in RevokeTokenInput) (interface{}, error) {
	if in.Token == "" || in.ClientID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	rt, err := store.GetRefreshTokenByValue(in.Token)
	if err != nil {
		return response.EmptyResponse(), nil
	}

	// Silently succeed when the token belongs to a different client —
	// do not delete, do not reveal the token's existence.
	if rt.ClientID != in.ClientID {
		return response.EmptyResponse(), nil
	}

	// Verify ClientSecret when the client has one configured.
	if in.ClientSecret != "" {
		client, err := store.GetUserPoolClient(rt.UserPoolID, in.ClientID)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		if client.ClientSecret != "" && subtle.ConstantTimeCompare([]byte(client.ClientSecret), []byte(in.ClientSecret)) != 1 {
			return nil, ErrNotAuthorized
		}
	}

	// Best-effort deletion; a concurrent revocation may have already removed
	// the token, which is not an error.
	_ = store.DeleteRefreshToken(rt.UserPoolID, rt.UserID, in.Token)

	return response.EmptyResponse(), nil
}

// getTokensFromRefreshTokenCore exchanges a refresh token for new access
// and ID tokens.
func (s *CognitoService) getTokensFromRefreshTokenCore(reqCtx *request.RequestContext, in GetTokensFromRefreshTokenInput) (interface{}, error) {
	if in.RefreshToken == "" || in.ClientID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	rt, err := store.GetRefreshTokenByValue(in.RefreshToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	if time.Now().After(rt.Expires) {
		return nil, ErrNotAuthorized
	}

	if rt.ClientID != in.ClientID {
		return nil, ErrNotAuthorized
	}

	user, err := store.GetUserByID(rt.UserID)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// When DeviceKey is provided, verify the device is remembered for this
	// user and include it in the token claims for device-aware sessions.
	if in.DeviceKey != "" {
		device, err := store.GetDevice(rt.UserPoolID, user.ID, in.DeviceKey)
		if err != nil || device.DeviceRememberedStatus != "remembered" {
			return nil, ErrNotAuthorized
		}
	}

	// CreateTokens fires TokenGenerationRefreshTokens internally and applies
	// the trigger result to the token claims. ClientMetadata is forwarded.
	accessToken, idToken, _, expiresIn, err := s.CreateTokens(reqCtx, rt.UserPoolID, user.ID, in.ClientID, TokenGenerationRefreshTokens, in.ClientMetadata)
	if err != nil {
		return nil, err
	}

	return authResultNoRefresh(accessToken, idToken, expiresIn), nil
}

// getUserAttributeVerificationCodeCore generates and stores a verification
// code for one of the caller's attributes.
func (s *CognitoService) getUserAttributeVerificationCodeCore(reqCtx *request.RequestContext, in GetUserAttributeVerificationCodeInput) (interface{}, error) {
	if in.AccessToken == "" || in.AttributeName == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, in.AccessToken)
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

	if user.Attributes[in.AttributeName] == "" {
		return nil, ErrInvalidParameter
	}

	code, err := generateConfirmationCode()
	if err != nil {
		return nil, ErrInternalError
	}
	if user.AttributeVerificationCodes == nil {
		user.AttributeVerificationCodes = make(map[string]*cognitostore.AttributeVerification)
	}
	user.AttributeVerificationCodes[in.AttributeName] = &cognitostore.AttributeVerification{
		Code:   code,
		Expiry: time.Now().UTC().Add(verificationCodeTTL),
	}
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	deliveryMedium := "EMAIL"
	if in.AttributeName == "phone_number" {
		deliveryMedium = "SMS"
	}

	return map[string]interface{}{
		"CodeDeliveryDetails": map[string]interface{}{
			"Destination":    "***",
			"DeliveryMedium": deliveryMedium,
			"AttributeName":  in.AttributeName,
		},
	}, nil
}

// verifyUserAttributeCore verifies an attribute with a confirmation code.
func (s *CognitoService) verifyUserAttributeCore(reqCtx *request.RequestContext, in VerifyUserAttributeInput) (interface{}, error) {
	if in.AccessToken == "" || in.AttributeName == "" || in.Code == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, in.AccessToken)
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

	av := user.AttributeVerificationCodes[in.AttributeName]
	if av == nil || av.Code == "" || subtle.ConstantTimeCompare([]byte(av.Code), []byte(in.Code)) != 1 {
		return nil, ErrCodeMismatch
	}

	if time.Now().After(av.Expiry) {
		return nil, ErrExpiredCode
	}

	delete(user.AttributeVerificationCodes, in.AttributeName)
	if user.Attributes == nil {
		user.Attributes = make(map[string]string)
	}
	user.Attributes[in.AttributeName+"_verified"] = "true"
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// resendConfirmationCodeCore reissues the sign-up confirmation code. The
// response is masked so the operation cannot be used to enumerate accounts.
func (s *CognitoService) resendConfirmationCodeCore(reqCtx *request.RequestContext, in ResendConfirmationCodeInput) (interface{}, error) {
	if in.ClientID == "" || in.Username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	userPool, err := store.GetUserPoolByClientID(in.ClientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	user, err := store.GetUser(userPool.ID, in.Username)
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

// getUserAuthFactorsCore returns the configured authentication factors for
// the access-token caller.
func (s *CognitoService) getUserAuthFactorsCore(reqCtx *request.RequestContext, in GetUserAuthFactorsInput) (interface{}, error) {
	if in.AccessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, in.AccessToken)
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

// adminGetUserAuthFactorsCore returns the configured authentication factors
// for a user as viewed by an administrator.
func (s *CognitoService) adminGetUserAuthFactorsCore(reqCtx *request.RequestContext, in AdminGetUserAuthFactorsInput) (interface{}, error) {
	if in.UserPoolID == "" || in.Username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(in.UserPoolID, in.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return computeUserAuthFactors(user), nil
}
