package iot

import (
	"time"

	"github.com/google/uuid"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Jobs-ecosystem Core: job templates, managed job templates, per-thing job
// executions, and OTA updates. These resources are persisted as generic
// records under "jobTemplate/", "jobExecution/<jobId>/<thingName>" and
// "otaUpdate/" keys.
// ---------------------------------------------------------------------------

// CreateJobTemplateInput carries the fields for CreateJobTemplate.
type CreateJobTemplateInput struct {
	JobTemplateID string
	Description   string
	Document      string
}

// CreateJobTemplateResult is the transport-agnostic result of
// CreateJobTemplate.
type CreateJobTemplateResult struct {
	JobTemplateID  string
	JobTemplateARN string
}

// CancelJobExecutionInput carries the fields for CancelJobExecution.
type CancelJobExecutionInput struct {
	JobID     string
	ThingName string
	Force     bool
}

// DescribeJobExecutionInput carries the fields for DescribeJobExecution.
type DescribeJobExecutionInput struct {
	JobID     string
	ThingName string
}

// CreateOTAUpdateInput carries the fields for CreateOTAUpdate.
type CreateOTAUpdateInput struct {
	OtaUpdateID                   string
	Description                   string
	Targets                       []string
	Protocols                     []string
	TargetSelection               string
	AwsJobExecutionsRolloutConfig string
	AwsJobPresignedUrlConfig      string
	AwsJobAbortConfig             string
	AwsJobTimeoutConfig           string
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

// createJobTemplateCore validates and persists a job template record.
func (s *IoTService) createJobTemplateCore(store iotstore.IotStoreInterface, in CreateJobTemplateInput) (*CreateJobTemplateResult, error) {
	if in.JobTemplateID == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{
		"jobTemplateId": in.JobTemplateID,
		"description":   in.Description,
		"document":      in.Document,
		"createdAt":     time.Now().Unix(),
	}
	if err := store.PutGeneric("jobTemplate/"+in.JobTemplateID, rec); err != nil {
		return nil, err
	}
	return &CreateJobTemplateResult{
		JobTemplateID:  in.JobTemplateID,
		JobTemplateARN: iotstore.BuildJobTemplateARN(store.GetAccountID(), store.GetRegion(), in.JobTemplateID),
	}, nil
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

// describeManagedJobTemplateCore shapes the synthetic managed-template
// description. Managed templates are AWS-provided and not persisted; the
// response is a stub describing the requested template.
func (s *IoTService) describeManagedJobTemplateCore(accountID, region, name string) map[string]interface{} {
	return map[string]interface{}{
		"templateName": name,
		"templateArn":  iotstore.BuildJobTemplateARN(accountID, region, name),
		"description":  "AWS-provided managed job template",
		"platform":     "Linux",
		"pathToDefine": "",
	}
}

// cancelJobExecutionCore cancels a per-thing job execution record.
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
	rec["status"] = "CANCELED"
	rec["forceCanceled"] = in.Force
	rec["lastUpdatedAt"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return err
	}
	return nil
}

// deleteJobExecutionCore removes a per-thing job execution record.
func (s *IoTService) deleteJobExecutionCore(store iotstore.IotStoreInterface, jobID, thingName string) error {
	if jobID == "" || thingName == "" {
		return iotstore.ErrMissingParam
	}
	if _, err := store.GetJob(jobID); err != nil {
		return err
	}
	key := jobExecutionRecordKey(jobID, thingName)
	exists, err := store.GetGenericExists(key, &map[string]interface{}{})
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrJobExecutionNotFound
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

// createOTAUpdateCore validates and persists an OTA update record.
func (s *IoTService) createOTAUpdateCore(store iotstore.IotStoreInterface, in CreateOTAUpdateInput) (*CreateOTAUpdateResult, error) {
	if in.OtaUpdateID == "" {
		return nil, iotstore.ErrMissingParam
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
	return &CreateOTAUpdateResult{
		OtaUpdateID:     in.OtaUpdateID,
		OtaUpdateArn:    rec["otaUpdateArn"].(string),
		OtaUpdateStatus: "CREATE_COMPLETE",
		AwsIotJobID:     otaID,
		AwsIotJobArn:    awsIotJobArn,
	}, nil
}

// deleteOTAUpdateCore removes an OTA update record.
func (s *IoTService) deleteOTAUpdateCore(store iotstore.IotStoreInterface, name string) error {
	exists, err := store.GetGenericExists("otaUpdate/"+name, &map[string]interface{}{})
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrJobNotFound
	}
	if err := store.DeleteGeneric("otaUpdate/" + name); err != nil {
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
	for _, k := range []string{"awsJobExecutionsRolloutConfig", "awsJobPresignedUrlConfig"} {
		if v, ok := rec[k]; ok && v != nil && v != "" {
			info[k] = v
		}
	}
	return map[string]interface{}{"otaUpdateInfo": info}, nil
}

// listOTAUpdatesCore lists OTA update summaries.
func (s *IoTService) listOTAUpdatesCore(store iotstore.IotStoreInterface) ([]map[string]interface{}, error) {
	items, err := store.ListGeneric("otaUpdate/")
	if err != nil {
		return nil, err
	}
	summaries := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		summaries = append(summaries, map[string]interface{}{
			"otaUpdateId":  rec["otaUpdateId"],
			"otaUpdateArn": rec["otaUpdateArn"],
			"creationDate": rec["creationDate"],
		})
	}
	return summaries, nil
}
