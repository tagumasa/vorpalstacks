package cognitoidentity

import (
	"time"

	"vorpalstacks/internal/common/request"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input/Result DTOs — identity lifecycle operations
// ---------------------------------------------------------------------------

// GetIdInput carries every field that GetId needs, in a format independent
// of the wire protocol (HTTP Query/JSON vs gRPC-Web).
type GetIdInput struct {
	IdentityPoolID string
	AccountID      string
	Logins         map[string]string
}

// GetCredentialsForIdentityInput carries every field that
// GetCredentialsForIdentity needs. Logins carries the caller-provided
// provider tokens (nil when the member is absent).
type GetCredentialsForIdentityInput struct {
	IdentityID    string
	CustomRoleARN string
	Logins        map[string]string
}

// CredentialsResult is the transport-agnostic credential payload issued for
// a Cognito identity by the enhanced authflow.
type CredentialsResult struct {
	IdentityID      string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	ExpirationUnix  int64
}

// DescribeIdentityInput carries every field that DescribeIdentity needs.
type DescribeIdentityInput struct {
	IdentityID string
}

// IdentityResult is the transport-agnostic representation of a stored
// identity, consumed by the HTTP serialisers.
type IdentityResult struct {
	ID               string
	CreationDate     time.Time
	LastModifiedDate time.Time
	Logins           map[string]string
}

// GetOpenIdTokenInput carries every field that GetOpenIdToken needs. Logins
// are accepted for wire compatibility and validated, but never persisted.
type GetOpenIdTokenInput struct {
	IdentityID string
	Logins     map[string]string
}

// OpenIdTokenResult is the transport-agnostic OpenID Connect token payload.
type OpenIdTokenResult struct {
	IdentityID string
	Token      string
}

// DeleteIdentitiesInput carries every field that DeleteIdentities needs.
type DeleteIdentitiesInput struct {
	IdentityIDs []string
}

// ListIdentitiesInput carries every field that ListIdentities needs.
type ListIdentitiesInput struct {
	IdentityPoolID     string
	MaxResultsProvided bool
	MaxResults         int
	NextToken          string
	// HideDisabled is accepted for SPEC compliance. Edge identities have no
	// disabled state, so the filter has no effect.
	HideDisabled bool
}

// ListIdentitiesResult is the transport-agnostic page of identities.
type ListIdentitiesResult struct {
	IdentityPoolID string
	Identities     []IdentityResult
	NextToken      string
}

// UnlinkIdentityInput carries every field that UnlinkIdentity needs.
type UnlinkIdentityInput struct {
	IdentityID     string
	LoginsToRemove []string
	Logins         map[string]string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// getIdCore is the single entry point for GetId, shared by the HTTP API and
// any internal caller. It returns the resolved identity ID.
func (s *CognitoIdentityService) getIdCore(reqCtx *request.RequestContext, in GetIdInput) (string, error) {
	if !validateIdentityPoolId(in.IdentityPoolID) {
		return "", ErrInvalidParameter
	}

	// AccountId is used for cross-account identity pool delegation in AWS.
	// It is accepted and validated for SPEC compliance; the edge environment
	// operates single-account so no cross-account check is enforced.
	if in.AccountID != "" {
		if !validateAccountId(in.AccountID) {
			return "", ErrInvalidParameter
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return "", err
	}

	if _, err := store.GetIdentityPool(in.IdentityPoolID); err != nil {
		return "", mapStoreError(err, cognitoidentitystore.ErrIdentityPoolNotFound)
	}

	if !validateMapSize(len(in.Logins), 10) || !validateLoginsKeys(in.Logins) {
		return "", ErrInvalidParameter
	}
	if !validateLoginsValues(in.Logins) {
		return "", ErrInvalidParameter
	}

	// For authenticated identities (Logins provided), AWS reuses the existing
	// identity whose logins match. Only create a new identity when no match
	// exists; the store resolves both paths under the pool's key lock so
	// concurrent callers cannot create duplicate identities.
	identity, err := store.GetOrCreateIdentityByLogins(in.IdentityPoolID, in.Logins)
	if err != nil {
		return "", ErrInternalError
	}

	return identity.ID, nil
}

// getCredentialsForIdentityCore is the single entry point for
// GetCredentialsForIdentity.
func (s *CognitoIdentityService) getCredentialsForIdentityCore(reqCtx *request.RequestContext, in GetCredentialsForIdentityInput) (*CredentialsResult, error) {
	if !validateIdentityId(in.IdentityID) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetIdentityByID(in.IdentityID)
	if err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityNotFound)
	}

	if in.CustomRoleARN != "" && !validateRoleARN(in.CustomRoleARN) {
		return nil, ErrInvalidParameter
	}

	// When the caller provides fresh provider tokens via Logins, persist them
	// onto the identity so that subsequent role selection and credential
	// issuance reflect the current authentication state.
	if len(in.Logins) > 0 {
		if !validateMapSize(len(in.Logins), 10) || !validateLoginsKeys(in.Logins) {
			return nil, ErrInvalidParameter
		}
		if !validateLoginsValues(in.Logins) {
			return nil, ErrInvalidParameter
		}
		if identity.Logins == nil {
			identity.Logins = make(map[string]string)
		}
		for k, v := range in.Logins {
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
	roleArn := in.CustomRoleARN
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

	result, err := s.credentialIssuer.IssueSession(roleArn, in.IdentityID, credentialSessionDurationSeconds)
	if err != nil {
		return nil, ErrInternalError
	}

	return &CredentialsResult{
		IdentityID:      in.IdentityID,
		AccessKeyID:     result.AccessKeyID,
		SecretAccessKey: result.SecretAccessKey,
		SessionToken:    result.SessionToken,
		ExpirationUnix:  result.Expiration.Unix(),
	}, nil
}

// describeIdentityCore is the single entry point for DescribeIdentity.
func (s *CognitoIdentityService) describeIdentityCore(reqCtx *request.RequestContext, in DescribeIdentityInput) (*IdentityResult, error) {
	if !validateIdentityId(in.IdentityID) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetIdentityByID(in.IdentityID)
	if err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityNotFound)
	}

	return &IdentityResult{
		ID:               identity.ID,
		CreationDate:     identity.CreationDate,
		LastModifiedDate: identity.LastModifiedDate,
		Logins:           identity.Logins,
	}, nil
}

// getOpenIdTokenCore is the single entry point for GetOpenIdToken.
func (s *CognitoIdentityService) getOpenIdTokenCore(reqCtx *request.RequestContext, in GetOpenIdTokenInput) (*OpenIdTokenResult, error) {
	if !validateIdentityId(in.IdentityID) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetIdentityByID(in.IdentityID)
	if err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityNotFound)
	}

	// Logins are accepted for wire compatibility but NOT persisted. In AWS,
	// GetOpenIdToken verifies the caller's provider tokens against the actual
	// identity provider before issuing a token. The edge environment cannot
	// perform external provider verification, so the parameter is accepted
	// without side effects to prevent identity takeover via Logins injection.
	if len(in.Logins) > 0 {
		if !validateMapSize(len(in.Logins), 10) || !validateLoginsKeys(in.Logins) {
			return nil, ErrInvalidParameter
		}
		if !validateLoginsValues(in.Logins) {
			return nil, ErrInvalidParameter
		}
	}

	token, err := s.tokenMgr.generateOpenIdToken(in.IdentityID, identity.IdentityPoolID, openIdTokenTTLSeconds, nil, nil)
	if err != nil {
		return nil, ErrInternalError
	}

	return &OpenIdTokenResult{
		IdentityID: in.IdentityID,
		Token:      token,
	}, nil
}

// deleteIdentitiesCore is the single entry point for DeleteIdentities. It
// returns the identity IDs that could not be deleted.
func (s *CognitoIdentityService) deleteIdentitiesCore(reqCtx *request.RequestContext, in DeleteIdentitiesInput) ([]string, error) {
	if len(in.IdentityIDs) == 0 {
		return nil, ErrInvalidParameter
	}
	if len(in.IdentityIDs) > maxIdentityIdsToDelete {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var unprocessed []string
	for _, id := range in.IdentityIDs {
		identity, err := store.GetIdentityByID(id)
		if err != nil {
			unprocessed = append(unprocessed, id)
			continue
		}
		if err := store.DeleteIdentity(identity.IdentityPoolID, id); err != nil {
			unprocessed = append(unprocessed, id)
		}
	}

	return unprocessed, nil
}

// listIdentitiesCore is the single entry point for ListIdentities.
func (s *CognitoIdentityService) listIdentitiesCore(reqCtx *request.RequestContext, in ListIdentitiesInput) (*ListIdentitiesResult, error) {
	if !validateIdentityPoolId(in.IdentityPoolID) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetIdentityPool(in.IdentityPoolID); err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityPoolNotFound)
	}

	if !in.MaxResultsProvided {
		return nil, ErrInvalidParameter
	}
	if !validateQueryLimit(in.MaxResults) {
		return nil, ErrInvalidParameter
	}
	if !validatePaginationKey(in.NextToken) {
		return nil, ErrInvalidParameter
	}

	identities, token, err := store.ListIdentitiesByPool(in.IdentityPoolID, in.MaxResults, in.NextToken)
	if err != nil {
		return nil, ErrInternalError
	}

	result := &ListIdentitiesResult{
		IdentityPoolID: in.IdentityPoolID,
		Identities:     make([]IdentityResult, 0, len(identities)),
	}
	for _, identity := range identities {
		result.Identities = append(result.Identities, IdentityResult{
			ID:               identity.ID,
			CreationDate:     identity.CreationDate,
			LastModifiedDate: identity.LastModifiedDate,
			Logins:           identity.Logins,
		})
	}
	result.NextToken = token

	return result, nil
}

// unlinkIdentityCore is the single entry point for UnlinkIdentity.
func (s *CognitoIdentityService) unlinkIdentityCore(reqCtx *request.RequestContext, in UnlinkIdentityInput) error {
	if !validateIdentityId(in.IdentityID) {
		return ErrInvalidParameter
	}

	if len(in.LoginsToRemove) == 0 {
		return ErrInvalidParameter
	}

	// Logins provides the caller's provider tokens for authorization. AWS
	// requires at least one provider token matching the identity's linked
	// providers before allowing an unlink operation.
	if len(in.Logins) == 0 {
		return ErrNotAuthorized
	}
	if !validateMapSize(len(in.Logins), 10) || !validateLoginsKeys(in.Logins) {
		return ErrInvalidParameter
	}
	if !validateLoginsValues(in.Logins) {
		return ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	identity, err := store.GetIdentityByID(in.IdentityID)
	if err != nil {
		return mapStoreError(err, cognitoidentitystore.ErrIdentityNotFound)
	}

	// Verify the caller holds a token for at least one provider linked to
	// this identity. Both the provider name AND the token value must match
	// to prevent impersonation via provider-name-only knowledge.
	providerMatch := false
	for provider, tokenValue := range in.Logins {
		if storedValue, exists := identity.Logins[provider]; exists && storedValue == tokenValue {
			providerMatch = true
			break
		}
	}
	if !providerMatch {
		return ErrNotAuthorized
	}

	if err := store.UnlinkLogins(identity.IdentityPoolID, in.IdentityID, in.LoginsToRemove); err != nil {
		return mapStoreError(err, cognitoidentitystore.ErrIdentityNotFound)
	}

	return nil
}
