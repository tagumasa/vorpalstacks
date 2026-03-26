// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

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

func validateLSIPartitionKey(tableKeySchema []*dbstore.KeySchemaElement, lsi []*dbstore.LocalSecondaryIndex) error {
	if len(lsi) == 0 {
		return nil
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
			return ErrInvalidParameter
		}
	}
	return nil
}

func validateKeySchema(keySchema []*dbstore.KeySchemaElement) error {
	if len(keySchema) > 2 {
		return ErrInvalidParameter
	}

	hasHash := false
	for _, elem := range keySchema {
		if elem.AttributeName == "" {
			return ErrInvalidParameter
		}
		switch elem.KeyType {
		case dbstore.KeyTypeHash:
			if hasHash {
				return ErrInvalidParameter
			}
			hasHash = true
		case dbstore.KeyTypeRange:
		default:
			return ErrInvalidParameter
		}
	}

	if !hasHash {
		return ErrInvalidParameter
	}
	return nil
}

func validateAttributeDefinitions(keySchema []*dbstore.KeySchemaElement, attrDefs []*dbstore.AttributeDefinition) error {
	defMap := make(map[string]bool)
	for _, def := range attrDefs {
		if def.AttributeName == "" {
			return ErrInvalidParameter
		}
		switch def.AttributeType {
		case dbstore.ScalarAttributeTypeS, dbstore.ScalarAttributeTypeN, dbstore.ScalarAttributeTypeB:
		default:
			return ErrInvalidParameter
		}
		defMap[def.AttributeName] = true
	}

	for _, elem := range keySchema {
		if !defMap[elem.AttributeName] {
			return ErrInvalidParameter
		}
	}
	return nil
}

func validateAllKeyAttributesInDefs(keySchema []*dbstore.KeySchemaElement, gsi []*dbstore.GlobalSecondaryIndex, lsi []*dbstore.LocalSecondaryIndex, attrDefs []*dbstore.AttributeDefinition) error {
	defMap := make(map[string]bool)
	for _, def := range attrDefs {
		defMap[def.AttributeName] = true
	}

	for _, elem := range keySchema {
		if !defMap[elem.AttributeName] {
			return ErrInvalidParameter
		}
	}

	for _, idx := range gsi {
		for _, elem := range idx.KeySchema {
			if !defMap[elem.AttributeName] {
				return ErrInvalidParameter
			}
		}
	}

	for _, idx := range lsi {
		for _, elem := range idx.KeySchema {
			if !defMap[elem.AttributeName] {
				return ErrInvalidParameter
			}
		}
	}

	return nil
}
