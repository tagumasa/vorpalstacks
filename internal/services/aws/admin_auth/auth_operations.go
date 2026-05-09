package admin_auth

import (
	"context"

	pb "vorpalstacks/internal/pb/aws/admin_auth"
	"vorpalstacks/internal/pb/aws/common"
	iamstore "vorpalstacks/internal/store/aws/iam"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
	"vorpalstacks/pkg/vsjwt"

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

	jwtUser, err := s.resolveJWTUser(username)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

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

// resolveJWTUser builds a JWTUser for the given username.
// For the root user it returns a RootUserAdapter; for IAM users it looks up
// groups and attached policies from the IAM store.
func (s *AdminAuthService) resolveJWTUser(username string) (vsjwt.JWTUser, error) {
	if username == iamstore.RootUserName {
		return NewRootUserAdapter(s.accountID), nil
	}

	user, err := s.userReader.Get(username)
	if err != nil {
		return nil, err
	}

	groups, err := s.groupReader.ListGroupsForUser(username)
	if err != nil {
		return nil, err
	}

	policies, err := s.policyReader.ListAttachedPolicies("user", username)
	if err != nil {
		return nil, err
	}

	return NewUserAdapter(user, groups, policies), nil
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

	jwtUser, err := s.resolveJWTUser(rt.Username)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.deleteRefreshToken(refreshToken); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return s.generateLoginResponse(jwtUser, rt.Username)
}

// Logout invalidates the given access token by deleting any associated refresh tokens.
func (s *AdminAuthService) Logout(
	ctx context.Context,
	req *connect.Request[pb.LogoutRequest],
) (*connect.Response[common.Empty], error) {
	accessToken := req.Msg.GetAccessToken()
	if accessToken == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, ErrMissingAccessToken)
	}

	claims, err := s.tokenGenerator.ValidateToken(accessToken)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrInvalidToken)
	}

	if claims.TokenUse == "refresh" {
		if err := s.deleteRefreshToken(accessToken); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
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
			Arn:        arnutil.NewARNBuilder(s.accountID, "").IAM().Root(),
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
func (s *AdminAuthService) generateLoginResponse(jwtUser vsjwt.JWTUser, username string) (*connect.Response[pb.LoginResponse], error) {
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
