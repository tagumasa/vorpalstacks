package iam

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateSAMLProvider creates a SAML identity provider for the account.
func (s *IAMService) CreateSAMLProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &CreateSAMLProviderInput{
		Name:                    request.GetStringParam(req.Parameters, "Name"),
		SAMLMetadataDocument:    request.GetStringParam(req.Parameters, "SAMLMetadataDocument"),
		AssertionEncryptionMode: request.GetStringParam(req.Parameters, "AssertionEncryptionMode"),
		AddPrivateKey:           request.GetStringParam(req.Parameters, "AddPrivateKey"),
		Tags:                    tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"),
	}
	providerArn, err := s.createSAMLProviderCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"SAMLProviderArn": providerArn,
	}, nil
}

// GetSAMLProvider retrieves information about a SAML identity provider.
func (s *IAMService) GetSAMLProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	provider, err := s.getSAMLProviderCore(store, request.GetStringParam(req.Parameters, "SAMLProviderArn"))
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"SAMLProviderArn":      provider.Arn,
		"SAMLMetadataDocument": provider.SAMLMetadataDocument,
		"CreateDate":           provider.CreateDate.Format(timeutils.ISO8601SimpleFormat),
	}

	// Providers recorded before the UUID member existed carry no value;
	// the member is emitted only when the record holds one.
	if provider.UUID != "" {
		resp["SAMLProviderUUID"] = provider.UUID
	}

	if provider.ValidUntil != nil {
		resp["ValidUntil"] = provider.ValidUntil.Format(timeutils.ISO8601SimpleFormat)
	}

	if provider.AssertionEncryptionMode != "" {
		resp["AssertionEncryptionMode"] = provider.AssertionEncryptionMode
	}

	if len(provider.PrivateKeys) > 0 {
		// PrivateKeyList — each entry contains KeyId and Timestamp
		// metadata only.  Raw private key material is sensitive and not
		// returned (Smithy privateKeyType has sensitive trait).
		pkList := make([]interface{}, len(provider.PrivateKeys))
		for i, pk := range provider.PrivateKeys {
			pkList[i] = map[string]interface{}{
				"KeyId":     pk.KeyId,
				"Timestamp": pk.AddedAt.Format(timeutils.ISO8601SimpleFormat),
			}
		}
		resp["PrivateKeyList"] = pkList
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
	providerList, err := s.listSAMLProvidersCore(store)
	if err != nil {
		return nil, err
	}

	list := make([]interface{}, len(providerList))
	for i, provider := range providerList {
		entry := map[string]interface{}{
			"Arn":        provider.Arn,
			"CreateDate": provider.CreateDate.Format(timeutils.ISO8601SimpleFormat),
		}
		if provider.ValidUntil != nil {
			entry["ValidUntil"] = provider.ValidUntil.Format(timeutils.ISO8601SimpleFormat)
		}
		list[i] = entry
	}

	return map[string]interface{}{
		"SAMLProviderList": list,
	}, nil
}

// UpdateSAMLProvider updates the metadata document for the specified SAML provider.
func (s *IAMService) UpdateSAMLProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &UpdateSAMLProviderInput{
		ProviderArn:             request.GetStringParam(req.Parameters, "SAMLProviderArn"),
		SAMLMetadataDocument:    request.GetStringParam(req.Parameters, "SAMLMetadataDocument"),
		AssertionEncryptionMode: request.GetStringParam(req.Parameters, "AssertionEncryptionMode"),
		AddPrivateKey:           request.GetStringParam(req.Parameters, "AddPrivateKey"),
		RemovePrivateKey:        request.GetStringParam(req.Parameters, "RemovePrivateKey"),
	}
	providerArn, err := s.updateSAMLProviderCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"SAMLProviderArn": providerArn,
	}, nil
}

// DeleteSAMLProvider deletes a SAML identity provider.
func (s *IAMService) DeleteSAMLProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteSAMLProviderCore(store, request.GetStringParam(req.Parameters, "SAMLProviderArn")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

var samlProviderTagOps = tagOps[*iamstore.SAMLProvider]{
	paramName:  "SAMLProviderArn",
	emptyErr:   NewValidationError("SAMLProviderArn"),
	notFoundFn: func(n string) error { return NewNoSuchEntityError("SAML provider", n) },
	getFn:      func(s *iamstore.IAMStore, n string) (*iamstore.SAMLProvider, error) { return s.SAMLProviders().Get(n) },
	putFn:      func(s *iamstore.IAMStore, r *iamstore.SAMLProvider) error { return s.SAMLProviders().Put(r) },
	tagsFn:     func(r *iamstore.SAMLProvider) *[]tags.Tag { return &r.Tags },
}

// TagSAMLProvider adds tags to a SAML provider.
func (s *IAMService) TagSAMLProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagResource(ctx, s, reqCtx, req, samlProviderTagOps)
}

// UntagSAMLProvider removes tags from a SAML provider.
func (s *IAMService) UntagSAMLProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return untagResource(ctx, s, reqCtx, req, samlProviderTagOps)
}

// ListSAMLProviderTags lists the tags attached to a SAML provider.
func (s *IAMService) ListSAMLProviderTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return listResourceTags(ctx, s, reqCtx, req, samlProviderTagOps)
}
