package lambda

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// atomicLayerStore builds a LayerStore over a fresh empty storage.
func atomicLayerStore(t *testing.T) *LayerStore {
	t.Helper()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("new region storage manager: %v", err)
	}
	st, err := mgr.GetStorage("us-east-1")
	if err != nil {
		t.Fatalf("get storage: %v", err)
	}
	return NewLayerStore(st, "000000000000", "us-east-1")
}

// TestLayerStoreConcurrentPublishKeepsEveryVersion pins the atomic layer
// publish: concurrent PublishLayerVersion requests on one layer must all
// persist their own version — the former snapshot-taking publish let the
// second writer overwrite the first's whole layer record, silently
// dropping versions.
func TestLayerStoreConcurrentPublishKeepsEveryVersion(t *testing.T) {
	s := atomicLayerStore(t)
	if _, err := s.Create(&Layer{LayerName: "shared"}); err != nil {
		t.Fatalf("create layer: %v", err)
	}

	const publishers = 16
	var wg sync.WaitGroup
	for i := 0; i < publishers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			desc := fmt.Sprintf("published-by-%d", i)
			if _, err := s.PublishVersionAtomically("shared", func(*Layer) (*LayerVersion, error) {
				return &LayerVersion{Description: desc}, nil
			}); err != nil {
				t.Errorf("publish %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	layer, err := s.Get("shared")
	if err != nil {
		t.Fatalf("reload layer: %v", err)
	}
	if len(layer.Versions) != publishers {
		t.Fatalf("stored %d versions, want %d", len(layer.Versions), publishers)
	}
	seen := make(map[int64]bool)
	for _, v := range layer.Versions {
		seen[v.Version] = true
	}
	for n := int64(1); n <= publishers; n++ {
		if !seen[n] {
			t.Fatalf("version %d missing", n)
		}
	}
}

// TestLayerStorePublishVersusPolicyUpdateKeepsBoth pins the same
// no-lost-update rule across operations: version publishes running
// concurrently with policy-statement additions must leave both effects in
// the stored record.
func TestLayerStorePublishVersusPolicyUpdateKeepsBoth(t *testing.T) {
	s := atomicLayerStore(t)
	if _, err := s.Create(&Layer{LayerName: "mixed"}); err != nil {
		t.Fatalf("create layer: %v", err)
	}
	if _, err := s.PublishVersionAtomically("mixed", func(*Layer) (*LayerVersion, error) {
		return &LayerVersion{}, nil
	}); err != nil {
		t.Fatalf("seed version 1: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 8; i++ {
			if _, err := s.PublishVersionAtomically("mixed", func(*Layer) (*LayerVersion, error) {
				return &LayerVersion{}, nil
			}); err != nil {
				t.Errorf("publish: %v", err)
			}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 8; i++ {
			if err := s.AddPolicy("mixed", 1, &LayerPolicy{
				Id:        fmt.Sprintf("stmt-%d", i),
				Action:    "lambda:GetLayerVersion",
				Principal: "*",
			}); err != nil {
				t.Errorf("add policy: %v", err)
			}
		}
	}()
	wg.Wait()

	layer, err := s.Get("mixed")
	if err != nil {
		t.Fatalf("reload layer: %v", err)
	}
	if len(layer.Versions) != 9 {
		t.Fatalf("stored %d versions, want 9", len(layer.Versions))
	}
	v, err := s.GetVersion("mixed", 1)
	if err != nil {
		t.Fatalf("get version 1: %v", err)
	}
	if len(v.Policies) != 8 {
		t.Fatalf("version 1 carries %d policies, want 8", len(v.Policies))
	}
}

// TestFunctionStoreConcurrentPublishKeepsEveryVersion pins the atomic
// publish on the function store: the caller-loaded snapshot the former
// PublishVersion mutated could overwrite a concurrent publish or
// configuration update with stale state.
func TestFunctionStoreConcurrentPublishKeepsEveryVersion(t *testing.T) {
	s := provisionedConcurrencyStore(t)

	const publishers = 16
	var wg sync.WaitGroup
	for i := 0; i < publishers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := s.PublishVersionAtomically("fn", "", ""); err != nil {
				t.Errorf("publish: %v", err)
			}
		}()
	}
	wg.Wait()

	function, err := s.Get("fn")
	if err != nil {
		t.Fatalf("reload function: %v", err)
	}
	if len(function.Versions) != publishers {
		t.Fatalf("stored %d versions, want %d", len(function.Versions), publishers)
	}
}

// TestFunctionStorePublishVersionAtomicallyChecksRevisionPrecondition
// pins the in-lock optimistic-locking precondition: a RevisionId captured
// before a concurrent revision bump must fail the publish instead of
// silently creating the version.
func TestFunctionStorePublishVersionAtomicallyChecksRevisionPrecondition(t *testing.T) {
	s := provisionedConcurrencyStore(t)

	stale, err := s.Get("fn")
	if err != nil {
		t.Fatalf("load function: %v", err)
	}
	if _, err := s.UpdateAtomically("fn", func(fn *Function) error {
		fn.Description = "bumped elsewhere"
		return nil
	}); err != nil {
		t.Fatalf("bump revision: %v", err)
	}

	if _, err := s.PublishVersionAtomically("fn", "", stale.RevisionId); !errors.Is(err, ErrRevisionMismatch) {
		t.Fatalf("stale revision publish: expected ErrRevisionMismatch, got %v", err)
	}
	function, err := s.Get("fn")
	if err != nil {
		t.Fatalf("reload function: %v", err)
	}
	if len(function.Versions) != 0 {
		t.Fatalf("rejected publish left %d versions, want 0", len(function.Versions))
	}
}
