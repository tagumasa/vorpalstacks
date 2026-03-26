package admin_auth

import (
	"context"

	pb "vorpalstacks/internal/pb/aws/admin_auth"
	"vorpalstacks/internal/pb/aws/common"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"
)

// Login authenticates a user with username and password and returns tokens.
func (s *Service) Login(
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

	accessToken, err := s.tokenGenerator.GenerateAccessToken(jwtUser, DefaultClientID, AccessTokenDurationSec)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	idToken, err := s.tokenGenerator.GenerateIDToken(jwtUser, DefaultClientID, IDTokenDurationSec)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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

// RefreshToken refreshes access and ID tokens using a refresh token.
func (s *Service) RefreshToken(
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

	jwtUser := NewUserAdapter(user, groups, policies)

	newAccessToken, err := s.tokenGenerator.GenerateAccessToken(jwtUser, DefaultClientID, AccessTokenDurationSec)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	newIDToken, err := s.tokenGenerator.GenerateIDToken(jwtUser, DefaultClientID, IDTokenDurationSec)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	newRefreshToken := s.tokenGenerator.GenerateRefreshToken()

	if err := s.deleteRefreshToken(refreshToken); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if err := s.saveRefreshToken(newRefreshToken, rt.Username); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := connect.NewResponse(&pb.LoginResponse{
		AccessToken:      newAccessToken,
		RefreshToken:     newRefreshToken,
		IdToken:          newIDToken,
		ExpiresIn:        int32(AccessTokenDurationSec),
		TokenType:        "Bearer",
		RefreshExpiresIn: int32(RefreshTokenDurationSec),
	})

	return resp, nil
}

// Logout invalidates the given access token.
func (s *Service) Logout(
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
func (s *Service) GetCurrentUser(
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
