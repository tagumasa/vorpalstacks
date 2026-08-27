package iot

import (
	"fmt"
	"log/slog"
	"time"

	"vorpalstacks/internal/services/aws/iot/rules"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Topic-rule Core. Persistence goes through the rule store; after every
// mutation the regional rules executor is re-synchronised so the data plane
// matches the stored rule set.
// ---------------------------------------------------------------------------

// TopicRuleInput carries the fields of a topic-rule create or replace.
// Actions keeps the flat action map keyed by action type; ErrorAction is
// the single error action map.
type TopicRuleInput struct {
	RuleName         string
	SQL              string
	Description      string
	AwsIotSqlVersion string
	RuleDisabled     bool
	Actions          map[string]interface{}
	ErrorAction      map[string]interface{}
}

// CreateTopicRuleResult is the transport-agnostic result of CreateTopicRule.
type CreateTopicRuleResult struct {
	RuleArn  string
	RuleName string
}

// actionsToList converts a flat action map (keyed by type) into a list of
// single-key maps, matching the AWS IoT rule action representation.
func actionsToList(m map[string]interface{}) []map[string]interface{} {
	if len(m) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(m))
	for k, v := range m {
		result = append(result, map[string]interface{}{k: v})
	}
	return result
}

// createTopicRuleCore validates the SQL statement, persists the rule, and
// registers it with the regional executor.
func (s *IoTService) createTopicRuleCore(store iotstore.IotStoreInterface, in TopicRuleInput) (*CreateTopicRuleResult, error) {
	if in.RuleName == "" {
		return nil, iotstore.ErrMissingParam
	}

	rule := &iotstore.TopicRule{
		RuleName:         in.RuleName,
		SQL:              in.SQL,
		Description:      in.Description,
		AwsIotSqlVersion: in.AwsIotSqlVersion,
		CreatedAt:        fmt.Sprintf("%d", time.Now().UTC().Unix()),
		RuleDisabled:     in.RuleDisabled,
		Actions:          in.Actions,
	}
	if in.ErrorAction != nil {
		rule.ErrorAction = in.ErrorAction
	}

	// Validate the IoT SQL statement before persisting.
	if rule.SQL != "" {
		if _, err := rules.NewParser(rule.SQL).Parse(); err != nil {
			return nil, iotstore.ErrSqlParse
		}
	}

	created, err := store.CreateRule(rule)
	if err != nil {
		return nil, err
	}

	if executor := s.ExecutorForRegion(store.GetRegion()); !created.RuleDisabled && executor != nil && len(created.Actions) > 0 {
		if err := executor.AddRule(created.RuleName, created.TopicPattern, created.SQL, actionsToList(created.Actions), created.ErrorAction); err != nil {
			slog.Warn("rule created but executor registration failed", "rule", created.RuleName, "error", err)
		}
	}

	return &CreateTopicRuleResult{
		RuleArn:  created.ARN,
		RuleName: created.RuleName,
	}, nil
}

// describeTopicRuleCore retrieves a single rule by name.
func (s *IoTService) describeTopicRuleCore(store iotstore.IotStoreInterface, ruleName string) (*iotstore.TopicRule, error) {
	if ruleName == "" {
		return nil, iotstore.ErrMissingParam
	}
	return store.GetRule(ruleName)
}

// replaceTopicRuleCore validates the new definition, applies it, and
// re-registers the rule with the regional executor.
func (s *IoTService) replaceTopicRuleCore(store iotstore.IotStoreInterface, in TopicRuleInput) error {
	if in.RuleName == "" {
		return iotstore.ErrMissingParam
	}

	opts := iotstore.RuleUpdateOpts{
		SQL:              in.SQL,
		Description:      in.Description,
		AwsIotSqlVersion: in.AwsIotSqlVersion,
		RuleDisabled:     &in.RuleDisabled,
		Actions:          in.Actions,
		ErrorAction:      in.ErrorAction,
	}

	// Validate the IoT SQL statement before persisting.
	if opts.SQL != "" {
		if _, err := rules.NewParser(opts.SQL).Parse(); err != nil {
			return iotstore.ErrSqlParse
		}
	}

	updated, err := store.UpdateRule(in.RuleName, opts)
	if err != nil {
		return err
	}

	if executor := s.ExecutorForRegion(store.GetRegion()); executor != nil {
		executor.RemoveRule(in.RuleName)
		if !updated.RuleDisabled && updated.SQL != "" && len(updated.Actions) > 0 {
			if err := executor.AddRule(updated.RuleName, updated.TopicPattern, updated.SQL, actionsToList(updated.Actions), updated.ErrorAction); err != nil {
				slog.Warn("rule replaced but executor registration failed", "rule", in.RuleName, "error", err)
			}
		}
	}
	return nil
}

// deleteTopicRuleCore removes a rule, cleans up its tags, and deregisters
// it from the regional executor.
func (s *IoTService) deleteTopicRuleCore(store iotstore.IotStoreInterface, ruleName string) error {
	if ruleName == "" {
		return iotstore.ErrMissingParam
	}

	arn := iotstore.BuildRuleARN(store.GetAccountID(), store.GetRegion(), ruleName)
	_ = store.DeleteAllTags(arn)

	if err := store.DeleteRule(ruleName); err != nil {
		return err
	}

	if executor := s.ExecutorForRegion(store.GetRegion()); executor != nil {
		executor.RemoveRule(ruleName)
	}
	return nil
}

// enableTopicRuleCore re-enables a rule and re-registers it with the
// regional executor.
func (s *IoTService) enableTopicRuleCore(store iotstore.IotStoreInterface, ruleName string) error {
	if ruleName == "" {
		return iotstore.ErrMissingParam
	}

	disabled := false
	opts := iotstore.RuleUpdateOpts{RuleDisabled: &disabled}
	updated, err := store.UpdateRule(ruleName, opts)
	if err != nil {
		return err
	}

	if executor := s.ExecutorForRegion(store.GetRegion()); executor != nil && updated.SQL != "" && len(updated.Actions) > 0 {
		if err := executor.AddRule(updated.RuleName, updated.TopicPattern, updated.SQL, actionsToList(updated.Actions), updated.ErrorAction); err != nil {
			slog.Warn("rule enabled but executor registration failed", "rule", ruleName, "error", err)
		}
	}
	return nil
}

// disableTopicRuleCore disables a rule and deregisters it from the regional
// executor.
func (s *IoTService) disableTopicRuleCore(store iotstore.IotStoreInterface, ruleName string) error {
	if ruleName == "" {
		return iotstore.ErrMissingParam
	}

	disabled := true
	opts := iotstore.RuleUpdateOpts{RuleDisabled: &disabled}
	if _, err := store.UpdateRule(ruleName, opts); err != nil {
		return err
	}

	if executor := s.ExecutorForRegion(store.GetRegion()); executor != nil {
		executor.RemoveRule(ruleName)
	}
	return nil
}
