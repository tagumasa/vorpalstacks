package wafv2

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// AssociateWebACL associates the specified web ACL with a regional resource.
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

	if err := stores.associations.Associate(webACLArn, resourceArn); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// DisassociateWebACL removes the association between a web ACL and a regional resource.
func (s *WAFv2Service) DisassociateWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, validationError("ResourceArn is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	assoc, err := stores.associations.GetByResourceArn(resourceArn)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL association")
		}
		return nil, err
	}

	if err := stores.associations.Disassociate(assoc.WebACLArn, resourceArn); err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL association")
		}
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListResourcesForWebACL returns all regional resources associated with the specified web ACL.
func (s *WAFv2Service) ListResourcesForWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	webACLArn := request.GetStringParam(req.Parameters, "WebACLArn")
	if webACLArn == "" {
		return nil, validationError("WebACLArn is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	associations, err := stores.associations.GetByWebACLArn(webACLArn)
	if err != nil {
		return nil, err
	}

	resources := make([]interface{}, 0, len(associations))
	for _, assoc := range associations {
		resources = append(resources, map[string]interface{}{
			"ResourceArn": assoc.ResourceArn,
		})
	}

	return map[string]interface{}{
		"ResourceArns": resources,
	}, nil
}
