package iot

import (
	"context"
	"strconv"
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

	credentialDuration := int64(3600)
	if cdStr := request.GetParamCaseInsensitive(req.Parameters, "credentialDurationSeconds"); cdStr != "" {
		if parsed, err := strconv.ParseInt(cdStr, 10, 64); err == nil {
			credentialDuration = parsed
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
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

	return roleAliasResponse(created), nil
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
		return nil, err
	}

	// DescribeRoleAlias output shape wraps fields in roleAliasDescription.
	// CreateRoleAlias and UpdateRoleAlias use flat shapes, so we wrap here
	// rather than in roleAliasResponse.
	return map[string]interface{}{
		"roleAliasDescription": roleAliasResponse(ra),
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

	var opts iotstore.RoleAliasUpdateOpts
	opts.RoleARN = request.GetParamCaseInsensitive(req.Parameters, "roleArn")
	if cdStr := request.GetParamCaseInsensitive(req.Parameters, "credentialDurationSeconds"); cdStr != "" {
		if parsed, err := strconv.ParseInt(cdStr, 10, 64); err == nil {
			opts.DurationSeconds = parsed
		}
	}

	ra, err := store.UpdateRoleAlias(roleAlias, opts)
	if err != nil {
		return nil, err
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

	arn := iotstore.BuildRoleAliasARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), roleAlias)
	_ = store.DeleteAllTags(arn)

	if err := store.DeleteRoleAlias(roleAlias); err != nil {
		return nil, err
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
