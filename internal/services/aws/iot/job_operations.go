package iot

import (
	"context"
	"encoding/json"
	"fmt"

	"vorpalstacks/internal/common/request"
)

// CreateJob creates a new IoT job targeting one or more things.
func (s *IoTService) CreateJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tags, tagsErr := parseTags(req.Parameters)

	in := CreateJobInput{
		JobID:                      request.GetParamCaseInsensitive(req.Parameters, "jobId"),
		Description:                request.GetParamCaseInsensitive(req.Parameters, "description"),
		Document:                   request.GetParamCaseInsensitive(req.Parameters, "document"),
		Targets:                    request.GetStringList(req.Parameters, "targets"),
		TargetSelection:            request.GetParamCaseInsensitive(req.Parameters, "targetSelection"),
		NamespaceID:                request.GetParamCaseInsensitive(req.Parameters, "namespaceId"),
		JobTemplateARN:             request.GetParamCaseInsensitive(req.Parameters, "jobTemplateArn"),
		PresignedUrlConfig:         request.GetParamCaseInsensitive(req.Parameters, "presignedUrlConfig"),
		JobExecutionsRolloutConfig: request.GetParamCaseInsensitive(req.Parameters, "jobExecutionsRolloutConfig"),
		AbortConfig:                request.GetParamCaseInsensitive(req.Parameters, "abortConfig"),
		TimeoutConfig:              request.GetParamCaseInsensitive(req.Parameters, "timeoutConfig"),
		JobExecutionsRetryConfig:   request.GetParamCaseInsensitive(req.Parameters, "jobExecutionsRetryConfig"),
		DocumentParameters:         request.GetParamCaseInsensitive(req.Parameters, "documentParameters"),
		SchedulingConfig:           request.GetParamCaseInsensitive(req.Parameters, "schedulingConfig"),
		DestinationPackageVersions: request.GetParamCaseInsensitive(req.Parameters, "destinationPackageVersions"),
		Tags:                       tags,
		TagsMalformed:              tagsErr != nil,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createJobCore(store, in)
	if err != nil {
		return nil, err
	}
	return jobResponse(result.Job), nil
}

// DescribeJob retrieves details of a job including its status and configuration.
func (s *IoTService) DescribeJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.describeJobCore(store, request.GetParamCaseInsensitive(req.Parameters, "jobId"))
	if err != nil {
		return nil, err
	}

	// DescribeJobResponse members are documentSource and job; the document
	// itself is served by GetJobDocument.
	return map[string]interface{}{
		"job": jobResponse(result.Job),
	}, nil
}

// DeleteJob removes a job from the registry.
func (s *IoTService) DeleteJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteJobCore(store, request.GetParamCaseInsensitive(req.Parameters, "jobId")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

// ListJobs returns all jobs, optionally filtered by status.
func (s *IoTService) ListJobs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := parseListOptions(req.Parameters)
	result, err := s.listJobsCore(store, ListJobsInput{
		Status:    request.GetParamCaseInsensitive(req.Parameters, "status"),
		NextToken: opts.Marker,
		MaxItems:  opts.MaxItems,
	})
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(result.Jobs))
	for _, j := range result.Jobs {
		items = append(items, jobResponse(j))
	}

	return listResponse("jobs", items, result.NextToken), nil
}

// CancelJob cancels a running or queued job.
func (s *IoTService) CancelJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.cancelJobCore(store, CancelJobInput{
		JobID:      request.GetParamCaseInsensitive(req.Parameters, "jobId"),
		ReasonCode: request.GetParamCaseInsensitive(req.Parameters, "reasonCode"),
		Comment:    request.GetParamCaseInsensitive(req.Parameters, "comment"),
	})
	if err != nil {
		return nil, err
	}

	// The response shape is jobArn/jobId/description.
	return map[string]interface{}{
		"jobId":       result.Job.JobID,
		"jobArn":      result.Job.JobARN,
		"description": result.Job.Description,
	}, nil
}

// GetJobDocument retrieves the job document content.
func (s *IoTService) GetJobDocument(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.getJobDocumentCore(store, request.GetParamCaseInsensitive(req.Parameters, "jobId"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"document": result.Document,
	}, nil
}

// UpdateJob modifies a job's description or status.
func (s *IoTService) UpdateJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := UpdateJobInput{
		JobID:                      request.GetParamCaseInsensitive(req.Parameters, "jobId"),
		Description:                request.GetParamCaseInsensitive(req.Parameters, "description"),
		NamespaceID:                request.GetParamCaseInsensitive(req.Parameters, "namespaceId"),
		PresignedUrlConfig:         request.GetParamCaseInsensitive(req.Parameters, "presignedUrlConfig"),
		JobExecutionsRolloutConfig: request.GetParamCaseInsensitive(req.Parameters, "jobExecutionsRolloutConfig"),
		AbortConfig:                request.GetParamCaseInsensitive(req.Parameters, "abortConfig"),
		TimeoutConfig:              request.GetParamCaseInsensitive(req.Parameters, "timeoutConfig"),
		JobExecutionsRetryConfig:   request.GetParamCaseInsensitive(req.Parameters, "jobExecutionsRetryConfig"),
	}
	if err := s.updateJobCore(store, in); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

// parseTags extracts a tags map from a tags parameter.
// Handles both JSON-encoded string form and direct []interface{} from JSON body.
// Returns an error if the JSON-encoded string form is malformed.
func parseTags(params map[string]interface{}) (map[string]string, error) {
	result := make(map[string]string)

	if raw, ok := params["tags"]; ok {
		switch v := raw.(type) {
		case []interface{}:
			for _, t := range v {
				if m, ok := t.(map[string]interface{}); ok {
					key, _ := m["Key"].(string)
					val, _ := m["Value"].(string)
					if key != "" {
						result[key] = val
					}
				}
			}
			return result, nil
		case string:
			if v == "" {
				return nil, nil
			}
			var tagList []map[string]interface{}
			if err := json.Unmarshal([]byte(v), &tagList); err != nil {
				return nil, fmt.Errorf("malformed tags JSON: %w", err)
			}
			for _, t := range tagList {
				key, _ := t["Key"].(string)
				val, _ := t["Value"].(string)
				if key != "" {
					result[key] = val
				}
			}
		}
	}
	return result, nil
}

func (s *IoTService) AssociateTargetsWithJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.associateTargetsWithJobCore(store, AssociateTargetsWithJobInput{
		JobID:   request.GetParamCaseInsensitive(req.Parameters, "jobId"),
		Targets: request.GetStringList(req.Parameters, "targets"),
		Comment: request.GetParamCaseInsensitive(req.Parameters, "comment"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"jobArn":      result.JobARN,
		"jobId":       result.JobID,
		"description": result.Description,
	}, nil
}
