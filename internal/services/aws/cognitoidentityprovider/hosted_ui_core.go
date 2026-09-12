package cognitoidentityprovider

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// hostedUITokenResult carries the token set a hosted-UI grant mints.
type hostedUITokenResult struct {
	AccessToken  string
	IDToken      string
	RefreshToken string
	ExpiresIn    int64
}

// hostedUIAuthOutcome reports how a hosted-UI credential verification ended:
// either authentication completed and tokens were minted, or the pool's
// policy replaced token issuance with a challenge (NEW_PASSWORD_REQUIRED,
// MFA challenges) that the hosted UI pages cannot answer.
type hostedUIAuthOutcome struct {
	Tokens        *hostedUITokenResult
	ChallengeName string
}

// hostedUIError is the failure shape the hosted-UI endpoints report: an HTTP
// status, an OAuth2 error code (token-endpoint responses) and a
// human-readable description (browser pages). Both hosted-UI planes render
// the same value, so the Core decides the semantics once.
type hostedUIError struct {
	status int
	code   string
	desc   string
}

func (e *hostedUIError) Error() string { return e.code + ": " + e.desc }

func hostedUIServerError() *hostedUIError {
	return &hostedUIError{http.StatusInternalServerError, "server_error", "Internal error"}
}

// hostedUIFlowParams are the authorisation-request members the hosted-UI
// pages carry from form to form until the flow completes at a redirect:
// the client and its redirect target, the response mode, the OAuth state
// and OIDC nonce the client expects echoed back, and the PKCE challenge
// pair issued with the authorisation request.
type hostedUIFlowParams struct {
	ClientID            string
	RedirectURI         string
	ResponseType        string
	State               string
	Nonce               string
	CodeChallenge       string
	CodeChallengeMethod string
}

// queryValues and formValues abstract the request accessors behind the
// flow-parameter readers; url.Values and *http.Request satisfy them.
type queryValues interface {
	Get(string) string
}

type formValues interface {
	FormValue(string) string
}

func hostedUIFlowParamsFromQuery(r queryValues) hostedUIFlowParams {
	return hostedUIFlowParams{
		ClientID:            r.Get("client_id"),
		RedirectURI:         r.Get("redirect_uri"),
		ResponseType:        r.Get("response_type"),
		State:               r.Get("state"),
		Nonce:               r.Get("nonce"),
		CodeChallenge:       r.Get("code_challenge"),
		CodeChallengeMethod: r.Get("code_challenge_method"),
	}
}

func hostedUIFlowParamsFromForm(r formValues) hostedUIFlowParams {
	return hostedUIFlowParams{
		ClientID:            r.FormValue("client_id"),
		RedirectURI:         r.FormValue("redirect_uri"),
		ResponseType:        r.FormValue("response_type"),
		State:               r.FormValue("state"),
		Nonce:               r.FormValue("nonce"),
		CodeChallenge:       r.FormValue("code_challenge"),
		CodeChallengeMethod: r.FormValue("code_challenge_method"),
	}
}

// pageQuery assembles the flow members as page-link query values; only the
// members the request carried survive, so a plain /login link never grows
// empty parameters.
func (f hostedUIFlowParams) pageQuery() url.Values {
	q := url.Values{}
	setIfSet := func(key, value string) {
		if value != "" {
			q.Set(key, value)
		}
	}
	setIfSet("client_id", f.ClientID)
	setIfSet("redirect_uri", f.RedirectURI)
	setIfSet("response_type", f.ResponseType)
	setIfSet("state", f.State)
	setIfSet("nonce", f.Nonce)
	setIfSet("code_challenge", f.CodeChallenge)
	setIfSet("code_challenge_method", f.CodeChallengeMethod)
	return q
}

// authCodeEntry binds a single-use hosted-UI authorization code to the
// authenticated user, the pool behind the domain, the app client that
// requested the sign-in, and the authorisation-request members the token
// endpoint must re-verify at exchange: the redirect URI the code was issued
// for, the OIDC nonce to bind into the ID token, and the PKCE challenge.
type authCodeEntry struct {
	poolID        string
	userID        string
	clientID      string
	redirectURI   string
	nonce         string
	codeChallenge string
	expires       time.Time
}

func (s *CognitoService) startAuthCodeCleanup() {
	s.authCodeCleanupOnce.Do(func() {
		// The sweeper is a service-owned background worker: it registers
		// with the lifecycle wait group and exits on Close()'s context
		// cancellation like the import worker, so a closed service leaves
		// no goroutine behind.
		s.bgWg.Add(1)
		go func() {
			defer s.bgWg.Done()
			ticker := time.NewTicker(hostedUIAuthCodeSweepEvery)
			defer ticker.Stop()
			for {
				select {
				case <-s.bgCtx.Done():
					return
				case <-ticker.C:
					now := time.Now()
					s.authCodes.Range(func(key, value interface{}) bool {
						if entry, ok := value.(authCodeEntry); ok && now.After(entry.expires) {
							s.authCodes.Delete(key)
						}
						return true
					})
				}
			}
		}()
	})
}

// resolveDomainToPoolID maps a hosted-UI domain to its user pool. The stored
// domain entry carries the owning pool, so a single lookup answers.
func (s *CognitoService) resolveDomainToPoolID(domain string) (string, error) {
	reqCtx := request.NewRequestContext(context.Background(), s.storageManager, s.accountID, s.region)
	store, err := s.store(reqCtx)
	if err != nil {
		return "", err
	}
	entry, err := store.GetUserPoolDomain(domain)
	if err != nil {
		if !errors.Is(err, cognitostore.ErrUserPoolDomainNotFound) {
			return "", fmt.Errorf("failed to resolve domain %s: %w", domain, err)
		}
		return "", fmt.Errorf("domain %s not found", domain)
	}
	if entry == nil || entry.UserPoolID == "" {
		return "", fmt.Errorf("domain %s not found", domain)
	}
	return entry.UserPoolID, nil
}

// isRegisteredRedirectURI returns true if the redirectURI matches the client's
// DefaultRedirectURI or one of its registered CallbackURLs.
func isRegisteredRedirectURI(client *cognitostore.UserPoolClient, redirectURI string) bool {
	if client.DefaultRedirectURI != "" && client.DefaultRedirectURI == redirectURI {
		return true
	}
	for _, cb := range client.CallbackURLs {
		if cb == redirectURI {
			return true
		}
	}
	return false
}

// hostedUIClientCore resolves the app client a hosted-UI request names. The
// client is required and must belong to the domain's user pool, and a
// supplied redirect URI must be one the client registered — the
// authorisation-request validation the hosted UI fronts for.
func (s *CognitoService) hostedUIClientCore(reqCtx *request.RequestContext, poolID, clientID, redirectURI string) (*cognitostore.UserPoolClient, *hostedUIError) {
	if clientID == "" {
		return nil, &hostedUIError{http.StatusBadRequest, "invalid_request", "Missing client_id"}
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, hostedUIServerError()
	}
	client, err := store.GetUserPoolClient(poolID, clientID)
	if err != nil {
		if !errors.Is(err, cognitostore.ErrClientNotFound) {
			return nil, hostedUIServerError()
		}
		return nil, &hostedUIError{http.StatusUnauthorized, "invalid_client", "Invalid client_id"}
	}
	if redirectURI != "" && !isRegisteredRedirectURI(client, redirectURI) {
		return nil, &hostedUIError{http.StatusBadRequest, "invalid_request", "Invalid redirect_uri"}
	}
	return client, nil
}

// hostedUIAuthorizeCore validates an /oauth2/authorize request before any
// page renders. The client and its redirect URI are validated first: AWS
// documents that "If client_id and redirect_uri are valid, but the request
// parameters aren't formatted correctly, the authentication server
// redirects the error to the client's redirect_uri", so every formatting
// failure returns a non-empty oauthError the endpoint redirects with,
// while a client or redirect failure returns a hostedUIError to render as
// an error page (no redirect target is yet proven trustworthy).
// The response type must be the documented code/token pair; the PKCE pair
// must arrive complete (a challenge without its method is a documented
// invalid_request) and S256 is the only method the endpoint supports.
func (s *CognitoService) hostedUIAuthorizeCore(reqCtx *request.RequestContext, poolID string, f hostedUIFlowParams) (*cognitostore.UserPoolClient, string, *hostedUIError) {
	if f.RedirectURI == "" {
		return nil, "", &hostedUIError{http.StatusBadRequest, "invalid_request", "Missing redirect_uri"}
	}
	client, herr := s.hostedUIClientCore(reqCtx, poolID, f.ClientID, f.RedirectURI)
	if herr != nil {
		return nil, "", herr
	}
	if f.ResponseType != "code" && f.ResponseType != "token" {
		return client, "invalid_request", nil
	}
	// The requested grant must be one of the client's allowed OAuth flows;
	// the violation is the documented unauthorized_client ("Client is not
	// allowed for code grant flow or for refreshing tokens").
	requiredFlow := "code"
	if f.ResponseType == "token" {
		requiredFlow = "implicit"
	}
	if !clientAllowsOAuthFlow(client, requiredFlow) {
		return client, "unauthorized_client", nil
	}
	if f.CodeChallenge != "" && f.CodeChallengeMethod == "" {
		return client, "invalid_request", nil
	}
	if f.CodeChallengeMethod != "" && f.CodeChallengeMethod != "S256" {
		return client, "invalid_request", nil
	}
	return client, "", nil
}

// clientAllowsOAuthFlow reports whether the client's AllowedOAuthFlows
// permits the named flow (code, implicit or client_credentials) — the
// app-client opt-in every OAuth endpoint grant is gated on.
func clientAllowsOAuthFlow(client *cognitostore.UserPoolClient, flow string) bool {
	for _, f := range client.AllowedOAuthFlows {
		if f == flow {
			return true
		}
	}
	return false
}

// hostedUIClientAuthenticationCore authenticates the app client on a token
// endpoint grant. The client must belong to the domain's pool, and a client
// holding a secret — primary or a rotated descriptor — must present it; a
// client without any secret is public and authenticates by ID alone.
func (s *CognitoService) hostedUIClientAuthenticationCore(reqCtx *request.RequestContext, poolID, clientID, clientSecret string) (*cognitostore.UserPoolClient, *hostedUIError) {
	if clientID == "" {
		return nil, &hostedUIError{http.StatusBadRequest, "invalid_request", "Missing client_id"}
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, hostedUIServerError()
	}
	client, err := store.GetUserPoolClient(poolID, clientID)
	if err != nil {
		if !errors.Is(err, cognitostore.ErrClientNotFound) {
			return nil, hostedUIServerError()
		}
		return nil, &hostedUIError{http.StatusBadRequest, "invalid_client", "Client authentication failed."}
	}
	if client.ClientSecret != "" || len(client.ClientSecrets) > 0 {
		if !clientSecretMatches(client, clientSecret) {
			return nil, &hostedUIError{http.StatusBadRequest, "invalid_client", "Client authentication failed."}
		}
	}
	return client, nil
}

// hostedUIClientCredentialsResult carries the machine token minted for the
// client_credentials grant.
type hostedUIClientCredentialsResult struct {
	accessToken string
	expiresIn   int64
}

// hostedUIClientCredentialsTokenCore serves the documented client_credentials
// grant: the confidential client authenticates and the machine token issues
// through the API-plane client-token core, so the grant's flow gating, scope
// adjudication and trigger share one path with the API operation.
func (s *CognitoService) hostedUIClientCredentialsTokenCore(ctx context.Context, reqCtx *request.RequestContext, poolID, clientID, clientSecret, scopeParam string) (*hostedUIClientCredentialsResult, *hostedUIError) {
	client, herr := s.hostedUIClientAuthenticationCore(reqCtx, poolID, clientID, clientSecret)
	if herr != nil {
		return nil, herr
	}
	// The grant is a confidential-client grant: a public client holds no
	// secret to authenticate with, whatever the body carries.
	if client.ClientSecret == "" && len(client.ClientSecrets) == 0 {
		return nil, &hostedUIError{http.StatusBadRequest, "invalid_client", "Client authentication failed."}
	}

	out, err := s.getClientTokenCore(ctx, reqCtx, GetClientTokenInput{
		ClientID: clientID,
		Secret:   clientSecret,
		Scopes:   strings.Fields(scopeParam),
	})
	if err != nil {
		switch {
		case errors.Is(err, ErrOperationNotEnabled):
			return nil, &hostedUIError{http.StatusBadRequest, "unauthorized_client", "Client is not authorized for this grant."}
		case errors.Is(err, ErrInvalidParameter):
			return nil, &hostedUIError{http.StatusBadRequest, "invalid_scope", "Requested scope is invalid."}
		case errors.Is(err, ErrNotAuthorized):
			return nil, &hostedUIError{http.StatusBadRequest, "invalid_client", "Client authentication failed."}
		default:
			return nil, hostedUIServerError()
		}
	}

	result, _ := out.(map[string]interface{})
	auth, _ := result["ClientAuthenticationResult"].(map[string]interface{})
	accessToken, _ := auth["AccessToken"].(string)
	expiresIn, _ := auth["ExpiresIn"].(int64)
	if accessToken == "" {
		return nil, hostedUIServerError()
	}
	return &hostedUIClientCredentialsResult{accessToken: accessToken, expiresIn: expiresIn}, nil
}

// verifyPKCECodeChallenge checks a code_verifier against the S256 challenge
// captured at issue time: BASE64URL(SHA256(verifier)) must equal the
// challenge, compared in constant time.
func verifyPKCECodeChallenge(challenge, verifier string) bool {
	sum := sha256.Sum256([]byte(verifier))
	computed := base64.RawURLEncoding.EncodeToString(sum[:])
	return subtle.ConstantTimeCompare([]byte(computed), []byte(challenge)) == 1
}

// hostedUITokensFromAuthResult reads the AuthenticationResult map of a
// shared authenticateUser or refreshAuthToken response into the hosted-UI
// token set. A refresh grant echoes the caller's refresh token — its map
// carries no RefreshToken member — so the caller supplies the value to
// hold; an ordinary sign-in passes the empty string and takes the minted
// one. The boolean is false when the response carries no
// AuthenticationResult at all.
func hostedUITokensFromAuthResult(m map[string]interface{}, refreshToken string) (*hostedUITokenResult, bool) {
	authResult, ok := m["AuthenticationResult"].(map[string]interface{})
	if !ok {
		return nil, false
	}
	tokens := &hostedUITokenResult{RefreshToken: refreshToken}
	tokens.AccessToken, _ = authResult["AccessToken"].(string)
	tokens.IDToken, _ = authResult["IdToken"].(string)
	if d, ok := authResult["ExpiresIn"].(int64); ok {
		tokens.ExpiresIn = d
	}
	return tokens, true
}

// hostedUIAuthenticateCore runs the hosted UI's credential verification
// through authenticateUser, so user migration, imported-password hashes,
// FORCE_CHANGE_PASSWORD handling, Pre/PostAuthentication triggers, MFA
// enforcement and auth-event recording apply on the hosted UI exactly as on
// the API plane. The authorize flow's nonce is bound into the ID token.
func (s *CognitoService) hostedUIAuthenticateCore(ctx context.Context, reqCtx *request.RequestContext, poolID, clientID, username, password, nonce string) (*hostedUIAuthOutcome, *hostedUIError) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, hostedUIServerError()
	}
	pool, err := store.GetUserPool(poolID)
	if err != nil {
		return nil, hostedUIServerError()
	}

	res, err := s.authenticateUser(ctx, reqCtx, poolID, clientID, username, password, "", pool.LambdaConfig, nil, nil, nonce)
	if err != nil {
		return nil, hostedUIAuthError(err)
	}

	m, ok := res.(map[string]interface{})
	if !ok {
		return nil, hostedUIServerError()
	}
	if _, challenged := m["ChallengeName"]; challenged {
		name, _ := m["ChallengeName"].(string)
		return &hostedUIAuthOutcome{ChallengeName: name}, nil
	}
	tokens, ok := hostedUITokensFromAuthResult(m, "")
	if !ok {
		return nil, hostedUIServerError()
	}
	return &hostedUIAuthOutcome{Tokens: tokens}, nil
}

// hostedUIAuthError maps an authenticateUser failure onto the OAuth2 token
// endpoint's error vocabulary: credential failures are invalid_grant with
// AWS's descriptions (the token endpoint's error family is documented as
// HTTP 400), everything else is an internal error.
func hostedUIAuthError(err error) *hostedUIError {
	switch {
	case errors.Is(err, ErrNotAuthorized), errors.Is(err, ErrIncorrectPassword):
		return &hostedUIError{http.StatusBadRequest, "invalid_grant", "Incorrect username or password."}
	case errors.Is(err, ErrUserNotConfirmed):
		return &hostedUIError{http.StatusBadRequest, "invalid_grant", "User is not confirmed."}
	default:
		return hostedUIServerError()
	}
}

// issueHostedUIAuthCode mints a single-use authorization code for a browser
// sign-in that just completed, binding the code to the domain's pool, the
// client that requested the sign-in and the authenticated user, and to the
// authorisation-request members the exchange must re-verify: the redirect
// URI the code was issued for, the OIDC nonce, and the PKCE challenge.
func (s *CognitoService) issueHostedUIAuthCode(reqCtx *request.RequestContext, poolID, clientID, username, redirectURI, nonce, codeChallenge string) (string, *hostedUIError) {
	store, err := s.store(reqCtx)
	if err != nil {
		return "", hostedUIServerError()
	}
	user, err := store.GetUser(poolID, username)
	if err != nil {
		return "", hostedUIServerError()
	}
	code, err := generateAuthCode()
	if err != nil {
		return "", hostedUIServerError()
	}
	s.authCodes.Store(code, authCodeEntry{
		poolID:        poolID,
		userID:        user.ID,
		clientID:      clientID,
		redirectURI:   redirectURI,
		nonce:         nonce,
		codeChallenge: codeChallenge,
		expires:       time.Now().Add(hostedUIAuthCodeTTL),
	})
	s.startAuthCodeCleanup()
	return code, nil
}

// hostedUIAuthorizationCodeCore exchanges a hosted-UI authorization code for
// tokens. The exchange loads and consumes the code in one step (single use
// holds even when a later check fails), then enforces its bindings: the code
// must have been issued for the domain's pool, the requesting client and the
// redirect URI the authorisation request named — AWS documents a mismatching
// redirect_uri as unauthorized_client with error_description invalid_redirect
// — a confidential client must present its secret, and a code issued against
// a PKCE challenge is redeemed only with the matching verifier.
func (s *CognitoService) hostedUIAuthorizationCodeCore(reqCtx *request.RequestContext, poolID, clientID, clientSecret, code, redirectURI, codeVerifier string) (*hostedUITokenResult, *hostedUIError) {
	if code == "" {
		return nil, &hostedUIError{http.StatusBadRequest, "invalid_request", "Missing code parameter"}
	}
	// AWS token-endpoint contract: "You must provide this parameter if
	// grant_type is authorization_code."
	if redirectURI == "" {
		return nil, &hostedUIError{http.StatusBadRequest, "invalid_request", "Missing redirect_uri"}
	}
	raw, ok := s.authCodes.LoadAndDelete(code)
	if !ok {
		return nil, &hostedUIError{http.StatusBadRequest, "invalid_grant", "Invalid authorization code."}
	}
	entry := raw.(authCodeEntry)
	if time.Now().After(entry.expires) {
		return nil, &hostedUIError{http.StatusBadRequest, "invalid_grant", "Authorization code has expired."}
	}
	if entry.poolID != poolID {
		return nil, &hostedUIError{http.StatusBadRequest, "invalid_grant", "Authorization code was not issued for this user pool."}
	}
	if entry.clientID == "" || entry.clientID != clientID {
		return nil, &hostedUIError{http.StatusBadRequest, "invalid_grant", "Client ID mismatch."}
	}
	if redirectURI != entry.redirectURI {
		return nil, &hostedUIError{http.StatusBadRequest, "unauthorized_client", "invalid_redirect"}
	}

	client, herr := s.hostedUIClientAuthenticationCore(reqCtx, poolID, clientID, clientSecret)
	if herr != nil {
		return nil, herr
	}
	// The code grant must be among the client's allowed OAuth flows.
	if !clientAllowsOAuthFlow(client, "code") {
		return nil, &hostedUIError{http.StatusBadRequest, "unauthorized_client", "Client is not allowed for code grant flow."}
	}

	if entry.codeChallenge != "" {
		if codeVerifier == "" {
			return nil, &hostedUIError{http.StatusBadRequest, "invalid_request", "Missing code_verifier"}
		}
		if !verifyPKCECodeChallenge(entry.codeChallenge, codeVerifier) {
			return nil, &hostedUIError{http.StatusBadRequest, "invalid_grant", "PKCE verification failed."}
		}
	}

	// The OIDC nonce rides into the ID token, per the authorize-endpoint
	// contract: "The nonce value that you provide is included in the ID
	// token that Amazon Cognito issues."
	var idClaims map[string]interface{}
	if entry.nonce != "" {
		idClaims = map[string]interface{}{"nonce": entry.nonce}
	}
	accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokensWithIDClaims(reqCtx, entry.poolID, entry.userID, entry.clientID, TokenGenerationHostedAuth, nil, idClaims)
	if err != nil {
		return nil, &hostedUIError{http.StatusInternalServerError, "server_error", "Failed to create tokens."}
	}
	return &hostedUITokenResult{AccessToken: accessToken, IDToken: idToken, RefreshToken: refreshToken, ExpiresIn: expiresIn}, nil
}

// hostedUIRefreshTokenCore refreshes tokens through the shared
// refreshAuthToken Core with the domain's pool and the requesting client
// bound to the stored token: a refresh token issued by another pool or
// client cannot mint tokens here, and a client holding a secret must
// authenticate like on every other grant. Token-domain failures are
// invalid_grant (HTTP 400 per the token-endpoint error map); infrastructure
// failures are server errors. The original refresh token is echoed in the
// result — a refresh grant never issues a new one.
func (s *CognitoService) hostedUIRefreshTokenCore(reqCtx *request.RequestContext, poolID, clientID, clientSecret, refreshToken string) (*hostedUITokenResult, *hostedUIError) {
	if refreshToken == "" {
		return nil, &hostedUIError{http.StatusBadRequest, "invalid_request", "Missing refresh_token"}
	}
	client, herr := s.hostedUIClientAuthenticationCore(reqCtx, poolID, clientID, clientSecret)
	if herr != nil {
		return nil, herr
	}
	// Refresh tokens exist only for code-grant sessions, so the refresh
	// grant is gated on the same flow ("Client is not allowed for code
	// grant flow or for refreshing tokens.").
	if !clientAllowsOAuthFlow(client, "code") {
		return nil, &hostedUIError{http.StatusBadRequest, "unauthorized_client", "Client is not allowed for code grant flow or for refreshing tokens."}
	}

	res, err := s.refreshAuthToken(reqCtx, poolID, clientID, refreshToken, "", false)
	if err != nil {
		if errors.Is(err, ErrNotAuthorized) {
			return nil, &hostedUIError{http.StatusBadRequest, "invalid_grant", "Invalid refresh token."}
		}
		return nil, hostedUIServerError()
	}
	m, ok := res.(map[string]interface{})
	if !ok {
		return nil, hostedUIServerError()
	}
	tokens, ok := hostedUITokensFromAuthResult(m, refreshToken)
	if !ok {
		return nil, hostedUIServerError()
	}
	return tokens, nil
}

// hostedUILogoutCore validates a logout redirect and returns the redirect
// target. The endpoint takes either a logout_uri (an Allowed sign-out URL
// for the client) or a redirect_uri (an Allowed callback URL, together with
// the parameters of a well-formed authorize request — the user is
// redirected to the sign-in page with those parameters appended); when both
// are present the logout_uri wins ("Amazon Cognito will utilize the
// logout_uri parameter exclusively, overriding the redirect_uri
// parameter"). Registering the URIs keeps /logout from becoming an open
// redirect.
func (s *CognitoService) hostedUILogoutCore(reqCtx *request.RequestContext, poolID, clientID, logoutURI, redirectURI, responseType string) (string, *hostedUIError) {
	if clientID == "" {
		return "", &hostedUIError{http.StatusBadRequest, "invalid_request", "Missing client_id"}
	}
	if logoutURI == "" && redirectURI == "" {
		return "", &hostedUIError{http.StatusBadRequest, "invalid_request", "Missing logout_uri"}
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return "", hostedUIServerError()
	}
	client, err := store.GetUserPoolClient(poolID, clientID)
	if err != nil {
		if !errors.Is(err, cognitostore.ErrClientNotFound) {
			return "", hostedUIServerError()
		}
		return "", &hostedUIError{http.StatusBadRequest, "invalid_request", "Invalid client_id"}
	}
	if logoutURI != "" {
		for _, uri := range client.LogoutURLs {
			if uri == logoutURI {
				return logoutURI, nil
			}
		}
		return "", &hostedUIError{http.StatusBadRequest, "invalid_request", "Invalid logout_uri"}
	}
	// The redirect_uri form re-invites the user to sign in: response_type is
	// required alongside it, and the URI must be one of the client's
	// registered callback URLs.
	if responseType == "" {
		return "", &hostedUIError{http.StatusBadRequest, "invalid_request", "Missing response_type"}
	}
	if responseType != "code" && responseType != "token" {
		return "", &hostedUIError{http.StatusBadRequest, "invalid_request", "Invalid response_type"}
	}
	for _, uri := range client.CallbackURLs {
		if uri == redirectURI {
			return "/login", nil
		}
	}
	return "", &hostedUIError{http.StatusBadRequest, "invalid_request", "Invalid redirect_uri"}
}

// hostedUISignUpErrorDescription extracts the human-readable message of a
// signUpCore failure for the sign-up page.
func hostedUISignUpErrorDescription(err error) string {
	var awsErr *awserrors.AWSError
	if errors.As(err, &awsErr) && awsErr.Message != "" {
		return awsErr.Message
	}
	return "Sign-up failed."
}

// hostedUIConfirmErrorDescription extracts the human-readable message of a
// confirmSignUpCore failure for the confirmation page.
func hostedUIConfirmErrorDescription(err error) string {
	var awsErr *awserrors.AWSError
	if errors.As(err, &awsErr) && awsErr.Message != "" {
		return awsErr.Message
	}
	return "Confirmation failed."
}
