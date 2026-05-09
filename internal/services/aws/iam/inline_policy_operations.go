package iam

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// PutUserPolicy creates or updates an inline policy for a user.
func (s *IAMService) PutUserPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return putInlinePolicy(ctx, s, reqCtx, req, userInlinePolicyOps)
}

// GetUserPolicy retrieves an inline policy for a user.
func (s *IAMService) GetUserPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return getInlinePolicy(ctx, s, reqCtx, req, userInlinePolicyOps)
}

// DeleteUserPolicy deletes an inline policy from a user.
func (s *IAMService) DeleteUserPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return deleteInlinePolicy(ctx, s, reqCtx, req, userInlinePolicyOps)
}

// ListUserPolicies lists inline policies for a user.
func (s *IAMService) ListUserPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return listInlinePolicies(ctx, s, reqCtx, req, userInlinePolicyOps)
}

// PutGroupPolicy creates or updates an inline policy for a group.
func (s *IAMService) PutGroupPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return putInlinePolicy(ctx, s, reqCtx, req, groupInlinePolicyOps)
}

// GetGroupPolicy retrieves an inline policy for a group.
func (s *IAMService) GetGroupPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return getInlinePolicy(ctx, s, reqCtx, req, groupInlinePolicyOps)
}

// DeleteGroupPolicy deletes an inline policy from a group.
func (s *IAMService) DeleteGroupPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return deleteInlinePolicy(ctx, s, reqCtx, req, groupInlinePolicyOps)
}

// ListGroupPolicies lists inline policies for a group.
func (s *IAMService) ListGroupPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return listInlinePolicies(ctx, s, reqCtx, req, groupInlinePolicyOps)
}

// PutRolePolicy creates or updates an inline policy for a role.
func (s *IAMService) PutRolePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return putInlinePolicy(ctx, s, reqCtx, req, roleInlinePolicyOps)
}

// GetRolePolicy retrieves an inline policy for a role.
func (s *IAMService) GetRolePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return getInlinePolicy(ctx, s, reqCtx, req, roleInlinePolicyOps)
}

// DeleteRolePolicy deletes an inline policy from a role.
func (s *IAMService) DeleteRolePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return deleteInlinePolicy(ctx, s, reqCtx, req, roleInlinePolicyOps)
}

// ListRolePolicies lists inline policies for a role.
func (s *IAMService) ListRolePolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return listInlinePolicies(ctx, s, reqCtx, req, roleInlinePolicyOps)
}
