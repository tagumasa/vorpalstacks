package wafv2

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
	wafstore "vorpalstacks/internal/store/aws/waf"
	awserrors "vorpalstacks/internal/utils/aws/errors"
)

var wafv2GlobalAssocKey struct{}

type wafv2Stores struct {
	webACLs          *wafstore.WebACLStore
	ruleGroups       *wafstore.RuleGroupStore
	ipSets           *wafstore.IPSetStore
	regexPatternSets *wafstore.RegexPatternSetStore
	associations     *wafstore.WebACLAssociationStore
	loggingConfigs   *wafstore.LoggingStore
	arnBuilder       *wafstore.ARNBuilder
}

// WAFv2Service implements the AWS WAF v2 API operations.
type WAFv2Service struct {
	stores sync.Map // region → *wafv2Stores
}

// NewWAFv2Service creates a new WAFv2Service instance.
func NewWAFv2Service(accountID, region string) *WAFv2Service {
	return &WAFv2Service{}
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
			arnBuilder:       wafstore.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()),
		}, nil
	})
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

	d.RegisterHandlerForService("wafv2", "GetChangeToken", s.GetChangeToken)
	d.RegisterHandlerForService("wafv2", "GetChangeTokenStatus", s.GetChangeTokenStatus)

	d.RegisterHandlerForService("wafv2", "ListAvailableManagedRuleGroups", s.ListAvailableManagedRuleGroups)
	d.RegisterHandlerForService("wafv2", "DescribeManagedRuleGroup", s.DescribeManagedRuleGroup)
	d.RegisterHandlerForService("wafv2", "ListAvailableManagedRuleGroupVersions", s.ListAvailableManagedRuleGroupVersions)
}

// GetChangeToken generates and returns a new change token for WAFv2 write operations.
func (s *WAFv2Service) GetChangeToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	token := uuid.New().String()
	return map[string]interface{}{
		"ChangeToken": token,
	}, nil
}

// GetChangeTokenStatus returns the current status of the most recent change token.
func (s *WAFv2Service) GetChangeTokenStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"ChangeTokenStatus": "PROVISIONED",
	}, nil
}

// CheckCapacity returns the capacity consumed by the specified rules.
func (s *WAFv2Service) CheckCapacity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"Capacity": 10,
	}, nil
}

func isCloudFrontResource(resourceArn string) bool {
	return strings.Contains(resourceArn, ":cloudfront:")
}

func (s *WAFv2Service) associationStoreFor(reqCtx *request.RequestContext, resourceArn string) (*wafstore.WebACLAssociationStore, error) {
	if isCloudFrontResource(resourceArn) {
		if cached, ok := s.stores.Load(wafv2GlobalAssocKey); ok {
			if typed, ok := cached.(*wafstore.WebACLAssociationStore); ok {
				return typed, nil
			}
		}
		globalStorage, err := reqCtx.GetGlobalStorage()
		if err != nil {
			return nil, err
		}
		store := wafstore.NewWebACLAssociationStore(globalStorage)
		if actual, loaded := s.stores.LoadOrStore(wafv2GlobalAssocKey, store); loaded {
			if typed, ok := actual.(*wafstore.WebACLAssociationStore); ok {
				return typed, nil
			}
		}
		return store, nil
	}
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return stores.associations, nil
}

func (s *WAFv2Service) allAssociationStores(reqCtx *request.RequestContext) ([]*wafstore.WebACLAssociationStore, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result := []*wafstore.WebACLAssociationStore{stores.associations}
	globalStorage, err := reqCtx.GetGlobalStorage()
	if err == nil {
		if cached, ok := s.stores.Load(wafv2GlobalAssocKey); ok {
			if typed, ok := cached.(*wafstore.WebACLAssociationStore); ok {
				result = append(result, typed)
				return result, nil
			}
		}
		store := wafstore.NewWebACLAssociationStore(globalStorage)
		if actual, loaded := s.stores.LoadOrStore(wafv2GlobalAssocKey, store); loaded {
			if typed, ok := actual.(*wafstore.WebACLAssociationStore); ok {
				result = append(result, typed)
				return result, nil
			}
		}
		result = append(result, store)
	}
	return result, nil
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

func validationError(msg string) error {
	return awserrors.NewAWSError("ValidationException", msg, http.StatusBadRequest)
}

func lockTokenError() error {
	return awserrors.NewAWSError("WAFInvalidLockTokenException", "The provided lock token is not valid", http.StatusBadRequest)
}
