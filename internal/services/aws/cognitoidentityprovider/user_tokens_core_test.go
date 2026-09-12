package cognitoidentityprovider

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	"vorpalstacks/pkg/vsjwt"
)

func newTokensTestEnv(t *testing.T) (*CognitoService, *request.RequestContext, cognitostore.CognitoStoreInterface, *cognitostore.UserPool) {
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
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("tokens-pool", "us-east-1"))
	if err != nil {
		t.Fatal(err)
	}
	return svc, reqCtx, store, pool
}

// Token issuance fails closed when the app client record cannot be read:
// defaulting to the standard validity periods would mint tokens that
// violate the client's configuration.
func TestCreateTokensFailsClosedWhenClientMissing(t *testing.T) {
	svc, reqCtx, store, pool := newTokensTestEnv(t)
	user := cognitostore.NewUser(pool.ID, "no-client-user")
	user.UserStatus = "CONFIRMED"
	if err := store.CreateUser(user); err != nil {
		t.Fatal(err)
	}

	_, _, _, _, err := svc.CreateTokens(reqCtx, pool.ID, user.ID, "absent-client", TokenGenerationAuthentication, nil)
	if err == nil {
		t.Fatal("expected token issuance to fail when the app client record cannot be read")
	}
}

// Token issuance fails closed when the user record cannot be read instead
// of misreporting an IO failure as user-not-found.
func TestCreateTokensFailsClosedWhenUserMissing(t *testing.T) {
	svc, reqCtx, store, pool := newTokensTestEnv(t)
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID: pool.ID,
		ClientID:   "tokens-client",
		ClientName: "tokens-client",
	}); err != nil {
		t.Fatal(err)
	}

	_, _, _, _, err := svc.CreateTokens(reqCtx, pool.ID, "absent-user-id", "tokens-client", TokenGenerationAuthentication, nil)
	if err == nil {
		t.Fatal("expected token issuance to fail when the user record cannot be read")
	}
}

// The access token's scope claim follows the session's origin: API sign-in
// carries only the self-service scope even on a client with granted OAuth
// scopes; the hosted-UI origin carries the granted scopes, on the claim and
// the stored access-token record alike.
func TestCreateTokensCarriesGrantedScope(t *testing.T) {
	svc, reqCtx, store, pool := newTokensTestEnv(t)
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID:         pool.ID,
		ClientID:           "scoped-client",
		ClientName:         "scoped-client",
		AllowedOAuthScopes: []string{"openid", "profile"},
	}); err != nil {
		t.Fatal(err)
	}
	user := cognitostore.NewUser(pool.ID, "scoped-user")
	user.UserStatus = "CONFIRMED"
	if err := store.CreateUser(user); err != nil {
		t.Fatal(err)
	}

	accessToken, _, _, _, err := svc.CreateTokens(reqCtx, pool.ID, user.ID, "scoped-client", TokenGenerationAuthentication, nil)
	if err != nil {
		t.Fatal(err)
	}

	publicKey, err := vsjwt.DecodePublicKeyFromPEM(pool.JwtPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	issuer := fmt.Sprintf("https://%s/%s", cognitoIdpHost("us-east-1"), pool.ID)
	validator, err := vsjwt.NewManagerWithPublicKey(publicKey, pool.JwtKeyID, issuer)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := validator.ValidateTokenForUse(accessToken, "access", "")
	if err != nil {
		t.Fatal(err)
	}
	if claims.Scope != apiSignInScope {
		t.Fatalf("API sign-in scope claim = %q, want %q", claims.Scope, apiSignInScope)
	}
	record, err := store.GetAccessTokenByValue(accessToken)
	if err != nil {
		t.Fatal(err)
	}
	if record.Scope != apiSignInScope {
		t.Fatalf("stored API sign-in scope = %q, want %q", record.Scope, apiSignInScope)
	}

	accessToken, _, _, _, err = svc.CreateTokensWithIDClaims(reqCtx, pool.ID, user.ID, "scoped-client", TokenGenerationHostedAuth, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	claims, err = validator.ValidateTokenForUse(accessToken, "access", "")
	if err != nil {
		t.Fatal(err)
	}
	if claims.Scope != "openid profile" {
		t.Fatalf("hosted-UI scope claim = %q, want %q", claims.Scope, "openid profile")
	}
	record, err = store.GetAccessTokenByValue(accessToken)
	if err != nil {
		t.Fatal(err)
	}
	if record.Scope != "openid profile" {
		t.Fatalf("stored access token scope = %q, want %q", record.Scope, "openid profile")
	}
}

// A hosted-UI session for a client without configured OAuth scopes grants
// the openid default; the API origin carries the self-service scope alone.
func TestCreateTokensDefaultsScopeToOpenID(t *testing.T) {
	svc, reqCtx, store, pool := newTokensTestEnv(t)
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID: pool.ID,
		ClientID:   "plain-client",
		ClientName: "plain-client",
	}); err != nil {
		t.Fatal(err)
	}
	user := cognitostore.NewUser(pool.ID, "plain-user")
	user.UserStatus = "CONFIRMED"
	if err := store.CreateUser(user); err != nil {
		t.Fatal(err)
	}

	accessToken, _, _, _, err := svc.CreateTokensWithIDClaims(reqCtx, pool.ID, user.ID, "plain-client", TokenGenerationHostedAuth, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	record, err := store.GetAccessTokenByValue(accessToken)
	if err != nil {
		t.Fatal(err)
	}
	if record.Scope != "openid" {
		t.Fatalf("stored hosted-UI scope = %q, want openid", record.Scope)
	}

	accessToken, _, _, _, err = svc.CreateTokens(reqCtx, pool.ID, user.ID, "plain-client", TokenGenerationAuthentication, nil)
	if err != nil {
		t.Fatal(err)
	}
	record, err = store.GetAccessTokenByValue(accessToken)
	if err != nil {
		t.Fatal(err)
	}
	if record.Scope != apiSignInScope {
		t.Fatalf("stored API sign-in scope = %q, want %q", record.Scope, apiSignInScope)
	}
}

// createUserIOFailureStore fails every user creation with a generic store
// error, standing in for a store whose writes fail.
type createUserIOFailureStore struct {
	cognitostore.CognitoStoreInterface
}

func (f *createUserIOFailureStore) CreateUser(u *cognitostore.User) error {
	return errors.New("simulated store IO failure")
}

// A store failure at user creation reports InternalError, never
// UsernameExists: a transient fault must not read as a taken username.
func TestSignUpCoreMapsStoreIOFailureToInternalError(t *testing.T) {
	svc, reqCtx, store, pool := newTokensTestEnv(t)
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID: pool.ID,
		ClientID:   "io-client",
		ClientName: "io-client",
	}); err != nil {
		t.Fatal(err)
	}
	svc.stores.Store("us-east-1", &createUserIOFailureStore{CognitoStoreInterface: store})

	_, err := svc.signUpCore(context.Background(), reqCtx, SignUpInput{
		ClientID:       "io-client",
		Username:       "io-user",
		Password:       "SchemaPass123!",
		UserAttributes: map[string]string{"email": "io@example.com"},
	})
	if !errors.Is(err, ErrInternalError) {
		t.Fatalf("sign-up under a store IO failure returned %v, want ErrInternalError", err)
	}
}

// The issued ID token carries the IAM roles of the user's groups — the
// claims identity-pool Token role mappings read — ordered by ascending
// precedence with the highest-precedence group's role as
// cognito:preferred_role, and the claim resolver returns them from the
// validated token alone. The standard members resolve too: the token
// carries the app client ID as its registered aud claim (the value the
// principal-tag default mappings tag sessions with) and the user's
// identifier as sub.
func TestIDTokenCarriesGroupRoleClaims(t *testing.T) {
	svc, reqCtx, store, pool := newTokensTestEnv(t)
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID: pool.ID,
		ClientID:   "roles-client",
		ClientName: "roles-client",
	}); err != nil {
		t.Fatal(err)
	}
	engineerRole := "arn:aws:iam::000000000000:role/engineer-role"
	staffRole := "arn:aws:iam::000000000000:role/staff-role"
	lowPrecedence, highPrecedence := 2, 1
	if err := store.CreateGroup(&cognitostore.Group{UserPoolID: pool.ID, Name: "staff", RoleArn: staffRole, Precedence: &lowPrecedence}); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateGroup(&cognitostore.Group{UserPoolID: pool.ID, Name: "engineers", RoleArn: engineerRole, Precedence: &highPrecedence}); err != nil {
		t.Fatal(err)
	}

	user := cognitostore.NewUser(pool.ID, "role-holder")
	user.Attributes = map[string]string{"sub": user.ID}
	if err := store.CreateUser(user); err != nil {
		t.Fatal(err)
	}
	// The staff membership lands first: the claim order must come from
	// precedence, not membership order.
	if err := store.AddUserToGroup(pool.ID, "staff", "role-holder"); err != nil {
		t.Fatal(err)
	}
	if err := store.AddUserToGroup(pool.ID, "engineers", "role-holder"); err != nil {
		t.Fatal(err)
	}

	accessToken, idToken, _, _, err := svc.CreateTokens(reqCtx, pool.ID, user.ID, "roles-client", TokenGenerationAuthentication, nil)
	if err != nil {
		t.Fatalf("create tokens: %v", err)
	}

	claims, err := svc.IDTokenClaimsForPool(context.Background(), "us-east-1", pool.ID, idToken)
	if err != nil {
		t.Fatalf("resolve ID token claims: %v", err)
	}
	if claims["cognito:roles"] != engineerRole+","+staffRole {
		t.Fatalf("cognito:roles = %q, want precedence-ordered %s,%s", claims["cognito:roles"], engineerRole, staffRole)
	}
	if claims["cognito:preferred_role"] != engineerRole {
		t.Fatalf("cognito:preferred_role = %q, want the precedence-1 group's role %s", claims["cognito:preferred_role"], engineerRole)
	}
	if claims["aud"] != "roles-client" {
		t.Fatalf("aud = %q, want the app client ID %q", claims["aud"], "roles-client")
	}
	if claims["sub"] != user.ID {
		t.Fatalf("sub = %q, want the user's identifier %q", claims["sub"], user.ID)
	}

	// A tampered token and an access token presented as an ID token both
	// fail validation.
	if _, err := svc.IDTokenClaimsForPool(context.Background(), "us-east-1", pool.ID, idToken+"x"); !errors.Is(err, ErrNotAuthorized) {
		t.Fatalf("tampered ID token returned %v, want NotAuthorizedException", err)
	}
	if _, err := svc.IDTokenClaimsForPool(context.Background(), "us-east-1", pool.ID, accessToken); !errors.Is(err, ErrNotAuthorized) {
		t.Fatalf("access token accepted as an ID token: %v", err)
	}
}

// The role-claim derivation fails closed on a group read failure: a store
// error during token minting must not silently drop cognito:roles (the
// downstream role resolution would degrade to its ambiguous fallback). A
// group that no longer exists is skipped — deletion cascades membership
// removal — and the remaining groups still contribute.
func TestGroupRoleClaimsFailClosedOnStoreErrors(t *testing.T) {
	_, _, store, pool := newTokensTestEnv(t)
	precedence := 1
	if err := store.CreateGroup(&cognitostore.Group{
		UserPoolID: pool.ID, Name: "admins", RoleArn: "arn:aws:iam::000000000000:role/Admins", Precedence: &precedence,
	}); err != nil {
		t.Fatal(err)
	}

	roles, preferred, err := groupRoleClaims(store, pool.ID, []string{"admins"})
	if err != nil {
		t.Fatalf("groupRoleClaims: %v", err)
	}
	if len(roles) != 1 || roles[0] != "arn:aws:iam::000000000000:role/Admins" || preferred != roles[0] {
		t.Fatalf("claims = %v/%s, want the admins role", roles, preferred)
	}

	failing := &groupReadFailingStore{CognitoStoreInterface: store}
	if _, _, err := groupRoleClaims(failing, pool.ID, []string{"admins"}); err == nil {
		t.Fatal("group read failure swallowed by the claim derivation")
	}

	roles, _, err = groupRoleClaims(store, pool.ID, []string{"gone-group", "admins"})
	if err != nil {
		t.Fatalf("absent group must be skipped, got %v", err)
	}
	if len(roles) != 1 {
		t.Fatalf("claims = %v, want the surviving group's role alone", roles)
	}
}

// groupReadFailingStore stands in for a store whose group reads fail.
type groupReadFailingStore struct {
	cognitostore.CognitoStoreInterface
}

func (f *groupReadFailingStore) GetGroup(string, string) (*cognitostore.Group, error) {
	return nil, errors.New("group store unavailable")
}
