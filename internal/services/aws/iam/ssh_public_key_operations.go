package iam

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// UploadSSHPublicKey uploads an SSH public key for the specified IAM user.
func (s *IAMService) UploadSSHPublicKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := &UploadSSHPublicKeyInput{
		UserName:         request.GetStringParam(req.Parameters, "UserName"),
		SSHPublicKeyBody: request.GetStringParam(req.Parameters, "SSHPublicKeyBody"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.uploadSSHPublicKeyCore(store, input)
	if err != nil {
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
	input := &GetSSHPublicKeyInput{
		UserName:       request.GetStringParam(req.Parameters, "UserName"),
		SSHPublicKeyId: request.GetStringParam(req.Parameters, "SSHPublicKeyId"),
		Encoding:       request.GetStringParam(req.Parameters, "Encoding"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.getSSHPublicKeyCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"SSHPublicKey": s.sshPublicKeyToResponse(key, input.Encoding),
	}, nil
}

// UpdateSSHPublicKey changes the status of the specified SSH public key to Active or Inactive.
func (s *IAMService) UpdateSSHPublicKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := &UpdateSSHPublicKeyInput{
		UserName:       request.GetStringParam(req.Parameters, "UserName"),
		SSHPublicKeyId: request.GetStringParam(req.Parameters, "SSHPublicKeyId"),
		Status:         request.GetStringParam(req.Parameters, "Status"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.updateSSHPublicKeyCore(store, input); err != nil {
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
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listSSHPublicKeysCore(store, userName, marker, maxItems)
	if err != nil {
		return nil, err
	}

	keyList := make([]interface{}, len(result.Keys))
	for i, key := range result.Keys {
		keyList[i] = s.sshPublicKeyToResponse(key, "")
	}

	resp := map[string]interface{}{
		"SSHPublicKeys": keyList,
		"IsTruncated":   result.IsTruncated,
	}
	if result.NextMarker != "" {
		resp["Marker"] = result.NextMarker
	}
	return resp, nil
}

// DeleteSSHPublicKey deletes the specified SSH public key for the specified IAM user.
func (s *IAMService) DeleteSSHPublicKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := &UpdateSSHPublicKeyInput{
		UserName:       request.GetStringParam(req.Parameters, "UserName"),
		SSHPublicKeyId: request.GetStringParam(req.Parameters, "SSHPublicKeyId"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteSSHPublicKeyCore(store, input); err != nil {
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
