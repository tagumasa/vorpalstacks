package resilience

import (
	"sync"
	"testing"
	"time"
)

// The containment defer pattern every call site uses: recover() called
// directly by the deferred closure. This pin drives a real panic through
// the pattern and asserts the handler runs and the goroutine survives.
func TestDirectRecoverPatternContains(t *testing.T) {
	ran := make(chan struct{}, 1)
	func() {
		defer func() {
			if r := recover(); r != nil {
				LogPanic("pin direct recover", r)
				ran <- struct{}{}
			}
		}()
		panic("pin: boom")
	}()
	select {
	case <-ran:
	default:
		t.Fatal("recover handler did not run")
	}
}

// RestartAfterPanic must relaunch the loop and leave the WaitGroup
// balanced: the Add inside it replaces the membership the relaunched
// loop consumes with its own deferred Done, and the panicking frame's
// own Done still runs afterwards.
func TestRestartAfterPanicRelaunches(t *testing.T) {
	var wg sync.WaitGroup
	restarted := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				RestartAfterPanic("pin restart", r, &wg, func() {
					defer wg.Done()
					close(restarted)
				})
			}
		}()
		panic("pin: restart me")
	}()
	select {
	case <-restarted:
	case <-time.After(5 * time.Second):
		t.Fatal("restart function was never relaunched")
	}
	settled := make(chan struct{})
	go func() {
		wg.Wait()
		close(settled)
	}()
	select {
	case <-settled:
	case <-time.After(5 * time.Second):
		t.Fatal("WaitGroup never settled after restart")
	}
}
