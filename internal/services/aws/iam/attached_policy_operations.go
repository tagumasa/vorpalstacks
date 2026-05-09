package iam

import (
	"context"

	"vorpalstacks/internal/common/request"
)

const (
	// PrincipalTypeUser represents an IAM user principal type.
	PrincipalTypeUser = "user"
	// PrincipalTypeGroup represents an IAM group principal type.
	PrincipalTypeGroup = "group"
	// PrincipalTypeRole represents an IAM role principal type.
	PrincipalTypeRole = "role"
)

// AttachUserPolicy attaches a policy to a user.
func (s *IAMService) AttachUserPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return attachPolicy(ctx, s, reqCtx, req, userAttachOps)
}

// DetachUserPolicy detaches a policy from a user.
func (s *IAMService) DetachUserPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return detachPolicy(ctx, s, reqCtx, req, userAttachOps)
}

// ListAttachedUserPolicies lists the policies attached to a user.
func (s *IAMService) ListAttachedUserPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return listAttachedPolicies(ctx, s, reqCtx, req, userAttachOps)
}

// AttachGroupPolicy attaches a policy to a group.
func (s *IAMService) AttachGroupPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return attachPolicy(ctx, s, reqCtx, req, groupAttachOps)
}

// DetachGroupPolicy detaches a policy from a group.
func (s *IAMService) DetachGroupPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return detachPolicy(ctx, s, reqCtx, req, groupAttachOps)
}

// ListAttachedGroupPolicies lists the policies attached to a group.
func (s *IAMService) ListAttachedGroupPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return listAttachedPolicies(ctx, s, reqCtx, req, groupAttachOps)
}

// AttachRolePolicy attaches a policy to a role.
func (s *IAMService) AttachRolePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return attachPolicy(ctx, s, reqCtx, req, roleAttachOps)
}

// DetachRolePolicy detaches a policy from a role.
func (s *IAMService) DetachRolePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return detachPolicy(ctx, s, reqCtx, req, roleAttachOps)
}

// ListAttachedRolePolicies lists the policies attached to a role.
func (s *IAMService) ListAttachedRolePolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return listAttachedPolicies(ctx, s, reqCtx, req, roleAttachOps)
}
