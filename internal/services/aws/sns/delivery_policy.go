package sns

// The HTTP/S delivery policy: parsing, validation, default resolution,
// schedule computation and effective-policy rendering. A delivery policy is
// "composed of a retry policy, throttle policy and a request policy" with
// nine attributes in total (Amazon SNS message delivery retries); the
// subscription plane stores it as
// {"healthyRetryPolicy":…,"throttlePolicy":…,"requestPolicy":…} and the
// topic plane as the same three policies under "http" with
// default-prefixed names plus disableSubscriptionOverrides.

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	snsstore "vorpalstacks/internal/store/aws/sns"
)

// ---------------------------------------------------------------------------
// Wire documents
// ---------------------------------------------------------------------------

// Every numeric and vocabulary field is a pointer so an omitted field (the
// system default applies) stays distinguishable from an explicit zero —
// numRetries omitted means the default of three, numRetries set to 0 means
// no retries at all.
type wireRetryPolicy struct {
	MinDelayTarget     *int    `json:"minDelayTarget"`
	MaxDelayTarget     *int    `json:"maxDelayTarget"`
	NumRetries         *int    `json:"numRetries"`
	NumNoDelayRetries  *int    `json:"numNoDelayRetries"`
	NumMinDelayRetries *int    `json:"numMinDelayRetries"`
	NumMaxDelayRetries *int    `json:"numMaxDelayRetries"`
	BackoffFunction    *string `json:"backoffFunction"`
}

type wireThrottlePolicy struct {
	MaxReceivesPerSecond *int `json:"maxReceivesPerSecond"`
}

type wireRequestPolicy struct {
	HeaderContentType *string `json:"headerContentType"`
}

// wireSubscriptionDeliveryPolicy is the subscription-level document shape.
// The optional members marshal omitempty so the effective-policy rendering
// omits unset throttle and request policies instead of emitting nulls.
type wireSubscriptionDeliveryPolicy struct {
	HealthyRetryPolicy *wireRetryPolicy    `json:"healthyRetryPolicy,omitempty"`
	ThrottlePolicy     *wireThrottlePolicy `json:"throttlePolicy,omitempty"`
	RequestPolicy      *wireRequestPolicy  `json:"requestPolicy,omitempty"`
}

// wireTopicDeliveryPolicy is the topic-level document shape: the same three
// policies under "http", named as the topic defaults they define.
type wireTopicDeliveryPolicy struct {
	HTTP *struct {
		DefaultHealthyRetryPolicy    *wireRetryPolicy    `json:"defaultHealthyRetryPolicy,omitempty"`
		DisableSubscriptionOverrides bool                `json:"disableSubscriptionOverrides"`
		DefaultThrottlePolicy        *wireThrottlePolicy `json:"defaultThrottlePolicy,omitempty"`
		DefaultRequestPolicy         *wireRequestPolicy  `json:"defaultRequestPolicy,omitempty"`
	} `json:"http,omitempty"`
}

// ---------------------------------------------------------------------------
// Resolved policy
// ---------------------------------------------------------------------------

// deliveryPolicy is a fully resolved delivery policy: every retry field
// carries its system default when the document omitted it, so the delivery
// engine and the effective-policy rendering read one shape.
type deliveryPolicy struct {
	minDelaySeconds              int
	maxDelaySeconds              int
	numRetries                   int
	numNoDelayRetries            int
	numMinDelayRetries           int
	numMaxDelayRetries           int
	backoffFunction              string
	maxReceivesPerSecond         int // 0 = no throttling (documented default)
	headerContentType            string
	disableSubscriptionOverrides bool // topic-plane switch, false on subscription documents
}

// validBackoffFunctions is the retry-backoff vocabulary: "One of four
// options: arithmetic, exponential, geometric, linear. Default: linear".
var validBackoffFunctions = map[string]bool{
	"linear":      true,
	"arithmetic":  true,
	"geometric":   true,
	"exponential": true,
}

// validHeaderContentTypes is the union of the two documented content-type
// sets — the narrower non-raw set (application/json, text/plain) and the
// raw-delivery set (text/css, text/csv, text/html, text/plain, text/xml,
// application/atom+xml, application/json, application/octet-stream,
// application/soap+xml, application/x-www-form-urlencoded,
// application/xhtml+xml, application/xml). Which set applies depends on the
// subscription's raw setting at delivery time, which attribute-write order
// cannot pin down, so the union gates both.
var validHeaderContentTypes = map[string]bool{
	"text/css":                          true,
	"text/csv":                          true,
	"text/html":                         true,
	"text/plain":                        true,
	"text/xml":                          true,
	"application/atom+xml":              true,
	"application/json":                  true,
	"application/octet-stream":          true,
	"application/soap+xml":              true,
	"application/x-www-form-urlencoded": true,
	"application/xhtml+xml":             true,
	"application/xml":                   true,
}

// systemDefaultDeliveryPolicy returns the policy the platform applies when
// neither the subscription nor the topic defines one — every field at its
// documented system default.
func systemDefaultDeliveryPolicy() *deliveryPolicy {
	return &deliveryPolicy{
		minDelaySeconds:      snsstore.DefaultDeliveryPolicyMinDelaySeconds,
		maxDelaySeconds:      snsstore.DefaultDeliveryPolicyMaxDelaySeconds,
		numRetries:           snsstore.DefaultDeliveryPolicyNumRetries,
		numNoDelayRetries:    0,
		numMinDelayRetries:   0,
		numMaxDelayRetries:   0,
		backoffFunction:      "linear",
		maxReceivesPerSecond: 0,
	}
}

// decodeDeliveryPolicyDocument decodes one delivery-policy JSON document
// strictly: a field name outside the documented nine is rejected rather
// than silently ignored, because an ignored typo would make the effective
// policy differ from the caller's intent.
func decodeDeliveryPolicyDocument(raw string, dst interface{}) error {
	dec := json.NewDecoder(strings.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return NewInvalidParameter(fmt.Sprintf("Invalid delivery policy: %s", err.Error()))
	}
	var extra json.Token
	if err := dec.Decode(&extra); err == nil {
		return NewInvalidParameter("Invalid delivery policy: unexpected trailing content")
	}
	return nil
}

// resolveRetryPolicy validates one wire retry policy against the documented
// constraints and fills every omitted field with its system default:
// minDelayTarget 1..maxDelayTarget (default 20), maxDelayTarget
// minDelayTarget..3600 (default 20), numRetries 0..100 (default 3), the
// three phase counts 0 or greater (default 0), and backoffFunction one of
// linear/arithmetic/geometric/exponential (default linear). The phase
// counts together may not exceed numRetries — the backoff phase takes
// "numRetries - numNoDelayRetries - numMinDelayRetries - numMaxDelayRetries".
func resolveRetryPolicy(name string, wire *wireRetryPolicy) (minDelay, maxDelay, numRetries, noDelay, minDelayRetries, maxDelayRetries int, backoff string, err error) {
	pick := func(v *int, def int) int {
		if v == nil {
			return def
		}
		return *v
	}
	minDelay = pick(wire.MinDelayTarget, snsstore.DefaultDeliveryPolicyMinDelaySeconds)
	maxDelay = pick(wire.MaxDelayTarget, snsstore.DefaultDeliveryPolicyMaxDelaySeconds)
	numRetries = pick(wire.NumRetries, snsstore.DefaultDeliveryPolicyNumRetries)
	noDelay = pick(wire.NumNoDelayRetries, 0)
	minDelayRetries = pick(wire.NumMinDelayRetries, 0)
	maxDelayRetries = pick(wire.NumMaxDelayRetries, 0)
	backoff = "linear"
	if wire.BackoffFunction != nil {
		backoff = *wire.BackoffFunction
	}

	if minDelay < snsstore.MinDeliveryPolicyDelaySeconds {
		return 0, 0, 0, 0, 0, 0, "", NewInvalidParameter(fmt.Sprintf(
			"Invalid delivery policy: %s.minDelayTarget must be at least %d seconds (got %d)",
			name, snsstore.MinDeliveryPolicyDelaySeconds, minDelay))
	}
	if maxDelay > snsstore.MaxDeliveryPolicySeconds {
		return 0, 0, 0, 0, 0, 0, "", NewInvalidParameter(fmt.Sprintf(
			"Invalid delivery policy: %s.maxDelayTarget must be at most %d seconds (got %d)",
			name, snsstore.MaxDeliveryPolicySeconds, maxDelay))
	}
	if minDelay > maxDelay {
		return 0, 0, 0, 0, 0, 0, "", NewInvalidParameter(fmt.Sprintf(
			"Invalid delivery policy: %s.minDelayTarget (%d) must not exceed maxDelayTarget (%d)",
			name, minDelay, maxDelay))
	}
	if numRetries < 0 || numRetries > snsstore.MaxDeliveryPolicyRetries {
		return 0, 0, 0, 0, 0, 0, "", NewInvalidParameter(fmt.Sprintf(
			"Invalid delivery policy: %s.numRetries must be between 0 and %d (got %d)",
			name, snsstore.MaxDeliveryPolicyRetries, numRetries))
	}
	if noDelay < 0 || minDelayRetries < 0 || maxDelayRetries < 0 {
		return 0, 0, 0, 0, 0, 0, "", NewInvalidParameter(fmt.Sprintf(
			"Invalid delivery policy: %s retry phase counts must be 0 or greater", name))
	}
	if noDelay+minDelayRetries+maxDelayRetries > numRetries {
		return 0, 0, 0, 0, 0, 0, "", NewInvalidParameter(fmt.Sprintf(
			"Invalid delivery policy: %s retry phase counts (%d + %d + %d) exceed numRetries (%d)",
			name, noDelay, minDelayRetries, maxDelayRetries, numRetries))
	}
	if !validBackoffFunctions[backoff] {
		return 0, 0, 0, 0, 0, 0, "", NewInvalidParameter(fmt.Sprintf(
			"Invalid delivery policy: %s.backoffFunction %q is not one of arithmetic, exponential, geometric, linear",
			name, backoff))
	}
	return minDelay, maxDelay, numRetries, noDelay, minDelayRetries, maxDelayRetries, backoff, nil
}

// resolveDeliveryPolicyFields assembles a resolved policy from validated
// wire fields and enforces the whole-document constraint the field checks
// cannot see: the computed retry schedule's total delay must stay within
// "the total policy retry time for an HTTP/S endpoint cannot be greater
// than 3,600 seconds".
func resolveDeliveryPolicyFields(retryName string, retry *wireRetryPolicy, throttle *wireThrottlePolicy, request *wireRequestPolicy, disableOverrides bool) (*deliveryPolicy, error) {
	policy := systemDefaultDeliveryPolicy()
	policy.disableSubscriptionOverrides = disableOverrides

	if retry != nil {
		minDelay, maxDelay, numRetries, noDelay, minDelayRetries, maxDelayRetries, backoff, err :=
			resolveRetryPolicy(retryName, retry)
		if err != nil {
			return nil, err
		}
		policy.minDelaySeconds = minDelay
		policy.maxDelaySeconds = maxDelay
		policy.numRetries = numRetries
		policy.numNoDelayRetries = noDelay
		policy.numMinDelayRetries = minDelayRetries
		policy.numMaxDelayRetries = maxDelayRetries
		policy.backoffFunction = backoff
	}
	if throttle != nil && throttle.MaxReceivesPerSecond != nil {
		if *throttle.MaxReceivesPerSecond < snsstore.MinThrottleReceivesPerSecond {
			return nil, NewInvalidParameter(fmt.Sprintf(
				"Invalid delivery policy: throttlePolicy.maxReceivesPerSecond must be %d or greater (got %d)",
				snsstore.MinThrottleReceivesPerSecond, *throttle.MaxReceivesPerSecond))
		}
		policy.maxReceivesPerSecond = *throttle.MaxReceivesPerSecond
	}
	if request != nil && request.HeaderContentType != nil {
		if !validHeaderContentTypes[*request.HeaderContentType] {
			return nil, NewInvalidParameter(fmt.Sprintf(
				"Invalid delivery policy: requestPolicy.headerContentType %q is not a supported content type",
				*request.HeaderContentType))
		}
		policy.headerContentType = *request.HeaderContentType
	}

	var total time.Duration
	for _, d := range policy.retrySchedule() {
		total += d
	}
	if total > time.Duration(snsstore.MaxDeliveryPolicySeconds)*time.Second {
		return nil, NewInvalidParameter(fmt.Sprintf(
			"Invalid delivery policy: total retry time %s exceeds the hard limit of %d seconds",
			total.Truncate(time.Second), snsstore.MaxDeliveryPolicySeconds))
	}
	return policy, nil
}

// parseSubscriptionDeliveryPolicy validates and resolves a
// subscription-level DeliveryPolicy value; the empty value is the documented
// unset form and yields nil.
func parseSubscriptionDeliveryPolicy(value string) (*deliveryPolicy, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	var wire wireSubscriptionDeliveryPolicy
	if err := decodeDeliveryPolicyDocument(value, &wire); err != nil {
		return nil, err
	}
	return resolveDeliveryPolicyFields("healthyRetryPolicy",
		wire.HealthyRetryPolicy, wire.ThrottlePolicy, wire.RequestPolicy, false)
}

// parseTopicDeliveryPolicy validates and resolves a topic-level
// DeliveryPolicy value; the empty value is the documented unset form and
// yields nil. A document whose "http" member is absent carries no settings
// and resolves to the system defaults.
func parseTopicDeliveryPolicy(value string) (*deliveryPolicy, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	var wire wireTopicDeliveryPolicy
	if err := decodeDeliveryPolicyDocument(value, &wire); err != nil {
		return nil, err
	}
	if wire.HTTP == nil {
		return systemDefaultDeliveryPolicy(), nil
	}
	return resolveDeliveryPolicyFields("defaultHealthyRetryPolicy",
		wire.HTTP.DefaultHealthyRetryPolicy, wire.HTTP.DefaultThrottlePolicy,
		wire.HTTP.DefaultRequestPolicy, wire.HTTP.DisableSubscriptionOverrides)
}

// ---------------------------------------------------------------------------
// Effective-policy resolution and rendering
// ---------------------------------------------------------------------------

// resolveEffectiveSubscriptionPolicy selects the policy one HTTP/S delivery
// runs under: the subscription's own policy, unless the topic's policy
// disables subscription overrides; the topic's default when the
// subscription defines none; the system default when neither does.
func resolveEffectiveSubscriptionPolicy(subPolicy, topicPolicy *deliveryPolicy) *deliveryPolicy {
	if topicPolicy != nil && topicPolicy.disableSubscriptionOverrides {
		return topicPolicy
	}
	if subPolicy != nil {
		return subPolicy
	}
	if topicPolicy != nil {
		return topicPolicy
	}
	return systemDefaultDeliveryPolicy()
}

// renderWireRetryPolicy projects a resolved policy's retry fields back into
// the wire shape with every member present — the
// EffectiveDeliveryPolicy rendering "taking system defaults into account".
func renderWireRetryPolicy(p *deliveryPolicy) *wireRetryPolicy {
	return &wireRetryPolicy{
		MinDelayTarget:     &p.minDelaySeconds,
		MaxDelayTarget:     &p.maxDelaySeconds,
		NumRetries:         &p.numRetries,
		NumNoDelayRetries:  &p.numNoDelayRetries,
		NumMinDelayRetries: &p.numMinDelayRetries,
		NumMaxDelayRetries: &p.numMaxDelayRetries,
		BackoffFunction:    &p.backoffFunction,
	}
}

func (p *deliveryPolicy) renderThrottle() *wireThrottlePolicy {
	if p.maxReceivesPerSecond == 0 {
		return nil
	}
	return &wireThrottlePolicy{MaxReceivesPerSecond: &p.maxReceivesPerSecond}
}

func (p *deliveryPolicy) renderRequest() *wireRequestPolicy {
	if p.headerContentType == "" {
		return nil
	}
	return &wireRequestPolicy{HeaderContentType: &p.headerContentType}
}

// effectiveSubscriptionJSON renders the effective delivery policy in the
// subscription-plane document shape (GetSubscriptionAttributes'
// EffectiveDeliveryPolicy).
func (p *deliveryPolicy) effectiveSubscriptionJSON() (string, error) {
	doc := wireSubscriptionDeliveryPolicy{
		HealthyRetryPolicy: renderWireRetryPolicy(p),
		ThrottlePolicy:     p.renderThrottle(),
		RequestPolicy:      p.renderRequest(),
	}
	b, err := json.Marshal(doc)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// effectiveTopicJSON renders the effective delivery policy in the
// topic-plane document shape (GetTopicAttributes' EffectiveDeliveryPolicy).
func (p *deliveryPolicy) effectiveTopicJSON() (string, error) {
	doc := wireTopicDeliveryPolicy{
		HTTP: &struct {
			DefaultHealthyRetryPolicy    *wireRetryPolicy    `json:"defaultHealthyRetryPolicy,omitempty"`
			DisableSubscriptionOverrides bool                `json:"disableSubscriptionOverrides"`
			DefaultThrottlePolicy        *wireThrottlePolicy `json:"defaultThrottlePolicy,omitempty"`
			DefaultRequestPolicy         *wireRequestPolicy  `json:"defaultRequestPolicy,omitempty"`
		}{
			DefaultHealthyRetryPolicy:    renderWireRetryPolicy(p),
			DisableSubscriptionOverrides: p.disableSubscriptionOverrides,
			DefaultThrottlePolicy:        p.renderThrottle(),
			DefaultRequestPolicy:         p.renderRequest(),
		},
	}
	b, err := json.Marshal(doc)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ---------------------------------------------------------------------------
// Retry schedule
// ---------------------------------------------------------------------------

// retrySchedule computes the delay before every retry attempt, in attempt
// order: numNoDelayRetries immediate retries, numMinDelayRetries retries at
// the minimum delay, the backoff phase (the retries left of numRetries)
// interpolating from the minimum to the maximum delay by the backoff
// function, and numMaxDelayRetries retries at the maximum delay. The
// initial delivery attempt is not part of this schedule.
func (p *deliveryPolicy) retrySchedule() []time.Duration {
	schedule := make([]time.Duration, 0, p.numRetries)
	for i := 0; i < p.numNoDelayRetries; i++ {
		schedule = append(schedule, 0)
	}
	for i := 0; i < p.numMinDelayRetries; i++ {
		schedule = append(schedule, time.Duration(p.minDelaySeconds)*time.Second)
	}
	backoffCount := p.numRetries - p.numNoDelayRetries - p.numMinDelayRetries - p.numMaxDelayRetries
	for i := 0; i < backoffCount; i++ {
		schedule = append(schedule, backoffDelay(
			p.minDelaySeconds, p.maxDelaySeconds, i, backoffCount, p.backoffFunction))
	}
	for i := 0; i < p.numMaxDelayRetries; i++ {
		schedule = append(schedule, time.Duration(p.maxDelaySeconds)*time.Second)
	}
	return schedule
}

// backoffDelay interpolates one backoff-phase delay. The current AWS
// documentation describes the four functions by shape — every curve
// "starts near the minimum delay and approaches the maximum delay by the
// last retry", linear "Increases steadily with each attempt", exponential
// "reaching the maximum delay the quickest", and "Arithmetic and Geometric
// ... steeper than linear but less rapid than exponential" — without
// publishing the formulas. These implementations match the documented
// shapes: each curve is the linear position x = attempt/(count−1) bent by
// a concave transform 1−(1−x)^p with p = 2 (exponential), 1.5 (arithmetic)
// and 1.25 (geometric), so every non-linear curve rises faster than linear
// and exponential rises fastest of all, while all of them start exactly at
// the minimum and end exactly at the maximum.
func backoffDelay(minSeconds, maxSeconds, attempt, count int, fn string) time.Duration {
	if count <= 1 {
		return time.Duration(minSeconds) * time.Second
	}
	x := float64(attempt) / float64(count-1)
	switch fn {
	case "arithmetic":
		x = 1 - math.Pow(1-x, 1.5)
	case "geometric":
		x = 1 - math.Pow(1-x, 1.25)
	case "exponential":
		x = 1 - math.Pow(1-x, 2)
	}
	delay := float64(minSeconds) + float64(maxSeconds-minSeconds)*x
	return time.Duration(math.Round(delay)) * time.Second
}
