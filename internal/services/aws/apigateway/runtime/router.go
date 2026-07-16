// Package runtime provides API Gateway runtime functionality for vorpalstacks.
package runtime

import (
	"fmt"
	"strings"

	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
)

// RouteMatch represents a matched route with method information.
type RouteMatch struct {
	ResourceId string
	Path       string
	HttpMethod string
	Params     map[string]string
	Method     *apigatewaystore.Method
}

// RuntimeRouter handles API Gateway routing.
type RuntimeRouter struct{}

// NewRuntimeRouter creates a new runtime router.
func NewRuntimeRouter() *RuntimeRouter {
	return &RuntimeRouter{}
}

// ResourceInfo contains information about a resource and its methods.
type ResourceInfo struct {
	Id              string
	Path            string
	ParentId        string
	ResourceMethods map[string]*apigatewaystore.Method
}

// MatchWithMethods matches a request path and method against available resources.
func (r *RuntimeRouter) MatchWithMethods(resources []ResourceInfo, path string, httpMethod string) (*RouteMatch, error) {
	path = normalizePath(path)

	for _, resource := range resources {
		if resource.Path == path {
			method := resource.ResourceMethods[httpMethod]
			if method == nil {
				method = resource.ResourceMethods[strings.ToUpper(httpMethod)]
			}
			if method != nil {
				return &RouteMatch{
					ResourceId: resource.Id,
					Path:       path,
					HttpMethod: httpMethod,
					Params:     make(map[string]string),
					Method:     method,
				}, nil
			}
		}
	}

	for _, resource := range resources {
		if params := matchPathWithParams(resource.Path, path); params != nil {
			method := resource.ResourceMethods[httpMethod]
			if method == nil {
				method = resource.ResourceMethods[strings.ToUpper(httpMethod)]
			}
			if method != nil {
				return &RouteMatch{
					ResourceId: resource.Id,
					Path:       resource.Path,
					HttpMethod: httpMethod,
					Params:     params,
					Method:     method,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("no matching route for %s %s", httpMethod, path)
}

func normalizePath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

func matchPathWithParams(pattern string, path string) map[string]string {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	params := make(map[string]string)

	patternIdx := 0
	pathIdx := 0

	for patternIdx < len(patternParts) && pathIdx < len(pathParts) {
		patternPart := patternParts[patternIdx]
		pathPart := pathParts[pathIdx]

		if strings.HasPrefix(patternPart, "{") && strings.HasSuffix(patternPart, "}") {
			paramName := patternPart[1 : len(patternPart)-1]

			if strings.HasSuffix(paramName, "+") {
				paramName = paramName[:len(paramName)-1]
				remainingParts := pathParts[pathIdx:]
				params[paramName] = strings.Join(remainingParts, "/")
				return params
			}

			params[paramName] = pathPart
			patternIdx++
			pathIdx++
		} else if patternPart == pathPart {
			patternIdx++
			pathIdx++
		} else {
			return nil
		}
	}

	if patternIdx == len(patternParts) && pathIdx == len(pathParts) {
		return params
	}

	return nil
}
