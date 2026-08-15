package iot

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func (s *IoTService) StartAuditMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{
		"taskId":    taskId,
		"status":    "IN_PROGRESS",
		"startTime": time.Now().UTC().Unix(),
	}
	if err := store.PutGeneric("auditMitigationTask/"+taskId, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) CancelAuditMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "auditMitigationTask/" + taskId
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditMitigationTaskNotFound
	}
	// Cancellation transitions the task to CANCELED (the
	// AuditMitigationActionsTaskStatus enum value); the record stays
	// queryable via DescribeAuditMitigationActionsTask.
	rec["status"] = "CANCELED"
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeAuditMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("auditMitigationTask/"+taskId, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditMitigationTaskNotFound
	}
	return rec, nil
}
func (s *IoTService) ListAuditMitigationActionsExecutions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return paginatedMaps("taskExecutions", []map[string]interface{}{}, req.Parameters), nil
}
func (s *IoTService) ListAuditMitigationActionsTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("auditMitigationTask/")
	if err != nil {
		return nil, err
	}
	tasks := make([]map[string]interface{}, 0, len(items))
	tasks = append(tasks, items...)
	return paginatedMaps("tasks", tasks, req.Parameters), nil
}
