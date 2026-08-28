package appsync

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

// typeInRequestedFormat returns the definition and format members a list
// response should carry for the requested output serialisation. When the
// stored serialisation differs from the requested one the definition is
// converted; a definition that cannot be converted is echoed unchanged in
// its stored serialisation.
func typeInRequestedFormat(t *appsyncstore.Type, format string) (string, string) {
	if t.Definition == "" || t.Format == format {
		return t.Definition, t.Format
	}
	switch {
	case t.Format == "SDL" && format == "JSON":
		if converted, ok := sdlTypeDefinitionToJSON(t.Definition); ok {
			return converted, format
		}
	case t.Format == "JSON" && format == "SDL":
		if converted, ok := jsonTypeDefinitionToSDL(t.Definition); ok {
			return converted, format
		}
	}
	return t.Definition, t.Format
}

// sdlTypeDefinitionToJSON parses one SDL type definition and serialises it
// as the introspection __Type JSON object.
func sdlTypeDefinitionToJSON(definition string) (string, bool) {
	schema, err := gqlparser.LoadSchema(&ast.Source{
		Name:  "schema.graphql",
		Input: definition,
	})
	if err != nil {
		return "", false
	}

	typeName := extractTypeName(definition)
	if typeName == "" {
		return "", false
	}
	def := schema.Types[typeName]
	if def == nil {
		return "", false
	}

	engine := &graphQLEngine{}
	data, err := json.Marshal(engine.buildTypeObject(schema, def))
	if err != nil {
		return "", false
	}
	return string(data), true
}

// jsonTypeDefinitionToSDL reconstructs SDL for one type from its
// introspection __Type JSON object.
func jsonTypeDefinitionToSDL(definition string) (string, bool) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(definition), &obj); err != nil {
		return "", false
	}

	name, _ := obj["name"].(string)
	kind, _ := obj["kind"].(string)
	if name == "" || kind == "" {
		return "", false
	}

	var b strings.Builder
	switch kind {
	case "SCALAR":
		fmt.Fprintf(&b, "scalar %s", name)

	case "OBJECT", "INTERFACE":
		keyword := "type"
		if kind == "INTERFACE" {
			keyword = "interface"
		}
		fmt.Fprintf(&b, "%s %s", keyword, name)
		if ifaces, ok := jsonTypeNameList(obj["interfaces"]); ok && len(ifaces) > 0 {
			fmt.Fprintf(&b, " implements %s", strings.Join(ifaces, " & "))
		}
		b.WriteString(" {\n")
		fields, ok := jsonFieldList(obj["fields"])
		if !ok {
			return "", false
		}
		for _, f := range fields {
			writeJSONDescription(&b, f)
			fName, _ := f["name"].(string)
			ref, ok := jsonTypeRefString(f["type"])
			if !ok {
				return "", false
			}
			b.WriteString("  " + fName)
			if args, ok := jsonInputValueList(f["args"]); ok && len(args) > 0 {
				b.WriteString("(")
				for i, a := range args {
					if i > 0 {
						b.WriteString(", ")
					}
					aName, _ := a["name"].(string)
					aRef, ok := jsonTypeRefString(a["type"])
					if !ok {
						return "", false
					}
					if def, ok := a["defaultValue"].(string); ok && def != "" {
						fmt.Fprintf(&b, "%s: %s = %s", aName, aRef, def)
					} else {
						fmt.Fprintf(&b, "%s: %s", aName, aRef)
					}
				}
				b.WriteString(")")
			}
			fmt.Fprintf(&b, ": %s\n", ref)
		}
		b.WriteString("}")

	case "ENUM":
		fmt.Fprintf(&b, "enum %s {\n", name)
		values, ok := jsonEnumValueList(obj["enumValues"])
		if !ok {
			return "", false
		}
		for _, v := range values {
			writeJSONDescription(&b, v)
			vName, _ := v["name"].(string)
			fmt.Fprintf(&b, "  %s\n", vName)
		}
		b.WriteString("}")

	case "INPUT_OBJECT":
		fmt.Fprintf(&b, "input %s {\n", name)
		fields, ok := jsonInputValueList(obj["inputFields"])
		if !ok {
			return "", false
		}
		for _, f := range fields {
			writeJSONDescription(&b, f)
			fName, _ := f["name"].(string)
			ref, ok := jsonTypeRefString(f["type"])
			if !ok {
				return "", false
			}
			fmt.Fprintf(&b, "  %s: %s\n", fName, ref)
		}
		b.WriteString("}")

	case "UNION":
		members, ok := jsonTypeNameList(obj["possibleTypes"])
		if !ok || len(members) == 0 {
			return "", false
		}
		fmt.Fprintf(&b, "union %s = %s", name, strings.Join(members, " | "))

	default:
		return "", false
	}

	return b.String(), true
}

// writeJSONDescription emits a non-empty introspection description member as
// an SDL triple-quoted block.
func writeJSONDescription(b *strings.Builder, obj map[string]interface{}) {
	if desc, ok := obj["description"].(string); ok && desc != "" {
		fmt.Fprintf(b, "  \"\"\"%s\"\"\"\n", desc)
	}
}

// jsonTypeNameList extracts the plain type names from an interfaces or
// possibleTypes array.
func jsonTypeNameList(raw interface{}) ([]string, bool) {
	items, ok := raw.([]interface{})
	if !ok {
		if raw == nil {
			return nil, true
		}
		return nil, false
	}
	names := make([]string, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, false
		}
		name, _ := m["name"].(string)
		if name == "" {
			return nil, false
		}
		names = append(names, name)
	}
	return names, true
}

// jsonFieldList extracts the field array of an object or interface type.
func jsonFieldList(raw interface{}) ([]map[string]interface{}, bool) {
	return jsonObjectList(raw)
}

// jsonInputValueList extracts an args or inputFields array.
func jsonInputValueList(raw interface{}) ([]map[string]interface{}, bool) {
	return jsonObjectList(raw)
}

// jsonEnumValueList extracts an enumValues array.
func jsonEnumValueList(raw interface{}) ([]map[string]interface{}, bool) {
	return jsonObjectList(raw)
}

// jsonObjectList coerces an introspection JSON array into a map slice; a
// nil member is an empty list, anything non-array is a conversion failure.
func jsonObjectList(raw interface{}) ([]map[string]interface{}, bool) {
	items, ok := raw.([]interface{})
	if !ok {
		if raw == nil {
			return nil, true
		}
		return nil, false
	}
	out := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, false
		}
		out = append(out, m)
	}
	return out, true
}

// jsonTypeRefString renders a __Type reference (NON_NULL/LIST wrappers over
// ofType) back into SDL type syntax.
func jsonTypeRefString(raw interface{}) (string, bool) {
	ref, ok := raw.(map[string]interface{})
	if !ok || ref == nil {
		return "", false
	}
	kind, _ := ref["kind"].(string)
	switch kind {
	case "NON_NULL":
		inner, ok := jsonTypeRefString(ref["ofType"])
		if !ok {
			return "", false
		}
		return inner + "!", true
	case "LIST":
		inner, ok := jsonTypeRefString(ref["ofType"])
		if !ok {
			return "", false
		}
		return "[" + inner + "]", true
	case "":
		return "", false
	default:
		name, _ := ref["name"].(string)
		if name == "" {
			return "", false
		}
		return name, true
	}
}
