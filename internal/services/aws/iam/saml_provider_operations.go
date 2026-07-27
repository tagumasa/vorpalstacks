package iam

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"regexp"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/timeutils"
)

var x509CertDataPattern = regexp.MustCompile(`(?s)<(?:ds:)?X509Certificate>([^<]+)</(?:ds:)?X509Certificate>`)
var whitespacePattern = regexp.MustCompile(`\s+`)

func extractValidUntilFromSAMLMetadata(metadata string) *time.Time {
	matches := x509CertDataPattern.FindStringSubmatch(metadata)
	if len(matches) < 2 {
		return nil
	}
	certData := whitespacePattern.ReplaceAllString(matches[1], "")

	derBytes, err := base64.StdEncoding.DecodeString(certData)
	if err != nil {
		pemBlock, _ := pem.Decode([]byte("-----BEGIN CERTIFICATE-----\n" + certData + "\n-----END CERTIFICATE-----"))
		if pemBlock == nil {
			return nil
		}
		derBytes = pemBlock.Bytes
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil
	}

	notAfter := cert.NotAfter.UTC()
	return &notAfter
}

// CreateSAMLProvider creates a SAML identity provider for the account.
func (s *IAMService) CreateSAMLProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, NewValidationError("Name")
	}
	if !samlProviderNamePattern.MatchString(name) {
		return nil, NewInvalidInputError("Name", "must be 1 to 128 alphanumeric characters or any of ._-")
	}
	metadataDocument := request.GetStringParam(req.Parameters, "SAMLMetadataDocument")
	if metadataDocument == "" {
		return nil, NewValidationError("SAMLMetadataDocument")
	}
	if len(metadataDocument) < 1000 || len(metadataDocument) > 10000000 {
		return nil, NewInvalidInputError("SAMLMetadataDocument", "must be between 1000 and 10000000 characters")
	}

	newTags := tags.ParseTagsWithQueryFallback(req.Parameters, "Tags")
	if err := validateNewTags(newTags); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	validUntil := extractValidUntilFromSAMLMetadata(metadataDocument)
	provider, err := store.SAMLProviders().Create(name, metadataDocument, validUntil, newTags)
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

	validUntil := extractValidUntilFromSAMLMetadata(metadataDocument)
	if err := store.SAMLProviders().Update(providerArn, metadataDocument, validUntil); err != nil {
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

var samlProviderTagOps = tagOps[*iamstore.SAMLProvider]{
	paramName:  "SAMLProviderArn",
	emptyErr:   NewValidationError("SAMLProviderArn"),
	notFoundFn: func(n string) error { return NewNoSuchEntityError("SAML provider", n) },
	getFn:      func(s *iamstore.IAMStore, n string) (*iamstore.SAMLProvider, error) { return s.SAMLProviders().Get(n) },
	putFn:      func(s *iamstore.IAMStore, r *iamstore.SAMLProvider) error { return s.SAMLProviders().Put(r) },
	tagsFn:     func(r *iamstore.SAMLProvider) *[]types.Tag { return &r.Tags },
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
