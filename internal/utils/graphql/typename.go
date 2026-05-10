// Package graphql provides utility functions for GraphQL SDL processing.
package graphql

import "strings"

// graphQLTypePrefixes lists the keyword prefixes that introduce a named type in SDL.
var graphQLTypePrefixes = []string{
	"type ", "input ", "enum ", "interface ", "union ", "scalar ",
	"extend type ", "extend input ", "extend interface ", "extend enum ", "extend union ", "extend scalar ",
}

// ExtractTypeName parses a GraphQL SDL definition and extracts the type name.
// Handles formats like "type Post { ... }", "input PostInput { ... }", "enum Status { ... }",
// "extend type Post { ... }", etc.
func ExtractTypeName(def string) string {
	for _, prefix := range graphQLTypePrefixes {
		if idx := strings.Index(def, prefix); idx != -1 {
			rest := def[idx+len(prefix):]
			rest = strings.TrimLeft(rest, " ")
			end := strings.IndexAny(rest, " \t\n{(")
			if end > 0 {
				return rest[:end]
			}
		}
	}
	return ""
}
