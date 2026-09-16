package eventbridge

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/logs"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	"vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Input / Result structs (transport-agnostic)
// ---------------------------------------------------------------------------

// StartReplayInput carries the parameters for StartReplay. Destination and
// the event time bounds hold the raw wire values because their validation
// runs after the archive existence check, so the error precedence must be
// preserved at the Core layer.
type StartReplayInput struct {
	ReplayName     string
	EventSourceArn string
	DescriptionSet bool
	Description    string
	Destination    interface{}
	EventStartTime interface{}
	EventEndTime   interface{}
	Region         string
}

// ListReplaysInput carries the parameters for ListReplays.
type ListReplaysInput struct {
	NamePrefix     string
	State          eventsstore.ReplayState
	EventSourceArn string
	Limit          int32
	NextToken      string
}

// ---------------------------------------------------------------------------
// Core functions
// ---------------------------------------------------------------------------

// startReplayCore validates input, creates the replay record and launches
// the asynchronous replay worker.
func (s *EventsService) startReplayCore(ctx context.Context, store *eventsstore.EventsStore, reqCtxRegion string, input StartReplayInput) (*eventsstore.Replay, error) {
	replayName := input.ReplayName
	if replayName == "" {
		return nil, awserrors.NewValidationException("ReplayName is required")
	}
	if !validateReplayName(replayName) {
		return nil, awserrors.NewValidationException("ReplayName must be 1-64 characters matching [.-_A-Za-z0-9]")
	}

	eventSourceArn := input.EventSourceArn
	if eventSourceArn == "" {
		return nil, awserrors.NewValidationException("EventSourceArn is required")
	}

	archiveName := arn.ExtractArchiveNameFromARN(eventSourceArn)
	if archiveName == "" {
		return nil, awserrors.NewValidationException("Invalid EventSourceArn")
	}

	archive, err := store.GetArchive(ctx, archiveName)
	if err != nil {
		return nil, mapStoreError(err, archiveName)
	}

	var destination *eventsstore.ReplayDestination
	if destMap, ok := input.Destination.(map[string]interface{}); ok {
		destination = &eventsstore.ReplayDestination{}
		if arn, ok := destMap["Arn"].(string); ok {
			destination.Arn = arn
		}
		if filterArns, ok := destMap["FilterArns"].([]interface{}); ok {
			for _, fa := range filterArns {
				if faStr, ok := fa.(string); ok {
					destination.FilterArns = append(destination.FilterArns, faStr)
				}
			}
		}
	}
	if destination == nil || destination.Arn == "" {
		return nil, awserrors.NewValidationException("Destination.Arn is required")
	}

	destEventBusName := arn.ExtractEventBusNameFromARN(destination.Arn)
	if destEventBusName == "" {
		return nil, awserrors.NewValidationException("Invalid Destination.Arn")
	}

	// "The ARN of the event bus to replay event to. You can replay events
	// only to the event bus specified to create the archive" (Smithy
	// ReplayDestination.Arn member documentation): the destination bus is
	// constrained to the archive's source bus, so a replay onto any other
	// bus — including an empty or cross-account one — is rejected instead of
	// silently delivering archived events to an unrelated bus.
	if destEventBusName != archive.EventBusName {
		return nil, awserrors.NewValidationException(
			"Destination.Arn must identify the event bus the archive was created on; events can only be replayed to the source event bus")
	}

	// The source bus must still exist: an archive outliving its deleted bus
	// cannot receive a replay onto that bus.
	if _, err := store.GetEventBus(ctx, destEventBusName); err != nil {
		return nil, mapStoreError(err, destEventBusName)
	}

	var eventStartTime, eventEndTime time.Time
	if startTimeVal := input.EventStartTime; startTimeVal != nil {
		switch v := startTimeVal.(type) {
		case float64:
			eventStartTime = time.Unix(int64(v), 0)
		case string:
			if unix, err := strconv.ParseInt(v, 10, 64); err == nil {
				eventStartTime = time.Unix(unix, 0)
			} else if t, err := time.Parse(time.RFC3339, v); err == nil {
				eventStartTime = t
			}
		}
	}
	if endTimeVal := input.EventEndTime; endTimeVal != nil {
		switch v := endTimeVal.(type) {
		case float64:
			eventEndTime = time.Unix(int64(v), 0)
		case string:
			if unix, err := strconv.ParseInt(v, 10, 64); err == nil {
				eventEndTime = time.Unix(unix, 0)
			} else if t, err := time.Parse(time.RFC3339, v); err == nil {
				eventEndTime = t
			}
		}
	}

	if eventStartTime.IsZero() || eventEndTime.IsZero() {
		return nil, awserrors.NewValidationException("EventStartTime and EventEndTime are required")
	}

	// Only events that occurred between EventStartTime and EventEndTime are
	// replayed (model member documentation), so a window whose start does
	// not precede its end is empty by construction — rejected rather than
	// replayed as a silently-completing no-op.
	if !eventStartTime.Before(eventEndTime) {
		return nil, awserrors.NewValidationException("EventStartTime must be before EventEndTime")
	}

	replay := &eventsstore.Replay{
		Name:           replayName,
		EventSourceARN: eventSourceArn,
		Destination:    destination,
		EventStartTime: eventStartTime,
		EventEndTime:   eventEndTime,
		State:          eventsstore.ReplayStateStarting,
	}

	if input.DescriptionSet {
		if !validateDescription(input.Description) {
			return nil, errDescriptionTooLong()
		}
		replay.Description = input.Description
	}

	// "You can have a maximum of ten active concurrent replays per account
	// per AWS Region" (the archives user guide); the model carries
	// LimitExceededException on StartReplay for the breach. Active means a
	// non-terminal state: STARTING, RUNNING or CANCELLING. The count and the
	// create run as one store-level locked section, so concurrent
	// StartReplay calls cannot both pass the gate.
	if err := store.CreateReplayCapped(ctx, replay, eventsstore.MaxConcurrentReplays); err != nil {
		if errors.Is(err, eventsstore.ErrReplayCapReached) {
			return nil, awserrors.NewLimitExceededException(fmt.Sprintf(
				"The maximum of %d active concurrent replays per account per region has been reached", eventsstore.MaxConcurrentReplays))
		}
		return nil, mapStoreError(err, replayName)
	}

	replayCtx, cancel := context.WithCancel(context.Background())
	s.replayCancels.Store(replayName, cancel)
	s.replayWg.Add(1)
	go func() {
		defer s.replayWg.Done()
		defer func() {
			if r := recover(); r != nil {
				logs.Error("eventbridge: replay goroutine panicked",
					logs.String("replayName", replayName),
					logs.Any("panic", r))
				_ = store.MutateReplay(context.Background(), replayName, func(current *eventsstore.Replay) error {
					// A cancellation that arrived before the panic owns the
					// terminal state; it only needs its CANCELLING record
					// converged to CANCELLED here.
					if current.State == eventsstore.ReplayStateCancelled {
						return nil
					}
					if current.State == eventsstore.ReplayStateCancelling {
						current.State = eventsstore.ReplayStateCancelled
						current.ReplayEndTime = time.Now().UTC()
						return nil
					}
					current.State = eventsstore.ReplayStateFailed
					current.StateReason = fmt.Sprintf("Internal error: %v", r)
					current.ReplayEndTime = time.Now().UTC()
					return nil
				})
				s.replayCancels.Delete(replayName)
			}
		}()
		s.executeReplay(replayCtx, reqCtxRegion, replay, archive, destEventBusName, store)
	}()

	return replay, nil
}

// replayScanPageSize bounds one page of the archive scan during a replay:
// the walk streams pages instead of materialising the archive's entire
// filtered set, and progress (EventLastReplayedTime) persists once per page.
const replayScanPageSize = 100

// executeReplay drives an asynchronous replay: transitions the record to
// running, walks the archive's events for the window in pages and replays
// them onto the destination bus, persisting progress per page, and records
// the terminal state atomically — a cancellation that arrived while
// replaying decides the terminal state by lock order, never by a stale
// pre-read record.
func (s *EventsService) executeReplay(ctx context.Context, region string, replay *eventsstore.Replay, archive *eventsstore.Archive, destEventBusName string, store *eventsstore.EventsStore) {
	// The stored cancel func is consumed as well as deleted: the WithCancel
	// child of Background stays registered with its parent until cancel
	// runs, so a replay that finishes without CancelReplay must still
	// release the context tree.
	defer func() {
		if v, ok := s.replayCancels.LoadAndDelete(replay.Name); ok {
			if cancelFn, ok := v.(context.CancelFunc); ok {
				cancelFn()
			}
		}
	}()

	if err := store.MutateReplay(ctx, replay.Name, func(current *eventsstore.Replay) error {
		// Only a STARTING record enters RUNNING: a replay cancelled before
		// the worker started keeps its CANCELLING state for the walk's
		// abort path to converge.
		if current.State == eventsstore.ReplayStateStarting {
			current.State = eventsstore.ReplayStateRunning
			current.ReplayStartTime = time.Now().UTC()
		}
		return nil
	}); err != nil {
		logs.Warn("failed to update replay state to running",
			logs.String("replayName", replay.Name),
			logs.Err(err))
	}

	// The destination's FilterArns restricts the delivery fan-out to the
	// named rules ("A list of ARNs for rules to replay events to", Smithy
	// ReplayDestination.FilterArns member documentation); an absent or empty
	// list leaves the fan-out unfiltered, as before.
	var filterRuleARNs map[string]bool
	if replay.Destination != nil {
		for _, filterARN := range replay.Destination.FilterArns {
			if filterRuleARNs == nil {
				filterRuleARNs = make(map[string]bool, len(replay.Destination.FilterArns))
			}
			filterRuleARNs[filterARN] = true
		}
	}

	replayedCount := int64(0)
	failedCount := int64(0)
	// lastReplayed tracks the event time of the most recently replayed
	// event: EventLastReplayedTime "indicates the time within the specified
	// time range associated with the last event replayed" (StartReplay
	// operation documentation) — the event's own timestamp, not wall clock.
	lastReplayed := time.Time{}
	nextToken := ""
	for {
		select {
		case <-ctx.Done():
			s.recordAbortedReplay(store, replay.Name)
			return
		default:
		}

		page, err := store.ListArchiveEvents(ctx, archive.Name, replay.EventStartTime, replay.EventEndTime, replayScanPageSize, nextToken)
		if err != nil {
			s.recordFailedReplay(store, replay.Name, "Failed to retrieve archived events: "+err.Error())
			return
		}

		for _, archivedEvent := range page.Events {
			select {
			case <-ctx.Done():
				s.recordAbortedReplay(store, replay.Name)
				return
			default:
			}

			if err := s.replayEventToBus(ctx, region, archivedEvent, destEventBusName, store, replay.Name, filterRuleARNs); err != nil {
				failedCount++
				logs.Warn("failed to replay event",
					logs.String("eventId", archivedEvent.ID),
					logs.String("replayName", replay.Name),
					logs.Err(err))
				continue
			}
			replayedCount++
			lastReplayed = archivedEvent.Timestamp
		}

		// Progress persists once per page and only advances a live replay:
		// a record that reached a cancelling or terminal state meanwhile
		// keeps its state and timestamps untouched.
		if !lastReplayed.IsZero() {
			if err := store.MutateReplay(ctx, replay.Name, func(current *eventsstore.Replay) error {
				switch current.State {
				case eventsstore.ReplayStateStarting, eventsstore.ReplayStateRunning:
					current.EventLastReplayedTime = lastReplayed
				}
				return nil
			}); err != nil {
				logs.Warn("failed to persist replay progress",
					logs.String("replayName", replay.Name),
					logs.Err(err))
			}
		}

		if page.NextToken == "" {
			break
		}
		nextToken = page.NextToken
	}

	// The completion write re-reads inside the record mutation: a
	// cancellation that arrived while replaying decides the terminal state
	// by lock order, never by a stale pre-read record.
	if err := store.MutateReplay(ctx, replay.Name, func(current *eventsstore.Replay) error {
		switch current.State {
		case eventsstore.ReplayStateCancelled:
			return nil
		case eventsstore.ReplayStateCancelling:
			// The cancellation won the race against the worker's finish;
			// the worker is done, so the terminal write converges the
			// record here.
			current.State = eventsstore.ReplayStateCancelled
			current.ReplayEndTime = time.Now().UTC()
			return nil
		}
		current.State = eventsstore.ReplayStateCompleted
		if failedCount > 0 {
			current.StateReason = fmt.Sprintf("%d event(s) failed to replay", failedCount)
		}
		current.ReplayEndTime = time.Now().UTC()
		return nil
	}); err != nil {
		logs.Warn("failed to update replay state to completed",
			logs.String("replayName", replay.Name),
			logs.Err(err))
	}

	logs.Info("replay completed",
		logs.String("replayName", replay.Name),
		logs.Int("eventsReplayed", int(replayedCount)),
		logs.Int("eventsFailed", int(failedCount)))
}

// recordAbortedReplay converges a replay whose worker stopped without
// finishing: a CANCELLING record converges to CANCELLED (a user
// cancellation whose terminal write belonged to this worker), while a
// STARTING or RUNNING record means the worker's context died without
// CancelReplay — an aborted run recorded FAILED, never an orphaned live
// state no goroutine will advance.
func (s *EventsService) recordAbortedReplay(store *eventsstore.EventsStore, replayName string) {
	if err := store.MutateReplay(context.Background(), replayName, func(current *eventsstore.Replay) error {
		switch current.State {
		case eventsstore.ReplayStateCancelling:
			current.State = eventsstore.ReplayStateCancelled
			current.ReplayEndTime = time.Now().UTC()
			return nil
		case eventsstore.ReplayStateStarting, eventsstore.ReplayStateRunning:
			current.State = eventsstore.ReplayStateFailed
			current.StateReason = "Replay worker aborted before completion"
			current.ReplayEndTime = time.Now().UTC()
			return nil
		}
		return nil
	}); err != nil {
		logs.Warn("failed to record aborted replay",
			logs.String("replayName", replayName),
			logs.Err(err))
	}
}

// recordFailedReplay writes the replay's FAILED terminal state with the
// retrieval fault as its reason.
func (s *EventsService) recordFailedReplay(store *eventsstore.EventsStore, replayName, reason string) {
	if err := store.MutateReplay(context.Background(), replayName, func(current *eventsstore.Replay) error {
		if current.State == eventsstore.ReplayStateCancelling || current.State == eventsstore.ReplayStateCancelled {
			return nil
		}
		current.State = eventsstore.ReplayStateFailed
		current.StateReason = reason
		current.ReplayEndTime = time.Now().UTC()
		return nil
	}); err != nil {
		logs.Warn("failed to update replay state to failed",
			logs.String("replayName", replayName),
			logs.Err(err))
	}
}

// replayEventToBus rebuilds an Event from an archived event envelope and
// delivers it onto the destination bus. The rebuilt event carries the
// replay-name stamp — EventBridge "adds a metadata field to the event,
// replay-name" while replaying (the archives user guide), which downstream
// consumers match with the documented {"replay-name":[{"exists":false}]}
// pattern — and the delivery fan-out is restricted to filterRuleARNs when
// the replay destination named FilterArns.
func (s *EventsService) replayEventToBus(ctx context.Context, region string, archivedEvent *eventsstore.ArchivedEvent, destEventBusName string, store *eventsstore.EventsStore, replayName string, filterRuleARNs map[string]bool) error {
	event := eventFromEnvelope(archivedEvent.Event, destEventBusName)
	event.ReplayName = replayName
	return s.deliverEventToBusRules(ctx, region, event, destEventBusName, store, 0, filterRuleARNs)
}

// describeReplayCore validates input and fetches the replay.
func (s *EventsService) describeReplayCore(ctx context.Context, store *eventsstore.EventsStore, replayName string) (*eventsstore.Replay, error) {
	if replayName == "" {
		return nil, awserrors.NewValidationException("ReplayName is required")
	}
	if !validateReplayName(replayName) {
		return nil, awserrors.NewValidationException("ReplayName must be 1-64 characters matching [.-_A-Za-z0-9]")
	}

	replay, err := store.GetReplay(ctx, replayName)
	if err != nil {
		return nil, mapStoreError(err, replayName)
	}
	return replay, nil
}

// listReplaysCore validates the query and lists the replays: the limit
// window, the NamePrefix charset (both per the model) and the State
// filter's enum membership.
func (s *EventsService) listReplaysCore(ctx context.Context, store *eventsstore.EventsStore, input ListReplaysInput) (*eventsstore.ReplayListResult, error) {
	if input.NamePrefix != "" {
		if err := validateListNamePrefix(input.NamePrefix, "replay"); err != nil {
			return nil, err
		}
	}
	if input.State != "" && !eventsstore.IsValidReplayState(input.State) {
		return nil, awserrors.NewValidationException(
			"State must be a member of the ReplayState enum: " + strings.Join(eventsstore.ReplayStateVocabulary(), ", "))
	}
	limit, err := normaliseListLimit(input.Limit)
	if err != nil {
		return nil, err
	}
	return store.ListReplays(ctx, input.NamePrefix, input.State, input.EventSourceArn, limit, input.NextToken)
}

// cancelReplayCore validates input, cancels a running replay and signals its
// worker goroutine.
func (s *EventsService) cancelReplayCore(ctx context.Context, store *eventsstore.EventsStore, replayName string) (*eventsstore.Replay, error) {
	if replayName == "" {
		return nil, awserrors.NewValidationException("ReplayName is required")
	}
	if !validateReplayName(replayName) {
		return nil, awserrors.NewValidationException("ReplayName must be 1-64 characters matching [.-_A-Za-z0-9]")
	}

	// The state gate and the CANCELLING write run as one atomic record
	// mutation: a replay that completed or failed between the core's read
	// and the write can no longer be cancelled by a stale pre-read record.
	// The CANCELLED terminal write belongs to the worker — it lands when
	// the replay loop observes the cancellation, converging the record it
	// finds in CANCELLING.
	var cancelling *eventsstore.Replay
	if err := store.MutateReplay(ctx, replayName, func(current *eventsstore.Replay) error {
		if current.State != eventsstore.ReplayStateRunning && current.State != eventsstore.ReplayStateStarting {
			return awserrors.NewIllegalStatusException("Replay cannot be cancelled in state: " + string(current.State))
		}

		current.State = eventsstore.ReplayStateCancelling
		current.StateReason = "Cancelled by user"
		cancelling = current
		return nil
	}); err != nil {
		return nil, mapStoreError(err, replayName)
	}

	if val, ok := s.replayCancels.LoadAndDelete(replayName); ok {
		if cancelFn, ok := val.(context.CancelFunc); ok {
			cancelFn()
		}
	}

	return cancelling, nil
}
