package iot

import (
	"testing"

	"vorpalstacks/internal/core/storage"
)

func newIotStore(t *testing.T) *IotStore {
	t.Helper()
	tmpDir := t.TempDir()
	st, err := storage.Open(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	return NewIotStore(st, "000000000000", "us-east-1", nil)
}

func TestCreateThingRejectsDeprecatedType(t *testing.T) {
	st := newIotStore(t)

	st.CreateThingType(&ThingType{ThingTypeName: "oldType", Deprecated: true})

	_, err := st.CreateThing(&Thing{ThingName: "t1", ThingTypeName: "oldType"})
	if err != ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}

	st.CreateThingType(&ThingType{ThingTypeName: "newType", Deprecated: false})
	thing, err := st.CreateThing(&Thing{ThingName: "t2", ThingTypeName: "newType"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if thing.ThingName != "t2" {
		t.Fatalf("expected t2, got %s", thing.ThingName)
	}
}

func TestUpdateThingRejectsDeprecatedType(t *testing.T) {
	st := newIotStore(t)

	st.CreateThingType(&ThingType{ThingTypeName: "typeA", Deprecated: false})
	_, err := st.CreateThing(&Thing{ThingName: "t1", ThingTypeName: "typeA"})
	if err != nil {
		t.Fatalf("create thing: %v", err)
	}

	st.SetThingTypeDeprecation("typeA", true)

	_, err = st.UpdateThing("t1", ThingUpdateOpts{ThingTypeName: "typeA"})
	if err != ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest for deprecated type, got %v", err)
	}
}

func TestUpdateJobRejectsInvalidStatus(t *testing.T) {
	st := newIotStore(t)

	st.CreateJob(&Job{JobID: "j1", Status: "QUEUED"})

	_, err := st.UpdateJob("j1", JobUpdateOpts{Status: "IN_PROGRESS"})
	if err != nil {
		t.Fatalf("unexpected error for valid status: %v", err)
	}

	_, err = st.UpdateJob("j1", JobUpdateOpts{Status: "BANANA"})
	if err != ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest for invalid status, got %v", err)
	}
}
