package cognitoidentityprovider

import (
	"context"
	"errors"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
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

// UpdateGroupInput carries the update members of UpdateGroup. A nil
// Description keeps the stored description (a non-nil empty string clears
// it); a non-empty RoleArn replaces the stored role (an empty string cannot
// clear it); a nil Precedence keeps the stored value.
type UpdateGroupInput struct {
	UserPoolID  string
	GroupName   string
	Description *string
	RoleArn     string
	Precedence  *int
}

// createGroupValidatedCore is the single create path for Cognito groups,
// shared by the HTTP API and the admin console: the required members, the
// group-name pattern, the precedence range and the IAM role trust validation
// for the group role, followed by the shared persistence path. The IAM
// validator is obtained inside the Core, so both planes validate
// identically.
func (s *CognitoService) createGroupValidatedCore(ctx context.Context, region string, in CreateGroupInput) (*cognitostore.Group, error) {
	if in.UserPoolID == "" || in.GroupName == "" {
		return nil, ErrInvalidParameter
	}
	if !validateUsernamePattern(in.GroupName) {
		return nil, ErrInvalidParameter
	}

	group := cognitostore.NewGroup(in.UserPoolID, in.GroupName)
	group.Description = in.Description
	group.RoleArn = in.RoleArn
	group.Precedence = in.Precedence
	if in.Precedence != nil && !validatePrecedence(*in.Precedence) {
		return nil, ErrInvalidParameter
	}

	if in.RoleArn != "" {
		if err := s.validateGroupRoleArn(ctx, in.RoleArn); err != nil {
			return nil, err
		}
	}

	return s.createGroupCore(region, group)
}

// validateGroupRoleArn validates the group role's trust policy for the
// Cognito service principal. Without an injected role provider the trust
// check is skipped, matching the scheduler Core's nil handling.
func (s *CognitoService) validateGroupRoleArn(ctx context.Context, roleArn string) error {
	validator := s.iamValidator()
	if validator == nil {
		return nil
	}
	return validator.ValidateRoleForServiceWithErrors(ctx, roleArn, iam.ServicePrincipalCognito, &iam.RoleErrorFactories{
		RoleNotFoundError:        iam.NewCognitoRoleError,
		RoleCannotBeAssumedError: iam.NewCognitoRoleError,
		InvalidArnError:          iam.NewCognitoRoleError,
	})
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
		if errors.Is(err, cognitostore.ErrGroupNotFound) {
			return nil, ErrGroupNotFound
		}
		return nil, ErrInternalError
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

	maxResults := applyListLimitDefaults(in.MaxResults)

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

// updateGroupCore applies the update members onto the stored group and
// persists it, returning the updated group for response serialisation. The
// required members, the precedence range and the IAM role trust validation
// for the group role all live here, so both planes validate identically.
func (s *CognitoService) updateGroupCore(ctx context.Context, region string, in UpdateGroupInput) (*cognitostore.Group, error) {
	if in.UserPoolID == "" || in.GroupName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	// The existence read keeps the ResourceNotFoundException precedence and
	// the role trust validation stays outside the store's critical section;
	// the field application itself is a serialised read-modify-write so a
	// concurrent membership change cannot be overwritten by a stale member
	// list.
	if _, err := store.GetGroup(in.UserPoolID, in.GroupName); err != nil {
		if errors.Is(err, cognitostore.ErrGroupNotFound) {
			return nil, ErrGroupNotFound
		}
		return nil, ErrInternalError
	}
	if in.RoleArn != "" {
		if err := s.validateGroupRoleArn(ctx, in.RoleArn); err != nil {
			return nil, err
		}
	}

	var updated *cognitostore.Group
	err = store.UpdateGroupFunc(in.UserPoolID, in.GroupName, func(g *cognitostore.Group) error {
		if in.Description != nil {
			g.Description = *in.Description
		}
		if in.RoleArn != "" {
			g.RoleArn = in.RoleArn
		}
		if in.Precedence != nil {
			if !validatePrecedence(*in.Precedence) {
				return ErrInvalidParameter
			}
			g.Precedence = in.Precedence
		}
		updated = g
		return nil
	})
	if err != nil {
		if errors.Is(err, cognitostore.ErrGroupNotFound) {
			return nil, ErrGroupNotFound
		}
		var awsErr *awserrors.AWSError
		if errors.As(err, &awsErr) {
			return nil, err
		}
		return nil, ErrInternalError
	}
	return updated, nil
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

	maxResults := applyListLimitDefaults(in.MaxResults)

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

	maxResults := applyListLimitDefaults(in.MaxResults)

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
