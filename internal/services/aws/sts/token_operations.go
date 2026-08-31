package sts

import (
	"context"

	"vorpalstacks/internal/common/request"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

// GetSessionToken returns a set of temporary credentials for an AWS account or IAM user.
func (s *STSService) GetSessionToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.getSessionTokenCore(reqCtx, WireInput{
		Parameters:  req.Parameters,
		AccessKeyID: callerAccessKeyID(req),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials": credentialsMap(*result),
	}, nil
}

// GetFederationToken returns a set of temporary security credentials for a federated user.
func (s *STSService) GetFederationToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.getFederationTokenCore(reqCtx, WireInput{
		Parameters:  req.Parameters,
		AccessKeyID: callerAccessKeyID(req),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials": credentialsMap(result.Credentials),
		"FederatedUser": map[string]interface{}{
			"FederatedUserId": reqCtx.GetAccountID() + ":" + result.Name,
			"Arn":             arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").STS().FederatedUser(result.Name),
		},
		"PackedPolicySize": result.PackedPolicySize,
	}, nil
}

// GetDelegatedAccessToken returns a set of temporary security credentials that represent an IAM identity centre user.
func (s *STSService) GetDelegatedAccessToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.getDelegatedAccessTokenCore(reqCtx, WireInput{Parameters: req.Parameters})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials":      credentialsMap(result.Credentials),
		"AssumedPrincipal": result.AssumedPrincipal,
		"PackedPolicySize": 0,
	}, nil
}

// GetWebIdentityToken returns a signed JSON Web Token (JWT) representing the
// calling AWS identity. The returned JWT can be used to authenticate with
// external services that support OIDC discovery. The token is signed using
// the caller-specified algorithm (RS256 or ES384).
//
// AWS spec: https://docs.aws.amazon.com/STS/latest/APIReference/API_GetWebIdentityToken.html
func (s *STSService) GetWebIdentityToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.getWebIdentityTokenCore(reqCtx, WireInput{
		Parameters:    req.Parameters,
		AccessKeyID:   callerAccessKeyID(req),
		SecurityToken: req.Headers.Get("X-Amz-Security-Token"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"WebIdentityToken": result.Token,
		"Expiration":       result.Expiration.Format(timeutils.ISO8601SimpleFormat),
	}, nil
}
