// Package waf provides WAF (Web Application Firewall) service operations for vorpalstacks.
package waf

import (
	"context"

	"github.com/google/uuid"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// wafStores holds the various WAF stores.
type wafStores struct {
	webACLs          *wafstore.WebACLStore
	ruleGroups       *wafstore.RuleGroupStore
	ipSets           *wafstore.IPSetStore
	regexPatternSets *wafstore.RegexPatternSetStore
	associations     *wafstore.WebACLAssociationStore
	loggingConfigs   *wafstore.LoggingStore
	arnBuilder       *wafstore.ARNBuilder
}

// WAFService provides AWS WAF (Web Application Firewall) operations.
type WAFService struct{}

// NewWAFService creates a new WAF service instance.
func NewWAFService(store storage.BasicStorage, accountID, region string) *WAFService {
	return &WAFService{}
}

func (s *WAFService) store(reqCtx *request.RequestContext) (*wafStores, error) {
	if stores := reqCtx.GetWAFStores(); stores != nil {
		return &wafStores{
			webACLs:          stores.WebACLStore().Raw(),
			ruleGroups:       stores.RuleGroupStore().Raw(),
			ipSets:           stores.IPSetStore().Raw(),
			regexPatternSets: stores.RegexPatternSetStore().Raw(),
			associations:     stores.AssociationStore().Raw(),
			loggingConfigs:   stores.LoggingStore().Raw(),
			arnBuilder:       wafstore.NewARNBuilder(reqCtx.GetAccountID(), ""),
		}, nil
	}
	storage, err := reqCtx.GetGlobalStorage()
	if err != nil {
		return nil, err
	}
	return &wafStores{
		webACLs:          wafstore.NewWebACLStore(storage, reqCtx.GetAccountID(), ""),
		ruleGroups:       wafstore.NewRuleGroupStore(storage, reqCtx.GetAccountID(), ""),
		ipSets:           wafstore.NewIPSetStore(storage, reqCtx.GetAccountID(), ""),
		regexPatternSets: wafstore.NewRegexPatternSetStore(storage, reqCtx.GetAccountID(), ""),
		associations:     wafstore.NewWebACLAssociationStore(storage),
		loggingConfigs:   wafstore.NewLoggingStore(storage),
		arnBuilder:       wafstore.NewARNBuilder(reqCtx.GetAccountID(), ""),
	}, nil
}

// RegisterHandlers registers the WAF service handlers with the dispatcher.
func (s *WAFService) RegisterHandlers(d *dispatcher.Dispatcher) {
	d.RegisterHandlerForService("waf", "CreateWebACL", s.CreateWebACL)
	d.RegisterHandlerForService("waf", "GetWebACL", s.GetWebACL)
	d.RegisterHandlerForService("waf", "ListWebACLs", s.ListWebACLs)
	d.RegisterHandlerForService("waf", "UpdateWebACL", s.UpdateWebACL)
	d.RegisterHandlerForService("waf", "DeleteWebACL", s.DeleteWebACL)
	d.RegisterHandlerForService("waf", "CheckCapacity", s.CheckCapacity)

	d.RegisterHandlerForService("waf", "CreateRuleGroup", s.CreateRuleGroup)
	d.RegisterHandlerForService("waf", "GetRuleGroup", s.GetRuleGroup)
	d.RegisterHandlerForService("waf", "ListRuleGroups", s.ListRuleGroups)
	d.RegisterHandlerForService("waf", "UpdateRuleGroup", s.UpdateRuleGroup)
	d.RegisterHandlerForService("waf", "DeleteRuleGroup", s.DeleteRuleGroup)

	d.RegisterHandlerForService("waf", "CreateRule", s.CreateRule)
	d.RegisterHandlerForService("waf", "GetRule", s.GetRule)
	d.RegisterHandlerForService("waf", "UpdateRule", s.UpdateRule)
	d.RegisterHandlerForService("waf", "DeleteRule", s.DeleteRule)

	d.RegisterHandlerForService("waf", "CreateIPSet", s.CreateIPSet)
	d.RegisterHandlerForService("waf", "GetIPSet", s.GetIPSet)
	d.RegisterHandlerForService("waf", "ListIPSets", s.ListIPSets)
	d.RegisterHandlerForService("waf", "UpdateIPSet", s.UpdateIPSet)
	d.RegisterHandlerForService("waf", "DeleteIPSet", s.DeleteIPSet)

	d.RegisterHandlerForService("waf", "CreateRegexPatternSet", s.CreateRegexPatternSet)
	d.RegisterHandlerForService("waf", "GetRegexPatternSet", s.GetRegexPatternSet)
	d.RegisterHandlerForService("waf", "ListRegexPatternSets", s.ListRegexPatternSets)
	d.RegisterHandlerForService("waf", "UpdateRegexPatternSet", s.UpdateRegexPatternSet)
	d.RegisterHandlerForService("waf", "DeleteRegexPatternSet", s.DeleteRegexPatternSet)

	d.RegisterHandlerForService("waf", "AssociateWebACL", s.AssociateWebACL)
	d.RegisterHandlerForService("waf", "DisassociateWebACL", s.DisassociateWebACL)
	d.RegisterHandlerForService("waf", "ListResourcesForWebACL", s.ListResourcesAssociatedToWebACL)

	d.RegisterHandlerForService("waf", "PutLoggingConfiguration", s.PutLoggingConfiguration)
	d.RegisterHandlerForService("waf", "GetLoggingConfiguration", s.GetLoggingConfiguration)
	d.RegisterHandlerForService("waf", "DeleteLoggingConfiguration", s.DeleteLoggingConfiguration)
	d.RegisterHandlerForService("waf", "ListLoggingConfigurations", s.ListLoggingConfigurations)

	d.RegisterHandlerForService("waf", "TagResource", s.TagResource)
	d.RegisterHandlerForService("waf", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("waf", "ListTagsForResource", s.ListTagsForResource)

	d.RegisterHandlerForService("waf", "GetChangeToken", s.GetChangeToken)
	d.RegisterHandlerForService("waf", "GetChangeTokenStatus", s.GetChangeTokenStatus)

	d.RegisterHandlerForService("waf", "ListAvailableManagedRuleGroups", s.ListAvailableManagedRuleGroups)
	d.RegisterHandlerForService("waf", "DescribeManagedRuleGroup", s.DescribeManagedRuleGroup)
	d.RegisterHandlerForService("waf", "ListAvailableManagedRuleGroupVersions", s.ListAvailableManagedRuleGroupVersions)
}

// GetChangeToken generates and returns a new change token for WAF resource updates.
func (s *WAFService) GetChangeToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	token := uuid.New().String()
	return map[string]interface{}{
		"ChangeToken": token,
	}, nil
}

// GetChangeTokenStatus returns the status of the most recently issued change token.
func (s *WAFService) GetChangeTokenStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"ChangeTokenStatus": "PROVISIONED",
	}, nil
}
