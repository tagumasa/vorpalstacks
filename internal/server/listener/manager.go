package listener

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"vorpalstacks/internal/common/headerorder"
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
	// TLS, when set, binds the listener as a TLS server whose certificate
	// is resolved per handshake through GetCertificate (SNI). A TLS
	// listener keeps the TLS configuration of its first registrant: the
	// handshake terminates before Host-based routing, so merged handlers
	// share one certificate resolver.
	TLS *ListenerTLSConfig
}

// ListenerTLSConfig carries the SNI certificate resolver for a TLS listener.
// GetCertificate mirrors tls.Config.GetCertificate; returning an error (or
// nil with no certificate) fails the handshake.
type ListenerTLSConfig struct {
	GetCertificate func(*tls.ClientHelloInfo) (*tls.Certificate, error)
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
	dispatch   *stableHandler
	fqdnRouter *fqdnrouter.Router
	handlers   []registeredHandler
}

// stableHandler wraps an http.Handler behind an atomic pointer so that
// the active handler can be swapped without racing the HTTP server's
// accept goroutine (which reads server.Handler on every request).
type stableHandler struct {
	current atomic.Pointer[http.Handler]
}

func (s *stableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := s.current.Load()
	if h == nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}
	(*h).ServeHTTP(w, r)
}

func (s *stableHandler) Store(h http.Handler) {
	s.current.Store(&h)
}

// Manager manages secondary HTTP listeners for auxiliary services.
// It supports port sharing: when multiple ListenerConfigs target the
// same port, their handlers are merged via fqdnrouter and dispatched
// based on the Host header suffix.
type Manager struct {
	mainPort  int
	started   bool
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

	sh := &stableHandler{}
	sh.Store(wrapWithMiddleware(cfg.Handler))
	ml := &managedListener{
		name:     cfg.Name,
		portKey:  cfg.PortKey,
		port:     port,
		handlers: []registeredHandler{rh},
		dispatch: sh,
		server: &http.Server{
			Addr:              fmt.Sprintf(":%d", port),
			Handler:           headerorder.TLSStateMiddleware(sh),
			ConnContext:       headerorder.ConnContext(nil),
			ReadHeaderTimeout: timeouts.ReadHeaderTimeout,
			ReadTimeout:       timeouts.ReadTimeout,
			WriteTimeout:      timeouts.WriteTimeout,
			IdleTimeout:       timeouts.IdleTimeout,
		},
	}
	if cfg.TLS != nil && cfg.TLS.GetCertificate != nil {
		// The TLS floor matches the platform default; per-distribution
		// protocol policies cannot lower the crypto/tls handshake floor
		// on a shared listener.
		ml.server.TLSConfig = &tls.Config{
			GetCertificate: cfg.TLS.GetCertificate,
			MinVersion:     tls.VersionTLS12,
		}
	}
	m.listeners[cfg.Name] = ml
	m.portIndex[port] = ml

	if m.started {
		logs.Info("Starting dynamic listener", logs.String("name", cfg.Name), logs.String("port", ml.server.Addr))
		go func() {
			if err := serveListener(ml); err != nil && err != http.ErrServerClosed {
				logs.Error("Dynamic listener error", logs.String("name", cfg.Name), logs.Err(err))
			}
		}()
	}
}

// serveListener runs the listener's HTTP server, terminating TLS when the
// listener carries a TLS configuration. Certificate files stay empty: the
// TLS configuration resolves every certificate through GetCertificate.
// The header-order capture wraps the listener (above the TLS layer) so
// order-sensitive consumers such as the WAF HeaderOrder component see the
// wire order on the per-service planes as well; the eager handshake and
// Request.TLS restoration in the headerorder ConnContext/TLSStateMiddleware
// pair compensate for the wrapper hiding the *tls.Conn from the HTTP
// server. Managed TLS listeners negotiate no ALPN protocol today; adding
// HTTP/2 negotiation would require revisiting this wrapper, because the
// server's HTTP/2 dispatch also relies on the concrete *tls.Conn type.
func serveListener(ml *managedListener) error {
	raw, err := net.Listen("tcp", ml.server.Addr)
	if err != nil {
		return err
	}
	capture := headerorder.NewListener(raw)
	if ml.server.TLSConfig != nil {
		capture = headerorder.NewListener(tls.NewListener(raw, ml.server.TLSConfig))
	}
	return ml.server.Serve(capture)
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
		ml.dispatch.Store(wrapWithMiddleware(fr))
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

	m.started = true

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
			if err := serveListener(ll); err != nil && err != http.ErrServerClosed {
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

// Unregister removes a listener by name and shuts down its HTTP server if no
// other handlers share the same port. Used for dynamic per-cluster listeners.
// When FQDN-mode port sharing is active, the fqdnrouter is rebuilt from the
// remaining handlers so requests for the removed handler's host suffix no
// longer reach the old handler.
func (m *Manager) Unregister(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ml, ok := m.listeners[name]
	if !ok {
		return
	}

	delete(m.listeners, name)
	for i, h := range ml.handlers {
		if h.name == name {
			ml.handlers = append(ml.handlers[:i], ml.handlers[i+1:]...)
			break
		}
	}

	if len(ml.handlers) == 0 {
		delete(m.portIndex, ml.port)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ml.server.Shutdown(ctx)
		return
	}

	// Rebuild the fqdnrouter from remaining handlers so the unregistered
	// handler's host suffix no longer routes to the old handler.
	m.rebuildFQDNRouter(ml)
}

// rebuildFQDNRouter recreates the fqdnrouter from the current handler list.
// Called after a handler is removed from a shared (FQDN-mode) port.
func (m *Manager) rebuildFQDNRouter(ml *managedListener) {
	if len(ml.handlers) == 0 {
		return
	}
	if len(ml.handlers) == 1 {
		// Single handler remaining — bypass fqdnrouter entirely.
		ml.fqdnRouter = nil
		ml.dispatch.Store(wrapWithMiddleware(ml.handlers[0].handler))
		return
	}
	fr := fqdnrouter.New()
	for _, h := range ml.handlers {
		if h.hostSuffix != "" {
			fr.Register(h.hostSuffix, h.handler)
		} else {
			fr.SetFallback(h.handler)
		}
	}
	ml.fqdnRouter = fr
	ml.dispatch.Store(wrapWithMiddleware(fr))
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
