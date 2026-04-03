package iam

import (
	"context"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/iam"
	"vorpalstacks/internal/pb/aws/iam/iamconnect"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

type AdminHandler struct {
	iamconnect.UnimplementedIAMServiceHandler
	store    storage.BasicStorage
	storeObj *iamstore.IAMStore
}

var _ iamconnect.IAMServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler(store storage.BasicStorage, accountID string) *AdminHandler {
	return &AdminHandler{
		store:    store,
		storeObj: iamstore.NewIAMStore(store, accountID),
	}
}

func (h *AdminHandler) ListUsers(ctx context.Context, req *connect.Request[pb.ListUsersRequest]) (*connect.Response[pb.ListUsersResponse], error) {
	maxItems := int(req.Msg.Maxitems)
	if maxItems <= 0 {
		maxItems = 100
	}

	result, err := h.storeObj.Users().List(req.Msg.Pathprefix, req.Msg.Marker, maxItems)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	users := make([]*pb.User, len(result.Users))
	for i, user := range result.Users {
		users[i] = toPbUser(user)
	}

	return connect.NewResponse(&pb.ListUsersResponse{
		Users:       users,
		Istruncated: result.IsTruncated,
		Marker:      result.Marker,
	}), nil
}

func (h *AdminHandler) ListRoles(ctx context.Context, req *connect.Request[pb.ListRolesRequest]) (*connect.Response[pb.ListRolesResponse], error) {
	maxItems := int(req.Msg.Maxitems)
	if maxItems <= 0 {
		maxItems = 100
	}

	result, err := h.storeObj.Roles().List(req.Msg.Pathprefix, req.Msg.Marker, maxItems)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	roles := make([]*pb.Role, len(result.Roles))
	for i, role := range result.Roles {
		roles[i] = toPbRole(role)
	}

	return connect.NewResponse(&pb.ListRolesResponse{
		Roles:       roles,
		Istruncated: result.IsTruncated,
		Marker:      result.Marker,
	}), nil
}

func (h *AdminHandler) ListPolicies(ctx context.Context, req *connect.Request[pb.ListPoliciesRequest]) (*connect.Response[pb.ListPoliciesResponse], error) {
	maxItems := int(req.Msg.Maxitems)
	if maxItems <= 0 {
		maxItems = 100
	}

	scope := "Local"
	if req.Msg.Scope == pb.PolicyScopeType_POLICY_SCOPE_TYPE_AWS {
		scope = "AWS"
	} else if req.Msg.Scope == pb.PolicyScopeType_POLICY_SCOPE_TYPE_ALL {
		scope = "All"
	}

	result, err := h.storeObj.Policies().List(scope, req.Msg.Pathprefix, req.Msg.Onlyattached, req.Msg.Marker, maxItems)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	policies := make([]*pb.Policy, len(result.Policies))
	for i, policy := range result.Policies {
		policies[i] = toPbPolicy(policy)
	}

	return connect.NewResponse(&pb.ListPoliciesResponse{
		Policies:    policies,
		Istruncated: result.IsTruncated,
		Marker:      result.Marker,
	}), nil
}

func toPbUser(user *iamstore.User) *pb.User {
	pbUser := &pb.User{
		Username:   user.UserName,
		Userid:     user.ID,
		Arn:        user.Arn,
		Path:       user.Path,
		Createdate: user.CreateDate.Format("2006-01-02T15:04:05Z"),
	}

	if user.PasswordLastUsed != nil {
		pbUser.Passwordlastused = user.PasswordLastUsed.Format("2006-01-02T15:04:05Z")
	}

	if user.PermissionsBoundary != nil {
		pbUser.Permissionsboundary = &pb.AttachedPermissionsBoundary{
			Permissionsboundaryarn:  user.PermissionsBoundary.PermissionsBoundaryArn,
			Permissionsboundarytype: pb.PermissionsBoundaryAttachmentType_PERMISSIONS_BOUNDARY_ATTACHMENT_TYPE_POLICY,
		}
	}

	if len(user.Tags) > 0 {
		pbUser.Tags = make([]*pb.Tag, len(user.Tags))
		for i, tag := range user.Tags {
			pbUser.Tags[i] = &pb.Tag{Key: tag.Key, Value: tag.Value}
		}
	}

	return pbUser
}

func toPbRole(role *iamstore.Role) *pb.Role {
	pbRole := &pb.Role{
		Rolename:                 role.RoleName,
		Roleid:                   role.ID,
		Arn:                      role.Arn,
		Path:                     role.Path,
		Createdate:               role.CreateDate.Format("2006-01-02T15:04:05Z"),
		Assumerolepolicydocument: role.AssumeRolePolicyDocument,
		Description:              role.Description,
		Maxsessionduration:       int32(role.MaxSessionDuration),
	}

	if role.PermissionsBoundary != nil {
		pbRole.Permissionsboundary = &pb.AttachedPermissionsBoundary{
			Permissionsboundaryarn:  role.PermissionsBoundary.PermissionsBoundaryArn,
			Permissionsboundarytype: pb.PermissionsBoundaryAttachmentType_PERMISSIONS_BOUNDARY_ATTACHMENT_TYPE_POLICY,
		}
	}

	if role.RoleLastUsed != nil {
		pbRole.Rolelastused = &pb.RoleLastUsed{
			Region: role.RoleLastUsed.Region,
		}
		if role.RoleLastUsed.LastUsedDate != nil {
			pbRole.Rolelastused.Lastuseddate = role.RoleLastUsed.LastUsedDate.Format("2006-01-02T15:04:05Z")
		}
	}

	if len(role.Tags) > 0 {
		pbRole.Tags = make([]*pb.Tag, len(role.Tags))
		for i, tag := range role.Tags {
			pbRole.Tags[i] = &pb.Tag{Key: tag.Key, Value: tag.Value}
		}
	}

	return pbRole
}

func toPbPolicy(policy *iamstore.Policy) *pb.Policy {
	pbPolicy := &pb.Policy{
		Policyname:                    policy.PolicyName,
		Policyid:                      policy.ID,
		Arn:                           policy.Arn,
		Path:                          policy.Path,
		Createdate:                    policy.CreateDate.Format("2006-01-02T15:04:05Z"),
		Updatedate:                    policy.UpdateDate.Format("2006-01-02T15:04:05Z"),
		Defaultversionid:              policy.DefaultVersionId,
		Attachmentcount:               int32(policy.AttachmentCount),
		Isattachable:                  policy.IsAttachable,
		Description:                   policy.Description,
		Permissionsboundaryusagecount: 0,
	}

	if len(policy.Tags) > 0 {
		pbPolicy.Tags = make([]*pb.Tag, len(policy.Tags))
		for i, tag := range policy.Tags {
			pbPolicy.Tags[i] = &pb.Tag{Key: tag.Key, Value: tag.Value}
		}
	}

	return pbPolicy
}
