package iot

import (
	"vorpalstacks/internal/store/aws/common"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
)
func (s *IotStore) CreateRule(rule *TopicRule) (*TopicRule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if rule.RuleName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.TopicRule{}
	if err := s.rulesBase.GetProto(rule.RuleName, existing); err == nil {
		return nil, ErrRuleAlreadyExists
	}
	rule.ARN = BuildRuleARN(s.accountID, s.region, rule.RuleName)
	return rule, s.topicRulePS.Create(rule)
}

func (s *IotStore) GetRule(ruleName string) (*TopicRule, error) {
	return s.topicRulePS.Get(ruleName)
}

func (s *IotStore) DeleteRule(ruleName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.topicRulePS.DeleteIfExists(ruleName)
}

func (s *IotStore) ListRules(opts common.ListOptions) (*common.ListResult[TopicRule], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.rulesBase, opts, func() *pb.TopicRule { return &pb.TopicRule{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*TopicRule, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToRule(p))
	}
	return &common.ListResult[TopicRule]{Items: items, NextMarker: result.NextMarker}, nil
}

func (s *IotStore) UpdateRule(ruleName string, opts RuleUpdateOpts) (*TopicRule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, err := s.topicRulePS.Get(ruleName)
	if err != nil {
		return nil, ErrRuleNotFound
	}
	if opts.SQL != "" {
		existing.SQL = opts.SQL
	}
	if opts.Description != "" {
		existing.Description = opts.Description
	}
	if opts.AwsIotSqlVersion != "" {
		existing.AwsIotSqlVersion = opts.AwsIotSqlVersion
	}
	if opts.RuleDisabled != nil {
		existing.RuleDisabled = *opts.RuleDisabled
	}
	if opts.Actions != nil {
		existing.Actions = opts.Actions
	}
	if opts.ErrorAction != nil {
		existing.ErrorAction = opts.ErrorAction
	}
	return existing, s.topicRulePS.Update(existing)
}

// UpdateThingType persists changes to an existing thing type.
