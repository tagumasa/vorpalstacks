package cognitoidentity

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// GetId obtains a unique identity ID for a Cognito identity pool.
func (s *CognitoIdentityService) GetId(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identityID, err := s.getIdCore(reqCtx, GetIdInput{
		IdentityPoolID: req.GetParam("IdentityPoolId"),
		AccountID:      req.GetParam("AccountId"),
		Logins:         parseMapParam(req, "Logins"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"IdentityId": identityID,
	}, nil
}

// GetCredentialsForIdentity returns temporary credentials for an identity.
// In the enhanced authflow, this is functionally equivalent to calling
// GetOpenIdToken followed by AssumeRoleWithWebIdentity.
func (s *CognitoIdentityService) GetCredentialsForIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.getCredentialsForIdentityCore(reqCtx, GetCredentialsForIdentityInput{
		IdentityID:    req.GetParam("IdentityId"),
		CustomRoleARN: req.GetParam("CustomRoleArn"),
		Logins:        parseMapParam(req, "Logins"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"IdentityId": result.IdentityID,
		"Credentials": map[string]interface{}{
			"AccessKeyId":  result.AccessKeyID,
			"SecretKey":    result.SecretAccessKey,
			"SessionToken": result.SessionToken,
			"Expiration":   result.ExpirationUnix,
		},
	}, nil
}

// DescribeIdentity returns information about a Cognito identity.
func (s *CognitoIdentityService) DescribeIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.describeIdentityCore(reqCtx, DescribeIdentityInput{
		IdentityID: req.GetParam("IdentityId"),
	})
	if err != nil {
		return nil, err
	}

	return identityResultToHTTP(result), nil
}

// identityResultToHTTP serialises an IdentityResult into the member format
// shared by the DescribeIdentity and ListIdentities responses.
func identityResultToHTTP(r *IdentityResult) map[string]interface{} {
	return map[string]interface{}{
		"IdentityId":       r.ID,
		"CreationDate":     r.CreationDate.Unix(),
		"LastModifiedDate": r.LastModifiedDate.Unix(),
		"Logins":           formatLoginKeys(r.Logins),
	}
}
