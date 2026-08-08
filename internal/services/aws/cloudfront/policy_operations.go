package cloudfront

import (
	"context"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	arnutil "vorpalstacks/internal/utils/aws/arn"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
	"vorpalstacks/internal/utils/aws/types"
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
		return nil, awserrors.NewAWSError("InvalidArgument", "Name is required", 400)
	}

	params, err := parseParametersInCacheKey(request.GetMapParam(configMap, "ParametersInCacheKeyAndForwardedToOrigin"))
	if err != nil {
		return nil, err
	}

	config := &cloudfrontstore.CachePolicyConfig{
		Name:                                     name,
		Comment:                                  request.GetStringParam(configMap, "Comment"),
		DefaultTTL:                               int64(request.GetIntParam(configMap, "DefaultTTL")),
		MaxTTL:                                   int64(request.GetIntParam(configMap, "MaxTTL")),
		MinTTL:                                   int64(request.GetIntParam(configMap, "MinTTL")),
		ParametersInCacheKeyParametersInCacheKey: params,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, _ := store.cachePolicies.GetByName(name)
	if existing != nil {
		return nil, awserrors.NewAWSError("CachePolicyAlreadyExists", "Cache policy with this name already exists", 409)
	}

	cachePolicy, err := store.cachePolicies.Create(name, "", config)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ETag": cachePolicy.ETag,
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
		return nil, awserrors.NewAWSError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cachePolicy, err := store.cachePolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchCachePolicy", "Cache policy not found", 404)
		}
		return nil, err
	}

	return map[string]interface{}{
		"ETag": cachePolicy.ETag,
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
		if cp.IsManaged {
			policyType = "managed"
		}
		items = append(items, map[string]interface{}{
			"Type": policyType,
			"CachePolicy": map[string]interface{}{
				"Id":                cp.ID,
				"LastModifiedTime":  cp.ModifiedAt.Format(time.RFC3339),
				"CachePolicyConfig": cp.CachePolicyConfig,
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
		return nil, awserrors.NewAWSError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cachePolicy, err := store.cachePolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchCachePolicy", "Cache policy not found", 404)
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
		return nil, awserrors.NewAWSError("InvalidArgument", "Id is required", 400)
	}

	ifMatch := getIfMatch(req)
	if ifMatch == "" {
		return nil, awserrors.NewAWSError("InvalidIfMatchVersion",
			"The If-Match version is missing or not valid", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.cachePolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchCachePolicy", "Cache policy not found", 404)
		}
		return nil, err
	}

	if ifMatch != "*" && existing.ETag != ifMatch {
		return nil, awserrors.NewAWSError("PreconditionFailed", preconditionFailedETagMsg, 412)
	}

	configMap := request.GetMapParam(req.Parameters, "CachePolicyConfig")
	if configMap == nil {
		configMap = req.Parameters
	}

	params, err := parseParametersInCacheKey(request.GetMapParam(configMap, "ParametersInCacheKeyAndForwardedToOrigin"))
	if err != nil {
		return nil, err
	}

	config := &cloudfrontstore.CachePolicyConfig{
		Name:                                     request.GetStringParam(configMap, "Name"),
		Comment:                                  request.GetStringParam(configMap, "Comment"),
		DefaultTTL:                               int64(request.GetIntParam(configMap, "DefaultTTL")),
		MaxTTL:                                   int64(request.GetIntParam(configMap, "MaxTTL")),
		MinTTL:                                   int64(request.GetIntParam(configMap, "MinTTL")),
		ParametersInCacheKeyParametersInCacheKey: params,
	}

	if config.Name != existing.Name {
		dup, _ := store.cachePolicies.GetByName(config.Name)
		if dup != nil {
			return nil, awserrors.NewAWSError("CachePolicyAlreadyExists", "Cache policy with this name already exists", 409)
		}
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
		return nil, awserrors.NewAWSError("InvalidArgument", "Id is required", 400)
	}

	ifMatch := getIfMatch(req)
	if ifMatch == "" {
		return nil, awserrors.NewAWSError("InvalidIfMatchVersion",
			"The If-Match version is missing or not valid", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.cachePolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchCachePolicy", "Cache policy not found", 404)
		}
		return nil, err
	}

	if ifMatch != "*" && existing.ETag != ifMatch {
		return nil, awserrors.NewAWSError("PreconditionFailed", preconditionFailedETagMsg, 412)
	}

	if isCachePolicyAttached(store, id) {
		return nil, awserrors.NewAWSError("CachePolicyInUse",
			"Cannot delete this cache policy because it is attached to one or more distributions", 409)
	}

	err = store.cachePolicies.Delete(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchCachePolicy", "Cache policy not found", 404)
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
		return nil, awserrors.NewAWSError("InvalidArgument", "Name is required", 400)
	}

	cookiesCfg, err := parseCookiesConfig(request.GetMapParam(configMap, "CookiesConfig"))
	if err != nil {
		return nil, err
	}
	headersCfg, err := parseORPHeadersConfig(request.GetMapParam(configMap, "HeadersConfig"))
	if err != nil {
		return nil, err
	}
	queryStringsCfg, err := parseORPQueryStringsConfig(request.GetMapParam(configMap, "QueryStringsConfig"))
	if err != nil {
		return nil, err
	}

	config := &cloudfrontstore.OriginRequestPolicyConfig{
		Name:               name,
		Comment:            request.GetStringParam(configMap, "Comment"),
		CookiesConfig:      cookiesCfg,
		HeadersConfig:      headersCfg,
		QueryStringsConfig: queryStringsCfg,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, _ := store.originRequestPolicies.GetByName(name)
	if existing != nil {
		return nil, awserrors.NewAWSError("OriginRequestPolicyAlreadyExists", "Origin request policy with this name already exists", 409)
	}

	policy, err := store.originRequestPolicies.Create(name, "", config)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ETag": policy.ETag,
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
		return nil, awserrors.NewAWSError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := store.originRequestPolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchOriginRequestPolicy", "Origin request policy not found", 404)
		}
		return nil, err
	}

	return map[string]interface{}{
		"ETag": policy.ETag,
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
		if p.IsManaged {
			policyType = "managed"
		}
		items = append(items, map[string]interface{}{
			"Type": policyType,
			"OriginRequestPolicy": map[string]interface{}{
				"Id":                        p.ID,
				"LastModifiedTime":          p.ModifiedAt.Format(time.RFC3339),
				"OriginRequestPolicyConfig": p.OriginRequestPolicyConfig,
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
		return nil, awserrors.NewAWSError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := store.originRequestPolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchOriginRequestPolicy", "Origin request policy not found", 404)
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
		return nil, awserrors.NewAWSError("InvalidArgument", "Id is required", 400)
	}

	ifMatch := getIfMatch(req)
	if ifMatch == "" {
		return nil, awserrors.NewAWSError("InvalidIfMatchVersion",
			"The If-Match version is missing or not valid", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.originRequestPolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchOriginRequestPolicy", "Origin request policy not found", 404)
		}
		return nil, err
	}

	if ifMatch != "*" && existing.ETag != ifMatch {
		return nil, awserrors.NewAWSError("PreconditionFailed", preconditionFailedETagMsg, 412)
	}

	configMap := request.GetMapParam(req.Parameters, "OriginRequestPolicyConfig")
	if configMap == nil {
		configMap = req.Parameters
	}

	cookiesCfg, err := parseCookiesConfig(request.GetMapParam(configMap, "CookiesConfig"))
	if err != nil {
		return nil, err
	}
	headersCfg, err := parseORPHeadersConfig(request.GetMapParam(configMap, "HeadersConfig"))
	if err != nil {
		return nil, err
	}
	queryStringsCfg, err := parseORPQueryStringsConfig(request.GetMapParam(configMap, "QueryStringsConfig"))
	if err != nil {
		return nil, err
	}

	config := &cloudfrontstore.OriginRequestPolicyConfig{
		Name:               request.GetStringParam(configMap, "Name"),
		Comment:            request.GetStringParam(configMap, "Comment"),
		CookiesConfig:      cookiesCfg,
		HeadersConfig:      headersCfg,
		QueryStringsConfig: queryStringsCfg,
	}

	if config.Name != existing.Name {
		dup, _ := store.originRequestPolicies.GetByName(config.Name)
		if dup != nil {
			return nil, awserrors.NewAWSError("OriginRequestPolicyAlreadyExists", "Origin request policy with this name already exists", 409)
		}
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
		return nil, awserrors.NewAWSError("InvalidArgument", "Id is required", 400)
	}

	ifMatch := getIfMatch(req)
	if ifMatch == "" {
		return nil, awserrors.NewAWSError("InvalidIfMatchVersion",
			"The If-Match version is missing or not valid", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.originRequestPolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchOriginRequestPolicy", "Origin request policy not found", 404)
		}
		return nil, err
	}

	if ifMatch != "*" && existing.ETag != ifMatch {
		return nil, awserrors.NewAWSError("PreconditionFailed", preconditionFailedETagMsg, 412)
	}

	if isOriginRequestPolicyAttached(store, id) {
		return nil, awserrors.NewAWSError("OriginRequestPolicyInUse",
			"Cannot delete this origin request policy because it is attached to one or more distributions", 409)
	}

	err = store.originRequestPolicies.Delete(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchOriginRequestPolicy", "Origin request policy not found", 404)
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
		return nil, awserrors.NewAWSError("InvalidArgument", "Resource is required", 400)
	}

	resourceKey := arn
	if !strings.HasPrefix(strings.ToLower(arn), "arn:") {
		resourceKey = arnutil.NewARNBuilder("", "").CloudFront().Distribution(arn)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tags, err := store.tags.ListTagsForResource(resourceKey)
	if err != nil {
		return nil, awserrors.NewAWSError("InternalError", err.Error(), 500)
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
		return nil, awserrors.NewAWSError("InvalidArgument", "Resource is required", 400)
	}

	resourceKey := arn
	if !strings.HasPrefix(strings.ToLower(arn), "arn:") {
		resourceKey = arnutil.NewARNBuilder("", "").CloudFront().Distribution(arn)
	}

	var tags []types.Tag
	tagsMap := request.GetMapParam(req.Parameters, "Tags")
	if tagsMap != nil {
		tags = parseXMLTags(tagsMap)
	}
	if len(tags) == 0 {
		parsedTags := tagutil.ParseTags(req.Parameters, "Tags")
		for _, t := range parsedTags {
			tags = append(tags, types.Tag(t))
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.tags.Tag(resourceKey, tags); err != nil {
		return nil, awserrors.NewAWSError("InternalError", err.Error(), 500)
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
		return nil, awserrors.NewAWSError("InvalidArgument", "Resource is required", 400)
	}

	resourceKey := arn
	if !strings.HasPrefix(strings.ToLower(arn), "arn:") {
		resourceKey = arnutil.NewARNBuilder("", "").CloudFront().Distribution(arn)
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
		if err := store.tags.Untag(resourceKey, tagKeys); err != nil {
			return nil, awserrors.NewAWSError("InternalError", err.Error(), 500)
		}
	}

	return response.EmptyResponse(), nil
}

func parseParametersInCacheKey(m map[string]interface{}) (*cloudfrontstore.ParametersInCacheKey, error) {
	if m == nil {
		return nil, nil
	}
	qsCfg, err := parseCachePolicyQueryStringConfig(request.GetMapParam(m, "QueryStringsConfig"))
	if err != nil {
		return nil, err
	}
	cookieCfg, err := parseCachePolicyCookieConfig(request.GetMapParam(m, "CookiesConfig"))
	if err != nil {
		return nil, err
	}
	headerCfg, err := parseCachePolicyHeaderConfig(request.GetMapParam(m, "HeadersConfig"))
	if err != nil {
		return nil, err
	}
	return &cloudfrontstore.ParametersInCacheKey{
		EnableAcceptEncodingGzip:   request.GetBoolParam(m, "EnableAcceptEncodingGzip"),
		EnableAcceptEncodingBrotli: request.GetBoolParam(m, "EnableAcceptEncodingBrotli"),
		QueryStringsConfig:         qsCfg,
		CookiesConfig:              cookieCfg,
		HeadersConfig:              headerCfg,
	}, nil
}

func parseCachePolicyQueryStringConfig(m map[string]interface{}) (*cloudfrontstore.QueryStringConfig, error) {
	if m == nil {
		return nil, nil
	}
	behavior := request.GetStringParam(m, "QueryStringBehavior")
	if err := validateBehavior("QueryStringBehavior", behavior, isValidCachePolicyQueryStringBehavior); err != nil {
		return nil, err
	}
	cfg := &cloudfrontstore.QueryStringConfig{
		QueryStringBehavior: behavior,
	}
	if qsMap := request.GetMapParam(m, "QueryStrings"); qsMap != nil {
		cfg.QueryStrings = &cloudfrontstore.QueryStrings{
			Quantity: request.GetIntParam(qsMap, "Quantity"),
		}
		parseStringItems(qsMap, "Items", &cfg.QueryStrings.Items)
	}
	return cfg, nil
}

func parseCachePolicyCookieConfig(m map[string]interface{}) (*cloudfrontstore.CookieConfig, error) {
	if m == nil {
		return nil, nil
	}
	behavior := request.GetStringParam(m, "CookieBehavior")
	if err := validateBehavior("CookieBehavior", behavior, isValidCachePolicyCookieBehavior); err != nil {
		return nil, err
	}
	cfg := &cloudfrontstore.CookieConfig{
		CookieBehavior: behavior,
	}
	if cMap := request.GetMapParam(m, "Cookies"); cMap != nil {
		cfg.Cookies = &cloudfrontstore.Cookies{
			Quantity: request.GetIntParam(cMap, "Quantity"),
		}
		parseStringItems(cMap, "Items", &cfg.Cookies.Items)
	}
	return cfg, nil
}

func parseCachePolicyHeaderConfig(m map[string]interface{}) (*cloudfrontstore.HeaderConfig, error) {
	if m == nil {
		return nil, nil
	}
	behavior := request.GetStringParam(m, "HeaderBehavior")
	if err := validateBehavior("HeaderBehavior", behavior, isValidCachePolicyHeaderBehavior); err != nil {
		return nil, err
	}
	cfg := &cloudfrontstore.HeaderConfig{
		HeaderBehavior: behavior,
	}
	if hMap := request.GetMapParam(m, "Headers"); hMap != nil {
		cfg.Headers = &cloudfrontstore.Headers{
			Quantity: request.GetIntParam(hMap, "Quantity"),
		}
		parseStringItems(hMap, "Items", &cfg.Headers.Items)
	}
	return cfg, nil
}
func parseCookiesConfig(m map[string]interface{}) (*cloudfrontstore.CookiesConfig, error) {
	if m == nil {
		return nil, nil
	}
	behavior := request.GetStringParam(m, "CookieBehavior")
	if err := validateBehavior("CookieBehavior", behavior, isValidORPCookieBehavior); err != nil {
		return nil, err
	}
	cfg := &cloudfrontstore.CookiesConfig{
		CookieBehavior: behavior,
	}
	if cMap := request.GetMapParam(m, "Cookies"); cMap != nil {
		cfg.Cookies = &cloudfrontstore.Cookies{
			Quantity: request.GetIntParam(cMap, "Quantity"),
		}
		parseStringItems(cMap, "Items", &cfg.Cookies.Items)
	}
	return cfg, nil
}

func parseORPHeadersConfig(m map[string]interface{}) (*cloudfrontstore.HeadersConfig, error) {
	if m == nil {
		return nil, nil
	}
	behavior := request.GetStringParam(m, "HeaderBehavior")
	if err := validateBehavior("HeaderBehavior", behavior, isValidORPHeaderBehavior); err != nil {
		return nil, err
	}
	cfg := &cloudfrontstore.HeadersConfig{
		HeaderBehavior: behavior,
	}
	if hMap := request.GetMapParam(m, "Headers"); hMap != nil {
		cfg.Headers = &cloudfrontstore.Headers{
			Quantity: request.GetIntParam(hMap, "Quantity"),
		}
		parseStringItems(hMap, "Items", &cfg.Headers.Items)
	}
	return cfg, nil
}

func parseORPQueryStringsConfig(m map[string]interface{}) (*cloudfrontstore.QueryStringsConfig, error) {
	if m == nil {
		return nil, nil
	}
	behavior := request.GetStringParam(m, "QueryStringBehavior")
	if err := validateBehavior("QueryStringBehavior", behavior, isValidORPQueryStringBehavior); err != nil {
		return nil, err
	}
	cfg := &cloudfrontstore.QueryStringsConfig{
		QueryStringBehavior: behavior,
	}
	if qsMap := request.GetMapParam(m, "QueryStrings"); qsMap != nil {
		cfg.QueryStrings = &cloudfrontstore.QueryStrings{
			Quantity: request.GetIntParam(qsMap, "Quantity"),
		}
		parseStringItems(qsMap, "Items", &cfg.QueryStrings.Items)
	}
	return cfg, nil
}

func parseStringItems(m map[string]interface{}, key string, out *[]string) {
	if items, ok := m[key]; ok {
		switch v := items.(type) {
		case []interface{}:
			for _, item := range v {
				if s, ok := item.(string); ok {
					*out = append(*out, s)
				}
			}
		case map[string]interface{}:
			for _, val := range v {
				if arr, ok := val.([]interface{}); ok {
					for _, item := range arr {
						if s, ok := item.(string); ok {
							*out = append(*out, s)
						}
					}
				} else if s, ok := val.(string); ok {
					*out = append(*out, s)
				}
			}
		}
	}
}

// isCachePolicyAttached checks whether any distribution references the cache
// policy ID in its DefaultCacheBehaviour or additional CacheBehaviours.
func isCachePolicyAttached(store *cloudfrontStores, policyID string) bool {
	result, err := store.distributions.List("", 10000)
	if err != nil {
		return false
	}
	for _, dist := range result.Distributions {
		if dist.DistributionConfig == nil {
			continue
		}
		cfg := dist.DistributionConfig
		if cfg.DefaultCacheBehavior != nil && cfg.DefaultCacheBehavior.CachePolicyId == policyID {
			return true
		}
		if cfg.CacheBehaviors != nil {
			for _, cb := range cfg.CacheBehaviors.Items {
				if cb != nil && cb.CachePolicyId == policyID {
					return true
				}
			}
		}
	}
	return false
}

// isOriginRequestPolicyAttached checks whether any distribution references
// the origin request policy ID in its cache behaviours.
func isOriginRequestPolicyAttached(store *cloudfrontStores, policyID string) bool {
	result, err := store.distributions.List("", 10000)
	if err != nil {
		return false
	}
	for _, dist := range result.Distributions {
		if dist.DistributionConfig == nil {
			continue
		}
		cfg := dist.DistributionConfig
		if cfg.DefaultCacheBehavior != nil && cfg.DefaultCacheBehavior.OriginRequestPolicyId == policyID {
			return true
		}
		if cfg.CacheBehaviors != nil {
			for _, cb := range cfg.CacheBehaviors.Items {
				if cb != nil && cb.OriginRequestPolicyId == policyID {
					return true
				}
			}
		}
	}
	return false
}

// isResponseHeadersPolicyAttached checks whether any distribution references
// the response headers policy ID in its cache behaviours.
func isResponseHeadersPolicyAttached(store *cloudfrontStores, policyID string) bool {
	result, err := store.distributions.List("", 10000)
	if err != nil {
		return false
	}
	for _, dist := range result.Distributions {
		if dist.DistributionConfig == nil {
			continue
		}
		cfg := dist.DistributionConfig
		if cfg.DefaultCacheBehavior != nil && cfg.DefaultCacheBehavior.ResponseHeadersPolicyId == policyID {
			return true
		}
		if cfg.CacheBehaviors != nil {
			for _, cb := range cfg.CacheBehaviors.Items {
				if cb != nil && cb.ResponseHeadersPolicyId == policyID {
					return true
				}
			}
		}
	}
	return false
}

// isOriginAccessControlAttached checks whether any distribution references
// the origin access control ID in its origins.
func isOriginAccessControlAttached(store *cloudfrontStores, oacID string) bool {
	result, err := store.distributions.List("", 10000)
	if err != nil {
		return false
	}
	for _, dist := range result.Distributions {
		if dist.DistributionConfig == nil {
			continue
		}
		for _, origin := range dist.DistributionConfig.Origins.Items {
			if origin != nil && origin.OriginAccessControlId == oacID {
				return true
			}
		}
	}
	return false
}
