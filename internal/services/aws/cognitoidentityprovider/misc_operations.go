package cognitoidentityprovider

import (
	"context"
	"encoding/base64"
	"sync"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// ===================== UI Customization =====================

// GetUICustomization retrieves the UI customisation for a user pool/client.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetUICustomization.html
func (s *CognitoService) GetUICustomization(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	clientID := req.GetParam("ClientId")
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	ui, err := store.GetUICustomization(userPoolID, clientID)
	if err != nil {
		ui = &cognitostore.UICustomization{UserPoolID: userPoolID, ClientID: clientID}
	}

	return map[string]interface{}{
		"UICustomization": formatUICustomization(ui),
	}, nil
}

// SetUICustomization sets the UI customisation for a user pool/client.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_SetUICustomization.html
func (s *CognitoService) SetUICustomization(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	clientID := req.GetParam("ClientId")
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	ui, err := store.GetUICustomization(userPoolID, clientID)
	if err != nil {
		ui = &cognitostore.UICustomization{UserPoolID: userPoolID, ClientID: clientID}
	}

	if css := req.GetParam("CSS"); css != "" {
		ui.CSS = css
	}
	if imageFile, ok := req.Parameters["ImageFile"].(string); ok && imageFile != "" {
		decoded, err := base64.StdEncoding.DecodeString(imageFile)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		if !validateImageFileSize(decoded) {
			return nil, ErrInvalidParameter
		}
		ui.ImageFile = decoded
	}

	if err := store.SaveUICustomization(ui); err != nil {
		return nil, ErrInternalError
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
	userPoolID := req.GetParam("UserPoolId")
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	userRaw, ok := req.Parameters["User"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	providerName := getStringParam(userRaw, "ProviderName")
	providerAttrValue := getStringParam(userRaw, "ProviderAttributeValue")
	if providerName == "" || providerAttrValue == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByProvider(userPoolID, providerName, providerAttrValue)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.ProviderName = ""
	user.ProviderAttributeName = ""
	user.ProviderAttributeValue = ""
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// AdminLinkProviderForUser links a federated provider to an existing user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminLinkProviderForUser.html
func (s *CognitoService) AdminLinkProviderForUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	destRaw, ok := req.Parameters["DestinationUser"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	srcRaw, ok := req.Parameters["SourceUser"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	destUsername := getStringParam(destRaw, "ProviderAttributeValue")
	destProviderName := getStringParam(destRaw, "ProviderName")
	if destProviderName == "" {
		return nil, ErrInvalidParameter
	}
	srcProviderName := getStringParam(srcRaw, "ProviderName")
	srcProviderAttrName := getStringParam(srcRaw, "ProviderAttributeName")
	srcProviderAttrValue := getStringParam(srcRaw, "ProviderAttributeValue")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(userPoolID, destUsername)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.ProviderName = srcProviderName
	user.ProviderAttributeName = srcProviderAttrName
	user.ProviderAttributeValue = srcProviderAttrValue
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// ===================== Misc Small Operations =====================

// AddCustomAttributes adds custom schema attributes to a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AddCustomAttributes.html
func (s *CognitoService) AddCustomAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if attrs, ok := req.Parameters["CustomAttributes"].([]interface{}); ok {
		for _, a := range attrs {
			m, ok := a.(map[string]interface{})
			if !ok {
				continue
			}
			attrName := getStringParam(m, "Name")
			attrType := getStringParam(m, "AttributeDataType")
			if attrName == "" || attrType == "" {
				return nil, ErrInvalidParameter
			}
			if err := validateCustomAttributeName(attrName); err != nil {
				return nil, err
			}
			if !validateAttributeDataType(attrType) {
				return nil, ErrInvalidParameter
			}
			for _, existing := range pool.SchemaAttributes {
				if existing.Name == attrName {
					return nil, ErrInvalidParameter
				}
			}
			newAttr := cognitostore.SchemaAttributeType{
				Name:              attrName,
				AttributeDataType: attrType,
			}
			if dev, ok := m["DeveloperOnlyAttribute"].(bool); ok {
				newAttr.DeveloperOnlyAttribute = dev
			}
			if mut, ok := m["Mutable"].(bool); ok {
				newAttr.Mutable = mut
			}
			if reqVal, ok := m["Required"].(bool); ok {
				newAttr.Required = reqVal
			}
			if nac, ok := m["NumberAttributeConstraints"].(map[string]interface{}); ok {
				nc := &cognitostore.NumberAttributeConstraints{}
				if v, ok := nac["MinValue"].(string); ok {
					nc.MinValue = v
				}
				if v, ok := nac["MaxValue"].(string); ok {
					nc.MaxValue = v
				}
				newAttr.NumberAttributeConstraints = nc
			}
			if sac, ok := m["StringAttributeConstraints"].(map[string]interface{}); ok {
				sc := &cognitostore.StringAttributeConstraints{}
				if v, ok := sac["MinLength"].(string); ok {
					sc.MinLength = v
				}
				if v, ok := sac["MaxLength"].(string); ok {
					sc.MaxLength = v
				}
				newAttr.StringAttributeConstraints = sc
			}
			pool.SchemaAttributes = append(pool.SchemaAttributes, newAttr)
		}
	}

	if err := store.UpdateUserPool(pool); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// GetIdentityProviderByIdentifier retrieves an IdP by its identifier (domain or DNS name).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetIdentityProviderByIdentifier.html
func (s *CognitoService) GetIdentityProviderByIdentifier(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	idpIdentifier := req.GetParam("IdpIdentifier")
	if userPoolID == "" || idpIdentifier == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	providers, err := store.ListIdentityProviders(userPoolID)
	if err != nil {
		return nil, ErrInternalError
	}

	for _, idp := range providers {
		for _, ident := range idp.IdpIdentifiers {
			if ident == idpIdentifier {
				return map[string]interface{}{"IdentityProvider": formatIdentityProvider(idp)}, nil
			}
		}
	}

	return nil, ErrResourceNotFound
}

// GetSigningCertificate returns the user pool's signing certificate.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetSigningCertificate.html
func (s *CognitoService) GetSigningCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{
		"Certificate": pool.JwtPublicKey,
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
