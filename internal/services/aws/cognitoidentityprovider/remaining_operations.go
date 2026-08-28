package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// ===================== WebAuthn =====================

// StartWebAuthnRegistration starts a WebAuthn credential registration flow.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_StartWebAuthnRegistration.html
func (s *CognitoService) StartWebAuthnRegistration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.startWebAuthnRegistrationCore(reqCtx, StartWebAuthnRegistrationInput{
		AccessToken: getAccessToken(req),
	})
}

// CompleteWebAuthnRegistration completes a WebAuthn credential registration.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CompleteWebAuthnRegistration.html
func (s *CognitoService) CompleteWebAuthnRegistration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.completeWebAuthnRegistrationCore(reqCtx, CompleteWebAuthnRegistrationInput{
		AccessToken: getAccessToken(req),
		Params:      req.Parameters,
	})
}

// ListWebAuthnCredentials lists registered WebAuthn credentials.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListWebAuthnCredentials.html
func (s *CognitoService) ListWebAuthnCredentials(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.listWebAuthnCredentialsCore(reqCtx, ListWebAuthnCredentialsInput{
		AccessToken: getAccessToken(req),
		NextToken:   req.GetParam("NextToken"),
		Params:      req.Parameters,
	})
}

// DeleteWebAuthnCredential deletes a WebAuthn credential.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteWebAuthnCredential.html
func (s *CognitoService) DeleteWebAuthnCredential(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.deleteWebAuthnCredentialCore(reqCtx, DeleteWebAuthnCredentialInput{
		AccessToken:  getAccessToken(req),
		CredentialID: req.GetParam("CredentialId"),
	})
}

// ===================== Managed Login Branding =====================

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

// ===================== Terms =====================

// CreateTerms creates a terms document.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateTerms.html
func (s *CognitoService) CreateTerms(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.createTermsCore(reqCtx, CreateTermsInput{
		UserPoolID:  req.GetParam("UserPoolId"),
		ClientID:    req.GetParam("ClientId"),
		TermsName:   req.GetParam("TermsName"),
		TermsSource: req.GetParam("TermsSource"),
		Enforcement: req.GetParam("Enforcement"),
		Params:      req.Parameters,
	})
}

// DescribeTerms describes a terms document.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DescribeTerms.html
func (s *CognitoService) DescribeTerms(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.describeTermsCore(reqCtx, DescribeTermsInput{
		UserPoolID: req.GetParam("UserPoolId"),
		TermsID:    req.GetParam("TermsId"),
	})
}

// ListTerms lists terms documents for a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListTerms.html
func (s *CognitoService) ListTerms(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.listTermsCore(reqCtx, ListTermsInput{
		UserPoolID: req.GetParam("UserPoolId"),
		NextToken:  req.GetParam("NextToken"),
		Params:     req.Parameters,
	})
}

// UpdateTerms updates a terms document.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateTerms.html
func (s *CognitoService) UpdateTerms(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateTermsCore(reqCtx, UpdateTermsInput{
		UserPoolID:  req.GetParam("UserPoolId"),
		TermsID:     req.GetParam("TermsId"),
		TermsName:   req.GetParam("TermsName"),
		TermsSource: req.GetParam("TermsSource"),
		Enforcement: req.GetParam("Enforcement"),
		Params:      req.Parameters,
	})
}

// DeleteTerms deletes a terms document.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteTerms.html
func (s *CognitoService) DeleteTerms(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.deleteTermsCore(reqCtx, DeleteTermsInput{
		UserPoolID: req.GetParam("UserPoolId"),
		TermsID:    req.GetParam("TermsId"),
	})
}

func formatTerms(t *cognitostore.Terms) map[string]interface{} {
	result := map[string]interface{}{
		"TermsId":    t.TermsID,
		"UserPoolId": t.UserPoolID,
		"TermsName":  t.TermsName,
	}
	if t.ClientID != "" {
		result["ClientId"] = t.ClientID
	}
	if t.TermsSource != "" {
		result["TermsSource"] = t.TermsSource
	}
	if t.Enforcement != "" {
		result["Enforcement"] = t.Enforcement
	}
	if t.Links != nil {
		result["Links"] = t.Links
	}
	if !t.CreationDate.IsZero() {
		result["CreationDate"] = t.CreationDate.Unix()
	}
	if !t.LastModifiedDate.IsZero() {
		result["LastModifiedDate"] = t.LastModifiedDate.Unix()
	}
	return result
}

// formatTermsDescription projects a terms document onto the smaller
// TermsDescriptionType shape that ListTerms returns: five members, without
// the UserPoolId/ClientId/TermsSource/Links that TermsType carries.
func formatTermsDescription(t *cognitostore.Terms) map[string]interface{} {
	result := map[string]interface{}{
		"TermsId":   t.TermsID,
		"TermsName": t.TermsName,
	}
	if t.Enforcement != "" {
		result["Enforcement"] = t.Enforcement
	}
	if !t.CreationDate.IsZero() {
		result["CreationDate"] = t.CreationDate.Unix()
	}
	if !t.LastModifiedDate.IsZero() {
		result["LastModifiedDate"] = t.LastModifiedDate.Unix()
	}
	return result
}

// ===================== Replicas =====================

// CreateUserPoolReplica creates a cross-region replica of a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateUserPoolReplica.html
func (s *CognitoService) CreateUserPoolReplica(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.createUserPoolReplicaCore(reqCtx, CreateUserPoolReplicaInput{
		UserPoolID: req.GetParam("UserPoolId"),
		RegionName: req.GetParam("RegionName"),
		Params:     req.Parameters,
	})
}

// ListUserPoolReplicas lists replicas for a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListUserPoolReplicas.html
func (s *CognitoService) ListUserPoolReplicas(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.listUserPoolReplicasCore(reqCtx, ListUserPoolReplicasInput{
		UserPoolID: req.GetParam("UserPoolId"),
		NextToken:  req.GetParam("NextToken"),
	})
}

// DeleteUserPoolReplica deletes a cross-region replica.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteUserPoolReplica.html
func (s *CognitoService) DeleteUserPoolReplica(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.deleteUserPoolReplicaCore(reqCtx, DeleteUserPoolReplicaInput{
		UserPoolID: req.GetParam("UserPoolId"),
		RegionName: req.GetParam("RegionName"),
	})
}

// UpdateUserPoolReplica updates a cross-region replica.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateUserPoolReplica.html
func (s *CognitoService) UpdateUserPoolReplica(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateUserPoolReplicaCore(reqCtx, UpdateUserPoolReplicaInput{
		UserPoolID: req.GetParam("UserPoolId"),
		RegionName: req.GetParam("RegionName"),
		Status:     req.GetParam("Status"),
	})
}

func formatUserPoolReplica(r *cognitostore.UserPoolReplica) map[string]interface{} {
	result := map[string]interface{}{
		"RegionName": r.RegionName,
		"Status":     r.Status,
		"Role":       r.Role,
	}
	if r.UserPoolArn != "" {
		result["UserPoolArn"] = r.UserPoolArn
	}
	return result
}
