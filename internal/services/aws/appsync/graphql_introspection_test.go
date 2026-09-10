package appsync

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

// A nullable list type must carry the element in Elem and leave the outer
// NamedType empty: isListType treats "Elem != nil && NamedType == \"\"" as the
// list marker, and any outer name routes list values down the composite
// completion path instead of iterating the elements.
func TestListTypeCarriesElementOnly(t *testing.T) {
	lt := listType("__Field")
	if lt.NamedType != "" {
		t.Fatalf("outer NamedType must be empty, got %q", lt.NamedType)
	}
	if lt.Elem == nil || lt.Elem.NamedType != "__Field" {
		t.Fatalf("element must carry the named type, got %+v", lt.Elem)
	}
	if !(&graphQLEngine{}).isListType(lt) {
		t.Fatal("isListType must recognise a listType value as a list")
	}
}

// The __Type meta-fields declared through listType (fields, interfaces,
// possibleTypes, enumValues, inputFields) must all resolve as list types.
func TestIntrospectionListMetaFieldsAreLists(t *testing.T) {
	schema := &ast.Schema{Types: map[string]*ast.Definition{}}
	injectIntrospectionTypes(schema)
	ty := schema.Types["__Type"]
	if ty == nil {
		t.Fatal("__Type not injected")
	}
	engine := &graphQLEngine{}
	for _, name := range []string{"fields", "interfaces", "possibleTypes", "enumValues", "inputFields"} {
		f := ty.Fields.ForName(name)
		if f == nil {
			t.Fatalf("__Type.%s missing", name)
		}
		if !engine.isListType(f.Type) {
			t.Fatalf("__Type.%s must be a list type, got %+v", name, f.Type)
		}
	}
}
