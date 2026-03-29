package http

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"

	"vorpalstacks/internal/server/http/router"
)

func (s *Server) registerRoutes(r chi.Router) {
	if s.chainGateway != nil && s.config.UseChainGateway {
		s.registerChainRoutes(r)
		return
	}

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Amz-Target") != "" {
				next.ServeHTTP(w, r)
				return
			}
			if r.URL.Query().Get("Action") != "" {
				next.ServeHTTP(w, r)
				return
			}
			if r.Method == "POST" && strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
				next.ServeHTTP(w, r)
				return
			}
			if strings.Contains(r.Header.Get("Content-Type"), "application/cbor") {
				next.ServeHTTP(w, r)
				return
			}
			if router.IsLambdaRestPath(r.URL.Path) || isApiGatewayPath(r.URL.Path) || isSchedulerPath(r.URL.Path) || isSESv2Path(r.URL.Path) || isRoute53Path(r.URL.Path) || isCloudFrontPath(r.URL.Path) || isCloudWatchCBORPath(r.URL.Path) {
				next.ServeHTTP(w, r)
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

	awsJSONHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Amz-Target") != "" {
			s.handleRequest(w, r)
			return
		}
		if router.IsLambdaRestPath(r.URL.Path) {
			s.handleRequest(w, r)
			return
		}
		apiGatewayRuntime := s.APIGatewayRuntimeHandler()
		if isApiGatewayRuntimePath(r.URL.Path) && apiGatewayRuntime != nil {
			apiGatewayRuntime.ServeHTTP(w, r)
			return
		}
		if isApiGatewayPath(r.URL.Path) {
			s.handleRequest(w, r)
			return
		}
		if r.URL.Query().Get("Action") != "" {
			s.handleRequest(w, r)
			return
		}
		if r.Method == "POST" && strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
			s.handleRequest(w, r)
			return
		}
		s.handleRequest(w, r)
	}

	r.HandleFunc("/", awsJSONHandler)
	r.HandleFunc("/*", awsJSONHandler)
}

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

		routePath := "/"
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

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	serviceName := s.extractServiceFromRequest(r)
	s.dispatcher.Dispatch(w, r, serviceName, nil)
}

func isApiGatewayPath(path string) bool {
	if isApiGatewayRuntimePath(path) {
		return false
	}
	return strings.HasPrefix(path, "/restapis") ||
		strings.HasPrefix(path, "/apikeys") ||
		strings.HasPrefix(path, "/usageplans") ||
		strings.HasPrefix(path, "/domainnames") ||
		strings.HasPrefix(path, "/vpclinks") ||
		strings.HasPrefix(path, "/apis") ||
		strings.HasPrefix(path, "/authorizers") ||
		(strings.HasPrefix(path, "/tags/") && strings.Contains(path, "apigateway"))
}

func isApiGatewayRuntimePath(path string) bool {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 4 && parts[0] == "restapis" && parts[3] == "_user_request_" {
		return true
	}
	return false
}

func isSchedulerPath(path string) bool {
	return strings.HasPrefix(path, "/schedule-groups") ||
		strings.HasPrefix(path, "/schedules") ||
		strings.HasPrefix(path, "/tags/arn:aws:scheduler")
}

func isSESv2Path(path string) bool {
	return strings.HasPrefix(path, "/v2/email/")
}

func isRoute53Path(path string) bool {
	return strings.HasPrefix(path, "/2013-04-01/")
}

func isCloudFrontPath(path string) bool {
	return strings.HasPrefix(path, "/2020-05-31/")
}

func isCloudWatchCBORPath(path string) bool {
	return strings.HasPrefix(path, "/service/GraniteServiceVersion20100801/")
}

func isS3StyleRoute(routePath string) bool {
	return strings.Contains(routePath, "{Bucket}") ||
		strings.Contains(routePath, "{Key+")
}

func extractPath(uri string) string {
	for i, c := range uri {
		if c == '?' {
			return uri[:i]
		}
	}
	return uri
}
