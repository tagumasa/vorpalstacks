package iot

import (
	"context"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func (s *IoTService) CreateDynamicThingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.CreateThingGroup(ctx, reqCtx, req)
}

func (s *IoTService) DeleteDynamicThingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.DeleteThingGroup(ctx, reqCtx, req)
}

func (s *IoTService) UpdateDynamicThingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.UpdateThingGroup(ctx, reqCtx, req)
}

func (s *IoTService) AddThingToBillingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	billingGroup := request.GetParamCaseInsensitive(req.Parameters, "billingGroupName")
	if thingName == "" || billingGroup == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.AddThingToBillingGroup(thingName, billingGroup); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) RemoveThingFromBillingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	billingGroup := request.GetParamCaseInsensitive(req.Parameters, "billingGroupName")
	if thingName == "" || billingGroup == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.RemoveThingFromBillingGroup(thingName, billingGroup); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ListThingsInBillingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	billingGroup := request.GetParamCaseInsensitive(req.Parameters, "billingGroupName")
	if billingGroup == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// AWS: returns ResourceNotFoundException for a non-existent billing group.
	if _, err := store.GetBillingGroup(billingGroup); err != nil {
		return nil, err
	}
	things, err := store.ListThingsInBillingGroup(billingGroup)
	if err != nil {
		return nil, err
	}
	return paginatedStrings("things", things, req.Parameters), nil
}

func (s *IoTService) RegisterThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	templateName := request.GetParamCaseInsensitive(req.Parameters, "templateName")
	if templateName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	if thingName == "" {
		thingName = "registered-" + uuid.New().String()[:8]
	}
	thing := &iotstore.Thing{
		ThingName:        thingName,
		Version:          1,
		CreationDate:     time.Now().UTC(),
		LastModifiedDate: time.Now().UTC(),
	}
	created, err := store.CreateThing(thing)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"thingName": created.ThingName,
		"thingArn":  created.ThingARN,
	}, nil
}

func (s *IoTService) StartThingRegistrationTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	taskID := uuid.New().String()
	rec := map[string]interface{}{
		"taskId":       taskID,
		"templateBody": request.GetParamCaseInsensitive(req.Parameters, "templateBody"),
		"status":       "Completed",
		"creationDate": time.Now().UTC().Unix(),
	}
	if err := store.PutGeneric("registrationTask/"+taskID, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"taskId": taskID,
	}, nil
}

func (s *IoTService) StopThingRegistrationTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskID := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("registrationTask/"+taskID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrThingRegistrationTaskNotFound
	}
	rec["status"] = "Cancelled"
	if err := store.PutGeneric("registrationTask/"+taskID, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DescribeThingRegistrationTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskID := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("registrationTask/"+taskID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrThingRegistrationTaskNotFound
	}
	status, _ := rec["status"].(string)
	return map[string]interface{}{
		"taskId":               taskID,
		"status":               status,
		"creationDate":         rec["creationDate"],
		"successfulExecutions": 0,
		"failedExecutions":     0,
		"percentageProgress":   100,
	}, nil
}

func (s *IoTService) ListThingRegistrationTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("registrationTask/")
	if err != nil {
		return nil, err
	}
	// Smithy: ListThingRegistrationTasksResponse.taskIds is list<TaskId> (string).
	taskIds := make([]string, 0, len(items))
	for _, item := range items {
		if id, ok := item["taskId"].(string); ok {
			taskIds = append(taskIds, id)
		}
	}
	return paginatedStrings("taskIds", taskIds, req.Parameters), nil
}

func (s *IoTService) ListThingRegistrationTaskReports(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskID := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskID == "" {
		return nil, iotstore.ErrMissingParam
	}
	return paginatedMaps("taskReports", []map[string]interface{}{}, req.Parameters), nil
}
