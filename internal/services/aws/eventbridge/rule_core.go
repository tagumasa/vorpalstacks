package eventbridge

import (
	"context"
	"fmt"
	"strconv"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
	tagutil "vorpalstacks/internal/common/tags"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// ---------------------------------------------------------------------------
// Input / Result structs (transport-agnostic)
// ---------------------------------------------------------------------------

// PutRuleInput carries the parameters for PutRule (create and upsert). The
// *Set flags distinguish an omitted member from an explicitly provided empty
// one: validation applies to supplied values alone, while the stored outcome
// of an upsert follows the model's replacement contract — omitted members
// clear rather than keep their stored value.
type PutRuleInput struct {
	Name                  string
	EventBusName          string
	EventBusNameProvided  bool
	CreatedBy             string
	DescriptionSet        bool
	Description           string
	EventPatternSet       bool
	EventPattern          string
	ScheduleExpressionSet bool
	ScheduleExpression    string
	StateSet              bool
	State                 string
	RoleArnSet            bool
	RoleArn               string
	Tags                  []tagutil.Tag
	IAMValidator          *iam.IAMValidator
}

// PutRuleResult holds the outcome of putRuleCore.
type PutRuleResult struct {
	RuleArn string
}

// DeleteRuleInput carries the parameters for DeleteRule.
type DeleteRuleInput struct {
	EventBusName         string
	EventBusNameProvided bool
	Name                 string
	Force                bool
}

// DescribeRuleInput carries the parameters for DescribeRule.
type DescribeRuleInput struct {
	EventBusName         string
	EventBusNameProvided bool
	Name                 string
}

// SetRuleStateInput carries the parameters for the shared EnableRule and
// DisableRule state transition.
type SetRuleStateInput struct {
	EventBusName         string
	EventBusNameProvided bool
	Name                 string
	State                eventsstore.RuleState
}

// DescribeRuleResult holds the outcome of describeRuleCore: the rule
// record. Tags are not a DescribeRuleResponse member — the tag read
// surface is ListTagsForResource.
type DescribeRuleResult struct {
	Rule *eventsstore.Rule
}

// ListRuleNamesByTargetInput carries the parameters for
// ListRuleNamesByTarget.
type ListRuleNamesByTargetInput struct {
	EventBusName         string
	EventBusNameProvided bool
	TargetArn            string
	Limit                int32
	NextToken            string
}

// ListRuleNamesByTargetResult holds the outcome of
// listRuleNamesByTargetCore.
type ListRuleNamesByTargetResult struct {
	RuleNames []string
	NextToken string
}

// ---------------------------------------------------------------------------
// Core functions
// ---------------------------------------------------------------------------

// putRuleCore creates or updates a rule on the specified event bus.
// Supports event patterns and schedule expressions (cron/rate). When the
// rule already exists the user's fields are applied through the store-level
// atomic mutation so a concurrent delivery-marker write can never be lost to
// this update's read-modify-write cycle.
func (s *EventsService) putRuleCore(ctx context.Context, store *eventsstore.EventsStore, input PutRuleInput) (*PutRuleResult, error) {
	if input.Name == "" {
		return nil, awserrors.NewValidationException("Rule name is required")
	}
	if !validateResourceName(input.Name, "rule") {
		return nil, awserrors.NewValidationException("Rule name must match the pattern and be 1-64 characters")
	}

	eventBusName, err := resolveEventBusNameCore(input.EventBusName, input.EventBusNameProvided)
	if err != nil {
		return nil, err
	}

	// Check if event bus exists, auto-create default event bus if needed
	if _, err := store.GetEventBus(ctx, eventBusName); err != nil {
		if err == eventsstore.ErrEventBusNotFound {
			if eventBusName == "default" {
				if createErr := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: eventBusName}); createErr != nil {
					return nil, createErr
				}
			} else {
				return nil, NewResourceNotFoundException("Event bus '" + eventBusName + "' does not exist")
			}
		} else {
			return nil, err
		}
	}

	// "You can only create scheduled rules using the default event bus"
	// (user guide, Creating an Amazon EventBridge rule that runs on a
	// schedule). The restriction is enforced for both the create and the
	// upsert path: a scheduled rule on a custom bus can never legitimately
	// exist, so moving an existing custom-bus rule onto a schedule is
	// rejected the same way.
	if input.ScheduleExpressionSet && eventBusName != "default" {
		return nil, awserrors.NewValidationException("Scheduled rules can only be created on the default event bus")
	}

	rule := &eventsstore.Rule{
		Name:         input.Name,
		EventBusName: eventBusName,
		CreatedBy:    input.CreatedBy,
	}

	if input.DescriptionSet {
		if !validateDescription(input.Description) {
			return nil, errDescriptionTooLong()
		}
		rule.Description = input.Description
	}

	if input.EventPatternSet {
		if !validateEventPatternLength(input.EventPattern) {
			return nil, awserrors.NewValidationException("EventPattern must be at most 4096 characters")
		}
		if err := validateEventPatternStructure(input.EventPattern); err != nil {
			return nil, err
		}
		rule.EventPattern = input.EventPattern
	}

	if input.ScheduleExpressionSet {
		if !isValidScheduleExpression(input.ScheduleExpression) {
			return nil, awserrors.NewValidationException("ScheduleExpression must be a valid rate() or cron() expression")
		}
		rule.ScheduleExpression = input.ScheduleExpression
	}

	// AWS requires a rule to contain at least an EventPattern or
	// ScheduleExpression.  A rule with neither is rejected with
	// ValidationException.
	if rule.EventPattern == "" && rule.ScheduleExpression == "" {
		return nil, awserrors.NewValidationException("A rule must contain at least an EventPattern or ScheduleExpression")
	}

	if input.StateSet {
		if !validateRuleState(input.State) {
			return nil, awserrors.NewValidationException("State must be 'ENABLED', 'DISABLED', or 'ENABLED_WITH_ALL_CLOUDTRAIL_MANAGEMENT_EVENTS'")
		}
		rule.State = eventsstore.RuleState(input.State)
	} else {
		rule.State = eventsstore.RuleStateEnabled
	}

	if input.RoleArnSet {
		if input.RoleArn != "" {
			if s.bus != nil {
				if rr := s.bus.RoleResolver(); rr != nil {
					if err := rr.ValidateRole(ctx, input.RoleArn); err != nil {
						return nil, err
					}
				}
			}
			if input.IAMValidator != nil {
				if err := input.IAMValidator.ValidateRoleForService(ctx, input.RoleArn, iam.ServicePrincipalEvents); err != nil {
					return nil, err
				}
			}
		}
		rule.RoleARN = input.RoleArn
	}

	// "Maximum number of rules an account can have per event bus" — 300
	// per the EventBridge quotas page; the model carries
	// LimitExceededException on PutRule for the breach. The gate lives in
	// the store's capped create (one locked section with the write), and
	// the upsert path below never hits it: the quota bounds the creation
	// of new rules alone.
	if err := store.CreateRuleCapped(ctx, rule, eventsstore.MaxRulesPerEventBus); err != nil {
		if err == eventsstore.ErrRuleCapReached {
			return nil, awserrors.NewLimitExceededException(fmt.Sprintf(
				"The maximum of %d rules per event bus has been reached", eventsstore.MaxRulesPerEventBus))
		}
		if err == eventsstore.ErrRuleAlreadyExists {
			// The user's fields are applied through the store-level atomic
			// mutation so a concurrent delivery-marker write can never be
			// lost to this update's read-modify-write cycle. Every member
			// was validated above, before the create attempt, so the
			// mutation assigns without re-validating — one validation
			// site per member, exercised by both the create and the
			// upsert path.
			//
			// "If you are updating an existing rule, the rule is replaced
			// with what you specify in this PutRule command. If you omit
			// arguments in PutRule, the old values for those arguments
			// are not kept. Instead, they are replaced with null values",
			// and "Rules are enabled by default" (model operation
			// documentation): the update replaces every request member
			// with the input-built rule's value — omitted members clear,
			// and an omitted State resets to ENABLED. The record's
			// identity members (ARN, CreatedAt, CreatedBy, ManagedBy) and
			// the scheduler bookkeeping survive the replacement.
			if err := store.MutateRule(ctx, eventBusName, input.Name, func(existingRule *eventsstore.Rule) error {
				existingRule.Description = rule.Description
				existingRule.EventPattern = rule.EventPattern
				existingRule.ScheduleExpression = rule.ScheduleExpression
				existingRule.RoleARN = rule.RoleARN
				existingRule.State = rule.State
				existingRule.LastModifiedAt = time.Now().UTC()
				return nil
			}); err != nil {
				return nil, err
			}
			existingRule, err := store.GetRule(ctx, eventBusName, input.Name)
			if err != nil {
				return nil, err
			}
			// "If you are updating an existing rule, any tags you specify
			// in the PutRule operation are ignored" (model operation
			// documentation) — tags ride the create path alone.
			return &PutRuleResult{RuleArn: existingRule.ARN}, nil
		}
		return nil, err
	}

	if len(input.Tags) > 0 {
		if err := store.TagStore.TagFromSlice(rule.ARN, input.Tags); err != nil {
			return nil, err
		}
	}

	return &PutRuleResult{RuleArn: rule.ARN}, nil
}

// deleteRuleCore validates input and removes a rule from the event bus.
// Rules with targets cannot be deleted until targets are removed, unless
// Force is set in which case the targets are cascade-deleted first.
func (s *EventsService) deleteRuleCore(ctx context.Context, store *eventsstore.EventsStore, input DeleteRuleInput) error {
	if input.Name == "" {
		return awserrors.NewValidationException("Rule name is required")
	}

	eventBusName, err := resolveEventBusNameCore(input.EventBusName, input.EventBusNameProvided)
	if err != nil {
		return err
	}

	// "If you call delete rule multiple times for the same rule, all calls
	// will succeed. When you call delete rule for a non-existent custom
	// eventbus, ResourceNotFoundException is returned" (model operation
	// documentation) — so the bus's existence is checked first and a missing
	// rule on an existing bus is a successful no-op. The default bus always
	// exists on AWS; here it may not have been created yet, but no rule can
	// exist without it either, so the delete still succeeds.
	if _, err := store.GetEventBus(ctx, eventBusName); err != nil {
		if err == eventsstore.ErrEventBusNotFound && eventBusName == "default" {
			return nil
		}
		return mapStoreError(err, eventBusName)
	}

	rule, err := store.GetRule(ctx, eventBusName, input.Name)
	if err != nil {
		if err == eventsstore.ErrRuleNotFound {
			return nil
		}
		return mapStoreError(err, input.Name)
	}

	force := input.Force

	// Check if rule has targets. When Force is false (the default), AWS
	// rejects the delete with a DependencyViolation-style error so that
	// callers must explicitly remove targets first. When Force is true the
	// targets are cascade-deleted before the rule itself is removed.
	targetsResult, err := store.ListTargetsByRule(ctx, eventBusName, input.Name, 1, "")
	if err != nil {
		return err
	}
	if len(targetsResult.Targets) > 0 && !force {
		return awserrors.NewValidationException("Rule '" + input.Name + "' has targets. Remove targets before deleting the rule, or retry with Force=true.")
	}
	if len(targetsResult.Targets) > 0 && force {
		// Cascade-delete targets before removing the rule.
		allToken := ""
		for {
			page, err := store.ListTargetsByRule(ctx, eventBusName, input.Name, 1000, allToken)
			if err != nil {
				return err
			}
			for _, t := range page.Targets {
				if err := store.DeleteTarget(ctx, eventBusName, input.Name, t.ID); err != nil {
					return err
				}
			}
			if page.NextToken == "" {
				break
			}
			allToken = page.NextToken
		}
	}

	if err := store.DeleteRule(ctx, eventBusName, input.Name); err != nil {
		return err
	}

	// Clean up scheduler state for the deleted rule.
	s.fireDedup.deleteLastFire(rule.ARN)

	return nil
}

// describeRuleCore validates input and fetches the rule with its tags.
func (s *EventsService) describeRuleCore(ctx context.Context, store *eventsstore.EventsStore, input DescribeRuleInput) (*DescribeRuleResult, error) {
	if input.Name == "" {
		return nil, awserrors.NewValidationException("Rule name is required")
	}

	eventBusName, err := resolveEventBusNameCore(input.EventBusName, input.EventBusNameProvided)
	if err != nil {
		return nil, err
	}

	rule, err := store.GetRule(ctx, eventBusName, input.Name)
	if err != nil {
		return nil, mapStoreError(err, input.Name)
	}

	return &DescribeRuleResult{Rule: rule}, nil
}

// setRuleStateCore validates input and transitions the rule to the given
// state. Shared by EnableRule and DisableRule.
func (s *EventsService) setRuleStateCore(ctx context.Context, store *eventsstore.EventsStore, input SetRuleStateInput) error {
	if input.Name == "" {
		return awserrors.NewValidationException("Rule name is required")
	}

	eventBusName, err := resolveEventBusNameCore(input.EventBusName, input.EventBusNameProvided)
	if err != nil {
		return err
	}

	if err := store.MutateRule(ctx, eventBusName, input.Name, func(rule *eventsstore.Rule) error {
		rule.State = input.State
		rule.LastModifiedAt = time.Now().UTC()
		return nil
	}); err != nil {
		return mapStoreError(err, input.Name)
	}
	return nil
}

// listRuleNamesByTargetCore scans all rules on the event bus, checks each
// rule's targets for a match, then applies offset pagination to the matched
// rule names.
func (s *EventsService) listRuleNamesByTargetCore(ctx context.Context, store *eventsstore.EventsStore, input ListRuleNamesByTargetInput) (*ListRuleNamesByTargetResult, error) {
	if input.TargetArn == "" {
		return nil, awserrors.NewValidationException("TargetArn is required")
	}

	eventBusName, err := resolveEventBusNameCore(input.EventBusName, input.EventBusNameProvided)
	if err != nil {
		return nil, err
	}

	limit, err := normaliseListLimit(input.Limit)
	if err != nil {
		return nil, err
	}

	var allRuleNames []string
	scanToken := ""
	for {
		rulesResult, err := store.ListRules(ctx, eventBusName, "", 100, scanToken)
		if err != nil {
			return nil, err
		}
		for _, rule := range rulesResult.Rules {
			// The per-rule target scan paginates like the identical scan
			// on the delivery path: references beyond the first page are
			// references too. (The 5-targets-per-rule API quota keeps
			// live rules single-page; the bound holds the scan correct
			// regardless.)
			targetToken := ""
			matched := false
			for !matched {
				targets, err := store.ListTargetsByRule(ctx, eventBusName, rule.Name, 100, targetToken)
				if err != nil {
					// Swallowing the error here would silently omit rules
					// from the result, under-reporting which rules reference
					// the target. Fail the whole listing instead.
					return nil, err
				}
				for _, t := range targets.Targets {
					if t.ARN == input.TargetArn {
						allRuleNames = append(allRuleNames, rule.Name)
						matched = true
						break
					}
				}
				if matched || targets.NextToken == "" {
					break
				}
				targetToken = targets.NextToken
			}
		}
		if rulesResult.NextToken == "" {
			break
		}
		scanToken = rulesResult.NextToken
	}

	start := 0
	if input.NextToken != "" {
		// AWS returns InvalidParameterException when the supplied
		// NextToken is not a recognised opaque cursor. Our pagination
		// scheme uses a zero-based integer offset, so reject any token
		// that fails to parse as a non-negative integer.
		idx, err := strconv.Atoi(input.NextToken)
		if err != nil || idx < 0 {
			return nil, awserrors.NewInvalidParameterException("Invalid NextToken: " + input.NextToken)
		}
		start = idx
	}
	end := start + int(limit)
	// Clamp both endpoints to the result-set length to prevent slice
	// bounds panics when a stale NextToken offset exceeds the current
	// length (e.g. rules were deleted between paginated calls).
	if start > len(allRuleNames) {
		start = len(allRuleNames)
	}
	if end > len(allRuleNames) {
		end = len(allRuleNames)
	}

	ruleNames := allRuleNames[start:end]

	result := &ListRuleNamesByTargetResult{RuleNames: ruleNames}
	if end < len(allRuleNames) {
		result.NextToken = strconv.Itoa(end)
	}
	return result, nil
}
