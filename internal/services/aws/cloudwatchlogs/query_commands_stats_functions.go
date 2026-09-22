package cloudwatchlogs

import (
	"math"
	"sort"
	"strings"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The stats command's aggregation evaluators: the arithmetic over
// aggregation calls, the shared argument views, the basic aggregation
// family and the timestamp result typing its passthrough functions
// carry. The time-series family lives beside these in its own file.

// evalAggExpr evaluates an aggregation expression over a group of rows.
// Arithmetic over aggregation calls is supported per the documented grammar.
func evalAggExpr(e exprNode, rows []queryResultRow, ctx *execContext) interface{} {
	switch n := e.(type) {
	case *funcNode:
		return computeAggregation(strings.ToLower(n.name), n.args, rows, ctx)
	case *binOpNode:
		l := evalAggExpr(n.l, rows, ctx)
		r := evalAggExpr(n.r, rows, ctx)
		ln, lok := asNumber(l)
		rn, rok := asNumber(r)
		if !lok || !rok {
			return nil
		}
		return applyArith(n.op, ln, rn)
	case *literalNode:
		return n.val
	}
	return nil
}

// aggArgs carries one aggregation call's argument expressions, the
// group's rows and the execution context; the column views are the
// interface every aggregation family reads its input through.
type aggArgs struct {
	args []exprNode
	rows []queryResultRow
	ctx  *execContext
}

func (a aggArgs) argValues() []interface{} {
	out := make([]interface{}, len(a.rows))
	for i := range a.rows {
		row := a.rows[i]
		if len(a.args) == 0 {
			out[i] = nil
			continue
		}
		out[i] = a.args[0].eval(&row, a.ctx)
	}
	return out
}

func (a aggArgs) numeric() []float64 {
	var out []float64
	for _, v := range a.argValues() {
		if f, ok := asNumber(v); ok {
			out = append(out, f)
		}
	}
	return out
}

// present collects the rows' argument values that exist. count(fieldName)
// "counts all records that include the specified field name" — inclusion
// is existence (a null value is the missing field), not non-emptiness:
// an empty-string value is a present value (stats function table).
func (a aggArgs) present() []interface{} {
	var out []interface{}
	for _, v := range a.argValues() {
		if v != nil {
			out = append(out, v)
		}
	}
	return out
}

// computeAggregation evaluates one aggregation function over the rows of a
// group: the family dispatchers decide the name in turn.
func computeAggregation(name string, args []exprNode, rows []queryResultRow, ctx *execContext) interface{} {
	a := aggArgs{args: args, rows: rows, ctx: ctx}
	if v, ok := evalBasicAggregations(name, a); ok {
		return v
	}
	if v, ok := evalTimeSeriesAggregations(name, a); ok {
		return v
	}
	return nil
}

// evalBasicAggregations holds the plain aggregation family's dispatch —
// counting, arithmetic, distinctness, statistics, ordering — false means
// the name is not this family's.
func evalBasicAggregations(name string, a aggArgs) (interface{}, bool) {
	switch name {
	case "count":
		if len(a.args) == 0 {
			return float64(len(a.rows)), true
		}
		return float64(len(a.present())), true
	case "avg":
		vals := a.numeric()
		if len(vals) == 0 {
			return nil, true
		}
		sum := 0.0
		for _, v := range vals {
			sum += v
		}
		return sum / float64(len(vals)), true
	case "sum":
		vals := a.numeric()
		if len(vals) == 0 {
			return float64(0), true
		}
		sum := 0.0
		for _, v := range vals {
			sum += v
		}
		return sum, true
	case "min":
		vals := a.argValues()
		if len(vals) == 0 {
			return nil, true
		}
		best := vals[0]
		for _, v := range vals[1:] {
			if compareForSort(v, best) < 0 {
				best = v
			}
		}
		return wrapTimestampResult(a.args, best), true
	case "max":
		vals := a.argValues()
		if len(vals) == 0 {
			return nil, true
		}
		best := vals[0]
		for _, v := range vals[1:] {
			if compareForSort(v, best) > 0 {
				best = v
			}
		}
		return wrapTimestampResult(a.args, best), true
	case "countdistinct":
		seen := make(map[string]bool)
		for _, v := range a.present() {
			seen[asString(v)] = true
		}
		return float64(len(seen)), true
	case "stddev":
		vals := a.numeric()
		if len(vals) == 0 {
			return nil, true
		}
		mean := 0.0
		for _, v := range vals {
			mean += v
		}
		mean /= float64(len(vals))
		variance := 0.0
		for _, v := range vals {
			variance += (v - mean) * (v - mean)
		}
		variance /= float64(len(vals))
		return math.Sqrt(variance), true
	case "variance":
		vals := a.numeric()
		if len(vals) == 0 {
			return nil, true
		}
		mean := 0.0
		for _, v := range vals {
			mean += v
		}
		mean /= float64(len(vals))
		variance := 0.0
		for _, v := range vals {
			variance += (v - mean) * (v - mean)
		}
		variance /= float64(len(vals))
		return variance, true
	case "pct":
		vals := a.numeric()
		if len(vals) == 0 || len(a.args) < 2 {
			return nil, true
		}
		p, ok := asNumber(a.args[1].eval(&a.rows[0], a.ctx))
		if !ok {
			return nil, true
		}
		return wrapTimestampResult(a.args, percentile(vals, p)), true
	case "values", "collect_values":
		seen := make(map[string]bool)
		var out []interface{}
		for _, v := range a.present() {
			s := asString(v)
			if !seen[s] {
				seen[s] = true
				out = append(out, s)
			}
		}
		if out == nil {
			out = []interface{}{}
		}
		return out, true
	case "topk":
		if len(a.args) < 2 || len(a.rows) == 0 {
			return nil, true
		}
		k, ok := asNumber(a.args[0].eval(&a.rows[0], a.ctx))
		if !ok || k < 1 || k > logsstore.MaxTopkK {
			return nil, true
		}
		// topk's field is the SECOND argument ("topk(k: number,
		// fieldName: LogField)"), unlike every other aggregation whose
		// field is the first — counting args[0] here counted the literal
		// k once per row.
		counts := make(map[string]int)
		for i := range a.rows {
			if v := a.args[1].eval(&a.rows[i], a.ctx); v != nil {
				counts[asString(v)]++
			}
		}
		type kv struct {
			k string
			v int
		}
		var list []kv
		for k, v := range counts {
			list = append(list, kv{k, v})
		}
		sort.Slice(list, func(i, j int) bool {
			if list[i].v != list[j].v {
				return list[i].v > list[j].v
			}
			return list[i].k < list[j].k
		})
		limit := int(k)
		if limit > len(list) {
			limit = len(list)
		}
		var out []interface{}
		for i := 0; i < limit; i++ {
			out = append(out, list[i].k)
		}
		return out, true
	case "earliest", "latest":
		bestRow := -1
		var bestTS float64
		for i := range a.rows {
			row := a.rows[i]
			ts, ok := asNumber(row.fields["@timestamp"])
			if !ok {
				continue
			}
			if bestRow < 0 ||
				(name == "earliest" && ts < bestTS) ||
				(name == "latest" && ts > bestTS) {
				bestRow = i
				bestTS = ts
			}
		}
		if bestRow < 0 || len(a.args) == 0 {
			return nil, true
		}
		row := a.rows[bestRow]
		return wrapTimestampResult(a.args, a.args[0].eval(&row, a.ctx)), true
	case "sortsfirst", "sortslast":
		vals := a.present()
		if len(vals) == 0 {
			return nil, true
		}
		best := vals[0]
		for _, v := range vals[1:] {
			cmp := compareForSort(v, best)
			if (name == "sortsfirst" && cmp < 0) || (name == "sortslast" && cmp > 0) {
				best = v
			}
		}
		return wrapTimestampResult(a.args, best), true
	}
	return nil, false
}

// isTimestampExpr reports whether the expression yields a timestamp-typed
// value: a direct event timestamp field, a datetime function whose
// documented result type is Timestamp, or the passthrough aggregation
// functions (documented result type LogField) applied to one of those.
func isTimestampExpr(e exprNode) bool {
	switch n := e.(type) {
	case *fieldNode:
		if len(n.segs) == 1 {
			switch strings.ToLower(n.segs[0].name) {
			case "@timestamp", "@ingestiontime":
				return true
			}
		}
	case *funcNode:
		switch strings.ToLower(n.name) {
		case "frommillis", "datefloor", "dateceil":
			return true
		case "earliest", "latest", "min", "max", "pct", "sortsfirst", "sortslast":
			return len(n.args) > 0 && isTimestampExpr(n.args[0])
		}
	}
	return false
}

// wrapTimestampResult renders a passthrough aggregation result as a
// timestamp when its argument is timestamp-typed, so the value keeps the
// documented Timestamp rendering in result rows while remaining numeric
// inside further expressions.
func wrapTimestampResult(args []exprNode, v interface{}) interface{} {
	if _, ok := v.(timestampValue); ok {
		return v
	}
	if len(args) == 0 || !isTimestampExpr(args[0]) {
		return v
	}
	if ms, ok := asNumber(v); ok {
		return timestampValue(int64(ms))
	}
	return v
}
