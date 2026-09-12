package cognitoidentity

import (
	"errors"

	"vorpalstacks/internal/common/request"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input/Result DTOs — developer identity operations
// ---------------------------------------------------------------------------

// GetOpenIdTokenForDeveloperIdentityInput carries every field that
// GetOpenIdTokenForDeveloperIdentity needs. PrincipalTagsProvided/Raw carry
// the presence flag and untyped wire value of the PrincipalTags member so the
// Core can type-check and validate it at the exact position the wire contract
// requires.
type GetOpenIdTokenForDeveloperIdentityInput struct {
	IdentityPoolID        string
	IdentityID            string
	Logins                map[string]string
	TokenDurationProvided bool
	TokenDuration         int
	PrincipalTagsProvided bool
	PrincipalTagsRaw      interface{}
}

// LookupDeveloperIdentityInput carries every field that
// LookupDeveloperIdentity needs.
type LookupDeveloperIdentityInput struct {
	IdentityPoolID          string
	IdentityID              string
	DeveloperUserIdentifier string
	MaxResultsProvided      bool
	MaxResults              int
	NextToken               string
}

// LookupDeveloperIdentityResult is the transport-agnostic lookup outcome.
type LookupDeveloperIdentityResult struct {
	MatchedIdentityID string
	DeveloperUserIDs  []string
	NextToken         string
}

// MergeDeveloperIdentitiesInput carries every field that
// MergeDeveloperIdentities needs.
type MergeDeveloperIdentitiesInput struct {
	IdentityPoolID            string
	DeveloperProviderName     string
	SourceUserIdentifier      string
	DestinationUserIdentifier string
}

// UnlinkDeveloperIdentityInput carries every field that
// UnlinkDeveloperIdentity needs.
type UnlinkDeveloperIdentityInput struct {
	IdentityID              string
	IdentityPoolID          string
	DeveloperProviderName   string
	DeveloperUserIdentifier string
}

// GetPrincipalTagAttributeMapInput carries every field that
// GetPrincipalTagAttributeMap needs.
type GetPrincipalTagAttributeMapInput struct {
	IdentityPoolID       string
	IdentityProviderName string
}

// SetPrincipalTagAttributeMapInput carries every field that
// SetPrincipalTagAttributeMap needs.
type SetPrincipalTagAttributeMapInput struct {
	IdentityPoolID       string
	IdentityProviderName string
	PrincipalTags        map[string]string
	UseDefaults          bool
}

// PrincipalTagAttributeMapResult is the transport-agnostic principal tag
// attribute map for an identity provider.
type PrincipalTagAttributeMapResult struct {
	IdentityPoolID       string
	IdentityProviderName string
	PrincipalTags        map[string]string
	UseDefaults          bool
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// getOpenIdTokenForDeveloperIdentityCore is the single entry point for
// GetOpenIdTokenForDeveloperIdentity: it registers (or retrieves) the
// developer identity and issues the OpenID Connect token for it.
func (s *CognitoIdentityService) getOpenIdTokenForDeveloperIdentityCore(reqCtx *request.RequestContext, in GetOpenIdTokenForDeveloperIdentityInput) (*OpenIdTokenResult, error) {
	if !validateIdentityPoolId(in.IdentityPoolID) {
		return nil, ErrInvalidParameter
	}

	if !validateMapSize(len(in.Logins), maxLoginsPerRequest) || !validateLoginsKeys(in.Logins) {
		return nil, ErrInvalidParameter
	}
	if !validateLoginsValues(in.Logins) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pool, err := store.GetIdentityPool(in.IdentityPoolID)
	if err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityPoolNotFound)
	}

	// The model documents Logins as public-provider pairs plus exactly one
	// developer-provider pair keyed by the pool's configured domain ("You
	// can only specify one developer provider as part of the Logins map,
	// which is linked to the identity pool"), so the domain must be
	// configured and present among the logins.
	if pool.DeveloperProviderName == "" {
		return nil, ErrInvalidParameter
	}
	devUserID, isDeveloper := in.Logins[pool.DeveloperProviderName]
	if !isDeveloper {
		return nil, ErrInvalidParameter
	}
	publicLogins := make(map[string]string, len(in.Logins)-1)
	for provider, token := range in.Logins {
		if provider != pool.DeveloperProviderName {
			publicLogins[provider] = token
		}
	}

	identityID := in.IdentityID
	if identityID != "" && !validateIdentityId(identityID) {
		return nil, ErrInvalidParameter
	}

	// TokenDuration controls token expiry (Smithy TokenDuration range).
	tokenDuration := int64(developerTokenDefaultTTLSeconds)
	if in.TokenDurationProvided {
		td := int64(in.TokenDuration)
		if !validateTokenDuration(td) {
			return nil, ErrInvalidParameter
		}
		tokenDuration = td
	}

	// PrincipalTags are embedded into the JWT as cognito:principal_tags
	// so that STS AssumeRoleWithWebIdentity propagates them as session tags.
	var principalTags map[string]string
	if in.PrincipalTagsProvided {
		if ptMap, ok := in.PrincipalTagsRaw.(map[string]interface{}); ok {
			if !validateMapSize(len(ptMap), maxPrincipalTags) {
				return nil, ErrInvalidParameter
			}
			principalTags = make(map[string]string, len(ptMap))
			for k, v := range ptMap {
				if !validatePrincipalTagName(k) {
					return nil, ErrInvalidParameter
				}
				value, ok := v.(string)
				if !ok || !validatePrincipalTagValue(value) {
					return nil, ErrInvalidParameter
				}
				principalTags[k] = value
			}
		}
	}

	// The store resolves the developer identity under its key lock: an
	// existing link is reused (a differing supplied IdentityId maps to
	// DeveloperUserAlreadyRegisteredException), otherwise a fresh identity
	// is created and linked in one critical section.
	resolved, err := store.EnsureDeveloperIdentity(in.IdentityPoolID, pool.DeveloperProviderName, devUserID, identityID)
	if err != nil {
		if errors.Is(err, cognitoidentitystore.ErrDeveloperIdentityConflict) {
			return nil, ErrDeveloperUserAlreadyRegistered
		}
		// Both a missing pool and a missing identity surface as the
		// operation's ResourceNotFoundException; anything else is a storage
		// failure.
		if errors.Is(err, cognitoidentitystore.ErrIdentityNotFound) ||
			errors.Is(err, cognitoidentitystore.ErrIdentityPoolNotFound) {
			return nil, ErrResourceNotFound
		}
		return nil, ErrInternalError
	}
	identityID = resolved

	// Public-provider entries are linked to the resolved identity — the
	// implicit linked account the model documents for multiple logins. A
	// login already linked to another identity is a conflict.
	if len(publicLogins) > 0 {
		if err := store.LinkLogins(in.IdentityPoolID, identityID, publicLogins); err != nil {
			if errors.Is(err, cognitoidentitystore.ErrLoginConflict) {
				return nil, ErrResourceConflict
			}
			if errors.Is(err, cognitoidentitystore.ErrIdentityNotFound) {
				return nil, ErrResourceNotFound
			}
			return nil, ErrInternalError
		}
	}

	token, err := s.tokenMgr.generateOpenIdToken(identityID, in.IdentityPoolID, tokenDuration, nil, principalTags)
	if err != nil {
		return nil, ErrInternalError
	}

	return &OpenIdTokenResult{
		IdentityID: identityID,
		Token:      token,
	}, nil
}

// lookupDeveloperIdentityCore is the single entry point for
// LookupDeveloperIdentity.
func (s *CognitoIdentityService) lookupDeveloperIdentityCore(reqCtx *request.RequestContext, in LookupDeveloperIdentityInput) (*LookupDeveloperIdentityResult, error) {
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

	if in.IdentityID != "" && !validateIdentityId(in.IdentityID) {
		return nil, ErrInvalidParameter
	}
	if in.IdentityID == "" && in.DeveloperUserIdentifier == "" {
		return nil, ErrInvalidParameter
	}
	if in.DeveloperUserIdentifier != "" && !validateDeveloperUserIdentifier(in.DeveloperUserIdentifier) {
		return nil, ErrInvalidParameter
	}
	maxResults := defaultLookupMaxResults
	if in.MaxResultsProvided {
		if !validateQueryLimit(in.MaxResults) {
			return nil, ErrInvalidParameter
		}
		maxResults = in.MaxResults
	}
	if !validatePaginationKey(in.NextToken) {
		return nil, ErrInvalidParameter
	}

	matchedIdentityID, devUserIDs, nextTokenOut, err := store.LookupDeveloperIdentity(in.IdentityPoolID, in.IdentityID, in.DeveloperUserIdentifier, maxResults, in.NextToken)
	if err != nil {
		return nil, ErrInternalError
	}

	return &LookupDeveloperIdentityResult{
		MatchedIdentityID: matchedIdentityID,
		DeveloperUserIDs:  devUserIDs,
		NextToken:         nextTokenOut,
	}, nil
}

// mergeDeveloperIdentitiesCore is the single entry point for
// MergeDeveloperIdentities. It returns the surviving destination identity ID.
func (s *CognitoIdentityService) mergeDeveloperIdentitiesCore(reqCtx *request.RequestContext, in MergeDeveloperIdentitiesInput) (string, error) {
	if !validateIdentityPoolId(in.IdentityPoolID) {
		return "", ErrInvalidParameter
	}
	if in.DeveloperProviderName == "" || !validateDeveloperProviderName(in.DeveloperProviderName) {
		return "", ErrInvalidParameter
	}
	if in.SourceUserIdentifier == "" || !validateDeveloperUserIdentifier(in.SourceUserIdentifier) {
		return "", ErrInvalidParameter
	}
	if in.DestinationUserIdentifier == "" || !validateDeveloperUserIdentifier(in.DestinationUserIdentifier) {
		return "", ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return "", err
	}

	// The store performs the merge under the pool lock with the developer
	// identity link moving before any identity record is destroyed.
	destIdentityID, err := store.MergeDeveloperIdentities(in.IdentityPoolID, in.DeveloperProviderName, in.SourceUserIdentifier, in.DestinationUserIdentifier)
	if err != nil {
		return "", mapStoreError(err, cognitoidentitystore.ErrIdentityNotFound)
	}

	return destIdentityID, nil
}

// unlinkDeveloperIdentityCore is the single entry point for
// UnlinkDeveloperIdentity.
func (s *CognitoIdentityService) unlinkDeveloperIdentityCore(reqCtx *request.RequestContext, in UnlinkDeveloperIdentityInput) error {
	if !validateIdentityId(in.IdentityID) {
		return ErrInvalidParameter
	}
	if !validateIdentityPoolId(in.IdentityPoolID) {
		return ErrInvalidParameter
	}
	if in.DeveloperProviderName == "" || !validateDeveloperProviderName(in.DeveloperProviderName) {
		return ErrInvalidParameter
	}
	if in.DeveloperUserIdentifier == "" || !validateDeveloperUserIdentifier(in.DeveloperUserIdentifier) {
		return ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	if err := store.UnlinkDeveloperIdentity(in.IdentityPoolID, in.DeveloperProviderName, in.DeveloperUserIdentifier); err != nil {
		return mapStoreError(err, cognitoidentitystore.ErrIdentityNotFound)
	}

	return nil
}

// getPrincipalTagAttributeMapCore is the single entry point for
// GetPrincipalTagAttributeMap. A provider without a stored map reports the
// AWS default behaviour (empty tags with UseDefaults=true).
func (s *CognitoIdentityService) getPrincipalTagAttributeMapCore(reqCtx *request.RequestContext, in GetPrincipalTagAttributeMapInput) (*PrincipalTagAttributeMapResult, error) {
	if !validateIdentityPoolId(in.IdentityPoolID) {
		return nil, ErrInvalidParameter
	}
	if in.IdentityProviderName == "" || !validateIdentityProviderNameLength(in.IdentityProviderName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetIdentityPool(in.IdentityPoolID); err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityPoolNotFound)
	}

	ptam, err := store.GetPrincipalTagAttributeMap(in.IdentityPoolID, in.IdentityProviderName)
	if err != nil {
		if !errors.Is(err, cognitoidentitystore.ErrIdentityNotFound) {
			return nil, ErrInternalError
		}
		return &PrincipalTagAttributeMapResult{
			IdentityPoolID:       in.IdentityPoolID,
			IdentityProviderName: in.IdentityProviderName,
			PrincipalTags:        map[string]string{},
			UseDefaults:          true,
		}, nil
	}

	return &PrincipalTagAttributeMapResult{
		IdentityPoolID:       ptam.IdentityPoolID,
		IdentityProviderName: ptam.IdentityProviderName,
		PrincipalTags:        ptam.PrincipalTags,
		UseDefaults:          ptam.UseDefaults,
	}, nil
}

// setPrincipalTagAttributeMapCore is the single entry point for
// SetPrincipalTagAttributeMap. It returns the stored attribute map.
func (s *CognitoIdentityService) setPrincipalTagAttributeMapCore(reqCtx *request.RequestContext, in SetPrincipalTagAttributeMapInput) (*PrincipalTagAttributeMapResult, error) {
	if !validateIdentityPoolId(in.IdentityPoolID) {
		return nil, ErrInvalidParameter
	}
	if in.IdentityProviderName == "" || !validateIdentityProviderNameLength(in.IdentityProviderName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetIdentityPool(in.IdentityPoolID); err != nil {
		return nil, mapStoreError(err, cognitoidentitystore.ErrIdentityPoolNotFound)
	}

	if len(in.PrincipalTags) > 0 {
		if !validateMapSize(len(in.PrincipalTags), maxPrincipalTags) {
			return nil, ErrInvalidParameter
		}
		for k, v := range in.PrincipalTags {
			if !validatePrincipalTagName(k) || !validatePrincipalTagValue(v) {
				return nil, ErrInvalidParameter
			}
		}
	}

	if err := store.SetPrincipalTagAttributeMap(in.IdentityPoolID, in.IdentityProviderName, in.PrincipalTags, in.UseDefaults); err != nil {
		return nil, ErrInternalError
	}

	return &PrincipalTagAttributeMapResult{
		IdentityPoolID:       in.IdentityPoolID,
		IdentityProviderName: in.IdentityProviderName,
		PrincipalTags:        in.PrincipalTags,
		UseDefaults:          in.UseDefaults,
	}, nil
}
