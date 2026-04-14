package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
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

	userAttrs := parseUserAttributes(req)
	userAttrs["sub"] = ""

	preSignUpResult, err := invokePreSignUp(ctx, s, PreSignUpSignUp, targetPool.ID, username, clientID, targetPool.LambdaConfig, userAttrs)
	if err != nil {
		return nil, ErrInternalError
	}

	if preSignUpResult.UserAttributes != nil {
		userAttrs = preSignUpResult.UserAttributes
	}
	delete(userAttrs, "sub")

	user := cognitostore.NewUser(targetPool.ID, username)
	user.Attributes = userAttrs

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrInternalError
	}
	user.PasswordHash = string(hash)

	if preSignUpResult.AutoConfirmUser {
		user.UserStatus = "CONFIRMED"
	}

	if err := store.CreateUser(user); err != nil {
		return nil, ErrUserAlreadyExists
	}

	if preSignUpResult.AutoConfirmUser {
		attrs := userAttributesMap(user)
		if err := invokePostConfirmation(ctx, s, PostConfirmationConfirmSignUp, targetPool.ID, username, clientID, targetPool.LambdaConfig, attrs); err != nil {
			logs.Warn("PostConfirmation trigger failed", logs.Err(err))
		}
	} else {
		if _, err := invokeCustomMessage(ctx, s, CustomMessageSignUp, targetPool.ID, username, clientID, targetPool.LambdaConfig, "####", userAttributesMap(user)); err != nil {
			logs.Warn("CustomMessage trigger failed", logs.Err(err))
		}
	}

	return map[string]interface{}{
		"UserConfirmed":       preSignUpResult.AutoConfirmUser,
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

	attrs := userAttributesMap(user)
	if err := invokePostConfirmation(ctx, s, PostConfirmationConfirmSignUp, targetPool.ID, username, clientID, targetPool.LambdaConfig, attrs); err != nil {
		logs.Warn("PostConfirmation trigger failed", logs.Err(err))
	}

	return response.EmptyResponse(), nil
}
