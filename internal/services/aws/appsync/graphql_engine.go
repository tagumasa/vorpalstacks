package appsync

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	appsyncstore "vorpalstacks/internal/store/aws/appsync"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/pkg/vtl"
)

// jsonUnmarshal is an alias for json.Unmarshal to avoid shadowing in template methods.
var jsonUnmarshal = json.Unmarshal

// graphqlExecutionResult is the standard GraphQL response envelope.
type graphqlExecutionResult struct {
	Data   interface{}    `json:"data,omitempty"`
	Errors []graphqlError `json:"errors,omitempty"`
}

// graphqlError represents a single error in a GraphQL response.
type graphqlError struct {
	Message   string      `json:"message"`
	Locations []gqlErrLoc `json:"locations,omitempty"`
	Path      interface{} `json:"path,omitempty"`
}

// gqlErrLoc represents a source location within a GraphQL error.
type gqlErrLoc struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// schemaCacheEntry holds a parsed schema and resolver map for a given API.
type schemaCacheEntry struct {
	schema      *ast.Schema
	sdl         string
	hash        [32]byte
	resolverMap map[string]map[string]*appsyncstore.Resolver
}

func schemaHash(sdl string) [32]byte {
	return sha256.Sum256([]byte(sdl))
}

// graphQLEngine orchestrates GraphQL query execution against a stored schema,
// dispatching field resolution to AppSync resolvers (unit or pipeline).
type graphQLEngine struct {
	store       *appsyncstore.AppSyncStore
	bus         BusPublisher
	schemaCache *sync.Map
}

// BusPublisher abstracts the event bus publish capability for WebSocket
// subscription fan-out. Uses an interface to allow nil-safe calls.
type BusPublisher interface {
	Publish(ctx context.Context, event interface{}) error
}

// busPublisherAdapter wraps an eventbus.Bus to satisfy the BusPublisher interface.
type busPublisherAdapter struct {
	bus eventbus.Bus
}

// Publish dispatches an event through the underlying event bus if it satisfies the Event interface.
func (a *busPublisherAdapter) Publish(ctx context.Context, event interface{}) error {
	if e, ok := event.(eventbus.Event); ok {
		return a.bus.Publish(ctx, e)
	}
	return nil
}

// newGraphQLEngine creates a new GraphQL execution engine scoped to the given store.
func newGraphQLEngine(store *appsyncstore.AppSyncStore, bus BusPublisher, schemaCache *sync.Map) *graphQLEngine {
	return &graphQLEngine{
		store:       store,
		bus:         bus,
		schemaCache: schemaCache,
	}
}

// Execute processes a GraphQL request and returns the execution result.
func (e *graphQLEngine) Execute(ctx context.Context, reqCtx *request.RequestContext, apiId string, gqlReq *graphqlRequest) *graphqlExecutionResult {
	entry, err := e.loadSchema(ctx, reqCtx, apiId)
	if err != nil {
		return &graphqlExecutionResult{
			Errors: []graphqlError{{Message: err.Error()}},
		}
	}
	schema := entry.schema

	doc, errs := gqlparser.LoadQueryWithRules(schema, gqlReq.Query, nil)
	if errs != nil {
		return &graphqlExecutionResult{
			Errors: convertGqlErrors(errs),
		}
	}

	if len(doc.Operations) == 0 {
		return &graphqlExecutionResult{
			Errors: []graphqlError{{Message: "No operations found in document"}},
		}
	}

	var op *ast.OperationDefinition
	if gqlReq.OperationName != "" {
		for _, o := range doc.Operations {
			if o.Name == gqlReq.OperationName {
				op = o
				break
			}
		}
		if op == nil {
			return &graphqlExecutionResult{
				Errors: []graphqlError{{Message: fmt.Sprintf("Operation %q not found", gqlReq.OperationName)}},
			}
		}
	} else if len(doc.Operations) == 1 {
		op = doc.Operations[0]
	} else {
		return &graphqlExecutionResult{
			Errors: []graphqlError{{Message: "Must provide operation name when multiple operations exist"}},
		}
	}

	data, execErrs := e.executeOperation(ctx, reqCtx, apiId, schema, op, gqlReq.Variables, entry.resolverMap, doc.Fragments)

	result := &graphqlExecutionResult{Data: data}
	if len(execErrs) > 0 {
		result.Errors = execErrs
	}
	return result
}

// loadSchema retrieves the SDL from the store, parses it with gqlparser,
// and caches the result. The cache is invalidated on each call by comparing
// the stored SDL with the cached version.
func (e *graphQLEngine) loadSchema(ctx context.Context, reqCtx *request.RequestContext, apiId string) (*schemaCacheEntry, error) {
	sdl := collectSchemaSDL(e.store, apiId)

	if sdl == "" {
		return nil, fmt.Errorf("no schema found for API %s", apiId)
	}

	h := schemaHash(sdl)

	if cached, ok := e.schemaCache.Load(apiId); ok {
		entry := cached.(*schemaCacheEntry)
		if entry.hash == h && entry.resolverMap != nil {
			return entry, nil
		}
	}

	schema, err := gqlparser.LoadSchema(&ast.Source{
		Name:  "schema.graphql",
		Input: sdl,
	})
	if err != nil {
		return nil, fmt.Errorf("schema parsing failed: %w", err)
	}

	injectIntrospectionTypes(schema)

	rMap, err := e.buildResolverMap(apiId)
	if err != nil {
		return nil, fmt.Errorf("failed to build resolver map: %w", err)
	}

	entry := &schemaCacheEntry{schema: schema, sdl: sdl, hash: h, resolverMap: rMap}
	e.schemaCache.Store(apiId, entry)
	return entry, nil
}

// buildResolverMap scans all resolvers for the given API and builds a
// nested map: resolverMap[typeName][fieldName] -> *appsyncstore.Resolver.
func (e *graphQLEngine) buildResolverMap(apiId string) (map[string]map[string]*appsyncstore.Resolver, error) {
	resolvers, _, err := e.store.ListResolvers(apiId, "", common.ListOptions{MaxItems: 10000})
	if err != nil {
		return nil, err
	}

	rMap := make(map[string]map[string]*appsyncstore.Resolver)
	for _, r := range resolvers {
		if rMap[r.TypeName] == nil {
			rMap[r.TypeName] = make(map[string]*appsyncstore.Resolver)
		}
		rMap[r.TypeName][r.FieldName] = r
	}
	return rMap, nil
}

// executeOperation dispatches to the appropriate root type handler based on
// the operation kind (query, mutation, subscription).
func (e *graphQLEngine) executeOperation(
	ctx context.Context,
	reqCtx *request.RequestContext,
	apiId string,
	schema *ast.Schema,
	op *ast.OperationDefinition,
	variables map[string]interface{},
	resolverMap map[string]map[string]*appsyncstore.Resolver,
	fragments ast.FragmentDefinitionList,
) (interface{}, []graphqlError) {
	switch op.Operation {
	case ast.Query:
		if op.SelectionSet == nil {
			return nil, nil
		}
		return e.resolveSelectionSet(ctx, reqCtx, apiId, schema, "Query", nil, op.SelectionSet, variables, resolverMap, fragments)
	case ast.Mutation:
		if op.SelectionSet == nil {
			return nil, nil
		}
		return e.resolveSelectionSet(ctx, reqCtx, apiId, schema, "Mutation", nil, op.SelectionSet, variables, resolverMap, fragments)
	case ast.Subscription:
		if op.SelectionSet == nil {
			return nil, nil
		}
		return e.resolveSelectionSet(ctx, reqCtx, apiId, schema, "Subscription", nil, op.SelectionSet, variables, resolverMap, fragments)
	default:
		return nil, []graphqlError{{Message: fmt.Sprintf("Unsupported operation type: %s", op.Operation)}}
	}
}

// resolveSelectionSet walks a SelectionSet and resolves each field.
// For root fields, it looks up the resolver from resolverMap.
// For child fields on object results, it extracts from the parent source.
func (e *graphQLEngine) resolveSelectionSet(
	ctx context.Context,
	reqCtx *request.RequestContext,
	apiId string,
	schema *ast.Schema,
	parentTypeName string,
	parentSource interface{},
	selectionSet ast.SelectionSet,
	variables map[string]interface{},
	resolverMap map[string]map[string]*appsyncstore.Resolver,
	fragments ast.FragmentDefinitionList,
) (map[string]interface{}, []graphqlError) {
	result := make(map[string]interface{})
	var errs []graphqlError

	for _, sel := range selectionSet {
		switch s := sel.(type) {
		case *ast.Field:
			fieldName := s.Name
			alias := fieldName
			if s.Alias != "" {
				alias = s.Alias
			}

			if fieldName == "__typename" {
				result[alias] = parentTypeName
				continue
			}

			parentType := schema.Types[parentTypeName]
			if parentType == nil || parentType.Fields == nil {
				logs.Warn("GraphQL type not found in schema",
					logs.String("type", parentTypeName),
					logs.String("field", fieldName))
				continue
			}

			fieldDef := parentType.Fields.ForName(fieldName)
			if fieldDef == nil {
				logs.Warn("GraphQL field not found in schema",
					logs.String("type", parentTypeName),
					logs.String("field", fieldName))
				continue
			}

			resolved, fieldErrs := e.resolveField(
				ctx, reqCtx, apiId, schema,
				parentTypeName, fieldName, parentSource,
				s, variables, resolverMap,
			)

			if len(fieldErrs) > 0 {
				errs = append(errs, fieldErrs...)
			}

			if resolved != nil && s.SelectionSet != nil {
				resolved = e.completeValue(
					ctx, reqCtx, apiId, schema,
					fieldDef.Type, resolved,
					s.SelectionSet, variables, resolverMap, fragments,
				)
			}

			result[alias] = resolved

		case *ast.FragmentSpread:
			fragDef := fragments.ForName(s.Name)
			if fragDef != nil {
				fragResult, fragErrs := e.resolveSelectionSet(
					ctx, reqCtx, apiId, schema,
					parentTypeName, parentSource,
					fragDef.SelectionSet, variables, resolverMap, fragments,
				)
				for k, v := range fragResult {
					result[k] = v
				}
				errs = append(errs, fragErrs...)
			}

		case *ast.InlineFragment:
			if s.TypeCondition != "" {
				if parentSource != nil {
					if srcMap, ok := parentSource.(map[string]interface{}); ok {
						if typename, ok := srcMap["__typename"].(string); ok && typename != "" && typename != s.TypeCondition {
							continue
						}
					}
				}
				fragResult, fragErrs := e.resolveSelectionSet(
					ctx, reqCtx, apiId, schema,
					s.TypeCondition, parentSource,
					s.SelectionSet, variables, resolverMap, fragments,
				)
				for k, v := range fragResult {
					result[k] = v
				}
				errs = append(errs, fragErrs...)
			} else {
				fragResult, fragErrs := e.resolveSelectionSet(
					ctx, reqCtx, apiId, schema,
					parentTypeName, parentSource,
					s.SelectionSet, variables, resolverMap, fragments,
				)
				for k, v := range fragResult {
					result[k] = v
				}
				errs = append(errs, fragErrs...)
			}
		}
	}

	return result, errs
}

// resolveField resolves a single field. If a resolver exists for the
// (parentType, fieldName) pair, it executes the resolver pipeline.
// Otherwise, it extracts the value from the parent source.
func (e *graphQLEngine) resolveField(
	ctx context.Context,
	reqCtx *request.RequestContext,
	apiId string,
	schema *ast.Schema,
	parentTypeName string,
	fieldName string,
	parentSource interface{},
	field *ast.Field,
	variables map[string]interface{},
	resolverMap map[string]map[string]*appsyncstore.Resolver,
) (interface{}, []graphqlError) {

	args, err := e.resolveArguments(field, variables)
	if err != nil {
		return nil, []graphqlError{{Message: fmt.Sprintf("Failed to resolve arguments for %s.%s: %v", parentTypeName, fieldName, err)}}
	}

	if resolved, handled := e.resolveIntrospectionField(schema, parentTypeName, fieldName, parentSource, args); handled {
		return resolved, nil
	}

	resolver, hasResolver := resolverMap[parentTypeName][fieldName]

	if !hasResolver {
		if parentSource != nil {
			return e.extractFromSource(parentSource, fieldName)
		}
		return nil, nil
	}

	fieldType := e.resolveFieldType(schema, field.Definition.Type)
	isList := e.isListType(field.Definition.Type)

	if resolver.Kind == "PIPELINE" && resolver.PipelineConfig != nil && len(resolver.PipelineConfig.Functions) > 0 {
		return e.executePipelineResolver(ctx, reqCtx, apiId, schema, resolver, parentTypeName, fieldName, parentSource, args, fieldType, isList, resolverMap, variables)
	}

	return e.executeUnitResolver(ctx, reqCtx, apiId, schema, resolver, parentTypeName, fieldName, parentSource, args, fieldType, isList)
}

// executeUnitResolver runs a single resolver: VTL request template → DataSource dispatch → VTL response template.
func (e *graphQLEngine) executeUnitResolver(
	ctx context.Context,
	reqCtx *request.RequestContext,
	apiId string,
	schema *ast.Schema,
	resolver *appsyncstore.Resolver,
	parentTypeName string,
	fieldName string,
	parentSource interface{},
	args map[string]interface{},
	fieldType string,
	isList bool,
) (interface{}, []graphqlError) {

	engine := vtl.NewEngine()
	engine.AppSyncCtx = &vtl.AppSyncContext{
		Args:   args,
		Source: parentSource,
		Stash:  make(map[string]interface{}),
		Info: &vtl.AppSyncFieldInfo{
			FieldName:           fieldName,
			ParentTypeName:      parentTypeName,
			SelectionSetGraphQL: "",
			SelectionSetList:    []string{},
		},
	}

	var dsPayload interface{}

	if resolver.RequestMappingTemplate != "" {
		result, err := engine.Transform(resolver.RequestMappingTemplate)
		if err != nil {
			return nil, []graphqlError{{Message: fmt.Sprintf("Request template error: %v", err)}}
		}

		if engine.AppSyncCtx != nil && len(engine.AppSyncCtx.Errors) > 0 {
			var errs []graphqlError
			for _, ae := range engine.AppSyncCtx.Errors {
				errs = append(errs, graphqlError{Message: ae.Message})
			}
			return nil, errs
		}

		if result != "" {
			var parsed interface{}
			if jsonErr := jsonUnmarshal([]byte(result), &parsed); jsonErr != nil {
				dsPayload = result
			} else {
				dsPayload = parsed
			}
		}
	}

	dsResult, dsErr := e.dispatchDataSource(ctx, reqCtx, apiId, resolver.DataSourceName, dsPayload)

	if dsErr != nil {
		return nil, []graphqlError{{Message: dsErr.Error()}}
	}

	if resolver.ResponseMappingTemplate != "" {
		engine.AppSyncCtx.Result = dsResult
		respStr, err := engine.Transform(resolver.ResponseMappingTemplate)
		if err != nil {
			return nil, []graphqlError{{Message: fmt.Sprintf("Response template error: %v", err)}}
		}

		if engine.AppSyncCtx != nil && len(engine.AppSyncCtx.Errors) > 0 {
			var errs []graphqlError
			for _, ae := range engine.AppSyncCtx.Errors {
				errs = append(errs, graphqlError{Message: ae.Message})
			}
			return nil, errs
		}

		if respStr != "" {
			var parsed interface{}
			if jsonErr := jsonUnmarshal([]byte(respStr), &parsed); jsonErr != nil {
				return respStr, nil
			}
			return parsed, nil
		}
	}

	return dsResult, nil
}

// executePipelineResolver runs a pipeline resolver: optional before template →
// function chain (each with request/response templates) → optional after template.
func (e *graphQLEngine) executePipelineResolver(
	ctx context.Context,
	reqCtx *request.RequestContext,
	apiId string,
	schema *ast.Schema,
	resolver *appsyncstore.Resolver,
	parentTypeName string,
	fieldName string,
	parentSource interface{},
	args map[string]interface{},
	fieldType string,
	isList bool,
	resolverMap map[string]map[string]*appsyncstore.Resolver,
	variables map[string]interface{},
) (interface{}, []graphqlError) {

	stash := make(map[string]interface{})

	engine := vtl.NewEngine()
	engine.AppSyncCtx = &vtl.AppSyncContext{
		Args:   args,
		Source: parentSource,
		Stash:  stash,
		Info: &vtl.AppSyncFieldInfo{
			FieldName:           fieldName,
			ParentTypeName:      parentTypeName,
			SelectionSetGraphQL: "",
			SelectionSetList:    []string{},
		},
	}

	if resolver.RequestMappingTemplate != "" {
		_, _ = engine.Transform(resolver.RequestMappingTemplate)
		if engine.AppSyncCtx != nil && len(engine.AppSyncCtx.Errors) > 0 {
			var errs []graphqlError
			for _, ae := range engine.AppSyncCtx.Errors {
				errs = append(errs, graphqlError{Message: ae.Message})
			}
			return nil, errs
		}
	}

	var result interface{} = map[string]interface{}{}

	for _, functionId := range resolver.PipelineConfig.Functions {
		fn, err := e.store.GetFunction(apiId, functionId)
		if err != nil {
			return nil, []graphqlError{{Message: fmt.Sprintf("Function %s not found: %v", functionId, err)}}
		}

		fnEngine := vtl.NewEngine()
		fnEngine.AppSyncCtx = &vtl.AppSyncContext{
			Args:   args,
			Source: parentSource,
			Stash:  stash,
			Result: result,
			Info: &vtl.AppSyncFieldInfo{
				FieldName:           fieldName,
				ParentTypeName:      parentTypeName,
				SelectionSetGraphQL: "",
				SelectionSetList:    []string{},
			},
		}

		var fnPayload interface{}
		if fn.RequestMappingTemplate != "" {
			fnResult, err := fnEngine.Transform(fn.RequestMappingTemplate)
			if err != nil {
				return nil, []graphqlError{{Message: fmt.Sprintf("Function %s request template error: %v", functionId, err)}}
			}
			if fnEngine.AppSyncCtx != nil && len(fnEngine.AppSyncCtx.Errors) > 0 {
				var errs []graphqlError
				for _, ae := range fnEngine.AppSyncCtx.Errors {
					errs = append(errs, graphqlError{Message: ae.Message})
				}
				return nil, errs
			}
			if fnResult != "" {
				var parsed interface{}
				if jsonErr := jsonUnmarshal([]byte(fnResult), &parsed); jsonErr != nil {
					fnPayload = fnResult
				} else {
					fnPayload = parsed
				}
			}
		}

		dsResult, dsErr := e.dispatchDataSource(ctx, reqCtx, apiId, fn.DataSourceName, fnPayload)
		if dsErr != nil {
			return nil, []graphqlError{{Message: dsErr.Error()}}
		}

		if fn.ResponseMappingTemplate != "" {
			fnEngine.AppSyncCtx.Result = dsResult
			respStr, err := fnEngine.Transform(fn.ResponseMappingTemplate)
			if err != nil {
				return nil, []graphqlError{{Message: fmt.Sprintf("Function %s response template error: %v", functionId, err)}}
			}
			if fnEngine.AppSyncCtx != nil && len(fnEngine.AppSyncCtx.Errors) > 0 {
				var errs []graphqlError
				for _, ae := range fnEngine.AppSyncCtx.Errors {
					errs = append(errs, graphqlError{Message: ae.Message})
				}
				return nil, errs
			}
			if respStr != "" {
				var parsed interface{}
				if jsonErr := jsonUnmarshal([]byte(respStr), &parsed); jsonErr != nil {
					result = respStr
				} else {
					result = parsed
				}
			} else {
				result = dsResult
			}
		} else {
			result = dsResult
		}
	}

	if resolver.ResponseMappingTemplate != "" {
		engine.AppSyncCtx.Result = result
		engine.AppSyncCtx.Stash = stash
		respStr, err := engine.Transform(resolver.ResponseMappingTemplate)
		if err != nil {
			return nil, []graphqlError{{Message: fmt.Sprintf("Pipeline after template error: %v", err)}}
		}
		if engine.AppSyncCtx != nil && len(engine.AppSyncCtx.Errors) > 0 {
			var errs []graphqlError
			for _, ae := range engine.AppSyncCtx.Errors {
				errs = append(errs, graphqlError{Message: ae.Message})
			}
			return nil, errs
		}
		if respStr != "" {
			var parsed interface{}
			if jsonErr := jsonUnmarshal([]byte(respStr), &parsed); jsonErr != nil {
				return respStr, nil
			}
			return parsed, nil
		}
	}

	return result, nil
}

// resolveArguments extracts field arguments using the gqlparser Value(vars) method.
func (e *graphQLEngine) resolveArguments(field *ast.Field, variables map[string]interface{}) (map[string]interface{}, error) {
	if len(field.Arguments) == 0 {
		return nil, nil
	}

	args := make(map[string]interface{})
	for _, arg := range field.Arguments {
		if arg.Value == nil {
			continue
		}
		val, err := arg.Value.Value(variables)
		if err != nil {
			return nil, err
		}
		args[arg.Name] = val
	}
	return args, nil
}

// resolveFieldType extracts the named type string from a gqlparser Type,
// unwrapping NonNull and List wrappers.
func (e *graphQLEngine) resolveFieldType(schema *ast.Schema, t *ast.Type) string {
	if t == nil {
		return "String"
	}
	named := t.NamedType
	if named != "" && schema.Types[named] != nil {
		return named
	}
	return "String"
}

// isListType checks whether a gqlparser Type represents a list type.
// In gqlparser v2, a list type has Elem != nil and NamedType == "".
func (e *graphQLEngine) isListType(t *ast.Type) bool {
	for t != nil {
		if t.Elem != nil && t.NamedType == "" {
			return true
		}
		t = t.Elem
	}
	return false
}

// extractFromSource extracts a field value from a parent source map.
func (e *graphQLEngine) extractFromSource(source interface{}, fieldName string) (interface{}, []graphqlError) {
	switch s := source.(type) {
	case map[string]interface{}:
		if v, ok := s[fieldName]; ok {
			return v, nil
		}
		return nil, nil
	default:
		return nil, nil
	}
}

// completeValue recursively resolves child fields on a resolved value when
// the field type is a composite (object, interface, or union) type.
// This is the "field completion" step in the GraphQL execution algorithm:
// after a resolver returns a value, composite-typed fields must have their
// own selection sets resolved against the returned data.
func (e *graphQLEngine) completeValue(
	ctx context.Context,
	reqCtx *request.RequestContext,
	apiId string,
	schema *ast.Schema,
	fieldType *ast.Type,
	value interface{},
	selectionSet ast.SelectionSet,
	variables map[string]interface{},
	resolverMap map[string]map[string]*appsyncstore.Resolver,
	fragments ast.FragmentDefinitionList,
) interface{} {
	if fieldType == nil || value == nil {
		return value
	}

	namedType := fieldType.NamedType
	if e.isListType(fieldType) {
		list, ok := value.([]interface{})
		if !ok {
			return value
		}
		elemType := fieldType.Elem
		for elemType.Elem != nil {
			elemType = elemType.Elem
		}
		if elemType == nil || elemType.NamedType == "" || !e.isCompositeType(schema, elemType.NamedType) {
			return value
		}
		result := make([]interface{}, len(list))
		for i, item := range list {
			if item != nil {
				result[i] = e.completeValue(ctx, reqCtx, apiId, schema, elemType, item, selectionSet, variables, resolverMap, fragments)
			}
		}
		return result
	}

	if !e.isCompositeType(schema, namedType) {
		return value
	}

	childResult, _ := e.resolveSelectionSet(
		ctx, reqCtx, apiId, schema,
		namedType, value,
		selectionSet, variables, resolverMap, fragments,
	)
	return childResult
}

// isCompositeType reports whether a named type is an object, interface, or union
// type that requires further field resolution. Scalar and enum types do not.
func (e *graphQLEngine) isCompositeType(schema *ast.Schema, typeName string) bool {
	if typeName == "" {
		return false
	}
	def := schema.Types[typeName]
	if def == nil {
		return false
	}
	switch def.Kind {
	case ast.Object, ast.Interface, ast.Union:
		return true
	default:
		return false
	}
}

// convertGqlErrors converts gqlparser error list to our graphqlError format.
func convertGqlErrors(errs gqlerror.List) []graphqlError {
	var result []graphqlError
	for _, e := range errs {
		locs := make([]gqlErrLoc, 0, len(e.Locations))
		for _, loc := range e.Locations {
			locs = append(locs, gqlErrLoc{
				Line:   loc.Line,
				Column: loc.Column,
			})
		}
		result = append(result, graphqlError{
			Message:   e.Message,
			Locations: locs,
		})
	}
	return result
}
