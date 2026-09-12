package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
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
	return s.getUserAttributeVerificationCodeCore(ctx, reqCtx, GetUserAttributeVerificationCodeInput{
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
	return s.resendConfirmationCodeCore(ctx, reqCtx, ResendConfirmationCodeInput{
		ClientID:       getClientId(req),
		Username:       getUsername(req),
		ClientMetadata: parseClientMetadata(req),
	})
}

// GetUserAuthFactors returns the configured authentication factors for a user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetUserAuthFactors.html
func (s *CognitoService) GetUserAuthFactors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.getUserAuthFactorsCore(reqCtx, GetUserAuthFactorsInput{AccessToken: getAccessToken(req)})
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
