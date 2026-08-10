package rules

import (
	"math"
)

// BinaryOp evaluates a binary operator given pre-evaluated left and right operands.
type BinaryOp func(left, right interface{}) (interface{}, error)

var binaryOps map[string]BinaryOp

func init() {
	binaryOps = map[string]BinaryOp{
		"=":           opEqual,
		"==":          opEqual,
		"!=":          opNotEqual,
		"<>":          opNotEqual,
		"<":           opLess,
		">":           opGreater,
		"<=":          opLessEqual,
		">=":          opGreaterEqual,
		"LIKE":        opLike,
		"+":           opAdd,
		"-":           opSub,
		"*":           opMul,
		"/":           opDiv,
		"%":           opMod,
		"IS NULL":     opIsNull,
		"IS NOT NULL": opIsNotNull,
	}
}

func opEqual(left, right interface{}) (interface{}, error) {
	if isUnknown(left) || isUnknown(right) {
		return unknownValue{}, nil
	}
	return compareEqual(left, right), nil
}

func opNotEqual(left, right interface{}) (interface{}, error) {
	if isUnknown(left) || isUnknown(right) {
		return false, nil
	}
	return !compareEqual(left, right), nil
}

func opLess(left, right interface{}) (interface{}, error) {
	if isUnknown(left) || isUnknown(right) {
		return false, nil
	}
	return compareNumeric(left, right, -1), nil
}

func opGreater(left, right interface{}) (interface{}, error) {
	if isUnknown(left) || isUnknown(right) {
		return false, nil
	}
	return compareNumeric(left, right, 1), nil
}

func opLessEqual(left, right interface{}) (interface{}, error) {
	if isUnknown(left) || isUnknown(right) {
		return false, nil
	}
	return !compareNumeric(left, right, 1), nil
}

func opGreaterEqual(left, right interface{}) (interface{}, error) {
	if isUnknown(left) || isUnknown(right) {
		return false, nil
	}
	return !compareNumeric(left, right, -1), nil
}

func opLike(left, right interface{}) (interface{}, error) {
	if isUnknown(left) || isUnknown(right) {
		return false, nil
	}
	return matchLike(toString(left), toString(right)), nil
}

func opAdd(left, right interface{}) (interface{}, error) {
	if isUnknown(left) || isUnknown(right) {
		return unknownValue{}, nil
	}
	// AWS IoT SQL overloads + for string concatenation when either
	// operand is a string.  Only when both operands are numeric does
	// the operator perform arithmetic addition.
	if isStringOperand(left) || isStringOperand(right) {
		return toString(left) + toString(right), nil
	}
	return toFloat(left) + toFloat(right), nil
}

func opSub(left, right interface{}) (interface{}, error) {
	if isUnknown(left) || isUnknown(right) {
		return unknownValue{}, nil
	}
	return toFloat(left) - toFloat(right), nil
}

func opMul(left, right interface{}) (interface{}, error) {
	if isUnknown(left) || isUnknown(right) {
		return unknownValue{}, nil
	}
	return toFloat(left) * toFloat(right), nil
}

func opDiv(left, right interface{}) (interface{}, error) {
	if isUnknown(left) || isUnknown(right) {
		return unknownValue{}, nil
	}
	if toFloat(right) == 0 {
		// AWS IoT SQL returns Undefined for division by zero so the
		// rule continues with a non-matching WHERE clause rather than
		// causing the entire rule to error.
		return unknownValue{}, nil
	}
	return toFloat(left) / toFloat(right), nil
}

func opMod(left, right interface{}) (interface{}, error) {
	if isUnknown(left) || isUnknown(right) {
		return unknownValue{}, nil
	}
	if toFloat(right) == 0 {
		return unknownValue{}, nil
	}
	return math.Mod(toFloat(left), toFloat(right)), nil
}

func opIsNull(left, right interface{}) (interface{}, error) {
	// AWS IoT SQL distinguishes Null (JSON null) from Undefined (field
	// absent).  IS NULL returns true only for an actual Null value.
	return left == nil, nil
}

func opIsNotNull(left, right interface{}) (interface{}, error) {
	// IS NOT NULL is the logical complement of IS NULL: returns true
	// for any value that is not JSON null, including Undefined.
	return left != nil, nil
}

// isStringOperand reports whether the value is a Go string.  This is
// used to decide whether the + operator performs string concatenation
// (AWS IoT SQL overloads + for strings) or numeric addition.
func isStringOperand(v interface{}) bool {
	switch v.(type) {
	case string:
		return true
	default:
		return false
	}
}
