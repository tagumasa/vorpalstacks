package cloudwatchlogs

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/utils/aws/arn"
)

// The limit parameter of DescribeImportTaskBatches: "Valid Range:
// Minimum value of 1. Maximum value of 50. Default: 10".
const (
	defaultImportTaskBatchesLimit = 10
	maxImportTaskBatchesLimit     = 50
)

// --- Export Task Core ---

// CreateExportTaskInput holds parameters for CreateExportTask.
type CreateExportTaskInput struct {
	TaskName            string
	LogGroupName        string
	Destination         string
	From                int
	To                  int
	LogStreamNamePrefix string
	DestinationPrefix   string
	Region              string
}

// createExportTaskCore validates input and creates an export task.
func (s *LogsService) createExportTaskCore(store *logsstore.Store, input *CreateExportTaskInput) (string, error) {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return "", err
	}
	if err := validateExportDestinationBucket(input.Destination); err != nil {
		return "", err
	}
	// taskName: "Length Constraints: Minimum length of 1. Maximum length
	// of 512." (CreateExportTask API reference) — the absent member is the
	// empty string the wire cannot distinguish from a present empty one.
	if utf8.RuneCountInString(input.TaskName) > logsstore.MaxExportTaskNameLength {
		return "", NewLogsError("InvalidParameterException",
			fmt.Sprintf("taskName must be between 1 and %d characters", logsstore.MaxExportTaskNameLength), 400)
	}
	// from/to target Timestamp, whose range trait floors at zero — "Valid
	// Range: Minimum value of 0" — so from == 0 (the epoch start) is a
	// valid window start.
	if input.From < 0 || input.To < 0 || input.From >= input.To {
		return "", NewLogsError("InvalidParameterException",
			"from and to must be non-negative timestamps with from < to", 400)
	}
	if err := validateLogStreamNamePrefix(input.LogStreamNamePrefix); err != nil {
		return "", err
	}
	if err := validateDestinationPrefix(input.DestinationPrefix); err != nil {
		return "", err
	}

	lg, err := store.GetLogGroup(input.LogGroupName)
	if err != nil {
		return "", mapStoreError(err)
	}
	// "You must specify a time that is not earlier than when this log
	// group was created." (to, CreateExportTask API reference).
	if int64(input.To) < lg.CreatedAt.UnixMilli() {
		return "", NewLogsError("InvalidParameterException",
			"to must not be earlier than when this log group was created", 400)
	}

	// The documented quota: "Each account can only have one active
	// (RUNNING or PENDING) export task at a time" (CreateExportTask
	// documentation, LimitExceededException declared on the operation).
	// The quotas registry fixes the scope: "Active export task — Each
	// supported Region: 1" — the Region's active slot, not the account's
	// across Regions (an export running in another Region leaves this
	// Region's slot open). The census and the RUNNING record's write run
	// as one critical section so two concurrent creations cannot both
	// observe the quota open.
	s.taskAdmissionMu.Lock()
	// A RUNNING task past the documented 24-hour deadline expires before
	// the census, so a timed-out task cannot hold the Region's single
	// active slot.
	s.expireTimedOutExportTasks()
	if hasActiveExportTask(store, "") {
		s.taskAdmissionMu.Unlock()
		return "", ErrLimitExceeded
	}

	taskId := fmt.Sprintf("export-%d", time.Now().UnixNano())

	task := &logsstore.ExportTask{
		TaskId:              taskId,
		TaskName:            input.TaskName,
		LogGroupName:        input.LogGroupName,
		LogStreamNamePrefix: input.LogStreamNamePrefix,
		From:                int64(input.From),
		To:                  int64(input.To),
		Destination:         input.Destination,
		DestinationPrefix:   input.DestinationPrefix,
		Status:              logsstore.ExportStatusRunning,
		CreationTime:        time.Now().UTC().UnixMilli(),
	}

	if err := store.PutExportTask(task); err != nil {
		s.taskAdmissionMu.Unlock()
		return "", mapStoreError(err)
	}
	s.taskAdmissionMu.Unlock()

	s.spawnTask(func(ctx context.Context) {
		s.executeExportTask(ctx, input.Region, input.LogGroupName, input.LogStreamNamePrefix,
			int64(input.From), int64(input.To), input.Destination, input.DestinationPrefix, taskId)
	})

	return taskId, nil
}

// hasActiveExportTask reports whether the Region's store holds a PENDING
// or RUNNING export task other than the excluded one. A listing failure
// reads as no active task — the census is best-effort, as it was when it
// walked the region set.
func hasActiveExportTask(store *logsstore.Store, excludeTaskId string) bool {
	tasks, err := store.ListExportTasks("")
	if err != nil {
		return false
	}
	for _, t := range tasks {
		if t.TaskId == excludeTaskId {
			continue
		}
		if t.Status == logsstore.ExportStatusRunning || t.Status == logsstore.ExportStatusPending {
			return true
		}
	}
	return false
}

// DescribeExportTasksInput holds parameters for DescribeExportTasks.
type DescribeExportTasksInput struct {
	TaskId     string
	StatusCode string
	NextToken  string
	Limit      int32
}

// describeExportTasksCore validates input and lists export tasks.
func (s *LogsService) describeExportTasksCore(store *logsstore.Store, input *DescribeExportTasksInput) ([]*logsstore.ExportTask, string, error) {
	// statusCode is the ExportTaskStatusCode enumeration; a value outside
	// the vocabulary rejects.
	if input.StatusCode != "" && !validExportStatusCode[input.StatusCode] {
		return nil, "", NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid statusCode: %s. Valid values: CANCELLED, COMPLETED, FAILED, PENDING, PENDING_CANCEL, RUNNING", input.StatusCode), 400)
	}

	// A RUNNING task past the documented 24-hour deadline expires before
	// the listing resolves, so the read reports the timed-out state.
	s.expireTimedOutExportTasks()

	if input.TaskId != "" {
		task, err := store.GetExportTask(input.TaskId)
		if err != nil {
			// "Specifying a task ID filters the results to one or zero
			// export tasks" — an unknown id is the zero side of that
			// filter (an empty list, not an error: the operation declares
			// no ResourceNotFoundException), and the statusCode filter
			// keeps applying to the resolved task: a status mismatch is
			// the same zero side.
			if err == logsstore.ErrResourceNotFound {
				return nil, "", nil
			}
			return nil, "", mapStoreError(err)
		}
		if input.StatusCode != "" && task.Status != input.StatusCode {
			return nil, "", nil
		}
		return []*logsstore.ExportTask{task}, "", nil
	}

	limit, err := validateListLimit(input.Limit, logsstore.DefaultDescribeLimit, logsstore.MaxDescribeLimit)
	if err != nil {
		return nil, "", err
	}

	allTasks, err := store.ListExportTasks(input.StatusCode)
	if err != nil {
		return nil, "", mapStoreError(err)
	}

	// The page marker rides the scoped listing vocabulary: the token
	// carries the request identity and its documented expiry rather than
	// a bare task id.
	scope := listingScope("exporttasks", input.StatusCode)
	result, err := paginateScopedListing(scope, input.NextToken, allTasks, int(limit), func(t *logsstore.ExportTask) string {
		return t.TaskId
	})
	if err != nil {
		return nil, "", err
	}
	return result.Items, result.NextMarker, nil
}

// cancelExportTaskCore validates input and cancels an export task. The
// terminal-state gate and the status write run as one atomic
// read-modify-write, so a worker completion landing mid-cancel cannot
// be overwritten and a completed task cannot be cancelled. The
// canceller writes PENDING_CANCEL, the live transitional state the
// status vocabulary documents: the running worker honours the request
// at its next checkpoint and writes the terminal CANCELLED, so a
// client never observes a terminal status while the worker may still
// be mid-read. A repeat cancel of a PENDING_CANCEL task is idempotent.
func (s *LogsService) cancelExportTaskCore(store *logsstore.Store, taskId string) error {
	if taskId == "" {
		return errRequiredMember("taskId")
	}

	err := store.MutateExportTask(taskId, func(task *logsstore.ExportTask) error {
		if isTerminalExportStatus(task.Status) {
			return NewLogsError("InvalidOperationException",
				fmt.Sprintf("Cannot cancel export task in %s state", task.Status), 400)
		}
		if task.Status == logsstore.ExportStatusPendingCancel {
			return nil
		}
		task.Status = logsstore.ExportStatusPendingCancel
		task.StatusMessage = "Cancelled by user"
		return nil
	})
	return mapStoreError(err)
}

// --- Import Task Core ---

// CreateImportTaskInput holds parameters for CreateImportTask.
type CreateImportTaskInput struct {
	ImportSourceArn string
	ImportRoleArn   string
	ImportFilter    map[string]interface{}
	Region          string
	AccountID       string
}

// createImportTaskCore validates input and creates an import task.
func (s *LogsService) createImportTaskCore(store *logsstore.Store, input *CreateImportTaskInput) (*logsstore.ImportTask, error) {
	if input.ImportSourceArn == "" {
		return nil, errRequiredMember("importSourceArn")
	}
	// importRoleArn is a required member: "The ARN of the IAM role that
	// grants CloudWatch Logs permission to import" (member
	// documentation). The platform performs the S3 reads through its own
	// invoker rather than assuming the role, so the member is validated
	// for form and stored for the response's fidelity.
	if input.ImportRoleArn == "" {
		return nil, NewLogsError("InvalidParameterException",
			"The importRoleArn member is required", 400)
	}
	if err := validateIAMRoleArn(input.ImportRoleArn); err != nil {
		return nil, err
	}

	// The documented source form is the CloudTrail Lake event data store
	// ("The only supported source is a CloudTrail Lake event data store" —
	// the CreateImportTask example is an eventdatastore ARN); the S3
	// object form is the platform's additional source, imported through
	// the S3 invoker. Both address shapes parse here to the execution
	// plan their executor reads.
	parsedArn, err := arn.ParseARN(input.ImportSourceArn)
	if err != nil {
		return nil, NewLogsError("InvalidParameterException",
			"importSourceArn must be a CloudTrail Lake event data store ARN or an S3 ARN", 400)
	}
	var bucket, s3Key, edsID string
	switch parsedArn.Service {
	case "s3":
		bucket = parsedArn.AccountID
		s3Key = parsedArn.Resource
		if idx := strings.Index(s3Key, "/"); idx >= 0 {
			bucket = s3Key[:idx]
			s3Key = s3Key[idx+1:]
		} else {
			bucket = s3Key
			s3Key = ""
		}
		if bucket == "" {
			return nil, NewLogsError("InvalidParameterException",
				"importSourceArn must contain an S3 bucket", 400)
		}
	case "cloudtrail":
		resource := parsedArn.Resource
		if idx := strings.Index(resource, "/"); idx >= 0 {
			edsID = resource[idx+1:]
		}
		if !strings.HasPrefix(resource, "eventdatastore/") || edsID == "" {
			return nil, NewLogsError("InvalidParameterException",
				"importSourceArn must address a CloudTrail Lake event data store (eventdatastore/<id>)", 400)
		}
	default:
		return nil, NewLogsError("InvalidParameterException",
			"importSourceArn must be a CloudTrail Lake event data store ARN or an S3 ARN", 400)
	}

	if input.ImportFilter != nil {
		startMs, hasStart, err := importFilterTimestamp(input.ImportFilter, "startEventTime")
		if err != nil {
			return nil, err
		}
		endMs, hasEnd, err := importFilterTimestamp(input.ImportFilter, "endEventTime")
		if err != nil {
			return nil, err
		}
		if hasStart && hasEnd && startMs > endMs {
			return nil, NewLogsError("InvalidParameterException",
				"startEventTime must be less than or equal to endEventTime", 400)
		}
	}

	// The documented account-level concurrency quota: "There can be no
	// more than 3 active imports per account at a given time"
	// (CreateImportTask operation documentation). Active imports are
	// counted across every configured region — the account's imports,
	// not one region's — because this quota's only source is the
	// operation documentation's "per account": the quotas registry
	// carries no row for it, unlike the export quota whose registry row
	// is "Each supported Region: 1" and whose census is the Region's
	// store alone. The rejection is ThrottlingException, the one
	// quota-shaped error the operation declares ("The request was
	// throttled because of quota limits"); LimitExceededException is the
	// quota error of other operations' error lists and is undeclared
	// here, so emitting it would be an error identity the operation does
	// not carry.
	// The admission window runs as one critical section: the census and
	// the RUNNING record's write cannot interleave with another
	// creation's, so the quota cannot be double-admitted.
	s.taskAdmissionMu.Lock()
	if s.countActiveImportTasks() >= logsstore.MaxActiveImportTasks {
		s.taskAdmissionMu.Unlock()
		return nil, NewLogsError("ThrottlingException",
			fmt.Sprintf("There can be no more than %d active imports per account at a given time", logsstore.MaxActiveImportTasks), 400)
	}

	// The managed group names the source under the documented CloudTrail
	// Lake import prefix: the event data store's id or the S3 bucket.
	sourceName := bucket
	if edsID != "" {
		sourceName = edsID
	}
	managedLogGroupName := deriveManagedLogGroupName(sourceName)
	// The managed group the task creates is a service-named group like
	// any other: the derived name must satisfy the log-group name rules
	// before anything creates it.
	if err := validateLogGroupName(managedLogGroupName); err != nil {
		s.taskAdmissionMu.Unlock()
		return nil, err
	}
	importDestinationArn := arn.NewARNBuilder(input.AccountID, input.Region).CloudWatch().LogGroup(managedLogGroupName)

	if _, err := store.GetLogGroup(managedLogGroupName); err != nil {
		// The managed group is a log group like any other and takes one
		// of the Region's slots — the admission the CreateLogGroup core
		// runs ("You can create up to 1,000,000 log groups per Region
		// per account"), rejected here with this operation's declared
		// quota error identity.
		groupCount, err := store.CountLogGroups()
		if err != nil {
			s.taskAdmissionMu.Unlock()
			return nil, mapStoreError(err)
		}
		if groupCount >= logsstore.MaxLogGroupsPerRegion {
			s.taskAdmissionMu.Unlock()
			return nil, NewLogsError("ThrottlingException",
				fmt.Sprintf("You can create up to %d log groups per Region per account", logsstore.MaxLogGroupsPerRegion), 400)
		}
		lg := logsstore.NewLogGroup(managedLogGroupName, input.Region, input.AccountID)
		if err := store.CreateLogGroup(lg); err != nil {
			if existing, gErr := store.GetLogGroup(managedLogGroupName); gErr != nil || existing == nil {
				s.taskAdmissionMu.Unlock()
				return nil, mapStoreError(err)
			}
		}
	}

	importId := fmt.Sprintf("import-%d", time.Now().UnixNano())

	task := &logsstore.ImportTask{
		ImportId:             importId,
		ImportSourceArn:      input.ImportSourceArn,
		ImportRoleArn:        input.ImportRoleArn,
		LogGroupName:         managedLogGroupName,
		ImportStatus:         logsstore.ImportStatusInProgress,
		ImportDestinationArn: importDestinationArn,
		ImportFilter:         input.ImportFilter,
		CreationTime:         time.Now().UTC().UnixMilli(),
	}

	if err := store.PutImportTask(task); err != nil {
		s.taskAdmissionMu.Unlock()
		return nil, mapStoreError(err)
	}
	s.taskAdmissionMu.Unlock()

	s.spawnTask(func(ctx context.Context) {
		if edsID != "" {
			s.executeImportTaskFromEDS(ctx, input.Region, importId, input.ImportSourceArn, managedLogGroupName)
			return
		}
		s.executeImportTask(ctx, input.Region, importId, bucket, s3Key, managedLogGroupName)
	})

	return task, nil
}

// importFilterTimestamp reads one importFilter window member: the value
// targets Timestamp, whose range trait floors at zero, so a non-numeric or
// negative value rejects instead of silently becoming the unbounded
// window the executor would treat it as. The third return reports whether
// the member was present.
func importFilterTimestamp(filter map[string]interface{}, member string) (int64, bool, error) {
	raw, present := filter[member]
	if !present {
		return 0, false, nil
	}
	value, ok := toInt64(raw)
	if !ok || value < 0 {
		return 0, false, NewLogsError("InvalidParameterException",
			fmt.Sprintf("%s must be a non-negative number of milliseconds", member), 400)
	}
	return value, true, nil
}

// countActiveImportTasks counts the account's IN_PROGRESS import tasks
// across every configured region — the account-wide scope this quota's
// only source states ("per account"; the quotas registry carries no row
// for it). The region set comes from the storage manager, not the store
// cache.
func (s *LogsService) countActiveImportTasks() int {
	count := 0
	if s.storageManager == nil {
		return count
	}
	for _, region := range s.storageManager.GetActiveRegions() {
		store, err := s.getLogsStoreByRegion(region)
		if err != nil {
			continue
		}
		tasks, err := store.ListImportTasks("", "", "")
		if err != nil {
			continue
		}
		for _, t := range tasks {
			if t.ImportStatus == logsstore.ImportStatusInProgress {
				count++
			}
		}
	}
	return count
}

// describeImportTasksCore validates input and lists import tasks.
func (s *LogsService) describeImportTasksCore(store *logsstore.Store, importId, status, sourceArn, nextToken string, limit int32) ([]*logsstore.ImportTask, string, error) {
	// importStatus is the ImportStatus enumeration; a value outside the
	// vocabulary rejects.
	if status != "" && !validImportStatus[status] {
		return nil, "", NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid importStatus: %s. Valid values: IN_PROGRESS, CANCELLED, COMPLETED, FAILED", status), 400)
	}
	// The optional importId filter still carries the ImportId traits: a
	// malformed id is a parameter error, not an empty result.
	if importId != "" {
		if err := validateImportId(importId); err != nil {
			return nil, "", err
		}
	}

	limit, err := validateListLimit(limit, logsstore.DefaultDescribeLimit, logsstore.MaxDescribeLimit)
	if err != nil {
		return nil, "", err
	}

	allTasks, err := store.ListImportTasks(importId, status, sourceArn)
	if err != nil {
		return nil, "", mapStoreError(err)
	}

	// The page marker rides the scoped listing vocabulary: the token
	// carries the request identity and its documented expiry rather than
	// a bare import id.
	scope := listingScope("importtasks", importId, status, sourceArn)
	result, err := paginateScopedListing(scope, nextToken, allTasks, int(limit), func(t *logsstore.ImportTask) string {
		return t.ImportId
	})
	if err != nil {
		return nil, "", err
	}
	return result.Items, result.NextMarker, nil
}

// DescribeImportTaskBatchesInput holds parameters for
// DescribeImportTaskBatches.
type DescribeImportTaskBatchesInput struct {
	ImportId          string
	BatchImportStatus []string
	NextToken         string
	Limit             int32
}

// DescribeImportTaskBatchesResult holds the batches page.
type DescribeImportTaskBatchesResult struct {
	Task          *logsstore.ImportTask
	ImportBatches []*logsstore.ImportBatch
	NextToken     string
}

// describeImportTaskBatchesCore validates input, resolves the task and
// serves one page of its batch records, optionally filtered by status.
// Batches are grouped by event time (the granularity the operation
// documentation names), ordered by their bucket boundary; the page walks
// positionally, which is exact because a task's batch list is immutable
// once the task reaches a terminal status.
func (s *LogsService) describeImportTaskBatchesCore(store *logsstore.Store, input *DescribeImportTaskBatchesInput) (*DescribeImportTaskBatchesResult, error) {
	if err := validateImportId(input.ImportId); err != nil {
		return nil, err
	}
	task, err := store.GetImportTask(input.ImportId)
	if err != nil {
		return nil, mapStoreError(err)
	}

	// limit: "Valid Range: Minimum value of 1. Maximum value of 50.
	// Default: 10".
	limit, err := validateListLimit(input.Limit, defaultImportTaskBatchesLimit, maxImportTaskBatchesLimit)
	if err != nil {
		return nil, err
	}

	// batchImportStatus is the BatchImportStatus enumeration; a value
	// outside the vocabulary rejects.
	statusSet := make(map[string]bool, len(input.BatchImportStatus))
	for _, st := range input.BatchImportStatus {
		if !validImportStatus[st] {
			return nil, NewLogsError("InvalidParameterException",
				fmt.Sprintf("Invalid batchImportStatus: %s. Valid values: IN_PROGRESS, CANCELLED, COMPLETED, FAILED", st), 400)
		}
		statusSet[st] = true
	}
	var filtered []*logsstore.ImportBatch
	for _, b := range task.ImportBatches {
		if len(statusSet) > 0 && !statusSet[b.Status] {
			continue
		}
		filtered = append(filtered, b)
	}

	// The page marker rides the scoped listing vocabulary: the token
	// carries the request identity and its documented expiry rather than
	// a bare positional index.
	statusKey := ""
	if len(input.BatchImportStatus) > 0 {
		sorted := make([]string, len(input.BatchImportStatus))
		copy(sorted, input.BatchImportStatus)
		sort.Strings(sorted)
		statusKey = strings.Join(sorted, "\x1f")
	}
	scope := listingScope("importtaskbatches", input.ImportId, statusKey)
	page, err := paginateScopedListingByPosition(scope, input.NextToken, filtered, int(limit))
	if err != nil {
		return nil, err
	}
	return &DescribeImportTaskBatchesResult{
		Task:          task,
		ImportBatches: page.Items,
		NextToken:     page.NextMarker,
	}, nil
}

// cancelImportTaskCore validates input, cancels an import task, and re-reads
// the persisted record for the response. The terminal-state gate and the
// CANCELLED write run as one atomic read-modify-write (the export twin's
// pattern), so a worker completion landing mid-cancel cannot be overwritten
// by a stale snapshot. A read failure degrades to a nil task, to which the
// caller responds with the minimal member set.
func (s *LogsService) cancelImportTaskCore(store *logsstore.Store, importId string) (*logsstore.ImportTask, error) {
	if err := validateImportId(importId); err != nil {
		return nil, err
	}

	err := store.MutateImportTask(importId, func(task *logsstore.ImportTask) error {
		if task.ImportStatus == logsstore.ImportStatusCompleted || task.ImportStatus == logsstore.ImportStatusFailed || task.ImportStatus == logsstore.ImportStatusCancelled {
			return NewLogsError("InvalidOperationException",
				fmt.Sprintf("Cannot cancel import task in %s state", task.ImportStatus), 400)
		}
		task.ImportStatus = logsstore.ImportStatusCancelled
		return nil
	})
	if err != nil {
		return nil, mapStoreError(err)
	}
	saved, _ := store.GetImportTask(importId)
	return saved, nil
}
