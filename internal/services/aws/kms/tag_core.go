package kms

// tag_core.go carries the Core functions of the KMS resource-tagging
// family. Cores resolve and authorise the key, apply the AWS tag limits
// against the accumulated tag count, and mutate the tag store; handlers
// parse the wire members and serialise the paginated results.

import (
	"vorpalstacks/internal/common/pagination"
	tagutil "vorpalstacks/internal/common/tags"

	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// TagListResult is the paginated Core result of ListResourceTags.
type TagListResult struct {
	Items       []tagutil.Tag
	IsTruncated bool
	NextMarker  string
}

// tagResourceCore is the single entry point for tagging a key. The
// accumulated-count check (existing non-overwritten tags plus new tags
// within the AWS per-resource limit) runs against the live tag store
// before the write.
func (s *KMSService) tagResourceCore(stores *kmsStores, principal, keyID string, tags []tagutil.Tag) error {
	key, err := s.resolveKeyByKeyID(stores, keyID)
	if err != nil {
		return err
	}

	if err := s.authorizeOperation(stores, principal, "TagResource", key.KeyID, nil); err != nil {
		return err
	}
	// AWS: tagging is rejected when the key is PendingDeletion.
	if key.KeyState == kmsstore.KeyStatePendingDeletion {
		return ErrKeyPendingDeletion
	}
	if len(tags) == 0 {
		return nil
	}

	if err := validateKMSTags(tags); err != nil {
		return err
	}

	// Check accumulated tag count: existing tags minus duplicates in new tags
	// plus new tags must not exceed the limit.
	existingTags, err := stores.keys.TagStore.ListAsSlice(key.KeyID)
	if err != nil {
		return err
	}
	newKeySet := make(map[string]bool, len(tags))
	for _, t := range tags {
		newKeySet[t.Key] = true
	}
	nonOverwritten := 0
	for _, et := range existingTags {
		if !newKeySet[et.Key] {
			nonOverwritten++
		}
	}
	if nonOverwritten+len(tags) > tagutil.MaxTagsPerResource {
		return ErrTagException
	}

	return stores.keys.TagStore.Tag(key.KeyID, tagutil.ToMap(tags))
}

// untagResourceCore is the single entry point for removing tags from a key.
func (s *KMSService) untagResourceCore(stores *kmsStores, principal, keyID string, tagKeys []string) error {
	key, err := s.resolveKeyByKeyID(stores, keyID)
	if err != nil {
		return err
	}

	if err := s.authorizeOperation(stores, principal, "UntagResource", key.KeyID, nil); err != nil {
		return err
	}
	if key.KeyState == kmsstore.KeyStatePendingDeletion {
		return ErrKeyPendingDeletion
	}
	if len(tagKeys) == 0 {
		return nil
	}

	return stores.keys.TagStore.Untag(key.KeyID, tagKeys)
}

// listResourceTagsCore is the single entry point for listing a key's tags.
// The marker validation deliberately runs after the tag fetch, matching
// the original failure precedence where a stale marker on a missing key
// surfaces the key error first.
func (s *KMSService) listResourceTagsCore(stores *kmsStores, principal, keyID, marker string, maxItems int) (*TagListResult, error) {
	key, err := s.resolveKeyByKeyID(stores, keyID)
	if err != nil {
		return nil, err
	}

	if err := s.authorizeOperation(stores, principal, "ListResourceTags", key.KeyID, nil); err != nil {
		return nil, err
	}
	tags, err := stores.keys.TagStore.ListAsSlice(key.KeyID)
	if err != nil {
		return nil, err
	}

	if err := validateMarkerLength(marker); err != nil {
		return nil, err
	}

	result := pagination.PaginateSlice(tags, marker, maxItems, func(t tagutil.Tag) string {
		return t.Key
	})
	return &TagListResult{
		Items:       result.Items,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}
