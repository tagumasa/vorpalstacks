package eventbridge

import (
	"context"

	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

func isValidTargetARN(arn string) bool {
	if arn == "" {
		return false
	}
	_, service, _, _, _ := svcarn.SplitARN(arn)
	validServices := map[string]bool{
		"lambda":        true,
		"sqs":           true,
		"sns":           true,
		"events":        true,
		"kinesis":       true,
		"firehose":      true,
		"stepfunctions": true,
		"states":        true,
		"ecs":           true,
	}
	return validServices[service]
}

// PutTargets adds targets to a rule in EventBridge.
func (s *EventsService) PutTargets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	ruleName := request.GetParamLowerFirst(req.Parameters, "Rule")
	if ruleName == "" {
		return nil, NewValidationException("Rule name is required")
	}

	eventBusName := request.GetParamLowerFirst(req.Parameters, "EventBusName")
	if eventBusName == "" {
		eventBusName = "default"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Check if event bus exists
	if _, err := store.GetEventBus(ctx, eventBusName); err != nil {
		if err == eventsstore.ErrEventBusNotFound {
			return nil, NewResourceNotFoundException("Event bus '" + eventBusName + "' does not exist")
		}
		return nil, err
	}

	_, err = store.GetRule(ctx, eventBusName, ruleName)
	if err != nil {
		if err == eventsstore.ErrRuleNotFound {
			return nil, NewResourceNotFoundException("Rule '" + ruleName + "' does not exist on event bus '" + eventBusName + "'")
		}
		return nil, err
	}

	targets, ok := req.Parameters["Targets"].([]interface{})
	if !ok {
		targets, ok = req.Parameters["targets"].([]interface{})
	}
	if !ok || len(targets) == 0 {
		return nil, NewValidationException("Targets are required")
	}

	// Check for duplicate target IDs
	seenIDs := make(map[string]bool)
	for _, t := range targets {
		targetMap, ok := t.(map[string]interface{})
		if !ok {
			continue
		}
		targetID, _ := targetMap["Id"].(string)
		if targetID != "" && seenIDs[targetID] {
			return nil, NewValidationException("Duplicate target ID: " + targetID)
		}
		seenIDs[targetID] = true
	}

	failedEntries := make([]map[string]interface{}, 0)
	failedCount := int32(0)

	for _, t := range targets {
		targetMap, ok := t.(map[string]interface{})
		if !ok {
			continue
		}

		targetID, _ := targetMap["Id"].(string)
		targetArn, _ := targetMap["Arn"].(string)

		if targetID == "" || targetArn == "" {
			failedEntries = append(failedEntries, map[string]interface{}{
				"TargetId":     targetID,
				"ErrorCode":    "ValidationException",
				"ErrorMessage": "Target ID and ARN are required",
			})
			failedCount++
			continue
		}

		if !isValidTargetARN(targetArn) {
			failedEntries = append(failedEntries, map[string]interface{}{
				"TargetId":     targetID,
				"ErrorCode":    "ValidationException",
				"ErrorMessage": "Invalid target ARN",
			})
			failedCount++
			continue
		}

		target := &eventsstore.Target{
			ID:           targetID,
			RuleName:     ruleName,
			EventBusName: eventBusName,
			ARN:          targetArn,
		}

		if input, ok := targetMap["Input"].(string); ok {
			target.Input = input
		}

		if inputPath, ok := targetMap["InputPath"].(string); ok {
			target.InputPath = inputPath
		}

		if roleArn, ok := targetMap["RoleArn"].(string); ok {
			if roleArn != "" {
				if s.bus != nil {
					if rr := s.bus.RoleResolver(); rr != nil {
						if err := rr.ValidateRole(ctx, roleArn); err != nil {
							failedEntries = append(failedEntries, map[string]interface{}{
								"TargetId":     targetID,
								"ErrorCode":    "ValidationException",
								"ErrorMessage": err.Error(),
							})
							failedCount++
							continue
						}
					}
				}
				if validator := reqCtx.GetIAMValidator(); validator != nil {
					if err := validator.ValidateRoleForService(ctx, roleArn, iam.ServicePrincipalEvents); err != nil {
						failedEntries = append(failedEntries, map[string]interface{}{
							"TargetId":     targetID,
							"ErrorCode":    "ValidationException",
							"ErrorMessage": err.Error(),
						})
						failedCount++
						continue
					}
				}
			}
			target.RoleARN = roleArn
		}

		if inputTransformer, ok := targetMap["InputTransformer"].(map[string]interface{}); ok {
			target.InputTransformer = &eventsstore.InputTransformer{}
			if paths, ok := inputTransformer["InputPathsMap"].(map[string]interface{}); ok {
				target.InputTransformer.InputPathsMap = make(map[string]string)
				for k, v := range paths {
					if vs, ok := v.(string); ok {
						target.InputTransformer.InputPathsMap[k] = vs
					}
				}
			}
			if template, ok := inputTransformer["InputTemplate"].(string); ok {
				target.InputTransformer.InputTemplate = template
			}
		}

		if dlConfig, ok := targetMap["DeadLetterConfig"].(map[string]interface{}); ok {
			target.DeadLetterConfig = &eventsstore.DeadLetterConfig{}
			if arn, ok := dlConfig["Arn"].(string); ok {
				target.DeadLetterConfig.Arn = arn
			}
		}

		if retryPolicy, ok := targetMap["RetryPolicy"].(map[string]interface{}); ok {
			target.RetryPolicy = &eventsstore.RetryPolicy{}
			if maxAge, ok := retryPolicy["MaximumEventAgeInSeconds"].(float64); ok {
				target.RetryPolicy.MaximumEventAgeInSeconds = int32(maxAge)
			}
			if maxRetry, ok := retryPolicy["MaximumRetryAttempts"].(float64); ok {
				target.RetryPolicy.MaximumRetryAttempts = int32(maxRetry)
			}
		}

		if sqsParams, ok := targetMap["SqsParameters"].(map[string]interface{}); ok {
			target.SqsParameters = &eventsstore.SqsParameters{}
			if groupId, ok := sqsParams["MessageGroupId"].(string); ok {
				target.SqsParameters.MessageGroupId = groupId
			}
		}

		if httpParams, ok := targetMap["HttpParameters"].(map[string]interface{}); ok {
			target.HttpParameters = &eventsstore.HttpParameters{}
			if headers, ok := httpParams["HeaderParameters"].(map[string]interface{}); ok {
				target.HttpParameters.HeaderParameters = make(map[string]string)
				for k, v := range headers {
					if vs, ok := v.(string); ok {
						target.HttpParameters.HeaderParameters[k] = vs
					}
				}
			}
			if paths, ok := httpParams["PathParameterValues"].([]interface{}); ok {
				for _, p := range paths {
					if ps, ok := p.(string); ok {
						target.HttpParameters.PathParameterValues = append(target.HttpParameters.PathParameterValues, ps)
					}
				}
			}
			if qs, ok := httpParams["QueryStringParameters"].(map[string]interface{}); ok {
				target.HttpParameters.QueryStringParameters = make(map[string]string)
				for k, v := range qs {
					if vs, ok := v.(string); ok {
						target.HttpParameters.QueryStringParameters[k] = vs
					}
				}
			}
		}

		if err := store.PutTarget(ctx, target); err != nil {
			failedEntries = append(failedEntries, map[string]interface{}{
				"TargetId":     targetID,
				"ErrorCode":    "InternalFailure",
				"ErrorMessage": err.Error(),
			})
			failedCount++
		}
	}

	return map[string]interface{}{
		"FailedEntryCount": failedCount,
		"FailedEntries":    failedEntries,
	}, nil
}

// RemoveTargets removes targets from a rule in EventBridge.
func (s *EventsService) RemoveTargets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	ruleName := request.GetParamLowerFirst(req.Parameters, "Rule")
	if ruleName == "" {
		return nil, NewValidationException("Rule name is required")
	}

	eventBusName := request.GetParamLowerFirst(req.Parameters, "EventBusName")
	if eventBusName == "" {
		eventBusName = "default"
	}

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

	if len(targetIDs) == 0 {
		return nil, NewValidationException("Target IDs are required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	failedEntries := make([]map[string]interface{}, 0)
	failedCount := int32(0)

	for _, targetID := range targetIDs {
		if err := store.DeleteTarget(ctx, eventBusName, ruleName, targetID); err != nil {
			failedEntries = append(failedEntries, map[string]interface{}{
				"TargetId":     targetID,
				"ErrorCode":    "InternalFailure",
				"ErrorMessage": err.Error(),
			})
			failedCount++
		}
	}

	return map[string]interface{}{
		"FailedEntryCount": failedCount,
		"FailedEntries":    failedEntries,
	}, nil
}

// ListTargetsByRule lists targets for a rule in EventBridge.
func (s *EventsService) ListTargetsByRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	ruleName := request.GetParamLowerFirst(req.Parameters, "Rule")
	if ruleName == "" {
		return nil, NewValidationException("Rule name is required")
	}

	eventBusName := request.GetParamLowerFirst(req.Parameters, "EventBusName")
	if eventBusName == "" {
		eventBusName = "default"
	}
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")

	if limit == 0 {
		limit = 100
	}
	if limit < 1 || limit > 100 {
		return nil, NewValidationException("Limit must be between 1 and 100")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.ListTargetsByRule(ctx, eventBusName, ruleName, limit, nextToken)
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
	}

	response := map[string]interface{}{
		"Targets": targets,
	}

	if result.NextToken != "" {
		response["NextToken"] = result.NextToken
	}

	return response, nil
}
