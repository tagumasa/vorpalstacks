package waf

import (
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/types"
)

const ipSetBucketName = "waf_ip_sets"

var ipSetAccessor = wafResourceAccessor[IPSet]{
	getIDFn:        func(r *IPSet) string { return r.ID },
	getNameFn:      func(r *IPSet) string { return r.Name },
	getARNFn:       func(r *IPSet) string { return r.ARN },
	setARNFn:       func(r *IPSet, arn string) { r.ARN = arn },
	getLockTokenFn: func(r *IPSet) string { return r.LockToken },
	setLockTokenFn: func(r *IPSet, lt string) { r.LockToken = lt },
	setModifiedFn:  func(r *IPSet) { r.ModifiedAt = time.Now() },
}

// IPSetStore provides storage for WAF IP Sets.
type IPSetStore struct {
	*ResourceStore[IPSet]
}

// NewIPSetStore creates a new IP Set store.
func NewIPSetStore(store storage.BasicStorage, accountId, region string) *IPSetStore {
	return &IPSetStore{
		ResourceStore: NewResourceStore[IPSet](store, ipSetBucketName, NewARNBuilder(accountId, region), ipSetAccessor),
	}
}

// Create creates a new IP Set.
func (s *IPSetStore) Create(id, name, description, ipAddressVersion string, addresses []string) (*IPSet, error) {
	if existing, _ := s.FindByName(name); existing != nil {
		return nil, ErrAlreadyExists
	}
	ipSet := &IPSet{
		ID:               id,
		Name:             name,
		Description:      description,
		IPAddressVersion: ipAddressVersion,
		Addresses:        addresses,
		Tags:             []types.Tag{},
		CreatedAt:        time.Now(),
		ModifiedAt:       time.Now(),
	}
	ipSet.ARN = s.arnBuilder.BuildIPSetARN(id)
	SetTimestamps(&ipSetAccessor, ipSet)
	if err := s.Put(id, ipSet, "create_ip_set"); err != nil {
		return nil, err
	}
	return ipSet, nil
}

// Update updates an existing IP Set.
func (s *IPSetStore) Update(id, lockToken string, addresses []string) (*IPSet, error) {
	return s.UpdateWithLockToken(id, lockToken, func(ipSet *IPSet) error {
		ipSet.Addresses = addresses
		return nil
	}, "update_ip_set")
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
