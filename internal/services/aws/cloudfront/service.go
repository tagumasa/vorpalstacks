// Package cloudfront provides AWS CloudFront service operations for vorpalstacks.
package cloudfront

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

var managedPoliciesOnce sync.Once

// cloudfrontStores holds the various CloudFront stores.
type cloudfrontStores struct {
	distributions           *cloudfrontstore.DistributionStore
	cachePolicies           *cloudfrontstore.CachePolicyStore
	originRequestPolicies   *cloudfrontstore.OriginRequestPolicyStore
	originAccessControls    *cloudfrontstore.OriginAccessControlStore
	responseHeadersPolicies *cloudfrontstore.ResponseHeadersPolicyStore
	tags                    *cloudfrontstore.TagStore
	arnBuilder              *cloudfrontstore.ARNBuilder
}

// CloudFrontService provides AWS CloudFront operations.
type CloudFrontService struct {
	accountID string
}

// NewCloudFrontService creates a new CloudFront service instance.
func NewCloudFrontService(accountID string) *CloudFrontService {
	return &CloudFrontService{
		accountID: accountID,
	}
}

// store returns the CloudFront stores for the given request context.
func (s *CloudFrontService) store(reqCtx *request.RequestContext) (*cloudfrontStores, error) {
	if stores := reqCtx.GetCloudFrontStores(); stores != nil {
		raw := stores.Raw()
		return &cloudfrontStores{
			distributions:           raw.Distributions().Raw(),
			cachePolicies:           raw.CachePolicies().Raw(),
			originRequestPolicies:   raw.OriginRequestPolicies().Raw(),
			originAccessControls:    raw.OriginAccessControls().Raw(),
			responseHeadersPolicies: raw.ResponseHeadersPolicies().Raw(),
			tags:                    raw.Tags().Raw(),
			arnBuilder:              raw.ARNBuilder(),
		}, nil
	}
	storage, err := reqCtx.GetGlobalStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to get global storage: %w", err)
	}
	arnBuilder := cloudfrontstore.NewARNBuilder(s.accountID)
	cacheStore := cloudfrontstore.NewCachePolicyStore(storage, s.accountID)
	orpStore := cloudfrontstore.NewOriginRequestPolicyStore(storage, s.accountID)
	managedPoliciesOnce.Do(func() {
		cloudfrontstore.SeedManagedPolicies(cacheStore, orpStore)
	})
	return &cloudfrontStores{
		distributions:           cloudfrontstore.NewDistributionStore(storage, s.accountID),
		cachePolicies:           cacheStore,
		originRequestPolicies:   orpStore,
		originAccessControls:    cloudfrontstore.NewOriginAccessControlStore(storage, s.accountID),
		responseHeadersPolicies: cloudfrontstore.NewResponseHeadersPolicyStore(storage, s.accountID),
		tags:                    cloudfrontstore.NewTagStore(storage),
		arnBuilder:              arnBuilder,
	}, nil
}

func (s *CloudFrontService) wafAssociationStore(reqCtx *request.RequestContext) (*wafstore.WebACLAssociationStore, error) {
	storage, err := reqCtx.GetGlobalStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to get global storage for WAF association: %w", err)
	}
	return wafstore.NewWebACLAssociationStore(storage), nil
}

func (s *CloudFrontService) syncWAFAssociation(reqCtx *request.RequestContext, webACLId, distributionArn string) error {
	if webACLId == "" {
		return nil
	}
	assocStore, err := s.wafAssociationStore(reqCtx)
	if err != nil {
		return err
	}
	return assocStore.Associate(webACLId, distributionArn)
}

func (s *CloudFrontService) removeWAFAssociation(reqCtx *request.RequestContext, webACLId, distributionArn string) error {
	if webACLId == "" {
		return nil
	}
	assocStore, err := s.wafAssociationStore(reqCtx)
	if err != nil {
		return err
	}
	return assocStore.Disassociate(webACLId, distributionArn)
}

// AccountId returns the account ID.
func (s *CloudFrontService) AccountId() string {
	return s.accountID
}

// generateETag generates a random ETag value.
func (s *CloudFrontService) generateETag() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		bytes = make([]byte, 16)
		for i := range bytes {
			bytes[i] = byte(i % 256)
		}
	}
	return hex.EncodeToString(bytes)
}

// RegisterHandlers registers the CloudFront handlers with the dispatcher.
func (s *CloudFrontService) RegisterHandlers(d *dispatcher.Dispatcher) {
	d.RegisterHandlerForService("cloudfront", "CreateDistribution", s.CreateDistribution)
	d.RegisterHandlerForService("cloudfront", "CreateDistributionWithTags", s.CreateDistributionWithTags)
	d.RegisterHandlerForService("cloudfront", "GetDistribution", s.GetDistribution)
	d.RegisterHandlerForService("cloudfront", "GetDistributionConfig", s.GetDistributionConfig)
	d.RegisterHandlerForService("cloudfront", "ListDistributions", s.ListDistributions)
	d.RegisterHandlerForService("cloudfront", "UpdateDistribution", s.UpdateDistribution)
	d.RegisterHandlerForService("cloudfront", "DeleteDistribution", s.DeleteDistribution)

	d.RegisterHandlerForService("cloudfront", "CreateInvalidation", s.CreateInvalidation)
	d.RegisterHandlerForService("cloudfront", "GetInvalidation", s.GetInvalidation)
	d.RegisterHandlerForService("cloudfront", "ListInvalidations", s.ListInvalidations)

	d.RegisterHandlerForService("cloudfront", "CreateCachePolicy", s.CreateCachePolicy)
	d.RegisterHandlerForService("cloudfront", "GetCachePolicy", s.GetCachePolicy)
	d.RegisterHandlerForService("cloudfront", "ListCachePolicies", s.ListCachePolicies)
	d.RegisterHandlerForService("cloudfront", "DeleteCachePolicy", s.DeleteCachePolicy)

	d.RegisterHandlerForService("cloudfront", "CreateOriginRequestPolicy", s.CreateOriginRequestPolicy)
	d.RegisterHandlerForService("cloudfront", "GetOriginRequestPolicy", s.GetOriginRequestPolicy)
	d.RegisterHandlerForService("cloudfront", "ListOriginRequestPolicies", s.ListOriginRequestPolicies)
	d.RegisterHandlerForService("cloudfront", "DeleteOriginRequestPolicy", s.DeleteOriginRequestPolicy)

	d.RegisterHandlerForService("cloudfront", "CreateOriginAccessControl", s.CreateOriginAccessControl)
	d.RegisterHandlerForService("cloudfront", "GetOriginAccessControl", s.GetOriginAccessControl)
	d.RegisterHandlerForService("cloudfront", "UpdateOriginAccessControl", s.UpdateOriginAccessControl)
	d.RegisterHandlerForService("cloudfront", "DeleteOriginAccessControl", s.DeleteOriginAccessControl)
	d.RegisterHandlerForService("cloudfront", "ListOriginAccessControls", s.ListOriginAccessControls)

	d.RegisterHandlerForService("cloudfront", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("cloudfront", "TagResource", s.TagResource)
	d.RegisterHandlerForService("cloudfront", "UntagResource", s.UntagResource)

	d.RegisterHandlerForService("cloudfront", "CreateResponseHeadersPolicy", s.CreateResponseHeadersPolicy)
	d.RegisterHandlerForService("cloudfront", "GetResponseHeadersPolicy", s.GetResponseHeadersPolicy)
	d.RegisterHandlerForService("cloudfront", "GetResponseHeadersPolicyConfig", s.GetResponseHeadersPolicyConfig)
	d.RegisterHandlerForService("cloudfront", "UpdateResponseHeadersPolicy", s.UpdateResponseHeadersPolicy)
	d.RegisterHandlerForService("cloudfront", "DeleteResponseHeadersPolicy", s.DeleteResponseHeadersPolicy)
	d.RegisterHandlerForService("cloudfront", "ListResponseHeadersPolicies", s.ListResponseHeadersPolicies)
}
