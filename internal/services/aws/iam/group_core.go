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

// getGroupCore returns the IAM group with the given name; an empty name
// is rejected as a validation error.
func (s *IAMService) getGroupCore(store *iamstore.IAMStore, groupName string) (*iamstore.Group, error) {
	if groupName == "" {
		return nil, NewValidationError("GroupName")
	}
	group, err := store.Groups().Get(groupName)
	if err != nil {
		return nil, storeReadError(err, iamstore.ErrGroupNotFound, NewNoSuchGroupError(groupName))
	}
	return group, nil
}

// listGroupsCore returns a paginated list of IAM groups.
func (s *IAMService) listGroupsCore(store *iamstore.IAMStore, pathPrefix, marker string, maxItems int) (*iamstore.GroupListResult, error) {
	return store.Groups().List(pathPrefix, marker, maxItems)
}

// listUsersInGroupCore returns the member users of a group in membership
// order. A membership whose user cannot be read — an infrastructure fault,
// or a user that vanished between the membership listing and the read —
// fails the listing instead of silently omitting the user.
func (s *IAMService) listUsersInGroupCore(store *iamstore.IAMStore, groupName string) ([]*iamstore.User, error) {
	userNames, err := store.UserGroups().ListUsersInGroup(groupName)
	if err != nil {
		return nil, err
	}
	users := make([]*iamstore.User, 0, len(userNames))
	for _, userName := range userNames {
		user, err := store.Users().Get(userName)
		if err != nil {
			return nil, storeListError(err)
		}
		users = append(users, user)
	}
	return users, nil
}

// updateGroupCore validates input and renames/repaths an IAM group.
// Returns the updated group or an IAM-formatted error.
func (s *IAMService) updateGroupCore(store *iamstore.IAMStore, input *UpdateGroupInput) (*iamstore.Group, error) {
	return updateEntityCore("Group", input.GroupName, input.NewPath, input.NewGroupName,
		func() error { return store.RenameGroup(input.GroupName, input.NewGroupName, input.NewPath) },
		iamstore.ErrGroupAlreadyExists,
		func(name string) error { return NewGroupAlreadyExistsError(name) },
		validateEntityName128,
		store.Groups().Get,
	)
}

// deleteGroupCore validates input and deletes an IAM group.
// When input.Cascade is true, all dependent resources (inline/attached
// policies, user memberships) are removed first via cascadeDeleteGroup.
// When false, a DeleteConflict error is returned if any dependent resource
// remains.
func (s *IAMService) deleteGroupCore(store *iamstore.IAMStore, input *DeleteGroupInput) error {
	if input.GroupName == "" {
		return NewValidationError("GroupName")
	}
	// The group is resolved rather than probed: the resolved read keeps an
	// outage from masquerading as a missing group.
	if _, err := store.Groups().Get(input.GroupName); err != nil {
		return storeReadError(err, iamstore.ErrGroupNotFound, NewNoSuchGroupError(input.GroupName))
	}

	if input.Cascade {
		return cascadeDeleteGroup(store, input.GroupName)
	}

	// Conflict detection — AWS API path.
	userCount, err := store.UserGroups().CountUsersInGroup(input.GroupName)
	if err != nil {
		return err
	}
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
