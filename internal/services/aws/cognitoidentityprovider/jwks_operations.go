package cognitoidentityprovider

import (
	"fmt"

	"vorpalstacks/internal/services/aws/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	"vorpalstacks/pkg/vsjwt"
)

// GetJWKS returns the JSON Web Key Set (JWKS) for a user pool.
func (s *CognitoService) GetJWKS(reqCtx *request.RequestContext, userPoolID string) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	privateKey, err := vsjwt.DecodePrivateKeyFromPEM(userPool.JwtPrivateKey)
	if err != nil {
		return nil, ErrInternalError
	}

	issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", reqCtx.GetRegion(), userPoolID)
	jwtManager := vsjwt.NewManager(privateKey, userPool.JwtKeyID, issuer)
	return jwtManager.GetJWKS(), nil
}

// GetUserPoolByID returns a user pool by its ID.
func (s *CognitoService) GetUserPoolByID(reqCtx *request.RequestContext, userPoolID string) (*cognitostore.UserPool, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return store.GetUserPool(userPoolID)
}

// GetRegion returns the region from the request context.
func (s *CognitoService) GetRegion(reqCtx *request.RequestContext) string {
	return reqCtx.GetRegion()
}

// ListUserPoolsRaw returns a list of user pools.
func (s *CognitoService) ListUserPoolsRaw(reqCtx *request.RequestContext) ([]*cognitostore.UserPool, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return store.ListUserPools()
}
