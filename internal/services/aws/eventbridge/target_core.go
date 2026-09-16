package eventbridge

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	svcarn "vorpalstacks/internal/utils/aws/arn"
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
	Region               string
	IAMValidator         *iam.IAMValidator
}

// TargetsMutationResult holds the outcome of putTargetsCore and
// removeTargetsCore: the per-entry failure report both operations return.
type TargetsMutationResult struct {
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

// validateTargetARNAcceptance adjudicates the accept/deliver matrix at
// PutTargets: a target ARN is accepted only when the platform has a
// delivery path for its service AND resource form. Every rejection names
// the reason in the per-entry failure.
//
// Accepted (delivery implemented in attemptDelivery):
//
//	lambda, sqs, sns, states (Step Functions), logs, kinesis — the six
//	  invoker/bus delivery arms.
//	events — the event-bus/ resource form (same- or cross-account bus
//	  delivery with the hop-depth guard) and the api-destination/ form
//	  (HTTPS invocation of the destination's endpoint, authorised by its
//	  connection); every other events-service resource form (rule, archive,
//	  a bare name) has no delivery path.
//	appsync — a GraphQL endpoint ARN (apis/<apiId>[/endpoints/GRAPHQL])
//	  with AppSyncParameters.GraphQLOperation; delivery invokes the
//	  mutation with the transformed payload as variables.
//	firehose — accepted because the Firehose service is a planned release
//	  blocker 4 deliverable; until it lands, delivery fails fast to the
//	  terminal handling instead of burning the retry budget.
//
// Rejected (no delivery path can exist on this platform):
//
//	ecs — the ECS service is not implemented (Basic Policy future
//	  expansion list).
//	ssm — SSM Run Command requires the SendCommand operation and an
//	  instance execution plane; the SSM service implements Parameter
//	  Store only (docs/services.md).
func validateTargetARNAcceptance(arn string) error {
	if arn == "" {
		return fmt.Errorf("Target ARN must not be empty")
	}
	_, service, _, _, resource := svcarn.SplitARN(arn)
	switch service {
	case "lambda", "sqs", "sns", "states", "logs", "kinesis", "firehose":
		return nil
	case "events":
		if !strings.HasPrefix(resource, "event-bus/") && !strings.HasPrefix(resource, "api-destination/") {
			return fmt.Errorf("Target ARN %s: only event-bus/ and api-destination/ resources of the events service have a delivery path", arn)
		}
		return nil
	case "appsync":
		if extractAppSyncApiIDFromARN(arn) == "" {
			return fmt.Errorf("Target ARN %s: AppSync targets must name a GraphQL API endpoint (apis/<apiId>)", arn)
		}
		return nil
	case "ecs":
		return fmt.Errorf("Target ARN %s: ECS task targets are not supported (the ECS service is not available on this platform)", arn)
	case "ssm":
		return fmt.Errorf("Target ARN %s: SSM Run Command targets are not supported (the SSM service implements Parameter Store only on this platform)", arn)
	default:
		return fmt.Errorf("Target ARN %s: service %q has no delivery path on this platform", arn, service)
	}
}

// putTargetsCore validates the rule and target entries, enforces the
// per-rule target quota and stores the targets, reporting per-entry
// failures.
func (s *EventsService) putTargetsCore(ctx context.Context, store *eventsstore.EventsStore, input PutTargetsInput) (*TargetsMutationResult, error) {
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

		if err := validateTargetARNAcceptance(targetArn); err != nil {
			failedEntries = append(failedEntries, map[string]interface{}{
				"TargetId":     targetID,
				"ErrorCode":    "ValidationException",
				"ErrorMessage": err.Error(),
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

		if err := validateTargetInputConfiguration(target); err != nil {
			failedEntries = append(failedEntries, map[string]interface{}{
				"TargetId":     targetID,
				"ErrorCode":    "ValidationException",
				"ErrorMessage": err.Error(),
			})
			failedCount++
			continue
		}

		if dlConfig, ok := targetMap["DeadLetterConfig"].(map[string]interface{}); ok {
			target.DeadLetterConfig = &eventsstore.DeadLetterConfig{}
			if arn, ok := dlConfig["Arn"].(string); ok {
				if err := validateDeadLetterQueueARN(arn, input.Region); err != nil {
					failedEntries = append(failedEntries, map[string]interface{}{
						"TargetId":     targetID,
						"ErrorCode":    "ValidationException",
						"ErrorMessage": err.Error(),
					})
					failedCount++
					continue
				}
				target.DeadLetterConfig.Arn = arn
			}
		}

		if retryPolicy, ok := targetMap["RetryPolicy"].(map[string]interface{}); ok {
			target.RetryPolicy = &eventsstore.RetryPolicy{}
			// The age member distinguishes "explicitly supplied" from
			// "omitted": an explicit value outside 60-86400 (zero
			// included) is a per-entry validation failure; an omitted
			// member keeps the deployment default.
			_, maxAgeProvided := retryPolicy["MaximumEventAgeInSeconds"]
			if maxAge, ok := retryPolicy["MaximumEventAgeInSeconds"].(float64); ok {
				target.RetryPolicy.MaximumEventAgeInSeconds = int32(maxAge)
			}
			if maxRetry, ok := retryPolicy["MaximumRetryAttempts"].(float64); ok {
				target.RetryPolicy.MaximumRetryAttempts = int32(maxRetry)
			}
			if !validateRetryPolicy(target.RetryPolicy, maxAgeProvided) {
				failedEntries = append(failedEntries, map[string]interface{}{
					"TargetId":  targetID,
					"ErrorCode": "ValidationException",
					"ErrorMessage": fmt.Sprintf("RetryPolicy: MaximumRetryAttempts must be 0-%d, MaximumEventAgeInSeconds must be %d-%d when provided",
						eventsstore.RetryPolicyMaxRetryAttempts,
						eventsstore.RetryPolicyMinEventAgeSeconds,
						eventsstore.RetryPolicyMaxEventAgeSeconds),
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

		if asp, ok := targetMap["AppSyncParameters"].(map[string]interface{}); ok {
			target.AppSyncParameters = &eventsstore.AppSyncParameters{}
			if op, ok := asp["GraphQLOperation"].(string); ok {
				target.AppSyncParameters.GraphQLOperation = op
			}
		}
		// An AppSync target without an operation document cannot invoke
		// anything: the mutation is the target's entire delivery payload.
		if _, service, _, _, _ := svcarn.SplitARN(targetArn); service == "appsync" &&
			(target.AppSyncParameters == nil || target.AppSyncParameters.GraphQLOperation == "") {
			failedEntries = append(failedEntries, map[string]interface{}{
				"TargetId":     targetID,
				"ErrorCode":    "ValidationException",
				"ErrorMessage": "AppSync targets require AppSyncParameters.GraphQLOperation (the mutation to invoke)",
			})
			failedCount++
			continue
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

	return &TargetsMutationResult{
		FailedEntryCount: failedCount,
		FailedEntries:    failedEntries,
	}, nil
}

// removeTargetsCore validates the rule and deletes the requested target IDs,
// reporting per-entry failures.
func (s *EventsService) removeTargetsCore(ctx context.Context, store *eventsstore.EventsStore, input RemoveTargetsInput) (*TargetsMutationResult, error) {
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
		// Removing a target ID that no longer exists succeeds: the API
		// reference frames a successful RemoveTargets as "the target(s)
		// listed in the request are removed" and documents no per-entry
		// error code for a missing target, so an idempotent removal is
		// the contract — not an InternalFailure entry carrying the raw
		// store error string.
		if err := store.DeleteTarget(ctx, eventBusName, ruleName, targetID); err != nil && err != eventsstore.ErrTargetNotFound {
			failedEntries = append(failedEntries, map[string]interface{}{
				"TargetId":     targetID,
				"ErrorCode":    "InternalFailure",
				"ErrorMessage": err.Error(),
			})
			failedCount++
		}
	}

	return &TargetsMutationResult{
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

	limit, err := normaliseListLimit(input.Limit)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetRule(ctx, eventBusName, ruleName); err != nil {
		return nil, awserrors.NewResourceNotFoundException("Rule", ruleName)
	}
	return store.ListTargetsByRule(ctx, eventBusName, ruleName, limit, input.NextToken)
}
