package waf

// Package waf provides WAF (Web Application Firewall) data store implementations
// for vorpalstacks.

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const ruleGroupBucketName = "waf_rule_groups"
const ruleKeyPrefix = "rule_"

// RuleGroupStore provides storage for WAF Rule Groups.
type RuleGroupStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
}

// NewRuleGroupStore creates a new Rule Group store.
func NewRuleGroupStore(store storage.BasicStorage, accountId, region string) *RuleGroupStore {
	return &RuleGroupStore{
		BaseStore:  common.NewBaseStore(store.Bucket(ruleGroupBucketName), "waf"),
		arnBuilder: NewARNBuilder(accountId, region),
	}
}

// Get retrieves a Rule Group by its ID.
func (s *RuleGroupStore) Get(id string) (*RuleGroup, error) {
	var ruleGroup RuleGroup
	if err := s.BaseStore.Get(id, &ruleGroup); err != nil {
		return nil, NewStoreError("get_rule_group", err)
	}
	return &ruleGroup, nil
}

// GetByARN retrieves a Rule Group by its ARN.
func (s *RuleGroupStore) GetByARN(arn string) (*RuleGroup, error) {
	var found *RuleGroup
	err := s.ForEach(func(key string, value []byte) error {
		var rg RuleGroup
		if err := json.Unmarshal(value, &rg); err != nil {
			return err
		}
		if rg.ARN == arn {
			found = &rg
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("get_rule_group_by_arn", err)
	}
	if found == nil {
		return nil, NewStoreError("get_rule_group_by_arn", ErrNotFound)
	}
	return found, nil
}

// Create creates a new Rule Group.
func (s *RuleGroupStore) Create(id, name, description string, capacity int64, rules []*Rule, visibilityConfig *VisibilityConfig) (*RuleGroup, error) {
	ruleGroup := &RuleGroup{
		ID:               id,
		Name:             name,
		Description:      description,
		Capacity:         capacity,
		Rules:            rules,
		VisibilityConfig: visibilityConfig,
		ARN:              s.arnBuilder.BuildRuleGroupARN(id),
		LockToken:        GenerateLockToken(),
		CreatedAt:        time.Now(),
		ModifiedAt:       time.Now(),
	}

	if err := s.BaseStore.Put(id, ruleGroup); err != nil {
		return nil, NewStoreError("create_rule_group", err)
	}
	return ruleGroup, nil
}

// Update updates an existing Rule Group.
func (s *RuleGroupStore) Update(id, lockToken string, capacity int64, rules []*Rule, visibilityConfig *VisibilityConfig) (*RuleGroup, error) {
	ruleGroup, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if lockToken != "" && ruleGroup.LockToken != lockToken {
		return nil, NewStoreError("delete_rule_group", ErrLockTokenMismatch)
	}

	ruleGroup.Capacity = capacity
	ruleGroup.Rules = rules
	ruleGroup.VisibilityConfig = visibilityConfig
	ruleGroup.ModifiedAt = time.Now()
	ruleGroup.LockToken = GenerateLockToken()

	if err := s.BaseStore.Put(id, ruleGroup); err != nil {
		return nil, NewStoreError("update_rule_group", err)
	}
	return ruleGroup, nil
}

// Delete deletes a Rule Group.
func (s *RuleGroupStore) Delete(id, lockToken string) error {
	ruleGroup, err := s.Get(id)
	if err != nil {
		return err
	}

	if ruleGroup.LockToken != lockToken {
		return NewStoreError("delete_rule_group", ErrLockTokenMismatch)
	}

	if err := s.BaseStore.Delete(id); err != nil {
		return NewStoreError("delete_rule_group", err)
	}
	return nil
}

// List returns a paginated list of Rule Groups.
func (s *RuleGroupStore) List(marker string, maxItems int) (*RuleGroupListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var ruleGroups []*RuleGroup
	count := 0
	started := marker == ""
	hasMore := false
	var lastKey string

	err := s.ForEach(func(key string, value []byte) error {
		if len(key) > len(ruleKeyPrefix) && key[:len(ruleKeyPrefix)] == ruleKeyPrefix {
			return nil
		}
		var rg RuleGroup
		if err := json.Unmarshal(value, &rg); err != nil {
			return err
		}

		if !started {
			if rg.ID == marker {
				started = true
			}
			return nil
		}

		if count < maxItems {
			ruleGroups = append(ruleGroups, &rg)
			lastKey = key
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_rule_groups", err)
	}

	var nextMarker string
	if hasMore {
		nextMarker = lastKey
	}

	return &RuleGroupListResult{
		RuleGroups:  ruleGroups,
		IsTruncated: hasMore,
		NextMarker:  nextMarker,
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
