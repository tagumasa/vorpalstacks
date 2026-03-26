// Package common provides common AWS store utilities for vorpalstacks.
package common

import (
	"vorpalstacks/internal/core/storage"
)

// TagStore provides common tag CRUD operations
type TagStore struct {
	*BaseStore
}

// NewTagStore creates a new tag store for the given service
func NewTagStore(store storage.BasicStorage, serviceName string) *TagStore {
	return &TagStore{
		BaseStore: NewBaseStore(store.Bucket(serviceName+"-tags"), serviceName+"-tags"),
	}
}

// NewTagStoreWithRegion creates a new tag store with region isolation
func NewTagStoreWithRegion(store storage.BasicStorage, serviceName, region string) *TagStore {
	bucketName := serviceName + "-tags-" + region
	return &TagStore{
		BaseStore: NewBaseStore(store.Bucket(bucketName), bucketName),
	}
}

// ListTags returns tags as map[string]string
func (t *TagStore) ListTags(resourceKey string) (map[string]string, error) {
	var tags map[string]string
	if err := t.Get(resourceKey, &tags); err != nil {
		return make(map[string]string), nil
	}
	if tags == nil {
		tags = make(map[string]string)
	}
	return tags, nil
}

// ListTagsAsSlice returns tags as []Tag
func (t *TagStore) ListTagsAsSlice(resourceKey string) ([]Tag, error) {
	tags, err := t.ListTags(resourceKey)
	if err != nil {
		return nil, err
	}
	var result []Tag
	for k, v := range tags {
		result = append(result, Tag{Key: k, Value: v})
	}
	return result, nil
}

// TagResource adds or updates tags for a resource
func (t *TagStore) TagResource(resourceKey string, newTags map[string]string) error {
	existingTags, err := t.ListTags(resourceKey)
	if err != nil {
		return err
	}
	if existingTags == nil {
		existingTags = make(map[string]string)
	}
	for k, v := range newTags {
		existingTags[k] = v
	}
	return t.Put(resourceKey, existingTags)
}

// TagResourceFromSlice adds tags from []Tag format
func (t *TagStore) TagResourceFromSlice(resourceKey string, tags []Tag) error {
	tagMap := make(map[string]string)
	for _, tag := range tags {
		tagMap[tag.Key] = tag.Value
	}
	return t.TagResource(resourceKey, tagMap)
}

// UntagResource removes specified tag keys from a resource
func (t *TagStore) UntagResource(resourceKey string, tagKeys []string) error {
	tags, err := t.ListTags(resourceKey)
	if err != nil {
		return err
	}
	for _, key := range tagKeys {
		delete(tags, key)
	}
	if len(tags) == 0 {
		return t.Delete(resourceKey)
	}
	return t.Put(resourceKey, tags)
}
