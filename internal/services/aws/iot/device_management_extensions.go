package iot

import (
	"context"
	"strings"

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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.addThingToBillingGroupCore(store, thingName, billingGroup); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) RemoveThingFromBillingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	billingGroup := request.GetParamCaseInsensitive(req.Parameters, "billingGroupName")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.removeThingFromBillingGroupCore(store, thingName, billingGroup); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ListThingsInBillingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	billingGroup := request.GetParamCaseInsensitive(req.Parameters, "billingGroupName")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	things, err := s.listThingsInBillingGroupCore(store, billingGroup)
	if err != nil {
		return nil, err
	}
	return paginatedStrings("things", things, req.Parameters)
}

func (s *IoTService) RegisterThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	parameters := map[string]string{}
	for k, v := range request.GetMapParamCaseInsensitive(req.Parameters, "parameters") {
		if str, ok := v.(string); ok {
			parameters[k] = str
		}
	}
	in := RegisterThingInput{
		TemplateBody: request.GetParamCaseInsensitive(req.Parameters, "templateBody"),
		Parameters:   parameters,
	}
	result, err := s.registerThingCore(store, s.caForReq(reqCtx), in)
	if err != nil {
		return nil, err
	}
	resp := map[string]interface{}{
		"resourceArns": result.ResourceArns,
	}
	if result.CertificatePem != "" {
		resp["certificatePem"] = result.CertificatePem
	}
	return resp, nil
}

func (s *IoTService) StartThingRegistrationTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := StartThingRegistrationTaskInput{
		TemplateBody:    request.GetParamCaseInsensitive(req.Parameters, "templateBody"),
		InputFileBucket: request.GetParamCaseInsensitive(req.Parameters, "inputFileBucket"),
		InputFileKey:    request.GetParamCaseInsensitive(req.Parameters, "inputFileKey"),
		RoleArn:         request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
	}
	if reqCtx != nil {
		in.Region = reqCtx.GetRegion()
	}
	taskID, err := s.startThingRegistrationTaskCore(store, in)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"taskId": taskID,
	}, nil
}

func (s *IoTService) StopThingRegistrationTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskID := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.stopThingRegistrationTaskCore(store, taskID); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DescribeThingRegistrationTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskID := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.describeThingRegistrationTaskCore(store, taskID)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"taskId":             result.TaskID,
		"status":             result.Status,
		"creationDate":       result.CreationDate,
		"lastModifiedDate":   result.LastModifiedDate,
		"templateBody":       result.TemplateBody,
		"inputFileBucket":    result.InputFileBucket,
		"inputFileKey":       result.InputFileKey,
		"roleArn":            result.RoleArn,
		"message":            result.Message,
		"successCount":       result.SuccessCount,
		"failureCount":       result.FailureCount,
		"percentageProgress": result.PercentageProgress,
	}, nil
}

func (s *IoTService) ListThingRegistrationTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Smithy: ListThingRegistrationTasksResponse.taskIds is list<TaskId> (string).
	taskIds, err := s.listThingRegistrationTasksCore(store)
	if err != nil {
		return nil, err
	}
	return paginatedStrings("taskIds", taskIds, req.Parameters)
}

func (s *IoTService) ListThingRegistrationTaskReports(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskID := request.GetParamCaseInsensitive(req.Parameters, "taskId")
	if taskID == "" {
		return nil, iotstore.ErrMissingParam
	}
	// The Smithy ReportType enum's wire values are ERRORS and RESULTS.
	reportType := strings.ToUpper(request.GetParamCaseInsensitive(req.Parameters, "reportType"))
	if reportType != "ERRORS" && reportType != "RESULTS" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	links, err := s.thingRegistrationReportLinks(reqCtx, req, store, taskID, reportType)
	if err != nil {
		return nil, err
	}
	// The response shape is reportType/resourceLinks/nextToken; the link
	// expiry lives inside the presigned URLs, not as a response member.
	return map[string]interface{}{
		"reportType":    reportType,
		"resourceLinks": links,
	}, nil
}
