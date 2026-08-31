package cognitoidentity

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// DeleteIdentities deletes the identities from the specified identity pool.
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
		for _, id := range unprocessed {
			items = append(items, map[string]interface{}{
				"IdentityId": id,
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

// GetPrincipalTagAttributeMap retrieves the principal tag attribute map for an identity provider.
func (s *CognitoIdentityService) GetPrincipalTagAttributeMap(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.getPrincipalTagAttributeMapCore(reqCtx, GetPrincipalTagAttributeMapInput{
		IdentityPoolID:       req.GetParam("IdentityPoolId"),
		IdentityProviderName: req.GetParam("IdentityProviderName"),
	})
	if err != nil {
		return nil, err
	}

	return principalTagAttributeMapToHTTP(result), nil
}

// SetPrincipalTagAttributeMap sets the principal tag attribute map for an identity provider.
func (s *CognitoIdentityService) SetPrincipalTagAttributeMap(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.setPrincipalTagAttributeMapCore(reqCtx, SetPrincipalTagAttributeMapInput{
		IdentityPoolID:       req.GetParam("IdentityPoolId"),
		IdentityProviderName: req.GetParam("IdentityProviderName"),
		PrincipalTags:        parseMapParam(req, "PrincipalTags"),
		UseDefaults:          getBoolParam(req, "UseDefaults"),
	})
	if err != nil {
		return nil, err
	}

	return principalTagAttributeMapToHTTP(result), nil
}

// principalTagAttributeMapToHTTP serialises a PrincipalTagAttributeMapResult
// into the member format shared by the Get/SetPrincipalTagAttributeMap
// responses.
func principalTagAttributeMapToHTTP(r *PrincipalTagAttributeMapResult) map[string]interface{} {
	return map[string]interface{}{
		"IdentityPoolId":       r.IdentityPoolID,
		"IdentityProviderName": r.IdentityProviderName,
		"PrincipalTags":        r.PrincipalTags,
		"UseDefaults":          r.UseDefaults,
	}
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
