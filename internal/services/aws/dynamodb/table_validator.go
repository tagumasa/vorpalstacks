package dynamodb

import (
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// validateIndexNameUniqueness reports whether all GSI and LSI index names
// within a table definition are distinct. Returns the typed
// ErrIndexAlreadyExists sentinel so the caller can surface
// ResourceInUseException with a precise message rather than the generic
// ValidationException.
func validateIndexNameUniqueness(gsi []*dbstore.GlobalSecondaryIndex, lsi []*dbstore.LocalSecondaryIndex) error {
	indexNames := make(map[string]bool)
	for _, idx := range gsi {
		if idx.IndexName == "" {
			return ErrInvalidParameter
		}
		if indexNames[idx.IndexName] {
			return ErrIndexAlreadyExists
		}
		indexNames[idx.IndexName] = true
	}
	for _, idx := range lsi {
		if idx.IndexName == "" {
			return ErrInvalidParameter
		}
		if indexNames[idx.IndexName] {
			return ErrIndexAlreadyExists
		}
		indexNames[idx.IndexName] = true
	}
	return nil
}

// validateLSIPartitionKey reports whether every LSI shares the same hash
// key as the base table. Returns true when no LSIs are defined.
func validateLSIPartitionKey(tableKeySchema []*dbstore.KeySchemaElement, lsi []*dbstore.LocalSecondaryIndex) bool {
	if len(lsi) == 0 {
		return true
	}

	var tablePartitionKey string
	for _, elem := range tableKeySchema {
		if elem.KeyType == dbstore.KeyTypeHash {
			tablePartitionKey = elem.AttributeName
			break
		}
	}

	for _, idx := range lsi {
		var idxPartitionKey string
		for _, elem := range idx.KeySchema {
			if elem.KeyType == dbstore.KeyTypeHash {
				idxPartitionKey = elem.AttributeName
				break
			}
		}
		if idxPartitionKey != tablePartitionKey {
			return false
		}
	}
	return true
}

// validateKeySchema reports whether keySchema satisfies the DynamoDB
// constraints: at most 2 elements, exactly one HASH, optional RANGE,
// non-empty attribute names, and only HASH/RANGE key types.
func validateKeySchema(keySchema []*dbstore.KeySchemaElement) bool {
	if len(keySchema) > 2 {
		return false
	}

	hasHash := false
	for _, elem := range keySchema {
		if elem.AttributeName == "" {
			return false
		}
		switch elem.KeyType {
		case dbstore.KeyTypeHash:
			if hasHash {
				return false
			}
			hasHash = true
		case dbstore.KeyTypeRange:
		default:
			return false
		}
	}

	return hasHash
}

// validateAttributeDefinitions reports whether every key-schema attribute
// has a matching entry in attrDefs with a valid scalar type
// (S, N, or B), and whether all attrDefs entries have non-empty names and
// valid types.
func validateAttributeDefinitions(keySchema []*dbstore.KeySchemaElement, attrDefs []*dbstore.AttributeDefinition) bool {
	defMap := make(map[string]bool)
	for _, def := range attrDefs {
		if def.AttributeName == "" {
			return false
		}
		switch def.AttributeType {
		case dbstore.ScalarAttributeTypeS, dbstore.ScalarAttributeTypeN, dbstore.ScalarAttributeTypeB:
		default:
			return false
		}
		defMap[def.AttributeName] = true
	}

	for _, elem := range keySchema {
		if !defMap[elem.AttributeName] {
			return false
		}
	}
	return true
}

// validateAllKeyAttributesInDefs reports whether every key attribute in the
// base table key schema, GSI key schemas, and LSI key schemas has a matching
// entry in attrDefs.
func validateAllKeyAttributesInDefs(keySchema []*dbstore.KeySchemaElement, gsi []*dbstore.GlobalSecondaryIndex, lsi []*dbstore.LocalSecondaryIndex, attrDefs []*dbstore.AttributeDefinition) bool {
	defMap := make(map[string]bool)
	for _, def := range attrDefs {
		defMap[def.AttributeName] = true
	}

	for _, elem := range keySchema {
		if !defMap[elem.AttributeName] {
			return false
		}
	}

	for _, idx := range gsi {
		for _, elem := range idx.KeySchema {
			if !defMap[elem.AttributeName] {
				return false
			}
		}
	}

	for _, idx := range lsi {
		for _, elem := range idx.KeySchema {
			if !defMap[elem.AttributeName] {
				return false
			}
		}
	}

	return true
}
