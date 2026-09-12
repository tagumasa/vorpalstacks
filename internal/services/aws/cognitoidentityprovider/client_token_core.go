package cognitoidentityprovider

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	"vorpalstacks/pkg/vsjwt"
)

// Core for the machine-to-machine (M2M) token family. GetClientToken issues
// an access token that authorises an application rather than a user: the app
// client presents its ID and secret, and receives a token carrying the custom
// resource-server scopes it requested.

// Smithy ClientSecretType: length {min: 24, max: 64}, pattern ^[\w+]+$.
const (
	clientSecretMinLength = 24
	clientSecretMaxLength = 64
)

// clientSecretPattern is the Smithy ClientSecretType pattern.
var clientSecretPattern = regexp.MustCompile(`^[\w+]+$`)

// maxClientTokenScopes is the Smithy ScopeListType length bound {min: 0,
// max: 50}.
const maxClientTokenScopes = 50

// clientTokenAuthFlow is the ExplicitAuthFlows value an app client must
// carry — and carry alone — to request machine-to-machine tokens. The AWS
// developer guide records the value; the vendored model enum lags it.
const clientTokenAuthFlow = "ALLOW_CLIENT_TOKEN_AUTH"

// GetClientTokenInput carries the wire parameters of GetClientToken.
type GetClientTokenInput struct {
	ClientID       string
	Secret         string
	Scopes         []string
	ClientMetadata map[string]string
}

// m2mTokenPrincipal adapts an app client to the vsjwt JWTUser contract: the
// machine identity's subject is the client itself, and it carries no user
// attributes or groups.
type m2mTokenPrincipal struct {
	clientID string
}

func (p m2mTokenPrincipal) GetID() string                           { return p.clientID }
func (p m2mTokenPrincipal) GetUsername() string                     { return "" }
func (p m2mTokenPrincipal) GetGroups() []string                     { return nil }
func (p m2mTokenPrincipal) GetEmail() string                        { return "" }
func (p m2mTokenPrincipal) GetEmailVerified() bool                  { return false }
func (p m2mTokenPrincipal) GetCustomClaims() map[string]interface{} { return nil }

// clientSecretMatches reports whether the presented secret equals the
// client's primary secret or one of its active secret descriptors. The
// comparison is constant-time per candidate.
func clientSecretMatches(client *cognitostore.UserPoolClient, secret string) bool {
	if client.ClientSecret != "" &&
		subtle.ConstantTimeCompare([]byte(secret), []byte(client.ClientSecret)) == 1 {
		return true
	}
	for _, d := range client.ClientSecrets {
		if subtle.ConstantTimeCompare([]byte(secret), []byte(d.ClientSecretValue)) == 1 {
			return true
		}
	}
	return false
}

// secretHashFor computes the client proof the AWS developer guide
// documents ("Computing secret hash values"):
// Base64 ( HMAC_SHA256 ( "Client Secret Key", "Username" + "Client Id" ) ).
// The username element is whatever sign-in identifier the client presented
// — the guide permits any user pool sign-in attribute there for sign-in
// operations.
func secretHashFor(secret, username, clientID string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(username + clientID))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// verifyClientSecretHash enforces the model's SECRET_HASH contract: an app
// client that holds a secret must present SECRET_HASH — matching one of its
// active secrets — in the auth parameters of every InitiateAuth flow and in
// every challenge response ("You must provide a SECRET_HASH parameter in
// all challenge responses to an app client that has a client secret"); a
// client without a secret is exempt. Absence or mismatch is
// NotAuthorizedException, per the model's NotAuthorizedException note.
func verifyClientSecretHash(client *cognitostore.UserPoolClient, username, provided string) error {
	if client == nil || (client.ClientSecret == "" && len(client.ClientSecrets) == 0) {
		return nil
	}
	if username == "" || provided == "" {
		return ErrNotAuthorized
	}
	if client.ClientSecret != "" &&
		subtle.ConstantTimeCompare([]byte(provided), []byte(secretHashFor(client.ClientSecret, username, client.ClientID))) == 1 {
		return nil
	}
	for _, d := range client.ClientSecrets {
		if subtle.ConstantTimeCompare([]byte(provided), []byte(secretHashFor(d.ClientSecretValue, username, client.ClientID))) == 1 {
			return nil
		}
	}
	return ErrNotAuthorized
}

// verifyAuthPlaneSecretHash resolves the app client behind clientID and
// enforces its SECRET_HASH proof on the authentication planes. A client
// that cannot be resolved holds no secret to enforce — the calling flow's
// own client handling produces its error immediately after.
func (s *CognitoService) verifyAuthPlaneSecretHash(reqCtx *request.RequestContext, clientID, username, provided string) error {
	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}
	pool, err := store.GetUserPoolByClientID(clientID)
	if err != nil {
		return nil
	}
	client, err := store.GetUserPoolClient(pool.ID, clientID)
	if err != nil {
		return nil
	}
	return verifyClientSecretHash(client, username, provided)
}

// activeResourceServerScopes returns the set of custom scopes the pool's
// resource servers currently define, each in its full
// resource-server-identifier/scope-name form. A scope deleted from its
// resource server leaves the set, so a client still associating it requests
// an inactive scope: the request proceeds without it.
func activeResourceServerScopes(store cognitostore.CognitoStoreInterface, userPoolID string) (map[string]bool, error) {
	servers, err := store.ListResourceServers(userPoolID)
	if err != nil {
		return nil, err
	}
	scopes := make(map[string]bool)
	for _, rs := range servers {
		for _, sc := range rs.Scopes {
			scopes[rs.Identifier+"/"+sc.ScopeName] = true
		}
	}
	return scopes, nil
}

// getClientTokenCore issues an M2M access token for a confidential app
// client. The client is identified by its ID alone — the operation takes no
// UserPoolId — and authorises with an active client secret. Granted scopes
// are custom resource-server scopes: an omitted Scopes member grants the
// client's configured scopes, a requested scope the client does not
// associate fails the authentication, and a scope no longer defined on any
// resource server is skipped as inactive.
func (s *CognitoService) getClientTokenCore(ctx context.Context, reqCtx *request.RequestContext, in GetClientTokenInput) (interface{}, error) {
	// GetClientTokenRequest marks ClientId and Secret required; the secret
	// carries the ClientSecretType length and pattern.
	if in.ClientID == "" || in.Secret == "" {
		return nil, ErrInvalidParameter
	}
	if len(in.Secret) < clientSecretMinLength || len(in.Secret) > clientSecretMaxLength ||
		!clientSecretPattern.MatchString(in.Secret) {
		return nil, ErrInvalidParameter
	}
	if len(in.Scopes) > maxClientTokenScopes {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// An unknown client maps to the operation's InvalidClientException; a
	// storage failure behind the index lookup is an internal error, exactly
	// as the client read below distinguishes them.
	pool, err := store.GetUserPoolByClientID(in.ClientID)
	if err != nil {
		if errors.Is(err, cognitostore.ErrClientNotFound) || errors.Is(err, cognitostore.ErrUserPoolNotFound) {
			return nil, ErrClientNotFound
		}
		return nil, ErrInternalError
	}
	client, err := store.GetUserPoolClient(pool.ID, in.ClientID)
	if err != nil {
		if errors.Is(err, cognitostore.ErrClientNotFound) {
			return nil, ErrClientNotFound
		}
		return nil, ErrInternalError
	}

	// The operation is enabled only for a client configured for
	// machine-to-machine authorisation: it must hold a secret and carry the
	// client-token flow as its only authentication flow.
	hasSecret := client.ClientSecret != "" || len(client.ClientSecrets) > 0
	flowEnabled := len(client.ExplicitAuthFlows) == 1 && client.ExplicitAuthFlows[0] == clientTokenAuthFlow
	if !hasSecret || !flowEnabled {
		return nil, ErrOperationNotEnabled
	}
	if !clientSecretMatches(client, in.Secret) {
		return nil, ErrNotAuthorized
	}

	activeScopes, err := activeResourceServerScopes(store, pool.ID)
	if err != nil {
		return nil, ErrInternalError
	}

	var granted []string
	if len(in.Scopes) == 0 {
		// The scopes configured for the app client, narrowed to the custom
		// scopes its resource servers still define.
		for _, sc := range client.AllowedOAuthScopes {
			if activeScopes[sc] {
				granted = append(granted, sc)
			}
		}
	} else {
		associated := make(map[string]bool, len(client.AllowedOAuthScopes))
		for _, sc := range client.AllowedOAuthScopes {
			associated[sc] = true
		}
		for _, sc := range in.Scopes {
			if associated[sc] {
				// A scope the client associates whose resource server no
				// longer defines it is inactive: the request proceeds
				// without it.
				if activeScopes[sc] {
					granted = append(granted, sc)
				}
				continue
			}
			if activeScopes[sc] {
				// Defined in the pool but not granted to this client: the
				// authentication fails.
				return nil, ErrNotAuthorized
			}
			// A scope no resource server in the pool defines is not a value
			// the Scopes member accepts.
			return nil, ErrInvalidParameter
		}
	}
	scopeClaim := strings.Join(granted, " ")

	customClaims := map[string]interface{}{"scope": scopeClaim}

	// The pre token generation trigger fires for M2M tokens only when the
	// pool configures it at event version V3_0. A machine identity has no
	// user attributes or groups, so only the claim overrides apply.
	if cfg := pool.LambdaConfig; cfg != nil && cfg.PreTokenGenerationConfig != nil &&
		cfg.PreTokenGenerationConfig.LambdaVersion == "V3_0" {
		result, terr := invokePreTokenGeneration(ctx, s, TokenGenerationClientCredentials,
			pool.ID, in.ClientID, in.ClientID, cfg, nil, nil, in.ClientMetadata)
		if terr != nil {
			return nil, ErrUserLambdaValidation
		}
		for k, v := range result.ClaimsToAddOrOverride {
			customClaims[k] = v
		}
		for _, k := range result.ClaimsToSuppress {
			delete(customClaims, k)
		}
	}

	privateKey, err := vsjwt.DecodePrivateKeyFromPEM(pool.JwtPrivateKey)
	if err != nil {
		return nil, ErrInternalError
	}
	issuer := fmt.Sprintf("https://%s/%s", cognitoIdpHost(reqCtx.GetRegion()), pool.ID)
	jwtManager, err := vsjwt.NewManager(privateKey, pool.JwtKeyID, issuer)
	if err != nil {
		return nil, ErrInternalError
	}

	// Token validity follows the client's configured access-token validity,
	// exactly as user token issuance does.
	validity := resolveClientTokenValidity(client)
	atValiditySec := int64(validity.accessMinutes) * 60

	accessToken, err := jwtManager.GenerateAccessTokenWithClaims(
		m2mTokenPrincipal{clientID: in.ClientID}, in.ClientID, atValiditySec, customClaims)
	if err != nil {
		return nil, ErrInternalError
	}

	// The stored record pins the token to the pool and client for
	// revocation, with an empty user ID marking the client principal.
	at := cognitostore.NewAccessToken(pool.ID, "", in.ClientID, scopeClaim, time.Now().Add(time.Duration(validity.accessMinutes)*time.Minute))
	at.Token = accessToken
	if err := store.CreateAccessToken(at); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"ClientAuthenticationResult": map[string]interface{}{
			"AccessToken": accessToken,
			"ExpiresIn":   atValiditySec,
			"TokenType":   "Bearer",
		},
	}, nil
}
