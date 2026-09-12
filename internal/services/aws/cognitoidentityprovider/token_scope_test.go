package cognitoidentityprovider

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// jwtScopeClaim extracts the scope claim from a JWT's payload without
// verifying the signature — the tests assert claim content, not provenance.
func jwtScopeClaim(t *testing.T, token string) string {
	t.Helper()
	parts := splitToken(token)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("decode JWT payload: %v", err)
	}
	var claims struct {
		Scope string `json:"scope"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		t.Fatalf("parse JWT payload: %v", err)
	}
	return claims.Scope
}

func splitToken(token string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(token); i++ {
		if token[i] == '.' {
			parts = append(parts, token[start:i])
			start = i + 1
		}
	}
	return append(parts, token[start:])
}

// The access token's scope claim follows the session's origin: a token from
// API sign-in carries only aws.cognito.signin.user.admin; a hosted-UI
// session carries the client's granted OAuth scopes; a refresh reissues the
// session's original scope.
func TestAccessTokenScopeFollowsIssuanceOrigin(t *testing.T) {
	env := newChallengeTestEnv(t)

	// API sign-in — even though the client carries OAuth scopes of its own.
	if err := env.store.UpdateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID:         env.pool.ID,
		ClientID:           challengeTestClientID,
		ClientName:         "test-client",
		AllowedOAuthScopes: []string{"openid", "email"},
	}); err != nil {
		t.Fatal(err)
	}

	accessToken, _, _, _, err := env.svc.CreateTokens(env.reqCtx, env.pool.ID, env.user.ID, challengeTestClientID, TokenGenerationAuthentication, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := jwtScopeClaim(t, accessToken); got != apiSignInScope {
		t.Fatalf("API sign-in scope = %q, want %q", got, apiSignInScope)
	}

	// Hosted-UI session — the client's granted scopes.
	accessToken, _, _, _, err = env.svc.CreateTokensWithIDClaims(env.reqCtx, env.pool.ID, env.user.ID, challengeTestClientID, TokenGenerationHostedAuth, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := jwtScopeClaim(t, accessToken); got != "openid email" {
		t.Fatalf("hosted-UI scope = %q, want %q", got, "openid email")
	}

	// A client without configured scopes falls back to openid on the hosted
	// origin, and the API origin still carries the self-service scope alone.
	otherClient := &cognitostore.UserPoolClient{
		UserPoolID: env.pool.ID,
		ClientID:   "scope-less-client",
		ClientName: "scope-less",
	}
	if err := env.store.CreateUserPoolClient(otherClient); err != nil {
		t.Fatal(err)
	}
	accessToken, _, _, _, err = env.svc.CreateTokensWithIDClaims(env.reqCtx, env.pool.ID, env.user.ID, otherClient.ClientID, TokenGenerationHostedAuth, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := jwtScopeClaim(t, accessToken); got != "openid" {
		t.Fatalf("hosted-UI fallback scope = %q, want openid", got)
	}
	accessToken, _, _, _, err = env.svc.CreateTokens(env.reqCtx, env.pool.ID, env.user.ID, otherClient.ClientID, TokenGenerationAuthentication, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := jwtScopeClaim(t, accessToken); got != apiSignInScope {
		t.Fatalf("API sign-in scope (scope-less client) = %q, want %q", got, apiSignInScope)
	}

	// Refresh — the session's original scope is reissued verbatim.
	accessToken, _, _, _, err = env.svc.CreateTokensWithSessionScope(env.reqCtx, env.pool.ID, env.user.ID, challengeTestClientID, TokenGenerationRefreshTokens, "openid email", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := jwtScopeClaim(t, accessToken); got != "openid email" {
		t.Fatalf("refresh scope = %q, want the session's original %q", got, "openid email")
	}
}
