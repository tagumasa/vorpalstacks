package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// ---------------------------------------------------------------------------
// Command operations (remote command execution metadata management).
// Command definitions store a payload template, mandatory parameters, and
// an optional preprocessor Lambda ARN. Command executions are created by
// the MQTT data-plane when a device triggers a command; the control-plane
// API provides CRUD for definitions and read/delete for execution records.
// ---------------------------------------------------------------------------

func (s *IoTService) CreateCommand(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.createCommandCore(store, CreateCommandInput{
		CommandID:           request.GetParamCaseInsensitive(req.Parameters, "commandId"),
		Namespace:           request.GetParamCaseInsensitive(req.Parameters, "namespace"),
		DisplayName:         request.GetParamCaseInsensitive(req.Parameters, "displayName"),
		Description:         request.GetParamCaseInsensitive(req.Parameters, "description"),
		Payload:             req.Parameters["payload"],
		PayloadTemplate:     request.GetParamCaseInsensitive(req.Parameters, "payloadTemplate"),
		Preprocessor:        req.Parameters["preprocessor"],
		MandatoryParameters: req.Parameters["mandatoryParameters"],
		RoleArn:             request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"commandId":  result.CommandID,
		"commandArn": result.CommandArn,
	}, nil
}

func (s *IoTService) GetCommand(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getCommandCore(store, request.GetParamCaseInsensitive(req.Parameters, "commandId"))
}

func (s *IoTService) UpdateCommand(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	_, deprecatedProvided := req.Parameters["deprecated"]
	result, err := s.updateCommandCore(store, UpdateCommandInput{
		CommandID:          request.GetParamCaseInsensitive(req.Parameters, "commandId"),
		DisplayName:        request.GetParamCaseInsensitive(req.Parameters, "displayName"),
		Description:        request.GetParamCaseInsensitive(req.Parameters, "description"),
		Deprecated:         request.GetBoolParam(req.Parameters, "deprecated"),
		DeprecatedProvided: deprecatedProvided,
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"commandId":     result.CommandID,
		"displayName":   result.DisplayName,
		"description":   result.Description,
		"deprecated":    result.Deprecated,
		"lastUpdatedAt": result.LastUpdatedAt,
	}, nil
}

func (s *IoTService) DeleteCommand(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteCommandCore(store, request.GetParamCaseInsensitive(req.Parameters, "commandId")); err != nil {
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
	commands, err := s.listCommandsCore(store, ListCommandsInput{
		Namespace: request.GetParamCaseInsensitive(req.Parameters, "namespace"),
	})
	if err != nil {
		return nil, err
	}
	return paginatedMaps("commands", commands, req.Parameters)
}

// --- CommandExecution ---

func (s *IoTService) GetCommandExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Smithy: targetArn is httpQuery+required for both Get and Delete.
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getCommandExecutionCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "executionId"),
		request.GetParamCaseInsensitive(req.Parameters, "targetArn"))
}

func (s *IoTService) DeleteCommandExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteCommandExecutionCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "executionId"),
		request.GetParamCaseInsensitive(req.Parameters, "targetArn")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ListCommandExecutions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	executions, err := s.listCommandExecutionsCore(store, ListCommandExecutionsInput{
		TargetArn:  request.GetParamCaseInsensitive(req.Parameters, "targetArn"),
		CommandArn: request.GetParamCaseInsensitive(req.Parameters, "commandArn"),
		Status:     request.GetParamCaseInsensitive(req.Parameters, "status"),
	})
	if err != nil {
		return nil, err
	}
	return paginatedMaps("commandExecutions", executions, req.Parameters)
}
