package cloudwatchlogs

import (
	"errors"
	"sync"
	"time"
)

// --- Resource Policy ---

func (s *Store) PutResourcePolicy(policy *ResourcePolicy) error {
	policy.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.resourcePolicyKey(policy.PolicyName), policy)
}

func (s *Store) GetResourcePolicy(policyName string) (*ResourcePolicy, error) {
	return getJSONRecord[ResourcePolicy](s, s.resourcePolicyKey(policyName), ErrResourceNotFound)
}

// resourcePolicyMu serialises the resource-policy admission cycle. The
// optimistic-revision guard and the documented policy census ("An account
// can have a maximum of 10 policies without resourceARN and one per
// LogGroup resourceARN") both decide on reads the bare PutResourcePolicy
// write never checks, so the read, the admission checks and the write
// must commit under one mutex: concurrent Puts quoting the same
// expectedRevisionId admit exactly one, and the census cannot be raced
// past its ceiling.
var resourcePolicyMu sync.Mutex

// MutateResourcePolicy commits one resource-policy write atomically: fn
// receives the persisted record for the name, or nil when none exists,
// and returns the record to persist alongside a nil error — an error
// aborts with no write, and a nil record writes nothing. The storage
// read error surfaces instead of masquerading as "no existing policy",
// so an outage cannot silently skip the revision requirement.
func (s *Store) MutateResourcePolicy(policyName string, fn func(existing *ResourcePolicy) (*ResourcePolicy, error)) error {
	resourcePolicyMu.Lock()
	defer resourcePolicyMu.Unlock()
	existing, err := s.GetResourcePolicy(policyName)
	if err != nil && !errors.Is(err, ErrResourceNotFound) {
		return err
	}
	rp, err := fn(existing)
	if err != nil {
		return err
	}
	if rp == nil {
		return nil
	}
	return s.PutResourcePolicy(rp)
}

// DeleteResourcePolicyAdmitted removes the named policy after fn admits
// the delete against the persisted record. The revision guard reads the
// stored record and the deletion commits under the admission mutex, so a
// concurrent Put cannot slip a freshly minted revision between the guard
// and the delete.
func (s *Store) DeleteResourcePolicyAdmitted(policyName string, fn func(existing *ResourcePolicy) error) error {
	resourcePolicyMu.Lock()
	defer resourcePolicyMu.Unlock()
	existing, err := s.GetResourcePolicy(policyName)
	if err != nil {
		return err
	}
	if err := fn(existing); err != nil {
		return err
	}
	return s.DeleteResourcePolicy(policyName)
}

func (s *Store) DeleteResourcePolicy(policyName string) error {
	if !s.Exists(s.resourcePolicyKey(policyName)) {
		return ErrResourceNotFound
	}
	return s.Delete(s.resourcePolicyKey(policyName))
}

func (s *Store) ListResourcePolicies(resourceArn string) ([]*ResourcePolicy, error) {
	return listJSONRecords[ResourcePolicy](s, keyPrefixResourcePolicy, "resource policy record", func(p *ResourcePolicy) bool {
		return resourceArn == "" || p.ResourceArn == resourceArn
	})
}

// --- Query-result KMS association ---

// QueryResultKmsKey records the account-level AssociateKmsKey association
// for stored query results (the resourceIdentifier's query-result:* ARN
// form). One record per region store — the ARN is region-and-account
// scoped, and the platform's stores are the region boundary.
type QueryResultKmsKey struct {
	KmsKeyId        string `json:"kmsKeyId"`
	LastUpdatedTime int64  `json:"lastUpdatedTime"`
}

// PutQueryResultKmsKey records (or replaces) the query-result key
// association.
func (s *Store) PutQueryResultKmsKey(kmsKeyId string) error {
	return s.Put(keyPrefixQueryResultKms+"default", &QueryResultKmsKey{
		KmsKeyId:        kmsKeyId,
		LastUpdatedTime: time.Now().UTC().UnixMilli(),
	})
}

// GetQueryResultKmsKey returns the recorded association, or
// ErrResourceNotFound when no key is associated.
func (s *Store) GetQueryResultKmsKey() (*QueryResultKmsKey, error) {
	return getJSONRecord[QueryResultKmsKey](s, keyPrefixQueryResultKms+"default", ErrResourceNotFound)
}

// DeleteQueryResultKmsKey removes the association; a missing record
// deletes nothing (the disassociate path clears like the log-group form).
func (s *Store) DeleteQueryResultKmsKey() error {
	return s.Delete(keyPrefixQueryResultKms + "default")
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
	scanPrefix := keyPrefixAccountPolicy
	if policyType != "" {
		scanPrefix = keyPrefixAccountPolicy + escapePath(policyType) + ":"
	}
	return listJSONRecords[AccountPolicy](s, scanPrefix, "account policy record", func(p *AccountPolicy) bool {
		return policyName == "" || p.PolicyName == policyName
	})
}

// --- Data Protection Policy ---

// PutDataProtectionPolicy stores the group's data protection policy. The
// write holds the group's write lock — the DeleteLogGroup teardown lock —
// and re-checks the group inside it: the service layer's earlier group
// check runs outside any critical section, so without the re-check a
// teardown interleaving this put would leave the policy permanently
// outliving its deleted group (the metric-filter family's pattern).
func (s *Store) PutDataProtectionPolicy(policy *DataProtectionPolicy) error {
	gLock := s.groupLock(policy.LogGroupIdentifier)
	gLock.Lock()
	defer gLock.Unlock()

	if _, err := s.GetLogGroup(policy.LogGroupIdentifier); err != nil {
		return err
	}
	policy.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.dataProtectionPolicyKey(policy.LogGroupIdentifier), policy)
}

func (s *Store) GetDataProtectionPolicy(logGroupIdentifier string) (*DataProtectionPolicy, error) {
	return getJSONRecord[DataProtectionPolicy](s, s.dataProtectionPolicyKey(logGroupIdentifier), ErrResourceNotFound)
}

func (s *Store) DeleteDataProtectionPolicy(logGroupIdentifier string) error {
	key := s.dataProtectionPolicyKey(logGroupIdentifier)
	if !s.Exists(key) {
		return ErrResourceNotFound
	}
	return s.Delete(key)
}

// ResourcePolicy represents a CloudWatch Logs resource policy.
type ResourcePolicy struct {
	PolicyName      string `json:"policyName"`
	PolicyDocument  string `json:"policyDocument"`
	ResourceArn     string `json:"resourceArn,omitempty"`
	PolicyScope     string `json:"policyScope,omitempty"`
	RevisionId      string `json:"revisionId,omitempty"`
	LastUpdatedTime int64  `json:"lastUpdatedTime"`
}

// AccountPolicy represents a CloudWatch Logs account-level policy.
type AccountPolicy struct {
	PolicyName        string `json:"policyName"`
	PolicyDocument    string `json:"policyDocument"`
	PolicyType        string `json:"policyType"`
	Scope             string `json:"scope,omitempty"`
	SelectionCriteria string `json:"selectionCriteria,omitempty"`
	AccountId         string `json:"accountId,omitempty"`
	LastUpdatedTime   int64  `json:"lastUpdatedTime"`
}

// DataProtectionPolicy represents a CloudWatch Logs data protection policy.
type DataProtectionPolicy struct {
	LogGroupIdentifier string `json:"logGroupIdentifier"`
	PolicyDocument     string `json:"policyDocument"`
	LastUpdatedTime    int64  `json:"lastUpdatedTime"`
}

// StorageTierPolicy is the account-level storage tier policy: the tier
// the account's log data is held in (the StorageTier enum, STANDARD or
// INTELLIGENT_TIERING) and the stamp of its last update. The tiering has
// no API-observable output surface of its own — every read operation
// serves events the same way — so the policy's faithful platform form is
// the stored round-trip.
type StorageTierPolicy struct {
	StorageTier     string `json:"storageTier"`
	LastUpdatedTime int64  `json:"lastUpdatedTime"`
}
