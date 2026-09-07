package dynamodb

import (
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	"vorpalstacks/pkg/sqlparser"
)

func objectLiteralToAttributesWithParams(obj *sqlparser.ObjectLiteral, params *partiQLParams) (map[string]*dbstore.AttributeValue, error) {
	result := make(map[string]*dbstore.AttributeValue)
	for _, prop := range obj.Properties {
		key := extractPropertyKey(prop)
		if key == "" {
			continue
		}
		value, err := exprToAttributeValueWithParams(prop.Value, params)
		if err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, nil
}

// resolvePlaceholder materialises the parameter bound to a statement
// placeholder as an AttributeValue. The tokenizer encodes `?` as :v<N>,
// numbered per statement. DynamoDB requires one parameter value per
// placeholder, so a reference with no bound parameter — or a parameter
// that is not a valid AttributeValue — is a ValidationException, never a
// stored value.
func resolvePlaceholder(val []byte, params *partiQLParams) (*dbstore.AttributeValue, error) {
	if params == nil {
		return nil, ErrInvalidParameter
	}
	idx, err := strconv.Atoi(strings.TrimPrefix(string(val), ":v"))
	if err != nil || idx < 1 || idx > len(params.Parameters) {
		return nil, ErrInvalidParameter
	}
	parsed := parseAttributeValue(params.Parameters[idx-1])
	if parsed == nil {
		return nil, ErrInvalidParameter
	}
	return parsed, nil
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

// exprToAttributeValueWithParams materialises a PartiQL value expression —
// a literal, a bound parameter, or any nesting of lists and objects
// containing them — as an AttributeValue.
func exprToAttributeValueWithParams(expr sqlparser.Expr, params *partiQLParams) (*dbstore.AttributeValue, error) {
	switch e := expr.(type) {
	case *sqlparser.SQLVal:
		if e.Type == sqlparser.ValArg {
			return resolvePlaceholder(e.Val, params)
		}
	case *sqlparser.ObjectLiteral:
		m, err := objectLiteralToAttributesWithParams(e, params)
		if err != nil {
			return nil, err
		}
		return &dbstore.AttributeValue{M: m}, nil
	case sqlparser.ValTuple:
		l, err := tupleToAttributeListWithParams(e, params)
		if err != nil {
			return nil, err
		}
		return &dbstore.AttributeValue{L: l}, nil
	}
	// The literal materialiser is the only path that builds stored numbers
	// without going through the wire AttributeValue parser, so the DynamoDB
	// Number contract is enforced there.
	v, err := exprToAttributeValue(expr)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// exprToAttributeValue materialises literal value expressions only;
// placeholders must go through exprToAttributeValueWithParams. Number
// literals are held to the DynamoDB Number contract — an out-of-range or
// over-precision literal is a ValidationException, never a stored value:
// the storage-key number encoder renders a fixed fractional width that is
// injective only over validated numbers, so an unvalidated literal could
// collide with another item's key by differing solely beyond the rendered
// digits.
func exprToAttributeValue(expr sqlparser.Expr) (*dbstore.AttributeValue, error) {
	switch e := expr.(type) {
	case *sqlparser.SQLVal:
		switch e.Type {
		case sqlparser.StrVal:
			s := string(e.Val)
			return &dbstore.AttributeValue{S: &s}, nil
		case sqlparser.IntVal, sqlparser.FloatVal:
			s := string(e.Val)
			if !isValidDynamoDBNumber(s) {
				return nil, ErrInvalidParameter
			}
			return &dbstore.AttributeValue{N: &s}, nil
		}
	case *sqlparser.NullVal:
		return &dbstore.AttributeValue{NULL: proto.Bool(true)}, nil
	case sqlparser.BoolVal:
		b := bool(e)
		return &dbstore.AttributeValue{BOOL: &b}, nil
	case *sqlparser.UnaryExpr:
		// The grammar reduces a sign-prefixed number literal to a unary
		// expression around the bare literal; only that shape is a Number —
		// any other unary expression keeps the NULL fallback.
		if e.Operator != sqlparser.UMinusStr && e.Operator != sqlparser.UPlusStr {
			return &dbstore.AttributeValue{NULL: proto.Bool(true)}, nil
		}
		inner, ok := e.Expr.(*sqlparser.SQLVal)
		if !ok || (inner.Type != sqlparser.IntVal && inner.Type != sqlparser.FloatVal) {
			return &dbstore.AttributeValue{NULL: proto.Bool(true)}, nil
		}
		s := string(inner.Val)
		if e.Operator == sqlparser.UMinusStr {
			s = "-" + s
		}
		if !isValidDynamoDBNumber(s) {
			return nil, ErrInvalidParameter
		}
		return &dbstore.AttributeValue{N: &s}, nil
	}
	return &dbstore.AttributeValue{NULL: proto.Bool(true)}, nil
}

func tupleToAttributeListWithParams(tuple sqlparser.ValTuple, params *partiQLParams) ([]*dbstore.AttributeValue, error) {
	result := make([]*dbstore.AttributeValue, 0, len(tuple))
	for _, item := range tuple {
		value, err := exprToAttributeValueWithParams(item, params)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}
