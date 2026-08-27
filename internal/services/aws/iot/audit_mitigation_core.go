package iot

import (
	"time"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Audit Mitigation Actions Task Core. Task records persist under
// "auditMitigationTask/<taskId>" so that Cancel/Describe resolve the
// identifier and return ResourceNotFoundException for unknown task ids,
// matching the Smithy error trait set on each operation.
// ---------------------------------------------------------------------------

// StartAuditMitigationActionsTaskInput carries the parsed
// StartAuditMitigationActionsTask request. Target and
// AuditCheckToActionsMapping keep the raw wire structures.
type StartAuditMitigationActionsTaskInput struct {
	TaskID                     string
	Target                     map[string]interface{}
	AuditCheckToActionsMapping map[string]interface{}
}

// startAuditMitigationActionsTaskCore persists an audit mitigation task
// record and returns the task id.
func (s *IoTService) startAuditMitigationActionsTaskCore(store iotstore.IotStoreInterface, in StartAuditMitigationActionsTaskInput) (string, error) {
	if in.TaskID == "" {
		return "", iotstore.ErrMissingParam
	}
	if in.Target == nil || in.AuditCheckToActionsMapping == nil {
		return "", iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{
		"taskId":                     in.TaskID,
		"status":                     "IN_PROGRESS",
		"startTime":                  time.Now().UTC().Unix(),
		"target":                     in.Target,
		"auditCheckToActionsMapping": in.AuditCheckToActionsMapping,
	}
	if err := store.PutGeneric("auditMitigationTask/"+in.TaskID, rec); err != nil {
		return "", err
	}
	return in.TaskID, nil
}

// cancelAuditMitigationActionsTaskCore transitions an audit mitigation
// task to CANCELED (the AuditMitigationActionsTaskStatus enum value); the
// record stays queryable via DescribeAuditMitigationActionsTask. An
// unknown task id yields ErrAuditMitigationTaskNotFound.
func (s *IoTService) cancelAuditMitigationActionsTaskCore(store iotstore.IotStoreInterface, taskId string) error {
	if taskId == "" {
		return iotstore.ErrMissingParam
	}
	key := "auditMitigationTask/" + taskId
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrAuditMitigationTaskNotFound
	}
	rec["status"] = "CANCELED"
	return store.PutGeneric(key, rec)
}

// AuditMitigationTaskRecord is the projected
// DescribeAuditMitigationActionsTask response payload
// (DescribeAuditMitigationActionsTaskResponse members only).
type AuditMitigationTaskRecord struct {
	TaskStatus                 interface{}
	StartTime                  interface{}
	EndTime                    interface{}
	Target                     interface{}
	AuditCheckToActionsMapping interface{}
}

// describeAuditMitigationActionsTaskCore loads an audit mitigation task
// record. An unknown task id yields ErrAuditMitigationTaskNotFound.
func (s *IoTService) describeAuditMitigationActionsTaskCore(store iotstore.IotStoreInterface, taskId string) (*AuditMitigationTaskRecord, error) {
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("auditMitigationTask/"+taskId, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditMitigationTaskNotFound
	}
	return &AuditMitigationTaskRecord{
		TaskStatus:                 rec["status"],
		StartTime:                  rec["startTime"],
		EndTime:                    rec["endTime"],
		Target:                     rec["target"],
		AuditCheckToActionsMapping: rec["auditCheckToActionsMapping"],
	}, nil
}

// AuditMitigationTaskListItem is one ListAuditMitigationActionsTasks entry
// (AuditMitigationActionsTaskMetadata member set).
type AuditMitigationTaskListItem struct {
	TaskID     interface{}
	StartTime  interface{}
	TaskStatus interface{}
}

// listAuditMitigationActionsTasksCore lists every audit mitigation task
// record.
func (s *IoTService) listAuditMitigationActionsTasksCore(store iotstore.IotStoreInterface) ([]AuditMitigationTaskListItem, error) {
	items, err := store.ListGeneric("auditMitigationTask/")
	if err != nil {
		return nil, err
	}
	tasks := make([]AuditMitigationTaskListItem, 0, len(items))
	for _, rec := range items {
		tasks = append(tasks, AuditMitigationTaskListItem{
			TaskID:     rec["taskId"],
			StartTime:  rec["startTime"],
			TaskStatus: rec["status"],
		})
	}
	return tasks, nil
}
