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
			_ = s.UpdateIdentityPoolFunc(pool.ID, func(p *IdentityPool) error {
				p.DeveloperProviderName = "login.example.com"
				return nil
			})
		}()
		close(start)
		wg.Wait()

		if _, err := s.GetIdentityPool(pool.ID); !errors.Is(err, ErrIdentityPoolNotFound) {
			t.Fatalf("round %d: pool record survived or resurrected after deletion: %v", round, err)
		}
	}
}

// Concurrent UpdateIdentityPoolFunc mutations must all land: the pool lock
// serialises each read-modify-write cycle, so no caller's field change can be
// discarded by another caller's write.
func TestUpdateIdentityPoolFuncConcurrentMutationsAllLand(t *testing.T) {
	s := newTestStore(t)
	pool := NewIdentityPool("serial-updates", false, "us-east-1")
	if _, err := s.CreateIdentityPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}

	const mutations = 20
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < mutations; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_ = s.UpdateIdentityPoolFunc(pool.ID, func(p *IdentityPool) error {
				if p.SupportedLoginProviders == nil {
					p.SupportedLoginProviders = make(map[string]string)
				}
				p.SupportedLoginProviders[fmt.Sprintf("provider-%d", i)] = fmt.Sprintf("token-%d", i)
				return nil
			})
		}(i)
	}
	close(start)
	wg.Wait()

	got, err := s.GetIdentityPool(pool.ID)
	if err != nil {
		t.Fatalf("get pool after mutations: %v", err)
	}
	if len(got.SupportedLoginProviders) != mutations {
		t.Fatalf("lost updates: %d of %d provider entries survived", len(got.SupportedLoginProviders), mutations)
	}
}

// A mutation racing the pool deletion must observe ErrIdentityPoolNotFound
// rather than silently resurrecting the deleted record.
func TestUpdateIdentityPoolFuncDeleteRaceReturnsPoolNotFound(t *testing.T) {
	s := newTestStore(t)
	const rounds = 60
	for round := 0; round < rounds; round++ {
		pool := NewIdentityPool(fmt.Sprintf("func-race-%d", round), false, "us-east-1")
		if _, err := s.CreateIdentityPool(pool); err != nil {
			t.Fatalf("round %d: create pool: %v", round, err)
		}

		var wg sync.WaitGroup
		start := make(chan struct{})
		var updateErr error
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-start
			_ = s.DeleteIdentityPool(pool.ID)
		}()
		go func() {
			defer wg.Done()
			<-start
			updateErr = s.UpdateIdentityPoolFunc(pool.ID, func(p *IdentityPool) error {
				p.DeveloperProviderName = "login.example.com"
				return nil
			})
		}()
		close(start)
		wg.Wait()

		if _, err := s.GetIdentityPool(pool.ID); !errors.Is(err, ErrIdentityPoolNotFound) {
			t.Fatalf("round %d: pool record survived or resurrected after deletion: %v", round, err)
		}
		if updateErr != nil && !errors.Is(updateErr, ErrIdentityPoolNotFound) {
			t.Fatalf("round %d: unexpected update error: %v", round, updateErr)
		}
	}
}
