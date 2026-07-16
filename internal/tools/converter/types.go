// Package converter provides types and utilities for parsing Smithy model files.
// Smithy is a language for defining services and SDKs, and this package handles
// parsing and processing of Smithy model definitions.
package converter

// SmithyModel represents a complete parsed Smithy model file.
// It contains the Smithy version, metadata, and all shape definitions.
type SmithyModel struct {
	Smithy   string                 `json:"smithy"`
	Version  string                 `json:"version"`
	Metadata map[string]interface{} `json:"metadata"`
	Shapes   map[string]*Shape      `json:"shapes"`
}

// Shape represents a Smithy shape definition.
// Shapes are the fundamental building blocks for data types in Smithy models,
// including structures, unions, lists, maps, and enums.
type Shape struct {
	Type    string                 `json:"type"`
	Traits  map[string]interface{} `json:"traits"`
	Member  *Member                `json:"member,omitempty"`
	Members map[string]*Member     `json:"members,omitempty"`
	Key     *Member                `json:"key,omitempty"`
	Value   *Member                `json:"value,omitempty"`
	Enum    []EnumValue            `json:"enum,omitempty"`
	Input   *ShapeRef              `json:"input,omitempty"`
	Output  *ShapeRef              `json:"output,omitempty"`
}

// ShapeRef represents a reference to another shape, typically used for
// input and output shapes in operation definitions.
type ShapeRef struct {
	Target string `json:"target,omitempty"`
}

// Member represents a member (field) of a structure or union shape.
// Members define the properties of complex types and include their target shape.
type Member struct {
	Target string                 `json:"target"`
	Traits map[string]interface{} `json:"traits"`
}

// EnumValue represents a single value in a string enum shape.
// Enum values define the possible values that an enum type can take.
type EnumValue struct {
	Value  string                 `json:"value"`
	Name   string                 `json:"name,omitempty"`
	Traits map[string]interface{} `json:"traits,omitempty"`
}

// HTTPBinding represents HTTP binding information for an operation.
// This includes the HTTP method, URI, and optional status code override.
type HTTPBinding struct {
	Method string `json:"method"`
	URI    string `json:"uri"`
	Code   int    `json:"code"`
}
