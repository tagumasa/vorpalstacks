package iam

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// CreateAccountAlias sets the specified alias for the AWS account.
func (s *IAMService) CreateAccountAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	alias := request.GetStringParam(req.Parameters, "AccountAlias")
	if alias == "" {
		return nil, NewValidationError("AccountAlias")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.AccountAlias().Put(alias); err != nil {
		return nil, err
	}

	return nil, nil
}

// DeleteAccountAlias removes the account alias for the AWS account.
func (s *IAMService) DeleteAccountAlias(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.AccountAlias().Delete(); err != nil {
		return nil, err
	}

	return nil, nil
}

// ListAccountAliases lists the account alias associated with the AWS account.
func (s *IAMService) ListAccountAliases(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	alias, err := store.AccountAlias().Get()
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
