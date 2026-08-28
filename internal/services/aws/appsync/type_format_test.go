package appsync

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

func TestTypeFormatSDLToJSON(t *testing.T) {
	definition := "type Post { id: ID! title: String tags: [String!]! }"
	converted, ok := sdlTypeDefinitionToJSON(definition)
	if !ok {
		t.Fatalf("conversion failed for %q", definition)
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(converted), &obj); err != nil {
		t.Fatalf("converted definition is not JSON: %v (%q)", err, converted)
	}
	if obj["name"] != "Post" || obj["kind"] != "OBJECT" {
		t.Fatalf("unexpected name/kind: %v", obj)
	}
	fields, ok := obj["fields"].([]interface{})
	if !ok || len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %v", obj["fields"])
	}
}

func TestTypeFormatJSONToSDL(t *testing.T) {
	enumJSON, ok := sdlTypeDefinitionToJSON("enum Status { ACTIVE INACTIVE }")
	if !ok {
		t.Fatal("SDL to JSON conversion failed for the enum definition")
	}
	back, ok := jsonTypeDefinitionToSDL(enumJSON)
	if !ok {
		t.Fatalf("JSON to SDL conversion failed for %q", enumJSON)
	}
	if !strings.HasPrefix(strings.TrimSpace(back), "enum Status") {
		t.Fatalf("expected enum definition, got %q", back)
	}

	// An argument default value must survive the full round trip.
	definition := "type Query2 { posts(status: String = \"ACTIVE\", limit: Int): [String!] }"
	jsonDef, ok := sdlTypeDefinitionToJSON(definition)
	if !ok {
		t.Fatalf("SDL to JSON conversion failed for %q", definition)
	}
	back, ok = jsonTypeDefinitionToSDL(jsonDef)
	if !ok {
		t.Fatalf("JSON to SDL conversion failed for %q", jsonDef)
	}

	schema, err := gqlparser.LoadSchema(&ast.Source{Name: "schema.graphql", Input: back})
	if err != nil {
		t.Fatalf("reconstructed SDL does not parse: %v (%q)", err, back)
	}
	q2 := schema.Types["Query2"]
	if q2 == nil {
		t.Fatalf("Query2 missing from reconstructed schema: %q", back)
	}
	posts := q2.Fields.ForName("posts")
	if posts == nil {
		t.Fatalf("posts field lost: %q", back)
	}
	status := posts.Arguments.ForName("status")
	if status == nil || status.DefaultValue == nil {
		t.Fatalf("status argument or its default value lost: %q", back)
	}
	if posts.Type == nil || posts.Type.Elem == nil || !posts.Type.Elem.NonNull || posts.Type.Elem.NamedType != "String" {
		t.Fatalf("posts return type lost: %q", back)
	}
}

func TestTypeFormatUnconvertibleFallsBack(t *testing.T) {
	stored := &appsyncstore.Type{
		Name:       "Broken",
		Format:     "JSON",
		Definition: "not json at all",
	}
	definition, format := typeInRequestedFormat(stored, "SDL")
	if definition != stored.Definition || format != stored.Format {
		t.Fatalf("expected stored serialisation fallback, got %q/%s", definition, format)
	}
}

func TestTypeFormatSameFormatSkipsConversion(t *testing.T) {
	stored := &appsyncstore.Type{
		Name:       "Post",
		Format:     "SDL",
		Definition: "type Post { id: ID! }",
	}
	definition, format := typeInRequestedFormat(stored, "SDL")
	if definition != stored.Definition || format != "SDL" {
		t.Fatalf("expected stored serialisation, got %q/%s", definition, format)
	}
}
