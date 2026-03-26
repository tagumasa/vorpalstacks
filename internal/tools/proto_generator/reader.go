// Package proto_generator provides utilities for generating Protocol Buffer definitions
// from Smithy models. This package handles type mapping and code generation.
package proto_generator

import "context"

// ServiceInfo represents service metadata
type ServiceInfo struct {
	ID             int64
	Name           string
	ShapeID        string
	Version        string
	Protocol       string
	EndpointPrefix string
}

// ShapeInfo represents a shape (structure, enum, list, map, etc.)
type ShapeInfo struct {
	ID              int64
	ShapeID         string
	Type            string
	Documentation   string
	TargetShapeID   *int64
	ListMemberID    *int64
	MapKeyShapeID   *int64
	MapValueShapeID *int64
}

// OperationInfo represents an API operation with HTTP binding
type OperationInfo struct {
	ID            int64
	Name          string
	HTTPMethod    string
	HTTPURIPath   string
	HTTPURIQuery  string
	Documentation string
	InputShapeID  *int64
	OutputShapeID *int64
	// AWS Protocol specific
	ContentType    string
	XAmzTarget     string
	Protocol       string
	SmithyProtocol string
	XAmznQueryMode bool
	// Error information
	Errors []ErrorInfo
}

// ErrorInfo represents operation error details
type ErrorInfo struct {
	ID                int64
	ShapeID           string
	Name              string
	ErrorFault        string
	HTTPStatusCode    int
	ErrorCodeOverride string
}

// MemberInfo represents a member field with HTTP binding
type MemberInfo struct {
	ID            int64
	Name          string
	TargetShapeID int64
	IsRequired    bool
	// HTTP binding
	HTTPLabel               bool
	HTTPQuery               string
	HTTPHeader              string
	HTTPPayload             bool
	JSONName                string
	HTTPPrefix              string
	HTTPLocationName        string
	DefaultPayloadMediaType string
	// Documentation
	Documentation string
}

// EnumValueInfo represents an enum value
type EnumValueInfo struct {
	Name  string
	Index int
}

// ProtocolInfo represents protocol-specific information
type ProtocolInfo struct {
	Type           string
	ContentType    string
	XAmzTarget     string
	SmithyProtocol string
	XAmznQueryMode bool
}

// ServiceReader interface for reading service information
type ServiceReader interface {
	FindServiceByName(ctx context.Context, name string) (*ServiceInfo, error)
	ListServices(ctx context.Context) ([]*ServiceInfo, error)
}

// ShapeReader interface for reading shape information
type ShapeReader interface {
	FindShapeByID(ctx context.Context, id int64) (*ShapeInfo, error)
	FindShapesByServiceID(ctx context.Context, serviceID int64) ([]*ShapeInfo, error)
}

// OperationReader interface for reading operation information
type OperationReader interface {
	FindOperationsByServiceID(ctx context.Context, serviceID int64) ([]*OperationInfo, error)
}

// MemberReader interface for reading member information
type MemberReader interface {
	FindMembersByShapeID(ctx context.Context, shapeID int64) ([]*MemberInfo, error)
}

// EnumValueReader interface for reading enum values
type EnumValueReader interface {
	FindEnumValuesByShapeID(ctx context.Context, shapeID int64) ([]*EnumValueInfo, error)
}

// ProtoDataReader combines all reader interfaces
type ProtoDataReader interface {
	ServiceReader
	ShapeReader
	OperationReader
	MemberReader
	EnumValueReader
}
