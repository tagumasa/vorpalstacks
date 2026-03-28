package cognitoidentityprovider

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

func (s *CognitoService) CreateUserPoolDomain(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domain := getParam(req, "Domain")
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

	cfDomain := fmt.Sprintf("%s.auth.%s.amazoncognito.com", domain, s.region)
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

func (s *CognitoService) DescribeUserPoolDomain(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domain := getParam(req, "Domain")
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

func (s *CognitoService) DeleteUserPoolDomain(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domain := getParam(req, "Domain")
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

func (s *CognitoService) CreateResourceServer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	identifier := getParam(req, "Identifier")
	name := getParam(req, "Name")
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

	if err := store.CreateResourceServer(rs); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ResourceServer": map[string]interface{}{
			"UserPoolId": rs.UserPoolID,
			"Identifier": rs.Identifier,
			"Name":       rs.Name,
			"Scopes":     []interface{}{},
		},
	}, nil
}

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
		scopes := make([]interface{}, 0, len(rs.Scopes))
		for _, sc := range rs.Scopes {
			scopes = append(scopes, map[string]interface{}{
				"ScopeName":        sc.ScopeName,
				"ScopeDescription": sc.ScopeDescription,
			})
		}
		items = append(items, map[string]interface{}{
			"UserPoolId": rs.UserPoolID,
			"Identifier": rs.Identifier,
			"Name":       rs.Name,
			"Scopes":     scopes,
		})
	}

	return map[string]interface{}{
		"ResourceServers": items,
	}, nil
}

func (s *CognitoService) CreateIdentityProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	providerName := getParam(req, "ProviderName")
	providerType := getParam(req, "ProviderType")
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

	result := map[string]interface{}{
		"UserPoolId":   ip.UserPoolID,
		"ProviderName": ip.ProviderName,
		"ProviderType": ip.ProviderType,
	}
	if ip.ProviderDetails != nil {
		result["ProviderDetails"] = ip.ProviderDetails
	}

	return map[string]interface{}{
		"IdentityProvider": result,
	}, nil
}

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
		result := map[string]interface{}{
			"UserPoolId":   ip.UserPoolID,
			"ProviderName": ip.ProviderName,
			"ProviderType": ip.ProviderType,
		}
		if ip.ProviderDetails != nil {
			result["ProviderDetails"] = ip.ProviderDetails
		}
		items = append(items, result)
	}

	return map[string]interface{}{
		"Providers": items,
	}, nil
}

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
