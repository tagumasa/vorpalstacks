package listener

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"vorpalstacks/internal/core/logs"
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
}

type managedListener struct {
	name    string
	portKey string
	server  *http.Server
	handler http.Handler
}

// Manager manages secondary HTTP listeners for auxiliary services.
type Manager struct {
	mainPort      int
	listeners     map[string]*managedListener
	mu            sync.Mutex
	shutdownHooks []func(context.Context) error
}

// NewManager creates a new Manager for the given main port.
func NewManager(mainPort int) *Manager {
	return &Manager{
		mainPort:  mainPort,
		listeners: make(map[string]*managedListener),
	}
}

// Register registers a secondary HTTP listener with the given configuration.
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

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Handle("/*", cfg.Handler)

	defaultTimeouts := &ListenerTimeouts{
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	timeouts := cfg.Timeouts
	if timeouts == nil {
		timeouts = defaultTimeouts
	}

	m.listeners[cfg.Name] = &managedListener{
		name:    cfg.Name,
		portKey: cfg.PortKey,
		server: &http.Server{
			Addr:              fmt.Sprintf(":%d", port),
			Handler:           r,
			ReadHeaderTimeout: timeouts.ReadHeaderTimeout,
			ReadTimeout:       timeouts.ReadTimeout,
			WriteTimeout:      timeouts.WriteTimeout,
			IdleTimeout:       timeouts.IdleTimeout,
		},
		handler: cfg.Handler,
	}
}

// Start launches all registered secondary listeners in background goroutines.
func (m *Manager) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, l := range m.listeners {
		ll := l
		logs.Info("Starting secondary listener", logs.String("name", name), logs.String("port", ll.server.Addr))
		go func() {
			if err := ll.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logs.Error("Secondary listener error", logs.String("name", name), logs.Err(err))
			}
		}()
	}
}

// Shutdown gracefully shuts down all secondary listeners and executes registered shutdown hooks.
func (m *Manager) Shutdown(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, l := range m.listeners {
		if err := l.server.Shutdown(ctx); err != nil {
			logs.Error("Secondary listener shutdown error", logs.String("name", name), logs.Err(err))
		}
	}

	for i := len(m.shutdownHooks) - 1; i >= 0; i-- {
		if err := m.shutdownHooks[i](ctx); err != nil {
			logs.Error("Listener shutdown hook error", logs.Err(err))
		}
	}
}

// IsRunning reports whether a listener with the given name is currently registered.
func (m *Manager) IsRunning(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.listeners[name]
	return ok
}
