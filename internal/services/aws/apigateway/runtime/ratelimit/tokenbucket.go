// Package ratelimit provides the token-bucket rate limiter shared by the
// API Gateway runtime's stage throttling and API-key throttling.
package ratelimit

import (
	"sync"
	"time"
)

// TokenBucket is a token-bucket limiter: Allow consumes one token and
// returns false when the bucket is empty; tokens replenish continuously at
// refillRate per second up to the burst capacity.
type TokenBucket struct {
	mu         sync.Mutex
	tokens     float64
	maxTokens  float64
	refillRate float64
	lastRefill time.Time
}

// New creates a bucket with the given sustained refill rate (tokens per
// second) and burst capacity.
func New(refillRate, burst float64) *TokenBucket {
	return &TokenBucket{
		tokens:     burst,
		maxTokens:  burst,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Update adjusts the rate and burst limits on an existing limiter so
// configuration changes take effect without a restart: a burst increase
// replenishes immediately, a decrease clamps the current tokens.
func (r *TokenBucket) Update(refillRate, burst float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	oldMax := r.maxTokens
	r.refillRate = refillRate
	r.maxTokens = burst
	if burst > oldMax {
		r.tokens = burst
	} else if r.tokens > r.maxTokens {
		r.tokens = r.maxTokens
	}
}

// Allow attempts to consume a token, returning false when the rate limit
// has been exceeded.
func (r *TokenBucket) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastRefill).Seconds()
	r.tokens = min(r.maxTokens, r.tokens+elapsed*r.refillRate)
	r.lastRefill = now

	if r.tokens >= 1 {
		r.tokens--
		return true
	}
	return false
}
