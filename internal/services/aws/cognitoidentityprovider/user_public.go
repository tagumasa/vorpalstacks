package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// SignUp registers a new user in the specified user pool.
func (s *CognitoService) SignUp(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.signUpCore(ctx, reqCtx, SignUpInput{
		ClientID:       getClientId(req),
		Username:       getUsername(req),
		Password:       getPassword(req),
		UserAttributes: parseUserAttributes(req),
		ValidationData: parseValidationData(req),
		ClientMetadata: parseClientMetadata(req),
	})
}

// ConfirmSignUp confirms a user's registration with the confirmation code.
func (s *CognitoService) ConfirmSignUp(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.confirmSignUpCore(ctx, reqCtx, ConfirmSignUpInput{
		ClientID:         getClientId(req),
		Username:         getUsername(req),
		ConfirmationCode: getConfirmationCode(req),
	})
}

// AdminConfirmSignUp confirms a user's registration as an administrator.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminConfirmSignUp.html
func (s *CognitoService) AdminConfirmSignUp(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminConfirmSignUpCore(ctx, reqCtx, AdminConfirmSignUpInput{
		UserPoolID: getUserPoolID(req),
		Username:   getUsername(req),
	})
}
