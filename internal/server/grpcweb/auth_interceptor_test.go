package grpcweb

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"vorpalstacks/pkg/vsjwt"
)

// TestAuthHTTPMiddlewareRejectsIDToken verifies that the HTTP auth middleware
// rejects ID tokens (token_use="id") even when the signature and issuer are
// valid. Only access tokens should be accepted for admin API authentication.
func TestAuthHTTPMiddlewareRejectsIDToken(t *testing.T) {
	privateKey, err := vsjwt.GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	issuer := "https://test-issuer.com"
	manager := vsjwt.MustNewManager(privateKey, "key-id", issuer)

	user := &testJWTUser{
		id:       "user-1",
		username: "test",
	}

	idToken, err := manager.GenerateIDToken(user, "client", 3600)
	if err != nil {
		t.Fatalf("failed to generate ID token: %v", err)
	}

	accessToken, err := manager.GenerateAccessToken(user, "client", 3600)
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := newAuthHTTPMiddleware(manager, handler)

	// ID token should be rejected
	req := httptest.NewRequest("POST", "/admin_config.AdminConfigService/GetConfig", nil)
	req.Header.Set("Authorization", "Bearer "+idToken)
	rec := httptest.NewRecorder()
	mw.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("ID token: expected 401, got %d", rec.Code)
	}

	// Access token should be accepted
	req2 := httptest.NewRequest("POST", "/admin_config.AdminConfigService/GetConfig", nil)
	req2.Header.Set("Authorization", "Bearer "+accessToken)
	rec2 := httptest.NewRecorder()
	mw.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Errorf("Access token: expected 200, got %d", rec2.Code)
	}
}

// TestAuthHTTPMiddlewareNoAuthPath verifies that admin auth paths bypass
// authentication entirely.
func TestAuthHTTPMiddlewareNoAuthPath(t *testing.T) {
	privateKey, _ := vsjwt.GenerateRSAKeyPair()
	manager := vsjwt.MustNewManager(privateKey, "key-id", "https://test-issuer.com")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := newAuthHTTPMiddleware(manager, handler)

	req := httptest.NewRequest("POST", "/admin_auth.AdminAuthService/Login", nil)
	rec := httptest.NewRecorder()
	mw.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("no-auth path: expected 200, got %d", rec.Code)
	}
}

// TestAuthHTTPMiddlewareMissingToken verifies that requests without a token
// are rejected.
func TestAuthHTTPMiddlewareMissingToken(t *testing.T) {
	privateKey, _ := vsjwt.GenerateRSAKeyPair()
	manager := vsjwt.MustNewManager(privateKey, "key-id", "https://test-issuer.com")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := newAuthHTTPMiddleware(manager, handler)

	req := httptest.NewRequest("POST", "/admin_config.AdminConfigService/GetConfig", nil)
	rec := httptest.NewRecorder()
	mw.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("missing token: expected 401, got %d", rec.Code)
	}
}

type testJWTUser struct {
	id       string
	username string
}

func (u *testJWTUser) GetID() string                           { return u.id }
func (u *testJWTUser) GetUsername() string                     { return u.username }
func (u *testJWTUser) GetGroups() []string                     { return nil }
func (u *testJWTUser) GetEmail() string                        { return "" }
func (u *testJWTUser) GetCustomClaims() map[string]interface{} { return nil }
