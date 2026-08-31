package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"
	"regexp"
	"strings"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// aliasNamePattern mirrors the AWS AliasName constraint: after the required
// "alias/" prefix, only alphanumerics, forward slashes, underscores, and
// hyphens are permitted. Total length 1-256 characters.
var aliasNamePattern = regexp.MustCompile(`^alias/[a-zA-Z0-9/_-]+$`)

// validateAliasName enforces the AWS AliasName format rules:
//   - must be 1-256 characters
//   - must start with "alias/"
//   - must not start with "alias/aws/" (reserved for AWS-managed aliases)
//   - the suffix after "alias/" may contain only [a-zA-Z0-9/_-]
func validateAliasName(aliasName string) error {
	if aliasName == "" || len(aliasName) > 256 {
		return ErrInvalidAliasName
	}
	if !strings.HasPrefix(aliasName, "alias/") {
		return ErrInvalidAliasName
	}
	if strings.HasPrefix(aliasName, "alias/aws/") {
		return ErrInvalidAliasName
	}
	if !aliasNamePattern.MatchString(aliasName) {
		return ErrInvalidAliasName
	}
	return nil
}

// CreateAlias creates a display name for a KMS key.
func (s *KMSService) CreateAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.createAliasCore(stores, s.resolveCallerPrincipal(reqCtx, req),
		request.GetStringParam(req.Parameters, "AliasName"),
		request.GetStringParam(req.Parameters, "TargetKeyId")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteAlias deletes the specified alias.
func (s *KMSService) DeleteAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteAliasCore(stores, s.resolveCallerPrincipal(reqCtx, req),
		request.GetStringParam(req.Parameters, "AliasName")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListAliases retrieves a list of aliases in the account and region.
func (s *KMSService) ListAliases(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listAliasesCore(stores,
		pagination.GetMarker(req.Parameters),
		pagination.GetMaxItems(req.Parameters, 100),
		s.getKeyID(req.Parameters))
	if err != nil {
		return nil, err
	}

	aliases := make([]map[string]interface{}, len(result.Aliases))
	for i, a := range result.Aliases {
		aliasEntry := map[string]interface{}{
			"AliasName": a.AliasName,
			"AliasArn":  a.AliasArn,
		}
		if a.TargetKeyID != "" {
			aliasEntry["TargetKeyId"] = a.TargetKeyID
		}
		if !a.CreationDate.IsZero() {
			aliasEntry["CreationDate"] = a.CreationDate.Unix()
		}
		if !a.LastUpdatedDate.IsZero() {
			aliasEntry["LastUpdatedDate"] = a.LastUpdatedDate.Unix()
		}
		aliases[i] = aliasEntry
	}

	response := map[string]interface{}{
		"Aliases":   aliases,
		"Truncated": result.IsTruncated,
	}
	if result.IsTruncated {
		response["NextMarker"] = result.NextMarker
	}

	return response, nil
}

// UpdateAlias updates the alias to point to a new KMS key.
func (s *KMSService) UpdateAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.updateAliasCore(stores, s.resolveCallerPrincipal(reqCtx, req),
		request.GetStringParam(req.Parameters, "AliasName"),
		request.GetStringParam(req.Parameters, "TargetKeyId")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
