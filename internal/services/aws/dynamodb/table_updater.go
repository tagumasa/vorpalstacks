// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func mergeAttributeDefinitions(existing, newDefs []*dbstore.AttributeDefinition) []*dbstore.AttributeDefinition {
	attrMap := make(map[string]*dbstore.AttributeDefinition)
	for _, a := range existing {
		attrMap[a.AttributeName] = a
	}
	for _, a := range newDefs {
		attrMap[a.AttributeName] = a
	}
	result := make([]*dbstore.AttributeDefinition, 0, len(attrMap))
	for _, a := range attrMap {
		result = append(result, a)
	}
	return result
}

func deepCopyTable(t *dbstore.Table) *dbstore.Table {
	if t == nil {
		return nil
	}
	cp := *t
	if t.KeySchema != nil {
		cp.KeySchema = make([]*dbstore.KeySchemaElement, len(t.KeySchema))
		for i, ks := range t.KeySchema {
			kse := *ks
			cp.KeySchema[i] = &kse
		}
	}
	if t.AttributeDefinitions != nil {
		cp.AttributeDefinitions = make([]*dbstore.AttributeDefinition, len(t.AttributeDefinitions))
		for i, ad := range t.AttributeDefinitions {
			ade := *ad
			cp.AttributeDefinitions[i] = &ade
		}
	}
	if t.ProvisionedThroughput != nil {
		pt := *t.ProvisionedThroughput
		cp.ProvisionedThroughput = &pt
	}
	if t.GlobalSecondaryIndexes != nil {
		cp.GlobalSecondaryIndexes = make([]*dbstore.GlobalSecondaryIndex, len(t.GlobalSecondaryIndexes))
		for i, gsi := range t.GlobalSecondaryIndexes {
			g := *gsi
			if gsi.KeySchema != nil {
				g.KeySchema = make([]*dbstore.KeySchemaElement, len(gsi.KeySchema))
				for j, ks := range gsi.KeySchema {
					kse := *ks
					g.KeySchema[j] = &kse
				}
			}
			if gsi.Projection != nil {
				p := *gsi.Projection
				if gsi.Projection.NonKeyAttributes != nil {
					p.NonKeyAttributes = make([]string, len(gsi.Projection.NonKeyAttributes))
					copy(p.NonKeyAttributes, gsi.Projection.NonKeyAttributes)
				}
				g.Projection = &p
			}
			if gsi.ProvisionedThroughput != nil {
				pt := *gsi.ProvisionedThroughput
				g.ProvisionedThroughput = &pt
			}
			cp.GlobalSecondaryIndexes[i] = &g
		}
	}
	if t.LocalSecondaryIndexes != nil {
		cp.LocalSecondaryIndexes = make([]*dbstore.LocalSecondaryIndex, len(t.LocalSecondaryIndexes))
		for i, lsi := range t.LocalSecondaryIndexes {
			l := *lsi
			if lsi.KeySchema != nil {
				l.KeySchema = make([]*dbstore.KeySchemaElement, len(lsi.KeySchema))
				for j, ks := range lsi.KeySchema {
					kse := *ks
					l.KeySchema[j] = &kse
				}
			}
			if lsi.Projection != nil {
				p := *lsi.Projection
				if lsi.Projection.NonKeyAttributes != nil {
					p.NonKeyAttributes = make([]string, len(lsi.Projection.NonKeyAttributes))
					copy(p.NonKeyAttributes, lsi.Projection.NonKeyAttributes)
				}
				l.Projection = &p
			}
			cp.LocalSecondaryIndexes[i] = &l
		}
	}
	if t.StreamSpecification != nil {
		ss := *t.StreamSpecification
		cp.StreamSpecification = &ss
	}
	if t.SSEDescription != nil {
		sd := *t.SSEDescription
		cp.SSEDescription = &sd
	}
	if t.TimeToLive != nil {
		ttl := *t.TimeToLive
		cp.TimeToLive = &ttl
	}
	if t.PointInTimeRecovery != nil {
		pitr := *t.PointInTimeRecovery
		cp.PointInTimeRecovery = &pitr
	}
	return &cp
}

func applyGSIUpdates(tableARN string, existing []*dbstore.GlobalSecondaryIndex, updates []interface{}) []*dbstore.GlobalSecondaryIndex {
	gsiMap := make(map[string]*dbstore.GlobalSecondaryIndex)
	for _, g := range existing {
		gsiMap[g.IndexName] = g
	}

	for _, u := range updates {
		update, ok := u.(map[string]interface{})
		if !ok {
			continue
		}

		if create, ok := update["Create"].(map[string]interface{}); ok {
			idxName := request.GetStringParam(create, "IndexName")
			if idxName != "" {
				gsiMap[idxName] = &dbstore.GlobalSecondaryIndex{
					IndexName:             idxName,
					IndexArn:              tableARN + "/index/" + idxName,
					KeySchema:             parseKeySchema(create),
					Projection:            parseProjection(create),
					ProvisionedThroughput: parseProvisionedThroughput(create),
					IndexStatus:           dbstore.IndexStatusActive,
				}
			}
		}

		if updateGSI, ok := update["Update"].(map[string]interface{}); ok {
			idxName := request.GetStringParam(updateGSI, "IndexName")
			if idx, exists := gsiMap[idxName]; exists {
				if provThroughput := parseProvisionedThroughput(updateGSI); provThroughput != nil {
					idx.ProvisionedThroughput = provThroughput
				}
				idx.IndexStatus = dbstore.IndexStatusActive
			}
		}

		if deleteGSI, ok := update["Delete"].(map[string]interface{}); ok {
			idxNameToDelete := request.GetStringParam(deleteGSI, "IndexName")
			delete(gsiMap, idxNameToDelete)
		}
	}

	result := make([]*dbstore.GlobalSecondaryIndex, 0, len(gsiMap))
	for _, g := range gsiMap {
		result = append(result, g)
	}
	return result
}
