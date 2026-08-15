package cloudwatchlogs

import (
	"encoding/json"
	"time"
)

// --- Resource Policy ---

func (s *Store) PutResourcePolicy(policy *ResourcePolicy) error {
	policy.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.resourcePolicyKey(policy.PolicyName), policy)
}

func (s *Store) GetResourcePolicy(policyName string) (*ResourcePolicy, error) {
	var p ResourcePolicy
	if err := s.Get(s.resourcePolicyKey(policyName), &p); err != nil {
		return nil, ErrResourceNotFound
	}
	return &p, nil
}

func (s *Store) DeleteResourcePolicy(policyName string) error {
	if !s.Exists(s.resourcePolicyKey(policyName)) {
		return ErrResourceNotFound
	}
	return s.Delete(s.resourcePolicyKey(policyName))
}

func (s *Store) ListResourcePolicies(resourceArn string) ([]*ResourcePolicy, error) {
	var policies []*ResourcePolicy
	if err := s.ScanPrefix("resource-policy:", func(key string, value []byte) error {
		var p ResourcePolicy
		if err := json.Unmarshal(value, &p); err != nil {
			return nil
		}
		if resourceArn == "" || p.ResourceArn == resourceArn {
			policies = append(policies, &p)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return policies, nil
}

// --- Account Policy ---

func (s *Store) PutAccountPolicy(policy *AccountPolicy) error {
	policy.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.accountPolicyKey(policy.PolicyType, policy.PolicyName), policy)
}

func (s *Store) DeleteAccountPolicyEntry(policyType, policyName string) error {
	key := s.accountPolicyKey(policyType, policyName)
	if !s.Exists(key) {
		return ErrResourceNotFound
	}
	return s.Delete(key)
}

func (s *Store) ListAccountPolicies(policyType, policyName string) ([]*AccountPolicy, error) {
	var policies []*AccountPolicy
	scanPrefix := "account-policy:"
	if policyType != "" {
		scanPrefix = "account-policy:" + policyType + ":"
	}
	if err := s.ScanPrefix(scanPrefix, func(key string, value []byte) error {
		var p AccountPolicy
		if err := json.Unmarshal(value, &p); err != nil {
			return nil
		}
		if policyName != "" && p.PolicyName != policyName {
			return nil
		}
		policies = append(policies, &p)
		return nil
	}); err != nil {
		return nil, err
	}
	return policies, nil
}

// --- Data Protection Policy ---

func (s *Store) PutDataProtectionPolicy(policy *DataProtectionPolicy) error {
	policy.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.dataProtectionPolicyKey(policy.LogGroupIdentifier), policy)
}

func (s *Store) GetDataProtectionPolicy(logGroupIdentifier string) (*DataProtectionPolicy, error) {
	var p DataProtectionPolicy
	if err := s.Get(s.dataProtectionPolicyKey(logGroupIdentifier), &p); err != nil {
		return nil, ErrResourceNotFound
	}
	return &p, nil
}

func (s *Store) DeleteDataProtectionPolicy(logGroupIdentifier string) error {
	key := s.dataProtectionPolicyKey(logGroupIdentifier)
	if !s.Exists(key) {
		return ErrResourceNotFound
	}
	return s.Delete(key)
}
