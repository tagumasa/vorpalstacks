package scheduler

import (
	"context"
	"regexp"
	"strconv"

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
	// Parse target and flexible time window from request parameters.
	target, err := parseTarget(req.Parameters)
	if err != nil {
		return nil, err
	}

	flexibleTimeWindow, err := parseFlexibleTimeWindow(req.Parameters)
	if err != nil {
		return nil, err
	}
	if flexibleTimeWindow == nil {
		flexibleTimeWindow = &schedulerstore.FlexibleTimeWindow{Mode: schedulerstore.FlexibleTimeWindowModeOff}
	}

	// Build the common spec and run full validation (H1 shared layer).
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

	validated, err := ValidateScheduleFields(spec)
	if err != nil {
		return nil, err
	}

	if target.RoleArn != "" {
		validator := reqCtx.GetIAMValidator()
		if err := validator.ValidateRoleForService(ctx, target.RoleArn, iam.ServicePrincipalScheduler); err != nil {
			return nil, err
		}
	}

	if err := s.validateVpcConfig(ctx, reqCtx.GetRegion(), target); err != nil {
		return nil, err
	}

	groupName := spec.GroupName
	if groupName == "" {
		groupName = "default"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// ClientToken idempotency: if the same token was used within the TTL
	// window, return the previously created resource's ARN (M5).
	clientToken := request.GetStringParam(req.Parameters, "ClientToken")
	tokenClaimed := false
	if clientToken != "" {
		if err := validateClientToken(clientToken); err != nil {
			return nil, err
		}
		expectedArn := store.BuildScheduleARN(groupName, spec.Name)
		if entry, created := store.ClientTokens().LookupOrClaim(clientToken, expectedArn, "schedule"); !created {
			return map[string]interface{}{
				"ScheduleArn": entry.ResourceArn,
			}, nil
		}
		tokenClaimed = true
	}

	if groupName != "default" {
		if _, err := store.GetScheduleGroup(ctx, groupName); err != nil {
			if tokenClaimed {
				store.ClientTokens().Release(clientToken)
			}
			if err == schedulerstore.ErrScheduleGroupNotFound {
				return nil, ErrScheduleGroupNotFound
			}
			return nil, ErrInternalServer
		}
	}

	schedule := &schedulerstore.Schedule{
		Name:                       spec.Name,
		GroupName:                  groupName,
		ScheduleExpression:         spec.ScheduleExpression,
		Target:                     target,
		FlexibleTimeWindow:         flexibleTimeWindow,
		State:                      validated.State,
		ScheduleExpressionTimezone: spec.ScheduleExpressionTimezone,
		Description:                spec.Description,
		KmsKeyArn:                  spec.KmsKeyArn,
		StartDate:                  validated.StartDate,
		EndDate:                    validated.EndDate,
		ActionAfterCompletion:      validated.ActionAfterCompletion,
	}

	if err := store.CreateSchedule(ctx, schedule); err != nil {
		if tokenClaimed {
			store.ClientTokens().Release(clientToken)
		}
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

	// Parse target and flexible time window from request parameters.
	target, err := parseTarget(req.Parameters)
	if err != nil {
		return nil, err
	}

	flexibleTimeWindow, err := parseFlexibleTimeWindow(req.Parameters)
	if err != nil {
		return nil, err
	}
	if flexibleTimeWindow == nil {
		flexibleTimeWindow = &schedulerstore.FlexibleTimeWindow{Mode: schedulerstore.FlexibleTimeWindowModeOff}
	}

	// Build the common spec and run full validation (H1 shared layer).
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

	validated, err := ValidateScheduleFields(spec)
	if err != nil {
		return nil, err
	}

	if target.RoleArn != "" {
		validator := reqCtx.GetIAMValidator()
		if err := validator.ValidateRoleForService(ctx, target.RoleArn, iam.ServicePrincipalScheduler); err != nil {
			return nil, err
		}
	}

	if err := s.validateVpcConfig(ctx, reqCtx.GetRegion(), target); err != nil {
		return nil, err
	}

	// Apply validated fields to the existing schedule.
	existing.ScheduleExpression = spec.ScheduleExpression
	existing.Target = target
	existing.FlexibleTimeWindow = flexibleTimeWindow
	existing.Description = spec.Description
	existing.ScheduleExpressionTimezone = spec.ScheduleExpressionTimezone
	existing.KmsKeyArn = spec.KmsKeyArn
	existing.State = validated.State
	existing.ActionAfterCompletion = validated.ActionAfterCompletion
	existing.StartDate = validated.StartDate
	existing.EndDate = validated.EndDate

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
		parsed, err := strconv.Atoi(mr)
		if err != nil {
			return nil, ErrValidation
		}
		if parsed < 1 || parsed > 100 {
			return nil, ErrValidation
		}
		maxResults = int32(parsed)
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
