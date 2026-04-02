package wafv2

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

func (s *WAFv2Service) AssociateWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	webACLArn := request.GetStringParam(req.Parameters, "WebACLArn")
	if webACLArn == "" {
		return nil, validationError("WebACLArn is required")
	}

	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, validationError("ResourceArn is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = stores.webACLs.GetByARN(webACLArn)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL")
		}
		return nil, err
	}

	assocStore, err := s.associationStoreFor(reqCtx, resourceArn)
	if err != nil {
		return nil, err
	}

	if err := assocStore.Associate(webACLArn, resourceArn); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *WAFv2Service) DisassociateWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, validationError("ResourceArn is required")
	}

	assocStore, err := s.associationStoreFor(reqCtx, resourceArn)
	if err != nil {
		return nil, err
	}

	assoc, err := assocStore.GetByResourceArn(resourceArn)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL association")
		}
		return nil, err
	}

	if err := assocStore.Disassociate(assoc.WebACLArn, resourceArn); err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL association")
		}
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *WAFv2Service) ListResourcesForWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	webACLArn := request.GetStringParam(req.Parameters, "WebACLArn")
	if webACLArn == "" {
		return nil, validationError("WebACLArn is required")
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
				seen[assoc.ResourceArn] = true
				resources = append(resources, assoc.ResourceArn)
			}
		}
	}

	return map[string]interface{}{
		"ResourceArns": resources,
	}, nil
}

func (s *WAFv2Service) GetWebACLForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, validationError("ResourceArn is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
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

	webACL, err := stores.webACLs.GetByARN(assoc.WebACLArn)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL")
		}
		return nil, err
	}

	return map[string]interface{}{
		"WebACL": map[string]interface{}{
			"Id":               webACL.ID,
			"Name":             webACL.Name,
			"ARN":              webACL.ARN,
			"Description":      webACL.Description,
			"DefaultAction":    convertActionToResponse(webACL.DefaultAction),
			"Rules":            convertRulesToResponse(webACL.Rules),
			"VisibilityConfig": convertVisibilityConfigToResponse(webACL.VisibilityConfig),
			"Capacity":         webACL.Capacity,
			"Scope":            webACL.Scope,
			"LockToken":        webACL.LockToken,
		},
	}, nil
}
