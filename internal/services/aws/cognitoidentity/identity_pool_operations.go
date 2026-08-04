package cognitoidentity

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

// CreateIdentityPool creates a new Cognito identity pool.
func (s *CognitoIdentityService) CreateIdentityPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	input := CreateIdentityPoolInput{
		IdentityPoolName:               req.GetParam("IdentityPoolName"),
		AllowUnauthenticatedIdentities: getBoolParam(req, "AllowUnauthenticatedIdentities"),
		Region:                         reqCtx.GetRegion(),
	}

	if providers, perr := parseCognitoIdentityProviders(req); perr != nil {
		return nil, perr
	} else {
		for _, p := range providers {
			input.CognitoIdentityProviders = append(input.CognitoIdentityProviders, ProviderOut{
				ProviderName:         p.ProviderName,
				ClientID:             p.ClientID,
				ServerSideTokenCheck: p.ServerSideTokenCheck,
			})
		}
	}

	input.DeveloperProviderName = req.GetParam("DeveloperProviderName")
	input.SupportedLoginProviders = parseMapParam(req, "SupportedLoginProviders")
	input.OpenIdConnectProviderARNs = getStringSliceParam(req, "OpenIdConnectProviderARNs")
	input.SamlProviderARNs = getStringSliceParam(req, "SamlProviderARNs")

	if allowClassic, ok := req.Parameters["AllowClassicFlow"]; ok {
		if b, ok := allowClassic.(bool); ok {
			input.AllowClassicFlow = b
			input.AllowClassicFlowProvided = true
		}
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "IdentityPoolTags"))
	if len(tags) > 0 {
		input.Tags = tags
		input.TagsProvided = true
	}

	result, err := s.createIdentityPoolCore(store, input)
	if err != nil {
		return nil, err
	}

	return poolOutToHTTP(result), nil
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

	result, err := s.describeIdentityPoolCore(store, poolID)
	if err != nil {
		return nil, err
	}
	return poolOutToHTTP(result), nil
}

// DeleteIdentityPool deletes a Cognito identity pool.
func (s *CognitoIdentityService) DeleteIdentityPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := getIdentityPoolID(req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteIdentityPoolCore(store, poolID); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListIdentityPools returns a list of Cognito identity pools.
func (s *CognitoIdentityService) ListIdentityPools(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	items, nextToken, err := s.listIdentityPoolsShortCore(store, ListIdentityPoolsInput{
		MaxResults: request.GetIntParam(req.Parameters, "MaxResults"),
		NextToken:  request.GetStringParam(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	identityPools := make([]map[string]interface{}, 0, len(items))
	for _, pool := range items {
		identityPools = append(identityPools, map[string]interface{}{
			"IdentityPoolId":   pool.ID,
			"IdentityPoolName": pool.Name,
		})
	}

	resp := map[string]interface{}{
		"IdentityPools": identityPools,
	}
	if nextToken != "" {
		resp["NextToken"] = nextToken
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

	// IdentityPoolName is @required in the Smithy IdentityPool shape.
	poolName := req.GetParam("IdentityPoolName")
	if poolName == "" {
		return nil, ErrInvalidParameter
	}
	if err := validateIdentityPoolName(poolName); err != nil {
		return nil, err
	}
	pool.Name = poolName

	// AllowUnauthenticatedIdentities is @required in the Smithy shape.
	allowUnauthVal, hasUnauth := req.Parameters["AllowUnauthenticatedIdentities"]
	if !hasUnauth {
		return nil, ErrInvalidParameter
	}
	if b, ok := allowUnauthVal.(bool); ok {
		pool.AllowUnauthenticatedIdentities = b
	} else {
		return nil, ErrInvalidParameter
	}

	if allowClassic, ok := req.Parameters["AllowClassicFlow"]; ok {
		if b, ok := allowClassic.(bool); ok {
			pool.AllowClassicFlow = b
		}
	}
	if providerName := req.GetParam("DeveloperProviderName"); providerName != "" {
		if err := validateDeveloperProviderName(providerName); err != nil {
			return nil, err
		}
		pool.DeveloperProviderName = providerName
	}
	if providers, err := parseCognitoIdentityProviders(req); err != nil {
		return nil, err
	} else if len(providers) > 0 {
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

	pool.Tags = updatedTags
	return poolOutToHTTP(poolToOut(pool)), nil
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
		rolesMap, ok := rolesVal.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}
		if err := validateRoleKeys(rolesMap); err != nil {
			return nil, err
		}
		if v, ok := rolesMap["authenticated"].(string); ok {
			authRole = v
		}
		if v, ok := rolesMap["unauthenticated"].(string); ok {
			unauthRole = v
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
	mappings, err := parseRoleMappings(req)
	if err != nil {
		return nil, err
	}

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
