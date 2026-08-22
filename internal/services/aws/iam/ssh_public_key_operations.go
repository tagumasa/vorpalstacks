package iam

import (
	"context"
	"errors"
	"fmt"
	"unicode/utf8"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// UploadSSHPublicKey uploads an SSH public key for the specified IAM user.
func (s *IAMService) UploadSSHPublicKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, NewValidationError("UserName")
	}
	sshPublicKeyBody := request.GetStringParam(req.Parameters, "SSHPublicKeyBody")
	if sshPublicKeyBody == "" {
		return nil, NewValidationError("SSHPublicKeyBody")
	}
	// publicKeyMaterialType carries a Latin-1 pattern, so lengths count
	// Unicode characters.
	if utf8.RuneCountInString(sshPublicKeyBody) > maxSSHPublicKeyLength {
		return nil, NewInvalidInputError("SSHPublicKeyBody", fmt.Sprintf("must be 1 to %d characters", maxSSHPublicKeyLength))
	}

	parsedKey, err := parseSSHPublicKey(sshPublicKeyBody)
	if err != nil {
		return nil, ErrInvalidPublicKey
	}
	canonicalBody := canonicalSSHPublicKeyBody(parsedKey)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	key, err := store.SSHPublicKeys().UploadWithGuards(userName, canonicalBody)
	if err != nil {
		if errors.Is(err, iamstore.ErrDuplicateSSHPublicKey) {
			return nil, ErrDuplicateSSHPublicKey
		}
		if errors.Is(err, iamstore.ErrSSHPublicKeyLimitExceeded) {
			return nil, ErrLimitExceededSSHPublicKeys
		}
		return nil, err
	}

	// The output shape referenced by UploadSSHPublicKey includes the key
	// body; return it in the canonical SSH form that was stored.
	return map[string]interface{}{
		"SSHPublicKey": s.sshPublicKeyToResponse(key, "SSH"),
	}, nil
}

// GetSSHPublicKey retrieves the specified SSH public key for the specified IAM user.
func (s *IAMService) GetSSHPublicKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	keyId := request.GetStringParam(req.Parameters, "SSHPublicKeyId")
	if keyId == "" {
		return nil, NewValidationError("SSHPublicKeyId")
	}
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, NewValidationError("UserName")
	}
	encoding := request.GetStringParam(req.Parameters, "Encoding")
	if encoding == "" {
		return nil, NewValidationError("Encoding")
	}
	if encoding != "SSH" && encoding != "PEM" {
		return nil, NewInvalidInputError("Encoding", "must be SSH or PEM")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	key, err := store.SSHPublicKeys().Get(keyId)
	if err != nil {
		return nil, NewNoSuchEntityError("SSH public key", keyId)
	}

	return map[string]interface{}{
		"SSHPublicKey": s.sshPublicKeyToResponse(key, encoding),
	}, nil
}

// UpdateSSHPublicKey changes the status of the specified SSH public key to Active or Inactive.
func (s *IAMService) UpdateSSHPublicKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	keyId := request.GetStringParam(req.Parameters, "SSHPublicKeyId")
	if keyId == "" {
		return nil, NewValidationError("SSHPublicKeyId")
	}
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, NewValidationError("UserName")
	}
	status := request.GetStringParam(req.Parameters, "Status")
	if status == "" {
		return nil, NewValidationError("Status")
	}
	if status != "Active" && status != "Inactive" {
		return nil, NewInvalidInputError("Status", "must be Active or Inactive")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}
	key, err := store.SSHPublicKeys().Get(keyId)
	if err != nil {
		return nil, NewNoSuchEntityError("SSH public key", keyId)
	}
	// The named user must own the key; otherwise the operation reports the
	// key as not existing for that user.
	if key.UserName != userName {
		return nil, NewNoSuchEntityError("SSH public key", keyId)
	}

	if err := store.SSHPublicKeys().UpdateStatus(keyId, status); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListSSHPublicKeys returns information about the SSH public keys associated with the specified IAM user.
func (s *IAMService) ListSSHPublicKeys(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	userName, err := resolveUserName(reqCtx, userName)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	keys, err := store.SSHPublicKeys().ListByUserName(userName)
	if err != nil {
		return nil, err
	}

	keyList := make([]interface{}, len(keys))
	for i, key := range keys {
		keyList[i] = s.sshPublicKeyToResponse(key, "")
	}

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	paged := pagination.PaginateSlice(keyList, marker, maxItems, func(item interface{}) string {
		if m, ok := item.(map[string]interface{}); ok {
			if id, ok := m["SSHPublicKeyId"].(string); ok {
				return id
			}
		}
		return ""
	})

	resp := map[string]interface{}{
		"SSHPublicKeys": paged.Items,
		"IsTruncated":   paged.IsTruncated,
	}
	if paged.NextMarker != "" {
		resp["Marker"] = paged.NextMarker
	}
	return resp, nil
}

// DeleteSSHPublicKey deletes the specified SSH public key for the specified IAM user.
func (s *IAMService) DeleteSSHPublicKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	keyId := request.GetStringParam(req.Parameters, "SSHPublicKeyId")
	if keyId == "" {
		return nil, NewValidationError("SSHPublicKeyId")
	}
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, NewValidationError("UserName")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}
	key, err := store.SSHPublicKeys().Get(keyId)
	if err != nil {
		return nil, NewNoSuchEntityError("SSH public key", keyId)
	}
	// The named user must own the key; otherwise the operation reports the
	// key as not existing for that user.
	if key.UserName != userName {
		return nil, NewNoSuchEntityError("SSH public key", keyId)
	}

	if err := store.SSHPublicKeys().Delete(keyId); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// sshPublicKeyToResponse renders an SSH public key. An empty encoding
// omits the body (listing responses); "SSH" returns the canonical single
// line form; "PEM" returns a SubjectPublicKeyInfo PEM block.
func (s *IAMService) sshPublicKeyToResponse(key *iamstore.SSHPublicKey, encoding string) map[string]interface{} {
	resp := map[string]interface{}{
		"SSHPublicKeyId": key.SSHPublicKeyId,
		"UserName":       key.UserName,
		"Fingerprint":    key.Fingerprint,
		"Status":         key.Status,
		"UploadDate":     key.UploadDate.Format(timeutils.ISO8601SimpleFormat),
	}
	switch encoding {
	case "SSH":
		resp["SSHPublicKeyBody"] = key.SSHPublicKeyBody
	case "PEM":
		if parsed, err := parseSSHPublicKey(key.SSHPublicKeyBody); err == nil {
			if pemBody, err := sshPublicKeyBodyPEM(parsed); err == nil {
				resp["SSHPublicKeyBody"] = pemBody
			}
		}
	}
	return resp
}
