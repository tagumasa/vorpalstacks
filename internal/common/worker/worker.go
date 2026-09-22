// Package worker carries the background ticker-worker shell the
// scheduled engines share: the TEST_MODE cadence resolution, the
// supervised ticker loop, and the configured-region iteration over the
// storage manager. What one pass does — the due computation, the
// boundary claim, the execution record — is the owning engine's
// strategy and stays with it.
package worker

import (
	"context"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
)

// Cadence resolves a ticker worker's interval: the production cadence,
// or the test cadence when the process runs under TEST_MODE, so
// integration tests observe a full evaluation pass in seconds instead
// of waiting out the production interval.
func Cadence(production, test time.Duration) time.Duration {
	if os.Getenv("TEST_MODE") == "true" {
		return test
	}
	return production
}

// TickerLoop returns one supervised worker's loop body: tick fires once
// per interval until the done channel closes. The returned closure is
// the body RunSupervised (or an equivalent lifecycle) runs; done is the
// caller's termination signal — a service context for service-owned
// workers, a stop channel for engine-owned ones.
func TickerLoop(done <-chan struct{}, interval time.Duration, tick func()) func() {
	return func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				tick()
			}
		}
	}
}

// respawnBackoff bounds the supervisor's restart rate after a panic: a
// body that panics before its first wait (an engine's immediate first
// pass) would otherwise spin the supervisor at CPU speed.
const respawnBackoff = 1 * time.Second

// RunSupervised runs fn on its own goroutine and restarts it when it
// panics. A normal return ends the worker; a panic is logged with the
// worker's name and stack, and the body restarts after the respawn
// backoff. ctx cancellation ends the supervision; wg settles when the
// goroutine exits for either reason.
func RunSupervised(ctx context.Context, wg *sync.WaitGroup, name string, fn func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			panicked := true
			func() {
				defer func() {
					if r := recover(); r != nil {
						logs.Error("PANIC in "+name+", restarting",
							logs.Any("panic", r),
							logs.String("stack", string(debug.Stack())))
					}
				}()
				fn()
				panicked = false
			}()
			if !panicked || ctx.Err() != nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(respawnBackoff):
			}
		}
	}()
}

// ActiveRegions returns the regions background work must cover: every
// region the storage manager holds storage for, nil-safe for a service
// without one. Workers iterate this set — never a per-service store
// cache, which fills only when API traffic constructs a region's store
// — so background enforcement in a region does not depend on unrelated
// foreground traffic having touched it: an ENABLED scheduled query in a
// quiet region still fires, its schedule boundaries are not skipped
// silently, and retention still purges. The empty set falls back to
// the platform default region inside the manager itself.
func ActiveRegions(sm *storage.RegionStorageManager) []string {
	if sm == nil {
		return nil
	}
	return sm.GetActiveRegions()
}
