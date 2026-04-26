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
	return objectLiteralToAttributesWithParams(obj, nil)
}

func objectLiteralToAttributesWithParams(obj *sqlparser.ObjectLiteral, params *partiQLParams) map[string]*dbstore.AttributeValue {
	result := make(map[string]*dbstore.AttributeValue)
	for _, prop := range obj.Properties {
		key := extractPropertyKey(prop)
		if key == "" {
			continue
		}
		result[key] = exprToAttributeValueWithParams(prop.Value, params)
	}
	return result
}

func exprToAttributeValueWithParams(expr sqlparser.Expr, params *partiQLParams) *dbstore.AttributeValue {
	if e, ok := expr.(*sqlparser.SQLVal); ok && e.Type == sqlparser.ValArg && params != nil {
		if strings.HasPrefix(string(e.Val), ":") {
			idxStr := strings.TrimPrefix(string(e.Val), ":v")
			if idx, err := strconv.Atoi(idxStr); err == nil && idx > 0 && idx <= len(params.Parameters) {
				return paramToAttributeValue(params.Parameters[idx-1])
			}
		}
	}
	return exprToAttributeValue(expr)
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
