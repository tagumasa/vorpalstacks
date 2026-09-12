package cognitoidentity

import (
	"errors"
	"testing"

	"vorpalstacks/internal/common/request"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
)

// Updating a pool that no longer exists must surface
// ResourceNotFoundException, not InternalErrorException: the not-found
// sentinel raised inside the locked mutation propagates through the Core
// mapping.
func TestUpdateIdentityPoolCoreDeletedPoolMapsResourceNotFound(t *testing.T) {
	svc, real := newMergeTestService(t)
	pool := cognitoidentitystore.NewIdentityPool("gone-pool", false, "us-east-1")
	if _, err := real.CreateIdentityPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}
	if err := real.DeleteIdentityPool(pool.ID); err != nil {
		t.Fatalf("delete pool: %v", err)
	}

	_, err := svc.updateIdentityPoolCore(&request.RequestContext{Region: "us-east-1"}, UpdateIdentityPoolInput{
		IdentityPoolID:      pool.ID,
		PoolName:            "gone-pool-renamed",
		AllowUnauthProvided: true,
		AllowUnauthRaw:      false,
	})
	if err == nil {
		t.Fatal("update of a deleted pool unexpectedly succeeded")
	}
	if !errors.Is(err, ErrResourceNotFound) {
		t.Fatalf("update of a deleted pool returned %v, want ResourceNotFoundException", err)
	}
}
