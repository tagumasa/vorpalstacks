package iam

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// CreateAccountAlias sets the specified alias for the AWS account.
func (s *IAMService) CreateAccountAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.createAccountAliasCore(store, request.GetStringParam(req.Parameters, "AccountAlias")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// DeleteAccountAlias removes the account alias for the AWS account.
func (s *IAMService) DeleteAccountAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteAccountAliasCore(store, request.GetStringParam(req.Parameters, "AccountAlias")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListAccountAliases lists the account alias associated with the AWS account.
func (s *IAMService) ListAccountAliases(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	alias, err := s.listAccountAliasesCore(store)
	if err != nil {
		return nil, err
	}

	aliases := []interface{}{}
	if alias != nil && alias.AccountAlias != "" {
		aliases = append(aliases, alias.AccountAlias)
	}

	return map[string]interface{}{
		"AccountAliases": aliases,
		"IsTruncated":    false,
	}, nil
}
