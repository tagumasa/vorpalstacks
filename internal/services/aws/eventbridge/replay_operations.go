package eventbridge

import (
	"context"

	"vorpalstacks/internal/common/request"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// parseStartReplayInput reads the StartReplay wire request into the
// transport-agnostic Core input. The Destination and event time bounds stay
// raw because their validation runs at the Core layer after the archive
// existence check.
func parseStartReplayInput(req *request.ParsedRequest) StartReplayInput {
	input := StartReplayInput{
		ReplayName:     request.GetParamLowerFirst(req.Parameters, "ReplayName"),
		EventSourceArn: request.GetParamLowerFirst(req.Parameters, "EventSourceArn"),
		Destination:    req.Parameters["Destination"],
	}
	if startTimeVal, ok := req.Parameters["EventStartTime"]; ok {
		input.EventStartTime = startTimeVal
	}
	if endTimeVal, ok := req.Parameters["EventEndTime"]; ok {
		input.EventEndTime = endTimeVal
	}
	if desc, ok := req.Parameters["Description"].(string); ok {
		input.DescriptionSet = true
		input.Description = desc
	}
	return input
}

// StartReplay starts an event replay from an archive.
func (s *EventsService) StartReplay(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := parseStartReplayInput(req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	replay, err := s.startReplayCore(ctx, store, reqCtx.GetRegion(), input)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"ReplayArn": replay.ARN,
		"State":     string(replay.State),
	}
	if !replay.ReplayStartTime.IsZero() {
		result["ReplayStartTime"] = replay.ReplayStartTime.Unix()
	}
	if replay.StateReason != "" {
		result["StateReason"] = replay.StateReason
	}
	return result, nil
}

// DescribeReplay returns information about a replay.
func (s *EventsService) DescribeReplay(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	replayName := request.GetParamLowerFirst(req.Parameters, "ReplayName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	replay, err := s.describeReplayCore(ctx, store, replayName)
	if err != nil {
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
	input := ListReplaysInput{
		NamePrefix:     request.GetParamLowerFirst(req.Parameters, "NamePrefix"),
		EventSourceArn: request.GetParamLowerFirst(req.Parameters, "EventSourceArn"),
		Limit:          int32(request.GetIntParam(req.Parameters, "Limit")),
		NextToken:      request.GetParamLowerFirst(req.Parameters, "NextToken"),
	}
	if stateStr := request.GetParamLowerFirst(req.Parameters, "State"); stateStr != "" {
		input.State = eventsstore.ReplayState(stateStr)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listReplaysCore(ctx, store, input)
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	replay, err := s.cancelReplayCore(ctx, store, replayName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ReplayArn":   replay.ARN,
		"State":       string(replay.State),
		"StateReason": replay.StateReason,
	}, nil
}
