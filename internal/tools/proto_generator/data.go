// Package proto_generator provides utilities for generating Protocol Buffer definitions
// from Smithy models. This package handles type mapping and code generation.
package proto_generator

// ProtoTemplateData holds data for proto template generation
type ProtoTemplateData struct {
	PackageName string
	ServiceName string
	GoPackage   string
	DartPackage string
	Operations  []OperationData
	Shapes      []ShapeData
}

// OperationData holds operation information for template
type OperationData struct {
	Name           string
	InputShape     string
	OutputShape    string
	Comment        string
	HTTPMethod     string
	HTTPURIPath    string
	HTTPURIQuery   string
	ContentType    string
	XAmzTarget     string
	Protocol       string
	SmithyProtocol string
	XAmznQueryMode bool
	Errors         []ErrorData
}

// ErrorData holds error information for template
type ErrorData struct {
	Name              string
	ErrorFault        string
	HTTPStatusCode    int
	ErrorCodeOverride string
}

// ShapeData holds shape information for template
type ShapeData struct {
	Name       string
	Fields     []FieldData
	Comment    string
	IsEnum     bool
	EnumValues []*EnumValueInfo
}

// FieldData holds field information for template
type FieldData struct {
	Name                    string
	NameLower               string
	Type                    string
	Number                  int
	IsRequired              bool
	HTTPLabel               bool
	HTTPQuery               string
	HTTPHeader              string
	HTTPPayload             bool
	JSONName                string
	HTTPPrefix              string
	HTTPLocationName        string
	DefaultPayloadMediaType string
}
