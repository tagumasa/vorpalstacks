package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

func (s *IoTService) CreateRoleAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	duration, durationProvided := parseOptionalInt64Param(req.Parameters, "credentialDurationSeconds")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	created, err := s.createRoleAliasCore(store, CreateRoleAliasInput{
		RoleAlias:                 request.GetParamCaseInsensitive(req.Parameters, "roleAlias"),
		RoleARN:                   request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		CredentialDurationSeconds: duration,
		DurationProvided:          durationProvided,
	})
	if err != nil {
		return nil, err
	}

	return roleAliasResponse(created), nil
}

func (s *IoTService) DescribeRoleAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ra, err := s.describeRoleAliasCore(store, request.GetParamCaseInsensitive(req.Parameters, "roleAlias"))
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
	duration, durationProvided := parseOptionalInt64Param(req.Parameters, "credentialDurationSeconds")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ra, err := s.updateRoleAliasCore(store, UpdateRoleAliasInput{
		RoleAlias:                 request.GetParamCaseInsensitive(req.Parameters, "roleAlias"),
		RoleARN:                   request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		CredentialDurationSeconds: duration,
		DurationProvided:          durationProvided,
	})
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteRoleAliasCore(store, request.GetParamCaseInsensitive(req.Parameters, "roleAlias")); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListRoleAliases(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := parseListOptions(req.Parameters)
	result, err := s.listRoleAliasesCore(store, opts.Marker, opts.MaxItems)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"roleAliases": result.RoleAliases,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
	}
	return resp, nil
}

// parseOptionalInt64Param parses an optional integer wire parameter,
// reporting whether a valid number was supplied. Numeric wire values arrive
// as JSON numbers, so the integer reader is required — the string reader
// silently drops them.
func parseOptionalInt64Param(params map[string]interface{}, key string) (int64, bool) {
	v, ok := request.GetIntParamCaseInsensitive(params, key)
	if !ok {
		return 0, false
	}
	return int64(v), true
}
