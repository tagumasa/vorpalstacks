package wafv2

import (
	"context"
	"fmt"
	"strings"

	"vorpalstacks/internal/common/request"
	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	wafstore "vorpalstacks/internal/store/aws/waf"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// CreateWebACLInput is the transport-agnostic input for creating a
// WebACL. The admin console sets only Name/Description/Scope; the HTTP
// API additionally provides DefaultAction, Rules, VisibilityConfig and
// other rich fields. Store-typed fields are safe because the admin
// handler constructs the struct by field name and never references
// the store packages.
type CreateWebACLInput struct {
	Name        string
	Description string
	Scope       string

	// Rich fields used by the HTTP API; left nil by the admin console.
	DefaultAction          *wafstore.Action
	Rules                  []*wafstore.Rule
	VisibilityConfig       *wafstore.VisibilityConfig
	CustomResponseBodies   interface{}
	CaptchaConfig          interface{}
	ChallengeConfig        interface{}
	TokenDomains           interface{}
	AssociationConfig      interface{}
	ApplicationConfig      interface{}
	MonetizationConfig     interface{}
	DataProtectionConfig   interface{}
	OnSourceDDoSProtection interface{}
	Tags                   []types.Tag
}

// ListWebACLsInput is the transport-agnostic input for listing WebACLs.
type ListWebACLsInput struct {
	Scope      string
	Limit      int
	NextMarker string
}

// createWebACLCore is the single entry point for creating a WebACL,
// shared by the HTTP API and the admin gRPC-Web handler.
func (s *WAFv2Service) createWebACLCore(stores *wafv2Stores, input CreateWebACLInput) (*wafstore.WebACL, error) {
	if err := validateEntityName(input.Name); err != nil {
		return nil, err
	}
	if err := validateScope(input.Scope); err != nil {
		return nil, err
	}
	if err := validateEntityDescription(input.Description); err != nil {
		return nil, err
	}
	if err := validateTokenDomains(input.TokenDomains); err != nil {
		return nil, err
	}
	if err := validateCustomResponseBodies(input.CustomResponseBodies); err != nil {
		return nil, err
	}
	if err := validateImmunityConfig(input.CaptchaConfig, "Captcha"); err != nil {
		return nil, err
	}
	if err := validateImmunityConfig(input.ChallengeConfig, "Challenge"); err != nil {
		return nil, err
	}
	if err := validateMonetizationConfig(input.MonetizationConfig); err != nil {
		return nil, err
	}
	if err := validateMonetizeRules(input.Rules, input.Scope, input.MonetizationConfig); err != nil {
		return nil, err
	}
	if err := validateRuleGroupReferenceOverrides(stores, input.Rules); err != nil {
		return nil, err
	}
	// DefaultAction and VisibilityConfig are @required on
	// CreateWebACLRequest in the Smithy model; the admin console
	// transport synthesises defaults before calling this core so the
	// single contract stays identical for both transports.
	if err := validateDefaultAction(input.DefaultAction); err != nil {
		return nil, err
	}
	if err := validateVisibilityConfig(input.VisibilityConfig); err != nil {
		return nil, err
	}

	// Capacity is a read-only consumed-capacity value on the WebACL
	// shape: compute it from the submitted rules and enforce the AWS
	// WCU quota rather than persisting a placeholder.
	capacity := calculateRulesCapacity(input.Rules)
	if capacity > wafstore.MaxWebACLCapacity {
		return nil, limitsExceededError(capacity)
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	webACL := &wafstore.WebACL{
		ID:                     id,
		Name:                   input.Name,
		Description:            input.Description,
		Scope:                  input.Scope,
		Capacity:               capacity,
		Rules:                  input.Rules,
		DefaultAction:          input.DefaultAction,
		VisibilityConfig:       input.VisibilityConfig,
		CustomResponseBodies:   input.CustomResponseBodies,
		CaptchaConfig:          input.CaptchaConfig,
		ChallengeConfig:        input.ChallengeConfig,
		TokenDomains:           input.TokenDomains,
		AssociationConfig:      input.AssociationConfig,
		ApplicationConfig:      input.ApplicationConfig,
		MonetizationConfig:     input.MonetizationConfig,
		DataProtectionConfig:   input.DataProtectionConfig,
		OnSourceDDoSProtection: input.OnSourceDDoSProtection,
	}

	webACL, err = stores.webACLs.Create(webACL)
	if err != nil {
		if wafstore.IsAlreadyExists(err) {
			return nil, newAPIError("WAFDuplicateItemException", fmt.Sprintf("AWS WAF couldn't perform the operation because some resource in your request is a duplicate of an existing one: %s", input.Name), 400)
		}
		return nil, err
	}

	if len(input.Tags) > 0 {
		if err := stores.tags.TagFromSlice(webACL.ARN, input.Tags); err != nil {
			logs.Warn("failed to persist tags for WebACL", logs.String("id", webACL.ID), logs.Err(err))
		}
	}

	return webACL, nil
}

// deleteWebACLCore is the single entry point for deleting a WebACL,
// shared by the HTTP API and the admin gRPC-Web handler. A WebACL that
// is still associated with a resource cannot be deleted (AWS returns
// WAFAssociatedItemException); association stores are supplied by the
// transport layer because the global-scope store needs the request
// context.
func (s *WAFv2Service) deleteWebACLCore(stores *wafv2Stores, assocStores []*wafstore.WebACLAssociationStore, id, lockToken string) (*wafstore.WebACL, error) {
	if id == "" {
		return nil, invalidParamError("Id is required")
	}
	if lockToken == "" {
		return nil, invalidParamError("LockToken is required")
	}

	existing, err := stores.webACLs.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL")
		}
		return nil, err
	}

	if err := ensureNotAssociated(assocStores, existing.ARN); err != nil {
		return nil, err
	}

	deleted, err := stores.webACLs.Delete(id, lockToken)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL")
		}
		if wafstore.IsLockTokenMismatch(err) {
			return nil, lockTokenError()
		}
		return nil, err
	}

	if deleted.ARN != "" {
		if err := stores.tags.Delete(deleted.ARN); err != nil {
			logs.Warn("failed to clean up tags for deleted WebACL", logs.String("id", id), logs.Err(err))
		}
	}

	return deleted, nil
}

// listWebACLsCore is the single entry point for listing WebACLs,
// shared by the HTTP API and the admin gRPC-Web handler.
func (s *WAFv2Service) listWebACLsCore(stores *wafv2Stores, input ListWebACLsInput) (*wafstore.WebACLListResult, error) {
	if err := validateScope(input.Scope); err != nil {
		return nil, err
	}

	limit := input.Limit
	if limit <= 0 {
		limit = 100
	}

	result, err := stores.webACLs.List(input.NextMarker, limit, input.Scope)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// getWebACLCore is the single entry point for retrieving a WebACL by ID.
// The request context is taken directly because the required-Id check
// precedes the store acquisition in the original failure precedence. The
// scope routes the store bundle: the CloudFront scope lives in the
// us-east-1 store whatever region the call arrives from.
func (s *WAFv2Service) getWebACLCore(reqCtx *request.RequestContext, id, scope string) (*wafstore.WebACL, error) {
	if id == "" {
		return nil, invalidParamError("Id is required")
	}

	stores, err := s.storeForScope(reqCtx, scope)
	if err != nil {
		return nil, err
	}
	webACL, err := stores.webACLs.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL")
		}
		return nil, err
	}

	return webACL, nil
}

// WebACLUpdateInput is the transport-agnostic input for updating a WebACL.
// DefaultAction, VisibilityConfig and Rules travel raw because their wire
// conversion must run after the store acquisition, matching the original
// failure precedence; the nine optional members travel raw with their
// presence-based apply semantics.
type WebACLUpdateInput struct {
	Id                        string
	LockToken                 string
	Scope                     string
	DefaultActionRaw          interface{}
	VisibilityConfigRaw       interface{}
	RulesRaw                  interface{}
	Description               string
	CustomResponseBodies      interface{}
	CaptchaConfig             interface{}
	ChallengeConfig           interface{}
	TokenDomains              interface{}
	AssociationConfig         interface{}
	ApplicationConfig         interface{}
	MonetizationConfig        interface{}
	DataProtectionConfig      interface{}
	OnSourceDDoSProtectionCfg interface{}
}

// updateWebACLCore is the single entry point for updating a WebACL,
// returning the next lock token.
func (s *WAFv2Service) updateWebACLCore(reqCtx *request.RequestContext, in WebACLUpdateInput) (string, error) {
	if in.Id == "" {
		return "", invalidParamError("Id is required")
	}

	if in.LockToken == "" {
		return "", invalidParamError("LockToken is required")
	}

	stores, err := s.storeForScope(reqCtx, in.Scope)
	if err != nil {
		return "", err
	}

	// UpdateWebACL is a full-replace operation per the Smithy model
	// documentation ("This operation completely replaces the mutable
	// specifications"): DefaultAction and VisibilityConfig are required
	// on every call, Rules omitted means an empty rule list, and
	// capacity is never accepted from the request — the model has no
	// Capacity member on UpdateWebACLRequest — but is always recomputed
	// from the resulting rule set.
	vcRaw := in.VisibilityConfigRaw
	if vcRaw == nil {
		return "", invalidParamError("VisibilityConfig is required")
	}
	vcMap, ok := vcRaw.(map[string]interface{})
	if !ok {
		return "", invalidParamError("VisibilityConfig must be an object")
	}
	visibilityConfig := convertVisibilityConfig(vcMap)
	if err := validateVisibilityConfig(visibilityConfig); err != nil {
		return "", err
	}

	daRaw := in.DefaultActionRaw
	if daRaw == nil {
		return "", invalidParamError("DefaultAction is required")
	}
	daMap, ok := daRaw.(map[string]interface{})
	if !ok {
		return "", invalidParamError("DefaultAction must be an object")
	}
	daAction := convertAction(daMap)
	if err := validateDefaultAction(daAction); err != nil {
		return "", err
	}

	var rules []*wafstore.Rule
	if in.RulesRaw != nil {
		parsed, pErr := parseRules(in.RulesRaw)
		if pErr != nil {
			return "", pErr
		}
		rules = parsed
	}

	capacity := calculateRulesCapacity(rules)
	if capacity > wafstore.MaxWebACLCapacity {
		return "", limitsExceededError(capacity)
	}

	if in.TokenDomains != nil {
		if err := validateTokenDomains(in.TokenDomains); err != nil {
			return "", err
		}
	}
	if in.CustomResponseBodies != nil {
		if err := validateCustomResponseBodies(in.CustomResponseBodies); err != nil {
			return "", err
		}
	}
	if err := validateImmunityConfig(in.CaptchaConfig, "Captcha"); err != nil {
		return "", err
	}
	if err := validateImmunityConfig(in.ChallengeConfig, "Challenge"); err != nil {
		return "", err
	}
	if err := validateMonetizationConfig(in.MonetizationConfig); err != nil {
		return "", err
	}
	// Monetize rules may ride on a configuration that an earlier update
	// left on the web ACL, so the effective configuration falls back to
	// the stored one when this call omits it.
	effectiveMonetization := in.MonetizationConfig
	if effectiveMonetization == nil && rulesContainMonetize(rules) {
		if existing, getErr := stores.webACLs.Get(in.Id); getErr == nil {
			effectiveMonetization = existing.MonetizationConfig
		}
	}
	if err := validateMonetizeRules(rules, in.Scope, effectiveMonetization); err != nil {
		return "", err
	}
	if err := validateRuleGroupReferenceOverrides(stores, rules); err != nil {
		return "", err
	}

	updated, err := stores.webACLs.Update(in.Id, in.LockToken, capacity, rules, daAction, visibilityConfig, in.Description, func(webACL *wafstore.WebACL) {
		if in.CustomResponseBodies != nil {
			webACL.CustomResponseBodies = in.CustomResponseBodies
		}
		if in.CaptchaConfig != nil {
			webACL.CaptchaConfig = in.CaptchaConfig
		}
		if in.ChallengeConfig != nil {
			webACL.ChallengeConfig = in.ChallengeConfig
		}
		if in.TokenDomains != nil {
			webACL.TokenDomains = in.TokenDomains
		}
		if in.AssociationConfig != nil {
			webACL.AssociationConfig = in.AssociationConfig
		}
		if in.ApplicationConfig != nil {
			webACL.ApplicationConfig = in.ApplicationConfig
		}
		if in.MonetizationConfig != nil {
			webACL.MonetizationConfig = in.MonetizationConfig
		}
		if in.DataProtectionConfig != nil {
			webACL.DataProtectionConfig = in.DataProtectionConfig
		}
		if in.OnSourceDDoSProtectionCfg != nil {
			webACL.OnSourceDDoSProtection = in.OnSourceDDoSProtectionCfg
		}
	})
	if err != nil {
		if wafstore.IsLockTokenMismatch(err) {
			return "", lockTokenError()
		}
		return "", err
	}

	return updated.LockToken, nil
}

// WebACLExistsForInvoker reports whether a Web ACL with the given ARN or
// ID exists. It backs the cross-service WebACLExists invoker method
// (e.g. CloudFront validating a distribution's WebACLId). A WAFv2 ARN
// carries the owning region in its fourth field; a bare WAF Classic
// style ID is looked up in the service's default region.
func (s *WAFv2Service) WebACLExistsForInvoker(ctx context.Context, webACLIdOrArn string) bool {
	_ = ctx
	if webACLIdOrArn == "" {
		return false
	}
	region := s.region
	if strings.HasPrefix(webACLIdOrArn, "arn:") {
		_, _, arnRegion, _, resource := svcarn.SplitARN(webACLIdOrArn)
		if arnRegion != "" {
			region = arnRegion
		}
		stores, err := s.GetStoresForRegion(region)
		if err != nil {
			return false
		}
		webACL, err := stores.webACLs.GetByARN(webACLIdOrArn)
		if err == nil && webACL != nil {
			return true
		}
		// Fall back to the ID tail for ARNs whose stored form differs from
		// the caller's rendering (for example a region-less ARN).
		idParts := strings.Split(resource, "/")
		id := idParts[len(idParts)-1]
		webACL, err = stores.webACLs.Get(id)
		return err == nil && webACL != nil
	}
	stores, err := s.GetStoresForRegion(region)
	if err != nil {
		return false
	}
	webACL, err := stores.webACLs.Get(webACLIdOrArn)
	return err == nil && webACL != nil
}
