package runtime

import (
	"testing"
	"time"
)

func TestStageRateLimiterBurstAllow(t *testing.T) {
	// Burst of 5, high refill rate
	limiter := newStageRateLimiter(10000, 5)

	for i := 0; i < 5; i++ {
		if !limiter.allow() {
			t.Fatalf("request %d should be allowed within burst", i)
		}
	}

	if limiter.allow() {
		t.Fatal("6th request should be denied (burst exhausted)")
	}
}

func TestStageRateLimiterUpdateIncreasesBurst(t *testing.T) {
	limiter := newStageRateLimiter(10000, 2)

	if !limiter.allow() {
		t.Fatal("1st request should be allowed")
	}
	if !limiter.allow() {
		t.Fatal("2nd request should be allowed")
	}
	if limiter.allow() {
		t.Fatal("3rd request should be denied (burst=2)")
	}

	// Update to higher burst — tokens are 0 but maxTokens increases
	limiter.update(10000, 10)

	// Still denied because tokens haven't refilled yet, but allow time for
	// a tiny refill from the high rate
	time.Sleep(2 * time.Millisecond)
	if !limiter.allow() {
		t.Fatal("should be allowed after update + small wait (rate=10000)")
	}
}

func TestStageRateLimiterUpdateDecreasesBurst(t *testing.T) {
	limiter := newStageRateLimiter(10000, 100)

	// Consume some tokens
	for i := 0; i < 10; i++ {
		if !limiter.allow() {
			t.Fatalf("request %d should be allowed", i)
		}
	}

	// Update to lower burst — tokens should be capped
	limiter.update(10000, 5)

	// maxTokens is now 5, but tokens might still be up to 5 after capping.
	// The update caps tokens: if tokens > maxTokens, tokens = maxTokens.
	// tokens was ~90, so after update it's 5.
	for i := 0; i < 5; i++ {
		if !limiter.allow() {
			t.Fatalf("request %d should be allowed after update to burst=5", i)
		}
	}

	if limiter.allow() {
		t.Fatal("should be denied after exhausting updated burst")
	}
}

func TestStageRateLimiterRefill(t *testing.T) {
	// Rate: 100 tokens/sec, Burst: 1
	limiter := newStageRateLimiter(100, 1)

	if !limiter.allow() {
		t.Fatal("1st request should be allowed")
	}
	if limiter.allow() {
		t.Fatal("2nd request should be denied immediately")
	}

	// Wait 20ms → should refill ~2 tokens
	time.Sleep(20 * time.Millisecond)

	if !limiter.allow() {
		t.Fatal("should be allowed after refill")
	}
}
