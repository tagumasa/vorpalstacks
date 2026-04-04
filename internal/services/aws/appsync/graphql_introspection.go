package appsync

import (
	"fmt"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

// injectIntrospectionTypes adds the standard GraphQL introspection types to the
// schema so that __schema, __type, and __typename queries can be executed.
// gqlparser injects __schema and __type field definitions on the Query type but
// does NOT provide the corresponding type definitions (__Schema, __Type, etc.),
// nor does it provide any runtime resolution logic.
func injectIntrospectionTypes(schema *ast.Schema) {
	if schema.Types["__Schema"] != nil {
		return
	}

	schema.Types["__Schema"] = &ast.Definition{
		Kind: ast.Object,
		Name: "__Schema",
		Fields: ast.FieldList{
			{Name: "types", Type: nonNullListType("__Type")},
			{Name: "queryType", Type: nonNullType("__Type")},
			{Name: "mutationType", Type: namedType("__Type")},
			{Name: "subscriptionType", Type: namedType("__Type")},
			{Name: "directives", Type: nonNullListType("__Directive")},
		},
	}

	schema.Types["__Type"] = &ast.Definition{
		Kind: ast.Object,
		Name: "__Type",
		Fields: ast.FieldList{
			{Name: "kind", Type: nonNullType("__TypeKind")},
			{Name: "name", Type: namedType("String")},
			{Name: "description", Type: namedType("String")},
			{Name: "fields", Type: listType("__Field"), Arguments: ast.ArgumentDefinitionList{{Name: "includeDeprecated", Type: namedType("Boolean")}}},
			{Name: "interfaces", Type: listType("__Type")},
			{Name: "possibleTypes", Type: listType("__Type")},
			{Name: "enumValues", Type: listType("__EnumValue"), Arguments: ast.ArgumentDefinitionList{{Name: "includeDeprecated", Type: namedType("Boolean")}}},
			{Name: "inputFields", Type: listType("__InputValue")},
			{Name: "ofType", Type: namedType("__Type")},
		},
	}

	schema.Types["__Field"] = &ast.Definition{
		Kind: ast.Object,
		Name: "__Field",
		Fields: ast.FieldList{
			{Name: "name", Type: nonNullType("String")},
			{Name: "description", Type: namedType("String")},
			{Name: "args", Type: nonNullListType("__InputValue")},
			{Name: "type", Type: nonNullType("__Type")},
			{Name: "isDeprecated", Type: nonNullType("Boolean")},
			{Name: "deprecationReason", Type: namedType("String")},
		},
	}

	schema.Types["__InputValue"] = &ast.Definition{
		Kind: ast.Object,
		Name: "__InputValue",
		Fields: ast.FieldList{
			{Name: "name", Type: nonNullType("String")},
			{Name: "description", Type: namedType("String")},
			{Name: "type", Type: nonNullType("__Type")},
			{Name: "defaultValue", Type: namedType("String")},
		},
	}

	schema.Types["__EnumValue"] = &ast.Definition{
		Kind: ast.Object,
		Name: "__EnumValue",
		Fields: ast.FieldList{
			{Name: "name", Type: nonNullType("String")},
			{Name: "description", Type: namedType("String")},
			{Name: "isDeprecated", Type: nonNullType("Boolean")},
			{Name: "deprecationReason", Type: namedType("String")},
		},
	}

	schema.Types["__Directive"] = &ast.Definition{
		Kind: ast.Object,
		Name: "__Directive",
		Fields: ast.FieldList{
			{Name: "name", Type: nonNullType("String")},
			{Name: "description", Type: namedType("String")},
			{Name: "locations", Type: nonNullListType("__DirectiveLocation")},
			{Name: "args", Type: nonNullListType("__InputValue")},
		},
	}

	schema.Types["__TypeKind"] = &ast.Definition{
		Kind:        ast.Enum,
		Name:        "__TypeKind",
		Description: "An enum describing what kind of type a given __Type is.",
		EnumValues: ast.EnumValueList{
			{Name: "SCALAR"},
			{Name: "OBJECT"},
			{Name: "INTERFACE"},
			{Name: "UNION"},
			{Name: "ENUM"},
			{Name: "INPUT_OBJECT"},
			{Name: "LIST"},
			{Name: "NON_NULL"},
		},
	}

	schema.Types["__DirectiveLocation"] = &ast.Definition{
		Kind:        ast.Enum,
		Name:        "__DirectiveLocation",
		Description: "An enum describing valid locations where a directive can be placed.",
		EnumValues: ast.EnumValueList{
			{Name: "QUERY"},
			{Name: "MUTATION"},
			{Name: "SUBSCRIPTION"},
			{Name: "FIELD"},
			{Name: "FRAGMENT_DEFINITION"},
			{Name: "FRAGMENT_SPREAD"},
			{Name: "INLINE_FRAGMENT"},
			{Name: "VARIABLE_DEFINITION"},
			{Name: "SCHEMA"},
			{Name: "SCALAR"},
			{Name: "OBJECT"},
			{Name: "FIELD_DEFINITION"},
			{Name: "ARGUMENT_DEFINITION"},
			{Name: "INTERFACE"},
			{Name: "UNION"},
			{Name: "ENUM"},
			{Name: "ENUM_VALUE"},
			{Name: "INPUT_OBJECT"},
			{Name: "INPUT_FIELD_DEFINITION"},
		},
	}
}

// namedType creates a nullable named type reference.
func namedType(name string) *ast.Type {
	return &ast.Type{NamedType: name}
}

// nonNullType creates a NonNull wrapper around a named type.
func nonNullType(name string) *ast.Type {
	return &ast.Type{NamedType: name, NonNull: true}
}

// listType creates a nullable list of a named type.
func listType(elem string) *ast.Type {
	return &ast.Type{NamedType: elem, Elem: &ast.Type{NamedType: elem}}
}

// nonNullListType creates a NonNull list of NonNull element type.
func nonNullListType(elem string) *ast.Type {
	return &ast.Type{NonNull: true, Elem: &ast.Type{NamedType: elem, NonNull: true, Elem: &ast.Type{NamedType: elem}}}
}

// resolveIntrospectionField handles __typename, __schema, and __type fields.
// Returns (resolved, handled) — if handled is true, the caller should use the
// resolved value instead of the normal resolver pipeline.
func (e *graphQLEngine) resolveIntrospectionField(
	schema *ast.Schema,
	parentTypeName string,
	fieldName string,
	parentSource interface{},
	args map[string]interface{},
) (interface{}, bool) {
	switch fieldName {
	case "__typename":
		if parentSource != nil {
			return parentTypeName, true
		}
		return parentTypeName, true

	case "__schema":
		return e.buildSchemaObject(schema), true

	case "__type":
		if args == nil {
			return nil, true
		}
		name, ok := args["name"].(string)
		if !ok || name == "" {
			return nil, true
		}
		typeDef := schema.Types[name]
		if typeDef == nil {
			return nil, true
		}
		return e.buildTypeObject(schema, typeDef), true
	}

	return nil, false
}

// buildSchemaObject constructs the __Schema object for introspection.
func (e *graphQLEngine) buildSchemaObject(schema *ast.Schema) map[string]interface{} {
	types := make([]interface{}, 0)
	for _, name := range sortedTypeNames(schema) {
		def := schema.Types[name]
		if strings.HasPrefix(name, "__") {
			continue
		}
		types = append(types, e.buildTypeObject(schema, def))
	}

	result := map[string]interface{}{
		"types":      types,
		"queryType":  e.buildTypeObject(schema, schema.Query),
		"directives": e.buildDirectives(schema),
	}
	if schema.Mutation != nil {
		result["mutationType"] = e.buildTypeObject(schema, schema.Mutation)
	}
	if schema.Subscription != nil {
		result["subscriptionType"] = e.buildTypeObject(schema, schema.Subscription)
	}
	return result
}

// buildTypeObject constructs a __Type representation from an ast.Definition.
func (e *graphQLEngine) buildTypeObject(schema *ast.Schema, def *ast.Definition) map[string]interface{} {
	if def == nil {
		return nil
	}

	result := map[string]interface{}{
		"name": def.Name,
		"kind": typeKindString(def.Kind),
	}
	if def.Description != "" {
		result["description"] = def.Description
	}

	switch def.Kind {
	case ast.Object, ast.Interface:
		result["fields"] = e.buildFields(schema, def.Fields, false)
		interfaces := make([]interface{}, 0)
		for _, ifaceName := range def.Interfaces {
			if ifaceDef := schema.Types[ifaceName]; ifaceDef != nil {
				interfaces = append(interfaces, e.buildTypeObject(schema, ifaceDef))
			}
		}
		result["interfaces"] = interfaces

	case ast.Union:
		possibleTypes := make([]interface{}, 0)
		for _, typeName := range def.Types {
			if typeDef := schema.Types[typeName]; typeDef != nil {
				possibleTypes = append(possibleTypes, e.buildTypeObject(schema, typeDef))
			}
		}
		result["possibleTypes"] = possibleTypes

	case ast.Enum:
		result["enumValues"] = e.buildEnumValues(def.EnumValues, false)

	case ast.InputObject:
		result["inputFields"] = e.buildInputFields(schema, def.Fields)
	}

	return result
}

// buildTypeRefObject constructs a __Type representation from an ast.Type reference.
// This handles wrapper types (NonNull, List) by recursing.
func (e *graphQLEngine) buildTypeRefObject(schema *ast.Schema, t *ast.Type) map[string]interface{} {
	if t == nil {
		return nil
	}

	if t.NonNull {
		inner := e.buildTypeRefObject(schema, t.Elem)
		if inner == nil {
			return nil
		}
		inner["kind"] = "NON_NULL"
		return inner
	}

	if t.Elem != nil {
		inner := e.buildTypeRefObject(schema, t.Elem)
		if inner == nil {
			return nil
		}
		inner["kind"] = "LIST"
		return inner
	}

	named := t.NamedType
	if named == "" {
		return map[string]interface{}{"kind": "SCALAR", "name": nil}
	}

	def := schema.Types[named]
	if def == nil {
		return map[string]interface{}{"kind": "SCALAR", "name": named}
	}

	result := e.buildTypeObject(schema, def)
	return result
}

// buildFields constructs __Field objects from an ast.FieldList.
func (e *graphQLEngine) buildFields(schema *ast.Schema, fields ast.FieldList, includeDeprecated bool) []interface{} {
	result := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		if !includeDeprecated && f.Directives.ForName("deprecated") != nil {
			continue
		}
		fieldObj := map[string]interface{}{
			"name": f.Name,
		}
		if f.Description != "" {
			fieldObj["description"] = f.Description
		}
		fieldObj["args"] = e.buildInputValues(schema, f.Arguments)
		fieldObj["type"] = e.buildTypeRefObject(schema, f.Type)
		depDir := f.Directives.ForName("deprecated")
		if depDir != nil {
			fieldObj["isDeprecated"] = true
			if reason := depDir.Arguments.ForName("reason"); reason != nil && reason.Value != nil {
				fieldObj["deprecationReason"] = fmt.Sprintf("%v", reason.Value)
			} else {
				fieldObj["deprecationReason"] = "No longer supported"
			}
		} else {
			fieldObj["isDeprecated"] = false
		}
		result = append(result, fieldObj)
	}
	return result
}

// buildInputValues constructs __InputValue objects from an ast.ArgumentDefinitionList.
func (e *graphQLEngine) buildInputValues(schema *ast.Schema, args ast.ArgumentDefinitionList) []interface{} {
	result := make([]interface{}, 0, len(args))
	for _, arg := range args {
		iv := map[string]interface{}{
			"name": arg.Name,
		}
		if arg.Description != "" {
			iv["description"] = arg.Description
		}
		iv["type"] = e.buildTypeRefObject(schema, arg.Type)
		if arg.DefaultValue != nil {
			iv["defaultValue"] = fmt.Sprintf("%v", arg.DefaultValue)
		}
		result = append(result, iv)
	}
	return result
}

// buildInputFields constructs __InputValue objects from an ast.FieldList (used for InputObject fields).
func (e *graphQLEngine) buildInputFields(schema *ast.Schema, fields ast.FieldList) []interface{} {
	result := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		iv := map[string]interface{}{
			"name": f.Name,
		}
		if f.Description != "" {
			iv["description"] = f.Description
		}
		iv["type"] = e.buildTypeRefObject(schema, f.Type)
		result = append(result, iv)
	}
	return result
}

// buildEnumValues constructs __EnumValue objects from an ast.EnumValueList.
func (e *graphQLEngine) buildEnumValues(values ast.EnumValueList, includeDeprecated bool) []interface{} {
	result := make([]interface{}, 0, len(values))
	for _, v := range values {
		if !includeDeprecated && v.Directives.ForName("deprecated") != nil {
			continue
		}
		ev := map[string]interface{}{
			"name": v.Name,
		}
		if v.Description != "" {
			ev["description"] = v.Description
		}
		ev["isDeprecated"] = v.Directives.ForName("deprecated") != nil
		result = append(result, ev)
	}
	return result
}

// buildDirectives constructs __Directive objects from the schema's directive list.
func (e *graphQLEngine) buildDirectives(schema *ast.Schema) []interface{} {
	result := make([]interface{}, 0, len(schema.Directives))
	for _, d := range schema.Directives {
		dir := map[string]interface{}{
			"name": d.Name,
		}
		if d.Description != "" {
			dir["description"] = d.Description
		}
		locs := make([]interface{}, 0, len(d.Locations))
		for _, loc := range d.Locations {
			locs = append(locs, string(loc))
		}
		dir["locations"] = locs
		dir["args"] = e.buildInputValues(schema, d.Arguments)
		result = append(result, dir)
	}
	return result
}

// typeKindString converts an ast.DefinitionKind to the GraphQL introspection __TypeKind enum value.
func typeKindString(kind ast.DefinitionKind) string {
	switch kind {
	case ast.Scalar:
		return "SCALAR"
	case ast.Object:
		return "OBJECT"
	case ast.Interface:
		return "INTERFACE"
	case ast.Union:
		return "UNION"
	case ast.Enum:
		return "ENUM"
	case ast.InputObject:
		return "INPUT_OBJECT"
	default:
		return "SCALAR"
	}
}

// sortedTypeNames returns the sorted list of user-defined type names from the schema,
// excluding introspection types (prefixed with __).
func sortedTypeNames(schema *ast.Schema) []string {
	names := make([]string, 0, len(schema.Types))
	for name := range schema.Types {
		if !strings.HasPrefix(name, "__") {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}
