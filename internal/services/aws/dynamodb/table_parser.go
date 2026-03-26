// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"regexp"

	"vorpalstacks/internal/services/aws/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

var tableNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

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

func parseGlobalSecondaryIndexes(params map[string]interface{}) []*dbstore.GlobalSecondaryIndex {
	gsis, ok := params["GlobalSecondaryIndexes"].([]interface{})
	if !ok {
		return nil
	}

	var result []*dbstore.GlobalSecondaryIndex
	for _, g := range gsis {
		if gm, ok := g.(map[string]interface{}); ok {
			idx := &dbstore.GlobalSecondaryIndex{
				IndexName:             request.GetStringParam(gm, "IndexName"),
				KeySchema:             parseKeySchema(gm),
				Projection:            parseProjection(gm),
				ProvisionedThroughput: parseProvisionedThroughput(gm),
				IndexStatus:           dbstore.IndexStatusActive,
			}
			result = append(result, idx)
		}
	}
	return result
}

func parseLocalSecondaryIndexes(params map[string]interface{}) []*dbstore.LocalSecondaryIndex {
	lsis, ok := params["LocalSecondaryIndexes"].([]interface{})
	if !ok {
		return nil
	}

	var result []*dbstore.LocalSecondaryIndex
	for _, l := range lsis {
		if lm, ok := l.(map[string]interface{}); ok {
			idx := &dbstore.LocalSecondaryIndex{
				IndexName:  request.GetStringParam(lm, "IndexName"),
				KeySchema:  parseKeySchema(lm),
				Projection: parseProjection(lm),
			}
			result = append(result, idx)
		}
	}
	return result
}

func parseProjection(params map[string]interface{}) *dbstore.Projection {
	proj, ok := params["Projection"].(map[string]interface{})
	if !ok {
		return &dbstore.Projection{ProjectionType: "ALL"}
	}

	p := &dbstore.Projection{
		ProjectionType: request.GetStringParam(proj, "ProjectionType"),
	}
	if p.ProjectionType == "" {
		p.ProjectionType = "ALL"
	}

	if nkAs, ok := proj["NonKeyAttributes"].([]interface{}); ok {
		for _, nk := range nkAs {
			if nks, ok := nk.(string); ok {
				p.NonKeyAttributes = append(p.NonKeyAttributes, nks)
			}
		}
	}
	return p
}

func parseStreamSpecification(params map[string]interface{}) *dbstore.StreamSpecification {
	ss, ok := params["StreamSpecification"].(map[string]interface{})
	if !ok {
		return nil
	}

	enabled := false
	if e, ok := ss["StreamEnabled"].(bool); ok {
		enabled = e
	}

	return &dbstore.StreamSpecification{
		StreamEnabled:  enabled,
		StreamViewType: dbstore.StreamViewType(request.GetStringParam(ss, "StreamViewType")),
	}
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
