// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"sort"

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
	sort.Slice(result, func(i, j int) bool {
		return result[i].AttributeName < result[j].AttributeName
	})
	return result
}

// applyGSIUpdates applies GlobalSecondaryIndexUpdates to the existing index
// list and returns the updated list plus the names of indexes this request
// removes. A name that is both deleted and re-created in the same request
// is not returned as deleted: the final schema still contains it, so its
// entries must survive for the backfill to rebuild on top of.
func applyGSIUpdates(tableARN string, existing []*dbstore.GlobalSecondaryIndex, updates []interface{}) ([]*dbstore.GlobalSecondaryIndex, []string, error) {
	gsiMap := make(map[string]*dbstore.GlobalSecondaryIndex)
	for _, g := range existing {
		gsiMap[g.IndexName] = g
	}
	deletedNames := []string{}

	for _, u := range updates {
		update, ok := u.(map[string]interface{})
		if !ok {
			continue
		}

		if create, ok := update["Create"].(map[string]interface{}); ok {
			if !validateGSICreateRequired(create) {
				return nil, nil, ErrInvalidParameter
			}
			idxName := request.GetStringParam(create, "IndexName")
			keySchema := parseKeySchema(create)
			proj, err := parseProjection(create)
			if err != nil {
				return nil, nil, err
			}
			if !validateProjectionRequired(create["Projection"].(map[string]interface{})) {
				return nil, nil, ErrInvalidParameter
			}
			gsiPT := parseProvisionedThroughput(create)
			if gsiPT != nil && !validateProvisionedThroughputValues(gsiPT) {
				return nil, nil, ErrInvalidParameter
			}
			gsiMap[idxName] = &dbstore.GlobalSecondaryIndex{
				IndexName:             idxName,
				IndexArn:              tableARN + "/index/" + idxName,
				KeySchema:             keySchema,
				Projection:            proj,
				ProvisionedThroughput: gsiPT,
				IndexStatus:           dbstore.IndexStatusActive,
			}
		}

		if updateGSI, ok := update["Update"].(map[string]interface{}); ok {
			idxName := request.GetStringParam(updateGSI, "IndexName")
			if idxName == "" {
				return nil, nil, ErrInvalidParameter
			}
			if !validateResourceName(idxName) {
				return nil, nil, ErrInvalidParameter
			}
			if idx, exists := gsiMap[idxName]; exists {
				if provThroughput := parseProvisionedThroughput(updateGSI); provThroughput != nil {
					if !validateProvisionedThroughputValues(provThroughput) {
						return nil, nil, ErrInvalidParameter
					}
					idx.ProvisionedThroughput = provThroughput
				}
				idx.IndexStatus = dbstore.IndexStatusActive
			} else {
				return nil, nil, ErrIndexNotFound
			}
		}

		if deleteGSI, ok := update["Delete"].(map[string]interface{}); ok {
			idxNameToDelete := request.GetStringParam(deleteGSI, "IndexName")
			if idxNameToDelete == "" {
				return nil, nil, ErrInvalidParameter
			}
			if err := validateGSIDeleteExists(gsiMap, idxNameToDelete); err != nil {
				return nil, nil, err
			}
			delete(gsiMap, idxNameToDelete)
			deletedNames = append(deletedNames, idxNameToDelete)
		}
	}

	// An index deleted and re-created in the same request still exists in
	// the final schema; only names absent from the final map have their
	// entries cleaned up.
	finalDeletes := deletedNames[:0]
	for _, name := range deletedNames {
		if _, stillPresent := gsiMap[name]; !stillPresent {
			finalDeletes = append(finalDeletes, name)
		}
	}

	result := make([]*dbstore.GlobalSecondaryIndex, 0, len(gsiMap))
	for _, g := range gsiMap {
		result = append(result, g)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].IndexName < result[j].IndexName
	})
	return result, finalDeletes, nil
}

// applyVectorIndexUpdates applies VectorIndexUpdates to the existing vector
// index list and returns the updated list plus the names this request creates
// and removes. One UpdateTable request may add or remove exactly one vector
// index, so more than one update element is rejected.
func applyVectorIndexUpdates(tableARN string, existing []*dbstore.VectorIndex, updates []interface{}) ([]*dbstore.VectorIndex, []string, []string, error) {
	if len(updates) > 1 {
		return nil, nil, nil, ErrInvalidParameter
	}

	viMap := make(map[string]*dbstore.VectorIndex)
	for _, v := range existing {
		viMap[v.IndexName] = v
	}
	createdNames := []string{}
	deletedNames := []string{}

	for _, u := range updates {
		update, ok := u.(map[string]interface{})
		if !ok {
			return nil, nil, nil, ErrInvalidParameter
		}

		if create, ok := update["Create"].(map[string]interface{}); ok {
			idx, err := parseVectorIndex(create)
			if err != nil {
				return nil, nil, nil, err
			}
			if _, exists := viMap[idx.IndexName]; exists {
				return nil, nil, nil, ErrIndexAlreadyExists
			}
			idx.IndexArn = tableARN + "/index/" + idx.IndexName
			viMap[idx.IndexName] = idx
			createdNames = append(createdNames, idx.IndexName)
		}

		if deleteVI, ok := update["Delete"].(map[string]interface{}); ok {
			idxName := request.GetStringParam(deleteVI, "IndexName")
			if idxName == "" {
				return nil, nil, nil, ErrInvalidParameter
			}
			if _, exists := viMap[idxName]; !exists {
				return nil, nil, nil, ErrIndexNotFound
			}
			delete(viMap, idxName)
			deletedNames = append(deletedNames, idxName)
		}
	}

	result := make([]*dbstore.VectorIndex, 0, len(viMap))
	for _, v := range viMap {
		result = append(result, v)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].IndexName < result[j].IndexName
	})
	return result, createdNames, deletedNames, nil
}
