package rules

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type Evaluator struct {
	Payload   map[string]interface{}
	Topic     string
	ClientID  string
	Timestamp time.Time
}

func NewEvaluator(payload map[string]interface{}, topic, clientID string) *Evaluator {
	return &Evaluator{
		Payload:   payload,
		Topic:     topic,
		ClientID:  clientID,
		Timestamp: time.Now().UTC(),
	}
}

func (e *Evaluator) Eval(expr Expr) (interface{}, error) {
	switch v := expr.(type) {
	case *StringLiteral:
		return v.Value, nil
	case *NumberLiteral:
		return v.Value, nil
	case *BoolLiteral:
		return v.Value, nil
	case *NullLiteral:
		return nil, nil
	case *Identifier:
		return e.resolveIdentifier(v.Name)
	case *JsonPath:
		return e.resolveJsonPath(v.Path)
	case *StarExpr:
		return e.Payload, nil
	case *BinaryExpr:
		return e.evalBinary(v)
	case *UnaryExpr:
		return e.evalUnary(v)
	case *FunctionCall:
		return e.evalFunction(v)
	}
	return nil, fmt.Errorf("evaluator: unknown expression type %T", expr)
}

func (e *Evaluator) resolveIdentifier(name string) (interface{}, error) {
	switch strings.ToLower(name) {
	case "clientid", "client_id":
		return e.ClientID, nil
	case "topic":
		return e.Topic, nil
	}
	if val, ok := e.Payload[name]; ok {
		return val, nil
	}
	return nil, nil
}

func (e *Evaluator) resolveJsonPath(path []string) (interface{}, error) {
	if len(path) == 0 {
		return nil, nil
	}
	current, ok := e.Payload[path[0]]
	if !ok {
		return nil, nil
	}
	for _, part := range path[1:] {
		current = navigateJSON(current, part)
		if current == nil {
			return nil, nil
		}
	}
	return current, nil
}

func navigateJSON(data interface{}, key string) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		return v[key]
	case []interface{}:
		idx, err := strconv.Atoi(key)
		if err != nil {
			return nil
		}
		if idx >= 0 && idx < len(v) {
			return v[idx]
		}
		return nil
	}
	return nil
}

func (e *Evaluator) evalBinary(expr *BinaryExpr) (interface{}, error) {
	switch expr.Op {
	case "AND":
		left, err := e.Eval(expr.Left)
		if err != nil {
			return nil, err
		}
		if !toBool(left) {
			return false, nil
		}
		right, err := e.Eval(expr.Right)
		if err != nil {
			return nil, err
		}
		return toBool(right), nil
	case "OR":
		left, err := e.Eval(expr.Left)
		if err != nil {
			return nil, err
		}
		if toBool(left) {
			return true, nil
		}
		right, err := e.Eval(expr.Right)
		if err != nil {
			return nil, err
		}
		return toBool(right), nil
	case "NOT":
		right, err := e.Eval(expr.Right)
		if err != nil {
			return nil, err
		}
		return !toBool(right), nil
	}

	left, err := e.Eval(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := e.Eval(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Op {
	case "=":
		return compareEqual(left, right), nil
	case "!=", "<>":
		return !compareEqual(left, right), nil
	case "<":
		return compareNumeric(left, right, -1), nil
	case ">":
		return compareNumeric(left, right, 1), nil
	case "<=":
		return !compareNumeric(left, right, 1), nil
	case ">=":
		return !compareNumeric(left, right, -1), nil
	case "LIKE":
		return matchLike(toString(left), toString(right)), nil
	case "+":
		return toFloat(left) + toFloat(right), nil
	case "-":
		return toFloat(left) - toFloat(right), nil
	case "*":
		return toFloat(left) * toFloat(right), nil
	case "/":
		if toFloat(right) == 0 {
			return nil, fmt.Errorf("evaluator: division by zero")
		}
		return toFloat(left) / toFloat(right), nil
	case "%":
		return math.Mod(toFloat(left), toFloat(right)), nil
	case "IS NULL":
		return left == nil, nil
	case "IS NOT NULL":
		return left != nil, nil
	default:
		return nil, fmt.Errorf("evaluator: unknown operator %s", expr.Op)
	}
}

func (e *Evaluator) evalUnary(expr *UnaryExpr) (interface{}, error) {
	switch expr.Op {
	case "NOT":
		val, err := e.Eval(expr.Operand)
		if err != nil {
			return nil, err
		}
		return !toBool(val), nil
	case "-":
		val, err := e.Eval(expr.Operand)
		if err != nil {
			return nil, err
		}
		return -toFloat(val), nil
	default:
		return nil, fmt.Errorf("evaluator: unknown unary operator %s", expr.Op)
	}
}

func (e *Evaluator) evalFunction(call *FunctionCall) (interface{}, error) {
	name := strings.ToLower(call.Name)
	switch name {
	case "topic":
		return e.Topic, nil
	case "timestamp":
		return e.Timestamp.Unix(), nil
	case "clientid":
		return e.ClientID, nil
	case "concat":
		var parts []string
		for _, arg := range call.Args {
			val, err := e.Eval(arg)
			if err != nil {
				return nil, err
			}
			parts = append(parts, toString(val))
		}
		return strings.Join(parts, ""), nil
	case "length":
		if len(call.Args) < 1 {
			return nil, fmt.Errorf("evaluator: length() requires 1 argument")
		}
		val, err := e.Eval(call.Args[0])
		if err != nil {
			return nil, err
		}
		switch v := val.(type) {
		case string:
			return float64(len(v)), nil
		case []interface{}:
			return float64(len(v)), nil
		case map[string]interface{}:
			return float64(len(v)), nil
		default:
			return float64(len(toString(v))), nil
		}
	case "upper":
		if len(call.Args) < 1 {
			return nil, nil
		}
		val, err := e.Eval(call.Args[0])
		if err != nil {
			return nil, err
		}
		return strings.ToUpper(toString(val)), nil
	case "lower":
		if len(call.Args) < 1 {
			return nil, nil
		}
		val, err := e.Eval(call.Args[0])
		if err != nil {
			return nil, err
		}
		return strings.ToLower(toString(val)), nil
	case "trim":
		if len(call.Args) < 1 {
			return nil, nil
		}
		val, err := e.Eval(call.Args[0])
		if err != nil {
			return nil, err
		}
		return strings.TrimSpace(toString(val)), nil
	case "replace":
		if len(call.Args) < 3 {
			return nil, nil
		}
		str, _ := e.Eval(call.Args[0])
		old, _ := e.Eval(call.Args[1])
		repl, _ := e.Eval(call.Args[2])
		return strings.ReplaceAll(toString(str), toString(old), toString(repl)), nil
	case "substring":
		if len(call.Args) < 2 {
			return nil, nil
		}
		str, _ := e.Eval(call.Args[0])
		startVal, _ := e.Eval(call.Args[1])
		start := int(toFloat(startVal))
		s := toString(str)
		length := len(s) - start
		if len(call.Args) >= 3 {
			lenVal, _ := e.Eval(call.Args[2])
			length = int(toFloat(lenVal))
		}
		if start < 0 {
			return "", nil
		}
		if start+length > len(s) {
			return s[start:], nil
		}
		if length < 0 {
			return s[start:], nil
		}
		return s[start : start+length], nil
	case "cast":
		return nil, fmt.Errorf("evaluator: cast() not yet implemented")
	case "abs":
		if len(call.Args) < 1 {
			return nil, nil
		}
		val, _ := e.Eval(call.Args[0])
		return math.Abs(toFloat(val)), nil
	case "ceil":
		if len(call.Args) < 1 {
			return nil, nil
		}
		val, _ := e.Eval(call.Args[0])
		return math.Ceil(toFloat(val)), nil
	case "floor":
		if len(call.Args) < 1 {
			return nil, nil
		}
		val, _ := e.Eval(call.Args[0])
		return math.Floor(toFloat(val)), nil
	default:
		return nil, fmt.Errorf("evaluator: unknown function %s", name)
	}
}

func toBool(val interface{}) bool {
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return v != ""
	case float64:
		return v != 0
	case int:
		return v != 0
	case int64:
		return v != 0
	case json.Number:
		n, _ := v.Int64()
		return n != 0
	default:
		return true
	}
}

func toFloat(val interface{}) float64 {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case json.Number:
		f, _ := v.Float64()
		return f
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	case bool:
		if v {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func toString(val interface{}) string {
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case bool:
		return strconv.FormatBool(v)
	case json.Number:
		return v.String()
	default:
		b, _ := json.Marshal(val)
		return string(b)
	}
}

func compareEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return toString(a) == toString(b)
}

func compareNumeric(a, b interface{}, direction int) bool {
	aFloat := toFloat(a)
	bFloat := toFloat(b)
	switch direction {
	case -1:
		return aFloat < bFloat
	case 1:
		return aFloat > bFloat
	default:
		return false
	}
}

func matchLike(text, pattern string) bool {
	parts := strings.Split(pattern, "%")
	if len(parts) == 1 {
		return text == pattern
	}
	result := true
	for i, part := range parts {
		if part == "" {
			continue
		}
		idx := strings.Index(text, part)
		if idx < 0 {
			return false
		}
		if i == 0 && idx > 0 {
			return false
		}
		if i == len(parts)-1 && idx+len(part) != len(text) {
			return false
		}
		text = text[idx+len(part):]
	}
	return result
}

func EvaluateSQL(sql string, payload map[string]interface{}, topic, clientID string) (*SelectExpr, map[string]interface{}, error) {
	parser := NewParser(sql)
	expr, err := parser.Parse()
	if err != nil {
		return nil, nil, fmt.Errorf("parse error: %w", err)
	}
	eval := NewEvaluator(payload, topic, clientID)
	result := make(map[string]interface{})
	for i, field := range expr.Fields {
		name := field.Alias
		if name == "" {
			switch f := field.Expression.(type) {
			case *Identifier:
				name = f.Name
			case *JsonPath:
				name = strings.Join(f.Path, ".")
			default:
				name = fmt.Sprintf("col%d", i)
			}
		}
		val, err := eval.Eval(field.Expression)
		if err != nil {
			return nil, nil, fmt.Errorf("eval error on field %s: %w", name, err)
		}
		result[name] = val
	}
	return expr, result, nil
}

func MatchTopicFilter(ruleTopic, messageTopic string) bool {
	return topicMatches(ruleTopic, messageTopic)
}

func topicMatches(pattern, topic string) bool {
	patternParts := strings.Split(pattern, "/")
	topicParts := strings.Split(topic, "/")
	return matchTopicParts(patternParts, topicParts, 0, 0)
}

func matchTopicParts(pattern, topic []string, pi, ti int) bool {
	for pi < len(pattern) && ti < len(topic) {
		if pattern[pi] == "#" {
			if pi != len(pattern)-1 {
				return false
			}
			return true
		}
		if pattern[pi] == "+" || pattern[pi] == topic[ti] {
			pi++
			ti++
			continue
		}
		return false
	}
	if pi < len(pattern) && pattern[pi] == "#" && pi == len(pattern)-1 {
		return true
	}
	return pi == len(pattern) && ti == len(topic)
}
