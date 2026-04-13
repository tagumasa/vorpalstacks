package cognitoidentityprovider

import (
	"fmt"
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
func (s *CognitoService) CreateTokens(reqCtx *request.RequestContext, userPoolID, userID, clientID string) (accessToken, idToken, refreshToken string, expiresIn int64) {
	store, err := s.store(reqCtx)
	if err != nil {
		return "", "", "", 3600
	}
	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return "", "", "", 3600
	}

	privateKey, err := vsjwt.DecodePrivateKeyFromPEM(userPool.JwtPrivateKey)
	if err != nil {
		return "", "", "", 3600
	}

	issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", reqCtx.GetRegion(), userPoolID)
	jwtManager := vsjwt.NewManager(privateKey, userPool.JwtKeyID, issuer)

	user, _ := store.GetUserByID(userID)
	if user == nil {
		return "", "", "", 3600
	}

	attrs := userAttributesMap(user)
	tokGenResult, tokGenErr := invokePreTokenGeneration(
		reqCtx, s, TokenGenerationAuthentication, userPoolID, user.Username, clientID,
		userPool.LambdaConfig, attrs, user.Groups,
	)
	if tokGenErr != nil {
		return "", "", "", 3600
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

	accessToken, err = jwtManager.GenerateAccessToken(user, clientID, 3600)
	if err != nil {
		return "", "", "", 3600
	}

	at := cognitostore.NewAccessToken(userPoolID, user.ID, clientID, accessToken, time.Now().Add(time.Hour))
	if err := store.CreateAccessToken(at); err != nil {
		return "", "", "", 3600
	}

	idToken, err = jwtManager.GenerateIDToken(user, clientID, 3600)
	if err != nil {
		return "", "", "", 3600
	}

	refreshToken = jwtManager.GenerateRefreshToken()

	rt := cognitostore.NewRefreshToken(userPoolID, user.ID, clientID, "openid", time.Now().Add(30*24*time.Hour))
	rt.Token = refreshToken
	if err := store.CreateRefreshToken(rt); err != nil {
		return "", "", "", 3600
	}

	return accessToken, idToken, refreshToken, 3600
}

// ValidateAccessToken validates an access token and returns the user ID.
func (s *CognitoService) ValidateAccessToken(reqCtx *request.RequestContext, tokenString string) (string, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return "", err
	}
	userPools, err := store.ListUserPools()
	if err != nil || len(userPools) == 0 {
		return "", ErrResourceNotFound
	}

	for _, pool := range userPools {
		publicKey, err := vsjwt.DecodePublicKeyFromPEM(pool.JwtPublicKey)
		if err != nil {
			continue
		}

		issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", reqCtx.GetRegion(), pool.ID)
		jwtManager := vsjwt.NewManagerWithPublicKey(publicKey, pool.JwtKeyID, issuer)
		claims, err := jwtManager.ValidateToken(tokenString)
		if err == nil && claims.TokenUse == "access" {
			return claims.Subject, nil
		}
	}

	return "", ErrNotAuthorized
}
