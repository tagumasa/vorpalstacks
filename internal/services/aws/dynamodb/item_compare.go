package dynamodb

import (
	"bytes"
	"math/big"
	"strings"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func compareBetween(attr, low, high *dbstore.AttributeValue) bool {
	cLow, okLow := compareOrderedValues(attr, low)
	cHigh, okHigh := compareOrderedValues(attr, high)
	return okLow && okHigh && cLow >= 0 && cHigh <= 0
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

// isValidComparisonOperator returns true for recognised DynamoDB
// ConditionExpression comparison operators. Legacy ComparisonOperator
// API names (NE, LT, LE, GT, GE, !=) are intentionally excluded.
func isValidComparisonOperator(op string) bool {
	switch op {
	case "=", "<>", "<", "<=", ">", ">=":
		return true
	}
	return false
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
		c, ok := compareOrderedValues(attr, value)
		return ok && c < 0
	case "<=", "LE":
		c, ok := compareOrderedValues(attr, value)
		return ok && c <= 0
	case ">", "GT":
		c, ok := compareOrderedValues(attr, value)
		return ok && c > 0
	case ">=", "GE":
		c, ok := compareOrderedValues(attr, value)
		return ok && c >= 0
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

// compareOrderedValues applies DynamoDB ordering rules: only same-type
// S/N/B operand pairs are ordered (numbers numerically, binaries bytewise);
// every other pair — mismatched types, BOOL, sets, lists, maps, NULL — is
// incomparable and reports ok=false.
func compareOrderedValues(a, b *dbstore.AttributeValue) (int, bool) {
	if a.S != nil && b.S != nil {
		return strings.Compare(*a.S, *b.S), true
	}
	if a.N != nil && b.N != nil {
		return compareNumbers(a, b), true
	}
	if a.B != nil && b.B != nil {
		return bytes.Compare(a.B, b.B), true
	}
	return 0, false
}

// genericCompare orders two attributes for the key-condition engines,
// whose operands are schema-typed; incomparable pairs compare as equal.
func genericCompare(a, b *dbstore.AttributeValue) int {
	c, _ := compareOrderedValues(a, b)
	return c
}

func normalizeNumberString(n string) string {
	num, ok := new(big.Rat).SetString(n)
	if !ok {
		return n
	}
	return num.RatString()
}
