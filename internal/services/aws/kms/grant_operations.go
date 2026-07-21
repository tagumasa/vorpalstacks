package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"
	"errors"
	"regexp"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
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

	granteePrincipal := request.GetStringParam(req.Parameters, "GranteePrincipal")
	if granteePrincipal == "" {
		// AWS: GranteePrincipal is a required field; missing is a
		// ValidationException, not AccessDenied.
		return nil, ErrValidation
	}

	retiringPrincipal := request.GetStringParam(req.Parameters, "RetiringPrincipal")
	name := request.GetStringParam(req.Parameters, "Name")
	// AWS: Name is optional but if present must be 1-256 chars matching
	// the grantNamePattern (alnum, colon, slash, underscore, hyphen).
	if name != "" && !grantNamePattern.MatchString(name) {
		return nil, ErrValidation
	}

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
	// AWS: Operations is a required field for CreateGrant.
	if len(operations) == 0 {
		return nil, ErrValidation
	}

	constraints, err := parseGrantConstraints(req.Parameters)
	if err != nil {
		return nil, err
	}

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
		grants[i] = s.buildGrantResponse(g, key.Arn)
	}

	return s.buildGrantsListResponse(grants, result.IsTruncated, result.NextMarker), nil
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

	var grants []map[string]interface{}
	grants = make([]map[string]interface{}, 0)
	for _, g := range result.Grants {
		key, err := stores.keys.Get(g.KeyID)
		if err != nil {
			// A grant whose key cannot be resolved is an orphan — typically
			// the result of a partially-failed cascade-delete. Skip the
			// entry (matching AWS behaviour where retired/deleted keys do
			// not appear in ListRetirableGrants) but log loudly so the
			// operator detects the data-integrity issue. The previous
			// code returned ErrKMSInternal here, which broke the entire
			// list response for a single orphaned grant.
			logs.Error("ListRetirableGrants: skipping orphaned grant (key not found)", logs.String("keyId", g.KeyID), logs.String("grantId", g.GrantID), logs.Err(err))
			continue
		}

		grant := s.buildGrantResponse(g, key.Arn)
		grants = append(grants, grant)
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

	grantID := request.GetStringParam(req.Parameters, "GrantId")
	if grantID == "" {
		// AWS: GrantId is required; missing is a ValidationException,
		// not NotFoundException.
		return nil, ErrValidation
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
		// AWS: at least one of GrantId or GrantToken is required;
		// missing both is a ValidationException, not NotFoundException.
		return nil, ErrValidation
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

func parseGrantConstraints(params map[string]interface{}) (*kmsstore.GrantConstraints, error) {
	var constraints *kmsstore.GrantConstraints

	if c, ok := params["Constraints"]; ok {
		if cmap, ok := c.(map[string]interface{}); ok {
			// Reject unknown constraint members up-front so that future
			// AWS spec additions surface as ValidationException rather
			// than being silently dropped (over-authorising the grant).
			for k := range cmap {
				switch k {
				case "EncryptionContextEquals", "EncryptionContextSubset", "SourceArn":
				default:
					return nil, ErrValidation
				}
			}
			if ecEquals, ok := cmap["EncryptionContextEquals"]; ok {
				ecMap, ok := ecEquals.(map[string]interface{})
				if !ok {
					// Smithy: EncryptionContextEquals is map<string,string>.
					// A non-map value is a malformed request.
					return nil, ErrValidation
				}
				if constraints == nil {
					constraints = &kmsstore.GrantConstraints{}
				}
				constraints.EncryptionContextEquals = make(map[string]string)
				for k, v := range ecMap {
					vs, ok := v.(string)
					if !ok {
						// Smithy: EncryptionContextValue is a string.
						// Non-string values are ValidationException
						// rather than being silently dropped (which
						// would weaken the constraint).
						return nil, ErrValidation
					}
					constraints.EncryptionContextEquals[k] = vs
				}
			}
			if ecSubset, ok := cmap["EncryptionContextSubset"]; ok {
				ecMap, ok := ecSubset.(map[string]interface{})
				if !ok {
					return nil, ErrValidation
				}
				if constraints == nil {
					constraints = &kmsstore.GrantConstraints{}
				}
				constraints.EncryptionContextSubset = make(map[string]string)
				for k, v := range ecMap {
					vs, ok := v.(string)
					if !ok {
						return nil, ErrValidation
					}
					constraints.EncryptionContextSubset[k] = vs
				}
			}
			// Smithy com.amazonaws.kms#GrantConstraints has three members;
			// SourceArn was previously dropped silently, over-authorising
			// grants whose callers depended on the constraint.
			if sourceArnVal, ok := cmap["SourceArn"]; ok {
				sourceArn, ok := sourceArnVal.(string)
				if !ok || sourceArn == "" {
					return nil, ErrValidation
				}
				if constraints == nil {
					constraints = &kmsstore.GrantConstraints{}
				}
				constraints.SourceArn = sourceArn
			}
		}
	}

	return constraints, nil
}
