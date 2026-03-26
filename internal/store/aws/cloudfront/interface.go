package cloudfront

import "vorpalstacks/internal/store/aws/common"

// DistributionStoreInterface defines operations for managing CloudFront distributions.
type DistributionStoreInterface interface {
	Get(id string) (*Distribution, error)
	GetByCallerReference(callerRef string) (*Distribution, error)
	Raw() *DistributionStore
}

// OriginAccessControlStoreInterface defines operations for managing origin access controls.
type OriginAccessControlStoreInterface interface {
	Get(id string) (*OriginAccessControl, error)
	Raw() *OriginAccessControlStore
}

// CachePolicyStoreInterface defines operations for managing cache policies.
type CachePolicyStoreInterface interface {
	Get(id string) (*CachePolicy, error)
	Raw() *CachePolicyStore
}

// OriginRequestPolicyStoreInterface defines operations for managing origin request policies.
type OriginRequestPolicyStoreInterface interface {
	Get(id string) (*OriginRequestPolicy, error)
	Raw() *OriginRequestPolicyStore
}

// ResponseHeadersPolicyStoreInterface defines operations for managing response headers policies.
type ResponseHeadersPolicyStoreInterface interface {
	Get(id string) (*ResponseHeadersPolicy, error)
	Raw() *ResponseHeadersPolicyStore
}

// TagStoreInterface defines operations for managing CloudFront tags.
type TagStoreInterface interface {
	ListTagsForResource(arn string) ([]common.Tag, error)
	TagResource(arn string, tags []common.Tag) error
	UntagResource(arn string, tagKeys []string) error
	Raw() *TagStore
}

// CloudFrontStoresInterface defines access to all CloudFront stores.
type CloudFrontStoresInterface interface {
	Distributions() DistributionStoreInterface
	OriginAccessControls() OriginAccessControlStoreInterface
	CachePolicies() CachePolicyStoreInterface
	OriginRequestPolicies() OriginRequestPolicyStoreInterface
	ResponseHeadersPolicies() ResponseHeadersPolicyStoreInterface
	Tags() TagStoreInterface
	ARNBuilder() *ARNBuilder
	Raw() *CloudFrontStores
}

// CloudFrontStores provides access to all CloudFront stores.
type CloudFrontStores struct {
	distributions           *DistributionStore
	cachePolicies           *CachePolicyStore
	originRequestPolicies   *OriginRequestPolicyStore
	originAccessControls    *OriginAccessControlStore
	responseHeadersPolicies *ResponseHeadersPolicyStore
	tags                    *TagStore
	arnBuilder              *ARNBuilder
}

// NewCloudFrontStores creates a new CloudFrontStores with the given stores.
func NewCloudFrontStores(
	distributions *DistributionStore,
	cachePolicies *CachePolicyStore,
	originRequestPolicies *OriginRequestPolicyStore,
	originAccessControls *OriginAccessControlStore,
	responseHeadersPolicies *ResponseHeadersPolicyStore,
	tags *TagStore,
	arnBuilder *ARNBuilder,
) *CloudFrontStores {
	return &CloudFrontStores{
		distributions:           distributions,
		cachePolicies:           cachePolicies,
		originRequestPolicies:   originRequestPolicies,
		originAccessControls:    originAccessControls,
		responseHeadersPolicies: responseHeadersPolicies,
		tags:                    tags,
		arnBuilder:              arnBuilder,
	}
}

// Distributions returns the distribution store.
func (s *CloudFrontStores) Distributions() DistributionStoreInterface {
	return s.distributions
}

// OriginAccessControls returns the origin access control store.
func (s *CloudFrontStores) OriginAccessControls() OriginAccessControlStoreInterface {
	return s.originAccessControls
}

// CachePolicies returns the cache policy store.
func (s *CloudFrontStores) CachePolicies() CachePolicyStoreInterface {
	return s.cachePolicies
}

// OriginRequestPolicies returns the origin request policy store.
func (s *CloudFrontStores) OriginRequestPolicies() OriginRequestPolicyStoreInterface {
	return s.originRequestPolicies
}

// ResponseHeadersPolicies returns the response headers policy store.
func (s *CloudFrontStores) ResponseHeadersPolicies() ResponseHeadersPolicyStoreInterface {
	return s.responseHeadersPolicies
}

// Tags returns the tag store.
func (s *CloudFrontStores) Tags() TagStoreInterface {
	return s.tags
}

// ARNBuilder returns the ARN builder.
func (s *CloudFrontStores) ARNBuilder() *ARNBuilder {
	return s.arnBuilder
}

// Raw returns the CloudFront stores.
func (s *CloudFrontStores) Raw() *CloudFrontStores {
	return s
}

var _ DistributionStoreInterface = (*DistributionStore)(nil)
var _ OriginAccessControlStoreInterface = (*OriginAccessControlStore)(nil)
var _ CachePolicyStoreInterface = (*CachePolicyStore)(nil)
var _ OriginRequestPolicyStoreInterface = (*OriginRequestPolicyStore)(nil)
var _ ResponseHeadersPolicyStoreInterface = (*ResponseHeadersPolicyStore)(nil)
var _ TagStoreInterface = (*TagStore)(nil)
var _ CloudFrontStoresInterface = (*CloudFrontStores)(nil)
