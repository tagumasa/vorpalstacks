package cloudwatchlogs

import (
	"sync"
	"time"
)

// taskRecordMu serialises read-modify-write cycles on task records: the
// worker's status writes, the canceller's transition and the quota
// admission's terminal check all decide on the persisted record under
// this mutex, so a status decided by one side cannot be overwritten by
// a stale snapshot from another.
var taskRecordMu sync.Mutex

// --- Export Task ---

// MutateExportTask loads the export task, applies fn, and persists the
// result as one atomic read-modify-write; an fn error aborts without
// persisting. Every export-task status decision (the worker's
// transitions, the canceller's gate) runs through here.
func (s *Store) MutateExportTask(taskId string, fn func(*ExportTask) error) error {
	taskRecordMu.Lock()
	defer taskRecordMu.Unlock()
	task, err := s.GetExportTask(taskId)
	if err != nil {
		return err
	}
	if err := fn(task); err != nil {
		return err
	}
	return s.putExportTaskLocked(task)
}

// PutExportTask writes an export task record under the task family mutex:
// an external full-record write (a test seed, a reconciliation pass)
// cannot interleave with a mutate cycle's read-modify-write.
func (s *Store) PutExportTask(task *ExportTask) error {
	taskRecordMu.Lock()
	defer taskRecordMu.Unlock()
	return s.putExportTaskLocked(task)
}

func (s *Store) putExportTaskLocked(task *ExportTask) error {
	return s.Put(s.exportTaskKey(task.TaskId), task)
}

func (s *Store) GetExportTask(taskId string) (*ExportTask, error) {
	return getJSONRecord[ExportTask](s, s.exportTaskKey(taskId), ErrResourceNotFound)
}

func (s *Store) DeleteExportTask(taskId string) error {
	taskRecordMu.Lock()
	defer taskRecordMu.Unlock()
	return s.Delete(s.exportTaskKey(taskId))
}

func (s *Store) ListExportTasks(statusCode string) ([]*ExportTask, error) {
	return listJSONRecords[ExportTask](s, keyPrefixExportTask, "export task record", func(t *ExportTask) bool {
		return statusCode == "" || t.Status == statusCode
	})
}

// --- Import Task ---

// MutateImportTask is the import plane's equivalent of MutateExportTask.
func (s *Store) MutateImportTask(importId string, fn func(*ImportTask) error) error {
	taskRecordMu.Lock()
	defer taskRecordMu.Unlock()
	task, err := s.GetImportTask(importId)
	if err != nil {
		return err
	}
	if err := fn(task); err != nil {
		return err
	}
	return s.putImportTaskLocked(task)
}

// PutImportTask writes an import task record under the task family mutex,
// stamping the update time inside the critical section.
func (s *Store) PutImportTask(task *ImportTask) error {
	taskRecordMu.Lock()
	defer taskRecordMu.Unlock()
	return s.putImportTaskLocked(task)
}

func (s *Store) putImportTaskLocked(task *ImportTask) error {
	task.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.importTaskKey(task.ImportId), task)
}

func (s *Store) GetImportTask(importId string) (*ImportTask, error) {
	return getJSONRecord[ImportTask](s, s.importTaskKey(importId), ErrResourceNotFound)
}

func (s *Store) ListImportTasks(importId, status, sourceArn string) ([]*ImportTask, error) {
	return listJSONRecords[ImportTask](s, keyPrefixImportTask, "import task record", func(t *ImportTask) bool {
		if importId != "" && t.ImportId != importId {
			return false
		}
		if status != "" && t.ImportStatus != status {
			return false
		}
		return sourceArn == "" || t.ImportSourceArn == sourceArn
	})
}

// ExportTask represents a CloudWatch Logs export-to-S3 task.
type ExportTask struct {
	TaskId              string                 `json:"taskId"`
	TaskName            string                 `json:"taskName"`
	LogGroupName        string                 `json:"logGroupName"`
	LogStreamNamePrefix string                 `json:"logStreamNamePrefix,omitempty"`
	From                int64                  `json:"from"`
	To                  int64                  `json:"to"`
	Destination         string                 `json:"destination"`
	DestinationPrefix   string                 `json:"destinationPrefix,omitempty"`
	Status              string                 `json:"status"`
	StatusMessage       string                 `json:"statusMessage,omitempty"`
	ExecutionInfo       map[string]interface{} `json:"executionInfo,omitempty"`
	CreationTime        int64                  `json:"creationTime"`
}

// ImportTask represents a CloudWatch Logs import-from-S3 task.
type ImportTask struct {
	ImportId             string                 `json:"importId"`
	ImportSourceArn      string                 `json:"importSourceArn"`
	ImportRoleArn        string                 `json:"importRoleArn,omitempty"`
	LogGroupName         string                 `json:"logGroupName"`
	ImportStatus         string                 `json:"importStatus"`
	ImportDestinationArn string                 `json:"importDestinationArn,omitempty"`
	ImportStatistics     map[string]interface{} `json:"importStatistics,omitempty"`
	ImportFilter         map[string]interface{} `json:"importFilter,omitempty"`
	ErrorMessage         string                 `json:"errorMessage,omitempty"`
	// ImportBatches holds the task's per-batch records (events grouped
	// by event time, the granularity the DescribeImportTaskBatches
	// documentation names) and is append-only until the task reaches a
	// terminal status, immutable after.
	ImportBatches   []*ImportBatch `json:"importBatches,omitempty"`
	CreationTime    int64          `json:"creationTime"`
	LastUpdatedTime int64          `json:"lastUpdatedTime"`
}

// ImportBatch is one import batch: the subset of a task's imported
// events that share an event-time bucket, with its own status and error
// message (the ImportBatch wire shape).
type ImportBatch struct {
	BatchId      string `json:"batchId"`
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	// HourStartMs is the bucket boundary the batch groups by; it stays
	// off the wire and orders the batch listing.
	HourStartMs int64 `json:"hourStartMs,omitempty"`
}

// Export task record statuses (the internal vocabulary of the export
// task records; PENDING and PENDING_CANCEL are the documented statuses
// the platform's synchronous worker start cannot reach).
const (
	ExportStatusRunning       = "RUNNING"
	ExportStatusPending       = "PENDING"
	ExportStatusCompleted     = "COMPLETED"
	ExportStatusFailed        = "FAILED"
	ExportStatusCancelled     = "CANCELLED"
	ExportStatusPendingCancel = "PENDING_CANCEL"
)

// Import task record statuses.
const (
	ImportStatusInProgress = "IN_PROGRESS"
	ImportStatusCompleted  = "COMPLETED"
	ImportStatusFailed     = "FAILED"
	ImportStatusCancelled  = "CANCELLED"
)
