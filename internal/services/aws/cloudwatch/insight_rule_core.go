package cloudwatch

import (
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	"vorpalstacks/internal/store/aws/common"
)

// PutInsightRuleInput holds parameters for PutInsightRule.
type PutInsightRuleInput struct {
	RuleName               string
	RuleState              string
	RuleDefinition         string
	ApplyOnTransformedLogs bool
	Tags                   map[string]string
}

// DeleteInsightRulesInput holds parameters for DeleteInsightRules.
type DeleteInsightRulesInput struct {
	RuleNames []string
}

// SetInsightRuleStateInput holds parameters for EnableInsightRules and
// DisableInsightRules.
type SetInsightRuleStateInput struct {
	RuleNames []string
	State     string
}

// DescribeInsightRulesInput holds parameters for DescribeInsightRules.
type DescribeInsightRulesInput struct {
	NextToken   string
	MaxResults  int
	ManagedRule bool
}

// ListManagedInsightRulesInput holds parameters for ListManagedInsightRules.
type ListManagedInsightRulesInput struct {
	ResourceARN string
	NextToken   string
	MaxResults  int
}

// PutManagedInsightRuleItem represents a single managed rule entry.
type PutManagedInsightRuleItem struct {
	TemplateName string
	ResourceARN  string
	Tags         map[string]string
}

// PutManagedInsightRulesInput holds parameters for PutManagedInsightRules.
type PutManagedInsightRulesInput struct {
	ManagedRules []PutManagedInsightRuleItem
}

// putInsightRuleCore validates input and creates or updates an insight rule.
func (s *CloudWatchService) putInsightRuleCore(stores *cloudwatchStores, input *PutInsightRuleInput) error {
	if err := validateInsightRuleName(input.RuleName); err != nil {
		return err
	}

	state := input.RuleState
	if state == "" {
		state = "ENABLED"
	}
	if err := validateRuleState(state); err != nil {
		return err
	}

	if err := validateInsightRuleDefinition(input.RuleDefinition); err != nil {
		return err
	}

	rule := &cwstore.InsightRule{
		Name:                   input.RuleName,
		State:                  state,
		Definition:             input.RuleDefinition,
		ApplyOnTransformedLogs: input.ApplyOnTransformedLogs,
		Tags:                   input.Tags,
	}

	_, err := stores.insightRules.PutInsightRule(rule)
	return err
}

// deleteInsightRulesCore deletes insight rules and returns names not found.
func (s *CloudWatchService) deleteInsightRulesCore(stores *cloudwatchStores, input *DeleteInsightRulesInput) ([]string, error) {
	if len(input.RuleNames) == 0 {
		return nil, awserrors.NewMissingParameter("RuleNames is required")
	}
	_, notFound, err := stores.insightRules.DeleteInsightRules(input.RuleNames)
	if err != nil {
		return nil, fmt.Errorf("failed to delete insight rules: %w", err)
	}
	return notFound, nil
}

// setInsightRuleStateCore enables or disables insight rules.
func (s *CloudWatchService) setInsightRuleStateCore(stores *cloudwatchStores, input *SetInsightRuleStateInput) ([]string, error) {
	if len(input.RuleNames) == 0 {
		return nil, awserrors.NewMissingParameter("RuleNames is required")
	}
	_, notFound, err := stores.insightRules.SetRuleState(input.RuleNames, input.State)
	if err != nil {
		return nil, fmt.Errorf("failed to set insight rule state: %w", err)
	}
	return notFound, nil
}

// describeInsightRulesCore lists insight rules with pagination.
func (s *CloudWatchService) describeInsightRulesCore(stores *cloudwatchStores, input *DescribeInsightRulesInput) ([]*cwstore.InsightRule, string, error) {
	opts := common.ListOptions{
		Marker:   input.NextToken,
		MaxItems: input.MaxResults,
	}
	result, err := stores.insightRules.ListInsightRulesPaginated(input.ManagedRule, opts)
	if err != nil {
		return nil, "", fmt.Errorf("failed to list insight rules: %w", err)
	}
	return result.Items, result.NextMarker, nil
}

// getInsightRuleCore retrieves a single insight rule by name.
func (s *CloudWatchService) getInsightRuleCore(stores *cloudwatchStores, ruleName string) (*cwstore.InsightRule, error) {
	if ruleName == "" {
		return nil, awserrors.NewMissingParameter("RuleName is required")
	}
	rule, err := stores.insightRules.GetInsightRule(ruleName)
	if err != nil {
		return nil, awserrors.NewResourceNotFoundException("InsightRule", ruleName)
	}
	return rule, nil
}

// listManagedInsightRulesCore lists managed insight rules.
func (s *CloudWatchService) listManagedInsightRulesCore(stores *cloudwatchStores, input *ListManagedInsightRulesInput) ([]*cwstore.InsightRule, string, error) {
	opts := common.ListOptions{
		Marker:   input.NextToken,
		MaxItems: input.MaxResults,
	}
	result, err := stores.insightRules.ListInsightRulesPaginated(true, opts)
	if err != nil {
		return nil, "", fmt.Errorf("failed to list managed insight rules: %w", err)
	}

	if input.ResourceARN != "" {
		filtered := make([]*cwstore.InsightRule, 0, len(result.Items))
		for _, r := range result.Items {
			if r.ResourceARN == input.ResourceARN {
				filtered = append(filtered, r)
			}
		}
		result.Items = filtered
	}

	return result.Items, result.NextMarker, nil
}

// putManagedInsightRulesCore creates managed insight rules.
func (s *CloudWatchService) putManagedInsightRulesCore(stores *cloudwatchStores, input *PutManagedInsightRulesInput) []map[string]interface{} {
	var failures []map[string]interface{}

	for _, mr := range input.ManagedRules {
		if mr.TemplateName == "" || mr.ResourceARN == "" {
			failures = append(failures, map[string]interface{}{
				"FailureName":        mr.TemplateName,
				"ExceptionName":      "InvalidParameterValueException",
				"FailureDescription": "TemplateName and ResourceARN are required",
			})
			continue
		}

		ruleName := fmt.Sprintf("ManagedRule:%s:%s", mr.TemplateName, mr.ResourceARN)

		rule := &cwstore.InsightRule{
			Name:         ruleName,
			State:        "ENABLED",
			ManagedRule:  true,
			TemplateName: mr.TemplateName,
			ResourceARN:  mr.ResourceARN,
			Tags:         mr.Tags,
		}

		if _, err := stores.insightRules.PutInsightRule(rule); err != nil {
			failures = append(failures, map[string]interface{}{
				"FailureName":        mr.TemplateName,
				"ExceptionName":      "InternalServiceFault",
				"FailureDescription": err.Error(),
			})
		}
	}
	return failures
}
