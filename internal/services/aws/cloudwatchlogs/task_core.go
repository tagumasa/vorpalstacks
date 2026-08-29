package cloudwatchlogs

import (
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/pagination"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/utils/aws/arn"
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
	if input.From <= 0 || input.To <= 0 || input.From >= input.To {
		return "", NewLogsError("InvalidParameterException",
			"from and to must be positive timestamps with from < to", 400)
	}
	if err := validateDestinationPrefix(input.DestinationPrefix); err != nil {
		return "", err
	}

	if _, err := store.GetLogGroup(input.LogGroupName); err != nil {
		return "", mapStoreError(err)
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
		Status:              "RUNNING",
		CreationTime:        time.Now().UTC().UnixMilli(),
	}

	if err := store.PutExportTask(task); err != nil {
		return "", mapStoreError(err)
	}

	go s.executeExportTask(input.Region, input.LogGroupName, input.LogStreamNamePrefix,
		int64(input.From), int64(input.To), input.Destination, input.DestinationPrefix, taskId)

	return taskId, nil
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
	if input.TaskId != "" {
		task, err := store.GetExportTask(input.TaskId)
		if err != nil {
			return nil, "", mapStoreError(err)
		}
		return []*logsstore.ExportTask{task}, "", nil
	}

	limit := input.Limit
	if limit <= 0 {
		limit = 50
	}

	allTasks, err := store.ListExportTasks(input.StatusCode)
	if err != nil {
		return nil, "", mapStoreError(err)
	}

	result := pagination.PaginateSlice(allTasks, input.NextToken, int(limit), func(t *logsstore.ExportTask) string {
		return t.TaskId
	})
	return result.Items, result.NextMarker, nil
}

// cancelExportTaskCore validates input and cancels an export task.
func (s *LogsService) cancelExportTaskCore(store *logsstore.Store, taskId string) error {
	if taskId == "" {
		return ErrMissingParameter
	}

	task, err := store.GetExportTask(taskId)
	if err != nil {
		return mapStoreError(err)
	}

	if task.Status == "COMPLETED" || task.Status == "FAILED" || task.Status == "CANCELLED" {
		return NewLogsError("InvalidOperationException",
			fmt.Sprintf("Cannot cancel export task in %s state", task.Status), 400)
	}

	task.Status = "CANCELLED"
	task.StatusMessage = "Cancelled by user"
	if err := store.PutExportTask(task); err != nil {
		return mapStoreError(err)
	}
	return nil
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
		return nil, ErrMissingParameter
	}

	parsedArn, err := arn.ParseARN(input.ImportSourceArn)
	if err != nil || parsedArn.Service != "s3" {
		return nil, NewLogsError("InvalidParameterException",
			"importSourceArn must be a valid S3 ARN", 400)
	}

	bucket := parsedArn.AccountID
	s3Key := parsedArn.Resource
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

	if input.ImportFilter != nil {
		startTime, hasStart := input.ImportFilter["startEventTime"]
		endTime, hasEnd := input.ImportFilter["endEventTime"]
		if hasStart && hasEnd {
			startMs, okS := toInt64(startTime)
			endMs, okE := toInt64(endTime)
			if okS && okE && startMs > endMs {
				return nil, NewLogsError("InvalidParameterException",
					"startEventTime must be less than or equal to endEventTime", 400)
			}
		}
	}

	managedLogGroupName := deriveManagedLogGroupName(input.ImportSourceArn)
	importDestinationArn := fmt.Sprintf("arn:aws:logs:%s:%s:log-group:%s", input.Region, input.AccountID, managedLogGroupName)

	if _, err := store.GetLogGroup(managedLogGroupName); err != nil {
		lg := logsstore.NewLogGroup(managedLogGroupName, input.Region, input.AccountID)
		if err := store.CreateLogGroup(lg); err != nil {
			if existing, gErr := store.GetLogGroup(managedLogGroupName); gErr != nil || existing == nil {
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
		ImportStatus:         "IN_PROGRESS",
		ImportDestinationArn: importDestinationArn,
		ImportFilter:         input.ImportFilter,
		CreationTime:         time.Now().UTC().UnixMilli(),
	}

	if err := store.PutImportTask(task); err != nil {
		return nil, mapStoreError(err)
	}

	go s.executeImportTask(input.Region, importId, bucket, s3Key, managedLogGroupName)

	return task, nil
}

// describeImportTasksCore validates input and lists import tasks.
func (s *LogsService) describeImportTasksCore(store *logsstore.Store, importId, status, sourceArn, nextToken string, limit int32) ([]*logsstore.ImportTask, string, error) {
	if limit <= 0 {
		limit = 50
	}

	allTasks, err := store.ListImportTasks(importId, status, sourceArn)
	if err != nil {
		return nil, "", mapStoreError(err)
	}

	result := pagination.PaginateSlice(allTasks, nextToken, int(limit), func(t *logsstore.ImportTask) string {
		return t.ImportId
	})
	return result.Items, result.NextMarker, nil
}

// describeImportTaskBatchesCore validates input and resolves the import task
// whose batches are being described. batchImportStatus is accepted at the
// handler layer but batch-level status tracking is not implemented; callers
// render an empty importBatches list.
func (s *LogsService) describeImportTaskBatchesCore(store *logsstore.Store, importId string) (*logsstore.ImportTask, error) {
	if importId == "" {
		return nil, ErrMissingParameter
	}
	task, err := store.GetImportTask(importId)
	if err != nil {
		return nil, mapStoreError(err)
	}
	return task, nil
}

// cancelImportTaskCore validates input, cancels an import task, and re-reads
// the persisted record for the response. A read failure degrades to a nil
// task, to which the caller responds with the minimal member set.
func (s *LogsService) cancelImportTaskCore(store *logsstore.Store, importId string) (*logsstore.ImportTask, error) {
	if importId == "" {
		return nil, ErrMissingParameter
	}

	task, err := store.GetImportTask(importId)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if task.ImportStatus == "COMPLETED" || task.ImportStatus == "FAILED" || task.ImportStatus == "CANCELLED" {
		return nil, NewLogsError("InvalidOperationException",
			fmt.Sprintf("Cannot cancel import task in %s state", task.ImportStatus), 400)
	}

	task.ImportStatus = "CANCELLED"
	task.LastUpdatedTime = time.Now().UTC().UnixMilli()
	if err := store.PutImportTask(task); err != nil {
		return nil, mapStoreError(err)
	}
	saved, _ := store.GetImportTask(importId)
	return saved, nil
}
