package resilience

import (
	"log/slog"
	"runtime/debug"
	"sync"
)

// RecoverPanic recovers from a panic, logs it with a full stack trace, and
// returns the recovered value. Returns nil if no panic occurred. Use in defer:
//
//	defer func() { if resilience.RecoverPanic("component") { /* handle */ } }()
func RecoverPanic(component string) (recovered any) {
	if r := recover(); r != nil {
		slog.Error("PANIC in "+component, "panic", r, "stack", string(debug.Stack()))
		recovered = r
	}
	return
}

// RecoverAndRestart recovers from a panic, logs it, calls wg.Add(1), and
// launches restart as a new goroutine. Use in defer for long-running loops:
//
//	defer func() { resilience.RecoverAndRestart("evalLoop", &e.wg, e.evalLoop) }()
func RecoverAndRestart(component string, wg *sync.WaitGroup, restart func()) {
	if r := recover(); r != nil {
		slog.Error("PANIC in "+component+", restarting", "panic", r, "stack", string(debug.Stack()))
		wg.Add(1)
		go restart()
	}
}
