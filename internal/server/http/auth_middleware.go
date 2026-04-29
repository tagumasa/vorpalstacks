package http

import (
	"io"
	"net/http"
	"strings"

	"vorpalstacks/internal/common/auth"
	"vorpalstacks/internal/server/http/classifier"
)

// SignatureMiddleware returns an HTTP middleware that verifies AWS Signature
// Version 4 on incoming requests. It supports both static credentials and
// temporary STS session credentials via the optional SessionResolver.
func SignatureMiddleware(cfg SignatureConfig, c *classifier.Classifier, sessionResolver auth.SessionResolver) func(http.Handler) http.Handler {
	if !cfg.Enabled {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}
	}

	staticProvider := auth.NewStaticCredentialsProvider(
		cfg.AccessKeyID,
		cfg.SecretAccessKey,
		cfg.Region,
		"",
	)
	staticVerifier := auth.NewSignatureV4Verifier(staticProvider)

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

			if isUnauthenticatedOperation(r) {
				next.ServeHTTP(w, r)
				return
			}

			if err := staticVerifier.VerifyRequest(r, cr.ServiceName, cfg.Region); err != nil {
				if trySessionCredentials(sessionResolver, cfg, cr.ServiceName, r, w, next) {
					return
				}
				http.Error(w, "Forbidden: Invalid signature", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func trySessionCredentials(
	sessionResolver auth.SessionResolver,
	cfg SignatureConfig,
	serviceName string,
	r *http.Request,
	w http.ResponseWriter,
	next http.Handler,
) bool {
	if sessionResolver == nil {
		return false
	}

	accessKeyId := extractAccessKeyIdFromAuth(r.Header.Get("Authorization"))
	if accessKeyId == "" || !strings.HasPrefix(accessKeyId, "ASIA") {
		return false
	}

	session, err := sessionResolver.ResolveSession(accessKeyId)
	if err != nil || session == nil {
		return false
	}

	securityToken := r.Header.Get("X-Amz-Security-Token")
	if securityToken != session.SessionToken {
		return false
	}

	sessionProvider := auth.NewStaticCredentialsProvider(
		session.AccessKeyID,
		session.SecretAccessKey,
		cfg.Region,
		"",
	)
	sessionVerifier := auth.NewSignatureV4Verifier(sessionProvider)
	if err := sessionVerifier.VerifyRequest(r, serviceName, cfg.Region); err != nil {
		return false
	}

	r.Header.Set("X-Amz-Access-Key", session.AccessKeyID)
	next.ServeHTTP(w, r)
	return true
}

func extractAccessKeyIdFromAuth(authHeader string) string {
	if !strings.HasPrefix(authHeader, "AWS4-HMAC-SHA256 ") {
		return ""
	}
	for _, part := range strings.Split(strings.TrimPrefix(authHeader, "AWS4-HMAC-SHA256 "), ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "Credential=") {
			cred := strings.TrimPrefix(part, "Credential=")
			fields := strings.Split(cred, "/")
			if len(fields) > 0 {
				return fields[0]
			}
		}
	}
	return ""
}

var unauthenticatedOps = map[string]bool{
	"AssumeRoleWithSAML":        true,
	"AssumeRoleWithWebIdentity": true,
}

func isUnauthenticatedOperation(r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
		return false
	}

	var bodyBytes []byte
	if r.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(r.Body)
		if err != nil {
			return false
		}
	}

	r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	for _, pair := range strings.Split(string(bodyBytes), "&") {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 && kv[0] == "Action" {
			decoded := strings.ReplaceAll(kv[1], "+", " ")
			return unauthenticatedOps[decoded]
		}
	}
	return false
}
