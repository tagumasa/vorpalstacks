package cloudwatchlogs

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// The fillmissing command: constant fills for the bin intervals a stats
// by bin() grouping skipped, completing the time series.

// --- fillmissing ---

type fillmissingCommand struct {
	fills  []fillConst
	headTk token
}

type fillConst struct {
	value string
	field string
}

func (c *fillmissingCommand) name() string { return "fillmissing" }

func parseFillmissingCommand(args []token, head token) (command, error) {
	c := &fillmissingCommand{headTk: head}
	i := 0
	for i < len(args) {
		t := args[i]
		if !(t.kind == tokIdent && strings.EqualFold(t.text, "with")) {
			return nil, newQueryCompileError(
				fmt.Sprintf("Syntax error at '%s': fillmissing expects with <value> for <field>", t.raw), t.start, t.end)
		}
		if i+4 > len(args) {
			return nil, newQueryCompileError("Syntax error: incomplete with clause", head.start, head.end)
		}
		val := args[i+1]
		if val.kind != tokIdent && val.kind != tokNumber && val.kind != tokString {
			return nil, newQueryCompileError("Syntax error: with requires a value", val.start, val.end)
		}
		if args[i+2].kind != tokIdent || !strings.EqualFold(args[i+2].text, "for") {
			return nil, newQueryCompileError("Syntax error: expected for", args[i+2].start, args[i+2].end)
		}
		fld := args[i+3]
		if fld.kind != tokIdent && fld.kind != tokBacktickIdent {
			return nil, newQueryCompileError("Syntax error: for requires a field", fld.start, fld.end)
		}
		c.fills = append(c.fills, fillConst{value: val.text, field: fld.text})
		i += 4
		if i < len(args) && args[i].kind == tokComma {
			i++
		}
	}
	return c, nil
}

func (c *fillmissingCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	info := ctx.lastBins
	if info == nil || info.dur <= 0 {
		return rows
	}
	existing := make(map[string]bool, len(rows))
	for _, r := range rows {
		existing[r.fields["@bin"]] = true
	}
	var out []queryResultRow
	out = append(out, rows...)
	for bin := info.minBin; bin <= info.maxBin; bin += info.dur {
		key := formatResultTimestamp(strconv.FormatInt(bin, 10))
		if existing[key] {
			continue
		}
		synth := queryResultRow{}
		synth.set("@bin", key)
		for _, f := range c.fills {
			synth.set(f.field, f.value)
		}
		out = append(out, synth)
	}
	// Bins are stored in their timestamp rendering; the fixed-width layout
	// orders chronologically, with the numeric form kept as a fallback for
	// defensive sorting of unformatted values.
	binOrder := func(v string) int64 {
		if ms, ok := parseResultTimestamp(v); ok {
			return ms
		}
		ms, _ := asNumber(v)
		return int64(ms)
	}
	sort.SliceStable(out, func(i, j int) bool {
		return binOrder(out[i].fields["@bin"]) < binOrder(out[j].fields["@bin"])
	})
	return out
}
