package neptunedata

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/utils/ntriples"
	pb "vorpalstacks/internal/pb/storage/storage_neptune"
	"vorpalstacks/internal/core/storage/graphengine"
)

type loaderStats struct {
	totalRecords int64
	succeeded    int64
	failed       int64
	nodesLoaded  int64
	edgesLoaded  int64
}

// StartLoaderJob initiates a bulk load job for loading data into the Neptune
// graph from the specified source location.
func (s *NeptuneDataService) StartLoaderJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	body := req.Body
	var params struct {
		Source      string `json:"source"`
		Format      string `json:"format"`
		IamRoleArn  string `json:"iamRoleArn"`
		Mode        string `json:"mode"`
		Parallelism string `json:"parallelism"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}
	if params.Source == "" {
		return nil, missingParameter("source")
	}
	if params.Format == "" {
		return nil, missingParameter("format")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	loadId := generateQueryID()

	job := &pb.LoaderJob{
		LoadId:     loadId,
		Status:     "LOAD_IN_PROGRESS",
		Source:     params.Source,
		Format:     params.Format,
		SubmitTime: timestamppb.Now(),
	}
	if err := store.CreateLoaderJob(job); err != nil {
		return nil, err
	}

	region := reqCtx.GetRegion()
	s.loaderWg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("loader goroutine panicked", logs.Any("panic", r))
			}
		}()
		defer s.loaderWg.Done()
		s.runLoaderJob(region, loadId, params.Source, params.Format)
	}()

	return map[string]interface{}{
		"status": "200",
		"payload": map[string]interface{}{
			"loadId": loadId,
		},
	}, nil
}

// GetLoaderJobStatus returns the current status and statistics of a bulk load
// job.
func (s *NeptuneDataService) GetLoaderJobStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	loadId := getPathParam(req, "loadId")
	if loadId == "" {
		return nil, missingParameter("loadId")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	job, err := store.GetLoaderJob(loadId)
	if err != nil || job == nil {
		return nil, bulkLoadNotFound(loadId)
	}

	failed := job.GetTotalErrors()
	total := job.GetTotalRecords()

	var totalTimeMs float64
	if job.GetEndTime() != nil && job.GetSubmitTime() != nil {
		totalTimeMs = float64(job.GetEndTime().AsTime().Sub(job.GetSubmitTime().AsTime()).Milliseconds())
	}

	payload := map[string]interface{}{
		"loadId": job.GetLoadId(),
		"status": job.GetStatus(),
		"feedCount": map[string]interface{}{
			"total":     1,
			"succeeded": mapBool(total == 0 || failed == 0),
			"failed":    mapBool(failed > 0),
		},
		"overallStatus": map[string]interface{}{
			"totalRecords": total,
			"totalTime":    totalTimeMs,
		},
	}

	return map[string]interface{}{
		"status":  "200",
		"payload": payload,
	}, nil
}

func mapBool(b bool) int {
	if b {
		return 1
	}
	return 0
}

// ListLoaderJobs returns the load IDs of all submitted bulk load jobs,
// optionally including queued loads.
func (s *NeptuneDataService) ListLoaderJobs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	includeQueuedLoads := request.GetBoolParam(req.Parameters, "includeQueuedLoads")
	limit := request.GetIntParam(req.Parameters, "limit")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	jobs, err := store.ListLoaderJobs()
	if err != nil {
		return nil, err
	}

	loadIds := make([]string, 0, len(jobs))
	for _, job := range jobs {
		if !includeQueuedLoads && job.GetStatus() == "LOAD_QUEUED" {
			continue
		}
		loadIds = append(loadIds, job.GetLoadId())
	}

	if limit > 0 && len(loadIds) > limit {
		loadIds = loadIds[:limit]
	}

	return map[string]interface{}{
		"status": "200",
		"payload": map[string]interface{}{
			"loadIds": loadIds,
		},
	}, nil
}

// CancelLoaderJob cancels a running or queued bulk load job and marks its
// status as cancelled.
func (s *NeptuneDataService) CancelLoaderJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	loadId := getPathParam(req, "loadId")
	if loadId == "" {
		return nil, missingParameter("loadId")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	job, err := store.GetLoaderJob(loadId)
	if err != nil || job == nil {
		return nil, bulkLoadNotFound(loadId)
	}
	st := job.GetStatus()
	if st == "LOAD_COMPLETED" || st == "LOAD_FAILED" || st == "CANCELLED" {
		return nil, badRequest(fmt.Sprintf("cannot cancel loader job in terminal state: %s", st))
	}
	job.Status = "CANCELLED"
	job.EndTime = timestamppb.Now()
	if err := store.UpdateLoaderJob(job); err != nil {
		logs.Warn("failed to persist loader job cancellation", logs.String("loadId", loadId), logs.Err(err))
	}

	return map[string]interface{}{
		"status": "200",
	}, nil
}

// runLoaderJob executes a bulk load job asynchronously. Supports file:// sources
// with CSV and ntriples formats. S3 sources fail with an appropriate error in
// standalone mode. The job status is persisted to Pebble storage on completion
// or failure.
func (s *NeptuneDataService) runLoaderJob(region, loadID, source, format string) {
	time.Sleep(100 * time.Millisecond)

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		logs.Warn("loader job failed to get store", logs.String("loadId", loadID), logs.Err(err))
		return
	}

	job, err := store.GetLoaderJob(loadID)
	if err != nil || job == nil {
		return
	}

	if job.GetStatus() == "CANCELLED" {
		return
	}

	stats := &loaderStats{}
	var loadErr string

	switch {
	case strings.HasPrefix(source, "s3://"):
		loadErr = fmt.Sprintf("source %s is not accessible in standalone mode", source)
	case strings.HasPrefix(source, "file://"):
		loadErr = s.loadFromFile(job, loadID, source, format, stats)
	default:
		loadErr = fmt.Sprintf("unsupported source URI scheme: %s", source)
	}

	job.EndTime = timestamppb.Now()
	job.TotalRecords = stats.totalRecords
	job.TotalErrors = stats.failed

	if loadErr != "" {
		job.Status = "LOAD_FAILED"
		if job.Details == nil {
			job.Details = make(map[string]string)
		}
		job.Details["error"] = loadErr
	} else {
		job.Status = "LOAD_COMPLETED"
		job.OverallStatus = fmt.Sprintf("%d records loaded (%d nodes, %d edges)", stats.totalRecords, stats.nodesLoaded, stats.edgesLoaded)
		if stats.failed > 0 {
			job.OverallStatus = fmt.Sprintf("%d records loaded, %d errors (%d nodes, %d edges)", stats.succeeded, stats.failed, stats.nodesLoaded, stats.edgesLoaded)
		}
	}

	if updateErr := store.UpdateLoaderJob(job); updateErr != nil {
		logs.Warn("failed to update loader job", logs.String("loadId", loadID), logs.Err(updateErr))
	}
}

func (s *NeptuneDataService) loadFromFile(job *pb.LoaderJob, loadID, source, format string, stats *loaderStats) string {
	filePath := strings.TrimPrefix(source, "file://")
	if filePath == "" {
		return "empty file path in file:// source"
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Sprintf("failed to open file %s: %v", filePath, err)
	}
	defer f.Close()

	if s.graphDB == nil {
		return "graph database not available"
	}

	writer := graphengine.GraphWriter(s.graphDB)

	switch strings.ToLower(format) {
	case "csv":
		return s.loadCSV(f, writer, stats)
	case "ntriples", "ntriplesrdf":
		return s.loadNTriples(f, writer, stats)
	default:
		return fmt.Sprintf("unsupported format: %s", format)
	}
}

func (s *NeptuneDataService) loadCSV(f *os.File, writer graphengine.GraphWriter, stats *loaderStats) string {
	r := csv.NewReader(f)
	r.LazyQuotes = true
	r.FieldsPerRecord = -1

	header, err := r.Read()
	if err != nil {
		return fmt.Sprintf("failed to read CSV header: %v", err)
	}

	if len(header) == 0 {
		return "empty CSV header"
	}

	hasFromTo := false
	idIdx := -1
	labelIdx := -1
	fromIdx := -1
	toIdx := -1
	propIndices := make(map[int]string)

	for i, col := range header {
		col = strings.TrimSpace(col)
		switch col {
		case "~id":
			idIdx = i
		case "~label":
			labelIdx = i
		case "~from":
			fromIdx = i
			hasFromTo = true
		case "~to":
			toIdx = i
			hasFromTo = true
		default:
			if col != "" {
				propIndices[i] = col
			}
		}
	}

	if hasFromTo {
		return s.loadCSVEdges(r, writer, stats, idIdx, labelIdx, fromIdx, toIdx, propIndices)
	}
	return s.loadCSVNodes(r, writer, stats, idIdx, labelIdx, propIndices)
}

func (s *NeptuneDataService) loadCSVNodes(r *csv.Reader, writer graphengine.GraphWriter, stats *loaderStats, idIdx, labelIdx int, propIndices map[int]string) string {
	idMap := make(map[string]graphengine.NodeID)
	type batchEntry struct {
		Labels []string
		Props  graphengine.Props
		OrigID string
	}
	var batch []batchEntry
	batchSize := 500

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			stats.failed++
			stats.totalRecords++
			continue
		}

		stats.totalRecords++

		var nodeID string
		if idIdx >= 0 && idIdx < len(record) {
			nodeID = strings.TrimSpace(record[idIdx])
		}

		labels := []string{}
		if labelIdx >= 0 && labelIdx < len(record) {
			for _, l := range strings.Split(record[labelIdx], ";") {
				l = strings.TrimSpace(l)
				if l != "" {
					labels = append(labels, l)
				}
			}
		}

		props := make(graphengine.Props)
		for idx, name := range propIndices {
			if idx < len(record) {
				val := strings.TrimSpace(record[idx])
				props[name] = parseCSVValue(val)
			}
		}

		batch = append(batch, batchEntry{Labels: labels, Props: props, OrigID: nodeID})

		if len(batch) >= batchSize {
			items := make([]struct {
				Labels []string
				Props  graphengine.Props
			}, len(batch))
			for i, e := range batch {
				items[i] = struct {
					Labels []string
					Props  graphengine.Props
				}{Labels: e.Labels, Props: e.Props}
			}
			ids, err := writer.AddNodeBatch(items)
			if err != nil {
				stats.failed += int64(len(batch))
				logs.Warn("failed to add node batch", logs.Err(err))
			} else {
				stats.succeeded += int64(len(ids))
				stats.nodesLoaded += int64(len(ids))
				for i, assignedID := range ids {
					if i < len(batch) && batch[i].OrigID != "" {
						idMap[batch[i].OrigID] = assignedID
					}
				}
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		items := make([]struct {
			Labels []string
			Props  graphengine.Props
		}, len(batch))
		for i, e := range batch {
			items[i] = struct {
				Labels []string
				Props  graphengine.Props
			}{Labels: e.Labels, Props: e.Props}
		}
		ids, err := writer.AddNodeBatch(items)
		if err != nil {
			stats.failed += int64(len(batch))
			logs.Warn("failed to add final node batch", logs.Err(err))
		} else {
			stats.succeeded += int64(len(ids))
			stats.nodesLoaded += int64(len(ids))
			for i, assignedID := range ids {
				if i < len(batch) && batch[i].OrigID != "" {
					idMap[batch[i].OrigID] = assignedID
				}
			}
		}
	}

	return ""
}

func (s *NeptuneDataService) loadCSVEdges(r *csv.Reader, writer graphengine.GraphWriter, stats *loaderStats, idIdx, labelIdx, fromIdx, toIdx int, propIndices map[int]string) string {
	if fromIdx < 0 || toIdx < 0 {
		return "edge CSV requires ~from and ~to columns"
	}

	batch := make([]graphengine.Edge, 0, 500)
	batchSize := 500

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			stats.failed++
			stats.totalRecords++
			continue
		}

		stats.totalRecords++

		fromStr := ""
		toStr := ""
		edgeLabel := ""

		if fromIdx < len(record) {
			fromStr = strings.TrimSpace(record[fromIdx])
		}
		if toIdx < len(record) {
			toStr = strings.TrimSpace(record[toIdx])
		}
		if labelIdx >= 0 && labelIdx < len(record) {
			edgeLabel = strings.TrimSpace(record[labelIdx])
		}

		props := make(graphengine.Props)
		for idx, name := range propIndices {
			if idx < len(record) {
				val := strings.TrimSpace(record[idx])
				props[name] = parseCSVValue(val)
			}
		}

		edgeID := graphengine.NodeID(0)
		if idIdx >= 0 && idIdx < len(record) {
			val := strings.TrimSpace(record[idIdx])
			if n, err := strconv.ParseUint(val, 10, 64); err == nil {
				edgeID = graphengine.NodeID(n)
			}
		}

		fromNodeID := graphengine.NodeID(0)
		if n, err := strconv.ParseUint(fromStr, 10, 64); err == nil {
			fromNodeID = graphengine.NodeID(n)
		}
		toNodeID := graphengine.NodeID(0)
		if n, err := strconv.ParseUint(toStr, 10, 64); err == nil {
			toNodeID = graphengine.NodeID(n)
		}

		batch = append(batch, graphengine.Edge{
			ID:    graphengine.EdgeID(edgeID),
			From:  fromNodeID,
			To:    toNodeID,
			Label: edgeLabel,
			Props: props,
		})

		if len(batch) >= batchSize {
			ids, err := writer.AddEdgeBatch(batch)
			if err != nil {
				stats.failed += int64(len(batch))
				logs.Warn("failed to add edge batch", logs.Err(err))
			} else {
				stats.succeeded += int64(len(ids))
				stats.edgesLoaded += int64(len(ids))
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		ids, err := writer.AddEdgeBatch(batch)
		if err != nil {
			stats.failed += int64(len(batch))
			logs.Warn("failed to add final edge batch", logs.Err(err))
		} else {
			stats.succeeded += int64(len(ids))
			stats.edgesLoaded += int64(len(ids))
		}
	}

	return ""
}

func (s *NeptuneDataService) loadNTriples(f *os.File, writer graphengine.GraphWriter, stats *loaderStats) string {
	type pendingEdgeEntry struct {
		fromExt string
		toExt   string
		label   string
	}

	idMap := make(map[string]graphengine.NodeID)
	pendingURIs := make(map[string]bool)
	var batch []struct {
		Labels []string
		Props  graphengine.Props
	}
	batchSize := 500

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	lastSubject := ""
	var pendingEdges []pendingEdgeEntry

	flushBatch := func() {
		if len(batch) > 0 {
			ids, err := writer.AddNodeBatch(batch)
			if err != nil {
				stats.failed += int64(len(batch))
				logs.Warn("failed to add ntriples node batch", logs.Err(err))
			} else {
				stats.succeeded += int64(len(ids))
				stats.nodesLoaded += int64(len(ids))
				for i, id := range ids {
					if i < len(batch) && batch[i].Props != nil {
						if origID, ok := batch[i].Props["__uri"].(string); ok {
							idMap[origID] = id
							delete(batch[i].Props, "__uri")
							delete(pendingURIs, origID)
						}
					}
				}
			}
			batch = batch[:0]
		}
	}

	flushEdges := func() {
		if len(pendingEdges) == 0 {
			return
		}
		edges := make([]graphengine.Edge, 0, len(pendingEdges))
		for _, pe := range pendingEdges {
			fromID, ok := idMap[pe.fromExt]
			if !ok {
				continue
			}
			toID, ok := idMap[pe.toExt]
			if !ok {
				continue
			}
			edges = append(edges, graphengine.Edge{
				From:  fromID,
				To:    toID,
				Label: pe.label,
			})
		}
		if len(edges) > 0 {
			ids, err := writer.AddEdgeBatch(edges)
			if err != nil {
				stats.failed += int64(len(edges))
				logs.Warn("failed to add ntriples edge batch", logs.Err(err))
			} else {
				stats.succeeded += int64(len(ids))
				stats.edgesLoaded += int64(len(ids))
			}
		}
		pendingEdges = pendingEdges[:0]
	}

	ensureNode := func(uri string, labels []string) {
		if _, exists := idMap[uri]; exists {
			return
		}
		if pendingURIs[uri] {
			return
		}
		pendingURIs[uri] = true
		props := graphengine.Props{"__uri": uri}
		if strings.HasPrefix(uri, "\"") {
			props["value"] = strings.Trim(uri, "\"")
		}
		batch = append(batch, struct {
			Labels []string
			Props  graphengine.Props
		}{Labels: labels, Props: props})
		if len(batch) >= batchSize {
			flushBatch()
		}
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		subj, pred, obj, ok := ntriples.ParseLine(line)
		if !ok {
			stats.failed++
			stats.totalRecords++
			continue
		}

		stats.totalRecords++

		if lastSubject != "" && subj != lastSubject {
			flushBatch()
			flushEdges()
		}
		lastSubject = subj

		ensureNode(subj, []string{"Resource"})

		isLiteral := !strings.HasPrefix(obj, "<")
		objLabels := []string{"Resource"}
		if isLiteral {
			objLabels = []string{"Literal"}
		}
		ensureNode(obj, objLabels)

		predLabel := ntriples.ExtractLocalName(pred)
		pendingEdges = append(pendingEdges, pendingEdgeEntry{
			fromExt: subj,
			toExt:   obj,
			label:   predLabel,
		})
	}

	flushBatch()
	flushEdges()

	return ""
}

func parseCSVValue(val string) interface{} {
	if val == "" {
		return nil
	}
	if strings.EqualFold(val, "true") {
		return true
	}
	if strings.EqualFold(val, "false") {
		return false
	}
	if n, err := strconv.ParseInt(val, 10, 64); err == nil {
		return n
	}
	if n, err := strconv.ParseFloat(val, 64); err == nil {
		return n
	}
	return val
}
