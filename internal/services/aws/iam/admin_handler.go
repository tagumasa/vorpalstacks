package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/iam"
	"vorpalstacks/internal/pb/aws/iam/iamconnect"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/timeutils"
)

const defaultMaxItemsValue = 100

func pbTagsToStoreTags(pbTags []*pb.Tag) []types.Tag {
	if len(pbTags) == 0 {
		return nil
	}
	tags := make([]types.Tag, len(pbTags))
	for i, t := range pbTags {
		tags[i] = types.Tag{Key: t.Key, Value: t.Value}
	}
	return tags
}

func defaultMaxItems(n int32) int {
	v := int(n)
	if v <= 0 {
		return defaultMaxItemsValue
	}
	return v
}

var _ iamconnect.IAMServiceHandler = (*AdminHandler)(nil)

// AdminHandler implements the IAM admin console gRPC-Web handler.
type AdminHandler struct {
	iamconnect.UnimplementedIAMServiceHandler
	service *IAMService
}

// NewAdminHandler creates a new IAM admin handler.
func NewAdminHandler(svc *IAMService) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) getStore() (*iamstore.IAMStore, error) {
	return h.service.GetStoreForRegion("")
}

func storeErr(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, iamstore.ErrUserNotFound),
		errors.Is(err, iamstore.ErrRoleNotFound),
		errors.Is(err, iamstore.ErrGroupNotFound),
		errors.Is(err, iamstore.ErrPolicyNotFound),
		errors.Is(err, iamstore.ErrAccessKeyNotFound),
		errors.Is(err, iamstore.ErrLoginProfileNotFound),
		errors.Is(err, iamstore.ErrInstanceProfileNotFound),
		errors.Is(err, iamstore.ErrMFADeviceNotFound),
		errors.Is(err, iamstore.ErrPasswordPolicyNotFound),
		errors.Is(err, iamstore.ErrServerCertificateNotFound),
		errors.Is(err, iamstore.ErrSAMLProviderNotFound),
		errors.Is(err, iamstore.ErrOpenIDConnectProviderNotFound),
		errors.Is(err, iamstore.ErrSigningCertificateNotFound),
		errors.Is(err, iamstore.ErrSSHPublicKeyNotFound),
		errors.Is(err, iamstore.ErrServiceSpecificCredentialNotFound):
		return connect.NewError(connect.CodeNotFound, err)
	case errors.Is(err, iamstore.ErrUserAlreadyExists),
		errors.Is(err, iamstore.ErrRoleAlreadyExists),
		errors.Is(err, iamstore.ErrGroupAlreadyExists),
		errors.Is(err, iamstore.ErrPolicyAlreadyExists),
		errors.Is(err, iamstore.ErrLoginProfileExists),
		errors.Is(err, iamstore.ErrInstanceProfileAlreadyExists),
		errors.Is(err, iamstore.ErrRoleAlreadyInInstanceProfile),
		errors.Is(err, iamstore.ErrServerCertificateAlreadyExists),
		errors.Is(err, iamstore.ErrSAMLProviderAlreadyExists),
		errors.Is(err, iamstore.ErrOpenIDConnectProviderAlreadyExists),
		errors.Is(err, iamstore.ErrUserAlreadyInGroup):
		return connect.NewError(connect.CodeAlreadyExists, err)
	case errors.Is(err, iamstore.ErrUserNotInGroup),
		errors.Is(err, iamstore.ErrRoleNotInInstanceProfile):
		return connect.NewError(connect.CodeFailedPrecondition, err)
	case errors.Is(err, iamstore.ErrInvalidPassword),
		errors.Is(err, iamstore.ErrInvalidAccessKeyStatus):
		return connect.NewError(connect.CodeInvalidArgument, err)
	default:
		return connect.NewError(connect.CodeInternal, err)
	}
}

// --- User operations ---

// GetUser returns a single IAM user by name.
func (h *AdminHandler) GetUser(ctx context.Context, req *connect.Request[pb.GetUserRequest]) (*connect.Response[pb.GetUserResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	if req.Msg.Username == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("UserName is required"))
	}
	user, err := stores.Users().Get(req.Msg.Username)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(&pb.GetUserResponse{User: toPbUser(user)}), nil
}

// ListUsers returns a paginated list of IAM users via the admin console gRPC-Web interface.
func (h *AdminHandler) ListUsers(ctx context.Context, req *connect.Request[pb.ListUsersRequest]) (*connect.Response[pb.ListUsersResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	maxItems := defaultMaxItems(req.Msg.Maxitems)

	result, err := stores.Users().List(req.Msg.Pathprefix, req.Msg.Marker, maxItems)
	if err != nil {
		return nil, storeErr(err)
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

// CreateUser creates a new IAM user via the admin console gRPC-Web interface.
func (h *AdminHandler) CreateUser(ctx context.Context, req *connect.Request[pb.CreateUserRequest]) (*connect.Response[pb.CreateUserResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	if req.Msg.Username == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("UserName is required"))
	}

	tags := pbTagsToStoreTags(req.Msg.Tags)

	user, err := stores.Users().Create(req.Msg.Username, req.Msg.Path, stores.AccountID(), tags)
	if err != nil {
		return nil, storeErr(err)
	}

	return connect.NewResponse(&pb.CreateUserResponse{
		User: toPbUser(user),
	}), nil
}

// UpdateUser updates an existing IAM user.
func (h *AdminHandler) UpdateUser(ctx context.Context, req *connect.Request[pb.UpdateUserRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	if req.Msg.Username == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("UserName is required"))
	}

	if err := stores.RenameUser(req.Msg.Username, req.Msg.Newusername, req.Msg.Newpath); err != nil {
		return nil, storeErr(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// DeleteUser deletes an IAM user via the admin console gRPC-Web interface.
func (h *AdminHandler) DeleteUser(ctx context.Context, req *connect.Request[pb.DeleteUserRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	if req.Msg.Username == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("UserName is required"))
	}

	if err := stores.Users().Delete(req.Msg.Username); err != nil {
		return nil, storeErr(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// --- Role operations ---

// GetRole returns a single IAM role by name.
func (h *AdminHandler) GetRole(ctx context.Context, req *connect.Request[pb.GetRoleRequest]) (*connect.Response[pb.GetRoleResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	if req.Msg.Rolename == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("RoleName is required"))
	}
	role, err := stores.Roles().Get(req.Msg.Rolename)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(&pb.GetRoleResponse{Role: toPbRole(role)}), nil
}

// ListRoles returns a paginated list of IAM roles via the admin console gRPC-Web interface.
func (h *AdminHandler) ListRoles(ctx context.Context, req *connect.Request[pb.ListRolesRequest]) (*connect.Response[pb.ListRolesResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	maxItems := defaultMaxItems(req.Msg.Maxitems)

	result, err := stores.Roles().List(req.Msg.Pathprefix, req.Msg.Marker, maxItems)
	if err != nil {
		return nil, storeErr(err)
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

// CreateRole creates a new IAM role via the admin console gRPC-Web interface.
func (h *AdminHandler) CreateRole(ctx context.Context, req *connect.Request[pb.CreateRoleRequest]) (*connect.Response[pb.CreateRoleResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	if req.Msg.Rolename == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("RoleName is required"))
	}
	if req.Msg.Assumerolepolicydocument == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("AssumeRolePolicyDocument is required"))
	}

	tags := pbTagsToStoreTags(req.Msg.Tags)

	maxSessionDuration := int(req.Msg.Maxsessionduration)
	if maxSessionDuration == 0 {
		maxSessionDuration = 3600
	}

	role, err := stores.Roles().Create(
		req.Msg.Rolename,
		req.Msg.Path,
		stores.AccountID(),
		req.Msg.Assumerolepolicydocument,
		req.Msg.Description,
		maxSessionDuration,
		tags,
	)
	if err != nil {
		return nil, storeErr(err)
	}

	return connect.NewResponse(&pb.CreateRoleResponse{
		Role: toPbRole(role),
	}), nil
}

// UpdateRole updates an existing IAM role.
func (h *AdminHandler) UpdateRole(ctx context.Context, req *connect.Request[pb.UpdateRoleRequest]) (*connect.Response[pb.UpdateRoleResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	if req.Msg.Rolename == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("RoleName is required"))
	}

	if err := stores.UpdateRoleFields(req.Msg.Rolename, req.Msg.Description, int(req.Msg.Maxsessionduration)); err != nil {
		return nil, storeErr(err)
	}

	return connect.NewResponse(&pb.UpdateRoleResponse{}), nil
}

// DeleteRole deletes an IAM role via the admin console gRPC-Web interface.
func (h *AdminHandler) DeleteRole(ctx context.Context, req *connect.Request[pb.DeleteRoleRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	if req.Msg.Rolename == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("RoleName is required"))
	}

	if err := stores.Roles().Delete(req.Msg.Rolename); err != nil {
		return nil, storeErr(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// --- Policy operations ---

// GetPolicy returns a single IAM policy by ARN.
func (h *AdminHandler) GetPolicy(ctx context.Context, req *connect.Request[pb.GetPolicyRequest]) (*connect.Response[pb.GetPolicyResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	if req.Msg.Policyarn == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("PolicyArn is required"))
	}
	policy, err := stores.Policies().Get(req.Msg.Policyarn)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(&pb.GetPolicyResponse{Policy: toPbPolicy(policy)}), nil
}

// ListPolicies returns a paginated list of IAM policies via the admin console gRPC-Web interface.
func (h *AdminHandler) ListPolicies(ctx context.Context, req *connect.Request[pb.ListPoliciesRequest]) (*connect.Response[pb.ListPoliciesResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	maxItems := defaultMaxItems(req.Msg.Maxitems)

	scope := "Local"
	if req.Msg.Scope == pb.PolicyScopeType_POLICY_SCOPE_TYPE_AWS {
		scope = "AWS"
	} else if req.Msg.Scope == pb.PolicyScopeType_POLICY_SCOPE_TYPE_ALL {
		scope = "All"
	}

	result, err := stores.Policies().List(scope, req.Msg.Pathprefix, req.Msg.Onlyattached, req.Msg.Marker, maxItems)
	if err != nil {
		return nil, storeErr(err)
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

// CreatePolicy creates a new IAM managed policy.
func (h *AdminHandler) CreatePolicy(ctx context.Context, req *connect.Request[pb.CreatePolicyRequest]) (*connect.Response[pb.CreatePolicyResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	if req.Msg.Policyname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("PolicyName is required"))
	}
	if req.Msg.Policydocument == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("PolicyDocument is required"))
	}

	tags := pbTagsToStoreTags(req.Msg.Tags)

	policy, err := stores.Policies().Create(req.Msg.Policyname, req.Msg.Path, stores.AccountID(), req.Msg.Policydocument, req.Msg.Description, tags)
	if err != nil {
		return nil, storeErr(err)
	}

	return connect.NewResponse(&pb.CreatePolicyResponse{Policy: toPbPolicy(policy)}), nil
}

// DeletePolicy deletes an IAM managed policy.
func (h *AdminHandler) DeletePolicy(ctx context.Context, req *connect.Request[pb.DeletePolicyRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	if req.Msg.Policyarn == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("PolicyArn is required"))
	}

	if err := stores.Policies().Delete(req.Msg.Policyarn); err != nil {
		return nil, storeErr(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// --- Group operations ---

// GetGroup returns a single IAM group by name.
func (h *AdminHandler) GetGroup(ctx context.Context, req *connect.Request[pb.GetGroupRequest]) (*connect.Response[pb.GetGroupResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	if req.Msg.Groupname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("GroupName is required"))
	}
	group, err := stores.Groups().Get(req.Msg.Groupname)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(&pb.GetGroupResponse{Group: toPbGroup(group)}), nil
}

// ListGroups returns a paginated list of IAM groups.
func (h *AdminHandler) ListGroups(ctx context.Context, req *connect.Request[pb.ListGroupsRequest]) (*connect.Response[pb.ListGroupsResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	maxItems := defaultMaxItems(req.Msg.Maxitems)

	result, err := stores.Groups().List(req.Msg.Pathprefix, req.Msg.Marker, maxItems)
	if err != nil {
		return nil, storeErr(err)
	}

	groups := make([]*pb.Group, len(result.Groups))
	for i, group := range result.Groups {
		groups[i] = toPbGroup(group)
	}

	return connect.NewResponse(&pb.ListGroupsResponse{
		Groups:      groups,
		Istruncated: result.IsTruncated,
		Marker:      result.Marker,
	}), nil
}

// CreateGroup creates a new IAM group.
func (h *AdminHandler) CreateGroup(ctx context.Context, req *connect.Request[pb.CreateGroupRequest]) (*connect.Response[pb.CreateGroupResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	if req.Msg.Groupname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("GroupName is required"))
	}

	group, err := stores.Groups().Create(req.Msg.Groupname, req.Msg.Path, stores.AccountID())
	if err != nil {
		return nil, storeErr(err)
	}

	return connect.NewResponse(&pb.CreateGroupResponse{Group: toPbGroup(group)}), nil
}

// UpdateGroup updates an existing IAM group.
func (h *AdminHandler) UpdateGroup(ctx context.Context, req *connect.Request[pb.UpdateGroupRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	if req.Msg.Groupname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("GroupName is required"))
	}

	if err := stores.RenameGroup(req.Msg.Groupname, req.Msg.Newgroupname, req.Msg.Newpath); err != nil {
		return nil, storeErr(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// DeleteGroup deletes an IAM group.
func (h *AdminHandler) DeleteGroup(ctx context.Context, req *connect.Request[pb.DeleteGroupRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, storeErr(err)
	}
	if req.Msg.Groupname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("GroupName is required"))
	}

	if err := stores.Groups().Delete(req.Msg.Groupname); err != nil {
		return nil, storeErr(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// --- Convert functions ---

func toPbUser(user *iamstore.User) *pb.User {
	pbUser := &pb.User{
		Username:   user.UserName,
		Userid:     user.ID,
		Arn:        user.Arn,
		Path:       user.Path,
		Createdate: user.CreateDate.Format(timeutils.ISO8601UTCFormat),
	}

	if user.PasswordLastUsed != nil {
		pbUser.Passwordlastused = user.PasswordLastUsed.Format(timeutils.ISO8601UTCFormat)
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
		Createdate:               role.CreateDate.Format(timeutils.ISO8601UTCFormat),
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
			pbRole.Rolelastused.Lastuseddate = role.RoleLastUsed.LastUsedDate.Format(timeutils.ISO8601UTCFormat)
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
		Policyname:       policy.PolicyName,
		Policyid:         policy.ID,
		Arn:              policy.Arn,
		Path:             policy.Path,
		Createdate:       policy.CreateDate.Format(timeutils.ISO8601UTCFormat),
		Updatedate:       policy.UpdateDate.Format(timeutils.ISO8601UTCFormat),
		Defaultversionid: policy.DefaultVersionId,
		Attachmentcount:  int32(policy.AttachmentCount),
		Isattachable:     policy.IsAttachable,
		Description:      policy.Description,
	}

	if len(policy.Tags) > 0 {
		pbPolicy.Tags = make([]*pb.Tag, len(policy.Tags))
		for i, tag := range policy.Tags {
			pbPolicy.Tags[i] = &pb.Tag{Key: tag.Key, Value: tag.Value}
		}
	}

	return pbPolicy
}

func toPbGroup(group *iamstore.Group) *pb.Group {
	pbGroup := &pb.Group{
		Groupid:    group.ID,
		Groupname:  group.GroupName,
		Arn:        group.Arn,
		Path:       group.Path,
		Createdate: group.CreateDate.Format(timeutils.ISO8601UTCFormat),
	}

	return pbGroup
}

// NewConnectHandler creates a gRPC-Web connect handler for the Iam admin console.
func NewConnectHandler(svc *IAMService) (string, http.Handler) {
	return iamconnect.NewIAMServiceHandler(NewAdminHandler(svc))
}
