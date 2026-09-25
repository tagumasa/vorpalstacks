package dynamodb

import (
	"bytes"
	"math/big"
	"strings"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// filterByCondition applies a request's FilterExpression — parsed and
// validated before the read ran — to the scanned page.
func filterByCondition(items []*dbstore.Item, cond *compiledCondition) []*dbstore.Item {
	if cond == nil {
		return items
	}
	var result []*dbstore.Item
	for _, item := range items {
		if cond.matches(item) {
			result = append(result, item)
		}
	}
	return result
}

// compareAttributeValues evaluates one comparison under the expression
// grammar's six comparators (= <> < <= > >=). Legacy ComparisonOperator API
// names never reach this point: the legacy parameters are translated to the
// grammar's operators before evaluation, and the operator gates reject
// anything else.
func compareAttributeValues(attr *dbstore.AttributeValue, op string, value *dbstore.AttributeValue) bool {
	if attr == nil || value == nil {
		return false
	}

	switch op {
	case "=":
		return attributeValuesEqual(attr, value)
	case "<>":
		return !attributeValuesEqual(attr, value)
	case "<":
		c, ok := compareOrderedValues(attr, value)
		return ok && c < 0
	case "<=":
		c, ok := compareOrderedValues(attr, value)
		return ok && c <= 0
	case ">":
		c, ok := compareOrderedValues(attr, value)
		return ok && c > 0
	case ">=":
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
		aSet := make(map[string]bool)
		for _, elem := range a.BS {
			aSet[string(elem)] = true
		}
		for _, elem := range b.BS {
			if !aSet[string(elem)] {
				return false
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

// normalizeNumberString renders one number spelling in the canonical
// decimal DynamoDB form. Number-set membership keys and appended set
// members both use it, so two spellings of one value ("1" and "01",
// "1.5" and "1.50") are one set member, and a member stored through an
// ADD is always a valid DynamoDB number — big.Rat's own fraction form
// ("3/2" for 1.5) is not part of the number grammar and must never
// reach the stored item.
func normalizeNumberString(n string) string {
	num, ok := new(big.Rat).SetString(n)
	if !ok {
		return n
	}
	return normalizeNumber(num)
}
