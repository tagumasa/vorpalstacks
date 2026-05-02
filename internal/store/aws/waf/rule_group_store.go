package waf

import (
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/types"
)

const ruleGroupBucketName = "waf_rule_groups"
const ruleKeyPrefix = "rule_"

var ruleGroupAccessor = wafResourceAccessor[RuleGroup]{
	getIDFn:        func(r *RuleGroup) string { return r.ID },
	getNameFn:      func(r *RuleGroup) string { return r.Name },
	getARNFn:       func(r *RuleGroup) string { return r.ARN },
	setARNFn:       func(r *RuleGroup, arn string) { r.ARN = arn },
	getLockTokenFn: func(r *RuleGroup) string { return r.LockToken },
	setLockTokenFn: func(r *RuleGroup, lt string) { r.LockToken = lt },
	setModifiedFn:  func(r *RuleGroup) { r.ModifiedAt = time.Now() },
}

// RuleGroupStore provides storage for WAF Rule Groups.
type RuleGroupStore struct {
	*ResourceStore[RuleGroup]
}

// NewRuleGroupStore creates a new Rule Group store.
func NewRuleGroupStore(store storage.BasicStorage, accountId, region string) *RuleGroupStore {
	return &RuleGroupStore{
		ResourceStore: NewResourceStore[RuleGroup](store, ruleGroupBucketName, NewARNBuilder(accountId, region), ruleGroupAccessor),
	}
}

// Create creates a new Rule Group.
func (s *RuleGroupStore) Create(id, name, description string, capacity int64, rules []*Rule, visibilityConfig *VisibilityConfig) (*RuleGroup, error) {
	if existing, _ := s.FindByName(name); existing != nil {
		return nil, ErrAlreadyExists
	}
	ruleGroup := &RuleGroup{
		ID:               id,
		Name:             name,
		Description:      description,
		Capacity:         capacity,
		Rules:            rules,
		VisibilityConfig: visibilityConfig,
		Tags:             []types.Tag{},
		CreatedAt:        time.Now(),
		ModifiedAt:       time.Now(),
	}
	ruleGroup.ARN = s.arnBuilder.BuildRuleGroupARN(id)
	SetTimestamps(&ruleGroupAccessor, ruleGroup)
	if err := s.Put(id, ruleGroup, "create_rule_group"); err != nil {
		return nil, err
	}
	return ruleGroup, nil
}

// Update updates an existing Rule Group.
func (s *RuleGroupStore) Update(id, lockToken string, capacity int64, rules []*Rule, visibilityConfig *VisibilityConfig) (*RuleGroup, error) {
	return s.UpdateWithLockToken(id, lockToken, func(ruleGroup *RuleGroup) error {
		ruleGroup.Capacity = capacity
		ruleGroup.Rules = rules
		ruleGroup.VisibilityConfig = visibilityConfig
		return nil
	}, "update_rule_group")
}

// List returns a paginated list of Rule Groups.
func (s *RuleGroupStore) List(marker string, maxItems int) (*RuleGroupListResult, error) {
	result, err := common.List[RuleGroup](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, nil)
	if err != nil {
		return nil, NewStoreError("list_rule_groups", err)
	}
	return &RuleGroupListResult{
		RuleGroups:  result.Items,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}

// CreateRule creates a standalone WAF Rule stored within the Rule Group bucket.
func (s *RuleGroupStore) CreateRule(id string, rule *Rule) error {
	key := ruleKeyPrefix + id
	rule.ID = id
	rule.ARN = s.arnBuilder.BuildRuleARN(id)
	return s.BaseStore.Put(key, rule)
}

// GetRule retrieves a standalone WAF Rule by its ID.
func (s *RuleGroupStore) GetRule(id string) (*Rule, error) {
	var rule Rule
	if err := s.BaseStore.Get(ruleKeyPrefix+id, &rule); err != nil {
		return nil, NewStoreError("get_rule", err)
	}
	return &rule, nil
}

// DeleteRule permanently deletes a standalone WAF Rule by its ID.
func (s *RuleGroupStore) DeleteRule(id string) error {
	return s.BaseStore.Delete(ruleKeyPrefix + id)
}

// ListRules returns a list of standalone WAF rules stored in the Rule Group bucket.
func (s *RuleGroupStore) ListRules(limit int) ([]*Rule, error) {
	if limit <= 0 {
		limit = 100
	}
	result, err := common.List[Rule](s.BaseStore, common.ListOptions{Prefix: ruleKeyPrefix, MaxItems: limit}, nil)
	if err != nil {
		return nil, NewStoreError("list_rules", err)
	}
	return result.Items, nil
}
