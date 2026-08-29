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

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// CreateImportTask creates a task to import log events from an S3 source.
func (s *LogsService) CreateImportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	importSourceArn := request.GetParamLowerFirst(req.Parameters, "ImportSourceArn")

	var importFilter map[string]interface{}
	if filter, ok := req.Parameters["importFilter"]; ok {
		if m, ok := filter.(map[string]interface{}); ok {
			importFilter = m
		}
	}

	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	task, err := s.createImportTaskCore(store, &CreateImportTaskInput{
		ImportSourceArn: importSourceArn,
		ImportRoleArn:   request.GetParamLowerFirst(req.Parameters, "ImportRoleArn"),
		ImportFilter:    importFilter,
		Region:          reqCtx.GetRegion(),
		AccountID:       s.accountID,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"importId":             task.ImportId,
		"importDestinationArn": task.ImportDestinationArn,
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
	limit, err := validateListLimit(int32(request.GetIntParam(req.Parameters, "Limit")), 50, 50)
	if err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	items, nextMarker, err := s.describeImportTasksCore(store,
		request.GetParamLowerFirst(req.Parameters, "ImportId"),
		request.GetParamLowerFirst(req.Parameters, "ImportStatus"),
		request.GetParamLowerFirst(req.Parameters, "ImportSourceArn"),
		request.GetParamLowerFirst(req.Parameters, "NextToken"),
		limit)
	if err != nil {
		return nil, err
	}

	tasks := make([]map[string]interface{}, len(items))
	for i, t := range items {
		tasks[i] = formatImportTask(t)
	}

	resp := map[string]interface{}{
		"imports": tasks,
	}
	if nextMarker != "" {
		resp["nextToken"] = nextMarker
	}

	return resp, nil
}

// CancelImportTask cancels a running import task.
func (s *LogsService) CancelImportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	importId := request.GetParamLowerFirst(req.Parameters, "ImportId")

	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	task, err := s.cancelImportTaskCore(store, importId)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"importId":     importId,
		"importStatus": "CANCELLED",
	}
	if task != nil {
		resp["creationTime"] = task.CreationTime
		resp["lastUpdatedTime"] = task.LastUpdatedTime
		if task.ImportStatistics != nil {
			resp["importStatistics"] = task.ImportStatistics
		}
	}
	return resp, nil
}

// DescribeImportTaskBatches lists import task batches.
func (s *LogsService) DescribeImportTaskBatches(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// batchImportStatus is accepted but batch-level status tracking is not
	// implemented; the response always returns an empty importBatches list.

	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	task, err := s.describeImportTaskBatchesCore(store, request.GetParamLowerFirst(req.Parameters, "ImportId"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"importId":        task.ImportId,
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

// toInt64 attempts to convert an interface{} (float64 from JSON or json.Number)
// to an int64 value, returning false if the conversion fails.
func toInt64(v interface{}) (int64, bool) {
	switch n := v.(type) {
	case float64:
		return int64(n), true
	case int64:
		return n, true
	case int:
		return int64(n), true
	}
	return 0, false
}
