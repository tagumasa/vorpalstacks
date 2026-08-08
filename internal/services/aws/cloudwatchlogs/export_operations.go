package cloudwatchlogs

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// CreateExportTask creates a task to export log events to an S3 bucket.
func (s *LogsService) CreateExportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskName := request.GetParamLowerFirst(req.Parameters, "TaskName")
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	destination := request.GetParamLowerFirst(req.Parameters, "Destination")
	fromTime := request.GetIntParam(req.Parameters, "From")
	toTime := request.GetIntParam(req.Parameters, "To")

	if err := validateLogGroupName(logGroupName); err != nil {
		return nil, err
	}
	if err := validateExportDestinationBucket(destination); err != nil {
		return nil, err
	}

	if fromTime <= 0 || toTime <= 0 || fromTime >= toTime {
		return nil, NewLogsError("InvalidParameterException",
			"from and to must be positive timestamps with from < to", 400)
	}

	logStreamNamePrefix := request.GetParamLowerFirst(req.Parameters, "LogStreamNamePrefix")
	destinationPrefix := request.GetParamLowerFirst(req.Parameters, "DestinationPrefix")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetLogGroup(logGroupName)
	if err != nil {
		return nil, mapStoreError(err)
	}

	taskId := fmt.Sprintf("export-%d", time.Now().UnixNano())
	region := reqCtx.GetRegion()

	task := &logsstore.ExportTask{
		TaskId:              taskId,
		TaskName:            taskName,
		LogGroupName:        logGroupName,
		LogStreamNamePrefix: logStreamNamePrefix,
		From:                int64(fromTime),
		To:                  int64(toTime),
		Destination:         destination,
		DestinationPrefix:   destinationPrefix,
		Status:              "RUNNING",
		CreationTime:        time.Now().UTC().UnixMilli(),
	}

	if err := store.PutExportTask(task); err != nil {
		return nil, mapStoreError(err)
	}

	go s.executeExportTask(region, logGroupName, logStreamNamePrefix, int64(fromTime), int64(toTime), destination, destinationPrefix, taskId)

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
	taskId := request.GetParamLowerFirst(req.Parameters, "TaskId")
	statusCode := request.GetParamLowerFirst(req.Parameters, "StatusCode")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	limit, err := validateListLimit(int32(request.GetIntParam(req.Parameters, "Limit")), 50, 50)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if taskId != "" {
		task, err := store.GetExportTask(taskId)
		if err != nil {
			return nil, mapStoreError(err)
		}
		return map[string]interface{}{
			"exportTasks": []map[string]interface{}{formatExportTask(task)},
		}, nil
	}

	allTasks, err := store.ListExportTasks(statusCode)
	if err != nil {
		return nil, mapStoreError(err)
	}

	result := pagination.PaginateSlice(allTasks, nextToken, int(limit), func(t *logsstore.ExportTask) string {
		return t.TaskId
	})

	tasks := make([]map[string]interface{}, len(result.Items))
	for i, t := range result.Items {
		tasks[i] = formatExportTask(t)
	}

	resp := map[string]interface{}{
		"exportTasks": tasks,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
	}

	return resp, nil
}

// CancelExportTask cancels a running export task.
func (s *LogsService) CancelExportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamLowerFirst(req.Parameters, "TaskId")
	if taskId == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	task, err := store.GetExportTask(taskId)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if task.Status == "COMPLETED" || task.Status == "FAILED" || task.Status == "CANCELLED" {
		return nil, NewLogsError("InvalidOperationException",
			fmt.Sprintf("Cannot cancel export task in %s state", task.Status), 400)
	}

	task.Status = "CANCELLED"
	task.StatusMessage = "Cancelled by user"
	if err := store.PutExportTask(task); err != nil {
		return nil, mapStoreError(err)
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
