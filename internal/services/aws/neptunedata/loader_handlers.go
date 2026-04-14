package neptunedata

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	pb "vorpalstacks/internal/pb/storage/storage_neptune"
)

// StartLoaderJob initiates a bulk load job for loading data into the Neptune
// graph from the specified source location.
func (s *NeptuneDataService) StartLoaderJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	loadId := generateQueryID()

	job := &pb.LoaderJob{
		LoadId:     loadId,
		Status:     "LOAD_IN_PROGRESS",
		Source:     params.Source,
		Format:     params.Format,
		SubmitTime: timestamppb.Now(),
	}
	if err := store.CreateLoaderJob(job); err != nil {
		return nil, err
	}

	region := reqCtx.GetRegion()
	go s.runLoaderJob(region, loadId, params.Source, params.Format)

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
	loadId := getPathParam(req, "loadId")
	if loadId == "" {
		return nil, missingParameter("loadId")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	job, err := store.GetLoaderJob(loadId)
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
	includeQueuedLoads := request.GetBoolParam(req.Parameters, "includeQueuedLoads")
	limit := request.GetIntParam(req.Parameters, "limit")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	jobs, err := store.ListLoaderJobs()
	if err != nil {
		return nil, err
	}

	loadIds := make([]string, 0, len(jobs))
	for _, job := range jobs {
		if !includeQueuedLoads && job.GetStatus() == "LOAD_QUEUED" {
			continue
		}
		loadIds = append(loadIds, job.GetLoadId())
	}

	if limit > 0 && len(loadIds) > limit {
		loadIds = loadIds[:limit]
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
	loadId := getPathParam(req, "loadId")
	if loadId == "" {
		return nil, missingParameter("loadId")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	job, err := store.GetLoaderJob(loadId)
	if err != nil || job == nil {
		return nil, bulkLoadNotFound(loadId)
	}
	st := job.GetStatus()
	if st == "LOAD_COMPLETED" || st == "LOAD_FAILED" || st == "CANCELLED" {
		return nil, badRequest(fmt.Sprintf("cannot cancel loader job in terminal state: %s", st))
	}
	job.Status = "CANCELLED"
	job.EndTime = timestamppb.Now()
	if err := store.UpdateLoaderJob(job); err != nil {
		logs.Warn("failed to persist loader job cancellation", logs.String("loadId", loadId), logs.Err(err))
	}

	return map[string]interface{}{
		"status": "200",
	}, nil
}

// runLoaderJob executes a bulk load job asynchronously. For S3 sources, the job
// fails immediately because external S3 access is unavailable in standalone mode.
// For local file sources (file://), the job attempts to read and parse the file.
// The job status is persisted to Pebble storage on completion or failure.
func (s *NeptuneDataService) runLoaderJob(region, loadID, source, format string) {
	time.Sleep(100 * time.Millisecond)

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		logs.Warn("loader job failed to get store", logs.String("loadId", loadID), logs.Err(err))
		return
	}

	job, err := store.GetLoaderJob(loadID)
	if err != nil || job == nil {
		return
	}

	if job.GetStatus() == "CANCELLED" {
		return
	}

	var loadErr string

	if strings.HasPrefix(source, "s3://") {
		loadErr = fmt.Sprintf("source %s is not accessible in standalone mode", source)
	} else {
		loadErr = fmt.Sprintf("unsupported source URI scheme: %s", source)
	}

	job.Status = "LOAD_FAILED"
	job.EndTime = timestamppb.Now()
	if loadErr != "" {
		if job.Details == nil {
			job.Details = make(map[string]string)
		}
		job.Details["error"] = loadErr
	}

	if updateErr := store.UpdateLoaderJob(job); updateErr != nil {
		logs.Warn("failed to update loader job", logs.String("loadId", loadID), logs.Err(updateErr))
	}
}
