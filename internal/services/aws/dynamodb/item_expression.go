package dynamodb

import (
	"fmt"
	"math/big"
	"strings"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func applyUpdateExpression(attrs map[string]*dbstore.AttributeValue, expr string, names map[string]string, values map[string]*dbstore.AttributeValue) error {
	_, err := applyUpdateExpressionWithTracking(attrs, expr, names, values)
	return err
}

func applyUpdateExpressionWithTracking(attrs map[string]*dbstore.AttributeValue, expr string, names map[string]string, values map[string]*dbstore.AttributeValue) ([]string, error) {
	var updatedAttrs []string
	tokens := tokenizeUpdateExpression(expr)
	i := 0
	for i < len(tokens) {
		action := tokens[i]
		switch action {
		case "SET":
			i++
			for i < len(tokens) && tokens[i] != "REMOVE" && tokens[i] != "ADD" && tokens[i] != "DELETE" {
				if i >= len(tokens) {
					return nil, ErrInvalidParameter
				}
				if i+1 >= len(tokens) || tokens[i+1] != "=" {
					return nil, ErrInvalidParameter
				}
				if i+2 >= len(tokens) {
					return nil, ErrInvalidParameter
				}
				path := resolveName(tokens[i], names)
				exprTokens := []string{tokens[i+2]}
				j := i + 3
				for j < len(tokens) && (tokens[j] == "+" || tokens[j] == "-" || tokens[j] == "*") && j+1 < len(tokens) {
					exprTokens = append(exprTokens, tokens[j], tokens[j+1])
					j += 2
				}
				value, valErr := resolveValueWithIfNotExists(strings.Join(exprTokens, " "), values, names, attrs)
				if valErr != nil {
					return nil, valErr
				}
				if value != nil {
					if err := setNestedValue(attrs, path, value); err != nil {
						return nil, err
					}
					updatedAttrs = append(updatedAttrs, getTopLevelAttr(path))
				}
				i = j
				if i < len(tokens) && tokens[i] == "," {
					i++
				}
			}
		case "REMOVE":
			i++
			for i < len(tokens) && tokens[i] != "SET" && tokens[i] != "ADD" && tokens[i] != "DELETE" {
				if i >= len(tokens) {
					return nil, ErrInvalidParameter
				}
				path := resolveName(tokens[i], names)
				removeNestedValue(attrs, path)
				updatedAttrs = append(updatedAttrs, getTopLevelAttr(path))
				i++
				if i < len(tokens) && tokens[i] == "," {
					i++
				}
			}
		case "ADD":
			i++
			for i < len(tokens) && tokens[i] != "SET" && tokens[i] != "REMOVE" && tokens[i] != "DELETE" {
				if i+1 >= len(tokens) {
					return nil, ErrInvalidParameter
				}
				attrName := resolveName(tokens[i], names)
				value := resolveValue(tokens[i+1], values, names)
				if value != nil {
					if changed, err := applyAddActionWithTracking(attrs, attrName, value); err != nil {
						return nil, err
					} else if changed {
						updatedAttrs = append(updatedAttrs, attrName)
					}
				}
				i += 2
				if i < len(tokens) && tokens[i] == "," {
					i++
				}
			}
		case "DELETE":
			i++
			for i < len(tokens) && tokens[i] != "SET" && tokens[i] != "REMOVE" && tokens[i] != "ADD" {
				if i+1 >= len(tokens) {
					return nil, ErrInvalidParameter
				}
				attrName := resolveName(tokens[i], names)
				value := resolveValue(tokens[i+1], values, names)
				if value != nil {
					if applyDeleteActionWithTracking(attrs, attrName, value) {
						updatedAttrs = append(updatedAttrs, attrName)
					}
				}
				i += 2
				if i < len(tokens) && tokens[i] == "," {
					i++
				}
			}
		default:
			return nil, ErrInvalidParameter
		}
	}
	return updatedAttrs, nil
}

func getTopLevelAttr(path string) string {
	if idx := strings.Index(path, "."); idx != -1 {
		return path[:idx]
	}
	if idx := strings.Index(path, "["); idx != -1 {
		return path[:idx]
	}
	return path
}

func setNestedValue(attrs map[string]*dbstore.AttributeValue, path string, value *dbstore.AttributeValue) error {
	parts := parsePathParts(path)
	if len(parts) == 0 {
		return nil
	}

	if len(parts) == 1 {
		attrs[parts[0].name] = value
		return nil
	}

	current, exists := attrs[parts[0].name]
	nextPartIsIndex := len(parts) > 1 && parts[1].isIndex
	if !exists {
		if nextPartIsIndex {
			return ErrInvalidParameter
		}
		current = dbstore.MapValue(make(map[string]*dbstore.AttributeValue))
		attrs[parts[0].name] = current
	} else if nextPartIsIndex && current.L == nil {
		return ErrInvalidParameter
	} else if !nextPartIsIndex && current.M == nil {
		current = dbstore.MapValue(make(map[string]*dbstore.AttributeValue))
		attrs[parts[0].name] = current
	}

	return setNestedValueRecursive(current, parts[1:], value)
}

func setNestedValueRecursive(current *dbstore.AttributeValue, parts []pathPart, value *dbstore.AttributeValue) error {
	if len(parts) == 0 {
		return nil
	}

	part := parts[0]
	isLast := len(parts) == 1

	if part.isIndex {
		if current.L == nil || part.index >= len(current.L) {
			return ErrInvalidParameter
		}
		if isLast {
			current.L[part.index] = value
			return nil
		} else {
			return setNestedValueRecursive(current.L[part.index], parts[1:], value)
		}
	} else {
		if current.M == nil {
			if current.S != nil || current.N != nil || current.B != nil ||
				current.BOOL != nil || current.NULL != nil || current.L != nil ||
				current.SS != nil || current.NS != nil || current.BS != nil {
				return ErrInvalidParameter
			}
			current.M = make(map[string]*dbstore.AttributeValue)
		}
		if isLast {
			current.M[part.name] = value
			return nil
		} else {
			if _, ok := current.M[part.name]; !ok {
				current.M[part.name] = dbstore.MapValue(make(map[string]*dbstore.AttributeValue))
			}
			return setNestedValueRecursive(current.M[part.name], parts[1:], value)
		}
	}
}

func removeNestedValue(attrs map[string]*dbstore.AttributeValue, path string) {
	parts := parsePathParts(path)
	if len(parts) == 0 {
		return
	}

	if len(parts) == 1 {
		delete(attrs, parts[0].name)
		return
	}

	current, exists := attrs[parts[0].name]
	if !exists {
		return
	}

	removeNestedValueRecursive(current, parts[1:])
}

func removeNestedValueRecursive(current *dbstore.AttributeValue, parts []pathPart) {
	if len(parts) == 0 {
		return
	}

	part := parts[0]
	isLast := len(parts) == 1

	if part.isIndex {
		if current.L == nil || part.index >= len(current.L) {
			return
		}
		if isLast {
			newList := make([]*dbstore.AttributeValue, 0, len(current.L)-1)
			for i, v := range current.L {
				if i != part.index {
					newList = append(newList, v)
				}
			}
			current.L = newList
		} else {
			removeNestedValueRecursive(current.L[part.index], parts[1:])
		}
	} else {
		if current.M == nil {
			return
		}
		if isLast {
			delete(current.M, part.name)
		} else {
			if next, ok := current.M[part.name]; ok {
				removeNestedValueRecursive(next, parts[1:])
			}
		}
	}
}

type pathPart struct {
	name    string
	index   int
	isIndex bool
}

func parsePathParts(path string) []pathPart {
	var parts []pathPart
	current := ""
	i := 0

	for i < len(path) {
		c := path[i]
		switch c {
		case '.':
			if current != "" {
				parts = append(parts, pathPart{name: current})
				current = ""
			}
			i++
		case '[':
			if current != "" {
				parts = append(parts, pathPart{name: current})
				current = ""
			}
			j := i + 1
			for j < len(path) && path[j] != ']' {
				j++
			}
			idxStr := path[i+1 : j]
			idx := 0
			for _, ch := range idxStr {
				if ch >= '0' && ch <= '9' {
					idx = idx*10 + int(ch-'0')
				}
			}
			parts = append(parts, pathPart{index: idx, isIndex: true})
			i = j + 1
		default:
			current += string(c)
			i++
		}
	}

	if current != "" {
		parts = append(parts, pathPart{name: current})
	}

	return parts
}

func addNumbers(a, b string) string {
	numA, okA := new(big.Rat).SetString(a)
	numB, okB := new(big.Rat).SetString(b)
	if !okA || !okB {
		return ""
	}
	numA.Add(numA, numB)
	return normalizeNumber(numA)
}

func normalizeNumber(r *big.Rat) string {
	num := r.Num()
	den := r.Denom()
	if den.Sign() == 0 {
		return "0"
	}
	if den.Cmp(big.NewInt(1)) == 0 {
		return num.String()
	}
	quotient := new(big.Int).Quo(num, den)
	remainder := new(big.Int).Rem(num, den)
	absRemainder := new(big.Int).Abs(remainder)
	if remainder.Sign() == 0 {
		return quotient.String()
	}
	scaledNum := new(big.Int).Mul(absRemainder, new(big.Int).Exp(big.NewInt(10), big.NewInt(38), nil))
	scaledQuotient := new(big.Int).Div(scaledNum, den)
	if scaledQuotient.Sign() == 0 {
		return quotient.String()
	}
	decimalStr := scaledQuotient.String()
	for len(decimalStr) < 38 {
		decimalStr = "0" + decimalStr
	}
	decimalStr = strings.TrimRight(decimalStr, "0")
	if quotient.Sign() == 0 && num.Sign() < 0 {
		return "-" + quotient.String() + "." + decimalStr
	}
	return quotient.String() + "." + decimalStr
}

func applyAddActionWithTracking(attrs map[string]*dbstore.AttributeValue, attrName string, value *dbstore.AttributeValue) (bool, error) {
	existing, exists := attrs[attrName]
	if !exists {
		attrs[attrName] = value
		return true, nil
	}

	if value.N != nil && existing.N != nil {
		result := addNumbers(*existing.N, *value.N)
		if result != "" {
			attrs[attrName] = &dbstore.AttributeValue{N: &result}
			return true, nil
		}
		return false, nil
	}

	if value.SS != nil && existing.SS != nil {
		existingSet := make(map[string]bool)
		for _, s := range existing.SS {
			existingSet[s] = true
		}
		for _, s := range value.SS {
			if !existingSet[s] {
				existing.SS = append(existing.SS, s)
			}
		}
		return len(value.SS) > 0, nil
	}

	if value.NS != nil && existing.NS != nil {
		existingSet := make(map[string]bool)
		for _, n := range existing.NS {
			existingSet[normalizeNumberString(n)] = true
		}
		for _, n := range value.NS {
			normalized := normalizeNumberString(n)
			if !existingSet[normalized] {
				existing.NS = append(existing.NS, normalized)
			}
		}
		return len(value.NS) > 0, nil
	}

	if value.BS != nil && existing.BS != nil {
		existingSet := make(map[string]bool)
		for _, b := range existing.BS {
			existingSet[string(b)] = true
		}
		for _, b := range value.BS {
			if !existingSet[string(b)] {
				existing.BS = append(existing.BS, b)
			}
		}
		return len(value.BS) > 0, nil
	}

	return false, ErrInvalidParameter
}

func applyDeleteActionWithTracking(attrs map[string]*dbstore.AttributeValue, attrName string, value *dbstore.AttributeValue) bool {
	existing, exists := attrs[attrName]
	if !exists {
		return false
	}

	if value.SS != nil && existing.SS != nil {
		toDelete := make(map[string]bool)
		for _, s := range value.SS {
			toDelete[s] = true
		}
		var newSS []string
		for _, s := range existing.SS {
			if !toDelete[s] {
				newSS = append(newSS, s)
			}
		}
		if len(newSS) == 0 {
			delete(attrs, attrName)
		} else {
			existing.SS = newSS
		}
		return true
	}

	if value.NS != nil && existing.NS != nil {
		toDelete := make(map[string]bool)
		for _, n := range value.NS {
			toDelete[normalizeNumberString(n)] = true
		}
		var newNS []string
		for _, n := range existing.NS {
			if !toDelete[normalizeNumberString(n)] {
				newNS = append(newNS, n)
			}
		}
		if len(newNS) == 0 {
			delete(attrs, attrName)
		} else {
			existing.NS = newNS
		}
		return true
	}

	if value.BS != nil && existing.BS != nil {
		toDelete := make(map[string]bool)
		for _, b := range value.BS {
			toDelete[string(b)] = true
		}
		var newBS [][]byte
		for _, b := range existing.BS {
			if !toDelete[string(b)] {
				newBS = append(newBS, b)
			}
		}
		if len(newBS) == 0 {
			delete(attrs, attrName)
		} else {
			existing.BS = newBS
		}
		return true
	}
	return false
}

func applyAttributeUpdatesWithTracking(attrs map[string]*dbstore.AttributeValue, updates interface{}) ([]string, error) {
	updatesMap, ok := updates.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	var updatedAttrs []string
	for attrName, update := range updatesMap {
		updateAction, ok := update.(map[string]interface{})
		if !ok {
			continue
		}

		if action, ok := updateAction["Action"].(string); ok {
			switch action {
			case "PUT":
				if value, ok := updateAction["Value"]; ok {
					attrs[attrName] = parseAttributeValue(value)
					updatedAttrs = append(updatedAttrs, attrName)
				}
			case "DELETE":
				delete(attrs, attrName)
				updatedAttrs = append(updatedAttrs, attrName)
			case "ADD":
				if value, ok := updateAction["Value"]; ok {
					addValue := parseAttributeValue(value)
					if addValue == nil {
						continue
					}
					if _, exists := attrs[attrName]; !exists {
						attrs[attrName] = addValue
						updatedAttrs = append(updatedAttrs, attrName)
						continue
					}
					if changed, err := applyAddActionWithTracking(attrs, attrName, addValue); err != nil {
						return nil, err
					} else if changed {
						updatedAttrs = append(updatedAttrs, attrName)
					}
				}
			}
		}
	}
	return updatedAttrs, nil
}

func tokenizeUpdateExpression(expr string) []string {
	var tokens []string
	current := ""
	inQuotes := false
	parenDepth := 0

	for _, c := range expr {
		switch c {
		case ' ', '\t', '\n':
			if !inQuotes && parenDepth == 0 && current != "" {
				tokens = append(tokens, current)
				current = ""
			}
		case ',':
			if !inQuotes && parenDepth == 0 {
				if current != "" {
					tokens = append(tokens, current)
					current = ""
				}
				tokens = append(tokens, ",")
			} else {
				current += string(c)
			}
		case '"':
			inQuotes = !inQuotes
			current += string(c)
		case '(':
			parenDepth++
			current += string(c)
		case ')':
			parenDepth--
			current += string(c)
		default:
			current += string(c)
		}
	}

	if current != "" {
		tokens = append(tokens, current)
	}

	return tokens
}

func resolveName(token string, names map[string]string) string {
	if names != nil {
		if resolved, ok := names[token]; ok {
			return resolved
		}
	}
	return token
}

func resolveValue(token string, values map[string]*dbstore.AttributeValue, names map[string]string) *dbstore.AttributeValue {
	if values != nil {
		if v, ok := values[token]; ok {
			return v
		}
	}
	return nil
}

func resolveValueWithIfNotExists(token string, values map[string]*dbstore.AttributeValue, names map[string]string, attrs map[string]*dbstore.AttributeValue) (*dbstore.AttributeValue, error) {
	if result := evaluateArithmeticExpression(token, values, names, attrs); result != nil {
		return result, nil
	}

	if strings.HasPrefix(token, "if_not_exists(") {
		closeParen := strings.Index(token, ")")
		if closeParen != -1 {
			inner := token[14:closeParen]
			parts := strings.SplitN(inner, ",", 2)
			if len(parts) == 2 {
				path := strings.TrimSpace(parts[0])
				defaultVal := strings.TrimSpace(parts[1])
				attrName := resolveName(path, names)
				if existing, ok := attrs[attrName]; ok && existing != nil {
					return existing, nil
				}
				return resolveValueWithIfNotExists(defaultVal, values, names, attrs)
			}
		}
	}

	if strings.HasPrefix(token, "list_append(") {
		closeParen := findMatchingCloseParen(token, 11)
		if closeParen != -1 {
			inner := token[12:closeParen]
			parts := splitListAppendArgs(inner)
			if len(parts) == 2 {
				list1 := resolveValueOrPath(strings.TrimSpace(parts[0]), values, names, attrs)
				list2 := resolveValueOrPath(strings.TrimSpace(parts[1]), values, names, attrs)
				if list1 != nil && list1.L != nil && list2 != nil && list2.L != nil {
					combined := make([]*dbstore.AttributeValue, 0, len(list1.L)+len(list2.L))
					combined = append(combined, list1.L...)
					combined = append(combined, list2.L...)
					return dbstore.ListValue(combined), nil
				}
				if list1 == nil && list2 != nil && list2.L != nil {
					return list2, nil
				}
				if list1 != nil && list1.L != nil {
					return list1, nil
				}
				if (list1 != nil && list1.L == nil) || (list2 != nil && list2.L == nil) {
					return nil, fmt.Errorf("TYPE_MISMATCH: Type mismatch for attribute to update")
				}
			}
		}
	}

	return resolveValue(token, values, names), nil
}

func evaluateArithmeticExpression(token string, values map[string]*dbstore.AttributeValue, names map[string]string, attrs map[string]*dbstore.AttributeValue) *dbstore.AttributeValue {
	token = strings.TrimSpace(token)

	parts := splitArithmeticExpression(token)
	if len(parts) == 1 {
		return nil
	}

	var result *big.Rat
	for i := 0; i < len(parts); i++ {
		part := strings.TrimSpace(parts[i])
		if part == "" {
			continue
		}

		var op string
		var operand string
		if i == 0 {
			op = "+"
			operand = part
		} else if part == "+" || part == "-" || part == "*" {
			op = part
			i++
			if i < len(parts) {
				operand = strings.TrimSpace(parts[i])
			}
		} else {
			operand = part
			op = "+"
		}

		if operand == "" {
			continue
		}

		var val *big.Rat
		if valAttr := resolveValueOrPath(operand, values, names, attrs); valAttr != nil && valAttr.N != nil {
			if r, ok := new(big.Rat).SetString(*valAttr.N); ok {
				val = r
			}
		}

		if val == nil {
			return nil
		}

		if result == nil {
			result = val
		} else {
			switch op {
			case "+":
				result = new(big.Rat).Add(result, val)
			case "-":
				result = new(big.Rat).Sub(result, val)
			case "*":
				result = new(big.Rat).Mul(result, val)
			}
		}
	}

	if result != nil {
		str := normalizeNumber(result)
		return dbstore.NumberValue(str)
	}
	return nil
}

func splitArithmeticExpression(expr string) []string {
	var parts []string
	current := ""
	depth := 0
	inString := false

	for i := 0; i < len(expr); i++ {
		c := expr[i]
		switch c {
		case '"':
			inString = !inString
			current += string(c)
		case '(', '[', '{':
			depth++
			current += string(c)
		case ')', ']', '}':
			depth--
			current += string(c)
		case ' ', '\t':
			if inString || depth > 0 {
				current += string(c)
			} else if current != "" {
				parts = append(parts, current)
				current = ""
			}
		case '+', '-', '*':
			if depth == 0 && !inString {
				if current != "" {
					parts = append(parts, current)
					current = ""
				}
				parts = append(parts, string(c))
			} else {
				current += string(c)
			}
		default:
			current += string(c)
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}

func findMatchingCloseParen(s string, openPos int) int {
	depth := 0
	for i := openPos; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func splitListAppendArgs(s string) []string {
	depth := 0
	for i, c := range s {
		if c == '(' {
			depth++
		} else if c == ')' {
			depth--
		} else if c == ',' && depth == 0 {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}

func resolveValueOrPath(token string, values map[string]*dbstore.AttributeValue, names map[string]string, attrs map[string]*dbstore.AttributeValue) *dbstore.AttributeValue {
	if strings.HasPrefix(token, "if_not_exists(") {
		val, _ := resolveValueWithIfNotExists(token, values, names, attrs)
		return val
	}
	if strings.HasPrefix(token, ":") {
		return resolveValue(token, values, names)
	}
	attrName := resolveName(token, names)
	if existing, ok := attrs[attrName]; ok {
		return existing
	}
	return nil
}
