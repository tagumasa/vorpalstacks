package kms

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

// validateKMSTags checks tag count, individual key/value lengths and the
// aws: key reservation against the AWS-wide tag limits.
func validateKMSTags(tags []tagutil.Tag) error {
	if v, _ := tagutil.CheckTags(tags, tagutil.StandardLimits()); v != tagutil.OK {
		return ErrTagException
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

	tagList := tagutil.ParseTagsWithKeyNames(req.Parameters, "Tags", "TagKey", "TagValue")
	if err := s.tagResourceCore(stores, s.resolveCallerPrincipal(reqCtx, req), s.getKeyID(req.Parameters), tagList); err != nil {
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

	tagKeys := tagutil.ParseTagKeysAsSlice(req.Parameters, "TagKeys")
	if err := s.untagResourceCore(stores, s.resolveCallerPrincipal(reqCtx, req), s.getKeyID(req.Parameters), tagKeys); err != nil {
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

	result, err := s.listResourceTagsCore(stores, s.resolveCallerPrincipal(reqCtx, req), s.getKeyID(req.Parameters),
		pagination.GetMarker(req.Parameters),
		pagination.GetMaxItems(req.Parameters, 100))
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"Tags":      tagutil.ToResponseWithKeyNames(result.Items, "TagKey", "TagValue"),
		"Truncated": result.IsTruncated,
	}
	if result.NextMarker != "" {
		response["NextMarker"] = result.NextMarker
	}

	return response, nil
}
