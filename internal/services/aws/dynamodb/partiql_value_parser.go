package dynamodb

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	"vorpalstacks/pkg/sqlparser"
)

func objectLiteralToAttributes(obj *sqlparser.ObjectLiteral) map[string]*dbstore.AttributeValue {
	result := make(map[string]*dbstore.AttributeValue)
	for _, prop := range obj.Properties {
		key := extractPropertyKey(prop)
		if key == "" {
			continue
		}
		result[key] = exprToAttributeValue(prop.Value)
	}
	return result
}

func extractPropertyKey(prop *sqlparser.ObjectProperty) string {
	switch k := prop.Key.(type) {
	case *sqlparser.SQLVal:
		return string(k.Val)
	case *sqlparser.ColName:
		return k.Name.String()
	}
	return ""
}

func exprToAttributeValue(expr sqlparser.Expr) *dbstore.AttributeValue {
	switch e := expr.(type) {
	case *sqlparser.SQLVal:
		switch e.Type {
		case sqlparser.StrVal:
			s := string(e.Val)
			return &dbstore.AttributeValue{S: &s}
		case sqlparser.IntVal:
			s := string(e.Val)
			return &dbstore.AttributeValue{N: &s}
		case sqlparser.FloatVal:
			s := string(e.Val)
			return &dbstore.AttributeValue{N: &s}
		case sqlparser.ValArg:
			s := string(e.Val)
			return &dbstore.AttributeValue{S: &s}
		}
	case *sqlparser.NullVal:
		return &dbstore.AttributeValue{NULL: ptrBool(true)}
	case *sqlparser.BoolVal:
		b := bool(*e)
		return &dbstore.AttributeValue{BOOL: &b}
	case *sqlparser.ObjectLiteral:
		m := objectLiteralToAttributes(e)
		return &dbstore.AttributeValue{M: m}
	case sqlparser.ValTuple:
		l := tupleToAttributeList(e)
		return &dbstore.AttributeValue{L: l}
	}
	return &dbstore.AttributeValue{NULL: ptrBool(true)}
}

func tupleToAttributeList(tuple sqlparser.ValTuple) []*dbstore.AttributeValue {
	var result []*dbstore.AttributeValue
	for _, item := range tuple {
		result = append(result, exprToAttributeValue(item))
	}
	return result
}

func parsePartiQLValue(valueStr string) map[string]*dbstore.AttributeValue {
	result := make(map[string]*dbstore.AttributeValue)

	valueStr = strings.TrimSpace(valueStr)
	if !strings.HasPrefix(valueStr, "{") || !strings.HasSuffix(valueStr, "}") {
		return result
	}
	valueStr = strings.Trim(valueStr, "{}")

	pairs := splitPartiQLPairs(valueStr)
	for _, pair := range pairs {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(strings.Trim(parts[0], `"'`))
		value := strings.TrimSpace(parts[1])

		result[key] = parsePartiQLValuePart(value)
	}

	return result
}

func splitPartiQLPairs(s string) []string {
	var pairs []string
	depth := 0
	current := ""
	for _, ch := range s {
		switch ch {
		case '{', '[', '(':
			depth++
			current += string(ch)
		case '}', ']', ')':
			depth--
			current += string(ch)
		case ',':
			if depth == 0 {
				pairs = append(pairs, strings.TrimSpace(current))
				current = ""
			} else {
				current += string(ch)
			}
		default:
			current += string(ch)
		}
	}
	if strings.TrimSpace(current) != "" {
		pairs = append(pairs, strings.TrimSpace(current))
	}
	return pairs
}

func parsePartiQLValuePart(value string) *dbstore.AttributeValue {
	value = strings.TrimSpace(value)

	if value == "null" {
		return &dbstore.AttributeValue{NULL: ptrBool(true)}
	}

	if value == "true" || value == "false" {
		b := value == "true"
		return &dbstore.AttributeValue{BOOL: &b}
	}

	if strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`) {
		s := strings.Trim(value, `'`)
		return &dbstore.AttributeValue{S: &s}
	}

	if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
		s := strings.Trim(value, `"`)
		return &dbstore.AttributeValue{S: &s}
	}

	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return &dbstore.AttributeValue{N: &value}
	}

	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return &dbstore.AttributeValue{N: &value}
	}

	if strings.HasPrefix(value, "{") {
		nested := parsePartiQLValue(value)
		return &dbstore.AttributeValue{M: nested}
	}

	if strings.HasPrefix(value, "[") {
		list := parsePartiQLList(value)
		return &dbstore.AttributeValue{L: list}
	}

	return &dbstore.AttributeValue{S: &value}
}

func parsePartiQLList(value string) []*dbstore.AttributeValue {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, "[") || !strings.HasSuffix(value, "]") {
		return nil
	}
	value = strings.Trim(value, "[]")

	if strings.TrimSpace(value) == "" {
		return nil
	}

	var result []*dbstore.AttributeValue
	depth := 0
	current := ""
	for _, ch := range value {
		switch ch {
		case '{', '[', '(':
			depth++
			current += string(ch)
		case '}', ']', ')':
			depth--
			current += string(ch)
		case ',':
			if depth == 0 {
				result = append(result, parsePartiQLValuePart(strings.TrimSpace(current)))
				current = ""
			} else {
				current += string(ch)
			}
		default:
			current += string(ch)
		}
	}
	if strings.TrimSpace(current) != "" {
		result = append(result, parsePartiQLValuePart(strings.TrimSpace(current)))
	}

	return result
}

func paramToAttributeValue(param interface{}) *dbstore.AttributeValue {
	switch v := param.(type) {
	case map[string]interface{}:
		if s, ok := v["S"].(string); ok {
			return &dbstore.AttributeValue{S: &s}
		}
		if n, ok := v["N"].(string); ok {
			return &dbstore.AttributeValue{N: &n}
		}
		if b, ok := v["BOOL"].(bool); ok {
			return &dbstore.AttributeValue{BOOL: &b}
		}
		if null, ok := v["NULL"].(bool); ok && null {
			return &dbstore.AttributeValue{NULL: ptrBool(true)}
		}
	case string:
		return &dbstore.AttributeValue{S: &v}
	case float64:
		s := strconv.FormatFloat(v, 'f', -1, 64)
		return &dbstore.AttributeValue{N: &s}
	case int:
		s := strconv.FormatInt(int64(v), 10)
		return &dbstore.AttributeValue{N: &s}
	case int64:
		s := strconv.FormatInt(v, 10)
		return &dbstore.AttributeValue{N: &s}
	case bool:
		return &dbstore.AttributeValue{BOOL: &v}
	}
	s := fmt.Sprintf("%v", param)
	return &dbstore.AttributeValue{S: &s}
}

func compareNumberStrings(a, b string) int {
	numA, okA := new(big.Rat).SetString(a)
	numB, okB := new(big.Rat).SetString(b)
	if !okA || !okB {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	}
	return numA.Cmp(numB)
}

func numbersEqual(a, b string) bool {
	return compareNumberStrings(a, b) == 0
}
