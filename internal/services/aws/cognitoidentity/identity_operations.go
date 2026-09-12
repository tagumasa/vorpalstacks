package cognitoidentity

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
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

// DeleteIdentities deletes the identities resolved by their IDs; the IDs may
// belong to different pools.
func (s *CognitoIdentityService) DeleteIdentities(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	unprocessed, err := s.deleteIdentitiesCore(reqCtx, DeleteIdentitiesInput{
		IdentityIDs: getStringSliceParam(req, "IdentityIdsToDelete"),
	})
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	if len(unprocessed) > 0 {
		items := make([]map[string]interface{}, 0, len(unprocessed))
		for _, entry := range unprocessed {
			items = append(items, map[string]interface{}{
				"IdentityId": entry.IdentityID,
				"ErrorCode":  entry.ErrorCode,
			})
		}
		result["UnprocessedIdentityIds"] = items
	}

	return result, nil
}

// ListIdentities lists the identities in an identity pool.
func (s *CognitoIdentityService) ListIdentities(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_, maxResultsProvided := req.Parameters["MaxResults"]
	result, err := s.listIdentitiesCore(reqCtx, ListIdentitiesInput{
		IdentityPoolID:     req.GetParam("IdentityPoolId"),
		MaxResultsProvided: maxResultsProvided,
		MaxResults:         request.GetIntParam(req.Parameters, "MaxResults"),
		NextToken:          request.GetStringParam(req.Parameters, "NextToken"),
		HideDisabled:       getBoolParam(req, "HideDisabled"),
	})
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(result.Identities))
	for i := range result.Identities {
		items = append(items, identityResultToHTTP(&result.Identities[i]))
	}

	resp := map[string]interface{}{
		"IdentityPoolId": result.IdentityPoolID,
		"Identities":     items,
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}

	return resp, nil
}

// GetOpenIdToken gets an OpenID token for a Cognito identity.
func (s *CognitoIdentityService) GetOpenIdToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.getOpenIdTokenCore(reqCtx, GetOpenIdTokenInput{
		IdentityID: req.GetParam("IdentityId"),
		Logins:     parseMapParam(req, "Logins"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"IdentityId": result.IdentityID,
		"Token":      result.Token,
	}, nil
}

// UnlinkIdentity unlinks login providers from a Cognito identity.
func (s *CognitoIdentityService) UnlinkIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.unlinkIdentityCore(reqCtx, UnlinkIdentityInput{
		IdentityID:     req.GetParam("IdentityId"),
		LoginsToRemove: getStringSliceParam(req, "LoginsToRemove"),
		Logins:         parseMapParam(req, "Logins"),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
