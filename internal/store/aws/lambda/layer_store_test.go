package lambda

import (
	"errors"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// layerStore builds a LayerStore over fresh region storage.
func layerStore(t *testing.T) *LayerStore {
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

// TestLayerCreateDuplicate pins the duplicate sentinel of Create: creating
// a layer whose name already exists answers ErrLayerAlreadyExists, the
// sentinel the service's concurrent-first-publish recovery matches. A
// storage fault or a different sentinel would silently turn every race
// loser into a conflict failure.
func TestLayerCreateDuplicate(t *testing.T) {
	s := layerStore(t)

	created, err := s.Create(&Layer{LayerName: "race-layer"})
	if err != nil {
		t.Fatalf("first create: %v", err)
	}
	if created.LayerArn == "" {
		t.Fatal("create did not assign the layer ARN")
	}

	_, err = s.Create(&Layer{LayerName: "race-layer"})
	if !errors.Is(err, ErrLayerAlreadyExists) {
		t.Fatalf("second create error %v, want ErrLayerAlreadyExists", err)
	}
}

// TestLayerMutationsPropagateReadFaults pins the error discrimination of
// the in-lock layer re-reads: a storage fault (here a record that fails
// to unmarshal) must propagate as itself, never collapsed into the
// layer-not-found sentinel a client would read as a 404.
func TestLayerMutationsPropagateReadFaults(t *testing.T) {
	s := layerStore(t)
	if err := s.BaseStore.PutRaw("fault-layer", []byte("not json")); err != nil {
		t.Fatalf("seed corrupt record: %v", err)
	}

	ops := map[string]func() error{
		"DeleteVersion": func() error {
			return s.DeleteVersion("fault-layer", 1)
		},
		"AddPolicy": func() error {
			return s.AddPolicy("fault-layer", 1, &LayerPolicy{Id: "s1"})
		},
		"RemovePolicy": func() error {
			return s.RemovePolicy("fault-layer", 1, "s1")
		},
	}
	for name, op := range ops {
		err := op()
		if err == nil {
			t.Fatalf("%s: corrupt record produced no error", name)
		}
		if errors.Is(err, ErrLayerNotFound) {
			t.Fatalf("%s: storage fault collapsed into ErrLayerNotFound: %v", name, err)
		}
	}
}
