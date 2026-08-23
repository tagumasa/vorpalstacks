package cognitoidentity

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
)

// DeleteIdentities deletes the identities from the specified identity pool.
func (s *CognitoIdentityService) DeleteIdentities(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identityIDs := getStringSliceParam(req, "IdentityIdsToDelete")
	if len(identityIDs) == 0 {
		return nil, ErrInvalidParameter
	}
	if len(identityIDs) > maxIdentityIdsToDelete {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var unprocessed []string
	for _, id := range identityIDs {
		identity, err := store.GetIdentityByID(id)
		if err != nil {
			unprocessed = append(unprocessed, id)
			continue
		}
		if err := store.DeleteIdentity(identity.IdentityPoolID, id); err != nil {
			unprocessed = append(unprocessed, id)
		}
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
	poolID := req.GetParam("IdentityPoolId")
	if !validateIdentityPoolId(poolID) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetIdentityPool(poolID); err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityPoolNotFound)
	}

	_, ok := req.Parameters["MaxResults"]
	if !ok {
		return nil, ErrInvalidParameter
	}
	maxResults := request.GetIntParam(req.Parameters, "MaxResults")
	if !validateQueryLimit(maxResults) {
		return nil, ErrInvalidParameter
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")
	if !validatePaginationKey(nextToken) {
		return nil, ErrInvalidParameter
	}
	// HideDisabled is accepted for SPEC compliance. Edge identities have no
	// disabled state, so the filter has no effect.
	_ = getBoolParam(req, "HideDisabled")

	identities, token, err := store.ListIdentitiesByPool(poolID, maxResults, nextToken)
	if err != nil {
		return nil, ErrInternalError
	}

	items := make([]map[string]interface{}, 0, len(identities))
	for _, identity := range identities {
		items = append(items, map[string]interface{}{
			"IdentityId":       identity.ID,
			"CreationDate":     identity.CreationDate.Unix(),
			"LastModifiedDate": identity.LastModifiedDate.Unix(),
			"Logins":           formatLoginKeys(identity.Logins),
		})
	}

	result := map[string]interface{}{
		"IdentityPoolId": poolID,
		"Identities":     items,
	}
	if token != "" {
		result["NextToken"] = token
	}

	return result, nil
}

// GetOpenIdToken gets an OpenID token for a Cognito identity.
func (s *CognitoIdentityService) GetOpenIdToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identityID := req.GetParam("IdentityId")
	if !validateIdentityId(identityID) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetIdentityByID(identityID)
	if err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityNotFound)
	}

	// Logins are accepted for wire compatibility but NOT persisted. In AWS,
	// GetOpenIdToken verifies the caller's provider tokens against the actual
	// identity provider before issuing a token. The edge environment cannot
	// perform external provider verification, so the parameter is accepted
	// without side effects to prevent identity takeover via Logins injection.
	if logins := parseMapParam(req, "Logins"); len(logins) > 0 {
		if !validateMapSize(len(logins), 10) || !validateLoginsKeys(logins) {
			return nil, ErrInvalidParameter
		}
		if !validateLoginsValues(logins) {
			return nil, ErrInvalidParameter
		}
	}

	token, err := s.tokenMgr.generateOpenIdToken(identityID, identity.IdentityPoolID, openIdTokenTTLSeconds, nil, nil)
	if err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"IdentityId": identityID,
		"Token":      token,
	}, nil
}

// GetOpenIdTokenForDeveloperIdentity registers (or retrieves) a developer identity and returns an OpenID token.
func (s *CognitoIdentityService) GetOpenIdTokenForDeveloperIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := req.GetParam("IdentityPoolId")
	if !validateIdentityPoolId(poolID) {
		return nil, ErrInvalidParameter
	}

	logins := parseMapParam(req, "Logins")
	if len(logins) == 0 {
		return nil, ErrInvalidParameter
	}
	if !validateLoginsValues(logins) {
		return nil, ErrInvalidParameter
	}

	// AWS expects exactly one Login entry per request (1 developer user).
	if len(logins) != 1 {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetIdentityPool(poolID); err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityPoolNotFound)
	}

	identityID := req.GetParam("IdentityId")
	if identityID != "" && !validateIdentityId(identityID) {
		return nil, ErrInvalidParameter
	}

	// TokenDuration controls token expiry (range 1-86400 seconds per AWS spec).
	tokenDuration := int64(developerTokenDefaultTTLSeconds)
	if _, ok := req.Parameters["TokenDuration"]; ok {
		td := int64(request.GetIntParam(req.Parameters, "TokenDuration"))
		if !validateTokenDuration(td) {
			return nil, ErrInvalidParameter
		}
		tokenDuration = td
	}

	// PrincipalTags are embedded into the JWT as cognito:principal_tags
	// so that STS AssumeRoleWithWebIdentity propagates them as session tags.
	var principalTags map[string]string
	if ptVal, ok := req.Parameters["PrincipalTags"]; ok {
		if ptMap, ok := ptVal.(map[string]interface{}); ok {
			if !validateMapSize(len(ptMap), 50) {
				return nil, ErrInvalidParameter
			}
			principalTags = make(map[string]string, len(ptMap))
			for k, v := range ptMap {
				if !validatePrincipalTagName(k) {
					return nil, ErrInvalidParameter
				}
				s, ok := v.(string)
				if !ok || !validatePrincipalTagValue(s) {
					return nil, ErrInvalidParameter
				}
				principalTags[k] = s
			}
		}
	}

	for providerName, devUserID := range logins {
		// The store resolves the developer identity under its key lock: an
		// existing link is reused (a differing supplied IdentityId maps to
		// DeveloperUserAlreadyRegisteredException), otherwise a fresh identity
		// is created and linked in one critical section.
		resolved, err := store.EnsureDeveloperIdentity(poolID, providerName, devUserID, identityID)
		if err != nil {
			if errors.Is(err, cognitoidentitystore.ErrDeveloperIdentityConflict) {
				return nil, ErrDeveloperUserAlreadyRegistered
			}
			if errors.Is(err, cognitoidentitystore.ErrIdentityNotFound) {
				return nil, ErrResourceNotFound
			}
			return nil, ErrInternalError
		}
		identityID = resolved
	}

	token, err := s.tokenMgr.generateOpenIdToken(identityID, poolID, tokenDuration, nil, principalTags)
	if err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"IdentityId": identityID,
		"Token":      token,
	}, nil
}

// GetPrincipalTagAttributeMap retrieves the principal tag attribute map for an identity provider.
func (s *CognitoIdentityService) GetPrincipalTagAttributeMap(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := req.GetParam("IdentityPoolId")
	if !validateIdentityPoolId(poolID) {
		return nil, ErrInvalidParameter
	}
	providerName := req.GetParam("IdentityProviderName")
	if providerName == "" || !validateIdentityProviderNameLength(providerName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetIdentityPool(poolID); err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityPoolNotFound)
	}

	ptam, err := store.GetPrincipalTagAttributeMap(poolID, providerName)
	if err != nil {
		if !errors.Is(err, cognitoidentitystore.ErrIdentityNotFound) {
			return nil, ErrInternalError
		}
		return map[string]interface{}{
			"IdentityPoolId":       poolID,
			"IdentityProviderName": providerName,
			"PrincipalTags":        map[string]string{},
			"UseDefaults":          true,
		}, nil
	}

	return map[string]interface{}{
		"IdentityPoolId":       ptam.IdentityPoolID,
		"IdentityProviderName": ptam.IdentityProviderName,
		"PrincipalTags":        ptam.PrincipalTags,
		"UseDefaults":          ptam.UseDefaults,
	}, nil
}

// SetPrincipalTagAttributeMap sets the principal tag attribute map for an identity provider.
func (s *CognitoIdentityService) SetPrincipalTagAttributeMap(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := req.GetParam("IdentityPoolId")
	if !validateIdentityPoolId(poolID) {
		return nil, ErrInvalidParameter
	}
	providerName := req.GetParam("IdentityProviderName")
	if providerName == "" || !validateIdentityProviderNameLength(providerName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetIdentityPool(poolID); err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityPoolNotFound)
	}

	principalTags := parseMapParam(req, "PrincipalTags")
	if len(principalTags) > 0 {
		if !validateMapSize(len(principalTags), 50) {
			return nil, ErrInvalidParameter
		}
		for k, v := range principalTags {
			if !validatePrincipalTagName(k) || !validatePrincipalTagValue(v) {
				return nil, ErrInvalidParameter
			}
		}
	}
	useDefaults := getBoolParam(req, "UseDefaults")

	if err := store.SetPrincipalTagAttributeMap(poolID, providerName, principalTags, useDefaults); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"IdentityPoolId":       poolID,
		"IdentityProviderName": providerName,
		"PrincipalTags":        principalTags,
		"UseDefaults":          useDefaults,
	}, nil
}

// LookupDeveloperIdentity looks up a developer identity identifier and returns the mapped identity IDs.
func (s *CognitoIdentityService) LookupDeveloperIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := req.GetParam("IdentityPoolId")
	if !validateIdentityPoolId(poolID) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetIdentityPool(poolID); err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityPoolNotFound)
	}

	identityID := req.GetParam("IdentityId")
	if identityID != "" && !validateIdentityId(identityID) {
		return nil, ErrInvalidParameter
	}
	devUserID := req.GetParam("DeveloperUserIdentifier")
	if identityID == "" && devUserID == "" {
		return nil, ErrInvalidParameter
	}
	if devUserID != "" && !validateDeveloperUserIdentifier(devUserID) {
		return nil, ErrInvalidParameter
	}
	maxResults := defaultLookupMaxResults
	if _, ok := req.Parameters["MaxResults"]; ok {
		n := request.GetIntParam(req.Parameters, "MaxResults")
		if !validateQueryLimit(n) {
			return nil, ErrInvalidParameter
		}
		maxResults = n
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")
	if !validatePaginationKey(nextToken) {
		return nil, ErrInvalidParameter
	}

	matchedIdentityID, devUserIDs, nextTokenOut, err := store.LookupDeveloperIdentity(poolID, identityID, devUserID, maxResults, nextToken)
	if err != nil {
		return nil, ErrInternalError
	}

	return lookupDeveloperIdentityResult(matchedIdentityID, devUserIDs, nextTokenOut), nil
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
	poolID := req.GetParam("IdentityPoolId")
	if !validateIdentityPoolId(poolID) {
		return nil, ErrInvalidParameter
	}
	providerName := req.GetParam("DeveloperProviderName")
	if providerName == "" || !validateDeveloperProviderName(providerName) {
		return nil, ErrInvalidParameter
	}
	sourceUserID := req.GetParam("SourceUserIdentifier")
	if sourceUserID == "" || !validateDeveloperUserIdentifier(sourceUserID) {
		return nil, ErrInvalidParameter
	}
	destUserID := req.GetParam("DestinationUserIdentifier")
	if destUserID == "" || !validateDeveloperUserIdentifier(destUserID) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// The store performs the merge under the pool lock with the developer
	// identity link moving before any identity record is destroyed.
	destIdentityID, err := store.MergeDeveloperIdentities(poolID, providerName, sourceUserID, destUserID)
	if err != nil {
		if errors.Is(err, cognitoidentitystore.ErrIdentityNotFound) {
			return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityNotFound)
		}
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"IdentityId": destIdentityID,
	}, nil
}

// UnlinkDeveloperIdentity unlinks a developer identity from a Cognito identity.
func (s *CognitoIdentityService) UnlinkDeveloperIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identityID := req.GetParam("IdentityId")
	if !validateIdentityId(identityID) {
		return nil, ErrInvalidParameter
	}
	poolID := req.GetParam("IdentityPoolId")
	if !validateIdentityPoolId(poolID) {
		return nil, ErrInvalidParameter
	}
	providerName := req.GetParam("DeveloperProviderName")
	if providerName == "" || !validateDeveloperProviderName(providerName) {
		return nil, ErrInvalidParameter
	}
	devUserID := req.GetParam("DeveloperUserIdentifier")
	if devUserID == "" || !validateDeveloperUserIdentifier(devUserID) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.UnlinkDeveloperIdentity(poolID, providerName, devUserID); err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityNotFound)
	}

	return response.EmptyResponse(), nil
}

// UnlinkIdentity unlinks login providers from a Cognito identity.
func (s *CognitoIdentityService) UnlinkIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identityID := req.GetParam("IdentityId")
	if !validateIdentityId(identityID) {
		return nil, ErrInvalidParameter
	}

	loginsToRemove := getStringSliceParam(req, "LoginsToRemove")
	if len(loginsToRemove) == 0 {
		return nil, ErrInvalidParameter
	}

	// Logins provides the caller's provider tokens for authorization. AWS
	// requires at least one provider token matching the identity's linked
	// providers before allowing an unlink operation.
	logins := parseMapParam(req, "Logins")
	if len(logins) == 0 {
		return nil, ErrNotAuthorized
	}
	if !validateMapSize(len(logins), 10) || !validateLoginsKeys(logins) {
		return nil, ErrInvalidParameter
	}
	if !validateLoginsValues(logins) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetIdentityByID(identityID)
	if err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityNotFound)
	}

	// Verify the caller holds a token for at least one provider linked to
	// this identity. Both the provider name AND the token value must match
	// to prevent impersonation via provider-name-only knowledge.
	providerMatch := false
	for provider, tokenValue := range logins {
		if storedValue, exists := identity.Logins[provider]; exists && storedValue == tokenValue {
			providerMatch = true
			break
		}
	}
	if !providerMatch {
		return nil, ErrNotAuthorized
	}

	if err := store.UnlinkLogins(identity.IdentityPoolID, identityID, loginsToRemove); err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityNotFound)
	}

	return response.EmptyResponse(), nil
}
