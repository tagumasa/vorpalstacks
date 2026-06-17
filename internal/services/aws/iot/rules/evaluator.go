package rules

import (
	"encoding/json"
	"fmt"
	"regexp"
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

type unknownValue struct{}

func (unknownValue) String() string { return "UNKNOWN" }

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
	case *InExpr:
		return e.evalIn(v)
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
	return unknownValue{}, nil
}

func (e *Evaluator) resolveJsonPath(path []string) (interface{}, error) {
	if len(path) == 0 {
		return nil, nil
	}
	current, ok := e.Payload[path[0]]
	if !ok {
		return unknownValue{}, nil
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
	}

	left, err := e.Eval(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := e.Eval(expr.Right)
	if err != nil {
		return nil, err
	}

	op, ok := binaryOps[expr.Op]
	if !ok {
		return nil, fmt.Errorf("evaluator: unknown operator %s", expr.Op)
	}
	return op(left, right)
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

func (e *Evaluator) evalIn(expr *InExpr) (interface{}, error) {
	val, err := e.Eval(expr.Expr)
	if err != nil {
		return nil, err
	}
	for _, v := range expr.Values {
		candidate, err := e.Eval(v)
		if err != nil {
			return nil, err
		}
		if compareEqual(val, candidate) {
			if expr.Not {
				return false, nil
			}
			return true, nil
		}
	}
	if expr.Not {
		return true, nil
	}
	return false, nil
}

func (e *Evaluator) evalFunction(call *FunctionCall) (interface{}, error) {
	name := strings.ToLower(call.Name)
	handler, ok := functionHandlers[name]
	if !ok {
		return nil, fmt.Errorf("evaluator: unknown function %s", name)
	}
	return handler(e, call.Args)
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

func isUnknown(v interface{}) bool {
	_, ok := v.(unknownValue)
	return ok
}

func compareEqual(a, b interface{}) bool {
	if isUnknown(a) || isUnknown(b) {
		return false
	}
	if a == nil && b == nil {
		return false
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
	var sb strings.Builder
	sb.WriteString("^")
	for _, ch := range pattern {
		switch ch {
		case '%':
			sb.WriteString(".*")
		case '_':
			sb.WriteString(".")
		case '\\', '.', '+', '*', '?', '(', ')', '[', ']', '{', '}', '^', '$', '|':
			sb.WriteByte('\\')
			sb.WriteRune(ch)
		default:
			sb.WriteRune(ch)
		}
	}
	sb.WriteString("$")
	re, err := regexp.Compile(sb.String())
	if err != nil {
		return false
	}
	return re.MatchString(text)
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
			return pi == len(pattern)-1
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

func evalExprName(expr Expr) string {
	switch e := expr.(type) {
	case *Identifier:
		return e.Name
	case *StringLiteral:
		return e.Value
	default:
		return fmt.Sprintf("%v", expr)
	}
}

func (e *Evaluator) castValue(val interface{}, typeName string) (interface{}, error) {
	switch typeName {
	case "STRING":
		return toString(val), nil
	case "NUMERIC", "NUMBER", "DOUBLE", "FLOAT":
		return toFloat(val), nil
	case "BOOLEAN", "BOOL":
		return toBool(val), nil
	case "TIMESTAMP", "DATETIME":
		s := toString(val)
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04:05", s)
			if err != nil {
				return nil, fmt.Errorf("cast to TIMESTAMP: cannot parse %q", s)
			}
		}
		return t, nil
	case "INT", "INTEGER":
		return int64(toFloat(val)), nil
	default:
		return nil, fmt.Errorf("cast: unsupported type %q", typeName)
	}
}
