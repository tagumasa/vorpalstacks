package iot

import (
	"context"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func (s *IoTService) CreateJobTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "jobTemplateId")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{
		"jobTemplateId": name,
		"description":   request.GetParamCaseInsensitive(req.Parameters, "description"),
		"document":      request.GetParamCaseInsensitive(req.Parameters, "document"),
		"createdAt":     time.Now().Unix(),
	}
	if err := store.PutGeneric("jobTemplate/"+name, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"jobTemplateId":  name,
		"jobTemplateArn": iotstore.BuildJobTemplateARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
	}, nil
}

func (s *IoTService) DeleteJobTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "jobTemplateId")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	exists, err := store.GetGenericExists("jobTemplate/"+name, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrJobTemplateNotFound
	}
	if err := store.DeleteGeneric("jobTemplate/" + name); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DescribeJobTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "jobTemplateId")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("jobTemplate/"+name, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrJobTemplateNotFound
	}
	return rec, nil
}

func (s *IoTService) ListJobTemplates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("jobTemplate/")
	if err != nil {
		return nil, err
	}
	return paginatedMaps("jobTemplates", items, req.Parameters), nil
}

func (s *IoTService) DescribeManagedJobTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "templateName")
	return map[string]interface{}{
		"templateName": name,
		"templateArn":  iotstore.BuildJobTemplateARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
		"description":  "AWS-provided managed job template",
		"platform":     "Linux",
		"pathToDefine": "",
	}, nil
}

func (s *IoTService) ListManagedJobTemplates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return paginatedMaps("managedJobTemplates", []map[string]interface{}{}, req.Parameters), nil
}

func (s *IoTService) CancelJobExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	jobID := request.GetParamCaseInsensitive(req.Parameters, "jobId")
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	if jobID == "" || thingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Validate the parent job exists; AWS returns ResourceNotFoundException.
	if _, err := store.GetJob(jobID); err != nil {
		return nil, err
	}
	key := "jobExecution/" + jobID + "/" + thingName
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrJobExecutionNotFound
	}
	rec["status"] = "CANCELED"
	rec["forceCanceled"] = request.GetBoolParam(req.Parameters, "force")
	rec["lastUpdatedAt"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DeleteJobExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	jobID := request.GetParamCaseInsensitive(req.Parameters, "jobId")
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	if jobID == "" || thingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetJob(jobID); err != nil {
		return nil, err
	}
	key := "jobExecution/" + jobID + "/" + thingName
	exists, err := store.GetGenericExists(key, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrJobExecutionNotFound
	}
	if err := store.DeleteGeneric(key); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DescribeJobExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	jobID := request.GetParamCaseInsensitive(req.Parameters, "jobId")
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	if jobID == "" || thingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetJob(jobID); err != nil {
		return nil, err
	}
	key := "jobExecution/" + jobID + "/" + thingName
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrJobExecutionNotFound
	}
	return map[string]interface{}{
		"execution": map[string]interface{}{
			"jobId":           jobID,
			"status":          rec["status"],
			"executionNumber": rec["executionNumber"],
			"queuedAt":        rec["queuedAt"],
			"startedAt":       rec["startedAt"],
			"lastUpdatedAt":   rec["lastUpdatedAt"],
			"thingArn":        iotstore.BuildThingARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), thingName),
			"versionNumber":   rec["versionNumber"],
		},
	}, nil
}

func (s *IoTService) ListJobExecutionsForJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	jobID := request.GetParamCaseInsensitive(req.Parameters, "jobId")
	if jobID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Validate the parent job exists; AWS returns ResourceNotFoundException.
	if _, err := store.GetJob(jobID); err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("jobExecution/" + jobID + "/")
	if err != nil {
		return nil, err
	}
	summaries := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		thingName, _ := rec["thingName"].(string)
		summaries = append(summaries, map[string]interface{}{
			"thingArn": iotstore.BuildThingARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), thingName),
			"jobExecutionSummary": map[string]interface{}{
				"status":          rec["status"],
				"executionNumber": rec["executionNumber"],
				"queuedAt":        rec["queuedAt"],
				"versionNumber":   rec["versionNumber"],
			},
		})
	}
	return paginatedMaps("executionSummaries", summaries, req.Parameters), nil
}

func (s *IoTService) ListJobExecutionsForThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Scan all job-execution records and filter by thingName.
	allItems, err := store.ListGeneric("jobExecution/")
	if err != nil {
		return nil, err
	}
	summaries := make([]map[string]interface{}, 0)
	for _, rec := range allItems {
		rThing, _ := rec["thingName"].(string)
		if rThing != thingName {
			continue
		}
		jobID, _ := rec["jobId"].(string)
		summaries = append(summaries, map[string]interface{}{
			"jobArn": iotstore.BuildJobARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), jobID),
			"jobExecutionSummary": map[string]interface{}{
				"status":          rec["status"],
				"executionNumber": rec["executionNumber"],
				"queuedAt":        rec["queuedAt"],
				"versionNumber":   rec["versionNumber"],
			},
		})
	}
	return paginatedMaps("executionSummaries", summaries, req.Parameters), nil
}

func (s *IoTService) CreateOTAUpdate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "otaUpdateId")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	otaID := uuid.New().String()
	rec := map[string]interface{}{
		"otaUpdateId":  name,
		"otaUpdateArn": iotstore.BuildOTAUpdateARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
		"status":       "CREATE_COMPLETE",
		"createdAt":    time.Now().Unix(),
	}
	if err := store.PutGeneric("otaUpdate/"+name, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"otaUpdateId":     name,
		"otaUpdateArn":    iotstore.BuildOTAUpdateARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
		"otaUpdateStatus": "CREATE_COMPLETE",
		"awsIotJobId":     otaID,
	}, nil
}

func (s *IoTService) DeleteOTAUpdate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "otaUpdateId")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	exists, err := store.GetGenericExists("otaUpdate/"+name, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrJobNotFound
	}
	if err := store.DeleteGeneric("otaUpdate/" + name); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) GetOTAUpdate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "otaUpdateId")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("otaUpdate/"+name, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrJobNotFound
	}
	// AWS wraps the OTA update in otaUpdateInfo; returning it flat leaves the
	// SDK's OtaUpdateInfo field nil.
	return map[string]interface{}{"otaUpdateInfo": rec}, nil
}

func (s *IoTService) ListOTAUpdates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("otaUpdate/")
	if err != nil {
		return nil, err
	}
	return paginatedMaps("otaUpdates", items, req.Parameters), nil
}
