package rules

import (
	"fmt"
)

func ValidateExpr(expr Expr) error {
	return validateExprNode(expr)
}

func validateExprNode(expr Expr) error {
	switch v := expr.(type) {
	case *StringLiteral, *NumberLiteral, *BoolLiteral, *NullLiteral, *StarExpr, *Identifier, *JsonPath:
		return nil
	case *BinaryExpr:
		if err := validateExprNode(v.Left); err != nil {
			return err
		}
		return validateExprNode(v.Right)
	case *UnaryExpr:
		return validateExprNode(v.Operand)
	case *FunctionCall:
		for _, arg := range v.Args {
			if err := validateExprNode(arg); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("rules: unrecognised expression type %T", expr)
	}
}
