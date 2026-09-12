package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// Handlers for the managed-login branding family; validation and store
// access live in branding_core.go.

// CreateManagedLoginBranding creates a managed login branding configuration.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateManagedLoginBranding.html
func (s *CognitoService) CreateManagedLoginBranding(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.createManagedLoginBrandingCore(reqCtx, CreateManagedLoginBrandingInput{
		UserPoolID: req.GetParam("UserPoolId"),
		ClientID:   req.GetParam("ClientId"),
		Params:     req.Parameters,
	})
}

// DescribeManagedLoginBranding describes a managed login branding configuration.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DescribeManagedLoginBranding.html
func (s *CognitoService) DescribeManagedLoginBranding(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.describeManagedLoginBrandingCore(reqCtx, DescribeManagedLoginBrandingInput{
		UserPoolID:             req.GetParam("UserPoolId"),
		ManagedLoginBrandingID: req.GetParam("ManagedLoginBrandingId"),
	})
}

// DescribeManagedLoginBrandingByClient describes branding by client ID.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DescribeManagedLoginBrandingByClient.html
func (s *CognitoService) DescribeManagedLoginBrandingByClient(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.describeManagedLoginBrandingByClientCore(reqCtx, DescribeManagedLoginBrandingByClientInput{
		UserPoolID: req.GetParam("UserPoolId"),
		ClientID:   req.GetParam("ClientId"),
	})
}

// UpdateManagedLoginBranding updates a managed login branding configuration.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateManagedLoginBranding.html
func (s *CognitoService) UpdateManagedLoginBranding(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateManagedLoginBrandingCore(reqCtx, UpdateManagedLoginBrandingInput{
		UserPoolID:             req.GetParam("UserPoolId"),
		ManagedLoginBrandingID: req.GetParam("ManagedLoginBrandingId"),
		Params:                 req.Parameters,
	})
}

// DeleteManagedLoginBranding deletes a managed login branding configuration.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteManagedLoginBranding.html
func (s *CognitoService) DeleteManagedLoginBranding(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.deleteManagedLoginBrandingCore(reqCtx, DeleteManagedLoginBrandingInput{
		UserPoolID:             req.GetParam("UserPoolId"),
		ManagedLoginBrandingID: req.GetParam("ManagedLoginBrandingId"),
	})
}

func formatManagedLoginBranding(b *cognitostore.ManagedLoginBranding) map[string]interface{} {
	result := map[string]interface{}{
		"ManagedLoginBrandingId":   b.ManagedLoginBrandingId,
		"UserPoolId":               b.UserPoolID,
		"UseCognitoProvidedValues": b.UseCognitoProvidedValues,
	}
	if b.ClientID != "" {
		result["ClientId"] = b.ClientID
	}
	if b.Settings != nil {
		result["Settings"] = b.Settings
	}
	if len(b.Assets) > 0 {
		assets := make([]map[string]interface{}, 0, len(b.Assets))
		for _, a := range b.Assets {
			assets = append(assets, map[string]interface{}{
				"Category":  a.Category,
				"Color":     a.Color,
				"Extension": a.Extension,
				"Bytes":     a.Bytes,
			})
		}
		result["Assets"] = assets
	}
	if !b.CreationDate.IsZero() {
		result["CreationDate"] = b.CreationDate.Unix()
	}
	if !b.LastModifiedDate.IsZero() {
		result["LastModifiedDate"] = b.LastModifiedDate.Unix()
	}
	return result
}
