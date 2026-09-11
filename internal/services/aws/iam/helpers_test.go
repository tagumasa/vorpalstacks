package iam

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vorpalstacks/internal/core/storage"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// The permissions-boundary usage tallies move only after the entity
// persisted: a failed put leaves the entity and both policies' counts
// exactly as they were, so a policy cannot pass DeletePolicy's conflict
// check while an entity still references it. On success the previous
// policy's tally decrements and the new one increments; a same-ARN attach
// is idempotent.
func TestBoundaryUsageCountsFollowPersistedState(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()
	store := iamstore.NewIAMStore(st, "123456789012")

	oldArn := "arn:aws:iam::123456789012:policy/boundary-old"
	newArn := "arn:aws:iam::123456789012:policy/boundary-new"
	for _, policy := range []*iamstore.Policy{
		{Arn: oldArn, PolicyName: "boundary-old", AccountId: "123456789012", Path: "/", PermissionsBoundaryUsageCount: 1},
		{Arn: newArn, PolicyName: "boundary-new", AccountId: "123456789012", Path: "/"},
	} {
		if err := store.Policies().Put(policy); err != nil {
			t.Fatalf("put policy %s: %v", policy.PolicyName, err)
		}
	}

	userBoundary := func(u *iamstore.User) **iamstore.PermissionsBoundary { return &u.PermissionsBoundary }
	count := func(t *testing.T, arn string) int {
		t.Helper()
		policy, err := store.Policies().Get(arn)
		require.NoError(t, err)
		return policy.PermissionsBoundaryUsageCount
	}

	t.Run("failed attach put leaves both counts unchanged", func(t *testing.T) {
		user := &iamstore.User{UserName: "bound-user", PermissionsBoundary: &iamstore.PermissionsBoundary{
			PermissionsBoundaryType: "Policy", PermissionsBoundaryArn: oldArn,
		}}
		putFails := func(*iamstore.User) error { return errors.New("simulated entity persist failure") }

		err := attachPermissionsBoundaryCore(store, user, newArn, userBoundary, putFails)
		require.Error(t, err)
		assert.Equal(t, 1, count(t, oldArn), "the previous policy's tally must not move before the entity persisted")
		assert.Equal(t, 0, count(t, newArn), "the new policy's tally must not move before the entity persisted")
	})

	t.Run("failed delete put leaves the count unchanged", func(t *testing.T) {
		user := &iamstore.User{UserName: "bound-user", PermissionsBoundary: &iamstore.PermissionsBoundary{
			PermissionsBoundaryType: "Policy", PermissionsBoundaryArn: oldArn,
		}}
		putFails := func(*iamstore.User) error { return errors.New("simulated entity persist failure") }

		err := deletePermissionsBoundaryCore(store, user, userBoundary, putFails)
		require.Error(t, err)
		assert.Equal(t, 1, count(t, oldArn), "the bound policy's tally must not move before the entity persisted")
	})

	t.Run("successful attach moves both counts", func(t *testing.T) {
		user := &iamstore.User{UserName: "bound-user", PermissionsBoundary: &iamstore.PermissionsBoundary{
			PermissionsBoundaryType: "Policy", PermissionsBoundaryArn: oldArn,
		}}
		put := func(u *iamstore.User) error { return store.Users().Put(u) }

		require.NoError(t, attachPermissionsBoundaryCore(store, user, newArn, userBoundary, put))
		assert.Equal(t, 0, count(t, oldArn))
		assert.Equal(t, 1, count(t, newArn))

		// Same-ARN attach is idempotent: the tallies stay put.
		require.NoError(t, attachPermissionsBoundaryCore(store, user, newArn, userBoundary, put))
		assert.Equal(t, 0, count(t, oldArn))
		assert.Equal(t, 1, count(t, newArn))
	})

	t.Run("successful delete clears the count", func(t *testing.T) {
		user := &iamstore.User{UserName: "bound-user", PermissionsBoundary: &iamstore.PermissionsBoundary{
			PermissionsBoundaryType: "Policy", PermissionsBoundaryArn: newArn,
		}}
		put := func(u *iamstore.User) error { return store.Users().Put(u) }

		require.NoError(t, deletePermissionsBoundaryCore(store, user, userBoundary, put))
		assert.Equal(t, 0, count(t, newArn))
	})
}
