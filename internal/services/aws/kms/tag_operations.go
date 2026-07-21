package kms

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/utils/aws/types"

	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// KMS tag limits per AWS spec.
const (
	maxKMSTags     = 50
	maxTagKeyLen   = 128
	maxTagValueLen = 256
)

// validateKMSTags checks tag count and individual key/value lengths.
func validateKMSTags(tags []types.Tag) error {
	if len(tags) > maxKMSTags {
		return ErrTagException
	}
	for _, t := range tags {
		if len(t.Key) == 0 || len(t.Key) > maxTagKeyLen {
			return ErrTagException
		}
		if len(t.Value) > maxTagValueLen {
			return ErrTagException
		}
		// AWS KMS tag keys cannot start with "aws:"
		if strings.HasPrefix(strings.ToLower(t.Key), "aws:") {
			return ErrTagException
		}
	}
	return nil
}

// TagResource adds one or more tags to a KMS key.
// Tags are used to identify and organise your keys.
func (s *KMSService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveKey(stores, req.Parameters)
	if err != nil {
		return nil, err
	}

	if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "TagResource", key.KeyID, nil); err != nil {
		return nil, err
	}
	// AWS: tagging is rejected when the key is PendingDeletion.
	if key.KeyState == kmsstore.KeyStatePendingDeletion {
		return nil, ErrKeyPendingDeletion
	}
	tagList := tagutil.ParseTagsWithKeyNames(req.Parameters, "Tags", "TagKey", "TagValue")
	if len(tagList) == 0 {
		return response.EmptyResponse(), nil
	}

	if err := validateKMSTags(tagList); err != nil {
		return nil, err
	}

	// Check accumulated tag count: existing tags minus duplicates in new tags
	// plus new tags must not exceed the limit.
	existingTags, err := stores.keys.TagStore.ListAsSlice(key.KeyID)
	if err != nil {
		return nil, err
	}
	newKeySet := make(map[string]bool, len(tagList))
	for _, t := range tagList {
		newKeySet[t.Key] = true
	}
	nonOverwritten := 0
	for _, et := range existingTags {
		if !newKeySet[et.Key] {
			nonOverwritten++
		}
	}
	if nonOverwritten+len(tagList) > maxKMSTags {
		return nil, ErrTagException
	}

	if err := stores.keys.TagStore.Tag(key.KeyID, tagutil.ToMap(tagList)); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes one or more tags from a KMS key.
func (s *KMSService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveKey(stores, req.Parameters)
	if err != nil {
		return nil, err
	}

	if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "UntagResource", key.KeyID, nil); err != nil {
		return nil, err
	}
	if key.KeyState == kmsstore.KeyStatePendingDeletion {
		return nil, ErrKeyPendingDeletion
	}
	tagKeys := tagutil.ParseTagKeysAsSlice(req.Parameters, "TagKeys")
	if len(tagKeys) == 0 {
		return response.EmptyResponse(), nil
	}

	if err := stores.keys.TagStore.Untag(key.KeyID, tagKeys); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListResourceTags returns all tags associated with a KMS key.
// Results can be paginated using the Marker and MaxItems parameters.
func (s *KMSService) ListResourceTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveKey(stores, req.Parameters)
	if err != nil {
		return nil, err
	}

	if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "ListResourceTags", key.KeyID, nil); err != nil {
		return nil, err
	}
	tags, err := stores.keys.TagStore.ListAsSlice(key.KeyID)
	if err != nil {
		return nil, err
	}

	marker := pagination.GetMarker(req.Parameters)
	maxItems := pagination.GetMaxItems(req.Parameters, 100)

	result := pagination.PaginateSlice(tags, marker, maxItems, func(t types.Tag) string {
		return t.Key
	})

	response := map[string]interface{}{
		"Tags":      tagutil.ToResponseWithKeyNames(result.Items, "TagKey", "TagValue"),
		"Truncated": result.IsTruncated,
	}
	if result.NextMarker != "" {
		response["NextMarker"] = result.NextMarker
	}

	return response, nil
}
