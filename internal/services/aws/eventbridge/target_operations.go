package eventbridge

import (
	"context"
	"strconv"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

const maxTargetsPerRule = 5

// isValidTargetARN validates target ARNs at PutTargets time. The service
// segment of the ARN must match a supported delivery type in the
// dispatchToTarget switch (event_operations.go).
//
// Currently supported (delivery implemented):
//
//	lambda, sqs, sns, events, states (Step Functions), logs, kinesis
//
// Future targets (infrastructure exists, delivery pending service implementation):
//
//	ecs — EcsParameters types/parsing/validation fully implemented. Delivery
//	       stub exists. Enable by implementing deliverToECS when the ECS
//	       service is available on this platform.
//	firehose — Delivery stub exists, no sub-parameters required (Smithy model
//	           has no FirehoseParameters). Enable by implementing
//	           deliverToFirehose when the Firehose service is available.
//
// Platform-implemented, delivery not yet wired:
//
//	ssm — RunCommandParameters type exists, SSM service exists. Delivery
//	       TODO: call SSM StartAutomationExecution.
//	appsync — AppSyncParameters type exists, AppSync service exists. Delivery
//	          TODO: call AppSync GraphQL API.
//
// Out of scope (permanently unsupported on this edge/on-prem platform):
//
//	sagemaker — ML pipeline service (types stripped)
//	batch — Batch processing service (types stripped)
//	redshift — Data warehouse service (types stripped)
//	codebuild — CI/CD build service
//	codepipeline — CI/CD pipeline orchestration
//	inspector — Security assessment service
func isValidTargetARN(arn string) bool {
	if arn == "" {
		return false
	}
	_, service, _, _, _ := svcarn.SplitARN(arn)
	validServices := map[string]bool{
		"lambda":   true,
		"sqs":      true,
		"sns":      true,
		"events":   true,
		"ecs":      true,
		"firehose": true,
		"kinesis":  true,
		"states":   true,
		"logs":     true,
		"ssm":      true,
		"appsync":  true,
	}
	return validServices[service]
}

// PutTargets adds targets to a rule in EventBridge.
func (s *EventsService) PutTargets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	ruleName := request.GetParamLowerFirst(req.Parameters, "Rule")
	if ruleName == "" {
		return nil, awserrors.NewValidationException("Rule name is required")
	}

	eventBusName, err := resolveEventBusName(req)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Check if event bus exists
	if _, err := store.GetEventBus(ctx, eventBusName); err != nil {
		return nil, mapStoreError(err, eventBusName)
	}

	_, err = store.GetRule(ctx, eventBusName, ruleName)
	if err != nil {
		return nil, mapStoreError(err, ruleName)
	}

	targets, ok := req.Parameters["Targets"].([]interface{})
	if !ok {
		targets, ok = req.Parameters["targets"].([]interface{})
	}
	if !ok || len(targets) == 0 {
		return nil, awserrors.NewValidationException("Targets are required")
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
			return nil, awserrors.NewValidationException("Duplicate target ID: " + targetID)
		}
		seenIDs[targetID] = true
	}

	// Enforce the 5-targets-per-rule limit (AWS quota).
	existingTargets := make(map[string]bool)
	existToken := ""
	for {
		existingResult, err := store.ListTargetsByRule(ctx, eventBusName, ruleName, 100, existToken)
		if err != nil {
			return nil, awserrors.NewInternalFailureException("Failed to list existing targets for rule '" + ruleName + "': " + err.Error())
		}
		for _, et := range existingResult.Targets {
			existingTargets[et.ID] = true
		}
		if existingResult.NextToken == "" {
			break
		}
		existToken = existingResult.NextToken
	}
	newCount := 0
	for _, t := range targets {
		targetMap, ok := t.(map[string]interface{})
		if !ok {
			continue
		}
		targetID, _ := targetMap["Id"].(string)
		if !existingTargets[targetID] {
			newCount++
		}
	}
	if len(existingTargets)+newCount > maxTargetsPerRule {
		return nil, awserrors.NewValidationException(
			"Rule '" + ruleName + "' already has the maximum of " +
				strconv.Itoa(maxTargetsPerRule) + " targets. " +
				"Remove a target before adding new ones.")
	}

	failedEntries := make([]map[string]interface{}, 0)
	failedCount := int32(0)

	for _, t := range targets {
		targetMap, ok := t.(map[string]interface{})
		if !ok {
			failedEntries = append(failedEntries, map[string]interface{}{
				"TargetId":     "",
				"ErrorCode":    "ValidationException",
				"ErrorMessage": "Target entry must be an object",
			})
			failedCount++
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
			if !validateRetryPolicy(target.RetryPolicy) {
				failedEntries = append(failedEntries, map[string]interface{}{
					"TargetId":     targetID,
					"ErrorCode":    "ValidationException",
					"ErrorMessage": "RetryPolicy: MaximumRetryAttempts must be 0-185, MaximumEventAgeInSeconds must be 60-86400",
				})
				failedCount++
				continue
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

		if kinesisParams, ok := targetMap["KinesisParameters"].(map[string]interface{}); ok {
			target.KinesisParameters = &eventsstore.KinesisParameters{}
			if pkPath, ok := kinesisParams["PartitionKeyPath"].(string); ok {
				target.KinesisParameters.PartitionKeyPath = pkPath
			}
		}

		if rcp, ok := targetMap["RunCommandParameters"].(map[string]interface{}); ok {
			target.RunCommandParameters = parseRunCommandParameters(rcp)
		}
		if asp, ok := targetMap["AppSyncParameters"].(map[string]interface{}); ok {
			target.AppSyncParameters = &eventsstore.AppSyncParameters{}
			if op, ok := asp["GraphQLOperation"].(string); ok {
				target.AppSyncParameters.GraphQLOperation = op
			}
		}
		if ecs, ok := targetMap["EcsParameters"].(map[string]interface{}); ok {
			ecsParams, err := parseEcsParameters(ecs)
			if err != nil {
				failedEntries = append(failedEntries, map[string]interface{}{
					"TargetId":     targetID,
					"ErrorCode":    "ValidationException",
					"ErrorMessage": err.Error(),
				})
				failedCount++
				continue
			}
			target.EcsParameters = ecsParams
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
		return nil, awserrors.NewValidationException("Rule name is required")
	}

	eventBusName, err := resolveEventBusName(req)
	if err != nil {
		return nil, err
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
		return nil, awserrors.NewValidationException("Target IDs are required")
	}

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

	_, err = store.GetRule(ctx, eventBusName, ruleName)
	if err != nil {
		return nil, mapStoreError(err, ruleName)
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
		return nil, awserrors.NewValidationException("Rule name is required")
	}

	eventBusName, err := resolveEventBusName(req)
	if err != nil {
		return nil, err
	}
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")

	if limit == 0 {
		limit = 100
	}
	if limit < 1 || limit > 100 {
		return nil, awserrors.NewValidationException("Limit must be between 1 and 100")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetRule(ctx, eventBusName, ruleName); err != nil {
		return nil, awserrors.NewResourceNotFoundException("Rule", ruleName)
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
		if t.KinesisParameters != nil {
			targets[i]["KinesisParameters"] = map[string]interface{}{
				"PartitionKeyPath": t.KinesisParameters.PartitionKeyPath,
			}
		}
		if t.RunCommandParameters != nil {
			rcTargets := make([]map[string]interface{}, len(t.RunCommandParameters.RunCommandTargets))
			for j, rct := range t.RunCommandParameters.RunCommandTargets {
				rcTargets[j] = map[string]interface{}{
					"Key":    rct.Key,
					"Values": rct.Values,
				}
			}
			targets[i]["RunCommandParameters"] = map[string]interface{}{
				"RunCommandTargets": rcTargets,
			}
		}
		if t.AppSyncParameters != nil {
			targets[i]["AppSyncParameters"] = map[string]interface{}{
				"GraphQLOperation": t.AppSyncParameters.GraphQLOperation,
			}
		}
		if t.EcsParameters != nil {
			ecsp := map[string]interface{}{}
			if t.EcsParameters.TaskDefinitionArn != "" {
				ecsp["TaskDefinitionArn"] = t.EcsParameters.TaskDefinitionArn
			}
			if t.EcsParameters.TaskCount != 0 {
				ecsp["TaskCount"] = t.EcsParameters.TaskCount
			}
			if t.EcsParameters.LaunchType != "" {
				ecsp["LaunchType"] = t.EcsParameters.LaunchType
			}
			if t.EcsParameters.NetworkConfiguration != nil {
				ecsp["NetworkConfiguration"] = t.EcsParameters.NetworkConfiguration
			}
			if t.EcsParameters.PlatformVersion != "" {
				ecsp["PlatformVersion"] = t.EcsParameters.PlatformVersion
			}
			if t.EcsParameters.Group != "" {
				ecsp["Group"] = t.EcsParameters.Group
			}
			if len(t.EcsParameters.CapacityProviderStrategy) > 0 {
				ecsp["CapacityProviderStrategy"] = t.EcsParameters.CapacityProviderStrategy
			}
			if t.EcsParameters.EnableECSManagedTags {
				ecsp["EnableECSManagedTags"] = true
			}
			if t.EcsParameters.EnableExecuteCommand {
				ecsp["EnableExecuteCommand"] = true
			}
			if len(t.EcsParameters.PlacementConstraints) > 0 {
				ecsp["PlacementConstraints"] = t.EcsParameters.PlacementConstraints
			}
			if len(t.EcsParameters.PlacementStrategy) > 0 {
				ecsp["PlacementStrategy"] = t.EcsParameters.PlacementStrategy
			}
			if t.EcsParameters.PropagateTags != "" {
				ecsp["PropagateTags"] = t.EcsParameters.PropagateTags
			}
			if t.EcsParameters.ReferenceId != "" {
				ecsp["ReferenceId"] = t.EcsParameters.ReferenceId
			}
			targets[i]["EcsParameters"] = ecsp
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

// parseRunCommandParameters builds RunCommandParameters from a request map.
// At least one RunCommandTarget with both Key and Values is required.
func parseRunCommandParameters(m map[string]interface{}) *eventsstore.RunCommandParameters {
	out := &eventsstore.RunCommandParameters{}
	if targets, ok := m["RunCommandTargets"].([]interface{}); ok {
		for _, t := range targets {
			tm, ok := t.(map[string]interface{})
			if !ok {
				continue
			}
			rct := eventsstore.RunCommandTarget{
				Key: getStringField(tm, "Key"),
			}
			if vals, ok := tm["Values"].([]interface{}); ok {
				for _, v := range vals {
					if vs, ok := v.(string); ok {
						rct.Values = append(rct.Values, vs)
					}
				}
			}
			out.RunCommandTargets = append(out.RunCommandTargets, rct)
		}
	}
	return out
}

// parseEcsParameters captures ECS task target parameters.  ECS delivery is
// not available on this platform; parameters are persisted for SDK parity.
func parseEcsParameters(m map[string]interface{}) (*eventsstore.EcsParameters, error) {
	out := &eventsstore.EcsParameters{
		TaskDefinitionArn: getStringField(m, "TaskDefinitionArn"),
		LaunchType:        getStringField(m, "LaunchType"),
		PlatformVersion:   getStringField(m, "PlatformVersion"),
		Group:             getStringField(m, "Group"),
		PropagateTags:     getStringField(m, "PropagateTags"),
		ReferenceId:       getStringField(m, "ReferenceId"),
	}
	if v, ok := m["TaskCount"].(float64); ok {
		out.TaskCount = int32(v)
	}
	if v, ok := m["EnableECSManagedTags"].(bool); ok {
		out.EnableECSManagedTags = v
	}
	if v, ok := m["EnableExecuteCommand"].(bool); ok {
		out.EnableExecuteCommand = v
	}
	if lt := out.LaunchType; lt != "" && !validateLaunchType(lt) {
		return nil, awserrors.NewValidationException("LaunchType must be one of: EC2, FARGATE, EXTERNAL")
	}
	if nc, ok := m["NetworkConfiguration"].(map[string]interface{}); ok {
		out.NetworkConfiguration = nc
	}
	if cps, ok := m["CapacityProviderStrategy"].([]interface{}); ok {
		for _, cp := range cps {
			if cpm, ok := cp.(map[string]interface{}); ok {
				out.CapacityProviderStrategy = append(out.CapacityProviderStrategy, cpm)
			}
		}
	}
	if pcs, ok := m["PlacementConstraints"].([]interface{}); ok {
		for _, pc := range pcs {
			if pcm, ok := pc.(map[string]interface{}); ok {
				out.PlacementConstraints = append(out.PlacementConstraints, pcm)
			}
		}
	}
	if pss, ok := m["PlacementStrategy"].([]interface{}); ok {
		for _, ps := range pss {
			if psm, ok := ps.(map[string]interface{}); ok {
				out.PlacementStrategy = append(out.PlacementStrategy, psm)
			}
		}
	}
	return out, nil
}
