package scheduler

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// namePattern matches the AWS Scheduler Name/GroupName constraint:
// 1-64 chars of alphanumeric, hyphen, underscore, and period.
var namePattern = regexp.MustCompile(`^[0-9a-zA-Z-_.]{1,64}$`)

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

func getListGroupName(params map[string]interface{}) string {
	groupName := request.GetStringParam(params, "GroupName")
	if groupName == "" {
		groupName = request.GetStringParam(params, "groupName")
	}
	if groupName == "" {
		groupName = request.GetStringParam(params, "ScheduleGroup")
	}
	if groupName == "" {
		groupName = "default"
	}
	return groupName
}

// CreateSchedule creates a new schedule in EventBridge Scheduler.
func (s *SchedulerService) CreateSchedule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, ErrValidation
	}
	if !namePattern.MatchString(name) {
		return nil, ErrValidation
	}

	scheduleExpression := request.GetStringParam(req.Parameters, "ScheduleExpression")
	if scheduleExpression == "" {
		return nil, ErrValidation
	}

	if !isValidScheduleExpression(scheduleExpression) {
		return nil, ErrInvalidScheduleExpression
	}

	target, err := parseTarget(req.Parameters)
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, ErrInvalidTarget
	}

	if err := s.validateVpcConfig(ctx, reqCtx, target); err != nil {
		return nil, err
	}

	if target.RoleArn != "" {
		validator := reqCtx.GetIAMValidator()
		if err := validator.ValidateRoleForService(ctx, target.RoleArn, iam.ServicePrincipalScheduler); err != nil {
			return nil, err
		}
	}

	flexibleTimeWindow, err := parseFlexibleTimeWindow(req.Parameters)
	if err != nil {
		return nil, err
	}
	if flexibleTimeWindow == nil {
		flexibleTimeWindow = &schedulerstore.FlexibleTimeWindow{Mode: schedulerstore.FlexibleTimeWindowModeOff}
	}

	groupName := request.GetStringParam(req.Parameters, "GroupName")
	if groupName == "" {
		groupName = "default"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if groupName != "default" {
		if _, err := store.GetScheduleGroup(ctx, groupName); err != nil {
			if err == schedulerstore.ErrScheduleGroupNotFound {
				return nil, ErrScheduleGroupNotFound
			}
			return nil, ErrInternalServer
		}
	}

	schedule := &schedulerstore.Schedule{
		Name:                       name,
		GroupName:                  groupName,
		ScheduleExpression:         scheduleExpression,
		Target:                     target,
		FlexibleTimeWindow:         flexibleTimeWindow,
		State:                      schedulerstore.ScheduleStateEnabled,
		ScheduleExpressionTimezone: request.GetStringParam(req.Parameters, "ScheduleExpressionTimezone"),
		Description:                request.GetStringParam(req.Parameters, "Description"),
		KmsKeyArn:                  request.GetStringParam(req.Parameters, "KmsKeyArn"),
	}

	if startDateStr := request.GetStringParam(req.Parameters, "StartDate"); startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			schedule.StartDate = &t
		} else {
			return nil, ErrInvalidDate
		}
	}
	if endDateStr := request.GetStringParam(req.Parameters, "EndDate"); endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			schedule.EndDate = &t
		} else {
			return nil, ErrInvalidDate
		}
	}

	if stateStr := request.GetStringParam(req.Parameters, "State"); stateStr != "" {
		if stateStr != "ENABLED" && stateStr != "DISABLED" {
			return nil, ErrInvalidState
		}
		schedule.State = schedulerstore.ScheduleState(stateStr)
	}

	if actionStr := request.GetStringParam(req.Parameters, "ActionAfterCompletion"); actionStr != "" {
		if actionStr != "DELETE" && actionStr != "NONE" {
			return nil, ErrInvalidActionAfterCompletion
		}
		schedule.ActionAfterCompletion = schedulerstore.ActionAfterCompletion(actionStr)
	}

	if err := store.CreateSchedule(ctx, schedule); err != nil {
		if err == schedulerstore.ErrScheduleAlreadyExists {
			return nil, ErrScheduleAlreadyExists
		}
		return nil, ErrInternalServer
	}

	return map[string]interface{}{
		"ScheduleArn": schedule.ARN,
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

	if err := store.DeleteSchedule(ctx, groupName, name); err != nil {
		if err == schedulerstore.ErrScheduleNotFound {
			return nil, ErrScheduleNotFound
		}
		return nil, ErrInternalServer
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

	schedule, err := store.GetSchedule(ctx, groupName, name)
	if err != nil {
		if err == schedulerstore.ErrScheduleNotFound {
			return nil, ErrScheduleNotFound
		}
		return nil, ErrInternalServer
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

	existing, err := store.GetSchedule(ctx, groupName, name)
	if err != nil {
		if err == schedulerstore.ErrScheduleNotFound {
			return nil, ErrScheduleNotFound
		}
		return nil, ErrInternalServer
	}

	// UpdateSchedule is a PUT operation (full replacement). All fields from
	// the request replace the existing values, including empty strings which
	// clear the field. Required fields (ScheduleExpression, Target,
	// FlexibleTimeWindow) must be provided and are validated below.
	scheduleExpression := request.GetStringParam(req.Parameters, "ScheduleExpression")
	if scheduleExpression == "" {
		return nil, ErrValidation
	}
	if !isValidScheduleExpression(scheduleExpression) {
		return nil, ErrInvalidScheduleExpression
	}
	existing.ScheduleExpression = scheduleExpression

	target, err := parseTarget(req.Parameters)
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, ErrInvalidTarget
	}
	if err := s.validateVpcConfig(ctx, reqCtx, target); err != nil {
		return nil, err
	}
	if target.RoleArn != "" {
		validator := reqCtx.GetIAMValidator()
		if err := validator.ValidateRoleForService(ctx, target.RoleArn, iam.ServicePrincipalScheduler); err != nil {
			return nil, err
		}
	}
	existing.Target = target

	flexibleTimeWindow, err := parseFlexibleTimeWindow(req.Parameters)
	if err != nil {
		return nil, err
	}
	if flexibleTimeWindow == nil {
		flexibleTimeWindow = &schedulerstore.FlexibleTimeWindow{Mode: schedulerstore.FlexibleTimeWindowModeOff}
	}
	existing.FlexibleTimeWindow = flexibleTimeWindow

	// Optional fields — full PUT replacement. Empty/unset applies the
	// AWS default, not preservation of the previous value.
	existing.Description = request.GetStringParam(req.Parameters, "Description")
	existing.ScheduleExpressionTimezone = request.GetStringParam(req.Parameters, "ScheduleExpressionTimezone")
	existing.KmsKeyArn = request.GetStringParam(req.Parameters, "KmsKeyArn")

	stateStr := request.GetStringParam(req.Parameters, "State")
	if stateStr == "" {
		existing.State = schedulerstore.ScheduleStateEnabled
	} else if stateStr != "ENABLED" && stateStr != "DISABLED" {
		return nil, ErrInvalidState
	} else {
		existing.State = schedulerstore.ScheduleState(stateStr)
	}

	actionStr := request.GetStringParam(req.Parameters, "ActionAfterCompletion")
	if actionStr == "" {
		existing.ActionAfterCompletion = schedulerstore.ActionAfterCompletionNone
	} else if actionStr != "DELETE" && actionStr != "NONE" {
		return nil, ErrInvalidActionAfterCompletion
	} else {
		existing.ActionAfterCompletion = schedulerstore.ActionAfterCompletion(actionStr)
	}
	if startDateStr := request.GetStringParam(req.Parameters, "StartDate"); startDateStr != "" {
		t, parseErr := time.Parse(time.RFC3339, startDateStr)
		if parseErr != nil {
			return nil, ErrInvalidDate
		}
		existing.StartDate = &t
	} else {
		existing.StartDate = nil
	}
	if endDateStr := request.GetStringParam(req.Parameters, "EndDate"); endDateStr != "" {
		t, parseErr := time.Parse(time.RFC3339, endDateStr)
		if parseErr != nil {
			return nil, ErrInvalidDate
		}
		existing.EndDate = &t
	} else {
		existing.EndDate = nil
	}

	if err := store.UpdateSchedule(ctx, existing); err != nil {
		return nil, ErrInternalServer
	}

	return map[string]interface{}{
		"ScheduleArn": existing.ARN,
	}, nil
}

// ListSchedules lists schedules in EventBridge Scheduler.
func (s *SchedulerService) ListSchedules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := getListGroupName(req.Parameters)
	namePrefix := request.GetStringParam(req.Parameters, "NamePrefix")
	stateFilter := schedulerstore.ScheduleState(request.GetStringParam(req.Parameters, "State"))
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")
	maxResults := int32(100)
	if mr := request.GetStringParam(req.Parameters, "MaxResults"); mr != "" {
		if parsed, err := strconv.Atoi(mr); err == nil {
			if parsed < 1 || parsed > 100 {
				return nil, ErrValidation
			}
			maxResults = int32(parsed)
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.ListSchedules(ctx, groupName, namePrefix, stateFilter, maxResults, nextToken)
	if err != nil {
		return nil, ErrInternalServer
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
			item["CreationDate"] = formatEpochSeconds(*sch.CreationDate)
		}
		if sch.LastModificationDate != nil {
			item["LastModificationDate"] = formatEpochSeconds(*sch.LastModificationDate)
		}
		schedules[i] = item
		if sch.Target != nil {
			schedules[i]["Target"] = map[string]interface{}{
				"Arn": sch.Target.Arn,
			}
		}
	}

	response := map[string]interface{}{
		"Schedules": schedules,
	}
	pagination.SetNextToken(response, "NextToken", result.NextToken)

	return response, nil
}
