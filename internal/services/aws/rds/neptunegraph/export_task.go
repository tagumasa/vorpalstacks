package neptunegraph

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
)

// StartExportTask initiates a bulk export of graph data to the specified S3 destination.
func (s *NeptuneGraphService) StartExportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &StartExportTaskInput{
		GraphIdentifier:  request.GetStringParam(req.Parameters, "graphIdentifier"),
		Format:           strings.ToUpper(request.GetStringParam(req.Parameters, "format")),
		KmsKeyIdentifier: request.GetStringParam(req.Parameters, "kmsKeyIdentifier"),
		RoleArn:          request.GetStringParam(req.Parameters, "roleArn"),
		Destination:      request.GetStringParam(req.Parameters, "destination"),
		ParquetType:      strings.ToUpper(request.GetStringParam(req.Parameters, "parquetType")),
		HasExportFilter:  request.HasParam(req.Parameters, "exportFilter"),
		ExportFilter:     parseExportFilter(req.Parameters),
	}

	task, err := s.startExportTaskCore(store, in)
	if err != nil {
		return nil, err
	}
	return exportTaskToResponse(task), nil
}

// GetExportTask retrieves an export task by its identifier.
func (s *NeptuneGraphService) GetExportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	task, err := s.getExportTaskCore(store, request.GetStringParam(req.Parameters, "taskIdentifier"))
	if err != nil {
		return nil, err
	}

	return exportTaskToResponse(task), nil
}

// ListExportTasks returns a paginated list of export task summaries, optionally filtered by graph.
func (s *NeptuneGraphService) ListExportTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &ListExportTasksInput{
		MaxItems:        clampMaxResults(request.GetIntParam(req.Parameters, "maxResults")),
		Marker:          request.GetStringParam(req.Parameters, "nextToken"),
		GraphIdentifier: request.GetStringParam(req.Parameters, "graphIdentifier"),
	}

	res, err := s.listExportTasksCore(store, in)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(res.Tasks))
	for _, t := range res.Tasks {
		items = append(items, exportTaskSummaryToResponse(t))
	}

	result := map[string]interface{}{
		"tasks": items,
	}
	if res.Truncated {
		result["nextToken"] = res.NextToken
	}
	return result, nil
}

// CancelExportTask cancels an in-progress export task, transitioning it to CANCELLED state.
func (s *NeptuneGraphService) CancelExportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &CancelExportTaskInput{
		TaskIdentifier: request.GetStringParam(req.Parameters, "taskIdentifier"),
	}

	task, err := s.cancelExportTaskCore(store, in)
	if err != nil {
		return nil, err
	}
	return exportTaskSummaryToResponse(task), nil
}
