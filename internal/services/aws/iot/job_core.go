package iot

import (
	"time"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Job Core. Jobs distribute documents to targeted things; per-target
// execution records are materialised by the store alongside the job record.
// ---------------------------------------------------------------------------

// CreateJobInput carries every CreateJobRequest member. TagsMalformed
// preserves the wire reality that the tags member arrives as a possibly
// malformed JSON string; the Core rejects it after the targetSelection
// validation, matching the historical check order.
type CreateJobInput struct {
	JobID                      string
	Description                string
	Document                   string
	Targets                    []string
	TargetSelection            string
	NamespaceID                string
	JobTemplateARN             string
	PresignedUrlConfig         string
	JobExecutionsRolloutConfig string
	AbortConfig                string
	TimeoutConfig              string
	JobExecutionsRetryConfig   string
	DocumentParameters         string
	SchedulingConfig           string
	DestinationPackageVersions string
	Tags                       map[string]string
	TagsMalformed              bool
}

// UpdateJobInput carries the fields that UpdateJob applies. All members are
// optional strings; an empty value leaves the stored field untouched, same as
// the wire semantics.
type UpdateJobInput struct {
	JobID                      string
	Description                string
	NamespaceID                string
	PresignedUrlConfig         string
	JobExecutionsRolloutConfig string
	AbortConfig                string
	TimeoutConfig              string
	JobExecutionsRetryConfig   string
}

// AssociateTargetsWithJobInput carries the fields for
// AssociateTargetsWithJob.
type AssociateTargetsWithJobInput struct {
	JobID   string
	Targets []string
	Comment string
}

// ListJobsInput holds the parameters for listing jobs: the status filter
// plus the target-selection and thing-group filters the model carries as
// query members.
type ListJobsInput struct {
	Status          string
	TargetSelection string
	ThingGroupName  string
	ThingGroupID    string
	NextToken       string
	MaxItems        int
}

// CreateJobResult is the transport-agnostic result of CreateJob.
type CreateJobResult struct {
	Job *iotstore.Job
}

// DescribeJobResult is the transport-agnostic result of DescribeJob.
type DescribeJobResult struct {
	Job *iotstore.Job
}

// CancelJobInput carries the CancelJobRequest members: the optional
// reasonCode and comment are recorded on the job.
type CancelJobInput struct {
	JobID      string
	ReasonCode string
	Comment    string
}

// CancelJobResult is the transport-agnostic result of CancelJob (the
// response carries the model's jobArn/jobId/description members).
type CancelJobResult struct {
	Job *iotstore.Job
}

// GetJobDocumentResult is the transport-agnostic result of GetJobDocument.
type GetJobDocumentResult struct {
	Document string
}

// AssociateTargetsWithJobResult is the transport-agnostic result of
// AssociateTargetsWithJob.
type AssociateTargetsWithJobResult struct {
	JobARN      string
	JobID       string
	Description string
}

// ListJobsResult is the transport-agnostic result of ListJobs.
type ListJobsResult struct {
	Jobs      []*iotstore.Job
	NextToken string
}

// createJobCore validates and persists a new job.
func (s *IoTService) createJobCore(store iotstore.IotStoreInterface, in CreateJobInput) (*CreateJobResult, error) {
	if in.JobID == "" {
		return nil, iotstore.ErrMissingParam
	}

	job := &iotstore.Job{
		JobID:                      in.JobID,
		Description:                in.Description,
		Document:                   in.Document,
		Targets:                    in.Targets,
		Status:                     "IN_PROGRESS",
		TargetSelection:            in.TargetSelection,
		NamespaceID:                in.NamespaceID,
		JobTemplateARN:             in.JobTemplateARN,
		PresignedUrlConfig:         in.PresignedUrlConfig,
		JobExecutionsRolloutConfig: in.JobExecutionsRolloutConfig,
		AbortConfig:                in.AbortConfig,
		TimeoutConfig:              in.TimeoutConfig,
		JobExecutionsRetryConfig:   in.JobExecutionsRetryConfig,
		DocumentParameters:         in.DocumentParameters,
		SchedulingConfig:           in.SchedulingConfig,
		DestinationPackageVersions: in.DestinationPackageVersions,
		CreatedAt:                  time.Now().UTC(),
		LastUpdatedAt:              time.Now().UTC(),
	}
	if job.TargetSelection == "" {
		// An omitted targetSelection is documented as SNAPSHOT.
		job.TargetSelection = "SNAPSHOT"
	}

	if job.TargetSelection != "" {
		if err := ValidateTargetSelection(job.TargetSelection); err != nil {
			return nil, err
		}
	}

	if in.TagsMalformed {
		return nil, iotstore.ErrValidation
	}
	if len(in.Tags) > 0 {
		job.Tags = in.Tags
	}

	created, err := store.CreateJob(job)
	if err != nil {
		return nil, err
	}

	// Per-target execution records are written by the store inside
	// CreateJob, under the same lock and before the job record, so the
	// job never becomes visible with missing executions.

	return &CreateJobResult{Job: created}, nil
}

// describeJobCore retrieves a single job by ID.
func (s *IoTService) describeJobCore(store iotstore.IotStoreInterface, jobID string) (*DescribeJobResult, error) {
	if jobID == "" {
		return nil, iotstore.ErrMissingParam
	}
	job, err := store.GetJob(jobID)
	if err != nil {
		return nil, err
	}
	return &DescribeJobResult{Job: job}, nil
}

// deleteJobCore removes a job and cleans up its tags.
func (s *IoTService) deleteJobCore(store iotstore.IotStoreInterface, jobID string) error {
	if jobID == "" {
		return iotstore.ErrMissingParam
	}

	arn := iotstore.BuildJobARN(store.GetAccountID(), store.GetRegion(), jobID)
	_ = store.DeleteAllTags(arn)

	if err := store.DeleteJob(jobID); err != nil {
		return err
	}
	return nil
}

// listJobsCore lists jobs with the status, target-selection, and
// thing-group query filters. An unknown thing group matches no jobs.
func (s *IoTService) listJobsCore(store iotstore.IotStoreInterface, in ListJobsInput) (*ListJobsResult, error) {
	if in.TargetSelection != "" {
		if err := ValidateTargetSelection(in.TargetSelection); err != nil {
			return nil, err
		}
	}
	filters := iotstore.JobListFilters{
		Status:          in.Status,
		TargetSelection: in.TargetSelection,
	}
	switch {
	case in.ThingGroupName != "":
		// An unresolvable group name filters everything out.
		if group, err := store.GetThingGroup(in.ThingGroupName); err == nil {
			filters.ThingGroupARN = group.GroupARN
		} else {
			return &ListJobsResult{Jobs: []*iotstore.Job{}}, nil
		}
	case in.ThingGroupID != "":
		opts := iotstoreListOpts(0, "")
		resolved := false
		for !resolved {
			groups, err := store.ListThingGroups(opts, "")
			if err != nil {
				return nil, err
			}
			for _, g := range groups.Items {
				if g.GroupID == in.ThingGroupID {
					filters.ThingGroupARN = g.GroupARN
					resolved = true
					break
				}
			}
			if resolved || groups.NextMarker == "" {
				break
			}
			opts = iotstoreListOpts(0, groups.NextMarker)
		}
		if !resolved {
			return &ListJobsResult{Jobs: []*iotstore.Job{}}, nil
		}
	}
	maxItems := in.MaxItems
	if maxItems <= 0 {
		maxItems = 100
	}
	opts := iotstoreListOpts(maxItems, in.NextToken)
	jobs, err := store.ListJobs(opts, filters)
	if err != nil {
		return nil, err
	}
	return &ListJobsResult{
		Jobs:      jobs.Items,
		NextToken: jobs.NextMarker,
	}, nil
}

// cancelJobCore cancels a running or scheduled job, recording the optional
// reason code and comment supplied by the caller.
func (s *IoTService) cancelJobCore(store iotstore.IotStoreInterface, in CancelJobInput) (*CancelJobResult, error) {
	if in.JobID == "" {
		return nil, iotstore.ErrMissingParam
	}

	job, err := store.GetJob(in.JobID)
	if err != nil {
		return nil, err
	}
	if job.Status != "IN_PROGRESS" && job.Status != "SCHEDULED" {
		return nil, iotstore.ErrInvalidRequest
	}

	opts := iotstore.JobUpdateOpts{Status: "CANCELED", ReasonCode: in.ReasonCode, Comment: in.Comment}
	job, err = store.UpdateJob(in.JobID, opts)
	if err != nil {
		return nil, err
	}
	return &CancelJobResult{Job: job}, nil
}

// getJobDocumentCore retrieves the job document content.
func (s *IoTService) getJobDocumentCore(store iotstore.IotStoreInterface, jobID string) (*GetJobDocumentResult, error) {
	if jobID == "" {
		return nil, iotstore.ErrMissingParam
	}
	job, err := store.GetJob(jobID)
	if err != nil {
		return nil, err
	}
	return &GetJobDocumentResult{Document: job.Document}, nil
}

// updateJobCore applies the supplied fields to a job. Only the
// UpdateJobRequest members are accepted; CancelJobRequest carries the
// status/reasonCode/comment members instead.
func (s *IoTService) updateJobCore(store iotstore.IotStoreInterface, in UpdateJobInput) error {
	if in.JobID == "" {
		return iotstore.ErrMissingParam
	}

	opts := iotstore.JobUpdateOpts{
		Description:                in.Description,
		NamespaceID:                in.NamespaceID,
		PresignedUrlConfig:         in.PresignedUrlConfig,
		JobExecutionsRolloutConfig: in.JobExecutionsRolloutConfig,
		AbortConfig:                in.AbortConfig,
		TimeoutConfig:              in.TimeoutConfig,
		JobExecutionsRetryConfig:   in.JobExecutionsRetryConfig,
	}
	_, err := store.UpdateJob(in.JobID, opts)
	return err
}

// associateTargetsWithJobCore merges new targets into a job.
func (s *IoTService) associateTargetsWithJobCore(store iotstore.IotStoreInterface, in AssociateTargetsWithJobInput) (*AssociateTargetsWithJobResult, error) {
	if in.JobID == "" {
		return nil, iotstore.ErrMissingParam
	}

	// The store merges the targets and materialises a QUEUED execution
	// record for each genuinely new thing under one lock, so the job's
	// target list and its executions cannot drift apart.
	job, err := store.AssociateJobTargets(in.JobID, in.Targets, in.Comment)
	if err != nil {
		return nil, err
	}
	return &AssociateTargetsWithJobResult{
		JobARN:      job.JobARN,
		JobID:       job.JobID,
		Description: job.Description,
	}, nil
}
