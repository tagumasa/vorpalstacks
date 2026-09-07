// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func parseKeySchema(params map[string]interface{}) []*dbstore.KeySchemaElement {
	schema, ok := params["KeySchema"].([]interface{})
	if !ok {
		return nil
	}

	var result []*dbstore.KeySchemaElement
	for _, s := range schema {
		if sm, ok := s.(map[string]interface{}); ok {
			elem := &dbstore.KeySchemaElement{
				AttributeName: request.GetStringParam(sm, "AttributeName"),
				KeyType:       dbstore.KeyType(request.GetStringParam(sm, "KeyType")),
			}
			result = append(result, elem)
		}
	}
	return result
}

func parseAttributeDefinitions(params map[string]interface{}) []*dbstore.AttributeDefinition {
	defs, ok := params["AttributeDefinitions"].([]interface{})
	if !ok {
		return nil
	}

	var result []*dbstore.AttributeDefinition
	for _, d := range defs {
		if dm, ok := d.(map[string]interface{}); ok {
			def := &dbstore.AttributeDefinition{
				AttributeName: request.GetStringParam(dm, "AttributeName"),
				AttributeType: dbstore.ScalarAttributeType(request.GetStringParam(dm, "AttributeType")),
			}
			result = append(result, def)
		}
	}
	return result
}

// parseProvisionedThroughput decodes the ProvisionedThroughput member
// verbatim, returning nil only when the member is absent. Value bounds are
// validated by the Core (validateProvisionedThroughputValues), so a present
// member with out-of-range units still reaches the Core's rejection instead
// of masquerading as an omitted one.
func parseProvisionedThroughput(params map[string]interface{}) *dbstore.ProvisionedThroughput {
	pt, ok := params["ProvisionedThroughput"].(map[string]interface{})
	if !ok {
		return nil
	}

	return &dbstore.ProvisionedThroughput{
		ReadCapacityUnits:  request.GetInt64Param(pt, "ReadCapacityUnits"),
		WriteCapacityUnits: request.GetInt64Param(pt, "WriteCapacityUnits"),
	}
}

func parseGlobalSecondaryIndexes(params map[string]interface{}) ([]*dbstore.GlobalSecondaryIndex, error) {
	gsis, ok := params["GlobalSecondaryIndexes"].([]interface{})
	if !ok {
		return nil, nil
	}

	var result []*dbstore.GlobalSecondaryIndex
	for _, g := range gsis {
		if gm, ok := g.(map[string]interface{}); ok {
			idxName := request.GetStringParam(gm, "IndexName")
			if idxName == "" {
				return nil, ErrInvalidParameter
			}
			if !validateResourceName(idxName) {
				return nil, ErrInvalidParameter
			}
			proj, err := parseProjection(gm)
			if err != nil {
				return nil, err
			}
			if !validateProjectionRequired(gm["Projection"].(map[string]interface{})) {
				return nil, ErrInvalidParameter
			}
			idx := &dbstore.GlobalSecondaryIndex{
				IndexName:             idxName,
				KeySchema:             parseKeySchema(gm),
				Projection:            proj,
				ProvisionedThroughput: parseProvisionedThroughput(gm),
				IndexStatus:           dbstore.IndexStatusActive,
			}
			result = append(result, idx)
		}
	}
	return result, nil
}

func parseLocalSecondaryIndexes(params map[string]interface{}) ([]*dbstore.LocalSecondaryIndex, error) {
	lsis, ok := params["LocalSecondaryIndexes"].([]interface{})
	if !ok {
		return nil, nil
	}

	var result []*dbstore.LocalSecondaryIndex
	for _, l := range lsis {
		if lm, ok := l.(map[string]interface{}); ok {
			idxName := request.GetStringParam(lm, "IndexName")
			if idxName == "" {
				return nil, ErrInvalidParameter
			}
			if !validateResourceName(idxName) {
				return nil, ErrInvalidParameter
			}
			proj, err := parseProjection(lm)
			if err != nil {
				return nil, err
			}
			if !validateProjectionRequired(lm["Projection"].(map[string]interface{})) {
				return nil, ErrInvalidParameter
			}
			idx := &dbstore.LocalSecondaryIndex{
				IndexName:  idxName,
				KeySchema:  parseKeySchema(lm),
				Projection: proj,
			}
			result = append(result, idx)
		}
	}
	return result, nil
}

// parseVectorIndexes parses the CreateTable VectorIndexes member. Absence of
// the member yields nil (no vector indexes); field-level shape checks mirror
// the GSI/LSI parsers, while cross-index rules (name uniqueness across index
// families, shared-attribute dimension consistency) are enforced by the core.
func parseVectorIndexes(params map[string]interface{}) ([]*dbstore.VectorIndex, error) {
	rawList, ok := params["VectorIndexes"].([]interface{})
	if !ok {
		return nil, nil
	}
	result := make([]*dbstore.VectorIndex, 0, len(rawList))
	for _, v := range rawList {
		vm, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}
		idx, err := parseVectorIndex(vm)
		if err != nil {
			return nil, err
		}
		result = append(result, idx)
	}
	return result, nil
}

// parseVectorIndex parses one VectorIndex definition — a CreateTable
// VectorIndexes element or an UpdateTable VectorIndexUpdates Create action,
// which carry identical members.
func parseVectorIndex(vm map[string]interface{}) (*dbstore.VectorIndex, error) {
	idxName := request.GetStringParam(vm, "IndexName")
	if idxName == "" || !validateResourceName(idxName) {
		return nil, ErrInvalidParameter
	}
	vecAttr, ok := vm["VectorAttribute"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	attrName := request.GetStringParam(vecAttr, "AttributeName")
	if attrName == "" || len(attrName) > dbstore.VectorAttributeNameMax {
		return nil, ErrInvalidParameter
	}
	dims := request.GetInt64Param(vm, "Dimensions")
	if dims < dbstore.VectorDimensionsMin || dims > dbstore.VectorDimensionsMax {
		return nil, ErrInvalidParameter
	}
	distanceFn := request.GetStringParam(vm, "DistanceFunction")
	if !validateVectorDistanceFunction(distanceFn) {
		return nil, ErrInvalidParameter
	}
	proj, err := parseProjection(vm)
	if err != nil {
		return nil, err
	}
	if !validateProjectionRequired(vm["Projection"].(map[string]interface{})) {
		return nil, ErrInvalidParameter
	}
	searchSchema, err := parseSearchSchema(vm)
	if err != nil {
		return nil, err
	}
	return &dbstore.VectorIndex{
		IndexName:           idxName,
		VectorAttributeName: attrName,
		Dimensions:          dims,
		DistanceFunction:    distanceFn,
		Projection:          proj,
		SearchSchema:        searchSchema,
		IndexStatus:         dbstore.IndexStatusActive,
	}, nil
}

// parseSearchSchema parses the optional SearchSchema member: a non-empty
// list of {AttributeName, SearchSchemaElementType} elements with at most one
// HASH element and no duplicate attribute names.
func parseSearchSchema(vm map[string]interface{}) ([]*dbstore.SearchSchemaElement, error) {
	rawList, ok := vm["SearchSchema"].([]interface{})
	if !ok {
		return nil, nil
	}
	result := make([]*dbstore.SearchSchemaElement, 0, len(rawList))
	seen := make(map[string]bool)
	hasHash := false
	inlineFilters := 0
	for _, e := range rawList {
		em, ok := e.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}
		attrName := request.GetStringParam(em, "AttributeName")
		if attrName == "" {
			return nil, ErrInvalidParameter
		}
		elemType := request.GetStringParam(em, "SearchSchemaElementType")
		if elemType != "HASH" && elemType != "INLINE_FILTER" {
			return nil, ErrInvalidParameter
		}
		if seen[attrName] {
			return nil, ErrInvalidParameter
		}
		seen[attrName] = true
		if elemType == "HASH" {
			if hasHash {
				return nil, ErrInvalidParameter
			}
			hasHash = true
		} else {
			inlineFilters++
			if inlineFilters > dbstore.VectorInlineFiltersMax {
				return nil, ErrInvalidParameter
			}
		}
		result = append(result, &dbstore.SearchSchemaElement{
			AttributeName:           attrName,
			SearchSchemaElementType: elemType,
		})
	}
	if len(result) == 0 {
		return nil, ErrInvalidParameter
	}
	return result, nil
}

func parseProjection(params map[string]interface{}) (*dbstore.Projection, error) {
	proj, ok := params["Projection"].(map[string]interface{})
	if !ok {
		// Projection is required for GSI and LSI definitions.
		return nil, ErrInvalidParameter
	}

	p := &dbstore.Projection{
		ProjectionType: request.GetStringParam(proj, "ProjectionType"),
	}
	if p.ProjectionType == "" {
		p.ProjectionType = ProjectionTypeAll
	}
	// Validate the ProjectionType enum (ALL, KEYS_ONLY, INCLUDE).
	if !validateProjectionType(p.ProjectionType) {
		return nil, ErrInvalidParameter
	}

	if nkAs, ok := proj["NonKeyAttributes"].([]interface{}); ok {
		for _, nk := range nkAs {
			if nks, ok := nk.(string); ok {
				p.NonKeyAttributes = append(p.NonKeyAttributes, nks)
			}
		}
	}
	return p, nil
}

func parseStreamSpecification(params map[string]interface{}) (*dbstore.StreamSpecification, error) {
	ss, ok := params["StreamSpecification"].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	enabled, err := validateBoolParam(ss, "StreamEnabled", false)
	if err != nil {
		return nil, err
	}

	return &dbstore.StreamSpecification{
		StreamEnabled:  enabled,
		StreamViewType: dbstore.StreamViewType(request.GetStringParam(ss, "StreamViewType")),
	}, nil
}

func parseSSESpecification(ss map[string]interface{}) (*dbstore.SSEDescription, error) {
	enabled := false
	if e, ok := ss["Enabled"].(bool); ok {
		enabled = e
	}
	if !enabled {
		return nil, nil
	}

	sseType := dbstore.SSEType(request.GetStringParam(ss, "SSEType"))
	if sseType == "" {
		sseType = dbstore.SSETypeAES256
	}

	if sseType != dbstore.SSETypeAES256 && sseType != dbstore.SSETypeKMS {
		return nil, ErrInvalidParameter
	}

	kmsMasterKeyId := request.GetStringParam(ss, "KMSMasterKeyId")
	if sseType == dbstore.SSETypeKMS && kmsMasterKeyId == "" {
		return nil, ErrInvalidParameter
	}

	return &dbstore.SSEDescription{
		Status:          "ENABLED",
		SSEType:         sseType,
		KMSMasterKeyArn: kmsMasterKeyId,
	}, nil
}
