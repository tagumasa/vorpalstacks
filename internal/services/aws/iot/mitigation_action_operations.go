package iot

import (
	"context"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---- Mitigation Actions --------------------------------------------

func (s *IoTService) CreateMitigationAction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := request.GetMapParamCaseInsensitive(req.Parameters, "actionParams")
	tagList := tags.ParseTagsWithQueryFallback(req.Parameters, "tags")
	recTags := make(map[string]string, len(tagList))
	for _, t := range tagList {
		recTags[t.Key] = t.Value
	}
	rec, err := s.bulkCreate(reqCtx, "mitigationAction", req, "actionName", map[string]interface{}{
		"actionType":   deriveMitigationActionType(params),
		"actionParams": params,
		"roleArn":      request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		"tags":         recTags,
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"actionArn": iotstore.BuildMitigationActionARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"actionId":  uuid.New().String(),
	}, nil
}

// deriveMitigationActionType infers the action type from the actionParams keys.
// AWS derives the type automatically based on which params member is set.
func deriveMitigationActionType(params map[string]interface{}) string {
	if _, ok := params["addThingsToThingGroupParams"]; ok {
		return "ADD_THINGS_TO_THING_GROUP"
	}
	if _, ok := params["enableIoTLoggingParams"]; ok {
		return "ENABLE_IOT_LOGGING"
	}
	if _, ok := params["publishFindingToSnsParams"]; ok {
		return "PUBLISH_FINDING_TO_SNS"
	}
	if _, ok := params["addThingsToCertPoolParams"]; ok {
		return "ADD_THINGS_TO_CERTIFICATE_POOL"
	}
	if _, ok := params["replaceCACertificateParams"]; ok {
		return "REPLACE_CA_CERTIFICATE"
	}
	if _, ok := params["updateDeviceCertificateParams"]; ok {
		return "UPDATE_DEVICE_CERTIFICATE"
	}
	return ""
}
func (s *IoTService) DeleteMitigationAction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.bulkDelete(reqCtx, "mitigationAction", req, "actionName"); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeMitigationAction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, _, exists, err := s.bulkGet(reqCtx, "mitigationAction", req, "actionName")
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrMitigationActionNotFound
	}
	return map[string]interface{}{
		"actionName":       rec["name"],
		"actionType":       rec["actionType"],
		"actionParams":     rec["actionParams"],
		"roleArn":          rec["roleArn"],
		"actionArn":        iotstore.BuildMitigationActionARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"creationDate":     rec["creationDate"],
		"lastModifiedDate": rec["lastModifiedDate"],
	}, nil
}
func (s *IoTService) ListMitigationActions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	items, err := s.bulkList(reqCtx, "mitigationAction")
	if err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		name, _ := item["name"].(string)
		out = append(out, map[string]interface{}{
			"actionName":   name,
			"actionArn":    iotstore.BuildMitigationActionARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
			"creationDate": item["creationDate"],
		})
	}
	return paginatedMaps("actionIdentifiers", out, req.Parameters), nil
}
func (s *IoTService) UpdateMitigationAction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, exists, err := s.bulkUpdate(reqCtx, "mitigationAction", req, "actionName", map[string]interface{}{
		"roleArn":      request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		"actionParams": request.GetMapParamCaseInsensitive(req.Parameters, "actionParams"),
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrMitigationActionNotFound
	}
	return map[string]interface{}{
		"actionName":       rec["name"],
		"actionType":       rec["actionType"],
		"actionParams":     rec["actionParams"],
		"roleArn":          rec["roleArn"],
		"actionArn":        iotstore.BuildMitigationActionARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"lastModifiedDate": rec["lastModifiedDate"],
	}, nil
}

// ---- Detect Mitigation Actions Tasks --------------------------------
// These handlers persist task records so that Cancel/Describe can resolve the
// identifier and return ResourceNotFoundException for unknown task ids,
// matching the Smithy error trait set on each operation.

func (s *IoTService) StartDetectMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskId == "" {
		taskId = uuid.New().String()
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{
		"taskId":         taskId,
		"status":         "IN_PROGRESS",
		"startTime":      time.Now().UTC().Unix(),
		"target":         request.GetParamCaseInsensitive(req.Parameters, "target"),
		"actions":        request.GetParamCaseInsensitive(req.Parameters, "actions"),
		"violationEvent": request.GetParamCaseInsensitive(req.Parameters, "violationEvent"),
	}
	if err := store.PutGeneric("detectMitigationTask/"+taskId, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{"taskId": taskId}, nil
}
func (s *IoTService) CancelDetectMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "detectMitigationTask/" + taskId
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrDetectMitigationTaskNotFound
	}
	rec["status"] = "CANCELED"
	rec["endTime"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeDetectMitigationActionsTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("detectMitigationTask/"+taskId, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrDetectMitigationTaskNotFound
	}
	return map[string]interface{}{
		"taskId":         rec["taskId"],
		"status":         rec["status"],
		"startTime":      rec["startTime"],
		"endTime":        rec["endTime"],
		"target":         rec["target"],
		"actions":        rec["actions"],
		"violationEvent": rec["violationEvent"],
	}, nil
}
func (s *IoTService) ListDetectMitigationActionsExecutions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return paginatedMaps("taskExecutions", []map[string]interface{}{}, req.Parameters), nil
}
func (s *IoTService) ListDetectMitigationActionsTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("detectMitigationTask/")
	if err != nil {
		return nil, err
	}
	tasks := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		tasks = append(tasks, map[string]interface{}{
			"taskId":        rec["taskId"],
			"taskStatus":    rec["status"],
			"taskStartTime": rec["startTime"],
		})
	}
	return paginatedMaps("tasks", tasks, req.Parameters), nil
}
