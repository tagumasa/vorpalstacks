package eventbridge

import (
	"context"
	"strconv"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// ---------------------------------------------------------------------------
// Input / Result structs (transport-agnostic)
// ---------------------------------------------------------------------------

// PutTargetsInput carries the parameters for PutTargets. Targets holds the
// raw wire entries because every parse or validation failure of a single
// entry becomes a per-entry failure record rather than a request error, so
// the per-entry loop must run at the Core layer in its original order.
type PutTargetsInput struct {
	EventBusName         string
	EventBusNameProvided bool
	Rule                 string
	Targets              []interface{}
	IAMValidator         *iam.IAMValidator
}

// PutTargetsResult holds the outcome of putTargetsCore.
type PutTargetsResult struct {
	FailedEntryCount int32
	FailedEntries    []map[string]interface{}
}

// RemoveTargetsInput carries the parameters for RemoveTargets.
type RemoveTargetsInput struct {
	EventBusName         string
	EventBusNameProvided bool
	Rule                 string
	Ids                  []string
}

// ListTargetsByRuleInput carries the parameters for ListTargetsByRule.
type ListTargetsByRuleInput struct {
	EventBusName         string
	EventBusNameProvided bool
	Rule                 string
	Limit                int32
	NextToken            string
}

// ---------------------------------------------------------------------------
// Core functions
// ---------------------------------------------------------------------------

// putTargetsCore validates the rule and target entries, enforces the
// per-rule target quota and stores the targets, reporting per-entry
// failures.
func (s *EventsService) putTargetsCore(ctx context.Context, store *eventsstore.EventsStore, input PutTargetsInput) (*PutTargetsResult, error) {
	ruleName := input.Rule
	if ruleName == "" {
		return nil, awserrors.NewValidationException("Rule name is required")
	}
	eventBusName, err := resolveEventBusNameCore(input.EventBusName, input.EventBusNameProvided)
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

	targets := input.Targets
	if len(targets) == 0 {
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
				if input.IAMValidator != nil {
					if err := input.IAMValidator.ValidateRoleForService(ctx, roleArn, iam.ServicePrincipalEvents); err != nil {
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

	return &PutTargetsResult{
		FailedEntryCount: failedCount,
		FailedEntries:    failedEntries,
	}, nil
}

// removeTargetsCore validates the rule and deletes the requested target IDs,
// reporting per-entry failures.
func (s *EventsService) removeTargetsCore(ctx context.Context, store *eventsstore.EventsStore, input RemoveTargetsInput) (*PutTargetsResult, error) {
	ruleName := input.Rule
	if ruleName == "" {
		return nil, awserrors.NewValidationException("Rule name is required")
	}
	eventBusName, err := resolveEventBusNameCore(input.EventBusName, input.EventBusNameProvided)
	if err != nil {
		return nil, err
	}

	targetIDs := input.Ids
	if len(targetIDs) == 0 {
		return nil, awserrors.NewValidationException("Target IDs are required")
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

	return &PutTargetsResult{
		FailedEntryCount: failedCount,
		FailedEntries:    failedEntries,
	}, nil
}

// listTargetsByRuleCore validates the rule and limit and lists the rule's
// targets.
func (s *EventsService) listTargetsByRuleCore(ctx context.Context, store *eventsstore.EventsStore, input ListTargetsByRuleInput) (*eventsstore.TargetListResult, error) {
	ruleName := input.Rule
	if ruleName == "" {
		return nil, awserrors.NewValidationException("Rule name is required")
	}
	eventBusName, err := resolveEventBusNameCore(input.EventBusName, input.EventBusNameProvided)
	if err != nil {
		return nil, err
	}

	limit := input.Limit
	if limit == 0 {
		limit = 100
	}
	if limit < 1 || limit > 100 {
		return nil, awserrors.NewValidationException("Limit must be between 1 and 100")
	}

	if _, err := store.GetRule(ctx, eventBusName, ruleName); err != nil {
		return nil, awserrors.NewResourceNotFoundException("Rule", ruleName)
	}
	return store.ListTargetsByRule(ctx, eventBusName, ruleName, limit, input.NextToken)
}
