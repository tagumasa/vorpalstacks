// Package hooks provides hook functionality for extending component behavior.
package hooks

import (
	"context"
	"sort"
	"sync"
)

// Manager manages the registration and execution of hooks.
type Manager struct {
	mu    sync.RWMutex
	hooks map[string][]hookEntry
}

var globalManager = NewManager()

// NewManager creates a new hook manager.
func NewManager() *Manager {
	return &Manager{
		hooks: make(map[string][]hookEntry),
	}
}

// Register registers a hook for a given event.
func (m *Manager) Register(event string, h Hook, opts ...HookOption) {
	cfg := &hookConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	entry := hookEntry{
		hook:     h,
		priority: cfg.priority,
		name:     cfg.name,
	}

	m.hooks[event] = append(m.hooks[event], entry)
	sort.Slice(m.hooks[event], func(i, j int) bool {
		return m.hooks[event][i].priority > m.hooks[event][j].priority
	})
}

// Run executes all hooks registered for a given event synchronously.
func (m *Manager) Run(ctx context.Context, event string) error {
	m.mu.RLock()
	hooks := make([]hookEntry, len(m.hooks[event]))
	copy(hooks, m.hooks[event])
	m.mu.RUnlock()

	for _, entry := range hooks {
		if err := entry.hook(ctx); err != nil {
			return err
		}
	}
	return nil
}

// RunAsync executes all hooks registered for a given event asynchronously.
func (m *Manager) RunAsync(ctx context.Context, event string) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		errCh <- m.Run(ctx, event)
	}()
	return errCh
}

// GetHooks returns all hooks registered for a given event.
func (m *Manager) GetHooks(event string) []Hook {
	m.mu.RLock()
	defer m.mu.RUnlock()

	hooks := make([]Hook, len(m.hooks[event]))
	for i, entry := range m.hooks[event] {
		hooks[i] = entry.hook
	}
	return hooks
}

// Clear removes all registered hooks.
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hooks = make(map[string][]hookEntry)
}

// ClearEvent removes all hooks registered for a specific event.
func (m *Manager) ClearEvent(event string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.hooks, event)
}
