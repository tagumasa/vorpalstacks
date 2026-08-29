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

	// Per-target execution records are written before the job record so
	// the create is crash-safe in one direction: the keys are
	// deterministic, so an interrupted create leaves overwriteable
	// orphans rather than a job missing executions, and a retry only
	// fails with ErrJobAlreadyExists once the job itself is durable.
	if err := s.createJobExecutionsLocked(job.JobID, job.Targets); err != nil {
		return nil, err
	}

	if err := s.jobPS.Create(job); err != nil {
		return nil, err
	}
	return job, nil
}

// createJobExecutionsLocked writes a QUEUED execution record per target
// thing. The records are keyed deterministically, so re-running the call
// for an already-recorded target is an overwrite with identical content,
// which keeps partial failures recoverable by a plain retry.
// Callers must hold s.mu.
func (s *IotStore) createJobExecutionsLocked(jobID string, targets []string) error {
	now := time.Now().UTC().Unix()
	for _, targetARN := range targets {
		thingName := ThingNameFromARN(targetARN)
		if thingName == "" {
			continue
		}
		execRec := map[string]interface{}{
			"jobId":           jobID,
			"thingName":       thingName,
			"status":          "QUEUED",
			"executionNumber": int64(1),
			"queuedAt":        now,
			"versionNumber":   int64(1),
		}
		if err := s.genericKVBase.Put("jobExecution/"+jobID+"/"+thingName, execRec); err != nil {
			return err
		}
	}
	return nil
}

// AssociateJobTargets merges additional targets into an existing job and
// materialises a QUEUED execution record for each genuinely new target,
// so associated things show up in job execution listings the same way
// targets of a fresh CreateJob do. The merge, the execution records and
// the job update happen under one lock, keeping the job's target list and
// its executions consistent.
func (s *IotStore) AssociateJobTargets(jobID string, newTargets []string, comment string) (*Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, err := s.jobPS.Get(jobID)
	if err != nil {
		return nil, err
	}
	if job == nil || job.JobID == "" {
		return nil, ErrJobNotFound
	}

	seen := make(map[string]bool, len(job.Targets))
	for _, t := range job.Targets {
		seen[t] = true
	}
	var added []string
	for _, t := range newTargets {
		if !seen[t] {
			job.Targets = append(job.Targets, t)
			seen[t] = true
			added = append(added, t)
		}
	}
	if comment != "" && job.Description == "" {
		job.Description = comment
	}

	if err := s.createJobExecutionsLocked(job.JobID, added); err != nil {
		return nil, err
	}

	job.LastUpdatedAt = time.Now().UTC()
	if err := s.jobPS.Update(job); err != nil {
		return nil, err
	}
	return job, nil
}

func (s *IotStore) GetJob(jobID string) (*Job, error) {
	return s.jobPS.Get(jobID)
}

func (s *IotStore) DeleteJob(jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Job executions are stored as generic-KV records keyed
	// "jobExecution/<jobID>/<thingName>"; deleting the job must remove
	// them as well, otherwise executions of a deleted job keep showing
	// up in ListJobExecutionsForThing.
	if err := s.genericKVBase.ScanPrefix("jobExecution/"+jobID+"/", func(key string, _ []byte) error {
		return s.genericKVBase.Delete(key)
	}); err != nil {
		return err
	}
	return s.jobPS.DeleteIfExists(jobID)
}

// JobListFilters carries the ListJobs query filters: the status and
// target-selection enum filters plus the thing-group target, which matches
// a job whose target list carries the group's ARN.
type JobListFilters struct {
	Status          string
	TargetSelection string
	ThingGroupARN   string
}

func (s *IotStore) ListJobs(opts common.ListOptions, filters JobListFilters) (*common.ListResult[Job], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var filter func(*pb.Job) bool
	hasStatus := filters.Status != ""
	hasTargetSelection := filters.TargetSelection != ""
	hasGroup := filters.ThingGroupARN != ""
	if hasStatus || hasTargetSelection || hasGroup {
		filter = func(j *pb.Job) bool {
			if hasStatus && j.Status != filters.Status {
				return false
			}
			if hasTargetSelection && j.TargetSelection != filters.TargetSelection {
				return false
			}
			if hasGroup && !jobTargetsContains(j.Targets, filters.ThingGroupARN) {
				return false
			}
			return true
		}
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

// jobTargetsContains reports whether the target list carries the given ARN.
func jobTargetsContains(targets []string, arn string) bool {
	for _, t := range targets {
		if t == arn {
			return true
		}
	}
	return false
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
	if opts.ReasonCode != "" {
		existing.ReasonCode = opts.ReasonCode
	}
	if opts.Comment != "" {
		existing.Comment = opts.Comment
	}
	existing.LastUpdatedAt = time.Now().UTC()
	return existing, s.jobPS.Update(existing)
}
