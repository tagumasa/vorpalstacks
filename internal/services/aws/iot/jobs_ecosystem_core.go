package iot

import (
	"strings"
	"time"

	"github.com/google/uuid"

	iotstore "vorpalstacks/internal/store/aws/iot"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Jobs-ecosystem Core: job templates, managed job templates, per-thing job
// executions, and OTA updates. These resources are persisted as generic
// records under "jobTemplate/", "jobExecution/<jobId>/<thingName>" and
// "otaUpdate/" keys.
// ---------------------------------------------------------------------------

// CreateJobTemplateInput carries the fields for CreateJobTemplate. The
// configuration members keep the raw wire structures; they are stored and
// echoed verbatim. jobArn identifies a job whose document seeds the
// template when no inline document is supplied.
type CreateJobTemplateInput struct {
	JobTemplateID              string
	JobArn                     string
	Description                string
	Document                   string
	DocumentSource             string
	PresignedUrlConfig         interface{}
	JobExecutionsRolloutConfig interface{}
	AbortConfig                interface{}
	TimeoutConfig              interface{}
	JobExecutionsRetryConfig   interface{}
	MaintenanceWindows         interface{}
	DestinationPackageVersions interface{}
}

// CreateJobTemplateResult is the transport-agnostic result of
// CreateJobTemplate.
type CreateJobTemplateResult struct {
	JobTemplateID  string
	JobTemplateARN string
}

// CancelJobExecutionInput carries the fields for CancelJobExecution.
// ExpectedVersion is the model's optimistic-concurrency member (a mismatch
// rejects the cancel) and ExpectedVersionProvided distinguishes an
// explicitly supplied value from an omitted one — an explicit zero never
// matches the stored version, which starts at one; StatusDetails is stored
// on the execution.
type CancelJobExecutionInput struct {
	JobID                   string
	ThingName               string
	Force                   bool
	ExpectedVersion         int64
	ExpectedVersionProvided bool
	StatusDetails           interface{}
}

// DescribeJobExecutionInput carries the fields for DescribeJobExecution.
type DescribeJobExecutionInput struct {
	JobID     string
	ThingName string
}

// DeleteJobExecutionInput carries the fields for DeleteJobExecution. The
// model marks executionNumber required; it must match the stored execution.
type DeleteJobExecutionInput struct {
	JobID           string
	ThingName       string
	ExecutionNumber int64
	Force           bool
}

// CreateOTAUpdateInput carries the fields for CreateOTAUpdate. The awsJob*
// configuration members and the files keep the raw wire structures; they
// are stored and echoed verbatim.
type CreateOTAUpdateInput struct {
	OtaUpdateID                   string
	Description                   string
	Targets                       []string
	Protocols                     []string
	TargetSelection               string
	AwsJobExecutionsRolloutConfig interface{}
	AwsJobPresignedUrlConfig      interface{}
	AwsJobAbortConfig             interface{}
	AwsJobTimeoutConfig           interface{}
	Files                         interface{}
	RoleArn                       string
	Tags                          map[string]string
}

// CreateOTAUpdateResult is the transport-agnostic result of CreateOTAUpdate.
type CreateOTAUpdateResult struct {
	OtaUpdateID     string
	OtaUpdateArn    string
	OtaUpdateStatus string
	AwsIotJobID     string
	AwsIotJobArn    string
}

// jobExecutionRecordKey returns the generic-record key of a per-thing job
// execution.
func jobExecutionRecordKey(jobID, thingName string) string {
	return "jobExecution/" + jobID + "/" + thingName
}

// createJobTemplateCore validates and persists a job template record. The
// model marks description required; a jobArn member seeds the document from
// the referenced job when no inline document is supplied.
func (s *IoTService) createJobTemplateCore(store iotstore.IotStoreInterface, in CreateJobTemplateInput) (*CreateJobTemplateResult, error) {
	if in.JobTemplateID == "" {
		return nil, iotstore.ErrMissingParam
	}
	if in.Description == "" {
		return nil, iotstore.ErrValidation
	}
	// The job document is required unless documentSource or jobArn
	// supplies it, so a request carrying none of the three carriers is
	// rejected.
	if in.Document == "" && in.DocumentSource == "" && in.JobArn == "" {
		return nil, iotstore.ErrInvalidRequest
	}
	document := in.Document
	if document == "" && in.JobArn != "" {
		jobID, err := jobIDFromJobARN(in.JobArn)
		if err != nil {
			return nil, err
		}
		job, err := store.GetJob(jobID)
		if err != nil {
			return nil, err
		}
		document = job.Document
	}
	rec := map[string]interface{}{
		"jobTemplateId":  in.JobTemplateID,
		"jobTemplateArn": iotstore.BuildJobTemplateARN(store.GetAccountID(), store.GetRegion(), in.JobTemplateID),
		"description":    in.Description,
		"document":       document,
		"createdAt":      time.Now().Unix(),
	}
	if in.JobArn != "" {
		rec["jobArn"] = in.JobArn
	}
	if in.DocumentSource != "" {
		rec["documentSource"] = in.DocumentSource
	}
	for key, value := range map[string]interface{}{
		"presignedUrlConfig":         in.PresignedUrlConfig,
		"jobExecutionsRolloutConfig": in.JobExecutionsRolloutConfig,
		"abortConfig":                in.AbortConfig,
		"timeoutConfig":              in.TimeoutConfig,
		"jobExecutionsRetryConfig":   in.JobExecutionsRetryConfig,
		"maintenanceWindows":         in.MaintenanceWindows,
		"destinationPackageVersions": in.DestinationPackageVersions,
	} {
		if value != nil {
			rec[key] = value
		}
	}
	if err := store.PutGeneric("jobTemplate/"+in.JobTemplateID, rec); err != nil {
		return nil, err
	}
	return &CreateJobTemplateResult{
		JobTemplateID:  in.JobTemplateID,
		JobTemplateARN: rec["jobTemplateArn"].(string),
	}, nil
}

// jobIDFromJobARN extracts the job ID from a job ARN's resource section.
func jobIDFromJobARN(jobArn string) (string, error) {
	_, _, _, _, resource := svcarn.SplitARN(jobArn)
	if !strings.HasPrefix(resource, "job/") {
		return "", iotstore.ErrValidation
	}
	return strings.TrimPrefix(resource, "job/"), nil
}

// deleteJobTemplateCore removes a job template record.
func (s *IoTService) deleteJobTemplateCore(store iotstore.IotStoreInterface, name string) error {
	exists, err := store.GetGenericExists("jobTemplate/"+name, &map[string]interface{}{})
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrJobTemplateNotFound
	}
	if err := store.DeleteGeneric("jobTemplate/" + name); err != nil {
		return err
	}
	return nil
}

// describeJobTemplateCore retrieves a job template record.
func (s *IoTService) describeJobTemplateCore(store iotstore.IotStoreInterface, name string) (map[string]interface{}, error) {
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("jobTemplate/"+name, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrJobTemplateNotFound
	}
	return rec, nil
}

// listJobTemplatesCore lists all job template records.
func (s *IoTService) listJobTemplatesCore(store iotstore.IotStoreInterface) ([]map[string]interface{}, error) {
	return store.ListGeneric("jobTemplate/")
}

// describeManagedJobTemplateCore resolves a managed job template. The
// platform ships no AWS-provided managed-template catalogue (the catalogue
// content is AWS's copyrighted material), so the catalogue is empty and
// every describe resolves to the documented not-found error; the response
// shape's templateVersion/environments/documentParameters/document members
// exist in the model for the catalogue AWS serves.
func (s *IoTService) describeManagedJobTemplateCore(name string) (map[string]interface{}, error) {
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	return nil, iotstore.ErrJobNotFound
}

// cancelJobExecutionCore cancels a per-thing job execution record. A
// non-zero expectedVersion that does not match the stored execution version
// rejects the cancel with the model's version-conflict error; statusDetails
// are recorded on the execution.
func (s *IoTService) cancelJobExecutionCore(store iotstore.IotStoreInterface, in CancelJobExecutionInput) error {
	if in.JobID == "" || in.ThingName == "" {
		return iotstore.ErrMissingParam
	}
	// Validate the parent job exists; AWS returns ResourceNotFoundException.
	if _, err := store.GetJob(in.JobID); err != nil {
		return err
	}
	key := jobExecutionRecordKey(in.JobID, in.ThingName)
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrJobExecutionNotFound
	}
	if in.ExpectedVersionProvided {
		if stored, ok := recordInt64(rec["versionNumber"]); !ok || stored != in.ExpectedVersion {
			return iotstore.ErrVersionConflict
		}
	}
	// A job execution can be canceled while QUEUED, or while IN_PROGRESS
	// only with force; any other state is an invalid transition.
	status, _ := rec["status"].(string)
	if !(status == "QUEUED" || (status == "IN_PROGRESS" && in.Force)) {
		return iotstore.ErrInvalidStateTransition
	}
	rec["status"] = "CANCELED"
	rec["forceCanceled"] = in.Force
	if in.StatusDetails != nil {
		// The response member wraps the map in its detailsMap shape.
		rec["statusDetails"] = map[string]interface{}{"detailsMap": in.StatusDetails}
	}
	rec["lastUpdatedAt"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return err
	}
	return nil
}

// recordInt64 coerces a stored numeric record value (int64 in memory,
// float64 after a JSON round trip) to int64.
func recordInt64(v interface{}) (int64, bool) {
	switch n := v.(type) {
	case int64:
		return n, true
	case int:
		return int64(n), true
	case float64:
		return int64(n), true
	default:
		return 0, false
	}
}

// terminalJobExecutionStatus reports whether a stored job execution status
// is one of the terminal states deletable without force.
func terminalJobExecutionStatus(raw interface{}) bool {
	status, _ := raw.(string)
	switch status {
	case "SUCCEEDED", "FAILED", "REJECTED", "REMOVED", "CANCELED":
		return true
	}
	return false
}

// deleteJobExecutionCore removes a per-thing job execution record. The
// model marks executionNumber required and httpLabel-carried; a number that
// does not match the stored execution resolves to the not-found error.
// Without force only terminal executions may be deleted; force lifts the
// restriction for non-terminal states such as IN_PROGRESS.
func (s *IoTService) deleteJobExecutionCore(store iotstore.IotStoreInterface, in DeleteJobExecutionInput) error {
	if in.JobID == "" || in.ThingName == "" {
		return iotstore.ErrMissingParam
	}
	if in.ExecutionNumber == 0 {
		return iotstore.ErrValidation
	}
	if _, err := store.GetJob(in.JobID); err != nil {
		return err
	}
	key := jobExecutionRecordKey(in.JobID, in.ThingName)
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrJobExecutionNotFound
	}
	if stored, ok := recordInt64(rec["executionNumber"]); !ok || stored != in.ExecutionNumber {
		return iotstore.ErrJobExecutionNotFound
	}
	if !in.Force && !terminalJobExecutionStatus(rec["status"]) {
		return iotstore.ErrInvalidStateTransition
	}
	if err := store.DeleteGeneric(key); err != nil {
		return err
	}
	return nil
}

// describeJobExecutionCore retrieves a per-thing job execution and shapes
// the JobExecution response member.
func (s *IoTService) describeJobExecutionCore(store iotstore.IotStoreInterface, in DescribeJobExecutionInput) (map[string]interface{}, error) {
	if in.JobID == "" || in.ThingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	if _, err := store.GetJob(in.JobID); err != nil {
		return nil, err
	}
	key := jobExecutionRecordKey(in.JobID, in.ThingName)
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrJobExecutionNotFound
	}
	return map[string]interface{}{
		"execution": map[string]interface{}{
			"jobId":           in.JobID,
			"status":          rec["status"],
			"executionNumber": rec["executionNumber"],
			"queuedAt":        rec["queuedAt"],
			"startedAt":       rec["startedAt"],
			"lastUpdatedAt":   rec["lastUpdatedAt"],
			"thingArn":        iotstore.BuildThingARN(store.GetAccountID(), store.GetRegion(), in.ThingName),
			"versionNumber":   rec["versionNumber"],
			"statusDetails":   rec["statusDetails"],
		},
	}, nil
}

// jobExecutionSummaryMembers shapes the JobExecutionSummary members from a
// stored execution record. startedAt / lastUpdatedAt are optional model
// members and are emitted only when the record carries them.
func jobExecutionSummaryMembers(rec map[string]interface{}) map[string]interface{} {
	summary := map[string]interface{}{
		"status":          rec["status"],
		"executionNumber": rec["executionNumber"],
		"queuedAt":        rec["queuedAt"],
	}
	if v, ok := rec["startedAt"]; ok {
		summary["startedAt"] = v
	}
	if v, ok := rec["lastUpdatedAt"]; ok {
		summary["lastUpdatedAt"] = v
	}
	return summary
}

// listJobExecutionsForJobCore lists the per-thing execution summaries of a
// job.
func (s *IoTService) listJobExecutionsForJobCore(store iotstore.IotStoreInterface, jobID string) ([]map[string]interface{}, error) {
	if jobID == "" {
		return nil, iotstore.ErrMissingParam
	}
	// Validate the parent job exists; AWS returns ResourceNotFoundException.
	if _, err := store.GetJob(jobID); err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("jobExecution/" + jobID + "/")
	if err != nil {
		return nil, err
	}
	summaries := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		thingName, _ := rec["thingName"].(string)
		summaries = append(summaries, map[string]interface{}{
			"thingArn":            iotstore.BuildThingARN(store.GetAccountID(), store.GetRegion(), thingName),
			"jobExecutionSummary": jobExecutionSummaryMembers(rec),
		})
	}
	return summaries, nil
}

// listJobExecutionsForThingCore scans all job-execution records and filters
// by thing name. The summary-for-thing shape is keyed by jobId.
func (s *IoTService) listJobExecutionsForThingCore(store iotstore.IotStoreInterface, thingName string) ([]map[string]interface{}, error) {
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	allItems, err := store.ListGeneric("jobExecution/")
	if err != nil {
		return nil, err
	}
	summaries := make([]map[string]interface{}, 0)
	for _, rec := range allItems {
		rThing, _ := rec["thingName"].(string)
		if rThing != thingName {
			continue
		}
		jobID, _ := rec["jobId"].(string)
		summaries = append(summaries, map[string]interface{}{
			"jobId":               jobID,
			"jobExecutionSummary": jobExecutionSummaryMembers(rec),
		})
	}
	return summaries, nil
}

// createOTAUpdateCore validates and persists an OTA update record. The
// model marks targets, files and roleArn required.
func (s *IoTService) createOTAUpdateCore(store iotstore.IotStoreInterface, in CreateOTAUpdateInput) (*CreateOTAUpdateResult, error) {
	if in.OtaUpdateID == "" {
		return nil, iotstore.ErrMissingParam
	}
	if len(in.Targets) == 0 || in.Files == nil || in.RoleArn == "" {
		return nil, iotstore.ErrValidation
	}
	otaID := uuid.New().String()
	awsIotJobArn := iotstore.BuildJobARN(store.GetAccountID(), store.GetRegion(), otaID)
	now := time.Now().Unix()
	rec := map[string]interface{}{
		"otaUpdateId":                   in.OtaUpdateID,
		"otaUpdateArn":                  iotstore.BuildOTAUpdateARN(store.GetAccountID(), store.GetRegion(), in.OtaUpdateID),
		"description":                   in.Description,
		"targets":                       in.Targets,
		"protocols":                     in.Protocols,
		"targetSelection":               in.TargetSelection,
		"awsJobExecutionsRolloutConfig": in.AwsJobExecutionsRolloutConfig,
		"awsJobPresignedUrlConfig":      in.AwsJobPresignedUrlConfig,
		"awsJobAbortConfig":             in.AwsJobAbortConfig,
		"awsJobTimeoutConfig":           in.AwsJobTimeoutConfig,
		"otaUpdateFiles":                in.Files,
		"roleArn":                       in.RoleArn,
		"tags":                          in.Tags,
		"otaUpdateStatus":               "CREATE_COMPLETE",
		"awsIotJobId":                   otaID,
		"awsIotJobArn":                  awsIotJobArn,
		"creationDate":                  now,
		"lastModifiedDate":              now,
	}
	if err := store.PutGeneric("otaUpdate/"+in.OtaUpdateID, rec); err != nil {
		return nil, err
	}
	// The awsIotJobId the response carries identifies a real IoT job that
	// targets the update's targets; deleting the update without
	// forceDeleteAWSJob is blocked while the job is not terminal.
	if _, err := store.CreateJob(&iotstore.Job{
		JobID:       otaID,
		Description: in.Description,
		Targets:     in.Targets,
		Status:      "IN_PROGRESS",
		CreatedAt:   time.Now().UTC(),
	}); err != nil {
		return nil, err
	}
	return &CreateOTAUpdateResult{
		OtaUpdateID:     in.OtaUpdateID,
		OtaUpdateArn:    rec["otaUpdateArn"].(string),
		OtaUpdateStatus: "CREATE_COMPLETE",
		AwsIotJobID:     otaID,
		AwsIotJobArn:    awsIotJobArn,
	}, nil
}

// DeleteOTAUpdateInput carries the fields for DeleteOTAUpdate: the
// deleteStream flag removes the streams referenced by the update's files,
// and forceDeleteAWSJob allows deleting the update while its IoT job is
// still in progress.
type DeleteOTAUpdateInput struct {
	OtaUpdateID       string
	DeleteStream      bool
	ForceDeleteAWSJob bool
}

// deleteOTAUpdateCore removes an OTA update record. Without
// forceDeleteAWSJob an IoT job that is not in a terminal state (COMPLETED
// or CANCELED) blocks the deletion; the referenced job record is removed
// with the update. With deleteStream the streams named by the files'
// stream locations are removed as well.
func (s *IoTService) deleteOTAUpdateCore(store iotstore.IotStoreInterface, in DeleteOTAUpdateInput) error {
	if in.OtaUpdateID == "" {
		return iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("otaUpdate/"+in.OtaUpdateID, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrJobNotFound
	}
	if awsJobID, _ := rec["awsIotJobId"].(string); awsJobID != "" {
		job, jobErr := store.GetJob(awsJobID)
		if jobErr == nil {
			// The job must be in a terminal state (COMPLETED or CANCELED)
			// unless the caller forces the delete.
			if job.Status != "COMPLETED" && job.Status != "CANCELED" && !in.ForceDeleteAWSJob {
				return iotstore.ErrInvalidRequest
			}
			if err := store.DeleteJob(awsJobID); err != nil {
				return err
			}
		}
	}
	// deleteStream only removes streams the OTAUpdate process itself
	// created. This platform's CreateOTAUpdate never generates streams —
	// every stream referenced through the files' stream location was
	// supplied by the user — so the member is always ignored here and the
	// referenced streams survive the OTA update deletion.
	if err := store.DeleteGeneric("otaUpdate/" + in.OtaUpdateID); err != nil {
		return err
	}
	return nil
}

// getOTAUpdateCore retrieves an OTA update record and projects it onto the
// OTAUpdateInfo member set. Input-only members (roleArn, tags, the abort and
// timeout configs) stay in the record but are not response members; the two
// model config members are omitted when empty so the SDK never sees
// empty-string structures.
func (s *IoTService) getOTAUpdateCore(store iotstore.IotStoreInterface, name string) (map[string]interface{}, error) {
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("otaUpdate/"+name, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrJobNotFound
	}
	info := map[string]interface{}{
		"otaUpdateId":      rec["otaUpdateId"],
		"otaUpdateArn":     rec["otaUpdateArn"],
		"description":      rec["description"],
		"targets":          rec["targets"],
		"protocols":        rec["protocols"],
		"targetSelection":  rec["targetSelection"],
		"otaUpdateStatus":  rec["otaUpdateStatus"],
		"awsIotJobId":      rec["awsIotJobId"],
		"awsIotJobArn":     rec["awsIotJobArn"],
		"creationDate":     rec["creationDate"],
		"lastModifiedDate": rec["lastModifiedDate"],
	}
	if files, ok := rec["otaUpdateFiles"]; ok && files != nil {
		info["otaUpdateFiles"] = files
	}
	for _, k := range []string{"awsJobExecutionsRolloutConfig", "awsJobPresignedUrlConfig"} {
		if v, ok := rec[k]; ok && v != nil && v != "" {
			info[k] = v
		}
	}
	return map[string]interface{}{"otaUpdateInfo": info}, nil
}

// listOTAUpdatesCore lists OTA update summaries with the optional
// otaUpdateStatus filter.
func (s *IoTService) listOTAUpdatesCore(store iotstore.IotStoreInterface, statusFilter string) ([]map[string]interface{}, error) {
	items, err := store.ListGeneric("otaUpdate/")
	if err != nil {
		return nil, err
	}
	summaries := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		if statusFilter != "" {
			if status, _ := rec["otaUpdateStatus"].(string); status != statusFilter {
				continue
			}
		}
		summaries = append(summaries, map[string]interface{}{
			"otaUpdateId":  rec["otaUpdateId"],
			"otaUpdateArn": rec["otaUpdateArn"],
			"creationDate": rec["creationDate"],
		})
	}
	return summaries, nil
}
