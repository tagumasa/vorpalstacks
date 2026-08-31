package kms

// grant_core.go carries the Core functions of the KMS grant family. The
// write Cores replay the original validation ladder (required members,
// principal/name/ID patterns, constraint shapes, dry-run gate) before
// minting a token or touching the grant store; the list Cores keep the
// documented pagination and the GrantId direct-fetch short circuit.

import (
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// GrantCreateInput carries the raw wire parameters of CreateGrant
// (GranteePrincipal, RetiringPrincipal, Name, Operations, GrantTokens,
// Constraints, DryRun) so the Core can validate every member in the
// original order.
type GrantCreateInput struct {
	Params map[string]interface{}
}

// createGrantCore is the single entry point for creating a grant. It
// returns the created grant and its token so the caller can serialise
// both.
func (s *KMSService) createGrantCore(stores *kmsStores, key *kmsstore.Key, in GrantCreateInput) (*kmsstore.Grant, string, error) {
	granteePrincipal := request.GetStringParam(in.Params, "GranteePrincipal")
	if granteePrincipal == "" {
		// AWS: GranteePrincipal is a required field; missing is a
		// ValidationException, not AccessDenied.
		return nil, "", ErrValidation
	}
	if err := validatePrincipalId(granteePrincipal); err != nil {
		return nil, "", err
	}

	retiringPrincipal := request.GetStringParam(in.Params, "RetiringPrincipal")
	if retiringPrincipal != "" {
		if err := validatePrincipalId(retiringPrincipal); err != nil {
			return nil, "", err
		}
	}

	name := request.GetStringParam(in.Params, "Name")
	// AWS: Name is optional but if present must be 1-256 chars matching
	// the grantNamePattern (alnum, colon, slash, underscore, hyphen).
	if name != "" && !grantNamePattern.MatchString(name) {
		return nil, "", ErrValidation
	}

	var operations []string
	if ops, ok := in.Params["Operations"]; ok {
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
		return nil, "", ErrValidation
	}

	// Smithy GrantTokenList: length 0-10.
	if gt, ok := in.Params["GrantTokens"]; ok {
		if gtList, ok := gt.([]interface{}); ok {
			if err := validateGrantTokenListSize(gtList); err != nil {
				return nil, "", err
			}
		}
	}

	constraints, err := parseGrantConstraints(in.Params)
	if err != nil {
		return nil, "", err
	}

	if err := checkKMSDryRun(in.Params); err != nil {
		return nil, "", err
	}

	grantToken, err := kmsstore.GenerateGrantToken()
	if err != nil {
		return nil, "", err
	}

	grant, err := stores.grants.CreateWithToken(key.KeyID, granteePrincipal, retiringPrincipal, operations, name, constraints, grantToken)
	if err != nil {
		return nil, "", err
	}

	return grant, grantToken, nil
}

// GrantListResult is the Core result of the grant list operations: the
// store grants plus the pagination markers.
type GrantListResult struct {
	Grants      []*kmsstore.Grant
	IsTruncated bool
	NextMarker  string
}

// listGrantsCore is the single entry point for listing grants on a key.
// When GrantId is specified it fetches the single grant directly (with
// the key-ownership check) instead of paginating; GrantId uniquely
// identifies a grant, so post-filter pagination is unnecessary and would
// return incorrect IsTruncated/NextMarker values. The GranteeServicePrincipal
// filter is applied with the original branch semantics: the GrantId path
// yields a non-nil empty list on mismatch, the list path starts from a
// nil list — callers must preserve that nil-ness when serialising.
func (s *KMSService) listGrantsCore(stores *kmsStores, key *kmsstore.Key, marker string, maxItems int, granteePrincipal, grantIDFilter, granteeServicePrincipal string) (*GrantListResult, error) {
	if err := validateMarkerLength(marker); err != nil {
		return nil, err
	}
	if grantIDFilter != "" {
		if err := validateGrantIdLength(grantIDFilter); err != nil {
			return nil, err
		}
	}

	if grantIDFilter != "" {
		grant, err := stores.grants.Get(grantIDFilter)
		if err != nil {
			if kmsstore.IsNotFound(err) {
				return nil, ErrGrantNotFound
			}
			return nil, err
		}
		if grant.KeyID != key.KeyID {
			return nil, ErrGrantNotFound
		}
		if granteeServicePrincipal != "" && grant.GranteePrincipal != granteeServicePrincipal {
			return &GrantListResult{Grants: []*kmsstore.Grant{}}, nil
		}
		return &GrantListResult{Grants: []*kmsstore.Grant{grant}}, nil
	}

	result, err := stores.grants.List(key.KeyID, granteePrincipal, marker, maxItems)
	if err != nil {
		return nil, err
	}

	// Apply GranteeServicePrincipal filter (post-filter is acceptable
	// here because GranteeServicePrincipal is a coarse filter that rarely
	// splits across pages in practice).
	var grants []*kmsstore.Grant
	for _, g := range result.Grants {
		if granteeServicePrincipal != "" && g.GranteePrincipal != granteeServicePrincipal {
			continue
		}
		grants = append(grants, g)
	}

	return &GrantListResult{
		Grants:      grants,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}

// RetirableGrantEntry pairs a retirable grant with its key ARN, the two
// members ListRetirableGrants serialises per entry.
type RetirableGrantEntry struct {
	Grant  *kmsstore.Grant
	KeyArn string
}

// RetirableGrantsResult is the Core result of ListRetirableGrants.
type RetirableGrantsResult struct {
	Grants      []RetirableGrantEntry
	IsTruncated bool
	NextMarker  string
}

// listRetirableGrantsCore is the single entry point for listing grants a
// retiring principal can retire. A grant whose key cannot be resolved is
// an orphan — typically the result of a partially-failed cascade-delete —
// and is skipped so a single orphan cannot break the list response.
func (s *KMSService) listRetirableGrantsCore(stores *kmsStores, retiringPrincipal, marker string, maxItems int) (*RetirableGrantsResult, error) {
	if retiringPrincipal == "" {
		return nil, ErrValidation
	}

	if err := validateMarkerLength(marker); err != nil {
		return nil, err
	}

	result, err := stores.grants.ListByRetiringPrincipal(retiringPrincipal, marker, maxItems)
	if err != nil {
		return nil, err
	}

	grants := make([]RetirableGrantEntry, 0, len(result.Grants))
	for _, g := range result.Grants {
		key, err := stores.keys.Get(g.KeyID)
		if err != nil {
			// Skip the entry (matching AWS behaviour where retired/deleted
			// keys do not appear in ListRetirableGrants) but log loudly so
			// the operator detects the data-integrity issue. The previous
			// code returned ErrKMSInternal here, which broke the entire
			// list response for a single orphaned grant.
			logs.Error("ListRetirableGrants: skipping orphaned grant (key not found)", logs.String("keyId", g.KeyID), logs.String("grantId", g.GrantID), logs.Err(err))
			continue
		}

		grants = append(grants, RetirableGrantEntry{Grant: g, KeyArn: key.Arn})
	}

	return &RetirableGrantsResult{
		Grants:      grants,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}

// GrantRevokeInput carries the RevokeGrant members; Params transports the
// raw wire parameters so the Core can evaluate the DryRun gate in its
// original position.
type GrantRevokeInput struct {
	GrantID string
	Params  map[string]interface{}
}

// revokeGrantCore is the single entry point for revoking a grant.
func (s *KMSService) revokeGrantCore(stores *kmsStores, key *kmsstore.Key, in GrantRevokeInput) error {
	if in.GrantID == "" {
		// AWS: GrantId is required; missing is a ValidationException,
		// not NotFoundException.
		return ErrValidation
	}
	if err := validateGrantIdLength(in.GrantID); err != nil {
		return err
	}

	grant, err := stores.grants.Get(in.GrantID)
	if err != nil {
		if kmsstore.IsNotFound(err) {
			return ErrGrantNotFound
		}
		return err
	}

	if grant.KeyID != key.KeyID {
		return ErrGrantNotFound
	}

	if err := checkKMSDryRun(in.Params); err != nil {
		return err
	}

	return stores.grants.Delete(in.GrantID)
}

// GrantRetireInput carries the RetireGrant members. RetireGrant resolves
// the grant by GrantId or GrantToken and optionally authorises against a
// caller-supplied key, so Params transports the full wire parameters.
type GrantRetireInput struct {
	Principal string
	Params    map[string]interface{}
}

// retireGrantCore is the single entry point for retiring a grant.
func (s *KMSService) retireGrantCore(stores *kmsStores, in GrantRetireInput) error {
	grantID := request.GetStringParam(in.Params, "GrantId")
	grantToken := request.GetStringParam(in.Params, "GrantToken")

	var grant *kmsstore.Grant
	if grantID != "" {
		if err := validateGrantIdLength(grantID); err != nil {
			return err
		}
		var err error
		grant, err = stores.grants.Get(grantID)
		if err != nil {
			if kmsstore.IsNotFound(err) {
				return ErrGrantNotFound
			}
			return err
		}
	} else if grantToken != "" {
		if err := validateGrantTokenLength(grantToken); err != nil {
			return err
		}
		var err error
		grant, err = stores.grants.GetByToken(grantToken)
		if err != nil {
			if kmsstore.IsNotFound(err) {
				return ErrGrantNotFound
			}
			return err
		}
		grantID = grant.GrantID
	} else {
		// AWS: at least one of GrantId or GrantToken is required;
		// missing both is a ValidationException, not NotFoundException.
		return ErrValidation
	}

	keyID := s.getKeyID(in.Params)
	if keyID != "" {
		key, err := s.resolveKey(stores, in.Params)
		if err != nil {
			return err
		}
		if err := s.authorizeOperation(stores, in.Principal, "RetireGrant", key.KeyID, nil); err != nil {
			return err
		}
		if grant.KeyID != key.KeyID {
			return ErrGrantNotFound
		}
	}

	if err := checkKMSDryRun(in.Params); err != nil {
		return err
	}

	return stores.grants.Delete(grantID)
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
				if err := validateGrantSourceArn(sourceArn); err != nil {
					return nil, err
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
