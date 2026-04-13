package appsync

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// StartSchemaCreation initiates schema creation for a GraphQL API.
// The definition is base64-encoded SDL.
// POST /v1/apis/{apiId}/schemacreation
func (s *AppSyncService) StartSchemaCreation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	definitionB64 := request.GetStringParam(req.Parameters, "definition")
	if definitionB64 == "" {
		return nil, NewBadRequestException("definition is required")
	}

	decodedDef, err := base64.StdEncoding.DecodeString(definitionB64)
	if err != nil {
		decodedDef, err = base64.RawStdEncoding.DecodeString(definitionB64)
		if err != nil {
			return nil, NewBadRequestException("definition is not valid base64")
		}
	}

	_, err = store.GetGraphqlApiById(apiId)
	if err != nil {
		return nil, NewNotFoundException(fmt.Sprintf("GraphQL API with ID %s", apiId))
	}

	status := &appsyncstore.SchemaCreationStatus{
		ApiId:      apiId,
		Status:     "PROCESSING",
		Details:    "",
		Definition: string(decodedDef),
	}

	if err := store.SaveSchemaCreationStatus(apiId, status); err != nil {
		return nil, ErrInternalFailureException
	}

	// Capture the decoded definition in a closure variable to avoid
	// re-reading from the store, which introduces a race condition.
	defStr := string(decodedDef)
	go func() {
		completed := &appsyncstore.SchemaCreationStatus{
			ApiId:      apiId,
			Status:     "SUCCESS",
			Details:    "The schema was successfully created.",
			Definition: defStr,
		}
		_ = store.SaveSchemaCreationStatus(apiId, completed)
	}()

	return map[string]interface{}{
		"status": "PROCESSING",
	}, nil
}

// GetSchemaCreationStatus retrieves the status of a schema creation operation.
// GET /v1/apis/{apiId}/schemacreation
func (s *AppSyncService) GetSchemaCreationStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	status, err := store.GetSchemaCreationStatus(apiId)
	if err != nil {
		return map[string]interface{}{
			"status":  "NOT_APPLICABLE",
			"details": "No schema creation has been initiated for this API.",
		}, nil
	}

	response := map[string]interface{}{
		"status": status.Status,
	}
	if status.Details != "" {
		response["details"] = status.Details
	}
	return response, nil
}

// GetIntrospectionSchema returns the introspection schema for a GraphQL API.
// The response is raw bytes (SDL or JSON), not JSON-wrapped.
// GET /v1/apis/{apiId}/schema?format=SDL|JSON&includeDirectives=true|false
func (s *AppSyncService) GetIntrospectionSchema(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	_, err = store.GetGraphqlApiById(apiId)
	if err != nil {
		return mapStoreError(err)
	}

	format := request.GetStringParam(req.Parameters, "format")
	if format == "" {
		format = "SDL"
	}

	includeDirectives := request.GetBoolParam(req.Parameters, "includeDirectives")

	schemaSDL := collectSchemaSDL(store, apiId)
	if schemaSDL == "" {
		schemaSDL = buildIntrospectionSchema(includeDirectives)
	}

	if strings.EqualFold(format, "JSON") {
		return schemaJSONResponse(schemaSDL, includeDirectives), nil
	}

	return schemaSDL, nil
}

// collectSchemaSDL builds a complete SDL string from the schema creation status
// and any individual type definitions stored via CreateType.
func collectSchemaSDL(store *appsyncstore.AppSyncStore, apiId string) string {
	schemaStatus, err := store.GetSchemaCreationStatus(apiId)
	if err != nil {
		return ""
	}

	sdl := schemaStatus.Definition
	if sdl == "" {
		return ""
	}

	types, _, err := store.ListTypes(apiId, storecommon.ListOptions{MaxItems: 10000})
	if err != nil || len(types) == 0 {
		return sdl
	}

	for _, t := range types {
		if t.Definition != "" && !strings.Contains(sdl, t.Definition) {
			sdl += "\n\n" + t.Definition
		}
	}

	return sdl
}

// buildIntrospectionSchema generates a default introspection SDL.
func buildIntrospectionSchema(includeDirectives bool) string {
	schema := `schema {
  query: Query
  mutation: Mutation
  subscription: Subscription
}

type Query {
  _empty: String
}

type Mutation {
  _empty: String
}

type Subscription {
  _empty: String
}`

	if includeDirectives {
		schema += `

directive @aws_subscribe(mutations: [String!]!) on FIELD_DEFINITION

directive @aws_auth(cognito_groups: [String!]) on OBJECT | FIELD_DEFINITION

directive @aws_cognito_user_pools(cognito_groups: [String!]) on OBJECT | FIELD_DEFINITION

directive @aws_iam on OBJECT | FIELD_DEFINITION

directive @aws_api_key on OBJECT | FIELD_DEFINITION

directive @aws_oidc on OBJECT | FIELD_DEFINITION

directive @deprecated(reason: String) on FIELD_DEFINITION | ENUM_VALUE

directive @skip(if: Boolean!) on FIELD | FRAGMENT_SPREAD | INLINE_FRAGMENT

directive @include(if: Boolean!) on FIELD | FRAGMENT_SPREAD | INLINE_FRAGMENT

directive @connection(keyName: String!, fields: [String!]!) on FIELD_DEFINITION

directive @auth(rules: [AuthRule!]!) on OBJECT | FIELD_DEFINITION

directive @model(queries: ModelQueryInput, mutations: ModelMutationInput, subscriptions: ModelSubscriptionInput, timestamps: ModelTimestampInput) on OBJECT

scalar AWSDate

scalar AWSTime

scalar AWSDateTime

scalar AWSTimestamp

scalar AWSEmail

scalar AWSJSON

scalar AWSURL

scalar AWSPhone

scalar AWSIPAddress`
	}

	return schema
}

// schemaJSONResponse wraps the introspection schema in a JSON introspection response.
func schemaJSONResponse(sdl string, includeDirectives bool) map[string]interface{} {
	types := []interface{}{
		map[string]interface{}{
			"kind":          "OBJECT",
			"name":          "Query",
			"description":   "Root query type",
			"fields":        []interface{}{},
			"inputFields":   nil,
			"interfaces":    []interface{}{},
			"enumValues":    nil,
			"possibleTypes": nil,
		},
		map[string]interface{}{
			"kind":          "OBJECT",
			"name":          "Mutation",
			"description":   "Root mutation type",
			"fields":        []interface{}{},
			"inputFields":   nil,
			"interfaces":    []interface{}{},
			"enumValues":    nil,
			"possibleTypes": nil,
		},
		map[string]interface{}{
			"kind":          "OBJECT",
			"name":          "Subscription",
			"description":   "Root subscription type",
			"fields":        []interface{}{},
			"inputFields":   nil,
			"interfaces":    []interface{}{},
			"enumValues":    nil,
			"possibleTypes": nil,
		},
	}

	if includeDirectives {
		types = append(types,
			map[string]interface{}{
				"kind":          "SCALAR",
				"name":          "AWSDate",
				"description":   "The AWSDate scalar type represents a valid extended ISO 8601 Date string.",
				"fields":        nil,
				"inputFields":   nil,
				"interfaces":    nil,
				"enumValues":    nil,
				"possibleTypes": nil,
			},
			map[string]interface{}{
				"kind":          "SCALAR",
				"name":          "AWSTime",
				"description":   "The AWSTime scalar type represents a valid extended ISO 8601 Time string.",
				"fields":        nil,
				"inputFields":   nil,
				"interfaces":    nil,
				"enumValues":    nil,
				"possibleTypes": nil,
			},
			map[string]interface{}{
				"kind":          "SCALAR",
				"name":          "AWSDateTime",
				"description":   "The AWSDateTime scalar type represents a valid extended ISO 8601 DateTime string.",
				"fields":        nil,
				"inputFields":   nil,
				"interfaces":    nil,
				"enumValues":    nil,
				"possibleTypes": nil,
			},
		)
	}

	return map[string]interface{}{
		"__schema": map[string]interface{}{
			"types":            types,
			"queryType":        map[string]interface{}{"name": "Query"},
			"mutationType":     map[string]interface{}{"name": "Mutation"},
			"subscriptionType": map[string]interface{}{"name": "Subscription"},
			"directives":       []interface{}{},
		},
	}
}
