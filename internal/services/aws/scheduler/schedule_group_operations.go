package scheduler

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateScheduleGroup creates a new schedule group in EventBridge Scheduler.
func (s *SchedulerService) CreateScheduleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.createScheduleGroupCore(ctx, store, &CreateScheduleGroupInput{
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteScheduleGroupCore(ctx, store, &DeleteScheduleGroupInput{
		Name: request.GetStringParam(req.Parameters, "Name"),
		// DeleteScheduleGroupInput.ClientToken is the httpQuery
		// "clientToken" binding the model declares.
		ClientToken: request.GetStringParam(req.Parameters, "clientToken"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// GetScheduleGroup retrieves a schedule group from EventBridge Scheduler.
func (s *SchedulerService) GetScheduleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	group, err := s.getScheduleGroupCore(ctx, store, &GetScheduleGroupInput{
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
	maxResults, err := parseMaxResultsParam(req.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listScheduleGroupsCore(ctx, store, &ListScheduleGroupsInput{
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
			"State": string(g.State),
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
