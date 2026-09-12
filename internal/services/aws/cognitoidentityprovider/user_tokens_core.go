package cognitoidentityprovider

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	"vorpalstacks/pkg/vsjwt"
)

// clientTokenValidity holds the effective token-validity periods an
// issuance applies, resolved from the app client's configuration.
type clientTokenValidity struct {
	accessMinutes int
	idMinutes     int
	refreshDays   int
}

// resolveClientTokenValidity reads the client's configured validity periods,
// falling back to the platform defaults for any validity the client leaves
// unset. Both issuance paths — user tokens and machine-to-machine client
// tokens — resolve their periods through here so the fallback semantics
// cannot drift apart.
func resolveClientTokenValidity(client *cognitostore.UserPoolClient) clientTokenValidity {
	v := clientTokenValidity{
		accessMinutes: cognitostore.DefaultAccessTokenValidityMinutes,
		idMinutes:     cognitostore.DefaultIDTokenValidityMinutes,
		refreshDays:   cognitostore.DefaultRefreshTokenValidityDays,
	}
	if client.AccessTokenValidity > 0 {
		v.accessMinutes = client.AccessTokenValidity
	}
	if client.IDTokenValidity > 0 {
		v.idMinutes = client.IDTokenValidity
	}
	if client.RefreshTokenValidity > 0 {
		v.refreshDays = client.RefreshTokenValidity
	}
	return v
}

// CreateTokens creates access, ID, and refresh tokens for the specified user.
// The triggerSource parameter controls which PreTokenGeneration Lambda event
// fires (TokenGenerationAuthentication, TokenGenerationRefreshTokens, or
// TokenGenerationHostedAuth). The clientMetadata parameter is forwarded to
// the trigger so Lambda can customise claims.
func (s *CognitoService) CreateTokens(reqCtx *request.RequestContext, userPoolID, userID, clientID, triggerSource string, clientMetadata map[string]string) (accessToken, idToken, refreshToken string, expiresIn int64, err error) {
	return s.createTokens(reqCtx, userPoolID, userID, clientID, triggerSource, clientMetadata, nil, "")
}

// CreateTokensWithIDClaims issues the same token set as CreateTokens with
// additional ID-token claims. The hosted-UI authorisation flow binds the
// OIDC nonce into the ID token through this hook; nil idTokenClaims issues
// the standard claim set.
func (s *CognitoService) CreateTokensWithIDClaims(reqCtx *request.RequestContext, userPoolID, userID, clientID, triggerSource string, clientMetadata map[string]string, idTokenClaims map[string]interface{}) (accessToken, idToken, refreshToken string, expiresIn int64, err error) {
	return s.createTokens(reqCtx, userPoolID, userID, clientID, triggerSource, clientMetadata, idTokenClaims, "")
}

// CreateTokensWithSessionScope issues the same token set as CreateTokens,
// reissuing the scope of the session that minted the refresh token: a
// refresh grant — whichever plane carries it — neither widens nor narrows
// the originally granted scope.
func (s *CognitoService) CreateTokensWithSessionScope(reqCtx *request.RequestContext, userPoolID, userID, clientID, triggerSource, sessionScope string, clientMetadata map[string]string) (accessToken, idToken, refreshToken string, expiresIn int64, err error) {
	return s.createTokens(reqCtx, userPoolID, userID, clientID, triggerSource, clientMetadata, nil, sessionScope)
}

// apiSignInScope is the scope claim of an access token issued through the
// Amazon Cognito user pools API: per the token documentation, "A token from
// Amazon Cognito API sign-in only contains the scope
// aws.cognito.signin.user.admin."
const apiSignInScope = "aws.cognito.signin.user.admin"

// createTokens is the single token-issuance path. The scope carried on the
// access and refresh tokens follows the session's origin (see the scope
// derivation below); sessionScope, when non-empty, reissues a stored
// session's original scope.
func (s *CognitoService) createTokens(reqCtx *request.RequestContext, userPoolID, userID, clientID, triggerSource string, clientMetadata map[string]string, idTokenClaims map[string]interface{}, sessionScope string) (accessToken, idToken, refreshToken string, expiresIn int64, err error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return "", "", "", 0, err
	}
	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("user pool not found: %w", err)
	}

	privateKey, err := vsjwt.DecodePrivateKeyFromPEM(userPool.JwtPrivateKey)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to decode JWT private key: %w", err)
	}

	issuer := fmt.Sprintf("https://%s/%s", cognitoIdpHost(reqCtx.GetRegion()), userPoolID)
	jwtManager, err := vsjwt.NewManager(privateKey, userPool.JwtKeyID, issuer)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to create JWT manager: %w", err)
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to load user for token issuance: %w", err)
	}
	if user == nil {
		return "", "", "", 0, fmt.Errorf("user not found: %s", userID)
	}

	attrs := userAttributesMap(user)
	// The ID token derives its IAM role claims from the user's groups: the
	// roles of every group carrying one, with the highest-precedence group's
	// role as cognito:preferred_role — the claims identity-pool Token role
	// mappings read from the token. A PreTokenGeneration override below
	// replaces the derived set.
	if user.Attributes == nil {
		user.Attributes = make(map[string]string)
	}
	groupRoles, preferredRole, err := groupRoleClaims(store, userPoolID, user.Groups)
	if err != nil {
		return "", "", "", 0, ErrInternalError
	}
	if len(groupRoles) > 0 {
		user.Attributes["cognito:roles"] = strings.Join(groupRoles, ",")
		if preferredRole != "" {
			user.Attributes["cognito:preferred_role"] = preferredRole
		}
	}
	tokGenResult, tokGenErr := invokePreTokenGeneration(
		reqCtx, s, triggerSource, userPoolID, user.Username, clientID,
		userPool.LambdaConfig, attrs, user.Groups, clientMetadata,
	)
	if tokGenErr != nil {
		return "", "", "", 0, fmt.Errorf("pre-token-generation trigger failed: %w", tokGenErr)
	}

	if tokGenResult != nil {
		// A user without stored attributes still carries a nil map after
		// the JSON round-trip; claim application writes into it.
		if user.Attributes == nil {
			user.Attributes = make(map[string]string)
		}
		for k, v := range tokGenResult.ClaimsToAddOrOverride {
			user.Attributes[k] = v
		}
		for _, k := range tokGenResult.ClaimsToSuppress {
			delete(user.Attributes, k)
		}
		if len(tokGenResult.GroupsToOverride) > 0 {
			user.Groups = tokGenResult.GroupsToOverride
		}
		// The role members of groupOverrideDetails carry into the same
		// cognito:roles / cognito:preferred_role claims AWS derives from
		// group role ARNs; here the trigger is the sole source, applied
		// when it supplies them. The mutations are in-memory — they shape
		// this token's claims without persisting to the user record.
		if len(tokGenResult.IAMRolesToOverride) > 0 {
			user.Attributes["cognito:roles"] = strings.Join(tokGenResult.IAMRolesToOverride, ",")
		}
		if tokGenResult.PreferredRole != "" {
			user.Attributes["cognito:preferred_role"] = tokGenResult.PreferredRole
		}
	}

	// Determine token validity from client configuration. An issuance must
	// fail closed when the client record cannot be read — defaulting to the
	// standard validity periods would mint tokens that violate the client's
	// configuration. The one client fetch also answers the scope carried on
	// the access and refresh tokens below — reading the record twice would
	// double the cost of every token issuance.
	client, err := store.GetUserPoolClient(userPoolID, clientID)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to load user pool client for token issuance: %w", err)
	}
	validity := resolveClientTokenValidity(client)
	atValiditySec := int64(validity.accessMinutes) * 60
	idValiditySec := int64(validity.idMinutes) * 60
	atExpiry := time.Duration(validity.accessMinutes) * time.Minute
	idExpiry := time.Duration(validity.idMinutes) * time.Minute
	rtExpiry := time.Duration(validity.refreshDays) * 24 * time.Hour

	// The access token's scope claim follows the session's origin. A token
	// from API sign-in carries only the user self-service scope; a hosted-UI
	// session carries the client's granted OAuth scopes (openid as the
	// default, the token endpoint's no-scope behaviour). A refresh reissues
	// the scope of the session that minted the refresh token, carried in as
	// sessionScope. The access and refresh token records store the same
	// granted scope.
	scope := apiSignInScope
	if triggerSource == TokenGenerationHostedAuth {
		scope = "openid"
		if len(client.AllowedOAuthScopes) > 0 {
			scope = strings.Join(client.AllowedOAuthScopes, " ")
		}
	}
	if sessionScope != "" {
		scope = sessionScope
	}

	accessToken, err = jwtManager.GenerateAccessTokenWithScope(user, clientID, scope, atValiditySec)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to generate access token: %w", err)
	}

	at := cognitostore.NewAccessToken(userPoolID, user.ID, clientID, scope, time.Now().Add(atExpiry))
	at.Token = accessToken
	if err := store.CreateAccessToken(at); err != nil {
		return "", "", "", 0, fmt.Errorf("failed to store access token: %w", err)
	}

	idToken, err = jwtManager.GenerateIDTokenWithClaims(user, clientID, idValiditySec, idTokenClaims)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to generate ID token: %w", err)
	}

	// Store ID token so it can be validated and revoked.
	it := cognitostore.NewIDToken(userPoolID, user.ID, clientID, "", time.Now().Add(idExpiry), user.Groups)
	it.Token = idToken
	if err := store.CreateIDToken(it); err != nil {
		return "", "", "", 0, fmt.Errorf("failed to store ID token: %w", err)
	}

	refreshToken = jwtManager.GenerateRefreshToken()

	rt := cognitostore.NewRefreshToken(userPoolID, user.ID, clientID, scope, time.Now().Add(rtExpiry))
	rt.Token = refreshToken
	if err := store.CreateRefreshToken(rt); err != nil {
		return "", "", "", 0, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return accessToken, idToken, refreshToken, atValiditySec, nil
}

// ValidateAccessToken validates an access token and returns the user ID.
func (s *CognitoService) ValidateAccessToken(reqCtx *request.RequestContext, tokenString string) (string, error) {
	record, err := s.validateAccessTokenRecord(reqCtx, tokenString)
	if err != nil {
		return "", err
	}
	return record.UserID, nil
}

// validateAccessTokenRecord validates an access token and returns the stored
// token record, whose ClientID pins the app client that was issued the token.
func (s *CognitoService) validateAccessTokenRecord(reqCtx *request.RequestContext, tokenString string) (*cognitostore.AccessToken, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// The stored token record pins the owning pool, so the signature can be
	// verified against that pool's key alone instead of scanning every pool.
	record, err := store.GetAccessTokenByValue(tokenString)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	pool, err := store.GetUserPool(record.UserPoolID)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// A disabled user's tokens stop authenticating immediately, even ones
	// minted before the disable: AdminDisableUser deactivates the profile
	// and revokes its access.
	user, err := store.GetUserByID(record.UserID)
	if err != nil || !user.Enabled {
		return nil, ErrNotAuthorized
	}

	if _, err := validateAccessTokenSignature(pool, reqCtx.GetRegion(), tokenString); err != nil {
		return nil, err
	}
	return record, nil
}

// validateAccessTokenSignature verifies an access token against a specific
// pool's signing key and returns the subject (user ID).
func validateAccessTokenSignature(pool *cognitostore.UserPool, region, tokenString string) (string, error) {
	publicKey, err := vsjwt.DecodePublicKeyFromPEM(pool.JwtPublicKey)
	if err != nil {
		return "", ErrNotAuthorized
	}

	issuer := fmt.Sprintf("https://%s/%s", cognitoIdpHost(region), pool.ID)
	jwtManager, err := vsjwt.NewManagerWithPublicKey(publicKey, pool.JwtKeyID, issuer)
	if err != nil {
		return "", ErrNotAuthorized
	}

	claims, err := jwtManager.ValidateTokenForUse(tokenString, "access", "")
	if err != nil {
		return "", ErrNotAuthorized
	}

	return claims.Subject, nil
}

// ValidateTokenForPool validates a Cognito access token for a specific
// user pool, designed for cross-service consumers (e.g. API Gateway
// COGNITO_USER_POOLS authorizer) via the eventbus CognitoTokenValidator
// interface.
func (s *CognitoService) ValidateTokenForPool(ctx context.Context, region, userPoolID, accessToken string) (string, error) {
	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return "", err
	}

	pool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return "", ErrResourceNotFound
	}

	// Verify the token has not been revoked.
	record, err := store.GetAccessTokenByValue(accessToken)
	if err != nil {
		return "", ErrNotAuthorized
	}

	// A client-principal token (machine-to-machine issuance) carries no user
	// record: an empty user ID skips the user fetch and its enabled check,
	// leaving the signature validation below as the authorisation.
	if record.UserID != "" {
		// A disabled user's tokens no longer authorise, even for cross-service
		// consumers such as API Gateway COGNITO_USER_POOLS authorisers.
		user, err := store.GetUserByID(record.UserID)
		if err != nil || !user.Enabled {
			return "", ErrNotAuthorized
		}
	}

	return validateAccessTokenSignature(pool, region, accessToken)
}

// IDTokenClaimsForPool validates a Cognito user-pool ID token for a specific
// pool and returns the string claims identity-pool role mappings resolve
// against — the cross-service CognitoIDTokenClaimResolver contract. The
// signature, issuer, expiry and token_use=id are verified against the pool's
// own key; the returned set carries the token's custom claims (including
// cognito:roles and cognito:preferred_role) plus the standard string
// members — aud and sub included, feeding the principal-tag default
// mappings (the app client ID and the user ID). Non-string claim values are
// skipped: role mappings match strings.
func (s *CognitoService) IDTokenClaimsForPool(ctx context.Context, region, userPoolID, idToken string) (map[string]string, error) {
	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}
	pool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	publicKey, err := vsjwt.DecodePublicKeyFromPEM(pool.JwtPublicKey)
	if err != nil {
		return nil, ErrNotAuthorized
	}
	issuer := fmt.Sprintf("https://%s/%s", cognitoIdpHost(region), userPoolID)
	jwtManager, err := vsjwt.NewManagerWithPublicKey(publicKey, pool.JwtKeyID, issuer)
	if err != nil {
		return nil, ErrNotAuthorized
	}
	claims, err := jwtManager.ValidateTokenForUse(idToken, "id", "")
	if err != nil {
		return nil, ErrNotAuthorized
	}
	result := make(map[string]string, len(claims.CustomClaims())+3)
	for k, v := range claims.CustomClaims() {
		if str, ok := v.(string); ok {
			result[k] = str
		}
	}
	if claims.Username != "" {
		result["cognito:username"] = claims.Username
	}
	// An ID token carries its app client ID as the registered aud claim —
	// the client_id member belongs to the access-token shape and stays empty
	// on ID tokens.
	if len(claims.Audience) > 0 {
		result["aud"] = claims.Audience[0]
	}
	if claims.Subject != "" {
		result["sub"] = claims.Subject
	}
	if claims.Email != "" {
		result["email"] = claims.Email
	}
	if len(claims.Groups) > 0 {
		result["cognito:groups"] = strings.Join(claims.Groups, ",")
	}
	return result, nil
}

// groupRoleClaims derives the IAM role claims the ID token carries from the
// user's groups: every group with a role ARN contributes to cognito:roles,
// and the highest-precedence group's role becomes cognito:preferred_role.
// The precedence order is ascending numeric value — the pool's group
// precedence contract — and a group without a precedence value cannot
// nominate the preferred role. A group that no longer exists is skipped
// (deletion cascades membership removal); any other read failure is
// returned so token minting fails closed instead of silently issuing a
// token without its role claims.
func groupRoleClaims(store cognitostore.CognitoStoreInterface, userPoolID string, groups []string) (roles []string, preferredRole string, err error) {
	type groupRole struct {
		precedence  int
		hasPriority bool
		role        string
	}
	var collected []groupRole
	for _, name := range groups {
		group, gerr := store.GetGroup(userPoolID, name)
		if errors.Is(gerr, cognitostore.ErrGroupNotFound) {
			continue
		}
		if gerr != nil {
			return nil, "", gerr
		}
		if group.RoleArn == "" {
			continue
		}
		entry := groupRole{role: group.RoleArn}
		if group.Precedence != nil {
			entry.precedence = *group.Precedence
			entry.hasPriority = true
		}
		collected = append(collected, entry)
	}
	if len(collected) == 0 {
		return nil, "", nil
	}
	sort.SliceStable(collected, func(i, j int) bool {
		// Groups with a precedence rank ahead of those without; equal or
		// absent precedence keeps the membership order.
		if collected[i].hasPriority != collected[j].hasPriority {
			return collected[i].hasPriority
		}
		return collected[i].precedence < collected[j].precedence
	})
	for _, entry := range collected {
		roles = append(roles, entry.role)
	}
	if collected[0].hasPriority {
		preferredRole = collected[0].role
	}
	return roles, preferredRole, nil
}
