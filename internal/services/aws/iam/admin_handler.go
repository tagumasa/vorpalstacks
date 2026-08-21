package iam

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	svcerrors "vorpalstacks/internal/common/errors"
	types "vorpalstacks/internal/common/tags"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/iam"
	"vorpalstacks/internal/pb/aws/iam/iamconnect"
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

// --- User operations ---

// GetUser returns a single IAM user by name.
func (h *AdminHandler) GetUser(ctx context.Context, req *connect.Request[pb.GetUserRequest]) (*connect.Response[pb.GetUserResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	if req.Msg.Username == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("UserName is required"))
	}
	user, err := h.service.getUserCore(stores, req.Msg.Username)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pb.GetUserResponse{User: toPbUser(user)}), nil
}

// ListUsers returns a paginated list of IAM users via the admin console gRPC-Web interface.
func (h *AdminHandler) ListUsers(ctx context.Context, req *connect.Request[pb.ListUsersRequest]) (*connect.Response[pb.ListUsersResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	maxItems := defaultMaxItems(req.Msg.GetMaxitems())

	result, err := h.service.listUsersCore(stores, req.Msg.Pathprefix, req.Msg.Marker, maxItems)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	users := make([]*pb.User, len(result.Users))
	for i, user := range result.Users {
		users[i] = toPbUser(user)
	}

	return connect.NewResponse(&pb.ListUsersResponse{
		Users:       users,
		Istruncated: proto.Bool(result.IsTruncated),
		Marker:      result.Marker,
	}), nil
}

// CreateUser creates a new IAM user via the admin console gRPC-Web interface.
func (h *AdminHandler) CreateUser(ctx context.Context, req *connect.Request[pb.CreateUserRequest]) (*connect.Response[pb.CreateUserResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	input := &CreateUserInput{
		UserName: req.Msg.Username,
		Path:     req.Msg.Path,
		Tags:     pbTagsToStoreTags(req.Msg.Tags),
	}
	user, err := h.service.createUserCore(stores, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateUserResponse{
		User: toPbUser(user),
	}), nil
}

// UpdateUser updates an existing IAM user.
func (h *AdminHandler) UpdateUser(ctx context.Context, req *connect.Request[pb.UpdateUserRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	input := &UpdateUserInput{
		UserName:    req.Msg.Username,
		NewPath:     req.Msg.Newpath,
		NewUserName: req.Msg.Newusername,
	}
	if _, err := h.service.updateUserCore(stores, input); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// DeleteUser deletes an IAM user via the admin console gRPC-Web interface.
func (h *AdminHandler) DeleteUser(ctx context.Context, req *connect.Request[pb.DeleteUserRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	input := &DeleteUserInput{
		UserName: req.Msg.Username,
		Cascade:  true,
	}
	if err := h.service.deleteUserCore(stores, input); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// --- Role operations ---

// GetRole returns a single IAM role by name.
func (h *AdminHandler) GetRole(ctx context.Context, req *connect.Request[pb.GetRoleRequest]) (*connect.Response[pb.GetRoleResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	if req.Msg.Rolename == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("RoleName is required"))
	}
	role, err := h.service.getRoleCore(stores, req.Msg.Rolename)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pb.GetRoleResponse{Role: toPbRole(role)}), nil
}

// ListRoles returns a paginated list of IAM roles via the admin console gRPC-Web interface.
func (h *AdminHandler) ListRoles(ctx context.Context, req *connect.Request[pb.ListRolesRequest]) (*connect.Response[pb.ListRolesResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	maxItems := defaultMaxItems(req.Msg.GetMaxitems())

	result, err := h.service.listRolesCore(stores, req.Msg.Pathprefix, req.Msg.Marker, maxItems)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	roles := make([]*pb.Role, len(result.Roles))
	for i, role := range result.Roles {
		roles[i] = toPbRole(role)
	}

	return connect.NewResponse(&pb.ListRolesResponse{
		Roles:       roles,
		Istruncated: proto.Bool(result.IsTruncated),
		Marker:      result.Marker,
	}), nil
}

// CreateRole creates a new IAM role via the admin console gRPC-Web interface.
func (h *AdminHandler) CreateRole(ctx context.Context, req *connect.Request[pb.CreateRoleRequest]) (*connect.Response[pb.CreateRoleResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	input := &CreateRoleInput{
		RoleName:                 req.Msg.Rolename,
		Path:                     req.Msg.Path,
		AssumeRolePolicyDocument: req.Msg.Assumerolepolicydocument,
		Description:              req.Msg.Description,
		MaxSessionDuration:       int(req.Msg.GetMaxsessionduration()),
		Tags:                     pbTagsToStoreTags(req.Msg.Tags),
	}
	role, err := h.service.createRoleCore(stores, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateRoleResponse{
		Role: toPbRole(role),
	}), nil
}

// UpdateRole updates an existing IAM role.
func (h *AdminHandler) UpdateRole(ctx context.Context, req *connect.Request[pb.UpdateRoleRequest]) (*connect.Response[pb.UpdateRoleResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	input := &UpdateRoleInput{
		RoleName:           req.Msg.Rolename,
		Description:        req.Msg.Description,
		MaxSessionDuration: int(req.Msg.GetMaxsessionduration()),
	}
	if _, err := h.service.updateRoleCore(stores, input); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.UpdateRoleResponse{}), nil
}

// DeleteRole deletes an IAM role via the admin console gRPC-Web interface.
func (h *AdminHandler) DeleteRole(ctx context.Context, req *connect.Request[pb.DeleteRoleRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	input := &DeleteRoleInput{
		RoleName: req.Msg.Rolename,
		Cascade:  true,
	}
	if err := h.service.deleteRoleCore(stores, input); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// --- Policy operations ---

// GetPolicy returns a single IAM policy by ARN.
func (h *AdminHandler) GetPolicy(ctx context.Context, req *connect.Request[pb.GetPolicyRequest]) (*connect.Response[pb.GetPolicyResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	if req.Msg.Policyarn == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("PolicyArn is required"))
	}
	policy, err := h.service.getPolicyCore(stores, req.Msg.Policyarn)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pb.GetPolicyResponse{Policy: toPbPolicy(policy)}), nil
}

// ListPolicies returns a paginated list of IAM policies via the admin console gRPC-Web interface.
func (h *AdminHandler) ListPolicies(ctx context.Context, req *connect.Request[pb.ListPoliciesRequest]) (*connect.Response[pb.ListPoliciesResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	maxItems := defaultMaxItems(req.Msg.GetMaxitems())

	scope := "Local"
	if req.Msg.Scope == pb.PolicyScopeType_POLICY_SCOPE_TYPE_AWS {
		scope = "AWS"
	} else if req.Msg.Scope == pb.PolicyScopeType_POLICY_SCOPE_TYPE_ALL {
		scope = "All"
	}

	result, err := h.service.listPoliciesCore(stores, scope, req.Msg.Pathprefix, req.Msg.Marker, req.Msg.GetOnlyattached(), maxItems)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	policies := make([]*pb.Policy, len(result.Policies))
	for i, policy := range result.Policies {
		policies[i] = toPbPolicy(policy)
	}

	return connect.NewResponse(&pb.ListPoliciesResponse{
		Policies:    policies,
		Istruncated: proto.Bool(result.IsTruncated),
		Marker:      result.Marker,
	}), nil
}

// CreatePolicy creates a new IAM managed policy.
func (h *AdminHandler) CreatePolicy(ctx context.Context, req *connect.Request[pb.CreatePolicyRequest]) (*connect.Response[pb.CreatePolicyResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	input := &CreatePolicyInput{
		PolicyName:     req.Msg.Policyname,
		Path:           req.Msg.Path,
		PolicyDocument: req.Msg.Policydocument,
		Description:    req.Msg.Description,
		Tags:           pbTagsToStoreTags(req.Msg.Tags),
	}
	policy, err := h.service.createPolicyCore(stores, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreatePolicyResponse{Policy: toPbPolicy(policy)}), nil
}

// DeletePolicy deletes an IAM managed policy.
func (h *AdminHandler) DeletePolicy(ctx context.Context, req *connect.Request[pb.DeletePolicyRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	input := &DeletePolicyInput{
		PolicyArn: req.Msg.Policyarn,
	}
	if err := h.service.deletePolicyCore(stores, input); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// --- Group operations ---

// GetGroup returns a single IAM group by name.
func (h *AdminHandler) GetGroup(ctx context.Context, req *connect.Request[pb.GetGroupRequest]) (*connect.Response[pb.GetGroupResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	if req.Msg.Groupname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("GroupName is required"))
	}
	group, err := h.service.getGroupCore(stores, req.Msg.Groupname)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pb.GetGroupResponse{Group: toPbGroup(group)}), nil
}

// ListGroups returns a paginated list of IAM groups.
func (h *AdminHandler) ListGroups(ctx context.Context, req *connect.Request[pb.ListGroupsRequest]) (*connect.Response[pb.ListGroupsResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	maxItems := defaultMaxItems(req.Msg.GetMaxitems())

	result, err := h.service.listGroupsCore(stores, req.Msg.Pathprefix, req.Msg.Marker, maxItems)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	groups := make([]*pb.Group, len(result.Groups))
	for i, group := range result.Groups {
		groups[i] = toPbGroup(group)
	}

	return connect.NewResponse(&pb.ListGroupsResponse{
		Groups:      groups,
		Istruncated: proto.Bool(result.IsTruncated),
		Marker:      result.Marker,
	}), nil
}

// CreateGroup creates a new IAM group.
func (h *AdminHandler) CreateGroup(ctx context.Context, req *connect.Request[pb.CreateGroupRequest]) (*connect.Response[pb.CreateGroupResponse], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	input := &CreateGroupInput{
		GroupName: req.Msg.Groupname,
		Path:      req.Msg.Path,
	}
	group, err := h.service.createGroupCore(stores, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateGroupResponse{Group: toPbGroup(group)}), nil
}

// UpdateGroup updates an existing IAM group.
func (h *AdminHandler) UpdateGroup(ctx context.Context, req *connect.Request[pb.UpdateGroupRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	input := &UpdateGroupInput{
		GroupName:    req.Msg.Groupname,
		NewPath:      req.Msg.Newpath,
		NewGroupName: req.Msg.Newgroupname,
	}
	if _, err := h.service.updateGroupCore(stores, input); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// DeleteGroup deletes an IAM group.
func (h *AdminHandler) DeleteGroup(ctx context.Context, req *connect.Request[pb.DeleteGroupRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStore()
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	input := &DeleteGroupInput{
		GroupName: req.Msg.Groupname,
		Cascade:   true,
	}
	if err := h.service.deleteGroupCore(stores, input); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Iam admin console.
func NewConnectHandler(svc *IAMService) (string, http.Handler) {
	return iamconnect.NewIAMServiceHandler(NewAdminHandler(svc))
}
