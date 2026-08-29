package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// ---- Audit (task/findings) ------------------------------------------

func (s *IoTService) DescribeAccountAuditConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.describeAccountAuditConfigurationCore(store)
}
func (s *IoTService) UpdateAccountAuditConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := UpdateAccountAuditConfigurationInput{
		RoleArn:                               request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		AuditCheckConfigurations:              request.GetMapParamCaseInsensitive(req.Parameters, "auditCheckConfigurations"),
		AuditNotificationTargetConfigurations: request.GetMapParamCaseInsensitive(req.Parameters, "auditNotificationTargetConfigurations"),
	}
	if err := s.updateAccountAuditConfigurationCore(store, in); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DeleteAccountAuditConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteAccountAuditConfigurationCore(store,
		request.GetBoolParam(req.Parameters, "deleteScheduledAudits")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) StartOnDemandAuditTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	taskId, err := s.startOnDemandAuditTaskCore(store, StartOnDemandAuditTaskInput{
		TargetCheckNames: request.GetStringList(req.Parameters, "targetCheckNames"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"taskId": taskId}, nil
}
func (s *IoTService) CancelAuditTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.cancelAuditTaskCore(store, request.GetParamCaseInsensitive(req.Parameters, "taskId")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeAuditTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	details, err := s.describeAuditTaskCore(store, request.GetParamCaseInsensitive(req.Parameters, "taskId"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"taskStatus":    details.TaskStatus,
		"taskType":      details.TaskType,
		"taskStartTime": details.TaskStartTime,
		"auditDetails":  details.AuditDetails,
	}, nil
}
func (s *IoTService) ListAuditTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := s.listAuditTasksCore(store, ListAuditTasksInput{
		StartTime:         timestampMemberParam(req.Parameters, "startTime"),
		EndTime:           timestampMemberParam(req.Parameters, "endTime"),
		StartTimeProvided: request.HasParam(req.Parameters, "startTime"),
		EndTimeProvided:   request.HasParam(req.Parameters, "endTime"),
		TaskType:          request.GetParamCaseInsensitive(req.Parameters, "taskType"),
		TaskStatus:        request.GetParamCaseInsensitive(req.Parameters, "taskStatus"),
	})
	if err != nil {
		return nil, err
	}
	tasks := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		tasks = append(tasks, map[string]interface{}{
			"taskId":     item.TaskID,
			"taskStatus": item.TaskStatus,
			"taskType":   item.TaskType,
		})
	}
	return paginatedMaps("tasks", tasks, req.Parameters)
}
func (s *IoTService) DescribeAuditFinding(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec, err := s.describeAuditFindingCore(store, request.GetParamCaseInsensitive(req.Parameters, "findingId"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"finding": rec}, nil
}
func (s *IoTService) ListAuditFindings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	var listSuppressed *bool
	if request.HasParam(req.Parameters, "listSuppressedFindings") {
		v := request.GetBoolParam(req.Parameters, "listSuppressedFindings")
		listSuppressed = &v
	}
	findings, err := s.listAuditFindingsCore(store, ListAuditFindingsInput{
		TaskID:             request.GetParamCaseInsensitive(req.Parameters, "taskId"),
		CheckName:          request.GetParamCaseInsensitive(req.Parameters, "checkName"),
		ResourceIdentifier: resourceIdentifierFromParams(req.Parameters),
		StartTime:          timestampMemberParam(req.Parameters, "startTime"),
		EndTime:            timestampMemberParam(req.Parameters, "endTime"),
		StartTimeProvided:  request.HasParam(req.Parameters, "startTime"),
		EndTimeProvided:    request.HasParam(req.Parameters, "endTime"),
		ListSuppressed:     listSuppressed,
	})
	if err != nil {
		return nil, err
	}
	return paginatedMaps("findings", findings, req.Parameters)
}
func (s *IoTService) ListRelatedResourcesForAuditFinding(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	resources, err := s.listRelatedResourcesForAuditFindingCore(store, request.GetParamCaseInsensitive(req.Parameters, "findingId"))
	if err != nil {
		return nil, err
	}
	return paginatedMaps("relatedResources", resources, req.Parameters)
}
