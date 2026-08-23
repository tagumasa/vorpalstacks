package cognitoidentity

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
)

// GetId obtains a unique identity ID for a Cognito identity pool.
func (s *CognitoIdentityService) GetId(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := req.GetParam("IdentityPoolId")
	if !validateIdentityPoolId(poolID) {
		return nil, ErrInvalidParameter
	}

	// AccountId is used for cross-account identity pool delegation in AWS.
	// It is accepted and validated for SPEC compliance; the edge environment
	// operates single-account so no cross-account check is enforced.
	if accountID := req.GetParam("AccountId"); accountID != "" {
		if !validateAccountId(accountID) {
			return nil, ErrInvalidParameter
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetIdentityPool(poolID)
	if err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityPoolNotFound)
	}

	logins := parseMapParam(req, "Logins")
	if !validateMapSize(len(logins), 10) || !validateLoginsKeys(logins) {
		return nil, ErrInvalidParameter
	}
	if !validateLoginsValues(logins) {
		return nil, ErrInvalidParameter
	}

	// For authenticated identities (Logins provided), AWS reuses the existing
	// identity whose logins match. Only create a new identity when no match
	// exists; the store resolves both paths under the pool's key lock so
	// concurrent callers cannot create duplicate identities.
	identity, err := store.GetOrCreateIdentityByLogins(poolID, logins)
	if err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"IdentityId": identity.ID,
	}, nil
}

// GetCredentialsForIdentity returns temporary credentials for an identity.
// In the enhanced authflow, this is functionally equivalent to calling
// GetOpenIdToken followed by AssumeRoleWithWebIdentity.
func (s *CognitoIdentityService) GetCredentialsForIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identityID := req.GetParam("IdentityId")
	if !validateIdentityId(identityID) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetIdentityByID(identityID)
	if err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityNotFound)
	}

	customRoleArn := req.GetParam("CustomRoleArn")
	if customRoleArn != "" && !validateRoleARN(customRoleArn) {
		return nil, ErrInvalidParameter
	}

	// When the caller provides fresh provider tokens via Logins, persist them
	// onto the identity so that subsequent role selection and credential
	// issuance reflect the current authentication state.
	if logins := parseMapParam(req, "Logins"); len(logins) > 0 {
		if !validateMapSize(len(logins), 10) || !validateLoginsKeys(logins) {
			return nil, ErrInvalidParameter
		}
		if !validateLoginsValues(logins) {
			return nil, ErrInvalidParameter
		}
		if identity.Logins == nil {
			identity.Logins = make(map[string]string)
		}
		for k, v := range logins {
			identity.Logins[k] = v
		}
		identity.LastModifiedDate = time.Now().UTC()
		if err := store.PutIdentity(identity); err != nil {
			return nil, ErrInternalError
		}
	}

	authRole, unauthRole, _, err := store.GetIdentityPoolRoles(identity.IdentityPoolID)
	if err != nil {
		return nil, ErrInvalidIdentityPoolConfig
	}

	// Determine which role to assume: authenticated if the identity has logins,
	// unauthenticated otherwise. CustomRoleArn takes precedence when provided.
	roleArn := customRoleArn
	if roleArn == "" {
		if len(identity.Logins) > 0 {
			roleArn = authRole
		} else {
			roleArn = unauthRole
		}
	}
	if roleArn == "" {
		return nil, ErrInvalidIdentityPoolConfig
	}

	if s.credentialIssuer == nil {
		return nil, ErrInternalError
	}

	result, err := s.credentialIssuer.IssueSession(roleArn, identityID, credentialSessionDurationSeconds)
	if err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"IdentityId": identityID,
		"Credentials": map[string]interface{}{
			"AccessKeyId":  result.AccessKeyID,
			"SecretKey":    result.SecretAccessKey,
			"SessionToken": result.SessionToken,
			"Expiration":   result.Expiration.Unix(),
		},
	}, nil
}

// DescribeIdentity returns information about a Cognito identity.
func (s *CognitoIdentityService) DescribeIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identityID := req.GetParam("IdentityId")
	if !validateIdentityId(identityID) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetIdentityByID(identityID)
	if err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityNotFound)
	}

	return map[string]interface{}{
		"IdentityId":       identity.ID,
		"CreationDate":     identity.CreationDate.Unix(),
		"LastModifiedDate": identity.LastModifiedDate.Unix(),
		"Logins":           formatLoginKeys(identity.Logins),
	}, nil
}
