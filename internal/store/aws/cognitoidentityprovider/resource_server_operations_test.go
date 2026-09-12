package cognitoidentityprovider

import (
	"errors"
	"testing"
)

// A missing resource server or identity provider reports its own sentinel —
// not the pool's — so callers can tell which entity the request named.
func TestMissingEntitiesReportTheirOwnNotFoundSentinels(t *testing.T) {
	s := newUserPoolTestStore(t)
	pool, err := s.CreateUserPool(NewUserPool("sentinel-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	if _, err := s.GetResourceServer(pool.ID, "missing-server"); !errors.Is(err, ErrResourceServerNotFound) {
		t.Fatalf("GetResourceServer: got %v, want ErrResourceServerNotFound", err)
	}
	if err := s.DeleteResourceServer(pool.ID, "missing-server"); !errors.Is(err, ErrResourceServerNotFound) {
		t.Fatalf("DeleteResourceServer: got %v, want ErrResourceServerNotFound", err)
	}
	if _, err := s.GetIdentityProvider(pool.ID, "MissingProvider"); !errors.Is(err, ErrIdentityProviderNotFound) {
		t.Fatalf("GetIdentityProvider: got %v, want ErrIdentityProviderNotFound", err)
	}
	if err := s.DeleteIdentityProvider(pool.ID, "MissingProvider"); !errors.Is(err, ErrIdentityProviderNotFound) {
		t.Fatalf("DeleteIdentityProvider: got %v, want ErrIdentityProviderNotFound", err)
	}

	// An existing entity still resolves through both paths.
	rs := &ResourceServer{UserPoolID: pool.ID, Identifier: "https://example.com", Name: "example"}
	if err := s.CreateResourceServer(rs); err != nil {
		t.Fatalf("create resource server: %v", err)
	}
	if _, err := s.GetResourceServer(pool.ID, rs.Identifier); err != nil {
		t.Fatalf("get existing resource server: %v", err)
	}
	if err := s.DeleteResourceServer(pool.ID, rs.Identifier); err != nil {
		t.Fatalf("delete existing resource server: %v", err)
	}
}
