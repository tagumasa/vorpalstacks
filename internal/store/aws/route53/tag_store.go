package route53

import (
	"net/http"

	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// MaxTagsPerResource is Route 53's own per-resource tag total: the
// ChangeTagsForResource API reference caps a hosted zone or health check
// at ten tags ("You can add a maximum of 10 tags to a health check or a
// hosted zone"), diverging from the AWS-wide fifty most services
// document.
const MaxTagsPerResource = 10

// TagStore manages Route 53 resource tags.
type TagStore struct {
	*common.TagStore
}

// NewTagStore creates a new TagStore.
func NewTagStore(store storage.BasicStorage) *TagStore {
	return &TagStore{
		TagStore: common.NewTagStore(store, "route53", common.TagBudget{
			MaxKeys:  MaxTagsPerResource,
			Exceeded: common.NewAWSError("InvalidInput", "Too many tags.", http.StatusBadRequest),
		}),
	}
}

// ListTagsForResource returns tags for a resource.
func (s *TagStore) ListTagsForResource(resourceKey string) ([]types.Tag, error) {
	return s.TagStore.ListAsSlice(resourceKey)
}

// TagResource adds or updates tags for a Route 53 resource.
func (s *TagStore) Tag(resourceKey string, tags []types.Tag) error {
	return s.TagStore.TagFromSlice(resourceKey, tags)
}
