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

func parseProvisionedThroughput(params map[string]interface{}) *dbstore.ProvisionedThroughput {
	pt, ok := params["ProvisionedThroughput"].(map[string]interface{})
	if !ok {
		return nil
	}

	readUnits := request.GetInt64Param(pt, "ReadCapacityUnits")
	writeUnits := request.GetInt64Param(pt, "WriteCapacityUnits")

	if readUnits < 1 || writeUnits < 1 {
		return nil
	}

	return &dbstore.ProvisionedThroughput{
		ReadCapacityUnits:  readUnits,
		WriteCapacityUnits: writeUnits,
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
			if err := validateIndexName(idxName); err != nil {
				return nil, err
			}
			proj, err := parseProjection(gm)
			if err != nil {
				return nil, err
			}
			if err := validateProjectionRequired(gm["Projection"].(map[string]interface{})); err != nil {
				return nil, err
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
			if err := validateIndexName(idxName); err != nil {
				return nil, err
			}
			proj, err := parseProjection(lm)
			if err != nil {
				return nil, err
			}
			if err := validateProjectionRequired(lm["Projection"].(map[string]interface{})); err != nil {
				return nil, err
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
	if err := validateProjectionType(p.ProjectionType); err != nil {
		return nil, err
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
