package neptunegraph

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage/graphengine"
	storecommon "vorpalstacks/internal/store/aws/common"
	ngstore "vorpalstacks/internal/store/aws/neptunegraph"
	"vorpalstacks/internal/utils/ntriples"
)

// CreateGraphUsingImportTask creates a new graph and initiates a bulk import task from the specified source.
func (s *NeptuneGraphService) CreateGraphUsingImportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphName := request.GetStringParam(req.Parameters, "graphName")
	if graphName == "" || strings.HasPrefix(graphName, "g-") || !graphNameRegex.MatchString(graphName) || len(graphName) > 63 {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphName")
	}

	roleArn := request.GetStringParam(req.Parameters, "roleArn")
	if roleArn == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "roleArn")
	}

	source := request.GetStringParam(req.Parameters, "source")
	if source == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "source")
	}

	region := reqCtx.GetRegion()
	graphID := generateID("g-")
	taskID := generateID("t-")
	now := time.Now().UTC()

	graph := &ngstore.Graph{
		Id:                 graphID,
		Name:               graphName,
		Arn:                s.arnBuilder.NeptuneGraph().Graph(graphID),
		Status:             "IMPORTING",
		ProvisionedMemory:  int32Ptr(128),
		ReplicaCount:       int32Ptr(1),
		DeletionProtection: request.GetBoolParam(req.Parameters, "deletionProtection"),
		PublicConnectivity: request.GetBoolParam(req.Parameters, "publicConnectivity"),
		KmsKeyIdentifier:   request.GetStringParam(req.Parameters, "kmsKeyIdentifier"),
		BuildNumber:        "1.0.20250313",
		CreateTime:         &now,
		AccountID:          s.accountID,
		Region:             region,
	}

	if vsc := parseVectorSearchConfig(req.Parameters); vsc != nil {
		graph.VectorSearchConfiguration = vsc
	}

	if request.HasParam(req.Parameters, "replicaCount") {
		rc := request.GetIntParam(req.Parameters, "replicaCount")
		if rc >= 0 && rc <= maxReplicaCount {
			graph.ReplicaCount = int32Ptr(int32(rc))
		}
	}

	if err := store.CreateGraph(graph); err != nil {
		if ngstore.IsAlreadyExists(err) {
			return nil, newConflictException("CONCURRENT_MODIFICATION")
		}
		return nil, err
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
		DeletionProtection: request.GetBoolParam(req.Parameters, "deletionProtection"),
		KmsKeyIdentifier:   request.GetStringParam(req.Parameters, "kmsKeyIdentifier"),
		PublicConnectivity: request.GetBoolParam(req.Parameters, "publicConnectivity"),
		Format:             request.GetStringParam(req.Parameters, "format"),
		ParquetType:        request.GetStringParam(req.Parameters, "parquetType"),
		BlankNodeHandling:  request.GetStringParam(req.Parameters, "blankNodeHandling"),
		FailOnError:        request.GetBoolParam(req.Parameters, "failOnError"),
		StartTime:          &now,
	}

	if request.HasParam(req.Parameters, "replicaCount") {
		rc := request.GetIntParam(req.Parameters, "replicaCount")
		task.ReplicaCount = int32Ptr(int32(rc))
	}
	if request.HasParam(req.Parameters, "minProvisionedMemory") {
		task.MinProvisionedMemory = int32Ptr(int32(request.GetIntParam(req.Parameters, "minProvisionedMemory")))
	}
	if request.HasParam(req.Parameters, "maxProvisionedMemory") {
		task.MaxProvisionedMemory = int32Ptr(int32(request.GetIntParam(req.Parameters, "maxProvisionedMemory")))
	}
	if request.HasParam(req.Parameters, "importOptions") {
		task.ImportOptions = parseImportOptions(req.Parameters)
	}
	if graph.VectorSearchConfiguration != nil {
		task.VectorSearchConfiguration = graph.VectorSearchConfiguration
	}

	if err := store.CreateImportTask(task); err != nil {
		return nil, err
	}

	s.taskWg.Add(1)
	go s.advanceImportTask(store, taskID, graphID)

	return importTaskToResponse(task), nil
}

// GetImportTask retrieves an import task by its identifier.
func (s *NeptuneGraphService) GetImportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	taskID := request.GetStringParam(req.Parameters, "taskIdentifier")
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

	return importTaskToResponse(task), nil
}

// ListImportTasks returns a paginated list of import task summaries.
func (s *NeptuneGraphService) ListImportTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := storecommon.ListOptions{
		MaxItems: clampMaxResults(request.GetIntParam(req.Parameters, "maxResults")),
		Marker:   request.GetStringParam(req.Parameters, "nextToken"),
	}

	tasks, nextToken, truncated, err := store.ListImportTasks(opts)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(tasks))
	for _, t := range tasks {
		items = append(items, importTaskSummaryToResponse(t))
	}

	result := map[string]interface{}{
		"tasks": items,
	}
	if truncated {
		result["nextToken"] = nextToken
	}
	return result, nil
}

// CancelImportTask cancels an in-progress import task, transitioning it to CANCELLED state.
func (s *NeptuneGraphService) CancelImportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	taskID := request.GetStringParam(req.Parameters, "taskIdentifier")
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

	if task.Status == "SUCCEEDED" || task.Status == "FAILED" || task.Status == "CANCELLED" {
		return importTaskSummaryToResponse(task), nil
	}

	task.Status = "CANCELLED"
	task.StatusReason = "Cancelled by user"
	if err := store.TryAdvanceImportTask(taskID, task.Status, func(t *ngstore.ImportTask) {
		t.Status = "CANCELLED"
		t.StatusReason = "Cancelled by user"
	}); err != nil {
		logs.Warn("failed to cancel import task", logs.String("taskId", taskID), logs.Err(err))
	}

	return importTaskSummaryToResponse(task), nil
}

// StartImportTask initiates a bulk import task on an existing graph from the specified source.
func (s *NeptuneGraphService) StartImportTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	roleArn := request.GetStringParam(req.Parameters, "roleArn")
	if roleArn == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "roleArn")
	}

	source := request.GetStringParam(req.Parameters, "source")
	if source == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "source")
	}

	taskID := generateID("t-")
	now := time.Now().UTC()

	task := &ngstore.ImportTask{
		TaskId:            taskID,
		GraphId:           graphID,
		Status:            "INITIALIZING",
		Source:            source,
		RoleArn:           roleArn,
		Format:            request.GetStringParam(req.Parameters, "format"),
		ParquetType:       request.GetStringParam(req.Parameters, "parquetType"),
		BlankNodeHandling: request.GetStringParam(req.Parameters, "blankNodeHandling"),
		FailOnError:       request.GetBoolParam(req.Parameters, "failOnError"),
		StartTime:         &now,
	}

	if request.HasParam(req.Parameters, "importOptions") {
		task.ImportOptions = parseImportOptions(req.Parameters)
	}

	if err := store.CreateImportTask(task); err != nil {
		return nil, err
	}

	s.taskWg.Add(1)
	go s.advanceImportTask(store, taskID, graphID)

	return importTaskToResponse(task), nil
}

func (s *NeptuneGraphService) advanceImportTask(store *ngstore.NeptuneGraphStore, taskID, graphID string) {
	defer s.taskWg.Done()
	defer func() { resilience.RecoverPanic("NeptuneGraph advanceImportTask") }()

	task, err := store.GetImportTask(taskID)
	if err != nil {
		logs.Warn("failed to get import task", logs.String("taskId", taskID), logs.Err(err))
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
				ProgressPercentage: int32Ptr(0),
				StartTime:          t.StartTime,
				TimeElapsedSeconds: int64Ptr(sinceStart),
				StatementCount:     int64Ptr(0),
				ErrorCount:         int32Ptr(1),
				Status:             stringPtr("FAILED"),
			}
		})
		if err != nil {
			logs.Warn("failed to advance import task to FAILED", logs.String("taskId", taskID), logs.Err(err))
		}
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
				ProgressPercentage: int32Ptr(0),
				StartTime:          t.StartTime,
				TimeElapsedSeconds: int64Ptr(sinceStart),
				StatementCount:     int64Ptr(0),
				ErrorCount:         int32Ptr(1),
				Status:             stringPtr("FAILED"),
			}
		})
		if err != nil {
			logs.Warn("failed to advance import task to FAILED", logs.String("taskId", taskID), logs.Err(err))
		}
		return
	}

	err = store.TryAdvanceImportTask(taskID, "INITIALIZING", func(t *ngstore.ImportTask) {
		t.Status = "IN_PROGRESS"
	})
	if err != nil {
		logs.Warn("failed to advance import task to IN_PROGRESS", logs.String("taskId", taskID), logs.Err(err))
		return
	}

	var statementCount, dictionaryCount, errorCount int64

	switch {
	case format == "" || format == "csv":
		statementCount, dictionaryCount, errorCount = s.importCSV(entry.db, filePath)
	case format == "json":
		statementCount, dictionaryCount, errorCount = s.importJSON(entry.db, filePath)
	case format == "ntriples" || format == "nquad":
		statementCount, dictionaryCount, errorCount = s.importRDF(entry.db, filePath)
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
		ProgressPercentage:   int32Ptr(100),
		StartTime:            task.StartTime,
		TimeElapsedSeconds:   int64Ptr(sinceStart),
		StatementCount:       int64Ptr(statementCount),
		DictionaryEntryCount: int64Ptr(dictionaryCount),
		ErrorCount:           int32Ptr(int32(errorCount)),
		Status:               stringPtr(finalStatus),
	}

	err = store.TryAdvanceImportTask(taskID, "IN_PROGRESS", func(t *ngstore.ImportTask) {
		t.Status = finalStatus
		t.StatusReason = statusReason
		t.ImportTaskDetails = details
	})
	if err != nil {
		logs.Warn("failed to advance import task to final state", logs.String("taskId", taskID), logs.Err(err))
		return
	}

	if finalStatus == "SUCCEEDED" {
		graph, err := store.GetGraph(graphID)
		if err == nil && graph.Status == "IMPORTING" {
			graph.Status = "AVAILABLE"
			if err := store.UpdateGraph(graph); err != nil {
				logs.Warn("Failed to update graph status to AVAILABLE", logs.Err(err))
			}
		}
	} else if finalStatus == "FAILED" {
		graph, err := store.GetGraph(graphID)
		if err == nil && graph.Status == "IMPORTING" {
			graph.Status = "AVAILABLE"
			graph.StatusReason = "Import failed"
			if err := store.UpdateGraph(graph); err != nil {
				logs.Warn("Failed to update graph status after failed import", logs.Err(err))
			}
		}
	}
}

func (s *NeptuneGraphService) importCSV(db *graphengine.DB, filePath string) (statementCount, dictionaryCount, errorCount int64) {
	f, err := os.Open(filePath)
	if err != nil {
		errorCount++
		logs.Warn("failed to open import file", logs.String("path", filePath), logs.Err(err))
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.LazyQuotes = true

	header, err := r.Read()
	if err != nil {
		errorCount++
		logs.Warn("failed to read CSV header", logs.Err(err))
		return
	}

	headerLower := make([]string, len(header))
	for i, h := range header {
		headerLower[i] = strings.ToLower(strings.TrimSpace(h))
	}

	hasFromTo := false
	for _, h := range headerLower {
		if h == "~from" || h == "~to" {
			hasFromTo = true
			break
		}
	}

	isEdge := hasFromTo

	idIndex := -1
	labelIndex := -1
	fromIndex := -1
	toIndex := -1
	for i, h := range headerLower {
		switch h {
		case "~id":
			idIndex = i
		case "~label":
			labelIndex = i
		case "~from":
			fromIndex = i
		case "~to":
			toIndex = i
		}
	}

	propIndices := make(map[string]int)
	for i, h := range headerLower {
		if strings.HasPrefix(h, "~") {
			continue
		}
		propIndices[h] = i
	}

	var nodeBatch []struct {
		Labels []string
		Props  graphengine.Props
	}
	var edgeBatch []graphengine.Edge
	extToInternal := make(map[string]graphengine.NodeID)

	flushNodeBatch := func() {
		if len(nodeBatch) == 0 {
			return
		}
		ids, err := db.AddNodeBatch(nodeBatch)
		if err != nil {
			errorCount += int64(len(nodeBatch))
			logs.Warn("AddNodeBatch failed", logs.Err(err))
		} else {
			for i, id := range ids {
				if idIndex >= 0 && i < len(nodeBatch) {
					if extID, ok := nodeBatch[i].Props["~id"].(string); ok {
						extToInternal[extID] = id
						delete(nodeBatch[i].Props, "~id")
					}
				}
			}
			dictionaryCount += int64(len(ids))
		}
		nodeBatch = nil
	}

	for {
		record, err := r.Read()
		if err != nil {
			if err.Error() != "EOF" {
				errorCount++
			}
			break
		}

		statementCount++

		if isEdge {
			flushNodeBatch()

			if fromIndex < 0 || toIndex < 0 {
				errorCount++
				continue
			}

			fromExt := strings.TrimSpace(record[fromIndex])
			toExt := strings.TrimSpace(record[toIndex])
			label := ""
			if labelIndex >= 0 && labelIndex < len(record) {
				label = strings.TrimSpace(record[labelIndex])
			}

			fromID, ok := extToInternal[fromExt]
			if !ok {
				errorCount++
				continue
			}
			toID, ok := extToInternal[toExt]
			if !ok {
				errorCount++
				continue
			}

			props := make(graphengine.Props)
			for k, idx := range propIndices {
				if idx < len(record) && record[idx] != "" {
					props[k] = strings.TrimSpace(record[idx])
				}
			}

			edgeBatch = append(edgeBatch, graphengine.Edge{
				From:  fromID,
				To:    toID,
				Label: label,
				Props: props,
			})
		} else {
			var labels []string
			if labelIndex >= 0 && labelIndex < len(record) {
				for _, l := range strings.Split(strings.TrimSpace(record[labelIndex]), ":") {
					if l != "" {
						labels = append(labels, strings.TrimSpace(l))
					}
				}
			}

			props := make(graphengine.Props)
			if idIndex >= 0 && idIndex < len(record) {
				props["~id"] = strings.TrimSpace(record[idIndex])
			}
			for k, idx := range propIndices {
				if idx < len(record) && record[idx] != "" {
					props[k] = strings.TrimSpace(record[idx])
				}
			}

			nodeBatch = append(nodeBatch, struct {
				Labels []string
				Props  graphengine.Props
			}{Labels: labels, Props: props})
		}

		if len(nodeBatch) >= 500 {
			flushNodeBatch()
		}

		if len(edgeBatch) >= 500 {
			_, err := db.AddEdgeBatch(edgeBatch)
			if err != nil {
				errorCount += int64(len(edgeBatch))
				logs.Warn("AddEdgeBatch failed", logs.Err(err))
			} else {
				dictionaryCount += int64(len(edgeBatch))
			}
			edgeBatch = nil
		}
	}

	if len(nodeBatch) > 0 {
		flushNodeBatch()
	}

	if len(edgeBatch) > 0 {
		_, err := db.AddEdgeBatch(edgeBatch)
		if err != nil {
			errorCount += int64(len(edgeBatch))
			logs.Warn("AddEdgeBatch failed", logs.Err(err))
		} else {
			dictionaryCount += int64(len(edgeBatch))
		}
	}

	return
}

func (s *NeptuneGraphService) importJSON(db *graphengine.DB, filePath string) (statementCount, dictionaryCount, errorCount int64) {
	f, err := os.Open(filePath)
	if err != nil {
		errorCount++
		logs.Warn("failed to open import file", logs.String("path", filePath), logs.Err(err))
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)

	var nodeBatch []struct {
		Labels []string
		Props  graphengine.Props
	}
	var edgeBatch []graphengine.Edge
	extToInternal := make(map[string]graphengine.NodeID)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var doc map[string]interface{}
		if err := json.Unmarshal([]byte(line), &doc); err != nil {
			errorCount++
			continue
		}

		statementCount++

		if fromVal, hasFrom := doc["~from"]; hasFrom {
			fromExt := fmt.Sprintf("%v", fromVal)
			toExt := fmt.Sprintf("%v", doc["~to"])
			label := fmt.Sprintf("%v", doc["~type"])

			fromID, ok := extToInternal[fromExt]
			if !ok {
				errorCount++
				continue
			}
			toID, ok := extToInternal[toExt]
			if !ok {
				errorCount++
				continue
			}

			props := make(graphengine.Props)
			for k, v := range doc {
				if strings.HasPrefix(k, "~") {
					continue
				}
				props[k] = v
			}

			edgeBatch = append(edgeBatch, graphengine.Edge{
				From:  fromID,
				To:    toID,
				Label: label,
				Props: props,
			})
		} else {
			var labels []string
			if lbl, ok := doc["~label"]; ok {
				switch l := lbl.(type) {
				case string:
					for _, part := range strings.Split(l, ":") {
						if part != "" {
							labels = append(labels, part)
						}
					}
				case []interface{}:
					for _, item := range l {
						if s, ok := item.(string); ok {
							labels = append(labels, s)
						}
					}
				}
			}

			props := make(graphengine.Props)
			for k, v := range doc {
				if strings.HasPrefix(k, "~") {
					continue
				}
				props[k] = v
			}

			nodeBatch = append(nodeBatch, struct {
				Labels []string
				Props  graphengine.Props
			}{Labels: labels, Props: props})
		}

		if len(nodeBatch) >= 500 {
			ids, err := db.AddNodeBatch(nodeBatch)
			if err != nil {
				errorCount += int64(len(nodeBatch))
				logs.Warn("AddNodeBatch failed", logs.Err(err))
			} else {
				for i, id := range ids {
					if i < len(nodeBatch) {
						if extID, ok := nodeBatch[i].Props["~id"].(string); ok {
							extToInternal[extID] = id
							delete(nodeBatch[i].Props, "~id")
						}
					}
				}
				dictionaryCount += int64(len(ids))
			}
			nodeBatch = nil
		}

		if len(edgeBatch) >= 500 {
			_, err := db.AddEdgeBatch(edgeBatch)
			if err != nil {
				errorCount += int64(len(edgeBatch))
				logs.Warn("AddEdgeBatch failed", logs.Err(err))
			} else {
				dictionaryCount += int64(len(edgeBatch))
			}
			edgeBatch = nil
		}
	}

	if len(nodeBatch) > 0 {
		ids, err := db.AddNodeBatch(nodeBatch)
		if err != nil {
			errorCount += int64(len(nodeBatch))
			logs.Warn("AddNodeBatch failed", logs.Err(err))
		} else {
			for i, id := range ids {
				if i < len(nodeBatch) {
					if extID, ok := nodeBatch[i].Props["~id"].(string); ok {
						extToInternal[extID] = id
						delete(nodeBatch[i].Props, "~id")
					}
				}
			}
			dictionaryCount += int64(len(ids))
		}
	}

	if len(edgeBatch) > 0 {
		_, err := db.AddEdgeBatch(edgeBatch)
		if err != nil {
			errorCount += int64(len(edgeBatch))
			logs.Warn("AddEdgeBatch failed", logs.Err(err))
		} else {
			dictionaryCount += int64(len(edgeBatch))
		}
	}

	return
}

func (s *NeptuneGraphService) importRDF(db *graphengine.DB, filePath string) (statementCount, dictionaryCount, errorCount int64) {
	f, err := os.Open(filePath)
	if err != nil {
		errorCount++
		logs.Warn("failed to open import file", logs.String("path", filePath), logs.Err(err))
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)

	extToInternal := make(map[string]graphengine.NodeID)

	var nodeBatch []struct {
		Labels []string
		Props  graphengine.Props
	}
	var edgeBatch []graphengine.Edge

	ensureNode := func(extID string) graphengine.NodeID {
		if id, ok := extToInternal[extID]; ok {
			return id
		}
		nodeID, err := db.AddNode([]string{"Resource"}, graphengine.Props{"uri": extID})
		if err != nil {
			return 0
		}
		extToInternal[extID] = nodeID
		dictionaryCount++
		return nodeID
	}

	flushNodes := func() {
		if len(nodeBatch) == 0 {
			return
		}
		ids, err := db.AddNodeBatch(nodeBatch)
		if err != nil {
			errorCount += int64(len(nodeBatch))
			logs.Warn("AddNodeBatch failed", logs.Err(err))
		} else {
			for i, id := range ids {
				if i < len(nodeBatch) {
					if uri, ok := nodeBatch[i].Props["uri"].(string); ok {
						extToInternal[uri] = id
					}
				}
			}
			dictionaryCount += int64(len(ids))
		}
		nodeBatch = nil
	}

	flushEdges := func() {
		if len(edgeBatch) == 0 {
			return
		}
		_, err := db.AddEdgeBatch(edgeBatch)
		if err != nil {
			errorCount += int64(len(edgeBatch))
			logs.Warn("AddEdgeBatch failed", logs.Err(err))
		} else {
			dictionaryCount += int64(len(edgeBatch))
		}
		edgeBatch = nil
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		subject, predicate, obj, ok := ntriples.ParseLine(line)
		if !ok {
			errorCount++
			continue
		}

		statementCount++

		subjID := ensureNode(subject)

		isURI := strings.HasPrefix(obj, "<")
		if isURI {
			objURI := strings.Trim(obj, "<>")
			objID := ensureNode(objURI)

			predLabel := ntriples.ExtractLocalName(predicate)

			edgeBatch = append(edgeBatch, graphengine.Edge{
				From:  subjID,
				To:    objID,
				Label: predLabel,
			})
		} else {
			predKey := ntriples.ExtractLocalName(predicate)

			node, err := db.GetNode(subjID)
			if err == nil && node != nil {
				updatedProps := make(graphengine.Props)
				for k, v := range node.Props {
					updatedProps[k] = v
				}
				if strings.HasPrefix(obj, "\"") {
					closing := strings.Index(obj[1:], "\"")
					if closing >= 0 {
						obj = obj[1 : closing+1]
					}
				}
				updatedProps[predKey] = obj
				_ = db.UpdateNode(subjID, updatedProps)
			}
		}

		if len(edgeBatch) >= 500 {
			flushEdges()
		}
	}

	flushNodes()
	flushEdges()
	return
}

func parseImportOptions(params map[string]interface{}) *ngstore.ImportOptions {
	v, ok := params["importOptions"]
	if !ok || v == nil {
		return nil
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	neptune, ok := m["neptune"]
	if !ok || neptune == nil {
		return nil
	}
	nm, ok := neptune.(map[string]interface{})
	if !ok {
		return nil
	}
	opts := &ngstore.ImportOptions{
		Neptune: &ngstore.NeptuneImportOptions{},
	}
	if s, ok := nm["s3ExportPath"].(string); ok {
		opts.Neptune.S3ExportPath = s
	}
	if s, ok := nm["s3ExportKmsKeyId"].(string); ok {
		opts.Neptune.S3ExportKmsKeyId = s
	}
	if b, ok := nm["preserveDefaultVertexLabels"].(bool); ok {
		opts.Neptune.PreserveDefaultVertexLabels = &b
	}
	if b, ok := nm["preserveEdgeIds"].(bool); ok {
		opts.Neptune.PreserveEdgeIds = &b
	}
	return opts
}
