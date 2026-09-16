package eventbridge

import (
	"context"
	"errors"
	"testing"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// seedReplayStore returns a store with one bus and one enabled archive, the
// minimal stage for replay lifecycle transitions.
func seedReplayStore(t *testing.T) *eventsstore.EventsStore {
	t.Helper()
	ctx := context.Background()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := eventsstore.NewEventsStore(st, "000000000000", "us-east-1")
	if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: "life-bus"}); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateArchive(ctx, &eventsstore.Archive{
		Name:         "life-archive",
		EventBusName: "life-bus",
		State:        eventsstore.ArchiveStateEnabled,
	}); err != nil {
		t.Fatal(err)
	}
	return store
}

// seedReplayRecord writes a replay record directly in the given state.
func seedReplayRecord(t *testing.T, store *eventsstore.EventsStore, name string, state eventsstore.ReplayState) {
	t.Helper()
	if err := store.CreateReplay(context.Background(), &eventsstore.Replay{
		Name:           name,
		EventSourceARN: "arn:aws:events:us-east-1:000000000000:archive/life-archive",
		Destination: &eventsstore.ReplayDestination{
			Arn: "arn:aws:events:us-east-1:000000000000:event-bus/life-bus",
		},
		EventStartTime: time.Now().Add(-time.Hour),
		EventEndTime:   time.Now().Add(time.Hour),
		State:          state,
	}); err != nil {
		t.Fatal(err)
	}
}

func replayState(t *testing.T, store *eventsstore.EventsStore, name string) eventsstore.ReplayState {
	t.Helper()
	replay, err := store.GetReplay(context.Background(), name)
	if err != nil {
		t.Fatalf("get replay %s: %v", name, err)
	}
	return replay.State
}

// TestReplayCancelStateMachine pins the cancellation transitions: a live
// replay (STARTING or RUNNING) moves to CANCELLING; every other state —
// including CANCELLING itself and the terminals — rejects CancelReplay with
// IllegalStatusException.
func TestReplayCancelStateMachine(t *testing.T) {
	ctx := context.Background()
	store := seedReplayStore(t)
	svc := NewEventsService(nil, "000000000000")

	for _, state := range []eventsstore.ReplayState{
		eventsstore.ReplayStateStarting,
		eventsstore.ReplayStateRunning,
	} {
		name := "cancel-from-" + string(state)
		seedReplayRecord(t, store, name, state)
		replay, err := svc.cancelReplayCore(ctx, store, name)
		if err != nil {
			t.Fatalf("cancel from %s must succeed: %v", state, err)
		}
		if replay.State != eventsstore.ReplayStateCancelling {
			t.Fatalf("cancel from %s must land in CANCELLING, got %s", state, replay.State)
		}
		if replay.StateReason != "Cancelled by user" {
			t.Fatalf("the cancelling record must carry the user-cancel reason, got %q", replay.StateReason)
		}
	}

	for _, state := range []eventsstore.ReplayState{
		eventsstore.ReplayStateCancelling,
		eventsstore.ReplayStateCancelled,
		eventsstore.ReplayStateCompleted,
		eventsstore.ReplayStateFailed,
	} {
		name := "cancel-reject-" + string(state)
		seedReplayRecord(t, store, name, state)
		_, err := svc.cancelReplayCore(ctx, store, name)
		var awsErr *awserrors.AWSError
		if !errors.As(err, &awsErr) || awsErr.Code != "IllegalStatusException" {
			t.Fatalf("cancel in state %s must be rejected as IllegalStatusException, got %v", state, err)
		}
		if got := replayState(t, store, name); got != state {
			t.Fatalf("the rejected cancel must leave the record in %s, got %s", state, got)
		}
	}
}

// TestReplayWorkerConvergesCancelling pins the worker side of the state
// machine: a replay whose worker finishes after a cancellation landed (the
// completion write races a CANCELLING record) converges to CANCELLED, and a
// worker aborted by its context in a live state records FAILED — never an
// orphaned STARTING/RUNNING no goroutine will advance.
func TestReplayWorkerConvergesCancelling(t *testing.T) {
	store := seedReplayStore(t)
	svc := NewEventsService(nil, "000000000000")
	archive, err := store.GetArchive(context.Background(), "life-archive")
	if err != nil {
		t.Fatal(err)
	}

	// The completion write racing a CANCELLING record converges CANCELLED.
	seedReplayRecord(t, store, "race-completion", eventsstore.ReplayStateCancelling)
	svc.executeReplay(context.Background(), "us-east-1", mustReplay(t, store, "race-completion"), archive, "life-bus", store)
	if got := replayState(t, store, "race-completion"); got != eventsstore.ReplayStateCancelled {
		t.Fatalf("a completion racing a cancellation must converge to CANCELLED, got %s", got)
	}

	// The abort path: CANCELLING converges CANCELLED.
	seedReplayRecord(t, store, "abort-cancelling", eventsstore.ReplayStateCancelling)
	svc.recordAbortedReplay(store, "abort-cancelling")
	if got := replayState(t, store, "abort-cancelling"); got != eventsstore.ReplayStateCancelled {
		t.Fatalf("an aborted CANCELLING replay must converge to CANCELLED, got %s", got)
	}

	// The abort path: STARTING/RUNNING records FAILED with a reason.
	for _, state := range []eventsstore.ReplayState{
		eventsstore.ReplayStateStarting,
		eventsstore.ReplayStateRunning,
	} {
		name := "abort-" + string(state)
		seedReplayRecord(t, store, name, state)
		svc.recordAbortedReplay(store, name)
		if got := replayState(t, store, name); got != eventsstore.ReplayStateFailed {
			t.Fatalf("an aborted %s replay must record FAILED, got %s", state, got)
		}
		replay, err := store.GetReplay(context.Background(), name)
		if err != nil || replay.StateReason == "" {
			t.Fatalf("the aborted replay must carry a state reason, got %+v (%v)", replay, err)
		}
	}

	// Terminal records are untouched by the abort path.
	seedReplayRecord(t, store, "abort-terminal", eventsstore.ReplayStateCompleted)
	svc.recordAbortedReplay(store, "abort-terminal")
	if got := replayState(t, store, "abort-terminal"); got != eventsstore.ReplayStateCompleted {
		t.Fatalf("a terminal record must stay terminal through the abort path, got %s", got)
	}
}

func mustReplay(t *testing.T, store *eventsstore.EventsStore, name string) *eventsstore.Replay {
	t.Helper()
	replay, err := store.GetReplay(context.Background(), name)
	if err != nil {
		t.Fatal(err)
	}
	return replay
}

// TestReplayProgressPersistsPerPage pins the paged archive walk and its
// progress contract: a replay of an archive whose window holds more events
// than one page replays every event (the walk follows the continuation
// token), EventLastReplayedTime ends at the last replayed event's own
// timestamp ("the time within the specified time range associated with the
// last event replayed", StartReplay operation documentation), the replayed
// events do not re-enter the archive, and the worker's cancel registration
// is consumed when the replay finishes.
func TestReplayProgressPersistsPerPage(t *testing.T) {
	ctx := context.Background()
	svc, store, invoker := newReplayRoutingFixture(t)

	total := replayScanPageSize + 1
	base := time.Now().UTC().Add(-30 * time.Minute)
	for i := 0; i < total; i++ {
		eventTime := base.Add(time.Duration(i) * time.Second)
		event := &eventsstore.Event{
			ID:         eventIDForPage(i),
			Version:    "0",
			DetailType: "PagingTest",
			Source:     "com.example.paging",
			Account:    "000000000000",
			Time:       eventTime,
			Region:     "us-east-1",
			Detail:     map[string]interface{}{"i": i},
		}
		if err := store.StoreArchiveEvent(ctx, "routing-archive", &eventsstore.ArchivedEvent{
			ID:        event.ID,
			Event:     eventEnvelope(event),
			Timestamp: eventTime,
		}); err != nil {
			t.Fatal(err)
		}
	}
	lastEventTime := base.Add(time.Duration(total-1) * time.Second)

	replayName := "paging-replay"
	if _, err := svc.startReplayCore(ctx, store, "us-east-1", StartReplayInput{
		ReplayName:     replayName,
		EventSourceArn: "arn:aws:events:us-east-1:000000000000:archive/routing-archive",
		Destination: map[string]interface{}{
			"Arn": "arn:aws:events:us-east-1:000000000000:event-bus/routing-bus",
		},
		EventStartTime: float64(base.Add(-time.Minute).Unix()),
		EventEndTime:   float64(time.Now().Add(time.Minute).Unix()),
	}); err != nil {
		t.Fatalf("start replay: %v", err)
	}

	// Both rules receive every replayed event: total x 2 deliveries.
	waitForStreamDeliveries(t, invoker, total*2)
	counts := streamCounts(invoker.streams())
	if counts["replay-filtered-rule"] != total || counts["replay-other-rule"] != total {
		t.Fatalf("the paged walk must replay every event to both rules, got %v", counts)
	}

	deadline := time.Now().Add(5 * time.Second)
	for {
		replay := mustReplay(t, store, replayName)
		if replay.State == eventsstore.ReplayStateCompleted {
			if replay.EventLastReplayedTime.IsZero() {
				t.Fatalf("the completed replay must carry EventLastReplayedTime, got %+v", replay)
			}
			if !replay.EventLastReplayedTime.Equal(lastEventTime) {
				t.Fatalf("EventLastReplayedTime must be the last replayed event's own time %v, got %v",
					lastEventTime, replay.EventLastReplayedTime)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("the paged replay did not complete, state %s", replay.State)
		}
		time.Sleep(10 * time.Millisecond)
	}

	// The replayed events never re-enter the archive they were replayed
	// from (still exactly the directly-stored total).
	if got := archivedCount(t, store, "routing-archive"); got != total {
		t.Fatalf("the replay must not grow the archive, want %d events, got %d", total, got)
	}

	// The worker consumed its cancel registration on the way out.
	if _, stillRegistered := svc.replayCancels.Load(replayName); stillRegistered {
		t.Fatal("the finished replay's cancel registration must be consumed, not merely forgotten")
	}
}

func eventIDForPage(i int) string {
	return "paging-event-" + string(rune('a'+i%26)) + string(rune('0'+i/26))
}

// TestStartReplayWindowAndDestinationValidation pins the start-gate
// validation: a window whose start does not precede its end is rejected, and
// a replay onto an archive whose source bus no longer exists reports the
// missing bus instead of silently "completing".
func TestStartReplayWindowAndDestinationValidation(t *testing.T) {
	ctx := context.Background()
	store := seedReplayStore(t)
	svc := NewEventsService(nil, "000000000000")

	now := time.Now().UTC()
	baseInput := func() StartReplayInput {
		return StartReplayInput{
			ReplayName:     "gate-replay",
			EventSourceArn: "arn:aws:events:us-east-1:000000000000:archive/life-archive",
			Destination: map[string]interface{}{
				"Arn": "arn:aws:events:us-east-1:000000000000:event-bus/life-bus",
			},
		}
	}

	// start == end and start > end are both rejected.
	for _, window := range []struct {
		start, end time.Time
	}{
		{now, now},
		{now, now.Add(-time.Minute)},
	} {
		input := baseInput()
		input.EventStartTime = float64(window.start.Unix())
		input.EventEndTime = float64(window.end.Unix())
		_, err := svc.startReplayCore(ctx, store, "us-east-1", input)
		var awsErr *awserrors.AWSError
		if !errors.As(err, &awsErr) || awsErr.Code != "ValidationException" {
			t.Fatalf("an empty or reversed window must be rejected as ValidationException, got %v", err)
		}
	}

	// The archive's source bus is gone: the replay reports the missing bus.
	input := baseInput()
	input.EventStartTime = float64(now.Add(-time.Hour).Unix())
	input.EventEndTime = float64(now.Add(time.Hour).Unix())
	if err := store.DeleteEventBus(ctx, "life-bus"); err != nil {
		t.Fatal(err)
	}
	_, err := svc.startReplayCore(ctx, store, "us-east-1", input)
	var awsErr *awserrors.AWSError
	if !errors.As(err, &awsErr) || awsErr.Code != "ResourceNotFoundException" {
		t.Fatalf("a replay onto a deleted source bus must report ResourceNotFoundException, got %v", err)
	}
}

// TestStartReplayConcurrencyCap pins the active-replay quota: with the
// maximum of ten active replays (STARTING, RUNNING or CANCELLING — terminal
// states never count) already recorded, StartReplay is rejected with
// LimitExceededException.
func TestStartReplayConcurrencyCap(t *testing.T) {
	ctx := context.Background()
	store := seedReplayStore(t)
	svc := NewEventsService(nil, "000000000000")

	activeStates := []eventsstore.ReplayState{
		eventsstore.ReplayStateStarting,
		eventsstore.ReplayStateRunning,
		eventsstore.ReplayStateCancelling,
	}
	for i := 0; i < eventsstore.MaxConcurrentReplays; i++ {
		seedReplayRecord(t, store, "cap-active-"+eventIDForPage(i), activeStates[i%len(activeStates)])
	}
	// Terminal replays never count towards the cap.
	seedReplayRecord(t, store, "cap-done", eventsstore.ReplayStateCompleted)

	// The seeded store holds exactly the cap of active replays (the capped
	// create counts the same states).
	active := 0
	listToken := ""
	for {
		result, err := store.ListReplays(ctx, "", "", "", eventsstore.ListLimitMaximum, listToken)
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range result.Replays {
			switch r.State {
			case eventsstore.ReplayStateStarting, eventsstore.ReplayStateRunning, eventsstore.ReplayStateCancelling:
				active++
			}
		}
		if result.NextToken == "" {
			break
		}
		listToken = result.NextToken
	}
	if active != eventsstore.MaxConcurrentReplays {
		t.Fatalf("the seeded store must hold exactly %d active replays, got %d",
			eventsstore.MaxConcurrentReplays, active)
	}

	now := time.Now().UTC()
	_, err := svc.startReplayCore(ctx, store, "us-east-1", StartReplayInput{
		ReplayName:     "cap-replay",
		EventSourceArn: "arn:aws:events:us-east-1:000000000000:archive/life-archive",
		Destination: map[string]interface{}{
			"Arn": "arn:aws:events:us-east-1:000000000000:event-bus/life-bus",
		},
		EventStartTime: float64(now.Add(-time.Hour).Unix()),
		EventEndTime:   float64(now.Add(time.Hour).Unix()),
	})
	var awsErr *awserrors.AWSError
	if !errors.As(err, &awsErr) || awsErr.Code != "LimitExceededException" {
		t.Fatalf("the eleventh active replay must be rejected as LimitExceededException, got %v", err)
	}
	if _, err := store.GetReplay(ctx, "cap-replay"); !errors.Is(err, eventsstore.ErrReplayNotFound) {
		t.Fatalf("the rejected replay must not be recorded, got %v", err)
	}
}

// TestReplayRecoverySweep pins the startup convergence: replays a previous
// process left in non-terminal states reach a terminal state when the
// region's store first materialises in this process, while already-terminal
// records are untouched.
func TestReplayRecoverySweep(t *testing.T) {
	store := seedReplayStore(t)
	svc := NewEventsService(nil, "000000000000")

	seedReplayRecord(t, store, "sweep-starting", eventsstore.ReplayStateStarting)
	seedReplayRecord(t, store, "sweep-running", eventsstore.ReplayStateRunning)
	seedReplayRecord(t, store, "sweep-cancelling", eventsstore.ReplayStateCancelling)
	seedReplayRecord(t, store, "sweep-completed", eventsstore.ReplayStateCompleted)

	svc.recoverAbandonedReplays(store)

	for name, want := range map[string]eventsstore.ReplayState{
		"sweep-starting":   eventsstore.ReplayStateFailed,
		"sweep-running":    eventsstore.ReplayStateFailed,
		"sweep-cancelling": eventsstore.ReplayStateCancelled,
		"sweep-completed":  eventsstore.ReplayStateCompleted,
	} {
		if got := replayState(t, store, name); got != want {
			t.Fatalf("recovery must converge %s to %s, got %s", name, want, got)
		}
	}
}

// TestReplayRetentionSweep pins the 90-day replay record expiry: records
// older than the retention are deleted, fresh records stay, and a record
// without a creation stamp predating the field is left in place rather than
// deleted by a fabricated timestamp.
func TestReplayRetentionSweep(t *testing.T) {
	ctx := context.Background()
	store := seedReplayStore(t)
	svc := NewEventsService(nil, "000000000000")
	// The retention worker ranges the service's region stores, exactly as
	// the hourly tick finds them in production.
	svc.SetEventsStore("us-east-1", store)
	now := time.Now().UTC()

	seed := func(name string, age time.Duration) {
		t.Helper()
		seedReplayRecord(t, store, name, eventsstore.ReplayStateCompleted)
		if err := store.MutateReplay(ctx, name, func(current *eventsstore.Replay) error {
			current.CreatedAt = now.Add(-age)
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	seed("ttl-old", eventsstore.ReplayRetentionDays*24*time.Hour+time.Hour)
	seed("ttl-fresh", 24*time.Hour)
	seedReplayRecord(t, store, "ttl-unstamped", eventsstore.ReplayStateCompleted) // zero CreatedAt

	svc.purgeExpiredReplays(ctx, now)

	if _, err := store.GetReplay(ctx, "ttl-old"); !errors.Is(err, eventsstore.ErrReplayNotFound) {
		t.Fatalf("a replay past the 90-day retention must be deleted, got %v", err)
	}
	for _, name := range []string{"ttl-fresh", "ttl-unstamped"} {
		if _, err := store.GetReplay(ctx, name); err != nil {
			t.Fatalf("the %s replay must survive the retention sweep: %v", name, err)
		}
	}
}

// TestReplayToListItemShape pins the ListReplays item against the Smithy
// Replay shape (the ReplayList member type): exactly EventEndTime,
// EventLastReplayedTime, EventSourceArn, EventStartTime, ReplayEndTime,
// ReplayName, ReplayStartTime, State, StateReason — no ARN, description
// or destination, which are DescribeReplayResponse members.
func TestReplayToListItemShape(t *testing.T) {
	start := time.Date(2026, 9, 16, 10, 0, 0, 0, time.UTC)
	replay := &eventsstore.Replay{
		Name:                  "shape-probe",
		ARN:                   "arn:aws:events:us-east-1:000000000000:rule/shape-probe",
		State:                 eventsstore.ReplayStateCompleted,
		StateReason:           "completed",
		EventSourceARN:        "arn:aws:events:us-east-1:000000000000:archive/shape-src",
		Description:           "describe-only member",
		EventStartTime:        start,
		EventEndTime:          start.Add(time.Hour),
		ReplayStartTime:       start.Add(2 * time.Hour),
		ReplayEndTime:         start.Add(3 * time.Hour),
		EventLastReplayedTime: start.Add(time.Hour),
		Destination:           &eventsstore.ReplayDestination{Arn: "arn:aws:events:us-east-1:000000000000:event-bus/shape-bus"},
	}

	item := replayToListItem(replay)
	want := map[string]interface{}{
		"ReplayName":            "shape-probe",
		"State":                 string(eventsstore.ReplayStateCompleted),
		"StateReason":           "completed",
		"EventSourceArn":        replay.EventSourceARN,
		"EventStartTime":        start.Unix(),
		"EventEndTime":          start.Add(time.Hour).Unix(),
		"ReplayStartTime":       start.Add(2 * time.Hour).Unix(),
		"ReplayEndTime":         start.Add(3 * time.Hour).Unix(),
		"EventLastReplayedTime": start.Add(time.Hour).Unix(),
	}
	if len(item) != len(want) {
		t.Fatalf("list item has %d members, want %d: %v", len(item), len(want), item)
	}
	for k, v := range want {
		if got, ok := item[k]; !ok || got != v {
			t.Errorf("list item %s = %v (%v), want %v", k, got, ok, v)
		}
	}
}
