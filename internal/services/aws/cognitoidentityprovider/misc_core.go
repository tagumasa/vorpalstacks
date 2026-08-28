package cognitoidentityprovider

import (
	"encoding/base64"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// Core functions for the remaining single-operation families: UI
// customisation, provider linking, custom schema attributes, identity
// provider lookup by identifier and the signing certificate. The handlers
// extract the wire members; validation and store access live here.

// SetUICustomizationInput carries the wire members of SetUICustomization.
// ImageFileProvided distinguishes an absent image from an empty one.
type SetUICustomizationInput struct {
	Region            string
	UserPoolID        string
	ClientID          string
	CSS               string
	ImageFile         string
	ImageFileProvided bool
}

// getUICustomizationCore loads the UI customisation for a pool/client pair,
// falling back to a fresh record when none is stored.
func (s *CognitoService) getUICustomizationCore(region, userPoolID, clientID string) (*cognitostore.UICustomization, error) {
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
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
	return ui, nil
}

// setUICustomizationCore applies the CSS and image members onto the stored
// UI customisation. The image payload is base64-decoded and size-validated.
func (s *CognitoService) setUICustomizationCore(in SetUICustomizationInput) (*cognitostore.UICustomization, error) {
	if in.UserPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(in.Region)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(in.UserPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	ui, err := store.GetUICustomization(in.UserPoolID, in.ClientID)
	if err != nil {
		ui = &cognitostore.UICustomization{UserPoolID: in.UserPoolID, ClientID: in.ClientID}
	}

	if in.CSS != "" {
		ui.CSS = in.CSS
	}
	if in.ImageFileProvided {
		decoded, err := base64.StdEncoding.DecodeString(in.ImageFile)
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

	return ui, nil
}

// adminDisableProviderForUserCore detaches the federated provider identity
// from the referenced user.
func (s *CognitoService) adminDisableProviderForUserCore(region, userPoolID string, user map[string]interface{}) error {
	if userPoolID == "" {
		return ErrInvalidParameter
	}
	if user == nil {
		return ErrInvalidParameter
	}

	providerName := getStringParam(user, "ProviderName")
	providerAttrValue := getStringParam(user, "ProviderAttributeValue")
	if providerName == "" || providerAttrValue == "" {
		return ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	u, err := store.GetUserByProvider(userPoolID, providerName, providerAttrValue)
	if err != nil {
		return ErrUserNotFound
	}

	u.ProviderName = ""
	u.ProviderAttributeName = ""
	u.ProviderAttributeValue = ""
	if err := store.UpdateUser(u); err != nil {
		return ErrInternalError
	}

	return nil
}

// adminLinkProviderForUserCore links a federated source identity onto the
// local destination user.
func (s *CognitoService) adminLinkProviderForUserCore(region, userPoolID string, destinationUser, sourceUser map[string]interface{}) error {
	if userPoolID == "" {
		return ErrInvalidParameter
	}
	if destinationUser == nil {
		return ErrInvalidParameter
	}
	if sourceUser == nil {
		return ErrInvalidParameter
	}

	destUsername := getStringParam(destinationUser, "ProviderAttributeValue")
	destProviderName := getStringParam(destinationUser, "ProviderName")
	if destProviderName == "" {
		return ErrInvalidParameter
	}
	srcProviderName := getStringParam(sourceUser, "ProviderName")
	srcProviderAttrName := getStringParam(sourceUser, "ProviderAttributeName")
	srcProviderAttrValue := getStringParam(sourceUser, "ProviderAttributeValue")

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	u, err := store.GetUser(userPoolID, destUsername)
	if err != nil {
		return ErrUserNotFound
	}

	u.ProviderName = srcProviderName
	u.ProviderAttributeName = srcProviderAttrName
	u.ProviderAttributeValue = srcProviderAttrValue
	if err := store.UpdateUser(u); err != nil {
		return ErrInternalError
	}

	return nil
}

// addCustomAttributesCore appends the requested custom attribute schemas to
// the pool, rejecting invalid names, non-member data types and duplicates of
// existing schema attributes.
func (s *CognitoService) addCustomAttributesCore(region, userPoolID string, customAttributes []interface{}) error {
	if userPoolID == "" {
		return ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}

	pool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return ErrResourceNotFound
	}

	for _, a := range customAttributes {
		m, ok := a.(map[string]interface{})
		if !ok {
			continue
		}
		attrName := getStringParam(m, "Name")
		attrType := getStringParam(m, "AttributeDataType")
		if attrName == "" || attrType == "" {
			return ErrInvalidParameter
		}
		if err := validateCustomAttributeName(attrName); err != nil {
			return err
		}
		if !validateAttributeDataType(attrType) {
			return ErrInvalidParameter
		}
		for _, existing := range pool.SchemaAttributes {
			if existing.Name == attrName {
				return ErrInvalidParameter
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

	if err := store.UpdateUserPool(pool); err != nil {
		return ErrInternalError
	}

	return nil
}

// getIdentityProviderByIdentifierCore resolves the identity provider that
// declares the given IdpIdentifier.
func (s *CognitoService) getIdentityProviderByIdentifierCore(region, userPoolID, idpIdentifier string) (*cognitostore.IdentityProvider, error) {
	if userPoolID == "" || idpIdentifier == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
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
				return idp, nil
			}
		}
	}

	return nil, ErrResourceNotFound
}

// getSigningCertificateCore returns the pool's JWT signing key.
func (s *CognitoService) getSigningCertificateCore(region, userPoolID string) (string, error) {
	if userPoolID == "" {
		return "", ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return "", err
	}

	pool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return "", ErrResourceNotFound
	}

	return pool.JwtPublicKey, nil
}
