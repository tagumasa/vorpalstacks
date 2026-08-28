// Transport-agnostic Core functions for IAM groups: validation and store
// operations shared by the AWS-compatible HTTP API handlers and the admin
// gRPC-Web handler (the xxxCore pattern).
package iam

import (
	"errors"

	iamstore "vorpalstacks/internal/store/aws/iam"
)

// CreateGroupInput holds the parameters for creating an IAM group.
type CreateGroupInput struct {
	GroupName string
	Path      string
}

// UpdateGroupInput holds the parameters for renaming/updating an IAM group.
type UpdateGroupInput struct {
	GroupName    string
	NewPath      string
	NewGroupName string
}

// DeleteGroupInput holds the parameters for deleting an IAM group.
// When Cascade is true, all dependent resources are removed before the group
// record (admin handler behaviour).  When false, a DeleteConflict error is
// returned if any dependent resource remains (AWS API behaviour).
type DeleteGroupInput struct {
	GroupName string
	Cascade   bool
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

// listUsersInGroupCore returns the member users of a group in membership
// order; users that disappear mid-walk are skipped.
func (s *IAMService) listUsersInGroupCore(store *iamstore.IAMStore, groupName string) ([]*iamstore.User, error) {
	userNames, err := store.UserGroups().ListUsersInGroup(groupName)
	if err != nil {
		return nil, err
	}
	users := make([]*iamstore.User, 0, len(userNames))
	for _, userName := range userNames {
		if user, err := store.Users().Get(userName); err == nil {
			users = append(users, user)
		}
	}
	return users, nil
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
	if input.NewGroupName != "" {
		if err := validateEntityName128(input.NewGroupName, "NewGroupName"); err != nil {
			return nil, err
		}
	}

	if err := store.RenameGroup(input.GroupName, input.NewGroupName, input.NewPath); err != nil {
		if errors.Is(err, iamstore.ErrGroupAlreadyExists) {
			return nil, ErrGroupAlreadyExists
		}
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
