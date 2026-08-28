package appsync

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
)

// graphqlExecutionInput carries the parsed GraphQL-over-HTTP request into the
// execution Core. Body stays raw so the JSON parse keeps its original
// position after authorisation.
type graphqlExecutionInput struct {
	ApiId        string
	APIKeyHeader string
	Body         []byte
}

// graphqlWireError describes a pre-execution failure that the handler
// serialises as a GraphQL-style HTTP error response.
type graphqlWireError struct {
	HTTPStatus int
	ErrorType  string
	Message    string
}

// executeGraphQLCore authorises and executes a GraphQL request on the data
// plane. A non-nil graphqlWireError aborts execution before the engine runs.
func (s *AppSyncService) executeGraphQLCore(ctx context.Context, reqCtx *request.RequestContext, store *appsyncstore.AppSyncStore, in *graphqlExecutionInput) (*graphqlExecutionResult, *graphqlWireError) {
	if in.ApiId == "" {
		return nil, &graphqlWireError{http.StatusBadRequest, "BadRequestException", "apiId is required"}
	}

	api, err := store.GetGraphqlApiById(in.ApiId)
	if err != nil {
		return nil, &graphqlWireError{http.StatusNotFound, "NotFoundException", fmt.Sprintf("GraphQL API %s not found", in.ApiId)}
	}

	// Enforce per-authentication-type access control on the data-plane
	// endpoint. Management API calls are already gated by the dispatcher's
	// IAM authoriser; this check protects the GraphQL query path.
	if wireErr := authorizeGraphQLRequest(store, in.APIKeyHeader, api); wireErr != nil {
		return nil, wireErr
	}

	var gqlReq graphqlRequest
	if in.Body != nil {
		if err := json.Unmarshal(in.Body, &gqlReq); err != nil {
			return nil, &graphqlWireError{http.StatusBadRequest, "BadRequestException", "Request body must be valid JSON"}
		}
	}

	if gqlReq.Query == "" {
		return nil, &graphqlWireError{http.StatusBadRequest, "BadRequestException", "query is required"}
	}

	if gqlReq.Variables == nil {
		gqlReq.Variables = make(map[string]interface{})
	}

	engine := newGraphQLEngine(store, wrapBus(s.bus), &s.schemaCache)
	return engine.Execute(ctx, reqCtx, in.ApiId, &gqlReq), nil
}

// authorizeGraphQLRequest enforces the API's configured authentication types
// against the client's chosen auth method.
func authorizeGraphQLRequest(store *appsyncstore.AppSyncStore, apiKeyHeader string, api *appsyncstore.GraphqlApi) *graphqlWireError {
	// Build the set of configured auth types (primary + additional).
	configured := map[string]bool{api.AuthenticationType: true}
	for _, p := range api.AdditionalAuthenticationProviders {
		configured[p.AuthenticationType] = true
	}

	// Detect the client's chosen auth method from the API key header.
	if apiKeyHeader != "" {
		if configured["API_KEY"] {
			return authorizeGraphQLAPIKey(store, apiKeyHeader, api)
		}
		return &graphqlWireError{http.StatusUnauthorized, "UnauthorizedException", "API key authentication is not configured for this API."}
	}

	// No API key header — fall back to IAM or other configured types.
	if configured["AWS_IAM"] {
		// IAM authentication is enforced by the dispatcher authoriser; if
		// the request reached here it already passed IAM checks.
		return nil
	}

	// OPENID_CONNECT, AMAZON_COGNITO_USER_POOLS, AWS_LAMBDA:
	// JWT and Lambda-authoriser validation is not yet implemented.
	// Fail closed to prevent unauthenticated access.
	return &graphqlWireError{http.StatusUnauthorized, "UnauthorizedException", "You are not authorized to make this call."}
}

// authorizeGraphQLAPIKey validates an API-key-authenticated request against
// the stored, unexpired API keys of the API.
func authorizeGraphQLAPIKey(store *appsyncstore.AppSyncStore, keyValue string, api *appsyncstore.GraphqlApi) *graphqlWireError {
	if keyValue == "" {
		return &graphqlWireError{http.StatusUnauthorized, "UnauthorizedException", "You are not authorized to make this call."}
	}

	apiKey, err := store.GetApiKey(api.ApiId, keyValue)
	if err != nil {
		return &graphqlWireError{http.StatusUnauthorized, "UnauthorizedException", "You are not authorized to make this call."}
	}

	// Check expiry: AWS API keys have a 1-year default validity. Expired
	// keys must reject the request.
	if apiKey.Expires > 0 && time.Now().Unix() > apiKey.Expires {
		return &graphqlWireError{http.StatusUnauthorized, "UnauthorizedException", "The API key has expired."}
	}

	return nil
}
