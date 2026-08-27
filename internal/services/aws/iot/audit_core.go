package iot

import (
	"time"

	"github.com/google/uuid"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Audit Core (account audit configuration, on-demand audit tasks,
// audit findings). Audit task lifecycle mirrors the Detect Mitigation task
// pattern: Start persists the task id, Cancel/Describe enforce
// ResourceNotFoundException for unknown ids. Audit findings are not
// generated without a Defender engine, so DescribeAuditFinding always
// returns NotFound for an arbitrary id.
// ---------------------------------------------------------------------------

// describeAccountAuditConfigurationCore loads the account audit
// configuration, returning the documented default shape when none has been
// written.
func (s *IoTService) describeAccountAuditConfigurationCore(store iotstore.IotStoreInterface) (map[string]interface{}, error) {
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("config/accountAudit", &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return map[string]interface{}{
			"auditCheckConfigurations":              map[string]interface{}{},
			"auditNotificationTargetConfigurations": map[string]interface{}{},
		}, nil
	}
	return rec, nil
}

// UpdateAccountAuditConfigurationInput carries the parsed
// UpdateAccountAuditConfiguration request.
type UpdateAccountAuditConfigurationInput struct {
	RoleArn                               string
	AuditCheckConfigurations              map[string]interface{}
	AuditNotificationTargetConfigurations map[string]interface{}
}

// updateAccountAuditConfigurationCore applies an omission-based partial
// update to the account audit configuration: unspecified checks keep their
// existing enabled/disabled state, and an empty roleArn leaves the stored
// one unchanged.
func (s *IoTService) updateAccountAuditConfigurationCore(store iotstore.IotStoreInterface, in UpdateAccountAuditConfigurationInput) error {
	rec := map[string]interface{}{}
	if _, err := store.GetGenericExists("config/accountAudit", &rec); err != nil {
		return err
	}
	if in.RoleArn != "" {
		rec["roleArn"] = in.RoleArn
	}
	for _, entry := range []struct {
		key   string
		value map[string]interface{}
	}{
		{"auditCheckConfigurations", in.AuditCheckConfigurations},
		{"auditNotificationTargetConfigurations", in.AuditNotificationTargetConfigurations},
	} {
		if entry.value == nil {
			continue
		}
		existing, _ := rec[entry.key].(map[string]interface{})
		if existing == nil {
			existing = map[string]interface{}{}
		}
		for k, v := range entry.value {
			existing[k] = v
		}
		rec[entry.key] = existing
	}
	return store.PutGeneric("config/accountAudit", rec)
}

// deleteAccountAuditConfigurationCore removes the account audit
// configuration.
func (s *IoTService) deleteAccountAuditConfigurationCore(store iotstore.IotStoreInterface) error {
	return store.DeleteGeneric("config/accountAudit")
}

// StartOnDemandAuditTaskInput carries the parsed StartOnDemandAuditTask
// request.
type StartOnDemandAuditTaskInput struct {
	TargetCheckNames []string
}

// startOnDemandAuditTaskCore persists a new on-demand audit task record and
// returns its task id.
func (s *IoTService) startOnDemandAuditTaskCore(store iotstore.IotStoreInterface, in StartOnDemandAuditTaskInput) (string, error) {
	if len(in.TargetCheckNames) == 0 {
		return "", iotstore.ErrMissingParam
	}
	taskId := uuid.New().String()
	// AuditDetails is a map keyed by check name; without a Defender engine
	// every requested check simply carries an empty detail object.
	auditDetails := make(map[string]interface{}, len(in.TargetCheckNames))
	for _, check := range in.TargetCheckNames {
		auditDetails[check] = map[string]interface{}{}
	}
	rec := map[string]interface{}{
		"taskId":           taskId,
		"status":           "IN_PROGRESS",
		"startTime":        time.Now().UTC().Unix(),
		"targetCheckNames": in.TargetCheckNames,
		"auditDetails":     auditDetails,
	}
	if err := store.PutGeneric("auditTask/"+taskId, rec); err != nil {
		return "", err
	}
	return taskId, nil
}

// cancelAuditTaskCore transitions an audit task to CANCELED. An unknown
// task id yields ErrAuditTaskNotFound.
func (s *IoTService) cancelAuditTaskCore(store iotstore.IotStoreInterface, taskId string) error {
	if taskId == "" {
		return iotstore.ErrMissingParam
	}
	key := "auditTask/" + taskId
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrAuditTaskNotFound
	}
	rec["status"] = "CANCELED"
	rec["endTime"] = time.Now().UTC().Unix()
	return store.PutGeneric(key, rec)
}

// AuditTaskDetails is the projected DescribeAuditTask response payload
// (DescribeAuditTaskResponse members only; on-demand tasks always report
// the ON_DEMAND_AUDIT_TASK type).
type AuditTaskDetails struct {
	TaskStatus    interface{}
	TaskType      string
	TaskStartTime interface{}
	AuditDetails  interface{}
}

// describeAuditTaskCore loads an audit task record. An unknown task id
// yields ErrAuditTaskNotFound.
func (s *IoTService) describeAuditTaskCore(store iotstore.IotStoreInterface, taskId string) (*AuditTaskDetails, error) {
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("auditTask/"+taskId, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditTaskNotFound
	}
	return &AuditTaskDetails{
		TaskStatus:    rec["status"],
		TaskType:      "ON_DEMAND_AUDIT_TASK",
		TaskStartTime: rec["startTime"],
		AuditDetails:  rec["auditDetails"],
	}, nil
}

// AuditTaskListItem is one ListAuditTasks entry.
type AuditTaskListItem struct {
	TaskID     interface{}
	TaskStatus interface{}
	TaskType   string
}

// listAuditTasksCore lists every audit task record.
func (s *IoTService) listAuditTasksCore(store iotstore.IotStoreInterface) ([]AuditTaskListItem, error) {
	items, err := store.ListGeneric("auditTask/")
	if err != nil {
		return nil, err
	}
	tasks := make([]AuditTaskListItem, 0, len(items))
	for _, rec := range items {
		tasks = append(tasks, AuditTaskListItem{
			TaskID:     rec["taskId"],
			TaskStatus: rec["status"],
			TaskType:   "ON_DEMAND_AUDIT_TASK",
		})
	}
	return tasks, nil
}

// describeAuditFindingCore loads an audit finding record. No Defender
// engine generates findings, so any caller-supplied id is unknown to the
// platform; AWS returns ResourceNotFoundException.
func (s *IoTService) describeAuditFindingCore(store iotstore.IotStoreInterface, findingId string) (map[string]interface{}, error) {
	if findingId == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("auditFinding/"+findingId, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditFindingNotFound
	}
	return rec, nil
}

// listAuditFindingsCore lists every audit finding record.
func (s *IoTService) listAuditFindingsCore(store iotstore.IotStoreInterface) ([]map[string]interface{}, error) {
	return store.ListGeneric("auditFinding/")
}

// listRelatedResourcesForAuditFindingCore loads the related resources
// recorded on an audit finding. An unknown finding id yields
// ErrAuditFindingNotFound.
func (s *IoTService) listRelatedResourcesForAuditFindingCore(store iotstore.IotStoreInterface, findingID string) ([]map[string]interface{}, error) {
	if findingID == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("auditFinding/"+findingID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditFindingNotFound
	}
	resources := []map[string]interface{}{}
	if raw, ok := rec["relatedResources"].([]interface{}); ok {
		for _, r := range raw {
			if m, ok := r.(map[string]interface{}); ok {
				resources = append(resources, m)
			}
		}
	}
	return resources, nil
}
