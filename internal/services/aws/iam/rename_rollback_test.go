package iam

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"vorpalstacks/internal/core/storage"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// A reverse migration step that fails during RenameUser's rollback must
// surface alongside the forward failure: the returned error carries the
// joined pair, so the split state between the two names is observable by
// the caller instead of hidden behind the original forward error.
func TestRenameUserJoinsReverseMigrationFailureIntoRollbackError(t *testing.T) {
	const oldName = "rollback-user"
	const newName = "rollback-user-new"

	// Two targeted faults: the membership listing fails the forward
	// migration, and the login-profile restore Put — issued only by the
	// reverse of that step, keyed by the old name — fails on the way
	// back. The restore fault arms after the fixtures are laid down (the
	// fixture Put targets the same key), and the forward login-profile
	// Put targets the new name, so the migration reaches the faulted
	// forward step first.
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })

	profilesFault := &readFault{err: errors.New("simulated login profile restore failure")}
	faulted := faultingStorage{
		inner: faultingStorage{
			inner:      st,
			bucketName: "iam_user_groups",
			fault: &readFault{
				failScan: true,
				err:      errors.New("simulated membership listing failure"),
			},
		},
		bucketName: "iam_login_profiles",
		fault:      profilesFault,
	}
	store := iamstore.NewIAMStore(faulted, "123456789012")

	if _, err := store.Users().Create(oldName, "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if _, err := store.LoginProfiles().Create(oldName, "P@ssw0rd-Rollback1", false); err != nil {
		t.Fatalf("create login profile: %v", err)
	}
	if _, err := store.AccessKeys().Create(oldName); err != nil {
		t.Fatalf("create access key: %v", err)
	}
	profilesFault.failPutKey = oldName

	err = store.RenameUser(oldName, newName, "")
	if err == nil {
		t.Fatal("expected the rename to fail on the faulted membership listing")
	}
	assert.Contains(t, err.Error(), "simulated membership listing failure")
	assert.Contains(t, err.Error(), "simulated login profile restore failure")

	// The reverse steps that could run all ran: the access key is back
	// under the old name, while the login profile — whose restore failed —
	// remains under the new one, matching the joined-error split state.
	keys, err := store.AccessKeys().ListByUserName(oldName)
	if err != nil {
		t.Fatalf("list access keys: %v", err)
	}
	assert.Len(t, keys, 1)
	if !store.LoginProfiles().Exists(newName) {
		t.Fatal("the unrestored login profile must remain under the new name")
	}
}

// A migration step that fails mid-way through a group rename must leave
// the already-migrated memberships and policy keys back under the old
// name: the reverse steps run in reverse order and no new group record
// survives, so the old group is complete and usable again.
func TestRenameGroupReversesMigratedResources(t *testing.T) {
	const oldName = "rollback-group"
	const newName = "rollback-group-new"

	// The attached-policy migration is the faulted forward step; the
	// memberships and the inline-policy migration run before it and must
	// be reversed on the way back.
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	faulted := faultingStorage{
		inner:      st,
		bucketName: "iam_attached_policies",
		fault: &readFault{
			failScan: true,
			err:      errors.New("simulated attached policy migration failure"),
		},
	}
	store := iamstore.NewIAMStore(faulted, "123456789012")

	if _, err := store.Groups().Create(oldName, "/", "123456789012"); err != nil {
		t.Fatalf("create group: %v", err)
	}
	for _, userName := range []string{"group-member-a", "group-member-b"} {
		if _, err := store.Users().Create(userName, "/", "123456789012", nil); err != nil {
			t.Fatalf("create user %s: %v", userName, err)
		}
		if err := store.UserGroups().AddUserToGroup(userName, oldName); err != nil {
			t.Fatalf("add member %s: %v", userName, err)
		}
	}
	if err := store.InlinePolicies().Put("group", oldName, "InlinePolicy", `{"Version":"2012-10-17","Statement":[]}`); err != nil {
		t.Fatalf("put inline policy: %v", err)
	}

	err = store.RenameGroup(oldName, newName, "")
	if err == nil {
		t.Fatal("expected the rename to fail on the faulted attached-policy migration")
	}
	assert.Contains(t, err.Error(), "simulated attached policy migration failure")

	// Everything the forward pass moved is back under the old name, and
	// the faulted step's own resources never left it.
	members, err := store.UserGroups().ListUsersInGroup(oldName)
	if err != nil {
		t.Fatalf("list members: %v", err)
	}
	assert.ElementsMatch(t, []string{"group-member-a", "group-member-b"}, members)
	if newMembers, err := store.UserGroups().ListUsersInGroup(newName); err != nil {
		t.Fatalf("list new-name members: %v", err)
	} else {
		assert.Empty(t, newMembers, "no membership may remain under the dead new name")
	}

	inline, err := store.InlinePolicies().List("group", oldName)
	if err != nil {
		t.Fatalf("list inline policies: %v", err)
	}
	assert.Equal(t, []string{"InlinePolicy"}, inline)
	if newInline, err := store.InlinePolicies().List("group", newName); err != nil {
		t.Fatalf("list new-name inline policies: %v", err)
	} else {
		assert.Empty(t, newInline)
	}

	assert.True(t, store.Groups().Exists(oldName), "the old group record must remain")
	assert.False(t, store.Groups().Exists(newName), "no new group record may survive a failed rename")
}
