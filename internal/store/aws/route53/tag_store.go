package route53

import (
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// TagStore manages Route 53 resource tags.
type TagStore struct {
	*common.TagStore
}

// NewTagStore creates a new TagStore.
func NewTagStore(store storage.BasicStorage) *TagStore {
	return &TagStore{
		TagStore: common.NewTagStore(store, "route53"),
	}
}

// ListTagsForResource returns tags for a resource.
func (s *TagStore) ListTagsForResource(resourceKey string) ([]common.Tag, error) {
	return s.TagStore.ListTagsAsSlice(resourceKey)
}

// TagResource adds or updates tags for a Route 53 resource.
func (s *TagStore) TagResource(resourceKey string, tags []common.Tag) error {
	return s.TagStore.TagResourceFromSlice(resourceKey, tags)
}
