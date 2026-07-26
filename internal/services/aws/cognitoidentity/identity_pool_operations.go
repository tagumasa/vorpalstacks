package cognitoidentity

import (
	"context"
	"errors"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// CreateIdentityPool creates a new Cognito identity pool.
func (s *CognitoIdentityService) CreateIdentityPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolName := req.GetParam("IdentityPoolName")
	if poolName == "" {
		return nil, ErrInvalidParameter
	}

	allowUnauthenticated := getBoolParam(req, "AllowUnauthenticatedIdentities")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pool := cognitoidentitystore.NewIdentityPool(poolName, allowUnauthenticated, reqCtx.GetRegion())

	if providers := parseCognitoIdentityProviders(req); len(providers) > 0 {
		pool.CognitoIdentityProviders = providers
	}
	if providerName := req.GetParam("DeveloperProviderName"); providerName != "" {
		pool.DeveloperProviderName = providerName
	}
	if loginProviders := parseMapParam(req, "SupportedLoginProviders"); len(loginProviders) > 0 {
		pool.SupportedLoginProviders = loginProviders
	}
	if oidcArns := getStringSliceParam(req, "OpenIdConnectProviderARNs"); len(oidcArns) > 0 {
		pool.OpenIdConnectProviderARNs = oidcArns
	}
	if samlArns := getStringSliceParam(req, "SamlProviderARNs"); len(samlArns) > 0 {
		pool.SamlProviderARNs = samlArns
	}
	if allowClassic, ok := req.Parameters["AllowClassicFlow"]; ok {
		if b, ok := allowClassic.(bool); ok {
			pool.AllowClassicFlow = b
		}
	}

	created, err := store.CreateIdentityPool(pool)
	if err != nil {
		if errors.Is(err, cognitoidentitystore.ErrIdentityPoolAlreadyExists) {
			return nil, ErrResourceInUse
		}
		return nil, ErrInternalError
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "IdentityPoolTags"))
	if len(tags) > 0 {
		if err := store.Tag(created.Arn, tags); err != nil {
			logs.Error("Failed to tag identity pool, attempting cleanup", logs.String("poolId", created.ID), logs.Err(err))
			if delErr := store.DeleteIdentityPool(created.ID); delErr != nil {
				logs.Error("Failed to cleanup identity pool after tag failure", logs.String("poolId", created.ID), logs.Err(delErr))
			}
			return nil, ErrInternalError
		}
	}

	if len(tags) > 0 {
		return formatIdentityPoolWithTags(created, tags), nil
	}
	return formatIdentityPool(created), nil
}

// DescribeIdentityPool returns details about a Cognito identity pool.
func (s *CognitoIdentityService) DescribeIdentityPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := getIdentityPoolID(req)
	if poolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pool, err := store.GetIdentityPool(poolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	tags, _ := store.List(pool.Arn)
	if len(tags) > 0 {
		return formatIdentityPoolWithTags(pool, tags), nil
	}

	return formatIdentityPool(pool), nil
}

// DeleteIdentityPool deletes a Cognito identity pool.
func (s *CognitoIdentityService) DeleteIdentityPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := getIdentityPoolID(req)
	if poolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteIdentityPool(poolID); err != nil {
		return nil, ErrResourceNotFound
	}

	return response.EmptyResponse(), nil
}

// ListIdentityPools returns a list of Cognito identity pools.
func (s *CognitoIdentityService) ListIdentityPools(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	maxResults := request.GetIntParam(req.Parameters, "MaxResults")
	if maxResults <= 0 || maxResults > 60 {
		maxResults = 60
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")

	opts := storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   nextToken,
	}

	result, err := store.ListIdentityPools(opts)
	if err != nil {
		return nil, ErrInternalError
	}

	identityPools := make([]map[string]interface{}, 0, len(result.Items))
	for _, pool := range result.Items {
		identityPools = append(identityPools, map[string]interface{}{
			"IdentityPoolId":   pool.ID,
			"IdentityPoolName": pool.Name,
		})
	}

	resp := map[string]interface{}{
		"IdentityPools": identityPools,
	}
	if result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}

	return resp, nil
}

// UpdateIdentityPool updates a Cognito identity pool.
func (s *CognitoIdentityService) UpdateIdentityPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := getIdentityPoolID(req)
	if poolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pool, err := store.GetIdentityPool(poolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if poolName := req.GetParam("IdentityPoolName"); poolName != "" {
		pool.Name = poolName
	}
	if allowUnauth, ok := req.Parameters["AllowUnauthenticatedIdentities"]; ok {
		if b, ok := allowUnauth.(bool); ok {
			pool.AllowUnauthenticatedIdentities = b
		}
	}
	if allowClassic, ok := req.Parameters["AllowClassicFlow"]; ok {
		if b, ok := allowClassic.(bool); ok {
			pool.AllowClassicFlow = b
		}
	}
	if providerName := req.GetParam("DeveloperProviderName"); providerName != "" {
		pool.DeveloperProviderName = providerName
	}
	if providers := parseCognitoIdentityProviders(req); len(providers) > 0 {
		pool.CognitoIdentityProviders = providers
	}
	if loginProviders := parseMapParam(req, "SupportedLoginProviders"); len(loginProviders) > 0 {
		pool.SupportedLoginProviders = loginProviders
	}
	if oidcArns := getStringSliceParam(req, "OpenIdConnectProviderARNs"); len(oidcArns) > 0 {
		pool.OpenIdConnectProviderARNs = oidcArns
	}
	if samlArns := getStringSliceParam(req, "SamlProviderARNs"); len(samlArns) > 0 {
		pool.SamlProviderARNs = samlArns
	}

	var updatedTags map[string]string
	if _, ok := req.Parameters["IdentityPoolTags"]; ok {
		updatedTags = tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "IdentityPoolTags"))
		existingTags, _ := store.List(pool.Arn)
		var keysToRemove []string
		removedTags := make(map[string]string)
		for k, v := range existingTags {
			if _, keep := updatedTags[k]; !keep {
				keysToRemove = append(keysToRemove, k)
				removedTags[k] = v
			}
		}
		if len(keysToRemove) > 0 {
			if err := store.Untag(pool.Arn, keysToRemove); err != nil {
				return nil, ErrInternalError
			}
		}
		if err := store.Tag(pool.Arn, updatedTags); err != nil {
			// Rollback: re-apply the tags that were removed to avoid
			// leaving the resource in a partially-untagged state.
			if len(removedTags) > 0 {
				_ = store.Tag(pool.Arn, removedTags)
			}
			return nil, ErrInternalError
		}
		pool.Tags = updatedTags
	} else {
		updatedTags, _ = store.List(pool.Arn)
	}

	pool.LastModifiedDate = time.Now().UTC()

	if err := store.UpdateIdentityPool(pool); err != nil {
		return nil, ErrInternalError
	}

	return formatIdentityPoolWithTags(pool, updatedTags), nil
}

// GetIdentityPoolRoles returns the roles for a Cognito identity pool.
func (s *CognitoIdentityService) GetIdentityPoolRoles(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := getIdentityPoolID(req)
	if poolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	authRole, unauthRole, mappings, err := store.GetIdentityPoolRoles(poolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	result := map[string]interface{}{
		"IdentityPoolId": poolID,
	}

	roles := map[string]interface{}{}
	if authRole != "" {
		roles["authenticated"] = authRole
	}
	if unauthRole != "" {
		roles["unauthenticated"] = unauthRole
	}
	if len(roles) > 0 {
		result["Roles"] = roles
	}
	if len(mappings) > 0 {
		result["RoleMappings"] = formatRoleMappings(mappings)
	}

	return result, nil
}

// SetIdentityPoolRoles sets the roles for a Cognito identity pool.
func (s *CognitoIdentityService) SetIdentityPoolRoles(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := getIdentityPoolID(req)
	if poolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	authRole, unauthRole := "", ""
	if rolesVal, ok := req.Parameters["Roles"]; ok {
		if m, ok := rolesVal.(map[string]interface{}); ok {
			if v, ok := m["authenticated"].(string); ok {
				authRole = v
			}
			if v, ok := m["unauthenticated"].(string); ok {
				unauthRole = v
			}
		}
	} else {
		// Roles is semantically required by AWS. Absent Roles would silently
		// clear all existing roles — a destructive operation that AWS rejects
		// with InvalidParameterException.
		return nil, ErrInvalidParameter
	}
	if authRole == "" && unauthRole == "" {
		return nil, ErrInvalidParameter
	}
	mappings := parseRoleMappings(req)

	if err := store.SetIdentityPoolRoles(poolID, authRole, unauthRole, mappings); err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{
		"IdentityPoolId": poolID,
	}, nil
}

func getIdentityPoolID(req *request.ParsedRequest) string {
	return req.GetParam("IdentityPoolId")
}
