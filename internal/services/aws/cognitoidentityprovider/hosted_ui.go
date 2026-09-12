package cognitoidentityprovider

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strings"

	waf "vorpalstacks/internal/common/invokers/waf"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/server/fqdnrouter"
)

// HostedUIHandler serves the Cognito hosted UI pages for login, sign-up, and OAuth2 flows.
func (s *CognitoService) HostedUIHandler(w http.ResponseWriter, r *http.Request) {
	if waf.ServeWAFTokenExchange(r.Context(), s.waf.currentInspector(), w, r) {
		return
	}

	if s.storageManager == nil {
		http.Error(w, "Cognito hosted UI not available", http.StatusServiceUnavailable)
		return
	}

	domain := fqdnrouter.ResourceIDFromContext(r.Context())
	if domain == "" {
		domain = s.extractDomain(r.Host)
	}
	if domain == "" {
		http.Error(w, "Could not determine domain from Host", http.StatusBadRequest)
		return
	}

	poolID, err := s.resolveDomainToPoolID(domain)
	if err != nil || poolID == "" {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	// WAF inspection covers every hosted UI endpoint; the hosted UI
	// forwards no request body to AWS WAF, so only headers and the path
	// are inspected.
	if poolARN := s.poolARNByID(poolID); poolARN != "" {
		if s.enforceWAFOnHostedUI(w, r, poolARN) {
			return
		}
	}

	path := r.URL.Path

	switch {
	case path == "/login" || path == "/signin":
		if r.Method == http.MethodPost {
			s.handleLoginSubmit(w, r, poolID)
		} else {
			s.renderLoginPage(w, r, hostedUIFlowParamsFromQuery(r.URL.Query()), poolID, "")
		}
	case path == "/signup" || path == "/register":
		if r.Method == http.MethodPost {
			s.handleSignUpSubmit(w, r, poolID)
		} else {
			s.renderSignUpPage(w, r, hostedUIFlowParamsFromQuery(r.URL.Query()), poolID, "")
		}
	case path == "/confirm":
		if r.Method == http.MethodPost {
			s.handleConfirmSubmit(w, r, poolID)
		} else {
			s.renderConfirmPage(w, r, hostedUIFlowParamsFromQuery(r.URL.Query()), "")
		}
	case path == "/oauth2/authorize":
		s.handleAuthorize(w, r, poolID)
	case path == "/oauth2/token":
		s.handleTokenEndpoint(w, r, poolID)
	case path == "/logout":
		q := r.URL.Query()
		reqCtx := request.NewRequestContext(r.Context(), s.storageManager, s.accountID, s.region)
		target, herr := s.hostedUILogoutCore(reqCtx, poolID, q.Get("client_id"), q.Get("logout_uri"), q.Get("redirect_uri"), q.Get("response_type"))
		if herr != nil {
			http.Error(w, herr.desc, herr.status)
			return
		}
		if target == "/login" {
			// The redirect_uri form re-invites the user to the sign-in
			// page with the original authorize parameters appended ("The
			// logout endpoint appends the parameters in your original
			// request to the redirect destination").
			q.Del("logout_uri")
			http.Redirect(w, r, "/login?"+q.Encode(), http.StatusFound)
			return
		}
		http.Redirect(w, r, target, http.StatusFound)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (s *CognitoService) extractDomain(host string) string {
	// Host labels fold to lowercase, matching the production router's
	// lowercasing of the Host header before domain extraction.
	host = strings.ToLower(strings.Split(host, ":")[0])
	parts := strings.Split(host, ".")
	if len(parts) >= 2 {
		return parts[0]
	}
	return ""
}

func writeOAuthError(w http.ResponseWriter, status int, errCode, desc string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error":             errCode,
		"error_description": desc,
	})
}

func writeTokenJSON(w http.ResponseWriter, accessToken, idToken, refreshToken string, expiresIn int64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":  accessToken,
		"token_type":    "Bearer",
		"expires_in":    expiresIn,
		"id_token":      idToken,
		"refresh_token": refreshToken,
	})
}

// writeClientCredentialsTokenJSON renders the client_credentials grant's
// response: the machine access token alone, with no ID or refresh token.
func writeClientCredentialsTokenJSON(w http.ResponseWriter, accessToken string, expiresIn int64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   expiresIn,
	})
}

// hostedUIErrorLine renders an error message for the hosted UI pages; an
// empty message renders nothing.
func hostedUIErrorLine(message string) string {
	if message == "" {
		return ""
	}
	return `<p class="error">` + html.EscapeString(message) + `</p>`
}

// handleAuthorize fronts the /oauth2/authorize endpoint: it validates the
// authorisation request in the Core before any page renders, redirects the
// documented formatting errors to the registered redirect URI, and
// otherwise serves the login page carrying the request's members.
func (s *CognitoService) handleAuthorize(w http.ResponseWriter, r *http.Request, poolID string) {
	f := hostedUIFlowParamsFromQuery(r.URL.Query())
	reqCtx := request.NewRequestContext(r.Context(), s.storageManager, s.accountID, s.region)

	_, oauthErr, herr := s.hostedUIAuthorizeCore(reqCtx, poolID, f)
	if herr != nil {
		// The client or redirect URI failed validation: no redirect target
		// is proven trustworthy, so the failure renders as an error page.
		http.Error(w, herr.desc, herr.status)
		return
	}
	if oauthErr != "" {
		// AWS documents the formatting-error class as a bare error redirect:
		// "HTTP 1.1 302 Found Location: https://client_redirect_uri?error=invalid_request".
		eq := url.Values{}
		eq.Set("error", oauthErr)
		sep := "?"
		if strings.Contains(f.RedirectURI, "?") {
			sep = "&"
		}
		http.Redirect(w, r, f.RedirectURI+sep+eq.Encode(), http.StatusFound)
		return
	}
	s.renderLoginPage(w, r, f, poolID, "")
}

// parseBasicClientCredentials decodes client_secret_basic authentication:
// Authorization: Basic base64(client_id:client_secret). A malformed header
// reports ok=false and is ignored — the grant then fails on the missing
// credentials like any unauthenticated request.
func parseBasicClientCredentials(header string) (clientID, clientSecret string, ok bool) {
	const prefix = "Basic "
	if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return "", "", false
	}
	raw, err := base64.StdEncoding.DecodeString(header[len(prefix):])
	if err != nil {
		return "", "", false
	}
	parts := strings.SplitN(string(raw), ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func (s *CognitoService) handleLoginSubmit(w http.ResponseWriter, r *http.Request, poolID string) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not parse form data", http.StatusBadRequest)
		return
	}
	f := hostedUIFlowParamsFromForm(r)
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		s.renderLoginPage(w, r, f, poolID, "")
		return
	}

	reqCtx := request.NewRequestContext(r.Context(), s.storageManager, s.accountID, s.region)

	if _, herr := s.hostedUIClientCore(reqCtx, poolID, f.ClientID, f.RedirectURI); herr != nil {
		s.renderLoginPage(w, r, f, poolID, herr.desc)
		return
	}

	outcome, herr := s.hostedUIAuthenticateCore(r.Context(), reqCtx, poolID, f.ClientID, username, password, f.Nonce)
	if herr != nil {
		s.renderLoginPage(w, r, f, poolID, herr.desc)
		return
	}
	// The hosted UI pages collect a username and password only; when the
	// pool's policy replaces token issuance with a challenge the pages
	// cannot answer, the sign-in is refused rather than weakened.
	if outcome.ChallengeName != "" {
		s.renderLoginPage(w, r, f, poolID, "Sign-in requires an additional challenge that this page does not support.")
		return
	}

	if f.RedirectURI == "" {
		f.RedirectURI = "/"
	}

	if f.ResponseType == "token" {
		frag := url.Values{}
		frag.Set("access_token", outcome.Tokens.AccessToken)
		frag.Set("id_token", outcome.Tokens.IDToken)
		frag.Set("token_type", "Bearer")
		frag.Set("expires_in", fmt.Sprintf("%d", outcome.Tokens.ExpiresIn))
		if f.State != "" {
			frag.Set("state", f.State)
		}
		http.Redirect(w, r, f.RedirectURI+"#"+frag.Encode(), http.StatusFound)
		return
	}

	code, herr := s.issueHostedUIAuthCode(reqCtx, poolID, f.ClientID, username, f.RedirectURI, f.Nonce, f.CodeChallenge)
	if herr != nil {
		http.Error(w, herr.desc, herr.status)
		return
	}

	q := url.Values{}
	q.Set("code", code)
	if f.State != "" {
		q.Set("state", f.State)
	}
	sep := "?"
	if strings.Contains(f.RedirectURI, "?") {
		sep = "&"
	}
	http.Redirect(w, r, f.RedirectURI+sep+q.Encode(), http.StatusFound)
}

// handleSignUpSubmit creates the account through signUpCore — the same
// registration path as the SignUp API. The default flow leaves the account
// awaiting its confirmation code, so the browser is sent to the
// confirmation page carrying the authorisation request's members; an
// auto-confirmed registration continues straight to the login page.
func (s *CognitoService) handleSignUpSubmit(w http.ResponseWriter, r *http.Request, poolID string) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not parse form data", http.StatusBadRequest)
		return
	}
	f := hostedUIFlowParamsFromForm(r)
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	reqCtx := request.NewRequestContext(r.Context(), s.storageManager, s.accountID, s.region)

	// The sign-up belongs to the pool behind this domain; a client from
	// another pool must not register users through it.
	if _, herr := s.hostedUIClientCore(reqCtx, poolID, f.ClientID, ""); herr != nil {
		s.renderSignUpPage(w, r, f, poolID, herr.desc)
		return
	}

	attrs := map[string]string{}
	if email != "" {
		attrs["email"] = email
	}
	result, err := s.signUpCore(r.Context(), reqCtx, SignUpInput{
		ClientID:       f.ClientID,
		Username:       username,
		Password:       password,
		UserAttributes: attrs,
	})
	if err != nil {
		s.renderSignUpPage(w, r, f, poolID, hostedUISignUpErrorDescription(err))
		return
	}

	confirmed := false
	if m, ok := result.(map[string]interface{}); ok {
		confirmed, _ = m["UserConfirmed"].(bool)
	}
	target := "/login?" + f.pageQuery().Encode()
	if !confirmed {
		target = "/confirm?" + f.pageQuery().Encode()
	}
	http.Redirect(w, r, target, http.StatusFound)
}

// handleConfirmSubmit completes the hosted-UI sign-up flow: the
// confirmation code the user received is checked through confirmSignUpCore
// — the same path as the ConfirmSignUp API — and a confirmed account
// continues to the login page with the authorisation request's members.
func (s *CognitoService) handleConfirmSubmit(w http.ResponseWriter, r *http.Request, poolID string) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not parse form data", http.StatusBadRequest)
		return
	}
	f := hostedUIFlowParamsFromForm(r)
	username := r.FormValue("username")
	confirmationCode := r.FormValue("confirmation_code")

	reqCtx := request.NewRequestContext(r.Context(), s.storageManager, s.accountID, s.region)

	if _, err := s.confirmSignUpCore(r.Context(), reqCtx, ConfirmSignUpInput{
		ClientID:         f.ClientID,
		Username:         username,
		ConfirmationCode: confirmationCode,
	}); err != nil {
		s.renderConfirmPage(w, r, f, hostedUIConfirmErrorDescription(err))
		return
	}
	http.Redirect(w, r, "/login?"+f.pageQuery().Encode(), http.StatusFound)
}

// hostedUIHiddenFields renders the authorisation-request members as hidden
// form inputs so a page submission carries the flow forward. client_id must
// name the app client, not the pool, or the submitted sign-in can never
// validate.
func hostedUIHiddenFields(f hostedUIFlowParams) string {
	return fmt.Sprintf(`<input type="hidden" name="client_id" value="%s">
<input type="hidden" name="redirect_uri" value="%s">
<input type="hidden" name="response_type" value="%s">
<input type="hidden" name="state" value="%s">
<input type="hidden" name="nonce" value="%s">
<input type="hidden" name="code_challenge" value="%s">
<input type="hidden" name="code_challenge_method" value="%s">`,
		html.EscapeString(f.ClientID), html.EscapeString(f.RedirectURI), html.EscapeString(f.ResponseType),
		html.EscapeString(f.State), html.EscapeString(f.Nonce),
		html.EscapeString(f.CodeChallenge), html.EscapeString(f.CodeChallengeMethod))
}

func (s *CognitoService) renderLoginPage(w http.ResponseWriter, r *http.Request, f hostedUIFlowParams, poolID, errMsg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Sign In</title>
<style>
body { font-family: sans-serif; max-width: 400px; margin: 60px auto; padding: 0 20px; }
h1 { color: #333; }
input { display: block; width: 100%%; padding: 8px; margin: 8px 0; box-sizing: border-box; }
button { width: 100%%; padding: 10px; background: #1597D3; color: white; border: none; cursor: pointer; margin-top: 10px; }
button:hover { background: #1274A3; }
.error { color: red; font-size: 0.9em; }
</style>
</head>
<body>
<h1>Sign in</h1>
%s
<form method="POST" action="/login">
%s
<label>Username</label>
<input type="text" name="username" required>
<label>Password</label>
<input type="password" name="password" required>
<button type="submit">Sign In</button>
</form>
</body>
</html>`, hostedUIErrorLine(errMsg), hostedUIHiddenFields(f))
}

func (s *CognitoService) renderSignUpPage(w http.ResponseWriter, r *http.Request, f hostedUIFlowParams, poolID, errMsg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Sign Up</title>
<style>
body { font-family: sans-serif; max-width: 400px; margin: 60px auto; padding: 0 20px; }
h1 { color: #333; }
input { display: block; width: 100%%; padding: 8px; margin: 8px 0; box-sizing: border-box; }
button { width: 100%%; padding: 10px; background: #1597D3; color: white; border: none; cursor: pointer; margin-top: 10px; }
button:hover { background: #1274A3; }
.error { color: red; font-size: 0.9em; }
</style>
</head>
<body>
<h1>Sign Up</h1>
%s
<form method="POST" action="/signup">
%s
<label>Username</label>
<input type="text" name="username" required>
<label>Password</label>
<input type="password" name="password" required>
<label>Email</label>
<input type="email" name="email">
<button type="submit">Sign Up</button>
</form>
</body>
</html>`, hostedUIErrorLine(errMsg), hostedUIHiddenFields(f))
}

// renderConfirmPage serves the confirmation-code step of the hosted-UI
// sign-up flow: the account exists but cannot sign in until the delivered
// code is entered, per the documented confirmation process ("The user
// enters the confirmation code in the app ... sets the user's account to
// the confirmed state ... and the user can sign in").
func (s *CognitoService) renderConfirmPage(w http.ResponseWriter, r *http.Request, f hostedUIFlowParams, errMsg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Confirm your account</title>
<style>
body { font-family: sans-serif; max-width: 400px; margin: 60px auto; padding: 0 20px; }
h1 { color: #333; }
input { display: block; width: 100%%; padding: 8px; margin: 8px 0; box-sizing: border-box; }
button { width: 100%%; padding: 10px; background: #1597D3; color: white; border: none; cursor: pointer; margin-top: 10px; }
button:hover { background: #1274A3; }
.error { color: red; font-size: 0.9em; }
</style>
</head>
<body>
<h1>Confirm your account</h1>
<p>Enter the confirmation code that was sent to your email address or phone number.</p>
%s
<form method="POST" action="/confirm">
%s
<label>Username</label>
<input type="text" name="username" required>
<label>Confirmation code</label>
<input type="text" name="confirmation_code" required>
<button type="submit">Confirm</button>
</form>
</body>
</html>`, hostedUIErrorLine(errMsg), hostedUIHiddenFields(f))
}

func (s *CognitoService) handleTokenEndpoint(w http.ResponseWriter, r *http.Request, poolID string) {
	if r.Method != http.MethodPost {
		writeOAuthError(w, http.StatusMethodNotAllowed, "invalid_request", "Method not allowed")
		return
	}

	if err := r.ParseForm(); err != nil {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "Could not parse form data")
		return
	}

	grantType := r.FormValue("grant_type")
	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")
	// client_secret_basic: the Authorization header carries the same pair
	// the body may carry; an explicit body parameter wins so the
	// client_secret_post form stays authoritative.
	if basicID, basicSecret, ok := parseBasicClientCredentials(r.Header.Get("Authorization")); ok {
		if clientID == "" {
			clientID = basicID
		}
		if clientSecret == "" {
			clientSecret = basicSecret
		}
	}
	reqCtx := request.NewRequestContext(r.Context(), s.storageManager, s.accountID, s.region)

	switch grantType {
	case "authorization_code":
		res, herr := s.hostedUIAuthorizationCodeCore(reqCtx, poolID, clientID, clientSecret,
			r.FormValue("code"), r.FormValue("redirect_uri"), r.FormValue("code_verifier"))
		if herr != nil {
			writeOAuthError(w, herr.status, herr.code, herr.desc)
			return
		}
		writeTokenJSON(w, res.AccessToken, res.IDToken, res.RefreshToken, res.ExpiresIn)

	case "client_credentials":
		// The client itself is the principal: the confidential client
		// authenticates (client_secret_basic or client_secret_post) and
		// the machine token issues through the same path the API-plane
		// client-token operation uses.
		res, herr := s.hostedUIClientCredentialsTokenCore(r.Context(), reqCtx, poolID, clientID, clientSecret, r.FormValue("scope"))
		if herr != nil {
			writeOAuthError(w, herr.status, herr.code, herr.desc)
			return
		}
		writeClientCredentialsTokenJSON(w, res.accessToken, res.expiresIn)

	case "refresh_token":
		res, herr := s.hostedUIRefreshTokenCore(reqCtx, poolID, clientID, clientSecret, r.FormValue("refresh_token"))
		if herr != nil {
			writeOAuthError(w, herr.status, herr.code, herr.desc)
			return
		}
		writeTokenJSON(w, res.AccessToken, res.IDToken, res.RefreshToken, res.ExpiresIn)

	default:
		writeOAuthError(w, http.StatusBadRequest, "unsupported_grant_type",
			fmt.Sprintf("Grant type '%s' is not supported", grantType))
	}
}
