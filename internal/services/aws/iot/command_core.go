package iot

import (
	"time"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Command Core (remote command execution metadata management).
// Command definitions store a payload template, mandatory parameters, and
// an optional preprocessor Lambda ARN. Command executions are created by
// the MQTT data-plane when a device triggers a command; the control-plane
// API provides CRUD for definitions and read/delete for execution records.
// Records live under "iot-command/<commandId>" and
// "iot-command-execution/<executionId>" keys.
// ---------------------------------------------------------------------------

// CreateCommandInput carries the fields for CreateCommand. Payload,
// Preprocessor and MandatoryParameters keep the raw wire values (nested
// structures).
type CreateCommandInput struct {
	CommandID           string
	Namespace           string
	DisplayName         string
	Description         string
	Payload             interface{}
	PayloadTemplate     string
	Preprocessor        interface{}
	MandatoryParameters interface{}
	RoleArn             string
}

// CreateCommandResult is the transport-agnostic result of CreateCommand.
type CreateCommandResult struct {
	CommandID  string
	CommandArn string
}

// UpdateCommandInput carries the fields for UpdateCommand. DeprecatedProvided
// distinguishes an explicitly supplied flag from an omitted one.
type UpdateCommandInput struct {
	CommandID          string
	DisplayName        string
	Description        string
	Deprecated         bool
	DeprecatedProvided bool
}

// UpdateCommandResult is the transport-agnostic result of UpdateCommand.
type UpdateCommandResult struct {
	CommandID     string
	DisplayName   interface{}
	Description   interface{}
	Deprecated    interface{}
	LastUpdatedAt interface{}
}

// ListCommandsInput holds the parameters for listing commands.
type ListCommandsInput struct {
	Namespace string
}

// ListCommandExecutionsInput holds the parameters for listing command
// executions.
type ListCommandExecutionsInput struct {
	TargetArn  string
	CommandArn string
	Status     string
}

// createCommandCore validates and persists a command definition, rejecting
// duplicate command IDs.
func (s *IoTService) createCommandCore(store iotstore.IotStoreInterface, in CreateCommandInput) (*CreateCommandResult, error) {
	if in.CommandID == "" {
		return nil, iotstore.ErrMissingParam
	}
	// Reject duplicate command IDs.
	cmdExists, err := store.GetGenericExists("iot-command/"+in.CommandID, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if cmdExists {
		return nil, iotstore.ErrCommandAlreadyExists
	}
	now := time.Now().UTC().Unix()
	commandArn := iotstore.BuildCommandARN(store.GetAccountID(), store.GetRegion(), in.CommandID)
	rec := map[string]interface{}{
		"commandId":           in.CommandID,
		"commandArn":          commandArn,
		"namespace":           in.Namespace,
		"displayName":         in.DisplayName,
		"description":         in.Description,
		"payload":             in.Payload,
		"payloadTemplate":     in.PayloadTemplate,
		"preprocessor":        in.Preprocessor,
		"mandatoryParameters": in.MandatoryParameters,
		"roleArn":             in.RoleArn,
		"createdAt":           now,
		"lastUpdatedAt":       now,
		"deprecated":          false,
		"pendingDeletion":     false,
	}
	if err := store.PutGeneric("iot-command/"+in.CommandID, rec); err != nil {
		return nil, err
	}
	return &CreateCommandResult{
		CommandID:  in.CommandID,
		CommandArn: commandArn,
	}, nil
}

// getCommandCore retrieves a command definition record.
func (s *IoTService) getCommandCore(store iotstore.IotStoreInterface, cmdID string) (map[string]interface{}, error) {
	if cmdID == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("iot-command/"+cmdID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCommandNotFound
	}
	return rec, nil
}

// updateCommandCore applies the supplied fields to an existing command
// definition.
func (s *IoTService) updateCommandCore(store iotstore.IotStoreInterface, in UpdateCommandInput) (*UpdateCommandResult, error) {
	if in.CommandID == "" {
		return nil, iotstore.ErrMissingParam
	}
	key := "iot-command/" + in.CommandID
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCommandNotFound
	}
	if in.DisplayName != "" {
		rec["displayName"] = in.DisplayName
	}
	if in.Description != "" {
		rec["description"] = in.Description
	}
	if in.DeprecatedProvided {
		rec["deprecated"] = in.Deprecated
	}
	rec["lastUpdatedAt"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return &UpdateCommandResult{
		CommandID:     in.CommandID,
		DisplayName:   rec["displayName"],
		Description:   rec["description"],
		Deprecated:    rec["deprecated"],
		LastUpdatedAt: rec["lastUpdatedAt"],
	}, nil
}

// deleteCommandCore removes a command definition record.
func (s *IoTService) deleteCommandCore(store iotstore.IotStoreInterface, cmdID string) error {
	if cmdID == "" {
		return iotstore.ErrMissingParam
	}
	exists, err := store.GetGenericExists("iot-command/"+cmdID, &map[string]interface{}{})
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrCommandNotFound
	}
	return store.DeleteGeneric("iot-command/" + cmdID)
}

// listCommandsCore lists command summaries with an optional namespace
// filter.
func (s *IoTService) listCommandsCore(store iotstore.IotStoreInterface, in ListCommandsInput) ([]map[string]interface{}, error) {
	items, err := store.ListGeneric("iot-command/")
	if err != nil {
		return nil, err
	}
	commands := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		if in.Namespace != "" {
			if ns, _ := rec["namespace"].(string); ns != in.Namespace {
				continue
			}
		}
		commands = append(commands, map[string]interface{}{
			"commandId":       rec["commandId"],
			"commandArn":      rec["commandArn"],
			"displayName":     rec["displayName"],
			"deprecated":      rec["deprecated"],
			"createdAt":       rec["createdAt"],
			"lastUpdatedAt":   rec["lastUpdatedAt"],
			"pendingDeletion": rec["pendingDeletion"],
		})
	}
	return commands, nil
}

// getCommandExecutionCore retrieves a command execution record. The model
// marks targetArn required alongside executionId; both must be present.
func (s *IoTService) getCommandExecutionCore(store iotstore.IotStoreInterface, execID, targetArn string) (map[string]interface{}, error) {
	if execID == "" || targetArn == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("iot-command-execution/"+execID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCommandExecutionNotFound
	}
	return rec, nil
}

// deleteCommandExecutionCore removes a command execution record. The model
// marks targetArn required alongside executionId; both must be present.
func (s *IoTService) deleteCommandExecutionCore(store iotstore.IotStoreInterface, execID, targetArn string) error {
	if execID == "" || targetArn == "" {
		return iotstore.ErrMissingParam
	}
	exists, err := store.GetGenericExists("iot-command-execution/"+execID, &map[string]interface{}{})
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrCommandExecutionNotFound
	}
	return store.DeleteGeneric("iot-command-execution/" + execID)
}

// listCommandExecutionsCore lists command execution summaries with the
// targetArn / commandArn / status body filters.
func (s *IoTService) listCommandExecutionsCore(store iotstore.IotStoreInterface, in ListCommandExecutionsInput) ([]map[string]interface{}, error) {
	items, err := store.ListGeneric("iot-command-execution/")
	if err != nil {
		return nil, err
	}
	executions := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		if in.TargetArn != "" {
			if ta, _ := rec["targetArn"].(string); ta != in.TargetArn {
				continue
			}
		}
		if in.CommandArn != "" {
			if ca, _ := rec["commandArn"].(string); ca != in.CommandArn {
				continue
			}
		}
		if in.Status != "" {
			if st, _ := rec["status"].(string); st != in.Status {
				continue
			}
		}
		executions = append(executions, map[string]interface{}{
			"executionId": rec["executionId"],
			"commandArn":  rec["commandArn"],
			"targetArn":   rec["targetArn"],
			"status":      rec["status"],
			"createdAt":   rec["createdAt"],
			"startedAt":   rec["startedAt"],
			"completedAt": rec["completedAt"],
		})
	}
	return executions, nil
}
