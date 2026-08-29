package eventbridge

import (
	"context"
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
// one so the merge semantics survive the transport boundary.
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

// DescribeRuleResult holds the outcome of describeRuleCore: the rule record
// plus its tags.
type DescribeRuleResult struct {
	Rule *eventsstore.Rule
	Tags []tagutil.Tag
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
		if !isValidEventPattern(input.EventPattern) {
			return nil, awserrors.NewInvalidEventPatternException("EventPattern must be valid JSON")
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

	if err := store.CreateRule(ctx, rule); err != nil {
		if err == eventsstore.ErrRuleAlreadyExists {
			// The user's fields are applied through the store-level atomic
			// mutation so a concurrent delivery-marker write can never be
			// lost to this update's read-modify-write cycle.
			if err := store.MutateRule(ctx, eventBusName, input.Name, func(existingRule *eventsstore.Rule) error {
				if input.DescriptionSet {
					if !validateDescription(input.Description) {
						return errDescriptionTooLong()
					}
					existingRule.Description = input.Description
				}
				if input.EventPatternSet {
					if !validateEventPatternLength(input.EventPattern) {
						return awserrors.NewValidationException("EventPattern must be at most 4096 characters")
					}
					if !isValidEventPattern(input.EventPattern) {
						return awserrors.NewInvalidEventPatternException("EventPattern must be valid JSON")
					}
					existingRule.EventPattern = input.EventPattern
				}
				if input.ScheduleExpressionSet {
					existingRule.ScheduleExpression = input.ScheduleExpression
				}
				if input.RoleArnSet {
					existingRule.RoleARN = input.RoleArn
				}
				if input.StateSet {
					existingRule.State = eventsstore.RuleState(input.State)
				}
				existingRule.LastModifiedAt = time.Now().UTC()
				return nil
			}); err != nil {
				return nil, err
			}
			existingRule, err := store.GetRule(ctx, eventBusName, input.Name)
			if err != nil {
				return nil, err
			}
			if len(input.Tags) > 0 {
				if err := store.TagStore.TagFromSlice(existingRule.ARN, input.Tags); err != nil {
					return nil, err
				}
			}
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

	rule, err := store.GetRule(ctx, eventBusName, input.Name)
	if err != nil {
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
	lastFireTimes.Delete(rule.ARN)

	_ = store.TagStore.Delete(rule.ARN)

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

	result := &DescribeRuleResult{Rule: rule}
	if tagSlice, err := store.TagStore.ListAsSlice(rule.ARN); err == nil && len(tagSlice) > 0 {
		result.Tags = tagSlice
	}
	return result, nil
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

	limit := input.Limit
	if limit == 0 {
		limit = 100
	}
	if limit < 1 || limit > 100 {
		return nil, awserrors.NewValidationException("Limit must be between 1 and 100")
	}

	var allRuleNames []string
	scanToken := ""
	for {
		rulesResult, err := store.ListRules(ctx, eventBusName, "", 100, scanToken)
		if err != nil {
			return nil, err
		}
		for _, rule := range rulesResult.Rules {
			targets, err := store.ListTargetsByRule(ctx, eventBusName, rule.Name, 100, "")
			if err != nil {
				// Swallowing the error here would silently omit rules
				// from the result, under-reporting which rules reference
				// the target. Fail the whole listing instead.
				return nil, err
			}
			for _, t := range targets.Targets {
				if t.ARN == input.TargetArn {
					allRuleNames = append(allRuleNames, rule.Name)
					break
				}
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
