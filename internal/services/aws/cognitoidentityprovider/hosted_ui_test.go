package cognitoidentityprovider

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"golang.org/x/crypto/bcrypt"
)

const (
	hostedUITestDomain    = "authpool"
	hostedUITestOtherHost = "otherpool"
	hostedUITestClientID  = "hosted-client"
	// hostedUITestSecondCLI is a second public client in the same pool as
	// hostedUITestClientID, so a token issued to one can be presented
	// under the other to pin the refresh token's client binding.
	hostedUITestSecondCLI = "hosted-client-2"
	hostedUITestOtherCLI  = "other-client"
	// hostedUITestSecretClient carries a secret, so every grant it makes on
	// the token endpoint must authenticate with it.
	hostedUITestSecretClient = "hosted-secret-client"
	hostedUITestSecret       = "UnitTestSecret1234567890abcdef"
	// hostedUITestM2MClient carries a secret and the machine-token flow
	// alone, so the client_credentials grant serves it.
	hostedUITestM2MClient = "hosted-m2m-client"
	hostedUITestM2MSecret = "M2MSecret1234567890abcdefABCDEF"
)

// hostedUITestEnv seeds two pools, each with its own domain and app client,
// plus a confirmed native-credential user in the first pool. The domains
// differ so tests can pin that a grant served on one domain never accepts
// another pool's tokens or clients.
type hostedUITestEnv struct {
	svc       *CognitoService
	reqCtx    *request.RequestContext
	store     cognitostore.CognitoStoreInterface
	pool      *cognitostore.UserPool
	otherPool *cognitostore.UserPool
}

func newHostedUITestEnv(t *testing.T) *hostedUITestEnv {
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

	pool, err := store.CreateUserPool(cognitostore.NewUserPool("hostedpool", "us-east-1"))
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SetUserPoolDomain(hostedUITestDomain, &cognitostore.UserPoolDomain{
		Domain:     hostedUITestDomain,
		UserPoolID: pool.ID,
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID:         pool.ID,
		ClientID:           hostedUITestClientID,
		ClientName:         "hosted",
		CallbackURLs:       []string{"https://app.example.com/callback"},
		LogoutURLs:         []string{"https://app.example.com/signed-out"},
		DefaultRedirectURI: "https://app.example.com/callback",
		AllowedOAuthFlows:  []string{"code", "implicit"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID:         pool.ID,
		ClientID:           hostedUITestSecondCLI,
		ClientName:         "hosted-second",
		CallbackURLs:       []string{"https://app.example.com/callback"},
		DefaultRedirectURI: "https://app.example.com/callback",
		AllowedOAuthFlows:  []string{"code"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID:         pool.ID,
		ClientID:           hostedUITestSecretClient,
		ClientName:         "hosted-confidential",
		ClientSecret:       hostedUITestSecret,
		CallbackURLs:       []string{"https://app.example.com/callback"},
		LogoutURLs:         []string{"https://app.example.com/signed-out"},
		DefaultRedirectURI: "https://app.example.com/callback",
		AllowedOAuthFlows:  []string{"code"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID:         pool.ID,
		ClientID:           hostedUITestM2MClient,
		ClientName:         "hosted-m2m",
		ClientSecret:       hostedUITestM2MSecret,
		ExplicitAuthFlows:  []string{clientTokenAuthFlow},
		CallbackURLs:       []string{"https://app.example.com/callback"},
		DefaultRedirectURI: "https://app.example.com/callback",
	}); err != nil {
		t.Fatal(err)
	}

	otherPool, err := store.CreateUserPool(cognitostore.NewUserPool("otherpool", "us-east-1"))
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SetUserPoolDomain(hostedUITestOtherHost, &cognitostore.UserPoolDomain{
		Domain:     hostedUITestOtherHost,
		UserPoolID: otherPool.ID,
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID: otherPool.ID,
		ClientID:   hostedUITestOtherCLI,
		ClientName: "other",
	}); err != nil {
		t.Fatal(err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("OldPass123!"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	user := cognitostore.NewUser(pool.ID, "alice")
	user.UserStatus = "CONFIRMED"
	user.PasswordHash = string(hash)
	saltHex, verifierHex, err := computeSrpVerifier(pool.ID, "alice", "OldPass123!")
	if err != nil {
		t.Fatal(err)
	}
	user.SrpSalt = saltHex
	user.SrpVerifier = verifierHex
	if err := store.CreateUser(user); err != nil {
		t.Fatal(err)
	}

	return &hostedUITestEnv{svc: svc, reqCtx: reqCtx, store: store, pool: pool, otherPool: otherPool}
}

// serve drives the hosted UI handler with the given Host label; the first
// Host component is the hosted-UI domain.
func (e *hostedUITestEnv) serve(t *testing.T, method, target, host, form string) (*http.Response, string) {
	t.Helper()
	var req *http.Request
	if form != "" {
		req = httptest.NewRequest(method, target, strings.NewReader(form))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req = httptest.NewRequest(method, target, nil)
	}
	req.Host = host + ".auth.example.com"
	rec := httptest.NewRecorder()
	e.svc.HostedUIHandler(rec, req)
	res := rec.Result()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	return res, string(body)
}

// authorizationCodeGrant drives the browser code flow end to end — login
// form submit, then the token-endpoint exchange — returning the minted
// token payload.
func (e *hostedUITestEnv) authorizationCodeGrant(t *testing.T, host, clientID string) map[string]interface{} {
	t.Helper()
	form := url.Values{}
	form.Set("username", "alice")
	form.Set("password", "OldPass123!")
	form.Set("client_id", clientID)
	form.Set("redirect_uri", "https://app.example.com/callback")
	form.Set("response_type", "code")
	res, _ := e.serve(t, http.MethodPost, "/login", host, form.Encode())
	if res.StatusCode != http.StatusFound {
		t.Fatalf("login submit status = %d, want 302", res.StatusCode)
	}
	location, err := url.Parse(res.Header.Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	code := location.Query().Get("code")
	if code == "" {
		t.Fatalf("login redirect carries no authorization code: %s", res.Header.Get("Location"))
	}

	exchange := url.Values{}
	exchange.Set("grant_type", "authorization_code")
	exchange.Set("client_id", clientID)
	exchange.Set("code", code)
	exchange.Set("redirect_uri", "https://app.example.com/callback")
	res, body := e.serve(t, http.MethodPost, "/oauth2/token", host, exchange.Encode())
	if res.StatusCode != http.StatusOK {
		t.Fatalf("authorization_code exchange status = %d: %s", res.StatusCode, body)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("token exchange response is not JSON: %s", body)
	}
	return payload
}

// The rendered login form must carry the authorisation request's client_id,
// not the pool ID — a form that submits the pool ID as its client can never
// validate.
func TestHostedUILoginFormCarriesClientIDFromQuery(t *testing.T) {
	env := newHostedUITestEnv(t)

	res, body := env.serve(t, http.MethodGet, "/login?client_id="+hostedUITestClientID, hostedUITestDomain, "")
	if res.StatusCode != http.StatusOK {
		t.Fatalf("login page status = %d", res.StatusCode)
	}
	if !strings.Contains(body, `name="client_id" value="`+hostedUITestClientID+`"`) {
		t.Fatalf("login form does not carry the query client_id:\n%s", body)
	}
	if strings.Contains(body, `value="`+env.pool.ID+`"`) {
		t.Fatalf("login form leaks the pool ID into a client_id field:\n%s", body)
	}
}

// The browser code flow: a submitted sign-in redirects with a single-use
// authorization code, and the token endpoint exchanges it for tokens bound
// to the client that received the code.
func TestHostedUIAuthorizationCodeGrantRoundTrip(t *testing.T) {
	env := newHostedUITestEnv(t)

	form := url.Values{}
	form.Set("username", "alice")
	form.Set("password", "OldPass123!")
	form.Set("client_id", hostedUITestClientID)
	form.Set("redirect_uri", "https://app.example.com/callback")
	form.Set("response_type", "code")
	res, _ := env.serve(t, http.MethodPost, "/login", hostedUITestDomain, form.Encode())
	if res.StatusCode != http.StatusFound {
		t.Fatalf("login submit status = %d, want 302", res.StatusCode)
	}
	location, err := url.Parse(res.Header.Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	code := location.Query().Get("code")
	if code == "" {
		t.Fatalf("login redirect carries no authorization code: %s", res.Header.Get("Location"))
	}

	exchange := url.Values{}
	exchange.Set("grant_type", "authorization_code")
	exchange.Set("client_id", hostedUITestClientID)
	exchange.Set("code", code)
	exchange.Set("redirect_uri", "https://app.example.com/callback")
	res, body := env.serve(t, http.MethodPost, "/oauth2/token", hostedUITestDomain, exchange.Encode())
	if res.StatusCode != http.StatusOK {
		t.Fatalf("authorization_code exchange status = %d: %s", res.StatusCode, body)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["access_token"] == "" || payload["id_token"] == "" || payload["refresh_token"] == "" {
		t.Fatalf("exchange response is missing tokens: %s", body)
	}

	// The code is single use: a second exchange must not mint again.
	res, body = env.serve(t, http.MethodPost, "/oauth2/token", hostedUITestDomain, exchange.Encode())
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("replayed code exchange status = %d, want 400", res.StatusCode)
	}
	if !strings.Contains(body, "invalid_grant") {
		t.Fatalf("replayed code exchange error = %s", body)
	}
}

// The token endpoint offers exactly the documented grants: the
// resource-owner password grant is not among them, so a request for it is
// refused with unsupported_grant_type.
func TestHostedUITokenEndpointRejectsPasswordGrant(t *testing.T) {
	env := newHostedUITestEnv(t)

	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("client_id", hostedUITestClientID)
	form.Set("username", "alice")
	form.Set("password", "OldPass123!")
	res, body := env.serve(t, http.MethodPost, "/oauth2/token", hostedUITestDomain, form.Encode())
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("password grant status = %d, want 400: %s", res.StatusCode, body)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("password grant response is not JSON: %s", body)
	}
	if payload["error"] != "unsupported_grant_type" {
		t.Fatalf("password grant error = %v", payload)
	}
}

// The client_credentials grant authenticates the confidential client (body
// secret or client_secret_basic) and returns the machine access token
// alone — no ID or refresh token — while wrong secrets and public clients
// are refused.
func TestHostedUIClientCredentialsGrant(t *testing.T) {
	env := newHostedUITestEnv(t)

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", hostedUITestM2MClient)
	form.Set("client_secret", hostedUITestM2MSecret)
	res, body := env.serve(t, http.MethodPost, "/oauth2/token", hostedUITestDomain, form.Encode())
	if res.StatusCode != http.StatusOK {
		t.Fatalf("client_credentials grant status = %d: %s", res.StatusCode, body)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("client_credentials response is not JSON: %s", body)
	}
	if payload["access_token"] == "" || payload["token_type"] != "Bearer" {
		t.Fatalf("client_credentials response is missing the access token: %v", payload)
	}
	if _, present := payload["refresh_token"]; present {
		t.Fatalf("client_credentials response carries a refresh token: %v", payload)
	}
	if _, present := payload["id_token"]; present {
		t.Fatalf("client_credentials response carries an ID token: %v", payload)
	}

	// client_secret_basic carries the same pair in the Authorization
	// header.
	basicForm := url.Values{}
	basicForm.Set("grant_type", "client_credentials")
	basic := base64.StdEncoding.EncodeToString([]byte(hostedUITestM2MClient + ":" + hostedUITestM2MSecret))
	req := httptest.NewRequest(http.MethodPost, "/oauth2/token", strings.NewReader(basicForm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+basic)
	req.Host = hostedUITestDomain + ".auth.example.com"
	rec := httptest.NewRecorder()
	env.svc.HostedUIHandler(rec, req)
	res2 := rec.Result()
	body2, err := io.ReadAll(res2.Body)
	if err != nil {
		t.Fatal(err)
	}
	if res2.StatusCode != http.StatusOK {
		t.Fatalf("client_secret_basic grant status = %d: %s", res2.StatusCode, body2)
	}
	var basicPayload map[string]interface{}
	if err := json.Unmarshal(body2, &basicPayload); err != nil {
		t.Fatalf("client_secret_basic response is not JSON: %s", body2)
	}
	if basicPayload["access_token"] == "" {
		t.Fatalf("client_secret_basic response carries no access token: %s", body2)
	}

	// A wrong secret fails client authentication.
	form.Set("client_secret", "WrongSecret1234567890abcdefXYZ")
	res, body = env.serve(t, http.MethodPost, "/oauth2/token", hostedUITestDomain, form.Encode())
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("wrong-secret grant status = %d, want 400", res.StatusCode)
	}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("wrong-secret response is not JSON: %s", body)
	}
	if payload["error"] != "invalid_client" {
		t.Fatalf("wrong-secret grant error = %v", payload)
	}

	// A public client holds no secret and cannot use the grant.
	publicForm := url.Values{}
	publicForm.Set("grant_type", "client_credentials")
	publicForm.Set("client_id", hostedUITestClientID)
	res, body = env.serve(t, http.MethodPost, "/oauth2/token", hostedUITestDomain, publicForm.Encode())
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("public-client grant status = %d, want 400: %s", res.StatusCode, body)
	}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("public-client response is not JSON: %s", body)
	}
	if payload["error"] != "invalid_client" {
		t.Fatalf("public-client grant error = %v", payload)
	}
}

// A refresh token is scoped to the pool and client that issued it: the
// hosted UI refresh grant refuses tokens from another pool (served on that
// pool's own domain) and clients other than the token's own, while the
// token's own pool and client still refresh.
func TestHostedUIRefreshGrantRejectsForeignPoolAndClient(t *testing.T) {
	env := newHostedUITestEnv(t)

	payload := env.authorizationCodeGrant(t, hostedUITestDomain, hostedUITestClientID)
	refreshToken, _ := payload["refresh_token"].(string)
	if refreshToken == "" {
		t.Fatalf("authorization_code grant returned no refresh token: %v", payload)
	}

	refreshForm := func(host, clientID, token string) (*http.Response, map[string]interface{}) {
		form := url.Values{}
		form.Set("grant_type", "refresh_token")
		form.Set("client_id", clientID)
		form.Set("refresh_token", token)
		res, body := env.serve(t, http.MethodPost, "/oauth2/token", host, form.Encode())
		var p map[string]interface{}
		if err := json.Unmarshal([]byte(body), &p); err != nil {
			t.Fatalf("refresh grant response is not JSON: %s", body)
		}
		return res, p
	}

	// The token's own client, but presented on another pool's domain: the
	// client does not exist on that pool's endpoint, so the client
	// authentication itself refuses.
	res, p := refreshForm(hostedUITestOtherHost, hostedUITestClientID, refreshToken)
	if res.StatusCode != http.StatusBadRequest || p["error"] != "invalid_client" {
		t.Fatalf("cross-pool refresh accepted: status=%d payload=%v", res.StatusCode, p)
	}

	// A second client of the same pool: it authenticates, but the token's
	// client binding must still refuse.
	res, p = refreshForm(hostedUITestDomain, hostedUITestSecondCLI, refreshToken)
	if res.StatusCode != http.StatusBadRequest || p["error"] != "invalid_grant" {
		t.Fatalf("cross-client refresh accepted: status=%d payload=%v", res.StatusCode, p)
	}

	// The token's own pool and client refresh, echoing the same refresh
	// token.
	res, p = refreshForm(hostedUITestDomain, hostedUITestClientID, refreshToken)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("own refresh rejected: status=%d payload=%v", res.StatusCode, p)
	}
	if p["access_token"] == "" || p["id_token"] == "" {
		t.Fatalf("refresh grant response is missing tokens: %v", p)
	}
	if p["refresh_token"] != refreshToken {
		t.Fatalf("refresh grant reissued the refresh token: %v", p)
	}
}

// /logout redirects only to a URI the client registered as a logout URL.
func TestHostedUILogoutValidatesRegisteredURI(t *testing.T) {
	env := newHostedUITestEnv(t)

	res, _ := env.serve(t, http.MethodGet,
		"/logout?client_id="+hostedUITestClientID+"&logout_uri=https://app.example.com/signed-out",
		hostedUITestDomain, "")
	if res.StatusCode != http.StatusFound {
		t.Fatalf("registered logout status = %d, want 302", res.StatusCode)
	}
	if got := res.Header.Get("Location"); got != "https://app.example.com/signed-out" {
		t.Fatalf("registered logout redirect = %s", got)
	}

	res, _ = env.serve(t, http.MethodGet,
		"/logout?client_id="+hostedUITestClientID+"&logout_uri=https://evil.example.net/steal",
		hostedUITestDomain, "")
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("unregistered logout_uri status = %d, want 400", res.StatusCode)
	}

	res, _ = env.serve(t, http.MethodGet, "/logout?logout_uri=https://app.example.com/signed-out", hostedUITestDomain, "")
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("client-less logout status = %d, want 400", res.StatusCode)
	}
}

// The sign-up page's POST creates the account through the shared
// registration Core: the user exists unconfirmed with a sign-up code, and a
// registration a second time is refused.
func TestHostedUISignUpCreatesUnconfirmedUser(t *testing.T) {
	env := newHostedUITestEnv(t)

	form := url.Values{}
	form.Set("username", "newuser")
	form.Set("password", "NewPass123!")
	form.Set("email", "newuser@example.com")
	form.Set("client_id", hostedUITestClientID)
	res, _ := env.serve(t, http.MethodPost, "/signup", hostedUITestDomain, form.Encode())
	if res.StatusCode != http.StatusFound {
		t.Fatalf("sign-up submit status = %d, want 302", res.StatusCode)
	}
	if got := res.Header.Get("Location"); !strings.HasPrefix(got, "/confirm?") {
		t.Fatalf("sign-up redirect = %s, want /confirm", got)
	}

	user, err := env.store.GetUser(env.pool.ID, "newuser")
	if err != nil {
		t.Fatalf("sign-up did not create the user: %v", err)
	}
	if user.UserStatus == "CONFIRMED" {
		t.Fatal("sign-up confirmed the user without a confirmation round")
	}
	if user.SignUpCode == "" {
		t.Fatal("sign-up left the user without a confirmation code")
	}
	if user.Attributes["email"] != "newuser@example.com" {
		t.Fatalf("sign-up dropped the email attribute: %v", user.Attributes)
	}

	// A duplicate registration is refused and does not overwrite the user.
	res, body := env.serve(t, http.MethodPost, "/signup", hostedUITestDomain, form.Encode())
	if res.StatusCode != http.StatusOK || !strings.Contains(body, "class=\"error\"") {
		t.Fatalf("duplicate sign-up did not re-render with an error: status=%d", res.StatusCode)
	}

	// A client from another pool must not register users through this
	// domain's sign-up page.
	form.Set("client_id", hostedUITestOtherCLI)
	res, _ = env.serve(t, http.MethodPost, "/signup", hostedUITestDomain, form.Encode())
	if res.StatusCode != http.StatusOK {
		t.Fatalf("foreign-client sign-up status = %d", res.StatusCode)
	}
	if _, err := env.store.GetUser(env.otherPool.ID, "newuser"); err == nil {
		t.Fatal("foreign-client sign-up created the user in another pool")
	}
}

// decodeIDTokenClaims decodes the payload segment of a JWT without
// verifying the signature — the tests inspect claims of tokens the service
// itself just minted.
func decodeIDTokenClaims(t *testing.T, token string) map[string]interface{} {
	t.Helper()
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("token is not a three-segment JWT: %s", token)
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("token payload is not base64url: %v", err)
	}
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		t.Fatalf("token payload is not JSON: %v", err)
	}
	return claims
}

// s256Challenge derives the PKCE challenge from a verifier exactly as a
// public client would.
func s256Challenge(t *testing.T, verifier string) string {
	t.Helper()
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// The authorize endpoint validates the request before any page renders:
// with a valid client and registered redirect URI, formatting failures
// (missing or invalid response_type, a challenge without its method, a
// non-S256 method) redirect the bare error to the redirect URI, while a
// failed client or redirect renders an error page. A valid request serves
// the login form carrying every flow member.
func TestHostedUIAuthorizeValidatesRequest(t *testing.T) {
	env := newHostedUITestEnv(t)

	const base = "/oauth2/authorize?client_id=" + hostedUITestClientID + "&redirect_uri=https://app.example.com/callback&state=st-123"

	for _, target := range []string{
		base + "&response_type=bogus",
		"/oauth2/authorize?client_id=" + hostedUITestClientID + "&redirect_uri=https://app.example.com/callback",
		base + "&response_type=code&code_challenge=some-challenge",
		base + "&response_type=code&code_challenge_method=plain&code_challenge=some-challenge",
	} {
		res, _ := env.serve(t, http.MethodGet, target, hostedUITestDomain, "")
		if res.StatusCode != http.StatusFound {
			t.Fatalf("authorize %s status = %d, want 302 error redirect", target, res.StatusCode)
		}
		if got := res.Header.Get("Location"); got != "https://app.example.com/callback?error=invalid_request" {
			t.Fatalf("authorize %s error redirect = %s", target, got)
		}
	}

	// An unknown client or unregistered redirect leaves no trustworthy
	// redirect target: the endpoint answers with an error page.
	res, _ := env.serve(t, http.MethodGet,
		"/oauth2/authorize?client_id=no-such-client&redirect_uri=https://app.example.com/callback&response_type=code",
		hostedUITestDomain, "")
	if res.StatusCode == http.StatusFound {
		t.Fatalf("authorize with unknown client redirected: %s", res.Header.Get("Location"))
	}
	res, _ = env.serve(t, http.MethodGet,
		"/oauth2/authorize?client_id="+hostedUITestClientID+"&redirect_uri=https://evil.example.net/steal&response_type=code",
		hostedUITestDomain, "")
	if res.StatusCode == http.StatusFound {
		t.Fatalf("authorize with unregistered redirect redirected: %s", res.Header.Get("Location"))
	}
	res, _ = env.serve(t, http.MethodGet,
		"/oauth2/authorize?client_id="+hostedUITestClientID+"&response_type=code",
		hostedUITestDomain, "")
	if res.StatusCode == http.StatusFound {
		t.Fatalf("authorize without redirect_uri redirected: %s", res.Header.Get("Location"))
	}

	// A valid request renders the login form with the flow members as
	// hidden fields.
	res, body := env.serve(t, http.MethodGet,
		base+"&response_type=code&nonce=n-42&code_challenge_method=S256&code_challenge=chal-42",
		hostedUITestDomain, "")
	if res.StatusCode != http.StatusOK {
		t.Fatalf("valid authorize status = %d", res.StatusCode)
	}
	for _, hidden := range []string{
		`name="client_id" value="` + hostedUITestClientID + `"`,
		`name="redirect_uri" value="https://app.example.com/callback"`,
		`name="response_type" value="code"`,
		`name="state" value="st-123"`,
		`name="nonce" value="n-42"`,
		`name="code_challenge" value="chal-42"`,
		`name="code_challenge_method" value="S256"`,
	} {
		if !strings.Contains(body, hidden) {
			t.Fatalf("login form missing hidden field %s:\n%s", hidden, body)
		}
	}
}

// The code flow round-trip echoes the state on the redirect and binds the
// nonce into the ID token the exchange issues.
func TestHostedUIAuthorizationCodeEchoesStateAndBindsNonce(t *testing.T) {
	env := newHostedUITestEnv(t)

	form := url.Values{}
	form.Set("username", "alice")
	form.Set("password", "OldPass123!")
	form.Set("client_id", hostedUITestClientID)
	form.Set("redirect_uri", "https://app.example.com/callback")
	form.Set("response_type", "code")
	form.Set("state", "st-abc")
	form.Set("nonce", "n-nonce")
	res, _ := env.serve(t, http.MethodPost, "/login", hostedUITestDomain, form.Encode())
	if res.StatusCode != http.StatusFound {
		t.Fatalf("login submit status = %d, want 302", res.StatusCode)
	}
	location, err := url.Parse(res.Header.Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	if location.Query().Get("state") != "st-abc" {
		t.Fatalf("code redirect dropped the state: %s", res.Header.Get("Location"))
	}
	code := location.Query().Get("code")
	if code == "" {
		t.Fatalf("code redirect carries no code: %s", res.Header.Get("Location"))
	}

	exchange := url.Values{}
	exchange.Set("grant_type", "authorization_code")
	exchange.Set("client_id", hostedUITestClientID)
	exchange.Set("code", code)
	exchange.Set("redirect_uri", "https://app.example.com/callback")
	res, body := env.serve(t, http.MethodPost, "/oauth2/token", hostedUITestDomain, exchange.Encode())
	if res.StatusCode != http.StatusOK {
		t.Fatalf("exchange status = %d: %s", res.StatusCode, body)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatal(err)
	}
	claims := decodeIDTokenClaims(t, payload["id_token"].(string))
	if claims["nonce"] != "n-nonce" {
		t.Fatalf("ID token nonce claim = %v, want n-nonce", claims["nonce"])
	}
}

// The implicit (token) flow returns the tokens and state in the fragment
// and binds the nonce into the ID token.
func TestHostedUIImplicitGrantEchoesStateAndBindsNonce(t *testing.T) {
	env := newHostedUITestEnv(t)

	form := url.Values{}
	form.Set("username", "alice")
	form.Set("password", "OldPass123!")
	form.Set("client_id", hostedUITestClientID)
	form.Set("redirect_uri", "https://app.example.com/callback")
	form.Set("response_type", "token")
	form.Set("state", "st-imp")
	form.Set("nonce", "n-imp")
	res, _ := env.serve(t, http.MethodPost, "/login", hostedUITestDomain, form.Encode())
	if res.StatusCode != http.StatusFound {
		t.Fatalf("implicit login status = %d, want 302", res.StatusCode)
	}
	location := res.Header.Get("Location")
	frag, err := url.ParseQuery(location[strings.Index(location, "#")+1:])
	if err != nil {
		t.Fatalf("implicit fragment is not a query: %s", location)
	}
	if frag.Get("state") != "st-imp" {
		t.Fatalf("implicit fragment dropped the state: %s", location)
	}
	if frag.Get("access_token") == "" || frag.Get("id_token") == "" {
		t.Fatalf("implicit fragment is missing tokens: %s", location)
	}
	claims := decodeIDTokenClaims(t, frag.Get("id_token"))
	if claims["nonce"] != "n-imp" {
		t.Fatalf("ID token nonce claim = %v, want n-imp", claims["nonce"])
	}
}

// A code issued against a PKCE challenge is redeemed only with the matching
// verifier: a missing verifier is invalid_request, a wrong one
// invalid_grant, and the correct verifier completes the exchange.
func TestHostedUIPKCE(t *testing.T) {
	env := newHostedUITestEnv(t)

	verifier := "a-pkce-verifier-with-enough-entropy-to-be-realistic-1234567890"
	form := url.Values{}
	form.Set("username", "alice")
	form.Set("password", "OldPass123!")
	form.Set("client_id", hostedUITestClientID)
	form.Set("redirect_uri", "https://app.example.com/callback")
	form.Set("response_type", "code")
	form.Set("code_challenge", s256Challenge(t, verifier))
	form.Set("code_challenge_method", "S256")
	res, _ := env.serve(t, http.MethodPost, "/login", hostedUITestDomain, form.Encode())
	if res.StatusCode != http.StatusFound {
		t.Fatalf("PKCE login submit status = %d, want 302", res.StatusCode)
	}
	location, err := url.Parse(res.Header.Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	code := location.Query().Get("code")
	if code == "" {
		t.Fatalf("PKCE login redirect carries no code: %s", res.Header.Get("Location"))
	}

	exchange := func(verifier string) (*http.Response, map[string]interface{}) {
		form := url.Values{}
		form.Set("grant_type", "authorization_code")
		form.Set("client_id", hostedUITestClientID)
		form.Set("code", code)
		form.Set("redirect_uri", "https://app.example.com/callback")
		if verifier != "" {
			form.Set("code_verifier", verifier)
		}
		res, body := env.serve(t, http.MethodPost, "/oauth2/token", hostedUITestDomain, form.Encode())
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(body), &payload); err != nil {
			t.Fatalf("PKCE exchange response is not JSON: %s", body)
		}
		return res, payload
	}

	// The code is consumed by the first attempt, so each outcome needs a
	// fresh code.
	freshCode := func(t *testing.T) string {
		t.Helper()
		res, _ := env.serve(t, http.MethodPost, "/login", hostedUITestDomain, form.Encode())
		loc, err := url.Parse(res.Header.Get("Location"))
		if err != nil {
			t.Fatal(err)
		}
		return loc.Query().Get("code")
	}

	code = freshCode(t)
	res, payload := exchange("")
	if res.StatusCode != http.StatusBadRequest || payload["error"] != "invalid_request" {
		t.Fatalf("verifier-less exchange = status %d %v, want 400 invalid_request", res.StatusCode, payload)
	}

	code = freshCode(t)
	res, payload = exchange("wrong-verifier-wrong-verifier-wrong-verifier-wrong")
	if res.StatusCode != http.StatusBadRequest || payload["error"] != "invalid_grant" {
		t.Fatalf("wrong-verifier exchange = status %d %v, want 400 invalid_grant", res.StatusCode, payload)
	}

	code = freshCode(t)
	res, payload = exchange(verifier)
	if res.StatusCode != http.StatusOK || payload["access_token"] == "" {
		t.Fatalf("correct-verifier exchange = status %d %v", res.StatusCode, payload)
	}
}

// The exchange requires the redirect_uri and enforces its binding: a
// mismatching URI is unauthorized_client with the documented
// invalid_redirect description.
func TestHostedUIAuthorizationCodeRedirectURIBinding(t *testing.T) {
	env := newHostedUITestEnv(t)

	loginForm := url.Values{}
	loginForm.Set("username", "alice")
	loginForm.Set("password", "OldPass123!")
	loginForm.Set("client_id", hostedUITestClientID)
	loginForm.Set("redirect_uri", "https://app.example.com/callback")
	loginForm.Set("response_type", "code")
	res, _ := env.serve(t, http.MethodPost, "/login", hostedUITestDomain, loginForm.Encode())
	location, err := url.Parse(res.Header.Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	code := location.Query().Get("code")
	if code == "" {
		t.Fatal("login redirect carries no code")
	}

	exchange := url.Values{}
	exchange.Set("grant_type", "authorization_code")
	exchange.Set("client_id", hostedUITestClientID)
	exchange.Set("code", code)
	res, body := env.serve(t, http.MethodPost, "/oauth2/token", hostedUITestDomain, exchange.Encode())
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("redirect-less exchange status = %d: %s", res.StatusCode, body)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["error"] != "invalid_request" {
		t.Fatalf("redirect-less exchange error = %v", payload)
	}

	// A fresh code for the mismatch case.
	res, _ = env.serve(t, http.MethodPost, "/login", hostedUITestDomain, loginForm.Encode())
	location, err = url.Parse(res.Header.Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	code = location.Query().Get("code")

	exchange.Set("code", code)
	exchange.Set("redirect_uri", "https://app.example.com/other")
	res, body = env.serve(t, http.MethodPost, "/oauth2/token", hostedUITestDomain, exchange.Encode())
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("mismatched-redirect exchange status = %d: %s", res.StatusCode, body)
	}
	payload = nil
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["error"] != "unauthorized_client" || payload["error_description"] != "invalid_redirect" {
		t.Fatalf("mismatched-redirect exchange error = %v", payload)
	}
}

// A confidential client authenticates on every token-endpoint grant with
// its secret — body (client_secret_post) or Basic header
// (client_secret_basic) — and a public client is unaffected.
func TestHostedUIConfidentialClientAuthenticatesOnAllGrants(t *testing.T) {
	env := newHostedUITestEnv(t)

	// Token exchange for the confidential client: the authorization code
	// redeems with the client secret, and the minted refresh token feeds
	// the refresh-grant checks below.
	codeLoginForm := url.Values{}
	codeLoginForm.Set("username", "alice")
	codeLoginForm.Set("password", "OldPass123!")
	codeLoginForm.Set("client_id", hostedUITestSecretClient)
	codeLoginForm.Set("redirect_uri", "https://app.example.com/callback")
	codeLoginForm.Set("response_type", "code")
	res, body := env.serve(t, http.MethodPost, "/login", hostedUITestDomain, codeLoginForm.Encode())
	if res.StatusCode != http.StatusFound {
		t.Fatalf("confidential login submit status = %d, want 302: %s", res.StatusCode, body)
	}
	redirectURL, perr := url.Parse(res.Header.Get("Location"))
	if perr != nil {
		t.Fatal(perr)
	}
	authCode := redirectURL.Query().Get("code")
	if authCode == "" {
		t.Fatal("confidential login redirect carries no code")
	}

	codeForm := url.Values{}
	codeForm.Set("grant_type", "authorization_code")
	codeForm.Set("client_id", hostedUITestSecretClient)
	codeForm.Set("code", authCode)
	codeForm.Set("redirect_uri", "https://app.example.com/callback")
	codeForm.Set("client_secret", hostedUITestSecret)
	res, body = env.serve(t, http.MethodPost, "/oauth2/token", hostedUITestDomain, codeForm.Encode())
	if res.StatusCode != http.StatusOK || !strings.Contains(body, "access_token") {
		t.Fatalf("secret-bearing code exchange = status %d %s", res.StatusCode, body)
	}
	var passwordPayload map[string]interface{}
	if err := json.Unmarshal([]byte(body), &passwordPayload); err != nil {
		t.Fatal(err)
	}

	// Refresh grant: the confidential client's refresh token needs the
	// secret too.
	refreshForm := func(secret string) (*http.Response, string) {
		form := url.Values{}
		form.Set("grant_type", "refresh_token")
		form.Set("client_id", hostedUITestSecretClient)
		form.Set("refresh_token", passwordPayload["refresh_token"].(string))
		if secret != "" {
			form.Set("client_secret", secret)
		}
		return env.serve(t, http.MethodPost, "/oauth2/token", hostedUITestDomain, form.Encode())
	}
	res, body = refreshForm("")
	if res.StatusCode != http.StatusBadRequest || !strings.Contains(body, "invalid_client") {
		t.Fatalf("secret-less refresh grant = status %d %s, want 400 invalid_client", res.StatusCode, body)
	}
	res, body = refreshForm(hostedUITestSecret)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("secret-bearing refresh grant = status %d %s", res.StatusCode, body)
	}

	// Authorization-code grant: the code issued to the confidential client
	// redeems only with the secret, presented here via the Basic header.
	loginForm := url.Values{}
	loginForm.Set("username", "alice")
	loginForm.Set("password", "OldPass123!")
	loginForm.Set("client_id", hostedUITestSecretClient)
	loginForm.Set("redirect_uri", "https://app.example.com/callback")
	loginForm.Set("response_type", "code")
	res, _ = env.serve(t, http.MethodPost, "/login", hostedUITestDomain, loginForm.Encode())
	location, err := url.Parse(res.Header.Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	code := location.Query().Get("code")
	if code == "" {
		t.Fatal("confidential login redirect carries no code")
	}

	exchange := url.Values{}
	exchange.Set("grant_type", "authorization_code")
	exchange.Set("client_id", hostedUITestSecretClient)
	exchange.Set("code", code)
	exchange.Set("redirect_uri", "https://app.example.com/callback")
	res, body = env.serve(t, http.MethodPost, "/oauth2/token", hostedUITestDomain, exchange.Encode())
	if res.StatusCode != http.StatusBadRequest || !strings.Contains(body, "invalid_client") {
		t.Fatalf("secret-less code exchange = status %d %s, want 400 invalid_client", res.StatusCode, body)
	}

	res, _ = env.serve(t, http.MethodPost, "/login", hostedUITestDomain, loginForm.Encode())
	location, err = url.Parse(res.Header.Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	code = location.Query().Get("code")
	exchange.Set("code", code)

	req := httptest.NewRequest(http.MethodPost, "/oauth2/token", strings.NewReader(exchange.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(
		[]byte(hostedUITestSecretClient+":"+hostedUITestSecret)))
	req.Host = hostedUITestDomain + ".auth.example.com"
	rec := httptest.NewRecorder()
	env.svc.HostedUIHandler(rec, req)
	res = rec.Result()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK || !strings.Contains(string(bodyBytes), "access_token") {
		t.Fatalf("Basic-auth code exchange = status %d %s", res.StatusCode, string(bodyBytes))
	}
}

// The hosted-UI sign-up flow completes through the confirmation page: the
// sign-up redirect carries the flow to /confirm, the page submits the
// delivered code through the shared ConfirmSignUp Core, a wrong code
// re-renders with an error, and the confirmed account continues to /login
// and can sign in.
func TestHostedUISignUpConfirmationRoundTrip(t *testing.T) {
	env := newHostedUITestEnv(t)

	form := url.Values{}
	form.Set("username", "confirmme")
	form.Set("password", "NewPass123!")
	form.Set("email", "confirmme@example.com")
	form.Set("client_id", hostedUITestClientID)
	form.Set("redirect_uri", "https://app.example.com/callback")
	form.Set("response_type", "code")
	form.Set("state", "st-signup")
	res, _ := env.serve(t, http.MethodPost, "/signup", hostedUITestDomain, form.Encode())
	if res.StatusCode != http.StatusFound {
		t.Fatalf("sign-up submit status = %d, want 302", res.StatusCode)
	}
	redirect := res.Header.Get("Location")
	if !strings.HasPrefix(redirect, "/confirm?") {
		t.Fatalf("sign-up redirect = %s, want /confirm", redirect)
	}
	if !strings.Contains(redirect, "state=st-signup") {
		t.Fatalf("sign-up redirect dropped the state: %s", redirect)
	}

	res, body := env.serve(t, http.MethodGet, redirect, hostedUITestDomain, "")
	if res.StatusCode != http.StatusOK || !strings.Contains(body, `name="confirmation_code"`) {
		t.Fatalf("confirmation page did not render a code form: status=%d", res.StatusCode)
	}

	confirmForm := url.Values{}
	confirmForm.Set("username", "confirmme")
	confirmForm.Set("confirmation_code", "00000000")
	confirmForm.Set("client_id", hostedUITestClientID)
	confirmForm.Set("redirect_uri", "https://app.example.com/callback")
	confirmForm.Set("response_type", "code")
	confirmForm.Set("state", "st-signup")
	res, body = env.serve(t, http.MethodPost, "/confirm", hostedUITestDomain, confirmForm.Encode())
	if res.StatusCode != http.StatusOK || !strings.Contains(body, "class=\"error\"") {
		t.Fatalf("wrong-code confirmation did not re-render with an error: status=%d", res.StatusCode)
	}
	user, err := env.store.GetUser(env.pool.ID, "confirmme")
	if err != nil {
		t.Fatal(err)
	}
	if user.UserStatus == "CONFIRMED" {
		t.Fatal("wrong confirmation code confirmed the user")
	}

	confirmForm.Set("confirmation_code", user.SignUpCode)
	res, _ = env.serve(t, http.MethodPost, "/confirm", hostedUITestDomain, confirmForm.Encode())
	if res.StatusCode != http.StatusFound {
		t.Fatalf("correct-code confirmation status = %d, want 302", res.StatusCode)
	}
	if got := res.Header.Get("Location"); !strings.HasPrefix(got, "/login?") || !strings.Contains(got, "state=st-signup") {
		t.Fatalf("confirmation redirect = %s, want /login with the carried state", got)
	}

	user, err = env.store.GetUser(env.pool.ID, "confirmme")
	if err != nil {
		t.Fatal(err)
	}
	if user.UserStatus != "CONFIRMED" {
		t.Fatalf("confirmed user status = %s", user.UserStatus)
	}
}

// The auth-code sweeper is a service-owned goroutine: after Close() it must
// exit, so constructing and closing services leaves no goroutine behind.
func TestHostedUIAuthCodeSweeperStopsOnClose(t *testing.T) {
	env := newHostedUITestEnv(t)

	before := runtime.NumGoroutine()

	form := url.Values{}
	form.Set("username", "alice")
	form.Set("password", "OldPass123!")
	form.Set("client_id", hostedUITestClientID)
	form.Set("redirect_uri", "https://app.example.com/callback")
	form.Set("response_type", "code")
	res, _ := env.serve(t, http.MethodPost, "/login", hostedUITestDomain, form.Encode())
	if res.StatusCode != http.StatusFound {
		t.Fatalf("login submit status = %d, want 302", res.StatusCode)
	}

	env.svc.Close()

	// Close() waits on the wait group, but give slow schedulers a brief
	// grace window before declaring the sweep goroutine leaked.
	deadline := time.Now().Add(2 * time.Second)
	for runtime.NumGoroutine() > before && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if after := runtime.NumGoroutine(); after > before {
		t.Fatalf("goroutines after Close: before=%d after=%d", before, after)
	}
}

// The OAuth endpoint grants are gated on the client's allowed OAuth flows:
// a client without the code flow gets unauthorized_client on the authorize
// endpoint (redirect form), the code grant and the refresh grant.
func TestHostedUIFlowGates(t *testing.T) {
	env := newHostedUITestEnv(t)

	// The M2M client holds a secret but no code flow.
	res, _ := env.serve(t, http.MethodGet,
		"/oauth2/authorize?client_id="+hostedUITestM2MClient+"&redirect_uri=https://app.example.com/callback&response_type=code",
		hostedUITestDomain, "")
	if res.StatusCode != http.StatusFound {
		t.Fatalf("authorize on a flows-less client status = %d, want 302", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.Contains(loc, "error=unauthorized_client") {
		t.Fatalf("authorize on a flows-less client redirected to %s, want error=unauthorized_client", loc)
	}

	// The code grant consumes the code but never redeems it for a
	// flows-less client.
	env.svc.authCodes.Store("gate-code", authCodeEntry{
		poolID:      env.pool.ID,
		userID:      "alice",
		clientID:    hostedUITestM2MClient,
		redirectURI: "https://app.example.com/callback",
		expires:     time.Now().Add(time.Minute),
	})
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", hostedUITestM2MClient)
	form.Set("client_secret", hostedUITestM2MSecret)
	form.Set("code", "gate-code")
	form.Set("redirect_uri", "https://app.example.com/callback")
	res, body := env.serve(t, http.MethodPost, "/oauth2/token", hostedUITestDomain, form.Encode())
	if res.StatusCode != http.StatusBadRequest || !strings.Contains(body, "unauthorized_client") {
		t.Fatalf("code grant on a flows-less client = %d %s, want 400 unauthorized_client", res.StatusCode, body)
	}

	// The refresh grant is gated the same way, with a refresh token the
	// client legitimately holds.
	alice, err := env.store.GetUser(env.pool.ID, "alice")
	if err != nil {
		t.Fatal(err)
	}
	_, _, refreshToken, _, err := env.svc.CreateTokens(env.reqCtx, env.pool.ID, alice.ID, hostedUITestM2MClient, TokenGenerationAuthentication, nil)
	if err != nil {
		t.Fatal(err)
	}
	form = url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", hostedUITestM2MClient)
	form.Set("client_secret", hostedUITestM2MSecret)
	form.Set("refresh_token", refreshToken)
	res, body = env.serve(t, http.MethodPost, "/oauth2/token", hostedUITestDomain, form.Encode())
	if res.StatusCode != http.StatusBadRequest || !strings.Contains(body, "unauthorized_client") {
		t.Fatalf("refresh grant on a flows-less client = %d %s, want 400 unauthorized_client", res.StatusCode, body)
	}
}

// The /logout redirect_uri form re-invites the user to the sign-in page
// with the original authorize parameters appended; response_type is
// required with it, the redirect_uri must be a registered callback URL,
// and a present logout_uri overrides it.
func TestHostedUILogoutRedirectURIForm(t *testing.T) {
	env := newHostedUITestEnv(t)

	res, _ := env.serve(t, http.MethodGet,
		"/logout?client_id="+hostedUITestClientID+"&redirect_uri=https://app.example.com/callback&response_type=code&state=st-1",
		hostedUITestDomain, "")
	if res.StatusCode != http.StatusFound {
		t.Fatalf("redirect_uri logout status = %d, want 302", res.StatusCode)
	}
	loc := res.Header.Get("Location")
	if !strings.HasPrefix(loc, "/login?") {
		t.Fatalf("redirect_uri logout target = %s, want the sign-in page", loc)
	}
	parsed, err := url.ParseRequestURI(loc)
	if err != nil {
		t.Fatal(err)
	}
	q := parsed.Query()
	if q.Get("redirect_uri") != "https://app.example.com/callback" || q.Get("response_type") != "code" || q.Get("state") != "st-1" || q.Get("client_id") != hostedUITestClientID {
		t.Fatalf("redirect_uri logout lost the original parameters: %s", loc)
	}

	res, _ = env.serve(t, http.MethodGet,
		"/logout?client_id="+hostedUITestClientID+"&redirect_uri=https://app.example.com/callback",
		hostedUITestDomain, "")
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("redirect_uri logout without response_type status = %d, want 400", res.StatusCode)
	}

	res, _ = env.serve(t, http.MethodGet,
		"/logout?client_id="+hostedUITestClientID+"&redirect_uri=https://evil.example.net/steal&response_type=code",
		hostedUITestDomain, "")
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("unregistered redirect_uri logout status = %d, want 400", res.StatusCode)
	}

	res, _ = env.serve(t, http.MethodGet,
		"/logout?client_id="+hostedUITestClientID+"&logout_uri=https://app.example.com/signed-out&redirect_uri=https://app.example.com/callback&response_type=code",
		hostedUITestDomain, "")
	if res.StatusCode != http.StatusFound || res.Header.Get("Location") != "https://app.example.com/signed-out" {
		t.Fatalf("logout_uri must override redirect_uri: %d %s", res.StatusCode, res.Header.Get("Location"))
	}
}
