package scheduler

import (
	"errors"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// TestGetDistinguishesNotFoundFromStorageFailure pins the error taxonomy of
// the record Gets: an absent resource is the not-found sentinel, while a
// storage fault is returned as-is. Collapsing the latter into the sentinel
// would report a server fault to clients as a 404 and make the engine's
// retry probe discard pending retry records on transient I/O errors.
func TestGetDistinguishesNotFoundFromStorageFailure(t *testing.T) {
	store := newCompletionTestStore(t)

	if _, err := store.GetScheduleGroup(t.Context(), "absent-group"); !errors.Is(err, ErrScheduleGroupNotFound) {
		t.Errorf("GetScheduleGroup on absent group: error = %v, want ErrScheduleGroupNotFound", err)
	}
	if _, err := store.GetSchedule(t.Context(), "default", "absent-schedule"); !errors.Is(err, ErrScheduleNotFound) {
		t.Errorf("GetSchedule on absent schedule: error = %v, want ErrScheduleNotFound", err)
	}

	// A store over a closed database: every Get is a storage fault and
	// must surface as an error other than the not-found sentinel.
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	faulted := NewSchedulerStore(st, "000000000000", "us-east-1")
	if err := st.Close(); err != nil {
		t.Fatalf("close storage: %v", err)
	}

	_, groupErr := faulted.GetScheduleGroup(t.Context(), "any-group")
	if groupErr == nil {
		t.Error("GetScheduleGroup on a faulted store returned nil, want a storage error")
	} else if errors.Is(groupErr, ErrScheduleGroupNotFound) {
		t.Error("GetScheduleGroup collapsed a storage failure into ErrScheduleGroupNotFound")
	}
	_, scheduleErr := faulted.GetSchedule(t.Context(), "default", "any-schedule")
	if scheduleErr == nil {
		t.Error("GetSchedule on a faulted store returned nil, want a storage error")
	} else if errors.Is(scheduleErr, ErrScheduleNotFound) {
		t.Error("GetSchedule collapsed a storage failure into ErrScheduleNotFound")
	}
}
