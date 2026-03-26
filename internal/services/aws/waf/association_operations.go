package waf

import (
	"context"
	"fmt"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// AssociateWebACL associates a web ACL with a resource.
func (s *WAFService) AssociateWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	webACLArn := request.GetStringParam(req.Parameters, "WebACLArn")
	if webACLArn == "" {
		return nil, NewAPIError("InvalidParameterValue", "WebACLArn is required", 400)
	}

	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, NewAPIError("InvalidParameterValue", "ResourceArn is required", 400)
	}

	err = stores.associations.Associate(webACLArn, resourceArn)
	if err != nil {
		return nil, fmt.Errorf("failed to associate web acl: %w", err)
	}

	return response.EmptyResponse(), nil
}

// DisassociateWebACL disassociates a web ACL from a resource.
func (s *WAFService) DisassociateWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, NewAPIError("InvalidParameterValue", "ResourceArn is required", 400)
	}

	assoc, err := stores.associations.GetByResourceArn(resourceArn)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, NewAPIError("WAFNonexistentItemException", "WebACL association not found", 400)
		}
		return nil, err
	}

	err = stores.associations.Disassociate(assoc.WebACLArn, resourceArn)
	if err != nil {
		return nil, fmt.Errorf("failed to disassociate web acl: %w", err)
	}

	return response.EmptyResponse(), nil
}

// ListResourcesAssociatedToWebACL returns the resources associated with a web ACL.
func (s *WAFService) ListResourcesAssociatedToWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	webACLArn := request.GetStringParam(req.Parameters, "WebACLArn")
	if webACLArn == "" {
		webACLArn = request.GetStringParam(req.Parameters, "Id")
	}
	if webACLArn == "" {
		return nil, NewAPIError("InvalidParameterValue", "WebACLArn is required", 400)
	}

	associations, err := stores.associations.GetByWebACLArn(webACLArn)
	if err != nil {
		return nil, err
	}

	resources := make([]string, 0, len(associations))
	for _, assoc := range associations {
		resources = append(resources, assoc.ResourceArn)
	}

	return map[string]interface{}{
		"ResourceArns": resources,
	}, nil
}
