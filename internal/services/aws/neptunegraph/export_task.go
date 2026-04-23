package neptunegraph

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	storecommon "vorpalstacks/internal/store/aws/common"
	ngstore "vorpalstacks/internal/store/aws/neptunegraph"
	"vorpalstacks/pkg/graphengine"
)

// StartExportTask initiates a bulk export of graph data to the specified S3 destination.
func (s *NeptuneGraphService) StartExportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	_, err = store.GetGraph(graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}

	destination := request.GetStringParam(req.Parameters, "destination")
	if destination == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "destination")
	}

	format := request.GetStringParam(req.Parameters, "format")
	if format == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "format")
	}

	kmsKeyID := request.GetStringParam(req.Parameters, "kmsKeyIdentifier")
	if kmsKeyID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "kmsKeyIdentifier")
	}

	roleArn := request.GetStringParam(req.Parameters, "roleArn")
	if roleArn == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "roleArn")
	}

	taskID := generateID("t-")
	now := time.Now().UTC()

	task := &ngstore.ExportTask{
		TaskId:           taskID,
		GraphId:          graphID,
		Status:           "INITIALIZING",
		Format:           format,
		ParquetType:      request.GetStringParam(req.Parameters, "parquetType"),
		Destination:      destination,
		RoleArn:          roleArn,
		KmsKeyIdentifier: kmsKeyID,
		StartTime:        &now,
	}

	if request.HasParam(req.Parameters, "exportFilter") {
		task.ExportFilter = parseExportFilter(req.Parameters)
	}

	if err := store.CreateExportTask(task); err != nil {
		return nil, err
	}

	s.taskWg.Add(1)
	go s.advanceExportTask(store, taskID)

	return exportTaskToResponse(task), nil
}

// GetExportTask retrieves an export task by its identifier.
func (s *NeptuneGraphService) GetExportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	taskID := request.GetStringParam(req.Parameters, "taskIdentifier")
	if taskID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "taskIdentifier")
	}

	task, err := store.GetExportTask(taskID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("export task", taskID)
		}
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

	opts := storecommon.ListOptions{
		MaxItems: clampMaxResults(request.GetIntParam(req.Parameters, "maxResults")),
		Marker:   request.GetStringParam(req.Parameters, "nextToken"),
	}

	graphID := request.GetStringParam(req.Parameters, "graphIdentifier")
	tasks, nextToken, truncated, err := store.ListExportTasks(opts, graphID)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(tasks))
	for _, t := range tasks {
		items = append(items, exportTaskSummaryToResponse(t))
	}

	result := map[string]interface{}{
		"tasks": items,
	}
	if truncated {
		result["nextToken"] = nextToken
	}
	return result, nil
}

// CancelExportTask cancels an in-progress export task, transitioning it to CANCELLED state.
func (s *NeptuneGraphService) CancelExportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	taskID := request.GetStringParam(req.Parameters, "taskIdentifier")
	if taskID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "taskIdentifier")
	}

	task, err := store.GetExportTask(taskID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("export task", taskID)
		}
		return nil, err
	}

	if task.Status == "SUCCEEDED" || task.Status == "FAILED" || task.Status == "CANCELLED" {
		return exportTaskSummaryToResponse(task), nil
	}

	task.Status = "CANCELLING"
	if err := store.UpdateExportTask(task); err != nil {
		logs.Warn("failed to update export task status to CANCELLING", logs.String("taskId", taskID), logs.Err(err))
	}

	task.Status = "CANCELLED"
	task.StatusReason = "Cancelled by user"
	if err := store.UpdateExportTask(task); err != nil {
		logs.Warn("failed to update export task status to CANCELLED", logs.String("taskId", taskID), logs.Err(err))
	}

	return exportTaskSummaryToResponse(task), nil
}

func (s *NeptuneGraphService) advanceExportTask(store *ngstore.NeptuneGraphStore, taskID string) {
	defer s.taskWg.Done()
	defer func() { resilience.RecoverPanic("NeptuneGraph advanceExportTask") }()

	task, err := store.GetExportTask(taskID)
	if err != nil {
		logs.Warn("failed to get export task", logs.String("taskId", taskID), logs.Err(err))
		return
	}

	err = store.TryAdvanceExportTask(taskID, "INITIALIZING", func(t *ngstore.ExportTask) {
		t.Status = "EXPORTING"
	})
	if err != nil {
		logs.Warn("failed to advance export task to EXPORTING", logs.String("taskId", taskID), logs.Err(err))
		return
	}

	s.enginesMu.RLock()
	entry, ok := s.activeEngines[task.GraphId]
	s.enginesMu.RUnlock()
	if !ok {
		failExportTask(store, taskID, task, "graph engine not found for graphId: "+task.GraphId)
		return
	}

	dest := task.Destination
	format := strings.ToUpper(task.Format)

	if strings.HasPrefix(dest, "s3://") {
		failExportTask(store, taskID, task, "S3 destination is not accessible in standalone mode")
		return
	}

	if !strings.HasPrefix(dest, "file://") {
		failExportTask(store, taskID, task, "unsupported destination scheme: "+dest)
		return
	}

	filePath := strings.TrimPrefix(dest, "file://")
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		failExportTask(store, taskID, task, fmt.Sprintf("failed to create export directory: %v", err))
		return
	}

	var nodeCount, edgeCount int64
	if format == "CSV" || format == "CSV+BINARY" {
		nodeCount, edgeCount, err = exportGraphCSV(entry.db, filePath)
	} else {
		failExportTask(store, taskID, task, "unsupported export format: "+task.Format)
		return
	}

	if err != nil {
		failExportTask(store, taskID, task, fmt.Sprintf("export failed: %v", err))
		return
	}

	task, err = store.GetExportTask(taskID)
	if err != nil || task.Status == "CANCELLED" {
		return
	}

	err = store.TryAdvanceExportTask(taskID, "EXPORTING", func(t *ngstore.ExportTask) {
		t.Status = "SUCCEEDED"
		now := time.Now().UTC()
		sinceStart := int64(now.Sub(*t.StartTime).Seconds())
		details := t.ExportTaskDetails
		if details == nil {
			details = &ngstore.ExportTaskDetails{}
		}
		details.ProgressPercentage = int32Ptr(100)
		details.StartTime = t.StartTime
		details.TimeElapsedSeconds = int64Ptr(sinceStart)
		details.NumVerticesWritten = int64Ptr(nodeCount)
		details.NumEdgesWritten = int64Ptr(edgeCount)
		t.ExportTaskDetails = details
	})
	if err != nil {
		logs.Warn("failed to advance export task to SUCCEEDED", logs.String("taskId", taskID), logs.Err(err))
	}
}

func failExportTask(store *ngstore.NeptuneGraphStore, taskID string, task *ngstore.ExportTask, reason string) {
	logs.Warn("export task failed", logs.String("taskId", taskID), logs.String("reason", reason))
	now := time.Now().UTC()
	sinceStart := int64(0)
	if task.StartTime != nil {
		sinceStart = int64(now.Sub(*task.StartTime).Seconds())
	}
	store.TryAdvanceExportTask(taskID, "EXPORTING", func(t *ngstore.ExportTask) {
		t.Status = "FAILED"
		t.StatusReason = reason
		t.ExportTaskDetails = &ngstore.ExportTaskDetails{
			ProgressPercentage: int32Ptr(0),
			StartTime:          t.StartTime,
			TimeElapsedSeconds: int64Ptr(sinceStart),
		}
	})
}

func exportGraphCSV(db *graphengine.DB, filePath string) (int64, int64, error) {
	nodesFile := filePath
	if !strings.HasSuffix(strings.ToLower(filePath), ".csv") {
		nodesFile = filePath + "_nodes.csv"
	}
	edgesFile := strings.TrimSuffix(nodesFile, "_nodes.csv") + "_edges.csv"
	if nodesFile == filePath {
		edgesFile = filePath + "_edges.csv"
	}

	var nodeCount, edgeCount int64

	nf, err := os.Create(nodesFile)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create nodes file: %w", err)
	}
	defer nf.Close()

	nodeW := csv.NewWriter(nf)
	defer nodeW.Flush()

	ef, err := os.Create(edgesFile)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create edges file: %w", err)
	}
	defer ef.Close()

	edgeW := csv.NewWriter(ef)
	defer edgeW.Flush()

	var propKeys []string

	err = db.ForEachNode(func(node *graphengine.Node) error {
		propKeys = sortedPropKeys(node.Props)
		record := make([]string, 0, 2+len(propKeys))
		record = append(record, fmt.Sprintf("%d", node.ID))
		record = append(record, strings.Join(node.Labels, ";"))
		for _, k := range propKeys {
			record = append(record, formatPropValue(node.Props[k]))
		}
		if err := nodeW.Write(record); err != nil {
			return err
		}
		nodeCount++
		return nil
	})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to enumerate nodes: %w", err)
	}

	err = db.ForEachEdge(func(edge *graphengine.Edge) error {
		propKeys = sortedPropKeys(edge.Props)
		record := make([]string, 0, 4+len(propKeys))
		record = append(record, fmt.Sprintf("%d", edge.ID))
		record = append(record, edge.Label)
		record = append(record, fmt.Sprintf("%d", edge.From))
		record = append(record, fmt.Sprintf("%d", edge.To))
		for _, k := range propKeys {
			record = append(record, formatPropValue(edge.Props[k]))
		}
		if err := edgeW.Write(record); err != nil {
			return err
		}
		edgeCount++
		return nil
	})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to enumerate edges: %w", err)
	}

	return nodeCount, edgeCount, nil
}

func sortedPropKeys(props graphengine.Props) []string {
	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
}

func formatPropValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		b, err := json.Marshal(val)
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
		return string(b)
	}
}

func parseExportFilter(params map[string]interface{}) *ngstore.ExportFilter {
	v, ok := params["exportFilter"]
	if !ok || v == nil {
		return nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var f ngstore.ExportFilter
	if err := json.Unmarshal(data, &f); err != nil {
		return nil
	}
	return &f
}
