package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// Handlers for the UI customisation family; validation and store access
// live in ui_customization_core.go.

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
