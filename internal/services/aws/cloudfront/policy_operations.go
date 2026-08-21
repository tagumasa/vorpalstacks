package cloudfront

import (
	"context"
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// CreateCachePolicy creates a cache policy.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_CreateCachePolicy.html
func (s *CloudFrontService) CreateCachePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	cachePolicy, err := s.createCachePolicyCore(store, CreateCachePolicyInput{Config: config})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ETag":        cachePolicy.ETag,
		"CachePolicy": formatCachePolicy(cachePolicy),
	}, nil
}

// GetCachePolicy returns a cache policy.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_GetCachePolicy.html
func (s *CloudFrontService) GetCachePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cachePolicy, err := s.getCachePolicyCore(store, request.GetStringParam(req.Parameters, "Id"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ETag":        cachePolicy.ETag,
		"CachePolicy": formatCachePolicy(cachePolicy),
	}, nil
}

// ListCachePolicies lists cache policies.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_ListCachePolicies.html
func (s *CloudFrontService) ListCachePolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listCachePoliciesCore(store, ListPoliciesInput{
		Marker:     request.GetStringParam(req.Parameters, "Marker"),
		MaxItems:   request.GetIntParam(req.Parameters, "MaxItems"),
		TypeFilter: request.GetStringParam(req.Parameters, "Type"),
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Policies))
	for _, cp := range result.Policies {
		policyType := "custom"
		if cp.IsManaged {
			policyType = "managed"
		}
		items = append(items, map[string]interface{}{
			"Type":        policyType,
			"CachePolicy": formatCachePolicySummary(cp),
		})
	}

	cpList := map[string]interface{}{
		"Items":    protocol.XMLElements{ElementName: "CachePolicySummary", Items: items},
		"Quantity": len(items),
		"MaxItems": result.EffectiveMaxItems,
	}
	if result.NextMarker != "" {
		cpList["NextMarker"] = result.NextMarker
	}
	return map[string]interface{}{"CachePolicyList": cpList}, nil
}

// GetCachePolicyConfig returns the configuration of a cache policy.
func (s *CloudFrontService) GetCachePolicyConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cachePolicy, err := s.getCachePolicyCore(store, request.GetStringParam(req.Parameters, "Id"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"CachePolicyConfig": cachePolicy.CachePolicyConfig,
		"ETag":              cachePolicy.ETag,
	}, nil
}

// UpdateCachePolicy updates a cache policy.
func (s *CloudFrontService) UpdateCachePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	cachePolicy, err := s.updateCachePolicyCore(store, UpdateCachePolicyInput{
		Id:      request.GetStringParam(req.Parameters, "Id"),
		IfMatch: getIfMatch(req),
		Config:  config,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"CachePolicy": formatCachePolicy(cachePolicy),
		"ETag":        cachePolicy.ETag,
	}, nil
}

// DeleteCachePolicy deletes a cache policy.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_DeleteCachePolicy.html
func (s *CloudFrontService) DeleteCachePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteCachePolicyCore(store,
		request.GetStringParam(req.Parameters, "Id"), getIfMatch(req)); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// CreateOriginRequestPolicy creates an origin request policy.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_CreateOriginRequestPolicy.html
func (s *CloudFrontService) CreateOriginRequestPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	config, err := parseOriginRequestPolicyConfig(req)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	policy, err := s.createOriginRequestPolicyCore(store, CreateOriginRequestPolicyInput{Config: config})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ETag":                policy.ETag,
		"OriginRequestPolicy": formatOriginRequestPolicy(policy),
	}, nil
}

// GetOriginRequestPolicy returns an origin request policy.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_GetOriginRequestPolicy.html
func (s *CloudFrontService) GetOriginRequestPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := s.getOriginRequestPolicyCore(store, request.GetStringParam(req.Parameters, "Id"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ETag":                policy.ETag,
		"OriginRequestPolicy": formatOriginRequestPolicy(policy),
	}, nil
}

// ListOriginRequestPolicies lists origin request policies.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_ListOriginRequestPolicies.html
func (s *CloudFrontService) ListOriginRequestPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listOriginRequestPoliciesCore(store, ListPoliciesInput{
		Marker:     request.GetStringParam(req.Parameters, "Marker"),
		MaxItems:   request.GetIntParam(req.Parameters, "MaxItems"),
		TypeFilter: request.GetStringParam(req.Parameters, "Type"),
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Policies))
	for _, p := range result.Policies {
		policyType := "custom"
		if p.IsManaged {
			policyType = "managed"
		}
		items = append(items, map[string]interface{}{
			"Type":                policyType,
			"OriginRequestPolicy": formatOriginRequestPolicySummary(p),
		})
	}

	orpList := map[string]interface{}{
		"Items":    protocol.XMLElements{ElementName: "OriginRequestPolicySummary", Items: items},
		"Quantity": len(items),
		"MaxItems": result.EffectiveMaxItems,
	}
	if result.NextMarker != "" {
		orpList["NextMarker"] = result.NextMarker
	}
	return map[string]interface{}{"OriginRequestPolicyList": orpList}, nil
}

// GetOriginRequestPolicyConfig returns the configuration of an origin request policy.
func (s *CloudFrontService) GetOriginRequestPolicyConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := s.getOriginRequestPolicyCore(store, request.GetStringParam(req.Parameters, "Id"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"OriginRequestPolicyConfig": policy.OriginRequestPolicyConfig,
		"ETag":                      policy.ETag,
	}, nil
}

// UpdateOriginRequestPolicy updates an origin request policy.
func (s *CloudFrontService) UpdateOriginRequestPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	config, err := parseOriginRequestPolicyConfig(req)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	policy, err := s.updateOriginRequestPolicyCore(store, UpdateOriginRequestPolicyInput{
		Id:      request.GetStringParam(req.Parameters, "Id"),
		IfMatch: getIfMatch(req),
		Config:  config,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"OriginRequestPolicy": formatOriginRequestPolicy(policy),
		"ETag":                policy.ETag,
	}, nil
}

// DeleteOriginRequestPolicy deletes an origin request policy.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_DeleteOriginRequestPolicy.html
func (s *CloudFrontService) DeleteOriginRequestPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteOriginRequestPolicyCore(store,
		request.GetStringParam(req.Parameters, "Id"), getIfMatch(req)); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// parseOriginRequestPolicyConfig parses the OriginRequestPolicyConfig
// request payload into the store configuration type.
func parseOriginRequestPolicyConfig(req *request.ParsedRequest) (*cloudfrontstore.OriginRequestPolicyConfig, error) {
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

	return &cloudfrontstore.OriginRequestPolicyConfig{
		Name:               request.GetStringParam(configMap, "Name"),
		Comment:            request.GetStringParam(configMap, "Comment"),
		CookiesConfig:      cookiesCfg,
		HeadersConfig:      headersCfg,
		QueryStringsConfig: queryStringsCfg,
	}, nil
}

// formatCachePolicy renders a cache policy for Get/Create/Update responses.
func formatCachePolicy(cp *cloudfrontstore.CachePolicy) map[string]interface{} {
	return map[string]interface{}{
		"Id":                cp.ID,
		"ARN":               cp.ARN,
		"Name":              cp.Name,
		"CachePolicyConfig": cp.CachePolicyConfig,
		"LastModifiedTime":  cp.ModifiedAt.Format(time.RFC3339),
	}
}

// formatCachePolicySummary renders a cache policy for list summaries.
func formatCachePolicySummary(cp *cloudfrontstore.CachePolicy) map[string]interface{} {
	return map[string]interface{}{
		"Id":                cp.ID,
		"LastModifiedTime":  cp.ModifiedAt.Format(time.RFC3339),
		"CachePolicyConfig": cp.CachePolicyConfig,
	}
}

// formatOriginRequestPolicy renders an origin request policy for
// Get/Create/Update responses.
func formatOriginRequestPolicy(p *cloudfrontstore.OriginRequestPolicy) map[string]interface{} {
	return map[string]interface{}{
		"Id":                        p.ID,
		"ARN":                       p.ARN,
		"Name":                      p.Name,
		"OriginRequestPolicyConfig": p.OriginRequestPolicyConfig,
		"LastModifiedTime":          p.ModifiedAt.Format(time.RFC3339),
	}
}

// formatOriginRequestPolicySummary renders an origin request policy for
// list summaries.
func formatOriginRequestPolicySummary(p *cloudfrontstore.OriginRequestPolicy) map[string]interface{} {
	return map[string]interface{}{
		"Id":                        p.ID,
		"LastModifiedTime":          p.ModifiedAt.Format(time.RFC3339),
		"OriginRequestPolicyConfig": p.OriginRequestPolicyConfig,
	}
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
	if !isValidResourceArn(arn) {
		return nil, awserrors.NewAWSError("InvalidArgument", "Resource must be a CloudFront resource ARN: "+arn, 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tags, err := store.tags.ListTagsForResource(arn)
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
	if !isValidResourceArn(arn) {
		return nil, awserrors.NewAWSError("InvalidArgument", "Resource must be a CloudFront resource ARN: "+arn, 400)
	}

	var tags []tagutil.Tag
	tagsMap := request.GetMapParam(req.Parameters, "Tags")
	if tagsMap != nil {
		tags = parseXMLTags(tagsMap)
	}
	if len(tags) == 0 {
		parsedTags := tagutil.ParseTags(req.Parameters, "Tags")
		for _, t := range parsedTags {
			tags = append(tags, tagutil.Tag(t))
		}
	}
	if len(tags) == 0 {
		return nil, awserrors.NewAWSError("InvalidArgument", "At least one tag is required", 400)
	}
	for _, t := range tags {
		if !isValidTagKey(t.Key) {
			return nil, awserrors.NewAWSError("InvalidArgument", fmt.Sprintf("Invalid tag key: %q", t.Key), 400)
		}
		if !isValidTagValue(t.Value) {
			return nil, awserrors.NewAWSError("InvalidArgument", fmt.Sprintf("Invalid tag value for key %q", t.Key), 400)
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.tags.Tag(arn, tags); err != nil {
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
	if !isValidResourceArn(arn) {
		return nil, awserrors.NewAWSError("InvalidArgument", "Resource must be a CloudFront resource ARN: "+arn, 400)
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

	if len(tagKeys) == 0 {
		return nil, awserrors.NewAWSError("InvalidArgument", "At least one tag key is required", 400)
	}
	for _, k := range tagKeys {
		if !isValidTagKey(k) {
			return nil, awserrors.NewAWSError("InvalidArgument", fmt.Sprintf("Invalid tag key: %q", k), 400)
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.tags.Untag(arn, tagKeys); err != nil {
		return nil, awserrors.NewAWSError("InternalError", err.Error(), 500)
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
// A store failure is reported as an error so deletion can fail closed.
func isCachePolicyAttached(store *cloudfrontStores, policyID string) (bool, error) {
	return scanDistributions(store, func(dist *cloudfrontstore.Distribution) bool {
		if dist.DistributionConfig == nil {
			return false
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
		return false
	})
}

// isOriginRequestPolicyAttached checks whether any distribution references
// the origin request policy ID in its cache behaviours.
func isOriginRequestPolicyAttached(store *cloudfrontStores, policyID string) (bool, error) {
	return scanDistributions(store, func(dist *cloudfrontstore.Distribution) bool {
		if dist.DistributionConfig == nil {
			return false
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
		return false
	})
}

// isResponseHeadersPolicyAttached checks whether any distribution references
// the response headers policy ID in its cache behaviours.
func isResponseHeadersPolicyAttached(store *cloudfrontStores, policyID string) (bool, error) {
	return scanDistributions(store, func(dist *cloudfrontstore.Distribution) bool {
		if dist.DistributionConfig == nil {
			return false
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
		return false
	})
}

// isOriginAccessControlAttached checks whether any distribution references
// the origin access control ID in its origins.
func isOriginAccessControlAttached(store *cloudfrontStores, oacID string) (bool, error) {
	return scanDistributions(store, func(dist *cloudfrontstore.Distribution) bool {
		if dist.DistributionConfig == nil {
			return false
		}
		for _, origin := range dist.DistributionConfig.Origins.Items {
			if origin != nil && origin.OriginAccessControlId == oacID {
				return true
			}
		}
		return false
	})
}
