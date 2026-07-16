// Package hooks provides hook functionality for extending component behavior.
package hooks

import (
	"context"
)

// Hook represents a function that can be registered and executed in response to events.
// It takes a context and returns an error if the hook execution fails.
type Hook func(ctx context.Context) error

// HookOption configures a hook when it is registered.
type HookOption func(*hookConfig)

type hookConfig struct {
	priority int
	name     string
}

// WithPriority sets the priority for a hook.
func WithPriority(priority int) HookOption {
	return func(c *hookConfig) {
		c.priority = priority
	}
}

// WithName sets the name for a hook.
func WithName(name string) HookOption {
	return func(c *hookConfig) {
		c.name = name
	}
}

type hookEntry struct {
	hook     Hook
	priority int
	name     string
}

// Register registers a hook for the given event.
func Register(event string, h Hook, opts ...HookOption) {
	globalManager.Register(event, h, opts...)
}

// Run runs all hooks registered for the given event.
func Run(ctx context.Context, event string) error {
	return globalManager.Run(ctx, event)
}

// Clear clears all registered hooks.
func Clear() {
	globalManager.Clear()
}

// On registers a hook for the given event (alias for Register).
func On(event string, h Hook, opts ...HookOption) {
	Register(event, h, opts...)
}
