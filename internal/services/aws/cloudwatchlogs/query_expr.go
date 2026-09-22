package cloudwatchlogs

import (
	"encoding/json"
	"math"
	"regexp"
	"strconv"
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
	x       exprNode
	list    []exprNode // literal list elements, or a subquery
	operand exprNode   // the array-valued right-hand side: "put the array after in"
	neg     bool
}

type likeNode struct {
	x      exprNode
	substr string         // a quoted pattern: literal substring containment
	re     *regexp.Regexp // a regex-mode like pattern
	neg    bool
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
	name := n.segs[0].name
	var cur interface{}
	if sv, ok := (*row).structs[name]; ok {
		// The structural view: dot notation traverses the decoded tree of
		// a field that was structure at ingest or creation time.
		cur = sv
	} else {
		raw, ok := (*row).fields[name]
		if !ok {
			// Fields without the @ prefix resolve against the @-prefixed
			// discoverable fields, e.g. timestamp matches @timestamp.
			raw, ok = (*row).fields["@"+name]
		}
		if !ok {
			return nil
		}
		// A plain string value is an opaque leaf even when its content is
		// valid JSON — jsonParse is the documented way into it.
		cur = raw
	}
	return descendPath(cur, n.segs[1:])
}

// descendPath walks structure segments: maps by key, lists by index.
// String values are opaque leaves — "Dot notation traverses only fields
// that are stored as structurally nested JSON objects at ingest time";
// the explicit way into a JSON-encoded string is jsonParse.
func descendPath(start interface{}, segs []pathSeg) interface{} {
	cur := start
	for _, seg := range segs {
		if cur == nil {
			return nil
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

// parseJSONValue decodes JSON object or array text; other strings carry no
// structure.
func parseJSONValue(s string) (interface{}, bool) {
	trimmed := strings.TrimSpace(s)
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		var decoded interface{}
		if err := json.Unmarshal([]byte(trimmed), &decoded); err == nil {
			return decoded, true
		}
	}
	return nil, false
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
		// "Arithmetic operators accept numeric data types as arguments
		// and return numeric results" — a non-numeric operand is not
		// concatenated or coerced; the operation yields null like the
		// other operators' non-numeric cases.
		ln, lok := asNumber(l)
		rn, rok := asNumber(r)
		if !lok || !rok {
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
	if n.operand != nil {
		// The array form: "To check for elements in an array, put the
		// array after in" — the operand evaluates to a list whose elements
		// are compared against the left side.
		for _, v := range arrayValues(n.operand.eval(row, ctx)) {
			if valuesEqual(x, v) {
				found = true
				break
			}
		}
		if n.neg {
			return !found
		}
		return found
	}
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

// arrayValues coerces a value to its element slice: typed lists pass
// through, and string values decode from JSON on demand — the same lazy
// structure decoding field access rides.
func arrayValues(v interface{}) []interface{} {
	switch x := v.(type) {
	case []interface{}:
		return x
	case string:
		if decoded, ok := parseJSONValue(x); ok {
			if list, ok := decoded.([]interface{}); ok {
				return list
			}
		}
	}
	return nil
}

func (n *likeNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	v := asString(n.x.eval(row, ctx))
	var matched bool
	switch {
	case n.re != nil:
		matched = n.re.MatchString(v)
	case n.substr != "":
		// "You can use the keyword phrases like and not like to match
		// substrings" — a quoted pattern is a literal substring, not an
		// anchored full-string glob.
		matched = strings.Contains(v, n.substr)
	}
	if n.neg {
		return !matched
	}
	return matched
}

func (n *funcNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	return callQueryFunction(strings.ToLower(n.name), n.args, row, ctx)
}

// --- expression rendering ---

func exprDisplay(args []exprNode) string {
	parts := make([]string, 0, len(args))
	for _, a := range args {
		parts = append(parts, exprText(a))
	}
	return strings.Join(parts, ", ")
}

// pathText renders a field path's segments: names joined by dots, index
// segments as bracketed integers.
func pathText(segs []pathSeg) string {
	var b strings.Builder
	for i, s := range segs {
		if i > 0 {
			b.WriteByte('.')
		}
		if s.isIdx {
			b.WriteString("[")
			b.WriteString(strconv.Itoa(s.index))
			b.WriteString("]")
		} else {
			b.WriteString(s.name)
		}
	}
	return b.String()
}

// exprText renders an expression back to its written form — the name a
// projection or aggregation column carries when the query gives it no
// explicit alias, so every node kind must render: a kind that falls
// through names its column with the empty string.
func exprText(e exprNode) string {
	switch n := e.(type) {
	case *literalNode:
		return asString(n.val)
	case *fieldNode:
		return pathText(n.segs)
	case *funcNode:
		return n.name + "(" + exprDisplay(n.args) + ")"
	case *binOpNode:
		return exprText(n.l) + " " + n.op + " " + exprText(n.r)
	case *unaryNode:
		if n.op == "not" {
			return "not " + exprText(n.x)
		}
		return n.op + exprText(n.x)
	case *inNode:
		op := " in "
		if n.neg {
			op = " not in "
		}
		if n.operand != nil {
			return exprText(n.x) + op + exprText(n.operand)
		}
		return exprText(n.x) + op + "[" + exprDisplay(n.list) + "]"
	case *likeNode:
		pattern := "'" + n.substr + "'"
		if n.re != nil {
			pattern = "/" + n.re.String() + "/"
		}
		op := " like "
		if n.neg {
			op = " not like "
		}
		return exprText(n.x) + op + pattern
	case *accessNode:
		var b strings.Builder
		b.WriteString(exprText(n.base))
		for _, s := range n.segs {
			if s.isIdx {
				b.WriteString("[")
				b.WriteString(strconv.Itoa(s.index))
				b.WriteString("]")
			} else {
				b.WriteString(".")
				b.WriteString(s.name)
			}
		}
		return b.String()
	case *arrayLiteralNode:
		return "[" + exprDisplay(n.items) + "]"
	case *subqueryNode:
		parts := make([]string, 0, len(n.toks))
		for _, t := range n.toks {
			parts = append(parts, t.raw)
		}
		return "(" + strings.Join(parts, " ") + ")"
	}
	return ""
}
