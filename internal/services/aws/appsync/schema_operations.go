package appsync

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/utils/graphql"
)

// StartSchemaCreation initiates schema creation for a GraphQL API.
// The definition is base64-encoded SDL.
// POST /v1/apis/{apiId}/schemacreation
func (s *AppSyncService) StartSchemaCreation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
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

	defStr := string(decodedDef)
	s.schemaWg.Add(1)
	go func() {
		defer s.schemaWg.Done()
		defer func() { resilience.RecoverPanic("appsync schema creation") }()

		_, parseErr := gqlparser.LoadSchema(&ast.Source{
			Name:  "schema.graphql",
			Input: defStr,
		})

		if parseErr != nil {
			errMsg := parseErr.Error()
			completed := &appsyncstore.SchemaCreationStatus{
				ApiId:      apiId,
				Status:     "FAILED",
				Details:    errMsg,
				Definition: defStr,
			}
			saveSchemaStatusWithRetry(store, apiId, completed)
			return
		}

		completed := &appsyncstore.SchemaCreationStatus{
			ApiId:      apiId,
			Status:     "SUCCESS",
			Details:    "The schema was successfully created.",
			Definition: defStr,
		}
		saveSchemaStatusWithRetry(store, apiId, completed)
	}()

	return map[string]interface{}{
		"status": "PROCESSING",
	}, nil
}

// saveSchemaStatusWithRetry attempts to persist schema creation status with
// up to 3 retries on failure. This prevents the status from being stuck in
// PROCESSING if a transient Pebble error occurs during the goroutine save.
func saveSchemaStatusWithRetry(store *appsyncstore.AppSyncStore, apiId string, status *appsyncstore.SchemaCreationStatus) {
	for attempt := 0; attempt < 3; attempt++ {
		if err := store.SaveSchemaCreationStatus(apiId, status); err != nil {
			logs.Warn("failed to persist schema creation status",
				logs.String("apiId", apiId),
				logs.Int("attempt", attempt+1),
				logs.Err(err))
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}
		return
	}
	logs.Error("failed to persist schema creation status after 3 retries",
		logs.String("apiId", apiId),
		logs.String("status", status.Status))
}

// GetSchemaCreationStatus retrieves the status of a schema creation operation.
// GET /v1/apis/{apiId}/schemacreation
func (s *AppSyncService) GetSchemaCreationStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
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
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	_, err = store.GetGraphqlApiById(apiId)
	if err != nil {
		return mapStoreError(err)
	}

	// Check if schema creation has failed — AWS returns GraphQLSchemaException
	// when the introspection schema is in an invalid state.
	if status, err := store.GetSchemaCreationStatus(apiId); err == nil && status.Status == "FAILED" {
		return nil, NewGraphQLSchemaException(status.Details)
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
		return schemaJSONFromSDL(schemaSDL, includeDirectives), nil
	}

	if !includeDirectives {
		schemaSDL = stripDirectivesFromSDL(schemaSDL)
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

	types, err := store.GetAllTypesForApi(apiId)
	if err != nil || len(types) == 0 {
		return sdl
	}

	for _, t := range types {
		if t.Definition != "" && !typeDefInSDL(sdl, t.Definition) {
			sdl += "\n\n" + t.Definition
		}
	}

	return sdl
}

var typeNamePrefixes = []string{"type ", "input ", "enum ", "interface ", "union ", "scalar "}

func typeDefInSDL(sdl, def string) bool {
	typeName := extractTypeName(def)
	if typeName == "" {
		return strings.Contains(sdl, def)
	}
	for _, prefix := range typeNamePrefixes {
		if strings.Contains(def, prefix) {
			return strings.Contains(sdl, prefix+typeName+" ") ||
				strings.Contains(sdl, prefix+typeName+"{") ||
				strings.Contains(sdl, prefix+typeName+"\n")
		}
	}
	return strings.Contains(sdl, def)
}

func extractTypeName(def string) string {
	return graphql.ExtractTypeName(def)
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

// schemaJSONFromSDL parses the SDL using gqlparser and builds a complete
// GraphQL introspection JSON response using the engine's buildSchemaObject.
// When includeDirectives is false, the directives array is omitted.
func schemaJSONFromSDL(sdl string, includeDirectives bool) map[string]interface{} {
	if sdl == "" {
		return map[string]interface{}{
			"__schema": map[string]interface{}{
				"types":      []interface{}{},
				"directives": []interface{}{},
			},
		}
	}

	schema, err := gqlparser.LoadSchema(&ast.Source{
		Name:  "schema.graphql",
		Input: sdl,
	})
	if err != nil {
		return map[string]interface{}{
			"__schema": map[string]interface{}{
				"types":      []interface{}{},
				"directives": []interface{}{},
			},
		}
	}

	engine := &graphQLEngine{}
	schemaObj := engine.buildSchemaObject(schema)
	if !includeDirectives {
		schemaObj["directives"] = []interface{}{}
	}
	return map[string]interface{}{
		"__schema": schemaObj,
	}
}

// stripDirectivesFromSDL removes directive definitions (e.g.
// "directive @aws_api_key on FIELD_DEFINITION") from SDL output.
// Directive applications on type fields (e.g. "@deprecated") are preserved.
// Handles multi-line directive definitions by tracking the start
// ("directive @") until the "on" keyword that terminates the locations list.
func stripDirectivesFromSDL(sdl string) string {
	lines := strings.Split(sdl, "\n")
	filtered := make([]string, 0, len(lines))
	inDirectiveDef := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "directive @") {
			inDirectiveDef = true
		}
		if inDirectiveDef {
			// Directive definitions end with the "on" keyword
			// followed by one or more locations.
			if strings.Contains(trimmed, " on ") || strings.HasSuffix(trimmed, " on") {
				inDirectiveDef = false
			}
			continue
		}
		filtered = append(filtered, line)
	}
	return strings.Join(filtered, "\n")
}
