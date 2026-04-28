package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// CreateGrant creates a grant for a KMS key.
func (s *KMSService) CreateGrant(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "CreateGrant", nil)
	if err != nil {
		return nil, err
	}

	granteePrincipal := request.GetStringParam(req.Parameters, "GranteePrincipal")
	if granteePrincipal == "" {
		return nil, ErrAccessDenied
	}

	retiringPrincipal := request.GetStringParam(req.Parameters, "RetiringPrincipal")
	name := request.GetStringParam(req.Parameters, "Name")

	var operations []string
	if ops, ok := req.Parameters["Operations"]; ok {
		if opList, ok := ops.([]interface{}); ok {
			for _, op := range opList {
				if opStr, ok := op.(string); ok {
					operations = append(operations, opStr)
				}
			}
		}
	}

	constraints := parseGrantConstraints(req.Parameters)

	grantToken, err := kmsstore.GenerateGrantToken()
	if err != nil {
		return nil, err
	}

	grant, err := stores.grants.CreateWithToken(key.KeyID, granteePrincipal, retiringPrincipal, operations, name, constraints, grantToken)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"GrantId":    grant.GrantID,
		"GrantToken": grantToken,
		"KeyId":      key.Arn,
	}, nil
}

// ListGrants retrieves grants for the specified KMS key.
func (s *KMSService) ListGrants(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "ListGrants", nil)
	if err != nil {
		return nil, err
	}
	marker := pagination.GetMarker(req.Parameters)
	maxItems := pagination.GetMaxItems(req.Parameters, 100)
	granteePrincipal := request.GetStringParam(req.Parameters, "GranteePrincipal")

	result, err := stores.grants.List(key.KeyID, granteePrincipal, marker, maxItems)
	if err != nil {
		return nil, err
	}

	grants := make([]map[string]interface{}, len(result.Grants))
	for i, g := range result.Grants {
		grant := map[string]interface{}{
			"KeyId":            key.Arn,
			"GrantId":          g.GrantID,
			"GranteePrincipal": g.GranteePrincipal,
			"Operations":       g.Operations,
			"IssuingAccount":   g.IssuingAccount,
			"CreationDate":     g.CreationDate.Unix(),
		}
		if g.Name != "" {
			grant["Name"] = g.Name
		}
		if g.RetiringPrincipal != "" {
			grant["RetiringPrincipal"] = g.RetiringPrincipal
		}
		if g.Constraints != nil {
			grant["Constraints"] = g.Constraints
		}
		grants[i] = grant
	}

	response := map[string]interface{}{
		"Grants":    grants,
		"Truncated": result.IsTruncated,
	}
	if result.IsTruncated {
		response["NextMarker"] = result.NextMarker
	}

	return response, nil
}

// ListRetirableGrants retrieves grants that can be retired by the specified retiring principal.
func (s *KMSService) ListRetirableGrants(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	retiringPrincipal := request.GetStringParam(req.Parameters, "RetiringPrincipal")
	if retiringPrincipal == "" {
		return nil, ErrValidation
	}

	marker := pagination.GetMarker(req.Parameters)
	maxItems := pagination.GetMaxItems(req.Parameters, 100)

	result, err := stores.grants.ListByRetiringPrincipal(retiringPrincipal, marker, maxItems)
	if err != nil {
		return nil, err
	}

	grants := make([]map[string]interface{}, len(result.Grants))
	for i, g := range result.Grants {
		key, err := stores.keys.Get(g.KeyID)
		if err != nil {
			continue
		}

		grant := map[string]interface{}{
			"KeyId":            key.Arn,
			"GrantId":          g.GrantID,
			"GranteePrincipal": g.GranteePrincipal,
			"Operations":       g.Operations,
			"IssuingAccount":   g.IssuingAccount,
			"CreationDate":     g.CreationDate.Unix(),
		}
		if g.Name != "" {
			grant["Name"] = g.Name
		}
		if g.RetiringPrincipal != "" {
			grant["RetiringPrincipal"] = g.RetiringPrincipal
		}
		if g.Constraints != nil {
			grant["Constraints"] = g.Constraints
		}
		grants[i] = grant
	}

	response := map[string]interface{}{
		"Grants":    grants,
		"Truncated": result.IsTruncated,
	}
	if result.IsTruncated {
		response["NextMarker"] = result.NextMarker
	}

	return response, nil
}

// RevokeGrant revokes a grant from a KMS key.
func (s *KMSService) RevokeGrant(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "RevokeGrant", nil)
	if err != nil {
		return nil, err
	}

	grantID := request.GetStringParam(req.Parameters, "GrantId")
	if grantID == "" {
		return nil, ErrGrantNotFound
	}

	grant, err := stores.grants.Get(grantID)
	if err != nil {
		if errors.Is(err, kmsstore.ErrGrantNotFound) {
			return nil, ErrGrantNotFound
		}
		return nil, err
	}

	if grant.KeyID != key.KeyID {
		return nil, ErrGrantNotFound
	}

	if err := stores.grants.Delete(grantID); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// RetireGrant retires a grant from a KMS key.
func (s *KMSService) RetireGrant(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	grantID := request.GetStringParam(req.Parameters, "GrantId")
	grantToken := request.GetStringParam(req.Parameters, "GrantToken")

	var grant *kmsstore.Grant
	if grantID != "" {
		var err error
		grant, err = stores.grants.Get(grantID)
		if err != nil {
			if errors.Is(err, kmsstore.ErrGrantNotFound) {
				return nil, ErrGrantNotFound
			}
			return nil, err
		}
	} else if grantToken != "" {
		var err error
		grant, err = stores.grants.GetByToken(grantToken)
		if err != nil {
			if errors.Is(err, kmsstore.ErrGrantNotFound) {
				return nil, ErrGrantNotFound
			}
			return nil, err
		}
		grantID = grant.GrantID
	} else {
		return nil, ErrGrantNotFound
	}

	keyID := s.getKeyID(req.Parameters)
	if keyID != "" {
		key, err := s.resolveKey(stores, req.Parameters)
		if err != nil {
			return nil, err
		}
		if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "RetireGrant", key.KeyID, nil); err != nil {
			return nil, err
		}
		if grant.KeyID != key.KeyID {
			return nil, ErrGrantNotFound
		}
	}

	if err := stores.grants.Delete(grantID); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func parseGrantConstraints(params map[string]interface{}) *kmsstore.GrantConstraints {
	var constraints *kmsstore.GrantConstraints

	if c, ok := params["Constraints"]; ok {
		if cmap, ok := c.(map[string]interface{}); ok {
			if ecEquals, ok := cmap["EncryptionContextEquals"]; ok {
				if ecMap, ok := ecEquals.(map[string]interface{}); ok {
					if constraints == nil {
						constraints = &kmsstore.GrantConstraints{}
					}
					constraints.EncryptionContextEquals = make(map[string]string)
					for k, v := range ecMap {
						if vs, ok := v.(string); ok {
							constraints.EncryptionContextEquals[k] = vs
						}
					}
				}
			}
			if ecSubset, ok := cmap["EncryptionContextSubset"]; ok {
				if ecMap, ok := ecSubset.(map[string]interface{}); ok {
					if constraints == nil {
						constraints = &kmsstore.GrantConstraints{}
					}
					constraints.EncryptionContextSubset = make(map[string]string)
					for k, v := range ecMap {
						if vs, ok := v.(string); ok {
							constraints.EncryptionContextSubset[k] = vs
						}
					}
				}
			}
		}
	}

	return constraints
}
