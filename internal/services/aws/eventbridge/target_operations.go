package eventbridge

import (
	"context"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

const maxTargetsPerRule = 5

// isValidTargetARN validates target ARNs at PutTargets time. The service
// segment of the ARN must match a supported delivery type in the
// dispatchToTarget switch (delivery_core.go).
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
