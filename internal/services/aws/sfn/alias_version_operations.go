package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// PublishStateMachineVersion creates a new version of a state machine by snapshotting
// its current definition. Each successive call increments the version number.
func (s *StepFunctionService) PublishStateMachineVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	smArn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")
	description := request.GetParamLowerFirst(req.Parameters, "description")

	if smArn == "" {
		return nil, NewInvalidArnException("stateMachineArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetStateMachine(ctx, smArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + smArn)
		}
		return nil, err
	}

	version, err := store.PublishStateMachineVersion(ctx, smArn, description)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"stateMachineVersionArn": version.StateMachineVersionArn,
		"creationDate":           version.CreationDate.Unix(),
	}, nil
}

// DescribeStateMachineVersion returns details about a specific state machine version.
func (s *StepFunctionService) DescribeStateMachineVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	versionArn := request.GetParamLowerFirst(req.Parameters, "stateMachineVersionArn")

	if versionArn == "" {
		return nil, NewInvalidArnException("stateMachineVersionArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	version, err := store.GetStateMachineVersion(ctx, versionArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineVersionNotFound) {
			return nil, awserrors.NewAWSError("StateMachineVersionNotFound", "State Machine Version Does not exist: "+versionArn, 400)
		}
		return nil, err
	}

	return map[string]interface{}{
		"stateMachineVersionArn": version.StateMachineVersionArn,
		"stateMachineArn":        version.StateMachineArn,
		"version":                version.Version,
		"creationDate":           version.CreationDate.Unix(),
		"description":            version.Description,
	}, nil
}

// DeleteStateMachineVersion removes a previously published state machine version.
func (s *StepFunctionService) DeleteStateMachineVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	versionArn := request.GetParamLowerFirst(req.Parameters, "stateMachineVersionArn")

	if versionArn == "" {
		return nil, NewInvalidArnException("stateMachineVersionArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteStateMachineVersion(ctx, versionArn); err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineVersionNotFound) {
			return nil, awserrors.NewAWSError("StateMachineVersionNotFound", "State Machine Version Does not exist: "+versionArn, 400)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListStateMachineVersions returns a paginated list of all published versions
// for a given state machine, ordered by version number descending.
func (s *StepFunctionService) ListStateMachineVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	smArn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")
	limit := int32(request.GetIntParam(req.Parameters, "maxResults"))
	if limit == 0 {
		limit = 100
	}
	nextToken := request.GetParamLowerFirst(req.Parameters, "nextToken")

	if smArn == "" {
		return nil, NewInvalidArnException("stateMachineArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.ListStateMachineVersions(ctx, smArn, limit, nextToken)
	if err != nil {
		return nil, err
	}

	versions := make([]map[string]interface{}, len(result.Versions))
	for i, v := range result.Versions {
		versions[i] = map[string]interface{}{
			"stateMachineVersionArn": v.StateMachineVersionArn,
			"creationDate":           v.CreationDate.Unix(),
		}
	}

	resp := map[string]interface{}{
		"stateMachineVersions": versions,
	}
	if result.NextToken != "" {
		resp["nextToken"] = result.NextToken
	}

	return resp, nil
}

// CreateStateMachineAlias creates an alias that routes traffic to one or more
// state machine versions. The state machine ARN is derived from the version ARN
// in the routing configuration when not provided directly.
func (s *StepFunctionService) CreateStateMachineAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "name")
	description := request.GetParamLowerFirst(req.Parameters, "description")

	if name == "" {
		return nil, NewInvalidName("name is required")
	}

	routingConfig, err := parseRoutingConfiguration(req)
	if err != nil {
		return nil, err
	}
	if len(routingConfig) == 0 {
		return nil, NewInvalidDefinitionException("routingConfiguration is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	smArn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")
	if smArn == "" {
		smArn = extractStateMachineArnFromVersionArn(routingConfig[0].StateMachineVersionArn)
	}
	if smArn == "" {
		return nil, NewInvalidArnException("stateMachineArn is required")
	}

	_, err = store.GetStateMachine(ctx, smArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + smArn)
		}
		return nil, err
	}

	alias := &sfnstore.StateMachineAlias{
		StateMachineArn:      smArn,
		Name:                 name,
		Description:          description,
		RoutingConfiguration: routingConfig,
	}

	if err := store.CreateStateMachineAlias(ctx, alias); err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineAliasAlreadyExists) {
			return nil, awserrors.NewAWSError("StateMachineAliasAlreadyExists", fmt.Sprintf("State Machine Alias already exists: %s", name), 400)
		}
		return nil, err
	}

	return map[string]interface{}{
		"stateMachineAliasArn": alias.StateMachineAliasArn,
		"creationDate":         alias.CreationDate.Unix(),
	}, nil
}

// DescribeStateMachineAlias returns full details of a state machine alias,
// including its routing configuration.
func (s *StepFunctionService) DescribeStateMachineAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	aliasArn := request.GetParamLowerFirst(req.Parameters, "stateMachineAliasArn")

	if aliasArn == "" {
		return nil, NewInvalidArnException("stateMachineAliasArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	alias, err := store.GetStateMachineAlias(ctx, aliasArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineAliasNotFound) {
			return nil, awserrors.NewAWSError("StateMachineAliasDoesNotExist", "State Machine Alias Does not exist: "+aliasArn, 400)
		}
		return nil, err
	}

	resp := map[string]interface{}{
		"stateMachineAliasArn": alias.StateMachineAliasArn,
		"name":                 alias.Name,
		"creationDate":         alias.CreationDate.Unix(),
		"updateDate":           alias.UpdateDate.Unix(),
	}

	if alias.Description != "" {
		resp["description"] = alias.Description
	}
	if len(alias.RoutingConfiguration) > 0 {
		resp["routingConfiguration"] = formatRoutingConfiguration(alias.RoutingConfiguration)
	}

	return resp, nil
}

// UpdateStateMachineAlias modifies the description and/or routing configuration
// of an existing state machine alias.
func (s *StepFunctionService) UpdateStateMachineAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	aliasArn := request.GetParamLowerFirst(req.Parameters, "stateMachineAliasArn")
	description := request.GetParamLowerFirst(req.Parameters, "description")

	if aliasArn == "" {
		return nil, NewInvalidArnException("stateMachineAliasArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	alias, err := store.GetStateMachineAlias(ctx, aliasArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineAliasNotFound) {
			return nil, awserrors.NewAWSError("StateMachineAliasDoesNotExist", "State Machine Alias Does not exist: "+aliasArn, 400)
		}
		return nil, err
	}

	routingConfig, err := parseRoutingConfiguration(req)
	if err != nil {
		return nil, err
	}

	if description != "" {
		alias.Description = description
	}
	if len(routingConfig) > 0 {
		alias.RoutingConfiguration = routingConfig
	}

	if err := store.UpdateStateMachineAlias(ctx, alias); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"updateDate": alias.UpdateDate.Unix(),
	}, nil
}

// DeleteStateMachineAlias removes a state machine alias.
func (s *StepFunctionService) DeleteStateMachineAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	aliasArn := request.GetParamLowerFirst(req.Parameters, "stateMachineAliasArn")

	if aliasArn == "" {
		return nil, NewInvalidArnException("stateMachineAliasArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteStateMachineAlias(ctx, aliasArn); err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineAliasNotFound) {
			return nil, awserrors.NewAWSError("StateMachineAliasDoesNotExist", "State Machine Alias Does not exist: "+aliasArn, 400)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListStateMachineAliases returns a paginated list of aliases for a given
// state machine.
func (s *StepFunctionService) ListStateMachineAliases(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	smArn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")
	limit := int32(request.GetIntParam(req.Parameters, "maxResults"))
	if limit == 0 {
		limit = 100
	}
	nextToken := request.GetParamLowerFirst(req.Parameters, "nextToken")

	if smArn == "" {
		return nil, NewInvalidArnException("stateMachineArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.ListStateMachineAliases(ctx, smArn, limit, nextToken)
	if err != nil {
		return nil, err
	}

	aliases := make([]map[string]interface{}, len(result.Aliases))
	for i, a := range result.Aliases {
		aliases[i] = map[string]interface{}{
			"stateMachineAliasArn": a.StateMachineAliasArn,
			"creationDate":         a.CreationDate.Unix(),
		}
	}

	resp := map[string]interface{}{
		"stateMachineAliases": aliases,
	}
	if result.NextToken != "" {
		resp["nextToken"] = result.NextToken
	}

	return resp, nil
}

// UpdateMapRun modifies the concurrency and failure tolerance settings of a
// running map execution.
func (s *StepFunctionService) UpdateMapRun(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	mapRunArn := request.GetParamLowerFirst(req.Parameters, "mapRunArn")

	if mapRunArn == "" {
		return nil, NewInvalidName("mapRunArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	mr, err := store.GetMapRun(ctx, mapRunArn)
	if err != nil {
		return nil, NewMapRunDoesNotExist("Map Run does not exist: " + mapRunArn)
	}

	if maxConcurrency := request.GetIntParam(req.Parameters, "maxConcurrency"); maxConcurrency > 0 {
		mr.MaxConcurrency = int64(maxConcurrency)
	}
	if toleratedCount := request.GetInt64Param(req.Parameters, "toleratedFailureCount"); toleratedCount > 0 {
		mr.ToleratedFailureCount = toleratedCount
	}
	if toleratedPct := float32(request.GetFloatParam(req.Parameters, "toleratedFailurePercentage")); toleratedPct > 0 {
		mr.ToleratedFailurePercentage = toleratedPct
	}

	if err := store.UpdateMapRun(ctx, mr); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// parseRoutingConfiguration extracts the routing configuration from the request
// parameters. Handles both JSON string (URL-encoded) and native []interface{}
// (JSON-RPC) formats.
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
		return nil, nil
	}

	var totalWeight int32
	config := make([]sfnstore.RoutingConfiguration, 0, len(rawConfig))
	for _, item := range rawConfig {
		versionArn, _ := item["stateMachineVersionArn"].(string)
		weightVal, _ := item["weight"].(float64)
		weight := int32(weightVal)

		if versionArn == "" {
			return nil, NewInvalidDefinitionException("each routingConfiguration entry must have stateMachineVersionArn")
		}

		config = append(config, sfnstore.RoutingConfiguration{
			StateMachineVersionArn: versionArn,
			Weight:                 weight,
		})
		totalWeight += weight
	}

	if totalWeight != 100 && len(config) > 0 {
		return nil, NewInvalidDefinitionException(fmt.Sprintf("routingConfiguration weights must sum to 100, got %d", totalWeight))
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

// extractStateMachineArnFromVersionArn derives the state machine ARN from a
// version ARN by stripping the trailing version number segment.
func extractStateMachineArnFromVersionArn(versionArn string) string {
	idx := strings.LastIndex(versionArn, ":")
	if idx < 0 {
		return ""
	}
	afterColon := versionArn[idx+1:]
	if _, err := strconv.Atoi(afterColon); err == nil {
		return versionArn[:idx]
	}
	return versionArn
}
