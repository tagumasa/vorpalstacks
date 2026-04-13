package appsync

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/common/request"
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

	_, err = store.GetGraphqlApiById(apiId)
	if err != nil {
		return s.graphqlErrorResponse(http.StatusNotFound, "NotFoundException", fmt.Sprintf("GraphQL API %s not found", apiId)), nil
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

	engine := newGraphQLEngine(store, s.lambdaInvoker, wrapBus(s.bus), &s.schemaCache)
	result := engine.Execute(ctx, reqCtx, apiId, &gqlReq)

	headers := http.Header{}
	headers.Set("Content-Type", "application/graphql-response+json")

	status := http.StatusOK
	if result.Errors != nil && len(result.Errors) > 0 {
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

// wrapBus wraps an eventbus.Bus into a BusPublisher for use by the GraphQL engine.
// Returns nil-safe adapter that silently drops events when bus is nil.
func wrapBus(bus eventbus.Bus) BusPublisher {
	if bus == nil {
		return nil
	}
	return &busPublisherAdapter{bus: bus}
}
