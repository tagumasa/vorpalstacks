package neptunegraph

// Import task Core functions: the single validation and persistence path for
// the import task operations, plus the async import worker.

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage/graphengine"
	storecommon "vorpalstacks/internal/store/aws/common"
	ngstore "vorpalstacks/internal/store/aws/rds/neptunegraph"
)

// CreateGraphUsingImportTaskInput carries the wire-parsed
// CreateGraphUsingImportTask request. VectorSearchConfig holds the raw
// vectorSearchConfiguration wire value.
type CreateGraphUsingImportTaskInput struct {
	GraphName               string
	RoleArn                 string
	Source                  string
	Format                  string
	ParquetType             string
	BlankNodeHandling       string
	KmsKeyIdentifier        string
	DeletionProtection      bool
	PublicConnectivity      bool
	FailOnError             bool
	VectorSearchConfig      interface{}
	HasReplicaCount         bool
	ReplicaCount            int
	HasMinProvisionedMemory bool
	MinProvisionedMemory    int
	HasMaxProvisionedMemory bool
	MaxProvisionedMemory    int
	HasImportOptions        bool
	ImportOptions           *ngstore.ImportOptions
	Tags                    map[string]string
	Region                  string
}

// GetImportTaskInput carries the wire-parsed GetImportTask request.
type GetImportTaskInput struct {
	TaskIdentifier string
}

// ListImportTasksInput carries the wire-parsed ListImportTasks request with
// the already-clamped page size.
type ListImportTasksInput struct {
	MaxItems int
	Marker   string
}

// ListImportTasksResult carries the import task records for one list page.
type ListImportTasksResult struct {
	Tasks     []*ngstore.ImportTask
	NextToken string
	Truncated bool
}

// CancelImportTaskInput carries the wire-parsed CancelImportTask request.
type CancelImportTaskInput struct {
	TaskIdentifier string
}

// StartImportTaskInput carries the wire-parsed StartImportTask request.
type StartImportTaskInput struct {
	GraphIdentifier   string
	RoleArn           string
	Source            string
	Format            string
	ParquetType       string
	BlankNodeHandling string
	FailOnError       bool
	HasImportOptions  bool
	ImportOptions     *ngstore.ImportOptions
}

// createGraphUsingImportTaskCore validates the request, creates the graph and
// task records, opens the engine and starts the async import worker.
func (s *NeptuneGraphService) createGraphUsingImportTaskCore(ctx context.Context, store *ngstore.NeptuneGraphStore, in *CreateGraphUsingImportTaskInput) (*ngstore.ImportTask, error) {
	graphName := in.GraphName
	if err := validateGraphName(graphName); err != nil {
		return nil, err
	}

	roleArn := in.RoleArn
	if err := validateRoleArn(roleArn); err != nil {
		return nil, err
	}

	source := in.Source
	if source == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "source is required")
	}

	format := in.Format
	if err := validateImportFormat(format); err != nil {
		return nil, err
	}

	parquetType := in.ParquetType
	if err := validateParquetType(parquetType); err != nil {
		return nil, err
	}

	blankNodeHandling := in.BlankNodeHandling
	if err := validateBlankNodeHandling(blankNodeHandling); err != nil {
		return nil, err
	}

	graphID := generateID("g-")
	taskID := generateID("t-")
	now := time.Now().UTC()

	graph := &ngstore.Graph{
		Id:                 graphID,
		Name:               graphName,
		Arn:                s.arnBuilder.NeptuneGraph().Graph(graphID),
		Status:             "IMPORTING",
		ProvisionedMemory:  proto.Int32(128),
		ReplicaCount:       proto.Int32(1),
		DeletionProtection: in.DeletionProtection,
		PublicConnectivity: in.PublicConnectivity,
		KmsKeyIdentifier:   in.KmsKeyIdentifier,
		BuildNumber:        neptuneGraphBuildNumber,
		CreateTime:         &now,
		AccountID:          s.accountID,
		Region:             in.Region,
	}

	if vsc, err := parseVectorSearchConfigValue(in.VectorSearchConfig); err != nil {
		return nil, err
	} else if vsc != nil {
		graph.VectorSearchConfiguration = vsc
	}

	if in.HasReplicaCount {
		rc := in.ReplicaCount
		if err := validateReplicaCount(rc); err != nil {
			return nil, err
		}
		graph.ReplicaCount = proto.Int32(int32(rc))
	}

	if err := store.CreateGraph(graph); err != nil {
		if ngstore.IsAlreadyExists(err) {
			return nil, newConflictException("CONCURRENT_MODIFICATION")
		}
		return nil, err
	}

	if tags := in.Tags; len(tags) > 0 {
		if err := store.AddTags(graph.Arn, tags); err != nil {
			logs.Warn("failed to store tags for import-created graph", logs.String("graphId", graphID), logs.Err(err))
		}
	}

	bucket, err := s.graphBucket(graphID)
	if err != nil {
		return nil, newInternalServerException(err)
	}
	db, err := graphengine.New(bucket, s.engineOptions())
	if err != nil {
		logs.Warn("failed to open graph engine", logs.String("graphId", graphID), logs.Err(err))
	} else {
		s.enginesMu.Lock()
		s.activeEngines[graphID] = &engineEntry{db: db}
		s.enginesMu.Unlock()
	}

	task := &ngstore.ImportTask{
		TaskId:             taskID,
		GraphId:            graphID,
		Status:             "INITIALIZING",
		Source:             source,
		RoleArn:            roleArn,
		GraphName:          graphName,
		DeletionProtection: in.DeletionProtection,
		KmsKeyIdentifier:   in.KmsKeyIdentifier,
		PublicConnectivity: in.PublicConnectivity,
		Format:             format,
		ParquetType:        parquetType,
		BlankNodeHandling:  blankNodeHandling,
		FailOnError:        in.FailOnError,
		StartTime:          &now,
	}

	if in.HasReplicaCount {
		task.ReplicaCount = proto.Int32(int32(in.ReplicaCount))
	}
	if in.HasMinProvisionedMemory {
		mem := in.MinProvisionedMemory
		if err := validateProvisionedMemory(mem, false); err != nil {
			return nil, err
		}
		task.MinProvisionedMemory = proto.Int32(int32(mem))
	}
	if in.HasMaxProvisionedMemory {
		mem := in.MaxProvisionedMemory
		if err := validateProvisionedMemory(mem, false); err != nil {
			return nil, err
		}
		task.MaxProvisionedMemory = proto.Int32(int32(mem))
	}
	if in.HasImportOptions {
		opts := in.ImportOptions
		if err := validateImportOptions(opts); err != nil {
			return nil, err
		}
		task.ImportOptions = opts
	}
	if graph.VectorSearchConfiguration != nil {
		task.VectorSearchConfiguration = graph.VectorSearchConfiguration
	}

	if err := store.CreateImportTask(task); err != nil {
		return nil, err
	}

	s.taskWg.Add(1)
	go s.advanceImportTask(store, taskID, graphID)

	return task, nil
}

// getImportTaskCore retrieves an import task record, mapping a missing task
// to the documented ResourceNotFoundException.
func (s *NeptuneGraphService) getImportTaskCore(store *ngstore.NeptuneGraphStore, taskID string) (*ngstore.ImportTask, error) {
	if taskID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "taskIdentifier")
	}

	task, err := store.GetImportTask(taskID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("import task", taskID)
		}
		return nil, err
	}

	return task, nil
}

// listImportTasksCore reads one page of import task records, returning the
// store error unchanged so both planes keep their own error mapping.
func (s *NeptuneGraphService) listImportTasksCore(store *ngstore.NeptuneGraphStore, in *ListImportTasksInput) (*ListImportTasksResult, error) {
	tasks, nextToken, truncated, err := store.ListImportTasks(storecommon.ListOptions{
		MaxItems: in.MaxItems,
		Marker:   in.Marker,
	})
	if err != nil {
		return nil, err
	}
	return &ListImportTasksResult{Tasks: tasks, NextToken: nextToken, Truncated: truncated}, nil
}

// cancelImportTaskCore cancels an in-progress import task, transitioning it
// to CANCELLING state.
func (s *NeptuneGraphService) cancelImportTaskCore(store *ngstore.NeptuneGraphStore, in *CancelImportTaskInput) (*ngstore.ImportTask, error) {
	taskID := in.TaskIdentifier
	if taskID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "taskIdentifier")
	}

	task, err := store.GetImportTask(taskID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("import task", taskID)
		}
		return nil, err
	}

	if task.Status == "SUCCEEDED" || task.Status == "FAILED" || task.Status == "CANCELLED" || task.Status == "CANCELLING" {
		return task, nil
	}

	originalStatus := task.Status
	task.Status = "CANCELLING"
	task.StatusReason = "Cancelled by user"
	if err := store.TryAdvanceImportTask(taskID, originalStatus, func(t *ngstore.ImportTask) {
		t.Status = "CANCELLING"
		t.StatusReason = "Cancelled by user"
	}); err != nil {
		logs.Warn("failed to cancel import task", logs.String("taskId", taskID), logs.Err(err))
	}

	return task, nil
}

// startImportTaskCore validates the request and starts a bulk import task on
// an existing graph.
func (s *NeptuneGraphService) startImportTaskCore(store *ngstore.NeptuneGraphStore, in *StartImportTaskInput) (*ngstore.ImportTask, error) {
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

	roleArn := in.RoleArn
	if err := validateRoleArn(roleArn); err != nil {
		return nil, err
	}

	source := in.Source
	if source == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "source is required")
	}

	format := in.Format
	if err := validateImportFormat(format); err != nil {
		return nil, err
	}

	parquetType := in.ParquetType
	if err := validateParquetType(parquetType); err != nil {
		return nil, err
	}

	blankNodeHandling := in.BlankNodeHandling
	if err := validateBlankNodeHandling(blankNodeHandling); err != nil {
		return nil, err
	}

	taskID := generateID("t-")
	now := time.Now().UTC()

	task := &ngstore.ImportTask{
		TaskId:            taskID,
		GraphId:           graph.Id,
		Status:            "INITIALIZING",
		Source:            source,
		RoleArn:           roleArn,
		Format:            format,
		ParquetType:       parquetType,
		BlankNodeHandling: blankNodeHandling,
		FailOnError:       in.FailOnError,
		StartTime:         &now,
	}

	if in.HasImportOptions {
		opts := in.ImportOptions
		if err := validateImportOptions(opts); err != nil {
			return nil, err
		}
		task.ImportOptions = opts
	}

	if err := store.CreateImportTask(task); err != nil {
		return nil, err
	}

	s.taskWg.Add(1)
	go s.advanceImportTask(store, taskID, graph.Id)

	return task, nil
}

func (s *NeptuneGraphService) advanceImportTask(store *ngstore.NeptuneGraphStore, taskID, graphID string) {
	defer s.taskWg.Done()
	defer func() { resilience.RecoverPanic("NeptuneGraph advanceImportTask") }()

	task, err := store.GetImportTask(taskID)
	if err != nil {
		logs.Warn("failed to get import task", logs.String("taskId", taskID), logs.Err(err))
		finaliseImportedGraph(store, graphID, false, "failed to get import task")
		return
	}

	source := task.Source
	format := strings.ToLower(task.Format)

	if strings.HasPrefix(strings.ToLower(source), "s3://") {
		err = store.TryAdvanceImportTask(taskID, "INITIALIZING", func(t *ngstore.ImportTask) {
			t.Status = "FAILED"
			t.StatusReason = "S3 sources are not accessible in standalone mode"
			now := time.Now().UTC()
			sinceStart := int64(now.Sub(*t.StartTime).Seconds())
			t.ImportTaskDetails = &ngstore.ImportTaskDetails{
				ProgressPercentage: proto.Int32(0),
				StartTime:          t.StartTime,
				TimeElapsedSeconds: proto.Int64(sinceStart),
				StatementCount:     proto.Int64(0),
				ErrorCount:         proto.Int32(1),
				Status:             proto.String("FAILED"),
			}
		})
		if err != nil {
			logs.Warn("failed to advance import task to FAILED", logs.String("taskId", taskID), logs.Err(err))
		}
		finaliseImportedGraph(store, graphID, false, "Import failed")
		return
	}

	filePath := strings.TrimPrefix(source, "file://")
	if filePath == source {
		filePath = source
	}

	s.enginesMu.RLock()
	entry, ok := s.activeEngines[graphID]
	s.enginesMu.RUnlock()

	if !ok {
		logs.Warn("engine not found for graph during import", logs.String("graphId", graphID))
		err = store.TryAdvanceImportTask(taskID, "INITIALIZING", func(t *ngstore.ImportTask) {
			t.Status = "FAILED"
			t.StatusReason = "Graph engine not available"
			now := time.Now().UTC()
			sinceStart := int64(now.Sub(*t.StartTime).Seconds())
			t.ImportTaskDetails = &ngstore.ImportTaskDetails{
				ProgressPercentage: proto.Int32(0),
				StartTime:          t.StartTime,
				TimeElapsedSeconds: proto.Int64(sinceStart),
				StatementCount:     proto.Int64(0),
				ErrorCount:         proto.Int32(1),
				Status:             proto.String("FAILED"),
			}
		})
		if err != nil {
			logs.Warn("failed to advance import task to FAILED", logs.String("taskId", taskID), logs.Err(err))
		}
		finaliseImportedGraph(store, graphID, false, "Import failed")
		return
	}

	err = store.TryAdvanceImportTask(taskID, "INITIALIZING", func(t *ngstore.ImportTask) {
		t.Status = "IMPORTING"
	})
	if err != nil {
		logs.Warn("failed to advance import task to IMPORTING", logs.String("taskId", taskID), logs.Err(err))
		// Detect concurrent CancelImportTask: status may have transitioned
		// from INITIALIZING to CANCELLING before the goroutine started.
		current, getErr := store.GetImportTask(taskID)
		if getErr == nil && current.Status == "CANCELLING" {
			_ = store.TryAdvanceImportTask(taskID, "CANCELLING", func(t *ngstore.ImportTask) {
				t.Status = "CANCELLED"
			})
			finaliseImportedGraph(store, graphID, false, "Import cancelled")
		} else {
			finaliseImportedGraph(store, graphID, false, "failed to advance to IMPORTING")
		}
		return
	}

	var statementCount, dictionaryCount, errorCount int64

	switch {
	case format == "" || format == "csv":
		statementCount, dictionaryCount, errorCount = s.importCSV(entry.db, filePath)
	case format == "open_cypher":
		statementCount, dictionaryCount, errorCount = s.importCSV(entry.db, filePath)
	case format == "ntriples":
		statementCount, dictionaryCount, errorCount = s.importRDF(entry.db, filePath)
	case format == "parquet":
		errorCount = 1
		logs.Warn("parquet import not supported in standalone mode", logs.String("format", format), logs.String("taskId", taskID))
	default:
		errorCount = 1
		logs.Warn("unsupported import format", logs.String("format", format), logs.String("taskId", taskID))
	}

	finalStatus := "SUCCEEDED"
	statusReason := ""
	if errorCount > 0 && task.FailOnError {
		finalStatus = "FAILED"
		statusReason = fmt.Sprintf("%d errors during import", errorCount)
	}

	now := time.Now().UTC()
	sinceStart := int64(now.Sub(*task.StartTime).Seconds())
	details := &ngstore.ImportTaskDetails{
		ProgressPercentage:   proto.Int32(100),
		StartTime:            task.StartTime,
		TimeElapsedSeconds:   proto.Int64(sinceStart),
		StatementCount:       proto.Int64(statementCount),
		DictionaryEntryCount: proto.Int64(dictionaryCount),
		ErrorCount:           proto.Int32(int32(errorCount)),
		Status:               proto.String(finalStatus),
	}

	err = store.TryAdvanceImportTask(taskID, "IMPORTING", func(t *ngstore.ImportTask) {
		t.Status = finalStatus
		t.StatusReason = statusReason
		t.ImportTaskDetails = details
	})
	if err != nil {
		logs.Warn("failed to advance import task to final state", logs.String("taskId", taskID), logs.Err(err))
		// Detect concurrent CancelImportTask: status may have transitioned
		// from IMPORTING to CANCELLING while the import was running.
		current, getErr := store.GetImportTask(taskID)
		if getErr == nil && current.Status == "CANCELLING" {
			_ = store.TryAdvanceImportTask(taskID, "CANCELLING", func(t *ngstore.ImportTask) {
				t.Status = "CANCELLED"
			})
			finaliseImportedGraph(store, graphID, false, "Import cancelled")
		} else {
			finaliseImportedGraph(store, graphID, finalStatus == "SUCCEEDED", statusReason)
		}
		return
	}

	finaliseImportedGraph(store, graphID, finalStatus == "SUCCEEDED", statusReason)
}

// finaliseImportedGraph transitions a graph from IMPORTING to AVAILABLE or FAILED
// based on import outcome. All early-return paths in advanceImportTask must call
// this to prevent the graph from being stuck in IMPORTING state forever.
func finaliseImportedGraph(store *ngstore.NeptuneGraphStore, graphID string, succeeded bool, reason string) {
	graph, err := store.GetGraph(graphID)
	if err != nil {
		return
	}
	if graph.Status != "IMPORTING" {
		return
	}
	if succeeded {
		graph.Status = "AVAILABLE"
	} else {
		graph.Status = "FAILED"
		graph.StatusReason = reason
	}
	if err := store.UpdateGraph(graph); err != nil {
		logs.Warn("failed to update graph status after import", logs.String("graphId", graphID), logs.Err(err))
	}
}
