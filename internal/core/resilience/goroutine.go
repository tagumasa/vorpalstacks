package resilience

import (
	"log/slog"
	"runtime/debug"
	"sync"
)

// LogPanic logs a recovered panic value with a full stack trace. It is the
// logging half of a containment defer; the recover() call itself must live
// directly in the deferred closure, because Go only honours a recover()
// made by the deferred function itself — a helper calling recover() on the
// caller's behalf returns nil and the panic kills the process:
//
//	defer func() {
//		if r := recover(); r != nil {
//			resilience.LogPanic("component", r)
//		}
//	}()
func LogPanic(component string, r any) {
	slog.Error("PANIC in "+component, "panic", r, "stack", string(debug.Stack()))
}

// RestartAfterPanic logs a recovered panic value and relaunches restart as
// a new goroutine, keeping the caller's WaitGroup counter balanced: the
// site's own deferred Done still runs after this returns, so the Add here
// replaces the membership the restarted goroutine needs. As with LogPanic,
// the recover() call must live directly in the deferred closure:
//
//	defer func() {
//		if r := recover(); r != nil {
//			resilience.RestartAfterPanic("loop", r, &e.wg, e.run)
//		}
//	}()
func RestartAfterPanic(component string, r any, wg *sync.WaitGroup, restart func()) {
	slog.Error("PANIC in "+component+", restarting", "panic", r, "stack", string(debug.Stack()))
	wg.Add(1)
	go restart()
}
