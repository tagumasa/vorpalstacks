package appsync

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/eventbus"
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
func (s *AppSyncService) HandleGraphQLExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return s.graphqlErrorResponse(http.StatusInternalServerError, "InternalFailureException", "Failed to access store"), nil
	}

	in := &graphqlExecutionInput{
		ApiId: request.GetStringParam(req.Parameters, "apiId"),
		Body:  req.Body,
	}
	if req.Headers != nil {
		in.APIKeyHeader = req.Headers.Get("x-api-key")
	}

	result, wireErr := s.executeGraphQLCore(ctx, reqCtx, store, in)
	if wireErr != nil {
		return s.graphqlErrorResponse(wireErr.HTTPStatus, wireErr.ErrorType, wireErr.Message), nil
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/graphql-response+json")

	status := http.StatusOK
	if len(result.Errors) > 0 {
		for _, e := range result.Errors {
			if e.ErrorType == "INTERNAL_FAILURE" {
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

// graphqlErrorResponse builds a GraphQL-style HTTP error response body.
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

// wrapBus wraps an eventbus.Bus into a BusPublisher for use by the GraphQL engine.
// Returns nil-safe adapter that silently drops events when bus is nil.
func wrapBus(bus eventbus.Bus) BusPublisher {
	if bus == nil {
		return nil
	}
	return &busPublisherAdapter{bus: bus}
}
