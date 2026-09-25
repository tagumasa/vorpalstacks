package dynamodb

import (
	"math/big"
	"strings"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// updateOperand is one operand of a SET action's value expression: a
// document path into the item, a materialised value (a resolved
// placeholder or a literal), one of the two SET functions (if_not_exists,
// list_append), or an arithmetic combination of number-valued operands.
// Both expression surfaces — the UpdateExpression string plane and the
// PartiQL clause plane — compile into this one model and evaluate it
// against one semantics, the DynamoDB update-expression operand grammar.
//
// Evaluation resolves a document path that addresses nothing to
// (nil, nil): a missing path is legal inside if_not_exists and a missing
// operand elsewhere is the caller's validation error, so the nil travels
// rather than being rejected here.
type updateOperand interface {
	eval(attrs map[string]*dbstore.AttributeValue) (*dbstore.AttributeValue, error)
}

// updatePath reads a document path from the item, resolving expression
// attribute names through the names map.
type updatePath struct {
	path  string
	names map[string]string
}

func (p updatePath) eval(attrs map[string]*dbstore.AttributeValue) (*dbstore.AttributeValue, error) {
	return getNestedAttributeValue(attrs, p.path, p.names)
}

// updateValue carries a value materialised at compile time: a resolved
// placeholder or a literal.
type updateValue struct {
	value *dbstore.AttributeValue
}

func (v updateValue) eval(attrs map[string]*dbstore.AttributeValue) (*dbstore.AttributeValue, error) {
	return v.value, nil
}

// updateIfNotExists keeps the value at the documented path when the item
// carries one ("if the item does not contain an attribute at the specified
// path, then if_not_exists evaluates to operand; otherwise, it evaluates
// to path") and otherwise evaluates the fallback operand.
type updateIfNotExists struct {
	path     updatePath
	fallback updateOperand
}

func (o updateIfNotExists) eval(attrs map[string]*dbstore.AttributeValue) (*dbstore.AttributeValue, error) {
	existing, err := o.path.eval(attrs)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	return o.fallback.eval(attrs)
}

// updateListAppend concatenates the lists its two operands resolve to.
type updateListAppend struct {
	first  updateOperand
	second updateOperand
}

func (o updateListAppend) eval(attrs map[string]*dbstore.AttributeValue) (*dbstore.AttributeValue, error) {
	list1, err := o.first.eval(attrs)
	if err != nil {
		return nil, err
	}
	list2, err := o.second.eval(attrs)
	if err != nil {
		return nil, err
	}
	if list1 == nil || list1.L == nil || list2 == nil || list2.L == nil {
		// Both operands must evaluate to lists: an operand that resolved
		// to nothing (a path naming an attribute the item does not hold)
		// or to a non-list value is an evaluation error, never a silent
		// single-operand result.
		return nil, ErrTypeMismatch
	}
	combined := make([]*dbstore.AttributeValue, 0, len(list1.L)+len(list2.L))
	combined = append(combined, list1.L...)
	combined = append(combined, list2.L...)
	return dbstore.ListValue(combined), nil
}

// updateSetFunction is a PartiQL SET-clause SET function — set_add folds
// its value operand into the set (or number) the assignment's target
// holds, set_delete subtracts from it. The function has no value of its
// own: the SET applier routes it through the ADD and DELETE action
// appliers, whose semantics they are, so nested use (inside list_append
// or if_not_exists) is a validation error. path carries the path the
// first argument names: the only documented form addresses the
// assignment's own target, and the applier enforces that correspondence.
type updateSetFunction struct {
	add     bool
	path    string
	operand updateOperand
}

func (f updateSetFunction) eval(attrs map[string]*dbstore.AttributeValue) (*dbstore.AttributeValue, error) {
	return nil, ErrInvalidParameter
}

// arithmeticTerm is one link of an arithmetic chain: the operand and the
// operator that joins it to the accumulated result. The first term carries
// "+".
type arithmeticTerm struct {
	op      string
	operand updateOperand
}

// updateArithmetic folds a +/- chain of number-valued operands.
type updateArithmetic struct {
	terms []arithmeticTerm
}

func (o updateArithmetic) eval(attrs map[string]*dbstore.AttributeValue) (*dbstore.AttributeValue, error) {
	var result *big.Rat
	for _, term := range o.terms {
		attr, err := term.operand.eval(attrs)
		if err != nil {
			return nil, err
		}
		// Arithmetic operands must be numbers present in the item or the
		// value map.
		if attr == nil || attr.N == nil {
			return nil, ErrInvalidParameter
		}
		val, ok := new(big.Rat).SetString(*attr.N)
		if !ok {
			return nil, ErrInvalidParameter
		}
		if result == nil {
			result = val
			continue
		}
		if term.op == "-" {
			result = new(big.Rat).Sub(result, val)
		} else {
			result = new(big.Rat).Add(result, val)
		}
	}
	if result == nil {
		return nil, ErrInvalidParameter
	}
	return dbstore.NumberValue(normalizeNumber(result)), nil
}

// parseUpdateOperand compiles one SET-action value operand of the
// UpdateExpression string plane into the shared operand model. Value
// placeholders resolve at compile time — an undefined placeholder is a
// validation error before any evaluation runs.
func parseUpdateOperand(s string, values map[string]*dbstore.AttributeValue, names map[string]string) (updateOperand, error) {
	s = strings.TrimSpace(s)
	parts := splitArithmeticExpression(s)
	if len(parts) == 1 {
		return parseUpdateOperandTerm(s, values, names)
	}

	// An arithmetic chain: the parts alternate operands and +/-
	// separators, exactly the decomposition the splitter produces at
	// parenthesis depth zero.
	var terms []arithmeticTerm
	for i := 0; i < len(parts); i++ {
		part := strings.TrimSpace(parts[i])
		if part == "" {
			continue
		}
		op, operandStr := "+", part
		if i != 0 && (part == "+" || part == "-") {
			op = part
			i++
			if i < len(parts) {
				operandStr = strings.TrimSpace(parts[i])
			} else {
				operandStr = ""
			}
		}
		if operandStr == "" {
			continue
		}
		term, err := parseUpdateOperandTerm(operandStr, values, names)
		if err != nil {
			return nil, err
		}
		terms = append(terms, arithmeticTerm{op: op, operand: term})
	}
	if len(terms) == 0 {
		return nil, ErrInvalidParameter
	}
	return updateArithmetic{terms: terms}, nil
}

// parseUpdateOperandTerm compiles one operand without a top-level
// arithmetic operator: one of the two SET functions, a value placeholder,
// or a bare document path.
func parseUpdateOperandTerm(s string, values map[string]*dbstore.AttributeValue, names map[string]string) (updateOperand, error) {
	if strings.HasPrefix(s, "if_not_exists(") {
		closeParen := findMatchingCloseParen(s, 13)
		if closeParen == -1 {
			return nil, ErrInvalidParameter
		}
		path, fallback, ok := splitBalancedTwoArgs(s[14:closeParen])
		if !ok {
			return nil, ErrInvalidParameter
		}
		fallbackOperand, err := parseUpdateOperand(fallback, values, names)
		if err != nil {
			return nil, err
		}
		return updateIfNotExists{
			path:     updatePath{path: strings.TrimSpace(path), names: names},
			fallback: fallbackOperand,
		}, nil
	}
	if strings.HasPrefix(s, "list_append(") {
		closeParen := findMatchingCloseParen(s, 11)
		if closeParen == -1 {
			return nil, ErrInvalidParameter
		}
		first, second, ok := splitBalancedTwoArgs(s[12:closeParen])
		if !ok {
			return nil, ErrInvalidParameter
		}
		firstOperand, err := parseUpdateOperand(first, values, names)
		if err != nil {
			return nil, err
		}
		secondOperand, err := parseUpdateOperand(second, values, names)
		if err != nil {
			return nil, err
		}
		return updateListAppend{first: firstOperand, second: secondOperand}, nil
	}
	if strings.HasPrefix(s, ":") {
		if v := resolveValue(s, values, names); v != nil {
			return updateValue{value: v}, nil
		}
		// An undefined value placeholder is a validation error.
		return nil, ErrInvalidParameter
	}
	// A bare operand is a document path to the item.
	return updatePath{path: s, names: names}, nil
}

// splitBalancedTwoArgs splits a two-argument function body at the comma
// that sits at parenthesis depth zero, so an argument carrying nested
// function calls keeps its own commas.
func splitBalancedTwoArgs(s string) (string, string, bool) {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				return s[:i], s[i+1:], true
			}
		}
	}
	return "", "", false
}
