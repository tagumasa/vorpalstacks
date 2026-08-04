package cognitoidentityprovider

import (
	"errors"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Result structs
// ---------------------------------------------------------------------------

// ListUsersInput carries the pagination and filter parameters for ListUsers.
type ListUsersInput struct {
	UserPoolID string
	MaxResults int
	NextToken  string
	Filter     string
}

// ListUsersResult is the paginated result of ListUsers.
type ListUsersResult struct {
	Users     []*cognitostore.User
	NextToken string
}

// ---------------------------------------------------------------------------
// Core functions
// ---------------------------------------------------------------------------

// listUsersCore lists users in a user pool with pagination and optional filter.
func (s *CognitoService) listUsersCore(region string, in ListUsersInput) (*ListUsersResult, error) {
	if in.UserPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetUserPool(in.UserPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	maxResults := in.MaxResults
	if maxResults <= 0 || maxResults > 60 {
		maxResults = 60
	}

	var filterFunc func(*cognitostore.User) bool
	if in.Filter != "" {
		filterFunc = func(user *cognitostore.User) bool {
			return matchUserFilter(user, in.Filter)
		}
	}

	result, err := store.ListUsersPaginated(in.UserPoolID, storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   in.NextToken,
	}, filterFunc)
	if err != nil {
		return nil, err
	}

	return &ListUsersResult{
		Users:     result.Items,
		NextToken: result.NextMarker,
	}, nil
}

// adminGetUserCore retrieves a user by username. Returns the store-level User
// for format conversion by callers.
func (s *CognitoService) adminGetUserCore(region, userPoolID, username string) (*cognitostore.User, error) {
	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// adminDeleteUserCore deletes a user and their tokens.
func (s *CognitoService) adminDeleteUserCore(region, userPoolID, username string) error {
	if userPoolID == "" || username == "" {
		return ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return ErrUserNotFound
	}

	if err := store.DeleteUserTokens(userPoolID, user.ID); err != nil {
		return ErrInternalError
	}

	if err := store.DeleteUser(userPoolID, username); err != nil {
		return ErrUserNotFound
	}
	return nil
}

// adminEnableUserCore enables a user.
func (s *CognitoService) adminEnableUserCore(region, userPoolID, username string) error {
	return s.setUserEnabledCore(region, userPoolID, username, true)
}

// adminDisableUserCore disables a user.
func (s *CognitoService) adminDisableUserCore(region, userPoolID, username string) error {
	return s.setUserEnabledCore(region, userPoolID, username, false)
}

// setUserEnabledCore sets the Enabled flag on a user and persists the change.
func (s *CognitoService) setUserEnabledCore(region, userPoolID, username string, enabled bool) error {
	if userPoolID == "" || username == "" {
		return ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		if errors.Is(err, cognitostore.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	user.Enabled = enabled
	return store.UpdateUser(user)
}
