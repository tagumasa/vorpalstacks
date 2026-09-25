package dynamodb

import (
	"fmt"
	"math/big"
	"net/http"
	"sort"
	"strings"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func applyUpdateExpression(attrs map[string]*dbstore.AttributeValue, expr string, names map[string]string, values map[string]*dbstore.AttributeValue) error {
	_, err := applyUpdateExpressionWithTracking(attrs, expr, names, values)
	return err
}

// extractUpdatedPaths analyses an UpdateExpression and returns the top-level
// attribute names that would be modified, without executing the update.
func extractUpdatedPaths(expr string, names map[string]string) ([]string, error) {
	tokens := tokenizeUpdateExpression(expr)
	var paths []string
	i := 0
	for i < len(tokens) {
		action := tokens[i]
		switch action {
		case "SET":
			i++
			for i < len(tokens) && tokens[i] != "REMOVE" && tokens[i] != "ADD" && tokens[i] != "DELETE" {
				paths = append(paths, topLevelSegment(tokens[i], names))
				j := i + 3
				for j < len(tokens) && (tokens[j] == "+" || tokens[j] == "-" || tokens[j] == "*") && j+1 < len(tokens) {
					j += 2
				}
				i = j
				next, sepErr := requireActionSeparator(tokens, i, "SET")
				if sepErr != nil {
					return nil, sepErr
				}
				i = next
			}
		case "REMOVE":
			i++
			for i < len(tokens) && tokens[i] != "SET" && tokens[i] != "ADD" && tokens[i] != "DELETE" {
				paths = append(paths, topLevelSegment(tokens[i], names))
				i++
				next, sepErr := requireActionSeparator(tokens, i, "REMOVE")
				if sepErr != nil {
					return nil, sepErr
				}
				i = next
			}
		case "ADD":
			i++
			for i < len(tokens) && tokens[i] != "SET" && tokens[i] != "REMOVE" && tokens[i] != "DELETE" {
				// The applier rejects document paths under ADD, so this
				// reduction is the identity on accepted expressions; it
				// keeps the key-attribute guard's input shape uniform
				// across clause types.
				paths = append(paths, topLevelSegment(tokens[i], names))
				i += 2
				next, sepErr := requireActionSeparator(tokens, i, "ADD")
				if sepErr != nil {
					return nil, sepErr
				}
				i = next
			}
		case "DELETE":
			i++
			for i < len(tokens) && tokens[i] != "SET" && tokens[i] != "REMOVE" && tokens[i] != "ADD" {
				// Same uniform-shape rule as ADD.
				paths = append(paths, topLevelSegment(tokens[i], names))
				i += 2
				next, sepErr := requireActionSeparator(tokens, i, "DELETE")
				if sepErr != nil {
					return nil, sepErr
				}
				i = next
			}
		default:
			return nil, ErrInvalidParameter
		}
	}
	return paths, nil
}

// applyUpdateExpressionWithTracking applies an UpdateExpression to attrs
// and reports the top-level attributes its actions touched. Every action
// evaluates against the item as it was before the expression — the
// documented basis ("DynamoDB evaluates every action against the item's
// attribute values as they were before the update") — so the walk's reads
// (SET operands, ADD and DELETE bases) resolve against a deep copy taken
// at entry, while the writes apply to attrs itself in expression order.
func applyUpdateExpressionWithTracking(attrs map[string]*dbstore.AttributeValue, expr string, names map[string]string, values map[string]*dbstore.AttributeValue) ([]string, error) {
	before := copyAttributes(attrs)
	var updatedAttrs []string
	tokens := tokenizeUpdateExpression(expr)
	i := 0
	for i < len(tokens) {
		action := tokens[i]
		switch action {
		case "SET":
			i++
			// Assignments whose target path ends in a list index are held
			// until the end of the clause and applied in ascending element
			// number — the documented ordering when one SET operation adds
			// multiple list elements, composed with the append-at-end rule
			// for out-of-range indices. The operand still evaluates at
			// encounter, exactly as an immediately-applied assignment's.
			var listAssignments []listElementAssignment
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
				pathParts, nameErr := resolveDocPathParts(tokens[i], names)
				if nameErr != nil {
					return nil, nameErr
				}
				if len(pathParts) == 0 {
					return nil, ErrInvalidParameter
				}
				exprTokens := []string{tokens[i+2]}
				j := i + 3
				for j+1 < len(tokens) {
					if tokens[j] == "*" {
						// Arithmetic in a SET action supports only the plus
						// and minus operators.
						return nil, ErrInvalidParameter
					}
					if tokens[j] != "+" && tokens[j] != "-" {
						break
					}
					exprTokens = append(exprTokens, tokens[j], tokens[j+1])
					j += 2
				}
				operand, opErr := parseUpdateOperand(strings.Join(exprTokens, " "), values, names)
				if opErr != nil {
					return nil, opErr
				}
				value, valErr := operand.eval(before)
				if valErr != nil {
					return nil, valErr
				}
				if value == nil {
					// Every operand of a SET action must resolve to a
					// value; a reference to an attribute that does not
					// exist is a validation error.
					return nil, ErrInvalidParameter
				}
				if last := pathParts[len(pathParts)-1]; last.isIndex {
					listAssignments = append(listAssignments, listElementAssignment{path: pathParts, value: value})
				} else if err := setNestedValue(attrs, pathParts, value); err != nil {
					return nil, err
				}
				updatedAttrs = append(updatedAttrs, pathParts[0].name)
				i = j
				next, sepErr := requireActionSeparator(tokens, i, "SET")
				if sepErr != nil {
					return nil, sepErr
				}
				i = next
			}
			sort.SliceStable(listAssignments, func(a, b int) bool {
				return listAssignments[a].path[len(listAssignments[a].path)-1].index <
					listAssignments[b].path[len(listAssignments[b].path)-1].index
			})
			for _, la := range listAssignments {
				if err := setNestedValue(attrs, la.path, la.value); err != nil {
					return nil, err
				}
			}
		case "REMOVE":
			i++
			for i < len(tokens) && tokens[i] != "SET" && tokens[i] != "ADD" && tokens[i] != "DELETE" {
				if i >= len(tokens) {
					return nil, ErrInvalidParameter
				}
				pathParts, nameErr := resolveDocPathParts(tokens[i], names)
				if nameErr != nil {
					return nil, nameErr
				}
				if len(pathParts) == 0 {
					return nil, ErrInvalidParameter
				}
				if rmErr := removeNestedValue(attrs, pathParts); rmErr != nil {
					return nil, rmErr
				}
				updatedAttrs = append(updatedAttrs, pathParts[0].name)
				i++
				next, sepErr := requireActionSeparator(tokens, i, "REMOVE")
				if sepErr != nil {
					return nil, sepErr
				}
				i = next
			}
		case "ADD":
			i++
			for i < len(tokens) && tokens[i] != "SET" && tokens[i] != "REMOVE" && tokens[i] != "DELETE" {
				if i+1 >= len(tokens) {
					return nil, ErrInvalidParameter
				}
				if isDocumentPathToken(tokens[i]) {
					return nil, topLevelActionPathError("ADD", tokens[i])
				}
				attrName, nameErr := resolveNameStrict(tokens[i], names)
				if nameErr != nil {
					return nil, nameErr
				}
				value := resolveValue(tokens[i+1], values, names)
				if value == nil {
					// An undefined value placeholder is a validation error.
					return nil, ErrInvalidParameter
				}
				if changed, err := applyAddAction(attrs, before, attrName, value); err != nil {
					return nil, err
				} else if changed {
					updatedAttrs = append(updatedAttrs, attrName)
				}
				i += 2
				next, sepErr := requireActionSeparator(tokens, i, "ADD")
				if sepErr != nil {
					return nil, sepErr
				}
				i = next
			}
		case "DELETE":
			i++
			for i < len(tokens) && tokens[i] != "SET" && tokens[i] != "REMOVE" && tokens[i] != "ADD" {
				if i+1 >= len(tokens) {
					return nil, ErrInvalidParameter
				}
				if isDocumentPathToken(tokens[i]) {
					return nil, topLevelActionPathError("DELETE", tokens[i])
				}
				attrName, nameErr := resolveNameStrict(tokens[i], names)
				if nameErr != nil {
					return nil, nameErr
				}
				value := resolveValue(tokens[i+1], values, names)
				if value == nil {
					// An undefined value placeholder is a validation error.
					return nil, ErrInvalidParameter
				}
				changed, delErr := applyDeleteAction(attrs, before, attrName, value)
				if delErr != nil {
					return nil, delErr
				}
				if changed {
					updatedAttrs = append(updatedAttrs, attrName)
				}
				i += 2
				next, sepErr := requireActionSeparator(tokens, i, "DELETE")
				if sepErr != nil {
					return nil, sepErr
				}
				i = next
			}
		default:
			return nil, ErrInvalidParameter
		}
	}
	return updatedAttrs, nil
}

// requireActionSeparator enforces the update-expression grammar between
// the actions of one clause: a clause's actions are separated by commas,
// so after one action the next token is the separator (consumed — and the
// separator itself must be followed by an action of the same clause, so a
// trailing comma, or one directly before a clause keyword, is a syntax
// error), another clause's keyword (this clause ends), or the end of the
// expression. Anything else is a syntax error — including the current
// clause's own keyword, whose second occurrence the grammar forbids (each
// action keyword may appear only once).
func requireActionSeparator(tokens []string, i int, clause string) (int, error) {
	if i >= len(tokens) {
		return i, nil
	}
	if tokens[i] == "," {
		next := i + 1
		if next >= len(tokens) || isUpdateClauseKeyword(tokens[next]) {
			return i, ErrInvalidParameter
		}
		return next, nil
	}
	if tokens[i] != clause && isUpdateClauseKeyword(tokens[i]) {
		return i, nil
	}
	return i, ErrInvalidParameter
}

func isUpdateClauseKeyword(token string) bool {
	switch token {
	case "SET", "REMOVE", "ADD", "DELETE":
		return true
	default:
		return false
	}
}

// listElementAssignment holds a SET assignment whose target path ends in a
// list index until its clause completes; the clause then applies the held
// assignments in ascending element number.
type listElementAssignment struct {
	path  []docPathPart
	value *dbstore.AttributeValue
}

// isDocumentPathToken reports whether a raw UpdateExpression path token
// addresses below the top level: a '.' or '[' separator. Attribute names
// containing those characters are only expressible through expression
// attribute names, so any separator in the raw token is a real path step.
func isDocumentPathToken(token string) bool {
	return strings.ContainsAny(token, ".[")
}

// topLevelActionPathError rejects a document path under the ADD and DELETE
// actions, which address top-level attributes only: indexing the item by
// the dotted string would create an attribute literally named "a.b".
func topLevelActionPathError(action, token string) error {
	return NewAPIError("com.amazon.coral.validate#ValidationException",
		fmt.Sprintf("One or more parameter values were invalid: %s action requires a top-level attribute path, but %q is a document path", action, token),
		http.StatusBadRequest)
}

// topLevelSegment returns the top-level attribute name of a raw
// document-path token: the first parsed segment with its expression
// attribute name resolved when defined. Resolution on the parsed segment
// keeps an alias standing for a name that itself contains '.' or '[' a
// single top-level attribute. An undefined alias stays literal and an
// unparseable token is returned as-is — this analysis pass must not
// pre-empt the apply pass's rejection of either.
func topLevelSegment(token string, names map[string]string) string {
	parts, err := resolveDocPathPartsLenient(token, names)
	if err != nil || len(parts) == 0 {
		return token
	}
	return parts[0].name
}

func addNumbers(a, b string) (string, error) {
	numA, okA := new(big.Rat).SetString(a)
	numB, okB := new(big.Rat).SetString(b)
	if !okA || !okB {
		return "", fmt.Errorf("invalid number operand: %q, %q", a, b)
	}
	numA.Add(numA, numB)
	return normalizeNumber(numA), nil
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

// combineAddValues returns the ADD action's combination of a stored value
// and its operand: numbers add numerically, same-kind sets union their
// members. A nil result with a nil error marks the no-change outcome of an
// empty set operand. A pair with no compatible type combination answers
// ErrTypeMismatch, the documented ValidationException wording every
// update plane answers.
func combineAddValues(existing, value *dbstore.AttributeValue) (*dbstore.AttributeValue, error) {
	if value.N != nil && existing.N != nil {
		result, addErr := addNumbers(*existing.N, *value.N)
		if addErr != nil {
			return nil, addErr
		}
		return &dbstore.AttributeValue{N: &result}, nil
	}
	if value.SS != nil && existing.SS != nil {
		return dbstore.StringSet(addStringSetMembers(existing.SS, value.SS, stringSetCanonical)), nil
	}
	if value.NS != nil && existing.NS != nil {
		return dbstore.NumberSet(addStringSetMembers(existing.NS, value.NS, normalizeNumberString)), nil
	}
	if value.BS != nil && existing.BS != nil {
		return dbstore.BinarySet(addBinarySetMembers(existing.BS, value.BS)), nil
	}
	return nil, ErrTypeMismatch
}

// applyAddAction applies one DynamoDB ADD action: a number adds
// numerically, a set unions its members, and an absent attribute is created
// from the operand. The stored value the action combines against is read
// from preUpdate — the item as it was before the whole expression, the
// documented evaluation basis — while the combined result lands in attrs;
// planes that apply one action at a time pass the same map for both. Every
// branch assigns a fresh value into attrs — the stored attribute is never
// mutated in place, so an old-image snapshot taken before the clause
// application cannot alias the updated item. A pair sharing no compatible
// type answers ErrTypeMismatch; changed reports whether the action
// affected the item, for the calling plane's tracking.
func applyAddAction(attrs, preUpdate map[string]*dbstore.AttributeValue, attrName string, value *dbstore.AttributeValue) (bool, error) {
	existing, exists := preUpdate[attrName]
	if !exists {
		attrs[attrName] = value
		return true, nil
	}
	combined, err := combineAddValues(existing, value)
	if err != nil {
		return false, err
	}
	// An empty set operand leaves the stored set unchanged; every other
	// combination the core answers is a change.
	changed := !(value.SS != nil && len(value.SS) == 0) &&
		!(value.NS != nil && len(value.NS) == 0) &&
		!(value.BS != nil && len(value.BS) == 0)
	if combined != nil {
		attrs[attrName] = combined
	}
	return changed, nil
}

// combineDeleteValues returns the DELETE action's subtraction of the
// operand set's members from a stored set of the same kind; emptied
// reports an empty result, which the caller answers by deleting the
// attribute. A type-incompatible pair answers ErrTypeMismatch, per
// combineAddValues.
func combineDeleteValues(existing, value *dbstore.AttributeValue) (result *dbstore.AttributeValue, emptied bool, err error) {
	// The DELETE action removes elements from a set; the operand must be a
	// set of the same type as the stored attribute.
	if (value.SS != nil && existing.SS == nil) || (value.NS != nil && existing.NS == nil) || (value.BS != nil && existing.BS == nil) {
		return nil, false, ErrTypeMismatch
	}
	if value.SS == nil && value.NS == nil && value.BS == nil {
		return nil, false, ErrTypeMismatch
	}
	if value.SS != nil {
		remaining := removeStringSetMembers(existing.SS, value.SS, stringSetCanonical)
		if len(remaining) == 0 {
			return nil, true, nil
		}
		return dbstore.StringSet(remaining), false, nil
	}
	if value.NS != nil {
		remaining := removeStringSetMembers(existing.NS, value.NS, normalizeNumberString)
		if len(remaining) == 0 {
			return nil, true, nil
		}
		return dbstore.NumberSet(remaining), false, nil
	}
	remaining := removeBinarySetMembers(existing.BS, value.BS)
	if len(remaining) == 0 {
		return nil, true, nil
	}
	return dbstore.BinarySet(remaining), false, nil
}

// applyDeleteAction applies one DynamoDB DELETE action: it removes the
// operand set's elements from a stored set of the same type, deletes the
// attribute when the subtraction empties it, and treats an absent attribute
// as a no-op. The stored value is read from preUpdate exactly as
// applyAddAction reads its base; the fresh-assignment rule and the
// ErrTypeMismatch contract follow it too.
func applyDeleteAction(attrs, preUpdate map[string]*dbstore.AttributeValue, attrName string, value *dbstore.AttributeValue) (bool, error) {
	existing, exists := preUpdate[attrName]
	if !exists {
		return false, nil
	}
	result, emptied, err := combineDeleteValues(existing, value)
	if err != nil {
		return false, err
	}
	if emptied {
		delete(attrs, attrName)
		return true, nil
	}
	if result != nil {
		attrs[attrName] = result
	}
	return true, nil
}

// applyAttributeUpdatesWithTracking applies the legacy AttributeUpdates
// map. Every member is validated — an unknown Action, a Value missing where
// the action requires one (Value is optional only for DELETE), or a Value
// that does not parse rejects the whole update instead of skipping the
// member or storing a nil attribute. An omitted Action defaults to PUT (the
// model documents PUT as the default). DELETE removes the attribute when no
// Value is supplied and subtracts the Value set's elements from the stored
// set when one is; ADD performs the number/set arithmetic.
func applyAttributeUpdatesWithTracking(attrs map[string]*dbstore.AttributeValue, updates interface{}) ([]string, error) {
	updatesMap, ok := updates.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	parseUpdateValue := func(attrName, action string, value interface{}) (*dbstore.AttributeValue, error) {
		av, err := parseAttributeValue(value)
		if err != nil {
			return nil, err
		}
		if av == nil {
			return nil, NewAPIError("com.amazon.coral.validate#ValidationException",
				fmt.Sprintf("One or more parameter values were invalid: missing or malformed Value for AttributeUpdates member %s under action %s", attrName, action),
				http.StatusBadRequest)
		}
		return av, nil
	}

	var updatedAttrs []string
	for attrName, update := range updatesMap {
		updateAction, ok := update.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}

		action, _ := updateAction["Action"].(string)
		if action == "" {
			action = "PUT"
		}
		switch action {
		case "PUT":
			value, hasValue := updateAction["Value"]
			if !hasValue {
				return nil, NewAPIError("com.amazon.coral.validate#ValidationException",
					fmt.Sprintf("One or more parameter values were invalid: Value must be specified for AttributeUpdates action PUT on attribute %s", attrName),
					http.StatusBadRequest)
			}
			av, err := parseUpdateValue(attrName, action, value)
			if err != nil {
				return nil, err
			}
			attrs[attrName] = av
			updatedAttrs = append(updatedAttrs, attrName)
		case "DELETE":
			if value, hasValue := updateAction["Value"]; hasValue {
				delValue, err := parseUpdateValue(attrName, action, value)
				if err != nil {
					return nil, err
				}
				changed, err := applyDeleteAction(attrs, attrs, attrName, delValue)
				if err != nil {
					return nil, err
				}
				if changed {
					updatedAttrs = append(updatedAttrs, attrName)
				}
			} else {
				delete(attrs, attrName)
				updatedAttrs = append(updatedAttrs, attrName)
			}
		case "ADD":
			value, hasValue := updateAction["Value"]
			if !hasValue {
				return nil, NewAPIError("com.amazon.coral.validate#ValidationException",
					fmt.Sprintf("One or more parameter values were invalid: Value must be specified for AttributeUpdates action ADD on attribute %s", attrName),
					http.StatusBadRequest)
			}
			addValue, err := parseUpdateValue(attrName, action, value)
			if err != nil {
				return nil, err
			}
			changed, err := applyAddAction(attrs, attrs, attrName, addValue)
			if err != nil {
				return nil, err
			}
			if changed {
				updatedAttrs = append(updatedAttrs, attrName)
			}
		default:
			return nil, NewAPIError("com.amazon.coral.validate#ValidationException",
				fmt.Sprintf("Invalid AttributeUpdates Action value %q: must be one of ADD, DELETE, PUT", action),
				http.StatusBadRequest)
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

func resolveNameStrict(token string, names map[string]string) (string, error) {
	if strings.HasPrefix(token, "#") {
		if names == nil {
			return "", ErrInvalidParameter
		}
		if resolved, ok := names[token]; ok {
			return resolved, nil
		}
		return "", ErrInvalidParameter
	}
	return token, nil
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
