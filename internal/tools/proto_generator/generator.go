// Package proto_generator provides protobuf code generation for vorpalstacks.
package proto_generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

// ProtoGenerator generates proto files from Smithy models
type ProtoGenerator struct {
	reader     ProtoDataReader
	typeMapper *TypeMapper
}

// NewProtoGenerator creates a new ProtoGenerator
func NewProtoGenerator(reader ProtoDataReader) *ProtoGenerator {
	return &ProtoGenerator{
		reader:     reader,
		typeMapper: NewTypeMapper(),
	}
}

func (g *ProtoGenerator) getServiceFilename(service *ServiceInfo) string {
	if service.Name != "" {
		return sanitizePackageName(strings.ToLower(strings.ReplaceAll(service.Name, " ", "")))
	}
	if service.EndpointPrefix != "" {
		return sanitizePackageName(service.EndpointPrefix)
	}
	return "unknown"
}

// GenerateProto generates a proto file
func (g *ProtoGenerator) GenerateProto(ctx context.Context, serviceName, outputDir string, tmpl *template.Template) error {
	service, err := g.reader.FindServiceByName(ctx, serviceName)
	if err != nil {
		return fmt.Errorf("service not found: %w", err)
	}
	if service == nil {
		return fmt.Errorf("service %q not found", serviceName)
	}

	operations, err := g.reader.FindOperationsByServiceID(ctx, service.ID)
	if err != nil {
		return fmt.Errorf("failed to get operations: %w", err)
	}

	shapes, err := g.reader.FindShapesByServiceID(ctx, service.ID)
	if err != nil {
		return fmt.Errorf("failed to get shapes: %w", err)
	}

	data, err := g.buildTemplateData(ctx, service, operations, shapes, tmpl == flutterGrpcTemplate)
	if err != nil {
		return fmt.Errorf("failed to build template data: %w", err)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil { // #nosec G301
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filename := g.getServiceFilename(service)
	outputPath := filepath.Join(outputDir, fmt.Sprintf("%s.proto", filename))

	if err := g.writeProtoFile(outputPath, data, tmpl); err != nil {
		return fmt.Errorf("failed to write proto file: %w", err)
	}

	fmt.Printf("Generated: %s\n", outputPath)
	return nil
}

// GenerateInternalProto generates internal proto
func (g *ProtoGenerator) GenerateInternalProto(ctx context.Context, serviceName, outputDir string) error {
	return g.GenerateProto(ctx, serviceName, outputDir, internalGrpcTemplate)
}

// GenerateFlutterProto generates Flutter proto
func (g *ProtoGenerator) GenerateFlutterProto(ctx context.Context, serviceName, outputDir string) error {
	return g.GenerateProto(ctx, serviceName, outputDir, flutterGrpcTemplate)
}

func (g *ProtoGenerator) buildTemplateData(ctx context.Context, service *ServiceInfo, operations []*OperationInfo, shapes []*ShapeInfo, forFlutter bool) (ProtoTemplateData, error) {
	if service == nil {
		return ProtoTemplateData{}, fmt.Errorf("service cannot be nil")
	}
	if service.Name == "" {
		return ProtoTemplateData{}, fmt.Errorf("service name cannot be empty")
	}
	packageName := sanitizePackageName(strings.ToLower(strings.ReplaceAll(service.Name, " ", "")))
	if packageName == "" {
		packageName = sanitizePackageName(service.EndpointPrefix)
	}
	serviceName := sanitizeServiceName(service.Name)

	goPackage := fmt.Sprintf("vorpalstacks/internal/pb/aws/%s", packageName)
	dartPackage := ""
	if forFlutter {
		dartPackage = fmt.Sprintf("vorpalstacks.proto.%s", packageName)
	}

	opData := make([]OperationData, len(operations))
	for i, op := range operations {
		inputShape := ""
		outputShape := ""
		if op.InputShapeID != nil {
			shape, err := g.reader.FindShapeByID(ctx, *op.InputShapeID)
			if err != nil {
				return ProtoTemplateData{}, fmt.Errorf("failed to find input shape for operation %s: %w", op.Name, err)
			}
			if shape != nil {
				inputShape = toProtoMessageName(shape.ShapeID)
			}
		}
		if op.OutputShapeID != nil {
			shape, err := g.reader.FindShapeByID(ctx, *op.OutputShapeID)
			if err != nil {
				return ProtoTemplateData{}, fmt.Errorf("failed to find output shape for operation %s: %w", op.Name, err)
			}
			if shape != nil {
				outputShape = toProtoMessageName(shape.ShapeID)
			}
		}

		// Convert errors
		errors := make([]ErrorData, len(op.Errors))
		for j, err := range op.Errors {
			errors[j] = ErrorData{
				Name:              err.Name,
				ErrorFault:        err.ErrorFault,
				HTTPStatusCode:    err.HTTPStatusCode,
				ErrorCodeOverride: err.ErrorCodeOverride,
			}
		}

		opData[i] = OperationData{
			Name:           op.Name,
			InputShape:     inputShape,
			OutputShape:    outputShape,
			Comment:        op.Documentation,
			HTTPMethod:     op.HTTPMethod,
			HTTPURIPath:    op.HTTPURIPath,
			HTTPURIQuery:   op.HTTPURIQuery,
			ContentType:    op.ContentType,
			XAmzTarget:     op.XAmzTarget,
			Protocol:       op.Protocol,
			SmithyProtocol: op.SmithyProtocol,
			XAmznQueryMode: op.XAmznQueryMode,
			Errors:         errors,
		}
	}

	shapeData := make([]ShapeData, 0)
	for _, shape := range shapes {
		if shape.Type == "structure" || shape.Type == "union" || shape.Type == "error" {
			fields, err := g.buildFieldsForShape(ctx, shape)
			if err != nil {
				return ProtoTemplateData{}, fmt.Errorf("failed to build fields for shape %s: %w", shape.ShapeID, err)
			}
			shapeData = append(shapeData, ShapeData{
				Name:    toProtoMessageName(shape.ShapeID),
				Fields:  fields,
				Comment: shape.Documentation,
			})
		} else if shape.Type == "enum" {
			enumValues, err := g.reader.FindEnumValuesByShapeID(ctx, shape.ID)
			if err != nil {
				return ProtoTemplateData{}, fmt.Errorf("failed to find enum values for shape %s: %w", shape.ShapeID, err)
			}
			enumName := toProtoMessageName(shape.ShapeID)
			prefixedValues := make([]*EnumValueInfo, len(enumValues))
			prefix := toEnumPrefix(enumName)
			for i, v := range enumValues {
				prefixedValues[i] = &EnumValueInfo{
					Name:  prefix + "_" + v.Name,
					Index: i,
				}
			}
			shapeData = append(shapeData, ShapeData{
				Name:       enumName,
				IsEnum:     true,
				EnumValues: prefixedValues,
			})
		}
	}

	return ProtoTemplateData{
		PackageName: packageName,
		ServiceName: fmt.Sprintf("%sService", serviceName),
		GoPackage:   goPackage,
		DartPackage: dartPackage,
		Operations:  opData,
		Shapes:      shapeData,
	}, nil
}

func (g *ProtoGenerator) buildFieldsForShape(ctx context.Context, shape *ShapeInfo) ([]FieldData, error) {
	members, err := g.reader.FindMembersByShapeID(ctx, shape.ID)
	if err != nil {
		return nil, err
	}

	fields := make([]FieldData, 0, len(members))
	for _, member := range members {
		var protoType string

		if member.TargetShapeID > 0 {
			targetShape, err := g.reader.FindShapeByID(ctx, member.TargetShapeID)
			if err != nil {
				return nil, fmt.Errorf("failed to find target shape for member %s: %w", member.Name, err)
			}
			if targetShape == nil {
				protoType = "string"
			} else {
				protoType = g.resolveProtoType(ctx, targetShape)
			}
		} else {
			protoType = "string"
		}

		fields = append(fields, FieldData{
			Name:                    member.Name,
			NameLower:               toProtoFieldName(member.Name),
			Type:                    protoType,
			Number:                  stableFieldNumber(member.Name),
			IsRequired:              member.IsRequired,
			HTTPLabel:               member.HTTPLabel,
			HTTPQuery:               member.HTTPQuery,
			HTTPHeader:              member.HTTPHeader,
			HTTPPayload:             member.HTTPPayload,
			JSONName:                member.JSONName,
			HTTPPrefix:              member.HTTPPrefix,
			HTTPLocationName:        member.HTTPLocationName,
			DefaultPayloadMediaType: member.DefaultPayloadMediaType,
		})
	}

	return fields, nil
}

func (g *ProtoGenerator) resolveProtoType(ctx context.Context, shape *ShapeInfo) string {
	if shape.ListMemberID != nil {
		memberShape, err := g.reader.FindShapeByID(ctx, *shape.ListMemberID)
		if err != nil || memberShape == nil {
			return "repeated string"
		}
		memberType := g.resolveProtoType(ctx, memberShape)
		if strings.HasPrefix(memberType, "repeated ") || strings.HasPrefix(memberType, "map<") {
			return memberType
		}
		return "repeated " + memberType
	}

	if shape.MapKeyShapeID != nil && shape.MapValueShapeID != nil {
		valueShape, err := g.reader.FindShapeByID(ctx, *shape.MapValueShapeID)
		if err != nil || valueShape == nil {
			return "map<string, string>"
		}
		valueType := g.resolveProtoType(ctx, valueShape)
		if strings.HasPrefix(valueType, "repeated ") || strings.HasPrefix(valueType, "map<") {
			valueType = "string"
		}
		return fmt.Sprintf("map<string, %s>", valueType)
	}

	switch shape.Type {
	case "structure", "union", "error":
		return toProtoMessageName(shape.ShapeID)
	case "enum":
		return toProtoMessageName(shape.ShapeID)
	default:
		return g.typeMapper.MapSmithyType(shape.Type)
	}
}

func toProtoFieldName(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	var result []rune
	prevIsLower := false
	for i, r := range runes {
		if i == 0 {
			result = append(result, unicode.ToLower(r))
			prevIsLower = true
		} else if unicode.IsUpper(r) {
			if prevIsLower {
				result = append(result, unicode.ToLower(r))
			} else {
				result = append(result, unicode.ToLower(r))
			}
			prevIsLower = false
		} else {
			result = append(result, r)
			prevIsLower = unicode.IsLower(r)
		}
	}
	return string(result)
}

func toProtoMessageName(shapeID string) string {
	if idx := strings.Index(shapeID, "#"); idx >= 0 {
		return shapeID[idx+1:]
	}
	return shapeID
}

func toEnumPrefix(name string) string {
	var result []rune
	for i, r := range name {
		if i == 0 {
			result = append(result, r)
			continue
		}
		if r >= 'A' && r <= 'Z' {
			if len(result) > 0 && result[len(result)-1] != '_' {
				result = append(result, '_')
			}
			result = append(result, r)
		} else {
			result = append(result, r)
		}
	}
	return strings.ToUpper(string(result))
}

func sanitizePackageName(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	return name
}

func sanitizeServiceName(name string) string {
	parts := strings.Fields(name)
	var result string
	for _, part := range parts {
		if len(part) > 0 {
			result += strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return result
}

func (g *ProtoGenerator) writeProtoFile(path string, data ProtoTemplateData, tmpl *template.Template) error {
	// #nosec G304
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

func stableFieldNumber(fieldName string) int {
	h := fnv32a(fieldName)
	num := int(h%536870912) + 1
	if num >= 19000 && num <= 19999 {
		num = num%19000 + 1
		if num < 1 {
			num = 1
		}
	}
	return num
}

func fnv32a(s string) uint32 {
	h := uint32(2166136261)
	for _, c := range s {
		h ^= uint32(c)
		h *= 16777619
	}
	return h
}
