package cognitoidentityprovider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/config"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	"vorpalstacks/internal/store/aws/common"
)

// CreateUserPoolDomain creates a new domain for a user pool.
func (s *CognitoService) CreateUserPoolDomain(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domain := req.GetParam("Domain")
	userPoolID := getUserPoolID(req)
	if domain == "" || userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	cognitoSuffix := config.GetString("endpoints.cognito_suffix")
	cfDomain := fmt.Sprintf("%s.auth.%s", domain, strings.Replace(cognitoSuffix, "{region}", reqCtx.GetRegion(), 1))
	domainEntry := &cognitostore.UserPoolDomain{
		Domain:           domain,
		UserPoolID:       userPoolID,
		CloudFrontDomain: cfDomain,
		CreatedDate:      time.Now().UTC(),
		Status:           "ACTIVE",
	}
	if v, ok := req.Parameters["ManagedLoginVersion"]; ok {
		switch n := v.(type) {
		case int:
			mlv := n
			domainEntry.ManagedLoginVersion = &mlv
		case float64:
			mlv := int(n)
			domainEntry.ManagedLoginVersion = &mlv
		}
	}
	if cdc, ok := req.Parameters["CustomDomainConfig"].(map[string]interface{}); ok {
		cfg := &cognitostore.CustomDomainConfig{}
		if v, ok := cdc["CertificateArn"].(string); ok {
			cfg.CertificateArn = v
		}
		if v, ok := cdc["SecurityPolicy"].(string); ok {
			if !validateSecurityPolicy(v) {
				return nil, ErrInvalidParameter
			}
			cfg.SecurityPolicy = v
		}
		domainEntry.CustomDomainConfig = cfg
	}
	if rt, ok := req.Parameters["Routing"].(map[string]interface{}); ok {
		routing := &cognitostore.Routing{}
		if fo, ok := rt["Failover"].(map[string]interface{}); ok {
			failover := &cognitostore.FailoverType{}
			if v, ok := fo["SecondaryRegion"].(string); ok {
				failover.SecondaryRegion = v
			}
			if v, ok := fo["PrimaryRoute53HealthCheckId"].(string); ok {
				failover.PrimaryRoute53HealthCheckId = v
			}
			routing.Failover = failover
		}
		domainEntry.Routing = routing
	}
	if err := store.SetUserPoolDomain(domain, domainEntry); err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"CloudFrontDomain": cfDomain,
	}
	if domainEntry.ManagedLoginVersion != nil {
		resp["ManagedLoginVersion"] = *domainEntry.ManagedLoginVersion
	}
	if domainEntry.Routing != nil && domainEntry.Routing.Failover != nil {
		resp["Routing"] = routingToMap(domainEntry.Routing)
	}

	return resp, nil
}

// routingToMap serialises a Routing store struct into the Smithy Routing JSON shape.
func routingToMap(r *cognitostore.Routing) map[string]interface{} {
	m := map[string]interface{}{}
	if r.Failover != nil {
		m["Failover"] = map[string]interface{}{
			"SecondaryRegion":             r.Failover.SecondaryRegion,
			"PrimaryRoute53HealthCheckId": r.Failover.PrimaryRoute53HealthCheckId,
		}
	}
	return m
}

// buildDomainDescription constructs a DomainDescriptionType map from a store
// UserPoolDomain, including all Smithy-defined members.
func buildDomainDescription(d *cognitostore.UserPoolDomain) map[string]interface{} {
	desc := map[string]interface{}{
		"Domain":                 d.Domain,
		"UserPoolId":             d.UserPoolID,
		"CloudFrontDistribution": d.CloudFrontDomain,
		"Status":                 d.Status,
	}
	if d.ManagedLoginVersion != nil {
		desc["ManagedLoginVersion"] = *d.ManagedLoginVersion
	}
	if d.CustomDomainConfig != nil {
		cfg := map[string]interface{}{
			"CertificateArn": d.CustomDomainConfig.CertificateArn,
		}
		if d.CustomDomainConfig.SecurityPolicy != "" {
			cfg["SecurityPolicy"] = d.CustomDomainConfig.SecurityPolicy
		}
		desc["CustomDomainConfig"] = cfg
	}
	if d.Routing != nil && d.Routing.Failover != nil {
		desc["Routing"] = routingToMap(d.Routing)
	}
	return desc
}

// DescribeUserPoolDomain returns information about a user pool domain.
func (s *CognitoService) DescribeUserPoolDomain(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domain := req.GetParam("Domain")
	if domain == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	domainEntry, err := store.GetUserPoolDomain(domain)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{
		"DomainDescription": buildDomainDescription(domainEntry),
	}, nil
}

// DeleteUserPoolDomain deletes a domain from a user pool.
func (s *CognitoService) DeleteUserPoolDomain(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domain := req.GetParam("Domain")
	if domain == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteUserPoolDomain(domain); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UpdateUserPoolDomain updates the configuration for a user pool domain.
func (s *CognitoService) UpdateUserPoolDomain(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domain := req.GetParam("Domain")
	userPoolID := getUserPoolID(req)
	if domain == "" || userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	existing, err := store.GetUserPoolDomain(domain)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	cognitoSuffix := config.GetString("endpoints.cognito_suffix")
	cfDomain := fmt.Sprintf("%s.auth.%s", domain, strings.Replace(cognitoSuffix, "{region}", reqCtx.GetRegion(), 1))
	domainEntry := &cognitostore.UserPoolDomain{
		Domain:           domain,
		UserPoolID:       userPoolID,
		CloudFrontDomain: cfDomain,
		CreatedDate:      existing.CreatedDate,
	}
	if err := store.SetUserPoolDomain(domain, domainEntry); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"CloudFrontDomain": cfDomain,
	}, nil
}

// CreateResourceServer creates a new resource server for a user pool.
func (s *CognitoService) CreateResourceServer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	identifier := req.GetParam("Identifier")
	name := req.GetParam("Name")
	if userPoolID == "" || identifier == "" || name == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	rs := &cognitostore.ResourceServer{
		UserPoolID: userPoolID,
		Identifier: identifier,
		Name:       name,
		Scopes:     []cognitostore.ResourceServerScope{},
	}

	if scopes, ok := req.Parameters["Scopes"].([]interface{}); ok {
		for _, sc := range scopes {
			if m, ok := sc.(map[string]interface{}); ok {
				scopeName, _ := m["ScopeName"].(string)
				scopeDesc, _ := m["ScopeDescription"].(string)
				if scopeName == "" {
					return nil, ErrInvalidParameter
				}
				rs.Scopes = append(rs.Scopes, cognitostore.ResourceServerScope{
					ScopeName:        scopeName,
					ScopeDescription: scopeDesc,
				})
			}
		}
	}

	if err := store.CreateResourceServer(rs); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ResourceServer": formatResourceServer(rs),
	}, nil
}

// DescribeResourceServer returns details of a specified resource server in a user pool.
func (s *CognitoService) DescribeResourceServer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	identifier := req.GetParam("Identifier")
	if userPoolID == "" || identifier == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	rs, err := store.GetResourceServer(userPoolID, identifier)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{
		"ResourceServer": formatResourceServer(rs),
	}, nil
}

// UpdateResourceServer updates the name and scopes of a specified resource server in a user pool.
func (s *CognitoService) UpdateResourceServer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	identifier := req.GetParam("Identifier")
	if userPoolID == "" || identifier == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	rs, err := store.GetResourceServer(userPoolID, identifier)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if name := req.GetParam("Name"); name != "" {
		rs.Name = name
	}

	if scopes, ok := req.Parameters["Scopes"].([]interface{}); ok {
		var newScopes []cognitostore.ResourceServerScope
		for _, sc := range scopes {
			if m, ok := sc.(map[string]interface{}); ok {
				scopeName, _ := m["ScopeName"].(string)
				scopeDesc, _ := m["ScopeDescription"].(string)
				if scopeName == "" {
					return nil, ErrInvalidParameter
				}
				newScopes = append(newScopes, cognitostore.ResourceServerScope{
					ScopeName:        scopeName,
					ScopeDescription: scopeDesc,
				})
			}
		}
		rs.Scopes = newScopes
	}

	if err := store.UpdateResourceServer(rs); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"ResourceServer": formatResourceServer(rs),
	}, nil
}

// DeleteResourceServer deletes a specified resource server from a user pool.
func (s *CognitoService) DeleteResourceServer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	identifier := req.GetParam("Identifier")
	if userPoolID == "" || identifier == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteResourceServer(userPoolID, identifier); err != nil {
		return nil, ErrResourceNotFound
	}

	return response.EmptyResponse(), nil
}

// ListResourceServers lists all resource servers in a user pool.
func (s *CognitoService) ListResourceServers(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	maxResults := request.GetIntParam(req.Parameters, "MaxResults")
	if maxResults <= 0 || maxResults > 50 {
		maxResults = 50
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")

	opts := common.ListOptions{
		MaxItems: maxResults,
		Marker:   nextToken,
	}

	result, err := store.ListResourceServersPaginated(userPoolID, opts)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Items))
	for _, rs := range result.Items {
		items = append(items, formatResourceServer(rs))
	}

	resp := map[string]interface{}{
		"ResourceServers": items,
	}
	if result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}

	return resp, nil
}

// CreateIdentityProvider adds a new identity provider to a user pool.
func (s *CognitoService) CreateIdentityProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	providerName := req.GetParam("ProviderName")
	providerType := req.GetParam("ProviderType")

	var providerDetails map[string]string
	if pd, ok := req.Parameters["ProviderDetails"].(map[string]interface{}); ok {
		providerDetails = make(map[string]string)
		for k, v := range pd {
			if vs, ok := v.(string); ok {
				providerDetails[k] = vs
			}
		}
	}

	var attributeMapping map[string]string
	if am, ok := req.Parameters["AttributeMapping"].(map[string]interface{}); ok {
		attributeMapping = make(map[string]string)
		for k, v := range am {
			if vs, ok := v.(string); ok {
				attributeMapping[k] = vs
			}
		}
	}

	ip, err := s.createIdentityProviderFromInputCore(reqCtx.GetRegion(), CreateIdentityProviderInput{
		UserPoolID:       userPoolID,
		ProviderName:     providerName,
		ProviderType:     providerType,
		ProviderDetails:  providerDetails,
		AttributeMapping: attributeMapping,
		IdpIdentifiers:   getStringSliceParam(req, "IdpIdentifiers"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"IdentityProvider": formatIdentityProvider(ip),
	}, nil
}

// DescribeIdentityProvider returns details of a specified identity provider in a user pool.
func (s *CognitoService) DescribeIdentityProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	ip, err := s.describeIdentityProviderCore(reqCtx.GetRegion(), getUserPoolID(req), req.GetParam("ProviderName"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"IdentityProvider": formatIdentityProvider(ip)}, nil
}

// UpdateIdentityProvider updates the configuration of a specified identity provider in a user pool.
func (s *CognitoService) UpdateIdentityProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	providerName := req.GetParam("ProviderName")
	if userPoolID == "" || providerName == "" {
		return nil, ErrInvalidParameter
	}

	ip, err := s.describeIdentityProviderCore(reqCtx.GetRegion(), userPoolID, providerName)
	if err != nil {
		return nil, err
	}

	if providerType := req.GetParam("ProviderType"); providerType != "" {
		if !validateProviderType(providerType) {
			return nil, ErrInvalidParameter
		}
		ip.ProviderType = providerType
	}

	if pd, ok := req.Parameters["ProviderDetails"].(map[string]interface{}); ok {
		providerDetails := make(map[string]string)
		for k, v := range pd {
			if vs, ok := v.(string); ok {
				providerDetails[k] = vs
			}
		}
		ip.ProviderDetails = providerDetails
	}

	if am, ok := req.Parameters["AttributeMapping"].(map[string]interface{}); ok {
		ip.AttributeMapping = make(map[string]string)
		for k, v := range am {
			if vs, ok := v.(string); ok {
				ip.AttributeMapping[k] = vs
			}
		}
	}
	if ids := getStringSliceParam(req, "IdpIdentifiers"); len(ids) > 0 {
		ip.IdpIdentifiers = ids
	}

	if err := s.updateIdentityProviderCore(reqCtx.GetRegion(), ip); err != nil {
		return nil, err
	}

	return map[string]interface{}{"IdentityProvider": formatIdentityProvider(ip)}, nil
}

// DeleteIdentityProvider deletes a specified identity provider from a user pool.
func (s *CognitoService) DeleteIdentityProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.deleteIdentityProviderCore(reqCtx.GetRegion(), getUserPoolID(req), req.GetParam("ProviderName")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListIdentityProviders lists all identity providers in a user pool.
func (s *CognitoService) ListIdentityProviders(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.listIdentityProvidersCore(reqCtx.GetRegion(), ListIdentityProvidersInput{
		UserPoolID: getUserPoolID(req),
		MaxResults: request.GetIntParam(req.Parameters, "MaxResults"),
		NextToken:  request.GetStringParam(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Providers))
	for _, ip := range result.Providers {
		items = append(items, formatIdentityProvider(ip))
	}

	resp := map[string]interface{}{
		"Providers": items,
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}

	return resp, nil
}

// GetCSVHeader returns the CSV headers for importing users into a user pool.
func (s *CognitoService) GetCSVHeader(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	pool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	csvHeader := []string{
		"cognito:username", "name", "given_name", "family_name", "middle_name",
		"nickname", "preferred_username", "profile", "picture", "website",
		"email", "email_verified", "gender", "birthdate", "zoneinfo",
		"locale", "phone_number", "phone_number_verified", "address", "updated_at",
	}
	for _, sa := range pool.SchemaAttributes {
		if sa.Name != "" {
			csvHeader = append(csvHeader, sa.Name)
		}
	}

	return map[string]interface{}{
		"CSVHeader": csvHeader,
	}, nil
}

// DescribeRiskConfiguration describes the risk configuration for a user pool.
func (s *CognitoService) DescribeRiskConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	clientID := req.GetParam("ClientId")

	cfg, err := store.GetRiskConfiguration(userPoolID, clientID)
	if err == nil && cfg != nil {
		return map[string]interface{}{
			"RiskConfiguration": formatRiskConfiguration(cfg),
		}, nil
	}

	if clientID != "" {
		cfg, err = store.GetRiskConfiguration(userPoolID, "")
		if err == nil && cfg != nil {
			return map[string]interface{}{
				"RiskConfiguration": formatRiskConfiguration(cfg),
			}, nil
		}
	}

	return map[string]interface{}{
		"RiskConfiguration": formatRiskConfiguration(&cognitostore.RiskConfiguration{UserPoolID: userPoolID}),
	}, nil
}

// formatResourceServer converts a ResourceServer store model to the API response map.
func formatResourceServer(rs *cognitostore.ResourceServer) map[string]interface{} {
	result := map[string]interface{}{
		"UserPoolId":       rs.UserPoolID,
		"Identifier":       rs.Identifier,
		"Name":             rs.Name,
		"Scopes":           []interface{}{},
		"CreationDate":     rs.CreationDate.Unix(),
		"LastModifiedDate": rs.LastModifiedDate.Unix(),
	}
	if len(rs.Scopes) > 0 {
		scopes := make([]interface{}, 0, len(rs.Scopes))
		for _, sc := range rs.Scopes {
			scopes = append(scopes, map[string]interface{}{
				"ScopeName":        sc.ScopeName,
				"ScopeDescription": sc.ScopeDescription,
			})
		}
		result["Scopes"] = scopes
	}
	return result
}

// formatIdentityProvider converts an IdentityProvider store model to the API response map.
func formatIdentityProvider(ip *cognitostore.IdentityProvider) map[string]interface{} {
	result := map[string]interface{}{
		"UserPoolId":       ip.UserPoolID,
		"ProviderName":     ip.ProviderName,
		"ProviderType":     ip.ProviderType,
		"CreationDate":     ip.CreationDate.Unix(),
		"LastModifiedDate": ip.LastModifiedDate.Unix(),
	}
	if ip.ProviderDetails != nil {
		result["ProviderDetails"] = ip.ProviderDetails
	}
	if ip.AttributeMapping != nil {
		result["AttributeMapping"] = ip.AttributeMapping
	}
	if len(ip.IdpIdentifiers) > 0 {
		result["IdpIdentifiers"] = ip.IdpIdentifiers
	}
	return result
}
