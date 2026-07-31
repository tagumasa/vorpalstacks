package wafv2

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// AssociateWebACL associates a WebACL with the specified resource ARN.
func (s *WAFv2Service) AssociateWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	webACLArn := request.GetStringParam(req.Parameters, "WebACLArn")
	if webACLArn == "" {
		return nil, invalidParamError("WebACLArn is required")
	}

	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, invalidParamError("ResourceArn is required")
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

// DisassociateWebACL removes the WebACL association from the specified resource ARN.
func (s *WAFv2Service) DisassociateWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, invalidParamError("ResourceArn is required")
	}

	assocStore, err := s.associationStoreFor(reqCtx, resourceArn)
	if err != nil {
		return nil, err
	}

	if _, err := assocStore.GetByResourceArn(resourceArn); err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL association")
		}
		return nil, err
	}

	if err := assocStore.Disassociate(resourceArn); err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL association")
		}
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListResourcesForWebACL returns all resource ARNs associated with the specified WebACL.
// If ResourceType is provided, results are filtered to only include resources
// of that type.
func (s *WAFv2Service) ListResourcesForWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	webACLArn := request.GetStringParam(req.Parameters, "WebACLArn")
	if webACLArn == "" {
		return nil, invalidParamError("WebACLArn is required")
	}

	resourceType := request.GetStringParam(req.Parameters, "ResourceType")

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

	return map[string]interface{}{
		"ResourceArns": resources,
	}, nil
}

// matchesResourceType checks whether a resource ARN matches the given
// AWS ResourceType enum value.
func matchesResourceType(resourceArn, resourceType string) bool {
	switch resourceType {
	case "APPLICATION_LOAD_BALANCER":
		return strings.Contains(resourceArn, ":elasticloadbalancing:")
	case "API_GATEWAY":
		return strings.Contains(resourceArn, ":apigateway:")
	case "APPSYNC":
		return strings.Contains(resourceArn, ":appsync:")
	case "COGNITIO_USER_POOL":
		return strings.Contains(resourceArn, ":cognito-idp:")
	case "APP_RUNNER_SERVICE":
		return strings.Contains(resourceArn, ":runner:")
	case "VERIFIED_ACCESS_INSTANCE":
		return strings.Contains(resourceArn, ":ec2:") && strings.Contains(resourceArn, "verified-access")
	case "AMPLIFY":
		return strings.Contains(resourceArn, ":amplify:")
	case "AGENTCORE_GATEWAY":
		return strings.Contains(resourceArn, ":agentcore:")
	default:
		return true
	}
}

// GetWebACLForResource returns the WebACL associated with the specified resource ARN.
func (s *WAFv2Service) GetWebACLForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, invalidParamError("ResourceArn is required")
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
