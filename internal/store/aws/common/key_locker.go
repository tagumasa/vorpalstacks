package common

import "sync"

// KeyLocker provides per-key mutual exclusion using a sharded sync.Map of
// mutexes. It is designed to protect Get-Modify-Put cycles in stores backed
// by Pebble (where individual operations are safe but read-modify-write
// sequences are not).
//
// Each key gets its own mutex, so concurrent operations on different keys do
// not block each other. Mutex entries are created on demand via
// sync.Map.LoadOrStore and retained for the lifetime of the KeyLocker.
type KeyLocker struct {
	mu sync.Map
}

// Lock acquires an exclusive lock for the given key. The caller must call
// Unlock with the same key when the critical section is done.
func (kl *KeyLocker) Lock(key string) {
	mu := kl.getMutex(key)
	mu.Lock()
}

// Unlock releases the exclusive lock for the given key.
func (kl *KeyLocker) Unlock(key string) {
	mu := kl.getMutex(key)
	mu.Unlock()
}

// WithLock acquires the lock for key, runs fn, then releases the lock.
// It returns the error returned by fn.
func (kl *KeyLocker) WithLock(key string, fn func() error) error {
	kl.Lock(key)
	defer kl.Unlock(key)
	return fn()
}

// Delete removes the mutex entry for a key. This prevents unbounded growth
// when keys are transient (e.g. S3 multipart upload IDs that are deleted
// after completion).
func (kl *KeyLocker) Delete(key string) {
	kl.mu.Delete(key)
}

// DeleteByPrefix removes all mutex entries whose key starts with the given
// prefix. Useful for cleaning up all locks for a resource that has been
// deleted (e.g. all object locks for a deleted S3 bucket).
func (kl *KeyLocker) DeleteByPrefix(prefix string) {
	kl.mu.Range(func(key, _ any) bool {
		if k, ok := key.(string); ok && len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			kl.mu.Delete(key)
		}
		return true
	})
}

// Range calls fn for each key in the locker. If fn returns false, iteration
// stops. Use with Delete to implement custom cleanup logic.
func (kl *KeyLocker) Range(fn func(key string) bool) {
	kl.mu.Range(func(key, _ any) bool {
		if k, ok := key.(string); ok {
			return fn(k)
		}
		return true
	})
}

func (kl *KeyLocker) getMutex(key string) *sync.Mutex {
	v, _ := kl.mu.LoadOrStore(key, &sync.Mutex{})
	return v.(*sync.Mutex)
}
