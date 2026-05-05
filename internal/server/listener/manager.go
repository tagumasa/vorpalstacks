package listener

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/server/fqdnrouter"
)

// ListenerTimeouts allows per-listener timeout overrides.
// When nil, the Manager applies sensible defaults (15s read, 30s write, 120s idle).
// Set individual fields to zero to disable that timeout (e.g. for WebSocket).
type ListenerTimeouts struct {
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

// ListenerConfig defines the configuration for a secondary HTTP listener.
type ListenerConfig struct {
	Name        string
	PortKey     string
	DefaultPort int
	Handler     http.Handler
	Timeouts    *ListenerTimeouts
	HostSuffix  string
}

// registeredHandler records a single handler bound to a listener port.
type registeredHandler struct {
	name       string
	hostSuffix string
	handler    http.Handler
}

// managedListener represents a bound HTTP server on a specific port.
// When multiple services share the same port (FQDN mode), fqdnRouter
// is non-nil and dispatches based on Host header. Single-handler ports
// use the handler directly, bypassing fqdnrouter entirely.
type managedListener struct {
	name       string
	portKey    string
	port       int
	server     *http.Server
	fqdnRouter *fqdnrouter.Router
	handlers   []registeredHandler
}

// Manager manages secondary HTTP listeners for auxiliary services.
// It supports port sharing: when multiple ListenerConfigs target the
// same port, their handlers are merged via fqdnrouter and dispatched
// based on the Host header suffix.
type Manager struct {
	mainPort  int
	listeners map[string]*managedListener
	portIndex map[int]*managedListener
	mu        sync.Mutex
}

// NewManager creates a new Manager for the given main port.
func NewManager(mainPort int) *Manager {
	return &Manager{
		mainPort:  mainPort,
		listeners: make(map[string]*managedListener),
		portIndex: make(map[int]*managedListener),
	}
}

func defaultTimeouts() *ListenerTimeouts {
	return &ListenerTimeouts{
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
}

// wrapWithMiddleware applies the standard middleware chain (Recoverer, RequestID)
// around the given handler.
func wrapWithMiddleware(h http.Handler) http.Handler {
	return middleware.Recoverer(middleware.RequestID(h))
}

// Register registers a secondary HTTP listener with the given configuration.
// If another listener already uses the same port, the handlers are merged
// via fqdnrouter (FQDN mode). The fqdnrouter is only created when a second
// handler registers on the same port; single-handler ports bypass it entirely.
func (m *Manager) Register(cfg ListenerConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()

	port := cfg.DefaultPort
	if port == 0 || port == m.mainPort {
		return
	}

	if _, exists := m.listeners[cfg.Name]; exists {
		return
	}

	rh := registeredHandler{
		name:       cfg.Name,
		hostSuffix: cfg.HostSuffix,
		handler:    cfg.Handler,
	}

	if existing, ok := m.portIndex[port]; ok {
		m.mergeIntoExisting(existing, rh)
		return
	}

	timeouts := cfg.Timeouts
	if timeouts == nil {
		timeouts = defaultTimeouts()
	}

	ml := &managedListener{
		name:     cfg.Name,
		portKey:  cfg.PortKey,
		port:     port,
		handlers: []registeredHandler{rh},
		server: &http.Server{
			Addr:              fmt.Sprintf(":%d", port),
			Handler:           wrapWithMiddleware(cfg.Handler),
			ReadHeaderTimeout: timeouts.ReadHeaderTimeout,
			ReadTimeout:       timeouts.ReadTimeout,
			WriteTimeout:      timeouts.WriteTimeout,
			IdleTimeout:       timeouts.IdleTimeout,
		},
	}
	m.listeners[cfg.Name] = ml
	m.portIndex[port] = ml
}

// mergeIntoExisting adds a handler to an existing listener that already
// has a port bound. When this is called for the first time on a port,
// the fqdnrouter is created and the existing handler is migrated.
func (m *Manager) mergeIntoExisting(ml *managedListener, rh registeredHandler) {
	if ml.fqdnRouter == nil {
		fr := fqdnrouter.New()
		for _, h := range ml.handlers {
			if h.hostSuffix != "" {
				fr.Register(h.hostSuffix, h.handler)
			} else {
				fr.SetFallback(h.handler)
			}
		}
		ml.fqdnRouter = fr
		ml.server.Handler = wrapWithMiddleware(fr)
	}

	if rh.hostSuffix != "" {
		ml.fqdnRouter.Register(rh.hostSuffix, rh.handler)
	} else {
		ml.fqdnRouter.SetFallback(rh.handler)
	}

	ml.handlers = append(ml.handlers, rh)
	m.listeners[rh.name] = ml
}

// Start launches all registered secondary listeners in background goroutines.
// When multiple names share the same managedListener (FQDN mode), only one
// HTTP server is started per port.
func (m *Manager) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()

	seen := make(map[*managedListener]bool)
	for _, l := range m.listeners {
		if seen[l] {
			continue
		}
		seen[l] = true
		ll := l
		names := m.handlerNames(ll)
		logs.Info("Starting secondary listener", logs.String("name", names), logs.String("port", ll.server.Addr))
		go func() {
			if err := ll.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logs.Error("Secondary listener error", logs.String("name", names), logs.Err(err))
			}
		}()
	}
}

// Shutdown gracefully shuts down all secondary listeners.
func (m *Manager) Shutdown(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	seen := make(map[*managedListener]bool)
	for _, l := range m.listeners {
		if seen[l] {
			continue
		}
		seen[l] = true
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		if err := l.server.Shutdown(shutdownCtx); err != nil {
			logs.Error("Secondary listener shutdown error", logs.String("name", m.handlerNames(l)), logs.Err(err))
		}
		cancel()
	}
}

// IsRunning reports whether a listener with the given name is currently registered.
func (m *Manager) IsRunning(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.listeners[name]
	return ok
}

// handlerNames returns a comma-separated list of handler names sharing
// the same managedListener, for logging purposes.
func (m *Manager) handlerNames(ml *managedListener) string {
	names := ""
	for _, h := range ml.handlers {
		if names != "" {
			names += ","
		}
		names += h.name
	}
	return names
}
