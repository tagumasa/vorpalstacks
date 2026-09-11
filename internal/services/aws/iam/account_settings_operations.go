package iam

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// EnableOutboundWebIdentityFederation enables outbound web identity federation for the account.
func (s *IAMService) EnableOutboundWebIdentityFederation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	issuerIdentifier, err := s.enableOutboundWebIdentityFederationCore(store)
	if err != nil {
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
	if err := s.disableOutboundWebIdentityFederationCore(store); err != nil {
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
	settings, err := s.getOutboundWebIdentityFederationCore(store)
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.setSecurityTokenServicePreferencesCore(store, request.GetStringParam(req.Parameters, "GlobalEndpointTokenVersion")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// GetAccountProperties retrieves the account-level properties for the
// account (Namespace/PropertyName key-value pairs).
func (s *IAMService) GetAccountProperties(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	properties, err := s.getAccountPropertiesCore(store)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Properties": properties,
	}, nil
}

// PutAccountProperties sets account-level properties for the account. The
// wire format is Properties.entry.N.key / Properties.entry.N.value pairs.
func (s *IAMService) PutAccountProperties(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	properties := map[string]string{}
	for i := 1; ; i++ {
		key := request.GetStringParam(req.Parameters, fmt.Sprintf("Properties.entry.%d.key", i))
		if key == "" {
			break
		}
		properties[key] = request.GetStringParam(req.Parameters, fmt.Sprintf("Properties.entry.%d.value", i))
	}

	if err := s.putAccountPropertiesCore(store, properties); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}
