package cognitoidentityprovider

import (
	"errors"

	"vorpalstacks/internal/common/request"
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
	if maxResults <= 0 || maxResults > listLimitMax {
		maxResults = listLimitMax
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

// userByAccessToken resolves the caller's user record from an access token:
// the token is required, validated, and the user must exist. It is the
// shared resolution step of the token-authenticated user operations.
func (s *CognitoService) userByAccessToken(reqCtx *request.RequestContext, accessToken string) (*cognitostore.User, error) {
	if accessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// getUserByAccessTokenCore resolves the caller for GetUser.
func (s *CognitoService) getUserByAccessTokenCore(reqCtx *request.RequestContext, accessToken string) (*cognitostore.User, error) {
	return s.userByAccessToken(reqCtx, accessToken)
}

// deleteUserByAccessTokenCore removes the caller's user record and every
// token minted for it.
func (s *CognitoService) deleteUserByAccessTokenCore(reqCtx *request.RequestContext, accessToken string) error {
	user, err := s.userByAccessToken(reqCtx, accessToken)
	if err != nil {
		return err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}
	if err := store.DeleteUser(user.UserPoolID, user.Username); err != nil {
		return ErrInternalError
	}
	if err := store.DeleteUserTokens(user.UserPoolID, user.ID); err != nil {
		return ErrInternalError
	}
	return nil
}

// deleteUserAttributesByAccessTokenCore removes the named attributes from
// the caller's record; a record without attributes is a no-op.
func (s *CognitoService) deleteUserAttributesByAccessTokenCore(reqCtx *request.RequestContext, accessToken string, attrNames []string) error {
	user, err := s.userByAccessToken(reqCtx, accessToken)
	if err != nil {
		return err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	if user.Attributes == nil {
		return nil
	}

	for _, name := range attrNames {
		delete(user.Attributes, name)
	}

	if err := store.UpdateUser(user); err != nil {
		return ErrInternalError
	}
	return nil
}

// updateUserAttributesByAccessTokenCore merges the supplied attributes into
// the caller's record.
func (s *CognitoService) updateUserAttributesByAccessTokenCore(reqCtx *request.RequestContext, accessToken string, attrs map[string]string) error {
	user, err := s.userByAccessToken(reqCtx, accessToken)
	if err != nil {
		return err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	if user.Attributes == nil {
		user.Attributes = make(map[string]string)
	}

	for k, v := range attrs {
		user.Attributes[k] = v
	}

	if err := store.UpdateUser(user); err != nil {
		return ErrInternalError
	}
	return nil
}
