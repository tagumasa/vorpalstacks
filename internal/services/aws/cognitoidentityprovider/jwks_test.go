package cognitoidentityprovider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// The JWKS endpoint is per-pool: a request without a userPoolId is rejected
// instead of answered with an arbitrary pool's keys, an unknown pool is a
// 404, and a real pool serves its signing key.
func TestJWKSHandlerRequiresPoolAndServesPerPool(t *testing.T) {
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewCognitoService("000000000000", "us-east-1")
	svc.SetStorageManager(mgr)
	reqCtx := request.NewRequestContext(context.Background(), mgr, "000000000000", "us-east-1")
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("jwkspool", "us-east-1"))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(svc.JWKSHandler)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/.well-known/jwks.json", nil))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("request without userPoolId: status %d, want 400", rec.Code)
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/.well-known/jwks.json?userPoolId=us-east-1_nonexistent", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("request for unknown pool: status %d, want 404", rec.Code)
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/.well-known/jwks.json?userPoolId="+pool.ID, nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("request for existing pool: status %d, want 200 (body %s)", rec.Code, rec.Body.String())
	}
	var jwks struct {
		Keys []map[string]interface{} `json:"keys"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &jwks); err != nil {
		t.Fatalf("decode jwks response: %v", err)
	}
	if len(jwks.Keys) != 1 {
		t.Fatalf("expected exactly one signing key, got %d", len(jwks.Keys))
	}
}
