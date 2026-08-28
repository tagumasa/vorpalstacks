package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// InitiateAuth initiates the authentication flow.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_InitiateAuth.html
func (s *CognitoService) InitiateAuth(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := InitiateAuthInput{
		AuthFlow:       req.GetParam("AuthFlow"),
		ClientID:       getClientId(req),
		Params:         req.Parameters,
		ValidationData: parseValidationData(req),
		ClientMetadata: parseClientMetadata(req),
	}
	return s.initiateAuthCore(ctx, reqCtx, in)
}

// AdminInitiateAuth initiates the admin authentication flow.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminInitiateAuth.html
func (s *CognitoService) AdminInitiateAuth(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := AdminInitiateAuthInput{
		UserPoolID:     getUserPoolID(req),
		AuthFlow:       req.GetParam("AuthFlow"),
		ClientID:       getClientId(req),
		Params:         req.Parameters,
		ValidationData: parseValidationData(req),
		ClientMetadata: parseClientMetadata(req),
	}
	return s.adminInitiateAuthCore(ctx, reqCtx, in)
}
