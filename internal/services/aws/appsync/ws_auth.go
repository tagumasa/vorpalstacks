package appsync

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
)

// headerSubprotocolPrefix marks the WebSocket subprotocol carrying the
// Base64URL-encoded connection credentials of the Events API handshake.
const headerSubprotocolPrefix = "header-"

// extractSubprotocolAuth decodes the credentials carried by the
// header-<base64url> WebSocket subprotocol: {"host":..., "x-api-key":...}
// for API key connections, or the SigV4 header set for IAM. It returns nil
// when no decodable auth subprotocol is present.
//
// The protocol list is read from the raw header values rather than via the
// websocket helper: clients may send the subprotocols as one comma-joined
// header line or as repeated header lines, and only the raw values see both
// forms.
func extractSubprotocolAuth(r *http.Request) map[string]string {
	for _, protoList := range r.Header.Values("Sec-WebSocket-Protocol") {
		for _, proto := range strings.Split(protoList, ",") {
			proto = strings.TrimSpace(proto)
			if !strings.HasPrefix(proto, headerSubprotocolPrefix) {
				continue
			}
			encoded := strings.TrimPrefix(proto, headerSubprotocolPrefix)
			decoded, err := base64.RawURLEncoding.DecodeString(encoded)
			if err != nil {
				continue
			}
			var auth map[string]string
			if json.Unmarshal(decoded, &auth) == nil {
				return auth
			}
		}
	}
	return nil
}

// verifyConnectionAuth enforces the connection-level authorisation of the
// Events API handshake: the subprotocol credentials must satisfy one of the
// API's configured ConnectionAuthModes before the WebSocket upgrade is
// accepted.
func (s *EventServer) verifyConnectionAuth(ctx context.Context, apiId string, auth map[string]string) bool {
	if apiId == "" || auth == nil || s.storeLookup == nil {
		return false
	}

	store, err := s.storeLookup(apiId)
	if err != nil || store == nil {
		return false
	}

	api, apiErr := store.GetApiById(apiId)
	if apiErr != nil || api == nil || api.EventConfig == nil {
		return false
	}

	generic := make(map[string]interface{}, len(auth))
	for k, v := range auth {
		generic[k] = v
	}

	for _, mode := range api.EventConfig.ConnectionAuthModes {
		switch mode.AuthType {
		case "API_KEY":
			if s.verifyAPIKey(store, apiId, generic) {
				return true
			}
		case "AWS_IAM":
			if s.verifySubprotocolIAM(apiId, auth) {
				return true
			}
		case "AMAZON_COGNITO_USER_POOLS":
			if s.verifyCognito(ctx, api.EventConfig.AuthProviders, generic) {
				return true
			}
		case "AWS_LAMBDA":
			if s.verifyLambdaAuthorizer(ctx, api.EventConfig.AuthProviders, apiId, generic) {
				return true
			}
		case "OPENID_CONNECT":
			// Verifying an OIDC token requires reaching the external
			// identity provider that issued it; the platform integrates no
			// external IdPs (recorded in docs/services.md), so credentials
			// are never accepted unverified.
		}
	}
	return false
}

// verifySubprotocolIAM verifies the SigV4 signature carried inside the
// header- subprotocol over the documented synthetic request: a POST to
// https://<host>/event with body {} and the subprotocol's header set.
func (s *EventServer) verifySubprotocolIAM(apiId string, auth map[string]string) bool {
	if s.sigVerifier == nil {
		return false
	}
	host := auth["host"]
	if host == "" {
		return false
	}

	req, err := http.NewRequest(http.MethodPost, "https://"+host+"/event", strings.NewReader("{}"))
	if err != nil {
		return false
	}
	for k, v := range auth {
		req.Header.Set(k, v)
	}

	if err := s.sigVerifier.VerifyRequest(req, "appsync", s.lookupRegion(apiId)); err != nil {
		logs.Warn("Events API connection IAM verification failed",
			logs.String("apiId", apiId), logs.Err(err))
		return false
	}
	return true
}

// apiKeyExpired reports whether an API key has passed its expiry epoch.
// An Expires value of 0 means the key never expires.
func apiKeyExpired(expires int64) bool {
	return expires > 0 && time.Now().Unix() > expires
}
