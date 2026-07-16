package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

// registerRoutes sets up the HTTP routing table, installing the classify
// middleware and registering service-specific routes. When the chain gateway
// is enabled, it delegates entirely to registerChainRoutes.
func (s *Server) registerRoutes(r chi.Router) {
	if s.chainGateway != nil && s.config.UseChainGateway {
		s.registerChainRoutes(r)
		return
	}

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/.well-known/health" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"status":"ok"}`))
				return
			}

			cr, err := s.classifier.Classify(r)
			if err != nil || cr == nil {
				if s3Handler := s.S3Handler(); s3Handler != nil {
					s3Handler.ServeHTTP(w, r)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			if cr.ServiceName == "s3" {
				if s3Handler := s.S3Handler(); s3Handler != nil {
					s3Handler.ServeHTTP(w, r)
					return
				}
			}

			if cr.ServiceName != "" {
				s.dispatcher.DispatchClassified(w, r, cr)
				return
			}

			if s3Handler := s.S3Handler(); s3Handler != nil {
				s3Handler.ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	services, err := s.dispatcher.ListServices()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to list services: %v\n", err)
	} else {
		for _, svc := range services {
			s.registerServiceRoutes(r, svc)
		}
	}

	if apiGatewayRuntime := s.APIGatewayRuntimeHandler(); apiGatewayRuntime != nil {
		r.HandleFunc("/restapis/{restApiId}/{stageName}/_user_request_/*", func(w http.ResponseWriter, req *http.Request) {
			apiGatewayRuntime.ServeHTTP(w, req)
		})
	}

	if jwksHandler := s.JWKSHandler(); jwksHandler != nil {
		r.Handle("/.well-known/jwks.json", jwksHandler)
	}

	r.HandleFunc("/", s.handleClassifiedRequest)
	r.HandleFunc("/*", s.handleClassifiedRequest)
}

// registerChainRoutes sets up the chain gateway as the sole request handler.
func (s *Server) registerChainRoutes(r chi.Router) {
	if apiGatewayRuntime := s.APIGatewayRuntimeHandler(); apiGatewayRuntime != nil {
		r.HandleFunc("/restapis/{restApiId}/{stageName}/_user_request_/*", func(w http.ResponseWriter, req *http.Request) {
			apiGatewayRuntime.ServeHTTP(w, req)
		})
	}

	if jwksHandler := s.JWKSHandler(); jwksHandler != nil {
		r.Handle("/.well-known/jwks.json", jwksHandler)
	}

	r.HandleFunc("/", s.chainGateway.ServeHTTP)
	r.HandleFunc("/*", s.chainGateway.ServeHTTP)
}

// registerServiceRoutes registers chi routes for each operation of the given
// service. S3-style routes (containing {Bucket} or {Key+}) are skipped
// because S3 has its own dedicated handler.
func (s *Server) registerServiceRoutes(r chi.Router, serviceName string) {
	if serviceName == "s3" {
		return
	}

	operations, err := s.dispatcher.ListOperations(serviceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to list operations for service %s: %v\n", serviceName, err)
		return
	}

	for _, op := range operations {
		operation := op
		if operation.HTTPMethod == nil || *operation.HTTPMethod == "" {
			continue
		}

		var routePath string
		if operation.HTTPURIPath != nil && *operation.HTTPURIPath != "" {
			routePath = *operation.HTTPURIPath
		} else if operation.HTTPURI != nil && *operation.HTTPURI != "" {
			routePath = extractPath(*operation.HTTPURI)
		} else {
			continue
		}

		if isS3StyleRoute(routePath) {
			continue
		}

		r.Method(*operation.HTTPMethod, routePath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.dispatcher.Dispatch(w, r, serviceName, operation)
		}))
	}
}

// handleClassifiedRequest is the chi fallback handler that classifies
// unrecognised routes and dispatches them via DispatchClassified.
func (s *Server) handleClassifiedRequest(w http.ResponseWriter, r *http.Request) {
	cr, err := s.classifier.Classify(r)
	if err != nil || cr == nil {
		if s3Handler := s.S3Handler(); s3Handler != nil {
			s3Handler.ServeHTTP(w, r)
			return
		}
		s.dispatcher.Dispatch(w, r, "", nil)
		return
	}
	if cr.ServiceName == "s3" {
		if s3Handler := s.S3Handler(); s3Handler != nil {
			s3Handler.ServeHTTP(w, r)
			return
		}
	}
	s.dispatcher.DispatchClassified(w, r, cr)
}

// isS3StyleRoute reports whether the route path contains S3-style path
// parameters ({Bucket} or {Key+}).
func isS3StyleRoute(routePath string) bool {
	for _, s := range []string{"{Bucket}", "{Key+"} {
		if len(routePath) >= len(s) {
			for i := 0; i <= len(routePath)-len(s); i++ {
				if routePath[i:i+len(s)] == s {
					return true
				}
			}
		}
	}
	return false
}

// extractPath strips query parameters from a URI, returning only the path portion.
func extractPath(uri string) string {
	for i, c := range uri {
		if c == '?' {
			return uri[:i]
		}
	}
	return uri
}
