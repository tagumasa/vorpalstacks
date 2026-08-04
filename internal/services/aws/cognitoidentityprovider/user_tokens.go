package cognitoidentityprovider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	"vorpalstacks/pkg/vsjwt"
)

// GetUserPool returns the first user pool from the store for the given request context.
func (s *CognitoService) GetUserPool(reqCtx *request.RequestContext) *cognitostore.UserPool {
	store, err := s.store(reqCtx)
	if err != nil || store == nil {
		return nil
	}
	pools, _ := store.ListUserPools()
	if len(pools) > 0 {
		return pools[0]
	}
	return nil
}

// CreateTokens creates access, ID, and refresh tokens for the specified user.
// The triggerSource parameter controls which PreTokenGeneration Lambda event
// fires (TokenGenerationAuthentication, TokenGenerationRefreshTokens, or
// TokenGenerationHostedAuth). The clientMetadata parameter is forwarded to
// the trigger so Lambda can customise claims.
func (s *CognitoService) CreateTokens(reqCtx *request.RequestContext, userPoolID, userID, clientID, triggerSource string, clientMetadata map[string]string) (accessToken, idToken, refreshToken string, expiresIn int64, err error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return "", "", "", 0, err
	}
	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("user pool not found: %w", err)
	}

	privateKey, err := vsjwt.DecodePrivateKeyFromPEM(userPool.JwtPrivateKey)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to decode JWT private key: %w", err)
	}

	issuer := fmt.Sprintf("https://%s/%s", cognitoIdpHost(reqCtx.GetRegion()), userPoolID)
	jwtManager, err := vsjwt.NewManager(privateKey, userPool.JwtKeyID, issuer)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to create JWT manager: %w", err)
	}

	user, _ := store.GetUserByID(userID)
	if user == nil {
		return "", "", "", 0, fmt.Errorf("user not found: %s", userID)
	}

	attrs := userAttributesMap(user)
	tokGenResult, tokGenErr := invokePreTokenGeneration(
		reqCtx, s, triggerSource, userPoolID, user.Username, clientID,
		userPool.LambdaConfig, attrs, user.Groups, clientMetadata,
	)
	if tokGenErr != nil {
		return "", "", "", 0, fmt.Errorf("pre-token-generation trigger failed: %w", tokGenErr)
	}

	if tokGenResult != nil {
		for k, v := range tokGenResult.ClaimsToAddOrOverride {
			user.Attributes[k] = v
		}
		for _, k := range tokGenResult.ClaimsToSuppress {
			delete(user.Attributes, k)
		}
		if len(tokGenResult.GroupsToOverride) > 0 {
			user.Groups = tokGenResult.GroupsToOverride
		}
	}

	// Determine token validity from client configuration.
	atValidityMin := 60
	idValidityMin := 60
	rtValidityDays := 30
	if client, cerr := store.GetUserPoolClient(userPoolID, clientID); cerr == nil && client != nil {
		if client.AccessTokenValidity > 0 {
			atValidityMin = client.AccessTokenValidity
		}
		if client.IDTokenValidity > 0 {
			idValidityMin = client.IDTokenValidity
		}
		if client.RefreshTokenValidity > 0 {
			rtValidityDays = client.RefreshTokenValidity
		}
	}
	atValiditySec := int64(atValidityMin) * 60
	idValiditySec := int64(idValidityMin) * 60
	atExpiry := time.Duration(atValidityMin) * time.Minute
	idExpiry := time.Duration(idValidityMin) * time.Minute
	rtExpiry := time.Duration(rtValidityDays) * 24 * time.Hour

	accessToken, err = jwtManager.GenerateAccessToken(user, clientID, atValiditySec)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to generate access token: %w", err)
	}

	at := cognitostore.NewAccessToken(userPoolID, user.ID, clientID, "", time.Now().Add(atExpiry))
	at.Token = accessToken
	if err := store.CreateAccessToken(at); err != nil {
		return "", "", "", 0, fmt.Errorf("failed to store access token: %w", err)
	}

	idToken, err = jwtManager.GenerateIDToken(user, clientID, idValiditySec)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to generate ID token: %w", err)
	}

	// Store ID token so it can be validated and revoked.
	it := cognitostore.NewIDToken(userPoolID, user.ID, clientID, "", time.Now().Add(idExpiry), user.Groups)
	it.Token = idToken
	if err := store.CreateIDToken(it); err != nil {
		return "", "", "", 0, fmt.Errorf("failed to store ID token: %w", err)
	}

	refreshToken = jwtManager.GenerateRefreshToken()

	// Use the client's AllowedOAuthScopes if configured, falling back to
	// "openid" as the default scope.
	scope := "openid"
	if client, cerr := store.GetUserPoolClient(userPoolID, clientID); cerr == nil && client != nil && len(client.AllowedOAuthScopes) > 0 {
		scope = strings.Join(client.AllowedOAuthScopes, " ")
	}

	rt := cognitostore.NewRefreshToken(userPoolID, user.ID, clientID, scope, time.Now().Add(rtExpiry))
	rt.Token = refreshToken
	if err := store.CreateRefreshToken(rt); err != nil {
		return "", "", "", 0, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return accessToken, idToken, refreshToken, atValiditySec, nil
}

// ValidateAccessToken validates an access token and returns the user ID.
func (s *CognitoService) ValidateAccessToken(reqCtx *request.RequestContext, tokenString string) (string, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return "", err
	}

	_, err = store.GetAccessTokenByValue(tokenString)
	if err != nil {
		return "", ErrNotAuthorized
	}

	userPools, err := store.ListUserPools()
	if err != nil || len(userPools) == 0 {
		return "", ErrNotAuthorized
	}

	for _, pool := range userPools {
		publicKey, err := vsjwt.DecodePublicKeyFromPEM(pool.JwtPublicKey)
		if err != nil {
			continue
		}

		issuer := fmt.Sprintf("https://%s/%s", cognitoIdpHost(reqCtx.GetRegion()), pool.ID)
		jwtManager, err := vsjwt.NewManagerWithPublicKey(publicKey, pool.JwtKeyID, issuer)
		if err != nil {
			continue
		}
		claims, err := jwtManager.ValidateToken(tokenString)
		if err == nil && claims.TokenUse == "access" {
			return claims.Subject, nil
		}
	}

	return "", ErrNotAuthorized
}

// ValidateTokenForPool validates a Cognito access token for a specific
// user pool, designed for cross-service consumers (e.g. API Gateway
// COGNITO_USER_POOLS authorizer) via the eventbus CognitoTokenValidator
// interface.
func (s *CognitoService) ValidateTokenForPool(ctx context.Context, region, userPoolID, accessToken string) (string, error) {
	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return "", err
	}

	pool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return "", ErrResourceNotFound
	}

	// Verify the token has not been revoked.
	if _, err := store.GetAccessTokenByValue(accessToken); err != nil {
		return "", ErrNotAuthorized
	}

	publicKey, err := vsjwt.DecodePublicKeyFromPEM(pool.JwtPublicKey)
	if err != nil {
		return "", ErrNotAuthorized
	}

	issuer := fmt.Sprintf("https://%s/%s", cognitoIdpHost(region), pool.ID)
	jwtManager, err := vsjwt.NewManagerWithPublicKey(publicKey, pool.JwtKeyID, issuer)
	if err != nil {
		return "", ErrNotAuthorized
	}

	claims, err := jwtManager.ValidateToken(accessToken)
	if err != nil || claims.TokenUse != "access" {
		return "", ErrNotAuthorized
	}

	return claims.Subject, nil
}
