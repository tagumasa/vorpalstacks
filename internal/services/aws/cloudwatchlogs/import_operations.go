package cloudwatchlogs

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/utils/aws/arn"
)

// CreateImportTask creates a task to import log events from an S3 source.
func (s *LogsService) CreateImportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	importSourceArn := request.GetParamLowerFirst(req.Parameters, "ImportSourceArn")

	if importSourceArn == "" {
		return nil, ErrMissingParameter
	}

	parsedArn, err := arn.ParseARN(importSourceArn)
	if err != nil || parsedArn.Service != "s3" {
		return nil, NewLogsError("InvalidParameterException",
			"importSourceArn must be a valid S3 ARN", 400)
	}

	bucket := parsedArn.AccountID
	s3Key := parsedArn.Resource
	if idx := strings.Index(s3Key, "/"); idx >= 0 {
		bucket = s3Key[:idx]
		s3Key = s3Key[idx+1:]
	} else {
		bucket = s3Key
		s3Key = ""
	}

	if bucket == "" {
		return nil, NewLogsError("InvalidParameterException",
			"importSourceArn must contain an S3 bucket", 400)
	}

	var importFilter map[string]interface{}
	if filter, ok := req.Parameters["importFilter"]; ok {
		if m, ok := filter.(map[string]interface{}); ok {
			importFilter = m
		}
	}

	managedLogGroupName := deriveManagedLogGroupName(importSourceArn)
	region := reqCtx.GetRegion()
	importDestinationArn := fmt.Sprintf("arn:aws:logs:%s:%s:log-group:%s", region, s.accountID, managedLogGroupName)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetLogGroup(managedLogGroupName); err != nil {
		lg := logsstore.NewLogGroup(managedLogGroupName, region, s.accountID)
		if err := store.CreateLogGroup(lg); err != nil {
			if existing, gErr := store.GetLogGroup(managedLogGroupName); gErr != nil || existing == nil {
				return nil, mapStoreError(err)
			}
		}
	}

	importId := fmt.Sprintf("import-%d", time.Now().UnixNano())
	roleArn := request.GetParamLowerFirst(req.Parameters, "ImportRoleArn")

	task := &logsstore.ImportTask{
		ImportId:             importId,
		ImportSourceArn:      importSourceArn,
		ImportRoleArn:        roleArn,
		LogGroupName:         managedLogGroupName,
		ImportStatus:         "IN_PROGRESS",
		ImportDestinationArn: importDestinationArn,
		ImportFilter:         importFilter,
		CreationTime:         time.Now().UTC().UnixMilli(),
	}

	if err := store.PutImportTask(task); err != nil {
		return nil, mapStoreError(err)
	}

	go s.executeImportTask(region, importId, bucket, s3Key, managedLogGroupName)

	return map[string]interface{}{
		"importId":             importId,
		"importDestinationArn": importDestinationArn,
		"creationTime":         task.CreationTime,
	}, nil
}

func (s *LogsService) executeImportTask(region, importId, bucket, s3Key, logGroupName string) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("PANIC in import task",
				logs.String("importId", importId),
				logs.Any("panic", r))
			s.updateImportTaskStatus(region, importId, "FAILED", fmt.Sprintf("panic: %v", r))
		}
	}()

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		s.updateImportTaskStatus(region, importId, "FAILED", fmt.Sprintf("store error: %v", err))
		return
	}

	if s.bus == nil {
		s.updateImportTaskStatus(region, importId, "FAILED", "event bus not configured")
		return
	}

	s3Invoker := s.bus.S3Invoker()
	if s3Invoker == nil {
		s.updateImportTaskStatus(region, importId, "FAILED", "S3 invoker not configured")
		return
	}

	keys := []string{s3Key}
	if s3Key == "" {
		listed, err := s3Invoker.ListObjects(context.Background(), region, bucket, "", 1000)
		if err != nil {
			s.updateImportTaskStatus(region, importId, "FAILED", fmt.Sprintf("S3 list error: %v", err))
			return
		}
		keys = listed
	}

	totalBytes := int64(0)
	for _, key := range keys {
		if key == "" {
			continue
		}
		data, err := s3Invoker.GetObject(context.Background(), region, bucket, key, 100*1024*1024)
		if err != nil {
			logs.Warn("Failed to get S3 object during import",
				logs.String("bucket", bucket),
				logs.String("key", key),
				logs.Err(err))
			continue
		}

		content, err := decompressIfNeeded(data)
		if err != nil {
			logs.Warn("Failed to decompress S3 object during import",
				logs.String("key", key),
				logs.Err(err))
			continue
		}

		events := parseImportContent(content)
		if len(events) == 0 {
			continue
		}

		streamName := fmt.Sprintf("import-%s-stream", importId)
		ls := logsstore.NewLogStream(streamName, logGroupName)
		if err := store.CreateLogStream(ls); err != nil {
			logs.Warn("Failed to create import log stream",
				logs.String("importId", importId),
				logs.String("stream", streamName),
				logs.Err(err))
			continue
		}

		_, err = store.PutLogEvents(logGroupName, streamName, events)
		if err != nil {
			logs.Warn("Failed to put imported log events",
				logs.String("importId", importId),
				logs.Err(err))
			continue
		}
		totalBytes += int64(len(data))
	}

	s.updateImportTaskStatus(region, importId, "COMPLETED", "")
	s.updateImportStats(region, importId, totalBytes)
}

func decompressIfNeeded(data []byte) ([]byte, error) {
	if len(data) >= 2 && data[0] == 0x1f && data[1] == 0x8b {
		gr, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		defer gr.Close()
		return io.ReadAll(gr)
	}
	return data, nil
}

func parseImportContent(content []byte) []logsstore.LogEntry {
	var events []logsstore.LogEntry
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var record map[string]interface{}
		if json.Unmarshal([]byte(line), &record) == nil {
			ts := int64(0)
			if t, ok := record["timestamp"].(float64); ok {
				ts = int64(t)
			}
			msg := ""
			if m, ok := record["message"].(string); ok {
				msg = m
			} else {
				msg = line
			}
			if ts == 0 {
				ts = time.Now().UnixMilli()
			}
			events = append(events, logsstore.LogEntry{Timestamp: ts, Message: msg})
		} else {
			events = append(events, logsstore.LogEntry{
				Timestamp: time.Now().UnixMilli(),
				Message:   line,
			})
		}
	}
	return events
}

func (s *LogsService) updateImportTaskStatus(region, importId, status, message string) {
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return
	}
	task, err := store.GetImportTask(importId)
	if err != nil {
		return
	}
	task.ImportStatus = status
	if message != "" {
		task.ErrorMessage = message
	}
	task.LastUpdatedTime = time.Now().UTC().UnixMilli()
	if err := store.PutImportTask(task); err != nil {
		logs.Error("Failed to persist import task status update",
			logs.String("importId", importId),
			logs.String("status", status),
			logs.Err(err))
	}
}

func (s *LogsService) updateImportStats(region, importId string, bytesImported int64) {
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return
	}
	task, err := store.GetImportTask(importId)
	if err != nil {
		return
	}
	task.ImportStatistics = map[string]interface{}{
		"bytesImported": bytesImported,
	}
	if err := store.PutImportTask(task); err != nil {
		logs.Error("Failed to persist import task statistics",
			logs.String("importId", importId),
			logs.Err(err))
	}
}

// DescribeImportTasks lists import tasks.
func (s *LogsService) DescribeImportTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	importId := request.GetParamLowerFirst(req.Parameters, "ImportId")
	status := request.GetParamLowerFirst(req.Parameters, "ImportStatus")
	sourceArn := request.GetParamLowerFirst(req.Parameters, "ImportSourceArn")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	limit, err := validateListLimit(int32(request.GetIntParam(req.Parameters, "Limit")), 50, 50)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	allTasks, err := store.ListImportTasks(importId, status, sourceArn)
	if err != nil {
		return nil, mapStoreError(err)
	}

	result := pagination.PaginateSlice(allTasks, nextToken, int(limit), func(t *logsstore.ImportTask) string {
		return t.ImportId
	})

	tasks := make([]map[string]interface{}, len(result.Items))
	for i, t := range result.Items {
		tasks[i] = formatImportTask(t)
	}

	resp := map[string]interface{}{
		"imports": tasks,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
	}

	return resp, nil
}

// CancelImportTask cancels a running import task.
func (s *LogsService) CancelImportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	importId := request.GetParamLowerFirst(req.Parameters, "ImportId")
	if importId == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	task, err := store.GetImportTask(importId)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if task.ImportStatus == "COMPLETED" || task.ImportStatus == "FAILED" || task.ImportStatus == "CANCELLED" {
		return nil, NewLogsError("InvalidOperationException",
			fmt.Sprintf("Cannot cancel import task in %s state", task.ImportStatus), 400)
	}

	task.ImportStatus = "CANCELLED"
	task.LastUpdatedTime = time.Now().UTC().UnixMilli()
	if err := store.PutImportTask(task); err != nil {
		return nil, mapStoreError(err)
	}

	resp := map[string]interface{}{
		"importId":        importId,
		"importStatus":    "CANCELLED",
		"creationTime":    task.CreationTime,
		"lastUpdatedTime": task.LastUpdatedTime,
	}
	if task.ImportStatistics != nil {
		resp["importStatistics"] = task.ImportStatistics
	}
	return resp, nil
}

// DescribeImportTaskBatches lists import task batches.
func (s *LogsService) DescribeImportTaskBatches(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	importId := request.GetParamLowerFirst(req.Parameters, "ImportId")
	if importId == "" {
		return nil, ErrMissingParameter
	}

	// batchImportStatus is accepted but batch-level status tracking is not
	// implemented; the response always returns an empty importBatches list.

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	task, err := store.GetImportTask(importId)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"importId":        importId,
		"importSourceArn": task.ImportSourceArn,
		"importBatches":   []interface{}{},
	}, nil
}

func formatImportTask(t *logsstore.ImportTask) map[string]interface{} {
	result := map[string]interface{}{
		"importId":             t.ImportId,
		"importSourceArn":      t.ImportSourceArn,
		"importStatus":         t.ImportStatus,
		"importDestinationArn": t.ImportDestinationArn,
		"creationTime":         t.CreationTime,
		"lastUpdatedTime":      t.LastUpdatedTime,
	}
	if t.ImportStatistics != nil {
		result["importStatistics"] = t.ImportStatistics
	}
	if t.ImportFilter != nil {
		result["importFilter"] = t.ImportFilter
	}
	if t.ErrorMessage != "" {
		result["errorMessage"] = t.ErrorMessage
	}
	return result
}

func deriveManagedLogGroupName(importSourceArn string) string {
	edsId := importSourceArn
	if idx := strings.LastIndex(edsId, "/"); idx >= 0 {
		edsId = edsId[idx+1:]
	}
	return "/aws/imported/cloudtrail-lake/" + edsId
}
