package dynamodb

import (
	"bytes"
	"sort"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)





func getSortKeyName(table *dbstore.Table, indexName string) string {
	if indexName == "" {
		for _, ks := range table.KeySchema {
			if ks.KeyType == dbstore.KeyTypeRange {
				return ks.AttributeName
			}
		}
		return ""
	}
	for _, gsi := range table.GlobalSecondaryIndexes {
		if gsi.IndexName == indexName {
			for _, ks := range gsi.KeySchema {
				if ks.KeyType == dbstore.KeyTypeRange {
					return ks.AttributeName
				}
			}
			return ""
		}
	}
	for _, lsi := range table.LocalSecondaryIndexes {
		if lsi.IndexName == indexName {
			for _, ks := range lsi.KeySchema {
				if ks.KeyType == dbstore.KeyTypeRange {
					return ks.AttributeName
				}
			}
			return ""
		}
	}
	return ""
}

func getSortKeyType(table *dbstore.Table, sortKeyName string) string {
	for _, ad := range table.AttributeDefinitions {
		if ad.AttributeName == sortKeyName {
			return string(ad.AttributeType)
		}
	}
	return ""
}

func sortItemsBySortKeyWithIndex(table *dbstore.Table, items []*dbstore.Item, indexName string) {
	sortKeyName := getSortKeyName(table, indexName)
	if sortKeyName == "" {
		return
	}

	var basePK, baseSK string
	for _, ks := range table.KeySchema {
		if ks.KeyType == dbstore.KeyTypeHash {
			basePK = ks.AttributeName
		} else if ks.KeyType == dbstore.KeyTypeRange {
			baseSK = ks.AttributeName
		}
	}

	sort.Slice(items, func(i, j int) bool {
		avI := items[i].Attributes[sortKeyName]
		avJ := items[j].Attributes[sortKeyName]
		if avI == nil || avJ == nil {
			return false
		}

		cmp := compareAttributeValuesGeneric(avI, avJ)

		if cmp == 0 && basePK != "" {
			pkI := items[i].Attributes[basePK]
			pkJ := items[j].Attributes[basePK]
			if pkI != nil && pkJ != nil {
				cmp = compareAttributeValuesGeneric(pkI, pkJ)
			}
		}

		if cmp == 0 && baseSK != "" {
			skI := items[i].Attributes[baseSK]
			skJ := items[j].Attributes[baseSK]
			if skI != nil && skJ != nil {
				cmp = compareAttributeValuesGeneric(skI, skJ)
			}
		}

		return cmp < 0
	})
}

func compareAttributeValuesGeneric(avI, avJ *dbstore.AttributeValue) int {
	if avI.N != nil && avJ.N != nil {
		return compareNumberStrings(*avI.N, *avJ.N)
	}
	if avI.S != nil && avJ.S != nil {
		if *avI.S < *avJ.S {
			return -1
		} else if *avI.S > *avJ.S {
			return 1
		}
		return 0
	}
	if avI.B != nil && avJ.B != nil {
		return bytes.Compare(avI.B, avJ.B)
	}
	return 0
}

func sortItemsReverseBySortKeyWithIndex(table *dbstore.Table, items []*dbstore.Item, indexName string) {
	sortKeyName := getSortKeyName(table, indexName)
	if sortKeyName == "" {
		return
	}

	var basePK, baseSK string
	for _, ks := range table.KeySchema {
		if ks.KeyType == dbstore.KeyTypeHash {
			basePK = ks.AttributeName
		} else if ks.KeyType == dbstore.KeyTypeRange {
			baseSK = ks.AttributeName
		}
	}

	sort.Slice(items, func(i, j int) bool {
		avI := items[i].Attributes[sortKeyName]
		avJ := items[j].Attributes[sortKeyName]
		if avI == nil || avJ == nil {
			return false
		}

		cmp := compareAttributeValuesGeneric(avI, avJ)

		if cmp == 0 && basePK != "" {
			pkI := items[i].Attributes[basePK]
			pkJ := items[j].Attributes[basePK]
			if pkI != nil && pkJ != nil {
				cmp = compareAttributeValuesGeneric(pkI, pkJ)
			}
		}

		if cmp == 0 && baseSK != "" {
			skI := items[i].Attributes[baseSK]
			skJ := items[j].Attributes[baseSK]
			if skI != nil && skJ != nil {
				cmp = compareAttributeValuesGeneric(skI, skJ)
			}
		}

		return cmp > 0
	})
}
