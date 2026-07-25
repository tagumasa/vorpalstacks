package cloudwatch

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

func insightRuleBucketName(region string) string {
	return "cw_insight_rules-" + region
}

// InsightRuleStore provides storage operations for CloudWatch
// Contributor Insights rules.
type InsightRuleStore struct {
	*common.BaseStore
	mu sync.Mutex
}

// NewInsightRuleStore creates a new InsightRuleStore instance.
func NewInsightRuleStore(store storage.BasicStorage, region string) *InsightRuleStore {
	return &InsightRuleStore{
		BaseStore: common.NewBaseStore(store.Bucket(insightRuleBucketName(region)), "cloudwatch-insight-rules"),
	}
}

// PutInsightRule creates or updates an insight rule.
func (s *InsightRuleStore) PutInsightRule(rule *InsightRule) (*InsightRule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := insightRuleKey(rule.Name)

	existing := &InsightRule{}
	if err := s.BaseStore.Get(key, existing); err == nil {
		existing.State = rule.State
		existing.Definition = rule.Definition
		existing.Schema = rule.Schema
		existing.ApplyOnTransformedLogs = rule.ApplyOnTransformedLogs
		if rule.Tags != nil {
			existing.Tags = rule.Tags
		}
		existing.UpdatedAt = time.Now().UTC()
		if err := s.BaseStore.Put(key, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	now := time.Now().UTC()
	if rule.State == "" {
		rule.State = "ENABLED"
	}
	rule.CreatedAt = now
	rule.UpdatedAt = now
	if err := s.BaseStore.Put(key, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

// DeleteInsightRules deletes rules by name. Returns successfully deleted
// and not-found name lists.
func (s *InsightRuleStore) DeleteInsightRules(names []string) ([]string, []string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var deleted, notFound []string
	for _, name := range names {
		key := insightRuleKey(name)
		if !s.BaseStore.Exists(key) {
			notFound = append(notFound, name)
			continue
		}
		if err := s.BaseStore.Delete(key); err != nil {
			return deleted, notFound, err
		}
		deleted = append(deleted, name)
	}
	return deleted, notFound, nil
}

// SetRuleState enables or disables the named rules.
func (s *InsightRuleStore) SetRuleState(names []string, state string) ([]string, []string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var updated, notFound []string
	for _, name := range names {
		key := insightRuleKey(name)
		rule := &InsightRule{}
		if err := s.BaseStore.Get(key, rule); err != nil {
			notFound = append(notFound, name)
			continue
		}
		rule.State = state
		rule.UpdatedAt = time.Now().UTC()
		if err := s.BaseStore.Put(key, rule); err != nil {
			return updated, notFound, err
		}
		updated = append(updated, name)
	}
	return updated, notFound, nil
}

// ListInsightRules returns all insight rules, optionally filtered by
// managed status.
func (s *InsightRuleStore) ListInsightRules(managedOnly bool) ([]*InsightRule, error) {
	var rules []*InsightRule
	err := s.BaseStore.ScanPrefix("insight_rule:", func(key string, value []byte) error {
		var rule InsightRule
		if err := json.Unmarshal(value, &rule); err != nil {
			return err
		}
		if managedOnly && !rule.ManagedRule {
			return nil
		}
		rules = append(rules, &rule)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// ListInsightRulesPaginated returns a paginated list of insight rules,
// optionally filtered by managed status.
func (s *InsightRuleStore) ListInsightRulesPaginated(managedOnly bool, opts common.ListOptions) (*common.ListResult[InsightRule], error) {
	opts.Prefix = "insight_rule:"
	var filter func(*InsightRule) bool
	if managedOnly {
		filter = func(r *InsightRule) bool { return r.ManagedRule }
	}
	return common.List[InsightRule](s.BaseStore, opts, filter)
}

// GetInsightRule returns a rule by name.
func (s *InsightRuleStore) GetInsightRule(name string) (*InsightRule, error) {
	var rule InsightRule
	key := insightRuleKey(name)
	if err := s.BaseStore.Get(key, &rule); err != nil {
		return nil, fmt.Errorf("insight rule not found: %s", name)
	}
	return &rule, nil
}

func insightRuleKey(name string) string {
	return "insight_rule:" + name
}
