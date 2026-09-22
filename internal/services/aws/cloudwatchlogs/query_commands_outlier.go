package cloudwatchlogs

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// The statistical and sequence commands: outlier (with the percentile
// helper it shares) and sessionize — each classifies or segments the
// pipeline's current row set.

// --- outlier ---

type outlierCommand struct {
	action   string // remove or transform
	param    float64
	useLower bool
	mark     bool
	fields   []string
	head     token
}

func (c *outlierCommand) name() string { return "outlier" }

func parseOutlierCommand(args []token, head token) (command, error) {
	c := &outlierCommand{head: head, action: "transform", param: 2.5}
	for i := 0; i < len(args); i++ {
		t := args[i]
		if t.kind != tokIdent {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
		lower := strings.ToLower(t.text)
		isKV := i+2 < len(args) && args[i+1].kind == tokOp && args[i+1].text == "="
		switch {
		case lower == "action" && isKV:
			c.action = strings.ToLower(args[i+2].text)
			i += 2
		case lower == "param" && isKV:
			f, err := strconv.ParseFloat(args[i+2].text, 64)
			if err != nil {
				return nil, newQueryCompileError(fmt.Sprintf("Invalid param '%s'", args[i+2].text), args[i+2].start, args[i+2].end)
			}
			c.param = f
			i += 2
		case lower == "uselower" && isKV:
			c.useLower = strings.EqualFold(args[i+2].text, "true")
			i += 2
		case lower == "mark" && isKV:
			c.mark = strings.EqualFold(args[i+2].text, "true")
			i += 2
		default:
			c.fields = append(c.fields, t.text)
		}
	}
	if c.action != "remove" && c.action != "transform" {
		return nil, newQueryCompileError(fmt.Sprintf("Invalid outlier action '%s'", c.action), head.start, head.end)
	}
	if len(c.fields) == 0 {
		return nil, newQueryCompileError("Syntax error: outlier requires at least one field", head.start, head.end)
	}
	return c, nil
}

func (c *outlierCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	bounds := make(map[string][2]float64, len(c.fields))
	for _, f := range c.fields {
		var vals []float64
		for i := range rows {
			if v, ok := asNumber(readField(rows[i], f, ctx)); ok {
				vals = append(vals, v)
			}
		}
		if len(vals) == 0 {
			continue
		}
		q1 := percentile(vals, 25)
		q3 := percentile(vals, 75)
		iqr := q3 - q1
		bounds[f] = [2]float64{q1 - c.param*iqr, q3 + c.param*iqr}
	}
	out := make([]queryResultRow, 0, len(rows))
	for i := range rows {
		row := &rows[i]
		isOutlier := false
		for _, f := range c.fields {
			b, ok := bounds[f]
			if !ok {
				continue
			}
			v, ok := asNumber(readField(*row, f, ctx))
			if !ok {
				continue
			}
			high := v > b[1]
			low := c.useLower && v < b[0]
			if c.mark {
				row.set(f+"_outlier", strconv.FormatBool(high || low))
			}
			if high || low {
				isOutlier = true
				if c.action == "transform" {
					if high {
						row.set(f, formatNumber(b[1]))
					} else {
						row.set(f, formatNumber(b[0]))
					}
					isOutlier = false
				}
			}
		}
		if c.action == "remove" && isOutlier {
			continue
		}
		out = append(out, *row)
	}
	return out
}

// percentile computes the nearest-rank percentile of a value slice. The
// slice is sorted in place.
func percentile(vals []float64, p float64) float64 {
	sort.Float64s(vals)
	if len(vals) == 1 {
		return vals[0]
	}
	rank := int(math.Ceil(p / 100 * float64(len(vals))))
	if rank < 1 {
		rank = 1
	}
	if rank > len(vals) {
		rank = len(vals)
	}
	return vals[rank-1]
}

// --- sessionize ---

type sessionizeCommand struct {
	fields  []string
	maxSpan int64
	out     string
	head    token
}

func (c *sessionizeCommand) name() string { return "sessionize" }

func parseSessionizeCommand(args []token, head token) (command, error) {
	c := &sessionizeCommand{head: head, maxSpan: 30 * 60 * 1000, out: "session_id"}
	for i := 0; i < len(args); i++ {
		t := args[i]
		if t.kind != tokIdent {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
		lower := strings.ToLower(t.text)
		switch {
		case lower == "maxspan" && i+1 < len(args):
			text, span, ok := periodArgAt(args, i+1)
			if !ok {
				return nil, newQueryCompileError(fmt.Sprintf("Invalid maxspan '%s'", args[i+1].text), args[i+1].start, args[i+1].end)
			}
			ms, ok2 := parsePeriodMillis(text)
			if !ok2 {
				return nil, newQueryCompileError(fmt.Sprintf("Invalid maxspan '%s'", text), args[i+1].start, args[i+1].end)
			}
			c.maxSpan = ms
			i += span
		case lower == "as":
			if i+1 >= len(args) {
				return nil, newQueryCompileError("Syntax error: as requires a name", head.start, head.end)
			}
			c.out = args[i+1].text
			i++
		default:
			c.fields = append(c.fields, t.text)
		}
	}
	if len(c.fields) == 0 {
		return nil, newQueryCompileError("Syntax error: sessionize requires identity fields", head.start, head.end)
	}
	return c, nil
}

func (c *sessionizeCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	// Sessions are formed per identity value combination: events ordered by
	// @timestamp start a new session when the inactivity gap exceeds maxspan.
	type sessionState struct {
		lastTS int64
		id     int
	}
	sessions := make(map[string]*sessionState)
	counter := 0
	idx := make([]int, len(rows))
	for i := range idx {
		idx[i] = i
	}
	sort.SliceStable(idx, func(a, b int) bool {
		va, _ := asNumber(rows[idx[a]].fields["@timestamp"])
		vb, _ := asNumber(rows[idx[b]].fields["@timestamp"])
		return va < vb
	})
	for _, i := range idx {
		var key strings.Builder
		for _, f := range c.fields {
			key.WriteString(readField(rows[i], f, ctx))
			key.WriteByte(0)
		}
		ts, _ := asNumber(rows[i].fields["@timestamp"])
		st, ok := sessions[key.String()]
		if !ok || int64(ts)-st.lastTS > c.maxSpan {
			counter++
			st = &sessionState{id: counter}
			sessions[key.String()] = st
		}
		st.lastTS = int64(ts)
		rows[i].set(c.out, strconv.Itoa(st.id))
	}
	return rows
}
