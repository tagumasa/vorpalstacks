package iam

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateOpenIDConnectProvider creates an IAM entity to describe an identity provider (IdP) that supports OpenID Connect (OIDC).
func (s *IAMService) CreateOpenIDConnectProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	url := request.GetStringParam(req.Parameters, "Url")
	if url == "" {
		return nil, NewValidationError("Url")
	}

	thumbprintList := request.GetStringList(req.Parameters, "ThumbprintList")
	clientIdList := request.GetStringList(req.Parameters, "ClientIDList")
	newTags := tags.ParseTagsWithQueryFallback(req.Parameters, "Tags")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	provider, err := store.OpenIDConnectProviders().Create(url, thumbprintList, clientIdList, newTags)
	if err != nil {
		if errors.Is(err, iamstore.ErrOpenIDConnectProviderAlreadyExists) {
			return nil, NewEntityAlreadyExistsError("OpenID Connect Provider " + url)
		}
		return nil, err
	}

	return map[string]interface{}{
		"OpenIDConnectProviderArn": provider.Arn,
	}, nil
}

// GetOpenIDConnectProvider retrieves information about an OpenID Connect (OIDC) provider.
func (s *IAMService) GetOpenIDConnectProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	providerArn := request.GetStringParam(req.Parameters, "OpenIDConnectProviderArn")
	if providerArn == "" {
		return nil, NewValidationError("OpenIDConnectProviderArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	provider, err := store.OpenIDConnectProviders().Get(providerArn)
	if err != nil {
		return nil, NewNoSuchEntityError("OpenID Connect provider", providerArn)
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
	result, err := store.OpenIDConnectProviders().List()
	if err != nil {
		return nil, err
	}

	providerList := make([]interface{}, len(result.OpenIDConnectProviderList))
	for i, provider := range result.OpenIDConnectProviderList {
		providerList[i] = map[string]interface{}{
			"Arn":        provider.Arn,
			"CreateDate": provider.CreateDate.Format(timeutils.ISO8601SimpleFormat),
		}
	}

	return map[string]interface{}{
		"OpenIDConnectProviderList": providerList,
	}, nil
}

// UpdateOpenIDConnectProviderThumbprint replaces the existing list of thumbprints with a new list for the specified OpenID Connect (OIDC) provider.
func (s *IAMService) UpdateOpenIDConnectProviderThumbprint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	providerArn := request.GetStringParam(req.Parameters, "OpenIDConnectProviderArn")
	if providerArn == "" {
		return nil, NewValidationError("OpenIDConnectProviderArn")
	}

	thumbprintList := request.GetStringList(req.Parameters, "ThumbprintList")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.OpenIDConnectProviders().Exists(providerArn) {
		return nil, NewNoSuchEntityError("OpenID Connect provider", providerArn)
	}

	if err := store.OpenIDConnectProviders().Update(providerArn, thumbprintList, nil); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// AddClientIDToOpenIDConnectProvider adds a new client ID to the list of client IDs for the specified OpenID Connect (OIDC) provider.
func (s *IAMService) AddClientIDToOpenIDConnectProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	providerArn := request.GetStringParam(req.Parameters, "OpenIDConnectProviderArn")
	if providerArn == "" {
		return nil, NewValidationError("OpenIDConnectProviderArn")
	}
	clientId := request.GetStringParam(req.Parameters, "ClientID")
	if clientId == "" {
		return nil, NewValidationError("ClientID")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	provider, err := store.OpenIDConnectProviders().Get(providerArn)
	if err != nil {
		return nil, NewNoSuchEntityError("OpenID Connect provider", providerArn)
	}

	for _, existing := range provider.ClientIDList {
		if existing == clientId {
			return nil, NewInvalidInputError("ClientID", "already exists")
		}
	}

	provider.ClientIDList = append(provider.ClientIDList, clientId)
	if err := store.OpenIDConnectProviders().Put(provider); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// RemoveClientIDFromOpenIDConnectProvider removes the specified client ID from the list of client IDs for the specified OpenID Connect (OIDC) provider.
func (s *IAMService) RemoveClientIDFromOpenIDConnectProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	providerArn := request.GetStringParam(req.Parameters, "OpenIDConnectProviderArn")
	if providerArn == "" {
		return nil, NewValidationError("OpenIDConnectProviderArn")
	}
	clientId := request.GetStringParam(req.Parameters, "ClientID")
	if clientId == "" {
		return nil, NewValidationError("ClientID")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	provider, err := store.OpenIDConnectProviders().Get(providerArn)
	if err != nil {
		return nil, NewNoSuchEntityError("OpenID Connect provider", providerArn)
	}

	found := false
	newList := make([]string, 0, len(provider.ClientIDList))
	for _, existing := range provider.ClientIDList {
		if existing == clientId {
			found = true
			continue
		}
		newList = append(newList, existing)
	}

	if !found {
		return nil, NewNoSuchEntityError("ClientID", clientId)
	}

	provider.ClientIDList = newList
	if err := store.OpenIDConnectProviders().Put(provider); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteOpenIDConnectProvider deletes an OpenID Connect (OIDC) identity provider.
func (s *IAMService) DeleteOpenIDConnectProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	providerArn := request.GetStringParam(req.Parameters, "OpenIDConnectProviderArn")
	if providerArn == "" {
		return nil, NewValidationError("OpenIDConnectProviderArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.OpenIDConnectProviders().Exists(providerArn) {
		return nil, NewNoSuchEntityError("OpenID Connect provider", providerArn)
	}
	if err := store.OpenIDConnectProviders().Delete(providerArn); err != nil {
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
	tagsFn: func(r *iamstore.OpenIDConnectProvider) *[]types.Tag { return &r.Tags },
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
