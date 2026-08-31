package cloudfront

import (
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	types "vorpalstacks/internal/common/tags"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// ---------------------------------------------------------------------------
// Cache policy cores — single validation and persistence path shared by
// every entry point (HTTP API and future admin gRPC handlers).
// ---------------------------------------------------------------------------

// CreateCachePolicyInput carries the parsed cache policy configuration.
type CreateCachePolicyInput struct {
	Config *cloudfrontstore.CachePolicyConfig
}

// UpdateCachePolicyInput carries the parameters for updating a cache policy.
type UpdateCachePolicyInput struct {
	Id      string
	IfMatch string
	Config  *cloudfrontstore.CachePolicyConfig
}

// ListPoliciesInput carries the pagination parameters and the optional
// managed|custom filter shared by the policy list operations.
type ListPoliciesInput struct {
	Marker     string
	MaxItems   int
	TypeFilter string
}

// CachePoliciesResult is the transport-agnostic list result for cache
// policies, already filtered by the requested policy type.
type CachePoliciesResult struct {
	Policies          []*cloudfrontstore.CachePolicy
	EffectiveMaxItems int
	NextMarker        string
}

// OriginRequestPoliciesResult is the transport-agnostic list result for
// origin request policies, already filtered by the requested policy type.
type OriginRequestPoliciesResult struct {
	Policies          []*cloudfrontstore.OriginRequestPolicy
	EffectiveMaxItems int
	NextMarker        string
}

func (s *CloudFrontService) createCachePolicyCore(stores *cloudfrontStores, in CreateCachePolicyInput) (*cloudfrontstore.CachePolicy, error) {
	if in.Config == nil || in.Config.Name == "" {
		return nil, invalidArgument("Name is required")
	}
	existing, _ := stores.cachePolicies.GetByName(in.Config.Name)
	if existing != nil {
		return nil, awserrors.NewAWSError("CachePolicyAlreadyExists", "Cache policy with this name already exists", 409)
	}
	return stores.cachePolicies.Create(in.Config.Name, "", in.Config)
}

func (s *CloudFrontService) getCachePolicyCore(stores *cloudfrontStores, id string) (*cloudfrontstore.CachePolicy, error) {
	if err := requireID(id); err != nil {
		return nil, err
	}
	cp, err := stores.cachePolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchCachePolicy", "Cache policy not found", 404)
		}
		return nil, err
	}
	return cp, nil
}

func (s *CloudFrontService) updateCachePolicyCore(stores *cloudfrontStores, in UpdateCachePolicyInput) (*cloudfrontstore.CachePolicy, error) {
	if err := requireID(in.Id); err != nil {
		return nil, err
	}
	if in.Config == nil || in.Config.Name == "" {
		return nil, invalidArgument("Name is required")
	}
	existing, err := stores.cachePolicies.Get(in.Id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchCachePolicy", "Cache policy not found", 404)
		}
		return nil, err
	}
	if err := verifyIfMatch(in.IfMatch, existing.ETag); err != nil {
		return nil, err
	}
	if err := ensureNameAvailable(existing.Name, in.Config.Name,
		func(name string) bool { dup, _ := stores.cachePolicies.GetByName(name); return dup != nil },
		awserrors.NewAWSError("CachePolicyAlreadyExists", "Cache policy with this name already exists", 409)); err != nil {
		return nil, err
	}
	return stores.cachePolicies.Update(in.Id, in.Config)
}

func (s *CloudFrontService) deleteCachePolicyCore(stores *cloudfrontStores, id, ifMatch string) error {
	if err := requireID(id); err != nil {
		return err
	}
	existing, err := stores.cachePolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchCachePolicy", "Cache policy not found", 404)
		}
		return err
	}
	if err := verifyIfMatch(ifMatch, existing.ETag); err != nil {
		return err
	}
	attached, err := isCachePolicyAttached(stores, id)
	if err != nil {
		return awserrors.NewAWSError("InternalError", "Failed to check cache policy usage: "+err.Error(), 500)
	}
	if attached {
		return awserrors.NewAWSError("CachePolicyInUse",
			"Cannot delete this cache policy because it is attached to one or more distributions", 409)
	}
	if err := stores.cachePolicies.Delete(id); err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchCachePolicy", "Cache policy not found", 404)
		}
		return err
	}
	return nil
}

func (s *CloudFrontService) listCachePoliciesCore(stores *cloudfrontStores, in ListPoliciesInput) (*CachePoliciesResult, error) {
	if in.TypeFilter != "" && !isValidPolicyListType(in.TypeFilter) {
		return nil, invalidArgument("Invalid Type: " + in.TypeFilter)
	}
	maxItems := resolveListMaxItems(in.MaxItems)
	result, err := stores.cachePolicies.List(in.Marker, maxItems)
	if err != nil {
		return nil, err
	}
	policies := make([]*cloudfrontstore.CachePolicy, 0, len(result.CachePolicies))
	for _, cp := range result.CachePolicies {
		if policyMatchesListType(cp.IsManaged, in.TypeFilter) {
			policies = append(policies, cp)
		}
	}
	return &CachePoliciesResult{
		Policies:          policies,
		EffectiveMaxItems: maxItems,
		NextMarker:        result.NextMarker,
	}, nil
}

// ---------------------------------------------------------------------------
// Origin request policy cores
// ---------------------------------------------------------------------------

// CreateOriginRequestPolicyInput carries the parsed origin request policy
// configuration.
type CreateOriginRequestPolicyInput struct {
	Config *cloudfrontstore.OriginRequestPolicyConfig
}

// UpdateOriginRequestPolicyInput carries the parameters for updating an
// origin request policy.
type UpdateOriginRequestPolicyInput struct {
	Id      string
	IfMatch string
	Config  *cloudfrontstore.OriginRequestPolicyConfig
}

func (s *CloudFrontService) createOriginRequestPolicyCore(stores *cloudfrontStores, in CreateOriginRequestPolicyInput) (*cloudfrontstore.OriginRequestPolicy, error) {
	if in.Config == nil || in.Config.Name == "" {
		return nil, invalidArgument("Name is required")
	}
	existing, _ := stores.originRequestPolicies.GetByName(in.Config.Name)
	if existing != nil {
		return nil, awserrors.NewAWSError("OriginRequestPolicyAlreadyExists", "Origin request policy with this name already exists", 409)
	}
	return stores.originRequestPolicies.Create(in.Config.Name, "", in.Config)
}

func (s *CloudFrontService) getOriginRequestPolicyCore(stores *cloudfrontStores, id string) (*cloudfrontstore.OriginRequestPolicy, error) {
	if err := requireID(id); err != nil {
		return nil, err
	}
	p, err := stores.originRequestPolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchOriginRequestPolicy", "Origin request policy not found", 404)
		}
		return nil, err
	}
	return p, nil
}

func (s *CloudFrontService) updateOriginRequestPolicyCore(stores *cloudfrontStores, in UpdateOriginRequestPolicyInput) (*cloudfrontstore.OriginRequestPolicy, error) {
	if err := requireID(in.Id); err != nil {
		return nil, err
	}
	if in.Config == nil || in.Config.Name == "" {
		return nil, invalidArgument("Name is required")
	}
	existing, err := stores.originRequestPolicies.Get(in.Id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchOriginRequestPolicy", "Origin request policy not found", 404)
		}
		return nil, err
	}
	if err := verifyIfMatch(in.IfMatch, existing.ETag); err != nil {
		return nil, err
	}
	if err := ensureNameAvailable(existing.Name, in.Config.Name,
		func(name string) bool { dup, _ := stores.originRequestPolicies.GetByName(name); return dup != nil },
		awserrors.NewAWSError("OriginRequestPolicyAlreadyExists", "Origin request policy with this name already exists", 409)); err != nil {
		return nil, err
	}
	return stores.originRequestPolicies.Update(in.Id, in.Config)
}

func (s *CloudFrontService) deleteOriginRequestPolicyCore(stores *cloudfrontStores, id, ifMatch string) error {
	if err := requireID(id); err != nil {
		return err
	}
	existing, err := stores.originRequestPolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchOriginRequestPolicy", "Origin request policy not found", 404)
		}
		return err
	}
	if err := verifyIfMatch(ifMatch, existing.ETag); err != nil {
		return err
	}
	attached, err := isOriginRequestPolicyAttached(stores, id)
	if err != nil {
		return awserrors.NewAWSError("InternalError", "Failed to check origin request policy usage: "+err.Error(), 500)
	}
	if attached {
		return awserrors.NewAWSError("OriginRequestPolicyInUse",
			"Cannot delete this origin request policy because it is attached to one or more distributions", 409)
	}
	if err := stores.originRequestPolicies.Delete(id); err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchOriginRequestPolicy", "Origin request policy not found", 404)
		}
		return err
	}
	return nil
}

func (s *CloudFrontService) listOriginRequestPoliciesCore(stores *cloudfrontStores, in ListPoliciesInput) (*OriginRequestPoliciesResult, error) {
	if in.TypeFilter != "" && !isValidPolicyListType(in.TypeFilter) {
		return nil, invalidArgument("Invalid Type: " + in.TypeFilter)
	}
	maxItems := resolveListMaxItems(in.MaxItems)
	result, err := stores.originRequestPolicies.List(in.Marker, maxItems)
	if err != nil {
		return nil, err
	}
	policies := make([]*cloudfrontstore.OriginRequestPolicy, 0, len(result.OriginRequestPolicies))
	for _, p := range result.OriginRequestPolicies {
		if policyMatchesListType(p.IsManaged, in.TypeFilter) {
			policies = append(policies, p)
		}
	}
	return &OriginRequestPoliciesResult{
		Policies:          policies,
		EffectiveMaxItems: maxItems,
		NextMarker:        result.NextMarker,
	}, nil
}

// ---------------------------------------------------------------------------
// Tag cores — single validation and persistence path for the CloudFront
// resource tagging operations.
// ---------------------------------------------------------------------------

// ListTagsForResourceInput carries the resource ARN for
// ListTagsForResource.
type ListTagsForResourceInput struct {
	Resource string
}

// TagResourceInput carries the resource ARN and the parsed tags for
// TagResource.
type TagResourceInput struct {
	Resource string
	Tags     []types.Tag
}

// UntagResourceInput carries the resource ARN and the parsed tag keys for
// UntagResource.
type UntagResourceInput struct {
	Resource string
	TagKeys  []string
}

// validateTagResourceArn enforces the shared resource-ARN contract of the
// three tagging operations: the parameter is required and must match the
// CloudFront ARN pattern.
func validateTagResourceArn(arn string) error {
	if arn == "" {
		return awserrors.NewAWSError("InvalidArgument", "Resource is required", 400)
	}
	if !isValidResourceArn(arn) {
		return awserrors.NewAWSError("InvalidArgument", "Resource must be a CloudFront resource ARN: "+arn, 400)
	}
	return nil
}

// listTagsForResourceCore is the single entry point for reading the tags
// of a CloudFront resource.
func (s *CloudFrontService) listTagsForResourceCore(stores *cloudfrontStores, in ListTagsForResourceInput) ([]types.Tag, error) {
	if err := validateTagResourceArn(in.Resource); err != nil {
		return nil, err
	}

	tags, err := stores.tags.ListTagsForResource(in.Resource)
	if err != nil {
		return nil, awserrors.NewAWSError("InternalError", err.Error(), 500)
	}
	return tags, nil
}

// tagResourceCore is the single entry point for applying tags to a
// CloudFront resource: it validates the resource ARN, the tag set, and
// every tag key and value before persisting.
func (s *CloudFrontService) tagResourceCore(stores *cloudfrontStores, in TagResourceInput) error {
	if err := validateTagResourceArn(in.Resource); err != nil {
		return err
	}

	if len(in.Tags) == 0 {
		return awserrors.NewAWSError("InvalidArgument", "At least one tag is required", 400)
	}
	for _, t := range in.Tags {
		if !isValidTagKey(t.Key) {
			return awserrors.NewAWSError("InvalidArgument", fmt.Sprintf("Invalid tag key: %q", t.Key), 400)
		}
		if !isValidTagValue(t.Value) {
			return awserrors.NewAWSError("InvalidArgument", fmt.Sprintf("Invalid tag value for key %q", t.Key), 400)
		}
	}

	if err := stores.tags.Tag(in.Resource, in.Tags); err != nil {
		return awserrors.NewAWSError("InternalError", err.Error(), 500)
	}
	return nil
}

// untagResourceCore is the single entry point for removing tags from a
// CloudFront resource: it validates the resource ARN and every tag key
// before persisting.
func (s *CloudFrontService) untagResourceCore(stores *cloudfrontStores, in UntagResourceInput) error {
	if err := validateTagResourceArn(in.Resource); err != nil {
		return err
	}

	if len(in.TagKeys) == 0 {
		return awserrors.NewAWSError("InvalidArgument", "At least one tag key is required", 400)
	}
	for _, k := range in.TagKeys {
		if !isValidTagKey(k) {
			return awserrors.NewAWSError("InvalidArgument", fmt.Sprintf("Invalid tag key: %q", k), 400)
		}
	}

	if err := stores.tags.Untag(in.Resource, in.TagKeys); err != nil {
		return awserrors.NewAWSError("InternalError", err.Error(), 500)
	}
	return nil
}
