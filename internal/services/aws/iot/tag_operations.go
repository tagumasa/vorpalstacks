package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
)

// TagResource adds or overwrites tags on an IoT resource.
func (s *IoTService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, s.iotTagConfig(store))
}

// UntagResource removes tags from an IoT resource.
func (s *IoTService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, s.iotTagConfig(store))
}

// ListTagsForResource lists tags attached to an IoT resource.
func (s *IoTService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, s.iotTagConfig(store))
}

func (s *IoTService) ListActiveViolations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	var listSuppressed *bool
	if request.HasParam(req.Parameters, "listSuppressedAlerts") {
		v := request.GetBoolParam(req.Parameters, "listSuppressedAlerts")
		listSuppressed = &v
	}
	violations, err := s.listActiveViolationsCore(store, ListActiveViolationsInput{
		ThingName:            request.GetParamCaseInsensitive(req.Parameters, "thingName"),
		SecurityProfileName:  request.GetParamCaseInsensitive(req.Parameters, "securityProfileName"),
		BehaviorCriteriaType: request.GetParamCaseInsensitive(req.Parameters, "behaviorCriteriaType"),
		VerificationState:    request.GetParamCaseInsensitive(req.Parameters, "verificationState"),
		ListSuppressedAlerts: listSuppressed,
	})
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0, len(violations))
	for _, v := range violations {
		items = append(items, activeViolationResponse(v))
	}
	return paginatedMaps("activeViolations", items, req.Parameters)
}

func (s *IoTService) ListViolationEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	var listSuppressed *bool
	if request.HasParam(req.Parameters, "listSuppressedAlerts") {
		v := request.GetBoolParam(req.Parameters, "listSuppressedAlerts")
		listSuppressed = &v
	}
	events, err := s.listViolationEventsCore(store, parseListOptions(req.Parameters), ListViolationEventsInput{
		StartTime:            timestampMemberParam(req.Parameters, "startTime"),
		EndTime:              timestampMemberParam(req.Parameters, "endTime"),
		StartTimeProvided:    request.HasParam(req.Parameters, "startTime"),
		EndTimeProvided:      request.HasParam(req.Parameters, "endTime"),
		SecurityProfileName:  request.GetParamCaseInsensitive(req.Parameters, "securityProfileName"),
		ThingName:            request.GetParamCaseInsensitive(req.Parameters, "thingName"),
		BehaviorCriteriaType: request.GetParamCaseInsensitive(req.Parameters, "behaviorCriteriaType"),
		VerificationState:    request.GetParamCaseInsensitive(req.Parameters, "verificationState"),
		ListSuppressedAlerts: listSuppressed,
	})
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0, len(events))
	for _, e := range events {
		items = append(items, violationEventResponse(e))
	}
	return map[string]interface{}{
		"violationEvents": items,
	}, nil
}

// No ML behaviour-model training pipeline exists in this platform; an empty
// list matches the AWS wire shape.
func (s *IoTService) GetBehaviorModelTrainingSummaries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"summaries": []map[string]interface{}{},
	}, nil
}
