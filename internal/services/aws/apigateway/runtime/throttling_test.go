package runtime

import (
	"testing"
	"time"

	"vorpalstacks/internal/services/aws/apigateway/runtime/ratelimit"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
)

func TestStageRateLimiterBurstAllow(t *testing.T) {
	// Burst of 5, high refill rate
	limiter := ratelimit.New(10000, 5)

	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Fatalf("request %d should be allowed within burst", i)
		}
	}

	if limiter.Allow() {
		t.Fatal("6th request should be denied (burst exhausted)")
	}
}

func TestStageRateLimiterUpdateIncreasesBurst(t *testing.T) {
	limiter := ratelimit.New(10000, 2)

	if !limiter.Allow() {
		t.Fatal("1st request should be allowed")
	}
	if !limiter.Allow() {
		t.Fatal("2nd request should be allowed")
	}
	if limiter.Allow() {
		t.Fatal("3rd request should be denied (burst=2)")
	}

	// Update to higher burst — tokens are 0 but maxTokens increases
	limiter.Update(10000, 10)

	// Still denied because tokens haven't refilled yet, but allow time for
	// a tiny refill from the high rate
	time.Sleep(2 * time.Millisecond)
	if !limiter.Allow() {
		t.Fatal("should be allowed after update + small wait (rate=10000)")
	}
}

func TestStageRateLimiterUpdateDecreasesBurst(t *testing.T) {
	limiter := ratelimit.New(10000, 100)

	// Consume some tokens
	for i := 0; i < 10; i++ {
		if !limiter.Allow() {
			t.Fatalf("request %d should be allowed", i)
		}
	}

	// Update to lower burst — tokens should be capped
	limiter.Update(10000, 5)

	// maxTokens is now 5, but tokens might still be up to 5 after capping.
	// The update caps tokens: if tokens > maxTokens, tokens = maxTokens.
	// tokens was ~90, so after update it's 5.
	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Fatalf("request %d should be allowed after update to burst=5", i)
		}
	}

	if limiter.Allow() {
		t.Fatal("should be denied after exhausting updated burst")
	}
}

func TestStageRateLimiterRefill(t *testing.T) {
	// Rate: 100 tokens/sec, Burst: 1
	limiter := ratelimit.New(100, 1)

	if !limiter.Allow() {
		t.Fatal("1st request should be allowed")
	}
	if limiter.Allow() {
		t.Fatal("2nd request should be denied immediately")
	}

	// Wait 20ms → should refill ~2 tokens
	time.Sleep(20 * time.Millisecond)

	if !limiter.Allow() {
		t.Fatal("should be allowed after refill")
	}
}

// TestCheckStageThrottlingRawMethodSettingsKeys pins that real-slash keys
// (the spelling whole-map values carry when supplied through CreateStage or
// the whole-row replace) enforce on the execution plane. The documented
// per-setting patch rows store the as-addressed escaped form instead — that
// spelling is pinned by TestCheckStageThrottlingEscapedMethodSettingsKeys.
func TestCheckStageThrottlingRawMethodSettingsKeys(t *testing.T) {
	s := &RuntimeServer{}
	stage := &apigatewaystore.Stage{
		StageName: "raw-keys",
		MethodSettings: map[string]*apigatewaystore.MethodSetting{
			"/pets/GET": {ThrottlingRateLimit: 1, ThrottlingBurstLimit: 1},
		},
	}

	match := &RouteMatch{Path: "/pets", HttpMethod: "GET"}
	if err := s.checkStageThrottling(stage, match); err != nil {
		t.Fatalf("first request within burst should be allowed, got %v", err)
	}
	if err := s.checkStageThrottling(stage, match); err == nil {
		t.Fatal("second request should be throttled by the /pets/GET setting")
	}

	// A different path must not match the exact-key setting: unlimited.
	other := &RouteMatch{Path: "/other", HttpMethod: "GET"}
	for i := 0; i < 5; i++ {
		if err := s.checkStageThrottling(stage, other); err != nil {
			t.Fatalf("request %d on an unconfigured path should not be throttled, got %v", i, err)
		}
	}
}

// TestCheckStageThrottlingWildcardSettingKeys pins the documented wildcard
// entries (*/* and */method), which carry no path component and therefore
// apply to every route.
func TestCheckStageThrottlingWildcardSettingKeys(t *testing.T) {
	s := &RuntimeServer{}
	stage := &apigatewaystore.Stage{
		StageName: "wildcard",
		MethodSettings: map[string]*apigatewaystore.MethodSetting{
			"*/*": {ThrottlingRateLimit: 1, ThrottlingBurstLimit: 1},
		},
	}

	match := &RouteMatch{Path: "/anything", HttpMethod: "POST"}
	if err := s.checkStageThrottling(stage, match); err != nil {
		t.Fatalf("first request within burst should be allowed, got %v", err)
	}
	if err := s.checkStageThrottling(stage, match); err == nil {
		t.Fatal("second request should be throttled by the */* setting")
	}
}

// TestCheckStageThrottlingEscapedMethodSettingsKeys pins the as-addressed
// keys the documented per-setting patch rows store ("~1pets~1{petId}/GET"
// and the per-path method wildcard "~1orders/*"): the escaped spelling of
// the matched route must enforce exactly like the raw one.
func TestCheckStageThrottlingEscapedMethodSettingsKeys(t *testing.T) {
	s := &RuntimeServer{}
	stage := &apigatewaystore.Stage{
		StageName: "escaped-keys",
		MethodSettings: map[string]*apigatewaystore.MethodSetting{
			"~1pets~1{petId}/GET": {ThrottlingRateLimit: 1, ThrottlingBurstLimit: 1},
			"~1orders/*":          {ThrottlingRateLimit: 1, ThrottlingBurstLimit: 1},
		},
	}

	match := &RouteMatch{Path: "/pets/{petId}", HttpMethod: "GET"}
	if err := s.checkStageThrottling(stage, match); err != nil {
		t.Fatalf("first request within burst should be allowed, got %v", err)
	}
	if err := s.checkStageThrottling(stage, match); err == nil {
		t.Fatal("second request should be throttled by the ~1pets~1{petId}/GET setting")
	}

	// The per-path method wildcard covers every method under the path.
	wildcard := &RouteMatch{Path: "/orders", HttpMethod: "POST"}
	if err := s.checkStageThrottling(stage, wildcard); err != nil {
		t.Fatalf("first wildcard-covered request should be allowed, got %v", err)
	}
	if err := s.checkStageThrottling(stage, wildcard); err == nil {
		t.Fatal("second request should be throttled by the ~1orders/* setting")
	}

	// A different path must not match either setting: unlimited.
	other := &RouteMatch{Path: "/other", HttpMethod: "GET"}
	for i := 0; i < 5; i++ {
		if err := s.checkStageThrottling(stage, other); err != nil {
			t.Fatalf("request %d on an unconfigured path should not be throttled, got %v", i, err)
		}
	}
}
