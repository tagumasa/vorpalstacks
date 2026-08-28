package neptunegraph

// Export task Core functions: the single validation and persistence path for
// the export task operations, plus the async export worker.

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage/graphengine"
	storecommon "vorpalstacks/internal/store/aws/common"
	ngstore "vorpalstacks/internal/store/aws/rds/neptunegraph"
)

// StartExportTaskInput carries the wire-parsed StartExportTask request.
type StartExportTaskInput struct {
	GraphIdentifier  string
	Format           string
	KmsKeyIdentifier string
	RoleArn          string
	Destination      string
	ParquetType      string
	HasExportFilter  bool
	ExportFilter     *ngstore.ExportFilter
}

// GetExportTaskInput carries the wire-parsed GetExportTask request.
type GetExportTaskInput struct {
	TaskIdentifier string
}

// ListExportTasksInput carries the wire-parsed ListExportTasks request with
// the already-clamped page size.
type ListExportTasksInput struct {
	MaxItems        int
	Marker          string
	GraphIdentifier string
}

// ListExportTasksResult carries the export task records for one list page.
type ListExportTasksResult struct {
	Tasks     []*ngstore.ExportTask
	NextToken string
	Truncated bool
}

// CancelExportTaskInput carries the wire-parsed CancelExportTask request.
type CancelExportTaskInput struct {
	TaskIdentifier string
}

// startExportTaskCore validates the request and starts a bulk export of
// graph data to the specified destination.
func (s *NeptuneGraphService) startExportTaskCore(store *ngstore.NeptuneGraphStore, in *StartExportTaskInput) (*ngstore.ExportTask, error) {
	graphID := in.GraphIdentifier
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier")
	}

	graph, err := s.resolveGraphIdentifier(store, graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}

	format := in.Format
	if err := validateExportFormat(format); err != nil {
		return nil, err
	}

	kmsKeyID := in.KmsKeyIdentifier
	if err := validateKmsKeyArn(kmsKeyID, true); err != nil {
		return nil, err
	}

	roleArn := in.RoleArn
	if err := validateRoleArn(roleArn); err != nil {
		return nil, err
	}

	destination := in.Destination
	if err := validateDestination(destination); err != nil {
		return nil, err
	}

	parquetType := in.ParquetType
	if err := validateParquetType(parquetType); err != nil {
		return nil, err
	}

	taskID := generateID("t-")
	now := time.Now().UTC()

	task := &ngstore.ExportTask{
		TaskId:           taskID,
		GraphId:          graph.Id,
		Status:           "INITIALIZING",
		Format:           format,
		ParquetType:      parquetType,
		Destination:      destination,
		RoleArn:          roleArn,
		KmsKeyIdentifier: kmsKeyID,
		StartTime:        &now,
	}

	if in.HasExportFilter {
		task.ExportFilter = in.ExportFilter
	}

	if err := store.CreateExportTask(task); err != nil {
		return nil, err
	}

	s.taskWg.Add(1)
	go s.advanceExportTask(store, taskID)

	return task, nil
}

// getExportTaskCore retrieves an export task record, mapping a missing task
// to the documented ResourceNotFoundException.
func (s *NeptuneGraphService) getExportTaskCore(store *ngstore.NeptuneGraphStore, taskID string) (*ngstore.ExportTask, error) {
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

	return task, nil
}

// listExportTasksCore reads one page of export task records, optionally
// filtered by graph, returning the store error unchanged so both planes keep
// their own error mapping.
func (s *NeptuneGraphService) listExportTasksCore(store *ngstore.NeptuneGraphStore, in *ListExportTasksInput) (*ListExportTasksResult, error) {
	tasks, nextToken, truncated, err := store.ListExportTasks(storecommon.ListOptions{
		MaxItems: in.MaxItems,
		Marker:   in.Marker,
	}, in.GraphIdentifier)
	if err != nil {
		return nil, err
	}
	return &ListExportTasksResult{Tasks: tasks, NextToken: nextToken, Truncated: truncated}, nil
}

// cancelExportTaskCore cancels an in-progress export task, transitioning it
// to CANCELLING state.
func (s *NeptuneGraphService) cancelExportTaskCore(store *ngstore.NeptuneGraphStore, in *CancelExportTaskInput) (*ngstore.ExportTask, error) {
	taskID := in.TaskIdentifier
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

	if task.Status == "SUCCEEDED" || task.Status == "FAILED" || task.Status == "CANCELLED" || task.Status == "CANCELLING" {
		return task, nil
	}

	originalStatus := task.Status
	task.Status = "CANCELLING"
	task.StatusReason = "Cancelled by user"
	if err := store.TryAdvanceExportTask(taskID, originalStatus, func(t *ngstore.ExportTask) {
		t.Status = "CANCELLING"
		t.StatusReason = "Cancelled by user"
	}); err != nil {
		logs.Warn("failed to cancel export task", logs.String("taskId", taskID), logs.Err(err))
	}

	return task, nil
}

func (s *NeptuneGraphService) advanceExportTask(store *ngstore.NeptuneGraphStore, taskID string) {
	defer s.taskWg.Done()
	defer func() { resilience.RecoverPanic("NeptuneGraph advanceExportTask") }()

	task, err := store.GetExportTask(taskID)
	if err != nil {
		logs.Error("failed to get export task", logs.String("taskId", taskID), logs.Err(err))
		return
	}

	err = store.TryAdvanceExportTask(taskID, "INITIALIZING", func(t *ngstore.ExportTask) {
		t.Status = "EXPORTING"
	})
	if err != nil {
		logs.Warn("failed to advance export task to EXPORTING", logs.String("taskId", taskID), logs.Err(err))
		// Detect concurrent CancelExportTask: status may have transitioned
		// from INITIALIZING to CANCELLING before the goroutine started.
		current, getErr := store.GetExportTask(taskID)
		if getErr == nil {
			if current.Status == "CANCELLING" {
				_ = store.TryAdvanceExportTask(taskID, "CANCELLING", func(t *ngstore.ExportTask) {
					t.Status = "CANCELLED"
				})
			} else if current.Status != "CANCELLED" {
				// Task is still INITIALIZING (INITIALIZING->EXPORTING
				// failed). failExportTask hardcodes "EXPORTING" as the
				// expected state, which would silently fail here because
				// the actual state is INITIALIZING. Transition directly
				// from the current state to FAILED instead.
				now := time.Now().UTC()
				sinceStart := int64(0)
				if current.StartTime != nil {
					sinceStart = int64(now.Sub(*current.StartTime).Seconds())
				}
				_ = store.TryAdvanceExportTask(taskID, current.Status, func(t *ngstore.ExportTask) {
					t.Status = "FAILED"
					t.StatusReason = "failed to advance to EXPORTING"
					t.ExportTaskDetails = &ngstore.ExportTaskDetails{
						ProgressPercentage: proto.Int32(0),
						StartTime:          t.StartTime,
						TimeElapsedSeconds: proto.Int64(sinceStart),
					}
				})
			}
		}
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
	if format == "CSV" {
		nodeCount, edgeCount, err = exportGraphCSV(entry.db, filePath, task.ExportFilter)
	} else {
		failExportTask(store, taskID, task, "unsupported export format: "+task.Format)
		return
	}

	if err != nil {
		failExportTask(store, taskID, task, fmt.Sprintf("export failed: %v", err))
		return
	}

	task, err = store.GetExportTask(taskID)
	if err != nil {
		return
	}
	// If CancelExportTask set CANCELLING during the export, complete the
	// transition to CANCELLED here. The TryAdvanceExportTask at line 251
	// would fail anyway (status is no longer EXPORTING), so we handle
	// the cancellation explicitly before returning.
	if task.Status == "CANCELLING" {
		_ = store.TryAdvanceExportTask(taskID, "CANCELLING", func(t *ngstore.ExportTask) {
			t.Status = "CANCELLED"
		})
		return
	}
	if task.Status == "CANCELLED" {
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
		details.ProgressPercentage = proto.Int32(100)
		details.StartTime = t.StartTime
		details.TimeElapsedSeconds = proto.Int64(sinceStart)
		details.NumVerticesWritten = proto.Int64(nodeCount)
		details.NumEdgesWritten = proto.Int64(edgeCount)
		t.ExportTaskDetails = details
	})
	if err != nil {
		logs.Warn("failed to advance export task to SUCCEEDED", logs.String("taskId", taskID), logs.Err(err))
		// Detect concurrent CancelExportTask: status may have transitioned
		// from EXPORTING to CANCELLING while the export was running.
		current, getErr := store.GetExportTask(taskID)
		if getErr == nil && current.Status == "CANCELLING" {
			_ = store.TryAdvanceExportTask(taskID, "CANCELLING", func(t *ngstore.ExportTask) {
				t.Status = "CANCELLED"
			})
		}
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
			ProgressPercentage: proto.Int32(0),
			StartTime:          t.StartTime,
			TimeElapsedSeconds: proto.Int64(sinceStart),
		}
	})
}

// shouldExportNode returns true if the node passes the vertex filter.
// If no filter is set, all nodes pass.
func shouldExportNode(node *graphengine.Node, filter *ngstore.ExportFilter) bool {
	if filter == nil || len(filter.VertexFilter) == 0 {
		return true
	}
	for _, label := range node.Labels {
		if _, ok := filter.VertexFilter[label]; ok {
			return true
		}
	}
	return false
}

// shouldExportEdge returns true if the edge passes the edge filter.
// If no filter is set, all edges pass.
func shouldExportEdge(edge *graphengine.Edge, filter *ngstore.ExportFilter) bool {
	if filter == nil || len(filter.EdgeFilter) == 0 {
		return true
	}
	_, ok := filter.EdgeFilter[edge.Label]
	return ok
}

// filterNodeProps returns the properties to export for a node after applying
// the ExportFilter. If the filter specifies properties for the node's label,
// only those properties are included (renamed via SourcePropertyName if set).
func filterNodeProps(node *graphengine.Node, filter *ngstore.ExportFilter) graphengine.Props {
	if filter == nil || len(filter.VertexFilter) == 0 {
		return node.Props
	}
	result := make(graphengine.Props)
	for _, label := range node.Labels {
		elem, ok := filter.VertexFilter[label]
		if !ok {
			continue
		}
		for propName, attrs := range elem.Properties {
			if val, hasProp := node.Props[propName]; hasProp {
				outKey := propName
				if attrs.SourcePropertyName != nil && *attrs.SourcePropertyName != "" {
					outKey = *attrs.SourcePropertyName
				}
				result[outKey] = formatPropByType(val, attrs.OutputType)
			}
		}
	}
	// If the label matched but no properties were specified, export all.
	if len(result) == 0 {
		for _, label := range node.Labels {
			if elem, ok := filter.VertexFilter[label]; ok && len(elem.Properties) == 0 {
				return node.Props
			}
		}
	}
	return result
}

// filterEdgeProps returns the properties to export for an edge after applying
// the ExportFilter.
func filterEdgeProps(edge *graphengine.Edge, filter *ngstore.ExportFilter) graphengine.Props {
	if filter == nil || len(filter.EdgeFilter) == 0 {
		return edge.Props
	}
	elem, ok := filter.EdgeFilter[edge.Label]
	if !ok {
		return edge.Props
	}
	if len(elem.Properties) == 0 {
		return edge.Props
	}
	result := make(graphengine.Props)
	for propName, attrs := range elem.Properties {
		if val, hasProp := edge.Props[propName]; hasProp {
			outKey := propName
			if attrs.SourcePropertyName != nil && *attrs.SourcePropertyName != "" {
				outKey = *attrs.SourcePropertyName
			}
			result[outKey] = formatPropByType(val, attrs.OutputType)
		}
	}
	return result
}

// formatPropByType converts a property value according to the specified output type.
func formatPropByType(val interface{}, outputType *string) interface{} {
	if outputType == nil {
		return val
	}
	switch strings.ToUpper(*outputType) {
	case "STRING":
		return fmt.Sprintf("%v", val)
	case "NUMBER":
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case int64:
			return float64(v)
		default:
			return val
		}
	default:
		return val
	}
}

func exportGraphCSV(db *graphengine.DB, filePath string, filter *ngstore.ExportFilter) (int64, int64, error) {
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

	// Collect property keys from filtered nodes and edges.
	var propKeys []string
	allPropKeys := make(map[string]bool)
	_ = db.ForEachNode(func(node *graphengine.Node) error {
		if !shouldExportNode(node, filter) {
			return nil
		}
		props := filterNodeProps(node, filter)
		for k := range props {
			allPropKeys[k] = true
		}
		return nil
	})
	_ = db.ForEachEdge(func(edge *graphengine.Edge) error {
		if !shouldExportEdge(edge, filter) {
			return nil
		}
		props := filterEdgeProps(edge, filter)
		for k := range props {
			allPropKeys[k] = true
		}
		return nil
	})
	for k := range allPropKeys {
		propKeys = append(propKeys, k)
	}
	sort.Strings(propKeys)

	nodeHeader := make([]string, 0, 2+len(propKeys))
	nodeHeader = append(nodeHeader, "~id", "~label")
	nodeHeader = append(nodeHeader, propKeys...)
	if err := nodeW.Write(nodeHeader); err != nil {
		return 0, 0, fmt.Errorf("failed to write node header: %w", err)
	}

	edgeHeader := make([]string, 0, 4+len(propKeys))
	edgeHeader = append(edgeHeader, "~id", "~label", "~from", "~to")
	edgeHeader = append(edgeHeader, propKeys...)
	if err := edgeW.Write(edgeHeader); err != nil {
		return 0, 0, fmt.Errorf("failed to write edge header: %w", err)
	}

	err = db.ForEachNode(func(node *graphengine.Node) error {
		if !shouldExportNode(node, filter) {
			return nil
		}
		props := filterNodeProps(node, filter)
		record := make([]string, 0, 2+len(propKeys))
		record = append(record, fmt.Sprintf("%d", node.ID))
		record = append(record, strings.Join(node.Labels, ";"))
		for _, k := range propKeys {
			record = append(record, formatPropValue(props[k]))
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
		if !shouldExportEdge(edge, filter) {
			return nil
		}
		props := filterEdgeProps(edge, filter)
		record := make([]string, 0, 4+len(propKeys))
		record = append(record, fmt.Sprintf("%d", edge.ID))
		record = append(record, edge.Label)
		record = append(record, fmt.Sprintf("%d", edge.From))
		record = append(record, fmt.Sprintf("%d", edge.To))
		for _, k := range propKeys {
			record = append(record, formatPropValue(props[k]))
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
