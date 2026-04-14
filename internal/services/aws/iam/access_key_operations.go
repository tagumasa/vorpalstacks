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
func (s *IAMService) CreateAccessKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	currentKeyCount, err := store.AccessKeys().CountByUserName(userName)
	if err != nil {
		return nil, err
	}
	if currentKeyCount >= MaxAccessKeysPerUser {
		return nil, ErrAccessKeyLimitExceeded
	}

	key, err := store.AccessKeys().Create(userName)
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
	if accessKeyId == "" {
		return nil, ErrNoSuchAccessKey
	}

	userName := request.GetStringParam(req.Parameters, "UserName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := store.AccessKeys().Get(accessKeyId)
	if err != nil {
		return nil, NewNoSuchAccessKeyError(accessKeyId)
	}

	if userName != "" && key.UserName != userName {
		return nil, NewNoSuchAccessKeyError(accessKeyId)
	}

	if err := store.AccessKeys().Delete(accessKeyId); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListAccessKeys lists the access keys for the specified user.
func (s *IAMService) ListAccessKeys(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := getMaxItems(req.Parameters)

	keys, err := store.AccessKeys().ListByUserName(userName)
	if err != nil {
		return nil, err
	}

	result := pagination.PaginateSlice(keys, marker, maxItems, func(k *iamstore.AccessKey) string {
		return k.AccessKeyId
	})

	accessKeyMetadata := make([]interface{}, len(result.Items))
	for i, key := range result.Items {
		accessKeyMetadata[i] = s.accessKeyToMetadata(key)
	}

	response := map[string]interface{}{
		"AccessKeyMetadata": accessKeyMetadata,
		"IsTruncated":       result.IsTruncated,
	}

	if result.IsTruncated && len(result.Items) > 0 {
		response["Marker"] = result.NextMarker
	}

	return response, nil
}

// GetAccessKeyLastUsed returns information about when the access key was last used.
func (s *IAMService) GetAccessKeyLastUsed(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessKeyId := request.GetStringParam(req.Parameters, "AccessKeyId")
	if accessKeyId == "" {
		return nil, ErrNoSuchAccessKey
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := store.AccessKeys().Get(accessKeyId)
	if err != nil {
		return nil, NewNoSuchAccessKeyError(accessKeyId)
	}

	response := map[string]interface{}{
		"UserName": key.UserName,
	}

	if key.LastUsedDate != nil {
		response["AccessKeyLastUsed"] = map[string]interface{}{
			"LastUsedDate": key.LastUsedDate.Format(timeutils.ISO8601SimpleFormat),
			"Region":       key.LastUsedRegion,
			"ServiceName":  key.LastUsedService,
		}
	} else {
		response["AccessKeyLastUsed"] = map[string]interface{}{
			"Region":      "N/A",
			"ServiceName": "N/A",
		}
	}

	return response, nil
}

// UpdateAccessKey updates the status of the specified access key.
func (s *IAMService) UpdateAccessKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessKeyId := request.GetStringParam(req.Parameters, "AccessKeyId")
	if accessKeyId == "" {
		return nil, ErrNoSuchAccessKey
	}

	userName := request.GetStringParam(req.Parameters, "UserName")
	status := request.GetStringParam(req.Parameters, "Status")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := store.AccessKeys().Get(accessKeyId)
	if err != nil {
		return nil, NewNoSuchAccessKeyError(accessKeyId)
	}

	if userName != "" && key.UserName != userName {
		return nil, NewNoSuchAccessKeyError(accessKeyId)
	}

	newStatus := iamstore.AccessKeyStatus(status)
	if newStatus != iamstore.AccessKeyStatusActive && newStatus != iamstore.AccessKeyStatusInactive {
		return nil, NewInvalidInputError("Status", "must be Active or Inactive")
	}

	if err := store.AccessKeys().UpdateStatus(accessKeyId, newStatus); err != nil {
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
