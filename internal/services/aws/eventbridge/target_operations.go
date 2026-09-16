package eventbridge

import (
	"context"

	"vorpalstacks/internal/common/request"
)

const maxTargetsPerRule = 5

// parseTargetEntries extracts the Targets wire list in its two casing
// variants.
func parseTargetEntries(req *request.ParsedRequest) []interface{} {
	if targets, ok := req.Parameters["Targets"].([]interface{}); ok {
		return targets
	}
	if targets, ok := req.Parameters["targets"].([]interface{}); ok {
		return targets
	}
	return nil
}

// parseTargetIds extracts the Ids wire list in its two casing variants.
func parseTargetIds(req *request.ParsedRequest) []string {
	var targetIDs []string
	if ids, ok := req.Parameters["Ids"].([]interface{}); ok {
		for _, id := range ids {
			if idStr, ok := id.(string); ok {
				targetIDs = append(targetIDs, idStr)
			}
		}
	}
	if ids, ok := req.Parameters["ids"].([]interface{}); ok {
		for _, id := range ids {
			if idStr, ok := id.(string); ok {
				targetIDs = append(targetIDs, idStr)
			}
		}
	}
	return targetIDs
}

// PutTargets adds targets to a rule in EventBridge.
func (s *EventsService) PutTargets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	eventBusName, eventBusNameProvided := eventBusNameParam(req)

	input := PutTargetsInput{
		EventBusName:         eventBusName,
		EventBusNameProvided: eventBusNameProvided,
		Rule:                 request.GetParamLowerFirst(req.Parameters, "Rule"),
		Targets:              parseTargetEntries(req),
		Region:               reqCtx.GetRegion(),
		IAMValidator:         reqCtx.GetIAMValidator(),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.putTargetsCore(ctx, store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"FailedEntryCount": result.FailedEntryCount,
		"FailedEntries":    result.FailedEntries,
	}, nil
}

// RemoveTargets removes targets from a rule in EventBridge.
func (s *EventsService) RemoveTargets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	eventBusName, eventBusNameProvided := eventBusNameParam(req)

	var targetIDs = parseTargetIds(req)

	// Force is currently accepted for SDK parity. The vorpalstacks
	// RemoveTargets implementation always removes the requested target IDs
	// (no extra pre-conditions to bypass), so the flag does not alter
	// behaviour here, but accepting it avoids spurious ValidationException
	// responses for SDK clients that pass Force=true.
	_, _ = req.Parameters["Force"].(bool)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.removeTargetsCore(ctx, store, RemoveTargetsInput{
		EventBusName:         eventBusName,
		EventBusNameProvided: eventBusNameProvided,
		Rule:                 request.GetParamLowerFirst(req.Parameters, "Rule"),
		Ids:                  targetIDs,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"FailedEntryCount": result.FailedEntryCount,
		"FailedEntries":    result.FailedEntries,
	}, nil
}

// ListTargetsByRule lists targets for a rule in EventBridge.
func (s *EventsService) ListTargetsByRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	eventBusName, eventBusNameProvided := eventBusNameParam(req)

	input := ListTargetsByRuleInput{
		EventBusName:         eventBusName,
		EventBusNameProvided: eventBusNameProvided,
		Rule:                 request.GetParamLowerFirst(req.Parameters, "Rule"),
		Limit:                int32(request.GetIntParam(req.Parameters, "Limit")),
		NextToken:            request.GetParamLowerFirst(req.Parameters, "NextToken"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listTargetsByRuleCore(ctx, store, input)
	if err != nil {
		return nil, err
	}

	targets := make([]map[string]interface{}, len(result.Targets))
	for i, t := range result.Targets {
		targets[i] = map[string]interface{}{
			"Id":  t.ID,
			"Arn": t.ARN,
		}
		if t.RoleARN != "" {
			targets[i]["RoleArn"] = t.RoleARN
		}
		if t.Input != "" {
			targets[i]["Input"] = t.Input
		}
		if t.InputPath != "" {
			targets[i]["InputPath"] = t.InputPath
		}
		if t.InputTransformer != nil {
			targets[i]["InputTransformer"] = map[string]interface{}{
				"InputPathsMap": t.InputTransformer.InputPathsMap,
				"InputTemplate": t.InputTransformer.InputTemplate,
			}
		}
		if t.DeadLetterConfig != nil {
			targets[i]["DeadLetterConfig"] = map[string]interface{}{
				"Arn": t.DeadLetterConfig.Arn,
			}
		}
		if t.RetryPolicy != nil {
			targets[i]["RetryPolicy"] = map[string]interface{}{
				"MaximumEventAgeInSeconds": t.RetryPolicy.MaximumEventAgeInSeconds,
				"MaximumRetryAttempts":     t.RetryPolicy.MaximumRetryAttempts,
			}
		}
		if t.SqsParameters != nil {
			targets[i]["SqsParameters"] = map[string]interface{}{
				"MessageGroupId": t.SqsParameters.MessageGroupId,
			}
		}
		if t.HttpParameters != nil {
			targets[i]["HttpParameters"] = map[string]interface{}{
				"HeaderParameters":      t.HttpParameters.HeaderParameters,
				"PathParameterValues":   t.HttpParameters.PathParameterValues,
				"QueryStringParameters": t.HttpParameters.QueryStringParameters,
			}
		}
		if t.KinesisParameters != nil {
			targets[i]["KinesisParameters"] = map[string]interface{}{
				"PartitionKeyPath": t.KinesisParameters.PartitionKeyPath,
			}
		}
		if t.AppSyncParameters != nil {
			targets[i]["AppSyncParameters"] = map[string]interface{}{
				"GraphQLOperation": t.AppSyncParameters.GraphQLOperation,
			}
		}
	}

	response := map[string]interface{}{
		"Targets": targets,
	}

	if result.NextToken != "" {
		response["NextToken"] = result.NextToken
	}

	return response, nil
}
