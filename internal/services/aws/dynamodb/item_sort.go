package dynamodb

import (
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
	sortItemsBySortKeyWithIndexDirection(table, items, indexName, true)
}

func sortItemsReverseBySortKeyWithIndex(table *dbstore.Table, items []*dbstore.Item, indexName string) {
	sortItemsBySortKeyWithIndexDirection(table, items, indexName, false)
}

func sortItemsBySortKeyWithIndexDirection(table *dbstore.Table, items []*dbstore.Item, indexName string, ascending bool) {
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

		cmp := genericCompare(avI, avJ)

		if cmp == 0 && basePK != "" {
			pkI := items[i].Attributes[basePK]
			pkJ := items[j].Attributes[basePK]
			if pkI != nil && pkJ != nil {
				cmp = genericCompare(pkI, pkJ)
			}
		}

		if cmp == 0 && baseSK != "" {
			skI := items[i].Attributes[baseSK]
			skJ := items[j].Attributes[baseSK]
			if skI != nil && skJ != nil {
				cmp = genericCompare(skI, skJ)
			}
		}

		if ascending {
			return cmp < 0
		}
		return cmp > 0
	})
}
