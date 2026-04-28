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
	}
	if err := store.SetUserPoolDomain(domain, domainEntry); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"CloudFrontDomain": cfDomain,
	}, nil
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
		"DomainDescription": map[string]interface{}{
			"Domain":           domainEntry.Domain,
			"UserPoolId":       domainEntry.UserPoolID,
			"CloudFrontDomain": domainEntry.CloudFrontDomain,
			"CreatedDate":      domainEntry.CreatedDate.Unix(),
		},
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

	cognitoSuffix := config.GetString("endpoints.cognito_suffix")
	cfDomain := fmt.Sprintf("%s.auth.%s", domain, strings.Replace(cognitoSuffix, "{region}", reqCtx.GetRegion(), 1))
	domainEntry := &cognitostore.UserPoolDomain{
		Domain:           domain,
		UserPoolID:       userPoolID,
		CloudFrontDomain: cfDomain,
		CreatedDate:      time.Now().UTC(),
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
				rs.Scopes = append(rs.Scopes, cognitostore.ResourceServerScope{
					ScopeName:        fmt.Sprintf("%v", m["ScopeName"]),
					ScopeDescription: fmt.Sprintf("%v", m["ScopeDescription"]),
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
				newScopes = append(newScopes, cognitostore.ResourceServerScope{
					ScopeName:        fmt.Sprintf("%v", m["ScopeName"]),
					ScopeDescription: fmt.Sprintf("%v", m["ScopeDescription"]),
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

	servers, err := store.ListResourceServers(userPoolID)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(servers))
	for _, rs := range servers {
		items = append(items, formatResourceServer(rs))
	}

	return map[string]interface{}{
		"ResourceServers": items,
	}, nil
}

// CreateIdentityProvider adds a new identity provider to a user pool.
func (s *CognitoService) CreateIdentityProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	providerName := req.GetParam("ProviderName")
	providerType := req.GetParam("ProviderType")
	if userPoolID == "" || providerName == "" || providerType == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	var providerDetails map[string]string
	if pd, ok := req.Parameters["ProviderDetails"].(map[string]interface{}); ok {
		providerDetails = make(map[string]string)
		for k, v := range pd {
			if vs, ok := v.(string); ok {
				providerDetails[k] = vs
			}
		}
	}

	ip := &cognitostore.IdentityProvider{
		UserPoolID:      userPoolID,
		ProviderName:    providerName,
		ProviderType:    providerType,
		ProviderDetails: providerDetails,
	}

	if err := store.CreateIdentityProvider(ip); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"IdentityProvider": formatIdentityProvider(ip),
	}, nil
}

// DescribeIdentityProvider returns details of a specified identity provider in a user pool.
func (s *CognitoService) DescribeIdentityProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	providerName := req.GetParam("ProviderName")
	if userPoolID == "" || providerName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ip, err := store.GetIdentityProvider(userPoolID, providerName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{
		"IdentityProvider": formatIdentityProvider(ip),
	}, nil
}

// UpdateIdentityProvider updates the configuration of a specified identity provider in a user pool.
func (s *CognitoService) UpdateIdentityProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	providerName := req.GetParam("ProviderName")
	if userPoolID == "" || providerName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ip, err := store.GetIdentityProvider(userPoolID, providerName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if providerType := req.GetParam("ProviderType"); providerType != "" {
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

	if err := store.UpdateIdentityProvider(ip); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"IdentityProvider": formatIdentityProvider(ip),
	}, nil
}

// DeleteIdentityProvider deletes a specified identity provider from a user pool.
func (s *CognitoService) DeleteIdentityProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	providerName := req.GetParam("ProviderName")
	if userPoolID == "" || providerName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteIdentityProvider(userPoolID, providerName); err != nil {
		return nil, ErrResourceNotFound
	}

	return response.EmptyResponse(), nil
}

// ListIdentityProviders lists all identity providers in a user pool.
func (s *CognitoService) ListIdentityProviders(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	providers, err := store.ListIdentityProviders(userPoolID)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(providers))
	for _, ip := range providers {
		items = append(items, formatIdentityProvider(ip))
	}

	return map[string]interface{}{
		"Providers": items,
	}, nil
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

	csvHeader := []string{
		"cognito:username", "name", "given_name", "family_name", "middle_name",
		"nickname", "preferred_username", "profile", "picture", "website",
		"email", "email_verified", "gender", "birthdate", "zoneinfo",
		"locale", "phone_number", "phone_number_verified", "address", "updated_at",
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

	return map[string]interface{}{
		"RiskConfiguration": map[string]interface{}{
			"UserPoolId": userPoolID,
			"CompromisedCredentialsRiskConfiguration": map[string]interface{}{
				"EventFilter": []string{},
				"Actions": map[string]interface{}{
					"EventAction": "NO_ACTION",
				},
			},
			"AccountTakeoverRiskConfiguration": map[string]interface{}{
				"NotifyConfiguration": map[string]interface{}{
					"From": "NO_ACTION",
				},
				"Actions": map[string]interface{}{
					"HighAction": map[string]interface{}{
						"EventAction": "NO_ACTION",
						"Notify":      false,
					},
					"MediumAction": map[string]interface{}{
						"EventAction": "NO_ACTION",
						"Notify":      false,
					},
					"LowAction": map[string]interface{}{
						"EventAction": "NO_ACTION",
						"Notify":      false,
					},
				},
			},
			"RiskExceptionConfiguration": map[string]interface{}{
				"BlockedIPRangeList":    []string{},
				"NotBlockedIPRangeList": []string{},
			},
			"LastModifiedDate": time.Now().UTC().Unix(),
		},
	}, nil
}

// formatResourceServer converts a ResourceServer store model to the API response map.
func formatResourceServer(rs *cognitostore.ResourceServer) map[string]interface{} {
	result := map[string]interface{}{
		"UserPoolId": rs.UserPoolID,
		"Identifier": rs.Identifier,
		"Name":       rs.Name,
		"Scopes":     []interface{}{},
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
		"UserPoolId":   ip.UserPoolID,
		"ProviderName": ip.ProviderName,
		"ProviderType": ip.ProviderType,
	}
	if ip.ProviderDetails != nil {
		result["ProviderDetails"] = ip.ProviderDetails
	}
	return result
}
