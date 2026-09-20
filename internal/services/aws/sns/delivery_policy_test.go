package sns

import (
	"strings"
	"testing"
	"time"

	snsstore "vorpalstacks/internal/store/aws/sns"
)

// TestDeliveryPolicyDefaults pins the system defaults every omitted field
// resolves to — three retries, twenty-second delays, linear backoff, no
// throttling — on both planes' documents.
func TestDeliveryPolicyDefaults(t *testing.T) {
	for name, parse := range map[string]func(string) (*deliveryPolicy, error){
		"subscription": parseSubscriptionDeliveryPolicy,
		"topic":        parseTopicDeliveryPolicy,
	} {
		policy, err := parse("{}")
		if err != nil {
			t.Fatalf("%s plane empty document: %v", name, err)
		}
		if policy.numRetries != snsstore.DefaultDeliveryPolicyNumRetries ||
			policy.minDelaySeconds != snsstore.DefaultDeliveryPolicyMinDelaySeconds ||
			policy.maxDelaySeconds != snsstore.DefaultDeliveryPolicyMaxDelaySeconds ||
			policy.backoffFunction != "linear" || policy.maxReceivesPerSecond != 0 {
			t.Errorf("%s plane defaults = %+v, want (20s, 20s, 3 retries, linear, no throttle)", name, policy)
		}
	}

	sub, err := parseSubscriptionDeliveryPolicy(`{"healthyRetryPolicy":{"numRetries":5}}`)
	if err != nil {
		t.Fatalf("partial subscription policy: %v", err)
	}
	if sub.numRetries != 5 || sub.minDelaySeconds != snsstore.DefaultDeliveryPolicyMinDelaySeconds {
		t.Errorf("partial policy = %+v, want numRetries 5 with the 20s default delay", sub)
	}

	topic, err := parseTopicDeliveryPolicy(
		`{"http":{"defaultHealthyRetryPolicy":{"numRetries":7},"disableSubscriptionOverrides":true}}`)
	if err != nil {
		t.Fatalf("partial topic policy: %v", err)
	}
	if topic.numRetries != 7 || !topic.disableSubscriptionOverrides {
		t.Errorf("topic policy = %+v, want numRetries 7 with overrides disabled", topic)
	}
}

// TestDeliveryPolicyValidation rejects the documented constraint
// violations on both planes: delay bounds, the retry ceiling, phase counts
// beyond numRetries, the backoff vocabulary, the throttle floor, the
// content-type vocabulary and the 3,600-second total-retry-time hard
// limit — plus unknown fields, which a strict document rejects rather than
// silently ignoring.
func TestDeliveryPolicyValidation(t *testing.T) {
	invalidSub := []string{
		`{"healthyRetryPolicy":{"minDelayTarget":0}}`,
		`{"healthyRetryPolicy":{"minDelayTarget":30,"maxDelayTarget":20}}`,
		`{"healthyRetryPolicy":{"maxDelayTarget":3601}}`,
		`{"healthyRetryPolicy":{"numRetries":101}}`,
		`{"healthyRetryPolicy":{"numRetries":2,"numNoDelayRetries":3}}`,
		`{"healthyRetryPolicy":{"backoffFunction":"quadratic"}}`,
		`{"throttlePolicy":{"maxReceivesPerSecond":0}}`,
		`{"requestPolicy":{"headerContentType":"application/yaml"}}`,
		`{"healthyRetryPolicy":{"minDelayTarget":1,"maxDelayTarget":3600,"numRetries":100}}`,
		`{"healthyPolicy":{}}`,
		`not json`,
	}
	for _, v := range invalidSub {
		if _, err := parseSubscriptionDeliveryPolicy(v); err == nil {
			t.Errorf("subscription policy accepted: %s", v)
		}
	}
	invalidTopic := []string{
		`{"http":{"defaultHealthyRetryPolicy":{"numRetries":101}}}`,
		`{"https":{}}`,
	}
	for _, v := range invalidTopic {
		if _, err := parseTopicDeliveryPolicy(v); err == nil {
			t.Errorf("topic policy accepted: %s", v)
		}
	}

	// The empty value is the documented unset form on both planes.
	if policy, err := parseSubscriptionDeliveryPolicy(""); err != nil || policy != nil {
		t.Errorf("empty subscription policy = (%v, %v), want (nil, nil)", policy, err)
	}
	if policy, err := parseTopicDeliveryPolicy(""); err != nil || policy != nil {
		t.Errorf("empty topic policy = (%v, %v), want (nil, nil)", policy, err)
	}
}

// TestRetryScheduleLadder pins the four-phase ladder: no-delay retries,
// minimum-delay retries, the interpolated backoff phase (the retries left
// of numRetries) and maximum-delay retries, in attempt order.
func TestRetryScheduleLadder(t *testing.T) {
	// The documented example shape: 50 retries split as 3 immediate, 2 at
	// the 1s minimum, the backoff phase, 35 at the 60s maximum.
	policy, err := parseSubscriptionDeliveryPolicy(
		`{"healthyRetryPolicy":{"minDelayTarget":1,"maxDelayTarget":60,"numRetries":50,"numNoDelayRetries":3,"numMinDelayRetries":2,"numMaxDelayRetries":35,"backoffFunction":"exponential"}}`)
	if err != nil {
		t.Fatalf("parse example ladder: %v", err)
	}
	schedule := policy.retrySchedule()
	if len(schedule) != 50 {
		t.Fatalf("schedule length = %d, want numRetries 50", len(schedule))
	}
	for i := 0; i < 3; i++ {
		if schedule[i] != 0 {
			t.Errorf("no-delay retry %d = %s, want 0", i, schedule[i])
		}
	}
	for i := 3; i < 5; i++ {
		if schedule[i] != time.Second {
			t.Errorf("minimum-delay retry %d = %s, want 1s", i, schedule[i])
		}
	}
	for i := 45; i < 50; i++ {
		if schedule[i] != 60*time.Second {
			t.Errorf("maximum-delay retry %d = %s, want 60s", i, schedule[i])
		}
	}
	// The exponential backoff phase starts at the minimum and reaches the
	// maximum on its last step, never leaving the [min, max] band.
	for i, d := range schedule[5:45] {
		if d < time.Second || d > 60*time.Second {
			t.Errorf("backoff delay %d = %s, outside the [1s, 60s] band", i, d)
		}
	}
	if schedule[5] != time.Second {
		t.Errorf("first backoff delay = %s, want the 1s minimum", schedule[5])
	}
	if schedule[44] != 60*time.Second {
		t.Errorf("last backoff delay = %s, want the 60s maximum", schedule[44])
	}
}

// TestBackoffFunctionsShape pins the documented shape ordering: every
// curve starts at the minimum and ends at the maximum, grows monotonically
// and — "exponential … reaching the maximum delay the quickest", "Linear …
// Increases steadily", "Arithmetic and Geometric … steeper than linear but
// less rapid than exponential" — sits between them.
func TestBackoffFunctionsShape(t *testing.T) {
	const min, max, count = 5, 260, 10
	for _, fn := range []string{"linear", "arithmetic", "geometric", "exponential"} {
		var prev time.Duration
		for i := 0; i < count; i++ {
			d := backoffDelay(min, max, i, count, fn)
			if d < time.Duration(min)*time.Second || d > time.Duration(max)*time.Second {
				t.Errorf("%s delay %d = %s, outside the [%ds, %ds] band", fn, i, d, min, max)
			}
			if i > 0 && d < prev {
				t.Errorf("%s delay %d = %s regressed below %s", fn, i, d, prev)
			}
			if i == 0 && d != time.Duration(min)*time.Second {
				t.Errorf("%s first delay = %s, want the %ds minimum", fn, d, min)
			}
			if i == count-1 && d != time.Duration(max)*time.Second {
				t.Errorf("%s last delay = %s, want the %ds maximum", fn, d, max)
			}
			prev = d
		}
	}

	linearMid := backoffDelay(min, max, 5, count, "linear")
	expMid := backoffDelay(min, max, 5, count, "exponential")
	if expMid <= linearMid {
		t.Errorf("exponential mid delay %s does not outpace linear %s", expMid, linearMid)
	}
}

// TestEffectivePolicyResolution pins the override precedence: the
// subscription's policy wins unless the topic disables overrides; the
// topic's default applies when the subscription sets none; the system
// default applies when neither does.
func TestEffectivePolicyResolution(t *testing.T) {
	sub, err := parseSubscriptionDeliveryPolicy(`{"healthyRetryPolicy":{"numRetries":5}}`)
	if err != nil {
		t.Fatalf("sub policy: %v", err)
	}
	topic, err := parseTopicDeliveryPolicy(
		`{"http":{"defaultHealthyRetryPolicy":{"numRetries":9}}}`)
	if err != nil {
		t.Fatalf("topic policy: %v", err)
	}
	topicLocking, err := parseTopicDeliveryPolicy(
		`{"http":{"defaultHealthyRetryPolicy":{"numRetries":9},"disableSubscriptionOverrides":true}}`)
	if err != nil {
		t.Fatalf("locking topic policy: %v", err)
	}

	if got := resolveEffectiveSubscriptionPolicy(sub, topic); got.numRetries != 5 {
		t.Errorf("subscription override = %d retries, want 5", got.numRetries)
	}
	if got := resolveEffectiveSubscriptionPolicy(sub, topicLocking); got.numRetries != 9 {
		t.Errorf("disableSubscriptionOverrides = %d retries, want the topic's 9", got.numRetries)
	}
	if got := resolveEffectiveSubscriptionPolicy(nil, topic); got.numRetries != 9 {
		t.Errorf("topic default = %d retries, want 9", got.numRetries)
	}
	if got := resolveEffectiveSubscriptionPolicy(nil, nil); got.numRetries != snsstore.DefaultDeliveryPolicyNumRetries {
		t.Errorf("no policies = %d retries, want the system default %d", got.numRetries, snsstore.DefaultDeliveryPolicyNumRetries)
	}
}

// TestEffectivePolicyRendering pins the synthesised EffectiveDeliveryPolicy
// JSON on both planes: every documented retry member present at its
// resolved value, in the plane's own document shape.
func TestEffectivePolicyRendering(t *testing.T) {
	def := systemDefaultDeliveryPolicy()

	subJSON, err := def.effectiveSubscriptionJSON()
	if err != nil {
		t.Fatalf("render subscription: %v", err)
	}
	for _, want := range []string{
		`"healthyRetryPolicy":{`,
		`"minDelayTarget":20`, `"maxDelayTarget":20`, `"numRetries":3`,
		`"numNoDelayRetries":0`, `"numMinDelayRetries":0`, `"numMaxDelayRetries":0`,
		`"backoffFunction":"linear"`,
	} {
		if !strings.Contains(subJSON, want) {
			t.Errorf("subscription effective JSON missing %s: %s", want, subJSON)
		}
	}
	if strings.Contains(subJSON, "throttlePolicy") || strings.Contains(subJSON, "requestPolicy") {
		t.Errorf("unset throttle/request policies rendered: %s", subJSON)
	}

	topicJSON, err := def.effectiveTopicJSON()
	if err != nil {
		t.Fatalf("render topic: %v", err)
	}
	for _, want := range []string{`"http":{`, `"defaultHealthyRetryPolicy":{`, `"disableSubscriptionOverrides":false`} {
		if !strings.Contains(topicJSON, want) {
			t.Errorf("topic effective JSON missing %s: %s", want, topicJSON)
		}
	}

	withThrottle, err := parseSubscriptionDeliveryPolicy(
		`{"throttlePolicy":{"maxReceivesPerSecond":10},"requestPolicy":{"headerContentType":"application/json"}}`)
	if err != nil {
		t.Fatalf("throttle policy: %v", err)
	}
	rendered, err := withThrottle.effectiveSubscriptionJSON()
	if err != nil {
		t.Fatalf("render throttle: %v", err)
	}
	if !strings.Contains(rendered, `"maxReceivesPerSecond":10`) ||
		!strings.Contains(rendered, `"headerContentType":"application/json"`) {
		t.Errorf("set throttle/request policies not rendered: %s", rendered)
	}
}
