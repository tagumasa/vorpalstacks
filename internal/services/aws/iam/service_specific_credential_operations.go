package iam

import (
	"context"

	"vorpalstacks/internal/common/request"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateServiceSpecificCredential generates a new service-specific credential for the specified user and service.
func (s *IAMService) CreateServiceSpecificCredential(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, ErrNoSuchUser
	}
	serviceName := request.GetStringParam(req.Parameters, "ServiceName")
	if serviceName == "" {
		return nil, NewValidationError("ServiceName")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	cred, err := store.ServiceSpecificCredentials().Create(userName, serviceName)
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
	if credentialId == "" {
		return nil, NewValidationError("ServiceSpecificCredentialId")
	}
	userName := request.GetStringParam(req.Parameters, "UserName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if userName != "" && !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}
	if !store.ServiceSpecificCredentials().Exists(credentialId) {
		return nil, NewNoSuchEntityError("service-specific credential", credentialId)
	}

	if err := store.ServiceSpecificCredentials().Delete(credentialId); err != nil {
		return nil, err
	}

	return nil, nil
}

// ListServiceSpecificCredentials lists all service-specific credentials for the specified user.
func (s *IAMService) ListServiceSpecificCredentials(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, ErrNoSuchUser
	}
	serviceName := request.GetStringParam(req.Parameters, "ServiceName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	creds, err := store.ServiceSpecificCredentials().ListByUserName(userName)
	if err != nil {
		return nil, err
	}

	credList := make([]interface{}, 0)
	for _, cred := range creds {
		if serviceName != "" && cred.ServiceName != serviceName {
			continue
		}
		credList = append(credList, s.serviceSpecificCredentialToResponse(cred, false))
	}

	return map[string]interface{}{
		"ServiceSpecificCredentials": credList,
	}, nil
}

// ResetServiceSpecificCredential resets the password for a service-specific credential.
func (s *IAMService) ResetServiceSpecificCredential(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	credentialId := request.GetStringParam(req.Parameters, "ServiceSpecificCredentialId")
	if credentialId == "" {
		return nil, NewValidationError("ServiceSpecificCredentialId")
	}
	userName := request.GetStringParam(req.Parameters, "UserName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if userName != "" && !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	cred, err := store.ServiceSpecificCredentials().Get(credentialId)
	if err != nil {
		return nil, NewNoSuchEntityError("service-specific credential", credentialId)
	}

	newCred, err := store.ServiceSpecificCredentials().Create(cred.UserName, cred.ServiceName)
	if err != nil {
		return nil, err
	}

	if err := store.ServiceSpecificCredentials().Delete(credentialId); err != nil {
		return nil, err
	}

	newCred.ServiceSpecificCredentialId = credentialId
	if err := store.ServiceSpecificCredentials().Put(newCred); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ServiceSpecificCredential": s.serviceSpecificCredentialToResponse(newCred, true),
	}, nil
}

// UpdateServiceSpecificCredential sets the status of a service-specific credential to Active or Inactive.
func (s *IAMService) UpdateServiceSpecificCredential(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	credentialId := request.GetStringParam(req.Parameters, "ServiceSpecificCredentialId")
	if credentialId == "" {
		return nil, NewValidationError("ServiceSpecificCredentialId")
	}
	status := request.GetStringParam(req.Parameters, "Status")
	if status == "" {
		return nil, NewValidationError("Status")
	}
	if status != "Active" && status != "Inactive" {
		return nil, NewInvalidInputError("Status", "must be Active or Inactive")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.ServiceSpecificCredentials().Exists(credentialId) {
		return nil, NewNoSuchEntityError("service-specific credential", credentialId)
	}

	if err := store.ServiceSpecificCredentials().UpdateStatus(credentialId, status); err != nil {
		return nil, err
	}

	return nil, nil
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
	if includePassword {
		resp["ServicePassword"] = cred.ServicePassword
	}
	return resp
}
