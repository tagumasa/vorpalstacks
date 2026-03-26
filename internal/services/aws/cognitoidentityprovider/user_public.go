package cognitoidentityprovider

import (
	"context"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"golang.org/x/crypto/bcrypt"
)

// SignUp registers a new user in the specified user pool.
func (s *CognitoService) SignUp(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	clientID := getClientId(req)
	username := getUsername(req)
	password := getPassword(req)

	if clientID == "" || username == "" || password == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	targetPool, err := store.GetUserPoolByClientID(clientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if err := validatePassword(password, targetPool.PasswordPolicy); err != nil {
		return nil, ErrPasswordPolicyViolation
	}

	user := cognitostore.NewUser(targetPool.ID, username)
	user.Attributes = parseUserAttributes(req)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrInternalError
	}
	user.PasswordHash = string(hash)

	if err := store.CreateUser(user); err != nil {
		return nil, ErrUserAlreadyExists
	}

	return map[string]interface{}{
		"UserConfirmed":       false,
		"CodeDeliveryDetails": map[string]interface{}{},
		"UserSub":             user.ID,
	}, nil
}

// ConfirmSignUp confirms a user's registration with the confirmation code.
func (s *CognitoService) ConfirmSignUp(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	clientID := getClientId(req)
	username := getUsername(req)
	confirmationCode := getConfirmationCode(req)

	if clientID == "" || username == "" || confirmationCode == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	targetPool, err := store.GetUserPoolByClientID(clientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	user, err := store.GetUser(targetPool.ID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.UserStatus == "CONFIRMED" {
		return nil, ErrUserAlreadyConfirmed
	}

	user.UserStatus = "CONFIRMED"
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}
