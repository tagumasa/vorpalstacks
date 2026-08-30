package apigateway

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// testInvokePayload parses the request-simulation wire members shared by the
// test-invoke operations.
func testInvokePayload(req *request.ParsedRequest) *TestInvokePayload {
	p := &TestInvokePayload{
		Body:                request.GetStringParam(req.Parameters, "body"),
		PathWithQueryString: request.GetStringParam(req.Parameters, "pathWithQueryString"),
		Headers:             make(map[string]string),
		MultiValueHeaders:   make(map[string][]string),
		StageVariables:      make(map[string]string),
		AdditionalContext:   make(map[string]string),
	}
	if h, ok := req.Parameters["headers"].(map[string]interface{}); ok {
		for k, v := range h {
			if vs, ok := v.(string); ok {
				p.Headers[k] = vs
			}
		}
	}
	if mvh, ok := req.Parameters["multiValueHeaders"].(map[string]interface{}); ok {
		for k, v := range mvh {
			if arr, ok := v.([]interface{}); ok {
				for _, item := range arr {
					if vs, ok := item.(string); ok {
						p.MultiValueHeaders[k] = append(p.MultiValueHeaders[k], vs)
					}
				}
			}
		}
	}
	if sv, ok := req.Parameters["stageVariables"].(map[string]interface{}); ok {
		for k, v := range sv {
			if vs, ok := v.(string); ok {
				p.StageVariables[k] = vs
			}
		}
	}
	if ac, ok := req.Parameters["additionalContext"].(map[string]interface{}); ok {
		for k, v := range ac {
			if vs, ok := v.(string); ok {
				p.AdditionalContext[k] = vs
			}
		}
	}
	return p
}

// TestInvokeMethod simulates an API method invocation for testing purposes.
func (s *APIGatewayService) TestInvokeMethod(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	p := testInvokePayload(req)

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.testInvokeMethodCore(ctx, stores, apiId, resourceId, httpMethod, p)
}

// TestInvokeAuthorizer simulates an authoriser invocation for testing purposes.
func (s *APIGatewayService) TestInvokeAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	authorizerId := request.GetStringParam(req.Parameters, "authorizerId")
	p := testInvokePayload(req)
	// body and pathWithQueryString are accepted per Smithy model but not used
	// in this stub authorizer test implementation.

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.testInvokeAuthorizerCore(stores, apiId, authorizerId, p)
}
