// Package portalloc provides dynamic port allocation for services
// operating in Individual mode. Ports are allocated from a configurable
// range and persisted via the PortStore interface.
package portalloc

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"vorpalstacks/internal/common/serviceports"
)

// PortStore is the interface used by Allocator to persist port allocations.
// Implemented by storeconfig.Store.
type PortStore interface {
	GetResourcePort(serviceKey, resourceID string) (int, error)
	SetResourcePort(serviceKey, resourceID string, port int) error
	DeleteResourcePort(serviceKey, resourceID string) error
	ForEach(func(key string, raw []byte) error) error
}

// Allocator manages dynamic port allocation within [start, end].
// It is safe for concurrent use.
type Allocator struct {
	store     PortStore
	start     int
	end       int
	allocated map[string]int // "serviceKey.resourceID" → port
	freePool  []int
	mu        sync.Mutex
}

// New creates a new Allocator backed by the given PortStore.
func New(store PortStore) *Allocator {
	return NewWithRange(store, serviceports.DynamicRangeStart, serviceports.DynamicRangeEnd)
}

// NewWithRange creates a new Allocator with a custom port range.
func NewWithRange(store PortStore, start, end int) *Allocator {
	free := make([]int, 0, end-start+1)
	for p := start; p <= end; p++ {
		free = append(free, p)
	}
	return &Allocator{
		store:     store,
		start:     start,
		end:       end,
		allocated: make(map[string]int),
		freePool:  free,
	}
}

func allocKey(serviceKey, resourceID string) string {
	return serviceKey + "." + resourceID
}

// LoadExisting reads all existing port allocations from the store and
// rebuilds the free pool. Call this once at startup before any Allocate calls.
func (a *Allocator) LoadExisting(knownPrefixes []string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.store.ForEach(func(key string, raw []byte) error {
		if !a.isResourcePortKey(key, knownPrefixes) {
			return nil
		}
		var entry struct {
			Value interface{} `json:"value"`
		}
		if err := json.Unmarshal(raw, &entry); err != nil {
			return nil
		}
		port := 0
		switch v := entry.Value.(type) {
		case float64:
			port = int(v)
		case int:
			port = v
		default:
			return nil
		}
		if port < a.start || port > a.end {
			return nil
		}
		a.allocated[key] = port
		a.removeFromFreePool(port)
		return nil
	})
}

func (a *Allocator) isResourcePortKey(key string, prefixes []string) bool {
	for _, pfx := range prefixes {
		if len(key) > len(pfx) && key[:len(pfx)+1] == pfx+"." {
			return true
		}
	}
	return false
}

func (a *Allocator) removeFromFreePool(port int) {
	for i, p := range a.freePool {
		if p == port {
			a.freePool = append(a.freePool[:i], a.freePool[i+1:]...)
			return
		}
	}
}

// Allocate assigns a free port to the given service resource, persists it,
// and returns the port number.
func (a *Allocator) Allocate(serviceKey, resourceID string) (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	k := allocKey(serviceKey, resourceID)
	if p, ok := a.allocated[k]; ok {
		return p, nil
	}

	if len(a.freePool) == 0 {
		return 0, fmt.Errorf("portalloc: no free ports in range [%d, %d]", a.start, a.end)
	}

	port := a.freePool[0]
	a.freePool = a.freePool[1:]
	a.allocated[k] = port

	if err := a.store.SetResourcePort(serviceKey, resourceID, port); err != nil {
		delete(a.allocated, k)
		a.freePool = append(a.freePool, port)
		sort.Ints(a.freePool)
		return 0, fmt.Errorf("portalloc: persist allocation: %w", err)
	}

	return port, nil
}

// Release returns a port to the free pool and removes the allocation.
func (a *Allocator) Release(serviceKey, resourceID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	k := allocKey(serviceKey, resourceID)
	port, ok := a.allocated[k]
	if !ok {
		return nil
	}

	if err := a.store.DeleteResourcePort(serviceKey, resourceID); err != nil {
		return fmt.Errorf("portalloc: delete allocation: %w", err)
	}

	delete(a.allocated, k)
	a.freePool = append(a.freePool, port)
	sort.Ints(a.freePool)
	return nil
}

// Get returns the allocated port for a service resource, or auto-allocates one.
func (a *Allocator) Get(serviceKey, resourceID string) (int, error) {
	a.mu.Lock()
	if p, ok := a.allocated[allocKey(serviceKey, resourceID)]; ok {
		a.mu.Unlock()
		return p, nil
	}
	a.mu.Unlock()
	return a.Allocate(serviceKey, resourceID)
}

// Allocations returns a snapshot of all current allocations.
func (a *Allocator) Allocations() map[string]int {
	a.mu.Lock()
	defer a.mu.Unlock()
	cp := make(map[string]int, len(a.allocated))
	for k, v := range a.allocated {
		cp[k] = v
	}
	return cp
}
