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
// configuration; with deleteScheduledAudits set it also removes every
// scheduled audit record.
func (s *IoTService) deleteAccountAuditConfigurationCore(store iotstore.IotStoreInterface, deleteScheduledAudits bool) error {
	if err := store.DeleteGeneric("config/accountAudit"); err != nil {
		return err
	}
	if !deleteScheduledAudits {
		return nil
	}
	items, err := store.ListGeneric("scheduledAudit/")
	if err != nil {
		return err
	}
	for _, rec := range items {
		name, _ := rec["name"].(string)
		if name == "" {
			continue
		}
		if err := store.DeleteGeneric("scheduledAudit/" + name); err != nil {
			return err
		}
	}
	return nil
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

// ListAuditTasksInput carries the parsed ListAuditTasks request. The model
// marks both time-range members required, and the type/status filters carry
// the AuditTaskType/AuditTaskStatus enums.
type ListAuditTasksInput struct {
	StartTime         int64
	EndTime           int64
	StartTimeProvided bool
	EndTimeProvided   bool
	TaskType          string
	TaskStatus        string
}

// auditTaskTypes is the AuditTaskType enum member set (wire values).
var auditTaskTypes = map[string]bool{
	"ON_DEMAND_AUDIT_TASK": true, "SCHEDULED_AUDIT_TASK": true,
}

// auditTaskStatuses is the AuditTaskStatus enum member set.
var auditTaskStatuses = map[string]bool{
	"IN_PROGRESS": true, "COMPLETED": true, "FAILED": true, "CANCELED": true,
}

// listAuditTasksCore lists audit tasks inside the required [startTime,
// endTime] window, filtered by the optional type/status members. Task
// records carry no type member; they are all on-demand tasks started
// explicitly, so the SCHEDULED_AUDIT_TASK filter matches none of them.
func (s *IoTService) listAuditTasksCore(store iotstore.IotStoreInterface, in ListAuditTasksInput) ([]AuditTaskListItem, error) {
	if !in.StartTimeProvided || !in.EndTimeProvided {
		return nil, iotstore.ErrInvalidRequest
	}
	if in.TaskType != "" && !auditTaskTypes[in.TaskType] {
		return nil, iotstore.ErrInvalidRequest
	}
	if in.TaskStatus != "" && !auditTaskStatuses[in.TaskStatus] {
		return nil, iotstore.ErrInvalidRequest
	}
	items, err := store.ListGeneric("auditTask/")
	if err != nil {
		return nil, err
	}
	tasks := make([]AuditTaskListItem, 0, len(items))
	for _, rec := range items {
		if !auditTaskMatchesFilters(rec, in) {
			continue
		}
		tasks = append(tasks, AuditTaskListItem{
			TaskID:     rec["taskId"],
			TaskStatus: rec["status"],
			TaskType:   "ON_DEMAND_AUDIT_TASK",
		})
	}
	return tasks, nil
}

// auditTaskMatchesFilters applies the time window and the optional
// type/status filters to one stored audit-task record. Records without a
// type member are on-demand tasks.
func auditTaskMatchesFilters(rec map[string]interface{}, in ListAuditTasksInput) bool {
	start := recordEpoch(rec["startTime"])
	if start < in.StartTime || start > in.EndTime {
		return false
	}
	if in.TaskType != "" {
		taskType, _ := rec["taskType"].(string)
		if taskType == "" {
			taskType = "ON_DEMAND_AUDIT_TASK"
		}
		if taskType != in.TaskType {
			return false
		}
	}
	if in.TaskStatus != "" && rec["status"] != in.TaskStatus {
		return false
	}
	return true
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

// ListAuditFindingsInput carries the parsed ListAuditFindings request. The
// documented contract requires exactly one of the taskId or the full
// startTime/endTime range; carrying both together is rejected, and so is
// carrying neither.
type ListAuditFindingsInput struct {
	TaskID             string
	CheckName          string
	ResourceIdentifier map[string]interface{}
	StartTime          int64
	EndTime            int64
	StartTimeProvided  bool
	EndTimeProvided    bool
	ListSuppressed     *bool
}

// listAuditFindingsCore lists audit findings filtered by the optional
// taskId/checkName/resourceIdentifier members, the finding-time range and
// the suppressed-flag selector.
func (s *IoTService) listAuditFindingsCore(store iotstore.IotStoreInterface, in ListAuditFindingsInput) ([]map[string]interface{}, error) {
	// Either the taskId or the startTime and endTime pair must be
	// specified, but not both; the checkName filter is not a substitute
	// for either option and half of the pair does not satisfy it.
	if in.TaskID != "" && (in.StartTimeProvided || in.EndTimeProvided) {
		return nil, iotstore.ErrInvalidRequest
	}
	if in.TaskID == "" && !(in.StartTimeProvided && in.EndTimeProvided) {
		return nil, iotstore.ErrInvalidRequest
	}
	items, err := store.ListGeneric("auditFinding/")
	if err != nil {
		return nil, err
	}
	findings := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		if !auditFindingMatchesFilters(rec, in) {
			continue
		}
		findings = append(findings, rec)
	}
	return findings, nil
}

// auditFindingMatchesFilters applies the ListAuditFindings filter set to one
// stored finding record. A record without an isSuppressed member counts as
// not suppressed; the omitted listSuppressedFindings selector returns both
// kinds.
func auditFindingMatchesFilters(rec map[string]interface{}, in ListAuditFindingsInput) bool {
	if in.TaskID != "" && rec["taskId"] != in.TaskID {
		return false
	}
	if in.CheckName != "" && rec["checkName"] != in.CheckName {
		return false
	}
	if len(in.ResourceIdentifier) > 0 {
		stored, _ := rec["resourceIdentifier"].(map[string]interface{})
		for k, v := range in.ResourceIdentifier {
			if stored[k] != v {
				return false
			}
		}
	}
	if in.StartTimeProvided || in.EndTimeProvided {
		findingTime := recordEpoch(rec["findingTime"])
		if in.StartTimeProvided && findingTime < in.StartTime {
			return false
		}
		if in.EndTimeProvided && findingTime > in.EndTime {
			return false
		}
	}
	if in.ListSuppressed != nil {
		suppressed, _ := rec["isSuppressed"].(bool)
		if suppressed != *in.ListSuppressed {
			return false
		}
	}
	return true
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
