package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// Handlers for the terms-document family; validation and store access live
// in terms_core.go.

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

// DescribeTermsByClient describes the terms document an app client holds
// under a terms name.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DescribeTermsByClient.html
func (s *CognitoService) DescribeTermsByClient(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.describeTermsByClientCore(reqCtx, DescribeTermsByClientInput{
		UserPoolID: req.GetParam("UserPoolId"),
		ClientID:   req.GetParam("ClientId"),
		TermsName:  req.GetParam("TermsName"),
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
