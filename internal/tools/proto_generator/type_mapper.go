// Package proto_generator provides utilities for generating Protocol Buffer definitions
// from Smithy models. This package handles type mapping and code generation.
package proto_generator

// TypeMapper converts Smithy primitive types to Protocol Buffer type equivalents.
type TypeMapper struct{}

// NewTypeMapper creates a new TypeMapper instance for converting Smithy types
// to Protocol Buffer types.
func NewTypeMapper() *TypeMapper {
	return &TypeMapper{}
}

// MapSmithyType converts a Smithy primitive type to its Protocol Buffer equivalent.
// This mapping standardises the type translation from Smithy's type system to
// Protocol Buffer's type system for code generation purposes.
func (m *TypeMapper) MapSmithyType(smithyType string) string {
	switch smithyType {
	case "string":
		return "string"
	case "integer":
		return "int32"
	case "long":
		return "int64"
	case "float":
		return "float"
	case "double":
		return "double"
	case "boolean":
		return "bool"
	case "blob":
		return "bytes"
	case "timestamp":
		return "string"
	case "biginteger":
		return "string"
	case "bigdecimal":
		return "string"
	case "byte":
		return "int32"
	case "short":
		return "int32"
	default:
		return "string"
	}
}
