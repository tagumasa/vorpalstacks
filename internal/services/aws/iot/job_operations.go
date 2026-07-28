package iot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// CreateJob creates a new IoT job targeting one or more things.
func (s *IoTService) CreateJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	jobID := request.GetParamCaseInsensitive(req.Parameters, "jobId")
	if jobID == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	job := &iotstore.Job{
		JobID:                      jobID,
		Description:                request.GetParamCaseInsensitive(req.Parameters, "description"),
		Document:                   request.GetParamCaseInsensitive(req.Parameters, "document"),
		Targets:                    request.GetStringList(req.Parameters, "targets"),
		Status:                     "IN_PROGRESS",
		TargetSelection:            request.GetParamCaseInsensitive(req.Parameters, "targetSelection"),
		ReasonCode:                 request.GetParamCaseInsensitive(req.Parameters, "reasonCode"),
		Comment:                    request.GetParamCaseInsensitive(req.Parameters, "comment"),
		NamespaceID:                request.GetParamCaseInsensitive(req.Parameters, "namespaceId"),
		JobTemplateARN:             request.GetParamCaseInsensitive(req.Parameters, "jobTemplateArn"),
		PresignedUrlConfig:         request.GetParamCaseInsensitive(req.Parameters, "presignedUrlConfig"),
		JobExecutionsRolloutConfig: request.GetParamCaseInsensitive(req.Parameters, "jobExecutionsRolloutConfig"),
		AbortConfig:                request.GetParamCaseInsensitive(req.Parameters, "abortConfig"),
		TimeoutConfig:              request.GetParamCaseInsensitive(req.Parameters, "timeoutConfig"),
		JobExecutionsRetryConfig:   request.GetParamCaseInsensitive(req.Parameters, "jobExecutionsRetryConfig"),
		DocumentParameters:         request.GetParamCaseInsensitive(req.Parameters, "documentParameters"),
		SchedulingConfig:           request.GetParamCaseInsensitive(req.Parameters, "schedulingConfig"),
		ScheduledJobRollouts:       request.GetParamCaseInsensitive(req.Parameters, "scheduledJobRollouts"),
		DestinationPackageVersions: request.GetParamCaseInsensitive(req.Parameters, "destinationPackageVersions"),
		CreatedAt:                  time.Now().UTC(),
		LastUpdatedAt:              time.Now().UTC(),
	}

	if job.TargetSelection != "" {
		if err := ValidateTargetSelection(job.TargetSelection); err != nil {
			return nil, err
		}
	}

	tags, err := parseTags(req.Parameters)
	if err != nil {
		return nil, iotstore.ErrValidation
	}
	if len(tags) > 0 {
		job.Tags = tags
	}

	created, err := store.CreateJob(job)
	if err != nil {
		return nil, err
	}

	// Create a job-execution record for each target thing so that
	// ListJobExecutionsForJob, ListJobExecutionsForThing and
	// DescribeJobExecution return real data matching AWS behaviour.
	now := time.Now().UTC().Unix()
	for _, targetARN := range created.Targets {
		thingName := iotstore.ThingNameFromARN(targetARN)
		if thingName == "" {
			continue
		}
		execKey := "jobExecution/" + jobID + "/" + thingName
		execRec := map[string]interface{}{
			"jobId":           jobID,
			"thingName":       thingName,
			"status":          "QUEUED",
			"executionNumber": int64(1),
			"queuedAt":        now,
			"versionNumber":   int64(1),
		}
		if err := store.PutGeneric(execKey, execRec); err != nil {
			return nil, err
		}
	}

	return jobResponse(created), nil
}

// DescribeJob retrieves details of a job including its status and configuration.
func (s *IoTService) DescribeJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	jobID := request.GetParamCaseInsensitive(req.Parameters, "jobId")
	if jobID == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	job, err := store.GetJob(jobID)
	if err != nil {
		return nil, err
	}

	jobMap := jobResponse(job)
	jobMap["version"] = job.Version
	return map[string]interface{}{
		"job":      jobMap,
		"document": job.Document,
	}, nil
}

// DeleteJob removes a job from the registry.
func (s *IoTService) DeleteJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	jobID := request.GetParamCaseInsensitive(req.Parameters, "jobId")
	if jobID == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	arn := iotstore.BuildJobARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), jobID)
	_ = store.DeleteAllTags(arn)

	if err := store.DeleteJob(jobID); err != nil {
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

	jobs, err := store.ListJobs(parseListOptions(req.Parameters), request.GetParamCaseInsensitive(req.Parameters, "status"))
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(jobs.Items))
	for _, j := range jobs.Items {
		result = append(result, jobResponse(j))
	}

	return listResponse("jobs", result, jobs.NextMarker), nil
}

// CancelJob cancels a running or queued job.
func (s *IoTService) CancelJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	jobID := request.GetParamCaseInsensitive(req.Parameters, "jobId")
	if jobID == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	job, err := store.GetJob(jobID)
	if err != nil {
		return nil, err
	}
	if job.Status != "IN_PROGRESS" && job.Status != "SCHEDULED" {
		return nil, iotstore.ErrInvalidRequest
	}

	opts := iotstore.JobUpdateOpts{Status: "CANCELED"}
	job, err = store.UpdateJob(jobID, opts)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"jobId":         job.JobID,
		"jobArn":        job.JobARN,
		"description":   job.Description,
		"status":        job.Status,
		"lastUpdatedAt": job.LastUpdatedAt.Unix(),
	}, nil
}

// GetJobDocument retrieves the job document content.
func (s *IoTService) GetJobDocument(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	jobID := request.GetParamCaseInsensitive(req.Parameters, "jobId")
	if jobID == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	job, err := store.GetJob(jobID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"document": job.Document,
	}, nil
}

// UpdateJob modifies a job's description or status.
func (s *IoTService) UpdateJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	jobID := request.GetParamCaseInsensitive(req.Parameters, "jobId")
	if jobID == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := iotstore.JobUpdateOpts{
		Description:                request.GetParamCaseInsensitive(req.Parameters, "description"),
		Status:                     request.GetParamCaseInsensitive(req.Parameters, "status"),
		ReasonCode:                 request.GetParamCaseInsensitive(req.Parameters, "reasonCode"),
		Comment:                    request.GetParamCaseInsensitive(req.Parameters, "comment"),
		NamespaceID:                request.GetParamCaseInsensitive(req.Parameters, "namespaceId"),
		PresignedUrlConfig:         request.GetParamCaseInsensitive(req.Parameters, "presignedUrlConfig"),
		JobExecutionsRolloutConfig: request.GetParamCaseInsensitive(req.Parameters, "jobExecutionsRolloutConfig"),
		AbortConfig:                request.GetParamCaseInsensitive(req.Parameters, "abortConfig"),
		TimeoutConfig:              request.GetParamCaseInsensitive(req.Parameters, "timeoutConfig"),
		JobExecutionsRetryConfig:   request.GetParamCaseInsensitive(req.Parameters, "jobExecutionsRetryConfig"),
		SchedulingConfig:           request.GetParamCaseInsensitive(req.Parameters, "schedulingConfig"),
		ScheduledJobRollouts:       request.GetParamCaseInsensitive(req.Parameters, "scheduledJobRollouts"),
	}
	_, err = store.UpdateJob(jobID, opts)
	if err != nil {
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
