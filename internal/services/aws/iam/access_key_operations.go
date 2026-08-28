// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateAccessKey creates a new access key for the specified user.
// The per-user access key quota is enforced atomically inside
// the store layer to prevent race conditions on concurrent requests.
func (s *IAMService) CreateAccessKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	userName, err := resolveUserName(reqCtx, userName)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.createAccessKeyCore(store, userName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"AccessKey": map[string]interface{}{
			"UserName":        key.UserName,
			"AccessKeyId":     key.AccessKeyId,
			"Status":          string(key.Status),
			"SecretAccessKey": key.SecretAccessKey,
			"CreateDate":      key.CreateDate.Format(timeutils.ISO8601SimpleFormat),
		},
	}, nil
}

// DeleteAccessKey deletes the specified access key.
func (s *IAMService) DeleteAccessKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessKeyId := request.GetStringParam(req.Parameters, "AccessKeyId")
	userName := request.GetStringParam(req.Parameters, "UserName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteAccessKeyCore(store, accessKeyId, userName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListAccessKeys lists the access keys for the specified user.
func (s *IAMService) ListAccessKeys(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	result, err := s.listAccessKeysCore(store, userName, marker, maxItems)
	if err != nil {
		return nil, err
	}

	accessKeyMetadata := make([]interface{}, len(result.Keys))
	for i, key := range result.Keys {
		accessKeyMetadata[i] = s.accessKeyToMetadata(key)
	}

	resp := map[string]interface{}{
		"AccessKeyMetadata": accessKeyMetadata,
		"IsTruncated":       result.IsTruncated,
	}

	if result.IsTruncated && len(result.Keys) > 0 {
		resp["Marker"] = result.NextMarker
	}

	return resp, nil
}

// GetAccessKeyLastUsed returns information about when the access key was last used.
func (s *IAMService) GetAccessKeyLastUsed(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessKeyId := request.GetStringParam(req.Parameters, "AccessKeyId")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.getAccessKeyLastUsedCore(store, accessKeyId)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"UserName": key.UserName,
	}

	if key.LastUsedDate != nil {
		resp["AccessKeyLastUsed"] = map[string]interface{}{
			"LastUsedDate": key.LastUsedDate.Format(timeutils.ISO8601SimpleFormat),
			"Region":       key.LastUsedRegion,
			"ServiceName":  key.LastUsedService,
		}
	} else {
		resp["AccessKeyLastUsed"] = map[string]interface{}{
			"Region":      "N/A",
			"ServiceName": "N/A",
		}
	}

	return resp, nil
}

// UpdateAccessKey updates the status of the specified access key.
func (s *IAMService) UpdateAccessKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := &UpdateAccessKeyInput{
		AccessKeyId: request.GetStringParam(req.Parameters, "AccessKeyId"),
		UserName:    request.GetStringParam(req.Parameters, "UserName"),
		Status:      request.GetStringParam(req.Parameters, "Status"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.updateAccessKeyCore(store, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func (s *IAMService) accessKeyToMetadata(key *iamstore.AccessKey) map[string]interface{} {
	return map[string]interface{}{
		"UserName":    key.UserName,
		"AccessKeyId": key.AccessKeyId,
		"Status":      string(key.Status),
		"CreateDate":  key.CreateDate.Format(timeutils.ISO8601SimpleFormat),
	}
}
