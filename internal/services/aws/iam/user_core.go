// Transport-agnostic Core functions for IAM users: validation and store
// operations shared by the AWS-compatible HTTP API handlers and the admin
// gRPC-Web handler (the xxxCore pattern).
package iam

import (
	"errors"

	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// CreateUserInput holds the parameters for creating an IAM user.
type CreateUserInput struct {
	UserName               string
	Path                   string
	PermissionsBoundaryArn string
	Tags                   []tags.Tag
}

// UpdateUserInput holds the parameters for renaming/updating an IAM user.
type UpdateUserInput struct {
	UserName    string
	NewPath     string
	NewUserName string
}

// DeleteUserInput holds the parameters for deleting an IAM user.
// When Cascade is true, all dependent resources are removed before the user
// record (admin handler behaviour).  When false, a DeleteConflict error is
// returned if any dependent resource remains (AWS API behaviour).
type DeleteUserInput struct {
	UserName string
	Cascade  bool
}

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
	if input.NewUserName != "" {
		if err := validateEntityName(input.NewUserName, "NewUserName"); err != nil {
			return nil, err
		}
	}

	if err := store.RenameUser(input.UserName, input.NewUserName, input.NewPath); err != nil {
		if errors.Is(err, iamstore.ErrUserAlreadyExists) {
			return nil, ErrUserAlreadyExists
		}
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
