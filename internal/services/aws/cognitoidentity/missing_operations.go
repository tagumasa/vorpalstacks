package cognitoidentity

import (
	"context"
	"time"

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
	if poolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetIdentityPool(poolID); err != nil {
		return nil, ErrResourceNotFound
	}

	maxResults := 60
	if _, ok := req.Parameters["MaxResults"]; ok {
		n := request.GetIntParam(req.Parameters, "MaxResults")
		if !validateQueryLimit(n) {
			return nil, ErrInvalidParameter
		}
		maxResults = n
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")
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
	if identityID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetIdentityByID(identityID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	// Logins are accepted for wire compatibility but NOT persisted. In AWS,
	// GetOpenIdToken verifies the caller's provider tokens against the actual
	// identity provider before issuing a token. The edge environment cannot
	// perform external provider verification, so the parameter is accepted
	// without side effects to prevent identity takeover via Logins injection.
	_ = parseMapParam(req, "Logins")

	token, err := s.tokenMgr.generateOpenIdToken(identityID, identity.IdentityPoolID, 900, nil, nil)
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
	if poolID == "" {
		return nil, ErrInvalidParameter
	}

	logins := parseMapParam(req, "Logins")
	if len(logins) == 0 {
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
		return nil, ErrResourceNotFound
	}

	identityID := req.GetParam("IdentityId")

	// TokenDuration controls token expiry (range 1-86400 seconds per AWS spec).
	tokenDuration := int64(900) // default 15 minutes
	if _, ok := req.Parameters["TokenDuration"]; ok {
		td := request.GetIntParam(req.Parameters, "TokenDuration")
		if td < 1 || td > 86400 {
			return nil, ErrInvalidParameter
		}
		tokenDuration = int64(td)
	}

	// PrincipalTags are embedded into the JWT as cognito:principal_tags
	// so that STS AssumeRoleWithWebIdentity propagates them as session tags.
	var principalTags map[string]string
	if ptVal, ok := req.Parameters["PrincipalTags"]; ok {
		if ptMap, ok := ptVal.(map[string]interface{}); ok {
			if len(ptMap) > 50 {
				return nil, ErrInvalidParameter
			}
			principalTags = make(map[string]string, len(ptMap))
			for k, v := range ptMap {
				if s, ok := v.(string); ok {
					principalTags[k] = s
				}
			}
		}
	}

	for providerName, devUserID := range logins {
		existing, err := store.GetDeveloperIdentity(poolID, providerName, devUserID)
		if err == nil && existing.IdentityID != "" {
			// The developer identity already exists. If the caller supplied
			// an IdentityId that differs from the existing link, AWS returns
			// DeveloperUserAlreadyRegisteredException.
			if identityID != "" && existing.IdentityID != identityID {
				return nil, ErrDeveloperUserAlreadyRegistered
			}
			identityID = existing.IdentityID
			break
		}

		if identityID == "" {
			identity := cognitoidentitystore.NewIdentity(poolID)
			identityID = identity.ID
			if err := store.CreateIdentity(identity); err != nil {
				return nil, ErrInternalError
			}
		}

		di := &cognitoidentitystore.DeveloperIdentity{
			DeveloperUserIdentifier: devUserID,
			DeveloperProviderName:   providerName,
			IdentityPoolID:          poolID,
			IdentityID:              identityID,
		}
		if err := store.LinkDeveloperIdentity(di); err != nil {
			return nil, ErrInternalError
		}
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
	if poolID == "" {
		return nil, ErrInvalidParameter
	}
	providerName := req.GetParam("IdentityProviderName")
	if providerName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetIdentityPool(poolID); err != nil {
		return nil, ErrResourceNotFound
	}

	ptam, err := store.GetPrincipalTagAttributeMap(poolID, providerName)
	if err != nil {
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
	if poolID == "" {
		return nil, ErrInvalidParameter
	}
	providerName := req.GetParam("IdentityProviderName")
	if providerName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetIdentityPool(poolID); err != nil {
		return nil, ErrResourceNotFound
	}

	principalTags := parseMapParam(req, "PrincipalTags")
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
	if poolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetIdentityPool(poolID); err != nil {
		return nil, ErrResourceNotFound
	}

	identityID := req.GetParam("IdentityId")
	devUserID := req.GetParam("DeveloperUserIdentifier")
	maxResults := 60
	if _, ok := req.Parameters["MaxResults"]; ok {
		n := request.GetIntParam(req.Parameters, "MaxResults")
		if !validateQueryLimit(n) {
			return nil, ErrInvalidParameter
		}
		maxResults = n
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")

	matchedIdentityID, devUserIDs, nextTokenOut, err := store.LookupDeveloperIdentity(poolID, identityID, devUserID, maxResults, nextToken)
	if err != nil {
		return nil, ErrInternalError
	}

	result := map[string]interface{}{
		"IdentityPoolId":              poolID,
		"DeveloperUserIdentifierList": devUserIDs,
	}
	if matchedIdentityID != "" {
		result["IdentityId"] = matchedIdentityID
	}
	if nextTokenOut != "" {
		result["NextToken"] = nextTokenOut
	}

	return result, nil
}

// MergeDeveloperIdentities merges two developer user identities.
func (s *CognitoIdentityService) MergeDeveloperIdentities(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := req.GetParam("IdentityPoolId")
	if poolID == "" {
		return nil, ErrInvalidParameter
	}
	sourceUserID := req.GetParam("SourceUserIdentifier")
	if sourceUserID == "" {
		return nil, ErrInvalidParameter
	}
	destUserID := req.GetParam("DestinationUserIdentifier")
	if destUserID == "" {
		return nil, ErrInvalidParameter
	}
	providerName := req.GetParam("DeveloperProviderName")
	if providerName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sourceDI, err := store.GetDeveloperIdentity(poolID, providerName, sourceUserID)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	destDI, err := store.GetDeveloperIdentity(poolID, providerName, destUserID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if sourceDI.IdentityID != "" && destDI.IdentityID != "" && sourceDI.IdentityID != destDI.IdentityID {
		// Merge the source identity's logins into the destination identity so that
		// public provider links (Facebook, Google, etc.) are not lost.
		sourceIdentity, srcErr := store.GetIdentity(poolID, sourceDI.IdentityID)
		if srcErr == nil && sourceIdentity != nil {
			destIdentity, dstErr := store.GetIdentity(poolID, destDI.IdentityID)
			if dstErr == nil && destIdentity != nil {
				if destIdentity.Logins == nil {
					destIdentity.Logins = make(map[string]string)
				}
				for provider, token := range sourceIdentity.Logins {
					if _, exists := destIdentity.Logins[provider]; !exists {
						destIdentity.Logins[provider] = token
					}
				}
				destIdentity.LastModifiedDate = time.Now().UTC()
				destKey := cognitoidentitystore.IdentityPoolIdentityKey(poolID, destDI.IdentityID)
				_ = store.Identities().Put(destKey, destIdentity)
			}
		}

		sourceKey := cognitoidentitystore.IdentityPoolIdentityKey(poolID, sourceDI.IdentityID)
		if err := store.Identities().Delete(sourceKey); err != nil {
			return nil, ErrInternalError
		}
	}

	destDI, err = store.GetDeveloperIdentity(poolID, providerName, destUserID)
	if err != nil {
		return nil, ErrInternalError
	}
	if err := store.LinkDeveloperIdentity(&cognitoidentitystore.DeveloperIdentity{
		DeveloperUserIdentifier: sourceUserID,
		DeveloperProviderName:   providerName,
		IdentityPoolID:          poolID,
		IdentityID:              destDI.IdentityID,
	}); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"IdentityId": destDI.IdentityID,
	}, nil
}

// UnlinkDeveloperIdentity unlinks a developer identity from a Cognito identity.
func (s *CognitoIdentityService) UnlinkDeveloperIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identityID := req.GetParam("IdentityId")
	if identityID == "" {
		return nil, ErrInvalidParameter
	}
	poolID := req.GetParam("IdentityPoolId")
	if poolID == "" {
		return nil, ErrInvalidParameter
	}
	providerName := req.GetParam("DeveloperProviderName")
	if providerName == "" {
		return nil, ErrInvalidParameter
	}
	devUserID := req.GetParam("DeveloperUserIdentifier")
	if devUserID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.UnlinkDeveloperIdentity(poolID, providerName, devUserID); err != nil {
		return nil, ErrResourceNotFound
	}

	return response.EmptyResponse(), nil
}

// UnlinkIdentity unlinks login providers from a Cognito identity.
func (s *CognitoIdentityService) UnlinkIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identityID := req.GetParam("IdentityId")
	if identityID == "" {
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetIdentityByID(identityID)
	if err != nil {
		return nil, ErrResourceNotFound
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
		return nil, ErrResourceNotFound
	}

	return response.EmptyResponse(), nil
}
