package waf

// Package waf provides WAF (Web Application Firewall) data store implementations
// for vorpalstacks.

import (
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const regexPatternSetBucketName = "waf_regex_pattern_sets"

// RegexPatternSetStore provides storage for WAF Regex Pattern Sets.
type RegexPatternSetStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	mu         sync.Mutex
}

// NewRegexPatternSetStore creates a new Regex Pattern Set store.
func NewRegexPatternSetStore(store storage.BasicStorage, accountId, region string) *RegexPatternSetStore {
	return &RegexPatternSetStore{
		BaseStore:  common.NewBaseStore(store.Bucket(regexPatternSetBucketName), "waf"),
		arnBuilder: NewARNBuilder(accountId, region),
	}
}

// Get retrieves a Regex Pattern Set by its ID.
func (s *RegexPatternSetStore) Get(id string) (*RegexPatternSet, error) {
	var regexPatternSet RegexPatternSet
	if err := s.BaseStore.Get(id, &regexPatternSet); err != nil {
		return nil, NewStoreError("get_regex_pattern_set", err)
	}
	return &regexPatternSet, nil
}

// GetByARN retrieves a Regex Pattern Set by its ARN.
func (s *RegexPatternSetStore) GetByARN(arn string) (*RegexPatternSet, error) {
	return common.FindFirst[RegexPatternSet](s.BaseStore, func(r *RegexPatternSet) bool { return r.ARN == arn })
}

// Create creates a new Regex Pattern Set.
func (s *RegexPatternSetStore) Create(id, name, description string, regularPatterns []string) (*RegexPatternSet, error) {
	regexPatternSet := &RegexPatternSet{
		ID:              id,
		Name:            name,
		Description:     description,
		RegularPatterns: regularPatterns,
		ARN:             s.arnBuilder.BuildRegexPatternSetARN(id),
		LockToken:       GenerateLockToken(),
		CreatedAt:       time.Now(),
		ModifiedAt:      time.Now(),
	}

	if err := s.BaseStore.Put(id, regexPatternSet); err != nil {
		return nil, NewStoreError("create_regex_pattern_set", err)
	}
	return regexPatternSet, nil
}

// Update updates an existing Regex Pattern Set.
func (s *RegexPatternSetStore) Update(id, lockToken string, regularPatterns []string) (*RegexPatternSet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	regexPatternSet, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if regexPatternSet.LockToken != lockToken {
		return nil, NewStoreError("update_regex_pattern_set", ErrLockTokenMismatch)
	}

	regexPatternSet.RegularPatterns = regularPatterns
	regexPatternSet.ModifiedAt = time.Now()
	regexPatternSet.LockToken = GenerateLockToken()

	if err := s.BaseStore.Put(id, regexPatternSet); err != nil {
		return nil, NewStoreError("update_regex_pattern_set", err)
	}
	return regexPatternSet, nil
}

// Delete deletes a Regex Pattern Set.
func (s *RegexPatternSetStore) Delete(id, lockToken string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	regexPatternSet, err := s.Get(id)
	if err != nil {
		return err
	}

	if lockToken != "" && regexPatternSet.LockToken != lockToken {
		return NewStoreError("delete_regex_pattern_set", ErrLockTokenMismatch)
	}

	if err := s.BaseStore.Delete(id); err != nil {
		return NewStoreError("delete_regex_pattern_set", err)
	}
	return nil
}

// List returns a paginated list of Regex Pattern Sets.
func (s *RegexPatternSetStore) List(marker string, maxItems int) (*RegexPatternSetListResult, error) {
	result, err := common.List[RegexPatternSet](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, nil)
	if err != nil {
		return nil, NewStoreError("list_regex_pattern_sets", err)
	}
	return &RegexPatternSetListResult{
		RegexPatternSets: result.Items,
		IsTruncated:      result.IsTruncated,
		NextMarker:       result.NextMarker,
	}, nil
}
