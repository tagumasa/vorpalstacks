package cognitoidentityprovider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"vorpalstacks/internal/common/request"
)

// HostedUIHandler serves the Cognito hosted UI pages for login, sign-up, and OAuth2 flows.
func (s *CognitoService) HostedUIHandler(w http.ResponseWriter, r *http.Request) {
	if s.storageManager == nil {
		http.Error(w, "Cognito hosted UI not available", http.StatusServiceUnavailable)
		return
	}

	domain := s.extractDomain(r.Host)
	if domain == "" {
		http.Error(w, "Could not determine domain from Host", http.StatusBadRequest)
		return
	}

	poolID, err := s.resolveDomainToPoolID(domain)
	if err != nil || poolID == "" {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	path := r.URL.Path

	switch {
	case path == "/login" || path == "/signin":
		s.renderLoginPage(w, r, poolID)
	case path == "/signup" || path == "/register":
		s.renderSignUpPage(w, r, poolID)
	case path == "/oauth2/authorize":
		s.renderLoginPage(w, r, poolID)
	case path == "/oauth2/token":
		s.handleTokenEndpoint(w, r, poolID)
	case path == "/logout":
		redirectURL := r.URL.Query().Get("logout_uri")
		if redirectURL == "" {
			redirectURL = "/"
		}
		http.Redirect(w, r, redirectURL, http.StatusFound)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (s *CognitoService) extractDomain(host string) string {
	host = strings.Split(host, ":")[0]
	parts := strings.Split(host, ".")
	if len(parts) >= 2 {
		return parts[0]
	}
	return ""
}

func (s *CognitoService) resolveDomainToPoolID(domain string) (string, error) {
	ctx := request.NewRequestContext(context.Background(), s.storageManager, s.accountID, s.region)
	pools, _ := s.ListUserPoolsRaw(ctx)
	for _, pool := range pools {
		store, err := s.store(ctx)
		if err != nil {
			continue
		}
		domainEntry, err := store.GetUserPoolDomain(domain)
		if err == nil && domainEntry != nil && domainEntry.UserPoolID == pool.ID {
			return pool.ID, nil
		}
	}
	return "", fmt.Errorf("domain %s not found", domain)
}

func (s *CognitoService) renderLoginPage(w http.ResponseWriter, r *http.Request, poolID string) {
	_ = r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	responseType := r.URL.Query().Get("response_type")

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
<form method="POST" action="/login">
<input type="hidden" name="client_id" value="%s">
<input type="hidden" name="redirect_uri" value="%s">
<input type="hidden" name="response_type" value="%s">
<label>Username</label>
<input type="text" name="username" required>
<label>Password</label>
<input type="password" name="password" required>
<button type="submit">Sign In</button>
</form>
</body>
</html>`, poolID, redirectURI, responseType)
}

func (s *CognitoService) renderSignUpPage(w http.ResponseWriter, r *http.Request, poolID string) {
	_ = r.URL.Query().Get("client_id")

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
</style>
</head>
<body>
<h1>Sign Up</h1>
<form method="POST" action="/signup">
<input type="hidden" name="client_id" value="%s">
<label>Username</label>
<input type="text" name="username" required>
<label>Password</label>
<input type="password" name="password" required>
<label>Email</label>
<input type="email" name="email">
<button type="submit">Sign Up</button>
</form>
</body>
</html>`, "")
}

func (s *CognitoService) handleTokenEndpoint(w http.ResponseWriter, r *http.Request, poolID string) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"invalid_request"}`, http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	grantType := r.FormValue("grant_type")

	if grantType != "authorization_code" && grantType != "client_credentials" && grantType != "refresh_token" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"unsupported_grant_type","error_description":"Grant type '%s' is not supported in mock mode"}`, grantType)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
"access_token":"mock-access-token-%s",
"token_type":"Bearer",
"expires_in":3600,
"id_token":"mock-id-token-%s",
"refresh_token":"mock-refresh-token-%s"
}`, poolID, poolID, poolID)
}
