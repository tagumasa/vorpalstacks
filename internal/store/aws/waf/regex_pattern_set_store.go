package waf

import (
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/types"
)

const regexPatternSetBucketName = "waf_regex_pattern_sets"

var regexPatternSetAccessor = wafResourceAccessor[RegexPatternSet]{
	getIDFn:        func(r *RegexPatternSet) string { return r.ID },
	getARNFn:       func(r *RegexPatternSet) string { return r.ARN },
	setARNFn:       func(r *RegexPatternSet, arn string) { r.ARN = arn },
	getLockTokenFn: func(r *RegexPatternSet) string { return r.LockToken },
	setLockTokenFn: func(r *RegexPatternSet, lt string) { r.LockToken = lt },
	setModifiedFn:  func(r *RegexPatternSet) { r.ModifiedAt = time.Now() },
}

// RegexPatternSetStore provides storage for WAF Regex Pattern Sets.
type RegexPatternSetStore struct {
	*ResourceStore[RegexPatternSet]
}

// NewRegexPatternSetStore creates a new Regex Pattern Set store.
func NewRegexPatternSetStore(store storage.BasicStorage, accountId, region string) *RegexPatternSetStore {
	return &RegexPatternSetStore{
		ResourceStore: NewResourceStore[RegexPatternSet](store, regexPatternSetBucketName, NewARNBuilder(accountId, region), regexPatternSetAccessor),
	}
}

// Create creates a new Regex Pattern Set.
func (s *RegexPatternSetStore) Create(id, name, description string, regularPatterns []string) (*RegexPatternSet, error) {
	regexPatternSet := &RegexPatternSet{
		ID:              id,
		Name:            name,
		Description:     description,
		RegularPatterns: regularPatterns,
		Tags:            []types.Tag{},
		CreatedAt:       time.Now(),
		ModifiedAt:      time.Now(),
	}
	regexPatternSet.ARN = s.arnBuilder.BuildRegexPatternSetARN(id)
	SetTimestamps(&regexPatternSetAccessor, regexPatternSet)
	if err := s.Put(id, regexPatternSet, "create_regex_pattern_set"); err != nil {
		return nil, err
	}
	return regexPatternSet, nil
}

// Update updates an existing Regex Pattern Set.
func (s *RegexPatternSetStore) Update(id, lockToken string, regularPatterns []string) (*RegexPatternSet, error) {
	return s.UpdateWithLockToken(id, lockToken, func(regexPatternSet *RegexPatternSet) error {
		regexPatternSet.RegularPatterns = regularPatterns
		return nil
	}, "update_regex_pattern_set")
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
