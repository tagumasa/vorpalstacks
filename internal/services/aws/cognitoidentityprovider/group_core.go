package cognitoidentityprovider

import (
	"errors"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Result structs
// ---------------------------------------------------------------------------

// ListGroupsInput carries the pagination parameters for ListGroups.
type ListGroupsInput struct {
	UserPoolID string
	MaxResults int
	NextToken  string
}

// ListGroupsResult is the paginated result of ListGroups.
type ListGroupsResult struct {
	Groups    []*cognitostore.Group
	NextToken string
}

// ListUsersInGroupInput carries pagination parameters for ListUsersInGroup.
type ListUsersInGroupInput struct {
	UserPoolID string
	GroupName  string
	MaxResults int
	NextToken  string
}

// ListUsersInGroupResult is the paginated result of ListUsersInGroup.
type ListUsersInGroupResult struct {
	Users     []*cognitostore.User
	NextToken string
}

// AdminListGroupsForUserInput carries pagination parameters.
type AdminListGroupsForUserInput struct {
	UserPoolID string
	Username   string
	MaxResults int
	NextToken  string
}

// AdminListGroupsForUserResult is the paginated result.
type AdminListGroupsForUserResult struct {
	Groups    []*cognitostore.Group
	NextToken string
}

// ---------------------------------------------------------------------------
// Core functions
// ---------------------------------------------------------------------------

// CreateGroupInput carries every field needed to create a Cognito group in a
// wire-protocol-independent format.
type CreateGroupInput struct {
	UserPoolID  string
	GroupName   string
	Description string
	RoleArn     string
	Precedence  *int
}

// createGroupFromInputCore creates a group from transport-agnostic input.
func (s *CognitoService) createGroupFromInputCore(region string, in CreateGroupInput) (*cognitostore.Group, error) {
	group := cognitostore.NewGroup(in.UserPoolID, in.GroupName)
	group.Description = in.Description
	group.RoleArn = in.RoleArn
	group.Precedence = in.Precedence
	return s.createGroupCore(region, group)
}

// createGroupCore creates a new Cognito group. The caller is responsible for
// any IAM role validation prior to calling this method.
func (s *CognitoService) createGroupCore(region string, group *cognitostore.Group) (*cognitostore.Group, error) {
	if group.UserPoolID == "" || group.Name == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetUserPool(group.UserPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	if err := store.CreateGroup(group); err != nil {
		return nil, ErrGroupAlreadyExists
	}
	return group, nil
}

// getGroupCore retrieves a group by name.
func (s *CognitoService) getGroupCore(region, userPoolID, groupName string) (*cognitostore.Group, error) {
	if userPoolID == "" || groupName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	group, err := store.GetGroup(userPoolID, groupName)
	if err != nil {
		return nil, ErrGroupNotFound
	}
	return group, nil
}

// deleteGroupCore deletes a group by name.
func (s *CognitoService) deleteGroupCore(region, userPoolID, groupName string) error {
	if userPoolID == "" || groupName == "" {
		return ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	if err := store.DeleteGroup(userPoolID, groupName); err != nil {
		return ErrGroupNotFound
	}
	return nil
}

// listGroupsCore lists groups in a user pool with pagination.
func (s *CognitoService) listGroupsCore(region string, in ListGroupsInput) (*ListGroupsResult, error) {
	if in.UserPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	maxResults := in.MaxResults
	if maxResults <= 0 || maxResults > listLimitMax {
		maxResults = listLimitMax
	}

	result, err := store.ListGroupsPaginated(in.UserPoolID, storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   in.NextToken,
	})
	if err != nil {
		return nil, err
	}

	return &ListGroupsResult{
		Groups:    result.Items,
		NextToken: result.NextMarker,
	}, nil
}

// updateGroupCore persists updates to a group. The caller is responsible for
// any IAM role validation prior to calling this method.
func (s *CognitoService) updateGroupCore(region string, group *cognitostore.Group) error {
	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	if err := store.UpdateGroup(group); err != nil {
		return ErrInternalError
	}
	return nil
}

// adminAddUserToGroupCore adds a user to a group.
func (s *CognitoService) adminAddUserToGroupCore(region, userPoolID, groupName, username string) error {
	if userPoolID == "" || groupName == "" || username == "" {
		return ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	if err := store.AddUserToGroup(userPoolID, groupName, username); err != nil {
		if errors.Is(err, cognitostore.ErrGroupNotFound) {
			return ErrGroupNotFound
		}
		if errors.Is(err, cognitostore.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

// adminRemoveUserFromGroupCore removes a user from a group.
func (s *CognitoService) adminRemoveUserFromGroupCore(region, userPoolID, groupName, username string) error {
	if userPoolID == "" || groupName == "" || username == "" {
		return ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	if err := store.RemoveUserFromGroup(userPoolID, groupName, username); err != nil {
		if errors.Is(err, cognitostore.ErrGroupNotFound) {
			return ErrGroupNotFound
		}
		if errors.Is(err, cognitostore.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

// listUsersInGroupCore lists users in a group with pagination.
func (s *CognitoService) listUsersInGroupCore(region string, in ListUsersInGroupInput) (*ListUsersInGroupResult, error) {
	if in.UserPoolID == "" || in.GroupName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	maxResults := in.MaxResults
	if maxResults <= 0 || maxResults > listLimitMax {
		maxResults = listLimitMax
	}

	result, err := store.ListUsersInGroupPaginated(in.UserPoolID, in.GroupName, storecommon.ListOptions{
		Marker:   in.NextToken,
		MaxItems: maxResults,
	})
	if err != nil {
		return nil, err
	}

	return &ListUsersInGroupResult{
		Users:     result.Items,
		NextToken: result.NextMarker,
	}, nil
}

// adminListGroupsForUserCore lists groups for a user with pagination.
func (s *CognitoService) adminListGroupsForUserCore(region string, in AdminListGroupsForUserInput) (*AdminListGroupsForUserResult, error) {
	if in.UserPoolID == "" || in.Username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	maxResults := in.MaxResults
	if maxResults <= 0 || maxResults > listLimitMax {
		maxResults = listLimitMax
	}

	result, err := store.ListGroupsForUserPaginated(in.UserPoolID, in.Username, storecommon.ListOptions{
		Marker:   in.NextToken,
		MaxItems: maxResults,
	})
	if err != nil {
		return nil, err
	}

	return &AdminListGroupsForUserResult{
		Groups:    result.Items,
		NextToken: result.NextMarker,
	}, nil
}
