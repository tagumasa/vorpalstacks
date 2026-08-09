package route53

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const cidrCollectionBucketName = "route53_cidr_collections"

// CidrCollectionStore manages Route 53 CIDR collections.
type CidrCollectionStore struct {
	*common.BaseStore
}

// NewCidrCollectionStore creates a new CidrCollectionStore.
func NewCidrCollectionStore(store storage.BasicStorage) *CidrCollectionStore {
	return &CidrCollectionStore{
		BaseStore: common.NewBaseStore(store.Bucket(cidrCollectionBucketName), "route53"),
	}
}

// Create creates a new CIDR collection.
func (s *CidrCollectionStore) Create(c *CidrCollection) error {
	if s.Exists(c.ID) {
		return NewStoreError("create_cidr_collection", common.ErrAlreadyExists)
	}
	c.CreatedAt = time.Now()
	c.Version = 1
	if c.Locations == nil {
		c.Locations = make(map[string][]string)
	}
	if err := s.BaseStore.Put(c.ID, c); err != nil {
		return NewStoreError("create_cidr_collection", err)
	}
	return nil
}

// Get retrieves a CIDR collection by its ID.
func (s *CidrCollectionStore) Get(id string) (*CidrCollection, error) {
	var c CidrCollection
	if err := s.BaseStore.Get(id, &c); err != nil {
		return nil, NewStoreError("get_cidr_collection", err)
	}
	return &c, nil
}

// Delete removes a CIDR collection.
func (s *CidrCollectionStore) Delete(id string) error {
	if !s.Exists(id) {
		return NewStoreError("delete_cidr_collection", common.ErrNotFound)
	}
	return s.BaseStore.Delete(id)
}

// List returns CIDR collections with pagination.
func (s *CidrCollectionStore) List(nextToken string, maxResults int) (*CidrCollectionListResult, error) {
	result, err := common.List[CidrCollection](s.BaseStore, common.ListOptions{
		Marker:   nextToken,
		MaxItems: maxResults,
	}, nil)
	if err != nil {
		return nil, NewStoreError("list_cidr_collections", err)
	}
	return &CidrCollectionListResult{
		Collections: result.Items,
		IsTruncated: result.IsTruncated,
		NextToken:   result.NextMarker,
	}, nil
}

// Update updates a CIDR collection.
func (s *CidrCollectionStore) Update(c *CidrCollection) error {
	existing, err := s.Get(c.ID)
	if err != nil {
		return NewStoreError("update_cidr_collection", err)
	}
	c.Version = existing.Version + 1
	c.CreatedAt = existing.CreatedAt
	if c.Locations == nil {
		c.Locations = existing.Locations
	}
	return s.BaseStore.Put(c.ID, c)
}

// Exists checks whether a CIDR collection exists.
func (s *CidrCollectionStore) Exists(id string) bool {
	return s.BaseStore.Exists(id)
}

// GetByName looks up a CIDR collection by name.
func (s *CidrCollectionStore) GetByName(name string) (*CidrCollection, error) {
	return common.FindFirst[CidrCollection](s.BaseStore, func(c *CidrCollection) bool {
		return c.Name == name
	})
}

// FindLocationForIP returns the LocationName whose CIDR blocks contain
// the given IP address, or empty string if no match. Longest prefix match
// is used to select the most specific location.
func (s *CidrCollectionStore) FindLocationForIP(collectionID, ip string) (string, error) {
	c, err := s.Get(collectionID)
	if err != nil {
		return "", err
	}
	return findCidrLocation(c.Locations, ip), nil
}

// findCidrLocation performs longest-prefix-match CIDR lookup.
func findCidrLocation(locations map[string][]string, ip string) string {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ""
	}
	bestMatch := ""
	bestPrefixLen := -1
	for locName, cidrs := range locations {
		for _, cidr := range cidrs {
			_, ipNet, err := net.ParseCIDR(cidr)
			if err != nil {
				continue
			}
			if ipNet.Contains(parsedIP) {
				prefixLen, _ := ipNet.Mask.Size()
				if prefixLen > bestPrefixLen {
					bestPrefixLen = prefixLen
					bestMatch = locName
				}
			}
		}
	}
	return bestMatch
}

// unmarshalCidrCollection deserialises a CIDR collection from JSON bytes.
func unmarshalCidrCollection(data []byte) (*CidrCollection, error) {
	var c CidrCollection
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("unmarshal cidr collection: %w", err)
	}
	return &c, nil
}
