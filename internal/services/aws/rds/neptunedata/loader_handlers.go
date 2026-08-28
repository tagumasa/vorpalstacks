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

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage/graphengine"
	pb "vorpalstacks/internal/pb/storage/storage_neptune"
	"vorpalstacks/internal/utils/ntriples"
)

// StartLoaderJob initiates a bulk load job for loading data into the Neptune
// graph from the specified source location.
func (s *NeptuneDataService) StartLoaderJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	body := req.Body
	var params struct {
		Source                            string            `json:"source"`
		Format                            string            `json:"format"`
		Region                            string            `json:"region"`
		IamRoleArn                        string            `json:"iamRoleArn"`
		Mode                              string            `json:"mode"`
		FailOnError                       *bool             `json:"failOnError"`
		Parallelism                       string            `json:"parallelism"`
		ParserConfiguration               map[string]string `json:"parserConfiguration"`
		UpdateSingleCardinalityProperties *bool             `json:"updateSingleCardinalityProperties"`
		QueueRequest                      *bool             `json:"queueRequest"`
		Dependencies                      []string          `json:"dependencies"`
		UserProvidedEdgeIds               *bool             `json:"userProvidedEdgeIds"`
		EdgeOnlyLoad                      *bool             `json:"edgeOnlyLoad"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}

	var clusterDB *graphengine.DB
	if reader := reqCtx.GraphReader(); reader != nil {
		if db, ok := reader.(*graphengine.DB); ok {
			clusterDB = db
		}
	}

	return s.startLoaderJobCore(&StartLoaderJobInput{
		Source:                            params.Source,
		Format:                            params.Format,
		S3BucketRegion:                    params.Region,
		IamRoleArn:                        params.IamRoleArn,
		Mode:                              params.Mode,
		Parallelism:                       params.Parallelism,
		FailOnError:                       params.FailOnError,
		ParserConfiguration:               params.ParserConfiguration,
		UpdateSingleCardinalityProperties: params.UpdateSingleCardinalityProperties,
		QueueRequest:                      params.QueueRequest,
		Dependencies:                      params.Dependencies,
		UserProvidedEdgeIds:               params.UserProvidedEdgeIds,
		EdgeOnlyLoad:                      params.EdgeOnlyLoad,
		ClusterRegion:                     reqCtx.GetRegion(),
		ClusterDB:                         clusterDB,
	})
}

// GetLoaderJobStatus returns the current status and statistics of a bulk load
// job in the AWS Neptune loader status response format.
func (s *NeptuneDataService) GetLoaderJobStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	// Smithy query params: details, errors, page, errorsPerPage
	return s.getLoaderJobStatusCore(&GetLoaderJobStatusInput{
		LoadId:           getPathParam(req, "loadId"),
		Details:          request.GetBoolParam(req.Parameters, "details"),
		Errors:           request.GetBoolParam(req.Parameters, "errors"),
		HasPage:          request.HasParam(req.Parameters, "page"),
		Page:             request.GetIntParam(req.Parameters, "page"),
		HasErrorsPerPage: request.HasParam(req.Parameters, "errorsPerPage"),
		ErrorsPerPage:    request.GetIntParam(req.Parameters, "errorsPerPage"),
		Region:           reqCtx.GetRegion(),
	})
}

// buildFeedCountList creates the AWS Neptune loader feedCount list.
// Format: [{LOAD_COMPLETED: N}] or [{LOAD_IN_PROGRESS: N}] etc.
func buildFeedCountList(job *pb.LoaderJob) []map[string]interface{} {
	status := job.GetStatus()
	count := int64(1)
	if status == "LOAD_COMPLETED" {
		count = 1
	}
	entry := map[string]interface{}{}
	entry[status] = count
	return []map[string]interface{}{entry}
}

// ListLoaderJobs returns the load IDs of all submitted bulk load jobs,
// optionally including queued loads.
func (s *NeptuneDataService) ListLoaderJobs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	return s.listLoaderJobsCore(&ListLoaderJobsInput{
		IncludeQueuedLoads: request.GetBoolParam(req.Parameters, "includeQueuedLoads"),
		HasLimit:           request.HasParam(req.Parameters, "limit"),
		Limit:              request.GetIntParam(req.Parameters, "limit"),
		Region:             reqCtx.GetRegion(),
	})
}

// CancelLoaderJob cancels a running or queued bulk load job and marks its
// status as cancelled.
func (s *NeptuneDataService) CancelLoaderJob(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	return s.cancelLoaderJobCore(&CancelLoaderJobInput{
		LoadId: getPathParam(req, "loadId"),
		Region: reqCtx.GetRegion(),
	})
}

// formatOverallStatus builds the human-readable overall status string
// matching the AWS Neptune loader status response format.
func formatOverallStatus(job *pb.LoaderJob, stats *loaderStats) string {
	if job.GetStatus() == "LOAD_FAILED" {
		return fmt.Sprintf("LOAD_FAILED: %s", job.GetErrorLog())
	}
	if stats.failed > 0 {
		return fmt.Sprintf("%d records loaded, %d errors (%d nodes, %d edges)",
			stats.succeeded, stats.failed, stats.nodesLoaded, stats.edgesLoaded)
	}
	return fmt.Sprintf("%d records loaded (%d nodes, %d edges)",
		stats.totalRecords, stats.nodesLoaded, stats.edgesLoaded)
}

// normalizeS3HTTPSURI converts an HTTPS S3 URL to the canonical s3://
// format. Accepts both path-style (s3.region.amazonaws.com/bucket/key)
// and virtual-host-style (bucket.s3.region.amazonaws.com/key) URLs.
func normalizeS3HTTPSURI(uri string) string {
	uri = strings.TrimPrefix(uri, "https://")
	uri = strings.TrimPrefix(uri, "http://")

	slashIdx := strings.Index(uri, "/")
	if slashIdx < 0 {
		return "s3://" + uri
	}

	host := uri[:slashIdx]
	path := uri[slashIdx+1:]

	if strings.HasPrefix(host, "s3.") {
		return "s3://" + path
	}

	if strings.HasPrefix(host, "s3-") {
		return "s3://" + path
	}

	dotIdx := strings.Index(host, ".s3")
	if dotIdx > 0 {
		bucket := host[:dotIdx]
		return "s3://" + bucket + "/" + path
	}

	return "s3://" + path
}

// parseS3URI extracts bucket and prefix from an s3:// URI.
func parseS3URI(uri string) (bucket, prefix string) {
	uri = strings.TrimPrefix(uri, "s3://")
	parts := strings.SplitN(uri, "/", 2)
	bucket = parts[0]
	if len(parts) > 1 {
		prefix = parts[1]
	}
	return
}

func (s *NeptuneDataService) loadCSV(f *os.File, writer graphengine.GraphWriter, stats *loaderStats, cancelCh chan struct{}) string {
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

	idMap := make(map[string]graphengine.NodeID)
	if hasFromTo {
		return s.loadCSVEdges(r, writer, stats, idIdx, labelIdx, fromIdx, toIdx, propIndices, idMap, cancelCh)
	}
	return s.loadCSVNodes(r, writer, stats, idIdx, labelIdx, propIndices, cancelCh)
}

type csvNodeEntry struct {
	Labels []string
	Props  graphengine.Props
	OrigID string
}

func flushCSVNodeBatch(writer graphengine.GraphWriter, batch []csvNodeEntry, stats *loaderStats, idMap map[string]graphengine.NodeID) []csvNodeEntry {
	if len(batch) == 0 {
		return batch
	}
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
	return batch[:0]
}

func resolveNodeID(str string, idMap map[string]graphengine.NodeID) graphengine.NodeID {
	if idMap != nil {
		if nid, ok := idMap[str]; ok {
			return nid
		}
	}
	if n, err := strconv.ParseUint(str, 10, 64); err == nil {
		return graphengine.NodeID(n)
	}
	return graphengine.NodeID(0)
}

func (s *NeptuneDataService) loadCSVNodes(r *csv.Reader, writer graphengine.GraphWriter, stats *loaderStats, idIdx, labelIdx int, propIndices map[int]string, cancelCh chan struct{}) string {
	idMap := make(map[string]graphengine.NodeID)
	var batch []csvNodeEntry
	batchSize := 500

	for {
		select {
		case <-cancelCh:
			return "loader job cancelled"
		default:
		}

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

		batch = append(batch, csvNodeEntry{Labels: labels, Props: props, OrigID: nodeID})

		if len(batch) >= batchSize {
			batch = flushCSVNodeBatch(writer, batch, stats, idMap)
		}
	}

	if len(batch) > 0 {
		flushCSVNodeBatch(writer, batch, stats, idMap)
	}

	return ""
}

func (s *NeptuneDataService) loadCSVEdges(r *csv.Reader, writer graphengine.GraphWriter, stats *loaderStats, idIdx, labelIdx, fromIdx, toIdx int, propIndices map[int]string, idMap map[string]graphengine.NodeID, cancelCh chan struct{}) string {
	if fromIdx < 0 || toIdx < 0 {
		return "edge CSV requires ~from and ~to columns"
	}

	batch := make([]graphengine.Edge, 0, 500)
	batchSize := 500

	for {
		select {
		case <-cancelCh:
			return "loader job cancelled"
		default:
		}

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

		fromNodeID := resolveNodeID(fromStr, idMap)
		toNodeID := resolveNodeID(toStr, idMap)

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

func (s *NeptuneDataService) loadNTriples(f *os.File, writer graphengine.GraphWriter, stats *loaderStats, cancelCh chan struct{}) string {
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
		select {
		case <-cancelCh:
			return "loader job cancelled"
		default:
		}

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
