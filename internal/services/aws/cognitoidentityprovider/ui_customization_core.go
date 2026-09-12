package cognitoidentityprovider

import (
	"encoding/base64"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// Core functions for the UI customisation family. The handlers extract the
// wire members; validation and store access live here.

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
		return nil, ErrInternalError
	}
	if ui == nil {
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
		return nil, ErrInternalError
	}
	if ui == nil {
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
