package http

import (
	"net/http"

	"vorpalstacks/internal/server/http/classifier"
	"vorpalstacks/internal/common/auth"
)

// SignatureMiddleware returns an HTTP middleware that verifies AWS Signature
// Version 4 on incoming requests. It uses the classifier to determine the
// target service name for signature verification. When signature verification
// is disabled, it returns a pass-through middleware.
func SignatureMiddleware(cfg SignatureConfig, c *classifier.Classifier) func(http.Handler) http.Handler {
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

			cr, err := c.Classify(r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			if cr == nil || cr.ServiceName == "" {
				http.Error(w, "Forbidden: Could not determine service", http.StatusForbidden)
				return
			}

			if err := verifier.VerifyRequest(r, cr.ServiceName, cfg.Region); err != nil {
				http.Error(w, "Forbidden: Invalid signature", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
