package neptunedata

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
)

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

	s.mu.Lock()
	s.loaderJobs[loadId] = &loaderJobState{
		LoadId:    loadId,
		Status:    "LOAD_COMPLETED",
		Source:    params.Source,
		Format:    params.Format,
		StartedAt: time.Now(),
	}
	s.mu.Unlock()

	return map[string]interface{}{
		"status": "200",
		"payload": map[string]interface{}{
			"loadId": loadId,
		},
	}, nil
}

func (s *NeptuneDataService) GetLoaderJobStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	loadId := getPathParam(req, "loadId")
	if loadId == "" {
		return nil, missingParameter("loadId")
	}
	details := request.GetBoolParam(req.Parameters, "details")
	errors := request.GetBoolParam(req.Parameters, "errors")
	errorsPerPage := request.GetIntParam(req.Parameters, "errorsPerPage")
	page := request.GetIntParam(req.Parameters, "page")

	_ = details
	_ = errors
	_ = errorsPerPage
	_ = page

	s.mu.RLock()
	job, ok := s.loaderJobs[loadId]
	s.mu.RUnlock()

	if !ok {
		return nil, bulkLoadNotFound(loadId)
	}

	payload := map[string]interface{}{
		"loadId":        job.LoadId,
		"status":        job.Status,
		"feedCount":     map[string]interface{}{"total": 0, "succeeded": 0, "failed": 0},
		"overallStatus": map[string]interface{}{"totalRecords": 0, "totalTime": 0},
	}

	return map[string]interface{}{
		"status":  "200",
		"payload": payload,
	}, nil
}

func (s *NeptuneDataService) ListLoaderJobs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	_ = request.GetBoolParam(req.Parameters, "includeQueuedLoads")
	_ = request.GetIntParam(req.Parameters, "limit")

	s.mu.RLock()
	loadIds := make([]string, 0, len(s.loaderJobs))
	for _, job := range s.loaderJobs {
		loadIds = append(loadIds, job.LoadId)
	}
	s.mu.RUnlock()

	return map[string]interface{}{
		"status": "200",
		"payload": map[string]interface{}{
			"loadIds": loadIds,
		},
	}, nil
}

func (s *NeptuneDataService) CancelLoaderJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	loadId := getPathParam(req, "loadId")
	if loadId == "" {
		return nil, missingParameter("loadId")
	}

	s.mu.Lock()
	job, ok := s.loaderJobs[loadId]
	if !ok {
		s.mu.Unlock()
		return nil, bulkLoadNotFound(loadId)
	}
	job.Status = "CANCELLED"
	job.FinishedAt = time.Now()
	s.mu.Unlock()

	return map[string]interface{}{
		"status": "200",
	}, nil
}
