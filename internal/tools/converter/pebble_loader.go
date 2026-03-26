// Package converter provides types and utilities for parsing Smithy model files.
package converter

import (
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/store/api"
)

// PebbleLoader loads Smithy model data into the runtime storage system.
// It converts parsed Smithy models into the internal API representation.
type PebbleLoader struct {
	serviceStore   *api.ServiceStore
	operationStore *api.OperationStore
	shapeStore     *api.ShapeStore
	memberStore    *api.MemberStore
	// New stores for traits, errors, and enum values
	shapeTraitStore     *api.ShapeTraitStore
	memberTraitStore    *api.MemberTraitStore
	operationErrorStore *api.OperationErrorStore
	enumValueStore      *api.EnumValueStore

	shapeIDCounter  int64
	memberIDCounter int64
	traitIDCounter  int64
	errorIDCounter  int64
	enumIDCounter   int64

	// shapeIDMap maps Smithy shape IDs to internal IDs
	shapeIDMap map[string]int64
}

// NewPebbleLoader creates a new PebbleLoader with the provided stores.
func NewPebbleLoader(
	serviceStore *api.ServiceStore,
	operationStore *api.OperationStore,
	shapeStore *api.ShapeStore,
	memberStore *api.MemberStore,
	shapeTraitStore *api.ShapeTraitStore,
	memberTraitStore *api.MemberTraitStore,
	operationErrorStore *api.OperationErrorStore,
	enumValueStore *api.EnumValueStore,
) *PebbleLoader {
	return &PebbleLoader{
		serviceStore:        serviceStore,
		operationStore:      operationStore,
		shapeStore:          shapeStore,
		memberStore:         memberStore,
		shapeTraitStore:     shapeTraitStore,
		memberTraitStore:    memberTraitStore,
		operationErrorStore: operationErrorStore,
		enumValueStore:      enumValueStore,
		shapeIDMap:          make(map[string]int64),
	}
}

// LoadModel loads a parsed Smithy model into the runtime storage.
// It processes the service, operations, and shapes, storing them in the provided stores.
func (l *PebbleLoader) LoadModel(parser *Parser) error {
	model := parser.GetModel()
	if model == nil {
		return fmt.Errorf("model not parsed")
	}

	serviceShape, serviceShapeID, err := parser.GetServiceShape()
	if err != nil {
		return fmt.Errorf("failed to get service shape: %w", err)
	}

	serviceName := parser.GetServiceName()
	service := l.buildService(serviceShape, serviceShapeID, serviceName, parser)
	if err := l.serviceStore.Put(service); err != nil {
		return fmt.Errorf("failed to store service: %w", err)
	}

	// First pass: process all operations (builds shapeIDMap)
	for shapeID, shape := range model.Shapes {
		if shape.Type == "operation" {
			if err := l.processOperation(shapeID, shape, service, parser); err != nil {
				return fmt.Errorf("failed to process operation %s: %w", shapeID, err)
			}
		}
	}

	// Second pass: process traits, errors, and enum values
	for shapeID, shape := range model.Shapes {
		// Process shape traits (excluding operation shapes which are handled separately)
		if shape.Type != "operation" && shape.Type != "service" {
			internalID, exists := l.shapeIDMap[shapeID]
			if exists {
				if err := l.processShapeTraits(shape, internalID); err != nil {
					return fmt.Errorf("failed to process shape traits for %s: %w", shapeID, err)
				}
			}
		}

		// Process enum values for string/enum types
		if shape.Type == "string" || shape.Type == "enum" {
			internalID, exists := l.shapeIDMap[shapeID]
			if exists {
				if err := l.processEnumValues(shape, internalID); err != nil {
					return fmt.Errorf("failed to process enum values for %s: %w", shapeID, err)
				}
			}
		}
	}

	return nil
}

func (l *PebbleLoader) buildService(shape *Shape, shapeID, serviceName string, parser *Parser) *api.Service {
	now := time.Now().Unix()
	sdkID, arnNamespace, cloudFormationName, cloudTrailEventSource, endpointPrefix := extractServiceInfo(shape.Traits)
	return &api.Service{
		ID:                    l.nextID(),
		ShapeID:               shapeID,
		Name:                  serviceName,
		Version:               extractStringTrait(shape.Traits, "smithy.api#version"),
		SDKID:                 sdkID,
		ARNNamespace:          arnNamespace,
		CloudFormationName:    cloudFormationName,
		CloudTrailEventSource: cloudTrailEventSource,
		EndpointPrefix:        endpointPrefix,
		Protocol:              string(parser.GetProtocol()),
		CreatedAt:             now,
	}
}

func (l *PebbleLoader) processOperation(shapeID string, shape *Shape, service *api.Service, parser *Parser) error {
	now := time.Now().Unix()

	httpBinding := parser.GetHTTPBinding(shape.Traits)
	var httpMethod, httpURI *string
	var httpCode *int
	if httpBinding != nil {
		httpMethod = &httpBinding.Method
		httpURI = &httpBinding.URI
		httpCode = &httpBinding.Code
	}

	httpPath, httpQuery := parser.ParseURI(derefString(httpURI, ""))

	var inputShapeID, outputShapeID *int64
	if shape.Input != nil {
		id := l.nextID()
		inputShapeID = &id
		if err := l.processShapeRecursive(shape.Input.Target, parser); err != nil {
			return err
		}
	}
	if shape.Output != nil {
		id := l.nextID()
		outputShapeID = &id
		if err := l.processShapeRecursive(shape.Output.Target, parser); err != nil {
			return err
		}
	}

	contentType := parser.GetContentTypeForOperation(shape.Traits)
	xAmzTarget := parser.GetXAmzTarget(shapeID)

	op := &api.Operation{
		ID:             l.nextID(),
		ShapeID:        shapeID,
		ServiceID:      service.ID,
		Name:           l.extractOperationName(shapeID),
		Documentation:  stringPtrOrNil(parser.GetDocumentation(shape.Traits)),
		HTTPMethod:     httpMethod,
		HTTPURI:        httpURI,
		HTTPCode:       httpCode,
		IsReadonly:     isReadonlyOperation(httpBinding),
		InputShapeID:   inputShapeID,
		OutputShapeID:  outputShapeID,
		ContentType:    stringPtrOrNil(contentType),
		XAmzTarget:     stringPtrOrNil(xAmzTarget),
		SmithyProtocol: stringPtrOrNil(string(parser.GetProtocol())),
		XAmznQueryMode: parser.GetXAmznQueryMode(parser.GetModel().Shapes[service.ShapeID]),
		HTTPURIPath:    stringPtrOrNil(httpPath),
		HTTPURIQuery:   stringPtrOrNil(httpQuery),
		CreatedAt:      now,
	}

	if err := l.operationStore.Put(service.Name, op); err != nil {
		return err
	}

	// Process operation errors
	if err := l.processOperationErrors(shape, op.ID, parser); err != nil {
		return fmt.Errorf("failed to process operation errors: %w", err)
	}

	return nil
}

func (l *PebbleLoader) processShapeRecursive(shapeID string, parser *Parser) error {
	model := parser.GetModel()
	shape, exists := model.Shapes[shapeID]
	if !exists {
		return nil
	}

	// Check if already processed
	if _, alreadyProcessed := l.shapeIDMap[shapeID]; alreadyProcessed {
		return nil
	}

	now := time.Now().Unix()
	internalID := l.nextID()
	l.shapeIDMap[shapeID] = internalID

	apiShape := &api.Shape{
		ID:            internalID,
		ShapeID:       shapeID,
		Type:          shape.Type,
		Documentation: stringPtrOrNil(parser.GetDocumentation(shape.Traits)),
		CreatedAt:     now,
	}

	switch shape.Type {
	case "list":
		if shape.Member != nil {
			id := l.nextID()
			apiShape.ListMemberShapeID = &id
			if err := l.processShapeRecursive(shape.Member.Target, parser); err != nil {
				return err
			}
		}
	case "map":
		if shape.Key != nil {
			id := l.nextID()
			apiShape.MapKeyShapeID = &id
			if err := l.processShapeRecursive(shape.Key.Target, parser); err != nil {
				return err
			}
		}
		if shape.Value != nil {
			id := l.nextID()
			apiShape.MapValueShapeID = &id
			if err := l.processShapeRecursive(shape.Value.Target, parser); err != nil {
				return err
			}
		}
	case "structure":
		for memberName, member := range shape.Members {
			m := &api.Member{
				ID:                      l.nextMemberID(),
				ShapeID:                 apiShape.ID,
				Name:                    memberName,
				TargetShapeID:           l.nextID(),
				Documentation:           stringPtrOrNil(parser.GetDocumentation(member.Traits)),
				IsRequired:              extractBoolTrait(member.Traits, "smithy.api#required"),
				HTTPLabel:               extractBoolTrait(member.Traits, "smithy.api#httpLabel"),
				HTTPQuery:               stringPtrOrNil(extractStringTrait(member.Traits, "smithy.api#httpQuery")),
				HTTPHeader:              stringPtrOrNil(extractStringTrait(member.Traits, "smithy.api#httpHeader")),
				HTTPPayload:             extractBoolTrait(member.Traits, "smithy.api#httpPayload"),
				JSONName:                stringPtrOrNil(extractStringTrait(member.Traits, "smithy.api#jsonName")),
				HTTPPrefix:              stringPtrOrNil(extractStringTrait(member.Traits, "smithy.api#httpPrefix")),
				HTTPLocationName:        stringPtrOrNil(parser.GetHTTPLocationName(member.Traits)),
				DefaultPayloadMediaType: stringPtrOrNil(extractStringTrait(member.Traits, "smithy.http#defaultPayloadMediaType")),
				CreatedAt:               now,
			}
			if err := l.memberStore.Put(m); err != nil {
				return fmt.Errorf("failed to store member %s: %w", memberName, err)
			}
			// Process member traits
			if err := l.processMemberTraits(member, m.ID); err != nil {
				return fmt.Errorf("failed to process member traits for %s: %w", memberName, err)
			}
			if err := l.processShapeRecursive(member.Target, parser); err != nil {
				return err
			}
		}
	}

	return l.shapeStore.Put(apiShape)
}

func (l *PebbleLoader) extractOperationName(shapeID string) string {
	parts := strings.Split(shapeID, "#")
	if len(parts) > 1 {
		return parts[1]
	}
	return shapeID
}

func (l *PebbleLoader) nextID() int64 {
	l.shapeIDCounter++
	return l.shapeIDCounter
}

func (l *PebbleLoader) nextMemberID() int64 {
	l.memberIDCounter++
	return l.memberIDCounter
}

func stringPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func derefString(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}

func isReadonlyOperation(httpBinding *HTTPBinding) bool {
	if httpBinding == nil {
		return false
	}
	method := strings.ToUpper(httpBinding.Method)
	return method == "GET" || method == "HEAD" || method == "OPTIONS"
}

// processShapeTraits stores shape traits that are not skipped.
func (l *PebbleLoader) processShapeTraits(shape *Shape, internalID int64) error {
	if shape.Traits == nil {
		return nil
	}

	now := time.Now().Unix()
	for traitName, traitValue := range shape.Traits {
		if !ShouldStoreShapeTrait(traitName) {
			continue
		}

		trait := &api.ShapeTrait{
			ID:         l.nextTraitID(),
			ShapeID:    internalID,
			TraitName:  traitName,
			TraitValue: formatTraitValue(traitValue),
			CreatedAt:  now,
		}
		if err := l.shapeTraitStore.Put(trait); err != nil {
			return fmt.Errorf("failed to store shape trait %s: %w", traitName, err)
		}
	}

	return nil
}

// processMemberTraits stores member traits that are not skipped.
func (l *PebbleLoader) processMemberTraits(member *Member, internalID int64) error {
	if member.Traits == nil {
		return nil
	}

	now := time.Now().Unix()
	for traitName, traitValue := range member.Traits {
		if !ShouldStoreMemberTrait(traitName) {
			continue
		}

		trait := &api.MemberTrait{
			ID:         l.nextMemberTraitID(),
			MemberID:   internalID,
			TraitName:  traitName,
			TraitValue: formatTraitValue(traitValue),
			CreatedAt:  now,
		}
		if err := l.memberTraitStore.Put(trait); err != nil {
			return fmt.Errorf("failed to store member trait %s: %w", traitName, err)
		}
	}

	return nil
}

// processOperationErrors stores operation error definitions.
func (l *PebbleLoader) processOperationErrors(shape *Shape, operationID int64, parser *Parser) error {
	if shape.Traits == nil {
		return nil
	}

	errorsTrait := shape.Traits["smithy.api#errors"]
	if errorsTrait == nil {
		// Try alternate key
		errorsTrait = shape.Traits["errors"]
	}
	if errorsTrait == nil {
		return nil
	}

	targets := ParseErrorsTrait(errorsTrait)
	model := parser.GetModel()

	now := time.Now().Unix()
	for _, target := range targets {
		errorShape, ok := model.Shapes[target]
		if !ok {
			continue
		}

		errorShapeID, exists := l.shapeIDMap[target]
		if !exists {
			// Process the error shape if not yet processed
			if err := l.processShapeRecursive(target, parser); err != nil {
				continue
			}
			errorShapeID = l.shapeIDMap[target]
		}

		errorFault := extractErrorFault(errorShape.Traits)
		if errorFault == "" {
			errorFault = parser.GetErrorFault(errorShape.Traits)
		}

		httpStatusCode := extractHTTPStatusCode(errorShape.Traits)
		if httpStatusCode == 0 {
			httpStatusCode = parser.GetHTTPErrorCode(errorShape.Traits)
		}

		errorCodeOverride := extractErrorCodeOverride(errorShape.Traits)
		if errorCodeOverride == "" {
			errorCodeOverride = parser.GetErrorCodeOverride(errorShape.Traits)
		}

		opError := &api.OperationError{
			ID:                l.nextErrorID(),
			OperationID:       operationID,
			ErrorShapeID:      errorShapeID,
			ErrorFault:        errorFault,
			HTTPStatusCode:    httpStatusCode,
			ErrorCodeOverride: errorCodeOverride,
			CreatedAt:         now,
		}
		if err := l.operationErrorStore.Put(opError); err != nil {
			return fmt.Errorf("failed to store operation error for %s: %w", target, err)
		}
	}

	return nil
}

// processEnumValues stores enum value definitions for string/enum shapes.
func (l *PebbleLoader) processEnumValues(shape *Shape, internalID int64) error {
	if shape.Members == nil && len(shape.Enum) == 0 {
		return nil
	}

	now := time.Now().Unix()

	// Handle enum members (newer Smithy format)
	for memberName, member := range shape.Members {
		enumVal := extractEnumValue(member.Traits)
		if enumVal == "" {
			enumVal = memberName // Use name as value if no enumValue trait
		}

		ev := &api.ShapeEnumValue{
			ID:            l.nextEnumID(),
			ShapeID:       internalID,
			Name:          memberName,
			Value:         enumVal,
			Documentation: extractDocumentation(member.Traits),
			CreatedAt:     now,
		}
		if err := l.enumValueStore.Put(ev); err != nil {
			return fmt.Errorf("failed to store enum value %s: %w", memberName, err)
		}
	}

	// Handle legacy enum format
	for _, enumVal := range shape.Enum {
		name := enumVal.Name
		if name == "" {
			name = enumVal.Value
		}

		ev := &api.ShapeEnumValue{
			ID:            l.nextEnumID(),
			ShapeID:       internalID,
			Name:          name,
			Value:         enumVal.Value,
			Documentation: extractDocumentation(enumVal.Traits),
			CreatedAt:     now,
		}
		if err := l.enumValueStore.Put(ev); err != nil {
			return fmt.Errorf("failed to store enum value %s: %w", name, err)
		}
	}

	return nil
}

func (l *PebbleLoader) nextTraitID() int64 {
	l.traitIDCounter++
	return l.traitIDCounter
}

func (l *PebbleLoader) nextMemberTraitID() int64 {
	l.memberIDCounter++
	return l.memberIDCounter
}

func (l *PebbleLoader) nextErrorID() int64 {
	l.errorIDCounter++
	return l.errorIDCounter
}

func (l *PebbleLoader) nextEnumID() int64 {
	l.enumIDCounter++
	return l.enumIDCounter
}
