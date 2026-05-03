package admin_auth

import (
	"context"
	"fmt"

	pb "vorpalstacks/internal/pb/aws/admin_auth"
	"vorpalstacks/internal/pb/aws/common"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"
)

// Login authenticates an IAM user with username and password and returns tokens.
func (s *AdminAuthService) Login(
	ctx context.Context,
	req *connect.Request[pb.LoginRequest],
) (*connect.Response[pb.LoginResponse], error) {
	username := req.Msg.GetUsername()
	password := req.Msg.GetPassword()

	if username == "" || password == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, ErrMissingCredentials)
	}

	valid, err := s.passwordVerifier.VerifyPassword(username, password)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if !valid {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrInvalidCredentials)
	}

	user, err := s.userReader.Get(username)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	groups, err := s.groupReader.ListGroupsForUser(username)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	policies, err := s.policyReader.ListAttachedPolicies("user", username)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	jwtUser := NewUserAdapter(user, groups, policies)

	return s.generateLoginResponse(jwtUser, username)
}

// LoginRoot authenticates the root user with password only.
// The root user is identified by the special login profile with UserName = iamstore.RootUserName.
func (s *AdminAuthService) LoginRoot(
	ctx context.Context,
	req *connect.Request[pb.LoginRootRequest],
) (*connect.Response[pb.LoginResponse], error) {
	password := req.Msg.GetPassword()

	if password == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, ErrMissingPassword)
	}

	valid, err := s.passwordVerifier.VerifyPassword(iamstore.RootUserName, password)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if !valid {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrInvalidCredentials)
	}

	jwtUser := NewRootUserAdapter(s.accountID)

	return s.generateLoginResponse(jwtUser, iamstore.RootUserName)
}

// IsRootInitialized returns whether the root user has been set up.
// It checks for the existence of a login profile for the root user.
func (s *AdminAuthService) IsRootInitialized(
	ctx context.Context,
	req *connect.Request[common.Empty],
) (*connect.Response[pb.IsRootInitializedResponse], error) {
	return connect.NewResponse(&pb.IsRootInitializedResponse{
		IsInitialized: s.loginProfileCheck.Exists(iamstore.RootUserName),
	}), nil
}

// InitialSetup creates the root user when the system has not yet been initialised.
// It creates a login profile (password) and an access key for the root user.
// Returns JWT tokens and the root access key credentials on success.
func (s *AdminAuthService) InitialSetup(
	ctx context.Context,
	req *connect.Request[pb.InitialSetupRequest],
) (*connect.Response[pb.InitialSetupResponse], error) {
	password := req.Msg.GetPassword()

	if password == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, ErrMissingPassword)
	}
	if len(password) < 8 {
		return nil, connect.NewError(connect.CodeInvalidArgument, ErrPasswordTooShort)
	}

	if s.loginProfileCheck.Exists(iamstore.RootUserName) {
		return nil, connect.NewError(connect.CodeFailedPrecondition, ErrAlreadyInitialised)
	}

	_, err := s.passwordCreator.Create(iamstore.RootUserName, password, false)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	accessKey, err := s.accessKeyCreator.Create(iamstore.RootUserName)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	jwtUser := NewRootUserAdapter(s.accountID)

	accessToken, err := s.tokenGenerator.GenerateAccessToken(jwtUser, DefaultClientID, AccessTokenDurationSec)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	idToken, err := s.tokenGenerator.GenerateIDToken(jwtUser, DefaultClientID, IDTokenDurationSec)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	refreshToken := s.tokenGenerator.GenerateRefreshToken()

	if err := s.saveRefreshToken(refreshToken, iamstore.RootUserName); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := connect.NewResponse(&pb.InitialSetupResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		IdToken:          idToken,
		ExpiresIn:        int32(AccessTokenDurationSec),
		TokenType:        "Bearer",
		RefreshExpiresIn: int32(RefreshTokenDurationSec),
		AccessKeyId:      accessKey.AccessKeyId,
		SecretAccessKey:  accessKey.SecretAccessKey,
	})

	return resp, nil
}

// RefreshToken refreshes access and ID tokens using a refresh token.
func (s *AdminAuthService) RefreshToken(
	ctx context.Context,
	req *connect.Request[pb.RefreshTokenRequest],
) (*connect.Response[pb.LoginResponse], error) {
	refreshToken := req.Msg.GetRefreshToken()
	if refreshToken == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, ErrMissingRefreshToken)
	}

	rt, err := s.getRefreshToken(refreshToken)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrInvalidRefreshToken)
	}

	var jwtUser interface {
		GetID() string
		GetUsername() string
		GetEmail() string
		GetGroups() []string
		GetCustomClaims() map[string]interface{}
	}

	if rt.Username == iamstore.RootUserName {
		jwtUser = NewRootUserAdapter(s.accountID)
	} else {
		user, err := s.userReader.Get(rt.Username)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		groups, err := s.groupReader.ListGroupsForUser(rt.Username)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		policies, err := s.policyReader.ListAttachedPolicies("user", rt.Username)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		jwtUser = NewUserAdapter(user, groups, policies)
	}

	// Generate tokens using the adapter that satisfies vsjwt.JWTUser
	var accessToken, idToken string
	switch u := jwtUser.(type) {
	case *RootUserAdapter:
		accessToken, err = s.tokenGenerator.GenerateAccessToken(u, DefaultClientID, AccessTokenDurationSec)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		idToken, err = s.tokenGenerator.GenerateIDToken(u, DefaultClientID, IDTokenDurationSec)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	case *UserAdapter:
		accessToken, err = s.tokenGenerator.GenerateAccessToken(u, DefaultClientID, AccessTokenDurationSec)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		idToken, err = s.tokenGenerator.GenerateIDToken(u, DefaultClientID, IDTokenDurationSec)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	newRefreshToken := s.tokenGenerator.GenerateRefreshToken()

	if err := s.deleteRefreshToken(refreshToken); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if err := s.saveRefreshToken(newRefreshToken, rt.Username); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := connect.NewResponse(&pb.LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     newRefreshToken,
		IdToken:          idToken,
		ExpiresIn:        int32(AccessTokenDurationSec),
		TokenType:        "Bearer",
		RefreshExpiresIn: int32(RefreshTokenDurationSec),
	})

	return resp, nil
}

// Logout invalidates the given access token.
func (s *AdminAuthService) Logout(
	ctx context.Context,
	req *connect.Request[pb.LogoutRequest],
) (*connect.Response[common.Empty], error) {
	accessToken := req.Msg.GetAccessToken()
	if accessToken == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, ErrMissingAccessToken)
	}

	_, err := s.tokenGenerator.ValidateToken(accessToken)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrInvalidToken)
	}

	return connect.NewResponse(&common.Empty{}), nil
}

// GetCurrentUser returns the current user information based on the access token.
func (s *AdminAuthService) GetCurrentUser(
	ctx context.Context,
	req *connect.Request[pb.GetCurrentUserRequest],
) (*connect.Response[pb.GetCurrentUserResponse], error) {
	accessToken := req.Msg.GetAccessToken()
	if accessToken == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, ErrMissingAccessToken)
	}

	claims, err := s.tokenGenerator.ValidateToken(accessToken)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrInvalidToken)
	}

	username := claims.Subject

	if username == iamstore.RootUserName {
		return connect.NewResponse(&pb.GetCurrentUserResponse{
			Username:   iamstore.RootUserName,
			Arn:        fmt.Sprintf("arn:aws:iam::%s:root", s.accountID),
			UserId:     iamstore.RootUserName,
			CreateDate: "",
		}), nil
	}

	user, err := s.userReader.Get(username)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, ErrUserNotFound)
	}

	groups, err := s.groupReader.ListGroupsForUser(username)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	policies, err := s.policyReader.ListAttachedPolicies("user", username)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := connect.NewResponse(&pb.GetCurrentUserResponse{
		Username:         user.UserName,
		Arn:              user.Arn,
		UserId:           user.ID,
		Groups:           groups,
		AttachedPolicies: policies,
		CreateDate:       user.CreateDate.Format(timeutils.ISO8601SimpleFormat),
	})

	return resp, nil
}

// generateLoginResponse is a helper that generates JWT tokens and builds a LoginResponse.
func (s *AdminAuthService) generateLoginResponse(jwtUser interface {
	GetID() string
	GetUsername() string
}, username string) (*connect.Response[pb.LoginResponse], error) {
	var accessToken, idToken string
	var err error

	switch u := jwtUser.(type) {
	case *RootUserAdapter:
		accessToken, err = s.tokenGenerator.GenerateAccessToken(u, DefaultClientID, AccessTokenDurationSec)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		idToken, err = s.tokenGenerator.GenerateIDToken(u, DefaultClientID, IDTokenDurationSec)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	case *UserAdapter:
		accessToken, err = s.tokenGenerator.GenerateAccessToken(u, DefaultClientID, AccessTokenDurationSec)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		idToken, err = s.tokenGenerator.GenerateIDToken(u, DefaultClientID, IDTokenDurationSec)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	refreshToken := s.tokenGenerator.GenerateRefreshToken()

	if err := s.saveRefreshToken(refreshToken, username); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := connect.NewResponse(&pb.LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		IdToken:          idToken,
		ExpiresIn:        int32(AccessTokenDurationSec),
		TokenType:        "Bearer",
		RefreshExpiresIn: int32(RefreshTokenDurationSec),
	})

	return resp, nil
}
