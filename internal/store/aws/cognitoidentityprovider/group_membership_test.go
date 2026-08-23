package cognitoidentityprovider

import (
	"errors"
	"testing"
)

// seedGroupMemberPool creates a pool, a user stored under the canonical
// username, and a group, returning the pool and group names.
func seedGroupMemberPool(t *testing.T, s *CognitoStore, username string) (poolID, groupName string) {
	t.Helper()
	pool, err := s.CreateUserPool(NewUserPool("membership-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	user := NewUser(pool.ID, username)
	if err := s.CreateUser(user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	group := NewGroup(pool.ID, "staff")
	if err := s.CreateGroup(group); err != nil {
		t.Fatalf("create group: %v", err)
	}
	return pool.ID, "staff"
}

// A case-insensitive pool treats differently-cased spellings as the same
// user, so the membership record must hold the canonical stored username.
func TestAddUserToGroupStoresCanonicalUsername(t *testing.T) {
	s := newImportJobTestStore(t)
	poolID, groupName := seedGroupMemberPool(t, s, "alice")

	if err := s.AddUserToGroup(poolID, groupName, "ALICE"); err != nil {
		t.Fatalf("add with different case: %v", err)
	}

	group, err := s.GetGroup(poolID, groupName)
	if err != nil {
		t.Fatalf("load group: %v", err)
	}
	if len(group.Members) != 1 || group.Members[0] != "alice" {
		t.Fatalf("members = %v, want [alice]", group.Members)
	}
	user, err := s.GetUser(poolID, "alice")
	if err != nil {
		t.Fatalf("load user: %v", err)
	}
	if len(user.Groups) != 1 || user.Groups[0] != groupName {
		t.Fatalf("user groups = %v, want [%s]", user.Groups, groupName)
	}
}

// The duplicate check must match across case variants of the same user.
func TestAddUserToGroupDuplicateDetectionIsCaseInsensitive(t *testing.T) {
	s := newImportJobTestStore(t)
	poolID, groupName := seedGroupMemberPool(t, s, "alice")

	if err := s.AddUserToGroup(poolID, groupName, "alice"); err != nil {
		t.Fatalf("initial add: %v", err)
	}
	err := s.AddUserToGroup(poolID, groupName, "ALICE")
	if !errors.Is(err, ErrUserAlreadyInGroup) {
		t.Fatalf("second add error = %v, want ErrUserAlreadyInGroup", err)
	}
	group, err := s.GetGroup(poolID, groupName)
	if err != nil {
		t.Fatalf("load group: %v", err)
	}
	if len(group.Members) != 1 {
		t.Fatalf("members = %v, want a single entry", group.Members)
	}
}

// A nonexistent user must fail without leaving a member record behind.
func TestAddUserToGroupNonexistentUserLeavesNoMember(t *testing.T) {
	s := newImportJobTestStore(t)
	poolID, groupName := seedGroupMemberPool(t, s, "alice")

	err := s.AddUserToGroup(poolID, groupName, "ghost-user")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("error = %v, want ErrUserNotFound", err)
	}
	group, err := s.GetGroup(poolID, groupName)
	if err != nil {
		t.Fatalf("load group: %v", err)
	}
	if len(group.Members) != 0 {
		t.Fatalf("members = %v, want none", group.Members)
	}
}

// Removing membership must resolve the username like every other user
// lookup, so a differently-cased spelling still finds the member.
func TestRemoveUserFromGroupResolvesUsernameCase(t *testing.T) {
	s := newImportJobTestStore(t)
	poolID, groupName := seedGroupMemberPool(t, s, "alice")

	if err := s.AddUserToGroup(poolID, groupName, "alice"); err != nil {
		t.Fatalf("add: %v", err)
	}
	if err := s.RemoveUserFromGroup(poolID, groupName, "ALICE"); err != nil {
		t.Fatalf("remove with different case: %v", err)
	}

	group, err := s.GetGroup(poolID, groupName)
	if err != nil {
		t.Fatalf("load group: %v", err)
	}
	if len(group.Members) != 0 {
		t.Fatalf("members = %v, want none", group.Members)
	}
	user, err := s.GetUser(poolID, "alice")
	if err != nil {
		t.Fatalf("load user: %v", err)
	}
	if len(user.Groups) != 0 {
		t.Fatalf("user groups = %v, want none", user.Groups)
	}
}

// Deleting a user must clear a membership that was created through a
// differently-cased spelling — no ghost member may survive.
func TestDeleteUserRemovesMembershipAddedWithDifferentCase(t *testing.T) {
	s := newImportJobTestStore(t)
	poolID, groupName := seedGroupMemberPool(t, s, "alice")

	if err := s.AddUserToGroup(poolID, groupName, "ALICE"); err != nil {
		t.Fatalf("add with different case: %v", err)
	}
	if err := s.DeleteUser(poolID, "alice"); err != nil {
		t.Fatalf("delete user: %v", err)
	}

	users, err := s.ListUsersInGroup(poolID, groupName)
	if err != nil {
		t.Fatalf("list users in group: %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("group still lists %d user(s) after deletion", len(users))
	}
	group, err := s.GetGroup(poolID, groupName)
	if err != nil {
		t.Fatalf("load group: %v", err)
	}
	if len(group.Members) != 0 {
		t.Fatalf("members = %v, want none", group.Members)
	}
}
