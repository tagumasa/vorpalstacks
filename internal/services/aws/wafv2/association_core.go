package wafv2

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/request"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// associationStoreFor resolves the association store that owns the given
// resource ARN: CloudFront-scoped resources live in the global storage,
// everything else in the request region (relocated verbatim from
// service.go — store acquisition for the association plane lives on the
// Core layer).
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

// allAssociationStores returns every association store the request can
// see: the request region's store plus the global CloudFront-scope store
// (relocated verbatim from service.go).
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

// associateWebACLCore is the single entry point for associating a WebACL
// with a resource ARN. The request context is taken directly because the
// required-member checks precede the store acquisition in the original
// failure precedence.
func (s *WAFv2Service) associateWebACLCore(ctx context.Context, reqCtx *request.RequestContext, webACLArn, resourceArn string) error {
	if webACLArn == "" {
		return invalidParamError("WebACLArn is required")
	}

	if resourceArn == "" {
		return invalidParamError("ResourceArn is required")
	}

	// The WebACL lives in the partition its ARN names — us-east-1 for the
	// CloudFront scope — whatever region the call arrives from.
	stores, err := s.GetStoresForRegion(s.arnRegion(webACLArn))
	if err != nil {
		return err
	}

	if _, err = stores.webACLs.GetByARN(webACLArn); err != nil {
		if wafstore.IsNotFound(err) {
			return notFoundError("WebACL")
		}
		return err
	}

	if err := s.ensureAssociableResource(ctx, reqCtx.GetRegion(), resourceArn); err != nil {
		return err
	}

	assocStore, err := s.associationStoreFor(reqCtx, resourceArn)
	if err != nil {
		return err
	}

	return assocStore.Associate(webACLArn, resourceArn)
}

// disassociateWebACLCore is the single entry point for removing the
// WebACL association from a resource ARN.
func (s *WAFv2Service) disassociateWebACLCore(reqCtx *request.RequestContext, resourceArn string) error {
	if resourceArn == "" {
		return invalidParamError("ResourceArn is required")
	}

	assocStore, err := s.associationStoreFor(reqCtx, resourceArn)
	if err != nil {
		return err
	}

	if _, err := assocStore.GetByResourceArn(resourceArn); err != nil {
		if wafstore.IsNotFound(err) {
			return notFoundError("WebACL association")
		}
		return err
	}

	if err := assocStore.Disassociate(resourceArn); err != nil {
		if wafstore.IsNotFound(err) {
			return notFoundError("WebACL association")
		}
		return err
	}

	return nil
}

// getWebACLForResourceCore is the single entry point for retrieving the
// WebACL associated with a resource ARN.
func (s *WAFv2Service) getWebACLForResourceCore(reqCtx *request.RequestContext, resourceArn string) (*wafstore.WebACL, error) {
	if resourceArn == "" {
		return nil, invalidParamError("ResourceArn is required")
	}

	assocStore, err := s.associationStoreFor(reqCtx, resourceArn)
	if err != nil {
		return nil, err
	}

	assoc, err := assocStore.GetByResourceArn(resourceArn)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL association for the specified resource")
		}
		return nil, err
	}

	// The associated WebACL is loaded from the partition its own ARN
	// names, which is where AssociateWebACL verified it exists.
	stores, err := s.GetStoresForRegion(s.arnRegion(assoc.WebACLArn))
	if err != nil {
		return nil, err
	}

	webACL, err := stores.webACLs.GetByARN(assoc.WebACLArn)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL")
		}
		return nil, err
	}

	return webACL, nil
}

// listResourcesForWebACLCore is the single entry point for listing the
// resource ARNs associated with a WebACL, optionally filtered by the
// ResourceType enum.
func (s *WAFv2Service) listResourcesForWebACLCore(reqCtx *request.RequestContext, webACLArn, resourceType string) ([]string, error) {
	if webACLArn == "" {
		return nil, invalidParamError("WebACLArn is required")
	}

	if resourceType != "" && !validResourceTypes[resourceType] {
		return nil, invalidParamError(fmt.Sprintf("Unsupported ResourceType: %s", resourceType))
	}

	associationStores, err := s.allAssociationStores(reqCtx)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	resources := make([]string, 0)
	for _, assocStore := range associationStores {
		associations, err := assocStore.GetByWebACLArn(webACLArn)
		if err != nil {
			return nil, err
		}
		for _, assoc := range associations {
			if !seen[assoc.ResourceArn] {
				if resourceType == "" || matchesResourceType(assoc.ResourceArn, resourceType) {
					seen[assoc.ResourceArn] = true
					resources = append(resources, assoc.ResourceArn)
				}
			}
		}
	}

	return resources, nil
}
