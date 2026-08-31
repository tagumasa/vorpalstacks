package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"
	"regexp"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// grantNamePattern mirrors the AWS GrantNameType constraint: 1-256 chars
// of alphanumerics plus colon, slash, underscore, and hyphen. Length bounds
// come from the smithy.api#length trait on com.amazonaws.kms#GrantNameType.
var grantNamePattern = regexp.MustCompile(`^[a-zA-Z0-9:/_-]{1,256}$`)

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

	grant, grantToken, err := s.createGrantCore(stores, key, GrantCreateInput{Params: req.Parameters})
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

	result, err := s.listGrantsCore(stores, key,
		pagination.GetMarker(req.Parameters),
		pagination.GetMaxItems(req.Parameters, 100),
		request.GetStringParam(req.Parameters, "GranteePrincipal"),
		request.GetStringParam(req.Parameters, "GrantId"),
		request.GetStringParam(req.Parameters, "GranteeServicePrincipal"))
	if err != nil {
		return nil, err
	}

	// The Core applies the GranteeServicePrincipal filter with the
	// original branch semantics; preserve the nil-ness of the grant list
	// when mapping to the wire shape (a nil list serialises as null, an
	// empty one as []).
	var grants []map[string]interface{}
	if result.Grants != nil {
		grants = make([]map[string]interface{}, 0, len(result.Grants))
		for _, g := range result.Grants {
			grants = append(grants, s.buildGrantResponse(g, key.Arn))
		}
	}

	return s.buildGrantsListResponse(grants, result.IsTruncated, result.NextMarker), nil
}

// ListRetirableGrants retrieves grants that can be retired by the specified retiring principal.
func (s *KMSService) ListRetirableGrants(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listRetirableGrantsCore(stores,
		request.GetStringParam(req.Parameters, "RetiringPrincipal"),
		pagination.GetMarker(req.Parameters),
		pagination.GetMaxItems(req.Parameters, 100))
	if err != nil {
		return nil, err
	}

	grants := make([]map[string]interface{}, 0, len(result.Grants))
	for _, entry := range result.Grants {
		grants = append(grants, s.buildGrantResponse(entry.Grant, entry.KeyArn))
	}

	return s.buildGrantsListResponse(grants, result.IsTruncated, result.NextMarker), nil
}

// buildGrantResponse constructs a grant response map from a store grant.
func (s *KMSService) buildGrantResponse(g *kmsstore.Grant, keyArn string) map[string]interface{} {
	grant := map[string]interface{}{
		"KeyId":            keyArn,
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
		// AWS expects Constraints as an object with EncryptionContextEquals,
		// EncryptionContextSubset, and/or SourceArn fields. Emitting the raw
		// store struct produces an empty Constraints object on the wire
		// because the struct's json tags use snake_case while AWS uses
		// PascalCase member names. Build the map explicitly and include
		// SourceArn when set so that round-tripping a grant preserves the
		// constraint.
		constraintsMap := map[string]interface{}{}
		if len(g.Constraints.EncryptionContextEquals) > 0 {
			constraintsMap["EncryptionContextEquals"] = g.Constraints.EncryptionContextEquals
		}
		if len(g.Constraints.EncryptionContextSubset) > 0 {
			constraintsMap["EncryptionContextSubset"] = g.Constraints.EncryptionContextSubset
		}
		if g.Constraints.SourceArn != "" {
			constraintsMap["SourceArn"] = g.Constraints.SourceArn
		}
		if len(constraintsMap) > 0 {
			grant["Constraints"] = constraintsMap
		}
	}
	return grant
}

// buildGrantsListResponse constructs a paginated grants list response.
func (s *KMSService) buildGrantsListResponse(grants []map[string]interface{}, isTruncated bool, nextMarker string) map[string]interface{} {
	response := map[string]interface{}{
		"Grants":    grants,
		"Truncated": isTruncated,
	}
	if isTruncated {
		response["NextMarker"] = nextMarker
	}
	return response
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

	if err := s.revokeGrantCore(stores, key, GrantRevokeInput{
		GrantID: request.GetStringParam(req.Parameters, "GrantId"),
		Params:  req.Parameters,
	}); err != nil {
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

	if err := s.retireGrantCore(stores, GrantRetireInput{
		Principal: s.resolveCallerPrincipal(reqCtx, req),
		Params:    req.Parameters,
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
