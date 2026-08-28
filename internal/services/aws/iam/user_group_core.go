// Transport-agnostic Core functions for IAM user-group membership:
// validation and store operations shared by the AWS-compatible HTTP API
// handlers and any admin surface (the xxxCore pattern).
package iam

import (
	"errors"

	iamstore "vorpalstacks/internal/store/aws/iam"
)

// UserGroupMembershipInput holds the parameters for the add-user-to-group
// and remove-user-from-group operations.
type UserGroupMembershipInput struct {
	UserName  string
	GroupName string
}

// addUserToGroupCore validates input and adds a user to a group. Adding an
// existing membership is an idempotent success; the documented error
// surface has no duplicate-membership error.
func (s *IAMService) addUserToGroupCore(store *iamstore.IAMStore, input *UserGroupMembershipInput) error {
	if input.UserName == "" {
		return NewValidationError("UserName")
	}
	if input.GroupName == "" {
		return NewValidationError("GroupName")
	}

	if !store.Users().Exists(input.UserName) {
		return NewNoSuchUserError(input.UserName)
	}

	if !store.Groups().Exists(input.GroupName) {
		return NewNoSuchGroupError(input.GroupName)
	}

	// The store sentinel covers concurrent duplicate additions racing past
	// the pre-check.
	if !store.UserGroups().IsUserInGroup(input.UserName, input.GroupName) {
		if err := store.UserGroups().AddUserToGroup(input.UserName, input.GroupName); err != nil {
			if errors.Is(err, iamstore.ErrUserAlreadyInGroup) {
				return nil
			}
			return err
		}
	}

	return nil
}

// removeUserFromGroupCore validates input and removes a user from a group.
func (s *IAMService) removeUserFromGroupCore(store *iamstore.IAMStore, input *UserGroupMembershipInput) error {
	if input.UserName == "" {
		return NewValidationError("UserName")
	}
	if input.GroupName == "" {
		return NewValidationError("GroupName")
	}

	if !store.Users().Exists(input.UserName) {
		return NewNoSuchUserError(input.UserName)
	}

	if !store.Groups().Exists(input.GroupName) {
		return NewNoSuchGroupError(input.GroupName)
	}

	if !store.UserGroups().IsUserInGroup(input.UserName, input.GroupName) {
		return NewUserNotInGroupError(input.UserName, input.GroupName)
	}

	return store.UserGroups().RemoveUserFromGroup(input.UserName, input.GroupName)
}
