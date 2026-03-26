package iam

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/iam"
	"vorpalstacks/internal/pb/aws/iam/iamconnect"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// AdminHandler provides IAM service administration functionality.
// It implements the IAMServiceHandler interface for gRPC-Web communication.
type AdminHandler struct {
	store    storage.BasicStorage
	storeObj *iamstore.IAMStore
}

var _ iamconnect.IAMServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new IAM AdminHandler.
// It initialises the handler with the provided storage and account ID.
func NewAdminHandler(store storage.BasicStorage, accountID string) *AdminHandler {
	return &AdminHandler{
		store:    store,
		storeObj: iamstore.NewIAMStore(store, accountID),
	}
}

// ListUsers lists IAM users.
// It returns users matching the specified path prefix with pagination support.
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

// ListRoles lists IAM roles.
// It returns roles matching the specified path prefix with pagination support.
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

// ListPolicies lists IAM policies.
// It returns policies matching the specified scope and path prefix with pagination support.
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

// AcceptDelegationRequest accepts a delegation request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AcceptDelegationRequest(ctx context.Context, req *connect.Request[pb.AcceptDelegationRequestRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// AddClientIDToOpenIDConnectProvider adds a client ID to an OpenID Connect provider.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AddClientIDToOpenIDConnectProvider(ctx context.Context, req *connect.Request[pb.AddClientIDToOpenIDConnectProviderRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// AddRoleToInstanceProfile adds a role to an instance profile.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AddRoleToInstanceProfile(ctx context.Context, req *connect.Request[pb.AddRoleToInstanceProfileRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// AddUserToGroup adds a user to a group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AddUserToGroup(ctx context.Context, req *connect.Request[pb.AddUserToGroupRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// AssociateDelegationRequest associates a delegation request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AssociateDelegationRequest(ctx context.Context, req *connect.Request[pb.AssociateDelegationRequestRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// AttachGroupPolicy attaches a policy to a group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AttachGroupPolicy(ctx context.Context, req *connect.Request[pb.AttachGroupPolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// AttachRolePolicy attaches a policy to a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AttachRolePolicy(ctx context.Context, req *connect.Request[pb.AttachRolePolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// AttachUserPolicy attaches a policy to a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AttachUserPolicy(ctx context.Context, req *connect.Request[pb.AttachUserPolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ChangePassword changes the password for a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ChangePassword(ctx context.Context, req *connect.Request[pb.ChangePasswordRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateAccessKey creates an access key for a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateAccessKey(ctx context.Context, req *connect.Request[pb.CreateAccessKeyRequest]) (*connect.Response[pb.CreateAccessKeyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateAccountAlias creates an account alias.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateAccountAlias(ctx context.Context, req *connect.Request[pb.CreateAccountAliasRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateDelegationRequest creates a delegation request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateDelegationRequest(ctx context.Context, req *connect.Request[pb.CreateDelegationRequestRequest]) (*connect.Response[pb.CreateDelegationRequestResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateGroup creates a new group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateGroup(ctx context.Context, req *connect.Request[pb.CreateGroupRequest]) (*connect.Response[pb.CreateGroupResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateInstanceProfile creates an instance profile.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateInstanceProfile(ctx context.Context, req *connect.Request[pb.CreateInstanceProfileRequest]) (*connect.Response[pb.CreateInstanceProfileResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateLoginProfile creates a login profile for a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateLoginProfile(ctx context.Context, req *connect.Request[pb.CreateLoginProfileRequest]) (*connect.Response[pb.CreateLoginProfileResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateOpenIDConnectProvider creates an OpenID Connect provider.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateOpenIDConnectProvider(ctx context.Context, req *connect.Request[pb.CreateOpenIDConnectProviderRequest]) (*connect.Response[pb.CreateOpenIDConnectProviderResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreatePolicy creates a new managed policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreatePolicy(ctx context.Context, req *connect.Request[pb.CreatePolicyRequest]) (*connect.Response[pb.CreatePolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreatePolicyVersion creates a new version of a managed policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreatePolicyVersion(ctx context.Context, req *connect.Request[pb.CreatePolicyVersionRequest]) (*connect.Response[pb.CreatePolicyVersionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateRole creates a new role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateRole(ctx context.Context, req *connect.Request[pb.CreateRoleRequest]) (*connect.Response[pb.CreateRoleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateSAMLProvider creates a SAML provider.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateSAMLProvider(ctx context.Context, req *connect.Request[pb.CreateSAMLProviderRequest]) (*connect.Response[pb.CreateSAMLProviderResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateServiceLinkedRole creates a service-linked role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateServiceLinkedRole(ctx context.Context, req *connect.Request[pb.CreateServiceLinkedRoleRequest]) (*connect.Response[pb.CreateServiceLinkedRoleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateServiceSpecificCredential creates a service-specific credential.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateServiceSpecificCredential(ctx context.Context, req *connect.Request[pb.CreateServiceSpecificCredentialRequest]) (*connect.Response[pb.CreateServiceSpecificCredentialResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateUser creates a new IAM user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateUser(ctx context.Context, req *connect.Request[pb.CreateUserRequest]) (*connect.Response[pb.CreateUserResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateVirtualMFADevice creates a virtual MFA device.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateVirtualMFADevice(ctx context.Context, req *connect.Request[pb.CreateVirtualMFADeviceRequest]) (*connect.Response[pb.CreateVirtualMFADeviceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeactivateMFADevice deactivates an MFA device.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeactivateMFADevice(ctx context.Context, req *connect.Request[pb.DeactivateMFADeviceRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteAccessKey deletes an access key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteAccessKey(ctx context.Context, req *connect.Request[pb.DeleteAccessKeyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteAccountAlias deletes an account alias.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteAccountAlias(ctx context.Context, req *connect.Request[pb.DeleteAccountAliasRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteAccountPasswordPolicy deletes the account password policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteAccountPasswordPolicy(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteGroup deletes a group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteGroup(ctx context.Context, req *connect.Request[pb.DeleteGroupRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteGroupPolicy deletes a policy from a group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteGroupPolicy(ctx context.Context, req *connect.Request[pb.DeleteGroupPolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteInstanceProfile deletes an instance profile.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteInstanceProfile(ctx context.Context, req *connect.Request[pb.DeleteInstanceProfileRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteLoginProfile deletes a login profile for a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteLoginProfile(ctx context.Context, req *connect.Request[pb.DeleteLoginProfileRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteOpenIDConnectProvider deletes an OpenID Connect provider.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteOpenIDConnectProvider(ctx context.Context, req *connect.Request[pb.DeleteOpenIDConnectProviderRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeletePolicy deletes a managed policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeletePolicy(ctx context.Context, req *connect.Request[pb.DeletePolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeletePolicyVersion deletes a version of a managed policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeletePolicyVersion(ctx context.Context, req *connect.Request[pb.DeletePolicyVersionRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteRole deletes a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteRole(ctx context.Context, req *connect.Request[pb.DeleteRoleRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteRolePermissionsBoundary deletes the permissions boundary for a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteRolePermissionsBoundary(ctx context.Context, req *connect.Request[pb.DeleteRolePermissionsBoundaryRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteRolePolicy deletes a policy from a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteRolePolicy(ctx context.Context, req *connect.Request[pb.DeleteRolePolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteSAMLProvider deletes a SAML provider.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteSAMLProvider(ctx context.Context, req *connect.Request[pb.DeleteSAMLProviderRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteSSHPublicKey deletes an SSH public key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteSSHPublicKey(ctx context.Context, req *connect.Request[pb.DeleteSSHPublicKeyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteServerCertificate deletes a server certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteServerCertificate(ctx context.Context, req *connect.Request[pb.DeleteServerCertificateRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteServiceLinkedRole deletes a service-linked role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteServiceLinkedRole(ctx context.Context, req *connect.Request[pb.DeleteServiceLinkedRoleRequest]) (*connect.Response[pb.DeleteServiceLinkedRoleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteServiceSpecificCredential deletes a service-specific credential.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteServiceSpecificCredential(ctx context.Context, req *connect.Request[pb.DeleteServiceSpecificCredentialRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteSigningCertificate deletes a signing certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteSigningCertificate(ctx context.Context, req *connect.Request[pb.DeleteSigningCertificateRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteUser deletes a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteUser(ctx context.Context, req *connect.Request[pb.DeleteUserRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteUserPermissionsBoundary deletes the permissions boundary for a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteUserPermissionsBoundary(ctx context.Context, req *connect.Request[pb.DeleteUserPermissionsBoundaryRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteUserPolicy deletes a policy from a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteUserPolicy(ctx context.Context, req *connect.Request[pb.DeleteUserPolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteVirtualMFADevice deletes a virtual MFA device.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteVirtualMFADevice(ctx context.Context, req *connect.Request[pb.DeleteVirtualMFADeviceRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DetachGroupPolicy detaches a policy from a group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DetachGroupPolicy(ctx context.Context, req *connect.Request[pb.DetachGroupPolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DetachRolePolicy detaches a policy from a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DetachRolePolicy(ctx context.Context, req *connect.Request[pb.DetachRolePolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DetachUserPolicy detaches a policy from a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DetachUserPolicy(ctx context.Context, req *connect.Request[pb.DetachUserPolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DisableOrganizationsRootCredentialsManagement disables root credentials management.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DisableOrganizationsRootCredentialsManagement(ctx context.Context, req *connect.Request[pb.DisableOrganizationsRootCredentialsManagementRequest]) (*connect.Response[pb.DisableOrganizationsRootCredentialsManagementResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DisableOrganizationsRootSessions disables root sessions.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DisableOrganizationsRootSessions(ctx context.Context, req *connect.Request[pb.DisableOrganizationsRootSessionsRequest]) (*connect.Response[pb.DisableOrganizationsRootSessionsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DisableOutboundWebIdentityFederation disables outbound web identity federation.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DisableOutboundWebIdentityFederation(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// EnableMFADevice enables an MFA device for a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) EnableMFADevice(ctx context.Context, req *connect.Request[pb.EnableMFADeviceRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// EnableOrganizationsRootCredentialsManagement enables root credentials management.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) EnableOrganizationsRootCredentialsManagement(ctx context.Context, req *connect.Request[pb.EnableOrganizationsRootCredentialsManagementRequest]) (*connect.Response[pb.EnableOrganizationsRootCredentialsManagementResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// EnableOrganizationsRootSessions enables root sessions.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) EnableOrganizationsRootSessions(ctx context.Context, req *connect.Request[pb.EnableOrganizationsRootSessionsRequest]) (*connect.Response[pb.EnableOrganizationsRootSessionsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// EnableOutboundWebIdentityFederation enables outbound web identity federation.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) EnableOutboundWebIdentityFederation(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[pb.EnableOutboundWebIdentityFederationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GenerateCredentialReport generates a credential report.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GenerateCredentialReport(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[pb.GenerateCredentialReportResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GenerateOrganizationsAccessReport generates an organisations access report.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GenerateOrganizationsAccessReport(ctx context.Context, req *connect.Request[pb.GenerateOrganizationsAccessReportRequest]) (*connect.Response[pb.GenerateOrganizationsAccessReportResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GenerateServiceLastAccessedDetails generates service last accessed details.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GenerateServiceLastAccessedDetails(ctx context.Context, req *connect.Request[pb.GenerateServiceLastAccessedDetailsRequest]) (*connect.Response[pb.GenerateServiceLastAccessedDetailsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetAccessKeyLastUsed returns information about when an access key was last used.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetAccessKeyLastUsed(ctx context.Context, req *connect.Request[pb.GetAccessKeyLastUsedRequest]) (*connect.Response[pb.GetAccessKeyLastUsedResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetAccountAuthorizationDetails returns account authorisation details.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetAccountAuthorizationDetails(ctx context.Context, req *connect.Request[pb.GetAccountAuthorizationDetailsRequest]) (*connect.Response[pb.GetAccountAuthorizationDetailsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetAccountPasswordPolicy returns the account password policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetAccountPasswordPolicy(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[pb.GetAccountPasswordPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetAccountSummary returns account summary information.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetAccountSummary(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[pb.GetAccountSummaryResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetContextKeysForCustomPolicy returns context keys for a custom policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetContextKeysForCustomPolicy(ctx context.Context, req *connect.Request[pb.GetContextKeysForCustomPolicyRequest]) (*connect.Response[pb.GetContextKeysForPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetContextKeysForPrincipalPolicy returns context keys for a principal policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetContextKeysForPrincipalPolicy(ctx context.Context, req *connect.Request[pb.GetContextKeysForPrincipalPolicyRequest]) (*connect.Response[pb.GetContextKeysForPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetCredentialReport returns a credential report.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetCredentialReport(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[pb.GetCredentialReportResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetDelegationRequest returns information about a delegation request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDelegationRequest(ctx context.Context, req *connect.Request[pb.GetDelegationRequestRequest]) (*connect.Response[pb.GetDelegationRequestResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetGroup returns information about a group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetGroup(ctx context.Context, req *connect.Request[pb.GetGroupRequest]) (*connect.Response[pb.GetGroupResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetGroupPolicy returns information about a policy attached to a group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetGroupPolicy(ctx context.Context, req *connect.Request[pb.GetGroupPolicyRequest]) (*connect.Response[pb.GetGroupPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetHumanReadableSummary returns a human-readable summary of service last accessed details.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetHumanReadableSummary(ctx context.Context, req *connect.Request[pb.GetHumanReadableSummaryRequest]) (*connect.Response[pb.GetHumanReadableSummaryResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetInstanceProfile returns information about an instance profile.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetInstanceProfile(ctx context.Context, req *connect.Request[pb.GetInstanceProfileRequest]) (*connect.Response[pb.GetInstanceProfileResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetLoginProfile returns information about a login profile for a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetLoginProfile(ctx context.Context, req *connect.Request[pb.GetLoginProfileRequest]) (*connect.Response[pb.GetLoginProfileResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetMFADevice returns information about an MFA device.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetMFADevice(ctx context.Context, req *connect.Request[pb.GetMFADeviceRequest]) (*connect.Response[pb.GetMFADeviceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetOpenIDConnectProvider returns information about an OpenID Connect provider.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetOpenIDConnectProvider(ctx context.Context, req *connect.Request[pb.GetOpenIDConnectProviderRequest]) (*connect.Response[pb.GetOpenIDConnectProviderResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetOrganizationsAccessReport returns an organisations access report.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetOrganizationsAccessReport(ctx context.Context, req *connect.Request[pb.GetOrganizationsAccessReportRequest]) (*connect.Response[pb.GetOrganizationsAccessReportResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetOutboundWebIdentityFederationInfo returns information about outbound web identity federation.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetOutboundWebIdentityFederationInfo(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[pb.GetOutboundWebIdentityFederationInfoResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetPolicy returns information about a managed policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetPolicy(ctx context.Context, req *connect.Request[pb.GetPolicyRequest]) (*connect.Response[pb.GetPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetPolicyVersion returns information about a version of a managed policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetPolicyVersion(ctx context.Context, req *connect.Request[pb.GetPolicyVersionRequest]) (*connect.Response[pb.GetPolicyVersionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetRole returns information about a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRole(ctx context.Context, req *connect.Request[pb.GetRoleRequest]) (*connect.Response[pb.GetRoleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetRolePolicy returns information about a policy attached to a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRolePolicy(ctx context.Context, req *connect.Request[pb.GetRolePolicyRequest]) (*connect.Response[pb.GetRolePolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetSAMLProvider returns information about a SAML provider.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetSAMLProvider(ctx context.Context, req *connect.Request[pb.GetSAMLProviderRequest]) (*connect.Response[pb.GetSAMLProviderResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetSSHPublicKey returns information about an SSH public key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetSSHPublicKey(ctx context.Context, req *connect.Request[pb.GetSSHPublicKeyRequest]) (*connect.Response[pb.GetSSHPublicKeyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetServerCertificate returns information about a server certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetServerCertificate(ctx context.Context, req *connect.Request[pb.GetServerCertificateRequest]) (*connect.Response[pb.GetServerCertificateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetServiceLastAccessedDetails returns service last accessed details.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetServiceLastAccessedDetails(ctx context.Context, req *connect.Request[pb.GetServiceLastAccessedDetailsRequest]) (*connect.Response[pb.GetServiceLastAccessedDetailsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetServiceLastAccessedDetailsWithEntities returns service last accessed details with entities.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetServiceLastAccessedDetailsWithEntities(ctx context.Context, req *connect.Request[pb.GetServiceLastAccessedDetailsWithEntitiesRequest]) (*connect.Response[pb.GetServiceLastAccessedDetailsWithEntitiesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetServiceLinkedRoleDeletionStatus returns the deletion status of a service-linked role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetServiceLinkedRoleDeletionStatus(ctx context.Context, req *connect.Request[pb.GetServiceLinkedRoleDeletionStatusRequest]) (*connect.Response[pb.GetServiceLinkedRoleDeletionStatusResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetUser returns information about an IAM user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetUser(ctx context.Context, req *connect.Request[pb.GetUserRequest]) (*connect.Response[pb.GetUserResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetUserPolicy returns information about a policy attached to a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetUserPolicy(ctx context.Context, req *connect.Request[pb.GetUserPolicyRequest]) (*connect.Response[pb.GetUserPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListAccessKeys lists access keys for a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListAccessKeys(ctx context.Context, req *connect.Request[pb.ListAccessKeysRequest]) (*connect.Response[pb.ListAccessKeysResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListAccountAliases lists account aliases.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListAccountAliases(ctx context.Context, req *connect.Request[pb.ListAccountAliasesRequest]) (*connect.Response[pb.ListAccountAliasesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListAttachedGroupPolicies lists policies attached to a group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListAttachedGroupPolicies(ctx context.Context, req *connect.Request[pb.ListAttachedGroupPoliciesRequest]) (*connect.Response[pb.ListAttachedGroupPoliciesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListAttachedRolePolicies lists policies attached to a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListAttachedRolePolicies(ctx context.Context, req *connect.Request[pb.ListAttachedRolePoliciesRequest]) (*connect.Response[pb.ListAttachedRolePoliciesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListAttachedUserPolicies lists policies attached to a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListAttachedUserPolicies(ctx context.Context, req *connect.Request[pb.ListAttachedUserPoliciesRequest]) (*connect.Response[pb.ListAttachedUserPoliciesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListDelegationRequests lists delegation requests.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDelegationRequests(ctx context.Context, req *connect.Request[pb.ListDelegationRequestsRequest]) (*connect.Response[pb.ListDelegationRequestsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListEntitiesForPolicy lists entities that a policy is attached to.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListEntitiesForPolicy(ctx context.Context, req *connect.Request[pb.ListEntitiesForPolicyRequest]) (*connect.Response[pb.ListEntitiesForPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListGroupPolicies lists inline policies attached to a group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListGroupPolicies(ctx context.Context, req *connect.Request[pb.ListGroupPoliciesRequest]) (*connect.Response[pb.ListGroupPoliciesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListGroups lists IAM groups.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListGroups(ctx context.Context, req *connect.Request[pb.ListGroupsRequest]) (*connect.Response[pb.ListGroupsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListGroupsForUser lists groups that a user belongs to.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListGroupsForUser(ctx context.Context, req *connect.Request[pb.ListGroupsForUserRequest]) (*connect.Response[pb.ListGroupsForUserResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListInstanceProfileTags lists tags for an instance profile.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListInstanceProfileTags(ctx context.Context, req *connect.Request[pb.ListInstanceProfileTagsRequest]) (*connect.Response[pb.ListInstanceProfileTagsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListInstanceProfiles lists instance profiles.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListInstanceProfiles(ctx context.Context, req *connect.Request[pb.ListInstanceProfilesRequest]) (*connect.Response[pb.ListInstanceProfilesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListInstanceProfilesForRole lists instance profiles for a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListInstanceProfilesForRole(ctx context.Context, req *connect.Request[pb.ListInstanceProfilesForRoleRequest]) (*connect.Response[pb.ListInstanceProfilesForRoleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListMFADeviceTags lists tags for an MFA device.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListMFADeviceTags(ctx context.Context, req *connect.Request[pb.ListMFADeviceTagsRequest]) (*connect.Response[pb.ListMFADeviceTagsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListMFADevices lists MFA devices for a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListMFADevices(ctx context.Context, req *connect.Request[pb.ListMFADevicesRequest]) (*connect.Response[pb.ListMFADevicesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListOpenIDConnectProviderTags lists tags for an OpenID Connect provider.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListOpenIDConnectProviderTags(ctx context.Context, req *connect.Request[pb.ListOpenIDConnectProviderTagsRequest]) (*connect.Response[pb.ListOpenIDConnectProviderTagsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListOpenIDConnectProviders lists OpenID Connect providers.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListOpenIDConnectProviders(ctx context.Context, req *connect.Request[pb.ListOpenIDConnectProvidersRequest]) (*connect.Response[pb.ListOpenIDConnectProvidersResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListOrganizationsFeatures lists organisations features.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListOrganizationsFeatures(ctx context.Context, req *connect.Request[pb.ListOrganizationsFeaturesRequest]) (*connect.Response[pb.ListOrganizationsFeaturesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListPoliciesGrantingServiceAccess lists policies that grant access to a service.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListPoliciesGrantingServiceAccess(ctx context.Context, req *connect.Request[pb.ListPoliciesGrantingServiceAccessRequest]) (*connect.Response[pb.ListPoliciesGrantingServiceAccessResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListPolicyTags lists tags for a policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListPolicyTags(ctx context.Context, req *connect.Request[pb.ListPolicyTagsRequest]) (*connect.Response[pb.ListPolicyTagsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListPolicyVersions lists versions of a policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListPolicyVersions(ctx context.Context, req *connect.Request[pb.ListPolicyVersionsRequest]) (*connect.Response[pb.ListPolicyVersionsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListRolePolicies lists inline policies attached to a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListRolePolicies(ctx context.Context, req *connect.Request[pb.ListRolePoliciesRequest]) (*connect.Response[pb.ListRolePoliciesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListRoleTags lists tags for a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListRoleTags(ctx context.Context, req *connect.Request[pb.ListRoleTagsRequest]) (*connect.Response[pb.ListRoleTagsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListSAMLProviderTags lists tags for a SAML provider.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListSAMLProviderTags(ctx context.Context, req *connect.Request[pb.ListSAMLProviderTagsRequest]) (*connect.Response[pb.ListSAMLProviderTagsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListSAMLProviders lists SAML providers.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListSAMLProviders(ctx context.Context, req *connect.Request[pb.ListSAMLProvidersRequest]) (*connect.Response[pb.ListSAMLProvidersResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListSSHPublicKeys lists SSH public keys for a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListSSHPublicKeys(ctx context.Context, req *connect.Request[pb.ListSSHPublicKeysRequest]) (*connect.Response[pb.ListSSHPublicKeysResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListServerCertificateTags lists tags for a server certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListServerCertificateTags(ctx context.Context, req *connect.Request[pb.ListServerCertificateTagsRequest]) (*connect.Response[pb.ListServerCertificateTagsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListServerCertificates lists server certificates.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListServerCertificates(ctx context.Context, req *connect.Request[pb.ListServerCertificatesRequest]) (*connect.Response[pb.ListServerCertificatesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListServiceSpecificCredentials lists service-specific credentials for a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListServiceSpecificCredentials(ctx context.Context, req *connect.Request[pb.ListServiceSpecificCredentialsRequest]) (*connect.Response[pb.ListServiceSpecificCredentialsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListSigningCertificates lists signing certificates for a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListSigningCertificates(ctx context.Context, req *connect.Request[pb.ListSigningCertificatesRequest]) (*connect.Response[pb.ListSigningCertificatesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListUserPolicies lists inline policies attached to a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListUserPolicies(ctx context.Context, req *connect.Request[pb.ListUserPoliciesRequest]) (*connect.Response[pb.ListUserPoliciesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListUserTags lists tags for a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListUserTags(ctx context.Context, req *connect.Request[pb.ListUserTagsRequest]) (*connect.Response[pb.ListUserTagsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListVirtualMFADevices lists virtual MFA devices.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListVirtualMFADevices(ctx context.Context, req *connect.Request[pb.ListVirtualMFADevicesRequest]) (*connect.Response[pb.ListVirtualMFADevicesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PutGroupPolicy adds or updates an inline policy for a group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutGroupPolicy(ctx context.Context, req *connect.Request[pb.PutGroupPolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PutRolePermissionsBoundary sets the permissions boundary for a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutRolePermissionsBoundary(ctx context.Context, req *connect.Request[pb.PutRolePermissionsBoundaryRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PutRolePolicy adds or updates an inline policy for a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutRolePolicy(ctx context.Context, req *connect.Request[pb.PutRolePolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PutUserPermissionsBoundary sets the permissions boundary for a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutUserPermissionsBoundary(ctx context.Context, req *connect.Request[pb.PutUserPermissionsBoundaryRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PutUserPolicy adds or updates an inline policy for a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutUserPolicy(ctx context.Context, req *connect.Request[pb.PutUserPolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// RejectDelegationRequest rejects a delegation request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RejectDelegationRequest(ctx context.Context, req *connect.Request[pb.RejectDelegationRequestRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// RemoveClientIDFromOpenIDConnectProvider removes a client ID from an OpenID Connect provider.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RemoveClientIDFromOpenIDConnectProvider(ctx context.Context, req *connect.Request[pb.RemoveClientIDFromOpenIDConnectProviderRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// RemoveRoleFromInstanceProfile removes a role from an instance profile.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RemoveRoleFromInstanceProfile(ctx context.Context, req *connect.Request[pb.RemoveRoleFromInstanceProfileRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// RemoveUserFromGroup removes a user from a group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RemoveUserFromGroup(ctx context.Context, req *connect.Request[pb.RemoveUserFromGroupRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ResetServiceSpecificCredential resets a service-specific credential.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ResetServiceSpecificCredential(ctx context.Context, req *connect.Request[pb.ResetServiceSpecificCredentialRequest]) (*connect.Response[pb.ResetServiceSpecificCredentialResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ResyncMFADevice synchronizes an MFA device.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ResyncMFADevice(ctx context.Context, req *connect.Request[pb.ResyncMFADeviceRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// SendDelegationToken sends a delegation token.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) SendDelegationToken(ctx context.Context, req *connect.Request[pb.SendDelegationTokenRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// SetDefaultPolicyVersion sets the default version of a policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) SetDefaultPolicyVersion(ctx context.Context, req *connect.Request[pb.SetDefaultPolicyVersionRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// SetSecurityTokenServicePreferences sets preferences for the Security Token Service.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) SetSecurityTokenServicePreferences(ctx context.Context, req *connect.Request[pb.SetSecurityTokenServicePreferencesRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// SimulateCustomPolicy simulates whether a custom policy would allow an operation.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) SimulateCustomPolicy(ctx context.Context, req *connect.Request[pb.SimulateCustomPolicyRequest]) (*connect.Response[pb.SimulatePolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// SimulatePrincipalPolicy simulates whether a principal can perform an operation.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) SimulatePrincipalPolicy(ctx context.Context, req *connect.Request[pb.SimulatePrincipalPolicyRequest]) (*connect.Response[pb.SimulatePolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// TagInstanceProfile adds tags to an instance profile.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagInstanceProfile(ctx context.Context, req *connect.Request[pb.TagInstanceProfileRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// TagMFADevice adds tags to an MFA device.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagMFADevice(ctx context.Context, req *connect.Request[pb.TagMFADeviceRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// TagOpenIDConnectProvider adds tags to an OpenID Connect provider.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagOpenIDConnectProvider(ctx context.Context, req *connect.Request[pb.TagOpenIDConnectProviderRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// TagPolicy adds tags to a policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagPolicy(ctx context.Context, req *connect.Request[pb.TagPolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// TagRole adds tags to a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagRole(ctx context.Context, req *connect.Request[pb.TagRoleRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// TagSAMLProvider adds tags to a SAML provider.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagSAMLProvider(ctx context.Context, req *connect.Request[pb.TagSAMLProviderRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// TagServerCertificate adds tags to a server certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagServerCertificate(ctx context.Context, req *connect.Request[pb.TagServerCertificateRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// TagUser adds tags to a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagUser(ctx context.Context, req *connect.Request[pb.TagUserRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UntagInstanceProfile removes tags from an instance profile.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagInstanceProfile(ctx context.Context, req *connect.Request[pb.UntagInstanceProfileRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UntagMFADevice removes tags from an MFA device.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagMFADevice(ctx context.Context, req *connect.Request[pb.UntagMFADeviceRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UntagOpenIDConnectProvider removes tags from an OpenID Connect provider.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagOpenIDConnectProvider(ctx context.Context, req *connect.Request[pb.UntagOpenIDConnectProviderRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UntagPolicy removes tags from a policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagPolicy(ctx context.Context, req *connect.Request[pb.UntagPolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UntagRole removes tags from a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagRole(ctx context.Context, req *connect.Request[pb.UntagRoleRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UntagSAMLProvider removes tags from a SAML provider.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagSAMLProvider(ctx context.Context, req *connect.Request[pb.UntagSAMLProviderRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UntagServerCertificate removes tags from a server certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagServerCertificate(ctx context.Context, req *connect.Request[pb.UntagServerCertificateRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UntagUser removes tags from a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagUser(ctx context.Context, req *connect.Request[pb.UntagUserRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateAccessKey updates the status of an access key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateAccessKey(ctx context.Context, req *connect.Request[pb.UpdateAccessKeyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateAccountPasswordPolicy updates the account password policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateAccountPasswordPolicy(ctx context.Context, req *connect.Request[pb.UpdateAccountPasswordPolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateAssumeRolePolicy updates the trust policy for a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateAssumeRolePolicy(ctx context.Context, req *connect.Request[pb.UpdateAssumeRolePolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateDelegationRequest updates a delegation request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateDelegationRequest(ctx context.Context, req *connect.Request[pb.UpdateDelegationRequestRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateGroup updates the name or path of a group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateGroup(ctx context.Context, req *connect.Request[pb.UpdateGroupRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateLoginProfile updates the password for a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateLoginProfile(ctx context.Context, req *connect.Request[pb.UpdateLoginProfileRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateOpenIDConnectProviderThumbprint updates the thumbprint for an OpenID Connect provider.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateOpenIDConnectProviderThumbprint(ctx context.Context, req *connect.Request[pb.UpdateOpenIDConnectProviderThumbprintRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateRole updates a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateRole(ctx context.Context, req *connect.Request[pb.UpdateRoleRequest]) (*connect.Response[pb.UpdateRoleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateRoleDescription updates the description of a role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateRoleDescription(ctx context.Context, req *connect.Request[pb.UpdateRoleDescriptionRequest]) (*connect.Response[pb.UpdateRoleDescriptionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateSAMLProvider updates a SAML provider.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateSAMLProvider(ctx context.Context, req *connect.Request[pb.UpdateSAMLProviderRequest]) (*connect.Response[pb.UpdateSAMLProviderResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateSSHPublicKey updates the status of an SSH public key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateSSHPublicKey(ctx context.Context, req *connect.Request[pb.UpdateSSHPublicKeyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateServerCertificate updates the name or path of a server certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateServerCertificate(ctx context.Context, req *connect.Request[pb.UpdateServerCertificateRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateServiceSpecificCredential updates the status of a service-specific credential.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateServiceSpecificCredential(ctx context.Context, req *connect.Request[pb.UpdateServiceSpecificCredentialRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateSigningCertificate updates the status of a signing certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateSigningCertificate(ctx context.Context, req *connect.Request[pb.UpdateSigningCertificateRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateUser updates the name or path of a user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateUser(ctx context.Context, req *connect.Request[pb.UpdateUserRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UploadServerCertificate uploads a server certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UploadServerCertificate(ctx context.Context, req *connect.Request[pb.UploadServerCertificateRequest]) (*connect.Response[pb.UploadServerCertificateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UploadSigningCertificate uploads a signing certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UploadSigningCertificate(ctx context.Context, req *connect.Request[pb.UploadSigningCertificateRequest]) (*connect.Response[pb.UploadSigningCertificateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UploadSSHPublicKey uploads an SSH public key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UploadSSHPublicKey(ctx context.Context, req *connect.Request[pb.UploadSSHPublicKeyRequest]) (*connect.Response[pb.UploadSSHPublicKeyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}
