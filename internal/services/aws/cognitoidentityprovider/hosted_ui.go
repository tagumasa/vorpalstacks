package cognitoidentityprovider

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/server/fqdnrouter"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"golang.org/x/crypto/bcrypt"
)

// HostedUIHandler serves the Cognito hosted UI pages for login, sign-up, and OAuth2 flows.
func (s *CognitoService) HostedUIHandler(w http.ResponseWriter, r *http.Request) {
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

	path := r.URL.Path

	switch {
	case path == "/login" || path == "/signin":
		if r.Method == http.MethodPost {
			s.handleLoginSubmit(w, r, poolID)
		} else {
			s.renderLoginPage(w, r, poolID)
		}
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

type authCodeEntry struct {
	poolID   string
	userID   string
	clientID string
	expires  time.Time
}

func generateAuthCode() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *CognitoService) startAuthCodeCleanup() {
	s.authCodeCleanupOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(1 * time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				now := time.Now()
				s.authCodes.Range(func(key, value interface{}) bool {
					if entry, ok := value.(authCodeEntry); ok && now.After(entry.expires) {
						s.authCodes.Delete(key)
					}
					return true
				})
			}
		}()
	})
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

func (s *CognitoService) handleLoginSubmit(w http.ResponseWriter, r *http.Request, poolID string) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not parse form data", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	clientID := r.FormValue("client_id")
	redirectURI := r.FormValue("redirect_uri")
	responseType := r.FormValue("response_type")

	if username == "" || password == "" {
		s.renderLoginPage(w, r, poolID)
		return
	}

	ctx := request.NewRequestContext(context.Background(), s.storageManager, s.accountID, s.region)
	store, err := s.store(ctx)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if clientID != "" {
		client, err := store.GetUserPoolClient(poolID, clientID)
		if err != nil || client == nil {
			http.Error(w, "Invalid client_id", http.StatusBadRequest)
			return
		}
		if redirectURI != "" && !isRegisteredRedirectURI(client, redirectURI) {
			http.Error(w, "Invalid redirect_uri", http.StatusBadRequest)
			return
		}
	}

	user, err := store.GetUser(poolID, username)
	if err != nil || !user.Enabled || user.UserStatus != "CONFIRMED" {
		s.renderLoginPage(w, r, poolID)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.renderLoginPage(w, r, poolID)
		return
	}

	code, err := generateAuthCode()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	s.authCodes.Store(code, authCodeEntry{
		poolID:   poolID,
		userID:   user.ID,
		clientID: clientID,
		expires:  time.Now().Add(5 * time.Minute),
	})
	s.startAuthCodeCleanup()

	if redirectURI == "" {
		redirectURI = "/"
	}

	if responseType == "token" {
		accessToken, idToken, _, expiresIn, err := s.CreateTokens(ctx, poolID, user.ID, clientID, TokenGenerationHostedAuth, nil)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		frag := fmt.Sprintf("access_token=%s&id_token=%s&expires_in=%d&token_type=Bearer", accessToken, idToken, expiresIn)
		http.Redirect(w, r, redirectURI+"#"+frag, http.StatusFound)
		return
	}

	q := url.Values{}
	q.Set("code", code)
	sep := "?"
	if strings.Contains(redirectURI, "?") {
		sep = "&"
	}
	http.Redirect(w, r, redirectURI+sep+q.Encode(), http.StatusFound)
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
</html>`, html.EscapeString(poolID), html.EscapeString(redirectURI), html.EscapeString(responseType))
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
</html>`, html.EscapeString(poolID))
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

	switch grantType {
	case "authorization_code":
		code := r.FormValue("code")
		if code == "" {
			writeOAuthError(w, http.StatusBadRequest, "invalid_request", "Missing code parameter")
			return
		}
		raw, ok := s.authCodes.LoadAndDelete(code)
		if !ok {
			writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "Invalid authorization code.")
			return
		}
		entry := raw.(authCodeEntry)
		if time.Now().After(entry.expires) {
			writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "Authorization code has expired.")
			return
		}
		if clientID != entry.clientID {
			writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "Client ID mismatch.")
			return
		}
		ctx := request.NewRequestContext(context.Background(), s.storageManager, s.accountID, s.region)
		if entry.clientID != "" {
			store, err := s.store(ctx)
			if err == nil {
				client, err := store.GetUserPoolClient(entry.poolID, entry.clientID)
				if err == nil && client != nil && client.ClientSecret != "" {
					clientSecret := r.FormValue("client_secret")
					if subtle.ConstantTimeCompare([]byte(clientSecret), []byte(client.ClientSecret)) != 1 {
						writeOAuthError(w, http.StatusUnauthorized, "invalid_client", "Client authentication failed.")
						return
					}
				}
			}
		}
		accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokens(ctx, entry.poolID, entry.userID, entry.clientID, TokenGenerationHostedAuth, nil)
		if err != nil {
			writeOAuthError(w, http.StatusInternalServerError, "server_error", "Failed to create tokens.")
			return
		}
		writeTokenJSON(w, accessToken, idToken, refreshToken, expiresIn)

	case "password":
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == "" || password == "" {
			writeOAuthError(w, http.StatusBadRequest, "invalid_request", "Missing username or password")
			return
		}

		ctx := request.NewRequestContext(context.Background(), s.storageManager, s.accountID, s.region)
		store, err := s.store(ctx)
		if err != nil {
			writeOAuthError(w, http.StatusInternalServerError, "server_error", "Internal error")
			return
		}

		user, err := store.GetUser(poolID, username)
		if err != nil {
			writeOAuthError(w, http.StatusUnauthorized, "invalid_grant", "Incorrect username or password.")
			return
		}

		if !user.Enabled {
			writeOAuthError(w, http.StatusUnauthorized, "invalid_grant", "Incorrect username or password.")
			return
		}

		if user.UserStatus != "CONFIRMED" {
			writeOAuthError(w, http.StatusUnauthorized, "invalid_grant", "User is not confirmed.")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			writeOAuthError(w, http.StatusUnauthorized, "invalid_grant", "Incorrect username or password.")
			return
		}

		accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokens(ctx, poolID, user.ID, clientID, TokenGenerationHostedAuth, nil)
		if err != nil {
			writeOAuthError(w, http.StatusInternalServerError, "server_error", "Failed to create tokens.")
			return
		}

		writeTokenJSON(w, accessToken, idToken, refreshToken, expiresIn)

	case "refresh_token":
		refreshToken := r.FormValue("refresh_token")
		if refreshToken == "" {
			writeOAuthError(w, http.StatusBadRequest, "invalid_request", "Missing refresh_token")
			return
		}

		ctx := request.NewRequestContext(context.Background(), s.storageManager, s.accountID, s.region)
		store, err := s.store(ctx)
		if err != nil {
			writeOAuthError(w, http.StatusInternalServerError, "server_error", "Internal error")
			return
		}

		storedToken, err := store.GetRefreshTokenByValue(refreshToken)
		if err != nil {
			errCode := "invalid_grant"
			errDesc := "Invalid refresh token."
			if errors.Is(err, cognitostore.ErrTokenExpired) {
				errCode = "expired_grant"
				errDesc = "Refresh token has expired."
			}
			writeOAuthError(w, http.StatusUnauthorized, errCode, errDesc)
			return
		}

		user, err := store.GetUserByID(storedToken.UserID)
		if err != nil {
			writeOAuthError(w, http.StatusUnauthorized, "invalid_grant", "Invalid refresh token.")
			return
		}

		accessToken, idToken, _, expiresIn, err := s.CreateTokens(ctx, poolID, user.ID, clientID, TokenGenerationRefreshTokens, nil)
		if err != nil {
			writeOAuthError(w, http.StatusInternalServerError, "server_error", "Failed to create tokens.")
			return
		}

		writeTokenJSON(w, accessToken, idToken, refreshToken, expiresIn)

	default:
		writeOAuthError(w, http.StatusBadRequest, "unsupported_grant_type",
			fmt.Sprintf("Grant type '%s' is not supported", grantType))
	}
}
