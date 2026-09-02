package waf

import (
	"time"

	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const ruleGroupBucketName = "waf_rule_groups"
const ruleKeyPrefix = "rule_"

var ruleGroupAccessor = wafResourceAccessor[RuleGroup]{
	getIDFn:        func(r *RuleGroup) string { return r.ID },
	getNameFn:      func(r *RuleGroup) string { return r.Name },
	getScopeFn:     func(r *RuleGroup) string { return r.Scope },
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
func (s *RuleGroupStore) Create(ruleGroup *RuleGroup) (*RuleGroup, error) {
	if existing, _ := s.FindByNameAndScope(ruleGroup.Name, ruleGroup.Scope); existing != nil {
		return nil, ErrAlreadyExists
	}
	if ruleGroup.Scope == "" {
		ruleGroup.Scope = "REGIONAL"
	}
	if ruleGroup.Tags == nil {
		ruleGroup.Tags = []types.Tag{}
	}
	ruleGroup.ARN = s.arnBuilder.BuildRuleGroupARN(ruleGroup.Name, ruleGroup.ID, ruleGroup.Scope)
	SetTimestamps(&ruleGroupAccessor, ruleGroup)
	if err := s.Put(ruleGroup.ID, ruleGroup, "create_rule_group"); err != nil {
		return nil, err
	}
	return ruleGroup, nil
}

// Update updates an existing Rule Group. Capacity is not a parameter:
// UpdateRuleGroupRequest has no Capacity member in the Smithy model, so
// a rule group's capacity is fixed at creation and only the rules and
// visibility configuration are mutable.
func (s *RuleGroupStore) Update(id, lockToken string, rules []*Rule, visibilityConfig *VisibilityConfig) (*RuleGroup, error) {
	return s.UpdateWithLockToken(id, lockToken, func(ruleGroup *RuleGroup) error {
		ruleGroup.Rules = rules
		ruleGroup.VisibilityConfig = visibilityConfig
		return nil
	}, "update_rule_group")
}

// List returns a paginated list of Rule Groups filtered by scope.
func (s *RuleGroupStore) List(marker string, maxItems int, scope string) (*RuleGroupListResult, error) {
	var filter common.FilterFunc[RuleGroup]
	if scope != "" {
		filter = func(rg *RuleGroup) bool { return rg.Scope == scope }
	}
	result, err := common.List[RuleGroup](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, filter)
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
