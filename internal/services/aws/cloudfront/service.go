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
	invalidations           *cloudfrontstore.InvalidationStore
	arnBuilder              *cloudfrontstore.ARNBuilder
}

// CloudFrontService provides AWS CloudFront operations.
type CloudFrontService struct {
	accountID           string
	storageManager      *storage.RegionStorageManager
	stores              sync.Map // global (no region) — single cached instance
	seedManagedPolicies sync.Once
	distributionServer  *DistributionServer
	wafInvoker          eventbus.WAFInvoker
	acmInvoker          eventbus.ACMInvoker
}

// NewCloudFrontService creates a new CloudFront service instance.
func NewCloudFrontService(accountID string) *CloudFrontService {
	return &CloudFrontService{
		accountID: accountID,
	}
}

// SetStorageManager injects the storage manager needed for lazy store
// creation and the distribution server. CloudFront is a global service,
// so no region parameter is required.
func (s *CloudFrontService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// InitDistributionServer creates the DistributionServer owned by this service.
// CloudFront is a global service, so the distribution server reads from the
// global Pebble DB. After the management API creates the global store
// instance, SetDistributionStore is called to share it with the server,
// eliminating any second store backed by a different DB.
func (s *CloudFrontService) InitDistributionServer() *DistributionServer {
	s.distributionServer = NewDistributionServer(s.storageManager, s.accountID)
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

// createStores builds a complete cloudfrontStores from the given global
// storage. Called by both store() and GetStoreForRegion() to ensure a
// single code path for store construction.
func (s *CloudFrontService) createStores(st storage.BasicStorage) *cloudfrontStores {
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
		publicKeys:              cloudfrontstore.NewPublicKeyStore(st, s.accountID),
		keyGroups:               cloudfrontstore.NewKeyGroupStore(st, s.accountID),
		tags:                    cloudfrontstore.NewTagStore(st),
		invalidations:           cloudfrontstore.NewInvalidationStore(st),
		arnBuilder:              arnBuilder,
	}
	if s.distributionServer != nil {
		s.distributionServer.SetDistributionStore(stores.distributions)
	}
	return stores
}

func (s *CloudFrontService) store(reqCtx *request.RequestContext) (*cloudfrontStores, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, "global", func() (*cloudfrontStores, error) {
		st, err := reqCtx.GetGlobalStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get global storage: %w", err)
		}
		return s.createStores(st), nil
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
	stores := s.createStores(st)
	actual, _ := s.stores.LoadOrStore("global", stores)
	return actual.(*cloudfrontStores), nil
}

// SetWAFInvoker injects the WAF invoker for cross-service WebACL association.
func (s *CloudFrontService) SetWAFInvoker(invoker eventbus.WAFInvoker) {
	s.wafInvoker = invoker
}

// SetACMInvoker injects the ACM invoker for cross-service certificate usage
// tracking. When a distribution references an ACM certificate, the invoker
// records the association so that DeleteCertificate can enforce InUseBy.
func (s *CloudFrontService) SetACMInvoker(invoker eventbus.ACMInvoker) {
	s.acmInvoker = invoker
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
	d.RegisterHandlerForService("cloudfront", "ListDistributionsByCachePolicyId", s.ListDistributionsByCachePolicyId)
	d.RegisterHandlerForService("cloudfront", "ListDistributionsByOriginRequestPolicyId", s.ListDistributionsByOriginRequestPolicyId)
	d.RegisterHandlerForService("cloudfront", "ListDistributionsByResponseHeadersPolicyId", s.ListDistributionsByResponseHeadersPolicyId)
	d.RegisterHandlerForService("cloudfront", "ListDistributionsByKeyGroup", s.ListDistributionsByKeyGroup)
	d.RegisterHandlerForService("cloudfront", "UpdateDistribution", s.UpdateDistribution)
	d.RegisterHandlerForService("cloudfront", "DeleteDistribution", s.DeleteDistribution)
	d.RegisterHandlerForService("cloudfront", "AssociateDistributionWebACL", s.AssociateDistributionWebACL)
	d.RegisterHandlerForService("cloudfront", "DisassociateDistributionWebACL", s.DisassociateDistributionWebACL)
	d.RegisterHandlerForService("cloudfront", "CopyDistribution", s.CopyDistribution)

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
