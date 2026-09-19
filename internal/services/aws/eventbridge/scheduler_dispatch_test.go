package eventbridge

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// streamRecordingInvoker records the stream name of every PutRecord so a
// test can assert which rule's targets received a delivery.
type streamRecordingInvoker struct {
	recordingKinesisInvoker
	mu     sync.Mutex
	stream []string
}

func (r *streamRecordingInvoker) PutRecord(ctx context.Context, region, streamName, partitionKey string, data []byte) (string, string, error) {
	r.mu.Lock()
	r.stream = append(r.stream, streamName)
	r.mu.Unlock()
	return r.recordingKinesisInvoker.PutRecord(ctx, region, streamName, partitionKey, data)
}

func (r *streamRecordingInvoker) streams() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.stream...)
}

// TestScheduledRulesFireOnlyTheirOwnTargets pins the scheduler's dispatch
// scope: two scheduled rules on one bus each deliver their scheduled event to
// their own targets exactly once — never to each other's, and never through
// bus-wide matching to pattern-less bystanders.
func TestScheduledRulesFireOnlyTheirOwnTargets(t *testing.T) {
	ctx := context.Background()

	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := eventsstore.NewEventsStore(st, "000000000000", "us-east-1")
	if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: "sched-bus"}); err != nil {
		t.Fatal(err)
	}

	// A pattern-less bystander rule must not receive the scheduled fires;
	// its target stream would record any cross-delivery.
	rules := []struct {
		name         string
		pattern      string
		schedule     string
		targetStream string
	}{
		{"minutely-alpha", "", "rate(1 minute)", "sched-alpha"},
		{"minutely-beta", "", "rate(1 minute)", "sched-beta"},
		{"bystander", `{"source":["never.scheduled"]}`, "", "sched-bystander"},
	}
	now := time.Now().UTC().Truncate(time.Minute)
	for _, r := range rules {
		rule := &eventsstore.Rule{
			Name:         r.name,
			EventBusName: "sched-bus",
			EventPattern: r.pattern,
			State:        eventsstore.RuleStateEnabled,
		}
		if r.schedule != "" {
			rule.ScheduleExpression = r.schedule
		}
		if err := store.CreateRule(ctx, rule); err != nil {
			t.Fatal(err)
		}
		if r.targetStream != "" {
			if err := store.PutTarget(ctx, &eventsstore.Target{
				ID:           "stream",
				RuleName:     r.name,
				EventBusName: "sched-bus",
				ARN:          "arn:aws:kinesis:us-east-1:000000000000:stream/" + r.targetStream,
			}); err != nil {
				t.Fatal(err)
			}
		}
	}
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	invoker := &streamRecordingInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	if err := bus.Start(ctx); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })

	svc := NewEventsService(mgr, "000000000000")
	svc.SetEventsStore("us-east-1", store)
	svc.SetEventBus(bus)
	t.Cleanup(svc.Close)

	// The rate rules anchor to their creation, which falls somewhere inside
	// the truncated minute; a tick two minutes past that minute therefore
	// has exactly one elapsed boundary per rule regardless of the second
	// the test started at.
	svc.tickScheduledRules(ctx, now.Add(2*time.Minute))

	deadline := time.Now().Add(5 * time.Second)
	for len(invoker.streams()) < 2 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(200 * time.Millisecond)

	streams := invoker.streams()
	if len(streams) != 2 {
		t.Fatalf("expected exactly two deliveries (one per scheduled rule), got %v", streams)
	}
	counts := map[string]int{}
	for _, s := range streams {
		counts[s]++
	}
	if counts["sched-alpha"] != 1 || counts["sched-beta"] != 1 {
		t.Fatalf("each scheduled rule must deliver to its own stream exactly once, got %v", counts)
	}
	if counts["sched-bystander"] != 0 {
		t.Fatalf("the pattern-less bystander must not receive scheduled fires, got %v", counts)
	}

	// The delivered payload is the scheduled-event envelope.
	payloads := invoker.payloads()
	decoded, err := base64.StdEncoding.DecodeString(string(payloads[0]))
	if err != nil {
		t.Fatalf("the delivered payload must be base64-encoded: %v", err)
	}
	var envelope map[string]interface{}
	if err := json.Unmarshal(decoded, &envelope); err != nil {
		t.Fatalf("the delivered payload must be a JSON envelope: %v", err)
	}
	if envelope["detail-type"] != "Scheduled Event" || envelope["source"] != "aws.events" {
		t.Fatalf("the scheduled fire must carry the scheduled-event envelope, got %v", envelope)
	}
}

func ruleARNFor(busName, ruleName string) string {
	return "arn:aws:events:us-east-1:000000000000:rule/" + busName + "/" + ruleName
}

// TestScheduledFireReachesPatternSiblings pins the bus visibility the
// events reference documents: "EventBridge itself emits the following
// events. These events are automatically sent to the default event bus as
// with any other AWS service." — an enabled rule whose pattern matches
// source aws.events receives the fire alongside the firing rule's own
// targets, with the reference's sample envelope: the firing rule's ARN in
// resources and an empty detail object.
func TestScheduledFireReachesPatternSiblings(t *testing.T) {
	ctx := context.Background()

	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := eventsstore.NewEventsStore(st, "000000000000", "us-east-1")
	if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: "sched-sibling-bus"}); err != nil {
		t.Fatal(err)
	}
	rules := []struct {
		name         string
		pattern      string
		schedule     string
		targetStream string
	}{
		{"minutely-alpha", "", "rate(1 minute)", "sched-alpha"},
		{"aws-events-listener", `{"source":["aws.events"]}`, "", "sched-listener"},
	}
	now := time.Now().UTC().Truncate(time.Minute)
	for _, r := range rules {
		rule := &eventsstore.Rule{
			Name:         r.name,
			EventBusName: "sched-sibling-bus",
			EventPattern: r.pattern,
			State:        eventsstore.RuleStateEnabled,
		}
		if r.schedule != "" {
			rule.ScheduleExpression = r.schedule
		}
		if err := store.CreateRule(ctx, rule); err != nil {
			t.Fatal(err)
		}
		if err := store.PutTarget(ctx, &eventsstore.Target{
			ID:           "stream",
			RuleName:     r.name,
			EventBusName: "sched-sibling-bus",
			ARN:          "arn:aws:kinesis:us-east-1:000000000000:stream/" + r.targetStream,
		}); err != nil {
			t.Fatal(err)
		}
	}
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	invoker := &streamRecordingInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	if err := bus.Start(ctx); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })

	svc := NewEventsService(mgr, "000000000000")
	svc.SetEventsStore("us-east-1", store)
	svc.SetEventBus(bus)
	t.Cleanup(svc.Close)

	svc.tickScheduledRules(ctx, now.Add(2*time.Minute))

	waitForStreamDeliveries(t, invoker, 2)
	counts := streamCounts(invoker.streams())
	if counts["sched-alpha"] != 1 || counts["sched-listener"] != 1 {
		t.Fatalf("the fire must reach the firing rule's own target and the matching sibling exactly once each, got %v", counts)
	}

	// Every delivered envelope is the reference's sample shape: the firing
	// rule's ARN in resources, an empty detail object.
	alphaARN := ruleARNFor("sched-sibling-bus", "minutely-alpha")
	for _, payload := range invoker.payloads() {
		decoded, err := base64.StdEncoding.DecodeString(string(payload))
		if err != nil {
			t.Fatalf("the delivered payload must be base64-encoded: %v", err)
		}
		var envelope map[string]interface{}
		if err := json.Unmarshal(decoded, &envelope); err != nil {
			t.Fatalf("the delivered payload must be a JSON envelope: %v", err)
		}
		if envelope["detail-type"] != "Scheduled Event" || envelope["source"] != "aws.events" {
			t.Fatalf("the fire must carry the scheduled-event envelope, got %v", envelope)
		}
		resources, _ := envelope["resources"].([]interface{})
		if len(resources) != 1 || resources[0] != alphaARN {
			t.Fatalf("the envelope must carry the firing rule's ARN in resources, got %v", envelope["resources"])
		}
		detail, _ := envelope["detail"].(map[string]interface{})
		if len(detail) != 0 {
			t.Fatalf("the scheduled-event detail must be empty, got %v", detail)
		}
	}
}

// panickingRulesBucket poisons the rule listing: a corrupted rule record
// surfaces as a panic inside the storage iteration.
type panickingRulesBucket struct {
	storage.Bucket
}

func (b *panickingRulesBucket) ForEach(fn func(k, v []byte) error) error {
	panic("poisoned rule record")
}

// panickingStorage routes the rules bucket through the panicking view.
type panickingStorage struct {
	storage.BasicStorage
}

func (s *panickingStorage) Bucket(name string) storage.Bucket {
	if name == "events-rules-us-east-1" {
		return &panickingRulesBucket{}
	}
	return s.BasicStorage.Bucket(name)
}

// TestSchedulerTickSurvivesPanic pins the tick's panic boundary: a panic in
// the rule listing is contained and the scheduler stays usable for the next
// tick instead of dying for the process lifetime.
func TestSchedulerTickSurvivesPanic(t *testing.T) {
	ctx := context.Background()

	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := eventsstore.NewEventsStore(&panickingStorage{BasicStorage: st}, "000000000000", "us-east-1")

	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewEventsService(mgr, "000000000000")
	svc.SetEventsStore("us-east-1", store)

	// The guarded tick must return normally; an unguarded panic would fail
	// the test run.
	svc.tickScheduledRulesGuarded(ctx, time.Now().UTC())
	// A second guarded tick keeps working after the poisoned one.
	svc.tickScheduledRulesGuarded(ctx, time.Now().UTC().Add(time.Minute))
}

// TestUnevaluableScheduleNeverFiresAndLogsOnce pins the unevaluable-expression
// path: such a rule never fires, and the once-per-expression marker is set so
// the warning is not repeated every tick.
func TestUnevaluableScheduleNeverFiresAndLogsOnce(t *testing.T) {
	arn := ruleARNFor("sched-bus", "unevaluable-probe")
	key := arn + "\x00" + "not-a-schedule-expression"
	unevaluableScheduleLogged.Delete(key)
	defer unevaluableScheduleLogged.Delete(key)

	var dedup scheduleFireDedup
	now := time.Now().UTC()
	if dedup.shouldFireSchedule(arn, "not-a-schedule-expression", now, now) {
		t.Fatal("an unevaluable schedule expression must never fire")
	}
	if dedup.shouldFireSchedule(arn, "not-a-schedule-expression", now.Add(time.Minute), now) {
		t.Fatal("an unevaluable schedule expression must never fire on later ticks")
	}
	if _, marked := unevaluableScheduleLogged.Load(key); !marked {
		t.Fatal("the unevaluable expression must be marked as logged exactly once")
	}
}
