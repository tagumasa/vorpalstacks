package scheduler

import (
	"context"
	"strconv"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/utils/timeutils"
)

func getScheduleNameAndGroup(params map[string]interface{}) (name, groupName string, err error) {
	name = request.GetStringParam(params, "Name")
	if name == "" {
		name = request.GetStringParam(params, "name")
	}
	if name == "" {
		return "", "", ErrValidation
	}
	groupName = request.GetStringParam(params, "GroupName")
	if groupName == "" {
		groupName = request.GetStringParam(params, "groupName")
	}
	if groupName == "" {
		groupName = "default"
	}
	return name, groupName, nil
}

// getListGroupName extracts the GroupName filter for ListSchedules. An
// absent parameter means no group filter: per the ListSchedulesInput
// member documentation, the group filter only applies "if specified", so
// an unfiltered list must return schedules from every group (the store
// treats an empty group name as no prefix filter).
func getListGroupName(params map[string]interface{}) string {
	groupName := request.GetStringParam(params, "GroupName")
	if groupName == "" {
		groupName = request.GetStringParam(params, "groupName")
	}
	if groupName == "" {
		groupName = request.GetStringParam(params, "ScheduleGroup")
	}
	return groupName
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
	name, groupName, err := getScheduleNameAndGroup(req.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteScheduleCore(ctx, store, &DeleteScheduleInput{
		Name:      name,
		GroupName: groupName,
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetSchedule retrieves a schedule from EventBridge Scheduler.
func (s *SchedulerService) GetSchedule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name, groupName, err := getScheduleNameAndGroup(req.Parameters)
	if err != nil {
		return nil, err
	}

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
	name, groupName, err := getScheduleNameAndGroup(req.Parameters)
	if err != nil {
		return nil, err
	}

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
	maxResults := int32(DefaultListMaxResults)
	if mr := request.GetStringParam(req.Parameters, "MaxResults"); mr != "" {
		parsed, err := strconv.Atoi(mr)
		if err != nil {
			return nil, ErrValidation
		}
		if parsed < 1 || parsed > MaxListMaxResults {
			return nil, ErrValidation
		}
		maxResults = int32(parsed)
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
