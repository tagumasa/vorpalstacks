package cognitoidentityprovider

import (
	"context"
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

	maxResults := applyListLimitDefaults(in.MaxResults)

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
	_, err := s.setUserEnabledCore(region, userPoolID, username, true)
	return err
}

// adminDisableUserCore disables a user and revokes every token minted for
// them — the model documents AdminDisableUser as deactivating the profile
// "and revokes all access tokens for the user", so a disabled user's
// existing tokens and refresh grants stop working immediately.
func (s *CognitoService) adminDisableUserCore(region, userPoolID, username string) error {
	user, err := s.setUserEnabledCore(region, userPoolID, username, false)
	if err != nil {
		return err
	}
	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}
	if err := store.DeleteUserTokens(userPoolID, user.ID); err != nil {
		return ErrInternalError
	}
	return nil
}

// setUserEnabledCore sets the Enabled flag on a user, persists the change
// and returns the updated record.
func (s *CognitoService) setUserEnabledCore(region, userPoolID, username string, enabled bool) (*cognitostore.User, error) {
	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		if errors.Is(err, cognitostore.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user.Enabled = enabled
	if err := store.UpdateUser(user); err != nil {
		return nil, err
	}
	return user, nil
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
	// A disabled user keeps appearing in listings but their tokens no
	// longer authenticate — validateAccessTokenRecord already rejects the
	// token, so this guards callers that resolve the user by other means.
	if !user.Enabled {
		return nil, ErrNotAuthorized
	}
	return user, nil
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
	// Tokens are the dependent record: they go first, so a failed user
	// delete can never leave live credentials behind (the order the admin
	// delete path uses).
	if err := store.DeleteUserTokens(user.UserPoolID, user.ID); err != nil {
		return ErrInternalError
	}
	if err := store.DeleteUser(user.UserPoolID, user.Username); err != nil {
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
func (s *CognitoService) updateUserAttributesByAccessTokenCore(ctx context.Context, reqCtx *request.RequestContext, accessToken string, attrs map[string]string) error {
	user, err := s.userByAccessToken(reqCtx, accessToken)
	if err != nil {
		return err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	userPool, err := store.GetUserPool(user.UserPoolID)
	if err != nil {
		return ErrInternalError
	}

	if err := validateUserAttributesAgainstSchema(userPool, attrs, false); err != nil {
		return ErrInvalidParameter
	}

	if user.Attributes == nil {
		user.Attributes = make(map[string]string)
	}

	for k, v := range attrs {
		user.Attributes[k] = v
	}

	// An updated email/phone_number enters the verification round: the flag
	// drops, a VerifyUserAttribute code is minted and the CustomMessage
	// UpdateUserAttribute trigger fires.
	if err := issueAttributeUpdateVerification(ctx, s, user, userPool, attrs); err != nil {
		return err
	}

	if err := store.UpdateUser(user); err != nil {
		if errors.Is(err, cognitostore.ErrAliasExists) {
			return ErrAliasExists
		}
		return ErrInternalError
	}
	return nil
}
