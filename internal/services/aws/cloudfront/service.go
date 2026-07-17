// Package cloudfront provides AWS CloudFront service operations for vorpalstacks.
package cloudfront

import (
	"fmt"
	"net/http"
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// cloudfrontStores holds the various CloudFront stores.
type cloudfrontStores struct {
	distributions           *cloudfrontstore.DistributionStore
	cachePolicies           *cloudfrontstore.CachePolicyStore
	originRequestPolicies   *cloudfrontstore.OriginRequestPolicyStore
	originAccessControls    *cloudfrontstore.OriginAccessControlStore
	responseHeadersPolicies *cloudfrontstore.ResponseHeadersPolicyStore
	publicKeys              *cloudfrontstore.PublicKeyStore
	keyGroups               *cloudfrontstore.KeyGroupStore
	tags                    *cloudfrontstore.TagStore
	arnBuilder              *cloudfrontstore.ARNBuilder
}

// CloudFrontService provides AWS CloudFront operations.
type CloudFrontService struct {
	accountID           string
	region              string
	storageManager      *storage.RegionStorageManager
	stores              sync.Map // global (no region) — single cached instance
	seedManagedPolicies sync.Once
	distributionServer  *DistributionServer
	wafInvoker          eventbus.WAFInvoker
}

// NewCloudFrontService creates a new CloudFront service instance.
func NewCloudFrontService(accountID string) *CloudFrontService {
	return &CloudFrontService{
		accountID: accountID,
	}
}

// SetRegionAndStorage injects the region and storage manager needed for
// creating the distribution server with a shared store.
func (s *CloudFrontService) SetRegionAndStorage(region string, sm *storage.RegionStorageManager) {
	s.region = region
	s.storageManager = sm
}

// InitDistributionServer creates the DistributionServer owned by this service.
// The distribution store is populated lazily on first management API call;
// until then the server uses its own fallback store backed by the same Pebble.
func (s *CloudFrontService) InitDistributionServer() *DistributionServer {
	s.distributionServer = NewDistributionServer(s.storageManager, s.accountID, s.region)
	return s.distributionServer
}

// DistributionHandler returns an http.Handler for the distribution server,
// or nil if it has not been initialised.
func (s *CloudFrontService) DistributionHandler() http.Handler {
	if s.distributionServer == nil {
		return nil
	}
	return http.HandlerFunc(s.distributionServer.HandleRequest)
}

func (s *CloudFrontService) store(reqCtx *request.RequestContext) (*cloudfrontStores, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, "global", func() (*cloudfrontStores, error) {
		storage, err := reqCtx.GetGlobalStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get global storage: %w", err)
		}
		arnBuilder := cloudfrontstore.NewARNBuilder(s.accountID)
		cacheStore := cloudfrontstore.NewCachePolicyStore(storage, s.accountID)
		orpStore := cloudfrontstore.NewOriginRequestPolicyStore(storage, s.accountID)
		s.seedManagedPolicies.Do(func() {
			cloudfrontstore.SeedManagedPolicies(cacheStore, orpStore)
		})
		return &cloudfrontStores{
			distributions:           cloudfrontstore.NewDistributionStore(storage, s.accountID),
			cachePolicies:           cacheStore,
			originRequestPolicies:   orpStore,
			originAccessControls:    cloudfrontstore.NewOriginAccessControlStore(storage, s.accountID),
			responseHeadersPolicies: cloudfrontstore.NewResponseHeadersPolicyStore(storage, s.accountID),
			publicKeys:              cloudfrontstore.NewPublicKeyStore(storage, s.accountID),
			keyGroups:               cloudfrontstore.NewKeyGroupStore(storage, s.accountID),
			tags:                    cloudfrontstore.NewTagStore(storage),
			arnBuilder:              arnBuilder,
		}, nil
	})
}

// GetStoreForRegion returns the cached CloudFront stores. CloudFront is a
// global service, so the region parameter is ignored and the single cached
// instance is always returned.
func (s *CloudFrontService) GetStoreForRegion(_ string) (*cloudfrontStores, error) {
	if v, ok := s.stores.Load("global"); ok {
		return v.(*cloudfrontStores), nil
	}
	if s.storageManager == nil {
		return nil, fmt.Errorf("cloudfront storage manager not initialised")
	}
	st, err := s.storageManager.GetGlobalStorage()
	if err != nil {
		return nil, err
	}
	arnBuilder := cloudfrontstore.NewARNBuilder(s.accountID)
	cacheStore := cloudfrontstore.NewCachePolicyStore(st, s.accountID)
	orpStore := cloudfrontstore.NewOriginRequestPolicyStore(st, s.accountID)
	s.seedManagedPolicies.Do(func() {
		cloudfrontstore.SeedManagedPolicies(cacheStore, orpStore)
	})
	stores := &cloudfrontStores{
		distributions:           cloudfrontstore.NewDistributionStore(st, s.accountID),
		cachePolicies:           cacheStore,
		originRequestPolicies:   orpStore,
		originAccessControls:    cloudfrontstore.NewOriginAccessControlStore(st, s.accountID),
		responseHeadersPolicies: cloudfrontstore.NewResponseHeadersPolicyStore(st, s.accountID),
		tags:                    cloudfrontstore.NewTagStore(st),
		arnBuilder:              arnBuilder,
	}
	actual, _ := s.stores.LoadOrStore("global", stores)
	return actual.(*cloudfrontStores), nil
}

// SetWAFInvoker injects the WAF invoker for cross-service WebACL association.
func (s *CloudFrontService) SetWAFInvoker(invoker eventbus.WAFInvoker) {
	s.wafInvoker = invoker
}

// AccountId returns the account ID.
func (s *CloudFrontService) AccountId() string {
	return s.accountID
}

// RegisterHandlers registers all CloudFront API handlers with the dispatcher.
func (s *CloudFrontService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("cloudfront", "CreateDistribution", s.CreateDistribution)
	d.RegisterHandlerForService("cloudfront", "CreateDistributionWithTags", s.CreateDistributionWithTags)
	d.RegisterHandlerForService("cloudfront", "GetDistribution", s.GetDistribution)
	d.RegisterHandlerForService("cloudfront", "GetDistributionConfig", s.GetDistributionConfig)
	d.RegisterHandlerForService("cloudfront", "ListDistributions", s.ListDistributions)
	d.RegisterHandlerForService("cloudfront", "ListDistributionsByWebACLId", s.ListDistributionsByWebACLId)
	d.RegisterHandlerForService("cloudfront", "UpdateDistribution", s.UpdateDistribution)
	d.RegisterHandlerForService("cloudfront", "DeleteDistribution", s.DeleteDistribution)

	d.RegisterHandlerForService("cloudfront", "CreateInvalidation", s.CreateInvalidation)
	d.RegisterHandlerForService("cloudfront", "GetInvalidation", s.GetInvalidation)
	d.RegisterHandlerForService("cloudfront", "ListInvalidations", s.ListInvalidations)

	d.RegisterHandlerForService("cloudfront", "CreateCachePolicy", s.CreateCachePolicy)
	d.RegisterHandlerForService("cloudfront", "GetCachePolicy", s.GetCachePolicy)
	d.RegisterHandlerForService("cloudfront", "GetCachePolicyConfig", s.GetCachePolicyConfig)
	d.RegisterHandlerForService("cloudfront", "ListCachePolicies", s.ListCachePolicies)
	d.RegisterHandlerForService("cloudfront", "UpdateCachePolicy", s.UpdateCachePolicy)
	d.RegisterHandlerForService("cloudfront", "DeleteCachePolicy", s.DeleteCachePolicy)

	d.RegisterHandlerForService("cloudfront", "CreateOriginRequestPolicy", s.CreateOriginRequestPolicy)
	d.RegisterHandlerForService("cloudfront", "GetOriginRequestPolicy", s.GetOriginRequestPolicy)
	d.RegisterHandlerForService("cloudfront", "GetOriginRequestPolicyConfig", s.GetOriginRequestPolicyConfig)
	d.RegisterHandlerForService("cloudfront", "ListOriginRequestPolicies", s.ListOriginRequestPolicies)
	d.RegisterHandlerForService("cloudfront", "UpdateOriginRequestPolicy", s.UpdateOriginRequestPolicy)
	d.RegisterHandlerForService("cloudfront", "DeleteOriginRequestPolicy", s.DeleteOriginRequestPolicy)

	d.RegisterHandlerForService("cloudfront", "CreateOriginAccessControl", s.CreateOriginAccessControl)
	d.RegisterHandlerForService("cloudfront", "GetOriginAccessControl", s.GetOriginAccessControl)
	d.RegisterHandlerForService("cloudfront", "GetOriginAccessControlConfig", s.GetOriginAccessControlConfig)
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

	d.RegisterHandlerForService("cloudfront", "CreatePublicKey", s.CreatePublicKey)
	d.RegisterHandlerForService("cloudfront", "GetPublicKey", s.GetPublicKey)
	d.RegisterHandlerForService("cloudfront", "GetPublicKeyConfig", s.GetPublicKeyConfig)
	d.RegisterHandlerForService("cloudfront", "UpdatePublicKey", s.UpdatePublicKey)
	d.RegisterHandlerForService("cloudfront", "DeletePublicKey", s.DeletePublicKey)
	d.RegisterHandlerForService("cloudfront", "ListPublicKeys", s.ListPublicKeys)

	d.RegisterHandlerForService("cloudfront", "CreateKeyGroup", s.CreateKeyGroup)
	d.RegisterHandlerForService("cloudfront", "GetKeyGroup", s.GetKeyGroup)
	d.RegisterHandlerForService("cloudfront", "GetKeyGroupConfig", s.GetKeyGroupConfig)
	d.RegisterHandlerForService("cloudfront", "UpdateKeyGroup", s.UpdateKeyGroup)
	d.RegisterHandlerForService("cloudfront", "DeleteKeyGroup", s.DeleteKeyGroup)
	d.RegisterHandlerForService("cloudfront", "ListKeyGroups", s.ListKeyGroups)
}
