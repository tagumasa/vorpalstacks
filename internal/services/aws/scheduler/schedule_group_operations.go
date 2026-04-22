package scheduler

import (
	"context"
	"strconv"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// CreateScheduleGroup creates a new schedule group in EventBridge Scheduler.
func (s *SchedulerService) CreateScheduleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, ErrValidation
	}

	group := &schedulerstore.ScheduleGroup{
		Name: name,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.CreateScheduleGroup(ctx, group); err != nil {
		if err == schedulerstore.ErrScheduleGroupAlreadyExists {
			return nil, ErrScheduleGroupAlreadyExists
		}
		logs.Debug("Failed to create schedule group", logs.String("name", name), logs.String("error", err.Error()))
		return nil, ErrInternalServer
	}

	parsedTags := tags.ParseTags(req.Parameters, "Tags")
	if len(parsedTags) > 0 {
		arn := group.ARN
		if err := store.TagFromSlice(arn, parsedTags); err != nil {
			logs.Debug("Failed to tag schedule group", logs.String("arn", arn), logs.String("error", err.Error()))
		}
	}

	return map[string]interface{}{
		"ScheduleGroupArn": group.ARN,
	}, nil
}

// DeleteScheduleGroup deletes a schedule group from EventBridge Scheduler.
func (s *SchedulerService) DeleteScheduleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, ErrValidation
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteScheduleGroup(ctx, name); err != nil {
		if err == schedulerstore.ErrScheduleGroupNotFound {
			return nil, ErrScheduleGroupNotFound
		}
		if err == schedulerstore.ErrScheduleGroupNotEmpty {
			return nil, ErrScheduleGroupNotEmpty
		}
		logs.Debug("Failed to delete schedule group", logs.String("name", name), logs.String("error", err.Error()))
		return nil, ErrInternalServer
	}

	return response.EmptyResponse(), nil
}

// GetScheduleGroup retrieves a schedule group from EventBridge Scheduler.
func (s *SchedulerService) GetScheduleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, ErrValidation
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	group, err := store.GetScheduleGroup(ctx, name)
	if err != nil {
		if err == schedulerstore.ErrScheduleGroupNotFound {
			return nil, ErrScheduleGroupNotFound
		}
		logs.Debug("Failed to get schedule group", logs.String("name", name), logs.String("error", err.Error()))
		return nil, ErrInternalServer
	}

	return map[string]interface{}{
		"Arn":                  group.ARN,
		"Name":                 group.Name,
		"State":                string(group.State),
		"CreationDate":         formatEpochSeconds(group.CreationDate),
		"LastModificationDate": formatEpochSeconds(group.LastModificationDate),
	}, nil
}

// ListScheduleGroups lists schedule groups in EventBridge Scheduler.
func (s *SchedulerService) ListScheduleGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namePrefix := request.GetStringParam(req.Parameters, "NamePrefix")
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

	result, err := store.ListScheduleGroups(ctx, namePrefix, maxResults, nextToken)
	if err != nil {
		logs.Debug("Failed to list schedule groups", logs.String("error", err.Error()))
		return nil, ErrInternalServer
	}

	groups := make([]map[string]interface{}, len(result.ScheduleGroups))
	for i, g := range result.ScheduleGroups {
		item := map[string]interface{}{
			"Arn":   g.Arn,
			"Name":  g.Name,
			"State": string(g.State),
		}
		if g.CreationDate != nil {
			item["CreationDate"] = formatEpochSeconds(*g.CreationDate)
		}
		if g.LastModificationDate != nil {
			item["LastModificationDate"] = formatEpochSeconds(*g.LastModificationDate)
		}
		groups[i] = item
	}

	response := map[string]interface{}{
		"ScheduleGroups": groups,
	}
	pagination.SetNextToken(response, "NextToken", result.NextToken)

	return response, nil
}
