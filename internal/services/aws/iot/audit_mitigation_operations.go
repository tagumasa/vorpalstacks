package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

func (s *IoTService) StartAuditMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	taskId, err := s.startAuditMitigationActionsTaskCore(store, StartAuditMitigationActionsTaskInput{
		TaskID:                     request.GetParamCaseInsensitive(req.Parameters, "taskId"),
		Target:                     request.GetMapParamCaseInsensitive(req.Parameters, "target"),
		AuditCheckToActionsMapping: request.GetMapParamCaseInsensitive(req.Parameters, "auditCheckToActionsMapping"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"taskId": taskId}, nil
}
func (s *IoTService) CancelAuditMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.cancelAuditMitigationActionsTaskCore(store, request.GetParamCaseInsensitive(req.Parameters, "taskId")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeAuditMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec, err := s.describeAuditMitigationActionsTaskCore(store, request.GetParamCaseInsensitive(req.Parameters, "taskId"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"taskStatus":                 rec.TaskStatus,
		"startTime":                  rec.StartTime,
		"endTime":                    rec.EndTime,
		"target":                     rec.Target,
		"auditCheckToActionsMapping": rec.AuditCheckToActionsMapping,
	}, nil
}
func (s *IoTService) ListAuditMitigationActionsExecutions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return paginatedMaps("taskExecutions", []map[string]interface{}{}, req.Parameters)
}
func (s *IoTService) ListAuditMitigationActionsTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := s.listAuditMitigationActionsTasksCore(store)
	if err != nil {
		return nil, err
	}
	tasks := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		tasks = append(tasks, map[string]interface{}{
			"taskId":     item.TaskID,
			"startTime":  item.StartTime,
			"taskStatus": item.TaskStatus,
		})
	}
	return paginatedMaps("tasks", tasks, req.Parameters)
}
