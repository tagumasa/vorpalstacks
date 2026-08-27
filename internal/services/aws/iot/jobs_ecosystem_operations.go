package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
)

func (s *IoTService) CreateJobTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.createJobTemplateCore(store, CreateJobTemplateInput{
		JobTemplateID: request.GetParamCaseInsensitive(req.Parameters, "jobTemplateId"),
		Description:   request.GetParamCaseInsensitive(req.Parameters, "description"),
		Document:      request.GetParamCaseInsensitive(req.Parameters, "document"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"jobTemplateId":  result.JobTemplateID,
		"jobTemplateArn": result.JobTemplateARN,
	}, nil
}

func (s *IoTService) DeleteJobTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteJobTemplateCore(store, request.GetParamCaseInsensitive(req.Parameters, "jobTemplateId")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DescribeJobTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.describeJobTemplateCore(store, request.GetParamCaseInsensitive(req.Parameters, "jobTemplateId"))
}

func (s *IoTService) ListJobTemplates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := s.listJobTemplatesCore(store)
	if err != nil {
		return nil, err
	}
	return paginatedMaps("jobTemplates", items, req.Parameters)
}

func (s *IoTService) DescribeManagedJobTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.describeManagedJobTemplateCore(
		reqCtx.GetAccountID(), reqCtx.GetRegion(),
		request.GetParamCaseInsensitive(req.Parameters, "templateName")), nil
}

func (s *IoTService) ListManagedJobTemplates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return paginatedMaps("managedJobTemplates", []map[string]interface{}{}, req.Parameters)
}

func (s *IoTService) CancelJobExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.cancelJobExecutionCore(store, CancelJobExecutionInput{
		JobID:     request.GetParamCaseInsensitive(req.Parameters, "jobId"),
		ThingName: request.GetParamCaseInsensitive(req.Parameters, "thingName"),
		Force:     request.GetBoolParam(req.Parameters, "force"),
	}); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DeleteJobExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteJobExecutionCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "jobId"),
		request.GetParamCaseInsensitive(req.Parameters, "thingName")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DescribeJobExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.describeJobExecutionCore(store, DescribeJobExecutionInput{
		JobID:     request.GetParamCaseInsensitive(req.Parameters, "jobId"),
		ThingName: request.GetParamCaseInsensitive(req.Parameters, "thingName"),
	})
}

func (s *IoTService) ListJobExecutionsForJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	summaries, err := s.listJobExecutionsForJobCore(store, request.GetParamCaseInsensitive(req.Parameters, "jobId"))
	if err != nil {
		return nil, err
	}
	return paginatedMaps("executionSummaries", summaries, req.Parameters)
}

func (s *IoTService) ListJobExecutionsForThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	summaries, err := s.listJobExecutionsForThingCore(store, request.GetParamCaseInsensitive(req.Parameters, "thingName"))
	if err != nil {
		return nil, err
	}
	return paginatedMaps("executionSummaries", summaries, req.Parameters)
}

func (s *IoTService) CreateOTAUpdate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tagList := tags.ParseTagsWithQueryFallback(req.Parameters, "tags")
	recTags := make(map[string]string, len(tagList))
	for _, t := range tagList {
		recTags[t.Key] = t.Value
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.createOTAUpdateCore(store, CreateOTAUpdateInput{
		OtaUpdateID:                   request.GetParamCaseInsensitive(req.Parameters, "otaUpdateId"),
		Description:                   request.GetParamCaseInsensitive(req.Parameters, "description"),
		Targets:                       request.GetStringList(req.Parameters, "targets"),
		Protocols:                     request.GetStringList(req.Parameters, "protocols"),
		TargetSelection:               request.GetParamCaseInsensitive(req.Parameters, "targetSelection"),
		AwsJobExecutionsRolloutConfig: request.GetParamCaseInsensitive(req.Parameters, "awsJobExecutionsRolloutConfig"),
		AwsJobPresignedUrlConfig:      request.GetParamCaseInsensitive(req.Parameters, "awsJobPresignedUrlConfig"),
		AwsJobAbortConfig:             request.GetParamCaseInsensitive(req.Parameters, "awsJobAbortConfig"),
		AwsJobTimeoutConfig:           request.GetParamCaseInsensitive(req.Parameters, "awsJobTimeoutConfig"),
		RoleArn:                       request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		Tags:                          recTags,
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"otaUpdateId":     result.OtaUpdateID,
		"otaUpdateArn":    result.OtaUpdateArn,
		"otaUpdateStatus": result.OtaUpdateStatus,
		"awsIotJobId":     result.AwsIotJobID,
		"awsIotJobArn":    result.AwsIotJobArn,
	}, nil
}

func (s *IoTService) DeleteOTAUpdate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteOTAUpdateCore(store, request.GetParamCaseInsensitive(req.Parameters, "otaUpdateId")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) GetOTAUpdate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getOTAUpdateCore(store, request.GetParamCaseInsensitive(req.Parameters, "otaUpdateId"))
}

func (s *IoTService) ListOTAUpdates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := s.listOTAUpdatesCore(store)
	if err != nil {
		return nil, err
	}
	return paginatedMaps("otaUpdates", items, req.Parameters)
}
