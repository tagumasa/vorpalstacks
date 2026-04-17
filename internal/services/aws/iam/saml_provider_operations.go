package iam

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateSAMLProvider creates a SAML identity provider for the account.
func (s *IAMService) CreateSAMLProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, NewValidationError("Name")
	}
	metadataDocument := request.GetStringParam(req.Parameters, "SAMLMetadataDocument")
	if metadataDocument == "" {
		return nil, NewValidationError("SAMLMetadataDocument")
	}

	newTags := tags.ParseTagsWithQueryFallback(req.Parameters, "Tags")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	provider, err := store.SAMLProviders().Create(name, metadataDocument, nil, newTags)
	if err != nil {
		if errors.Is(err, iamstore.ErrSAMLProviderAlreadyExists) {
			return nil, NewEntityAlreadyExistsError("SAML Provider " + name)
		}
		return nil, err
	}

	return map[string]interface{}{
		"SAMLProviderArn": provider.Arn,
	}, nil
}

// GetSAMLProvider retrieves information about a SAML identity provider.
func (s *IAMService) GetSAMLProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	providerArn := request.GetStringParam(req.Parameters, "SAMLProviderArn")
	if providerArn == "" {
		return nil, NewValidationError("SAMLProviderArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	provider, err := store.SAMLProviders().Get(providerArn)
	if err != nil {
		return nil, NewNoSuchEntityError("SAML provider", providerArn)
	}

	resp := map[string]interface{}{
		"SAMLProviderArn":      provider.Arn,
		"SAMLMetadataDocument": provider.SAMLMetadataDocument,
		"CreateDate":           provider.CreateDate.Format(timeutils.ISO8601SimpleFormat),
	}

	if provider.ValidUntil != nil {
		resp["ValidUntil"] = provider.ValidUntil.Format(timeutils.ISO8601SimpleFormat)
	}

	if provider.Tags != nil {
		resp["Tags"] = tags.ToResponse(provider.Tags)
	}

	return resp, nil
}

// ListSAMLProviders lists the SAML providers in the account.
func (s *IAMService) ListSAMLProviders(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.SAMLProviders().List()
	if err != nil {
		return nil, err
	}

	providerList := make([]interface{}, len(result.SAMLProviders))
	for i, provider := range result.SAMLProviders {
		entry := map[string]interface{}{
			"Arn":        provider.Arn,
			"CreateDate": provider.CreateDate.Format(timeutils.ISO8601SimpleFormat),
		}
		if provider.ValidUntil != nil {
			entry["ValidUntil"] = provider.ValidUntil.Format(timeutils.ISO8601SimpleFormat)
		}
		providerList[i] = entry
	}

	return map[string]interface{}{
		"SAMLProviderList": providerList,
	}, nil
}

// UpdateSAMLProvider updates the metadata document for the specified SAML provider.
func (s *IAMService) UpdateSAMLProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	providerArn := request.GetStringParam(req.Parameters, "SAMLProviderArn")
	if providerArn == "" {
		return nil, NewValidationError("SAMLProviderArn")
	}
	metadataDocument := request.GetStringParam(req.Parameters, "SAMLMetadataDocument")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.SAMLProviders().Exists(providerArn) {
		return nil, NewNoSuchEntityError("SAML provider", providerArn)
	}

	if err := store.SAMLProviders().Update(providerArn, metadataDocument, nil); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"SAMLProviderArn": providerArn,
	}, nil
}

// DeleteSAMLProvider deletes a SAML identity provider.
func (s *IAMService) DeleteSAMLProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	providerArn := request.GetStringParam(req.Parameters, "SAMLProviderArn")
	if providerArn == "" {
		return nil, NewValidationError("SAMLProviderArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.SAMLProviders().Exists(providerArn) {
		return nil, NewNoSuchEntityError("SAML provider", providerArn)
	}
	if err := store.SAMLProviders().Delete(providerArn); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// TagSAMLProvider adds tags to a SAML provider.
func (s *IAMService) TagSAMLProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	providerArn := request.GetStringParam(req.Parameters, "SAMLProviderArn")
	if providerArn == "" {
		return nil, NewValidationError("SAMLProviderArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	provider, err := store.SAMLProviders().Get(providerArn)
	if err != nil {
		return nil, NewNoSuchEntityError("SAML provider", providerArn)
	}

	provider.Tags = tags.Apply(provider.Tags, tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"))
	if err := store.SAMLProviders().Put(provider); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UntagSAMLProvider removes tags from a SAML provider.
func (s *IAMService) UntagSAMLProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	providerArn := request.GetStringParam(req.Parameters, "SAMLProviderArn")
	if providerArn == "" {
		return nil, NewValidationError("SAMLProviderArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	provider, err := store.SAMLProviders().Get(providerArn)
	if err != nil {
		return nil, NewNoSuchEntityError("SAML provider", providerArn)
	}

	provider.Tags = tags.RemoveByTagKeys(provider.Tags, tags.ParseTagKeysWithQueryFallback(req.Parameters, "TagKeys"))
	if err := store.SAMLProviders().Put(provider); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListSAMLProviderTags lists the tags attached to a SAML provider.
func (s *IAMService) ListSAMLProviderTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	providerArn := request.GetStringParam(req.Parameters, "SAMLProviderArn")
	if providerArn == "" {
		return nil, NewValidationError("SAMLProviderArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	provider, err := store.SAMLProviders().Get(providerArn)
	if err != nil {
		return nil, NewNoSuchEntityError("SAML provider", providerArn)
	}

	return map[string]interface{}{
		"Tags":        tags.ToResponse(provider.Tags),
		"IsTruncated": false,
	}, nil
}
