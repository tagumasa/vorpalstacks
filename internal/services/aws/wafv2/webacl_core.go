package wafv2

import (
	"fmt"

	"vorpalstacks/internal/core/logs"
	wafstore "vorpalstacks/internal/store/aws/waf"
	"vorpalstacks/internal/utils/aws/types"
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
	if input.DefaultAction != nil {
		if err := validateDefaultAction(input.DefaultAction); err != nil {
			return nil, err
		}
	}

	capacity := int64(1500)

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
// shared by the HTTP API and the admin gRPC-Web handler.
func (s *WAFv2Service) deleteWebACLCore(stores *wafv2Stores, id, lockToken string) (*wafstore.WebACL, error) {
	if id == "" {
		return nil, invalidParamError("Id is required")
	}
	if lockToken == "" {
		return nil, invalidParamError("LockToken is required")
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
