// Package http provides HTTP server functionality for vorpalstacks.
package http

import (
	"net/http"

	"vorpalstacks/internal/services/aws/common/auth"
)

// SignatureMiddleware returns a middleware that verifies AWS signature version 4 on incoming requests.
func SignatureMiddleware(cfg SignatureConfig) func(http.Handler) http.Handler {
	if !cfg.Enabled {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}
	}

	provider := auth.NewStaticCredentialsProvider(
		cfg.AccessKeyID,
		cfg.SecretAccessKey,
		cfg.Region,
		"",
	)
	verifier := auth.NewSignatureV4Verifier(provider)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			service := extractServiceFromRequestFallback(r)
			if service == "" {
				http.Error(w, "Forbidden: Could not determine service", http.StatusForbidden)
				return
			}

			if err := verifier.VerifyRequest(r, service, cfg.Region); err != nil {
				http.Error(w, "Forbidden: Invalid signature", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
