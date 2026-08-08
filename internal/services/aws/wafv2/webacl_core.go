package wafv2

import (
	"fmt"

	"vorpalstacks/internal/core/logs"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// AdminCreateWebACLInput is the transport-agnostic input for creating a
// WebACL from the admin console. It contains only primitive fields so
// that the admin handler does not need to import store packages.
type AdminCreateWebACLInput struct {
	Name        string
	Description string
	Scope       string
}

// AdminListWebACLsInput is the transport-agnostic input for listing
// WebACLs from the admin console.
type AdminListWebACLsInput struct {
	Scope      string
	Limit      int
	NextMarker string
}

// createWebACLCore creates a WebACL using a service-layer DTO. It is
// used by the admin console gRPC-Web handler. The HTTP API handler
// (CreateWebACL in web_acl_operations.go) has its own richer path that
// accepts DefaultAction, Rules, VisibilityConfig, and other fields.
func (s *WAFv2Service) createWebACLCore(stores *wafv2Stores, input AdminCreateWebACLInput) (*wafstore.WebACL, error) {
	if err := validateEntityName(input.Name); err != nil {
		return nil, err
	}
	if err := validateScope(input.Scope); err != nil {
		return nil, err
	}
	if err := validateEntityDescription(input.Description); err != nil {
		return nil, err
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	webACL := &wafstore.WebACL{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		Scope:       input.Scope,
		Capacity:    1500,
	}

	webACL, err = stores.webACLs.Create(webACL)
	if err != nil {
		if wafstore.IsAlreadyExists(err) {
			return nil, newAPIError("WAFDuplicateItemException", fmt.Sprintf("AWS WAF couldn't perform the operation because some resource in your request is a duplicate of an existing one: %s", input.Name), 400)
		}
		return nil, err
	}

	return webACL, nil
}

// deleteWebACLCore deletes a WebACL and cleans up associated tags. It is
// used by the admin console gRPC-Web handler.
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

// listWebACLsCore lists WebACLs using a service-layer DTO. It is used
// by the admin console gRPC-Web handler.
func (s *WAFv2Service) listWebACLsCore(stores *wafv2Stores, input AdminListWebACLsInput) (*wafstore.WebACLListResult, error) {
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
