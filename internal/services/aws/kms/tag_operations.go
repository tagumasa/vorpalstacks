package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/utils/aws/types"
)

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
	tagList := tagutil.ParseTagsWithKeyNames(req.Parameters, "Tags", "TagKey", "TagValue")
	if len(tagList) == 0 {
		return response.EmptyResponse(), nil
	}

	if err := stores.keys.AddTags(key.KeyID, tagList); err != nil {
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
	tagKeys := tagutil.ParseTagKeysAsSlice(req.Parameters, "TagKeys")
	if len(tagKeys) == 0 {
		return response.EmptyResponse(), nil
	}

	if err := stores.keys.RemoveTags(key.KeyID, tagKeys); err != nil {
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
	tags, err := stores.keys.ListTags(key.KeyID)
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
