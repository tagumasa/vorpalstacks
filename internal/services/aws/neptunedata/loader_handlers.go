package neptunedata

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "vorpalstacks/internal/pb/storage/storage_neptune"
	"vorpalstacks/internal/services/aws/common/request"
)

// StartLoaderJob initiates a bulk load job for loading data into the Neptune
// graph from the specified source location.
func (s *NeptuneDataService) StartLoaderJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	body := req.Body
	var params struct {
		Source      string `json:"source"`
		Format      string `json:"format"`
		IamRoleArn  string `json:"iamRoleArn"`
		Mode        string `json:"mode"`
		Parallelism string `json:"parallelism"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}
	if params.Source == "" {
		return nil, missingParameter("source")
	}
	if params.Format == "" {
		return nil, missingParameter("format")
	}

	loadId := generateQueryID()

	job := &pb.LoaderJob{
		LoadId:     loadId,
		Status:     "LOAD_COMPLETED",
		Source:     params.Source,
		Format:     params.Format,
		SubmitTime: timestamppb.Now(),
	}
	if err := s.store.CreateLoaderJob(job); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"status": "200",
		"payload": map[string]interface{}{
			"loadId": loadId,
		},
	}, nil
}

// GetLoaderJobStatus returns the current status and statistics of a bulk load
// job.
func (s *NeptuneDataService) GetLoaderJobStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	loadId := getPathParam(req, "loadId")
	if loadId == "" {
		return nil, missingParameter("loadId")
	}

	job, err := s.store.GetLoaderJob(loadId)
	if err != nil || job == nil {
		return nil, bulkLoadNotFound(loadId)
	}

	payload := map[string]interface{}{
		"loadId":        job.GetLoadId(),
		"status":        job.GetStatus(),
		"feedCount":     map[string]interface{}{"total": 0, "succeeded": 0, "failed": 0},
		"overallStatus": map[string]interface{}{"totalRecords": 0, "totalTime": 0},
	}

	return map[string]interface{}{
		"status":  "200",
		"payload": payload,
	}, nil
}

// ListLoaderJobs returns the load IDs of all submitted bulk load jobs,
// optionally including queued loads.
func (s *NeptuneDataService) ListLoaderJobs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	_ = request.GetBoolParam(req.Parameters, "includeQueuedLoads")
	_ = request.GetIntParam(req.Parameters, "limit")

	jobs, err := s.store.ListLoaderJobs()
	if err != nil {
		return nil, err
	}

	loadIds := make([]string, 0, len(jobs))
	for _, job := range jobs {
		loadIds = append(loadIds, job.GetLoadId())
	}

	return map[string]interface{}{
		"status": "200",
		"payload": map[string]interface{}{
			"loadIds": loadIds,
		},
	}, nil
}

// CancelLoaderJob cancels a running or queued bulk load job and marks its
// status as cancelled.
func (s *NeptuneDataService) CancelLoaderJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	loadId := getPathParam(req, "loadId")
	if loadId == "" {
		return nil, missingParameter("loadId")
	}

	job, err := s.store.GetLoaderJob(loadId)
	if err != nil || job == nil {
		return nil, bulkLoadNotFound(loadId)
	}
	job.Status = "CANCELLED"
	job.EndTime = timestamppb.Now()
	_ = s.store.UpdateLoaderJob(job)

	return map[string]interface{}{
		"status": "200",
	}, nil
}
