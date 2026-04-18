// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"encoding/json"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/aws/types"
)

const (
	// MaxAccessKeysPerUser is the maximum number of access keys a user can have.
	MaxAccessKeysPerUser = 2
	// MaxPolicyVersions is the maximum number of policy versions allowed.
	MaxPolicyVersions = 5
)

// validatePolicyDocument checks if a policy document is valid JSON.
func validatePolicyDocument(document string) bool {
	if document == "" {
		return false
	}
	var js interface{}
	return json.Unmarshal([]byte(document), &js) == nil
}

type tagOps[T any] struct {
	paramName  string
	emptyErr   error
	notFoundFn func(string) error
	getFn      func(*iamstore.IAMStore, string) (T, error)
	putFn      func(*iamstore.IAMStore, T) error
	tagsFn     func(T) *[]types.Tag
}

func tagResource[T any](ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops tagOps[T]) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, ops.paramName)
	if name == "" {
		return nil, ops.emptyErr
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	res, err := ops.getFn(store, name)
	if err != nil {
		return nil, ops.notFoundFn(name)
	}
	currentTags := ops.tagsFn(res)
	*currentTags = tags.Apply(*currentTags, tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"))
	if err := ops.putFn(store, res); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func untagResource[T any](ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops tagOps[T]) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, ops.paramName)
	if name == "" {
		return nil, ops.emptyErr
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	res, err := ops.getFn(store, name)
	if err != nil {
		return nil, ops.notFoundFn(name)
	}
	currentTags := ops.tagsFn(res)
	*currentTags = tags.RemoveByTagKeys(*currentTags, tags.ParseTagKeysWithQueryFallback(req.Parameters, "TagKeys"))
	if err := ops.putFn(store, res); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func listResourceTags[T any](ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops tagOps[T]) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, ops.paramName)
	if name == "" {
		return nil, ops.emptyErr
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	res, err := ops.getFn(store, name)
	if err != nil {
		return nil, ops.notFoundFn(name)
	}
	return map[string]interface{}{
		"Tags":        tags.ToResponse(*ops.tagsFn(res)),
		"IsTruncated": false,
	}, nil
}
