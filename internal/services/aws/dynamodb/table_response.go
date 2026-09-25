// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// buildTableDescription renders the TableDescription response member.
// replicas carries the table's replication-group rendering; a standalone
// table passes nil for the empty list.
func (s *DynamoDBService) buildTableDescription(table *dbstore.Table, replicas []interface{}) map[string]interface{} {
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
	// The table UUID is emitted only once minted; records written before
	// the identity existed carry none and the optional member stays absent.
	if table.TableId != "" {
		desc["TableId"] = table.TableId
	}

	desc["ProvisionedThroughput"] = buildProvisionedThroughputDescriptionResponse(table.ProvisionedThroughput, true)

	if len(table.GlobalSecondaryIndexes) > 0 {
		desc["GlobalSecondaryIndexes"] = buildGSIsResponse(table.GlobalSecondaryIndexes)
	}

	if len(table.LocalSecondaryIndexes) > 0 {
		desc["LocalSecondaryIndexes"] = buildLSIsResponse(table.LocalSecondaryIndexes)
	}

	if len(table.VectorIndexes) > 0 {
		desc["VectorIndexes"] = buildVectorIndexesResponse(table.VectorIndexes)
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
			"Status":          string(table.SSEDescription.Status),
			"SSEType":         string(table.SSEDescription.SSEType),
			"KMSMasterKeyArn": table.SSEDescription.KMSMasterKeyArn,
		}
	}

	if odt := buildOnDemandThroughputResponse(table.OnDemandThroughput); odt != nil {
		desc["OnDemandThroughput"] = odt
	}

	if table.WarmThroughput != nil {
		// Warm throughput applies synchronously on this single-store
		// platform, so ACTIVE is the live state rather than a placeholder.
		desc["WarmThroughput"] = map[string]interface{}{
			"ReadUnitsPerSecond":  table.WarmThroughput.ReadUnitsPerSecond,
			"WriteUnitsPerSecond": table.WarmThroughput.WriteUnitsPerSecond,
			"Status":              "ACTIVE",
		}
	}

	if table.RestoreSummary != nil {
		restoreSummary := map[string]interface{}{
			"SourceTableArn":    table.RestoreSummary.SourceTableArn,
			"RestoreDateTime":   table.RestoreSummary.RestoreDateTime.Unix(),
			"RestoreInProgress": table.RestoreSummary.RestoreInProgress,
		}
		if table.RestoreSummary.SourceBackupArn != "" {
			restoreSummary["SourceBackupArn"] = table.RestoreSummary.SourceBackupArn
		}
		desc["RestoreSummary"] = restoreSummary
	}

	tableClass := string(table.TableClass)
	if tableClass == "" {
		tableClass = "STANDARD"
	}
	desc["TableClassSummary"] = map[string]interface{}{
		"TableClass": tableClass,
	}

	// Replicas is the global-table replication rendering; a standalone
	// table carries none and the member stays absent — the protocol omits
	// unset members rather than carrying an empty list.
	if replicas != nil {
		desc["Replicas"] = replicas
	}

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
		idx["ProvisionedThroughput"] = buildProvisionedThroughputDescriptionResponse(g.ProvisionedThroughput, false)
		if odt := buildOnDemandThroughputResponse(g.OnDemandThroughput); odt != nil {
			idx["OnDemandThroughput"] = odt
		}
		if g.WarmThroughput != nil {
			idx["WarmThroughput"] = map[string]interface{}{
				"ReadUnitsPerSecond":  g.WarmThroughput.ReadUnitsPerSecond,
				"WriteUnitsPerSecond": g.WarmThroughput.WriteUnitsPerSecond,
			}
		}
		result[i] = idx
	}
	return result
}

// buildOnDemandThroughputResponse renders the on-demand maximum pair
// presence-aware: a member the request never carried stays unset in the
// stored record (zero) and is rendered absent, never as zero — only a
// set limit or the removal sentinel reaches the wire.
func buildOnDemandThroughputResponse(odt *dbstore.OnDemandThroughput) map[string]interface{} {
	if odt == nil {
		return nil
	}
	resp := map[string]interface{}{}
	if odt.MaxReadRequestUnits != 0 {
		resp["MaxReadRequestUnits"] = odt.MaxReadRequestUnits
	}
	if odt.MaxWriteRequestUnits != 0 {
		resp["MaxWriteRequestUnits"] = odt.MaxWriteRequestUnits
	}
	if len(resp) == 0 {
		return nil
	}
	return resp
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

func buildVectorIndexesResponse(vis []*dbstore.VectorIndex) []map[string]interface{} {
	result := make([]map[string]interface{}, len(vis))
	for i, v := range vis {
		idx := map[string]interface{}{
			"IndexName":        v.IndexName,
			"IndexArn":         v.IndexArn,
			"VectorAttribute":  map[string]interface{}{"AttributeName": v.VectorAttributeName},
			"Dimensions":       v.Dimensions,
			"DistanceFunction": string(v.DistanceFunction),
			"Projection":       buildProjectionResponse(v.Projection),
			"IndexStatus":      string(v.IndexStatus),
			"ItemCount":        v.ItemCount,
			"IndexSizeBytes":   v.IndexSizeBytes,
		}
		if len(v.SearchSchema) > 0 {
			schema := make([]map[string]interface{}, len(v.SearchSchema))
			for j, e := range v.SearchSchema {
				schema[j] = map[string]interface{}{
					"AttributeName":           e.AttributeName,
					"SearchSchemaElementType": string(e.SearchSchemaElementType),
				}
			}
			idx["SearchSchema"] = schema
		}
		if v.Backfilling {
			idx["Backfilling"] = true
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
		"ProjectionType": string(p.ProjectionType),
	}
	if len(p.NonKeyAttributes) > 0 {
		resp["NonKeyAttributes"] = p.NonKeyAttributes
	}
	return resp
}
