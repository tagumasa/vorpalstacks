package dynamodb

import (
	"bytes"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

type sortKeyCondition struct {
	op       string
	attrName string
	value    *dbstore.AttributeValue
	valueTo  *dbstore.AttributeValue
}

// queryKeyConditionNotSupported is the ValidationException every malformed
// key condition answers with. The documented grammar: the partition key
// appears as an equality condition alone, the sort key condition uses one
// of =, <, <=, >, >=, BETWEEN, begins_with, the two combine with AND, and
// no non-key attribute may appear.
func queryKeyConditionNotSupported() error {
	return NewAPIError("com.amazon.coral.validate#ValidationException",
		"Query key condition not supported", http.StatusBadRequest)
}

func extractPrimaryKeyCondition(table *dbstore.Table, expr string, names map[string]string, values map[string]*dbstore.AttributeValue) (string, *dbstore.AttributeValue, *sortKeyCondition, error) {
	var hashKeyName, sortKeyName string
	for _, ks := range table.KeySchema {
		if ks.KeyType == dbstore.KeyTypeHash {
			hashKeyName = ks.AttributeName
		} else if ks.KeyType == dbstore.KeyTypeRange {
			sortKeyName = ks.AttributeName
		}
	}
	hashValue, hashAttr, sortCond, err := extractKeyConditionFromSchema(hashKeyName, sortKeyName, expr, names, values)
	if err != nil {
		return "", nil, nil, err
	}
	if err := validateKeyConditionValueTypes(table, hashKeyName, hashAttr, sortKeyName, sortCond); err != nil {
		return "", nil, nil, err
	}
	return hashValue, hashAttr, sortCond, nil
}

func extractIndexKeyCondition(table *dbstore.Table, indexName, expr string, names map[string]string, values map[string]*dbstore.AttributeValue) (string, *dbstore.AttributeValue, *sortKeyCondition, error) {
	var hashKeyName, sortKeyName string
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
	hashValue, hashAttr, sortCond, err := extractKeyConditionFromSchema(hashKeyName, sortKeyName, expr, names, values)
	if err != nil {
		return "", nil, nil, err
	}
	if err := validateKeyConditionValueTypes(table, hashKeyName, hashAttr, sortKeyName, sortCond); err != nil {
		return "", nil, nil, err
	}
	return hashValue, hashAttr, sortCond, nil
}

// conditionTypeMismatchError is the ValidationException a key condition
// whose resolved value type does not match the key's schema type answers
// with.
func conditionTypeMismatchError() error {
	return NewAPIError("com.amazon.coral.validate#ValidationException",
		"One or more parameter values were invalid: Condition parameter type does not match schema type", http.StatusBadRequest)
}

// validateKeyConditionValueTypes checks the resolved condition values
// against the attribute definitions of the key attributes they address
// (the table's own schema, or the queried index's — both draw their types
// from the same table attribute definitions). A wrong-typed value is
// encoded or compared by its own type and would silently match nothing,
// so it is rejected instead. begins_with excludes a Number-typed sort key
// alone ("You cannot use this function with a sort key that is of type
// Number" — documented on the Query key-condition grammar page; the legacy
// KeyConditions face types both the operand and the target "String or
// Binary (not a Number or a set type)"), so String and Binary sort keys
// both carry it.
func validateKeyConditionValueTypes(table *dbstore.Table, hashName string, hashVal *dbstore.AttributeValue, sortName string, sortCond *sortKeyCondition) error {
	defs := attributeTypeDefinitions(table)
	check := func(name string, val *dbstore.AttributeValue) error {
		if name == "" || val == nil {
			return nil
		}
		expected, hasDef := defs[name]
		if !hasDef {
			return nil
		}
		if attributeValueType(val) != string(expected) {
			return conditionTypeMismatchError()
		}
		return nil
	}
	if err := check(hashName, hashVal); err != nil {
		return err
	}
	if sortCond == nil {
		return nil
	}
	if sortCond.op == "begins_with" {
		if expected, hasDef := defs[sortName]; hasDef && expected == dbstore.ScalarAttributeTypeN {
			return queryKeyConditionNotSupported()
		}
	}
	if err := check(sortName, sortCond.value); err != nil {
		return err
	}
	return check(sortName, sortCond.valueTo)
}

// caseInsensitiveMatchAt reports whether the runes starting at expr[pos:]
// match opUpper case-insensitively, and how many input bytes the match
// consumed. opUpper is compared rune by rune against each input rune's
// simple uppercase mapping, so the two strings' byte indices never have to
// agree — a precomputed strings.ToUpper copy of expr would shift its bytes
// wherever a rune's uppercase form changes byte length (e.g. 'ı' → 'I').
//
// Key-plane machinery: the KeyConditionExpression decomposition still
// splits its expression string; the condition plane parses to a typed
// tree instead (condition_ast.go).
func caseInsensitiveMatchAt(expr, opUpper string, pos int) (bool, int) {
	i, j := pos, 0
	for j < len(opUpper) {
		if i >= len(expr) {
			return false, 0
		}
		r, size := utf8.DecodeRuneInString(expr[i:])
		opRune, opSize := utf8.DecodeRuneInString(opUpper[j:])
		if unicode.ToUpper(r) != opRune {
			return false, 0
		}
		i += size
		j += opSize
	}
	return true, i - pos
}

// splitByLogicalOp splits an expression at its top-level occurrences of
// the given logical operator (case-insensitive, quote- and depth-aware).
// When splitting by AND, the operator occurring inside a BETWEEN's bounds
// is the range separator, not a conjunction. Key-plane machinery for the
// KeyConditionExpression decomposition, as above.
func splitByLogicalOp(expr string, op string) []string {
	depth := 0
	inString := false
	var parts []string
	current := ""
	opUpper := strings.ToUpper(op)
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

			if depth == 0 {
				if matched, matchLen := caseInsensitiveMatchAt(expr, opUpper, i); matched {
					// When splitting by " AND ", skip ANDs that are part
					// of a BETWEEN expression. The accumulated text ending
					// with "<token> BETWEEN <token>" means this AND is the
					// BETWEEN separator, not a logical conjunction.
					if opUpper == " AND " {
						fields := strings.Fields(strings.TrimSpace(current))
						if len(fields) >= 2 && strings.EqualFold(fields[len(fields)-2], "BETWEEN") {
							current += expr[i : i+1]
							i++
							continue
						}
					}
					if trimmed := strings.TrimSpace(current); trimmed != "" {
						parts = append(parts, trimmed)
					}
					current = ""
					i += matchLen
					continue
				}
			}
		}

		// ch is a raw byte compared against ASCII delimiters only; the
		// accumulation must carry the byte itself, not string(ch), which
		// would re-encode the byte value as a code point and mangle any
		// multi-byte UTF-8 rune.
		current += expr[i : i+1]
		i++
	}

	if trimmed := strings.TrimSpace(current); trimmed != "" {
		parts = append(parts, trimmed)
	}

	return parts
}

// extractKeyConditionFromSchema parses a KeyConditionExpression against one
// key schema (the table's, or the queried index's). The expression is split
// into its top-level AND conjuncts; every conjunct must be exactly one
// condition — a partition-key equality or a sort-key condition with a
// documented operator — and anything else (an OR, a duplicate condition, an
// unsupported operator, a condition whose value placeholder does not
// resolve, a non-key attribute, begins_with on the partition key) is the
// documented "Query key condition not supported" rejection. A conjunct the
// schema cannot name leaves the hash value empty; the caller answers that
// shape with the missed-key-schema rejection.
func extractKeyConditionFromSchema(hashKeyName, sortKeyName, expr string, names map[string]string, values map[string]*dbstore.AttributeValue) (string, *dbstore.AttributeValue, *sortKeyCondition, error) {
	// A top-level OR is not a key condition: the grammar combines the
	// partition and sort conditions with AND alone.
	if orParts := splitByLogicalOp(expr, " OR "); len(orParts) > 1 {
		return "", nil, nil, queryKeyConditionNotSupported()
	}

	var hashKeyValue string
	var hashKeyAttr *dbstore.AttributeValue
	var sortCond *sortKeyCondition
	sawHash, sawSort := false, false

	for _, conjunct := range splitByLogicalOp(expr, " AND ") {
		tokens := tokenizeExpression(conjunct)
		if len(tokens) == 0 {
			return "", nil, nil, queryKeyConditionNotSupported()
		}

		// Function form: begins_with(key-attribute, :value) glues into a
		// single token. The function name is case-sensitive, so the match
		// is exact.
		if strings.HasPrefix(tokens[0], "begins_with(") {
			funcToken := tokens[0][len("begins_with("):]
			commaIdx := strings.Index(funcToken, ",")
			if commaIdx == -1 {
				return "", nil, nil, queryKeyConditionNotSupported()
			}
			funcAttrName := resolveName(strings.TrimSpace(funcToken[:commaIdx]), names)
			if funcAttrName == hashKeyName {
				// The partition key admits equality alone.
				return "", nil, nil, queryKeyConditionNotSupported()
			}
			if funcAttrName != sortKeyName {
				return "", nil, nil, queryKeyConditionNotSupported()
			}
			if sawSort {
				return "", nil, nil, queryKeyConditionNotSupported()
			}
			valueToken := strings.TrimSuffix(strings.TrimSpace(funcToken[commaIdx+1:]), ")")
			val := resolveValue(valueToken, values, names)
			if val == nil {
				return "", nil, nil, queryKeyConditionNotSupported()
			}
			sortCond = &sortKeyCondition{op: "begins_with", attrName: sortKeyName, value: val}
			sawSort = true
			continue
		}

		attrName := resolveName(tokens[0], names)
		if attrName == hashKeyName {
			if len(tokens) != 3 || tokens[1] != "=" {
				return "", nil, nil, queryKeyConditionNotSupported()
			}
			if sawHash {
				return "", nil, nil, queryKeyConditionNotSupported()
			}
			val := resolveValue(tokens[2], values, names)
			if val == nil {
				return "", nil, nil, queryKeyConditionNotSupported()
			}
			// The hash value is the storage partition prefix; it must be
			// the store's key encoding or the prefix scan never matches
			// stored keys of non-string types.
			hashKeyValue = dbstore.EncodeKeyValue(val)
			hashKeyAttr = val
			sawHash = true
			continue
		}

		if attrName == sortKeyName {
			if len(tokens) < 3 {
				return "", nil, nil, queryKeyConditionNotSupported()
			}
			op := strings.ToLower(tokens[1])
			switch op {
			case "=", "<", "<=", ">", ">=":
				// begins_with exists in the case-sensitive function form
				// alone; the documented comparator set carries no infix
				// operator for it.
				if len(tokens) != 3 || sawSort {
					return "", nil, nil, queryKeyConditionNotSupported()
				}
				val := resolveValue(tokens[2], values, names)
				if val == nil {
					return "", nil, nil, queryKeyConditionNotSupported()
				}
				sortCond = &sortKeyCondition{op: op, attrName: sortKeyName, value: val}
				sawSort = true
			case "between":
				if len(tokens) != 5 || !strings.EqualFold(tokens[3], "AND") || sawSort {
					return "", nil, nil, queryKeyConditionNotSupported()
				}
				valFrom := resolveValue(tokens[2], values, names)
				valTo := resolveValue(tokens[4], values, names)
				if valFrom == nil || valTo == nil {
					return "", nil, nil, queryKeyConditionNotSupported()
				}
				sortCond = &sortKeyCondition{op: op, attrName: sortKeyName, value: valFrom, valueTo: valTo}
				sawSort = true
			default:
				// Unsupported operators (contains, IN, <>, attribute
				// functions, ...) on the sort key are rejected, not stored
				// verbatim to filter every item out.
				return "", nil, nil, queryKeyConditionNotSupported()
			}
			continue
		}

		// No other attribute may appear in a key condition.
		return "", nil, nil, queryKeyConditionNotSupported()
	}

	return hashKeyValue, hashKeyAttr, sortCond, nil
}

func isGSI(table *dbstore.Table, indexName string) bool {
	for _, gsi := range table.GlobalSecondaryIndexes {
		if gsi.IndexName == indexName {
			return true
		}
	}
	return false
}

// sortKeyConditionMatches reports whether an item's sort-key attribute
// satisfies the sort-key condition of a Query. A nil or value-less
// condition matches everything; an item lacking the sort-key attribute
// never matches.
func sortKeyConditionMatches(item *dbstore.Item, cond *sortKeyCondition) bool {
	if cond == nil || cond.value == nil || cond.attrName == "" {
		return true
	}
	av, ok := item.Attributes[cond.attrName]
	if !ok || av == nil {
		return false
	}
	return compareWithCondition(av, cond)
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
		if attr.B != nil && cond.value.B != nil {
			return bytes.HasPrefix(attr.B, cond.value.B)
		}
	case "between":
		if cond.valueTo == nil {
			return false
		}
		return genericCompare(attr, cond.value) >= 0 && genericCompare(attr, cond.valueTo) <= 0
	}
	return false
}
