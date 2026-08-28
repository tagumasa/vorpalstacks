package neptunegraph

// Tag Core functions: the single validation and persistence path for the
// resource tag operations.

import (
	ngstore "vorpalstacks/internal/store/aws/rds/neptunegraph"
)

// ListTagsForResourceInput carries the wire-parsed ListTagsForResource request.
type ListTagsForResourceInput struct {
	ResourceArn string
}

// TagResourceInput carries the wire-parsed TagResource request.
type TagResourceInput struct {
	ResourceArn string
	Tags        map[string]string
}

// UntagResourceInput carries the wire-parsed UntagResource request.
type UntagResourceInput struct {
	ResourceArn string
	TagKeys     []string
}

// listTagsForResourceCore returns all tags associated with a resource ARN.
func (s *NeptuneGraphService) listTagsForResourceCore(store *ngstore.NeptuneGraphStore, in *ListTagsForResourceInput) (map[string]string, error) {
	resourceArn := in.ResourceArn
	if resourceArn == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "resourceArn")
	}

	return store.GetTags(resourceArn)
}

// tagResourceCore validates the tag constraints and adds or updates tags on
// the resource.
func (s *NeptuneGraphService) tagResourceCore(store *ngstore.NeptuneGraphStore, in *TagResourceInput) (map[string]string, error) {
	resourceArn := in.ResourceArn
	if resourceArn == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "resourceArn")
	}

	tags := in.Tags
	if err := validateTags(tags); err != nil {
		return nil, err
	}

	if err := store.AddTags(resourceArn, tags); err != nil {
		return nil, err
	}

	return tags, nil
}

// untagResourceCore removes the given tag keys from a resource's tag set and
// returns the remaining tags.
func (s *NeptuneGraphService) untagResourceCore(store *ngstore.NeptuneGraphStore, in *UntagResourceInput) (map[string]string, error) {
	resourceArn := in.ResourceArn
	if resourceArn == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "resourceArn")
	}

	if err := store.RemoveTags(resourceArn, in.TagKeys); err != nil {
		return nil, err
	}

	tags, _ := store.GetTags(resourceArn)
	return tags, nil
}
