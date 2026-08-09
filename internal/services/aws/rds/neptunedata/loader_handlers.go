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
	"vorpalstacks/internal/core/storage/graphengine"
	pb "vorpalstacks/internal/pb/storage/storage_neptune"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
	"vorpalstacks/internal/utils/ntriples"
)

type loaderStats struct {
	totalRecords           int64
	succeeded              int64
	failed                 int64
	nodesLoaded            int64
	edgesLoaded            int64
	duplicates             int64
	parsingErrors          int64
	datatypeMismatchErrors int64
	insertErrors           int64
}

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

	if params.Source == "" {
		return nil, missingParameter("source")
	}
	if params.Format == "" {
		return nil, missingParameter("format")
	}
	if !validateLoaderFormat(params.Format) {
		return nil, invalidParameter(fmt.Sprintf("invalid format: %s (valid values: csv, opencypher, ntriples, nquads, rdfxml, turtle)", params.Format))
	}
	if params.Region == "" {
		return nil, missingParameter("region")
	}
	if params.IamRoleArn == "" {
		return nil, missingParameter("iamRoleArn")
	}
	if !validateIamRoleArn(params.IamRoleArn) {
		return nil, invalidParameter(fmt.Sprintf("iamRoleArn must be a valid IAM role ARN (arn:aws:iam::<account>:role/<name>): %s", params.IamRoleArn))
	}
	if !validateLoaderMode(params.Mode) {
		return nil, invalidParameter(fmt.Sprintf("invalid mode: %s (valid values: RESUME, NEW, AUTO)", params.Mode))
	}
	if !validateLoaderParallelism(params.Parallelism) {
		return nil, invalidParameter(fmt.Sprintf("invalid parallelism: %s (valid values: LOW, MEDIUM, HIGH, OVERSUBSCRIBE)", params.Parallelism))
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	// Validate that all dependencies exist and do not form a cycle.
	if len(params.Dependencies) > 0 {
		visited := make(map[string]bool)
		done := make(map[string]bool)
		for _, depID := range params.Dependencies {
			if err := validateDependencyChain(store, depID, visited, done); err != nil {
				return nil, err
			}
		}
	}

	loadId := generateQueryID()

	failOnError := true
	if params.FailOnError != nil {
		failOnError = *params.FailOnError
	}

	updateSingleCard := false
	if params.UpdateSingleCardinalityProperties != nil {
		updateSingleCard = *params.UpdateSingleCardinalityProperties
	}

	queueRequest := false
	if params.QueueRequest != nil {
		queueRequest = *params.QueueRequest
	}

	userProvidedEdgeIds := false
	if params.UserProvidedEdgeIds != nil {
		userProvidedEdgeIds = *params.UserProvidedEdgeIds
	}

	edgeOnlyLoad := false
	if params.EdgeOnlyLoad != nil {
		edgeOnlyLoad = *params.EdgeOnlyLoad
	}

	initialStatus := "LOAD_IN_PROGRESS"
	if queueRequest {
		initialStatus = "LOAD_QUEUED"
	}

	job := &pb.LoaderJob{
		LoadId:                            loadId,
		Status:                            initialStatus,
		Source:                            params.Source,
		Format:                            params.Format,
		SubmitTime:                        timestamppb.Now(),
		S3BucketRegion:                    params.Region,
		IamRoleArn:                        params.IamRoleArn,
		Mode:                              params.Mode,
		FailOnError:                       failOnError,
		Parallelism:                       params.Parallelism,
		ParserConfiguration:               params.ParserConfiguration,
		UpdateSingleCardinalityProperties: updateSingleCard,
		QueueRequest:                      queueRequest,
		Dependencies:                      params.Dependencies,
		UserProvidedEdgeIds:               userProvidedEdgeIds,
		EdgeOnlyLoad:                      edgeOnlyLoad,
		FullUri:                           params.Source,
		RunNumber:                         1,
		RetryNumber:                       0,
	}
	if err := store.CreateLoaderJob(job); err != nil {
		return nil, err
	}

	// Use the Neptune cluster's region (from the request context) for
	// store operations, not params.Region which is the S3 bucket region.
	// The S3 bucket region is already persisted in the job proto
	// (job.GetS3BucketRegion()) and is used by loadFromS3 for S3 calls.
	clusterRegion := reqCtx.GetRegion()
	var clusterDB *graphengine.DB
	if reader := reqCtx.GraphReader(); reader != nil {
		if db, ok := reader.(*graphengine.DB); ok {
			clusterDB = db
		}
	}

	if queueRequest {
		s.startLoaderDispatcher()
	} else {
		s.launchLoaderJob(clusterRegion, job, clusterDB)
	}

	return map[string]interface{}{
		"status": "200",
		"payload": map[string]interface{}{
			"loadId": loadId,
		},
	}, nil
}

// GetLoaderJobStatus returns the current status and statistics of a bulk load
// job in the AWS Neptune loader status response format.
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

	// Smithy query params: details, errors, page, errorsPerPage
	detailsFlag := request.GetBoolParam(req.Parameters, "details")
	errorsFlag := request.GetBoolParam(req.Parameters, "errors")
	page := request.GetIntParam(req.Parameters, "page")
	if page <= 0 {
		page = 1
	}
	errorsPerPage := request.GetIntParam(req.Parameters, "errorsPerPage")
	if errorsPerPage <= 0 {
		errorsPerPage = 10
	}

	failed := job.GetTotalErrors()
	total := job.GetTotalRecords()
	succeeded := total - failed
	if succeeded < 0 {
		succeeded = 0
	}

	var totalTimeMs float64
	if job.GetEndTime() != nil && job.GetSubmitTime() != nil {
		totalTimeMs = float64(job.GetEndTime().AsTime().Sub(job.GetSubmitTime().AsTime()).Milliseconds())
	}
	if job.GetEndTime() == nil && job.GetSubmitTime() != nil && job.GetStatus() == "LOAD_IN_PROGRESS" {
		totalTimeMs = float64(time.Since(job.GetSubmitTime().AsTime()).Milliseconds())
	}

	// Build feedCount in AWS Neptune list format: [{LOAD_COMPLETED: N}, ...]
	feedCount := buildFeedCountList(job)

	// Build overallStatus matching AWS Neptune format
	overallStatus := map[string]interface{}{
		"fullUri":                job.GetFullUri(),
		"runNumber":              int(job.GetRunNumber()),
		"retryNumber":            int(job.GetRetryNumber()),
		"status":                 job.GetStatus(),
		"totalTimeSpent":         totalTimeMs,
		"totalRecords":           total,
		"loaded":                 succeeded,
		"inserts":                succeeded,
		"errors":                 failed,
		"drops":                  0,
		"duplicates":             job.GetTotalDuplicates(),
		"parsingErrors":          job.GetParsingErrors(),
		"datatypeMismatchErrors": job.GetDatatypeMismatchErrors(),
		"insertErrors":           job.GetInsertErrors(),
	}

	// Build error structure when errors flag is set and errors exist.
	// Split the error log by newlines into individual entries and paginate
	// using the page and errorsPerPage query parameters.
	if errorsFlag && job.GetErrorLog() != "" {
		errorLines := strings.Split(strings.TrimSpace(job.GetErrorLog()), "\n")
		totalErrors := len(errorLines)
		startIdx := (page - 1) * errorsPerPage
		endIdx := startIdx + errorsPerPage
		if startIdx > totalErrors {
			startIdx = totalErrors
		}
		if endIdx > totalErrors {
			endIdx = totalErrors
		}

		logs := make([]map[string]interface{}, 0, endIdx-startIdx)
		for _, line := range errorLines[startIdx:endIdx] {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			logs = append(logs, map[string]interface{}{
				"errorCode":    "LoaderError",
				"errorMessage": line,
			})
		}

		errorEntry := map[string]interface{}{
			"startIndex": startIdx + 1,
			"endIndex":   endIdx,
			"loadId":     loadId,
			"errorLogs":  logs,
		}
		overallStatus["errors"] = []map[string]interface{}{errorEntry}
	}

	payload := map[string]interface{}{
		"loadId":    job.GetLoadId(),
		"status":    job.GetStatus(),
		"feedCount": feedCount,
	}

	// overallStatus is always included per AWS Neptune specification.
	// See: https://docs.aws.amazon.com/neptune/latest/userguide/load-api-reference-status-examples.html
	payload["overallStatus"] = overallStatus

	// failedFeeds is only included when details=true, per AWS Neptune
	// specification. Without details, the response contains feedCount and
	// overallStatus only.
	if detailsFlag && len(job.GetFailedFeeds()) > 0 {
		failedFeeds := make([]map[string]interface{}, 0, len(job.GetFailedFeeds()))
		for _, feed := range job.GetFailedFeeds() {
			failedFeeds = append(failedFeeds, map[string]interface{}{
				"source": feed,
				"status": "LOAD_FAILED",
			})
		}
		payload["failedFeeds"] = failedFeeds
	}

	return map[string]interface{}{
		"status":  "200",
		"payload": payload,
	}, nil
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
	if st == "LOAD_COMPLETED" || st == "LOAD_FAILED" || st == "CANCELLED" || st == "LOAD_CANCELLED_DUE_TO_ERROR" || st == "LOAD_FAILED_BECAUSE_DEPENDENCY_NOT_SATISFIED" {
		return nil, badRequest(fmt.Sprintf("cannot cancel loader job in terminal state: %s", st))
	}
	job.Status = "CANCELLED"
	job.EndTime = timestamppb.Now()
	if err := store.UpdateLoaderJob(job); err != nil {
		logs.Warn("failed to persist loader job cancellation", logs.String("loadId", loadId), logs.Err(err))
	}

	s.loaderMu.Lock()
	if ch, ok := s.loaderCancelChs[loadId]; ok {
		close(ch)
		delete(s.loaderCancelChs, loadId)
	}
	s.loaderMu.Unlock()

	// Trigger the dispatcher to pick up any queued jobs that may now
	// have a free slot due to this cancellation.
	s.startLoaderDispatcher()

	return map[string]interface{}{
		"status": "200",
	}, nil
}

// launchLoaderJob starts the async loader goroutine for a non-queued job.
// clusterRegion is the Neptune cluster's region (for store lookups); the S3
// bucket region is read from job.GetS3BucketRegion() inside runLoaderJob.
func (s *NeptuneDataService) launchLoaderJob(clusterRegion string, job *pb.LoaderJob, clusterDB *graphengine.DB) {
	cancelCh := make(chan struct{})
	s.loaderMu.Lock()
	s.loaderCancelChs[job.GetLoadId()] = cancelCh
	s.loaderMu.Unlock()

	s.loaderWg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("loader goroutine panicked", logs.Any("panic", r))
			}
		}()
		defer s.loaderWg.Done()
		s.runLoaderJob(clusterRegion, job, clusterDB, cancelCh)
	}()
}

// runLoaderJob executes a bulk load job asynchronously. Supports s3:// and
// https://s3.*.amazonaws.com sources with CSV and ntriples formats. The job
// status is persisted to Pebble storage on completion or failure.
// clusterRegion is the Neptune cluster's region for store operations; S3
// bucket region is obtained from job.GetS3BucketRegion() for S3 calls.
func (s *NeptuneDataService) runLoaderJob(clusterRegion string, job *pb.LoaderJob, clusterDB *graphengine.DB, cancelCh chan struct{}) {
	loadID := job.GetLoadId()
	source := job.GetSource()
	format := job.GetFormat()

	time.Sleep(100 * time.Millisecond)

	store, err := s.GetStoreForRegion(clusterRegion)
	if err != nil {
		logs.Warn("loader job failed to get store", logs.String("loadId", loadID), logs.Err(err))
		return
	}

	current, err := store.GetLoaderJob(loadID)
	if err != nil || current == nil {
		return
	}
	if current.GetStatus() == "CANCELLED" {
		return
	}

	// S3 bucket region may differ from the cluster region.
	s3Region := job.GetS3BucketRegion()

	stats := &loaderStats{}
	var loadErr string

	switch {
	case strings.HasPrefix(source, "s3://"):
		loadErr = s.loadFromS3(s3Region, job, loadID, source, format, stats, clusterDB, cancelCh, store)
	case strings.HasPrefix(source, "https://s3"), strings.HasPrefix(source, "http://s3"):
		s3uri := normalizeS3HTTPSURI(source)
		loadErr = s.loadFromS3(s3Region, job, loadID, s3uri, format, stats, clusterDB, cancelCh, store)
	default:
		loadErr = fmt.Sprintf("unsupported source URI scheme: %s (only s3:// and https://s3.*.amazonaws.com are supported)", source)
	}

	job.EndTime = timestamppb.Now()
	job.TotalRecords = stats.totalRecords
	job.TotalErrors = stats.failed
	job.TotalDuplicates = stats.duplicates
	job.ParsingErrors = stats.parsingErrors
	job.DatatypeMismatchErrors = stats.datatypeMismatchErrors
	job.InsertErrors = stats.insertErrors

	// Re-read the job from the store to detect a concurrent cancel.
	current, err = store.GetLoaderJob(loadID)
	if err == nil && current != nil && current.GetStatus() == "CANCELLED" {
		job.Status = "CANCELLED"
	} else if loadErr != "" {
		job.Status = "LOAD_FAILED"
		if job.Details == nil {
			job.Details = make(map[string]string)
		}
		job.Details["error"] = loadErr
		if job.ErrorLog == "" {
			job.ErrorLog = loadErr
		}
	} else {
		job.Status = "LOAD_COMPLETED"
	}

	job.OverallStatus = formatOverallStatus(job, stats)

	if updateErr := store.UpdateLoaderJob(job); updateErr != nil {
		logs.Warn("failed to update loader job", logs.String("loadId", loadID), logs.Err(updateErr))
	}

	s.loaderMu.Lock()
	delete(s.loaderCancelChs, loadID)
	s.loaderMu.Unlock()

	// Trigger the dispatcher to pick up any queued jobs.
	s.startLoaderDispatcher()
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

// loadFromS3 reads objects from S3 via the S3Reader invoker and delegates to
// format-specific loaders. Supports CSV, ntriples and nquads formats.
// Progress metrics are persisted per-file for real-time monitoring.
// s3Region is the S3 bucket region (may differ from the Neptune cluster region).
func (s *NeptuneDataService) loadFromS3(s3Region string, job *pb.LoaderJob, loadID, source, format string, stats *loaderStats, clusterDB *graphengine.DB, cancelCh chan struct{}, store *neptunestore.NeptuneStore) string {
	if s.s3Invoker == nil {
		return fmt.Sprintf("S3 service not available for loading from %s", source)
	}

	bucket, prefix := parseS3URI(source)

	keys, err := s.s3Invoker.ListObjects(context.Background(), s3Region, bucket, prefix, 0)
	if err != nil {
		return fmt.Sprintf("failed to list S3 objects at %s: %v", source, err)
	}
	if len(keys) == 0 {
		return fmt.Sprintf("no objects found at %s", source)
	}

	if clusterDB == nil {
		return "graph database not available"
	}

	writer := graphengine.GraphWriter(clusterDB)

	for fileIdx, key := range keys {
		select {
		case <-cancelCh:
			return "loader job cancelled"
		default:
		}

		data, err := s.s3Invoker.GetObject(context.Background(), s3Region, bucket, key, 0)
		if err != nil {
			stats.failed++
			stats.totalRecords++
			if job.GetFailOnError() {
				return fmt.Sprintf("failed to get S3 object %s/%s: %v", bucket, key, err)
			}
			logs.Warn("loader: failed to get S3 object, continuing",
				logs.String("loadId", loadID), logs.String("key", key), logs.Err(err))
			continue
		}

		tmpFile, err := os.CreateTemp("", "neptune-load-*.csv")
		if err != nil {
			return fmt.Sprintf("failed to create temp file: %v", err)
		}
		tmpPath := tmpFile.Name()

		if _, err := tmpFile.Write(data); err != nil {
			tmpFile.Close()
			os.Remove(tmpPath)
			return fmt.Sprintf("failed to write temp file: %v", err)
		}
		tmpFile.Close()

		f, err := os.Open(tmpPath)
		if err != nil {
			os.Remove(tmpPath)
			return fmt.Sprintf("failed to open temp file: %v", err)
		}

		var loadErr string
		switch format {
		case "csv":
			loadErr = s.loadCSV(f, writer, stats, cancelCh)
		case "ntriples", "nquads":
			loadErr = s.loadNTriples(f, writer, stats, cancelCh)
		case "opencypher", "rdfxml", "turtle":
			f.Close()
			os.Remove(tmpPath)
			return fmt.Sprintf("format %s is not yet supported for loading", format)
		default:
			f.Close()
			os.Remove(tmpPath)
			return fmt.Sprintf("unsupported format: %s", format)
		}
		f.Close()
		os.Remove(tmpPath)

		// Persist progress after each file for real-time monitoring.
		job.TotalRecords = stats.totalRecords
		job.TotalErrors = stats.failed
		job.TotalDuplicates = stats.duplicates
		job.ParsingErrors = stats.parsingErrors
		job.InsertErrors = stats.insertErrors
		_ = store.UpdateLoaderJob(job)

		if loadErr != "" {
			if job.GetFailOnError() {
				return loadErr
			}
			stats.failed++
			logs.Warn("loader: file error, continuing (failOnError=false)",
				logs.String("loadId", loadID), logs.Int("fileIdx", fileIdx), logs.String("key", key))
		}
	}

	return ""
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

// validateDependencyChain performs a DFS traversal of the dependency graph
// starting from depID. It verifies that each dependency exists and detects
// cycles that would cause the dispatcher to loop indefinitely.
//
// The visited map tracks nodes on the current recursion stack for cycle
// detection. The done map memoises nodes that have been fully explored so
// that diamond dependencies (A→B→C, A→C) do not re-traverse C. This keeps
// the algorithm O(V+E) instead of O(2^V) for pathological DAGs.
func validateDependencyChain(store *neptunestore.NeptuneStore, depID string, visited, done map[string]bool) error {
	if done[depID] {
		return nil
	}
	if visited[depID] {
		return invalidParameter(fmt.Sprintf("circular dependency detected involving load: %s", depID))
	}
	visited[depID] = true
	defer delete(visited, depID)

	dep, err := store.GetLoaderJob(depID)
	if err != nil {
		return internalFailure(fmt.Sprintf("failed to check dependency %s: %v", depID, err))
	}
	if dep == nil {
		return bulkLoadNotFound(depID)
	}
	for _, transitiveDep := range dep.GetDependencies() {
		if err := validateDependencyChain(store, transitiveDep, visited, done); err != nil {
			return err
		}
	}
	done[depID] = true
	return nil
}
