// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func (s *DynamoDBService) buildTableDescription(table *dbstore.Table) map[string]interface{} {
	desc := map[string]interface{}{
		"TableName":                 table.Name,
		"TableArn":                  table.ARN,
		"TableStatus":               string(table.Status),
		"CreationDateTime":          table.CreationDateTime.Unix(),
		"KeySchema":                 buildKeySchemaResponse(table.KeySchema),
		"AttributeDefinitions":      buildAttributeDefinitionsResponse(table.AttributeDefinitions),
		"TableSizeBytes":            table.TableSizeBytes,
		"ItemCount":                 table.ItemCount,
		"DeletionProtectionEnabled": table.DeletionProtectionEnabled,
		"BillingModeSummary": map[string]interface{}{
			"BillingMode": string(table.BillingMode),
		},
	}

	if table.ProvisionedThroughput != nil {
		pt := map[string]interface{}{
			"ReadCapacityUnits":  table.ProvisionedThroughput.ReadCapacityUnits,
			"WriteCapacityUnits": table.ProvisionedThroughput.WriteCapacityUnits,
		}
		if !table.ProvisionedThroughput.LastDecreaseDateTime.IsZero() {
			pt["LastDecreaseDateTime"] = table.ProvisionedThroughput.LastDecreaseDateTime.Unix()
		}
		if !table.ProvisionedThroughput.LastIncreaseDateTime.IsZero() {
			pt["LastIncreaseDateTime"] = table.ProvisionedThroughput.LastIncreaseDateTime.Unix()
		}
		pt["NumberOfDecreasesToday"] = table.ProvisionedThroughput.NumberOfDecreasesToday
		desc["ProvisionedThroughput"] = pt
	} else {
		desc["ProvisionedThroughput"] = map[string]interface{}{
			"ReadCapacityUnits":      0,
			"WriteCapacityUnits":     0,
			"NumberOfDecreasesToday": 0,
		}
	}

	if len(table.GlobalSecondaryIndexes) > 0 {
		desc["GlobalSecondaryIndexes"] = buildGSIsResponse(table.GlobalSecondaryIndexes)
	}

	if len(table.LocalSecondaryIndexes) > 0 {
		desc["LocalSecondaryIndexes"] = buildLSIsResponse(table.LocalSecondaryIndexes)
	}

	if table.StreamSpecification != nil {
		desc["StreamSpecification"] = map[string]interface{}{
			"StreamEnabled":  table.StreamSpecification.StreamEnabled,
			"StreamViewType": string(table.StreamSpecification.StreamViewType),
		}
		if table.StreamSpecification.StreamEnabled && table.StreamArn != "" {
			desc["LatestStreamArn"] = table.StreamArn
			desc["LatestStreamLabel"] = table.LatestStreamLabel
		}
	}

	if table.SSEDescription != nil {
		desc["SSEDescription"] = map[string]interface{}{
			"Status":          table.SSEDescription.Status,
			"SSEType":         string(table.SSEDescription.SSEType),
			"KMSMasterKeyArn": table.SSEDescription.KMSMasterKeyArn,
		}
	}

	desc["WarmThroughput"] = map[string]interface{}{
		"ReadUnitsPerSecond":  0,
		"WriteUnitsPerSecond": 0,
		"Status":              "ACTIVE",
	}

	tableClass := table.TableClass
	if tableClass == "" {
		tableClass = "STANDARD"
	}
	desc["TableClassSummary"] = map[string]interface{}{
		"TableClass": tableClass,
	}

	desc["Replicas"] = []interface{}{}

	return desc
}

func buildKeySchemaResponse(schema []*dbstore.KeySchemaElement) []map[string]interface{} {
	result := make([]map[string]interface{}, len(schema))
	for i, s := range schema {
		result[i] = map[string]interface{}{
			"AttributeName": s.AttributeName,
			"KeyType":       string(s.KeyType),
		}
	}
	return result
}

func buildAttributeDefinitionsResponse(defs []*dbstore.AttributeDefinition) []map[string]interface{} {
	result := make([]map[string]interface{}, len(defs))
	for i, d := range defs {
		result[i] = map[string]interface{}{
			"AttributeName": d.AttributeName,
			"AttributeType": string(d.AttributeType),
		}
	}
	return result
}

func buildGSIsResponse(gsis []*dbstore.GlobalSecondaryIndex) []map[string]interface{} {
	result := make([]map[string]interface{}, len(gsis))
	for i, g := range gsis {
		idx := map[string]interface{}{
			"IndexName":      g.IndexName,
			"IndexArn":       g.IndexArn,
			"KeySchema":      buildKeySchemaResponse(g.KeySchema),
			"Projection":     buildProjectionResponse(g.Projection),
			"IndexStatus":    string(g.IndexStatus),
			"ItemCount":      g.ItemCount,
			"IndexSizeBytes": g.IndexSizeBytes,
		}
		if g.ProvisionedThroughput != nil {
			idx["ProvisionedThroughput"] = map[string]interface{}{
				"ReadCapacityUnits":  g.ProvisionedThroughput.ReadCapacityUnits,
				"WriteCapacityUnits": g.ProvisionedThroughput.WriteCapacityUnits,
			}
		} else {
			idx["ProvisionedThroughput"] = map[string]interface{}{
				"ReadCapacityUnits":  0,
				"WriteCapacityUnits": 0,
			}
		}
		result[i] = idx
	}
	return result
}

func buildLSIsResponse(lsis []*dbstore.LocalSecondaryIndex) []map[string]interface{} {
	result := make([]map[string]interface{}, len(lsis))
	for i, l := range lsis {
		idx := map[string]interface{}{
			"IndexName":  l.IndexName,
			"KeySchema":  buildKeySchemaResponse(l.KeySchema),
			"Projection": buildProjectionResponse(l.Projection),
			"ItemCount":  l.ItemCount,
		}
		if l.IndexSizeBytes > 0 {
			idx["IndexSizeBytes"] = l.IndexSizeBytes
		}
		result[i] = idx
	}
	return result
}

func buildProjectionResponse(p *dbstore.Projection) map[string]interface{} {
	if p == nil {
		return map[string]interface{}{"ProjectionType": "ALL"}
	}
	resp := map[string]interface{}{
		"ProjectionType": p.ProjectionType,
	}
	if len(p.NonKeyAttributes) > 0 {
		resp["NonKeyAttributes"] = p.NonKeyAttributes
	}
	return resp
}
