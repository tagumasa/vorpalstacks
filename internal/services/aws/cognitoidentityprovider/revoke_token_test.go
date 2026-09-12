package cognitoidentityprovider

import (
	"context"
	"errors"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// newRevokeTestEnv builds a service whose store holds a pool, a
// confidential client, a public client and a user with two live refresh
// tokens.
func newRevokeTestEnv(t *testing.T) (*CognitoService, *request.RequestContext, *cognitostore.UserPool, string, string) {
	t.Helper()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	reqCtx := request.NewRequestContext(context.Background(), mgr, "000000000000", "us-east-1")
	svc := NewCognitoService("000000000000", "us-east-1")
	svc.SetStorageManager(mgr)
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}

	pool, err := store.CreateUserPool(cognitostore.NewUserPool("revokepool", "us-east-1"))
	if err != nil {
		t.Fatal(err)
	}
	for _, client := range []*cognitostore.UserPoolClient{
		{UserPoolID: pool.ID, ClientID: "confidential-client", ClientName: "confidential", ClientSecret: "s3cret"},
		{UserPoolID: pool.ID, ClientID: "public-client", ClientName: "public"},
	} {
		if err := store.CreateUserPoolClient(client); err != nil {
			t.Fatal(err)
		}
	}
	user := cognitostore.NewUser(pool.ID, "revoke-user")
	user.UserStatus = "CONFIRMED"
	if err := store.CreateUser(user); err != nil {
		t.Fatal(err)
	}
	expires := time.Now().Add(time.Hour)
	for _, issued := range []struct {
		token  string
		client string
	}{
		{token: "confidential-refresh-token", client: "confidential-client"},
		{token: "public-refresh-token", client: "public-client"},
	} {
		rt := cognitostore.NewRefreshToken(pool.ID, user.ID, issued.client, "openid", expires)
		rt.Token = issued.token
		if err := store.CreateRefreshToken(rt); err != nil {
			t.Fatal(err)
		}
	}
	return svc, reqCtx, pool, "confidential-refresh-token", "public-refresh-token"
}

// A confidential client's refresh token cannot be revoked without the
// client secret: omitting the secret is a mismatch, and the refused request
// leaves the token in place. The correct secret revokes it.
func TestRevokeTokenConfidentialClientRequiresSecret(t *testing.T) {
	svc, reqCtx, _, confidentialToken, _ := newRevokeTestEnv(t)
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.revokeTokenCore(reqCtx, RevokeTokenInput{Token: confidentialToken, ClientID: "confidential-client"})
	if !errors.Is(err, ErrNotAuthorized) {
		t.Fatalf("revocation without the secret returned %v, want NotAuthorizedException", err)
	}
	if _, err := store.GetRefreshTokenByValue(confidentialToken); err != nil {
		t.Fatal("the refused request destroyed the refresh token")
	}

	_, err = svc.revokeTokenCore(reqCtx, RevokeTokenInput{Token: confidentialToken, ClientID: "confidential-client", ClientSecret: "wrong"})
	if !errors.Is(err, ErrNotAuthorized) {
		t.Fatalf("revocation with a wrong secret returned %v, want NotAuthorizedException", err)
	}
	if _, err := store.GetRefreshTokenByValue(confidentialToken); err != nil {
		t.Fatal("the refused request destroyed the refresh token")
	}

	if _, err := svc.revokeTokenCore(reqCtx, RevokeTokenInput{Token: confidentialToken, ClientID: "confidential-client", ClientSecret: "s3cret"}); err != nil {
		t.Fatalf("revocation with the correct secret failed: %v", err)
	}
	if _, err := store.GetRefreshTokenByValue(confidentialToken); err == nil {
		t.Fatal("the refresh token survived a correct-secret revocation")
	}
}

// A public client has no secret to present; its tokens revoke without one.
func TestRevokeTokenPublicClientRevokesWithoutSecret(t *testing.T) {
	svc, reqCtx, _, _, publicToken := newRevokeTestEnv(t)
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := svc.revokeTokenCore(reqCtx, RevokeTokenInput{Token: publicToken, ClientID: "public-client"}); err != nil {
		t.Fatalf("public-client revocation without a secret failed: %v", err)
	}
	if _, err := store.GetRefreshTokenByValue(publicToken); err == nil {
		t.Fatal("the refresh token survived a public-client revocation")
	}
}
