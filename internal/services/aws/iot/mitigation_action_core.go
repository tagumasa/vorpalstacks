package iot

import (
	"time"

	"github.com/google/uuid"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Mitigation Action Core (audit finding mitigation actions and ML Detect
// mitigation tasks). Mitigation action definitions live under the
// generic-KV "mitigationAction/" prefix; Detect task records under
// "detectMitigationTask/<taskId>". The task records persist so that
// Cancel/Describe resolve the identifier and return
// ResourceNotFoundException for unknown task ids, matching the Smithy error
// trait set on each operation.
// ---------------------------------------------------------------------------

// CreateMitigationActionInput carries the parsed CreateMitigationAction
// request.
type CreateMitigationActionInput struct {
	ActionName   string
	RoleArn      string
	ActionParams map[string]interface{}
	Tags         map[string]string
}

// CreateMitigationActionResult is the transport-agnostic result of
// CreateMitigationAction.
type CreateMitigationActionResult struct {
	ActionArn string
	ActionID  string
}

// createMitigationActionCore validates and persists a mitigation action.
// The roleArn and actionParams members are required by the model; the
// actionId is minted once at create and stored so Describe/Update echo the
// same identifier.
func (s *IoTService) createMitigationActionCore(store iotstore.IotStoreInterface, in CreateMitigationActionInput) (*CreateMitigationActionResult, error) {
	if in.RoleArn == "" {
		return nil, iotstore.ErrMissingParam
	}
	if in.ActionParams == nil {
		return nil, iotstore.ErrMissingParam
	}
	actionID := uuid.New().String()
	rec, err := s.bulkCreateCore(store, "mitigationAction", in.ActionName, map[string]interface{}{
		"actionType":   deriveMitigationActionType(in.ActionParams),
		"actionParams": in.ActionParams,
		"roleArn":      in.RoleArn,
		"actionId":     actionID,
	})
	if err != nil {
		return nil, err
	}
	arn := iotstore.BuildMitigationActionARN(store.GetAccountID(), store.GetRegion(), bulkName(rec))
	if len(in.Tags) > 0 {
		if err := store.TagResource(arn, in.Tags); err != nil {
			return nil, err
		}
	}
	return &CreateMitigationActionResult{
		ActionArn: arn,
		ActionID:  actionID,
	}, nil
}

// deleteMitigationActionCore removes a mitigation action record and its
// tags.
func (s *IoTService) deleteMitigationActionCore(store iotstore.IotStoreInterface, name string) error {
	arn := iotstore.BuildMitigationActionARN(store.GetAccountID(), store.GetRegion(), name)
	_ = store.DeleteAllTags(arn)
	return s.bulkDeleteCore(store, "mitigationAction", name)
}

// deriveMitigationActionType infers the action type from the actionParams
// keys. AWS derives the type automatically based on which params member is
// set; the keys and returned values are the MitigationActionParams members
// and MitigationActionType enum members.
func deriveMitigationActionType(params map[string]interface{}) string {
	if _, ok := params["updateDeviceCertificateParams"]; ok {
		return "UPDATE_DEVICE_CERTIFICATE"
	}
	if _, ok := params["updateCACertificateParams"]; ok {
		return "UPDATE_CA_CERTIFICATE"
	}
	if _, ok := params["addThingsToThingGroupParams"]; ok {
		return "ADD_THINGS_TO_THING_GROUP"
	}
	if _, ok := params["replaceDefaultPolicyVersionParams"]; ok {
		return "REPLACE_DEFAULT_POLICY_VERSION"
	}
	if _, ok := params["enableIoTLoggingParams"]; ok {
		return "ENABLE_IOT_LOGGING"
	}
	if _, ok := params["publishFindingToSnsParams"]; ok {
		return "PUBLISH_FINDING_TO_SNS"
	}
	return ""
}

// MitigationActionRecord is the persisted mitigation action record plus its
// ARN.
type MitigationActionRecord struct {
	Rec map[string]interface{}
	Arn string
}

// describeMitigationActionCore loads a mitigation action record and computes
// its ARN. An unknown name yields ErrMitigationActionNotFound.
func (s *IoTService) describeMitigationActionCore(store iotstore.IotStoreInterface, name string) (*MitigationActionRecord, error) {
	rec, exists, err := s.bulkGetCore(store, "mitigationAction", name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrMitigationActionNotFound
	}
	return &MitigationActionRecord{
		Rec: rec,
		Arn: iotstore.BuildMitigationActionARN(store.GetAccountID(), store.GetRegion(), name),
	}, nil
}

// MitigationActionListItem is one ListMitigationActions entry.
type MitigationActionListItem struct {
	Name         string
	Arn          string
	CreationDate interface{}
}

// listMitigationActionsCore lists every mitigation action with its ARN.
func (s *IoTService) listMitigationActionsCore(store iotstore.IotStoreInterface) ([]MitigationActionListItem, error) {
	items, err := s.bulkListCore(store, "mitigationAction")
	if err != nil {
		return nil, err
	}
	out := make([]MitigationActionListItem, 0, len(items))
	for _, item := range items {
		name, _ := item["name"].(string)
		out = append(out, MitigationActionListItem{
			Name:         name,
			Arn:          iotstore.BuildMitigationActionARN(store.GetAccountID(), store.GetRegion(), name),
			CreationDate: item["creationDate"],
		})
	}
	return out, nil
}

// UpdateMitigationActionInput carries the parsed UpdateMitigationAction
// request. RoleArnProvided distinguishes an explicitly supplied roleArn
// from an omitted one so a partial update keeps the stored role.
type UpdateMitigationActionInput struct {
	ActionName      string
	RoleArn         string
	RoleArnProvided bool
	ActionParams    map[string]interface{}
}

// UpdateMitigationActionResult is the transport-agnostic result of
// UpdateMitigationAction.
type UpdateMitigationActionResult struct {
	Rec map[string]interface{}
	Arn string
}

// updateMitigationActionCore merges the supplied fields into an existing
// mitigation action. An unknown name yields
// ErrMitigationActionNotFound.
func (s *IoTService) updateMitigationActionCore(store iotstore.IotStoreInterface, in UpdateMitigationActionInput) (*UpdateMitigationActionResult, error) {
	merge := map[string]interface{}{
		"actionParams": in.ActionParams,
	}
	if in.RoleArnProvided {
		merge["roleArn"] = in.RoleArn
	}
	rec, exists, err := s.bulkUpdateCore(store, "mitigationAction", in.ActionName, merge)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrMitigationActionNotFound
	}
	return &UpdateMitigationActionResult{
		Rec: rec,
		Arn: iotstore.BuildMitigationActionARN(store.GetAccountID(), store.GetRegion(), bulkName(rec)),
	}, nil
}

// StartDetectMitigationActionsTaskInput carries the parsed
// StartDetectMitigationActionsTask request. Target keeps the raw wire
// structure; Actions is the mitigation action name list.
type StartDetectMitigationActionsTaskInput struct {
	TaskID  string
	Target  map[string]interface{}
	Actions []string
}

// startDetectMitigationActionsTaskCore persists a Detect mitigation task
// record. The taskId, target and actions members are required by the model.
func (s *IoTService) startDetectMitigationActionsTaskCore(store iotstore.IotStoreInterface, in StartDetectMitigationActionsTaskInput) (string, error) {
	if in.TaskID == "" {
		return "", iotstore.ErrMissingParam
	}
	if in.Target == nil {
		return "", iotstore.ErrMissingParam
	}
	if len(in.Actions) == 0 {
		return "", iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{
		"taskId":    in.TaskID,
		"status":    "IN_PROGRESS",
		"startTime": time.Now().UTC().Unix(),
		"target":    in.Target,
		"actions":   in.Actions,
	}
	if err := store.PutGeneric("detectMitigationTask/"+in.TaskID, rec); err != nil {
		return "", err
	}
	return in.TaskID, nil
}

// cancelDetectMitigationActionsTaskCore transitions a Detect mitigation
// task to CANCELED. An unknown task id yields
// ErrDetectMitigationTaskNotFound.
func (s *IoTService) cancelDetectMitigationActionsTaskCore(store iotstore.IotStoreInterface, taskId string) error {
	if taskId == "" {
		return iotstore.ErrMissingParam
	}
	key := "detectMitigationTask/" + taskId
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrDetectMitigationTaskNotFound
	}
	rec["status"] = "CANCELED"
	rec["endTime"] = time.Now().UTC().Unix()
	return store.PutGeneric(key, rec)
}

// DetectMitigationTaskSummary is the projected
// DetectMitigationActionsTaskSummary member set carried by the describe
// response's taskSummary member.
type DetectMitigationTaskSummary struct {
	TaskID        interface{}
	TaskStatus    interface{}
	TaskStartTime interface{}
	TaskEndTime   interface{}
	Target        interface{}
}

// describeDetectMitigationActionsTaskCore loads a Detect mitigation task
// record projected onto the summary member set. An unknown task id yields
// ErrDetectMitigationTaskNotFound.
func (s *IoTService) describeDetectMitigationActionsTaskCore(store iotstore.IotStoreInterface, taskId string) (*DetectMitigationTaskSummary, error) {
	if taskId == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("detectMitigationTask/"+taskId, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrDetectMitigationTaskNotFound
	}
	return &DetectMitigationTaskSummary{
		TaskID:        rec["taskId"],
		TaskStatus:    rec["status"],
		TaskStartTime: rec["startTime"],
		TaskEndTime:   rec["endTime"],
		Target:        rec["target"],
	}, nil
}

// DetectMitigationTaskListItem is one ListDetectMitigationActionsTasks
// entry.
type DetectMitigationTaskListItem struct {
	TaskID        interface{}
	TaskStatus    interface{}
	TaskStartTime interface{}
}

// listDetectMitigationActionsTasksCore lists every Detect mitigation task.
func (s *IoTService) listDetectMitigationActionsTasksCore(store iotstore.IotStoreInterface) ([]DetectMitigationTaskListItem, error) {
	items, err := store.ListGeneric("detectMitigationTask/")
	if err != nil {
		return nil, err
	}
	tasks := make([]DetectMitigationTaskListItem, 0, len(items))
	for _, rec := range items {
		tasks = append(tasks, DetectMitigationTaskListItem{
			TaskID:        rec["taskId"],
			TaskStatus:    rec["status"],
			TaskStartTime: rec["startTime"],
		})
	}
	return tasks, nil
}
