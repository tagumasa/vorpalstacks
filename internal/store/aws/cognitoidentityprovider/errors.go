package cognitoidentityprovider

import (
	"errors"

	"vorpalstacks/internal/store/aws/common"
)

var (
	// ErrUserPoolNotFound is returned when the specified Cognito user pool
	// does not exist.
	ErrUserPoolNotFound = errors.New("user pool not found")

	// ErrUserPoolAlreadyExists is returned when attempting to create a user pool
	// that already exists.
	ErrUserPoolAlreadyExists = errors.New("user pool already exists")

	// ErrInvalidUserPoolName is returned when the user pool name is not valid.
	ErrInvalidUserPoolName = errors.New("invalid user pool name")

	// ErrUserNotFound is returned when the specified user does not exist.
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists is returned when attempting to create a user
	// that already exists.
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrInvalidUsername is returned when the username is not valid.
	ErrInvalidUsername = errors.New("invalid username")

	// ErrInvalidPassword is returned when the password is not valid.
	ErrInvalidPassword = errors.New("invalid password")

	// ErrGroupNotFound is returned when the specified group does not exist.
	ErrGroupNotFound = errors.New("group not found")

	// ErrGroupAlreadyExists is returned when attempting to create a group
	// that already exists.
	ErrGroupAlreadyExists = errors.New("group already exists")

	// ErrInvalidGroupName is returned when the group name is not valid.
	ErrInvalidGroupName = errors.New("invalid group name")

	// ErrUserNotInGroup is returned when attempting to remove a user from
	// a group they are not a member of.
	ErrUserNotInGroup = errors.New("user not in group")

	// ErrUserAlreadyInGroup is returned when attempting to add a user to
	// a group they are already a member of.
	ErrUserAlreadyInGroup = errors.New("user already in group")

	// ErrTokenNotFound is returned when the specified token does not exist.
	ErrTokenNotFound = errors.New("token not found")

	// ErrTokenExpired is returned when the token has expired.
	ErrTokenExpired = errors.New("token expired")

	// ErrInvalidToken is returned when the token is not valid.
	ErrInvalidToken = errors.New("invalid token")

	// ErrInvalidConfirmationCode is returned when the confirmation code is not valid.
	ErrInvalidConfirmationCode = errors.New("invalid confirmation code")

	// ErrUserNotConfirmed is returned when the user has not been confirmed.
	ErrUserNotConfirmed = errors.New("user not confirmed")

	// ErrUserAlreadyConfirmed is returned when the user is already confirmed.
	ErrUserAlreadyConfirmed = errors.New("user already confirmed")

	// ErrIncorrectPassword is returned when the password is incorrect.
	ErrIncorrectPassword = errors.New("incorrect password")

	// ErrPasswordPolicyViolation is returned when the password does not meet
	// the password policy requirements.
	ErrPasswordPolicyViolation = errors.New("password policy violation")

	// ErrStoreClosed is returned when the store has been closed.
	ErrStoreClosed = errors.New("store is closed")

	// ErrClientNotFound is returned when the specified client does not exist.
	ErrClientNotFound = errors.New("client not found")

	// ErrClientAlreadyExists is returned when attempting to create a client
	// that already exists.
	ErrClientAlreadyExists = errors.New("client already exists")

	// ErrInvalidParameter is returned when a parameter is not valid.
	ErrInvalidParameter = common.ErrInvalidParameter

	// ErrResourceAlreadyExists is returned when a resource already exists.
	ErrResourceAlreadyExists = errors.New("resource already exists")
)
