package dynamodb

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// getNestedAttributeValue resolves a document path (e.g. "a.b.c" or
// "a[0].b") against an item's attributes, navigating Map and List
// AttributeValue types. Returns nil if any segment is missing.
func getNestedAttributeValue(attrs map[string]*dbstore.AttributeValue, path string) *dbstore.AttributeValue {
	if attrs == nil || path == "" {
		return nil
	}
	// Fast path: top-level attribute (no dots or brackets).
	if !strings.ContainsAny(path, ".[") {
		return attrs[path]
	}
	parts := splitDocPath(path)
	var current *dbstore.AttributeValue
	ok := false
	for i, part := range parts {
		if i == 0 {
			current, ok = attrs[part.name]
			if !ok || current == nil {
				return nil
			}
		} else {
			if current == nil {
				return nil
			}
			if part.name != "" {
				if current.M != nil {
					v, exists := current.M[part.name]
					if !exists || v == nil {
						return nil
					}
					current = v
				} else {
					return nil
				}
			}
		}
		if part.index >= 0 && current != nil {
			if current.L == nil || part.index >= len(current.L) {
				return nil
			}
			current = current.L[part.index]
		}
	}
	return current
}

type docPathPart struct {
	name  string
	index int
}

func splitDocPath(path string) []docPathPart {
	var parts []docPathPart
	segment := ""
	i := 0
	for i < len(path) {
		ch := path[i]
		if ch == '.' {
			if segment != "" {
				parts = append(parts, docPathPart{name: segment, index: -1})
				segment = ""
			}
			i++
			continue
		}
		if ch == '[' {
			if segment != "" {
				parts = append(parts, docPathPart{name: segment, index: -1})
				segment = ""
			}
			// Parse the index until ']'
			j := i + 1
			for j < len(path) && path[j] != ']' {
				j++
			}
			if j < len(path) {
				idxStr := path[i+1 : j]
				idx, err := strconv.Atoi(idxStr)
				if err == nil {
					if len(parts) > 0 && parts[len(parts)-1].index == -1 {
						parts[len(parts)-1].index = idx
					} else {
						parts = append(parts, docPathPart{name: "", index: idx})
					}
				}
				i = j + 1
				continue
			}
		}
		segment += string(ch)
		i++
	}
	if segment != "" {
		parts = append(parts, docPathPart{name: segment, index: -1})
	}
	return parts
}

// attributeNestedExists checks whether a document path exists in the item.
func attributeNestedExists(attrs map[string]*dbstore.AttributeValue, path string) bool {
	return getNestedAttributeValue(attrs, path) != nil
}

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
	if len(itemKey) != len(searchKey) {
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
	opUpper := strings.ToUpper(op)
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
				// N1: When splitting by " AND ", skip ANDs that are part
				// of a BETWEEN expression. The accumulated text ending
				// with "<token> BETWEEN <token>" means this AND is the
				// BETWEEN separator, not a logical conjunction.
				if opUpper == " AND " {
					fields := strings.Fields(strings.TrimSpace(current))
					if len(fields) >= 2 && strings.EqualFold(fields[len(fields)-2], "BETWEEN") {
						current += string(ch)
						i++
						continue
					}
				}
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
	expr = strings.TrimSpace(expr)

	// NOT handling — after OR/AND split, so NOT applies only to the
	// immediate operand. "NOT a AND b" is split to ["NOT a", "b"],
	// and each part is evaluated separately. NOT negates only its own
	// operand, not the entire expression.
	if len(expr) > 4 && strings.EqualFold(expr[:4], "NOT ") {
		inner := strings.TrimSpace(expr[4:])
		result, err := evaluateConditionExpr(item, inner, names, values)
		if err != nil {
			return false, err
		}
		return !result, nil
	}
	if len(expr) > 4 && strings.EqualFold(expr[:3], "NOT") && expr[3] == '(' {
		inner := expr[3:]
		if strings.HasPrefix(inner, "(") && strings.HasSuffix(inner, ")") {
			inner = strings.TrimSpace(inner[1 : len(inner)-1])
		}
		result, err := evaluateConditionExpr(item, inner, names, values)
		if err != nil {
			return false, err
		}
		return !result, nil
	}

	tokens := tokenizeExpression(expr)
	if len(tokens) == 0 {
		return true, nil
	}

	if len(tokens) == 1 {
		token := tokens[0]
		for _, funcName := range []string{"attribute_exists", "attribute_not_exists", "attribute_type", "begins_with", "contains"} {
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
						exists := attributeNestedExists(item.Attributes, attrName)
						if funcName == "attribute_exists" {
							return exists, nil
						}
						return !exists, nil
					}
				}
				if funcName == "attribute_type" {
					return evaluateAttributeType(item, args, names, values)
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
				exists := attributeNestedExists(item.Attributes, attrName)
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

	if len(tokens) >= 5 && strings.EqualFold(tokens[1], "BETWEEN") && strings.EqualFold(tokens[3], "AND") {
		attrName := resolveName(tokens[0], names)
		attr := getNestedAttributeValue(item.Attributes, attrName)
		if attr == nil {
			return false, nil
		}
		low := resolveValue(tokens[2], values, names)
		high := resolveValue(tokens[4], values, names)
		if low == nil || high == nil {
			return false, nil
		}
		return compareAttributeValues(attr, ">=", low) && compareAttributeValues(attr, "<=", high), nil
	}

	if len(tokens) >= 4 && strings.EqualFold(tokens[1], "IN") {
		attrName := resolveName(tokens[0], names)
		attr := getNestedAttributeValue(item.Attributes, attrName)
		if attr == nil {
			return false, nil
		}
		for i := 2; i < len(tokens); i++ {
			t := tokens[i]
			if t == "(" || t == ")" || t == "," {
				continue
			}
			v := resolveValue(t, values, names)
			if v != nil && attributeValuesEqual(attr, v) {
				return true, nil
			}
		}
		return false, nil
	}

	if len(tokens) >= 3 {
		if strings.HasPrefix(tokens[0], "size(") && strings.HasSuffix(tokens[0], ")") {
			pathStr := tokens[0][5 : len(tokens[0])-1]
			attrName := resolveName(pathStr, names)
			attr := getNestedAttributeValue(item.Attributes, attrName)
			if attr == nil {
				return false, nil
			}
			size, ok := computeAttributeSize(attr)
			if !ok {
				return false, nil
			}
			sizeStr := strconv.Itoa(size)
			sizeAttr := &dbstore.AttributeValue{N: &sizeStr}
			value := resolveValue(tokens[2], values, names)
			return compareAttributeValues(sizeAttr, tokens[1], value), nil
		}

		attrName := resolveName(tokens[0], names)
		op := tokens[1]
		if !isValidComparisonOperator(op) {
			return false, fmt.Errorf("unsupported comparison operator: %s", op)
		}
		value := resolveValue(tokens[2], values, names)

		attr := getNestedAttributeValue(item.Attributes, attrName)
		if attr == nil {
			if op == "=" && value != nil && value.NULL != nil && *value.NULL {
				return true, nil
			}
			return false, nil
		}

		return compareAttributeValues(attr, op, value), nil
	}

	return false, fmt.Errorf("unsupported condition expression: %s", expr)
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
	attr := getNestedAttributeValue(item.Attributes, attrName)
	if attr == nil {
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
		if attr.L != nil {
			return listContainsValue(attr.L, checkValue), nil
		}
	}

	return false, nil
}

func evaluateAttributeType(item *dbstore.Item, args []string, names map[string]string, values map[string]*dbstore.AttributeValue) (bool, error) {
	if len(args) < 2 {
		return false, nil
	}

	path := strings.TrimSpace(args[0])
	typeToken := strings.TrimSpace(args[1])

	attrName := resolveName(path, names)
	attr := getNestedAttributeValue(item.Attributes, attrName)
	if attr == nil {
		return false, nil
	}

	var typeStr string
	if strings.HasPrefix(typeToken, ":") {
		checkValue := resolveValue(typeToken, values, names)
		if checkValue != nil && checkValue.S != nil {
			typeStr = *checkValue.S
		}
	} else {
		typeStr = strings.Trim(typeToken, "'")
	}

	return getAttributeTypeName(attr) == typeStr, nil
}

func getAttributeTypeName(attr *dbstore.AttributeValue) string {
	switch {
	case attr.S != nil:
		return "S"
	case attr.SS != nil:
		return "SS"
	case attr.N != nil:
		return "N"
	case attr.NS != nil:
		return "NS"
	case attr.B != nil:
		return "B"
	case attr.BS != nil:
		return "BS"
	case attr.BOOL != nil:
		return "BOOL"
	case attr.NULL != nil:
		return "NULL"
	case attr.L != nil:
		return "L"
	case attr.M != nil:
		return "M"
	}
	return ""
}

func computeAttributeSize(attr *dbstore.AttributeValue) (int, bool) {
	switch {
	case attr.S != nil:
		return len(*attr.S), true
	case attr.B != nil:
		return len(attr.B), true
	case attr.SS != nil:
		return len(attr.SS), true
	case attr.NS != nil:
		return len(attr.NS), true
	case attr.BS != nil:
		return len(attr.BS), true
	case attr.L != nil:
		return len(attr.L), true
	case attr.M != nil:
		return len(attr.M), true
	}
	return 0, false
}

func listContainsValue(list []*dbstore.AttributeValue, target *dbstore.AttributeValue) bool {
	for _, elem := range list {
		if attributeValuesEqual(elem, target) {
			return true
		}
	}
	return false
}
