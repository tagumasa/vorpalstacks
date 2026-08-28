package appsync

import (
	"context"
	"strings"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"

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

	in := startSchemaCreationInput{
		ApiId:      request.GetStringParam(req.Parameters, "apiId"),
		Definition: request.GetStringParam(req.Parameters, "definition"),
	}

	status, err := s.startSchemaCreationCore(store, in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"status": status,
	}, nil
}

// GetSchemaCreationStatus retrieves the current schema creation status of a
// GraphQL API.
// GET /v1/apis/{apiId}/schemacreation
func (s *AppSyncService) GetSchemaCreationStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	result, err := s.getSchemaCreationStatusCore(store, request.GetStringParam(req.Parameters, "apiId"))
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"status": result.Status,
	}
	if result.Details != "" {
		response["details"] = result.Details
	}
	return response, nil
}

// GetIntrospectionSchema retrieves the introspection schema of a GraphQL API.
// GET /v1/apis/{apiId}/schema
func (s *AppSyncService) GetIntrospectionSchema(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	in := getIntrospectionSchemaInput{
		ApiId:             request.GetStringParam(req.Parameters, "apiId"),
		Format:            request.GetStringParam(req.Parameters, "format"),
		IncludeDirectives: request.GetBoolParam(req.Parameters, "includeDirectives"),
	}

	return s.getIntrospectionSchemaCore(store, in)
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
