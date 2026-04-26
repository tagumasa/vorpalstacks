package dynamodb

import (
	"bytes"
	"math/big"
	"strings"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func filterByKeyCondition(items []*dbstore.Item, expr string, names map[string]string, values map[string]*dbstore.AttributeValue) []*dbstore.Item {
	var result []*dbstore.Item
	for _, item := range items {
		if evaluateKeyCondition(item, expr, names, values) {
			result = append(result, item)
		}
	}
	return result
}

func evaluateKeyCondition(item *dbstore.Item, expr string, names map[string]string, values map[string]*dbstore.AttributeValue) bool {
	tokens := tokenizeExpression(expr)
	if len(tokens) < 3 {
		return true
	}

	attrName := resolveName(tokens[0], names)
	op := tokens[1]

	attr, ok := item.Attributes[attrName]
	if !ok {
		return false
	}

	if strings.ToUpper(op) == "BETWEEN" && len(tokens) >= 5 {
		val1 := resolveValue(tokens[2], values, names)
		val2 := resolveValue(tokens[4], values, names)
		return compareBetween(attr, val1, val2)
	}

	if strings.ToUpper(op) == "IN" && len(tokens) >= 4 {
		var inValues []*dbstore.AttributeValue
		for i := 2; i < len(tokens); i++ {
			if tokens[i] != "(" && tokens[i] != ")" && tokens[i] != "," {
				inValues = append(inValues, resolveValue(tokens[i], values, names))
			}
		}
		return compareIn(attr, inValues)
	}

	value := resolveValue(tokens[2], values, names)
	return compareAttributeValues(attr, op, value)
}

func compareBetween(attr, low, high *dbstore.AttributeValue) bool {
	return genericCompare(attr, low) >= 0 && genericCompare(attr, high) <= 0
}

func compareIn(attr *dbstore.AttributeValue, inValues []*dbstore.AttributeValue) bool {
	for _, v := range inValues {
		if attributeValuesEqual(attr, v) {
			return true
		}
	}
	return false
}

func filterByExpression(items []*dbstore.Item, expr string, names map[string]string, values map[string]*dbstore.AttributeValue) []*dbstore.Item {
	var result []*dbstore.Item
	for _, item := range items {
		if evaluateFilterExpression(item, expr, names, values) {
			result = append(result, item)
		}
	}
	return result
}

func evaluateFilterExpression(item *dbstore.Item, expr string, names map[string]string, values map[string]*dbstore.AttributeValue) bool {
	if expr == "" {
		return true
	}

	result, err := evaluateConditionExpr(item, expr, names, values)
	if err != nil {
		return false
	}
	return result
}

func compareAttributeValues(attr *dbstore.AttributeValue, op string, value *dbstore.AttributeValue) bool {
	if attr == nil || value == nil {
		return false
	}

	switch op {
	case "=":
		return attributeValuesEqual(attr, value)
	case "<>", "!=", "NE":
		return !attributeValuesEqual(attr, value)
	case "<", "LT":
		return genericCompare(attr, value) < 0
	case "<=", "LE":
		return genericCompare(attr, value) <= 0
	case ">", "GT":
		return genericCompare(attr, value) > 0
	case ">=", "GE":
		return genericCompare(attr, value) >= 0
	}

	return false
}

func attributeValuesEqual(a, b *dbstore.AttributeValue) bool {
	if a.NULL != nil && b.NULL != nil {
		return *a.NULL == *b.NULL
	}
	if a.NULL != nil || b.NULL != nil {
		return false
	}

	if a.S != nil && b.S != nil {
		return *a.S == *b.S
	}
	if a.S != nil || b.S != nil {
		return false
	}

	if a.N != nil && b.N != nil {
		return compareNumbers(a, b) == 0
	}
	if a.N != nil || b.N != nil {
		return false
	}

	if a.B != nil && b.B != nil {
		if len(a.B) != len(b.B) {
			return false
		}
		for i := range a.B {
			if a.B[i] != b.B[i] {
				return false
			}
		}
		return true
	}
	if a.B != nil || b.B != nil {
		return false
	}

	if a.BOOL != nil && b.BOOL != nil {
		return *a.BOOL == *b.BOOL
	}
	if a.BOOL != nil || b.BOOL != nil {
		return false
	}

	if a.SS != nil && b.SS != nil {
		if len(a.SS) != len(b.SS) {
			return false
		}
		aSet := make(map[string]bool)
		for _, s := range a.SS {
			aSet[s] = true
		}
		for _, s := range b.SS {
			if !aSet[s] {
				return false
			}
		}
		return true
	}
	if a.SS != nil || b.SS != nil {
		return false
	}

	if a.NS != nil && b.NS != nil {
		if len(a.NS) != len(b.NS) {
			return false
		}
		aSet := make(map[string]bool)
		for _, n := range a.NS {
			normalized := normalizeNumberString(n)
			aSet[normalized] = true
		}
		for _, n := range b.NS {
			normalized := normalizeNumberString(n)
			if !aSet[normalized] {
				return false
			}
		}
		return true
	}
	if a.NS != nil || b.NS != nil {
		return false
	}

	if a.BS != nil && b.BS != nil {
		if len(a.BS) != len(b.BS) {
			return false
		}
		for i := range a.BS {
			if len(a.BS[i]) != len(b.BS[i]) {
				return false
			}
			for j := range a.BS[i] {
				if a.BS[i][j] != b.BS[i][j] {
					return false
				}
			}
		}
		return true
	}
	if a.BS != nil || b.BS != nil {
		return false
	}

	if a.M != nil && b.M != nil {
		if len(a.M) != len(b.M) {
			return false
		}
		for k, v := range a.M {
			bv, ok := b.M[k]
			if !ok || !attributeValuesEqual(v, bv) {
				return false
			}
		}
		return true
	}
	if a.M != nil || b.M != nil {
		return false
	}

	if a.L != nil && b.L != nil {
		if len(a.L) != len(b.L) {
			return false
		}
		for i := range a.L {
			if !attributeValuesEqual(a.L[i], b.L[i]) {
				return false
			}
		}
		return true
	}
	if a.L != nil || b.L != nil {
		return false
	}

	return false
}

func compareNumbers(a, b *dbstore.AttributeValue) int {
	if a.N == nil || b.N == nil {
		return 0
	}
	numA, okA := new(big.Rat).SetString(*a.N)
	numB, okB := new(big.Rat).SetString(*b.N)
	if !okA || !okB {
		return 0
	}
	return numA.Cmp(numB)
}

func genericCompare(a, b *dbstore.AttributeValue) int {
	if a.S != nil && b.S != nil {
		return strings.Compare(*a.S, *b.S)
	}
	if a.N != nil && b.N != nil {
		return compareNumbers(a, b)
	}
	if a.B != nil && b.B != nil {
		return bytes.Compare(a.B, b.B)
	}
	return 0
}

func normalizeNumberString(n string) string {
	num, ok := new(big.Rat).SetString(n)
	if !ok {
		return n
	}
	return num.RatString()
}
