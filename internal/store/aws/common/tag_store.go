package common

import (
	"bytes"
	"strings"
	"sync"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/utils/aws/types"
)

const indexSeparator = "\x00"

// TagStore manages resource tags using a Loki-style inverted index.
// The main store holds tag maps keyed by resource key, while the index store
// enables efficient lookup of resources by tag key or tag key-value pair.
type TagStore struct {
	main  *BaseStore
	index *BaseStore
	mu    sync.Mutex
}

// NewTagStore creates a TagStore with region-agnostic bucket names.
func NewTagStore(store storage.BasicStorage, serviceName string) *TagStore {
	mainName := serviceName + "-tags"
	indexName := serviceName + "-tag-idx"
	return &TagStore{
		main:  NewBaseStore(store.Bucket(mainName), mainName),
		index: NewBaseStore(store.Bucket(indexName), indexName),
	}
}

// NewTagStoreWithRegion creates a TagStore with region-scoped bucket names.
func NewTagStoreWithRegion(store storage.BasicStorage, serviceName, region string) *TagStore {
	mainName := serviceName + "-tags-" + region
	indexName := serviceName + "-tag-idx-" + region
	return &TagStore{
		main:  NewBaseStore(store.Bucket(mainName), mainName),
		index: NewBaseStore(store.Bucket(indexName), indexName),
	}
}

// List returns all tags for the given resource as a key-value map.
// Returns an empty map if the resource has no tags.
func (t *TagStore) List(resourceKey string) (map[string]string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	var tags map[string]string
	if err := t.main.Get(resourceKey, &tags); err != nil {
		return make(map[string]string), nil
	}
	if tags == nil {
		tags = make(map[string]string)
	}
	return tags, nil
}

// ListAsSlice returns all tags for the given resource as a slice of Tag structs.
func (t *TagStore) ListAsSlice(resourceKey string) ([]types.Tag, error) {
	tags, err := t.List(resourceKey)
	if err != nil {
		return nil, err
	}
	result := make([]types.Tag, 0, len(tags))
	for k, v := range tags {
		result = append(result, types.Tag{Key: k, Value: v})
	}
	return result, nil
}

// Tag merges the given tags into the existing tags for a resource.
// Existing tag keys are overwritten; new tag keys are added.
// For overwritten keys, stale index entries are removed before inserting new ones.
func (t *TagStore) Tag(resourceKey string, newTags map[string]string) error {
	if len(newTags) == 0 {
		return nil
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	existing, err := t.listUnlocked(resourceKey)
	if err != nil {
		return err
	}
	for k, v := range newTags {
		existing[k] = v
	}
	if err := t.main.Put(resourceKey, existing); err != nil {
		return err
	}

	suffix := indexSeparator + resourceKey
	for k, v := range newTags {
		oldIdxPrefix := k + "="
		if err := t.index.ScanPrefix(oldIdxPrefix, func(idxKey string, _ []byte) error {
			if strings.HasSuffix(idxKey, suffix) {
				return t.index.Delete(idxKey)
			}
			return nil
		}); err != nil {
			return err
		}
		idxKey := k + "=" + v + indexSeparator + resourceKey
		if err := t.index.Put(idxKey, []byte{0x01}); err != nil {
			return err
		}
	}
	return nil
}

// TagFromSlice converts a slice of Tag structs to a map and delegates to Tag.
func (t *TagStore) TagFromSlice(resourceKey string, tags []types.Tag) error {
	tagMap := make(map[string]string, len(tags))
	for _, tag := range tags {
		tagMap[tag.Key] = tag.Value
	}
	return t.Tag(resourceKey, tagMap)
}

// Untag removes the specified tag keys from a resource.
// If all tags are removed, the resource entry is deleted from the main store.
// The corresponding inverted index entries are also removed.
func (t *TagStore) Untag(resourceKey string, tagKeys []string) error {
	if len(tagKeys) == 0 {
		return nil
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	tags, err := t.listUnlocked(resourceKey)
	if err != nil {
		return err
	}

	removeSet := make(map[string]bool, len(tagKeys))
	for _, k := range tagKeys {
		delete(tags, k)
		removeSet[k] = true
	}

	if len(tags) == 0 {
		if err := t.main.Delete(resourceKey); err != nil {
			return err
		}
	} else {
		if err := t.main.Put(resourceKey, tags); err != nil {
			return err
		}
	}

	suffix := indexSeparator + resourceKey
	for k := range removeSet {
		prefix := k + "="
		if err := t.index.ScanPrefix(prefix, func(idxKey string, _ []byte) error {
			if strings.HasSuffix(idxKey, suffix) {
				return t.index.Delete(idxKey)
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

// Delete removes all tags and index entries for a resource.
func (t *TagStore) Delete(resourceKey string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := t.main.Delete(resourceKey); err != nil {
		return err
	}

	suffix := indexSeparator + resourceKey
	return t.index.ScanPrefix("", func(idxKey string, _ []byte) error {
		if strings.HasSuffix(idxKey, suffix) {
			return t.index.Delete(idxKey)
		}
		return nil
	})
}

// FindByTag returns all resource keys that have a tag with the given key (any value).
func (t *TagStore) FindByTag(tagKey string) ([]string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	prefix := tagKey + "="
	result := make([]string, 0)
	err := t.index.ScanPrefix(prefix, func(idxKey string, _ []byte) error {
		idx := bytes.Index([]byte(idxKey), []byte(indexSeparator))
		if idx >= 0 {
			result = append(result, idxKey[idx+1:])
		}
		return nil
	})
	return result, err
}

// FindByTagValue returns all resource keys that have a tag matching both the given key and value.
func (t *TagStore) FindByTagValue(tagKey, tagValue string) ([]string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	prefix := tagKey + "=" + tagValue + indexSeparator
	result := make([]string, 0)
	err := t.index.ScanPrefix(prefix, func(idxKey string, _ []byte) error {
		result = append(result, idxKey[len(prefix):])
		return nil
	})
	return result, err
}

// RebuildIndex deletes and recreates the inverted index from the main tag store.
// Use this after bulk tag migrations or data corruption recovery.
func (t *TagStore) RebuildIndex() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := t.index.DeleteByPrefix(""); err != nil {
		return err
	}

	return t.main.ForEach(func(resourceKey string, value []byte) error {
		var tags map[string]string
		if err := t.main.Get(resourceKey, &tags); err != nil {
			return err
		}
		for k, v := range tags {
			idxKey := k + "=" + v + indexSeparator + resourceKey
			if err := t.index.Put(idxKey, []byte{0x01}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (t *TagStore) listUnlocked(resourceKey string) (map[string]string, error) {
	var tags map[string]string
	if err := t.main.Get(resourceKey, &tags); err != nil {
		return make(map[string]string), nil
	}
	if tags == nil {
		tags = make(map[string]string)
	}
	return tags, nil
}
