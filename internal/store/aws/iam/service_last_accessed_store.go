package iam

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const serviceLastAccessedBucketName = "iam_service_last_accessed_jobs"

// ServiceLastAccessedJob represents a report generation job for IAM Service Last Accessed Details.
type ServiceLastAccessedJob struct {
	JobID                      string                      `json:"jobId"`
	Arn                        string                      `json:"arn"`
	JobType                    string                      `json:"jobType"`
	JobStatus                  string                      `json:"jobStatus"`
	JobCreationTime            time.Time                   `json:"jobCreationTime"`
	JobCompletionTime          *time.Time                  `json:"jobCompletionTime,omitempty"`
	Error                      string                      `json:"error,omitempty"`
	Granularity                string                      `json:"granularity,omitempty"`
	TrackedActionsLastAccessed []TrackedActionLastAccessed `json:"trackedActionsLastAccessed,omitempty"`
	ServicesLastAccessed       []ServiceLastAccessed       `json:"servicesLastAccessed,omitempty"`
}

// TrackedActionLastAccessed contains details about the last time a specific IAM action was accessed.
type TrackedActionLastAccessed struct {
	ActionName         string     `json:"actionName"`
	LastAccessedDate   *time.Time `json:"lastAccessedDate,omitempty"`
	LastAccessedRegion string     `json:"lastAccessedRegion,omitempty"`
	ServiceNamespace   string     `json:"serviceNamespace"`
	EntityPath         string     `json:"entityPath"`
}

// ServiceLastAccessed contains details about the last time an AWS service was accessed by an entity.
type ServiceLastAccessed struct {
	ServiceName                string                      `json:"serviceName"`
	ServiceNamespace           string                      `json:"serviceNamespace"`
	LastAuthenticated          *time.Time                  `json:"lastAuthenticated,omitempty"`
	LastAuthenticatedRegion    string                      `json:"lastAuthenticatedRegion,omitempty"`
	TotalAuthenticatedEntities int                         `json:"totalAuthenticatedEntities"`
	TrackedActionsLastAccessed []TrackedActionLastAccessed `json:"trackedActionsLastAccessed,omitempty"`
}

// ServiceLastAccessedDetailsJobStore provides persistent storage for Service Last Accessed Details jobs.
type ServiceLastAccessedDetailsJobStore struct {
	*common.BaseStore
}

// NewServiceLastAccessedDetailsJobStore creates a new store backed by the given storage.
func NewServiceLastAccessedDetailsJobStore(store storage.BasicStorage) *ServiceLastAccessedDetailsJobStore {
	return &ServiceLastAccessedDetailsJobStore{
		BaseStore: common.NewBaseStore(store.Bucket(serviceLastAccessedBucketName), "iam"),
	}
}

// Get retrieves a job by its ID. Returns nil if no job exists with the given ID.
func (s *ServiceLastAccessedDetailsJobStore) Get(jobID string) (*ServiceLastAccessedJob, error) {
	data, err := s.Bucket().Get([]byte(jobID))
	if err != nil {
		return nil, NewStoreError("get_service_last_accessed_job", err)
	}
	if data == nil {
		return nil, nil
	}
	var job ServiceLastAccessedJob
	if err := json.Unmarshal(data, &job); err != nil {
		return nil, NewStoreError("get_service_last_accessed_job", err)
	}
	return &job, nil
}

// Put persists a job to storage keyed by its JobID.
func (s *ServiceLastAccessedDetailsJobStore) Put(job *ServiceLastAccessedJob) error {
	return s.BaseStore.Put(job.JobID, job)
}
