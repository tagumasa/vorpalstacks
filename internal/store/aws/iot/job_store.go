package iot

import (
	"time"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
	"vorpalstacks/internal/store/aws/common"
)

func (s *IotStore) CreateJob(job *Job) (*Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job.JobID == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.Job{}
	if err := s.jobsBase.GetProto(job.JobID, existing); err == nil {
		return nil, ErrJobAlreadyExists
	}
	job.JobARN = BuildJobARN(s.accountID, s.region, job.JobID)
	return job, s.jobPS.Create(job)
}

func (s *IotStore) GetJob(jobID string) (*Job, error) {
	return s.jobPS.Get(jobID)
}

func (s *IotStore) DeleteJob(jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.jobPS.DeleteIfExists(jobID)
}

func (s *IotStore) ListJobs(opts common.ListOptions, statusFilter string) (*common.ListResult[Job], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var filter func(*pb.Job) bool
	if statusFilter != "" {
		filter = func(j *pb.Job) bool { return j.Status == statusFilter }
	}
	result, err := common.ListProto(s.jobsBase, opts, func() *pb.Job { return &pb.Job{} }, filter)
	if err != nil {
		return nil, err
	}
	items := make([]*Job, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToJob(p))
	}
	return &common.ListResult[Job]{Items: items, NextMarker: result.NextMarker}, nil
}

func (s *IotStore) UpdateJob(jobID string, opts JobUpdateOpts) (*Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, err := s.jobPS.Get(jobID)
	if err != nil {
		return nil, err
	}
	if opts.Description != "" {
		existing.Description = opts.Description
	}
	if opts.Targets != nil {
		existing.Targets = opts.Targets
	}
	if opts.Status != "" {
		switch opts.Status {
		case "IN_PROGRESS", "CANCELED", "COMPLETED", "DELETION_IN_PROGRESS", "SCHEDULED":
			existing.Status = opts.Status
		default:
			return nil, ErrInvalidRequest
		}
	}
	if opts.ReasonCode != "" {
		existing.ReasonCode = opts.ReasonCode
	}
	if opts.Comment != "" {
		existing.Comment = opts.Comment
	}
	if opts.NamespaceID != "" {
		existing.NamespaceID = opts.NamespaceID
	}
	if opts.PresignedUrlConfig != "" {
		existing.PresignedUrlConfig = opts.PresignedUrlConfig
	}
	if opts.JobExecutionsRolloutConfig != "" {
		existing.JobExecutionsRolloutConfig = opts.JobExecutionsRolloutConfig
	}
	if opts.AbortConfig != "" {
		existing.AbortConfig = opts.AbortConfig
	}
	if opts.TimeoutConfig != "" {
		existing.TimeoutConfig = opts.TimeoutConfig
	}
	if opts.JobExecutionsRetryConfig != "" {
		existing.JobExecutionsRetryConfig = opts.JobExecutionsRetryConfig
	}
	if opts.SchedulingConfig != "" {
		existing.SchedulingConfig = opts.SchedulingConfig
	}
	if opts.ScheduledJobRollouts != "" {
		existing.ScheduledJobRollouts = opts.ScheduledJobRollouts
	}
	existing.LastUpdatedAt = time.Now().UTC()
	return existing, s.jobPS.Update(existing)
}
