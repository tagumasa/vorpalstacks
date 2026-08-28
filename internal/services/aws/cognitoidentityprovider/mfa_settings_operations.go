package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// AdminSetUserMFAPreference sets the MFA preferences for a user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminSetUserMFAPreference.html
func (s *CognitoService) AdminSetUserMFAPreference(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminSetUserMFAPreferenceCore(reqCtx, AdminSetUserMFAPreferenceInput{
		UserPoolID: req.GetParam("UserPoolId"),
		Username:   getUsername(req),
		Params:     req.Parameters,
	})
}

// SetUserMFAPreference sets the MFA preferences for the authenticated user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_SetUserMFAPreference.html
func (s *CognitoService) SetUserMFAPreference(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.setUserMFAPreferenceCore(reqCtx, SetUserMFAPreferenceInput{
		AccessToken: getAccessToken(req),
		Params:      req.Parameters,
	})
}

// AdminSetUserSettings sets the legacy MFA settings for a user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminSetUserSettings.html
func (s *CognitoService) AdminSetUserSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminSetUserSettingsCore(reqCtx, AdminSetUserSettingsInput{
		UserPoolID: req.GetParam("UserPoolId"),
		Username:   getUsername(req),
		Params:     req.Parameters,
	})
}

// SetUserSettings sets the legacy MFA settings for the authenticated user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_SetUserSettings.html
func (s *CognitoService) SetUserSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.setUserSettingsCore(reqCtx, SetUserSettingsInput{
		AccessToken: getAccessToken(req),
		Params:      req.Parameters,
	})
}
