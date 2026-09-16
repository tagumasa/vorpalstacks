package scheduler

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/utils/timeutils"
)

// getScheduleNameAndGroup extracts the identifier pair for GetSchedule and
// DeleteSchedule: Name is the URI label (the router carries labels under
// the modelled member name) and GroupName is the httpQuery "groupName"
// binding both operations carry in the model. UpdateSchedule differs — its
// GroupName is a body member — and reads its own pair. Required-member
// rejection, the Name/GroupName patterns, and the default-group resolution
// live in the Core (resolveScheduleIdentifier); the shared request helper's
// all-lowercase fallback is architecture-level behaviour, not a handler
// spelling chain.
func getScheduleNameAndGroup(params map[string]interface{}) (name, groupName string) {
	name = request.GetStringParam(params, "Name")
	groupName = request.GetStringParam(params, "groupName")
	return name, groupName
}

// getListGroupName extracts the GroupName filter for ListSchedules, whose
// model binding is the httpQuery "ScheduleGroup". An absent parameter means
// no group filter: per the ListSchedulesInput member documentation, the
// group filter only applies "if specified", so an unfiltered list must
// return schedules from every group (the store treats an empty group name
// as no prefix filter).
func getListGroupName(params map[string]interface{}) string {
	return request.GetStringParam(params, "ScheduleGroup")
}

// CreateSchedule creates a new schedule in EventBridge Scheduler.
func (s *SchedulerService) CreateSchedule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	target, err := parseTarget(req.Parameters)
	if err != nil {
		return nil, err
	}

	flexibleTimeWindow, err := parseFlexibleTimeWindow(req.Parameters)
	if err != nil {
		return nil, err
	}

	spec := &ScheduleSpec{
		Name:                       request.GetStringParam(req.Parameters, "Name"),
		GroupName:                  request.GetStringParam(req.Parameters, "GroupName"),
		ScheduleExpression:         request.GetStringParam(req.Parameters, "ScheduleExpression"),
		ScheduleExpressionTimezone: request.GetStringParam(req.Parameters, "ScheduleExpressionTimezone"),
		Description:                request.GetStringParam(req.Parameters, "Description"),
		State:                      request.GetStringParam(req.Parameters, "State"),
		KmsKeyArn:                  request.GetStringParam(req.Parameters, "KmsKeyArn"),
		StartDate:                  request.GetStringParam(req.Parameters, "StartDate"),
		EndDate:                    request.GetStringParam(req.Parameters, "EndDate"),
		ActionAfterCompletion:      request.GetStringParam(req.Parameters, "ActionAfterCompletion"),
		Target:                     target,
		FlexibleTimeWindow:         flexibleTimeWindow,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createScheduleCore(ctx, store, &CreateScheduleInput{
		Spec:         spec,
		ClientToken:  request.GetStringParam(req.Parameters, "ClientToken"),
		Region:       reqCtx.GetRegion(),
		IAMValidator: reqCtx.GetIAMValidator(),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ScheduleArn": result.ScheduleArn,
	}, nil
}

// DeleteSchedule deletes a schedule from EventBridge Scheduler.
func (s *SchedulerService) DeleteSchedule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name, groupName := getScheduleNameAndGroup(req.Parameters)
	// DeleteScheduleInput.ClientToken is the httpQuery "clientToken"
	// binding the model declares (the Pascal try was a dead key: query
	// keys are case-sensitive).
	clientToken := request.GetStringParam(req.Parameters, "clientToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteScheduleCore(ctx, store, &DeleteScheduleInput{
		Name:        name,
		GroupName:   groupName,
		ClientToken: clientToken,
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetSchedule retrieves a schedule from EventBridge Scheduler.
func (s *SchedulerService) GetSchedule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name, groupName := getScheduleNameAndGroup(req.Parameters)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	schedule, err := s.getScheduleCore(ctx, store, &GetScheduleInput{
		Name:      name,
		GroupName: groupName,
	})
	if err != nil {
		return nil, err
	}

	return scheduleToResponse(schedule), nil
}

// UpdateSchedule updates an existing schedule in EventBridge Scheduler.
func (s *SchedulerService) UpdateSchedule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// UpdateScheduleInput binds Name as the URI label and GroupName as a
	// body member (Pascal) — unlike the Get/Delete httpQuery "groupName".
	name := request.GetStringParam(req.Parameters, "Name")
	groupName := request.GetStringParam(req.Parameters, "GroupName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	target, err := parseTarget(req.Parameters)
	if err != nil {
		return nil, err
	}

	flexibleTimeWindow, err := parseFlexibleTimeWindow(req.Parameters)
	if err != nil {
		return nil, err
	}

	spec := &ScheduleSpec{
		Name:                       name,
		GroupName:                  groupName,
		ScheduleExpression:         request.GetStringParam(req.Parameters, "ScheduleExpression"),
		ScheduleExpressionTimezone: request.GetStringParam(req.Parameters, "ScheduleExpressionTimezone"),
		Description:                request.GetStringParam(req.Parameters, "Description"),
		State:                      request.GetStringParam(req.Parameters, "State"),
		KmsKeyArn:                  request.GetStringParam(req.Parameters, "KmsKeyArn"),
		StartDate:                  request.GetStringParam(req.Parameters, "StartDate"),
		EndDate:                    request.GetStringParam(req.Parameters, "EndDate"),
		ActionAfterCompletion:      request.GetStringParam(req.Parameters, "ActionAfterCompletion"),
		Target:                     target,
		FlexibleTimeWindow:         flexibleTimeWindow,
	}

	result, err := s.updateScheduleCore(ctx, store, &UpdateScheduleInput{
		Spec:         spec,
		Region:       reqCtx.GetRegion(),
		IAMValidator: reqCtx.GetIAMValidator(),
		ClientToken:  request.GetStringParam(req.Parameters, "ClientToken"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ScheduleArn": result.ScheduleArn,
	}, nil
}

// ListSchedules lists schedules in EventBridge Scheduler.
func (s *SchedulerService) ListSchedules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := getListGroupName(req.Parameters)
	namePrefix := request.GetStringParam(req.Parameters, "NamePrefix")
	stateFilter := request.GetStringParam(req.Parameters, "State")
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")
	maxResults, err := parseMaxResultsParam(req.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listSchedulesCore(ctx, store, &ListSchedulesInput{
		GroupName:  groupName,
		NamePrefix: namePrefix,
		State:      stateFilter,
		MaxResults: maxResults,
		NextToken:  nextToken,
	})
	if err != nil {
		return nil, err
	}

	schedules := make([]map[string]interface{}, len(result.Schedules))
	for i, sch := range result.Schedules {
		item := map[string]interface{}{
			"Arn":       sch.Arn,
			"Name":      sch.Name,
			"GroupName": sch.GroupName,
			"State":     string(sch.State),
		}
		if sch.CreationDate != nil {
			item["CreationDate"] = timeutils.FormatEpochSeconds(*sch.CreationDate)
		}
		if sch.LastModificationDate != nil {
			item["LastModificationDate"] = timeutils.FormatEpochSeconds(*sch.LastModificationDate)
		}
		schedules[i] = item
		if sch.Target != nil {
			schedules[i]["Target"] = map[string]interface{}{
				"Arn": sch.Target.Arn,
			}
		}
	}

	resp := map[string]interface{}{
		"Schedules": schedules,
	}
	pagination.SetNextToken(resp, "NextToken", result.NextToken)

	return resp, nil
}
