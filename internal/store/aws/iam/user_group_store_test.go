package iam

import (
	"fmt"
	"sync"
	"testing"

	"vorpalstacks/internal/core/storage"
)

func userGroupTestStore(t *testing.T) *IAMStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return NewIAMStore(st, "123456789012")
}

// Every mutation of one user's membership set runs under the user's
// membership lock, so a group rename migration concurrent with
// quota-checked additions can never leave the user over quota: the
// migration's add/delete pair is atomic against the quota read.
func TestMembershipQuotaHoldsUnderConcurrentRename(t *testing.T) {
	store := userGroupTestStore(t)
	const user = "lock-user"
	if _, err := store.Users().Create(user, "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}
	// groupNames[1] never exists up front: the ping-pong rename creates
	// it on the first hop and the groups alternate names thereafter.
	groupNames := make([]string, MaxIAMGroupsPerUser+2)
	for i := range groupNames {
		groupNames[i] = fmt.Sprintf("lock-grp-%02d", i)
		if i == 1 {
			continue
		}
		if _, err := store.Groups().Create(groupNames[i], "/", "123456789012"); err != nil {
			t.Fatalf("create group: %v", err)
		}
	}

	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Repeatedly rename group[0] to group[1] and back — each migration
	// moves the user's membership under the user's lock.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			from, to := groupNames[0], groupNames[1]
			if i%2 == 1 {
				from, to = to, from
			}
			if err := store.RenameGroup(from, to, ""); err != nil {
				t.Errorf("rename: %v", err)
				return
			}
		}
		close(stop)
	}()

	// Concurrently drive the user to the quota through the checked add
	// path; every accepted add must keep the user within the quota.
	for _, group := range groupNames[2:] {
		wg.Add(1)
		go func(group string) {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				if err := store.UserGroups().AddUserToGroupWithLimit(user, group, MaxIAMGroupsPerUser); err != nil {
					return // quota reached or raced — both legitimate
				}
			}
		}(group)
	}
	wg.Wait()

	groups, err := store.UserGroups().ListGroupsForUser(user)
	if err != nil {
		t.Fatalf("list groups: %v", err)
	}
	if len(groups) > MaxIAMGroupsPerUser {
		t.Fatalf("quota violated: user holds %d groups, limit %d: %v", len(groups), MaxIAMGroupsPerUser, groups)
	}
}

// MigrateUserGroup moves one membership add-first under the user's lock;
// the from-membership is gone and the to-membership carries the user.
func TestMigrateUserGroup(t *testing.T) {
	store := userGroupTestStore(t)
	if _, err := store.Users().Create("mig-user", "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if _, err := store.Groups().Create("mig-from", "/", "123456789012"); err != nil {
		t.Fatalf("create group: %v", err)
	}
	if _, err := store.Groups().Create("mig-to", "/", "123456789012"); err != nil {
		t.Fatalf("create group: %v", err)
	}
	if err := store.UserGroups().AddUserToGroup("mig-user", "mig-from"); err != nil {
		t.Fatalf("add: %v", err)
	}

	if err := store.UserGroups().MigrateUserGroup("mig-user", "mig-from", "mig-to"); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if store.UserGroups().IsUserInGroup("mig-user", "mig-from") {
		t.Fatal("the from-membership must be gone")
	}
	if !store.UserGroups().IsUserInGroup("mig-user", "mig-to") {
		t.Fatal("the to-membership must exist")
	}

	if err := store.UserGroups().MigrateUserGroup("mig-user", "mig-from", "mig-to"); err == nil {
		t.Fatal("migrating an absent membership must fail")
	}
}
