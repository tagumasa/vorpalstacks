package rules

import (
	"fmt"
	"math"
)

// BinaryOp evaluates a binary operator given pre-evaluated left and right operands.
type BinaryOp func(left, right interface{}) (interface{}, error)

var binaryOps map[string]BinaryOp

func init() {
	binaryOps = map[string]BinaryOp{
		"=":           opEqual,
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
	return compareEqual(left, right), nil
}

func opNotEqual(left, right interface{}) (interface{}, error) {
	return !compareEqual(left, right), nil
}

func opLess(left, right interface{}) (interface{}, error) {
	return compareNumeric(left, right, -1), nil
}

func opGreater(left, right interface{}) (interface{}, error) {
	return compareNumeric(left, right, 1), nil
}

func opLessEqual(left, right interface{}) (interface{}, error) {
	return !compareNumeric(left, right, 1), nil
}

func opGreaterEqual(left, right interface{}) (interface{}, error) {
	return !compareNumeric(left, right, -1), nil
}

func opLike(left, right interface{}) (interface{}, error) {
	return matchLike(toString(left), toString(right)), nil
}

func opAdd(left, right interface{}) (interface{}, error) {
	return toFloat(left) + toFloat(right), nil
}

func opSub(left, right interface{}) (interface{}, error) {
	return toFloat(left) - toFloat(right), nil
}

func opMul(left, right interface{}) (interface{}, error) {
	return toFloat(left) * toFloat(right), nil
}

func opDiv(left, right interface{}) (interface{}, error) {
	if toFloat(right) == 0 {
		return nil, fmt.Errorf("evaluator: division by zero")
	}
	return toFloat(left) / toFloat(right), nil
}

func opMod(left, right interface{}) (interface{}, error) {
	return math.Mod(toFloat(left), toFloat(right)), nil
}

func opIsNull(left, right interface{}) (interface{}, error) {
	return left == nil, nil
}

func opIsNotNull(left, right interface{}) (interface{}, error) {
	return left != nil, nil
}
