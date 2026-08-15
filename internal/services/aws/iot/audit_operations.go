package iot

import (
	"context"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---- Audit (task/findings) ------------------------------------------
// Audit task lifecycle mirrors the Detect Mitigation task pattern: Start
// persists the task id, Cancel/Describe enforce ResourceNotFoundException for
// unknown ids. Audit findings are not generated without a Defender engine,
// so DescribeAuditFinding always returns NotFound for an arbitrary id.

func (s *IoTService) DescribeAccountAuditConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("config/accountAudit", &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return map[string]interface{}{
			"auditCheckConfigurations":              map[string]interface{}{},
			"auditNotificationTargetConfigurations": map[string]interface{}{},
		}, nil
	}
	return rec, nil
}
func (s *IoTService) UpdateAccountAuditConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{
		"roleArn":                               request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		"auditCheckConfigurations":              request.GetMapParamCaseInsensitive(req.Parameters, "auditCheckConfigurations"),
		"auditNotificationTargetConfigurations": request.GetMapParamCaseInsensitive(req.Parameters, "auditNotificationTargetConfigurations"),
	}
	if err := store.PutGeneric("config/accountAudit", rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DeleteAccountAuditConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteGeneric("config/accountAudit"); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) StartOnDemandAuditTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := uuid.New().String()
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{
		"taskId":         taskId,
		"status":         "IN_PROGRESS",
		"startTime":      time.Now().UTC().Unix(),
		"targetAccounts": request.GetParamCaseInsensitive(req.Parameters, "targetAccounts"),
		"auditChecks":    request.GetParamCaseInsensitive(req.Parameters, "auditChecks"),
	}
	if err := store.PutGeneric("auditTask/"+taskId, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{"taskId": taskId}, nil
}
func (s *IoTService) CancelAuditTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "auditTask/" + taskId
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditTaskNotFound
	}
	rec["status"] = "CANCELED"
	rec["endTime"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeAuditTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("auditTask/"+taskId, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditTaskNotFound
	}
	return map[string]interface{}{
		"taskId":           rec["taskId"],
		"taskStatus":       rec["status"],
		"auditTaskDetails": rec,
	}, nil
}
func (s *IoTService) ListAuditTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("auditTask/")
	if err != nil {
		return nil, err
	}
	tasks := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		tasks = append(tasks, map[string]interface{}{
			"taskId":     rec["taskId"],
			"taskStatus": rec["status"],
			"taskType":   "ON_DEMAND_AUDIT_TASK",
		})
	}
	return paginatedMaps("tasks", tasks, req.Parameters), nil
}
func (s *IoTService) DescribeAuditFinding(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	findingId := request.GetParamCaseInsensitive(req.Parameters, "findingId")
	if findingId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("auditFinding/"+findingId, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		// No Defender engine generates findings, so any caller-supplied id
		// is unknown to the platform. AWS returns ResourceNotFoundException.
		return nil, iotstore.ErrAuditFindingNotFound
	}
	return map[string]interface{}{"finding": rec}, nil
}
func (s *IoTService) ListAuditFindings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("auditFinding/")
	if err != nil {
		return nil, err
	}
	findings := make([]map[string]interface{}, 0, len(items))
	findings = append(findings, items...)
	return paginatedMaps("findings", findings, req.Parameters), nil
}
func (s *IoTService) ListRelatedResourcesForAuditFinding(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	findingID := request.GetParamCaseInsensitive(req.Parameters, "findingId")
	if findingID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("auditFinding/"+findingID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditFindingNotFound
	}
	resources := []map[string]interface{}{}
	if raw, ok := rec["relatedResources"].([]interface{}); ok {
		for _, r := range raw {
			if m, ok := r.(map[string]interface{}); ok {
				resources = append(resources, m)
			}
		}
	}
	return paginatedMaps("relatedResources", resources, req.Parameters), nil
}
