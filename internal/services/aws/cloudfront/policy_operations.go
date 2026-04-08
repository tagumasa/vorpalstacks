package cloudfront

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
	"vorpalstacks/internal/store/aws/common"
)

// CreateCachePolicy creates a cache policy.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_CreateCachePolicy.html
func (s *CloudFrontService) CreateCachePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	configMap := request.GetMapParam(req.Parameters, "CachePolicyConfig")
	if configMap == nil {
		configMap = req.Parameters
	}

	name := request.GetStringParam(configMap, "Name")
	if name == "" {
		return nil, NewAPIError("InvalidArgument", "Name is required", 400)
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
		policyType := "custom"
		if cp.ETag == "managed" {
			policyType = "managed"
		}
		items = append(items, map[string]interface{}{
			"Type": policyType,
			"CachePolicy": map[string]interface{}{
				"Id":               cp.ID,
				"LastModifiedTime": cp.ModifiedAt.Format(time.RFC3339),
				"CachePolicyConfig": map[string]interface{}{
					"Name":       cp.CachePolicyConfig.Name,
					"Comment":    cp.CachePolicyConfig.Comment,
					"DefaultTTL": cp.CachePolicyConfig.DefaultTTL,
					"MaxTTL":     cp.CachePolicyConfig.MaxTTL,
					"MinTTL":     cp.CachePolicyConfig.MinTTL,
				},
			},
		})
	}

	return map[string]interface{}{
		"CachePolicyList": map[string]interface{}{
			"Items":       protocol.XMLElements{ElementName: "CachePolicySummary", Items: items},
			"IsTruncated": result.IsTruncated,
			"NextMarker":  result.NextMarker,
			"Quantity":    len(result.CachePolicies),
			"MaxItems":    maxItems,
		},
	}, nil
}

// GetCachePolicyConfig returns the configuration of a cache policy.
func (s *CloudFrontService) GetCachePolicyConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
		"CachePolicyConfig": cachePolicy.CachePolicyConfig,
		"ETag":              cachePolicy.ETag,
	}, nil
}

// UpdateCachePolicy updates a cache policy.
func (s *CloudFrontService) UpdateCachePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidArgument", "Id is required", 400)
	}

	ifMatch := getIfMatch(req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.cachePolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, NewAPIError("NoSuchCachePolicy", "Cache policy not found", 404)
		}
		return nil, err
	}

	if ifMatch != "" && ifMatch != "*" && existing.ETag != ifMatch {
		return nil, NewAPIError("PreconditionFailed", "The If-Match header does not match the current ETag", 412)
	}

	configMap := request.GetMapParam(req.Parameters, "CachePolicyConfig")
	if configMap == nil {
		configMap = req.Parameters
	}

	config := &cloudfrontstore.CachePolicyConfig{
		Name:       request.GetStringParam(configMap, "Name"),
		Comment:    request.GetStringParam(configMap, "Comment"),
		DefaultTTL: int64(request.GetIntParam(configMap, "DefaultTTL")),
		MaxTTL:     int64(request.GetIntParam(configMap, "MaxTTL")),
		MinTTL:     int64(request.GetIntParam(configMap, "MinTTL")),
	}

	cachePolicy, err := store.cachePolicies.Update(id, config)
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
		"ETag": cachePolicy.ETag,
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
	configMap := request.GetMapParam(req.Parameters, "OriginRequestPolicyConfig")
	if configMap == nil {
		configMap = req.Parameters
	}

	name := request.GetStringParam(configMap, "Name")
	if name == "" {
		return nil, NewAPIError("InvalidArgument", "Name is required", 400)
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
		policyType := "custom"
		if p.ETag == "managed" {
			policyType = "managed"
		}
		items = append(items, map[string]interface{}{
			"Type": policyType,
			"OriginRequestPolicy": map[string]interface{}{
				"Id":               p.ID,
				"LastModifiedTime": p.ModifiedAt.Format(time.RFC3339),
				"OriginRequestPolicyConfig": map[string]interface{}{
					"Name":    p.OriginRequestPolicyConfig.Name,
					"Comment": p.OriginRequestPolicyConfig.Comment,
				},
			},
		})
	}

	return map[string]interface{}{
		"OriginRequestPolicyList": map[string]interface{}{
			"Items":       protocol.XMLElements{ElementName: "OriginRequestPolicySummary", Items: items},
			"IsTruncated": result.IsTruncated,
			"NextMarker":  result.NextMarker,
			"Quantity":    len(result.OriginRequestPolicies),
			"MaxItems":    maxItems,
		},
	}, nil
}

// GetOriginRequestPolicyConfig returns the configuration of an origin request policy.
func (s *CloudFrontService) GetOriginRequestPolicyConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
		"OriginRequestPolicyConfig": policy.OriginRequestPolicyConfig,
		"ETag":                      policy.ETag,
	}, nil
}

// UpdateOriginRequestPolicy updates an origin request policy.
func (s *CloudFrontService) UpdateOriginRequestPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidArgument", "Id is required", 400)
	}

	ifMatch := getIfMatch(req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.originRequestPolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, NewAPIError("NoSuchOriginRequestPolicy", "Origin request policy not found", 404)
		}
		return nil, err
	}

	if ifMatch != "" && ifMatch != "*" && existing.ETag != ifMatch {
		return nil, NewAPIError("PreconditionFailed", "The If-Match header does not match the current ETag", 412)
	}

	configMap := request.GetMapParam(req.Parameters, "OriginRequestPolicyConfig")
	if configMap == nil {
		configMap = req.Parameters
	}

	config := &cloudfrontstore.OriginRequestPolicyConfig{
		Name:    request.GetStringParam(configMap, "Name"),
		Comment: request.GetStringParam(configMap, "Comment"),
	}

	policy, err := store.originRequestPolicies.Update(id, config)
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
		"ETag": policy.ETag,
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

	resourceKey := arn
	if !strings.HasPrefix(strings.ToLower(arn), "arn:") {
		resourceKey = "arn:aws:cloudfront:::distribution/" + arn
	}
	fmt.Fprintf(os.Stderr, "[DBG ListTagsForResource] arn=%s resourceKey=%s\n", arn, resourceKey)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tags, err := store.tags.ListTagsForResource(resourceKey)
	if err != nil {
		return nil, NewAPIError("ListTags", err.Error(), 500)
	}
	fmt.Fprintf(os.Stderr, "[DBG ListTagsForResource] found %d tags for resourceKey=%s\n", len(tags), resourceKey)
	for _, t := range tags {
		fmt.Fprintf(os.Stderr, "[DBG ListTagsForResource]   tag: %s=%s\n", t.Key, t.Value)
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

	resourceKey := arn
	if !strings.HasPrefix(strings.ToLower(arn), "arn:") {
		resourceKey = "arn:aws:cloudfront:::distribution/" + arn
	}

	var tags []common.Tag
	tagsMap := request.GetMapParam(req.Parameters, "Tags")
	if tagsMap != nil {
		if tagInner, ok := tagsMap["Tag"]; ok {
			switch tv := tagInner.(type) {
			case []interface{}:
				for _, t := range tv {
					if m, ok := t.(map[string]interface{}); ok {
						tags = append(tags, common.Tag{Key: request.GetStringParam(m, "Key"), Value: request.GetStringParam(m, "Value")})
					}
				}
			case map[string]interface{}:
				tags = append(tags, common.Tag{Key: request.GetStringParam(tv, "Key"), Value: request.GetStringParam(tv, "Value")})
			}
		} else if items := tagsMap["Items"]; items != nil {
			switch v := items.(type) {
			case map[string]interface{}:
				if tagItems, ok := v["Tag"]; ok {
					switch tv := tagItems.(type) {
					case []interface{}:
						for _, t := range tv {
							if m, ok := t.(map[string]interface{}); ok {
								tags = append(tags, common.Tag{Key: request.GetStringParam(m, "Key"), Value: request.GetStringParam(m, "Value")})
							}
						}
					case map[string]interface{}:
						tags = append(tags, common.Tag{Key: request.GetStringParam(tv, "Key"), Value: request.GetStringParam(tv, "Value")})
					}
				}
			}
		}
	}
	if len(tags) == 0 {
		var keys []string
		for k := range req.Parameters {
			keys = append(keys, k)
		}
		logs.Debug("CloudFront TagResource: no tags parsed", logs.String("arn", arn), logs.String("param_keys", strings.Join(keys, ",")))
		parsedTags := tagutil.ParseTags(req.Parameters, "Tags")
		for _, t := range parsedTags {
			tags = append(tags, common.Tag{Key: t.Key, Value: t.Value})
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.tags.TagResource(resourceKey, tags); err != nil {
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

	resourceKey := arn
	if !strings.HasPrefix(strings.ToLower(arn), "arn:") {
		resourceKey = "arn:aws:cloudfront:::distribution/" + arn
	}

	tagKeys := tagutil.ParseTagKeysWithQueryFallback(req.Parameters, "TagKeys")
	if len(tagKeys) == 0 {
		if tagKeysMap := request.GetMapParam(req.Parameters, "TagKeys"); tagKeysMap != nil {
			if keyVal, ok := tagKeysMap["Key"]; ok {
				switch kv := keyVal.(type) {
				case string:
					tagKeys = append(tagKeys, kv)
				case []interface{}:
					for _, k := range kv {
						if s, ok := k.(string); ok {
							tagKeys = append(tagKeys, s)
						}
					}
				}
			} else if items := tagKeysMap["Items"]; items != nil {
				switch v := items.(type) {
				case map[string]interface{}:
					if keyItems, ok := v["Key"]; ok {
						switch kv := keyItems.(type) {
						case []interface{}:
							for _, k := range kv {
								if s, ok := k.(string); ok {
									tagKeys = append(tagKeys, s)
								}
							}
						case string:
							tagKeys = append(tagKeys, kv)
						}
					}
				}
			}
		}
	}
	if len(tagKeys) == 0 {
		if items := request.GetMapParam(req.Parameters, "Items"); items != nil {
			if keyItems, ok := items["Key"]; ok {
				switch kv := keyItems.(type) {
				case []interface{}:
					for _, k := range kv {
						if s, ok := k.(string); ok {
							tagKeys = append(tagKeys, s)
						}
					}
				case string:
					tagKeys = append(tagKeys, kv)
				}
			}
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if len(tagKeys) > 0 {
		if err := store.tags.UntagResource(resourceKey, tagKeys); err != nil {
			return nil, NewAPIError("UntagResource", err.Error(), 500)
		}
	}

	return response.EmptyResponse(), nil
}
