package iam

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
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

// TagOpenIDConnectProvider adds tags to an OpenID Connect (OIDC) provider.
func (s *IAMService) TagOpenIDConnectProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &TagResourceInput{
		ResourceName: request.GetStringParam(req.Parameters, "OpenIDConnectProviderArn"),
		Tags:         tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"),
	}
	if err := tagResourceCore(store, oidcProviderTagOps, input); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// UntagOpenIDConnectProvider removes tags from an OpenID Connect (OIDC) provider.
func (s *IAMService) UntagOpenIDConnectProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &UntagResourceInput{
		ResourceName: request.GetStringParam(req.Parameters, "OpenIDConnectProviderArn"),
		TagKeys:      tags.ParseTagKeysWithQueryFallback(req.Parameters, "TagKeys"),
	}
	if err := untagResourceCore(store, oidcProviderTagOps, input); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListOpenIDConnectProviderTags lists the tags attached to an OpenID Connect (OIDC) provider.
func (s *IAMService) ListOpenIDConnectProviderTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &ListResourceTagsInput{
		ResourceName: request.GetStringParam(req.Parameters, "OpenIDConnectProviderArn"),
		Marker:       request.GetStringParam(req.Parameters, "Marker"),
		MaxItems:     pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems),
	}
	result, err := listResourceTagsCore(store, oidcProviderTagOps, input)
	if err != nil {
		return nil, err
	}
	resp := map[string]interface{}{
		"Tags":        tags.ToResponse(result.Tags),
		"IsTruncated": result.IsTruncated,
	}
	if result.Marker != "" {
		resp["Marker"] = result.Marker
	}
	return resp, nil
}
