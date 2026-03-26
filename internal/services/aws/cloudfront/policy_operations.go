package cloudfront

import (
	"context"
	"time"

	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// CreateCachePolicy creates a cache policy.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_CreateCachePolicy.html
func (s *CloudFrontService) CreateCachePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, NewAPIError("InvalidArgument", "Name is required", 400)
	}

	configMap := request.GetMapParam(req.Parameters, "CachePolicyConfig")
	if configMap == nil {
		return nil, NewAPIError("InvalidArgument", "CachePolicyConfig is required", 400)
	}

	config := &cloudfrontstore.CachePolicyConfig{
		Name:       name,
		Comment:    request.GetStringParam(configMap, "Comment"),
		DefaultTTL: int64(request.GetIntParam(configMap, "DefaultTTL")),
		MaxTTL:     int64(request.GetIntParam(configMap, "MaxTTL")),
		MinTTL:     int64(request.GetIntParam(configMap, "MinTTL")),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cachePolicy, err := store.cachePolicies.Create(name, "", config)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"CachePolicy": map[string]interface{}{
			"Id":                cachePolicy.ID,
			"ARN":               cachePolicy.ARN,
			"Name":              cachePolicy.Name,
			"CachePolicyConfig": cachePolicy.CachePolicyConfig,
			"LastModifiedTime":  cachePolicy.ModifiedAt.Format(time.RFC3339),
		},
	}, nil
}

// GetCachePolicy returns a cache policy.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_GetCachePolicy.html
func (s *CloudFrontService) GetCachePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cachePolicy, err := store.cachePolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, NewAPIError("NoSuchCachePolicy", "Cache policy not found", 404)
		}
		return nil, err
	}

	return map[string]interface{}{
		"CachePolicy": map[string]interface{}{
			"Id":                cachePolicy.ID,
			"ARN":               cachePolicy.ARN,
			"Name":              cachePolicy.Name,
			"CachePolicyConfig": cachePolicy.CachePolicyConfig,
			"LastModifiedTime":  cachePolicy.ModifiedAt.Format(time.RFC3339),
		},
	}, nil
}

// ListCachePolicies lists cache policies.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_ListCachePolicies.html
func (s *CloudFrontService) ListCachePolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	if maxItems == 0 {
		maxItems = 100
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.cachePolicies.List(marker, maxItems)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.CachePolicies))
	for _, cp := range result.CachePolicies {
		items = append(items, map[string]interface{}{
			"Id":               cp.ID,
			"ARN":              cp.ARN,
			"Name":             cp.Name,
			"LastModifiedTime": cp.ModifiedAt.Format(time.RFC3339),
		})
	}

	return map[string]interface{}{
		"CachePolicyList": map[string]interface{}{
			"Items":       items,
			"IsTruncated": result.IsTruncated,
			"NextMarker":  result.NextMarker,
		},
	}, nil
}

// DeleteCachePolicy deletes a cache policy.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_DeleteCachePolicy.html
func (s *CloudFrontService) DeleteCachePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	err = store.cachePolicies.Delete(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, NewAPIError("NoSuchCachePolicy", "Cache policy not found", 404)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// CreateOriginRequestPolicy creates an origin request policy.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_CreateOriginRequestPolicy.html
func (s *CloudFrontService) CreateOriginRequestPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, NewAPIError("InvalidArgument", "Name is required", 400)
	}

	configMap := request.GetMapParam(req.Parameters, "OriginRequestPolicyConfig")
	if configMap == nil {
		return nil, NewAPIError("InvalidArgument", "OriginRequestPolicyConfig is required", 400)
	}

	config := &cloudfrontstore.OriginRequestPolicyConfig{
		Name:    name,
		Comment: request.GetStringParam(configMap, "Comment"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := store.originRequestPolicies.Create(name, "", config)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"OriginRequestPolicy": map[string]interface{}{
			"Id":                        policy.ID,
			"ARN":                       policy.ARN,
			"Name":                      policy.Name,
			"OriginRequestPolicyConfig": policy.OriginRequestPolicyConfig,
			"LastModifiedTime":          policy.ModifiedAt.Format(time.RFC3339),
		},
	}, nil
}

// GetOriginRequestPolicy returns an origin request policy.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_GetOriginRequestPolicy.html
func (s *CloudFrontService) GetOriginRequestPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := store.originRequestPolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, NewAPIError("NoSuchOriginRequestPolicy", "Origin request policy not found", 404)
		}
		return nil, err
	}

	return map[string]interface{}{
		"OriginRequestPolicy": map[string]interface{}{
			"Id":                        policy.ID,
			"ARN":                       policy.ARN,
			"Name":                      policy.Name,
			"OriginRequestPolicyConfig": policy.OriginRequestPolicyConfig,
			"LastModifiedTime":          policy.ModifiedAt.Format(time.RFC3339),
		},
	}, nil
}

// ListOriginRequestPolicies lists origin request policies.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_ListOriginRequestPolicies.html
func (s *CloudFrontService) ListOriginRequestPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	if maxItems == 0 {
		maxItems = 100
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.originRequestPolicies.List(marker, maxItems)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.OriginRequestPolicies))
	for _, p := range result.OriginRequestPolicies {
		items = append(items, map[string]interface{}{
			"Id":               p.ID,
			"ARN":              p.ARN,
			"Name":             p.Name,
			"LastModifiedTime": p.ModifiedAt.Format(time.RFC3339),
		})
	}

	return map[string]interface{}{
		"OriginRequestPolicyList": map[string]interface{}{
			"Items":       items,
			"IsTruncated": result.IsTruncated,
			"NextMarker":  result.NextMarker,
		},
	}, nil
}

// DeleteOriginRequestPolicy deletes an origin request policy.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_DeleteOriginRequestPolicy.html
func (s *CloudFrontService) DeleteOriginRequestPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	err = store.originRequestPolicies.Delete(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, NewAPIError("NoSuchOriginRequestPolicy", "Origin request policy not found", 404)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists tags for a CloudFront resource.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_ListTagsForResource.html
func (s *CloudFrontService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetStringParam(req.Parameters, "Resource")
	if arn == "" {
		arn = request.GetStringParam(req.Parameters, "ResourceARN")
	}
	if arn == "" {
		return nil, NewAPIError("InvalidArgument", "Resource is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tags, err := store.tags.ListTagsForResource(arn)
	if err != nil {
		return nil, NewAPIError("ListTags", err.Error(), 500)
	}

	tagItems := tagutil.ToResponse(tags)
	var items interface{}
	if len(tagItems) > 0 {
		tagSlice := make([]interface{}, len(tagItems))
		for i, t := range tagItems {
			tagSlice[i] = t
		}
		items = protocol.XMLElements{ElementName: "Tag", Items: tagSlice}
	}

	return map[string]interface{}{
		"Tags": map[string]interface{}{
			"Items": items,
		},
	}, nil
}

// TagResource tags a CloudFront resource.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_TagResource.html
func (s *CloudFrontService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetStringParam(req.Parameters, "Resource")
	if arn == "" {
		arn = request.GetStringParam(req.Parameters, "ResourceARN")
	}
	if arn == "" {
		return nil, NewAPIError("InvalidArgument", "Resource is required", 400)
	}

	tags := tagutil.ParseTags(req.Parameters, "Tags")
	if len(tags) == 0 {
		tags = tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.tags.TagResource(arn, tags); err != nil {
		return nil, NewAPIError("TagResource", err.Error(), 500)
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from a CloudFront resource.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_UntagResource.html
func (s *CloudFrontService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetStringParam(req.Parameters, "Resource")
	if arn == "" {
		arn = request.GetStringParam(req.Parameters, "ResourceARN")
	}
	if arn == "" {
		return nil, NewAPIError("InvalidArgument", "Resource is required", 400)
	}

	tagKeys := tagutil.ParseTagKeysWithQueryFallback(req.Parameters, "TagKeys")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if len(tagKeys) > 0 {
		if err := store.tags.UntagResource(arn, tagKeys); err != nil {
			return nil, NewAPIError("UntagResource", err.Error(), 500)
		}
	}

	return response.EmptyResponse(), nil
}
