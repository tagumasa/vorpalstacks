package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
)

// ---- Mitigation Actions --------------------------------------------

func (s *IoTService) CreateMitigationAction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tagList := tags.ParseTagsWithQueryFallback(req.Parameters, "tags")
	recTags := make(map[string]string, len(tagList))
	for _, t := range tagList {
		recTags[t.Key] = t.Value
	}
	result, err := s.createMitigationActionCore(store, CreateMitigationActionInput{
		ActionName:   request.GetParamCaseInsensitive(req.Parameters, "actionName"),
		RoleArn:      request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		ActionParams: request.GetMapParamCaseInsensitive(req.Parameters, "actionParams"),
		Tags:         recTags,
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"actionArn": result.ActionArn,
		"actionId":  result.ActionID,
	}, nil
}

func (s *IoTService) DeleteMitigationAction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteMitigationActionCore(store, request.GetParamCaseInsensitive(req.Parameters, "actionName")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeMitigationAction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec, err := s.describeMitigationActionCore(store, request.GetParamCaseInsensitive(req.Parameters, "actionName"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"actionName":       rec.Rec["name"],
		"actionType":       rec.Rec["actionType"],
		"actionParams":     rec.Rec["actionParams"],
		"roleArn":          rec.Rec["roleArn"],
		"actionArn":        rec.Arn,
		"actionId":         rec.Rec["actionId"],
		"creationDate":     rec.Rec["creationDate"],
		"lastModifiedDate": rec.Rec["lastModifiedDate"],
	}, nil
}
func (s *IoTService) ListMitigationActions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := s.listMitigationActionsCore(store)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		out = append(out, map[string]interface{}{
			"actionName":   item.Name,
			"actionArn":    item.Arn,
			"creationDate": item.CreationDate,
		})
	}
	return paginatedMaps("actionIdentifiers", out, req.Parameters)
}
func (s *IoTService) UpdateMitigationAction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.updateMitigationActionCore(store, UpdateMitigationActionInput{
		ActionName:      request.GetParamCaseInsensitive(req.Parameters, "actionName"),
		RoleArn:         request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		RoleArnProvided: hasParam(req.Parameters, "roleArn"),
		ActionParams:    request.GetMapParamCaseInsensitive(req.Parameters, "actionParams"),
	})
	if err != nil {
		return nil, err
	}
	// UpdateMitigationActionResponse carries only actionArn and actionId.
	return map[string]interface{}{
		"actionArn": result.Arn,
		"actionId":  result.Rec["actionId"],
	}, nil
}

// ---- Detect Mitigation Actions Tasks --------------------------------
// These handlers persist task records so that Cancel/Describe can resolve the
// identifier and return ResourceNotFoundException for unknown task ids,
// matching the Smithy error trait set on each operation.

func (s *IoTService) StartDetectMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	taskId, err := s.startDetectMitigationActionsTaskCore(store, StartDetectMitigationActionsTaskInput{
		TaskID:  request.GetParamCaseInsensitive(req.Parameters, "taskId"),
		Target:  request.GetMapParamCaseInsensitive(req.Parameters, "target"),
		Actions: request.GetStringList(req.Parameters, "actions"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"taskId": taskId}, nil
}
func (s *IoTService) CancelDetectMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.cancelDetectMitigationActionsTaskCore(store, request.GetParamCaseInsensitive(req.Parameters, "taskId")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeDetectMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec, err := s.describeDetectMitigationActionsTaskCore(store, request.GetParamCaseInsensitive(req.Parameters, "taskId"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"taskSummary": map[string]interface{}{
			"taskId":        rec.TaskID,
			"taskStatus":    rec.TaskStatus,
			"taskStartTime": rec.TaskStartTime,
			"taskEndTime":   rec.TaskEndTime,
			"target":        rec.Target,
		},
	}, nil
}
func (s *IoTService) ListDetectMitigationActionsExecutions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return paginatedMaps("taskExecutions", []map[string]interface{}{}, req.Parameters)
}
func (s *IoTService) ListDetectMitigationActionsTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := s.listDetectMitigationActionsTasksCore(store)
	if err != nil {
		return nil, err
	}
	tasks := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		tasks = append(tasks, map[string]interface{}{
			"taskId":        item.TaskID,
			"taskStatus":    item.TaskStatus,
			"taskStartTime": item.TaskStartTime,
		})
	}
	return paginatedMaps("tasks", tasks, req.Parameters)
}
