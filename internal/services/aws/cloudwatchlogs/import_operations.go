package cloudwatchlogs

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	svcarn "vorpalstacks/internal/utils/aws/arn"
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

// importStreamName derives the log stream an imported S3 object's events
// land in: the object key with the two characters a stream name forbids
// (':' and '*') replaced, capped to the stream-name length limit with a
// digest suffix so long keys stay unique. AWS documents no stream layout
// for the managed import group; one stream per source object keeps every
// object's events (the fix for the constant-stream-name data loss).
func importStreamName(key string) string {
	name := strings.Map(func(r rune) rune {
		if r == ':' || r == '*' {
			return '_'
		}
		return r
	}, key)
	if len(name) > logsstore.MaxLogStreamNameLength {
		digest := hex.EncodeToString(sha256sum(key))
		keep := logsstore.MaxLogStreamNameLength - 45
		// Trim whole runes only: slicing mid-rune would mint an invalid
		// UTF-8 stream name.
		for len(name) > keep {
			_, size := utf8.DecodeLastRuneInString(name)
			name = name[:len(name)-size]
		}
		name = name + "-" + digest[:44]
	}
	return name
}

func sha256sum(s string) []byte {
	sum := sha256.Sum256([]byte(s))
	return sum[:]
}

// isTerminalImportStatus reports whether an import task status is final;
// terminal states win races against the worker's own writes.
func isTerminalImportStatus(status string) bool {
	return status == logsstore.ImportStatusCompleted || status == logsstore.ImportStatusFailed || status == logsstore.ImportStatusCancelled
}

func (s *LogsService) executeImportTask(ctx context.Context, region, importId, bucket, s3Key, logGroupName string) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("PANIC in import task",
				logs.String("importId", importId),
				logs.Any("panic", r))
			s.updateImportTaskStatus(region, importId, logsstore.ImportStatusFailed, fmt.Sprintf("panic: %v", r))
		}
	}()

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		s.updateImportTaskStatus(region, importId, logsstore.ImportStatusFailed, fmt.Sprintf("store error: %v", err))
		return
	}

	if s.eventBus() == nil {
		s.updateImportTaskStatus(region, importId, logsstore.ImportStatusFailed, "event bus not configured")
		return
	}

	s3Invoker := s.eventBus().S3Invoker()
	if s3Invoker == nil {
		s.updateImportTaskStatus(region, importId, logsstore.ImportStatusFailed, "S3 invoker not configured")
		return
	}

	task, err := store.GetImportTask(importId)
	if err != nil {
		s.updateImportTaskStatus(region, importId, logsstore.ImportStatusFailed, fmt.Sprintf("task read error: %v", err))
		return
	}
	filterStart, filterEnd := importFilterWindow(task.ImportFilter)

	keys := []string{s3Key}
	if s3Key == "" {
		listed, err := fetchAllS3ObjectKeys(ctx, s3Invoker, region, bucket, "")
		if err != nil {
			s.updateImportTaskStatus(region, importId, logsstore.ImportStatusFailed, fmt.Sprintf("S3 list error: %v", err))
			return
		}
		keys = listed
	}

	// Batches group the imported events by event-time bucket — the
	// granularity the DescribeImportTaskBatches documentation names.
	// hourFailures carries the first failure each bucket observed, so a
	// partially-failed import stamps the affected batches FAILED with
	// their errorMessage instead of reporting everything COMPLETED.
	batchIndex := make(map[int64]*logsstore.ImportBatch)
	hourFailures := make(map[int64]string)
	totalBytes := int64(0)
	succeeded, failed := 0, 0

	bucketEvents := func(events []logsstore.LogEntry) {
		for _, e := range events {
			hourStart := e.Timestamp / (60 * 60 * 1000) * (60 * 60 * 1000)
			if _, ok := batchIndex[hourStart]; !ok {
				batchIndex[hourStart] = &logsstore.ImportBatch{
					BatchId:     fmt.Sprintf("batch-%d", hourStart),
					Status:      logsstore.ImportStatusInProgress,
					HourStartMs: hourStart,
				}
			}
		}
	}

	for _, key := range keys {
		if key == "" {
			continue
		}

		// Cancellation checkpoint between objects: a task the canceller
		// moved to CANCELLED stops importing and keeps that state.
		if current, err := store.GetImportTask(importId); err == nil && isTerminalImportStatus(current.ImportStatus) {
			s.finalizeImportTask(region, importId, current.ImportStatus, "", totalBytes, batchIndex, hourFailures)
			return
		}
		if err := ctx.Err(); err != nil {
			s.finalizeImportTask(region, importId, logsstore.ImportStatusFailed, fmt.Sprintf("service shutting down: %v", err), totalBytes, batchIndex, hourFailures)
			return
		}

		data, err := s3Invoker.GetObject(ctx, region, bucket, key, logsstore.MaxImportObjectBytes)
		if err != nil {
			failed++
			logs.Warn("Failed to get S3 object during import",
				logs.String("bucket", bucket),
				logs.String("key", key),
				logs.Err(err))
			continue
		}

		content, err := decompressIfNeeded(data)
		if err != nil {
			failed++
			logs.Warn("Failed to decompress S3 object during import",
				logs.String("key", key),
				logs.Err(err))
			continue
		}

		events := applyImportFilter(parseImportContent(content), filterStart, filterEnd)
		if len(events) == 0 {
			continue
		}

		// The object's events join their batches whether or not the
		// write succeeds — a failed write leaves its buckets FAILED
		// rather than absent.
		bucketEvents(events)
		streamName := importStreamName(key)
		ls := logsstore.NewLogStream(streamName, logGroupName)
		if err := store.CreateLogStream(ls); err != nil && !errors.Is(err, logsstore.ErrLogStreamAlreadyExists) {
			failed++
			logs.Warn("Failed to create import log stream",
				logs.String("importId", importId),
				logs.String("stream", streamName),
				logs.Err(err))
			for _, e := range events {
				hourStart := e.Timestamp / (60 * 60 * 1000) * (60 * 60 * 1000)
				if hourFailures[hourStart] == "" {
					hourFailures[hourStart] = fmt.Sprintf("stream %s create error: %v", streamName, err)
				}
			}
			s.persistImportProgress(region, importId, totalBytes, batchIndex)
			continue
		}

		if _, err := store.PutLogEvents(logGroupName, streamName, events); err != nil {
			failed++
			logs.Warn("Failed to put imported log events",
				logs.String("importId", importId),
				logs.String("stream", streamName),
				logs.Err(err))
			for _, e := range events {
				hourStart := e.Timestamp / (60 * 60 * 1000) * (60 * 60 * 1000)
				if hourFailures[hourStart] == "" {
					hourFailures[hourStart] = fmt.Sprintf("stream %s write error: %v", streamName, err)
				}
			}
			s.persistImportProgress(region, importId, totalBytes, batchIndex)
			continue
		}
		// Imported events are ingested events: the field-index bounds
		// scan covers them too (no transformed copies exist on this
		// path — the import applies no transformer).
		s.scanIngestedFieldIndexes(store, logGroupName, streamName, events, nil)
		succeeded++
		totalBytes += int64(len(data))
		// Progress is visible while the task runs: the batches report
		// IN_PROGRESS and a mid-flight cancellation observes the
		// statistics gathered so far.
		s.persistImportProgress(region, importId, totalBytes, batchIndex)
	}

	status := logsstore.ImportStatusCompleted
	message := ""
	if succeeded == 0 && failed > 0 {
		status = logsstore.ImportStatusFailed
		message = fmt.Sprintf("all %d source objects failed to import", failed)
	} else if failed > 0 {
		message = fmt.Sprintf("%d of %d source objects failed to import", failed, succeeded+failed)
	}
	s.finalizeImportTask(region, importId, status, message, totalBytes, batchIndex, hourFailures)
}

// executeImportTaskFromEDS runs the documented source form: a CloudTrail
// Lake event data store walked through the cross-service invoker in
// ascending event-time order, every record's own JSON payload ingested as
// one log event into the managed group. The batching, progress
// persistence and finalisation mirror the S3 executor's; one stream holds
// the whole walk because the event data store is one source (the S3
// executor's one-stream-per-object split has no object granularity here).
func (s *LogsService) executeImportTaskFromEDS(ctx context.Context, region, importId, edsArn, logGroupName string) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("PANIC in import task",
				logs.String("importId", importId),
				logs.Any("panic", r))
			s.updateImportTaskStatus(region, importId, logsstore.ImportStatusFailed, fmt.Sprintf("panic: %v", r))
		}
	}()

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		s.updateImportTaskStatus(region, importId, logsstore.ImportStatusFailed, fmt.Sprintf("store error: %v", err))
		return
	}
	if s.eventBus() == nil {
		s.updateImportTaskStatus(region, importId, logsstore.ImportStatusFailed, "event bus not configured")
		return
	}
	ctInvoker := s.eventBus().CloudTrailInvoker()
	if ctInvoker == nil {
		s.updateImportTaskStatus(region, importId, logsstore.ImportStatusFailed, "CloudTrail invoker not configured")
		return
	}
	task, err := store.GetImportTask(importId)
	if err != nil {
		s.updateImportTaskStatus(region, importId, logsstore.ImportStatusFailed, fmt.Sprintf("task read error: %v", err))
		return
	}
	filterStart, filterEnd := importFilterWindow(task.ImportFilter)
	var startTime, endTime *time.Time
	if filterStart > 0 {
		t := time.UnixMilli(filterStart)
		startTime = &t
	}
	if filterEnd > 0 {
		t := time.UnixMilli(filterEnd)
		endTime = &t
	}

	streamName := importStreamName("eventdatastore/" + edsSourceID(edsArn))
	batchIndex := make(map[int64]*logsstore.ImportBatch)
	hourFailures := make(map[int64]string)
	totalBytes := int64(0)
	imported := 0

	if err := store.CreateLogStream(logsstore.NewLogStream(streamName, logGroupName)); err != nil && !errors.Is(err, logsstore.ErrLogStreamAlreadyExists) {
		s.finalizeImportTask(region, importId, logsstore.ImportStatusFailed,
			fmt.Sprintf("stream %s create error: %v", streamName, err), totalBytes, batchIndex, hourFailures)
		return
	}

	nextToken := ""
	for {
		// Cancellation checkpoint between pages: a task the canceller
		// moved to CANCELLED stops importing and keeps that state.
		if current, err := store.GetImportTask(importId); err == nil && isTerminalImportStatus(current.ImportStatus) {
			s.finalizeImportTask(region, importId, current.ImportStatus, "", totalBytes, batchIndex, hourFailures)
			return
		}
		if err := ctx.Err(); err != nil {
			s.finalizeImportTask(region, importId, logsstore.ImportStatusFailed, fmt.Sprintf("service shutting down: %v", err), totalBytes, batchIndex, hourFailures)
			return
		}
		page, token, err := ctInvoker.LookupEDSEvents(ctx, edsArn, startTime, endTime, importEDSPageSize, nextToken)
		if err != nil {
			hourFailures[0] = fmt.Sprintf("event data store walk error: %v", err)
			status := logsstore.ImportStatusCompleted
			if imported == 0 {
				status = logsstore.ImportStatusFailed
			}
			s.finalizeImportTask(region, importId, status, hourFailures[0], totalBytes, batchIndex, hourFailures)
			return
		}
		var entries []logsstore.LogEntry
		for _, e := range page {
			entries = append(entries, logsstore.LogEntry{
				Timestamp: e.EventTime.UnixMilli(),
				Message:   e.Payload,
			})
			totalBytes += int64(len(e.Payload))
		}
		if len(entries) > 0 {
			for _, e := range entries {
				hourStart := e.Timestamp / (60 * 60 * 1000) * (60 * 60 * 1000)
				if _, ok := batchIndex[hourStart]; !ok {
					batchIndex[hourStart] = &logsstore.ImportBatch{
						BatchId:     fmt.Sprintf("batch-%d", hourStart),
						Status:      logsstore.ImportStatusInProgress,
						HourStartMs: hourStart,
					}
				}
			}
			// A page's entries split into batch-contract-sized writes: a
			// page of sizeable records can exceed the PutLogEvents batch
			// bound as a whole.
			for _, batch := range splitDeliveryBatches(entries) {
				if _, err := store.PutLogEvents(logGroupName, streamName, batch); err != nil {
					logs.Warn("Failed to put imported event data store events",
						logs.String("importId", importId),
						logs.String("stream", streamName),
						logs.Err(err))
					for _, e := range batch {
						hourStart := e.Timestamp / (60 * 60 * 1000) * (60 * 60 * 1000)
						if hourFailures[hourStart] == "" {
							hourFailures[hourStart] = fmt.Sprintf("stream %s write error: %v", streamName, err)
						}
					}
					continue
				}
				// Ingested events feed the field-index bounds scan (no
				// transformed copies exist on the import path).
				s.scanIngestedFieldIndexes(store, logGroupName, streamName, batch, nil)
				imported += len(batch)
			}
			s.persistImportProgress(region, importId, totalBytes, batchIndex)
		}
		if token == "" {
			break
		}
		nextToken = token
	}

	status := logsstore.ImportStatusCompleted
	message := ""
	if imported == 0 && len(hourFailures) > 0 {
		status = logsstore.ImportStatusFailed
		message = "the event data store events failed to import"
	} else if len(hourFailures) > 0 {
		message = "some event data store events failed to import"
	}
	s.finalizeImportTask(region, importId, status, message, totalBytes, batchIndex, hourFailures)
}

// importEDSPageSize is the event data store walk's page size: small
// enough that a page's entries fit one PutLogEvents batch whatever the
// payload sizes, bounded away from the walk's contract-free interior
// (the CloudTrail store serves any positive bound).
const importEDSPageSize = 1000

// edsSourceID extracts the event data store id from its ARN for the
// managed group's stream name.
func edsSourceID(edsArn string) string {
	_, _, _, _, resource := svcarn.SplitARN(edsArn)
	id, _ := strings.CutPrefix(resource, "eventdatastore/")
	return id
}

// applyImportFilter keeps the events inside the task's declared window
// (importFilter startEventTime/endEventTime, epoch milliseconds). Events
// the filter excludes are not imported.
func applyImportFilter(events []logsstore.LogEntry, startMs, endMs int64) []logsstore.LogEntry {
	if startMs == 0 && endMs == 0 {
		return events
	}
	filtered := make([]logsstore.LogEntry, 0, len(events))
	for _, e := range events {
		if startMs > 0 && e.Timestamp < startMs {
			continue
		}
		if endMs > 0 && e.Timestamp > endMs {
			continue
		}
		filtered = append(filtered, e)
	}
	return filtered
}

// importFilterWindow extracts the validated window from the task's
// importFilter member.
func importFilterWindow(filter map[string]interface{}) (startMs, endMs int64) {
	if filter == nil {
		return 0, 0
	}
	if v, ok := filter["startEventTime"]; ok {
		if n, ok := toInt64(v); ok {
			startMs = n
		}
	}
	if v, ok := filter["endEventTime"]; ok {
		if n, ok := toInt64(v); ok {
			endMs = n
		}
	}
	return startMs, endMs
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
	s.finalizeImportTask(region, importId, status, message, -1, nil, nil)
}

// importBatchSnapshot renders the in-memory batch index as the ordered
// batch list the task record persists.
func importBatchSnapshot(batchIndex map[int64]*logsstore.ImportBatch) []*logsstore.ImportBatch {
	hours := make([]int64, 0, len(batchIndex))
	for h := range batchIndex {
		hours = append(hours, h)
	}
	sort.Slice(hours, func(i, j int) bool { return hours[i] < hours[j] })
	batches := make([]*logsstore.ImportBatch, 0, len(hours))
	for _, h := range hours {
		batches = append(batches, batchIndex[h])
	}
	return batches
}

// persistImportProgress records the task's in-flight batch set and
// statistics so the batch listing observes live batches and a mid-flight
// cancellation reports "Statistics about the import progress at the time
// of cancellation": the batches carry IN_PROGRESS while the task runs.
// The write never touches ImportStatus, so it cannot race a
// cancellation's terminal transition.
func (s *LogsService) persistImportProgress(region, importId string, bytesImported int64, batchIndex map[int64]*logsstore.ImportBatch) {
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return
	}
	err = store.MutateImportTask(importId, func(task *logsstore.ImportTask) error {
		task.ImportBatches = importBatchSnapshot(batchIndex)
		if bytesImported > 0 {
			task.ImportStatistics = map[string]interface{}{
				"bytesImported": bytesImported,
			}
		}
		task.LastUpdatedTime = time.Now().UTC().UnixMilli()
		return nil
	})
	if err != nil {
		logs.Warn("Failed to persist import progress",
			logs.String("importId", importId), logs.Err(err))
	}
}

// finalizeImportTask writes an import task's terminal transition in one
// record: status, error message, statistics and the batch outcomes. A
// task that already reached a terminal status keeps it (a canceller won
// the race); read failures are logged with the task id instead of
// returning silently. bytesImported < 0 leaves the existing statistics
// untouched (the early-failure paths have nothing to report yet); a
// batch whose bucket observed a write failure stamps FAILED with its
// errorMessage ("Only present when status is FAILED"), the rest follow
// the task's terminal status.
func (s *LogsService) finalizeImportTask(region, importId, status, message string, bytesImported int64, batchIndex map[int64]*logsstore.ImportBatch, hourFailures map[int64]string) {
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		logs.Error("Failed to resolve store for import task status update",
			logs.String("importId", importId), logs.String("status", status), logs.Err(err))
		return
	}
	// The keep-or-write decision runs inside the store's atomic task
	// mutate: the canceller's CANCELLED and the worker's terminal status
	// cannot interleave, and a terminal record keeps its status (an
	// untouched record persists byte-identical).
	err = store.MutateImportTask(importId, func(task *logsstore.ImportTask) error {
		if isTerminalImportStatus(task.ImportStatus) && task.ImportStatus != status {
			logs.Warn("Import task already terminal, keeping existing status",
				logs.String("importId", importId),
				logs.String("existing", task.ImportStatus),
				logs.String("attempted", status))
			return nil
		}
		task.ImportStatus = status
		if message != "" {
			task.ErrorMessage = message
		}
		if bytesImported >= 0 {
			task.ImportStatistics = map[string]interface{}{
				"bytesImported": bytesImported,
			}
		}
		if batchIndex != nil {
			batches := importBatchSnapshot(batchIndex)
			for _, b := range batches {
				if msg := hourFailures[b.HourStartMs]; msg != "" {
					b.Status = logsstore.ImportStatusFailed
					b.ErrorMessage = msg
					continue
				}
				b.Status = status
			}
			task.ImportBatches = batches
		}
		task.LastUpdatedTime = time.Now().UTC().UnixMilli()
		return nil
	})
	if err != nil {
		logs.Error("Failed to persist import task status update",
			logs.String("importId", importId),
			logs.String("status", status),
			logs.Err(err))
	}
}

// DescribeImportTasks lists import tasks.
func (s *LogsService) DescribeImportTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	items, nextMarker, err := s.describeImportTasksCore(store,
		request.GetParamLowerFirst(req.Parameters, "ImportId"),
		request.GetParamLowerFirst(req.Parameters, "ImportStatus"),
		request.GetParamLowerFirst(req.Parameters, "ImportSourceArn"),
		request.GetParamLowerFirst(req.Parameters, "NextToken"),
		int32(request.GetIntParam(req.Parameters, "Limit")))
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
		"importStatus": logsstore.ImportStatusCancelled,
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
	store, err := s.getLogsStoreByRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	result, err := s.describeImportTaskBatchesCore(store, &DescribeImportTaskBatchesInput{
		ImportId:          request.GetParamLowerFirst(req.Parameters, "ImportId"),
		BatchImportStatus: request.GetStringList(req.Parameters, "BatchImportStatus"),
		NextToken:         request.GetParamLowerFirst(req.Parameters, "NextToken"),
		Limit:             int32(request.GetIntParam(req.Parameters, "Limit")),
	})
	if err != nil {
		return nil, err
	}

	batches := make([]map[string]interface{}, 0, len(result.ImportBatches))
	for _, b := range result.ImportBatches {
		entry := map[string]interface{}{
			"batchId": b.BatchId,
			"status":  b.Status,
		}
		if b.ErrorMessage != "" {
			entry["errorMessage"] = b.ErrorMessage
		}
		batches = append(batches, entry)
	}

	resp := map[string]interface{}{
		"importId":        result.Task.ImportId,
		"importSourceArn": result.Task.ImportSourceArn,
		"importBatches":   batches,
	}
	if result.NextToken != "" {
		resp["nextToken"] = result.NextToken
	}
	return resp, nil
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

// deriveManagedLogGroupName names the managed destination group of an
// import after the source bucket: the bucket is the source resource the
// ARN addresses, and a bucket-only ARN carries no slash at all — the
// post-slash remainder would leak the whole ARN (colons included) into a
// name the log-group charset forbids. Bucket names ([a-z0-9.-]) sit
// wholly inside that charset.
func deriveManagedLogGroupName(bucket string) string {
	return "/aws/imported/cloudtrail-lake/" + bucket
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
