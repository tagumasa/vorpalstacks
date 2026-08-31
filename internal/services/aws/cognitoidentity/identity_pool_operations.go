package cognitoidentity

import (
	"context"

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

	providers, perr := parseCognitoIdentityProviders(req.Parameters["CognitoIdentityProviders"])
	if perr != nil {
		return nil, perr
	}
	input.CognitoIdentityProviders = providers

	input.DeveloperProviderName = req.GetParam("DeveloperProviderName")
	input.SupportedLoginProviders = parseMapParam(req, "SupportedLoginProviders")
	input.OpenIdConnectProviderARNs = getStringSliceParam(req, "OpenIdConnectProviderARNs")
	input.SamlProviderARNs = getStringSliceParam(req, "SamlProviderARNs")

	if allowClassic, ok := req.Parameters["AllowClassicFlow"]; ok {
		b, ok := allowClassic.(bool)
		if !ok {
			return nil, ErrInvalidParameter
		}
		input.AllowClassicFlow = b
		input.AllowClassicFlowProvided = true
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.describeIdentityPoolCore(store, getIdentityPoolID(req))
	if err != nil {
		return nil, err
	}
	return poolOutToHTTP(result), nil
}

// DeleteIdentityPool deletes a Cognito identity pool.
func (s *CognitoIdentityService) DeleteIdentityPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteIdentityPoolCore(store, getIdentityPoolID(req)); err != nil {
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

	_, maxResultsProvided := req.Parameters["MaxResults"]
	items, nextToken, err := s.listIdentityPoolsShortCore(store, ListIdentityPoolsInput{
		MaxResults:         request.GetIntParam(req.Parameters, "MaxResults"),
		MaxResultsProvided: maxResultsProvided,
		NextToken:          request.GetStringParam(req.Parameters, "NextToken"),
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
	allowUnauthRaw, allowUnauthProvided := req.Parameters["AllowUnauthenticatedIdentities"]
	allowClassicRaw, allowClassicProvided := req.Parameters["AllowClassicFlow"]

	input := UpdateIdentityPoolInput{
		IdentityPoolID:            getIdentityPoolID(req),
		PoolName:                  req.GetParam("IdentityPoolName"),
		AllowUnauthProvided:       allowUnauthProvided,
		AllowUnauthRaw:            allowUnauthRaw,
		AllowClassicProvided:      allowClassicProvided,
		AllowClassicRaw:           allowClassicRaw,
		DeveloperProviderName:     req.GetParam("DeveloperProviderName"),
		ProvidersRaw:              req.Parameters["CognitoIdentityProviders"],
		SupportedLoginProviders:   parseMapParam(req, "SupportedLoginProviders"),
		OpenIdConnectProviderARNs: getStringSliceParam(req, "OpenIdConnectProviderARNs"),
		SamlProviderARNs:          getStringSliceParam(req, "SamlProviderARNs"),
	}
	if _, ok := req.Parameters["IdentityPoolTags"]; ok {
		input.TagsProvided = true
		input.Tags = tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "IdentityPoolTags"))
	}

	result, err := s.updateIdentityPoolCore(reqCtx, input)
	if err != nil {
		return nil, err
	}

	return poolOutToHTTP(result), nil
}

// GetIdentityPoolRoles returns the roles for a Cognito identity pool.
func (s *CognitoIdentityService) GetIdentityPoolRoles(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.getIdentityPoolRolesCore(reqCtx, getIdentityPoolID(req))
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"IdentityPoolId": result.IdentityPoolID,
	}

	roles := map[string]interface{}{}
	if result.AuthRole != "" {
		roles["authenticated"] = result.AuthRole
	}
	if result.UnauthRole != "" {
		roles["unauthenticated"] = result.UnauthRole
	}
	if len(roles) > 0 {
		resp["Roles"] = roles
	}
	if len(result.RoleMappings) > 0 {
		resp["RoleMappings"] = formatRoleMappings(result.RoleMappings)
	}

	return resp, nil
}

// SetIdentityPoolRoles sets the roles for a Cognito identity pool.
func (s *CognitoIdentityService) SetIdentityPoolRoles(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := getIdentityPoolID(req)
	rolesRaw, rolesProvided := req.Parameters["Roles"]

	if err := s.setIdentityPoolRolesCore(reqCtx, SetIdentityPoolRolesInput{
		IdentityPoolID:  poolID,
		RolesProvided:   rolesProvided,
		RolesRaw:        rolesRaw,
		RoleMappingsRaw: req.Parameters["RoleMappings"],
	}); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"IdentityPoolId": poolID,
	}, nil
}

func getIdentityPoolID(req *request.ParsedRequest) string {
	return req.GetParam("IdentityPoolId")
}
