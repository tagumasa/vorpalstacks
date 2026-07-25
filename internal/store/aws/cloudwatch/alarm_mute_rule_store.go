package cloudwatch

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

func alarmMuteRuleBucketName(region string) string {
	return "cw_alarm_mute_rules-" + region
}

// AlarmMuteRuleStore provides storage operations for CloudWatch alarm
// mute rules.
type AlarmMuteRuleStore struct {
	*common.BaseStore
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
	mu         sync.Mutex
}

// NewAlarmMuteRuleStore creates a new AlarmMuteRuleStore instance.
func NewAlarmMuteRuleStore(store storage.BasicStorage, accountID, region string) *AlarmMuteRuleStore {
	return &AlarmMuteRuleStore{
		BaseStore:  common.NewBaseStore(store.Bucket(alarmMuteRuleBucketName(region)), "cloudwatch-alarm-mute-rules"),
		arnBuilder: svcarn.NewARNBuilder(accountID, region),
		accountID:  accountID,
		region:     region,
	}
}

// PutAlarmMuteRule creates or updates an alarm mute rule.
func (s *AlarmMuteRuleStore) PutAlarmMuteRule(rule *AlarmMuteRule) (*AlarmMuteRule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := alarmMuteRuleKey(rule.Name)
	now := time.Now().UTC()

	existing := &AlarmMuteRule{}
	if err := s.BaseStore.Get(key, existing); err == nil {
		existing.Description = rule.Description
		existing.ScheduleExpr = rule.ScheduleExpr
		existing.MutedAlarmNames = rule.MutedAlarmNames
		existing.StartDate = rule.StartDate
		existing.ExpireDate = rule.ExpireDate
		existing.MuteType = rule.MuteType
		if rule.Tags != nil {
			existing.Tags = rule.Tags
		}
		existing.Status = computeMuteRuleStatus(rule.StartDate, rule.ExpireDate, now)
		existing.UpdatedAt = now
		if err := s.BaseStore.Put(key, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	rule.ARN = s.arnBuilder.Build("cloudwatch", "mute-rule:"+rule.Name)
	rule.Status = computeMuteRuleStatus(rule.StartDate, rule.ExpireDate, now)
	if rule.MuteType == "" {
		rule.MuteType = "AUTOMATIC"
	}
	rule.CreatedAt = now
	rule.UpdatedAt = now
	if err := s.BaseStore.Put(key, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

// DeleteAlarmMuteRule deletes a rule by name.
func (s *AlarmMuteRuleStore) DeleteAlarmMuteRule(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := alarmMuteRuleKey(name)
	if !s.BaseStore.Exists(key) {
		return fmt.Errorf("%w: %s", ErrAlarmMuteRuleNotFound, name)
	}
	return s.BaseStore.Delete(key)
}

// ErrAlarmMuteRuleNotFound is returned when a mute rule does not exist.
var ErrAlarmMuteRuleNotFound = fmt.Errorf("alarm mute rule not found")

// GetAlarmMuteRule returns a rule by name.
func (s *AlarmMuteRuleStore) GetAlarmMuteRule(name string) (*AlarmMuteRule, error) {
	var rule AlarmMuteRule
	key := alarmMuteRuleKey(name)
	if err := s.BaseStore.Get(key, &rule); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrAlarmMuteRuleNotFound, name)
	}
	// Recompute status so callers always see the current value.
	rule.Status = computeMuteRuleStatus(rule.StartDate, rule.ExpireDate, time.Now().UTC())
	return &rule, nil
}

// ListAlarmMuteRules returns all alarm mute rules, optionally filtered
// by alarm name and statuses.  Status is recomputed at query time so
// that expired rules are correctly classified without a background
// reaper.
func (s *AlarmMuteRuleStore) ListAlarmMuteRules(alarmName string, statuses []string) ([]*AlarmMuteRule, error) {
	var rules []*AlarmMuteRule
	now := time.Now().UTC()
	err := s.BaseStore.ScanPrefix("mute_rule:", func(key string, value []byte) error {
		var rule AlarmMuteRule
		if err := json.Unmarshal(value, &rule); err != nil {
			return err
		}

		// Recompute status before filtering so time-based transitions
		// (SCHEDULED → ACTIVE → EXPIRED) are always accurate.
		rule.Status = computeMuteRuleStatus(rule.StartDate, rule.ExpireDate, now)

		// Filter by alarm name if provided.
		if alarmName != "" {
			found := false
			for _, n := range rule.MutedAlarmNames {
				if n == alarmName {
					found = true
					break
				}
			}
			if !found {
				return nil
			}
		}

		// Filter by statuses if provided.
		if len(statuses) > 0 {
			matched := false
			for _, st := range statuses {
				if rule.Status == st {
					matched = true
					break
				}
			}
			if !matched {
				return nil
			}
		}

		rules = append(rules, &rule)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// ListAlarmMuteRulesPaginated returns a paginated list of alarm mute
// rules, optionally filtered by alarm name and statuses.
func (s *AlarmMuteRuleStore) ListAlarmMuteRulesPaginated(alarmName string, statuses []string, opts common.ListOptions) (*common.ListResult[AlarmMuteRule], error) {
	opts.Prefix = "mute_rule:"
	now := time.Now().UTC()

	var filter func(*AlarmMuteRule) bool
	if alarmName != "" || len(statuses) > 0 {
		filter = func(r *AlarmMuteRule) bool {
			r.Status = computeMuteRuleStatus(r.StartDate, r.ExpireDate, now)
			if alarmName != "" {
				found := false
				for _, n := range r.MutedAlarmNames {
					if n == alarmName {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
			if len(statuses) > 0 {
				matched := false
				for _, st := range statuses {
					if r.Status == st {
						matched = true
						break
					}
				}
				if !matched {
					return false
				}
			}
			return true
		}
	} else {
		filter = func(r *AlarmMuteRule) bool {
			r.Status = computeMuteRuleStatus(r.StartDate, r.ExpireDate, now)
			return true
		}
	}
	return common.List[AlarmMuteRule](s.BaseStore, opts, filter)
}

// IsAlarmMuted checks if the given alarm name is muted by any ACTIVE
// alarm mute rule.
func (s *AlarmMuteRuleStore) IsAlarmMuted(alarmName string) bool {
	rules, err := s.ListAlarmMuteRules(alarmName, []string{"ACTIVE"})
	if err != nil {
		return false
	}
	return len(rules) > 0
}

// computeMuteRuleStatus determines the status of a mute rule based on
// the current time relative to the start and expiry dates.
func computeMuteRuleStatus(start, expire, now time.Time) string {
	if !start.IsZero() && now.Before(start) {
		return "SCHEDULED"
	}
	if !expire.IsZero() && now.After(expire) {
		return "EXPIRED"
	}
	return "ACTIVE"
}

func alarmMuteRuleKey(name string) string {
	return "mute_rule:" + name
}
