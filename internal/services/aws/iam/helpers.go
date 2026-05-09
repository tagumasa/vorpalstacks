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

// buildAttachedManagedPolicies returns the attached managed policy list for
// a principal, collapsing the copy-paste across User/Group/Role in
// GetAccountAuthorizationDetails.
func buildAttachedManagedPolicies(store *iamstore.IAMStore, principalType, principalName string) []interface{} {
	arns, _ := store.AttachedPolicies().ListAttachedPolicies(principalType, principalName)
	policies := make([]interface{}, 0, len(arns))
	for _, arn := range arns {
		if p, err := store.Policies().Get(arn); err == nil {
			policies = append(policies, map[string]interface{}{
				"PolicyName": p.PolicyName,
				"PolicyArn":  p.Arn,
			})
		}
	}
	return policies
}

// buildInlinePolicyList returns the inline policy name list for a principal.
func buildInlinePolicyList(store *iamstore.IAMStore, principalType, principalName string) []interface{} {
	names, _ := store.InlinePolicies().List(principalType, principalName)
	list := make([]interface{}, 0, len(names))
	for _, pn := range names {
		list = append(list, map[string]interface{}{
			"PolicyName": pn,
		})
	}
	return list
}

// listAllUsers paginates through all users matching pathPrefix.
func listAllUsers(store *iamstore.IAMStore, pathPrefix string) ([]*iamstore.User, error) {
	var all []*iamstore.User
	marker := ""
	for {
		result, err := store.Users().List(pathPrefix, marker, 1000)
		if err != nil {
			return nil, err
		}
		all = append(all, result.Users...)
		if !result.IsTruncated {
			break
		}
		marker = result.Marker
	}
	return all, nil
}

// listAllGroups paginates through all groups matching pathPrefix.
func listAllGroups(store *iamstore.IAMStore, pathPrefix string) ([]*iamstore.Group, error) {
	var all []*iamstore.Group
	marker := ""
	for {
		result, err := store.Groups().List(pathPrefix, marker, 1000)
		if err != nil {
			return nil, err
		}
		all = append(all, result.Groups...)
		if !result.IsTruncated {
			break
		}
		marker = result.Marker
	}
	return all, nil
}

// listAllRoles paginates through all roles matching pathPrefix.
func listAllRoles(store *iamstore.IAMStore, pathPrefix string) ([]*iamstore.Role, error) {
	var all []*iamstore.Role
	marker := ""
	for {
		result, err := store.Roles().List(pathPrefix, marker, 1000)
		if err != nil {
			return nil, err
		}
		all = append(all, result.Roles...)
		if !result.IsTruncated {
			break
		}
		marker = result.Marker
	}
	return all, nil
}

// listAllPolicies paginates through all policies matching the given filters.
func listAllPolicies(store *iamstore.IAMStore, scope, pathPrefix string, onlyAttached bool) ([]*iamstore.Policy, error) {
	var all []*iamstore.Policy
	marker := ""
	for {
		result, err := store.Policies().List(scope, pathPrefix, onlyAttached, marker, 1000)
		if err != nil {
			return nil, err
		}
		all = append(all, result.Policies...)
		if !result.IsTruncated {
			break
		}
		marker = result.Marker
	}
	return all, nil
}
