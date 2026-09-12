package cognitoidentityprovider

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/pkg/vsjwt"
)

func newUserPoolTestStore(t *testing.T) *CognitoStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return NewCognitoStore(st, "000000000000", "us-east-1")
}

// Concurrent UpdateUserPoolFunc mutations must all land: the pool lock
// serialises each read-modify-write cycle, so no caller's field change can be
// discarded by another caller's write.
func TestUpdateUserPoolFuncConcurrentMutationsAllLand(t *testing.T) {
	s := newUserPoolTestStore(t)
	pool, err := s.CreateUserPool(NewUserPool("serial-updates", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	const mutations = 12
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < mutations; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_ = s.UpdateUserPoolFunc(pool.ID, func(p *UserPool) error {
				p.SchemaAttributes = append(p.SchemaAttributes, SchemaAttributeType{
					Name:              fmt.Sprintf("custom:attr-%d", i),
					AttributeDataType: "String",
				})
				return nil
			})
		}(i)
	}
	close(start)
	wg.Wait()

	got, err := s.GetUserPool(pool.ID)
	if err != nil {
		t.Fatalf("get pool after mutations: %v", err)
	}
	if len(got.SchemaAttributes) != mutations {
		t.Fatalf("concurrent mutations lost writes: got %d schema attributes, want %d", len(got.SchemaAttributes), mutations)
	}
}

// A pool mutation racing the pool deletion must surface ErrUserPoolNotFound
// rather than silently resurrecting the deleted record.
func TestUpdateUserPoolFuncDeleteRaceReturnsPoolNotFound(t *testing.T) {
	s := newUserPoolTestStore(t)
	// Each round creates its own pool, and pool creation generates a fresh
	// RSA signing key, so the round count is kept low to bound the runtime.
	const rounds = 8
	for round := 0; round < rounds; round++ {
		pool, err := s.CreateUserPool(NewUserPool(fmt.Sprintf("race-%d", round), "us-east-1"))
		if err != nil {
			t.Fatalf("round %d: create pool: %v", round, err)
		}

		var wg sync.WaitGroup
		start := make(chan struct{})
		var updateErr error
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-start
			_ = s.DeleteUserPool(pool.ID)
		}()
		go func() {
			defer wg.Done()
			<-start
			updateErr = s.UpdateUserPoolFunc(pool.ID, func(p *UserPool) error {
				p.MfaConfiguration = "ON"
				return nil
			})
		}()
		close(start)
		wg.Wait()

		if updateErr != nil && !errors.Is(updateErr, ErrUserPoolNotFound) {
			t.Fatalf("round %d: update returned %v, want nil or ErrUserPoolNotFound", round, updateErr)
		}
		if _, err := s.GetUserPool(pool.ID); !errors.Is(err, ErrUserPoolNotFound) {
			t.Fatalf("round %d: pool record survived or resurrected after deletion: %v", round, err)
		}
	}
}

// Mutating a pool that does not exist must report ErrUserPoolNotFound, and a
// read-modify-write after deletion must not resurrect the record either.
func TestUpdateUserPoolFuncMissingPoolReturnsNotFound(t *testing.T) {
	s := newUserPoolTestStore(t)
	if err := s.UpdateUserPoolFunc("us-east-1_missing", func(p *UserPool) error { return nil }); !errors.Is(err, ErrUserPoolNotFound) {
		t.Fatalf("UpdateUserPoolFunc on missing pool returned %v, want ErrUserPoolNotFound", err)
	}

	pool, err := s.CreateUserPool(NewUserPool("gone", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	if err := s.DeleteUserPool(pool.ID); err != nil {
		t.Fatalf("delete pool: %v", err)
	}
	if err := s.UpdateUserPoolFunc(pool.ID, func(p *UserPool) error { return nil }); !errors.Is(err, ErrUserPoolNotFound) {
		t.Fatalf("UpdateUserPoolFunc after deletion returned %v, want ErrUserPoolNotFound", err)
	}
	if s.Exists(pool.ID) {
		t.Fatal("UpdateUserPoolFunc resurrected the deleted pool record")
	}
}

// A created pool must carry decodable key material: the public key is what
// every token validation and the JWKS route read, so an encode failure at
// creation has to abort the pool rather than persist an empty key.
func TestCreateUserPoolKeyMaterialDecodable(t *testing.T) {
	s := newUserPoolTestStore(t)
	pool, err := s.CreateUserPool(NewUserPool("keymaterial", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	if pool.JwtPrivateKey == "" || pool.JwtPublicKey == "" {
		t.Fatal("created pool carries empty key material")
	}
	if _, err := vsjwt.DecodePrivateKeyFromPEM(pool.JwtPrivateKey); err != nil {
		t.Errorf("private key not decodable: %v", err)
	}
	if _, err := vsjwt.DecodePublicKeyFromPEM(pool.JwtPublicKey); err != nil {
		t.Errorf("public key not decodable: %v", err)
	}
}
