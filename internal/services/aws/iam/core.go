// Package iam provides IAM service operations for vorpalstacks.
//
// This file contains transport-agnostic core functions (the xxxCore pattern)
// that encapsulate validation and store operations for IAM entities.  Both
// the AWS-compatible HTTP API handlers (in *_operations.go) and the admin
// gRPC-Web handler (in admin_handler.go) delegate to these functions,
// ensuring that validation is applied identically at every entry point.
package iam

import (
	"errors"

	"vorpalstacks/internal/common/pagination"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/timeutils"
)

// ---------------------------------------------------------------------------
// Input structs (transport-agnostic)
// ---------------------------------------------------------------------------

// CreateUserInput holds the parameters for creating an IAM user.
type CreateUserInput struct {
	UserName               string
	Path                   string
	PermissionsBoundaryArn string
	Tags                   []types.Tag
}

// CreateRoleInput holds the parameters for creating an IAM role.
type CreateRoleInput struct {
	RoleName                 string
	Path                     string
	AssumeRolePolicyDocument string
	Description              string
	MaxSessionDuration       int
	PermissionsBoundaryArn   string
	Tags                     []types.Tag
}

// CreateGroupInput holds the parameters for creating an IAM group.
type CreateGroupInput struct {
	GroupName string
	Path      string
}

// CreatePolicyInput holds the parameters for creating an IAM managed policy.
type CreatePolicyInput struct {
	PolicyName     string
	Path           string
	PolicyDocument string
	Description    string
	Tags           []types.Tag
}

// ---------------------------------------------------------------------------
// Create core functions
// ---------------------------------------------------------------------------

// createUserCore validates input and creates an IAM user in the store.
// Returns the created user or an IAM-formatted error.
func (s *IAMService) createUserCore(store *iamstore.IAMStore, input *CreateUserInput) (*iamstore.User, error) {
	if input.UserName == "" {
		return nil, NewInvalidInputError("UserName", "cannot be empty")
	}
	if err := validateEntityName(input.UserName, "UserName"); err != nil {
		return nil, err
	}

	path := input.Path
	if path == "" {
		path = "/"
	}
	if !validatePath(path) {
		return nil, NewInvalidInputError("Path", "must be a valid path starting and ending with /")
	}

	if err := validateNewTags(input.Tags); err != nil {
		return nil, err
	}

	user, err := store.Users().Create(input.UserName, path, store.AccountID(), input.Tags)
	if err != nil {
		if errors.Is(err, iamstore.ErrUserAlreadyExists) {
			return nil, NewUserAlreadyExistsError(input.UserName)
		}
		return nil, err
	}

	// Apply permissions boundary if specified at creation time (Smithy
	// CreateUserInput.PermissionsBoundary).
	if input.PermissionsBoundaryArn != "" {
		if err := putUserPermissionsBoundaryCore(store, user, input.PermissionsBoundaryArn); err != nil {
			return nil, err
		}
	}

	return user, nil
}

// createRoleCore validates input and creates an IAM role in the store.
// Returns the created role or an IAM-formatted error.
func (s *IAMService) createRoleCore(store *iamstore.IAMStore, input *CreateRoleInput) (*iamstore.Role, error) {
	if input.RoleName == "" {
		return nil, NewInvalidInputError("RoleName", "cannot be empty")
	}
	if err := validateEntityName(input.RoleName, "RoleName"); err != nil {
		return nil, err
	}

	path := input.Path
	if path == "" {
		path = "/"
	}
	if !validatePath(path) {
		return nil, NewInvalidInputError("Path", "must be a valid path starting and ending with /")
	}

	if input.AssumeRolePolicyDocument == "" {
		return nil, ErrMalformedPolicyDocument
	}
	if !validateTrustPolicyDocument(input.AssumeRolePolicyDocument) {
		return nil, ErrMalformedPolicyDocument
	}

	if !validateRoleDescription(input.Description) {
		return nil, NewInvalidInputError("Description", "must be 0 to 1000 characters; allowed: tab, LF, CR, printable ASCII, Latin-1 supplement")
	}

	maxSessionDuration := input.MaxSessionDuration
	if maxSessionDuration == 0 {
		maxSessionDuration = 3600
	}
	if !validateRoleMaxSessionDuration(maxSessionDuration) {
		return nil, NewInvalidInputError("MaxSessionDuration", "must be between 3600 and 43200 seconds")
	}

	if err := validateNewTags(input.Tags); err != nil {
		return nil, err
	}

	role, err := store.Roles().Create(
		input.RoleName, path, store.AccountID(),
		input.AssumeRolePolicyDocument, input.Description,
		maxSessionDuration, input.Tags,
	)
	if err != nil {
		if errors.Is(err, iamstore.ErrRoleAlreadyExists) {
			return nil, NewRoleAlreadyExistsError(input.RoleName)
		}
		return nil, err
	}

	// Apply permissions boundary if specified at creation time (Smithy
	// CreateRoleInput.PermissionsBoundary).
	if input.PermissionsBoundaryArn != "" {
		if err := putRolePermissionsBoundaryCore(store, role, input.PermissionsBoundaryArn); err != nil {
			return nil, err
		}
	}

	return role, nil
}

// putUserPermissionsBoundaryCore atomically sets a permissions boundary
// on an IAM user.  It validates the ARN, checks the policy exists,
// handles old-boundary decrement and same-ARN idempotency, persists the
// user, and increments the new policy's usage count.
//
// Used by both createUserCore (when PermissionsBoundaryArn is specified
// at creation time) and the PutUserPermissionsBoundary HTTP operation.
// Consolidating the logic here prevents the create-time vs update-time
// drift that previously existed.
func putUserPermissionsBoundaryCore(store *iamstore.IAMStore, user *iamstore.User, pbArn string) error {
	if err := validateIAMPolicyArn(pbArn); err != nil {
		return err
	}
	if !store.Policies().Exists(pbArn) {
		return NewNoSuchPolicyError(pbArn)
	}
	// Idempotent: same ARN already set — nothing to do.
	if user.PermissionsBoundary != nil && user.PermissionsBoundary.PermissionsBoundaryArn == pbArn {
		return nil
	}
	// Decrement the previous boundary's usage count, if any.  At create
	// time user.PermissionsBoundary is nil so this is a no-op there.
	if user.PermissionsBoundary != nil && user.PermissionsBoundary.PermissionsBoundaryArn != "" {
		_ = store.Policies().DecrementPermissionsBoundaryUsageCount(user.PermissionsBoundary.PermissionsBoundaryArn)
	}
	user.PermissionsBoundary = &iamstore.PermissionsBoundary{
		PermissionsBoundaryType: "Policy",
		PermissionsBoundaryArn:  pbArn,
	}
	if err := store.Users().Put(user); err != nil {
		return err
	}
	_ = store.Policies().IncrementPermissionsBoundaryUsageCount(pbArn)
	return nil
}

// putRolePermissionsBoundaryCore is the role equivalent of
// putUserPermissionsBoundaryCore.
func putRolePermissionsBoundaryCore(store *iamstore.IAMStore, role *iamstore.Role, pbArn string) error {
	if err := validateIAMPolicyArn(pbArn); err != nil {
		return err
	}
	if !store.Policies().Exists(pbArn) {
		return NewNoSuchPolicyError(pbArn)
	}
	if role.PermissionsBoundary != nil && role.PermissionsBoundary.PermissionsBoundaryArn == pbArn {
		return nil
	}
	if role.PermissionsBoundary != nil && role.PermissionsBoundary.PermissionsBoundaryArn != "" {
		_ = store.Policies().DecrementPermissionsBoundaryUsageCount(role.PermissionsBoundary.PermissionsBoundaryArn)
	}
	role.PermissionsBoundary = &iamstore.PermissionsBoundary{
		PermissionsBoundaryType: "Policy",
		PermissionsBoundaryArn:  pbArn,
	}
	if err := store.Roles().Put(role); err != nil {
		return err
	}
	_ = store.Policies().IncrementPermissionsBoundaryUsageCount(pbArn)
	return nil
}

// createGroupCore validates input and creates an IAM group in the store.
// Returns the created group or an IAM-formatted error.
func (s *IAMService) createGroupCore(store *iamstore.IAMStore, input *CreateGroupInput) (*iamstore.Group, error) {
	if input.GroupName == "" {
		return nil, NewInvalidInputError("GroupName", "cannot be empty")
	}
	if err := validateEntityName128(input.GroupName, "GroupName"); err != nil {
		return nil, err
	}

	path := input.Path
	if path == "" {
		path = "/"
	}
	if !validatePath(path) {
		return nil, NewInvalidInputError("Path", "must be a valid path starting and ending with /")
	}

	group, err := store.Groups().Create(input.GroupName, path, store.AccountID())
	if err != nil {
		if errors.Is(err, iamstore.ErrGroupAlreadyExists) {
			return nil, NewGroupAlreadyExistsError(input.GroupName)
		}
		return nil, err
	}
	return group, nil
}

// createPolicyCore validates input and creates an IAM managed policy in the
// store.  Returns the created policy or an IAM-formatted error.
func (s *IAMService) createPolicyCore(store *iamstore.IAMStore, input *CreatePolicyInput) (*iamstore.Policy, error) {
	if input.PolicyName == "" {
		return nil, NewInvalidInputError("PolicyName", "cannot be empty")
	}
	if err := validateEntityName128(input.PolicyName, "PolicyName"); err != nil {
		return nil, err
	}

	path := input.Path
	if path == "" {
		path = "/"
	}
	if !validatePath(path) {
		return nil, NewInvalidInputError("Path", "must be a valid path starting and ending with /")
	}

	if !validatePolicyDocument(input.PolicyDocument) {
		return nil, ErrMalformedPolicyDocument
	}

	if err := validateNewTags(input.Tags); err != nil {
		return nil, err
	}

	policy, err := store.Policies().Create(input.PolicyName, path, store.AccountID(), input.PolicyDocument, input.Description, input.Tags)
	if err != nil {
		if errors.Is(err, iamstore.ErrPolicyAlreadyExists) {
			return nil, NewPolicyAlreadyExistsError(input.PolicyName)
		}
		return nil, err
	}
	return policy, nil
}

// ---------------------------------------------------------------------------
// Get/List core functions: route admin handler and HTTP API through cores.
//
// These functions encapsulate the store lookup + not-found error mapping so
// that both the admin gRPC-Web handler and the AWS-compatible HTTP API
// share a single code path.  Future validation or filtering only needs to
// be added here.
// ---------------------------------------------------------------------------

// getUserCore returns the IAM user with the given name.  Callers must
// validate that userName is non-empty before calling.
func (s *IAMService) getUserCore(store *iamstore.IAMStore, userName string) (*iamstore.User, error) {
	user, err := store.Users().Get(userName)
	if err != nil {
		return nil, NewNoSuchUserError(userName)
	}
	return user, nil
}

// listUsersCore returns a paginated list of IAM users.
func (s *IAMService) listUsersCore(store *iamstore.IAMStore, pathPrefix, marker string, maxItems int) (*iamstore.UserListResult, error) {
	return store.Users().List(pathPrefix, marker, maxItems)
}

// getRoleCore returns the IAM role with the given name.  Callers must
// validate that roleName is non-empty before calling.
func (s *IAMService) getRoleCore(store *iamstore.IAMStore, roleName string) (*iamstore.Role, error) {
	role, err := store.Roles().Get(roleName)
	if err != nil {
		return nil, NewNoSuchRoleError(roleName)
	}
	return role, nil
}

// listRolesCore returns a paginated list of IAM roles.
func (s *IAMService) listRolesCore(store *iamstore.IAMStore, pathPrefix, marker string, maxItems int) (*iamstore.RoleListResult, error) {
	return store.Roles().List(pathPrefix, marker, maxItems)
}

// getPolicyCore returns the IAM managed policy with the given ARN.
// Callers must validate that policyArn is non-empty before calling.
func (s *IAMService) getPolicyCore(store *iamstore.IAMStore, policyArn string) (*iamstore.Policy, error) {
	policy, err := store.Policies().Get(policyArn)
	if err != nil {
		return nil, NewNoSuchPolicyError(policyArn)
	}
	return policy, nil
}

// listPoliciesCore returns a paginated list of IAM managed policies.
func (s *IAMService) listPoliciesCore(store *iamstore.IAMStore, scope, pathPrefix, marker string, onlyAttached bool, maxItems int) (*iamstore.PolicyListResult, error) {
	return store.Policies().List(scope, pathPrefix, onlyAttached, marker, maxItems)
}

// getGroupCore returns the IAM group with the given name.  Callers must
// validate that groupName is non-empty before calling.
func (s *IAMService) getGroupCore(store *iamstore.IAMStore, groupName string) (*iamstore.Group, error) {
	group, err := store.Groups().Get(groupName)
	if err != nil {
		return nil, NewNoSuchGroupError(groupName)
	}
	return group, nil
}

// listGroupsCore returns a paginated list of IAM groups.
func (s *IAMService) listGroupsCore(store *iamstore.IAMStore, pathPrefix, marker string, maxItems int) (*iamstore.GroupListResult, error) {
	return store.Groups().List(pathPrefix, marker, maxItems)
}

// updateAssumeRolePolicyCore validates the trust policy document and
// replaces the role's AssumeRolePolicyDocument field.  Used by the HTTP
// API UpdateAssumeRolePolicy operation; routing through core keeps the
// validation identical if a future admin handler mirrors the op.
func (s *IAMService) updateAssumeRolePolicyCore(store *iamstore.IAMStore, roleName, policyDocument string) error {
	role, err := s.getRoleCore(store, roleName)
	if err != nil {
		return err
	}
	if !validateTrustPolicyDocument(policyDocument) {
		return ErrMalformedPolicyDocument
	}
	role.AssumeRolePolicyDocument = policyDocument
	return store.Roles().Put(role)
}

// ---------------------------------------------------------------------------
// Update input structs
// ---------------------------------------------------------------------------

// UpdateUserInput holds the parameters for renaming/updating an IAM user.
type UpdateUserInput struct {
	UserName    string
	NewPath     string
	NewUserName string
}

// UpdateRoleInput holds the parameters for updating an IAM role's
// description and/or max session duration.
type UpdateRoleInput struct {
	RoleName           string
	Description        string
	MaxSessionDuration int
}

// UpdateGroupInput holds the parameters for renaming/updating an IAM group.
type UpdateGroupInput struct {
	GroupName    string
	NewPath      string
	NewGroupName string
}

// ---------------------------------------------------------------------------
// Update core functions
// ---------------------------------------------------------------------------

// updateUserCore validates input and renames/repaths an IAM user.
// Returns the updated user or an IAM-formatted error.
func (s *IAMService) updateUserCore(store *iamstore.IAMStore, input *UpdateUserInput) (*iamstore.User, error) {
	if input.UserName == "" {
		return nil, NewValidationError("UserName")
	}

	if input.NewPath == "" && input.NewUserName == "" {
		return nil, NewInvalidInputError("UpdateUser", "at least one of NewPath or NewUserName must be specified")
	}
	if input.NewPath != "" && !validatePath(input.NewPath) {
		return nil, NewInvalidInputError("NewPath", "must be a valid path starting and ending with /")
	}

	if err := store.RenameUser(input.UserName, input.NewUserName, input.NewPath); err != nil {
		return nil, err
	}

	targetName := input.UserName
	if input.NewUserName != "" {
		targetName = input.NewUserName
	}
	user, err := store.Users().Get(targetName)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// updateRoleCore validates input and updates an IAM role's description and/or
// max session duration.  Returns the updated role or an IAM-formatted error.
func (s *IAMService) updateRoleCore(store *iamstore.IAMStore, input *UpdateRoleInput) (*iamstore.Role, error) {
	if input.RoleName == "" {
		return nil, NewValidationError("RoleName")
	}

	maxSessionDuration := input.MaxSessionDuration
	if maxSessionDuration > 0 {
		if !validateRoleMaxSessionDuration(maxSessionDuration) {
			return nil, NewInvalidInputError("MaxSessionDuration", "must be between 3600 and 43200 seconds")
		}
	}

	if !validateRoleDescription(input.Description) {
		return nil, NewInvalidInputError("Description", "must be 0 to 1000 characters; allowed: tab, LF, CR, printable ASCII, Latin-1 supplement")
	}

	if err := store.UpdateRoleFields(input.RoleName, input.Description, maxSessionDuration); err != nil {
		return nil, err
	}

	role, err := store.Roles().Get(input.RoleName)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// updateGroupCore validates input and renames/repaths an IAM group.
// Returns the updated group or an IAM-formatted error.
func (s *IAMService) updateGroupCore(store *iamstore.IAMStore, input *UpdateGroupInput) (*iamstore.Group, error) {
	if input.GroupName == "" {
		return nil, NewValidationError("GroupName")
	}

	if input.NewPath == "" && input.NewGroupName == "" {
		return nil, NewInvalidInputError("UpdateGroup", "at least one of NewPath or NewGroupName must be specified")
	}
	if input.NewPath != "" && !validatePath(input.NewPath) {
		return nil, NewInvalidInputError("NewPath", "must be a valid path starting and ending with /")
	}

	if err := store.RenameGroup(input.GroupName, input.NewGroupName, input.NewPath); err != nil {
		return nil, err
	}

	targetName := input.GroupName
	if input.NewGroupName != "" {
		targetName = input.NewGroupName
	}
	group, err := store.Groups().Get(targetName)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// ---------------------------------------------------------------------------
// Delete input structs
// ---------------------------------------------------------------------------

// DeleteUserInput holds the parameters for deleting an IAM user.
// When Cascade is true, all dependent resources are removed before the user
// record (admin handler behaviour).  When false, a DeleteConflict error is
// returned if any dependent resource remains (AWS API behaviour).
type DeleteUserInput struct {
	UserName string
	Cascade  bool
}

// DeleteRoleInput holds the parameters for deleting an IAM role.
type DeleteRoleInput struct {
	RoleName string
	Cascade  bool
}

// DeleteGroupInput holds the parameters for deleting an IAM group.
type DeleteGroupInput struct {
	GroupName string
	Cascade   bool
}

// DeletePolicyInput holds the parameters for deleting an IAM managed policy.
// Both AWS API and admin handler enforce the AttachmentCount > 0 check.
type DeletePolicyInput struct {
	PolicyArn string
}

// ---------------------------------------------------------------------------
// Delete core functions
// ---------------------------------------------------------------------------

// deleteUserCore validates input and deletes an IAM user.
// When input.Cascade is true, all dependent resources (login profile, access
// keys, certificates, MFA devices, inline/attached policies, group
// memberships) are removed first via cascadeDeleteUser.  When false, a
// DeleteConflict error is returned if any dependent resource remains.
// In both modes the permissions-boundary usage count is decremented
// best-effort before the user record is removed.
func (s *IAMService) deleteUserCore(store *iamstore.IAMStore, input *DeleteUserInput) error {
	if input.UserName == "" {
		return NewValidationError("UserName")
	}
	if !store.Users().Exists(input.UserName) {
		return NewNoSuchUserError(input.UserName)
	}

	if input.Cascade {
		return cascadeDeleteUser(store, input.UserName)
	}

	// Conflict detection — AWS API path.
	if store.LoginProfiles().Exists(input.UserName) {
		return NewDeleteConflictError("Cannot delete entity, must delete login profile first.")
	}

	keyCount, err := store.AccessKeys().CountByUserName(input.UserName)
	if err != nil {
		return err
	}
	if keyCount > 0 {
		return NewDeleteConflictError("Cannot delete entity, must delete access keys first.")
	}

	certs, err := store.SigningCertificates().ListByUserName(input.UserName)
	if err != nil {
		return err
	}
	if len(certs) > 0 {
		return NewDeleteConflictError("Cannot delete entity, must delete signing certificates first.")
	}

	sshKeyCount, err := store.SSHPublicKeys().CountByUserName(input.UserName)
	if err != nil {
		return err
	}
	if sshKeyCount > 0 {
		return NewDeleteConflictError("Cannot delete entity, must delete SSH public keys first.")
	}

	svcCreds, err := store.ServiceSpecificCredentials().ListByUserName(input.UserName)
	if err != nil {
		return err
	}
	if len(svcCreds) > 0 {
		return NewDeleteConflictError("Cannot delete entity, must delete service-specific credentials first.")
	}

	mfaResult, err := store.MFADevices().ListForUser(input.UserName, "", 1)
	if err != nil {
		return err
	}
	if len(mfaResult.MFADevices) > 0 {
		return NewDeleteConflictError("Cannot delete entity, must deactivate MFA devices first.")
	}

	inlinePolicies, err := store.InlinePolicies().List(PrincipalTypeUser, input.UserName)
	if err != nil {
		return err
	}
	if len(inlinePolicies) > 0 {
		return NewDeleteConflictError("Cannot delete entity, must delete inline policies first.")
	}

	attachedPolicies, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeUser, input.UserName)
	if err != nil {
		return err
	}
	if len(attachedPolicies) > 0 {
		return NewDeleteConflictError("Cannot delete entity, must detach managed policies first.")
	}

	groups, err := store.UserGroups().ListGroupsForUser(input.UserName)
	if err != nil {
		return err
	}
	if len(groups) > 0 {
		return NewDeleteConflictError("Cannot delete entity, must remove user from groups first.")
	}

	// Decrement permissions boundary usage count before the user record is
	// removed. AWS allows deleting an entity that still has a permissions
	// boundary attached (it is not a deletion prerequisite), so the policy
	// counter must be adjusted here to avoid drift. Best-effort, matching the
	// PutUserPermissionsBoundary / DeleteUserPermissionsBoundary pattern.
	if user, gErr := store.Users().Get(input.UserName); gErr == nil {
		if user.PermissionsBoundary != nil && user.PermissionsBoundary.PermissionsBoundaryArn != "" {
			_ = store.Policies().DecrementPermissionsBoundaryUsageCount(user.PermissionsBoundary.PermissionsBoundaryArn)
		}
	}

	return store.Users().Delete(input.UserName)
}

// deleteRoleCore validates input and deletes an IAM role.
// When input.Cascade is true, all dependent resources (inline/attached
// policies, instance profile associations) are removed first via
// cascadeDeleteRole.  When false, a DeleteConflict error is returned if any
// dependent resource remains.
// In both modes the permissions-boundary usage count is decremented
// best-effort before the role record is removed.
func (s *IAMService) deleteRoleCore(store *iamstore.IAMStore, input *DeleteRoleInput) error {
	if input.RoleName == "" {
		return ErrNoSuchRole
	}
	if !store.Roles().Exists(input.RoleName) {
		return NewNoSuchRoleError(input.RoleName)
	}

	if input.Cascade {
		return cascadeDeleteRole(store, input.RoleName)
	}

	// Conflict detection — AWS API path.
	instanceProfiles, err := store.InstanceProfiles().ListForRole(input.RoleName, "", 1)
	if err != nil {
		return err
	}
	if len(instanceProfiles.InstanceProfiles) > 0 {
		return NewDeleteRoleConflictError("Cannot delete entity, must remove role from instance profile first.")
	}

	inlinePolicies, err := store.InlinePolicies().List(PrincipalTypeRole, input.RoleName)
	if err != nil {
		return err
	}
	if len(inlinePolicies) > 0 {
		return NewDeleteRoleConflictError("Cannot delete entity, must delete policies first.")
	}

	attachedPolicies, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeRole, input.RoleName)
	if err != nil {
		return err
	}
	if len(attachedPolicies) > 0 {
		return NewDeleteRoleConflictError("Cannot delete entity, must detach policies first.")
	}

	// Decrement permissions boundary usage count before the role record is
	// removed. AWS allows deleting an entity that still has a permissions
	// boundary attached (it is not a deletion prerequisite), so the policy
	// counter must be adjusted here to avoid drift. Best-effort, matching the
	// PutRolePermissionsBoundary / DeleteRolePermissionsBoundary pattern.
	if role, gErr := store.Roles().Get(input.RoleName); gErr == nil {
		if role.PermissionsBoundary != nil && role.PermissionsBoundary.PermissionsBoundaryArn != "" {
			_ = store.Policies().DecrementPermissionsBoundaryUsageCount(role.PermissionsBoundary.PermissionsBoundaryArn)
		}
	}

	return store.Roles().Delete(input.RoleName)
}

// deleteGroupCore validates input and deletes an IAM group.
// When input.Cascade is true, all dependent resources (inline/attached
// policies, user memberships) are removed first via cascadeDeleteGroup.
// When false, a DeleteConflict error is returned if any dependent resource
// remains.
func (s *IAMService) deleteGroupCore(store *iamstore.IAMStore, input *DeleteGroupInput) error {
	if input.GroupName == "" {
		return ErrNoSuchGroup
	}
	if !store.Groups().Exists(input.GroupName) {
		return NewNoSuchGroupError(input.GroupName)
	}

	if input.Cascade {
		return cascadeDeleteGroup(store, input.GroupName)
	}

	// Conflict detection — AWS API path.
	userCount := store.UserGroups().CountUsersInGroup(input.GroupName)
	if userCount > 0 {
		return NewDeleteGroupConflictError("Cannot delete entity, must remove users from group first.")
	}

	inlinePolicies, err := store.InlinePolicies().List(PrincipalTypeGroup, input.GroupName)
	if err != nil {
		return err
	}
	if len(inlinePolicies) > 0 {
		return NewDeleteGroupConflictError("Cannot delete entity, must delete policies first.")
	}

	attachedPolicies, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeGroup, input.GroupName)
	if err != nil {
		return err
	}
	if len(attachedPolicies) > 0 {
		return NewDeleteGroupConflictError("Cannot delete entity, must detach policies first.")
	}

	return store.Groups().Delete(input.GroupName)
}

// deletePolicyCore validates input and deletes an IAM managed policy.
// Rejects if the policy has active attachments (AttachmentCount > 0).
// All policy versions are cleaned up before the policy record is removed.
func (s *IAMService) deletePolicyCore(store *iamstore.IAMStore, input *DeletePolicyInput) error {
	if input.PolicyArn == "" {
		return ErrNoSuchPolicy
	}

	policy, err := store.Policies().Get(input.PolicyArn)
	if err != nil {
		return NewNoSuchPolicyError(input.PolicyArn)
	}

	if policy.AttachmentCount > 0 {
		return NewDeletePolicyConflictError(input.PolicyArn)
	}

	if err := store.Policies().DeleteAllVersions(input.PolicyArn); err != nil {
		return err
	}

	return store.Policies().Delete(input.PolicyArn)
}

// deleteUserPermissionsBoundaryCore removes the permissions boundary
// from a user and decrements the usage count on the previously-bound
// policy.  Consolidates the logic that was previously inline in the
// HTTP API handler so that future admin-handler paths can delegate
// here as well.
func (s *IAMService) deleteUserPermissionsBoundaryCore(store *iamstore.IAMStore, userName string) error {
	user, err := store.Users().Get(userName)
	if err != nil {
		return NewNoSuchUserError(userName)
	}

	if user.PermissionsBoundary != nil && user.PermissionsBoundary.PermissionsBoundaryArn != "" {
		_ = store.Policies().DecrementPermissionsBoundaryUsageCount(user.PermissionsBoundary.PermissionsBoundaryArn)
	}

	user.PermissionsBoundary = nil
	return store.Users().Put(user)
}

// deleteRolePermissionsBoundaryCore removes the permissions boundary
// from a role and decrements the usage count on the previously-bound
// policy.  Mirrors deleteUserPermissionsBoundaryCore for roles.
func (s *IAMService) deleteRolePermissionsBoundaryCore(store *iamstore.IAMStore, roleName string) error {
	role, err := store.Roles().Get(roleName)
	if err != nil {
		return NewNoSuchRoleError(roleName)
	}

	if role.PermissionsBoundary != nil && role.PermissionsBoundary.PermissionsBoundaryArn != "" {
		_ = store.Policies().DecrementPermissionsBoundaryUsageCount(role.PermissionsBoundary.PermissionsBoundaryArn)
	}

	role.PermissionsBoundary = nil
	return store.Roles().Put(role)
}

// listGroupsForUserCore retrieves the list of groups that a user
// belongs to.  Consolidates the store-direct logic so that both the
// HTTP API and future admin-handler paths delegate here.
func (s *IAMService) listGroupsForUserCore(store *iamstore.IAMStore, userName string) ([]*iamstore.Group, error) {
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	groupNames, err := store.UserGroups().ListGroupsForUser(userName)
	if err != nil {
		return nil, err
	}

	groups := make([]*iamstore.Group, 0, len(groupNames))
	for _, groupName := range groupNames {
		if group, err := store.Groups().Get(groupName); err == nil {
			groups = append(groups, group)
		}
	}
	return groups, nil
}

// AccountAuthorizationDetailsInput holds the parameters for
// GetAccountAuthorizationDetails.
type AccountAuthorizationDetailsInput struct {
	Filters  map[string]bool
	Marker   string
	MaxItems int
}

// getAccountAuthorizationDetailsCore aggregates all IAM entities (users,
// groups, roles, local managed policies) with their inline and attached
// policy relationships.  The logic is consolidated in core.go so that
// both the HTTP API and future admin-handler paths delegate here.
func (s *IAMService) getAccountAuthorizationDetailsCore(store *iamstore.IAMStore, input *AccountAuthorizationDetailsInput) (interface{}, error) {
	filters := input.Filters
	if len(filters) == 0 {
		filters = map[string]bool{
			"User":               true,
			"Group":              true,
			"Role":               true,
			"LocalManagedPolicy": true,
		}
	}

	marker := input.Marker
	maxItems := input.MaxItems

	type section struct {
		key    string
		filter string
		items  []interface{}
		marker func(i int) string
	}
	sections := make([]section, 0, 4)

	if filters["User"] {
		users, err := listAllUsers(store, "")
		if err != nil {
			return nil, err
		}
		items := make([]interface{}, 0, len(users))
		for _, user := range users {
			detail := map[string]interface{}{
				"UserId":     user.ID,
				"Path":       user.Path,
				"UserName":   user.UserName,
				"Arn":        user.Arn,
				"CreateDate": user.CreateDate.Format(timeutils.ISO8601SimpleFormat),
			}
			groupNames, _ := store.UserGroups().ListGroupsForUser(user.UserName)
			groupList := make([]interface{}, 0, len(groupNames))
			for _, gn := range groupNames {
				groupList = append(groupList, gn)
			}
			detail["GroupList"] = groupList
			detail["AttachedManagedPolicies"] = buildAttachedManagedPolicies(store, PrincipalTypeUser, user.UserName)
			detail["UserPolicyList"] = buildInlinePolicyList(store, PrincipalTypeUser, user.UserName)
			items = append(items, detail)
		}
		sections = append(sections, section{
			key: "UserDetailList", filter: "User", items: items,
			marker: func(i int) string { return "user:" + users[i].UserName },
		})
	}

	if filters["Group"] {
		groups, err := listAllGroups(store, "")
		if err != nil {
			return nil, err
		}
		items := make([]interface{}, 0, len(groups))
		for _, group := range groups {
			detail := map[string]interface{}{
				"GroupId":    group.ID,
				"Path":       group.Path,
				"GroupName":  group.GroupName,
				"Arn":        group.Arn,
				"CreateDate": group.CreateDate.Format(timeutils.ISO8601SimpleFormat),
			}
			detail["GroupPolicyList"] = buildInlinePolicyList(store, PrincipalTypeGroup, group.GroupName)
			detail["AttachedManagedPolicies"] = buildAttachedManagedPolicies(store, PrincipalTypeGroup, group.GroupName)
			items = append(items, detail)
		}
		sections = append(sections, section{
			key: "GroupDetailList", filter: "Group", items: items,
			marker: func(i int) string { return "group:" + groups[i].GroupName },
		})
	}

	if filters["Role"] {
		roles, err := listAllRoles(store, "")
		if err != nil {
			return nil, err
		}
		items := make([]interface{}, 0, len(roles))
		for _, role := range roles {
			detail := map[string]interface{}{
				"RoleId":                   role.ID,
				"Path":                     role.Path,
				"RoleName":                 role.RoleName,
				"Arn":                      role.Arn,
				"CreateDate":               role.CreateDate.Format(timeutils.ISO8601SimpleFormat),
				"AssumeRolePolicyDocument": role.AssumeRolePolicyDocument,
			}
			detail["RolePolicyList"] = buildInlinePolicyList(store, PrincipalTypeRole, role.RoleName)
			detail["AttachedManagedPolicies"] = buildAttachedManagedPolicies(store, PrincipalTypeRole, role.RoleName)
			items = append(items, detail)
		}
		sections = append(sections, section{
			key: "RoleDetailList", filter: "Role", items: items,
			marker: func(i int) string { return "role:" + roles[i].RoleName },
		})
	}

	if filters["LocalManagedPolicy"] {
		policies, err := listAllPolicies(store, "Local", "", false)
		if err != nil {
			return nil, err
		}
		items := make([]interface{}, 0, len(policies))
		for _, policy := range policies {
			items = append(items, map[string]interface{}{
				"PolicyName":       policy.PolicyName,
				"PolicyId":         policy.ID,
				"Arn":              policy.Arn,
				"Path":             policy.Path,
				"DefaultVersionId": policy.DefaultVersionId,
			})
		}
		sections = append(sections, section{
			key: "Policies", filter: "LocalManagedPolicy", items: items,
			marker: func(i int) string { return "policy:" + policies[i].Arn },
		})
	}

	// Apply pagination across all sections combined.
	// Marker format: "<sectionType>:<itemName>" (e.g. "user:admin", "role:MyRole").
	resp := map[string]interface{}{}
	skipUntilMarker := marker != ""
	count := 0
	isTruncated := false
	nextMarker := ""

	for _, sec := range sections {
		// Always emit the key so AWS SDK deserialisation gets an empty list
		// instead of null for skipped/partial sections.
		secItems := make([]interface{}, 0)

		for i, item := range sec.items {
			itemMarker := sec.marker(i)

			if skipUntilMarker {
				if itemMarker == marker {
					skipUntilMarker = false
				}
				continue
			}

			if count >= maxItems {
				isTruncated = true
				nextMarker = itemMarker
				break
			}
			secItems = append(secItems, item)
			count++
		}

		resp[sec.key] = secItems

		if isTruncated {
			break
		}
	}

	if skipUntilMarker {
		for _, sec := range sections {
			resp[sec.key] = []interface{}{}
		}
	}

	resp["IsTruncated"] = isTruncated
	if isTruncated && nextMarker != "" {
		resp["Marker"] = nextMarker
	}

	return resp, nil
}

// ---------------------------------------------------------------------------
// PolicyVersion core functions.
// ---------------------------------------------------------------------------

// getPolicyVersionCore retrieves a specific policy version.
func (s *IAMService) getPolicyVersionCore(store *iamstore.IAMStore, policyArn, versionId string) (*iamstore.PolicyVersion, error) {
	version, err := store.Policies().GetVersion(policyArn, versionId)
	if err != nil {
		return nil, NewNoSuchPolicyVersionError(versionId)
	}
	return version, nil
}

// deletePolicyVersionCore deletes a non-default policy version.
func (s *IAMService) deletePolicyVersionCore(store *iamstore.IAMStore, policyArn, versionId string) error {
	policy, err := store.Policies().Get(policyArn)
	if err != nil {
		return NewNoSuchPolicyError(policyArn)
	}

	if policy.DefaultVersionId == versionId {
		return NewDeleteConflictError("Cannot delete the default policy version.")
	}

	if err := store.Policies().DeleteVersion(policyArn, versionId); err != nil {
		return NewNoSuchPolicyVersionError(versionId)
	}
	return nil
}

// listPolicyVersionsCore returns a paginated list of policy versions.
func (s *IAMService) listPolicyVersionsCore(store *iamstore.IAMStore, policyArn, marker string, maxItems int) (*iamstore.PolicyVersionListResult, error) {
	if !store.Policies().Exists(policyArn) {
		return nil, NewNoSuchPolicyError(policyArn)
	}
	return store.Policies().ListVersions(policyArn, marker, maxItems)
}

// setDefaultPolicyVersionCore sets the default version for a policy.
func (s *IAMService) setDefaultPolicyVersionCore(store *iamstore.IAMStore, policyArn, versionId string) error {
	if !store.Policies().Exists(policyArn) {
		return NewNoSuchPolicyError(policyArn)
	}
	if err := store.Policies().SetDefaultVersion(policyArn, versionId); err != nil {
		return NewNoSuchPolicyVersionError(versionId)
	}
	return nil
}

// ListEntitiesForPolicyResult holds the aggregated entity lists returned
// by listEntitiesForPolicyCore.
type ListEntitiesForPolicyResult struct {
	PolicyUsers  []map[string]interface{}
	PolicyGroups []map[string]interface{}
	PolicyRoles  []map[string]interface{}
	IsTruncated  bool
	Marker       string
}

// listEntitiesForPolicyCore lists all principals that the specified
// managed policy is attached to, with optional entity-type filtering
// and cross-type pagination.
func (s *IAMService) listEntitiesForPolicyCore(store *iamstore.IAMStore, policyArn, entityFilter, marker string, maxItems int) (*ListEntitiesForPolicyResult, error) {
	if !store.Policies().Exists(policyArn) {
		return nil, NewNoSuchPolicyError(policyArn)
	}

	refs, err := store.AttachedPolicies().ListPrincipalsForPolicy(policyArn)
	if err != nil {
		return nil, err
	}

	type entityEntry struct {
		entityType string
		name       string
		data       map[string]interface{}
	}

	combined := make([]entityEntry, 0)

	for _, ref := range refs {
		switch ref.PrincipalType {
		case PrincipalTypeUser:
			if entityFilter != "" && entityFilter != "User" {
				continue
			}
			if user, err := store.Users().Get(ref.PrincipalName); err == nil {
				entry := map[string]interface{}{
					"UserName": user.UserName,
					"UserId":   user.ID,
					"Arn":      user.Arn,
				}
				combined = append(combined, entityEntry{"User", user.UserName, entry})
			}
		case PrincipalTypeGroup:
			if entityFilter != "" && entityFilter != "Group" {
				continue
			}
			if group, err := store.Groups().Get(ref.PrincipalName); err == nil {
				entry := map[string]interface{}{
					"GroupName": group.GroupName,
					"GroupId":   group.ID,
					"Arn":       group.Arn,
				}
				combined = append(combined, entityEntry{"Group", group.GroupName, entry})
			}
		case PrincipalTypeRole:
			if entityFilter != "" && entityFilter != "Role" {
				continue
			}
			if role, err := store.Roles().Get(ref.PrincipalName); err == nil {
				entry := map[string]interface{}{
					"RoleName": role.RoleName,
					"RoleId":   role.ID,
					"Arn":      role.Arn,
				}
				combined = append(combined, entityEntry{"Role", role.RoleName, entry})
			}
		}
	}

	paged := pagination.PaginateSlice(combined, marker, maxItems, func(item entityEntry) string {
		return item.entityType + ":" + item.name
	})

	result := &ListEntitiesForPolicyResult{
		PolicyUsers:  make([]map[string]interface{}, 0),
		PolicyGroups: make([]map[string]interface{}, 0),
		PolicyRoles:  make([]map[string]interface{}, 0),
		IsTruncated:  paged.IsTruncated,
		Marker:       paged.NextMarker,
	}

	for _, entry := range paged.Items {
		switch entry.entityType {
		case "User":
			result.PolicyUsers = append(result.PolicyUsers, entry.data)
		case "Group":
			result.PolicyGroups = append(result.PolicyGroups, entry.data)
		case "Role":
			result.PolicyRoles = append(result.PolicyRoles, entry.data)
		}
	}

	return result, nil
}
