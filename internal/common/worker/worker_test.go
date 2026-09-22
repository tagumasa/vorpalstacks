package worker

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

// Cadence picks the production interval outside TEST_MODE and the test
// interval under it.
func TestCadenceHonoursTestModeBoundary(t *testing.T) {
	if got := Cadence(5*time.Minute, time.Second); got != 5*time.Minute {
		t.Fatalf("production cadence = %v, want 5m", got)
	}
	t.Setenv("TEST_MODE", "true")
	if got := Cadence(5*time.Minute, time.Second); got != time.Second {
		t.Fatalf("test cadence = %v, want 1s", got)
	}
}

// TickerLoop fires tick once per interval and returns when done closes.
func TestTickerLoopTicksUntilDone(t *testing.T) {
	done := make(chan struct{})
	var ticks atomic.Int32
	loop := TickerLoop(done, 5*time.Millisecond, func() { ticks.Add(1) })

	finished := make(chan struct{})
	go func() {
		loop()
		close(finished)
	}()

	// Two intervals must elapse before the stop: at least one tick and
	// a clean return afterwards.
	time.Sleep(25 * time.Millisecond)
	close(done)
	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("loop did not return after done closed")
	}
	if ticks.Load() == 0 {
		t.Fatal("tick never fired")
	}
}

// RunSupervised restarts a panicking body until it returns normally,
// and the wait group settles on that normal return.
func TestRunSupervisedRespawnsAfterPanic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var runs atomic.Int32
	var wg sync.WaitGroup
	RunSupervised(ctx, &wg, "test worker", func() {
		if runs.Add(1) < 3 {
			panic("boom")
		}
		// The third run returns normally: the supervision ends there.
	})

	finished := make(chan struct{})
	go func() {
		wg.Wait()
		close(finished)
	}()
	select {
	case <-finished:
	case <-time.After(5 * time.Second):
		t.Fatal("supervisor did not settle after the body returned normally")
	}
	if got := runs.Load(); got != 3 {
		t.Fatalf("body ran %d times, want 3 (two panics, one clean return)", got)
	}
}

// RunSupervised stops restarting when the context is cancelled during
// the respawn backoff.
func TestRunSupervisedStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var runs atomic.Int32
	var wg sync.WaitGroup
	RunSupervised(ctx, &wg, "test worker", func() {
		runs.Add(1)
		panic("always panics")
	})
	cancel()

	finished := make(chan struct{})
	go func() {
		wg.Wait()
		close(finished)
	}()
	select {
	case <-finished:
	case <-time.After(5 * time.Second):
		t.Fatal("supervisor did not settle after context cancellation")
	}
	// The body may have run once or twice (the in-flight attempt plus
	// at most one backoff race), but the cancellation must stop the
	// restart cycle — the settle above is the assertion.
	if runs.Load() > 2 {
		t.Fatalf("body ran %d times after cancellation, restart cycle did not stop", runs.Load())
	}
}

// ActiveRegions is nil-safe and iterates the storage manager's regions.
func TestActiveRegionsNilSafeAndStorageBacked(t *testing.T) {
	if got := ActiveRegions(nil); got != nil {
		t.Fatalf("ActiveRegions(nil) = %v, want nil", got)
	}
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("region storage manager: %v", err)
	}
	t.Cleanup(func() { sm.Close() })
	if _, err := sm.GetStorage("us-west-2"); err != nil {
		t.Fatalf("get storage: %v", err)
	}
	regions := ActiveRegions(sm)
	if len(regions) != 1 || regions[0] != "us-west-2" {
		t.Fatalf("ActiveRegions = %v, want [us-west-2]", regions)
	}
}
