package iot

import (
	"context"
	"encoding/json"
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
		JobID:           jobID,
		Description:     request.GetParamCaseInsensitive(req.Parameters, "description"),
		Document:        request.GetParamCaseInsensitive(req.Parameters, "document"),
		Targets:         request.GetStringList(req.Parameters, "targets"),
		Status:          "QUEUED",
		TargetSelection: request.GetParamCaseInsensitive(req.Parameters, "targetSelection"),
		CreatedAt:       time.Now().UTC(),
		LastUpdatedAt:   time.Now().UTC(),
	}

	if job.TargetSelection != "" {
		if err := ValidateTargetSelection(job.TargetSelection); err != nil {
			return nil, err
		}
	}

	tags := parseTags(req.Parameters)
	if len(tags) > 0 {
		job.Tags = tags
	}

	created, err := store.CreateJob(job)
	if err != nil {
		return nil, err
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
		return nil, iotstore.ErrJobNotFound
	}

	jobMap := map[string]interface{}{
		"jobId":           job.JobID,
		"jobArn":          job.JobARN,
		"description":     job.Description,
		"status":          job.Status,
		"targetSelection": job.TargetSelection,
		"targets":         job.Targets,
		"version":         job.Version,
		"createdAt":       job.CreatedAt.Unix(),
		"lastUpdatedAt":   job.LastUpdatedAt.Unix(),
	}
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
		result = append(result, map[string]interface{}{
			"jobId":            j.JobID,
			"jobArn":           j.JobARN,
			"description":      j.Description,
			"status":           j.Status,
			"targetSelection":  j.TargetSelection,
			"creationDate":     j.CreatedAt.Unix(),
			"lastModifiedDate": j.LastUpdatedAt.Unix(),
		})
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

	opts := iotstore.JobUpdateOpts{Status: "CANCELLED"}
	job, err := store.UpdateJob(jobID, opts)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"jobId":            job.JobID,
		"jobArn":           job.JobARN,
		"description":      job.Description,
		"status":           job.Status,
		"lastModifiedDate": job.LastUpdatedAt.Unix(),
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
		return nil, iotstore.ErrJobNotFound
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
		Description: request.GetParamCaseInsensitive(req.Parameters, "description"),
		Status:      request.GetParamCaseInsensitive(req.Parameters, "status"),
	}
	_, err = store.UpdateJob(jobID, opts)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// parseTags extracts a tags map from a tags parameter.
// Handles both JSON-encoded string form and direct []interface{} from JSON body.
func parseTags(params map[string]interface{}) map[string]string {
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
			return result
		case string:
			if v == "" {
				return nil
			}
			var tagList []map[string]interface{}
			if err := json.Unmarshal([]byte(v), &tagList); err != nil {
				return nil
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
	return result
}
