package listener

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"vorpalstacks/internal/core/logs"
)

type ListenerConfig struct {
	Name        string
	PortKey     string
	DefaultPort int
	Handler     http.Handler
}

type managedListener struct {
	name    string
	portKey string
	server  *http.Server
	handler http.Handler
}

type Manager struct {
	mainPort      int
	listeners     map[string]*managedListener
	mu            sync.Mutex
	shutdownHooks []func(context.Context) error
}

func NewManager(mainPort int) *Manager {
	return &Manager{
		mainPort:  mainPort,
		listeners: make(map[string]*managedListener),
	}
}

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

	m.listeners[cfg.Name] = &managedListener{
		name:    cfg.Name,
		portKey: cfg.PortKey,
		server: &http.Server{
			Addr:              fmt.Sprintf(":%d", port),
			Handler:           r,
			ReadHeaderTimeout: 5 * 1e9,
			ReadTimeout:       15 * 1e9,
			WriteTimeout:      30 * 1e9,
			IdleTimeout:       120 * 1e9,
		},
		handler: cfg.Handler,
	}
}

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

func (m *Manager) IsRunning(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.listeners[name]
	return ok
}
