package iot

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Command operations (remote command execution metadata management).
// Command definitions store a payload template, mandatory parameters, and
// an optional preprocessor Lambda ARN. Command executions are created by
// the MQTT data-plane when a device triggers a command; the control-plane
// API provides CRUD for definitions and read/delete for execution records.
// ---------------------------------------------------------------------------

func (s *IoTService) CreateCommand(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	cmdID := request.GetParamCaseInsensitive(req.Parameters, "commandId")
	if cmdID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Reject duplicate command IDs.
	cmdExists, err := store.GetGenericExists("iot-command/"+cmdID, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if cmdExists {
		return nil, iotstore.ErrCommandAlreadyExists
	}
	now := time.Now().UTC().Unix()
	rec := map[string]interface{}{
		"commandId":           cmdID,
		"commandArn":          iotstore.BuildCommandARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), cmdID),
		"namespace":           request.GetParamCaseInsensitive(req.Parameters, "namespace"),
		"displayName":         request.GetParamCaseInsensitive(req.Parameters, "displayName"),
		"description":         request.GetParamCaseInsensitive(req.Parameters, "description"),
		"payload":             req.Parameters["payload"],
		"payloadTemplate":     request.GetParamCaseInsensitive(req.Parameters, "payloadTemplate"),
		"preprocessor":        req.Parameters["preprocessor"],
		"mandatoryParameters": req.Parameters["mandatoryParameters"],
		"roleArn":             request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		"createdAt":           now,
		"lastUpdatedAt":       now,
		"deprecated":          false,
		"pendingDeletion":     false,
	}
	if err := store.PutGeneric("iot-command/"+cmdID, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"commandId":  cmdID,
		"commandArn": rec["commandArn"],
	}, nil
}

func (s *IoTService) GetCommand(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	cmdID := request.GetParamCaseInsensitive(req.Parameters, "commandId")
	if cmdID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
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

func (s *IoTService) UpdateCommand(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	cmdID := request.GetParamCaseInsensitive(req.Parameters, "commandId")
	if cmdID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "iot-command/" + cmdID
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCommandNotFound
	}
	if dn := request.GetParamCaseInsensitive(req.Parameters, "displayName"); dn != "" {
		rec["displayName"] = dn
	}
	if desc := request.GetParamCaseInsensitive(req.Parameters, "description"); desc != "" {
		rec["description"] = desc
	}
	if _, ok := req.Parameters["deprecated"]; ok {
		rec["deprecated"] = request.GetBoolParam(req.Parameters, "deprecated")
	}
	rec["lastUpdatedAt"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"commandId":     cmdID,
		"displayName":   rec["displayName"],
		"description":   rec["description"],
		"deprecated":    rec["deprecated"],
		"lastUpdatedAt": rec["lastUpdatedAt"],
	}, nil
}

func (s *IoTService) DeleteCommand(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	cmdID := request.GetParamCaseInsensitive(req.Parameters, "commandId")
	if cmdID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	exists, err := store.GetGenericExists("iot-command/"+cmdID, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCommandNotFound
	}
	if err := store.DeleteGeneric("iot-command/" + cmdID); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"statusCode": 200,
	}, nil
}

func (s *IoTService) ListCommands(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("iot-command/")
	if err != nil {
		return nil, err
	}
	nsFilter := request.GetParamCaseInsensitive(req.Parameters, "namespace")
	commands := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		if nsFilter != "" {
			if ns, _ := rec["namespace"].(string); ns != nsFilter {
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
	return paginatedMaps("commands", commands, req.Parameters), nil
}

// --- CommandExecution ---

func (s *IoTService) GetCommandExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	execID := request.GetParamCaseInsensitive(req.Parameters, "executionId")
	if execID == "" {
		return nil, iotstore.ErrMissingParam
	}
	// Smithy: targetArn is httpQuery+required for both Get and Delete.
	targetArn := request.GetParamCaseInsensitive(req.Parameters, "targetArn")
	if targetArn == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
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

func (s *IoTService) DeleteCommandExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	execID := request.GetParamCaseInsensitive(req.Parameters, "executionId")
	if execID == "" {
		return nil, iotstore.ErrMissingParam
	}
	targetArn := request.GetParamCaseInsensitive(req.Parameters, "targetArn")
	if targetArn == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	exists, err := store.GetGenericExists("iot-command-execution/"+execID, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCommandExecutionNotFound
	}
	if err := store.DeleteGeneric("iot-command-execution/" + execID); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ListCommandExecutions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("iot-command-execution/")
	if err != nil {
		return nil, err
	}
	// Apply body filters: targetArn, commandArn, status, namespace.
	targetFilter := request.GetParamCaseInsensitive(req.Parameters, "targetArn")
	commandFilter := request.GetParamCaseInsensitive(req.Parameters, "commandArn")
	statusFilter := request.GetParamCaseInsensitive(req.Parameters, "status")
	executions := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		if targetFilter != "" {
			if ta, _ := rec["targetArn"].(string); ta != targetFilter {
				continue
			}
		}
		if commandFilter != "" {
			if ca, _ := rec["commandArn"].(string); ca != commandFilter {
				continue
			}
		}
		if statusFilter != "" {
			if st, _ := rec["status"].(string); st != statusFilter {
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
	return paginatedMaps("commandExecutions", executions, req.Parameters), nil
}
