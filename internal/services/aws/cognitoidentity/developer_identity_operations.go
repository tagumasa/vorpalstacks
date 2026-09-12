package cognitoidentity

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// developer_identity_operations.go — HTTP-plane handlers for developer
// identities: token issuance, lookup, merge and unlink.

// GetOpenIdTokenForDeveloperIdentity registers (or retrieves) a developer identity and returns an OpenID token.
func (s *CognitoIdentityService) GetOpenIdTokenForDeveloperIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tokenDurationProvided := false
	tokenDuration := 0
	if _, ok := req.Parameters["TokenDuration"]; ok {
		tokenDurationProvided = true
		tokenDuration = request.GetIntParam(req.Parameters, "TokenDuration")
	}
	principalTagsRaw, principalTagsProvided := req.Parameters["PrincipalTags"]

	result, err := s.getOpenIdTokenForDeveloperIdentityCore(reqCtx, GetOpenIdTokenForDeveloperIdentityInput{
		IdentityPoolID:        req.GetParam("IdentityPoolId"),
		IdentityID:            req.GetParam("IdentityId"),
		Logins:                parseMapParam(req, "Logins"),
		TokenDurationProvided: tokenDurationProvided,
		TokenDuration:         tokenDuration,
		PrincipalTagsProvided: principalTagsProvided,
		PrincipalTagsRaw:      principalTagsRaw,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"IdentityId": result.IdentityID,
		"Token":      result.Token,
	}, nil
}

// LookupDeveloperIdentity looks up a developer identity identifier and returns the mapped identity IDs.
func (s *CognitoIdentityService) LookupDeveloperIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_, maxResultsProvided := req.Parameters["MaxResults"]
	result, err := s.lookupDeveloperIdentityCore(reqCtx, LookupDeveloperIdentityInput{
		IdentityPoolID:          req.GetParam("IdentityPoolId"),
		IdentityID:              req.GetParam("IdentityId"),
		DeveloperUserIdentifier: req.GetParam("DeveloperUserIdentifier"),
		MaxResultsProvided:      maxResultsProvided,
		MaxResults:              request.GetIntParam(req.Parameters, "MaxResults"),
		NextToken:               request.GetStringParam(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	return lookupDeveloperIdentityResult(result.MatchedIdentityID, result.DeveloperUserIDs, result.NextToken), nil
}

// lookupDeveloperIdentityResult builds the LookupDeveloperIdentity response
// from the store lookup outcome, carrying the model's response members only:
// the developer user identifiers, plus the matched identity ID and the page
// token when present.
func lookupDeveloperIdentityResult(matchedIdentityID string, devUserIDs []string, nextToken string) map[string]interface{} {
	result := map[string]interface{}{
		"DeveloperUserIdentifierList": devUserIDs,
	}
	if matchedIdentityID != "" {
		result["IdentityId"] = matchedIdentityID
	}
	if nextToken != "" {
		result["NextToken"] = nextToken
	}
	return result
}

// MergeDeveloperIdentities merges two developer user identities.
func (s *CognitoIdentityService) MergeDeveloperIdentities(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	destIdentityID, err := s.mergeDeveloperIdentitiesCore(reqCtx, MergeDeveloperIdentitiesInput{
		IdentityPoolID:            req.GetParam("IdentityPoolId"),
		DeveloperProviderName:     req.GetParam("DeveloperProviderName"),
		SourceUserIdentifier:      req.GetParam("SourceUserIdentifier"),
		DestinationUserIdentifier: req.GetParam("DestinationUserIdentifier"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"IdentityId": destIdentityID,
	}, nil
}

// UnlinkDeveloperIdentity unlinks a developer identity from a Cognito identity.
func (s *CognitoIdentityService) UnlinkDeveloperIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.unlinkDeveloperIdentityCore(reqCtx, UnlinkDeveloperIdentityInput{
		IdentityID:              req.GetParam("IdentityId"),
		IdentityPoolID:          req.GetParam("IdentityPoolId"),
		DeveloperProviderName:   req.GetParam("DeveloperProviderName"),
		DeveloperUserIdentifier: req.GetParam("DeveloperUserIdentifier"),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
