package iam

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/request"
)

// EnableOutboundWebIdentityFederation enables outbound web identity federation for the account.
func (s *IAMService) EnableOutboundWebIdentityFederation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	settings, err := store.AccountSettings().Get()
	if err != nil {
		return nil, err
	}

	issuerIdentifier := fmt.Sprintf("https://oidc.%s.amazonaws.com", s.accountID)
	settings.OutboundWebIdentityFederationEnabled = true
	settings.IssuerIdentifier = issuerIdentifier

	if err := store.AccountSettings().Put(settings); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"IssuerIdentifier": issuerIdentifier,
	}, nil
}

// DisableOutboundWebIdentityFederation disables outbound web identity federation for the account.
func (s *IAMService) DisableOutboundWebIdentityFederation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	settings, err := store.AccountSettings().Get()
	if err != nil {
		return nil, err
	}

	settings.OutboundWebIdentityFederationEnabled = false
	settings.IssuerIdentifier = ""

	if err := store.AccountSettings().Put(settings); err != nil {
		return nil, err
	}

	return nil, nil
}

// GetOutboundWebIdentityFederationInfo retrieves the outbound web identity federation configuration for the account.
func (s *IAMService) GetOutboundWebIdentityFederationInfo(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	settings, err := store.AccountSettings().Get()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"IssuerIdentifier":  settings.IssuerIdentifier,
		"JwtVendingEnabled": settings.OutboundWebIdentityFederationEnabled,
	}, nil
}

// SetSecurityTokenServicePreferences sets the global endpoint token version preference for the account.
func (s *IAMService) SetSecurityTokenServicePreferences(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tokenVersion := request.GetStringParam(req.Parameters, "GlobalEndpointTokenVersion")
	if tokenVersion == "" {
		return nil, NewValidationError("GlobalEndpointTokenVersion")
	}

	if tokenVersion != "v1Token" && tokenVersion != "v2Token" {
		return nil, fmt.Errorf("InvalidParameterValue: %s is not a valid value for GlobalEndpointTokenVersion", tokenVersion)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	settings, err := store.AccountSettings().Get()
	if err != nil {
		return nil, err
	}

	settings.GlobalEndpointTokenVersion = tokenVersion

	if err := store.AccountSettings().Put(settings); err != nil {
		return nil, err
	}

	return nil, nil
}
