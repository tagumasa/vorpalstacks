// Package smithy_reader provides functionality for reading Smithy JSON models directly.
package smithy_reader

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	proto_generator "vorpalstacks/internal/tools/proto_generator"
)

// SmithyReader reads Smithy JSON models directly
type SmithyReader struct {
	model      *SmithyModel
	serviceKey string
	service    *Shape
	protocol   string
}

// NewSmithyReader creates a new SmithyReader from a JSON file
func NewSmithyReader(jsonPath string) (*SmithyReader, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read smithy file: %w", err)
	}

	var model SmithyModel
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, fmt.Errorf("failed to parse smithy json: %w", err)
	}

	var serviceKey string
	var service *Shape
	for key, shape := range model.Shapes {
		if shape.Type == "service" {
			serviceKey = key
			service = &shape
			break
		}
	}

	if service == nil {
		return nil, fmt.Errorf("no service found in smithy model")
	}

	reader := &SmithyReader{
		model:      &model,
		serviceKey: serviceKey,
		service:    service,
	}

	// Detect protocol from service
	reader.protocol = reader.detectProtocol()

	return reader, nil
}

// GetServiceName returns the service name
func (r *SmithyReader) GetServiceName() string {
	if r.service.Traits != nil {
		if apiService, ok := r.service.Traits["aws.api#service"].(map[string]interface{}); ok {
			if sdkId, ok := apiService["sdkId"].(string); ok && sdkId != "" {
				return sdkId
			}
		}
	}
	if idx := strings.Index(r.serviceKey, "#"); idx >= 0 {
		return r.serviceKey[idx+1:]
	}
	return r.serviceKey
}

// FindServiceByName finds a service by name
func (r *SmithyReader) FindServiceByName(ctx context.Context, name string) (*proto_generator.ServiceInfo, error) {
	serviceName := r.GetServiceName()
	if !strings.EqualFold(serviceName, name) && !strings.EqualFold(r.serviceKey, name) {
		return nil, nil
	}

	endpointPrefix := strings.ToLower(serviceName)
	sdkId := serviceName
	if r.service.Traits != nil {
		if apiService, ok := r.service.Traits["aws.api#service"].(map[string]interface{}); ok {
			if ep, ok := apiService["endpointPrefix"].(string); ok {
				endpointPrefix = ep
			}
			if sid, ok := apiService["sdkId"].(string); ok {
				sdkId = sid
			}
		}
	}

	return &proto_generator.ServiceInfo{
		ID:             1,
		Name:           sdkId,
		ShapeID:        r.serviceKey,
		Version:        r.service.Version,
		Protocol:       r.protocol,
		EndpointPrefix: endpointPrefix,
	}, nil
}

// ListServices lists available services
func (r *SmithyReader) ListServices(ctx context.Context) ([]*proto_generator.ServiceInfo, error) {
	serviceName := r.GetServiceName()
	endpointPrefix := strings.ToLower(serviceName)
	sdkId := serviceName
	if r.service.Traits != nil {
		if apiService, ok := r.service.Traits["aws.api#service"].(map[string]interface{}); ok {
			if ep, ok := apiService["endpointPrefix"].(string); ok {
				endpointPrefix = ep
			}
			if sid, ok := apiService["sdkId"].(string); ok {
				sdkId = sid
			}
		}
	}

	return []*proto_generator.ServiceInfo{
		{
			ID:             1,
			Name:           sdkId,
			ShapeID:        r.serviceKey,
			Version:        r.service.Version,
			Protocol:       r.protocol,
			EndpointPrefix: endpointPrefix,
		},
	}, nil
}

// FindShapeByID finds a shape by ID
func (r *SmithyReader) FindShapeByID(ctx context.Context, id int64) (*proto_generator.ShapeInfo, error) {
	return r.FindShapeByHashID(id)
}

// FindShapeByKey finds a shape by key
func (r *SmithyReader) FindShapeByKey(key string) (*proto_generator.ShapeInfo, error) {
	shape, ok := r.model.Shapes[key]
	if !ok {
		return nil, nil
	}

	var listMemberID, mapKeyID, mapValueID *int64
	if shape.Member != nil {
		tmp := int64(hashKey(shape.Member.Target))
		listMemberID = &tmp
	}
	if shape.Key != nil {
		tmp := int64(hashKey(shape.Key.Target))
		mapKeyID = &tmp
	}
	if shape.Value != nil {
		tmp := int64(hashKey(shape.Value.Target))
		mapValueID = &tmp
	}

	return &proto_generator.ShapeInfo{
		ID:              int64(hashKey(key)),
		ShapeID:         key,
		Type:            shape.Type,
		Documentation:   getDocumentation(shape.Traits),
		ListMemberID:    listMemberID,
		MapKeyShapeID:   mapKeyID,
		MapValueShapeID: mapValueID,
	}, nil
}

// FindShapesByServiceID finds all shapes for a service
func (r *SmithyReader) FindShapesByServiceID(ctx context.Context, serviceID int64) ([]*proto_generator.ShapeInfo, error) {
	var shapes []*proto_generator.ShapeInfo

	for key, shape := range r.model.Shapes {
		switch shape.Type {
		case "structure", "union", "error", "enum", "list", "map":
			var listMemberID, mapKeyID, mapValueID *int64
			if shape.Member != nil {
				tmp := int64(hashKey(shape.Member.Target))
				listMemberID = &tmp
			}
			if shape.Key != nil {
				tmp := int64(hashKey(shape.Key.Target))
				mapKeyID = &tmp
			}
			if shape.Value != nil {
				tmp := int64(hashKey(shape.Value.Target))
				mapValueID = &tmp
			}

			shapes = append(shapes, &proto_generator.ShapeInfo{
				ID:              int64(hashKey(key)),
				ShapeID:         key,
				Type:            shape.Type,
				Documentation:   getDocumentation(shape.Traits),
				ListMemberID:    listMemberID,
				MapKeyShapeID:   mapKeyID,
				MapValueShapeID: mapValueID,
			})
		}
	}

	sort.Slice(shapes, func(i, j int) bool {
		return shapes[i].ShapeID < shapes[j].ShapeID
	})

	return shapes, nil
}

// FindOperationsByServiceID finds all operations for a service
func (r *SmithyReader) FindOperationsByServiceID(ctx context.Context, serviceID int64) ([]*proto_generator.OperationInfo, error) {
	seenOps := make(map[string]bool)
	var operations []*proto_generator.OperationInfo

	collectOperation := func(opRef OperationRef) {
		if seenOps[opRef.Target] {
			return
		}

		opShape, ok := r.model.Shapes[opRef.Target]
		if !ok || opShape.Type != "operation" {
			return
		}

		seenOps[opRef.Target] = true

		opName := extractName(opRef.Target)
		var inputID, outputID *int64
		if opShape.Input != nil {
			tmp := int64(hashKey(opShape.Input.Target))
			inputID = &tmp
		}
		if opShape.Output != nil {
			tmp := int64(hashKey(opShape.Output.Target))
			outputID = &tmp
		}

		httpMethod, httpURIPath, httpURIQuery := "", "", ""
		if httpTrait, ok := opShape.Traits["smithy.api#http"].(map[string]interface{}); ok {
			if m, ok := httpTrait["method"].(string); ok {
				httpMethod = m
			}
			if u, ok := httpTrait["uri"].(string); ok {
				httpURIPath, httpURIQuery = parseURI(u)
			}
		}

		if httpMethod == "" && (r.protocol == "awsQuery" || r.protocol == "ec2Query") {
			httpMethod = "POST"
			httpURIPath = "/"
		}

		contentType := r.getContentTypeForOperation(opShape.Traits)
		xAmzTarget := r.getXAmzTarget(opRef.Target)
		smithyProtocol := r.getSmithyProtocol()
		xAmznQueryMode := r.getXAmznQueryMode()
		errors := r.extractOperationErrors(opShape)

		operations = append(operations, &proto_generator.OperationInfo{
			ID:             int64(hashKey(opRef.Target)),
			Name:           opName,
			HTTPMethod:     httpMethod,
			HTTPURIPath:    httpURIPath,
			HTTPURIQuery:   httpURIQuery,
			Documentation:  getDocumentation(opShape.Traits),
			InputShapeID:   inputID,
			OutputShapeID:  outputID,
			ContentType:    contentType,
			XAmzTarget:     xAmzTarget,
			Protocol:       r.protocol,
			SmithyProtocol: smithyProtocol,
			XAmznQueryMode: xAmznQueryMode,
			Errors:         errors,
		})
	}

	for _, opRef := range r.service.Operations {
		collectOperation(opRef)
	}

	for _, resRef := range r.service.Resources {
		resShape, ok := r.model.Shapes[resRef.Target]
		if !ok || resShape.Type != "resource" {
			continue
		}

		if resShape.List != nil {
			collectOperation(OperationRef{Target: resShape.List.Target})
		}
		if resShape.Put != nil {
			collectOperation(OperationRef{Target: resShape.Put.Target})
		}
		if resShape.Create != nil {
			collectOperation(OperationRef{Target: resShape.Create.Target})
		}
		if resShape.Read != nil {
			collectOperation(OperationRef{Target: resShape.Read.Target})
		}
		if resShape.Update != nil {
			collectOperation(OperationRef{Target: resShape.Update.Target})
		}
		if resShape.Delete != nil {
			collectOperation(OperationRef{Target: resShape.Delete.Target})
		}
		for _, opRef := range resShape.Operations {
			collectOperation(opRef)
		}
	}

	return operations, nil
}

// FindMembersByShapeID finds all members for a shape
func (r *SmithyReader) FindMembersByShapeID(ctx context.Context, shapeID int64) ([]*proto_generator.MemberInfo, error) {
	key := r.findKeyByHash(shapeID)
	if key == "" {
		return nil, nil
	}

	shape, ok := r.model.Shapes[key]
	if !ok || shape.Members == nil {
		return nil, nil
	}

	var members []*proto_generator.MemberInfo
	for memberName, member := range shape.Members {
		targetID := int64(hashKey(member.Target))

		// Extract traits
		isRequired := hasTrait(member.Traits, "smithy.api#required")
		httpLabel := hasTrait(member.Traits, "smithy.api#httpLabel")
		httpPayload := hasTrait(member.Traits, "smithy.api#httpPayload")

		httpQuery := getStringTrait(member.Traits, "smithy.api#httpQuery")
		httpHeader := getStringTrait(member.Traits, "smithy.api#httpHeader")
		jsonName := getStringTrait(member.Traits, "smithy.api#jsonName")
		httpPrefix := getStringTrait(member.Traits, "smithy.api#httpPrefix")
		httpLocationName := getStringTrait(member.Traits, "smithy.api#httpLocationName")
		defaultPayloadMediaType := getStringTrait(member.Traits, "smithy.http#defaultPayloadMediaType")
		documentation := getDocumentation(member.Traits)

		members = append(members, &proto_generator.MemberInfo{
			ID:                      int64(hashKey(key + "#" + memberName)),
			Name:                    memberName,
			TargetShapeID:           targetID,
			IsRequired:              isRequired,
			HTTPLabel:               httpLabel,
			HTTPQuery:               httpQuery,
			HTTPHeader:              httpHeader,
			HTTPPayload:             httpPayload,
			JSONName:                jsonName,
			HTTPPrefix:              httpPrefix,
			HTTPLocationName:        httpLocationName,
			DefaultPayloadMediaType: defaultPayloadMediaType,
			Documentation:           documentation,
		})
	}

	sort.Slice(members, func(i, j int) bool {
		return members[i].Name < members[j].Name
	})

	return members, nil
}

// FindEnumValuesByShapeID finds all enum values for a shape
func (r *SmithyReader) FindEnumValuesByShapeID(ctx context.Context, shapeID int64) ([]*proto_generator.EnumValueInfo, error) {
	key := r.findKeyByHash(shapeID)
	if key == "" {
		return nil, nil
	}

	shape, ok := r.model.Shapes[key]
	if !ok {
		return nil, nil
	}

	var enumValues []*proto_generator.EnumValueInfo
	type enumEntry struct {
		name  string
		index int
	}
	var entries []enumEntry

	if len(shape.Enum) > 0 {
		for i, name := range shape.Enum {
			entries = append(entries, enumEntry{name: name, index: stableEnumIndex(name, i)})
		}
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].index != entries[j].index {
				return entries[i].index < entries[j].index
			}
			return entries[i].name < entries[j].name
		})
		for _, e := range entries {
			enumValues = append(enumValues, &proto_generator.EnumValueInfo{
				Name:  toProtoEnumName(e.name),
				Index: e.index,
			})
		}
		return enumValues, nil
	}

	if shape.Members != nil && shape.Type == "enum" {
		for memberName := range shape.Members {
			entries = append(entries, enumEntry{name: memberName, index: stableEnumIndex(memberName, 0)})
		}
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].index != entries[j].index {
				return entries[i].index < entries[j].index
			}
			return entries[i].name < entries[j].name
		})
		for _, e := range entries {
			enumValues = append(enumValues, &proto_generator.EnumValueInfo{
				Name:  toProtoEnumName(e.name),
				Index: e.index,
			})
		}
	}

	return enumValues, nil
}

func stableEnumIndex(name string, fallback int) int {
	h := uint32(2166136261)
	for _, c := range name {
		h ^= uint32(c)
		h *= 16777619
	}
	return int(h%1000) + 1
}

// FindShapeByHashID finds a shape by hash ID
func (r *SmithyReader) FindShapeByHashID(hashID int64) (*proto_generator.ShapeInfo, error) {
	key := r.findKeyByHash(hashID)
	if key == "" {
		return nil, nil
	}
	return r.FindShapeByKey(key)
}

// findKeyByHash finds a shape key by its hash
func (r *SmithyReader) findKeyByHash(hashID int64) string {
	for key := range r.model.Shapes {
		if int64(hashKey(key)) == hashID {
			return key
		}
	}
	return ""
}

// detectProtocol detects the protocol from the service shape
func (r *SmithyReader) detectProtocol() string {
	if r.service.Traits == nil {
		return ""
	}

	protocolTraits := []string{
		"aws.protocols#rpcv2Cbor",
		"aws.protocols#restXml",
		"aws.protocols#restJson1",
		"aws.protocols#awsJson1_0",
		"aws.protocols#awsJson1_1",
		"aws.protocols#awsQuery",
		"aws.protocols#ec2Query",
	}

	for _, protocol := range protocolTraits {
		if traitValue, exists := r.service.Traits[protocol]; exists {
			if b, ok := traitValue.(bool); ok && !b {
				continue
			}
			switch protocol {
			case "aws.protocols#rpcv2Cbor":
				return "rpcv2Cbor"
			case "aws.protocols#restXml":
				return "restXml"
			case "aws.protocols#restJson1":
				return "restJson1"
			case "aws.protocols#awsJson1_0":
				return "awsJson1_0"
			case "aws.protocols#awsJson1_1":
				return "awsJson1_1"
			case "aws.protocols#awsQuery":
				return "awsQuery"
			case "aws.protocols#ec2Query":
				return "ec2Query"
			}
		}
	}

	return ""
}

// getContentType returns the Content-Type for the protocol
func (r *SmithyReader) getContentType() string {
	switch r.protocol {
	case "restXml":
		return "application/xml"
	case "restJson1":
		return "application/json"
	case "awsJson1_0":
		return "application/x-amz-json-1.0"
	case "awsJson1_1":
		return "application/x-amz-json-1.1"
	case "awsQuery", "ec2Query":
		return "application/x-www-form-urlencoded"
	case "rpcv2Cbor":
		return "application/cbor"
	default:
		return ""
	}
}

// getContentTypeForOperation returns the Content-Type for a specific operation
func (r *SmithyReader) getContentTypeForOperation(operationTraits map[string]interface{}) string {
	// Check for defaultPayloadMediaType trait
	if mediaType := getStringTrait(operationTraits, "smithy.http#defaultPayloadMediaType"); mediaType != "" {
		return mediaType
	}

	// Check for service-specific Content-Type overrides
	if r.shouldUseApplicationJSON() {
		return "application/json"
	}

	return r.getContentType()
}

// shouldUseApplicationJSON checks if the service should use application/json
func (r *SmithyReader) shouldUseApplicationJSON() bool {
	servicesUsingApplicationJSON := map[string]bool{
		"com.amazonaws.bedrockruntime#AmazonBedrockFrontendService": true,
		"com.amazonaws.eks#AmazonEKS":                               true,
		"com.amazonaws.sagemakerruntime#AmazonSageMakerRuntime":     true,
	}

	return servicesUsingApplicationJSON[r.serviceKey]
}

// getXAmzTarget generates the X-Amz-Target header value
func (r *SmithyReader) getXAmzTarget(operationShapeID string) string {
	if r.protocol != "awsJson1_0" && r.protocol != "awsJson1_1" {
		return ""
	}

	serviceShapeName := extractName(r.serviceKey)
	if serviceShapeName == "" {
		return ""
	}

	operationShapeName := extractName(operationShapeID)
	if operationShapeName == "" {
		return ""
	}

	return fmt.Sprintf("%s.%s", serviceShapeName, operationShapeName)
}

// getSmithyProtocol returns the smithy-protocol header value
func (r *SmithyReader) getSmithyProtocol() string {
	if r.protocol == "rpcv2Cbor" {
		return "rpc-v2-cbor"
	}
	return ""
}

// getXAmznQueryMode returns the X-Amzn-Query-Mode header value
func (r *SmithyReader) getXAmznQueryMode() bool {
	if r.service.Traits == nil {
		return false
	}

	if queryCompatibleTrait, ok := r.service.Traits["aws.api#awsQueryCompatible"]; ok {
		if b, ok := queryCompatibleTrait.(bool); ok {
			return b
		}
		if _, ok := queryCompatibleTrait.(map[string]interface{}); ok {
			return true
		}
	}

	return false
}

// extractOperationErrors extracts error information from operation
func (r *SmithyReader) extractOperationErrors(opShape Shape) []proto_generator.ErrorInfo {
	var errors []proto_generator.ErrorInfo

	// Extract errors from traits
	var errorsTrait interface{}
	var ok bool

	if errorsTrait, ok = opShape.Traits["smithy.api#errors"]; !ok {
		if errorsTrait, ok = opShape.Traits["errors"]; !ok {
			return errors
		}
	}

	// Handle different error formats
	var errorTargets []string
	switch v := errorsTrait.(type) {
	case []interface{}:
		for _, e := range v {
			if target, ok := e.(string); ok {
				errorTargets = append(errorTargets, target)
			} else if m, ok := e.(map[string]interface{}); ok {
				if target, ok := m["target"].(string); ok {
					errorTargets = append(errorTargets, target)
				}
			}
		}
	case []string:
		errorTargets = v
	case string:
		errorTargets = []string{v}
	}

	for _, target := range errorTargets {
		errorShape, ok := r.model.Shapes[target]
		if !ok {
			continue
		}

		// Extract error information
		errorFault := getStringTrait(errorShape.Traits, "smithy.api#error")
		httpStatusCode := getIntTrait(errorShape.Traits, "smithy.api#httpError")
		if httpStatusCode == 0 {
			httpStatusCode = getIntTrait(errorShape.Traits, "smithy.api#httpCode")
		}
		errorCodeOverride := getStringTrait(errorShape.Traits, "aws.api#errorCode")

		errors = append(errors, proto_generator.ErrorInfo{
			ID:                int64(hashKey(target)),
			ShapeID:           target,
			Name:              extractName(target),
			ErrorFault:        errorFault,
			HTTPStatusCode:    httpStatusCode,
			ErrorCodeOverride: errorCodeOverride,
		})
	}

	return errors
}

// Helper functions

func hashKey(s string) uint32 {
	h := uint32(0)
	for _, c := range s {
		h = h*31 + uint32(c)
	}
	return h
}

func extractName(key string) string {
	if idx := strings.Index(key, "#"); idx >= 0 {
		return key[idx+1:]
	}
	return key
}

func getDocumentation(traits map[string]interface{}) string {
	if traits == nil {
		return ""
	}
	if doc, ok := traits["smithy.api#documentation"].(string); ok {
		return doc
	}
	return ""
}

func getStringTrait(traits map[string]interface{}, traitName string) string {
	if traits == nil {
		return ""
	}
	if val, ok := traits[traitName]; ok {
		if str, ok := val.(string); ok {
			return str
		}
		if m, ok := val.(map[string]interface{}); ok {
			if value, ok := m["value"].(string); ok {
				return value
			}
		}
	}
	return ""
}

func getIntTrait(traits map[string]interface{}, traitName string) int {
	if traits == nil {
		return 0
	}
	if val, ok := traits[traitName]; ok {
		if f, ok := val.(float64); ok {
			return int(f)
		}
		if i, ok := val.(int); ok {
			return i
		}
		if m, ok := val.(map[string]interface{}); ok {
			if value, ok := m["value"].(float64); ok {
				return int(value)
			}
		}
	}
	return 0
}

func hasTrait(traits map[string]interface{}, traitName string) bool {
	if traits == nil {
		return false
	}
	if val, ok := traits[traitName]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
		return true
	}
	return false
}

func parseURI(uri string) (path, query string) {
	if uri == "" {
		return "", ""
	}

	// Split on ? for query parameters
	parts := strings.SplitN(uri, "?", 2)
	path = parts[0]
	if len(parts) > 1 {
		query = parts[1]
	}

	return path, query
}

func toProtoEnumName(name string) string {
	result := make([]byte, 0, len(name))
	for i, c := range name {
		switch {
		case c >= 'a' && c <= 'z':
			result = append(result, byte(c))
		case c >= 'A' && c <= 'Z':
			result = append(result, byte(c))
		case c >= '0' && c <= '9':
			if i == 0 {
				result = append(result, '_', byte(c))
			} else {
				result = append(result, byte(c))
			}
		default:
			if i > 0 && (len(result) == 0 || result[len(result)-1] != '_') {
				result = append(result, '_')
			}
		}
	}
	if len(result) == 0 {
		return "UNKNOWN"
	}
	name = string(result)
	if name[0] >= '0' && name[0] <= '9' {
		name = "_" + name
	}
	return strings.ToUpper(name)
}
