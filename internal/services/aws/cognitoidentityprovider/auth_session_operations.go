package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// GlobalSignOut signs out a user from all devices.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GlobalSignOut.html
func (s *CognitoService) GlobalSignOut(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.globalSignOutCore(reqCtx, GlobalSignOutInput{AccessToken: getAccessToken(req)})
}

// ChangePassword changes the password for a user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ChangePassword.html
func (s *CognitoService) ChangePassword(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.changePasswordCore(reqCtx, ChangePasswordInput{
		AccessToken:      getAccessToken(req),
		PreviousPassword: getPreviousPassword(req),
		NewPassword:      getNewPassword(req),
	})
}

// ForgotPassword initiates the forgot password flow.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ForgotPassword.html
func (s *CognitoService) ForgotPassword(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.forgotPasswordCore(ctx, reqCtx, ForgotPasswordInput{
		ClientID: getClientId(req),
		Username: getUsername(req),
	})
}

// ConfirmForgotPassword confirms the forgot password flow.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ConfirmForgotPassword.html
func (s *CognitoService) ConfirmForgotPassword(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.confirmForgotPasswordCore(reqCtx, ConfirmForgotPasswordInput{
		ClientID:         getClientId(req),
		Username:         getUsername(req),
		Password:         getPassword(req),
		ConfirmationCode: getConfirmationCode(req),
	})
}
