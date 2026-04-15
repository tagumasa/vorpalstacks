package cognitoidentity

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"

	"github.com/google/uuid"
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
	poolID := getParam(req, "IdentityPoolId")
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

	maxResults := request.GetIntParam(req.Parameters, "MaxResults")
	if maxResults <= 0 || maxResults > 60 {
		maxResults = 60
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")

	identities, token, err := store.ListIdentitiesByPool(poolID, maxResults, nextToken)
	if err != nil {
		return nil, ErrInternalError
	}

	items := make([]map[string]interface{}, 0, len(identities))
	for _, identity := range identities {
		logins := make([]string, 0, len(identity.Logins))
		for k := range identity.Logins {
			logins = append(logins, k)
		}
		items = append(items, map[string]interface{}{
			"IdentityId":       identity.ID,
			"CreationDate":     identity.CreationDate.Unix(),
			"LastModifiedDate": identity.LastModifiedDate.Unix(),
			"Logins":           logins,
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
	identityID := getParam(req, "IdentityId")
	if identityID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetIdentityByID(identityID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if logins := parseLogins(req); len(logins) > 0 {
		identity, err := store.GetIdentityByID(identityID)
		if err != nil {
			return nil, ErrResourceNotFound
		}
		for k, v := range logins {
			identity.Logins[k] = v
		}
		key := cognitoidentitystore.IdentityPoolIdentityKey(identity.IdentityPoolID, identity.ID)
		if err := store.Identities().Put(key, identity); err != nil {
			return nil, fmt.Errorf("failed to update identity logins: %w", err)
		}
	}

	token := uuid.New().String()

	return map[string]interface{}{
		"IdentityId": identityID,
		"Token":      token,
	}, nil
}

// GetOpenIdTokenForDeveloperIdentity registers (or retrieves) a developer identity and returns an OpenID token.
func (s *CognitoIdentityService) GetOpenIdTokenForDeveloperIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := getParam(req, "IdentityPoolId")
	if poolID == "" {
		return nil, ErrInvalidParameter
	}

	logins := parseLogins(req)
	if len(logins) == 0 {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetIdentityPool(poolID); err != nil {
		return nil, ErrResourceNotFound
	}

	identityID := getParam(req, "IdentityId")

	for providerName, devUserID := range logins {
		existing, err := store.GetDeveloperIdentity(poolID, providerName, devUserID)
		if err == nil && existing.IdentityID != "" {
			identityID = existing.IdentityID
			continue
		}

		if identityID == "" {
			identity := cognitoidentitystore.NewIdentity(poolID)
			identityID = identity.ID
			if err := store.CreateIdentity(identity); err != nil {
				return nil, fmt.Errorf("failed to create identity: %w", err)
			}
		}

		di := &cognitoidentitystore.DeveloperIdentity{
			DeveloperUserIdentifier: devUserID,
			DeveloperProviderName:   providerName,
			IdentityPoolID:          poolID,
			IdentityID:              identityID,
		}
		if err := store.LinkDeveloperIdentity(di); err != nil {
			return nil, fmt.Errorf("failed to link developer identity: %w", err)
		}
	}

	token := uuid.New().String()

	return map[string]interface{}{
		"IdentityId": identityID,
		"Token":      token,
	}, nil
}

// GetPrincipalTagAttributeMap retrieves the principal tag attribute map for an identity provider.
func (s *CognitoIdentityService) GetPrincipalTagAttributeMap(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := getParam(req, "IdentityPoolId")
	if poolID == "" {
		return nil, ErrInvalidParameter
	}
	providerName := getParam(req, "IdentityProviderName")
	if providerName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
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
	poolID := getParam(req, "IdentityPoolId")
	if poolID == "" {
		return nil, ErrInvalidParameter
	}
	providerName := getParam(req, "IdentityProviderName")
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
	poolID := getParam(req, "IdentityPoolId")
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

	identityID := getParam(req, "IdentityId")
	devUserID := getParam(req, "DeveloperUserIdentifier")
	maxResults := request.GetIntParam(req.Parameters, "MaxResults")
	if maxResults <= 0 {
		maxResults = 60
	}

	matchedIdentityID, devUserIDs, _, err := store.LookupDeveloperIdentity(poolID, identityID, devUserID, maxResults)
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

	return result, nil
}

// MergeDeveloperIdentities merges two developer user identities.
func (s *CognitoIdentityService) MergeDeveloperIdentities(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := getParam(req, "IdentityPoolId")
	if poolID == "" {
		return nil, ErrInvalidParameter
	}
	sourceUserID := getParam(req, "SourceUserIdentifier")
	if sourceUserID == "" {
		return nil, ErrInvalidParameter
	}
	destUserID := getParam(req, "DestinationUserIdentifier")
	if destUserID == "" {
		return nil, ErrInvalidParameter
	}
	providerName := getParam(req, "DeveloperProviderName")
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
		sourceKey := cognitoidentitystore.IdentityPoolIdentityKey(poolID, sourceDI.IdentityID)
		if err := store.Identities().Delete(sourceKey); err != nil {
			return nil, fmt.Errorf("failed to delete source identity: %w", err)
		}
	}

	destDI, err = store.GetDeveloperIdentity(poolID, providerName, destUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to re-fetch destination developer identity after merge: %w", err)
	}
	if err := store.LinkDeveloperIdentity(&cognitoidentitystore.DeveloperIdentity{
		DeveloperUserIdentifier: sourceUserID,
		DeveloperProviderName:   providerName,
		IdentityPoolID:          poolID,
		IdentityID:              destDI.IdentityID,
	}); err != nil {
		return nil, fmt.Errorf("failed to link merged developer identity: %w", err)
	}

	return map[string]interface{}{
		"IdentityId": destDI.IdentityID,
	}, nil
}

// UnlinkDeveloperIdentity unlinks a developer identity from a Cognito identity.
func (s *CognitoIdentityService) UnlinkDeveloperIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identityID := getParam(req, "IdentityId")
	if identityID == "" {
		return nil, ErrInvalidParameter
	}
	poolID := getParam(req, "IdentityPoolId")
	if poolID == "" {
		return nil, ErrInvalidParameter
	}
	providerName := getParam(req, "DeveloperProviderName")
	if providerName == "" {
		return nil, ErrInvalidParameter
	}
	devUserID := getParam(req, "DeveloperUserIdentifier")
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
	identityID := getParam(req, "IdentityId")
	if identityID == "" {
		return nil, ErrInvalidParameter
	}

	loginsToRemove := getStringSliceParam(req, "LoginsToRemove")
	if len(loginsToRemove) == 0 {
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

	if err := store.UnlinkLogins(identity.IdentityPoolID, identityID, loginsToRemove); err != nil {
		return nil, ErrResourceNotFound
	}

	return response.EmptyResponse(), nil
}
