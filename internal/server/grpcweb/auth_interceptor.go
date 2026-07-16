package grpcweb

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"vorpalstacks/pkg/vsjwt"
)

type contextKey string

const claimsContextKey contextKey = "cognito_claims"

var (
	errMissingAuth = errors.New("missing or malformed authorization header")
	errInvalidTok  = errors.New("invalid or expired token")
)

// noAuthPathPrefixes lists RPC path prefixes that bypass authentication.
// The admin auth service must be accessible without a token because clients
// need to log in before they have one.
var noAuthPathPrefixes = []string{
	"/admin_auth.AdminAuthService/",
}

// ClaimsFromContext extracts CognitoClaims that were set by the auth
// interceptor or HTTP middleware.
func ClaimsFromContext(ctx context.Context) (*vsjwt.CognitoClaims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*vsjwt.CognitoClaims)
	return claims, ok
}

// NewAuthInterceptor creates a ConnectRPC UnaryInterceptorFunc that validates
// JWT access tokens from the Authorization: Bearer header.  RPCs whose
// procedure matches a prefix in noAuthPathPrefixes are allowed through without
// authentication.
func NewAuthInterceptor(jwtMgr *vsjwt.Manager) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			for _, prefix := range noAuthPathPrefixes {
				if strings.HasPrefix(req.Spec().Procedure, prefix) {
					return next(ctx, req)
				}
			}

			token, ok := extractBearerToken(req.Header().Get("Authorization"))
			if !ok {
				return nil, connect.NewError(connect.CodeUnauthenticated, errMissingAuth)
			}

			claims, err := jwtMgr.ValidateToken(token)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, errInvalidTok)
			}

			ctx = context.WithValue(ctx, claimsContextKey, claims)
			return next(ctx, req)
		}
	}
}

// newAuthHTTPMiddleware wraps an http.Handler with JWT authentication.
// Requests whose URL path starts with a noAuth prefix are allowed through
// without validation.  Validated claims are stored in the request context so
// downstream ConnectRPC handlers can retrieve them via ClaimsFromContext.
func newAuthHTTPMiddleware(jwtMgr *vsjwt.Manager, wrapped http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, prefix := range noAuthPathPrefixes {
			if strings.HasPrefix(r.URL.Path, prefix) {
				wrapped.ServeHTTP(w, r)
				return
			}
		}

		token, ok := extractBearerToken(r.Header.Get("Authorization"))
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := jwtMgr.ValidateToken(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), claimsContextKey, claims)
		wrapped.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractBearerToken parses a "Bearer <token>" string and returns the token.
func extractBearerToken(header string) (string, bool) {
	if header == "" {
		return "", false
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}
	token := strings.TrimPrefix(header, prefix)
	if token == "" {
		return "", false
	}
	return token, true
}
