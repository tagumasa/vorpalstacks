package cloudfront

import (
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const cloudfrontTagBucketName = "cloudfront_tags"

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
func (s *TagStore) ListTagsForResource(resourceKey string) ([]common.Tag, error) {
	return s.TagStore.ListTagsAsSlice(resourceKey)
}

// TagResource adds tags to a CloudFront resource.
func (s *TagStore) TagResource(resourceKey string, tags []common.Tag) error {
	return s.TagStore.TagResourceFromSlice(resourceKey, tags)
}

// Raw returns the underlying tag store.
func (s *TagStore) Raw() *TagStore {
	return s
}
