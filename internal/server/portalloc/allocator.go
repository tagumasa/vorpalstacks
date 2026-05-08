// Package portalloc provides dynamic port allocation for services
// operating in Individual mode, and bind address resolution for console ports.
// All state is persisted in Pebble via PortStore — no in-memory cache.
package portalloc

import (
	"fmt"
	"net"
	"sync"

	"vorpalstacks/internal/common/serviceports"
)

// PortStore is the interface used by Allocator to persist port allocations.
// Implemented by storeconfig.Store.
type PortStore interface {
	GetResourcePort(serviceKey, resourceID string) (int, error)
	SetResourcePort(serviceKey, resourceID string, port int) error
	DeleteResourcePort(serviceKey, resourceID string) error
	ListResourcePorts(serviceKey string) (map[string]int, error)
}

// Allocator manages dynamic port allocation within [start, end].
// Pebble is the single source of truth. It is safe for concurrent use.
type Allocator struct {
	store PortStore
	start int
	end   int
	mu    sync.Mutex
}

// New creates a new Allocator backed by the given PortStore.
func New(store PortStore) *Allocator {
	return NewWithRange(store, serviceports.DynamicRangeStart, serviceports.DynamicRangeEnd)
}

// NewWithRange creates a new Allocator with a custom port range.
func NewWithRange(store PortStore, start, end int) *Allocator {
	return &Allocator{
		store: store,
		start: start,
		end:   end,
	}
}

// Get returns the allocated port for a service resource. If no port is allocated
// yet, it picks the first free port in the range, persists it, and returns it.
// The same (serviceKey, resourceID) always returns the same port across restarts.
func (a *Allocator) Get(serviceKey, resourceID string) (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	port, err := a.store.GetResourcePort(serviceKey, resourceID)
	if err == nil && port > 0 {
		return port, nil
	}

	port, err = a.findFreePort(serviceKey)
	if err != nil {
		return 0, err
	}

	if err := a.store.SetResourcePort(serviceKey, resourceID, port); err != nil {
		return 0, fmt.Errorf("portalloc: persist allocation: %w", err)
	}

	return port, nil
}

// Release removes the port allocation for a service resource.
func (a *Allocator) Release(serviceKey, resourceID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.store.DeleteResourcePort(serviceKey, resourceID)
}

// findFreePort scans the range [start, end] and returns the first port not
// allocated to any resource under the same serviceKey. Caller must hold a.mu.
func (a *Allocator) findFreePort(serviceKey string) (int, error) {
	used := make(map[int]bool)
	existing, err := a.store.ListResourcePorts(serviceKey)
	if err == nil {
		for _, port := range existing {
			if port >= a.start && port <= a.end {
				used[port] = true
			}
		}
	}

	for p := a.start; p <= a.end; p++ {
		if !used[p] {
			return p, nil
		}
	}

	return 0, fmt.Errorf("portalloc: no free ports in range [%d, %d]", a.start, a.end)
}

// BindAddrStore is the interface for reading bind address configuration.
type BindAddrStore interface {
	GetString(key string) string
}

// ResolveBindAddr resolves the console bind address from the config store.
// Default is 0.0.0.0 (bind to all interfaces). Fallback on any error: 127.0.0.1.
func ResolveBindAddr(store BindAddrStore) string {
	if store == nil {
		return "127.0.0.1"
	}
	mode := store.GetString("server.bind_mode")
	if mode == "" {
		mode = "all"
	}
	switch mode {
	case "localhost":
		return "127.0.0.1"
	case "all":
		return "0.0.0.0"
	case "interface":
		iface := store.GetString("server.bind_interface")
		if iface == "" {
			return "127.0.0.1"
		}
		ip := net.ParseIP(iface)
		if ip == nil {
			return "127.0.0.1"
		}
		ifaces, err := net.Interfaces()
		if err != nil {
			return "127.0.0.1"
		}
		for _, i := range ifaces {
			addrs, err := i.Addrs()
			if err != nil {
				continue
			}
			for _, a := range addrs {
				var ifaceIP net.IP
				switch v := a.(type) {
				case *net.IPNet:
					ifaceIP = v.IP
				case *net.IPAddr:
					ifaceIP = v.IP
				}
				if ifaceIP != nil && ifaceIP.Equal(ip) {
					return iface
				}
			}
		}
		return "127.0.0.1"
	default:
		return "127.0.0.1"
	}
}
