package cloudfront

import (
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const invalidationBucketName = "cloudfront_invalidations"

// Invalidation represents a CloudFront cache invalidation.
type Invalidation struct {
	ID              string    `json:"id"`
	DistributionID  string    `json:"distributionId"`
	Status          string    `json:"status"`
	CreateTime      time.Time `json:"createTime"`
	CallerReference string    `json:"callerReference"`
	Paths           []string  `json:"paths"`
}

// InvalidationStore provides persisted storage for CloudFront invalidations.
type InvalidationStore struct {
	*common.BaseStore
	mu sync.Mutex
}

// NewInvalidationStore creates a new InvalidationStore instance.
func NewInvalidationStore(store storage.BasicStorage) *InvalidationStore {
	return &InvalidationStore{
		BaseStore: common.NewBaseStore(store.Bucket(invalidationBucketName), "cloudfront"),
	}
}

func invalidationKey(distID, invID string) string {
	return fmt.Sprintf("%s/%s", distID, invID)
}

// Create persists a new invalidation and returns the created entity with a
// generated ID.
func (s *InvalidationStore) Create(distID, callerRef string, paths []string) (*Invalidation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, err := generateDistributionID()
	if err != nil {
		return nil, NewStoreError("create_invalidation", err)
	}

	inv := &Invalidation{
		ID:              id,
		DistributionID:  distID,
		Status:          "InProgress",
		CreateTime:      time.Now().UTC(),
		CallerReference: callerRef,
		Paths:           paths,
	}

	key := invalidationKey(distID, id)
	if err := s.BaseStore.Put(key, inv); err != nil {
		return nil, NewStoreError("create_invalidation", err)
	}

	return inv, nil
}

// Update persists an updated invalidation (e.g. status transition).
func (s *InvalidationStore) Update(inv *Invalidation) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := invalidationKey(inv.DistributionID, inv.ID)
	return s.BaseStore.Put(key, inv)
}

// Get retrieves a specific invalidation by distribution ID and invalidation ID.
func (s *InvalidationStore) Get(distID, invID string) (*Invalidation, error) {
	var inv Invalidation
	if err := s.BaseStore.Get(invalidationKey(distID, invID), &inv); err != nil {
		return nil, NewStoreError("get_invalidation", err)
	}
	return &inv, nil
}

// InvalidationListResult represents the result of listing invalidations.
type InvalidationListResult struct {
	Invalidations []*Invalidation
	IsTruncated   bool
	NextMarker    string
}

// List returns a paginated list of invalidations for a given distribution.
func (s *InvalidationStore) List(distID, marker string, maxItems int) (*InvalidationListResult, error) {
	prefix := distID + "/"
	listMarker := ""
	if marker != "" {
		listMarker = invalidationKey(distID, marker)
	}
	result, err := common.List[Invalidation](s.BaseStore, common.ListOptions{
		Prefix:   prefix,
		Marker:   listMarker,
		MaxItems: maxItems,
	}, nil)
	if err != nil {
		return nil, NewStoreError("list_invalidations", err)
	}

	return &InvalidationListResult{
		Invalidations: result.Items,
		IsTruncated:   result.IsTruncated,
		NextMarker:    result.NextMarker,
	}, nil
}

// DeleteByDistribution removes all invalidations for a given distribution ID.
func (s *InvalidationStore) DeleteByDistribution(distID string) error {
	prefix := distID + "/"
	return s.BaseStore.DeleteByPrefix(prefix)
}
