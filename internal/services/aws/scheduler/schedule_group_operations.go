package scheduler

import (
	"context"
	"strconv"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateScheduleGroup creates a new schedule group in EventBridge Scheduler.
func (s *SchedulerService) CreateScheduleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.createScheduleGroupCore(ctx, reqCtx, &CreateScheduleGroupInput{
		Name:        request.GetStringParam(req.Parameters, "Name"),
		Tags:        tags.ParseTags(req.Parameters, "Tags"),
		ClientToken: request.GetStringParam(req.Parameters, "ClientToken"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"ScheduleGroupArn": result.ScheduleGroupArn,
	}, nil
}

// DeleteScheduleGroup deletes a schedule group from EventBridge Scheduler.
func (s *SchedulerService) DeleteScheduleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.deleteScheduleGroupCore(ctx, reqCtx, &DeleteScheduleGroupInput{
		Name: request.GetStringParam(req.Parameters, "Name"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// GetScheduleGroup retrieves a schedule group from EventBridge Scheduler.
func (s *SchedulerService) GetScheduleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	group, err := s.getScheduleGroupCore(ctx, reqCtx, &GetScheduleGroupInput{
		Name: request.GetStringParam(req.Parameters, "Name"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"Arn":                  group.Arn,
		"Name":                 group.Name,
		"State":                group.State,
		"CreationDate":         timeutils.FormatEpochSeconds(group.CreationDate),
		"LastModificationDate": timeutils.FormatEpochSeconds(group.LastModificationDate),
	}, nil
}

// ListScheduleGroups lists schedule groups in EventBridge Scheduler.
func (s *SchedulerService) ListScheduleGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var maxResults int32
	if mr := request.GetStringParam(req.Parameters, "MaxResults"); mr != "" {
		parsed, err := strconv.Atoi(mr)
		if err != nil {
			return nil, ErrValidation
		}
		maxResults = int32(parsed)
	}

	result, err := s.listScheduleGroupsCore(ctx, reqCtx, &ListScheduleGroupsInput{
		NamePrefix: request.GetStringParam(req.Parameters, "NamePrefix"),
		MaxResults: maxResults,
		NextToken:  pagination.GetMarker(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	groups := make([]map[string]interface{}, len(result.ScheduleGroups))
	for i, g := range result.ScheduleGroups {
		item := map[string]interface{}{
			"Arn":   g.Arn,
			"Name":  g.Name,
			"State": g.State,
		}
		if g.CreationDate != nil {
			item["CreationDate"] = timeutils.FormatEpochSeconds(*g.CreationDate)
		}
		if g.LastModificationDate != nil {
			item["LastModificationDate"] = timeutils.FormatEpochSeconds(*g.LastModificationDate)
		}
		groups[i] = item
	}

	resp := map[string]interface{}{
		"ScheduleGroups": groups,
	}
	pagination.SetNextToken(resp, "NextToken", result.NextToken)

	return resp, nil
}
