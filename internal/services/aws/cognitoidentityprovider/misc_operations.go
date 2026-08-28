package cognitoidentityprovider

import (
	"context"
	"sync"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// ===================== UI Customization =====================

// GetUICustomization retrieves the UI customisation for a user pool/client.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetUICustomization.html
func (s *CognitoService) GetUICustomization(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	ui, err := s.getUICustomizationCore(reqCtx.GetRegion(), req.GetParam("UserPoolId"), req.GetParam("ClientId"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"UICustomization": formatUICustomization(ui),
	}, nil
}

// SetUICustomization sets the UI customisation for a user pool/client.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_SetUICustomization.html
func (s *CognitoService) SetUICustomization(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := SetUICustomizationInput{
		Region:     reqCtx.GetRegion(),
		UserPoolID: req.GetParam("UserPoolId"),
		ClientID:   req.GetParam("ClientId"),
		CSS:        req.GetParam("CSS"),
	}
	if v, ok := req.Parameters["ImageFile"].(string); ok && v != "" {
		in.ImageFile = v
		in.ImageFileProvided = true
	}

	ui, err := s.setUICustomizationCore(in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"UICustomization": formatUICustomization(ui),
	}, nil
}

func formatUICustomization(ui *cognitostore.UICustomization) map[string]interface{} {
	result := map[string]interface{}{
		"UserPoolId": ui.UserPoolID,
		"CSS":        ui.CSS,
		"CSSVersion": ui.CSSVersion,
	}
	if ui.ClientID != "" {
		result["ClientId"] = ui.ClientID
	}
	if !ui.CreationDate.IsZero() {
		result["CreationDate"] = ui.CreationDate.Unix()
	}
	if !ui.LastModifiedDate.IsZero() {
		result["LastModifiedDate"] = ui.LastModifiedDate.Unix()
	}
	return result
}

// ===================== Provider User Linking =====================

// AdminDisableProviderForUser disables a federated provider for a user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminDisableProviderForUser.html
func (s *CognitoService) AdminDisableProviderForUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	user, _ := req.Parameters["User"].(map[string]interface{})

	if err := s.adminDisableProviderForUserCore(reqCtx.GetRegion(), req.GetParam("UserPoolId"), user); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// AdminLinkProviderForUser links a federated provider to an existing user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminLinkProviderForUser.html
func (s *CognitoService) AdminLinkProviderForUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	destinationUser, _ := req.Parameters["DestinationUser"].(map[string]interface{})
	sourceUser, _ := req.Parameters["SourceUser"].(map[string]interface{})

	if err := s.adminLinkProviderForUserCore(reqCtx.GetRegion(), req.GetParam("UserPoolId"), destinationUser, sourceUser); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ===================== Misc Small Operations =====================

// AddCustomAttributes adds custom schema attributes to a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AddCustomAttributes.html
func (s *CognitoService) AddCustomAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var attrs []interface{}
	if v, ok := req.Parameters["CustomAttributes"].([]interface{}); ok {
		attrs = v
	}

	if err := s.addCustomAttributesCore(reqCtx.GetRegion(), req.GetParam("UserPoolId"), attrs); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetIdentityProviderByIdentifier retrieves an IdP by its identifier (domain or DNS name).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetIdentityProviderByIdentifier.html
func (s *CognitoService) GetIdentityProviderByIdentifier(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	idp, err := s.getIdentityProviderByIdentifierCore(reqCtx.GetRegion(), req.GetParam("UserPoolId"), req.GetParam("IdpIdentifier"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"IdentityProvider": formatIdentityProvider(idp)}, nil
}

// GetSigningCertificate returns the user pool's signing certificate.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetSigningCertificate.html
func (s *CognitoService) GetSigningCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	certificate, err := s.getSigningCertificateCore(reqCtx.GetRegion(), req.GetParam("UserPoolId"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Certificate": certificate,
	}, nil
}

// ===================== Provisioned Limits =====================

// provisionedLimits is a simple in-memory store since this is rarely used.
var (
	provisionedLimits   = make(map[string]int)
	provisionedLimitsMu sync.RWMutex
)

func provisionedLimitKey(class, resourceType string) string {
	return class + ":" + resourceType
}

// GetProvisionedLimit retrieves a provisioned limit.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetProvisionedLimit.html
func (s *CognitoService) GetProvisionedLimit(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limitDef, ok := req.Parameters["LimitDefinition"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	limitClass := getStringParam(limitDef, "LimitClass")
	attrs, _ := limitDef["Attributes"].(map[string]interface{})
	resourceType := getStringParam(attrs, "ResourceType")

	key := provisionedLimitKey(limitClass, resourceType)
	provisionedLimitsMu.RLock()
	limit := provisionedLimits[key]
	provisionedLimitsMu.RUnlock()
	if limit == 0 {
		limit = defaultProvisionedLimit
	}

	return map[string]interface{}{
		"Limit": map[string]interface{}{
			"LimitDefinition": limitDef,
			"CurrentValue":    limit,
			"EffectiveLimit":  limit,
		},
	}, nil
}

// UpdateProvisionedLimit updates a provisioned limit.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateProvisionedLimit.html
func (s *CognitoService) UpdateProvisionedLimit(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limitDef, ok := req.Parameters["LimitDefinition"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	limitClass := getStringParam(limitDef, "LimitClass")
	attrs, _ := limitDef["Attributes"].(map[string]interface{})
	resourceType := getStringParam(attrs, "ResourceType")

	requestedValue := request.GetIntParam(req.Parameters, "RequestedLimitValue")
	if requestedValue <= 0 {
		return nil, ErrInvalidParameter
	}

	key := provisionedLimitKey(limitClass, resourceType)
	provisionedLimitsMu.Lock()
	provisionedLimits[key] = requestedValue
	provisionedLimitsMu.Unlock()

	return map[string]interface{}{
		"Limit": map[string]interface{}{
			"LimitDefinition": limitDef,
			"CurrentValue":    requestedValue,
			"EffectiveLimit":  requestedValue,
		},
	}, nil
}
