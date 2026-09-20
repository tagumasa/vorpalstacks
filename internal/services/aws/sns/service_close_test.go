package sns

// Pins for the delivery spawn path's closure against Close's WaitGroup
// drain: an Add racing a Wait that observed a zero counter is the misuse
// the WaitGroup spec forbids, so spawns hold the read lock while Close
// marks the closure under the write lock.

import (
	"sync"
	"sync/atomic"
	"testing"
)

// TestRunTrackedRefusesSpawnsAfterClose pins the closed state: a spawn
// before Close runs and is drained by Close's Wait; a spawn after Close is
// refused (dropped with a warning) instead of joining the drained counter.
func TestRunTrackedRefusesSpawnsAfterClose(t *testing.T) {
	svc, _ := newFanoutTestService(t)

	var ran atomic.Int32
	svc.runTracked("pre-close", func() { ran.Add(1) })
	svc.Close()

	svc.runTracked("post-close", func() { ran.Add(1) })
	if got := ran.Load(); got != 1 {
		t.Fatalf("executions = %d, want 1 (pre-close ran under Close's drain, post-close refused)", got)
	}
}

// TestRunTrackedConcurrentWithClose pins the interleaving under the race
// detector: goroutines spawn deliveries while Close drains — every spawn
// either joins the counter before the closure or is refused after it, so
// no Add races the Wait and Close returns with nothing still scheduled.
func TestRunTrackedConcurrentWithClose(t *testing.T) {
	for round := 0; round < 20; round++ {
		svc, _ := newFanoutTestService(t)

		var wg sync.WaitGroup
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					svc.runTracked("race pin", func() {})
				}
			}()
		}
		svc.Close()
		wg.Wait()
	}
}
