package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"
	"errors"
	"strings"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

func (s *KMSService) mapAliasStoreError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, kmsstore.ErrAliasAlreadyExists):
		return ErrAliasAlreadyExists
	case errors.Is(err, kmsstore.ErrAliasNotFound):
		return ErrAliasNotFound
	case errors.Is(err, kmsstore.ErrAliasNameReserved):
		return ErrInvalidAliasName
	case errors.Is(err, kmsstore.ErrInvalidKeyState):
		return ErrKeyPendingDeletion
	default:
		return err
	}
}

// CreateAlias creates a display name for a KMS key.
func (s *KMSService) CreateAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	aliasName := request.GetStringParam(req.Parameters, "AliasName")
	targetKeyID := request.GetStringParam(req.Parameters, "TargetKeyId")

	if aliasName == "" {
		return nil, ErrInvalidAliasName
	}

	if !strings.HasPrefix(aliasName, "alias/") {
		return nil, ErrInvalidAliasName
	}

	if strings.HasPrefix(aliasName, "alias/aws/") {
		return nil, ErrInvalidAliasName
	}

	if len(aliasName) > 256 {
		return nil, ErrInvalidAliasName
	}

	if targetKeyID == "" {
		return nil, ErrKeyNotFound
	}

	key, err := s.resolveKeyByKeyID(stores, targetKeyID)
	if err != nil {
		return nil, err
	}
	if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "CreateAlias", key.KeyID, nil); err != nil {
		return nil, err
	}
	if key.KeyState == kmsstore.KeyStatePendingDeletion {
		return nil, ErrKeyPendingDeletion
	}
	_, err = stores.aliases.Create(aliasName, key.KeyID)
	if err != nil {
		return nil, s.mapAliasStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// DeleteAlias deletes the specified alias.
func (s *KMSService) DeleteAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	aliasName := request.GetStringParam(req.Parameters, "AliasName")
	if aliasName == "" {
		return nil, ErrInvalidAliasName
	}

	if strings.HasPrefix(aliasName, "alias/aws/") {
		return nil, ErrInvalidAliasName
	}

	if !stores.aliases.Exists(aliasName) {
		return nil, ErrAliasNotFound
	}

	if err := stores.aliases.Delete(aliasName); err != nil {
		return nil, s.mapAliasStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// ListAliases retrieves a list of aliases in the account and region.
func (s *KMSService) ListAliases(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	marker := pagination.GetMarker(req.Parameters)
	maxItems := pagination.GetMaxItems(req.Parameters, 100)
	keyID := s.getKeyID(req.Parameters)

	result, err := stores.aliases.List(marker, maxItems, keyID)
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
	aliasName := request.GetStringParam(req.Parameters, "AliasName")
	targetKeyID := request.GetStringParam(req.Parameters, "TargetKeyId")

	if aliasName == "" {
		return nil, ErrInvalidAliasName
	}

	if !strings.HasPrefix(aliasName, "alias/") {
		return nil, ErrInvalidAliasName
	}

	if strings.HasPrefix(aliasName, "alias/aws/") {
		return nil, ErrInvalidAliasName
	}

	if len(aliasName) > 256 {
		return nil, ErrInvalidAliasName
	}

	if targetKeyID == "" {
		return nil, ErrKeyNotFound
	}

	key, err := s.resolveKeyByKeyID(stores, targetKeyID)
	if err != nil {
		return nil, err
	}
	if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "UpdateAlias", key.KeyID, nil); err != nil {
		return nil, err
	}
	if key.KeyState == kmsstore.KeyStatePendingDeletion {
		return nil, ErrKeyPendingDeletion
	}
	if err := stores.aliases.UpdateTarget(aliasName, key.KeyID); err != nil {
		return nil, s.mapAliasStoreError(err)
	}

	return response.EmptyResponse(), nil
}
