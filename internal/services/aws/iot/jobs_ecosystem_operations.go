package iot

import (
	"context"
	"strconv"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func (s *IoTService) CreateJobTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.createJobTemplateCore(store, CreateJobTemplateInput{
		JobTemplateID:              request.GetParamCaseInsensitive(req.Parameters, "jobTemplateId"),
		JobArn:                     request.GetParamCaseInsensitive(req.Parameters, "jobArn"),
		Description:                request.GetParamCaseInsensitive(req.Parameters, "description"),
		Document:                   request.GetParamCaseInsensitive(req.Parameters, "document"),
		DocumentSource:             request.GetParamCaseInsensitive(req.Parameters, "documentSource"),
		PresignedUrlConfig:         req.Parameters["presignedUrlConfig"],
		JobExecutionsRolloutConfig: req.Parameters["jobExecutionsRolloutConfig"],
		AbortConfig:                req.Parameters["abortConfig"],
		TimeoutConfig:              req.Parameters["timeoutConfig"],
		JobExecutionsRetryConfig:   req.Parameters["jobExecutionsRetryConfig"],
		MaintenanceWindows:         req.Parameters["maintenanceWindows"],
		DestinationPackageVersions: req.Parameters["destinationPackageVersions"],
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
		request.GetParamCaseInsensitive(req.Parameters, "templateName"))
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
		JobID:                   request.GetParamCaseInsensitive(req.Parameters, "jobId"),
		ThingName:               request.GetParamCaseInsensitive(req.Parameters, "thingName"),
		Force:                   request.GetBoolParam(req.Parameters, "force"),
		ExpectedVersion:         request.GetInt64Param(req.Parameters, "expectedVersion"),
		ExpectedVersionProvided: request.HasParam(req.Parameters, "expectedVersion"),
		StatusDetails:           req.Parameters["statusDetails"],
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
	// executionNumber travels as a URI path label, so it arrives as a
	// string on the wire.
	numberRaw := request.GetParamCaseInsensitive(req.Parameters, "executionNumber")
	executionNumber := int64(0)
	if numberRaw != "" {
		parsed, parseErr := strconv.ParseInt(numberRaw, 10, 64)
		if parseErr != nil {
			return nil, iotstore.ErrValidation
		}
		executionNumber = parsed
	}
	if err := s.deleteJobExecutionCore(store, DeleteJobExecutionInput{
		JobID:           request.GetParamCaseInsensitive(req.Parameters, "jobId"),
		ThingName:       request.GetParamCaseInsensitive(req.Parameters, "thingName"),
		ExecutionNumber: executionNumber,
		Force:           request.GetBoolParam(req.Parameters, "force"),
	}); err != nil {
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
		AwsJobExecutionsRolloutConfig: req.Parameters["awsJobExecutionsRolloutConfig"],
		AwsJobPresignedUrlConfig:      req.Parameters["awsJobPresignedUrlConfig"],
		AwsJobAbortConfig:             req.Parameters["awsJobAbortConfig"],
		AwsJobTimeoutConfig:           req.Parameters["awsJobTimeoutConfig"],
		Files:                         req.Parameters["files"],
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
	if err := s.deleteOTAUpdateCore(store, DeleteOTAUpdateInput{
		OtaUpdateID:       request.GetParamCaseInsensitive(req.Parameters, "otaUpdateId"),
		DeleteStream:      request.GetBoolParam(req.Parameters, "deleteStream"),
		ForceDeleteAWSJob: request.GetBoolParam(req.Parameters, "forceDeleteAWSJob"),
	}); err != nil {
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
	items, err := s.listOTAUpdatesCore(store, request.GetParamCaseInsensitive(req.Parameters, "otaUpdateStatus"))
	if err != nil {
		return nil, err
	}
	return paginatedMaps("otaUpdates", items, req.Parameters)
}
