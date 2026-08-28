package iam

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateOpenIDConnectProvider creates an IAM entity to describe an identity provider (IdP) that supports OpenID Connect (OIDC).
func (s *IAMService) CreateOpenIDConnectProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &CreateOpenIDConnectProviderInput{
		Url:            request.GetStringParam(req.Parameters, "Url"),
		ThumbprintList: request.GetStringList(req.Parameters, "ThumbprintList"),
		ClientIDList:   request.GetStringList(req.Parameters, "ClientIDList"),
		Tags:           tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"),
	}
	providerArn, err := s.createOpenIDConnectProviderCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"OpenIDConnectProviderArn": providerArn,
	}, nil
}

// GetOpenIDConnectProvider retrieves information about an OpenID Connect (OIDC) provider.
func (s *IAMService) GetOpenIDConnectProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	provider, err := s.getOpenIDConnectProviderCore(store, request.GetStringParam(req.Parameters, "OpenIDConnectProviderArn"))
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"OpenIDConnectProviderArn": provider.Arn,
		"Url":                      provider.URL,
		"CreateDate":               provider.CreateDate.Format(timeutils.ISO8601SimpleFormat),
	}

	if provider.ThumbprintList != nil {
		resp["ThumbprintList"] = provider.ThumbprintList
	}
	if provider.ClientIDList != nil {
		resp["ClientIDList"] = provider.ClientIDList
	}
	if provider.LastModifiedDate != nil {
		resp["LastModifiedDate"] = provider.LastModifiedDate.Format(timeutils.ISO8601SimpleFormat)
	}
	if provider.Tags != nil {
		resp["Tags"] = tags.ToResponse(provider.Tags)
	}

	return resp, nil
}

// ListOpenIDConnectProviders lists the OpenID Connect (OIDC) providers in the account.
func (s *IAMService) ListOpenIDConnectProviders(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	providerList, err := s.listOpenIDConnectProvidersCore(store)
	if err != nil {
		return nil, err
	}

	list := make([]interface{}, len(providerList))
	for i, provider := range providerList {
		// Smithy OpenIDConnectProviderListEntry contains only Arn.
		// CreateDate is intentionally omitted for spec compliance.
		list[i] = map[string]interface{}{
			"Arn": provider.Arn,
		}
	}

	return map[string]interface{}{
		"OpenIDConnectProviderList": list,
	}, nil
}

// UpdateOpenIDConnectProviderThumbprint replaces the existing list of thumbprints with a new list for the specified OpenID Connect (OIDC) provider.
func (s *IAMService) UpdateOpenIDConnectProviderThumbprint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &UpdateOpenIDConnectProviderThumbprintInput{
		ProviderArn:    request.GetStringParam(req.Parameters, "OpenIDConnectProviderArn"),
		ThumbprintList: request.GetStringList(req.Parameters, "ThumbprintList"),
	}
	if err := s.updateOpenIDConnectProviderThumbprintCore(store, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// AddClientIDToOpenIDConnectProvider adds a new client ID to the list of client IDs for the specified OpenID Connect (OIDC) provider.
func (s *IAMService) AddClientIDToOpenIDConnectProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &OpenIDConnectProviderClientIDInput{
		ProviderArn: request.GetStringParam(req.Parameters, "OpenIDConnectProviderArn"),
		ClientID:    request.GetStringParam(req.Parameters, "ClientID"),
	}
	if err := s.addClientIDToOpenIDConnectProviderCore(store, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// RemoveClientIDFromOpenIDConnectProvider removes the specified client ID from the list of client IDs for the specified OpenID Connect (OIDC) provider.
func (s *IAMService) RemoveClientIDFromOpenIDConnectProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &OpenIDConnectProviderClientIDInput{
		ProviderArn: request.GetStringParam(req.Parameters, "OpenIDConnectProviderArn"),
		ClientID:    request.GetStringParam(req.Parameters, "ClientID"),
	}
	if err := s.removeClientIDFromOpenIDConnectProviderCore(store, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteOpenIDConnectProvider deletes an OpenID Connect (OIDC) identity provider.
func (s *IAMService) DeleteOpenIDConnectProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteOpenIDConnectProviderCore(store, request.GetStringParam(req.Parameters, "OpenIDConnectProviderArn")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

var oidcProviderTagOps = tagOps[*iamstore.OpenIDConnectProvider]{
	paramName:  "OpenIDConnectProviderArn",
	emptyErr:   NewValidationError("OpenIDConnectProviderArn"),
	notFoundFn: func(n string) error { return NewNoSuchEntityError("OpenID Connect provider", n) },
	getFn: func(s *iamstore.IAMStore, n string) (*iamstore.OpenIDConnectProvider, error) {
		return s.OpenIDConnectProviders().Get(n)
	},
	putFn: func(s *iamstore.IAMStore, r *iamstore.OpenIDConnectProvider) error {
		return s.OpenIDConnectProviders().Put(r)
	},
	tagsFn: func(r *iamstore.OpenIDConnectProvider) *[]tags.Tag { return &r.Tags },
}

// TagOpenIDConnectProvider adds tags to an OpenID Connect (OIDC) provider.
func (s *IAMService) TagOpenIDConnectProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagResource(ctx, s, reqCtx, req, oidcProviderTagOps)
}

// UntagOpenIDConnectProvider removes tags from an OpenID Connect (OIDC) provider.
func (s *IAMService) UntagOpenIDConnectProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return untagResource(ctx, s, reqCtx, req, oidcProviderTagOps)
}

// ListOpenIDConnectProviderTags lists the tags attached to an OpenID Connect (OIDC) provider.
func (s *IAMService) ListOpenIDConnectProviderTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return listResourceTags(ctx, s, reqCtx, req, oidcProviderTagOps)
}
