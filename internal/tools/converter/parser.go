// Package converter provides types and utilities for parsing Smithy model files.
package converter

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
)

// ProtocolType represents the Smithy protocol type used by AWS services.
// Different protocols determine how requests and responses are formatted.
type ProtocolType string

// Protocol constants define the various Smithy protocol types supported by AWS services.
const (
	ProtocolRestXML    ProtocolType = "aws.protocols#restXml"
	ProtocolRestJSON1  ProtocolType = "aws.protocols#restJson1"
	ProtocolAwsJSON1_0 ProtocolType = "aws.protocols#awsJson1_0"
	ProtocolAwsJSON1_1 ProtocolType = "aws.protocols#awsJson1_1"
	ProtocolAwsQuery   ProtocolType = "aws.protocols#awsQuery"
	ProtocolEc2Query   ProtocolType = "aws.protocols#ec2Query"
	ProtocolRPCV2Cbor  ProtocolType = "aws.protocols#rpcv2Cbor"
)

// Parser parses Smithy JSON models and provides methods for extracting service information,
// protocol details, HTTP bindings, and trait values.
type Parser struct {
	model       *SmithyModel
	serviceID   string
	serviceName string
	protocol    ProtocolType
}

// NewParser creates a new Parser instance for parsing Smithy models.
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile reads and parses a Smithy JSON model file.
// It returns the parsed model and extracts service information and protocol details.
func (p *Parser) ParseFile(filePath string) (*SmithyModel, error) {
	// #nosec G304
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var model SmithyModel
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	p.model = &model

	// Extract service information
	p.extractServiceInfo()
	p.detectProtocol()

	return &model, nil
}

// GetModel returns the currently parsed Smithy model.
// Returns nil if no model has been parsed yet.
func (p *Parser) GetModel() *SmithyModel {
	return p.model
}

// GetServiceShape returns the service shape and its shape ID from the parsed model.
// The service shape defines the service's metadata, including the protocol and endpoint configuration.
func (p *Parser) GetServiceShape() (*Shape, string, error) {
	if p.model == nil {
		return nil, "", fmt.Errorf("model not parsed yet")
	}

	// Find the service shape (usually the first shape with type "service")
	for shapeID, shape := range p.model.Shapes {
		if shape.Type == "service" {
			return shape, shapeID, nil
		}
	}

	return nil, "", fmt.Errorf("service shape not found")
}

// GetServiceName returns the AWS service name extracted from the model.
// This is typically derived from the service shape ID namespace.
func (p *Parser) GetServiceName() string {
	if p.serviceName != "" {
		return p.serviceName
	}

	if _, shapeID, err := p.GetServiceShape(); err == nil {
		parts := strings.Split(shapeID, "#")
		if len(parts) > 0 {
			namespacePart := parts[0]
			namespaceParts := strings.Split(namespacePart, ".")
			for i := len(namespaceParts) - 1; i >= 0; i-- {
				if namespaceParts[i] != "" && namespaceParts[i] != "amazonaws" && namespaceParts[i] != "com" {
					p.serviceName = namespaceParts[i]
					return p.serviceName
				}
			}
		}
	}

	return "unknown"
}

// GetServiceID returns the service identifier from the model.
// This is typically in the format "com.amazonaws.service#ServiceName".
func (p *Parser) GetServiceID() string {
	if p.serviceID != "" {
		return p.serviceID
	}

	if _, shapeID, err := p.GetServiceShape(); err == nil {
		p.serviceID = shapeID
	}

	return p.serviceID
}

// GetProtocol returns the detected protocol type for the service.
func (p *Parser) GetProtocol() ProtocolType {
	return p.protocol
}

// detectProtocol examines the service shape to determine the protocol type.
// It checks for protocol traits in order of specificity.
func (p *Parser) detectProtocol() {
	serviceShape, _, err := p.GetServiceShape()
	if err != nil {
		return
	}

	// Check for protocol traits in order of specificity
	// Order matters: more specific protocols should be checked first
	protocolTraits := []ProtocolType{
		ProtocolRPCV2Cbor,
		ProtocolRestXML,
		ProtocolRestJSON1,
		ProtocolAwsJSON1_0,
		ProtocolAwsJSON1_1,
		ProtocolAwsQuery,
		ProtocolEc2Query,
	}

	for _, protocol := range protocolTraits {
		// Check if the protocol trait exists (can be boolean, map, or any value)
		if traitValue, exists := serviceShape.Traits[string(protocol)]; exists {
			if _, ok := traitValue.(bool); ok && !traitValue.(bool) {
				// Skip if explicitly set to false
				continue
			}
			p.protocol = protocol
			return
		}
	}
}

// GetContentType returns the Content-Type for the protocol
func (p *Parser) GetContentType() string {
	switch p.protocol {
	case ProtocolRestXML:
		return "application/xml"
	case ProtocolRestJSON1:
		return "application/json"
	case ProtocolAwsJSON1_0:
		return "application/x-amz-json-1.0"
	case ProtocolAwsJSON1_1:
		return "application/x-amz-json-1.1"
	case ProtocolAwsQuery, ProtocolEc2Query:
		return "application/x-www-form-urlencoded"
	case ProtocolRPCV2Cbor:
		return "application/cbor"
	default:
		return ""
	}
}

// GetContentTypeForOperation returns the Content-Type for a specific operation
// This handles special cases like Lambda Invoke which uses application/octet-stream
func (p *Parser) GetContentTypeForOperation(operationTraits map[string]interface{}) string {
	// Check for defaultPayloadMediaType trait (used by Lambda Invoke)
	if mediaType := p.GetDefaultPayloadMediaType(operationTraits); mediaType != "" {
		return mediaType
	}

	// Check for service-specific Content-Type overrides
	// Some services (Bedrock Runtime, EKS, etc.) use application/json instead of
	// the protocol default (application/x-amz-json-1.1 for AWS JSON 1.1)
	if p.shouldUseApplicationJSON() {
		return "application/json"
	}

	// Fall back to protocol default
	return p.GetContentType()
}

// shouldUseApplicationJSON checks if the service should use application/json
// instead of the protocol default Content-Type
func (p *Parser) shouldUseApplicationJSON() bool {
	if p.serviceID == "" {
		return false
	}

	// Services that use application/json despite using AWS JSON 1.1 protocol
	// Note: These services may use different protocols (REST-JSON 1.1, AWS JSON 1.1, etc.)
	// The actual Content-Type depends on the protocol and service-specific settings
	servicesUsingApplicationJSON := map[string]bool{
		"com.amazonaws.bedrockruntime#AmazonBedrockFrontendService": true,
		"com.amazonaws.eks#AmazonEKS":                               true,
		"com.amazonaws.sagemakerruntime#AmazonSageMakerRuntime":     true,
	}

	return servicesUsingApplicationJSON[p.serviceID]
}

// GetXAmzTarget generates the X-Amz-Target header value
func (p *Parser) GetXAmzTarget(operationShapeID string) string {
	// X-Amz-Target is only used for AWS JSON 1.0 and 1.1
	if p.protocol != ProtocolAwsJSON1_0 && p.protocol != ProtocolAwsJSON1_1 {
		return ""
	}

	// Format: {ServiceShapeName}.{OperationShapeName}
	// Extract service shape name from service ID (e.g., "com.amazonaws.lambda#LambdaService" -> "LambdaService")
	serviceShapeName := p.extractShapeName(p.serviceID)
	if serviceShapeName == "" {
		return ""
	}

	// Extract operation shape name from operation shape ID (e.g., "com.amazonaws.lambda#Invoke" -> "Invoke")
	operationShapeName := p.extractShapeName(operationShapeID)
	if operationShapeName == "" {
		return ""
	}

	return fmt.Sprintf("%s.%s", serviceShapeName, operationShapeName)
}

// GetSmithyProtocol returns the smithy-protocol header value
func (p *Parser) GetSmithyProtocol() string {
	switch p.protocol {
	case ProtocolRPCV2Cbor:
		return "rpc-v2-cbor"
	default:
		return ""
	}
}

// GetXAmznQueryMode returns the X-Amzn-Query-Mode header value
func (p *Parser) GetXAmznQueryMode(serviceShape *Shape) bool {
	if serviceShape == nil || serviceShape.Traits == nil {
		return false
	}

	// Check for aws.api#awsQueryCompatible trait (can be boolean or map)
	if queryCompatibleTrait, ok := serviceShape.Traits["aws.api#awsQueryCompatible"]; ok {
		if b, ok := queryCompatibleTrait.(bool); ok {
			return b
		} else if _, ok := queryCompatibleTrait.(map[string]interface{}); ok {
			return true
		}
	}

	return false
}

// ParseURI parses the URI and splits it into path and query parts
func (p *Parser) ParseURI(uri string) (path, query string) {
	if uri == "" {
		return "", ""
	}

	// Parse the URI
	parsedURL, err := url.Parse(uri)
	if err != nil {
		// If parsing fails, return the URI as path
		// This handles simple URI paths like "/resources/{id}"
		return uri, ""
	}

	// Extract path and query
	path = parsedURL.Path
	query = parsedURL.Query().Encode()

	// If the original URI contains query parameters but parsedURL.Query is empty,
	// try to extract them manually
	if query == "" && strings.Contains(uri, "?") {
		parts := strings.SplitN(uri, "?", 2)
		if len(parts) == 2 {
			query = parts[1]
		}
	}

	return path, query
}

// GetDefaultPayloadMediaType extracts the default payload media type from traits
func (p *Parser) GetDefaultPayloadMediaType(traits map[string]interface{}) string {
	return extractStringTrait(traits, "smithy.http#defaultPayloadMediaType")
}

// GetHTTPLocationName extracts the HTTP location name from traits
func (p *Parser) GetHTTPLocationName(traits map[string]interface{}) string {
	if traits == nil {
		return ""
	}

	if _, ok := traits["smithy.api#httpLabel"]; ok {
		return "label"
	}

	if val := extractStringTrait(traits, "smithy.api#httpQuery"); val != "" {
		return val
	}

	if val := extractStringTrait(traits, "smithy.api#httpHeader"); val != "" {
		return val
	}

	if _, ok := traits["smithy.api#httpPayload"]; ok {
		return "payload"
	}

	return ""
}

// GetHTTPPrefix extracts the HTTP prefix from traits
func (p *Parser) GetHTTPPrefix(traits map[string]interface{}) string {
	return extractStringTrait(traits, "smithy.api#httpPrefix")
}

// GetErrorFault extracts the error fault from traits
func (p *Parser) GetErrorFault(traits map[string]interface{}) string {
	return extractStringTrait(traits, "smithy.api#error")
}

// GetHTTPErrorCode extracts the HTTP error code from traits
func (p *Parser) GetHTTPErrorCode(traits map[string]interface{}) int {
	if val := extractIntTrait(traits, "smithy.api#httpError"); val != 0 {
		return val
	}
	return extractIntTrait(traits, "smithy.api#httpCode")
}

// GetErrorCodeOverride extracts the error code override from traits
func (p *Parser) GetErrorCodeOverride(traits map[string]interface{}) string {
	if traits == nil {
		return ""
	}

	return extractStringTrait(traits, "aws.api#errorCode")
}

// GetDocumentation extracts documentation from traits
func (p *Parser) GetDocumentation(traits map[string]interface{}) string {
	if val := extractStringTrait(traits, "docs"); val != "" {
		return val
	}
	return extractStringTrait(traits, "documentation")
}

// GetHTTPBinding extracts HTTP binding from traits
func (p *Parser) GetHTTPBinding(traits map[string]interface{}) *HTTPBinding {
	return extractHTTPBinding(traits, []string{"smithy.api#http", "http"})
}

// GetTraitValue extracts a trait value
func (p *Parser) GetTraitValue(traits map[string]interface{}, traitName string) interface{} {
	if traits == nil {
		return nil
	}

	return traits[traitName]
}

// GetTraitString extracts a string trait value
func (p *Parser) GetTraitString(traits map[string]interface{}, traitName string) string {
	return extractStringTrait(traits, traitName)
}

// GetTraitBool extracts a boolean trait value
func (p *Parser) GetTraitBool(traits map[string]interface{}, traitName string) bool {
	return extractBoolTrait(traits, traitName)
}

// extractShapeName extracts the shape name from a shape ID
func (p *Parser) extractShapeName(shapeID string) string {
	return extractShapeName(shapeID)
}

// extractServiceInfo extracts service information from the model
func (p *Parser) extractServiceInfo() {
	_, shapeID, _ := p.GetServiceShape()
	p.serviceID = shapeID
	p.serviceName = p.GetServiceName()
}
