package iam

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/config"
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

	// Enabling an already enabled feature fails with FeatureEnabled.
	if settings.OutboundWebIdentityFederationEnabled {
		return nil, ErrFeatureEnabled
	}

	// Derive the OIDC issuer identifier from the configured base URL
	// rather than hardcoding amazonaws.com.  For edge/on-prem
	// deployments the base URL reflects the actual deployment domain.
	issuerIdentifier := deriveOIDCIssuer(s.accountID)

	settings.OutboundWebIdentityFederationEnabled = true
	settings.IssuerIdentifier = issuerIdentifier

	if err := store.AccountSettings().Put(settings); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"IssuerIdentifier": issuerIdentifier,
	}, nil
}

// deriveOIDCIssuer builds the OIDC issuer URL from the configured
// endpoints.base_url, falling back to the amazonaws.com format when
// the base URL cannot be parsed.
//
// The host segment preserves non-default ports so that edge deployments
// exposed on a non-443 port produce a usable issuer URL (e.g.
// https://oidc.<account>.edge.example:8443). Default ports (80 for http,
// 443 for https) are stripped because including them in the issuer would
// diverge from the canonical URL form expected by JWT consumers.
func deriveOIDCIssuer(accountID string) string {
	baseURL := config.BaseURL()
	if baseURL == "" {
		return fmt.Sprintf("https://oidc.%s.amazonaws.com", accountID)
	}
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Hostname() == "" {
		return fmt.Sprintf("https://oidc.%s.amazonaws.com", accountID)
	}
	host := parsed.Hostname()
	if port := parsed.Port(); port != "" && port != "80" && port != "443" {
		host = host + ":" + port
	}
	return fmt.Sprintf("https://oidc.%s.%s", accountID, host)
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

	// Disabling a feature that is not enabled fails with FeatureDisabled.
	if !settings.OutboundWebIdentityFederationEnabled {
		return nil, ErrFeatureDisabled
	}

	settings.OutboundWebIdentityFederationEnabled = false
	settings.IssuerIdentifier = ""

	if err := store.AccountSettings().Put(settings); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
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

	// Reading the configuration of a disabled feature fails with
	// FeatureDisabled.
	if !settings.OutboundWebIdentityFederationEnabled {
		return nil, ErrFeatureDisabled
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
		return nil, errors.NewAWSError("InvalidInput", tokenVersion+" is not a valid value for GlobalEndpointTokenVersion. Must be v1Token or v2Token.", http.StatusBadRequest)
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

	return response.EmptyResponse(), nil
}
