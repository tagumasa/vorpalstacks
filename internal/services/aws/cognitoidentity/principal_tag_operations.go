package cognitoidentity

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// principal_tag_operations.go — HTTP-plane handlers for the per-provider
// principal tag attribute maps.

// GetPrincipalTagAttributeMap retrieves the principal tag attribute map for an identity provider.
func (s *CognitoIdentityService) GetPrincipalTagAttributeMap(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.getPrincipalTagAttributeMapCore(reqCtx, GetPrincipalTagAttributeMapInput{
		IdentityPoolID:       req.GetParam("IdentityPoolId"),
		IdentityProviderName: req.GetParam("IdentityProviderName"),
	})
	if err != nil {
		return nil, err
	}

	return principalTagAttributeMapToHTTP(result), nil
}

// SetPrincipalTagAttributeMap sets the principal tag attribute map for an identity provider.
func (s *CognitoIdentityService) SetPrincipalTagAttributeMap(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.setPrincipalTagAttributeMapCore(reqCtx, SetPrincipalTagAttributeMapInput{
		IdentityPoolID:       req.GetParam("IdentityPoolId"),
		IdentityProviderName: req.GetParam("IdentityProviderName"),
		PrincipalTags:        parseMapParam(req, "PrincipalTags"),
		UseDefaults:          getBoolParam(req, "UseDefaults"),
	})
	if err != nil {
		return nil, err
	}

	return principalTagAttributeMapToHTTP(result), nil
}

// principalTagAttributeMapToHTTP serialises a PrincipalTagAttributeMapResult
// into the member format shared by the Get/SetPrincipalTagAttributeMap
// responses.
func principalTagAttributeMapToHTTP(r *PrincipalTagAttributeMapResult) map[string]interface{} {
	return map[string]interface{}{
		"IdentityPoolId":       r.IdentityPoolID,
		"IdentityProviderName": r.IdentityProviderName,
		"PrincipalTags":        r.PrincipalTags,
		"UseDefaults":          r.UseDefaults,
	}
}
