package cognitoidentityprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"vorpalstacks/internal/common/request"

	"golang.org/x/crypto/bcrypt"
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error":             "invalid_request",
			"error_description": "Method not allowed",
		})
		return
	}

	if err := r.ParseForm(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":             "invalid_request",
			"error_description": "Could not parse form data",
		})
		return
	}

	grantType := r.FormValue("grant_type")
	clientID := r.FormValue("client_id")

	switch grantType {
	case "authorization_code", "password":
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == "" || password == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_request",
				"error_description": "Missing username or password",
			})
			return
		}

		ctx := request.NewRequestContext(context.Background(), s.storageManager, s.accountID, s.region)
		store, err := s.store(ctx)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "server_error",
				"error_description": "Internal error",
			})
			return
		}

		user, err := store.GetUser(poolID, username)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_grant",
				"error_description": "Incorrect username or password.",
			})
			return
		}

		if !user.Enabled {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_grant",
				"error_description": "Incorrect username or password.",
			})
			return
		}

		if user.UserStatus != "CONFIRMED" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_grant",
				"error_description": "User is not confirmed.",
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_grant",
				"error_description": "Incorrect username or password.",
			})
			return
		}

		accessToken, idToken, refreshToken, expiresIn := s.CreateTokens(ctx, poolID, user.ID, clientID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token":  accessToken,
			"token_type":    "Bearer",
			"expires_in":    expiresIn,
			"id_token":      idToken,
			"refresh_token": refreshToken,
		})

	case "refresh_token":
		refreshToken := r.FormValue("refresh_token")
		if refreshToken == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_request",
				"error_description": "Missing refresh_token",
			})
			return
		}

		ctx := request.NewRequestContext(context.Background(), s.storageManager, s.accountID, s.region)
		store, err := s.store(ctx)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "server_error",
				"error_description": "Internal error",
			})
			return
		}

		storedToken, err := store.GetRefreshTokenByValue(refreshToken)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_grant",
				"error_description": "Invalid refresh token.",
			})
			return
		}

		user, err := store.GetUserByID(storedToken.UserID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "invalid_grant",
				"error_description": "Invalid refresh token.",
			})
			return
		}

		accessToken, idToken, _, expiresIn := s.CreateTokens(ctx, poolID, user.ID, clientID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token":  accessToken,
			"token_type":    "Bearer",
			"expires_in":    expiresIn,
			"id_token":      idToken,
			"refresh_token": refreshToken,
		})

	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":             "unsupported_grant_type",
			"error_description": fmt.Sprintf("Grant type '%s' is not supported", grantType),
		})
	}
}
