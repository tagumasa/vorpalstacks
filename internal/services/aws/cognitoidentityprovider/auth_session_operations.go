package cognitoidentityprovider

import (
	"context"
	"crypto/rand"
	"fmt"

	"vorpalstacks/internal/common/request"
)

func generateConfirmationCode() (string, error) {
	const maxCode = 1000000
	const limit = (1 << 24) / maxCode * maxCode
	for {
		b := make([]byte, 3)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		n := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
		if n < limit {
			return fmt.Sprintf("%06d", n%maxCode), nil
		}
	}
}

// SignOut signs out a user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_SignOut.html
func (s *CognitoService) SignOut(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.signOutCore(reqCtx, SignOutInput{AccessToken: getAccessToken(req)})
}

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
