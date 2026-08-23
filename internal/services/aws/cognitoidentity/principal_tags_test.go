package cognitoidentity

import (
	"context"
	"errors"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
)

func newPrincipalTagTestPool(t *testing.T, store *cognitoidentitystore.CognitoIdentityStore) *cognitoidentitystore.IdentityPool {
	t.Helper()
	pool := cognitoidentitystore.NewIdentityPool("principal-tag-pool", false, "us-east-1")
	if _, err := store.CreateIdentityPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}
	return pool
}

// Principal tag names must respect the model's PrincipalTagID length bound
// of 1-128 characters on the token-issuing path.
func TestGetOpenIdTokenForDeveloperIdentityRejectsOversizedPrincipalTagName(t *testing.T) {
	svc, real := newMergeTestService(t)
	pool := newPrincipalTagTestPool(t, real)

	_, err := svc.GetOpenIdTokenForDeveloperIdentity(context.Background(),
		&request.RequestContext{Region: "us-east-1"},
		&request.ParsedRequest{Parameters: map[string]interface{}{
			"IdentityPoolId": pool.ID,
			"Logins":         map[string]interface{}{"login.example.com": "user-1"},
			"PrincipalTags":  map[string]interface{}{strings.Repeat("a", 129): "value"},
		}})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("oversized principal tag name returned %v, want ErrInvalidParameter", err)
	}
}

// Principal tag names must respect the model's PrincipalTagID length bound
// of 1-128 characters on the attribute-map path.
func TestSetPrincipalTagAttributeMapRejectsOversizedPrincipalTagName(t *testing.T) {
	svc, real := newMergeTestService(t)
	pool := newPrincipalTagTestPool(t, real)

	_, err := svc.SetPrincipalTagAttributeMap(context.Background(),
		&request.RequestContext{Region: "us-east-1"},
		&request.ParsedRequest{Parameters: map[string]interface{}{
			"IdentityPoolId":       pool.ID,
			"IdentityProviderName": "login.example.com",
			"PrincipalTags":        map[string]interface{}{strings.Repeat("a", 129): "value"},
		}})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("oversized principal tag name returned %v, want ErrInvalidParameter", err)
	}
}

// Tag keys on identity pool updates must pass the same key validation as
// pool creation.
func TestUpdateIdentityPoolRejectsOversizedTagKey(t *testing.T) {
	svc, real := newMergeTestService(t)
	pool := newPrincipalTagTestPool(t, real)

	_, err := svc.UpdateIdentityPool(context.Background(),
		&request.RequestContext{Region: "us-east-1"},
		&request.ParsedRequest{Parameters: map[string]interface{}{
			"IdentityPoolId":                 pool.ID,
			"IdentityPoolName":               "principal-tag-pool",
			"AllowUnauthenticatedIdentities": true,
			"IdentityPoolTags":               map[string]interface{}{strings.Repeat("t", 129): "value"},
		}})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("oversized tag key returned %v, want ErrInvalidParameter", err)
	}
}
