package cognitoidentityprovider

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

// Concurrent UpdateUserFunc mutations must all land: recordMu serialises
// each read-modify-write cycle, so no caller's change can be discarded by
// another caller's whole-record write.
func TestUpdateUserFuncConcurrentMutationsAllLand(t *testing.T) {
	s := newUserPoolTestStore(t)
	pool, err := s.CreateUserPool(NewUserPool("user-serial-updates", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	user := NewUser(pool.ID, "racer")
	// A non-empty attribute map: the record's omitempty JSON drops an empty
	// map, and the mutations below write into the map the read returns.
	user.Attributes = map[string]string{"seed": "1"}
	if err := s.CreateUser(user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	const racers = 8
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < racers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_ = s.UpdateUserFunc(pool.ID, "racer", func(u *User) error {
				u.Attributes[fmt.Sprintf("note%d", i)] = fmt.Sprintf("value%d", i)
				return nil
			})
		}(i)
	}
	close(start)
	wg.Wait()

	got, err := s.GetUser(pool.ID, "racer")
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	for i := 0; i < racers; i++ {
		if got.Attributes[fmt.Sprintf("note%d", i)] != fmt.Sprintf("value%d", i) {
			t.Fatalf("mutation %d lost; attributes: %v", i, got.Attributes)
		}
	}
}

// UpdateUserFunc on a user that does not exist reports ErrUserNotFound and
// never invokes the mutation.
func TestUpdateUserFuncMissingUserReturnsNotFound(t *testing.T) {
	s := newUserPoolTestStore(t)
	pool, err := s.CreateUserPool(NewUserPool("missing-user", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	mutated := false
	if err := s.UpdateUserFunc(pool.ID, "ghost", func(u *User) error {
		mutated = true
		return nil
	}); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("UpdateUserFunc on missing user returned %v, want ErrUserNotFound", err)
	}
	if mutated {
		t.Fatal("mutation ran for a missing user")
	}
}

// A concurrent attribute update cannot discard a membership change: both
// writes hold recordMu, so the attribute update re-reads the member list
// instead of overwriting it with a stale snapshot.
func TestUpdateGroupFuncConcurrentUpdateKeepsMembership(t *testing.T) {
	s := newUserPoolTestStore(t)
	pool, err := s.CreateUserPool(NewUserPool("group-serial-updates", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	if err := s.CreateGroup(&Group{UserPoolID: pool.ID, Name: "raced"}); err != nil {
		t.Fatalf("create group: %v", err)
	}
	member := NewUser(pool.ID, "member")
	member.Attributes = map[string]string{}
	if err := s.CreateUser(member); err != nil {
		t.Fatalf("create member: %v", err)
	}

	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		_ = s.UpdateGroupFunc(pool.ID, "raced", func(g *Group) error {
			g.Description = "updated"
			return nil
		})
	}()
	go func() {
		defer wg.Done()
		<-start
		_ = s.AddUserToGroup(pool.ID, "raced", "member")
	}()
	close(start)
	wg.Wait()

	got, err := s.GetGroup(pool.ID, "raced")
	if err != nil {
		t.Fatalf("get group: %v", err)
	}
	if got.Description != "updated" {
		t.Fatalf("attribute update lost: %q", got.Description)
	}
	if len(got.Members) != 1 || got.Members[0] != "member" {
		t.Fatalf("membership lost: %v", got.Members)
	}
}
