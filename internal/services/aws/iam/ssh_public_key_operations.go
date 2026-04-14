package iam

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// UploadSSHPublicKey uploads an SSH public key for the specified IAM user.
func (s *IAMService) UploadSSHPublicKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, ErrNoSuchUser
	}
	sshPublicKeyBody := request.GetStringParam(req.Parameters, "SSHPublicKeyBody")
	if sshPublicKeyBody == "" {
		return nil, NewValidationError("SSHPublicKeyBody")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	key, err := store.SSHPublicKeys().Upload(userName, sshPublicKeyBody)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"SSHPublicKey": s.sshPublicKeyToResponse(key, true),
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
		return nil, ErrNoSuchUser
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
		"SSHPublicKey": s.sshPublicKeyToResponse(key, true),
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
		return nil, ErrNoSuchUser
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
	if !store.SSHPublicKeys().Exists(keyId) {
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
	if userName == "" {
		return nil, ErrNoSuchUser
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
		keyList[i] = s.sshPublicKeyToResponse(key, false)
	}

	return map[string]interface{}{
		"SSHPublicKeys": keyList,
		"IsTruncated":   false,
	}, nil
}

// DeleteSSHPublicKey deletes the specified SSH public key for the specified IAM user.
func (s *IAMService) DeleteSSHPublicKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	keyId := request.GetStringParam(req.Parameters, "SSHPublicKeyId")
	if keyId == "" {
		return nil, NewValidationError("SSHPublicKeyId")
	}
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, ErrNoSuchUser
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}
	if !store.SSHPublicKeys().Exists(keyId) {
		return nil, NewNoSuchEntityError("SSH public key", keyId)
	}

	if err := store.SSHPublicKeys().Delete(keyId); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func (s *IAMService) sshPublicKeyToResponse(key *iamstore.SSHPublicKey, includeBody bool) map[string]interface{} {
	resp := map[string]interface{}{
		"SSHPublicKeyId": key.SSHPublicKeyId,
		"UserName":       key.UserName,
		"Fingerprint":    key.Fingerprint,
		"Status":         key.Status,
		"UploadDate":     key.UploadDate.Format(timeutils.ISO8601SimpleFormat),
	}
	if includeBody {
		resp["SSHPublicKeyBody"] = key.SSHPublicKeyBody
	}
	return resp
}
