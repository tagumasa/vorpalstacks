package waf

// Package waf provides WAF (Web Application Firewall) data store implementations
// for vorpalstacks.

import (
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const ipSetBucketName = "waf_ip_sets"

// IPSetStore provides storage for WAF IP Sets.
type IPSetStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	mu         sync.Mutex
}

// NewIPSetStore creates a new IP Set store.
func NewIPSetStore(store storage.BasicStorage, accountId, region string) *IPSetStore {
	return &IPSetStore{
		BaseStore:  common.NewBaseStore(store.Bucket(ipSetBucketName), "waf"),
		arnBuilder: NewARNBuilder(accountId, region),
	}
}

// Get retrieves an IP Set by its ID.
func (s *IPSetStore) Get(id string) (*IPSet, error) {
	var ipSet IPSet
	if err := s.BaseStore.Get(id, &ipSet); err != nil {
		return nil, NewStoreError("get_ip_set", err)
	}
	return &ipSet, nil
}

// GetByARN retrieves an IP Set by its ARN.
func (s *IPSetStore) GetByARN(arn string) (*IPSet, error) {
	return common.FindFirst[IPSet](s.BaseStore, func(i *IPSet) bool { return i.ARN == arn })
}

// Create creates a new IP Set.
func (s *IPSetStore) Create(id, name, description, ipAddressVersion string, addresses []string) (*IPSet, error) {
	ipSet := &IPSet{
		ID:               id,
		Name:             name,
		Description:      description,
		IPAddressVersion: ipAddressVersion,
		Addresses:        addresses,
		ARN:              s.arnBuilder.BuildIPSetARN(id),
		LockToken:        GenerateLockToken(),
		CreatedAt:        time.Now(),
		ModifiedAt:       time.Now(),
	}

	if err := s.BaseStore.Put(id, ipSet); err != nil {
		return nil, NewStoreError("create_ip_set", err)
	}
	return ipSet, nil
}

// Update updates an existing IP Set.
func (s *IPSetStore) Update(id, lockToken string, addresses []string) (*IPSet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ipSet, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if ipSet.LockToken != lockToken {
		return nil, NewStoreError("update_ip_set", ErrLockTokenMismatch)
	}

	ipSet.Addresses = addresses
	ipSet.ModifiedAt = time.Now()
	ipSet.LockToken = GenerateLockToken()

	if err := s.BaseStore.Put(id, ipSet); err != nil {
		return nil, NewStoreError("update_ip_set", err)
	}
	return ipSet, nil
}

// Delete deletes an IP Set.
func (s *IPSetStore) Delete(id, lockToken string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ipSet, err := s.Get(id)
	if err != nil {
		return err
	}

	if lockToken != "" && ipSet.LockToken != lockToken {
		return NewStoreError("delete_ip_set", ErrLockTokenMismatch)
	}

	if err := s.BaseStore.Delete(id); err != nil {
		return NewStoreError("delete_ip_set", err)
	}
	return nil
}

// List returns a paginated list of IP Sets.
func (s *IPSetStore) List(marker string, maxItems int) (*IPSetListResult, error) {
	result, err := common.List[IPSet](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, nil)
	if err != nil {
		return nil, NewStoreError("list_ip_sets", err)
	}
	return &IPSetListResult{
		IPSets:      result.Items,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}
