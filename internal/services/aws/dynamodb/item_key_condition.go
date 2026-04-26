package dynamodb

import (
	"encoding/base64"
	"strings"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

type sortKeyCondition struct {
	op       string
	attrName string
	value    *dbstore.AttributeValue
	valueTo  *dbstore.AttributeValue
}

func extractPrimaryKeyCondition(table *dbstore.Table, expr string, names map[string]string, values map[string]*dbstore.AttributeValue) (string, *sortKeyCondition) {
	var hashKeyName, sortKeyName string
	var hashKeyValue string
	var sortCond *sortKeyCondition

	for _, ks := range table.KeySchema {
		if ks.KeyType == dbstore.KeyTypeHash {
			hashKeyName = ks.AttributeName
		} else if ks.KeyType == dbstore.KeyTypeRange {
			sortKeyName = ks.AttributeName
		}
	}

	tokens := tokenizeExpression(expr)
	for i := 0; i < len(tokens); i++ {
		attrName := resolveName(tokens[i], names)
		if attrName == hashKeyName && i+2 < len(tokens) && tokens[i+1] == "=" {
			val := resolveValue(tokens[i+2], values, names)
			if val != nil && val.S != nil {
				hashKeyValue = *val.S
			} else if val != nil && val.N != nil {
				hashKeyValue = *val.N
			} else if val != nil && val.B != nil {
				hashKeyValue = base64.StdEncoding.EncodeToString(val.B)
			}
			i += 2
		} else if strings.HasPrefix(strings.ToLower(tokens[i]), "begins_with(") {
			funcToken := tokens[i][len("begins_with("):]
			commaIdx := strings.Index(funcToken, ",")
			if commaIdx != -1 {
				funcAttrName := resolveName(strings.TrimSpace(funcToken[:commaIdx]), names)
				if funcAttrName == sortKeyName {
					rest := strings.TrimSuffix(strings.TrimSpace(funcToken[commaIdx+1:]), ")")
					var valueToken string
					if rest == "" && i+1 < len(tokens) {
						i++
						valueToken = strings.TrimSuffix(strings.TrimSpace(tokens[i]), ")")
					} else {
						valueToken = rest
					}
					val := resolveValue(valueToken, values, names)
					if val != nil {
						sortCond = &sortKeyCondition{op: "begins_with", attrName: sortKeyName, value: val}
					}
				}
			}
		} else if attrName == sortKeyName && i+2 < len(tokens) {
			op := strings.ToLower(tokens[i+1])
			if op == "between" {
				if i+4 < len(tokens) && strings.ToUpper(tokens[i+3]) == "AND" {
					valFrom := resolveValue(tokens[i+2], values, names)
					valTo := resolveValue(tokens[i+4], values, names)
					if valFrom != nil && valTo != nil {
						sortCond = &sortKeyCondition{
							op:       op,
							attrName: sortKeyName,
							value:    valFrom,
							valueTo:  valTo,
						}
					}
					i += 4
				} else {
					i++
				}
			} else if op == "begins_with" {
				val := resolveValue(tokens[i+2], values, names)
				if val != nil {
					sortCond = &sortKeyCondition{op: op, attrName: sortKeyName, value: val}
				}
				i += 2
			} else {
				val := resolveValue(tokens[i+2], values, names)
				if val != nil {
					sortCond = &sortKeyCondition{op: op, attrName: sortKeyName, value: val}
				}
				i += 2
			}
		}
	}

	return hashKeyValue, sortCond
}

func extractIndexKeyCondition(table *dbstore.Table, indexName, expr string, names map[string]string, values map[string]*dbstore.AttributeValue) (string, *sortKeyCondition) {
	var hashKeyName, sortKeyName string
	var hashKeyValue string
	var sortCond *sortKeyCondition

	for _, gsi := range table.GlobalSecondaryIndexes {
		if gsi.IndexName == indexName {
			for _, ks := range gsi.KeySchema {
				if ks.KeyType == dbstore.KeyTypeHash {
					hashKeyName = ks.AttributeName
				} else if ks.KeyType == dbstore.KeyTypeRange {
					sortKeyName = ks.AttributeName
				}
			}
			break
		}
	}
	for _, lsi := range table.LocalSecondaryIndexes {
		if lsi.IndexName == indexName {
			for _, ks := range table.KeySchema {
				if ks.KeyType == dbstore.KeyTypeHash {
					hashKeyName = ks.AttributeName
					break
				}
			}
			for _, ks := range lsi.KeySchema {
				if ks.KeyType == dbstore.KeyTypeRange {
					sortKeyName = ks.AttributeName
					break
				}
			}
			break
		}
	}

	tokens := tokenizeExpression(expr)
	for i := 0; i < len(tokens); i++ {
		attrName := resolveName(tokens[i], names)
		if attrName == hashKeyName && i+2 < len(tokens) && tokens[i+1] == "=" {
			val := resolveValue(tokens[i+2], values, names)
			if val != nil && val.S != nil {
				hashKeyValue = *val.S
			} else if val != nil && val.N != nil {
				hashKeyValue = *val.N
			} else if val != nil && val.B != nil {
				hashKeyValue = base64.StdEncoding.EncodeToString(val.B)
			}
			i += 2
		} else if strings.HasPrefix(strings.ToLower(tokens[i]), "begins_with(") {
			funcToken := tokens[i][len("begins_with("):]
			commaIdx := strings.Index(funcToken, ",")
			if commaIdx != -1 {
				funcAttrName := resolveName(strings.TrimSpace(funcToken[:commaIdx]), names)
				if funcAttrName == sortKeyName {
					rest := strings.TrimSuffix(strings.TrimSpace(funcToken[commaIdx+1:]), ")")
					var valueToken string
					if rest == "" && i+1 < len(tokens) {
						i++
						valueToken = strings.TrimSuffix(strings.TrimSpace(tokens[i]), ")")
					} else {
						valueToken = rest
					}
					val := resolveValue(valueToken, values, names)
					if val != nil {
						sortCond = &sortKeyCondition{op: "begins_with", attrName: sortKeyName, value: val}
					}
				}
			}
		} else if attrName == sortKeyName && i+2 < len(tokens) {
			op := strings.ToLower(tokens[i+1])
			if op == "between" {
				if i+4 < len(tokens) && strings.ToUpper(tokens[i+3]) == "AND" {
					valFrom := resolveValue(tokens[i+2], values, names)
					valTo := resolveValue(tokens[i+4], values, names)
					if valFrom != nil && valTo != nil {
						sortCond = &sortKeyCondition{
							op:       op,
							attrName: sortKeyName,
							value:    valFrom,
							valueTo:  valTo,
						}
					}
					i += 4
				} else {
					i++
				}
			} else if op == "begins_with" {
				val := resolveValue(tokens[i+2], values, names)
				if val != nil {
					sortCond = &sortKeyCondition{op: op, attrName: sortKeyName, value: val}
				}
				i += 2
			} else {
				val := resolveValue(tokens[i+2], values, names)
				if val != nil {
					sortCond = &sortKeyCondition{op: op, attrName: sortKeyName, value: val}
				}
				i += 2
			}
		}
	}

	return hashKeyValue, sortCond
}

func isGSI(table *dbstore.Table, indexName string) bool {
	for _, gsi := range table.GlobalSecondaryIndexes {
		if gsi.IndexName == indexName {
			return true
		}
	}
	return false
}

func filterBySortKeyCondition(items []*dbstore.Item, cond *sortKeyCondition) []*dbstore.Item {
	if cond == nil || cond.value == nil || cond.attrName == "" {
		return items
	}
	var result []*dbstore.Item
	for _, item := range items {
		av, ok := item.Attributes[cond.attrName]
		if !ok || av == nil {
			continue
		}
		if compareWithCondition(av, cond) {
			result = append(result, item)
		}
	}
	return result
}

func compareWithCondition(attr *dbstore.AttributeValue, cond *sortKeyCondition) bool {
	if attr == nil || cond == nil || cond.value == nil {
		return false
	}
	switch cond.op {
	case "=":
		return attributeValuesEqual(attr, cond.value)
	case "<":
		return genericCompare(attr, cond.value) < 0
	case "<=":
		return genericCompare(attr, cond.value) <= 0
	case ">":
		return genericCompare(attr, cond.value) > 0
	case ">=":
		return genericCompare(attr, cond.value) >= 0
	case "begins_with":
		if attr.S != nil && cond.value.S != nil {
			return strings.HasPrefix(*attr.S, *cond.value.S)
		}
	case "between":
		if cond.valueTo == nil {
			return false
		}
		return genericCompare(attr, cond.value) >= 0 && genericCompare(attr, cond.valueTo) <= 0
	}
	return false
}
