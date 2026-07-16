package cognitoidentity

import (
	"context"

	"vorpalstacks/internal/common/request"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
)

// GetId obtains a unique identity ID for a Cognito identity pool.
func (s *CognitoIdentityService) GetId(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolID := req.GetParam("IdentityPoolId")
	if poolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetIdentityPool(poolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	logins := parseMapParam(req, "Logins")

	// For authenticated identities (Logins provided), AWS reuses the existing
	// identity whose logins match. Only create a new identity when no match exists.
	if len(logins) > 0 {
		if existing, err := store.FindIdentityByLogins(poolID, logins); err == nil && existing != nil {
			return map[string]interface{}{
				"IdentityId": existing.ID,
			}, nil
		}
	}

	identity := cognitoidentitystore.NewIdentity(poolID)
	if len(logins) > 0 {
		identity.Logins = logins
	}

	if err := store.CreateIdentity(identity); err != nil {
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
	if identityID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetIdentityByID(identityID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	customRoleArn := req.GetParam("CustomRoleArn")
	_ = parseMapParam(req, "Logins")

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

	result, err := s.credentialIssuer.IssueSession(roleArn, identityID, 3600)
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
	if identityID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetIdentityByID(identityID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{
		"IdentityId":       identity.ID,
		"CreationDate":     identity.CreationDate.Unix(),
		"LastModifiedDate": identity.LastModifiedDate.Unix(),
		"Logins":           formatLoginKeys(identity.Logins),
	}, nil
}
