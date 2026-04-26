package dynamodb

import (
	"bytes"
	"strings"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func skipToKeyMap(items []*dbstore.Item, exclusiveStartKey map[string]*dbstore.AttributeValue, table *dbstore.Table, indexName string) []*dbstore.Item {
	if exclusiveStartKey == nil {
		return items
	}
	for i, item := range items {
		if itemKeyMatches(item.Key, exclusiveStartKey) {
			if i+1 < len(items) {
				return items[i+1:]
			}
			return nil
		}
	}
	for i, item := range items {
		if itemKeySortsAfter(item.Key, exclusiveStartKey, table, indexName) {
			return items[i:]
		}
	}
	return nil
}

func getHashKeyName(table *dbstore.Table) string {
	for _, ks := range table.KeySchema {
		if ks.KeyType == dbstore.KeyTypeHash {
			return ks.AttributeName
		}
	}
	return ""
}

func getHashKeyNameForIndex(table *dbstore.Table, indexName string) string {
	if indexName == "" {
		return getHashKeyName(table)
	}
	for _, gsi := range table.GlobalSecondaryIndexes {
		if gsi.IndexName == indexName {
			for _, ks := range gsi.KeySchema {
				if ks.KeyType == dbstore.KeyTypeHash {
					return ks.AttributeName
				}
			}
			return ""
		}
	}
	for _, lsi := range table.LocalSecondaryIndexes {
		if lsi.IndexName == indexName {
			for _, ks := range lsi.KeySchema {
				if ks.KeyType == dbstore.KeyTypeHash {
					return ks.AttributeName
				}
			}
			return ""
		}
	}
	return ""
}

func itemKeySortsAfter(itemKey, startKey map[string]*dbstore.AttributeValue, table *dbstore.Table, indexName string) bool {
	if len(startKey) == 0 {
		return true
	}
	hashKeyName := getHashKeyNameForIndex(table, indexName)
	if hashKeyName == "" {
		return true
	}
	startVal, ok := startKey[hashKeyName]
	if ok {
		itemVal, ok := itemKey[hashKeyName]
		if !ok {
			return false
		}
		cmp := genericCompare(itemVal, startVal)
		if cmp != 0 {
			return cmp > 0
		}
	}
	sortKeyName := getSortKeyName(table, indexName)
	if sortKeyName != "" {
		if startVal, ok := startKey[sortKeyName]; ok {
			itemVal, ok := itemKey[sortKeyName]
			if !ok {
				return false
			}
			cmp := genericCompare(itemVal, startVal)
			return cmp > 0
		}
	}
	return true
}

func itemKeyMatches(itemKey, searchKey map[string]*dbstore.AttributeValue) bool {
	if itemKey == nil || searchKey == nil {
		return false
	}
	for k, v := range itemKey {
		searchV, ok := searchKey[k]
		if !ok {
			return false
		}
		if !attributeValuesEqual(v, searchV) {
			return false
		}
	}
	return true
}

func evaluateConditionExpression(item *dbstore.Item, conditionExpr string, exprAttrNames map[string]string, exprAttrValues map[string]*dbstore.AttributeValue) (bool, error) {
	if conditionExpr == "" {
		return true, nil
	}

	return evaluateConditionExpr(item, conditionExpr, exprAttrNames, exprAttrValues)
}

func evaluateConditionExpr(item *dbstore.Item, expr string, names map[string]string, values map[string]*dbstore.AttributeValue) (bool, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return true, nil
	}

	if len(expr) > 4 && strings.EqualFold(expr[:4], "NOT ") {
		inner := strings.TrimSpace(expr[4:])
		result, err := evaluateConditionExpr(item, inner, names, values)
		if err != nil {
			return false, err
		}
		return !result, nil
	}

	orParts := splitByLogicalOp(expr, " OR ")
	if len(orParts) > 1 {
		for _, part := range orParts {
			result, err := evaluateConditionExpr(item, strings.TrimSpace(part), names, values)
			if err != nil {
				return false, err
			}
			if result {
				return true, nil
			}
		}
		return false, nil
	}

	andParts := splitByLogicalOp(expr, " AND ")
	if len(andParts) > 1 {
		for _, part := range andParts {
			result, err := evaluateConditionExpr(item, strings.TrimSpace(part), names, values)
			if err != nil {
				return false, err
			}
			if !result {
				return false, nil
			}
		}
		return true, nil
	}

	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		depth := 0
		matchingPos := -1
		for i, ch := range expr {
			if ch == '(' {
				depth++
			} else if ch == ')' {
				depth--
				if depth == 0 {
					matchingPos = i
					break
				}
			}
		}
		if matchingPos == len(expr)-1 {
			return evaluateConditionExpr(item, expr[1:len(expr)-1], names, values)
		}
	}

	return evaluateSimpleCondition(item, expr, names, values)
}

func splitByLogicalOp(expr string, op string) []string {
	depth := 0
	inString := false
	var parts []string
	current := ""
	upperExpr := strings.ToUpper(expr)
	opUpper := strings.ToUpper(strings.TrimSpace(op))
	opLen := len(opUpper)
	i := 0

	for i < len(expr) {
		ch := expr[i]

		if ch == '\'' && (i == 0 || expr[i-1] != '\'') {
			if inString && i+1 < len(expr) && expr[i+1] == '\'' {
				current += "''"
				i += 2
				continue
			}
			inString = !inString
		}

		if !inString {
			if ch == '(' {
				depth++
			} else if ch == ')' {
				depth--
			}

			if depth == 0 && i+opLen <= len(upperExpr) && upperExpr[i:i+opLen] == opUpper {
				if trimmed := strings.TrimSpace(current); trimmed != "" {
					parts = append(parts, trimmed)
				}
				current = ""
				i += opLen
				continue
			}
		}

		current += string(ch)
		i++
	}

	if trimmed := strings.TrimSpace(current); trimmed != "" {
		parts = append(parts, trimmed)
	}

	return parts
}

func evaluateSimpleCondition(item *dbstore.Item, expr string, names map[string]string, values map[string]*dbstore.AttributeValue) (bool, error) {
	tokens := tokenizeExpression(expr)
	if len(tokens) == 0 {
		return true, nil
	}

	if len(tokens) == 1 {
		token := tokens[0]
		for _, funcName := range []string{"attribute_exists", "attribute_not_exists", "begins_with", "contains"} {
			prefix := funcName + "("
			if strings.HasPrefix(token, prefix) && strings.HasSuffix(token, ")") {
				argStr := token[len(prefix) : len(token)-1]
				args := strings.Split(argStr, ",")
				for i, arg := range args {
					args[i] = strings.TrimSpace(arg)
				}
				if funcName == "attribute_exists" || funcName == "attribute_not_exists" {
					if len(args) >= 1 {
						attrName := resolveName(args[0], names)
						_, exists := item.Attributes[attrName]
						if funcName == "attribute_exists" {
							return exists, nil
						}
						return !exists, nil
					}
				}
				if funcName == "begins_with" || funcName == "contains" {
					return evaluateFunctionCondition(item, []string{funcName, "(" + argStr + ")"}, names, values)
				}
			}
		}
	}

	if len(tokens) >= 2 {
		funcName := tokens[0]
		if funcName == "attribute_exists" || funcName == "attribute_not_exists" {
			attrToken := tokens[1]
			if len(attrToken) >= 2 && attrToken[0] == '(' && attrToken[len(attrToken)-1] == ')' {
				attrName := attrToken[1 : len(attrToken)-1]
				attrName = resolveName(attrName, names)
				_, exists := item.Attributes[attrName]
				if funcName == "attribute_exists" {
					return exists, nil
				}
				return !exists, nil
			}
		}

		if funcName == "begins_with" || funcName == "contains" {
			return evaluateFunctionCondition(item, tokens, names, values)
		}
	}

	if len(tokens) >= 3 {
		attrName := resolveName(tokens[0], names)
		op := tokens[1]
		value := resolveValue(tokens[2], values, names)

		attr, exists := item.Attributes[attrName]
		if !exists {
			if op == "=" && value != nil && value.NULL != nil && *value.NULL {
				return true, nil
			}
			return false, nil
		}

		return compareAttributeValues(attr, op, value), nil
	}

	return true, nil
}

func evaluateFunctionCondition(item *dbstore.Item, tokens []string, names map[string]string, values map[string]*dbstore.AttributeValue) (bool, error) {
	funcName := tokens[0]
	if len(tokens) < 2 {
		return false, nil
	}

	argStr := strings.Join(tokens[1:], " ")
	argStr = strings.Trim(argStr, "()")
	args := strings.Split(argStr, ",")
	if len(args) < 2 {
		return false, nil
	}

	path := strings.TrimSpace(args[0])
	valToken := strings.TrimSpace(args[1])

	attrName := resolveName(path, names)
	attr, exists := item.Attributes[attrName]
	if !exists {
		return false, nil
	}

	var checkValue *dbstore.AttributeValue
	if strings.HasPrefix(valToken, ":") {
		checkValue = resolveValue(valToken, values, names)
	} else {
		checkValue = &dbstore.AttributeValue{S: &valToken}
	}

	switch funcName {
	case "begins_with":
		if attr.S != nil && checkValue.S != nil {
			return strings.HasPrefix(*attr.S, *checkValue.S), nil
		}
	case "contains":
		if attr.S != nil && checkValue.S != nil {
			return strings.Contains(*attr.S, *checkValue.S), nil
		}
		if attr.SS != nil && checkValue.S != nil {
			for _, s := range attr.SS {
				if s == *checkValue.S {
					return true, nil
				}
			}
			return false, nil
		}
		if attr.NS != nil && checkValue.N != nil {
			for _, n := range attr.NS {
				normalizedN := normalizeNumberString(n)
				normalizedCheck := normalizeNumberString(*checkValue.N)
				if normalizedN == normalizedCheck {
					return true, nil
				}
			}
			return false, nil
		}
		if attr.BS != nil && checkValue.B != nil {
			for _, b := range attr.BS {
				if bytes.Equal(b, checkValue.B) {
					return true, nil
				}
			}
			return false, nil
		}
	}

	return false, nil
}
