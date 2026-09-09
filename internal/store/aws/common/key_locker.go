package common

import "sync"

// KeyLocker provides per-key mutual exclusion. It is designed to protect
// Get-Modify-Put cycles in stores backed by Pebble (where individual
// operations are safe but read-modify-write sequences are not).
//
// Each key gets its own mutex, so concurrent operations on different keys do
// not block each other. Entries are reference-counted: Lock counts a
// goroutine in (including while it waits on the mutex) and Unlock counts it
// out, removing the entry when the last one leaves. A goroutine therefore
// always counts itself in under the locker's map mutex before it can touch
// the entry's mutex, which is what makes cleanup safe for transient keys
// (e.g. S3 multipart upload IDs): deleting an entry whose last holder has
// gone can never strand a waiter on an orphaned mutex, and two goroutines
// can never hold two different mutexes for the same key. This contract
// exists because an earlier unlock-then-delete pattern did break mutual
// exclusion exactly that way.
type KeyLocker struct {
	mu   sync.Mutex
	keys map[string]*keyEntry
}

type keyEntry struct {
	mu   sync.Mutex
	refs int
}

// Lock acquires an exclusive lock for the given key. The caller must call
// Unlock with the same key when the critical section is done.
func (kl *KeyLocker) Lock(key string) {
	kl.mu.Lock()
	if kl.keys == nil {
		kl.keys = map[string]*keyEntry{}
	}
	e, ok := kl.keys[key]
	if !ok {
		e = &keyEntry{}
		kl.keys[key] = e
	}
	e.refs++
	kl.mu.Unlock()
	e.mu.Lock()
}

// Unlock releases the exclusive lock for the given key. The entry is
// removed once no holder or waiter remains.
func (kl *KeyLocker) Unlock(key string) {
	kl.mu.Lock()
	e, ok := kl.keys[key]
	if !ok {
		kl.mu.Unlock()
		panic("keylocker: Unlock without a matching Lock for key " + key)
	}
	e.refs--
	if e.refs == 0 {
		delete(kl.keys, key)
	}
	kl.mu.Unlock()
	e.mu.Unlock()
}

// WithLock acquires the lock for key, runs fn, then releases the lock.
// It returns the error returned by fn.
func (kl *KeyLocker) WithLock(key string, fn func() error) error {
	kl.Lock(key)
	defer kl.Unlock(key)
	return fn()
}

// Delete removes the mutex entry for an idle key. An entry with active
// holders or waiters is left in place and self-cleans at its last Unlock,
// so deletion can never hand two goroutines different mutexes for one key.
func (kl *KeyLocker) Delete(key string) {
	kl.mu.Lock()
	defer kl.mu.Unlock()
	if e, ok := kl.keys[key]; ok && e.refs == 0 {
		delete(kl.keys, key)
	}
}

// DeleteByPrefix removes the idle mutex entries whose key starts with the
// given prefix. Entries with active holders or waiters self-clean at their
// last Unlock.
func (kl *KeyLocker) DeleteByPrefix(prefix string) {
	kl.mu.Lock()
	defer kl.mu.Unlock()
	for k, e := range kl.keys {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix && e.refs == 0 {
			delete(kl.keys, k)
		}
	}
}

// Range calls fn for each key in the locker. If fn returns false, iteration
// stops.
func (kl *KeyLocker) Range(fn func(key string) bool) {
	kl.mu.Lock()
	keys := make([]string, 0, len(kl.keys))
	for k := range kl.keys {
		keys = append(keys, k)
	}
	kl.mu.Unlock()
	for _, k := range keys {
		if !fn(k) {
			return
		}
	}
}
