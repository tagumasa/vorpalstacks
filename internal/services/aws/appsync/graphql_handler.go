package appsync

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/eventbus"
	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

// graphqlRequest represents the incoming GraphQL-over-HTTP request body.
type graphqlRequest struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	OperationName string                 `json:"operationName,omitempty"`
}

// graphqlResponse implements response.StreamableResponse to bypass the
// REST-JSON response wrapping pipeline. The dispatcher's writeResponse
// detects StreamableResponse and writes raw bytes with custom headers.
type graphqlResponse struct {
	body    []byte
	headers http.Header
	status  int
}

// GetStream returns the response body as an io.Reader for streaming.
func (r *graphqlResponse) GetStream() io.Reader {
	return bytes.NewReader(r.body)
}

// GetStreamHeaders returns the HTTP headers for the streamed response.
func (r *graphqlResponse) GetStreamHeaders() http.Header {
	return r.headers
}

// GetStreamStatusCode returns the HTTP status code for the streamed response.
func (r *graphqlResponse) GetStreamStatusCode() int {
	return r.status
}

// HandleGraphQLExecution processes POST /v1/apis/{apiId}/graphql.
// This is not a REST-JSON SDK operation — it is a raw GraphQL-over-HTTP endpoint.
// The parser returns sentinel operation name "GraphQLExecution" which routes here.
// Returns a graphqlResponse (StreamableResponse) to bypass REST-JSON wrapping.
func (s *AppSyncService) HandleGraphQLExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return s.graphqlErrorResponse(http.StatusInternalServerError, "InternalFailureException", "Failed to access store"), nil
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return s.graphqlErrorResponse(http.StatusBadRequest, "BadRequestException", "apiId is required"), nil
	}

	api, err := store.GetGraphqlApiById(apiId)
	if err != nil {
		return s.graphqlErrorResponse(http.StatusNotFound, "NotFoundException", fmt.Sprintf("GraphQL API %s not found", apiId)), nil
	}

	// Enforce per-authentication-type access control on the data-plane
	// endpoint. Management API calls are already gated by the dispatcher's
	// IAM authoriser (Phase 1); this check protects the GraphQL query path.
	if authResp := s.authorizeGraphQLRequest(store, req, api); authResp != nil {
		return authResp, nil
	}

	var gqlReq graphqlRequest
	if req.Body != nil {
		if err := json.Unmarshal(req.Body, &gqlReq); err != nil {
			return s.graphqlErrorResponse(http.StatusBadRequest, "BadRequestException", "Request body must be valid JSON"), nil
		}
	}

	if gqlReq.Query == "" {
		return s.graphqlErrorResponse(http.StatusBadRequest, "BadRequestException", "query is required"), nil
	}

	if gqlReq.Variables == nil {
		gqlReq.Variables = make(map[string]interface{})
	}

	engine := newGraphQLEngine(store, wrapBus(s.bus), &s.schemaCache)
	result := engine.Execute(ctx, reqCtx, apiId, &gqlReq)

	headers := http.Header{}
	headers.Set("Content-Type", "application/graphql-response+json")

	status := http.StatusOK
	if len(result.Errors) > 0 {
		for _, e := range result.Errors {
			if e.Message == "InternalFailureException" {
				status = http.StatusInternalServerError
				break
			}
		}
	}

	body, err := json.Marshal(result)
	if err != nil {
		return s.graphqlErrorResponse(http.StatusInternalServerError, "InternalFailureException", "Failed to marshal response"), nil
	}

	return &graphqlResponse{
		body:    body,
		headers: headers,
		status:  status,
	}, nil
}

// graphqlErrorResponse creates a graphqlResponse containing a single GraphQL error.
func (s *AppSyncService) graphqlErrorResponse(status int, errType, message string) *graphqlResponse {
	headers := http.Header{}
	headers.Set("Content-Type", "application/graphql-response+json")

	result := &graphqlExecutionResult{
		Errors: []graphqlError{{
			Message:   message,
			ErrorType: errType,
			Locations: nil,
			Path:      nil,
		}},
	}

	body, _ := json.Marshal(result)
	return &graphqlResponse{
		body:    body,
		headers: headers,
		status:  status,
	}
}

// authorizeGraphQLRequest validates that the caller is authenticated to
// execute queries against the GraphQL data-plane endpoint. The check
// depends on the API's configured authentication type:
//
//   - API_KEY: requires a valid, non-expired key in the x-api-key header.
//   - AWS_IAM: handled by the dispatcher's IAM authoriser; no additional
//     check is needed here.
//   - OPENID_CONNECT / AMAZON_COGNITO_USER_POOLS / AWS_LAMBDA: deferred
//     to future JWT and Lambda-authoriser integration work.
//
// Returns a non-nil *graphqlResponse when the request is denied; nil when
// the caller is authorised.
func (s *AppSyncService) authorizeGraphQLRequest(store *appsyncstore.AppSyncStore, req *request.ParsedRequest, api *appsyncstore.GraphqlApi) *graphqlResponse {
	switch api.AuthenticationType {
	case "API_KEY":
		return s.authorizeAPIKey(store, req, api)
	case "AWS_IAM":
		// IAM authentication is enforced by the dispatcher authoriser.
		return nil
	default:
		// OPENID_CONNECT, AMAZON_COGNITO_USER_POOLS, AWS_LAMBDA:
		// not yet implemented; allow requests to preserve existing behaviour.
		return nil
	}
}

// authorizeAPIKey validates the x-api-key header against stored API keys.
func (s *AppSyncService) authorizeAPIKey(store *appsyncstore.AppSyncStore, req *request.ParsedRequest, api *appsyncstore.GraphqlApi) *graphqlResponse {
	keyValue := ""
	if req.Headers != nil {
		keyValue = req.Headers.Get("x-api-key")
	}
	if keyValue == "" {
		return s.graphqlErrorResponse(http.StatusUnauthorized, "UnauthorizedException", "You are not authorized to make this call.")
	}

	apiKey, err := store.GetApiKey(api.ApiId, keyValue)
	if err != nil {
		return s.graphqlErrorResponse(http.StatusUnauthorized, "UnauthorizedException", "You are not authorized to make this call.")
	}

	// Check expiry: AWS API keys have a 1-year default validity. Expired
	// keys must reject the request.
	if apiKey.Expires > 0 && time.Now().Unix() > apiKey.Expires {
		return s.graphqlErrorResponse(http.StatusUnauthorized, "UnauthorizedException", "The API key has expired.")
	}

	return nil
}

// wrapBus wraps an eventbus.Bus into a BusPublisher for use by the GraphQL engine.
// Returns nil-safe adapter that silently drops events when bus is nil.
func wrapBus(bus eventbus.Bus) BusPublisher {
	if bus == nil {
		return nil
	}
	return &busPublisherAdapter{bus: bus}
}
