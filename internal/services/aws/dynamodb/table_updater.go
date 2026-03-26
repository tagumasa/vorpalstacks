// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"vorpalstacks/internal/services/aws/common/request"
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
