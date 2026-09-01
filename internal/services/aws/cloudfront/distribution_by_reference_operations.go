package cloudfront

import (
	"context"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// listDistributionsByReference implements the four
// ListDistributionsByXxxId operations that return the IDs of every
// distribution referencing a cache policy, origin request policy,
// response headers policy, or key group. The response uses the shared
// DistributionIdList shape. Every operation models a 404 error for an
// unknown target ID, so requireTarget verifies the referenced resource
// exists before distributions are collected.
func (s *CloudFrontService) listDistributionsByReference(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest, param string,
	requireTarget func(*cloudfrontStores, string) error,
	matches func(*cloudfrontstore.DistributionConfig) bool) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, param)
	if err := requireReferenceIdCore(id, param); err != nil {
		return nil, err
	}

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := resolveListMaxItems(request.GetIntParam(req.Parameters, "MaxItems"))

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := requireTarget(store, id); err != nil {
		return nil, err
	}

	matched, err := collectDistributions(store, func(dist *cloudfrontstore.Distribution) bool {
		return dist.DistributionConfig != nil && matches(dist.DistributionConfig)
	})
	if err != nil {
		return nil, awserrors.NewAWSError("InternalError", "Failed to list distributions by reference: "+err.Error(), 500)
	}

	skipCount := 0
	if marker != "" {
		for i, distID := range matched {
			if distID == marker {
				skipCount = i + 1
				break
			}
		}
	}

	paged := matched[skipCount:]
	isTruncated := len(paged) > maxItems
	if isTruncated {
		paged = paged[:maxItems]
	}

	nextMarker := ""
	if isTruncated && len(paged) > 0 {
		nextMarker = paged[len(paged)-1]
	}

	items := make([]interface{}, len(paged))
	for i, distID := range paged {
		items[i] = distID
	}

	idList := map[string]interface{}{
		"Marker":      marker,
		"MaxItems":    maxItems,
		"IsTruncated": isTruncated,
		"Quantity":    len(paged),
		"Items":       protocol.XMLElements{ElementName: "DistributionId", Items: items},
	}
	if nextMarker != "" {
		idList["NextMarker"] = nextMarker
	}
	return map[string]interface{}{"DistributionIdList": idList}, nil
}

// behaviourField extracts the value of a policy reference from the
// default cache behaviour and the additional cache behaviours of a
// distribution config.
func behaviourField(cfg *cloudfrontstore.DistributionConfig, get func(*cloudfrontstore.CacheBehavior) string) []string {
	var values []string
	if cfg.DefaultCacheBehavior != nil {
		values = append(values, get(cfg.DefaultCacheBehavior))
	}
	if cfg.CacheBehaviors != nil {
		for _, cb := range cfg.CacheBehaviors.Items {
			if cb != nil {
				values = append(values, get(cb))
			}
		}
	}
	return values
}

// requireCachePolicyForList rejects the by-cache-policy listing when the
// referenced cache policy does not exist. The API models this as the
// NoSuchCachePolicy error with HTTP status 404.
func requireCachePolicyForList(stores *cloudfrontStores, id string) error {
	if _, err := stores.cachePolicies.Get(id); err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchCachePolicy", "The cache policy does not exist.", 404)
		}
		return err
	}
	return nil
}

// requireOriginRequestPolicyForList rejects the by-origin-request-policy
// listing when the referenced origin request policy does not exist. The
// API models this as the NoSuchOriginRequestPolicy error with HTTP
// status 404.
func requireOriginRequestPolicyForList(stores *cloudfrontStores, id string) error {
	if _, err := stores.originRequestPolicies.Get(id); err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchOriginRequestPolicy", "The origin request policy does not exist.", 404)
		}
		return err
	}
	return nil
}

// requireResponseHeadersPolicyForList rejects the by-response-headers-
// policy listing when the referenced response headers policy does not
// exist. The API models this as the NoSuchResponseHeadersPolicy error
// with HTTP status 404.
func requireResponseHeadersPolicyForList(stores *cloudfrontStores, id string) error {
	if _, err := stores.responseHeadersPolicies.Get(id); err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchResponseHeadersPolicy", "The response headers policy does not exist.", 404)
		}
		return err
	}
	return nil
}

// requireKeyGroupForList rejects the by-key-group listing when the
// referenced key group does not exist. The API models this as the
// NoSuchResource error with HTTP status 404.
func requireKeyGroupForList(stores *cloudfrontStores, id string) error {
	if _, err := stores.keyGroups.Get(id); err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchResource", "A resource that was specified is not valid: "+id, 404)
		}
		return err
	}
	return nil
}

// ListDistributionsByCachePolicyId lists the IDs of distributions with a
// cache behaviour associated with the specified cache policy.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_ListDistributionsByCachePolicyId.html
func (s *CloudFrontService) ListDistributionsByCachePolicyId(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.listDistributionsByReference(ctx, reqCtx, req, "CachePolicyId", requireCachePolicyForList, func(cfg *cloudfrontstore.DistributionConfig) bool {
		for _, v := range behaviourField(cfg, func(cb *cloudfrontstore.CacheBehavior) string { return cb.CachePolicyId }) {
			if v == request.GetStringParam(req.Parameters, "CachePolicyId") {
				return true
			}
		}
		return false
	})
}

// ListDistributionsByOriginRequestPolicyId lists the IDs of distributions
// with a cache behaviour associated with the specified origin request
// policy.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_ListDistributionsByOriginRequestPolicyId.html
func (s *CloudFrontService) ListDistributionsByOriginRequestPolicyId(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.listDistributionsByReference(ctx, reqCtx, req, "OriginRequestPolicyId", requireOriginRequestPolicyForList, func(cfg *cloudfrontstore.DistributionConfig) bool {
		for _, v := range behaviourField(cfg, func(cb *cloudfrontstore.CacheBehavior) string { return cb.OriginRequestPolicyId }) {
			if v == request.GetStringParam(req.Parameters, "OriginRequestPolicyId") {
				return true
			}
		}
		return false
	})
}

// ListDistributionsByResponseHeadersPolicyId lists the IDs of
// distributions with a cache behaviour associated with the specified
// response headers policy.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_ListDistributionsByResponseHeadersPolicyId.html
func (s *CloudFrontService) ListDistributionsByResponseHeadersPolicyId(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.listDistributionsByReference(ctx, reqCtx, req, "ResponseHeadersPolicyId", requireResponseHeadersPolicyForList, func(cfg *cloudfrontstore.DistributionConfig) bool {
		for _, v := range behaviourField(cfg, func(cb *cloudfrontstore.CacheBehavior) string { return cb.ResponseHeadersPolicyId }) {
			if v == request.GetStringParam(req.Parameters, "ResponseHeadersPolicyId") {
				return true
			}
		}
		return false
	})
}

// ListDistributionsByKeyGroup lists the IDs of distributions whose cache
// behaviours trust the specified key group.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_ListDistributionsByKeyGroup.html
func (s *CloudFrontService) ListDistributionsByKeyGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	keyGroupID := request.GetStringParam(req.Parameters, "KeyGroupId")
	return s.listDistributionsByReference(ctx, reqCtx, req, "KeyGroupId", requireKeyGroupForList, func(cfg *cloudfrontstore.DistributionConfig) bool {
		for _, v := range behaviourField(cfg, func(cb *cloudfrontstore.CacheBehavior) string {
			if cb.TrustedKeyGroups == nil {
				return ""
			}
			for _, kgID := range cb.TrustedKeyGroups.Items {
				if kgID == keyGroupID {
					return kgID
				}
			}
			return ""
		}) {
			if v == keyGroupID {
				return true
			}
		}
		return false
	})
}
