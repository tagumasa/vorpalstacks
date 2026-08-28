package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// RespondToAuthChallenge responds to an authentication challenge.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_RespondToAuthChallenge.html
func (s *CognitoService) RespondToAuthChallenge(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := RespondToAuthChallengeInput{
		ClientID:       getClientId(req),
		ChallengeName:  req.GetParam("ChallengeName"),
		Session:        req.GetParam("Session"),
		Params:         req.Parameters,
		ClientMetadata: parseClientMetadata(req),
	}
	return s.respondToAuthChallengeCore(ctx, reqCtx, in)
}

// AdminRespondToAuthChallenge responds to an admin authentication challenge.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminRespondToAuthChallenge.html
func (s *CognitoService) AdminRespondToAuthChallenge(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := AdminRespondToAuthChallengeInput{
		UserPoolID:     getUserPoolID(req),
		ClientID:       getClientId(req),
		ChallengeName:  req.GetParam("ChallengeName"),
		Session:        req.GetParam("Session"),
		Params:         req.Parameters,
		ClientMetadata: parseClientMetadata(req),
	}
	return s.adminRespondToAuthChallengeCore(ctx, reqCtx, in)
}
