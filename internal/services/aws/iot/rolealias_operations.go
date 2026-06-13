package iot

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func (s *IoTService) CreateRoleAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleAlias := request.GetParamCaseInsensitive(req.Parameters, "roleAlias")
	roleARN := request.GetParamCaseInsensitive(req.Parameters, "roleArn")
	if roleAlias == "" || roleARN == "" {
		return nil, iotstore.ErrMissingParam
	}

	credentialDuration := request.GetParamCaseInsensitive(req.Parameters, "credentialDurationSeconds")
	if credentialDuration == "" {
		credentialDuration = "3600"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetRoleAlias(roleAlias); err == nil {
		return nil, iotstore.ErrResourceAlreadyExists
	}

	ra := &iotstore.RoleAlias{
		RoleAlias:                 roleAlias,
		RoleARN:                   roleARN,
		CredentialDurationSeconds: credentialDuration,
		Owner:                     s.accountID,
		CreationDate:              time.Now().UTC(),
		LastModifiedDate:          time.Now().UTC(),
	}

	created, err := store.CreateRoleAlias(ra)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"roleAlias":                 created.RoleAlias,
		"roleAliasArn":              created.RoleAliasARN,
		"roleArn":                   created.RoleARN,
		"credentialDurationSeconds": created.CredentialDurationSeconds,
		"creationDate":              created.CreationDate.Unix(),
		"lastModifiedDate":          created.LastModifiedDate.Unix(),
	}, nil
}

func (s *IoTService) DescribeRoleAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleAlias := request.GetParamCaseInsensitive(req.Parameters, "roleAlias")
	if roleAlias == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ra, err := store.GetRoleAlias(roleAlias)
	if err != nil {
		return nil, iotstore.ErrRoleAliasNotFound
	}

	return map[string]interface{}{
		"roleAliasDescription": map[string]interface{}{
			"roleAlias":                 ra.RoleAlias,
			"roleAliasArn":              ra.RoleAliasARN,
			"roleArn":                   ra.RoleARN,
			"credentialDurationSeconds": ra.CredentialDurationSeconds,
			"creationDate":              ra.CreationDate.Unix(),
			"lastModifiedDate":          ra.LastModifiedDate.Unix(),
		},
	}, nil
}

func (s *IoTService) UpdateRoleAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleAlias := request.GetParamCaseInsensitive(req.Parameters, "roleAlias")
	if roleAlias == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ra, err := store.GetRoleAlias(roleAlias)
	if err != nil {
		return nil, iotstore.ErrRoleAliasNotFound
	}

	if roleARN := request.GetParamCaseInsensitive(req.Parameters, "roleArn"); roleARN != "" {
		ra.RoleARN = roleARN
	}
	if dur := request.GetParamCaseInsensitive(req.Parameters, "credentialDurationSeconds"); dur != "" {
		ra.CredentialDurationSeconds = dur
	}
	ra.LastModifiedDate = time.Now().UTC()

	if err := store.UpdateRoleAlias(ra); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{
		"roleAlias":        ra.RoleAlias,
		"roleAliasArn":     ra.RoleAliasARN,
		"lastModifiedDate": ra.LastModifiedDate.Unix(),
	}, nil
}

func (s *IoTService) DeleteRoleAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleAlias := request.GetParamCaseInsensitive(req.Parameters, "roleAlias")
	if roleAlias == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteRoleAlias(roleAlias); err != nil {
		return nil, iotstore.ErrRoleAliasNotFound
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListRoleAliases(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	aliases, err := store.ListRoleAliases(parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(aliases.Items))
	for _, a := range aliases.Items {
		names = append(names, a.RoleAlias)
	}

	resp := map[string]interface{}{
		"roleAliases": names,
	}
	if aliases.NextMarker != "" {
		resp["nextToken"] = aliases.NextMarker
	}
	return resp, nil
}
