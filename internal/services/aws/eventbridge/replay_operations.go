package eventbridge

import (
	"context"
	"strconv"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/common/request"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	"vorpalstacks/internal/utils/aws/arn"
)

// StartReplay starts an event replay from an archive.
func (s *EventsService) StartReplay(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	replayName := request.GetParamLowerFirst(req.Parameters, "ReplayName")
	if replayName == "" {
		return nil, NewValidationException("ReplayName is required")
	}

	eventSourceArn := request.GetParamLowerFirst(req.Parameters, "EventSourceArn")
	if eventSourceArn == "" {
		return nil, NewValidationException("EventSourceArn is required")
	}

	archiveName := arn.ExtractArchiveNameFromARN(eventSourceArn)
	if archiveName == "" {
		return nil, NewValidationException("Invalid EventSourceArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	archive, err := store.GetArchive(ctx, archiveName)
	if err != nil {
		if err == eventsstore.ErrArchiveNotFound {
			return nil, NewResourceNotFoundException("Archive '" + archiveName + "' does not exist")
		}
		return nil, err
	}

	var destination *eventsstore.ReplayDestination
	if destMap, ok := req.Parameters["Destination"].(map[string]interface{}); ok {
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
		return nil, NewValidationException("Destination.Arn is required")
	}

	destEventBusName := arn.ExtractEventBusNameFromARN(destination.Arn)
	if destEventBusName == "" {
		return nil, NewValidationException("Invalid Destination.Arn")
	}

	var eventStartTime, eventEndTime time.Time
	if startTimeVal, ok := req.Parameters["EventStartTime"]; ok {
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
	if endTimeVal, ok := req.Parameters["EventEndTime"]; ok {
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
		return nil, NewValidationException("EventStartTime and EventEndTime are required")
	}

	replay := &eventsstore.Replay{
		Name:           replayName,
		EventSourceARN: eventSourceArn,
		Destination:    destination,
		EventStartTime: eventStartTime,
		EventEndTime:   eventEndTime,
		State:          eventsstore.ReplayStateStarting,
	}

	if desc, ok := req.Parameters["Description"].(string); ok {
		replay.Description = desc
	}

	if err := store.CreateReplay(ctx, replay); err != nil {
		if err == eventsstore.ErrReplayAlreadyExists {
			return nil, NewResourceAlreadyExistsException("Replay '" + replayName + "' already exists")
		}
		return nil, err
	}

	replayCtx, cancel := context.WithCancel(context.Background())
	s.replayCancels.Store(replayName, cancel)
	go s.executeReplay(replayCtx, reqCtx.GetRegion(), replay, archive, destEventBusName, store)

	result := map[string]interface{}{
		"ReplayArn":      replay.ARN,
		"ReplayName":     replay.Name,
		"State":          string(replay.State),
		"EventSourceArn": replay.EventSourceARN,
		"Destination": map[string]interface{}{
			"Arn":        destination.Arn,
			"FilterArns": destination.FilterArns,
		},
		"EventStartTime": replay.EventStartTime.Unix(),
		"EventEndTime":   replay.EventEndTime.Unix(),
	}
	if replay.StateReason != "" {
		result["StateReason"] = replay.StateReason
	}
	return result, nil
}

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

func (s *EventsService) replayEventToBus(ctx context.Context, region string, archivedEvent *eventsstore.ArchivedEvent, destEventBusName string, store *eventsstore.EventsStore) error {
	eventMap := archivedEvent.Event

	event := &eventsstore.Event{
		ID:           getStringFromMap(eventMap, "id"),
		Version:      getStringFromMap(eventMap, "version"),
		DetailType:   getStringFromMap(eventMap, "detail-type"),
		Source:       getStringFromMap(eventMap, "source"),
		Account:      getStringFromMap(eventMap, "account"),
		Region:       getStringFromMap(eventMap, "region"),
		EventBusName: destEventBusName,
	}

	if timeStr := getStringFromMap(eventMap, "time"); timeStr != "" {
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

func getStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// DescribeReplay returns information about a replay.
func (s *EventsService) DescribeReplay(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	replayName := request.GetParamLowerFirst(req.Parameters, "ReplayName")
	if replayName == "" {
		return nil, NewValidationException("ReplayName is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	replay, err := store.GetReplay(ctx, replayName)
	if err != nil {
		if err == eventsstore.ErrReplayNotFound {
			return nil, NewResourceNotFoundException("Replay '" + replayName + "' does not exist")
		}
		return nil, err
	}

	result := map[string]interface{}{
		"ReplayName":     replay.Name,
		"ReplayArn":      replay.ARN,
		"State":          string(replay.State),
		"EventSourceArn": replay.EventSourceARN,
		"EventStartTime": replay.EventStartTime.Unix(),
		"EventEndTime":   replay.EventEndTime.Unix(),
	}

	if replay.Description != "" {
		result["Description"] = replay.Description
	}
	if replay.StateReason != "" {
		result["StateReason"] = replay.StateReason
	}
	if replay.Destination != nil {
		result["Destination"] = map[string]interface{}{
			"Arn":        replay.Destination.Arn,
			"FilterArns": replay.Destination.FilterArns,
		}
	}
	if !replay.ReplayStartTime.IsZero() {
		result["ReplayStartTime"] = replay.ReplayStartTime.Unix()
	}
	if !replay.ReplayEndTime.IsZero() {
		result["ReplayEndTime"] = replay.ReplayEndTime.Unix()
	}
	if !replay.EventLastReplayedTime.IsZero() {
		result["EventLastReplayedTime"] = replay.EventLastReplayedTime.Unix()
	}

	return result, nil
}

// ListReplays returns a list of replays.
func (s *EventsService) ListReplays(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namePrefix := request.GetParamLowerFirst(req.Parameters, "NamePrefix")

	var state eventsstore.ReplayState
	if stateStr := request.GetParamLowerFirst(req.Parameters, "State"); stateStr != "" {
		state = eventsstore.ReplayState(stateStr)
	}

	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	if limit == 0 {
		limit = 50
	}
	if limit > 100 {
		return nil, NewValidationException("Limit must be between 1 and 100")
	}

	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.ListReplays(ctx, namePrefix, state, limit, nextToken)
	if err != nil {
		return nil, err
	}

	replays := make([]map[string]interface{}, 0, len(result.Replays))
	for _, replay := range result.Replays {
		r := map[string]interface{}{
			"ReplayName":     replay.Name,
			"ReplayArn":      replay.ARN,
			"State":          string(replay.State),
			"EventSourceArn": replay.EventSourceARN,
			"EventStartTime": replay.EventStartTime.Unix(),
			"EventEndTime":   replay.EventEndTime.Unix(),
		}
		if replay.Destination != nil {
			r["Destination"] = map[string]interface{}{
				"Arn":        replay.Destination.Arn,
				"FilterArns": replay.Destination.FilterArns,
			}
		}
		replays = append(replays, r)
	}

	response := map[string]interface{}{
		"Replays": replays,
	}
	if result.NextToken != "" {
		response["NextToken"] = result.NextToken
	}

	return response, nil
}

// CancelReplay cancels a running replay.
func (s *EventsService) CancelReplay(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	replayName := request.GetParamLowerFirst(req.Parameters, "ReplayName")
	if replayName == "" {
		return nil, NewValidationException("ReplayName is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	replay, err := store.GetReplay(ctx, replayName)
	if err != nil {
		if err == eventsstore.ErrReplayNotFound {
			return nil, NewResourceNotFoundException("Replay '" + replayName + "' does not exist")
		}
		return nil, err
	}

	if replay.State != eventsstore.ReplayStateRunning && replay.State != eventsstore.ReplayStateStarting {
		return nil, NewValidationException("Replay cannot be cancelled in state: " + string(replay.State))
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

	return map[string]interface{}{
		"ReplayArn":   replay.ARN,
		"State":       string(replay.State),
		"StateReason": replay.StateReason,
	}, nil
}
