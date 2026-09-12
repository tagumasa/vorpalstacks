package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// Handlers for the WebAuthn credential family; the ceremonies themselves
// live in webauthn_core.go.

// StartWebAuthnRegistration starts a WebAuthn credential registration flow.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_StartWebAuthnRegistration.html
func (s *CognitoService) StartWebAuthnRegistration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.startWebAuthnRegistrationCore(reqCtx, StartWebAuthnRegistrationInput{
		AccessToken: getAccessToken(req),
	})
}

// CompleteWebAuthnRegistration completes a WebAuthn credential registration.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CompleteWebAuthnRegistration.html
func (s *CognitoService) CompleteWebAuthnRegistration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.completeWebAuthnRegistrationCore(reqCtx, CompleteWebAuthnRegistrationInput{
		AccessToken: getAccessToken(req),
		Params:      req.Parameters,
	})
}

// ListWebAuthnCredentials lists registered WebAuthn credentials.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListWebAuthnCredentials.html
func (s *CognitoService) ListWebAuthnCredentials(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.listWebAuthnCredentialsCore(reqCtx, ListWebAuthnCredentialsInput{
		AccessToken: getAccessToken(req),
		NextToken:   req.GetParam("NextToken"),
		Params:      req.Parameters,
	})
}

// DeleteWebAuthnCredential deletes a WebAuthn credential.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteWebAuthnCredential.html
func (s *CognitoService) DeleteWebAuthnCredential(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.deleteWebAuthnCredentialCore(reqCtx, DeleteWebAuthnCredentialInput{
		AccessToken:  getAccessToken(req),
		CredentialID: req.GetParam("CredentialId"),
	})
}
