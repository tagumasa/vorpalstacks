package iam

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const slRoleDeletionTaskBucketName = "iam_sl_role_deletion_tasks"

// SLRoleDeletionTask represents an async service-linked role deletion
// task.  The task is created with status IN_PROGRESS, updated to
// SUCCEEDED or FAILED when the background cleanup goroutine completes.
type SLRoleDeletionTask struct {
	TaskID         string    `json:"taskId"`
	RoleName       string    `json:"roleName"`
	Status         string    `json:"status"`
	DeletionFailed bool      `json:"deletionFailed"`
	ErrorReason    string    `json:"errorReason,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

// SLRoleDeletionTaskStore manages service-linked role deletion tasks in
// persistent storage so that task status survives server restarts.
type SLRoleDeletionTaskStore struct {
	*common.BaseStore
}

// NewSLRoleDeletionTaskStore creates a new store for SL role deletion tasks.
func NewSLRoleDeletionTaskStore(store storage.BasicStorage) *SLRoleDeletionTaskStore {
	return &SLRoleDeletionTaskStore{
		BaseStore: common.NewBaseStore(store.Bucket(slRoleDeletionTaskBucketName), "iam"),
	}
}

// Put stores a deletion task.
func (s *SLRoleDeletionTaskStore) Put(task *SLRoleDeletionTask) error {
	return s.BaseStore.Put(task.TaskID, task)
}

// Get retrieves a deletion task by its ID.
func (s *SLRoleDeletionTaskStore) Get(taskID string) (*SLRoleDeletionTask, error) {
	var task SLRoleDeletionTask
	if err := s.BaseStore.Get(taskID, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

// List returns all deletion tasks.
func (s *SLRoleDeletionTaskStore) List() ([]*SLRoleDeletionTask, error) {
	var tasks []*SLRoleDeletionTask
	err := s.ForEach(func(k string, v []byte) error {
		var task SLRoleDeletionTask
		if err := json.Unmarshal(v, &task); err != nil {
			return err
		}
		tasks = append(tasks, &task)
		return nil
	})
	return tasks, err
}

// Delete removes a deletion task.
func (s *SLRoleDeletionTaskStore) Delete(taskID string) error {
	return s.BaseStore.Delete(taskID)
}
