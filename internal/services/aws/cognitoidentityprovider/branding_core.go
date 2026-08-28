package cognitoidentityprovider

import (
	"encoding/base64"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"github.com/google/uuid"
)

// CreateManagedLoginBrandingInput carries the wire parameters of
// CreateManagedLoginBranding. Params holds the raw request parameter map;
// the UseCognitoProvidedValues flag, the nested Settings structure and the
// Assets list are read from it inside the Core.
type CreateManagedLoginBrandingInput struct {
	UserPoolID string
	ClientID   string
	Params     map[string]interface{}
}

// DescribeManagedLoginBrandingInput carries the wire parameters of
// DescribeManagedLoginBranding.
type DescribeManagedLoginBrandingInput struct {
	UserPoolID             string
	ManagedLoginBrandingID string
}

// DescribeManagedLoginBrandingByClientInput carries the wire parameters of
// DescribeManagedLoginBrandingByClient.
type DescribeManagedLoginBrandingByClientInput struct {
	UserPoolID string
	ClientID   string
}

// UpdateManagedLoginBrandingInput carries the wire parameters of
// UpdateManagedLoginBranding. Params holds the raw request parameter map for
// the update members.
type UpdateManagedLoginBrandingInput struct {
	UserPoolID             string
	ManagedLoginBrandingID string
	Params                 map[string]interface{}
}

// DeleteManagedLoginBrandingInput carries the wire parameters of
// DeleteManagedLoginBranding.
type DeleteManagedLoginBrandingInput struct {
	UserPoolID             string
	ManagedLoginBrandingID string
}

// createManagedLoginBrandingCore creates a managed login branding
// configuration.
func (s *CognitoService) createManagedLoginBrandingCore(reqCtx *request.RequestContext, in CreateManagedLoginBrandingInput) (interface{}, error) {
	// CreateManagedLoginBrandingRequest marks UserPoolId and ClientId
	// required.
	if in.UserPoolID == "" || in.ClientID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(in.UserPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	// An app client can only ever have one assigned branding style.
	if _, err := store.GetManagedLoginBrandingByClient(in.UserPoolID, in.ClientID); err == nil {
		return nil, ErrManagedLoginBrandingExists
	}

	// ManagedLoginBrandingIdType is a v4 UUID pattern; the minted
	// identifier must be a value AWS itself could have issued.
	brandingID := uuid.NewString()
	b := &cognitostore.ManagedLoginBranding{
		ManagedLoginBrandingId: brandingID,
		UserPoolID:             in.UserPoolID,
		ClientID:               in.ClientID,
	}
	if v, ok := in.Params["UseCognitoProvidedValues"].(bool); ok {
		b.UseCognitoProvidedValues = v
	}
	if settings, ok := in.Params["Settings"]; ok {
		if m, ok := settings.(map[string]interface{}); ok {
			b.Settings = m
		}
	}
	if err := parseBrandingAssets(in.Params, b); err != nil {
		return nil, err
	}

	if err := store.SaveManagedLoginBranding(b); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{"ManagedLoginBranding": formatManagedLoginBranding(b)}, nil
}

// describeManagedLoginBrandingCore describes a managed login branding
// configuration.
func (s *CognitoService) describeManagedLoginBrandingCore(reqCtx *request.RequestContext, in DescribeManagedLoginBrandingInput) (interface{}, error) {
	if in.UserPoolID == "" || in.ManagedLoginBrandingID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	b, err := store.GetManagedLoginBranding(in.UserPoolID, in.ManagedLoginBrandingID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{"ManagedLoginBranding": formatManagedLoginBranding(b)}, nil
}

// describeManagedLoginBrandingByClientCore describes branding by client ID.
func (s *CognitoService) describeManagedLoginBrandingByClientCore(reqCtx *request.RequestContext, in DescribeManagedLoginBrandingByClientInput) (interface{}, error) {
	if in.UserPoolID == "" || in.ClientID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	b, err := store.GetManagedLoginBrandingByClient(in.UserPoolID, in.ClientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{"ManagedLoginBranding": formatManagedLoginBranding(b)}, nil
}

// updateManagedLoginBrandingCore updates a managed login branding
// configuration.
func (s *CognitoService) updateManagedLoginBrandingCore(reqCtx *request.RequestContext, in UpdateManagedLoginBrandingInput) (interface{}, error) {
	if in.UserPoolID == "" || in.ManagedLoginBrandingID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	b, err := store.GetManagedLoginBranding(in.UserPoolID, in.ManagedLoginBrandingID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if v, ok := in.Params["UseCognitoProvidedValues"].(bool); ok {
		b.UseCognitoProvidedValues = v
	}
	if settings, ok := in.Params["Settings"]; ok {
		if m, ok := settings.(map[string]interface{}); ok {
			b.Settings = m
		}
	}
	if err := parseBrandingAssets(in.Params, b); err != nil {
		return nil, err
	}

	if err := store.SaveManagedLoginBranding(b); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{"ManagedLoginBranding": formatManagedLoginBranding(b)}, nil
}

// deleteManagedLoginBrandingCore deletes a managed login branding
// configuration.
func (s *CognitoService) deleteManagedLoginBrandingCore(reqCtx *request.RequestContext, in DeleteManagedLoginBrandingInput) (interface{}, error) {
	if in.UserPoolID == "" || in.ManagedLoginBrandingID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteManagedLoginBranding(in.UserPoolID, in.ManagedLoginBrandingID); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

var validAssetCategories = map[string]bool{
	"AUTH_APP_GRAPHIC": true, "EMAIL_GRAPHIC": true, "FAVICON_ICO": true,
	"FAVICON_SVG": true, "FORM_BACKGROUND": true, "FORM_LOGO": true,
	"IDP_BUTTON_ICON": true, "PAGE_BACKGROUND": true, "PAGE_FOOTER_BACKGROUND": true,
	"PAGE_FOOTER_LOGO": true, "PAGE_HEADER_BACKGROUND": true, "PAGE_HEADER_LOGO": true,
	"PASSKEY_GRAPHIC": true, "PASSWORD_GRAPHIC": true, "SMS_GRAPHIC": true,
}

var validAssetExtensions = map[string]bool{
	"ICO": true, "JPEG": true, "PNG": true, "SVG": true, "WEBP": true,
}

var validColorModes = map[string]bool{
	"LIGHT": true, "DARK": true, "DYNAMIC": true,
}

var assetMagicBytes = map[string][]byte{
	"PNG":  {0x89, 0x50, 0x4E, 0x47},
	"JPEG": {0xFF, 0xD8, 0xFF},
	"ICO":  {0x00, 0x00, 0x01, 0x00},
	"GIF":  {0x47, 0x49, 0x46, 0x38},
}

func parseBrandingAssets(params map[string]interface{}, b *cognitostore.ManagedLoginBranding) error {
	if rawAssets, ok := params["Assets"].([]interface{}); ok {
		b.Assets = nil
		for _, a := range rawAssets {
			m, ok := a.(map[string]interface{})
			if !ok {
				continue
			}
			category := getStringParam(m, "Category")
			extension := getStringParam(m, "Extension")
			colorMode := getStringParam(m, "Color")
			// Validate enum values
			if category != "" && !validAssetCategories[category] {
				return ErrInvalidParameter
			}
			if extension != "" && !validAssetExtensions[extension] {
				return ErrInvalidParameter
			}
			if colorMode != "" && !validColorModes[colorMode] {
				return ErrInvalidParameter
			}
			// Validate Bytes magic bytes when extension is known
			bytesVal := getStringParam(m, "Bytes")
			if bytesVal != "" && extension != "" {
				if magic, ok := assetMagicBytes[extension]; ok {
					decoded, err := base64.StdEncoding.DecodeString(bytesVal)
					if err == nil && len(decoded) >= len(magic) {
						for i, b := range magic {
							if decoded[i] != b {
								return ErrInvalidParameter
							}
						}
					}
				}
			}
			asset := cognitostore.BrandingAsset{
				Category:  category,
				Color:     colorMode,
				Extension: extension,
				Bytes:     bytesVal,
			}
			b.Assets = append(b.Assets, asset)
		}
	}
	return nil
}
