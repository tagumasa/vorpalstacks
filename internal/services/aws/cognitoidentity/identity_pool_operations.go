package cognitoidentity

import (
	"context"
	"errors"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
)

// CreateIdentityPool creates a new Cognito identity pool.
func (s *CognitoIdentityService) CreateIdentityPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolName := getParam(req, "IdentityPoolName")
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
	if providerName := getParam(req, "DeveloperProviderName"); providerName != "" {
		pool.DeveloperProviderName = providerName
	}
	if loginProviders := parseSupportedLoginProviders(req); len(loginProviders) > 0 {
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
		if err := store.TagResource(created.Arn, tags); err != nil {
			_ = store.DeleteIdentityPool(created.ID)
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

	tags, _ := store.ListTags(pool.Arn)
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

	pools, err := store.ListIdentityPools()
	if err != nil {
		return nil, ErrInternalError
	}

	maxResults := request.GetIntParam(req.Parameters, "MaxResults")
	if maxResults <= 0 || maxResults > 60 {
		maxResults = 60
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")

	identityPools := make([]map[string]interface{}, 0)
	started := nextToken == ""
	for _, pool := range pools {
		if !started {
			if pool.ID == nextToken {
				started = true
			}
			continue
		}
		identityPools = append(identityPools, map[string]interface{}{
			"IdentityPoolId":   pool.ID,
			"IdentityPoolName": pool.Name,
		})
		if len(identityPools) >= maxResults {
			break
		}
	}

	result := map[string]interface{}{
		"IdentityPools": identityPools,
	}
	if len(identityPools) >= maxResults && len(identityPools) > 0 {
		result["NextToken"] = identityPools[len(identityPools)-1]["IdentityPoolId"]
	}

	return result, nil
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

	if poolName := getParam(req, "IdentityPoolName"); poolName != "" {
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
	if providerName := getParam(req, "DeveloperProviderName"); providerName != "" {
		pool.DeveloperProviderName = providerName
	}
	if providers := parseCognitoIdentityProviders(req); len(providers) > 0 {
		pool.CognitoIdentityProviders = providers
	}
	if loginProviders := parseSupportedLoginProviders(req); len(loginProviders) > 0 {
		pool.SupportedLoginProviders = loginProviders
	}
	if oidcArns := getStringSliceParam(req, "OpenIdConnectProviderARNs"); len(oidcArns) > 0 {
		pool.OpenIdConnectProviderARNs = oidcArns
	}
	if samlArns := getStringSliceParam(req, "SamlProviderARNs"); len(samlArns) > 0 {
		pool.SamlProviderARNs = samlArns
	}

	if err := store.UpdateIdentityPool(pool); err != nil {
		return nil, ErrInternalError
	}

	return formatIdentityPool(pool), nil
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
	if roles, ok := req.Parameters["Roles"]; ok {
		if m, ok := roles.(map[string]interface{}); ok {
			if v, ok := m["authenticated"].(string); ok {
				authRole = v
			}
			if v, ok := m["unauthenticated"].(string); ok {
				unauthRole = v
			}
		}
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
	return getParam(req, "IdentityPoolId")
}
