// Package api provides storage implementations for API metadata including services, operations, and shapes.
package api

import "time"

// Service represents an AWS service definition stored in the runtime.
type Service struct {
	ID                    int64  `json:"id"`
	ShapeID               string `json:"shape_id"`
	Name                  string `json:"name"`
	Version               string `json:"version"`
	SDKID                 string `json:"sdk_id"`
	ARNNamespace          string `json:"arn_namespace"`
	CloudFormationName    string `json:"cloud_formation_name"`
	CloudTrailEventSource string `json:"cloud_trail_event_source"`
	EndpointPrefix        string `json:"endpoint_prefix"`
	Protocol              string `json:"protocol"`
	CreatedAt             int64  `json:"created_at"`
}

// Operation represents an API operation defined in a Smithy model.
type Operation struct {
	ID             int64   `json:"id"`
	ShapeID        string  `json:"shape_id"`
	ServiceID      int64   `json:"service_id"`
	Name           string  `json:"name"`
	Documentation  *string `json:"documentation,omitempty"`
	HTTPMethod     *string `json:"http_method,omitempty"`
	HTTPURI        *string `json:"http_uri,omitempty"`
	HTTPCode       *int    `json:"http_code,omitempty"`
	IsReadonly     bool    `json:"is_readonly"`
	InputShapeID   *int64  `json:"input_shape_id,omitempty"`
	OutputShapeID  *int64  `json:"output_shape_id,omitempty"`
	ContentType    *string `json:"content_type,omitempty"`
	XAmzTarget     *string `json:"x_amz_target,omitempty"`
	SmithyProtocol *string `json:"smithy_protocol,omitempty"`
	XAmznQueryMode bool    `json:"x_amzn_query_mode"`
	HTTPURIPath    *string `json:"http_uri_path,omitempty"`
	HTTPURIQuery   *string `json:"http_uri_query,omitempty"`
	CreatedAt      int64   `json:"created_at"`
}

// Shape represents a data structure defined in a Smithy model.
type Shape struct {
	ID                int64   `json:"id"`
	ShapeID           string  `json:"shape_id"`
	Type              string  `json:"type"`
	Documentation     *string `json:"documentation,omitempty"`
	ListMemberShapeID *int64  `json:"list_member_shape_id,omitempty"`
	MapKeyShapeID     *int64  `json:"map_key_shape_id,omitempty"`
	MapValueShapeID   *int64  `json:"map_value_shape_id,omitempty"`
	ServiceID         *int64  `json:"service_id,omitempty"`
	CreatedAt         int64   `json:"created_at"`
}

// Member represents a structure member (field) in a Smithy model.
type Member struct {
	ID                      int64   `json:"id"`
	ShapeID                 int64   `json:"shape_id"`
	Name                    string  `json:"name"`
	TargetShapeID           int64   `json:"target_shape_id"`
	Documentation           *string `json:"documentation,omitempty"`
	IsRequired              bool    `json:"is_required"`
	HTTPLabel               bool    `json:"http_label"`
	HTTPQuery               *string `json:"http_query,omitempty"`
	HTTPHeader              *string `json:"http_header,omitempty"`
	HTTPPayload             bool    `json:"http_payload"`
	JSONName                *string `json:"json_name,omitempty"`
	HTTPPrefix              *string `json:"http_prefix,omitempty"`
	HTTPLocationName        *string `json:"http_location_name,omitempty"`
	DefaultPayloadMediaType *string `json:"default_payload_media_type,omitempty"`
	CreatedAt               int64   `json:"created_at"`
}

// OperationError represents an error response that can be returned by an operation.
type OperationError struct {
	ID                int64  `json:"id"`
	OperationID       int64  `json:"operation_id"`
	ErrorShapeID      int64  `json:"error_shape_id"`
	ErrorFault        string `json:"error_fault"`
	HTTPStatusCode    int    `json:"http_status_code"`
	ErrorCodeOverride string `json:"error_code_override"`
	CreatedAt         int64  `json:"created_at"`
}

// ShapeTrait represents a trait annotation attached to a shape.
type ShapeTrait struct {
	ID         int64  `json:"id"`
	ShapeID    int64  `json:"shape_id"`
	TraitName  string `json:"trait_name"`
	TraitValue string `json:"trait_value"`
	CreatedAt  int64  `json:"created_at"`
}

// MemberTrait represents a trait annotation attached to a structure member.
type MemberTrait struct {
	ID         int64  `json:"id"`
	MemberID   int64  `json:"member_id"`
	TraitName  string `json:"trait_name"`
	TraitValue string `json:"trait_value"`
	CreatedAt  int64  `json:"created_at"`
}

// ShapeEnumValue represents a possible value for a string enum shape.
type ShapeEnumValue struct {
	ID            int64  `json:"id"`
	ShapeID       int64  `json:"shape_id"`
	Name          string `json:"name"`
	Value         string `json:"value"`
	Documentation string `json:"documentation"`
	CreatedAt     int64  `json:"created_at"`
}

// Resource represents an AWS resource defined in a Smithy model.
type Resource struct {
	ID        int64  `json:"id"`
	ShapeID   string `json:"shape_id"`
	ServiceID int64  `json:"service_id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}

// ServiceMode represents the mode of operation for a service in vorpalstacks.
type ServiceMode int

// Service mode constants defining how vorpalstacks handles requests.
const (
	ServiceModeImplemented    ServiceMode = iota // Real implementation
	ServiceModeFallback                          // Auto-generated responses for unimplemented operations
	ServiceModeErrorInjection                    // Error injection for testing
)

// String returns the string representation of a ServiceMode.
func (m ServiceMode) String() string {
	switch m {
	case ServiceModeImplemented:
		return "IMPLEMENTED"
	case ServiceModeFallback:
		return "FALLBACK"
	case ServiceModeErrorInjection:
		return "ERROR_INJECTION"
	default:
		return "UNKNOWN"
	}
}

// ParseServiceMode converts a string to a ServiceMode.
func ParseServiceMode(s string) ServiceMode {
	switch s {
	case "IMPLEMENTED":
		return ServiceModeImplemented
	case "FALLBACK":
		return ServiceModeFallback
	case "ERROR_INJECTION":
		return ServiceModeErrorInjection
	default:
		return ServiceModeFallback
	}
}

// ErrorConfig defines the error response configuration for error injection mode.
type ErrorConfig struct {
	HTTPStatusCode int    `json:"http_status_code"`
	Code           string `json:"code"`
	Message        string `json:"message"`
}

// ServiceConfig holds the configuration for a service in vorpalstacks.
type ServiceConfig struct {
	ServiceName string       `json:"service_name"`
	Mode        ServiceMode  `json:"mode"`
	Error       *ErrorConfig `json:"error,omitempty"`
	Enabled     bool         `json:"enabled"`
	UpdatedAt   int64        `json:"updated_at"`
}

// NewServiceConfig creates a new service configuration with default values.
func NewServiceConfig(serviceName string, mode ServiceMode) *ServiceConfig {
	return &ServiceConfig{
		ServiceName: serviceName,
		Mode:        mode,
		Enabled:     true,
		UpdatedAt:   time.Now().Unix(),
	}
}
