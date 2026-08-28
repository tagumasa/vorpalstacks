package neptunedata

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage/graphengine"
	pb "vorpalstacks/internal/pb/storage/storage_neptune"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
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

// Loader job bounds from the neptune-data model: the ListLoaderJobs limit is
// a positive integer no greater than 100, which is also the documented
// default when the member is omitted; the GetLoaderJobStatus pagination
// members are PositiveInteger values (minimum 1, defaults 1 and 10).
const (
	MinListLoaderJobsLimit = 1
	MaxListLoaderJobsLimit = 100
	MinLoaderPage          = 1
	MinLoaderErrorsPerPage = 1
)

// StartLoaderJobInput carries the parsed StartLoaderJob wire payload plus the
// request-context state the loader Core needs. The pointer members preserve
// wire presence so the Core can apply the documented defaults.
type StartLoaderJobInput struct {
	Source                            string
	Format                            string
	S3BucketRegion                    string
	IamRoleArn                        string
	Mode                              string
	Parallelism                       string
	FailOnError                       *bool
	ParserConfiguration               map[string]string
	UpdateSingleCardinalityProperties *bool
	QueueRequest                      *bool
	Dependencies                      []string
	UserProvidedEdgeIds               *bool
	EdgeOnlyLoad                      *bool
	// ClusterRegion is the Neptune cluster's region (from the request
	// context); S3BucketRegion above is the S3 bucket region from the wire.
	ClusterRegion string
	// ClusterDB is the cluster's graph engine, extracted from the request
	// context graph reader by the handler.
	ClusterDB *graphengine.DB
}

// startLoaderJobCore validates the load request, persists the loader job and
// launches it (directly or through the queued dispatcher).
func (s *NeptuneDataService) startLoaderJobCore(in *StartLoaderJobInput) (map[string]interface{}, error) {
	if in.Source == "" {
		return nil, missingParameter("source")
	}
	if in.Format == "" {
		return nil, missingParameter("format")
	}
	if !validateLoaderFormat(in.Format) {
		return nil, invalidParameter(fmt.Sprintf("invalid format: %s (valid values: csv, opencypher, ntriples, nquads, rdfxml, turtle)", in.Format))
	}
	if in.S3BucketRegion == "" {
		return nil, missingParameter("region")
	}
	if in.IamRoleArn == "" {
		return nil, missingParameter("iamRoleArn")
	}
	if !validateIamRoleArn(in.IamRoleArn) {
		return nil, invalidParameter(fmt.Sprintf("iamRoleArn must be a valid IAM role ARN (arn:aws:iam::<account>:role/<name>): %s", in.IamRoleArn))
	}
	if !validateLoaderMode(in.Mode) {
		return nil, invalidParameter(fmt.Sprintf("invalid mode: %s (valid values: RESUME, NEW, AUTO)", in.Mode))
	}
	if !validateLoaderParallelism(in.Parallelism) {
		return nil, invalidParameter(fmt.Sprintf("invalid parallelism: %s (valid values: LOW, MEDIUM, HIGH, OVERSUBSCRIBE)", in.Parallelism))
	}

	store, err := s.GetStoreForRegion(in.ClusterRegion)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	// Validate that all dependencies exist and do not form a cycle.
	if len(in.Dependencies) > 0 {
		visited := make(map[string]bool)
		done := make(map[string]bool)
		for _, depID := range in.Dependencies {
			if err := validateDependencyChain(store, depID, visited, done); err != nil {
				return nil, err
			}
		}
	}

	loadId := generateQueryID()

	failOnError := true
	if in.FailOnError != nil {
		failOnError = *in.FailOnError
	}

	updateSingleCard := false
	if in.UpdateSingleCardinalityProperties != nil {
		updateSingleCard = *in.UpdateSingleCardinalityProperties
	}

	queueRequest := false
	if in.QueueRequest != nil {
		queueRequest = *in.QueueRequest
	}

	userProvidedEdgeIds := false
	if in.UserProvidedEdgeIds != nil {
		userProvidedEdgeIds = *in.UserProvidedEdgeIds
	}

	edgeOnlyLoad := false
	if in.EdgeOnlyLoad != nil {
		edgeOnlyLoad = *in.EdgeOnlyLoad
	}

	initialStatus := "LOAD_IN_PROGRESS"
	if queueRequest {
		initialStatus = "LOAD_QUEUED"
	}

	job := &pb.LoaderJob{
		LoadId:                            loadId,
		Status:                            initialStatus,
		Source:                            in.Source,
		Format:                            in.Format,
		SubmitTime:                        timestamppb.Now(),
		S3BucketRegion:                    in.S3BucketRegion,
		IamRoleArn:                        in.IamRoleArn,
		Mode:                              in.Mode,
		FailOnError:                       failOnError,
		Parallelism:                       in.Parallelism,
		ParserConfiguration:               in.ParserConfiguration,
		UpdateSingleCardinalityProperties: updateSingleCard,
		QueueRequest:                      queueRequest,
		Dependencies:                      in.Dependencies,
		UserProvidedEdgeIds:               userProvidedEdgeIds,
		EdgeOnlyLoad:                      edgeOnlyLoad,
		FullUri:                           in.Source,
		RunNumber:                         1,
		RetryNumber:                       0,
	}
	if err := store.CreateLoaderJob(job); err != nil {
		return nil, err
	}

	// Use the Neptune cluster's region (from the request context) for
	// store operations, not the S3 bucket region. The S3 bucket region is
	// already persisted in the job proto (job.GetS3BucketRegion()) and is
	// used by loadFromS3 for S3 calls.
	clusterRegion := in.ClusterRegion
	clusterDB := in.ClusterDB

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

// GetLoaderJobStatusInput carries the parsed GetLoaderJobStatus wire payload.
// The Has* flags distinguish an explicitly sent value from an omitted member
// so the Core can reject out-of-range values instead of defaulting them.
type GetLoaderJobStatusInput struct {
	LoadId           string
	Details          bool
	Errors           bool
	HasPage          bool
	Page             int
	HasErrorsPerPage bool
	ErrorsPerPage    int
	Region           string
}

// getLoaderJobStatusCore loads the job record and builds the AWS Neptune
// loader status response.
func (s *NeptuneDataService) getLoaderJobStatusCore(in *GetLoaderJobStatusInput) (map[string]interface{}, error) {
	if in.LoadId == "" {
		return nil, missingParameter("loadId")
	}
	if in.HasPage && in.Page < MinLoaderPage {
		return nil, invalidParameter(fmt.Sprintf("page must be a positive integer: %d", in.Page))
	}
	if in.HasErrorsPerPage && in.ErrorsPerPage < MinLoaderErrorsPerPage {
		return nil, invalidParameter(fmt.Sprintf("errorsPerPage must be a positive integer: %d", in.ErrorsPerPage))
	}

	store, err := s.GetStoreForRegion(in.Region)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	job, err := store.GetLoaderJob(in.LoadId)
	if err != nil || job == nil {
		return nil, bulkLoadNotFound(in.LoadId)
	}

	page := in.Page
	if page <= 0 {
		page = 1
	}
	errorsPerPage := in.ErrorsPerPage
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
	if in.Errors && job.GetErrorLog() != "" {
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
			"loadId":     in.LoadId,
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
	if in.Details && len(job.GetFailedFeeds()) > 0 {
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

// ListLoaderJobsInput carries the parsed ListLoaderJobs wire payload. HasLimit
// distinguishes an explicitly sent limit from an omitted member.
type ListLoaderJobsInput struct {
	IncludeQueuedLoads bool
	HasLimit           bool
	Limit              int
	Region             string
}

// listLoaderJobsCore lists the load IDs of all submitted bulk load jobs,
// optionally including queued loads. An explicit limit outside the documented
// 1-100 window is rejected; an omitted limit caps the list at the documented
// default of 100.
func (s *NeptuneDataService) listLoaderJobsCore(in *ListLoaderJobsInput) (map[string]interface{}, error) {
	limit := MaxListLoaderJobsLimit
	if in.HasLimit {
		if in.Limit < MinListLoaderJobsLimit || in.Limit > MaxListLoaderJobsLimit {
			return nil, invalidParameter(fmt.Sprintf("limit must be a positive integer no greater than %d: %d", MaxListLoaderJobsLimit, in.Limit))
		}
		limit = in.Limit
	}

	store, err := s.GetStoreForRegion(in.Region)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	jobs, err := store.ListLoaderJobs()
	if err != nil {
		return nil, err
	}

	loadIds := make([]string, 0, len(jobs))
	for _, job := range jobs {
		if !in.IncludeQueuedLoads && job.GetStatus() == "LOAD_QUEUED" {
			continue
		}
		loadIds = append(loadIds, job.GetLoadId())
	}

	if len(loadIds) > limit {
		loadIds = loadIds[:limit]
	}

	return map[string]interface{}{
		"status": "200",
		"payload": map[string]interface{}{
			"loadIds": loadIds,
		},
	}, nil
}

// CancelLoaderJobInput carries the parsed CancelLoaderJob wire payload.
type CancelLoaderJobInput struct {
	LoadId string
	Region string
}

// cancelLoaderJobCore cancels a running or queued bulk load job and marks its
// status as cancelled.
func (s *NeptuneDataService) cancelLoaderJobCore(in *CancelLoaderJobInput) (map[string]interface{}, error) {
	if in.LoadId == "" {
		return nil, missingParameter("loadId")
	}

	store, err := s.GetStoreForRegion(in.Region)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	job, err := store.GetLoaderJob(in.LoadId)
	if err != nil || job == nil {
		return nil, bulkLoadNotFound(in.LoadId)
	}
	st := job.GetStatus()
	if st == "LOAD_COMPLETED" || st == "LOAD_FAILED" || st == "CANCELLED" || st == "LOAD_CANCELLED_DUE_TO_ERROR" || st == "LOAD_FAILED_BECAUSE_DEPENDENCY_NOT_SATISFIED" {
		return nil, badRequest(fmt.Sprintf("cannot cancel loader job in terminal state: %s", st))
	}
	job.Status = "CANCELLED"
	job.EndTime = timestamppb.Now()
	if err := store.UpdateLoaderJob(job); err != nil {
		logs.Warn("failed to persist loader job cancellation", logs.String("loadId", in.LoadId), logs.Err(err))
	}

	s.loaderMu.Lock()
	if ch, ok := s.loaderCancelChs[in.LoadId]; ok {
		close(ch)
		delete(s.loaderCancelChs, in.LoadId)
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

// startLoaderDispatcher launches the loader job dispatcher if it is not
// already running. The dispatcher scans for LOAD_QUEUED jobs whose
// dependencies are satisfied and starts them, respecting the maxConcurrentLoads
// limit. It exits when no more queued jobs remain.
func (s *NeptuneDataService) startLoaderDispatcher() {
	if s.dispatcherRunning.Swap(true) {
		return
	}
	go s.loaderDispatchLoop()
}

// loaderDispatchLoop repeatedly scans for dispatchable queued loader jobs
// until none remain. Uses a short poll interval to keep latency low without
// busy-waiting.
func (s *NeptuneDataService) loaderDispatchLoop() {
	defer s.dispatcherRunning.Store(false)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		dispatched, err := s.dispatchQueuedJobs()
		if err != nil {
			logs.Warn("loader dispatcher error", logs.Err(err))
		}
		if dispatched == 0 {
			jobs, err := s.countQueuedJobs()
			if err != nil || jobs == 0 {
				return
			}
		}
	}
}

// countQueuedJobs counts LOAD_QUEUED jobs across all regional stores.
func (s *NeptuneDataService) countQueuedJobs() (int, error) {
	count := 0
	s.stores.Range(func(_, val any) bool {
		store, ok := val.(*neptunestore.NeptuneStore)
		if !ok {
			return true
		}
		jobs, err := store.ListLoaderJobs()
		if err != nil {
			return true
		}
		for _, job := range jobs {
			if job.GetStatus() == "LOAD_QUEUED" {
				count++
			}
		}
		return true
	})
	return count, nil
}

// dispatchQueuedJobs scans all stores for LOAD_QUEUED jobs whose dependencies
// are satisfied and launches them, respecting the concurrency limit.
// Returns the number of jobs dispatched in this pass.
func (s *NeptuneDataService) dispatchQueuedJobs() (int, error) {
	dispatched := 0

	s.stores.Range(func(regionVal, val any) bool {
		region, rok := regionVal.(string)
		store, ok := val.(*neptunestore.NeptuneStore)
		if !ok || !rok {
			return true
		}

		jobs, err := store.ListLoaderJobs()
		if err != nil {
			return true
		}

		running := 0
		for _, job := range jobs {
			if job.GetStatus() == "LOAD_IN_PROGRESS" {
				running++
			}
		}

		for _, job := range jobs {
			if job.GetStatus() != "LOAD_QUEUED" {
				continue
			}
			if running >= maxConcurrentLoads {
				break
			}

			if !s.dependenciesSatisfied(store, job) {
				continue
			}

			clusterDB := s.GetClusterEngineForRegion(region)

			// Re-read the job from the store to detect a concurrent
			// CancelLoaderJob before transitioning to LOAD_IN_PROGRESS.
			// Without this, the dispatcher's snapshot can overwrite a
			// CANCELLED status set between ListLoaderJobs and
			// UpdateLoaderJob.
			current, err := store.GetLoaderJob(job.GetLoadId())
			if err != nil || current == nil || current.GetStatus() != "LOAD_QUEUED" {
				continue
			}

			current.Status = "LOAD_IN_PROGRESS"
			_ = store.UpdateLoaderJob(current)

			s.launchLoaderJob(region, current, clusterDB)
			running++
			dispatched++
		}
		return true
	})

	return dispatched, nil
}

// dependenciesSatisfied checks whether all dependency load IDs are in a
// LOAD_COMPLETED state. If any dependency failed, the job is marked
// LOAD_FAILED_BECAUSE_DEPENDENCY_NOT_SATISFIED.
func (s *NeptuneDataService) dependenciesSatisfied(store *neptunestore.NeptuneStore, job *pb.LoaderJob) bool {
	for _, depID := range job.GetDependencies() {
		dep, err := store.GetLoaderJob(depID)
		if err != nil || dep == nil {
			job.Status = "LOAD_FAILED_BECAUSE_DEPENDENCY_NOT_SATISFIED"
			job.EndTime = timestamppb.Now()
			job.ErrorLog = fmt.Sprintf("dependency %s not found", depID)
			_ = store.UpdateLoaderJob(job)
			return false
		}
		depStatus := dep.GetStatus()
		if depStatus == "LOAD_FAILED" || depStatus == "CANCELLED" || depStatus == "LOAD_FAILED_BECAUSE_DEPENDENCY_NOT_SATISFIED" {
			job.Status = "LOAD_FAILED_BECAUSE_DEPENDENCY_NOT_SATISFIED"
			job.EndTime = timestamppb.Now()
			job.ErrorLog = fmt.Sprintf("dependency %s is in failed state: %s", depID, depStatus)
			_ = store.UpdateLoaderJob(job)
			return false
		}
		if depStatus != "LOAD_COMPLETED" {
			return false
		}
	}
	return true
}
