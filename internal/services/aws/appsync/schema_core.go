package appsync

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

// startSchemaCreationInput carries the parsed StartSchemaCreation request.
// Definition holds the base64-encoded SDL exactly as supplied on the wire.
type startSchemaCreationInput struct {
	ApiId      string
	Definition string
}

// getIntrospectionSchemaInput carries the parsed GetIntrospectionSchema
// request.
type getIntrospectionSchemaInput struct {
	ApiId             string
	Format            string
	IncludeDirectives bool
}

// schemaCreationStatusResult is the Core result of GetSchemaCreationStatus:
// a missing status resolves to the documented NOT_APPLICABLE fallback.
type schemaCreationStatusResult struct {
	Status  string
	Details string
}

// startSchemaCreationCore validates and persists the PROCESSING schema
// status and validates the SDL asynchronously.
func (s *AppSyncService) startSchemaCreationCore(store *appsyncstore.AppSyncStore, in startSchemaCreationInput) (string, error) {
	if in.ApiId == "" {
		return "", NewBadRequestException("apiId is required")
	}
	if in.Definition == "" {
		return "", NewBadRequestException("definition is required")
	}

	decodedDef, err := base64.StdEncoding.DecodeString(in.Definition)
	if err != nil {
		decodedDef, err = base64.RawStdEncoding.DecodeString(in.Definition)
		if err != nil {
			return "", NewBadRequestException("definition is not valid base64")
		}
	}

	_, err = store.GetGraphqlApiById(in.ApiId)
	if err != nil {
		return "", NewNotFoundException(fmt.Sprintf("GraphQL API with ID %s", in.ApiId))
	}

	status := &appsyncstore.SchemaCreationStatus{
		ApiId:      in.ApiId,
		Status:     "PROCESSING",
		Details:    "",
		Definition: string(decodedDef),
	}

	if err := store.SaveSchemaCreationStatus(in.ApiId, status); err != nil {
		return "", ErrInternalFailureException
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
				ApiId:      in.ApiId,
				Status:     "FAILED",
				Details:    errMsg,
				Definition: defStr,
			}
			saveSchemaStatusWithRetry(store, in.ApiId, completed)
			return
		}

		completed := &appsyncstore.SchemaCreationStatus{
			ApiId:      in.ApiId,
			Status:     "SUCCESS",
			Details:    "The schema was successfully created.",
			Definition: defStr,
		}
		saveSchemaStatusWithRetry(store, in.ApiId, completed)
	}()

	return "PROCESSING", nil
}

// getSchemaCreationStatusCore reads the current schema creation status of a
// GraphQL API.
func (s *AppSyncService) getSchemaCreationStatusCore(store *appsyncstore.AppSyncStore, apiId string) (schemaCreationStatusResult, error) {
	if apiId == "" {
		return schemaCreationStatusResult{}, NewBadRequestException("apiId is required")
	}

	status, err := store.GetSchemaCreationStatus(apiId)
	if err != nil {
		return schemaCreationStatusResult{
			Status:  "NOT_APPLICABLE",
			Details: "No schema creation has been initiated for this API.",
		}, nil
	}

	return schemaCreationStatusResult{
		Status:  status.Status,
		Details: status.Details,
	}, nil
}

// getIntrospectionSchemaCore resolves the introspection schema of a GraphQL
// API in the requested output serialisation. The format member is required
// on GetIntrospectionSchema and validated against the SDL and JSON enum
// values.
func (s *AppSyncService) getIntrospectionSchemaCore(store *appsyncstore.AppSyncStore, in getIntrospectionSchemaInput) (interface{}, error) {
	if in.ApiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}
	if in.Format == "" {
		return nil, NewBadRequestException("format is required")
	}
	if !validateTypeFormat(in.Format) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid format: %s. Valid values: SDL, JSON", in.Format))
	}

	_, err := store.GetGraphqlApiById(in.ApiId)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	// Check if schema creation has failed — AWS returns GraphQLSchemaException
	// when the introspection schema is in an invalid state.
	if status, err := store.GetSchemaCreationStatus(in.ApiId); err == nil && status.Status == "FAILED" {
		return nil, NewGraphQLSchemaException(status.Details)
	}

	format := in.Format

	schemaSDL := collectSchemaSDL(store, in.ApiId)
	if schemaSDL == "" {
		schemaSDL = buildIntrospectionSchema(in.IncludeDirectives)
	}

	if strings.EqualFold(format, "JSON") {
		return schemaJSONFromSDL(schemaSDL, in.IncludeDirectives), nil
	}

	if !in.IncludeDirectives {
		schemaSDL = stripDirectivesFromSDL(schemaSDL)
	}

	return schemaSDL, nil
}

// saveSchemaStatusWithRetry persists a schema creation status with up to
// three attempts, escalating to an error log on final failure.
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

// collectSchemaSDL assembles the stored schema SDL plus every stored type
// definition that is not already part of it.
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
