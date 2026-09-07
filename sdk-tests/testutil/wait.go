package testutil

import (
	"fmt"
	"time"
)

// waitFor polls cond every interval until it holds or timeout elapses.
// It replaces blind fixed sleeps: the wait ends the moment the condition
// is observed instead of after a worst-case delay. A final check after
// the deadline keeps conditions that flip during the last interval from
// being reported as timeouts.
func waitFor(interval, timeout time.Duration, cond func() bool) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return nil
		}
		time.Sleep(interval)
	}
	if cond() {
		return nil
	}
	return fmt.Errorf("condition not met within %v", timeout)
}
