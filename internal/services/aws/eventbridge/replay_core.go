package eventbridge

import (
	"context"
	"fmt"
	"strconv"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
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

	if err := store.CreateReplay(ctx, replay); err != nil {
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
				current, err := store.GetReplay(context.Background(), replayName)
				if err == nil && current.State != eventsstore.ReplayStateCancelled {
					current.State = eventsstore.ReplayStateFailed
					current.StateReason = fmt.Sprintf("Internal error: %v", r)
					current.ReplayEndTime = time.Now().UTC()
					_ = store.UpdateReplay(context.Background(), current)
				}
				s.replayCancels.Delete(replayName)
			}
		}()
		s.executeReplay(replayCtx, reqCtxRegion, replay, archive, destEventBusName, store)
	}()

	return replay, nil
}

// executeReplay drives an asynchronous replay: transitions the record to
// running, replays the archived events onto the destination bus and marks the
// replay completed unless a cancellation arrived meanwhile.
func (s *EventsService) executeReplay(ctx context.Context, region string, replay *eventsstore.Replay, archive *eventsstore.Archive, destEventBusName string, store *eventsstore.EventsStore) {
	defer s.replayCancels.Delete(replay.Name)

	replay.State = eventsstore.ReplayStateRunning
	replay.ReplayStartTime = time.Now().UTC()
	if err := store.UpdateReplay(ctx, replay); err != nil {
		logs.Warn("failed to update replay state to running",
			logs.String("replayName", replay.Name),
			logs.Err(err))
	}

	events, err := store.GetArchiveEvents(ctx, archive.Name, replay.EventStartTime, replay.EventEndTime)
	if err != nil {
		replay.State = eventsstore.ReplayStateFailed
		replay.StateReason = "Failed to retrieve archived events: " + err.Error()
		replay.ReplayEndTime = time.Now().UTC()
		if updateErr := store.UpdateReplay(ctx, replay); updateErr != nil {
			logs.Warn("failed to update replay state to failed",
				logs.String("replayName", replay.Name),
				logs.Err(updateErr))
		}
		return
	}

	replayedCount := int64(0)
	for _, archivedEvent := range events {
		select {
		case <-ctx.Done():
			logs.Info("replay cancelled during event processing",
				logs.String("replayName", replay.Name))
			return
		default:
		}

		if err := s.replayEventToBus(ctx, region, archivedEvent, destEventBusName, store); err != nil {
			logs.Warn("failed to replay event",
				logs.String("eventId", archivedEvent.ID),
				logs.String("replayName", replay.Name),
				logs.Err(err))
			continue
		}
		replayedCount++
	}

	// Re-fetch to avoid overwriting a cancellation that arrived while replaying.
	if current, err := store.GetReplay(ctx, replay.Name); err == nil &&
		current.State != eventsstore.ReplayStateCancelled {
		current.State = eventsstore.ReplayStateCompleted
		current.ReplayEndTime = time.Now().UTC()
		if updateErr := store.UpdateReplay(ctx, current); updateErr != nil {
			logs.Warn("failed to update replay state to completed",
				logs.String("replayName", replay.Name),
				logs.Err(updateErr))
		}
	}

	logs.Info("replay completed",
		logs.String("replayName", replay.Name),
		logs.Int("eventsReplayed", int(replayedCount)))
}

// replayEventToBus rebuilds an Event from an archived event envelope and
// delivers it onto the destination bus.
func (s *EventsService) replayEventToBus(ctx context.Context, region string, archivedEvent *eventsstore.ArchivedEvent, destEventBusName string, store *eventsstore.EventsStore) error {
	eventMap := archivedEvent.Event

	event := &eventsstore.Event{
		ID:           request.GetStringParam(eventMap, "id"),
		Version:      request.GetStringParam(eventMap, "version"),
		DetailType:   request.GetStringParam(eventMap, "detail-type"),
		Source:       request.GetStringParam(eventMap, "source"),
		Account:      request.GetStringParam(eventMap, "account"),
		Region:       request.GetStringParam(eventMap, "region"),
		EventBusName: destEventBusName,
	}

	if timeStr := request.GetStringParam(eventMap, "time"); timeStr != "" {
		if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
			event.Time = t
		}
	}

	if detail, ok := eventMap["detail"].(map[string]interface{}); ok {
		event.Detail = detail
	}

	if resources, ok := eventMap["resources"].([]interface{}); ok {
		for _, r := range resources {
			if rStr, ok := r.(string); ok {
				event.Resources = append(event.Resources, rStr)
			}
		}
	}

	return s.deliverEventWithStore(ctx, region, event, destEventBusName, store)
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

// listReplaysCore applies the documented limit window and lists the replays.
func (s *EventsService) listReplaysCore(ctx context.Context, store *eventsstore.EventsStore, input ListReplaysInput) (*eventsstore.ReplayListResult, error) {
	limit := input.Limit
	if limit < 0 || limit > 100 {
		return nil, awserrors.NewValidationException("Limit must be between 0 and 100")
	}
	if limit == 0 {
		limit = 50
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

	replay, err := store.GetReplay(ctx, replayName)
	if err != nil {
		return nil, mapStoreError(err, replayName)
	}

	if replay.State != eventsstore.ReplayStateRunning && replay.State != eventsstore.ReplayStateStarting {
		return nil, awserrors.NewIllegalStatusException("Replay cannot be cancelled in state: " + string(replay.State))
	}

	replay.State = eventsstore.ReplayStateCancelled
	replay.StateReason = "Cancelled by user"
	replay.ReplayEndTime = time.Now().UTC()

	if err := store.UpdateReplay(ctx, replay); err != nil {
		return nil, err
	}

	if val, ok := s.replayCancels.LoadAndDelete(replayName); ok {
		if cancelFn, ok := val.(context.CancelFunc); ok {
			cancelFn()
		}
	}

	return replay, nil
}
