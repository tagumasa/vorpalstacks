// Package logs provides logging functionality for vorpalstacks.
package logs

import (
	"sync"
	"time"
)

// RotationConfig defines configuration for log rotation.
type RotationConfig struct {
	MaxAge     time.Duration
	MaxSize    int64
	MaxEntries int64
	Compress   bool
}

// DefaultRotationConfig returns the default rotation configuration.
func DefaultRotationConfig() *RotationConfig {
	return &RotationConfig{
		MaxAge:     DefaultLogMaxAge,
		MaxSize:    DefaultMaxLogSize,
		MaxEntries: DefaultMaxLogEntries,
		Compress:   false,
	}
}

// Rotator handles automatic log rotation.
type Rotator struct {
	store    *LogStore
	config   *RotationConfig
	ticker   *time.Ticker
	stopCh   chan struct{}
	stopOnce sync.Once
	wg       sync.WaitGroup
}

// NewRotator creates a new Rotator with the given store and configuration.
func NewRotator(store *LogStore, config *RotationConfig) *Rotator {
	if config == nil {
		config = DefaultRotationConfig()
	}
	return &Rotator{
		store:  store,
		config: config,
		stopCh: make(chan struct{}),
	}
}

// Start begins the automatic rotation process.
func (r *Rotator) Start() {
	r.ticker = time.NewTicker(DefaultRotationInterval)
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for {
			select {
			case <-r.ticker.C:
				_ = r.Rotate()
			case <-r.stopCh:
				return
			}
		}
	}()
}

// Stop halts the automatic rotation process.
func (r *Rotator) Stop() {
	r.stopOnce.Do(func() {
		if r.ticker != nil {
			r.ticker.Stop()
		}
		close(r.stopCh)
		r.wg.Wait()
	})
}

// Rotate performs a single rotation operation.
func (r *Rotator) Rotate() error {
	if r.config.MaxAge > 0 {
		cutoff := time.Now().Add(-r.config.MaxAge)
		return r.store.DeleteBefore(cutoff)
	}
	return nil
}

// RunOnce performs a single rotation and returns.
func (r *Rotator) RunOnce() error {
	return r.Rotate()
}
