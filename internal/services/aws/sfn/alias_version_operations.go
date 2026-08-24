package sfn

import (
	"context"
	"encoding/json"
	"fmt"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// PublishStateMachineVersion creates a new version of a state machine by
// snapshotting its current definition; the validation and optimistic
// revision check live in the alias/version Core.
func (s *StepFunctionService) PublishStateMachineVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.publishStateMachineVersionCore(ctx, store, PublishStateMachineVersionInput{
		StateMachineArn: request.GetParamLowerFirst(req.Parameters, "stateMachineArn"),
		Description:     request.GetParamLowerFirst(req.Parameters, "description"),
		RevisionId:      request.GetParamLowerFirst(req.Parameters, "revisionId"),
	})
}

// DeleteStateMachineVersion removes a previously published state machine
// version; a version an alias still routes to is rejected with
// ConflictException.
func (s *StepFunctionService) DeleteStateMachineVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteStateMachineVersionCore(ctx, store, request.GetParamLowerFirst(req.Parameters, "stateMachineVersionArn")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListStateMachineVersions returns a paginated list of all published
// versions for a given state machine.
func (s *StepFunctionService) ListStateMachineVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit, err := parsePageLimit(req)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.listStateMachineVersionsCore(ctx, store,
		request.GetParamLowerFirst(req.Parameters, "stateMachineArn"),
		limit,
		request.GetParamLowerFirst(req.Parameters, "nextToken"))
}

// CreateStateMachineAlias creates an alias that routes traffic to one or
// two versions of a state machine; the state machine is derived from the
// routing configuration's version ARNs.
func (s *StepFunctionService) CreateStateMachineAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	routingConfig, err := parseRoutingConfiguration(req)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.createStateMachineAliasCore(ctx, store, CreateStateMachineAliasInput{
		Name:          request.GetParamLowerFirst(req.Parameters, "name"),
		Description:   request.GetParamLowerFirst(req.Parameters, "description"),
		RoutingConfig: routingConfig,
	})
}

// DescribeStateMachineAlias returns full details of a state machine alias,
// including its routing configuration.
func (s *StepFunctionService) DescribeStateMachineAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.describeStateMachineAliasCore(ctx, store, request.GetParamLowerFirst(req.Parameters, "stateMachineAliasArn"))
}

// UpdateStateMachineAlias modifies the description and/or routing
// configuration of an existing state machine alias.
func (s *StepFunctionService) UpdateStateMachineAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	routingConfig, err := parseRoutingConfiguration(req)
	if err != nil {
		return nil, err
	}
	description := request.GetParamLowerFirst(req.Parameters, "description")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.updateStateMachineAliasCore(ctx, store, UpdateStateMachineAliasInput{
		StateMachineAliasArn: request.GetParamLowerFirst(req.Parameters, "stateMachineAliasArn"),
		Description:          description,
		DescriptionProvided:  description != "",
		RoutingConfig:        routingConfig,
		RoutingProvided:      len(routingConfig) > 0,
	})
}

// DeleteStateMachineAlias removes a state machine alias.
func (s *StepFunctionService) DeleteStateMachineAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteStateMachineAliasCore(ctx, store, request.GetParamLowerFirst(req.Parameters, "stateMachineAliasArn")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListStateMachineAliases returns a paginated list of aliases for a given
// state machine.
func (s *StepFunctionService) ListStateMachineAliases(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit, err := parsePageLimit(req)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.listStateMachineAliasesCore(ctx, store,
		request.GetParamLowerFirst(req.Parameters, "stateMachineArn"),
		limit,
		request.GetParamLowerFirst(req.Parameters, "nextToken"))
}

// parseRoutingConfiguration extracts the routing configuration from the request
// parameters. Handles both JSON string (URL-encoded) and native []interface{}
// (JSON-RPC) formats. Weight-sum and per-entry validation run in the Core.
func parseRoutingConfiguration(req *request.ParsedRequest) ([]sfnstore.RoutingConfiguration, error) {
	rawValue := req.Parameters["routingConfiguration"]
	if rawValue == nil {
		return nil, nil
	}

	var rawConfig []map[string]interface{}

	switch v := rawValue.(type) {
	case string:
		if v == "" {
			return nil, nil
		}
		if err := json.Unmarshal([]byte(v), &rawConfig); err != nil {
			return nil, NewInvalidDefinitionException("routingConfiguration is not valid JSON")
		}
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				rawConfig = append(rawConfig, m)
			}
		}
	default:
		return nil, NewInvalidDefinitionException("routingConfiguration must be a JSON array, got " + fmt.Sprintf("%T", rawValue))
	}

	config := make([]sfnstore.RoutingConfiguration, 0, len(rawConfig))
	for _, item := range rawConfig {
		versionArn, _ := item["stateMachineVersionArn"].(string)
		weightVal, _ := item["weight"].(float64)

		if versionArn == "" {
			return nil, NewInvalidDefinitionException("each routingConfiguration entry must have stateMachineVersionArn")
		}

		config = append(config, sfnstore.RoutingConfiguration{
			StateMachineVersionArn: versionArn,
			Weight:                 int32(weightVal),
		})
	}

	return config, nil
}

// formatRoutingConfiguration converts routing configuration to the response map
// format expected by the SDK.
func formatRoutingConfiguration(config []sfnstore.RoutingConfiguration) []map[string]interface{} {
	result := make([]map[string]interface{}, len(config))
	for i, rc := range config {
		result[i] = map[string]interface{}{
			"stateMachineVersionArn": rc.StateMachineVersionArn,
			"weight":                 rc.Weight,
		}
	}
	return result
}
