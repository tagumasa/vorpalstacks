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
		if err := attachUserPermissionsBoundaryCore(store, user, input.PermissionsBoundaryArn); err != nil {
			return nil, err
		}
	}

	return user, nil
}

// attachUserPermissionsBoundaryCore atomically sets a permissions boundary
// on an already-resolved IAM user: the shared attach sequence in
// attachPermissionsBoundaryCore, parameterised over the user's boundary
// field and the user store's persistence.
//
// Used by both createUserCore (when PermissionsBoundaryArn is specified
// at creation time) and the putUserPermissionsBoundaryCore operation Core.
// Consolidating the logic here prevents the create-time vs update-time
// drift that previously existed.
func attachUserPermissionsBoundaryCore(store *iamstore.IAMStore, user *iamstore.User, pbArn string) error {
	return attachPermissionsBoundaryCore(store, user, pbArn,
		func(u *iamstore.User) **iamstore.PermissionsBoundary { return &u.PermissionsBoundary },
		store.Users().Put,
	)
}

// putUserPermissionsBoundaryCore is the operation Core for the
// PutUserPermissionsBoundary API: it validates the two required members in
// the wire-contract order (UserName first, then PermissionsBoundary), so a
// request omitting both reports the user name, resolves the user, and
// attaches the boundary via attachUserPermissionsBoundaryCore.
func (s *IAMService) putUserPermissionsBoundaryCore(store *iamstore.IAMStore, userName, pbArn string) error {
	if userName == "" {
		return NewValidationError("UserName")
	}
	if pbArn == "" {
		return NewValidationError("PermissionsBoundary")
	}
	user, err := s.getUserCore(store, userName)
	if err != nil {
		return err
	}
	return attachUserPermissionsBoundaryCore(store, user, pbArn)
}

// getUserCore returns the IAM user with the given name; an empty name is
// rejected as a validation error (the data-plane GetUser path resolves the
// caller's own name before reaching this core).
func (s *IAMService) getUserCore(store *iamstore.IAMStore, userName string) (*iamstore.User, error) {
	if userName == "" {
		return nil, NewValidationError("UserName")
	}
	user, err := store.Users().Get(userName)
	if err != nil {
		return nil, storeReadError(err, iamstore.ErrUserNotFound, NewNoSuchUserError(userName))
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
	return updateEntityCore("User", input.UserName, input.NewPath, input.NewUserName,
		func() error { return store.RenameUser(input.UserName, input.NewUserName, input.NewPath) },
		iamstore.ErrUserAlreadyExists,
		func(name string) error { return NewUserAlreadyExistsError(name) },
		validateEntityName,
		store.Users().Get,
	)
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
	// The user is resolved rather than probed: the resolved read keeps an
	// outage from masquerading as a missing user, and the resolved record
	// carries the permissions boundary the deletion must decrement.
	user, err := store.Users().Get(input.UserName)
	if err != nil {
		return storeReadError(err, iamstore.ErrUserNotFound, NewNoSuchUserError(input.UserName))
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
	// removed (see decrementBoundaryUsageCount).
	decrementBoundaryUsageCount(store, user.PermissionsBoundary)

	return store.Users().Delete(input.UserName)
}

// deleteUserPermissionsBoundaryCore removes the permissions boundary
// from a user and decrements the usage count on the previously-bound
// policy.  Consolidates the logic that was previously inline in the
// HTTP API handler so that future admin-handler paths can delegate
// here as well.
func (s *IAMService) deleteUserPermissionsBoundaryCore(store *iamstore.IAMStore, userName string) error {
	if userName == "" {
		return NewValidationError("UserName")
	}
	user, err := store.Users().Get(userName)
	if err != nil {
		return storeReadError(err, iamstore.ErrUserNotFound, NewNoSuchUserError(userName))
	}
	return deletePermissionsBoundaryCore(store, user,
		func(u *iamstore.User) **iamstore.PermissionsBoundary { return &u.PermissionsBoundary },
		store.Users().Put,
	)
}

// listGroupsForUserCore retrieves the list of groups that a user
// belongs to.  Consolidates the store-direct logic so that both the
// HTTP API and future admin-handler paths delegate here. An empty user
// name is rejected as a validation error; a membership whose group
// cannot be read — an infrastructure fault, or a group that vanished
// between the membership listing and the read — fails the listing
// instead of silently omitting the group.
func (s *IAMService) listGroupsForUserCore(store *iamstore.IAMStore, userName string) ([]*iamstore.Group, error) {
	if userName == "" {
		return nil, NewValidationError("UserName")
	}
	// The user is resolved rather than probed: the resolved read keeps an
	// outage from masquerading as a missing user.
	if _, err := store.Users().Get(userName); err != nil {
		return nil, storeReadError(err, iamstore.ErrUserNotFound, NewNoSuchUserError(userName))
	}

	groupNames, err := store.UserGroups().ListGroupsForUser(userName)
	if err != nil {
		return nil, err
	}

	groups := make([]*iamstore.Group, 0, len(groupNames))
	for _, groupName := range groupNames {
		group, err := store.Groups().Get(groupName)
		if err != nil {
			return nil, storeListError(err)
		}
		groups = append(groups, group)
	}
	return groups, nil
}
