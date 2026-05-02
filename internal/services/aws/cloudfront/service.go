// Package cloudfront provides AWS CloudFront service operations for vorpalstacks.
package cloudfront

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
	storecommon "vorpalstacks/internal/store/aws/common"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

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
	accountID           string
	region              string
	storageManager      *storage.RegionStorageManager
	stores              sync.Map // global (no region) — single cached instance
	seedManagedPolicies sync.Once
	distributionServer  *DistributionServer
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
			tags:                    cloudfrontstore.NewTagStore(storage),
			arnBuilder:              arnBuilder,
		}, nil
	})
}

var globalWAFAssocKey struct{}

func (s *CloudFrontService) wafAssociationStore(reqCtx *request.RequestContext) (*wafstore.WebACLAssociationStore, error) {
	if cached, ok := s.stores.Load(globalWAFAssocKey); ok {
		if typed, ok := cached.(*wafstore.WebACLAssociationStore); ok {
			return typed, nil
		}
	}
	storage, err := reqCtx.GetGlobalStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to get global storage for WAF association: %w", err)
	}
	store := wafstore.NewWebACLAssociationStore(storage)
	if actual, loaded := s.stores.LoadOrStore(globalWAFAssocKey, store); loaded {
		if typed, ok := actual.(*wafstore.WebACLAssociationStore); ok {
			return typed, nil
		}
	}
	return store, nil
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
func (s *CloudFrontService) ListKeyGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxItems := int(request.GetIntParam(req.Parameters, "MaxItems"))
	if maxItems == 0 {
		maxItems = 100
	}
	return map[string]interface{}{
		"KeyGroupList": map[string]interface{}{
			"MaxItems":    maxItems,
			"Quantity":    0,
			"IsTruncated": false,
		},
	}, nil
}

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

	d.RegisterHandlerForService("cloudfront", "ListKeyGroups", s.ListKeyGroups)
}
