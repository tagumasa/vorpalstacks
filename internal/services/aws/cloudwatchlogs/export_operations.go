package cloudwatchlogs

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
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

// exportTaskRuntime is the documented export deadline: "Export tasks time
// out after 24 hours." (Exporting log data to Amazon S3, CloudWatch Logs
// User Guide). A variable so the executor's deadline observatory (the
// tests) can shorten it; the value itself is the store constant.
var exportTaskRuntime = logsstore.MaxExportTaskRuntime

// exportTaskTimeoutMessage is the statusMessage a timed-out task carries.
const exportTaskTimeoutMessage = "Export task timed out after 24 hours"

// expireTimedOutExportTasks transitions the account's RUNNING and PENDING
// export tasks past the documented 24-hour deadline to FAILED — "Export
// tasks time out after 24 hours." (Exporting log data to Amazon S3,
// CloudWatch Logs User Guide). A task whose executor died with the
// process (a restart orphans the RUNNING record) must not report RUNNING
// forever or hold its Region's single active-task slot.
func (s *LogsService) expireTimedOutExportTasks() {
	if s.storageManager == nil {
		return
	}
	cutoff := time.Now().UTC().UnixMilli() - int64(exportTaskRuntime/time.Millisecond)
	for _, region := range s.storageManager.GetActiveRegions() {
		store, err := s.getLogsStoreByRegion(region)
		if err != nil {
			continue
		}
		tasks, err := store.ListExportTasks("")
		if err != nil {
			continue
		}
		for _, t := range tasks {
			if t.Status != logsstore.ExportStatusRunning && t.Status != logsstore.ExportStatusPending {
				continue
			}
			if t.CreationTime == 0 || t.CreationTime > cutoff {
				continue
			}
			s.updateExportTaskStatus(region, t.TaskId, logsstore.ExportStatusFailed, exportTaskTimeoutMessage)
		}
	}
}

func (s *LogsService) executeExportTask(ctx context.Context, region, logGroupName, streamPrefix string, fromTime, toTime int64, bucket, prefix, taskId string) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("PANIC in export task",
				logs.String("taskId", taskId),
				logs.Any("panic", r))
			s.updateExportTaskStatus(region, taskId, logsstore.ExportStatusFailed, fmt.Sprintf("panic: %v", r))
		}
	}()

	started := time.Now()
	deadlineExceeded := func() bool {
		return time.Since(started) > exportTaskRuntime
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		s.updateExportTaskStatus(region, taskId, logsstore.ExportStatusFailed, fmt.Sprintf("store error: %v", err))
		return
	}

	if s.eventBus() == nil || s.eventBus().S3Invoker() == nil {
		// An export that cannot deliver is a failed export, not a silent
		// no-op that reports COMPLETED.
		s.updateExportTaskStatus(region, taskId, logsstore.ExportStatusFailed, "S3 invoker not configured")
		return
	}
	s3Invoker := s.eventBus().S3Invoker()

	streams, err := fetchAllLogStreams(store, logGroupName, streamPrefix)
	if err != nil {
		s.updateExportTaskStatus(region, taskId, logsstore.ExportStatusFailed, fmt.Sprintf("stream listing error: %v", err))
		return
	}

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)

	for _, ls := range streams {
		// "Export tasks time out after 24 hours." — the deadline holds
		// across the read legs too: a window whose reads outgrow it fails
		// the task instead of running without bound.
		if deadlineExceeded() {
			s.updateExportTaskStatus(region, taskId, logsstore.ExportStatusFailed, exportTaskTimeoutMessage)
			return
		}
		// Cancellation checkpoint between streams: a task the canceller
		// moved to PENDING_CANCEL stops reading and takes its terminal
		// CANCELLED here.
		if task, err := store.GetExportTask(taskId); err == nil && task.Status == logsstore.ExportStatusPendingCancel {
			s.updateExportTaskStatus(region, taskId, logsstore.ExportStatusCancelled, "Cancelled by user")
			return
		}
		events, err := fetchAllLogEventsEndInclusive(store, logGroupName, ls.Name, fromTime, toTime)
		if err != nil {
			// A stream that cannot be read is a failed export, not a
			// silently smaller one: the delivered object would omit the
			// stream's data while the task reports COMPLETED.
			s.updateExportTaskStatus(region, taskId, logsstore.ExportStatusFailed,
				fmt.Sprintf("stream %s read error: %v", ls.Name, err))
			return
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
		s.updateExportTaskStatus(region, taskId, logsstore.ExportStatusFailed, fmt.Sprintf("gzip error: %v", err))
		return
	}
	if deadlineExceeded() {
		s.updateExportTaskStatus(region, taskId, logsstore.ExportStatusFailed, exportTaskTimeoutMessage)
		return
	}

	// Cancellation checkpoint: a task the canceller already moved to a
	// terminal state keeps it, and a PENDING_CANCEL task takes its
	// terminal CANCELLED here — either way the upload is skipped and the
	// finaliser below refuses to overwrite a terminal status.
	if task, err := store.GetExportTask(taskId); err == nil {
		if task.Status == logsstore.ExportStatusPendingCancel {
			s.updateExportTaskStatus(region, taskId, logsstore.ExportStatusCancelled, "Cancelled by user")
			return
		}
		if isTerminalExportStatus(task.Status) {
			logs.Info("Export task already terminal before upload, skipping delivery",
				logs.String("taskId", taskId),
				logs.String("status", task.Status))
			return
		}
	}
	if err := ctx.Err(); err != nil {
		s.updateExportTaskStatus(region, taskId, logsstore.ExportStatusFailed, fmt.Sprintf("service shutting down: %v", err))
		return
	}

	// destinationPrefix: "The prefix used as the start of the key for
	// every object exported. If you don't specify a value, the default
	// is exportedlogs."
	if prefix == "" {
		prefix = "exportedlogs"
	}
	s3Key := prefix + "/" + taskId + "/exportedlogs.gz"

	if err := s3Invoker.PutObject(ctx, region, bucket, s3Key, buf.Bytes(), "application/x-gzip"); err != nil {
		s.updateExportTaskStatus(region, taskId, logsstore.ExportStatusFailed, fmt.Sprintf("S3 upload error: %v", err))
		return
	}

	s.updateExportTaskStatus(region, taskId, logsstore.ExportStatusCompleted, "")
}

// isTerminalExportStatus reports whether an export task status is final:
// terminal states win races against the worker's own writes (a task the
// canceller moved to CANCELLED stays cancelled). PENDING_CANCEL is not
// terminal — the documented status vocabulary carries it as a live
// transitional state, so a later worker write still supersedes it.
func isTerminalExportStatus(status string) bool {
	return status == logsstore.ExportStatusCompleted || status == logsstore.ExportStatusFailed || status == logsstore.ExportStatusCancelled
}

// updateExportTaskStatus writes a status transition for an export task.
// A task that already reached a terminal status keeps it (the documented
// cancel flow relies on the running worker never overwriting CANCELLED);
// the keep-or-write decision runs inside the store's atomic task mutate,
// so the canceller's transition and the worker's completion cannot
// interleave. Read failures are logged with the task id instead of
// returning silently, so a wedged task is visible in the server log.
func (s *LogsService) updateExportTaskStatus(region, taskId, status, message string) {
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		logs.Error("Failed to resolve store for export task status update",
			logs.String("taskId", taskId), logs.String("status", status), logs.Err(err))
		return
	}
	err = store.MutateExportTask(taskId, func(task *logsstore.ExportTask) error {
		if isTerminalExportStatus(task.Status) && task.Status != status {
			logs.Warn("Export task already terminal, keeping existing status",
				logs.String("taskId", taskId),
				logs.String("existing", task.Status),
				logs.String("attempted", status))
			// An untouched record persists byte-identical: the terminal
			// status keeps.
			return nil
		}
		task.Status = status
		task.StatusMessage = message
		if isTerminalExportStatus(status) {
			if task.ExecutionInfo == nil {
				task.ExecutionInfo = make(map[string]interface{})
			}
			task.ExecutionInfo["completionTime"] = time.Now().UTC().UnixMilli()
		}
		return nil
	})
	if err != nil {
		logs.Error("Failed to persist export task status update",
			logs.String("taskId", taskId),
			logs.String("status", status),
			logs.Err(err))
	}
}

// DescribeExportTasks lists export tasks.
func (s *LogsService) DescribeExportTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return describeStoreFamilyHandler(s, reqCtx, req, "exportTasks",
		func(store *logsstore.Store, nextToken string, limit int32) ([]*logsstore.ExportTask, string, error) {
			return s.describeExportTasksCore(store, &DescribeExportTasksInput{
				TaskId:     request.GetParamLowerFirst(req.Parameters, "TaskId"),
				StatusCode: request.GetParamLowerFirst(req.Parameters, "StatusCode"),
				NextToken:  nextToken,
				Limit:      limit,
			})
		}, formatExportTask)
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

	return response.EmptyResponse(), nil
}

func formatExportTask(t *logsstore.ExportTask) map[string]interface{} {
	// The ExportTask response shape carries taskId, taskName, logGroupName,
	// from, to, destination, destinationPrefix, status and executionInfo —
	// logStreamNamePrefix is a request member alone, so the stored prefix
	// never renders.
	result := map[string]interface{}{
		"taskId":       t.TaskId,
		"taskName":     t.TaskName,
		"logGroupName": t.LogGroupName,
		"from":         t.From,
		"to":           t.To,
		"destination":  t.Destination,
		"status":       map[string]interface{}{"code": t.Status, "message": t.StatusMessage},
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
