package cognitoidentity

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

// Deleting a pool concurrently with identity creation must not leave orphaned
// identity records or index entries under the deleted pool.
func TestDeleteIdentityPoolConcurrentGetIdLeavesNoOrphans(t *testing.T) {
	s := newTestStore(t)
	const rounds = 60
	for round := 0; round < rounds; round++ {
		pool := NewIdentityPool(fmt.Sprintf("orphan-%d", round), false, "us-east-1")
		if _, err := s.CreateIdentityPool(pool); err != nil {
			t.Fatalf("round %d: create pool: %v", round, err)
		}

		var wg sync.WaitGroup
		start := make(chan struct{})
		createdID := ""
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-start
			_ = s.DeleteIdentityPool(pool.ID)
		}()
		go func() {
			defer wg.Done()
			<-start
			identity, err := s.GetOrCreateIdentityByLogins(pool.ID, map[string]string{
				"graph.facebook.com": fmt.Sprintf("token-%d", round),
			})
			if err == nil {
				createdID = identity.ID
			}
		}()
		close(start)
		wg.Wait()

		ids, _, err := s.ListIdentitiesByPool(pool.ID, 60, "")
		if err != nil {
			t.Fatalf("round %d: list identities: %v", round, err)
		}
		if len(ids) != 0 {
			t.Fatalf("round %d: %d identity records survived the pool deletion", round, len(ids))
		}
		if createdID != "" {
			if _, err := s.GetIdentityByID(createdID); !errors.Is(err, ErrIdentityNotFound) {
				t.Fatalf("round %d: identity %s still resolvable after the pool deletion: %v", round, createdID, err)
			}
		}
	}
}

// A pool update racing the pool deletion must not resurrect the pool record.
func TestDeleteIdentityPoolConcurrentUpdateLeavesNoResurrectedPool(t *testing.T) {
	s := newTestStore(t)
	const rounds = 60
	for round := 0; round < rounds; round++ {
		pool := NewIdentityPool(fmt.Sprintf("resurrect-%d", round), false, "us-east-1")
		if _, err := s.CreateIdentityPool(pool); err != nil {
			t.Fatalf("round %d: create pool: %v", round, err)
		}

		var wg sync.WaitGroup
		start := make(chan struct{})
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-start
			_ = s.DeleteIdentityPool(pool.ID)
		}()
		go func() {
			defer wg.Done()
			<-start
			if p, err := s.GetIdentityPool(pool.ID); err == nil {
				p.DeveloperProviderName = "login.example.com"
				_ = s.UpdateIdentityPool(p)
			}
		}()
		close(start)
		wg.Wait()

		if _, err := s.GetIdentityPool(pool.ID); !errors.Is(err, ErrIdentityPoolNotFound) {
			t.Fatalf("round %d: pool record survived or resurrected after deletion: %v", round, err)
		}
	}
}
