// Package wafv2 provides AWS WAFv2 service operations for vorpalstacks.
// This package is part of the WAF service family and shares the underlying
// store with "vorpalstacks/internal/store/aws/waf" (IAM sub-service exemption
// declared for WAF WebACL association operations).
package wafv2

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"vorpalstacks/internal/common/defaults"

	"github.com/google/uuid"
	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/handler"
	waf "vorpalstacks/internal/common/invokers/waf"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	storecommon "vorpalstacks/internal/store/aws/common"
	wafstore "vorpalstacks/internal/store/aws/waf"
	"vorpalstacks/internal/utils/aws/arn"
)

var wafv2GlobalAssocKey struct{}

type wafv2Stores struct {
	webACLs          *wafstore.WebACLStore
	ruleGroups       *wafstore.RuleGroupStore
	ipSets           *wafstore.IPSetStore
	regexPatternSets *wafstore.RegexPatternSetStore
	associations     *wafstore.WebACLAssociationStore
	loggingConfigs   *wafstore.LoggingStore
	samples          *wafstore.SamplingStore
	tags             *storecommon.TagStore
	arnBuilder       *wafstore.ARNBuilder
}

// WAFv2Service implements the AWS WAF v2 API operations.
type WAFv2Service struct {
	accountID        string
	region           string
	stores           sync.Map // region → *wafv2Stores
	storageManager   *storage.RegionStorageManager
	resourceCheckers []waf.WebACLResourceChecker
	sampleSweepOnce  sync.Once
	challenges       *challengeRegistry
}

// NewWAFv2Service creates a new WAFv2Service instance.
func NewWAFv2Service(accountID, region string) *WAFv2Service {
	return &WAFv2Service{accountID: accountID, region: region, challenges: newChallengeRegistry()}
}

// SetStorageManager injects the region storage manager for lazy store creation.
func (s *WAFv2Service) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// RegisterWebACLResourceChecker registers a service that hosts resources
// eligible for WebACL association. AssociateWebACL rejects ARNs owned by a
// registered checker that cannot be resolved to an existing resource.
func (s *WAFv2Service) RegisterWebACLResourceChecker(checker waf.WebACLResourceChecker) {
	s.resourceCheckers = append(s.resourceCheckers, checker)
}

// ensureAssociableResource verifies that the ResourceArn resolves to an
// existing resource when its ARN service namespace is owned by a registered
// checker. AWS returns WAFUnavailableEntityException when it cannot retrieve
// the specified resource; namespaces without a checker are resource types
// this platform does not host, which keep stub-association semantics.
func (s *WAFv2Service) ensureAssociableResource(ctx context.Context, region, resourceArn string) error {
	parsed, err := arn.ParseARN(resourceArn)
	if err != nil {
		return invalidParamError("ResourceArn is not a valid ARN")
	}
	for _, checker := range s.resourceCheckers {
		if checker.WebACLResourceService() != parsed.Service {
			continue
		}
		if !checker.WebACLResourceExists(ctx, region, resourceArn) {
			return newAPIError("WAFUnavailableEntityException",
				fmt.Sprintf("WAF couldn't retrieve the resource %s. If you've just created the resource, wait a few minutes for the change to propagate and retry the operation.", resourceArn),
				http.StatusBadRequest)
		}
		return nil
	}
	return nil
}

// GetStoresForRegion returns the full wafv2Stores for the given region,
// creating a new set if not already cached.
func (s *WAFv2Service) GetStoresForRegion(region string) (*wafv2Stores, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, region, func() (*wafv2Stores, error) {
		if s.storageManager == nil {
			return nil, fmt.Errorf("wafv2 storage manager not initialised")
		}
		st, err := s.storageManager.GetStorage(region)
		if err != nil {
			return nil, fmt.Errorf("failed to get storage for region %s: %w", region, err)
		}
		return &wafv2Stores{
			webACLs:          wafstore.NewWebACLStore(st, s.accountID, region),
			ruleGroups:       wafstore.NewRuleGroupStore(st, s.accountID, region),
			ipSets:           wafstore.NewIPSetStore(st, s.accountID, region),
			regexPatternSets: wafstore.NewRegexPatternSetStore(st, s.accountID, region),
			associations:     wafstore.NewWebACLAssociationStore(st),
			loggingConfigs:   wafstore.NewLoggingStore(st),
			samples:          wafstore.NewSamplingStore(st),
			tags:             storecommon.NewTagStoreWithRegion(st, "wafv2", region),
			arnBuilder:       wafstore.NewARNBuilder(s.accountID, region),
		}, nil
	})
}

// GetWebACLStoreForRegion returns the cached WebACLStore for the given region,
// creating a new store group if not already cached.
func (s *WAFv2Service) GetWebACLStoreForRegion(region string) (*wafstore.WebACLStore, error) {
	if v, ok := s.stores.Load(region); ok {
		return v.(*wafv2Stores).webACLs, nil
	}
	if s.storageManager == nil {
		return nil, fmt.Errorf("wafv2 storage manager not initialised")
	}
	st, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage for region %s: %w", region, err)
	}
	stores := &wafv2Stores{
		webACLs:          wafstore.NewWebACLStore(st, s.accountID, region),
		ruleGroups:       wafstore.NewRuleGroupStore(st, s.accountID, region),
		ipSets:           wafstore.NewIPSetStore(st, s.accountID, region),
		regexPatternSets: wafstore.NewRegexPatternSetStore(st, s.accountID, region),
		associations:     wafstore.NewWebACLAssociationStore(st),
		loggingConfigs:   wafstore.NewLoggingStore(st),
		samples:          wafstore.NewSamplingStore(st),
		tags:             storecommon.NewTagStoreWithRegion(st, "wafv2", region),
		arnBuilder:       wafstore.NewARNBuilder(s.accountID, region),
	}
	actual, _ := s.stores.LoadOrStore(region, stores)
	return actual.(*wafv2Stores).webACLs, nil
}

func (s *WAFv2Service) store(reqCtx *request.RequestContext) (*wafv2Stores, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (*wafv2Stores, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, err
		}
		return &wafv2Stores{
			webACLs:          wafstore.NewWebACLStore(storage, reqCtx.GetAccountID(), reqCtx.GetRegion()),
			ruleGroups:       wafstore.NewRuleGroupStore(storage, reqCtx.GetAccountID(), reqCtx.GetRegion()),
			ipSets:           wafstore.NewIPSetStore(storage, reqCtx.GetAccountID(), reqCtx.GetRegion()),
			regexPatternSets: wafstore.NewRegexPatternSetStore(storage, reqCtx.GetAccountID(), reqCtx.GetRegion()),
			associations:     wafstore.NewWebACLAssociationStore(storage),
			loggingConfigs:   wafstore.NewLoggingStore(storage),
			samples:          wafstore.NewSamplingStore(storage),
			tags:             storecommon.NewTagStoreWithRegion(storage, "wafv2", reqCtx.GetRegion()),
			arnBuilder:       wafstore.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()),
		}, nil
	})
}

// isGlobalScope reports whether the Scope parameter names the CloudFront
// (global) partition of the wafv2 resource space.
func isGlobalScope(scope string) bool {
	return strings.EqualFold(scope, "CLOUDFRONT")
}

// storeRegionForScope maps a caller region and a wafv2 Scope to the
// region whose store bundle owns the scope's resources: the CloudFront
// scope is served by us-east-1, every other scope by the caller's
// region.
func storeRegionForScope(callerRegion, scope string) string {
	if isGlobalScope(scope) {
		return arn.WAFCloudFrontRegion
	}
	return callerRegion
}

// storeForScope returns the store bundle that owns resources created
// under the given Scope. The AWS WAFv2 API serves the CloudFront scope
// from the us-east-1 endpoint and its ARNs carry us-east-1 whatever
// region the call arrives from, so CLOUDFRONT-scope resources live in
// the us-east-1 store; regional resources live in the request region.
// Every ARN-based reader (inspection, sampling, the CloudFront WebACL
// validation) resolves the store from the ARN's own region, so routing
// writes through this rule keeps both sides on the same partition.
func (s *WAFv2Service) storeForScope(reqCtx *request.RequestContext, scope string) (*wafv2Stores, error) {
	return s.GetStoresForRegion(storeRegionForScope(reqCtx.GetRegion(), scope))
}

// RegisterHandlers registers all WAFv2 API operation handlers with the dispatcher.
func (s *WAFv2Service) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("wafv2", "CreateWebACL", s.CreateWebACL)
	d.RegisterHandlerForService("wafv2", "GetWebACL", s.GetWebACL)
	d.RegisterHandlerForService("wafv2", "ListWebACLs", s.ListWebACLs)
	d.RegisterHandlerForService("wafv2", "UpdateWebACL", s.UpdateWebACL)
	d.RegisterHandlerForService("wafv2", "DeleteWebACL", s.DeleteWebACL)
	d.RegisterHandlerForService("wafv2", "CheckCapacity", s.CheckCapacity)

	d.RegisterHandlerForService("wafv2", "CreateRuleGroup", s.CreateRuleGroup)
	d.RegisterHandlerForService("wafv2", "GetRuleGroup", s.GetRuleGroup)
	d.RegisterHandlerForService("wafv2", "ListRuleGroups", s.ListRuleGroups)
	d.RegisterHandlerForService("wafv2", "UpdateRuleGroup", s.UpdateRuleGroup)
	d.RegisterHandlerForService("wafv2", "DeleteRuleGroup", s.DeleteRuleGroup)

	d.RegisterHandlerForService("wafv2", "CreateIPSet", s.CreateIPSet)
	d.RegisterHandlerForService("wafv2", "GetIPSet", s.GetIPSet)
	d.RegisterHandlerForService("wafv2", "ListIPSets", s.ListIPSets)
	d.RegisterHandlerForService("wafv2", "UpdateIPSet", s.UpdateIPSet)
	d.RegisterHandlerForService("wafv2", "DeleteIPSet", s.DeleteIPSet)

	d.RegisterHandlerForService("wafv2", "CreateRegexPatternSet", s.CreateRegexPatternSet)
	d.RegisterHandlerForService("wafv2", "GetRegexPatternSet", s.GetRegexPatternSet)
	d.RegisterHandlerForService("wafv2", "ListRegexPatternSets", s.ListRegexPatternSets)
	d.RegisterHandlerForService("wafv2", "UpdateRegexPatternSet", s.UpdateRegexPatternSet)
	d.RegisterHandlerForService("wafv2", "DeleteRegexPatternSet", s.DeleteRegexPatternSet)

	d.RegisterHandlerForService("wafv2", "AssociateWebACL", s.AssociateWebACL)
	d.RegisterHandlerForService("wafv2", "DisassociateWebACL", s.DisassociateWebACL)
	d.RegisterHandlerForService("wafv2", "ListResourcesForWebACL", s.ListResourcesForWebACL)
	d.RegisterHandlerForService("wafv2", "GetWebACLForResource", s.GetWebACLForResource)

	d.RegisterHandlerForService("wafv2", "PutLoggingConfiguration", s.PutLoggingConfiguration)
	d.RegisterHandlerForService("wafv2", "GetLoggingConfiguration", s.GetLoggingConfiguration)
	d.RegisterHandlerForService("wafv2", "DeleteLoggingConfiguration", s.DeleteLoggingConfiguration)
	d.RegisterHandlerForService("wafv2", "ListLoggingConfigurations", s.ListLoggingConfigurations)

	d.RegisterHandlerForService("wafv2", "TagResource", s.TagResource)
	d.RegisterHandlerForService("wafv2", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("wafv2", "ListTagsForResource", s.ListTagsForResource)

	d.RegisterHandlerForService("wafv2", "ListAvailableManagedRuleGroups", s.ListAvailableManagedRuleGroups)
	d.RegisterHandlerForService("wafv2", "DescribeManagedRuleGroup", s.DescribeManagedRuleGroup)
	d.RegisterHandlerForService("wafv2", "ListAvailableManagedRuleGroupVersions", s.ListAvailableManagedRuleGroupVersions)

	d.RegisterHandlerForService("wafv2", "GetSampledRequests", s.GetSampledRequests)
	d.RegisterHandlerForService("wafv2", "GetRateBasedStatementManagedKeys", s.GetRateBasedStatementManagedKeys)
}

// CheckCapacity calculates the capacity consumed by the specified rules.
// Each statement type has a base WCU (Web ACL Capacity Unit) cost per
// AWS documentation. The total is the sum across all rules.
func (s *WAFv2Service) CheckCapacity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rules := convertRules(req.Parameters["Rules"])
	var total int64
	for _, rule := range rules {
		if rule != nil {
			total += calculateStatementCapacity(rule.Statement)
		}
	}
	return map[string]interface{}{
		"Capacity": total,
	}, nil
}

func isCloudFrontResource(resourceArn string) bool {
	return arn.GetServiceFromARN(resourceArn) == "cloudfront"
}

// allAssociationStoresForConnect is the admin-console counterpart of
// allAssociationStores (which lives on the Core layer in
// association_core.go). The connect handler only carries a region
// header, so the region association store is resolved through
// GetStoresForRegion and the global-scope store through the storage
// manager directly, with the same sync.Map caching as the HTTP path.
func (s *WAFv2Service) allAssociationStoresForConnect(header http.Header) ([]*wafstore.WebACLAssociationStore, error) {
	region := defaults.GetRegionFromHeader(header)
	stores, err := s.GetStoresForRegion(region)
	if err != nil {
		return nil, err
	}
	result := []*wafstore.WebACLAssociationStore{stores.associations}
	if cached, ok := s.stores.Load(wafv2GlobalAssocKey); ok {
		if typed, ok := cached.(*wafstore.WebACLAssociationStore); ok {
			return append(result, typed), nil
		}
	}
	if s.storageManager == nil {
		return nil, fmt.Errorf("wafv2 storage manager not initialised")
	}
	globalStorage, err := s.storageManager.GetGlobalStorage()
	if err != nil {
		return nil, err
	}
	store := wafstore.NewWebACLAssociationStore(globalStorage)
	if actual, loaded := s.stores.LoadOrStore(wafv2GlobalAssocKey, store); loaded {
		if typed, ok := actual.(*wafstore.WebACLAssociationStore); ok {
			return append(result, typed), nil
		}
	}
	return append(result, store), nil
}

func generateID() (string, error) {
	return uuid.New().String(), nil
}

func newAPIError(code, message string, httpStatus int) error {
	return awserrors.NewAWSError(code, message, httpStatus)
}

func notFoundError(resource string) error {
	return awserrors.NewAWSError("WAFNonexistentItemException", fmt.Sprintf("%s not found", resource), http.StatusNotFound)
}

func invalidParamError(msg string) error {
	return awserrors.NewAWSError("WAFInvalidParameterException", msg, http.StatusBadRequest)
}

func lockTokenError() error {
	return awserrors.NewAWSError("WAFInvalidLockTokenException", "The provided lock token is not valid", http.StatusBadRequest)
}

func limitsExceededError(consumed int64) error {
	return awserrors.NewAWSError("WAFLimitsExceededException",
		fmt.Sprintf("AWS WAF couldn't perform the operation because your web ACL capacity would exceed the maximum of %d WCUs (requested rules consume %d WCUs)", wafstore.MaxWebACLCapacity, consumed),
		http.StatusBadRequest)
}

func associatedItemError(msg string) error {
	return awserrors.NewAWSError("WAFAssociatedItemException",
		fmt.Sprintf("AWS WAF couldn't perform the operation because your resource is being used by another resource or it's associated with another resource. %s", msg),
		http.StatusBadRequest)
}

// ensureNotAssociated checks every association store (regional and
// global CloudFront scope) for associations referencing the given web
// ACL ARN. Deleting an associated WebACL must fail with
// WAFAssociatedItemException, mirroring the AWS behaviour.
func ensureNotAssociated(assocStores []*wafstore.WebACLAssociationStore, webACLArn string) error {
	if webACLArn == "" {
		return nil
	}
	for _, store := range assocStores {
		if store == nil {
			continue
		}
		assocs, err := store.GetByWebACLArn(webACLArn)
		if err != nil {
			return err
		}
		if len(assocs) > 0 {
			return associatedItemError(fmt.Sprintf("WebACL %s is still associated with %d resource(s).", webACLArn, len(assocs)))
		}
	}
	return nil
}

// calculateRulesCapacity sums the WCU capacity of all rules.
func calculateRulesCapacity(rules []*wafstore.Rule) int64 {
	var total int64
	for _, rule := range rules {
		if rule != nil {
			total += calculateStatementCapacity(rule.Statement)
		}
	}
	return total
}
