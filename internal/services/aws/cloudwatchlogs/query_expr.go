package cloudwatchlogs

import (
	"encoding/json"
	"math"
	"regexp"
	"strings"
)

// pathSeg is one segment of a field path: either a map key or a list index.
// Structure access follows the documented dot and bracket operators, e.g.
// json_message.users[1].action.
type pathSeg struct {
	name  string
	index int
	isIdx bool
}

// exprNode is a parsed expression. Commands parse their argument expressions
// once and evaluate the AST per row.
type exprNode interface {
	eval(row *queryResultRow, ctx *execContext) interface{}
}

type literalNode struct {
	val interface{}
}

type fieldNode struct {
	segs []pathSeg
}

type binOpNode struct {
	op   string
	l, r exprNode
}

type unaryNode struct {
	op string
	x  exprNode
}

type inNode struct {
	x    exprNode
	list []exprNode
	neg  bool
}

type likeNode struct {
	x       exprNode
	pattern string
	re      *regexp.Regexp
	neg     bool
}

type funcNode struct {
	name string
	args []exprNode
}

// accessNode accesses structure attributes on the result of an arbitrary
// expression, e.g. jsonParse(@message).status.
type accessNode struct {
	base exprNode
	segs []pathSeg
}

func (n *accessNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	return descendPath(n.base.eval(row, ctx), n.segs)
}

type arrayLiteralNode struct {
	items []exprNode
}

func (n *arrayLiteralNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	out := make([]interface{}, 0, len(n.items))
	for _, e := range n.items {
		out = append(out, e.eval(row, ctx))
	}
	return out
}

// --- evaluation ---

func (n *literalNode) eval(row *queryResultRow, ctx *execContext) interface{} { return n.val }

func (n *fieldNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	var cur interface{}
	name := n.segs[0].name
	raw, ok := (*row).fields[name]
	if !ok {
		// Fields without the @ prefix resolve against the @-prefixed
		// discoverable fields, e.g. timestamp matches @timestamp.
		raw, ok = (*row).fields["@"+name]
	}
	if !ok {
		return nil
	}
	cur = raw
	return descendPath(cur, n.segs[1:])
}

// descendPath walks structure segments: maps by key, lists by index.
// String values read from rows are decoded from JSON on demand so that
// jsonParse results round-trip through the string-typed rows.
func descendPath(start interface{}, segs []pathSeg) interface{} {
	cur := start
	for _, seg := range segs {
		if cur == nil {
			return nil
		}
		if s, isStr := cur.(string); isStr {
			trimmed := strings.TrimSpace(s)
			if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
				var decoded interface{}
				if err := json.Unmarshal([]byte(trimmed), &decoded); err == nil {
					cur = decoded
				}
			}
		}
		if seg.isIdx {
			list, ok := cur.([]interface{})
			if !ok || seg.index < 0 || seg.index >= len(list) {
				return nil
			}
			cur = list[seg.index]
			continue
		}
		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil
		}
		v, ok := m[seg.name]
		if !ok {
			return nil
		}
		cur = v
	}
	return cur
}

func (n *unaryNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	v := n.x.eval(row, ctx)
	switch n.op {
	case "not":
		return !truthy(v)
	case "-":
		if f, ok := asNumber(v); ok {
			return -f
		}
		return nil
	}
	return nil
}

func (n *binOpNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	switch n.op {
	case "and":
		return truthy(n.l.eval(row, ctx)) && truthy(n.r.eval(row, ctx))
	case "or":
		return truthy(n.l.eval(row, ctx)) || truthy(n.r.eval(row, ctx))
	}
	l := n.l.eval(row, ctx)
	r := n.r.eval(row, ctx)
	switch n.op {
	case "=", "!=":
		eq := valuesEqual(l, r)
		if n.op == "=" {
			return eq
		}
		return !eq
	case "<", "<=", ">", ">=":
		return compareValues(l, r, n.op)
	case "+", "-", "*", "/", "%", "^":
		// String concatenation for + when either side is a non-numeric
		// string, numeric arithmetic otherwise.
		ln, lok := asNumber(l)
		rn, rok := asNumber(r)
		if !lok || !rok {
			if n.op == "+" {
				return asString(l) + asString(r)
			}
			return nil
		}
		return applyArith(n.op, ln, rn)
	}
	return nil
}

func applyArith(op string, l, r float64) interface{} {
	switch op {
	case "+":
		return l + r
	case "-":
		return l - r
	case "*":
		return l * r
	case "/":
		if r == 0 {
			return nil
		}
		return l / r
	case "%":
		if int64(r) == 0 {
			return nil
		}
		return float64(int64(l) % int64(r))
	case "^":
		return math.Pow(l, r)
	}
	return nil
}

func (n *inNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	x := n.x.eval(row, ctx)
	found := false
	for _, e := range n.list {
		if sub, ok := e.(*subqueryNode); ok {
			// Subquery membership: the nested query materialises a set of
			// values that are compared against the outer field.
			for _, v := range sub.evalSet(ctx) {
				if valuesEqual(x, v) {
					found = true
					break
				}
			}
			if found {
				break
			}
			continue
		}
		if valuesEqual(x, e.eval(row, ctx)) {
			found = true
			break
		}
	}
	if n.neg {
		return !found
	}
	return found
}

func (n *likeNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	v := asString(n.x.eval(row, ctx))
	var matched bool
	if n.re != nil {
		matched = n.re.MatchString(v)
	} else {
		matched = globMatch(n.pattern, v)
	}
	if n.neg {
		return !matched
	}
	return matched
}

func (n *funcNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	return callQueryFunction(strings.ToLower(n.name), n.args, row, ctx)
}
