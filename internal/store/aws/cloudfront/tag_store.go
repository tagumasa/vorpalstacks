package cloudfront

import (
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/types"
)

// TagStore provides tag storage operations for CloudFront resources.
type TagStore struct {
	*common.TagStore
}

// NewTagStore creates a new TagStore instance for CloudFront resources.
func NewTagStore(store storage.BasicStorage) *TagStore {
	return &TagStore{
		TagStore: common.NewTagStore(store, "cloudfront"),
	}
}

// ListTagsForResource returns the tags for a CloudFront resource.
func (s *TagStore) ListTagsForResource(resourceKey string) ([]types.Tag, error) {
	return s.TagStore.ListAsSlice(resourceKey)
}

// TagResource adds tags to a CloudFront resource.
func (s *TagStore) Tag(resourceKey string, tags []types.Tag) error {
	return s.TagStore.TagFromSlice(resourceKey, tags)
}

// Raw returns the underlying tag store.
func (s *TagStore) Raw() *TagStore {
	return s
}
