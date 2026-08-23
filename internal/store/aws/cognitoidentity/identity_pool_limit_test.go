package cognitoidentity

import (
	"errors"
	"fmt"
	"testing"

	"vorpalstacks/internal/core/storage"
)

func newTestStore(t *testing.T) *CognitoIdentityStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return NewCognitoIdentityStore(st, "000000000000", "us-east-1")
}

// AWS documents a default resource quota of 1,000 identity pools per account
// ("Quotas in Amazon Cognito"). Pools below the quota must all be accepted and
// creation past the quota must be rejected.
func TestCreateIdentityPoolEnforcesAccountQuota(t *testing.T) {
	s := newTestStore(t)
	for i := 0; i < MaxIdentityPoolsPerAccount; i++ {
		pool := NewIdentityPool(fmt.Sprintf("pool-%d", i), false, "us-east-1")
		if _, err := s.CreateIdentityPool(pool); err != nil {
			t.Fatalf("pool %d within the quota was rejected: %v", i, err)
		}
	}

	excess := NewIdentityPool("pool-over-quota", false, "us-east-1")
	_, err := s.CreateIdentityPool(excess)
	if !errors.Is(err, ErrTooManyIdentityPools) {
		t.Fatalf("creation past the quota returned %v, want ErrTooManyIdentityPools", err)
	}
}
