package cloudwatchlogs

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// CreateExportTask creates a task to export log events to an S3 bucket.
func (s *LogsService) CreateExportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	taskId, err := s.createExportTaskCore(store, &CreateExportTaskInput{
		TaskName:            request.GetParamLowerFirst(req.Parameters, "TaskName"),
		LogGroupName:        request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		Destination:         request.GetParamLowerFirst(req.Parameters, "Destination"),
		From:                request.GetIntParam(req.Parameters, "From"),
		To:                  request.GetIntParam(req.Parameters, "To"),
		LogStreamNamePrefix: request.GetParamLowerFirst(req.Parameters, "LogStreamNamePrefix"),
		DestinationPrefix:   request.GetParamLowerFirst(req.Parameters, "DestinationPrefix"),
		Region:              reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"taskId": taskId,
	}, nil
}

func (s *LogsService) executeExportTask(region, logGroupName, streamPrefix string, fromTime, toTime int64, bucket, prefix, taskId string) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("PANIC in export task",
				logs.String("taskId", taskId),
				logs.Any("panic", r))
			s.updateExportTaskStatus(region, taskId, "FAILED", fmt.Sprintf("panic: %v", r))
		}
	}()

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		s.updateExportTaskStatus(region, taskId, "FAILED", fmt.Sprintf("store error: %v", err))
		return
	}

	streams, _, _ := store.ListLogStreams(logGroupName, streamPrefix, "", 1000)

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)

	for _, ls := range streams {
		events, _, _, err := store.GetLogEvents(logGroupName, ls.Name, fromTime, toTime, 10000, true, "")
		if err != nil {
			continue
		}
		for _, evt := range events {
			record := map[string]interface{}{
				"timestamp":     evt.Timestamp,
				"message":       evt.Message,
				"ingestionTime": evt.IngestionTime,
				"logStream":     ls.Name,
			}
			line, _ := json.Marshal(record)
			gw.Write(line)
			gw.Write([]byte("\n"))
		}
	}

	if err := gw.Close(); err != nil {
		s.updateExportTaskStatus(region, taskId, "FAILED", fmt.Sprintf("gzip error: %v", err))
		return
	}

	s3Key := taskId + "/exportedlogs.gz"
	if prefix != "" {
		s3Key = prefix + "/" + s3Key
	}

	if s.bus != nil {
		s3Invoker := s.bus.S3Invoker()
		if s3Invoker != nil {
			if err := s3Invoker.PutObject(context.Background(), region, bucket, s3Key, buf.Bytes(), "application/x-gzip"); err != nil {
				s.updateExportTaskStatus(region, taskId, "FAILED", fmt.Sprintf("S3 upload error: %v", err))
				return
			}
		}
	}

	s.updateExportTaskStatus(region, taskId, "COMPLETED", "")
}

func (s *LogsService) updateExportTaskStatus(region, taskId, status, message string) {
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return
	}
	task, err := store.GetExportTask(taskId)
	if err != nil {
		return
	}
	task.Status = status
	task.StatusMessage = message
	if status == "COMPLETED" || status == "FAILED" || status == "CANCELLED" {
		if task.ExecutionInfo == nil {
			task.ExecutionInfo = make(map[string]interface{})
		}
		task.ExecutionInfo["completionTime"] = time.Now().UTC().UnixMilli()
	}
	if err := store.PutExportTask(task); err != nil {
		logs.Error("Failed to persist export task status update",
			logs.String("taskId", task.TaskId),
			logs.String("status", status),
			logs.Err(err))
	}
}

// DescribeExportTasks lists export tasks.
func (s *LogsService) DescribeExportTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit, err := validateListLimit(int32(request.GetIntParam(req.Parameters, "Limit")), 50, 50)
	if err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	items, nextMarker, err := s.describeExportTasksCore(store, &DescribeExportTasksInput{
		TaskId:     request.GetParamLowerFirst(req.Parameters, "TaskId"),
		StatusCode: request.GetParamLowerFirst(req.Parameters, "StatusCode"),
		NextToken:  request.GetParamLowerFirst(req.Parameters, "NextToken"),
		Limit:      limit,
	})
	if err != nil {
		return nil, err
	}

	tasks := make([]map[string]interface{}, len(items))
	for i, t := range items {
		tasks[i] = formatExportTask(t)
	}

	resp := map[string]interface{}{
		"exportTasks": tasks,
	}
	if nextMarker != "" {
		resp["nextToken"] = nextMarker
	}

	return resp, nil
}

// CancelExportTask cancels a running export task.
func (s *LogsService) CancelExportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	if err := s.cancelExportTaskCore(store, request.GetParamLowerFirst(req.Parameters, "TaskId")); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func formatExportTask(t *logsstore.ExportTask) map[string]interface{} {
	result := map[string]interface{}{
		"taskId":       t.TaskId,
		"taskName":     t.TaskName,
		"logGroupName": t.LogGroupName,
		"from":         t.From,
		"to":           t.To,
		"destination":  t.Destination,
		"status":       map[string]interface{}{"code": t.Status, "message": t.StatusMessage},
		"creationTime": t.CreationTime,
	}
	if t.LogStreamNamePrefix != "" {
		result["logStreamNamePrefix"] = t.LogStreamNamePrefix
	}
	if t.DestinationPrefix != "" {
		result["destinationPrefix"] = t.DestinationPrefix
	}
	execInfo := map[string]interface{}{}
	if t.CreationTime > 0 {
		execInfo["creationTime"] = t.CreationTime
	}
	if t.ExecutionInfo != nil {
		if ct, ok := t.ExecutionInfo["completionTime"]; ok {
			execInfo["completionTime"] = ct
		}
	}
	result["executionInfo"] = execInfo
	return result
}
