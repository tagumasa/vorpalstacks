package iam

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateServiceSpecificCredential generates a new service-specific credential for the specified user and service.
func (s *IAMService) CreateServiceSpecificCredential(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := &CreateServiceSpecificCredentialInput{
		UserName:    request.GetStringParam(req.Parameters, "UserName"),
		ServiceName: request.GetStringParam(req.Parameters, "ServiceName"),
	}
	if _, ok := req.Parameters["CredentialAgeDays"]; ok {
		credentialAgeDays := request.GetIntParam(req.Parameters, "CredentialAgeDays")
		input.CredentialAgeDays = &credentialAgeDays
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cred, err := s.createServiceSpecificCredentialCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ServiceSpecificCredential": s.serviceSpecificCredentialToResponse(cred, true),
	}, nil
}

// DeleteServiceSpecificCredential deletes the specified service-specific credential.
func (s *IAMService) DeleteServiceSpecificCredential(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	credentialId := request.GetStringParam(req.Parameters, "ServiceSpecificCredentialId")
	userName := request.GetStringParam(req.Parameters, "UserName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteServiceSpecificCredentialCore(store, credentialId, userName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListServiceSpecificCredentials lists all service-specific credentials for the specified user.
func (s *IAMService) ListServiceSpecificCredentials(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	userName, err := resolveUserName(reqCtx, userName)
	if err != nil {
		return nil, err
	}
	serviceName := request.GetStringParam(req.Parameters, "ServiceName")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listServiceSpecificCredentialsCore(store, userName, serviceName, marker, maxItems)
	if err != nil {
		return nil, err
	}

	credList := make([]interface{}, len(result.Credentials))
	for i, cred := range result.Credentials {
		credList[i] = s.serviceSpecificCredentialToResponse(cred, false)
	}

	resp := map[string]interface{}{
		"ServiceSpecificCredentials": credList,
		"IsTruncated":                result.IsTruncated,
	}
	if result.NextMarker != "" {
		resp["Marker"] = result.NextMarker
	}
	return resp, nil
}

// ResetServiceSpecificCredential resets the password for a service-specific credential.
func (s *IAMService) ResetServiceSpecificCredential(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	credentialId := request.GetStringParam(req.Parameters, "ServiceSpecificCredentialId")
	userName := request.GetStringParam(req.Parameters, "UserName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cred, err := s.resetServiceSpecificCredentialCore(store, credentialId, userName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ServiceSpecificCredential": s.serviceSpecificCredentialToResponse(cred, true),
	}, nil
}

// UpdateServiceSpecificCredential sets the status of a service-specific credential to Active or Inactive.
func (s *IAMService) UpdateServiceSpecificCredential(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := &UpdateServiceSpecificCredentialInput{
		ServiceSpecificCredentialId: request.GetStringParam(req.Parameters, "ServiceSpecificCredentialId"),
		UserName:                    request.GetStringParam(req.Parameters, "UserName"),
		Status:                      request.GetStringParam(req.Parameters, "Status"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.updateServiceSpecificCredentialCore(store, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func (s *IAMService) serviceSpecificCredentialToResponse(cred *iamstore.ServiceSpecificCredential, includePassword bool) map[string]interface{} {
	resp := map[string]interface{}{
		"ServiceSpecificCredentialId":  cred.ServiceSpecificCredentialId,
		"ServiceName":                  cred.ServiceName,
		"ServiceUserName":              cred.ServiceSpecificCredentialName,
		"UserName":                     cred.UserName,
		"ServiceSpecificCredentialArn": cred.ServiceSpecificCredentialArn,
		"CreateDate":                   cred.CreateDate.Format(timeutils.ISO8601SimpleFormat),
		"Status":                       cred.Status,
	}
	if cred.ExpirationDate != nil {
		resp["ExpirationDate"] = cred.ExpirationDate.Format(timeutils.ISO8601SimpleFormat)
	}
	if includePassword {
		resp["ServicePassword"] = cred.ServicePassword
	}
	return resp
}
